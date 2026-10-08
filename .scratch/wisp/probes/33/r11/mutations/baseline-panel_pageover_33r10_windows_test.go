//go:build windows

package main

// Ticket 33 AC#13 item 2, leg 33-r10.
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
// every document it is handed.
//
// What the assertions ask - capabilities of the FINAL document, never "how many
// times was SetHtml called":
//
//	(a) the last document the control was handed is not the document the host's
//	    own probe step writes. The probe document is CAPTURED AT RUNTIME by
//	    calling firstRoundTripLocked against a throwaway sink, not copied from a
//	    literal in this file, so rewriting the probe page's markup cannot quiet the
//	    ruler (a content anchor, not a line-number anchor and not a string the
//	    product could drift away from unnoticed).
//	(b) the last document carries something a person could look at: after the
//	    scripts are stripped, its body still holds an element or text. The probe
//	    page's body is one inline script and nothing else, so a handover that ends
//	    there is red here too - a second, independent tooth.
//	(c) when the embed really resolves an entry, the last document carries that
//	    entry's own bytes and at least one of the element ids parsed out of it.
//	    That is the "最终文档里含 embed 入口的真内容" half, asked of the resolved
//	    entry through the same Go seam the product serves from (panel.Assets).
//
// Deliberately absent: no count of SetHtml calls, no line numbers, no ordering
// read out of the source text, no real window, no -tags winlive.

import (
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

// docSink33r10 is a webview2.WebView that keeps every document it is handed. It
// also stands in for the engine's one behaviour the cold probe depends on: the
// inline script of a freshly shown document runs the moment it is shown, so a
// bound no-argument function (the shape firstRoundTripLocked registers) is called
// there. That is modelling, not measuring - no assertion below reads anything but
// the recorded documents and the embed's own bytes.
type docSink33r10 struct {
	mu        sync.Mutex
	docs      []string
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

func (c *docSink33r10) lastDoc() (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.docs) == 0 {
		return "", false
	}
	return c.docs[len(c.docs)-1], true
}

func (c *docSink33r10) firstDoc() (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.docs) == 0 {
		return "", false
	}
	return c.docs[0], true
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
// firstRoundTripLocked writes right now, not by a string copied into this file.
func captureProbeDoc33r10(t *testing.T, assets *panel.Assets) string {
	t.Helper()
	sink := newDocSink33r10()
	m := NewPanelManager(nil, assets, t.TempDir())
	attachSink33r10(m, sink)
	rtMs := m.firstRoundTripLocked(context.Background(), time.Now())
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
// records documents, then asks of the LAST document: is it the probe's own
// document (it must not be), does it carry lookable content (it must), and does it
// carry the resolved entry's bytes and element ids (it must, when a bundle exists).
//
// AC#13's mutation is exactly the first two questions: put firstRoundTripLocked
// back after serveEntry and the last document becomes the probe page - different
// content, no lookable body - so the case is red rather than merely reordered.
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
	last, ok := sink.lastDoc()
	if !ok {
		t.Fatalf("the cold-start handover handed no document to the control at all (rtMs=%v): nothing was served, AC#13's subject moved", rtMs)
	}

	if last == probeDoc {
		t.Errorf("after a cold start the LAST document handed to the control is byte-for-byte the document the round-trip probe step writes (%d byte(s)). That is ticket 33 AC#13: the user ends on the probe page instead of the panel. Product side: the probe stays - it is where cold usable is decided - but coldStartPageHandover must hand the entry over AFTER it", len(probeDoc))
	}
	if !hasLookableContent(last) {
		t.Errorf("after a cold start the last document carries nothing a user could look at: with its scripts stripped its body is empty (head %q). The panel was never handed over, so the panel window would show a shell. Ticket 33 AC#13 item 2 asks the final document about its content, not about whether a call was made", head33r10(last))
	}
	if entryOK {
		if !strings.Contains(last, string(entryBytes)) {
			t.Errorf("after a cold start the last document does not carry the embedded entry's own %d byte(s): ticket 33 AC#13's 'final document contains the real embed entry content' is not met", len(entryBytes))
		}
		hits := 0
		for _, id := range entryIDs {
			if strings.Contains(last, `id="`+id+`"`) || strings.Contains(last, `id='`+id+`'`) {
				hits++
			}
		}
		if len(entryIDs) > 0 && hits == 0 {
			t.Errorf("after a cold start the last document contains NONE of the %d element id(s) the resolved entry declares %v (head %q): the document in the window is not the panel page", len(entryIDs), entryIDs, head33r10(last))
		}
		t.Logf("AC#13 33-r10 read: last document %d byte(s), %d of %d entry id(s) present, entry bytes carried=%v",
			len(last), hits, len(entryIDs), strings.Contains(last, string(entryBytes)))
	} else {
		t.Logf("AC#13 33-r10 read: this tree's embed resolves no entry, so the entry-content question has no subject here; the last document is %d byte(s) and carries lookable content=%v; it is not the probe document=%v",
			len(last), hasLookableContent(last), last != probeDoc)
	}
}

func head33r10(s string) string {
	r := []rune(s)
	if len(r) > 120 {
		return string(r[:120])
	}
	return s
}
