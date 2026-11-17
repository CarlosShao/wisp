package main

// Readface wave 2 (2026-11-17), ticket 38 AC#4's cmd/wisp teeth: the cancel
// door over the host's EXISTING cancel paths, driven against a real runtime
// (newRunFixture's real bridge, real roster, real spawned subagent) - the
// reuse-not-reinvention claim measured, not asserted.
//
// The applied-steps reading the ticket asks for comes off the surfaces that
// actually hold it after a mid-flight stop: the settled roster row (the
// child's own finalize files the terminal state) and the tool_call forensics
// rows the bridge booked under the child's task id - one real fs.read
// dispatched before the stop is the step that "already landed", which is
// D31's evidence half for a cancelled task.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/statemachine"
	"github.com/CarlosShao/wisp/internal/tools"
)

func cancelEnvelope(taskID string) string {
	b, err := json.Marshal(map[string]string{
		"method": panel.MethodRosterCancel, "requestId": "rw2-c",
		"source": panel.ComposerRequestSource, "rosterTaskId": taskID,
	})
	if err != nil {
		panic(err)
	}
	return string(b)
}

// TestRosterCancelDoorRefusesTheEmptySelectorAndTheUnboundSeam: the two
// refusals the door owns itself - an id-less cancel stops nothing and says
// so; a machine with no task pipeline refuses by name instead of dropping.
func TestRosterCancelDoorRefusesTheEmptySelectorAndTheUnboundSeam(t *testing.T) {
	door := newRosterCancelDoor(nil, "test")
	_, err := door.HandleRosterCancel(context.Background(), panel.ComposerRequest{
		Method: panel.MethodRosterCancel, RequestID: "rw2-e", Source: panel.ComposerRequestSource,
	})
	if err == nil || !strings.Contains(err.Error(), "rosterTaskId") {
		t.Fatalf("empty selector: err = %v, want the named rosterTaskId refusal", err)
	}
	door2 := newRosterCancelDoor(nil, "test")
	_, err = door2.HandleRosterCancel(context.Background(), panel.ComposerRequest{
		Method: panel.MethodRosterCancel, RequestID: "rw2-n", Source: panel.ComposerRequestSource,
		RosterTaskID: "child-x",
	})
	if !errors.Is(err, panel.ErrNoHandlerAttached) {
		t.Fatalf("unbound seam: err = %v, want ErrNoHandlerAttached", err)
	}
}

// TestRosterCancelDoorReceiptPassesTheHostRefusalThrough: a row the host
// cannot stop (unknown id) surfaces the roster's own fail-closed reason, not
// a door-invented one.
func TestRosterCancelDoorReceiptPassesTheHostRefusalThrough(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	f.reply = strings.NewReader("")
	var receipt string
	var receiptErr error
	door := newRosterCancelDoor(nil, "test")
	f.rtHook = func(rt *agentRuntime) {
		door.bindSeam(rt.cancelRosterRowWave2)
		receipt, receiptErr = door.HandleRosterCancel(context.Background(), panel.ComposerRequest{
			Method: panel.MethodRosterCancel, RequestID: "rw2-u", Source: panel.ComposerRequestSource,
			RosterTaskID: "no-such-row-w2",
		})
	}
	if code := f.run("读一下这张卡上的读数"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}
	if receiptErr == nil {
		t.Fatal("an unknown row id answered success: the door would be a fake stop")
	}
	// The reason must be the roster's own sentence (164 定案②'s fail-closed
	// pair), carried verbatim - that is the "复用既有语义" claim measured.
	if !strings.Contains(receiptErr.Error(), "查不到任务 no-such-row-w2") {
		t.Errorf("refusal = %v, want the roster's own unknown-row sentence", receiptErr)
	}
	if receipt != "" {
		t.Errorf("a refused cancel returned a receipt %q: only a stop that fired gets one", receipt)
	}
}

// TestRosterCancelDoorStopsARealSubagentAndTheForensicsCarryTheAppliedSteps
// is the production reading: a REAL spawned subagent (the assembled bridge's
// own task.spawn carrier), one REAL fs.read dispatched under its task id
// before the stop, then the door cancels it through the seam. The settled
// row and the forensics rows are the applied-steps report the ticket asks
// for, read off live objects - not re-derived by the door.
func TestRosterCancelDoorStopsARealSubagentAndTheForensicsCarryTheAppliedSteps(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	// A wide latency window: the child's echo streams in chunks for seconds,
	// so "cancel mid-flight" is a state, not a coin toss.
	f.srv.Control(t, "/__control/latency", `{"ms":300}`)
	// The applied step is a REAL file the fixture's allowlisted dir holds.
	applied := filepath.Join(f.dir, "applied-step-w2.txt")
	if err := os.WriteFile(applied, []byte("票38 取消读数的落地步骤"), 0o600); err != nil {
		t.Fatal(err)
	}
	sp := newSpawn197()
	door := newRosterCancelDoor(nil, "test")
	var (
		mu         sync.Mutex
		driven     *agentRuntime
		childID    string
		receipt    string
		cancelErr  error
		stoppedRow tools.TaskOutput
	)
	// The child's book is long enough that its echo keeps streaming past the
	// cancel: the mock echoes the message, so length buys window.
	longChild := "把这句话原样说出来：" + strings.Repeat("w", 420)
	f.rtHook = func(rt *agentRuntime) {
		driven = rt
		door.bindSeam(rt.cancelRosterRowWave2)
		go sp.launch(rt, longChild, 1)
		go func() {
			rootID, err := waitRootRow197(rt)
			if err != nil {
				sp.setErr(err)
				return
			}
			id, idErr := waitChildRow197(rt, rootID)
			if idErr != nil {
				sp.setErr(idErr)
				return
			}
			mu.Lock()
			childID = id
			mu.Unlock()
			// The step that lands BEFORE the stop: one real fs.read through
			// the assembled bridge under the child's own admitted identity -
			// the same shape dispatchTask47 drives task.* with. L0 under the
			// fixture's allowlisted dir, so no card and no waiting.
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			out, eerr := rt.bridge.Execute(ctx, agent.ToolRequest{
				TaskID: id, CorrelationID: id, CallID: "w2-applied-step",
				Name: "fs.read", Args: json.RawMessage(fmt.Sprintf(`{"path":%q}`, applied)),
			})
			if eerr != nil {
				sp.setErr(eerr)
				return
			}
			if out.IsError {
				sp.setErr(fmt.Errorf("fs.read under the child id refused: %s", out.Text))
				return
			}
			// THE production cancel: the panel door over the runtime's own
			// seam. The receipt names the row; the stop is the handle the
			// spawner attached, fired.
			receipt, cancelErr = door.HandleRosterCancel(context.Background(), panel.ComposerRequest{
				Method: panel.MethodRosterCancel, RequestID: "rw2-real",
				Source: panel.ComposerRequestSource, RosterTaskID: id,
			})
		}()
	}
	if code := f.run(carrier197RootTask); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}
	<-sp.done
	sp.mu.Lock()
	spawnErr := sp.err
	sp.mu.Unlock()
	mu.Lock()
	id, rcpt, cerr := childID, receipt, cancelErr
	mu.Unlock()
	if spawnErr != nil {
		t.Fatalf("the spawn/dispatch behind the cancel failed: %v", spawnErr)
	}
	if id == "" {
		t.Fatal("the run never produced a subagent row to cancel")
	}
	if cerr != nil {
		t.Fatalf("the panel cancel was refused: %v", cerr)
	}
	if !strings.Contains(rcpt, id) || !strings.Contains(rcpt, "发出停止") {
		t.Errorf("receipt = %q, want it to name the row and state the stop was issued", rcpt)
	}
	// The settled row: the child's own finalize files the terminal state - the
	// D43 name, not a door-invented word.
	rt := driven
	for i := 0; i < 1200; i++ {
		if out, ok := rt.tasks.Look(id); ok && out.State != "" {
			mu.Lock()
			stoppedRow = out
			mu.Unlock()
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if stoppedRow.State == "" {
		t.Fatalf("the cancelled child never settled (state empty after the wait): %+v", stoppedRow)
	}
	if !statemachine.Valid(stoppedRow.State) {
		t.Errorf("settled state %q is not one of D43's names (the row must not carry a private word)", stoppedRow.State)
	}
	// The audit's 10.3/10.2-3 tooth (readface-adversarial-20261117-1432.md): the
	// cancelled child's real terminal name is D43's Muted - a seam that answers
	// "stopped" without stopping (the M16 mutation: known id => true, nothing
	// halted) lets the child settle to its own finished state instead, which is
	// a different D43 word, and this line is what turns that red.
	if stoppedRow.State != statemachine.StateMuted {
		t.Errorf("settled state = %q, want D43's Muted (the cancelled child's real terminal name; any other word means the stop did not really fire)", stoppedRow.State)
	}
	if stoppedRow.Kind != tools.TaskKindSubagent {
		t.Errorf("settled kind = %q, want a subagent row", stoppedRow.Kind)
	}
	// The applied-steps reading, off the surfaces that hold it: the bridge
	// booked the completed fs.read under the child's task id before the stop,
	// and the cancelled task books no further rows. The run closed its own
	// store handle, so this reopens the same database the rows landed in
	// (openStore, the fixture's post-run read shape).
	st := f.openStore()
	rows, err := st.ListToolCallsByTask(context.Background(), id)
	if err != nil {
		t.Fatalf("forensics read: %v", err)
	}
	appliedRows, cancelledRows := 0, 0
	for _, r := range rows {
		switch r.Outcome {
		case "success":
			appliedRows++
			if r.Tool != "fs.read" {
				t.Errorf("applied row tool = %q, want fs.read", r.Tool)
			}
		case "cancelled":
			cancelledRows++
		}
	}
	if appliedRows != 1 {
		t.Errorf("applied steps = %d, want 1 (the fs.read that really landed before the stop): %+v", appliedRows, rows)
	}
	if cancelledRows != 0 {
		t.Errorf("cancelled forensics rows = %d, want 0 (the stop landed between calls)", cancelledRows)
	}
	t.Logf("取消后 applied-steps 读数：task=%s 终态=%s；已落地 %d 步（tool_call success，tool=fs.read，路径 %s）；"+
		"被中止 %d 步；回执：%s", id, string(stoppedRow.State), appliedRows, applied, cancelledRows, rcpt)
}

// TestRosterCancelDoorRootStopNeedsTheConfirmBit is the audit R1's second
// class (readface-adversarial-20261117-1432.md §9.1-6): roster.cancel can stop
// the ROOT task - stopping the whole run - and until today the root branch of
// the seam executed with zero confirmation. Now the root branch demands the
// SAME destructive bit the privacy doors gate on (ComposerRequest
// .PrivacyConfirmed, JSON key privacyConfirmed): driven against a REAL
// running root (the fixture's own task), a bit-less cancel of the root's id
// is refused BY NAME and the run keeps going; the same id WITH the bit gets
// the receipt and the run exits the cancelled way (runExitCode's rc=1 for
// StatusCancelled) - which is also what reddens the "returns true, stops
// nothing" mutation: a fake stop leaves the run finishing normally at rc=0.
func TestRosterCancelDoorRootStopNeedsTheConfirmBit(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	// A wide latency window: the root's echo streams in chunks for seconds, so
	// "cancel mid-run" is a state, not a coin toss.
	f.srv.Control(t, "/__control/latency", `{"ms":300}`)
	door := newRosterCancelDoor(nil, "test")
	var (
		mu         sync.Mutex
		rootID     string
		waitErr    error
		refusedErr error
		receipt    string
		confirmErr error
	)
	f.rtHook = func(rt *agentRuntime) {
		door.bindSeam(rt.cancelRosterRowWave2)
		go func() {
			id, err := waitRootRow197(rt)
			mu.Lock()
			rootID, waitErr = id, err
			mu.Unlock()
			if err != nil {
				return
			}
			// THE GATE'S NEGATIVE: no confirm bit, the root's own id. The
			// refusal must name the bit it needs, and the root's handle must
			// not be touched (the positive below still finds it live).
			_, refusedErr = door.HandleRosterCancel(context.Background(), panel.ComposerRequest{
				Method: panel.MethodRosterCancel, RequestID: "rw2-root-nobit",
				Source: panel.ComposerRequestSource, RosterTaskID: id,
			})
			// THE GATE'S POSITIVE: the same id WITH the bit - the stop is
			// issued and the receipt says so.
			receipt, confirmErr = door.HandleRosterCancel(context.Background(), panel.ComposerRequest{
				Method: panel.MethodRosterCancel, RequestID: "rw2-root-ok",
				Source: panel.ComposerRequestSource, RosterTaskID: id,
				PrivacyConfirmed: true,
			})
		}()
	}
	code := f.run(carrier197RootTask)
	mu.Lock()
	id, werr, rerr, rcpt, cerr := rootID, waitErr, refusedErr, receipt, confirmErr
	mu.Unlock()
	if werr != nil {
		t.Fatalf("the run never produced a root row to cancel: %v", werr)
	}
	if id == "" {
		t.Fatal("the run never produced a root row to cancel")
	}
	if rerr == nil {
		t.Fatal("a bit-less root cancel was executed: the root stop's confirm gate is gone")
	}
	if !strings.Contains(rerr.Error(), "privacyConfirmed") {
		t.Errorf("no-bit root refusal = %v, want it to name the confirm bit it needs", rerr)
	}
	if cerr != nil {
		t.Fatalf("the confirmed root cancel was refused: %v", cerr)
	}
	if !strings.Contains(rcpt, id) || !strings.Contains(rcpt, "发出停止") {
		t.Errorf("receipt = %q, want it to name the root row and state the stop was issued", rcpt)
	}
	// The stop REALLY landed: the root task ended cancelled, and the process
	// exit code reflects that class (run.go runExitCode: StatusCancelled ->
	// ClassCancelled -> 1). A seam that faked the stop would leave this run
	// completing normally at rc=0 and this line red.
	if code != 1 {
		t.Errorf("exit code after the confirmed root cancel = %d, want 1 (the cancelled class; 0 would mean the stop never fired)\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}
	t.Logf("根任务确认位读数：task=%s 无位拒绝=%v；带位回执：%s；exit=%d", id, rerr, rcpt, code)
}
