package panel

// Ticket 197 leg B - the stream layer: one subagent, one stream, and overflow
// that no longer buys its bound by folding two producers into one row.
//
// What each case pins, and which mistake it can catch:
//
//	TestSubagentStreamKeyIsTheOneSpelling          §0's key shape is one spelling;
//	    a second shape appearing later means two rows for one agent, or one row for
//	    two agents.
//	TestNPlusOneSubagentsEachGetTheirOwnStream     the structural ruler: N+1 agents
//	    past the key bound still produce N+1 rows, none of them somebody else's
//	    stream with text bolted on. This is the case the old merge semantic failed.
//	TestTwoSubagentsWritingAtOnceDoNotPolluteEachOther
//	    the positive control on interleaving: two agents writing into one process at
//	    the same time, each row exactly its own deltas. Had leg B kept "merge", this
//	    is the case that goes red.
//	TestOverflowTruncatesEachStreamsOwnHeadAndTail  the chosen branch (乙): the bound
//	    is kept by eliding the MIDDLE of each stream and saying so in that stream's
//	    own text. This is the ruler that can go red for the truncation branch.
//	TestHardCeilingDropsWholeKeysAndNamesThem      the bound's last resort drops a
//	    whole key and NAMES it instead of folding it into a live row.
//
// What this file does NOT test is named at the bottom.

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"
)

func TestSubagentStreamKeyIsTheOneSpelling(t *testing.T) {
	if got := SubagentStreamKey("t-42"); got != "subagent:t-42" {
		t.Errorf("SubagentStreamKey(\"t-42\") = %q, want the literal shape subagent:<taskID>", got)
	}
	// The task id goes through verbatim: the panel must be able to point back at
	// the roster row it came from. Lowercasing an id here would silently rename a
	// task, and a colon in an id must not be mistaken for a second separator.
	if got := SubagentStreamKey("Task-7:inner"); got != "subagent:Task-7:inner" {
		t.Errorf("verbatim passthrough broken: %q", got)
	}
	for _, empty := range []string{"", "   "} {
		if got := SubagentStreamKey(empty); got != "" {
			t.Errorf("SubagentStreamKey(%q) = %q, want \"\": a stream owned by no task is not a stream", empty, got)
		}
	}
	// The prefix is this package's constant, not a copy of a string someone typed
	// in the tools leg: if the two ever disagree, the panel reads an untagged row.
	if !strings.HasPrefix(SubagentStreamKey("x"), SubagentStreamKeyPrefix) {
		t.Error("SubagentStreamKey does not start with its own prefix constant")
	}
}

func TestNPlusOneSubagentsEachGetTheirOwnStream(t *testing.T) {
	const bound = 3
	sl := NewStreamLog(bound)

	agents := []string{"alpha", "bravo", "charlie", "delta"} // bound + 1
	for _, name := range agents {
		sl.Append(SubagentStreamKey(name), "["+name+" step1] ")
		sl.Append(SubagentStreamKey(name), "["+name+" step2]")
		sl.Close(SubagentStreamKey(name))
	}

	chunks := sl.Chunks()
	if len(chunks) != len(agents) {
		t.Fatalf("rows = %d for %d subagents, want one row per agent - %s",
			len(chunks), len(agents), describeRows(chunks))
	}
	byKey := map[string]ResultChunk{}
	for _, c := range chunks {
		if _, clash := byKey[c.CorrelationID]; clash {
			t.Errorf("key %q appears twice in the log, want one row per key", c.CorrelationID)
		}
		byKey[c.CorrelationID] = c
	}
	for _, name := range agents {
		key := SubagentStreamKey(name)
		c, ok := byKey[key]
		if !ok {
			t.Errorf("no row for %q: that agent's work is not on screen at all - %s", key, describeRows(chunks))
			continue
		}
		want := "[" + name + " step1] [" + name + " step2]"
		if c.Text != want {
			t.Errorf("row %q = %q, want exactly its own deltas %q", key, c.Text, want)
		}
		// The folded shape of the old rule: one row carrying another agent's
		// sentence. Assert it per row, not just on the row count.
		for _, other := range agents {
			if other != name && strings.Contains(c.Text, other) {
				t.Errorf("row %q contains %q's text: two agents' streams became one row: %q", key, other, c.Text)
			}
		}
		if !c.Done {
			t.Errorf("row %q done=false, want the Close that was called on THIS key", key)
		}
	}

	// Past the bound the log says it truncated (the truth), and yet short streams
	// lost nothing - "what does not break when this is fixed" is exactly this line:
	// the bound moved from row count to runes-per-row, it did not eat content.
	if !sl.Truncated() {
		t.Error("4 keys through a 3-key bound reported no truncation: the bound went unreported")
	}
	if got := sl.ElidedRunes(); got != 0 {
		t.Errorf("elided runes = %d, want 0 - every stream here fits its head+tail: %s", got, describeRows(chunks))
	}
	if got := sl.DroppedKeys(); len(got) != 0 {
		t.Errorf("dropped keys = %v, want none under the hard ceiling", got)
	}

	// And the packet still carries them as separate rows with the same four
	// top-level keys - this leg adds no contract key (that is 197-r3's carrier).
	pump := NewSnapshotPump(PumpSources{
		Results: func() []ResultChunk { return sl.Chunks() },
	})
	data, err := pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(data, &generic); err != nil {
		t.Fatal(err)
	}
	if len(generic) != 4 {
		t.Errorf("snapshot top-level keys = %d, want the four PanelSnapshot declares: %v", len(generic), generic)
	}
	var results []map[string]json.RawMessage
	if err := json.Unmarshal(generic["results"], &results); err != nil {
		t.Fatal(err)
	}
	if len(results) != len(agents) {
		t.Errorf("packet has %d result rows for %d subagents: %s", len(results), len(agents), describeRows(chunks))
	}
	for _, r := range results {
		if len(r) != 3 {
			t.Errorf("a result row carries %d JSON keys, want the three ResultChunk has: %v", len(r), r)
		}
	}
}

func TestTwoSubagentsWritingAtOnceDoNotPolluteEachOther(t *testing.T) {
	sl := NewStreamLog(DefaultStreamKeys)

	// Three writers, one of them the root task, interleaved on purpose. Each key is
	// written by exactly one goroutine, which is the shape a real fan-out has: the
	// claim under test is not "concurrency is safe" (Append holds a mutex), it is
	// "no row ends up holding another stream's sentence".
	var wg sync.WaitGroup
	spawn := func(key string, deltas []string) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { recover() }() // the owner/recover shape D22/D38b asks for
			for _, d := range deltas {
				sl.Append(key, d)
			}
			sl.Close(key)
		}()
	}
	spawn(SubagentStreamKey("a1"), []string{"A1 ", "A2 ", "A3 "})
	spawn(SubagentStreamKey("b2"), []string{"B1 ", "B2 ", "B3 "})
	spawn("task-root", []string{"R1 ", "R2 ", "R3 "})
	wg.Wait()

	want := map[string]string{
		SubagentStreamKey("a1"): "A1 A2 A3 ",
		SubagentStreamKey("b2"): "B1 B2 B3 ",
		"task-root":             "R1 R2 R3 ",
	}
	// One sentinel per stream: its deltas all start with that letter, so a row
	// containing another stream's sentinel is a row holding two agents' text.
	sentinel := map[string]string{
		SubagentStreamKey("a1"): "A",
		SubagentStreamKey("b2"): "B",
		"task-root":             "R",
	}
	chunks := sl.Chunks()
	if len(chunks) != len(want) {
		t.Fatalf("rows = %d, want %d - %s", len(chunks), len(want), describeRows(chunks))
	}
	for _, c := range chunks {
		exp, ok := want[c.CorrelationID]
		if !ok {
			t.Errorf("row for unknown key %q: %q", c.CorrelationID, c.Text)
			continue
		}
		if c.Text != exp {
			t.Errorf("row %q = %q, want %q", c.CorrelationID, c.Text, exp)
		}
		for other, marker := range sentinel {
			if other != c.CorrelationID && strings.Contains(c.Text, marker) {
				t.Errorf("row %q carries text from stream %q: %q", c.CorrelationID, other, c.Text)
			}
		}
		if !c.Done {
			t.Errorf("row %q done=false after its own Close", c.CorrelationID)
		}
	}
	if sl.Truncated() {
		t.Error("3 keys inside a 32-key bound reported truncation")
	}
}

func TestOverflowTruncatesEachStreamsOwnHeadAndTail(t *testing.T) {
	const bound = 2
	sl := NewStreamLog(bound)

	// Three long-lived subagents through a 2-key bound. 904 runes each: the middle
	// is what the bound is allowed to take, per stream, out of that stream's own
	// text - and the head/tail of one agent must never become another's row.
	heads := map[string]string{"s1": "HEAD-ONE", "s2": "HEAD-TWO", "s3": "HEAD-THREE"}
	tails := map[string]string{"s1": "TAIL-ONE", "s2": "TAIL-TWO", "s3": "TAIL-THREE"}
	for _, id := range []string{"s1", "s2", "s3"} {
		sl.Append(SubagentStreamKey(id), heads[id])
		sl.Append(SubagentStreamKey(id), strings.Repeat("汉", 900))
		sl.Append(SubagentStreamKey(id), tails[id])
		sl.Close(SubagentStreamKey(id))
	}

	if !sl.Truncated() {
		t.Fatal("3 keys through a 2-key bound did not report truncation")
	}
	chunks := sl.Chunks()
	if len(chunks) != 3 {
		t.Fatalf("rows = %d, want 3 (one per subagent, past the bound) - %s", len(chunks), describeRows(chunks))
	}
	totalElided := 0
	for _, c := range chunks {
		id := strings.TrimPrefix(c.CorrelationID, SubagentStreamKeyPrefix)
		if !strings.HasPrefix(c.Text, heads[id]) {
			t.Errorf("row %q lost its own head: %q", c.CorrelationID, truncateForReport(c.Text))
		}
		if !strings.HasSuffix(c.Text, tails[id]) {
			t.Errorf("row %q lost its own tail: %q", c.CorrelationID, truncateForReport(c.Text))
		}
		if !utf8.ValidString(c.Text) {
			t.Errorf("row %q is not valid UTF-8 - the elision cut a rune in half: %q",
				c.CorrelationID, truncateForReport(c.Text))
		}
		// Every row lost its own middle and only its own middle: the marker counts
		// this stream's runes, not some neighbour's, and it appears ONCE - a re-fold
		// over an already folded row would stack a second marker inside the first.
		produced := utf8.RuneCountInString(heads[id]) + 900 + utf8.RuneCountInString(tails[id])
		elided := produced - 2*streamTruncateKeepRunes
		wantMarker := fmt.Sprintf("...[truncated: %d runes elided]...", elided)
		if got := strings.Count(c.Text, wantMarker); got != 1 {
			t.Errorf("row %q: marker %q appears %d times, want exactly 1: %q",
				c.CorrelationID, wantMarker, got, truncateForReport(c.Text))
		}
		if got := strings.Count(c.Text, "runes elided"); got != 1 {
			t.Errorf("row %q carries %d truncation markers, want 1 (the fold was re-run over its own output): %q",
				c.CorrelationID, got, truncateForReport(c.Text))
		}
		totalElided += elided
		for other, h := range heads {
			if other != id && (strings.Contains(c.Text, h) || strings.Contains(c.Text, tails[other])) {
				t.Errorf("row %q holds stream %q's head or tail: two streams became one row", c.CorrelationID, other)
			}
		}
	}
	if sl.ElidedRunes() != totalElided {
		t.Errorf("ElidedRunes = %d, want %d - 'truncated' that cannot say how much it lost is the old lie with a label on it",
			sl.ElidedRunes(), totalElided)
	}
	if got := sl.DroppedKeys(); len(got) != 0 {
		t.Errorf("dropped keys = %v, want none: 3 keys is under the 2x ceiling", got)
	}
}

func TestHardCeilingDropsWholeKeysAndNamesThem(t *testing.T) {
	const bound = 1
	ceiling := bound * StreamKeyHardCeilingMultiple
	sl := NewStreamLog(bound)

	// 5 closed subagents through a 1-key bound: the log may hold at most 2 rows,
	// the rest are dropped - and a dropped stream is named, not folded into a row
	// the panel is still showing.
	for i := 1; i <= 5; i++ {
		key := SubagentStreamKey(fmt.Sprintf("t%d", i))
		sl.Append(key, fmt.Sprintf("only-t%d", i))
		sl.Close(key)
	}
	chunks := sl.Chunks()
	if len(chunks) > ceiling {
		t.Fatalf("rows = %d, want at most the %d-key hard ceiling - %s", len(chunks), ceiling, describeRows(chunks))
	}
	dropped := sl.DroppedKeys()
	if len(dropped) != 5-len(chunks) {
		t.Errorf("dropped = %v, want exactly the %d keys no row covers", dropped, 5-len(chunks))
	}
	seen := map[string]bool{}
	for _, c := range chunks {
		seen[c.CorrelationID] = true
		for _, d := range dropped {
			if d == c.CorrelationID {
				t.Errorf("key %q is named as dropped AND still on screen", d)
			}
			if strings.Contains(c.Text, "only-"+strings.TrimPrefix(d, SubagentStreamKeyPrefix)) {
				t.Errorf("dropped stream %q was folded into live row %q instead of being named: %q",
					d, c.CorrelationID, c.Text)
			}
		}
	}
	for _, c := range chunks {
		if !strings.HasPrefix(c.Text, "only-") {
			t.Errorf("surviving row %q = %q, want its own delta and nobody else's", c.CorrelationID, c.Text)
		}
	}
}

// --------------------------------------------------------------------------
// WHAT THIS FILE DID NOT MEASURE (ticket 197 leg B, honest gaps):
//
//  1. Nothing here drives the real fan-out: the tools leg (197-r1) is what feeds
//     these keys, and this package has no subagent producer in it. What is pinned
//     is the carrier's behaviour for the §0 key SHAPE, so a tools leg that keys
//     differently would still get one row per key - a wrong key would not be caught
//     here, only in whatever test feeds the pump from a real run.
//  2. The truncation marker and the dropped-key list are Go-side accessors only.
//     No case here asserts the panel SHOWS them, because this leg adds no
//     ResultChunk / snapshot key (frontend/** is a zero-write surface and
//     TestComposerContractTypesMatchFrontend reconciles that pair); the visible
//     affordance is 197-r3's carrier layer, and until it lands the page still reads
//     a truncated row as if it were a whole answer.
//  3. The bound's memory effect is asserted as "runes per row" and "rows per log",
//     never as bytes of a live snapshot: no case here marshals a saturated log and
//     measures the packet size, so a regression that keeps both counts and still
//     grows (e.g. many keys, each at the head+tail floor) would not be judged red.
//  4. Key growth is not tied to the §0 concurrency cap of 8: StreamLog has no idea
//     how many tasks the host considers active, so option 甲's "cap by active
//     task count" is not represented anywhere here, and nothing checks that a run
//     with more than 32 live subagents is even possible.
//  5. Close() on a key that never appeared still opens a row (pre-existing
//     behaviour, kept). No case pins the interaction of that with the hard ceiling
//     drop order, and Append("") with a non-empty text is still accepted - the
//     unattributed row that produces is unchanged by this leg on purpose.

func describeRows(chunks []ResultChunk) string {
	var b strings.Builder
	for _, c := range chunks {
		fmt.Fprintf(&b, " {%q text=%q done=%v}", c.CorrelationID, truncateForReport(c.Text), c.Done)
	}
	return b.String()
}

func truncateForReport(s string) string {
	r := []rune(s)
	if len(r) <= 48 {
		return s
	}
	return string(r[:24]) + "<...>" + string(r[len(r)-24:])
}
