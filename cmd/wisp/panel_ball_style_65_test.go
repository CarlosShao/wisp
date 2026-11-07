package main

// Ticket 65's assembly-root pin: run.go's PumpSources literal hands the pump
// the ballStyle reader, the same disk reading panel_commands_213_test.go makes
// for the commands reader. Without this line the snapshot's appearance key
// would silently vanish from every packet this process pushes - the exact
// "carrier exists, nothing reads it" hop ticket 200 AC#7 closed for its own
// dimension, and the shape the reader-gated tests in internal/panel cannot see
// from over there.

import (
	"os"
	"strings"
	"testing"
)

func Test65AssemblyRootCarriesTheBallStyleReader(t *testing.T) {
	runGo, err := os.ReadFile("run.go")
	if err != nil {
		t.Fatalf("read run.go: %v", err)
	}
	if !strings.Contains(string(runGo), "BallStyle: func() panel.BallStyleSection") {
		t.Error("run.go's PumpSources literal no longer carries the ballStyle reader: the snapshot's appearance key would silently vanish from every packet this process pushes (ticket 65)")
	}
}
