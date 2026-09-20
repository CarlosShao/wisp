//go:build windows

package proc

import (
	"fmt"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/observe"
)

// ExternalSampler implements observe.TreeReader for a SUBJECT process that is
// not the caller (ticket 66 AC#3, registry A15).
//
// Why this exists: the D32 `Sleeping` CPU gate (<=0.5% all-core) was measured
// by a sampler running INSIDE the measured tree, so the number was the
// observer's own burn - one NtQuerySystemInformation read costs ~1.3ms of CPU
// (docs/SLO.md appendix B.2), which at a 250ms cadence is 0.52% of one core's
// share, i.e. above the gate itself. Under that basis the gate has no
// definition, so the product-side CPU row moves out of the tree: the parent
// reads a child pid and the observer's cost lands on the observer.
//
// Every field comes from outside the subject: the private working set, thread
// count, handle count and CPU times come from the ONE system snapshot decoder
// (WalkSystemProcesses - shared with the tree sampler and cmd/balldebug, never
// re-implemented), commit charge / GDI / USER / write counters from a
// PROCESS_QUERY_LIMITED_INFORMATION handle on the subject. The scan buffer is
// allocated by the measuring process, so it can never inflate the subject's
// footprint (the pitfall docs/SLO.md §7 records for self-measurement).
//
// Fail-closed by design: the subject's counters are a fixed pid set, so a
// snapshot that cannot see a subject means the subject died or the view is
// untrustworthy - both are errors, never a zero-byte pass.
type ExternalSampler struct {
	pids []uint32
}

// NewExternalSampler binds the reader to one or more subject pids. Passing the
// calling process's own pid is rejected at read time (that is exactly the
// self-measurement A15 bans).
func NewExternalSampler(pids ...uint32) *ExternalSampler {
	cp := make([]uint32, len(pids))
	copy(cp, pids)
	return &ExternalSampler{pids: cp}
}

// PIDs returns the configured subject pids (copy).
func (r *ExternalSampler) PIDs() []uint32 {
	out := make([]uint32, len(r.pids))
	copy(out, r.pids)
	return out
}

var _ observe.TreeReader = (*ExternalSampler)(nil)

// ReadTree reads the subject set from outside it.
func (r *ExternalSampler) ReadTree() (observe.TreeMetrics, error) {
	if len(r.pids) == 0 {
		return observe.TreeMetrics{}, observe.New(observe.ClassResource,
			"proc: external sampler has no subject pid")
	}
	self := windows.GetCurrentProcessId()
	for _, pid := range r.pids {
		if pid == self {
			return observe.TreeMetrics{}, observe.New(observe.ClassResource,
				fmt.Sprintf("proc: external sampler must not measure the observer itself (pid %d)", self))
		}
	}

	snap, err := SystemProcessSnapshot()
	if err != nil {
		return observe.TreeMetrics{}, observe.Wrap(observe.ClassResource, err, "proc: external system process snapshot")
	}

	m := observe.TreeMetrics{PIDs: len(r.pids)}
	subject := make(map[uint32]bool, len(r.pids))
	for _, pid := range r.pids {
		si, ok := snap[pid]
		if !ok {
			return observe.TreeMetrics{}, observe.New(observe.ClassResource,
				fmt.Sprintf("proc: system snapshot does not contain subject %d", pid))
		}
		subject[pid] = true
		m.PrivateWorkingSetBytes += si.PrivateWorkingSet
		m.CPUTotalNanos += (si.KernelTime100ns + si.UserTime100ns) * 100
		m.Handles += int(si.HandleCount)
		m.Threads += int(si.ThreadCount)

		// Counters the snapshot does not carry. Opening the subject costs the
		// OBSERVER CPU, never the subject's.
		h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
		if err != nil {
			return observe.TreeMetrics{}, observe.New(observe.ClassResource,
				fmt.Sprintf("proc: external OpenProcess(%d): %v", pid, err))
		}
		m.GDIObjects += int(getGuiResources(h, grGDIObjects))
		m.USERObjects += int(getGuiResources(h, grUSERObjects))
		wops, ok := getWriteOpCount(h)
		_ = windows.CloseHandle(h)
		if !ok {
			return observe.TreeMetrics{}, observe.New(observe.ClassResource,
				fmt.Sprintf("proc: external GetProcessIoCounters(%d) failed", pid))
		}
		// The subject is not an observer, so its writes gate: SelfWriteOps
		// stays 0 (that field exists for the in-tree reader, which must not
		// fail its own window).
		m.WriteOps += wops
		if bytes, ok := privateBytesByPid(pid); ok {
			m.CommitBytes += bytes
		}
	}
	m.TCPConnections = treeTCPConnections(subject)
	return m, nil
}
