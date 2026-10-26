package cer

import "testing"

// One assertion group per frozen normalisation rule (cer-a1 section 3). The two
// rules that are defined by NOT doing something - no simplified/traditional
// conversion, no inverse text normalisation of digits - are tested by demanding
// that the score be WRONG: a converter in the pipeline would turn these cases
// green, and a green here means the gate forgives a real error.

func TestRuleOneDropsWhitespaceAndPunctuation(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"今天，天气！", "今天天气"},
		{"今天、天气；好吗？", "今天天气好吗"},
		{" hello,\tworld! ", "helloworld"},
		{"今天...天气", "今天天气"},
		{"(今天)【天气】", "今天天气"},
		{"a-b_c", "abc"},
	} {
		if got := Normalize(tc.in); got != tc.want {
			t.Fatalf("Normalize(%q)=%q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestRuleTwoFoldsFullWidth(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"ＡＢＣ１２３", "ABC123"},
		{"２０２６", "2026"},
		{"　今　天　", "今天"}, // U+3000 ideographic space
		{"！＂＃＄％＆", ""},  // full-width punctuation, all dropped by rule 1
	} {
		if got := Normalize(tc.in); got != tc.want {
			t.Fatalf("Normalize(%q)=%q, want %q", tc.in, got, tc.want)
		}
	}
	// The fold is what makes the two spellings score equal.
	res, err := Score("２０２６年，今天天气", "2026年今天天气")
	if err != nil {
		t.Fatalf("Score: %v", err)
	}
	if res.CER != 0 {
		t.Fatalf("full-width reference vs half-width hypothesis scored %s, want 0", res)
	}
}

func TestRuleThreeNeverConvertsSimplifiedAndTraditional(t *testing.T) {
	if got := Normalize("漢語"); got != "漢語" {
		t.Fatalf("Normalize(漢語)=%q: traditional characters must survive verbatim", got)
	}
	if got := Normalize("汉语"); got != "汉语" {
		t.Fatalf("Normalize(汉语)=%q", got)
	}
	// This MUST be wrong. 2 substitutions out of 2 references is 100%; a
	// normaliser that converted 漢->汉 would report 0% and quietly forgive
	// every traditional-script output.
	res, err := Score("汉语", "漢語")
	if err != nil {
		t.Fatalf("Score: %v", err)
	}
	if want := (Counts{Substitutions: 2}); res.Counts != want {
		t.Fatalf("simplified reference vs traditional hypothesis: %s, want %s", res.Counts, want)
	}
	if res.CER != 1.0 {
		t.Fatalf("simplified reference vs traditional hypothesis: CER=%v, want exactly 1.0 (100%%)", res.CER)
	}
	if res.Ref != "汉语" || res.Hyp != "漢語" {
		t.Fatalf("a conversion slipped into normalisation: ref=%q hyp=%q", res.Ref, res.Hyp)
	}
}

func TestRuleFourDigitsAreNotSpelledOut(t *testing.T) {
	if got := Normalize("2026"); got != "2026" {
		t.Fatalf("Normalize(2026)=%q, want verbatim", got)
	}
	// This MUST be wrong too: no inverse text normalisation, so the digit
	// string and its Chinese reading are four different characters.
	res, err := Score("2026", "二零二六")
	if err != nil {
		t.Fatalf("Score: %v", err)
	}
	if want := (Counts{Substitutions: 4}); res.Counts != want {
		t.Fatalf("2026 vs 二零二六: %s, want %s", res.Counts, want)
	}
	if res.CER != 1.0 {
		t.Fatalf("2026 vs 二零二六: CER=%v, want exactly 1.0", res.CER)
	}
	// The reverse direction has its own denominator, and must also be wrong.
	rev, err := Score("二零二六", "2026")
	if err != nil {
		t.Fatalf("Score: %v", err)
	}
	if rev.RefRunes != 4 || rev.CER != 1.0 {
		t.Fatalf("二零二六 vs 2026: %s, want 4 ref chars and 100%%", rev)
	}
}

func TestMarkupMarkersAreDroppedNotSpelledOut(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"今天<UNK>天气", "今天天气"},
		{"<sil>今天<noise>天气", "今天天气"},
		{"今天<天气", "今天天气"},       // a stray '<' is punctuation
		{"今天<>天气", "今天天气"},      // an empty marker is not a marker
		{"今天<UNK天气", "今天UNK天气"}, // an unclosed marker is not a marker: the text stays
	} {
		if got := Normalize(tc.in); got != tc.want {
			t.Fatalf("Normalize(%q)=%q, want %q", tc.in, got, tc.want)
		}
	}
	// Without this rule the marker would cost five errors no ASR could avoid.
	res, err := Score("今天<UNK>天气", "今天天气")
	if err != nil {
		t.Fatalf("Score: %v", err)
	}
	if res.CER != 0 {
		t.Fatalf("markup reference vs clean hypothesis scored %s, want 0", res)
	}
}

// TestNormalisationIsTheOnlyDifference keeps the two layers apart: the raw
// distance counts the punctuation, the score does not. If a future change moves
// normalisation into Distance, this is the test that notices.
func TestNormalisationIsTheOnlyDifference(t *testing.T) {
	if got, want := Distance("今天，天气！", "今天天气"), 2; got.Errors() != want {
		t.Fatalf("raw Distance sees %d edits, want %d (two punctuation characters)", got.Errors(), want)
	}
	res, err := Score("今天，天气！", "今天天气")
	if err != nil {
		t.Fatalf("Score: %v", err)
	}
	if res.CER != 0 {
		t.Fatalf("normalised score=%s, want 0", res)
	}
	// And the identity holds for all four rules at once.
	res, err = Score("２０２６年，今天<UNK>天气！", "2026年今天天气")
	if err != nil {
		t.Fatalf("Score: %v", err)
	}
	if res.CER != 0 || res.RefRunes != 9 {
		t.Fatalf("all rules at once: %s (ref %q, %d chars), want 0 over 9", res, res.Ref, res.RefRunes)
	}
}
