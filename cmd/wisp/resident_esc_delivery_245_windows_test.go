//go:build windows && winlive

package main

// Ticket 245 AC#9, live tier: the "a second process really receives Esc" ruler
// becomes a standing case in THIS repository, run against the SHIPPED no-args
// process rather than only against an in-process Ball or a hand-started harness.
//
// WHY THIS CASE EXISTS. Ticket 245-v1 measured the delivery half with an
// acceptance-side rig (E1 normal=1 / E2 self-claim=0+wm_hotkey=1 / E3 factory
// process up=1 / E4 the defect=0). The orchestrator's ruling on AC#9 (ledger
// A480) was that an acceptance-side rig guards no regression, so the four shapes
// have to be runnable from the tree. The second-process program itself is already
// in-repo as cmd/wisp/testdata/esclistener (ticket 246-r1); what was missing was
// the 245-named case that runs the four shapes IN ORDER, with the shipped
// resident leg as "the factory process", and with no human driving anything.
//
// THE FOUR SHAPES, mapped one-to-one onto E1..E4:
//
//	baseline   - nothing of ours running: injected bare Esc reaches the observer
//	             (E1, keydown_esc=1);
//	control    - the SAME observer registers bare Esc itself: the injection
//	             vanishes from its keydown stream and arrives as WM_HOTKEY
//	             (E2, 0 + 1). Without this pair the "1" below would prove nothing,
//	             and it is also the "0" shape (E4) the defect produces;
//	product    - the shipped wisp.exe, ball window up, nothing waiting -> the
//	             observer STILL receives the Esc (E3, 1);
//	after      - that process leaves through its own shutdown path -> the key is
//	             deliverable again (E4's "returned" reading, 1).
//
// THE OWNERSHIP SPLIT (ticket 245 AC#8, carried here on purpose). When a reading
// comes back 0 while the product is up, the case must not call it a product
// defect before asking WHO holds the key. It releases the suspect (the resident
// leaves), then asks Win32 again: still claimed = a foreign program owns it and
// this desktop cannot measure delivery (SKIP-LOUD, with attempts and the verbatim
// errno); free again = the resident itself held it (a hard product RED). That is
// the same two-branch judgment internal/ball's requireIdleRoster now makes.
//
// Run it alone, on a machine with a desktop and no other Wisp ball:
//
//	GOFLAGS=-tags=winlive sh scripts/wisp-cli-staged-tests.sh \
//	  -test.run 'TestLive245'
//
// Nothing here depends on a human pressing a key or starting a process: the
// observer takes the foreground itself, the test triggers the injection over a
// stdin line, and every child is reaped by t.Cleanup.

import (
	"fmt"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

// bareEscVK245 / modNoRepeat245 are the production default pair
// (internal/ball's DefaultHotkeys cancel = "Esc"): the key the whole ticket is
// about, spelled here so the probe asks the same question Win32 answered for the
// defect.
const (
	bareEscVK245   = 0x1B
	modNoRepeat245 = 0x4000
	escProbeID245  = 0x9245
)

var (
	escRailUser32       = windows.NewLazySystemDLL("user32.dll")
	pEscRailRegisterHot = escRailUser32.NewProc("RegisterHotKey")
	pEscRailUnregHot    = escRailUser32.NewProc("UnregisterHotKey")
)

// bareEscClaimed245 asks Win32 whether the bare Esc is claimed right now, on
// this thread, and returns the numbers a named skip must carry: how many
// attempts were made and the verbatim error each one produced. It unregisters a
// successful probe immediately, so a "free" answer leaves nothing behind.
func bareEscClaimed245(t *testing.T) (claimed bool, attempts int, detail string) {
	t.Helper()
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var errs []string
	for i := 0; i < 3; i++ {
		attempts++
		r, _, e := pEscRailRegisterHot.Call(0, uintptr(escProbeID245), uintptr(modNoRepeat245), uintptr(bareEscVK245))
		if r != 0 {
			pEscRailUnregHot.Call(0, uintptr(escProbeID245))
			return false, attempts, strings.Join(errs, "; ")
		}
		errno, _ := e.(syscall.Errno)
		errs = append(errs, fmt.Sprintf("attempt %d: errno %d %q", attempts, uint32(errno), e.Error()))
		time.Sleep(100 * time.Millisecond)
	}
	return true, attempts, strings.Join(errs, "; ")
}

// TestLive245TheShippedLegLeavesTheDesktopsBareEscAlone is AC#9's standing
// ruler: the four shapes above, in one timeline, against the real binary.
func TestLive245TheShippedLegLeavesTheDesktopsBareEscAlone(t *testing.T) {
	rig := buildEscListener246(t)

	// SHAPE 1 (E1): baseline with nothing of ours running. A zero here is not
	// this ticket's product verdict: either the desktop already has a foreign
	// owner (SKIP-LOUD with numbers) or the rig itself is broken (RED).
	base := observeEsc246(t, rig, "-watch")
	if base.keydown != 1 {
		claimed, attempts, detail := bareEscClaimed245(t)
		if claimed {
			t.Skipf("SKIP-LOUD (ticket 245 AC#9, FOREIGN OWNER at baseline): an injected bare Esc reached no "+
				"foreground window while nothing of ours was running (keydown_esc=%d, wm_hotkey=%d), and Win32 "+
				"says the key is claimed: RegisterHotKey(foreign id) failed %d/3 times with %s. This desktop "+
				"already has an owner for the bare Esc, so delivery cannot be measured here; no product reading "+
				"was taken.", base.keydown, base.hotkey, attempts, detail)
		}
		t.Fatalf("ticket 245 AC#9 RED (instrument): an injected bare Esc reached no foreground window while "+
			"nothing of ours was running (keydown_esc=%d, wm_hotkey=%d), and Win32 says the bare Esc is FREE "+
			"(%d/3 attempts, no 1409). The observer or the injection half of this rig failed, so no reading "+
			"below could mean anything. Observer raw:\n%s", base.keydown, base.hotkey, attempts, base.raw)
	}

	// SHAPE 2 (E2/E4): the ruler's positive control. A process that registers the
	// bare Esc globally (this observer itself) must lose the keydown and receive
	// WM_HOTKEY instead. This is the same reading the pre-245 product produced
	// (E4's 0) and the proof that a 1 below is not the rig being deaf.
	steal := observeEsc246(t, rig, "-steal")
	if steal.keydown != 0 || steal.hotkey != 1 {
		t.Fatalf("ticket 245 AC#9 RED (instrument blind): with its OWN global bare-Esc registration the "+
			"observer still received keydown_esc=%d (want 0) and wm_hotkey=%d (want 1); the delivery ruler "+
			"cannot see a stolen key, so no product reading below means anything. Observer raw:\n%s",
			steal.keydown, steal.hotkey, steal.raw)
	}

	// SHAPE 3 (E3): the shipped no-args process with its ball up and nothing
	// waiting. The observer must still get the Esc.
	exe := buildWispForTest(t)
	leg := bootResidentLeg(t, exe, t.TempDir())
	if !pollUntil127(400, func() bool {
		out := leg.stdout.String()
		return strings.Contains(out, ballUpClaim) || strings.Contains(out, ballAbsentClaim)
	}) {
		t.Fatalf("ticket 245 AC#9 RED: the shipped resident leg never printed its ball posture.\n%s", leg.console())
	}
	out := leg.stdout.String()
	if strings.Contains(out, ballAbsentClaim) {
		// The "the ball must appear" verdict belongs to ticket 228's live case;
		// here the absence only means the key-ownership measurement has no subject.
		t.Skipf("SKIP-LOUD (ticket 245 AC#9): the shipped leg reported %q, so there is no idle ball whose "+
			"key ownership could be measured on this host.\n%s", ballAbsentClaim, out)
	}
	if i := strings.Index(out, "hotkeys live"); i >= 0 {
		frag := out[i:]
		if j := strings.IndexByte(frag, '\n'); j >= 0 {
			frag = frag[:j]
		}
		t.Logf("AC#6① live reading (fresh, this run): %s", frag)
	}
	idle := observeEsc246(t, rig, "-watch")
	if idle.keydown != 1 {
		// Ownership split: release the suspect, then ask Win32 again.
		_ = leg.breakToLoop()
		if !leg.exitedWithin(400) {
			leg.stop()
			t.Fatalf("ticket 245 AC#9 RED: the resident could not be released to attribute the Esc "+
				"ownership (exit %s).\n%s", leg.exitStatus(), leg.console())
		}
		claimed, attempts, detail := bareEscClaimed245(t)
		if claimed {
			t.Skipf("SKIP-LOUD (ticket 245 AC#9, FOREIGN OWNER): the shipped leg was up and idle, yet an "+
				"injected bare Esc reached the observer only %d time(s); after the leg left through its own "+
				"shutdown path Win32 STILL says the key is claimed (%d/3 attempts: %s), so the owner is not "+
				"this process and this desktop cannot judge delivery. The resident's own verdict:\n%s",
				idle.keydown, attempts, detail, out)
		}
		t.Fatalf("ticket 245 AC#9 RED (PRODUCT): the shipped no-args process held the desktop's bare Esc "+
			"while idle. An injected Esc reached the observer %d time(s) while the leg was up (want 1), and "+
			"the same probe found the key FREE again the moment the leg left - so the claim was this "+
			"process's, which is exactly the defect ticket 245 exists to keep out. Resident verdict:\n%s",
			idle.keydown, out)
	}

	// SHAPE 4: take the product away through its own shutdown path; the key must
	// be deliverable again. This is the "returned" half of E4 and the half a
	// bookkeeping flag cannot fake.
	_ = leg.breakToLoop()
	if !leg.exitedWithin(400) {
		leg.stop()
		t.Fatalf("ticket 245 AC#9 RED: the resident did not exit through its own shutdown path within the "+
			"budget (exit %s).\n%s", leg.exitStatus(), leg.console())
	}
	after := observeEsc246(t, rig, "-watch")
	if after.keydown != 1 {
		claimed, attempts, detail := bareEscClaimed245(t)
		if claimed {
			t.Skipf("SKIP-LOUD (ticket 245 AC#9, FOREIGN OWNER after exit): after the shipped leg left, an "+
				"injected bare Esc reached no foreground window (keydown_esc=%d), and Win32 says the key is "+
				"claimed again (%d/3 attempts: %s) - a foreign owner appeared, so the return half cannot be "+
				"judged on this desktop.", after.keydown, attempts, detail)
		}
		t.Fatalf("ticket 245 AC#9 RED: after the shipped leg left, an injected bare Esc reached no "+
			"foreground window (keydown_esc=%d, want 1) while Win32 says the bare Esc is FREE (%d/3 attempts, "+
			"no 1409). The key was not handed back to the desktop.\nObserver raw:\n%s",
			after.keydown, attempts, after.raw)
	}

	t.Logf("AC#9 four shapes (E1/E2/E3/E4): baseline=%d control_steal=%d/%d shipped_idle=%d after_exit=%d",
		base.keydown, steal.keydown, steal.hotkey, idle.keydown, after.keydown)
}
