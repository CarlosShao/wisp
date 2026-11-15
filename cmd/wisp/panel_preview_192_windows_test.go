//go:build windows

package main

// Ticket 192's wiring leg, the criteria. Three readings, one per step of the
// dispatch, each on the PRODUCT path rather than on a rebuild of it:
//
//	step 1  the resource filter's own lookup answers a preview: status, byte
//	        count and the header block a request under panelOrigin would get -
//	        and the same lookup with no preview source answers 404, which is the
//	        pre-wiring shape this ticket was filed against
//	step 2  the resident assembly hands the host that source and the chain the
//	        door, i.e. the wiring is in the composition root and not only here
//	step 3  the page's own door (PanelManager.dispatchRaw -> ComposerDispatch ->
//	        the message handler) reaches PreviewResolver.QuoteForModel, and the
//	        stamp it leaves is visible to the same C25 engine the egress checks
//	        read
//
// Nothing here re-implements a preview rule: the classes, the caps, the
// sandbox policy and the proof a quote is really in the artifact are
// internal/panel's and are judged in internal/panel/preview_192_test.go. What
// these cases judge is assembly - "can the product reach it, and can a claim
// about it be read off a run".

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/risk"
)

// preview192HTML is the previewed document. Small enough that the HTML class
// returns it whole, which is what makes "the byte count the response carries" a
// number a reader can check against the file.
func preview192HTML() []byte {
	return []byte("<!DOCTYPE html>\n<html><head><title>t192</title></head>\n" +
		"<body><h1>artifact preview</h1><p>ticket 192 wiring leg</p></body></html>\n")
}

// preview192SeedArtifacts opens the store over dir, writes the named artifacts
// and CLOSES it again, so the source under test opens the one handle the
// process would own. Returns the bytes written per name.
func preview192SeedArtifacts(t *testing.T, dir string, files map[string][]byte) {
	t.Helper()
	store, err := memory.Open(dir)
	if err != nil {
		t.Fatalf("memory.Open(%q): %v", dir, err)
	}
	for name, data := range files {
		if err := store.PutArtifact(context.Background(), name, data); err != nil {
			_ = store.Close()
			t.Fatalf("PutArtifact(%q): %v", name, err)
		}
	}
	if err := store.Close(); err != nil {
		t.Fatalf("closing the seeding store: %v", err)
	}
}

// preview192Source builds the production source over dir and closes the store
// it opens when the case ends (production keeps it for the process's life; a
// temp root does not have that luxury on Windows).
func preview192Source(t *testing.T, dir string) *panelPreviewSource {
	t.Helper()
	src := newPanelPreviewSource192(dir, nil)
	t.Cleanup(func() {
		if src.store != nil {
			_ = src.store.Close()
		}
	})
	return src
}

// Test192PreviewLookupAnswersTheProductPath is step 1. Two readings from the
// same function, differing only in whether a preview source was handed in:
//
//	no source   404, 0 byte   <- what the product did before this leg
//	source      200, N bytes  <- the user can open the artifact
//
// and the header block is compared BYTE FOR BYTE with the resolver's own
// PreviewDocument.Headers, which is the only way to tell "the preview's
// sandbox policy travelled" apart from "the panel's fallback policy was
// published over a preview and it merely looks served".
func Test192PreviewLookupAnswersTheProductPath(t *testing.T) {
	dir := t.TempDir()
	want := preview192HTML()
	preview192SeedArtifacts(t, dir, map[string][]byte{"page.html": want})
	src := preview192Source(t, dir)

	uri := panelOrigin + panel.PreviewURL("page.html")

	// BEFORE: the filter as it was built for every pre-192 host.
	before := panelAssetLookup(nil, nil, uri)
	t.Logf("step 1 before: panelAssetLookup(no preview source, %q) = status %d %q, %d byte(s), %d header byte(s)",
		uri, before.Status, before.Reason, len(before.Data), len(before.Headers))
	if before.Status != 404 {
		t.Errorf("the no-source lookup answered %d, want the pre-wiring 404: the 'before' half of this reading is not the shape the ticket was filed against", before.Status)
	}
	if len(before.Data) != 0 {
		t.Errorf("the no-source lookup carried %d byte(s): the 404 this ticket is about is an EMPTY miss, not a served preview", len(before.Data))
	}

	// AFTER: the same call, one option further.
	after := panelAssetLookup(nil, src, uri)
	t.Logf("step 1 after:  panelAssetLookup(preview source, %q) = status %d %q, %d byte(s), content-type %q",
		uri, after.Status, after.Reason, len(after.Data), after.ContentType)
	t.Logf("step 1 CSP (verbatim, %d byte(s)): %q", len(after.Headers), after.Headers)
	if after.Status != 200 {
		t.Fatalf("the product path answered %d %q for a preview that exists: the page still cannot open it", after.Status, after.Reason)
	}
	if len(after.Data) != len(want) {
		t.Errorf("the preview response carries %d byte(s), want the artifact's own %d: the page would draw a different document than the one on disk", len(after.Data), len(want))
	}
	if !strings.Contains(string(after.Data), "<h1>artifact preview</h1>") {
		t.Errorf("the preview response bytes are not the document: %q", head192(after.Data, 96))
	}

	// The policy: the stack's own sandbox line, verbatim.
	if !strings.Contains(after.Headers, "Content-Security-Policy: sandbox; default-src 'none'") {
		t.Errorf("the preview response header block carries no 'Content-Security-Policy: sandbox; default-src 'none'': %q. Without the sandbox directive the previewed document keeps this window's origin and can postMessage as the composer - AC#2's whole subject", after.Headers)
	}
	// And it is NOT the panel's fallback: that policy is on every other
	// response this filter answers, and finding it here would mean the preview
	// inherited the panel's policy instead of its own.
	if strings.Contains(after.Headers, panelCSPHeaderValue) {
		t.Errorf("the preview response carries the PANEL's fallback policy (%q): the class policy did not travel, so the HTML sandbox is not in force", panelCSPHeaderValue)
	}
	fallback := panelResponseHeaders(200, after.ContentType)
	if fallback == after.Headers {
		t.Errorf("the preview response header block equals the panel's fallback set (%q): previewResponseHeaders' headers were dropped", fallback)
	}
	// The one decisive comparison: byte equality with what the resolver built.
	resolver, err := src.resolver192()
	if err != nil {
		t.Fatalf("the source opened no resolver: %v", err)
	}
	doc := resolver.Lookup(context.Background(), panel.PreviewURL("page.html"))
	if doc.Status != 200 {
		t.Fatalf("the resolver itself answered %d (%s): the comparison below has no subject", doc.Status, doc.Reason)
	}
	if wantHeaders := strings.Join(doc.Headers, "\r\n"); after.Headers != wantHeaders {
		t.Errorf("the response header block is not PreviewDocument.Headers:\n got %q\nwant %q", after.Headers, wantHeaders)
	}
	if len(doc.Headers) == 0 {
		t.Error("the resolver handed back no headers at all, so the equality above proves nothing")
	}
	// No filesystem path may reach the page (the rule the stack's own criterion
	// pins; repeated here because this is the response that leaves the process).
	if strings.Contains(after.Headers, dir) || strings.Contains(string(after.Data), dir) {
		t.Errorf("the preview response carries the artifacts directory path %q; the browser must never learn where the bytes live", dir)
	}
}

// Test192PreviewRefusalKeepsTheFallbackPolicy is the other half of step 1: a
// preview that cannot be served is refused with the resolver's sentence AND
// still served under a policy. A preview refusal carries no class, so no class
// policy applies and the host's fallback set must: the alternative is a body
// reaching the window with no policy at all.
func Test192PreviewRefusalKeepsTheFallbackPolicy(t *testing.T) {
	dir := t.TempDir()
	src := preview192Source(t, dir)
	uri := panelOrigin + panel.PreviewURL("no-such-artifact.txt")

	resp := panelAssetLookup(nil, src, uri)
	t.Logf("step 1 refusal: %q = status %d, %d body byte(s), header block %q",
		uri, resp.Status, len(resp.Data), resp.Headers)
	if resp.Status != 404 {
		t.Fatalf("a missing artifact answered %d, want 404", resp.Status)
	}
	if resp.Headers != "" {
		t.Errorf("a refusal carried a class header block (%q): the resolver clears Headers exactly so the host's fallback applies", resp.Headers)
	}
	if headers := panelResponseHeaders(resp.Status, resp.ContentType); !strings.Contains(headers, panelCSPHeaderValue) {
		t.Errorf("the fallback set a refusal would be served under carries no CSP: %q", headers)
	}
	if !strings.Contains(string(resp.Data), "不在 artifacts 目录里") {
		t.Errorf("the refusal body is %q, want the resolver's own sentence so the user reads why the preview is empty", head192(resp.Data, 120))
	}
}

// Test192ResidentAssemblyWiresThePreviewStack is step 2: the composition root,
// not the test, hands both surfaces the one source. It assembles through
// newResidentPanelManager - the call the resident boot makes - and reads what
// came out of it.
func Test192ResidentAssemblyWiresThePreviewStack(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	host, _, _, _, err := newResidentPanelManager(dir, nil)
	if err != nil {
		t.Fatalf("newResidentPanelManager refused to assemble: %v", err)
	}
	if host == nil {
		t.Fatal("newResidentPanelManager returned no host")
	}
	if host.preview == nil {
		t.Fatal("the resident host carries no preview source: every /__artifact/ request would 404 through the embed, which is the pre-192 shape")
	}
	src, ok := host.preview.(*panelPreviewSource)
	if !ok {
		t.Fatalf("the host's preview source is %T, want *panelPreviewSource (the process's own, not a rebuild)", host.preview)
	}
	if src.dir != dir {
		t.Errorf("the preview source serves %q, want the boot's data root %q: a preview must read the same artifacts tree every tool write lands in", src.dir, dir)
	}
	if host.disp == nil || host.disp.Message == nil {
		t.Fatal("the resident chain has no message handler attached, so panel.message.send would be refused by the router and QuoteForModel would have no production caller")
	}
	// Ticket 453 wrapped the socket: the chain's message handler is now the
	// composite door, and its PREVIEW half must still be the 192 door over the
	// same source. Both halves are read here so the wrapper cannot be wired to a
	// second preview source - or drop the preview branch - without this reddening.
	door, ok := host.disp.Message.(*panelMessageDoor453)
	if !ok {
		t.Fatalf("the resident chain's message socket is %T, want *panelMessageDoor453", host.disp.Message)
	}
	if door.preview == nil {
		t.Fatal("the composite message door carries no 192 preview door: preview quotes would lose their only QuoteForModel caller")
	}
	if door.preview.src != src {
		t.Error("the door and the window were handed two different preview sources: 'which artifact is that' would have two answers")
	}
	// The engine bind belongs to the boot line where the task pipeline appears;
	// before it the door must be unbound (fail-closed), after it bound.
	if got := src.marker192(); got != nil {
		t.Errorf("the assembly already holds a C25 marker (%T): the panel boot runs before the task pipeline exists and must not pretend to have one", got)
	}
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	host.BindPreviewProvenance(prov)
	if got := src.marker192(); got == nil {
		t.Fatal("BindPreviewProvenance left the door unbound: every quote would be refused for having no engine")
	}
}

// Test192MessageDoorStampsThroughQuoteForModel is step 3, end to end through the
// page's own door: an envelope the renderer can post, a real artifact on disk,
// a real C25 engine. Three things are read:
//
//	the door REFUSES the relay by name (ticket 35's socket does not exist here)
//	the stamp is live in the same scope the egress checks consult
//	text that is not in the artifact is refused and leaves no mark
//
// The audit line is also read, because a door that quotes an artifact must not
// echo the quote into the log (the 35r7 raw-leak family).
func Test192MessageDoorStampsThroughQuoteForModel(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	secret := "这一段来自产物，只应该经 C25 盖戳后离开面板"
	preview192SeedArtifacts(t, dir, map[string][]byte{
		"quoted.txt": []byte("前缀一行\n" + secret + "\n后缀一行\n"),
	})
	host, _, _, _, err := newResidentPanelManager(dir, nil)
	if err != nil {
		t.Fatalf("newResidentPanelManager refused to assemble: %v", err)
	}
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	src := host.preview.(*panelPreviewSource)
	t.Cleanup(func() {
		if src.store != nil {
			_ = src.store.Close()
		}
	})
	ctx := context.Background()

	// (a) UNBOUND ENGINE: the same envelope must be refused, not stamped. This
	// runs BEFORE the bind below, on the product shape - the resident boot's
	// panel exists for a while before the task pipeline's C25 engine does.
	reply, err := host.dispatchRaw(ctx, preview192Envelope(t, "r-192-unbound", "quoted.txt", secret))
	t.Logf("step 3 unbound: reply=%q err=%v", reply, err)
	if err == nil {
		t.Fatalf("a quote was accepted with no C25 engine bound (reply %q): the text would reach a model unstamped", reply)
	}
	if !errors.Is(err, panel.ErrPreviewRefused) {
		t.Errorf("the unbound refusal is not the preview stack's own: %v", err)
	}
	if taints := prov.ScopeTaints(panelPreviewQuoteScope); len(taints) != 0 {
		t.Errorf("the unbound refusal still left %d mark(s) in the engine", len(taints))
	}

	// (b) BOUND ENGINE, the positive control for the same envelope.
	host.BindPreviewProvenance(prov)
	reply, err = host.dispatchRaw(ctx, preview192Envelope(t, "r-192-quote", "quoted.txt", secret))
	t.Logf("step 3 quote: reply=%q err=%v", reply, err)
	if err == nil {
		t.Fatalf("the door reported the quote as fully handled (reply %q): the message relay is ticket 35's and does not exist here, so a success is a claim nobody can back", reply)
	}
	if !errors.Is(err, panel.ErrNoHandlerAttached) {
		t.Errorf("the refusal is not the named 'this door is not the relay' one: %v", err)
	}
	hit, ok := prov.CheckText(panelPreviewQuoteScope, risk.ChHTTP, secret)
	if !ok {
		t.Fatalf("after the quote the C25 egress scan on scope %q found nothing: QuoteForModel was not the door this text went through", panelPreviewQuoteScope)
	}
	if hit.SrcTool != panel.PreviewSourceTool {
		t.Errorf("the mark's source is %q, want %q (the artifact door, not fs.read)", hit.SrcTool, panel.PreviewSourceTool)
	}
	if hit.Origin != "artifact:quoted.txt" {
		t.Errorf("the mark's origin is %q, want artifact:quoted.txt", hit.Origin)
	}

	// (c) TEXT THAT IS NOT IN THE ARTIFACT: refused, and nothing stamped.
	laundered := "页面自己编的一句话，产物里没有"
	reply, err = host.dispatchRaw(ctx, preview192Envelope(t, "r-192-launder", "quoted.txt", laundered))
	t.Logf("step 3 laundered: reply=%q err=%v", reply, err)
	if err == nil {
		t.Fatalf("the door accepted text the artifact does not contain (reply %q): the page could choose what the C25 index believes", reply)
	}
	if !errors.Is(err, panel.ErrPreviewRefused) {
		t.Errorf("the laundering refusal is not ErrPreviewRefused: %v", err)
	}
	if _, ok := prov.CheckText(panelPreviewQuoteScope, risk.ChHTTP, laundered); ok {
		t.Error("the refused quote still left a mark in the process engine")
	}

	// (d) THE LOG: the quote must not be in it.
	// dispatchRaw's own replies and the audit sink are separate; the audit line
	// is asserted on the door directly, with its own sink, so this case does not
	// depend on where the boot routes audit output.
	var lines []string
	door := newPanelPreviewQuoteDoor(src, func(format string, args ...any) {
		lines = append(lines, fmt.Sprintf(format, args...))
	}, "t192")
	if derr := door.HandleMessageRequest(ctx, panel.ComposerRequest{
		Method: panel.MethodMessageSend, RequestID: "r-192-audit",
		Source: panel.ComposerRequestSource, Path: panel.PreviewURL("quoted.txt"), Text: secret,
	}); derr == nil {
		t.Fatal("the door handled the quote without refusing the relay")
	}
	joined := strings.Join(lines, "\n")
	t.Logf("step 3 audit: %q", joined)
	if len(lines) == 0 {
		t.Fatal("the door left no audit line: a refused or stamped quote would be untraceable")
	}
	if strings.Contains(joined, secret) {
		t.Errorf("the audit line carries the quoted text: %q", joined)
	}
	if !strings.Contains(joined, "quoted.txt") {
		t.Errorf("the audit line does not name the artifact the quote was proven against: %q", joined)
	}
}

// Test192MessageDoorIsNotASecondRelay reads the shapes the door must NOT answer:
// a plain message (ticket 35's socket, not this one) and a path that is not a
// preview address. Both refuse by name and stamp nothing.
func Test192MessageDoorIsNotASecondRelay(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	preview192SeedArtifacts(t, dir, map[string][]byte{"page.html": preview192HTML()})
	src := preview192Source(t, dir)
	prov := risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	src.bindMarker192(prov)
	door := newPanelPreviewQuoteDoor(src, nil, "t192")
	ctx := context.Background()

	cases := []struct {
		name string
		req  panel.ComposerRequest
	}{
		{"no path at all", panel.ComposerRequest{Method: panel.MethodMessageSend, RequestID: "r-192-plain", Source: panel.ComposerRequestSource, Text: "普通消息"}},
		{"a path that is not a preview", panel.ComposerRequest{Method: panel.MethodMessageSend, RequestID: "r-192-path", Source: panel.ComposerRequestSource, Path: "/index.html", Text: "普通消息"}},
		{"a preview of nothing", panel.ComposerRequest{Method: panel.MethodMessageSend, RequestID: "r-192-empty", Source: panel.ComposerRequestSource, Path: panel.PreviewURL("page.html"), Text: "  "}},
	}
	for _, tc := range cases {
		err := door.HandleMessageRequest(ctx, tc.req)
		if err == nil {
			t.Fatalf("%s: the door accepted %+v", tc.name, tc.req)
		}
		t.Logf("%s: %v", tc.name, err)
	}
	if taints := prov.ScopeTaints(panelPreviewQuoteScope); len(taints) != 0 {
		t.Errorf("%d mark(s) landed for requests that are not a stamped preview quote: %+v", len(taints), taints)
	}
}

// preview192Envelope builds the postMessage text the renderer would send for a
// preview quote: the composer envelope with the preview URL in `path` and the
// selection in `text`, which is the ONE association this leg could read without
// adding a field to the shared envelope.
func preview192Envelope(t *testing.T, requestID, artifact, text string) string {
	t.Helper()
	body, err := json.Marshal(map[string]string{
		"method":    panel.MethodMessageSend,
		"source":    panel.ComposerRequestSource,
		"requestId": requestID,
		"path":      panel.PreviewURL(artifact),
		"text":      text,
	})
	if err != nil {
		t.Fatalf("building the envelope: %v", err)
	}
	return string(body)
}

// head192 is the short, quoted form of a body a failure sentence can carry
// without dumping a whole artifact into the log.
func head192(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "..."
}

// Test192PreWiringReadingIsStillTheSameFunction keeps the "before" half honest
// for a reader who only runs one case: with no preview source, a preview path
// must fall through to the embed and be a miss - not a panic, not a nil deref,
// not a silent 200 with zero bytes. That is the shape the resident boot ships
// in when no artifacts tree was handed in (a host built by any pre-192 test
// still constructs exactly that).
func Test192PreWiringReadingIsStillTheSameFunction(t *testing.T) {
	uri := panelOrigin + panel.PreviewPathPrefix + "whatever.html"
	resp := panelAssetLookup(nil, nil, uri)
	if resp.Status != 404 {
		t.Errorf("a preview path with no preview source answered %d, want 404", resp.Status)
	}
	if resp.Headers != "" {
		t.Errorf("a miss carried header bytes (%q): the host's fallback set is what a miss is served under", resp.Headers)
	}
}
