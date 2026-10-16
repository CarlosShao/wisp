// Ticket 322, leg 322-r1, AC#2. Untagged on purpose: the promise pinned here lives in
// the portable open (winsec.go:87 -> O_WRONLY|O_CREATE|O_TRUNC), so it is measured on
// BOTH platforms - container/Linux and this host's Windows - rather than asserted for one
// and assumed for the other.
//
// WHY THESE PINS DID NOT EXIST. 322-a1 read every production caller of that branch and
// measured two things: four of them depend on "or replaces", and NOT ONE case anywhere
// in the repository calls PrivateFile twice on the same path. The promise was live and
// unpinned, so 乙's cheaper variant (refuse when the name exists) could have silently
// broken all four callers with the suite still green. This file gives each caller its
// own world and its own assertion, so a refusal-for-existing-name goes red on the name of
// the site it broke.
//
// POSITIVE CONTROL (see the evidence file for the run): delete os.O_TRUNC from
// winsec.go:87 and every sub-case below must go red, because each one checks that the
// new bytes replaced the old ones rather than being laid over them.
package winsec

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// privateModeIsReadable reports whether this platform can be ASKED for a mode at all.
// On Windows the mode argument is decorative - that is this package's whole reason for
// existing - so 0600 can never be the read-back there, and the promise that holds on
// this platform is the security descriptor, which acl_windows_test.go measures with
// icacls. Only the mode read-back is split by platform; the rewrite landing and not
// being refused is asserted on both, and that is the half the four callers depend on.
var privateModeIsReadable = runtime.GOOS != "windows"

func TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1(t *testing.T) {
	// Each row is one caller's actual shape, not a generic "write twice". The FIRST body
	// is always longer than the replacement: the control run that deletes os.O_TRUNC from
	// winsec.go:87 left a tail of the old bytes behind, and a pair of equal-length bodies
	// would have stayed green through it (measured - the first draft of this pin used two
	// 15-character values and passed with no truncation at all, which is a pin with no
	// teeth, so every row here is now asymmetric by length).
	for _, tc := range []struct{ label, promise, first, second, file string }{
		{
			"internal/secret/store.go:79 Store.Set", "storing the same id twice replaces the blob in place",
			"secret-value-one, longer than the replacement", "secret-value-two", "blob-322r1.bin",
		},
		{
			"internal/config/migrate.go:93", "this call is the repair path for a wide backup left by an older build",
			"a wide backup left by an older build, longer than the replacement", "narrow", "config-backup-322r1.toml",
		},
		{
			"internal/secret/migrate.go:189", "a pre-rename temp left by a crash is overwritten, not refused",
			"tmp body from the run that crashed", "tmp body from this run", "wisp-secrets-322r1.toml.tmp",
		},
		{
			"internal/agent/spill.go:166", "a leftover .retry temp is overwritten (D31 last-writer-wins)",
			"first spill that needed a retry", "retry spill bytes", "artifact-322r1.bin.retry",
		},
	} {
		t.Run(tc.label, func(t *testing.T) {
			dir := SealableTempDirForTest124(t)
			p := filepath.Join(dir, tc.file)
			// Every world starts from an ENTRY THAT ALREADY EXISTS, which is the whole
			// question: an existing name must not become a refusal.
			if err := os.WriteFile(p, []byte(tc.first), 0o600); err != nil {
				t.Fatalf("pre-stage the existing entry: %v", err)
			}
			if err := PrivateFile(p, []byte(tc.second), 0o600); err != nil {
				t.Errorf("AC#2 RED: PrivateFile refused an existing name (%v), which breaks %s", err, tc.promise)
			}
			if got, err := os.ReadFile(p); err != nil || string(got) != tc.second {
				t.Errorf("AC#2 RED: content is %q (err %v), want %q - %s", got, err, tc.second, tc.promise)
			}
			info, err := os.Lstat(p)
			if err != nil {
				t.Fatalf("read the entry back: %v", err)
			}
			if got := info.Mode().Perm(); privateModeIsReadable && got != 0o600 {
				t.Errorf("AC#2 RED: mode is %v, want the replacement to land private (0600)", got)
			} else if !privateModeIsReadable {
				t.Logf("AC#2' mode read-back here is %v and is decorative by design; the descriptor is the "+
					"promise on this platform, and acl_windows_test.go owns that leg", got)
			}
		})
	}
	// The Windows colour of the same promise is not assumed: on this platform the mode
	// argument is decorative (the package doc's whole reason for existing), so the claim
	// that survives here is the one that matters for the four callers - the rewrite lands
	// and the entry stays private by descriptor, not by mode. This case states it rather
	// than borrowing the POSIX read-back.
	t.Run("the entry stays private after the rewrite", func(t *testing.T) {
		p := filepath.Join(SealableTempDirForTest124(t), "private-after-rewrite-322r1.bin")
		if err := PrivateFile(p, []byte("OLD bytes longer than the repair"), 0o644); err != nil {
			t.Fatalf("first write: %v", err)
		}
		if err := PrivateFile(p, []byte("NEW"), 0o644); err != nil {
			t.Errorf("AC#2 RED: rewrite of a name the caller itself created was refused: %v", err)
		}
		if got, _ := os.ReadFile(p); string(got) != "NEW" {
			t.Errorf("AC#2 RED: content is %q, want NEW", got)
		}
		info, err := os.Lstat(p)
		if err != nil {
			t.Fatal(err)
		}
		if privateModeIsReadable && info.Mode().Perm() != 0o600 {
			t.Errorf("AC#2 RED: a rewrite asked for 0644 landed %v, want the private answer", info.Mode().Perm())
		} else if !privateModeIsReadable {
			t.Logf("AC#2' read-back here is %v (decorative by design); the descriptor is this platform's promise",
				info.Mode().Perm())
		}
	})
	// The other half of the split, already pinned for the exclusive branch by
	// TestAC5FailureIsNotSwallowedByTheHappyPath (private_fail_test.go:108). What is NEW
	// here is that the leaf flag added by this ticket must not have moved that answer:
	// O_EXCL against an occupied REGULAR name still has to read fs.ErrExist, not ELOOP
	// and not nil.
	t.Run("the exclusive branch still refuses an occupied regular name", func(t *testing.T) {
		p := filepath.Join(SealableTempDirForTest124(t), "artifact-322r1.bin")
		if err := PrivateFileExclusive(p, []byte("FIRST")); err != nil {
			t.Fatalf("first exclusive create: %v", err)
		}
		if err := PrivateFileExclusive(p, []byte("SECOND")); !errors.Is(err, fs.ErrExist) {
			t.Errorf("AC#2 RED: the second exclusive create answered %v, want fs.ErrExist kept", err)
		}
		if got, _ := os.ReadFile(p); string(got) != "FIRST" {
			t.Errorf("AC#2 RED: the refused exclusive create replaced the artifact with %q", got)
		}
	})
}
