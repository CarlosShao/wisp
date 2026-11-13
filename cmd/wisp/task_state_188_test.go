package main

// Ticket 188, step 1 (the 188-c2 verdict's first move): the root task's status
// dimension has a producer again.
//
// What was true until now: tools.TaskBackfill was wired without a State, so the
// roster's root row could only answer 「宿主没有登记这一维」 while every subagent
// row named a real D43 state (ticket 197). What this step changes is one wiring
// at the composition root: the finished run's ending is spelled by taskLogState
// (ticket 365's one-liner, in ticket 196's memory vocabulary) and translated
// through memory.TaskStateFor - the one read door into ticket 196's
// memory-word -> D43-name table - so the name that reaches the roster is a
// table lookup, not a switch written here.
//
// The criterion below runs the REAL composition root (mockllm, no spawns, no
// cards) to a completed ending and reads the root row three ways, because the
// ticket's defect was never "a field exists" - it was "no production writer":
//
//	the roster record itself   tools.TaskOutput.State, a D43 name
//	the host's own judge       TaskOutput.StateAnswer() -> (name, "")
//	the panel's reader         taskRosterState() -> statusKnown=true
//
// The mutation that proves the wiring is load-bearing is recorded in
// docs/evidence/s1/188-root-task-state.md §4: removing the two assembly lines
// (the rootState lookup and the State: field) turns this test red with the
// root row back on 「宿主没有登记这一维」.

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/statemachine"
	"github.com/CarlosShao/wisp/internal/tools"
)

// Test188RootRowFilesTheD43NameTheLoopSpelled is the positive half of step 1:
// a completed run's own roster row carries the D43 name that run's ending maps
// to (completed -> "done" -> Settling), the same name on all three surfaces,
// with the identity MarkRoot filed (kind=root, label=the task text) still
// intact under it.
func Test188RootRowFilesTheD43NameTheLoopSpelled(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	var driven *agentRuntime
	f.rtHook = func(rt *agentRuntime) { driven = rt }
	const taskText = "总结一下 这份笔记"
	if code := f.run(taskText); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}
	if driven == nil {
		t.Fatal("the run never assembled a runtime")
	}

	var rootID string
	for _, row := range driven.taskRosterState().Rows {
		if row.Kind == tools.TaskKindRoot {
			if rootID != "" {
				t.Fatalf("two root rows in one run (%s, %s)", rootID, row.TaskID)
			}
			rootID = row.TaskID
		}
	}
	if rootID == "" {
		t.Fatal("the roster carries no root row at all")
	}

	// The record itself: the state dimension is a D43 name, and it is the one
	// this ending maps to - not empty, and not a memory-side word.
	rec, ok := driven.tasks.Look(rootID)
	if !ok {
		t.Fatalf("the root row vanished before Look: %s", rootID)
	}
	if rec.State != statemachine.StateSettling {
		t.Errorf("root row state = %q, want %q (completed -> done -> Settling): an empty "+
			"state here is the pre-188 shape this step exists to close",
			rec.State, statemachine.StateSettling)
	}
	if !statemachine.Valid(rec.State) {
		t.Errorf("root row state %q is not one of D43's 20 names", rec.State)
	}
	// The host's own judge answers the same name, with no refusal reason.
	if state, reason := rec.StateAnswer(); state != statemachine.StateSettling || reason != "" {
		t.Errorf("StateAnswer() = (%q, %q), want (Settling, empty)", state, reason)
	}
	// The identity MarkRoot filed survives the backfill's Record.
	if rec.Kind != tools.TaskKindRoot {
		t.Errorf("kind = %q, want %q", rec.Kind, tools.TaskKindRoot)
	}
	if rec.Label != taskText {
		t.Errorf("label = %q, want the task text %q", rec.Label, taskText)
	}

	// The panel's reader says the same thing: known, the same word, no reason.
	seen := false
	for _, row := range driven.taskRosterState().Rows {
		if row.TaskID != rootID {
			continue
		}
		seen = true
		if !row.StatusKnown {
			t.Errorf("panel row statusKnown = false (reason %q): the roster has a name "+
				"and the reader dropped it", row.StatusReason)
		}
		if row.Status != string(statemachine.StateSettling) {
			t.Errorf("panel row status = %q, want %q", row.Status, statemachine.StateSettling)
		}
		if row.StatusReason != "" {
			t.Errorf("panel row statusReason = %q for a row the roster filed a name for",
				row.StatusReason)
		}
	}
	if !seen {
		t.Fatalf("the panel reader lost the root row %s it had listed before Look", rootID)
	}
}

// Test188FailedRunStillFilesItsD43State is the panel-status-sources step
// (2026-11-13): a run that ends WITHOUT a body (a classified provider failure -
// the loop's finish publishes a message, not a body) used to be refused whole
// by the backfill's empty-text gate, so the roster carried no row at all and
// every surface answered 「宿主没有登记这一维」 for an ending the run DID spell.
// The status dimension and the body are now decoupled (task_backfill.go), so
// the ending the host already knew reaches the panel.
//
// Three readings, same shape as the completed-run criterion above: the roster
// record, the host's own judge, the panel's reader. Mutation (recorded with
// readings in docs/evidence/s1/188-panel-status-sources.md §4): restoring the
// gate to `case text == "":` turns this leg red - the failed run's row
// disappears and the reader is back on the 未登记 answer.
func Test188FailedRunStillFilesItsD43State(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	var driven *agentRuntime
	f.rtHook = func(rt *agentRuntime) { driven = rt }
	// fail_next(3): every attempt fails, the retry ladder gives up, the task
	// ends classified (the same shape TestRunTextTaskFailNextIsClassified and
	// Test365RunKeepsItsWordForTheControls/failed use).
	f.srv.QueueFault(t, adaptertest.Fault{Status: 500, Times: 3})
	code := f.run("总结一下 这份笔记")
	out, errb := f.out.String(), f.err.String()
	if code != 1 {
		t.Fatalf("exit = %d, want 1 for a classified provider failure\nstdout:\n%s\nstderr:\n%s", code, out, errb)
	}
	if driven == nil {
		t.Fatal("the run never assembled a runtime")
	}
	rootID := f.taskID()
	rec, ok := driven.tasks.Look(rootID)
	if !ok {
		t.Fatalf("the failed run left NO roster row: the empty-text gate ate the record whole\nstdout:\n%s", out)
	}
	if rec.State != statemachine.StateError {
		t.Errorf("roster state = %q, want Error (failed -> error -> Error): the ending the run "+
			"spelled must survive its missing body", string(rec.State))
	}
	if rec.Text != "" {
		t.Errorf("roster text = %q, want empty (this ending has no body to file)", rec.Text)
	}
	if state, reason := rec.StateAnswer(); state != statemachine.StateError || reason != "" {
		t.Errorf("StateAnswer() = (%q, %q), want (Error, empty)", string(state), reason)
	}
	seen := false
	for _, row := range driven.taskRosterState().Rows {
		if row.TaskID != rootID {
			continue
		}
		seen = true
		if !row.StatusKnown {
			t.Errorf("panel row statusKnown = false (reason %q): the roster has the ending and the reader dropped it", row.StatusReason)
		}
		if row.Status != string(statemachine.StateError) {
			t.Errorf("panel row status = %q, want %q", row.Status, statemachine.StateError)
		}
		if row.StatusReason != "" {
			t.Errorf("panel row statusReason = %q for a row the roster filed a name for", row.StatusReason)
		}
	}
	if !seen {
		t.Fatalf("the panel reader lost the root row %s the roster filed", rootID)
	}
	t.Logf("failed run: roster state=%q text=%d bytes, panel statusKnown=%v status=%q",
		string(rec.State), len(rec.Text), seen, statemachine.StateError)
}
