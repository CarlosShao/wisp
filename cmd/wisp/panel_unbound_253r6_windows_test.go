//go:build windows

package main

// Ticket 253 AC#3 (`:18`) - the silent-drop branch must be loud, judged on the real host.
//
// THE HOLE THIS FILE MEASURES. 253's own words: 页面若还按旧形状直接 postMessage，宿主库会把
// 这次调用的 Promise 解析成 null 并返回、Go 侧零审计. Since ticket 299-r1 the RPC pipe is
// THIS package's own code (cmd/wisp/panel_control_windows.go's msgcb/callbinding, the
// library's webview.go:139-220 ported with the receiver renamed), so the branch that
// answers (nil, nil) for a method no door bound is in-repo reachable for the first time:
// the message reaches the host's message callback, `callbinding` finds nothing, and the
// reply script resolves a window._rpc[id] slot - with zero Go-side record. The other
// drop branch in the same function (a message that is not JSON at all) leaves a bare
// library-shaped stderr line.
//
// WHAT THE FAKE LEG IS. One old-shape envelope - the page's own composer envelope
// verbatim, fed straight into window.external.invoke, which is the pre-ticket-35 door
// and the one path the postMessage forwarding hook (panelPostMessageForwardInit) routes
// to the native exit by design: `inside` is true inside invoke, so the message reaches
// msgcb exactly as the census said it would - an RPC frame with id=0 and an unbound
// method. The page registers the slot the library's reply script addresses (_rpc[0],
// webview.go:131-135 gives a no-id envelope ID=0), because that slot is the only way
// "what the page got" is observable at all; without it the reply script throws a
// TypeError nobody sees, which IS the ticket's 用户视角 symptom.
//
// ★WHY THE PAGE'S OWN WORDS. Every verdict is built from the string the PAGE sends back
// through the product's door (the 33-r5/AC#14 pattern), never from "Dispatch was
// called" or "Eval returned" - see panel_resident_windows_test.go's P2 note. Before the
// fix this same rig reads `RESOLVED(null)`: the silent null, verbatim. The Go-side half
// is a slog record, because that is the channel the resident process keeps (logsink.go
// installs the JSONL sink over slog; a plain log.Printf line is stderr-only).
//
// ONE WINDOW, THREE READINGS, on purpose: on this box a second panel window created
// right after a previous one intermittently dies mid-create (measured; it reproduces on
// the pristine tree with two pre-existing tests back to back - logs/13), so this file
// keeps one windowed case and takes all three readings in it: the fake leg, the
// positive control on the same assembly, and the folded postMessage path.
//
// WHAT THIS FILE IS NOT. It does not model the pipe (panel_transport_35r2_test.go does,
// and stays its own ticket's ruler); it drives the shipping host through
// startResidentPanel, the same thread and pump the product uses. No new pump, no new
// binding - the report rides the door that already exists so no document is re-created
// behind the ruler's back (33-r5's own reason for that choice).

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
)

// Two page strings this file reuses from their own tickets rather than retyping:
// pageModeRequestEnvelope is 35-r1's verbatim page envelope (the "old shape" here - the
// message as the page writes it, with no {id,method,params} RPC wrapper), and
// envelope35r3UnregisteredName is 35-r3's verbatim panel.approval.request envelope (a
// name the roster refuses). Reusing them couples the fixtures loudly (a rename is a
// compile error) instead of silently drifting two copies apart.

// startPanelFor253r6 is startPanelForTest with the dispatch handed in, because AC#3's
// Go-side readings are the dispatch's own audit sink and the host's slog record; the
// 33-r5 helper builds its own dispatch and cannot carry either.
func startPanelFor253r6(t *testing.T, disp *panel.ComposerDispatch) (*residentPanel, *PanelManager) {
	t.Helper()
	dataPath := filepath.Join(os.TempDir(), "wisp-253r6-panel-profile")
	mgr := NewPanelManager(disp, builtinAssetsOrTestNil(t), dataPath)
	rp := startResidentPanel(observe.NewRegistry(), mgr)
	t.Cleanup(rp.stop)
	return rp, mgr
}

// auditSink253r6 records the dispatch chain's own audit lines, so the folded path's
// criterion ① (a readable Go-side record) is read off the shipping sink rather than
// inferred from a message count.
type auditSink253r6 struct {
	mu    sync.Mutex
	lines []string
}

func (a *auditSink253r6) auditf(format string, args ...any) {
	a.mu.Lock()
	a.lines = append(a.lines, fmt.Sprintf(format, args...))
	a.mu.Unlock()
}

func (a *auditSink253r6) all() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]string(nil), a.lines...)
}

// slogCapture253r6 is a whole-process slog handler installed for one test's lifetime
// (no test in this package runs in parallel), so the unbound-frame warning can be
// asserted verbatim instead of being read off a console nobody keeps.
type slogCapture253r6 struct {
	mu    sync.Mutex
	lines []string
}

func (c *slogCapture253r6) Enabled(context.Context, slog.Level) bool { return true }

func (c *slogCapture253r6) Handle(_ context.Context, r slog.Record) error {
	var b strings.Builder
	b.WriteString(r.Message)
	r.Attrs(func(a slog.Attr) bool {
		b.WriteString(" ")
		b.WriteString(a.Key)
		b.WriteString("=")
		b.WriteString(a.Value.String())
		return true
	})
	c.mu.Lock()
	c.lines = append(c.lines, b.String())
	c.mu.Unlock()
	return nil
}

func (c *slogCapture253r6) WithAttrs([]slog.Attr) slog.Handler { return c }
func (c *slogCapture253r6) WithGroup(string) slog.Handler      { return c }

func (c *slogCapture253r6) all() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]string(nil), c.lines...)
}

func captureSlog253r6(t *testing.T) *slogCapture253r6 {
	t.Helper()
	prev := slog.Default()
	c := &slogCapture253r6{}
	slog.SetDefault(slog.New(c))
	t.Cleanup(func() { slog.SetDefault(prev) })
	return c
}

// unboundFrameJS253r6 is the fake leg, as page-side JS:
//
//	(1) register window._rpc[0] - the slot the host's reply script addresses for an
//	    envelope with no id - and turn its first settle into a string the page can post;
//	(2) fire the old-shape envelope into window.external.invoke (the pre-35 door);
//	(3) post whatever came back through the product's own door, under its own requestId.
//
// The report goes LAST, inside the promise callback, so "silence" (nothing settles) is
// an awaitReport timeout - a failed measurement, not a pass.
func unboundFrameJS253r6(reportID, envelope string) string {
	return fmt.Sprintf(`(function () {
  var slot = new Promise(function (settle) {
    window._rpc = window._rpc || { nextSeq: 1 };
    window._rpc[0] = {
      resolve: function (v) { settle('RESOLVED(' + JSON.stringify(v === undefined ? null : v) + ')'); },
      reject: function (e) { settle('REJECTED(' + JSON.stringify(String(e)) + ')'); }
    };
  });
  slot.then(function (s) {
    window.wispDispatch(%s);
  });
  window.external.invoke(%s);
})();`, reportJSEnv(reportID, `s.replace(/["\\]/g, "'")`), strconv.Quote(envelope))
}

// foldedEnvelopeJS253r6 posts one envelope the way the page's own transport does
// (window.chrome.webview.postMessage - the call frontend/src/lib/panel.ts makes) and
// reports what the returned promise settled with. It asserts nothing about shape: the
// forwarding hook is what decides whether this call reaches a door.
func foldedEnvelopeJS253r6(reportID, envelope string) string {
	return fmt.Sprintf(`(function () {
  var call = window.chrome.webview.postMessage(%s);
  Promise.resolve(call).then(function (v) {
    window.wispDispatch(%s);
  }, function (e) {
    window.wispDispatch(%s);
  });
})();`, strconv.Quote(envelope),
		reportJSEnv(reportID, `('RESOLVED(' + JSON.stringify(v === undefined ? null : v) + ')').replace(/["\\]/g, "'")`),
		reportJSEnv(reportID, `('REJECTED(' + JSON.stringify(String(e)) + ')').replace(/["\\]/g, "'")`))
}

// Test253r6TheSilentDropBranchIsLoud is AC#3's whole judgement on the shipping host:
//
//	(A) the fake leg - one old-shape envelope fed straight into the RPC callback with no
//	    bound door to answer it. Three readings: the warning the host writes (criterion
//	    ①), what the page's promise settles with (criterion ②, the silent null's
//	    replacement), and that the router was never reached (this is the pipe's death,
//	    not a composer refusal).
//	(B) the positive control on the same assembly: the normal door still reaches the
//	    handler, settles with the accepted (empty) receipt, and draws no warning.
//	(C) the folded postMessage path the ticket's own words describe: an unregistered
//	    name posted on window.chrome.webview.postMessage is carried by the forwarding
//	    hook to the router, refused by the roster guard, and audited - so the fix above
//	    changed nothing about the path ticket 35/303 closed.
func Test253r6TheSilentDropBranchIsLoud(t *testing.T) {
	captured := captureSlog253r6(t)
	h := &recordingModeHandler{}
	sink := &auditSink253r6{}
	disp := &panel.ComposerDispatch{Mode: h, Audit: sink.auditf}
	rp, _ := startPanelFor253r6(t, disp)
	showAndWait(t, rp)

	// --- (A) the fake leg -----------------------------------------------------
	evalOnPanelThread(t, rp, unboundFrameJS253r6("ac3r6-unbound", pageModeRequestEnvelope))

	outcome := awaitReport(t, h, "ac3r6-unbound", "what the page got for the old-shape envelope")
	t.Logf("AC#3 fake leg (old envelope into the RPC callback), the page's own words: %q", outcome)
	if !strings.HasPrefix(outcome, "REJECTED(") {
		t.Errorf("AC#3 RED (criterion ②): the page's promise settled with %q, want a REJECTED(<sentence>) - a silent resolve(null) is the ticket's own symptom (253:12), and it is what this rig reads with the audit branch absent", outcome)
	}
	if !strings.Contains(outcome, panel.MethodModeRequest) {
		t.Errorf("AC#3 RED: the page's rejection does not name the method it sent (%s): %q - a refusal the page cannot match to its own call is not a readable receipt", panel.MethodModeRequest, outcome)
	}
	if got := countRequests(h, pageModeRequestID); got != 0 {
		t.Errorf("AC#3 RED: the router was reached %d time(s) by a frame that named no bound door - this reading must be the PIPE's death, not a composer refusal", got)
	}

	warned := false
	for _, line := range captured.all() {
		if strings.Contains(line, "inbound RPC frame names no bound door") {
			warned = true
			t.Logf("AC#3 criterion ①, the host's own record: %s", line)
			if !strings.Contains(line, panel.MethodModeRequest) {
				t.Errorf("AC#3 RED: the Go-side warning does not name the method (%s): %s", panel.MethodModeRequest, line)
			}
		}
	}
	if !warned {
		t.Errorf("AC#3 RED (criterion ①): nothing in the host's slog records says an inbound RPC frame named no bound door. Records seen: %v", captured.all())
	}

	// --- (B) positive control: the normal door, same assembly ------------------
	before := len(captured.all())
	evalOnPanelThread(t, rp, foldedEnvelopeJS253r6("ac3r6-control", pageModeRequestEnvelope))
	control := awaitReport(t, h, "ac3r6-control", "the normal door's receipt (positive control)")
	t.Logf("AC#3 positive control (the same envelope through window.chrome.webview.postMessage): %q", control)
	if !strings.Contains(control, "RESOLVED") {
		t.Errorf("AC#3 RED (positive control): a normal composer call through the page's own transport settled with %q, want a receipt - the audit branch must not have changed the bound-door path", control)
	}
	if got := countRequests(h, pageModeRequestID); got != 1 {
		t.Errorf("AC#3 RED (positive control): the normal call must be the only one that reached the handler; handler saw %d with requestId %s", got, pageModeRequestID)
	}
	for _, line := range captured.all()[before:] {
		if strings.Contains(line, "inbound RPC frame names no bound door") {
			t.Errorf("AC#3 RED (positive control): the unbound-method warning fired for a call the door answered: %s", line)
		}
	}

	// --- (C) the folded path the ticket names ---------------------------------
	evalOnPanelThread(t, rp, foldedEnvelopeJS253r6("ac3r6-fold-unregistered", envelope35r3UnregisteredName))
	refusal := awaitReport(t, h, "ac3r6-fold-unregistered", "the folded path's refusal sentence")
	t.Logf("AC#3 folded postMessage path (unregistered name), the page's own words: %q", refusal)
	if !strings.Contains(refusal, guard35r3Roster) {
		t.Errorf("AC#3 RED: the folded path settled with %q, want the roster guard's own sentence (%s) - the envelope must reach the router, not die in the pipe", refusal, guard35r3Roster)
	}
	audited := false
	for _, line := range sink.all() {
		if strings.Contains(line, "INBOUND-DISPATCH") && strings.Contains(line, "panel.approval.request") {
			audited = true
			t.Logf("AC#3 criterion ① on the folded path, the dispatch chain's audit line: %s", line)
		}
	}
	if !audited {
		t.Errorf("AC#3 RED: no INBOUND-DISPATCH audit line names panel.approval.request. Lines seen: %v", sink.all())
	}
}
