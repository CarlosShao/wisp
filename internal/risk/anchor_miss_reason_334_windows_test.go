//go:build windows

package risk

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// partialCandidate334 is the Windows half of ticket 334's world: the candidate
// whose spelling the OS refuses to prove. hostile334 (anchor_world_334_test.go)
// carries U+001B / U+202E, which no Windows path component may contain, so
// Lstat fails with "The filename, directory name, or volume label syntax is
// incorrect" rather than "no such file" - and pathExists reads exactly that
// difference as present-but-unproven (partial=true), which is what makes
// uncertainAnchorMiss fire.
//
// The premise is asserted, not assumed. The clause is only exercised while the
// OS keeps refusing to parse this name: if the failure ever became a plain
// absence, partial would be false, the clause would go inert and the cases in
// anchor_miss_reason_334_test.go would stop measuring what they claim. That is
// not hypothetical - on POSIX both bytes are legal filename bytes and the same
// spelling IS a plain absence, which is why that platform's maker is a dangling
// symlink (anchor_miss_reason_339_other_test.go).
func partialCandidate334(t *testing.T, dir string) string {
	t.Helper()
	cand := filepath.Join(dir, hostile334)
	_, err := os.Lstat(cand)
	if err == nil {
		t.Fatalf("setup: %q must not exist before the maker runs", cand)
	}
	if errors.Is(err, os.ErrNotExist) {
		t.Fatalf("premise: the hostile spelling must be UNPARSEABLE, not a plain absent file - pathExists would read it as absent, partial would be false and the fail-closed clause would go inert: %v", err)
	}
	return cand
}
