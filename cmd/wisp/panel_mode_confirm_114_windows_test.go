//go:build windows

package main

// Ticket 114 AC#3 (R-92-2's other half): the panel's mode chain, driven through
// the PRODUCTION doors of both legs, before and after the L2 confirmation leg was
// attached (ticket 114-r2).
//
// WHAT THESE CASES ARE FOR, in the ticket's own words. AC#2 put a gate in front of
// the panel's mode write that refuses a widening request when the injected Confirm
// is nil. That gate was correct and unreachable: both production assemblies passed
// nil, so "no L2 leg" was a constant and the panel could never ask for a looser档
// at all. 114-r2 made the leg a parameter of the shared chain and the resident boot
// hands it the orb's own approval gate. These cases pin the two readings that
// change has to buy, and the one reading that must NOT change:
//
//  1. the resident leg (panel_host_windows.go's dispatchRaw, i.e. the page's door):
//     ask_every_step -> ask_high_risk used to be REFUSED and now reaches the write
//     leg - with ZERO cards, because R20/M4 only cards the switch INTO auto_approve;
//  2. the same leg, same door: auto_approve asks exactly ONE card, and it is the
//     auto_approve one. That pair is the proof that wiring the leg did not widen
//     M4's reach (the ticket's hard boundary); run_mode101_test.go:250 keeps the
//     same pair for the run seam, and it must stay green;
//  3. the CLI pipe (panel_inbound.go's chain): every widening is still refused, and
//     that is not a leftover stub - the case after it RUNS the alternative and shows
//     why a leg there would be a hole rather than a fix.

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
	return func(_ context.Context, sw perm.Switch) error {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.to = append(c.to, sw.To)
		return c.fail
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
//
// It calls newResidentComposerDispatch (the production builder the resident boot
// uses) and then PanelManager.dispatchRaw (the door the page's postMessage reaches,
// panel_host_windows.go:948). Neither the chain nor the door is rebuilt here.
func TestTicket114r2ResidentPanelChainFromRefusalToTheWriteLeg(t *testing.T) {
	ctx := context.Background()

	// (1) BEFORE: the nil both assemblies passed until 114-r2.
	beforeDir := newInboundDataDir33(t, "ask_every_step")
	beforeDisp, _, err := newResidentComposerDispatch(beforeDir, auditQuiet114r2, nil)
	if err != nil {
		t.Fatalf("the resident chain did not assemble: %v", err)
	}
	beforeHost := NewPanelManager(beforeDisp, nil, filepath.Join(beforeDir, "panel-webview2"))
	if _, err := beforeHost.dispatchRaw(ctx, modeRequest114r2("r114r2-before-ae-ah", "ask_high_risk")); !errors.Is(err, panel.ErrNoL2Confirm) {
		t.Errorf("with no leg attached the resident door must refuse with ErrNoL2Confirm; got %v", err)
	}
	if got := readInboundCfg33(t, beforeDir); got != "ask_every_step" {
		t.Errorf("a refused widening moved the档 on disk: got %q, want ask_every_step", got)
	}

	// (2) AFTER: the leg the resident boot now hands in (resident_windows.go's
	// panelModeConfirm(ra.confirmModeSwitch), as the test's own confirm stub).
	afterDir := newInboundDataDir33(t, "ask_every_step")
	leg := &countingPermConfirm{}
	afterDisp, _, err := newResidentComposerDispatch(afterDir, auditQuiet114r2, panelModeConfirm(leg.fn()))
	if err != nil {
		t.Fatalf("the resident chain did not assemble with a leg: %v", err)
	}
	afterHost := NewPanelManager(afterDisp, nil, filepath.Join(afterDir, "panel-webview2"))
	if _, err := afterHost.dispatchRaw(ctx, modeRequest114r2("r114r2-ae-ah", "ask_high_risk")); err != nil {
		t.Errorf("a widening that M4 does not card is still refused with the leg attached: %v", err)
	}
	if n := leg.calls(); n != 0 {
		t.Errorf("R20/M4 was widened: ask_every_step -> ask_high_risk asked the L2 leg %d time(s), "+
			"want 0 (only the switch INTO auto_approve costs a card). Asked for: %v", n, leg.asked())
	}
	if got := readInboundCfg33(t, afterDir); got != "ask_high_risk" {
		t.Errorf("the request reached no write leg: the档 on disk reads %q, want ask_high_risk", got)
	}

	// (3) The card that DOES belong to this leg, through the same door: exactly one,
	// and it is the auto_approve one.
	if _, err := afterHost.dispatchRaw(ctx, modeRequest114r2("r114r2-aa", "auto_approve")); err != nil {
		t.Errorf("the auto_approve switch was refused even though the leg is attached: %v", err)
	}
	if got := leg.asked(); len(got) != 1 || got[0] != risk.ModeAutoApprove {
		t.Errorf("L2 confirmations asked = %v, want exactly [auto_approve] (R20/M4)", got)
	}
	if got := readInboundCfg33(t, afterDir); got != "auto_approve" {
		t.Errorf("an answered auto_approve switch did not reach the write leg: disk reads %q", got)
	}
}

// TestTicket114r2CLIPipeStillRefusesEveryWidening is reading 3: the CLI seam's
// production assembly, unchanged by 114-r2. It is here so that "one change, two
// legs" cannot be claimed later without a reader noticing that one of the two still
// refuses - and so that removing the nil from that call site reddens a named case
// instead of quietly opening a pipe.
func TestTicket114r2CLIPipeStillRefusesEveryWidening(t *testing.T) {
	ctx := context.Background()
	for _, to := range []string{"ask_high_risk", "auto_approve"} {
		t.Run(to, func(t *testing.T) {
			dir := newInboundDataDir33(t, "ask_every_step")
			disp, err := newPanelInboundDispatch(dir, auditQuiet114r2)
			if err != nil {
				t.Fatalf("the CLI chain did not assemble: %v", err)
			}
			if _, err := disp.Handle(ctx, modeRequest114r2("r114r2-cli-"+to, to)); !errors.Is(err, panel.ErrNoL2Confirm) {
				t.Errorf("the pipe has no card surface, so a widening to %s must be refused with ErrNoL2Confirm; got %v", to, err)
			}
			if got := readInboundCfg33(t, dir); got != "ask_every_step" {
				t.Errorf("the pipe widened the档 to %q with no card anywhere", got)
			}
		})
	}
}

// TestTicket114r2ANonNilLegAtTheCLISeamWouldUngateThePipe is the executable reason
// the CLI seam keeps its nil. ModeWriteHandler reads only nil-ness, so ANY non-nil
// value at that call site - including a leg that can never show a card, which is all
// a stdin/stdout pipe could ever be - switches the gate off. What follows the switch
// is not a card: it is a widening that reaches perm.Store.Set and lands, because M4
// cards only auto_approve. This case asserts that outcome on purpose; if someone
// later gives that seam a real card surface, it goes red and the decision has to be
// taken knowingly.
func TestTicket114r2ANonNilLegAtTheCLISeamWouldUngateThePipe(t *testing.T) {
	ctx := context.Background()
	dir := newInboundDataDir33(t, "ask_every_step")
	leg := &countingPermConfirm{fail: errors.New("这条管子上没有可以出示卡片的表面")}
	disp, err := newComposerDispatchChain(dir, auditQuiet114r2, panelInboundActor, nil, panelModeConfirm(leg.fn()))
	if err != nil {
		t.Fatalf("the chain did not assemble: %v", err)
	}

	if _, err := disp.Handle(ctx, modeRequest114r2("r114r2-cli-fake-leg", "ask_high_risk")); err != nil {
		t.Fatalf("this case exists to show the hole the nil closes; the leg was refused instead: %v", err)
	}
	if n := leg.calls(); n != 0 {
		t.Errorf("M4 asked %d card(s) for ask_every_step -> ask_high_risk; the widening was supposed to be silent", n)
	}
	if got := readInboundCfg33(t, dir); got != "ask_high_risk" {
		t.Errorf("expected the unguarded widening to land (that IS the hole); disk reads %q", got)
	}
}

// TestTicket114r2TheResidentAssemblyRootHandsTheChainTheOrbLeg is the wiring pin.
// The two cases above prove the chain USES a leg it is given; nothing in them says
// the resident boot gives it one, because a test that passes the leg itself cannot
// see that choice. This one reads the boot's own source and asks exactly that.
//
// It is a source scan with no fixtures for the reason panel_inbound_33_test.go
// gives for its own: reverting the argument is a one-token edit that compiles, and a
// chain back on nil refuses every widening that a real user would ask for - a
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
		t.Fatalf("the resident boot hands the panel chain %T instead of the L2 leg: a nil (or any non-call) there puts every widening back on the refusal path, and the panel could never ask for a looser档", call.Args[1])
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
// the resident one does (widening lands with no card, auto_approve asks once).
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
	if got := leg.asked(); len(got) != 1 || got[0] != risk.ModeAutoApprove {
		t.Errorf("L2 confirmations asked = %v, want exactly [auto_approve] (R20/M4)", got)
	}
	if got := h.readMode(t); got != risk.ModeAutoApprove {
		t.Errorf("the answered switch never reached the config file: a cold read says %v", got)
	}
}
