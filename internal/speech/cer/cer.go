// Package cer is the CER measurement half of ticket 15 (AC#2): character-level
// edit distance, the frozen normalisation, the sources.tsv reader for the
// committed baseline sets, and the two SPEC-10 section 4 gates (near_clean
// <=6%, noisy_far <=15%).
//
// It is pure Go: no cgo, no onnxruntime, no audio decoding. The audio and the
// reference transcriptions are data this package reads; the hypothesis comes
// from whoever ran the ASR. That split is deliberate - a gate that needs a
// model to compute its own denominator cannot be run in a unit test, and the
// baseline corpus itself is a separate leg (AISHELL-1, see
// docs/evidence/s1/cer-a1-research.md).
//
// Character means rune, not word: Chinese ASR is scored per character, so
// "今天天气" against "今天气" is one deletion out of four references and a CER
// of 25%, not one word error.
package cer

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

// Counts is one edit script's tally. The two halves matter separately: a
// substitution and an insertion inflate CER identically but mean different
// things about the recogniser (mis-heard vs. over-generated), and the report
// carries both.
type Counts struct {
	Substitutions int
	Deletions     int
	Insertions    int
}

// Errors is S + D + I, the numerator of CER.
func (c Counts) Errors() int { return c.Substitutions + c.Deletions + c.Insertions }

// String is the report spelling used by the runner and the evidence tables.
func (c Counts) String() string {
	return fmt.Sprintf("S=%d D=%d I=%d", c.Substitutions, c.Deletions, c.Insertions)
}

// addCounts sums two tallies.
func addCounts(a, b Counts) Counts {
	return Counts{
		Substitutions: a.Substitutions + b.Substitutions,
		Deletions:     a.Deletions + b.Deletions,
		Insertions:    a.Insertions + b.Insertions,
	}
}

// cell is one dynamic-programming state: the cost so far plus how it was
// reached. The split is carried through the table so the backtrace is exact
// instead of reconstructed from a second pass.
type cell struct {
	cost int
	sub  int
	del  int
	ins  int
}

// plus adds one step's cost and label.
func (c cell) plus(step cell) cell {
	return cell{
		cost: c.cost + step.cost,
		sub:  c.sub + step.sub,
		del:  c.del + step.del,
		ins:  c.ins + step.ins,
	}
}

// The three unit-cost steps. CER weights them equally; if a future gate wants a
// different weight it belongs here, not in the callers.
var (
	subStep = cell{cost: 1, sub: 1}
	delStep = cell{cost: 1, del: 1}
	insStep = cell{cost: 1, ins: 1}
)

// prefers reports whether a is the canonical of two cells with the same cost:
// more substitutions first, then fewer deletions, then fewer insertions. Every
// minimal alignment has the same total (that is what "minimal" means), so this
// only fixes which split a report shows; without a rule the split would depend
// on iteration order and no two runs would produce the same evidence table.
func prefers(a, b cell) bool {
	if a.cost != b.cost {
		return a.cost < b.cost
	}
	if a.sub != b.sub {
		return a.sub > b.sub
	}
	if a.del != b.del {
		return a.del < b.del
	}
	return a.ins < b.ins
}

// distance is the weighted Levenshtein distance between two rune slices, with
// the edit script recovered. Two rolling rows are enough: every cell only reads
// the row above and its own left neighbour.
func distance(ref, hyp []rune) Counts {
	if len(ref) == 0 {
		return Counts{Insertions: len(hyp)}
	}
	prev := make([]cell, len(hyp)+1)
	cur := make([]cell, len(hyp)+1)
	for j := range prev {
		prev[j] = cell{cost: j, ins: j}
	}
	for i := 1; i <= len(ref); i++ {
		cur[0] = cell{cost: i, del: i}
		for j := 1; j <= len(hyp); j++ {
			best := prev[j-1]
			if ref[i-1] != hyp[j-1] {
				best = best.plus(subStep)
			}
			if cand := prev[j].plus(delStep); prefers(cand, best) {
				best = cand
			}
			if cand := cur[j-1].plus(insStep); prefers(cand, best) {
				best = cand
			}
			cur[j] = best
		}
		prev, cur = cur, prev
	}
	c := prev[len(hyp)]
	return Counts{Substitutions: c.sub, Deletions: c.del, Insertions: c.ins}
}

// Distance is the raw character-level distance: no normalisation, no CER. It
// exists for two callers: the tests that pin the edit script cell by cell, and
// any diagnostic that needs to know how much of a score is normalisation.
// Scoring uses Score, which normalises first.
func Distance(ref, hyp string) Counts { return distance([]rune(ref), []rune(hyp)) }

// ErrEmptyReference is returned when a reference is empty after normalisation.
// Such a sentence has no denominator, so CER is undefined - reporting 0% would
// turn a broken baseline into a passing gate, which is the one failure mode
// this package must never have.
var ErrEmptyReference = errors.New("cer: reference text is empty after normalization")

// Result is one scored utterance, normalised forms included so a report can
// show what was actually compared.
type Result struct {
	// Ref is the normalised reference; Hyp the normalised hypothesis.
	Ref string
	Hyp string
	// RefRunes is the denominator: len([]rune(Ref)).
	RefRunes int
	// Counts is the edit script, Eq. RefRunes+Insertions == len([]rune(Hyp))+
	// Deletions for every minimal alignment.
	Counts Counts
	// CER is Counts.Errors()/RefRunes as a fraction (0.06 is 6%).
	CER float64
}

// String is the per-utterance report line.
func (r Result) String() string {
	return fmt.Sprintf("CER=%.2f%% (%d/%d) %s", r.CER*100, r.Counts.Errors(), r.RefRunes, r.Counts)
}

// Score normalises both texts and returns the character error rate of hyp with
// respect to ref. An empty normalised reference is an error, never a zero.
func Score(ref, hyp string) (Result, error) {
	nref, nhyp := Normalize(ref), Normalize(hyp)
	n := utf8.RuneCountInString(nref)
	if n == 0 {
		return Result{}, fmt.Errorf("%w (raw reference %q, raw hypothesis %q)",
			ErrEmptyReference, ref, hyp)
	}
	c := distance([]rune(nref), []rune(nhyp))
	return Result{
		Ref:      nref,
		Hyp:      nhyp,
		RefRunes: n,
		Counts:   c,
		CER:      float64(c.Errors()) / float64(n),
	}, nil
}
