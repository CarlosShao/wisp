//go:build windows

package tools

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/proc"
)

// app_launch_windows_test.go — AC row 4, live half: "app.launch launches
// (notepad fixture) under Job; AUMID + name + path paths work." This is
// the only ticket-22 leg that really starts a process, against the real
// C30 JobScope (KILL_ON_JOB_CLOSE), and it closes the Job (killing the
// tree) via t.Cleanup.

func requireNotepad(t *testing.T) string {
	t.Helper()
	for _, cand := range []string{
		filepath.Join(os.Getenv("SystemRoot"), "System32", "notepad.exe"),
		"notepad.exe",
	} {
		if p, err := exec.LookPath(cand); err == nil {
			if _, statErr := os.Stat(p); statErr == nil {
				return p
			}
		}
		// absolute candidate: stat directly
		if _, err := os.Stat(cand); err == nil {
			return cand
		}
	}
	t.Fatalf("notepad fixture not found (AC row 4 names it; this runner lacks it)")
	return ""
}

func Test22AppLaunchNotepadUnderRealJob(t *testing.T) {
	js, err := proc.OpenJobScope()
	if err != nil {
		t.Fatalf("OpenJobScope: %v", err)
	}
	t.Cleanup(func() { _ = js.Close() })

	w := &webShared{
		deps: WebDeps{Jobs: js},
		logf: func(string, ...any) {},
	}
	notepad := requireNotepad(t)

	// path arm
	res := execCall(t, appLaunch{w: w}, map[string]any{"app": notepad})
	if res.IsError {
		t.Fatalf("path launch: %s", res.Text)
	}
	if !strings.Contains(res.Text, "已纳入 Job Object") {
		t.Errorf("must report real Job membership: %s", res.Text)
	}
	pids, err := js.TreePIDs()
	if err != nil {
		t.Fatalf("TreePIDs: %v", err)
	}
	if len(pids) == 0 {
		t.Fatal("Job tree empty: notepad did not enter the Job")
	}

	// name arm (App Paths registry / PATH): "notepad"
	res = execCall(t, appLaunch{w: w}, map[string]any{"app": "notepad"})
	if res.IsError {
		t.Fatalf("name launch: %s", res.Text)
	}
	pids2, err := js.TreePIDs()
	if err != nil {
		t.Fatal(err)
	}
	if len(pids2) < 2 {
		t.Fatalf("second child missing from job tree: %v", pids2)
	}

	// C30's kill-on-close is the whole point: closing the Job must take
	// the tree with it. Close, then watch the last pid die.
	last := pids2[len(pids2)-1]
	if err := js.Close(); err != nil {
		t.Fatalf("Close job: %v", err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, last)
		if err != nil {
			return // gone: kill-on-close worked
		}
		var code uint32
		if err := windows.GetExitCodeProcess(h, &code); err == nil && code != stillActive {
			_ = windows.CloseHandle(h)
			return // exited
		}
		_ = windows.CloseHandle(h)
		time.Sleep(100 * time.Millisecond)
	}
	t.Errorf("pid %d survived Job close (C30 kill-on-close not observed)", last)
}

func Test22AppLaunchAUMIDClassifiesToWindowsArm(t *testing.T) {
	// classification + carrier-shape only — no real UWP activation here
	// (the activated app process is not Job-containable; see
	// launcher_windows.go header; a live Calculator launch would be a
	// side effect this suite must not own).
	if !windowsAUMIDPattern.MatchString("Microsoft.WindowsCalculator_8wekyb3d8bbwe!App") {
		t.Fatal("AUMID pattern must accept the canonical shape")
	}
	for _, bad := range []string{"", "!", "a!b!c", `x!\..y`, "sp ace!App", "microsoft.windows!app ", "shell:AppsFolder\\X!Y"} {
		if windowsAUMIDPattern.MatchString(bad) {
			t.Errorf("AUMID pattern accepted %q", bad)
		}
	}
	w := &webShared{deps: WebDeps{Jobs: nil}, logf: func(string, ...any) {}}
	target, err := w.resolveApp(context.Background(), "Microsoft.WindowsCalculator_8wekyb3d8bbwe!App")
	if err != nil || target.Kind != "aumid" {
		t.Fatalf("classification: %+v err=%v", target, err)
	}
}
