//go:build windows

package tools

import (
	"fmt"
	"image"
	"image/png"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Screen capture for ticket 23 (Windows).
//
// The whole file exists to produce PNG bytes that NEVER touch the disk. Two
// readers matter here and both are contract sentences:
//
//   - "截图只存内存 bytes_ref，不落盘、不进日志、不进 DB" (PLAN.md:2556). So this
//     file writes nothing, has no os import, and the bytes travel only in the
//     return value.
//   - "screen.capture 对 DRM 保护窗口截到黑屏" (PLAN.md:3236) is a PLATFORM
//     BEHAVIOUR: the capture of a protected window legitimately comes back
//     black, and the tool says so in its text instead of presenting a black
//     rectangle as the user's screen. The flag below is what carries that
//     fact; the detection is "every sampled pixel is (0,0,0)".
//
// The pixel source is GDI BitBlt into a 32-bit DIB section. PrintWindow with
// PW_RENDERFULLCONTENT would capture some layered windows that BitBlt misses,
// but it also renders direct-composition surfaces as black on several drivers,
// which would make the DRM signal indistinguishable from a plain failure - the
// one trade-off this ticket cannot afford, because AC#2 asserts them apart.

var (
	procCreateCompatibleDC     = gdi32.NewProc("CreateCompatibleDC")
	procCreateDIBSection       = gdi32.NewProc("CreateDIBSection")
	procSelectObject           = gdi32.NewProc("SelectObject")
	procDeleteObject           = gdi32.NewProc("DeleteObject")
	procDeleteDC               = gdi32.NewProc("DeleteDC")
	procBitBlt                 = gdi32.NewProc("BitBlt")
	procGetSystemMetrics       = user32.NewProc("GetSystemMetrics")
	procGetDeviceCaps          = gdi32.NewProc("GetDeviceCaps")
	procSetCaptureDpiAwareness = setDPIAwarenessSafe
)

var gdi32 = windows.NewLazySystemDLL("gdi32.dll")

// bitmapinfo mirrors BITMAPINFOHEADER's layout (40 bytes) as the DIB header.
type bitmapInfoHeader struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

const (
	biRGB         = 0
	dibRGBColors  = 0
	srcCopy       = 0x00CC0020
	titleBarPx    = 0
	captureMaxPx  = 8192 // D32-ish ceiling: refuse an absurd surface rather than allocate
	defaultScreen = 0
	SMWidth       = 0
	SMHeight      = 1
)

// realCaptureScreen grabs the whole virtual screen.
func realCaptureScreen() ([]byte, bool, error) {
	w := int(callInt(procGetSystemMetrics, SMWidth))
	h := int(callInt(procGetSystemMetrics, SMHeight))
	if w <= 0 || h <= 0 {
		return nil, false, fmt.Errorf("读取屏幕尺寸失败")
	}
	src := callUintptr(procGetDC, 0)
	if src == 0 {
		return nil, false, fmt.Errorf("取屏幕 DC 失败")
	}
	defer callUintptr(procReleaseDC, 0, src)
	return copyToPNG(src, 0, 0, w, h)
}

// realCaptureWindow grabs one window's client-frame image. An iconic (minimized)
// window has no surface to copy, so that is a refusal with a readable reason
// instead of a black frame that would be mistaken for DRM.
func realCaptureWindow(hwnd uintptr) ([]byte, bool, error) {
	if hwnd == 0 {
		return nil, false, fmt.Errorf("没有可截图的窗口")
	}
	if callBool(procIsIconic, hwnd) {
		return nil, false, fmt.Errorf("目标窗口已最小化，没有可截取画面（先 window.manage restore）")
	}
	r, ok := realWindowRect(hwnd)
	if !ok {
		return nil, false, fmt.Errorf("读不到窗口矩形")
	}
	w := int(r.Right - r.Left)
	h := int(r.Bottom - r.Top)
	if w <= 0 || h <= 0 {
		return nil, false, fmt.Errorf("窗口矩形为空（%dx%d）", w, h)
	}
	src := callUintptr(procGetWindowDC, hwnd)
	if src == 0 {
		return nil, false, fmt.Errorf("取窗口 DC 失败")
	}
	defer callUintptr(procReleaseDC, hwnd, src)
	return copyToPNG(src, int(r.Left), int(r.Top), w, h)
}

// copyToPNG BitBlts a rectangle of src into an offscreen DIB and PNG-encodes it.
// The DIB is destroyed before returning, so the only surviving copy of the
// pixels is the returned slice - which is what bytes_ref points at.
func copyToPNG(src uintptr, x, y, w, h int) ([]byte, bool, error) {
	if w > captureMaxPx || h > captureMaxPx {
		return nil, false, fmt.Errorf("截图区域 %dx%d 超过单次上限 %dx%d", w, h, captureMaxPx, captureMaxPx)
	}
	mem := callUintptr(procCreateCompatibleDC, src)
	if mem == 0 {
		return nil, false, fmt.Errorf("CreateCompatibleDC 失败")
	}
	defer callUintptr(procDeleteDC, mem)

	bi := bitmapInfoHeader{
		Size:        uint32(unsafe.Sizeof(bitmapInfoHeader{})),
		Width:       int32(w),
		Height:      -int32(h), // top-down DIB: row 0 is the top row
		Planes:      1,
		BitCount:    32,
		Compression: biRGB,
	}
	var bits unsafe.Pointer
	dib, _, _ := procCreateDIBSection.Call(mem, uintptr(unsafe.Pointer(&bi)), dibRGBColors,
		uintptr(unsafe.Pointer(&bits)), 0, 0)
	if dib == 0 || bits == nil {
		return nil, false, fmt.Errorf("CreateDIBSection 失败")
	}
	defer callUintptr(procDeleteObject, dib)
	old := callUintptr(procSelectObject, mem, dib)
	if old != 0 {
		defer callUintptr(procSelectObject, mem, old)
	}
	if !callBool(procBitBlt, mem, 0, 0, uintptr(w), uintptr(h), src, uintptr(x), uintptr(y), srcCopy) {
		return nil, false, fmt.Errorf("BitBlt 失败（该表面不可复制）")
	}

	rgba := unsafe.Slice((*byte)(bits), w*h*4)
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	allBlack := true
	for i := 0; i < w*h; i++ {
		b, g, r := rgba[i*4], rgba[i*4+1], rgba[i*4+2]
		if r != 0 || g != 0 || b != 0 {
			allBlack = false
		}
		img.Pix[i*4] = r
		img.Pix[i*4+1] = g
		img.Pix[i*4+2] = b
		img.Pix[i*4+3] = 0xFF
	}
	var buf bytesBuffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, false, fmt.Errorf("PNG 编码失败: %w", err)
	}
	// allBlack is reported as the DRM flag: a protected-content surface is
	// returned as an all-zero bitmap by the display driver (PLAN.md:3236). The
	// caller's text distinguishes it from an error; a genuinely black screen
	// reports the same annotation, which is honest - we cannot tell them apart
	// and the note says so.
	return buf.b, allBlack, nil
}

// bytesBuffer is a minimal in-memory writer for png.Encode (bytes.Buffer would
// do, but this keeps the file's import list to the packages that carry the
// point: image/png is the only encoder, and nothing here can reach a file).
type bytesBuffer struct{ b []byte }

func (w *bytesBuffer) Write(p []byte) (int, error) {
	w.b = append(w.b, p...)
	return len(p), nil
}

func callInt(proc *windows.LazyProc, args ...uintptr) int32 {
	return int32(callUintptr(proc, args...))
}
