package memory

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/CarlosShao/wisp/internal/observe"
)

// C23 read surface (ticket 44).
//
// dao_misc.go already carries the WRITE primitive (BumpCostDay: the upsert that
// adds one task's numbers to its day row) and the single-day/whole-table reads
// ticket 04 froze. What the CostMeter still needs is a range read it can hand
// the user as "this month", and the base-table sum the daily row is supposed to
// agree with. Both are queries over the same two frozen tables — no new column,
// no new table, no DDL.
//
// Range comparisons are string comparisons on purpose: cost_daily.day is
// 'YYYY-MM-DD' (validDay enforces the fixed width on every write), so lexical
// order IS chronological order and `day BETWEEN ? AND ?` needs no date
// arithmetic in SQL — which also keeps the day boundary in Go, where the
// meter's calendar lives, instead of forking it into SQLite's 'localtime'.

// CostDaysBetween returns the aggregates whose day label falls in
// [fromDay, toDay] inclusive, oldest first. An empty range is an empty result,
// not an error: "no cost yet today" is a fact a budget line has to render.
func (s *Store) CostDaysBetween(ctx context.Context, fromDay, toDay string) ([]observe.DayTotal, error) {
	if !validDay(fromDay) || !validDay(toDay) {
		return nil, fmt.Errorf("memory: cost_daily range needs YYYY-MM-DD endpoints (got %q..%q)", fromDay, toDay)
	}
	if fromDay > toDay {
		return nil, fmt.Errorf("memory: cost_daily range %s..%s is inverted", fromDay, toDay)
	}
	rows, err := s.reader.QueryContext(ctx,
		`SELECT day, tokens_in, tokens_out, cost_micros, tasks
		   FROM cost_daily WHERE day BETWEEN ? AND ? ORDER BY day`, fromDay, toDay)
	if err != nil {
		return nil, fmt.Errorf("memory: cost_daily range: %w", err)
	}
	defer rows.Close()
	var out []observe.DayTotal
	for rows.Next() {
		var d observe.DayTotal
		if err := rows.Scan(&d.Day, &d.TokensIn, &d.TokensOut, &d.Micros, &d.Tasks); err != nil {
			return nil, fmt.Errorf("memory: scan cost_daily range: %w", err)
		}
		out = append(out, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("memory: cost_daily range: %w", err)
	}
	return out, nil
}

// SumTasksEndedBetween sums the C23 columns of the task_log rows that ENDED in
// [fromUnix, toUnix) — the base table the daily aggregate is derived from, and
// the number CostMeter.Reconcile compares it against (ticket 44 AC#2).
//
// ended_at, not started_at: a task books its cost when it finishes (that is
// when the numbers exist), so the day it ends is the day its money lands in. A
// task that started yesterday and finished today belongs to today. Rows still
// running (ended_at NULL) are excluded rather than guessed at: they have no
// final cost, and an aggregate holding a guess is indistinguishable from one
// holding a real number.
func (s *Store) SumTasksEndedBetween(ctx context.Context, fromUnix, toUnix int64) (observe.DayTotal, error) {
	if toUnix <= fromUnix {
		return observe.DayTotal{}, fmt.Errorf("memory: task_log sum needs a non-empty range (got %d..%d)", fromUnix, toUnix)
	}
	var t observe.DayTotal
	err := s.reader.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(cost_tokens_in), 0),
		        COALESCE(SUM(cost_tokens_out), 0),
		        COALESCE(SUM(cost_amount_micro), 0),
		        COUNT(*)
		   FROM task_log WHERE ended_at >= ? AND ended_at < ?`, fromUnix, toUnix).
		Scan(&t.TokensIn, &t.TokensOut, &t.Micros, &t.Tasks)
	if err != nil && err != sql.ErrNoRows {
		return observe.DayTotal{}, fmt.Errorf("memory: sum task_log %d..%d: %w", fromUnix, toUnix, err)
	}
	return t, nil
}

// SetCostDay OVERWRITES one day's aggregate row. This is the repair mouth, not
// the hot path: dao_misc.go's BumpCostDay adds one task at a time (the shape the
// loop books through), and a recompute of a day whose rows are all in task_log
// is the only way to heal an aggregate that drifted from its base table —
// drift happens when a task_log row is closed by someone who never went through
// the loop's journal (cmd/wisp writes the row for a control word that never
// reached the loop, run.go's !opened branch).
//
// A day with no finished tasks writes a zero row rather than deleting the row:
// "today: 0" is a fact worth storing, and an absent row already means zero to
// every reader (CostDay returns a zero-value row when absent).
func (s *Store) SetCostDay(ctx context.Context, day string, tokensIn, tokensOut, costMicros, tasks int64) error {
	if !validDay(day) {
		return fmt.Errorf("memory: invalid cost_daily.day %q (want YYYY-MM-DD)", day)
	}
	if tokensIn < 0 || tokensOut < 0 || costMicros < 0 || tasks < 0 {
		return fmt.Errorf("memory: cost_daily %s: a negative aggregate is not an aggregate (%d/%d/%d/%d)",
			day, tokensIn, tokensOut, costMicros, tasks)
	}
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO cost_daily(day, tokens_in, tokens_out, cost_micros, tasks) VALUES(?,?,?,?,?)
			 ON CONFLICT(day) DO UPDATE SET
			   tokens_in=excluded.tokens_in,
			   tokens_out=excluded.tokens_out,
			   cost_micros=excluded.cost_micros,
			   tasks=excluded.tasks`,
			day, tokensIn, tokensOut, costMicros, tasks)
		return err
	})
	if err != nil {
		return fmt.Errorf("memory: set cost_daily %s: %w", day, err)
	}
	return nil
}

// Compile-time proof that the store is a CostMeter mouth (observe keeps the
// interfaces narrow; this is the line that says the narrow shape is satisfied).
var (
	_ observe.DayLedger  = (*Store)(nil)
	_ observe.DayReader  = (*Store)(nil)
	_ observe.TaskReader = (*Store)(nil)
	_ observe.DayWriter  = (*Store)(nil)
)
