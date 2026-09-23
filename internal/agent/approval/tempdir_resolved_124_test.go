package approval_test

// This file is ticket 124's AC#2b batch-2 conversion for
// internal/agent/approval.
//
// The batch test names its allowlist root with t.TempDir(); on an OS that
// spells the temp dir through a symlink (macOS /var -> /private/var, a
// container with TMPDIR=/varlink/...), the C26 judgment cannot confirm that
// tree on disk, authorizes nothing, and a declared L1 write silently escalates
// to L2 - which is the opposite of what the test pins. The root is resolved
// first through proc.SealableRoot, ticket 119's product-path discipline. The
// gate's own refusals are untouched; only where the root comes from.

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
