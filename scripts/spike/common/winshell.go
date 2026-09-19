package common

// Raw Win32/D2D/DWrite plumbing for the S0 spike (ticket 02). No cgo: all
// syscalls through x/sys/windows + LazyDLL. Vtable slot indices were taken
// from the mingw-w64 headers (d2d1.h / dwrite.h, MIDL order) and are listed
// in docs/evidence/s0/02-spike-report.md.
//
// Float-argument ABI note: syscall args are integer words. On win-x64 any
// float arg beyond position 4 travels on the stack as its 4-byte IEEE bits
// in the low half of an 8-byte slot, so passing math.Float32bits(f) as an
// uintptr is correct there. Register-passed floats (positions 1-4) would NOT
// work this way; none of the calls used here needs one.

import (
	"fmt"
	"math"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modUser32 = windows.NewLazySystemDLL("user32.dll")
	modShell  = windows.NewLazySystemDLL("shell32.dll")
	modD2D1   = windows.NewLazySystemDLL("d2d1.dll")
	modDWrite = windows.NewLazySystemDLL("dwrite.dll")
	modOle32  = windows.NewLazySystemDLL("ole32.dll")

	pRegisterClassExW           = modUser32.NewProc("RegisterClassExW")
	pCreateWindowExW            = modUser32.NewProc("CreateWindowExW")
	pDefWindowProcW             = modUser32.NewProc("DefWindowProcW")
	pShowWindow                 = modUser32.NewProc("ShowWindow")
	pSetLayeredWindowAttributes = modUser32.NewProc("SetLayeredWindowAttributes")
	pPeekMessageW               = modUser32.NewProc("PeekMessageW")
	pTranslateMessage           = modUser32.NewProc("TranslateMessage")
	pDispatchMessageW           = modUser32.NewProc("DispatchMessageW")
	pDestroyWindow              = modUser32.NewProc("DestroyWindow")
	pLoadIconW                  = modUser32.NewProc("LoadIconW")
	pLoadCursorW                = modUser32.NewProc("LoadCursorW")
	pRegisterHotKey             = modUser32.NewProc("RegisterHotKey")
	pUnregisterHotKey           = modUser32.NewProc("UnregisterHotKey")

	pShellNotifyIconW = modShell.NewProc("Shell_NotifyIconW")

	pD2D1CreateFactory   = modD2D1.NewProc("D2D1CreateFactory")
	pDWriteCreateFactory = modDWrite.NewProc("DWriteCreateFactory")
	pCoInitializeEx      = modOle32.NewProc("CoInitializeEx")
)

const (
	wsExLayered    = 0x00080000
	wsExTopmost    = 0x00000008
	wsExToolWindow = 0x00000080
	wsExNoActivate = 0x08000000

	swShownoactivate = 4
	lwaAlpha         = 0x02
	colorWindow      = 5 // COLOR_WINDOW
	idcArrow         = 32512
	idiApplication   = 32512

	wmApp = 0x8000

	nimAdd     = 0
	nimDelete  = 2
	nifMessage = 1
	nifIcon    = 2
	nifTip     = 4

	modAlt      = 0x0001
	modControl  = 0x0002
	modNoRepeat = 0x4000

	coinitApartmentMT = 0x0

	d2d1FactoryTypeSingleThreaded = 0
)

// ---------------------------------------------------------------- COM helpers

// comCall invokes slot `slot` (0-based, IUnknown::QueryInterface = 0) of the
// COM object at `this`.
func comCall(this uintptr, slot uintptr, args ...uintptr) (uintptr, uintptr) {
	vt := *(*uintptr)(unsafe.Pointer(this))
	fn := *(*uintptr)(unsafe.Pointer(vt + slot*8))
	r1, r2, _ := syscall.SyscallN(fn, append([]uintptr{this}, args...)...)
	return r1, r2
}

func guid(data4 [8]byte, d1 uint32, d2, d3 uint16) windows.GUID {
	return windows.GUID{Data1: d1, Data2: d2, Data3: d3, Data4: data4}
}

func f32bits(f float32) uintptr { return uintptr(math.Float32bits(f)) }

var (
	// IID_ID2D1Factory 06152247-6f50-465a-9245-118bfd3b6007
	iidID2D1Factory = guid([8]byte{0x92, 0x45, 0x11, 0x8b, 0xfd, 0x3b, 0x60, 0x07}, 0x06152247, 0x6f50, 0x465a)
	// IID_IDWriteFactory b859ee5a-d838-4b5b-a2e8-1adc7d93db48
	iidIDWriteFactory = guid([8]byte{0xa2, 0xe8, 0x1a, 0xdc, 0x7d, 0x93, 0xdb, 0x48}, 0xb859ee5a, 0xd838, 0x4b5b)
)

// ---------------------------------------------------------------- D2D window

type d2d1RenderTargetProps struct {
	typ       uint32
	format    uint32 // DXGI_FORMAT_UNKNOWN
	alphaMode uint32 // D2D1_ALPHA_MODE_UNKNOWN
	dpiX      float32
	dpiY      float32
	usage     uint32
	minLevel  uint32
}

type d2d1HwndRenderTargetProps struct {
	hwnd      windows.HWND
	pixelSize struct{ w, h uint32 }
	options   uint32
	pad       uint32
}

type d2d1BrushProps struct {
	opacity float32
	// D2D1_MATRIX_3X2_F identity
	m [6]float32
}

type d2d1ColorF struct{ r, g, b, a float32 }

type d2d1RectF struct{ l, t, r, b float32 }

// D2DWindow is a real visible layered window with a Direct2D Hwnd render
// target that has presented at least one frame containing DWrite text.
type D2DWindow struct {
	Hwnd      windows.HWND
	factory   uintptr
	rt        uintptr
	brush     uintptr
	writeFact uintptr
	textFmt   uintptr
}

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

// CreateD2DWindow builds and shows a 260x92 layered window, creates the D2D
// factory / HwndRenderTarget / solid brush and the DWrite factory / text
// format, then presents one frame (clear + fill + DrawText) and pumps the
// message queue so DWM composites it.
func CreateD2DWindow(title string) (*D2DWindow, error) {
	pCoInitializeEx.Call(0, coinitApartmentMT)

	clsName, err := windows.UTF16PtrFromString("wisp_spike_d2d")
	if err != nil {
		return nil, err
	}
	icon, _, _ := pLoadIconW.Call(0, idiApplication)
	cursor, _, _ := pLoadCursorW.Call(0, idcArrow)
	modh, _ := windows.GetModuleHandle(nil)
	wc := wndClassEx{
		size:      uint32(unsafe.Sizeof(wndClassEx{})),
		wndProc:   windows.NewCallback(defWndProc),
		instance:  modh,
		icon:      windows.Handle(icon),
		iconSm:    windows.Handle(icon),
		cursor:    windows.Handle(cursor),
		bkgnd:     windows.Handle(colorWindow + 1), // HBRUSH(COLOR_WINDOW+1)
		className: clsName,
	}
	r1, _, err := pRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
	if r1 == 0 {
		return nil, fmt.Errorf("RegisterClassExW: %v", err)
	}

	t16, _ := windows.UTF16PtrFromString(title)
	hwnd, _, err := pCreateWindowExW(
		wsExLayered|wsExTopmost|wsExToolWindow|wsExNoActivate,
		uintptr(unsafe.Pointer(clsName)),
		uintptr(unsafe.Pointer(t16)),
		0, // WS_OVERLAPPED+not visible; ShowWindow below
		40, 40, 260, 92,
		0, 0, uintptr(modh), 0)
	if hwnd == 0 {
		return nil, fmt.Errorf("CreateWindowExW: %v", err)
	}

	w := &D2DWindow{Hwnd: windows.HWND(hwnd)}

	// D2D1CreateFactory(factoryType, riid, *D2D1_FACTORY_OPTIONS, out) - 4 args
	var factory uintptr
	opts := uint32(0) // D2D1_DEBUG_LEVEL_NONE
	r1, _, err = pD2D1CreateFactory.Call(
		d2d1FactoryTypeSingleThreaded,
		uintptr(unsafe.Pointer(&iidID2D1Factory)),
		uintptr(unsafe.Pointer(&opts)),
		uintptr(unsafe.Pointer(&factory)))
	if r1 != 0 {
		return nil, fmt.Errorf("D2D1CreateFactory: %v", err)
	}
	w.factory = factory

	// CreateHwndRenderTarget = slot 14 of ID2D1Factory (after 3x IUnknown)
	rtProps := d2d1RenderTargetProps{}
	hwProps := d2d1HwndRenderTargetProps{
		hwnd:      windows.HWND(hwnd),
		pixelSize: struct{ w, h uint32 }{w: 260, h: 92},
	}
	var rt uintptr
	hr, _, _ := comCall(factory, 14,
		uintptr(unsafe.Pointer(&rtProps)),
		uintptr(unsafe.Pointer(&hwProps)),
		uintptr(unsafe.Pointer(&rt)))
	if hr != 0 {
		return nil, fmt.Errorf("CreateHwndRenderTarget hr=0x%x", hr)
	}
	w.rt = rt

	// CreateSolidColorBrush = slot 8 of ID2D1RenderTarget
	// (vtable: 3 GetFactory; 4 CreateBitmap; 5 CreateBitmapFromWicBitmap;
	//  6 CreateSharedBitmap; 7 CreateBitmapBrush; 8 CreateSolidColorBrush)
	brushColor := d2d1ColorF{r: 0.16, g: 0.16, b: 0.20, a: 1.0}
	brushProps := d2d1BrushProps{opacity: 1.0, m: [6]float32{1, 0, 0, 1, 0, 0}}
	var brush uintptr
	hr, _, _ = comCall(rt, 8,
		uintptr(unsafe.Pointer(&brushColor)),
		uintptr(unsafe.Pointer(&brushProps)),
		uintptr(unsafe.Pointer(&brush)))
	if hr != 0 {
		return nil, fmt.Errorf("CreateSolidColorBrush hr=0x%x", hr)
	}
	w.brush = brush

	// DWriteCreateFactory(DWRITE_FACTORY_TYPE_SHARED=0, riid, out)
	var dwf uintptr
	r1, _, err = pDWriteCreateFactory.Call(
		0,
		uintptr(unsafe.Pointer(&iidIDWriteFactory)),
		uintptr(unsafe.Pointer(&dwf)))
	if r1 != 0 {
		return nil, fmt.Errorf("DWriteCreateFactory: %v", err)
	}
	w.writeFact = dwf

	// CreateTextFormat = slot 15 of ID2D1Factory's counterpart IDWriteFactory:
	// (family, collection=0, weight=400, style=0, stretch=5, fontSize(float,
	// stack-passed as bits), locale, out)
	family, _ := windows.UTF16PtrFromString("Microsoft YaHei")
	locale, _ := windows.UTF16PtrFromString("zh-CN")
	var tf uintptr
	hr, _, _ = comCall(dwf, 15,
		uintptr(unsafe.Pointer(family)),
		0, // system font collection
		400,
		0,
		5,
		f32bits(18.0),
		uintptr(unsafe.Pointer(locale)),
		uintptr(unsafe.Pointer(&tf)))
	if hr != 0 {
		return nil, fmt.Errorf("CreateTextFormat hr=0x%x", hr)
	}
	w.textFmt = tf

	// Show + alpha so it is a real composited layered window
	pShowWindow.Call(hwnd, swShownoactivate)
	alpha := uint32(235)
	pSetLayeredWindowAttributes.Call(hwnd, 0, uintptr(alpha), lwaAlpha)

	w.RenderFrame("Wisp S0 spike")
	pumpMessages(6)
	return w, nil
}

// RenderFrame draws one full frame: background clear, accent fill, text.
func (w *D2DWindow) RenderFrame(text string) {
	if w.rt == 0 {
		return
	}
	bg := d2d1ColorF{r: 0.06, g: 0.06, b: 0.09, a: 1.0}
	rect := d2d1RectF{l: 12, t: 30, r: 248, b: 62}
	textRect := d2d1RectF{l: 12, t: 30, r: 248, b: 86}
	s16, _ := windows.UTF16PtrFromString(text)
	comCall(w.rt, 48, uintptr(0)) // BeginDraw (vtable slot 48)
	comCall(w.rt, 47, uintptr(unsafe.Pointer(&bg))) // Clear (47)
	accent := d2d1ColorF{r: 0.22, g: 0.51, b: 0.93, a: 1.0}
	_ = accent
	comCall(w.rt, 17, // FillRectangle (17)
		uintptr(unsafe.Pointer(&rect)), uintptr(unsafe.Pointer(w.brush)))
	comCall(w.rt, 27, // DrawText (27)
		uintptr(unsafe.Pointer(s16)),
		uintptr(uint32(len([]rune(text)))),
		w.textFmt,
		uintptr(unsafe.Pointer(&textRect)),
		w.brush,
		0, 0)
	comCall(w.rt, 49, 0, 0) // EndDraw (49)
}

// TearDown releases D2D/DWrite COM objects and destroys the window.
func (w *D2DWindow) TearDown() {
	release := func(p *uintptr) {
		if *p != 0 {
			comCall(*p, 2) // IUnknown::Release
			*p = 0
		}
	}
	release(&w.textFmt)
	release(&w.writeFact)
	release(&w.brush)
	release(&w.rt)
	release(&w.factory)
	if w.Hwnd != 0 {
		pDestroyWindow.Call(uintptr(w.Hwnd))
		w.Hwnd = 0
	}
}

func defWndProc(hwnd uintptr, msg uintptr, wParam, lParam uintptr) uintptr {
	r, _, _ := pDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

func pumpMessages(rounds int) {
	type msg struct {
		hwnd    windows.HWND
		message uint32
		wParam  uintptr
		lParam  uintptr
		time    uint32
		pt      struct{ x, y int32 }
	}
	var m msg
	for i := 0; i < rounds; i++ {
		for {
			r, _, _ := pPeekMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0, 1) // PM_REMOVE
			if r == 0 {
				break
			}
			pTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
			pDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// ---------------------------------------------------------------- tray icon

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
	uTimeout         uint32
	_                uint32
	szInfoTitle      [64]uint16
	dwInfoFlags      uint32
	_                uint32
	guidItem         windows.GUID
	hBalloonIcon     windows.Handle
}

// TrayIcon wraps a Shell_NotifyIconW tray entry owned by a window.
type TrayIcon struct {
	hwnd windows.HWND
	data notifyIconData
}

func AddTrayIcon(owner windows.HWND) (*TrayIcon, error) {
	icon, _, _ := pLoadIconW.Call(0, idiApplication)
	tip, _ := windows.UTF16FromString("wisp spike")
	d := notifyIconData{
		cbSize:           uint32(unsafe.Sizeof(notifyIconData{})),
		hWnd:             owner,
		uID:              0x5701,
		uFlags:           nifMessage | nifIcon | nifTip,
		uCallbackMessage: wmApp + 0x101,
		hIcon:            windows.Handle(icon),
	}
	copy(d.szTip[:], tip)
	t := &TrayIcon{hwnd: owner, data: d}
	r1, _, err := pShellNotifyIconW.Call(nimAdd, uintptr(unsafe.Pointer(&t.data)))
	if r1 == 0 {
		return nil, fmt.Errorf("Shell_NotifyIconW(NIM_ADD): %v", err)
	}
	return t, nil
}

func (t *TrayIcon) Remove() {
	if t == nil {
		return
	}
	pShellNotifyIconW.Call(nimDelete, uintptr(unsafe.Pointer(&t.data)))
}

// ---------------------------------------------------------------- hotkey

// RegisterGlobalHotKey tries a small list of low-conflict combos and returns
// the first that registers (owners may have bound others).
func RegisterGlobalHotKey(owner windows.HWND) (registered bool, mods, vk uint32) {
	type combo struct {
		mods, vk uint32
	}
	for _, c := range []combo{
		{modControl | modAlt | modNoRepeat, 0x59}, // Ctrl+Alt+Y
		{modControl | modAlt | modNoRepeat, 0x5A}, // Ctrl+Alt+Z
		{modControl | modAlt | modNoRepeat, 0xDB}, // Ctrl+Alt+[
	} {
		r1, _, _ := pRegisterHotKey.Call(uintptr(owner), 0x5702, uintptr(c.mods), uintptr(c.vk))
		if r1 != 0 {
			return true, c.mods, c.vk
		}
	}
	return false, 0, 0
}

func UnregisterGlobalHotKey(owner windows.HWND) {
	pUnregisterHotKey.Call(uintptr(owner), 0x5702)
}

// ---------------------------------------------------------------- job object

// AttachToKillOnCloseJob assigns the current process to a new Job Object
// with KILL_ON_JOB_CLOSE (the C30 semantics D32's tree accounting relies on).
func AttachToKillOnCloseJob() (windows.Handle, error) {
	job, err := windows.CreateJobObjectW(nil, nil)
	if err != nil {
		return 0, err
	}
	type basicLimitInfo struct {
		PerProcessUserTimeLimit int64
		PerJobUserTimeLimit     int64
		LimitFlags              uint32
		MinimumWorkingSetSize   uintptr
		MaximumWorkingSetSize   uintptr
		ActiveProcessLimit      uint32
		Affinity                uintptr
		PriorityClass           uint32
		SchedulingClass         uint32
	}
	type ioCounters struct {
		ReadOperationCount, WriteOperationCount uint64
		OtherOperationCount                     uint64
		ReadTransferCount, WriteTransferCount   uint64
		OtherTransferCount                      uint64
	}
	type extLimitInfo struct {
		BasicLimitInformation basicLimitInfo
		IoInfo                ioCounters
		ProcessMemoryLimit    uintptr
		JobMemoryLimit        uintptr
		PeakProcessMemoryUsed uintptr
		PeakJobMemoryUsed     uintptr
	}
	var li extLimitInfo
	li.BasicLimitInformation.LimitFlags = 0x2000 // JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	size := uint32(unsafe.Sizeof(li))
	if err := windows.SetInformationJobObject(job, 9, uintptr(unsafe.Pointer(&li)), size); err != nil {
		return 0, err
	}
	if err := windows.AssignProcessToJobObject(job, windows.CurrentProcess()); err != nil {
		return 0, err
	}
	return job, nil
}

// HiddenMessageWindow creates a hidden top-level window (owner of the tray
// icon and the hotkey) without any visible chrome.
func HiddenMessageWindow(title string) (windows.HWND, error) {
	clsName, _ := windows.UTF16PtrFromString("wisp_spike_msg")
	modh, _ := windows.GetModuleHandle(nil)
	cursor, _, _ := pLoadCursorW.Call(0, idcArrow)
	wc := wndClassEx{
		size:      uint32(unsafe.Sizeof(wndClassEx{})),
		wndProc:   windows.NewCallback(defWndProc),
		instance:  modh,
		cursor:    windows.Handle(cursor),
		bkgnd:     windows.Handle(colorWindow + 1),
		className: clsName,
	}
	if r1, _, err := pRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); r1 == 0 {
		return 0, fmt.Errorf("RegisterClassExW: %v", err)
	}
	t16, _ := windows.UTF16PtrFromString(title)
	hwnd, _, err := pCreateWindowExW(0,
		uintptr(unsafe.Pointer(clsName)), uintptr(unsafe.Pointer(t16)),
		0, 0, 0, 0, 0, 0, 0, uintptr(modh), 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("CreateWindowExW: %v", err)
	}
	return windows.HWND(hwnd), nil
}
