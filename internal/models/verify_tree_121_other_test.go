//go:build !windows

package models

// Ticket 121 AC#4's link-shaped sibling, POSIX leg.
//
// verifyNothingUnnamed (downloader.go) decides "is this entry a file the
// manifest named" from os.Lstat-equivalent type bits, not from following the
// entry. That distinction is only testable where symlinks are unprivileged and
// always available, which is not this repository's Windows leg: a symlink there
// needs Backup privilege or Developer Mode, so a case that creates one would be
// a green-or-flaky coin flip on a runner rather than an assertion. This is the
// same division of labour tickets 108/113/118/119 use for internal/winsec (their
// "_other_test.go" linking legs), and the Windows half of "who may write here"
// stays measured by icacls in handoff_window_109_windows_test.go.
//
// The hole is specific: verifyFileHash calls os.Stat, which follows a link, so
// before this change an attacker who could write the (inherit-wide, ticket 95)
// install directory could replace model.onnx with a link to a byte-identical
// copy somewhere else, keep every hash green, and move the file the loader
// actually opens outside the tree that was ever checked.

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAC44SymlinkWhereTheManifestNamesAFile is the named-path leg: the bytes the
// link reaches are the manifest's own bytes, so a hash-only check cannot see
// anything wrong, and the tree that was verified is not the tree that is read.
func TestAC44SymlinkWhereTheManifestNamesAFile(t *testing.T) {
	_, m, dir := poisonTree(t)

	real := filepath.Join(t.TempDir(), "genuine-copy")
	body, err := os.ReadFile(filepath.Join(dir, "model.onnx"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(real, body, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(dir, "model.onnx")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(real, filepath.Join(dir, "model.onnx")); err != nil {
		t.Fatal(err)
	}

	// The premise: hashing the named file still passes, because os.Stat follows
	// the link. Without this line the case below could be red for the wrong
	// reason (a missing file) and nobody would know from the log.
	if err := verifyFileHash(filepath.Join(dir, "model.onnx"), tinyArchiveMembers[0].SHA256, tinyArchiveMembers[0].SizeBytes); err != nil {
		t.Fatalf("the premise did not land, the swap is visible to the hash check and this case proves "+
			"nothing about the tree check: %v", err)
	}

	err = m.VerifyInstalled("kws-fixture")
	if err == nil {
		t.Fatalf("AC#4/R-109-2: verification passed on a model directory whose model.onnx is a symlink to "+
			"bytes outside it (%s) - the hashes are true and the file the loader opens is not in the tree "+
			"that was checked", real)
	}
	if !strings.Contains(err.Error(), "model.onnx") || !strings.Contains(err.Error(), "symlink") {
		t.Errorf("the refusal must name the entry and say it is a link, got: %v", err)
	}
	t.Logf("refused and named: %v", err)
}

// TestAC45SymlinkInAnUnnamedPlace is the same check from the other side: the
// manifest's files all stay exactly as installed, and one link is added where
// nothing was named.
func TestAC45SymlinkInAnUnnamedPlace(t *testing.T) {
	_, m, dir := poisonTree(t)
	outside := filepath.Join(t.TempDir(), "outside.bin")
	if err := os.WriteFile(outside, []byte("anything"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, "tokens.conf")); err != nil {
		t.Fatal(err)
	}
	err := m.VerifyInstalled("kws-fixture")
	if err == nil {
		t.Fatal("AC#4/R-109-2: an unnamed symlink in the install directory verified clean")
	}
	if !strings.Contains(err.Error(), "tokens.conf") {
		t.Errorf("the refusal must name tokens.conf, got: %v", err)
	}
	t.Logf("refused and named: %v", err)
}

// TestAC46EnsureStillInstallsACleanTree keeps the fix from becoming a
// self-inflicted refusal: a fresh install through the normal path - archive
// extraction into staging, then the move into <store>/<id> - must still hand
// back a directory this check accepts, on the platform where the extraction is
// exercised at all.
func TestAC46EnsureStillInstallsACleanTree(t *testing.T) {
	_, mgr, dir := poisonTree(t)
	if _, err := os.Stat(filepath.Join(dir, "dict", "inner.txt")); err != nil {
		t.Fatalf("named subdirectory missing after install: %v", err)
	}
	if err := mgr.VerifyInstalled("kws-fixture"); err != nil {
		t.Fatalf("reverse leg on POSIX: a freshly installed tree must verify, got: %v", err)
	}
	state, err := WireDownloading(mgr, newWalkMachine()).Run(context.Background(), "kws-fixture")
	if err != nil {
		t.Fatalf("the hand-off refused a clean install (state=%s): %v", state, err)
	}
}
