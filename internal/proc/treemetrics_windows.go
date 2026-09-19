//go:build windows

package proc

import (
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/observe"
)

// TreeSampler implements observe.TreeReader over the C30 Job Object
// (ticket 08): the full D42#10 metric set - private bytes, CPU time, GDI
// and USER objects, handles, threads, write ops, TCP connections - for the
// whole PROCESS TREE. It lives beside JobScope because the D32 memory
// metric's implementation (C30) is the single source for the tree scope;
// per-process reads use PROCESS_QUERY_LIMITED_INFORMATION only, so sampling
// never disturbs the children.
type TreeSampler struct {
	job     *JobScope
	selfPID uint32
}

// NewTreeSampler binds the sampler to a JobScope.
func NewTreeSampler(job *JobScope) *TreeSampler {
	return &TreeSampler{job: job, selfPID: windows.GetCurrentProcessId()}
}

var (
	moduser32                 = windows.NewLazySystemDLL("user32.dll")
	procGetGuiResources       = moduser32.NewProc("GetGuiResources")
	modkernel32               = windows.NewLazySystemDLL("kernel32.dll")
	procGetProcessHandleCount = modkernel32.NewProc("GetProcessHandleCount")
	procGetProcessIoCounters  = modkernel32.NewProc("GetProcessIoCounters")
	modiphlpapi               = windows.NewLazySystemDLL("iphlpapi.dll")
	procGetExtendedTcpTable   = modiphlpapi.NewProc("GetExtendedTcpTable")
)

const (
	// GetGuiResources object classes.
	grGDIObjects  = 0
	grUSERObjects = 1
	// GetExtendedTcpTable: AF_INET + TCP_TABLE_OWNER_PID_ALL.
	tcpTableOwnerPIDAll = 5
	tcpRowOwnerPIDSize  = 24 // MIB_TCPROW_OWNER_PID
	tcpOwningPidOffset  = 20
)

// ReadTree implements observe.TreeReader. Processes that vanish between the
// pid listing and their per-pid read are skipped (best-effort instant
// sample, matching TreePrivateBytes semantics).
func (r *TreeSampler) ReadTree() (observe.TreeMetrics, error) {
	if r.job == nil || r.job.Closed() {
		return observe.TreeMetrics{}, observe.New(observe.ClassResource, "proc: job scope closed")
	}
	pids, err := r.job.TreePIDs()
	if err != nil {
		return observe.TreeMetrics{}, observe.Wrap(observe.ClassResource, err, "proc: job pid list")
	}
	private, err := r.job.TreePrivateBytes()
	if err != nil {
		return observe.TreeMetrics{}, observe.Wrap(observe.ClassResource, err, "proc: tree private bytes")
	}

	inTree := make(map[uint32]bool, len(pids))
	for _, pid := range pids {
		inTree[pid] = true
	}

	m := observe.TreeMetrics{PIDs: len(pids), PrivateBytes: private}
	for _, pid := range pids {
		h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
		if err != nil {
			continue // vanished between listing and open
		}
		self := pid == r.selfPID
		// CPU: kernel+user Filetimes are DURATIONS in 100ns units (not
		// epochs) - convert manually, Filetime.Nanoseconds assumes an epoch.
		var ctime, etime, ktime, utime windows.Filetime
		if windows.GetProcessTimes(h, &ctime, &etime, &ktime, &utime) == nil {
			m.CPUTotalNanos += filetimeDuration(ktime) + filetimeDuration(utime)
		}
		// GDI / USER objects (D42#10: the leak class RSS cannot see).
		m.GDIObjects += int(getGuiResources(h, grGDIObjects))
		m.USERObjects += int(getGuiResources(h, grUSERObjects))
		// Handle count.
		if n, ok := getHandleCount(h); ok {
			m.Handles += n
		}
		// IO write operations (Sleeping zero-periodic-disk-write input).
		if wops, ok := getWriteOpCount(h); ok {
			if self {
				m.SelfWriteOps += wops
			} else {
				m.WriteOps += wops
			}
		}
		_ = windows.CloseHandle(h)
	}
	// Thread count: one toolhelp snapshot for the whole tree.
	m.Threads = treeThreadCount(inTree)
	// TCP connections owned by tree processes.
	m.TCPConnections = treeTCPConnections(inTree)
	return m, nil
}

// filetimeDuration converts a duration-style FILETIME (100ns units) to ns.
func filetimeDuration(ft windows.Filetime) int64 {
	return (int64(ft.HighDateTime)<<32 | int64(ft.LowDateTime)) * 100
}

func getGuiResources(h windows.Handle, class uint32) uint32 {
	r1, _, _ := procGetGuiResources.Call(uintptr(h), uintptr(class))
	return uint32(r1)
}

func getHandleCount(h windows.Handle) (int, bool) {
	var count uint32
	r1, _, _ := procGetProcessHandleCount.Call(uintptr(h), uintptr(unsafe.Pointer(&count)))
	if r1 == 0 {
		return 0, false
	}
	return int(count), true
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

// treeThreadCount counts threads owned by tree processes via one
// TH32CS_SNAPTHREAD snapshot.
func treeThreadCount(inTree map[uint32]bool) int {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return 0
	}
	defer windows.CloseHandle(snap)
	n := 0
	te := windows.ThreadEntry32{Size: uint32(unsafe.Sizeof(windows.ThreadEntry32{}))}
	if err := windows.Thread32First(snap, &te); err != nil {
		return 0
	}
	for {
		if inTree[te.OwnerProcessID] {
			n++
		}
		if err := windows.Thread32Next(snap, &te); err != nil {
			break
		}
	}
	return n
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
