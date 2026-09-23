package config

import (
	"testing"

	"github.com/CarlosShao/wisp/internal/proc"
)

// sealableTempDir124 is t.TempDir() run through the data-root discipline
// ticket 119 settled for the production route: the layer that reads the OS
// resolves what the OS answered before that spelling reaches the sealing floor
// (internal/proc/envfork.go's SealableRoot, the same call cmd/wisp
// resolveDataDir and doctor make).
//
// SaveFile seals the directory it writes into, so a config saved under an
// unresolved TMPDIR is refused before the bytes are ever compared: on macOS the
// OS answers TMPDIR through /var, and in a container through whatever link the
// harness put there. Ticket 124's ledger counts the cases in this package that
// were failing on that spelling rather than on what they assert.
//
// No second walk is added here - this is one call to the existing leaf. Ticket
// 125's R-125-2 already books the copies that exist and why merging them is not
// free, so this batch does not make that number worse.
//
// Deliberately NOT converted: c26_seam_posix_125_test.go. Those cases are ticket
// 125's seam probes; they look at the unresolved spelling on purpose and
// t.Skipf themselves when the harness's temp base reaches itself through a link,
// because a probe that silently ran on a different tree would be worth nothing.
// Resolving their base would turn a refusal-to-measure into a measurement of the
// wrong shape.
func sealableTempDir124(t *testing.T) string {
	t.Helper()
	return proc.SealableRoot(t.TempDir())
}
