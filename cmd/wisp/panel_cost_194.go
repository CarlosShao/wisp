package main

// Ticket 194's cost door gets its real leg here (ledger A1681's AC#5 boundary:
// the internal/panel batch landed the door with a nil socket that refused by
// name; this file is the cmd/wisp wave that attaches the reader).
//
// WHAT IT IS. panel.CostReader is the injected cost leg (cost_handlers.go: one
// read method, nothing else, so the door is read-only by construction). This
// adapter answers it from memory.Store.ListCostDaily (dao_misc.go) - the same
// C23 cost_daily table the task loop's BumpCostDay writes - so cost.summary
// returns the machine's real per-day rows instead of the named refusal. The
// mapping is one struct field per column, integers only; the adapter derives
// nothing, adds nothing, and reorders nothing: ListCostDaily already returns
// newest-day-first and panel's handler repeats that contract in its prose.
//
// WHY OPEN PER READ. The same shape every read outlet in this package already
// uses (sessions.go's list/search/outlets: memory.Open, one query, defer
// Close). The cost door is a user-triggered, rare read - not a hot path - and
// a handle held for the process's whole life would pin wisp.db open for a
// session that may never ask again, while a failed open would either fail the
// whole chain's assembly for a door the request never used or have to be
// cached as a permanent refusal. Per-read keeps every failure on the request
// that actually hit it: a database that will not open refuses THE COST REQUEST
// - loudly, through the handler's read-failed audit branch - and the next
// request starts clean.

import (
	"context"
	"fmt"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/panel"
)

// costReader194 adapts memory.Store.ListCostDaily into panel.CostReader. It
// holds the data root only; no database handle lives between reads.
type costReader194 struct {
	dir string
}

// newCostReader194 builds the adapter over the data root the chain was given.
// No database is opened here - see the file comment.
func newCostReader194(dir string) *costReader194 {
	return &costReader194{dir: dir}
}

// CostDays answers panel.CostReader: the real cost_daily rows, newest day
// first, one view row per stored row. An empty table answers an empty slice -
// zero rows is a true answer for a machine that has not run a task, not an
// error and not a stub.
func (r *costReader194) CostDays(ctx context.Context) ([]panel.CostDayRow, error) {
	st, err := memory.Open(r.dir)
	if err != nil {
		return nil, fmt.Errorf("成本存储不可用（%s）：%w", r.dir, err)
	}
	defer func() { _ = st.Close() }()
	days, err := st.ListCostDaily(ctx)
	if err != nil {
		return nil, fmt.Errorf("成本读取失败（%s）：%w", r.dir, err)
	}
	out := make([]panel.CostDayRow, 0, len(days))
	for _, d := range days {
		out = append(out, panel.CostDayRow{
			Day:        d.Day,
			TokensIn:   d.TokensIn,
			TokensOut:  d.TokensOut,
			CostMicros: d.CostMicros,
			Tasks:      d.Tasks,
		})
	}
	return out, nil
}
