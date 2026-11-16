package panel

// Ticket 33 AC#8's sweep ruler (33-lane3, 2026-11-16): the one-time "which pump
// branches execute zero times" reading, taken at HEAD d28826829 with the whole
// internal/panel package under -coverprofile (86.6% statements, ok). Two BEHAVIOUR
// blocks in pump.go came back zero-executed; AC#8 says each gets a test or a
// deletion, and both blocks carry a documented promise, so deletion is not on the
// table - these two cases are the tests:
//
//	pump.go Close's unknown-key branch ("the stream ended and nobody had noticed
//	it start is a fact the panel should see rather than lose") - zero runs: every
//	existing caller closes a key it had Appended first (pump_test.go:176/:226,
//	subagent_stream_197_test.go:162 all close their own keys);
//	pump.go Snapshot's empty-credential branch ("nobody read it is not the same
//	fact as nothing is recorded") - zero runs: no test drives src.Credential with
//	an empty answer, so the empty->CredentialUnknown rewrite never executed.
//
// Each case carries its own separating control so "the branch fired" cannot be
// confused with "this is just what the function always writes" - the same shape
// the 33-r14 clamp rulers use.

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// TestAC8CloseUnknownKeyStillRecordsADoneChunk33lane3 pins the unknown-key branch
// of StreamLog.Close: closing a key nobody appended still leaves a row - a done
// chunk with that key and no borrowed text - because a stream that ended unseen is
// a fact, not a non-event. The control half closes a key that WAS appended and
// asserts the two rows stay distinct, so this ruler cannot pass by an accident of
// the map's contents.
func TestAC8CloseUnknownKeyStillRecordsADoneChunk33lane3(t *testing.T) {
	sl := NewStreamLog(4)
	sl.Append("known", "streamed text")
	sl.Close("ghost") // never appended: the branch under test
	sl.Close("known")

	chunks := sl.Chunks()
	byKey := map[string]ResultChunk{}
	for _, c := range chunks {
		byKey[c.CorrelationID] = c
	}
	ghost, ok := byKey["ghost"]
	if !ok {
		t.Fatalf("Close on an unknown key left no row at all (%d chunks) - the "+
			"'the stream ended and nobody had noticed it start' fact is being dropped, "+
			"which is the exact loss the branch's own comment says it exists to prevent", len(chunks))
	}
	if !ghost.Done {
		t.Errorf("the unknown-key row reports done=false (%+v): a closed stream that was never seen must still read as ended", ghost)
	}
	if ghost.Text != "" {
		t.Errorf("the unknown-key row carries text %q - an unseen stream has no deltas to show and inventing one is a fabrication", ghost.Text)
	}
	known := byKey["known"]
	if !ok || known.Text != "streamed text" || !known.Done {
		t.Errorf("the control row degraded: %+v, want the appended text kept and done=true - the unknown-key handling must not disturb a key the log tracks", known)
	}
}

// TestAC8EmptyCredentialReaderAnswersUnknownNotSilent33lane3 pins the empty-answer
// branch of the pump's credential reader: a reader that EXISTS and answers "" must
// produce credentialState=unknown WITH credentialKnown=true, because "asked and
// empty" and "nobody asked" are two facts the page must be able to tell apart.
// Three arms: the empty answer (the zero-executed branch), a real answer travelling
// verbatim (the control), and no reader at all (credentialKnown stays false - the
// sentence ticket 216 AC#4 pinned).
func TestAC8EmptyCredentialReaderAnswersUnknownNotSilent33lane3(t *testing.T) {
	p := NewSnapshotPump(PumpSources{
		Mode:       func() risk.Mode { return risk.ModeAskEveryStep },
		Credential: func() CredentialState { return "" },
	})
	snap := p.Snapshot()
	if snap.Composer.Credential != CredentialUnknown {
		t.Errorf("an empty credential answer travelled verbatim as %q: the empty->unknown rewrite never fired, so the page sees an unrenderable state instead of the documented unknown", snap.Composer.Credential)
	}
	if !snap.Composer.CredentialKnown {
		t.Errorf("an empty credential answer reports credentialKnown=false: a reader that exists and answered must be distinguishable from no reader at all")
	}

	p2 := NewSnapshotPump(PumpSources{
		Mode:       func() risk.Mode { return risk.ModeAskEveryStep },
		Credential: func() CredentialState { return CredentialAllRecorded },
	})
	snap2 := p2.Snapshot()
	if snap2.Composer.Credential != CredentialAllRecorded || !snap2.Composer.CredentialKnown {
		t.Errorf("a real credential answer was rewritten: got %q known=%v, want all_recorded verbatim with known=true", snap2.Composer.Credential, snap2.Composer.CredentialKnown)
	}

	p3 := NewSnapshotPump(PumpSources{
		Mode: func() risk.Mode { return risk.ModeAskEveryStep },
	})
	snap3 := p3.Snapshot()
	if snap3.Composer.CredentialKnown {
		t.Errorf("with no credential reader at all credentialKnown=%v: nobody-asked must stay its own sentence, not borrow the empty-answer shape", snap3.Composer.CredentialKnown)
	}
}
