//go:build windows

package proc

import (
	"encoding/binary"
	"testing"
)

// 180 orphan-red repair: ticket 66's residual ① said
// "NtQuerySystemInformation: buffer never sufficient（200 次读 1 次）未修".
//
// The mechanism is not flakiness in the kernel, it is the retry rule: a
// STATUS_INFO_LENGTH_MISMATCH answers with the byte count that was needed AT
// THAT INSTANT, and the pre-fix loop took that number verbatim as the next
// buffer size. On a machine whose process/thread population is still moving
// (this repo's own parallel test legs are exactly that kind of load) the next
// call is always one entry-set behind the number it was sized for, so every
// attempt misses and the sixth reports "buffer never sufficient".
//
// These tests drive the real retry loop (systemProcessSnapshot) with a query
// that models a growing system, because the live kernel cannot be asked to grow
// on demand. The live path stays covered by TestSystemProcessSnapshotSeesSelf.

const snapshotGrowthTestPID = uint64(4242)

// writeSnapshotEntryForTest fills buf with ONE decodable
// SYSTEM_PROCESS_INFORMATION entry whose NextEntryOffset is 0 - the chain
// terminator that must still be delivered (A14). Offsets mirror the one decoder
// in systemprocs_windows.go.
func writeSnapshotEntryForTest(buf []byte, pid uint64) {
	for i := range buf {
		buf[i] = 0
	}
	binary.LittleEndian.PutUint32(buf[0:], 0)    // NextEntryOffset: last entry
	binary.LittleEndian.PutUint32(buf[4:], 3)    // NumberOfThreads
	binary.LittleEndian.PutUint64(buf[8:], 4096) // WorkingSetPrivateSize
	binary.LittleEndian.PutUint64(buf[40:], 17)  // UserTime
	binary.LittleEndian.PutUint64(buf[48:], 23)  // KernelTime
	binary.LittleEndian.PutUint64(buf[80:], pid) // UniqueProcessId
	binary.LittleEndian.PutUint32(buf[96:], 11)  // HandleCount
}

// TestSystemProcessSnapshotKeepsUpWithAGrowingSystem is the case that was red
// before the repair: it fails with "buffer never sufficient" when the loop
// sizes itself from the reported length instead of from a grown buffer.
func TestSystemProcessSnapshotKeepsUpWithAGrowingSystem(t *testing.T) {
	const growth = 24 << 10 // bytes of new process/thread entries between two calls

	// Negative control, test-local on purpose: the pre-fix rule, verbatim.
	legacyNeed := uint32(sysProcInitialSize)
	legacySize := uint32(sysProcInitialSize)
	legacyRead := false
	for attempt := 0; attempt < 6; attempt++ {
		legacyNeed += growth
		if legacySize >= legacyNeed {
			legacyRead = true
			break
		}
		legacySize = legacyNeed
	}
	if legacyRead {
		t.Fatalf("negative control converged: this model does not reproduce the pre-fix livelock, so the case proves nothing")
	}

	need := uint32(sysProcInitialSize)
	calls := 0
	snap, err := systemProcessSnapshot(func(buf []byte) (uint32, uint32) {
		calls++
		need += growth
		if uint32(len(buf)) < need {
			return statusInfoLengthMismatch, need
		}
		writeSnapshotEntryForTest(buf, snapshotGrowthTestPID)
		return 0, 0
	})
	if err != nil {
		t.Fatalf("a snapshot whose system is still growing must still be read: %v", err)
	}
	if calls > 3 {
		t.Errorf("snapshot took %d calls, want <= 3 (at-least-doubling must outrun the growth, not chase it)", calls)
	}
	si, ok := snap[uint32(snapshotGrowthTestPID)]
	if !ok {
		t.Fatalf("snapshot pid %d missing from %v", snapshotGrowthTestPID, snap)
	}
	if si.PrivateWorkingSet != 4096 || si.ThreadCount != 3 || si.HandleCount != 11 {
		t.Errorf("decoded entry = %+v, want the bytes written by the model", si)
	}
	if si.UserTime100ns != 17 || si.KernelTime100ns != 23 {
		t.Errorf("decoded CPU times = user %d kernel %d, want 17/23", si.UserTime100ns, si.KernelTime100ns)
	}
}

// TestSystemProcessSnapshotFailsClosedAtTheCap pins the other half: growth is
// bounded, and a snapshot that cannot fit the bound is an error, never an
// unbounded allocation or a spin.
func TestSystemProcessSnapshotFailsClosedAtTheCap(t *testing.T) {
	calls := 0
	_, err := systemProcessSnapshot(func(buf []byte) (uint32, uint32) {
		calls++
		return statusInfoLengthMismatch, sysProcMaxSize + 1
	})
	if err == nil {
		t.Fatal("an over-cap snapshot must fail closed, not keep allocating")
	}
	if calls > 2 {
		t.Errorf("over-cap retry took %d calls, want <= 2 (cap reached, refuse immediately)", calls)
	}
}

// TestSysProcNextSizeNeverShrinksOrExceedsTheCap is the policy's own boundary:
// every accepted step grows, nothing survives past the cap, and a need ahead of
// the doubling keeps headroom instead of landing exactly on the reported edge.
func TestSysProcNextSizeNeverShrinksOrExceedsTheCap(t *testing.T) {
	need := uint32(sysProcInitialSize) + 1
	next, ok := sysProcNextSize(sysProcInitialSize, need)
	if !ok {
		t.Fatal("a growable step was refused")
	}
	if next < sysProcInitialSize*2 {
		t.Errorf("next size = %d, want >= one doubling (%d)", next, sysProcInitialSize*2)
	}
	// The headroom is the whole fix: the reported need alone would be exactly
	// the size one call behind.
	big := uint32(8 << 20)
	nextBig, ok := sysProcNextSize(sysProcInitialSize, big)
	if !ok {
		t.Fatal("a growable step to the reported need was refused")
	}
	if nextBig <= big {
		t.Errorf("next size = %d, want > reported need %d (headroom for growth between two calls)", nextBig, big)
	}
	if nextBig > sysProcMaxSize {
		t.Errorf("next size = %d, want <= cap %d", nextBig, sysProcMaxSize)
	}
	if got, ok := sysProcNextSize(sysProcMaxSize, sysProcMaxSize); ok || got != sysProcMaxSize {
		t.Errorf("at the cap: (%d, %v), want (%d, false)", got, ok, sysProcMaxSize)
	}
	if got, ok := sysProcNextSize(sysProcMaxSize-1, sysProcMaxSize+1); !ok || got != sysProcMaxSize {
		t.Errorf("one step below the cap: (%d, %v), want (%d, true)", got, ok, sysProcMaxSize)
	}
}
