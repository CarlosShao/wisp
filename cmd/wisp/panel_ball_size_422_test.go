package main

// Ticket 422's assembly-root pin: run.go's PumpSources literal hands the pump
// the ballSize reader, the same disk reading panel_ball_style_65_test.go makes
// for the ballStyle reader. Without this line the snapshot's size key would
// silently vanish from every packet this process pushes - the exact
// "carrier exists, nothing reads it" hop the reader-gated tests in
// internal/panel cannot see from over there, and the assembly gap the 413
// contract (docs/evidence/s1/413-ball-size-slider-contract.md §2.1 piece 3)
// named for this leg.

import (
	"os"
	"strings"
	"testing"
)

func Test422AssemblyRootCarriesTheBallSizeReader(t *testing.T) {
	runGo, err := os.ReadFile("run.go")
	if err != nil {
		t.Fatalf("read run.go: %v", err)
	}
	if !strings.Contains(string(runGo), "BallSize: func() panel.BallSizeSection") {
		t.Error("run.go's PumpSources literal no longer carries the ballSize reader: the snapshot's size key would silently vanish from every packet this process pushes (ticket 422)")
	}
}
