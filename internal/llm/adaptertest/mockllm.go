package adaptertest

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

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// The fault suite over the REAL mockllm subprocess (the second runner of the
// golden format). Ticket 09 built the endpoints; ticket 11 runs all three
// adapters through them so fault behavior is proven against a live HTTP
// server, not only against the in-process replayer.

// Mockllm is a running mockllm server subprocess.
type Mockllm struct {
	cmd  *exec.Cmd
	Base string // http://127.0.0.1:port
	t    *testing.T
}

// StartMockllm builds (once per call) and launches mockllm against the shared
// golden fixture directory. Skipped under -short.
func StartMockllm(t *testing.T) *Mockllm {
	t.Helper()
	if testing.Short() {
		t.Skip("integration test (mockllm subprocess) skipped under -short")
	}
	goBin, err := exec.LookPath("go")
	if err != nil {
		goBin = filepath.Join(runtime.GOROOT(), "bin", "go"+exeSuffix())
		if _, statErr := os.Stat(goBin); statErr != nil {
			t.Skipf("go toolchain not found (PATH nor GOROOT/bin): %v", statErr)
		}
	}
	root := RepoRoot()
	exe := filepath.Join(t.TempDir(), "mockllm"+exeSuffix())
	build := exec.Command(goBin, "build", "-o", exe, ".")
	build.Dir = filepath.Join(root, "tools", "mockllm")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build tools/mockllm: %v\n%s", err, out)
	}

	proc := &Mockllm{t: t}
	proc.cmd = exec.Command(exe, "-addr", "127.0.0.1:0", "-print-addr",
		"-golden-dir", GoldenDir())
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
				proc.Base = "http://" + v
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
	case <-time.After(20 * time.Second):
		_ = proc.cmd.Process.Kill()
		t.Fatal("mockllm did not become ready in 20s")
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

// Control POSTs to a /__control endpoint.
func (m *Mockllm) Control(t *testing.T, route, body string) map[string]any {
	t.Helper()
	resp, err := http.Post(m.Base+route, "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	return m.decode(route, resp)
}

// RouteCount reads /__control/state and returns the request count for one
// route name (chat | messages | responses | models).
func (m *Mockllm) RouteCount(t *testing.T, route string) int {
	t.Helper()
	resp, err := http.Get(m.Base + "/__control/state")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	state := m.decode("/__control/state", resp)
	routes, _ := state["routes"].(map[string]any)
	n, _ := routes[route].(float64)
	return int(n)
}

func (m *Mockllm) decode(route string, resp *http.Response) map[string]any {
	m.t.Helper()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		m.t.Fatalf("%s: status %d: %s", route, resp.StatusCode, b)
	}
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		m.t.Fatal(err)
	}
	return out
}

// Fault is a queued injected failure.
type Fault struct {
	Status     int
	Times      int
	RetryAfter string
	Body       string
}

// QueueFault installs fail_next.
func (m *Mockllm) QueueFault(t *testing.T, f Fault) {
	t.Helper()
	body, err := json.Marshal(map[string]any{
		"status": f.Status, "times": f.Times,
		"retry_after": f.RetryAfter, "body": f.Body,
	})
	if err != nil {
		t.Fatal(err)
	}
	m.Control(t, "/__control/fail_next", string(body))
}

// GoldenRequestModel is the mockllm golden selector for a fixture name.
func GoldenRequestModel(name string) *llm.Request {
	req := BaseRequest()
	req.Model = "golden/" + name
	return req
}

// ---------------------------------------------------------------------------
// The shared fault suite
// ---------------------------------------------------------------------------

// FaultSuite runs the mockllm-backed fault cases for one adapter unit.
//
// Each unit declares how to reach the endpoint (Path) and which route name
// mockllm counts for it.
type FaultSuite struct {
	// Path is the mockllm endpoint suffix carrying this protocol (the part
	// after /v1), e.g. "/chat/completions".
	Path string
	// Route is the mockllm /__control/state counter name (chat|messages|responses).
	Route string
	// GoldenPrefix is the fixture-name prefix of this unit's golden files, so
	// the suite can ask mockllm to replay them (model "golden/<prefix><sc>").
	GoldenPrefix string
}

// RunFaultSuite executes the injected-fault cases that every adapter must
// survive identically: 429 retry-after honored, 5xx exhausted -> provider,
// 401 -> auth without retry, truncate -> mid-stream disconnect.
func RunFaultSuite(t *testing.T, u Unit, fs FaultSuite) {
	t.Helper()
	proc := StartMockllm(t)
	base := proc.Base + "/v1"

	newProvider := func(o Options) llm.LlmProvider { return u.New(t, base, o) }

	t.Run("fault429_retry_after_honored", func(t *testing.T) {
		proc.Control(t, "/__control/reset", "{}")
		proc.QueueFault(t, Fault{Status: 429, Times: 2, RetryAfter: "1"})
		p := newProvider(Options{})
		rp := llm.NewRetrying(p, RetryOptionsOf())
		start := time.Now()
		_, turn, err := Drain(t, rp, BaseRequest())
		elapsed := time.Since(start)
		if err != nil {
			t.Fatal(err)
		}
		if turn.Text == "" || turn.Stop != llm.StopEndTurn {
			t.Errorf("turn = %+v", turn)
		}
		// Measured threshold (same as the golden ladder): the two
		// retry-after:1 hints sum to 2s; a non-honoring backoff would return
		// in ~10ms with the injected 5ms base.
		if elapsed < 1900*time.Millisecond {
			t.Errorf("elapsed = %v, want >= 1.9s (two retry-after:1 hints honored)", elapsed)
		}
		if n := proc.RouteCount(t, fs.Route); n != 3 {
			t.Errorf("%s requests = %d, want 3", fs.Route, n)
		}
		proc.Control(t, "/__control/reset", "{}")
	})

	t.Run("fault5xx_exhausted_is_provider", func(t *testing.T) {
		proc.Control(t, "/__control/reset", "{}")
		proc.QueueFault(t, Fault{Status: 500, Times: 10})
		rp := llm.NewRetrying(newProvider(Options{}), RetryOptionsOf())
		_, _, err := Drain(t, rp, BaseRequest())
		if err == nil {
			t.Fatal("expected provider error after exhausting retries")
		}
		if ClassOf(err) != observe.ClassProvider {
			t.Errorf("class = %s, want provider", ClassOf(err))
		}
		if n := proc.RouteCount(t, fs.Route); n != 4 {
			t.Errorf("%s requests = %d, want 4 (initial + 3 retries)", fs.Route, n)
		}
		proc.Control(t, "/__control/reset", "{}")
	})

	t.Run("fault401_never_retried", func(t *testing.T) {
		proc.Control(t, "/__control/reset", "{}")
		proc.QueueFault(t, Fault{Status: 401, Times: 10})
		rp := llm.NewRetrying(newProvider(Options{}), RetryOptionsOf())
		_, _, err := Drain(t, rp, BaseRequest())
		if err == nil {
			t.Fatal("expected auth error")
		}
		if ClassOf(err) != observe.ClassAuth {
			t.Errorf("class = %s, want auth", ClassOf(err))
		}
		if n := proc.RouteCount(t, fs.Route); n != 1 {
			t.Errorf("%s requests = %d, want 1 (401 never retried)", fs.Route, n)
		}
		proc.Control(t, "/__control/reset", "{}")
	})

	t.Run("truncate_is_mid_stream_disconnect", func(t *testing.T) {
		proc.Control(t, "/__control/reset", "{}")
		proc.Control(t, "/__control/truncate", `{"chunks":2}`)
		_, turn, err := Drain(t, newProvider(Options{}), BaseRequest())
		if err == nil {
			t.Fatal("expected stream_disconnected error")
		}
		if CodeOf(err) != llm.CodeStreamDisconnect {
			t.Errorf("provider_code = %q, want %q", CodeOf(err), llm.CodeStreamDisconnect)
		}
		// Dialects chunk differently, so the truncation point moves: what MUST
		// hold everywhere is that a delivered partial is retained and marked,
		// and that the stream ends with the failure triple. (The golden
		// disconnect scenarios pin the partial-retention case exactly.)
		if turn.Text != "" && !turn.Incomplete {
			t.Error("truncated turn delivered text but is not marked incomplete")
		}
		if turn.Err == nil {
			t.Error("truncated stream emitted no Error event")
		}
		if turn.Stop != llm.StopError {
			t.Errorf("truncated stream stop = %s, want error", turn.Stop)
		}
		proc.Control(t, "/__control/reset", "{}")
	})

	t.Run("request_wire_shape_understood_by_mockllm", func(t *testing.T) {
		// The strongest available proof that THIS adapter's request body is
		// valid for its protocol: mockllm parses it and echoes the last user
		// turn back, so a mis-encoded messages/input/system field would break
		// the echo text rather than pass silently.
		proc.Control(t, "/__control/reset", "{}")
		req := BaseRequest()
		req.Messages = []llm.Message{{Role: llm.RoleUser,
			Content: []llm.Content{llm.TextPart{Text: "ping-ticket11"}}}}
		_, turn, err := Drain(t, newProvider(Options{}), req)
		if err != nil {
			t.Fatalf("mockllm rejected the adapter request: %v (turn=%+v)", err, turn)
		}
		if !strings.Contains(turn.Text, "ping-ticket11") {
			t.Errorf("mockllm echo = %q, want it to contain the user text "+
				"(the request body did not parse into the protocol's message shape)", turn.Text)
		}
		proc.Control(t, "/__control/reset", "{}")
	})
}

// StreamThroughMockllmGolden replays a golden fixture through mockllm for a
// unit (the second runner), returning the events.
func StreamThroughMockllmGolden(t *testing.T, u Unit, path, fixture string) []llm.StreamEvent {
	t.Helper()
	proc := StartMockllm(t)
	proc.Control(t, "/__control/reset", "{}")
	p := u.New(t, proc.Base+"/v1", Options{})
	events, _, err := Drain(t, p, GoldenRequestModel(fixture))
	if err != nil {
		t.Fatalf("mockllm golden %s: %v", fixture, err)
	}
	return events
}

// ContextOf is a tiny alias used by adapter tests that need a bounded probe
// context.
func ContextOf(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return ctx
}
