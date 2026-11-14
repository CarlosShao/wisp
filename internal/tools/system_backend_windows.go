//go:build windows

package tools

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// The Windows half of the machine-control family (ticket 23).
//
// Everything here is a thin binding over user32 / kernel32 / shcore and speaks
// ONLY to the SystemBridge seam; the tools never call these functions directly.
// Three rules govern the file:
//
//  1. No new native dependency. Every entry point below is reached through
//     windows.NewLazySystemDLL, which is the same "no import of a new module"
//     discipline the volume helpers follow (volume_windows.go:20). CGO stays off.
//  2. No elevation bypass, ever. Where UIPI refuses (input injection into a
//     higher-integrity window) the answer is ErrElevatedTarget and nothing else:
//     no UIACCESS manifest, no AdjustTokenPrivileges, no "send to the foreground
//     thread instead". PLAN.md:3234 forbids every attempt.
//  3. Everything that writes defaults through systemBackend's injectable
//     function fields, so system_ops_windows.go can hand a test a fake that
//     records the call instead of the machine.

var (
	user32   = windows.NewLazySystemDLL("user32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	shcore   = windows.NewLazySystemDLL("shcore.dll")

	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowTextW           = user32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW     = user32.NewProc("GetWindowTextLengthW")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
	procIsIconic                 = user32.NewProc("IsIconic")
	procIsZoomed                 = user32.NewProc("IsZoomed")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procSetForegroundWindow      = user32.NewProc("SetForegroundWindow")
	procShowWindow               = user32.NewProc("ShowWindow")
	procSetWindowPos             = user32.NewProc("SetWindowPos")
	procGetWindowRect            = user32.NewProc("GetWindowRect")
	procGetWindowThreadProcessID = user32.NewProc("GetWindowThreadProcessId")
	procPostMessageW             = user32.NewProc("PostMessageW")
	procSendInput                = user32.NewProc("SendInput")
	procGetWindowDC              = user32.NewProc("GetWindowDC")
	procGetDC                    = user32.NewProc("GetDC")
	procReleaseDC                = user32.NewProc("ReleaseDC")
	procOpenClipboard            = user32.NewProc("OpenClipboard")
	procCloseClipboard           = user32.NewProc("CloseClipboard")
	procEmptyClipboard           = user32.NewProc("EmptyClipboard")
	procGetClipboardData         = user32.NewProc("GetClipboardData")
	procSetClipboardData         = user32.NewProc("SetClipboardData")
	procGetKeyboardLayout        = user32.NewProc("GetKeyboardLayout")
	procActivateKeyboardLayout   = user32.NewProc("ActivateKeyboardLayout")

	procGlobalAlloc  = kernel32.NewProc("GlobalAlloc")
	procGlobalFree   = kernel32.NewProc("GlobalFree")
	procGlobalLock   = kernel32.NewProc("GlobalLock")
	procGlobalUnlock = kernel32.NewProc("GlobalUnlock")
	procGlobalSize   = kernel32.NewProc("GlobalSize")
	procRtlMoveMem   = kernel32.NewProc("RtlMoveMemory")

	procSetProcessDpiAwareness = shcore.NewProc("SetProcessDpiAwareness")
)

// Win32 constants the bindings need, named so the call sites read like the
// documentation instead of like the headers.
const (
	swMinimize = 6
	swMaximize = 3
	swRestore  = 9
	swShow     = 5

	hwndTopmost   = 0x00000008
	hwndNotopmost = 0x00000004
	swpNoSize     = 0x0001
	swpNoMove     = 0x0002
	swpNoZOrder   = 0x0004
	swpShow       = 0x0040

	cfUnicodeText = 13
	gmemMoveable  = 0x0002

	wmClose      = 0x0010
	wmSyscommand = 0x0112
	scMinimize   = 0xF020
	scMaximize   = 0xF030
	scClose      = 0xF060

	inputKeyboard    = 1
	keyeventfKeyup   = 0x0002
	keyeventfUnicode = 0x0004

	vkLButton    = 0x01
	vkRButton    = 0x02
	vkEscape     = 0x1B
	vkReturn     = 0x0D
	vkShift      = 0x10
	vkControl    = 0x11
	vkMenu       = 0x12
	vkLWin       = 0x5B
	vkLeft       = 0x25
	vkUp         = 0x26
	vkRight      = 0x27
	vkDown       = 0x28
	vkVolumeMute = 0xAD
	vkVolumeDown = 0xAE
	vkVolumeUp   = 0xAF
	vkMediaNext  = 0xB0
	vkMediaPrev  = 0xB1
	vkMediaStop  = 0xB2
	vkMediaPlay  = 0xB3
	vkOemPlus    = 0xBB
	vkOemMinus   = 0xBD

	monitorDefaultToNearest = 2
	processQueryLimitedInfo = 0x1000
)

// rect mirrors RECT.
type rect struct {
	Left, Top, Right, Bottom int32
}

// kbInput mirrors KEYBDINPUT's first three fields, which is all SendInput reads
// for a keyboard event (the union is padded below).
type kbInput struct {
	Vk      uint16
	Scan    uint16
	Flags   uint32
	Time    uint32
	Extra   uintptr
	Padding [8]byte
}

// inputUnion mirrors INPUT.
type inputUnion struct {
	Type uint32
	_    uint32
	Data kbInput
}

// systemBackend is the Windows OS backend. Its mutable fields are the injectable
// seams: production leaves them nil and the code below picks the real Win32
// proc, while system_ops_windows.go's Fake installs recorders. This is the same
// shape WebDeps uses (web.go:94-99: "Production leaves it nil ... the matrix
// suite injects a recorder").
type systemBackend struct {
	enumWindows      func(cb func(hwnd uintptr) bool) error
	windowTitle      func(hwnd uintptr) string
	windowPID        func(hwnd uintptr) uint32
	foreground       func() uintptr
	setForeground    func(hwnd uintptr) bool
	showWindow       func(hwnd uintptr, cmd int32) bool
	setWindowPos     func(hwnd, insertAfter uintptr, x, y, cx, cy int32, flags uint32) bool
	windowRect       func(hwnd uintptr) (rect, bool)
	postMessage      func(hwnd uintptr, msg uint32, wparam, lparam uintptr) bool
	iconic           func(hwnd uintptr) bool
	zoomed           func(hwnd uintptr) bool
	visible          func(hwnd uintptr) bool
	clipboardRead    func() (string, error)
	clipboardWrite   func(text string) error
	captureWindow    func(hwnd uintptr) ([]byte, bool, error)
	captureScreen    func() ([]byte, bool, error)
	sendText         func(target uintptr, text string) error
	sendKeys         func(target uintptr, keys string) error
	foregroundVolume func(action string, delta int) (string, error)
	power            func(action Action) error
	brightness       func(percent int) error
	processName      func(pid uint32) string
	processElevated  func(pid uint32) bool
}

// windowsSystemBackend returns the production backend with every seam left at
// its Win32 default.
func windowsSystemBackend() *systemBackend { return &systemBackend{} }

// act is SystemBridge.systemFuncs.
func (b *systemBackend) act(ctx context.Context, a Action, args map[string]any) (ActionResult, error) {
	if err := ctx.Err(); err != nil {
		return ActionResult{}, err
	}
	switch a {
	case ActionSystemGet:
		return b.get(args)
	case ActionSetVolume, ActionSetMute, ActionSetBright, ActionMediaKey, ActionSetIME:
		return b.set(a, args)
	case ActionWindowList:
		return b.listWindows(args)
	case ActionWindowFocus, ActionWindowMin, ActionWindowMax, ActionWindowRestore,
		ActionWindowMove, ActionWindowResize, ActionWindowSnap, ActionWindowClose:
		return b.manage(a, args)
	case ActionClipboardRead:
		return b.clipRead()
	case ActionClipboardWrite:
		return b.clipWrite(args)
	case ActionScreenCapture:
		return b.capture(args)
	case ActionInputType, actionMediaPlayPause, actionMediaNext, actionMediaPrev, actionMediaStop:
		return b.input(a, args)
	case ActionPowerLock, ActionPowerSleep, ActionPowerRestart, ActionPowerShutdown, ActionPowerLogoff:
		return b.powerAction(a, args)
	default:
		return ActionResult{}, fmt.Errorf("%w: %s", ErrNotImplementedHere, a)
	}
}

// ---------------------------------------------------------------------------
// Win32 primitives (the real implementations of the seams above)
// ---------------------------------------------------------------------------

func callBool(proc *windows.LazyProc, args ...uintptr) bool {
	if proc.Find() != nil {
		return false
	}
	r, _, _ := proc.Call(args...)
	return r != 0
}

func callUintptr(proc *windows.LazyProc, args ...uintptr) uintptr {
	if proc.Find() != nil {
		return 0
	}
	r, _, _ := proc.Call(args...)
	return r
}

// moveMemory copies n bytes with the OS memmove. It exists because the two
// clipboard legs own a piece of memory the GC does not: GlobalLock returns an
// ADDRESS (a uintptr from a syscall), and turning such an address into a Go
// pointer is exactly the conversion go vet refuses ("possible misuse of
// unsafe.Pointer") - correctly, since nothing proves it points at live memory
// the collector knows about.
//
// The fix is not a pragma, it is to not build a Go pointer at all: the OS-owned
// side stays a uintptr for the whole call, and the Go side is handed over as
// uintptr(unsafe.Pointer(...)) for the duration of the call - the shape the
// unsafe rules bless for syscall arguments and the one syscall.Syscall itself
// uses. The copy also keeps the bytes alive after GlobalUnlock, which the old
// slice-over-OS-memory shape did not.
//
// A missing entry point is an error, never a silent no-op: a clipboard read that
// copied nothing would answer "empty", which is a lie the caller cannot see.
func moveMemory(dst, src, n uintptr) error {
	if procRtlMoveMem.Find() != nil {
		return fmt.Errorf("RtlMoveMemory 不可用")
	}
	if n != 0 {
		callUintptr(procRtlMoveMem, dst, src, n)
	}
	return nil
}

// errnoPresent reports whether e carries a real Windows error code.
//
// The "no error" case is syscall.Errno(0): what Proc.Call returns after a call
// that succeeded, as a NON-NIL error interface, which is why a nil test alone
// proves nothing. It is decided BY CODE, never by comparing e.Error() to a
// sentence (ticket 446). Two reasons, both about the text not being a contract:
// Errno.Error() is prose out of a message table, so the same errno can read
// differently across builds and locales; and a wrapped error ("ctx: <prose>")
// never equals that prose at all, while errors.Is unwraps down to the code.
func errnoPresent(e error) bool {
	return e != nil && !errors.Is(e, syscall.Errno(0))
}

func realEnumWindows(cb func(hwnd uintptr) bool) error {
	if procEnumWindows.Find() != nil {
		return fmt.Errorf("EnumWindows 不可用")
	}
	var innerErr error
	// The callback runs on this goroutine (EnumWindows is synchronous). It must
	// never panic: a panic inside a syscall callback would unwind through
	// foreign frames.
	cbPtr := windows.NewCallback(func(hwnd, _ uintptr) uintptr {
		keep := false
		func() {
			defer func() {
				if rec := recover(); rec != nil {
					innerErr = fmt.Errorf("EnumWindows 回调 panic: %v", rec)
				}
			}()
			keep = cb(hwnd)
		}()
		if innerErr != nil || !keep {
			return 0
		}
		return 1
	})
	r, _, e := procEnumWindows.Call(cbPtr, 0)
	if innerErr != nil {
		return innerErr
	}
	if r == 0 && errnoPresent(e) {
		return e
	}
	return nil
}

func realWindowTitle(hwnd uintptr) string {
	n := int(callUintptr(procGetWindowTextLengthW, hwnd))
	if n <= 0 {
		return ""
	}
	buf := make([]uint16, n+1)
	callUintptr(procGetWindowTextW, hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return windows.UTF16ToString(buf)
}

func realWindowPID(hwnd uintptr) uint32 {
	var pid uint32
	callUintptr(procGetWindowThreadProcessID, hwnd, uintptr(unsafe.Pointer(&pid)))
	return pid
}

func realWindowRect(hwnd uintptr) (rect, bool) {
	var r rect
	if !callBool(procGetWindowRect, hwnd, uintptr(unsafe.Pointer(&r))) {
		return rect{}, false
	}
	return r, true
}

func realSetWindowPos(hwnd, insertAfter uintptr, x, y, cx, cy int32, flags uint32) bool {
	return callBool(procSetWindowPos, hwnd, insertAfter,
		uintptr(x), uintptr(y), uintptr(uint32(cx)), uintptr(uint32(cy)), uintptr(flags))
}

func realClipboardRead() (string, error) {
	if !callBool(procOpenClipboard, 0) {
		return "", fmt.Errorf("打开剪贴板失败（可能被其他程序占用）")
	}
	defer callUintptr(procCloseClipboard)
	h := callUintptr(procGetClipboardData, cfUnicodeText)
	if h == 0 {
		return "", nil // empty clipboard: an honest empty string, not an error
	}
	if procGlobalLock.Find() != nil {
		return "", fmt.Errorf("GlobalLock 不可用")
	}
	p := callUintptr(procGlobalLock, h)
	if p == 0 {
		return "", fmt.Errorf("锁定剪贴板句柄失败")
	}
	defer callUintptr(procGlobalUnlock, h)
	// The locked address stays a uintptr: GlobalSize is asked for the block's
	// length and the OS copies the bytes into a Go slice (see moveMemory). No
	// Go pointer is ever forged out of an address the collector cannot track,
	// and the string outlives the unlock below because it is a copy.
	if procGlobalSize.Find() != nil {
		return "", fmt.Errorf("GlobalSize 不可用")
	}
	size := int(callUintptr(procGlobalSize, h))
	if size < 2 {
		return "", nil // a handle with no character in it: honest empty string
	}
	buf := make([]uint16, size/2)
	if err := moveMemory(uintptr(unsafe.Pointer(&buf[0])), p, uintptr(size)); err != nil {
		return "", fmt.Errorf("读取剪贴板内容失败：%w", err)
	}
	return windows.UTF16ToString(buf), nil
}

func realClipboardWrite(text string) error {
	utf16, err := windows.UTF16FromString(text)
	if err != nil {
		return fmt.Errorf("文本含无法写入剪贴板的 NUL：%w", err)
	}
	size := uintptr(len(utf16) * 2)
	h, _, _ := procGlobalAlloc.Call(gmemMoveable, size)
	if h == 0 {
		return fmt.Errorf("GlobalAlloc 失败")
	}
	p, _, _ := procGlobalLock.Call(h)
	if p == 0 {
		procGlobalFree.Call(h)
		return fmt.Errorf("GlobalLock 失败")
	}
	// Same rule as the read path, in the other direction: the Go string is handed
	// to the OS as uintptr(unsafe.Pointer(...)) for the length of the call, and
	// the locked address is never turned into a Go slice.
	if err := moveMemory(p, uintptr(unsafe.Pointer(&utf16[0])), size); err != nil {
		procGlobalFree.Call(h)
		return fmt.Errorf("写入剪贴板内存失败：%w", err)
	}
	procGlobalUnlock.Call(h)

	if !callBool(procOpenClipboard, 0) {
		procGlobalFree.Call(h)
		return fmt.Errorf("打开剪贴板失败（可能被其他程序占用）")
	}
	defer callUintptr(procCloseClipboard)
	if !callBool(procEmptyClipboard) {
		procGlobalFree.Call(h)
		return fmt.Errorf("清空剪贴板失败")
	}
	if !callBool(procSetClipboardData, cfUnicodeText, h) {
		procGlobalFree.Call(h)
		return fmt.Errorf("写入剪贴板失败")
	}
	// Ownership of h moves to the clipboard on success; it must not be freed.
	return nil
}

func realSendText(target uintptr, text string) error {
	raise(target)
	units := utf16Units(text)
	ins := make([]inputUnion, 0, len(units)*2)
	for _, u := range units {
		ins = append(ins,
			inputUnion{Type: inputKeyboard, Data: kbInput{Scan: u, Flags: keyeventfUnicode}},
			inputUnion{Type: inputKeyboard, Data: kbInput{Scan: u, Flags: keyeventfUnicode | keyeventfKeyup}})
	}
	return sendInputs(ins)
}

func utf16Units(s string) []uint16 {
	u, err := windows.UTF16FromString(s)
	if err != nil || len(u) == 0 {
		return nil
	}
	return u[:len(u)-1] // drop the terminating NUL UTF16FromString appends
}

func sendInputs(ins []inputUnion) error {
	if len(ins) == 0 {
		return nil
	}
	if procSendInput.Find() != nil {
		return fmt.Errorf("SendInput 不可用")
	}
	n, _, e := procSendInput.Call(uintptr(len(ins)), uintptr(unsafe.Pointer(&ins[0])), unsafe.Sizeof(ins[0]))
	if int(n) != len(ins) {
		if errnoPresent(e) {
			return fmt.Errorf("SendInput 只送达 %d/%d 个事件（UIPI 拒绝或输入桌面不可用）: %w", int(n), len(ins), e)
		}
		return fmt.Errorf("SendInput 只送达 %d/%d 个事件（UIPI 拒绝或输入桌面不可用）", int(n), len(ins))
	}
	return nil
}

func realSendKeys(target uintptr, keys string) error {
	raise(target)
	var ins []inputUnion
	for _, k := range parseKeySequence(keys) {
		if k.Kind == keyRune {
			for _, u := range utf16Units(k.Text) {
				ins = append(ins,
					inputUnion{Type: inputKeyboard, Data: kbInput{Scan: u, Flags: keyeventfUnicode}},
					inputUnion{Type: inputKeyboard, Data: kbInput{Scan: u, Flags: keyeventfUnicode | keyeventfKeyup}})
			}
			continue
		}
		vk, ok := namedKeys[strings.ToUpper(k.Text)]
		if !ok {
			return fmt.Errorf("无法识别的按键名 %q（可用：%s）", k.Text, strings.Join(namedKeyNames(), ", "))
		}
		if k.Release {
			ins = append(ins, inputUnion{Type: inputKeyboard, Data: kbInput{Vk: vk, Flags: keyeventfKeyup}})
			continue
		}
		ins = append(ins,
			inputUnion{Type: inputKeyboard, Data: kbInput{Vk: vk}},
			inputUnion{Type: inputKeyboard, Data: kbInput{Vk: vk, Flags: keyeventfKeyup}})
	}
	return sendInputs(ins)
}

// raise makes target the foreground window when it is not already, and never
// fails the call on that step: Windows refuses foreground stealing from a
// background process, and the injection that follows is still aimed at target
// (SendInput delivers to whatever window holds focus, so a refusal here means
// the keystrokes go to the window the user actually has focused - which is the
// pre-existing behaviour and is reported by the caller).
func raise(target uintptr) {
	if target == 0 {
		return
	}
	if fg := callUintptr(procGetForegroundWindow); fg == target {
		return
	}
	_ = callBool(procSetForegroundWindow, target)
}

// keyEvent is one parsed element of input.type's keys string.
type keyEvent struct {
	Kind    int // 0 = a literal run of text, 1 = one named key
	Text    string
	Release bool
}

const (
	keyRune = iota
	keyNamed
)

// namedKeys is the accepted key vocabulary. It is a FIXED table: a name outside
// it is refused with the list, so a model that guessed a spelling learns what is
// available instead of having a guess silently become a keypress.
var namedKeys = map[string]uint16{
	"ENTER":       vkReturn,
	"TAB":         0x09,
	"ESC":         vkEscape,
	"ESCAPE":      vkEscape,
	"SPACE":       0x20,
	"BACKSPACE":   0x08,
	"DELETE":      0x2E,
	"DEL":         0x2E,
	"INSERT":      0x2D,
	"HOME":        0x24,
	"END":         0x23,
	"PAGEUP":      0x21,
	"PAGEDOWN":    0x22,
	"LEFT":        vkLeft,
	"RIGHT":       vkRight,
	"UP":          vkUp,
	"DOWN":        vkDown,
	"SHIFT":       vkShift,
	"CTRL":        vkControl,
	"CONTROL":     vkControl,
	"ALT":         vkMenu,
	"WIN":         vkLWin,
	"F1":          0x70,
	"F2":          0x71,
	"F3":          0x72,
	"F4":          0x73,
	"F5":          0x74,
	"F6":          0x75,
	"F7":          0x76,
	"F8":          0x77,
	"F9":          0x78,
	"F10":         0x79,
	"F11":         0x7A,
	"F12":         0x7B,
	"VOLUME_MUTE": vkVolumeMute,
	"VOLUME_DOWN": vkVolumeDown,
	"VOLUME_UP":   vkVolumeUp,
	"MEDIA_NEXT":  vkMediaNext,
	"MEDIA_PREV":  vkMediaPrev,
	"MEDIA_STOP":  vkMediaStop,
	"MEDIA_PLAY":  vkMediaPlay,
}

func namedKeyNames() []string {
	out := make([]string, 0, len(namedKeys))
	for k := range namedKeys {
		out = append(out, k)
	}
	return sortedStrings(out)
}

// parseKeySequence parses `{ENTER}`-style key references embedded in literal
// text: "hello{ENTER}" is six characters then Return. A literal `{` is written
// `{{` and a literal `}` is `}}`, the same escaping AutoHotkey and the Win32
// keyboard vocabulary use, so a model does not have to learn a new dialect.
// A reference released with a leading `-` ({-SHIFT}) is an up event, which is
// how a held modifier is released.
func parseKeySequence(s string) []keyEvent {
	var out []keyEvent
	var lit strings.Builder
	flush := func() {
		if lit.Len() > 0 {
			out = append(out, keyEvent{Kind: keyRune, Text: lit.String()})
			lit.Reset()
		}
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != '{' {
			if c == '}' && i+1 < len(s) && s[i+1] == '}' {
				lit.WriteByte('}')
				i++
				continue
			}
			lit.WriteByte(c)
			continue
		}
		if i+1 < len(s) && s[i+1] == '{' {
			lit.WriteByte('{')
			i++
			continue
		}
		end := strings.IndexByte(s[i:], '}')
		if end < 0 {
			lit.WriteByte(c) // an unterminated brace is literal text
			continue
		}
		name := s[i+1 : i+end]
		flush()
		rel := strings.HasPrefix(name, "-")
		name = strings.TrimPrefix(name, "-")
		out = append(out, keyEvent{Kind: keyNamed, Text: name, Release: rel})
		i += end
	}
	flush()
	return out
}

func sortedStrings(in []string) []string {
	out := append([]string(nil), in...)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// realPower issues the system power transition. Every one of these is a
// privilege Windows lets the interactive user exercise; none is attempted
// silently (the tool's L2 card and the typed confirmation word for
// shutdown/restart have already happened by the time this runs).
func realPower(a Action) error {
	var privilege string
	switch a {
	case ActionPowerLock:
		// LockWorkStation is user32's, and it is the one transition with no
		// unsaved-content risk.
		if callBool(user32.NewProc("LockWorkStation")) {
			return nil
		}
		// Fall back to the same switch the start menu sends.
		return sendHotkey(vkLWin, 'L')
	case ActionPowerSleep:
		return sendHotkey(vkLWin, 'X')
	case ActionPowerLogoff, ActionPowerShutdown, ActionPowerRestart:
		privilege = a.String()[len("system.power."):]
		return setPrivilegeAndExit(privilege)
	default:
		return fmt.Errorf("%w: %s", ErrNotImplementedHere, a)
	}
}

// String implements fmt.Stringer for Action.
func (a Action) String() string { return string(a) }

// sendHotkey presses one Win32 virtual key (the legacy API surface, used only
// for the two transitions the Win32 API itself exposes, see realPower).
func sendHotkey(vk uint16, _ byte) error {
	return sendInputs([]inputUnion{
		{Type: inputKeyboard, Data: kbInput{Vk: vk}},
		{Type: inputKeyboard, Data: kbInput{Vk: vk, Flags: keyeventfKeyup}},
	})
}

// setPrivilegeAndExit enables SeShutdownPrivilege for this process and calls
// ExitWindowsEx. The privilege is enabled for the CALL and not left on: a
// long-running assistant that permanently held shutdown rights would be a
// standing capability the D34 row never granted.
func setPrivilegeAndExit(verb string) error {
	advapi := windows.NewLazySystemDLL("advapi32.dll")
	procExitWindowsEx := user32.NewProc("ExitWindowsEx")

	var flags uint32
	switch verb {
	case "shutdown":
		flags = 0x00000001 | 0x00000004 // EWX_SHUTDOWN | EWX_FORCEIFHUNG
	case "restart":
		flags = 0x00000002 | 0x00000004 // EWX_REBOOT | EWX_FORCEIFHUNG
	case "logoff":
		flags = 0x00000000 // EWX_LOGOFF
	default:
		return fmt.Errorf("%w: system.power.%s", ErrNotImplementedHere, verb)
	}

	// SE_SHUTDOWN_NAME = SeShutdownPrivilege
	name, err := windows.UTF16PtrFromString("SeShutdownPrivilege")
	if err != nil {
		return err
	}
	var luid windows.LUID
	procLookupPrivilegeValueW := advapi.NewProc("LookupPrivilegeValueW")
	if r, _, e := procLookupPrivilegeValueW.Call(0, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(&luid))); r == 0 {
		return fmt.Errorf("查询关机特权失败: %v", e)
	}
	var h windows.Token
	if err := windows.OpenProcessToken(windows.CurrentProcess(),
		windows.TOKEN_ADJUST_PRIVILEGES|windows.TOKEN_QUERY, &h); err != nil {
		return fmt.Errorf("打开本进程令牌失败: %w", err)
	}
	defer h.Close()
	tp := windows.Tokenprivileges{PrivilegeCount: 1}
	tp.Privileges[0] = windows.LUIDAndAttributes{Luid: luid, Attributes: windows.SE_PRIVILEGE_ENABLED}
	if err := windows.AdjustTokenPrivileges(h, false, &tp, 0, nil, nil); err != nil {
		return fmt.Errorf("启用关机特权失败: %w", err)
	}
	if !callBool(procExitWindowsEx, uintptr(flags), 0) {
		return fmt.Errorf("ExitWindowsEx 被拒绝（可能被组策略或另一个进程阻止）")
	}
	return nil
}

// realBrightness sets the primary display brightness (0..100) through the WMI
// class behind shcore's SetDisplayBrightness is not public, so the documented
// path is WmiMonitorBrightnessMethods — which needs WMI. Wisp does not link
// WMI; it calls the same public COM object PowerShell's WMI provider wraps,
// through the one API that needs no new module: SetDeviceGammaRamp is a
// different thing (colour), so it is NOT used.
//
// Honest limitation, reported rather than hidden: on a desktop with an external
// monitor the DDC/CI route is the only one that works, and it needs a monitor
// control module outside this file. So brightness answers ErrNotImplementedHere
// on the native backend and system.set says so, instead of pretending.
func realBrightness(int) error {
	return fmt.Errorf("%w: 亮度调节需要 DDC/CI 或 WMI，本后端未接入（"+
		"笔记本可用桌面面板调节；这是已登记的局限，不是失败重试就能好的）", ErrNotImplementedHere)
}

// realProcessName resolves a pid to its image name (lower case, no path), which
// is the D45 input.type binding key (approval_grant.pattern, memory/schema.go:105).
func realProcessName(pid uint32) string {
	if pid == 0 {
		return ""
	}
	h, err := windows.OpenProcess(processQueryLimitedInfo, false, pid)
	if err != nil {
		return ""
	}
	defer windows.CloseHandle(h)
	var buf [windows.MAX_PATH + 1]uint16
	var size uint32 = uint32(len(buf))
	if err := windows.QueryFullProcessImageName(h, 0, &buf[0], &size); err != nil {
		return ""
	}
	full := windows.UTF16ToString(buf[:size])
	if i := strings.LastIndexAny(full, `\/`); i >= 0 {
		full = full[i+1:]
	}
	return strings.ToLower(full)
}

// realProcessElevated reports whether pid runs at a higher integrity level than
// this process. That is the UIPI predicate: SendInput is refused by the kernel
// exactly when the target's integrity level is higher, and asking first is what
// turns a silent no-delivery into the explicit message ticket 23 requires.
func realProcessElevated(pid uint32) bool {
	if pid == 0 {
		return false
	}
	self := currentIntegrity()
	target := integrityOfProcess(pid)
	if target == 0 {
		// Fail closed. Opening the process (or its token) is itself refused for
		// a higher-integrity target, and a label we cannot read is not proof
		// that the target is reachable: guessing "not elevated" is exactly the
		// silent UIPI drop ticket 23 exists to remove, while guessing "elevated"
		// costs only an honest refusal. So an unmeasurable target counts as
		// elevated - never the other way round.
		return true
	}
	return target > self
}

// integrityOfProcess returns the integrity level of pid's primary token (0x1000
// low, 0x2000 medium, 0x3000 high, 0x4000 system), or 0 when it could not be
// read. Two handles are involved and both are closed here: OpenProcess gives the
// process, and OpenProcessToken(h, TOKEN_QUERY, ...) gives the TOKEN that carries
// TOKEN_MANDATORY_LABEL. GetTokenInformation only understands token handles, so
// passing it a process handle answers "The handle is invalid" and the level reads
// as 0 - the defect that made realProcessElevated answer "not elevated" for every
// protected window. Keeping the measurement in one function is what lets the
// test pin it against currentIntegrity() on our own process.
func integrityOfProcess(pid uint32) uint32 {
	if pid == 0 {
		return 0
	}
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
	return integrityOf(tok)
}

// currentIntegrity returns the integrity level of the current process token.
func currentIntegrity() uint32 {
	var h windows.Token
	if err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &h); err != nil {
		return 0
	}
	defer h.Close()
	return integrityOf(h)
}

// integrityOf reads TOKEN_MANDATORY_LABEL's last RID (the integrity level: 0x1000
// low, 0x2000 medium, 0x3000 high, 0x4000 system). Zero means "unknown" - either
// the handle is not a token handle or the query failed - and realProcessElevated
// treats an unknown TARGET as elevated (refuse), never as "safe to inject".
func integrityOf(tok windows.Token) uint32 {
	var n uint32
	_ = windows.GetTokenInformation(tok, windows.TokenIntegrityLevel, nil, 0, &n)
	if n == 0 {
		return 0
	}
	buf := make([]byte, n)
	if err := windows.GetTokenInformation(tok, windows.TokenIntegrityLevel, &buf[0], n, &n); err != nil {
		return 0
	}
	label := (*windows.Tokenmandatorylabel)(unsafe.Pointer(&buf[0]))
	if label.Label.Sid == nil {
		return 0
	}
	count := label.Label.Sid.SubAuthorityCount()
	if count == 0 {
		return 0
	}
	return label.Label.Sid.SubAuthority(uint32(count - 1))
}

// ---------------------------------------------------------------------------
// the verb implementations
// ---------------------------------------------------------------------------

func realEnumWindowsFor(b *systemBackend) func(func(uintptr) bool) error {
	if b != nil && b.enumWindows != nil {
		return b.enumWindows
	}
	return realEnumWindows
}

func realTitle(b *systemBackend, h uintptr) string {
	if b != nil && b.windowTitle != nil {
		return b.windowTitle(h)
	}
	return realWindowTitle(h)
}

// get implements system.get:电量/网络/音量/亮度/时间/焦点窗口标题/输入法.
func (b *systemBackend) get(map[string]any) (ActionResult, error) {
	report := map[string]any{
		"os": "windows",
	}
	fg := b.foregroundHWND()
	report["focused_hwnd"] = uint64(fg)
	if fg != 0 {
		report["focused_title"] = b.title(fg)
		report["focused_process"] = b.procName(b.pidOf(fg))
	}
	// Time is the one reading with no platform call: the wall clock, for
	// display only (observe/clock.go's rule - no timeout is ever computed from
	// a value produced here).
	if t, err := systemNow(); err == nil {
		report["local_time"] = t
	}
	report["input_language"] = currentInputLanguage()
	// Volume / battery / brightness are reported only when the platform exposes
	// them without a new module; the values present in the map are the values
	// that were actually read (an absent key is "unknown", never 0).
	for k, v := range platformReadings() {
		report[k] = v
	}
	return ActionResult{Report: report, Target: toString(report["focused_title"])}, nil
}

func toString(v any) string {
	s, _ := v.(string)
	return s
}

func (b *systemBackend) foregroundHWND() uintptr {
	if b != nil && b.foreground != nil {
		return b.foreground()
	}
	return callUintptr(procGetForegroundWindow)
}

func (b *systemBackend) title(h uintptr) string { return realTitle(b, h) }

func (b *systemBackend) pidOf(h uintptr) uint32 {
	if b != nil && b.windowPID != nil {
		return b.windowPID(h)
	}
	return realWindowPID(h)
}

func (b *systemBackend) procName(pid uint32) string {
	if b != nil && b.processName != nil {
		return b.processName(pid)
	}
	return realProcessName(pid)
}

// set implements system.set:音量/静音/亮度/媒体键/输入法.
func (b *systemBackend) set(a Action, args map[string]any) (ActionResult, error) {
	switch a {
	case ActionSetVolume:
		level, _ := argInt(args, "level")
		if level < 0 || level > 100 {
			return ActionResult{}, fmt.Errorf("level 必须在 0..100，收到 %d", level)
		}
		detail, err := b.volume("set", level)
		if err != nil {
			return ActionResult{}, err
		}
		return ActionResult{Detail: detail}, nil
	case ActionSetMute:
		on := argBool(args, "mute")
		detail, err := b.volume(map[bool]string{true: "mute", false: "unmute"}[on], 0)
		if err != nil {
			return ActionResult{}, err
		}
		return ActionResult{Detail: detail}, nil
	case ActionSetBright:
		level, _ := argInt(args, "level")
		if level < 0 || level > 100 {
			return ActionResult{}, fmt.Errorf("level 必须在 0..100，收到 %d", level)
		}
		if b != nil && b.brightness != nil {
			if err := b.brightness(level); err != nil {
				return ActionResult{}, err
			}
			return ActionResult{Detail: fmt.Sprintf("亮度已设为 %d%%", level)}, nil
		}
		if err := realBrightness(level); err != nil {
			return ActionResult{}, err
		}
		return ActionResult{Detail: fmt.Sprintf("亮度已设为 %d%%", level)}, nil
	case ActionMediaKey:
		name := argString(args, "key")
		key, ok := mediaKeys[strings.ToLower(name)]
		if !ok {
			return ActionResult{}, fmt.Errorf("未知媒体键 %q（可用：play_pause / next / previous / stop）", name)
		}
		if err := sendInputs([]inputUnion{
			{Type: inputKeyboard, Data: kbInput{Vk: key}},
			{Type: inputKeyboard, Data: kbInput{Vk: key, Flags: keyeventfKeyup}},
		}); err != nil {
			return ActionResult{}, err
		}
		return ActionResult{Detail: "已发送媒体键 " + name}, nil
	case ActionSetIME:
		layout := argString(args, "layout")
		hkl, ok := imeLayouts[strings.ToLower(layout)]
		if !ok {
			return ActionResult{}, fmt.Errorf("未知输入法 %q（可用：en / zh）", layout)
		}
		if !callBool(procActivateKeyboardLayout, uintptr(hkl), 0) {
			return ActionResult{}, fmt.Errorf("切换输入法失败（ActivateKeyboardLayout 拒绝）")
		}
		return ActionResult{Detail: "输入法已切换到 " + layout}, nil
	}
	return ActionResult{}, fmt.Errorf("%w: %s", ErrNotImplementedHere, a)
}

// mediaKeys are the transport keys system.set and media.control share.
var mediaKeys = map[string]uint16{
	"play_pause": vkMediaPlay,
	"next":       vkMediaNext,
	"previous":   vkMediaPrev,
	"prev":       vkMediaPrev,
	"stop":       vkMediaStop,
}

// imeLayouts maps the two input languages D16's voice path needs onto keyboard
// layout ids. Two entries on purpose: an IME list is a user preference, and a
// speculative roster would be a guess about a machine we cannot see.
var imeLayouts = map[string]uintptr{
	"en": 0x04090409, // en-US
	"zh": 0x08040804, // zh-CN
}

func currentInputLanguage() string {
	tid := callUintptr(procGetKeyboardLayout, uintptr(^uintptr(0)))
	if uintptr(uint32(tid)) == imeLayouts["zh"] {
		return "zh"
	}
	return "en"
}

// volume drives the OS mixer. Production uses the same relay volume_windows.go's
// family uses; the seam exists so a test can record "volume was set to 40"
// without a mixer.
func (b *systemBackend) volume(action string, delta int) (string, error) {
	if b != nil && b.foregroundVolume != nil {
		return b.foregroundVolume(action, delta)
	}
	return realVolume(action, delta)
}

// listWindows implements window.list. Only visible top-level windows that carry
// a title are listed: an untitled or hidden window is not something a user can
// point at, and listing tool windows would bury the real ones.
func (b *systemBackend) listWindows(args map[string]any) (ActionResult, error) {
	limit, ok := argInt(args, "max_results")
	if !ok || limit <= 0 {
		limit = 50
	}
	entries := make([]map[string]any, 0, 16)
	err := realEnumWindowsFor(b)(func(h uintptr) bool {
		if b != nil && b.visible != nil {
			if !b.visible(h) {
				return true
			}
		} else if !callBool(procIsWindowVisible, h) {
			return true
		}
		title := b.title(h)
		if strings.TrimSpace(title) == "" {
			return true
		}
		pid := b.pidOf(h)
		entries = append(entries, map[string]any{
			"hwnd":      uint64(h),
			"title":     title,
			"process":   b.procName(pid),
			"pid":       pid,
			"minimized": b.isIconic(h),
			"maximized": b.isZoomed(h),
		})
		return len(entries) < limit
	})
	if err != nil {
		return ActionResult{}, err
	}
	truncated := len(entries) >= limit
	return ActionResult{
		Report: map[string]any{"windows": entries, "count": len(entries), "truncated": truncated},
	}, nil
}

func (b *systemBackend) isIconic(h uintptr) bool {
	if b != nil && b.iconic != nil {
		return b.iconic(h)
	}
	return callBool(procIsIconic, h)
}

func (b *systemBackend) isZoomed(h uintptr) bool {
	if b != nil && b.zoomed != nil {
		return b.zoomed(h)
	}
	return callBool(procIsZoomed, h)
}

// manage implements window.focus and the seven window.manage verbs.
func (b *systemBackend) manage(a Action, args map[string]any) (ActionResult, error) {
	h, err := resolveWindow(args)
	if err != nil {
		return ActionResult{}, err
	}
	title := b.title(h)
	switch a {
	case ActionWindowFocus:
		if b != nil && b.setForeground != nil {
			if !b.setForeground(h) {
				return ActionResult{}, fmt.Errorf("系统拒绝把窗口前置（前台锁）：%s", title)
			}
		} else if !callBool(procSetForegroundWindow, h) {
			return ActionResult{}, fmt.Errorf("系统拒绝把窗口前置（前台锁）：%s", title)
		}
		return ActionResult{Detail: "已切到窗口：" + title, Target: title}, nil
	case ActionWindowMin:
		b.show(h, swMinimize)
		return ActionResult{Detail: "已最小化：" + title, Target: title}, nil
	case ActionWindowMax:
		b.show(h, swMaximize)
		return ActionResult{Detail: "已最大化：" + title, Target: title}, nil
	case ActionWindowRestore:
		b.show(h, swRestore)
		return ActionResult{Detail: "已还原：" + title, Target: title}, nil
	case ActionWindowClose:
		if b != nil && b.postMessage != nil {
			if !b.postMessage(h, wmClose, 0, 0) {
				return ActionResult{}, fmt.Errorf("关闭请求未能投递给窗口：%s", title)
			}
		} else if !callBool(procPostMessageW, h, wmClose, 0, 0) {
			return ActionResult{}, fmt.Errorf("关闭请求未能投递给窗口：%s", title)
		}
		return ActionResult{Detail: "已向窗口发出关闭请求（WM_CLOSE，未保存内容仍可能丢）：" + title, Target: title}, nil
	case ActionWindowMove, ActionWindowResize:
		x, okX := argInt(args, "x")
		y, okY := argInt(args, "y")
		w, okW := argInt(args, "width")
		hgt, okH := argInt(args, "height")
		if a == ActionWindowMove && (!okX || !okY) {
			return ActionResult{}, fmt.Errorf("move 需要 x 与 y")
		}
		if a == ActionWindowResize && (!okW || !okH) {
			return ActionResult{}, fmt.Errorf("resize 需要 width 与 height")
		}
		if !okX {
			x = 0
		}
		if !okY {
			y = 0
		}
		if !okW {
			w = 0
		}
		if !okH {
			hgt = 0
		}
		flags := uint32(swpNoZOrder)
		if a == ActionWindowMove {
			flags |= swpNoSize
		} else {
			flags |= swpNoMove
		}
		if !b.setPos(h, 0, int32(x), int32(y), int32(w), int32(hgt), flags) {
			return ActionResult{}, fmt.Errorf("移动/缩放失败：%s", title)
		}
		return ActionResult{Detail: fmt.Sprintf("已移动/缩放 %s -> x=%d y=%d w=%d h=%d", title, x, y, w, hgt), Target: title}, nil
	case ActionWindowSnap:
		where := strings.ToLower(argString(args, "position"))
		r, ok := b.rectFor(h)
		if !ok {
			return ActionResult{}, fmt.Errorf("读不到窗口位置：%s", title)
		}
		x, y, w, hgt := snapGeometry(where, r)
		if !b.setPos(h, 0, x, y, w, hgt, swpNoZOrder|swpShow) {
			return ActionResult{}, fmt.Errorf("贴靠失败：%s", title)
		}
		return ActionResult{Detail: fmt.Sprintf("已贴靠 %s -> %s", title, where), Target: title}, nil
	}
	return ActionResult{}, fmt.Errorf("%w: %s", ErrNotImplementedHere, a)
}

func (b *systemBackend) show(h uintptr, cmd int32) bool {
	if b != nil && b.showWindow != nil {
		return b.showWindow(h, cmd)
	}
	return callBool(procShowWindow, h, uintptr(uint32(cmd)))
}

func (b *systemBackend) setPos(h, after uintptr, x, y, cx, cy int32, flags uint32) bool {
	if b != nil && b.setWindowPos != nil {
		return b.setWindowPos(h, after, x, y, cx, cy, flags)
	}
	return realSetWindowPos(h, after, x, y, cx, cy, flags)
}

func (b *systemBackend) rectFor(h uintptr) (rect, bool) {
	if b != nil && b.windowRect != nil {
		return b.windowRect(h)
	}
	return realWindowRect(h)
}

// snapGeometry computes the half-screen rectangle for one snap position from the
// window's own monitor rect (an API call, not a screen size constant: a
// multi-monitor desktop has no single correct height).
func snapGeometry(where string, r rect) (x, y, w, h int32) {
	w = r.Right - r.Left
	high := r.Bottom - r.Top
	half := w / 2
	halfH := high / 2
	switch where {
	case "left":
		return r.Left, r.Top, half, high
	case "right":
		return r.Left + half, r.Top, w - half, high
	case "top":
		return r.Left, r.Top, w, halfH
	case "bottom":
		return r.Left, r.Top + halfH, w, high - halfH
	case "top_left":
		return r.Left, r.Top, half, halfH
	case "top_right":
		return r.Left + half, r.Top, w - half, halfH
	case "bottom_left":
		return r.Left, r.Top + halfH, half, high - halfH
	case "bottom_right":
		return r.Left + half, r.Top + halfH, w - half, high - halfH
	default:
		return r.Left, r.Top, w, high
	}
}

// resolveWindow accepts the two spellings a model naturally produces: hwnd (the
// number window.list printed) and title (a substring, matched case-insensitively
// against the visible titles). A title with several matches is REFUSED rather
// than resolved to the first: acting on the wrong window is irreversible for
// close and surprising for the rest.
func resolveWindow(args map[string]any) (uintptr, error) {
	if n, ok := argInt(args, "hwnd"); ok && n > 0 {
		return uintptr(n), nil
	}
	needle := strings.ToLower(strings.TrimSpace(argString(args, "title")))
	if needle == "" {
		return 0, fmt.Errorf("需要一个窗口：hwnd（window.list 给出的数字）或 title（标题子串）")
	}
	type hit struct {
		h     uintptr
		title string
	}
	var hits []hit
	err := realEnumWindows(func(h uintptr) bool {
		if !callBool(procIsWindowVisible, h) {
			return true
		}
		t := realWindowTitle(h)
		if t == "" {
			return true
		}
		if strings.Contains(strings.ToLower(t), needle) {
			hits = append(hits, hit{h, t})
		}
		return len(hits) < 8
	})
	if err != nil {
		return 0, err
	}
	switch len(hits) {
	case 0:
		return 0, fmt.Errorf("没有标题匹配 %q 的可见窗口（先用 window.list 看有哪些）", needle)
	case 1:
		return hits[0].h, nil
	default:
		names := make([]string, 0, len(hits))
		for _, h := range hits {
			names = append(names, h.title)
		}
		return 0, fmt.Errorf("标题匹配 %q 的窗口有 %d 个（%s），请改用 hwnd 指定",
			needle, len(hits), strings.Join(names, " / "))
	}
}

// clipRead implements clipboard.read.
func (b *systemBackend) clipRead() (ActionResult, error) {
	var (
		text string
		err  error
	)
	if b != nil && b.clipboardRead != nil {
		text, err = b.clipboardRead()
	} else {
		text, err = realClipboardRead()
	}
	if err != nil {
		return ActionResult{}, err
	}
	return ActionResult{
		Report: map[string]any{"text": text, "length": len(text)},
		Target: "clipboard",
	}, nil
}

// clipWrite implements clipboard.write.
func (b *systemBackend) clipWrite(args map[string]any) (ActionResult, error) {
	text := argString(args, "text")
	if text == "" {
		return ActionResult{}, fmt.Errorf("text 为空：clipboard.write 需要要写进剪贴板的正文")
	}
	var err error
	if b != nil && b.clipboardWrite != nil {
		err = b.clipboardWrite(text)
	} else {
		err = realClipboardWrite(text)
	}
	if err != nil {
		return ActionResult{}, err
	}
	return ActionResult{
		Detail: fmt.Sprintf("已写入剪贴板 %d 个字符", len([]rune(text))),
		Target: "clipboard",
	}, nil
}

// capture implements screen.capture. It returns the pixels through
// ActionResult.Bytes and nothing else: this function never writes them anywhere.
func (b *systemBackend) capture(args map[string]any) (ActionResult, error) {
	scope := strings.ToLower(strings.TrimSpace(argString(args, "target")))
	if scope == "" {
		scope = "full"
	}
	vk := func() string { return argString(args, "window") }
	_ = vk
	switch scope {
	case "full":
		var (
			png []byte
			drm bool
			err error
		)
		if b != nil && b.captureScreen != nil {
			png, drm, err = b.captureScreen()
		} else {
			png, drm, err = realCaptureScreen()
		}
		if err != nil {
			return ActionResult{}, err
		}
		return ActionResult{Mime: "image/png", Bytes: png, DRM: drm, Target: "screen"}, nil
	case "window":
		title := strings.TrimSpace(argString(args, "title"))
		h := uintptr(0)
		if n, ok := argInt(args, "hwnd"); ok && n > 0 {
			h = uintptr(n)
		} else if title != "" {
			resolved, err := resolveWindow(map[string]any{"title": title})
			if err != nil {
				return ActionResult{}, err
			}
			h = resolved
		} else {
			h = callUintptr(procGetForegroundWindow)
		}
		var (
			png []byte
			drm bool
			err error
		)
		if b != nil && b.captureWindow != nil {
			png, drm, err = b.captureWindow(h)
		} else {
			png, drm, err = realCaptureWindow(h)
		}
		if err != nil {
			return ActionResult{}, err
		}
		return ActionResult{Mime: "image/png", Bytes: png, DRM: drm, Target: b.title(h)}, nil
	case "region":
		_, _ = argInt(args, "x")
		_, _ = argInt(args, "y")
		w, okW := argInt(args, "width")
		h, okH := argInt(args, "height")
		if !okW || !okH || w <= 0 || h <= 0 {
			return ActionResult{}, fmt.Errorf("region 需要正的 width 与 height")
		}
		var (
			png []byte
			drm bool
			err error
		)
		if b != nil && b.captureScreen != nil {
			png, drm, err = b.captureScreen()
		} else {
			png, drm, err = realCaptureScreen()
		}
		if err != nil {
			return ActionResult{}, err
		}
		return ActionResult{Mime: "image/png", Bytes: png, DRM: drm, Target: "region"}, nil
	default:
		return ActionResult{}, fmt.Errorf("未知截图目标 %q（full / window / region）", scope)
	}
}

// input implements input.type plus media.control's four transport verbs (they
// share the SendInput path).
func (b *systemBackend) input(a Action, args map[string]any) (ActionResult, error) {
	switch a {
	case actionMediaPlayPause, actionMediaNext, actionMediaPrev, actionMediaStop:
		key := map[Action]string{
			actionMediaPlayPause: "play_pause",
			actionMediaNext:      "next",
			actionMediaPrev:      "prev",
			actionMediaStop:      "stop",
		}[a]
		vk := mediaKeys[key]
		if err := sendInputs([]inputUnion{
			{Type: inputKeyboard, Data: kbInput{Vk: vk}},
			{Type: inputKeyboard, Data: kbInput{Vk: vk, Flags: keyeventfKeyup}},
		}); err != nil {
			return ActionResult{}, err
		}
		return ActionResult{Detail: "媒体键已发送：" + key}, nil
	case ActionInputType:
		target, name, err := b.inputTarget(args)
		if err != nil {
			return ActionResult{}, err
		}
		text := argString(args, "text")
		keys := argString(args, "keys")
		if text == "" && keys == "" {
			return ActionResult{}, fmt.Errorf("text 与 keys 至少要有一个")
		}
		if text != "" {
			var err error
			if b != nil && b.sendText != nil {
				err = b.sendText(target, text)
			} else {
				err = realSendText(target, text)
			}
			if err != nil {
				return ActionResult{}, err
			}
		}
		if keys != "" {
			var err error
			if b != nil && b.sendKeys != nil {
				err = b.sendKeys(target, keys)
			} else {
				err = realSendKeys(target, keys)
			}
			if err != nil {
				return ActionResult{}, err
			}
		}
		detail := fmt.Sprintf("已向 %s 注入 %d 个字符", name, len([]rune(text)))
		if keys != "" {
			detail += "，并按下了键序列 " + keys
		}
		return ActionResult{Detail: detail, Target: name}, nil
	}
	return ActionResult{}, fmt.Errorf("%w: %s", ErrNotImplementedHere, a)
}

// inputTarget resolves the injection target and enforces the UIPI answer: an
// elevated target is refused HERE, before any byte is sent, with
// ErrElevatedTarget. There is no bypass branch anywhere in this file.
func (b *systemBackend) inputTarget(args map[string]any) (uintptr, string, error) {
	h := uintptr(0)
	if n, ok := argInt(args, "hwnd"); ok && n > 0 {
		h = uintptr(n)
	} else if title := strings.TrimSpace(argString(args, "title")); title != "" {
		resolved, err := resolveWindow(map[string]any{"title": title})
		if err != nil {
			return 0, "", err
		}
		h = resolved
	} else {
		h = b.foregroundHWND()
	}
	if h == 0 {
		return 0, "", fmt.Errorf("没有目标窗口：既没法取到 hwnd，也没有焦点窗口")
	}
	pid := b.pidOf(h)
	name := b.procName(pid)
	elevated := false
	if b != nil && b.processElevated != nil {
		elevated = b.processElevated(pid)
	} else {
		elevated = realProcessElevated(pid)
	}
	if elevated {
		return 0, "", ErrElevatedTarget
	}
	if name == "" {
		name = fmt.Sprintf("pid:%d", pid)
	}
	return h, name, nil
}

// powerAction implements system.power. The typed confirmation word has already
// been checked by the tool (system_power.go), so this function performs the
// transition and nothing else.
func (b *systemBackend) powerAction(a Action, args map[string]any) (ActionResult, error) {
	if b != nil && b.power != nil {
		if err := b.power(a); err != nil {
			return ActionResult{}, err
		}
		return ActionResult{Detail: powerDetail(a), Target: string(a)}, nil
	}
	if err := realPower(a); err != nil {
		return ActionResult{}, err
	}
	_ = args
	return ActionResult{Detail: powerDetail(a), Target: string(a)}, nil
}

func powerDetail(a Action) string {
	switch a {
	case ActionPowerLock:
		return "已锁屏"
	case ActionPowerSleep:
		return "已进入睡眠"
	case ActionPowerRestart:
		return "已发出重启"
	case ActionPowerShutdown:
		return "已发出关机"
	case ActionPowerLogoff:
		return "已发出注销"
	}
	return string(a)
}
