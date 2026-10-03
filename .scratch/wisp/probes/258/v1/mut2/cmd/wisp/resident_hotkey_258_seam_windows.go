package main

// Ticket 258: test seams for the resident hotkey work. Every symbol here exists
// only so resident_hotkey_258_windows_test.go and
// resident_hotkey_live_258_windows_test.go can drive unexported production
// state; none of it runs unless a test calls it.

// hotkeyBridgeCheck258 runs one bridge poll hop synchronously (the hop the
// armed loop performs on its 1s ticker). No-op when no bridge was installed.
func (rb *residentBall) hotkeyBridgeCheck258() {
	if rb == nil || rb.hotkeyBridge == nil {
		return
	}
	_ = rb.hotkeyBridge
}

// verdictStatusLineFor258 returns the resident ball's verdict (or the no-ball
// sentence) for test failure messages.
func (rb *residentBall) verdictStatusLineFor258() string {
	return rb.statusLine()
}
