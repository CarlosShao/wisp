//go:build windows

// uipi23 measures the platform side of ticket 23 AC#3's second clause on a real
// elevated window that already exists on this desktop (Task Manager runs at High
// integrity while this probe runs at Medium). It touches nothing in the product:
// every reading is a Win32 call whose answer is printed.
//
// The question it answers: when a Medium process aims input at a High window,
// what does the kernel actually do? Two readings matter and they disagree:
//
//	SendInput's return value counts events INSERTED into the input queue, and it
//	can report full success even though UIPI drops the delivery at the target's
//	thread; AttachThreadInput (the documented way to couple input state across
//	processes) is the call UIPI refuses outright.
//
// Usage: uipi23.exe <hwnd-of-elevated-window>
package main

import (
	"fmt"
	"os"
	"strconv"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32   = windows.NewLazySystemDLL("user32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procGetForegroundWindow  = user32.NewProc("GetForegroundWindow")
	procSetForegroundWindow  = user32.NewProc("SetForegroundWindow")
	procGetGUIThreadInfo     = user32.NewProc("GetGUIThreadInfo")
	procGetWindowThreadID    = user32.NewProc("GetWindowThreadProcessId")
	procAttachThreadInput    = user32.NewProc("AttachThreadInput")
	procGetClassNameW        = user32.NewProc("GetClassNameW")
	procEnumChildWindows     = user32.NewProc("EnumChildWindows")
	procGetWindowRect        = user32.NewProc("GetWindowRect")
	procGetWindowTextLengthW = user32.NewProc("GetWindowTextLengthW")
	procSendInput            = user32.NewProc("SendInput")
	procIsWindowVisible      = user32.NewProc("IsWindowVisible")
	procGetCurrentThreadID   = kernel32.NewProc("GetCurrentThreadId")
)

const (
	processQueryLimitedInfo = 0x1000
	inputKeyboard           = 1
	keyeventfKeyup          = 0x0002
	keyeventfUnicode        = 0x0004
)

type kbInput struct {
	Vk      uint16
	Scan    uint16
	Flags   uint32
	Time    uint32
	Extra   uintptr
	Padding [8]byte
}

type inputUnion struct {
	Type uint32
	_    uint32
	Data kbInput
}

type guiThreadInfo struct {
	Size       uint32
	Flags      uint32
	Active     uintptr
	Focus      uintptr
	Capture    uintptr
	MenuOwner  uintptr
	MoveSize   uintptr
	Caret      uintptr
	RectLeft   int32
	RectTop    int32
	RectRight  int32
	RectBottom int32
}

type rect struct{ Left, Top, Right, Bottom int32 }

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: uipi23.exe <hwnd>")
		os.Exit(2)
	}
	v, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bad hwnd")
		os.Exit(2)
	}
	hwnd := uintptr(v)

	fmt.Printf("target_hwnd=0x%x\n", hwnd)
	fmt.Printf("target_pid=%d\n", pidOf(hwnd))
	fmt.Printf("injector_integrity=%d\n", selfIntegrity())
	shipped, shippedErr := shippedIntegrity(pidOf(hwnd))
	fmt.Printf("target_integrity_shipped_shape=%d (error=%v)\n", shipped, shippedErr)
	fmt.Printf("target_integrity_corrected=%d\n", correctedIntegrity(pidOf(hwnd)))
	fmt.Printf("target_thread=%d\n", threadOf(hwnd))

	thread := threadOf(hwnd)
	myThread, _, _ := procGetCurrentThreadID.Call()
	fmt.Printf("injector_thread=%d\n", myThread)

	before := focusOf(thread)
	fmt.Printf("focus_before=0x%x class=%q\n", before, className(before))
	fmt.Printf("foreground_before=0x%x\n", foreground())

	// The UIPI answer, first reading: coupling to a higher-integrity thread.
	r, _, e := procAttachThreadInput.Call(myThread, uintptr(thread), 1)
	fmt.Printf("attach_thread_input_result=%d err=%v\n", r, e)
	if r != 0 {
		procAttachThreadInput.Call(myThread, uintptr(thread), 0)
	}

	if ok, _, ferr := procSetForegroundWindow.Call(hwnd); ok != 0 {
		fmt.Printf("set_foreground_window=ok err=%v\n", ferr)
	} else {
		fmt.Printf("set_foreground_window=refused err=%v\n", ferr)
	}
	fmt.Printf("foreground_after_focus=0x%x\n", foreground())

	after := focusOf(thread)
	fmt.Printf("focus_after=0x%x class=%q\n", after, className(after))

	// The UIPI answer, second reading: insertion count vs delivery. 'Z' is one
	// character; if UIPI really drops it, the target's own controls cannot show
	// it and the keyboard focus it had is unchanged.
	n, sendErr := sendOne('Z')
	fmt.Printf("sendinput_inserted=%d err=%v\n", n, sendErr)
	fmt.Printf("focus_after_sendinput=0x%x class=%q\n", focusOf(thread), className(focusOf(thread)))

	fmt.Printf("child_controls_begin\n")
	enumChildren(hwnd)
	fmt.Printf("child_controls_end\n")
}

func pidOf(h uintptr) uint32 {
	var pid uint32
	procGetWindowThreadID.Call(h, uintptr(unsafe.Pointer(&pid)))
	return pid
}

func threadOf(h uintptr) uint32 {
	var pid uint32
	t, _, _ := procGetWindowThreadID.Call(h, uintptr(unsafe.Pointer(&pid)))
	return uint32(t)
}

func foreground() uintptr {
	h, _, _ := procGetForegroundWindow.Call()
	return h
}

func focusOf(thread uint32) uintptr {
	if thread == 0 {
		return 0
	}
	info := guiThreadInfo{Size: 0x48}
	r, _, _ := procGetGUIThreadInfo.Call(uintptr(thread), uintptr(unsafe.Pointer(&info)))
	if r == 0 {
		return 0
	}
	return info.Focus
}

func className(h uintptr) string {
	if h == 0 {
		return ""
	}
	var buf [256]uint16
	n, _, _ := procGetClassNameW.Call(h, uintptr(unsafe.Pointer(&buf[0])), 256)
	if int32(n) <= 0 {
		return ""
	}
	return windows.UTF16ToString(buf[:n])
}

func rectOf(h uintptr) string {
	var r rect
	procGetWindowRect.Call(h, uintptr(unsafe.Pointer(&r)))
	return fmt.Sprintf("%d,%d,%d,%d", r.Left, r.Top, r.Right-r.Left, r.Bottom-r.Top)
}

func enumChildren(parent uintptr) {
	cb := windows.NewCallback(func(h, _ uintptr) uintptr {
		if v, _, _ := procIsWindowVisible.Call(h); v == 0 {
			return 1
		}
		l, _, _ := procGetWindowTextLengthW.Call(h)
		fmt.Printf("child=0x%x class=%q rect=%s text_len_cross_process=%d\n", h, className(h), rectOf(h), int32(l))
		return 1
	})
	procEnumChildWindows.Call(parent, cb, 0)
}

func sendOne(r rune) (int, error) {
	ins := []inputUnion{
		{Type: inputKeyboard, Data: kbInput{Scan: uint16(r), Flags: keyeventfUnicode}},
		{Type: inputKeyboard, Data: kbInput{Scan: uint16(r), Flags: keyeventfUnicode | keyeventfKeyup}},
	}
	n, _, e := procSendInput.Call(uintptr(len(ins)), uintptr(unsafe.Pointer(&ins[0])), unsafe.Sizeof(ins[0]))
	if int(n) != len(ins) {
		return int(n), fmt.Errorf("inserted %d/%d: %w", int(n), len(ins), e)
	}
	return int(n), nil
}

// shippedIntegrity replays internal/tools' realProcessElevated exactly: it opens
// the process, then hands that PROCESS handle to GetTokenInformation, which
// wants a TOKEN handle.
func shippedIntegrity(pid uint32) (uint32, error) {
	h, err := windows.OpenProcess(processQueryLimitedInfo, false, pid)
	if err != nil {
		return 0, err
	}
	defer windows.CloseHandle(h)
	return integrityVia(windows.Token(h))
}

func correctedIntegrity(pid uint32) uint32 {
	h, err := windows.OpenProcess(processQueryLimitedInfo, false, pid)
	if err != nil {
		return 0
	}
	defer windows.CloseHandle(h)
	var tok windows.Token
	if err := windows.OpenProcessToken(h, windows.TOKEN_QUERY, &tok); err != nil {
		return 0
	}
	defer tok.Close()
	v, _ := integrityVia(tok)
	return v
}

func selfIntegrity() uint32 {
	var tok windows.Token
	if err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &tok); err != nil {
		return 0
	}
	defer tok.Close()
	v, _ := integrityVia(tok)
	return v
}

func integrityVia(tok windows.Token) (uint32, error) {
	var n uint32
	_ = windows.GetTokenInformation(tok, windows.TokenIntegrityLevel, nil, 0, &n)
	if n == 0 {
		return 0, fmt.Errorf("sizing call returned no size")
	}
	buf := make([]byte, n)
	if err := windows.GetTokenInformation(tok, windows.TokenIntegrityLevel, &buf[0], n, &n); err != nil {
		return 0, err
	}
	label := (*windows.Tokenmandatorylabel)(unsafe.Pointer(&buf[0]))
	count := label.Label.Sid.SubAuthorityCount()
	if count == 0 {
		return 0, fmt.Errorf("no sub-authorities")
	}
	return label.Label.Sid.SubAuthority(uint32(count - 1)), nil
}
