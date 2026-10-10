//go:build windows

package main

// Ticket 35 :63 "Inbound judge must ride the page's own envelope" - the (d) trio.
//
// WHAT THIS FILE ADDS AND WHAT ALREADY EXISTED. The :63 frame needs four things: (a)
// the raw envelope verbatim from what the page posts, (b) entry through the transport
// the page uses AFTER shape 甲 sub-shape ③ landed, (c) dispatchRaw's downstream reached
// carrying the correlationId, (d) an unregistered name still refused. The r1/r2 cases
// already own the arrival/positive side (a roster method that PASSES the guard order,
// so it exercises none of the three reject branches - see panel_transport_35r1_test.go
// and panel_transport_35r2_test.go, whose assertions this file does NOT touch). (d)
// asks for the three INBOUND GUARDS to be proven red, each pointing at its OWN branch,
// none sharing a sentence. This file is those three and only those three.
//
// ★THE EDGE (the whole point of :63). Each case drives the page's own hop the way the
// r2 behavioural yard does: installPanelTransport binds the door and queues the real
// panelPostMessageForwardInit forwarding hook, openDocument runs that hook verbatim in
// the embedded JS interpreter (chromium.go:112 shim + webview.go:462-478 Bind stub +
// the production hook string), and pagePost is literally
// window.chrome.webview.postMessage(<envelope>) - the exact call frontend/src/lib/panel.ts
// :179/:218 makes. We never call dispatchRaw from Go, and we never run
// w.Eval("window.wispDispatch(...)") (that is N6, the test playing the page, zero
// discriminating power for this edge). The ONLY thing that carries the raw string to the
// door is the forwarding statement, so deleting installPanelTransport's w.Init(...) hook
// reddens all three arrivals below (the :63 falsifiability clause).
//
// ★THE THREE BRANCHES (orchestrator's guard-order ruling, ticket :69; verified this leg
// against bridge.go, CR count 0 == HEAD). ParseComposerRequest's次第 is parse :128 ->
// roster :132/refuse :133 -> source :135/refuse :136-137 -> requestId :139/refuse :140-141.
// (The parse guard at :128-129 is deliberately NOT covered here - out of :63's trio.)
// Because the roster guard runs first, an unregistered name can ONLY ever die at :133,
// so (d)-1 uses panel.approval.request (the name the page posts TODAY, panel.ts:181).
// (d)-2 and (d)-3 use a ROSTER method (panel.mode.request) so they reach the later two
// guards at all. Each case asserts its OWN keyword is in dispatchRaw's reply AND the
// other two keywords are NOT - that is "each points at its own", not "watch it go red".
//
// HOW THE REFUSAL TEXT REACHES THE TEST. ComposerDispatch.Handle
// (internal/panel/composer_dispatch.go:153) turns a ParseComposerRequest error into
// RefusedEnvelopeForUser(req, err), which embeds that guard's error verbatim. The bind
// closure drops the error and hands only that reply string across the JS boundary
// (panel_host_windows.go:694-696, `reply, _ := m.dispatchRaw(...)`), so the yard's
// lastReply is exactly the page-visible refusal sentence.
//
// (d)-1 is NOT the existing nail. internal/panel/l2_grant_boundary_test.go (:1964/:1965,
// :1251-1267) already calls knownComposerMethod / ParseComposerRequest directly and
// proves the roster set rejects panel.approval.request. That path never enters this
// transport, never runs Handle, and produces no page-visible refusal. The NEW surface
// here is: the same name, arriving on the page's own postMessage edge, is refused by the
// ROSTER branch specifically, named in the reply the page gets back. It does not
// substitute for the old nail and the old nail does not cover this.

import (
	"context"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
)

// The three guard keywords, verbatim from bridge.go's three return statements. No two
// cases may assert the same sentence (:63's "三支不许共用一句文案").
const (
	guard35r3Roster    = `不是面板 composer 通路的能力入口` // bridge.go:133 (name %q ... 不是面板 composer 通路的能力入口)
	guard35r3Source    = `按伪造/串台拒绝`              // bridge.go:136-137 (来源 %q 不是 %q，按伪造/串台拒绝...)
	guard35r3RequestID = `缺少 requestId`          // bridge.go:140-141 (缺少 requestId，无法与审计/卡片对齐...)
)

// The three envelopes, all in the page's own posting shape.
//
// (d)-1: panel.approval.request verbatim from panel.ts:179-184
// ({"method","correlationId","outcome"} - it carries NO source and NO requestId, which
// is why, once the roster guard is no-op'd, it does not silently sail through the other
// two: it lands on the source guard next, a DIFFERENT sentence).
// (d)-2: a roster method (panel.mode.request, via the :211 sendRequest that stamps
// requestId + source) with source tampered to "panel-composer-x".
// (d)-3: the same roster envelope with requestId removed, source left valid.
const (
	envelope35r3UnregisteredName = `{"method":"panel.approval.request","correlationId":"pc-35r3-approval-1","outcome":"allow"}`
	envelope35r3WrongSource      = `{"method":"panel.mode.request","requestId":"pc-35r3-2-9e7a-4c1b","source":"panel-composer-x","to":"ask_every_step"}`
	envelope35r3MissingID        = `{"method":"panel.mode.request","source":"panel-composer","to":"ask_every_step"}`
)

// deliverViaPageEdge35r3 wires the shipping transport onto a fresh document and posts
// one envelope exactly the way the page does. It returns the yard plus the mode spy so
// the caller can assert both the reply text and that the refusal happened BEFORE routing
// (the guard is what answered, not the handler).
func deliverViaPageEdge35r3(t *testing.T, envelope string) (*fakeDoc35r2, *modeSpy35r1) {
	t.Helper()

	spy := &modeSpy35r1{}
	disp := &panel.ComposerDispatch{Mode: spy}
	mgr := NewPanelManager(disp, nil, "")

	bf := newFakeDoc35r2()
	if err := mgr.installPanelTransport(bf, context.Background()); err != nil {
		t.Fatalf("installPanelTransport failed: %v", err)
	}
	bf.openDocument()

	// The page's own hop: chrome.webview.postMessage(<envelope>). A throw here would
	// mean the envelope never left the page, so the door result below would be stale.
	if thrown := bf.pagePost(envelope); thrown != "" {
		t.Fatalf("inbound guard RED: the page's own postMessage threw in the document: %s. %s",
			thrown, bf.summary())
	}

	// ★EDGE PROOF: the raw string reached dispatchRaw's door exactly once ONLY because
	// the forwarding hook routed it into the bound door. If installPanelTransport's
	// w.Init(forwarding) is deleted, the envelope dies at msgcb's unbound branch
	// (doorRounds=0) and this fires - which is the :63 falsifiability clause.
	if bf.doorRounds != 1 {
		t.Fatalf("inbound guard RED: the page's postMessage drove the door %d time(s), want 1. "+
			"Without the forwarding statement the envelope is an unbound RPC frame and never "+
			"reaches dispatchRaw. %s", bf.doorRounds, bf.summary())
	}
	return bf, spy
}

// assertGuardRedness35r3 pins the "each points at its own" contract: the reply must carry
// the one named guard's sentence and must NOT carry either of the other two. The absent
// keywords are what stops this from being the "看它红" the frame forbids.
func assertGuardRedness35r3(t *testing.T, bf *fakeDoc35r2, spy *modeSpy35r1, mine string, others ...string) {
	t.Helper()

	if !strings.Contains(bf.lastReply, mine) {
		t.Errorf("inbound guard RED: dispatchRaw's reply does not name guard %q. reply=%q | %s",
			mine, bf.lastReply, bf.summary())
	}
	for _, not := range others {
		if strings.Contains(bf.lastReply, not) {
			t.Errorf("inbound guard RED: dispatchRaw's reply leaked a DIFFERENT guard %q alongside %q - "+
				"the three branches are not landing on distinct sentences. reply=%q",
				not, mine, bf.lastReply)
		}
	}
	// A guard refusal happens inside ParseComposerRequest, before dispatch runs, so the
	// mode handler downstream must never have been reached. This separates "a guard said
	// no" from "a handler said no".
	if len(spy.reqs) != 0 {
		t.Errorf("inbound guard RED: a refusal should never reach the router, but the mode handler ran %d time(s): %+v",
			len(spy.reqs), spy.reqs)
	}
}

// (d)-1 unregistered name: the roster branch (bridge.go:132) refuses panel.approval.request
// - the method the page posts TODAY - as it arrives on the page's own edge, and names the
// roster guard in the reply. keyword: guard35r3Roster only.
func TestInboundRosterGuardRefusesUnregisteredNameOnPageEdge(t *testing.T) {
	bf, spy := deliverViaPageEdge35r3(t, envelope35r3UnregisteredName)
	assertGuardRedness35r3(t, bf, spy, guard35r3Roster, guard35r3Source, guard35r3RequestID)
	t.Logf("(d)-1 roster guard on the page edge: door reply=%q", bf.lastReply)
}

// (d)-2 roster method, tampered source: the roster guard passes (panel.mode.request is
// registered), the source guard (bridge.go:135) then refuses. keyword: guard35r3Source only.
func TestInboundSourceGuardRefusesForeignSourceOnPageEdge(t *testing.T) {
	bf, spy := deliverViaPageEdge35r3(t, envelope35r3WrongSource)
	assertGuardRedness35r3(t, bf, spy, guard35r3Source, guard35r3Roster, guard35r3RequestID)
	t.Logf("(d)-2 source guard on the page edge: door reply=%q", bf.lastReply)
}

// (d)-3 roster method, missing requestId: the roster and source guards both pass (valid
// name + "panel-composer"), the requestId guard (bridge.go:139) then refuses.
// keyword: guard35r3RequestID only.
func TestInboundRequestIDGuardRefusesMissingIDOnPageEdge(t *testing.T) {
	bf, spy := deliverViaPageEdge35r3(t, envelope35r3MissingID)
	assertGuardRedness35r3(t, bf, spy, guard35r3RequestID, guard35r3Roster, guard35r3Source)
	t.Logf("(d)-3 requestId guard on the page edge: door reply=%q", bf.lastReply)
}
