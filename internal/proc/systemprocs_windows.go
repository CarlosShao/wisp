//go:build windows

package proc

import (
	"fmt"
	"unsafe"
)

// Single decoder for the system process snapshot (ticket 66, registry A14).
//
// There is exactly ONE walk of the variable-length
// SYSTEM_PROCESS_INFORMATION chain in this repo: WalkSystemProcesses. The
// tree sampler (internal/proc) and the out-of-tree readers (cmd/wisp slo,
// cmd/balldebug -diff) all go through it, so the ordering bug that A14
// records can never be reintroduced in a copy.
//
// Layout is the documented x64 SYSTEM_PROCESS_INFORMATION (Vista/Win7
// additions included), read from a buffer filled by
// NtQuerySystemInformation(SystemProcessInformation):
//
//	+0x00 NextEntryOffset      ULONG   (0 = LAST entry of the chain)
//	+0x04 NumberOfThreads      ULONG
//	+0x08 WorkingSetPrivateSize LARGE_INTEGER
//	+0x28 UserTime             LARGE_INTEGER (100ns, since process start)
//	+0x30 KernelTime           LARGE_INTEGER (100ns, since process start)
//	+0x50 UniqueProcessId      HANDLE
//	+0x60 HandleCount          ULONG
//
// THE ORDERING RULE (this is what A14 was): NextEntryOffset == 0 is the
// LAST ENTRY MARKER, not a stop-before-this-entry signal. An entry is
// therefore always decoded and handed to fn BEFORE the terminator test.
// The pre-fix code tested next == 0 first and broke out, which made the
// snapshot's tail process permanently invisible - and newly created
// processes sit at the tail of ActiveProcessLinks, i.e. the process doing
// the sampling. cmd/balldebug/diff_windows.go::privateWorkingSetFor had
// the order right all along; this function is that shape, generalised.

// SystemProcEntry is one decoded entry of a system process snapshot.
type SystemProcEntry struct {
	// PID is UniqueProcessId. The "System Idle Process" entry carries 0.
	PID uint64
	// HandleCount is the process handle count (+0x60).
	HandleCount uint32
	// ThreadCount is NumberOfThreads (+0x04).
	ThreadCount uint32
	// PrivateWorkingSet is WorkingSetPrivateSize in bytes - Task Manager's
	// "Memory (private working set)", the D32 gate unit (docs/SLO.md §7).
	PrivateWorkingSet int64
	// KernelTime100ns / UserTime100ns are cumulative CPU since process start.
	KernelTime100ns int64
	UserTime100ns   int64
}

// systemProcEntryFloor is the byte span the decoder reads per entry
// (+0x60 HandleCount + 4).
const systemProcEntryFloor = 112

// WalkSystemProcesses decodes every entry of a SystemProcessInformation
// snapshot buffer in chain order and hands each to fn. It is the only
// place in this repo that interprets that layout.
func WalkSystemProcesses(buf []byte, fn func(SystemProcEntry)) error {
	if len(buf) == 0 {
		return fmt.Errorf("system process snapshot buffer is empty")
	}
	if fn == nil {
		return fmt.Errorf("system process snapshot walker has no callback")
	}
	offset := uintptr(0)
	base := unsafe.Pointer(&buf[0])
	for {
		if offset+4 > uintptr(len(buf)) {
			return nil // buffer exhausted: chain ended without a terminator
		}
		next := *(*uint32)(unsafe.Add(base, offset))
		if offset+systemProcEntryFloor > uintptr(len(buf)) {
			return fmt.Errorf("system process entry truncated at offset %d", offset)
		}
		e := SystemProcEntry{
			ThreadCount:       *(*uint32)(unsafe.Add(base, offset+4)),
			PrivateWorkingSet: int64(*(*int64)(unsafe.Add(base, offset+8))),
			UserTime100ns:     int64(*(*int64)(unsafe.Add(base, offset+40))),
			KernelTime100ns:   int64(*(*int64)(unsafe.Add(base, offset+48))),
			PID:               *(*uint64)(unsafe.Add(base, offset+80)),
			HandleCount:       *(*uint32)(unsafe.Add(base, offset+96)),
		}
		fn(e) // <- last entry included: deliver first, test the terminator after
		if next == 0 {
			return nil
		}
		if next < 8 {
			return fmt.Errorf("system process entry at offset %d does not advance (NextEntryOffset=%d)", offset, next)
		}
		offset += uintptr(next)
	}
}
