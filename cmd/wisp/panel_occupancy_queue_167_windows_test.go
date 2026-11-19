//go:build windows

package main

// Ticket 167 AC#2 + AC#3's composition-root teeth (2026-11-19 leg), headless:
// no Win32 window is created, nothing touches the owner's ball process.
//
//   - AC#3's cancel tooth (取消排队不许把后面的挤掉): a REAL queue (approval.New,
//     the same options shape ticket 48's panel join uses), three REAL L2 cards
//     parked through the exported production entry (PendingApproval), then card
//     2 cancelled. The teeth the dispatch names are both aimed at the queue's
//     own record as LiveApprovals carries it:
//       重排 ⇒ 红  - card 3's ADMISSION number must still be 3 (not 2) after
//         card 2 leaves; a LiveApprovals that derives Seq from the live FIFO
//         place reds this test (mutation M-seq).
//       丢条 ⇒ 红  - after the cancel exactly two rows remain, card 1 and card
//         3; nothing behind the cancelled card may vanish (mutation M-drop).
//     Card 3's live FIFO place DOES move to 2 - that is the honest 当前第几条,
//     carried on its own key so the stable number never has to lie for it.
//
//   - AC#2's occupancy reader (分子分母任一未知 ⇒ 整枚未知): rt.contextOccupancy
//     returns nil while EITHER half is unmeasured - no usage event yet, or a
//     boot whose endpoint and provider declare no window. The mutation
//     "invent the denominator" (fall back to a constant window) reds
//     Test167OccupancyReaderStaysUnknownUntilBothHalvesAreMeasured.

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

func cardOf167(t *testing.T, v panel.NativeVerdict) map[string]json.RawMessage {
	t.Helper()
	data, err := json.Marshal(v.CardView())
	if err != nil {
		t.Fatalf("marshal card: %v", err)
	}
	var card map[string]json.RawMessage
	if err := json.Unmarshal(data, &card); err != nil {
		t.Fatalf("unmarshal card: %v", err)
	}
	return card
}

// Test167QueueSeqStableAcrossCancel drives AC#3 end to end: real queue, real
// cancel, the two numbers read off the SAME record the pump renders from
// (LiveApprovals -> nativeVerdictsFromGate -> CardView -> wire bytes).
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
		if it.Seq != uint64(i+1) {
			t.Errorf("row %d Seq = %d, want the admission number %d", i, it.Seq, i+1)
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

	// 重排 ⇒ 红：card 3 keeps its ADMISSION number 3. It does NOT become "the
	// second one" in the numbering - only its live FIFO place moves, and that
	// place travels on its own key.
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
	if live[1].Seq != 3 {
		t.Errorf("card 3 Seq after the cancel = %d, want 3 - the admission number must NOT renumber when an earlier card is cancelled", live[1].Seq)
	}
	if live[1].Position != 2 {
		t.Errorf("card 3 Position after the cancel = %d, want 2 - the live FIFO place IS the honest 当前第几条 and does move", live[1].Position)
	}
	if live[0].Seq != 1 || live[0].Position != 1 {
		t.Errorf("card 1 after the cancel = seq %d pos %d, want 1/1 (untouched)", live[0].Seq, live[0].Position)
	}

	// The wire: the same rows through the pump's own mapping carry both numbers
	// verbatim, so the page can show the stable number and the live place.
	verdicts := nativeVerdictsFromGate(g)
	if len(verdicts) != 2 {
		t.Fatalf("pump verdicts after the cancel = %d, want 2", len(verdicts))
	}
	var card3 panel.NativeVerdict
	for _, v := range verdicts {
		if v.CorrelationID == corrOf(3) {
			card3 = v
		}
	}
	if card3.QueueSeq != 3 || card3.Position != 2 {
		t.Errorf("pump verdict for %q = seq %d pos %d, want seq 3 pos 2", corrOf(3), card3.QueueSeq, card3.Position)
	}
	card := cardOf167(t, card3)
	if v, ok := card["seq"]; !ok || string(v) != "3" {
		t.Errorf("wire seq = %s (present %v), want 3", v, ok)
	}
	if v, ok := card["position"]; !ok || string(v) != "2" {
		t.Errorf("wire position = %s (present %v), want 2", v, ok)
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
