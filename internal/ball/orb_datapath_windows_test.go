//go:build windows

package ball

// Teeth for the per-process orb user-data folder (ticket 408): the shared
// folder made `go test ./cmd/wisp` self-collide on one WebView2 UDF until the
// whole package timed out. Four teeth live here:
//
//   - stability within one process (the same attach reads the path twice:
//     browser.DataPath and the site folder must land in the SAME folder),
//   - difference across processes, pinned to the pid fact (the child's own
//     pid must appear in the child's path) so dropping the pid segment turns
//     this tooth red,
//   - the LOCALAPPDATA root and its empty-variable fallback to the temp dir,
//   - the stale-entry sweep: a dead owner's dir+lock is removed, a live
//     owner's dir is spared, and a lockless folder is never touched.
//
// The mutation proof (remove the pid segment from orbDataDirFor => the
// cross-process tooth goes red) ran on 2026-11-08, see
// docs/evidence/s1/cmdwisp-udf-per-process.md.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/sys/windows"
)

// TestOrbDataPathIsStablePerProcessAndCarriesThePid is tooth 1: two calls in
// THIS process agree, the leaf is the per-process `<pid>-<8 hex>` shape, and
// the leaf is not the old shared folder itself.
func TestOrbDataPathIsStablePerProcessAndCarriesThePid(t *testing.T) {
	first, second := orbDataPath(), orbDataPath()
	if first != second {
		t.Fatalf("orbDataPath moved inside one process: %q then %q - the attach reads it twice (browser.DataPath and the site folder) and they must be the same folder", first, second)
	}
	if filepath.Base(first) == "orb" {
		t.Fatalf("orbDataPath returned the shared root %q itself: the per-process leaf is gone", first)
	}
	wantLeaf := strconv.Itoa(os.Getpid()) + "-"
	if !strings.Contains(filepath.Base(first), wantLeaf) {
		t.Fatalf("orbDataPath leaf %q does not carry this process's pid (%s): pid-level uniqueness is gone", filepath.Base(first), os.Getpid())
	}
	if !orbDataLeafRe.MatchString(filepath.Base(first)) {
		t.Fatalf("orbDataPath leaf %q does not match the per-process leaf shape <pid>-<8hex>", filepath.Base(first))
	}
}

// TestOrbDataPathDiffersAcrossProcesses is tooth 2: a freshly exec'd copy of
// this very test binary returns a DIFFERENT folder, and that folder carries
// the child's OWN pid (read from the exec'd process, not guessed) - so the
// pid segment is load-bearing, not decoration.
func TestOrbDataPathDiffersAcrossProcesses(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestOrbDataPathHelperProcess$", "-test.count=1")
	cmd.Env = append(os.Environ(), orbDataPathHelperEnv+"=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("helper test binary failed: %v\n%s", err, out)
	}
	childPath := orbDataPathHelperParse(string(out))
	if childPath == "" {
		t.Fatalf("no %s line in helper output:\n%s", orbDataPathHelperEnv, out)
	}
	if childPath == orbDataPath() {
		t.Fatalf("child process returned THIS process's folder %q: per-process uniqueness is gone", childPath)
	}
	if childPid := strconv.Itoa(cmd.Process.Pid); !strings.Contains(filepath.Base(childPath), childPid+"-") {
		t.Fatalf("child's path leaf %q does not carry the child's own pid %s: the pid segment is removable, tooth would stay green over a shared folder", filepath.Base(childPath), childPid)
	}
}

// TestOrbDataPathHelperProcess is tooth 2's re-exec arm: when run normally it
// is a no-op; with the helper env set it prints this process's orbDataPath.
func TestOrbDataPathHelperProcess(t *testing.T) {
	if os.Getenv(orbDataPathHelperEnv) != "1" {
		return
	}
	fmt.Printf("%s=%s\n", orbDataPathHelperEnv, orbDataPath())
}

const orbDataPathHelperEnv = "WISP_ORB_DATAPATH_HELPER"

func orbDataPathHelperParse(out string) string {
	for _, line := range strings.Split(out, "\n") {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(line), orbDataPathHelperEnv+"="); ok {
			return rest
		}
	}
	return ""
}

// TestOrbDataDirForSitsUnderBaseAndFallsBackToTemp is tooth 3: the folder
// sits INSIDE the given base, and an empty LOCALAPPDATA root lands under
// os.TempDir() (the old fallback, pinned here so it cannot quietly vanish).
func TestOrbDataDirForSitsUnderBaseAndFallsBackToTemp(t *testing.T) {
	base := t.TempDir()
	dir := orbDataDirFor(base)
	if !strings.HasPrefix(dir, base+string(filepath.Separator)) {
		t.Fatalf("orb dir %q is not under its base %q", dir, base)
	}
	if rel, err := filepath.Rel(base, dir); err != nil || strings.HasPrefix(rel, "..") {
		t.Fatalf("orb dir %q does not sit inside base %q (rel=%q err=%v)", dir, base, rel, err)
	}
	if leaf := filepath.Base(dir); !orbDataLeafRe.MatchString(leaf) {
		t.Fatalf("leaf %q does not match the per-process leaf shape <pid>-<8hex>", leaf)
	}

	t.Setenv("LOCALAPPDATA", "")
	if got := orbDataBase(); got != os.TempDir() {
		t.Fatalf("empty LOCALAPPDATA must fall back to the temp dir, got %q (temp %q)", got, os.TempDir())
	}
	if dir := orbDataDirFor(orbDataBase()); !strings.HasPrefix(dir, os.TempDir()+string(filepath.Separator)) {
		t.Fatalf("fallback dir %q is not under the temp dir %q", dir, os.TempDir())
	}
}

// TestOrbSweepRemovesDeadOwnerAndSparesLiveOwner pins the cleanup promise:
// an entry whose lock opens exclusively (dead owner) goes away with its lock;
// an entry whose lock is held (live owner) survives with its lock; and a
// folder with NO lock is never touched (conservative branch).
func TestOrbSweepRemovesDeadOwnerAndSparesLiveOwner(t *testing.T) {
	root := t.TempDir()

	deadDir := filepath.Join(root, "12345-deadbeef")
	if err := os.MkdirAll(deadDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "12345-deadbeef.lock"), nil, 0o600); err != nil {
		t.Fatal(err)
	}

	liveDir := filepath.Join(root, "67890-cafef00d")
	if err := os.MkdirAll(liveDir, 0o755); err != nil {
		t.Fatal(err)
	}
	liveLock, err := windows.CreateFile(filepath.Join(root, "67890-cafef00d.lock"),
		windows.GENERIC_READ|windows.GENERIC_WRITE, 0, nil,
		windows.OPEN_ALWAYS, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		t.Fatalf("hold the live owner's lock: %v", err)
	}
	defer windows.CloseHandle(liveLock)

	locklessDir := filepath.Join(root, "11111-00000000")
	if err := os.MkdirAll(locklessDir, 0o755); err != nil {
		t.Fatal(err)
	}

	foreignDir := filepath.Join(root, "EBWebView")
	if err := os.MkdirAll(foreignDir, 0o755); err != nil {
		t.Fatal(err)
	}

	orbSweepStaleDirs(root)

	if _, err := os.Stat(deadDir); !os.IsNotExist(err) {
		t.Errorf("dead owner's dir survived the sweep: %q", deadDir)
	}
	if _, err := os.Stat(filepath.Join(root, "12345-deadbeef.lock")); !os.IsNotExist(err) {
		t.Errorf("dead owner's lock survived the sweep")
	}
	if _, err := os.Stat(liveDir); err != nil {
		t.Errorf("live owner's dir was swept: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "67890-cafef00d.lock")); err != nil {
		t.Errorf("live owner's lock was swept: %v", err)
	}
	if _, err := os.Stat(locklessDir); err != nil {
		t.Errorf("lockless dir must be skipped, was swept: %v", err)
	}
	if _, err := os.Stat(foreignDir); err != nil {
		t.Errorf("foreign (non per-process) entry must never be touched, was swept: %v", err)
	}
}
