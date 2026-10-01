//go:build windows

package main

// Ticket 33's resident-host rulers (leg 33-r5). These are the three things the
// earlier legs could not say, each with its own instrument:
//
//   AC#13 - what the user ENDS UP LOOKING AT. The cold path used to show a
//           round-trip probe page last, so every cold start finished on a stub.
//           The ruler here asks the live document, through the page's own mouth,
//           whether the embedded entry's elements are present.
//   AC#14 - the Go -> page hop, in TWO separate nails, because 33-p1 §A measured
//           them arriving in different shapes: (i) the reply to an awaited JS
//           binding, which only the library's Run() delivers, and (ii) a Go-side
//           Eval push, which arrives even without it. One instrument may not
//           stand in for both (orchestrator ruling P2).
//   P1/P3 - the dedicated panel thread: STA by name, the library pump, and an
//           exit that can be OBSERVED rather than hoped for.
//
// Per P2, three readings that are all true today and still prove nothing are
// deliberately NOT used as evidence anywhere below: "Dispatch was called", "Eval
// returned no error", "the Go side's done channel closed". Every verdict here is
// built from a string the PAGE sent back.
//
// Real windows, default windows tier like the rest of this family: a machine that
// cannot create one is red, not skipped. winlive (the quiet-desktop tier, no CI
// job) owns only the timing-sensitive exit clause, see
// panel_host_windows_live_test.go.

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
	"golang.org/x/sys/windows"
)

// panelThreadWait bounds every "did the panel thread get there" wait in this file.
// A monotonic deadline (time.After / time.Since on a monotonic reading), never a
// wall-clock difference - ban #4.
const panelThreadWait = 15 * time.Second

// startPanelForTest brings up the production thread with a production-shaped
// manager (real embedded assets, a real user-data folder under the system temp
// dir, the real fail-closed dispatch chain) and returns it with the cleanups
// already registered.
func startPanelForTest(t *testing.T) (*residentPanel, *PanelManager) {
	t.Helper()
	dataPath := filepath.Join(os.TempDir(), "wisp-33r5-panel-profile")
	disp := &panel.ComposerDispatch{Mode: &recordingModeHandler{}}
	mgr := NewPanelManager(disp, builtinAssetsOrTestNil(t), dataPath)
	rp := startResidentPanel(observe.NewRegistry(), mgr)
	t.Cleanup(rp.stop)
	return rp, mgr
}

// builtinAssetsOrTestNil reads the embed through the same Go API the product uses.
// A tree without a built bundle yields nil, and each caller says out loud what it
// did with that (an un-built tree is not a green AC#13).
func builtinAssetsOrTestNil(t *testing.T) *panel.Assets {
	t.Helper()
	assets, err := panel.BuiltinAssets()
	if err != nil {
		t.Logf("panel.BuiltinAssets: %v (this tree carries no page bundle; AC#13's entry assertion names this shape below)", err)
		return nil
	}
	return assets
}

// waitPanelTrue polls a predicate the panel thread can make true, with a monotonic
// deadline. A timeout is reported as a FAILED MEASUREMENT (t.Fatalf), because "the
// thread never got there" and "the assertion holds" are different facts.
func waitPanelTrue(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(panelThreadWait)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out after %v waiting for %s - this is a failed measurement, not a pass", panelThreadWait, what)
}

// showAndWait drives the product's own opening path (a gesture posts a show
// request; the panel thread creates the window on itself) and returns only once the
// thread has finished Show - which is the state the product ends in, not the moment
// the HWND first appears: bringUp marks the window created before the cold round
// trip and the entry hand-off are done, so waiting on IsCreated alone would race
// with the rest of Show (a 33-r5 finding, not a hypothesis: the first run of this
// file's thread test read created=true / shown=false).
func showAndWait(t *testing.T, rp *residentPanel) {
	t.Helper()
	if !rp.RequestShow("test-harness") {
		t.Fatalf("RequestShow was refused by the panel thread (startUp error: %v)", rp.startUpErr())
	}
	waitPanelTrue(t, "the panel thread to finish Show", func() bool {
		return rp.mgr.IsShown() || rp.startUpErr() != nil
	})
	if err := rp.startUpErr(); err != nil {
		t.Fatalf("the panel thread refused to bring up the window: %v", err)
	}
	if !rp.mgr.IsShown() {
		t.Fatalf("the panel thread reported success without showing the window (created=%t)", rp.mgr.IsCreated())
	}
}

// editorOnPanelThread puts a window of THIS process in front of the panel, pumped by
// the panel thread itself (its Run() dispatches messages for every window on the
// thread, which is exactly the product's arrangement). AC#4's bookkeeping needs a
// prior that is not the panel, and a test binary that never received input has no
// such window until it makes one.
func editorOnPanelThread(t *testing.T, rp *residentPanel) windows.HWND {
	t.Helper()
	var (
		hwnd windows.HWND
		err  error
	)
	done := make(chan struct{})
	if !rp.post(func() {
		hwnd, err = createEditorWindow("wisp 33r5 prior window")
		if hwnd != 0 {
			pnlShowWindow.Call(uintptr(hwnd), pnlSwShow)
			pnlUpdateWindow.Call(uintptr(hwnd))
			pnlSetForeground.Call(uintptr(hwnd))
		}
		close(done)
	}) {
		t.Fatalf("the panel thread refused the editor request")
	}
	select {
	case <-done:
	case <-time.After(panelThreadWait):
		t.Fatalf("the panel thread never built the editor window within %v", panelThreadWait)
	}
	if err != nil {
		t.Skipf("the ruler could not build its own prior window (failed measurement, not a product verdict): %v", err)
	}
	if hwnd == 0 {
		t.Fatalf("CreateWindowExW returned 0 with no error code")
	}
	t.Cleanup(func() { destroyEditorWindow(hwnd) })
	return hwnd
}

// evalOnPanelThread runs a JS snippet in the live document ON the panel thread (a
// foreign goroutine calling into the control is exactly the threading bug this leg
// exists to avoid) and returns once the pump has taken the request.
//
// Eval returns no value, and that is on purpose in this file's argument: "Eval did
// not error" is one of the three readings ruling P2 refuses as evidence. Waiting
// here only proves the request reached the owning thread; whether the page saw
// anything is answered by the page, below.
func evalOnPanelThread(t *testing.T, rp *residentPanel, js string) {
	t.Helper()
	done := make(chan error, 1)
	posted := rp.post(func() {
		w := rp.mgr.currentWindow()
		if w == nil {
			done <- fmt.Errorf("no window on the panel thread")
			return
		}
		w.Eval(js)
		done <- nil
	})
	if !posted {
		t.Fatalf("the panel thread refused the Eval request (startUp error: %v)", rp.startUpErr())
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Eval on the panel thread: %v", err)
		}
	case <-time.After(panelThreadWait):
		t.Fatalf("the panel thread never ran the Eval request within %v", panelThreadWait)
	}
}

// ---------------------------------------------------------------------------
// AC#13
// ---------------------------------------------------------------------------

var entryIDRe = regexp.MustCompile(`id=["']([A-Za-z0-9_:\-.]{1,64})["']`)

// entryIDProbes extracts element ids the live document can be asked about. The ids
// come out of the SAME Go-side resolve the product serves from (panel.Assets), so
// this ruler never opens a frontend file, and the probe list is what makes the
// assertion about content rather than about a call being made.
func entryIDProbes(t *testing.T) []string {
	t.Helper()
	assets, err := panel.BuiltinAssets()
	if err != nil {
		return nil
	}
	data, _, err := assets.Resolve(panel.EntryFile)
	if err != nil {
		t.Logf("Resolve(entry): %v", err)
		return nil
	}
	seen := map[string]bool{}
	var out []string
	for _, m := range entryIDRe.FindAllStringSubmatch(string(data), -1) {
		id := m[1]
		if seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
		if len(out) >= 8 {
			break
		}
	}
	t.Logf("AC#13 probes from the resolved entry (%d bytes): %d id(s) %v", len(data), len(out), out)
	return out
}

// TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe is ticket 33 AC#13's ruler:
// after a real cold start the document in the window must be the embedded page,
// not the round-trip probe the host shows while it checks the message channel.
//
// The question is put to the page and answered by the page: Go asks the live
// document, id by id, whether the entry's own elements exist, and the document's
// answer comes back through a Go-side binding. That shape is what gives this
// instrument teeth against the exact defect (a SetHtml ordering change flips the
// answer from every-id-present to no-id-present); an assertion like "SetHtml was
// called" or "the entry bytes were read" could not tell the two orders apart.
func TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe(t *testing.T) {
	probes := entryIDProbes(t)
	if len(probes) == 0 {
		t.Skipf("AC#13 has no subject in this tree: the embed resolves no entry, so there is no page content to be covered. Named skip; the working-tree reading is the one that ran in docs/evidence/s1/33-panel-host-c27-r5.md")
	}
	rp, mgr := startPanelForTest(t)
	showAndWait(t, rp)

	reports := make(chan string, 4)
	bindReceiptOnThread(t, rp, panelProbeReportBinding, reports)
	recreateDocumentOnThread(t, rp, mgr)

	js := probeJS(panelProbeReportBinding, probes)
	evalOnPanelThread(t, rp, js)

	var answer string
	select {
	case answer = <-reports:
	case <-time.After(panelThreadWait):
		t.Fatalf("the page never answered the AC#13 probe within %v - the document in the window reported nothing, which is also what a stuck or blank page does", panelThreadWait)
	}
	present := 0
	for _, b := range answer {
		if b == '1' {
			present++
		}
	}
	t.Logf("AC#13 page answer (head %s): %q of %d probe id(s) present in the live document", gitHeadShortForTest(t), answer, len(probes))
	if present == 0 {
		t.Errorf("after a real cold start the live document contains NONE of the %d element ids the embedded entry declares (page said %q). That is the shape ticket 33 AC#13 names: the round-trip probe page is the last document shown, so the user sees a stub instead of the panel. Product side: bringUp must serve the entry AFTER the probe, and the probe stays because it is where cold 'usable' is decided", len(probes), answer)
	}
}

const panelProbeReportBinding = "wisp33r5Report"

// bindReceiptOnThread registers a one-string-argument binding that forwards what
// the page hands it to a Go channel. Bind runs on the panel thread for the same
// reason Eval does.
func bindReceiptOnThread(t *testing.T, rp *residentPanel, name string, sink chan string) {
	t.Helper()
	done := make(chan error, 1)
	if !rp.post(func() {
		w := rp.mgr.currentWindow()
		if w == nil {
			done <- fmt.Errorf("no window on the panel thread")
			return
		}
		done <- w.Bind(name, func(text string) string {
			select {
			case sink <- text:
			default:
			}
			return "seen"
		})
	}) {
		t.Fatalf("the panel thread refused the Bind request (startUp error: %v)", rp.startUpErr())
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Bind %q on the panel thread: %v", name, err)
		}
	case <-time.After(panelThreadWait):
		t.Fatalf("the panel thread never ran the Bind request within %v", panelThreadWait)
	}
}

// recreateDocumentOnThread re-serves the entry document. It is needed because a
// JS binding is injected into documents created AFTER the Bind call, and this
// binding is a ruler, not part of the product. Re-serving the entry keeps AC#13's
// subject intact: the document under test is still the embedded page.
func recreateDocumentOnThread(t *testing.T, rp *residentPanel, mgr *PanelManager) {
	t.Helper()
	done := make(chan struct{})
	if !rp.post(func() {
		if err := mgr.serveEntry(); err != nil {
			t.Errorf("re-serving the entry for the ruler: %v", err)
		}
		close(done)
	}) {
		t.Fatalf("the panel thread refused the re-serve request")
	}
	select {
	case <-done:
	case <-time.After(panelThreadWait):
		t.Fatalf("the panel thread never re-served the document within %v", panelThreadWait)
	}
}

// probeJS asks the live document, one by one, whether the entry's elements are
// there, and reports the answer bitmap back through the given binding.
func probeJS(binding string, ids []string) string {
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, fmt.Sprintf("b += (document.getElementById(%q) ? '1' : '0');", id))
	}
	return "(function(){ var b=''; " + strings.Join(parts, " ") +
		" if (window." + binding + ") window." + binding + "(b); })();"
}

// ---------------------------------------------------------------------------
// AC#14, nail 1: the reply to an awaited JS binding
// ---------------------------------------------------------------------------

// ac14AwaitJS makes the page do the thing C17 names ("回复必须按 correlationId
// 路由"): call the real door, AWAIT its return value, and report what arrived. A
// reply that never lands shows up as NOT_RESOLVED in the page's own words, which
// is the only form of evidence ruling P2 accepts.
func ac14AwaitJS(binding string) string {
	return `(async function(){
  var out = [];
  for (var i = 0; i < 3; i++) {
    var env = '{"method":"panel.mode.request","requestId":"ac14-' + i + '","source":"panel-composer","to":"ask_every_step"}';
    var got = 'NOT_RESOLVED';
    try {
      got = await Promise.race([
        Promise.resolve(window.` + binding + `(env)).then(function(v){ return 'REPLIED:' + v; }),
        new Promise(function(res){ setTimeout(function(){ res('TIMEOUT-2S'); }, 2000); })
      ]);
    } catch (e) { got = 'THREW:' + String(e); }
    out.push(i + '=' + got);
  }
  if (window.wisp33r5ReplyReport) window.wisp33r5ReplyReport(out.join('|'));
})();`
}

func TestAC14AwaitedBindingReplyReachesThePage(t *testing.T) {
	rp, mgr := startPanelForTest(t)
	showAndWait(t, rp)

	h := mgr.disp.Mode.(*recordingModeHandler)
	reports := make(chan string, 4)
	bindReceiptOnThread(t, rp, "wisp33r5ReplyReport", reports)
	recreateDocumentOnThread(t, rp, mgr)
	evalOnPanelThread(t, rp, ac14AwaitJS(panelDispatchBinding))

	var answer string
	select {
	case answer = <-reports:
	case <-time.After(panelThreadWait):
		t.Fatalf("the page never reported the AC#14 awaited values within %v. Without a report from the page there is NO evidence either way: the three shapes that are all true with no reply (Dispatch called, Eval no error, Go-side done closed) are not assertions", panelThreadWait)
	}
	t.Logf("AC#14 nail 1 (reply hop), page's own words: %q (mode handler calls seen by Go: %d)", answer, h.count())
	if !strings.Contains(answer, "REPLIED:") {
		t.Errorf("the awaited JS binding reply never reached the page: the page reported %q. Go did receive the calls (handler saw %d), so this is the H10 half, not H3: webview.Dispatch only queues a closure and posts a thread message, and only the library's Run() drains that queue (33-p1 §A, R25 vs R26). The panel thread must hand its pump to Run()", answer, h.count())
	}
	if strings.Contains(answer, "NOT_RESOLVED") || strings.Contains(answer, "TIMEOUT-2S") {
		t.Errorf("only part of the reply hop arrived: the page reported %q - every one of the three awaited calls has to resolve for C17's push-and-route to hold", answer)
	}
	if h.count() < 3 {
		t.Errorf("the page says it got replies but Go only saw %d handler call(s) of 3 - the two halves of this hop must agree", h.count())
	}
}

// ---------------------------------------------------------------------------
// AC#14, nail 2: the Go -> page Eval push (a DIFFERENT dimension)
// ---------------------------------------------------------------------------

func TestAC14GoSideEvalPushReachesThePage(t *testing.T) {
	rp, mgr := startPanelForTest(t)
	showAndWait(t, rp)

	reports := make(chan string, 4)
	bindReceiptOnThread(t, rp, "wisp33r5PushReport", reports)
	recreateDocumentOnThread(t, rp, mgr)

	// A value the page can only know if Go's push landed. This nail deliberately
	// does not await a binding: it asks "can Go put something into the document",
	// which 33-p1 R27 measured arriving even in a host whose reply hop was broken.
	const pushed = "PUSHED-33R5-OK"
	evalOnPanelThread(t, rp, fmt.Sprintf("document.title = %q;", pushed))
	evalOnPanelThread(t, rp, "if (window.wisp33r5PushReport) window.wisp33r5PushReport(document.title);")

	var answer string
	select {
	case answer = <-reports:
	case <-time.After(panelThreadWait):
		t.Fatalf("the page never reported its own title within %v - the push instrument got no answer", panelThreadWait)
	}
	t.Logf("AC#14 nail 2 (Eval push hop), page's own words: title=%q", answer)
	if answer != pushed {
		t.Errorf("Go's Eval push did not reach the document: title is %q, want %q. This is the push dimension, separate from the awaited-reply dimension asserted in TestAC14AwaitedBindingReplyReachesThePage - one arriving says nothing about the other", answer, pushed)
	}
}

// ---------------------------------------------------------------------------
// P1 / P3: the thread itself
// ---------------------------------------------------------------------------

// TestPanelThreadIsSTAAndExitsCleanly nails the two shapes ruling P1/P3 asked for:
// the thread asks for COINIT_APARTMENTTHREADED by name (and refuses to create a
// window otherwise), and Run() does not trap the process - stop() ends it, and the
// end is observable.
func TestPanelThreadIsSTAAndExitsCleanly(t *testing.T) {
	rp, _ := startPanelForTest(t)
	if err := rp.startUpErr(); err != nil {
		t.Fatalf("the panel thread failed before it was asked for anything: %v", err)
	}
	if rp.isCreated() {
		t.Errorf("starting the panel thread created a window: the panel must appear when the user asks (a gesture or a test), not at boot")
	}
	showAndWait(t, rp)
	hwnd := rp.mgr.windowHandle()
	if hwnd == 0 {
		t.Fatalf("created but HWND is 0")
	}
	if !rp.mgr.IsShown() {
		t.Errorf("IsShown false right after the show request ran on the panel thread")
	}
	rp.stop()

	select {
	case <-rp.finished:
	case <-time.After(panelThreadExitBudget + 5*time.Second):
		t.Fatalf("the panel thread did not exit after stop(): Run() has no way out in this shape, which is exactly what ruling P1 required a shipping path to have")
	}
	if rp.mgr.IsCreated() {
		t.Errorf("the panel thread exited with the window still marked created - teardown did not run on the owning thread")
	}
	if rp.isFinished() != true {
		t.Errorf("the registry handle reports the thread still running after its finished channel closed")
	}
	t.Logf("panel thread up and down: hwnd 0x%x, exits observed, shows=%d", hwnd, rp.shows.Load())
}

// TestPanelThreadNameIsNotInResidentRoster records a boundary this leg was told
// not to cross: adding a resident goroutine name, or moving the baseline, is a
// contract change (D38b). The panel thread is deliberately unlisted, which makes
// it CategoryUnknown and nothing else.
func TestPanelThreadNameIsNotInResidentRoster(t *testing.T) {
	for _, n := range observe.ResidentNames {
		if n == panelSTAName {
			t.Errorf("panel thread name %q is now in observe.ResidentNames - that roster is frozen at six names (D38b) and adding to it is a human-approved contract change, not something this leg may do by wiring a thread", panelSTAName)
		}
	}
	if got := observe.ClassifyGoroutine(panelSTAName); got != observe.CategoryUnknown {
		t.Errorf("ClassifyGoroutine(%q) = %q, want %q: the panel thread must stay outside the frozen resident roster", panelSTAName, got, observe.CategoryUnknown)
	}
}

// TestBallPanelGesturesReachThePanelThread is AC#1's "the user can open it" half at
// the seam this process actually has: the two panel gestures carry an executor now,
// and that executor hands over to the panel thread instead of running Win32 on the
// ball's ui-sta.
func TestBallPanelGesturesReachThePanelThread(t *testing.T) {
	rp, mgr := startPanelForTest(t)
	hs := panelHostHooks{showPanel: rp.RequestToggle}
	hs.requestPanelOpen("panel-hotkey")()
	waitPanelTrue(t, "the gesture's show request to end in a shown window", func() bool {
		return mgr.IsShown() || rp.startUpErr() != nil
	})
	if !mgr.IsShown() {
		t.Errorf("the panel gesture arrived but no window is shown (created=%t)", mgr.IsCreated())
	}
	if got := rp.toggles.Load(); got < 1 {
		t.Errorf("the panel gesture reached the host but the thread counted %d toggle(s)", got)
	}
	// Same gesture again: it now hides the panel and hands focus back.
	hs.requestPanelOpen("tray-open-panel")()
	waitPanelTrue(t, "the second gesture to hide the panel", func() bool { return !mgr.IsShown() })

	if hs2 := (panelHostHooks{}); hs2.showPanel == nil {
		// The pre-ticket-33 shape: with no executor the gesture must still be
		// recorded, which is what keeps a silent nil callback from coming back.
		hs2.requestPanelOpen("panel-hotkey")
	}
}

// TestPanelHostHooksWithoutExecutorStillRecord covers the nil-executor branch
// named in requestPanelOpen: a boot that could not assemble the panel must say so
// per gesture, not swallow the click.
func TestBallGestureWithoutPanelHostStillRecords(t *testing.T) {
	hs := panelHostHooks{}
	if hs.showPanel != nil {
		t.Fatalf("this fixture is supposed to have no executor")
	}
	fn := hs.requestPanelOpen("panel-hotkey")
	if fn == nil {
		t.Fatalf("requestPanelOpen returned nil for the no-executor shape - a nil callback is the silence this file exists to prevent")
	}
	fn()
}

// TestAC4PriorFocusSurvivesARefusedPanelSample nails the two rules 33-r5 added to
// the focus bookkeeping, in the shape that is decidable WITHOUT depending on
// Windows' foreground rights (which is what makes the end-to-end ruler in
// TestAC4FocusReturnToPriorWindowGap33r5 environment-gated on assertion 3):
//
//   1. a sample that turns out to be the panel itself must not overwrite the honest
//      prior - that overwrite IS the defect 33-v1 §A#27 found;
//   2. the recorded prior survives a Hide, so the next Hide still has a target (see
//      the reasoning on PanelManager.Hide).
//
// Under the pre-33-r5 code both statements are false by construction (Show assigned
// GetForegroundWindow() unconditionally, after the window already had the focus), so
// this case has teeth against the exact change.
func TestAC4PriorFocusSurvivesARefusedPanelSample(t *testing.T) {
	rp, mgr := startPanelForTest(t)
	editor := editorOnPanelThread(t, rp)
	showAndWait(t, rp)
	hwnd := mgr.windowHandle()
	if hwnd == 0 {
		t.Fatalf("created but HWND is 0")
	}
	mgr.mu.Lock()
	first := mgr.prevFocus
	mgr.mu.Unlock()
	if uintptr(first) == uintptr(editor) {
		t.Logf("the host recorded the ruler's own editor window 0x%x as the prior - exactly the value a real user's editor would be", first)
	}
	if first == 0 {
		t.Skipf("no prior foreground was recorded on this desktop session (Show sampled 0 or the panel itself), so the survival rule has no subject here. Named skip; the handles this run saw are logged by TestAC4FocusReturnToPriorWindowGap33r5")
	}
	if uintptr(first) == hwnd {
		t.Errorf("the host recorded the panel itself as the prior foreground (0x%x) even though Show samples before creating - that is the 33-v1 §A#27 defect, and setPriorFocusLocked exists to refuse it", first)
	}
	mgr.Hide()
	mgr.mu.Lock()
	afterHide := mgr.prevFocus
	restoreTo := mgr.lastRestoreTo
	mgr.mu.Unlock()
	if uintptr(afterHide) != uintptr(first) {
		t.Errorf("Hide changed the recorded prior from 0x%x to 0x%x - with that, a later Hide has no target and focus never goes back", first, afterHide)
	}
	if uintptr(restoreTo) != uintptr(first) {
		t.Errorf("Hide attempted to restore to 0x%x but records 0x%x - the handles the AC#4 readings are built from must be the ones the call used", first, restoreTo)
	}

	// Now force the refused-sample branch: put the panel in the foreground, then ask
	// Show to record. It must keep the honest prior instead of overwriting it.
	pnlShowWindow.Call(hwnd, pnlSwShow)
	pnlSetForeground.Call(hwnd)
	pnlUpdateWindow.Call(hwnd)
	if err := mgr.Show(context.Background()); err != nil {
		t.Fatalf("second Show: %v", err)
	}
	mgr.mu.Lock()
	kept := mgr.prevFocus
	mgr.mu.Unlock()
	if uintptr(kept) == hwnd {
		t.Errorf("Show overwrote the recorded prior with the panel's own HWND (0x%x) - the refusal in setPriorFocusLocked is the fix and this assertion is its tooth", kept)
	}
	if uintptr(kept) != uintptr(first) {
		t.Errorf("Show replaced the honest prior 0x%x with 0x%x while the panel itself was foreground; the last reachable window the user came from must survive one refused sample", first, kept)
	}
}
