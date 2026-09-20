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
		if unicode.IsSpace(r) {
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
type fragmentIndex struct {
	n       int      // effective window size in runes (>=1, contract floor 8)
	norm    []rune   // normalized source content
	winStrs []string // norm window spellings, parallel to hashes (verification)
	hashes  []uint64 // sorted, deduped window hashes
}

// newFragmentIndex indexes norm (already normalized) for >=n-rune substring
// queries. n must already be clamped by the caller. Empty/short inputs yield
// an empty index that can never match (a source shorter than the fragment
// floor has no >=8-char contiguous fragment — documented residual, see
// Provenance docs).
func newFragmentIndex(norm string, n int) *fragmentIndex {
	f := &fragmentIndex{n: n, norm: []rune(norm)}
	if n < 1 {
		n = 1
		f.n = n
	}
	rp := f.norm
	if len(rp) < n {
		f.hashes = nil
		return f
	}
	hset := make([]uint64, 0, len(rp)-n+1)
	ws := make([]string, 0, len(rp)-n+1)
	win := make([]rune, n)
	for i := 0; i+n <= len(rp); i++ {
		copy(win, rp[i:i+n])
		hset = append(hset, hashWindow(win))
		ws = append(ws, string(rp[i:i+n]))
	}
	// Sort hashes with parallel permutation so a hit can be verified against
	// the exact window spelling (collision-safe, keeps match semantics).
	sort.Sort(hashWindowPairs{h: hset, s: ws})
	// Dedupe (same hash implies same-window-length strings verified equal by
	// the Contains check anyway, so dedup on hash keeps the set minimal).
	k := 0
	for i := range hset {
		if i == 0 || hset[i] != hset[k-1] {
			hset[k], ws[k] = hset[i], ws[i]
			k++
		}
	}
	f.hashes = hset[:k]
	f.winStrs = ws[:k]
	return f
}

type hashWindowPairs struct {
	h []uint64
	s []string
}

func (p hashWindowPairs) Len() int { return len(p.h) }
func (p hashWindowPairs) Swap(i, j int) {
	p.h[i], p.h[j] = p.h[j], p.h[i]
	p.s[i], p.s[j] = p.s[j], p.s[i]
}
func (p hashWindowPairs) Less(i, j int) bool { return p.h[i] < p.h[j] }

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
	src := string(f.norm)
	win := make([]rune, n)
	for i := 0; i+n <= len(rp); i++ {
		copy(win, rp[i:i+n])
		h := hashWindow(win)
		if j := sort.Search(len(f.hashes), func(j int) bool { return f.hashes[j] >= h }); j < len(f.hashes) && f.hashes[j] == h {
			w := string(rp[i : i+n])
			// Verify: hash equality is a filter, substring truth decides.
			if strings.Contains(src, w) {
				return w, true
			}
			// Collision on a deduped set: probe neighbours for the real
			// entry (rare; keeps verification exact).
			for k := j - 1; k >= 0 && f.hashes[k] == h; k-- {
				if strings.Contains(src, w) {
					return w, true
				}
			}
			for k := j + 1; k < len(f.hashes) && f.hashes[k] == h; k++ {
				if strings.Contains(src, w) {
					return w, true
				}
			}
		}
	}
	return "", false
}
