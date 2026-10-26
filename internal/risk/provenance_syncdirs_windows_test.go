//go:build windows

package risk

import (
	"path/filepath"
	"testing"
)

// Ticket 75, item 3 (orchestrator's continuation list): the exfil leg of the
// write/sync gate used to be expressed with a hard-coded Windows literal —
// baseOptions injected `C:\Users\test\OneDrive` and TestFourChannelExfilSuite
// wrote `C:\Users\test\OneDrive\Notes\shared.md`. Against a POSIX t.TempDir()
// home that pair could only ever pass for the wrong reason: the separator bug
// made EVERY path sync-suspect, so the assertion was being satisfied by the
// blanket fail-closed net, not by root membership. ed74595 moved the shared
// fixture to a real, platform-shaped directory (that is the leg that runs on
// both platforms, and on POSIX it is the t.TempDir()-built equivalent the
// report asks for); THIS file is the windows tier of the same channel, i.e.
// the half whose subject matter is a path shape that only exists here: '/' and
// '\' being two spellings of one file.
//
// Why the negative half used to live only on Windows: writeGate fails closed
// whenever sync detection is incomplete, and on POSIX detection was incomplete
// while resolveHandle (pathresolver_other.go) was the DEFERRED stub. Ticket 347
// landed realpath there, so the negative half is decidable on both platforms now
// and provenance_syncdirs_other_test.go's branch, parked as a live branch
// instead of a Skip for exactly this moment, evaluates on every Linux run. This
// file stays Windows-only for its actual subject - two separator spellings of one
// file - not as a workaround; the POSIX equivalent is that file's
// backslash-inside-a-name case.

// TestExfilSyncWriteWindowsSpellingInvariant pins both directions of the
// channel on one platform where the verdict is decidable: every separator
// spelling of a real file inside a registry-grade sync root upgrades to
// ChSyncWrite, and no spelling of a real file outside it does. The engine comes
// from m7Engine, which since ticket 82 builds its root and its plain directory
// OUTSIDE the profile so that membership — not the under-profile suspect net,
// and not a SyncDetectionComplete() that POSIX could not report before ticket
// 347 — is what makes the positive and the negative mean "membership" rather
// than "everything is dirty"; the consequence of a confirmed grade (disarming
// that net for an in-profile target) is now pinned on BOTH platforms by
// TestSyncConfirmedGradesDisarmFallback in syncdirs_test.go (moved there by
// ticket 347), and this file keeps only its own separator-shape subject.
func TestExfilSyncWriteWindowsSpellingInvariant(t *testing.T) {
	p, syncTarget, plain := m7Engine(t)

	for _, tc := range []struct {
		name string
		want bool // want == true: must be flagged as ChSyncWrite
		path string
	}{
		{"sync root, native spelling", true, syncTarget},
		{"sync root, all-forward-slash", true, filepath.ToSlash(syncTarget)},
		{"sync root, both separators at once", true, filepath.Dir(syncTarget) + `\Notes/shared.md`},
		{
			"sync root, forward slash under a backslash parent", true,
			filepath.Dir(filepath.Dir(syncTarget)) + `/Notes\shared.md`,
		},
		{"plain dir, native spelling", false, plain},
		{"plain dir, all-forward-slash", false, filepath.ToSlash(plain)},
		{"plain dir, both separators at once", false, filepath.Dir(plain) + `\brand-new.md`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			hit, ok := p.Inspect("task-1", "fs.write", map[string]any{
				"path": tc.path, "content": "see " + marker,
			})
			if !tc.want {
				if ok {
					t.Errorf("false positive: %q is not under the injected sync root but was flagged (%+v)",
						tc.path, hit)
				}
				return
			}
			if !ok {
				t.Fatalf("ESCAPIABLE: writing tainted bytes to %q (inside a confirmed sync root) is not "+
					"an exfil channel", tc.path)
			}
			if hit.Channel != ChSyncWrite {
				t.Errorf("channel: got %q want %q", hit.Channel, ChSyncWrite)
			}
			if !p.IsSyncPath(tc.path).Sync {
				t.Errorf("IsSyncPath(%q) = not sync, although the channel fired: the two verdicts must "+
					"not disagree", tc.path)
			}
		})
	}
}
