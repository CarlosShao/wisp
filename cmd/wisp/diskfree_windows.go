//go:build windows

package main

// The platform half of the D42#5 disk guard (ticket 45's guards, ticket 46's
// seam leg): the free-space reading. internal/observe owns the DECISION and
// deliberately owns no platform call (guards_45.go's header: "the cmd wiring
// owns the platform call"), so the probe lives here, at the composition root,
// next to the only site that installs it (installLogSink).
//
// The call is GetDiskFreeSpaceEx, already reachable through
// golang.org/x/sys/windows, which go.mod pins. No new dependency.

import (
	"fmt"

	"golang.org/x/sys/windows"
)

// dataVolumeFreeBytes reads the free bytes available to this user on the volume
// that holds volumePath. volumePath is a DIRECTORY inside the volume (the data
// root or its log directory) - GetDiskFreeSpaceEx resolves the volume from it,
// so nothing here has to guess a drive letter or a mount point.
func dataVolumeFreeBytes(volumePath string) (int64, error) {
	p, err := windows.UTF16PtrFromString(volumePath)
	if err != nil {
		return 0, fmt.Errorf("wisp: disk guard probe path %q: %w", volumePath, err)
	}
	var freeToCaller, total, totalFree uint64
	if err := windows.GetDiskFreeSpaceEx(p, &freeToCaller, &total, &totalFree); err != nil {
		return 0, fmt.Errorf("wisp: GetDiskFreeSpaceEx(%q): %w", volumePath, err)
	}
	if freeToCaller > uint64(^uint64(0)>>1) {
		// A reading that cannot be an int64 is not a reading: refuse it rather
		// than wrap it into a negative "free" number the guard would read as
		// "below threshold" forever.
		return 0, fmt.Errorf("wisp: GetDiskFreeSpaceEx(%q) answered %d bytes, out of range", volumePath, freeToCaller)
	}
	return int64(freeToCaller), nil
}
