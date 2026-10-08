// Package panel test - ticket 35, cell ":49" Backpressure, leg 35-r8 (writer).
//
// The ticket's AC sentence reads "flood events under blocked consumer -> merges,
// no unbounded memory (heap cap asserted), content integrity kept." Its first and
// third thirds were re-ruled by ticket 197 leg B (pump.go, verbatim: "Overflow
// TRUNCATES and never merges"), so this file owns the ONE taste the cell still
// accounts for: is what a blocked consumer pins in memory actually bounded,
// measured as behaviour and not as the existence of a constant.
//
// Nailed elsewhere, therefore NOT redone here:
//   - the row-COUNT bound and the dropped naming, pump_test.go
//     TestTheStreamLogTruncatesInsteadOfMerging - it compares len(chunks) against
//     2*StreamKeyHardCeilingMultiple, which is a count against a constant, in bytes
//     of nothing;
//   - the per-row truncation ACCOUNTING, subagent_roster_197_test.go.
//
// No ruler in this package has ever read the runtime's memory stats (this leg's
// dedup grep: HeapAlloc / runtime.MemStats / MemStats across internal/panel = 0
// hits), so "heap cap asserted" had no instrument at all. This file is that
// instrument, and it reports what it finds:
//   - the fan-out bound is real, sub-linear and counted (GREEN);
//   - a flood that never crosses the key bound pins every rune it is given, and the
//     dropped-key ledger pins one name per key the flood opened (both RED, and both
//     are the honest reading: making them green needs a production decision, which
//     this leg may not make).
package panel

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"
)

// ---------------------------------------------------------------- measuring --

// floodMeasurement35r8 is what one structure holds after a flood, in the units the
// AC can be read in: the runes of display text it pins, how many rows that text is
// spread over, the runes of KEY NAMES it pins for the keys it stopped tracking, and
// the live heap the flood bought. Names are memory too - a structure that gets its
// rows down by remembering every row it dropped is not bounded.
type floodMeasurement35r8 struct {
	textRunes int
	rows      int
	nameRunes int
	names     int
	heapDelta uint64
}

// floodTarget35r8 is the seam the flood is driven through. The production log and
// the positive control's sink both answer it, so one harness and one predicate
// judge both - which is the only reason the control below can say anything.
type floodTarget35r8 interface {
	Push(key, delta string)
	Measure() (textRunes, rows, nameRunes, names int)
	Snapshot() []ResultChunk
}

// streamLogTarget35r8 wraps the production log.
type streamLogTarget35r8 struct{ log *StreamLog }

func (s streamLogTarget35r8) Push(key, delta string) { s.log.Append(key, delta) }

func (s streamLogTarget35r8) Measure() (int, int, int, int) {
	text, rows, names := 0, 0, 0
	for _, c := range s.log.Chunks() {
		text += len([]rune(c.Text))
		rows++
	}
	dropped := s.log.DroppedKeys()
	for _, d := range dropped {
		names += len([]rune(d))
	}
	return text, rows, names, len(dropped)
}

func (s streamLogTarget35r8) Snapshot() []ResultChunk { return s.log.Chunks() }

// unboundedSink35r8 is that same shape with the bound taken out on purpose: it
// merges nothing, clamps nothing, caps nothing and remembers every name. It exists
// so the positive control can prove the predicate discriminates. 35-v7's verdict
// names the failure shape this guards against (a ruler with no differential face
// reads green however badly the thing under it behaves).
type unboundedSink35r8 struct {
	texts map[string]string
	order []string
}

func newUnboundedSink35r8() *unboundedSink35r8 {
	return &unboundedSink35r8{texts: map[string]string{}}
}

func (u *unboundedSink35r8) Push(key, delta string) {
	if _, ok := u.texts[key]; !ok {
		u.texts[key] = delta
		u.order = append(u.order, key)
		return
	}
	u.texts[key] += delta
}

func (u *unboundedSink35r8) Measure() (int, int, int, int) {
	text := 0
	for _, k := range u.order {
		text += len([]rune(u.texts[k]))
	}
	return text, len(u.order), 0, 0
}

func (u *unboundedSink35r8) Snapshot() []ResultChunk {
	out := make([]ResultChunk, 0, len(u.order))
	for _, k := range u.order {
		out = append(out, ResultChunk{CorrelationID: k, Text: u.texts[k]})
	}
	return out
}

// heapInUse35r8 reads live heap, GC'd first so the number is retention and not
// allocator slack.
func heapInUse35r8() uint64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc
}

// flood35r8 drives keys x perKey deltas of deltaLen runes through a target while its
// consumer is provably blocked, and returns the retention. The block is not a
// comment: the exit is the real PumpSources.Out, it does a blocking send that the
// page never picks up, every refused attempt is required to still be counted by
// Publishes(), and zero attempts is itself a failure (a flood that never reached the
// exit proves nothing about a consumer).
func flood35r8(t *testing.T, target floodTarget35r8, keys, perKey, deltaLen int) floodMeasurement35r8 {
	t.Helper()

	page := make(chan []ResultChunk) // nobody ever reads from it
	reads := 0
	pump := NewSnapshotPump(PumpSources{
		Results: func() []ResultChunk { return target.Snapshot() },
		Out: func(_ Snapshot, _ []byte) error {
			select {
			case page <- nil:
				reads++
				return nil
			case <-time.After(10 * time.Millisecond):
				return fmt.Errorf("35r8: the page never picked the snapshot up")
			}
		},
	})

	before := heapInUse35r8()
	delta := strings.Repeat("x", deltaLen)
	const maxPublishAttempts = 24
	pushed, attempts, refused := 0, 0, 0
	for k := 0; k < keys; k++ {
		key := fmt.Sprintf("subagent:%d", k)
		for i := 0; i < perKey; i++ {
			target.Push(key, delta)
			pushed++
			if pushed%64 == 0 && attempts < maxPublishAttempts {
				attempts++
				if _, _, err := pump.Publish(); err != nil {
					refused++
				}
			}
		}
	}
	after := heapInUse35r8()

	if attempts == 0 {
		t.Fatalf("the flood never reached the snapshot exit (%d pushes), so the blocked consumer was never exercised", pushed)
	}
	if n := pump.Publishes(); n != attempts {
		t.Fatalf("blocked consumer: the pump counted %d publishes for %d attempts - a refused exit that stops being counted is a silent drop", n, attempts)
	}
	if refused != attempts {
		t.Fatalf("blocked consumer: %d of %d snapshots were picked up while no page was reading, want 0 (reads=%d)", attempts-refused, attempts, reads)
	}
	select {
	case p := <-page:
		_ = p
		t.Fatal("blocked consumer: the page channel has an unread snapshot in it, so something read it")
	default:
	}

	text, rows, nameRunes, names := target.Measure()
	return floodMeasurement35r8{
		textRunes: text, rows: rows, nameRunes: nameRunes, names: names,
		heapDelta: after - before,
	}
}

// retentionCap35r8 is the AC's shape: a cap whose every term comes from the bound
// and NONE of them from the flood size. Rows are bounded by the hard ceiling, each
// row keeps its own head and tail (streamTruncateKeepRunes each) plus that row's own
// inline elision marker, and the naming ledger is charged one row's width per key so
// a "bound" that only pays for names still gets caught.
func retentionCap35r8() floodMeasurement35r8 {
	rows := DefaultStreamKeys * StreamKeyHardCeilingMultiple
	perRow := 2*streamTruncateKeepRunes + 64 // "...[truncated: N runes elided]..."
	return floodMeasurement35r8{textRunes: rows * perRow, rows: rows}
}

// overCap35r8 is the one predicate this file's tests share. It returns the reading
// instead of failing the test so the positive control can ask the same question
// about a deliberately unbounded sink and require the opposite answer.
func overCap35r8(m floodMeasurement35r8, cap floodMeasurement35r8) string {
	var hits []string
	if m.rows > cap.rows {
		hits = append(hits, fmt.Sprintf("rows %d > the hard ceiling %d", m.rows, cap.rows))
	}
	if m.textRunes > cap.textRunes {
		hits = append(hits, fmt.Sprintf("pinned text %d runes > the cap %d", m.textRunes, cap.textRunes))
	}
	if m.nameRunes > cap.textRunes {
		hits = append(hits, fmt.Sprintf("pinned dropped-name runes %d > the cap %d", m.nameRunes, cap.textRunes))
	}
	heapCap := uint64(cap.textRunes)*4 + (1 << 20) // 1 MiB of slack for map and GC bookkeeping
	if m.heapDelta > heapCap {
		hits = append(hits, fmt.Sprintf("live heap + %d bytes > the cap %d", m.heapDelta, heapCap))
	}
	if len(hits) == 0 {
		return ""
	}
	return strings.Join(hits, "; ")
}

// ------------------------------------------------------------------- tests ---

// TestStreamLogFanOutFloodKeepsRetentionUnderCap35r8 is the taste the cell still
// owns, measured where a bound is supposed to exist: fan-out past the key bound,
// every stream long, consumer blocked. Retention must come back the same size
// whether 400 or 1600 streams were flooded, and it must still be holding rows with
// content and counted loss - a cap met by dropping everything silently would pass
// the numbers and fail the checks at the bottom.
func TestStreamLogFanOutFloodKeepsRetentionUnderCap35r8(t *testing.T) {
	cap := retentionCap35r8()

	one := NewStreamLog(DefaultStreamKeys)
	m1 := flood35r8(t, streamLogTarget35r8{one}, 400, 60, 160)
	four := NewStreamLog(DefaultStreamKeys)
	m4 := flood35r8(t, streamLogTarget35r8{four}, 1600, 60, 160)

	if bad := overCap35r8(m1, cap); bad != "" {
		t.Errorf("400 streams x 60 deltas x 160 runes through a %d-key log is not bounded: %s (flooded %d runes in)",
			DefaultStreamKeys, bad, 400*60*160)
	}
	if bad := overCap35r8(m4, cap); bad != "" {
		t.Errorf("1600 streams under the same log is not bounded: %s (flooded %d runes in)", bad, 1600*60*160)
	}
	if m4.textRunes > m1.textRunes*5/4+1000 {
		t.Errorf("retention scaled with the flood: 1x pinned %d runes, 4x pinned %d - a cap that grows with the event count is not a cap",
			m1.textRunes, m4.textRunes)
	}

	// The bound must be paid for with content it can name, not with silence.
	if m4.rows != cap.rows {
		t.Errorf("rows = %d, want exactly the hard ceiling %d: fewer rows than the ceiling means the flood was dropped rather than bounded", m4.rows, cap.rows)
	}
	for _, c := range four.Chunks() {
		if len([]rune(c.Text)) < 2*streamTruncateKeepRunes {
			t.Fatalf("row %q holds %d runes, below the head+tail it promises: %q", c.CorrelationID, len([]rune(c.Text)), c.Text)
		}
		if !strings.Contains(c.Text, "truncated") {
			t.Fatalf("row %q carries no elision marker although the log is truncating, so its loss is silent: %q", c.CorrelationID, c.Text)
		}
	}
	if !four.Truncated() || four.ElidedRunes() == 0 {
		t.Errorf("counting exits silent under a 1600-stream flood: Truncated=%v ElidedRunes=%d", four.Truncated(), four.ElidedRunes())
	}
	if got := len(four.DroppedKeys()); got != 1600-cap.rows {
		t.Errorf("dropped naming = %d keys, want %d: every key the log stopped tracking must be named", got, 1600-cap.rows)
	}
	if one.ElidedRunes() == 0 || len(one.DroppedKeys()) != 400-cap.rows {
		t.Errorf("the 1x flood's counting exits disagree with its rows: ElidedRunes=%d dropped=%d", one.ElidedRunes(), len(one.DroppedKeys()))
	}
	t.Logf("35r8 fan-out reading: 1x text=%d rows=%d names=%d heap+%dB | 4x text=%d rows=%d names=%d heap+%dB | cap text=%d rows=%d",
		m1.textRunes, m1.rows, m1.names, m1.heapDelta, m4.textRunes, m4.rows, m4.names, m4.heapDelta, cap.textRunes, cap.rows)
}

// TestStreamLogFloodBelowKeyBoundIsNotBounded35r8 is a RED deliverable on purpose,
// and it is the reading the census pointed at without measuring. This package's
// bound engages on KEY COUNT, not on VOLUME: pump.go's Append clamps a row only
// while s.truncated is set, and s.truncated is set only by enforceBoundLocked once
// the key count passes maxKeys. So a flood that never opens a 33rd stream - one
// long answer, or a handful, the ordinary shape of a run - accumulates every rune it
// is given while the page is not reading, and there is no heap cap to assert. The
// numbers are printed in the failure. Closing this needs a production decision
// (clamp per-row runes below the key bound too, with the elision counted, or refuse
// at the exit); this leg may not make it, and the red must NEVER be "closed" by
// shrinking the flood until the assertion passes.
func TestStreamLogFloodBelowKeyBoundIsNotBounded35r8(t *testing.T) {
	cap := retentionCap35r8()

	one := NewStreamLog(DefaultStreamKeys)
	m1 := flood35r8(t, streamLogTarget35r8{one}, 3, 200, 2000)
	four := NewStreamLog(DefaultStreamKeys)
	m4 := flood35r8(t, streamLogTarget35r8{four}, 3, 800, 2000)

	if m1.textRunes != 3*200*2000 || m4.textRunes != 3*800*2000 {
		t.Fatalf("measurement sanity: 1x pinned %d runes want %d, 4x pinned %d want %d - the harness stopped measuring the flood",
			m1.textRunes, 3*200*2000, m4.textRunes, 3*800*2000)
	}
	bad := overCap35r8(m4, cap)
	if bad == "" {
		t.Fatal("3 streams under a 32-key bound are now capped: production grew a per-row volume bound and this red ruler should turn green")
	}
	t.Errorf("flood under a blocked consumer with no fan-out (3 streams, so the key bound never engages): 1x pinned %d text runes over %d rows (live heap +%d bytes), 4x pinned %d runes over %d rows (live heap +%d bytes), against the cap this package can state (%d runes over %d rows). Retention grew %.1fx for a 4x flood and equals the flood exactly, so ticket 35 :49's \"no unbounded memory (heap cap asserted)\" has no counterpart in internal/panel below the key bound. %s",
		m1.textRunes, m1.rows, m1.heapDelta, m4.textRunes, m4.rows, m4.heapDelta,
		cap.textRunes, cap.rows, float64(m4.textRunes)/float64(m1.textRunes), bad)
}

// TestStreamLogDroppedNamingLedgerIsNotBounded35r8 is the second RED reading, on the
// face the existing row-count rulers structurally cannot see: past the hard ceiling
// the log keeps ONE NAME per key it stopped tracking, for as long as the process
// lives, and DroppedKeys() hands the whole ledger out on every read. Rows stay at
// the ceiling while the bytes scale with the flood, so "bounded" in the count sense
// and "bounded" in the AC's sense are two different claims.
func TestStreamLogDroppedNamingLedgerIsNotBounded35r8(t *testing.T) {
	cap := retentionCap35r8()

	one := NewStreamLog(DefaultStreamKeys)
	m1 := flood35r8(t, streamLogTarget35r8{one}, 4000, 1, 1)
	four := NewStreamLog(DefaultStreamKeys)
	m4 := flood35r8(t, streamLogTarget35r8{four}, 16000, 1, 1)

	if m1.names != 4000-cap.rows || m4.names != 16000-cap.rows {
		t.Fatalf("measurement sanity: 1x named %d keys want %d, 4x named %d want %d - the ledger being measured is not the one the log writes",
			m1.names, 4000-cap.rows, m4.names, 16000-cap.rows)
	}
	bad := overCap35r8(m4, cap)
	if bad == "" {
		t.Fatal("the dropped-name ledger is capped now: production bounded what it names and this red ruler should turn green")
	}
	t.Errorf("blocked consumer, fan-out past the ceiling: 1x opened 4000 streams and pins %d names (%d runes, live heap +%d bytes), 4x opened 16000 and pins %d names (%d runes, live heap +%d bytes) - rows stayed at the ceiling %d in both, so the memory the flood bought is the naming ledger itself. Its cap is the event count, not the bound. %s",
		m1.names, m1.nameRunes, m1.heapDelta, m4.names, m4.nameRunes, m4.heapDelta, cap.rows, bad)
}

// TestFloodRulerFiresOnAnUnboundedSink35r8 is the positive control the repo requires
// of every negative ruler (plant the shape, the ruler must answer). The same
// harness, the same predicate, a sink with the bound taken out on purpose: the
// predicate must say OVER cap, and the sink must still be holding exactly the bytes
// it was given so the "over" cannot come from a harness that stopped measuring.
// Neutralise overCap35r8 (return "") or the measurement (return 0) and this test is
// the one that goes red first.
func TestFloodRulerFiresOnAnUnboundedSink35r8(t *testing.T) {
	cap := retentionCap35r8()

	sink := newUnboundedSink35r8()
	m := flood35r8(t, sink, 100, 50, 160)

	if m.textRunes != 100*50*160 || m.rows != 100 {
		t.Fatalf("the harness is not measuring the planted sink: pinned %d runes over %d rows, want %d over 100 - a control that measures nothing proves nothing",
			m.textRunes, m.rows, 100*50*160)
	}
	if bad := overCap35r8(m, cap); bad == "" {
		t.Fatal("the ruler stayed silent on an unbounded sink: 800000 runes of pinned text with the bound removed on purpose read as within cap, so every red in this file is decoration")
	} else {
		t.Logf("35r8 positive control fired as required: %s (heap +%d bytes)", bad, m.heapDelta)
	}

	// And the predicate is not simply "the flood is big": the same 100 streams read
	// through the production log stay under cap, so the control's fire came from the
	// missing bound, not from the harness's own size.
	bounded := NewStreamLog(DefaultStreamKeys)
	mb := flood35r8(t, streamLogTarget35r8{bounded}, 100, 50, 160)
	if bad := overCap35r8(mb, cap); bad != "" {
		t.Errorf("the same-shaped flood is over cap through the production log (%s), so the control's fire discriminates nothing", bad)
	}
}
