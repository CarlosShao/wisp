//go:build windows

package main

// Ticket 319 AC#1/AC#3 (N1), the HEADLESS half of "a panel thread that died is
// not the end of this process's panel".
//
// The defect (ticket 319, root cause fixed by leg 319-a1: .scratch/wisp/probes/
// 319/a1/census.md section 1): the panel window being closed - the owner's X, a
// system SC_CLOSE, or host RequestDispose, which is the same road - runs
// WM_CLOSE -> DestroyWindow -> WM_DESTROY -> PostQuitMessage(0), the library
// pump Run() returns, the thread's loop() returns, rp.finished closes, and every
// later gesture is refused by post() forever. The 2026-10-14 live log has both
// halves: "panel thread ending / why: the library pump returned" at 15:41:45,
// then "panel gesture could not be handed over" for every tray/ball gesture
// afterwards.
//
// This file pins BOTH sides of the fix. The raw half - a gesture handed to the
// dead thread - must stay refused (that is today's only path, verbatim, and the
// refusal itself is not the bug). The holder half - the same gesture through the
// assembly root's panelHolder - must rebuild: a fresh thread on the same
// manager. The rebuild is where the nail has teeth: delete the isFinished()
// branch in panelHolder.toggle and the holder half goes red on the spot, both on
// "toggle returned true" and on "the spawn was asked a second time".
//
// Ticket 319-r3 adds the D1 half (N4, AC#2/AC#3): with the manager believing its
// window is on screen, the production second-launch hook must take the SHOW road
// and never the toggle/hide road, and panelHolder.show must rebuild a dead
// thread exactly the way toggle does. Ticket 319-r4 adds the raise half to the
// same pins: the activation must also ask for the named RequestRaise hop behind
// the show (and the rebuilt thread must get both), because a visible panel that
// is pushed behind another window only comes forward through that request. All
// of it is below.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
)

// panel319DataPath is the WebView2 user-data folder this file's fixtures point
// at. Nothing here creates a window (the headless half never asks bringUp), so
// the folder is a name and not a profile - kept unique-ish per go test run only
// where it matters, and shared here on purpose: two fixtures in one run would
// otherwise race over one WebView2 profile if a later edit ever gave them a
// window.
func panel319DataPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(os.TempDir(), "wisp-319r1-holder-profile")
}

// newPanelFixture319 builds a production manager and the production panel
// thread, with NO window ever requested (this is the shape every resident boot
// starts in).
func newPanelFixture319(t *testing.T) (*residentPanel, *PanelManager) {
	t.Helper()
	disp := &panel.ComposerDispatch{Mode: &recordingModeHandler{}}
	mgr := NewPanelManager(disp, nil, panel319DataPath(t))
	rp := startResidentPanel(observe.NewRegistry(), mgr)
	t.Cleanup(rp.stop)
	return rp, mgr
}

// newDeadPanelThread319 ends the fixture's thread through its own stop() and
// returns it only once isFinished() is true - the observable state a closed
// window leaves behind. stop() is used and not a real SC_CLOSE because this
// file is the headless tier (no window may exist here); the live replay of the
// close chain itself lives in panel_death_replay_319_winlive_test.go, and what
// post() sees is the same fact either way: the thread's goroutine has returned.
func newDeadPanelThread319(t *testing.T) *residentPanel {
	t.Helper()
	rp, _ := newPanelFixture319(t)
	rp.stop()
	if !rp.isFinished() {
		t.Fatalf("the fixture's thread did not finish after stop(); every assertion below would be about a live thread")
	}
	return rp
}

// TestDeadPanelThreadRefusesEveryGesture319 is the "today, unfixed" reading at
// the headless level, kept as a case so the refusal that IS correct (a dead
// thread takes no request) can never be quietly replaced by a delivery that
// pretends to have happened. It goes red if post() ever answers true for a
// thread whose goroutine has returned.
func TestDeadPanelThreadRefusesEveryGesture319(t *testing.T) {
	rp := newDeadPanelThread319(t)
	toggle := rp.RequestToggle("319-dead-thread-toggle")
	show := rp.RequestShow("319-dead-thread-show")
	if toggle || show {
		t.Errorf("a gesture handed to a dead panel thread was accepted (toggle=%v show=%v) - post() must refuse a thread whose goroutine has returned, or the caller reports a window move that never happened", toggle, show)
	}
	if rp.isFinished() != true {
		t.Errorf("isFinished() = false on a thread whose goroutine returned")
	}
	t.Logf("319 dead-thread raw reading: RequestToggle=%v RequestShow=%v isFinished=%v", toggle, show, rp.isFinished())
}

// TestPanelHolderRebuildsADeadThread319 is the holder half this file's header
// names: the same gesture the case above hands to the dead thread, handed
// instead to the assembly root's panelHolder, must come back taken, with a
// fresh thread on the same manager behind it.
//
// Mutation control (the teeth): drop the isFinished() term from
// panelHolder.toggle's rebuild condition (cmd/wisp/resident_windows.go) and this
// case goes red on both judgements the header lists - "toggle returned true" and
// "the spawn was asked a second time" - because the dead thread behind the
// holder keeps refusing every request.
//
// The spawn seam is the observation point. Call 1 goes through the production
// startResidentPanel (the eager first thread every resident boot starts) and is
// stopped into the dead state the field leaves behind; only the REBUILD's thread
// is a double (inertPanelThread319 below), so this tier still creates no window.
// What the double cannot prove - that a rebuilt thread really brings a window up
// - is panel_death_replay_319_winlive_test.go's second case, on a real desktop.
func TestPanelHolderRebuildsADeadThread319(t *testing.T) {
	reg := observe.NewRegistry()
	mgr := NewPanelManager(&panel.ComposerDispatch{Mode: &recordingModeHandler{}}, nil, panel319DataPath(t))
	h := newPanelHolder(reg, mgr)

	productionSpawn := h.spawn
	var (
		spawns  int
		sameMgr bool
		fresh   *residentPanel
	)
	h.spawn = func(r *observe.Registry, m *PanelManager) *residentPanel {
		spawns++
		if spawns == 1 {
			return productionSpawn(r, m)
		}
		sameMgr = m == mgr
		fresh = inertPanelThread319(r, m)
		return fresh
	}
	h.start()
	if spawns != 1 || h.rp == nil {
		t.Fatalf("the holder did not start its eager first thread through the seam (spawns=%d, holder thread started=%v)", spawns, h.rp != nil)
	}
	dead := h.rp
	dead.stop()
	if !dead.isFinished() {
		t.Fatalf("the first thread did not finish after stop(); every assertion below would be about a live thread")
	}

	// Both roads, side by side: the dead thread itself still refuses this gesture
	// (the case above), the holder takes it (the fix).
	if dead.RequestToggle("319-holder-dead-thread") {
		t.Errorf("a gesture was accepted by the dead thread itself - the refusal is the honest half of the defect and must stay")
	}
	toggle := h.toggle("319-holder-rebuild")
	if !toggle {
		t.Errorf("panelHolder.toggle returned false for a gesture after the panel thread had died - the dead thread's refusal reached the caller again, which is the 2026-10-14 field defect (every tray/ball gesture refused after 15:41:45); a dead thread must make the holder start a fresh one for this gesture")
	}
	if spawns != 2 {
		t.Errorf("the spawn was asked %d time(s) for one gesture after a dead thread, want exactly 2 (the eager first thread plus one rebuild)", spawns)
	}
	if h.rp == dead {
		t.Errorf("the holder hands out the dead thread again - the rebuild replaced nothing, so the next gesture walks the same road")
	}
	if spawns == 2 {
		// Everything below is about the rebuilt thread, so it is judged only once
		// the case above has established that there is one: without that gate a
		// missing rebuild would report six reds for one defect.
		if !sameMgr {
			t.Errorf("the rebuilt thread was asked for on a manager other than the one the boot assembled; the window an operator sees comes from the manager's state, so a second manager would be a second panel")
		}
		if fresh == nil {
			t.Fatalf("the rebuild's spawn returned no thread at all")
		}
		if fresh.isFinished() {
			t.Errorf("the rebuilt thread is already finished when the gesture is handed to it")
		}
		if n := len(fresh.tasks); n != 1 {
			t.Errorf("this gesture did not land in the rebuilt thread's own queue (it holds %d request(s), want the one toggle just took) - the holder answered yes with nowhere for the request to be", n)
		}
	}
	queued := -1
	if fresh != nil {
		queued = len(fresh.tasks)
	}
	t.Logf("319 holder rebuild reading: dead thread refused its own gesture; holder toggle=%v; spawns=%d; rebuilt thread queued requests=%d; same manager=%v",
		toggle, spawns, queued, sameMgr)
}

// inertPanelThread319 is what this headless tier lets the holder rebuild into:
// the production residentPanel type wired to the production manager, with its
// request queue open and no goroutine behind it. A gesture taken by it lands in
// that queue exactly the way a rebuilt thread's first gesture does - post()
// answers true and the closure can be seen waiting - and nothing runs it, so
// this file judges the holder's decision without ever asking bringUp for a
// window, which is this tier's contract (panel319DataPath above).
//
// isFinished() is false on it for the same reason it is false on a thread that
// just started: it reads the registry handle, and this double holds none.
// finished is closed up front instead, so a later edit that wires holder.stop()
// into this case reads "already gone" rather than waiting out
// panelThreadExitBudget.
func inertPanelThread319(reg *observe.Registry, mgr *PanelManager) *residentPanel {
	rp := &residentPanel{
		mgr:     mgr,
		reg:     reg,
		tasks:   make(chan func(), 16),
		entered: make(chan struct{}),
		stopNow: make(chan struct{}),
	}
	rp.finished = make(chan struct{})
	close(rp.finished)
	return rp
}

// TestSecondLaunchActivationDoesNotTakeTheHideRoad319 is N4's headless half
// (ticket 319-r3, the D1 fix; extended with the raise half by 319-r4): the
// production second-launch hook, driven against a holder whose manager believes
// its window is on screen, must take the SHOW road - and, since 319-r4, also the
// RAISE road behind it. The defect the first half pins is the one 319-v measured
// live on 2026-10-15 (shape A in .scratch/wisp/probes/319/v/03-ac2-live-adjudication.txt):
// the hook used to drive panelHolder.toggle, and toggle's contract on a visible
// panel is to HIDE it - a double click made the window disappear. The defect the
// second half covers is 319-r3's own stopped item: the show road alone re-asserts
// visibility of a visible panel without moving the foreground, so a second
// launch on an already-visible panel that is pushed behind another window did
// nothing (the 21:40:38 field shape).
//
// How a headless tier can hold a "visible" state at all: the state the executor
// branches on is the manager's own (IsShown() = created && shown). No window
// exists here - the fixture's thread is inertPanelThread319, which runs nothing -
// so the two booleans are set directly under the manager's lock, and what is
// judged is the executor's DECISION: which requests the activation hands to the
// thread.
//
// The control is what makes the hides assertion non-vacuous: the toggle road is
// driven first on the same state and this tier must SEE its hide (hides=1). Only
// then is the hook's own reading meaningful.
//
// Mutation control: make secondLaunchActivation319 drive h.toggle again (the D1
// regression) and this case goes red on three named judgements - "did not take
// the show road", "took the toggle road" and "asked the panel thread to HIDE".
// Drop panelHolder.show's RequestRaise call (the 319-r4 raise-hop regression) and
// it goes red on "never asked for the raise hop" and on the queue count.
func TestSecondLaunchActivationDoesNotTakeTheHideRoad319(t *testing.T) {
	reg := observe.NewRegistry()
	mgr := NewPanelManager(&panel.ComposerDispatch{Mode: &recordingModeHandler{}}, nil, panel319DataPath(t))
	h := newPanelHolder(reg, mgr)
	var (
		spawns int
		live   *residentPanel
	)
	h.spawn = func(r *observe.Registry, m *PanelManager) *residentPanel {
		spawns++
		live = inertPanelThread319(r, m)
		return live
	}
	h.start()
	if live == nil || h.rp != live {
		t.Fatalf("the holder did not start its eager first thread through the seam (spawns=%d, holder thread started=%v)", spawns, h.rp != nil)
	}

	// The state D1 is about, set directly because this tier creates no window.
	mgr.mu.Lock()
	mgr.created = true
	mgr.shown = true
	mgr.mu.Unlock()
	if !mgr.IsShown() {
		t.Fatalf("this tier failed to model a visible panel (mgr.IsShown() = false), so the toggle road's hide branch would not even be reachable and the assertions below would be about nothing")
	}

	// The control: the toggle road really does hide a visible panel, and this
	// tier's instrument can see it.
	if !h.toggle("319-d1-control-toggle") {
		t.Fatalf("the toggle control was refused by the holder around a live thread; the control reading would be about nothing")
	}
	if got := live.hides.Load(); got != 1 {
		t.Errorf("the toggle control did not take the hide road (hides=%d, want 1): the instrument this case rests on does not read hides, so its later hides reading would be vacuous", got)
	}
	if got := live.shows.Load(); got != 0 {
		t.Errorf("the toggle control took the show road (shows=%d) while the manager reported the panel visible; D1's premise is that toggle HIDES in this state", got)
	}

	// The fix under the same state: the production hook.
	before := [4]int64{live.shows.Load(), live.toggles.Load(), live.hides.Load(), live.raises.Load()}
	receipt := secondLaunchActivation319(h, nil)()
	after := [4]int64{live.shows.Load(), live.toggles.Load(), live.hides.Load(), live.raises.Load()}

	if after[0]-before[0] != 1 {
		t.Errorf("the second-launch activation did not take the show road (shows went %d -> %d, want exactly one more): D1's fix is that the activation requests a SHOW, the one request whose meaning is 'this window should be on screen and in front'", before[0], after[0])
	}
	if after[1] != before[1] {
		t.Errorf("the second-launch activation took the toggle road (toggles went %d -> %d): with a visible panel that road HIDES it, which is the D1 defect 319-v measured live", before[1], after[1])
	}
	if after[2] != before[2] {
		t.Errorf("the second-launch activation asked the panel thread to HIDE (hides went %d -> %d): a second launch on an already-visible panel must never hide it (ticket 319 AC#2)", before[2], after[2])
	}
	if after[3]-before[3] != 1 {
		t.Errorf("the second-launch activation never asked for the raise hop (raises went %d -> %d, want exactly one more): a visible panel that is pushed behind another window only comes forward through RequestRaise (ticket 319 AC#2's field shape, leg 319-r4)", before[3], after[3])
	}
	if n := len(live.tasks); n != 3 {
		t.Errorf("the requests did not land in the thread's own queue (it holds %d, want 3: the control's hide plus the activation's show and raise): the executor answered with a request the thread never received", n)
	}
	if !strings.Contains(receipt, "second-launch") {
		t.Errorf("the activation hook answered %q, which is not the hand-off branch sentence; the counters above say the request was taken, so the receipt and the request disagree", receipt)
	}
	t.Logf("319-r3/r4 headless visible-panel reading: control toggle hides=%d shows=%d; activation receipt=%q; after activation: shows=%d toggles=%d hides=%d raises=%d; queued requests=%d; spawns=%d",
		before[2], before[0], receipt, after[0], after[1], after[2], after[3], len(live.tasks), spawns)
}

// TestPanelHolderShowRebuildsADeadThread319 is N4's second headless half: the
// show executor carries the same dead-thread rebuild toggle has. A panel thread
// that ended (AC#1's field shape) must not make the second launch impossible -
// the activation must start a fresh thread and hand IT the show request, and
// (since 319-r4) the raise request behind it.
//
// Mutation control: drop the isFinished() term from panelHolder.show's rebuild
// condition (cmd/wisp/resident_windows.go) and this case goes red on "the spawn
// was asked 1 time(s)" and on the dead thread's own refusal reaching the caller.
// The copy of that term in toggle is NOT exercised here: delete it instead and
// TestPanelHolderRebuildsADeadThread319 above goes red, which is why the two
// executors keep two separate copies with two separate teeth. The raise-hop
// mutation (RequestRaise dropped from show, tested by -overlay in the 319-r4
// evidence) reddens the raises and queue assertions below.
func TestPanelHolderShowRebuildsADeadThread319(t *testing.T) {
	reg := observe.NewRegistry()
	mgr := NewPanelManager(&panel.ComposerDispatch{Mode: &recordingModeHandler{}}, nil, panel319DataPath(t))
	h := newPanelHolder(reg, mgr)
	productionSpawn := h.spawn
	var (
		spawns  int
		sameMgr bool
		fresh   *residentPanel
	)
	h.spawn = func(r *observe.Registry, m *PanelManager) *residentPanel {
		spawns++
		if spawns == 1 {
			return productionSpawn(r, m)
		}
		sameMgr = m == mgr
		fresh = inertPanelThread319(r, m)
		return fresh
	}
	h.start()
	if spawns != 1 || h.rp == nil {
		t.Fatalf("the holder did not start its eager first thread through the seam (spawns=%d, holder thread started=%v)", spawns, h.rp != nil)
	}
	dead := h.rp
	dead.stop()
	if !dead.isFinished() {
		t.Fatalf("the first thread did not finish after stop(); every assertion below would be about a live thread")
	}

	receipt := secondLaunchActivation319(h, nil)()
	if spawns != 2 {
		t.Errorf("the spawn was asked %d time(s) for one activation after a dead thread, want exactly 2 (the eager first thread plus one rebuild): the show executor's refusal reached the caller again, which is the 2026-10-14 field defect", spawns)
	}
	if h.rp == dead {
		t.Errorf("the holder hands out the dead thread again - the rebuild replaced nothing, so the next activation walks the same road")
	}
	if fresh == nil {
		t.Fatalf("the rebuild's spawn returned no thread at all")
	}
	if !sameMgr {
		t.Errorf("the rebuilt thread was asked for on a manager other than the one the boot assembled; the window an operator sees comes from the manager's state, so a second manager would be a second panel")
	}
	if got := fresh.shows.Load(); got != 1 {
		t.Errorf("the rebuilt thread did not get the show request (shows=%d, want 1): the executor answered yes with nowhere for the request to be", got)
	}
	if got := fresh.toggles.Load(); got != 0 {
		t.Errorf("the rebuilt thread got a toggle request (toggles=%d, want 0): the activation must not take the hide-capable road", got)
	}
	if got := fresh.hides.Load(); got != 0 {
		t.Errorf("the rebuilt thread was asked to HIDE (hides=%d, want 0): a second launch never hides", got)
	}
	if got := fresh.raises.Load(); got != 1 {
		t.Errorf("the rebuilt thread did not get the raise request (raises=%d, want 1): without it the rebuilt window comes back without the in-front hop the second launch is for (ticket 319 AC#2)", got)
	}
	if n := len(fresh.tasks); n != 2 {
		t.Errorf("this activation did not land in the rebuilt thread's own queue (it holds %d request(s), want the show and the raise just taken)", n)
	}
	if !strings.Contains(receipt, "second-launch") {
		t.Errorf("the activation hook answered %q, which is not the hand-off branch sentence; the counters above say the request was taken, so the receipt and the request disagree", receipt)
	}
	t.Logf("319-r3/r4 headless dead-thread show reading: dead thread finished=true; spawns=%d; rebuilt thread shows=%d toggles=%d hides=%d raises=%d queued=%d same manager=%v; receipt=%q",
		spawns, fresh.shows.Load(), fresh.toggles.Load(), fresh.hides.Load(), fresh.raises.Load(), len(fresh.tasks), sameMgr, receipt)
}
