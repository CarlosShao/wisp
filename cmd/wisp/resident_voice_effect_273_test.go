package main

// Ticket 273, wave 2: the pure unit test for effectFanout273 (see
// resident_voice_effect_273.go). No machine, no window, no disk: the only thing
// under test is the routing decision, which is the part R1 is about.
//
// WHY THESE THREE ASSERTIONS AND NOT "IT DOES NOT PANIC". The failure mode this
// seam can produce is not a crash, it is an effect that was booked but never
// acted on (feature looks dead) or an effect that was acted on twice (mic closes
// twice, TTS stops a second time). So each case measures the CALLS on both
// injected halves and compares the exact effect values:
//
//  1. every one of the four drive-grade names reaches both halves;
//  2. mic.keep-off - already executed inline by the resident voice leg - and
//     tts.stop reach audit only, and drive is NOT called (empty slice, not
//     "not nil");
//  3. a name the frozen table does not contain yet reaches audit only, which is
//     how a future D43 name stays audit-only until someone edits drive273.
//
// The fourth case pins the set itself: exactly the four names, each one spelled
// as a real side effect in statemachine.Table, and the two inline-executed names
// explicitly outside it. Without it, a fifth name could be added to drive273 and
// every case above would still be green - which is precisely the regression R1
// warns about.

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// fanoutCalls273 records what the sink under test handed each half. The two
// slices are the assertion surface; nothing here reads production state.
type fanoutCalls273 struct {
	drive []statemachine.Effect
	audit []statemachine.Effect
}

// fanout builds the injected value: both halves append, neither does work.
func (c *fanoutCalls273) fanout() effectFanout273 {
	return effectFanout273{
		drive: func(e statemachine.Effect) { c.drive = append(c.drive, e) },
		audit: func(e statemachine.Effect) { c.audit = append(c.audit, e) },
	}
}

// assertFanoutCalls273 compares one recorded half with the exact effects the
// sink was fed, in order. Comparing values (Effect is comparable) rejects a
// recipient that invents a D43 row or rewrites the name while routing.
func assertFanoutCalls273(t *testing.T, half string, got, want []statemachine.Effect) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s called %d time(s), want %d; calls: %v", half, len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("%s call %d = %+v, want %+v", half, i, got[i], want[i])
		}
	}
}

// TestEffectFanout273DriveNamesFireBothHalves is assertion 1. Each case is one
// member of drive273 with the frozen row it is fired from (table.go), so the
// fixture is the table's spelling rather than a convenience string.
func TestEffectFanout273DriveNamesFireBothHalves(t *testing.T) {
	cases := []struct {
		name string
		d43  int
	}{
		{"panel.stream-push", 15},
		{"panel.keep-alive", 26},
		{"conversation.ring-solid", 28},
		{"mic.enable", 28},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var calls fanoutCalls273
			want := statemachine.Effect{D43: tc.d43, Name: tc.name}
			calls.fanout().sink(want)

			assertFanoutCalls273(t, "drive", calls.drive, []statemachine.Effect{want})
			assertFanoutCalls273(t, "audit", calls.audit, []statemachine.Effect{want})
		})
	}
}

// TestEffectFanout273InlineExecutedNamesAreAuditOnly is assertion 2. mic.keep-off
// is the row #26 inline branch (the resident leg has already closed the mic);
// tts.stop is row #27's action and a second stop is audible. Both must be booked
// and neither may be driven, even though an executor is injected.
func TestEffectFanout273InlineExecutedNamesAreAuditOnly(t *testing.T) {
	cases := []struct {
		name string
		d43  int
	}{
		{"mic.keep-off", 26},
		{"tts.stop", 27},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var calls fanoutCalls273
			want := statemachine.Effect{D43: tc.d43, Name: tc.name}
			calls.fanout().sink(want)

			assertFanoutCalls273(t, "audit", calls.audit, []statemachine.Effect{want})
			// Explicitly empty, not merely "no mismatch": the whole point of
			// auditOnly is that the action half never sees this name.
			if len(calls.drive) != 0 {
				t.Fatalf("%s reached drive %d time(s), want 0 (inline already executes it; R1 forbids a second run); calls: %v",
					tc.name, len(calls.drive), calls.drive)
			}
		})
	}
}

// TestEffectFanout273UnknownNameIsAuditOnly is assertion 3. The name is
// deliberately one that is NOT in statemachine.Table: this is the default-deny
// path a future frozen row would take.
func TestEffectFanout273UnknownNameIsAuditOnly(t *testing.T) {
	const name = "panel.future-effect-273"
	var calls fanoutCalls273
	want := statemachine.Effect{D43: 15, Name: name}
	calls.fanout().sink(want)

	assertFanoutCalls273(t, "audit", calls.audit, []statemachine.Effect{want})
	if len(calls.drive) != 0 {
		t.Fatalf("an effect name outside drive273 reached drive %d time(s), want 0; default must be audit-only; calls: %v",
			len(calls.drive), calls.drive)
	}
}

// TestEffectFanout273DriveSetIsFrozen is the guard against the drive set growing
// silently. The expected names are spelled here, not read back from drive273, so
// adding a fifth name on purpose still requires editing a test that says why.
func TestEffectFanout273DriveSetIsFrozen(t *testing.T) {
	want := []string{
		"panel.stream-push",
		"panel.keep-alive",
		"conversation.ring-solid",
		"mic.enable",
	}
	if len(drive273) != len(want) {
		t.Fatalf("drive273 has %d name(s), want exactly %d (%v); every added name claims its inline path does NOT already act",
			len(drive273), len(want), want)
	}
	for _, name := range want {
		if !drive273[name] {
			t.Errorf("drive273 is missing %q: its action half would silently become auditOnly", name)
		}
	}

	// Every drive name must exist in the frozen table. This is what catches a
	// typo: a misspelled key would never match a fired effect and would demote a
	// drive-grade name to auditOnly with no red test anywhere else.
	frozen := map[string]bool{}
	for _, row := range statemachine.Table {
		for _, name := range row.SideEffects {
			frozen[name] = true
		}
	}
	for name := range drive273 {
		if !frozen[name] {
			t.Errorf("drive273 contains %q, which is not a side effect in statemachine.Table: a drive name that never fires is dead routing", name)
		}
	}

	// The two names the resident leg already executes inline must stay out of
	// the drive set; this is R1 stated as an assertion rather than a comment.
	for _, name := range []string{"mic.keep-off", "tts.stop"} {
		if drive273[name] {
			t.Errorf("%q is in drive273 but the resident voice leg already executes it inline: driving it here would run the action twice", name)
		}
	}
}

// TestEffectFanout273InjectionHalvesDegradeIndependently pins the two nil
// branches, because production will install the audit half before the drive
// half. A missing audit must not stop the drive, and a missing drive must not
// stop the audit; both missing is a no-op rather than a panic.
func TestEffectFanout273InjectionHalvesDegradeIndependently(t *testing.T) {
	driveEffect := statemachine.Effect{D43: 28, Name: "mic.enable"}
	auditOnlyEffect := statemachine.Effect{D43: 26, Name: "mic.keep-off"}

	var audited []statemachine.Effect
	effectFanout273{audit: func(e statemachine.Effect) { audited = append(audited, e) }}.sink(driveEffect)
	assertFanoutCalls273(t, "audit", audited, []statemachine.Effect{driveEffect})

	var driven []statemachine.Effect
	effectFanout273{drive: func(e statemachine.Effect) { driven = append(driven, e) }}.sink(driveEffect)
	assertFanoutCalls273(t, "drive", driven, []statemachine.Effect{driveEffect})

	// Both halves absent: callable, and it must not drive an auditOnly name.
	effectFanout273{}.sink(auditOnlyEffect)
}
