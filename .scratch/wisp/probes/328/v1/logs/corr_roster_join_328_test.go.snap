package agent

// Ticket 328 AC#2's positive control: the correlation id this case joins is READ
// off a dispatch the real loop made, never written in this file.
//
// The chain, end to end, with the mint point named at each hop:
//
//	newHarness("tool-then-text")  -> real golden SSE bytes through the genuine
//	openai-chat adapter (C5 seam, SPEC-05 §9 - not a function mock)
//	h.run(...)                    -> Loop.run -> executeCalls
//	                               -> loop.go:676 CorrelationID: callCorr(taskID,
//	                                  p.call.ID, i)   (the ONLY call site; :603 is
//	                                  the minter, :607 is the arm this replay takes)
//	h.echo().Calls()[0].Req       -> the ToolRequest the loop actually dispatched
//	panel.NativeVerdict{...}.CardView()
//	                               -> pump.go:98 -> approval.go:72, the same pure
//	                                  half the host's liveVerdicts uses, so the card
//	                                  here and the card on the wire cannot disagree
//	panel.TaskRosterSectionFrom  -> taskRosterSectionFrom in the same file, the
//	                                  render pump.go:333 calls (line numbers move
//	                                  with this ticket's own edit, so this chain
//	                                  names functions, not lines)
//
// Before ticket 328 the last hop answered false for exactly these bytes: the card
// carried "<taskID>#call_e1", the roster looked that string up as a task id, and no
// row was named that. The shapes this case cannot produce here (arm B's id-less
// call, and the queue's rewrite stacked on either arm) are paired in
// TestTheLoopMintsTheShapesTheRosterAccepts below and measured on the join side in
// internal/panel/subagent_corr_task_prefix_328_test.go.

import (
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/risk"
)

// rosterWithTwoRows files the loop's own task and one sibling row, so "exactly one
// row moved" is measurable rather than assumed. The Status half is the host's filed
// word carried verbatim by the render - this case asserts below that the join did
// not touch it (ticket 220's 外部对照 bound, inherited by 328).
func rosterWithTwoRows(taskID, otherID string) panel.TaskRosterState {
	return panel.TaskRosterState{
		Rows: []panel.TaskRow{
			{
				TaskID: taskID, Label: "本轮任务", Kind: "root",
				Status: "Thinking", StatusKnown: true, StreamKey: taskID,
			},
			{
				TaskID: otherID, ParentTaskID: taskID, Label: "另一枚任务",
				Kind: "subagent", Status: "Thinking", StatusKnown: true,
				StreamKey: panel.SubagentStreamKey(otherID),
			},
		},
		InFlightSlots: 1, PoolCap: 4,
	}
}

// cardForCorr is the host's own hop: the verdict fields are what the queue would
// carry for an out-of-workspace write (internal/risk R2), and the correlation id is
// the value read off the dispatch. No corr is invented here.
func cardForCorr(corr string) panel.ApprovalCardView {
	return panel.NativeVerdict{
		CorrelationID: corr, Tool: "fs.write", Args: []string{"C:/Windows/not-mine.txt"},
		Level: risk.L2, RulesHit: []risk.RuleID{risk.R2}, Reason: "写在工作区之外",
	}.CardView()
}

func rowOf(t *testing.T, sect *panel.TaskRosterSection, taskID string) panel.TaskRowView {
	t.Helper()
	for _, r := range sect.Rows {
		if r.TaskID == taskID {
			return r
		}
	}
	t.Fatalf("the render returned no row for %q: %+v", taskID, sect.Rows)
	return panel.TaskRowView{}
}

// TestACorrelationTheLoopActuallyMintedLightsItsRosterRow is AC#2's shot. The
// assertion that matters is the one on a correlation id this file did not write.
func TestACorrelationTheLoopActuallyMintedLightsItsRosterRow(t *testing.T) {
	h := newHarness(t, "tool-then-text")
	res := h.run("现在天气怎么样")
	if res.Status != StatusCompleted {
		t.Fatalf("fixture did not run: status=%s (%s)", res.Status, res.Message)
	}
	calls := h.echo().Calls()
	if len(calls) != 1 {
		t.Fatalf("tool dispatches = %d, want the one this replay makes", len(calls))
	}
	req := calls[0].Req
	if req.TaskID != res.TaskID {
		t.Fatalf("the dispatched request carries task %q, the loop reported %q", req.TaskID, res.TaskID)
	}
	// The bytes under test, read off the production dispatch. This is the guard
	// against the shape this ticket was filed over - a corr that happens to look
	// like the queue's numeric rewrite.
	corr := req.CorrelationID
	if corr == "" {
		t.Fatal("the loop dispatched a call with no correlation id")
	}
	if corr == req.TaskID || !strings.HasPrefix(corr, req.TaskID+"#") {
		t.Fatalf("corr %q is not the production shape <taskID>#<callID> for task %q: "+
			"the mint point this case is about moved", corr, req.TaskID)
	}
	if strings.TrimPrefix(corr, req.TaskID+"#") == "" ||
		strings.TrimPrefix(corr, req.TaskID+"#") != req.CallID {
		t.Fatalf("corr %q does not end in the call's own id %q", corr, req.CallID)
	}
	// And the minter is the one loop.go:676 calls: the same three inputs give the
	// same bytes. (Test285LoopDispatchesOneCorrPerCall owns the "two calls never
	// share a key" half; this line owns the "the roster is fed what the loop minted"
	// half.)
	if got := callCorr(req.TaskID, req.CallID, 0); got != corr {
		t.Fatalf("callCorr(%q,%q,0) = %q but the dispatched corr was %q", req.TaskID, req.CallID, got, corr)
	}

	sect := panel.TaskRosterSectionFrom(rosterWithTwoRows(res.TaskID, "sibling-of-"+res.TaskID),
		[]panel.ApprovalCardView{cardForCorr(corr)})
	row := rowOf(t, sect, res.TaskID)
	if !row.BlockedOnApproval {
		t.Errorf("the row for task %q reads not-blocked although a card of its own call (%q) is "+
			"pending: this is ticket 328's defect, measured on a correlation the loop minted", res.TaskID, corr)
	}
	other := rowOf(t, sect, "sibling-of-"+res.TaskID)
	if other.BlockedOnApproval {
		t.Errorf("the join lit a row no card names: %q blocked while only %q is pending",
			other.TaskID, corr)
	}
	if row.Status != "Thinking" || !row.StatusKnown {
		t.Errorf("blockedOnApproval took the running dimension with it: status=%q known=%v, want the "+
			"host's filed word carried untouched (AC#3: the D43 transition table is frozen)",
			row.Status, row.StatusKnown)
	}
}

// TestTheLoopMintsTheShapesTheRosterAccepts pairs the minter's two arms with the
// byte shapes internal/panel/subagent_corr_task_prefix_328_test.go measures on the
// join, so neither file is asserting a shape nobody produces. Arm B needs a call
// with no id at all, and the openai-chat adapter substitutes its own
// call_unnamed_<n> before the loop sees it (internal/llm/openaichat/toolasm.go:114),
// so no golden replay in this tree reaches that arm - it is exercised here on the
// minter itself, with the task id a real run reported.
func TestTheLoopMintsTheShapesTheRosterAccepts(t *testing.T) {
	h := newHarness(t, "tool-then-text")
	res := h.run("现在天气怎么样")
	if res.TaskID == "" {
		t.Fatal("the loop reported no task id to build a correlation on")
	}
	armB := callCorr(res.TaskID, "", 0)
	if armB != res.TaskID+"#call-0" {
		t.Fatalf("callCorr arm B = %q, want %q: the shape the panel case feeds the roster is "+
			"no longer the shape the loop mints", armB, res.TaskID+"#call-0")
	}
	armA := callCorr(res.TaskID, "call_e1", 0)
	if armA != res.TaskID+"#call_e1" {
		t.Fatalf("callCorr arm A = %q, want %q", armA, res.TaskID+"#call_e1")
	}
	for _, corr := range []string{armA, armB} {
		sect := panel.TaskRosterSectionFrom(rosterWithTwoRows(res.TaskID, "sibling-of-"+res.TaskID),
			[]panel.ApprovalCardView{cardForCorr(corr)})
		if !rowOf(t, sect, res.TaskID).BlockedOnApproval {
			t.Errorf("a correlation the loop mints (%q) does not light its own row", corr)
		}
		if rowOf(t, sect, "sibling-of-"+res.TaskID).BlockedOnApproval {
			t.Errorf("correlation %q lit a row it was not built on", corr)
		}
	}
}
