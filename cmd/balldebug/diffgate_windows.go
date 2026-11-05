//go:build windows

package main

// Ticket 71 AC#1/AC#2: the differential run gets a threshold gate.
//
// The measured table used to be pure print output: `-diff` printed
// px_delta_ge8, timers and cpu row by row and returned nil regardless of the
// numbers, so the ticket-62 run whose ball did not image at all (px_delta_ge8=0
// against a 12px Sleeping dot) exited exactly like an ideal run. A gate that
// never consulted its own evidence is the same shape as a scanner that scans
// zero files (ticket 67), so this one judges the table and says what it looked
// at even when it passes.
//
// Two entry points share one evaluator:
//
//	-diff -diff-gate               judge the table the run just wrote
//	-diff-gate-table <path>        judge an already recorded table (no desktop)
//
// The second one exists on purpose: a real -diff needs the interactive desktop
// to shoot alive/dead frames, so the gate would otherwise be untestable in CI
// and the positive control could only be a person reading output.

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// diffGateFloorPx8 is the smallest px_delta_ge8 a row may record and still
// count as "the ball imaged". It is a floor, not a quality bar, and it is
// deliberately far below every accepted measurement:
//
//   - SPEC-08 §2 fixes the Sleeping rest body at SizePx*SleepRestRatio and
//     documents the readable-count geometry: the widest thing the frame draws is
//     drawGlass's 1.5*R halo, so the >=8/255 count is that disc's area. At the
//     default 56px configuration the body is 34.72px and the prediction is
//     pi*(1.5*34.72/2)^2 = 2130 px; three independent runs recorded 2098, 2103
//     and 2120 (docs/evidence/s1/62-diff-glass/, 62-diff-signoff/, docs/SLO.md A.2).
//   - The same spec sets a floor on the rest body: SleepingRestMinPx = 30, whose
//     halo would still put pi*(1.5*30/2)^2 = 1590 px over the >=8/255 line.
//   - Every accepted row recorded in this repo is >= 1577 px (the weakest is a
//     docked Sleeping row at 1620 and a small-window run at 1577). The rejected
//     frame that started all this recorded exactly 0, and a stray cursor or
//     tooltip moves single-digit-to-tens of pixels.
//
// 800 is half of the smallest documented body's halo area (1590/2 = 795): any
// run that images at the smallest legal rest body clears it with a 2x margin,
// while "nothing imaged" cannot. AC#1 forbids lowering it to make a run green,
// so resolveGateMinPx8 refuses a smaller value.
const diffGateFloorPx8 = 800

// diffGateSource travels with every verdict. The number must never be the only
// thing a reader has: "where did 800 come from" is exactly the question that
// kept ticket 62's 44px error alive.
const diffGateSource = "SPEC-08 §2 (rest body 56*SleepRestRatio=34.72px, 1.5R halo => 2130 px predicted, 2098/2103/2120 recorded; SleepingRestMinPx=30 => 1590 px floor)"

// diffGateSpec is what a gate run demands of a table.
type diffGateSpec struct {
	minPx8 int      // floor for every row's px_delta_ge8 (0 = diffGateFloorPx8)
	want   []string // labels that must appear exactly once (nil = judge what is there)
}

// resolveGateMinPx8 turns the flag into the floor the gate applies. 0 means
// "the documented floor"; anything below the floor is refused rather than
// clamped, because silently clamping would let a script believe it raised or
// lowered a bar it did not (ticket 71 AC#1: no tuning runs green).
func resolveGateMinPx8(v int) (int, error) {
	if v == 0 {
		return diffGateFloorPx8, nil
	}
	if v < diffGateFloorPx8 {
		return 0, fmt.Errorf("-diff-gate-min-px8=%d is below the documented floor %d, and it may only be raised: %s",
			v, diffGateFloorPx8, diffGateSource)
	}
	return v, nil
}

// gateBaseState strips the -dock-<edge> suffix a variant label carries, so a
// docked row is judged against the state policy, not against a fake state name.
func gateBaseState(label string) string {
	if i := strings.Index(label, "-dock-"); i > 0 {
		return label[:i]
	}
	return label
}

// gateRequiresZeroTimers answers "does the shipped policy grant this state no
// animation timer at all". It reads internal/ball instead of carrying a copy of
// the static-state list: Acting is in the animated whitelist and really did
// record timers=yes (docs/evidence/s1/12-ball-walk), while Warm/Settling arm
// timers through the alpha and fade paths, so a hand-copied list would have
// redded a legitimate recorded run. A state's policy moving is then the only
// way this expectation moves, which is what "single source of truth" means.
func gateRequiresZeroTimers(label string) bool {
	s := statemachine.State(gateBaseState(label))
	return !ball.Animated(s) && ball.AnimationPolicy(s).Kind == ball.AnimNone
}

// diffGateViolations is the pure judgement: every reason this table is not
// evidence, in the order a reader can act on. Empty means the rows satisfy the
// spec - nothing is inferred from rows that are absent unless they were asked
// for by name.
func diffGateViolations(rows []stateRow, spec diffGateSpec) []string {
	minPx8 := spec.minPx8
	if minPx8 == 0 {
		minPx8 = diffGateFloorPx8
	}
	var out []string

	counts := make(map[string]int, len(rows))
	labels := make([]string, 0, len(rows))
	for _, r := range rows {
		if counts[r.state] == 0 {
			labels = append(labels, r.state)
		}
		counts[r.state]++
	}
	for _, label := range spec.want {
		if counts[label] == 0 {
			out = append(out, fmt.Sprintf("%s: requested row missing - this state is unproven, not clean", label))
		}
	}
	for _, label := range labels {
		if counts[label] > 1 {
			out = append(out, fmt.Sprintf("%s: %d rows for one state - a measurement must be unambiguous", label, counts[label]))
		}
	}
	for _, r := range rows {
		if r.changed8 < minPx8 {
			out = append(out, fmt.Sprintf("%s: px_delta_ge8=%d < %d - the ball did not image (the >=8/255 count is too small to be the glass body)",
				r.state, r.changed8, minPx8))
		}
		if gateRequiresZeroTimers(r.state) && r.timers != "no" {
			out = append(out, fmt.Sprintf("%s: timers=%s but SPEC-08 §2 grants this state no animation timer (\"?\" means the status file proved neither)",
				r.state, r.timers))
		}
	}
	return out
}

// gateDiffRows prints what it examined - the count is the point, a gate with no
// workload report is the failure mode this ticket exists for - and returns an
// error when any row violates the spec.
func gateDiffRows(w io.Writer, rows []stateRow, spec diffGateSpec) error {
	minPx8 := spec.minPx8
	if minPx8 == 0 {
		minPx8 = diffGateFloorPx8
	}
	if n := len(spec.want); n > 0 {
		fmt.Fprintf(w, "gate: examined %d rows (of %d requested), min px_delta_ge8=%d\n", len(rows), n, minPx8)
	} else {
		fmt.Fprintf(w, "gate: examined %d rows, min px_delta_ge8=%d\n", len(rows), minPx8)
	}
	fmt.Fprintf(w, "gate: floor source: %s\n", diffGateSource)

	v := diffGateViolations(rows, spec)
	if len(v) == 0 {
		fmt.Fprintf(w, "gate: PASS - every examined row imaged and every zero-timer state held timers=no\n")
		return nil
	}
	for i, s := range v {
		fmt.Fprintf(w, "gate: FAIL %02d: %s\n", i+1, s)
	}
	fmt.Fprintf(w, "gate: %d violation(s) - this run is not evidence\n", len(v))
	return fmt.Errorf("diff gate: %d violation(s) at min px_delta_ge8=%d", len(v), minPx8)
}

// gateDiffTable judges a recorded table. Nothing is shot here, so it works on
// any machine and on the evidence already in the repo - which is how AC#2's
// positive control is a process exit code instead of a claim.
func gateDiffTable(w io.Writer, path string, spec diffGateSpec) error {
	rows, err := readDiffTable(path)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return fmt.Errorf("diff gate: %s has 0 rows - nothing examined is not a pass", path)
	}
	fmt.Fprintf(w, "gate: table=%s\n", path)
	return gateDiffRows(w, rows, spec)
}

// readDiffTable reads back what writeDiffTable wrote. Columns are looked up by
// NAME, not by position: the gate must not silently judge the wrong number if a
// column is ever inserted, and a table missing the column it needs is refused
// rather than treated as clean.
func readDiffTable(path string) ([]stateRow, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read diff table: %w", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	if !sc.Scan() {
		if err := sc.Err(); err != nil {
			return nil, fmt.Errorf("read diff table %s: %w", path, err)
		}
		return nil, nil
	}
	header := strings.Split(strings.TrimRight(sc.Text(), "\r"), "\t")
	col := make(map[string]int, len(header))
	for i, h := range header {
		col[h] = i
	}
	for _, need := range []string{"state", "timers", "px_delta_ge8", "px_total"} {
		if _, ok := col[need]; !ok {
			return nil, fmt.Errorf("read diff table %s: no %q column - refusing to judge a table it cannot read", path, need)
		}
	}

	var rows []stateRow
	for sc.Scan() {
		line := strings.TrimRight(sc.Text(), "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}
		cells := strings.Split(line, "\t")
		need := 0
		for _, name := range []string{"state", "timers", "px_delta_ge8", "px_total"} {
			if col[name] > need {
				need = col[name]
			}
		}
		if len(cells) <= need {
			return nil, fmt.Errorf("read diff table %s: row %q has %d cells, needs %d", path, line, len(cells), need+1)
		}
		row := stateRow{state: cells[col["state"]], timers: cells[col["timers"]]}
		if row.changed8, err = strconv.Atoi(cells[col["px_delta_ge8"]]); err != nil {
			return nil, fmt.Errorf("read diff table %s: px_delta_ge8 %q: %w", path, cells[col["px_delta_ge8"]], err)
		}
		if row.pixels, err = strconv.Atoi(cells[col["px_total"]]); err != nil {
			return nil, fmt.Errorf("read diff table %s: px_total %q: %w", path, cells[col["px_total"]], err)
		}
		rows = append(rows, row)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read diff table %s: %w", path, err)
	}
	return rows, nil
}
