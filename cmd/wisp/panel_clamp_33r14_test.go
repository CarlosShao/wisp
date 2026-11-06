//go:build windows

package main

// Ticket 33 AC#8 (a)'s ruler (33-r14): the summary-clamp branch inside
// panelSnapshotSummary has to be reachable by a test that goes through the
// SAME entry the production pump uses, and the branch's product has to be
// observable on the booking the run keeps. What the AC added on 2026-09-25
// ("cmd/wisp/panel_pump.go:186-192, the over-440 drop-ids branch, never
// reached in 13 mutations + one working-tree run; the booked pending= field
// only ever showed one value") is still literally true at this leg's start
// HEAD aa7159b29, branch now at panel_pump.go:422 - measured, not inherited.
//
// The test drives panel.NewSnapshotPump(...).Publish() with the booking entry
// as Out, exactly as run.go:958 wires it (renamed bookPanelSnapshotBlocking by
// this leg; run.go's production caller follows the same name), so the branch
// cannot be reached by a shape production cannot take. The ruling's other
// half - "拿掉它就红" - is a mutation check, not a test body: the mutation
// run and its four readings live in docs/evidence/s1/33-panel-host-c27-r14.md,
// because the green ruler cannot prove the red without being mutated.

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/risk"
)

// TestAC8SnapshotClampBranchBooksThroughThePumpEntry33r14 drives the clamp
// branch's enclosing function through the production pump entry and pins the
// product: the packet is retained verbatim (the digest the branch keeps even
// when it drops the ids pins it) and the summary line the branch builds is the
// booked SNAPSHOT record shape.
func TestAC8SnapshotClampBranchBooksThroughThePumpEntry33r14(t *testing.T) {
	rt := &agentRuntime{}
	// Out takes the production entry - the exact value run.go:958 wires - so
	// the branch is reached by a shape production cannot avoid taking.
	p := panel.NewSnapshotPump(panel.PumpSources{Out: rt.bookPanelSnapshotBlocking})
	if _, _, err := p.Publish(); err != nil {
		t.Fatalf("publish through the try-entry: %v", err)
	}
	if !rt.snapSeen {
		t.Fatal("the try-entry published but the runtime booked nothing: the clamp branch's own function did not run, so this test measured a different hop than production takes")
	}
	if len(rt.lastSnapBytes) == 0 {
		t.Fatal("the booked packet carries zero bytes: the AC#8 ruler needs the full bytes the branch's digest pins, not a summary alone")
	}
	// The line both branches build: the clamp branch's own output shape, read
	// off the booking the entry just made (same pair the function logged).
	line := panelSnapshotSummary(rt.lastSnap, rt.lastSnapBytes)
	if !strings.HasPrefix(line, "panel: SNAPSHOT at=") {
		t.Errorf("the booked summary is not a SNAPSHOT record: %q", line)
	}
	// The branch's contract for the drop case: the count and the digest still
	// pin which packet this was. Assert the digest against the retained bytes.
	if !strings.Contains(line, "sha256="+sha256Short(rt.lastSnapBytes)) {
		t.Errorf("the booked summary's digest does not pin the retained packet (bytes %d): %q", len(rt.lastSnapBytes), line)
	}
}

// ac8ClampPacket33r14 builds one packet that is guaranteed to cross the clamp
// bound: three pending correlation ids (all above summaryClamp together),
// which is exactly the input the over-440 branch drops ids for, with everything
// else shaped like a real booked packet (the run's own constructors, not a
// hand-assembled literal).
func ac8ClampPacket33r14(t *testing.T) (panel.Snapshot, []byte) {
	t.Helper()
	pending := make([]panel.NativeVerdict, 0, 3)
	for i := 0; i < 3; i++ {
		// 120 runes each: three of them plus the line's fixed fields clear
		// summaryClamp=440 together (measured 2026-11-06: 45-rune ids stayed
		// under the bound, which is exactly the miss this length rules out).
		id := "ac8-clamp-33r14-" + strings.Repeat("0", 102) + string(rune('a'+i))
		pending = append(pending, panel.NativeVerdict{
			CorrelationID: id,
			Tool:          "shell.exec",
			Level:         risk.L1,
		})
	}
	cards := make([]panel.ApprovalCardView, 0, len(pending))
	for _, v := range pending {
		cards = append(cards, v.CardView())
	}
	snap := panel.NewSnapshot(cards, nil, panel.NewComposerState(risk.ModeAskEveryStep, panel.UnsetWorkspaceView(), nil, panel.MaxAttachmentBytes), time.Now())
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("marshal clamp packet: %v", err)
	}
	return snap, data
}

// TestAC8SnapshotClampBranchDropsIdsAndKeepsTheDigest33r14 is the over-bound
// case itself: the summary line this packet books MUST be the clamp branch's
// own output - ids collapsed to pending=-, depth and digest still pinning the
// packet - and the line must be short enough that the redaction layer never
// needed to cut it (the branch clamps itself rather than relying on someone
// else's rule).
func TestAC8SnapshotClampBranchDropsIdsAndKeepsTheDigest33r14(t *testing.T) {
	snap, data := ac8ClampPacket33r14(t)
	line := panelSnapshotSummary(snap, data)
	if !strings.Contains(line, "pending=-") {
		t.Errorf("a %d-byte packet was booked with ids still spelled out - the over-440 branch did not take the drop: %q", len(data), line)
	}
	if !strings.Contains(line, "depth=3") {
		t.Errorf("the clamped line must still pin how deep the queue was: %q", line)
	}
	if !strings.Contains(line, "sha256="+sha256Short(data)) {
		t.Errorf("the clamped line lost the digest that pins which packet this was: %q", line)
	}
	if len(line) > 440 {
		t.Errorf("the clamped line is %d runes: the branch clamps its own output and must not hand the redaction layer a line past summaryClamp: %q", len(line), line)
	}
}

// TestAC8SnapshotShortPacketKeepsItsIds33r14 is the negative control: the
// empty production packet (the one the entry test above books) is under the
// bound, so it MUST keep its ids. Without this half the drop ruler cannot
// tell "the branch fired" from "pending=- is just what this function writes".
func TestAC8SnapshotShortPacketKeepsItsIds33r14(t *testing.T) {
	rt := &agentRuntime{}
	p := panel.NewSnapshotPump(panel.PumpSources{Out: rt.bookPanelSnapshotBlocking})
	snap, _, err := p.Publish()
	if err != nil {
		t.Fatalf("publish: %v", err)
	}
	line := panelSnapshotSummary(snap, rt.lastSnapBytes)
	if strings.Contains(line, "pending=-") {
		t.Errorf("a short packet was booked with ids dropped - that is the over-bound shape, so this ruler cannot separate the branch from the function's own spelling: %q", line)
	}
}

// TestAC8TryEntryRefusesWhileTheBookingLockIsHeld33r14 pins the new entry's
// own behaviour: a caller that arrives while a publish holds the booking lock
// loses the race LOUDLY (reported, nothing booked), and the next caller books.
// The blocking entry keeps waiting - that difference is the whole reason the
// two names exist (ticket 33 AC#8's naming split).
func TestAC8TryEntryRefusesWhileTheBookingLockIsHeld33r14(t *testing.T) {
	rt := &agentRuntime{}
	rt.snapMu.Lock() // stand in for a publish in flight
	ok, err := rt.tryBookPanelSnapshot(panel.Snapshot{}, []byte("blocked"))
	rt.snapMu.Unlock()
	if ok || err != nil {
		t.Fatalf("a try-entry arrived while the booking lock was held and reported (%v, %v): losing the race must be reported, not booked", ok, err)
	}
	if rt.snapSeen {
		t.Fatal("the refused caller booked anyway: the miss produced a ledger line production cannot distinguish from a real publish")
	}
	ok, err = rt.tryBookPanelSnapshot(panel.Snapshot{}, []byte("after"))
	if err != nil || !ok {
		t.Fatalf("the next caller after the lock was released did not book (ok=%v, err=%v)", ok, err)
	}
	if !rt.snapSeen || string(rt.lastSnapBytes) != "after" {
		t.Fatalf("the accepted caller booked the wrong state: seen=%v bytes=%q", rt.snapSeen, rt.lastSnapBytes)
	}
}

// TestAC8BlockingEntryStillBooks33r14 keeps the production entry honest after
// the rename: Publish through the blocking entry (the one run.go wires) still
// books. This is the pair half of the try ruler above - the rename moved a
// name, not the production path.
func TestAC8BlockingEntryStillBooks33r14(t *testing.T) {
	rt := &agentRuntime{}
	p := panel.NewSnapshotPump(panel.PumpSources{Out: rt.bookPanelSnapshotBlocking})
	if _, _, err := p.Publish(); err != nil {
		t.Fatalf("publish through the blocking entry: %v", err)
	}
	if !rt.snapSeen || len(rt.lastSnapBytes) == 0 {
		t.Fatalf("the blocking entry stopped booking after the rename: seen=%v bytes=%d", rt.snapSeen, len(rt.lastSnapBytes))
	}
}
