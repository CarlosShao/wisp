//go:build windows

// Ticket 103 AC#1/AC#2 legs, written before the fix so the red is on record.
//
// This file deliberately lives in the *external* test package: the whole point of
// AC#1 is that the seam can be reached from outside internal/winsec, which is
// exactly what PROBE A of ticket 94's acceptance did. An internal test could not
// demonstrate the attack, and a guard only provable from inside the package is
// not a guard.
package winsec_test

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/CarlosShao/wisp/internal/winsec"
)

// rubberStampResolver is PROBE A's attacker: it approves every spelling it is
// handed, which is the one property a resolver must never have. It rewrites
// nothing, so a pass-through fake is also the *minimum* a fake needs to be
// caught by the floor; the stronger fake (one that happily rewrites through a
// junction) is ticket 102's Actable() leg, not this one.
type rubberStampResolver struct{}

func (rubberStampResolver) Resolve(in string) (string, error) { return in, nil }

// narrowOnlyResolver and its twin both refuse everything, so they are conforming
// by the conformance probe's own rule; they exist to give the single-use guard
// two *different* legitimate resolvers to refuse to swap.
type narrowOnlyResolver struct{}

func (narrowOnlyResolver) Resolve(string) (string, error) {
	return "", errors.New("narrowOnlyResolver: refuses to answer")
}

type narrowOnlyResolverB struct{}

func (narrowOnlyResolverB) Resolve(string) (string, error) {
	return "", errors.New("narrowOnlyResolverB: refuses to answer")
}

func resolverName(r winsec.C26Resolver) string {
	if r == nil {
		return "<floor>"
	}
	return fmt.Sprintf("%T", r)
}

// captureSeamLog swaps the default slog handler for a buffer and hands it back,
// because that is the only audit channel that exists *today*: an assertion on a
// log line a guard has not been written yet is a red test, and a red test is what
// this file is for.
func captureSeamLog(t *testing.T) *strings.Builder {
	t.Helper()
	var buf strings.Builder
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })
	return &buf
}

// makeJunction103 builds a directory junction with mklink /J, which needs no
// privilege and no developer mode. Everything it touches is under t.TempDir().
func makeJunction103(t *testing.T, link, target string) {
	t.Helper()
	cmd := exec.Command("cmd", "/c", "mklink", "/J", link, target)
	out, err := cmd.CombinedOutput()
	t.Logf("mklink /J %s -> %s: %s", filepath.Base(link), filepath.Base(target), strings.TrimSpace(string(out)))
	if err != nil {
		t.Fatalf("mklink /J failed: %v (%s)", err, out)
	}
	info, lErr := os.Lstat(link)
	if lErr != nil {
		t.Fatalf("lstat(%s) after mklink: %v", link, lErr)
	}
	attr, ok := info.Sys().(*syscall.Win32FileAttributeData)
	if !ok || attr.FileAttributes&syscall.FILE_ATTRIBUTE_REPARSE_POINT == 0 {
		t.Fatalf("%s is not a reparse point (mode %v, sys %#v), so this test would measure nothing", link, info.Mode(), info.Sys())
	}
	t.Logf("junction in place: %s (attributes 0x%x)", link, attr.FileAttributes)
}

// foreignVictimTree is somebody else's tree, sealed wide on purpose: an
// Everyone read grant is the marker whose presence or absence the assertions
// read. It returns the victim directory and the SID list it carries.
func foreignVictimTree(t *testing.T) (victim string, innocentFile string, sidsBefore []string) {
	t.Helper()
	parent := wideParent(t, "victim")
	sub := filepath.Join(parent, "sub")
	if err := os.Mkdir(sub, 0o700); err != nil {
		t.Fatal(err)
	}
	run(t, "icacls", sub, "/inheritance:r", "/grant:r",
		"*"+currentSID(t)+":(OI)(CI)(F)",
		"*"+systemSID+":(OI)(CI)(F)",
		"*"+adminsSID+":(OI)(CI)(F)",
		"*"+everyoneSID+":(OI)(CI)(RX)")
	f := filepath.Join(sub, "keep-me.txt")
	if err := os.WriteFile(f, []byte("not ours to seal, not ours to delete"), 0o600); err != nil {
		t.Fatal(err)
	}
	sids, names := aclSIDs(t, parent)
	t.Logf("victim DACL before: sids=%v names=%v", sids, names)
	if !containsSID(sids, everyoneSID) {
		t.Fatalf("fixture is not wide: %s is absent from %v, so this test would measure nothing", everyoneSID, sids)
	}
	return parent, f, sids
}

// ------------------------------------------------------------------ AC#1 ----

// TestAC1SeamRejectsARubberStampAndLeavesTheForeignDaclAlone is the two-legged
// judgement AC#1 demands: installing a resolver that says OK to everything must
// be refused, and "nothing happened" has to be *measured*, not assumed - so the
// victim's own DACL is read before and after at SID level.
func TestAC1SeamRejectsARubberStampAndLeavesTheForeignDaclAlone(t *testing.T) {
	incumbent := winsec.PathResolverInstalled()
	t.Cleanup(func() {
		winsec.SetPathResolver(nil)
		winsec.SetPathResolver(incumbent)
	})
	t.Logf("incumbent resolver: %s", resolverName(incumbent))

	victim, innocentFile, sidsBefore := foreignVictimTree(t)
	dataRoot := filepath.Join(t.TempDir(), "data")
	if err := winsec.PrivateDirAll(dataRoot, 0o700); err != nil {
		t.Fatalf("PrivateDirAll(%s): %v", dataRoot, err)
	}
	link := filepath.Join(dataRoot, "link")
	makeJunction103(t, link, victim)
	// The exact shape of PROBE A: a directory that does not exist yet, reached
	// *through* a junction in the tree we were allowed to write in.
	target := filepath.Join(link, "artifacts")

	// Leg 1: the registration itself must be refused, loudly.
	logged := captureSeamLog(t)
	winsec.SetPathResolver(rubberStampResolver{})
	got := winsec.PathResolverInstalled()
	if resolverName(got) == resolverName(rubberStampResolver{}) {
		t.Errorf("AC#1 leg 1 RED: a pass-through rubber stamp was installed into the sealing seam (now %s)", resolverName(got))
	}
	if !strings.Contains(logged.String(), "rubberStampResolver") {
		t.Errorf("AC#1 leg 1 RED: refusing the fake left no audit record naming it; log was: %s", logged.String())
	}

	// Leg 2: with the fake *not* installed, acting *through* the junction must be
	// refused for the right reason. The reason is asserted because today this
	// call already fails - for an unrelated reason ("link is not a directory"),
	// which is an accidental green and is exactly what acceptance scored as
	// "红在偶然原因" on ticket 94.
	if err := winsec.PrivateDirAll(target, 0o700); err == nil {
		t.Errorf("AC#1 leg 2 RED: PrivateDirAll(%s) through a junction returned nil", target)
	} else if !strings.Contains(strings.ToLower(err.Error()), "reparse") &&
		!errors.Is(err, winsec.ErrUnresolvedPath) {
		t.Errorf("AC#1 leg 2 RED: the refusal is not about the link, so this leg measures nothing: %v", err)
	} else {
		t.Logf("PrivateDirAll through the junction refused as: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(victim, "artifacts")); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("AC#1 leg 2 RED: the foreign tree gained an entry where none should exist: %v", err)
	}

	// Leg 3: SealFile on an *existing* foreign file reached through the junction,
	// i.e. the silent DACL rewrite of PROBE A, judged at SID level.
	through := filepath.Join(link, "sub", "keep-me.txt")
	fileSidsBefore, fileNamesBefore := aclSIDs(t, through)
	t.Logf("foreign file DACL before: sids=%v names=%v", fileSidsBefore, fileNamesBefore)
	if !containsSID(fileSidsBefore, everyoneSID) {
		t.Fatalf("fixture is not wide: %s absent from the foreign file's %v", everyoneSID, fileSidsBefore)
	}
	if err := winsec.SealFile(through); err == nil {
		t.Errorf("AC#1 leg 3 RED: SealFile(%s) through a junction returned nil", through)
	} else if !strings.Contains(strings.ToLower(err.Error()), "reparse") &&
		!errors.Is(err, winsec.ErrUnresolvedPath) {
		t.Errorf("AC#1 leg 3 RED: SealFile refused for an unrelated reason, so this leg measures nothing: %v", err)
	}
	fileSidsAfter, fileNamesAfter := aclSIDs(t, through)
	t.Logf("foreign file DACL after:  sids=%v names=%v", fileSidsAfter, fileNamesAfter)
	if !containsSID(fileSidsAfter, everyoneSID) {
		t.Errorf("AC#1 leg 3 RED: %s was stripped from the foreign file, i.e. its DACL was rewritten: %v", everyoneSID, fileSidsAfter)
	}
	if !equalSIDs(fileSidsAfter, fileSidsBefore) {
		t.Errorf("AC#1 leg 3 RED: foreign file DACL changed: before=%v after=%v", fileSidsBefore, fileSidsAfter)
	}

	// Leg 4: the victim *directory*'s own SID list, read before and after.
	sidsAfter, namesAfter := aclSIDs(t, victim)
	t.Logf("victim DACL after:  sids=%v names=%v", sidsAfter, namesAfter)
	if !containsSID(sidsAfter, everyoneSID) {
		t.Errorf("AC#1 leg 4 RED: %s (Everyone) read grant was stripped from the foreign tree: %v", everyoneSID, sidsAfter)
	}
	if !equalSIDs(sidsAfter, sidsBefore) {
		t.Errorf("AC#1 leg 4 RED: victim DACL changed: before=%v after=%v", sidsBefore, sidsAfter)
	}
	if _, err := os.Stat(innocentFile); err != nil {
		t.Errorf("AC#1 leg 4: the innocent file behind the junction must be untouched: %v", err)
	}
	t.Logf("icacls(victim) after the refused install:\n%s", icaclsRaw(t, victim))
}

// TestAC1SeamIsSingleUse is the other leg of AC#1: two *conforming* resolvers may
// not be swapped under a running seal - the second install is refused and the
// incumbent stays in place.
func TestAC1SeamIsSingleUse(t *testing.T) {
	incumbent := winsec.PathResolverInstalled()
	t.Cleanup(func() {
		winsec.SetPathResolver(nil)
		winsec.SetPathResolver(incumbent)
	})

	logged := captureSeamLog(t)
	winsec.SetPathResolver(nil)
	if got := winsec.PathResolverInstalled(); got != nil {
		t.Fatalf("SetPathResolver(nil) must fall back to the built-in floor, got %s", resolverName(got))
	}
	winsec.SetPathResolver(narrowOnlyResolver{})
	if got := resolverName(winsec.PathResolverInstalled()); got != resolverName(narrowOnlyResolver{}) {
		t.Fatalf("a conforming resolver was not installed onto a free seam (got %s): the guard is refusing everything", got)
	}
	// Now the seam is occupied by a different legitimate resolver: a second
	// install attempt must be refused.
	winsec.SetPathResolver(narrowOnlyResolverB{})
	if got := resolverName(winsec.PathResolverInstalled()); got != resolverName(narrowOnlyResolver{}) {
		t.Errorf("AC#1 RED: the seam was re-installed from %s to %s - it is not single-use", got, resolverName(narrowOnlyResolverB{}))
	}
	if !strings.Contains(logged.String(), "already") {
		t.Errorf("AC#1: the second install must leave an already-installed audit record, got: %s", logged.String())
	}
}

// TestAC1RefusedInstallLeavesTheSealWorking is the leg AC#3② needs to stay green:
// a refused registration is a refusal, not a wedge - after it the installed
// pipeline still seals clean paths and still refuses hostile ones.
func TestAC1RefusedInstallLeavesTheSealWorking(t *testing.T) {
	incumbent := winsec.PathResolverInstalled()
	t.Cleanup(func() {
		winsec.SetPathResolver(nil)
		winsec.SetPathResolver(incumbent)
	})
	victim, _, sidsBefore := foreignVictimTree(t)
	dataRoot := filepath.Join(t.TempDir(), "data")
	if err := winsec.PrivateDirAll(dataRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dataRoot, "link")
	makeJunction103(t, link, victim)

	winsec.SetPathResolver(rubberStampResolver{})

	clean := filepath.Join(dataRoot, "private")
	if err := winsec.PrivateDirAll(clean, 0o700); err != nil {
		t.Fatalf("a clean path must still seal after a refused registration: %v", err)
	}
	assertPrivateACL(t, clean)
	if _, err := os.Stat(clean); err != nil {
		t.Fatal(err)
	}
	if err := winsec.PrivateDirAll(filepath.Join(link, "artifacts"), 0o700); err == nil {
		t.Error("the junction path must still be refused after a refused registration")
	}
	sidsAfter, _ := aclSIDs(t, victim)
	if !equalSIDs(sidsAfter, sidsBefore) || !containsSID(sidsAfter, everyoneSID) {
		t.Errorf("the foreign DACL changed across the refused-install sequence: before=%v after=%v", sidsBefore, sidsAfter)
	}
}

// ------------------------------------------------------------------ AC#2 ----

// TestAC2RemoveUnlinkedRefusesAPathThroughAJunction is PROBE F as a tree case:
// the entry being removed is a plain file, but the spelling reaches it *through*
// a junction, so the object is in somebody else's tree. Today the call returns
// nil and the foreign file is gone.
func TestAC2RemoveUnlinkedRefusesAPathThroughAJunction(t *testing.T) {
	victim, innocentFile, _ := foreignVictimTree(t)
	dataRoot := filepath.Join(t.TempDir(), "data")
	if err := winsec.PrivateDirAll(dataRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dataRoot, "link")
	makeJunction103(t, link, victim)
	through := filepath.Join(link, "sub", "keep-me.txt")

	err := winsec.RemoveUnlinked(through)
	if err == nil {
		t.Errorf("AC#2 RED: RemoveUnlinked(%q) through a junction returned nil", through)
	} else if !errors.Is(err, winsec.ErrIsReparsePoint) {
		t.Errorf("AC#2: the refusal must name ErrIsReparsePoint so a reclaim loop can attribute it, got %v", err)
	}
	if _, statErr := os.Stat(innocentFile); statErr != nil {
		t.Errorf("AC#2 RED: the foreign file behind the junction was deleted (%v)", statErr)
	}
	if _, statErr := os.Stat(filepath.Join(victim, "sub")); statErr != nil {
		t.Errorf("AC#2: the foreign directory was damaged too: %v", statErr)
	}
	if _, lErr := os.Lstat(link); lErr != nil {
		t.Errorf("the junction itself must survive a refused removal: %v", lErr)
	}
}

// TestAC2RemoveUnlinkedStillUnlinksAStandaloneLink is the positive control that
// keeps AC#2 from being scored by a guard that refuses everything: a link
// standing where an artifact was expected is still unlinked, target intact.
func TestAC2RemoveUnlinkedStillUnlinksAStandaloneLink(t *testing.T) {
	victim, innocentFile, _ := foreignVictimTree(t)
	dataRoot := filepath.Join(t.TempDir(), "data")
	if err := winsec.PrivateDirAll(dataRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dataRoot, "stray-link")
	makeJunction103(t, link, victim)

	if err := winsec.RemoveUnlinked(link); err != nil {
		t.Fatalf("RemoveUnlinked of a standalone link must succeed: %v", err)
	}
	if _, err := os.Lstat(link); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("the link survived its own unlink: %v", err)
	}
	if _, err := os.Stat(innocentFile); err != nil {
		t.Errorf("unlinking the link must not touch its target: %v", err)
	}
}
