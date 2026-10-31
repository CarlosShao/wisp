//go:build windows

package main

import (
	"context"
	"log/slog"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/proc"
	"github.com/CarlosShao/wisp/internal/statemachine"
	"github.com/CarlosShao/wisp/internal/watchdog"
)

// Ticket 384 (the assembly seam ticket 42's verdict §7 named): the production
// wiring of internal/watchdog into the resident process. Until this file the
// package had zero production callers - the loop, the whitelist enforcement,
// the settle gate and the alert surface all existed and ran for nobody.
//
// What this file connects, one seam per paragraph:
//
//   - the REAL TreeReader: proc.NewTreeSampler over the C30 Job Object this
//     boot already opened (rt.Job) - the same reader shape `wisp slo` samples
//     with (slo_windows.go), so the watchdog gates on the same tree metrics
//     the SLO instrument measures, with no second metrics path;
//   - StateFn wired to the resident D43 machine (rv.machine, the machine the
//     voice leg dispatches): the D32 row the loop gates against follows the
//     machine's live state, per tick;
//   - OnAlert mapped onto the D43 table rows BY NAME: AlertOverrun in Armed
//     -> EvWatchdogOverrun (row #9, the Armed self-loop that stays Armed),
//     AlertFatal -> EvWatchdogFatal (row #35, AnyState -> WatchdogAlert).
//     Alerts without a table row stay log/UI material - dispatching one would
//     be an illegal transition the machine rejects by design;
//   - the panel-webview whitelist component registered against the real panel
//     host; and Start, which spawns the loop under the D38b roster name
//     "watchdog" - the slot SPEC-03 §4.3 reserves and ticket 223's config tick
//     already books on the `wisp run` side.
//
// 273 SINK STATUS, NAMED (the dispatch's first check, re-run 2026-10-31):
// ticket 273 is OPEN (AC#2..AC#6 unchecked; the three-shape cost table has
// picked no shape), so the production machines are built with Sink=nil and the
// machine's own no-op default sink receives every Effect. This leg does NOT
// invent a sink - that surface is 273's contract. What is real here: the
// transitions dispatch, the states they leave are readable via
// rv.machine.State() (WatchdogAlert included), and the mapped effects
// (kws.keep-alive-alert, ui.one-click-restart) are delivered to whatever sink
// the machine carries. When 273 lands the real receiver this file needs no
// change: it already dispatches through the machine's own Sink field.
//
// CONFIG DUTY, NAMED: Options.Config stays nil. The resident process holds no
// config.Manager (223's 1s poll arms on `wisp run`'s agentRuntime; this boot's
// config reads are per-use fresh reads), and the dispatch bans new config
// keys - so the loop runs its resource-duty half only, which is exactly the
// library's documented nil-Config shape.
func startResidentWatchdog(rt *proc.Runtime, rv *residentVoice, panel *PanelManager) (*watchdog.Loop, *observe.Root) {
	root := observe.NewRoot("watchdog-host")
	loop := watchdog.New(watchdog.Options{
		TreeReader: proc.NewTreeSampler(rt.Job),
		StateFn:    func() observe.SLOState { return sloStateFromD43Machine(rv) },
		Registry:   rt.Registry,
		OnAlert:    func(a watchdog.Alert) { dispatchWatchdogAlert(rv, a) },
	})
	if panel != nil {
		loop.Register(watchdog.ComponentPanelWebView, panelWebviewUnloadable{mgr: panel})
	}
	loop.Start(root.Ctx)
	return loop, root
}

// sloStateFromD43Machine reads the resident machine's current D43 state into
// the D32 sampling row the loop gates against. The four-row mapping is the
// correspondence observe's own doc pins (sampler.go: Sleeping/Armed/Warm/
// Conversation ARE D43 machine states, "names pinned by test to internal/
// statemachine"), so those four map by name and nothing else maps by
// invention: the machine's remaining states carry no D32 sampling row of
// their own (the six-row table samples the resident postures; PanelOpen and
// WorkPeak are SLO postures over the machine, not machine states, and nothing
// produces them here). They answer the zero SLOState, which the loop treats
// as "not a D32 sampling state; tick skipped" and logs (watchdog.go,
// enforceResources) - fail-closed, and never a session-active state gated
// against a smaller state's cap, the D32 16.3.2 revision 3 harmful shape
// (a 25MB Sleeping gate would auto-unload a session that is merely awake).
func sloStateFromD43Machine(rv *residentVoice) observe.SLOState {
	if rv == nil || rv.machine == nil {
		return observe.SLOSleeping
	}
	switch rv.machine.State() {
	case statemachine.StateSleeping:
		return observe.SLOSleeping
	case statemachine.StateArmed:
		return observe.SLOArmed
	case statemachine.StateWarm:
		return observe.SLOWarm
	case statemachine.StateConversation:
		return observe.SLOConversation
	default:
		return observe.SLOState("")
	}
}

// dispatchWatchdogAlert is the host-side D43 mapping - the half ticket 42's
// package doc reserves for the host ("the loop never touches the state
// machine ... the host owns the D43 mapping"). Two alert shapes have table
// rows and are dispatched BY NAME, row numbers and side effects verbatim:
//
//   - AlertOverrun with State==Armed -> EvWatchdogOverrun. D43 row #9
//     (table.go: From Armed, To Armed - the self-loop that STAYS Armed,
//     guard watchdog-overrun, side effect kws.keep-alive-alert). KWS is
//     structurally un-unloadable here (D32 revision 3; the whitelist never
//     carries it), so the row is the keep-alive alert, not an unload.
//   - AlertFatal -> EvWatchdogFatal. D43 row #35 (From AnyState, To
//     StateWatchdogAlert, side effect ui.one-click-restart - the one-click
//     restart affordance).
//
// Everything else - AlertOverrun raised from a state whose row set has no
// watchdog.overrun edge, and AlertLeak (D42#10) - has no D43 row: the
// mapping records it as log/UI material and dispatches nothing. Dispatching
// anyway would ask the machine for a transition the frozen table does not
// carry, which the machine rejects (undefined = stop) - the mapping does not
// manufacture rows.
//
// The dispatched effects land in the machine's sink, which is 273's surface
// (see the file header): today that is the package's no-op default, so the
// transition and the state it leaves are the visible half, and the effect
// delivery (kws.keep-alive-alert, ui.one-click-restart) becomes visible when
// 273's receiver lands.
func dispatchWatchdogAlert(rv *residentVoice, a watchdog.Alert) {
	if rv == nil || rv.machine == nil {
		return
	}
	switch {
	case a.Kind == watchdog.AlertOverrun && a.State == observe.SLOArmed:
		to, err := rv.machine.Dispatch(statemachine.EvWatchdogOverrun, nil)
		if err != nil {
			slog.Warn("watchdog host: D43 #9 dispatch rejected",
				"from", string(rv.machine.State()),
				"err", err.Error(), "error_class", string(observe.ClassInternal))
			return
		}
		slog.Warn("watchdog host: D43 #9 Armed overrun mapped",
			"event", string(statemachine.EvWatchdogOverrun),
			"to", string(to),
			"side_effect", "kws.keep-alive-alert",
			"measured", a.Measured, "limit", a.Limit, "detail", a.Detail)
	case a.Kind == watchdog.AlertFatal:
		to, err := rv.machine.Dispatch(statemachine.EvWatchdogFatal, nil)
		if err != nil {
			slog.Warn("watchdog host: D43 #35 dispatch rejected",
				"from", string(rv.machine.State()),
				"err", err.Error(), "error_class", string(observe.ClassInternal))
			return
		}
		slog.Warn("watchdog host: D43 #35 fatal mapped (WatchdogAlert + one-click restart)",
			"event", string(statemachine.EvWatchdogFatal),
			"to", string(to),
			"side_effect", "ui.one-click-restart",
			"measured", a.Measured, "limit", a.Limit, "detail", a.Detail)
	default:
		slog.Info("watchdog host: alert recorded (no D43 row; log/UI material)",
			"kind", string(a.Kind), "state", string(a.State),
			"metric", a.Metric, "measured", a.Measured, "limit", a.Limit)
	}
}

// panelWebviewUnloadable is the D32 revision 3 panel-WebView half of the
// safe-unloadable whitelist, wired to the REAL panel host. "The panel WebView
// if the user is not viewing it" reads here as: a manager whose window was
// never created (or a nil host - a boot whose panel did not assemble
// registers nothing at all) has nothing to drop; a SHOWN window may be on the
// user's retina, so the pass skips it (dropped=false - the skip branch ticket
// 42's audit test pins); and a hidden window is torn down for real. The
// teardown is PanelManager.Destroy, whose webview half is go-webview2's
// Destroy = one PostMessageW(hwnd, WM_CLOSE) - postable from any thread, so
// the watchdog goroutine may call it, with the panel thread's own pump
// running the close chain. Destroy (not Hide) is the unload semantics: a hide
// only moves the window off screen and leaves the msedgewebview2 child set
// resident, while the close chain ends the child processes and the tree's
// private bytes actually fall. The drop is recoverable by construction: a
// later Show rebuilds the window (the recreate path panel_host_windows.go
// documents) and panelHolder starts a fresh panel thread for the next
// gesture (ticket 319).
type panelWebviewUnloadable struct {
	mgr *PanelManager
}

func (p panelWebviewUnloadable) Unload(ctx context.Context) (bool, string) {
	m := p.mgr
	if m == nil {
		return false, "panel host not assembled"
	}
	if !m.IsCreated() {
		return false, "no panel window exists"
	}
	if m.IsShown() {
		return false, "panel is shown; the user may be viewing it - not dropped"
	}
	m.Destroy()
	return true, "hidden panel WebView destroyed (the next show rebuilds it)"
}
