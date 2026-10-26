//go:build windows

package speech

import (
	"syscall"
	"unsafe"
)

// processMemoryCounters26 mirrors PROCESS_MEMORY_COUNTERS from psapi.h. Only
// WorkingSetSize is read; the rest are there so the struct is the size the API
// expects (cb is checked by the callee).
type processMemoryCounters26 struct {
	cb                         uint32
	pageFaultCount             uint32
	peakWorkingSetSize         uintptr
	workingSetSize             uintptr
	quotaPeakPagedPoolUsage    uintptr
	quotaPagedPoolUsage        uintptr
	quotaPeakNonPagedPoolUsage uintptr
	quotaNonPagedPoolUsage     uintptr
	pagefileUsage              uintptr
	peakPagefileUsage          uintptr
}

// workingSetBytes reads this process's working set: the RSS reading SPEC-04's
// resource budget is written in. It reports ok=false rather than failing a test
// when the API is unavailable, so the same test body runs everywhere.
func workingSetBytes() (int64, bool) {
	proc := syscall.NewLazyDLL("psapi.dll").NewProc("GetProcessMemoryInfo")
	var m processMemoryCounters26
	m.cb = uint32(unsafe.Sizeof(m))
	h, err := syscall.GetCurrentProcess()
	if err != nil {
		return 0, false
	}
	r, _, _ := proc.Call(uintptr(h), uintptr(unsafe.Pointer(&m)), uintptr(m.cb))
	if r == 0 {
		return 0, false
	}
	return int64(m.workingSetSize), true
}
