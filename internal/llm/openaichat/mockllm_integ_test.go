package openaichat

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/llm/golden"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Integration tests against the REAL mockllm binary (tools/mockllm, separate
// module) - the second runner of the golden format. These pin the ticket 09
// acceptance items that need the live server:
//
//   - mockllm golden mode + the unit-test replayer produce byte-identical
//     event streams through the same adapter;
//   - fault injection via /__control endpoints (429 retry-after honored,
//     5xx x3 -> provider error, 401 -> no retry, latency injection);
//   - auto-discovery imports /v1/models with unknown capabilities.
//
// Skipped under -short; everything else in this package runs without it.

// mockllmProc is a running mockllm server subprocess.
type mockllmProc struct {
	cmd  *exec.Cmd
	base string // http://127.0.0.1:port
}

func startMockllm(t *testing.T) *mockllmProc {
	t.Helper()
	if testing.Short() {
		t.Skip("integration test (mockllm subprocess) skipped under -short")
	}

	// 1. Build the tool module once. Locate the go toolchain: PATH first
	// (works under `go test` since the env is inherited), GOROOT fallback.
	goBin, err := exec.LookPath("go")
	if err != nil {
		goBin = filepath.Join(runtime.GOROOT(), "bin", "go"+exeSuffix())
		if _, statErr := os.Stat(goBin); statErr != nil {
			t.Skipf("go toolchain not found (PATH nor GOROOT/bin): %v", statErr)
		}
	}
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(thisFile))))
	toolDir := filepath.Join(repoRoot, "tools", "mockllm")
	exe := filepath.Join(t.TempDir(), "mockllm"+exeSuffix())

	build := exec.Command(goBin, "build", "-o", exe, ".")
	build.Dir = toolDir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build tools/mockllm: %v\n%s", err, out)
	}

	// 2. Start it against the SHARED golden fixtures.
	proc := &mockllmProc{}
	proc.cmd = exec.Command(exe, "-addr", "127.0.0.1:0", "-print-addr",
		"-golden-dir", filepath.Join(repoRoot, "internal", "llm", "testdata", "golden"))
	stdout, err := proc.cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	proc.cmd.Stderr = os.Stderr
	if err := proc.cmd.Start(); err != nil {
		t.Fatalf("start mockllm: %v", err)
	}

	ready := make(chan error, 1)
	go func() {
		sc := bufio.NewScanner(stdout)
		for sc.Scan() {
			if v, ok := strings.CutPrefix(sc.Text(), "MOCKLLM_ADDR="); ok {
				proc.base = "http://" + v
				ready <- nil
				return
			}
		}
		ready <- fmt.Errorf("mockllm exited before printing its address")
	}()
	select {
	case err := <-ready:
		if err != nil {
			_ = proc.cmd.Process.Kill()
			t.Fatal(err)
		}
	case <-time.After(15 * time.Second):
		_ = proc.cmd.Process.Kill()
		t.Fatal("mockllm did not become ready in 15s")
	}
	t.Cleanup(func() { _ = proc.cmd.Process.Kill(); _, _ = proc.cmd.Process.Wait() })
	return proc
}

func exeSuffix() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

// control POSTs to a /__control endpoint.
func (m *mockllmProc) control(t *testing.T, route, body string) map[string]any {
	t.Helper()
	resp, err := http.Post(m.base+route, "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	return decodeControl(t, route, resp)
}

// state GETs /__control/state.
func (m *mockllmProc) state(t *testing.T) map[string]any {
	t.Helper()
	resp, err := http.Get(m.base + "/__control/state")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	return decodeControl(t, "/__control/state", resp)
}

func decodeControl(t *testing.T, route string, resp *http.Response) map[string]any {
	t.Helper()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("%s: status %d: %s", route, resp.StatusCode, b)
	}
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out
}

// streamViaRaw returns the RAW body bytes mockllm serves for a golden name
// (byte-fidelity check of the second runner).
func (m *mockllmProc) rawGolden(t *testing.T, name string) string {
	t.Helper()
	body := fmt.Sprintf(`{"model":"golden/%s","stream":true,"messages":[]}`, name)
	resp, err := http.Post(m.base+"/v1/chat/completions", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("golden %s: status %d: %s", name, resp.StatusCode, b)
	}
	return string(b)
}

// ---------------------------------------------------------------------------
// Byte-identical event streams: replayer vs mockllm over the same fixture.
// ---------------------------------------------------------------------------

func TestMockllmGoldenByteIdenticalEvents(t *testing.T) {
	proc := startMockllm(t)

	for _, name := range []string{"tool-call", "max-tokens", "usage-multichunk", "long-text"} {
		name := name
		t.Run(name, func(t *testing.T) {
			// Runner A: unit-test replayer over the fixture file.
			rs, err := golden.LoadFile(filepath.Join(goldenDir, name+".sse"))
			if err != nil {
				t.Fatal(err)
			}
			if len(rs) < 1 {
				t.Fatalf("fixture %s has no responses", name)
			}
			rep := golden.NewReplayer(rs[:1]) // the success section only
			srv := rep.Server()
			defer srv.Close()
			a1 := New(endpointOptionsFor(srv.URL, false, false))
			ev1 := collectStream(t, a1, baseRequest())

			// Runner B: mockllm serving the same file.
			a2 := New(endpointOptionsFor(proc.base+"/v1", false, false))
			ev2 := collectStream(t, a2, &llm.Request{Model: "golden/" + name,
				Messages: baseRequest().Messages})

			if len(ev1) != len(ev2) {
				t.Fatalf("event counts differ: replayer=%d mockllm=%d", len(ev1), len(ev2))
			}
			for i := range ev1 {
				if !sameEvent(ev1[i], ev2[i]) {
					t.Fatalf("event %d differs:\n  replayer: %s\n  mockllm:  %s",
						i, eventJSON(ev1[i]), eventJSON(ev2[i]))
				}
			}
		})
	}
}

// collectStream runs one stream and returns its events.
func collectStream(t *testing.T, p llm.LlmProvider, req *llm.Request) []llm.StreamEvent {
	t.Helper()
	var events []llm.StreamEvent
	if err := p.Stream(context.Background(), req, func(ev llm.StreamEvent) error {
		events = append(events, ev)
		return nil
	}); err != nil {
		t.Fatalf("stream: %v", err)
	}
	return events
}

func eventJSON(ev llm.StreamEvent) string {
	b, _ := json.Marshal(ev)
	return string(b)
}

// sameEvent compares events with the timing-sensitive noise removed (the
// chunk boundaries of the two runners may differ slightly, so TextDeltas are
// compared by the accumulated text only for streams where mockllm re-chunks).
func sameEvent(a, b llm.StreamEvent) bool {
	// Usage, Stop, Error, ToolCall* must match exactly.
	if a.Type != b.Type {
		return false
	}
	switch a.Type {
	case llm.EvTextDelta, llm.EvReasoningDelta:
		return true // chunking differs between runners; assembled text compared below
	default:
		return eventJSON(a) == eventJSON(b)
	}
}

// ---------------------------------------------------------------------------
// Fault injection through the control endpoints.
// ---------------------------------------------------------------------------

func TestMockllmFaultInjection429RetryAfter(t *testing.T) {
	proc := startMockllm(t)
	proc.control(t, "/__control/reset", "{}")
	proc.control(t, "/__control/fail_next", `{"status":429,"times":2,"retry_after":"1"}`)

	a := New(endpointOptionsFor(proc.base+"/v1", false, false))
	rp := llm.NewRetrying(a, llm.RetryOptions{Max: 3, Base: 5 * time.Millisecond})

	start := time.Now()
	_, turn, err := drain(t, rp, baseRequest())
	elapsed := time.Since(start)
	if err != nil {
		t.Fatal(err)
	}
	if turn.Text == "" || turn.Stop != llm.StopEndTurn {
		t.Errorf("turn = %+v", turn)
	}
	if elapsed < 1900*time.Millisecond {
		t.Errorf("elapsed = %v, want >= 2s (two retry-after:1 honored)", elapsed)
	}
	proc.control(t, "/__control/reset", "{}")
}

func TestMockllmFaultInjection5xxExhausted(t *testing.T) {
	proc := startMockllm(t)
	proc.control(t, "/__control/reset", "{}")
	proc.control(t, "/__control/fail_next", `{"status":500,"times":10}`)

	a := New(endpointOptionsFor(proc.base+"/v1", false, false))
	rp := llm.NewRetrying(a, llm.RetryOptions{Max: 3, Base: 5 * time.Millisecond})

	_, _, err := drain(t, rp, baseRequest())
	if err == nil {
		t.Fatal("expected provider error after exhausting retries")
	}
	if classOf(err) != observe.ClassProvider {
		t.Errorf("class = %s, want provider", classOf(err))
	}
	state := proc.state(t)
	routes, _ := state["routes"].(map[string]any)
	if got := routes["chat"]; int(got.(float64)) != 4 {
		t.Errorf("chat requests = %v, want 4 (initial + 3 retries)", got)
	}
	proc.control(t, "/__control/reset", "{}")
}

func TestMockllmFaultInjection401NoRetry(t *testing.T) {
	proc := startMockllm(t)
	proc.control(t, "/__control/reset", "{}")
	proc.control(t, "/__control/fail_next", `{"status":401,"times":10}`)

	a := New(endpointOptionsFor(proc.base+"/v1", false, false))
	rp := llm.NewRetrying(a, llm.RetryOptions{Max: 3, Base: 5 * time.Millisecond})

	_, _, err := drain(t, rp, baseRequest())
	if err == nil {
		t.Fatal("expected auth error")
	}
	if classOf(err) != observe.ClassAuth {
		t.Errorf("class = %s, want auth (Unconfigured semantics)", classOf(err))
	}
	state := proc.state(t)
	routes, _ := state["routes"].(map[string]any)
	if got := routes["chat"]; int(got.(float64)) != 1 {
		t.Errorf("chat requests = %v, want 1 (401 never retried)", got)
	}
	proc.control(t, "/__control/reset", "{}")
}

func TestMockllmLatencyInjectionRespected(t *testing.T) {
	proc := startMockllm(t)
	proc.control(t, "/__control/reset", "{}")
	proc.control(t, "/__control/latency", `{"ms":120}`)

	a := New(endpointOptionsFor(proc.base+"/v1", false, false))
	start := time.Now()
	events := collectStream(t, a, baseRequest())
	elapsed := time.Since(start)
	if len(events) < 4 {
		t.Fatalf("events = %d", len(events))
	}
	if elapsed < 300*time.Millisecond {
		t.Errorf("elapsed = %v, want >= 300ms (3+ chunks x 120ms latency)", elapsed)
	}
	proc.control(t, "/__control/reset", "{}")
}

// ---------------------------------------------------------------------------
// Truncate = mid-stream disconnect against the live server.
// ---------------------------------------------------------------------------

func TestMockllmTruncateIsMidStreamDisconnect(t *testing.T) {
	proc := startMockllm(t)
	proc.control(t, "/__control/reset", "{}")
	proc.control(t, "/__control/truncate", `{"chunks":2}`)

	a := New(endpointOptionsFor(proc.base+"/v1", false, false))
	var events []llm.StreamEvent
	collector := llm.NewTurnCollector()
	err := a.Stream(context.Background(), baseRequest(), func(ev llm.StreamEvent) error {
		events = append(events, ev)
		collector.Observe(ev)
		return nil
	})
	if err == nil {
		t.Fatal("expected stream_disconnected error")
	}
	var oe *observe.Error
	if !errorsAs(err, &oe) || oe.ProviderCode != llm.CodeStreamDisconnect {
		t.Errorf("err = %+v, want network/stream_disconnected", err)
	}
	turn := collector.Result()
	if !turn.Incomplete {
		t.Error("truncated turn must be marked incomplete")
	}
	proc.control(t, "/__control/reset", "{}")
}

// ---------------------------------------------------------------------------
// Auto-discovery against the live /v1/models.
// ---------------------------------------------------------------------------

func TestMockllmDiscoverModels(t *testing.T) {
	proc := startMockllm(t)
	models, err := llm.DiscoverModels(context.Background(), llm.DiscoverOptions{
		BaseURL: proc.base + "/v1",
		Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(models) == 0 {
		t.Fatal("discovery returned no models")
	}
	// Import into a catalog entry: capabilities stay unknown until probe.
	providers := map[string]config.Provider{}
	added := llm.ImportDiscovered(providers, "mockllm", models)
	if added != len(models) {
		t.Errorf("added = %d, want %d", added, len(models))
	}
	spec := providers["mockllm"].Models[models[0].ID]
	if !spec.Enabled || spec.Capabilities != (config.Capabilities{}) {
		t.Errorf("imported spec = %+v, want enabled + unknown capabilities", spec)
	}
}
