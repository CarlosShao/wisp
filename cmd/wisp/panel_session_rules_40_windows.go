//go:build windows

package main

// Readface wave 4 (2026-11-18), ticket 40 AC#2: the late bind of the panel
// session-rules seam. The panel chain is assembled before the task pipeline
// exists (resident_windows.go's ordering, the roster-cancel pattern), so the
// door leaves its seam nil at construction and this call - on the first line
// where both objects exist - hands it the runtime's own session ledger
// (agentRuntime.sessionRulesSeam40: ListBySession / Revoke / SessionID, the
// same doors ticket 49's console verbs walk).
//
// A boot with no panel, no dispatch chain or no task source leaves the seam
// nil. That is the honest state of the CLI seam and of a refused pipeline:
// the door then refuses reads and revokes by name instead of dropping them.
func bindPanelSessionRules(m *PanelManager, src *residentTaskSource) {
	if m == nil || m.disp == nil || src == nil || src.run == nil {
		return
	}
	door, ok := m.disp.SessionRules.(*sessionRulesDoor40)
	if !ok || door == nil {
		return
	}
	door.bindSeam(src.run.sessionRulesSeam40())
}
