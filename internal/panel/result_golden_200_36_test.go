// Ticket 36 AC#4's golden replay instrument (the Go half; render:stream200-36
// reads the same fixture file, so the two sides cannot disagree about what the
// golden says).
//
// WHAT IT PROVES, in the shape the 10-08 ruling demands (the ticket's original
// "renders with merge" sentence was re-cut by ticket 197 leg B, approval A406,
// code 94ea50f7: Overflow TRUNCATES and never merges - internal/panel/pump.go's
// StreamLog block says so verbatim, and this file writes its criteria in that
// shape; NOT one byte of pump.go's truncation rule is changeable to satisfy the
// old word):
//
//   1. NO-JANK BY DECIDABLE FACTS ONLY (no wall clock anywhere): the replay
//      counts its own incremental publishes - one per Append / Usage / Close
//      step, exactly how the production sink pushes - and asserts
//        a. the publish count equals the replay's own step count (no step
//           silently coalesced away),
//        b. EVERY intermediate packet stays inside the log's own hard-ceiling
//           row bound (DefaultStreamKeys * StreamKeyHardCeilingMultiple) - the
//           per-publish row cap the page can rely on,
//        c. every row of every packet contains exactly its OWN chunk's unique
//           marker and never another chunk's - a merge anywhere in the replay
//           is caught at the step it happens, not only at the end,
//        d. a key that leaves the packet is gone through the named-dropped
//           ledger only (DroppedCount accounts for every disappearance).
//   2. OVERFLOW IS TRUNCATION, NEVER MERGE: 200 golden chunks through a
//      DefaultStreamKeys=32 log crosses the key bound and the hard ceiling.
//      The 64 surviving rows must be the first 64 keys in arrival order, each
//      with its own golden text VERBATIM (the chunks are sized under the
//      truncated-mode window, so clamping costs nothing and any altered rune
//      is a defect), the 136 dropped keys must be counted with the last 64
//      named, and no dropped key's marker may appear in any surviving row.
//      A second fixture part feeds one 10,000-rune stream BELOW the key bound
//      and pins the per-row shape: head 4096 + the exact "[truncated: N runes
//      elided]" marker + tail 4096, the elided count exact, the middle's own
//      marker absent (that middle is NAMED as lost, never folded anywhere).
//   3. THE TOKEN COUNTS RECONCILE PER CHUNK: every chunk the golden gives a
//      usage row to must carry that row VERBATIM off the wire (all five
//      fields), and every chunk the golden gives none to must travel WITHOUT
//      the usage key - the explicit "未知" shape, never a fabricated zero.
//      The wire half is checked on the pump's own marshalled bytes: rows with
//      usage carry four JSON keys, rows without carry exactly three.
package panel

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type goldenUsage36 struct {
	TokensIn     int    `json:"tokensIn"`
	TokensOut    int    `json:"tokensOut"`
	CachedTokens int    `json:"cachedTokens"`
	Micros       int64  `json:"micros"`
	Currency     string `json:"currency"`
}

type goldenChunk36 struct {
	Key        string         `json:"key"`
	Deltas     []string       `json:"deltas"`
	Close      bool           `json:"close"`
	Usage      *goldenUsage36 `json:"usage"`
	UsageFirst bool           `json:"usageFirst"`
}

type goldenLongRow36 struct {
	Key         string   `json:"key"`
	Deltas      []string `json:"deltas"`
	TotalRunes  int      `json:"totalRunes"`
	HeadMarker  string   `json:"headMarker"`
	MidMarker   string   `json:"midMarker"`
	TailMarker  string   `json:"tailMarker"`
	ElidedExact int      `json:"elidedExact"`
}

type golden36 struct {
	Schema              string          `json:"schema"`
	MaxKeys             int             `json:"maxKeys"`
	ExpectSurvivingKeys int             `json:"expectSurvivingKeys"`
	ExpectDroppedCount  int             `json:"expectDroppedCount"`
	Chunks              []goldenChunk36 `json:"chunks"`
	LongRow             goldenLongRow36 `json:"longRow"`
}

func loadGolden36(t *testing.T) golden36 {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(panelRepoRoot(t), "internal", "panel", "testdata", "stream200-golden-36.json"))
	if err != nil {
		t.Fatalf("read the 200-chunk golden fixture: %v", err)
	}
	var g golden36
	if err := json.Unmarshal(data, &g); err != nil {
		t.Fatalf("parse the 200-chunk golden fixture: %v", err)
	}
	if len(g.Chunks) != 200 {
		t.Fatalf("the golden carries %d chunks, want the 200 the AC names", len(g.Chunks))
	}
	return g
}

// markerOf36 is each chunk's unique marker: the replay proves no row ever
// carries two of these, which is the merge detector in both directions.
func markerOf36(key string) string {
	return "tok" + strings.TrimPrefix(key, "task-")
}

func rowKeysOf36(rows []ResultChunk) []string {
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = r.CorrelationID
	}
	return out
}

func TestResultGoldenReplay200ChunksTruncatesNeverMergesAndReconcilesTokens36(t *testing.T) {
	g := loadGolden36(t)
	sl := NewStreamLog(g.MaxKeys)
	pump := NewSnapshotPump(PumpSources{Results: sl.Chunks})

	// The unique-marker index over all 200 golden chunks.
	markers := make(map[string]string, len(g.Chunks)) // marker -> key
	for _, c := range g.Chunks {
		markers[markerOf36(c.Key)] = c.Key
	}

	// The replay: one publish per step, the shape the production sink has.
	steps := 0
	snap := func() []ResultChunk {
		steps++
		rows, err := pumpSnapshotResults36(t, pump)
		if err != nil {
			t.Fatalf("step %d: %v", steps, err)
		}
		if len(rows) > DefaultStreamKeys*StreamKeyHardCeilingMultiple {
			t.Fatalf("step %d: packet carries %d rows, past the log's own hard-ceiling bound %d - the per-publish cap the page can rely on is gone",
				steps, len(rows), DefaultStreamKeys*StreamKeyHardCeilingMultiple)
		}
		for _, row := range rows {
			own := markerOf36(row.CorrelationID)
			if !strings.Contains(row.Text, own) {
				t.Fatalf("step %d: row %s lost its own marker %q (text %q)", steps, row.CorrelationID, own, row.Text)
			}
			for m := range markers {
				if m != own && strings.Contains(row.Text, m) {
					t.Fatalf("step %d: row %s carries foreign marker %q - two chunks merged into one row, the exact shape ticket 197 leg B re-cut away", steps, row.CorrelationID, m)
				}
			}
		}
		return rows
	}

	for _, c := range g.Chunks {
		if c.UsageFirst && c.Usage != nil {
			sl.Usage(c.Key, RunUsageRow{TokensIn: c.Usage.TokensIn, TokensOut: c.Usage.TokensOut,
				CachedTokens: c.Usage.CachedTokens, Micros: c.Usage.Micros, Currency: c.Usage.Currency})
			snap()
		}
		for _, d := range c.Deltas {
			sl.Append(c.Key, d)
			snap()
		}
		if !c.UsageFirst && c.Usage != nil {
			sl.Usage(c.Key, RunUsageRow{TokensIn: c.Usage.TokensIn, TokensOut: c.Usage.TokensOut,
				CachedTokens: c.Usage.CachedTokens, Micros: c.Usage.Micros, Currency: c.Usage.Currency})
			snap()
		}
		if c.Close {
			sl.Close(c.Key)
			snap()
		}
	}

	// 1a. Every step produced its own publish - nothing coalesced away.
	wantSteps := 0
	for _, c := range g.Chunks {
		if c.UsageFirst && c.Usage != nil {
			wantSteps++
		}
		wantSteps += len(c.Deltas)
		if !c.UsageFirst && c.Usage != nil {
			wantSteps++
		}
		if c.Close {
			wantSteps++
		}
	}
	if steps != wantSteps {
		t.Fatalf("the replay published %d packets for %d replay steps - an increment was lost or coalesced", steps, wantSteps)
	}

	// 2. The final shape: truncation named, never merged. The drops evict the
	// OLDEST keys, so the survivors are the last 64 golden chunks in arrival
	// order (keys 136..199), not the first.
	chunks := sl.Chunks()
	if len(chunks) != g.ExpectSurvivingKeys {
		t.Fatalf("surviving rows = %d, want %d (the log's own hard ceiling)", len(chunks), g.ExpectSurvivingKeys)
	}
	firstSurvivor := len(g.Chunks) - g.ExpectSurvivingKeys
	for i, row := range chunks {
		want := g.Chunks[firstSurvivor+i]
		if row.CorrelationID != want.Key {
			t.Fatalf("surviving row %d is %s, want %s - arrival order did not survive", i, row.CorrelationID, want.Key)
		}
		full := strings.Join(want.Deltas, "")
		if row.Text != full {
			t.Fatalf("row %s: text drifted from the golden\n got  %q\n want %q", want.Key, row.Text, full)
		}
	}
	if !sl.Truncated() {
		t.Error("200 keys through a 32-key log reported no truncation: the bound went unreported")
	}
	if got := sl.DroppedCount(); got != g.ExpectDroppedCount {
		t.Errorf("DroppedCount = %d, want %d - a disappeared row without an account is the merge lie by another door", got, g.ExpectDroppedCount)
	}
	// Trace the drops: appending key 64 evicts key 0, ..., appending key 199
	// evicts key 135 - drops 0..135 in order, 136 total. The naming ledger
	// keeps the LAST streamDroppedNamesCap of them, oldest retained first:
	// keys 136-64 .. 135, i.e. 72..135.
	dropped := sl.DroppedKeys()
	if len(dropped) != streamDroppedNamesCap {
		t.Fatalf("DroppedKeys names %d keys, want the ledger cap %d", len(dropped), streamDroppedNamesCap)
	}
	for i, k := range dropped {
		want := g.Chunks[g.ExpectDroppedCount-streamDroppedNamesCap+i].Key
		if k != want {
			t.Fatalf("DroppedKeys[%d] = %s, want %s (the last %d drops, oldest retained first)", i, k, want, streamDroppedNamesCap)
		}
	}
	// No dropped key's content may appear in any surviving row: "nobody is
	// showing this one" is the sentence, not "somebody folded it away". The
	// dropped keys are exactly the first firstSurvivor golden chunks.
	for _, row := range chunks {
		for i := 0; i < firstSurvivor; i++ {
			if m := markerOf36(g.Chunks[i].Key); strings.Contains(row.Text, m) {
				t.Fatalf("surviving row %s carries dropped chunk %s's marker %q - a dropped stream was merged, not named", row.CorrelationID, g.Chunks[i].Key, m)
			}
		}
	}

	// 3. The token reconciliation, per chunk, on the wire's own bytes.
	raw, err := pump.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(raw, &generic); err != nil {
		t.Fatal(err)
	}
	var wireRows []map[string]json.RawMessage
	if err := json.Unmarshal(generic["results"], &wireRows); err != nil {
		t.Fatal(err)
	}
	if len(wireRows) != len(chunks) {
		t.Fatalf("wire carries %d result rows, the log holds %d", len(wireRows), len(chunks))
	}
	reconciled := 0
	for i, row := range chunks {
		gu := g.Chunks[firstSurvivor+i].Usage
		if gu == nil {
			// The explicit-unknown shape: no usage event, no key - never a
			// zero dressed as data.
			if row.Usage != nil {
				t.Errorf("chunk %s carries usage %+v but the golden gave it none - a fabricated row is not a reading", row.CorrelationID, row.Usage)
			}
			if _, present := wireRows[i]["usage"]; present {
				t.Errorf("wire row %d carries a usage key where the golden says unknown: %s", i, string(wireRows[i]["usage"]))
			}
			if len(wireRows[i]) != 3 {
				t.Errorf("a no-usage wire row carries %d keys, want exactly the original three: %v", len(wireRows[i]), wireRows[i])
			}
			continue
		}
		reconciled++
		if row.Usage == nil {
			t.Fatalf("chunk %s: the golden gave a usage event and the wire lost it", row.CorrelationID)
		}
		if row.Usage.TokensIn != gu.TokensIn || row.Usage.TokensOut != gu.TokensOut ||
			row.Usage.CachedTokens != gu.CachedTokens || row.Usage.Micros != gu.Micros ||
			row.Usage.Currency != gu.Currency {
			t.Fatalf("chunk %s usage = %+v, want the golden's %+v verbatim", row.CorrelationID, row.Usage, gu)
		}
		if len(wireRows[i]) != 4 {
			t.Errorf("a usage-carrying wire row carries %d keys, want four (the original three + usage): %v", len(wireRows[i]), wireRows[i])
		}
	}
	if reconciled == 0 {
		t.Fatal("the instrument reconciled no usage rows - the fixture's usage half went unread")
	}
	t.Logf("golden replay: %d publishes for %d steps, %d surviving rows verbatim, %d dropped counted (%d named), %d usage rows reconciled, %d unknown rows carried no usage key",
		steps, wantSteps, len(chunks), sl.DroppedCount(), len(dropped), reconciled, len(chunks)-reconciled)
}

// pumpSnapshotResults36 reads one packet off the pump - the page's view of one
// publish - with the marshal bytes parsed back so the row count it asserts is
// the count the wire actually carries.
func pumpSnapshotResults36(t *testing.T, pump *SnapshotPump) ([]ResultChunk, error) {
	t.Helper()
	raw, err := pump.Marshal()
	if err != nil {
		return nil, err
	}
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(raw, &generic); err != nil {
		return nil, err
	}
	var rows []ResultChunk
	if err := json.Unmarshal(generic["results"], &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// TestResultGoldenLongRowTruncatesBelowTheKeyBound36 pins the per-row half of
// the same rule: one 10,000-rune stream inside a 32-key log never crosses the
// key count bound, so the WIDE per-row window applies - and the row still
// keeps head 4096 + the exact named marker + tail 4096, with the elided count
// exact and the middle's own marker absent (named as lost, never merged into
// a neighbour, never silently shortened without a count).
func TestResultGoldenLongRowTruncatesBelowTheKeyBound36(t *testing.T) {
	g := loadGolden36(t)
	lr := g.LongRow
	sl := NewStreamLog(DefaultStreamKeys)
	for _, d := range lr.Deltas {
		sl.Append(lr.Key, d)
	}
	chunks := sl.Chunks()
	if len(chunks) != 1 {
		t.Fatalf("one stream must stay one row, got %d rows - that is a merge", len(chunks))
	}
	row := chunks[0]
	runes := []rune(strings.Join(lr.Deltas, ""))
	if len(runes) != lr.TotalRunes {
		t.Fatalf("golden long row = %d runes, fixture says %d", len(runes), lr.TotalRunes)
	}
	elided := lr.TotalRunes - 2*streamRowKeepRunes
	if elided != lr.ElidedExact {
		t.Fatalf("computed elision %d disagrees with the fixture's own %d", elided, lr.ElidedExact)
	}
	expected := string(runes[:streamRowKeepRunes]) +
		"...[truncated: "+strconv.Itoa(elided)+" runes elided]..."+
		string(runes[len(runes)-streamRowKeepRunes:])
	if row.Text != expected {
		t.Fatalf("the long row's render drifted from the truncation shape\n got  (len %d)\n want (len %d)\n got head %q\n want head %q",
			len([]rune(row.Text)), len([]rune(expected)), truncate36(row.Text, 40), truncate36(expected, 40))
	}
	if !strings.Contains(row.Text, lr.HeadMarker) || !strings.Contains(row.Text, lr.TailMarker) {
		t.Error("the rendered row lost its own head or tail marker")
	}
	if strings.Contains(row.Text, lr.MidMarker) {
		t.Error("the rendered row carries the MIDDLE's marker - the elided middle leaked back in")
	}
	if sl.Truncated() {
		t.Error("one stream through a 32-key log is below the key bound: the log-level flag must stay false")
	}
	tr := sl.TruncationFor(lr.Key)
	if !tr.Truncated || tr.ElidedRunes != elided {
		t.Errorf("TruncationFor = %+v, want Truncated=true ElidedRunes=%d - 'truncated' that cannot say how much it lost is the lie with a label", tr, elided)
	}
	if sl.ElidedRunes() != elided {
		t.Errorf("ElidedRunes = %d, want %d", sl.ElidedRunes(), elided)
	}
	// No usage event was ever fed for this stream: the row must travel
	// without the key, the explicit unknown.
	if row.Usage != nil {
		t.Errorf("a stream with no usage event carries usage %+v - a fabricated zero's cousin", row.Usage)
	}
	raw, err := NewSnapshotPump(PumpSources{Results: sl.Chunks}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(raw, &generic); err != nil {
		t.Fatal(err)
	}
	var wireRows []map[string]json.RawMessage
	if err := json.Unmarshal(generic["results"], &wireRows); err != nil {
		t.Fatal(err)
	}
	if len(wireRows) != 1 || len(wireRows[0]) != 3 {
		t.Fatalf("the no-usage wire row carries %d keys over %d rows, want exactly the original three over one row: %v", len(wireRows[0]), len(wireRows), wireRows[0])
	}
}

func truncate36(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "..."
}
