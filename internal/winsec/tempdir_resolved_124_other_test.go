//go:build !windows

package winsec

import "testing"

// SealableTempDirForTest124 hands a case the directory t.TempDir() just made,
// spelled the way the filesystem reaches it, so the floor is never asked to seal
// a root that names itself through the harness's own link. This is ticket 124
// AC#2b batch 3b: 17 internal/winsec cases measure red under a symlinked TMPDIR
// and green under a plain one, and every one of the 17 refusals names /varlink,
// the link the harness built, never a link the case planted.
//
// It deliberately does NOT call internal/proc's SealableRoot, even though nothing
// here stops compiling if it did - measured on a pristine copy of the anchor,
// internal/proc imports only context, errors, fmt and internal/buildinfo, so
// there is no import cycle and go vet is rc=0 either way. The reason is a written
// boundary, not a linker limit: internal/proc/envfork.go:148 states "nothing in
// internal/winsec calls it", and internal/winsec/resolve.go:224-228 restates the
// same rule from this side, because the floor must not start depending on the
// layer above it or that sentence stops being checkable. Ticket 125 met this
// exact problem already and answered it by walking in-package (resolveProbeRoot),
// so this reuses that answer instead of writing a fourth copy of the walk.
//
// The helper is fixture INPUT, which is what keeps it clear of R-119-9 (ticket
// 119): that rule forbids deriving a case's EXPECTED value from the function
// under test, including idempotent combinations of it. Nothing converted with
// this helper does that, and the cases that compare the floor's answer against a
// declared string still compare it against a declared string.
//
// On Windows this is a separate, do-nothing file - see
// tempdir_resolved_124_windows_test.go for why the host half must not move.
func SealableTempDirForTest124(t *testing.T) string {
	t.Helper()
	return resolveProbeRoot(t.TempDir())
}
