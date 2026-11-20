//go:build windows

package ball

// Ticket 245 AC#9 - the standing in-repo ruler for "a second process really
// receives Esc". The four-shape rig (E1 normal delivery / E2 self-claim /
// E3 idle process still delivers / E4 the return) was built on the acceptance
// side and the orchestrator's ruling on AC#9 (ledger A480) was that an
// acceptance-side rig guards no regression: the shapes have to be runnable
// from the tree, without a human driving anything. This file is that ruler for
// the BALL ITSELF (cmd/wisp/resident_esc_delivery_245_windows_test.go is the
// shipped-binary variant of the same four shapes).
//
// THE FOUR SHAPES, one timeline on one real ball:
//
//	shape1  the ball is up and idle, nothing waiting: the injected bare Esc
//	        still reaches the second process (keydown_esc=1) - the same
//	        reading E3 measured on the shipped leg;
//	shape2  in-process steady state: HotkeyReport().Bindings() holds the
//	        three always-on slots, the bare cancel slot is standby, no bare
//	        key is in the set (the AC#6② judgment, idleRosterReds245);
//	shape3  the card waits (Confirming + TakeEscForCancel): the same injected
//	        Esc now reaches the observer ZERO times and this ball's cancel
//	        gesture fires once (the veto works - D43 #22, one press, two
//	        verdicts);
//	shape4  the card is gone (ReleaseEscAfterSession): the observer receives
//	        the Esc again (keydown_esc=1), and the ball's gesture does NOT
//	        fire - "handed back" measured on the desktop, not on a flag.
//
// THE INSTRUMENT'S OWN CONTROLS run first (E1/E2, the same two the
// acceptance-side rig used): with nothing of ours running the observer gets
// the key (E1), and a window that registers the bare Esc itself loses the
// keydown and gets WM_HOTKEY instead (E2). Without that pair, the "1" readings
// below would prove nothing.
//
// THE SECOND PROCESS. The observer is THIS test binary exec'd in helper mode
// (the classic os/exec re-exec shape, the same one orb_datapath_windows_test.go
// uses): exec.Command(os.Args[0], "-test.run=^TestHelperProcess$") with an env
// var, the helper creating its own real top-level window, taking the
// foreground, waiting for the trigger over stdin, injecting one physical Esc
// through keybd_event and counting the WM_KEYDOWN/WM_SYSKEYDOWN/WM_HOTKEY
// messages it was delivered. No human starts anything; every child is reaped
// by t.Cleanup.
//
// WHO HELD THE KEY (the AC#8 ownership split, carried here on purpose). A "0"
// while the ball is supposed to be idle is not automatically the defect: the
// two decisive discriminators run before any foreign-owner skip - the ball's
// own cancel gesture (a leaked registration delivers WM_HOTKEY to OUR window
// and only ours) and then the Win32 claim probe. ours = hard RED; foreign =
// named SKIP-LOUD with the attempt count and the verbatim errno; nobody =
// RED (instrument). A dead instrument is never allowed to read as a pass.
//
// DESKTOP UNAVAILABLE. This ruler needs an interactive desktop (a real
// foreground window and injected input). Where that cannot be built the case
// SKIPs LOUDLY, naming the shape and carrying the numbers it measured; it
// never reports a green it did not take.
//
// Run it as the dispatch's gate runs it:
//
//	go test ./internal/ball -count=1 -run "245|Esc" -v
//
// Nothing here depends on a human pressing a key or starting a process: the
// observer takes the foreground itself, the test triggers the injection over
// a stdin line, and the SPI foreground-lock tweak the settle loop makes is
// restored before the helper exits.

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// ---------------------------------------------------------------------------
// Win32 surfaces this file needs that the package does not already carry at
// the default tier. Everything else reuses the product's own lazy procs
// (pRegisterHotKey, pUnregisterHotKey, pCreateWindowExW, pSetWindowLongPtrW,
// pShowWindow, pSetWindowPos, pPeekMessageW, pTranslateMessage,
// pDispatchMessageW, pPostQuitMessage, pKillTimer, pSetTimer, pPostMessageW,
// pSetForegroundWindow, pDestroyWindow, procGetModuleHandleW) and its
// constants (vkEscape, modNoRepeat, wmHotkey, wmTimer, msg).
// ---------------------------------------------------------------------------

const (
	escRulerHelperEnv245     = "WISP_245_ESC_DELIVERY_HELPER" // watch | steal
	escRulerModeWatch245     = "watch"                        // plain observer: the injected Esc must be delivered
	escRulerModeSteal245     = "steal"                        // control: the observer registers bare Esc itself
	escRulerProbeID245       = 0x245E                         // spare id for the parent's Win32 claim probe (NULL hwnd)
	escRulerStealID245       = 0x245F                         // the observer's own registration id (steal mode)
	escRulerQuitTimerID245   = 0x2451                         // the helper's self-quit timer id
	escRulerQuitTimerMs245   = 2600                           // pump window after READY; the quit timer is the only clock
	escRulerKeyup245         = 0x0002                         // KEYEVENTF_KEYUP
	escRulerEscScan245       = 0x01                           // Esc scan code, as the measured rig sent it
	escRulerVKControl245     = 0x11                           // the benign Ctrl tap that lets SetForegroundWindow through
	escRulerCtrlScan245      = 0x1D
	escRulerWmKeydown245     = 0x0100
	escRulerWmSysKeydown245  = 0x0104
	escRulerWmClose245       = 0x0010
	escRulerWmQuit245        = 0x0012
	escRulerGwlpWndProc245   = ^uintptr(3) // GWLP_WNDPROC == -4
	escRulerHwndTopmost245   = ^uintptr(0) // HWND_TOPMOST == -1
	escRulerWsOverlapped245  = 0x00CF0000  // WS_OVERLAPPEDWINDOW
	escRulerWsVisible245     = 0x10000000  // WS_VISIBLE
	escRulerSwShow245        = 5           // SW_SHOW (activates)
	escRulerSwpShow245       = 0x0040      // SWP_SHOWWINDOW
	escRulerPmRemove245      = 1           // PM_REMOVE
	escRulerSpiGetFgLock245  = 0x2002      // SPI_GETFOREGROUNDLOCKTIMEOUT
	escRulerSpiSetFgLock245  = 0x2001      // SPI_SETFOREGROUNDLOCKTIMEOUT
	escRulerWaitObject0245   = uintptr(0)
	escRulerWaitTimeout245   = uintptr(0x0102)
	escRulerQsAllInput245    = 0x00FF
	escRulerBallClass245     = "WispBallWindow"
	escRulerProto245         = "ESC245RULER"
	escRulerSettleBudget245  = 1500 * time.Millisecond
	escRulerObserverW245     = 260
	escRulerObserverH245     = 120
	escRulerHelperBudget245  = 90 * time.Second
	escRulerVetoBudget245    = 2 * time.Second
	escRulerClaimAttempts245 = 3
)

var (
	escRulerGetForegroundWindow245    = modUser32.NewProc("GetForegroundWindow")
	escRulerKeybdEvent245             = modUser32.NewProc("keybd_event")
	escRulerMsgWaitForMultipleObjects = modUser32.NewProc("MsgWaitForMultipleObjects")
	escRulerCallWindowProcW245        = modUser32.NewProc("CallWindowProcW")
	escRulerSystemParametersInfoW245  = modUser32.NewProc("SystemParametersInfoW")
	escRulerEnumWindows245            = modUser32.NewProc("EnumWindows")
	escRulerGetClassNameW245          = modUser32.NewProc("GetClassNameW")
	escRulerGetWindowThreadProcessId  = modUser32.NewProc("GetWindowThreadProcessId")
	escRulerGetWindowTextW245         = modUser32.NewProc("GetWindowTextW")
	escRulerIsWindow245               = modUser32.NewProc("IsWindow")
	escRulerObsProc245                = windows.NewCallback(escRulerObsWndProc245)
	escRulerPrevProc245               uintptr
	escRulerCntKeyDown245             int
	escRulerCntSysKey245              int
	escRulerCntHotkey245              int
	escRulerCntOther245               int
	escRulerForegroundLockTouched245  bool
	escRulerForegroundLockOldValue245 uint32
	escRulerForegroundLockRestored245 bool
)

// escRulerErrnoText245 renders a Win32 error the way the named skips carry it
// (verbatim text, empty for a zero errno).
func escRulerErrnoText245(e error) string {
	if e == nil {
		return ""
	}
	if errno, ok := e.(syscall.Errno); ok && errno == 0 {
		return ""
	}
	return e.Error()
}

// ---------------------------------------------------------------------------
// The helper process (the second real GUI process).
// ---------------------------------------------------------------------------

// TestHelperProcess is the classic helper entry. Run normally it does nothing
// but name itself; the ruler exec's this same test binary with
// -test.run=^TestHelperProcess$ and the helper env set, and then this function
// is the observer. The helper lives only in this one _windows_test.go file.
func TestHelperProcess(t *testing.T) {
	mode := os.Getenv(escRulerHelperEnv245)
	if mode == "" {
		t.Skip("helper entry for TestLive245EscDeliveryRulerFourShapes (exec'd by the ruler with " +
			"-test.run=^TestHelperProcess$ and " + escRulerHelperEnv245 + "=watch|steal); not a standalone test")
		return
	}
	escRulerHelperRun245(mode)
}

// escRulerObsWndProc245 counts the Esc delivery the OS really handed this
// window. It forwards everything to the STATIC procedure it subclassed, so the
// window keeps behaving.
func escRulerObsWndProc245(hwnd, m, wParam, lParam uintptr) uintptr {
	switch m {
	case escRulerWmKeydown245:
		if wParam == vkEscape {
			escRulerCntKeyDown245++
		} else {
			escRulerCntOther245++
		}
	case escRulerWmSysKeydown245:
		if wParam == vkEscape {
			escRulerCntSysKey245++
		} else {
			escRulerCntOther245++
		}
	case wmHotkey:
		escRulerCntHotkey245++
	}
	r, _, _ := escRulerCallWindowProcW245.Call(escRulerPrevProc245, hwnd, m, wParam, lParam)
	return r
}

// escRulerHelperRun245 is the observer: one real top-level window (system
// class STATIC, subclassed so no window class of ours is registered anywhere),
// foreground, trigger over stdin, one injected Esc, counts, one READING line.
// Single-threaded on a locked OS thread: the wait is
// MsgWaitForMultipleObjects over the stdin pipe handle and the window's
// messages, and the only clock is the Win32 quit timer.
func escRulerHelperRun245(mode string) {
	// Never let a caller hang on this harness.
	time.AfterFunc(60*time.Second, func() { os.Exit(4) })

	steal := mode == escRulerModeSteal245
	runtime.LockOSThread()

	// Foreground policy: the OS refuses SetForegroundWindow to a background
	// process for a short lock timeout. Zeroing it is the documented way a
	// test harness takes focus on an interactive desktop; the old value is
	// restored on the way out.
	var oldLock uint32
	escRulerSystemParametersInfoW245.Call(escRulerSpiGetFgLock245, 0, uintptr(unsafe.Pointer(&oldLock)), 0)
	escRulerSystemParametersInfoW245.Call(escRulerSpiSetFgLock245, 0, 0, 0)
	escRulerForegroundLockTouched245 = true
	escRulerForegroundLockOldValue245 = oldLock
	defer escRulerRestoreForegroundLock245()

	inst, _, _ := procGetModuleHandleW.Call(0)
	cls, _ := windows.UTF16PtrFromString("STATIC")
	title, _ := windows.UTF16PtrFromString("Wisp 245 Esc delivery observer")
	hwnd, _, err := pCreateWindowExW.Call(0, uintptr(unsafe.Pointer(cls)), uintptr(unsafe.Pointer(title)),
		escRulerWsOverlapped245|escRulerWsVisible245, 40, 40,
		escRulerObserverW245, escRulerObserverH245, 0, 0, inst, 0)
	if hwnd == 0 {
		fmt.Printf("%s FATAL create_window %v\n", escRulerProto245, err)
		os.Stdout.Sync()
		os.Exit(3)
	}

	escRulerPrevProc245, _, _ = pSetWindowLongPtrW.Call(hwnd, escRulerGwlpWndProc245, escRulerObsProc245)

	pShowWindow.Call(hwnd, escRulerSwShow245)
	pSetWindowPos.Call(hwnd, escRulerHwndTopmost245, 0, 0, 0, 0, escRulerSwpShow245)
	// A benign Ctrl tap: the shell lets the process that produced the last
	// input event take the foreground.
	escRulerKeybdEvent245.Call(escRulerVKControl245, escRulerCtrlScan245, 0, 0)
	escRulerKeybdEvent245.Call(escRulerVKControl245, escRulerCtrlScan245, escRulerKeyup245, 0)

	// Bounded settle: retry the foreground until it is ours, and say so on the
	// READY line either way (a window without focus gets no WM_KEYDOWN
	// whatever the desktop does, so the parent must not judge a false mine).
	fgMine := escRulerSettleForeground245(hwnd)

	regErr := ""
	if steal {
		// The registration happens BEFORE the READY line: whoever is going to
		// own the key must own it before anything is counted.
		r, _, e := pRegisterHotKey.Call(hwnd, escRulerStealID245, modNoRepeat, vkEscape)
		if r == 0 {
			regErr = escRulerErrnoText245(e)
		}
	}

	// The quit timer is the pump's only clock: it bounds the whole observation
	// so a parent that never triggers still gets a READING (and a dead harness
	// can never read as a pass).
	pSetTimer.Call(hwnd, escRulerQuitTimerID245, escRulerQuitTimerMs245, 0)

	fmt.Printf("%s READY hwnd=0x%X mine=%v steal=%v reg_err=%q\n",
		escRulerProto245, hwnd, fgMine, steal, regErr)
	os.Stdout.Sync()

	injected, injectMine := escRulerHelperPump245(hwnd)

	if steal {
		pUnregisterHotKey.Call(hwnd, escRulerStealID245)
	}
	pKillTimer.Call(hwnd, escRulerQuitTimerID245)
	pSetWindowLongPtrW.Call(hwnd, escRulerGwlpWndProc245, escRulerPrevProc245)
	pDestroyWindow.Call(hwnd)
	escRulerRestoreForegroundLock245()

	fmt.Printf("%s READING keydown_esc=%d syskeydown_esc=%d wm_hotkey=%d other=%d injected=%v mine_at_inject=%v steal=%v reg_err=%q\n",
		escRulerProto245, escRulerCntKeyDown245, escRulerCntSysKey245, escRulerCntHotkey245,
		escRulerCntOther245, injected, injectMine, steal, regErr)
	os.Stdout.Sync()
}

// escRulerSettleForeground245 retries SetForegroundWindow within a bounded
// budget and reports whether the window ended up owning the foreground.
func escRulerSettleForeground245(hwnd uintptr) bool {
	settle := time.NewTimer(escRulerSettleBudget245)
	defer settle.Stop()
	fgMine := false
	for {
		if fg, _, _ := escRulerGetForegroundWindow245.Call(); fg == hwnd {
			return true
		}
		select {
		case <-settle.C:
			pSetForegroundWindow.Call(hwnd)
			time.Sleep(100 * time.Millisecond)
			fg, _, _ := escRulerGetForegroundWindow245.Call()
			return fg == hwnd
		default:
			pSetForegroundWindow.Call(hwnd)
			time.Sleep(80 * time.Millisecond)
			if fg, _, _ := escRulerGetForegroundWindow245.Call(); fg == hwnd {
				fgMine = true
			}
		}
		if fgMine {
			return true
		}
	}
}

// escRulerRestoreForegroundLock245 puts the system foreground-lock timeout
// back exactly once (the helper exits through two doors: the deferred call and
// the explicit one before READING).
func escRulerRestoreForegroundLock245() {
	if !escRulerForegroundLockTouched245 || escRulerForegroundLockRestored245 {
		return
	}
	escRulerForegroundLockRestored245 = true
	escRulerSystemParametersInfoW245.Call(escRulerSpiSetFgLock245, 0,
		uintptr(escRulerForegroundLockOldValue245), 0)
}

// escRulerHelperPump245 waits for the stdin trigger and the window's messages
// on one thread (MsgWaitForMultipleObjects over the pipe handle and
// QS_ALLINPUT), injects ONE physical Esc at the trigger when the foreground is
// still ours, and pumps until the quit timer. It reports whether the key was
// really pressed and whether the foreground was ours at that moment.
func escRulerHelperPump245(hwnd uintptr) (injected, injectMine bool) {
	waitHandles := [1]uintptr{uintptr(os.Stdin.Fd())}
	for {
		r, _, _ := escRulerMsgWaitForMultipleObjects.Call(1,
			uintptr(unsafe.Pointer(&waitHandles[0])), 0, 3000, escRulerQsAllInput245)
		switch r {
		case escRulerWaitObject0245: // stdin: the trigger line (or the pipe broke)
			if !injected {
				fg, _, _ := escRulerGetForegroundWindow245.Call()
				injectMine = fg == hwnd
				if injectMine {
					escRulerKeybdEvent245.Call(vkEscape, escRulerEscScan245, 0, 0)
					time.Sleep(60 * time.Millisecond)
					escRulerKeybdEvent245.Call(vkEscape, escRulerEscScan245, escRulerKeyup245, 0)
					injected = true
				}
			}
		case escRulerWaitObject0245 + 1: // window messages
			for {
				var m msg
				pr, _, _ := pPeekMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0, escRulerPmRemove245)
				if pr == 0 {
					break
				}
				if m.message == escRulerWmQuit245 {
					return injected, injectMine
				}
				if (m.message == wmTimer && m.wParam == escRulerQuitTimerID245) ||
					m.message == escRulerWmClose245 {
					pPostQuitMessage.Call(0)
					continue
				}
				pTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
				pDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
			}
		case escRulerWaitTimeout245:
			// keep waiting: the quit timer (a message) ends this loop
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// ---------------------------------------------------------------------------
// The parent side: spawn the observer, trigger it, read its protocol lines.
// ---------------------------------------------------------------------------

// escRulerReading245 is one observation's full protocol, plus the raw output
// the verdicts quote.
type escRulerReading245 struct {
	mode       string
	ready      bool
	readyMine  bool
	keydown    int
	syskey     int
	hotkey     int
	other      int
	injected   bool
	injectMine bool
	regErr     string
	fatal      string
	raw        string
}

// escRulerObserve245 runs one observation: exec this test binary in helper
// mode, wait for its READY line, write the trigger, wait for READING. The
// helper's own quit timer bounds every path, and exec.CommandContext's
// deadline backs that up so a wedged child can never wedge the ruler.
func escRulerObserve245(t *testing.T, mode string) escRulerReading245 {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), escRulerHelperBudget245)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestHelperProcess$", "-test.count=1")
	cmd.Env = append(os.Environ(), escRulerHelperEnv245+"="+mode)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("ticket 245 AC#9 RED (instrument): the observer's stdout pipe could not be built: %v", err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("ticket 245 AC#9 RED (instrument): the observer's stdin pipe could not be built: %v", err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("ticket 245 AC#9 RED (instrument): the observer helper could not be started: %v", err)
	}
	t.Cleanup(func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	})

	obs := escRulerReading245{mode: mode}
	var raw strings.Builder
	sc := bufio.NewScanner(stdout)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		raw.WriteString(line + "\n")
		if rest, ok := strings.CutPrefix(line, escRulerProto245+" READY"); ok {
			obs.ready = true
			escRulerParseFields245(rest, &obs)
			if _, werr := stdin.Write([]byte("\n")); werr != nil {
				t.Fatalf("ticket 245 AC#9 RED (instrument): the trigger line could not reach the observer "+
					"helper: %v\n%s", werr, raw.String())
			}
			continue
		}
		if rest, ok := strings.CutPrefix(line, escRulerProto245+" READING"); ok {
			escRulerParseFields245(rest, &obs)
			break
		}
		if rest, ok := strings.CutPrefix(line, escRulerProto245+" FATAL"); ok {
			obs.fatal = strings.TrimSpace(rest)
			break
		}
	}
	_ = stdin.Close()
	if werr := cmd.Wait(); werr != nil && obs.fatal == "" && !obs.injected {
		obs.fatal = fmt.Sprintf("helper did not produce a READING (wait: %v, stderr: %s)", werr, stderr.String())
	}
	obs.raw = raw.String()
	return obs
}

// escRulerParseFields245 reads the key=value fields of a READY or READING
// protocol line into the reading.
func escRulerParseFields245(rest string, obs *escRulerReading245) {
	for _, f := range strings.Fields(rest) {
		k, v, _ := strings.Cut(f, "=")
		switch k {
		case "mine", "mine_at_inject":
			b := v == "true"
			if k == "mine" {
				obs.readyMine = b
			} else {
				obs.injectMine = b
			}
		case "injected":
			obs.injected = v == "true"
		case "keydown_esc":
			obs.keydown, _ = strconv.Atoi(v)
		case "syskeydown_esc":
			obs.syskey, _ = strconv.Atoi(v)
		case "wm_hotkey":
			obs.hotkey, _ = strconv.Atoi(v)
		case "other":
			obs.other, _ = strconv.Atoi(v)
		case "reg_err":
			if u, err := strconv.Unquote(v); err == nil {
				obs.regErr = u
			} else {
				obs.regErr = v
			}
		}
	}
}

// escRulerRequireUsable245 turns the instrument's own health into a verdict
// before any shape is judged: a helper that aborted on create_window is the
// named no-desktop skip; a helper that never got the foreground is the named
// inconclusive skip; anything else it could not do is a RED instrument.
func escRulerRequireUsable245(t *testing.T, obs escRulerReading245, where string) {
	t.Helper()
	if obs.fatal != "" {
		if strings.Contains(obs.fatal, "create_window") {
			t.Skipf("SKIP-LOUD (ticket 245 AC#9, NO INTERACTIVE DESKTOP at %s): the observer helper could not "+
				"create its window (%s), so this session has no desktop the Esc delivery shapes could be "+
				"measured on. No product reading was taken.", where, obs.fatal)
		}
		t.Fatalf("ticket 245 AC#9 RED (instrument) at %s: the observer helper aborted: %s\n%s",
			where, obs.fatal, obs.raw)
	}
	if !obs.ready {
		t.Fatalf("ticket 245 AC#9 RED (instrument) at %s: the observer helper never reported READY, so no "+
			"reading below could mean anything.\n%s", where, obs.raw)
	}
	if !obs.injected {
		t.Skipf("SKIP-LOUD (ticket 245 AC#9, FOREGROUND REFUSED at %s): the observer window did not own the "+
			"foreground at injection time (READY mine=%v, mine_at_inject=%v), so the trigger could not be "+
			"delivered to a window we own and no reading was taken (keydown_esc=%d, wm_hotkey=%d). This is a "+
			"desktop session limit, not a product verdict.", where, obs.readyMine, obs.injectMine, obs.keydown, obs.hotkey)
	}
}

// escRulerProbeEscClaimed245 asks Win32 - not our own bookkeeping - whether
// the bare Esc is claimed right now, on a NULL-window registration bound to
// this thread, and returns the numbers a named skip must carry. A successful
// probe is unregistered immediately, so a "free" answer leaves nothing behind.
func escRulerProbeEscClaimed245() (claimed bool, attempts int, detail string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var errs []string
	for i := 0; i < escRulerClaimAttempts245; i++ {
		attempts++
		r, _, e := pRegisterHotKey.Call(0, escRulerProbeID245, modNoRepeat, vkEscape)
		if r != 0 {
			pUnregisterHotKey.Call(0, escRulerProbeID245)
			return false, attempts, strings.Join(errs, "; ")
		}
		errno, _ := e.(syscall.Errno)
		errs = append(errs, fmt.Sprintf("attempt %d: errno %d %q", attempts, uint32(errno), e.Error()))
		time.Sleep(100 * time.Millisecond)
	}
	return true, attempts, strings.Join(errs, "; ")
}

// escRulerWaitFor245 polls until the condition holds or the budget is spent.
// The budget is a timer, never a wall-clock difference used as a verdict.
func escRulerWaitFor245(budget time.Duration, cond func() bool) bool {
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

// ---------------------------------------------------------------------------
// The quiet-desktop precondition (the requireQuietBallDesktop law, default
// tier's own copy so the ruler is self-contained).
// ---------------------------------------------------------------------------

type escRulerBallWindow245 struct {
	hwnd  uintptr
	pid   uint32
	title string
}

// escRulerBallWindows245 lists every top-level window of the ball class with
// its pid and title. Class matching is exactly what the live-guard file is
// careful about, so the result is reported, never quietly ignored.
func escRulerBallWindows245() []escRulerBallWindow245 {
	var out []escRulerBallWindow245
	cb := windows.NewCallback(func(hwnd, lParam uintptr) uintptr {
		buf := make([]uint16, 64)
		escRulerGetClassNameW245.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		if windows.UTF16ToString(buf) != escRulerBallClass245 {
			return 1
		}
		var pid uint32
		escRulerGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
		tb := make([]uint16, 128)
		n, _, _ := escRulerGetWindowTextW245.Call(hwnd, uintptr(unsafe.Pointer(&tb[0])), uintptr(len(tb)))
		out = append(out, escRulerBallWindow245{
			hwnd: hwnd, pid: pid, title: windows.UTF16ToString(tb[:n]),
		})
		return 1
	})
	escRulerEnumWindows245.Call(cb, 0)
	return out
}

// escRulerRequireQuietBallDesktop245 fails the case when the desktop is not
// ours alone: another process's ball holds hotkey registrations and the same
// window class, and any number taken now would be meaningless. Skipping past
// pollution is what hid it once (live-guard law), so nothing here skips on it.
func escRulerRequireQuietBallDesktop245(t *testing.T) {
	t.Helper()
	mine := uint32(os.Getpid())
	for _, w := range escRulerBallWindows245() {
		if w.pid != mine {
			t.Fatalf("foreign Wisp ball window on the desktop: hwnd=0x%X pid=%d title=%q. The Esc delivery "+
				"ruler must own the desktop: another ball process holds the window class and the global "+
				"hotkeys, so any number taken now is meaningless. Kill pid %d (a leftover balldebug/wisp) and "+
				"re-run. t.Skip is not an option here.", w.hwnd, w.pid, w.title, w.pid)
		}
		t.Fatalf("a Wisp ball window from THIS test process is alive before the ruler created one: "+
			"hwnd=0x%X title=%q - a previous test leaked its window.", w.hwnd, w.title)
	}
}

// ---------------------------------------------------------------------------
// The ruler itself.
// ---------------------------------------------------------------------------

// TestLive245EscDeliveryRulerFourShapes is AC#9's standing ruler for the ball:
// the instrument controls (E1/E2) and then the four shapes in one timeline
// against one real ball with the production default bindings (the shipped bare
// "Esc" cancel). The final log line carries every reading in one place.
func TestLive245EscDeliveryRulerFourShapes(t *testing.T) {
	escRulerRequireQuietBallDesktop245(t)

	// INSTRUMENT E1: baseline with nothing of ours running. A zero here is not
	// a product verdict: either the desktop already has a foreign owner
	// (SKIP-LOUD with numbers) or the rig itself is broken (RED).
	base := escRulerObserve245(t, escRulerModeWatch245)
	escRulerRequireUsable245(t, base, "instrument baseline (E1)")
	if base.keydown != 1 {
		claimed, attempts, detail := escRulerProbeEscClaimed245()
		if claimed {
			t.Skipf("SKIP-LOUD (ticket 245 AC#9, FOREIGN OWNER at baseline): an injected bare Esc reached no "+
				"foreground window while nothing of ours was running (keydown_esc=%d, wm_hotkey=%d), and Win32 "+
				"says the key is claimed: RegisterHotKey(foreign id) failed %d/3 times with %s. This desktop "+
				"already has an owner for the bare Esc, so delivery cannot be measured here; no product "+
				"reading was taken.", base.keydown, base.hotkey, attempts, detail)
		}
		t.Fatalf("ticket 245 AC#9 RED (instrument): an injected bare Esc reached no foreground window while "+
			"nothing of ours was running (keydown_esc=%d, wm_hotkey=%d), and Win32 says the bare Esc is FREE "+
			"(%d/3 attempts, no 1409). The observer or the injection half of this rig failed, so no reading "+
			"below could mean anything. Observer raw:\n%s", base.keydown, base.hotkey, attempts, base.raw)
	}

	// INSTRUMENT E2: the ruler's positive control. A process that registers the
	// bare Esc globally (this observer itself) must lose the keydown and
	// receive WM_HOTKEY instead - the proof that a "1" below is not the rig
	// being deaf, and the same reading the pre-245 product produced.
	steal := escRulerObserve245(t, escRulerModeSteal245)
	escRulerRequireUsable245(t, steal, "instrument control (E2)")
	if steal.regErr != "" {
		if strings.Contains(steal.regErr, "1409") {
			t.Skipf("SKIP-LOUD (ticket 245 AC#9, FOREIGN OWNER during control): the observer could not "+
				"register the bare Esc for its own control (%s): an owner appeared between E1 and E2, so this "+
				"desktop cannot judge delivery. No product reading was taken.", steal.regErr)
		}
		t.Fatalf("ticket 245 AC#9 RED (instrument): the observer could not register the bare Esc for its own "+
			"control: %s. Observer raw:\n%s", steal.regErr, steal.raw)
	}
	if steal.keydown != 0 || steal.hotkey != 1 {
		t.Fatalf("ticket 245 AC#9 RED (instrument blind): with its OWN global bare-Esc registration the "+
			"observer still received keydown_esc=%d (want 0) and wm_hotkey=%d (want 1); the delivery ruler "+
			"cannot see a stolen key, so no product reading below means anything. Observer raw:\n%s",
			steal.keydown, steal.hotkey, steal.raw)
	}

	// The ball under measurement: the production default set, whose cancel slot
	// is the bare "Esc" this whole ticket is about. The ball window is created
	// hidden (StartHidden) - visibility plays no part in RegisterHotKey, and a
	// hidden ball never fights the observer for the foreground.
	var veto atomic.Int64
	b, err := New(Options{
		Initial:     statemachine.StateSleeping,
		Hotkeys:     DefaultHotkeys(),
		StartHidden: true,
		WindowTitle: fmt.Sprintf("Wisp245EscRuler-%d", os.Getpid()),
		Events:      Events{OnCancelHotkey: func() { veto.Add(1) }},
	})
	if err != nil {
		// E1/E2 just proved an interactive desktop exists, so this is not a
		// "no desktop" condition: a ball that cannot boot there is a product
		// failure and reads as one.
		t.Fatalf("ticket 245 AC#9 RED (PRODUCT): the ball itself failed to boot on a desktop the instrument "+
			"just proved usable (E1 baseline=%d, E2 control=%d/%d): %v", base.keydown, steal.keydown,
			steal.hotkey, err)
	}
	hwnd := b.DebugHWND()
	t.Cleanup(func() {
		b.Close()
		if r, _, _ := escRulerIsWindow245.Call(uintptr(hwnd)); r != 0 {
			t.Errorf("the ruler's ball window %v is still alive after Close", hwnd)
		}
		escRulerRequireQuietBallDesktop245(t)
	})

	// SHAPE 2 (in-process steady state): the idle roster is the three
	// always-on slots, the bare cancel slot is standby, no bare key anywhere -
	// the AC#6② judgment, read through the one nail every tier shares.
	rep := b.HotkeyReport()
	if reds := idleRosterReds245(rep); len(reds) != 0 {
		t.Fatalf("ticket 245 AC#9 shape2 RED: the idle ball's roster is not the standby one (want the three "+
			"always-on slots and no bare key): %v\n%+v", reds, rep.Bindings())
	}
	if n := len(rep.Live()); n != 3 {
		t.Fatalf("ticket 245 AC#9 shape2 RED: the idle registration set holds %d entries, want exactly the "+
			"three always-on slots (summon/mute/panel; the bare cancel slot is standby): %+v",
			n, rep.Bindings())
	}
	if line, ok := rep.Binding(hkCancel); !ok || line.Status != HotkeyStandby {
		t.Fatalf("ticket 245 AC#9 shape2 RED: the bare cancel slot reads %v while idle, want standby: %+v",
			line.Status, rep.Bindings())
	}
	if b.EscTakenOver() {
		t.Fatal("ticket 245 AC#9 shape2 RED: EscTakenOver() is true while idle")
	}
	t.Logf("AC#9 shape2 (steady state): Bindings() live=3 cancel=standby bare-held=false (summon/mute/panel)")

	// SHAPE 1 (idle delivery): the ball is up and idle, and the second process
	// still receives the Esc. The ball's own gesture is the first
	// ours-vs-foreign discriminator (a leaked idle registration delivers
	// WM_HOTKEY to OUR window and only ours).
	before1 := veto.Load()
	idle := escRulerObserve245(t, escRulerModeWatch245)
	escRulerRequireUsable245(t, idle, "shape1 (idle ball)")
	if idle.keydown != 1 || veto.Load() != before1 {
		if got := veto.Load(); got > before1 {
			t.Fatalf("ticket 245 AC#9 shape1 RED (PRODUCT): the IDLE ball held the bare Esc: the observer "+
				"received keydown_esc=%d (want 1) while this ball's cancel gesture fired %d time(s) - that "+
				"gesture can only come from a registration on this ball's own window, which is exactly the "+
				"defect ticket 245 exists to keep out. Report:\n%+v", idle.keydown, got-before1, rep.Bindings())
		}
		claimed, attempts, detail := escRulerProbeEscClaimed245()
		if claimed {
			t.Skipf("SKIP-LOUD (ticket 245 AC#9, FOREIGN OWNER at shape1): the ball is idle and its own set is "+
				"clean, yet the injected Esc reached the observer only %d time(s) and Win32 says the key is "+
				"claimed (%d/3 attempts: %s) - a foreign program owns the bare Esc, so delivery cannot be "+
				"judged on this desktop.", idle.keydown, attempts, detail)
		}
		t.Fatalf("ticket 245 AC#9 RED (instrument) at shape1: the idle ball's reading is keydown_esc=%d "+
			"(want 1) with no gesture here (veto delta=0) and no foreign claim (%d/3 attempts, no 1409) - the "+
			"rig measured nothing it can name.\n%s", idle.keydown, attempts, idle.raw)
	}
	t.Logf("AC#9 shape1 (idle delivery): observer keydown_esc=%d, ball veto delta=0", idle.keydown)

	// SHAPE 3 (the card waits): Confirming borrows the bare Esc; the SAME
	// injected Esc now reaches the observer ZERO times and this ball's cancel
	// gesture fires once - one press, both verdicts (borrow + veto).
	b.SetState(statemachine.StateConfirming)
	b.TakeEscForCancel()
	repB := b.HotkeyReport()
	lineB, okB := repB.Binding(hkCancel)
	if !okB || lineB.Status != HotkeyLive || !b.EscTakenOver() {
		if okB && lineB.Status.Attempted() {
			t.Skipf("SKIP-LOUD (ticket 245 AC#9, BORROW REFUSED at shape3): Win32 refused the card's borrow "+
				"(%v): another program on this desktop owns the bare Esc right now, so the borrowed shape "+
				"cannot be measured here. The report line: %+v", lineB.Err, lineB)
		}
		t.Fatalf("ticket 245 AC#9 shape3 RED: the card's borrow did not land (cancel line status %v, "+
			"EscTakenOver=%v): %+v", lineB.Status, b.EscTakenOver(), repB.Bindings())
	}
	if lineB.Acc.Mods != modNoRepeat || lineB.Acc.VK != vkEscape {
		t.Fatalf("ticket 245 AC#9 shape3 RED: the borrowed accelerator is %+v, want the production default "+
			"bare pair {MOD_NOREPEAT 0x4000, VK_ESCAPE 0x1B}: %+v", lineB.Acc, repB.Bindings())
	}
	if n := len(repB.Live()); n != 4 {
		t.Fatalf("ticket 245 AC#9 shape3 RED: the borrowed set holds %d live entries, want 4 (the three "+
			"always-on slots plus the borrowed cancel): %+v", n, repB.Bindings())
	}
	before3 := veto.Load()
	borrow := escRulerObserve245(t, escRulerModeWatch245)
	escRulerRequireUsable245(t, borrow, "shape3 (card waiting)")
	vetoOK := escRulerWaitFor245(escRulerVetoBudget245, func() bool { return veto.Load()-before3 == 1 })
	vetoGot := veto.Load() - before3
	if borrow.keydown != 0 || borrow.hotkey != 0 || vetoGot != 1 {
		if borrow.keydown != 0 || borrow.hotkey != 0 {
			t.Fatalf("ticket 245 AC#9 shape3 RED (PRODUCT): the waiting card did not take the bare Esc from "+
				"the desktop: the observer still received keydown_esc=%d and wm_hotkey=%d (want 0 and 0) while "+
				"the borrow is on this ball's own report (%+v) - a globally registered key Win32 hands to the "+
				"foreground anyway is the pre-245 defect, alive again", borrow.keydown, borrow.hotkey, repB.Bindings())
		}
		t.Fatalf("ticket 245 AC#9 shape3 RED: the borrowed Esc vanished: the observer received nothing "+
			"(keydown_esc=%d, wm_hotkey=%d, want 0/0) and this ball's cancel gesture fired %d time(s) (want "+
			"exactly 1, settled=%v) - the press went nowhere this ruler can name.\n%s",
			borrow.keydown, borrow.hotkey, vetoGot, vetoOK, borrow.raw)
	}
	t.Logf("AC#9 shape3 (borrow + veto): observer keydown_esc=0, ball veto delta=1 (the cancel gesture fired)")

	// SHAPE 4 (the card is gone): ReleaseEscAfterSession hands the key back,
	// and the desktop - not a flag - says so: the observer receives the Esc
	// again and the ball's gesture stays silent. This is the reading a
	// bookkeeping-only release cannot fake.
	b.SetState(statemachine.StateSleeping)
	b.ReleaseEscAfterSession()
	if b.EscTakenOver() {
		t.Fatal("ticket 245 AC#9 shape4 RED: EscTakenOver() is still true after the session end")
	}
	repR := b.HotkeyReport()
	if lineR, ok := repR.Binding(hkCancel); !ok || lineR.Status != HotkeyStandby {
		t.Fatalf("ticket 245 AC#9 shape4 RED: the cancel slot reads %v after the return, want standby: %+v",
			lineR.Status, repR.Bindings())
	}
	if reds := idleRosterReds245(repR); len(reds) != 0 {
		t.Fatalf("ticket 245 AC#9 shape4 RED: the roster after the return is not the idle one: %v\n%+v",
			reds, repR.Bindings())
	}
	before4 := veto.Load()
	back := escRulerObserve245(t, escRulerModeWatch245)
	escRulerRequireUsable245(t, back, "shape4 (card gone)")
	if back.keydown != 1 {
		if escRulerWaitFor245(escRulerVetoBudget245, func() bool { return veto.Load() > before4 }) {
			t.Fatalf("ticket 245 AC#9 shape4 RED (PRODUCT): the return did not hand the bare Esc back: the "+
				"observer received keydown_esc=%d (want 1) and this ball's cancel gesture fired after the "+
				"release (%d time(s) on this shape) - that gesture can only come from a registration on this "+
				"ball's own window, so the release dropped the bookkeeping but not the key. Report after the "+
				"return:\n%+v", back.keydown, veto.Load()-before4, repR.Bindings())
		}
		claimed, attempts, detail := escRulerProbeEscClaimed245()
		if claimed {
			t.Skipf("SKIP-LOUD (ticket 245 AC#9, FOREIGN OWNER after the return): after the card was gone and "+
				"the return reported clean, the injected Esc reached the observer only %d time(s) and Win32 "+
				"says the key is claimed (%d/3 attempts: %s) - a foreign owner appeared mid-run, so the "+
				"return half cannot be judged on this desktop.", back.keydown, attempts, detail)
		}
		t.Fatalf("ticket 245 AC#9 RED (instrument) at shape4: the returned reading is keydown_esc=%d (want 1) "+
			"with no gesture here and no foreign claim (%d/3 attempts, no 1409) - the rig measured nothing it "+
			"can name.\n%s", back.keydown, attempts, back.raw)
	}
	if got := veto.Load(); got != before4 {
		t.Fatalf("ticket 245 AC#9 shape4 RED: the ball received a cancel gesture %d time(s) while the "+
			"observer ALSO received the key (keydown_esc=%d) - two owners of one key, neither of them "+
			"explainable by the report:\n%+v", got-before4, back.keydown, repR.Bindings())
	}
	t.Logf("AC#9 shape4 (return): observer keydown_esc=%d, ball veto delta=0", back.keydown)

	t.Logf("AC#9 ruler four shapes: E1 baseline=%d E2 control=%d/%d shape1_idle=%d "+
		"shape2_roster=live3/cancel=standby shape3_borrow=observer0+veto1 shape4_return=%d",
		base.keydown, steal.keydown, steal.hotkey, idle.keydown, back.keydown)
}
