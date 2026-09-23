package llm_test

// This file is ticket 124's AC#2b batch-2 conversion for internal/llm.
//
// The probe-suite fixtures open a real memory.Store under t.TempDir(); on an
// OS that spells the temp dir through a symlink (macOS /var -> /private/var,
// a container with TMPDIR=/varlink/...), the sealing floor rightly refuses the
// unresolved spelling and the fixture never opens. The root is now resolved
// through proc.SealableRoot - ticket 119's product-path discipline - before it
// is handed down. Assertions are untouched; only where the root comes from.

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
