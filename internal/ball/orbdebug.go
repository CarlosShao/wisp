package ball

// Ticket-65 integration leg (A1594): the pure, platform-neutral debug face of
// the orb takeover. The debug harness (balldebug) prints this line so the
// owner can see WHICH surface is rendering the ball and WHICH upstream orb
// state the page was pushed to, without opening a debugger. No Win32, no
// WebGPU, no locks - everything here reads values the caller already holds.
//
// The two *Ball debug methods (DebugOrbAttached/DebugOrbState) live in
// orbdebug_windows.go: they read fields defined only in ball_windows.go, so
// they must not compile into non-Windows builds (ticket 418).

import (
	"fmt"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// OrbSurfaceName names the rendering surface from the attach decision. The
// two names are the whole truth: the WebView2 orb either took over ("webview2
// -orb") or it did not and the D2D fallback is what you are looking at.
func OrbSurfaceName(attached bool) string {
	if attached {
		return "webview2-orb"
	}
	return "d2d-fallback"
}

// OrbMappedState returns the upstream orb state name the integration maps s
// onto ("idle" or "thinking" - the page's own stateSeeds keys, orbstate.go).
func OrbMappedState(s statemachine.State) string { return orbStateFor(s) }

// OrbDebugLine is the one-line debug face: surface + mapped state. The
// attached flag comes from Ball.OrbAttached; s from the machine's current
// state. Pure so every claim in the evidence doc is re-derivable anywhere.
func OrbDebugLine(attached bool, s statemachine.State) string {
	return fmt.Sprintf("orb-surface=%s orb-state=%s", OrbSurfaceName(attached), orbStateFor(s))
}
