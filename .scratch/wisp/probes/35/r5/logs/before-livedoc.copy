//go:build windows && winlive

// Acceptance rig 35-v2 (验收台件 only - NOT production code, NOT a shipping judge,
// one file added by this leg, no existing assertion touched).
//
// WHAT IT MEASURES - ticket 35 AC#6 (`:52`) on a REAL WebView2 window: does a
// page-side call of window.chrome.webview.postMessage(<composer envelope>) reach the
// Go door once shape 甲③'s forwarding hook is installed? Today's yard for that edge,
// cmd/wisp/panel_transport_35r2_test.go, runs in a hand-written JS model that ASSUMES
// assigning over chrome.webview.postMessage takes effect (its jsObject.set can never
// fail). This file checks that assumption against the control itself.
//
// PRODUCT PATH ONLY. The window comes up through PanelManager.bringUp
// (cmd/wisp/panel_host_windows.go:318), the one place that calls
// installPanelTransport (:408 -> :693 -> w.Init(panelPostMessageForwardInit) at :700).
// The steady-state pump is the library's own WebView.Run() - the same call the resident
// process makes at panel_resident_windows.go:299 - on a thread locked for life. This
// file starts no pump of its own and calls no product helper the resident path lacks.
//
// WHY THIS IS NOT THE BANNED N6 SHAPE. Ticket `:64-65` bans
// w.Eval("window.wispDispatch(...)") as evidence: there the test plays the page *at the
// door*, so it stays green under every forwarding shape and has zero discriminating
// power. This rig never touches the door. The script it evaluates calls
// window.chrome.webview.postMessage - exactly the property the whole bug lived on - so
// a forwarding shape that re-enters (the landed sub-shape ①, commit fb2fb802) yields
// ZERO arrivals here rather than a green. That control is run from /d/tmp with
// `go test -overlay` over a copy of the old hook; no repo file is edited for it.
//
// EVERY JUDGEMENT IS MADE GO-SIDE, off the recorder attached to
// panel.ComposerDispatch.Mode. Three envelopes are posted, distinguishable only by
// requestId: pc-35v2-canary (transport works at all), pc-35v2-pagepost-1 (the page's
// own envelope shape from frontend/src/lib/panel.ts:246 + :211 arrives with its
// correlationId, source and mode target intact), and pc-35v2-DESC-... (a diagnostic
// lane: page-side facts about the property that have no other way out of the document).
// The diagnostic lane asserts nothing; it is a reading, and it is labelled as one.

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	webview2 "github.com/jchv/go-webview2"

	"github.com/CarlosShao/wisp/internal/panel"
)

// probeCanary35v2 / probePageShape35v2: one postMessage on the page's own transport.
const probeCanary35v2 = `(function () {
  window.chrome.webview.postMessage('{"method":"panel.mode.request","requestId":"pc-35v2-canary","source":"panel-composer","to":"ask_every_step"}');
})();`

// probePageShape35v2 is the envelope the composer actually posts
// (frontend/src/lib/panel.ts:246 sendRequest("panel.mode.request", { to }) with the
// requestId + source fields :211 adds) - same key order, same shapes.
const probePageShape35v2 = `(function () {
  window.chrome.webview.postMessage('{"method":"panel.mode.request","requestId":"pc-35v2-pagepost-1","source":"panel-composer","to":"ask_every_step"}');
})();`

// probeDescribe35v2 reads the property back in the document and carries what it finds
// out through the requestId of a roster envelope (the only thing that can reach Go).
// Tokens: OWN|INHERIT|NONE = is postMessage an own data property now (the hook's
// assignment took effect); RO|WR = writable; NC|CF = configurable; HOOK|NOHOOK = the
// hook's own idempotency stamp; NOERR|ERR = did an indirect call through the property
// throw. Asserted for nothing - it is a measurement.
const probeDescribe35v2 = `(function () {
  var cw = window.chrome.webview;
  var d = Object.getOwnPropertyDescriptor(cw, "postMessage");
  var tok = [];
  tok.push(d ? "OWN" : (cw.postMessage ? "INHERIT" : "NONE"));
  tok.push(d && d.writable === false ? "RO" : "WR");
  tok.push(d && d.configurable === false ? "NC" : "CF");
  tok.push(cw.__wispForwardInstalled === true ? "HOOK" : "NOHOOK");
  var err = "NOERR";
  try {
    var f = cw.postMessage;
    f.call(cw, '{"method":"panel.mode.request","requestId":"pc-35v2-desc-ping","source":"panel-composer","to":"ask_every_step"}');
  } catch (e) {
    err = "ERR";
  }
  tok.push(err);
  window.chrome.webview.postMessage('{"method":"panel.mode.request","requestId":"pc-35v2-DESC-' + tok.join("-") + '","source":"panel-composer","to":"ask_every_step"}');
})();`

// t35v2Recorder answers the mode route and remembers what the door delivered.
type t35v2Recorder struct {
	mu   sync.Mutex
	seen []panel.ComposerRequest
}

func (r *t35v2Recorder) HandleModeRequest(_ context.Context, req panel.ComposerRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seen = append(r.seen, req)
	return nil
}

func (r *t35v2Recorder) snapshot() []panel.ComposerRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]panel.ComposerRequest(nil), r.seen...)
}

func (r *t35v2Recorder) ids() []string {
	var out []string
	for _, req := range r.snapshot() {
		out = append(out, req.RequestID)
	}
	return out
}

func (r *t35v2Recorder) has(id string) bool {
	return strings.Contains(strings.Join(r.ids(), ","), id)
}

func (r *t35v2Recorder) request(id string) (panel.ComposerRequest, bool) {
	for _, req := range r.snapshot() {
		if req.RequestID == id {
			return req, true
		}
	}
	return panel.ComposerRequest{}, false
}

// post35v2 routes one page-side evaluation onto the window's own thread and waits a
// bounded number of seconds for the named envelope to come back through the door.
// Re-posting is bounded and reported, so "took N attempts" is part of the reading.
func post35v2(t *testing.T, wv webview2.WebView, r *t35v2Recorder, id, script string, attempts int) (bool, int) {
	t.Helper()
	for i := 1; i <= attempts; i++ {
		wv.Dispatch(func() { wv.Eval(script) })
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			if r.has(id) {
				return true, i
			}
			time.Sleep(50 * time.Millisecond)
		}
		t.Logf("35v2: %s had not arrived after attempt %d (%d door recordings so far)", id, i, len(r.ids()))
	}
	return r.has(id), attempts
}

func TestLive35v2PagePostMessageEnvelopeReachesTheGoDoor(t *testing.T) {
	// W4 by-product: dump the string Go really renders, so the Sprintf substitution is
	// inspected rather than trusted.
	rendered := panelPostMessageForwardInit
	t.Logf("35v2 RENDERED FORWARDING SCRIPT (Go side, %%[1]s already substituted):\n%s\n35v2 ---- end of rendered script ----", rendered)
	if strings.Contains(rendered, "%") {
		t.Errorf("35v2 the rendered script still carries a %% verb (unsubstituted template left in the string)")
	}
	for _, want := range []string{"window.wispDispatch", "var native = cw.postMessage", "native.call(cw, message)", "finally { inside = false; }"} {
		if !strings.Contains(rendered, want) {
			t.Errorf("35v2 rendered script is missing the shape the comment at panel_host_windows.go:638 promises: %q", want)
		}
	}

	dataPath := filepath.Join(os.TempDir(), "wisp-35v2-panel-profile")
	rec := &t35v2Recorder{}
	disp := &panel.ComposerDispatch{Mode: rec}
	mgr := NewPanelManager(disp, nil, dataPath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	self := uint32(os.Getpid())
	before := readWebviewTree(t, self)
	t.Logf("35v2 webview processes at start: tree=%d machineNamed=%d pids=%d", before.TreeWebview, before.MachineNamed, len(before.TreePIDs))

	ready := make(chan webview2.WebView, 1)
	pumpDone := make(chan struct{}, 1)
	go func() {
		// owner: 35-v2 acceptance live thread (bringUp + the library pump); recover below.
		defer func() {
			if p := recover(); p != nil {
				t.Errorf("35v2 live thread panicked: %v", p)
				ready <- nil
				pumpDone <- struct{}{}
			}
		}()
		runtime.LockOSThread()
		if err := mgr.bringUp(ctx); err != nil {
			t.Errorf("35v2 bringUp on the live thread: %v", err)
			ready <- nil
			return
		}
		wv := mgr.currentWindow()
		ready <- wv
		if wv == nil {
			return
		}
		wv.Run() // the library pump, exactly like panel_resident_windows.go:299
		pumpDone <- struct{}{}
	}()

	var wv webview2.WebView
	select {
	case got := <-ready:
		wv = got
	case <-time.After(60 * time.Second):
		t.Fatalf("35v2 NO LIVE READING: bring-up never handed the window over within 60 s")
	}
	if wv == nil {
		t.Fatalf("35v2 NO LIVE READING: the product bring-up produced no window on this box")
	}
	if !mgr.IsCreated() {
		t.Fatalf("35v2 host reports not created after bring-up")
	}
	t.Logf("35v2 cold bring-up measured on this box: %.3f ms, hwnd=%#x", mgr.LastColdMs(), mgr.windowHandle())

	okCanary, nCanary := post35v2(t, wv, rec, "pc-35v2-canary", probeCanary35v2, 4)
	t.Logf("35v2 canary arrival=%v after %d post(s); door has recorded %v", okCanary, nCanary, rec.ids())

	okPage, nPage := post35v2(t, wv, rec, "pc-35v2-pagepost-1", probePageShape35v2, 3)
	if !okPage {
		t.Errorf("35v2 AC#6 LIVE ARRIVAL NOT MEASURED: chrome.webview.postMessage from the page did not reach the Go door (%d recordings: %v)", len(rec.ids()), rec.ids())
	} else {
		req, _ := rec.request("pc-35v2-pagepost-1")
		t.Logf("35v2 DELIVERED TO GO: method=%q requestId=%q source=%q to=%q (attempt %d)", req.Method, req.RequestID, req.Source, req.To, nPage)
		if req.Method != "panel.mode.request" || req.Source != "panel-composer" || req.To != "ask_every_step" {
			t.Errorf("35v2 the envelope that arrived is not the page's envelope field for field: %+v", req)
		}
	}

	okDesc, _ := post35v2(t, wv, rec, "pc-35v2-DESC", probeDescribe35v2, 3)
	for _, id := range rec.ids() {
		if strings.HasPrefix(id, "pc-35v2-DESC") {
			t.Logf("35v2 PAGE-SIDE FACTS (diagnostic lane, asserted for nothing): %s", id)
		}
	}
	if !okDesc {
		t.Logf("35v2 diagnostic lane did not come back - the property facts are UNOBTAINED, the arrival reading above stands on its own")
	}

	// Every judgement above is Go-side; this line only prints what the door saw.
	t.Logf("35v2 all door recordings: %v", rec.ids())

	// Teardown: Destroy runs on the thread that owns the window; Run() returns on the
	// WM_QUIT the library posts from WM_DESTROY.
	wv.Dispatch(func() { mgr.Destroy() })
	select {
	case <-pumpDone:
	case <-time.After(25 * time.Second):
		t.Errorf("35v2 the library pump did not return after Destroy - the window may still be up")
	}
	time.Sleep(3 * time.Second)
	after := readWebviewTree(t, self)
	t.Logf("35v2 webview processes after teardown: tree=%d machineNamed=%d pids=%d (start tree=%d pids=%d)",
		after.TreeWebview, after.MachineNamed, len(after.TreePIDs), before.TreeWebview, len(before.TreePIDs))
	if len(after.TreePIDs) > len(before.TreePIDs) {
		t.Errorf("35v2 left %d tree processes behind against %d at start", len(after.TreePIDs), len(before.TreePIDs))
	}
	fmt.Fprint(os.Stderr, "")
}
