package panel

// Ticket 333 AC#4 (and AC#8's escape-not-delete), leg 333-r1. 327-r2's nail 3 checks
// that judge and writer AGREE over both families; agreement alone is blind to both arms
// going quiet at once, which is what a fast path that always answers "clean" plus a
// writer that therefore never runs looks like. This file adds the DIRECTION, over the
// FULL Cf family and the FULL Cc family rune by rune (no sampling dressed up as a
// family) plus 22 named boundaries, and it re-reads the widened set for the two
// properties the widening puts under load: nothing is deleted, and washing twice is a
// no-op. Named-boundary list = the 22 scalars 327-r2's nail 3 already samples.

import (
	"fmt"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

var boundary333r1 = []rune{
	0x0000, 0x001F, 0x007F, 0x0080, 0x009F, 0x00A0, 0x00AD, 0x061C, 0x180E,
	0x200B, 0x202A, 0x202E, 0x2060, 0x2061, 0x2066, 0x2069, 0x206A, 0x206B,
	0xFEFF, 'a', 0x4E2D, 0x1F600,
}

func judge333r1(r rune) bool { return needsDisplayWash("photo" + string(r) + ".png") }

func writer333r1(r rune) bool {
	return !strings.ContainsRune(scrubForDisplay("photo"+string(r)+"报告.png"), r)
}

func report333r1(t *testing.T, failures []string) {
	t.Helper()
	for i, f := range failures {
		if i >= 20 {
			t.Errorf("... and %d more wash failures", len(failures)-20)
			break
		}
		t.Error(f)
	}
}

// AC#4: fast path and writer are one set AND both point at "wash", over whole families.
func TestDisplayWashFastPathAndWriterBothAnswerWashAcrossBothFamilies(t *testing.T) {
	var failures []string
	cf, cc := 0, 0
	for r := rune(0); r <= utf8.MaxRune; r++ {
		if r >= 0xD800 && r <= 0xDFFF {
			continue
		}
		isCf, isCc := unicode.Is(unicode.Cf, r), unicode.IsControl(r)
		if !isCf && !isCc {
			continue
		}
		kind := "Cc"
		if isCf {
			kind, cf = "Cf", cf+1
		} else {
			cc++
		}
		j, w := judge333r1(r), writer333r1(r)
		if !j { // arm 1: the fast path must promise the wash for every family rune
			failures = append(failures,
				fmt.Sprintf("%s %s: needsDisplayWash stayed silent, so scrubForDisplay never runs", kind, name327r2(r)))
		}
		if !w { // arm 2, a different line: the writer must escape what was promised
			failures = append(failures,
				fmt.Sprintf("%s %s: scrubForDisplay left the rune visible", kind, name327r2(r)))
		}
		if j != w { // arm 3: and the two arms must not diverge (ticket 327 AC#2's one set)
			failures = append(failures,
				fmt.Sprintf("%s %s: judge=%v writer=%v, the two sides are not one set", kind, name327r2(r), j, w))
		}
	}
	for _, r := range boundary333r1 { // SAMPLE arm, named, at the family edges
		want := unicode.Is(unicode.Cf, r) || unicode.IsControl(r)
		if j, w := judge333r1(r), writer333r1(r); j != want || w != want {
			failures = append(failures, fmt.Sprintf("boundary %s: want wash=%v, judge=%v writer=%v",
				name327r2(r), want, judge333r1(r), writer333r1(r)))
		}
	}
	t.Logf("AC#4 scope: full Cf=%d, full Cc=%d, boundaries=%d, failures=%d",
		cf, cc, len(boundary333r1), len(failures))
	if cf == 0 || cc == 0 { // positive control: an arm that examined nothing is no pass
		t.Errorf("family arms examined nothing (Cf=%d Cc=%d): the loop never ran", cf, cc)
	}
	report333r1(t, failures)
}

// AC#4 with AC#8: over the whole category ticket 333 just widened, the shape is a
// visible escape - never a splice, never a no-op, never a double escape.
func TestDisplayWashEscapesAndNeverDeletesAcrossTheWholeCategory(t *testing.T) {
	var failures []string
	examined := 0
	for r := rune(0); r <= utf8.MaxRune; r++ {
		if r >= 0xD800 && r <= 0xDFFF || !unicode.Is(unicode.Cf, r) {
			continue
		}
		examined++
		in := "photo" + string(r) + "报告.png"
		out := scrubForDisplay(in)
		if out == in {
			failures = append(failures, fmt.Sprintf("%s: %q came back untouched", name327r2(r), in))
			continue
		}
		if strings.ContainsRune(out, r) {
			failures = append(failures, fmt.Sprintf("%s: still visible in %q", name327r2(r), out))
		}
		if want := "photo" + escapeOf(r) + "报告.png"; out != want {
			failures = append(failures, fmt.Sprintf("%s: washed to %q, want %q", name327r2(r), out, want))
		}
		if !strings.HasPrefix(out, "photo") || !strings.HasSuffix(out, "报告.png") {
			failures = append(failures, fmt.Sprintf("%s: the visible halves spliced into %q", name327r2(r), out))
		}
		if again := scrubForDisplay(out); again != out {
			failures = append(failures, fmt.Sprintf("%s: not idempotent, %q -> %q", name327r2(r), out, again))
		}
	}
	product := len(washSetOver(isDisplayWashCf))
	t.Logf("AC#8 category arm: examined=%d, product washes=%d, failures=%d", examined, product, len(failures))
	if examined != product { // ⛔ a sampled family passed off as the whole family
		t.Errorf("this arm ran %d runes but the product washes %d: the family was sampled, not run",
			examined, product)
	}
	report333r1(t, failures)
}
