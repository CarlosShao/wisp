package cer

import (
	"errors"
	"fmt"
	"maps"
	"slices"
)

// Condition is one of the two committed baseline sets. The names are frozen by
// SPEC-01:86-87 (repeated verbatim at SPEC-04:125) as directory names under
// testdata/asr-baseline/. The split is by ACOUSTIC CONDITION, not by dev/test:
// both sets are gates, and neither is a training split.
type Condition string

const (
	// NearClean is close-talk, quiet audio. Gate: CER <= 6%.
	NearClean Condition = "near_clean"
	// NoisyFar is noisy and far-field audio. Gate: CER <= 15%.
	NoisyFar Condition = "noisy_far"
)

// The two gates, as fractions, exactly as SPEC-10 section 4 words them
// (D44): near_clean <= 6%, noisy_far <= 15%.
const (
	NearCleanGate = 0.06
	NoisyFarGate  = 0.15
)

var (
	// ErrUnknownCondition rejects a condition with no gate.
	ErrUnknownCondition = errors.New("cer: unknown acoustic condition")
	// ErrNoUtterances rejects an empty set: no denominator, no verdict. A set
	// that failed to load must never look like a passing gate.
	ErrNoUtterances = errors.New("cer: no utterances to score")
	// ErrHypothesisMissing rejects a set scored without a hypothesis for one of
	// its utterances. Skipping it would shrink the denominator silently.
	ErrHypothesisMissing = errors.New("cer: no hypothesis for utterance")
	// ErrHypothesisUnknown rejects a hypothesis for an utterance the set does
	// not contain: the two halves disagree about what was measured.
	ErrHypothesisUnknown = errors.New("cer: hypothesis for an utterance outside the set")
)

// Gate returns the CER ceiling for the condition.
func (c Condition) Gate() (float64, error) {
	switch c {
	case NearClean:
		return NearCleanGate, nil
	case NoisyFar:
		return NoisyFarGate, nil
	default:
		return 0, fmt.Errorf("%w %q (want %q or %q)", ErrUnknownCondition, string(c), NearClean, NoisyFar)
	}
}

// Utterance is one baseline sentence: ID is the manifest's file path (which is
// also the join key to the ASR hypothesis) and Ref is the reference transcript
// the corpus ships with the audio.
type Utterance struct {
	ID  string
	Ref string
}

// Set is one acoustic condition's utterances, in manifest order. The order is
// part of the contract: a report is diffed between runs, so it must not depend
// on map iteration.
type Set struct {
	Condition  Condition
	Utterances []Utterance
}

// Row is one scored utterance in a report.
type Row struct {
	ID     string
	Result Result
}

// Aggregate is the corpus-level sum. The gate reads this rather than the mean
// of the per-file rates: micro-averaging weights every reference character
// once, which is what a corpus CER means when sentence lengths differ.
type Aggregate struct {
	Utterances int
	RefRunes   int
	Counts     Counts
}

// CER is the aggregate character error rate.
func (a Aggregate) CER() float64 {
	if a.RefRunes == 0 {
		return 0
	}
	return float64(a.Counts.Errors()) / float64(a.RefRunes)
}

// Verdict is the gate decision for one condition.
type Verdict struct {
	Condition  Condition
	Gate       float64
	CER        float64
	Pass       bool
	Utterances int
	RefRunes   int
	Errors     int
}

// String is the CI report line: one line per condition, the numbers the gate
// actually read.
func (v Verdict) String() string {
	word := "FAIL"
	if v.Pass {
		word = "PASS"
	}
	return fmt.Sprintf("%s CER=%.2f%% gate<=%.2f%% %s (%d utterances, %d ref chars, %d errors)",
		v.Condition, v.CER*100, v.Gate*100, word, v.Utterances, v.RefRunes, v.Errors)
}

// Judge applies the condition's gate to an aggregate. It refuses an empty
// aggregate instead of reporting 0%, and the comparison is inclusive (<=), which
// is how SPEC-10 spells "<=6%".
func (a Aggregate) Judge(c Condition) (Verdict, error) {
	if a.Utterances == 0 || a.RefRunes == 0 {
		return Verdict{}, fmt.Errorf("%w: condition %q has %d utterances and %d reference chars",
			ErrNoUtterances, string(c), a.Utterances, a.RefRunes)
	}
	gate, err := c.Gate()
	if err != nil {
		return Verdict{}, err
	}
	v := Verdict{
		Condition:  c,
		Gate:       gate,
		CER:        a.CER(),
		Utterances: a.Utterances,
		RefRunes:   a.RefRunes,
		Errors:     a.Counts.Errors(),
	}
	v.Pass = v.CER <= gate
	return v, nil
}

// Report is one condition's full result: per-file rows plus the aggregate the
// gate reads.
type Report struct {
	Condition Condition
	Rows      []Row
	Aggregate Aggregate
}

// Verdict applies the gate to this report.
func (r Report) Verdict() (Verdict, error) { return r.Aggregate.Judge(r.Condition) }

// Score scores every utterance of the set against the caller's hypotheses,
// keyed by utterance ID. Missing and unknown hypotheses are errors, not
// omissions: both would silently change what was measured.
func (s Set) Score(hyp map[string]string) (Report, error) {
	if len(s.Utterances) == 0 {
		return Report{}, fmt.Errorf("%w: condition %q", ErrNoUtterances, string(s.Condition))
	}
	if _, err := s.Condition.Gate(); err != nil {
		return Report{}, err
	}
	known := make(map[string]bool, len(s.Utterances))
	rep := Report{Condition: s.Condition, Rows: make([]Row, 0, len(s.Utterances))}
	for _, u := range s.Utterances {
		if known[u.ID] {
			return Report{}, fmt.Errorf("cer: duplicate utterance id %q", u.ID)
		}
		known[u.ID] = true
		h, ok := hyp[u.ID]
		if !ok {
			return Report{}, fmt.Errorf("%w: %s (set %q)", ErrHypothesisMissing, u.ID, string(s.Condition))
		}
		res, err := Score(u.Ref, h)
		if err != nil {
			return Report{}, fmt.Errorf("cer: utterance %s: %w", u.ID, err)
		}
		rep.Rows = append(rep.Rows, Row{ID: u.ID, Result: res})
		rep.Aggregate.Utterances++
		rep.Aggregate.RefRunes += res.RefRunes
		rep.Aggregate.Counts = addCounts(rep.Aggregate.Counts, res.Counts)
	}
	for _, id := range slices.Sorted(maps.Keys(hyp)) {
		if !known[id] {
			return Report{}, fmt.Errorf("%w: %s (set %q)", ErrHypothesisUnknown, id, string(s.Condition))
		}
	}
	return rep, nil
}
