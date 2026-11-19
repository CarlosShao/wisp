//go:build windows

package main

// Ticket 167 AC#2 + AC#3's composition-root teeth (2026-11-19 leg), headless:
// no Win32 window is created, nothing touches the owner's ball process.
//
//   - AC#3's cancel tooth (取消排队不许把后面的挤掉): a REAL queue (approval.New,
//     the same options shape ticket 48's panel join uses), three REAL L2 cards
//     parked through the exported production entry (PendingApproval), then card
//     2 cancelled. The teeth the dispatch names are aimed at the DISPLAY
//     number the panel's pump casts: since AC#7's contract axis the stable
//     number is minted on the panel side (panel.QueueSeqLedger, one per pump,
//     first sight of a correlation id), NOT carried on LiveApprovals - so the
//     test drives the same mapping production renders through
//     (nativeVerdictsFromGate -> panel pump -> CardView -> wire bytes) and
//     pins BOTH halves:
//       重排 ⇒ 红  - card 3's DISPLAY number must still be 3 after card 2
//         leaves (the panel mint keeps it: first sighting of corr 3 cast 3;
//         the rebuild must not renumber), and its live FIFO place moves to 2.
//       丢条 ⇒ 红  - after the cancel exactly two rows remain, card 1 and card
//         3; nothing behind the cancelled card may vanish (mutation M-drop).
//
//   - AC#2's occupancy reader (分子分母任一未知 ⇒ 整枚未知): rt.contextOccupancy
//     returns nil while EITHER half is unmeasured - no usage event yet, or a
//     boot whose endpoint and provider declare no window. The mutation
//     "invent the denominator" (fall back to a constant window) reds
//     Test167OccupancyReaderStaysUnknownUntilBothHalvesAreMeasured.
//
// THE 167-d REPAIR LEG (2026-11-19 afternoon, after the morning audit): the
// audit named the panel mint's original one-call shape (Assign per card, in
// traversal order) CONSTRUCTIVELY GREEN - a ledger that renumbered every
// snapshot (casting 1..n by the current order) was indistinguishable from the
// real one through any pump-driven test. The repair separates the mint into
// two calls (castAll hands the ledger the WHOLE pending set of one snapshot,
// then each card recalls its number), and this file gains two teeth:
// Test167SeqReadTwiceFromTheOutboundPayload reads the proposition twice off
// the pump's OUTBOUND payload (before/after the cancel; corr 3 same number,
// corr 1 unchanged), and Test167MutationRenumberOnEverySnapshot injects the
// "renumber 1..n by current order" mutation through GOFLAGS -overlay (the
// tree on disk stays at the repaired shape) and MUST go red with the named
// reading (card 3's seq falls from 3 to 2); restore then greens. Both
// readings are booked in docs/evidence/s1/167-tooth-repair-d.md.

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/tools"
)

// raiseThreeL2For167 parks three REAL L2 cards on a real gate, one at a time,
// waiting for each admission so the FIFO order below measures the queue and
// not the fixture's goroutine timing (the same discipline ticket 48's suspend
// fixture states at its own raise site).
func raiseThreeL2For167(t *testing.T, g *approval.Gate, corrOf func(int) string, askDone chan tools.Answer) {
	t.Helper()
	const n = 3
	for i := 1; i <= n; i++ {
		i := i
		g.AdmitTextTask(fmt.Sprintf("task167-p%d", i))
		h := observeSpawn48(fmt.Sprintf("task167-p%d", i), corrOf(i), g)
		observe.Default.Spawn("approval-waiter", "test", nil, func(ctx context.Context) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("approval-waiter panicked (D37b): %v", r)
				}
			}()
			a, _ := h.gate.PendingApproval(ctx, h.decision)
			askDone <- a
		})
		deadline := time.Now().Add(10 * time.Second)
		for time.Now().Before(deadline) {
			if g.Queue().Depth() == i {
				break
			}
			time.Sleep(5 * time.Millisecond)
		}
		if g.Queue().Depth() != i {
			t.Fatalf("queue never reached depth %d (now %d)", i, g.Queue().Depth())
		}
	}
}

func newGate167() *approval.Gate {
	return approval.New(approval.Options{
		UI: approval.UIFuncs{
			PromptFn: func(context.Context, approval.Prompt) error { return nil },
			UpdateFn: func(context.Context, approval.Event) error { return nil },
		},
		Logf: func(string, ...any) {},
	})
}

// Test167QueueSeqStableAcrossCancel drives AC#3 end to end: real queue, real
// cancel, the panel-side mint measured through the SAME path production
// renders (nativeVerdictsFromGate -> panel pump -> CardView -> wire bytes).
func Test167QueueSeqStableAcrossCancel(t *testing.T) {
	g := newGate167()
	corrOf := func(i int) string { return fmt.Sprintf("task167-p%d#call-167-%d", i, i) }
	askDone := make(chan tools.Answer, 8)
	raiseThreeL2For167(t, g, corrOf, askDone)

	live := g.Queue().LiveApprovals()
	if len(live) != 3 {
		t.Fatalf("live rows = %d, want 3", len(live))
	}
	for i, it := range live {
		if it.CorrelationID != corrOf(i+1) || it.Position != i+1 {
			t.Errorf("row %d = %q pos %d, want %q pos %d", i, it.CorrelationID, it.Position, corrOf(i+1), i+1)
		}
	}

	// The panel-side mint, over the mapping production renders from: all three
	// ids cast in first-seen order 1/2/3 - card 3's number is 3 BEFORE the
	// cancel, which is the fact the post-cancel half pins.
	verdicts := nativeVerdictsFromGate(g)
	if len(verdicts) != 3 {
		t.Fatalf("pump verdicts = %d, want 3", len(verdicts))
	}
	pump := panel.NewSnapshotPump(panel.PumpSources{
		Verdicts: func() []panel.NativeVerdict { return verdicts },
	})
	cardSeq := func(t *testing.T, p *panel.SnapshotPump, corr string) (uint64, int) {
		t.Helper()
		snap := p.Snapshot()
		for _, c := range snap.Pending {
			if c.CorrelationID == corr {
				return c.Seq, c.Position
			}
		}
		t.Fatalf("card %q not in the packet's pending section", corr)
		return 0, 0
	}
	for i := 1; i <= 3; i++ {
		if seq, pos := cardSeq(t, pump, corrOf(i)); seq != uint64(i) || pos != i {
			t.Errorf("card %d through the pump = seq %d pos %d, want the panel-cast seq %d and FIFO pos %d", i, seq, pos, i, i)
		}
	}

	// The cancel: card 2 (the middle one) leaves through the production native
	// refusal route - exactly what a user's dismiss does.
	if err := g.DecideFromNative(context.Background(), approval.Request{
		CorrelationID: corrOf(2), Allow: false, Reason: "票 167 取消牙：取消第 2 条",
	}); err != nil {
		t.Fatalf("cancel %q: %v", corrOf(2), err)
	}
	select {
	case a := <-askDone:
		if a != tools.AnswerReject {
			t.Errorf("cancelled card answered %s, want reject", a)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the cancelled card's parked ask never woke")
	}

	// 重排 ⇒ 红：after the cancel exactly two rows remain and NOTHING behind
	// the cancelled card was squeezed out.
	live = g.Queue().LiveApprovals()
	if len(live) != 2 {
		t.Fatalf("live rows after the cancel = %d, want exactly 2 (card 2 gone, NOTHING behind it squeezed out)", len(live))
	}
	if live[0].CorrelationID != corrOf(1) {
		t.Errorf("row 0 after the cancel = %q, want %q (card 1 untouched)", live[0].CorrelationID, corrOf(1))
	}
	if live[1].CorrelationID != corrOf(3) {
		t.Fatalf("row 1 after the cancel = %q, want %q (card 3 must still be on the queue, not dropped)", live[1].CorrelationID, corrOf(3))
	}

	// The same pump's packet now: card 3 keeps its DISPLAY number 3 (the panel
	// mint never renumbers a seen id) while its live FIFO place moved to 2 -
	// the honest 当前第几条 travels on its own key.
	verdicts = nativeVerdictsFromGate(g)
	if len(verdicts) != 2 {
		t.Fatalf("pump verdicts after the cancel = %d, want 2", len(verdicts))
	}
	if seq, pos := cardSeq(t, pump, corrOf(1)); seq != 1 || pos != 1 {
		t.Errorf("card 1 after the cancel = seq %d pos %d, want 1/1 (untouched)", seq, pos)
	}
	if seq, pos := cardSeq(t, pump, corrOf(3)); seq != 3 || pos != 2 {
		t.Errorf("card 3 after the cancel = seq %d pos %d, want seq 3 (NOT renumbered to 2) pos 2", seq, pos)
	}

	// And the WIRE: the same packet's JSON carries both numbers verbatim, so
	// the page can show the stable number and the live place.
	data, err := pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var packet struct {
		Pending []map[string]json.RawMessage `json:"pending"`
	}
	if err := json.Unmarshal(data, &packet); err != nil {
		t.Fatal(err)
	}
	if len(packet.Pending) != 2 {
		t.Fatalf("wire pending cards = %d, want 2", len(packet.Pending))
	}
	var card3 map[string]json.RawMessage
	for _, c := range packet.Pending {
		if string(c["correlationId"]) == fmt.Sprintf("%q", corrOf(3)) {
			card3 = c
		}
	}
	if card3 == nil {
		t.Fatalf("wire pending has no card %q", corrOf(3))
	}
	if v, ok := card3["seq"]; !ok || string(v) != "3" {
		t.Errorf("wire seq = %s (present %v), want 3", v, ok)
	}
	if v, ok := card3["position"]; !ok || string(v) != "2" {
		t.Errorf("wire position = %s (present %v), want 2", v, ok)
	}
}

// Test167SeqReadTwiceFromTheOutboundPayload is the independent support tooth
// the repair leg adds (2026-11-19 167-d): the stable-number proposition read
// TWICE from the pump's OUTBOUND payload (the Snapshot the page is rendered
// from), never from any ledger field - the ledger's map and counter are
// private, and this case cannot see them even if it wanted to. The two reads
// bracket the cancel: before card 2 leaves, corr 3 reads A; after, corr 3
// reads the SAME A and corr 1 still reads its own number.
func Test167SeqReadTwiceFromTheOutboundPayload(t *testing.T) {
	g := newGate167()
	corrOf := func(i int) string { return fmt.Sprintf("task167-p%d#call-167-%d", i, i) }
	askDone := make(chan tools.Answer, 8)
	raiseThreeL2For167(t, g, corrOf, askDone)

	verdicts := nativeVerdictsFromGate(g)
	if len(verdicts) != 3 {
		t.Fatalf("pump verdicts = %d, want 3", len(verdicts))
	}
	pump := panel.NewSnapshotPump(panel.PumpSources{
		Verdicts: func() []panel.NativeVerdict { return verdicts },
	})

	// The outbound-payload reader: one snapshot's pending section, keyed by
	// correlation id. Seq here is whatever the pump put on the packet, which is
	// the whole claim - the page's number is this number.
	payloadSeq := func(t *testing.T, p *panel.SnapshotPump, corr string) uint64 {
		t.Helper()
		for _, c := range p.Snapshot().Pending {
			if c.CorrelationID == corr {
				return c.Seq
			}
		}
		t.Fatalf("card %q not in the outbound payload's pending section", corr)
		return 0
	}

	// Read 1, BEFORE the cancel: corr 3's display number A comes off the
	// payload, not off the queue and not off any ledger field.
	seqA := payloadSeq(t, pump, corrOf(3))
	if seqA == 0 {
		t.Fatalf("card 3's payload seq before the cancel = 0: the outbound card travels without a number, nothing downstream is pinned")
	}
	if seq1 := payloadSeq(t, pump, corrOf(1)); seq1 == 0 {
		t.Fatalf("card 1's payload seq before the cancel = 0, nothing downstream is pinned")
	}

	// The cancel: card 2 (the middle one) leaves through the production native
	// refusal route - exactly what a user's dismiss does.
	if err := g.DecideFromNative(context.Background(), approval.Request{
		CorrelationID: corrOf(2), Allow: false, Reason: "票 167 出向载荷牙：取消第 2 条",
	}); err != nil {
		t.Fatalf("cancel %q: %v", corrOf(2), err)
	}
	select {
	case a := <-askDone:
		if a != tools.AnswerReject {
			t.Errorf("cancelled card answered %s, want reject", a)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the cancelled card's parked ask never woke")
	}
	if live := g.Queue().LiveApprovals(); len(live) != 2 {
		t.Fatalf("live rows after the cancel = %d, want 2", len(live))
	}

	// Read 2, AFTER the cancel, from a fresh payload of the SAME pump: corr 3
	// must read the SAME number A (取消不重排, read off the packet the page
	// gets) and corr 1 must still read its own number.
	if seq3 := payloadSeq(t, pump, corrOf(3)); seq3 != seqA {
		t.Errorf("card 3 payload seq after the cancel = %d, want the SAME %d it read before (renumbering a survivor is the lie AC#3 forbids)", seq3, seqA)
	}
	if seq1 := payloadSeq(t, pump, corrOf(1)); seq1 == 0 {
		t.Errorf("card 1 payload seq after the cancel = 0, want its own cast number unchanged")
	}
}

// Test167MutationRenumberOnEverySnapshot is the repair leg's mutation tooth
// (2026-11-19 167-d): with the ledger mutated to renumber 1..n by the current
// order on every snapshot, the end-to-end cancel read MUST go red with the
// named reading (card 3's seq falls from 3 to 2) - proving the packet-borne
// teeth above can actually fail. The mutation is injected via GOFLAGS
// -overlay (the same shape evidence file 219-ac5-reason-findback.md §8
// books), so the tree on disk stays at the repaired shape for every run;
// restored, this case is a plain green tooth over the packet. Both readings
// are booked in docs/evidence/s1/167-tooth-repair-d.md.
func Test167MutationRenumberOnEverySnapshot(t *testing.T) {
	g := newGate167()
	corrOf := func(i int) string { return fmt.Sprintf("task167-p%d#call-167-%d", i, i) }
	askDone := make(chan tools.Answer, 8)
	raiseThreeL2For167(t, g, corrOf, askDone)

	verdicts := nativeVerdictsFromGate(g)
	if len(verdicts) != 3 {
		t.Fatalf("pump verdicts = %d, want 3", len(verdicts))
	}
	pump := panel.NewSnapshotPump(panel.PumpSources{
		Verdicts: func() []panel.NativeVerdict { return verdicts },
	})
	cardSeq := func(t *testing.T, p *panel.SnapshotPump, corr string) uint64 {
		t.Helper()
		for _, c := range p.Snapshot().Pending {
			if c.CorrelationID == corr {
				return c.Seq
			}
		}
		t.Fatalf("card %q not in the packet's pending section", corr)
		return 0
	}
	if seq := cardSeq(t, pump, corrOf(3)); seq != 3 {
		t.Fatalf("card 3 payload seq before the cancel = %d, want 3", seq)
	}

	if err := g.DecideFromNative(context.Background(), approval.Request{
		CorrelationID: corrOf(2), Allow: false, Reason: "票 167 变异载具：取消第 2 条",
	}); err != nil {
		t.Fatalf("cancel %q: %v", corrOf(2), err)
	}
	select {
	case a := <-askDone:
		if a != tools.AnswerReject {
			t.Errorf("cancelled card answered %s, want reject", a)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the cancelled card's parked ask never woke")
	}

	verdicts = nativeVerdictsFromGate(g)
	if len(verdicts) != 2 {
		t.Fatalf("pump verdicts after the cancel = %d, want 2", len(verdicts))
	}
	// THE MUTATION-SENSITIVE READ: under the renumber-every-snapshot mutation
	// the two survivors renumber to 1 and 2 by their current order, so card 3
	// reads 2 here and this test - the same proposition the packet teeth pin -
	// goes RED. Under the repaired ledger it reads 3 and the test is GREEN.
	if seq := cardSeq(t, pump, corrOf(3)); seq != 3 {
		t.Errorf("card 3 payload seq after the cancel = %d, want 3 (the mutation made it fall to 2: 取消不重排 is broken and this tooth catches it)", seq)
	}
}

// Test167OccupancyReaderStaysUnknownUntilBothHalvesAreMeasured drives AC#2's
// reader: nil until the numerator (a usage event) AND the denominator (a
// resolved window) both exist, then both numbers verbatim.
func Test167OccupancyReaderStaysUnknownUntilBothHalvesAreMeasured(t *testing.T) {
	rt := &agentRuntime{
		runFeed:  panel.NewRunFeed(),
		endpoint: llm.Endpoint{Model: "m", ContextWindow: 0},
	}

	// Nothing measured yet: nil is the named unknown answer.
	if got := rt.contextOccupancy(); got != nil {
		t.Fatalf("reader with no usage and no window = %+v, want nil (unknown is named, not zero)", got)
	}

	// A numerator arrives; the denominator is still unknown. Either half
	// unknown ⇒ the whole dimension unknown: the reader must NOT hand the pump
	// a reading with an invented window (mutation: constant fallback reds here).
	rt.runFeed.Usage(500, 40, 20, 100, "CNY")
	if got := rt.contextOccupancy(); got != nil {
		t.Fatalf("reader with a usage event but an unknown window = %+v, want nil (no invented denominator)", got)
	}

	// The window resolves (the way a real boot's endpoint carries it); now both
	// halves exist and travel verbatim. UsedTokens is in+out of the SAME
	// cumulative row, carried not re-derived.
	rt.endpoint.ContextWindow = 128000
	got := rt.contextOccupancy()
	if got == nil {
		t.Fatal("reader with both halves measured = nil, want the reading")
	}
	if got.UsedTokens != 540 || got.WindowTokens != 128000 {
		t.Errorf("reading = used %d window %d, want 540/128000 (the two measured halves, verbatim)", got.UsedTokens, got.WindowTokens)
	}

	// And through the pump: the wire says known=true with the numbers.
	data, err := panel.NewSnapshotPump(panel.PumpSources{Occupancy: rt.contextOccupancy}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var wire map[string]json.RawMessage
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatal(err)
	}
	var cs map[string]json.RawMessage
	if err := json.Unmarshal(wire["composer"], &cs); err != nil {
		t.Fatal(err)
	}
	if string(cs["contextOccupancyKnown"]) != "true" {
		t.Errorf("wire contextOccupancyKnown = %s, want true", cs["contextOccupancyKnown"])
	}
	if string(cs["contextTokens"]) != "540" || string(cs["contextWindow"]) != "128000" {
		t.Errorf("wire numbers = %s / %s, want 540 / 128000", cs["contextTokens"], cs["contextWindow"])
	}

	// The reader-gated shape: a pump assembled without the reader says unknown
	// by name (composer's own flag), never by inventing state.
	data, err = panel.NewSnapshotPump(panel.PumpSources{}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatal(err)
	}
	var bare map[string]json.RawMessage
	if err := json.Unmarshal(wire["composer"], &bare); err != nil {
		t.Fatal(err)
	}
	if string(bare["contextOccupancyKnown"]) != "false" {
		t.Errorf("readerless pump contextOccupancyKnown = %s, want the named unknown false", bare["contextOccupancyKnown"])
	}
	if _, ok := bare["contextTokens"]; ok {
		t.Errorf("readerless pump carries %s on the wire - absent is the honest unknown, a 0 is not", bare["contextTokens"])
	}
	if _, ok := bare["contextWindow"]; ok {
		t.Errorf("readerless pump carries %s on the wire - absent is the honest unknown, an invented window is not", bare["contextWindow"])
	}
}
