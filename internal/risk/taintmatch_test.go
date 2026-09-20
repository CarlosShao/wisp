package risk

import (
	"strings"
	"testing"
)

// AC (ticket 19): "Normalization tests: whitespace/case/full-width transforms
// still match; <8-char does not."

func TestNormalizeTaintContractForm(t *testing.T) {
	cases := []struct{ in, want string }{
		{"AbC dEf", "abcdef"},       // whitespace strip + case fold
		{"ＡＢＣ\tｄｅｆ", "abcdef"},      // full-width letters, tab
		{"1234５６７８", "12345678"},    // full-width digits unified
		{"ａｂ\u3000Ｃ", "abc"},        // ideographic space dropped
		{"秘密 秘密\n秘密１２", "秘密秘密秘密12"}, // CJK + full-width digit
		{"NoChange", "nochange"},
		{"  \t\n ", ""},
	}
	for _, c := range cases {
		if got := normalizeTaint(c.in); got != c.want {
			t.Errorf("normalizeTaint(%q)=%q want %q", c.in, got, c.want)
		}
	}
}

func TestFragmentMatchExactBoundary(t *testing.T) {
	src := normalizeTaint("the secret value is SUPERSECRETVALUE")
	idx := newFragmentIndex(src, contractMinFragmentChars)

	// 11-char contiguous fragment: hit.
	if frag, ok := idx.contains(normalizeTaint("SUPERSECRETVALUE")); !ok || frag == "" {
		t.Fatalf("full fragment must match, got %q/%v", frag, ok)
	}
	// Exactly 8 chars of it: hit (contract floor is >=8).
	if _, ok := idx.contains(normalizeTaint("SUPERSEC")); !ok {
		t.Fatal("exactly-8 fragment must match")
	}
	// 7 chars: NO match (AC: <8-char does not).
	if _, ok := idx.contains(normalizeTaint("SUPERSE")); ok {
		t.Fatal("7-char fragment must not match")
	}
}

func TestFragmentMatchAcrossTransforms(t *testing.T) {
	// Tainted source as read (mixed spellings); full-width letters normalize
	// to ASCII and hyphens survive normalization.
	source := normalizeTaint("Token: Ａ-Ｂ-Ｃ-Ｄ-Ｅ-Ｆ-Ｇ-Ｈ-Ｉ-Ｊ")
	idx := newFragmentIndex(source, contractMinFragmentChars)

	// Exfil candidate with different whitespace/case/full-width spelling of
	// the same contiguous fragment (normalized): must hit.
	if _, ok := idx.contains(normalizeTaint("please look up token: a-b-c-d-e-f-g-h-i-j thanks")); !ok {
		t.Fatal("transformed duplicate must match")
	}
	// Only 7 normalized chars of the fragment leak (hyphens count): no hit.
	if _, ok := idx.contains(normalizeTaint("a-b-c-d")); ok {
		t.Fatal("sub-8 leak must not match")
	}
}

func TestFragmentMatchCJK(t *testing.T) {
	src := normalizeTaint("会议 纪要：下 季度 预算 上调 百分之 十五")
	idx := newFragmentIndex(src, contractMinFragmentChars)
	// 8 CJK runes contiguous (whitespace-insensitive): hit.
	if _, ok := idx.contains(normalizeTaint("转发：会议纪要下季度预算上调百分之十五")); !ok {
		t.Fatal("CJK >=8-rune fragment must match")
	}
	if _, ok := idx.contains(normalizeTaint("季度预算上")); ok {
		t.Fatal("CJK 5-rune fragment must not match")
	}
}

func TestFragmentIndexShortSourceNeverMatches(t *testing.T) {
	// Documented residual: a source shorter than the 8-char floor cannot
	// produce a contract-shaped fragment (e.g. a 7-char secret).
	idx := newFragmentIndex(normalizeTaint("abc123"), contractMinFragmentChars)
	if _, ok := idx.contains("abc123extra-not-real"); ok {
		t.Fatal("index with zero windows must never hit")
	}
	if len(idx.hashes) != 0 {
		t.Fatalf("expected empty hash set, got %d", len(idx.hashes))
	}
}

func TestFragmentHashCollisionCannotFakeHit(t *testing.T) {
	// Even if a hash collides, the Contains verification decides. Build an
	// index over one source and probe a disjoint parameter that shares no
	// 8-rune window: must not report a hit.
	src := newFragmentIndex(normalizeTaint(strings.Repeat("x", 64)+"ONLYREALFRAGMENT"), contractMinFragmentChars)
	if _, ok := src.contains(normalizeTaint("yyyyyyyyzzzzzzzz")); ok {
		t.Fatal("disjoint window must not match")
	}
}
