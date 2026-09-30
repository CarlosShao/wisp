// Command twobreak is a 228-v1 verification probe (read-only with respect to the
// repository: it starts a built wisp.exe in a temp data dir and sends it console
// control events).
//
// Question it answers with a real machine reading, not a reasoning step:
// does a SECOND console control event, sent while the first one is still being
// answered (i.e. while rb.stop() and rt.Shutdown() run), kill the process at
// 0xc000013a and eat the D38(e) trail - or is it swallowed?
//
// cmd/wisp/resident_windows.go registers signal.Notify(bootExit, os.Interrupt,
// SIGTERM) on a 1-slot buffered channel that only ONE reader looks at once (the
// non-blocking select at :126-133). After that nobody drains it, so the second
// and every later signal goes to a full buffer, which os/signal handles by
// DROPPING it. That is a claim about Go's runtime semantics; this probe measures
// the consequence instead of asserting it.
//
// It is its own module on purpose (a go.mod in this directory), so the root
// module's ./... never compiles it and the shared tree's gates stay untouched.
//
// Usage (from the repository root, with the sherpa DLLs on PATH):
//
//	go build -o "$TEMP/twobreak.exe" .scratch/wisp/probes/228/v1/twobreak
//	"$TEMP/twobreak.exe" -exe "$TEMP/wispprobe/wisp.exe" -rounds 3
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var (
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procGenerateConsoleCtrlEvent = kernel32.NewProc("GenerateConsoleCtrlEvent")
)

const (
	ctrlBreakEvent         = 1
	createNewProcessGroup  = 0x00000200
	shutdownRecordFragment = "shutdown step"
	exitReportFragment     = "exited through the D38(e) shutdown order"
	ballCreatedFragment    = "ball: the resident leg created the floating ball window"
	ballStoppedFragment    = "ball: resident leg destroyed the ball window, its tray icon and its hotkeys"
)

func main() {
	exe := flag.String("exe", "", "path to a built wisp.exe (required)")
	rounds := flag.Int("rounds", 1, "how many independent child runs")
	gap := flag.Duration("gap", 3*time.Millisecond, "delay between the FIRST and the SECOND ctrl-break")
	fileDelay := flag.Duration("file-delay", 0, "extra wait after the sink file appears, before the first break")
	flag.Parse()
	if *exe == "" {
		fmt.Println("twobreak: -exe is required")
		os.Exit(2)
	}
	for r := 1; r <= *rounds; r++ {
		oneRun(r, *exe, *gap, *fileDelay)
	}
}

func oneRun(round int, exe string, gap, fileDelay time.Duration) {
	dataDir, err := os.MkdirTemp("", "twobreak-data-")
	if err != nil {
		fmt.Printf("round %d: mktemp: %v\n", round, err)
		return
	}
	defer os.RemoveAll(dataDir)

	logPath := filepath.Join(dataDir, "console.txt")
	consoleFile, err := os.Create(logPath)
	if err != nil {
		fmt.Printf("round %d: create console capture: %v\n", round, err)
		return
	}
	defer consoleFile.Close()

	cmd := exec.Command(exe)
	cmd.Dir = filepath.Dir(exe)
	cmd.Env = append(os.Environ(), "WISP_ENV=test", "WISP_TEST_DATA_DIR="+dataDir)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createNewProcessGroup}
	cmd.Stdout = consoleFile
	cmd.Stderr = consoleFile
	if err := cmd.Start(); err != nil {
		fmt.Printf("round %d: start: %v\n", round, err)
		return
	}
	pid := uint32(cmd.Process.Pid)
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	sinkDir := filepath.Join(dataDir, "logs")
	deadline := time.Now().Add(20 * time.Second)
	sawFile := false
	for time.Now().Before(deadline) {
		if n, _ := countJSONL(sinkDir); n > 0 {
			sawFile = true
			break
		}
		time.Sleep(time.Millisecond)
	}
	if !sawFile {
		fmt.Printf("round %d: no sink file appeared, aborting this round\n", round)
		_ = cmd.Process.Kill()
		<-done
		return
	}
	bootStamp := time.Now()
	if fileDelay > 0 {
		time.Sleep(fileDelay)
	}
	firstRC := sendBreak(pid)
	firstStamp := time.Now()
	time.Sleep(gap)
	secondRC := sendBreak(pid)
	secondStamp := time.Now()

	waitErr, exited := func() (error, bool) {
		select {
		case e := <-done:
			return e, true
		case <-time.After(30 * time.Second):
			_ = cmd.Process.Kill()
			return <-done, false
		}
	}()
	consoleFile.Close()

	body := readFileOnce(sinkDir)
	console := readFileOnceDir(dataDir, "console.txt")

	fmt.Printf("round %d  pid=%d\n", round, pid)
	fmt.Printf("  sink file at boot+%s\n", firstStamp.Sub(bootStamp).Round(time.Microsecond))
	fmt.Printf("  first break rc=%v  second break rc=%v (gap %s)\n", firstRC, secondRC, secondStamp.Sub(firstStamp).Round(time.Microsecond))
	fmt.Printf("  reaped by own exit=%v  waitErr=%v  exitStatus=%s\n", exited, waitErr, exitCodeOf(waitErr))
	fmt.Printf("  console has %q: %v\n", exitReportFragment, strings.Contains(console, exitReportFragment))
	fmt.Printf("  persistent log: shutdown-step lines=%d  ball-created=%v  ball-stopped=%v  bytes=%d\n",
		strings.Count(body, shutdownRecordFragment),
		strings.Contains(body, ballCreatedFragment),
		strings.Contains(body, ballStoppedFragment), len(body))
	fmt.Printf("  console tail: %s\n", lastLine(console))
	fmt.Printf("  CONSOLE DUMP >>>\n%s\n  <<< END CONSOLE\n", strings.TrimRight(console, "\n"))
}

func sendBreak(pid uint32) string {
	r1, _, e1 := procGenerateConsoleCtrlEvent.Call(ctrlBreakEvent, uintptr(pid))
	if r1 == 0 {
		return "FAILED " + e1.Error()
	}
	return "ok"
}

func countJSONL(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".jsonl") {
			n++
		}
	}
	return n, nil
}

func readFileOnce(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "(no log dir: " + err.Error() + ")"
	}
	var all []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		all = append(all, string(b))
	}
	return strings.Join(all, "\n")
}

func readFileOnceDir(dir, name string) string {
	b, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return "(no console capture: " + err.Error() + ")"
	}
	return string(b)
}

func lastLine(s string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) == 0 {
		return "(empty)"
	}
	return lines[len(lines)-1]
}

func exitCodeOf(err error) string {
	if err == nil {
		return "exit 0"
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return fmt.Sprintf("code %d (%v)", ee.ExitCode(), err)
	}
	return err.Error()
}
