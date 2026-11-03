//go:build windows

package manifestslot

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

// childJob is the per-child Windows Job Object behind the D46 "kill the process
// TREE" rule (SPEC-07 §7, PLAN.md:2670).
//
// Why a Job and not Process.Kill: a wrapped CLI is free to spawn helpers
// (lark-cli starts node, git starts an ssh, ffmpeg starts its encoders). Killing
// only the direct child leaves those running, and the ticket's acceptance line
// is "hanging CLI killed at timeout (process-tree gone)". KILL_ON_JOB_CLOSE makes
// closing the handle terminate every process still inside, including
// grandchildren that inherited the job.
type childJob struct {
	mu     sync.Mutex
	handle windows.Handle
	exited windows.Handle
}

// openChildJob creates one Job Object with KILL_ON_JOB_CLOSE.
func openChildJob() (*childJob, error) {
	h, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("manifestslot: CreateJobObject: %w", err)
	}
	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
	info.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	if _, err := windows.SetInformationJobObject(h, windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info))); err != nil {
		_ = windows.CloseHandle(h)
		return nil, fmt.Errorf("manifestslot: SetInformationJobObject: %w", err)
	}
	return &childJob{handle: h}, nil
}

// Start starts the command INSIDE the Job and takes a wait handle so the caller
// can stop waiting the instant the CLI exits (exec.Cmd.Wait alone waits for the
// output copiers too, which a surviving grandchild could hold open forever).
//
// The child is created SUSPENDED and assigned to the Job BEFORE its first
// instruction runs:
//
//	CreateProcess(CREATE_SUSPENDED) -> AssignProcessToJobObject -> ResumeThread
//
// Starting it first and assigning afterwards leaves a window in which the CLI
// could already have spawned a helper that escaped the Job - and the ticket's
// acceptance line is "process-tree gone", not "we killed the ones we happened
// to catch". This is also what mainstream harnesses do on Windows (the
// community fix for the same window is the `CREATE_SUSPENDED` +
// `AssignProcessToJobObject` + resume-nudge sequence, e.g. node's
// libuv job handling and Rust's `CommandExt::creation_flags` users).
//
// If the assignment fails the child is killed before the error returns, so it
// never runs outside the Job.
func (j *childJob) Start(cmd *exec.Cmd) error {
	cmd.SysProcAttr = &windows.SysProcAttr{CreationFlags: windows.CREATE_SUSPENDED}
	if err := cmd.Start(); err != nil {
		return err
	}
	pid := uint32(cmd.Process.Pid)
	if err := j.assign(cmd.Process); err != nil {
		_ = cmd.Process.Kill()
		resumeSuspended(pid)
		return err
	}
	// Only now may the child run: its threads are still suspended at the
	// entry point, so this is the first instruction it executes.
	resumeSuspended(pid)
	h, err := windows.OpenProcess(windows.SYNCHRONIZE, false, pid)
	if err == nil {
		j.mu.Lock()
		j.exited = h
		j.mu.Unlock()
	}
	return nil
}

func (j *childJob) assign(p *os.Process) error {
	var assignErr error
	_ = p.WithHandle(func(handle uintptr) {
		if err := windows.AssignProcessToJobObject(j.handle, windows.Handle(handle)); err != nil {
			assignErr = fmt.Errorf("manifestslot: AssignProcessToJobObject(pid %d): %w", p.Pid, err)
		}
	})
	return assignErr
}

// resumeSuspended nudges every thread of a CREATE_SUSPENDED child so it can
// run. There is no ResumeProcess in the Win32 API, so the threads are walked
// through a Toolhelp snapshot: Go's os/exec does not keep the primary thread
// handle, and a child that never resumes would hang until the D46 budget kills
// it - which would look like "the CLI hangs" instead of "we forgot to resume".
//
// The walk is bounded: a freshly created process has one thread, and a thread
// id that shows up after this snapshot belongs to a thread the child created
// itself (those are never suspended). Each ResumeThread there is a no-op, so
// racing with a thread the child spawns after the resume is harmless.
func resumeSuspended(pid uint32) {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		return
	}
	defer windows.CloseHandle(snap)
	var te windows.ThreadEntry32
	te.Size = uint32(unsafe.Sizeof(te))
	for err := windows.Thread32First(snap, &te); err == nil; err = windows.Thread32Next(snap, &te) {
		if te.OwnerProcessID != pid {
			continue
		}
		h, oerr := windows.OpenThread(windows.THREAD_SUSPEND_RESUME, false, te.ThreadID)
		if oerr != nil {
			continue
		}
		_, _ = windows.ResumeThread(h)
		_ = windows.CloseHandle(h)
	}
}

// WaitExited blocks until the direct child exits.
func (j *childJob) WaitExited() {
	j.mu.Lock()
	h := j.exited
	j.mu.Unlock()
	if h == 0 {
		return
	}
	_, _ = windows.WaitForSingleObject(h, windows.INFINITE)
}

// Kill terminates every process in the Job.
func (j *childJob) Kill() {
	j.mu.Lock()
	h := j.handle
	j.mu.Unlock()
	if h != 0 {
		_ = windows.TerminateJobObject(h, 1)
	}
}

// Close releases the handles. Per KILL_ON_JOB_CLOSE this also terminates any
// process still inside the Job - the last line of defence against an orphan.
func (j *childJob) Close() {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.exited != 0 {
		_ = windows.CloseHandle(j.exited)
		j.exited = 0
	}
	if j.handle != 0 {
		_ = windows.CloseHandle(j.handle)
		j.handle = 0
	}
}
