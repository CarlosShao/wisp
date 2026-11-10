//go:build windows

package ball

// Ticket-65 integration leg (A1594) + ticket 418: the *Ball half of the orb
// debug face. These methods read fields defined only in ball_windows.go (orb,
// mu, curState), so they live under the windows constraint - the
// platform-neutral helpers stay in orbdebug.go and compile everywhere.

import (
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// DebugOrbAttached reports whether the WebView2 orb owns the pixels (debug
// harness/evidence only; false = the D2D fallback is what is on screen).
// Same shape as DebugHWND's debug-only contract.
func (b *Ball) DebugOrbAttached() bool { return b.orb.Load() != nil }

// DebugOrbState returns the ball's current state under its own lock (the
// read side of SetState for the debug face).
func (b *Ball) DebugOrbState() statemachine.State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.curState
}
