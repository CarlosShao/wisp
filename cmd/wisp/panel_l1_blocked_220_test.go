package main

// Ticket 220 AC#2's last hop, measured on the production assembly (leg 220-r2).
// 220-r1 (480b970d2) landed the render half - corr+taskID 双查, the pump's optional
// L1Windows reader (internal/panel/pump.go:200, consumed :323-327) and the gate's
// read-only enumeration (window_read.go:63) - but NO PRODUCER: the composition root
// never pointed its reader at that enumeration, so on a real `wisp run` it stayed
// nil and every L1 window stayed invisible (A408's disease one hop away). Nothing
// here builds a packet: every byte asserted came out of a Publish(), and whether it
// was the RUN's publish or this file's sample is pinned by the snapshotCount.
// DO NOT READ THE PANEL-SIDE FIXTURE AS THIS PATH: internal/panel's
// subagent_blocked_220_test.go:151-152 feeds the pump by hand and measures the
// render; a host whose reader is nil keeps that render green forever. No render
// semantics moved here (AC#4: still a boolean, no new D43 name).
// HOW A REAL WINDOW IS REACHED: L1 opens only when the mode does not silence it, so
// the fixture keeps the default ask_every_step and dispatches an in-allowlist fs.write
// (reversible write, risk/assessor.go:47) through the assembled bridge, under the ROOT
// LOOP'S OWN task id - the one id the roster has a row for. The readings run on their
// own goroutine because s.onRuntime runs BEFORE execute(), so a blocking hook would
// never see a root row (subagent_carrier_197_test.go:16-23). Window 2-3s
// (approval/queue.go:130-135), pinned below to its ceiling, so the readings sample a
// live instant instead of racing it.

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
)

var l1RootTask220 = "总结一下 " + strings.Repeat("w", 300)

// l1ProbeContent lands only if the window ran out unopposed: on the L1 route the
// timeout MEANS EXECUTE (SPEC-06 §2, gate.go:338 books it ANSWER-EXPIRED).
const l1ProbeContent = "written through a live L1 confirmation window (ticket 220 AC#2)\n"

// pinL1WindowCeiling220 appends [risk] l1_window_sec = 3, the top of SPEC-06 §2's own
// band: nothing is loosened (the default 2s is inside it) and permission_mode stays
// at ask_every_step, which is what keeps the L1 route unsilenced.
func pinL1WindowCeiling220(t *testing.T, dataDir string) {
	t.Helper()
	cfgPath := filepath.Join(dataDir, configFileName)
	old, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cfgPath, append(append([]byte(nil), old...),
		[]byte("\n[risk]\nl1_window_sec = 3\n")...), 0o600); err != nil {
		t.Fatal(err)
	}
}

// l1RowRead220 pulls one roster row's blockedOnApproval, whether that row is on the
// packet at all, and the depth of the answerable pending section, straight out of
// packet bytes: the way a page reads them, with no Go-side help and no *testing.T,
// because these readings run off the test goroutine.
func l1RowRead220(data []byte, taskID string) (found, blocked bool, pending int) {
	var wire struct {
		Pending []json.RawMessage `json:"pending"`
		Tasks   carrier197Wire    `json:"tasks"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return false, false, 0
	}
	for _, r := range wire.Tasks.Rows {
		if r.TaskID == taskID {
			return true, r.BlockedOnApproval, len(wire.Pending)
		}
	}
	return false, false, len(wire.Pending)
}

// l1Probe220 carries every reading the probe makes. It records, it never fails.
type l1Probe220 struct {
	rootID                                 string
	setupErr                               string
	wins                                   []approval.L1Window
	before, atLit, liveQueue, pendingAtLit int
	preFound, preBlocked                   bool
	litBlocked, postBlocked                bool
	litBytes                               []byte
	outcomeText                            string
	outcomeIsError, writeExecuted          bool
}

// probeL1Window220 takes the four readings off one assembled runtime: everything it
// touches is production state - the gate's own enumeration, the bridge this run
// built, the packet this run retained.
func probeL1Window220(rt *agentRuntime, target, corr string) l1Probe220 {
	var r l1Probe220
	rootID, err := waitRootRow197(rt)
	if err != nil {
		r.setupErr = err.Error()
		return r
	}
	r.rootID = rootID
	// The root loop owns this task's admission; re-admitting keeps the window
	// reachable for the whole probe. Nothing here revokes it - pulling the admission
	// out from under a running loop fails its later calls closed.
	rt.gate.AdmitTextTask(rootID)

	// (1) BEFORE: the row exists and reads unblocked while nothing waits. This is the
	// 负向尺 the ticket's 正控 rule demands, and the reading a false-blocking 反形
	// has to redden.
	rt.publishPanelSnapshot()
	_, data, seen := rt.lastPanelSnapshot()
	if !seen {
		r.setupErr = "the run retained no packet, so nothing below could be read off it"
		return r
	}
	r.preFound, r.preBlocked, _ = l1RowRead220(data, rootID)
	r.before = rt.snapshotCount()

	// (2) Reach a REAL L1 window through the assembled bridge.
	args, _ := json.Marshal(map[string]string{
		"path": filepath.ToSlash(target), "content": l1ProbeContent,
	})
	done := make(chan struct{})
	go func() {
		defer close(done)
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		out, e := rt.bridge.Execute(ctx, agent.ToolRequest{
			TaskID: rootID, CorrelationID: corr, CallID: "l1-220-1",
			Name: "fs.write", Args: args,
		})
		r.outcomeText, r.outcomeIsError = out.Text, out.IsError
		if e != nil {
			r.outcomeText = fmt.Sprintf("bridge error: %v (%s)", e, out.Text)
			r.outcomeIsError = true
		}
	}()
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		if r.wins = rt.gate.LiveL1Windows(); len(r.wins) > 0 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if len(r.wins) == 0 {
		<-done
		r.setupErr = "no L1 window ever opened for this call: the route was silenced, refused, or never asked"
		return r
	}

	// (3) THE LIT INSTANT. This probe publishes nothing between (1) and here, so a
	// retained packet that lights the row can only be one the RUN built:
	// consoleApprovalUI.Prompt publishes (run.go:1451) and PendingWindow prompts from
	// inside the open window (gate.go:283-294).
	litDeadline := time.Now().Add(2500 * time.Millisecond)
	for time.Now().Before(litDeadline) {
		_, data2, _ := rt.lastPanelSnapshot()
		r.atLit = rt.snapshotCount()
		if f, b, pend := l1RowRead220(data2, rootID); f && b {
			r.litBlocked, r.pendingAtLit = true, pend
			r.litBytes = append([]byte(nil), data2...)
			r.liveQueue = len(rt.gate.Queue().LiveApprovals())
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	<-done
	body, rerr := os.ReadFile(target)
	r.writeExecuted = rerr == nil && string(body) == l1ProbeContent

	// (4) AFTER: the window is gone, so the same reader must change its answer.
	for i := 0; i < 2000 && len(rt.gate.LiveL1Windows()) > 0; i++ {
		time.Sleep(5 * time.Millisecond)
	}
	rt.publishPanelSnapshot()
	_, post, _ := rt.lastPanelSnapshot()
	_, r.postBlocked, _ = l1RowRead220(post, rootID)
	return r
}

func TestRunPacketLightsTheRosterRowOfItsLiveL1Window(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	// Ticket 375: an L1 window only opens when somebody could answer, so this run
	// needs an answer source (an attached, silent console) or the guard refuses
	// the call before there is any window for the roster to light.
	f.reply = strings.NewReader("")
	pinL1WindowCeiling220(t, f.dir)
	f.srv.Control(t, "/__control/latency", `{"ms":2500}`)
	target := filepath.Join(f.dir, "l1-220-reversible-write.txt")
	corr := "l1-220-corr-" + fmt.Sprint(time.Now().UnixNano())

	var probe l1Probe220
	done := make(chan struct{})
	f.rtHook = func(rt *agentRuntime) {
		go func() {
			defer close(done)
			probe = probeL1Window220(rt, target, corr)
		}()
	}
	code := f.run(l1RootTask220)
	<-done
	if code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}
	if probe.setupErr != "" {
		t.Fatalf("the probe never reached a live window: %s\nstdout:\n%s\nstderr:\n%s",
			probe.setupErr, f.out.String(), f.err.String())
	}
	// Both ids are the gate's own readings - a copy, not a judgement (window_read.go:36-47).
	if len(probe.wins) != 1 || probe.wins[0].CorrelationID != corr || probe.wins[0].TaskID != probe.rootID {
		t.Fatalf("the enumeration read back %+v, want corr %q and task %q", probe.wins, corr, probe.rootID)
	}
	if !probe.preFound || probe.preBlocked {
		t.Fatalf("the pre-window packet carried no row for %q, or already lit it: the readings below "+
			"would describe a constant, not a state", probe.rootID)
	}
	if !probe.litBlocked {
		t.Errorf("PRE-CHANGE / WIRING READING: an L1 window was live on task %q (corr %q) and its roster "+
			"row still read blockedOnApproval=false - the pump's reader is only as good as the producer "+
			"the composition root points at it (cmd/wisp/run.go's PumpSources literal)", probe.rootID, corr)
	}
	if probe.atLit <= probe.before {
		t.Errorf("the row read as blocked while snapshotCount never moved past %d: this probe published once "+
			"(the pre-window sample), so a lit packet at or below that count is those bytes again", probe.before)
	}
	if probe.pendingAtLit != 0 || probe.liveQueue != 0 {
		t.Errorf("an L1 window reached the answerable section: snapshot.pending %d, LiveApprovals %d, want 0/0 "+
			"- a window is not a card, and listing it as one is the second door 220-r1 refused (ticket 146's "+
			"readings on LiveApprovals)", probe.pendingAtLit, probe.liveQueue)
	}
	if probe.postBlocked {
		t.Errorf("the row still read blocked after the window closed: the dimension must move with the state, " +
			"not print one word whatever happens (ticket 181 AC#7's 空转 family)")
	}
	if !probe.writeExecuted || probe.outcomeIsError {
		t.Errorf("the L1 route never ran this run's own write (executed=%v outcome=%q), so the readings above "+
			"describe a call that did not take this route", probe.writeExecuted, probe.outcomeText)
	}
	// The production surface's own words: it displayed the countdown and the gate
	// booked the expiry - both from the run, not from anything this file called.
	expired := "approval: ANSWER-EXPIRED corr=" + corr + " tool=fs.write decision=timeout->execute"
	for _, want := range []string{"[确认 L1 fs.write]", expired} {
		if !strings.Contains(f.out.String(), want) && !strings.Contains(f.err.String(), want) {
			t.Errorf("neither surface ever said %q, so the window's own course is unproven:\nout:%s\nerr:%s",
				want, f.out.String(), f.err.String())
		}
	}
	if len(ledgerSummaries145(t, f.dir)) == 0 {
		t.Error("the run booked no panel record, so the packets above never went through a production exit")
	}
	t.Logf("L1 window task=%s corr=%s | pre blocked=%v | lit at publish #%d (was #%d; pending %d, queue %d) "+
		"| post blocked=%v | lit bytes: %s", probe.rootID, corr, probe.preBlocked, probe.atLit, probe.before,
		probe.pendingAtLit, probe.liveQueue, probe.postBlocked, probe.litBytes)
}

// WHAT THIS FILE DOES NOT CLAIM: the transport is still ticket 33/35's, so nothing here
// shows a page rendering the lit row, only that the packet a page would read says it;
// the lit instant is polled, not frozen (the snapshotCount pins are what keep that from
// being a claim); and the reader is non-exported on purpose - zero new exported names,
// zero new panel.* names (A562 §2). The assembly root's call-site count is a shell
// reading in docs/evidence/s1/220-wiring.md, and this case is its behavioural twin: the
// leg's 去牙 mutation "delete the assembly line" reddens this case on its own.
