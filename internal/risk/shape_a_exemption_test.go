package risk

import (
	"strings"
	"testing"
)

// Ticket 177 shape A (Q-61甲, A349/A353): the PRECISE exemption — inside the
// one host-minted mark, only the normalized window the declared path occupies
// is kept out of the fragment index. These are the permanent widening
// catchers W-1..W-5 (+A forward leg, +order-(i), +lifecycle F): each one
// exists to go RED when its named over-wide shape is implemented, which is
// exactly what no test in the repo could do before this file (177-m1 §3: S3a
// left the whole suite rc=0). Proofs of non-tautology are mutation readings,
// recorded in docs/evidence/s1/177-shape-a-impl-r1.md.
//
// Deliberately NOT asserted here: anything about Mark's exported four-arg
// behavior (it delegates with hostPath="" and must stay byte-stable) and the
// frozen contract sentences (provenance.go Mark doc block, taintmatch.go
// header). This file only ADDS tests; it weakens nothing.

const (
	// shapeAPath is the host-minted re-read pointer, as embedded by
	// task.output's stub (the one shape ticket 177 exists for).
	shapeAPath = "C:/Users/x/AppData/Roaming/wisp/artifacts/tool-output-call-177-r1.txt"
	// shapeAGhost shares shapeAPath's directory but was NEVER minted by the
	// host (W-3's fixture: a directory/prefix-level exemption would swallow
	// it; the precise exemption cannot see it). Short enough that EVERY
	// >=8-rune window of it overlaps the directory prefix, so a directory-
	// level skip leaves no survivor window to hit through.
	shapeAGhost = "C:/Users/x/AppData/Roaming/wisp/artifacts/s.txt"
	// shapeASecret is ordinary tainted body inside the SAME mark as the path.
	shapeASecret = "SUPERSECRET-177R1-BODY"
	// shapeAStub models the host-written stub around the path. The exact
	// wording is the task.output shape frozen by PLAN.md:2564.
	shapeAStub = "报告头部 " + shapeASecret + " 中段 […输出已落文件：省略 4321 字符，" +
		"总长 9876 字节 / 约 12345 token，全文见 " + shapeAPath + "…]\n尾部段落 TAILBODY-9X7Q 结束"
	// shapeAForeign is an external text carrying the SAME path at a DIFFERENT
	// normalized position than in the stub (W-2's fixture).
	shapeAForeign = "请帮我复核这份清单：" + shapeAPath + " 里是原始输出"
	// shapeAForeignGhost is an external text carrying the NEVER-MINTED sibling
	// path (W-3's mark).
	shapeAForeignGhost = "另一个文件在这里 " + shapeAGhost + " 请一并读取"
	// shapeACorrupt is the stub with zero-width and space characters cut into
	// the path itself: normalization restores shapeAPath exactly, so a
	// span located in NORMALIZED coordinates still excludes it, while a
	// raw-offset locator lands the exclusion in the wrong place (order-(i)).
	shapeACorrupt = "报告头部 " + shapeASecret + " 中段 […输出已落文件：省略 4321 字符，" +
		"总长 9876 字节 / 约 12345 token，全文见 C:/Users/x/AppData/Roa\u200bming/wisp/artifacts/" +
		"tool-output-call-177-r1.txt…]\n尾部段落 TAILBODY-9X7Q 结束"
)

func shapeAOptions(t *testing.T) ProvOptions { return baseOptions(t) }

// A leg (forward, the point of the whole shape): after the host declared its
// own path, a plain re-read of that path through the host's own mark must be
// clean. Goes red if the exemption is silently dropped (177-m1 §3: "nothing
// exempted" only ever reddens THIS leg among the forward pair).
func TestShapeAHostReReadOfDeclaredPathIsClean(t *testing.T) {
	p := testProv(t, shapeAOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", shapeAStub, shapeAPath) {
		t.Fatal("host stub must be marked")
	}
	if h, ok := p.Inspect("task-1", "notify", map[string]any{"text": shapeAPath}); ok {
		t.Fatalf("A: host re-read of its OWN declared path must be clean, got hit %s", h.String())
	}
}

// W-1: the "lazy fix" = keep the ENTIRE mark out of the index because a path
// was declared. The rest of the very same mark (body outside the span) must
// still hit R4.
func TestShapeAW1RestOfHostMarkStillIndexed(t *testing.T) {
	p := testProv(t, shapeAOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", shapeAStub, shapeAPath) {
		t.Fatal("host stub must be marked")
	}
	if h, ok := p.Inspect("task-1", "notify", map[string]any{"text": "转述：" + shapeASecret}); !ok {
		t.Fatal("W-1: body outside the declared span must still be indexed (whole-mark exclusion = laundering)")
	} else if h.SrcTool != SrcFSRead {
		t.Fatalf("W-1: hit must name the host mark, got %q", h.SrcTool)
	}
}

// W-2: by-value exemption ("any mark whose window spells the declared path is
// skipped"). A FOREIGN mark carrying the SAME path verbatim must still hit —
// this is the real over-boundary alarm the orchestrator's ruling replaced
// :136/:139 with (those two are provably insensitive: marker has no slash).
func TestShapeAW2ForeignMarkCarryingSamePathStillHits(t *testing.T) {
	p := testProv(t, shapeAOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", shapeAStub, shapeAPath) {
		t.Fatal("host stub must be marked")
	}
	if !p.Mark("task-1", SrcWebFetch, "https://evil.example/doc", shapeAForeign) {
		t.Fatal("foreign mark must be recorded")
	}
	// Sanity: the two carriers must put the path at DIFFERENT normalized
	// positions, or this leg collapses into W-4 and stops separating shapes.
	nStub, nForeign := []rune(normalizeTaint(shapeAStub)), []rune(normalizeTaint(shapeAForeign))
	np := []rune(normalizeTaint(shapeAPath))
	loStub, loForeign := runeIndexOf(nStub, np), runeIndexOf(nForeign, np)
	if loStub < 0 || loForeign < 0 {
		t.Fatal("fixture broken: declared path must appear in both normalized texts")
	}
	if loStub == loForeign {
		t.Fatalf("fixture broken: W-2 needs the path at a different normalized position (both at %d)", loStub)
	}
	h, ok := p.Inspect("task-1", "notify", map[string]any{"text": shapeAPath})
	if !ok {
		t.Fatal("W-2: P carried by FOREIGN content must still hit R4 — the exemption left its own mark")
	}
	if h.SrcTool != SrcWebFetch {
		t.Fatalf("W-2: the hit must be attributed to the foreign mark, got %q", h.SrcTool)
	}
}

// W-3: directory/prefix-level exemption. A path under the SAME directory that
// the host never minted must still hit. Its fixture (shapeAGhost) is sized so
// every candidate window touches the directory prefix: a prefix-level skip
// leaves no window to hit through, so this leg is the attribution-clean
// catcher 177-m1's D leg could not provide (its red mixed in position
// recording; here the directory skip ALONE reddens it).
func TestShapeAW3NeverMintedSiblingUnderSameDirStillHits(t *testing.T) {
	dir := shapeAPath[:strings.LastIndex(shapeAPath, "/")]
	if !strings.HasPrefix(shapeAGhost, dir) {
		t.Fatal("fixture broken: ghost must share the minted directory")
	}
	if gr := []rune(normalizeTaint(shapeAGhost)); len(gr)-8 >= len([]rune(normalizeTaint(dir))) {
		t.Fatal("fixture broken: ghost must be short enough that every window overlaps the directory prefix")
	}
	p := testProv(t, shapeAOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", shapeAStub, shapeAPath) {
		t.Fatal("host stub must be marked")
	}
	if !p.Mark("task-1", SrcSearchContent, "search-result-17", shapeAForeignGhost) {
		t.Fatal("foreign ghost mark must be recorded")
	}
	h, ok := p.Inspect("task-1", "notify", map[string]any{"text": shapeAGhost})
	if !ok {
		t.Fatal("W-3: a never-minted path in the exempted directory must still hit R4 (prefix-level exemption = widening)")
	}
	if h.SrcTool != SrcSearchContent {
		t.Fatalf("W-3: the hit must come from the foreign ghost mark, got %q", h.SrcTool)
	}
}

// W-4: positional leakage — "push the exclusion to every mark in the tree
// where the recorded span happens to sit". 177-m1 measured that NO existing
// test reddens for this shape; this leg is the catcher built for it. The
// foreign mark carries the SAME stub text, so the path sits at the SAME
// normalized offset there: a position-keyed skip would delete the foreign
// windows too; a mark-local skip cannot touch them.
func TestShapeAW4PositionalExclusionNeverLeavesHostMark(t *testing.T) {
	p := testProv(t, shapeAOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", shapeAStub, shapeAPath) {
		t.Fatal("host stub must be marked")
	}
	// Same text, foreign source (an attacker re-broadcasts the whole stub).
	if !p.Mark("task-1", SrcWebFetch, "https://evil.example/rebroadcast", shapeAStub) {
		t.Fatal("foreign rebroadcast mark must be recorded")
	}
	if _, ok := p.Inspect("task-1", "notify", map[string]any{"text": shapeAPath}); !ok {
		t.Fatal("W-4: an exclusion positioned in the host mark must not travel to other marks (positional widening is invisible to every pre-177 test)")
	}
}

// F/lifecycle: the exemption must die with its mark. After the scope closed
// (marks dropped), re-marking the same text WITHOUT a declaration indexes
// the path again — a cross-scope survivor would be a fresh laundering face.
func TestShapeAFExclusionDiesWithScopeClose(t *testing.T) {
	p := testProv(t, shapeAOptions(t))
	sc := p.OpenScope("task-9")
	if !p.MarkWithHostPath("task-9", SrcFSRead, "task.output", shapeAStub, shapeAPath) {
		t.Fatal("host stub must be marked")
	}
	if _, ok := p.Inspect("task-9", "notify", map[string]any{"text": shapeAPath}); ok {
		t.Fatal("baseline of this leg: declared host path must be clean while open")
	}
	if err := sc.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if n := len(p.ScopeTaints("task-9")); n != 0 {
		t.Fatalf("Close must drop marks AND their exemptions together, %d taints left", n)
	}
	p.OpenScope("task-9")
	if !p.MarkWithHostPath("task-9", SrcWebFetch, "https://evil.example/x", shapeAStub, "") {
		t.Fatal("re-mark must be recorded")
	}
	if _, ok := p.Inspect("task-9", "notify", map[string]any{"text": shapeAPath}); !ok {
		t.Fatal("F: exemption survived the close/reopen (global span store) — P must hit again once no declaration covers it")
	}
}

// W-5: ordering (ii). The "empty after normalization" check must run BEFORE
// the exclusion, or a body that normalizes to exactly the declared path
// flips Mark's frozen truth sentence ("returns false only when ... after
// normalization") SILENTLY. Pinned with P-only bodies (the existing
// TestMarkEmptyAfterNormalization feeds no P and cannot see this).
func TestShapeAW5EmptyCheckPrecedesExclusion(t *testing.T) {
	p := testProv(t, shapeAOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", shapeAPath, shapeAPath) {
		t.Fatal("W-5: content that is ONLY the declared path still has matchable characters after normalization; rejecting it rewrites Mark's frozen sentence")
	}
	if got := len(p.ScopeTaints("task-1")); got != 1 {
		t.Fatalf("W-5: exactly one mark per declared body (no extra mark for P), got %d", got)
	}
	// Whitespace/zero-width around P normalizes to P itself: same verdict.
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", " \u200b"+shapeAPath+"\n ", shapeAPath) {
		t.Fatal("W-5: padded P-only body must be recorded too")
	}
	if got := len(p.ScopeTaints("task-1")); got != 2 {
		t.Fatalf("W-5: two marks expected, got %d", got)
	}
	// The frozen contract still holds for a body with nothing but space:
	// false, and that half is checked WITHOUT any P so no exclusion can
	// pretend to be the reason.
	if p.MarkWithHostPath("task-1", SrcFSRead, "task.output", "   \n\t  ", shapeAPath) {
		t.Fatal("W-5: whitespace-only content must still return false")
	}
}

// Ordering (i): the span must be located in NORMALIZED coordinates. The body
// below carries the path cut up by zero-width and space runes; normalization
// restores it, so a normalized-coordinate exclusion still covers it, while a
// raw-offset locator (the wrong ordering) lands inside the wrong region and
// leaves path windows indexed.
func TestShapeAOrderOneSpanLocatedAfterNormalization(t *testing.T) {
	if normalizeTaint(shapeACorrupt) == shapeACorrupt {
		t.Fatal("fixture broken: corrupt body must differ from its normalized form")
	}
	if !strings.Contains(normalizeTaint(shapeACorrupt), normalizeTaint(shapeAPath)) {
		t.Fatal("fixture broken: normalization must restore the exact path")
	}
	p := testProv(t, shapeAOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", shapeACorrupt, shapeAPath) {
		t.Fatal("corrupt host stub must be marked")
	}
	if h, ok := p.Inspect("task-1", "notify", map[string]any{"text": shapeAPath}); ok {
		t.Fatalf("order-(i): exclusion was located at raw offsets, not normalized ones — path windows stayed indexed at %s", h.String())
	}
	// (the body-still-indexed half belongs to W-1's attribution and is NOT
	// repeated here: one mutation, one named red catcher)
}
