//go:build !windows

package main

import (
	"fmt"
	"os"
)

// cmdSLO refuses to run off Windows, with a nonzero exit.
//
// WHY THIS IS A REFUSAL AND NOT A NO-OP. `wisp slo` does not report a
// property of the code under test; it reports a property of a *running
// Windows resident process*: the Job Object created by proc.Boot, the
// per-session single-instance mutex, and RSS/private-bytes read out of the
// process tree by proc.NewTreeSampler / proc.NewExternalSampler for a specific
// PID (ticket 66's out-of-tree observer). None of those handles exist on any
// other platform - internal/proc's whole measurement surface is
// //go:build windows.
//
// The tempting alternative was to leave slo.go untagged and give internal/proc
// a portable Runtime whose samplers return zeroed structs. That would have
// turned a compile error into a passing SLO gate whose every number is a
// fabrication, which is strictly worse than the red text (ticket 78 explicitly
// forbids it). So the measured path is tagged away instead: on !windows the
// instrument does not exist, and this function says so.
//
// Exit 2 is the same fail-closed code slo_windows.go already returns when a
// subject cannot be started, measured, or cannot hand back its own report
// ("silently degrading to the old in-tree basis is exactly what the ticket
// exists to prevent"). Callers - scripts/slo-check.ps1 and the CI SLO jobs -
// therefore see "no measurement" on every platform that cannot measure, rather
// than a report.
func cmdSLO(args []string) int {
	fmt.Fprintln(os.Stderr, "wisp slo: the SLO gate measures a Windows resident process (Job Object + per-PID tree sampler);")
	fmt.Fprintln(os.Stderr, "wisp slo: this platform has nothing to measure, so the command is refusing instead of reporting numbers.")
	fmt.Fprintln(os.Stderr, "wisp slo: exit 2 (fail closed) - not a test failure, and not a pass.")
	_ = args
	return 2
}
