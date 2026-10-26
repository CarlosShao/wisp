package cer

import (
	"errors"
	"fmt"
	"testing"
)

// The calculator is measured against an ALU table rather than against a memory
// of what Levenshtein "should" do: every cell of a small exhaustive grid is
// checked against a second, independently written answer (a top-down oracle
// that shares no table, no rolling rows and no edit-script bookkeeping), and
// against the invariants every edit distance must satisfy. A single bad cell
// fails the run and prints the grid.

// oracleCost is the same problem solved recursively, top-down, with no table
// reuse and no counts: it can only agree with distance() if the bottom-up
// matrix is initialised, indexed and rolled correctly.
func oracleCost(ref, hyp []rune) int {
	memo := make(map[[2]int]int, (len(ref)+1)*(len(hyp)+1))
	var rec func(i, j int) int
	rec = func(i, j int) int {
		if i == len(ref) {
			return len(hyp) - j
		}
		if j == len(hyp) {
			return len(ref) - i
		}
		key := [2]int{i, j}
		if v, ok := memo[key]; ok {
			return v
		}
		best := rec(i+1, j+1)
		if ref[i] != hyp[j] {
			best++
		}
		if d := rec(i+1, j) + 1; d < best {
			best = d
		}
		if d := rec(i, j+1) + 1; d < best {
			best = d
		}
		memo[key] = best
		return best
	}
	return rec(0, 0)
}

// strings over {甲,乙} up to maxLen runes, empty string included.
func gridStrings(maxLen int) []string {
	var out []string
	var grow func(prefix string)
	grow = func(prefix string) {
		out = append(out, prefix)
		if len([]rune(prefix)) == maxLen {
			return
		}
		grow(prefix + "甲")
		grow(prefix + "乙")
	}
	grow("")
	return out
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// TestEditScriptEachOperation is the four named cells of the ALU table: one row
// per operation, each with its exact S/D/I split, and one row proving the split
// is not just a total.
func TestEditScriptEachOperation(t *testing.T) {
	for _, tc := range []struct {
		name string
		ref  string
		hyp  string
		want Counts
	}{
		{"equal", "今天天气", "今天天气", Counts{}},
		{"substitution", "今天天气", "今天天汽", Counts{Substitutions: 1}},
		{"deletion", "今天天气", "今天气", Counts{Deletions: 1}},
		{"insertion", "今天气", "今天天气", Counts{Insertions: 1}},
		{"one-substitution", "人工智能", "人工智慧", Counts{Substitutions: 1}},
		{"two-substitutions", "甲乙丙", "丁戊丙", Counts{Substitutions: 2}},
		{"insert-tail", "今天天气", "今天天气啊", Counts{Insertions: 1}},
	} {
		got := Distance(tc.ref, tc.hyp)
		if got != tc.want {
			t.Fatalf("%s: Distance(%q,%q)=%s, want %s", tc.name, tc.ref, tc.hyp, got, tc.want)
		}
		// The split must be a decomposition of the distance, not a second
		// estimate of it.
		if got.Errors() != oracleCost([]rune(tc.ref), []rune(tc.hyp)) {
			t.Fatalf("%s: %s sums to %d, oracle says %d",
				tc.name, got, got.Errors(), oracleCost([]rune(tc.ref), []rune(tc.hyp)))
		}
	}
}

// TestHandWorkedExample is the example a reader can check by hand: four
// reference characters, one of them deleted, so one error in four references is
// 25%. It also pins the report line CI would print.
func TestHandWorkedExample(t *testing.T) {
	res, err := Score("今天天气", "今天气")
	if err != nil {
		t.Fatalf("Score: %v", err)
	}
	if res.RefRunes != 4 {
		t.Fatalf("RefRunes=%d, want 4", res.RefRunes)
	}
	if want := (Counts{Deletions: 1}); res.Counts != want {
		t.Fatalf("Counts=%s, want %s", res.Counts, want)
	}
	if res.CER != 0.25 {
		t.Fatalf("CER=%v, want exactly 0.25", res.CER)
	}
	if got, want := res.String(), "CER=25.00% (1/4) S=0 D=1 I=0"; got != want {
		t.Fatalf("report line %q, want %q", got, want)
	}
}

// TestALUGrid walks every pair of strings up to three runes over a two-symbol
// alphabet - 15 references times 15 hypotheses - and checks four things per
// cell: the cost matches the oracle, the split decomposes it, the distance
// bounds hold, and the CER identity holds. It then checks symmetry and the
// triangle inequality across the whole grid.
func TestALUGrid(t *testing.T) {
	strs := gridStrings(3)
	if len(strs) != 15 {
		t.Fatalf("grid has %d strings, want 15", len(strs))
	}
	cost := make([][]int, len(strs))
	for i := range cost {
		cost[i] = make([]int, len(strs))
	}
	cells := 0
	for i, ref := range strs {
		for j, hyp := range strs {
			got := Distance(ref, hyp)
			want := oracleCost([]rune(ref), []rune(hyp))
			cost[i][j] = got.Errors()
			cells++

			if got.Errors() != want {
				t.Fatalf("cell (%q,%q): distance %d, oracle %d", ref, hyp, got.Errors(), want)
			}
			n, m := len([]rune(ref)), len([]rune(hyp))
			if got.Errors() < absInt(n-m) || got.Errors() > max(n, m) {
				t.Fatalf("cell (%q,%q): %d edits is outside [|n-m|,max(n,m)]=[%d,%d]",
					ref, hyp, got.Errors(), absInt(n-m), max(n, m))
			}
			if n == 0 {
				if got != (Counts{Insertions: m}) {
					t.Fatalf("cell (%q,%q): empty reference must be %d insertions, got %s", ref, hyp, m, got)
				}
				if _, err := Score(ref, hyp); !errors.Is(err, ErrEmptyReference) {
					t.Fatalf("cell (%q,%q): empty reference scored, want ErrEmptyReference, got %v", ref, hyp, err)
				}
				continue
			}
			res, err := Score(ref, hyp)
			if err != nil {
				t.Fatalf("cell (%q,%q): Score: %v", ref, hyp, err)
			}
			if res.Ref != ref || res.Hyp != hyp {
				t.Fatalf("cell (%q,%q): normalisation is not identity on plain CJK (%q,%q)",
					ref, hyp, res.Ref, res.Hyp)
			}
			if wantCER := float64(got.Errors()) / float64(n); res.CER != wantCER {
				t.Fatalf("cell (%q,%q): CER=%v, want %d/%d=%v", ref, hyp, res.CER, got.Errors(), n, wantCER)
			}
		}
	}
	for i := range strs {
		for j := range strs {
			if cost[i][j] != cost[j][i] {
				t.Fatalf("asymmetric: d(%q,%q)=%d but d(%q,%q)=%d",
					strs[i], strs[j], cost[i][j], strs[j], strs[i], cost[j][i])
			}
			for k := range strs {
				if cost[i][j] > cost[i][k]+cost[k][j] {
					t.Fatalf("triangle inequality broken: d(%q,%q)=%d > d(%q,%q)+d(%q,%q)=%d",
						strs[i], strs[j], cost[i][j], strs[i], strs[k], strs[k], strs[j],
						cost[i][k]+cost[k][j])
				}
			}
		}
	}
	t.Logf("ALU grid: %d cells (%d strings, 3 runes, 2 symbols); symmetry and triangle inequality hold",
		cells, len(strs))
	// A compact excerpt for the evidence table: the cost matrix on the eight
	// shortest strings, so a reader can recompute any cell by hand.
	short := strs[:8]
	var b []byte
	b = fmt.Appendf(b, "cost matrix (rows=ref, cols=hyp) over %v:\n", short)
	b = fmt.Appendf(b, "        ")
	for _, s := range short {
		b = fmt.Appendf(b, "%4q", s)
	}
	for i, s := range short {
		b = fmt.Appendf(b, "\n%6q", s)
		for j := range short {
			b = fmt.Appendf(b, "%4d", cost[i][j])
		}
	}
	t.Logf("%s", b)
}

// TestPositiveControlReferenceEqualsHypothesis: a perfect transcript is exactly
// zero, no rounding, no off-by-one character. The companion cell at the end is
// the guard that keeps this control honest: one character apart must not be
// zero, or "0%" would prove nothing.
func TestPositiveControlReferenceEqualsHypothesis(t *testing.T) {
	for _, s := range []string{"今天天气", "人工智能", "甲乙丙", "ABC123"} {
		res, err := Score(s, s)
		if err != nil {
			t.Fatalf("Score(%q,%q): %v", s, s, err)
		}
		if res.Counts != (Counts{}) {
			t.Fatalf("Score(%q,%q): %s, want no edits", s, s, res.Counts)
		}
		if res.CER != 0 {
			t.Fatalf("Score(%q,%q): CER=%v, want exactly 0", s, s, res.CER)
		}
	}
	// Equal AFTER normalisation is also zero: the same sentence written with
	// punctuation and full-width digits.
	res, err := Score("今天，天气！", "今天 天气")
	if err != nil {
		t.Fatalf("Score: %v", err)
	}
	if res.CER != 0 || res.Counts.Errors() != 0 {
		t.Fatalf("normalised-equal pair scored %s, want 0", res)
	}
	// One character apart is not zero.
	res, err = Score("今天天气", "今天天汽")
	if err != nil {
		t.Fatalf("Score: %v", err)
	}
	if res.CER == 0 || res.Counts.Errors() != 1 {
		t.Fatalf("one substitution scored %v / %s, want >0 with 1 error", res.CER, res.Counts)
	}
}

// TestGates pins both SPEC-10 section 4 ceilings on both sides of the boundary,
// plus the inclusive edge. The aggregate is built from counts directly so the
// arithmetic is exact.
func TestGates(t *testing.T) {
	for _, tc := range []struct {
		cond       Condition
		refs       int
		errs       int
		wantPass   bool
		wantGate   float64
		wantReport string
	}{
		{NearClean, 10000, 583, true, 0.06, "near_clean CER=5.83% gate<=6.00% PASS"},
		{NearClean, 10000, 600, true, 0.06, "near_clean CER=6.00% gate<=6.00% PASS"},
		{NearClean, 10000, 601, false, 0.06, "near_clean CER=6.01% gate<=6.00% FAIL"},
		{NoisyFar, 10000, 1499, true, 0.15, "noisy_far CER=14.99% gate<=15.00% PASS"},
		{NoisyFar, 10000, 1500, true, 0.15, "noisy_far CER=15.00% gate<=15.00% PASS"},
		{NoisyFar, 10000, 1501, false, 0.15, "noisy_far CER=15.01% gate<=15.00% FAIL"},
	} {
		agg := Aggregate{Utterances: 1, RefRunes: tc.refs, Counts: Counts{Deletions: tc.errs}}
		v, err := agg.Judge(tc.cond)
		if err != nil {
			t.Fatalf("%s %d/%d: Judge: %v", tc.cond, tc.errs, tc.refs, err)
		}
		if v.Pass != tc.wantPass {
			t.Fatalf("%s %d/%d: Pass=%v, want %v (%s)", tc.cond, tc.errs, tc.refs, v.Pass, tc.wantPass, v)
		}
		if v.Gate != tc.wantGate {
			t.Fatalf("%s: Gate=%v, want %v", tc.cond, v.Gate, tc.wantGate)
		}
		if v.Errors != tc.errs || v.RefRunes != tc.refs || v.Utterances != 1 {
			t.Fatalf("%s: verdict lost the numbers it judged: %+v", tc.cond, v)
		}
		if got := v.String(); len(got) < len(tc.wantReport) || got[:len(tc.wantReport)] != tc.wantReport {
			t.Fatalf("%s: report line %q, want prefix %q", tc.cond, got, tc.wantReport)
		}
	}
	// No condition, no gate.
	if _, err := (Aggregate{Utterances: 1, RefRunes: 10}).Judge("dev"); !errors.Is(err, ErrUnknownCondition) {
		t.Fatalf("unknown condition: got %v, want ErrUnknownCondition", err)
	}
	// No denominator, no verdict - not a passing one.
	if v, err := (Aggregate{}).Judge(NearClean); !errors.Is(err, ErrNoUtterances) {
		t.Fatalf("empty aggregate: got %+v / %v, want ErrNoUtterances", v, err)
	}
}

// TestScoreSetAggregatesAcrossUtterances: the gate reads a micro-average, so a
// three-character sentence and a six-character sentence must not count the same.
func TestScoreSetAggregatesAcrossUtterances(t *testing.T) {
	set := Set{Condition: NearClean, Utterances: []Utterance{
		{ID: "near_clean/a.wav", Ref: "今天天气"},   // 4 reference characters
		{ID: "near_clean/b.wav", Ref: "人工智能啊呀"}, // 6 reference characters
	}}
	rep, err := set.Score(map[string]string{
		"near_clean/a.wav": "今天气",  // 1 deletion
		"near_clean/b.wav": "人工智能", // 2 deletions
	})
	if err != nil {
		t.Fatalf("Score: %v", err)
	}
	if len(rep.Rows) != 2 {
		t.Fatalf("report has %d rows, want 2", len(rep.Rows))
	}
	if rep.Rows[0].ID != "near_clean/a.wav" || rep.Rows[1].ID != "near_clean/b.wav" {
		t.Fatalf("report is not in manifest order: %s, %s", rep.Rows[0].ID, rep.Rows[1].ID)
	}
	if rep.Rows[0].Result.CER != 0.25 {
		t.Fatalf("row 0 CER=%v, want 0.25", rep.Rows[0].Result.CER)
	}
	if want := 2.0 / 6.0; rep.Rows[1].Result.CER != want {
		t.Fatalf("row 1 CER=%v, want %v", rep.Rows[1].Result.CER, want)
	}
	if rep.Aggregate.Utterances != 2 || rep.Aggregate.RefRunes != 10 {
		t.Fatalf("aggregate over %d utterances / %d refs, want 2/10",
			rep.Aggregate.Utterances, rep.Aggregate.RefRunes)
	}
	if want := 0.3; rep.Aggregate.CER() != want {
		t.Fatalf("aggregate CER=%v, want 3/10=%v (not the mean of 25%% and 33%%)",
			rep.Aggregate.CER(), want)
	}
	v, err := rep.Verdict()
	if err != nil {
		t.Fatalf("Verdict: %v", err)
	}
	if v.Pass {
		t.Fatalf("30%% passed the 6%% gate: %s", v)
	}
	if want := rep.Aggregate.CER(); v.CER != want {
		t.Fatalf("verdict CER=%v, want the aggregate %v", v.CER, want)
	}
}

// TestScoreSetRejectsIncompleteHypotheses: a missing or extra hypothesis changes
// what was measured, so both are errors rather than omissions.
func TestScoreSetRejectsIncompleteHypotheses(t *testing.T) {
	set := Set{Condition: NoisyFar, Utterances: []Utterance{
		{ID: "noisy_far/a.wav", Ref: "今天天气"},
		{ID: "noisy_far/b.wav", Ref: "人工智能"},
	}}
	if _, err := set.Score(map[string]string{"noisy_far/a.wav": "今天天气"}); !errors.Is(err, ErrHypothesisMissing) {
		t.Fatalf("missing hypothesis: got %v, want ErrHypothesisMissing", err)
	}
	if _, err := set.Score(map[string]string{
		"noisy_far/a.wav": "今天天气",
		"noisy_far/b.wav": "人工智能",
		"noisy_far/c.wav": "多的",
	}); !errors.Is(err, ErrHypothesisUnknown) {
		t.Fatalf("unknown hypothesis: got %v, want ErrHypothesisUnknown", err)
	}
	// An utterance that normalises to nothing cannot be scored.
	if _, err := (Set{Condition: NearClean, Utterances: []Utterance{
		{ID: "near_clean/x.wav", Ref: "，。"},
	}}).Score(map[string]string{"near_clean/x.wav": ""}); !errors.Is(err, ErrEmptyReference) {
		t.Fatalf("punctuation-only reference: got %v, want ErrEmptyReference", err)
	}
	// An empty set cannot be scored at all - it has no denominator.
	if _, err := (Set{Condition: NearClean}).Score(map[string]string{}); !errors.Is(err, ErrNoUtterances) {
		t.Fatalf("empty set: got %v, want ErrNoUtterances", err)
	}
}

// TestNoPathToAPassingVerdictWithoutADenominator is the anti-green assertion at
// the API level: every road to a Verdict without reference characters is an
// error, so no caller can obtain Pass=true from nothing.
func TestNoPathToAPassingVerdictWithoutADenominator(t *testing.T) {
	reports := []Report{
		{},
		{Condition: NearClean},
		{Condition: NoisyFar, Rows: []Row{}},
	}
	for i, rep := range reports {
		v, err := rep.Verdict()
		if err == nil {
			t.Fatalf("report %d produced a verdict %s without a denominator", i, v)
		}
		if !errors.Is(err, ErrNoUtterances) {
			t.Fatalf("report %d: got %v, want ErrNoUtterances", i, err)
		}
		if v.Pass {
			t.Fatalf("report %d: a failed verdict is still marked Pass", i)
		}
	}
}
