package panel

// Ticket 327 residual R1, in-repo by leg 327-r2: the off-tree bench caught a silently
// added U+206B, the shipped suite did not. APPEND-ONLY. Nails: (1) cardinality - the
// roster read ONLY from isDisplayWashCf over every scalar, compared against the
// ruling's list BOTH ways (extra=rune added, missing=rune removed, different lines;
// the literals are the EXPECTED side, never the ACTUAL side, so no hand-copy compared
// to itself); (2) exclusion - the eight named runes washed=false; (3) same-set -
// needsDisplayWash and scrubForDisplay agree: Cf/Cc run FULL family, the rest SAMPLE.

import (
	"fmt"
	"sort"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

var ruledRoster327r2 = []rune{ // EXPECTED side only (the ruling's list, hand-copied).
	0x00AD, 0x200B, 0x200C, 0x200D, 0x2060, 0xFEFF,
	0x202A, 0x202B, 0x202C, 0x202D, 0x202E, 0x2066, 0x2067, 0x2068, 0x2069,
}

var excludedRoster327r2 = []rune{0x2061, 0x2062, 0x2063, 0x2064, 0x061C, 0x180E, 0x206A, 0x206B} // off-ruling eight; U+206B was the omitted one (R1)

func name327r2(r rune) string { return fmt.Sprintf("U+%04X", r) }

func names327r2(m map[rune]bool) []string {
	out := make([]string, 0, len(m))
	for r := range m {
		out = append(out, name327r2(r))
	}
	sort.Strings(out)
	return out
}

// productWashedRunes is the ACTUAL side: the product predicate over every scalar.
func productWashedRunes() map[rune]bool {
	washed := map[rune]bool{}
	for r := rune(0); r <= utf8.MaxRune; r++ {
		if r >= 0xD800 && r <= 0xDFFF {
			continue
		}
		if isDisplayWashCf(r) {
			washed[r] = true
		}
	}
	return washed
}

// Nail 1: cardinality, read from the product.
func TestDisplayWashRosterCardinalityIsExactlyTheRuledFifteen(t *testing.T) {
	actual := productWashedRunes()
	ruled := map[rune]bool{}
	for _, r := range ruledRoster327r2 {
		ruled[r] = true
	}
	extra, missing := map[rune]bool{}, map[rune]bool{}
	for r := range actual {
		if !ruled[r] {
			extra[r] = true
		}
	}
	for r := range ruled {
		if !actual[r] {
			missing[r] = true
		}
	}
	t.Logf("nail 1: product-washed=%d, ruled-list=%d", len(actual), len(ruled))
	if len(extra) != 0 { // reddens on a silent ADD (shot 1)
		t.Errorf("the product washes Cf runes the ruling never named: %v (ruled=%d product=%d)",
			names327r2(extra), len(ruled), len(actual))
	}
	if len(missing) != 0 { // reddens on a silent REMOVE (shot 2), a different line
		t.Errorf("a ruled rune is no longer washed by the product: %v (ruled=%d product=%d)",
			names327r2(missing), len(ruled), len(actual))
	}
	if len(actual) != len(ruled) {
		t.Errorf("product roster cardinality=%d, want %d", len(actual), len(ruled))
	}
}

// Nail 2: the eight named exclusions, read off the product.
func TestDisplayWashRosterExcludesTheEightNamedRunes(t *testing.T) {
	for _, r := range excludedRoster327r2 {
		if !unicode.Is(unicode.Cf, r) {
			t.Errorf("%s is not unicode.Cf - the exclusion lost its meaning", name327r2(r))
		}
		if isDisplayWashCf(r) {
			t.Errorf("%s is off the ruled list but isDisplayWashCf washed it", name327r2(r))
		}
		in := "photo" + string(r) + ".png"
		if needsDisplayWash(in) {
			t.Errorf("%s is off the ruled list but needsDisplayWash fired on %q", name327r2(r), in)
		}
		if got := scrubForDisplay(in); got != in {
			t.Errorf("%s is off the ruled list but the wash changed it: %q -> %q", name327r2(r), in, got)
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
