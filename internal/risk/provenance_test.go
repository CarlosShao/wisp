package risk

import (
	"context"
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

func baseOptions() ProvOptions {
	return ProvOptions{
		NoProbe: true,
		SyncRoots: []SyncRoot{
			{Provider: "OneDrive", Path: `C:\Users\test\OneDrive`, Source: "options"},
		},
	}
}

// --- marking + inspection ----------------------------------------------------

func TestMarkInspectSourceAttribution(t *testing.T) {
	p := testProv(t, baseOptions())
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
	p := testProv(t, baseOptions())
	p.Mark("task-1", SrcWebFetch, "https://a", marker)
	if _, ok := p.Inspect("task-1", "notify", map[string]any{"text": "totally unrelated note"}); ok {
		t.Fatal("clean notify must not hit")
	}
}

func TestMarkEmptyAfterNormalization(t *testing.T) {
	p := testProv(t, baseOptions())
	if p.Mark("task-1", SrcFSRead, "/tmp/x", "   \n\t  ") {
		t.Fatal("content empty after normalization must be rejected")
	}
}

func TestMarkUnknownSourceToolFailClosed(t *testing.T) {
	var logged []string
	old := Logf
	Logf = func(f string, a ...any) { logged = append(logged, f) }
	defer func() { Logf = old }()

	p := testProv(t, baseOptions())
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
	for _, tool := range []string{SrcFSRead, SrcSearchContent, SrcClipboardRead,
		SrcSystemGet, SrcWebFetch, SrcDocRead, SrcScreenCapture, SrcTranscript} {
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
	p := testProv(t, baseOptions())
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
			"path": `C:\Users\test\OneDrive\Notes\shared.md`, "content": "see " + marker}, ChSyncWrite},
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
		"text": "ＭＡＲＫＥＲ QZX-98WVE7-ＴＡＩＮＴＥＤ-SOURCE"}); !ok {
		t.Error("ESCAPIABLE: transformed marker missed on the notify channel")
	} else if h.Channel != ChNotify {
		t.Errorf("transformed marker lost its channel label: %q", h.Channel)
	}
}

// TestR4EndToEndViaAssessor wires the engine through the frozen C19 seam and
// asserts the fused decision: L2, R4, session-override blocked, source named.
func TestR4EndToEndViaAssessor(t *testing.T) {
	p := testProv(t, baseOptions())
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
	p := testProv(t, baseOptions())
	p.Mark("task-1", SrcFSRead, "/secrets.txt", marker)

	if _, ok := p.Inspect("task-1", "fs.write", map[string]any{
		"path": `C:\Users\test\Documents\notes.md`, "content": marker}); ok {
		t.Fatal("local (non-sync) write of tainted content must not be an R4 exfil hit")
	}
	// Same via the tool-name-less frozen seam (shape-based): still no hit.
	if src, hit := p.Detector("task-1").TaintHit(map[string]any{
		"path": `C:\Users\test\Documents\notes.md`, "content": marker}); hit {
		t.Fatalf("adapter must mirror Inspect, got src=%q", src)
	}
	// Unresolvable/missing path on a write is unverifiable -> fail-closed gate.
	if _, ok := p.Inspect("task-1", "fs.write", map[string]any{"content": marker}); !ok {
		t.Fatal("fs.write without a verifiable path must fail-closed to the scan")
	}
}

func TestDetectorUnknownToolScansEverything(t *testing.T) {
	p := testProv(t, baseOptions())
	p.Mark("task-1", SrcClipboardRead, "", "clip "+marker)
	// A plugin tool with no channel entry: generic fail-closed scan of the
	// arbitrary string params.
	src, hit := p.Detector("task-1").TaintHit(map[string]any{
		"payload": map[string]any{"nested": []any{"...", "carries clip " + marker}}})
	if !hit || !strings.Contains(src, "clipboard.read") {
		t.Fatalf("unknown tool must be scanned fail-closed, got %q/%v", src, hit)
	}
}

// --- scope lifetime / DisposalScope (AC #5) ----------------------------------

func TestScopesNeverInherit(t *testing.T) {
	p := testProv(t, baseOptions())
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
	p := NewProvenance(baseOptions())
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
	p := NewProvenance(baseOptions())
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
	p := testProv(t, baseOptions())
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
		"path": `C:\somewhere\plain\a.txt`, "body": "see " + marker}); !ok {
		t.Fatal("ESCAPIABLE: fs.write body param outside writeChannelKeys")
	}
	// Clean calls on named channels stay clean.
	if _, ok := p.Inspect("task-1", "notify", map[string]any{"text": "dinner at 7", "url": "https://x"}); ok {
		t.Fatal("clean notify must not hit")
	}
}

// --- M-4: nesting budget must not be a silent miss ----------------------------

func TestDeepNestingFailsClosedNotSilent(t *testing.T) {
	p := testProv(t, baseOptions())
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
	p := testProv(t, baseOptions())
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
	p := testProv(t, baseOptions())
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
	p := testProv(t, baseOptions())
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
