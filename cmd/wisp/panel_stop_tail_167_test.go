package main

// Ticket 167 AC#4's own shot of the cancel family: the panel stop must land
// while the assistant is BUSY, and after it lands no chunk may still be written.
//
// The family's division of labour is deliberate (the dispatch names it: 与票
// 164 AC#4 同族，各自钉各自那发). Ticket 164's own nail
// (internal/tools/task_164_ac4_cancel_tail_test.go) pins the ROSTER-ROW tail:
// after a row settles, TaskBackfill refuses the next write. This file pins the
// ROOT text task's STREAM tail: after the root the panel stopped has settled,
// the results surface (the StreamLog the snapshot's results section reads,
// fed by consoleSink's EvTextDelta / EvDone branches, run.go) receives no
// further delta and the chunk carries its Done mark.
//
// The exits reused here are the measured ones (no second stop path was built):
//
//	panel roster.cancel      the panel chain's cancel door
//	                         (internal/panel/bridge.go MethodRosterCancel)
//	cancelRosterRowWave2     the seam, root branch: the root task's own
//	                         agent.RunningTask.Cancel (loop.go:342), already
//	                         confirm-gated (ticket 38's root test)
//	consoleSink / StreamLog  the one chunk surface the panel reads
//
// The busy proof is IN the test, not assumed: the cancel is issued only after
// at least one text delta has already landed for the root task, so "stop while
// busy" is a state the run was verifiably in, not a coin toss.

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/panel"
)

func rootChunk167(rt *agentRuntime, rootID string) (panel.ResultChunk, bool) {
	for _, c := range rt.stream.Chunks() {
		if c.CorrelationID == rootID {
			return c, true
		}
	}
	return panel.ResultChunk{}, false
}

// waitForFirstDelta167 blocks until the root task's stream row holds at least
// one delta - the observable "the assistant is busy writing".
func waitForFirstDelta167(rt *agentRuntime, rootID string) (panel.ResultChunk, bool) {
	for i := 0; i < 2400; i++ {
		if c, ok := rootChunk167(rt, rootID); ok && c.Text != "" {
			return c, true
		}
		time.Sleep(5 * time.Millisecond)
	}
	return panel.ResultChunk{}, false
}

func Test167AC4BusyStopLandsAndLeavesNoWritingTail(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	// A wide latency window: the root's echo streams in chunks for seconds, so
	// "cancel mid-stream" is a state, not a coin toss (the 38 wave's shape).
	f.srv.Control(t, "/__control/latency", `{"ms":300}`)
	door := newRosterCancelDoor(nil, "test")
	var (
		driven    *agentRuntime
		rootID    string
		rootIDErr error
		cancelOK  bool
		stopErr   error
	)
	f.rtHook = func(rt *agentRuntime) {
		driven = rt
		door.bindSeam(rt.cancelRosterRowWave2)
		go func() {
			id, err := waitRootRow197(rt)
			rootID, rootIDErr = id, err
			if err != nil {
				return
			}
			// BUSY PROOF: wait until the root's own stream row carries at
			// least one delta. A stop issued here is a stop issued while the
			// assistant was verifiably mid-reply - not in an idle gap.
			if _, ok := waitForFirstDelta167(rt, id); !ok {
				return
			}
			_, stopErr = door.HandleRosterCancel(t.Context(), panel.ComposerRequest{
				Method: panel.MethodRosterCancel, RequestID: "167-stop-1",
				Source: panel.ComposerRequestSource, RosterTaskID: id,
				PrivacyConfirmed: true,
			})
			cancelOK = stopErr == nil
		}()
	}
	code := f.run(carrier197RootTask)

	if rootIDErr != nil {
		t.Fatalf("the run never produced a root row to stop: %v", rootIDErr)
	}
	if stopErr != nil {
		t.Fatalf("the panel cancel was refused while the assistant was busy: %v", stopErr)
	}
	if !cancelOK {
		t.Fatal("the busy-state stop never fired: the door returned without a receipt")
	}
	if code != 1 {
		t.Errorf("exit after the busy stop = %d, want 1 (the cancelled class; 0 means the stop never really fired)\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}

	// The settle evidence, on the durable row: the root's task_log closed with
	// the cancelled state. (The root ROSTER row deliberately files no state for
	// bodyless endings - run.go's backfill refuses cancellations whole via the
	// empty-text gate - so the D43-Muted settle reading belongs to the roster
	// branch of the family, ticket 164's nail, not to this root-row shot. The
	// exit code above plus this closed row are the root's real terminal facts.)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tl, err := f.openStore().TaskLogByID(ctx, rootID)
	if err != nil {
		t.Fatalf("task_log for the stopped root: %v", err)
	}
	if tl.State != "cancelled" || tl.EndedAt == nil {
		t.Errorf("task_log = state %q ended=%v, want a closed cancelled row (the busy stop must really land)", tl.State, tl.EndedAt != nil)
	}

	// THE TAIL: after the task has settled (EvDone closed its chunk), no chunk
	// may still land. One delta racing the cancel BEFORE the settle is physics;
	// the tail this family forbids is the NEXT write, after the row is closed -
	// the same precise statement ticket 164's own nail states for the roster.
	before, ok := rootChunk167(driven, rootID)
	if !ok {
		t.Fatal("the stopped root left no chunk row on the results surface")
	}
	if !before.Done {
		t.Errorf("the root's chunk is not marked Done after the settle: EvDone's Close never ran (stream text=%q)", before.Text)
	}
	if before.Text == "" {
		t.Error("the busy proof failed: no delta ever landed before the stop, so 'stop while busy' was never under test")
	}
	time.Sleep(600 * time.Millisecond)
	after, ok := rootChunk167(driven, rootID)
	if !ok {
		t.Fatal("the chunk row vanished after the settle")
	}
	if after.Text != before.Text {
		t.Errorf("停止后仍有 chunk 落地：settle 时 %d 字，之后 %d 字（+%q）——取消之后不许有还在写的尾巴",
			len([]rune(before.Text)), len([]rune(after.Text)), strings.TrimPrefix(after.Text, before.Text))
	}
}
