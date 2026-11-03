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

func (j *childJob) Kill() {}

func (j *childJob) Close() {}
