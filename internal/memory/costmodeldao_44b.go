package memory

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Per-model daily aggregation DAO (ticket 44b, D35 addition owner-approved as
// ledger A1445; table created by the 4->5 migration step in schema.go).
//
// Same-batch rule (44b AC#2): BumpCostDayModel writes cost_daily AND
// cost_daily_model inside ONE storage transaction. A crash between the two
// INSERTs must never leave a day whose model rows disagree with its global row
// — the aggregate cache and its per-model split are only trustworthy when they
// move together, which is exactly what store.go's s.write gives (one tx).
//
// The cost_daily half mirrors dao_misc.go's BumpCostDay statement verbatim: one
// task, +tokens, +micros, +1 task. Callers that book through this method do
// NOT call BumpCostDay separately for the same task — that would count the
// task twice in the global row.

// BumpCostDayModel books one finished task into both aggregate tables in the
// same batch. Negative counters, a malformed day, or an empty provider/model
// are refused before anything is written (the same anti-negative shape the
// cost_daily DAO states; unknown price arrives as micros == 0 with real
// tokens — the tokens-only booking, marker is the reader's).
func (s *Store) BumpCostDayModel(ctx context.Context, day, provider, model string, tokensIn, tokensOut, costMicros int64) error {
	if !validDay(day) {
		return fmt.Errorf("memory: invalid cost_daily_model.day %q (want YYYY-MM-DD)", day)
	}
	if provider == "" || model == "" {
		return fmt.Errorf("memory: cost_daily_model %s: provider and model must both be named (got %q/%q)", day, provider, model)
	}
	if tokensIn < 0 || tokensOut < 0 || costMicros < 0 {
		return fmt.Errorf("memory: cost_daily_model %s %s/%s: a negative aggregate is not an aggregate (%d/%d/%d)",
			day, provider, model, tokensIn, tokensOut, costMicros)
	}
	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO cost_daily(day, tokens_in, tokens_out, cost_micros, tasks) VALUES(?,?,?,?,1)
			 ON CONFLICT(day) DO UPDATE SET
			   tokens_in=tokens_in+excluded.tokens_in,
			   tokens_out=tokens_out+excluded.tokens_out,
			   cost_micros=cost_micros+excluded.cost_micros,
			   tasks=tasks+1`,
			day, tokensIn, tokensOut, costMicros); err != nil {
			return err
		}
		_, err := tx.ExecContext(ctx,
			`INSERT INTO cost_daily_model(day, provider, model, tokens_in, tokens_out, amount_micro, tasks) VALUES(?,?,?,?,?,?,1)
			 ON CONFLICT(day, provider, model) DO UPDATE SET
			   tokens_in=tokens_in+excluded.tokens_in,
			   tokens_out=tokens_out+excluded.tokens_out,
			   amount_micro=amount_micro+excluded.amount_micro,
			   tasks=tasks+1`,
			day, provider, model, tokensIn, tokensOut, costMicros)
		return err
	})
	if err != nil {
		return fmt.Errorf("memory: bump cost_daily_model %s %s/%s: %w", day, provider, model, err)
	}
	return nil
}

// CostModelDaysBetween returns the per-model rows whose day label falls in
// [fromDay, toDay] inclusive, oldest first (day-first, then provider/model for
// a stable scan order). An empty range is an empty result, not an error:
// "no per-model cost yet" is the v5-day-zero fact the seed path must render.
func (s *Store) CostModelDaysBetween(ctx context.Context, fromDay, toDay string) ([]observe.CostModelDayTotal, error) {
	if !validDay(fromDay) || !validDay(toDay) {
		return nil, fmt.Errorf("memory: cost_daily_model range needs YYYY-MM-DD endpoints (got %q..%q)", fromDay, toDay)
	}
	if fromDay > toDay {
		return nil, fmt.Errorf("memory: cost_daily_model range %s..%s is inverted", fromDay, toDay)
	}
	rows, err := s.reader.QueryContext(ctx,
		`SELECT day, provider, model, tokens_in, tokens_out, amount_micro, tasks
		   FROM cost_daily_model WHERE day BETWEEN ? AND ? ORDER BY day, provider, model`, fromDay, toDay)
	if err != nil {
		return nil, fmt.Errorf("memory: cost_daily_model range: %w", err)
	}
	defer rows.Close()
	var out []observe.CostModelDayTotal
	for rows.Next() {
		var d observe.CostModelDayTotal
		if err := rows.Scan(&d.Day, &d.Provider, &d.Model, &d.TokensIn, &d.TokensOut, &d.Micros, &d.Tasks); err != nil {
			return nil, fmt.Errorf("memory: scan cost_daily_model range: %w", err)
		}
		out = append(out, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("memory: cost_daily_model range: %w", err)
	}
	return out, nil
}

// Compile-time proof that the store is the per-model mouth (same rule as
// costdao_44.go's block: observe keeps the interfaces narrow, this line says
// the narrow shape is satisfied).
var (
	_ observe.ModelLedger = (*Store)(nil)
	_ observe.ModelReader = (*Store)(nil)
)
