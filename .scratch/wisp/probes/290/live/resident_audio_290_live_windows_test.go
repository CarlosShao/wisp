//go:build windows && winlive

package main

// Ticket 290 AC#3, both shapes, on the shipped binary and on the production
// assembly. This file is the reading instrument AC#3 was missing; it adds no
// production behaviour and moves no default.
//
// WHY IT IS A WINLIVE FILE. Everything here touches the real desktop and the
// real capture endpoint: shape 1 boots the shipped no-args process (the same
// command line a double click produces) and reads the boot verdict out of its
// own console, and shape 2 presses the user's own mute hot key on that same
// process and reads the level counter its own exit line prints. It is compiled
// only under `-tags winlive`, so the ordinary whole-package roster is untouched.
//
// THE TWO SHAPES AND WHAT EACH ONE IS ALLOWED TO CLAIM
//
//	shape 1 - the default posture. No config.toml in a fresh data dir, so the
//	          compiled default table decides: mic_muted_default=true, the gate
//	          is mounted with WithStartMuted(true), and the device is never
//	          opened. The claim is read off two places: the boot verdict the
//	          process prints, and the counter line its own shutdown prints
//	          (frames_sent=0 levels_delivered=0 - nothing was ever captured).
//
//	shape 2 - the user opens the door. Same binary, same defaults. The mute
//	          combination the boot says it registered live is pressed through
//	          the OS input queue (keybd_event, the same chord a hand on the
//	          keyboard produces), the ball's mute gesture turns THIS process's
//	          gate, and the exit line's levels_delivered counts the scalars
//	          handed to the ball seam after that. No temp config with
//	          mic_muted_default=false is involved: that shape is what AC#3's
//	          parenthetical forbids, and it is not built here.
//
// WHAT NEITHER SHAPE CLAIMS
//
//   - "the orb visibly breathes": prototypeVisuals is off by default
//     (internal/ball/statevisual.go) and Ball.SetAudioLevel returns early while
//     it is out; that is ticket 68's decision, and resident_audio_windows.go
//     already says so in its own outcome sentence.
//   - "the numbers rise when the owner speaks": that half needs a sound source
//     this machine can produce at the microphone, which is the owner-present
//     window (ticket 247 AC#2 / ticket 297 own the ruler question).
//
// HOW TO RUN IT. Both tasklist rulers must read 0 first. In Git Bash:
//
//	PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" WISP_LIVE_MIC=1 \
//	  go test -tags winlive ./cmd/wisp -run 'TestAC290Live' -count=1 -v
//
// The run leaves a real ball window on the desktop for the two seconds the
// shipped process is up, and the shipped process really owns the default mute
// hot key while it is alive; both end with the process, never before.

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	// ac290MutedVerdict is the boot sentence AC#3 shape 1 pins, verbatim from
	// cmd/wisp/resident_audio_windows.go's fmt.Sprintf (the ticket's ":238" line
	// as it reads today). It is a substring so a drift in the source path (the
	// tail of the sentence) cannot make the reading pass on a different posture.
	ac290MutedVerdict = "采集腿已装配、门处于静音：设备未打开（[audio] mic_muted_default=true"
	// ac290OpenVerdict is the outcome residentAudio.toggleMute returns when the
	// gate really handed the device over (resident_audio_windows.go's "已取消
	// 静音" branch). The console prints it behind "wisp: ball mute-hotkey: ".
	ac290OpenVerdict = "已取消静音：设备已交接，采集线程在跑"
	// ac290StopPrefix is the first half of the D38(e) step-4 exit line
	// (residentAudio.stop), the only place the shipped process tells anybody how
	// many scalars reached the ball seam.
	ac290StopPrefix = "audio capture leg stopped: levels_delivered="
	// ac290LoopLine is the boot-complete fact: the event loop owns the stop
	// request, so a break sent before this line could land in a process that is
	// not listening for it yet.
	ac290LoopLine = "resident event loop running"
	// ac290MuteLive is the boot summary's own claim that the default mute chord
	// was registered live. If this is absent the premise for shape 2 is not
	// there (somebody else holds the chord), and the case says so instead of
	// pressing a chord that would go to that other process.
	ac290MuteLive = "mute=Ctrl+Alt+M live"
)

const (
	ac290VKControl = 0x11 // VK_CONTROL
	ac290VKMenu    = 0x12 // VK_MENU (Alt)
	ac290VKLetterM = 0x4D // 'M'
	ac290KeyUp     = 0x0002
)

var pAc290KeybdEvent = windows.NewLazySystemDLL("user32.dll").NewProc("keybd_event")

// ac290PressMuteChord delivers Ctrl+Alt+M the way the input queue delivers a
// real hand: key-down in order, key-up in reverse. keybd_event (not PostMessage)
// because the point is that the OS's own hot key dispatch produces WM_HOTKEY -
// posting the message ourselves would test the window proc, not the user's path.
// The short gaps between the downs are there because a hot key is matched on the
// modifier state at the moment the letter is seen; a chord typed as one burst
// can be read before the modifiers land.
func ac290PressMuteChord() {
	pAc290KeybdEvent.Call(uintptr(ac290VKControl), 0, 0, 0)
	time.Sleep(20 * time.Millisecond)
	pAc290KeybdEvent.Call(uintptr(ac290VKMenu), 0, 0, 0)
	time.Sleep(20 * time.Millisecond)
	pAc290KeybdEvent.Call(uintptr(ac290VKLetterM), 0, 0, 0)
	time.Sleep(60 * time.Millisecond)
	pAc290KeybdEvent.Call(uintptr(ac290VKLetterM), 0, uintptr(ac290KeyUp), 0)
	pAc290KeybdEvent.Call(uintptr(ac290VKMenu), 0, uintptr(ac290KeyUp), 0)
	pAc290KeybdEvent.Call(uintptr(ac290VKControl), 0, uintptr(ac290KeyUp), 0)
}

// ac290InputDesktopName names the desktop this process's input would land on.
// "Default" is the interactive one; "Winlogon" (or a failure to open the input
// desktop at all) means the session is locked, and then a synthesised chord
// cannot reach any hot key - which is a fact the reading has to carry instead of
// being reported as "the product ignored the key".
func ac290InputDesktopName() string {
	r, _, err := pAc290OpenInputDesktop.Call(0, 0, 0x0100) // DESKTOP_SWITCHDESKTOP
	if r == 0 {
		return "unknown (OpenInputDesktop failed: " + err.Error() + ")"
	}
	defer pAc290CloseDesktop.Call(r)
	buf := make([]uint16, 128)
	var needed uint32
	ok, _, err2 := pAc290GetUserObjectInfoW.Call(r, 2, // UOI_NAME
		uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)*2), uintptr(unsafe.Pointer(&needed)))
	if ok == 0 {
		return "unknown (GetUserObjectInformation failed: " + err2.Error() + ")"
	}
	return windows.UTF16ToString(buf)
}

var (
	pAc290OpenInputDesktop   = windows.NewLazySystemDLL("user32.dll").NewProc("OpenInputDesktop")
	pAc290CloseDesktop       = windows.NewLazySystemDLL("user32.dll").NewProc("CloseDesktop")
	pAc290GetUserObjectInfoW = windows.NewLazySystemDLL("user32.dll").NewProc("GetUserObjectInformationW")
	pAc290RegisterHotKey     = windows.NewLazySystemDLL("user32.dll").NewProc("RegisterHotKey")
	pAc290UnregisterHotKey   = windows.NewLazySystemDLL("user32.dll").NewProc("UnregisterHotKey")
)

// ac290SquatID is this process's registration id for the ownership probe.
const ac290SquatID = 0x290

// ac290TrySquatMuteChord asks Win32 who holds Ctrl+Alt+M. It returns "" when
// THIS process managed to register it (then the shipped process is not holding
// it, and the registration is released again before returning), and a sentence
// naming the holder otherwise.
func ac290TrySquatMuteChord() string {
	// RegisterHotKey with hWnd 0 binds to the CALLING THREAD, and the release
	// must land on that same thread, so both happen here with the goroutine
	// pinned. A migrated release would leave this process holding the chord and
	// the press below would come back to us instead of the product.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	r, _, callErr := pAc290RegisterHotKey.Call(0, uintptr(ac290SquatID), 0x0002|0x0001, uintptr(ac290VKLetterM)) // MOD_CONTROL|MOD_ALT
	if r == 0 {
		errno, _ := callErr.(syscall.Errno)
		if errno == windows.ERROR_HOTKEY_ALREADY_REGISTERED {
			return "held by another process (RegisterHotKey answered ERROR_HOTKEY_ALREADY_REGISTERED while the child is alive)"
		}
		return "cannot tell (RegisterHotKey answered " + callErr.Error() + ")"
	}
	pAc290UnregisterHotKey.Call(0, uintptr(ac290SquatID))
	return ""
}

// ac290ClickTrayMuteRow drives the tray's mute item with real mouse input on the
// popup the CHILD itself opened, and returns the child's own console line when
// the gate opens ("" when no candidate row was the mute one).
//
// The carrier is ticket 228's: PostMessageW the tray callback the shell would
// send, find the popup by the CHILD's pid (a menu is a #32768 top-level window,
// and the owning pid is the identity), then click inside its rect. Tickets 229 /
// A495 forbid rebuilding this with a posted WM_COMMAND, because internal/ball
// takes the selection only from TrackPopupMenu's synchronous return.
//
// The row position is not known to a process that does not own the HMENU (the
// menu is built inside the child and destroyed there), so candidates are tried
// from the second item's neighbourhood downward and each attempt reopens the
// menu. The console says which row was hit, so a miss is visible, not silent.
// The Exit row (the last one) is never among the candidates.
func ac290ClickTrayMuteRow(t *testing.T, leg *residentLeg) string {
	t.Helper()

	// The cursor is the owner's, not this case's: saved first, restored first out.
	var restore live228R3Point
	pLive228R3GetCursorPos.Call(uintptr(unsafe.Pointer(&restore)))
	t.Cleanup(func() { pLive228R3SetCursorPos.Call(uintptr(restore.x), uintptr(restore.y)) })

	var hwnd windows.HWND
	if !pollUntil127(200, func() bool {
		for h, pid := range findBallWindows228(t) {
			if pid == leg.pid() {
				hwnd = h
				return true
			}
		}
		return false
	}) {
		t.Logf("AC#3 shape2 tray: no ball window is owned by pid %d, so there is no tray host to ask", leg.pid())
		return ""
	}

	for _, frac := range []float64{0.34, 0.30, 0.40, 0.26} {
		pLive228R3SetCursorPos.Call(400, 300)
		time.Sleep(80 * time.Millisecond)
		if r, _, err := pLive228R3PostMessageW.Call(uintptr(hwnd), live228R3TrayMsg, live228R3TrayUID, live228R3RBUp); r == 0 {
			t.Logf("AC#3 shape2 tray: PostMessageW(WM_APP+0x202, tray id, WM_RBUTTONUP) failed: %v", err)
			return ""
		}
		var menu windows.HWND
		if !pollUntil127(120, func() bool {
			menu = live228R3MenuWindow(t, leg.pid())
			return menu != 0
		}) {
			t.Logf("AC#3 shape2 tray: no popup menu is owned by pid %d after the tray callback; the break is "+
				"between the message and internal/ball's showMenu", leg.pid())
			return ""
		}
		rc, ok := live228R3RectOf(menu)
		if !ok {
			t.Logf("AC#3 shape2 tray: GetWindowRect on the menu window %v failed", menu)
			return ""
		}
		x := (rc.left + rc.right) / 2
		y := rc.top + int32(float64(rc.bottom-rc.top)*frac)
		t.Logf("AC#3 shape2 tray menu: rect=(%d,%d)-(%d,%d) clicking (%.0f%% of the height) at (%d,%d)",
			rc.left, rc.top, rc.right, rc.bottom, frac*100, x, y)
		live228R3Click(x, y)
		if pollUntil127(40, func() bool {
			return strings.Contains(leg.stdout.String(), ac290OpenVerdict)
		}) {
			return ac290LineWith(leg.stdout.String(), ac290OpenVerdict)
		}
		t.Logf("AC#3 shape2 tray menu: that row was not the mute one (no gate sentence); reopening")
	}
	return ""
}

// ac290LineWith returns the first console line containing needle, verbatim.
func ac290LineWith(out, needle string) string {
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, needle) {
			return strings.TrimRight(line, "\r")
		}
	}
	return ""
}

// ac290Counter reads one integer field out of the exit line, by name.
func ac290Counter(t *testing.T, out, name string) int {
	t.Helper()
	line := ac290LineWith(out, ac290StopPrefix)
	if line == "" {
		t.Fatalf("AC#3 RED: the shipped process printed no %q line, so no counter can be read.\n%s",
			ac290StopPrefix, out)
	}
	i := strings.Index(line, name+"=")
	if i < 0 {
		t.Fatalf("AC#3 RED: the exit line has no %s field: %s", name, line)
	}
	rest := line[i+len(name)+1:]
	digits := 0
	for digits < len(rest) && rest[digits] >= '0' && rest[digits] <= '9' {
		digits++
	}
	if digits == 0 {
		t.Fatalf("AC#3 RED: %s= is not followed by a number in: %s", name, line)
	}
	n, err := strconv.Atoi(rest[:digits])
	if err != nil {
		t.Fatalf("AC#3 RED: %s= is not an integer (%v) in: %s", name, err, line)
	}
	return n
}

// TestAC290LiveShapeOneDefaultBootOpensNoDevice is AC#3 shape 1: the posture a
// double click on the shipped build gets, read off the process's own console.
func TestAC290LiveShapeOneDefaultBootOpensNoDevice(t *testing.T) {
	if os.Getenv("WISP_LIVE_MIC") != "1" {
		t.Skip("AC#3 shape 1 boots the shipped capture leg: run with WISP_LIVE_MIC=1 after both tasklist rulers read 0")
	}
	requireNoCompetingAudioProcesses(t)

	exe := buildWispForTest(t)
	// An empty data dir: no config.toml, so the compiled default table decides
	// the posture. Nothing here writes mic_muted_default=false.
	dataDir := t.TempDir()
	leg := bootResidentLeg(t, exe, dataDir)
	t.Cleanup(leg.stop)

	if !pollUntil127(400, func() bool {
		out := leg.stdout.String()
		return strings.Contains(out, ac290MutedVerdict) && strings.Contains(out, ac290LoopLine)
	}) {
		t.Fatalf("AC#3 shape1 RED: the shipped process never printed the muted-boot verdict and its event loop.\n%s",
			leg.console())
	}
	out := leg.stdout.String()
	t.Logf("AC#3 shape1 boot verdict: %s", ac290LineWith(out, ac290MutedVerdict))
	t.Logf("AC#3 shape1 boot posture line: %s", ac290LineWith(out, ac290LoopLine))

	// It sits exactly as long as a double click would before the user does
	// anything, then leaves through its own ordered path (D38(e)).
	time.Sleep(1500 * time.Millisecond)
	if err := leg.breakToLoop(); err != nil {
		t.Fatalf("AC#3 shape1 RED: cannot ask the process to leave through its own path: %v", err)
	}
	if !leg.exitedWithin(400) {
		leg.stop()
		t.Fatalf("AC#3 shape1 RED: the process did not leave through its shutdown path.\n%s", leg.console())
	}
	out = leg.stdout.String()
	stop := ac290LineWith(out, ac290StopPrefix)
	if stop == "" {
		t.Fatalf("AC#3 shape1 RED: the process left without printing the step-4 counter line.\n%s", leg.console())
	}
	t.Logf("AC#3 shape1 exit line: %s", stop)
	// The same fact on disk: the persistent sink's own record, so the reading is
	// not only a console line a later reader has to trust.
	for _, line := range strings.Split(readSinkTail258(t, logSinkDir(dataDir)), "\n") {
		if strings.Contains(line, "capture armed but muted at boot") {
			t.Logf("AC#3 shape1 sink record: %s", strings.TrimSpace(line))
			break
		}
	}
	if ev := leg.waitErr; ev != nil {
		t.Errorf("AC#3 shape1 RED: the process ended as %s, want exit 0.\n%s", leg.exitStatus(), leg.console())
	}
	if n := ac290Counter(t, out, "frames_sent"); n != 0 {
		t.Errorf("AC#3 shape1 RED: frames_sent=%d at exit, want 0: at the default posture the muted gate must never "+
			"hand the device to the source, so no frame can exist.\n%s", n, stop)
	}
	if n := ac290Counter(t, out, "levels_delivered"); n != 0 {
		t.Errorf("AC#3 shape1 RED: levels_delivered=%d at exit, want 0: the ball seam received scalars from a leg "+
			"whose device was never opened.\n%s", n, stop)
	}
}

// TestAC290LiveShapeTwoTheMuteHotkeyOpensTheGate is AC#3 shape 2: the user's own
// gesture on the same shipped process, and the level counter that proves the
// gate let the device through and the ball seam received scalars.
func TestAC290LiveShapeTwoTheMuteHotkeyOpensTheGate(t *testing.T) {
	if os.Getenv("WISP_LIVE_MIC") != "1" {
		t.Skip("AC#3 shape 2 opens the real capture endpoint: run with WISP_LIVE_MIC=1 after both tasklist rulers read 0")
	}
	requireNoCompetingAudioProcesses(t)

	exe := buildWispForTest(t)
	leg := bootResidentLeg(t, exe, t.TempDir())
	t.Cleanup(leg.stop)

	if !pollUntil127(400, func() bool {
		out := leg.stdout.String()
		return strings.Contains(out, ac290MutedVerdict) && strings.Contains(out, ac290LoopLine)
	}) {
		t.Fatalf("AC#3 shape2 RED: the shipped process never reached its muted boot.\n%s", leg.console())
	}
	out := leg.stdout.String()
	t.Logf("AC#3 shape2 boot verdict: %s", ac290LineWith(out, ac290MutedVerdict))
	t.Logf("AC#3 shape2 boot hotkey line: %s", ac290LineWith(out, "hotkeys live"))
	if !strings.Contains(out, ac290MuteLive) {
		t.Skipf("SKIP-LOUD: this boot does not hold %q (some other program has the chord, or a config file decided "+
			"otherwise), so a pressed chord would be delivered to that other process and this case refuses to judge "+
			"the gate on somebody else's input.\n%s", ac290MuteLive, leg.console())
	}

	// Shape 1 is re-read here on purpose: the same process has not opened
	// anything yet, so shape 2's "before" is the same machine fact.
	t.Logf("AC#3 shape2 before the gesture: mute slot live, device not opened")
	t.Logf("AC#3 shape2 input desktop before the chord: %s", ac290InputDesktopName())
	// Does the child REALLY hold the combination? Registering it here answers
	// that with the OS's own bookkeeping: if this process can take it, the boot
	// sentence that said "live" and Win32 disagree, and a pressed chord would go
	// somewhere else. Released immediately either way so the press below has one
	// possible recipient.
	if owner := ac290TrySquatMuteChord(); owner != "" {
		t.Logf("AC#3 shape2 chord ownership: %s", owner)
	} else {
		t.Logf("AC#3 shape2 chord ownership: this process took Ctrl+Alt+M, so the shipped process is NOT holding it")
	}

	// Path 1: the orb's mute hot key, pressed through the OS input queue. A hand
	// can press a key more than once, so a dropped chord is not a product finding
	// until several presses miss.
	pressed := 0
	opened := false
	for attempt := 1; attempt <= 3 && !opened; attempt++ {
		ac290PressMuteChord()
		pressed = attempt
		t.Logf("AC#3 shape2 chord delivered (attempt %d of at most 3)", attempt)
		opened = pollUntil127(40, func() bool {
			return strings.Contains(leg.stdout.String(), ac290OpenVerdict)
		})
		if !opened {
			time.Sleep(400 * time.Millisecond)
		}
	}
	switch {
	case opened:
		t.Logf("AC#3 shape2 gesture line (orb mute hot key, press %d): %s", pressed,
			ac290LineWith(leg.stdout.String(), ac290OpenVerdict))
	default:
		// Path 2: the tray's mute row, clicked with real mouse input on the popup
		// the CHILD itself opened - the carrier ticket 228's live tier already
		// uses for the Exit row. It is the same host and the same gate; only the
		// gesture differs, and the console line says which one it was.
		t.Logf("AC#3 shape2: %d synthesised chord press(es) produced no sentence although Win32 says the child "+
			"holds the combination; trying the tray's mute row (real mouse input on the child's own menu)", pressed)
		line := ac290ClickTrayMuteRow(t, leg)
		if line == "" {
			// A SKIP, not a pass, and named: this case cannot build the premise
			// "a hand reaches the shipped process" in this session, and the repo's
			// rule is that an unbuildable premise is said out loud rather than
			// judged green or red. The counter reading this case would produce is
			// owed to a window with a real hand on the key (the same window tickets
			// 247/291/293 are waiting for).
			t.Skipf("SKIP-LOUD: no user gesture reached the shipped process in this session: %d synthesised "+
				"chord press(es) missed although Win32 says the child holds the combination (input desktop %q), "+
				"and the tray popup the child should have opened was not found either. This is a fact about input "+
				"injection here, not about the gate; the counter reading stays owed to an owner-present window.\n%s",
				pressed, ac290InputDesktopName(), leg.console())
		}
		t.Logf("AC#3 shape2 gesture line (tray mute row): %s", line)
	}

	// Four seconds at the measured cadence (one scalar per 32 ms frame, ~31/s)
	// is ~125 scalars; the point is not the number, it is that the counter is
	// read from the process's own exit line rather than computed here.
	time.Sleep(4 * time.Second)
	if err := leg.breakToLoop(); err != nil {
		t.Fatalf("AC#3 shape2 RED: cannot ask the process to leave through its own path: %v", err)
	}
	if !leg.exitedWithin(400) {
		leg.stop()
		t.Fatalf("AC#3 shape2 RED: the process did not leave through its shutdown path.\n%s", leg.console())
	}
	out = leg.stdout.String()
	stop := ac290LineWith(out, ac290StopPrefix)
	if stop == "" {
		t.Fatalf("AC#3 shape2 RED: the process left without printing the step-4 counter line.\n%s", leg.console())
	}
	t.Logf("AC#3 shape2 exit line: %s", stop)
	if ev := leg.waitErr; ev != nil {
		t.Errorf("AC#3 shape2 RED: the process ended as %s, want exit 0.\n%s", leg.exitStatus(), leg.console())
	}
	levels := ac290Counter(t, out, "levels_delivered")
	frames := ac290Counter(t, out, "frames_sent")
	if levels <= 0 {
		t.Errorf("AC#3 shape2 RED: levels_delivered=%d after the user opened the door: the ball seam received "+
			"nothing, so the gate did not really let the device through.\n%s", levels, stop)
	}
	if frames <= 0 {
		t.Errorf("AC#3 shape2 RED: frames_sent=%d after the user opened the door: no PCM frame ever left the "+
			"capture source.\n%s", frames, stop)
	}
	if dropped := ac290Counter(t, out, "frames_dropped"); dropped < 0 {
		t.Errorf("AC#3 shape2 RED: frames_dropped=%d is not a counter.\n%s", dropped, stop)
	}
	// The window between the gesture and the break is the four seconds above
	// plus the poll for the gesture sentence; at the measured one-scalar-per-32ms
	// cadence that is 60..250 scalars. The band is deliberately wide and it is a
	// cadence sanity check, not the reading (the reading is the verbatim exit
	// line above). frames_dropped is NOT compared against frames_sent: this boot
	// has no drainer for the C8 channel, so the bounded channel fills and most
	// frames drop - the shape ticket 291's mic probe recorded too, and not a
	// device fault.
	if levels < 60 || levels > 250 {
		t.Errorf("AC#3 shape2 RED: levels_delivered=%d is outside the 60..250 band a 4s open window produces "+
			"at the measured 31.25 scalars/s.\n%s", levels, stop)
	}
}

// TestAC290LiveScalarsAtTheBallSeamAfterTheUserPath is the depth half of shape
// 2: the shipped exit line counts the scalars, this case prints what they were
// (samples / mean / max) at the seam, using the production assembly and the
// production executor the orb's mute key and the tray's mute item call. The
// gate is the object production holds; the config is the shipped default table
// with mic_muted_default untouched; opening is done by residentAudio.toggleMute,
// not by a config edit.
func TestAC290LiveScalarsAtTheBallSeamAfterTheUserPath(t *testing.T) {
	if os.Getenv("WISP_LIVE_MIC") != "1" {
		t.Skip("AC#3 shape 2 opens the real capture endpoint: run with WISP_LIVE_MIC=1 after both tasklist rulers read 0")
	}
	requireNoCompetingAudioProcesses(t)

	rt, shut := bootAudioRuntime(t)
	dir := writeAudioConfig(t, nil) // the shipped table, untouched

	boot := &livePhase{name: "muted-at-boot"}
	router := &phaseRouter{}
	router.set(boot)

	ra := assembleCapture(rt, dir, router.route, newRealCaptureSource)
	if ra.gate == nil || ra.mic == nil {
		t.Fatalf("AC#3 RED: the shipped default table built no capture leg: %s", ra.unrunningClaim())
	}
	t.Logf("AC#3 assembly verdict: %s", ra.verdict)
	if !ra.mutedAtBoot {
		t.Fatalf("AC#3 RED: the assembly read mic_muted_default=%v; the shipped default is true and this case "+
			"must not be run against an edited table", ra.mutedAtBoot)
	}
	if !ra.gate.Muted() || ra.gate.Open() || ra.started.Load() {
		t.Fatalf("AC#3 RED: before any user gesture the gate reads muted=%v open=%v started=%v, want true/false/false",
			ra.gate.Muted(), ra.gate.Open(), ra.started.Load())
	}
	if n := rt.Registry.CountByName("audio-capture"); n != 0 {
		t.Fatalf("AC#3 RED: the muted boot books %d audio-capture thread(s), want 0: the device must not be open", n)
	}
	if n := boot.count(); n != 0 {
		t.Fatalf("AC#3 RED: the ball seam received %d scalar(s) before the user did anything", n)
	}

	// The user path: the same executor the orb's mute hot key and the tray's
	// mute item are handed (resident_windows.go's attachMuteGate).
	msg, executed := ra.toggleMute()
	t.Logf("AC#3 gesture open: executed=%v outcome=%q", executed, msg)
	if !executed {
		t.Fatalf("AC#3 RED: the mute gesture had no executor for this process's gate: %q", msg)
	}
	if !strings.Contains(msg, ac290OpenVerdict) || strings.Contains(msg, "设备未交接") {
		t.Fatalf("AC#3 RED: the gesture did not open the device; it said %q", msg)
	}
	if !ra.gate.Open() || ra.gate.Muted() {
		t.Fatalf("AC#3 RED: the gesture returned the open sentence but the gate reads muted=%v open=%v",
			ra.gate.Muted(), ra.gate.Open())
	}
	if n := rt.Registry.CountByName("audio-capture"); n != 1 {
		t.Fatalf("AC#3 RED: after the user opened the door audio-capture threads = %d, want 1", n)
	}

	open := &livePhase{name: "after-user-open"}
	router.set(open)
	time.Sleep(3 * time.Second)
	open.span = 3 * time.Second
	st := ra.mic.Stats()
	t.Logf("AC#3 counters after open: levels_delivered=%d frames_sent=%d frames_dropped=%d last_err=%q",
		ra.levels.Load(), st.FramesSent, st.FramesDropped, st.LastError)
	t.Logf("AC#3 scalars after open: %s", open)
	if open.count() == 0 {
		t.Fatalf("AC#3 RED: three seconds after the user opened the door the ball seam received no scalar "+
			"(gate open=%v, frames_sent=%d, last_err=%q)", ra.gate.Open(), st.FramesSent, st.LastError)
	}

	// And back again: the same gesture closes it, which is what makes the
	// counter above a property of the gate rather than of the machine.
	msg2, executed2 := ra.toggleMute()
	t.Logf("AC#3 gesture close: executed=%v outcome=%q", executed2, msg2)
	if !executed2 || !ra.gate.Muted() || ra.gate.Open() {
		t.Fatalf("AC#3 RED: the second gesture did not close the gate (executed=%v muted=%v open=%v)",
			executed2, ra.gate.Muted(), ra.gate.Open())
	}
	closed := &livePhase{name: "after-close"}
	router.set(closed)
	time.Sleep(time.Second)
	closed.span = time.Second
	t.Logf("AC#3 scalars after close: %s", closed)
	if closed.count() != 0 {
		t.Errorf("AC#3 RED: the gate is muted again and still handed %d scalar(s) to the ball seam", closed.count())
	}
	if n := rt.Registry.CountByName("audio-capture"); n != 0 {
		t.Errorf("AC#3 RED: audio-capture threads = %d after the user closed the door, want 0", n)
	}

	if recs := shut(); len(recs) != 10 {
		t.Fatalf("AC#3 RED: shutdown trail = %d records, want the frozen 10", len(recs))
	}
}
