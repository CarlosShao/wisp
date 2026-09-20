package tools_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// THE COMPOSITION TEST (rulings A8/A11/A13/R12 shape).
//
// Segment 1's answer to "who calls this in production" was "nobody": cmd/wisp
// builds no loop and no gate. This file does not change that answer - it
// changes what is missing. Everything the fs write family needs from the
// approval layer is now composed HERE, in one place, with the REAL
// approval.Gate: the L1 block window, the veto channels, and the
// applied-steps report type that D31 renders. If any of those three were a
// comment instead of a call, these tests would not compile, let alone pass.
//
// It lives in package tools_test because internal/agent/approval imports
// internal/tools: the dependency direction is a fact, and a test that imports
// it from the other side is where a cycle would be born.

// uiSpy is the injected presentation surface: it records every card the gate
// displayed, which is how "the window really opened" is observed.
type uiSpy struct {
	mu      sync.Mutex
	prompts []approval.Prompt
	events  []approval.Event
	// onPrompt runs while the window is open (a veto can be delivered from here).
	onPrompt func(approval.Prompt)
}

func (u *uiSpy) Prompt(_ context.Context, p approval.Prompt) error {
	u.mu.Lock()
	u.prompts = append(u.prompts, p)
	hook := u.onPrompt
	u.mu.Unlock()
	if hook != nil {
		hook(p)
	}
	return nil
}

func (u *uiSpy) Update(_ context.Context, e approval.Event) error {
	u.mu.Lock()
	u.events = append(u.events, e)
	u.mu.Unlock()
	return nil
}

func (u *uiSpy) cards() []approval.Prompt {
	u.mu.Lock()
	defer u.mu.Unlock()
	return append([]approval.Prompt(nil), u.prompts...)
}

func (u *uiSpy) started() bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	for _, e := range u.events {
		if e.Kind == approval.EventStarted {
			return true
		}
	}
	return false
}

// compose wires the fs family over a real approval gate and its cancel bus.
func compose(t *testing.T, allowed string, deps tools.FSDeps) (*tools.Bridge, *approval.Gate, *uiSpy) {
	t.Helper()
	ui := &uiSpy{}
	g := approval.New(approval.Options{UI: ui})
	deps.Paths = tools.NewPathCanonicalizer([]string{allowed}, nil)
	reg := tools.NewRegistry()
	for _, e := range tools.BuiltinFSEntries(deps) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	b := tools.New(tools.Options{
		Registry:   reg,
		Paths:      deps.Paths,
		Gate:       g,
		Cancel:     g.ToolsCancelBus(),
		Provenance: risk.NewProvenance(risk.ProvOptions{NoProbe: true}),
		Logf:       func(string, ...any) {},
	})
	// D47: only a task the TEXT loop registered may reach a gate at all.
	t.Cleanup(g.AdmitTextTask("task-wire-1"))
	return b, g, ui
}

func call(b *tools.Bridge, t *testing.T, name, args string) (agent.ToolOutcome, error) {
	t.Helper()
	return b.Execute(t.Context(), agentReq(name, args))
}

// TestL1WriteGoesThroughTheRealBlockWindow is the headline: an fs.write that
// declares L1 is NOT refused for lack of a gate and NOT executed without one.
// The window opens, nobody vetoes, the window expires, and the file lands.
func TestL1WriteGoesThroughTheRealBlockWindow(t *testing.T) {
	dir := t.TempDir()
	b, _, ui := compose(t, mustCanon(t, dir), tools.FSDeps{})
	target := filepath.Join(dir, "written-through-the-window.txt")

	start := time.Now()
	out, err := call(b, t, "fs.write", fmt.Sprintf(
		`{"path":%q,"content":"hello from the L1 window"}`, filepath.ToSlash(target)))
	if err != nil {
		t.Fatal(err)
	}
	if out.IsError {
		t.Fatalf("the window ran out unopposed, which MEANS execute: %+v", out)
	}
	cards := ui.cards()
	if len(cards) != 1 {
		t.Fatalf("the gate displayed %d cards, want exactly the one L1 window", len(cards))
	}
	if cards[0].Level != "L1" || cards[0].Tool != "fs.write" {
		t.Errorf("card = %s %s, want L1 fs.write", cards[0].Level, cards[0].Tool)
	}
	if cards[0].Window < 2*time.Second {
		t.Errorf("the card advertised a %v window, below the 2s floor (SPEC-06 §2)", cards[0].Window)
	}
	if !ui.started() {
		t.Error("the gate never reported the handoff to execution (EventStarted)")
	}
	body, err := os.ReadFile(target)
	if err != nil || string(body) != "hello from the L1 window" {
		t.Fatalf("the file did not land: %v %q", err, body)
	}
	if time.Since(start) < 2*time.Second {
		t.Errorf("the call returned in %v: an unvetoed L1 window must BLOCK, not pass through",
			time.Since(start))
	}
}

// TestVetoInsideTheWindowWritesNothing is the same route with a user who said
// stop in time: nothing is created, and the refusal travels as a refusal.
func TestVetoInsideTheWindowWritesNothing(t *testing.T) {
	dir := t.TempDir()
	b, g, ui := compose(t, mustCanon(t, dir), tools.FSDeps{
		Hooks: tools.Hooks{},
	})
	ui.onPrompt = func(p approval.Prompt) {
		if err := g.Veto(approval.Veto{CorrelationID: p.CorrelationID, Channel: approval.ChannelBall}); err != nil {
			t.Errorf("veto: %v", err)
		}
	}
	target := filepath.Join(dir, "never-written.txt")
	out, err := call(b, t, "fs.write", fmt.Sprintf(`{"path":%q,"content":"x"}`, filepath.ToSlash(target)))
	if err != nil {
		t.Fatal(err)
	}
	if !out.IsError {
		t.Fatalf("a vetoed L1 call must not execute: %+v", out)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("the veto came before the window ended, yet the file exists: %v", err)
	}
	if ui.started() {
		t.Error("the gate reported a handoff to execution for a call it vetoed")
	}
}

// TestLateVetoRendersTheApprovalLayersAppliedStepsReport is D31's cancel
// semantics on the REAL report type: the veto arrives after the window handed
// off, so the answer must be approval.CancellationReport's wording - including
// the veto channel and the per-step 「已执行」 lines - and not a second shape
// the bridge invented on the way out.
func TestLateVetoRendersTheApprovalLayersAppliedStepsReport(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "half-written.txt")

	var mu sync.Mutex
	handoff := make(chan struct{})
	// Declared first: the AtStep closure below reaches g, and a :=-declared
	// variable is not in scope inside its own initializer.
	var b *tools.Bridge
	var g *approval.Gate
	var ui *uiSpy
	b, g, ui = compose(t, mustCanon(t, dir), tools.FSDeps{
		WriteChunk: 8,
		Hooks: tools.Hooks{AtStep: func(step string) {
			if step != "write:24" {
				return
			}
			mu.Lock()
			defer mu.Unlock()
			// The user hits 「取消」 while the write is halfway through staging:
			// past the window, so this is a stop request and not a rewind.
			if err := g.Veto(approval.Veto{
				CorrelationID: "corr-wire-1", Channel: approval.ChannelEsc,
			}); err == nil {
				t.Errorf("a veto after handoff must be reported as non-atomic (ErrAlreadyStarted), got nil")
			}
			select {
			case <-handoff:
			default:
				close(handoff)
			}
		}},
	})
	_ = ui

	out, err := call(b, t, "fs.write", fmt.Sprintf(
		`{"path":%q,"content":%q}`, filepath.ToSlash(target), strings.Repeat("NEWB", 40)))
	if err != nil {
		t.Fatal(err)
	}
	if !out.IsError {
		t.Fatalf("the call stopped mid-flight and must say so: %+v", out)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("D31: a file appeared at the target despite the stop before the rename: %v", err)
	}
	for _, want := range []string{
		"取消不是原子的", // the D31 sentence, from approval.CancellationReport
		"否决通道",    // the channel attribution only that type renders
		"已执行：",    // one line per step that landed
	} {
		if !strings.Contains(out.Text, want) {
			t.Errorf("the report is missing %q:\n%s", want, out.Text)
		}
	}
	if len(out.AppliedSteps) == 0 {
		t.Error("the outcome carries no ledger even though the tool applied steps")
	}
	for _, st := range out.AppliedSteps {
		if !strings.Contains(out.Text, st) {
			t.Errorf("applied step %q is not in the report the user sees:\n%s", st, out.Text)
		}
	}
	// And the same steps the report lists are the ones the outcome carries: the
	// bridge forwards approval's report, it does not rebuild one.
	joined := strings.Join(out.AppliedSteps, "\n")
	if !strings.Contains(joined, "创建临时文件") || !strings.Contains(joined, "停止") {
		t.Errorf("the ledger should name the staging file and the stop point:\n%s", joined)
	}
}

// TestApprovedL2ThatStopsMidWriteIsNotRenderedAsNeverStarted covers the L2
// route's half of the same rule: an approved overwrite that then fails or is
// cancelled has started, and the report must not say otherwise.
func TestApprovedL2ThatStopsMidWriteIsNotRenderedAsNeverStarted(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "overwrite.txt")
	if err := os.WriteFile(target, []byte("ORIGINAL"), 0o600); err != nil {
		t.Fatal(err)
	}
	b, g, ui := compose(t, mustCanon(t, dir), tools.FSDeps{
		WriteChunk: 8,
		Hooks: tools.Hooks{Kill: func(step string) error {
			if step == "write:24" {
				return fmt.Errorf("模拟进程在此刻被杀")
			}
			return nil
		}},
	})
	ui.onPrompt = func(p approval.Prompt) {
		if p.Grant == "" {
			return
		}
		if err := g.Native().Allow(t.Context(), p.CorrelationID, p.Grant); err != nil {
			t.Errorf("native allow: %v", err)
		}
	}
	out, err := call(b, t, "fs.write", fmt.Sprintf(
		`{"path":%q,"content":%q}`, filepath.ToSlash(target), strings.Repeat("NEWB", 40)))
	if err != nil {
		t.Fatal(err)
	}
	if !out.IsError {
		t.Fatalf("a killed write is a failure: %+v", out)
	}
	if got, err := os.ReadFile(target); err != nil || string(got) != "ORIGINAL" {
		t.Fatalf("the overwrite must not have landed: %q err=%v", got, err)
	}
	if strings.Contains(out.Text, "本次调用未进入执行阶段") {
		t.Errorf("the report claims nothing ran over a ledger of applied steps:\n%s", out.Text)
	}
	if !strings.Contains(out.Text, "以下步骤已生效") {
		t.Errorf("the report must list what landed:\n%s", out.Text)
	}
}

// TestFSReadOnlyNeverOpensACard keeps the cheap end cheap: an L0 read through
// the SAME composition must not block on a window.
func TestFSReadOnlyNeverOpensACard(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "note.txt")
	if err := os.WriteFile(target, []byte("quiet"), 0o600); err != nil {
		t.Fatal(err)
	}
	b, _, ui := compose(t, mustCanon(t, dir), tools.FSDeps{})
	out, err := call(b, t, "fs.read", fmt.Sprintf(`{"path":%q}`, filepath.ToSlash(target)))
	if err != nil || out.IsError || out.Text != "quiet" {
		t.Fatalf("out=%+v err=%v", out, err)
	}
	if len(ui.cards()) != 0 {
		t.Errorf("an L0 read displayed %d cards: the L0 route must not reach a gate", len(ui.cards()))
	}
}

// helpers ---------------------------------------------------------------

func mustCanon(t *testing.T, p string) string {
	t.Helper()
	c, err := tools.NewPathCanonicalizer(nil, nil).Canonicalize(p)
	if err != nil {
		t.Fatalf("Canonicalize(%q): %v", p, err)
	}
	return c
}

func agentReq(name, args string) agent.ToolRequest {
	return agent.ToolRequest{
		TaskID: "task-wire-1", CorrelationID: "corr-wire-1", CallID: "call-1",
		Name: name, Args: json.RawMessage(args),
	}
}
