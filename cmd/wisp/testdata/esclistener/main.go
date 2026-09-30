//go:build windows

// Command esclistener is a SECOND REAL PROCESS with its own window, message pump
// and keyboard input. It answers the one question an in-process probe cannot:
//
//	does another application on this desktop actually get the bare Esc key?
//
// WHY IT LIVES HERE. The rig was built on the acceptance side of ticket 245
// (.scratch/wisp/probes/245/v1/esclistener/main.go) and the orchestrator's
// reading of that table (ledger A480, ticket 245 AC#9) was: an acceptance-side
// rig guards no regression, so the "a second process really receives Esc" half
// has to be in the repository. Ticket 246 AC#3 makes the same shape load-bearing
// for its own judgement (借到 / 否决生效 / 归还后另一进程收到), so this is that
// move, with the interface changed so the injection is triggered by the test
// instead of by a stop watch.
//
// WHAT IT MEASURES. It creates a real top-level window (system class STATIC,
// subclassed so no window class of ours is registered anywhere), takes the
// foreground, waits for a line on stdin (or -delay, or EOF), injects one
// physical-level Esc through keybd_event, and counts the WM_KEYDOWN /
// WM_SYSKEYDOWN messages with VK_ESCAPE that THIS window was delivered.
//
//	-watch   no global hotkey here; the injected Esc must be delivered
//	-steal   register bare Esc as a global hotkey on THIS window first; the same
//	         injected Esc must then NOT be delivered, and this window must get
//	         one WM_HOTKEY instead
//
// -steal is the ruler's own positive control: if a window that registered the
// key itself still sees the keydown, this program is blind and no -watch reading
// means anything.
//
// Protocol on stdout (one line each, flushed):
//
//	READY hwnd=0x… foreground=0x… foreground_is_mine=true|false steal=true|false
//	READING keydown_esc=N syskeydown_esc=N wm_hotkey=N other_keys=N …
//
// It exits 0 in both modes: the numbers are the verdict, not the exit code.
//
// It is a test fixture, so it is under testdata/: the go tool's package walk and
// this repository's d22scan walk both skip that directory by rule, and the case
// that uses it builds it with an explicit `go build ./testdata/esclistener`.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	procCreateWindowExW  = user32.NewProc("CreateWindowExW")
	procShowWindow       = user32.NewProc("ShowWindow")
	procSetWindowLongPtr = user32.NewProc("SetWindowLongPtrW")
	procCallWindowProc   = user32.NewProc("CallWindowProcW")
	procGetMessageW      = user32.NewProc("GetMessageW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
	procDispatchMessage  = user32.NewProc("DispatchMessageW")
	procPostMessageW     = user32.NewProc("PostMessageW")
	procSetForeground    = user32.NewProc("SetForegroundWindow")
	procGetForeground    = user32.NewProc("GetForegroundWindow")
	procSetWindowPos     = user32.NewProc("SetWindowPos")
	procKeybdEvent       = user32.NewProc("keybd_event")
	procRegisterHotKey   = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey = user32.NewProc("UnregisterHotKey")
	procSysParamsInfo    = user32.NewProc("SystemParametersInfoW")
	procPostQuitMessage  = user32.NewProc("PostQuitMessage")
	procGetModuleHandle  = kernel32GetModuleHandleW()
)

// negative Win32 indices and handles as pointer-width unsigned values
var (
	gwlpWndProcVal = ^uintptr(3) // GWLP_WNDPROC == -4
	topMostVal     = ^uintptr(0) // HWND_TOPMOST == -1
)

func kernel32GetModuleHandleW() *syscall.LazyProc {
	k := syscall.NewLazyDLL("kernel32.dll")
	return k.NewProc("GetModuleHandleW")
}

const (
	wmDestroy    = 0x0002
	wmClose      = 0x0010
	wmKeydown    = 0x0100
	wmSysKeydown = 0x0104
	wmHotkey     = 0x0312

	vkEscape  = 0x1B
	vkControl = 0x11

	modNoRepeat      = 0x4000
	spareHKID        = 900
	swShow           = 5
	wsVisible        = 0x10000000
	wsOverlappedWin  = 0x00CF0000
	swpShow          = 0x0040
	spiSetFgLockTime = 0x2001 // SPI_SETFOREGROUNDLOCKTIMEOUT

	keyeventfKeyup = 0x0002
)

type msg struct {
	hwnd   uintptr
	msg    uint32
	wParam uintptr
	lParam uintptr
	time   uint32
	ptX    int32
	ptY    int32
}

var (
	cntKeyDown   atomic.Int64
	cntSysKey    atomic.Int64
	cntHotkey    atomic.Int64
	cntOther     atomic.Int64
	oldProc      uintptr
	listenerHwnd atomic.Uintptr
)

func wndProc(hwnd, m, wParam, lParam uintptr) uintptr {
	switch m {
	case wmKeydown:
		if wParam == vkEscape {
			cntKeyDown.Add(1)
		} else {
			cntOther.Add(1)
		}
	case wmSysKeydown:
		if wParam == vkEscape {
			cntSysKey.Add(1)
		} else {
			cntOther.Add(1)
		}
	case wmHotkey:
		cntHotkey.Add(1)
	case wmClose:
		procPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := procCallWindowProc.Call(oldProc, hwnd, m, wParam, lParam)
	return r
}

func wide(s string) *uint16 {
	p, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		panic(err)
	}
	return p
}

func main() {
	watch := flag.Bool("watch", false, "no global hotkey here; the injected Esc must be delivered")
	steal := flag.Bool("steal", false, "register bare Esc globally here; the injected Esc must NOT be delivered")
	hold := flag.Int("hold", 0, "extra milliseconds to stay up after registering (foreign-owner mode)")
	delay := flag.Int("delay", -1, "milliseconds to wait before injecting, instead of waiting on stdin")
	after := flag.Int("after", 1200, "milliseconds of pumping after the injection")
	flag.Parse()

	// Never let a caller hang on this harness.
	time.AfterFunc(60*time.Second, func() { os.Exit(4) })

	if !*watch && !*steal {
		fmt.Println("esclistener: give -watch or -steal")
		os.Exit(2)
	}

	cb := syscall.NewCallback(wndProc)
	done := make(chan result, 1)
	fork := func() {
		runtime.LockOSThread()
		done <- run(cb, *steal, *hold, *delay, *after)
	}
	go fork()
	r := <-done

	fmt.Printf("esclistener: hwnd=0x%X foreground=0x%X foreground_is_mine=%v steal=%v "+
		"keydown_esc=%d syskeydown_esc=%d wm_hotkey=%d other_keys=%d\n",
		r.hwnd, r.fg, r.fg == r.hwnd, r.steal, r.keydown, r.syskey, r.hotkey, r.other)
	if r.regErr != "" {
		fmt.Printf("esclistener: RegisterHotKey(bare Esc) said: %s\n", r.regErr)
	}
	if !r.fgIsMine {
		fmt.Println("esclistener: WARNING this window never owned the foreground, so the delivery " +
			"counting below is INCONCLUSIVE (a window without focus gets no WM_KEYDOWN whatever Wisp does)")
	}
}

type result struct {
	hwnd     uintptr
	fg       uintptr
	fgIsMine bool
	steal    bool
	keydown  int64
	syskey   int64
	hotkey   int64
	other    int64
	regErr   string
}

func run(cb uintptr, steal bool, holdMs, delayMs, afterMs int) result {
	// Foreground policy: the OS refuses SetForegroundWindow to a background
	// process for a short lock timeout. Zeroing it is the documented way a test
	// harness takes focus on an interactive desktop.
	procSysParamsInfo.Call(uintptr(spiSetFgLockTime), 0, 0, 0)

	inst, _, _ := procGetModuleHandle.Call(0)
	hwnd, _, err := procCreateWindowExW.Call(0, uintptr(unsafe.Pointer(wide("STATIC"))),
		uintptr(unsafe.Pointer(wide("Wisp 246 Esc observer"))),
		wsOverlappedWin|wsVisible, 40, 40, 420, 220, 0, 0, inst, 0)
	if hwnd == 0 {
		fmt.Printf("esclistener: CreateWindowExW failed: %v\n", err)
		os.Exit(3)
	}
	listenerHwnd.Store(hwnd)

	oldProc, _, _ = procSetWindowLongPtr.Call(hwnd, gwlpWndProcVal, cb)

	procShowWindow.Call(hwnd, swShow)
	procSetWindowPos.Call(hwnd, topMostVal, 40, 40, 0, 0, swpShow)
	// A benign Ctrl tap: the shell lets the process that produced the last input
	// event take the foreground.
	procKeybdEvent.Call(vkControl, 0x1D, 0, 0)
	procKeybdEvent.Call(vkControl, 0x1D, keyeventfKeyup, 0)
	procSetForeground.Call(hwnd)

	res := result{hwnd: hwnd, steal: steal}
	res.fg, _, _ = procGetForeground.Call()
	res.fgIsMine = res.fg == hwnd

	// The registration happens BEFORE the ready line: whoever is going to own
	// the key must own it before the test starts counting anything.
	if steal {
		r, _, e := procRegisterHotKey.Call(hwnd, spareHKID, modNoRepeat, vkEscape)
		if r == 0 {
			res.regErr = errnoText(e)
		}
	}

	// Pump on this thread while a helper waits for the trigger, so the window
	// answers messages the whole time (a frozen window would "receive" nothing
	// whatever the desktop did, and that would read as a pass).
	go func() {
		waitForTrigger(delayMs)
		procKeybdEvent.Call(vkEscape, 0x01, 0, 0)
		time.Sleep(60 * time.Millisecond)
		procKeybdEvent.Call(vkEscape, 0x01, keyeventfKeyup, 0)
		live := afterMs
		if holdMs > live {
			live = holdMs
		}
		time.Sleep(time.Duration(live) * time.Millisecond)
		procPostMessageW.Call(hwnd, wmClose, 0, 0)
	}()

	fmt.Printf("READY hwnd=0x%X foreground=0x%X foreground_is_mine=%v steal=%v\n",
		hwnd, res.fg, res.fgIsMine, steal)
	os.Stdout.Sync()

	var m msg
	for {
		r, _, e := procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) == 0 {
			break // WM_QUIT
		}
		if int32(r) == -1 {
			fmt.Printf("esclistener: GetMessageW error: %s\n", errnoText(e))
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&m)))
	}

	if steal {
		procUnregisterHotKey.Call(hwnd, spareHKID)
	}
	procSetWindowLongPtr.Call(hwnd, gwlpWndProcVal, oldProc)

	res.keydown = cntKeyDown.Load()
	res.syskey = cntSysKey.Load()
	res.hotkey = cntHotkey.Load()
	res.other = cntOther.Load()

	fmt.Printf("READING keydown_esc=%d syskeydown_esc=%d wm_hotkey=%d other_keys=%d foreground_is_mine=%v steal=%v\n",
		res.keydown, res.syskey, res.hotkey, res.other, res.fgIsMine, res.steal)
	os.Stdout.Sync()
	return res
}

// waitForTrigger is the deterministic half of the rig: the test decides the
// moment the key is pressed, so a 2-3 second L1 window never races a stop watch.
// With -delay set it is a plain sleep (for hand-running); otherwise it reads one
// line from stdin, and stdin at EOF (a caller that never writes and never closes)
// injects at once.
func waitForTrigger(delayMs int) {
	if delayMs >= 0 {
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
		return
	}
	_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
}

func errnoText(e error) string {
	if e == nil || e == syscall.Errno(0) {
		return ""
	}
	return e.Error()
}
