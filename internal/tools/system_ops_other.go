//go:build !windows

package tools

// The non-Windows half of the machine-control family (ticket 23).
//
// Wisp ships Windows-only (D7, SPEC-01 §2), but internal/tools must still BUILD
// on the Linux runner so the risk matrix, the refusal texts and the taint rules
// around system.* / window.* / clipboard.* / input.type stay testable where CI
// runs. This file is what makes that build possible: platformSystemBackend is
// the single symbol the Windows file owns (system_ops_windows.go:84, reached
// through SystemBridge.funcs at system_bridge.go:129), so answering here is the
// whole port - every other file of the family is platform-neutral.
//
// What it must never do is pretend. There is no xdotool / AppleScript shim
// behind this seam on purpose: a POSIX "system.set.volume" that reported
// 「音量已设为 40%」 without ever touching a mixer would turn the D34 risk table
// into fiction on exactly the platform where nobody can check it by hand, and
// the Windows legs (window handles, UIPI, SendInput) have no honest analogue
// here anyway.
//
// Returning nil + ErrNoSystemBackend therefore lands every verb on the family's
// existing fail-closed path: SystemBridge.funcs reports the refusal, Do hands
// it back untouched, and the tools render it as IsError text the model can act
// on (D37) - byte for byte the answer the zero-value SystemBridge already
// gives. "This build cannot reach the machine" is a true sentence; "the machine
// was touched" is not.
func platformSystemBackend() (systemFuncs, error) {
	return nil, ErrNoSystemBackend
}
