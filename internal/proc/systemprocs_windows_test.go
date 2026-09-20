//go:build windows

package proc

import (
	"encoding/binary"
	"os"
	"testing"
)

// Regression tests for registry A14 (ticket 66): the last entry of a
// SYSTEM_PROCESS_INFORMATION chain carries NextEntryOffset == 0, and the
// pre-fix loop tested that terminator BEFORE parsing the current entry, so
// the snapshot's tail process could never enter the map. Freshly created
// processes live at the tail of ActiveProcessLinks, which is exactly the
// process that samples itself -> `wisp slo -state X` failed closed with
// "system snapshot does not contain the sampling process".
//
// These tests construct the buffer by hand, so they do not depend on which
// process happens to be last in a live snapshot.

// snapshotStride is the per-entry pitch used by the fixtures below; real
// entries are larger (the image name follows the fixed header), and any
// pitch >= systemProcEntryFloor exercises the same code path.
const snapshotStride = 0xC0

type fixtureEntry struct {
	pid         uint64
	threads     uint32
	handles     uint32
	privateWS   int64
	user100ns   int64
	kernel100ns int64
	next        uint32
}

// buildSnapshot lays out es at stride intervals, last entry's terminator
// included, and returns a buffer long enough that a decoder running past
// the terminator would read real (non-zero) offsets instead of falling off
// the end.
func buildSnapshot(t *testing.T, stride int, es []fixtureEntry) []byte {
	t.Helper()
	if stride < systemProcEntryFloor {
		t.Fatalf("fixture stride %d is below the decoder floor %d", stride, systemProcEntryFloor)
	}
	buf := make([]byte, stride*len(es)+stride) // + one entry of tail slack
	for i, e := range es {
		off := i * stride
		binary.LittleEndian.PutUint32(buf[off+0:], e.next)
		binary.LittleEndian.PutUint32(buf[off+4:], e.threads)
		binary.LittleEndian.PutUint64(buf[off+8:], uint64(e.privateWS))
		binary.LittleEndian.PutUint64(buf[off+40:], uint64(e.user100ns))
		binary.LittleEndian.PutUint64(buf[off+48:], uint64(e.kernel100ns))
		binary.LittleEndian.PutUint64(buf[off+80:], e.pid)
		binary.LittleEndian.PutUint32(buf[off+96:], e.handles)
	}
	return buf
}

func entryFor(i int, pid uint64, last bool) fixtureEntry {
	e := fixtureEntry{
		pid:         pid,
		threads:     uint32(10 + i),
		handles:     uint32(200 + i),
		privateWS:   int64((3 << 20) + i*4096),
		user100ns:   int64(1000 + i),
		kernel100ns: int64(2000 + i),
	}
	if !last {
		e.next = snapshotStride
	}
	return e
}

// TestParseSystemProcessesKeepsLastSnapshotEntry is the A14 regression: the
// pid that is LAST in the chain must be in the map with its fields parsed.
// Asserting the fields (not just len(out)) is deliberate - a fix that only
// made the entry count right while truncating its data must not pass.
func TestParseSystemProcessesKeepsLastSnapshotEntry(t *testing.T) {
	const lastPID = 0x1F90 // 8080: a plausible tail-of-chain pid
	pids := []uint64{4, 1337, 4004, 60060, lastPID}
	es := make([]fixtureEntry, 0, len(pids))
	for i, pid := range pids {
		es = append(es, entryFor(i, pid, i == len(pids)-1))
	}
	if es[len(es)-1].next != 0 {
		t.Fatal("fixture must terminate the chain on the last entry")
	}

	out, err := parseSystemProcesses(buildSnapshot(t, snapshotStride, es))
	if err != nil {
		t.Fatalf("parseSystemProcesses: %v", err)
	}
	if len(out) != len(pids) {
		t.Fatalf("parsed %d of %d snapshot entries: %v", len(out), len(pids), keysOf(out))
	}
	last, ok := out[lastPID]
	if !ok {
		t.Fatalf("last snapshot entry pid %d missing from the map (A14 regression): %v", lastPID, keysOf(out))
	}
	want := es[len(es)-1]
	if last.PrivateWorkingSet != want.privateWS {
		t.Errorf("last entry private working set = %d, want %d", last.PrivateWorkingSet, want.privateWS)
	}
	if last.HandleCount != want.handles {
		t.Errorf("last entry handle count = %d, want %d", last.HandleCount, want.handles)
	}
	if last.ThreadCount != want.threads {
		t.Errorf("last entry thread count = %d, want %d", last.ThreadCount, want.threads)
	}
	if last.UserTime100ns != want.user100ns || last.KernelTime100ns != want.kernel100ns {
		t.Errorf("last entry times = user %d kernel %d, want user %d kernel %d",
			last.UserTime100ns, last.KernelTime100ns, want.user100ns, want.kernel100ns)
	}
	// Every earlier entry must still be parsed identically (the fix must not
	// shift fields between entries).
	for i, pid := range pids[:len(pids)-1] {
		got, ok := out[uint32(pid)]
		if !ok {
			t.Fatalf("entry %d (pid %d) missing: %v", i, pid, keysOf(out))
		}
		if got.PrivateWorkingSet != es[i].privateWS || got.HandleCount != es[i].handles ||
			got.ThreadCount != es[i].threads || got.UserTime100ns != es[i].user100ns ||
			got.KernelTime100ns != es[i].kernel100ns {
			t.Errorf("entry %d (pid %d) decoded as %+v, want %+v", i, pid, got, es[i])
		}
	}
}

// TestWalkSystemProcessesStopsAtTerminator is the other half of the
// ordering rule: parsing the terminator entry must not turn into parsing
// the whole buffer. The fixture carries a live entry AFTER the terminator
// that must never be delivered.
func TestWalkSystemProcessesStopsAtTerminator(t *testing.T) {
	stride := snapshotStride
	es := []fixtureEntry{
		entryFor(0, 100, false),
		entryFor(1, 200, true), // NextEntryOffset == 0: end of chain
		entryFor(2, 300, true), // unreachable: sits past the terminator
	}
	buf := buildSnapshot(t, stride, es)

	var seen []uint64
	if err := WalkSystemProcesses(buf, func(e SystemProcEntry) {
		seen = append(seen, e.PID)
	}); err != nil {
		t.Fatalf("WalkSystemProcesses: %v", err)
	}
	if len(seen) != 2 || seen[0] != 100 || seen[1] != 200 {
		t.Fatalf("walker delivered %v, want [100 200] (stop after the terminator entry)", seen)
	}
}

// TestParseSystemProcessesSkipsPIDZero keeps the "System Idle Process"
// entry out of the map while still counting it as chain-terminating when
// it is last.
func TestParseSystemProcessesSkipsPIDZero(t *testing.T) {
	es := []fixtureEntry{
		entryFor(0, 0, false),   // System Idle Process
		entryFor(1, 4242, true), // last entry, carries the terminator
	}
	out, err := parseSystemProcesses(buildSnapshot(t, snapshotStride, es))
	if err != nil {
		t.Fatalf("parseSystemProcesses: %v", err)
	}
	if _, ok := out[0]; ok {
		t.Errorf("pid 0 must not be addressable in the snapshot map")
	}
	if _, ok := out[4242]; !ok {
		t.Errorf("last entry pid 4242 missing (A14 regression): %v", keysOf(out))
	}
}

// TestParseSystemProcessesRejectsTruncatedEntry proves the per-entry bounds
// check survives the reordering: an entry header that runs off the end of
// the buffer is an error, not a silent short read.
func TestParseSystemProcessesRejectsTruncatedEntry(t *testing.T) {
	buf := make([]byte, snapshotStride+systemProcEntryFloor-1)
	binary.LittleEndian.PutUint32(buf[0:], snapshotStride) // entry 0 -> entry 1
	binary.LittleEndian.PutUint64(buf[80:], 1234)          // entry 0 pid
	// Entry 1 starts but its header does not fit.
	binary.LittleEndian.PutUint32(buf[snapshotStride:], snapshotStride)
	if _, err := parseSystemProcesses(buf); err == nil {
		t.Fatal("parseSystemProcesses accepted a truncated entry header")
	}
}

// TestParseSystemProcessesRejectsNonAdvancingEntry pins the infinite-loop
// guard: NextEntryOffset < 8 cannot move the walker forward.
func TestParseSystemProcessesRejectsNonAdvancingEntry(t *testing.T) {
	buf := make([]byte, snapshotStride*2)
	binary.LittleEndian.PutUint32(buf[0:], 4) // would rewind / stall
	binary.LittleEndian.PutUint64(buf[80:], 7)
	if _, err := parseSystemProcesses(buf); err == nil {
		t.Fatal("parseSystemProcesses accepted a non-advancing entry")
	}
}

// TestSystemProcessSnapshotSeesSelf is the live-kernel half of A14: whatever
// the kernel returns, the sampling process must be addressable in it. This
// is the assertion whose failure made `wisp slo -state X` fail closed.
func TestSystemProcessSnapshotSeesSelf(t *testing.T) {
	snap, err := SystemProcessSnapshot()
	if err != nil {
		t.Fatalf("SystemProcessSnapshot: %v", err)
	}
	self := uint32(os.Getpid())
	me, ok := snap[self]
	if !ok {
		t.Fatalf("system snapshot does not contain the sampling process %d (%d entries)", self, len(snap))
	}
	if me.PrivateWorkingSet <= 0 {
		t.Errorf("self private working set = %d, want > 0", me.PrivateWorkingSet)
	}
	if me.HandleCount == 0 || me.ThreadCount == 0 {
		t.Errorf("self decoded as handles=%d threads=%d, want both > 0", me.HandleCount, me.ThreadCount)
	}
	if me.UserTime100ns+me.KernelTime100ns <= 0 {
		t.Errorf("self cpu times decoded as user=%d kernel=%d", me.UserTime100ns, me.KernelTime100ns)
	}
}

// TestWalkSystemProcessesEmptyBuffer pins the fail-closed inputs.
func TestWalkSystemProcessesEmptyBuffer(t *testing.T) {
	if err := WalkSystemProcesses(nil, func(SystemProcEntry) {}); err == nil {
		t.Error("WalkSystemProcesses(nil) reported success")
	}
	if err := WalkSystemProcesses(make([]byte, snapshotStride), nil); err == nil {
		t.Error("WalkSystemProcesses accepted a nil callback")
	}
	if _, err := parseSystemProcesses(make([]byte, systemProcEntryFloor)); err == nil {
		t.Error("parseSystemProcesses accepted an all-zero buffer with no pid")
	}
}

func keysOf(m map[uint32]SysProcSample) []uint32 {
	out := make([]uint32, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
