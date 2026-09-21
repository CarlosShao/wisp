package models

// Ticket 109 (source: acceptor-ticket95's AC95-R1/R3/R4).
//
// Ticket 95's ruling - model bytes stay inherit-wide on purpose - stands, and
// the cases here defend it: content is public (hash-pinned in a signed
// manifest) and sealing this class is exactly what would break a second
// instance running under another account. What that ruling did not cover is
// the shape of the clock:
//
//	downloader.go  Ensure -> VerifyDir on the cache hit ... returns the directory
//	bridge.go      Run    -> reports "model available" to the state machine
//	then           the reader opens the files
//
// Verification happens *before* the hand-off, use happens *after* it, and
// nothing in between answers "is this still that file". A principal that can
// write the install directory - measured at SID level in
// handoff_window_109_windows_test.go - owns that whole span, and the swap is
// invisible to everything above this package.
//
// AC#2 closes it with direction (b), re-verify at the hand-off point, because
// direction (a), narrowing that one landing spot, would take away the read
// grant the "do not seal" ruling exists to preserve.
//
// Portable on purpose: no ACL text is parsed here, so the same case runs under
// Linux (see the ticket's Docker reading).

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// swapOneInstalledFile flips one byte of one installed file and returns what
// puts it back. This is the "another account wrote here" step, performed
// in-process because this host has no second account to run as: what the other
// account is *allowed* to do is the part icacls measures.
func swapOneInstalledFile(t *testing.T, dir string, path string) func() {
	t.Helper()
	p := filepath.Join(dir, filepath.FromSlash(path))
	orig, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read installed %s: %v", path, err)
	}
	swapped := append([]byte(nil), orig...)
	swapped[0] ^= 0x01
	if err := os.WriteFile(p, swapped, 0o600); err != nil {
		t.Fatalf("swap %s: %v", path, err)
	}
	t.Logf("swapped 1 byte of %s inside the post-verification window", path)
	return func() {
		if err := os.WriteFile(p, orig, 0o600); err != nil {
			t.Errorf("restore %s: %v", path, err)
		}
	}
}

func newWalkMachine() *statemachine.Machine {
	return statemachine.New(statemachine.Options{Initial: statemachine.StateFirstRun})
}

// TestAC1HandoffRefusesAModelSwappedAfterEnsureIsVerified is ticket 109's
// must-be-red-before-the-fix case.
//
// The swap fires from the Progress hook, which runs inside Ensure *after*
// VerifyDir has passed and while Ensure is announcing completion - so the file
// changes hands (meter passes) verified, and reaches the hand-off point
// tampered. A hand-off that reports the model available is the window.
func TestAC1HandoffRefusesAModelSwappedAfterEnsureIsVerified(t *testing.T) {
	srv := newCountingServer(t, map[string][]byte{
		"/kws-fixture/tiny-model.tar.bz2": mustRead(t, filepath.Join("testdata", "tiny-model.tar.bz2")),
	})
	m := manifestWithURLs(t, srv.URL())
	entry, err := m.FindModel("kws-fixture")
	if err != nil {
		t.Fatal(err)
	}
	store := filepath.Join(t.TempDir(), "models")
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	mgr, err := NewManager(Options{
		DataDir:         store,
		Manifest:        m,
		VerifySignature: true,
		Attempts:        2,
		BackoffBase:     0,
		Progress: func(ev ProgressEvent) {
			if ev.Phase != PhaseDone || ev.ModelID != "kws-fixture" {
				return
			}
			restore := swapOneInstalledFile(t, filepath.Join(store, "kws-fixture"), entry.InstalledFiles()[0].Path)
			t.Cleanup(restore)
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	bridge := WireDownloading(mgr, newWalkMachine())
	if _, err := bridge.Run(context.Background(), "kws-fixture"); err == nil {
		t.Fatalf("AC#1/AC#2 (AC95-R1): the hand-off reported a model available whose file was swapped after Ensure "+
			"verified it - the span between Ensure's return and the reader's open has no guard. state=%v",
			newWalkMachine().State())
	} else {
		t.Logf("hand-off refused the swapped model: %v", err)
	}
}
