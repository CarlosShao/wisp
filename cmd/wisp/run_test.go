package main

// Ticket 12 AC#1 + AC#4, and rulings A8/A12/A13's reachability proof, run
// against the REAL mockllm over the REAL composition root:
// config.toml -> secret.Store -> llm.Resolver -> provider -> approval.Gate ->
// tools.Bridge -> agent.Loop -> stdout + notification + exit code + rows.
//
// What each case pins:
//
//	TestRunTextTaskTextPathEndToEnd           streamed reply, notification,
//	                                          exit 0, task_log row
//	TestRunTextTaskFailNextIsClassified       fail_next(3) -> error_class,
//	                                          non-zero exit, visible message
//	TestHostDispatchThroughTheAssembledBridge a call dispatched by the HOST on
//	                                          the bridge this root assembled:
//	                                          allowlisted read executes, an
//	                                          out-of-allowlist read is refused
//	                                          as L2 (A12's criterion, now
//	                                          reachable from cmd/wisp), and the
//	                                          rows carry correlation ids
//	TestComposedGateBlocksAWriteForTwoSeconds an L1 fs.write through the SAME
//	                                          stack opens the real block window
//	                                          (A13: nothing like this could
//	                                          happen on a real machine before)
//	TestRunTextTaskKeyResolvesInTheStore      A8: the Authorization header the
//	                                          PROVIDER set came out of the store
//	TestMissingBlobFailsUnconfigured          ...and a missing blob fails
//	                                          through the Unconfigured path

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/secret"
)

// fakeStoreKey is the fixed fake constant this suite puts into the store. It is
// not a real credential, and it is the only place a key value appears in a test
// string here: every assertion names the HEADER, never a production secret.
const fakeStoreKey = "wire-fake-key-0123456789abcdef"

// taskIDRe pulls the loop's task id out of the run's own status line, so the
// row assertions query the task this very run created.
var taskIDRe = regexp.MustCompile(`任务 (\S+) 结束`)

type runFixture struct {
	t   *testing.T
	srv *adaptertest.Mockllm
	dir string
	out *bytes.Buffer
	err *bytes.Buffer
	// notices records every poster call (title + NUL + body) this run made.
	notices []string
	// rtHook receives the assembled stack while it is still open.
	rtHook func(*agentRuntime)
	// calls records what the bridge booked, read back after the run.
	mu    sync.Mutex
	calls []string
}

// newRunFixture writes a config.toml into a fresh data dir (allowlisting that
// dir as the only fs root), stores the fake key under dpapi:acme, and starts
// mockllm.
func newRunFixture(t *testing.T, protocol string) *runFixture {
	t.Helper()
	f := &runFixture{t: t, srv: adaptertest.StartMockllm(t), out: &bytes.Buffer{}, err: &bytes.Buffer{}}
	f.dir = t.TempDir()
	cfgPath := filepath.Join(f.dir, configFileName)
	body := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/m1"]

[llm.retry]
max = 1
backoff_ms = 1

[llm.providers.acme]
protocol = %q
base_url = %q
api_key_ref = "dpapi:acme"

[llm.providers.acme.models.m1]
context_window = 128000

[fs]
allowed_dirs = [%q]
`, protocol, f.srv.Base+"/v1", filepath.ToSlash(f.dir))
	if err := os.WriteFile(cfgPath, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := secret.NewStore(f.dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Store("dpapi:acme", fakeStoreKey); err != nil {
		t.Fatal(err)
	}
	return f
}

// run drives the composition root exactly as `wisp run` does, with only the
// notification poster swapped for a recorder.
func (f *runFixture) run(task string) int {
	f.t.Helper()
	return runTextTask(runSpec{
		argv:      []string{task},
		stdout:    f.out,
		stderr:    f.err,
		dataDir:   f.dir,
		onRuntime: f.rtHook,
		notify: func(title, body string) error {
			f.mu.Lock()
			defer f.mu.Unlock()
			f.notices = append(f.notices, title+"\x00"+body)
			return nil
		},
	})
}

// taskID reads the id this run reported.
func (f *runFixture) taskID() string {
	f.t.Helper()
	m := taskIDRe.FindStringSubmatch(f.out.String())
	if m == nil {
		f.t.Fatalf("the run printed no status line:\n%s", f.out.String())
	}
	return m[1]
}

// openStore reopens the run's store after the runtime closed it.
func (f *runFixture) openStore() *memory.Store {
	f.t.Helper()
	st, err := memory.Open(f.dir)
	if err != nil {
		f.t.Fatal(err)
	}
	f.t.Cleanup(func() { _ = st.Close() })
	return st
}

func TestRunTextTaskTextPathEndToEnd(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	code := f.run("总结一下 这份笔记讲了什么")
	out, errb := f.out.String(), f.err.String()
	if code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, out, errb)
	}
	// The reply streams: the provider's text arrives before the status line.
	replyAt := strings.Index(out, "echo:")
	statusAt := strings.Index(out, "结束")
	if replyAt < 0 || statusAt < 0 {
		t.Fatalf("the run must stream a reply and then close with a status line:\n%s", out)
	}
	if replyAt > statusAt {
		t.Errorf("the reply was not streamed before the status line:\n%s", out)
	}
	if !strings.Contains(out, "轮") || !strings.Contains(out, "成本") {
		t.Errorf("missing the status/cost line the CLI owes (ticket 12 Key constraints):\n%s", out)
	}
	// D10: the reply also arrives as a notification, because a CLI run has no
	// ball to pop.
	f.mu.Lock()
	notices := append([]string(nil), f.notices...)
	f.mu.Unlock()
	if len(notices) == 0 {
		t.Fatal("no notification was posted")
	}
	var posted bool
	for _, n := range notices {
		if strings.Contains(n, "echo:") {
			posted = true
		}
	}
	if !posted {
		t.Errorf("the reply was never presented as a notification: %v", notices)
	}

	// AC#4 (first half): a task_log row for this task, closed with its state.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id := f.taskID()
	tl, err := f.openStore().TaskLogByID(ctx, id)
	if err != nil {
		t.Fatalf("task_log for %s: %v", id, err)
	}
	if tl.State != "done" || tl.QueryText == "" || tl.EndedAt == nil {
		t.Errorf("task_log = %+v, want a closed done row carrying the query text", tl)
	}
}

func TestRunTextTaskFailNextIsClassified(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	// fail_next(3): every attempt fails, so the retry ladder has to give up and
	// the task has to end classified.
	f.srv.QueueFault(t, adaptertest.Fault{Status: 500, Times: 3})
	code := f.run("总结一下 这份笔记")
	out, errb := f.out.String(), f.err.String()
	if code == 0 {
		t.Fatalf("a run whose provider failed must not exit 0:\n%s\n%s", out, errb)
	}
	// A provider failure is neither Unconfigured (2) nor a loop brake (3): the
	// classified-failure code is 1.
	if code != 1 {
		t.Errorf("exit = %d, want 1 for a classified provider failure\n%s", code, errb)
	}
	if !strings.Contains(errb, "error_class=provider") {
		t.Errorf("the error class is not on stderr:\n%s", errb)
	}
	if !strings.Contains(out, "wisp run:") || !strings.Contains(out, "failed") {
		t.Errorf("the user must get a status line, not only an exit code:\n%s", out)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tl, err := f.openStore().TaskLogByID(ctx, f.taskID())
	if err != nil {
		t.Fatalf("a failed task must still close its task_log row: %v", err)
	}
	if tl.State != "error" || tl.ErrorClass != "provider" {
		t.Errorf("task_log = state %q class %q, want error/provider", tl.State, tl.ErrorClass)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.notices) == 0 {
		t.Error("a failed task that posts nothing is invisible on the desktop (D10)")
	}
}

// TestHostDispatchThroughTheAssembledBridge dispatches host-side calls on the
// bridge `wisp run` built, while the run's store is still open (the composition
// root closes it when the task ends, so this is the only window where the
// assembled stack is live).
func TestHostDispatchThroughTheAssembledBridge(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	source := filepath.Join(f.dir, "note.txt")
	if err := os.WriteFile(source, []byte("notes body"), 0o600); err != nil {
		t.Fatal(err)
	}
	var (
		readText  string
		readErr   error
		outsideOK bool
	)
	f.rtHook = func(rt *agentRuntime) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// (1) allowlisted read: executes through C3 + C19 + the L0 route.
		out, err := rt.bridge.Execute(ctx, agent.ToolRequest{
			TaskID: "host-task", CorrelationID: "host-corr", CallID: "c1",
			Name: "fs.read",
			Args: json.RawMessage(fmt.Sprintf(`{"path":%q}`, filepath.ToSlash(source))),
		})
		readText, readErr = out.Text, err
		if out.IsError {
			t.Errorf("an allowlisted read must execute, got: %s", out.Text)
		}
		// (2) outside the allowlist: R2 raises it to L2, the D47 registration
		// this caller never got means the gate refuses, and the refusal is
		// booked. This is the end-to-end shape A12 said only existed in-package.
		out2, err2 := rt.bridge.Execute(ctx, agent.ToolRequest{
			TaskID: "host-task-2", CorrelationID: "host-corr-2", CallID: "c2",
			Name: "fs.read",
			Args: json.RawMessage(`{"path":"C:/Windows/win.ini"}`),
		})
		outsideOK = err2 == nil && out2.IsError
	}
	if code := f.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\n%s\n%s", code, f.out.String(), f.err.String())
	}
	if readErr != nil {
		t.Fatal(readErr)
	}
	if readText != "notes body" {
		t.Errorf("read returned %q, want the file body", readText)
	}
	if !outsideOK {
		t.Error("the out-of-allowlist read was not refused as an error outcome")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	st := f.openStore()
	rows, err := st.ListToolCallsByTask(ctx, "host-task")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %+v, want the one authoritative row for the read", rows)
	}
	r := rows[0]
	if r.Tool != "fs.read" || r.RiskLevel != memory.RiskL0 ||
		r.Decision != agent.DecisionAllow || r.Outcome != agent.OutcomeSuccess {
		t.Errorf("row = %+v, want fs.read/L0/allow/success", r)
	}
	if r.CorrelationID != "host-corr" {
		t.Errorf("correlation_id = %q, want host-corr (C18)", r.CorrelationID)
	}
	rows2, err := st.ListToolCallsByTask(ctx, "host-task-2")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows2) != 1 {
		t.Fatalf("rows2 = %+v", rows2)
	}
	if rows2[0].RiskLevel != memory.RiskL2 || rows2[0].Decision != agent.DecisionReject {
		t.Errorf("out-of-allowlist row = %+v, want L2/reject (R2 through C26)", rows2[0])
	}
	if rows2[0].CorrelationID != "host-corr-2" {
		t.Errorf("correlation_id = %q, want host-corr-2", rows2[0].CorrelationID)
	}
}

// TestComposedGateBlocksAWriteForTwoSeconds runs an L1 fs.write through the
// stack `wisp run` assembles - the same console UI, the same channel registry,
// the same window length from [risk].l1_window_sec. Nothing like this could
// happen on a real machine before this ticket: the gate had no constructor
// outside its own tests (A13).
func TestComposedGateBlocksAWriteForTwoSeconds(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	target := filepath.Join(f.dir, "composed-write.txt")
	var (
		blocked    time.Duration
		text       string
		isErr      bool
		cards      int
		unAdmitted bool
	)
	f.rtHook = func(rt *agentRuntime) {
		args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"landed through the window"}`,
			filepath.ToSlash(target)))
		// (0) D47 first: a task the text loop never registered may not reach a
		// gate at all. This is the fail-closed half of the pairing, and it is
		// the state every host-side caller starts in.
		out0, err0 := rt.bridge.Execute(context.Background(), agent.ToolRequest{
			TaskID: "gate-task", CorrelationID: "gate-corr-0", CallID: "w0",
			Name: "fs.write", Args: args,
		})
		unAdmitted = err0 == nil && out0.IsError
		if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
			t.Errorf("an un-registered task wrote to disk: %v", statErr)
		}
		// (1) register it the way the loop does, and the same call routes into
		// the real L1 window.
		revoke := rt.gate.AdmitTextTask("gate-task")
		defer revoke()
		start := time.Now()
		out, err := rt.bridge.Execute(context.Background(), agent.ToolRequest{
			TaskID: "gate-task", CorrelationID: "gate-corr", CallID: "w1",
			Name: "fs.write",
			Args: args,
		})
		blocked = time.Since(start)
		text, isErr = out.Text, out.IsError
		if err != nil {
			t.Fatalf("gate route: %v", err)
		}
		cards = rt.windowCount()
	}
	if code := f.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\n%s\n%s", code, f.out.String(), f.err.String())
	}
	if !unAdmitted {
		t.Error("a task that never went through AdmitTextTask must be refused (D47)")
	}
	if isErr {
		t.Fatalf("an unvetoed L1 window means EXECUTE, got: %s", text)
	}
	body, err := os.ReadFile(target)
	if err != nil || string(body) != "landed through the window" {
		t.Fatalf("the write did not land: %v %q", err, body)
	}
	if blocked < 2*time.Second {
		t.Errorf("the call returned after %v: the 2s L1 window must actually block", blocked)
	}
	if cards != 1 {
		t.Errorf("the composed gate displayed %d cards, want exactly the one L1 window", cards)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := f.openStore().ListToolCallsByTask(ctx, "gate-task")
	if err != nil {
		t.Fatal(err)
	}
	// Two rows, and the pairing is the point: the un-registered call is booked
	// L1/reject/user_rejected with no window, the registered one L1/allow.
	var gated, refused *memory.ToolCall
	for i := range rows {
		switch rows[i].CorrelationID {
		case "gate-corr":
			gated = &rows[i]
		case "gate-corr-0":
			refused = &rows[i]
		}
	}
	if gated == nil || gated.RiskLevel != memory.RiskL1 ||
		gated.Decision != agent.DecisionAllow || gated.Outcome != agent.OutcomeSuccess {
		t.Fatalf("the gated write must be booked as L1/allow/success: %+v", rows)
	}
	if refused == nil || refused.Decision != agent.DecisionReject ||
		refused.ErrorClass != "user_rejected" || refused.Outcome != agent.OutcomeError {
		t.Fatalf("the D47 refusal must be booked as reject/user_rejected: %+v", rows)
	}
}

// TestRunTextTaskKeyResolvesInTheStore is A8's box. The request is issued by
// internal/llm's own adapter, and the assertion reads the header off the
// server's capture of that request - not off a client this test wrote.
func TestRunTextTaskKeyResolvesInTheStore(t *testing.T) {
	// openai-responses because it is the adapter whose route records its
	// request headers (mockllm /__control/last_request); the resolution path
	// under test is protocol-independent.
	f := newRunFixture(t, "openai-responses")
	if code := f.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\n%s\n%s", code, f.out.String(), f.err.String())
	}
	_, hdr := f.srv.LastRequest(t, "responses")
	got := hdr.Get("Authorization")
	if got == "" {
		t.Fatal("the provider sent no Authorization header at all")
	}
	if got != "Bearer "+fakeStoreKey {
		t.Error("the provider's Authorization header is not the value the store holds")
	}
	if strings.Contains(f.out.String(), fakeStoreKey) ||
		strings.Contains(f.err.String(), fakeStoreKey) {
		t.Error("the key value leaked into the run's output")
	}
}

// TestMissingBlobFailsUnconfiguredNeverSilently is A8's mutation criterion:
// point the ref at a blob that does not exist and the run must fail through the
// Unconfigured path. "Silently succeed" would look like exit 0 plus a request
// carrying no (or a fallback) credential, so the request count is asserted too.
func TestMissingBlobFailsUnconfiguredNeverSilently(t *testing.T) {
	f := newRunFixture(t, "openai-responses")
	st, err := secret.NewStore(f.dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Delete("dpapi:acme"); err != nil {
		t.Fatal(err)
	}
	f.srv.Reset(t)
	code := f.run("总结一下 这份笔记")
	if code == 0 {
		t.Fatalf("a missing credential blob must never exit 0:\n%s", f.out.String())
	}
	if code != 2 {
		t.Errorf("exit = %d, want 2 (the Unconfigured class)", code)
	}
	if !strings.Contains(f.err.String(), "Unconfigured") {
		t.Errorf("the failure must name the Unconfigured path:\n%s", f.err.String())
	}
	if n := f.srv.RouteCount(t, "responses"); n != 0 {
		t.Errorf("%d requests reached the provider with an unresolvable key", n)
	}
	if strings.Contains(f.err.String(), fakeStoreKey) ||
		strings.Contains(f.out.String(), fakeStoreKey) {
		t.Error("the failure output must not carry a key value")
	}
}
