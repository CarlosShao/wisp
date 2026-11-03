package agent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Ticket 44b (落地腿 44b-r1): the per-model consumption permanence on the
// quota side.
//
// What 44-r4 left named in observe/quota_44.go's persistence note ("the
// per-model buckets are session-scoped until the owner-approved schema step"):
// 44b IS that step (D35 addition, ledger A1445). The table and DAO live in
// internal/memory (cost_daily_model, migration 4->5); the same-batch write
// primitive and the ledger seed live in internal/observe (costmodel_44b.go);
// this file is the gate's use of both, so the pre-dispatch quota judgment
// reads the same numbers after a restart as it did before one (44b AC#3).
//
// Why the gate is the booking mouth for the model dimension: it is the one
// place "provider/model" is known with certainty (Book already keys by the
// served chain element, 票 44-r4). The loop's bookCostDay books the GLOBAL
// cost_daily row and knows no served key; a caller that has a served key
// calls BookModel instead of the pair (PostTask + nothing) — BookModel writes
// BOTH cost_daily and cost_daily_model in one storage transaction (memory's
// BumpCostDayModel), so "同批" is atomic at the disk, not a best-effort pair.
//
// No cmd/wisp change rides on this file: the wiring leg (the host calling
// BookModel at task end and SeedBaseline at boot) is the composition root's
// one-inch step, named in docs/evidence/s1/44br1-impl.md. Book itself keeps
// its 44-r4 signature: existing hosts compile and behave unchanged.

// SetModelStore points the gate at the per-model disk mouths. Nil ledger =
// the gate books no model rows (session-only, exactly 44-r4's behaviour);
// nil reader = the gate cannot rebuild baselines (fresh session each boot).
// The clock is injectable for the fixture rollover tests (same rule as
// observe.NewQuotaLedger); nil = time.Now.
func (g *QuotaGate) SetModelStore(l observe.ModelLedger, r observe.ModelReader, clock func() time.Time) {
	g.modelLedger = l
	g.models = r
	if clock != nil {
		g.clock = clock
	}
}

// BookModel posts one served task's rollup into cost_daily AND
// cost_daily_model (one storage batch), and into the session ledger through
// the same Book the gate has always used. The key is the "provider/model"
// chain spelling. Tokens-only (unknown price) arrives as micros == 0 with
// real token counts — the marker is the reader's, never a fabricated number.
//
// Failures split honestly: a disk refusal never mutates the session ledger
// (the disk write is attempted first; on error Book is not called, so the
// gate's remaining and the disk cannot drift by one task in either direction
// without the other half saying so).
func (g *QuotaGate) BookModel(ctx context.Context, key string, tokensIn, tokensOut, micros int64) error {
	if g.modelLedger == nil {
		return errors.New("agent: quota gate has no model ledger (SetModelStore was never given one)")
	}
	provider, model := splitQuotaKey(key)
	if provider == "" || model == "" {
		return fmt.Errorf("agent: BookModel key %q is not \"provider/model\"", key)
	}
	if _, err := g.bookModelDay(ctx, provider, model, tokensIn, tokensOut, micros); err != nil {
		return err
	}
	return g.Book(key, micros)
}

// bookModelDay is the ledger call factored out so the day label is computed
// once from the gate's clock (the same DayKey shape memory validates).
func (g *QuotaGate) bookModelDay(ctx context.Context, provider, model string, tokensIn, tokensOut, micros int64) (string, error) {
	day := observe.DayKey(g.now())
	if err := g.modelLedger.BumpCostDayModel(ctx, day, provider, model, tokensIn, tokensOut, micros); err != nil {
		return day, err
	}
	return day, nil
}

// SeedBaseline rebuilds the per-model already-consumed numbers from
// cost_daily_model (44b AC#3): the month's rows fold into the ledger's day and
// month buckets, plan-billed models also draw the provider credit bucket
// (billing comes from the gate's own cards — the same source JudgeModel uses,
// so seed and judgment cannot disagree on what a model bills as).
//
// It is additive: seed once at boot before any task books. Repeated calls
// would double-count, and the method says so by refusing to be invisible —
// it reports the rows it folded. Idempotent restarts come from the host
// calling it exactly once at composition, not from a second source of truth.
func (g *QuotaGate) SeedBaseline(ctx context.Context) (int, error) {
	if g.models == nil {
		return 0, errors.New("agent: quota gate has no model reader (SetModelStore was never given one)")
	}
	now := g.now()
	day := observe.DayKey(now)
	rows, err := g.models.CostModelDaysBetween(ctx, observe.MonthFirst(day), day)
	if err != nil {
		return 0, err
	}
	if err := g.ledger.SeedFrom(rows, day, g.cards); err != nil {
		return 0, err
	}
	if g.logW != nil && len(rows) > 0 {
		g.logW("agent: quota baseline rebuilt from cost_daily_model", "rows", len(rows), "day", day)
	}
	return len(rows), nil
}

// now reads the gate's injected clock (default time.Now; SetModelStore can
// override it for fixture rollover tests).
func (g *QuotaGate) now() time.Time {
	if g.clock != nil {
		return g.clock()
	}
	return time.Now()
}
