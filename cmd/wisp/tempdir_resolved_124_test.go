package main

// This file is ticket 124's AC#2b batch-4 conversion for cmd/wisp.
//
// Two test setups here name a fresh per-test directory and hand it to a store:
// newProvidersFixture runs it through memory.Open, and the secret failure-path
// cases run it through secretCmd's dataDir. Both hands reach the sealing floor.
// On an OS that spells the temp dir through a symlink (macOS /var ->
// /private/var, a container with TMPDIR=/varlink/...), the floor rightly
// refuses the unresolved spelling, so the case dies in setup and never reaches
// what it exists to pin. The root is resolved first through proc.SealableRoot,
// ticket 119's product-path discipline - the same delegation batches 1, 2, 3a
// installed in memory, config, tools, llm, perm, agent/approval and agent, and
// the same call this package's own resolveDataDir (doctor.go) and
// secretLayoutOf (secret.go) already make on the OS answer they receive.
//
// The production lines are not the issue and are not touched: they already
// resolve. What this batch fixes is the test that bypasses them by naming a
// dataDir directly, which is the only way a case can point the command at a
// throwaway tree.
//
// No assertion, threshold, refusal or golden in this package changed. Cases
// that are red in BOTH temp-dir shapes (the DPAPI family and the command-plane
// legs ticket 123 and ticket 119's R-119-7 booked) are out of scope here and
// stay exactly as red as they were.

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
