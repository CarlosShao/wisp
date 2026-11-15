//go:build !windows

package main

import "fmt"

// The non-Windows half of the disk-guard probe, stated rather than omitted: on a
// platform this package has no build for (AGENTS.md's logsink.go note: cmd/wisp
// has no Linux build because of the sherpa import main.go carries), there is no
// reading to take, and a guard that cannot read must not be installed - the
// observable consequence is a nil LogConfig.Disk, i.e. the pre-ticket-45 log
// writer, and installLogSink books that as its note. It never answers 0 free
// bytes, which the guard would read as "the disk is full" and hold every line.
func dataVolumeFreeBytes(volumePath string) (int64, error) {
	return 0, fmt.Errorf("wisp: no disk-free probe on this platform (asked about %q)", volumePath)
}
