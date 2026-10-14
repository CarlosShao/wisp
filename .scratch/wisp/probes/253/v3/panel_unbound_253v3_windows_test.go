//go:build windows

package main

// 253-v3 acceptance overlay test - NOT part of cmd/wisp's committed tree.
// Master copy: .scratch/wisp/probes/253/v3/panel_unbound_253v3_windows_test.go
// Overlayed into cmd/wisp/ for two runs (post-fix tree / pre-fix overlay) and
// removed again; every reading in the acceptance note comes from THIS file's
// four shapes, none of which 253-r6's own test file planted:
//
//	(A) an unbound frame with an EXPLICIT id (977) and a method r6 never used
//	    ("legacy.echo") plus a params array - does the refusal carry the frame's
//	    own numbers, and does the slot with that id (not 0) receive it?
//	(B) the SLOTLESS shape: the old envelope fired with NO _rpc[0] slot at all.
//	    The page arms a getter trap on that property, so "the host's reply script
//	    tried to address slot 0" and "the page holds nothing" are both readable.
//	    This is the shape 253-r6 declared structurally unreceivable on the page
//	    side; the reading below says what actually happens.
//	(C) a message that is NOT JSON at all (msgcb's first silent death) - r6 did
//	    not plant this one in a real window either.
//	(D) positive control: a DIRECT window.wispDispatch call (the bound door, no
//	    fold involved) settles with the receipt, reaches the handler once, and
//	    draws no new unbound-door warning.
//	(E) P3: window.wispDispatch momentarily deleted, so the forwarding hook's
//	    typeof guard routes the page's own postMessage to the native exit - the
//	    raw envelope reaches msgcb and must take the same refusal.
//
// One window, four readings, on purpose (the known two-windows-per-binary
// flake; see 253-r6's own note). All four report through the product's own
// door, so every verdict is the page's own words (the 33-r5/AC#14 pattern).

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
)

// id977ShapeJS253v3: register _rpc[977], fire the explicit-id unbound frame.
func id977ShapeJS253v3(reportID string) string {
	frame := `{"id":977,"method":"legacy.echo","params":["x"],"requestId":"253v3-id7-env"}`
	return fmt.Sprintf(`(function () {
  window._rpc = window._rpc || { nextSeq: 1 };
  var slot = new Promise(function (settle) {
    window._rpc[977] = {
      resolve: function (v) { settle('RESOLVED(' + JSON.stringify(v === undefined ? null : v) + ')'); },
      reject: function (e) { settle('REJECTED(' + JSON.stringify(String(e)) + ')'); }
    };
  });
  slot.then(function (s) {
    window.wispDispatch(%s);
  });
  window.external.invoke(%s);
})();`, reportJSEnv(reportID, `s.replace(/["\\]/g, "'")`), strconv.Quote(frame))
}

// slotlessShapeJS253v3: no slot registered; the getter trap marks the host's
// reply script addressing _rpc[0] and reports the marker through the door.
func slotlessShapeJS253v3(reportID string) string {
	return fmt.Sprintf(`(function () {
  window._rpc = window._rpc || { nextSeq: 1 };
  var seen = { hit: false };
  Object.defineProperty(window._rpc, '0', {
    configurable: true,
    get: function () {
      seen.hit = true;
      Promise.resolve().then(function () {
        window.wispDispatch(%s);
      });
      return undefined;
    }
  });
  window.external.invoke(%s);
})();`, reportJSEnv(reportID, `'getterHit=' + seen.hit`), strconv.Quote(pageModeRequestEnvelope))
}

// nonJSONShapeJS253v3: a plain non-JSON string through the same door, then a
// beacon report (the report arrives after the host has handled the string, so
// the slog record is already written when the report reaches the handler).
func nonJSONShapeJS253v3(reportID string) string {
	return fmt.Sprintf(`(function () {
  window.external.invoke('253v3 plain text, not an RPC frame');
  Promise.resolve().then(function () {
    window.wispDispatch(%s);
  });
})();`, reportJSEnv(reportID, `'sent-non-json'`))
}

// unfoldedShapeJS253v3: P3 - the world where the forwarding hook's typeof guard
// routes the page's own postMessage to the native exit (the door is momentarily
// NOT callable). The envelope must therefore arrive raw at msgcb, exactly like
// the old door's shape. The getter trap from (B) is deleted first, because it is
// a configurable accessor and would swallow the slot assignment below.
func unfoldedShapeJS253v3(reportID string) string {
	return fmt.Sprintf(`(function () {
  window._rpc = window._rpc || { nextSeq: 1 };
  try { delete window._rpc['0']; } catch (e) {}
  var slot = new Promise(function (settle) {
    window._rpc[0] = {
      resolve: function (v) { settle('RESOLVED(' + JSON.stringify(v === undefined ? null : v) + ')'); },
      reject: function (e) { settle('REJECTED(' + JSON.stringify(String(e)) + ')'); }
    };
  });
  slot.then(function (s) {
    window.wispDispatch(%s);
  });
  var saved = window.wispDispatch;
  delete window.wispDispatch;
  window.chrome.webview.postMessage(%s);
  window.wispDispatch = saved;
})();`, reportJSEnv(reportID, `s.replace(/["\\]/g, "'")`), strconv.Quote(pageModeRequestEnvelope))
}

// controlShapeJS253v3: the bound door called directly, with its own envelope
// requestId so the handler count is attributable.
func controlShapeJS253v3(reportID, envRequestID string) string {
	env := fmt.Sprintf(`{"method":"panel.mode.request","requestId":"%s","source":"panel-composer","to":"ask_every_step"}`, envRequestID)
	return fmt.Sprintf(`(function () {
  var call = window.wispDispatch(%s);
  Promise.resolve(call).then(function (v) {
    window.wispDispatch(%s);
  }, function (e) {
    window.wispDispatch(%s);
  });
})();`, strconv.Quote(env),
		reportJSEnv(reportID, `('RESOLVED(' + JSON.stringify(v === undefined ? null : v) + ')').replace(/["\\]/g, "'")`),
		reportJSEnv(reportID, `('REJECTED(' + JSON.stringify(String(e)) + ')').replace(/["\\]/g, "'")`))
}

func Test253v3NewShapesOnTheRealHost(t *testing.T) {
	captured := captureSlog253r6(t)
	h := &recordingModeHandler{}
	sink := &auditSink253r6{}
	disp := &panel.ComposerDispatch{Mode: h, Audit: sink.auditf}
	rp, _ := startPanelFor253r6(t, disp)
	showAndWait(t, rp)

	// --- (A) explicit id 977, method r6 never used -----------------------------
	evalOnPanelThread(t, rp, id977ShapeJS253v3("ac3v3-id7"))
	outA := awaitReport(t, h, "ac3v3-id7", "the id=977 legacy.echo frame")
	t.Logf("253v3 (A) id=977 legacy.echo, the page's own words: %q", outA)
	if !strings.HasPrefix(outA, "REJECTED(") {
		t.Errorf("253v3 RED (A): the page settled with %q, want REJECTED(<sentence>)", outA)
	}
	if !strings.Contains(outA, "legacy.echo") {
		t.Errorf("253v3 RED (A): the refusal does not name the method the frame carried: %q", outA)
	}
	if got := countRequests(h, "253v3-id7-env"); got != 0 {
		t.Errorf("253v3 RED (A): the router saw the unbound frame %d time(s)", got)
	}
	lineA := ""
	for _, line := range captured.all() {
		if strings.Contains(line, "names no bound door") && strings.Contains(line, "method=legacy.echo") {
			lineA = line
		}
	}
	if lineA == "" {
		t.Errorf("253v3 RED (A): no slog record names legacy.echo; records: %v", captured.all())
	} else {
		t.Logf("253v3 (A) Go-side record: %s", lineA)
		if !strings.Contains(lineA, "id=977") || !strings.Contains(lineA, "params=1") {
			t.Errorf("253v3 RED (A): the record must carry the frame's own id/params: %s", lineA)
		}
	}

	// --- (B) the slotless shape ------------------------------------------------
	evalOnPanelThread(t, rp, slotlessShapeJS253v3("ac3v3-slotless"))
	outB := awaitReport(t, h, "ac3v3-slotless", "the slotless shape's getter trap")
	t.Logf("253v3 (B) slotless, the page's own words: %q", outB)
	if outB != "getterHit=true" {
		t.Errorf("253v3 RED (B): the trap read %q - want getterHit=true (the host's reply script addressed slot 0) with the page holding nothing", outB)
	}
	if got := countRequests(h, pageModeRequestID); got != 0 {
		t.Errorf("253v3 RED (B): the router saw the slotless envelope %d time(s)", got)
	}
	lineB := ""
	for _, line := range captured.all() {
		if strings.Contains(line, "names no bound door") && strings.Contains(line, "method=panel.mode.request") {
			lineB = line
		}
	}
	if lineB == "" {
		t.Errorf("253v3 RED (B): no slog record for the slotless old envelope; records: %v", captured.all())
	} else {
		t.Logf("253v3 (B) Go-side record: %s", lineB)
	}

	// --- (C) the non-JSON message ----------------------------------------------
	evalOnPanelThread(t, rp, nonJSONShapeJS253v3("ac3v3-nonjson"))
	outC := awaitReport(t, h, "ac3v3-nonjson", "the non-JSON shape's beacon")
	t.Logf("253v3 (C) non-JSON, the page's own words: %q", outC)
	lineC := ""
	for _, line := range captured.all() {
		if strings.Contains(line, "not an RPC frame") {
			lineC = line
		}
	}
	if lineC == "" {
		t.Errorf("253v3 RED (C): no slog record for the non-JSON message; records: %v", captured.all())
	} else {
		t.Logf("253v3 (C) Go-side record: %s", lineC)
	}

	// --- (D) positive control: the bound door directly -------------------------
	warnBefore := 0
	for _, line := range captured.all() {
		if strings.Contains(line, "names no bound door") {
			warnBefore++
		}
	}
	evalOnPanelThread(t, rp, controlShapeJS253v3("ac3v3-control", "253v3-control-env"))
	outD := awaitReport(t, h, "ac3v3-control", "the bound door's receipt (positive control)")
	t.Logf("253v3 (D) positive control (direct window.wispDispatch): %q", outD)
	if !strings.Contains(outD, "RESOLVED") {
		t.Errorf("253v3 RED (D): the bound door settled with %q, want a receipt", outD)
	}
	if got := countRequests(h, "253v3-control-env"); got != 1 {
		t.Errorf("253v3 RED (D): the control envelope reached the handler %d time(s), want exactly 1", got)
	}
	warnAfter := 0
	for _, line := range captured.all() {
		if strings.Contains(line, "names no bound door") {
			warnAfter++
		}
	}
	if warnAfter != warnBefore {
		t.Errorf("253v3 RED (D): the unbound-door warning fired during the positive control (%d -> %d)", warnBefore, warnAfter)
	}

	// --- (E) P3: the door momentarily absent, postMessage goes native ----------
	warnBeforeE := warnAfter
	evalOnPanelThread(t, rp, unfoldedShapeJS253v3("ac3v3-unfolded"))
	outE := awaitReport(t, h, "ac3v3-unfolded", "the door-absent raw postMessage shape")
	t.Logf("253v3 (E) door absent, the page's own words: %q", outE)
	if !strings.HasPrefix(outE, "REJECTED(") {
		t.Errorf("253v3 RED (E): the page settled with %q, want REJECTED(<sentence>)", outE)
	}
	if !strings.Contains(outE, panel.MethodModeRequest) {
		t.Errorf("253v3 RED (E): the refusal does not name the method the envelope carried: %q", outE)
	}
	if got := countRequests(h, pageModeRequestID); got != 0 {
		t.Errorf("253v3 RED (E): the router saw the door-absent envelope %d time(s)", got)
	}
	warnAfterE := 0
	for _, line := range captured.all() {
		if strings.Contains(line, "names no bound door") {
			warnAfterE++
		}
	}
	if warnAfterE != warnBeforeE+1 {
		t.Errorf("253v3 RED (E): want exactly one new unbound-door record for the door-absent shape (%d -> %d)", warnBeforeE, warnAfterE)
	}
}
