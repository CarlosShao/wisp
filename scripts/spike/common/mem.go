// Package common provides measurement primitives for the S0 spike programs
// (ticket 02). Throwaway quality is acceptable; the JSON outputs are the
// deliverable (see docs/evidence/s0/02-spike-report.md).
package common

import (
	"runtime"
	"runtime/debug"
	"sort"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modpsapi   = windows.NewLazySystemDLL("psapi.dll")
	prock32    = windows.NewLazySystemDLL("kernel32.dll")
	moduser32  = windows.NewLazySystemDLL("user32.dll")
	procGetPMI = modpsapi.NewProc("GetProcessMemoryInfo")
	procGGR    = moduser32.NewProc("GetGuiResources")
	procGPHC   = prock32.NewProc("GetProcessHandleCount")
	procGSI    = prock32.NewProc("GetSystemInfo")
)

// Layout mirrors PROCESS_MEMORY_COUNTERS_EX: cb(4) + PageFaultCount(4) then
// 9 SIZE_T fields (natural 8-byte alignment) = 80 bytes on x64.
type processMemoryCountersEx struct {
	CB                         uint32
	PageFaultCount             uint32
	PeakWorkingSetSize         uintptr
	WorkingSetSize             uintptr
	QuotaPeakPagedPoolUsage    uintptr
	QuotaPagedPoolUsage        uintptr
	QuotaPeakNonPagedPoolUsage uintptr
	QuotaNonPagedPoolUsage     uintptr
	PagefileUsage              uintptr
	PeakPagefileUsage          uintptr
	PrivateUsage               uintptr // PROCESS_MEMORY_COUNTERS_EX extension
}

// MemSample is one instant of memory/resource accounting for THIS process.
//
// Definitions (all in bytes):
//   - WorkingSet: total resident pages (GetProcessMemoryInfo.WorkingSetSize).
//   - PrivateWorkingSet: pages resident AND not shared with any other process
//     (QueryWorkingSetEx, Shared bit clear per page) == Task Manager's
//     "Memory (private working set)". This is the D25/D32 decision metric.
//   - SharedWorkingSet: WorkingSet - PrivateWorkingSet (DLL images, font files etc).
//   - PrivateCommit: PrivateUsage from PROCESS_MEMORY_COUNTERS_EX (charge, may
//     exceed what is resident; reported for reference).
type MemSample struct {
	WorkingSet        uint64 `json:"workingSetMB_f"`
	PrivateWorkingSet uint64 `json:"privateWorkingSetMB_f"`
	SharedWorkingSet  uint64 `json:"sharedWorkingSetMB_f"`
	PrivateCommit     uint64 `json:"privateCommitMB_f"`
	GDIObjects        uint32 `json:"gdiObjects"`
	UserObjects       uint32 `json:"userObjects"`
	Handles           uint32 `json:"handles"`
	Threads           uint32 `json:"threads"`
}

// SampleMem takes one measurement of the current process.
func SampleMem() MemSample {
	// TRACE-removed

	var s MemSample

	h := windows.CurrentProcess()

	var pmc processMemoryCountersEx
	pmc.CB = uint32(unsafe.Sizeof(pmc))
	r1, _, _ := procGetPMI.Call(uintptr(h), uintptr(unsafe.Pointer(&pmc)), uintptr(pmc.CB))
	if r1 != 0 {
		s.WorkingSet = uint64(pmc.WorkingSetSize)
		s.PrivateCommit = uint64(pmc.PrivateUsage)
	}
	s.PrivateWorkingSet = privateWorkingSet(windows.GetCurrentProcessId())
	if s.PrivateWorkingSet > s.WorkingSet { // defensive: never exceed total
		s.PrivateWorkingSet = s.WorkingSet
	}
	s.SharedWorkingSet = s.WorkingSet - s.PrivateWorkingSet

	if r1, _, _ := procGGR.Call(uintptr(h), 0); r1 != 0xFFFFFFFF { // GR_GDIOBJECTS
		s.GDIObjects = uint32(r1)
	}
	if r1, _, _ := procGGR.Call(uintptr(h), 1); r1 != 0xFFFFFFFF { // GR_USEROBJECTS
		s.UserObjects = uint32(r1)
	}
	var hc uint32
	if r1, _, _ := procGPHC.Call(uintptr(h), uintptr(unsafe.Pointer(&hc))); r1 != 0 {
		s.Handles = hc
	}
	s.Threads = uint32(threadCount())
	return s
}

func threadCount() int {
	// GetProcessInformation(ProcessMemoryInformation) is overkill; count via
	// SystemProcessInformation is heavy. Use CreateToolhelp32Snapshot once.
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return 0
	}
	defer windows.CloseHandle(snap)
	pid := windows.GetCurrentProcessId()
	var te windows.ThreadEntry32
	te.Size = uint32(unsafe.Sizeof(te))
	n := 0
	err = windows.Thread32First(snap, &te)
	for err == nil {
		if te.OwnerProcessID == pid {
			n++
		}
		err = windows.Thread32Next(snap, &te)
	}
	return n
}

// privateWorkingSet returns the process's private working set in bytes via
// NtQuerySystemInformation(SystemProcessInformation).WorkingSetPrivateSize -
// the same kernel field Task Manager's "Memory (private working set)" column
// shows (this is also how Chromium computes it). QueryWorkingSetEx was tried
// first and silently wrote no data on this machine, so it is not used.
//
// SYSTEM_PROCESS_INFORMATION layout (x64, stable since Vista):
//   0x00 NextEntryOffset, 0x04 NumberOfThreads, 0x08 WorkingSetPrivateSize,
//   0x50 UniqueProcessId ...
func privateWorkingSet(targetPID uint32) uint64 {
	const systemProcessInformation = 5
	const statusInfoLengthMismatch = 0xC0000004
	ntdll := windows.NewLazySystemDLL("ntdll.dll")
	pNtQSI := ntdll.NewProc("NtQuerySystemInformation")

	for _, size := range []int{1 << 22, 8 << 22, 32 << 22} {
		buf := make([]byte, size)
		var retLen uint32
		r1, _, _ := pNtQSI.Call(systemProcessInformation,
			uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)),
			uintptr(unsafe.Pointer(&retLen)))
		if r1 == statusInfoLengthMismatch {
			continue
		}
		if r1 != 0 {
			return 0
		}
		off := uintptr(0)
		for off+0xA8 <= uintptr(retLen) {
			next := *(*uint32)(unsafe.Pointer(&buf[off]))
			pid := *(*uint64)(unsafe.Pointer(&buf[off+0x50]))
			private := *(*int64)(unsafe.Pointer(&buf[off+0x08]))
			if uint32(pid) == targetPID {
				if private < 0 {
					return 0
				}
				return uint64(private)
			}
			if next == 0 {
				break
			}
			off += uintptr(next)
		}
		return 0
	}
	return 0
}

type systemInfo struct {
	_                uint32 // union dwOemId / processor architecture
	PageSize         uint32
	MinAppAddress    uintptr
	MaxAppAddress    uintptr
	ActiveProcMask   uintptr
	NumProcessors    uint32
	ProcessorType    uint32
	AllocGranularity uint32
	_                uint16
	_                uint16
}

func osPageSize() uint32 {
	var si systemInfo
	procGSI.Call(uintptr(unsafe.Pointer(&si)))
	return si.PageSize
}

// SettleGC forces the Go runtime into its idle shape before sampling:
// two full GCs + FreeOSMemory (returns memory to the OS) + a short sleep.
// This is the same primitive C11 will use as its Dispose step 4.
func SettleGC() {
	runtime.GC()
	runtime.GC()
	debug.FreeOSMemory()
}

// SampleStable takes n samples with a pause and returns min/median/max per
// field. Memory numbers can jitter while DWM/DLLs settle; medians are used
// in the report, min/max bound them.
func SampleStable(n int, pauseMs int) (minS, medS, maxS MemSample) {
	if n < 1 {
		n = 1
	}
	samples := make([]MemSample, 0, n)
	for i := 0; i < n; i++ {
		samples = append(samples, SampleMem())
		SleepMs(pauseMs)
	}
	minS, maxS = samples[0], samples[0]
	for _, s := range samples {
		minS.WorkingSet = u64min(minS.WorkingSet, s.WorkingSet)
		minS.PrivateWorkingSet = u64min(minS.PrivateWorkingSet, s.PrivateWorkingSet)
		minS.SharedWorkingSet = u64min(minS.SharedWorkingSet, s.SharedWorkingSet)
		minS.PrivateCommit = u64min(minS.PrivateCommit, s.PrivateCommit)
		maxS.WorkingSet = u64max(maxS.WorkingSet, s.WorkingSet)
		maxS.PrivateWorkingSet = u64max(maxS.PrivateWorkingSet, s.PrivateWorkingSet)
		maxS.SharedWorkingSet = u64max(maxS.SharedWorkingSet, s.SharedWorkingSet)
		maxS.PrivateCommit = u64max(maxS.PrivateCommit, s.PrivateCommit)
	}
	medS = medianOf(samples)
	return
}

func medianOf(s []MemSample) MemSample {
	pick := func(f func(MemSample) uint64) uint64 {
		v := make([]uint64, len(s))
		for i, m := range s {
			v[i] = f(m)
		}
		sort.Slice(v, func(a, b int) bool { return v[a] < v[b] })
		return v[len(v)/2]
	}
	u := make([]uint32, len(s))
	pick32 := func(f func(MemSample) uint32) uint32 {
		for i, m := range s {
			u[i] = f(m)
		}
		sort.Slice(u, func(a, b int) bool { return u[a] < u[b] })
		return u[len(u)/2]
	}
	return MemSample{
		WorkingSet:        pick(func(m MemSample) uint64 { return m.WorkingSet }),
		PrivateWorkingSet: pick(func(m MemSample) uint64 { return m.PrivateWorkingSet }),
		SharedWorkingSet:  pick(func(m MemSample) uint64 { return m.SharedWorkingSet }),
		PrivateCommit:     pick(func(m MemSample) uint64 { return m.PrivateCommit }),
		GDIObjects:        pick32(func(m MemSample) uint32 { return m.GDIObjects }),
		UserObjects:       pick32(func(m MemSample) uint32 { return m.UserObjects }),
		Handles:           pick32(func(m MemSample) uint32 { return m.Handles }),
		Threads:           pick32(func(m MemSample) uint32 { return m.Threads }),
	}
}

func u64min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
func u64max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func SleepMs(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// PeakPrivateCommitMB reports PeakPagefileUsage (peak private commit charge)
// for the current process, in MB. Used to corroborate "peak = max(ASR,TTS),
// not sum" during serial switch cycles.
func PeakPrivateCommitMB() float64 {
	var pmc processMemoryCountersEx
	pmc.CB = uint32(unsafe.Sizeof(pmc))
	r1, _, _ := procGetPMI.Call(uintptr(windows.CurrentProcess()), uintptr(unsafe.Pointer(&pmc)), uintptr(pmc.CB))
	if r1 == 0 {
		return 0
	}
	return float64(pmc.PeakPagefileUsage) / (1 << 20)
}
