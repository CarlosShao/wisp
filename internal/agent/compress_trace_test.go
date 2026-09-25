package agent

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/llm"
)

// Ticket 139 AC#3: a compression pass that folded history MUST leave a trace,
// and the trace must be the one thing that says "how much did we drop".
//
// Two shapes of false green are in scope here and the file is shaped around
// killing both:
//
//   - A tautological judgement ("there is a log line after compression") that
//     also holds when nothing was compressed. TestCompressionTraceSilentWhenNothingFolded
//     is the fang: the same capture must stay EMPTY for a pass that folded
//     nothing, so a trace printed unconditionally fails here.
//   - A trace that exists only in a unit call while the production wiring
//     dropped the logger. TestCompressionTraceSurvivesTheLoopWiring drives a
//     real Loop.Run, so the record must come out of the compressor the loop
//     built itself (loop.go New()), not one the test assembled.
//
// Every threshold below is read back out of Budgets; none is a written number
// (AC#1(2) forbids citing the 12000 reference value as an actual trigger).

const traceMsg = "agent: history compressed"

// rulerSentinel is what the capture-instrument positive control emits. It has
// to be seen by the very same handler a test then reads compression records
// from, or a "no trace here" reading is worth nothing.
const rulerSentinel = "trace-capture-ruler-control"

// ---------------------------------------------------------------------------
// in-memory slog handler

type capturedRecord struct {
	msg  string
	attr map[string]any
}

// traceCapture keeps every record it is handed. It renders nothing: assertions
// read the message and the attribute set, which is what the AC asks for.
type traceCapture struct {
	mu   sync.Mutex
	recs []capturedRecord
}

func newTraceCapture() *traceCapture { return &traceCapture{} }

func (tc *traceCapture) Enabled(context.Context, slog.Level) bool { return true }

func (tc *traceCapture) Handle(_ context.Context, r slog.Record) error {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	attrs := map[string]any{}
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})
	tc.recs = append(tc.recs, capturedRecord{msg: r.Message, attr: attrs})
	return nil
}

func (tc *traceCapture) WithAttrs([]slog.Attr) slog.Handler { return tc }
func (tc *traceCapture) WithGroup(string) slog.Handler      { return tc }

func (tc *traceCapture) logger() *slog.Logger { return slog.New(tc) }

func (tc *traceCapture) all() []capturedRecord {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	out := make([]capturedRecord, len(tc.recs))
	copy(out, tc.recs)
	return out
}

// with returns the records carrying `msg`.
func (tc *traceCapture) with(msg string) []capturedRecord {
	var out []capturedRecord
	for _, r := range tc.all() {
		if r.msg == msg {
			out = append(out, r)
		}
	}
	return out
}

// assertRulerLive drives the capture with a record this package's production
// code did not produce, and fails if the capture did not see it. Called at the
// top of every test below: a 0-record reading is only meaningful behind it.
func (tc *traceCapture) assertRulerLive(t *testing.T) {
	t.Helper()
	before := len(tc.all())
	tc.logger().Info(rulerSentinel, "control", 1)
	if len(tc.all()) != before+1 {
		t.Fatalf("capture instrument is dead: %d records before, %d after emitting one",
			before, len(tc.all()))
	}
	if got := tc.with(rulerSentinel); len(got) != 1 {
		t.Fatalf("capture instrument lost the control record: %d hits for %q", len(got), rulerSentinel)
	}
}

// flat renders a record (message plus every attribute) for the privacy scan.
func (r capturedRecord) flat() string {
	parts := make([]string, 0, len(r.attr)+1)
	parts = append(parts, r.msg)
	for k, v := range r.attr {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(parts, " ")
}

func attrInt(t *testing.T, r capturedRecord, key string) int64 {
	t.Helper()
	v, ok := r.attr[key]
	if !ok {
		t.Fatalf("trace record has no %q attribute; got %v", key, r.attr)
	}
	switch n := v.(type) {
	case int:
		return int64(n)
	case int64:
		return n
	default:
		t.Fatalf("attribute %q is %T (%v), want an integer", key, v, v)
		return 0
	}
}

func attrBool(t *testing.T, r capturedRecord, key string) bool {
	t.Helper()
	v, ok := r.attr[key]
	if !ok {
		t.Fatalf("trace record has no %q attribute; got %v", key, r.attr)
	}
	b, ok := v.(bool)
	if !ok {
		t.Fatalf("attribute %q is %T (%v), want a bool", key, v, v)
	}
	return b
}

// requireTraceAttrs pins the attribute names AC#2 asks for, by name, on one
// record. It returns nothing: the caller reads them back with attrInt.
func requireTraceAttrs(t *testing.T, r capturedRecord) {
	t.Helper()
	for _, key := range []string{
		"tokens_before", "tokens_after", "threshold",
		"compressed_msgs", "kept_raw_rounds", "history_changed",
	} {
		if _, ok := r.attr[key]; !ok {
			t.Fatalf("trace record lacks attribute %q; present: %v", key, r.attr)
		}
	}
}

// ---------------------------------------------------------------------------
// the trace itself

// A fold that happens is booked: one record, and its counts say the same thing
// the report does. The privacy half of AC#2 is asserted in the same test -
// every round of the fixture history carries a sentinel string, and the record
// must not carry it.
func TestCompressionTraceBooksCountsOnSuccess(t *testing.T) {
	recs := newTraceCapture()
	recs.assertRulerLive(t)

	b := BudgetsFor(4096) // HistoryCompressTokens is derived: 384
	th := b.HistoryCompressTokens
	c := NewCompressor(b, nil, WithLogger(recs.logger()))

	const rounds = 6
	const sentinel = "PRIVACY-SENTINEL-do-not-log-this"
	hist := buildRoundHistory(rounds, 400)
	// Put the sentinel into the prose the compressor folds, so any leak of
	// original content into the trace is caught by name.
	for i := range hist {
		for j, p := range hist[i].Content {
			if tp, ok := p.(llm.TextPart); ok {
				hist[i].Content[j] = llm.TextPart{Text: tp.Text + " " + sentinel}
			}
		}
	}
	before := c.TotalTokens(hist)
	if !c.Need(hist) {
		t.Fatalf("setup: %d tokens does not clear the derived threshold %d", before, th)
	}

	out, rep, err := c.Compress(context.Background(), hist)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if !rep.Ran {
		t.Fatal("setup: pass did not fold, so the trace has nothing to report")
	}

	hits := recs.with(traceMsg)
	if len(hits) != 1 {
		t.Fatalf("%q records = %d, want exactly 1 (all records: %v)",
			traceMsg, len(hits), flattenRecordMsgs(recs.all()))
	}
	rec := hits[0]
	requireTraceAttrs(t, rec)

	if got := attrInt(t, rec, "tokens_before"); got != int64(before) || got != int64(rep.TokensBefore) {
		t.Errorf("tokens_before = %d, want the pre-pass total (%d / report %d)", got, before, rep.TokensBefore)
	}
	if got := attrInt(t, rec, "tokens_after"); got != int64(rep.TokensAfter) {
		t.Errorf("tokens_after = %d, want the report's %d", got, rep.TokensAfter)
	}
	// The booked number must describe the history this call returned.
	if got, want := int64(c.TotalTokens(out)), attrInt(t, rec, "tokens_after"); got != want {
		t.Errorf("re-measured output = %d tokens, trace booked %d", got, want)
	}
	// D15(4) folds until the budget holds OR only KeepRawRounds raw rounds are
	// left, and these fixtures hit the second exit: tokens_after legitimately
	// stays above threshold. That is why the record carries "threshold" beside
	// the two totals - a reader must be able to tell "it fits now" from "it
	// stopped at the floor" without opening the source.
	if c.Need(out) && rep.KeptRawRounds > b.KeepRawRounds {
		t.Errorf("output still needs compressing (%d tokens over threshold %d) with %d raw rounds left, want the floor at %d",
			c.TotalTokens(out), th, rep.KeptRawRounds, b.KeepRawRounds)
	}
	if got := attrInt(t, rec, "threshold"); got != int64(th) {
		t.Errorf("threshold = %d, want the budget-derived %d", got, th)
	}
	if got := attrInt(t, rec, "compressed_msgs"); got <= 0 {
		t.Errorf("compressed_msgs = %d, want > 0 for a pass that folded", got)
	}
	if got := attrInt(t, rec, "kept_raw_rounds"); got != int64(b.KeepRawRounds) {
		t.Errorf("kept_raw_rounds = %d, want %d", got, b.KeepRawRounds)
	}
	if !attrBool(t, rec, "history_changed") {
		t.Error("history_changed = false, want true for a pass that replaced the history")
	}
	// The counts must be able to answer "how much was dropped" on their own.
	if attrInt(t, rec, "tokens_after") >= attrInt(t, rec, "tokens_before") {
		t.Errorf("trace reports %d -> %d tokens, which does not read as a reduction",
			attrInt(t, rec, "tokens_before"), attrInt(t, rec, "tokens_after"))
	}

	// AC#2 hard rule: counts only, never the folded text.
	for _, r := range recs.with(traceMsg) {
		if strings.Contains(r.flat(), sentinel) {
			t.Errorf("trace leaked history content: %q appears in %q", sentinel, r.flat())
		}
	}
	if strings.Contains(rec.flat(), "问题") {
		t.Errorf("trace leaked a user turn body: %q", rec.flat())
	}
}

// The fang against a tautological judgement: a pass that folded nothing leaves
// nothing. If the trace were emitted on every Compress call this goes red even
// though the "compression happened" test is still green.
func TestCompressionTraceSilentWhenNothingFolded(t *testing.T) {
	recs := newTraceCapture()
	recs.assertRulerLive(t)

	b := BudgetsFor(4096)
	c := NewCompressor(b, nil, WithLogger(recs.logger()))
	hist := buildRoundHistory(2, 20) // far under the derived threshold
	if c.Need(hist) {
		t.Fatalf("setup: %d tokens clears threshold %d, want a no-op pass",
			c.TotalTokens(hist), b.HistoryCompressTokens)
	}

	out, rep, err := c.Compress(context.Background(), hist)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if rep.Ran {
		t.Fatal("setup: a no-op pass reported Ran=true")
	}
	if c.TotalTokens(out) != c.TotalTokens(hist) {
		t.Fatal("setup: no-op pass changed the history")
	}
	if got := recs.with(traceMsg); len(got) != 0 {
		t.Errorf("no-op pass left %d trace record(s), want none: %v",
			len(got), flattenRecordMsgs(got))
	}
}

// The trace must survive the wiring, not just the constructor: New(Options) owns
// the compressor the loop calls, so a real Run over a seeded history has to put
// the record on the logger the host installed. This is also the only test here
// that reads Result.Compression - as a cross-check that the line and the
// outwire say the same thing.
func TestCompressionTraceSurvivesTheLoopWiring(t *testing.T) {
	recs := newTraceCapture()
	recs.assertRulerLive(t)

	b := BudgetsFor(4096)
	h := newHarness(t, "text-reply",
		withConfig(func(c *Config) { c.ContextWindow = 4096 }),
		withLogger(recs.logger()))
	if got := h.loop.Budgets().HistoryCompressTokens; got != b.HistoryCompressTokens {
		t.Fatalf("loop threshold = %d, want the window-derived %d", got, b.HistoryCompressTokens)
	}

	// Seed enough raw rounds that the fold has oldest rounds to take: four
	// rounds over the derived threshold, appended through the loop's own
	// history so the pass runs on the real conversation.
	for _, m := range buildRoundHistory(4, 400) {
		h.loop.append(m)
	}
	seeded := h.loop.History()
	startTokens := NewCompressor(b, nil).TotalTokens(seeded)
	if startTokens <= b.HistoryCompressTokens {
		t.Fatalf("setup: seeded %d tokens, want over %d", startTokens, b.HistoryCompressTokens)
	}

	res := h.run("收个尾")
	if res.Status != StatusCompleted {
		t.Fatalf("run status = %s (%s), want completed: the trace test needs a clean pass",
			res.Status, res.Message)
	}

	hits := recs.with(traceMsg)
	if len(hits) == 0 {
		t.Fatalf("a real Run that folded history left no %q record; captured: %v",
			traceMsg, flattenRecordMsgs(recs.all()))
	}
	rec := hits[0]
	requireTraceAttrs(t, rec)
	if !attrBool(t, rec, "history_changed") {
		t.Error("history_changed = false on a run that replaced the history")
	}
	// The record and the Result outwire must agree, or there are two answers.
	if got, want := attrInt(t, rec, "tokens_before"), int64(res.Compression.TokensBefore); got != want {
		t.Errorf("trace tokens_before = %d but Result.Compression says %d", got, want)
	}
	if got, want := attrInt(t, rec, "tokens_after"), int64(res.Compression.TokensAfter); got != want {
		t.Errorf("trace tokens_after = %d but Result.Compression says %d", got, want)
	}
	if !res.Compression.Ran {
		t.Error("Result.Compression.Ran = false while the trace fired")
	}
	// And the history the next round sees is the shrunken one.
	if end := NewCompressor(b, nil).TotalTokens(h.loop.History()); end >= startTokens {
		t.Errorf("history after the run = %d tokens, want less than the seeded %d", end, startTokens)
	}
	if len(hits) != 1 {
		t.Errorf("%q records = %d, want 1 for a single-round run: %v",
			traceMsg, len(hits), flattenRecordMsgs(hits))
	}
}

// The trace must not be paid for with the D15(4) invariant it rides on: folding
// is what the compressor does, and adding a record cannot change what it hands
// back. compress_test.go owns the id-ledger assertions; this one only pins that
// a logger-injected compressor returns the same shape as the plain one.
func TestCompressionTraceDoesNotAlterTheFold(t *testing.T) {
	b := BudgetsFor(4096)
	hist := buildRoundHistory(6, 400)

	plain, pref, err := NewCompressor(b, nil).Compress(context.Background(), hist)
	if err != nil {
		t.Fatalf("baseline Compress: %v", err)
	}
	traced, trep, err := NewCompressor(b, nil, WithLogger(newTraceCapture().logger())).
		Compress(context.Background(), hist)
	if err != nil {
		t.Fatalf("traced Compress: %v", err)
	}
	if len(plain) != len(traced) {
		t.Fatalf("message count differs: plain %d, traced %d", len(plain), len(traced))
	}
	for i := range plain {
		if len(plain[i].Content) != len(traced[i].Content) {
			t.Fatalf("message %d part count differs", i)
		}
	}
	if pref.TokensBefore != trep.TokensBefore || pref.TokensAfter != trep.TokensAfter ||
		pref.CompressedMsgs != trep.CompressedMsgs {
		t.Errorf("report differs once a logger is attached: %+v vs %+v", pref, trep)
	}
}

func flattenRecordMsgs(rs []capturedRecord) []string {
	out := make([]string, 0, len(rs))
	for _, r := range rs {
		out = append(out, r.msg)
	}
	return out
}
