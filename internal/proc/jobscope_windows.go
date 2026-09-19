//go:build windows

package proc

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

// JobScope — contract C30 (D32, D38e step 9, SPEC-01 §2/§3).
//
// The main process owns ONE Job Object; every child process (WebView2 host
// with its 3-5 msedgewebview2 children, the path-X speech subprocess, the
// updater) is assigned to it:
//   - JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE makes the OS kill the whole tree if
//     the main process dies (no orphan msedgewebview2.exe resident);
//   - TreePrivateBytes() is THE D32 SLO metric (process-tree private memory,
//     not single-process RSS).
//
// The measurement is summed per process (psapi GetProcessMemoryInfo,
// PrivateUsage) over the job's process id list, which is the documented
// "process tree private bytes" accounting the D32 tables are written against.

// jobObjectBasicProcessIdList is JOB_OBJECT class 3 (winbase.h); x/sys/windows
// does not export it.
const jobObjectBasicProcessIdList = 3

// jobPidListCapacity bounds the pid-list buffer. Wisp trees are small
// (WebView2 <=5 children + speech subprocess + updater); 256 is generous and
// one allocation.
const jobPidListCapacity = 256

var (
	modpsapi                 = windows.NewLazySystemDLL("psapi.dll")
	procGetProcessMemoryInfo = modpsapi.NewProc("GetProcessMemoryInfo")
)

// processMemoryCountersEx mirrors PROCESS_MEMORY_COUNTERS_EX; x/sys/windows
// does not export it. PrivateUsage is the private (committed) bytes the D32
// metric sums.
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
	PrivateUsage               uintptr
}

// JobScope owns the process-wide Job Object. Create with OpenJobScope; the
// zero value is not usable.
type JobScope struct {
	mu     sync.Mutex
	handle windows.Handle
	closed bool
}

// OpenJobScope creates the Job Object with KILL_ON_JOB_CLOSE. The caller
// MUST Close it during shutdown (D38e step 9) - normally via the shutdown
// sequence, which treats it as non-skippable.
func OpenJobScope() (*JobScope, error) {
	h, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("proc: CreateJobObject: %w", err)
	}
	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
	info.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	if _, err := windows.SetInformationJobObject(
		h, windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info)),
	); err != nil {
		_ = windows.CloseHandle(h)
		return nil, fmt.Errorf("proc: SetInformationJobObject(KILL_ON_JOB_CLOSE): %w", err)
	}
	return &JobScope{handle: h}, nil
}

// Assign puts an already-started process into the Job.
func (j *JobScope) Assign(p *os.Process) error {
	if p == nil {
		return fmt.Errorf("proc: Assign: nil process")
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.closed {
		return fmt.Errorf("proc: Assign: job already closed")
	}
	var assignErr error
	_ = p.WithHandle(func(handle uintptr) {
		if err := windows.AssignProcessToJobObject(j.handle, windows.Handle(handle)); err != nil {
			assignErr = fmt.Errorf("proc: AssignProcessToJobObject(pid %d): %w", p.Pid, err)
		}
	})
	return assignErr
}

// StartInJob starts cmd and immediately assigns it to the Job. If the
// assignment fails the child is killed before the error returns, so a child
// can never escape the Job by accident.
func (j *JobScope) StartInJob(cmd *exec.Cmd) (*os.Process, error) {
	if cmd == nil {
		return nil, fmt.Errorf("proc: StartInJob: nil cmd")
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("proc: StartInJob: %w", err)
	}
	if err := j.Assign(cmd.Process); err != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		return nil, err
	}
	return cmd.Process, nil
}

// TreePrivateBytes returns the sum of private (committed) bytes over every
// process currently in the Job - the D32 SLO metric. Processes that vanish
// between listing and measurement are skipped (best-effort sample, matching
// what Task Manager shows for the tree at that instant).
func (j *JobScope) TreePrivateBytes() (int64, error) {
	pids, err := j.pids()
	if err != nil {
		return 0, err
	}
	var total int64
	for _, pid := range pids {
		pb, ok := privateBytesByPid(pid)
		if ok {
			total += pb
		}
	}
	return total, nil
}

// TreeProcessCount returns the number of processes currently in the Job
// (diagnostics companion of the D42#10 sampling set).
func (j *JobScope) TreeProcessCount() (int, error) {
	pids, err := j.pids()
	if err != nil {
		return 0, err
	}
	return len(pids), nil
}

// Close closes the Job handle. Per KILL_ON_JOB_CLOSE the OS kills every
// process still in the Job. Idempotent.
func (j *JobScope) Close() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.closed {
		return nil
	}
	j.closed = true
	if err := windows.CloseHandle(j.handle); err != nil {
		return fmt.Errorf("proc: close job: %w", err)
	}
	return nil
}

// Closed reports whether the Job handle has been closed.
func (j *JobScope) Closed() bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.closed
}

func (j *JobScope) pids() ([]uint32, error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.closed {
		return nil, fmt.Errorf("proc: job is closed")
	}
	entrySize := unsafe.Sizeof(uintptr(0))
	buf := make([]byte, 8+uintptr(jobPidListCapacity)*entrySize)
	var retLen uint32
	if err := windows.QueryInformationJobObject(
		j.handle, jobObjectBasicProcessIdList,
		uintptr(unsafe.Pointer(&buf[0])), uint32(len(buf)), &retLen,
	); err != nil {
		return nil, fmt.Errorf("proc: QueryInformationJobObject(pid list): %w", err)
	}
	// JOB_OBJECT_BASIC_PROCESS_ID_LIST layout (x64):
	// DWORD NumberOfAssignedProcesses; DWORD NumberOfProcessIdsInList;
	// ULONG_PTR ProcessIdList[].
	n := *(*uint32)(unsafe.Pointer(&buf[4]))
	if n > jobPidListCapacity {
		n = jobPidListCapacity
	}
	pids := make([]uint32, 0, n)
	base := unsafe.Pointer(&buf[8])
	for i := uint32(0); i < n; i++ {
		pid := *(*uintptr)(unsafe.Add(base, uintptr(i)*entrySize))
		if pid != 0 {
			pids = append(pids, uint32(pid))
		}
	}
	return pids, nil
}

func privateBytesByPid(pid uint32) (int64, bool) {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return 0, false // process vanished between listing and open
	}
	defer windows.CloseHandle(h)

	var pmc processMemoryCountersEx
	pmc.CB = uint32(unsafe.Sizeof(pmc))
	r1, _, callErr := procGetProcessMemoryInfo.Call(
		uintptr(h),
		uintptr(unsafe.Pointer(&pmc)),
		uintptr(pmc.CB),
	)
	if r1 == 0 {
		// err is always non-nil for LazyProc.Call; only r1 decides.
		_ = callErr
		return 0, false
	}
	return int64(pmc.PrivateUsage), true
}
