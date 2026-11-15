package main

// Ticket 451 - L1 profiles reach the prompt: the ONE producer side.
//
// THE GAP THIS CLOSES. Everything downstream of `agent.Config.Profiles` was
// already wired: internal/agent/loop.go's buildRequest forwards the field into
// `PromptInput.Profiles`, and internal/agent/prompt.go renders it as the D39 (3)
// profile section clipped to its scaled budget. What was missing was the
// producer: a repo-wide search for assignments to that field in production code
// found zero, so the L1 rows the extractor wrote were never read back for a
// prompt. Measured before this file existed: `agent.Config.Profiles` = 0
// assignments, and the AC#1 case failed with a request body that had no
// `## profile` section at all.
//
// WHO BUILDS THE CONFIG, WHO QUERIES. The composition root's execute()
// (cmd/wisp/run.go) is the single assembly point: it builds the per-task
// agent.Options/agent.Config for every task this process runs, including the
// snapshot task.spawn derives child loops from (rt.loopOpt). The read is
// loopProfileLines451 below, called from exactly that site against the same
// *memory.Store assembleRuntime opened (`memory.Open(s.dataDir)`, rt.store). It
// runs ONCE PER TASK, not once per boot: L1 extraction writes after a turn, so
// the next task must see the row the previous turn learned - a boot snapshot
// would have shown the model a memory one task stale forever.
//
// BUDGET (AC#3), THE D15 PATH. The lines are handed over as
// agent.Config.Profiles and take the path the assembler already owns:
//
//	agent.Config.Profiles
//	  -> PromptInput.Profiles            (internal/agent/loop.go buildRequest)
//	  -> profileText(..., SectionProfile) (internal/agent/prompt.go, clips to
//	                                       the scaled (3) section grant; 400
//	                                       tokens at the 128k reference
//	                                       window, D39)
//	  -> clipToTokens(..., budget)        (per-section hard clip)
//	  -> enforceTotal(..., PromptTotal)   (D39 total ceiling, suffix sections)
//
// Nothing is concatenated here: this file reads rows and renders exactly one
// line per row. The overflow policy is the budget's own, and it is nameable:
// ListProfiles returns newest-updated first (dao_profile.go), profileText walks
// that order and stops BEFORE the first line that would not fit, so the rows
// that fall off are the OLDEST ones (the tail of the list) - whole lines are
// dropped, never half a line. If enforceTotal has to shrink the section further
// (only when the prefix blocks alone squeeze the total ceiling), clipToTokens
// applies its cut marker to the whole section; that is pre-existing D39
// behaviour and this ticket does not touch it.
//
// A READ FAILURE IS DEGRADED, NOT FATAL, AND IT IS VISIBLE: a task must still
// run when the store read fails, but it must not silently pretend it had no
// memory. The audit line names the failure, and the prompt is built without the
// section (the same request shape as before this ticket, plus the line that
// says why).

import (
	"context"
	"strings"
	"time"
)

// profileReadBudget451 bounds one task's L1 read. It rides the monotonic-clock
// Timeout helper's shape (context.WithTimeout), never a wall-clock subtraction
// (D22 ban #4 / D42#9), and mirrors the other per-task store reads in execute()
// (5s).
const profileReadBudget451 = 5 * time.Second

// loopProfileLines451 renders the L1 profile rows as the prompt lines
// agent.Config.Profiles carries, newest-updated first (ListProfiles' order).
// The line shape is "slot: content": the slot enum ('pref.tone',
// 'habit.work_hours', 'fact.family', ...) is what tells the model which
// dimension the fact belongs to, and it is the column the user can delete one
// row by.
func (rt *agentRuntime) loopProfileLines451() []string {
	if rt == nil || rt.store == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), profileReadBudget451)
	defer cancel()
	rows, err := rt.store.ListProfiles(ctx)
	if err != nil {
		rt.auditf("wisp run: L1 画像读取失败，本轮 prompt 不带画像（不是没有画像）：%v", err)
		return nil
	}
	out := make([]string, 0, len(rows))
	for _, p := range rows {
		line := strings.TrimSpace(p.Slot + ": " + p.Content)
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}
