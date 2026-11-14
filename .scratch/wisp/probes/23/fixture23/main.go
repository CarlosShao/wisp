//go:build windows

// wisp23fixture is the scripted target window ticket 23 AC#3 names as an
// alternative to notepad ("a test target app (notepad/scripted window)").
//
// It is a probe-side instrument only: nothing here imports the product. Its one
// job is to be a focused, genuinely editable window that reports, from inside
// its own process, exactly what the keyboard injected into it:
//
//	status.txt  hwnd / pid / edit control / whether it owns the foreground and
//	            the keyboard focus (written once the window is up)
//	dump.txt    the edit control's content length, sha256 and the WM_CHAR /
//	            WM_KEYDOWN counters, rewritten on every delivered character and
//	            on a 50ms timer (so the reader never races the last keystroke)
//	text.txt    the edit control's content verbatim, UTF-8
//
// Usage: wisp23fixture.exe <out-dir> [lifetime-seconds]
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"time"
	"unicode/utf8"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32   = windows.NewLazySystemDLL("user32.dll")
	gdi32    = windows.NewLazySystemDLL("gdi32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procRegisterClassW        = user32.NewProc("RegisterClassW")
	procCreateWindowExW       = user32.NewProc("CreateWindowExW")
	procDefWindowProcW        = user32.NewProc("DefWindowProcW")
	procDestroyWindow         = user32.NewProc("DestroyWindow")
	procPostQuitMessage       = user32.NewProc("PostQuitMessage")
	procShowWindow            = user32.NewProc("ShowWindow")
	procUpdateWindow          = user32.NewProc("UpdateWindow")
	procSetForegroundWindow   = user32.NewProc("SetForegroundWindow")
	procGetForegroundWindow   = user32.NewProc("GetForegroundWindow")
	procBringWindowToTop      = user32.NewProc("BringWindowToTop")
	procSetFocus              = user32.NewProc("SetFocus")
	procGetFocus              = user32.NewProc("GetFocus")
	procSetTimer              = user32.NewProc("SetTimer")
	procKillTimer             = user32.NewProc("KillTimer")
	procGetMessageW           = user32.NewProc("GetMessageW")
	procTranslateMessage      = user32.NewProc("TranslateMessage")
	procDispatchMessageW      = user32.NewProc("DispatchMessageW")
	procGetWindowTextLengthW  = user32.NewProc("GetWindowTextLengthW")
	procGetWindowTextW        = user32.NewProc("GetWindowTextW")
	procSendMessageW          = user32.NewProc("SendMessageW")
	procMoveWindow            = user32.NewProc("MoveWindow")
	procLoadCursorW           = user32.NewProc("LoadCursorW")
	procGetWindowThreadProcID = user32.NewProc("GetWindowThreadProcessId")
	procAttachThreadInput     = user32.NewProc("AttachThreadInput")
	procGetWindowLongPtrW     = user32.NewProc("GetWindowLongPtrW")
	procSetWindowLongPtrW     = user32.NewProc("SetWindowLongPtrW")
	procCallWindowProcW       = user32.NewProc("CallWindowProcW")
	procGetStockObject        = gdi32.NewProc("GetStockObject")
	procGetModuleHandleW      = kernel32.NewProc("GetModuleHandleW")
	procGetCurrentThreadID    = kernel32.NewProc("GetCurrentThreadId")
	procSwitchToThisWindow    = user32.NewProc("SwitchToThisWindow")
	procAllowSetForeground    = user32.NewProc("AllowSetForegroundWindow")
)

const (
	wsOverlappedWindow = 0x00CF0000
	wsVisible          = 0x10000000
	wsChild            = 0x40000000
	wsClipChildren     = 0x02000000
	wsVScroll          = 0x00200000
	wsTabStop          = 0x00010000
	wsExClientEdge     = 0x00000200
	esMultiline        = 0x0004
	esAutoVScroll      = 0x0040
	esWantReturn       = 0x1000

	wmDestroy      = 0x0002
	wmClose        = 0x0010
	wmSetFocus     = 0x0007
	wmSize         = 0x0005
	wmKeyDown      = 0x0100
	wmChar         = 0x0102
	wmTimer        = 0x0113
	wmSetText      = 0x000C
	wmSetFont      = 0x0030
	gwlpWndProc    = ^uintptr(3)
	swShowNormal   = 1
	swRestore      = 9
	defaultGUIFont = 17

	timerDump    = 1
	timerLife    = 2
	dumpPeriodMS = 50

	windowTitle = "WISP23-FIXTURE input.type target"
	windowClass = "Wisp23FixtureWindow"
)

type wndClassW struct {
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   uintptr
	Icon       uintptr
	Cursor     uintptr
	Background uintptr
	MenuName   *uint16
	ClassName  *uint16
}

type msgT struct {
	Hwnd    uintptr
	Message uint32
	_       uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	PtX     int32
	PtY     int32
	Private uint32
}

var (
	topHwnd  uintptr
	editHwnd uintptr
	outDir   string
	started  time.Time

	origEditProc uintptr
	editProcCb   uintptr
	wndProcCb    uintptr
	charsSeen    int64
	keyDownsSeen int64
	dumps        int64
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: wisp23fixture <out-dir> [lifetime-seconds]")
		os.Exit(2)
	}
	outDir = os.Args[1]
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "mkdir out:", err)
		os.Exit(1)
	}
	lifetime := 120
	if len(os.Args) > 2 {
		if n, err := strconv.Atoi(os.Args[2]); err == nil && n > 0 {
			lifetime = n
		}
	}
	started = time.Now()

	inst, _, _ := procGetModuleHandleW.Call(0)
	cursor, _, _ := procLoadCursorW.Call(0, uintptr(32512)) // IDC_ARROW
	cls, _ := windows.UTF16PtrFromString(windowClass)
	title, _ := windows.UTF16PtrFromString(windowTitle)

	wndProcCb = windows.NewCallback(wndProc)
	wc := wndClassW{
		WndProc:    wndProcCb,
		Instance:   inst,
		Cursor:     cursor,
		Background: 5, // COLOR_WINDOW + 1
		ClassName:  cls,
	}
	if r, _, e := procRegisterClassW.Call(uintptr(unsafe.Pointer(&wc))); r == 0 {
		fmt.Fprintln(os.Stderr, "RegisterClassW:", e)
		os.Exit(1)
	}

	topHwnd, _, _ = procCreateWindowExW.Call(
		0, uintptr(unsafe.Pointer(cls)), uintptr(unsafe.Pointer(title)),
		wsOverlappedWindow|wsVisible|wsClipChildren,
		80, 80, 900, 600, 0, 0, inst, 0)
	if topHwnd == 0 {
		fmt.Fprintln(os.Stderr, "CreateWindowExW top failed")
		os.Exit(1)
	}

	editCls, _ := windows.UTF16PtrFromString("EDIT")
	editHwnd, _, _ = procCreateWindowExW.Call(
		wsExClientEdge, uintptr(unsafe.Pointer(editCls)), 0,
		wsChild|wsVisible|wsVScroll|wsTabStop|esMultiline|esAutoVScroll|esWantReturn,
		8, 8, 870, 540, topHwnd, 0, inst, 0)
	if editHwnd == 0 {
		fmt.Fprintln(os.Stderr, "CreateWindowExW edit failed")
		os.Exit(1)
	}
	font, _, _ := procGetStockObject.Call(defaultGUIFont)
	procSendMessageW.Call(editHwnd, wmSetFont, font, 1)

	// Subclass the edit control: the WM_CHAR/WM_KEYDOWN counters are the
	// independent "what actually arrived" numbers the reader compares the
	// buffer length against.
	editProcCb = windows.NewCallback(editProc)
	origEditProc, _, _ = procGetWindowLongPtrW.Call(editHwnd, gwlpWndProc)
	procSetWindowLongPtrW.Call(editHwnd, gwlpWndProc, editProcCb)

	procShowWindow.Call(topHwnd, swShowNormal)
	procUpdateWindow.Call(topHwnd)

	method := forceForeground(topHwnd)
	procSetFocus.Call(editHwnd)

	procSetTimer.Call(topHwnd, timerDump, dumpPeriodMS, 0)
	procSetTimer.Call(topHwnd, timerLife, uintptr(lifetime*1000), 0)

	writeStatus(method)
	dump()

	var m msgT
	for {
		r, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) <= 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
	dump()
	writeStatus("exited")
}

func forceForeground(top uintptr) string {
	fg, _, _ := procGetForegroundWindow.Call()
	if fg == top {
		return "already-foreground"
	}
	procBringWindowToTop.Call(top)
	if ok, _, _ := procSetForegroundWindow.Call(top); ok != 0 {
		fg, _, _ = procGetForegroundWindow.Call()
		if fg == top {
			return "SetForegroundWindow"
		}
	}
	var fgPID uint32
	fg, _, _ = procGetForegroundWindow.Call()
	fgThread, _, _ := procGetWindowThreadProcID.Call(fg, uintptr(unsafe.Pointer(&fgPID)))
	myThread, _, _ := procGetCurrentThreadID.Call()
	if fgThread != 0 && fgThread != myThread {
		procAttachThreadInput.Call(myThread, fgThread, 1)
		procShowWindow.Call(top, swRestore)
		procSetForegroundWindow.Call(top)
		procAttachThreadInput.Call(myThread, fgThread, 0)
		if fg, _, _ = procGetForegroundWindow.Call(); fg == top {
			return "AttachThreadInput+SetForegroundWindow"
		}
	}
	return "failed"
}

func wndProc(hwnd, msg, wparam, lparam uintptr) uintptr {
	switch msg {
	case wmClose:
		procDestroyWindow.Call(hwnd)
		return 0
	case wmDestroy:
		procPostQuitMessage.Call(0)
		return 0
	case wmSetFocus:
		procSetFocus.Call(editHwnd)
		return 0
	case wmSize:
		w := lparam & 0xFFFF
		h := (lparam >> 16) & 0xFFFF
		procMoveWindow.Call(editHwnd, 0, 0, w, h, 1)
		return 0
	case wmTimer:
		if wparam == timerLife {
			procDestroyWindow.Call(hwnd)
			return 0
		}
		dump()
		return 0
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, msg, wparam, lparam)
	return r
}

func editProc(hwnd, msg, wparam, lparam uintptr) uintptr {
	switch msg {
	case wmChar:
		atomic.AddInt64(&charsSeen, 1)
	case wmKeyDown:
		atomic.AddInt64(&keyDownsSeen, 1)
	}
	r, _, _ := procCallWindowProcW.Call(origEditProc, hwnd, msg, wparam, lparam)
	if msg == wmChar {
		dump()
	}
	return r
}

func editText() string {
	n, _, _ := procGetWindowTextLengthW.Call(editHwnd)
	buf := make([]uint16, int(n)+1)
	if len(buf) == 0 {
		return ""
	}
	procGetWindowTextW.Call(editHwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return windows.UTF16ToString(buf)
}

// drainCommand lets the reader reset the target between injections without any
// cross-process handle: writing "clear" into <out>/cmd.txt empties the edit
// control and zeroes the message counters.
func drainCommand() {
	cmdPath := filepath.Join(outDir, "cmd.txt")
	b, err := os.ReadFile(cmdPath)
	if err != nil {
		return
	}
	if string(bytesTrim(b)) != "clear" {
		return
	}
	empty, _ := windows.UTF16PtrFromString("")
	procSendMessageW.Call(editHwnd, wmSetText, 0, uintptr(unsafe.Pointer(empty)))
	atomic.StoreInt64(&charsSeen, 0)
	atomic.StoreInt64(&keyDownsSeen, 0)
	_ = os.Remove(cmdPath)
}

func bytesTrim(b []byte) string {
	s := string(b)
	for len(s) > 0 && (s[len(s)-1] == '\r' || s[len(s)-1] == '\n' || s[len(s)-1] == ' ') {
		s = s[:len(s)-1]
	}
	for len(s) > 0 && s[0] == ' ' {
		s = s[1:]
	}
	return s
}

func dump() {
	drainCommand()
	text := editText()
	sum := sha256.Sum256([]byte(text))
	fg, _, _ := procGetForegroundWindow.Call()
	focus, _, _ := procGetFocus.Call()
	isFg := 0
	if fg == topHwnd {
		isFg = 1
	}
	line := fmt.Sprintf(
		"ts=%s top=%d edit=%d pid=%d chars=%d runes=%d utf16_units=%d sha256=%s wm_char=%d wm_keydown=%d foreground_hwnd=%d is_foreground=%d focus_hwnd=%d dumps=%d\n",
		time.Now().Format(time.RFC3339Nano), topHwnd, editHwnd, os.Getpid(),
		len([]rune(text)), utf8.RuneCountInString(text), len(windows.StringToUTF16(text))-1,
		hex.EncodeToString(sum[:]), atomic.LoadInt64(&charsSeen), atomic.LoadInt64(&keyDownsSeen),
		fg, isFg, focus, atomic.AddInt64(&dumps, 1))
	_ = os.WriteFile(filepath.Join(outDir, "dump.txt"), []byte(line), 0o644)
	_ = os.WriteFile(filepath.Join(outDir, "text.txt"), []byte(text), 0o644)
}

func writeStatus(method string) {
	fg, _, _ := procGetForegroundWindow.Call()
	focus, _, _ := procGetFocus.Call()
	line := fmt.Sprintf(
		"ts=%s hwnd=%d edit=%d pid=%d title=%q class=%q foreground_method=%s is_foreground=%d focus_hwnd=%d elapsed_ms=%d\n",
		time.Now().Format(time.RFC3339Nano), topHwnd, editHwnd, os.Getpid(), windowTitle, windowClass,
		method, boolInt(fg == topHwnd), focus, time.Since(started).Milliseconds())
	_ = os.WriteFile(filepath.Join(outDir, "status.txt"), []byte(line), 0o644)
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
