//go:build !windows

package adaptertest

import (
	"os/exec"
	"testing"
)

// StartChildGuarded on a non-Windows platform is a PLAIN cmd.Start().
//
// This is NOT a guard and must not be read as one: the internal/proc JobScope
// this stands in for on Windows is behind //go:build windows (it is the C30
// Windows Job Object wrapper), so there is no equivalent handle to close here.
// The ubuntu/linux leg of the matrix therefore STILL LEAKS the child exactly
// the way ticket 304 measured for Windows before this change - a test binary
// killed by `go test -timeout` or by an external signal leaves its mockllm
// child running, because nothing but t.Cleanup is watching it. Ticket 304's
// fix does not cover that branch; it is named as an open reading
// (Q-304-r1-1), not as a solved case.
func StartChildGuarded(t *testing.T, cmd *exec.Cmd) error {
	t.Helper()
	return cmd.Start()
}
