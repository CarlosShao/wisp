//go:build windows

package winsec

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// everyoneToken is how SDDL spells S-1-1-0; icacls and the SID form are
// accepted too, because which one appears depends on name resolution, not on
// whether the grant exists.
func namesEveryone(s string) bool {
	return strings.Contains(s, "WD") || strings.Contains(s, "S-1-1-0") || strings.Contains(s, "Everyone")
}

// TestSealReportsThePrincipalsItCleared pins the semantics that ticket 89's
// acceptance found missing from the suite: winsec owns the grants on the tree it
// seals, and a grant somebody placed out of band is removed **loudly**.
//
// The loud half is the point. The alternative readings of this design choice are
// (a) allow foreign principals - which would un-do the PROTECTED leg the same
// acceptance report proved load-bearing - and (b) say nothing, which is what the
// code did: an operator's service-account grant vanished at the next memory.Open
// with no error and no record. So the whitelist stays closed and every removal
// of an *explicit* principal becomes a notice.
func TestSealReportsThePrincipalsItCleared(t *testing.T) {
	var got []narrowNotice
	orig := noticeNarrowed
	noticeNarrowed = func(n narrowNotice) { got = append(got, n) }
	t.Cleanup(func() { noticeNarrowed = orig })

	root := filepath.Join(t.TempDir(), "data")
	if err := PrivateDirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}

	// Leg 1 (the "nothing to report" control): re-sealing a tree that is already
	// private must stay quiet, or the notice is noise nobody reads.
	got = nil
	if err := SealDir(root); err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("a seal that removed nothing reported %d notice(s): %+v", len(got), got)
	}

	// Leg 2: an explicit grant on the root itself, and an explicit grant on one
	// child. Both are operator-shaped: `icacls /grant` writes into that object's
	// own DACL.
	mustExec(t, "icacls", root, "/grant", "*S-1-1-0:(OI)(CI)(RX)")
	child := filepath.Join(root, "shared-with-a-service-account.txt")
	if err := os.WriteFile(child, []byte("artifact"), 0o600); err != nil {
		t.Fatal(err)
	}
	mustExec(t, "icacls", child, "/grant", "*S-1-1-0:(RX)")
	// A third child that only *inherits* the root's foreign grant. The OS
	// recomputes it once the root is narrowed, so this is the noise case: the
	// notice is about grants removed from an object's own DACL.
	inherited := filepath.Join(root, "inherits-the-root.txt")
	if err := os.WriteFile(inherited, []byte("sidecar"), 0o600); err != nil {
		t.Fatal(err)
	}

	got = nil
	if err := SealDir(root); err != nil {
		t.Fatalf("SealDir over an explicitly widened tree: %v", err)
	}

	reported := map[string][]string{}
	for _, n := range got {
		reported[n.Path] = append(reported[n.Path], strings.Join(n.Principals, "|"))
	}
	for _, p := range []string{root, child} {
		var hit bool
		for path, principals := range reported {
			if noticeNamesTree(narrowNotice{Path: path}, p) {
				for _, s := range principals {
					hit = hit || namesEveryone(s)
				}
			}
		}
		if !hit {
			t.Errorf("seal cleared a grant on %s without reporting it; notices: %+v", filepath.Base(p), got)
		}
	}
	for _, n := range got {
		if noticeNamesTree(n, inherited) {
			t.Errorf("%s held only an *inherited* copy of the foreign grant and the OS recomputed it, "+
				"so reporting it is the noise this notice was supposed to avoid: %+v", filepath.Base(inherited), n)
		}
	}

	// And the removal really happened - the notice is not describing a no-op.
	for _, p := range []string{root, child} {
		if err := verifyPrivate(p); err != nil {
			t.Errorf("%s still is not private after the reported narrowing: %v", filepath.Base(p), err)
		}
	}
}

// TestSealNoticeIsRecordedByDefault is the "不能静默" half at the default seam:
// with no test hook installed, a narrowing goes into the log pipeline the
// application already captures, carrying the path and the principal.
func TestSealNoticeIsRecordedByDefault(t *testing.T) {
	var buf bytes.Buffer
	origDefault := slog.Default()
	orig := noticeNarrowed
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() {
		noticeNarrowed = orig
		slog.SetDefault(origDefault)
	})

	root := filepath.Join(t.TempDir(), "data")
	if err := PrivateDirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	mustExec(t, "icacls", root, "/grant", "*S-1-5-6:(OI)(CI)(RX)") // SERVICE, outside the whitelist
	if err := SealDir(root); err != nil {
		t.Fatalf("SealDir: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "level=WARN") || !strings.Contains(out, "cleared=") {
		t.Fatalf("the default notifier wrote nothing auditable: %q", out)
	}
	if !strings.Contains(out, "SU") {
		t.Errorf("the log line does not name the principal it cleared: %q", out)
	}
	if !strings.Contains(out, filepath.Base(root)) {
		t.Errorf("the log line does not carry the path an operator has to grep for: %q", out)
	}
	t.Logf("recorded notice: %s", strings.TrimSpace(out))
}
