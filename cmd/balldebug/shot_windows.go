//go:build windows

package main

// Go-side composite capture (debug evidence): BitBlt from the screen DC into
// a private top-down 32bpp DIB, then encode PNG. Equivalent to the PowerShell
// CopyFromScreen path but self-contained (works when console tooling is
// unreliable). Screen DC captures the COMPOSITED desktop (ball + background),
// which is what the SPEC-08 §2.1 human review needs.

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modUser32S = windows.NewLazySystemDLL("user32.dll")
	modGDI32S  = windows.NewLazySystemDLL("gdi32.dll")

	pGetWindowRectS  = modUser32S.NewProc("GetWindowRect")
	pGetDCS          = modUser32S.NewProc("GetDC")
	pReleaseDCS      = modUser32S.NewProc("ReleaseDC")
	pCreateCompatDC  = modGDI32S.NewProc("CreateCompatibleDC")
	pDeleteDCS       = modGDI32S.NewProc("DeleteDC")
	pCreateDIBSect   = modGDI32S.NewProc("CreateDIBSection")
	pSelectObjectS   = modGDI32S.NewProc("SelectObject")
	pDeleteObjectS   = modGDI32S.NewProc("DeleteObject")
	pBitBltS         = modGDI32S.NewProc("BitBlt")
	pSetWindowPosS   = modUser32S.NewProc("SetWindowPos")
	pGetModuleHandle = windows.NewLazySystemDLL("kernel32.dll").NewProc("GetModuleHandleW")
)

type shotRect struct{ l, t, r, b int32 }

type shotBMI struct {
	header struct {
		size, width, height    uint32
		planes, bitCount       uint16
		compression, sizeImage uint32
		xppm, yppm             uint32
		clrUsed, clrImportant  uint32
	}
	_ uint32
}

// goScreenShot captures the composited screen at the window rect to path.
func goScreenShot(hwnd windows.HWND, path string) {
	var wr shotRect
	r, _, _ := pGetWindowRectS.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&wr)))
	if r == 0 {
		fmt.Printf("go-shot: GetWindowRect failed\n")
		return
	}
	w, h := int(wr.r-wr.l), int(wr.b-wr.t)
	if w < 4 || h < 4 {
		fmt.Printf("go-shot: bad rect %d,%d %dx%d\n", wr.l, wr.t, w, h)
		return
	}

	hdcScreen, _, _ := pGetDCS.Call(0)
	if hdcScreen == 0 {
		fmt.Printf("go-shot: GetDC failed\n")
		return
	}
	defer pReleaseDCS.Call(0, hdcScreen)

	hdcMem, _, _ := pCreateCompatDC.Call(hdcScreen)
	if hdcMem == 0 {
		fmt.Printf("go-shot: CreateCompatibleDC failed\n")
		return
	}
	defer pDeleteDCS.Call(hdcMem)

	bmi := shotBMI{}
	bmi.header.size = 40
	bmi.header.width = uint32(w)
	bmi.header.height = uint32(h) // bottom-up: BitBlt writes bottom-up rows
	bmi.header.planes = 1
	bmi.header.bitCount = 32
	var bits uintptr
	hbm, _, err := pCreateDIBSect.Call(hdcMem, uintptr(unsafe.Pointer(&bmi)), 0, uintptr(unsafe.Pointer(&bits)), 0, 0)
	if hbm == 0 {
		fmt.Printf("go-shot: CreateDIBSection: %v\n", err)
		return
	}
	defer pDeleteObjectS.Call(hbm)
	prev, _, _ := pSelectObjectS.Call(hdcMem, hbm)
	defer pSelectObjectS.Call(hdcMem, prev)

	// SRCCOPY|CAPTUREBLT: layered windows are only composited into screen DC
	// captures when CAPTUREBLT is set.
	rr, _, _ := pBitBltS.Call(hdcMem, 0, 0, uintptr(w), uintptr(h), hdcScreen,
		uintptr(wr.l), uintptr(wr.t), 0x00CC0020|0x40000000)
	if rr == 0 {
		fmt.Printf("go-shot: BitBlt failed\n")
		return
	}

	// Bottom-up BGRA -> top-down RGBA PNG (no un-premultiply: the screen is
	// already composited).
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	src := (*[1 << 26]uint8)(winHeapPtrLocal(bits))
	row := w * 4
	for y := 0; y < h; y++ {
		srcRow := (h - 1 - y) * row // bottom-up source
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, color.RGBA{R: src[srcRow+x*4+2], G: src[srcRow+x*4+1], B: src[srcRow+x*4], A: 255})
		}
	}
	f, err := os.Create(path)
	if err == nil {
		err = png.Encode(f, img)
		f.Close()
	}
	if err != nil {
		fmt.Printf("go-shot: %v\n", err)
		return
	}
	fmt.Printf("go-shot: %s (%dx%d at %d,%d)\n", baseName(path), w, h, wr.l, wr.t)
}

func baseName(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '\\' || p[i] == '/' {
			return p[i+1:]
		}
	}
	return p
}

func winHeapPtrLocal(p uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&p))
}

// debugMoveWindow places the ball window without activation (debug flag).
func debugMoveWindow(hwnd windows.HWND, x, y int32) {
	const swpNoSize = 0x0001
	const swpNoActivate = 0x0010
	pSetWindowPosS.Call(uintptr(hwnd), 0, uintptr(x), uintptr(y), 0, 0, swpNoSize|swpNoActivate)
}
