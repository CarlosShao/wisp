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

// mockllmReadyBudget is how long StartMockllm waits for the child to print
// its address. It is spent through observe.Timeout, i.e. the monotonic clock
// (D22 ban #4 / D42#9): no wall-clock difference decides readiness.
const mockllmReadyBudget = 20 * time.Second

// mockllmJoinBudget is how long cleanup waits for the stdout reader spawn to
// drain once the child is gone.
const mockllmJoinBudget = 2 * time.Second

// mockllmRegistry is the named-goroutine registry for this helper's spawns
// (D22 ban #1 / D38b). It is deliberately not observe.Default: a test-support
// goroutine must not move the resident-roster numbers the product watchdog and
// its tests read.
var mockllmRegistry = observe.NewRegistry()

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
	// Cleanup is registered before the readiness wait so every failure path
	// below still reaps the child and joins the reader spawn.
	spawnRoot := observe.NewRoot("adaptertest:mockllm")
	ready := make(chan error, 1)
	h := mockllmRegistry.Spawn("mockllm-stdout-reader", "test", spawnRoot,
		func(_ context.Context) {
			sc := bufio.NewScanner(stdout)
			for sc.Scan() {
				if v, ok := strings.CutPrefix(sc.Text(), "MOCKLLM_ADDR="); ok {
					proc.Base = "http://" + v
					ready <- nil
					return
				}
			}
			ready <- fmt.Errorf("mockllm exited before printing its address")
		})
	t.Cleanup(func() {
		_ = proc.cmd.Process.Kill() // without this, Wait blocks forever: mockllm is a server
		_, _ = proc.cmd.Process.Wait()
		if left := spawnRoot.Wait(mockllmJoinBudget); left != 0 {
			t.Errorf("mockllm stdout reader still pending after cleanup: %d", left)
		}
	})
	budget := observe.NewTimeout(mockllmReadyBudget)
	timer := time.NewTimer(budget.Remaining())
	defer timer.Stop()

	var startErr error
	select {
	case startErr = <-ready:
	case <-h.Done():
		// The reader exited without a clean handoff. Either its recover
		// boundary caught a panic - which must reach the test as a failure,
		// not as a 20s hang - or it delivered and returned in the same
		// instant; the buffered channel tells the two apart.
		if perr := h.Err(); perr != nil {
			startErr = fmt.Errorf("mockllm stdout reader failed: %w", perr)
		} else {
			startErr = <-ready
		}
	case <-timer.C:
		startErr = fmt.Errorf("mockllm did not become ready in %v", mockllmReadyBudget)
	}
	if startErr != nil {
		_ = proc.cmd.Process.Kill()
		t.Fatal(startErr)
	}
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

// StreamGolden replays a golden fixture through THIS mockllm for a unit (the
// second runner), returning the events. req may be nil (then a BaseRequest is
// used); when given, its Model is overridden with the golden selector, so a
// test can send a REALISTIC request body (system prompt, breakpoints, tools)
// and still control the response bytes.
func (m *Mockllm) StreamGolden(t *testing.T, u Unit, fixture string, req *llm.Request) []llm.StreamEvent {
	t.Helper()
	m.Reset(t)
	p := u.New(t, m.Base+"/v1", Options{})
	if req == nil {
		req = GoldenRequestModel(fixture)
	} else {
		req.Model = "golden/" + fixture
	}
	events, _, err := Drain(t, p, req)
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

// LastRequest returns the captured request body (raw bytes) mockllm received
// on one route, plus its headers. Fails the test when nothing was captured
// (a stale/missing capture must never be read as evidence).
func (m *Mockllm) LastRequest(t *testing.T, route string) (string, http.Header) {
	t.Helper()
	resp, err := http.Get(m.Base + "/__control/last_request?route=" + route)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	st := m.decode("/__control/last_request", resp)
	present, _ := st["present"].(bool)
	if !present {
		t.Fatalf("mockllm captured no request on route %q", route)
	}
	body, _ := st["body"].(string)
	hdr := http.Header{}
	if raw, ok := st["header"].(map[string]any); ok {
		for k, v := range raw {
			if list, ok := v.([]any); ok {
				for _, s := range list {
					if str, ok := s.(string); ok {
						hdr.Add(k, str)
					}
				}
			}
		}
	}
	return body, hdr
}

// Reset clears every injection, counter and request capture.
func (m *Mockllm) Reset(t *testing.T) {
	t.Helper()
	m.Control(t, "/__control/reset", "{}")
}
