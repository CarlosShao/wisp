//go:build windows

package main

// Ticket 33 AC#13 item 2, leg 33-r10. RETOOLED BY 299-r1 (2026-10-13): the entry
// handover is a NAVIGATION to the entry's synthetic origin URL now
// (cmd/wisp/panel_host_windows.go serveEntry -> panelEntryURL), not a SetHtml of
// the entry bytes, because ticket 299 AC#1 had to give the document an origin its
// relative ./assets/* references can resolve against. The three questions below
// are the same three AC#13 asked; the third moved one seam out, onto the exact
// lookup the resource filter answers with (panelAssetLookup), so "the final
// handover carries the entry's own content" is still checked against the bytes
// the product serves, not against a call having been made.
//
// What this file is: the ruler AC#13 asks for that runs WITHOUT a window. Until
// now the only case that asked "which document does a cold start end on" was
// TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe in
// panel_resident_windows_test.go - it drives a real WebView2 browser through
// startPanelForTest, and in a checkout whose embed carries no bundle it skips for
// lack of a subject. So "the user ends on the panel, not on the round-trip probe"
// was pinned by nothing that CI can run. This file pins it on the shipped
// handover sequence itself (PanelManager.coldStartPageHandover, split out of
// bringUp by this leg for exactly that reason), against a control that records
// every handover act, in order.
//
// What the assertions ask - capabilities of the FINAL handover, never "how many
// times was SetHtml called":
//
//	(a) the LAST act handed to the control is a navigation to the entry's
//	    synthetic URL (panelEntryURL), not a document. The probe document is
//	    CAPTURED AT RUNTIME by calling firstRoundTrip against a throwaway sink,
//	    not copied from a literal in this file, so rewriting the probe page's
//	    markup cannot quiet the ruler (a content anchor, not a line-number anchor
//	    and not a string the product could drift away from unnoticed).
//	(b) the probe document is still the FIRST act: cold "usable" is decided on
//	    the round trip, and AC#13's fix is the ORDER (probe first, entry last),
//	    not the probe's removal.
//	(c) the entry URL the final act navigates to is answered, by the filter's own
//	    lookup, with the resolved entry's bytes: status 200, the entry's content
//	    type, lookable body content, and at least one of the entry's own element
//	    ids. That is the "最终交付里含 embed 入口的真内容" half, asked of the
//	    resolved entry through the same Go seam the product serves from
//	    (panel.Assets + panelAssetLookup).
//
// Deliberately absent: no count of SetHtml calls, no line numbers, no ordering
// read out of the source text, no real window, no -tags winlive.

import (
	"bytes"
	"context"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/panel"
	webview2 "github.com/jchv/go-webview2"
)

// docAct33r10 is one handover act, in the order the control received it: a
// document (SetHtml) or a navigation (Navigate). 299-r1 added the ordered log
// because the entry handover is a navigation now, and "which act was last" is
// the question AC#13's order fix is measured on.
type docAct33r10 struct {
	kind    string // "doc" | "nav"
	payload string
}

// docSink33r10 is a webview2.WebView that keeps every document and navigation it
// is handed, in order. It also stands in for the engine's one behaviour the cold
// probe depends on: the inline script of a freshly shown document runs the moment
// it is shown, so a bound no-argument function (the shape firstRoundTrip
// registers) is called there. That is modelling, not measuring - no assertion
// below reads anything but the recorded acts and the embed's own bytes.
type docSink33r10 struct {
	mu        sync.Mutex
	docs      []string
	acts      []docAct33r10
	initJS    []string
	evalJS    []string
	navs      []string
	destroyed int
	bindings  map[string]interface{}
}

func newDocSink33r10() *docSink33r10 {
	return &docSink33r10{bindings: map[string]interface{}{}}
}

func (c *docSink33r10) SetHtml(html string) {
	c.mu.Lock()
	c.docs = append(c.docs, html)
	c.acts = append(c.acts, docAct33r10{kind: "doc", payload: html})
	var pageCalls []func() string
	for _, f := range c.bindings {
		if g, ok := f.(func() string); ok {
			pageCalls = append(pageCalls, g)
		}
	}
	c.mu.Unlock()
	for _, g := range pageCalls {
		g()
	}
}

func (c *docSink33r10) Bind(name string, f interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.bindings == nil {
		c.bindings = map[string]interface{}{}
	}
	c.bindings[name] = f
	return nil
}

func (c *docSink33r10) Init(js string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.initJS = append(c.initJS, js)
}

func (c *docSink33r10) Eval(js string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evalJS = append(c.evalJS, js)
}

func (c *docSink33r10) Navigate(url string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.navs = append(c.navs, url)
	c.acts = append(c.acts, docAct33r10{kind: "nav", payload: url})
}

func (c *docSink33r10) Destroy() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.destroyed++
}

func (c *docSink33r10) Run()       {}
func (c *docSink33r10) Terminate() {}
func (c *docSink33r10) Dispatch(f func()) {
	if f != nil {
		f()
	}
}
func (c *docSink33r10) Window() unsafe.Pointer               { return nil }
func (c *docSink33r10) SetTitle(title string)                {}
func (c *docSink33r10) SetSize(w, h int, hint webview2.Hint) {}

func (c *docSink33r10) firstDoc() (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.docs) == 0 {
		return "", false
	}
	return c.docs[0], true
}

// actLog returns every handover act, in order.
func (c *docSink33r10) actLog() []docAct33r10 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]docAct33r10(nil), c.acts...)
}

// attachSink33r10 points a manager at a recorded control the way bringUp does
// after a successful create (m.w = w, created = true), without creating anything.
func attachSink33r10(m *PanelManager, c *docSink33r10) {
	m.mu.Lock()
	m.w = c
	m.created = true
	m.mu.Unlock()
}

var (
	body33r10Re    = regexp.MustCompile(`(?is)<body[^>]*>(.*?)</body>`)
	script33r10Re  = regexp.MustCompile(`(?is)<script[^>]*>.*?</script>|<script[^>]*>`)
	tag33r10Re     = regexp.MustCompile(`<[a-zA-Z]`)
	entryID33r10Re = regexp.MustCompile(`id=["']([A-Za-z0-9_:.\-.]{1,64})["']`)
)

// hasLookableContent reports whether a document's body holds anything other than
// scripts - an element or text a user could actually see. Empty or script-only
// bodies answer false, which is what shape the round-trip probe page is.
func hasLookableContent(doc string) bool {
	inner := doc
	if m := body33r10Re.FindStringSubmatch(doc); m != nil {
		inner = m[1]
	}
	rest := script33r10Re.ReplaceAllString(inner, "")
	if tag33r10Re.MatchString(rest) {
		return true
	}
	return strings.TrimSpace(rest) != ""
}

// entryFromEmbed33r10 reads the entry through the seam the product serves from.
// ok=false is the not-built world (a fresh checkout, where the embed carries only
// the tracked anchor file); it is a different subject, not a pass.
func entryFromEmbed33r10(t *testing.T) (data []byte, ids []string, ok bool) {
	t.Helper()
	assets, err := panel.BuiltinAssets()
	if err != nil {
		t.Logf("AC#13 33-r10 world: panel.BuiltinAssets: %v (not-built branch)", err)
		return nil, nil, false
	}
	if !assets.Built() {
		t.Logf("AC#13 33-r10 world: embed reports Built()=false (not-built branch)")
		return nil, nil, false
	}
	data, _, err = assets.Resolve(panel.EntryFile)
	if err != nil {
		t.Logf("AC#13 33-r10 world: Resolve(%s): %v (not-built branch)", panel.EntryFile, err)
		return nil, nil, false
	}
	seen := map[string]bool{}
	for _, m := range entryID33r10Re.FindAllStringSubmatch(string(data), -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			ids = append(ids, m[1])
		}
		if len(ids) >= 8 {
			break
		}
	}
	t.Logf("AC#13 33-r10 world: embed resolves %d entry byte(s), %d element id(s) %v", len(data), len(ids), ids)
	return data, ids, true
}

// captureProbeDoc33r10 asks the product's own probe step what document it shows,
// on a throwaway sink. The probe page is therefore identified by what
// firstRoundTrip writes right now, not by a string copied into this file.
func captureProbeDoc33r10(t *testing.T, assets *panel.Assets) string {
	t.Helper()
	sink := newDocSink33r10()
	m := NewPanelManager(nil, assets, t.TempDir())
	attachSink33r10(m, sink)
	rtMs := m.firstRoundTrip(context.Background(), time.Now())
	doc, ok := sink.firstDoc()
	if !ok {
		t.Fatalf("the probe step handed no document to the control at all (rtMs=%v): AC#13's subject moved, this ruler needs re-reading", rtMs)
	}
	t.Logf("AC#13 33-r10 capture: probe document is %d byte(s), round trip rtMs=%v", len(doc), rtMs)
	return doc
}

// TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe is ticket 33 AC#13's
// window-free content ruler. It drives the shipped cold-start handover
// (coldStartPageHandover, which is bringUp's tail verbatim) against a control that
// records every handover act in order, then asks of the LAST act: is it the probe
// document's own SetHtml (it must not be), is it the navigation to the entry's
// synthetic URL, and is that URL answered with the resolved entry's own bytes and
// element ids (it must be, when a bundle exists). The probe-first half is asked
// too, because AC#13 fixes an ORDER, not the probe's existence.
//
// AC#13's mutation is exactly the last-act question: put firstRoundTrip back after
// serveEntry and the last act becomes the probe's SetHtml - wrong kind, wrong
// payload - so the case is red rather than merely reordered.
func TestAC13ColdStartPageOverEndsOnEntryContentNotTheProbe(t *testing.T) {
	assets, err := panel.BuiltinAssets()
	if err != nil {
		assets = nil
	}
	probeDoc := captureProbeDoc33r10(t, assets)

	entryBytes, entryIDs, entryOK := entryFromEmbed33r10(t)

	sink := newDocSink33r10()
	m := NewPanelManager(nil, assets, t.TempDir())
	attachSink33r10(m, sink)

	rtMs := m.coldStartPageHandover(context.Background(), time.Now())
	acts := sink.actLog()
	if len(acts) == 0 {
		t.Fatalf("the cold-start handover handed the control nothing at all (rtMs=%v): nothing was served, AC#13's subject moved", rtMs)
	}
	last := acts[len(acts)-1]
	t.Logf("AC#13 33-r10 read: %d handover act(s), first=%s:%q last=%s:%q rtMs=%v entryOK=%v",
		len(acts), acts[0].kind, head33r10(acts[0].payload), last.kind, head33r10(last.payload), rtMs, entryOK)

	if entryOK {
		if last.kind != "nav" || last.payload != panelEntryURL {
			t.Errorf("after a cold start the LAST handover act is %s:%q, want a navigation to panelEntryURL %q (ticket 299 AC#1's supply shape: the page's relative ./assets/* references need the synthetic origin, so the entry is navigated to, not SetHtml'd). A last act that is a document is ticket 33 AC#13: the user ends on the probe page instead of the panel.",
				last.kind, head33r10(last.payload), panelEntryURL)
		}
		if acts[0].kind != "doc" || acts[0].payload != probeDoc {
			t.Errorf("the FIRST handover act is no longer the round-trip probe document (%s:%q) - the probe stays, it is where cold 'usable' is decided, and AC#13's fix is the ORDER (probe first, entry last), not the probe's removal",
				acts[0].kind, head33r10(acts[0].payload))
		}

		// The bytes the final act's URL is served from: the resource filter's own
		// lookup, asked the same question a request under panelOrigin asks.
		data, ctype, status, reason := panelAssetLookup(assets, panelEntryURL)
		t.Logf("AC#13 33-r10 serving seam: panelAssetLookup(%q) = status %d %q, %d byte(s), content-type %q", panelEntryURL, status, reason, len(data), ctype)
		if status != 200 {
			t.Errorf("the entry URL the final handover navigates to is answered %d %q by the filter's own lookup: the navigation would land on nothing", status, reason)
		}
		if !strings.Contains(strings.ToLower(ctype), "text/html") {
			t.Errorf("the entry URL resolves as content type %q, want an html document", ctype)
		}
		if !bytes.Equal(data, entryBytes) {
			t.Errorf("the entry URL serves %d byte(s) that are not the resolved entry's own %d byte(s): the final handover does not carry the embed's entry content", len(data), len(entryBytes))
		}
		if !hasLookableContent(string(data)) {
			t.Errorf("the bytes the entry URL serves carry nothing a user could look at: with the scripts stripped the body is empty. The panel window would show a shell")
		}
		hits := 0
		for _, id := range entryIDs {
			if strings.Contains(string(data), `id="`+id+`"`) || strings.Contains(string(data), `id='`+id+`'`) {
				hits++
			}
		}
		if len(entryIDs) > 0 && hits == 0 {
			t.Errorf("the entry URL serves content containing NONE of the %d element id(s) the resolved entry declares %v: the URL is not serving the panel page", len(entryIDs), entryIDs)
		}
		t.Logf("AC#13 33-r10 read: last act navigates to the entry URL and that URL serves %d of %d entry id(s), lookable=%v, bytes-equal-to-resolved-entry=%v",
			hits, len(entryIDs), hasLookableContent(string(entryBytes)), bytes.Equal(data, entryBytes))
	} else {
		t.Logf("AC#13 33-r10 read: this tree's embed resolves no entry, so the entry-content question has no subject here; the last handover act is %s:%q (the not-built notice branch)", last.kind, head33r10(last.payload))
	}
}

func head33r10(s string) string {
	r := []rune(s)
	if len(r) > 120 {
		return string(r[:120])
	}
	return s
}
