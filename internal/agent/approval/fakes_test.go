package approval_test

import (
	"context"
	"encoding/json"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// ---------------------------------------------------------------------------
// fake monotonic clock
// ---------------------------------------------------------------------------

// fakeClock is the ticket-11 pattern (internal/llm/ratelimit.go's injected
// now): one source for both the stamps and the timers, advanced by the test.
// Every timeout in the layer is a duration handed to After, so a test can
// prove "300s means auto-reject" without sleeping 300 seconds and without a
// single wall-clock comparison.
type fakeClock struct {
	mu     sync.Mutex
	now    time.Time
	due    []*fakeTimer
	maxDue time.Duration
}

type fakeTimer struct {
	at    time.Time
	c     chan time.Time
	fired bool
}

func newFakeClock() *fakeClock { return &fakeClock{now: time.Unix(0, 0)} }

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) After(d time.Duration) <-chan time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	t := &fakeTimer{at: c.now.Add(d), c: make(chan time.Time, 1)}
	if d <= 0 {
		t.c <- c.now
		t.fired = true
		return t.c
	}
	c.due = append(c.due, t)
	return t.c
}

// Advance moves the single time source forward and fires everything now due,
// oldest first, so a warn-then-deadline sequence is observed in order.
func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	if d > c.maxDue {
		c.maxDue = d
	}
	c.now = c.now.Add(d)
	var fire []*fakeTimer
	rest := c.due[:0]
	for _, t := range c.due {
		if !t.at.After(c.now) {
			fire = append(fire, t)
			continue
		}
		rest = append(rest, t)
	}
	sortTimers(fire)
	c.due = rest
	for _, t := range fire {
		if !t.fired {
			t.fired = true
			t.c <- t.at
		}
	}
	c.mu.Unlock()
}

func sortTimers(ts []*fakeTimer) {
	for i := 1; i < len(ts); i++ {
		for j := i; j > 0 && ts[j].at.Before(ts[j-1].at); j-- {
			ts[j], ts[j-1] = ts[j-1], ts[j]
		}
	}
}

// ---------------------------------------------------------------------------
// fake UI (segment 2 owns the real one)
// ---------------------------------------------------------------------------

type fakeUI struct {
	mu       sync.Mutex
	prompts  []approval.Prompt
	events   []approval.Event
	views    []approval.PanelItem
	prompted chan approval.Prompt
	err      error
}

func newFakeUI() *fakeUI {
	return &fakeUI{prompted: make(chan approval.Prompt, 16)}
}

func (u *fakeUI) setErr(err error) {
	u.mu.Lock()
	u.err = err
	u.mu.Unlock()
}

func (u *fakeUI) Prompt(_ context.Context, p approval.Prompt) error {
	u.mu.Lock()
	err := u.err
	u.prompts = append(u.prompts, p)
	u.mu.Unlock()
	u.prompted <- p
	return err
}

func (u *fakeUI) Update(_ context.Context, e approval.Event) error {
	u.mu.Lock()
	u.events = append(u.events, e)
	u.mu.Unlock()
	return nil
}

func (u *fakeUI) wait(t *testing.T) approval.Prompt {
	t.Helper()
	select {
	case p := <-u.prompted:
		return p
	case <-time.After(2 * time.Second):
		t.Fatal("确认界面从未收到提示（等待 2s 后放弃）")
		return approval.Prompt{}
	}
}

func (u *fakeUI) all() []approval.Prompt {
	u.mu.Lock()
	defer u.mu.Unlock()
	return append([]approval.Prompt(nil), u.prompts...)
}

func (u *fakeUI) ofKind(k approval.EventKind) []approval.Event {
	u.mu.Lock()
	defer u.mu.Unlock()
	var out []approval.Event
	for _, e := range u.events {
		if e.Kind == k {
			out = append(out, e)
		}
	}
	return out
}

func (u *fakeUI) count() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	return len(u.prompts)
}

// ---------------------------------------------------------------------------
// shared wiring
// ---------------------------------------------------------------------------

// spawn runs fn with a recover boundary owned by the test, and waits for it in
// cleanup. No test goroutine is allowed to outlive its test (the repo rule is
// owner + recover; this is the test-shaped version of it).
func spawn(t *testing.T, fn func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("test goroutine panicked: %v", r)
			}
		}()
		fn()
	}()
	t.Cleanup(func() {
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Errorf("派生的测试协程未在 3s 内退出（闸门存在无限等待）")
		}
	})
}

type answer struct {
	a   tools.Answer
	why string
}

// runWindow calls PendingWindow off the test goroutine and returns a channel
// that receives its answer.
func runWindow(t *testing.T, g *approval.Gate, ctx context.Context, d tools.Decision) <-chan answer {
	t.Helper()
	out := make(chan answer, 1)
	spawn(t, func() {
		a, why := g.PendingWindow(ctx, d)
		out <- answer{a, why}
	})
	return out
}

func runApproval(t *testing.T, g *approval.Gate, ctx context.Context, d tools.Decision) <-chan answer {
	t.Helper()
	out := make(chan answer, 1)
	spawn(t, func() {
		a, why := g.PendingApproval(ctx, d)
		out <- answer{a, why}
	})
	return out
}

// mustAnswer waits for one gate answer with a bounded test-side wait.
func mustAnswer(t *testing.T, ch <-chan answer) answer {
	t.Helper()
	select {
	case a := <-ch:
		return a
	case <-time.After(2 * time.Second):
		t.Fatal("闸门未在 2s 内给出答复（说明存在无限等待）")
		return answer{}
	}
}

func newGate(t *testing.T, ui approval.UI, opts ...approval.Options) (*approval.Gate, *fakeClock, *approval.ChannelRegistry) {
	t.Helper()
	clk := newFakeClock()
	ch := approval.NewChannels(approval.ChannelBall, approval.ChannelEsc)
	o := approval.Options{UI: ui, Clock: clk, Channels: ch}
	if len(opts) > 0 {
		if opts[0].UI != nil {
			o.UI = opts[0].UI
		}
		if opts[0].Clock != nil {
			o.Clock = opts[0].Clock
		}
		if opts[0].Channels != nil {
			o.Channels = opts[0].Channels
		}
		o.Window = opts[0].Window
		o.ApprovalTimeout = opts[0].ApprovalTimeout
		o.WarningLead = opts[0].WarningLead
		o.MaxPending = opts[0].MaxPending
		o.MaxTracked = opts[0].MaxTracked
		if opts[0].Logf != nil {
			o.Logf = opts[0].Logf
		}
	}
	g := approval.New(o)
	t.Cleanup(func() {
		if n := g.Queue().Depth(); n != 0 {
			t.Errorf("闸门退出时仍有 %d 个待审批项未清理", n)
		}
	})
	return g, clk, ch
}

const (
	testTask = "task-text-1"
	testCorr = "corr-1"
)

func l1Decision(paths ...string) tools.Decision {
	return tools.Decision{
		Tool: "fs.write", Provider: tools.KindBuiltin,
		Params:        map[string]any{"paths": pathsAny(paths)},
		Args:          json.RawMessage(`{"paths":["a"]}`),
		Level:         risk.L1,
		RulesHit:      []risk.RuleID{risk.R1},
		Reason:        "声明 L1：可逆写",
		Paths:         paths,
		CorrelationID: testCorr,
		TaskID:        testTask,
		CallID:        "call-1",
	}
}

func l2Decision(paths ...string) tools.Decision {
	d := l1Decision(paths...)
	d.Level = risk.L2
	d.RulesHit = []risk.RuleID{risk.R1, risk.R2}
	d.Reason = "目标越出 [fs] 授权目录（R2），需 L2 强确认"
	return d
}

func pathsAny(paths []string) []any {
	out := make([]any, 0, len(paths))
	for _, p := range paths {
		out = append(out, p)
	}
	return out
}

func manyPaths(n int) []string {
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, string(rune('A'+i%26))+"/dir"+itoa(i/26)+"/file"+itoa(i)+".txt")
	}
	return out
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b []byte
	for i > 0 {
		b = append([]byte{byte('0' + i%10)}, b...)
		i /= 10
	}
	return string(b)
}

func channelText(p approval.Prompt, ch approval.Channel) (approval.ChannelStatus, bool) {
	for _, s := range p.Channels {
		if s.Channel == ch {
			return s, true
		}
	}
	return approval.ChannelStatus{}, false
}

// ---------------------------------------------------------------------------
// a fake gated tool, so the assertions run through the real bridge
// ---------------------------------------------------------------------------

type fakeTool struct {
	name    string
	cap     tools.Capability
	decl    risk.Level
	pathKey string

	mu    sync.Mutex
	calls int
	ran   chan struct{}
	steps []string
}

func newFakeTool(name string, c tools.Capability, lvl risk.Level, pathKey string) *fakeTool {
	return &fakeTool{name: name, cap: c, decl: lvl, pathKey: pathKey, ran: make(chan struct{}, 8)}
}

func (f *fakeTool) Name() string { return f.name }

func (f *fakeTool) Description() string { return "票 21 测试用工具" }

func (f *fakeTool) Parameters() tools.JSONSchema {
	return json.RawMessage(`{"type":"object"}`)
}

func (f *fakeTool) Execute(ctx context.Context, params json.RawMessage, _ func(string)) (tools.Result, error) {
	f.mu.Lock()
	f.calls++
	steps := append([]string(nil), f.steps...)
	f.mu.Unlock()
	select {
	case f.ran <- struct{}{}:
	default:
	}
	return tools.Result{Text: "已执行 " + f.name, AppliedSteps: steps}, nil
}

func (f *fakeTool) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.calls
}

func (f *fakeTool) setSteps(s []string) {
	f.mu.Lock()
	f.steps = s
	f.mu.Unlock()
}

// bridgeFor wires the gate under the real bridge, which is the only dispatch
// path the loop has (agent.ToolProvider). Journal is nil so the unit test does
// not need a store; the loop-side booking is ticket 12's composition.
func bridgeFor(t *testing.T, g tools.Gate, ft *fakeTool, allowed []string) *tools.Bridge {
	t.Helper()
	reg := tools.NewRegistry()
	if err := reg.Register(tools.Entry{
		Tool: ft,
		Decl: tools.Decl{
			Capabilities: []tools.Capability{ft.cap},
			Needs:        []tools.Capability{ft.cap},
			Declared:     ft.decl,
			PathParams:   []string{ft.pathKey},
			Provider:     tools.KindBuiltin,
			Resident:     true,
		},
	}); err != nil {
		t.Fatalf("Register: %v", err)
	}
	paths := tools.NewPathCanonicalizer(allowed, nil)
	return tools.New(tools.Options{
		Registry: reg, Gate: g, Paths: paths,
		Logf: func(string, ...any) {},
	})
}

var _ agent.ToolProvider = (*tools.Bridge)(nil)

// agentToolRequest builds the request the loop hands the bridge. CorrelationID
// is what C18 routes replies by, so the tests name it explicitly.
// argsJSON renders a call's arguments for the fake tool.
func argsJSON(t *testing.T, paths ...string) string {
	t.Helper()
	b, err := json.Marshal(map[string]any{"paths": pathsAny(paths)})
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// slashPath is the forward-slash spelling a model would emit; C26 owns making
// it canonical, so tests never compare against what the OS handed them.
func slashPath(p string) string { return filepath.ToSlash(p) }

func agentToolRequest(corr, args string) agent.ToolRequest {
	return agent.ToolRequest{
		TaskID: testTask, CorrelationID: corr, CallID: "call-" + corr,
		Name: "fs.write", Args: json.RawMessage(args),
		Timeout: 5 * time.Second,
	}
}
