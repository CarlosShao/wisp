//go:build windows && winlive

package ball

// Live hotkey tests (ticket 64 A1/A1b/A1d/A3). These run against the REAL
// desktop and the REAL RegisterHotKey, which is the only way to answer the
// question the old suite never asked: is the registration set what the config
// says it is?
//
//	go test -tags winlive ./internal/ball/ -run TestLiveHotkey -v

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/statemachine"
	"golang.org/x/sys/windows"
)

// spareHKID is an id the ball never uses, for squatting a combination in this
// very process to prove what a real 1409 looks like.
const spareHKID = 900

// liveHotkeys is the test binding set. It deliberately avoids the product
// defaults that are known-occupied on this machine (Ctrl+Alt+W, measured) and
// keeps a modifier combination for cancel: an idle ball no longer registers
// cancel at all (ticket 245 - the production default is a bare Esc, and
// RegisterHotKey would take Esc from every other app on the desktop for the
// whole test run). Since ticket 260 形ⓐ that choice does a second job: the
// borrow registers whatever [hotkey] cancel says, so this fixture is also the
// non-default case, and the live rulers below have to read "the key this seed
// names was borrowed" rather than the pre-260 "the borrowed key was Esc"
// (requireEscBorrowed; the expectation itself is pinned in the default suite by
// hotkey_cancel_borrow_expect_260r2_test.go). F11/F12
// are also occupied here (measured 2026-09-20: real 1409), which is why the
// previous ticket's F9-F12 fixture "passed" with nothing registered.
func liveHotkeys() HotkeyConfig {
	return HotkeyConfig{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Ctrl+Alt+V", Panel: "Ctrl+Alt+B"}
}

// hotkeyFreeProbe asks Win32 - not our own bookkeeping - whether ONE SPECIFIC
// combination is unclaimed on this desktop right now: register it on an id the
// ball never uses, then drop it immediately. ERROR_HOTKEY_ALREADY_REGISTERED
// means somebody holds it, and while a borrow is live that somebody is us - for
// whichever key the borrow registered, which since ticket 260 形ⓐ is the key
// [hotkey] cancel names, not a constant this file may spell out.
//
// This is the ruler behind ticket 245 AC#2's SECOND half - "还". The takeover
// flag alone only proves we wrote a bool; what the desktop can still use is a
// fact about Win32, and it is the same fact the defect was made of.
//
// The call runs on the window's owning thread (Win32 refuses it anywhere else)
// and unregisters again before returning, so a probe that says "free" leaves
// nothing bound behind.
//
// What it cannot tell apart - and could not for the bare Esc either - is "we hold
// it" from "another program holds it". Callers pair it with the ball's own
// registration set and say which of the two readings they got.
func hotkeyFreeProbe(t *testing.T, b *Ball, acc Accelerator) (free bool, err error) {
	t.Helper()
	res := make(chan error, 1)
	b.sta.PostTask(func() {
		r, _, e := pRegisterHotKey.Call(uintptr(b.hwnd), uintptr(spareHKID),
			uintptr(acc.Mods), uintptr(acc.VK))
		if r == 0 {
			if errno, ok := e.(syscall.Errno); ok && errno != 0 {
				res <- errno
				return
			}
			res <- syscall.EINVAL
			return
		}
		pUnregisterHotKey.Call(uintptr(b.hwnd), uintptr(spareHKID))
		res <- nil
	})
	if e := <-res; e != nil {
		if errors.Is(e, windows.ERROR_HOTKEY_ALREADY_REGISTERED) {
			return false, nil
		}
		return false, e
	}
	return true, nil
}

// escBorrowProbe is the bare-Esc instance of hotkeyFreeProbe. It stays exactly as
// ticket 245 wrote it because two questions are still about the bare Esc and
// nothing else: an idle ball must not hold the desktop's Esc whatever the user
// configured (requireIdleRoster), and a user who left [hotkey] cancel at the
// default must find the borrow on that same Esc (condition ①, 260-r1).
func escBorrowProbe(t *testing.T, b *Ball) (free bool, err error) {
	t.Helper()
	return hotkeyFreeProbe(t, b, Accelerator{Mods: modNoRepeat, VK: vkEscape})
}

// requireEscBorrowed is the non-skippable form of "Confirming really took the
// cancel key". The name and the bare-Esc probe are ticket 245's; WHICH key that
// is has come from [hotkey] cancel since ticket 260 形ⓐ, so the expectation is
// derived from the ball's own seed (wantCancelBorrow260r2, in the default-档 file
// hotkey_cancel_borrow_expect_260r2_test.go) instead of being spelled as a
// VK_ESCAPE constant in this file. That is a sync to the ruled behaviour, not a
// loosening: every branch below still ends in t.Fatal/t.Error, and the two cells
// stay two separate cells -
//
//   - seed unset or "Esc" (condition ①, zero drift): the expectation is the raw
//     default pair {MOD_NOREPEAT 0x4000, VK_ESCAPE 0x1B}, and the desktop half of
//     the check still runs against the bare Esc and still requires it NOT free;
//   - any other seed (AC#1): the expectation is that combination, the desktop half
//     probes THAT combination, and the bare Esc must now be free - the whole point
//     of 形ⓐ is that a user who moved off Esc does not also lose Esc.
//
// It accepts exactly two readings either way: the borrowed combination is in OUR
// registration set (then the desktop probe must see it as taken), or we tried and
// Win32 refused it (a foreign app holds that combination on this machine -
// reported, never passed off as a borrow). A ball that never attempted the borrow
// reads as standby and fails here, which is the whole point: the alternative would
// be a ruler that cannot tell "no borrow needed" from "borrow code deleted".
func requireEscBorrowed(t *testing.T, b *Ball) {
	t.Helper()
	seed := b.ConfiguredHotkeys().Cancel
	wish := wantCancelBorrow260r2(seed)
	rep := b.HotkeyReport()
	if len(rep.Live()) != 4 {
		line, _ := rep.Binding(hkCancel)
		if line.Status == HotkeyStandby {
			t.Fatalf("ticket 245 RED: Confirming never attempted the cancel borrow (cancel line still "+
				"standby, live set holds %d of the 4 slots, [hotkey] cancel = %q): %+v",
				len(rep.Live()), seed, rep.Bindings())
		}
		if line.Status == HotkeyTaken || line.Status == HotkeyError {
			t.Logf("SKIP-LOUD: the cancel borrow was ATTEMPTED and refused by Win32 (%v) - another program on "+
				"this desktop owns %+v, so this run did not prove the borrowed key reaches the shell. "+
				"It did prove the ball tries (that is what ticket 245's borrow path is).", line.Err, wish.Acc)
			return
		}
		t.Fatalf("after the borrow the live set holds %d entries, want 4: %+v", len(rep.Live()), rep.Bindings())
	}
	if got := rep.Live()[hkCancel]; got != wish.Acc {
		switch {
		case wish.Refused != nil:
			t.Fatalf("the borrowed cancel slot holds %+v, want the default pair %+v: [hotkey] cancel %q is not a "+
				"parsable combination (%v), and condition ③ says the borrow takes the default and SAYS so: %+v",
				got, wish.Acc, seed, wish.Refused, rep.Bindings())
		case wish.Default:
			t.Fatalf("ticket 260 condition ① RED: [hotkey] cancel is the default (%q), so the borrowed slot must "+
				"hold the raw pair {mods 0x4000, VK 0x1B}, got %+v: %+v", seed, got, rep.Bindings())
		default:
			t.Fatalf("the borrowed cancel slot holds %+v, want the configured cancel key %q = %+v: after ticket 260 "+
				"形ⓐ the borrow registers what the config names, so this is either the pre-260 hard-coded Esc or a "+
				"key the config never named: %+v", got, seed, wish.Acc, rep.Bindings())
		}
	}
	// The desktop, not our own bookkeeping, says the borrowed combination is ours.
	if free, err := hotkeyFreeProbe(t, b, wish.Acc); err != nil {
		t.Fatalf("the probe for the borrowed %q (%+v) could not run: %v", seed, wish.Acc, err)
	} else if free {
		t.Fatalf("ticket 245 RED: the ball says it borrowed %q (%+v), but the desktop still has that combination "+
			"free - the borrow registered on some other key: %+v", seed, wish.Acc, rep.Bindings())
	}
	if !wish.Default {
		// The second cell of 形ⓐ, on the same live ball: the key nobody configured
		// must not be the one that got taken.
		if free, err := escBorrowProbe(t, b); err != nil {
			t.Fatalf("the bare-Esc probe could not run: %v", err)
		} else if !free {
			t.Errorf("bare Esc is claimed while [hotkey] cancel says %q and the borrow holds %+v - the ball is "+
				"holding two keys: %+v", seed, wish.Acc, rep.Bindings())
		}
	}
}

// TestLiveHotkeyRebindEndToEnd is the A1 proof on a running instance: a real
// config.Manager over a real config.toml, the same bridge the app installs,
// a hand on the file, and the Win32 registration set checked before and after.
func TestLiveHotkeyRebindEndToEnd(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	writeHotkeyConfig(t, path, "Ctrl+Alt+Q")
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("config.NewManager(%s): %v", path, err)
	}

	var mu sync.Mutex
	fired := map[string]int{}
	bootCfg := ApplyHotkeyDefaults(hotkeysFromConfig(mgr))
	b := newLiveBall(t, Options{
		Initial: statemachine.StateSleeping,
		Hotkeys: bootCfg,
		Events: Events{
			OnSummonHotkey: func() { bump(&mu, fired, "summon") },
			OnMuteHotkey:   func() { bump(&mu, fired, "mute") },
		},
	})

	// 1. The registration set at boot is asserted, not assumed (registry A1:
	// "四项默认热键全注册失败、测试仍报 PASS"). Since ticket 245 the idle set is
	// THREE, and the missing one is the cancel slot by design.
	boot := b.HotkeyReport()
	if !boot.AllLive() {
		t.Fatalf("boot registration incomplete: %+v", boot.Bindings())
	}
	requireIdleRoster(t, b)
	if got := boot.Live()[hkSummon].VK; got != 'Q' {
		t.Fatalf("summon registered as VK 0x%X, want 'Q' (0x51)", got)
	}
	if p := boot.Problems(); len(p) != 0 {
		t.Fatalf("problems on a clean boot: %v", p)
	}

	// 2. Install the bridge exactly as the host does: the OnReload hook for the
	//    push path, plus the Refresh hook that drives config.Manager poll.
	r := NewHotkeyReloader(b, b.ConfiguredHotkeys(), func() HotkeyConfig {
		return ApplyHotkeyDefaults(hotkeysFromConfig(mgr))
	})
	r.Refresh = func() error { _, err := mgr.CheckAndReload(); return err }
	mgr.OnReload = r.OnReload()

	// 3. Edit the file the way a user does, then run the host's poll tick
	//    (watchdog tick = CheckAndReload + Check). No new goroutine here: an
	//    unrostered resident goroutine is a D38b leak symptom, and observe
	//    says so - see the ticket's note for ticket 12 about where Run() may
	//    legally live.
	writeHotkeyConfig(t, path, "Ctrl+Alt+R")
	seen := waitFor(3*time.Second, func() bool {
		r.Check()
		live := b.RegisteredHotkeys()
		return len(live) == 3 && live[hkSummon].VK == 'R'
	})
	if !seen {
		live := b.RegisteredHotkeys()
		t.Fatalf("editing [hotkey] summon to Ctrl+Alt+R never reached Win32: set=%s rebinds=%d applied=%+v",
			accelSetString(live), r.Rebinds(), r.Applied())
	}
	if r.Rebinds() != 1 {
		t.Fatalf("the bridge rebound %d times for one edit, want exactly 1", r.Rebinds())
	}
	after := b.HotkeyReport()
	if !after.IsLive(hkSummon) || after.Live()[hkSummon].VK != 'R' {
		t.Fatalf("post-rebind report wrong: %+v", after.Bindings())
	}
	if b.ConfiguredHotkeys().Summon != "Ctrl+Alt+R" {
		t.Errorf("ConfiguredHotkeys = %q, want Ctrl+Alt+R", b.ConfiguredHotkeys().Summon)
	}

	// 4. The window procedure routes the live id to the consumer callback.
	sendMessage(b, wmHotkey, hkSummon, 0)
	if count(&mu, fired, "summon") != 1 {
		t.Fatalf("WM_HOTKEY(summon) did not reach OnSummonHotkey: %v", fired)
	}

	// 5. The physical press: inject the NEW combination (must fire) and the OLD
	//    one (must not). Injected input travels through the shell, so this is
	//    the real hotkey dispatch and not a posted message.
	if !injectBinding(t, "Ctrl+Alt+R") {
		t.Skipf("SKIP-LOUD: keybd_event injection never reached the input stream, so this run did NOT prove " +
			"the shell->WM_HOTKEY dispatch for the rebound key. Steps 1-4 prove the config edit changed the " +
			"Win32 registration set and that the window procedure routes it; no assertion in this file was skipped.")
	}
	if !waitFor(2*time.Second, func() bool { return count(&mu, fired, "summon") >= 2 }) {
		t.Fatalf("the newly registered Ctrl+Alt+R never summoned the ball (fired=%v)", fired)
	}
	if got := count(&mu, fired, "summon"); got != 2 {
		t.Fatalf("one physical press produced %d summons in total, want exactly 2 (1 posted + 1 injected)", got)
	}
	if !injectBinding(t, "Ctrl+Alt+Q") {
		t.Skipf("SKIP-LOUD: injection refused on the second press, so this run did not prove the retired " +
			"key is dead. The registration set above shows VK 'Q' is no longer ours.")
	}
	time.Sleep(300 * time.Millisecond)
	if got := count(&mu, fired, "summon"); got != 2 {
		t.Fatalf("the OLD Ctrl+Alt+Q still summons after the rebind (fired=%d, want the 2 from before)", got)
	}
}

// TestLiveHotkeyOccupiedVsNotAttempted splits the two failure families against
// the REAL API: a combination this process squats on a foreign id must come
// back as occupied (1409), an unset binding must come back as never attempted.
func TestLiveHotkeyOccupiedVsNotAttempted(t *testing.T) {
	b := newLiveBall(t, Options{Initial: statemachine.StateSleeping, Hotkeys: liveHotkeys()})

	acc, err := ParseAccelerator("Ctrl+Alt+U")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	// Squat: same window, different id - the shape of "somebody else already
	// owns this combination" as far as RegisterHotKey is concerned. Register
	// and unregister on the OWNING thread: from anywhere else Win32 refuses
	// with "Invalid window; it belongs to other thread".
	if e := uiRegisterHotkey(t, b, spareHKID, acc); e != nil {
		if errors.Is(e, windows.ERROR_HOTKEY_ALREADY_REGISTERED) {
			t.Skipf("SKIP-LOUD: Ctrl+Alt+U is already owned by another program on this machine, so the "+
				"in-process squat cannot be built and this run did NOT exercise the real occupied path (%v). "+
				"The classification itself is pinned by TestRegisterAllSplitsFailureFamilies.", e)
		}
		t.Fatalf("squat RegisterHotKey failed: %v", e)
	}

	cfg := liveHotkeys()
	cfg.Summon = "Ctrl+Alt+U" // squatted on this very window
	rep := b.RebindHotkeys(cfg)
	sum, ok := rep.Binding(hkSummon)
	if !ok {
		t.Fatalf("no summon line in %+v", rep.Bindings())
	}
	if sum.Status != HotkeyTaken {
		t.Fatalf("squatted binding reported %q, want occupied-by-another-app", sum.Status)
	}
	if !errors.Is(sum.Err, windows.ERROR_HOTKEY_ALREADY_REGISTERED) {
		t.Fatalf("occupied verdict carries err %v, want ERROR_HOTKEY_ALREADY_REGISTERED (1409)", sum.Err)
	}
	if !sum.Status.Attempted() {
		t.Error("HotkeyTaken must be an attempted status")
	}
	probs := strings.Join(rep.Problems(), "\n")
	if !strings.Contains(probs, "occupied by another program") || !strings.Contains(probs, "summon") {
		t.Fatalf("the occupied case produced no user-visible line: %q", probs)
	}
	if !errors.Is(sum.Err, windows.Errno(1409)) {
		t.Error("the occupied verdict must be errno 1409")
	}

	// The other slots survive, and the unset / junk ones are never attempted.
	cfg.Mute = ""
	cfg.Panel = "Ctrl+Alt+NotAKey+"
	rep2 := b.RebindHotkeys(cfg)
	// Cancel is the slot this test used to assert as "still live after its
	// neighbours broke". Ticket 245 moved that property to the borrow: an idle
	// ball must NOT hold cancel, so what has to survive a broken neighbour set is
	// the cancel path's ABILITY to take Esc when Confirming asks for it.
	if c, _ := rep2.Binding(hkCancel); c.Status != HotkeyStandby {
		t.Errorf("cancel line after siblings failed = %+v, want standby (present in the report, "+
			"never bound while idle)", c)
	}
	if rep2.IsLive(hkCancel) {
		t.Error("ticket 245 RED: the rebalanced idle set still holds the cancel slot")
	}
	if mt, _ := rep2.Binding(hkMute); mt.Status != HotkeyDisabled || mt.Status.Attempted() {
		t.Errorf("unset mute = %q, want disabled/never-attempted", mt.Status)
	}
	if pn, _ := rep2.Binding(hkPanel); pn.Status != HotkeyUnparsable || pn.Status.Attempted() {
		t.Errorf("junk panel = %q, want unparsable/never-attempted", pn.Status)
	}
	if rep2.AllLive() {
		t.Error("AllLive true while three slots are not live")
	}

	// Free the squat and the same key goes live again: rebind unregisters all
	// four slots first, so the occupied verdict cannot be our own leftover.
	uiUnregisterHotkey(t, b, spareHKID)
	rep3 := b.RebindHotkeys(liveHotkeys())
	if !rep3.AllLive() {
		t.Fatalf("after freeing the key the set is still incomplete: %+v", rep3.Bindings())
	}
	if got := rep3.Live()[hkSummon].VK; got != 'Q' {
		t.Fatalf("summon VK = 0x%X, want 'Q'", got)
	}
	// The cancel path's ability, asserted once the spare id is really free (the
	// probes inside these helpers register the borrowed combination - for this
	// fixture [hotkey] cancel's Ctrl+Alt+V, plus the bare Esc that the idle
	// roster insists is free - on that id, so nothing may still be squatted).
	// This is what survived from ticket 64's "cancel is not dropped when its
	// neighbours fail": the slot is idle-unbound by design, and the borrow that
	// binds it has to work anyway - then give that same key back.
	requireIdleRoster(t, b)
	b.TakeEscForCancel()
	requireEscBorrowed(t, b)
	b.ReleaseEscAfterSession()
	requireEscReturned(t, b)
}

// TestLiveMuteHotkeyEndToEnd is A1d: the registered mute key -> OnMuteHotkey ->
// EvMuteKey -> the machine and the ball both in Muted, and back out.
func TestLiveMuteHotkeyEndToEnd(t *testing.T) {
	m := statemachine.New(statemachine.Options{Initial: statemachine.StateSleeping})
	var mu sync.Mutex
	hits := map[string]int{}
	var b *Ball
	// The consumer law balldebug implements, in miniature: the hotkey callback
	// dispatches the machine event and mirrors the resulting state onto the
	// ball (the same SetState the app does after every dispatch).
	onMute := func() {
		h := b.sta
		h.PostTask(func() { bump(&mu, hits, "mute-cb") })
		if _, err := m.Dispatch(statemachine.EvMuteKey, &statemachine.Facts{KwsLoaded: true}); err != nil {
			h.PostTask(func() { t.Errorf("dispatch EvMuteKey: %v", err) })
			return
		}
		b.SetState(m.State())
		h.PostTask(func() { bump(&mu, hits, "state:"+string(m.State())) })
	}
	b = newLiveBall(t, Options{
		Initial: statemachine.StateSleeping,
		Hotkeys: liveHotkeys(),
		Events:  Events{OnMuteHotkey: onMute},
	})
	if !b.HotkeyReport().IsLive(hkMute) {
		t.Fatalf("the mute hotkey is not registered: %+v", b.HotkeyReport().Bindings())
	}
	// D43 #5: Sleeping -> Armed; then #8: Armed + EvMuteKey -> Muted.
	if _, err := m.Dispatch(statemachine.EvKwsEnabled, nil); err != nil {
		t.Fatalf("dispatch EvKwsEnabled: %v", err)
	}
	if m.State() != statemachine.StateArmed {
		t.Fatalf("machine in %s, want Armed before the mute key", m.State())
	}

	sendMessage(b, wmHotkey, hkMute, 0)
	if !waitFor(2*time.Second, func() bool { return count(&mu, hits, "mute-cb") == 1 }) {
		t.Fatalf("WM_HOTKEY(mute) never reached OnMuteHotkey: %v", hits)
	}
	if got := m.State(); got != statemachine.StateMuted {
		t.Fatalf("machine after the mute hotkey = %s, want Muted (D43 #8)", got)
	}
	if got := readState(t, b); got != statemachine.StateMuted {
		t.Fatalf("the ball renders %s, want Muted", got)
	}
	// Un-mute with the same key (D43 #10, KWS loaded -> Armed).
	if _, err := m.Dispatch(statemachine.EvMuteKey, &statemachine.Facts{KwsLoaded: true}); err != nil {
		t.Fatalf("second EvMuteKey: %v", err)
	}
	if m.State() != statemachine.StateArmed {
		t.Fatalf("machine after the second mute event = %s, want Armed (D43 #10)", m.State())
	}
	// And the physical press, when the session accepts injected input.
	if !injectBinding(t, liveHotkeys().Mute) {
		t.Skipf("SKIP-LOUD: keybd_event injection was refused, so this run did not prove the shell dispatch " +
			"for the mute key. Registration of hkMute, the window-procedure routing, the EvMuteKey transition " +
			"and the rendered state were all asserted above; the un-mute direction was driven directly.")
	}
	if !waitFor(2*time.Second, func() bool { return count(&mu, hits, "mute-cb") == 2 }) {
		t.Fatalf("the registered Ctrl+Alt+M never reached OnMuteHotkey (fired=%v)", hits)
	}
	if got := m.State(); got != statemachine.StateMuted {
		t.Fatalf("physical mute key left the machine in %s, want Muted", got)
	}
}

// TestLiveSleepingZeroTimerHandles is A3: Sleeping holds no live timer,
// measured from the messages Win32 really delivers to the window. The same
// probe must SEE a timer in Warm, or the zero in Sleeping proves nothing.
func TestLiveSleepingZeroTimerHandles(t *testing.T) {
	b := newLiveBall(t, Options{Initial: statemachine.StateSleeping, Hotkeys: liveHotkeys()})
	cnt, restore := subclassTimers(t, b)
	defer restore()

	b.SetState(statemachine.StateSleeping)
	time.Sleep(600 * time.Millisecond)
	if n := cnt.wmTimer.Load(); n != 0 {
		t.Fatalf("Sleeping received %d WM_TIMER messages in 600ms (zero-timer discipline broken)", n)
	}
	if b.DebugTimersAlive() {
		t.Error("DebugTimersAlive true in Sleeping")
	}
	if n := cnt.total.Load(); n == 0 {
		t.Fatal("the probe saw zero messages at all: the subclass is not installed, so the " +
			"Sleeping verdict above proved nothing")
	}
	if th := ballWindowsOnThread(b.DebugHWND()); th != 1 {
		t.Errorf("the UI thread owns %d ball windows in Sleeping, want exactly 1 (all classes: %v)",
			th, threadWindowClasses(b.DebugHWND()))
	}

	// Warm breathes: the probe must be able to see a real timer, and then see
	// it die on the way back to Sleeping.
	b.SetState(statemachine.StateWarm)
	if !waitFor(2*time.Second, func() bool { return cnt.wmTimer.Load() > 0 }) {
		t.Fatalf("Warm armed no timer the window ever received (saw %d messages total)", cnt.total.Load())
	}
	b.SetState(statemachine.StateSleeping)
	time.Sleep(200 * time.Millisecond) // let any in-flight tick land
	before := cnt.wmTimer.Load()
	time.Sleep(600 * time.Millisecond)
	if after := cnt.wmTimer.Load(); after != before {
		t.Fatalf("%d WM_TIMER messages still arriving 600ms after the return to Sleeping", after-before)
	}
}

// requireIdleRoster is AC#1 on a real window: the idle registration set is
// EXACTLY summon / mute / panel, the cancel slot is not in it, and - because a
// count alone could be three by accident - the bare Esc must be free on the
// desktop at that moment. The two halves are the same ruler run twice: the
// pre-245 ball read 4 live with Esc bound, and a ball that simply failed to
// register anything would read 0 or 1.
func requireIdleRoster(t *testing.T, b *Ball) {
	t.Helper()
	rep := b.HotkeyReport()
	if _, held := rep.Live()[hkCancel]; held {
		t.Fatalf("ticket 245 RED: the idle ball holds the cancel hot key: %+v", rep.Bindings())
	}
	if len(rep.Live()) != 3 {
		t.Fatalf("idle live set = %d entries, want exactly 3 (summon/mute/panel): %+v",
			len(rep.Live()), rep.Bindings())
	}
	if !rep.IsLive(hkSummon) || !rep.IsLive(hkMute) || !rep.IsLive(hkPanel) {
		t.Fatalf("idle live set is not summon/mute/panel: %+v", rep.Bindings())
	}
	if c, _ := rep.Binding(hkCancel); c.Status != HotkeyStandby {
		t.Errorf("idle cancel line = %+v, want standby", c)
	}
	free, err := escBorrowProbe(t, b)
	if err != nil {
		t.Fatalf("the bare-Esc probe could not run: %v", err)
	}
	if !free {
		t.Fatal("ticket 245 RED: bare Esc is claimed on this desktop while the ball is idle")
	}
}

// requireEscReturned is AC#2's "还" half on a real window: the roster is idle
// again AND Win32 itself reports the key that was borrowed unclaimed.
//
// Before ticket 260 those were the same probe - the borrow was always the bare
// Esc, so requireIdleRoster's bare-Esc reading doubled as the return check. Once
// the borrow follows [hotkey] cancel (形ⓐ), a non-default seed needs its OWN
// free-reading or a release that never unregistered would be invisible, which is
// precisely the zero-symptom shape this ticket family exists to keep loud. For
// the default seed the borrowed key IS the bare Esc, so this adds no second
// probe and the default档 reading stays bit-for-bit what ticket 245 shipped.
func requireEscReturned(t *testing.T, b *Ball) {
	t.Helper()
	seed := b.ConfiguredHotkeys().Cancel
	wish := wantCancelBorrow260r2(seed)
	if b.EscTakenOver() {
		t.Fatal("ticket 245 RED: the cancel key is still marked taken over after the session end")
	}
	requireIdleRoster(t, b)
	if wish.Default {
		return // requireIdleRoster just proved the borrowed key - the bare Esc - is free
	}
	free, err := hotkeyFreeProbe(t, b, wish.Acc)
	if err != nil {
		t.Fatalf("the probe for the key the borrow held (%q / %+v) could not run: %v", seed, wish.Acc, err)
	}
	if free {
		return
	}
	// Still claimed. If our own set names that id, the return never happened; if
	// the set is clean, either the UnregisterHotKey did not reach Win32 or
	// something else took the key mid-run - and neither reading is a pass.
	if _, ours := b.RegisteredHotkeys()[hkCancel]; ours {
		t.Fatalf("ticket 245/260 RED: the borrowed %q (%+v) is still in the ball's own registration set after the "+
			"session end: %s", seed, wish.Acc, accelSetString(b.RegisteredHotkeys()))
	}
	t.Errorf("the borrowed %q (%+v) is still claimed on this desktop after the session end while the ball's own "+
		"set is clean: the return never reached Win32, or the key was taken again mid-run - this run cannot tell "+
		"those apart and neither one is a pass", seed, wish.Acc)
}

// ------------------------------------------------------------------ helpers

// uiRegisterHotkey registers a combination for the ball window ON the thread
// that owns that window (Win32 refuses the call from anywhere else).
func uiRegisterHotkey(t *testing.T, b *Ball, id uint32, acc Accelerator) error {
	t.Helper()
	res := make(chan error, 1)
	b.sta.PostTask(func() {
		r, _, e := pRegisterHotKey.Call(uintptr(b.hwnd), uintptr(id), uintptr(acc.Mods), uintptr(acc.VK))
		if r == 0 {
			if errno, ok := e.(syscall.Errno); ok && errno != 0 {
				res <- errno
				return
			}
			res <- syscall.EINVAL
			return
		}
		res <- nil
	})
	return <-res
}

// uiUnregisterHotkey drops one id on the owning thread.
func uiUnregisterHotkey(t *testing.T, b *Ball, id uint32) {
	t.Helper()
	done := make(chan struct{})
	b.sta.PostTask(func() {
		pUnregisterHotKey.Call(uintptr(b.hwnd), uintptr(id))
		close(done)
	})
	<-done
}

// hotkeysFromConfig maps [hotkey] the way a host does (the adapter ticket 12
// installs next to NewHotkeyReloader).
func hotkeysFromConfig(m *config.Manager) HotkeyConfig {
	h := m.Config().Hotkey
	return HotkeyConfig{Summon: h.Summon, Mute: h.Mute, Cancel: h.Cancel, Panel: h.Panel}
}

// writeHotkeyConfig writes a loadable config.toml whose [hotkey] summon
// binding is `summon`.
func writeHotkeyConfig(t *testing.T, path, summon string) {
	t.Helper()
	c := config.NewDefaults()
	c.Hotkey = config.HotkeySection{
		Summon: summon,
		Mute:   "Ctrl+Alt+M",
		Cancel: "Ctrl+Alt+V",
		Panel:  "Ctrl+Alt+B",
	}
	if err := config.SaveFile(path, c); err != nil {
		t.Fatalf("SaveFile(%s): %v", path, err)
	}
	// Nudge the mtime forward: two writes inside the same clock tick would be
	// polled as "unchanged" by the manager's mtime+size check.
	now := time.Now()
	if err := os.Chtimes(path, now.Add(time.Second), now.Add(time.Second)); err != nil {
		t.Fatalf("chtimes(%s): %v", path, err)
	}
	if _, _, err := config.LoadFile(path, nil); err != nil {
		t.Fatalf("the written config.toml does not load: %v", err)
	}
}

func bump(mu *sync.Mutex, m map[string]int, k string) {
	mu.Lock()
	defer mu.Unlock()
	m[k]++
}

func count(mu *sync.Mutex, m map[string]int, k string) int {
	mu.Lock()
	defer mu.Unlock()
	return m[k]
}

// waitFor polls until the condition holds or the budget is spent. The budget
// is a timer, never a wall-clock difference used as a verdict.
func waitFor(budget time.Duration, cond func() bool) bool {
	deadline := time.NewTimer(budget)
	defer deadline.Stop()
	for {
		if cond() {
			return true
		}
		select {
		case <-deadline.C:
			return cond()
		default:
			time.Sleep(20 * time.Millisecond)
		}
	}
}

func accelSetString(m map[uint32]Accelerator) string {
	ks := make([]string, 0, len(m))
	for id, a := range m {
		ks = append(ks, fmt.Sprintf("%d:%04x/%02x", id, a.Mods, a.VK))
	}
	return strings.Join(ks, " ")
}

// injectBinding presses a configured combination through the keyboard driver
// and reports whether the injection reached the input stream.
func injectBinding(t *testing.T, binding string) bool {
	t.Helper()
	acc, err := ParseAccelerator(binding)
	if err != nil {
		t.Fatalf("ParseAccelerator(%q): %v", binding, err)
	}
	return injectHotkey(acc.Mods&^modNoRepeat, acc.VK)
}
