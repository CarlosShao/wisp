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
// input - a data root with a junction among its components - handed to the
// private directory API. The promise is "either refuse, or seal the tree the
// filesystem resolved to"; what is not allowed is sealing the tree the caller's
// *string* pointed at, because with a lexical Abs those are different trees.
//
// The load-bearing leg is the one whose target ALREADY EXISTS behind the link:
// that is where the pre-ticket-94 code silently succeeded and narrowed somebody
// else's directory. A missing leaf is refused by the old code too, for an
// accidental reason (Go's Lstat reports a junction as not-a-directory, so the
// climb errored out), and an assertion that passes both before and after the fix
// measures nothing - so it is kept as a second leg, never as the proof.
//
// Both resolver states are measured, because the guarantee must not depend on
// which packages a binary links.
func TestAC3JunctionInputIsRefusedNotSealed(t *testing.T) {
	// Somebody else's tree, widened on purpose. Children keep the ACL they were
	// born with, so a directory created under this parent carries the foreign
	// read grant, which makes "it was sealed" observable as its removal - the
	// only evidence this package accepts (see acl_windows_test.go's header).
	outside := wideParent(t, "someone-elses-tree")
	victim := filepath.Join(outside, "artifacts")
	if err := os.Mkdir(victim, 0o700); err != nil {
		t.Fatal(err)
	}
	innocent := filepath.Join(victim, "keep-me.txt")
	if err := os.WriteFile(innocent, []byte("not ours to seal"), 0o600); err != nil {
		t.Fatal(err)
	}
	victimBefore, _ := aclSIDs(t, victim)
	if !containsSID(victimBefore, everyoneSID) {
		t.Fatalf("the fixture seals nothing to refuse: no inherited foreign grant on %v", victimBefore)
	}
	outsideBefore, _ := aclSIDs(t, outside)

	own := filepath.Join(t.TempDir(), "data")
	if err := winsec.PrivateDirAll(own, 0o700); err != nil {
		t.Fatalf("PrivateDirAll of an ordinary tree: %v", err)
	}
	link := filepath.Join(own, "link")
	mustJunction(t, link, outside)

	t.Run("existing directory behind the link", func(t *testing.T) {
		if winsec.PathResolverInstalled() == nil {
			t.Fatal("C26 is not installed, so this leg measured the fallback instead")
		}
		assertRefusedAndUntouched(t, victim, link, outside, outsideBefore, victimBefore, innocent)
	})

	t.Run("missing directory under the link", func(t *testing.T) {
		if winsec.PathResolverInstalled() == nil {
			t.Fatal("C26 is not installed, so this leg measured the fallback instead")
		}
		neverMade := filepath.Join(link, "spill-artifacts")
		err := winsec.PrivateDirAll(neverMade, 0o700)
		if !errors.Is(err, risk.ErrReparseDenied) {
			t.Fatalf("no C26 refusal: %v", err)
		}
		if _, lErr := os.Lstat(neverMade); !errors.Is(lErr, os.ErrNotExist) {
			t.Errorf("a refusal still created a directory behind the link: %v", lErr)
		}
		if _, sErr := os.Stat(filepath.Join(outside, "spill-artifacts")); sErr == nil {
			t.Errorf("SEELED THE WRONG TREE: %s exists inside the link target", filepath.Join(outside, "spill-artifacts"))
		}
	})

	t.Run("with only the built-in verifier", func(t *testing.T) {
		installed := winsec.PathResolverInstalled()
		winsec.SetPathResolver(nil)
		t.Cleanup(func() { winsec.SetPathResolver(installed) })
		assertRefusedAndUntouched(t, victim, link, outside, outsideBefore, victimBefore, innocent)
		if !errors.Is(winsec.PrivateDirAll(filepath.Join(link, "spill-artifacts"), 0o700), winsec.ErrUnresolvedPath) {
			t.Error("the built-in verifier did not name its own refusal")
		}
	})

	// And the same tree, spelled without the link in it, still gets sealed: the
	// refusal is about traversing a reparse point, not about the platform being
	// unable to make this directory private. If this leg fails, the fix above is
	// refusing the world and the tests that passed are worth nothing.
	if err := winsec.PrivateDirAll(victim, 0o700); err != nil {
		t.Fatalf("PrivateDirAll of the real tree: %v", err)
	}
	assertPrivateACL(t, victim)
}

// assertRefusedAndUntouched is the AC#3 criterion in four measurements: it
// refused, it said which resolver refused, the foreign directory kept every ACE
// it was born with, and its contents are intact.
func assertRefusedAndUntouched(t *testing.T, victim, link, outside string, outsideBefore, victimBefore []string, innocent string) {
	t.Helper()
	unresolved := filepath.Join(link, "artifacts") // the same object, reached through the junction
	err := winsec.PrivateDirAll(unresolved, 0o700)
	if err == nil {
		// Errorf, not Fatalf: the measurements below are the part that shows
		// *what* it sealed, and "it returned success" alone under-reports the bug.
		t.Errorf("PrivateDirAll returned success for a path whose tree lives behind a junction: %s", unresolved)
	} else if !errors.Is(err, winsec.ErrUnresolvedPath) && !errors.Is(err, risk.ErrReparseDenied) {
		t.Errorf("the refusal names neither C26 nor the built-in verifier, so it refused by accident: %v", err)
	}
	after, afterNames := aclSIDs(t, victim)
	if !equalSIDs(after, victimBefore) {
		t.Errorf("WRONG TREE SEALED: %s went from %v to %v (%v)", victim, victimBefore, after, afterNames)
	}
	afterOutside, _ := aclSIDs(t, outside)
	if !equalSIDs(afterOutside, outsideBefore) {
		t.Errorf("the link target's own descriptor changed: %v -> %v", outsideBefore, afterOutside)
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
