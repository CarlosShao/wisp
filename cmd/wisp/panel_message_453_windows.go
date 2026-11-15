//go:build windows

package main

// Ticket 453: the late bind of the panel message relay.
//
// The panel chain is assembled at resident_windows.go:359, the task pipeline at
// :548. The door therefore leaves its relay nil at construction, and this call -
// on the first line where both objects exist - is what turns a composed panel
// message into a task. It is deliberately a separate function rather than three
// lines inline: the binding itself is measurable (a door whose relay is nil
// refuses by name; after this call it delivers), and the resident boot has no
// test that can run it end to end.
//
// A boot with no panel, no dispatch chain or no task source leaves the relay nil.
// That is the honest state of the CLI seam and of a refused pipeline: the door
// then refuses plain messages by name instead of dropping them.
func bindPanelMessageRelay453(m *PanelManager, src *residentTaskSource) {
	if m == nil || m.disp == nil || src == nil {
		return
	}
	door, ok := m.disp.Message.(*panelMessageDoor453)
	if !ok || door == nil {
		return
	}
	door.bindRelay453(src.SubmitPanelMessage453)
}
