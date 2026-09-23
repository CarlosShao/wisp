package observe

// Ticket 130 AC#3, the unit half. The acceptance half is a subprocess case in
// cmd/wisp/early_log_nail_130_windows_test.go, because the record that started
// this ticket is made by ANOTHER package's init() and no test inside this package
// can re-run that. What belongs here is the machinery's own contract: the order
// the buffer replays in, the two caps, the fact that the console copy runs no
// matter what the buffer decided, and the one property that cost a 600s hang to
// learn - the console handler must be concrete.
//
// Every case below builds its own earlyLogBuffer, so none of them touches the
// package-level one the init() installed. That global is shared by every other
// test in this package, and a case that drained it would reorder the console
// records the rest of the suite reads.

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"
)

// capture130Handler is the drain target: it keeps what it was handed so a case
// can read the replay back out in order.
type capture130Handler struct {
	mu    sync.Mutex
	level slog.Level
	recs  []slog.Record
}

func newCapture130(level slog.Level) *capture130Handler {
	return &capture130Handler{level: level}
}

func (c *capture130Handler) Enabled(_ context.Context, l slog.Level) bool { return l >= c.level }

func (c *capture130Handler) Handle(_ context.Context, r slog.Record) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.recs = append(c.recs, r.Clone())
	return nil
}

func (c *capture130Handler) WithAttrs([]slog.Attr) slog.Handler { return c }
func (c *capture130Handler) WithGroup(string) slog.Handler      { return c }

func (c *capture130Handler) snapshot() []slog.Record {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]slog.Record(nil), c.recs...)
}

func (c *capture130Handler) msgs() []string {
	var out []string
	for _, r := range c.snapshot() {
		out = append(out, r.Message)
	}
	return out
}

// held reports how many records this buffer still keeps. Direct field access on
// purpose: the buffer's own size is the thing under test, and the exported
// counter on the package-level buffer would read a different object.
func (b *earlyLogBuffer) held130() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.recs)
}

// attrsOf130 flattens a record's attributes, because slog.Record exposes them
// only through its iterator.
func attrsOf130(r slog.Record) []slog.Attr {
	var out []slog.Attr
	r.Attrs(func(a slog.Attr) bool {
		out = append(out, a)
		return true
	})
	return out
}

// clock130 hands out strictly increasing stamps: the replay-order assertion needs
// distinguishable times, and the wall clock on this host can repeat inside a
// tight loop.
var (
	clock130Mu sync.Mutex
	clock130N  int64
)

func clock130() time.Time {
	clock130Mu.Lock()
	defer clock130Mu.Unlock()
	clock130N++
	return time.Unix(1700000000, clock130N)
}

// emit130 puts one INFO record through a fresh buffering handler over b, the way
// the installed tee would.
func emit130(b *earlyLogBuffer, msg string, attrs ...slog.Attr) {
	h := &earlyBufferHandler{b: b}
	rec := slog.NewRecord(clock130(), slog.LevelInfo, msg, 0)
	rec.AddAttrs(attrs...)
	if err := h.Handle(context.Background(), rec); err != nil {
		panic(err)
	}
}

// TestEarlyLogBufferReplaysInCaptureOrderThenStops pins the two halves of the
// replay contract: capture order is the on-disk order (which is what puts a
// boot record ahead of the install booking in cmd/wisp), and the buffer is closed
// by the first drain, so a caller cannot replay the same record twice.
func TestEarlyLogBufferReplaysInCaptureOrderThenStops(t *testing.T) {
	b := &earlyLogBuffer{}
	emit130(b, "first")
	emit130(b, "second", slog.String("k", "v"))
	emit130(b, "third")

	sink := newCapture130(slog.LevelInfo)
	flushed, dropped := b.drain(sink)
	if flushed != 3 || dropped != 0 {
		t.Fatalf("drain returned flushed=%d dropped=%d, want 3 and 0 (captured: %v)", flushed, dropped, sink.msgs())
	}
	got := sink.msgs()
	if len(got) != 3 || got[0] != "first" || got[1] != "second" || got[2] != "third" {
		t.Fatalf("replay order = %v, want [first second third]: the order records were spoken in is the claim", got)
	}
	recs := sink.snapshot()
	if k := attrsOf130(recs[1]); len(k) != 1 || k[0].Key != "k" {
		t.Errorf("replayed record 1 lost its attributes: %v", k)
	}
	if !recs[0].Time.Before(recs[2].Time) {
		t.Errorf("replay did not preserve the stamps it captured: %s then %s", recs[0].Time, recs[2].Time)
	}

	// Closed after the first drain: the pipeline carries records from here on,
	// and a second drain must not invent a replay out of nothing.
	emit130(b, "after")
	if n := b.held130(); n != 0 {
		t.Errorf("buffer holds %d record(s) after the drain, want 0", n)
	}
	flushed2, dropped2 := b.drain(newCapture130(slog.LevelInfo))
	if flushed2 != 0 || dropped2 != 0 {
		t.Errorf("second drain returned %d/%d, want 0/0 (drain must be idempotent)", flushed2, dropped2)
	}
}

// TestEarlyLogBufferOverflowIsABookedFact pins the overflow policy: over the cap
// the buffer stops appending and counts, and the count reaches the destination as
// its own record. "Dropped N early records" must be a fact on disk, not a silent
// subtraction, and the record named has to be the OLDEST kept one, because in a
// boot order the earliest record is usually the cause.
func TestEarlyLogBufferOverflowIsABookedFact(t *testing.T) {
	b := &earlyLogBuffer{}
	for i := 0; i < earlyLogMaxRecords+5; i++ {
		emit130(b, fmt.Sprintf("boot-%02d", i))
	}
	sink := newCapture130(slog.LevelInfo)
	flushed, dropped := b.drain(sink)
	if flushed != earlyLogMaxRecords {
		t.Errorf("flushed = %d, want the capacity %d", flushed, earlyLogMaxRecords)
	}
	if dropped != 5 {
		t.Errorf("dropped = %d, want 5", dropped)
	}
	recs := sink.snapshot()
	if len(recs) != earlyLogMaxRecords+1 {
		t.Fatalf("captured %d records, want %d replayed plus the overflow record", len(recs), earlyLogMaxRecords)
	}
	if recs[0].Message != "boot-00" {
		t.Errorf("first replayed record = %q, want the oldest kept one", recs[0].Message)
	}
	over := recs[len(recs)-1]
	if over.Message != earlyLogOverflowMsg {
		t.Fatalf("last record = %q, want %q", over.Message, earlyLogOverflowMsg)
	}
	if over.Level != slog.LevelWarn {
		t.Errorf("overflow level = %q, want WARN", over.Level)
	}
	want := map[string]string{
		"dropped":    "5",
		"first_kept": "boot-00",
		"capacity":   fmt.Sprint(earlyLogMaxRecords),
	}
	for k, v := range want {
		if got := attrString130(over, k); got != v {
			t.Errorf("overflow attribute %s = %q, want %q", k, got, v)
		}
	}
}

// TestEarlyLogBufferByteCapAlsoStopsAppending pins the second cap, the one that
// bounds memory when a boot logs a few very large records. Three records half
// the byte budget each: the first two fit, the third is refused.
func TestEarlyLogBufferByteCapAlsoStopsAppending(t *testing.T) {
	b := &earlyLogBuffer{}
	big := strings.Repeat("x", earlyLogMaxBytes/3)
	for _, msg := range []string{"half-1", "half-2", "half-3"} {
		emit130(b, msg, slog.String("payload", big))
	}
	flushed, dropped := b.drain(newCapture130(slog.LevelInfo))
	if flushed != 2 || dropped != 1 {
		t.Fatalf("byte cap returned flushed=%d dropped=%d, want 2 and 1", flushed, dropped)
	}
}

// TestEarlyLogBufferCarriesTheWithChain covers the shape that makes this buffer
// safe for a captured logger: internal/memory/open.go and internal/agent/loop.go
// both take slog.Default() at construction time and keep it, so a logger built
// off the early tee must contribute the SAME record to the buffer that it
// contributes to the console - attributes and groups included, and in the order a
// handler prints them (chain attributes first, the call's own second).
func TestEarlyLogBufferCarriesTheWithChain(t *testing.T) {
	b := &earlyLogBuffer{}
	h := (&earlyBufferHandler{b: b}).
		withAttrs([]slog.Attr{slog.String("component", "memory")}).
		withGroup("ctx").
		withAttrs([]slog.Attr{slog.Int("attempt", 2)}).
		withGroup("") // an empty group name is a no-op, per slog's own rule

	rec := slog.NewRecord(clock130(), slog.LevelWarn, "retention job scheduled", 0)
	rec.AddAttrs(slog.String("job", "prune"))
	if err := h.Handle(context.Background(), rec); err != nil {
		t.Fatal(err)
	}
	sink := newCapture130(slog.LevelInfo)
	if flushed, _ := b.drain(sink); flushed != 1 {
		t.Fatalf("flushed = %d, want 1 (captured: %v)", flushed, sink.msgs())
	}
	got := sink.snapshot()[0]
	if got.Message != "retention job scheduled" || got.Level != slog.LevelWarn {
		t.Fatalf("replayed the wrong record: %q / %s", got.Message, got.Level)
	}
	outer := attrsOf130(got)
	if len(outer) != 2 {
		t.Fatalf("attribute count = %d, want component outside the group and one ctx group holding the rest: %v",
			len(outer), outer)
	}
	if outer[0].Key != "component" || outer[0].Value.String() != "memory" {
		t.Errorf("first attribute = %v, want component=memory: it was added before any group was open, so it "+
			"stays at the top level", outer[0])
	}
	group := outer[1]
	if group.Key != "ctx" || group.Value.Kind() != slog.KindGroup {
		t.Fatalf("second attribute = %q kind %v, want a group named ctx", group.Key, group.Value.Kind())
	}
	inner := group.Value.Group()
	if len(inner) != 2 {
		t.Fatalf("group ctx holds %d attributes, want attempt and job merged into the one group they share: %v",
			len(inner), inner)
	}
	if inner[0].Key != "attempt" || inner[0].Value.Int64() != 2 {
		t.Errorf("first attribute in ctx = %v, want attempt=2 (chain attributes come before the call's own)", inner[0])
	}
	if inner[1].Key != "job" || inner[1].Value.String() != "prune" {
		t.Errorf("second attribute in ctx = %v, want job=prune (the call's own attribute last)", inner[1])
	}
}

// TestEarlyTeeNeverStopsMirroring is the "a-with-mirror" promise stated as
// behaviour: the console copy is not a fallback for an unfilled buffer, it is the
// copy that exists for the case where no listener ever arrives. Before, during and
// after a drain, every enabled record reaches the mirror exactly once and the
// buffer keeps nothing extra.
func TestEarlyTeeNeverStopsMirroring(t *testing.T) {
	b := &earlyLogBuffer{}
	mirror := newCapture130(slog.LevelInfo)
	tee := earlyTee{buf: &earlyBufferHandler{b: b}, mirror: mirror}
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		if err := tee.Handle(ctx, slog.NewRecord(clock130(), slog.LevelWarn, "sealing verdict", 0)); err != nil {
			t.Fatalf("tee.Handle returned %v, want nothing propagated to the caller", err)
		}
	}
	if n := len(mirror.msgs()); n != 3 {
		t.Fatalf("mirror saw %d records before the drain, want 3", n)
	}
	if flushed, _ := b.drain(newCapture130(slog.LevelInfo)); flushed != 3 {
		t.Fatalf("drain flushed %d, want the 3 records the mirror also saw", flushed)
	}
	if err := tee.Handle(ctx, slog.NewRecord(clock130(), slog.LevelWarn, "after the drain", 0)); err != nil {
		t.Fatal(err)
	}
	if n := len(mirror.msgs()); n != 4 {
		t.Fatalf("mirror saw %d records after the drain, want 4: the console pass never switches off", n)
	}
	if n := b.held130(); n != 0 {
		t.Errorf("buffer holds %d record(s) after the drain, want 0 (the tee must not double-keep)", n)
	}
}

// TestEarlyTeeLevelGateIsTheConsoles pins that the buffer costs no record Go's
// own default would have dropped, and keeps none the console is not shown:
// Enabled delegates to the mirror, so the gate in front of this design is today's
// gate.
func TestEarlyTeeLevelGateIsTheConsoles(t *testing.T) {
	tee := earlyTee{
		buf:    &earlyBufferHandler{b: &earlyLogBuffer{}},
		mirror: slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelWarn}),
	}
	ctx := context.Background()
	if tee.Enabled(ctx, slog.LevelInfo) {
		t.Error("INFO enabled under a WARN mirror: the buffer would keep records the console never printed")
	}
	if !tee.Enabled(ctx, slog.LevelWarn) {
		t.Error("WARN disabled under a WARN mirror")
	}
}

// TestInstalledEarlyTeeMirrorsToAConcreteHandler guards the one property of
// init() that cannot be deduced from how it reads. Capturing Go's own default
// handler as the console mirror looks strictly better - the console keeps its
// exact format - and it deadlocks: that handler is a bridge into the log package,
// whose output is bridged back into slog through slog.Default(), which is by then
// this tee. internal/observe's TestUnknownNameIsLeakSymptom hung this package's
// gate for 600s on that cycle, and a case that only runs when someone "tidies"
// init() is worth having.
func TestInstalledEarlyTeeMirrorsToAConcreteHandler(t *testing.T) {
	mirror := earlyTeeInstalled.mirror
	if mirror == nil {
		t.Fatal("the tee installed by this package's init() has no console mirror at all")
	}
	got := fmt.Sprintf("%T", mirror)
	if strings.Contains(got, "defaultHandler") {
		t.Fatalf("the installed mirror is Go's stock default handler (%s): that is the log-bridge cycle, "+
			"see the comment on init()", got)
	}
	if want := fmt.Sprintf("%T", slog.NewTextHandler(io.Discard, nil)); got != want {
		t.Errorf("the installed mirror is %s, want the concrete console handler %s", got, want)
	}
	if earlyTeeInstalled.buf == nil || earlyTeeInstalled.buf.b != earlyLogs {
		t.Error("the installed tee does not buffer into the package-level early log buffer")
	}
}

func attrString130(r slog.Record, key string) string {
	var out string
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == key {
			out = a.Value.String()
			return false
		}
		return true
	})
	return out
}
