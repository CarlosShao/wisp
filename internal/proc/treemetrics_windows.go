//go:build windows

package proc

import (
	"fmt"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/observe"
)

// TreeSampler implements observe.TreeReader over the C30 Job Object
// (ticket 08): the full D42#10 metric set for the whole PROCESS TREE.
//
// Memory units (docs/SLO.md §7): the D32 table's numbers were measured as
// PRIVATE WORKING SET (Task Manager "Memory (private working set)",
// WorkingSetPrivateSize). This sampler therefore gates on private working
// set and records the commit charge (PrivateUsage, the C30 Job Object
// sum) alongside - commit is always >= private WS (reserved-but-resident
// Go heap), so gating on commit against SLO.md numbers would false-fail
// every state.
//
// Per-sample cost: ONE NtQuerySystemInformation(SystemProcessInformation)
// snapshot supplies private working set, thread count, handle count and
// CPU times for every pid (no per-pid opens, no toolhelp walk); per-pid
// opens happen only for GDI/USER objects and IO write counters. The
// observer's own CPU burn is therefore negligible at the 250ms cadence.
type TreeSampler struct {
	job     *JobScope
	selfPID uint32
}

// NewTreeSampler binds the sampler to a JobScope.
func NewTreeSampler(job *JobScope) *TreeSampler {
	return &TreeSampler{job: job, selfPID: windows.GetCurrentProcessId()}
}

var (
	moduser32                    = windows.NewLazySystemDLL("user32.dll")
	procGetGuiResources          = moduser32.NewProc("GetGuiResources")
	modkernel32                  = windows.NewLazySystemDLL("kernel32.dll")
	procGetProcessIoCounters     = modkernel32.NewProc("GetProcessIoCounters")
	modiphlpapi                  = windows.NewLazySystemDLL("iphlpapi.dll")
	procGetExtendedTcpTable      = modiphlpapi.NewProc("GetExtendedTcpTable")
	modntdll                     = windows.NewLazySystemDLL("ntdll.dll")
	procNtQuerySystemInformation = modntdll.NewProc("NtQuerySystemInformation")
)

const (
	// GetGuiResources object classes.
	grGDIObjects  = 0
	grUSERObjects = 1
	// GetExtendedTcpTable: AF_INET + TCP_TABLE_OWNER_PID_ALL.
	tcpTableOwnerPIDAll = 5
	tcpRowOwnerPIDSize  = 24 // MIB_TCPROW_OWNER_PID
	tcpOwningPidOffset  = 20
	// NtQuerySystemInformation class SystemProcessInformation.
	systemProcessInformation = 5
	statusInfoLengthMismatch = 0xC0000004
)

// ReadTree implements observe.TreeReader. The sampled tree is the MAIN
// PROCESS plus every process currently assigned to the Job (C30 assigns
// children; the wisp main process is the tree root and always included -
// assigning self to a KILL_ON_JOB_CLOSE job would kill the process at
// graceful shutdown). Processes that vanish between reads are skipped
// (best-effort instant sample).
func (r *TreeSampler) ReadTree() (observe.TreeMetrics, error) {
	if r.job == nil || r.job.Closed() {
		return observe.TreeMetrics{}, observe.New(observe.ClassResource, "proc: job scope closed")
	}
	pids, err := r.job.TreePIDs()
	if err != nil {
		return observe.TreeMetrics{}, observe.Wrap(observe.ClassResource, err, "proc: job pid list")
	}

	selfInJob := false
	inTree := make(map[uint32]bool, len(pids)+1)
	for _, pid := range pids {
		inTree[pid] = true
		if pid == r.selfPID {
			selfInJob = true
		}
	}
	inTree[r.selfPID] = true

	// System-wide snapshot: private WS, threads, handles, CPU per pid. The
	// tree root (self) MUST be present - a snapshot that cannot see the
	// sampling process itself is not a trustworthy view; retry once, then
	// fail the read (fail-closed: an unmeasurable window must never turn
	// into a silent 0-byte pass).
	var snap map[uint32]SysProcSample
	var lastErr error
	for attempt := 0; attempt < 60; attempt++ {
		snap, err = SystemProcessSnapshot()
		if err != nil {
			lastErr = err
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if _, ok := snap[r.selfPID]; ok {
			break
		}
		// Snapshot succeeded but cannot see the sampling process itself:
		// not a trustworthy view (fail-closed), retry briefly.
		lastErr = observe.New(observe.ClassResource,
			"proc: system snapshot does not contain the sampling process")
		snap = nil
		time.Sleep(100 * time.Millisecond)
	}
	if snap == nil {
		if lastErr == nil {
			lastErr = observe.New(observe.ClassResource, "proc: system snapshot unavailable")
		}
		return observe.TreeMetrics{}, observe.Wrap(observe.ClassResource, lastErr, "proc: system process snapshot")
	}

	m := observe.TreeMetrics{PIDs: len(inTree)}
	for pid := range inTree {
		si, ok := snap[pid]
		if !ok {
			continue // vanished between listing and snapshot
		}
		m.PrivateWorkingSetBytes += si.PrivateWorkingSet
		m.CPUTotalNanos += (si.KernelTime100ns + si.UserTime100ns) * 100
		m.Handles += int(si.HandleCount)
		m.Threads += int(si.ThreadCount)
	}

	// Commit charge: the C30 Job sum covers children; self is added unless
	// the caller already assigned it (dedupe, never double count).
	commit, err := r.job.TreePrivateBytes()
	if err != nil {
		return observe.TreeMetrics{}, observe.Wrap(observe.ClassResource, err, "proc: tree private bytes")
	}
	if !selfInJob {
		if selfBytes, ok := PrivateBytesByPID(r.selfPID); ok {
			commit += selfBytes
		}
	}
	m.CommitBytes = commit

	// Per-pid opens only where no snapshot field exists.
	for pid := range inTree {
		h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
		if err != nil {
			continue
		}
		self := pid == r.selfPID
		m.GDIObjects += int(getGuiResources(h, grGDIObjects))
		m.USERObjects += int(getGuiResources(h, grUSERObjects))
		if wops, ok := getWriteOpCount(h); ok {
			if self {
				m.SelfWriteOps += wops
			} else {
				m.WriteOps += wops
			}
		}
		_ = windows.CloseHandle(h)
	}
	// TCP connections owned by tree processes.
	m.TCPConnections = treeTCPConnections(inTree)
	return m, nil
}

func getGuiResources(h windows.Handle, class uint32) uint32 {
	r1, _, _ := procGetGuiResources.Call(uintptr(h), uintptr(class))
	return uint32(r1)
}

// ioCounters mirrors IO_COUNTERS (kernel32 GetProcessIoCounters).
type ioCounters struct {
	ReadTransferCount   uint64
	WriteTransferCount  uint64
	OtherTransferCount  uint64
	ReadOperationCount  uint64
	WriteOperationCount uint64
	OtherOperationCount uint64
}

func getWriteOpCount(h windows.Handle) (int64, bool) {
	var io ioCounters
	r1, _, _ := procGetProcessIoCounters.Call(uintptr(h), uintptr(unsafe.Pointer(&io)))
	if r1 == 0 {
		return 0, false
	}
	return int64(io.WriteOperationCount), true
}

// SysProcSample is one process's entry from the system snapshot. It is
// exported because the out-of-tree readers (cmd/balldebug -diff, `wisp slo`
// -subject measurement) need the same fields from the same snapshot: one
// decoder, one record type, no second implementation (ticket 66 / A14).
type SysProcSample struct {
	PrivateWorkingSet int64 // WorkingSetPrivateSize, bytes
	KernelTime100ns   int64
	UserTime100ns     int64
	HandleCount       uint32
	ThreadCount       uint32
}

// SystemProcessSnapshot reads SystemProcessInformation once and extracts
// the fields the D42#10 metric set needs per pid. The scan buffer is a Go
// slice: for the out-of-tree readers that is free (the scratch belongs to
// the measuring process, never the subject); the in-tree reader's own
// footprint is dominated by the image map, not this scratch.
func SystemProcessSnapshot() (map[uint32]SysProcSample, error) {
	size := uint32(1 << 20)
	for attempt := 0; attempt < 6; attempt++ {
		buf := make([]byte, size)
		retLen := uint32(0)
		status, _, _ := procNtQuerySystemInformation.Call(
			uintptr(systemProcessInformation),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(size),
			uintptr(unsafe.Pointer(&retLen)),
		)
		if status == 0 {
			return parseSystemProcesses(buf)
		}
		if status != uintptr(statusInfoLengthMismatch) {
			return nil, fmt.Errorf("NtQuerySystemInformation: NTSTATUS 0x%08x", uint32(status))
		}
		if retLen > size {
			size = retLen
		} else {
			size *= 2
		}
	}
	return nil, fmt.Errorf("NtQuerySystemInformation: buffer never sufficient (last %d bytes)", size)
}

// parseSystemProcesses walks the variable-length entries (through the one
// decoder in the repo, proc.WalkSystemProcesses - see systemprocs_windows.go
// for why the last entry must be delivered before the NextEntryOffset == 0
// test).
func parseSystemProcesses(buf []byte) (map[uint32]SysProcSample, error) {
	out := make(map[uint32]SysProcSample, 256)
	err := WalkSystemProcesses(buf, func(e SystemProcEntry) {
		if e.PID == 0 || e.PID > 0xFFFFFFFF {
			return // "System Idle Process" / non-pid handle values
		}
		out[uint32(e.PID)] = SysProcSample{
			PrivateWorkingSet: e.PrivateWorkingSet,
			KernelTime100ns:   e.KernelTime100ns,
			UserTime100ns:     e.UserTime100ns,
			HandleCount:       e.HandleCount,
			ThreadCount:       e.ThreadCount,
		}
	})
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("system process snapshot parsed zero entries")
	}
	return out, nil
}

// treeTCPConnections counts TCP connections owned by tree processes via
// GetExtendedTcpTable (TCP_TABLE_OWNER_PID_ALL). The table read retries on
// an insufficient buffer; after exhausted retries it reports 0 - a real
// long-lived connection persists across the sampler's re-reads every
// interval, so a transient table failure cannot mask it (documented
// limitation in the slo-check output schema).
func treeTCPConnections(inTree map[uint32]bool) int {
	size := uint32(32 << 10)
	for attempt := 0; attempt < 4; attempt++ {
		buf := make([]byte, size)
		r1, _, _ := procGetExtendedTcpTable.Call(
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&size)),
			0, // bOrder = FALSE (counting only, no sort needed)
			uintptr(windows.AF_INET),
			uintptr(tcpTableOwnerPIDAll),
			0,
		)
		if r1 == 0 {
			entries := *(*uint32)(unsafe.Pointer(&buf[0]))
			n := 0
			for i := uint32(0); i < entries; i++ {
				off := 4 + uintptr(i)*tcpRowOwnerPIDSize + tcpOwningPidOffset
				pid := *(*uint32)(unsafe.Add(unsafe.Pointer(&buf[0]), off))
				if inTree[pid] {
					n++
				}
			}
			return n
		}
		if size > (4 << 20) {
			break
		}
	}
	return 0
}
