package main

// Ticket 145's two-hop wiring, the teeth:
//
//   hop 1 - the console sink's event switch feeds rt.runFeed (Publish's
//           EvToolEnd / EvReasoningDelta / EvError / EvStuck branches);
//   hop 2 - the assembly root's PumpSources.Run reader hands the feed's
//           Section to the pump.
//
// The carrier itself (RunFeed / RunSection and its honesty rules) is
// internal/panel's run_feed_145.go, whose own tests pin the reader-gated
// shape and the four-key byte nails. This file pins the two hops production
// code owns: the events a sink actually receives reach the wire, a sink
// stripped of its feed lines writes nothing (mutation), and a pump built
// without the reader sends no "run" key at all - the shape the four-key byte
// nails stay green under. The production packet carrying no "run" key until
// a reader exists is the landed design, not a defect.

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
)

// runSectionProbe145 decodes the wire bytes and isolates the "run" key, so a
// test asserts the key the panel would actually receive, not a Go field a
// composer could fill behind the wire's back.
func runSectionProbe145(t *testing.T, data []byte) (present bool, sec *panel.RunSection) {
	t.Helper()
	var envelope struct {
		Run *panel.RunSection `json:"run"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		t.Fatalf("the wire bytes are not the snapshot's JSON: %v", err)
	}
	return envelope.Run != nil, envelope.Run
}

// sinkEvents145 publishes one event of each fed kind through the PRODUCTION
// consoleSink - the same Publish switch execute() hands the loop - so the
// teeth exercise the wiring as shipped, not a test double of it.
func sinkEvents145(sink consoleSink) {
	sink.Publish(agent.Event{Kind: agent.EvToolEnd, CallID: "call-1", ToolName: "fs.read", Outcome: "success", Text: "42 bytes"})
	sink.Publish(agent.Event{Kind: agent.EvReasoningDelta, Text: "thinking out loud"})
	sink.Publish(agent.Event{Kind: agent.EvError, Text: "the provider refused", Err: &observe.Error{Class: observe.ClassProvider, Detail: "http 429"}})
	sink.Publish(agent.Event{Kind: agent.EvStuck, Text: "the loop is repeating itself", RepeatLevel: 2})
}

func Test145SinkFeedLinesCarryTheEventsIntoTheRunKey(t *testing.T) {
	feed := panel.NewRunFeed()
	sink := consoleSink{out: io.Discard, runFeed: feed}
	sinkEvents145(sink)
	// The notice-shaped EvError (Text with no classified error) comes last on
	// purpose: run_feed_145.go's RunErrorRow doc says a notice writes no row,
	// so a lastError overwritten by the notice's text here would mean the
	// sink invented a row the event cannot back.
	sink.Publish(agent.Event{Kind: agent.EvError, Text: "budget notice, not a failure"})

	data, err := panel.NewSnapshotPump(panel.PumpSources{Run: feed.Section}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	present, run := runSectionProbe145(t, data)
	if !present || run == nil {
		t.Fatalf("the fed sink's packet carries no run key at all:\n%s", data)
	}
	if len(run.Tools) != 1 {
		t.Fatalf("tools = %+v, want exactly the one ended call", run.Tools)
	}
	tool := run.Tools[0]
	if tool.CallID != "call-1" || tool.Name != "fs.read" || tool.Outcome != "success" || tool.ResultText != "42 bytes" {
		t.Errorf("tool row lost the event's own values: %+v", tool)
	}
	if run.Reasoning != "thinking out loud" {
		t.Errorf("reasoning = %q, want the delta text verbatim", run.Reasoning)
	}
	if run.Error == nil {
		t.Fatal("the classified error wrote no row")
	}
	if run.Error.Class != string(observe.ClassProvider) || run.Error.Detail != "http 429" || run.Error.Text != "the provider refused" {
		t.Errorf("error row lost the event's own values: %+v", run.Error)
	}
	// The notice must not have replaced the classified row (and must not have
	// manufactured a class out of nothing).
	if run.Error.Text == "budget notice, not a failure" {
		t.Errorf("a notice-shaped EvError overwrote the classified error row: %+v", run.Error)
	}
	if run.Stuck == nil || run.Stuck.Text != "the loop is repeating itself" || run.Stuck.RepeatLevel != 2 {
		t.Errorf("stuck row lost the event's own values: %+v", run.Stuck)
	}
}

func Test145MutationSinkWithoutItsFeedLinesWritesNothing(t *testing.T) {
	// Hop-1 mutation: the same events through a sink whose feed line is gone
	// (runFeed nil - the shape a sink built without the wiring has). The
	// reader hop below stays INTACT on purpose, so the packet still carries
	// the key - but only as an empty feed, because the sink lines are the
	// events' only route into the holder. An empty feed that suddenly showed
	// values here would mean a second writer exists.
	feed := panel.NewRunFeed()
	sink := consoleSink{out: io.Discard, runFeed: nil}
	sinkEvents145(sink)

	data, err := panel.NewSnapshotPump(panel.PumpSources{Run: feed.Section}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	present, run := runSectionProbe145(t, data)
	if !present || run == nil {
		t.Fatalf("the reader exists, so the key must exist (empty feed = tools:[]):\n%s", data)
	}
	if len(run.Tools) != 0 || run.Reasoning != "" || run.Error != nil || run.Stuck != nil {
		t.Errorf("a feed-less sink still wrote rows - a second writer is hiding: %+v", run)
	}
	if !strings.Contains(string(data), `"tools":[]`) {
		t.Errorf("an empty feed must render tools:[] not null or absent:\n%s", data)
	}
}

func Test145MutationReaderRemovedSendsNoRunKeyAtAll(t *testing.T) {
	// Hop-2 mutation: the sink keeps feeding, but the assembly root's reader
	// is gone - the pump the four-key byte nails pin. The key must be absent
	// from the WIRE BYTES, not merely nil in Go: composer omitempty is the
	// last mile, and this asserts the bytes the panel receives.
	feed := panel.NewRunFeed()
	sink := consoleSink{out: io.Discard, runFeed: feed}
	sinkEvents145(sink)

	data, err := panel.NewSnapshotPump(panel.PumpSources{}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if bytes.Contains(data, []byte(`"run"`)) {
		t.Errorf("a pump without the Run reader still sent the key:\n%s", data)
	}
}
