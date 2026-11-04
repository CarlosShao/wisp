package main

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// Ticket 273 wave 2b-1: proves the PRODUCTION constructor wires the audit sink.
//
// buildResidentVoice (resident_voice_windows.go:1789) now hands its machine
// `Sink: effectFanout273{audit: voiceAudit}.sink`. This test drives that machine
// one legal, unguarded step - StateSleeping + EvSummon -> row #4 - whose two
// side effects (session.scope-create, speech.load-vad-asr) are both auditOnly,
// and asserts the injected audit func receives them.
//
// THE DISCRIMINATOR IS THE PRODUCTION SINK LINE. Delete `Sink:` at :1794 (or
// pass nil for voiceAudit) and the machine still fires the exact same two
// effects, but the injected audit is never called: got stays empty and this test
// goes red. That is the reverse criterion AC#4 asks for, anchored on the shipped
// constructor rather than on a machine the test built for itself.
func TestResidentVoiceMachine273AuditIsWired(t *testing.T) {
	var got []statemachine.Effect
	rv := buildResidentVoice(t.TempDir(), nil, nil, nil, func(e statemachine.Effect) {
		got = append(got, e)
	})
	if rv == nil || rv.machine == nil {
		t.Fatal("buildResidentVoice returned no machine")
	}
	defer rv.machine.Close()

	if _, err := rv.machine.Dispatch(statemachine.EvSummon, &statemachine.Facts{}); err != nil {
		t.Fatalf("EvSummon dispatch: %v", err)
	}

	// Row #4 fires both names; either one missing means the audit half is not
	// wired into the production machine.
	want := map[string]bool{"session.scope-create": false, "speech.load-vad-asr": false}
	for _, e := range got {
		if e.D43 == 4 {
			if _, ok := want[e.Name]; ok {
				want[e.Name] = true
			}
		}
	}
	for name, seen := range want {
		if !seen {
			t.Fatalf("audit sink not wired into the production machine: d43=4 name=%q never reached the injected audit; got=%+v", name, got)
		}
	}
}
