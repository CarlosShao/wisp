package agent

// This file is ticket 124's AC#2b batch-3a conversion for internal/agent.
//
// The spill, guard, forensics, truncation and golden-loop tests hand a fresh
// per-test directory either to memory.Open or to NewSpiller, and both hands
// reach the sealing floor. On an OS that spells the temp dir through a symlink
// (macOS /var -> /private/var, a container with TMPDIR=/varlink/...), the floor
// rightly refuses the unresolved spelling, so the test dies in setup and never
// reaches the behaviour it exists to pin. The root is resolved first through
// proc.SealableRoot, ticket 119's product-path discipline - the same delegation
// batches 1 and 2 installed in memory, config, tools, llm, perm and
// agent/approval. No assertion, threshold, or refusal in this package changed.

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
