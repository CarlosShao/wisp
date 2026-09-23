package perm

// This file is ticket 124's AC#2b batch-2 conversion for internal/perm.
//
// The ticket-90 persistence tests write config.toml and open a memory.Store
// under t.TempDir(); both hands reach the sealing floor. On an OS that spells
// the temp dir through a symlink (macOS /var -> /private/var, a container with
// TMPDIR=/varlink/...), the floor rightly refuses the unresolved spelling, so
// the test never gets to its assertion. The root is resolved first through
// proc.SealableRoot, ticket 119's product-path discipline. No assertion,
// threshold, or refusal in this package changed.

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/proc"
)

// sealableTempDir124 returns a fresh per-test directory with every symlink in
// its existing prefix resolved away.
func sealableTempDir124(t *testing.T) string {
	t.Helper()
	return proc.SealableRoot(t.TempDir())
}
