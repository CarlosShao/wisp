//go:build windows && winlive

package main

// Ticket 43, live tier: the shipped no-args process really receives the two
// lifecycle events and really dies the two deaths the ticket names.
//
//   - TestLive43QueryEndSessionExitsWithin5s: a WM_QUERYENDSESSION is delivered
//     to the child's own lifecycle window (SendMessageTimeoutW addressed to the
//     window whose owning pid is the child's - the message pump, the wndproc,
//     the fast path and the D38(e) walk are all real; the only simulated hop is
//     the one that delivers the message, which on a real shutdown is the OS).
//     The clock is monotonic (D42#9).
//   - TestLive43KillNineMidTaskMarksInterruptedOnNextBoot: the child is really
//     force-killed (TerminateProcess, the Windows shape of kill -9: no defers,
//     no cleanup), its run marker must survive, and the next boot must mark the
//     dead process's still-'running' task row as interrupted and never rerun
//     anything.
//
// Run it alone (the same rule the 228 live tier books):
//
//	go test -tags winlive ./cmd/wisp -run TestLive43 -v

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/memory"
	"golang.org/x/sys/windows"
)

var (
	live43User32          = windows.NewLazySystemDLL("user32.dll")
	pLive43FindWindowExW  = live43User32.NewProc("FindWindowExW")
	pLive43SendMessageTO  = live43User32.NewProc("SendMessageTimeoutW")
	pLive43GetWindowPID   = live43User32.NewProc("GetWindowThreadProcessId")
	live43SMTOAbortIfHung = uintptr(2)
)

// buildWispForTest43 is this file's build harness. It is buildWispForTest's
// shape (a real go build of ./cmd/wisp, native DLLs colocated next to the exe)
// with one difference the staging rule forces: the exe lands NEXT TO THIS TEST
// BINARY instead of in a fresh temp dir, because the shared harness resolves
// the repo root from its own working directory and this binary runs from the
// staging root, where System32's older onnxruntime.dll would otherwise be found
// first (28-r1's §6 measurement). WISP_43_REPO_ROOT names the checkout, the
// same way WISP_15E_TMP names a staging directory.
func buildWispForTest43(t *testing.T) string {
	t.Helper()
	repo := os.Getenv("WISP_43_REPO_ROOT")
	if repo == "" {
		t.Fatal("WISP_43_REPO_ROOT is not set: the staging runner must name the checkout")
	}
	testBin, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable: %v", err)
	}
	binDir := filepath.Dir(testBin)
	exe := filepath.Join(binDir, "wisp43-under-test.exe")
	goExe := "go"
	build := exec.Command(goExe, "build", "-o", exe, "./cmd/wisp")
	build.Dir = repo
	var buildLog strings.Builder
	build.Stdout, build.Stderr = &buildLog, &buildLog
	if err := build.Run(); err != nil {
		t.Fatalf("go build ./cmd/wisp failed: %v\n%s", err, buildLog.String())
	}
	dllDir := filepath.Join(repo, "third_party", "sherpa-onnx")
	dlls, err := filepath.Glob(filepath.Join(dllDir, "*.dll"))
	if err != nil || len(dlls) == 0 {
		t.Fatalf("no native DLLs in %s (glob err %v)", dllDir, err)
	}
	for _, d := range dlls {
		dst := filepath.Join(binDir, filepath.Base(d))
		if _, err := os.Stat(dst); err != nil {
			data, rerr := os.ReadFile(d)
			if rerr != nil {
				t.Fatalf("read %s: %v", d, rerr)
			}
			if werr := os.WriteFile(dst, data, 0o600); werr != nil {
				t.Fatalf("place %s next to the test binary: %v", filepath.Base(d), werr)
			}
		}
	}
	return exe
}

// findLifecycleWindow43 returns the child's lifecycle window (class match +
// owning pid match - the same two-fact identity rule the 228 live tier books).
func findLifecycleWindow43(t *testing.T, pid uint32) windows.HWND {
	t.Helper()
	cls, err := windows.UTF16PtrFromString(lifecycleWindowClass)
	if err != nil {
		t.Fatalf("UTF16PtrFromString: %v", err)
	}
	var hwnd windows.HWND
	for {
		r, _, _ := pLive43FindWindowExW.Call(0, uintptr(hwnd), uintptr(unsafe.Pointer(cls)), 0)
		if r == 0 {
			return 0
		}
		hwnd = windows.HWND(r)
		var owner uint32
		pLive43GetWindowPID.Call(r, uintptr(unsafe.Pointer(&owner)))
		if owner == pid {
			return hwnd
		}
	}
}

// logFiles43 reads every JSONL log file the child's sink wrote.
func logFiles43(t *testing.T, dataDir string) string {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(dataDir, logDirName))
	if err != nil {
		t.Fatalf("read log dir: %v", err)
	}
	var b strings.Builder
	for _, e := range entries {
		raw, err := os.ReadFile(filepath.Join(dataDir, logDirName, e.Name()))
		if err != nil {
			continue
		}
		b.Write(raw)
		b.WriteString("\n")
	}
	return b.String()
}

func TestLive43QueryEndSessionExitsWithin5s(t *testing.T) {
	exe := buildWispForTest43(t)
	dataDir := t.TempDir()
	leg := bootResidentLeg(t, exe, dataDir)

	var hwnd windows.HWND
	if !pollUntil127(200, func() bool {
		hwnd = findLifecycleWindow43(t, leg.pid())
		return hwnd != 0
	}) {
		t.Fatalf("the child never created its %s window.\n%s", lifecycleWindowClass, leg.console())
	}

	// The deliver: SendMessageTimeoutW to the child's window. Monotonic clock
	// around the whole hop (D42#9).
	start := time.Now()
	var res uintptr
	r, callErr, _ := pLive43SendMessageTO.Call(
		uintptr(hwnd), uintptr(wmQueryEndSession), 0, 0,
		live43SMTOAbortIfHung, uintptr(3000), uintptr(unsafe.Pointer(&res)))
	if r == 0 {
		t.Fatalf("SendMessageTimeoutW(WM_QUERYENDSESSION) failed or timed out: %v (hwnd=%v pid=%d).\n%s",
			callErr, hwnd, leg.pid(), leg.console())
	}
	if !leg.exitedWithin(100) {
		t.Fatalf("the child did not exit within 5s of WM_QUERYENDSESSION (elapsed %v).\n%s",
			time.Since(start), leg.console())
	}
	elapsed := time.Since(start)
	if elapsed > 5*time.Second {
		t.Fatalf("exit took %v, want <= 5s (D38 fast path)", elapsed)
	}
	if leg.exitStatus() != "exited 0" {
		t.Fatalf("exit status = %s, want a clean exit through the shutdown path", leg.exitStatus())
	}
	t.Logf("LIVE READING: WM_QUERYENDSESSION -> exit in %v", elapsed)

	// The order audit: the full 10-step trail the resident boot now lands on
	// its own sink, steps 4/5/9 present and in the frozen relative order. The
	// fast-path note itself is logged by the sequence DURING the walk (before
	// the deferred per-step trail), so its position is existence, not order.
	logs := logFiles43(t, dataDir)
	i1 := strings.Index(logs, `D38(e) step 1 (scheduler-close)`)
	i4 := strings.Index(logs, `D38(e) step 4 (stop-audio): ok`)
	i5 := strings.Index(logs, `D38(e) step 5 (release-speech-sessions)`)
	iFast := strings.Index(logs, "fast shutdown: skipping non-critical log/db flushes")
	i9 := strings.Index(logs, `D38(e) step 9 (close-job-object): ok`)
	i10 := strings.Index(logs, `D38(e) step 10 (exit)`)
	for _, at := range []struct {
		what string
		i    int
	}{
		{"step 1 (scheduler-close)", i1},
		{"step 4 (stop-audio) done", i4},
		{"step 5 (release-speech-sessions) done", i5},
		{"the fast-path skip of step 7", iFast},
		{"step 9 (close-job-object) done", i9},
		{"step 10 (exit)", i10},
	} {
		if at.i < 0 {
			t.Fatalf("the D38(e) trail is missing %q after the fast shutdown.\nlog tail:\n%s", at.what, tailOf43(logs, 4000))
		}
	}
	if !(i1 < i4 && i4 < i5 && i5 < i9 && i9 < i10) {
		t.Fatalf("the D38(e) trail is out of order: step1@%d step4@%d step5@%d step9@%d step10@%d", i1, i4, i5, i9, i10)
	}
}

// TestLive43KillNineMidTaskMarksInterruptedOnNextBoot is the crash-recovery
// cell, end to end: a still-'running' task row, a really force-killed process,
// a surviving marker, and the next boot that marks the row interrupted and
// reruns nothing.
func TestLive43KillNineMidTaskMarksInterruptedOnNextBoot(t *testing.T) {
	exe := buildWispForTest43(t)
	dataDir := t.TempDir()

	// The row the dying process left behind: opened, never finished. Planted
	// through the real store - the same row shape a killed task's StartTaskLog
	// leaves on disk.
	{
		st, err := memory.Open(dataDir)
		if err != nil {
			t.Fatalf("memory.Open: %v", err)
		}
		if err := st.StartTaskLog(t.Context(), memory.TaskLog{
			ID:    "task-43-killed",
			State: string(memory.TaskStateRunning), QueryText: "kill -9 时还在跑的任务",
		}); err != nil {
			t.Fatalf("StartTaskLog: %v", err)
		}
		if err := st.Close(); err != nil {
			t.Fatalf("store close: %v", err)
		}
	}

	leg1 := bootResidentLeg(t, exe, dataDir)
	if !pollUntil127(200, func() bool {
		return strings.Contains(leg1.stdout.String(), "resident event loop running")
	}) {
		leg1.stop()
		t.Fatalf("leg1 never reached its event loop.\n%s", leg1.console())
	}

	// kill -9, the Windows shape: TerminateProcess. No defers run, the marker
	// stays, the row stays open.
	if err := leg1.cmd.Process.Kill(); err != nil {
		t.Fatalf("force-kill: %v", err)
	}
	<-leg1.done
	if hasRunMarker(dataDir) == false {
		t.Fatal("the run marker is gone after a force kill: the next boot would miss the crash")
	}

	// Next boot: the sweep marks the row interrupted and says so on stdout.
	leg2 := bootResidentLeg(t, exe, dataDir)
	if !pollUntil127(200, func() bool {
		return strings.Contains(leg2.stdout.String(), "已标记为中断")
	}) {
		leg2.stop()
		t.Fatalf("the next boot never announced the interrupted marking.\n%s", leg2.console())
	}
	st, err := memory.Open(dataDir)
	if err != nil {
		t.Fatalf("memory.Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	tl, err := st.TaskLogByID(t.Context(), "task-43-killed")
	if err != nil {
		t.Fatalf("TaskLogByID: %v", err)
	}
	if tl.State != string(memory.TaskStateInterrupted) {
		t.Fatalf("task state after the next boot = %q, want interrupted", tl.State)
	}
	rows, err := st.ListTaskLogs(t.Context())
	if err != nil {
		t.Fatalf("ListTaskLogs: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("task rows = %d, want 1: the next boot opened no new task (no rerun)", len(rows))
	}

	// A CLEAN exit clears the marker (the boot's own defer), so the boot after
	// a clean exit must not sweep.
	leg2.breakToLoop()
	if !leg2.exitedWithin(200) {
		leg2.stop()
		t.Fatalf("leg2 did not leave through its own shutdown path.\n%s", leg2.console())
	}
	if hasRunMarker(dataDir) {
		t.Fatal("the run marker survived a clean exit")
	}
	n, err := recoverCrashedTasks(t.Context(), st, dataDir)
	if err != nil || n != 0 {
		t.Fatalf("sweep after a clean exit marked %d tasks (err=%v), want 0/nil", n, err)
	}
}

func tailOf43(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}
