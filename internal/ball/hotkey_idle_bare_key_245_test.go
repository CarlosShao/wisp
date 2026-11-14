//go:build windows

package ball

// Ticket 245 AC#6②: the idle-roster nail judges the BARE KEY, not the slot
// COUNT.
//
// The defect this ticket is about is not "the idle ball holds three slots" - it
// is "the idle ball holds a key that Win32 takes from every other application
// on the desktop while nothing is being confirmed". The orchestrator re-cut the
// criterion on AC#6② for exactly that reason (ledger A478/A480; the ruling is
// quoted in the ticket): a COUNT can read green with a bare key bound (any slot
// whose binding the count never looked at), and it reads RED for a user who
// configured [hotkey] cancel to a modifier combination and wants that
// combination live - which steals nothing from anybody.
//
// So the shared judgment is:
//
//   - a LIVE registration with no real modifier (Ctrl/Alt/Shift/Win) is the
//     defect, whichever slot carries it. The shipped cancel default is the bare
//     "Esc"; a bare "F9" configured for summon is the same defect under a
//     different name, and the old count-anchored nail could not see it;
//   - a registration WITH a modifier is not that defect. Under this nail a
//     user-configured "Ctrl+Alt+V" cancel may be held while idle.
//
// What makes the rule reachable is the 2026-11-14 AC#6② ruling: the idle pass
// never registers a BARE cancel binding (registerAllWith's skip), and it does
// register one that carries a real modifier. Both halves are pinned with attempt
// counts by TestConfiguredCancelIdleBindingFollowsModifierRuling245. The bare-key
// rule is the second, independent ruler: it reds on the shipped default (bare Esc)
// whether it arrives as a cancel borrow, as a configured cancel, or as a bare key
// in one of the other three slots.

import (
	"fmt"
	"strings"
	"testing"
)

// realModifiers245 are the Win32 modifier bits that make a hot key a
// COMBINATION. MOD_NOREPEAT is deliberately not one of them: it suppresses key
// repeat, and {MOD_NOREPEAT, VK_ESCAPE} is precisely the bare Esc that
// RegisterHotKey hands to nobody else.
const realModifiers245 = modAlt | modControl | modShift | modWin

// idleBareKeyHeld245 returns the first LIVE binding in rep whose accelerator
// carries no real modifier.
func idleBareKeyHeld245(rep HotkeyReport) (HotkeyBinding, bool) {
	for _, b := range rep.bindings {
		if b.Status != HotkeyLive {
			continue
		}
		if b.Acc.Mods&realModifiers245 == 0 {
			return b, true
		}
	}
	return HotkeyBinding{}, false
}

// idleRosterReds245 is the one idle-roster judgment every tier reads: an empty
// slice means "an idle ball may hold this set". It replaces the
// "len(Live()) != 3" anchor on the nails ticket 245 AC#3①/AC#6② named.
func idleRosterReds245(rep HotkeyReport) []string {
	var reds []string
	if b, held := idleBareKeyHeld245(rep); held {
		reds = append(reds, fmt.Sprintf("ticket 245 RED: the idle ball holds %s = %q as a bare key "+
			"(%+v, no Ctrl/Alt/Shift/Win), and RegisterHotKey takes a bare key from the WHOLE desktop - "+
			"every other application has lost it for as long as the ball is idle", b.Name, b.Binding, b.Acc))
	}
	for _, id := range []uint32{hkSummon, hkMute, hkPanel} {
		if rep.IsLive(id) {
			continue
		}
		if b, ok := rep.Binding(id); ok {
			reds = append(reds, fmt.Sprintf("the idle roster is missing %s (%q): %s",
				b.Name, b.Binding, b.Status))
		} else {
			reds = append(reds, fmt.Sprintf("the idle roster has no line for id %d", id))
		}
	}
	return reds
}

// escIdleProbe245 is AC#8's evidence shape: the same bare-Esc question
// escBorrowProbe asks, plus the numbers a reader needs when the answer is
// "somebody else owns it" - how many attempts were made and what Win32 said.
//
// WHERE THE ERRNO COMES FROM. escBorrowProbe maps exactly one Win32 answer to
// (free=false, err=nil): ERROR_HOTKEY_ALREADY_REGISTERED (1409). Every other
// refusal is returned as an error and is fatal at the call site, so a not-free
// reading is 1409 by construction; recording that constant is not a guess. Any
// future widening of the probe's error mapping must widen this struct too, which
// is why the attribution lives here and not at three call sites.
type escIdleProbe245 struct {
	attempts int
	free     bool
	errs     []error
}

// escIdleOwner245 is ticket 245 AC#8's three-way answer to "the bare-Esc probe
// says the key is claimed": free (the idle ball leaves the desktop its key),
// ours (the ball's own registration set holds a bare key - a PRODUCT defect), or
// foreign (our set is clean, so another program owns the key - an environment
// limit, never a product verdict).
//
// It lives in this default-tier file, not in hotkey_live_test.go, because the
// split itself is pure: hotkey_esc_idle_owner_245_test.go pins all three
// verdicts without a desktop, and the winlive ruler (requireIdleRoster) reads
// the same function.
type escIdleOwner245 int

const (
	escIdleFree245 escIdleOwner245 = iota
	escIdleBall245
	escIdleForeign245
)

// classifyEscIdleOwner245 is the split ticket 245 AC#8 asks for. Order matters:
// a report that names a bare key is OURS even if a foreign program is also
// squatting, because the product is already wrong and must not hide behind the
// environment.
func classifyEscIdleOwner245(rep HotkeyReport, p escIdleProbe245) escIdleOwner245 {
	if p.free {
		return escIdleFree245
	}
	if _, held := idleBareKeyHeld245(rep); held {
		return escIdleBall245
	}
	return escIdleForeign245
}

// liveLine245 builds one LIVE line for the synthetic controls below. It parses
// the binding with the production parser so the controls exercise the same
// accelerator shapes RegisterHotKey would see.
func liveLine245(t *testing.T, name string, id uint32, binding string) HotkeyBinding {
	t.Helper()
	acc, err := ParseAccelerator(binding)
	if err != nil {
		t.Fatalf("fixture binding %q does not parse: %v", binding, err)
	}
	return HotkeyBinding{Name: name, ID: id, Binding: binding, Status: HotkeyLive, Acc: acc}
}

// TestAC6IdlePinJudgesTheBareKeyNotTheCount245 is AC#6②'s completion criterion
// in one case: the count-blind-spot control, two positive controls (a bare key
// bound while idle must be red; a user-configured combination bound while idle
// must not be), and the shipped defaults read green.
func TestAC6IdlePinJudgesTheBareKeyNotTheCount245(t *testing.T) {
	// (1) THE COUNT-BLIND-SPOT. summon = "F9" is a binding the schema accepts
	// and Win32 grants desktop-wide. The idle set is still three entries with
	// cancel standing by, so the "exactly three" anchor this ticket replaced
	// reads GREEN while the desktop has lost F9.
	reg := &fakeRegistry{}
	bareSummon := HotkeyConfig{Summon: "F9", Mute: "Ctrl+Alt+M", Cancel: "Esc", Panel: "Ctrl+Alt+P"}
	rep := registerAllWith(0, bareSummon, reg.register)
	if n := len(rep.Live()); n != 3 {
		t.Fatalf("fixture reading = %d live entries, want 3 - this control only means something while the "+
			"count anchor still reads green: %+v", n, rep.Bindings())
	}
	if rep.IsLive(hkCancel) {
		t.Fatalf("fixture reading: cancel is supposed to be standby in this control: %+v", rep.Bindings())
	}
	if _, held := idleBareKeyHeld245(rep); !held {
		t.Fatalf("the bare-key rule did not see the bare F9 the idle pass registered: %+v", rep.Bindings())
	}
	if reds := idleRosterReds245(rep); len(reds) != 1 ||
		!strings.Contains(reds[0], "summon") || !strings.Contains(reds[0], "F9") {
		t.Fatalf("the idle-roster nail did not name the bare summon key: %v (report %+v)", reds, rep.Bindings())
	}

	// (2) POSITIVE CONTROL A: the shipped defect shape - the bare Esc held while
	// idle, here in the cancel slot. This must be red, and it must name the key.
	bareEsc := HotkeyReport{bindings: []HotkeyBinding{
		liveLine245(t, "summon", hkSummon, "Ctrl+Alt+Q"),
		liveLine245(t, "mute", hkMute, "Ctrl+Alt+M"),
		liveLine245(t, "cancel", hkCancel, "Esc"),
		liveLine245(t, "panel", hkPanel, "Ctrl+Alt+P"),
	}}
	reds := idleRosterReds245(bareEsc)
	if len(reds) != 1 || !strings.Contains(reds[0], "cancel") || !strings.Contains(reds[0], "Esc") {
		t.Fatalf("a bare Esc held while idle did not red with a sentence that names it: %v", reds)
	}

	// (3) POSITIVE CONTROL B: the same slot held as the user's own modifier
	// combination. It steals nothing from another application, so the nail must
	// NOT call it a defect (AC#6②: otherwise the person who configures this is
	// told their own configuration is a bug, or the implementation is pushed into
	// silently not binding it).
	modified := HotkeyReport{bindings: []HotkeyBinding{
		liveLine245(t, "summon", hkSummon, "Ctrl+Alt+Q"),
		liveLine245(t, "mute", hkMute, "Ctrl+Alt+M"),
		liveLine245(t, "cancel", hkCancel, "Ctrl+Alt+V"),
		liveLine245(t, "panel", hkPanel, "Ctrl+Alt+P"),
	}}
	if reds := idleRosterReds245(modified); len(reds) != 0 {
		t.Fatalf("a user-configured modifier combination held while idle was called a defect: %v", reds)
	}
	if _, held := idleBareKeyHeld245(modified); held {
		t.Fatalf("the bare-key rule flagged a modifier combination as bare: %+v", modified.Bindings())
	}

	// (4) The shipped defaults read green - the rule is not red-by-construction.
	def := registerAllWith(0, DefaultHotkeys(), (&fakeRegistry{}).register)
	if reds := idleRosterReds245(def); len(reds) != 0 {
		t.Fatalf("the shipped defaults did not read as a clean idle roster: %v", reds)
	}
}
