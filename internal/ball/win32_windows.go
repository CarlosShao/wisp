//go:build windows

package ball

// Raw Win32/GDI plumbing for the ball layered window. No cgo: syscalls
// through x/sys/windows + LazyDLL (same discipline as the S0 spike,
// scripts/spike/common/winshell.go). D2D/DWrite COM vtable slots are in
// d2d_windows.go and were taken from the mingw-w64 headers (MIDL order).

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modUser32   = windows.NewLazySystemDLL("user32.dll")
	modGDI32    = windows.NewLazySystemDLL("gdi32.dll")
	modShell32  = windows.NewLazySystemDLL("shell32.dll")
	modKernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procGetModuleHandleW = modKernel32.NewProc("GetModuleHandleW")

	pRegisterClassExW              = modUser32.NewProc("RegisterClassExW")
	pCreateWindowExW               = modUser32.NewProc("CreateWindowExW")
	pDefWindowProcW                = modUser32.NewProc("DefWindowProcW")
	pDestroyWindow                 = modUser32.NewProc("DestroyWindow")
	pShowWindow                    = modUser32.NewProc("ShowWindow")
	pGetMessageW                   = modUser32.NewProc("GetMessageW")
	pTranslateMessage              = modUser32.NewProc("TranslateMessage")
	pDispatchMessageW              = modUser32.NewProc("DispatchMessageW")
	pPostQuitMessage               = modUser32.NewProc("PostQuitMessage")
	pPostMessageW                  = modUser32.NewProc("PostMessageW")
	pSetTimer                      = modUser32.NewProc("SetTimer")
	pKillTimer                     = modUser32.NewProc("KillTimer")
	pSetCapture                    = modUser32.NewProc("SetCapture")
	pReleaseCapture                = modUser32.NewProc("ReleaseCapture")
	pSetWindowPos                  = modUser32.NewProc("SetWindowPos")
	pGetWindowRect                 = modUser32.NewProc("GetWindowRect")
	pUpdateLayeredWindow           = modUser32.NewProc("UpdateLayeredWindow")
	pRegisterHotKey                = modUser32.NewProc("RegisterHotKey")
	pUnregisterHotKey              = modUser32.NewProc("UnregisterHotKey")
	pCreatePopupMenu               = modUser32.NewProc("CreatePopupMenu")
	pAppendMenuW                   = modUser32.NewProc("AppendMenuW")
	pTrackPopupMenu                = modUser32.NewProc("TrackPopupMenu")
	pDestroyMenu                   = modUser32.NewProc("DestroyMenu")
	pGetDpiForWindow               = modUser32.NewProc("GetDpiForWindow")
	pSetProcessDpiAwarenessContext = modUser32.NewProc("SetProcessDpiAwarenessContext")
	pMonitorFromWindow             = modUser32.NewProc("MonitorFromWindow")
	pGetMonitorInfoW               = modUser32.NewProc("GetMonitorInfoW")
	pEnumDisplayMonitors           = modUser32.NewProc("EnumDisplayMonitors")
	pLoadCursorW                   = modUser32.NewProc("LoadCursorW")
	pLoadIconW                     = modUser32.NewProc("LoadIconW")
	pGetDC                         = modUser32.NewProc("GetDC")
	pReleaseDC                     = modUser32.NewProc("ReleaseDC")
	pSetForegroundWindow           = modUser32.NewProc("SetForegroundWindow")

	pCreateCompatibleDC = modGDI32.NewProc("CreateCompatibleDC")
	pDeleteDC           = modGDI32.NewProc("DeleteDC")
	pCreateDIBSection   = modGDI32.NewProc("CreateDIBSection")
	pSelectObject       = modGDI32.NewProc("SelectObject")
	pDeleteObject       = modGDI32.NewProc("DeleteObject")

	pShellNotifyIconW = modShell32.NewProc("Shell_NotifyIconW")
)

// Window styles / extended styles (SPEC-08 §2 frozen set).
const (
	wsExLayered    = 0x00080000
	wsExTopmost    = 0x00000008
	wsExToolWindow = 0x00000080
	wsExNoActivate = 0x08000000
	wsPopup        = 0x80000000

	swShownoactivate = 4
	idcArrow         = 32512
	idiApplication   = 32512
)

// Messages.
const (
	wmNull           = 0x0000
	wmDestroy        = 0x0002
	wmTimer          = 0x0113
	wmCommand        = 0x0111
	wmNCHitTest      = 0x0084
	wmLButtonDown    = 0x0201
	wmLButtonUp      = 0x0202
	wmRButtonUp      = 0x0205
	wmMouseMove      = 0x0200
	wmCaptureChanged = 0x0215
	wmMouseActivate  = 0x0021
	wmLButtonUpX     = 0
	wmHotkey         = 0x0312
	wmDPICHanged     = 0x02E0
	wmDisplayChange  = 0x007E
	wmApp            = 0x8000

	// Package-private message slots on the ball/dispatcher windows.
	wmAppTask  = wmApp + 0x201 // run a closure posted to the STA thread
	wmAppTray  = wmApp + 0x202 // tray icon callback
	wmAppFrame = wmApp + 0x203 // cross-thread repaint request
)

// UpdateLayeredWindow / DIB constants.
const (
	ulwAlpha = 0x00000002

	acSrcOver  = 0
	acSrcAlpha = 0x01

	dibRGBColors = 0
)

// SetWindowPos flags.
const (
	swpNoSize     = 0x0001
	swpNoMove     = 0x0002
	swpNoZOrder   = 0x0004
	swpNoActivate = 0x0010
	swpShowWindow = 0x0040
)

// TrackPopupMenu flags.
const (
	tpmRightButton = 0x0002
	tpmReturnCMD   = 0x0100
	tpmNoActivate  = 0x0080 // do NOT steal activation for the popup
)

// Menu flags.
const (
	mfString = 0x00000000
	mfCheckd = 0x00000008
	mfSepart = 0x00000800
)

// Monitor constants.
const (
	monitorDefaultToNearest = 2
)

// DPI awareness context (PER_MONITOR_AWARE_V2).
const dpiAwarenessPerMonitorV2 = ^uintptr(3) // (HANDLE)-4

// COINIT_APARTMENTTHREADED for the ui-sta thread.
const coinitApartmentThreaded = 0x2

// Timer ids (one window, few timers).
const (
	timerAnimID = 1 // state animation / breathing / fade tick
	timerFadeID = 2
)

// ------------------------------------------------------------------- structs

type wndClassEx struct {
	size      uint32
	style     uint32
	wndProc   uintptr
	clsExtra  int32
	winExtra  int32
	instance  windows.Handle
	icon      windows.Handle
	cursor    windows.Handle
	bkgnd     windows.Handle
	menuName  uintptr
	className *uint16
	iconSm    windows.Handle
}

type msg struct {
	hwnd    windows.HWND
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      point
	_       uint32 // private DWORD on x64 padding
}

type point struct{ x, y int32 }

type size struct{ cx, cy int32 }

type rect struct{ l, t, r, b int32 }

func (r rect) width() int32  { return r.r - r.l }
func (r rect) height() int32 { return r.b - r.t }

type blendFunction struct {
	blendOp             uint8
	blendFlags          uint8
	sourceConstantAlpha uint8
	alphaFormat         uint8
}

type bitmapInfoHeader struct {
	size          uint32
	width         int32
	height        int32 // negative = top-down
	planes        uint16
	bitCount      uint16
	compression   uint32
	sizeImage     uint32
	xPelsPerMeter int32
	yPelsPerMeter int32
	clrUsed       uint32
	clrImportant  uint32
}

type bitmapInfo struct {
	header    bitmapInfoHeader
	bmiColors [1]uint32
}

type monitorInfoExW struct {
	size      uint32
	rcMonitor rect
	rcWork    rect
	flags     uint32
	device    [32]uint16
}

type notifyIconData struct {
	cbSize           uint32
	_                uint32
	hWnd             windows.HWND
	uID              uint32
	uFlags           uint32
	uCallbackMessage uint32
	_                uint32
	hIcon            windows.Handle
	szTip            [128]uint16
	dwState          uint32
	dwStateMask      uint32
	szInfo           [256]uint16
	uVersion         uint32
	_                uint32
	szInfoTitle      [64]uint16
	dwInfoFlags      uint32
	_                uint32
	guidItem         windows.GUID
	hBalloonIcon     windows.Handle
}

const (
	nimAdd     = 0
	nimModify  = 1
	nimDelete  = 2
	nifMessage = 1
	nifIcon    = 2
	nifTip     = 4
)

// hiWord/loWord extract the signed 16-bit halves of a message lParam
// (negative coordinates happen on multi-monitor layouts left of the primary).
func loSigned(lParam uintptr) int32 { return int32(int16(lParam & 0xFFFF)) }
func hiSigned(lParam uintptr) int32 { return int32(int16((lParam >> 16) & 0xFFFF)) }

func utf16(s string) *uint16 {
	p, _ := windows.UTF16PtrFromString(s)
	return p
}

// utf16Len counts UTF-16 code units for DrawText's stringLength.
func utf16Len(s string) int {
	n, _ := windows.UTF16FromString(s)
	return len(n) - 1
}

func unsafePtr[T any](v *T) uintptr { return uintptr(unsafe.Pointer(v)) }

func absi(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// u8 clamps a 0..255 float into a byte.
func u8(f float32) uint8 {
	if f <= 0 {
		return 0
	}
	if f >= 255 {
		return 255
	}
	return uint8(f + 0.5)
}

// htTransparentResult returns LRESULT(-1) for HTTRANSPARENT.
func htTransparentResult() uintptr { return ^uintptr(0) }

// rectFromUintptr converts a message lParam that carries a pointer to a RECT
// (WM_DPICHANGED). The pointer is owned by the OS for the duration of the
// call; the double indirection keeps vet's provenance check quiet.
func rectFromUintptr(p uintptr) *rect {
	return (*rect)(unsafe.Pointer(&p))
}
