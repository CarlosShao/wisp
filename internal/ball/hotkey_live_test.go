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
// the bare "Esc" cancel default: registering Esc as a GLOBAL hotkey swallows
// Esc from every other app on the desktop for the whole test run, so the live
// tests bind a modifier combination instead (see the ticket's found-defect
// note). F11/F12 are also occupied here (measured 2026-09-20: real 1409),
// which is why the previous ticket's F9-F12 fixture "passed" with nothing
// registered.
func liveHotkeys() HotkeyConfig {
	return HotkeyConfig{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Ctrl+Alt+V", Panel: "Ctrl+Alt+B"}
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
	// "四项默认热键全注册失败、测试仍报 PASS").
	boot := b.HotkeyReport()
	if !boot.AllLive() {
		t.Fatalf("boot registration incomplete: %+v", boot.Bindings())
	}
	if len(boot.Live()) != 4 {
		t.Fatalf("live set = %d entries, want 4: %+v", len(boot.Live()), boot.Bindings())
	}
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
		return len(live) == 4 && live[hkSummon].VK == 'R'
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
	if !rep2.IsLive(hkCancel) {
		t.Errorf("cancel dropped after siblings failed: %+v", rep2.Bindings())
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
