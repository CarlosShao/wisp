package memory

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/proc"
)

// sealableTempDir124 is t.TempDir() run through the data-root discipline
// ticket 119 settled for the production route: the layer that reads the OS
// resolves what the OS answered before that spelling is handed to the sealing
// floor (internal/proc/envfork.go's SealableRoot, the same call cmd/wisp
// resolveDataDir and doctor make).
//
// Every store this package opens goes through Open, and Open seals its
// directory with winsec. winsec's placement floor refuses a root that reaches
// itself through a symlink - that refusal is ticket 113's rule and it is
// supposed to stand. What it was hitting, though, was not a caller's declared
// root but the OS's own answer for TMPDIR, which is spelled through a link on
// ordinary machines (/var on macOS, a linked TMPDIR in a container). Ticket 124
// counted the result: the store-backed cases here failed on the harness's temp
// spelling instead of on whatever each case asserts.
//
// This helper adds no resolution of its own - it is one call to the existing
// leaf. A second walk over the same question is exactly what ticket 125's
// R-125-2 booked as debt, so it does not get a fourth copy here.
//
// Scope, and the boundary it respects: this package does not test the floor, so
// no case here wants an unresolved root. Cases that pin the refusal itself live
// in internal/winsec (tickets 113/118/119/125), and they keep building their
// own spelled-through-a-link roots untouched. Where a root is handed to a second
// process - concurrent_test.go's WISP_CRASH_DIR - the resolved spelling is what
// travels: a child given the unresolved one cannot open the store at all, and
// the parent reports that as "the subprocess committed 0 rows", which is a
// different failure from the half-written transaction the case is chasing.
func sealableTempDir124(t *testing.T) string {
	t.Helper()
	return proc.SealableRoot(t.TempDir())
}
