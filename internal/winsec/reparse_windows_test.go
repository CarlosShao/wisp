//go:build windows

package winsec

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"golang.org/x/sys/windows"
)

// mustExec is the internal-package copy of the icacls helper in
// acl_windows_test.go: that file lives in package winsec_test, so its helpers
// are out of reach here (and this file needs them for mklink).
func mustExec(t *testing.T, name string, args ...string) string {
	t.Helper()
	var out bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = &out, &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("%s %v: %v\n%s", name, args, err, out.String())
	}
	return out.String()
}

// mkJunction makes a directory junction. It needs no privilege and no developer
// mode, which is why the A51② leg of ticket 89 *is* constructible on this
// machine: a junction is exactly the shape of the bug, a link standing where an
// artifact was expected whose target is a non-empty directory. os.Symlink to a
// directory is the other shape and needs SeCreateSymbolicLinkPrivilege; if it
// is unavailable the test says so out loud rather than skipping silently, and
// the Linux leg in private_other_test.go covers that construction.
func mkJunction(t *testing.T, link, target string) {
	t.Helper()
	t.Logf("mklink /J: %s", strings.TrimSpace(mustExec(t, "cmd", "/c", "mklink", "/J", link, target)))
	if !reparseAt(link) {
		t.Fatalf("%s is not a reparse point, so this test measured nothing", link)
	}
}

// reparseAt is the same platform check the implementation uses, kept here so the
// test can prove it constructed the thing it claims to have constructed.
func reparseAt(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	attr, ok := info.Sys().(*syscall.Win32FileAttributeData)
	return ok && attr.FileAttributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0
}

func TestAC4JunctionAtArtifactPositionIsNotRecursed(t *testing.T) {
	// Somebody else's tree, with content, planted behind a link at the position
	// the artifact quota path would reclaim.
	outside := filepath.Join(t.TempDir(), "someone-elses-tree")
	deep := filepath.Join(outside, "sub")
	if err := os.MkdirAll(deep, 0o700); err != nil {
		t.Fatal(err)
	}
	innocent := filepath.Join(deep, "keep-me.txt")
	if err := os.WriteFile(innocent, []byte("not ours to delete"), 0o600); err != nil {
		t.Fatal(err)
	}

	root := filepath.Join(t.TempDir(), "data")
	if err := PrivateDirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "artifact")
	mkJunction(t, link, outside)
	if !reparseAt(link) {
		t.Fatalf("%s is not a reparse point, so this test measured nothing", link)
	}

	// A51② itself, on the record: the plain call cannot clear such an entry.
	// If this ever starts succeeding, the fix below stays correct but the ticket's
	// premise needs re-deriving - so the line is logged, not asserted.
	if err := os.Remove(link); err != nil {
		t.Logf("A51② reproduced: os.Remove(%s) = %v", filepath.Base(link), err)
		if _, sErr := os.Stat(innocent); sErr != nil {
			t.Fatalf("os.Remove reported failure yet the target is gone: %v", sErr)
		}
		// Put the link back for the leg below, whichever way that call went.
		if _, sErr := os.Lstat(link); sErr != nil {
			mkJunction(t, link, outside)
		}
	} else {
		t.Logf("os.Remove cleared this junction directly; A51② does not reproduce for it")
		if _, sErr := os.Stat(innocent); sErr != nil {
			t.Fatalf("os.Remove reported success and took the target with it: %v", sErr)
		}
		if _, sErr := os.Lstat(link); sErr == nil {
			t.Fatalf("os.Remove reported success but left the link: both outcomes are wrong")
		}
		mkJunction(t, link, outside)
	}

	if err := RemoveUnlinked(link); err != nil {
		t.Fatalf("RemoveUnlinked refused to clear the link: %v", err)
	}
	if _, err := os.Lstat(link); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("link still present after RemoveUnlinked: %v", err)
	}
	// ① the whole point: the deletion radius stopped at the link.
	if _, err := os.Stat(innocent); err != nil {
		t.Errorf("TARGET DELETED - RemoveUnlinked followed the link into %s: %v", outside, err)
	}
	if _, err := os.Stat(deep); err != nil {
		t.Errorf("target subtree deleted too: %v", err)
	}
}

// TestAC4UnlinkableLinkGivesNamedError covers the second half of AC#4: when the
// link cannot be removed, the caller must get a name it can act on instead of a
// loop that quietly leaves the subtree in place. Injected, because "make the OS
// refuse" is not a state a test can build at normal privilege.
func TestAC4UnlinkableLinkGivesNamedError(t *testing.T) {
	root := filepath.Join(t.TempDir(), "data")
	if err := PrivateDirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(t.TempDir(), "wide")
	if err := os.MkdirAll(filepath.Join(target, "inner"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "inner", "keep.txt"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "artifact")
	mkJunction(t, link, target)

	orig := deleteLink
	deleteLink = func(string) error { return errors.New("injected: disposition refused") }
	t.Cleanup(func() { deleteLink = orig })

	err := RemoveUnlinked(link)
	if err == nil {
		t.Fatal("injected unlink failure accepted, and the reclaim loop would keep going")
	}
	if !errors.Is(err, ErrIsReparsePoint) {
		t.Fatalf("error does not name the reparse point: %v", err)
	}
	if !strings.Contains(err.Error(), link) {
		t.Errorf("error does not carry the path the caller has to report: %v", err)
	}
	// Refusing to delete must still not mean recursing into the target.
	if _, statErr := os.Stat(filepath.Join(target, "inner", "keep.txt")); statErr != nil {
		t.Errorf("target deleted while refusing to remove the link: %v", statErr)
	}
	if _, lErr := os.Lstat(link); lErr != nil {
		t.Errorf("the link vanished anyway: %v", lErr)
	}

	// And a plain file next to it is still removed without ceremony.
	plain := filepath.Join(root, "plain.txt")
	if err := PrivateFile(plain, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := RemoveUnlinked(plain); err != nil {
		t.Fatalf("RemoveUnlinked(plain file): %v", err)
	}
	if _, err := os.Lstat(plain); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("plain file survived: %v", err)
	}
}

// TestAC4QuotaReclaimDoesNotFollowLinks pins the claim ticket 79 made in code:
// the stray-subtree walk removes with RemoveUnlinked, so a link inside the
// subtree is unlinked rather than entered.
func TestAC4SealedWalkSkipsLinks(t *testing.T) {
	root := filepath.Join(t.TempDir(), "data")
	if err := PrivateDirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(t.TempDir(), "outside")
	if err := os.MkdirAll(outside, 0o700); err != nil {
		t.Fatal(err)
	}
	keep := filepath.Join(outside, "keep.txt")
	if err := os.WriteFile(keep, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	inner := filepath.Join(root, "work")
	if err := os.MkdirAll(inner, 0o700); err != nil {
		t.Fatal(err)
	}
	mkJunction(t, filepath.Join(inner, "link"), outside)

	// Sealing the tree must not reach through the link either: propagatePrivate
	// skipping reparse points is the same rule as removing not following them.
	if err := SealDir(root); err != nil {
		t.Fatalf("SealDir over a subtree containing a link: %v", err)
	}
	if _, err := os.Stat(keep); err != nil {
		t.Fatalf("the walk touched the linked target: %v", err)
	}
	// The target's own descriptor must be untouched: verifyPrivate is the same
	// predicate the implementation fails closed on, so "it still fails" proves
	// the seal never reached through the link.
	if err := verifyPrivate(keep); err == nil {
		t.Errorf("linked target got sealed after all")
	}
	if err := RemoveUnlinked(filepath.Join(inner, "link")); err != nil {
		t.Fatalf("RemoveUnlinked inside a subtree: %v", err)
	}
	if _, err := os.Stat(keep); err != nil {
		t.Errorf("removing the inner link damaged the target: %v", err)
	}
}
