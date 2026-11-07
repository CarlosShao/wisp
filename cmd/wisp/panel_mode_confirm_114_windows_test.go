//go:build windows

package main

// Ticket 114 AC#3 (R-92-2's other half): the panel's mode chain, driven through
// the PRODUCTION doors of both legs, before and after the L2 confirmation leg was
// attached (ticket 114-r2).
//
// WHAT THESE CASES ARE FOR, in the ticket's own words. AC#2 put a gate in front of
// the panel's mode write that refuses a loosening request when the injected Confirm
// is nil. That gate was correct and unreachable: both production assemblies passed
// nil, so "no L2 leg" was a constant and the panel could never ask for a looser档
// at all. 114-r2 made the leg a parameter of the shared chain and the resident boot
// hands it the orb's own approval gate. These cases pin the readings that change
// has to buy, and the one reading that must NOT change:
//
//  1. the resident leg (panel_host_windows.go's dispatchRaw, i.e. the page's door):
//     ask_every_step -> ask_high_risk is a LOOSENING, so with no leg attached it is
//     refused at the panel gate and nothing lands (R9, owner 2026-11-05, ledger
//     A1541 item 4); with the leg attached it reaches the write leg and asks
//     exactly ONE card. Ticket 372 had made that row land with zero cards; R9
//     restores the card and keeps 372's read-side downgrade untouched.
//  2. the same leg, same door: auto_approve asks its own card, and it is the
//     SECOND card of the pair - one per loosening, never one per process. That
//     pair is the proof that wiring the leg did not remove M4's card; run_mode101
//     keeps the same pair for the run seam.
//  3. the CLI pipe (panel_inbound.go's chain): a loosening is refused with
//     ErrNoL2Confirm in EVERY tier now, because the pipe has nowhere to show a
//     card. That is R9's knowingly accepted cost (docs/evidence/s1/R9-impl.md
//     §4): the pipe can only tighten. What it must never do is widen what the
//     process enforces - its Store has no leg, so its effective read stays
//     ask_every_step even when the file says ask_high_risk (372-r2's read-side
//     downgrade), which the second subtest still measures.

import (
	"context"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/perm"
	"github.com/CarlosShao/wisp/internal/risk"
)

// auditQuiet114r2 is the sink both assemblies take. Silent on purpose: these cases
// read the disk and the leg's call count, and a case that asserted on prose would
// pass for a leg that wrote the right line and did nothing.
func auditQuiet114r2(string, ...any) {}

// countingPermConfirm is a stand-in for the C18 confirmation leg with exactly the
// two properties these cases need from one: it records every switch it was asked
// about, and it can refuse. It is the same shape run_mode101_test.go's stub has,
// which is why the run seam's reading and these can be compared at all.
type countingPermConfirm struct {
	mu   sync.Mutex
	to   []risk.Mode
	fail error
}

func (c *countingPermConfirm) fn() perm.ConfirmFunc {
	return func(_ context.Context, sw perm.Switch) (perm.Card, error) {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.to = append(c.to, sw.To)
		if c.fail != nil {
			return perm.Card{}, c.fail
		}
		return perm.AnsweredCard("t114r2"), nil
	}
}

func (c *countingPermConfirm) calls() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.to)
}

func (c *countingPermConfirm) asked() []risk.Mode {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]risk.Mode(nil), c.to...)
}

// modeRequest114r2 is one panel.mode.request envelope, spelled the way the page
// spells it: bridge.go accepts exactly one sender identity, and anything else is
// refused by a door in front of the one under test.
func modeRequest114r2(rid, to string) string {
	return `{"method":"panel.mode.request","requestId":"` + rid + `","source":"` +
		inboundPanelSource + `","to":"` + to + `"}`
}

// TestTicket114r2ResidentPanelChainFromRefusalToTheWriteLeg is readings 1 and 2:
// the resident door, the nil it used to be handed, and the leg it is handed now.
// Its name is the reading 114-r2 asked for; the no-leg expectation inside is R9's
// (the middle-tier loosening is refused, and with a leg it costs one card).
//
// It calls newResidentComposerDispatch (the production builder the resident boot
// uses) and then PanelManager.dispatchRaw (the door the page's postMessage reaches,
// panel_host_windows.go:948). Neither the chain nor the door is rebuilt here.
func TestTicket114r2ResidentPanelChainFromRefusalToTheWriteLeg(t *testing.T) {
	ctx := context.Background()

	// (1) NO LEG: a request that LOOSENS the档 cannot be carded here, so R9 has the
	// panel gate refuse it before the write leg and the file stays where it was.
	// The read side is unchanged: this Store has no confirmation channel, so its
	// effective档 is ask_every_step (perm's store_372_test.go).
	beforeDir := newInboundDataDir33(t, "ask_every_step")
	beforeDisp, _, _, err := newResidentComposerDispatch(beforeDir, auditQuiet114r2, nil)
	if err != nil {
		t.Fatalf("the resident chain did not assemble: %v", err)
	}
	beforeHost := NewPanelManager(beforeDisp, nil, filepath.Join(beforeDir, "panel-webview2"))
	if _, err := beforeHost.dispatchRaw(ctx, modeRequest114r2("r114r2-before-ae-ah", "ask_high_risk")); !errors.Is(err, panel.ErrNoL2Confirm) {
		t.Errorf("with no leg attached, a loosening must be refused with ErrNoL2Confirm (R9); got %v", err)
	}
	if got := readInboundCfg33(t, beforeDir); got != "ask_every_step" {
		t.Errorf("the refused loosening still reached a write leg: disk reads %q, want ask_every_step", got)
	}

	// (2) AFTER: the leg the resident boot now hands in (resident_windows.go's
	// panelModeConfirm(ra.confirmModeSwitch), as the test's own confirm stub). The
	// same request now costs exactly one card and lands.
	afterDir := newInboundDataDir33(t, "ask_every_step")
	leg := &countingPermConfirm{}
	afterDisp, _, _, err := newResidentComposerDispatch(afterDir, auditQuiet114r2, panelModeConfirm(leg.fn()))
	if err != nil {
		t.Fatalf("the resident chain did not assemble with a leg: %v", err)
	}
	afterHost := NewPanelManager(afterDisp, nil, filepath.Join(afterDir, "panel-webview2"))
	if _, err := afterHost.dispatchRaw(ctx, modeRequest114r2("r114r2-ae-ah", "ask_high_risk")); err != nil {
		t.Errorf("a loosening with the leg attached was refused: %v", err)
	}
	if n := leg.calls(); n != 1 {
		t.Errorf("ask_every_step -> ask_high_risk asked the L2 leg %d time(s), want 1 (R9: every loosening "+
			"shows a card). Asked for: %v", n, leg.asked())
	}
	if got := readInboundCfg33(t, afterDir); got != "ask_high_risk" {
		t.Errorf("the request reached no write leg: the档 on disk reads %q, want ask_high_risk", got)
	}

	// (3) The auto_approve step through the same door: exactly one more card, so
	// the pair is one card per loosening rather than one card per process.
	if _, err := afterHost.dispatchRaw(ctx, modeRequest114r2("r114r2-aa", "auto_approve")); err != nil {
		t.Errorf("the auto_approve switch was refused even though the leg is attached: %v", err)
	}
	if got := leg.asked(); len(got) != 2 || got[0] != risk.ModeAskHighRisk || got[1] != risk.ModeAutoApprove {
		t.Errorf("L2 confirmations asked = %v, want exactly [ask_high_risk auto_approve] "+
			"(one per loosening, R9)", got)
	}
	if got := readInboundCfg33(t, afterDir); got != "auto_approve" {
		t.Errorf("an answered auto_approve switch did not reach the write leg: disk reads %q", got)
	}
}

// TestTicket114r2CLIPipeStillRefusesEveryWidening is reading 3: the CLI seam's
// production assembly. The NAME is 114-r2's reading ("every widening refused");
// after ticket 372 it refused only the auto_approve target, and since R9 it refuses
// every loosening again - for the opposite reason: not because a widening is
// refused for being wider, but because this pipe has no card surface. The two
// subtexts pin both halves so neither can change silently.
//
// What must NOT happen is that a tightening widens what the pipe's own process
// enforces - its Store has no leg, so its effective read stays ask_every_step
// whatever the file says (372-r2's read-side downgrade). The second subtest reads
// that back through the same config Manager -> perm.Store shape rather than
// trusting the dispatch's return.
func TestTicket114r2CLIPipeStillRefusesEveryWidening(t *testing.T) {
	ctx := context.Background()

	t.Run("auto_approve is refused", func(t *testing.T) {
		dir := newInboundDataDir33(t, "ask_every_step")
		disp, err := newPanelInboundDispatch(dir, auditQuiet114r2)
		if err != nil {
			t.Fatalf("the CLI chain did not assemble: %v", err)
		}
		if _, err := disp.Handle(ctx, modeRequest114r2("r114r2-cli-auto_approve", "auto_approve")); !errors.Is(err, panel.ErrNoL2Confirm) {
			t.Errorf("the pipe has no card surface, so a request for auto_approve must be refused with ErrNoL2Confirm; got %v", err)
		}
		if got := readInboundCfg33(t, dir); got != "ask_every_step" {
			t.Errorf("the refused auto_approve request moved the档 on disk: got %q, want ask_every_step", got)
		}
	})

	t.Run("ask_high_risk is refused too and the file it cannot reach stays strict", func(t *testing.T) {
		// Seeded at ask_high_risk (as a hand-edited config would be): the loosening
		// request is refused, so that value is NOT what this subtest is measuring -
		// the effective read is.
		dir := newInboundDataDir33(t, "ask_high_risk")
		disp, err := newPanelInboundDispatch(dir, auditQuiet114r2)
		if err != nil {
			t.Fatalf("the CLI chain did not assemble: %v", err)
		}
		if _, err := disp.Handle(ctx, modeRequest114r2("r114r2-cli-ask_every_step", "ask_every_step")); err != nil {
			t.Errorf("the pipe refused a tightening, which is never guarded: %v", err)
		}
		if got := readInboundCfg33(t, dir); got != "ask_every_step" {
			t.Errorf("the tightening reached no write leg: disk reads %q, want ask_every_step", got)
		}
		mgr, err := config.NewManager(filepath.Join(dir, configFileName), nil)
		if err != nil {
			t.Fatalf("cannot reopen the pipe's config: %v", err)
		}
		st, err := perm.New(perm.Options{Manager: mgr})
		if err != nil {
			t.Fatalf("perm.New: %v", err)
		}
		if got := st.PermissionMode(); got != risk.ModeAskEveryStep {
			t.Errorf("the pipe's process reads %s; with no leg it must read ask_every_step (372-r2)", got)
		}
	})

	t.Run("the loosening the pipe refuses leaves the file alone", func(t *testing.T) {
		dir := newInboundDataDir33(t, "ask_high_risk")
		disp, err := newPanelInboundDispatch(dir, auditQuiet114r2)
		if err != nil {
			t.Fatalf("the CLI chain did not assemble: %v", err)
		}
		if _, err := disp.Handle(ctx, modeRequest114r2("r114r2-cli-ah-aa", "auto_approve")); !errors.Is(err, panel.ErrNoL2Confirm) {
			t.Errorf("ask_high_risk -> auto_approve is a loosening on a pipe with no card surface; want "+
				"ErrNoL2Confirm, got %v", err)
		}
		if got := readInboundCfg33(t, dir); got != "ask_high_risk" {
			t.Errorf("the refused loosening moved the file: reads %q, want ask_high_risk", got)
		}
	})
}

// TestTicket114r2ANonNilLegAtTheCLISeamStillOwesAnAnsweredCard is the executable
// reason the CLI seam keeps its nil - and, since R9, also the reason a non-nil leg
// there would NOT silently ungate the pipe.
//
// Before R9 this case was named "...WouldUngateThePipe" and asserted the hole: any
// non-nil value switched the handler's gate off, and because M4 carded only
// auto_approve, a widening then landed without a card. R9 makes every loosening ask
// for one, so a leg that cannot show a card - all a stdin/stdout pipe could ever be
// - is REFUSED instead of silently applied. What the case pins now is that the
// failure is a refusal that names the unanswered card, and that nothing lands. If
// someone later gives this seam a real card surface, the same request applies and
// this case goes red, so the decision has to be taken knowingly.
func TestTicket114r2ANonNilLegAtTheCLISeamStillOwesAnAnsweredCard(t *testing.T) {
	ctx := context.Background()
	dir := newInboundDataDir33(t, "ask_every_step")
	leg := &countingPermConfirm{fail: errors.New("这条管子上没有可以出示卡片的表面")}
	disp, err := newComposerDispatchChain(dir, auditQuiet114r2, panelInboundActor, nil, panelModeConfirm(leg.fn()), nil)
	if err != nil {
		t.Fatalf("the chain did not assemble: %v", err)
	}

	if _, err := disp.Handle(ctx, modeRequest114r2("r114r2-cli-fake-leg", "ask_high_risk")); err == nil {
		t.Fatalf("a loosening was applied by a leg that cannot show a card (R9: the switch must be refused)")
	}
	if n := leg.calls(); n != 1 {
		t.Errorf("the gate asked the leg %d time(s) for ask_every_step -> ask_high_risk, want 1 "+
			"(R9: a loosening asks, and a leg that cannot answer is not an answer)", n)
	}
	if got := readInboundCfg33(t, dir); got != "ask_every_step" {
		t.Errorf("the refused loosening still landed: disk reads %q, want ask_every_step", got)
	}
}

// TestTicket114r2TheResidentAssemblyRootHandsTheChainTheOrbLeg is the wiring pin.
// The two cases above prove the chain USES a leg it is given; nothing in them says
// the resident boot gives it one, because a test that passes the leg itself cannot
// see that choice. This one reads the boot's own source and asks exactly that.
//
// It is a source scan with no fixtures for the reason panel_inbound_33_test.go
// gives for its own: reverting the argument is a one-token edit that compiles, and a
// chain back on nil refuses every loosening that a real user would ask for - a
// regression with no crash and no failing behaviour test.
func TestTicket114r2TheResidentAssemblyRootHandsTheChainTheOrbLeg(t *testing.T) {
	const boot = "resident_windows.go"
	src, err := os.ReadFile(boot)
	if err != nil {
		t.Fatalf("cannot read the resident boot: %v", err)
	}
	f, err := parser.ParseFile(token.NewFileSet(), boot, src, 0)
	if err != nil {
		t.Fatalf("cannot parse %s: %v", boot, err)
	}

	call := findPanelManagerCall114r2(f)
	if call == nil {
		t.Fatalf("no newResidentPanelManager(...) call in %s: the resident panel's chain is assembled somewhere this pin does not know, so the L2 leg it is handed is unmeasured", boot)
	}
	if len(call.Args) != 2 {
		t.Fatalf("newResidentPanelManager has %d argument(s) in %s; this pin reads two (the data root and the L2 leg)", len(call.Args), boot)
	}
	leg, ok := call.Args[1].(*ast.CallExpr)
	if !ok {
		t.Fatalf("the resident boot hands the panel chain %T instead of the L2 leg: a nil (or any non-call) there puts every loosening back on the refusal path, and the panel could never ask for a looser档", call.Args[1])
	}
	fn, ok := leg.Fun.(*ast.Ident)
	if !ok || fn.Name != "panelModeConfirm" {
		t.Fatalf("the resident boot's L2 leg is not routed through panelModeConfirm, so it is not the same adapter the run seam uses (run.go's rt.modeWrites) and the two cards can drift into two sentences")
	}
	if len(leg.Args) != 1 {
		t.Fatalf("panelModeConfirm takes exactly one argument; the boot passes %d", len(leg.Args))
	}
	sel, ok := leg.Args[0].(*ast.SelectorExpr)
	if !ok {
		t.Fatalf("panelModeConfirm's argument is %T, not a method value on this process's approval gate", leg.Args[0])
	}
	x, ok := sel.X.(*ast.Ident)
	if !ok || x.Name != "ra" || sel.Sel == nil || sel.Sel.Name != "confirmModeSwitch" {
		t.Fatalf("the resident boot's L2 leg is not ra.confirmModeSwitch, so it is not this process's own gate and the card would go somewhere other than the orb (found %s)", exprText114r2(sel))
	}
}

// findPanelManagerCall114r2 returns the newResidentPanelManager call, or nil.
func findPanelManagerCall114r2(f *ast.File) *ast.CallExpr {
	var found *ast.CallExpr
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || found != nil {
			return true
		}
		if id, ok := call.Fun.(*ast.Ident); ok && id.Name == "newResidentPanelManager" {
			found = call
		}
		return true
	})
	return found
}

// exprText114r2 is only used to make a failure sentence name what it found.
func exprText114r2(e ast.Expr) string {
	if sel, ok := e.(*ast.SelectorExpr); ok {
		if id, ok := sel.X.(*ast.Ident); ok && sel.Sel != nil {
			return id.Name + "." + sel.Sel.Name
		}
	}
	return "an unrecognised expression"
}

// TestTicket114r2RunSeamPanelHandlerStillCarriesTheSameLeg keeps the run seam's
// half honest: rt.modeWrites is the handler the ticket named, it is built with the
// adapter both legs now share, and one panel request through it behaves exactly as
// the resident one does (both loosening steps ask one card each).
//
// This is also the case that would catch the refactor's one risk: if
// panelModeConfirm returned nil for a usable leg, the field would look assembled and
// refuse everything, which no other assertion in this package looks at.
func TestTicket114r2RunSeamPanelHandlerStillCarriesTheSameLeg(t *testing.T) {
	h := t101boot(t, "")
	leg := &countingPermConfirm{}
	var handled error
	hook := func(rt *agentRuntime) {
		if rt.modeWrites == nil {
			handled = errors.New("the run seam assembled no panel mode handler at all")
			return
		}
		if rt.modeWrites.Confirm == nil {
			handled = errors.New("rt.modeWrites lost its L2 leg: panelModeConfirm dropped a usable confirm, " +
				"so every panel widening would be refused by the handler while the assembly looks complete")
			return
		}
		disp := &panel.ComposerDispatch{Mode: rt.modeWrites}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := disp.Handle(ctx, modeRequest114r2("r114r2-run-ae-ah", "ask_high_risk")); err != nil {
			handled = err
			return
		}
		_, handled = disp.Handle(ctx, modeRequest114r2("r114r2-run-aa", "auto_approve"))
	}
	code, log := h.start(t, leg.fn(), hook)
	if code != 0 {
		t.Fatalf("the run seam exited %d\n%s", code, log)
	}
	if handled != nil {
		t.Errorf("the run seam's panel handler did not answer the two requests: %v", handled)
	}
	want := []risk.Mode{risk.ModeAskHighRisk, risk.ModeAutoApprove}
	if got := leg.asked(); len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("L2 confirmations asked = %v, want exactly %v (one per loosening, R9)", got, want)
	}
	if got := h.readMode(t); got != risk.ModeAutoApprove {
		t.Errorf("the answered switch never reached the config file: a cold read says %v", got)
	}
}
