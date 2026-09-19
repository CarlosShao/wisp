// Command mockllm is the Wisp mock LLM server (ticket 09): an
// OpenAI/Anthropic-compatible mock with golden replay and fault injection.
// It is a separate Go module (stdlib only) so it can be built and run
// without the sherpa/sqlite cgo budget: `go run .` from this directory, or
// via compose (port 18080).
//
// Endpoints:
//
//	POST /v1/chat/completions   OpenAI chat, stream SSE or JSON
//	POST /v1/responses          OpenAI Responses (framework level)
//	POST /v1/messages           Anthropic Messages (framework level)
//	GET  /v1/models             model list (discovery target)
//	POST /__control/fail_next   inject N failing responses {status,times,retry_after,body}
//	POST /__control/latency     per-chunk delay ms {ms}
//	POST /__control/truncate    cut the stream after N chunks {chunks}
//	POST /__control/reset       clear all injections and counters
//	GET  /__control/state       dump state + request counters
//
// Golden mode: a request whose model is "golden/<name>" (or that carries the
// X-Wisp-Golden header <name>) replays testdata/golden/<name>.sse through
// the SAME format grammar as internal/llm/golden (one format, two runners;
// the parser here mirrors that package and the equivalence is pinned by the
// integration test).
//
// Deterministic synthesis (no golden): the mock echoes the last user text
// ("echo: ..."), replies "vision-ok" when an image_url part is present,
// "audio-ok" for an input_audio part, emits a tool call when tools are
// declared with tool_choice required/forced, truncates with
// finish_reason=length when max_tokens is smaller than the answer, and
// reports chunked usage. Anything else is stable by construction, which is
// what tests and probe runs need.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:18080", "listen address (use :0 for a test port)")
	goldenDir := flag.String("golden-dir", "testdata/golden", "directory of golden .sse files")
	models := flag.String("models", "mock-small,mock-large", "comma-separated models reported by /v1/models")
	printAddr := flag.Bool("print-addr", false, "print MOCKLLM_ADDR=<addr> once listening")
	flag.Parse()

	srv := NewServer(*goldenDir, splitCSV(*models))
	mux := http.NewServeMux()
	srv.Register(mux)

	ln, err := listen(*addr)
	if err != nil {
		log.Fatalf("mockllm: listen %s: %v", *addr, err)
	}
	if *printAddr {
		fmt.Fprintf(os.Stdout, "MOCKLLM_ADDR=%s\n", ln.Addr().String())
		os.Stdout.Sync()
	}
	log.Printf("mockllm: serving on http://%s (golden dir %s)", ln.Addr().String(), *goldenDir)
	if err := (&http.Server{Handler: mux}).Serve(ln); err != nil && err != http.ErrServerClosed {
		log.Fatalf("mockllm: serve: %v", err)
	}
}
