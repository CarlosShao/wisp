//go:build windows

package winsec_test

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/winsec"
	"golang.org/x/sys/windows"
)

// TestC26PipelineIsWiredIntoWinsec pins the wiring itself, which is the load
// bearing half of ticket 94: internal/winsec cannot import internal/risk
// (risk -> observe -> secret -> winsec is already a path, so the reverse edge is
// a cycle), so the pipeline reaches the sealing code from the other direction,
// from internal/risk/winsec_c26.go's init. If that file is deleted, or if the
// install is moved behind something that can be skipped, every seal below falls
// back to the built-in verifier and this test says so out loud.
func TestC26PipelineIsWiredIntoWinsec(t *testing.T) {
	r := winsec.PathResolverInstalled()
	if r == nil {
		t.Fatal("no C26 pipeline installed in winsec: internal/risk/winsec_c26.go's init did not run")
	}
	// The installed resolver must be C26 and not some stand-in: the proof is the
	// one behavior only the real pipeline has, denial of a non-excepted junction.
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.Mkdir(target, 0o700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "link")
	mustJunction(t, link, target)
	if _, err := r.Resolve(link); !errors.Is(err, risk.ErrReparseDenied) {
		t.Fatalf("installed resolver accepted a junction traversal: %v", err)
	}
}

// TestAC3JunctionInputIsRefusedNotSealed is ticket 94 AC#3: a real unresolved
// input - a data root whose component IS a junction - handed to the private
// directory API. The promise is "either refuse, or seal the tree the filesystem
// resolved to"; what is not allowed is sealing the tree the caller's *string*
// pointed at, which with filepath.Abs is the tree behind the link.
//
// Both resolver states are measured, because the guarantee must not depend on
// which packages a binary links.
func TestAC3JunctionInputIsRefusedNotSealed(t *testing.T) {
	// Somebody else's tree, widened on purpose: after a refused seal the foreign
	// read grant has to still be there, because "sealed" and "untouched" are only
	// distinguishable by their difference. icacls text is the evidence this
	// repository accepts (see acl_windows_test.go's header).
	outside := wideParent(t, "someone-elses-tree")
	innocent := filepath.Join(outside, "keep-me.txt")
	if err := os.WriteFile(innocent, []byte("not ours to seal"), 0o600); err != nil {
		t.Fatal(err)
	}
	before, _ := aclSIDs(t, outside)
	if !containsSID(before, everyoneSID) {
		t.Fatalf("the fixture seals nothing to refuse: no foreign grant in %v", before)
	}

	own := filepath.Join(t.TempDir(), "data")
	if err := winsec.PrivateDirAll(own, 0o700); err != nil {
		t.Fatalf("PrivateDirAll of an ordinary tree: %v", err)
	}
	link := filepath.Join(own, "link")
	mustJunction(t, link, outside)

	// The input is what a caller actually builds: join the root onto the name the
	// caller believes it owns, then go one level deeper that does not exist yet.
	unresolved := filepath.Join(link, "artifacts")

	t.Run("with C26 installed", func(t *testing.T) {
		if winsec.PathResolverInstalled() == nil {
			t.Fatal("C26 is not installed, so this leg measured the fallback instead")
		}
		assertRefusedAndUntouched(t, unresolved, outside, before, innocent)
	})

	t.Run("with only the built-in verifier", func(t *testing.T) {
		installed := winsec.PathResolverInstalled()
		winsec.SetPathResolver(nil)
		t.Cleanup(func() { winsec.SetPathResolver(installed) })
		err := winsec.PrivateDirAll(unresolved, 0o700)
		if !errors.Is(err, winsec.ErrUnresolvedPath) {
			t.Fatalf("the built-in verifier did not name its own refusal: %v", err)
		}
		assertRefusedAndUntouched(t, unresolved, outside, before, innocent)
	})

	// And the same tree, spelled without the link in it, still works: the refusal
	// is about the traversal, not about the platform being unable to seal here.
	if err := winsec.PrivateDirAll(filepath.Join(own, "artifacts-real"), 0o700); err != nil {
		t.Fatalf("PrivateDirAll of the real tree: %v", err)
	}
	assertPrivateACL(t, filepath.Join(own, "artifacts-real"))
}

// assertRefusedAndUntouched is the whole AC#3 criterion in three measurements: an
// error, nothing created behind the link, and the foreign tree still wide.
func assertRefusedAndUntouched(t *testing.T, unresolved string, outside string, before []string, innocent string) {
	t.Helper()
	err := winsec.PrivateDirAll(unresolved, 0o700)
	if err == nil {
		t.Fatalf("PrivateDirAll sealed through a junction: %s", unresolved)
	}
	if _, lErr := os.Lstat(unresolved); !errors.Is(lErr, os.ErrNotExist) {
		t.Errorf("a refusal still created the directory behind the link: %v", lErr)
	}
	if _, sErr := os.Stat(filepath.Join(outside, "artifacts")); sErr == nil {
		t.Errorf("SEELED THE WRONG TREE: %s\\artifacts exists inside the link target", filepath.Base(outside))
	}
	after, afterNames := aclSIDs(t, outside)
	if !equalSIDs(after, before) {
		t.Errorf("the foreign tree's descriptor changed on a call that refused it:\nbefore %v\nafter  %v (%v)",
			before, after, afterNames)
	}
	if _, sErr := os.Stat(innocent); sErr != nil {
		t.Errorf("target contents damaged: %v", sErr)
	}
}

func equalSIDs(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func containsSID(sids []string, want string) bool {
	for _, s := range sids {
		if strings.EqualFold(s, want) {
			return true
		}
	}
	return false
}

// mustJunction is the external-package copy of mkJunction in
// reparse_windows_test.go (that file is package winsec, so its helpers are out
// of reach), and it fails the same way: a junction that did not materialize
// means the test measured nothing, which is not a pass.
func mustJunction(t *testing.T, link, target string) {
	t.Helper()
	var out bytes.Buffer
	cmd := exec.Command("cmd", "/c", "mklink", "/J", link, target)
	cmd.Stdout, cmd.Stderr = &out, &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("mklink /J %s %s: %v\n%s", link, target, err, out.String())
	}
	info, err := os.Lstat(link)
	if err != nil {
		t.Fatalf("junction not created: %v", err)
	}
	attr, ok := info.Sys().(*syscall.Win32FileAttributeData)
	if !ok || attr.FileAttributes&windows.FILE_ATTRIBUTE_REPARSE_POINT == 0 {
		t.Fatalf("%s is not a reparse point, so this test measured nothing", link)
	}
	t.Logf("junction built: %s -> %s", link, target)
}
