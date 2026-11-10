package main

// Ticket 414 (R-145-1): the ORIGIN hop of run.done.stop, the teeth.
//
// The carrier hop is already pinned: panel_run_145_test.go's
// Test145RunningUsageAndDoneTravelTheWire publishes a HAND-BUILT EvDone with
// Stop through the production sink and nails "stop":"cancelled" on the wire.
// What no test in the repo asserted was where that word comes from — the
// 145-accept-r1 census's acc-r1-mut probe cut internal/agent/loop.go's
// `Stop: stop` (the finish() publish, loop.go:1134) and the WHOLE ./cmd/wisp
// suite stayed green. This file closes that gap.
//
// The shape: a REAL agent.Loop over the C5 golden seam (a scripted SSE whose
// finish_reason "stop" the genuine openai-chat adapter maps to
// llm.StopEndTurn), the loop's events through the PRODUCTION consoleSink into
// a RunFeed, and the wire bytes the panel would receive. The asserted word is
// the replayed provider stream's OWN stop word, verbatim:
//
//	golden SSE finish_reason=stop -> adapter EvStop end_turn -> loop.go
//	finish() EvDone.Stop -> run.go's TaskEnded(string(e.Stop)) -> feed ->
//	"done":{"stop":"end_turn"} on the wire.
//
// That is what makes the assertion an origin-hop tooth rather than the
// decorative "value non-nil" shape: a loop that stops publishing the
// provider's word (mutation: loop.go:1134 `Stop: stop` -> `Stop: ""`) empties
// the terminal event, and nothing between the loop and the wire invents a
// replacement, so this test turns red — a value hardcoded anywhere else on
// the path cannot reproduce the provider's own word either.

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/golden"
	_ "github.com/CarlosShao/wisp/internal/llm/openaichat" // protocol registration
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
)

// t414StopWire is one completed text turn: the provider says its own
// finish_reason=stop (mapped to llm.StopEndTurn) with two text deltas and a
// usage block, the same shape internal/agent's text-reply golden unit replays.
const t414StopWire = `# wisp golden sse v1
# @scenario ticket 414: the provider's own finish_reason is the only source of run.done.stop
# @response 200
data: {"choices":[{"index":0,"delta":{"role":"assistant","content":""},"finish_reason":null}]}

data: {"choices":[{"index":0,"delta":{"content":"好的，"},"finish_reason":null}]}

data: {"choices":[{"index":0,"delta":{"content":"已经把提醒设好了。"},"finish_reason":null}]}

data: {"choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}

data: {"choices":[],"usage":{"prompt_tokens":140,"completion_tokens":9,"total_tokens":149,"prompt_tokens_details":{"cached_tokens":96}}}

data: [DONE]
`

// TestRunDoneStopComesFromTheLoopsTerminalEvent is ticket 414's answer to
// R-145-1: run.done.stop is the loop's terminal EvDone.Stop — the provider's
// own stop word — and the origin hop loop.go finish() -> Event -> consoleSink
// -> RunFeed is a real, asserted hop.
func TestRunDoneStopComesFromTheLoopsTerminalEvent(t *testing.T) {
	srv, _, err := golden.ServerFromBytes([]byte(t414StopWire))
	if err != nil {
		t.Fatalf("parse golden: %v", err)
	}
	t.Cleanup(srv.Close)

	prov, err := llm.NewProvider(llm.EndpointOptions{
		Provider: "mock", Model: "mock-small", Protocol: "openai-chat",
		BaseURL:       srv.URL,
		ContextWindow: 128000,
		HTTPClient:    &http.Client{Transport: &http.Transport{Proxy: nil}},
	})
	if err != nil {
		t.Fatalf("build provider: %v", err)
	}

	feed := panel.NewRunFeed()
	loop, err := agent.New(agent.Options{
		Provider: prov,
		Tools:    agent.NewEchoProvider(),
		Sink:     consoleSink{out: io.Discard, runFeed: feed},
		Registry: observe.NewRegistry(),
		Logger:   slog.New(slog.NewTextHandler(io.Discard, nil)),
		Config: agent.Config{
			Model: "mock-small", ContextWindow: 128000,
			ArtifactsDir:                t.TempDir(),
			PassThroughUnclassifiedRisk: true,
			SteeringEnabled:             true,
			Currency:                    "CNY",
			PerToolTimeout:              2 * time.Second,
		},
	})
	if err != nil {
		t.Fatalf("agent.New: %v", err)
	}

	// The real loop runs the replayed turn to completion and publishes every
	// event through the production sink; the sink's EvDone branch is the
	// origin hop's last mile (TaskEnded(string(e.Stop))).
	res := loop.Run(context.Background(), "提醒我五分钟后关火")
	if res.Status != agent.StatusCompleted {
		t.Fatalf("status = %s (%s), want completed", res.Status, res.Message)
	}
	if res.Text != "好的，已经把提醒设好了。" {
		t.Fatalf("text = %q, want the replayed reply (the loop ran a different turn)", res.Text)
	}

	// The origin hop's reading: the done row's stop word is the provider's
	// own word, verbatim. Not "non-nil", not a loop-side status copy — the
	// exact word the replayed wire spoke (finish_reason=stop -> end_turn).
	done := feed.Section().Done
	if done == nil {
		t.Fatal("no done row: the terminal EvDone never reached the feed")
	}
	if done.Stop != string(llm.StopEndTurn) {
		t.Fatalf(`run.done.stop = %q, want %q (the replayed wire's own finish_reason, spoken by loop.go's terminal EvDone)`,
			done.Stop, string(llm.StopEndTurn))
	}

	// And the same reading on the WIRE BYTES the panel receives, through the
	// pump reader exactly as the assembly root assembles it.
	data, err := panel.NewSnapshotPump(panel.PumpSources{Run: feed.Section}).Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !bytes.Contains(data, []byte(`"done":{"stop":"end_turn"}`)) {
		t.Fatalf("the wire lost the loop's terminal stop word:\n%s", data)
	}
}
