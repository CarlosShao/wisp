//go:build !windows

package risk

import (
	"os"
	"path/filepath"
	"testing"
)

// Ticket 75, item 3, POSIX leg — the counterpart of
// TestExfilSyncWriteWindowsSpellingInvariant, built from t.TempDir() exactly as
// the root-cause report demands ("挪进 windows 层，并补一条用 t.TempDir() 构造的
// POSIX 等价用例，否则覆盖面会静默下降"). The subject matter here is a path shape
// that only exists on this kind of system: a backslash that is PART OF A NAME.
//
// What this leg can and cannot prove today, stated plainly:
//   - the POSITIVE half (a write into the injected sync root must upgrade to
//     ChSyncWrite) is meaningful and is asserted unconditionally;
//   - the NEGATIVE half (a plain write must NOT be flagged) is not decidable on
//     POSIX until ticket 55 lands, because sync detection is incomplete here
//     (pathresolver_other.go's resolveHandle is still the DEFERRED stub) and
//     writeGate fails closed on incomplete detection. That half is therefore
//     asserted ONLY when the engine reports detection complete, which is the
//     honest way of parking it: it starts biting the moment ticket 55 lands, and
//     nobody can read today's green as "the POSIX negative control passes".
//     TestWriteGatePlainLocalWriteNotFlagged and its 7 siblings in the same
//     package are that missing-detection red today, and this ticket does not
//     touch them (they belong to ticket 55).

func posixSyncEngine(t *testing.T) (p *Provenance, syncTarget, plainTarget, backslashNameInSyncRoot string) {
	t.Helper()
	base := t.TempDir()
	home := filepath.Join(base, "profile")
	root := filepath.Join(home, "OneDrive")
	work := filepath.Join(home, "work")
	for _, d := range []string{root, work} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	p = NewProvenance(ProvOptions{NoProbe: true, HomeDir: home, SyncRoots: []SyncRoot{
		{Provider: "OneDrive", Path: root, Source: "registry"},
	}})
	p.OpenScope("task-1")
	if !p.Mark("task-1", SrcWebFetch, "https://x", "quote: "+marker) {
		t.Fatal("mark rejected")
	}
	return p, root + "/Notes/out.md", work + "/brand-new.md", root + `/.env\staging`
}

func TestExfilSyncWritePosixSpellingInvariant(t *testing.T) {
	p, syncTarget, plain, backslashName := posixSyncEngine(t)

	// Positive, unconditionally: membership in the injected root, spelled the
	// POSIX way, is the channel. The old fixture could only claim this because
	// every path looked suspect; here the root and the target are real strings
	// that agree with each other.
	for _, tc := range []struct{ name, path string }{
		{"write into sync root", syncTarget},
		// A backslash is not a separator here, so this is ONE file name
		// (`.env\staging`) sitting directly inside the sync root: still a sync
		// write, and the name must not be cut in half on the way to the verdict.
		{"backslash inside a POSIX file name", backslashName},
	} {
		t.Run(tc.name, func(t *testing.T) {
			hit, ok := p.Inspect("task-1", "fs.write", map[string]any{
				"path": tc.path, "content": "see " + marker,
			})
			if !ok {
				t.Fatalf("ESCAPIABLE: %q is inside the injected sync root and was not an exfil channel", tc.path)
			}
			if hit.Channel != ChSyncWrite {
				t.Errorf("channel: got %q want %q", hit.Channel, ChSyncWrite)
			}
		})
	}

	// The reverse direction of the same property: a plain write is not a sync
	// write because of WHERE the bytes land — but see the header comment, this
	// needs complete detection before it can mean anything.
	if !p.SyncDetectionComplete() {
		t.Log("POSIX negative half DORMANT: sync detection is incomplete until ticket 55 " +
			"(realpath + lstat), and writeGate fails closed on incomplete detection; this " +
			"branch is a live assertion the moment ticket 55 lands")
		return
	}
	t.Run("plain write is not flagged", func(t *testing.T) {
		if hit, ok := p.Inspect("task-1", "fs.write", map[string]any{
			"path": plain, "content": "local bytes " + marker,
		}); ok {
			t.Errorf("false positive: plain POSIX write flagged as %q", hit.Channel)
		}
	})
}
