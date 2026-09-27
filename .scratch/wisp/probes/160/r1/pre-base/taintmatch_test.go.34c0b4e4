package risk

import (
	"sort"
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

func TestNormalizeTaintDropsIgnorableChars(t *testing.T) {
	// M-6: invisible characters are not whitespace, so they used to cut a
	// fragment's contiguity — a normalization gap in the *contract* sense, not
	// the 16.9#1 paraphrase residual. Spellings are built from code points so
	// nothing here depends on an invisible literal surviving an editor round
	// trip.
	const (
		zwsp = string(rune(0x200B)) // zero-width space
		zwnj = string(rune(0x200C)) // zero-width non-joiner
		zwj  = string(rune(0x200D)) // zero-width joiner
		shy  = string(rune(0x00AD)) // soft hyphen
		wj   = string(rune(0x2060)) // word joiner
		bom  = string(rune(0xFEFF)) // byte-order mark (a literal is illegal in Go source)
		mn   = string(rune(0x0301)) // combining acute (category Mn)
	)
	cases := []struct{ name, in, want string }{
		{"zwsp", "ab" + zwsp + "cd", "abcd"},
		{"zwnj+zwj", "ab" + zwnj + zwj + "cd", "abcd"},
		{"soft hyphen", "ab" + shy + "cd", "abcd"},
		{"word joiner", "ab" + wj + "cd", "abcd"},
		{"bom", "ab" + bom + "cd", "abcd"},
		{"combining marks", "ab" + mn + "c" + mn + "d", "abcd"},
		{"full-width + ignorable mix", "ＡＢ" + shy + "cd", "abcd"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := normalizeTaint(c.in); got != c.want {
				t.Errorf("normalizeTaint(%q)=%q want %q", c.in, got, c.want)
			}
		})
	}
	// And the evasion window is closed end to end in normalized space.
	idx := newFragmentIndex(normalizeTaint("token SUPERSECRETVALUE here"), contractMinFragmentChars)
	for _, probe := range []string{
		"SU" + zwsp + "PERSECRET",
		"SUPER" + shy + "SECRETVALUE",
		"sup" + zwsp + "er" + mn + "se" + wj + "cre" + bom + "t",
	} {
		if _, ok := idx.contains(normalizeTaint(probe)); !ok {
			t.Errorf("ESCAPIABLE: invisible-character split %q did not match", probe)
		}
	}
}

func TestFragmentIndexShortSourceNeverMatches(t *testing.T) {
	// Documented residual: a source shorter than the 8-char floor cannot
	// produce a contract-shaped fragment (e.g. a 7-char secret).
	src := normalizeTaint("abc123")
	idx := newFragmentIndex(src, contractMinFragmentChars)
	if len(idx.hashes) != 0 {
		t.Fatalf("expected empty hash set, got %d", len(idx.hashes))
	}
	// A parameter that shares more than the source has must still not hit.
	if _, ok := idx.contains("abc123extra-not-real"); ok {
		t.Fatal("index with zero windows must never hit")
	}
	// Guard against the degenerate assertion the first version of this test
	// was: the empty index has to be exercised against a populated one too.
	full := newFragmentIndex(normalizeTaint("abc123abcdef tail"), contractMinFragmentChars)
	if len(full.hashes) == 0 {
		t.Fatal("control index must be populated")
	}
	if _, ok := full.contains(normalizeTaint("abc123ab")); !ok {
		t.Fatal("control: an >=8 shared fragment must hit")
	}
}

// N-3: this used to probe a disjoint window (a tautology). It now injects a
// genuine hash hit for a window the source does NOT contain, i.e. a simulated
// 64-bit collision, which is the only case the verification step can fail.
func TestFragmentHashCollisionCannotFakeHit(t *testing.T) {
	src := normalizeTaint(strings.Repeat("x", 64) + "ONLYREALFRAGMENT")
	idx := newFragmentIndex(src, contractMinFragmentChars)

	ghost := normalizeTaint("QQQQQQQQ") // 8 runes, absent from src
	if strings.Contains(src, ghost) {
		t.Fatal("test premise broken: ghost window is in the source")
	}
	// Forge the collision: add the ghost's hash where a real 64-bit collision
	// would have landed.
	idx.hashes = append(idx.hashes, hashWindow([]rune(ghost)))
	sort.Slice(idx.hashes, func(i, j int) bool { return idx.hashes[i] < idx.hashes[j] })

	if _, ok := idx.contains(ghost); ok {
		t.Fatal("hash equality alone must never create a hit")
	}
	// The real fragment still hits, so the forged entry did not break lookup.
	if _, ok := idx.contains(normalizeTaint("ONLYREALFRAG")); !ok {
		t.Fatal("real fragment must still match")
	}
}

// --- M-5 memory/latency budget (go test -bench . -benchmem) ------------------

// BenchmarkMarkFullSource measures one Mark() of a MaxSourceRunes-sized
// tainted source: with the dead winStrs table removed the index should cost
// ~8 bytes per window (~2 MiB for 262144 runes), not ~12 MiB.
func BenchmarkMarkFullSource(b *testing.B) {
	src := strings.Repeat("敏感-token-内容 ", 20000) // > MaxSourceRunes once normalized
	p := NewProvenance(ProvOptions{NoProbe: true})
	p.OpenScope("s")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Mark("s", SrcFSRead, "/big", src)
		p.CloseScope("s")
		p.OpenScope("s")
	}
}

// BenchmarkScanOverTenSources measures the response-path cost of one Inspect
// with a 40k-rune parameter against ten 8k-rune tainted sources.
func BenchmarkScanOverTenSources(b *testing.B) {
	p := NewProvenance(ProvOptions{NoProbe: true})
	p.OpenScope("s")
	for i := 0; i < 10; i++ {
		p.Mark("s", SrcWebFetch, "u", strings.Repeat("filler content to index ", 500)+
			string(rune('a'+i))+strings.Repeat(" more", 40))
	}
	param := strings.Repeat("unrelated body text ", 2000) // ~40k runes
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := p.Inspect("s", "notify", map[string]any{"text": param}); ok {
			b.Fatal("clean param must not hit")
		}
	}
}
