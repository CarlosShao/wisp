//go:build windows

package adaptertest

import (
	"os/exec"
	"testing"

	"github.com/CarlosShao/wisp/internal/proc"
)

// StartChildGuarded starts cmd as an EXPLICIT MEMBER of a Windows Job Object
// created with JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE, reusing the C30 owner
// proc.JobScope rather than inventing a second job wrapper (ticket 304,
// shape A).
//
// Why this exists: t.Cleanup only runs when the test binary reaches its own
// teardown. Windows has no "parent dies, children die" rule, so a test binary
// killed by `go test -timeout` (a panic that skips cleanup) or by an external
// TerminateProcess (a leg killed by the fleet, a shell timeout, Ctrl+C) leaves
// its mockllm.exe child listening on a loopback port forever - ticket 304
// measured exactly one orphan per such death. A Job handle is closed by the OS
// when the owning process dies by ANY means, and KILL_ON_JOB_CLOSE then takes
// the members with it, so the reaping does not depend on any Go code reaching
// a teardown line.
//
// Scope is deliberately narrow: this job carries only the process started
// here, and nothing else is ever assigned to it (ticket 304 boundary 2 - no
// sweeping the surrounding process tree, no killing a process we did not
// start).
func StartChildGuarded(t *testing.T, cmd *exec.Cmd) error {
	t.Helper()
	job, err := proc.OpenJobScope()
	if err != nil {
		return err
	}
	// StartInJob starts AND assigns, and it kills plus waits the child itself
	// if the assignment fails, so a child can never escape the job by accident.
	if _, err := job.StartInJob(cmd); err != nil {
		_ = job.Close()
		return err
	}
	// Registered here rather than by the caller: closing this handle is what
	// does the killing, so it must be armed as soon as membership exists.
	// Cleanups run LIFO, so this closes before the caller's Kill/Wait step -
	// the child is already terminated by then and Kill's error is the
	// pre-existing "already dead" case the caller ignores.
	t.Cleanup(func() {
		if err := job.Close(); err != nil {
			t.Errorf("close mockllm job scope: %v", err)
		}
	})
	return nil
}
