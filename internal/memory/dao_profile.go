package memory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/CarlosShao/wisp/internal/observe"
)

// L1 profile DAO (D20/D35). The <=20-row cap is enforced HERE (application
// layer), with LRU-by-updated_at eviction that writes a log line — eviction
// must be visible, otherwise the user experiences "it forgot me" with no
// explanation (D35 rule on profile).

// UpsertProfile inserts or refreshes a profile slot and enforces the row cap,
// evicting the least-recently-updated rows (with a log line each) when the
// cap is exceeded. Returns the row id.
func (s *Store) UpsertProfile(ctx context.Context, slot, content string, source ProfileSource) (int64, error) {
	if slot == "" {
		return 0, fmt.Errorf("memory: profile.slot is required")
	}
	if content == "" {
		return 0, fmt.Errorf("memory: profile.content is required")
	}
	if !source.Valid() {
		return 0, fmt.Errorf("memory: invalid profile source %q (want extracted|manual)", source)
	}
	now := observe.NowWallUTC().Unix()

	var id int64
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx,
			`INSERT INTO profile(slot, content, source, updated_at) VALUES(?,?,?,?)
			 ON CONFLICT(slot) DO UPDATE SET content=excluded.content, source=excluded.source, updated_at=excluded.updated_at`,
			slot, content, string(source), now)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		if err != nil {
			return err
		}
		// Row cap: evict LRU rows (never the row just written — it has the
		// newest updated_at) until the cap holds.
		if err := evictProfileLRU(ctx, tx, s.logger, slot); err != nil {
			return err
		}
		// slot is UNIQUE: on conflict LastInsertId is unreliable; fetch it.
		if err := tx.QueryRowContext(ctx,
			`SELECT id FROM profile WHERE slot=?`, slot).Scan(&id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("memory: upsert profile %q: %w", slot, err)
	}
	return id, nil
}

// evictProfileLRU deletes the oldest-updated rows while the table holds more
// than L1MaxProfiles rows. Every eviction logs a warning line (D35: the
// eviction must be visible).
func evictProfileLRU(ctx context.Context, tx *sql.Tx, logger *slog.Logger, keepSlot string) error {
	for {
		var over int
		if err := tx.QueryRowContext(ctx, `SELECT count(*) - ? FROM profile`, L1MaxProfiles).Scan(&over); err != nil {
			return err
		}
		if over <= 0 {
			return nil
		}
		var (
			evictID   int64
			evictSlot string
			updatedAt int64
		)
		// Oldest updated_at first; skip the row just written (it is the
		// newest, but guard anyway when several rows share a timestamp).
		err := tx.QueryRowContext(ctx,
			`SELECT id, slot, updated_at FROM profile
			 WHERE slot != ?
			 ORDER BY updated_at ASC, id ASC
			 LIMIT 1`, keepSlot).Scan(&evictID, &evictSlot, &updatedAt)
		if errors.Is(err, sql.ErrNoRows) {
			// Only the just-written row remains beyond the cap: the cap
			// config is larger than the table can hold; nothing to evict.
			return nil
		}
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM profile WHERE id=?`, evictID); err != nil {
			return err
		}
		logger.Warn("profile LRU eviction",
			"component", "memory",
			"evicted_slot", evictSlot,
			"evicted_id", evictID,
			"evicted_updated_at", updatedAt,
			"reason", fmt.Sprintf("L1 profile cap reached (%d rows)", L1MaxProfiles))
	}
}

// ListProfiles returns all L1 rows, most recently updated first.
func (s *Store) ListProfiles(ctx context.Context) ([]Profile, error) {
	rows, err := s.reader.QueryContext(ctx,
		`SELECT id, slot, content, source, updated_at FROM profile ORDER BY updated_at DESC, id DESC`)
	if err != nil {
		return nil, fmt.Errorf("memory: list profiles: %w", err)
	}
	defer rows.Close()
	return scanProfiles(rows)
}

// ProfileBySlot returns the profile row for slot (ErrNotFound when absent).
func (s *Store) ProfileBySlot(ctx context.Context, slot string) (Profile, error) {
	var p Profile
	err := s.reader.QueryRowContext(ctx,
		`SELECT id, slot, content, source, updated_at FROM profile WHERE slot=?`, slot).
		Scan(&p.ID, &p.Slot, &p.Content, &p.Source, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Profile{}, fmt.Errorf("memory: profile slot %q: %w", slot, ErrNotFound)
	}
	if err != nil {
		return Profile{}, fmt.Errorf("memory: profile by slot: %w", err)
	}
	return p, nil
}

// DeleteProfile removes one profile row by id. Returns ErrNotFound when the
// id does not exist.
func (s *Store) DeleteProfile(ctx context.Context, id int64) error {
	return s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `DELETE FROM profile WHERE id=?`, id)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return fmt.Errorf("profile %d: %w", id, ErrNotFound)
		}
		return nil
	})
}

// PurgeProfiles removes every L1 row (privacy: one-click clear, SPEC-02 §5).
func (s *Store) PurgeProfiles(ctx context.Context) (int64, error) {
	var n int64
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `DELETE FROM profile`)
		if err != nil {
			return err
		}
		n, _ = res.RowsAffected()
		return nil
	})
	return n, err
}

func scanProfiles(rows *sql.Rows) ([]Profile, error) {
	var out []Profile
	for rows.Next() {
		var p Profile
		if err := rows.Scan(&p.ID, &p.Slot, &p.Content, &p.Source, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ErrNotFound is returned by delete/get-one operations when the target row
// (or artifact file) does not exist.
var ErrNotFound = errors.New("memory: not found")
