package risk

import (
	"strings"
	"testing"
)

// Ticket 183 (ledger A362/A363) — the referee this ticket was cut for.
//
// Ticket 177 shape A exempts the host-minted re-read pointer by POSITION: the
// normalized span the declared path occupies inside THIS mark stays out of the
// fragment index (provenance.go MarkWithHostPath, taintmatch.go newFragmentIndex
// skip spans). The index and the verification corpus are keyed by STRING SET
// instead: windows are hashed, sorted and DEDUPED (positions are gone), and
// contains() re-verifies every hash hit against strings.Contains over the whole
// normalized body. So one 8-rune spelling that the declared path shares with
// body text OUTSIDE the exempt span is enough to re-index the exemption's own
// windows and the model's reread of the host's own pointer hits R4 again.
//
// The legs below fix the two halves of that claim and cannot be satisfied by
// the shapes ticket 177 already pins: every shape_a_exemption_test.go fixture
// carries the declared path exactly ONCE and its body never spells a window of
// it (183-a1 §5 ①). Rulers for this file:
//
//	grep -n "^func Test" internal/risk/pointer_183_test.go
//	grep -rn "ptr183Twin\|ptr183Almost" internal/risk/pointer_183_test.go
//
// Deliberately NOT asserted here: anything touching Mark's frozen four-arg
// contract sentence, and any parameter-side by-value pass (ticket 183 AC#5 ①
// forbids it; the W-2/W-3 legs below are the catchers that would go red if a
// by-value or directory-level exemption ever replaced the landed one).

const (
	// ptr183Path is the host-minted pointer in the shape the real CLI actually
	// produced (backslash windows path, artifacts dir, tool-output-agent-task-
	// <uuid>.txt) - ticket 183's live reading, shortened only in the uuid tail.
	ptr183Path = `C:\Users\swq\AppData\Roaming\wisp\artifacts\tool-output-agent-task-c66c0634.txt`

	// ptr183Twin is a >=contractMinFragmentChars window of ptr183Path (the
	// sanity legs below re-derive this, they do not assume it) that the marked
	// body ALSO carries outside the declared span. On the real CLI the carrier
	// was the spilled output's own filler line "...-output-line." meeting the
	// artifact name "tool-output-agent-task-...".
	ptr183Twin = "-output-"

	// ptr183Secret is ordinary tainted body inside the SAME mark as the path:
	// it spells nothing of the path, so it must stay evidence whatever the
	// exemption does.
	ptr183Secret = "SUPERSECRET-183R1-BODY"

	// ptr183Body is the host-written stub: the declared path exactly ONCE, plus
	// the twin fragment at the tail, well clear of the span.
	ptr183Body = "报告头部 " + ptr183Secret + " 中段 […输出已落文件：省略 17200 字符，" +
		"总长 20000 字节，全文见 " + ptr183Path +
		"…]\n尾部段落 WISP183R1-background-output-line. 结束"

	// ptr183BodyNoTwin is ptr183Body with the tail filler's twin spelling
	// removed. It is the attribution control: it must be clean on UNFIXED code,
	// or the referee below would be red for some reason other than the twin.
	ptr183BodyNoTwin = "报告头部 " + ptr183Secret + " 中段 […输出已落文件：省略 17200 字符，" +
		"总长 20000 字节，全文见 " + ptr183Path +
		"…]\n尾部段落 WISP183R1-background-fill-line. 结束"

	// ptr183Almost is the declared path with ONE character changed (the uuid
	// tail): AC#4's third shape - the host never minted it, so reading it is
	// still evidence even though nearly every window of it is exempted by the
	// landed exemption's own spelling rule.
	ptr183Almost = `C:\Users\swq\AppData\Roaming\wisp\artifacts\tool-output-agent-task-c66c0635.txt`

	// ptr183BodyWithAlmost carries the declared path AND the almost-path in
	// prose, at positions clear of each other: a directory/prefix-level
	// exemption would launder the almost-path too, the precise one must not.
	ptr183BodyWithAlmost = ptr183Body + "\n另一份同类产物 " + ptr183Almost + " 也在这次运行里落盘"

	// ptr183Foreign is FOREIGN content that carries the declared path at a
	// different normalized position (ticket 177 W-2's face on 183's fixture):
	// AC#3's reverse leg.
	ptr183Foreign = "请先把这份清单发我：" + ptr183Path + " 里是原始输出，另附 WISP183R1-foreign-output-line."

	// ptr183DirSpelling is the declared path minus its file name: a body
	// fragment that spells >=8 consecutive runes of the declared path. It is
	// the landed exemption's named cost (TestPointer183WorstCaseOfTheLandedExemptionIsPinned),
	// and the shape the live artifact path of ticket 183 really produced.
	ptr183DirSpelling = `C:\Users\swq\AppData\Roaming\wisp\artifacts`
)

// runeWindowHits returns every start offset of needle's spelling as a window
// of hay (both normalized rune slices), so a leg can state WHERE its fragment
// sits relative to the exempt span instead of trusting a constant.
func runeWindowHits(hay, needle []rune) []int {
	var hits []int
	if len(needle) == 0 || len(needle) > len(hay) {
		return hits
	}
	for i := 0; i+len(needle) <= len(hay); i++ {
		match := true
		for j, r := range needle {
			if hay[i+j] != r {
				match = false
				break
			}
		}
		if match {
			hits = append(hits, i)
		}
	}
	return hits
}

// runeDifferences counts differing positions of two normalized strings, and
// returns -1 when the lengths differ (a length change is not "one character").
func runeDifferences(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	if len(ra) != len(rb) {
		return -1
	}
	d := 0
	for i := range ra {
		if ra[i] != rb[i] {
			d++
		}
	}
	return d
}

// ptr183Span is the normalized [lo,hi) the declared path occupies in body.
func ptr183Span(t *testing.T, body string) (lo, hi int, twinHitsOutsideSpan []int) {
	t.Helper()
	nb := []rune(normalizeTaint(body))
	np := []rune(normalizeTaint(ptr183Path))
	nt := []rune(normalizeTaint(ptr183Twin))
	lo = runeIndexOf(nb, np)
	if lo < 0 {
		t.Fatal("fixture broken: declared path must appear in the normalized body")
	}
	if len(runeWindowHits(np, nt)) == 0 {
		t.Fatalf("fixture broken: %q must be a window of the declared path", ptr183Twin)
	}
	hi = lo + len(np)
	for _, j := range runeWindowHits(nb, nt) {
		if j >= hi || j+len(nt) <= lo { // clear of the span, not merely overlapping
			twinHitsOutsideSpan = append(twinHitsOutsideSpan, j)
		}
	}
	return
}

// AC#0 REFEREE (dispatch 183-r1 §2): same mark, declared path present ONCE, and
// the body carrying a path-identical >=8 fragment outside the exempt span => a
// reread of the host's own pointer through fs.read's `path` parameter must not
// hit R4. RED on un-fixed code (that is the point of the leg); the reading is
// kept in docs/evidence/s1/183-pointer-exemption-r1.md §2.
func TestPointer183RefereeTwinFragmentOutsideSpanStaysExempted(t *testing.T) {
	lo, hi, twins := ptr183Span(t, ptr183Body)
	if len(twins) == 0 {
		t.Fatalf("fixture broken: no %q occurrence clear of span [%d,%d)", ptr183Twin, lo, hi)
	}
	p := testProv(t, baseOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", ptr183Body, ptr183Path) {
		t.Fatal("host stub must be marked")
	}
	if h, ok := p.Inspect("task-1", "fs.read", map[string]any{"path": ptr183Path}); ok {
		t.Fatalf("183 referee: the host's own pointer must stay re-readable even when the body spells one of its windows elsewhere, got hit %s", h.String())
	}
}

// Attribution control for the referee, green on UNFIXED code: the identical
// stub whose body carries NO twin fragment is already clean today (this is
// 183-a1's M-D1 reading, made permanent). Without this leg the referee could
// be red because fs.read/path is scanned at all, not because of the twin.
func TestPointer183RefereeControlBodyWithoutTwinIsClean(t *testing.T) {
	if _, _, twins := ptr183Span(t, ptr183BodyNoTwin); len(twins) != 0 {
		t.Fatalf("fixture broken: control body must carry no twin window outside the span, found %v", twins)
	}
	p := testProv(t, baseOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", ptr183BodyNoTwin, ptr183Path) {
		t.Fatal("host stub must be marked")
	}
	if h, ok := p.Inspect("task-1", "fs.read", map[string]any{"path": ptr183Path}); ok {
		t.Fatalf("control: twin-free host stub must be clean on unfixed code (otherwise the referee is red for the wrong reason), got hit %s", h.String())
	}
}

// Second defect face, named separately so a fix cannot be credited for it by
// accident: the SAME declared path present TWICE in the body. The span locator
// returns only the FIRST match (taintmatch.go runeIndexOf), so the second
// occurrence stays indexed and the exemption dies again. Red on unfixed code.
func TestPointer183PathAppearingTwiceInBodyStaysExempted(t *testing.T) {
	body := ptr183BodyNoTwin + "\n宿主再次给出同一条 " + ptr183Path
	nb, np := []rune(normalizeTaint(body)), []rune(normalizeTaint(ptr183Path))
	nt := []rune(normalizeTaint(ptr183Twin))
	spans := runeWindowHits(nb, np)
	if len(spans) < 2 {
		t.Fatalf("fixture broken: declared path must appear twice, found %d", len(spans))
	}
	// Every twin window of this body must sit INSIDE a path occurrence, so the
	// only thing this leg adds over the control is repetition (not the twin).
	for _, j := range runeWindowHits(nb, nt) {
		inside := false
		for _, lo := range spans {
			if j >= lo && j+len(nt) <= lo+len(np) {
				inside = true
			}
		}
		if !inside {
			t.Fatalf("fixture broken: twin window at %d is outside every path occurrence", j)
		}
	}
	p := testProv(t, baseOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", body, ptr183Path) {
		t.Fatal("host stub must be marked")
	}
	if h, ok := p.Inspect("task-1", "fs.read", map[string]any{"path": ptr183Path}); ok {
		t.Fatalf("183 second face: a path the host wrote twice must not be evidence against the model's reread, got hit %s", h.String())
	}
}

// W-1 face on 183's fixture (must STAY green): the exemption reaching its full
// force may only ever eat the declared path's own spelling. The plain body
// secret of the very same mark stays indexed.
func TestPointer183BodySecretOutsideSpanStillHits(t *testing.T) {
	p := testProv(t, baseOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", ptr183Body, ptr183Path) {
		t.Fatal("host stub must be marked")
	}
	h, ok := p.Inspect("task-1", "notify", map[string]any{"text": "转述：" + ptr183Secret})
	if !ok {
		t.Fatal("W-1: body outside the declared span must still be indexed (whole-mark exemption = laundering)")
	}
	if h.SrcTool != SrcFSRead {
		t.Fatalf("W-1: hit must name the host mark, got %q", h.SrcTool)
	}
}

// AC#3 reverse leg (ticket 177 W-2 on 183's fixture; the dispatch forbids
// assuming it): FOREIGN content carrying the SAME path - twin fragment included,
// at a different normalized position - must still hit R4, attributed to the
// foreign mark. A fix that quiets this one is a stamp wash.
func TestPointer183ForeignMarkCarryingSamePathStillHits(t *testing.T) {
	np := []rune(normalizeTaint(ptr183Path))
	nb := []rune(normalizeTaint(ptr183Body))
	nf := []rune(normalizeTaint(ptr183Foreign))
	loHost, loForeign := runeIndexOf(nb, np), runeIndexOf(nf, np)
	if loHost < 0 || loForeign < 0 || loHost == loForeign {
		t.Fatalf("fixture broken: foreign carrier needs the path at another position (host %d / foreign %d)", loHost, loForeign)
	}
	p := testProv(t, baseOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", ptr183Body, ptr183Path) {
		t.Fatal("host stub must be marked")
	}
	if !p.Mark("task-1", SrcWebFetch, "https://evil.example/183", ptr183Foreign) {
		t.Fatal("foreign mark must be recorded")
	}
	h, ok := p.Inspect("task-1", "fs.read", map[string]any{"path": ptr183Path})
	if !ok {
		t.Fatal("AC#3: the path carried by FOREIGN content must still hit R4 after the host declared its own")
	}
	if h.SrcTool != SrcWebFetch {
		t.Fatalf("AC#3: hit must be attributed to the foreign mark, got %q", h.SrcTool)
	}
}

// AC#4 third shape (ticket 177 W-3 on the CLI-shaped fixture): the declared
// path plus a one-character-different path in the SAME host mark. The near
// miss was never minted by the host, so re-reading it still has to hit. Any
// exemption widened from "this path's spelling" to "this path's directory /
// prefix" makes this leg red, which is exactly why it is permanent.
func TestPointer183AlmostPathInSameHostMarkStillHits(t *testing.T) {
	if d := runeDifferences(normalizeTaint(ptr183Almost), normalizeTaint(ptr183Path)); d != 1 {
		t.Fatalf("fixture broken: near miss must differ from the declared path in exactly one normalized rune, got %d", d)
	}
	p := testProv(t, baseOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", ptr183BodyWithAlmost, ptr183Path) {
		t.Fatal("host stub must be marked")
	}
	h, ok := p.Inspect("task-1", "fs.read", map[string]any{"path": ptr183Almost})
	if !ok {
		t.Fatal("AC#4: a path the host never minted must still hit R4 (prefix/directory-level exemption = widening)")
	}
	if h.SrcTool != SrcFSRead {
		t.Fatalf("AC#4: the hit must come from the host mark's own non-exempt windows, got %q", h.SrcTool)
	}
}

// Permanent semantic pin for the landed exemption (dispatch 183-r1 §3 "residual
// leg",钉语义 rather than a window count). The landed family is positional PLUS
// value inside the declaring mark, so the one thing it costs is stated here in
// both directions:
//
//	(1) cost side - body text of THIS mark that spells >=8 consecutive runes of
//	    the declared path (the artifacts directory is the realistic shape, and
//	    the real CLI's filler line "-output-" is the one that actually fired)
//	    stops being evidence against an exfil parameter. RED on unfixed code:
//	    at the start anchor this hit, because the positional span only covered
//	    the first occurrence.
//	(2) boundary side - a body fragment that spells NOTHING of the path is still
//	    evidence from this very mark, so the value rule can never degrade into a
//	    whole-mark exemption (that half is W-1's, re-asserted here on the same
//	    fixture this cost is measured on).
//	(3) confinement side - the same directory spelling carried by a FOREIGN mark
//	    is still evidence: the exemption never leaves the mark that declared it
//	    (W-2/W-4's face on this fixture).
func TestPointer183WorstCaseOfTheLandedExemptionIsPinned(t *testing.T) {
	body := ptr183Body + "\n本次运行还在同一目录 " + ptr183DirSpelling + " 下写过临时文件"
	if !strings.Contains(normalizeTaint(ptr183Path), normalizeTaint(ptr183DirSpelling)) {
		t.Fatal("fixture broken: the pinned cost must spell a real substring of the declared path")
	}
	if _, _, twins := ptr183Span(t, body); len(twins) == 0 {
		t.Fatal("fixture broken: the twin fragment must still be present for this to be the landed shape")
	}
	p := testProv(t, baseOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", body, ptr183Path) {
		t.Fatal("host stub must be marked")
	}
	if h, ok := p.Inspect("task-1", "notify", map[string]any{"text": "转述：" + ptr183DirSpelling}); ok {
		t.Fatalf("183 cost side moved (an extra card is now raised for text that only spells the host's own path): %s", h.String())
	}
	if _, ok := p.Inspect("task-1", "notify", map[string]any{"text": "转述：" + ptr183Secret}); !ok {
		t.Fatal("183 boundary side: a fragment spelling nothing of the declared path must stay evidence from this mark")
	}
	if !p.Mark("task-1", SrcWebFetch, "https://evil.example/dir", "这里也列了 "+ptr183DirSpelling+" 这个目录") {
		t.Fatal("foreign directory mark must be recorded")
	}
	h, ok := p.Inspect("task-1", "notify", map[string]any{"text": "转述：" + ptr183DirSpelling})
	if !ok {
		t.Fatal("183 confinement side: the directory spelling carried by FOREIGN content must still hit R4")
	}
	if h.SrcTool != SrcWebFetch {
		t.Fatalf("183 confinement side: hit must come from the foreign mark, got %q", h.SrcTool)
	}
}

// ptr183BodyAbsentPath is ticket 183 AC#9's fixture: the SAME host-stub family
// as the referee, but the declared path is NOT in the body. What the body does
// carry are windows OF that path - the drive+directory spelling
// (ptr183DirSpelling) and the "-output-" filler - which is exactly what a real
// stub says when it names the artifacts directory without spelling one file.
// Symbol ruler for the branch this pins:
//
//	grep -n "declared path is not in this content" internal/risk/provenance.go
//
// MarkWithHostPath promises that a declaration this body cannot prove excludes
// NOTHING (fail-closed). Before this leg, nothing tested it: 183-v2's mutation
// M3 (attach declaredNorm without requiring runeIndexOf to hit) left every leg
// in this file green with exit=0, so the promise lived only as a comment.
const ptr183BodyAbsentPath = "报告头部 " + ptr183Secret +
	" 中段 […输出已落文件：省略 17200 字符，总长 20000 字节，全文都在 " +
	ptr183DirSpelling + `\tool-output-` +
	" 这一批文件里，具体文件名被截断，下次给出…]\n尾部段落 WISP183R1-background-output-line. 结束"

// AC#9 - the eighth permanent leg, and the only teeth the fail-closed branch
// has. A mark whose body never spells the declared path may not use that path
// as an exemption: the model's reread of it still has to hit R4.
//
// GREEN on today's code (this leg tests that the leg EXISTS and that the branch
// holds, it does not self-certify a defect); RED under 183-v2's M3, which is
// what makes it a catcher rather than a comment. The mechanism M3 breaks: with
// declaredNorm attached unconditionally, EVERY >=8-rune window of the candidate
// path spells the declaration, so contains() skips all of them and the hit
// disappears - even though the body proves nothing about that path.
//
// The two fixture assertions below are the anti-vacuity halves: the leg is only
// measuring the fail-closed branch if (1) the declared path really is absent
// from the normalized body, and (2) the body really does spell at least one
// >=8-rune window of it (otherwise "still hits" would be green because nothing
// could ever hit).
func TestPointer183DeclaredPathAbsentFromBodyStillHits(t *testing.T) {
	nb := []rune(normalizeTaint(ptr183BodyAbsentPath))
	np := []rune(normalizeTaint(ptr183Path))
	nd := []rune(normalizeTaint(ptr183DirSpelling))
	if runeIndexOf(nb, np) != -1 {
		t.Fatal("fixture broken: this leg is the declared-path-ABSENT shape, the path must not occur in the body")
	}
	if len(nd) < contractMinFragmentChars || !strings.Contains(normalizeTaint(ptr183Path), normalizeTaint(ptr183DirSpelling)) {
		t.Fatalf("fixture broken: the pinned fragment must be a >=%d-rune window of the declared path (got %d runes)", contractMinFragmentChars, len(nd))
	}
	if len(runeWindowHits(nb, nd)) == 0 {
		t.Fatal("fixture broken: the body must spell at least one real window of the declared path, or this leg asserts nothing")
	}
	p := testProv(t, baseOptions(t))
	if !p.MarkWithHostPath("task-1", SrcFSRead, "task.output", ptr183BodyAbsentPath, ptr183Path) {
		t.Fatal("host stub must be marked")
	}
	h, ok := p.Inspect("task-1", "fs.read", map[string]any{"path": ptr183Path})
	if !ok {
		t.Fatal("AC#9 fail-closed: a body that never spells the declared path must not exempt a reread of it (if this goes red, MarkWithHostPath's absent-path branch became a comment again)")
	}
	if h.SrcTool != SrcFSRead {
		t.Fatalf("AC#9: the hit must come from this mark's own non-exempt windows, got %q", h.SrcTool)
	}
}
