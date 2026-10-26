package cer

import (
	"strings"
	"unicode"
)

// maxMarkupRunes bounds how far a markup marker such as <UNK> may extend. A
// real corpus marker is a handful of characters; the bound keeps a stray '<'
// in a hypothesis from swallowing the rest of the sentence.
const maxMarkupRunes = 24

// Normalize applies the fixed normalisation of the CER baseline. The four rules
// are frozen by docs/evidence/s1/cer-a1-research.md section 3 (and restated in
// the ticket 15 dispatch); every CER number either gate emits depends on them,
// so changing one is a contract change, not a refactor:
//
//  1. whitespace and punctuation are dropped: only letters, digits and marks
//     survive;
//  2. full-width ASCII (U+FF01..U+FF5E) is folded to half-width and the
//     ideographic space (U+3000) to a plain space, which rule 1 then drops;
//  3. simplified and traditional Chinese are NOT converted: 汉 and 漢 stay
//     different characters, because converting them would forgive a real
//     recognition error;
//  4. digits and Latin text are kept verbatim: no inverse text normalisation,
//     so "2026" and "二零二六" remain different strings.
//
// One more rule comes from the same note, because the corpora need it:
// angle-bracket markup such as <UNK>, <sil> and <noise> is dropped rather than
// spelled out. Without it the marker would score as several errors no ASR could
// have avoided.
//
// Case is deliberately not folded. Rule 4 says verbatim, and the near_clean
// slice excludes utterances containing Latin or digits (cer-a1 section 4), so
// that choice cannot move either gate.
func Normalize(s string) string {
	rs := []rune(s)
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		if r == '<' {
			if end, ok := markupEnd(rs, i); ok {
				i = end
			}
			// A stray '<' is punctuation, dropped by the loop.
			continue
		}
		if r = foldWidth(r); keepRune(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// markupEnd reports the index of the '>' that closes the marker opening at
// start, if this really is a marker: non-empty, no nested '<', no whitespace
// and within maxMarkupRunes.
func markupEnd(rs []rune, start int) (int, bool) {
	for j := start + 1; j < len(rs) && j <= start+maxMarkupRunes; j++ {
		switch {
		case rs[j] == '>':
			return j, j > start+1
		case rs[j] == '<' || unicode.IsSpace(rs[j]):
			return 0, false
		}
	}
	return 0, false
}

// foldWidth is rule 2. The full-width block U+FF01..U+FF5E is the ASCII block
// plus 0xFEE0, so one subtraction folds digits, Latin letters and the full-width
// forms of the ASCII punctuation that rule 1 is about to drop anyway.
func foldWidth(r rune) rune {
	switch {
	case r >= 0xFF01 && r <= 0xFF5E:
		return r - 0xFEE0
	case r == 0x3000:
		return ' '
	default:
		return r
	}
}

// keepRune is rule 1's other half. Letters include CJK ideographs, hiragana and
// Latin; digits include the ASCII ones the width fold just produced; marks keep
// combining accents attached instead of turning them into errors.
func keepRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsMark(r)
}
