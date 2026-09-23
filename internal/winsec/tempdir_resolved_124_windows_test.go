//go:build windows

package winsec

import "testing"

// SealableTempDirForTest124 is a documented no-op here, and that is load-bearing.
//
// Six of the seventeen cases ticket 124 AC#2b batch 3b converts live in
// private_fail_test.go, which carries no build constraint and therefore also runs
// on Windows, where this package's whole host denominator lives. Resolving a
// Windows temp path is not a spelling-preserving operation: the runner's own path
// is an 8.3 short name for a longer profile directory, which is exactly why
// internal/winsec/export_test.go carries TreeOwnershipProbeForTest instead of
// letting a laptop reproduce that shape. Handing those six cases a short-name
// spelling they never asked for would move host and CI readings under a ticket
// that is about a POSIX harness shape.
//
// The other eleven cases are behind //go:build !windows and never compile on this
// platform, so this file exists only to keep the six above compiling and reading
// exactly what they read before. The resolving twin is
// tempdir_resolved_124_other_test.go.
func SealableTempDirForTest124(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}
