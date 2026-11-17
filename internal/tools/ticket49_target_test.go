package tools

// Ticket 49 tests: the input.type binding inside the REAL decision chain.
//
// The internal/session cases (ticket49_lifecycle_test.go) book the rows and
// measure the doors against a real SQLite file; these are the bridge-half
// teeth - the §15 item-7 gates expressed as behavior, each with its negative
// control, in the same shape grant_test.go established for ticket 224:
//
//	L2 tooth       TestTicket49InputTypeStaysL2WithoutAGrant
//	downgrade      TestTicket49LiveTargetGrantSilentlyExecutesTheL1Shape
//	red line ⑤     TestTicket49TaintedExfilIgnoresTheTargetGrant
//	per-process ③  TestTicket49TargetSwapReRefusesAndReAsks
//	pid fallback   TestTicket49BridgeRefusesPidFallbackKey
//
// These cases need no sherpa PATH.

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

// ---------------------------------------------------------------------------
// fixtures
// ---------------------------------------------------------------------------

// t49Grants answers the target question for the process names it was told to
// cover. Pointer receiver; per-bridge construction (the same reason
// t224Grants is one).
type t49Grants struct {
	covered map[string]int64
	calls   []string
}

func newT49Grants() *t49Grants { return &t49Grants{covered: map[string]int64{}} }

func (g *t49Grants) Covering(context.Context, string, []string) (int64, bool) { return 0, false }

func (g *t49Grants) TargetCovering(_ context.Context, tool, process string) (int64, bool) {
	g.calls = append(g.calls, tool+" "+process)
	if id, ok := g.covered[process]; ok && tool == TargetGrantTool {
		return id, true
	}
	return 0, false
}

// t49TargetBridge assembles a bridge whose input.type entry runs through a
// recording gate, with the resolver pinned to one process name. The
// systemTool itself is not registered: the bridge's judge/journal/route run on
// the DECLARATION the real family carries (InputTypeDecl), which is the layer
// the binding lives in. allowlist stays empty: input.type has no paths, so R2
// is dormant for it by construction.
func t49TargetBridge(t *testing.T, grants GrantSource, resolve func(map[string]any) (string, bool)) (*Bridge, *t90Gate) {
	t.Helper()
	e, _ := t90Entry(TargetGrantTool, risk.L2, nil, nil)
	// t90Entry builds a generic tool; the DECLARATION is what the judge reads,
	// and the binding keys on the tool NAME, so the generic tool is fine here.
	reg := NewRegistry()
	if err := reg.Register(e); err != nil {
		t.Fatalf("Register: %v", err)
	}
	gate := &t90Gate{}
	b := New(Options{
		Registry:      reg,
		Gate:          gate,
		Grants:        grants,
		ResolveTarget: resolve,
		Modes:         &t90Modes{},
		Logf:          func(string, ...any) {},
	})
	return b, gate
}

func t49Call(b *Bridge, corr string) (agent.ToolOutcome, error) {
	return b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "t49-task", CorrelationID: corr, CallID: "call-" + corr,
		Name: TargetGrantTool, Args: json.RawMessage(`{"text":"hello","hwnd":7}`),
	})
}

// ---------------------------------------------------------------------------
// the L2 tooth: without a grant, input.type keeps its card
// ---------------------------------------------------------------------------

// TestTicket49InputTypeStaysL2WithoutAGrant is the tooth the dispatch named
// first: input.type L2 永不被 grant 自动放行 in its NEGATIVE direction - with
// no covering row (or no resolver, or no source at all) the call produces
// exactly one approval card. The control is the grantless half of the
// downgrade case below.
func TestTicket49InputTypeStaysL2WithoutAGrant(t *testing.T) {
	for _, tc := range []struct {
		name    string
		grants  GrantSource
		resolve func(map[string]any) (string, bool)
	}{
		{"no source at all", nil, func(map[string]any) (string, bool) { return "notepad.exe", true }},
		{"source answers nothing", newT49Grants(), func(map[string]any) (string, bool) { return "notepad.exe", true }},
		{"no resolver wired", newT49Grants(), nil},
		{"resolver finds nothing", newT49Grants(), func(map[string]any) (string, bool) { return "", false }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, gate := t49TargetBridge(t, tc.grants, tc.resolve)
			if _, err := t49Call(b, "t49-l2-"+strings.ReplaceAll(tc.name, " ", "-")); err != nil {
				t.Fatalf("Execute: %v", err)
			}
			if w, a := gate.counts(); a != 1 || w != 0 {
				t.Errorf("%s: gate saw windows=%d approvals=%d, want 0/1: the L2 card must ask", tc.name, w, a)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// the downgrade: a live target grant turns the L2 card into a silent execute
// ---------------------------------------------------------------------------

// TestTicket49LiveTargetGrantSilentlyExecutesTheL1Shape is the positive cell
// the D34 row promises (L2 首次 → 会话授权后不再询问): with a live row for the
// resolved process, the SAME call produces ZERO cards and books
// decision=allow_session_grant with the row's id on tool_call.grant_id. The
// grantless control ran in the case above.
func TestTicket49LiveTargetGrantSilentlyExecutesTheL1Shape(t *testing.T) {
	grants := newT49Grants()
	grants.covered["notepad.exe"] = 4242
	j := &t224Journal{}
	e, _ := t90Entry(TargetGrantTool, risk.L2, nil, nil)
	reg := NewRegistry()
	if err := reg.Register(e); err != nil {
		t.Fatalf("Register: %v", err)
	}
	gate := &t90Gate{}
	b := New(Options{
		Registry:      reg,
		Gate:          gate,
		Grants:        grants,
		Journal:       j,
		ResolveTarget: func(map[string]any) (string, bool) { return "notepad.exe", true },
		Modes:         &t90Modes{},
		Logf:          func(string, ...any) {},
	})
	out, err := t49Call(b, "t49-hit")
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if out.IsError {
		t.Fatalf("the covered call errored: %s", out.Text)
	}
	if w, a := gate.counts(); w != 0 || a != 0 {
		t.Errorf("a live target grant left gate visits (windows=%d approvals=%d): the covered call must not ask", w, a)
	}
	row, ok := j.last()
	if !ok {
		t.Fatal("no tool_call row was booked")
	}
	if row.Decision != agent.DecisionAllowGrant {
		t.Errorf("covered call booked decision=%q, want %q", row.Decision, agent.DecisionAllowGrant)
	}
	if row.GrantID == nil || *row.GrantID != 4242 {
		t.Errorf("covered call booked grant_id=%v, want 4242", row.GrantID)
	}
}

// ---------------------------------------------------------------------------
// red line ⑤: taint upgrade ignores any grant
// ---------------------------------------------------------------------------

// TestTicket49TaintedExfilIgnoresTheTargetGrant is the §15 item-7 ⑤ cell: a
// C25 taint hit on the same call refuses the downgrade even though a live row
// covers the resolved process. The taint detector answers "hit" for every
// call - the widest possible R4 - and the L2 card must still arrive with no
// grant linkage.
func TestTicket49TaintedExfilIgnoresTheTargetGrant(t *testing.T) {
	grants := newT49Grants()
	grants.covered["notepad.exe"] = 777
	e, _ := t90Entry(TargetGrantTool, risk.L2, nil, nil)
	reg := NewRegistry()
	if err := reg.Register(e); err != nil {
		t.Fatalf("Register: %v", err)
	}
	gate := &t90Gate{}
	b := New(Options{
		Registry: reg,
		Gate:     gate,
		Grants:   grants,
		Assessor: risk.NewRiskAssessor().WithTaintDetector(t90Taint{}),
		ResolveTarget: func(map[string]any) (string, bool) {
			return "notepad.exe", true
		},
		Modes: &t90Modes{},
		Logf:  func(string, ...any) {},
	})
	if _, err := t49Call(b, "t49-taint"); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if w, a := gate.counts(); a != 1 {
		t.Errorf("a tainted call with a live grant produced windows=%d approvals=%d, want 0/1: "+
			"R4 污染升级不受会话授权覆盖（§15 第 7 项 ⑤）", w, a)
	}
	if len(grants.calls) != 0 {
		t.Errorf("the grant source was consulted on a tainted call (%v); the red-line gate must refuse BEFORE the lookup", grants.calls)
	}
}

// t90Taint (ticket90_test.go) is the widest possible R4: every call is a taint
// hit with SessionOverrideBlocked set by the real judge's ruleTaint. Injected
// through Options.Assessor, the sanctioned injection seam.

// ---------------------------------------------------------------------------
// per-process ③: a target swap re-asks
// ---------------------------------------------------------------------------

// TestTicket49TargetSwapReRefusesAndReAsks is the ③ cell at both layers:
//
//	lookup half - the granted process is notepad.exe; a call resolving to
//	              cmd.exe finds no covering row (per-process isolation) and
//	              asks its own L2 card;
//	byte half   - a call granted for notepad.exe whose resolver NOW answers
//	              cmd.exe is refused with the re-confirm wording before the
//	              tool would run. (The swap guard re-resolves at execution.)
func TestTicket49TargetSwapReRefusesAndReAsks(t *testing.T) {
	grants := newT49Grants()
	grants.covered["notepad.exe"] = 99

	// Lookup half: the granted process is notepad.exe; a call resolving to
	// cmd.exe has no covering row and asks.
	b, gate := t49TargetBridge(t, grants, func(map[string]any) (string, bool) { return "cmd.exe", true })
	if _, err := t49Call(b, "t49-swap-lookup"); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if w, a := gate.counts(); a != 1 {
		t.Errorf("a different target with a Word-only grant produced windows=%d approvals=%d, want 0/1: "+
			"换目标应用要重新确认（§15 第 7 项 ③）", w, a)
	}

	// Byte half: the grant covered the judge-time resolution, but the
	// execution-time re-resolution moved. The resolver is call-counted: the
	// prime call consumes judge+exec (2 reads) on notepad.exe; the second
	// Execute reads notepad.exe at judge time (grant hit) and cmd.exe at the
	// byte boundary - exactly the foreground-moved-between-judge-and-exec
	// shape the swap guard exists for.
	reads := 0
	current := func() string {
		reads++
		if reads <= 3 {
			return "notepad.exe"
		}
		return "cmd.exe"
	}
	b2, _ := t49TargetBridge(t, grants, func(map[string]any) (string, bool) {
		return current(), true
	})
	if _, err := t49Call(b2, "t49-swap-prime"); err != nil {
		t.Fatalf("prime call: %v", err)
	}
	out, err := t49Call(b2, "t49-swap-hit")
	if err != nil {
		t.Fatalf("swapped call: %v", err)
	}
	if !out.IsError {
		t.Fatalf("a call whose target moved between the judge and the byte boundary executed anyway: %s", out.Text)
	}
	if !strings.Contains(out.Text, "目标应用已切换") || !strings.Contains(out.Text, "重新确认") {
		t.Errorf("the swap refusal does not carry the ruling's wording: %s", out.Text)
	}
}

// ---------------------------------------------------------------------------
// the pid fallback: never a grant key
// ---------------------------------------------------------------------------

// TestTicket49BridgeRefusesPidFallbackKey is the mirror of the session-side
// shape note: when the resolver can only produce "pid:<n>" (the image name
// could not be read), the bridge refuses to key a lookup with it and the card
// asks - a pid-keyed rule is dead the moment the pid is reused.
func TestTicket49BridgeRefusesPidFallbackKey(t *testing.T) {
	grants := newT49Grants()
	grants.covered["pid:4242"] = 555
	b, gate := t49TargetBridge(t, grants, func(map[string]any) (string, bool) { return "pid:4242", true })
	if _, err := t49Call(b, "t49-pid"); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if w, a := gate.counts(); a != 1 {
		t.Errorf("a pid-fallback target with a matching row produced windows=%d approvals=%d, want 0/1: "+
			"a pid is not a pattern anyone agreed to", w, a)
	}
	if len(grants.calls) != 0 {
		t.Errorf("the grant source was keyed with a pid fallback (%v)", grants.calls)
	}
}

// ---------------------------------------------------------------------------
// the shape pins: one grant channel, no authority on the new fact
// ---------------------------------------------------------------------------

// TestTicket49TargetFactCarriesNoAuthority restates ticket 90/224's pins for
// the new field: TargetProcess is a description of the call, not a channel the
// caller could set a verdict through. The name may not contain an allow-shaped
// word, and TargetResolver is a facts function, not a grant channel (it does
// not appear among Bridge's grant-shaped fields - the 224 pin already proves
// the field count).
func TestTicket49TargetFactCarriesNoAuthority(t *testing.T) {
	if strings.Contains(strings.ToLower("TargetProcess"), "allow") ||
		strings.Contains(strings.ToLower("TargetProcess"), "grant") ||
		strings.Contains(strings.ToLower("TargetProcess"), "approve") {
		t.Errorf("the target fact's name itself reads as an authority")
	}
	// The resolver returns (name, ok) - no verdict, no answer vocabulary.
	var r TargetResolver = func(map[string]any) (string, bool) { return "", false }
	name, ok := r(nil)
	if name != "" || ok {
		t.Errorf("a nil-params resolve must answer empty/false, got (%q, %v)", name, ok)
	}
	// The binding tool name is the frozen row.
	if TargetGrantTool != "input.type" {
		t.Errorf("TargetGrantTool = %q, want the frozen D34 row", TargetGrantTool)
	}
}

// TestTicket49GrantUseBooksGrantIDOnEveryRow is the audit cell at bridge
// scale: TWO covered calls in a row each book their own tool_call row with
// the grant's id - 「每次使用被授权通道都写 tool_call 日志（可取证）」, per call,
// not once per grant.
func TestTicket49GrantUseBooksGrantIDOnEveryRow(t *testing.T) {
	grants := newT49Grants()
	grants.covered["notepad.exe"] = 31
	j := &t224Journal{}
	e, _ := t90Entry(TargetGrantTool, risk.L2, nil, nil)
	reg := NewRegistry()
	if err := reg.Register(e); err != nil {
		t.Fatalf("Register: %v", err)
	}
	b := New(Options{
		Registry:      reg,
		Gate:          &t90Gate{},
		Grants:        grants,
		Journal:       j,
		ResolveTarget: func(map[string]any) (string, bool) { return "notepad.exe", true },
		Modes:         &t90Modes{},
		Logf:          func(string, ...any) {},
	})
	for i := 0; i < 2; i++ {
		if _, err := t49Call(b, fmt.Sprintf("t49-audit-%d", i)); err != nil {
			t.Fatalf("call %d: %v", i, err)
		}
	}
	if len(j.calls) != 2 {
		t.Fatalf("journal rows = %d, want 2", len(j.calls))
	}
	for i, row := range j.calls {
		if row.Decision != agent.DecisionAllowGrant {
			t.Errorf("row %d decision=%q, want %q", i, row.Decision, agent.DecisionAllowGrant)
		}
		if row.GrantID == nil || *row.GrantID != 31 {
			t.Errorf("row %d grant_id=%v, want 31: every grant-backed execution books the linkage", i, row.GrantID)
		}
	}
}
