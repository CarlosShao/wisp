//go:build windows && winlive

package main

// Ticket 302 AC#1, form 甲 (orchestrator dispatch 302-r1): the three page-reply
// rulers below judge a verdict that only exists when a real WebView2 window is up
// AND the page itself answers back through the product door. They sat in the cli
// tier (--scope=cli), whose hosted-runner world could not supply that receipt, so
// the tier kept evaluating criteria whose input never arrives - ticket 302's whole
// finding. The dispatch ruled form 甲: move the tier, not the assertion. That is
// the framing rule the repo already runs on (the speech tier's staging carrier,
// ticket 409, is the same principle one level up): a criterion is judged by a run
// whose world can supply what the criterion reads.
//
// What moved, verbatim (byte-verified against the pre-move file, assertions
// untouched): AC#13's cold-start ruler with its probe harness (entryIDRe /
// entryIDProbes / probeJS / latestReport / reportCount / awaitEntryReport /
// probeAttemptSettle / awaitSettledHandoverDocument), AC#14 nail 1 (the awaited
// binding reply hop, ac14AwaitJS + TestAC14AwaitedBindingReplyReachesThePage) and
// nail 2 (the Go-side Eval push hop). "no report ... from the page within 15s"
// still judges a missing receipt a red, exactly as written.
//
// Tier mechanics (ticket 111 AC#11's gate keeps this honest): no CI job runs
// -tags winlive, but ci.yml's `go vet -tags winlive ./cmd/wisp/` compiles this
// file on every push, so this is a placement change, not a deletion - the three
// names stay in the compiled face (go test -tags winlive -list) and are evaluated
// where a desktop exists (dev box / self-hosted runner), the same rule
// panel_host_windows_live_test.go's exit clause has run under since 33-r5.
//
// countRequests deliberately did NOT move: panel_unbound_253r6_windows_test.go
// uses it on the default tier, so it stayed in panel_resident_windows_test.go
// with the rest of that family's shared harness.

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/panel"
)

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
//
// RETOOLED BY 299-r1's READING SHAPE (not this probe's body): the entry handover is
// a real NAVIGATION to the entry's synthetic origin since ticket 299 AC#1 (the page's
// relative ./assets/* references need a base URI a SetHtml cannot give them), so the
// document is replaced asynchronously - ~0.5s after Show returns on this box - and a
// single probe fired at the old timing reads the still-current about:blank page and
// honestly answers "0". awaitEntryReport below re-asks, bounded; this probe stays a
// single-shot bitmap so "the ids are present" remains the document's own answer.
func probeJS(ids []string) string {
	parts := make([]string, 0, len(ids)+2)
	parts = append(parts, "var b='';")
	for _, id := range ids {
		parts = append(parts, fmt.Sprintf("b += (document.getElementById(%q) ? '1' : '0');", id))
	}
	parts = append(parts, "window.wispDispatch("+reportJSEnv("ac13-probe", "b")+");")
	return "(function(){ " + strings.Join(parts, " ") + " })();"
}

// latestReport returns the LAST report recorded under requestID, so a re-asked
// probe's answer is the one read.
func latestReport(h *recordingModeHandler, requestID string) (string, bool) {
	var out string
	found := false
	for _, req := range h.all() {
		if req.RequestID == requestID {
			out = req.To
			found = true
		}
	}
	return out, found
}

// reportCount counts the reports recorded under requestID.
func reportCount(h *recordingModeHandler, requestID string) int {
	n := 0
	for _, req := range h.all() {
		if req.RequestID == requestID {
			n++
		}
	}
	return n
}

// awaitEntryReport re-asks the live document, bounded, until a report carries at
// least one present id - or returns the last report the budget produced, which the
// caller judges. The budget is panelThreadWait (a monotonic deadline; each attempt
// waits at most probeAttemptSettle for its own answer before re-asking). This keeps
// ticket 33 AC#13's verdict exactly where it was - the ids must be present in the
// LIVE document, answered by the page itself - while tolerating the asynchronous
// handover 299-r1 introduced; when the entry never arrives, the bounded loop ends
// and the caller's assertion reddens on the last "0".
func awaitEntryReport(t *testing.T, rp *residentPanel, h *recordingModeHandler, probes []string) string {
	t.Helper()
	deadline := time.Now().Add(panelThreadWait)
	attempts := 0
	var last string
	for time.Now().Before(deadline) {
		attempts++
		seen := reportCount(h, "ac13-probe")
		evalOnPanelThread(t, rp, probeJS(probes))
		settle := time.Now().Add(probeAttemptSettle)
		for time.Now().Before(settle) && time.Now().Before(deadline) {
			if reportCount(h, "ac13-probe") > seen {
				if got, ok := latestReport(h, "ac13-probe"); ok {
					last = got
					if strings.Contains(got, "1") {
						t.Logf("AC#13 entry-id report arrived on probe attempt %d: %q", attempts, got)
						return got
					}
				}
				break
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
	t.Logf("AC#13 no probe attempt over %d attempt(s) saw a present id; last report: %q", attempts, last)
	return last
}

// probeAttemptSettle bounds how long one re-ask waits for its own answer before the
// next attempt. It is an instrument bound on a monotonic clock, never a product
// threshold, and a slow attempt is reported through the final reading, not retried
// forever.
const probeAttemptSettle = 2 * time.Second

// awaitSettledHandoverDocument re-asks the live document, bounded, until it reads
// as the state coldStartPageHandover hands over LAST: the entry navigation when the
// embed carries the bundle, the not-built notice otherwise. bringUp issues that
// handover just before Show returns and the document it starts replaces the probe
// page asynchronously - ~0.5s after Show returns on this box, the shape 299-r1
// measured for AC#13's probe. The AC#14 push nail fires the moment showAndWait
// unblocks, so its document.title write can land in the outgoing document while its
// report envelope runs in the incoming one: the page then honestly reports a title
// the push never got to keep (measured 2026-11-01 on the working tree, and the same
// red on CI run 327 two weeks earlier). The settled reading is the page's own
// fingerprint (href|readyState|title), read twice in a row identical, readyState
// complete, and a document that is no longer the transient probe page - the probe
// page is the one handover document that carries no title of its own. A push timed
// after that measures the push, not the handover's transit.
func awaitSettledHandoverDocument(t *testing.T, rp *residentPanel, h *recordingModeHandler) {
	t.Helper()
	fingerprintOf := func(seq int) string {
		id := fmt.Sprintf("ac14-docprobe-%d", seq)
		evalOnPanelThread(t, rp, "window.wispDispatch("+reportJSEnv(id,
			"location.href + '|' + document.readyState + '|' + document.title")+");")
		settle := time.Now().Add(probeAttemptSettle)
		for time.Now().Before(settle) {
			if got, ok := latestReport(h, id); ok {
				return got
			}
			time.Sleep(10 * time.Millisecond)
		}
		return ""
	}
	settled := func(fp string) bool {
		parts := strings.SplitN(fp, "|", 3)
		return len(parts) == 3 && parts[1] == "complete" &&
			(parts[0] == panelEntryURL || parts[2] != "")
	}
	prev := ""
	deadline := time.Now().Add(panelThreadWait)
	for seq := 1; ; seq++ {
		if !time.Now().Before(deadline) {
			t.Fatalf("the panel document never settled on the handover's end state within %v (last fingerprint %q) - a push timed now would measure the handover's transit and not the push, which is a failed measurement, not a pass", panelThreadWait, prev)
		}
		fp := fingerprintOf(seq)
		if fp != "" && fp == prev && settled(fp) {
			t.Logf("AC#14 push nail: the handover document settled on probe attempt %d: %q", seq, fp)
			return
		}
		prev = fp
		time.Sleep(10 * time.Millisecond)
	}
}

// TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe is ticket 33 AC#13's ruler: after a
// real cold start the document in the window must be the embedded page, not the
// round-trip probe the host shows while it checks the message channel.
//
// The question goes to the live document and the document answers: does it carry the
// entry's own elements? That form is what gives this instrument teeth against the exact
// defect - putting the probe step back after the entry handover flips the answer from
// every-id-present to none-present (the re-ask budget in awaitEntryReport spends its
// whole wait on the probe page in that world) - where "SetHtml was called" or "the
// entry bytes were read" would tell the two orders apart just as little as the old
// t.Logf did. 299-r1 retooled the reading shape, not the question: the handover is a
// navigation now, so the probe is re-asked within a bounded, monotonic budget instead
// of fired once.
func TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe(t *testing.T) {
	probes := entryIDProbes(t)
	if len(probes) == 0 {
		t.Skipf("AC#13 has no subject in this tree: the embed resolves no entry, so there is no page content that could be covered. Named skip; the working-tree reading is the one that ran, recorded in docs/evidence/s1/33-panel-host-c27-r5.md")
	}
	rp, mgr := startPanelForTest(t)
	showAndWait(t, rp)
	h := mgr.disp.Mode.(*recordingModeHandler)

	answer := awaitEntryReport(t, rp, h, probes)
	present := strings.Count(answer, "1")
	t.Logf("AC#13 page answer (head %s): %q - %d of %d probe id(s) present in the live document",
		gitHeadShortForTest(t), answer, present, len(probes))
	if present == 0 {
		t.Errorf("after a real cold start the live document contains NONE of the %d element ids the embedded entry declares (the page itself answered %q). That is ticket 33 AC#13: the document in the window is not the embedded entry. Since 299-r1 the entry arrives through a navigation plus the resource filter, so this reading is red either when the order is wrong (the round-trip probe page keeping the last word - the defect AC#13 fixed) or when the entry navigation delivers nothing; the fix AC#13 names is the ORDER: bringUp must hand the entry over AFTER the probe - the probe stays, because it is where cold usable is decided", len(probes), answer)
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

	// The handover's last document must be live before anything is pushed into it:
	// the push and its report envelope would otherwise straddle the asynchronous
	// document replacement and measure the transit instead of the push (the helper
	// names the measurement and the runs that showed it).
	awaitSettledHandoverDocument(t, rp, h)

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
