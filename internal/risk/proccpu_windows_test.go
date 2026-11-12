//go:build windows

package risk

import (
	"time"

	"golang.org/x/sys/windows"
)

// processCPUTime reports the user+kernel CPU this process has actually burned.
// It is the instrument the C26 per-call budget reads against: it does not
// advance while the process sits descheduled, so a measurement window that a
// concurrently running package stole CPU from cannot inflate it (ticket 86).
//
// Same two syscalls internal/speech/cgo_boundary_windows_test.go already uses
// for the same purpose; restated here so internal/risk's test build owns its
// own instrument.
func processCPUTime() (time.Duration, error) {
	proc := windows.CurrentProcess()
	var created, exited, kernel, user windows.Filetime
	if err := windows.GetProcessTimes(proc, &created, &exited, &kernel, &user); err != nil {
		return 0, err
	}
	return time.Duration(kernel.Nanoseconds() + user.Nanoseconds()), nil
}
