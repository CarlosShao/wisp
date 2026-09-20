//go:build windows

package main

// OS-level argv forensics for `wisp secret` (ticket 63, Key constraint:
// 明文不进 argv - 否则 Get-CimInstance Win32_Process 能看到命令行).
//
// Asserting on the []string the parent built would only prove our own call
// site is careful. This test therefore runs the REAL binary as a child, and
// while it is alive reads its command line out of the PEB - the same
// _RTL_USER_PROCESS_PARAMETERS.CommandLine field WMI/Get-CimInstance reports.
// The secret has already been handed to the child over its standard input, so
// it is resident in the child when the command line is captured: the assertion
// is made while a leak would be visible.
//
// TestProcessCommandLineProbeDetectsAPlantedValue is the control that makes the
// rest of this file trustworthy: it plants the fake key in a child's argv and
// requires the same helper to find it.

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
	"unicode/utf16"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/secret"
	"golang.org/x/sys/windows"
)

// amd64 offsets of the fields this test reads. They are pinned to the 64-bit
// layout on purpose: the test fails loudly on a mismatch (the round-trip of
// the expected arguments is itself an assertion), never by skipping.
const (
	processVMRead = 0x0010
	// PROCESS_QUERY_LIMITED_INFORMATION: enough for NtQueryInformationProcess
	// with ProcessBasicInformation on a same-user child.
	processQueryLimitedInformation = 0x1000

	processBasicInformationClass = 0
	// _PROCESS_BASIC_INFORMATION.PebBaseAddress (handle-sized slot 1).
	pbiPebSlot = 1
	// PEB.ProcessParameters.
	pebProcessParametersOffset = 0x20
	// _RTL_USER_PROCESS_PARAMETERS.CommandLine (a UNICODE_STRING).
	rtlUserProcessParametersCommandLineOffset = 0x70
	// UNICODE_STRING: USHORT Length at +0, PWSTR Buffer at +8 (amd64).
	unicodeStringLengthOffset = 0
	unicodeStringBufferOffset = 8
)

// procCommandLine reads a live process's command line the way WMI does. The
// error it returns on failure is worth keeping: "no command line captured" and
// "OpenProcess denied" are different diagnoses, and the caller must not have to
// guess which one it hit.
func procCommandLine(pid uint32) (string, error) {
	if unsafe.Sizeof(uintptr(0)) != 8 {
		return "", fmt.Errorf("this forensics helper is for 64-bit processes (uintptr is %d bytes)", unsafe.Sizeof(uintptr(0)))
	}
	h, err := windows.OpenProcess(processQueryLimitedInformation|processVMRead, false, pid)
	if err != nil {
		return "", fmt.Errorf("OpenProcess(%d): %w", pid, err)
	}
	defer func() { _ = windows.CloseHandle(h) }()

	nqip := windows.NewLazySystemDLL("ntdll.dll").NewProc("NtQueryInformationProcess")
	var pbi [6]uintptr
	var returned uint32
	r1, _, errno := nqip.Call(
		uintptr(h),
		processBasicInformationClass,
		uintptr(unsafe.Pointer(&pbi[0])),
		uintptr(unsafe.Sizeof(pbi)),
		uintptr(unsafe.Pointer(&returned)),
	)
	if r1 != 0 {
		return "", fmt.Errorf("NtQueryInformationProcess(ProcessBasicInformation): ntstatus 0x%x (%v)", r1, errno)
	}
	peb := pbi[pbiPebSlot]
	if peb == 0 {
		return "", errors.New("PEB address read as 0")
	}
	params, err := rpmUintptr(h, peb+pebProcessParametersOffset)
	if err != nil {
		return "", err
	}
	cmdAddr := params + rtlUserProcessParametersCommandLineOffset
	length, err := rpmUint16(h, cmdAddr+unicodeStringLengthOffset)
	if err != nil {
		return "", err
	}
	if length == 0 || length > 32768 {
		return "", fmt.Errorf("implausible UNICODE_STRING length %d at 0x%x", length, cmdAddr)
	}
	buffer, err := rpmUintptr(h, cmdAddr+unicodeStringBufferOffset)
	if err != nil {
		return "", err
	}
	utf16Bytes, err := rpmBytes(h, buffer, uintptr(length))
	if err != nil {
		return "", err
	}
	units := make([]uint16, length/2)
	for i := range units {
		units[i] = uint16(utf16Bytes[2*i]) | uint16(utf16Bytes[2*i+1])<<8
	}
	return string(utf16.Decode(units)), nil
}

func rpmUintptr(h windows.Handle, addr uintptr) (uintptr, error) {
	var v uintptr
	n, err := rpmBytes(h, addr, unsafe.Sizeof(v))
	if err != nil {
		return 0, err
	}
	return uintptr(bytesToUint64LE(n)), nil
}

func rpmUint16(h windows.Handle, addr uintptr) (uint16, error) {
	n, err := rpmBytes(h, addr, 2)
	if err != nil {
		return 0, err
	}
	return uint16(n[0]) | uint16(n[1])<<8, nil
}

func rpmBytes(h windows.Handle, addr uintptr, size uintptr) ([]byte, error) {
	if addr == 0 {
		return nil, fmt.Errorf("null address read from the target process")
	}
	buf := make([]byte, size)
	var read uintptr
	if err := windows.ReadProcessMemory(h, addr, &buf[0], size, &read); err != nil {
		return nil, fmt.Errorf("ReadProcessMemory(0x%x, %d): %w", addr, size, err)
	}
	if read != size {
		return nil, fmt.Errorf("ReadProcessMemory(0x%x): read %d of %d bytes", addr, read, size)
	}
	return buf, nil
}

func bytesToUint64LE(b []byte) uint64 {
	var v uint64
	for i := len(b) - 1; i >= 0; i-- {
		v = v<<8 | uint64(b[i])
	}
	return v
}

// buildWispForTest compiles the real binary (the ticket's claim is about the
// shipped command line, not a test harness's) and colocates the native DLLs
// the way scripts/build.ps1 does.
func buildWispForTest(t *testing.T) string {
	t.Helper()
	if runtime.GOARCH != "amd64" {
		t.Fatalf("GOARCH = %s, want amd64: the PEB offsets this test reads are 64-bit specific", runtime.GOARCH)
	}
	goExe := "go"
	if p, err := exec.LookPath("go"); err == nil {
		goExe = p
	} else if rooted := filepath.Join(runtime.GOROOT(), "bin", "go.exe"); fileExists(rooted) {
		goExe = rooted
	}
	dir := t.TempDir()
	exe := filepath.Join(dir, "wisp.exe")
	cmd := exec.Command(goExe, "build", "-o", exe, "./cmd/wisp")
	cmd.Dir = filepath.Join("..", "..")
	var buildLog bytes.Buffer
	cmd.Stdout, cmd.Stderr = &buildLog, &buildLog
	if err := cmd.Run(); err != nil {
		t.Fatalf("go build ./cmd/wisp failed (is mingw on PATH?): %v\n%s", err, buildLog.String())
	}
	if !fileExists(exe) {
		t.Fatalf("build produced no %s", exe)
	}
	dllDir := filepath.Join("..", "..", "third_party", "sherpa-onnx")
	dlls, err := filepath.Glob(filepath.Join(dllDir, "*.dll"))
	if err != nil || len(dlls) == 0 {
		t.Fatalf("no native DLLs in %s - run scripts/fetch-deps.ps1 first (glob err %v)", dllDir, err)
	}
	for _, d := range dlls {
		data, err := os.ReadFile(d)
		if err != nil {
			t.Fatalf("read %s: %v", d, err)
		}
		if err := os.WriteFile(filepath.Join(dir, filepath.Base(d)), data, 0o600); err != nil {
			t.Fatalf("place %s next to wisp.exe: %v", filepath.Base(d), err)
		}
	}
	return exe
}

func fileExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

// pollCommandLine reads a live process's command line with a bounded retry and
// returns the last error it hit. A fixed attempt budget, no wall-clock
// arithmetic (D42#9): before the PID is visible in the snapshot OpenProcess
// simply fails, and exhausting the budget is a hard failure at the call site,
// never a skip.
func pollCommandLine(pid uint32) (string, error) {
	var lastErr error
	for attempt := 0; attempt < 100; attempt++ {
		time.Sleep(50 * time.Millisecond)
		cl, err := procCommandLine(pid)
		if err == nil && cl != "" {
			return cl, nil
		}
		if err == nil {
			err = errors.New("empty command line")
		}
		lastErr = err
	}
	return "", lastErr
}

// runSetWithSecretOnStdin starts `wisp secret set <name> ...` with the fake key
// already written to the child's stdin but the pipe still open (so the child is
// alive, has the value in hand, and has not finished). The callback receives the
// OS-visible command line; closing stdin lets the run finish.
func runSetWithSecretOnStdin(t *testing.T, exe, dataDir, name string, extraArgs []string, secretValue string, inspect func(t *testing.T, cmdline string, args []string)) {
	t.Helper()
	args := append([]string{"secret", "set", name}, extraArgs...)
	cmd := exec.Command(exe, args...)
	cmd.Dir = filepath.Dir(exe)
	cmd.Env = append(os.Environ(), "WISP_ENV=test", "WISP_TEST_DATA_DIR="+dataDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("StdinPipe: %v", err)
	}
	if _, err := io.WriteString(stdin, secretValue+"\n"); err != nil {
		t.Fatalf("write stdin: %v", err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("start %s: %v", exe, err)
	}
	pid := uint32(cmd.Process.Pid)
	// The child stays blocked in its stdin read until the pipe closes, so the
	// capture cannot race it away: the command line is read while the plaintext
	// is already resident in the process.
	cmdline, err := pollCommandLine(pid)
	if err != nil {
		_ = stdin.Close()
		_ = cmd.Wait()
		t.Fatalf("could not read the child's command line from its PEB within the 5s budget: %v; stderr=%s", err, stderr.String())
	}
	inspect(t, cmdline, args)
	if err := stdin.Close(); err != nil {
		t.Fatalf("close stdin: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		t.Logf("child exited with %v; stdout=%s stderr=%s", err, stdout.String(), stderr.String())
	}
}

// TestSecretArgvCarriesNoSecret is the Get-CimInstance-level proof for the
// --from-stdin path.
func TestSecretArgvCarriesNoSecret(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	const name = "argvprobe"
	// Middle fragment included on purpose: a truncated/quoted form of the key
	// in argv still fails the test.
	middleFragment := fakeKey[6 : len(fakeKey)-6]

	t.Run("from-stdin", func(t *testing.T) {
		ran := false
		runSetWithSecretOnStdin(t, exe, dataDir, name, []string{"--from-stdin"}, fakeKey, func(t *testing.T, cmdline string, args []string) {
			ran = true
			// Sanity: this really is the child's command line, read from the
			// OS - the same field Get-CimInstance Win32_Process prints.
			for _, want := range []string{"wisp.exe", "secret", "set", name, "--from-stdin"} {
				if !strings.Contains(cmdline, want) {
					t.Fatalf("PEB command line %s does not even contain %q: the capture is not reading the child we started",
						diagnose(cmdline), want)
				}
			}
			// Past the program path the OS-visible argv must be exactly the
			// tokens that were asked for: anything appended (a self
			// re-invocation carrying the value, a wrapper) fails here too, even
			// if it happens to dodge the probe string.
			tail := cmdline[strings.Index(cmdline, "secret"):]
			if got, want := strings.TrimSpace(tail), strings.Join(args, " "); got != want {
				t.Errorf("the child's OS-visible argv is %q, want exactly %q", got, want)
			}
			if strings.Contains(cmdline, fakeKey) {
				t.Errorf("the child's OS-visible command line carries the plaintext secret: %s", diagnose(cmdline))
			}
			if strings.Contains(cmdline, middleFragment) {
				t.Errorf("the child's command line carries a fragment of the secret: %s", diagnose(cmdline))
			}
		})
		if !ran {
			t.Fatal("the inspection callback never ran")
		}
		// The run really stored it: without this, a command line free of the
		// secret would only prove the child died early.
		if !fileExists(filepath.Join(dataDir, "secrets", name)) {
			t.Fatal("the child never stored the blob, so the argv observation was vacuous")
		}
		st, err := secret.NewStore(dataDir)
		if err != nil {
			t.Fatalf("open store: %v", err)
		}
		got, err := st.Resolve("dpapi:" + name)
		if err != nil || got != fakeKey {
			t.Errorf("resolve after the child run: (%d chars, %v), want the value entered on stdin", len(got), err)
		}
	})

	// The interactive form, run with a piped stdin (no console): it must
	// refuse without storing anything. Its command line is the fixed argv
	// built here - `wisp secret set <name>` and nothing else - and the
	// no-value-slot property of that argv is pinned for the whole file by
	// TestSecretFlagsAreBoolOnly, because a console-less child exits before an
	// in-flight PEB read is possible.
	t.Run("interactive-without-console", func(t *testing.T) {
		other := t.TempDir()
		args := []string{"secret", "set", name}
		if strings.Contains(strings.Join(args, " "), fakeKey) {
			t.Fatal("test setup error: the argv built for the interactive path carries the secret")
		}
		cmd := exec.Command(exe, args...)
		cmd.Dir = filepath.Dir(exe)
		cmd.Env = append(os.Environ(), "WISP_ENV=test", "WISP_TEST_DATA_DIR="+other)
		var stdout, stderr bytes.Buffer
		cmd.Stdout, cmd.Stderr = &stdout, &stderr
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatalf("StdinPipe: %v", err)
		}
		if err := cmd.Start(); err != nil {
			t.Fatalf("start: %v", err)
		}
		if _, err := io.WriteString(stdin, fakeKey+"\n"); err != nil {
			t.Fatalf("write stdin: %v", err)
		}
		_ = stdin.Close()
		if err := cmd.Wait(); err == nil {
			t.Fatalf("the shipped binary read a piped stdin as console input; stdout=%s", stdout.String())
		}
		if strings.Contains(stdout.String(), fakeKey) || strings.Contains(stderr.String(), fakeKey) {
			t.Error("the console-less refusal echoed the piped value back")
		}
		if fileExists(filepath.Join(other, "secrets", name)) {
			t.Error("a piped stdin must not be treated as console input: the blob was stored")
		}
	})
}

// TestSecretRealBinaryRefusesValueFlagAtOSLevel drives the shipped exe with the
// one argv shape that would leak: a value-carrying flag. It must be refused,
// must store nothing, and must not repeat the value back.
func TestSecretRealBinaryRefusesValueFlag(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	cmd := exec.Command(exe, "secret", "set", "never", "--value="+fakeKey)
	cmd.Dir = filepath.Dir(exe)
	cmd.Env = append(os.Environ(), "WISP_ENV=test", "WISP_TEST_DATA_DIR="+dataDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if err == nil {
		t.Fatal("a value-carrying flag was accepted by the shipped binary")
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("run failed in an unexpected way: %v", err)
	}
	if exitErr.ExitCode() != 2 {
		t.Errorf("exit code = %d, want the usage code 2", exitErr.ExitCode())
	}
	if strings.Contains(stdout.String(), fakeKey) || strings.Contains(stderr.String(), fakeKey) {
		t.Error("the rejected flag's value was echoed back by the shipped binary")
	}
	st, serr := secret.NewStore(dataDir)
	if serr != nil {
		t.Fatalf("open store: %v", serr)
	}
	if blobs, err := st.Blobs(); err != nil || len(blobs) != 0 {
		t.Errorf("a rejected argv value produced %d blob(s), want 0 (err %v)", len(blobs), err)
	}
}

// diagnose reports only the shape of a command line that failed an assertion,
// so a test log cannot itself become a leak of what was being looked for.
func diagnose(cmdline string) string {
	return fmt.Sprintf("<%d chars, args=%q>", len(cmdline), strings.Fields(cmdline))
}

// ---------------------------------------------------------------------------
// Positive control for the forensics helper
// ---------------------------------------------------------------------------

// The child mode is gated on two independent markers (an environment variable
// and a positional argument) so a stray environment cannot make an ordinary
// `go test ./cmd/wisp/` run block on a stdin read.
const (
	argvProbeChildEnv   = "WISP_ARGV_PROBE_CHILD"
	argvProbeChildToken = "wisp-argv-probe-child"
	argvProbeHelperName = "TestProcessCommandLineProbeHelperProcess"
)

// TestProcessCommandLineProbeHelperProcess is not a test: it is the body of the
// child process the positive control below starts. Under an ordinary run it
// returns immediately; under the two child markers it blocks on stdin, keeping
// the values its parent planted in its argv resident while they are read.
func TestProcessCommandLineProbeHelperProcess(t *testing.T) {
	if os.Getenv(argvProbeChildEnv) != "1" || !containsString(os.Args, argvProbeChildToken) {
		return
	}
	if _, err := io.ReadAll(os.Stdin); err != nil {
		t.Errorf("probe child: read stdin: %v", err)
	}
}

func containsString(hay []string, needle string) bool {
	for _, h := range hay {
		if h == needle {
			return true
		}
	}
	return false
}

// TestProcessCommandLineProbeDetectsAPlantedValue is what makes every "argv
// carries no secret" assertion in this file evidence rather than absence of
// evidence: it starts a real child process with the fake key *as an argument*
// and requires the very same procCommandLine helper to report it. A helper that
// silently returned a truncated or reconstructed string would keep the leak
// tests green forever; this one cannot pass unless the probe reads argv the way
// Get-CimInstance does.
func TestProcessCommandLineProbeDetectsAPlantedValue(t *testing.T) {
	middleFragment := fakeKey[6 : len(fakeKey)-6]
	// Both planted tokens are positional: a Go test binary leaves positionals
	// alone, so the child still starts, reaches the helper above, and waits
	// there with the values in its command line.
	cmd := exec.Command(os.Args[0],
		"-test.run=^"+argvProbeHelperName+"$", "-test.timeout=120s",
		argvProbeChildToken, fakeKey, middleFragment)
	cmd.Env = append(os.Environ(), argvProbeChildEnv+"=1")
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("StdinPipe: %v", err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("start the probe child: %v", err)
	}
	cmdline, perr := pollCommandLine(uint32(cmd.Process.Pid))
	_ = stdin.Close()
	if perr != nil {
		_ = cmd.Wait()
		t.Fatalf("POSITIVE CONTROL FAILED: the probe could not read even a child whose argv deliberately carries the secret: %v (stderr=%s)",
			perr, stderr.String())
	}
	if err := cmd.Wait(); err != nil {
		t.Fatalf("the probe child exited with %v (stdout=%s stderr=%s)", err, stdout.String(), stderr.String())
	}
	if !strings.Contains(cmdline, fakeKey) {
		t.Errorf("POSITIVE CONTROL FAILED: the probe did not report the plaintext planted in the child's argv, so the no-secret assertions in this file are vacuous. Captured: %s",
			diagnose(cmdline))
	}
	if !strings.Contains(cmdline, middleFragment) {
		t.Errorf("POSITIVE CONTROL FAILED: the probe missed the middle fragment of the planted argv value: %s",
			diagnose(cmdline))
	}
}
