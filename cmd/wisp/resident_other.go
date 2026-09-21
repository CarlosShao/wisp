//go:build !windows

package main

import (
	"fmt"
	"os"
)

// runResident refuses to start the resident process off Windows.
//
// The no-args path of this binary boots the Windows resident runtime: a Job
// Object that kills the whole process tree with the parent (C30), a
// per-session named mutex for single instance (D42#7), and the event loop that
// the floating ball window feeds. All three are //go:build windows in
// internal/proc, and there is no portable substitute for a job object or a
// named session mutex - so this platform has no resident process to start.
//
// Exiting 2 is the same code runResident uses on Windows when it cannot even
// resolve the environment: the caller learns "this did not run", never "it ran
// and did nothing". A `return` here would be the false green ticket 78 forbids.
func runResident() {
	fmt.Fprintln(os.Stderr, "wisp: the resident process requires Windows (Job Object + per-session single instance); refusing to start.")
	os.Exit(2)
}
