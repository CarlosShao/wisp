package tools_test

// This file is the external-test-package twin of
// tempdir_resolved_124_test.go (ticket 124 AC#2b batch 2). loop_approval_test.go
// lives in package tools_test, which cannot see the inner package's helpers, so
// the same one-line delegation to proc.SealableRoot - ticket 119's product-path
// discipline, resolve what the OS answered before handing a root down - is
// spelled here once. It is not a second resolution algorithm: the resolution
// itself stays in internal/proc.

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
