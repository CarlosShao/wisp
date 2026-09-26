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
// Ticket 153 AC#1 adds the third shape those two bullets do not cover: a pass
// that clears the threshold and still folds nothing (one raw round against the
// KeepRawRounds floor - the shape the real `wisp run` leg has today). Without
// it the "print whenever Need() says yes" guard - mutation M5 of the ticket 139
// acceptance - keeps every nail here green.
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

// ---------------------------------------------------------------------------
// ticket 153 AC#1: the shape the four nails above do not cover - the history IS
// over the threshold and the pass still cannot fold a single round.
//
// TestCompressionTraceSilentWhenNothingFolded treats "did not fold" as if it
// only ever meant "did not clear the threshold": it hard-fails its own setup
// the moment Need() is true (see the t.Fatalf above), so the over-threshold
// no-fold shape exists in none of the four nails. That hole is why the guard
// shape the 139 acceptance mutation M5 replaced `if rep.Ran` with
// (`if c.Need(hist)`, "print whenever the threshold was crossed") kept all four
// green and the whole package green.
//
// The shape is also what the real `wisp run` leg has today: one task is one
// user message, groupRounds opens a round only at a user role, and D15(4)'s
// floor keeps the last KeepRawRounds rounds verbatim, so a single-round history
// has nothing to fold however large it grows. The 139 acceptance measured
// exactly that pair on one disk at one threshold: the ordinary tool-loop shape
// (~41 706 tokens over 3 072) left 0 records, the degenerate multi-ladder shape
// left 1 - the only variable being whether a 4th user round existed
// (139 acceptance r1 s7.2 / s8, which is also the section that retracted its
// own s1).
//
// This case is the fang, and it cannot collapse into the nail above: the setup
// asserts Need() is TRUE.

func TestCompressionTraceSilentWhenNothingFoldableOverThreshold(t *testing.T) {
	recs := newTraceCapture()
	recs.assertRulerLive(t)

	b := BudgetsFor(4096) // threshold read back out of the budget, never written
	th := b.HistoryCompressTokens
	c := NewCompressor(b, nil, WithLogger(recs.logger()))

	// One raw round, deliberately fat enough to clear the derived threshold.
	hist := buildRoundHistory(1, 2000)
	if got := len(rawRoundIndexes(groupRounds(hist))); got != 1 {
		t.Fatalf("setup: %d raw rounds, want exactly 1 - with more than KeepRawRounds the fold would be possible and this would no longer be the shape under test", got)
	}
	if !c.Need(hist) {
		t.Fatalf("setup: %d tokens does NOT clear the derived threshold %d - this case would silently collapse into TestCompressionTraceSilentWhenNothingFolded",
			c.TotalTokens(hist), th)
	}

	out, rep, err := c.Compress(context.Background(), hist)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	if rep.Ran {
		t.Fatal("setup: the pass folded, so this is not the over-threshold/nothing-foldable shape")
	}
	if len(out) != len(hist) || c.TotalTokens(out) != c.TotalTokens(hist) {
		t.Fatalf("a pass that folded nothing changed the history: %d msgs / %d tokens -> %d msgs / %d tokens",
			len(hist), c.TotalTokens(hist), len(out), c.TotalTokens(out))
	}

	hits := recs.with(traceMsg)
	for _, h := range hits {
		t.Errorf("over-threshold pass that folded nothing left a trace: %s", h.flat())
	}
	if len(hits) != 0 {
		t.Fatalf("%q records = %d, want none: a record reading \"compressed\" while nothing folded is a false positive, and %d raw round(s) against the KeepRawRounds floor of %d is exactly the case where no fold can happen (all records: %v)",
			traceMsg, len(hits), len(rawRoundIndexes(groupRounds(hist))), b.KeepRawRounds, flattenRecordMsgs(recs.all()))
	}
}

// ---------------------------------------------------------------------------
// ticket 153 AC#2: the trace has to say which task it belongs to
//
// 139's own delivery named this gap (impl r1 s2.2: "痕里没有 taskID") and the
// acceptance confirmed it was still open. The compressor cannot see a task id:
// agent.New builds one Compressor per Loop (loop.go) before any task exists,
// and the id is minted inside Loop.run. So the two readings below pin the two
// halves of the fix: a real Run attributes its fold to the id that Run mints
// (the same string Result.TaskID and every Event of that task carry - not a
// lookalike), and a pass nobody attributed stays unattributed instead of being
// handed a placeholder that would read as a real owner.

func TestCompressionTraceCarriesTheOwningTaskID(t *testing.T) {
	recs := newTraceCapture()
	recs.assertRulerLive(t)

	b := BudgetsFor(4096)
	h := newHarness(t, "text-reply",
		withConfig(func(c *Config) { c.ContextWindow = 4096 }),
		withLogger(recs.logger()))
	// Four raw rounds over the derived threshold: the fold has oldest rounds to
	// take, so this really is a pass that ran (see AC#1's case for the shape
	// where it cannot).
	for _, m := range buildRoundHistory(4, 400) {
		h.loop.append(m)
	}

	res := h.run("收个尾")
	if res.Status != StatusCompleted {
		t.Fatalf("run status = %s (%s), want completed: attribution needs a pass that ran",
			res.Status, res.Message)
	}
	hits := recs.with(traceMsg)
	if len(hits) != 1 {
		t.Fatalf("%q records = %d, want exactly 1 from one real Run (all: %v)",
			traceMsg, len(hits), flattenRecordMsgs(recs.all()))
	}
	rec := hits[0]

	v, ok := rec.attr["task"]
	if !ok {
		t.Fatalf("trace carries no \"task\" attribute; present: %v - an unattributable fold cannot be joined to the task that paid for it (D37 attribution)", rec.attr)
	}
	task, isString := v.(string)
	if !isString {
		t.Fatalf("\"task\" attribute is %T (%v), want the loop's task id as a string", v, v)
	}
	if task != res.TaskID {
		t.Errorf("trace task = %q, want the Result's own %q", task, res.TaskID)
	}
	// The record still answers to the fold it came from: same threshold the
	// loop froze for this window, same report the Result outwire carries.
	if got := attrInt(t, rec, "threshold"); got != int64(b.HistoryCompressTokens) {
		t.Errorf("trace threshold = %d, want the window-derived %d", got, b.HistoryCompressTokens)
	}
	if got, want := attrInt(t, rec, "tokens_before"), int64(res.Compression.TokensBefore); got != want {
		t.Errorf("trace tokens_before = %d, want the Result's %d", got, want)
	}
	// The shape check is not decoration: newTaskID mints a 36-char uuid-shaped
	// id, and this is what tells the reading apart from a hand-typed constant.
	if len(task) != 36 || strings.Count(task, "-") != 4 {
		t.Errorf("trace task = %q, want the uuid-shaped id newTaskID mints (36 chars, 4 dashes)", task)
	}
}

// The other half: attribution is per call, and nothing is ever invented to fill
// the gap. A shared Compressor that stored the id as a field would pass the test
// above and fail the ordering asserted here.
func TestCompressionTraceNeverInventsATaskID(t *testing.T) {
	recs := newTraceCapture()
	recs.assertRulerLive(t)

	b := BudgetsFor(4096)
	c := NewCompressor(b, nil, WithLogger(recs.logger()))
	hist := buildRoundHistory(6, 400)
	if _, rep, err := c.Compress(context.Background(), hist); err != nil || !rep.Ran {
		t.Fatalf("setup: the fixture does not fold (err=%v ran=%v)", err, rep.Ran)
	}

	// An untagged call still leaves the trace - just with nothing to say about
	// who owned it.
	untagged := recs.with(traceMsg)
	if len(untagged) != 1 {
		t.Fatalf("%q records = %d, want 1 for a pass that folded", traceMsg, len(untagged))
	}
	if v, present := untagged[0].attr["task"]; present {
		t.Fatalf("untagged call produced a \"task\" attribute %q (%T): absence is honest, an invented owner is not", v, v)
	}

	// Two ids from the real mint, three calls, one shared compressor: each
	// record must answer to the call that produced it.
	idA, idB := newTaskID(), newTaskID()
	if idA == idB {
		t.Fatalf("setup: newTaskID returned the same id twice (%q)", idA)
	}
	ctxA := withTraceTask(context.Background(), idA)
	ctxB := withTraceTask(context.Background(), idB)
	for _, ctx := range []context.Context{ctxA, ctxB, ctxA} {
		if _, _, err := c.Compress(ctx, hist); err != nil {
			t.Fatalf("Compress: %v", err)
		}
	}
	tagged := recs.with(traceMsg)[1:]
	want := []string{idA, idB, idA}
	if len(tagged) != len(want) {
		t.Fatalf("tagged calls left %d records, want %d", len(tagged), len(want))
	}
	for i, r := range tagged {
		got, present := r.attr["task"]
		if !present {
			t.Fatalf("tagged call %d left no \"task\" attribute: %v", i, r.attr)
		}
		if got != want[i] {
			t.Errorf("record %d carries task %v, want %q - the id rides on the call, not on the shared compressor", i, got, want[i])
		}
	}

	// And no placeholder value ever reaches a record, on any attribute: that is
	// how "cannot attribute" gets laundered into "attributed to nobody".
	for _, r := range recs.with(traceMsg) {
		for k, v := range r.attr {
			s, isString := v.(string)
			if !isString {
				continue
			}
			switch strings.ToLower(strings.TrimSpace(s)) {
			case "", "unknown", "none", "nil", "n/a", "-":
				t.Errorf("attribute %q carries the placeholder value %q on record %q", k, s, r.msg)
			}
		}
	}
}
