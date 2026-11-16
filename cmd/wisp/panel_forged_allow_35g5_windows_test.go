//go:build windows

package main

// Ticket 35 cell 票面 :47 ("Forged panel allow", today's line number) - the
// three behavioural teeth the dispatch order names:
//
//	① a test injects allow:true from the BRIDGE CALL FACE (the page's own
//	  postMessage edge, the same machinery the (d) trio of
//	  panel_inbound_guards_35r3_test.go rides) and the server answers with a
//	  NAMED refusal - the red sentence verbatim, never a silent drop;
//	② locally turning the refusal into an allow (overlay mutation, repo
//	  untouched) reddens exactly this test; restoring greens it. The mutation
//	  readings live in docs/evidence/s1/35-g5-forged-allow-and-resync-20261116.md;
//	③ the L1/L2 legal paths are unaffected: the positive control below proves
//	  the same pipe still delivers a legitimate roster request to its handler,
//	  and the native-second-confirm route keeps its own named ruler
//	  (TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject,
//	  approval_reply_201_test.go:344), re-run in the same wave.
//
// WHY THE REFUSAL IS STRUCTURAL, NOT FRONT-END SELF-RESTRAINT. The page can
// post ANY string through window.chrome.webview.postMessage; the transport
// forwards it verbatim into the bound door. `approval.decide` is NOT on the
// C17 roster (internal/panel/bridge.go's knownComposerMethod case line: the
// eight answered names are panel.mode.request / panel.workspace.request /
// panel.attachment.add / panel.message.send / config.get / config.set /
// config.models.list / cost.summary), so a forged approval.decide{allow:true}
// dies in ParseComposerRequest's roster guard before any handler can see it.
// This leg adds NO roster entry - the refusal IS the claim, and widening the
// roster to make this test pass would be the exact forbidden move (AGENTS.md
// §1.2 ban #6: no panel-side L2 "allow"; the frozen grant-boundary instrument
// refuses any answered route carrying an approval word).
//
// THE ENVELOPE IS THE STRONGEST FORGERY, NOT A WEAK ONE. It carries the AC's
// literal method name and an allow:true body, PLUS a valid source and
// requestId - so every guard AFTER the roster one would pass. That isolates
// the roster guard as the one structural answer: if the refusal moves to the
// source guard (mutation M-forged-1) or the request reaches the dispatch
// table (mutation M-forged-2), the named sentence below disappears and this
// test goes red. The page's own sender shape (panel.approval.request +
// outcome, frontend/src/lib/panel.ts:630-633) is already pinned on this edge
// by (d)-1 in panel_inbound_guards_35r3_test.go; that nail is referenced, not
// duplicated or weakened here.
//
// SCOPE NOTES. d22scan ban #6 walks frontend/ only - this file is a Go test
// asserting the server-side refusal, not an implementation of a panel allow.
// The mutation runs used `go test -c -overlay` with copies outside the
// repository; `git diff` over internal/panel/ and cmd/wisp/ production files
// stayed empty across every mutation run (readings in the evidence file).

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
)

// envelope35g5ForgedAllow is the forged shape: the AC's literal method name,
// an allow body, and a fully-formed composer envelope (valid source +
// requestId) so the roster guard is the ONLY thing between the forged allow
// and the dispatch table.
const envelope35g5ForgedAllow = `{"method":"approval.decide","correlationId":"pc-35g5-allow-1","allow":true,"source":"panel-composer","requestId":"pc-35g5-allow-1"}`

// envelope35g5LegitMode is the positive control: the SAME envelope shape with
// a roster method instead of the forged one. If the pipe refused everything,
// tooth ① above would be a tautology; this case keeps it honest.
const envelope35g5LegitMode = `{"method":"panel.mode.request","requestId":"pc-35g5-ok-1","source":"panel-composer","to":"ask_every_step"}`

// auditRecorder35g5 captures the dispatch's "[audit]" lines, so the test can
// prove the refusal was RECORDED (a named, audited refusal) and not a silent
// drop - the dispatch's ⛔ 静默丢弃 clause has two faces: the page must get the
// sentence back, and the audit must book it.
type auditRecorder35g5 struct {
	mu    sync.Mutex
	lines []string
}

func (a *auditRecorder35g5) auditf(format string, args ...any) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.lines = append(a.lines, fmt.Sprintf(format, args...))
}

func (a *auditRecorder35g5) snapshot() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]string(nil), a.lines...)
}

// deliverViaPageEdge35g5 wires the shipping transport onto a fresh document
// and posts one envelope exactly the way the page does - the 35r3 page-edge
// harness shape, with an audit recorder attached. It returns the yard, the
// mode spy and the audit recorder; the caller asserts reply text, handler
// traffic and the audit line.
func deliverViaPageEdge35g5(t *testing.T, envelope string) (*fakeDoc35r2, *modeSpy35r1, *auditRecorder35g5) {
	t.Helper()

	spy := &modeSpy35r1{}
	audit := &auditRecorder35g5{}
	disp := &panel.ComposerDispatch{Mode: spy, Audit: audit.auditf}
	mgr := NewPanelManager(disp, nil, "")

	bf := newFakeDoc35r2()
	if err := mgr.installPanelTransport(bf, context.Background()); err != nil {
		t.Fatalf("installPanelTransport failed: %v", err)
	}
	bf.openDocument()

	// The page's own hop: chrome.webview.postMessage(<envelope>). A throw here
	// would mean the envelope never left the page, so the door result below
	// would be stale (the same precondition the 35r3 harness states).
	if thrown := bf.pagePost(envelope); thrown != "" {
		t.Fatalf("35g5 RED: the page's own postMessage threw in the document: %s. %s",
			thrown, bf.summary())
	}

	// ★EDGE PROOF: the raw string reached dispatchRaw's door exactly once only
	// because the forwarding hook routed it into the bound door. Without the
	// forwarding statement the envelope dies at msgcb's unbound branch
	// (doorRounds=0) and this fires - the :63 falsifiability clause these
	// cases ride, restated so a deleted hook cannot masquerade as a refusal.
	if bf.doorRounds != 1 {
		t.Fatalf("35g5 RED: the page's postMessage drove the door %d time(s), want 1. "+
			"Without the forwarding statement the envelope never reaches dispatchRaw. %s",
			bf.doorRounds, bf.summary())
	}
	return bf, spy, audit
}

// TestPanelForgedApprovalDecideAllowIsRefusedByNameOnPageEdge35g5 is tooth ①:
// approval.decide{allow:true} posted from the page's own edge gets back the
// roster guard's named refusal sentence, verbatim, with the forged method
// named inside it - and nothing executed the allow.
func TestPanelForgedApprovalDecideAllowIsRefusedByNameOnPageEdge35g5(t *testing.T) {
	bf, spy, audit := deliverViaPageEdge35g5(t, envelope35g5ForgedAllow)

	// NOT a silent drop, face 1: the door answered and the page got a sentence
	// back - non-empty, carrying the refusal header and the forged method's
	// own name (RefusedEnvelopeForUser's "面板请求被拒绝 [approval.decide …]"
	// prefix plus the roster guard's cause).
	if bf.lastReply == "" {
		t.Fatalf("35g5 RED: the door answered with an EMPTY reply - a forged panel allow "+
			"must come back with the named refusal sentence, not a silence. %s", bf.summary())
	}
	if !strings.Contains(bf.lastReply, "面板请求被拒绝") {
		t.Errorf("35g5 RED: the reply is not the refusal sentence the page renders: %q", bf.lastReply)
	}
	if !strings.Contains(bf.lastReply, "approval.decide") {
		t.Errorf("35g5 RED: the refusal does not name the forged method approval.decide: %q", bf.lastReply)
	}
	// The red sentence, verbatim (internal/panel/bridge.go's roster guard,
	// shared constant with the (d) trio so the two files cannot drift apart).
	if !strings.Contains(bf.lastReply, guard35r3Roster) {
		t.Errorf("35g5 RED: the reply does not carry the roster guard's sentence %q; "+
			"a forged panel allow that loses its named refusal is a hole, not a style change. reply=%q",
			guard35r3Roster, bf.lastReply)
	}
	// The other two guard sentences must NOT be the answer: the envelope is
	// fully formed (valid source + requestId), so the roster guard is
	// specifically what refused it. A mutation that moves the refusal to a
	// different guard reddens through this prong too.
	if strings.Contains(bf.lastReply, guard35r3Source) {
		t.Errorf("35g5 RED: the SOURCE guard answered a fully-sourced envelope - the refusal moved: %q", bf.lastReply)
	}
	if strings.Contains(bf.lastReply, guard35r3RequestID) {
		t.Errorf("35g5 RED: the requestId guard answered an envelope that carries one - the refusal moved: %q", bf.lastReply)
	}

	// Nothing executed the allow: the refusal happens inside
	// ParseComposerRequest, before dispatch runs, so no handler can have been
	// reached - this is the "structural, not front-end self-restraint" claim
	// made checkable.
	if reqs := spy.snapshot(); len(reqs) != 0 {
		t.Errorf("35g5 RED: a forged approval.decide{allow:true} reached a handler %d time(s): %+v - "+
			"the allow was executed, not refused", len(reqs), reqs)
	}

	// NOT a silent drop, face 2: the audit booked the refusal and named the
	// method (composer_dispatch.go's record on the parse-refusal branch).
	lines := audit.snapshot()
	refusals := 0
	for _, line := range lines {
		if strings.Contains(line, "INBOUND-DISPATCH") {
			refusals++
			if !strings.Contains(line, `method="approval.decide"`) {
				t.Errorf("35g5 RED: the audit line for the refusal does not name the forged method: %q", line)
			}
		}
	}
	if refusals != 1 {
		t.Errorf("35g5 RED: exactly one INBOUND-DISPATCH audit line is owed for one refused envelope; got %d of %d line(s): %v",
			refusals, len(lines), lines)
	}
	t.Logf("35g5 forged allow on the page edge: door reply=%q audit=%d line(s)", bf.lastReply, len(lines))
}

// TestPanelForgedAllowPositiveControlLegitRequestStillDelivered35g5 is tooth
// ③'s bridge face: the same pipe, the same harness, the same envelope shape
// with a roster method instead of the forged name - the request reaches the
// mode handler, no refusal is recorded, and the door returns the mode route's
// own empty receipt. This is what keeps tooth ① from being a blanket refusal
// wearing a security costume: the door's verdict differs BY METHOD, which is
// the roster doing its job.
func TestPanelForgedAllowPositiveControlLegitRequestStillDelivered35g5(t *testing.T) {
	bf, spy, audit := deliverViaPageEdge35g5(t, envelope35g5LegitMode)

	if reqs := spy.snapshot(); len(reqs) != 1 {
		t.Fatalf("35g5 RED: a legitimate roster request reached the mode handler %d time(s), want exactly 1: %+v",
			len(reqs), reqs)
	} else if reqs[0].RequestID != "pc-35g5-ok-1" {
		t.Errorf("35g5 RED: the handler saw requestId %q, want the page's own pc-35g5-ok-1", reqs[0].RequestID)
	}
	if bf.lastReply != "" {
		t.Errorf("35g5 RED: a delivered mode request must return the mode route's empty receipt, got %q", bf.lastReply)
	}
	for _, line := range audit.snapshot() {
		if strings.Contains(line, "INBOUND-DISPATCH") {
			t.Errorf("35g5 RED: a delivered request must leave no INBOUND-DISPATCH refusal line; got %q", line)
		}
	}
	t.Logf("35g5 positive control: legit panel.mode.request delivered, 0 refusal line(s)")
}
