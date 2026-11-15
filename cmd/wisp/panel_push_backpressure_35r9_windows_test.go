//go:build windows

// Package main test - ticket 35, cell ":49" Backpressure, leg 35-g3 (writer).
//
// WHY THIS FILE EXISTS, AND WHY IT IS NOT THE 35r8 FILE. The four 35r8 rulers in
// internal/panel/backpressure_bound_35r8_test.go measure the StreamLog and one
// deliberately unbounded in-memory sink. They are good rulers and they pass, but
// they never touch the thing the ticket's backpressure cell names for the OUTBOUND
// hop: not one of them calls Dispatch, reads dispatchq, or causes one push to be
// queued. 35-recon measured the consequence (docs/evidence/s1/35-recon-20261115.md
// §格⑤: "35r8 四枚现跑全绿 ... 量的是 StreamLog；⛔ 没有一枚尺打在推队列上"), and
// cmd/wisp/panel_pump.go:502-505 says the same thing in production's own words. So
// these are FOUR NEW rulers over the push queue itself (plus two the cell needs
// besides: the on-disk reading and the heap arithmetic), and the old four are named
// in each test's header as the ones they do not replace, not as ones they weaken.
// The old file is NOT edited by this leg (internal/panel is read-only this wave);
// nothing here loosens an assertion over there.
//
// THE MEASUREMENT ERROR THIS FILE MUST NOT INHERIT. The 35r8 harness stores the
// heap reading in a uint64 and computes `after - before`. When the GC returns more
// than the flood allocated, that subtraction UNDERFLOWS to ~2^64, and the shared
// predicate (`m.heapDelta > heapCap`) then reads it as "far over the cap" - a false
// red produced by the ruler, not by the subject. It stayed invisible because only
// the 4x reading was asserted there while the 1x heap number was printed and never
// judged (backpressure_bound_35r8_test.go:64, :207, :237, and the m1-only prints at
// :303 / :371 / :411). This file does it the other way round:
//   - the delta is SIGNED (int64) and an underflow is an explicit, recorded fact
//     (heapDelta35r9 returns `underflow bool`), never a huge number;
//   - the predicate refuses to call an underflowed reading "over cap";
//   - BOTH the 1x and the 4x reading are asserted on every dimension, including
//     the heap one - no dimension is print-only here.
// TestPanelPushHeapDeltaUnderflowIsNotReadAsOverCap35r9 pins that arithmetic
// directly, so a regression to the unsigned shape reddens a named test instead of
// turning up as a mysterious red half a year later.
//
// THE SHAPE UNDER TEST (production, cmd/wisp/panel_host_windows.go): evalPanelJS
// keeps ONE pending whole-state sentence and supersedes the older one; the page's
// queue therefore holds one delivery closure, not one per publish. The flood is
// always driven through the REAL evalPanelJS - the production queue - and the
// positive control drives the SAME harness and the SAME predicate over a queue
// with the bound taken out on purpose, so the predicate is proven to discriminate.
//
// WHICH SENTENCE OF THE CELL THIS IMPLEMENTS, AND WHICH TWO IT MAY NOT. The cell's
// original three thirds - "merges", "no unbounded memory (heap cap asserted)",
// "content integrity kept" - were re-cut by ticket 197 leg B ("Overflow TRUNCATES
// and never merges"), and the orchestrator's 10-08 11:3x ruling says so in the cell
// itself (ticket :51). What still scores after that ruling is the SECOND third only.
// So:
//   - this leg implements "no unbounded memory" for the PUSH QUEUE, and it is the
//     only third any ruler here claims;
//   - NOTHING here merges two packets into one sentence, and the positive control
//     would not catch a merge because it does not ask that question - the
//     "never merges" rule is held where it belongs (internal/panel/pump.go:632 and
//     its named test), and these rulers only ensure this hop never becomes a second
//     place that merges;
//   - "content integrity kept", read as "an event is never lost", is NOT claimed
//     and is FALSE by design here: intermediate whole-state packets are superseded
//     and counted (pushCoalesced). What is claimed is narrower and checkable - the
//     delivered packet is one WHOLE packet, never a slice and never two packets
//     welded together, and the newest one wins. Read TestPanelPushQueueDelivers
//     TheNewestWholePacket35r9's own header for that distinction.

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"unsafe"

	webview2 "github.com/jchv/go-webview2"
)

// --------------------------------------------------------------- the consumer --

// blockedPage35r9 is a WebView2 stand-in whose "panel thread" never runs the
// closures it is handed until the test calls resume. That is the ticket's "blocked
// consumer" made literal: Dispatch still accepts (the pre-fix bug grew exactly
// here), and nothing reaches Eval until the block is lifted.
type blockedPage35r9 struct {
	mu      sync.Mutex
	pending []func()
	last    string
	evals   int
}

var _ webview2.WebView = (*blockedPage35r9)(nil)

func newBlockedPage35r9() *blockedPage35r9 { return &blockedPage35r9{} }

func (p *blockedPage35r9) Dispatch(f func()) {
	if f == nil {
		return
	}
	p.mu.Lock()
	p.pending = append(p.pending, f)
	p.mu.Unlock()
}

func (p *blockedPage35r9) Eval(js string) {
	p.mu.Lock()
	p.last = js
	p.evals++
	p.mu.Unlock()
}

func (p *blockedPage35r9) queued() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.pending)
}

func (p *blockedPage35r9) delivered() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.evals
}

func (p *blockedPage35r9) lastEval() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.last
}

// resume lifts the block: every queued closure runs, in order, on the test's
// goroutine - which is the panel thread for the purpose of this ruler.
func (p *blockedPage35r9) resume() int {
	p.mu.Lock()
	q := append([]func(){}, p.pending...)
	p.pending = nil
	p.mu.Unlock()
	for _, f := range q {
		f()
	}
	return len(q)
}

func (p *blockedPage35r9) Run()                                  {}
func (p *blockedPage35r9) Terminate()                            {}
func (p *blockedPage35r9) Destroy()                              {}
func (p *blockedPage35r9) Window() unsafe.Pointer                { return nil }
func (p *blockedPage35r9) SetTitle(title string)                 {}
func (p *blockedPage35r9) SetSize(w, h int, hint webview2.Hint)  {}
func (p *blockedPage35r9) Navigate(url string)                   {}
func (p *blockedPage35r9) SetHtml(html string)                   {}
func (p *blockedPage35r9) Init(js string)                        {}
func (p *blockedPage35r9) Bind(name string, f interface{}) error { return nil }

// attachPanelWindow35r9 points a manager at the stand-in the way bringUp does after
// a successful create (m.w = w, created = true), with no window and no WebView2.
func attachPanelWindow35r9(t *testing.T, m *PanelManager, w webview2.WebView) {
	t.Helper()
	m.mu.Lock()
	m.w = w
	m.created = true
	m.mu.Unlock()
}

// ---------------------------------------------------------------- the two queues --

// pushQueue35r9 is the seam the flood is driven through: the production queue and
// the deliberately unbounded control both answer it, which is the only reason one
// harness and one predicate can judge both (the shape 35r8's control uses too).
type pushQueue35r9 interface {
	push(js string)
	measure() pushQueueReading35r9
}

type pushQueueReading35r9 struct {
	waiting     int   // whole sentences waiting to be handed to the page
	wakeups     int   // delivery closures already sitting in the page's queue
	waitingSize int   // bytes those waiting sentences pin
	coalesced   int64 // pushes superseded before the page saw them
	skipped     int64 // pushes that had no window at all
}

// panelPushQueue35r9 is the PRODUCTION queue: every push goes through the real
// evalPanelJS, and every number comes out of the production counters.
type panelPushQueue35r9 struct {
	m       *PanelManager
	blocked *blockedPage35r9
}

func (q panelPushQueue35r9) push(js string) { q.m.evalPanelJS(js) }

func (q panelPushQueue35r9) measure() pushQueueReading35r9 {
	s := q.m.PushStats()
	return pushQueueReading35r9{
		waiting:     s.Pending,
		wakeups:     q.blocked.queued(),
		waitingSize: s.PendingBytes,
		coalesced:   s.Coalesced,
		skipped:     s.Skipped,
	}
}

// unboundedPushQueue35r9 is the pre-fix shape with the bound taken out on purpose:
// every push is handed to Dispatch immediately - exactly what evalPanelJS did before
// this leg - so the page queue grows one closure and one whole sentence per event.
// It exists so the positive control can prove the predicate discriminates between
// "bounded" and "the same flood, no bound".
type unboundedPushQueue35r9 struct {
	page *blockedPage35r9
	mu   sync.Mutex
	held []string
}

func (u *unboundedPushQueue35r9) push(js string) {
	u.mu.Lock()
	u.held = append(u.held, js)
	u.mu.Unlock()
	// The closure captures the whole sentence: that is the memory the pre-fix code
	// pinned per publish, and it is what the control has to reproduce.
	p := u.page
	p.Dispatch(func() { p.Eval(js) })
}

func (u *unboundedPushQueue35r9) measure() pushQueueReading35r9 {
	u.mu.Lock()
	defer u.mu.Unlock()
	size := 0
	for _, s := range u.held {
		size += len(s)
	}
	return pushQueueReading35r9{
		waiting:     len(u.held),
		wakeups:     u.page.queued(),
		waitingSize: size,
	}
}

// ------------------------------------------------------------------ measuring --

// pushSentence35r9 is one complete push sentence: the named global the page's
// receiver is called through, wrapping a packet whose tail is the sequence number.
// Every sentence is WHOLE - none of the rulers below ever store or compare a slice
// of one - which is what makes "the delivered packet is not a fragment and not two
// packets welded together" measurable on this hop.
func pushSentence35r9(seq, padBytes int) string {
	return panelSnapshotGlobal + `({"seq":` + fmt.Sprint(seq) + `,"pad":"` + strings.Repeat("x", padBytes) + `"})`
}

// pushFloodBytes35r9 is the exact number of bytes a flood of attempts sentences
// plants, because each sentence's length depends on its sequence number's width
// (seq 9 is one byte shorter than seq 10). A control that compared against
// attempts times one sentence's length would be a few bytes short and could call a
// faithful queue unfaithful - or hide a real overrun in that slack.
func pushFloodBytes35r9(attempts, padBytes int) int {
	total := 0
	for i := 0; i < attempts; i++ {
		total += len(pushSentence35r9(i, padBytes))
	}
	return total
}

// heapInUse35r9 reads live heap, GC'd first so the number is retention and not
// allocator slack (the same reading 35r8 took; only the arithmetic below differs).
func heapInUse35r9() uint64 {
	runtime.GC()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms.HeapAlloc
}

// heapDelta35r9 converts two HeapAlloc readings into a SIGNED delta plus an
// explicit underflow flag. This is the whole fix for the 35r8 measurement error: an
// unsigned subtraction turns "the GC gave memory back" into ~2^64 and the shared
// predicate reads that as "far over the cap". A negative movement is a READING of
// its own ("the flood pinned less than the baseline"), and it is reported as such.
func heapDelta35r9(before, after uint64) (int64, bool) {
	if after < before {
		return 0, true
	}
	return int64(after - before), false
}

// pushReading35r9 is one flood's reading, in the units the cell is scored in: how
// many whole sentences the queue still holds, how many bytes they pin, how much the
// page queue itself grew, and what the process heap did (signed, with the underflow
// flag carried next to it instead of folded into the number).
type pushReading35r9 struct {
	attempts      int
	sentenceBytes int
	waiting       int
	wakeups       int
	waitingSize   int
	coalesced     int64
	skipped       int64
	heapDelta     int64
	heapUnderflow bool
	newest        string
}

// pushCap35r9 is the cap every ruler in this file uses. Every term comes from the
// BOUND and none from the flood size: one waiting sentence, one wake-up closure in
// the page queue, and at most one whole sentence's bytes pinned. The heap term is
// measurement slack with the same shape 35r8 used (4x the packet plus 1 MiB for map
// and GC bookkeeping), not a loosened bound.
type pushCap35r9 struct {
	waiting int
	wakeups int
	pinned  int64
	heap    int64
}

func pushCapFor35r9(sentenceBytes int) pushCap35r9 {
	return pushCap35r9{
		waiting: panelPushPendingCap,
		// One wake-up per burst, independent of the cap: the queue schedules a
		// delivery only when none is in flight, so N pushes share one closure.
		wakeups: 1,
		pinned:  int64(sentenceBytes),
		heap:    int64(sentenceBytes)*4 + (1 << 20),
	}
}

// overPushCap35r9 is the one predicate this file's tests share. It returns the
// reading instead of failing the test so the control can ask the same question about
// an uncapped queue and require the opposite answer. An underflowed heap reading is
// NEVER "over cap": it is reported as unreadable and the other dimensions still
// judge the reading.
func overPushCap35r9(m pushReading35r9, bound pushCap35r9) string {
	var hits []string
	if m.waiting > bound.waiting {
		hits = append(hits, fmt.Sprintf("waiting sentences %d > the bound %d", m.waiting, bound.waiting))
	}
	if m.wakeups > bound.wakeups {
		hits = append(hits, fmt.Sprintf("page-queue wake-ups %d > the bound %d (a wake-up per event is the pre-fix shape)", m.wakeups, bound.wakeups))
	}
	if int64(m.waitingSize) > bound.pinned {
		hits = append(hits, fmt.Sprintf("pinned sentence bytes %d > one whole sentence %d", m.waitingSize, bound.pinned))
	}
	if !m.heapUnderflow && m.heapDelta > bound.heap {
		hits = append(hits, fmt.Sprintf("live heap +%d bytes > the cap %d", m.heapDelta, bound.heap))
	}
	if len(hits) == 0 {
		return ""
	}
	return strings.Join(hits, "; ")
}

// floodPush35r9 drives attempts whole sentences through a queue while its consumer
// is blocked, and reports what that pinned. It asserts the block is real (the page
// ran nothing) instead of trusting a comment, and it never reads the page's queue
// after resuming - the reading is taken at the moment of maximum retention.
func floodPush35r9(t *testing.T, q pushQueue35r9, page *blockedPage35r9, attempts, padBytes int) pushReading35r9 {
	t.Helper()
	before := heapInUse35r9()
	newest := ""
	for i := 0; i < attempts; i++ {
		js := pushSentence35r9(i, padBytes)
		q.push(js)
		newest = js
	}
	after := heapInUse35r9()

	if got := page.delivered(); got != 0 {
		t.Fatalf("the blocked consumer ran %d delivery/deliveries during a %d-push flood: this harness is not measuring a blocked consumer", got, attempts)
	}
	m := q.measure()
	delta, under := heapDelta35r9(before, after)
	return pushReading35r9{
		attempts:      attempts,
		sentenceBytes: len(newest),
		waiting:       m.waiting,
		wakeups:       m.wakeups,
		waitingSize:   m.waitingSize,
		coalesced:     m.coalesced,
		skipped:       m.skipped,
		heapDelta:     delta,
		heapUnderflow: under,
		newest:        newest,
	}
}

// --------------------------------------------------------------- log on disk --

// logRecordsUnder35r9 reads every JSONL record this test's data root wrote, so a
// case can assert what an operator with only the disk would find.
func logRecordsUnder35r9(t *testing.T, dir string) []map[string]any {
	t.Helper()
	var out []map[string]any
	files := jsonlFilesUnder(t, dir)
	if len(files) == 0 {
		t.Fatalf("no wisp-*.jsonl under %s: the sink wrote nothing to disk, so no counter can be read there", dir)
	}
	for _, p := range files {
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("reading %s: %v", p, err)
		}
		for _, line := range strings.Split(string(b), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || !strings.HasPrefix(line, "{") {
				continue
			}
			var rec map[string]any
			if err := json.Unmarshal([]byte(line), &rec); err != nil {
				t.Fatalf("record in %s is not one JSON object: %v (%q)", p, err, line)
			}
			out = append(out, rec)
		}
	}
	return out
}

// findLastRecord35r9 returns the LAST record with this message, which is the one the
// case under test just wrote. It matters because the log sink replays records that
// were buffered before it was installed (ticket 130's early-log buffer), so earlier
// tests' records land in this file too - a first-match reader would assert them.
func findLastRecord35r9(recs []map[string]any, msg string) map[string]any {
	var found map[string]any
	for _, rec := range recs {
		if got, _ := rec["msg"].(string); got == msg {
			found = rec
		}
	}
	return found
}

func countRecords35r9(recs []map[string]any, msg string) int {
	n := 0
	for _, rec := range recs {
		if got, _ := rec["msg"].(string); got == msg {
			n++
		}
	}
	return n
}

func intField35r9(t *testing.T, rec map[string]any, key string) int64 {
	t.Helper()
	switch v := rec[key].(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			t.Fatalf("field %q is not an integer: %v", key, err)
		}
		return n
	default:
		t.Fatalf("field %q is %T (%v), want a number", key, rec[key], rec[key])
		return 0
	}
}

// -------------------------------------------------------------------- rulers --

// TestPanelPushQueueFloodUnderBlockedConsumerStaysBounded35r9 is this file's
// replacement for 35r8's StreamLog boundedness rulers
// (TestStreamLogFanOutFloodKeepsRetentionUnderCap35r8 and
// TestStreamLogFloodBelowKeyBoundIsBounded35r8): same "flood under a blocked
// consumer, cap asserted", but the subject is the PUSH QUEUE. 400 and 1600 pushes
// of the same 4 KiB packet must pin one sentence each - not 400 and not 1600 - and
// the retention must not scale between the two floods.
func TestPanelPushQueueFloodUnderBlockedConsumerStaysBounded35r9(t *testing.T) {
	const padBytes = 4096
	// The cap's byte unit is the LONGEST sentence the flood can produce (seq has
	// more digits at 1600 than at 0), so the cap never runs short by a digit.
	bound := pushCapFor35r9(len(pushSentence35r9(1600-1, padBytes)))

	run := func(attempts int) (pushReading35r9, *blockedPage35r9, panelPushQueue35r9) {
		page := newBlockedPage35r9()
		m := NewPanelManager(nil, nil, "")
		attachPanelWindow35r9(t, m, page)
		q := panelPushQueue35r9{m: m, blocked: page}
		return floodPush35r9(t, q, page, attempts, padBytes), page, q
	}

	m1, page1, q1 := run(400)
	if bad := overPushCap35r9(m1, bound); bad != "" {
		t.Errorf("400 pushes under a blocked consumer are not bounded: %s", bad)
	}
	m4, page4, q4 := run(1600)
	if bad := overPushCap35r9(m4, bound); bad != "" {
		t.Errorf("1600 pushes under the same blocked consumer are not bounded: %s", bad)
	}
	if got := q4.m.PushStats().Cap; got != panelPushPendingCap {
		t.Fatalf("PushStats().Cap = %d, want the production bound %d: the ruler's cap must be the shipping one", got, panelPushPendingCap)
	}

	// The bound must be exact, not "small": one whole sentence, one wake-up. A
	// number below that would mean the newest packet was thrown away, which is the
	// one thing this queue may not do.
	for name, m := range map[string]pushReading35r9{"1x": m1, "4x": m4} {
		if m.waiting != panelPushPendingCap {
			t.Errorf("%s: %d sentence(s) waiting, want exactly the bound %d - latest wins, so the newest is always held", name, m.waiting, panelPushPendingCap)
		}
		if m.wakeups != 1 {
			t.Errorf("%s: %d delivery wake-up(s) in the page queue, want exactly 1 - the whole flood must share one wake-up, not one per event (the pre-fix shape)", name, m.wakeups)
		}
		if m.waitingSize != m.sentenceBytes {
			t.Errorf("%s: pinned %d bytes, want one whole sentence %d", name, m.waitingSize, m.sentenceBytes)
		}
		if want := int64(m.attempts - 1); m.coalesced != want {
			t.Errorf("%s: coalesced = %d, want %d - every superseded push must be counted, or the loss is silent", name, m.coalesced, want)
		}
	}
	if m4.heapUnderflow || m1.heapUnderflow {
		// Recorded, never passed off as a huge number: this is the 35r8 hole, and
		// it is a fact here rather than a reading. It is NOT a failure - the GC may
		// legitimately return more than one blocked-consumer flood allocated - but
		// the dimension is then reported unavailable instead of "far over cap".
		t.Logf("heap reading underflowed (1x underflow=%v, 4x underflow=%v): recorded as an underflow, not as a reading; the sentence/byte/wake-up dimensions above still judge the bound", m1.heapUnderflow, m4.heapUnderflow)
	}

	// Resume both consumers: exactly one delivery each, carrying the newest packet,
	// and the counters agree afterwards.
	for name, tc := range map[string]struct {
		page    *blockedPage35r9
		q       panelPushQueue35r9
		reading pushReading35r9
	}{"1x": {page1, q1, m1}, "4x": {page4, q4, m4}} {
		if ran := tc.page.resume(); ran != 1 {
			t.Errorf("%s: the page ran %d queued delivery/deliveries, want exactly 1 - the queue coalesced, it did not re-post", name, ran)
		}
		if got := tc.page.lastEval(); got != tc.reading.newest {
			t.Errorf("%s: the page received %q, want the newest whole sentence %q", name, got, tc.reading.newest)
		}
		st := tc.q.m.PushStats()
		if st.Pending != 0 || st.Delivered != 1 || st.Coalesced != int64(tc.reading.attempts-1) {
			t.Errorf("%s: after resume PushStats = %+v, want Pending 0 / Delivered 1 / Coalesced %d", name, st, tc.reading.attempts-1)
		}
	}
	t.Logf("35-g3 push queue reading: 1x waiting=%d bytes=%d wakeups=%d coalesced=%d heap+%dB | 4x waiting=%d bytes=%d wakeups=%d coalesced=%d heap+%dB | bound waiting=%d pinned=%d heap=%d",
		m1.waiting, m1.waitingSize, m1.wakeups, m1.coalesced, m1.heapDelta,
		m4.waiting, m4.waitingSize, m4.wakeups, m4.coalesced, m4.heapDelta,
		bound.waiting, bound.pinned, bound.heap)
}

// TestPanelPushQueueDeliversTheNewestWholePacket35r9 is the ruler for WHAT the
// queue hands over, and it is deliberately NOT the cell's superseded "content
// integrity kept" sentence. Under the orchestrator's 10-08 11:3x ruling (ticket :51)
// that sentence no longer scores: ticket 197 leg B re-cut it to "Overflow TRUNCATES
// and never merges", and a leg that claimed "no event is ever lost" would be
// claiming the opposite of the shipping rule. What this ruler does claim is
// narrower, and it is what the bound owes:
//   - the delivered packet is ONE whole sentence - the bytes that were pushed, byte
//     for byte, never a slice of one and never two packets welded together;
//   - it is the NEWEST one (the queue supersedes the oldest, never the newest);
//   - the intermediate whole-state packets it superseded are COUNTED as losses
//     (pushCoalesced), because they really are lost at the event level.
//
// It replaces 35r8's TestStreamLogDroppedNamingLedgerIsBounded35r8 only in the
// sense of being this hop's "what is retained" ruler; the ledger ruler's subject is
// untouched, and nothing here is allowed to go green by concatenating or rewriting
// sentences (that would be the merge rule 197 leg B forbids).
func TestPanelPushQueueDeliversTheNewestWholePacket35r9(t *testing.T) {
	const attempts, padBytes = 200, 4096
	page := newBlockedPage35r9()
	m := NewPanelManager(nil, nil, "")
	attachPanelWindow35r9(t, m, page)
	q := panelPushQueue35r9{m: m, blocked: page}

	r := floodPush35r9(t, q, page, attempts, padBytes)
	if r.waiting != 1 || r.waitingSize != r.sentenceBytes {
		t.Fatalf("the queue holds %d sentence(s) / %d bytes, want one whole sentence (%d bytes): the newest packet must be held, not a fragment of one", r.waiting, r.waitingSize, r.sentenceBytes)
	}
	if got := strings.Count(r.newest, panelSnapshotGlobal+"("); got != 1 {
		t.Fatalf("the newest sentence carries the named global %d times: the harness is not building whole sentences", got)
	}

	if ran := page.resume(); ran != 1 {
		t.Fatalf("the page ran %d queued delivery/deliveries, want exactly 1", ran)
	}
	got := page.lastEval()
	if got != r.newest {
		t.Fatalf("the page received %q, want the newest whole sentence %q: content must be the packet that was pushed, byte for byte", got, r.newest)
	}
	if n := strings.Count(got, panelSnapshotGlobal+"("); n != 1 {
		t.Fatalf("the delivered sentence carries %d calls, want 1: two packets in one sentence is a merge, which ticket 197 leg B and the 10-08 ruling both forbid", n)
	}
	if !strings.Contains(got, `"seq":`+fmt.Sprint(attempts-1)) {
		t.Fatalf("the delivered sentence is not the newest push (want seq %d): %q", attempts-1, got)
	}

	st := m.PushStats()
	if st.Delivered != 1 || st.Pending != 0 || st.PendingBytes != 0 {
		t.Errorf("after the page drained, PushStats = %+v, want Delivered 1 / Pending 0 / PendingBytes 0", st)
	}
	if st.Coalesced != attempts-1 {
		t.Errorf("coalesced = %d, want %d: an older packet that was superseded must be a counted fact", st.Coalesced, attempts-1)
	}
	if r.heapUnderflow {
		t.Logf("heap reading underflowed and is recorded as such (not as a huge number): delta=%d underflow=%v", r.heapDelta, r.heapUnderflow)
	}
	t.Logf("35-g3 packet reading: %d pushes -> 1 delivery of the newest whole packet, %d bytes, seq=%d, coalesced=%d", attempts, len(got), attempts-1, st.Coalesced)
}

// TestPanelPushQueueSupersessionsAreReadableOnDisk35r9 is this file's replacement
// for 35r8's TestFloodRulerFiresOnAnUnboundedSink35r8 as the "readable exit" ruler,
// with one difference of kind: the reading is taken from the log file ON DISK, not
// from a Go struct. That is what ticket 35 AC#5's census found missing on this side
// (35-recon: failedPost had zero non-test readers), so a push that did not reach the
// page has to be findable by someone holding only <dataDir>\logs\wisp-*.jsonl.
func TestPanelPushQueueSupersessionsAreReadableOnDisk35r9(t *testing.T) {
	const attempts, padBytes = 300, 512
	prev := slog.Default()
	dir := t.TempDir()
	sink, err := installLogSink(dir)
	if err != nil {
		t.Fatalf("installLogSink(%s): %v", dir, err)
	}

	page := newBlockedPage35r9()
	m := NewPanelManager(nil, nil, "")
	attachPanelWindow35r9(t, m, page)
	q := panelPushQueue35r9{m: m, blocked: page}
	r := floodPush35r9(t, q, page, attempts, padBytes)
	if r.coalesced == 0 {
		t.Fatalf("the flood superseded nothing (%d pushes, waiting=%d): this case would assert a record that was never due", attempts, r.waiting)
	}
	// Lift the block: the burst ends and its total is booked.
	if ran := page.resume(); ran != 1 {
		t.Fatalf("the page ran %d queued delivery/deliveries, want exactly 1", ran)
	}

	// Close the pipeline (it flushes and syncs) before reading the file, and restore
	// the process default so later tests are not writing into a closed sink.
	sink.close()
	slog.SetDefault(prev)

	recs := logRecordsUnder35r9(t, dir)
	// (a) A consumer that NEVER drains still leaves a trace: the first-drop record
	// exists, so "the queue filled" is on disk even if the burst never ends.
	dropped := findLastRecord35r9(recs, "panel.push.superseded")
	if dropped == nil {
		t.Fatalf("no panel.push.superseded record under %s: a permanently blocked consumer would leave no trace that the queue filled", dir)
	}
	if got, want := intField35r9(t, dropped, "pending_bytes"), int64(len(pushSentence35r9(1, padBytes))); got != want {
		t.Errorf("first-drop record pending_bytes = %d, want %d (the one whole sentence pending at that instant: the second push superseded the first)", got, want)
	}
	// (b) The burst's own total is on disk, which is the counter a reader needs.
	rec := findLastRecord35r9(recs, "panel.push.burst")
	if rec == nil {
		t.Fatalf("no panel.push.burst record under %s: the burst's total was never booked, so the log can only ever say 'one supersession'", dir)
	}
	if got, want := intField35r9(t, rec, "burst_superseded"), int64(attempts-1); got != want {
		t.Errorf("burst_superseded = %d, want %d - the record must carry the burst's real total", got, want)
	}
	if got, want := intField35r9(t, rec, "coalesced_total"), int64(attempts-1); got != want {
		t.Errorf("record coalesced_total = %d, want %d - the counter on disk must be the counter that ran", got, want)
	}
	if got := intField35r9(t, rec, "queue_cap"); got != panelPushPendingCap {
		t.Errorf("record queue_cap = %d, want the production bound %d", got, panelPushPendingCap)
	}
	if got, want := intField35r9(t, rec, "delivered_sentences"), int64(1); got != want {
		t.Errorf("record delivered_sentences = %d, want %d (the coalesced queue hands over one whole packet)", got, want)
	}
	if got, _ := rec["detail"].(string); !strings.Contains(got, "superseded") {
		t.Errorf("record detail = %q, want a sentence naming what happened", got)
	}
	t.Logf("35-g3 disk reading: %s / burst_superseded=%d coalesced_total=%d queue_cap=%d delivered_sentences=%d; first-drop record pending_bytes=%d",
		rec["msg"], intField35r9(t, rec, "burst_superseded"), intField35r9(t, rec, "coalesced_total"),
		intField35r9(t, rec, "queue_cap"), intField35r9(t, rec, "delivered_sentences"),
		intField35r9(t, dropped, "pending_bytes"))
}

// TestPanelPushRulerFiresOnAnUnboundedQueue35r9 is the positive control the repo
// requires of every bound: the same harness, the same predicate, a queue with the
// bound removed on purpose. The predicate must say OVER cap, the queue must be
// holding exactly the bytes it was given (so the "over" cannot come from a harness
// that stopped measuring), and the same-shaped flood through the production queue
// must stay under cap so the control's fire comes from the missing bound and not
// from the flood's size.
func TestPanelPushRulerFiresOnAnUnboundedQueue35r9(t *testing.T) {
	const padBytes = 4096
	bound := pushCapFor35r9(len(pushSentence35r9(1600-1, padBytes)))

	control := func(attempts int) (pushReading35r9, *blockedPage35r9, *unboundedPushQueue35r9) {
		page := newBlockedPage35r9()
		q := &unboundedPushQueue35r9{page: page}
		return floodPush35r9(t, q, page, attempts, padBytes), page, q
	}

	u1, _, _ := control(400)
	u4, _, _ := control(1600)
	// The harness really measured the planted shape: one held sentence per push.
	if u1.waiting != 400 || u4.waiting != 1600 {
		t.Fatalf("the harness is not measuring the planted queue: 1x held %d, 4x held %d, want 400 / 1600", u1.waiting, u4.waiting)
	}
	if want := pushFloodBytes35r9(400, padBytes); u1.waitingSize != want {
		t.Fatalf("the planted queue pins %d bytes, want exactly %d (400 whole sentences)", u1.waitingSize, want)
	}
	if bad := overPushCap35r9(u1, bound); bad == "" {
		t.Fatalf("the ruler stayed silent on an uncapped 400-push queue: every bounded reading in this file is decoration")
	} else {
		t.Logf("35-g3 positive control fired on 1x as required: %s", bad)
	}
	if bad := overPushCap35r9(u4, bound); bad == "" {
		t.Fatalf("the ruler stayed silent on an uncapped 1600-push queue")
	}
	if u4.waitingSize < 3*u1.waitingSize {
		t.Errorf("the uncapped queue pinned %d bytes at 1x and %d at 4x: it is not reproducing the unbounded shape", u1.waitingSize, u4.waitingSize)
	}
	if u4.heapUnderflow || u1.heapUnderflow {
		t.Logf("positive control heap underflowed (1x=%v 4x=%v): the byte dimension still fired, so the control stands", u1.heapUnderflow, u4.heapUnderflow)
	}

	// Differential: the SAME flood through the production queue stays under cap.
	page := newBlockedPage35r9()
	m := NewPanelManager(nil, nil, "")
	attachPanelWindow35r9(t, m, page)
	bounded := floodPush35r9(t, panelPushQueue35r9{m: m, blocked: page}, page, 400, padBytes)
	if bad := overPushCap35r9(bounded, bound); bad != "" {
		t.Errorf("the same-shaped flood is over cap through the production queue (%s), so the control's fire discriminates nothing", bad)
	}
}

// TestPanelPushHeapDeltaUnderflowIsNotReadAsOverCap35r9 pins the measurement fix
// itself. The 35r8 harness's uint64 subtraction turns "the GC returned memory" into
// a number near 2^64, and its predicate reads that as far over the heap cap - a
// false red. Here: the conversion is signed, an underflow is flagged, and the
// predicate refuses to call it over cap; the same predicate must still fire on a
// genuinely large positive delta, so the fix did not defang it.
func TestPanelPushHeapDeltaUnderflowIsNotReadAsOverCap35r9(t *testing.T) {
	if got, under := heapDelta35r9(1000, 400); under != true || got != 0 {
		t.Fatalf("heapDelta35r9(1000, 400) = (%d, %v), want (0, true): a shrink is an underflow, not a delta", got, under)
	}
	if got, under := heapDelta35r9(400, 1000); under != false || got != 600 {
		t.Fatalf("heapDelta35r9(400, 1000) = (%d, %v), want (600, false)", got, under)
	}
	if got, under := heapDelta35r9(7, 7); under != false || got != 0 {
		t.Fatalf("heapDelta35r9(7, 7) = (%d, %v), want (0, false)", got, under)
	}

	bound := pushCapFor35r9(4096)
	// A reading whose only sin is an underflowed heap delta: the predicate must not
	// turn that into "over cap" (that is exactly the 35r8 false red).
	underflowed := pushReading35r9{
		attempts: 100, sentenceBytes: 4096, waiting: 1, wakeups: 1,
		waitingSize: 4096, coalesced: 99,
		heapDelta: 0, heapUnderflow: true,
	}
	if bad := overPushCap35r9(underflowed, bound); bad != "" {
		t.Fatalf("an underflowed heap reading was judged OVER cap (%s): a measurement error must not read as a subject failure", bad)
	}
	// And the heap dimension still has teeth: a real, large growth fires.
	grown := pushReading35r9{attempts: 100, sentenceBytes: 4096, waiting: 1, wakeups: 1, waitingSize: 4096, coalesced: 99, heapDelta: bound.heap + 1}
	if bad := overPushCap35r9(grown, bound); bad == "" {
		t.Fatalf("the heap dimension stopped judging a real growth over the cap: the underflow fix defanged the ruler")
	}
	// The reading the OLD unsigned shape produced from the same two HeapAlloc values
	// is 2^64-600. In a signed field that value is −600 (below), which is why the
	// old shape is banned here rather than merely discouraged: it cannot make the
	// predicate fire for the right reason, only for an arithmetic accident.
	var wrappedSource uint64 // 0 - 600 wraps at run time, exactly as the old subtraction did
	wrappedSource -= 600
	if wrapped := int64(wrappedSource); wrapped >= 0 {
		t.Fatalf("int64(0-600 under uint64) = %d, expected a negative wrap: this case no longer documents the old shape's collapse", wrapped)
	}
	t.Logf("35-g3 heap underflow control: underflowed reading not judged over cap; grown reading judged OVER; the old unsigned value collapses to %d in a signed field", int64(wrappedSource))
}

// TestPanelPostRefusalHasAReadPointOnDisk35r9 closes ticket 35 AC#5's second half
// (the census's "failedPost, non-test readers = 0"): a refused post must be counted
// AND findable on disk, and the counter's second exit - the panel status line - must
// quote the same number. The refusal is driven through the REAL post() on a REAL
// residentPanel whose task channel has no room and no reader, so this is the
// production path, not a counter bump.
func TestPanelPostRefusalHasAReadPointOnDisk35r9(t *testing.T) {
	prev := slog.Default()
	dir := t.TempDir()
	sink, err := installLogSink(dir)
	if err != nil {
		t.Fatalf("installLogSink(%s): %v", dir, err)
	}

	mgr := NewPanelManager(nil, nil, "")
	rp := &residentPanel{
		mgr:   mgr,
		tasks: make(chan func()), // no capacity, no reader: every post() is refused
	}
	if rp.post(func() {}) {
		t.Fatalf("post() accepted a request although the task channel has no room and no reader")
	}
	if rp.post(func() {}) {
		t.Fatalf("post() accepted a second request on the same full channel")
	}
	if got := rp.failedPost.Load(); got != 2 {
		t.Fatalf("failedPost = %d after two refused posts, want 2", got)
	}
	if line := rp.statusLine(); !strings.Contains(line, "2 refused post(s)") {
		t.Errorf("statusLine = %q, want it to carry '2 refused post(s)': the boot report is the counter's readable second exit", line)
	}

	sink.close()
	slog.SetDefault(prev)

	recs := logRecordsUnder35r9(t, dir)
	// One record per refusal, so a reader can see the event as well as the total.
	if got := countRecords35r9(recs, "panel.post.refused"); got != 2 {
		t.Fatalf("%d panel.post.refused record(s) under %s, want 2 (one per refusal): a refused post is still a silence outside the test binary", got, dir)
	}
	rec := findLastRecord35r9(recs, "panel.post.refused")
	if got := intField35r9(t, rec, "failed_posts"); got != 2 {
		t.Errorf("record failed_posts = %d, want the counter's 2", got)
	}
	if got := intField35r9(t, rec, "queue_cap"); got != 0 {
		t.Errorf("record queue_cap = %d, want this channel's 0", got)
	}
	t.Logf("35-g3 post refusal reading: 2 records on disk, last failed_posts=%d, statusLine carries the same count", intField35r9(t, rec, "failed_posts"))
}
