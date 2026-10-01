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
	"runtime"
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
// How the page talks back in this file, and why it is the door that already exists.
//
// Every verdict below is built from a string the PAGE sent, per ruling P2. What the
// page does NOT need is a new binding: registering one means creating a fresh document
// after the Bind, and a ruler that re-creates the document can no longer tell what a
// COLD START ended on. 33-r5 found this out by running its own reverse control (swap
// the two SetHtml calls in bringUp) against a ruler that stayed GREEN - the re-serve
// was covering the mutation. So the reports ride window.wispDispatch, the product's own
// inbound door, as real composer envelopes the test's recording handler reads back.
// That door is bound before any page exists, so it is present in whatever document the
// host ends on, and the ruler now changes nothing about the document lifecycle.
// ---------------------------------------------------------------------------

// all returns every request the recording handler was actually reached with.
func (h *recordingModeHandler) all() []panel.ComposerRequest {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]panel.ComposerRequest(nil), h.reqs...)
}

// awaitReport waits for the page to send one report under the given requestId and
// returns the payload it carried. Silence within the budget is a failed measurement and
// is reported as one.
func awaitReport(t *testing.T, h *recordingModeHandler, requestID, what string) string {
	t.Helper()
	deadline := time.Now().Add(panelThreadWait)
	for time.Now().Before(deadline) {
		for _, req := range h.all() {
			if req.RequestID == requestID {
				return req.To
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("no report %q from the page within %v (what DID arrive at the door: %s). %s cannot be decided without the page's own answer, and the three shapes that are all true with no reply (Dispatch called, Eval returned, Go-side channel closed) are not assertions",
		requestID, panelThreadWait, describeRequests(h), what)
	return ""
}

// describeRequests renders what the door received, so a missing report can be told
// apart from a report that arrived under a different name.
func describeRequests(h *recordingModeHandler) string {
	parts := make([]string, 0, 8)
	for _, req := range h.all() {
		to := req.To
		if len(to) > 24 {
			to = to[:24] + "..."
		}
		parts = append(parts, req.RequestID+"="+to)
	}
	if len(parts) == 0 {
		return "nothing at all"
	}
	return strings.Join(parts, ", ")
}

// reportJSEnv builds the raw envelope the page posts to the product's door, with
// payload spliced into the to field.
func reportJSEnv(requestID, payload string) string {
	return "'{\"method\":\"panel.mode.request\",\"requestId\":\"" + requestID +
		"\",\"source\":\"panel-composer\",\"to\":\"' + (" + payload + ") + '\"}'"
}

// ---------------------------------------------------------------------------
// AC#13
// ---------------------------------------------------------------------------

var entryIDRe = regexp.MustCompile(`id=["']([A-Za-z0-9_:\-.]{1,64})["']`)

// entryIDProbes extracts element ids the live document can be asked about. The ids come
// out of the SAME Go-side resolve the product serves from (panel.Assets), so this ruler
// never opens a frontend file, and the probe list is what makes the assertion be about
// content rather than about a call having been made.
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

// probeJS asks the live document, id by id, whether the entry's elements are there, and
// posts the answer bitmap back through the door that is already in the page.
func probeJS(ids []string) string {
	parts := make([]string, 0, len(ids)+2)
	parts = append(parts, "var b='';")
	for _, id := range ids {
		parts = append(parts, fmt.Sprintf("b += (document.getElementById(%q) ? '1' : '0');", id))
	}
	parts = append(parts, "window.wispDispatch("+reportJSEnv("ac13-probe", "b")+");")
	return "(function(){ " + strings.Join(parts, " ") + " })();"
}

// TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe is ticket 33 AC#13's ruler: after a
// real cold start the document in the window must be the embedded page, not the
// round-trip probe the host shows while it checks the message channel.
//
// The question goes to the live document and the document answers: does it carry the
// entry's own elements? That form is what gives this instrument teeth against the exact
// defect - swapping the two SetHtml calls flips the answer from every-id-present to
// none-present - where "SetHtml was called" or "the entry bytes were read" would tell
// the two orders apart just as little as the old t.Logf did.
func TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe(t *testing.T) {
	probes := entryIDProbes(t)
	if len(probes) == 0 {
		t.Skipf("AC#13 has no subject in this tree: the embed resolves no entry, so there is no page content that could be covered. Named skip; the working-tree reading is the one that ran, recorded in docs/evidence/s1/33-panel-host-c27-r5.md")
	}
	rp, mgr := startPanelForTest(t)
	showAndWait(t, rp)
	h := mgr.disp.Mode.(*recordingModeHandler)

	evalOnPanelThread(t, rp, probeJS(probes))

	answer := awaitReport(t, h, "ac13-probe", "AC#13")
	present := strings.Count(answer, "1")
	t.Logf("AC#13 page answer (head %s): %q - %d of %d probe id(s) present in the live document",
		gitHeadShortForTest(t), answer, present, len(probes))
	if present == 0 {
		t.Errorf("after a real cold start the live document contains NONE of the %d element ids the embedded entry declares (the page itself answered %q). That is ticket 33 AC#13: the round-trip probe page is the last document shown, so the user sees a stub instead of the panel. Product side: bringUp must hand the entry over AFTER the probe - the probe stays, because it is where cold usable is decided", len(probes), answer)
	}
}

// ---------------------------------------------------------------------------
// AC#13's order-dependent red, root cause, nailed DETERMINISTICALLY (leg 33-r6)
// ---------------------------------------------------------------------------

// postQuitProc reaches the same user32 PostQuitMessage go-webview2's Terminate is
// (webview.go:381-383), so this test plants exactly the message the library plants.
var postQuitProc = pnlModUser32.NewProc("PostQuitMessage")

// TestAC13BringUpSurvivesAReusedThreadQuit is the deterministic form of the red that
// only appeared when a window-owning thread ran before AC13. The full-package shape
// is a race (the Go scheduler has to hand AC13's panel-sta a pooled thread that an
// earlier hostThreadHarness returned via UnlockOSThread while it still carried a
// WM_QUIT), so it is not a stable instrument. Here the same poison is planted on THIS
// locked thread on purpose, then the product's own bringUp runs on it.
//
// Why this is the red: go-webview2's window procedure calls w.Terminate() on
// WM_DESTROY (webview.go:242-243), and Terminate is a bare PostQuitMessage onto
// whichever thread pumped the close. A thread that hosted a now-destroyed panel
// window therefore carries a pending WM_QUIT into the next bring-up: the library's
// Embed loop (pkg/edge chromium.go:96-111) dequeues that quit from GetMessageW,
// breaks with inited still 0 and e.webview never assigned, and Init
// (chromium.go:130-136, no guard on the nil) panics with exactly the recovered string
// in the ticket's log. The fix is drainStaleQuitBeforeCreate at the top of bringUp.
//
// The reverse control: delete that one call from bringUp and this test fails on the
// planted quit (NewWithOptions panics before the window exists), reported by name -
// not by a 15-second timeout, because the panic is what is being pinned here.
func TestAC13BringUpSurvivesAReusedThreadQuit(t *testing.T) {
	dataPath := filepath.Join(os.TempDir(), "wisp-33r6-reused-thread-quit")
	disp := &panel.ComposerDispatch{Mode: &recordingModeHandler{}}
	mgr := NewPanelManager(disp, builtinAssetsOrTestNil(t), dataPath)

	type outcome struct {
		err     error
		panicked any
		created bool
	}
	res := make(chan outcome, 1)
	go func() {
		// owner: test bring-up on a deliberately reused (poisoned) thread; recover below.
		var o outcome
		defer func() {
			o.panicked = recover()
			res <- o
		}()
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		// Plant the exact pending quit a prior window owner would have left behind.
		postQuitProc.Call(0)
		o.err = mgr.bringUp(context.Background())
		// Read the verdict BEFORE teardown clears it.
		o.created = mgr.IsCreated()
		mgr.Destroy()
	}()

	o := <-res
	t.Logf("AC#13 reused-thread root cause: planted one WM_QUIT on this locked thread, then ran bringUp - panicked=%v err=%v created=%v",
		o.panicked, o.err, o.created)
	if o.panicked != nil {
		t.Errorf("bringUp on a thread carrying a pending WM_QUIT %v (the recovered panic) instead of draining it and creating the window. This is ticket 33 AC#13's order-dependent red reproduced deterministically: the library's Embed loop (pkg/edge chromium.go:96-111) dequeued the quit and Init (chromium.go:131) dereferenced the nil control. Fix: bringUp must call drainStaleQuitBeforeCreate before NewWithOptions", o.panicked)
	}
	if o.err != nil {
		t.Errorf("bringUp returned an error on the reused thread: %v (a created window is expected; with the drain in place a pre-existing quit is no reason to refuse to bring the panel up)", o.err)
	}
	if !o.created {
		t.Errorf("bringUp reported success but no window was created on the reused thread (created=%v) - the panel would be blank for the user", o.created)
	}
}

// ---------------------------------------------------------------------------
// AC#14, nail 1: the reply to an awaited JS binding
// ---------------------------------------------------------------------------

// ac14AwaitJS makes the page do the thing C17 names (a push whose reply is routed by
// correlationId): call the real door three times, AWAIT each return value, and report
// each outcome back through that same door. A reply that never lands is reported as
// NOT_RESOLVED in the page's own words, which is the only evidence shape ruling P2
// accepts.
//
// Two build notes, both learned the hard way on this box. (1) Every envelope is
// composed by Go, not by JS: in JS two adjacent string literals separated only by a
// line break MERGE, so the first version of this script asked for a requestId whose
// text was literally `ac14-" + i + "` and the ruler then waited for a report nobody
// had sent. (2) The script is one line and opens with a beacon, so "the page never
// answered" and "the page never ran what I sent it" are different readings.
func ac14AwaitJS() string {
	envs := make([]string, 0, 3)
	reports := make([]string, 0, 3)
	for i := 0; i < 3; i++ {
		envs = append(envs, fmt.Sprintf(`'{"method":"panel.mode.request","requestId":"ac14-%d","source":"panel-composer","to":"ask_every_step"}'`, i))
		reports = append(reports, fmt.Sprintf(`'ac14r-%d'`, i))
	}
	env := func(ridExpr, toExpr string) string {
		return "'{\"method\":\"panel.mode.request\",\"requestId\":\"' + " + ridExpr +
			" + '\",\"source\":\"panel-composer\",\"to\":\"' + " + toExpr + " + '\"}'"
	}
	return "(async function(){" +
		"var ENVS=[" + strings.Join(envs, ",") + "];" +
		"var REP=[" + strings.Join(reports, ",") + "];" +
		// The beacon: proves the script ran, before anything can await.
		"window.wispDispatch(" + env("'ac14r-beacon'", "'SCRIPT-RAN'") + ");" +
		"for (var i = 0; i < ENVS.length; i++) {" +
		"var got = 'NOT_RESOLVED';" +
		"try { got = await Promise.race([" +
		"Promise.resolve(window.wispDispatch(ENVS[i])).then(function(v){ return 'REPLIED'; })," +
		"new Promise(function(res){ setTimeout(function(){ res('TIMEOUT-2S'); }, 2000); })" +
		"]); } catch (e) { got = 'THREW'; }" +
		"window.wispDispatch(" + env("REP[i]", "got.replace(/[^A-Za-z0-9-]/g, '')") + ");" +
		"} })();"
}

func TestAC14AwaitedBindingReplyReachesThePage(t *testing.T) {
	rp, mgr := startPanelForTest(t)
	showAndWait(t, rp)
	h := mgr.disp.Mode.(*recordingModeHandler)

	evalOnPanelThread(t, rp, ac14AwaitJS())

	outcomes := make([]string, 0, 3)
	for i := 0; i < 3; i++ {
		outcomes = append(outcomes, awaitReport(t, h, fmt.Sprintf("ac14r-%d", i), "AC#14's reply hop"))
	}
	joined := strings.Join(outcomes, ",")
	calls := countRequests(h, "ac14-")
	t.Logf("AC#14 nail 1 (reply hop), page's own words: %q (Go's handler was reached by %d of the 3 real requests)", joined, calls)
	if strings.Contains(joined, "NOT_RESOLVED") || strings.Contains(joined, "TIMEOUT-2S") ||
		strings.Contains(joined, "THREW") || !strings.Contains(joined, "REPLIED") {
		t.Errorf("the awaited JS binding reply did not reach the page: the page reported %q. Go did receive the calls (%d reached the handler), so this is the H10 half and not H3: webview.Dispatch only appends a closure and posts a thread message, and only the library's Run() drains that queue (33-p1 §A, R25 versus R26). The panel thread has to hand its pump over to Run()", joined, calls)
	}
	if calls != 3 {
		t.Errorf("the page says replies arrived but Go was reached by %d of the 3 real requests - the two halves of this hop have to agree, and one instrument may not stand in for both", calls)
	}
}

// countRequests counts recorded requests whose requestId starts with prefix - how the
// ruler separates the page's three real calls from its three reports.
func countRequests(h *recordingModeHandler, prefix string) int {
	n := 0
	for _, req := range h.all() {
		if strings.HasPrefix(req.RequestID, prefix) {
			n++
		}
	}
	return n
}

// ---------------------------------------------------------------------------
// AC#14, nail 2: the Go -> page Eval push (a DIFFERENT dimension, its own nail)
// ---------------------------------------------------------------------------

// TestAC14GoSideEvalPushReachesThePage is the second nail and stays a separate test:
// 33-p1 §A measured the push dimension arriving (R27, 907ms) while the reply dimension
// did not (R25) in the same process, and ticket 33's own wording forbids blending the
// two into one sentence.
func TestAC14GoSideEvalPushReachesThePage(t *testing.T) {
	rp, mgr := startPanelForTest(t)
	showAndWait(t, rp)
	h := mgr.disp.Mode.(*recordingModeHandler)

	// A value the page can only know if Go's push landed. This nail never awaits a
	// binding, so nothing here can be satisfied by the reply hop, and vice versa.
	const pushed = "PUSHED-33R5-OK"
	evalOnPanelThread(t, rp, fmt.Sprintf("document.title = %q;", pushed))
	evalOnPanelThread(t, rp, "window.wispDispatch("+reportJSEnv("ac14-push", "document.title")+");")

	answer := awaitReport(t, h, "ac14-push", "AC#14's push hop")
	t.Logf("AC#14 nail 2 (Eval push hop), page's own words: title=%q", answer)
	if answer != pushed {
		t.Errorf("Go's Eval push did not reach the document: the page reports its title as %q, want %q. This is the push dimension, separate from the awaited-reply dimension in TestAC14AwaitedBindingReplyReachesThePage - one arriving says nothing about the other", answer, pushed)
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
//  1. a sample that turns out to be the panel itself must not overwrite the honest
//     prior - that overwrite IS the defect 33-v1 §A#27 found;
//  2. the recorded prior survives a Hide, so the next Hide still has a target (see
//     the reasoning on PanelManager.Hide).
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
