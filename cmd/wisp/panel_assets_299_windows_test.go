//go:build windows

package main

// Ticket 299 AC#2 (ⓐ has teeth, ⓑ counter-controlled): after a real cold start
// the page in the panel window must be able to fetch its OWN bundle - the very
// defect this ticket exists for was that ./assets/*.js and ./assets/*.css could
// not be fetched at all, so the window rendered white.
//
// THE JUDGEMENT IS THE PAGE'S OWN WORDS (AC#2 ⓐ, the shape ticket 33 AC#14's
// ruling P2 made law: "Dispatch was called" / "Eval returned no error" / "the
// Go-side channel closed" are all true with no reply and are not assertions).
// This file Evals a probe into whatever document is current; the probe READS the
// document (location.href, document.styleSheets and their parsed rule counts,
// #root's own text, the Resource Timing entries) and posts its findings through
// the product's own inbound door (window.wispDispatch - the door
// installPanelTransport binds). Everything asserted below is parsed out of that
// report; nothing asserts that a Go-side call happened.
//
// WHY THE PROBE IS RE-EVALUATED, BOUNDED. The entry handover is a navigation
// (ticket 299 AC#1), so a probe Eval'd before that navigation commits reads the
// previous document (the about:blank round-trip probe page) and honestly reports
// that state. The loop below re-Evals until the page reports the state the
// ticket's completion criterion names (href on the synthetic origin, >=1 parsed
// stylesheet, #root carrying text), or until the monotonic budget ends - at
// which point the LAST report is what gets asserted, verbatim, so a failure
// carries the page's own sentence instead of a timeout.
//
// THE COUNTER-CONTROL (AC#2 ⓑ) is a real run, not a promise: with the filter
// registration removed (go test -overlay over a copy of panel_control_windows.go
// with the AddWebResourceRequestedFilter call dropped - the repo untouched), this
// same case must fail with the page reporting it never got its js/css, or with
// no report arriving at all because the navigation cannot be answered. The
// transcript is in .scratch/wisp/probes/299/r1/30-teeth.md; "no red" is not
// treated as "controlled for".
//
// A BUNDLE IS THE SUBJECT. In a checkout whose embed carries only the tracked
// anchor this case skips by name (no built page, nothing to fetch), the same
// shape panel_pageover_33r10 uses.

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/panel"
)

const (
	// ac299ProbeID is the requestId the page's report rides. The door's parser
	// only requires a non-empty id; the recorder selects by this prefix.
	ac299ProbeID = "ac299-r1-assets"
	// ac299Budget bounds the re-evaluate loop. It is a monotonic deadline, not
	// a wall-clock timeout (ban #4), and 20s is ~10x the measured local load.
	ac299Budget = 20 * time.Second
)

// ac299ProbeJS is the page's probe: pure reads of the document's own state, one
// report out through the product's inbound door. The sanitise step exists only
// so a hostile location.href cannot break the JSON envelope; every character the
// facts above use survives it.
func ac299ProbeJS() string {
	return "(function(){" +
		"function snap(){" +
		"var t='';" +
		"t+='href='+location.href+';';" +
		"t+='base='+document.baseURI+';';" +
		"t+='ready='+document.readyState+';';" +
		"var sheetList=document.styleSheets, n=sheetList.length, rules=0;" +
		"for (var i=0;i<n;i++){ try { rules += (sheetList[i].cssRules ? sheetList[i].cssRules.length : 0); } catch (e) { rules += 0; } }" +
		"t+='sheets='+n+';rules='+rules+';';" +
		"var root=document.getElementById('root');" +
		"t+='root='+(root?'yes':'no')+';';" +
		"t+='roottext='+(root?(''+(root.textContent||'')).replace(/\\s+/g,' ').trim().length:-1)+';';" +
		"var res=performance.getEntriesByType('resource'), njs=0, ncss=0;" +
		"for (var j=0;j<res.length;j++){ if (/\\/assets\\/.*\\.js/.test(res[j].name)) njs++; if (/\\/assets\\/.*\\.css/.test(res[j].name)) ncss++; }" +
		"t+='resjs='+njs+';rescss='+ncss+';';" +
		"var s=(''+t).replace(/[^A-Za-z0-9 =;.:,\\/_-]/g,'_');" +
		"window.wispDispatch(" + reportJSEnv(ac299ProbeID, "s") + ");" +
		"}" +
		"if (document.readyState==='complete'){ setTimeout(snap, 300); }" +
		"else { window.addEventListener('load', function(){ setTimeout(snap, 300); }, {once:true}); setTimeout(snap, 4000); }" +
		"})();"
}

// ac299LatestReport returns the LAST report the door recorded under the probe id,
// so re-evaluations overwrite rather than race each other.
func ac299LatestReport(h *recordingModeHandler) (string, bool) {
	var out string
	found := false
	for _, req := range h.all() {
		if req.RequestID == ac299ProbeID {
			out = req.To
			found = true
		}
	}
	return out, found
}

// ac299Token reads one key=value token out of the page's report string.
func ac299Token(report, key string) (string, bool) {
	i := strings.Index(report, key+"=")
	if i < 0 {
		return "", false
	}
	rest := report[i+len(key)+1:]
	if j := strings.IndexByte(rest, ';'); j >= 0 {
		rest = rest[:j]
	}
	return rest, true
}

// ac299Int reads one integer token; a missing or unparseable token answers
// (0, false), which the caller reports as its own finding.
func ac299Int(report, key string) (int, bool) {
	v, ok := ac299Token(report, key)
	if !ok || v == "" {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// ac299Satisfied is the ticket's completion criterion, read off the page's own
// words: the document is the panel page on the synthetic origin, at least one
// stylesheet parsed with at least one rule (the .css arrived AND was accepted),
// and #root carries text (the .js arrived AND ran).
func ac299Satisfied(report string) bool {
	href, _ := ac299Token(report, "href")
	if href != panelEntryURL {
		return false
	}
	sheets, okS := ac299Int(report, "sheets")
	rules, okR := ac299Int(report, "rules")
	root, _ := ac299Token(report, "root")
	roottext, okT := ac299Int(report, "roottext")
	return okS && okR && okT && sheets >= 1 && rules > 0 && root == "yes" && roottext > 0
}

// TestPanelPageFetchesItsOwnAssets299 is AC#2 ⓐ on a real window: the panel page
// must fetch and execute its own js/css. Every assertion below reads tokens out
// of the page's own report (see the file header); the Go side contributes only
// the probe transport and the bounded re-evaluation loop.
func TestPanelPageFetchesItsOwnAssets299(t *testing.T) {
	assets := builtinAssetsOrTestNil(t)
	if assets == nil || !assets.Built() {
		t.Skipf("299 AC#2 has no subject in this tree: the embed carries no built page bundle (frontend/dist is the tracked anchor only), so the page's ./assets/* references do not exist to fetch")
	}
	if _, _, err := assets.Resolve(panel.EntryFile); err != nil {
		t.Skipf("299 AC#2 has no subject: the embed reports built but the entry does not resolve (%v)", err)
	}

	dataPath := filepath.Join(os.TempDir(), "wisp-299r1-panel-profile")
	h := &recordingModeHandler{}
	disp := &panel.ComposerDispatch{Mode: h}
	mgr := NewPanelManager(disp, assets, dataPath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hh := &hostThreadHarness{}
	if err := hh.bringUp(mgr, ctx); err != nil {
		t.Fatalf("bringUp on the test thread: %v (a window that cannot be created on this box is a red, not a skip)", err)
	}
	defer hh.stopHostThread()
	defer mgr.Destroy()
	if !mgr.IsCreated() {
		t.Fatalf("host reports not created after bring-up")
	}

	head := gitHeadShortForTest(t)
	deadline := time.Now().Add(ac299Budget)
	var last string
	attempts := 0
	for time.Now().Before(deadline) {
		attempts++
		hh.call(func() {
			if w := mgr.currentWindow(); w != nil {
				w.Eval(ac299ProbeJS())
			}
		})
		settle := time.Now().Add(2 * time.Second)
		for time.Now().Before(settle) && time.Now().Before(deadline) {
			if r, ok := ac299LatestReport(h); ok {
				last = r
				if ac299Satisfied(last) {
					break
				}
			}
			time.Sleep(50 * time.Millisecond)
		}
		if last != "" && ac299Satisfied(last) {
			break
		}
	}
	if last == "" {
		t.Fatalf("299 AC#2 RED (head %s, %d probe attempt(s)): the page never reported at all - the door recorded nothing under %q (what DID arrive: %s). The probe rides the product's own inbound door (window.wispDispatch), so a silent page means the document never ran the probe: the entry navigation was never answered, which is exactly the white window this ticket fixes",
			head, attempts, ac299ProbeID, describeRequests(h))
	}

	href, _ := ac299Token(last, "href")
	base, _ := ac299Token(last, "base")
	ready, _ := ac299Token(last, "ready")
	sheets, _ := ac299Int(last, "sheets")
	rules, _ := ac299Int(last, "rules")
	root, _ := ac299Token(last, "root")
	roottext, _ := ac299Int(last, "roottext")
	resjs, _ := ac299Int(last, "resjs")
	rescss, _ := ac299Int(last, "rescss")
	t.Logf("299 AC#2 page report (head %s, %d attempt(s)): %q | href=%q base=%q ready=%q sheets=%d rules=%d root=%q roottext=%d resjs=%d rescss=%d",
		head, attempts, last, href, base, ready, sheets, rules, root, roottext, resjs, rescss)

	if href != panelEntryURL {
		t.Errorf("the page is not on the panel's own synthetic origin: it reports href=%q, want %q (base=%q, ready=%q). The entry navigation either never committed or was answered by nothing - the page cannot resolve its relative ./assets/* references from anywhere else",
			href, panelEntryURL, base, ready)
	}
	if sheets < 1 || rules <= 0 {
		t.Errorf("the page reports %d stylesheet(s) carrying %d parsed rule(s): the ./assets/*.css bytes did not arrive (or did not parse), so the panel renders unstyled - the fetch half of ticket 299 AC#1", sheets, rules)
	}
	if root != "yes" || roottext <= 0 {
		t.Errorf("the page reports #root present=%q with %d character(s) of its own text: the ./assets/*.js bytes either did not arrive or did not run, so the pane shows none of the page's own text nodes - the exact completion criterion ticket 299 AC#1 words (resjs=%d rescss=%d)",
			root, roottext, resjs, rescss)
	}
}
