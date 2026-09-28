package main

// Ticket 197 AC#4's production judgement (leg 197-r4): a roster row that a live
// approval card is holding reads blockedOnApproval=true ON A PACKET THIS RUN
// BUILT.
//
// WHY THIS CASE EXISTED TO WRITE. 197-r3's file closed its list with "the join
// needs a card whose correlation id equals a CHILD task id, and nothing in this
// tree can make that happen from the outside". 197-r3b corrected that in its own
// §⑤1 - the sentence was too strong, the host already has the line:
// cmd/wisp/run.go:577-579 signs a real L2 card for a host-registered task
// (`AdmitTextTask` then `PendingApproval`), which is exactly how a mode switch
// gets its R20/M4 confirmation. So the state is reachable with the assembled
// stack and NO production change; what was missing was someone calling it.
//
// WHAT IS REUSED AND WHY IT IS NOT A MOCK. The child is spawned by the same
// assembled bridge carrier197 uses (its roster row is minted by the spawner and
// by the child loop's own admission hook, never by this file), the card is
// admitted and raised by the SAME gate the run assembled (rt.gate, the D47
// registration the loop's own admitTask also uses), and the packet is the one
// rt.pump.Publish built - triggered by consoleApprovalUI.Prompt, i.e. by the run
// itself at the moment the card goes on screen (run.go:1027-1032 says so in the
// source). Nothing below constructs a panel.Snapshot or a NativeVerdict, and
// nothing below answers a card.
//
// TWO READINGS, BOTH FROM THE QUEUE'S OWN STATE: while the card is up the row is
// blocked; after the queue no longer holds it the same row is not. That second
// half is what stops the field from being a constant, and it is also what makes
// the card's own answer an assertion: an unanswered L2 must resolve to a reject
// (C18/300s polarity), so this carrier cannot be read as an approval outlet -
// ticket 197 AC#5 and Q-49 丙 still hold with this case in the suite.

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/statemachine"
	"github.com/CarlosShao/wisp/internal/tools"
)

// blocked197Tool names the card's tool. It is NOT a registered C4 tool and no
// call route uses it - the same shape as modeSwitchToolName (run.go:553-556),
// which exists only so a card says what it is confirming. Naming a real tool
// here would put a call nobody made onto a security card.
const blocked197Tool = "ticket197.blocked.probe"

// blocked197Reason is the card's reason, carried verbatim into the packet by the
// same renderer the CLI diagnostic uses. It names what this case is, because a
// reader of the packet - or of a -v log - cannot otherwise tell a probe card
// from a real one.
const blocked197Reason = "票 197 AC#4 的正控：名册行要能读出「有一张授权卡正挂着」"

// TestRunPacketMarksTheRosterRowACardIsHolding is AC#4 on the production path.
func TestRunPacketMarksTheRosterRowACardIsHolding(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	// The latency switch is what makes "in flight" a window rather than a coin
	// toss: with it the child's row exists well before either loop finishes.
	f.srv.Control(t, "/__control/latency", `{"ms":80}`)

	sp := newSpawn197()
	var (
		mu       sync.Mutex
		driven   *agentRuntime
		childID  string
		spawnErr error
		cardAns  tools.Answer
		cardWhy  string
	)
	note := func(id string, err error) {
		mu.Lock()
		defer mu.Unlock()
		if id != "" {
			childID = id
		}
		if err != nil {
			spawnErr = err
		}
	}
	// The card's context is cancelled by the case, never by a user clicking
	// allow: an unanswered L2 is a reject, and that is the reading below.
	ctxCard, cancelCard := context.WithCancel(context.Background())
	defer cancelCard()
	cardDone := make(chan struct{})

	f.rtHook = func(rt *agentRuntime) {
		mu.Lock()
		driven = rt
		mu.Unlock()
		go sp.launch(rt, carrier197ChildPrompt, 1)
		go func() {
			defer close(cardDone)
			rootID, err := waitRootRow197(rt)
			if err != nil {
				note("", err)
				return
			}
			// The child's row is on this run's roster from its first model call
			// onward (tools.PublishSubagent files it before the loop's admission
			// hook even runs), so this waits for the IN-FLIGHT row: the shape
			// PLAN.md:1127's incident is about - a blocked child that every
			// interface quietly waits on.
			id, idErr := waitChildRow197(rt, rootID)
			if idErr != nil {
				note("", idErr)
				return
			}
			note(id, nil)
			// The host's own two entry points, in the order run.go's mode switch
			// uses them. The child loop's own registration may still be live here;
			// this one is what makes the card reachable independent of when the
			// child joins, and its revoke runs after the case has read both
			// readings.
			revoke := rt.gate.AdmitTextTask(id)
			defer revoke()
			ans, why := rt.gate.PendingApproval(ctxCard, tools.Decision{
				TaskID: id, CorrelationID: id, Tool: blocked197Tool,
				Level:  risk.L2,
				Reason: blocked197Reason,
				Mode:   rt.modes.PermissionMode(),
			})
			mu.Lock()
			cardAns, cardWhy = ans, why
			mu.Unlock()
		}()
	}

	if code := f.run(carrier197RootTask); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}
	<-sp.done
	sp.mu.Lock()
	rootID, callErr := sp.rootID, sp.err
	sp.mu.Unlock()
	mu.Lock()
	rt, id, err := driven, childID, spawnErr
	mu.Unlock()
	if callErr != nil {
		t.Fatalf("the task.spawn behind the card failed: %v", callErr)
	}
	if err != nil {
		t.Fatalf("the spawn behind the card failed: %v", err)
	}
	if id == "" {
		t.Fatal("this run never produced a subagent roster row to hold a card for")
	}

	// Read the packet the RUN published. No manual publish has happened at this
	// point, so a card in this packet can only have come from the queue being
	// live when the run's own trigger built it.
	var snap panel.Snapshot
	var data []byte
	for i := 0; i < 600; i++ {
		snap, data, _ = rt.lastPanelSnapshot()
		if len(snap.Pending) == 1 && snap.Pending[0].CorrelationID == id {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if len(snap.Pending) != 1 || snap.Pending[0].CorrelationID != id {
		t.Fatalf("no packet the run published carries the card on its queue (pending %+v, wanted one "+
			"card named %q); the run displayed %d card(s)", snap.Pending, id, rt.windowCount())
	}
	sect := packetTasks197(t, data)
	row := carrierRow(t, sect, id)
	if !row.BlockedOnApproval {
		t.Errorf("the packet carries a card for %q and its roster row still reads blockedOnApproval="+
			"false: 「被阻塞的子代理在所有界面上都不可见」 is the incident ticket 197 AC#4 names, and a "+
			"carrier that misses it is the reason that field exists", id)
	}
	// The row keeps every other fact it carried: this is a join, not a rewrite.
	if !row.StatusKnown || !statemachine.Valid(statemachine.State(row.Status)) {
		t.Errorf("a blocked row lost its status dimension: status=%q known=%v reason=%q, want one of "+
			"D43's 20 names (statemachine.Valid is the judge)", row.Status, row.StatusKnown, row.StatusReason)
	}
	if row.StreamKey != panel.SubagentStreamKey(id) {
		t.Errorf("streamKey = %q, want %q: 点进去那一页 must stay addressable on a blocked row",
			row.StreamKey, panel.SubagentStreamKey(id))
	}
	// The join is by correlation id, so exactly one row moves: the root is in the
	// same packet, is in flight, and has no card naming it. A flag that went
	// section-wide would dress every other row up as a waiting one.
	if rootID == "" {
		t.Fatal("the spawn recorded no root task id, so the non-blocked half below would prove nothing")
	}
	if rootRow := carrierRow(t, sect, rootID); rootRow.BlockedOnApproval {
		t.Errorf("the root row reads blocked although no card names it (the only card is %q): the join "+
			"is by correlation id, not by anything being in flight", id)
	}

	// The same packet was booked by the run's own exit: this is a production
	// reading, not a snapshot of a test's call (197-r3b's tie, reused verbatim).
	want := sha256Short(data)
	var booked bool
	for _, rec := range ledgerSummaries145(t, f.dir) {
		if rec["sha256"] == want {
			booked = true
		}
	}
	if !booked {
		t.Errorf("the blocked packet (%s) was never booked by a production publish", want)
	}

	// Second reading: take the card away and ask the same reader again. This half
	// DOES sample (publishPanelSnapshot is the run's own publisher, called here),
	// and it is what makes blockedOnApproval a report rather than a constant.
	cancelCard()
	<-cardDone
	rt.publishPanelSnapshot()
	after, afterData, seen := rt.lastPanelSnapshot()
	if !seen {
		t.Fatal("no packet after the card left the queue")
	}
	if len(after.Pending) != 0 {
		t.Fatalf("the queue still reports %d pending card(s) after its context ended: %+v",
			len(after.Pending), after.Pending)
	}
	afterRow := carrierRow(t, packetTasks197(t, afterData), id)
	if afterRow.BlockedOnApproval {
		t.Errorf("the card is gone and the row still reads blocked: the field would then be a constant, "+
			"and a constant is what ticket 181 AC#7 was filed over (rows after: %+v)", afterRow)
	}
	mu.Lock()
	ans, why := cardAns, cardWhy
	mu.Unlock()
	if ans != tools.AnswerReject {
		t.Errorf("the unanswered card resolved to %q (%q), want a reject: nothing on this path may "+
			"answer a card for the user (AC#5, Q-49 丙)", ans, why)
	}
	t.Logf("blocked row on the run's own packet: task=%s status=%s streamKey=%s card=%s pending=%d "+
		"bytes=%d sha256=%s | after the card: blocked=%v answer=%s why=%q",
		row.TaskID, row.Status, row.StreamKey, snap.Pending[0].Tool, len(snap.Pending),
		len(data), want, afterRow.BlockedOnApproval, ans, why)
}

// waitChildRow197 returns the first subagent roster row this run filed under
// rootID. Error, never Fatal: it runs on the card's goroutine.
func waitChildRow197(rt *agentRuntime, rootID string) (string, error) {
	for i := 0; i < 2400; i++ {
		if id := childRowAfterJoin197(rt, rootID, 1); id != "" {
			return id, nil
		}
		time.Sleep(5 * time.Millisecond)
	}
	return "", fmt.Errorf("the spawn this case dispatched never filed a subagent row under %q, "+
		"so there was no row to hang a card on", rootID)
}
