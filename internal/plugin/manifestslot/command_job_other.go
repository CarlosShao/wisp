//go:build !windows

package manifestslot

import (
	"os/exec"
)

// childJob is the non-Windows placeholder: without a Job Object the tree kill
// degrades to killing the direct child, and this file says so instead of
// pretending otherwise. Wisp ships Windows-only (D7); this exists so the package
// still builds under the cross-vet leg.
type childJob struct{ exited uintptr }

func openChildJob() (*childJob, error) {
	return nil, errNoJobObject
}

func (j *childJob) Start(cmd *exec.Cmd) error { return cmd.Start() }

func (j *childJob) WaitExited() {}

// Kill is a no-op here because it is UNREACHABLE, not because the tree kill is
// silent: osChild.job is only set when openChildJob succeeded, and on this
// platform it never does (openChildJob returns errNoJobObject). The reachable
// degradation notice lives in OSCommandSpawner.Spawn's fallback (command.go),
// which prints what the kill will and will not reach before the child starts.
func (j *childJob) Kill() {}

func (j *childJob) Close() {}
