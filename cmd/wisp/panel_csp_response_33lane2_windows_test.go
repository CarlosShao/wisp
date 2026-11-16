//go:build windows

// Ticket 33's CSP row, read the way the dispatch demands it be read: from the
// BROWSER's handling of the real response header, not from any source string.
//
// HOW THE HEADER IS READ WITHOUT A FORK CHANGE. The fork maps
// AddWebResourceResponseReceived's vtbl slot but wraps no method, so a Go-side
// response-view reader is not reachable from this package. The reachable real
// reading is the enforcement reading: the browser parses the served document's
// CSP response header, and a fetch the header forbids fires a
// securitypolicyviolation event on the document whose fields ARE the parsed
// header - e.effectiveDirective names the directive that cut the request and
// e.originalPolicy is the policy TEXT exactly as the browser received and
// parsed it from the response header. Asserting on e.originalPolicy is
// therefore a reading of the real served header, not of a Go constant; the
// 33-r14 rulers (panel_csp_33r14_windows_test.go) stay as the construction-side
// reconciliation, and this file is the browser-side half.
//
// The probe is deliberately minimal: register a violation listener, fire ONE
// same-origin fetch (which connect-src 'none' must cut before any network or
// filter call happens), and report href + every violation back through
// window.wispDispatch, the product's own inbound door, exactly like the AC#13
// and AC#14 probes do. No CSP in the real response header means no violation
// event, and the assertion reddens on an empty report - an absent policy and a
// broken probe are kept apart by the report's own HREF element.
package main

import (
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
)

// cspProbeRequestID is the requestId the CSP report rides, unique to this file
// so no other probe's answer can be mistaken for it.
const cspProbeRequestID = "csp-33lane2-read"

// cspProbeJS registers the violation listener, fires one same-origin fetch the
// policy must cut, and reports 'HREF@@<enc(href)>@@<enc(dir||policy)>...' back
// through the door that is already in the page. Elements are
// encodeURIComponent-encoded and joined with '@@', a character encodeURIComponent
// escapes, so encoded content can never smuggle the separator; the leading
// literal HREF element lets the reader tell "the page spoke" from silence.
func cspProbeJS() string {
	return "(function(){" +
		"var v=[];" +
		"document.addEventListener('securitypolicyviolation',function(e){" +
		"v.push(encodeURIComponent(e.effectiveDirective+'||'+e.originalPolicy));" +
		"});" +
		"try{fetch('/__csp_probe_33lane2__').catch(function(){});}catch(err){}" +
		"setTimeout(function(){" +
		"window.wispDispatch(" + reportJSEnv(cspProbeRequestID,
		"'HREF@@' + encodeURIComponent(location.href) + '@@' + v.join('@@')") + ");" +
		"},400);" +
		"})();"
}

// cspViolation is one browser-parsed CSP violation: the directive that cut the
// request and the policy text exactly as the browser received it in the real
// response header.
type cspViolation struct {
	directive string
	policy    string
}

// parseCSPReport splits the page's report: element 1 is the document href, the
// rest are encoded 'directive||policy' pairs. The report's raw text goes into
// the failure message verbatim so a malformed report is diagnosable.
func parseCSPReport(t *testing.T, payload string) (href string, violations []cspViolation) {
	t.Helper()
	parts := strings.Split(payload, "@@")
	if len(parts) < 2 || parts[0] != "HREF" {
		t.Fatalf("the CSP report is malformed (want a leading HREF element): %q", payload)
	}
	var err error
	href, err = url.QueryUnescape(parts[1])
	if err != nil {
		t.Fatalf("the report's href does not decode: %q (%v)", parts[1], err)
	}
	for _, part := range parts[2:] {
		dec, derr := url.QueryUnescape(part)
		if derr != nil {
			t.Errorf("a violation element does not decode: %q (%d element(s) total in %q)", part, len(parts)-2, payload)
			continue
		}
		cut := strings.Index(dec, "||")
		if cut < 0 {
			t.Errorf("a decoded violation element carries no '||' split: %q", dec)
			continue
		}
		violations = append(violations, cspViolation{directive: dec[:cut], policy: dec[cut+2:]})
	}
	return href, violations
}

// assertConnectSrcNoneReadFromTheRealHeader is the shared reading: among the
// policies the browser parsed out of THIS document's live CSP set, the policy
// that cut the probe fetch must carry connect-src 'none' verbatim. A document
// can carry MORE than one policy (the served response header plus a page-side
// meta tag, both enforced, intersection wins): every violation is logged with
// its full parsed policy, and the AC's row is satisfied when the policy that
// carries connect-src 'none' is among them. A policy WITHOUT a connect-src
// directive (its fetch cut falls back to default-src) is recorded, not failed -
// the AC's connect-src 'none' row is about the served response header, and that
// policy announces itself by spelling the directive out.
func assertConnectSrcNoneReadFromTheRealHeader(t *testing.T, where string, violations []cspViolation) {
	t.Helper()
	if len(violations) == 0 {
		t.Errorf("%s: the browser fired ZERO securitypolicyviolation events for the probe fetch - no response-header CSP reached this document, so 'CSP present on served documents' fails at the only surface that matters (the browser's own parse)", where)
		return
	}
	carried := 0
	for _, v := range violations {
		if strings.Contains(v.policy, "connect-src 'none'") {
			carried++
			t.Logf("%s: the connect-src 'none' policy the browser parsed out of the real response: %q", where, v.policy)
		} else {
			t.Logf("%s: an additional enforced policy (no explicit connect-src; its fetch cut falls back to default-src): %q", where, v.policy)
		}
	}
	if carried == 0 {
		t.Errorf("%s: %d violation(s) fired under connect-src but NONE of the browser-parsed policies carries connect-src 'none' verbatim (policies: %v) - the served response header's exfiltration cut did not reach the document", where, len(violations), violationPolicies(violations))
	}
}

func violationPolicies(vs []cspViolation) []string {
	out := make([]string, 0, len(vs))
	for _, v := range vs {
		out = append(out, v.policy)
	}
	return out
}

func violationDirectives(vs []cspViolation) []string {
	out := make([]string, 0, len(vs))
	for _, v := range vs {
		out = append(out, v.directive)
	}
	return out
}

// TestAC6CSPOnTheRealPageIsReadFromItsRealResponse33lane2 is the real-page
// reading: with the embed carrying a built bundle, the cold handover ends on the
// entry navigation served by the resource filter, and the browser's own
// enforcement of that document's response-header CSP must cut the probe fetch
// under connect-src 'none'. Named skip when this binary's embed carries no
// bundle (a fresh checkout): the reading has no subject there, and the not-built
// half lives in the notice test below.
func TestAC6CSPOnTheRealPageIsReadFromItsRealResponse33lane2(t *testing.T) {
	assets, err := panel.BuiltinAssets()
	built := err == nil && assets.Built()
	if !built {
		t.Skipf("named skip: this binary's embed carries no page bundle (err=%v, built=%t) - the real-page CSP reading has no subject; rebuild frontend/dist, rebuild the test binary, and re-run", err, built)
	}

	rp, mgr := startPanelForTest(t)
	showAndWait(t, rp)
	h := mgr.disp.Mode.(*recordingModeHandler)
	// The probe must run in the FINAL handover document, not the transient
	// round-trip probe page: settled means readyState complete, twice identical.
	awaitSettledHandoverDocument(t, rp, h)

	evalOnPanelThread(t, rp, cspProbeJS())
	payload := awaitReport(t, h, cspProbeRequestID, "the real page's CSP report")
	href, violations := parseCSPReport(t, payload)
	t.Logf("real-page CSP report: href=%q, %d violation(s), directives=%v", href, len(violations), violationDirectives(violations))

	if href != panelEntryURL {
		t.Errorf("the probe ran in %q, want the entry URL %q - the report must come from the document the filter served, not from a leftover page", href, panelEntryURL)
	}
	assertConnectSrcNoneReadFromTheRealHeader(t, "real page", violations)
}

// TestAC6CSPOnTheNotBuiltNoticeIsReadFromItsRealResponse33lane2 is the
// not-built-notice reading, and it is the one that needed a product change:
// the notice used to be a SetHtml document with no response to carry a header,
// so it was the one served document without the CSP. Now serveNotBuiltNotice
// navigates to panelNotBuiltURL, the filter answers it like any other response,
// and the browser's parse of THAT response header must again cut the probe
// fetch under connect-src 'none'. The href element also proves the notice
// really is the filter-served route (the synthetic origin), not an about:blank
// SetHtml leftover. Named skip when the embed IS built: the notice path is the
// cold start of an unbuilt tree, and this binary cannot reach it.
func TestAC6CSPOnTheNotBuiltNoticeIsReadFromItsRealResponse33lane2(t *testing.T) {
	assets, err := panel.BuiltinAssets()
	built := err == nil && assets.Built()
	if built {
		t.Skipf("named skip: this binary's embed carries a built bundle (err=%v), so coldStartPageHandover serves the entry, not the notice - rebuild the binary from a tree whose dist is the tracked placeholder to give the notice path a subject", err)
	}

	// A manager with NO assets at all forces serveEntry's not-built verdict on
	// every tree state, so this reading does not depend on what the embed baked.
	dataPath := uniqueProfileLeaf("wisp-33lane2-notice-profile")
	t.Cleanup(func() { _ = os.RemoveAll(dataPath) })
	disp := &panel.ComposerDispatch{Mode: &recordingModeHandler{}}
	mgr := NewPanelManager(disp, nil, dataPath)
	rp := startResidentPanel(observe.NewRegistry(), mgr)
	t.Cleanup(rp.stop)
	showAndWait(t, rp)
	h := mgr.disp.Mode.(*recordingModeHandler)

	// The handover's end state here is the notice: href=panelNotBuiltURL carries
	// a title, so the settled predicate accepts it.
	awaitSettledHandoverDocument(t, rp, h)

	evalOnPanelThread(t, rp, cspProbeJS())
	payload := awaitReport(t, h, cspProbeRequestID, "the not-built notice's CSP report")
	href, violations := parseCSPReport(t, payload)
	t.Logf("notice-page CSP report: href=%q, %d violation(s), directives=%v", href, len(violations), violationDirectives(violations))

	if href != panelNotBuiltURL {
		t.Errorf("the notice document's href is %q, want the filter-served notice route %q - a SetHtml base (about:blank) would read empty or about:blank here, and that document carries no response-header CSP", href, panelNotBuiltURL)
	}
	assertConnectSrcNoneReadFromTheRealHeader(t, "not-built notice", violations)
}
