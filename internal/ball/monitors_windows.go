//go:build windows

package ball

// Monitor enumeration for multi-monitor positioning (D42#3 topology). The
// window itself is Per-Monitor V2 DPI aware (SetProcessDpiAwarenessContext
// called before any window creation); WM_DPICHANGED re-scales the render
// target on the fly (DPI resource rebuild stubbed to full re-render for S1,
// per the ticket).

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// enablePerMonitorV2DPI opts the process into Per-Monitor V2 DPI awareness.
// Must run before the first window is created; failure (older OS, already
// set by the manifest) is non-fatal.
func enablePerMonitorV2DPI() error {
	r, _, err := pSetProcessDpiAwarenessContext.Call(dpiAwarenessPerMonitorV2)
	if r == 0 {
		return fmt.Errorf("SetProcessDpiAwarenessContext: %v", err)
	}
	return nil
}

// getDpiForWindow reads the window's current monitor DPI (per-monitor V2).
func getDpiForWindow(hwnd windows.HWND) uint32 {
	d, _, _ := pGetDpiForWindow.Call(uintptr(hwnd))
	if d == 0 {
		return 96
	}
	return uint32(d)
}

// enumMonitors lists the connected monitors (screen + work rects + device
// name + primary flag) for ResolvePosition.
func enumMonitors() []MonitorRect {
	var out []MonitorRect
	cb := windows.NewCallback(func(hMon uintptr, hdc uintptr, rc uintptr, data uintptr) uintptr {
		var mi monitorInfoExW
		mi.size = uint32(unsafe.Sizeof(mi))
		if r, _, _ := pGetMonitorInfoW.Call(hMon, uintptr(unsafe.Pointer(&mi))); r != 0 {
			out = append(out, MonitorRect{
				Device:  windows.UTF16ToString(mi.device[:]),
				Primary: mi.flags&1 != 0, // MONITORINFOF_PRIMARY
				Screen:  Rect{int(mi.rcMonitor.l), int(mi.rcMonitor.t), int(mi.rcMonitor.r), int(mi.rcMonitor.b)},
				Work:    Rect{int(mi.rcWork.l), int(mi.rcWork.t), int(mi.rcWork.r), int(mi.rcWork.b)},
			})
		}
		return 1 // continue enumeration
	})
	pEnumDisplayMonitors.Call(0, 0, cb, 0)
	return out
}

// monitorFromWindow returns the device name of the monitor the window is
// (mostly) on.
func monitorFromWindow(hwnd windows.HWND) string {
	h, _, _ := pMonitorFromWindow.Call(uintptr(hwnd), monitorDefaultToNearest)
	if h == 0 {
		return ""
	}
	var mi monitorInfoExW
	mi.size = uint32(unsafe.Sizeof(mi))
	if r, _, _ := pGetMonitorInfoW.Call(h, uintptr(unsafe.Pointer(&mi))); r != 0 {
		return windows.UTF16ToString(mi.device[:])
	}
	return ""
}

// moveWindow places the window at (x, y) sized (w, h) without activation.
func moveWindow(hwnd uintptr, x, y, w, h int32) {
	pSetWindowPos.Call(uintptr(hwnd), 0, // HWND_TOP (keep topmost via SWP_NOZORDER? we WANT topmost retained)
		uintptr(x), uintptr(y), uintptr(w), uintptr(h),
		uintptr(swpNoActivate|swpShowWindow))
}

// foregroundRect returns the virtual screen bounds (for off-screen checks).
