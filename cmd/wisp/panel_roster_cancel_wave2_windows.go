//go:build windows

package main

// Readface wave 2 (2026-11-17), ticket 38 AC#4: the late bind of the panel
// cancel seam. The panel chain is assembled before the task pipeline exists
// (resident_windows.go's ordering, the same gap ticket 453's relay crossed),
// so the door leaves its seam nil at construction and this call - on the
// first line where both objects exist - hands it the runtime's own cancel
// paths (agentRuntime.cancelRosterRowWave2: the root task's RunningTask
// handle, then the TaskRoster's one per-row stop path).
//
// A boot with no panel, no dispatch chain or no task source leaves the seam
// nil. That is the honest state of the CLI seam and of a refused pipeline:
// the door then refuses cancels by name instead of dropping them.
func bindPanelRosterCancel(m *PanelManager, src *residentTaskSource) {
	if m == nil || m.disp == nil || src == nil || src.run == nil {
		return
	}
	door, ok := m.disp.RosterCancel.(*rosterCancelDoor)
	if !ok || door == nil {
		return
	}
	door.bindSeam(src.run.cancelRosterRowWave2)
}
