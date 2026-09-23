package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Acceptance criterion 3 (D15(3) spill): single tool result over the scaled
// threshold -> written in full to artifacts\tool-output-<id>.txt while the
// context keeps head + tail + total length + path; raw over the scaled hard cap
// -> truncated=true. Every threshold scales by the provider context_window.
//
// The anti-hardcoding proof is behavioural, not a constant read: the SAME tool
// output must spill under a tiny window and must NOT spill under the 128k
// reference window, so a value that merely happens to be 4000 could not pass.

// TestSpillThresholdScalesWithWindow is the AC3 headline: it drives Spiller
// with budgets produced by BudgetsFor and shows one payload crossing the
// threshold only because the context window shrank.
func TestSpillThresholdScalesWithWindow(t *testing.T) {
	big := BudgetsFor(128000)
	tiny := BudgetsFor(4096)

	// Reference value frozen at 128k.
	if big.SpillTokens != 4000 {
		t.Fatalf("128k spill threshold = %d, want the frozen 4000", big.SpillTokens)
	}
	// A smaller window must shrink it, proportionally.
	if tiny.SpillTokens >= big.SpillTokens {
		t.Fatalf("4096-window threshold = %d, must be below the 128k %d", tiny.SpillTokens, big.SpillTokens)
	}
	wantTiny := (4096 * 4000) / 128000 // 128
	if tiny.SpillTokens != wantTiny {
		t.Errorf("4096-window threshold = %d, want the proportional %d", tiny.SpillTokens, wantTiny)
	}
	// Head/tail/raw-cap scale too (nothing is allowed to stay hardcoded).
	if tiny.SpillHeadTokens >= big.SpillHeadTokens || tiny.SpillTailTokens >= big.SpillTailTokens {
		t.Errorf("head/tail did not shrink: tiny %d/%d vs big %d/%d",
			tiny.SpillHeadTokens, tiny.SpillTailTokens, big.SpillHeadTokens, big.SpillTailTokens)
	}
	if tiny.RawOutputCapBytes >= big.RawOutputCapBytes {
		t.Errorf("raw cap did not shrink: tiny %d vs big %d", tiny.RawOutputCapBytes, big.RawOutputCapBytes)
	}

	// One payload sits between the two thresholds: > tiny, < big.
	payload := strings.Repeat("a", (wantTiny+20)*4) // a hair over the tiny threshold
	if ApproxTokens(payload) <= tiny.SpillTokens || ApproxTokens(payload) >= big.SpillTokens {
		t.Fatalf("payload tokens %d not between tiny %d and big %d",
			ApproxTokens(payload), tiny.SpillTokens, big.SpillTokens)
	}

	dirBig, dirTiny := sealableTempDir124(t), sealableTempDir124(t)
	if sp, err := NewSpiller(dirBig, big).Prepare("call_x", payload); err != nil || sp.Spilled {
		t.Errorf("128k window: spilled=%v err=%v, want NO spill (below threshold)", sp.Spilled, err)
	}
	sp, err := NewSpiller(dirTiny, tiny).Prepare("call_x", payload)
	if err != nil || !sp.Spilled {
		t.Fatalf("4096 window: spilled=%v err=%v, want spill (threshold shrank past the payload)", sp.Spilled, err)
	}
	// Nothing may leak into the big-window artifacts dir (proves it did not spill).
	if entries, _ := os.ReadDir(dirBig); len(entries) != 0 {
		t.Errorf("128k window wrote artifacts: %v", entries)
	}
}

// TestSpillTokenBoundary pins the D15(3) 4000-token reference boundary exactly.
func TestSpillTokenBoundary(t *testing.T) {
	b := BudgetsFor(128000) // SpillTokens = 4000
	dir := sealableTempDir124(t)
	sp := NewSpiller(dir, b)

	at := strings.Repeat("a", b.SpillTokens*4)     // exactly 4000 tokens
	over := strings.Repeat("a", b.SpillTokens*4+4) // 4001 tokens

	s1, err := sp.Prepare("call_at", at)
	if err != nil {
		t.Fatal(err)
	}
	if s1.Spilled {
		t.Errorf("exactly %d tokens spilled, want kept inline (boundary is >)", b.SpillTokens)
	}
	s2, err := sp.Prepare("call_over", over)
	if err != nil {
		t.Fatal(err)
	}
	if !s2.Spilled {
		t.Errorf("%d tokens not spilled, want spill (threshold crossed)", b.SpillTokens+1)
	}
}

// TestSpillArtifactAndStubShape checks the on-disk artifact holds the full
// output and the context stub carries head + tail + totals + path.
func TestSpillArtifactAndStubShape(t *testing.T) {
	b := BudgetsFor(128000)
	dir := sealableTempDir124(t)
	sp := NewSpiller(dir, b)

	// A distinctive, comfortably-over-threshold payload so head/tail are
	// unambiguous.
	head := "HEADSTART-"
	tail := "-TAILFINISH"
	fill := strings.Repeat("0123456789abcdef ", b.SpillTokens) // ~20k tokens worth
	full := head + fill + tail

	s, err := sp.Prepare("call_shape", full)
	if err != nil {
		t.Fatal(err)
	}
	if !s.Spilled {
		t.Fatal("expected spill")
	}
	if s.Path == "" {
		t.Fatal("no artifact path")
	}
	if s.Path != filepath.Join(dir, "tool-output-call_shape.txt") {
		t.Errorf("artifact path = %s, want the D15(3) name under the dir", s.Path)
	}
	// Artifact content is the FULL original output (model can re-read it whole).
	got, err := os.ReadFile(s.Path)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	if string(got) != full {
		t.Errorf("artifact holds %d bytes, want the full %d-byte output", len(got), len(full))
	}

	// Context stub: starts with the kept head, ends with the kept tail, and
	// names the total length + the path in the middle marker.
	if !strings.HasPrefix(s.Text, head) {
		t.Errorf("stub lost the head:\n%.80q", s.Text)
	}
	if !strings.HasSuffix(s.Text, tail) {
		t.Errorf("stub lost the tail: %q…", s.Text[max(0, len(s.Text)-80):])
	}
	for _, needle := range []string{"输出已落文件", fmt.Sprintf("%d 字节", s.TotalBytes), s.Path} {
		if !strings.Contains(s.Text, needle) {
			t.Errorf("stub missing %q", needle)
		}
	}
	if ApproxTokens(s.Text) >= ApproxTokens(full) {
		t.Errorf("stub is not smaller than the original (%d >= %d tokens)",
			ApproxTokens(s.Text), ApproxTokens(full))
	}
	if s.KeptHead != b.SpillHeadTokens || s.KeptTail != b.SpillTailTokens {
		t.Errorf("kept head/tail = %d/%d, want %d/%d",
			s.KeptHead, s.KeptTail, b.SpillHeadTokens, b.SpillTailTokens)
	}
}

// TestSpillRawHardCap checks the scaled 1MB ceiling and the never-split-a-rune
// cut, tested directly on CapRaw (mechanism) with a hand-built small budget.
func TestSpillRawHardCap(t *testing.T) {
	sp := NewSpiller(t.TempDir(), Budgets{RawOutputCapBytes: 100})

	// Under the cap: unchanged, not flagged.
	if c, trunc := sp.CapRaw(strings.Repeat("x", 90)); trunc || len(c) != 90 {
		t.Errorf("90<=100: got len %d trunc %v", len(c), trunc)
	}
	// Over the cap: cut to the ceiling and flagged.
	c, trunc := sp.CapRaw(strings.Repeat("x", 140))
	if !trunc || len(c) != 100 {
		t.Fatalf("over-cap: got len %d trunc %v, want 100/true", len(c), trunc)
	}
	// The cut never splits a multi-byte rune: a 2-byte é straddles byte 100.
	runey := strings.Repeat("x", 99) + "é" + strings.Repeat("y", 50) // é at [99:101]
	c, trunc = sp.CapRaw(runey)
	if !trunc {
		t.Fatal("expected truncation")
	}
	if !utf8StartSafe(c) { // c must be valid, complete UTF-8
		t.Errorf("cap split a rune: %q", c)
	}
	if len(c) != 99 {
		t.Errorf("rune-safe cut = %d bytes, want 99 (é excluded)", len(c))
	}
}

// TestSpillThroughLoop is the end-to-end proof: the same golden stream spills
// only when the task's context_window is tiny, and the spill lands on the real
// path (the tool result that feeds round 2 is the stub, and the artifact file
// exists).
func TestSpillThroughLoop(t *testing.T) {
	// 128k reference window: the ~300-token echo result is well under the 4000
	// threshold, so nothing spills.
	hBig := newHarness(t, "spill-tool")
	resBig := hBig.run("把长结果整理一下")
	if resBig.Status != StatusCompleted {
		t.Fatalf("big-window status = %s (%s)", resBig.Status, resBig.Message)
	}
	if len(resBig.ToolLog) != 1 || resBig.ToolLog[0].Spilled {
		t.Fatalf("128k window should NOT spill, got %+v", resBig.ToolLog)
	}

	// Tiny window (4096): the scaled threshold (128 tokens) drops below the same
	// result, so it must spill through the real loop path.
	dir := sealableTempDir124(t)
	hTiny := newHarness(t, "spill-tool", withConfig(func(c *Config) {
		c.ContextWindow = 4096
		c.ArtifactsDir = filepath.Join(dir, "artifacts")
		c.PerToolTimeout = 2 * time.Second
	}))
	if got := hTiny.loop.Budgets().SpillTokens; got != 128 {
		t.Fatalf("4096-window spill threshold = %d, want 128", got)
	}
	resTiny := hTiny.run("把长结果整理一下")
	if resTiny.Status != StatusCompleted {
		t.Fatalf("tiny-window status = %s (%s)", resTiny.Status, resTiny.Message)
	}
	if len(resTiny.ToolLog) != 1 || !resTiny.ToolLog[0].Spilled {
		t.Fatalf("tiny window must spill the same result, got %+v", resTiny.ToolLog)
	}
	log := resTiny.ToolLog[0]
	// Artifact written under the loop's artifacts dir, named by the call id.
	artifact := filepath.Join(dir, "artifacts", "tool-output-call_sp1.txt")
	if body, err := os.ReadFile(artifact); err != nil || len(body) == 0 {
		t.Errorf("spill artifact missing/empty (%v)", err)
	}
	if log.Artifact != artifact {
		t.Errorf("log artifact = %q, want %q", log.Artifact, artifact)
	}
	// The stub (not the raw output) is what fed the next round.
	if !strings.Contains(log.Text, "输出已落文件") || !strings.Contains(log.Text, artifact) {
		t.Errorf("context stub shape wrong:\n%s", log.Text)
	}
	if !historyHasResultContaining(hTiny.loop.History(), "输出已落文件") {
		t.Errorf("round-2 history carries the raw output, not the stub:\n%s",
			dumpHistory(hTiny.loop.History()))
	}
}

// TestSpillArtifactRespectsRawCap closes the escape the review found (MINOR-1):
// the artifact used to be written from the UNCAPPED output, so with a small
// window (scaled ceiling) the file could be many times the hard cap while the
// stub announced the capped length - and the truncated=true marker, baked into
// the text, was cut away by the tail window. PLAN:432's order is truncate
// FIRST, then land the file.
func TestSpillArtifactRespectsRawCap(t *testing.T) {
	dir := sealableTempDir124(t)
	b := Budgets{RawOutputCapBytes: 200, SpillTokens: 40, SpillHeadTokens: 10, SpillTailTokens: 10}
	sp := NewSpiller(dir, b)

	full := strings.Repeat("x", 800)
	s, err := sp.Prepare("call_cap", full)
	if err != nil {
		t.Fatal(err)
	}
	if !s.TruncatedRaw || !s.Spilled {
		t.Fatalf("truncatedRaw=%v spilled=%v, want both true (800B over a 200B cap)",
			s.TruncatedRaw, s.Spilled)
	}
	body, err := os.ReadFile(s.Path)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	if len(body) > b.RawOutputCapBytes {
		t.Errorf("artifact is %d bytes, over the %d-byte hard cap (escape reproduced)",
			len(body), b.RawOutputCapBytes)
	}
	if len(body) != s.TotalBytes {
		t.Errorf("artifact is %d bytes but the stub announces 总长 %d 字节", len(body), s.TotalBytes)
	}
	// The announced numbers must be the numbers in the stub text.
	for _, needle := range []string{fmt.Sprintf("%d 字节", len(body)), fmt.Sprintf("%d 字节", s.TotalBytes)} {
		if !strings.Contains(s.Text, needle) {
			t.Errorf("stub missing %q:\n%s", needle, s.Text)
		}
	}
	// The truncated=true marker survives a 10-token tail window: it is context
	// text, not artifact bytes.
	if !strings.Contains(s.Text, "truncated=true") {
		t.Errorf("stub lost the truncated=true marker:\n%s", s.Text)
	}
	if !strings.Contains(s.Text, "硬上限") {
		t.Errorf("stub must name the ceiling that cut the output:\n%s", s.Text)
	}
}

// ---------------------------------------------------------------------------

// utf8StartSafe is a small validity probe for CapRaw's rune-boundary test.
func utf8StartSafe(s string) bool {
	for i := 0; i < len(s); {
		switch b := s[i]; {
		case b < 0x80:
			i++
		case b&0xE0 == 0xC0:
			if i+1 >= len(s) || s[i+1]&0xC0 != 0x80 {
				return false
			}
			i += 2
		case b&0xF0 == 0xE0:
			if i+2 >= len(s) || s[i+1]&0xC0 != 0x80 || s[i+2]&0xC0 != 0x80 {
				return false
			}
			i += 3
		case b&0xF8 == 0xF0:
			if i+3 >= len(s) || s[i+1]&0xC0 != 0x80 || s[i+2]&0xC0 != 0x80 || s[i+3]&0xC0 != 0x80 {
				return false
			}
			i += 4
		default:
			return false // a stray continuation byte
		}
	}
	return true
}
