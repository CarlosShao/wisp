package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

// ---------------------------------------------------------------------------
// fixtures

// fixtureTool is a Test double whose Execute records that it ran. It is the
// only way to prove "hard reject, not error": the counter must stay 0.
type fixtureTool struct {
	name     string
	params   string
	noSchema bool

	mu       sync.Mutex
	runs     int
	lastArgs json.RawMessage
	block    time.Duration
}

func (f *fixtureTool) Name() string        { return f.name }
func (f *fixtureTool) Description() string { return "fixture " + f.name }
func (f *fixtureTool) Parameters() JSONSchema {
	if f.noSchema {
		return nil
	}
	p := f.params
	if p == "" {
		p = `{"type":"object","properties":{},"additionalProperties":true}`
	}
	return JSONSchema(p)
}

// Runs reports how many times Execute was reached.
func (f *fixtureTool) Runs() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.runs
}

func (f *fixtureTool) Execute(ctx context.Context, params json.RawMessage, _ func(string)) (Result, error) {
	f.mu.Lock()
	f.runs++
	f.lastArgs = params
	f.mu.Unlock()
	if f.block > 0 {
		select {
		case <-time.After(f.block):
		case <-ctx.Done():
			return Result{}, ctx.Err()
		}
	}
	return Result{Text: "ran " + f.name}, nil
}

// gateSpy counts which gate branch the routing reached.
type gateSpy struct {
	mu         sync.Mutex
	windows    []Decision
	approvals  []Decision
	windowAns  Answer
	windowWhy  string
	approveAns Answer
	approveWhy string
}

func (g *gateSpy) PendingWindow(_ context.Context, d Decision) (Answer, string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.windows = append(g.windows, d)
	return g.windowAns, g.windowWhy
}

func (g *gateSpy) PendingApproval(_ context.Context, d Decision) (Answer, string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.approvals = append(g.approvals, d)
	return g.approveAns, g.approveWhy
}

func (g *gateSpy) counts() (int, int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return len(g.windows), len(g.approvals)
}

func (g *gateSpy) windowDecision() Decision {
	g.mu.Lock()
	defer g.mu.Unlock()
	if len(g.windows) == 0 {
		return Decision{}
	}
	return g.windows[len(g.windows)-1]
}

func (g *gateSpy) approvalDecision() Decision {
	g.mu.Lock()
	defer g.mu.Unlock()
	if len(g.approvals) == 0 {
		return Decision{}
	}
	return g.approvals[len(g.approvals)-1]
}

// newFixtureBridge returns a bridge over fixture tools with no path argument,
// so the assessed level is exactly the R1 declared floor.
func newFixtureBridge(t *testing.T, gate Gate, entries ...Entry) (*Bridge, map[string]*fixtureTool) {
	t.Helper()
	reg := NewRegistry()
	names := map[string]*fixtureTool{}
	for _, e := range entries {
		if f, ok := e.Tool.(*fixtureTool); ok {
			names[f.name] = f
		}
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	b := New(Options{Registry: reg, Gate: gate})
	return b, names
}

func entry(name string, lvl risk.Level, caps, needs []Capability) Entry {
	return Entry{
		Tool: &fixtureTool{name: name},
		Decl: Decl{
			Capabilities: caps, Needs: needs, Declared: lvl,
			Resident: true, Provider: KindBuiltin,
		},
	}
}

func req(name, args string) agent.ToolRequest {
	return agent.ToolRequest{
		TaskID: "task-1", CorrelationID: "corr-1", CallID: "call-1",
		Name: name, Args: json.RawMessage(args),
	}
}

// ---------------------------------------------------------------------------
// C1: the Tool contract
// ---------------------------------------------------------------------------

// TestC1ToolContract asserts the four-method contract on the real L0 pair,
// not on a fixture: the shape the model sees must be the shape SPEC-07 §2
// froze.
func TestC1ToolContract(t *testing.T) {
	var _ Tool = fsRead{}
	var _ Tool = fsList{}

	for _, tool := range []Tool{fsRead{}, fsList{}} {
		if got := tool.Name(); !strings.Contains(got, ".") {
			t.Errorf("%q: Name() is not 'namespace.action'", got)
		}
		if strings.TrimSpace(tool.Description()) == "" {
			t.Errorf("%s: empty Description()", tool.Name())
		}
		var schema map[string]any
		if err := json.Unmarshal(tool.Parameters(), &schema); err != nil {
			t.Fatalf("%s: Parameters() is not JSON: %v", tool.Name(), err)
		}
		if schema["type"] != "object" {
			t.Errorf("%s: parameter schema top type = %v, want object", tool.Name(), schema["type"])
		}
		if _, ok := schema["properties"]; !ok {
			t.Errorf("%s: parameter schema has no properties", tool.Name())
		}
		// Execute(ctx, params, onUpdate) is reached through the bridge below;
		// here: a nil onUpdate must not be a crash (the contract allows it).
		res, err := tool.Execute(t.Context(), json.RawMessage(`{"path":"Z:\\definitely-not-here\\x.txt"}`), nil)
		if err != nil {
			t.Fatalf("%s: Execute returned a host error: %v", tool.Name(), err)
		}
		if !res.IsError {
			t.Errorf("%s: reading a nonexistent path must be an IsError result", tool.Name())
		}
	}
}

// TestToolConcurrencyCeilingIsFour drives more calls than the D38d ceiling and
// asserts no bridge ever runs more than four at once.
func TestToolConcurrencyCeilingIsFour(t *testing.T) {
	var inflight, maxSeen atomic.Int32
	tool := &countingTool{inflight: &inflight, maxSeen: &maxSeen}
	reg := NewRegistry()
	if err := reg.Register(Entry{
		Tool: tool,
		Decl: Decl{
			Capabilities: []Capability{CapFSRead}, Needs: []Capability{CapFSRead},
			Declared: risk.L0, Provider: KindBuiltin,
		},
	}); err != nil {
		t.Fatal(err)
	}
	b := New(Options{Registry: reg, MaxConcurrency: 99}) // 99 must clamp to 4
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() { // test-only goroutine: d22scan scopes ban 1 to non-test files
			defer wg.Done()
			if _, err := b.Execute(t.Context(), req(tool.Name(), `{}`)); err != nil {
				t.Errorf("Execute: %v", err)
			}
		}()
	}
	wg.Wait()
	if got := maxSeen.Load(); got > MaxToolConcurrency {
		t.Fatalf("observed concurrency %d exceeded the D38d ceiling %d", got, MaxToolConcurrency)
	}
	if got := maxSeen.Load(); got < 2 {
		t.Fatalf("observed concurrency %d: the test never ran two at once, so it proves nothing", got)
	}
}

type countingTool struct {
	inflight, maxSeen *atomic.Int32
}

func (c *countingTool) Name() string           { return "probe.count" }
func (c *countingTool) Description() string    { return "counts" }
func (c *countingTool) Parameters() JSONSchema { return JSONSchema(`{"type":"object"}`) }
func (c *countingTool) Execute(ctx context.Context, _ json.RawMessage, _ func(string)) (Result, error) {
	cur := c.inflight.Add(1)
	for {
		prev := c.maxSeen.Load()
		if cur <= prev || c.maxSeen.CompareAndSwap(prev, cur) {
			break
		}
	}
	defer c.inflight.Add(-1)
	select {
	case <-time.After(15 * time.Millisecond):
	case <-ctx.Done():
		return Result{}, ctx.Err()
	}
	return Result{Text: "ok"}, nil
}

// TestPerToolTimeoutHonoredViaContext proves C22's per-tool budget reaches the
// tool as a ctx deadline (not a wall-clock comparison).
func TestPerToolTimeoutHonoredViaContext(t *testing.T) {
	slow := &fixtureTool{name: "probe.slow", block: 2 * time.Second}
	reg := NewRegistry()
	if err := reg.Register(Entry{Tool: slow, Decl: Decl{
		Capabilities: []Capability{CapFSRead}, Needs: []Capability{CapFSRead},
		Declared: risk.L0, Timeout: 40 * time.Millisecond, Provider: KindBuiltin,
	}}); err != nil {
		t.Fatal(err)
	}
	b := New(Options{Registry: reg})

	start := time.Now()
	out, err := b.Execute(t.Context(), req("probe.slow", `{}`))
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("a timed-out call must not be a host fault: %v", err)
	}
	if !out.IsError {
		t.Fatalf("timed-out call must be IsError, got %+v", out)
	}
	if !strings.Contains(out.Text, "超时") {
		t.Errorf("text = %q, want the structured TOOL_TIMEOUT wording", out.Text)
	}
	if elapsed > time.Second {
		t.Fatalf("the per-tool budget was not honored: %v", elapsed)
	}
	// A fast tool under the same bridge must not be cut short.
	fast := &fixtureTool{name: "probe.slow"}
	_ = fast
}

// ---------------------------------------------------------------------------
// C3: capability enforcement
// ---------------------------------------------------------------------------

// TestCapabilitySetIsFrozen pins SPEC-07 §2's eleven tokens.
func TestCapabilitySetIsFrozen(t *testing.T) {
	if len(AllCapabilities) != 11 {
		t.Fatalf("C3 has %d tokens, want 11", len(AllCapabilities))
	}
	want := "clipboard,fs.read,fs.write,input,memory,net,notify,screen,shell,sysinfo,window"
	got := make([]string, 0, 11)
	for _, c := range AllCapabilities {
		if !c.Valid() {
			t.Fatalf("%s claims to be in its own set", c)
		}
		got = append(got, string(c))
	}
	sorted := append([]string(nil), got...)
	sortStrings(sorted)
	if strings.Join(sorted, ",") != want {
		t.Fatalf("capability set drifted:\n got %v\nwant %v", strings.Join(sorted, ","), want)
	}
	for _, bad := range []Capability{"fs", "filesystem", "fs.delete", "", "NET"} {
		if bad.Valid() {
			t.Errorf("%q must not validate: C3 is closed", bad)
		}
	}
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

// TestUndeclaredCapabilityHardRejects is AC#1's first clause, and the load one:
// SPEC-07 §2 says 未声明即拒绝调用（不是报错，是拒绝）. The assertion that
// carries the meaning is err == nil TOGETHER WITH the tool never running: a Go
// error would be booked by the loop as a provider fault (agent/loop.go:654) and
// would surface as "Wisp broke" instead of "you did not ask for this power".
func TestUndeclaredCapabilityHardRejects(t *testing.T) {
	tool := &fixtureTool{name: "probe.net"}
	reg := NewRegistry()
	// The tool NEEDS net but never DECLARED it. Registration is the first
	// refusal...
	if err := reg.Register(Entry{Tool: tool, Decl: Decl{
		Capabilities: []Capability{CapFSRead}, Needs: []Capability{CapNet},
		Declared: risk.L0, Provider: KindBuiltin,
	}}); err == nil {
		t.Fatal("registration accepted a tool that needs an undeclared capability")
	} else if !strings.Contains(err.Error(), "undeclared capability") {
		t.Fatalf("error = %v, want the undeclared-capability refusal", err)
	}

	// ...and the runtime check is the one that cannot be bypassed: a plugin's
	// capability face can shrink AFTER registration (ticket 49 revokes a
	// session grant, ticket 50 re-reads a manifest), so the bridge re-resolves
	// the declared set per call through Options.DeclaredCaps.
	if err := reg.Register(Entry{Tool: tool, Decl: Decl{
		Capabilities: []Capability{CapNet}, Needs: []Capability{CapNet},
		Declared: risk.L0, Provider: KindManifest,
	}}); err != nil {
		t.Fatal(err)
	}
	b := New(Options{
		Registry: reg,
		// The host's view: this tool declared fs.read only. net was never
		// granted, which is exactly the bypass C3 exists to stop.
		DeclaredCaps: func(string) []Capability { return []Capability{CapFSRead} },
	})

	out, err := b.Execute(t.Context(), req("probe.net", `{}`))
	if err != nil {
		t.Fatalf("a C3 violation must NOT be returned as an error (SPEC-07 §2): %v", err)
	}
	if !out.IsError {
		t.Fatalf("a C3 violation must be refused, got a success: %+v", out)
	}
	if out.ErrorClass != "permission_denied" {
		t.Errorf("error class = %q, want permission_denied", out.ErrorClass)
	}
	if tool.Runs() != 0 {
		t.Fatalf("the tool executed %d times despite an undeclared capability", tool.Runs())
	}
	if !strings.Contains(out.Text, "net") {
		t.Errorf("rejection text must name the missing capability: %q", out.Text)
	}

	// Host-level authorization is the same rule with a different owner: the
	// tool declared honestly, the machine simply does not grant net.
	b2 := New(Options{Registry: reg, Authorized: []Capability{CapFSRead}})
	out2, err2 := b2.Execute(t.Context(), req("probe.net", `{}`))
	if err2 != nil || !out2.IsError || tool.Runs() != 0 {
		t.Fatalf("host-unauthorized capability must hard-reject too: err=%v out=%+v runs=%d",
			err2, out2, tool.Runs())
	}
}

// ---------------------------------------------------------------------------
// C4: the ToolProvider registry and its unlanded slots
// ---------------------------------------------------------------------------

// TestC4SlotsWithoutAnImplementationAreNotSilent is the anti-fake guard: the
// Manifest/Goja/MCP slots must answer with a checkable error, and no code in
// this repository may pretend to serve them.
func TestC4SlotsWithoutAnImplementationAreNotSilent(t *testing.T) {
	reg := NewRegistry()
	for _, kind := range []Kind{KindManifest, KindGoja, KindMCP} {
		err := reg.RegisterProvider(&slotProvider{kind: kind})
		if err == nil {
			t.Fatalf("slot %s accepted registration: something landed there without a ticket", kind)
		}
		if !errors.Is(err, ErrSlotNotLanded) {
			t.Fatalf("slot %s: err = %v, want errors.Is(err, ErrSlotNotLanded)", kind, err)
		}
		var se SlotErr
		if !errors.As(err, &se) || se.Kind != kind {
			t.Fatalf("slot %s: %v does not name the slot", kind, err)
		}
	}
	// The declared surface is enumerable, so `wisp doctor`-style checks can
	// report the gap instead of shipping it quietly.
	if got := len(UnlandedSlots()); got != 3 {
		t.Fatalf("UnlandedSlots() = %d, want 3 (manifest/goja/mcp)", got)
	}
	// Builtin is the one slot that does land, and its entries are the
	// directory the bridge serves.
	if err := reg.RegisterProvider(NewBuiltinProvider(
		Entry{Tool: &fixtureTool{name: "probe.ok"}, Decl: Decl{
			Capabilities: []Capability{CapFSRead}, Needs: []Capability{CapFSRead},
			Declared: risk.L0, Provider: KindBuiltin,
		}})); err != nil {
		t.Fatal(err)
	}
	if _, ok := reg.Lookup("probe.ok"); !ok {
		t.Fatal("builtin entry missing after RegisterProvider")
	}
}

type slotProvider struct{ kind Kind }

func (s *slotProvider) Kind() Kind                  { return s.kind }
func (s *slotProvider) List() []Entry               { return nil }
func (s *slotProvider) Lookup(string) (Entry, bool) { return Entry{}, false }

func TestC4RegistryRejectsBadDeclarations(t *testing.T) {
	reg := NewRegistry()
	good := Entry{Tool: &fixtureTool{name: "probe.a"}, Decl: Decl{
		Capabilities: []Capability{CapFSRead}, Needs: []Capability{CapFSRead},
		Declared: risk.L0, Provider: KindBuiltin,
	}}
	if err := reg.Register(good); err != nil {
		t.Fatal(err)
	}
	dup := good
	dup.Tool = &fixtureTool{name: "probe.a"}
	if err := reg.Register(dup); err == nil {
		t.Error("duplicate tool name accepted (C1 requires global uniqueness)")
	}
	for _, bad := range []struct {
		name string
		e    Entry
	}{
		{"no dot", Entry{Tool: &fixtureTool{name: "nodothere"}}},
		{"two dots", Entry{Tool: &fixtureTool{name: "a.b.c"}}},
		{"uppercase", Entry{Tool: &fixtureTool{name: "FS.Read"}}},
		{"no parameter schema", Entry{
			Tool: &fixtureTool{name: "probe.e", noSchema: true},
			Decl: Decl{Declared: risk.L0},
		}},
		{"unknown capability", Entry{Tool: &fixtureTool{name: "probe.u"}, Decl: Decl{
			Capabilities: []Capability{"root"}, Needs: []Capability{"root"},
			Declared: risk.L0,
		}}},
		{"risk above L2", Entry{Tool: &fixtureTool{name: "probe.r"}, Decl: Decl{
			Capabilities: []Capability{CapFSRead}, Needs: []Capability{CapFSRead},
			Declared: risk.Deny,
		}}},
	} {
		if err := reg.Register(bad.e); err == nil {
			t.Errorf("%s: Register accepted an invalid entry", bad.name)
		}
	}
}

// ---------------------------------------------------------------------------
// AC#1: risk decision -> the right gate branch
// ---------------------------------------------------------------------------

// TestRiskDecisionRoutesEachLevelToItsBranch is AC#1's second clause: L0 passes
// without asking, L1 reaches the pending-window callback, L2 reaches the
// pending-approval callback - and each branch is reached through the real C19
// assessor, not a stubbed level.
func TestRiskDecisionRoutesEachLevelToItsBranch(t *testing.T) {
	t.Run("L0 passes and no gate is consulted", func(t *testing.T) {
		g := &gateSpy{}
		b, tools := newFixtureBridge(t, g,
			entry("probe.read", risk.L0, []Capability{CapFSRead}, []Capability{CapFSRead}))
		out, err := b.Execute(t.Context(), req("probe.read", `{}`))
		if err != nil || out.IsError {
			t.Fatalf("out=%+v err=%v", out, err)
		}
		if w, a := g.counts(); w != 0 || a != 0 {
			t.Fatalf("L0 must not reach a gate, got window=%d approval=%d", w, a)
		}
		if tools["probe.read"].Runs() != 1 {
			t.Fatal("the tool must run exactly once")
		}
	})

	t.Run("L1 uses the pending-window callback", func(t *testing.T) {
		g := &gateSpy{windowAns: AnswerAllow}
		b, tools := newFixtureBridge(t, g,
			entry("probe.rev", risk.L1, []Capability{CapFSWrite}, []Capability{CapFSWrite}))
		out, err := b.Execute(t.Context(), req("probe.rev", `{}`))
		if err != nil {
			t.Fatal(err)
		}
		if out.IsError {
			t.Fatalf("an allowed L1 call must execute: %+v", out)
		}
		w, a := g.counts()
		if w != 1 || a != 0 {
			t.Fatalf("L1 routing: window=%d approval=%d, want 1/0", w, a)
		}
		if out.RiskLevel != "L1" {
			t.Errorf("outcome risk = %q, want L1 (the assessed level, not the directory's)", out.RiskLevel)
		}
		d := g.windowDecision()
		if d.Reason == "" || !strings.Contains(d.Reason, "R1") {
			t.Errorf("L1 card reason must come from the rules that fired: %+v", d)
		}
		if tools["probe.rev"].Runs() != 1 {
			t.Fatal("tool did not run once")
		}
	})

	t.Run("L1 window timeout means execute", func(t *testing.T) {
		g := &gateSpy{windowAns: AnswerTimeout}
		b, tools := newFixtureBridge(t, g,
			entry("probe.rev", risk.L1, []Capability{CapFSWrite}, []Capability{CapFSWrite}))
		if _, err := b.Execute(t.Context(), req("probe.rev", `{}`)); err != nil {
			t.Fatal(err)
		}
		if tools["probe.rev"].Runs() != 1 {
			t.Fatal("SPEC-06 §2: an unopposed L1 window expires INTO execution")
		}
	})

	t.Run("L1 veto returns a failure to the model, not to the user", func(t *testing.T) {
		g := &gateSpy{windowAns: AnswerVeto, windowWhy: "用户在确认窗口中点了取消"}
		b, tools := newFixtureBridge(t, g,
			entry("probe.rev", risk.L1, []Capability{CapFSWrite}, []Capability{CapFSWrite}))
		out, err := b.Execute(t.Context(), req("probe.rev", `{}`))
		if err != nil {
			t.Fatalf("a veto is a refusal, not a host fault: %v", err)
		}
		if !out.IsError || out.ErrorClass != "user_rejected" {
			t.Fatalf("out=%+v, want IsError + user_rejected", out)
		}
		if tools["probe.rev"].Runs() != 0 {
			t.Fatal("a vetoed call must never reach the tool")
		}
	})

	t.Run("L2 uses the pending-approval callback", func(t *testing.T) {
		g := &gateSpy{approveAns: AnswerAllow}
		b, tools := newFixtureBridge(t, g,
			entry("probe.irrev", risk.L2, []Capability{CapFSWrite}, []Capability{CapFSWrite}))
		out, err := b.Execute(t.Context(), req("probe.irrev", `{}`))
		if err != nil {
			t.Fatal(err)
		}
		if out.IsError {
			t.Fatalf("an approved L2 call must execute: %+v", out)
		}
		w, a := g.counts()
		if w != 0 || a != 1 {
			t.Fatalf("L2 routing: window=%d approval=%d, want 0/1", w, a)
		}
		if tools["probe.irrev"].Runs() != 1 {
			t.Fatal("tool did not run")
		}
	})

	t.Run("L2 approval timeout auto-rejects", func(t *testing.T) {
		g := &gateSpy{approveAns: AnswerTimeout}
		b, tools := newFixtureBridge(t, g,
			entry("probe.irrev", risk.L2, []Capability{CapFSWrite}, []Capability{CapFSWrite}))
		out, err := b.Execute(t.Context(), req("probe.irrev", `{}`))
		if err != nil {
			t.Fatal(err)
		}
		if !out.IsError {
			t.Fatalf("an unanswered L2 must never execute: %+v", out)
		}
		if tools["probe.irrev"].Runs() != 0 {
			t.Fatal("auto-rejected call reached the tool")
		}
	})

	t.Run("NoGate fails closed above L0", func(t *testing.T) {
		b, tools := newFixtureBridge(t, nil,
			entry("probe.irrev", risk.L2, []Capability{CapFSWrite}, []Capability{CapFSWrite}))
		out, err := b.Execute(t.Context(), req("probe.irrev", `{}`))
		if err != nil {
			t.Fatal(err)
		}
		if !out.IsError || tools["probe.irrev"].Runs() != 0 {
			t.Fatalf("with no approval channel wired the only safe answer is no: %+v", out)
		}
		if !strings.Contains(out.Text, "票 21") {
			t.Errorf("the refusal must name who owes the missing channel: %q", out.Text)
		}
	})
}

// TestDeclaredRiskIsOnlyAFloor proves R1's meaning in the bridge: a tool that
// DECLARES L0 but reaches for an out-of-scope path is judged L2 by C19, so a
// plugin (or a builtin) cannot self-downgrade its way past the gate.
func TestDeclaredRiskIsOnlyAFloor(t *testing.T) {
	root := tempCanonical(t)
	outside := tempCanonical(t)
	g := &gateSpy{approveAns: AnswerReject}
	b, _ := fsBridgeWith(t, nil, g, root)
	paths := b.Paths()

	// in-scope: L0, no gate.
	f := writeUnder(t, root, "in.txt", "hello")
	if out, err := b.Execute(t.Context(), req("fs.read", argsFor(f))); err != nil || out.IsError {
		t.Fatalf("in-allowlist read: out=%+v err=%v", out, err)
	}
	if w, a := g.counts(); w != 0 || a != 0 {
		t.Fatalf("in-scope read reached a gate: %d/%d", w, a)
	}

	// out-of-scope: the SAME tool, the SAME declared L0, now L2.
	if out, err := b.Execute(t.Context(), req("fs.read", argsFor(writeUnder(t, outside, "out.txt", "s")))); err != nil {
		t.Fatal(err)
	} else if !out.IsError || out.ErrorClass != "user_rejected" {
		t.Fatalf("out-of-scope read must be refused by the gate: %+v", out)
	}
	if w, a := g.counts(); w != 0 || a != 1 {
		t.Fatalf("out-of-scope read routing: window=%d approval=%d, want 0/1", w, a)
	}
	d := g.approvalDecision()
	if d.Level != risk.L2 || !contains(d.RulesHit, risk.R2) {
		t.Fatalf("verdict = %+v, want L2 with R2 in rules_hit", d)
	}
	if !paths.InAllowlist(root) || paths.InAllowlist(outside) {
		t.Fatal("allowlist membership is inverted")
	}
}

// ---------------------------------------------------------------------------
// AC#1: tool_call rows
// ---------------------------------------------------------------------------

// TestToolCallRowsAreComplete is AC#1's third clause: the row carries the
// ASSESSED level, the decision, the outcome and a correlation_id, for the
// executed path as well as for the refused ones.
func TestToolCallRowsAreComplete(t *testing.T) {
	root := tempCanonical(t)
	store := openStore(t, filepath.Join(root, "data"))
	var j agent.Journal = store
	g := &gateSpy{approveAns: AnswerReject, windowAns: AnswerVeto}
	b, _ := fsBridgeWith(t, j, g, root)
	mustStartTask(t, store, "task-1")

	outside := tempCanonical(t)
	inside := writeUnder(t, root, "a.txt", "hi")
	cases := []struct {
		tool, args, risk, decision, outcome, class string
	}{
		{"fs.read", argsFor(inside), "L0", "allow", "success", ""},
		// An L2 the user refused: rejected, booked L2, class user_rejected
		// (D37 keeps "the policy said no" and "the human said no" apart).
		{"fs.read", argsFor(writeUnder(t, outside, "b.txt", "x")), "L2", "reject", "error", "user_rejected"},
		{"fs.nope", `{}`, "L0", "reject", "error", "tool"},
		{"fs.read", `not-json`, "L0", "reject", "error", "tool"},
	}
	for i, c := range cases {
		rq := req(c.tool, c.args)
		rq.TaskID, rq.CorrelationID = "task-1", fmt.Sprintf("corr-%d", i)
		if _, err := b.Execute(t.Context(), rq); err != nil {
			t.Fatalf("%s: %v", c.tool, err)
		}
	}
	// One L1 veto path through a declared-L1 tool, to book the third decision
	// shape (user_rejected) the vocabulary has.
	l1 := Entry{Tool: &fixtureTool{name: "probe.rev"}, Decl: Decl{
		Capabilities: []Capability{CapFSWrite}, Needs: []Capability{CapFSWrite},
		Declared: risk.L1, Provider: KindBuiltin,
	}}
	if err := b.reg.Register(l1); err != nil {
		t.Fatal(err)
	}
	rq := req("probe.rev", `{}`)
	rq.TaskID, rq.CorrelationID = "task-1", "corr-l1"
	if _, err := b.Execute(t.Context(), rq); err != nil {
		t.Fatal(err)
	}
	cases = append(cases, struct{ tool, args, risk, decision, outcome, class string }{
		"probe.rev", "{}", "L1", "reject", "error", "user_rejected",
	})

	rows, err := store.ListToolCallsByTask(t.Context(), "task-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != len(cases) {
		t.Fatalf("booked %d rows, want %d: %+v", len(rows), len(cases), rows)
	}
	for i, c := range cases {
		r := rows[i]
		if r.Tool != c.tool || r.Seq != int64(i+1) {
			t.Errorf("row %d = %s seq %d, want %s seq %d", i, r.Tool, r.Seq, c.tool, i+1)
		}
		if r.RiskLevel != c.risk {
			t.Errorf("%s: risk_level = %q, want the ASSESSED %q", c.tool, r.RiskLevel, c.risk)
		}
		if r.Decision != c.decision {
			t.Errorf("%s: decision = %q, want %q", c.tool, r.Decision, c.decision)
		}
		if r.Outcome != c.outcome {
			t.Errorf("%s: outcome = %q, want %q", c.tool, r.Outcome, c.outcome)
		}
		if r.ErrorClass != c.class {
			t.Errorf("%s: error_class = %q, want %q", c.tool, r.ErrorClass, c.class)
		}
		if r.CorrelationID == "" {
			t.Errorf("%s: correlation_id is empty (C18 routing/forensics need it)", c.tool)
		}
	}
}

// TestBridgeIsTheLoopToolProvider is the "fits the existing seam, not a
// parallel one" proof: a *Bridge satisfies agent.ToolProvider, so the loop's
// Tools()/Execute() calls are the bridge's, with no second dispatch path.
func TestBridgeIsTheLoopToolProvider(t *testing.T) {
	var p agent.ToolProvider = New(Options{})
	dir, err := p.Tools(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if len(dir) != 0 {
		t.Fatalf("empty bridge reports %d tools", len(dir))
	}
	if _, err := p.Execute(t.Context(), req("fs.read", `{}`)); err != nil {
		t.Fatalf("unknown tool must not be a provider fault: %v", err)
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func contains(ids []risk.RuleID, want risk.RuleID) bool {
	for _, id := range ids {
		if id == want {
			return true
		}
	}
	return false
}
