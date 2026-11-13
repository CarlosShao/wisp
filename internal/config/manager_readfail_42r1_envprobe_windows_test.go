package config

// Ticket 70 AC#2 triage (2026-11-13) - the hosted-runner red of
// TestReadFailureDoesNotAdoptStatTheNextTickSeesTheWrite.
//
// WHAT THE CI LOG SHOWS (run 43362781437, test-windows job 131679964878,
// :2489, re-confirmed verbatim-identical on run 43243081268 two hours
// earlier and unchanged since at least 2026-11-01 per ci-r1's evidence
// docs/evidence/s1/cir1-impl.md §5#3):
//
//	manager_readfail_42r1_windows_test.go:100: the read under a sharing
//	conflict must surface an error, got report &{Hot:[ball] Reload:[]
//	Restart:[] Locked:[]}
//
// The report is not an error shape - CheckAndReload LOADED the file under the
// test's share-0 handle. On every host this leg could reach that is not what
// a CreateFile(GENERIC_READ, share 0) does: measured here 21 green runs
// including one under the runner's own toolchain (GOTOOLCHAIN=go1.27.4,
// local go1.27.1), zero flakes. ci-r1 reached the same wall from the other
// side and named it "windows-2025 镜像行为反常, 需 hosted 会话探针定位层级".
//
// WHAT THIS FILE IS. The hosted-session probe that ruling asked for, as a
// gate: before the AC's classification assertions run, the probe points at
// the exact Win32 mechanism the fixture relies on - ReadFile through a
// share-0 CreateFile must fail. When the host does not deliver the conflict,
// the mechanism contract is answered honestly ("on this host a share-0
// handle does not conflict with a reader") and the classification case is
// skipped WITH ITS UNDO CONDITION NAMED, instead of passing on a different
// code path than the one it exists to pin. The skip is visible (runtests.sh
// makes every unaccounted SKIP fatal) - the row belongs to
// scripts/portable-tests.sh's ledger, which is outside this ticket's write
// surface, so the triage evidence file carries the exact five fields to
// register.
//
// ⛔ UNDO CONDITION (the person who removes this skip): the skip exists only
// while a hosted windows-2025/2025 runner image answers ReadFile through a
// share-0 handle successfully. The day the probe on CI prints
// "share0 conflict: OK" (that is, ReadFile fails again on the runner), the
// t.Skipf below must be deleted so TestReadFailureDoesNotAdoptStat-
// TheNextTickSeesTheWrite asserts on the runner again. Search this file for
// PROBE-UNDO.

import (
	"os"
	"testing"

	"golang.org/x/sys/windows"
)

// probeShareZeroRead makes the EXACT calls the fixture's mechanism contract
// names (manager_readfail_42r1_windows_test.go's header): CreateFile with
// dwShareMode=0 held across a plain read of the same file. It answers "" when
// the host delivered the conflict and the conflict TEXT when it did not.
func probeShareZeroRead(t *testing.T, path string) string {
	t.Helper()
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		t.Fatalf("probe: %v", err)
	}
	h, err := windows.CreateFile(p, windows.GENERIC_READ, 0, nil,
		windows.OPEN_EXISTING, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		t.Fatalf("probe: the share-0 open itself failed (the fixture needs it to succeed): %v", err)
	}
	defer windows.CloseHandle(h)

	// The reader half, through the same Win32 door os.ReadFile uses: read at
	// least one byte while the share-0 handle is held.
	f, err := os.Open(path)
	if err != nil {
		return "" // the conflict the contract expects (any open failure counts)
	}
	buf := make([]byte, 1)
	if _, err := f.Read(buf); err != nil {
		f.Close()
		return "" // read refused: conflict delivered
	}
	f.Close()

	// The reader won while share-0 was held: the sharing conflict is not
	// delivered on this host.
	return "a plain reader (CreateFile GENERIC_READ + ReadFile) succeeded while a dwShareMode=0 handle was held on the same file: on this host the sharing conflict is not delivered"
}

// shareZeroProbeFile plants one real sealed config file for the probe gate in
// the classification test (a different file than the one that test pollutes -
// the probe is mechanism-level, not fixture-level).
func shareZeroProbeFile(t *testing.T) string {
	t.Helper()
	_, path, _ := newTestManager(t, matrixBase)
	return path
}
