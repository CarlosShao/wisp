//go:build windows

package ball

// Ticket 245 AC#8: "somebody else owns the desktop's bare Esc" and "our own idle
// set holds it" are two different verdicts and must not share one red (or one
// skip). The live ruler cannot fabricate the environment half on demand, so the
// split itself is pinned here, in the default tier, against the same production
// judgments the live cases call:
//
//   - idleRosterReds245 names a bare key held by the ball (the product defect);
//   - classifyEscIdleOwner245 turns "the probe says claimed" into OURS only when
//     that product judgment fires, and FOREIGN otherwise.
//
// The live-tier half (a real foreign holder makes the four cases SKIP-LOUD with
// the probe's numbers) is in the ticket's evidence file; this case is the control
// that the classifier cannot silently collapse the two.

import (
	"testing"

	"golang.org/x/sys/windows"
)

// cleanIdleReport245 builds the roster an idle ball is allowed to hold: the
// three always-on slots live, cancel standing by on a user combination. The live
// lines are built by liveLine245 (hotkey_idle_bare_key_245_test.go) so this
// control parses accelerators with the production parser.
func cleanIdleReport245(t *testing.T) HotkeyReport {
	t.Helper()
	return HotkeyReport{bindings: []HotkeyBinding{
		liveLine245(t, "summon", hkSummon, "Ctrl+Alt+Q"),
		liveLine245(t, "mute", hkMute, "Ctrl+Alt+M"),
		liveLine245(t, "cancel", hkCancel, "Ctrl+Alt+V"),
		liveLine245(t, "panel", hkPanel, "Ctrl+Alt+P"),
	}}
}

// TestAC8IdleProbeSplitsForeignOwnerFromOurOwnSet245 is AC#8's completion
// criterion at the classification seam: free, foreign and ours get three
// different answers, and the foreign answer carries the attempt count + errno.
func TestAC8IdleProbeSplitsForeignOwnerFromOurOwnSet245(t *testing.T) {
	clean := cleanIdleReport245(t)
	if reds := idleRosterReds245(clean); len(reds) != 0 {
		t.Fatalf("fixture reading: the clean report was already red: %v", reds)
	}

	// (1) FREE: nothing claimed, nothing to classify.
	if got := classifyEscIdleOwner245(clean, escIdleProbe245{attempts: 1, free: true}); got != escIdleFree245 {
		t.Fatalf("a free bare-Esc probe classified as owner %d, want free (%d)", got, escIdleFree245)
	}

	// (2) FOREIGN: the probe was refused three times and the ball's own set holds
	// no bare key. This is an environment limit; requireIdleRoster SKIPs LOUD with
	// these numbers instead of calling it a product defect (AC#8 + the
	// orchestrator's 2026-11-14 dispatch).
	foreign := escIdleProbe245{
		attempts: 3,
		errs: []error{
			windows.ERROR_HOTKEY_ALREADY_REGISTERED,
			windows.ERROR_HOTKEY_ALREADY_REGISTERED,
			windows.ERROR_HOTKEY_ALREADY_REGISTERED,
		},
	}
	if got := classifyEscIdleOwner245(clean, foreign); got != escIdleForeign245 {
		t.Fatalf("a claimed bare Esc with a clean ball set classified as owner %d, want foreign (%d)", got, escIdleForeign245)
	}
	if foreign.attempts != 3 || len(foreign.errs) != 3 {
		t.Fatalf("fixture reading: the numbers the skip must carry are %d attempts / %d errs", foreign.attempts, len(foreign.errs))
	}
	if want := "Hot key is already registered."; windows.ERROR_HOTKEY_ALREADY_REGISTERED.Error() != want {
		t.Fatalf("errno 1409's text changed: got %q, want %q (the named skip quotes it verbatim)",
			windows.ERROR_HOTKEY_ALREADY_REGISTERED.Error(), want)
	}

	// (3) OURS: the same refused probe, but the ball's own set names a bare key.
	// This must stay a PRODUCT red even though the desktop is also occupied, so an
	// environment skip can never absorb our defect.
	ours := clean
	ours.bindings = append([]HotkeyBinding(nil), clean.bindings...)
	for i := range ours.bindings {
		if ours.bindings[i].ID == hkCancel {
			ours.bindings[i] = liveLine245(t, "cancel", hkCancel, "Esc")
		}
	}
	if reds := idleRosterReds245(ours); len(reds) != 1 {
		t.Fatalf("fixture reading: the bare cancel did not red the idle roster: %v", reds)
	}
	if got := classifyEscIdleOwner245(ours, foreign); got != escIdleBall245 {
		t.Fatalf("a claimed bare Esc that the ball's own set holds classified as owner %d, want ours (%d) - "+
			"the foreign-owner skip must never absorb the product's defect", got, escIdleBall245)
	}
}
