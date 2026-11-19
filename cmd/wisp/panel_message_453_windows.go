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

// bindFailedSendLedger453 hands the pipeline pump the message door's failed-send
// slot (ticket 167 AC#5). The door owns the ledger - it is where sends fail -
// and the pipeline pump (run.go's PumpSources literal) reads through the runtime,
// so this bind crosses the ordering the relay bind above crosses: on the first
// line where both the panel chain and the task source exist, the door's ledger
// is bound onto src.run and the snapshot's failedSend* keys start carrying the
// slot the door writes. A boot with no panel, no dispatch chain or no task
// source leaves the runtime unbound - the pump then answers Known=false and
// both keys stay off the wire, the honest no-send-face shape.
func bindFailedSendLedger453(m *PanelManager, src *residentTaskSource) {
	if m == nil || m.disp == nil || src == nil || src.run == nil {
		return
	}
	door, ok := m.disp.Message.(*panelMessageDoor453)
	if !ok || door == nil {
		return
	}
	src.run.bindFailedSendLedger(&door.failed)
}
