package main

// Ticket 145's two-hop wiring, the teeth:
//
//   hop 1 - the console sink's event switch feeds rt.runFeed (Publish's
//           EvToolStart / EvToolEnd / EvUsage / EvReasoningDelta / EvError /
//           EvStuck / EvDone branches);
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
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/memory"
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

// Test145RunningUsageAndDoneTravelTheWire is AC#6's answer for the three
// dimensions whose hop-1 lines this wave landed: the running chip exists
// because the loop now forwards the adapters' tool_call_start (the kind's
// first producer), the usage row exists because the accumulator's totals
// now book an event, and the running chips retire with the task - the
// wire carries what the events wrote and nothing else.
func Test145RunningUsageAndDoneTravelTheWire(t *testing.T) {
	feed := panel.NewRunFeed()
	sink := consoleSink{out: io.Discard, runFeed: feed}

	sink.Publish(agent.Event{Kind: agent.EvToolStart, CallID: "call-r", ToolName: "fs.read"})
	sink.Publish(agent.Event{
		Kind:       agent.EvUsage,
		Usage:      llm.Usage{InputTokens: 120, OutputTokens: 24, CachedTokens: 96},
		CostMicros: 3500, Currency: "CNY",
	})

	data, err := panel.NewSnapshotPump(panel.PumpSources{Run: feed.Section}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	present, run := runSectionProbe145(t, data)
	if !present || run == nil {
		t.Fatalf("the fed sink's packet carries no run key at all:\n%s", data)
	}
	if len(run.Running) != 1 || run.Running[0].CallID != "call-r" || run.Running[0].Name != "fs.read" {
		t.Errorf("running = %+v, want the start event's own call in flight", run.Running)
	}
	if run.Usage == nil || run.Usage.TokensIn != 120 || run.Usage.TokensOut != 24 ||
		run.Usage.CachedTokens != 96 || run.Usage.Micros != 3500 || run.Usage.Currency != "CNY" {
		t.Errorf("usage = %+v, want the accumulator's own totals verbatim", run.Usage)
	}
	for _, want := range []string{`"micros":3500`, `"currency":"CNY"`} {
		if !strings.Contains(string(data), want) {
			t.Errorf("the exit bytes lost the micros/currency pair member %s:\n%s", want, data)
		}
	}

	// The end retires its chip into an ended row; the task's end retires any
	// that never got one. The terminal event now carries the loop's own stop
	// word (ticket 145 row 13): a cancelled ending's cancelled, verbatim.
	sink.Publish(agent.Event{Kind: agent.EvToolEnd, CallID: "call-r", ToolName: "fs.read", Outcome: "success", Text: "42 bytes"})
	sink.Publish(agent.Event{Kind: agent.EvToolStart, CallID: "call-lost", ToolName: "net.fetch"})
	sink.Publish(agent.Event{Kind: agent.EvDone, TaskID: "task-1", Text: "done", Stop: llm.StopCancelled})

	data, err = panel.NewSnapshotPump(panel.PumpSources{Run: feed.Section}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	_, run = runSectionProbe145(t, data)
	if len(run.Running) != 0 {
		t.Errorf("running after done = %+v, want every chip retired with the task", run.Running)
	}
	if len(run.Tools) != 1 || run.Tools[0].CallID != "call-r" || run.Tools[0].Outcome != "success" {
		t.Errorf("tools = %+v, want the ended call kept as the record", run.Tools)
	}
	if run.Done == nil || run.Done.Stop != string(llm.StopCancelled) {
		t.Errorf("done row = %+v, want the terminal event's own stop word verbatim", run.Done)
	}
	if !strings.Contains(string(data), `"stop":"cancelled"`) {
		t.Errorf("the exit bytes lost the stop word:\n%s", data)
	}
}

// Test145ArgsJSONComesFromTheBookedForensicsRow is AC#6's answer for the
// expanded-args dimension: the argsJson on the wire is the booked forensics
// row's sanitized args, joined by the correlation id the event carries -
// not a digest the sink invented. A call with no booked row (the mutation
// arms below) travels with argsJson absent.
func Test145ArgsJSONComesFromTheBookedForensicsRow(t *testing.T) {
	dir := t.TempDir()
	st, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open(%s): %v", dir, err)
	}
	t.Cleanup(func() { _ = st.Close() })
	if _, err := st.InsertToolCall(context.Background(), memory.ToolCall{
		TaskID: "task-145", Tool: "fs.write",
		ArgsJSON:      `{"path":"a.txt","content":"x"}`,
		RiskLevel:     "L1",
		CorrelationID: "task-145#call-1",
	}); err != nil {
		t.Fatalf("InsertToolCall: %v", err)
	}

	feed := panel.NewRunFeed()
	sink := consoleSink{out: io.Discard, runFeed: feed, store: st}
	sink.Publish(agent.Event{
		Kind: agent.EvToolEnd, TaskID: "task-145", CallID: "call-1",
		ToolName: "fs.write", Outcome: "success", Text: "wrote",
		CorrelationID: "task-145#call-1",
	})

	data, err := panel.NewSnapshotPump(panel.PumpSources{Run: feed.Section}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(data), `"argsJson":"{\"path\":\"a.txt\",\"content\":\"x\"}"`) {
		t.Errorf("the wire lost the booked row's sanitized args:\n%s", data)
	}
	// Ticket 145 row 4's risk dimension rides the SAME booked row: the gate's
	// L1 decision travels verbatim as the row's level.
	if !strings.Contains(string(data), `"level":"L1"`) {
		t.Errorf("the wire lost the booked row's RiskLevel:\n%s", data)
	}

	// Mutation 1: no store (a sink built without the forensics face) - the
	// args must be absent, never a digest invented from the event alone.
	bare := panel.NewRunFeed()
	consoleSink{out: io.Discard, runFeed: bare}.Publish(agent.Event{
		Kind: agent.EvToolEnd, TaskID: "task-145", CallID: "call-1",
		ToolName: "fs.write", Outcome: "success", Text: "wrote",
		CorrelationID: "task-145#call-1",
	})
	bareData, err := panel.NewSnapshotPump(panel.PumpSources{Run: bare.Section}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(bareData), "argsJson") {
		t.Errorf("a store-less sink invented args from the event alone:\n%s", bareData)
	}
	if strings.Contains(string(bareData), `"level"`) {
		t.Errorf("a store-less sink invented a risk level from the event alone:\n%s", bareData)
	}

	// Mutation 2: a correlation id no booked row answers (the loop-rejected
	// and locally-answered calls) - absent is the truth here too.
	noRow := panel.NewRunFeed()
	sink.Publish(agent.Event{
		Kind: agent.EvToolEnd, TaskID: "task-145", CallID: "call-9",
		ToolName: "list_tools", Outcome: "success", Text: "answered in-loop",
		CorrelationID: "task-145#call-9",
	})
	noRowData, err := panel.NewSnapshotPump(panel.PumpSources{Run: noRow.Section}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(noRowData), "argsJson") {
		t.Errorf("an unbooked call grew args from some other row:\n%s", noRowData)
	}
	if strings.Contains(string(noRowData), `"level"`) {
		t.Errorf("an unbooked call grew a risk level from some other row:\n%s", noRowData)
	}
}

// Test145RetryableRidesTheClassPolicy is AC#6's answer for the retryable
// half-key: the wire value is observe.Error.Retryable()'s own answer for the
// class the event carried - provider retries, auth does not - and the last
// classified error's row is the one on the wire.
func Test145RetryableRidesTheClassPolicy(t *testing.T) {
	feed := panel.NewRunFeed()
	sink := consoleSink{out: io.Discard, runFeed: feed}
	sink.Publish(agent.Event{Kind: agent.EvError, Text: "502", Err: &observe.Error{Class: observe.ClassProvider, Detail: "http 502"}})
	sink.Publish(agent.Event{Kind: agent.EvError, Text: "key", Err: &observe.Error{Class: observe.ClassAuth, Detail: "http 401"}})

	data, err := panel.NewSnapshotPump(panel.PumpSources{Run: feed.Section}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	_, run := runSectionProbe145(t, data)
	if run == nil || run.Error == nil {
		t.Fatalf("the classified errors wrote no row:\n%s", data)
	}
	if run.Error.Class != string(observe.ClassAuth) {
		t.Errorf("error class = %q, want the LAST classified error's own class", run.Error.Class)
	}
	if run.Error.Retryable {
		t.Errorf("retryable = true for %s, want the class policy's own no - carrying the provider arm's answer would be a fabricated retry promise", run.Error.Class)
	}
	if !strings.Contains(string(data), `"retryable":false`) {
		t.Errorf("the exit bytes lost the class policy's no:\n%s", data)
	}
}
