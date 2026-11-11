//go:build windows

package main

// The appearance wave's assembly-root pin: run.go's pipeline pump and
// panel_resident_pump.go's degraded-shape pump are the two PumpSources
// assembly sites, and this wave's four readers must ride both - without
// their reader lines the snapshot keys would silently vanish from every
// packet this process pushes, the exact "carrier exists, nothing reads it"
// hop the internal/panel byte nails cannot see from over there (the same
// disk reading panel_ball_size_422_test.go makes for the ballSize reader).
//
// The reverse criterion: deleting any reader line from either literal
// reddens its case here; deleting the feed it reads reddens the config-side
// round-trip in the wave's cmd/wisp wire test.

import (
	"os"
	"strings"
	"testing"
)

func TestWaveAssemblyRootsCarryTheReaders(t *testing.T) {
	runGo, err := os.ReadFile("run.go")
	if err != nil {
		t.Fatalf("read run.go: %v", err)
	}
	for _, line := range []string{
		"BallOpacityIdle: func() panel.BallOpacityIdleSection",
		"BallClickThrough: func() panel.BallClickThroughSection",
		"PanelWidth: func() panel.PanelWidthSection",
		"PanelFontSize: func() panel.PanelFontSizeSection",
	} {
		if !strings.Contains(string(runGo), line) {
			t.Errorf("run.go's PumpSources literal no longer carries %q: the key would silently vanish from every packet this process pushes (the appearance wave)", line)
		}
	}
	residentGo, err := os.ReadFile("panel_resident_pump.go")
	if err != nil {
		t.Fatalf("read panel_resident_pump.go: %v", err)
	}
	for _, line := range []string{
		"src.BallOpacityIdle = func() panel.BallOpacityIdleSection",
		"src.BallClickThrough = func() panel.BallClickThroughSection",
		"src.PanelWidth = func() panel.PanelWidthSection",
		"src.PanelFontSize = func() panel.PanelFontSizeSection",
	} {
		if !strings.Contains(string(residentGo), line) {
			t.Errorf("panel_resident_pump.go no longer carries %q: the shipping double-click shape would push packets without the key (the appearance wave)", line)
		}
	}
}
