package main

// Ticket 329, form 甲: the booked packet's last hop from this run to the host.
//
// WHAT ALREADY EXISTED BEFORE THIS FILE. The pump builds the packet and books it:
// publishPanelSnapshot calls pump.Publish, Publish marshals the snapshot and hands
// the bytes to src.Out, and out.bookPanelSnapshot retains them on the runtime
// (cmd/wisp/panel_pump.go:360). What did NOT exist was any production reader of
// that retention: lastPanelSnapshot had zero non-test callers, so the bytes stopped
// at the ledger summary and never reached a page. That gap is ticket 329's subject,
// and this file closes the Go half of it.
//
// THE SHAPE WAS RULED, NOT INVENTED HERE (orchestrator ruling A1009). The push
// rides the host's existing Eval primitive - the same one msgcb uses in
// cmd/wisp/panel_control_windows.go for the four RPC receipts - as one named
// global call wrapping the pump's own marshal bytes verbatim. No new inbound
// Method* constant was added (the bridge.go roster stays at its pinned size),
// no new publish drive point was invented (this hop hangs ONLY inside publishPanelSnapshot, which every
// existing state-move trigger already calls), and no merge or backpressure rule is
// claimed here - the "how often may a panel be pushed" rule belongs to ticket 35's
// D38d line and this leg deliberately does not touch it. The outbound rows this hop
// lands are already frozen contract: docs/specs/SPEC-08-ui-ball-panel.md lists the
// Go-to-page pushes as outbound (no request), and C17 (docs/PLAN.md:1367) already
// requires the Go-to-front-end event push. This file 兑现 that text; it changes no
// word of it.
//
// WHAT THIS IS NOT: the page half. The renderer's receiver for
// window.wispOnPanelSnapshot does not exist yet (frontend/src/main.tsx renders
// without a snapshot and frontend/src/lib/panel.ts declares only page-to-Go doors),
// so until the machine owner's front-end agent lands that receiver a push here is
// a call into a page that ignores it. The exact contract - global name, argument
// shape, trigger timing - is written up in .scratch/wisp/probes/329/r1/page-contract.md
// for handoff. Nothing here may be read as "the screen already moves".
//
// NIL MEANS NIL, LOUDLY. A console `wisp run` has no panel window: rt.snapHost
// stays nil, every publish still books, and no push is attempted. The resident
// process attaches its one PanelManager (cmd/wisp/resident_windows.go), whose half
// of the hop posts the Eval to the panel thread and does nothing while no window
// is open - the packet stays booked in the run, which is exactly the pre-ticket-329
// state for that instant, never a silent second destination.

// panelSnapshotHost is the narrow host side of the hop: one method, shaped like
// the Eval the panelControl already exports. It is an interface (not *panelControl)
// because the pump runs on the gate's and the sink's goroutines while the control
// lives on the panel's own thread - whoever implements this owes that routing, and
// cmd/wisp's PanelManager.evalPanelJS is the one production implementation.
type panelSnapshotHost interface {
	// evalPanelJS routes one already-assembled JS sentence to this host's page
	// thread. It never blocks the caller and never rewrites the bytes it carries.
	evalPanelJS(js string)
}

// panelSnapshotGlobal is the named JS global a page receives a booked packet
// through. It is an outbound call site, not an inbound door: the page never asks,
// it only ever receives, which is why this name crosses no roster and no Method*
// constant exists for it.
const panelSnapshotGlobal = "window.wispOnPanelSnapshot"

// panelSnapshotEvalJS wraps the pump's own marshal bytes into the global call.
// The bytes go in verbatim - the same raw-JSON embedding the RPC resolve receipt
// at cmd/wisp/panel_control_windows.go:752 already does - because the AC#3 ruler
// for this hop is bit-for-bit equality between what the pump marshaled and what
// the host received. A re-encode here would be a second marshal, and a second
// marshal is a second truth.
func panelSnapshotEvalJS(data []byte) string {
	return panelSnapshotGlobal + "(" + string(data) + ")"
}

// setPanelSnapshotHost books the host this run pushes its booked packets to. The
// assembly root calls it exactly once, after assembleRuntime returned and before
// any task can publish (cmd/wisp/resident_windows.go); every other shape of run
// leaves the host nil.
func (rt *agentRuntime) setPanelSnapshotHost(h panelSnapshotHost) {
	if rt == nil {
		return
	}
	rt.snapHostMu.Lock()
	rt.snapHost = h
	rt.snapHostMu.Unlock()
}

// pushPanelSnapshotToHost is ticket 329's hop and the first PRODUCTION reader of
// lastPanelSnapshot: read the bytes this run just booked, and hand them to the
// attached host wrapped in the named global call.
//
// It reads what was booked rather than what Publish returned for one reason - the
// booked bytes are the bytes the ledger digests, so the durable record and the
// pushed bytes pin the same packet. It is called only from publishPanelSnapshot,
// immediately after the booking, which is the same instant the run already treats
// as "a state the panel shows moved".
func (rt *agentRuntime) pushPanelSnapshotToHost() {
	if rt == nil {
		return
	}
	_, data, seen := rt.lastPanelSnapshot()
	if !seen {
		return
	}
	rt.snapHostMu.Lock()
	h := rt.snapHost
	rt.snapHostMu.Unlock()
	if h == nil {
		return
	}
	h.evalPanelJS(panelSnapshotEvalJS(data))
}
