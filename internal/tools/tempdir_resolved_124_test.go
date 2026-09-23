package tools

// This file is ticket 124's AC#2b batch-2 conversion for internal/tools.
//
// The shape under test here is the one ticket 119 already fixed on the product
// path: when the harness OS spells the temp dir through a symlink (macOS
// /var -> /private/var, a container with TMPDIR=/varlink/...), handing that
// unresolved spelling to the sealing floor - or to the C26 allowlist judgment -
// is refused, because an ancestor is a link. The refusal is correct behaviour
// for a real path nobody declared; it is a harness bug when the path came from
// t.TempDir() itself. So the root handed to the pipeline is resolved first,
// through the same proc.SealableRoot the product's TestDataDir /
// DefaultLayout use. Nothing in winsec or the risk pipeline changed: only
// where the tests' roots come from.

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/proc"
)

// sealableTempDir124 is the resolved-OS-spelling sibling of tempRaw: a fresh
// per-test directory with the existing symlink prefix resolved away, still in
// the spelling the OS can open. Tests that hand a root to the bridge (which
// canonicalizes it themselves) use this instead of t.TempDir().
func sealableTempDir124(t *testing.T) string {
	t.Helper()
	return proc.SealableRoot(t.TempDir())
}

// sealableTempCanonical124 is the sibling of tempCanonical: same C26 run on
// top (those tests compare against the canonical form by design), but the OS
// answer that goes INTO mustCanonical is resolved first, per ticket 124.
func sealableTempCanonical124(t *testing.T) string {
	t.Helper()
	return mustCanonical(t, sealableTempDir124(t))
}
