package risk

import (
	"sort"
	"strings"
	"unicode"
)

// C25 fragment matching (SPEC-06 §5, D33/F4, 16.9#1).
//
// The contract judgment is deliberately coarse: a normalized, contiguous
// fragment of >=8 characters of a tainted source appearing in an exfil
// parameter is a hit. Normalization = whitespace stripped, case folded,
// full-width unified with half-width. Exact per-token taint tracking is
// REJECTED (16.9#1) and must not be resurrected here.
//
// Matching runs in normalized space: source content is normalized once at
// Mark time and indexed by all distinct windows of MinFragmentChars runes
// (FNV-1a 64 hashes, sorted + deduped, built eagerly so entries are
// immutable and lock-free to read). A candidate parameter is normalized and
// slid with the same window; a hash hit is always verified with a substring
// search so hash collisions can never create a false hit (they cost time,
// nothing else). >=8 is equivalent to "some 8-rune window of the parameter
// occurs in the source", so the index is exact for the contract predicate.

// contractMinFragmentChars is the SPEC-06 §5 contract floor (>=8 字符). The
// engine clamps any weaker configuration back to this value; stricter values
// (smaller windows) are permitted because 用户只能调严 (SPEC-06 §2) and extra
// confirms are the safe direction (16.9#1).
const contractMinFragmentChars = 8

// isIgnorableTaintRune reports the invisible characters that must be dropped
// by normalization. U+00AD/U+200B..U+200D/U+2060/U+FEFF and the combining
// marks (Mn) render as nothing (or as a diacritic on the previous rune) but
// are NOT unicode.IsSpace, so leaving them in place lets an exfiltrator cut a
// tainted fragment into sub-8-rune pieces with zero-width characters — a gap
// in the *normalized* text the contract matches on, not the semantic
// paraphrase residual of 16.9#1 (adversarial report M-6).
func isIgnorableTaintRune(r rune) bool {
	switch r {
	case 0x00AD, 0x200B, 0x200C, 0x200D, 0x2060, 0xFEFF:
		return true
	}
	return unicode.Is(unicode.Mn, r)
}

// normalizeTaint applies the contract normalization to s and returns the
// normalized form. Runes that fold onto themselves are kept verbatim.
func normalizeTaint(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		// Full-width forms -> half-width (FF01..FF5E shift by 0xFEE0; the
		// ideographic space U+3000 becomes a plain space and is then
		// dropped as whitespace below).
		if r >= 0xFF01 && r <= 0xFF5E {
			r -= 0xFEE0
		}
		if unicode.IsSpace(r) || isIgnorableTaintRune(r) {
			continue
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

// hashWindow is FNV-1a (64-bit) over the runes of w. Small fixed inputs on a
// hot loop: inline rather than encoding/binary or the hash/fnv interfaces.
func hashWindow(w []rune) uint64 {
	h := uint64(14695981039346656037)
	for _, r := range w {
		for i := 0; i < 4; i++ {
			h ^= uint64(r>>(uint(i)*8)) & 0xff
			h *= 1099511628211
		}
	}
	return h
}

// fragmentIndex holds one tainted source in normalized, indexed form.
// Immutable after construction (built in newFragmentIndex), so concurrent
// readers need no per-entry lock.
//
// Storage is the sorted hash set plus the normalized source string and nothing
// else: an earlier revision also kept a write-only []string of window
// spellings parallel to the hashes, which cost ~5x the index memory on a full
// source and was never read (verification is a substring search over src).
// Dropped per adversarial report M-5.
type fragmentIndex struct {
	n      int      // effective window size in runes (>=1, contract floor 8)
	src    string   // normalized source content (the verification corpus)
	nrunes int      // rune count of src (introspection only)
	hashes []uint64 // sorted, deduped window hashes
}

// newFragmentIndex indexes norm (already normalized) for >=n-rune substring
// queries. n must already be clamped by the caller. Empty/short inputs yield
// an empty index that can never match (a source shorter than the fragment
// floor has no >=8-char contiguous fragment — documented residual, see
// Provenance docs).
//
// The optional skip spans (ticket 177 shape A, Q-61甲/A349: the PRECISE
// exemption of host-minted paths only) are half-open rune ranges [lo,hi) in
// the NORMALIZED coordinate space of norm: a window that OVERLAPS a span is
// never hashed. Nothing else about the index changes — src stays the full
// normalized text and contains() keeps verifying every hash hit through
// strings.Contains, so no window the contract would call a fragment can be
// invented by the exclusion, and every window clear of the spans is indexed
// exactly as before. The spans are inputs to this one build: they are not
// stored on the Provenance or on any scope (a per-scope path roster is
// explicitly NOT approved; 177-c1 §2), so an exemption can only ever die with
// the mark it was declared for. Spans must be located AFTER normalization
// (see MarkWithHostPath) — a raw-byte offset into content is not a coordinate
// in norm.
func newFragmentIndex(norm string, n int, skip ...[2]int) *fragmentIndex {
	return newFragmentIndexSkippingValues(norm, n, nil, skip...)
}

// newFragmentIndexSkippingValues is newFragmentIndex plus a VALUE-keyed
// exemption (ticket 183 family F2): exempt holds the window spellings of one
// declared host-minted string, and a window whose spelling is in that set stays
// out of the index at EVERY position of THIS build. MUTATION 183-r1 F2,
// measurement only.
func newFragmentIndexSkippingValues(norm string, n int, exempt map[string]bool, skip ...[2]int) *fragmentIndex {
	rp := []rune(norm)
	f := &fragmentIndex{n: n, src: norm, nrunes: len(rp)}
	if n < 1 {
		f.n, n = 1, 1
	}
	if len(rp) < n {
		f.hashes = nil
		return f
	}
	hashes := make([]uint64, 0, len(rp)-n+1)
	win := make([]rune, n)
	for i := 0; i+n <= len(rp); i++ {
		if windowExcluded(i, n, skip) {
			continue // spans: windows overlapping a declared span stay out
		}
		copy(win, rp[i:i+n])
		if exempt != nil && exempt[string(win)] {
			continue // values: windows spelling the declared string stay out too
		}
		hashes = append(hashes, hashWindow(win))
	}
	sort.Slice(hashes, func(i, j int) bool { return hashes[i] < hashes[j] })
	// Dedupe: equal hashes come from equal windows here, and even a genuine
	// 64-bit collision is harmless because contains() re-derives the candidate
	// window and verifies it against src (hash equality only ever costs a
	// wasted Contains call, never a false hit).
	k := 0
	for i := range hashes {
		if i == 0 || hashes[i] != hashes[k-1] {
			hashes[k] = hashes[i]
			k++
		}
	}
	f.hashes = hashes[:k]
	return f
}

// windowExcluded reports whether the window [i, i+n) overlaps any skip span.
// Overlap (not containment) is deliberate: a window that straddles a span
// edge carries part of the exempted text, and shape A exempts a HOST-WRITTEN
// span, not a best-guess cut. Dropping the straddling windows only shrinks the
// indexed set toward the declared span (the safe direction); the rest of the
// mark keeps every window clear of it (ticket 177 W-1).
func windowExcluded(i, n int, skip [][2]int) bool {
	for _, s := range skip {
		if i < s[1] && s[0] < i+n {
			return true
		}
	}
	return false
}

// runeIndexOf returns the start of needle in hay, or -1. Both are rune
// slices in normalized space — locating a declared span with this and NOT
// with strings.Index over the raw content is the ticket 177 hard ordering
// (i): normalizeTaint drops whitespace and zero-width runes, so a raw offset
// is the wrong coordinate for the index built over norm.
func runeIndexOf(hay, needle []rune) int {
	if len(needle) == 0 || len(needle) > len(hay) {
		return -1
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
			return i
		}
	}
	return -1
}

// windowSpellings enumerates every n-rune window of norm as a set, i.e. the
// exact fragment spellings the declared string itself would produce. MUTATION
// 183-r1 F2, measurement only.
func windowSpellings(norm string, n int) map[string]bool {
	rp := []rune(norm)
	if n < 1 || len(rp) < n {
		return nil
	}
	out := make(map[string]bool, len(rp)-n+1)
	for i := 0; i+n <= len(rp); i++ {
		out[string(rp[i:i+n])] = true
	}
	return out
}

// contains reports whether any window of paramNorm (a normalized candidate)
// occurs in the indexed source, returning the matched fragment. paramNorm
// must be a prefix of the fully normalized candidate already windowed by the
// caller's scan budget; the scan stops after MaxScanRunes runes (see
// Provenance) — the caller enforces the budget, this scans what it is given.
func (f *fragmentIndex) contains(paramNorm string) (string, bool) {
	rp := []rune(paramNorm)
	n := f.n
	if len(rp) < n || len(f.hashes) == 0 {
		return "", false
	}
	win := make([]rune, n)
	for i := 0; i+n <= len(rp); i++ {
		copy(win, rp[i:i+n])
		h := hashWindow(win)
		if j := sort.Search(len(f.hashes), func(j int) bool { return f.hashes[j] >= h }); j < len(f.hashes) && f.hashes[j] == h {
			w := string(rp[i : i+n])
			// Verify: hash equality is a filter, substring truth decides.
			if strings.Contains(f.src, w) {
				return w, true
			}
		}
	}
	return "", false
}
