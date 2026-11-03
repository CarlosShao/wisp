//go:build windows

package tools

import "golang.org/x/sys/windows"

// DPI awareness for the window-measuring verbs. Measuring and setting window
// rectangles while the process is DPI-unaware makes Windows scale the numbers
// behind our back, so a "move to x=100" lands somewhere else and a captured
// window rect is off by the scale factor. Declaring per-monitor awareness once
// at first use is the documented fix (SetProcessDpiAwareness, shcore.dll).

// setDPIAwarenessSafe declares per-monitor DPI awareness once per process. The
// proc is absent on Windows 8, where PROCESS_DPI_UNAWARE is the only behaviour
// and the callers accept physical pixels - so absence is reported, never fatal.
func setDPIAwarenessSafe() error {
	if procSetProcessDpiAwareness.Find() != nil {
		return windows.ERROR_PROC_NOT_FOUND
	}
	// PROCESS_PER_MONITOR_DPI_AWARE = 2
	const perMonitor = 2
	r, _, e := procSetProcessDpiAwareness.Call(perMonitor)
	// E_ACCESSDENIED is the "already declared" answer (a host set it first),
	// which is the desired end state rather than a failure.
	if r == 0 && e != nil {
		if errno, ok := e.(windows.Errno); ok && errno == windows.ERROR_ACCESS_DENIED {
			return nil
		}
		return e
	}
	return nil
}
