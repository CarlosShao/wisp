package panel

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// Ticket 145's closing-leg tests for the run key. The reader-gated shape is
// the 213/187/65 procedure, so these teeth follow the same three questions:
//
//   - does the key travel ONLY from its reader (reader present -> key on the
//     wire; reader absent -> key absent, the exact four-key shape the byte
//     nails pin - this arm is the mutation tooth, and the fill line in
//     Snapshot() is the other half of the same tooth);
//   - does an empty feed render "tools":[] rather than null or a vanished
//     key (the commands/modelSwitch rule);
//   - does every value on the wire come from what the events wrote (AC#6's
//     value-from-true-source question) - the feed is scripted with the exact
//     values loop.go's publish sites hand an event, and the section is
//     asserted verbatim against that script, never against a helper's
//     restatement.

// TestRunSectionTravelsOnlyFromItsReader pins the key's gate: present with a
// reader, absent without one, "tools":[] when the feed is empty.
func TestRunSectionTravelsOnlyFromItsReader(t *testing.T) {
	fedData, fedErr := NewSnapshotPump(PumpSources{
		Run: func() *RunSection { return NewRunFeed().Section() },
	}).Marshal()
	if fedErr != nil {
		t.Fatal(fedErr)
	}
	if !strings.Contains(string(fedData), `"run":{`) {
		t.Errorf("a pump with a run reader sent no run key:\n%s", fedData)
	}
	if !strings.Contains(string(fedData), `"tools":[]`) {
		t.Errorf("an empty feed must render tools as [] (never null, never a vanished key):\n%s", fedData)
	}

	// The mutation tooth: pull the reader off the pump and the key must be
	// gone - the four-key byte nails' shape. A fill line that fired without
	// a reader, or a reader wired to nothing, shows up here.
	bareData, bareErr := NewSnapshotPump(PumpSources{}).Marshal()
	if bareErr != nil {
		t.Fatal(bareErr)
	}
	if strings.Contains(string(bareData), `"run"`) {
		t.Errorf("a pump without a run reader invented a run key:\n%s", bareData)
	}
}

// TestRunFeedCarriesOnlyWhatEventsWrote is AC#6's answer for this key: every
// field is asserted against the scripted event values verbatim, in order.
func TestRunFeedCarriesOnlyWhatEventsWrote(t *testing.T) {
	feed := NewRunFeed()
	// The values are what loop.go's EvToolEnd publish sites carry: the
	// journal's own outcome vocabulary, the call id, the name, the result
	// text.
	feed.ToolEnded("call-1", "fs.write", "success", "wrote 3 files")
	feed.ToolEnded("call-2", "net.fetch", "error", "boom")
	// A classified error event: class and detail are observe.Error's own
	// fields as the event carried them, text is the event's message.
	feed.Errored("provider", "http 502 from upstream", "调用失败")
	// The loop's own stuck sentence and the C22 rung.
	feed.Stuck("fs.write 已连续调用 3 次且参数相同，任务停在 Stuck 态。", 2)
	for _, delta := range []string{"先查", "目录结构", "再决定"} {
		feed.ReasoningDelta(delta)
	}

	sec := feed.Section()
	if sec == nil {
		t.Fatal("a live feed built a nil section - the key would vanish from the wire")
	}
	if len(sec.Tools) != 2 {
		t.Fatalf("tools = %#v, want the two ended rows", sec.Tools)
	}
	if sec.Tools[0].CallID != "call-1" || sec.Tools[0].Name != "fs.write" ||
		sec.Tools[0].Outcome != "success" || sec.Tools[0].ResultText != "wrote 3 files" {
		t.Errorf("row 0 = %#v, want the event's own values verbatim", sec.Tools[0])
	}
	if sec.Tools[1].Outcome != "error" || sec.Tools[1].ResultText != "boom" {
		t.Errorf("row 1 = %#v, want the event's own outcome and text", sec.Tools[1])
	}
	if sec.Error == nil || sec.Error.Class != "provider" ||
		sec.Error.Detail != "http 502 from upstream" || sec.Error.Text != "调用失败" {
		t.Errorf("error row = %#v, want class/detail/text verbatim", sec.Error)
	}
	if sec.Stuck == nil || sec.Stuck.RepeatLevel != 2 ||
		sec.Stuck.Text != "fs.write 已连续调用 3 次且参数相同，任务停在 Stuck 态。" {
		t.Errorf("stuck row = %#v, want the loop's sentence and rung verbatim", sec.Stuck)
	}
	if sec.Reasoning != "先查目录结构再决定" {
		t.Errorf("reasoning = %q, want the deltas in stream order", sec.Reasoning)
	}
	if sec.ToolsDropped != 0 {
		t.Errorf("toolsDropped = %d with nothing evicted", sec.ToolsDropped)
	}

	// The wire carries the same values, byte-level.
	data, err := json.Marshal(sec)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"outcome":"success"`, `"outcome":"error"`, `"class":"provider"`, `"repeatLevel":2`} {
		if !strings.Contains(string(data), want) {
			t.Errorf("the exit bytes lost %s:\n%s", want, data)
		}
	}

	// A notice-shaped error (no class, no detail, no text) writes nothing -
	// inventing a class for it would be the constant-as-data shape.
	quiet := NewRunFeed()
	quiet.Errored("", "", "")
	if s := quiet.Section(); s.Error != nil {
		t.Errorf("an eventless error wrote %#v - the row must stay nil until an event carries one", s.Error)
	}
}

// TestRunFeedBoundsAreHonest pins both bounds' honesty rule: eviction and
// truncation are counted, never silent.
func TestRunFeedBoundsAreHonest(t *testing.T) {
	feed := NewRunFeed()
	total := runToolRowsKept + 5
	for i := 0; i < total; i++ {
		feed.ToolEnded(fmt.Sprintf("call-%03d", i), "tool", "success", "")
	}
	sec := feed.Section()
	if len(sec.Tools) != runToolRowsKept {
		t.Fatalf("tools = %d rows, want the cap %d", len(sec.Tools), runToolRowsKept)
	}
	if sec.ToolsDropped != 5 {
		t.Errorf("toolsDropped = %d, want 5 - a bound without its count is silence", sec.ToolsDropped)
	}
	if sec.Tools[0].CallID != fmt.Sprintf("call-%03d", total-runToolRowsKept) {
		t.Errorf("the oldest kept row is %q, want call-%03d - eviction must drop the oldest, not the newest",
			sec.Tools[0].CallID, total-runToolRowsKept)
	}

	// Reasoning's tail window: past it, the elided count is inside the
	// marker and the retained text is the TAIL (the recent reasoning is what
	// the panel renders mid-stream).
	long := NewRunFeed()
	head := strings.Repeat("头", runReasoningKeepRunes)
	tail := strings.Repeat("尾", 100)
	long.ReasoningDelta(head + tail)
	rendered := long.Section().Reasoning
	wantElided := len([]rune(head+tail)) - runReasoningKeepRunes
	if !strings.Contains(rendered, fmt.Sprintf("truncated: %d runes elided", wantElided)) {
		t.Errorf("the marker's count is not the honest one (want %d):\n%q", wantElided, rendered)
	}
	if !strings.HasSuffix(rendered, tail) {
		t.Errorf("the window must keep the tail, got suffix ...%q", lastRunes(rendered, len(tail)))
	}
}

func lastRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[len(r)-n:])
}
