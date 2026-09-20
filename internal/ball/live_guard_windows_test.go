//go:build windows && winlive

package ball

// Shared plumbing for the live (real-desktop) tests. Ticket 64 rewrote the
// rules here after two measurement reports were silently wrong:
//
//  1. A measurement may never find a Wisp ball window it did not create.
//     requireQuietBallDesktop enumerates every top-level window of class
//     "WispBallWindow" and FAILS with the offending pid when one is not ours:
//     the stray balldebug.exe made TestBallLiveLifecycle go red for the wrong
//     reason and made "all four hotkeys failed" look like a product defect
//     (registry A5). Skipping past that is what hid it, so nothing here skips
//     on pollution - and nothing here skips on a failed assertion either.
//  2. Liveness is asserted by the test's OWN hwnd/title, never by class name.
//     Every live ball gets a unique window title (Options.WindowTitle).

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// Test-local Win32 entry points (production does not need any of these).
var (
	pIsWindow                 = modUser32.NewProc("IsWindow")
	pGetClassNameW            = modUser32.NewProc("GetClassNameW")
	pGetWindowThreadProcessID = modUser32.NewProc("GetWindowThreadProcessId")
	pGetWindowTextW           = modUser32.NewProc("GetWindowTextW")
	pSendMessageW             = modUser32.NewProc("SendMessageW")
	pWindowFromPoint          = modUser32.NewProc("WindowFromPoint")
	pGetForegroundWindow      = modUser32.NewProc("GetForegroundWindow")
	pSetWindowLongPtrW        = modUser32.NewProc("SetWindowLongPtrW")
	pGetWindowLongPtrW        = modUser32.NewProc("GetWindowLongPtrW")
	pCallWindowProcW          = modUser32.NewProc("CallWindowProcW")
	pEnumThreadWindows        = modUser32.NewProc("EnumThreadWindows")
	pKeybdEvent               = modUser32.NewProc("keybd_event")
	pGetAsyncKeyState         = modUser32.NewProc("GetAsyncKeyState")
	pGetGuiResources          = modUser32.NewProc("GetGuiResources")
)

const (
	gwlpWndProc    = ^uintptr(3) // -4
	irKeyup        = 0x0002
	vkControl      = 0x11
	vkMenu         = 0x12
	vkShift        = 0x10
	vkLWin         = 0x5B
	ballClassName  = "WispBallWindow"
	grUserObjects  = 1
	grGdiObjects   = 2
	keyStateDownLo = 0x8000
)

// wndProcProbe is the slot the subclassed window procedure writes into
// (one live test at a time; these tests do not run in parallel).
var (
	wndProcProbe   atomic.Pointer[timerCount]
	prevWndProc    atomic.Uintptr
	subclassActive atomic.Bool
)

type timerCount struct{ wmTimer, total atomic.Uint64 }

// ballWndProcProbe replaces the ball window procedure and forwards to the
// procedure that was really installed before, so the window keeps behaving.
func ballWndProcProbe(hwnd, msg, wParam, lParam uintptr) uintptr {
	if p := wndProcProbe.Load(); p != nil {
		p.total.Add(1)
		if msg == wmTimer {
			p.wmTimer.Add(1)
		}
	}
	r, _, _ := pCallWindowProcW.Call(prevWndProc.Load(), hwnd, msg, wParam, lParam)
	return r
}

var probePtr = windows.NewCallback(ballWndProcProbe)

// subclassTimers swaps the ball window procedure for the counting one and
// returns the counter plus a restore func. This is the handle-level-style
// measurement A3 asked for: it counts the WM_TIMER messages the OS actually
// delivers to this window, not a boolean the implementation keeps about
// itself.
func subclassTimers(t *testing.T, b *Ball) (*timerCount, func()) {
	t.Helper()
	if !subclassActive.CompareAndSwap(false, true) {
		t.Fatal("a previous live test left the ball window subclassed")
	}
	cnt := &timerCount{}
	wndProcProbe.Store(cnt)
	done := make(chan uintptr, 1)
	b.sta.PostTask(func() {
		r, _, _ := pSetWindowLongPtrW.Call(uintptr(b.hwnd), gwlpWndProc, probePtr)
		done <- r
	})
	prev := <-done
	if prev == 0 {
		wndProcProbe.Store(nil)
		subclassActive.Store(false)
		t.Fatalf("SetWindowLongPtrW(GWLP_WNDPROC) failed: the timer probe cannot be installed")
	}
	prevWndProc.Store(prev)
	return cnt, func() {
		b.sta.PostTask(func() {
			pSetWindowLongPtrW.Call(uintptr(b.hwnd), gwlpWndProc, prev)
		})
		wndProcProbe.Store(nil)
		subclassActive.Store(false)
	}
}

// enumerateBallWindows lists every top-level window of the ball class, with
// its pid and title. Class-name matching is exactly what the ticket says to
// be careful about, so the result is reported, never quietly ignored.
func enumerateBallWindows() []struct {
	hwnd  windows.HWND
	pid   uint32
	title string
} {
	type hit = struct {
		hwnd  windows.HWND
		pid   uint32
		title string
	}
	var out []hit
	cb := windows.NewCallback(func(hwnd, lParam uintptr) uintptr {
		buf := make([]uint16, 64)
		pGetClassNameW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		if windows.UTF16ToString(buf) != ballClassName {
			return 1
		}
		var pid uint32
		pGetWindowThreadProcessID.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
		tb := make([]uint16, 128)
		n, _, _ := pGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&tb[0])), uintptr(len(tb)))
		out = append(out, hit{hwnd: windows.HWND(hwnd), pid: pid,
			title: windows.UTF16ToString(tb[:n])})
		return 1
	})
	pEnumWindows.Call(cb, 0)
	return out
}

var pEnumWindows = modUser32.NewProc("EnumWindows")

// requireQuietBallDesktop fails the test when the desktop is not ours alone:
// another process's ball (a leftover demo harness) holds hotkey registrations
// and the same window class, and an in-process ball that outlived its test
// means the previous measurement leaked.
func requireQuietBallDesktop(t *testing.T) {
	t.Helper()
	mine := uint32(windows.GetCurrentProcessId())
	for _, w := range enumerateBallWindows() {
		if w.pid != mine {
			t.Fatalf("foreign Wisp ball window on the desktop: hwnd=%v pid=%d title=%q. "+
				"Live measurements must own the desktop: another process holds the window class "+
				"and the global hotkeys, so any number taken now is meaningless. "+
				"Kill pid %d (a leftover balldebug/wisp) and re-run. t.Skip is not an option here.",
				w.hwnd, w.pid, w.title, w.pid)
		}
		t.Fatalf("a Wisp ball window from THIS test process is still alive before the test created "+
			"it: hwnd=%v title=%q - a previous live test leaked its window.", w.hwnd, w.title)
	}
}

// liveBallTitle returns a window title no other run can guess, so every
// liveness assertion can be made about OUR window.
func liveBallTitle(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("WispLive-%d-%s", windows.GetCurrentProcessId(), t.Name())
}

// newLiveBall creates a real ball under the quiet-desktop precondition and
// tears it down. Assertions must use b.DebugHWND(); never FindWindowW by class.
func newLiveBall(t *testing.T, opts Options) *Ball {
	t.Helper()
	requireQuietBallDesktop(t)
	if opts.WindowTitle == "" {
		opts.WindowTitle = liveBallTitle(t)
	}
	b, err := New(opts)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() {
		hwnd := b.DebugHWND()
		b.Close()
		if r, _, _ := pIsWindow.Call(uintptr(hwnd)); r != 0 {
			t.Errorf("ball window %v still alive after Close", hwnd)
		}
		requireQuietBallDesktop(t)
	})
	return b
}

// windowAlive asks Win32, not Go, whether the handle is still a window.
func windowAlive(hwnd windows.HWND) bool {
	r, _, _ := pIsWindow.Call(uintptr(hwnd))
	return r != 0
}

// threadWindows lists the window classes owned by the UI thread. The Go
// runtime keeps a hidden helper window on a thread, so callers assert on the
// BALL class, not on the raw count.
func threadWindowClasses(hwnd windows.HWND) []string {
	tid, _, _ := pGetWindowThreadProcessID.Call(uintptr(hwnd), 0)
	var out []string
	cb := windows.NewCallback(func(h, _ uintptr) uintptr {
		buf := make([]uint16, 64)
		pGetClassNameW.Call(h, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		out = append(out, windows.UTF16ToString(buf))
		return 1
	})
	pEnumThreadWindows.Call(tid, cb, 0)
	return out
}

// ballWindowsOnThread counts WispBallWindow surfaces owned by this thread:
// one Ball must own exactly one.
func ballWindowsOnThread(hwnd windows.HWND) int {
	n := 0
	for _, c := range threadWindowClasses(hwnd) {
		if c == ballClassName {
			n++
		}
	}
	return n
}

// guiResources returns the process USER or GDI object count.
func guiResources(which uintptr) uint64 {
	r, _, _ := pGetGuiResources.Call(uintptr(windows.CurrentProcess()), which)
	return uint64(r)
}

// sendMouse posts a mouse message to the ball window and runs it through the
// real window procedure (the same route a physical click takes after
// WM_NCHITTEST has decided the pixel belongs to the orb).
func sendMouse(b *Ball, msg, wParam, lParam uintptr) uintptr {
	done := make(chan uintptr, 1)
	b.sta.PostTask(func() {
		r, _, _ := pSendMessageW.Call(uintptr(b.hwnd), msg, wParam, lParam)
		done <- r
	})
	return <-done
}

// sendMessage runs an arbitrary message through the live window procedure.
func sendMessage(b *Ball, msg, wParam, lParam uintptr) uintptr {
	return sendMouse(b, msg, wParam, lParam)
}

// lParamPoint packs screen/client coordinates the way mouse messages carry them.
func lParamPoint(x, y int32) uintptr {
	return uintptr(uint32(uint16(x)) | uint32(uint16(y))<<16)
}

// injectHotkey presses a combination through the keyboard driver. It returns
// false when the session refuses injected input (GetAsyncKeyState never sees
// the key), which is an environment limit - NOT a product verdict - so the
// caller reports it loudly instead of pretending the path was exercised.
func injectHotkey(mods uint32, vk uint32) bool {
	press := func(v uintptr, up bool) {
		f := uintptr(0)
		if up {
			f = irKeyup
		}
		pKeybdEvent.Call(v, 0, f, 0)
	}
	var seq []uintptr
	if mods&modControl != 0 {
		seq = append(seq, vkControl)
	}
	if mods&modAlt != 0 {
		seq = append(seq, vkMenu)
	}
	if mods&modShift != 0 {
		seq = append(seq, vkShift)
	}
	if mods&modWin != 0 {
		seq = append(seq, vkLWin)
	}
	seq = append(seq, uintptr(vk))
	for _, v := range seq {
		press(v, false)
	}
	// Did the injection reach the input stream at all?
	arrived := true
	for _, v := range seq {
		if !keyIsDown(v) {
			arrived = false
		}
	}
	for i := len(seq) - 1; i >= 0; i-- {
		press(seq[i], true)
	}
	return arrived
}

func keyIsDown(vk uintptr) bool {
	r, _, _ := pGetAsyncKeyState.Call(vk)
	return uint16(r&0xFFFF)&keyStateDownLo != 0
}

// readState asks the UI thread which state the ball is rendering right now.
// Posts are FIFO, so the answer is ordered against everything posted before it.
func readState(t *testing.T, b *Ball) statemachine.State {
	t.Helper()
	done := make(chan statemachine.State, 1)
	b.sta.PostTask(func() { done <- b.curState })
	select {
	case s := <-done:
		return s
	case <-time.After(5 * time.Second):
		t.Fatal("the UI thread never answered a state read (is the pump alive?)")
		return ""
	}
}

// uiUserObjects is the process USER handle count, the number A6 tracks.
func uiUserObjects() uint64 { return guiResources(grUserObjects) }
