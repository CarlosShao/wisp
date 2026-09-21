//go:build windows

package models

// Ticket 109, the half that has to be measured on NTFS rather than argued:
// who actually holds the write in the span between Ensure's return and the
// reader's open, in the only form this repository accepts as evidence (what
// `icacls` prints, resolved to SIDs).
//
// AC95-R1's claim is about an *inherited* grant on the install directory, so
// the fixture seeds the store root the way a real profile does - BUILTIN\Users
// with Modify, inheritable - and then reads back what the install directory and
// one installed file ended up holding. Two directions are pinned:
//
//	before/after the fix : the ACL is identical, i.e. closing the window did NOT
//	                       narrow that one landing spot (ticket 95's reuse
//	                       rationale survives, AC#2's reverse leg)
//	the write itself     : a swap made through that right is refused at the
//	                       hand-off point (AC#2's forward leg)
//
// And AC#4's nail for the grid ticket 95 left without a pin:
// downloader.go's install() creates <DataDir>/<id> with 0o755 (decorative on
// Windows), so "wide here" is a decision that now has a stated break line.

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// aceLinesNaming returns the ACE lines of path's descriptor that name the given
// account, with the path itself stripped the way the ticket 95 helper does.
func aceLinesNaming(t *testing.T, path, account string) []string {
	t.Helper()
	raw := icaclsRun(t, path)
	var out []string
	for _, line := range strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n") {
		line = strings.ReplaceAll(line, path, "")
		i := strings.LastIndex(line, ":(")
		if i <= 0 {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(line[:i]), account) {
			out = append(out, strings.TrimSpace(line))
		}
	}
	return out
}

func seededStore(t *testing.T) (store string, m *Manifest) {
	t.Helper()
	srv := newCountingServer(t, map[string][]byte{
		"/kws-fixture/tiny-model.tar.bz2": mustRead(t, filepath.Join("testdata", "tiny-model.tar.bz2")),
	})
	m = manifestWithURLs(t, srv.URL())
	store = filepath.Join(t.TempDir(), "models")
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	// "any other local account can write here", spelled the way a profile
	// directory spells it: inheritable Modify for BUILTIN\Users.
	icaclsRun(t, store, "/grant:r", usersSID+`:(OI)(CI)(M)`)
	if got := aceLinesNaming(t, store, `BUILTIN\Users`); len(got) == 0 {
		t.Fatalf("the seed did not land, this test would prove nothing; icacls says: %s", icaclsRun(t, store))
	} else {
		t.Logf("store root seeded with: %v", got)
	}
	return store, m
}

func TestAC1TheWriteWindowIsAnInheritedCrossAccountRight(t *testing.T) {
	store, m := seededStore(t)
	entry, err := m.FindModel("kws-fixture")
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := NewManager(Options{DataDir: store, Manifest: m, VerifySignature: true, Attempts: 2, BackoffBase: 0})
	if err != nil {
		t.Fatal(err)
	}
	dir, err := mgr.Ensure(context.Background(), "kws-fixture")
	if err != nil {
		t.Fatalf("install into a store another local account can write: %v", err)
	}
	target := filepath.Join(dir, filepath.FromSlash(entry.InstalledFiles()[0].Path))

	beforeDir := aceLinesNaming(t, dir, `BUILTIN\Users`)
	beforeFile := aceLinesNaming(t, target, `BUILTIN\Users`)
	t.Logf("BEFORE the hand-off guard - install dir %s: %v", filepath.Base(dir), beforeDir)
	t.Logf("BEFORE the hand-off guard - installed file %s: %v", filepath.Base(target), beforeFile)
	if len(beforeDir) == 0 || len(beforeFile) == 0 {
		t.Fatalf("AC95-R1's premise did not reproduce: no foreign principal on the install dir (%v) or the installed file (%v)", beforeDir, beforeFile)
	}
	for _, l := range append(append([]string{}, beforeDir...), beforeFile...) {
		if !strings.Contains(l, "(I)") {
			t.Errorf("the foreign grant reached this object without the inherited flag, i.e. somebody widened this one object by hand rather than by inheritance - that is ticket 89's family, not this ticket: %q", l)
		}
	}

	// The write, made through the right icacls just printed, at the moment the
	// window is open: after Ensure verified, before the hand-off announced it.
	swapped := false
	var restore func()
	hook := func(ev ProgressEvent) {
		if ev.Phase != PhaseDone || swapped {
			return
		}
		swapped = true
		restore = swapOneInstalledFile(t, dir, entry.InstalledFiles()[0].Path)
	}
	mgr2, err := NewManager(Options{
		DataDir: store, Manifest: m, VerifySignature: true, Attempts: 2, BackoffBase: 0, Progress: hook,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := WireDownloading(mgr2, newWalkMachine()).Run(context.Background(), "kws-fixture"); err == nil {
		t.Errorf("AC#2 forward leg: the hand-off announced a model available although a swap made through that inherited "+
			"Modify right was in flight - installed file %s is not the one the manifest pinned", filepath.Base(target))
	} else {
		t.Logf("hand-off refused the swap made through that right: %v", err)
	}
	if !swapped {
		t.Fatal("the fixture never reached the hand-off point, so the two readings above measure nothing")
	}
	if restore != nil {
		restore()
	}

	afterDir := aceLinesNaming(t, dir, `BUILTIN\Users`)
	afterFile := aceLinesNaming(t, target, `BUILTIN\Users`)
	t.Logf("AFTER the hand-off guard ran - install dir: %v", afterDir)
	t.Logf("AFTER the hand-off guard ran - installed file: %v", afterFile)
	// The reverse leg, at SID level: closing the window must not have touched
	// this landing spot's grants, or the reuse ticket 95 bought is gone.
	if strings.Join(beforeDir, "|") != strings.Join(afterDir, "|") || strings.Join(beforeFile, "|") != strings.Join(afterFile, "|") {
		t.Errorf("the closure narrowed the landing spot after all: dir %v -> %v, file %v -> %v", beforeDir, afterDir, beforeFile, afterFile)
	}
	if err := mgr.VerifyInstalled("kws-fixture"); err != nil {
		t.Errorf("after restore, the clean install no longer verifies: %v", err)
	}
}

// TestAC4InstallDirWidthIsDeliberateAndHasABreakLine nails the grid ticket 95
// ruled on but never pinned: downloader.go's install() calls
// os.MkdirAll(installDir, 0o755) - decorative on Windows - so the install
// directory is as wide as whatever stands above it. "Wide" is acceptable here
// in exactly one form, and this test states where it stops being acceptable:
//
//	OK    : the foreign grant is inherited, i.e. it is the *parent's* policy and
//	        reading is what the second instance needs
//	BREAK : an explicit foreign ACE on this object (a hand-placed grant on the
//	        install directory itself is out-of-band by ticket 89's definition)
//	BREAK : the foreign grant is gone altogether, i.e. somebody sealed this
//	        class and reversed ticket 95's ruling without re-opening it
func TestAC4InstallDirWidthIsDeliberateAndHasABreakLine(t *testing.T) {
	store, m := seededStore(t)
	mgr, err := NewManager(Options{DataDir: store, Manifest: m, VerifySignature: true, Attempts: 2, BackoffBase: 0})
	if err != nil {
		t.Fatal(err)
	}
	dir, err := mgr.Ensure(context.Background(), "kws-fixture")
	if err != nil {
		t.Fatal(err)
	}
	// A store root with NO foreign grant at all: the install dir then inherits
	// nothing wide, which is the "sealed" shape the ruling forbids silently.
	icaclsRun(t, store, "/remove", `BUILTIN\Users`)
	if got := aceLinesNaming(t, dir, `BUILTIN\Users`); len(got) != 0 {
		t.Fatalf("removing the parent's grant left an explicit foreign ACE on the install dir - out-of-band by ticket 89's definition: %v", got)
	}

	// Fresh store, the documented shape.
	store2, m2 := seededStore(t)
	mgr2, err := NewManager(Options{DataDir: store2, Manifest: m2, VerifySignature: true, Attempts: 2, BackoffBase: 0})
	if err != nil {
		t.Fatal(err)
	}
	dir2, err := mgr2.Ensure(context.Background(), "kws-fixture")
	if err != nil {
		t.Fatal(err)
	}
	lines := aceLinesNaming(t, dir2, `BUILTIN\Users`)
	if len(lines) == 0 {
		t.Fatalf("AC#4: the install dir of a store root that grants Users Modify inherits nothing - either this class got sealed "+
			"(re-open ticket 95's ruling) or the width arrived by hand: icacls %s", icaclsRun(t, dir2))
	}
	for _, l := range lines {
		if !strings.Contains(l, "(I)") {
			t.Errorf("AC#4 break line: explicit foreign grant on the install dir %q", l)
		}
	}
	// And reading by another account - the whole point of the width - still
	// works: the installed file is reachable through the directory it is in.
	if _, err := os.ReadFile(filepath.Join(dir2, filepath.FromSlash(tinyArchiveMembers[0].Path))); err != nil {
		t.Errorf("installed file unreadable through the inherit-wide dir, so the width bought nothing: %v", err)
	}
	t.Logf("AC#4 nail: install dir %s inherits %v (parent policy, guard is the re-verify)", filepath.Base(dir2), lines)
}
