package panel

// v1 verdict-leg rulers. They live ONLY in the out-of-repo export tree (second
// copy, D:/tmp/w327v1/bench-ex); the repository's tracked bytes are never touched.
//
//  1. TestV1CfSweepIsExactlyTheRuledFifteen: every scalar value whose category is
//     Cf, asked of the real scrubForDisplay. The expected set is hard-coded from
//     the ruling's enumeration (A1008 第 5 节第 1 条), NOT read from
//     isDisplayWashCf - so a roster that is too wide OR too narrow is a red,
//     including a widening to a rune that is neither ruled nor named as an
//     exclusion (the shape the delivered suite is blind to; see the m6 replay).
//  2. TestV1NamedExclusionsPassThroughRaw: the same question rune by rune for the
//     eight the census measured as reachable, U+206B included.
//  3. TestV1JudgeAndWriterAgreeOnEveryScalarValue: AC#2 stated as a property over
//     the whole scalar space. The writer is also asked on a string whose fast path
//     is forced open by an unrelated ESC, because the early return otherwise
//     hides writer-blindness (both sides then report "no change") - which is
//     exactly the shape ticket 327 AC#2 is about.

import (
	"fmt"
	"sort"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

// v1Ruled15 is the ruling's list, spelled out here independently of the code.
var v1Ruled15 = []rune{
	0x00AD,
	0x200B, 0x200C, 0x200D, 0x2060, 0xFEFF,
	0x202A, 0x202B, 0x202C, 0x202D, 0x202E,
	0x2066, 0x2067, 0x2068, 0x2069,
}

func v1Set(rs []rune) map[rune]bool {
	m := make(map[rune]bool, len(rs))
	for _, r := range rs {
		m[r] = true
	}
	return m
}

func v1Names(ms map[rune]bool) []string {
	out := make([]string, 0, len(ms))
	for r := range ms {
		out = append(out, fmt.Sprintf("U+%04X", r))
	}
	sort.Strings(out)
	return out
}

// v1Escapes reports whether the writer turns the rune into a visible escape when
// the fast path is known to fire (an ESC elsewhere in the string forces it open).
func v1Escapes(r rune) bool {
	clean := "photo" + string(r) + ".png"
	return !strings.ContainsRune(scrubForDisplay("a\x1b"+clean+"b"), r)
}

func v1EachScalar(fn func(r rune)) {
	for r := rune(0); r <= utf8.MaxRune; r++ {
		if r >= 0xD800 && r <= 0xDFFF {
			continue
		}
		fn(r)
	}
}

func TestV1CfSweepIsExactlyTheRuledFifteen(t *testing.T) {
	ruled, got := v1Set(v1Ruled15), map[rune]bool{}
	extra, missing := map[rune]bool{}, map[rune]bool{}
	totalCf := 0
	v1EachScalar(func(r rune) {
		if !unicode.Is(unicode.Cf, r) {
			return
		}
		totalCf++
		if !v1Escapes(r) {
			return
		}
		got[r] = true
		if !ruled[r] {
			extra[r] = true
		}
	})
	for r := range ruled {
		if !got[r] {
			missing[r] = true
		}
	}
	t.Logf("v1 sweep: unicode.Cf cardinality=%d, boundary-washed Cf cardinality=%d", totalCf, len(got))
	if len(extra) != 0 {
		t.Errorf("the boundary washes Cf runes the ruling did not name: %v", v1Names(extra))
	}
	if len(missing) != 0 {
		t.Errorf("the boundary leaves ruled Cf runes raw: %v", v1Names(missing))
	}
}

func TestV1NamedExclusionsPassThroughRaw(t *testing.T) {
	for _, r := range []rune{0x2061, 0x2062, 0x2063, 0x2064, 0x061C, 0x180E, 0x206A, 0x206B} {
		in := "photo" + string(r) + ".png"
		if !unicode.Is(unicode.Cf, r) {
			t.Errorf("%s is not Cf, the control lost its meaning", fmt.Sprintf("U+%04X", r))
		}
		if needsDisplayWash(in) {
			t.Errorf("%s is off the ruled list but needsDisplayWash fired", fmt.Sprintf("U+%04X", r))
		}
		if got := scrubForDisplay(in); got != in {
			t.Errorf("%s is off the ruled list but the wash changed it: %q -> %q", fmt.Sprintf("U+%04X", r), in, got)
		}
		t.Logf("exclusion control %s: Cf=true washed=%v", fmt.Sprintf("U+%04X", r), scrubForDisplay(in) != in)
	}
}

func TestV1JudgeAndWriterAgreeOnEveryScalarValue(t *testing.T) {
	ruled := v1Set(v1Ruled15)
	type bad struct {
		r   rune
		why string
	}
	var disagreements []bad
	silentWriter := map[rune]bool{}
	total := 0
	v1EachScalar(func(r rune) {
		if r == utf8.RuneError {
			return // the invalid-byte shape is ticket 216's own case
		}
		total++
		want := unicode.IsControl(r) || ruled[r]
		clean := "photo" + string(r) + ".png"
		judge := needsDisplayWash(clean)
		if judge != want {
			disagreements = append(disagreements, bad{r, "needsDisplayWash"})
		}
		if v1Escapes(r) != want {
			disagreements = append(disagreements, bad{r, "scrubForDisplay"})
		}
		if judge && scrubForDisplay(clean) == clean {
			silentWriter[r] = true
		}
	})
	t.Logf("v1 sweep: scalars=%d, rule-disagreements=%d, fastpath-fires-but-write-unchanged=%d",
		total, len(disagreements), len(silentWriter))
	for i, d := range disagreements {
		if i >= 20 {
			t.Errorf("... and %d more rune-shape disagreements", len(disagreements)-20)
			break
		}
		t.Errorf("%s disagrees with the ruled set at %s (Cc=%v Cf=%v)", d.why,
			fmt.Sprintf("U+%04X", d.r), unicode.IsControl(d.r), unicode.Is(unicode.Cf, d.r))
	}
	for r := range silentWriter {
		if unicode.IsControl(r) || ruled[r] {
			t.Errorf("the fast path fires on %s but the clean-string write left it raw: the two rulers are not one set",
				fmt.Sprintf("U+%04X", r))
		}
	}
}
