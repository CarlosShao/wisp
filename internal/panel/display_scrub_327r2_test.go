package panel

// 327 residual R1, in-repo by leg 327-r2; reshaped for ticket 333 by leg 333-r1, which
// deleted ruledRoster327r2 (nail 1's fifteen-literal EXPECTED side) and flipped nail 2
// from excluded to washed. Old test names, both flips and the eight runes itemised by
// name: docs/evidence/s1/333-r1-display-cf-category.md. Nail 1 is a RELATION now - the
// product's predicate over every scalar against this toolchain's unicode.Cf table, so
// a Go or Unicode table move reddens instead of silently reopening the hole. Nail 3 is
// the anchor's own text; whole-category idempotence moved to ticket 333's file beside it.

import (
	"fmt"
	"sort"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

// excludedRoster327r2 is 327's off-ruling eight (U+206B was the omitted one, R1).
// Ticket 333 turns all eight into WASHED; nail 2 reads them one by one.
var excludedRoster327r2 = []rune{0x2061, 0x2062, 0x2063, 0x2064, 0x061C, 0x180E, 0x206A, 0x206B}

func name327r2(r rune) string { return fmt.Sprintf("U+%04X", r) }

func names327r2(m map[rune]bool) []string {
	out := make([]string, 0, len(m))
	for r := range m {
		out = append(out, name327r2(r))
	}
	sort.Strings(out)
	return out
}

// washSetOver collects every scalar a predicate names, skipping the surrogates. Nail 1
// hands it two DIFFERENT predicates, so neither side is a copy of the other.
func washSetOver(names func(rune) bool) map[rune]bool {
	set := map[rune]bool{}
	for r := rune(0); r <= utf8.MaxRune; r++ {
		if r >= 0xD800 && r <= 0xDFFF {
			continue
		}
		if names(r) {
			set[r] = true
		}
	}
	return set
}

// Nail 1: the product's wash set IS the Cf table, in both directions. ACTUAL = the
// product's own predicate, EXPECTED = stdlib's category: two sources, not a self-copy.
func TestDisplayWashRosterIsExactlyTheCfCategory(t *testing.T) {
	actual := washSetOver(isDisplayWashCf)
	expected := washSetOver(func(r rune) bool { return unicode.Is(unicode.Cf, r) })
	extra, missing := map[rune]bool{}, map[rune]bool{}
	for r := range actual {
		if !expected[r] {
			extra[r] = true
		}
	}
	for r := range expected {
		if !actual[r] {
			missing[r] = true
		}
	}
	t.Logf("nail 1: product-washed=%d, unicode.Cf=%d", len(actual), len(expected))
	if len(extra) != 0 { // reddens if the product washes beyond the category
		t.Errorf("the product washes runes the Cf table does not name: %v (product=%d category=%d)",
			names327r2(extra), len(actual), len(expected))
	}
	if len(missing) != 0 { // reddens if the product skips any of the category, other line
		t.Errorf("a Cf-table rune is not washed by the product: %v (product=%d category=%d)",
			names327r2(missing), len(actual), len(expected))
	}
	if len(actual) != len(expected) {
		t.Errorf("product roster cardinality=%d, want the category's %d", len(actual), len(expected))
	}
}

// Nail 2: the eight the old ruling kept raw are washed now - escaped, never deleted.
func TestDisplayWashRosterWashesTheEightRunesTheOldRulingExcluded(t *testing.T) {
	for _, r := range excludedRoster327r2 {
		if !unicode.Is(unicode.Cf, r) {
			t.Errorf("%s is not unicode.Cf - this nail lost its reason for existing", name327r2(r))
		}
		if !isDisplayWashCf(r) { // direction flipped from false to true at ticket 333
			t.Errorf("%s is Cf but isDisplayWashCf still skips it", name327r2(r))
		}
		in := "photo" + string(r) + ".png"
		if !needsDisplayWash(in) {
			t.Errorf("%s is Cf but the fast path would return %q raw", name327r2(r), in)
		}
		got := scrubForDisplay(in)
		if strings.ContainsRune(got, r) {
			t.Errorf("%s survived as itself in %q: the shape is escape, not deletion", name327r2(r), got)
		}
		if want := "photo" + escapeOf(r) + ".png"; got != want {
			t.Errorf("%s washed to %q, want the boundary vocabulary's %q (unchanged input is the other failure)", name327r2(r), got, want)
		}
	}
}

// Nail 3: the two sides are one set; ESC forces the fast path open (⛔ writer-blindness
// hidden by "clean string needs no wash").
func TestDisplayWashJudgeAndWriterAgreeOnBothFamiliesAndBoundaries(t *testing.T) {
	writerEscapes := func(r rune) bool {
		return !strings.ContainsRune(scrubForDisplay("a\x1b photo"+string(r)+".png b"), r)
	}
	judgeFires := func(r rune) bool { return needsDisplayWash("photo" + string(r) + ".png") }
	var disagreements []string
	record := func(kind string, r rune) {
		if judgeFires(r) != writerEscapes(r) {
			disagreements = append(disagreements,
				fmt.Sprintf("%s %s: judge=%v writer=%v", kind, name327r2(r), judgeFires(r), writerEscapes(r)))
		}
	}
	cfTotal, ccTotal := 0, 0
	for r := rune(0); r <= utf8.MaxRune; r++ { // FULL-family arms: every Cf, every Cc
		if r >= 0xD800 && r <= 0xDFFF {
			continue
		}
		isCf, isCc := unicode.Is(unicode.Cf, r), unicode.IsControl(r)
		if !isCf && !isCc {
			continue
		}
		if isCf {
			cfTotal++
		}
		if isCc {
			ccTotal++
		}
		if isCf {
			record("full-family Cf", r)
		} else {
			record("full-family Cc", r)
		}
	}
	// SAMPLE arm: named boundary scalars (C0/C1 edges, family edges, ordinary runes).
	boundarySample := []rune{
		0x0000, 0x001F, 0x007F, 0x0080, 0x009F, 0x00A0, 0x00AD, 0x061C, 0x180E,
		0x200B, 0x202A, 0x202E, 0x2060, 0x2061, 0x2066, 0x2069, 0x206A, 0x206B,
		0xFEFF, 'a', 0x4E2D, 0x1F600,
	}
	for _, r := range boundarySample {
		record("boundary-sample", r)
	}
	t.Logf("nail 3 scope: full Cf=%d, full Cc=%d, boundary sample=%d, disagreements=%d",
		cfTotal, ccTotal, len(boundarySample), len(disagreements))
	for i, d := range disagreements {
		if i >= 20 {
			t.Errorf("... and %d more judge/writer disagreements", len(disagreements)-20)
			break
		}
		t.Errorf("judge and writer are not one set: %s", d)
	}
	// Positive control: a known invisible rune surfaces as an escape (writer arm alive).
	if out := scrubForDisplay("a\x1bp\u202eqb"); !strings.Contains(out, `\u202e`) {
		t.Errorf("positive control failed: U+202E did not surface as a visible escape in %q", out)
	}
}
