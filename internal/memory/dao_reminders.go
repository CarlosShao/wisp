package memory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/CarlosShao/wisp/internal/observe"
)

// reminders DAO (ticket 31, schema v6): one-shot in-process reminders.
// Lifecycle policy (who arms the timer, what fire does) belongs to
// internal/tools + internal/agent; this layer only stores rows and enforces
// the status roster at the boundary (same shape as taskLogStateRoster,
// ticket 196).

// reminderStatusRoster is the authority for reminders.status. The v6 DDL
// CHECK is GENERATED from it (buildDDLV6 in schema.go), so the SQL and the DAO
// validator cannot fork — the same rule ticket 196 set for task_log.state.
var reminderStatusRoster = []string{
	ReminderArmed,
	ReminderFired,
	ReminderCancelled,
}

func isReminderStatus(s string) bool {
	for _, v := range reminderStatusRoster {
		if v == s {
			return true
		}
	}
	return false
}

// Status transitions are one-way and terminal: armed -> fired | cancelled.
// MarkFired/Cancel only UPDATE rows still in armed, so a late timer cannot
// overwrite a cancel, and a double-fire cannot re-stamp fired_utc.

// InsertReminder stores one armed reminder row.
func (s *Store) InsertReminder(ctx context.Context, r Reminder) error {
	if r.ID == "" || r.Note == "" {
		return fmt.Errorf("memory: reminder id/note are required")
	}
	if r.FireAt <= 0 {
		return fmt.Errorf("memory: reminder fire_at_wallclock_utc must be a positive unix-seconds value")
	}
	if r.Status == "" {
		r.Status = ReminderArmed
	}
	if !isReminderStatus(r.Status) {
		return fmt.Errorf("memory: invalid reminder status %q (want armed|fired|cancelled)", r.Status)
	}
	if r.CreatedAt == 0 {
		r.CreatedAt = observe.NowWallUTC().Unix()
	}
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO reminders(id, fire_at_wallclock_utc, note, status, created_utc, fired_utc)
			 VALUES(?,?,?,?,?,?)`,
			r.ID, r.FireAt, r.Note, r.Status, r.CreatedAt, nullableInt64(r.FiredAt))
		return err
	})
	if err != nil {
		return fmt.Errorf("memory: insert reminder %s: %w", r.ID, err)
	}
	return nil
}

// ReminderByID returns one reminder row.
func (s *Store) ReminderByID(ctx context.Context, id string) (Reminder, error) {
	var r Reminder
	var fired sql.NullInt64
	err := s.reader.QueryRowContext(ctx,
		`SELECT id, fire_at_wallclock_utc, note, status, created_utc, fired_utc
		 FROM reminders WHERE id=?`, id).
		Scan(&r.ID, &r.FireAt, &r.Note, &r.Status, &r.CreatedAt, &fired)
	if errors.Is(err, sql.ErrNoRows) {
		return Reminder{}, fmt.Errorf("reminders %s: %w", id, ErrNotFound)
	}
	if err != nil {
		return Reminder{}, fmt.Errorf("memory: reminder by id: %w", err)
	}
	r.FiredAt = nullInt64Ptr(fired)
	return r, nil
}

// ListReminders returns reminders in fire-at order (all statuses; the tool
// layer filters what the user sees).
func (s *Store) ListReminders(ctx context.Context) ([]Reminder, error) {
	rows, err := s.reader.QueryContext(ctx,
		`SELECT id, fire_at_wallclock_utc, note, status, created_utc, fired_utc
		 FROM reminders ORDER BY fire_at_wallclock_utc, id`)
	if err != nil {
		return nil, fmt.Errorf("memory: list reminders: %w", err)
	}
	defer rows.Close()
	var out []Reminder
	for rows.Next() {
		var r Reminder
		var fired sql.NullInt64
		if err := rows.Scan(&r.ID, &r.FireAt, &r.Note, &r.Status, &r.CreatedAt, &fired); err != nil {
			return nil, err
		}
		r.FiredAt = nullInt64Ptr(fired)
		out = append(out, r)
	}
	return out, rows.Err()
}

// MarkReminderFired stamps the one-way armed -> fired transition. It affects
// zero rows when the reminder was already cancelled or already fired (the
// late-timer guard); callers treat RowsAffected==0 as "do not deliver".
func (s *Store) MarkReminderFired(ctx context.Context, id string) (bool, error) {
	now := observe.NowWallUTC().Unix()
	var applied bool
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx,
			`UPDATE reminders SET status='fired', fired_utc=? WHERE id=? AND status='armed'`, now, id)
		if err != nil {
			return err
		}
		n, _ := res.RowsAffected()
		applied = n > 0
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("memory: mark reminder fired %s: %w", id, err)
	}
	return applied, nil
}

// CancelReminder stamps the one-way armed -> cancelled transition.
func (s *Store) CancelReminder(ctx context.Context, id string) (bool, error) {
	var applied bool
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx,
			`UPDATE reminders SET status='cancelled' WHERE id=? AND status='armed'`, id)
		if err != nil {
			return err
		}
		n, _ := res.RowsAffected()
		applied = n > 0
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("memory: cancel reminder %s: %w", id, err)
	}
	return applied, nil
}

// CountMissedArmed returns how many reminders are armed with a wall-clock
// fire time already in the past at `now`. Ticket 31 rule: boot reports this
// count and nothing else - the rows STAY armed in the table (the session's
// reminder tool arms them again for the future), and no missed reminder is
// ever re-delivered (跨重启不补发).
func (s *Store) CountMissedArmed(ctx context.Context, now int64) (int, error) {
	var n int
	err := s.reader.QueryRowContext(ctx,
		`SELECT count(*) FROM reminders WHERE status='armed' AND fire_at_wallclock_utc <= ?`, now).
		Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("memory: count missed armed reminders: %w", err)
	}
	return n, nil
}
