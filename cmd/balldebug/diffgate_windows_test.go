//go:build windows

package main

// Ticket 71 AC#1/AC#2: the differential evidence run must FAIL when the ball
// did not image, instead of printing 0 and exiting 0. These tests pin the gate
// itself; the process-level positive control (a real recorded non-imaging run
// going red) is in docs/evidence/s1/71-balldebug-diff-gate.md.

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// gateRow builds the only three columns the gate reads.
func gateRow(state, timers string, px8 int) stateRow {
	return stateRow{state: state, timers: timers, changed8: px8}
}

// TestDiffGateRedOnANonImagingRun is AC#2 at the unit seam: the ticket-62
// baseline frame - the rejected 12px Sleeping dot - recorded px_delta_ge8=0
// and the old -diff printed that 0 and returned nil. The gate must refuse it.
func TestDiffGateRedOnANonImagingRun(t *testing.T) {
	rows := []stateRow{gateRow("Sleeping", "no", 0)}
	var out bytes.Buffer
	err := gateDiffRows(&out, rows, diffGateSpec{minPx8: diffGateFloorPx8, want: []string{"Sleeping"}})
	if err == nil {
		t.Fatalf("gate passed a non-imaging run (px_delta_ge8=0); output=%q", out.String())
	}
	if !strings.Contains(out.String(), "px_delta_ge8=0") {
		t.Errorf("red gate did not report the measured 0 it judged; output=%q", out.String())
	}
}

// TestDiffGateGreenOnRecordedGoodRuns pins the other direction: the numbers
// three accepted runs actually recorded (SPEC-08 §2: 2098/2103/2120 for the
// undocked Sleeping orb, 1620 docked, 3220-4856 animated) must all pass, so a
// narrow gate cannot become noise the next agent learns to ignore.
func TestDiffGateGreenOnRecordedGoodRuns(t *testing.T) {
	rows := []stateRow{
		gateRow("Sleeping", "no", 2120),
		gateRow("Sleeping-dock-right", "no", 1620),
		gateRow("Listening", "yes", 4856),
		gateRow("Listening-dock-right", "yes", 3230),
		gateRow("Speaking", "yes", 4848),
	}
	var out bytes.Buffer
	if err := gateDiffRows(&out, rows, diffGateSpec{minPx8: diffGateFloorPx8}); err != nil {
		t.Fatalf("gate rejected recorded good rows: %v; output=%q", err, out.String())
	}
	if !strings.Contains(out.String(), "5 rows") {
		t.Errorf("green gate did not self-report how many rows it examined; output=%q", out.String())
	}
}

// TestDiffGateZeroTimerOnlyWhereTheBallGrantsNone anchors the timers rule to
// the shipped policy instead of a copied list: docs/evidence/s1/12-ball-walk
// records Acting with timers=yes (the whitelist allows it) while Sleeping must
// hold timers=no (SPEC-08 §2). A hand-copied "static states" list would have
// redded that real run.
func TestDiffGateZeroTimerOnlyWhereTheBallGrantsNone(t *testing.T) {
	rows := []stateRow{
		gateRow("Sleeping", "no", 2120),
		gateRow("Acting", "yes", 4844), // recorded, legitimate: Animated(Acting)
		gateRow("Warm", "yes", 4625),   // recorded, legitimate: AnimAlpha path
		gateRow("Settling", "yes", 4692),
	}
	var out bytes.Buffer
	if err := gateDiffRows(&out, rows, diffGateSpec{minPx8: diffGateFloorPx8}); err != nil {
		t.Fatalf("gate redded states whose timers the policy allows: %v; output=%q", err, out.String())
	}

	// The same columns with a timer in a state the ball grants none: red.
	bad := []stateRow{gateRow("Sleeping", "yes", 2120)}
	var badOut bytes.Buffer
	if err := gateDiffRows(&badOut, bad, diffGateSpec{minPx8: diffGateFloorPx8}); err == nil {
		t.Fatalf("gate passed timers=yes in Sleeping; output=%q", badOut.String())
	}
	// "?" means the status file said neither true nor false - unproven is not
	// the same as "no", so it must fail too.
	unknown := []stateRow{gateRow("Sleeping", "?", 2120)}
	var unkOut bytes.Buffer
	if err := gateDiffRows(&unkOut, unknown, diffGateSpec{minPx8: diffGateFloorPx8}); err == nil {
		t.Fatalf("gate treated an unreadable timers column as a pass; output=%q", unkOut.String())
	}
}

// TestDiffGateMissingAndDuplicateRows: a short table is the shape that used to
// look like a clean run (the loop died before the row that would have failed).
func TestDiffGateMissingAndDuplicateRows(t *testing.T) {
	rows := []stateRow{gateRow("Sleeping", "no", 2120)}
	var out bytes.Buffer
	err := gateDiffRows(&out, rows, diffGateSpec{
		minPx8: diffGateFloorPx8,
		want:   []string{"Sleeping", "Listening"},
	})
	if err == nil {
		t.Fatalf("gate passed a table missing the requested Listening row; output=%q", out.String())
	}
	if !strings.Contains(out.String(), "Listening") {
		t.Errorf("gap report did not name the missing row; output=%q", out.String())
	}

	dup := []stateRow{gateRow("Sleeping", "no", 2120), gateRow("Sleeping", "no", 2120)}
	var dupOut bytes.Buffer
	if err := gateDiffRows(&dupOut, dup, diffGateSpec{minPx8: diffGateFloorPx8}); err == nil {
		t.Fatalf("gate passed a table that counts one state twice; output=%q", dupOut.String())
	}
}

// TestDiffTableRoundTripIsWhatTheGateReads: the gate judges a TABLE, so the
// reader must read exactly what writeDiffTable wrote. A column silently
// shifting position is how a gate ends up judging the wrong number - the
// failure mode ticket 71 is about.
func TestDiffTableRoundTripIsWhatTheGateReads(t *testing.T) {
	want := []stateRow{
		{state: "Sleeping", timers: "no", changed8: 2120, changed3: 2400, changed24: 900, pixels: 97344},
		{state: "Sleeping-dock-right", timers: "no", changed8: 1620, changed3: 1800, changed24: 700, pixels: 97344},
		{state: "Listening", timers: "yes", changed8: 4856, changed3: 5000, changed24: 2000, pixels: 97344},
	}
	path := filepath.Join(t.TempDir(), "diff-table.txt")
	if err := writeDiffTable(path, want); err != nil {
		t.Fatalf("writeDiffTable: %v", err)
	}
	got, err := readDiffTable(path)
	if err != nil {
		t.Fatalf("readDiffTable: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("read %d rows, wrote %d", len(got), len(want))
	}
	for i := range want {
		if got[i].state != want[i].state || got[i].timers != want[i].timers ||
			got[i].changed8 != want[i].changed8 || got[i].pixels != want[i].pixels {
			t.Errorf("row %d read back as %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestReadDiffTableRefusesATableItCannotJudge: no px column means no verdict,
// not a pass.
func TestReadDiffTableRefusesATableItCannotJudge(t *testing.T) {
	path := filepath.Join(t.TempDir(), "diff-table.txt")
	if err := os.WriteFile(path, []byte("state\ttimers\nSleeping\tno\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := readDiffTable(path); err == nil {
		t.Fatal("readDiffTable accepted a table without a px_delta_ge8 column")
	}
}

// TestDiffGateTableRefusesAnEmptyTable: "examined nothing" is not "clean" -
// the same sentence ticket 71 was opened with.
func TestDiffGateTableRefusesAnEmptyTable(t *testing.T) {
	path := filepath.Join(t.TempDir(), "diff-table.txt")
	if err := writeDiffTable(path, nil); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := gateDiffTable(&out, path, diffGateSpec{minPx8: diffGateFloorPx8}); err == nil {
		t.Fatalf("gate accepted a table with zero rows; output=%q", out.String())
	}
}

// TestGateFloorCannotBeLowered pins AC#1's "must not be tuned down to go
// green": the flag may raise the bar, never lower it below the documented
// floor, and 0 means "use the floor".
func TestGateFloorCannotBeLowered(t *testing.T) {
	for _, tc := range []struct {
		in      int
		want    int
		wantErr bool
	}{
		{0, diffGateFloorPx8, false},
		{diffGateFloorPx8, diffGateFloorPx8, false},
		{diffGateFloorPx8 + 1, diffGateFloorPx8 + 1, false},
		{diffGateFloorPx8 - 1, 0, true},
		{10, 0, true},
		{-5, 0, true},
	} {
		got, err := resolveGateMinPx8(tc.in)
		if tc.wantErr != (err != nil) {
			t.Errorf("resolveGateMinPx8(%d) err=%v, wantErr=%v", tc.in, err, tc.wantErr)
			continue
		}
		if !tc.wantErr && got != tc.want {
			t.Errorf("resolveGateMinPx8(%d)=%d, want %d", tc.in, got, tc.want)
		}
	}
}
