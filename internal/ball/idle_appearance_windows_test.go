//go:build windows

package ball

// The appearance wave's pure-core rulers: the two folds are functions of
// (options, state, current value) precisely so these can be judged without a
// window. The reverse criteria live in the folds' own callers - deleting the
// visualForLocked fold reddens the applyStateLocked wiring, deleting the
// applyClickThroughStyleLocked call reddens the state-edge wiring, and
// deleting the win32 procs reddens the build.
//
// The zero case in the opacity ruler is load-bearing beyond its own line: 0
// is a real value inside the schema band, and a plain float64 option would
// have swallowed it (the lying-option shape) - the pointer exists so the
// carrier cannot confuse "the file said 0" with "nobody read a file".

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

func waveFptr(v float64) *float64 { return &v }

// TestWaveIdleVisualOpacityIsSleepingOnly: Sleeping carries the configured
// value, every other state keeps its own visual's opacity, and nil means the
// standing value.
func TestWaveIdleVisualOpacityIsSleepingOnly(t *testing.T) {
	sleeping := VisualFor(DarkPalette(), 56, statemachine.StateSleeping, 0)
	warm := VisualFor(DarkPalette(), 56, statemachine.StateWarm, 0)
	fed := Options{OpacityIdle: waveFptr(0.9)}
	if got := idleVisualOpacity(fed, statemachine.StateSleeping, sleeping.Opacity); got != 0.9 {
		t.Fatalf("Sleeping opacity with a fed option = %v, want the configured 0.9", got)
	}
	if got := idleVisualOpacity(fed, statemachine.StateWarm, warm.Opacity); got != warm.Opacity {
		t.Fatalf("Warm opacity with a fed option = %v, want the state's own %v (the override is idle-only)", got, warm.Opacity)
	}
	if got := idleVisualOpacity(Options{}, statemachine.StateSleeping, sleeping.Opacity); got != sleeping.Opacity {
		t.Fatalf("Sleeping opacity with no reading = %v, want the standing %v", got, sleeping.Opacity)
	}
	if got := idleVisualOpacity(Options{OpacityIdle: waveFptr(0)}, statemachine.StateSleeping, sleeping.Opacity); got != 0 {
		t.Fatalf("Sleeping opacity with a configured 0 = %v, want 0 (a real band value must not be swallowed)", got)
	}
}

// TestWaveClickThroughExStyleFoldsByState: the bit rides only in Sleeping
// under a true option, is cleared in every other state, and a false option
// keeps the word byte-for-byte what the window was born with - the pre-wave
// behavior.
func TestWaveClickThroughExStyleFoldsByState(t *testing.T) {
	base := uintptr(wsExLayered | wsExTopmost | wsExToolWindow | wsExNoActivate)
	on := Options{ClickThrough: true}
	off := Options{}
	if got := clickThroughExStyle(base, on, statemachine.StateSleeping); got != base|wsExTransparent {
		t.Fatalf("Sleeping ex-style with the switch on = %#x, want the bit set (%#x)", got, base|wsExTransparent)
	}
	if got := clickThroughExStyle(base|wsExTransparent, on, statemachine.StateWarm); got != base {
		t.Fatalf("Warm ex-style = %#x, want the bit cleared (%#x): the active ball must stay clickable", got, base)
	}
	if got := clickThroughExStyle(base, off, statemachine.StateSleeping); got != base {
		t.Fatalf("Sleeping ex-style with the switch off = %#x, want the untouched base word (%#x)", got, base)
	}
	if got := clickThroughExStyle(base|wsExTransparent, on, statemachine.StateSleeping); got != base|wsExTransparent {
		t.Fatalf("re-applying the same answer moved the word: %#x", got)
	}
}
