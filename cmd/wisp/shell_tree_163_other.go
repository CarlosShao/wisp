//go:build !windows

package main

// The shell family's tree seam, non-Windows half (ticket 163's
// composition-root wiring). Wisp is Windows-only (D7), but this file exists so
// the untagged assembly keeps type-checking on other GOOSes - and so the
// failure there is LOUD, never a silent zero-value: a shell command started
// outside any kill-on-close Job containment is the exact posture ticket
// 100/R21 recorded as not landing. No process is the fail-closed answer.

import (
	"errors"

	"github.com/CarlosShao/wisp/internal/tools"
)

// newShellSessionTree refuses: there is no tree provider on this GOOS, so no
// shell process may start uncontained.
func newShellSessionTree() (tools.SessionTree, error) {
	return nil, errors.New("shell trees are Windows-only in this assembly (D7): no kill-on-close Job provider on this GOOS, refusing to start an uncontained process")
}
