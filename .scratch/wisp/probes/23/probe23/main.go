//go:build windows

// probe23 answers ticket 23 AC#3's unchecked cell on the real machine:
//
//	input.type: types 300-char text into fixture app focused window;
//	elevated-process injection -> explicit UIPI error; [...]
//
// It does not re-implement the product. Every injection below goes through the
// REAL roster the composition root builds (tools.BuiltinSystemEntries with a
// zero SystemDeps, whose nil Backend resolves to the platform backend through
// SystemBridge.funcs), i.e. the same tool objects, the same bridge, the same
// backend and the same SendInput path cmd/wisp registers. Only the measurements
// around them are probe code.
//
// Usage: probe23.exe <fixture-exe> [out-dir]
//
// Output: probe.log (append-only key=value lines) and summary.json in out-dir,
// plus the fixture's own status.txt/dump.txt/text.txt inside out-dir/fixture/.
package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/tools"
	"golang.org/x/sys/windows"
)

var (
	user32   = windows.NewLazySystemDLL("user32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procForegroundWindow    = user32.NewProc("GetForegroundWindow")
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	procGetGUIThreadInfo    = user32.NewProc("GetGUIThreadInfo")
	procGetWindowThreadProc = user32.NewProc("GetWindowThreadProcessId")
	procAttachThreadInput   = user32.NewProc("AttachThreadInput")
	procSetFocus            = user32.NewProc("SetFocus")
	procEnumWindows         = user32.NewProc("EnumWindows")
	procIsWindowVisible     = user32.NewProc("IsWindowVisible")
	procGetWindowTextLen    = user32.NewProc("GetWindowTextLengthW")
	procGetWindowTextW      = user32.NewProc("GetWindowTextW")
	procPostMessageW        = user32.NewProc("PostMessageW")
	procSendInput           = user32.NewProc("SendInput")
	procGetCurrentThreadID  = kernel32.NewProc("GetCurrentThreadId")
	procSwitchToThisWindow  = user32.NewProc("SwitchToThisWindow")
	procSetWindowPos        = user32.NewProc("SetWindowPos")
	procGetWindowRect       = user32.NewProc("GetWindowRect")
	procSetCursorPos        = user32.NewProc("SetCursorPos")
	procGetCursorPos        = user32.NewProc("GetCursorPos")
	procMouseEvent          = user32.NewProc("mouse_event")
	procOpenProcess         = kernel32.NewProc("OpenProcess")
	procQueryFullProcName   = kernel32.NewProc("QueryFullProcessImageNameW")
	procCloseHandle         = kernel32.NewProc("CloseHandle")
)

const (
	wmClose = 0x0010

	inputKeyboard    = 1
	keyeventfKeyup   = 0x0002
	keyeventfUnicode = 0x0004
	processQueryInfo = 0x1000
	threadQueryInfo  = 0x0040
	guithreadinfoLen = 0x48
)

// kbInput / inputUnion mirror the product's structs (system_backend_windows.go
// :130-146) so the independent raw-Instrument reading uses the same layout the
// product's SendInput call does.
type kbInput struct {
	Vk      uint16
	Scan    uint16
	Flags   uint32
	Time    uint32
	Extra   uintptr
	Padding [8]byte
}

type inputUnion struct {
	Type uint32
	_    uint32
	Data kbInput
}

type rect struct{ Left, Top, Right, Bottom int32 }

type guiThreadInfo struct {
	Size       uint32
	Flags      uint32
	Active     uintptr
	Focus      uintptr
	Capture    uintptr
	MenuOwner  uintptr
	MoveSize   uintptr
	Caret      uintptr
	RectLeft   int32
	RectTop    int32
	RectRight  int32
	RectBottom int32
}

type status struct {
	ts         string
	hwnd       uintptr
	edit       uintptr
	pid        uint32
	foreground string
	isFg       int
	focus      uintptr
}

var (
	outDir string
	logF   *os.File
	steps  = map[string]any{}
)

func logf(format string, a ...any) {
	line := fmt.Sprintf("ts=%s %s\n", time.Now().Format("2006-01-02T15:04:05.000"), fmt.Sprintf(format, a...))
	fmt.Print(line)
	if logF != nil {
		_, _ = logF.WriteString(line)
	}
}

func record(key string, v any) {
	steps[key] = v
	logf("%s=%v", key, v)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: probe23.exe <fixture-exe> [out-dir]")
		os.Exit(2)
	}
	fixtureExe := os.Args[1]
	if len(os.Args) > 2 {
		outDir = os.Args[2]
	} else {
		outDir = filepath.Join("logs", "run-"+time.Now().Format("20060102-150405"))
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "mkdir:", err)
		os.Exit(1)
	}
	fixtureDir := filepath.Join(outDir, "fixture")
	_ = os.MkdirAll(fixtureDir, 0o755)
	fixtureDirGlobal = fixtureDir
	f, err := os.OpenFile(filepath.Join(outDir, "probe.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "log:", err)
		os.Exit(1)
	}
	logF = f
	defer f.Close()

	logf("probe=probe23 pid=%d exe=%s", os.Getpid(), os.Args[0])
	logf("out_dir=%s", abs(outDir))
	record("injector_integrity", fmt.Sprintf("0x%x", selfIntegrity()))

	entries := tools.BuiltinSystemEntries(tools.SystemDeps{})
	record("roster_entries", len(entries))

	// ---- start the fixture window -------------------------------------
	cmd := exec.Command(fixtureExe, fixtureDir, "180")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		logf("FATAL fixture start: %v", err)
		os.Exit(1)
	}
	fixturePID := uint32(cmd.Process.Pid)
	defer func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	}()
	record("fixture_pid", fixturePID)

	st, err := waitStatus(fixtureDir, 20*time.Second)
	if err != nil {
		logf("FATAL fixture status: %v", err)
		os.Exit(1)
	}
	record("fixture_hwnd", fmt.Sprintf("0x%x (%d)", st.hwnd, st.hwnd))
	record("fixture_edit_hwnd", fmt.Sprintf("0x%x (%d)", st.edit, st.edit))
	record("fixture_foreground_method", st.foreground)
	record("fixture_is_foreground_at_start", st.isFg)
	record("fixture_process_name_probe_side", procName(st.pid))
	record("fixture_integrity", fmt.Sprintf("0x%x", pidIntegrity(st.pid)))

	// ---- the window must really hold foreground AND the keyboard focus --
	fgMethod, fg, focus, err := ensureFocus(st.hwnd, st.edit)
	record("focus_method", fgMethod)
	record("foreground_hwnd", fmt.Sprintf("0x%x", fg))
	record("keyboard_focus_hwnd", fmt.Sprintf("0x%x", focus))
	record("foreground_is_fixture", fg == st.hwnd)
	record("focus_is_edit", focus == st.edit)
	if err != nil || fg != st.hwnd || focus != st.edit {
		record("ABORTED", "focus could not be placed on the fixture edit control; nothing was typed")
		// PROBE23-UIPI-SPLIT (orchestrator 2026-11-14): the two halves of AC#3
		// have different premises. Typing 300 chars needs the fixture to own the
		// foreground; the elevated-target half does NOT - the refusal happens at
		// the integrity gate, before a single character is typed, so it needs no
		// focus at all. Measured on this desktop: an elevated window (Taskmgr,
		// High IL) held the foreground and a Medium-IL process cannot take it
		// back, which aborted the run and left the UIPI reading unmeasured even
		// though it was measurable. Run that half before exiting.
		record("fixture_half_skipped", "foreground held by a higher-integrity window")
		elevatedForm(entries, st, fixtureDir)
		_ = cmd.Process.Kill()
		writeSummary()
		os.Exit(3)
	}

	// Product-side proof of the same fact: system.get reports the focused
	// window title/process the product itself resolves.
	if res, ferr := callTool(entries, "system.get", map[string]any{}); ferr == nil {
		record("system_get_text", res.Text)
	}
	if res, ferr := callTool(entries, "window.focus", map[string]any{"hwnd": int(st.hwnd)}); ferr == nil {
		record("window_focus_text", res.Text)
		record("window_focus_is_error", res.IsError)
	} else {
		record("window_focus_host_error", ferr.Error())
	}

	// ---- control: the raw instrument can deliver one char at all --------
	clearFixture(fixtureDir)
	if n, sendErr := rawSendUnicode('Q'); sendErr != nil {
		record("raw_control_sendinput", fmt.Sprintf("n=%d err=%v", n, sendErr))
	} else {
		record("raw_control_sendinput", fmt.Sprintf("n=%d err=nil", n))
	}
	d := waitDump(fixtureDir, 3*time.Second, func(d dumpState) bool { return d.chars >= 1 })
	record("raw_control_chars_in_window", d.chars)
	record("raw_control_wm_char", d.wmChar)

	// ---- run A: 300 ASCII chars through the product tool ----------------
	textA := twoHundredNinetyNineEdge300()
	record("runA_expected_chars", len([]rune(textA)))
	record("runA_expected_utf16_units", len(utf16.Encode([]rune(textA))))
	record("runA_expected_sha256", sha(textA))
	clearFixture(fixtureDir)
	incBefore := waitDump(fixtureDir, 2*time.Second, func(d dumpState) bool { return d.chars == 0 })
	record("runA_chars_after_clear", incBefore.chars)
	runInjection(entries, "runA", textA, st)

	// ---- run B: 300 CJK chars (BMP, one UTF-16 unit each) ---------------
	textB := cjk300()
	record("runB_expected_chars", len([]rune(textB)))
	record("runB_expected_utf16_units", len(utf16.Encode([]rune(textB))))
	record("runB_expected_sha256", sha(textB))
	clearFixture(fixtureDir)
	waitDump(fixtureDir, 2*time.Second, func(d dumpState) bool { return d.chars == 0 })
	runInjection(entries, "runB", textB, st)

	// ---- elevated-process form -----------------------------------------
	elevatedForm(entries, st, fixtureDir)

	// ---- cleanup: close only the fixture we started ---------------------
	procPostMessageW.Call(st.hwnd, wmClose, 0, 0)
	time.Sleep(300 * time.Millisecond)
	record("fixture_close_posted", true)
	writeSummary()
}

type injectionResult struct {
	dump        dumpState
	textMatches bool
	elapsedMS   int64
}

func runInjection(entries []tools.Entry, run, text string, st status) injectionResult {
	t0 := time.Now()
	res, err := callTool(entries, "input.type", map[string]any{
		"text": text,
		"hwnd": int(st.hwnd),
	})
	elapsed := time.Since(t0).Milliseconds()
	if err != nil {
		record(run+"_host_error", err.Error())
	}
	record(run+"_is_error", res.IsError)
	record(run+"_text", res.Text)
	record(run+"_elapsed_ms", elapsed)

	want := len([]rune(text))
	d := waitDump(fixtureDirGlobal, 5*time.Second, func(d dumpState) bool { return d.chars >= want })
	got, rerr := os.ReadFile(filepath.Join(fixtureDirGlobal, "text.txt"))
	if rerr != nil {
		logf("%s readback error: %v", run, rerr)
	}
	record(run+"_readback_chars", d.chars)
	record(run+"_readback_runes", d.runes)
	record(run+"_readback_utf16_units", d.units)
	record(run+"_readback_sha256", d.sha)
	record(run+"_readback_wm_char", d.wmChar)
	record(run+"_readback_wm_keydown", d.wmKey)
	matches := string(got) == text
	_ = os.WriteFile(filepath.Join(outDir, run+"-readback.txt"), got, 0o644)
	record(run+"_readback_equals_expected", matches)
	record(run+"_summary", fmt.Sprintf("is_error=%v chars=%d runes=%d wm_char=%d matches=%v elapsed_ms=%d", res.IsError, d.chars, d.runes, d.wmChar, matches, elapsed))
	if !matches {
		record(run+"_readback_head", head(string(got), 80))
	}
	return injectionResult{dump: d, textMatches: matches, elapsedMS: elapsed}
}

var fixtureDirGlobal string

func ensureFocus(hwnd, focus uintptr) (string, uintptr, uintptr, error) {
	method := "already-focused"
	if foreground() != hwnd {
		if ok, _, _ := procSetForegroundWindow.Call(hwnd); ok != 0 && foreground() == hwnd {
			method = "SetForegroundWindow"
		} else if forceForegroundViaAttach(hwnd) {
			method = "AttachThreadInput+SetForegroundWindow"
		} else if clickActivate(hwnd) {
			method = "raise+click"
		} else {
			method = "foreground-failed"
		}
	}
	if foreground() == hwnd && focusHWND() != focus {
		if attachSetFocus(hwnd, focus) {
			method += "+AttachThreadInput+SetFocus"
		} else {
			method += "+set-focus-failed"
		}
	}
	fg, got := foreground(), focusHWND()
	if fg == hwnd && got == focus {
		return method, fg, got, nil
	}
	return method + "/unverified", fg, got, fmt.Errorf("could not move foreground and keyboard focus onto the fixture")
}

// clickActivate raises the target window and clicks its middle, which is the
// one foreground transition Windows always honours for a same-integrity window:
// a real user input event. The cursor is put back by the caller.
func clickActivate(hwnd uintptr) bool {
	const (
		hwndTopmost   = ^uintptr(0)
		hwndNoTopmost = ^uintptr(1)
		swpNoSize     = 0x0001
		swpNoMove     = 0x0002
		swpNoActivate = 0x0010
		mouseLeftDown = 0x0002
		mouseLeftUp   = 0x0004
	)
	var r rect
	procSetWindowPos.Call(hwnd, hwndTopmost, 0, 0, 0, 0, swpNoMove|swpNoSize|swpNoActivate)
	time.Sleep(150 * time.Millisecond)
	if ok, _, _ := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&r))); ok == 0 {
		return false
	}
	x := int32((r.Left + r.Right) / 2)
	y := int32((r.Top + r.Bottom) / 2)
	var saved struct{ X, Y int32 }
	_, _, _ = procGetCursorPos.Call(uintptr(unsafe.Pointer(&saved)))
	procSetCursorPos.Call(uintptr(x), uintptr(y))
	time.Sleep(80 * time.Millisecond)
	procMouseEvent.Call(mouseLeftDown, 0, 0, 0, 0)
	procMouseEvent.Call(mouseLeftUp, 0, 0, 0, 0)
	time.Sleep(250 * time.Millisecond)
	procSetWindowPos.Call(hwnd, hwndNoTopmost, 0, 0, 0, 0, swpNoMove|swpNoSize|swpNoActivate)
	_, _, _ = procSetCursorPos.Call(uintptr(saved.X), uintptr(saved.Y))
	return foreground() == hwnd
}

// forceForegroundViaAttach is the standard toolkit fallback for the Windows
// foreground lock: attach to the thread that currently owns the foreground,
// then SetForegroundWindow, then detach.
func forceForegroundViaAttach(hwnd uintptr) bool {
	// The shell path first: SwitchToThisWindow ignores the foreground lock.
	procSwitchToThisWindow.Call(hwnd, 1)
	if foreground() == hwnd {
		return true
	}
	fg := foreground()
	var pid uint32
	fgThread, _, _ := procGetWindowThreadProc.Call(fg, uintptr(unsafe.Pointer(&pid)))
	myThread, _, _ := procGetCurrentThreadID.Call()
	if fgThread == 0 || fgThread == myThread {
		return false
	}
	procAttachThreadInput.Call(myThread, fgThread, 1)
	procSetForegroundWindow.Call(hwnd)
	procAttachThreadInput.Call(myThread, fgThread, 0)
	return foreground() == hwnd
}

func elevatedForm(entries []tools.Entry, st status, fixtureDir string) {
	fixtureDirGlobal = fixtureDir
	self := selfIntegrity()
	type win struct {
		hwnd  uintptr
		pid   uint32
		name  string
		title string
		il    uint32
	}
	var higher []win
	var all []win
	enumTopWindows(func(h uintptr) bool {
		var pid uint32
		procGetWindowThreadProc.Call(h, uintptr(unsafe.Pointer(&pid)))
		if pid == 0 {
			return true
		}
		il := pidIntegrity(pid)
		w := win{hwnd: h, pid: pid, name: procName(pid), title: windowTitle(h), il: il}
		all = append(all, w)
		if il > self {
			higher = append(higher, w)
		}
		return true
	})
	record("top_windows_scanned", len(all))
	record("windows_above_our_integrity", len(higher))
	for i, w := range higher {
		if i >= 8 {
			break
		}
		record(fmt.Sprintf("elevated_candidate_%d", i),
			fmt.Sprintf("hwnd=0x%x pid=%d name=%s il=0x%x title=%q", w.hwnd, w.pid, w.name, w.il, w.title))
	}
	if len(higher) == 0 {
		record("elevated_form", "no window above our integrity level exists on this desktop; the UIPI form is unmeasurable without starting one")
		return
	}
	target := higher[0]
	for _, w := range higher {
		if strings.EqualFold(w.name, "taskmgr.exe") {
			target = w
			break
		}
	}
	record("elevated_target",
		fmt.Sprintf("hwnd=0x%x pid=%d name=%s il=0x%x title=%q", target.hwnd, target.pid, target.name, target.il, target.title))

	before := dumpOf(fixtureDir)
	res, err := callTool(entries, "input.type", map[string]any{"text": "wisp-uipi-probe", "hwnd": int(target.hwnd)})
	if err != nil {
		record("elevated_host_error", err.Error())
	}
	record("elevated_is_error", res.IsError)
	record("elevated_text", res.Text)

	time.Sleep(400 * time.Millisecond)
	after := dumpOf(fixtureDir)
	record("elevated_fixture_chars_unchanged", before.chars == after.chars)

	// Independent platform reading: put the elevated window in the foreground
	// and ask SendInput for one character. A medium process is not allowed to
	// reach a high-integrity window, so the kernel answers "0 inserted" with
	// ERROR_ACCESS_DENIED instead of silently dropping it.
	origFg := foreground()
	record("pre_elevated_foreground", fmt.Sprintf("0x%x", origFg))
	if _, ferr := callTool(entries, "window.focus", map[string]any{"hwnd": int(target.hwnd)}); ferr == nil {
		logf("window.focus on the elevated target returned without host error")
	} else {
		record("elevated_focus_host_error", ferr.Error())
	}
	fgNow := foreground()
	record("elevated_foreground_after_focus", fmt.Sprintf("0x%x", fgNow))
	if fgNow == target.hwnd {
		n, sendErr := rawSendUnicode('Z')
		if sendErr != nil {
			record("elevated_raw_sendinput", fmt.Sprintf("n=%d err=%v", n, sendErr))
		} else {
			record("elevated_raw_sendinput", fmt.Sprintf("n=%d err=nil", n))
		}
	} else {
		record("elevated_raw_sendinput",
			"not measured: the elevated window could not be made foreground, so SendInput would have targeted another window")
	}
	// Put the fixture back in front so the next readback is unambiguous, then
	// restore whatever was in front when this leg began.
	_, _, _, _ = ensureFocus(st.hwnd, st.edit)
	if origFg != 0 && origFg != st.hwnd {
		procSetForegroundWindow.Call(origFg)
	}
}

// ---------------------------------------------------------------------------
// product surface
// ---------------------------------------------------------------------------

func callTool(entries []tools.Entry, name string, args map[string]any) (tools.Result, error) {
	raw, err := json.Marshal(args)
	if err != nil {
		return tools.Result{}, err
	}
	for _, e := range entries {
		if e.Tool.Name() != name {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return e.Tool.Execute(ctx, raw, nil)
	}
	return tools.Result{}, fmt.Errorf("tool %q is not in the roster", name)
}

// ---------------------------------------------------------------------------
// independent Win32 instruments
// ---------------------------------------------------------------------------

func foreground() uintptr {
	h, _, _ := procForegroundWindow.Call()
	return h
}

func selfIntegrity() uint32 {
	return pidIntegrity(uint32(os.Getpid()))
}

func pidIntegrity(pid uint32) uint32 {
	h, err := windows.OpenProcess(processQueryInfo, false, pid)
	if err != nil {
		return 0
	}
	defer windows.CloseHandle(h)
	var tok windows.Token
	if err := windows.OpenProcessToken(h, windows.TOKEN_QUERY, &tok); err != nil {
		return 0
	}
	defer tok.Close()
	var n uint32
	_ = windows.GetTokenInformation(tok, windows.TokenIntegrityLevel, nil, 0, &n)
	if n == 0 {
		return 0
	}
	buf := make([]byte, n)
	if err := windows.GetTokenInformation(tok, windows.TokenIntegrityLevel, &buf[0], n, &n); err != nil {
		return 0
	}
	// TOKEN_MANDATORY_LABEL is one SID_AND_ATTRIBUTES: the SID is
	// S-1-16-<integrity rid> (8192 medium, 12288 high, 16384 system).
	label := (*windows.SIDAndAttributes)(unsafe.Pointer(&buf[0]))
	if label.Sid == nil {
		return 0
	}
	s := label.Sid.String()
	i := strings.LastIndex(s, "-")
	if i < 0 {
		return 0
	}
	v, _ := strconv.ParseUint(s[i+1:], 10, 32)
	return uint32(v)
}

func procName(pid uint32) string {
	h, err := windows.OpenProcess(processQueryInfo, false, pid)
	if err != nil {
		return ""
	}
	defer windows.CloseHandle(h)
	var buf [windows.MAX_PATH + 1]uint16
	size := uint32(len(buf))
	r, _, _ := procQueryFullProcName.Call(uintptr(h), 0, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)))
	if r == 0 {
		return ""
	}
	full := windows.UTF16ToString(buf[:size])
	if i := strings.LastIndexAny(full, `\/`); i >= 0 {
		full = full[i+1:]
	}
	return strings.ToLower(full)
}

func windowTitle(h uintptr) string {
	n, _, _ := procGetWindowTextLen.Call(h)
	if int32(n) <= 0 {
		return ""
	}
	buf := make([]uint16, int(n)+1)
	procGetWindowTextW.Call(h, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return windows.UTF16ToString(buf)
}

func enumTopWindows(cb func(uintptr) bool) {
	cbPtr := windows.NewCallback(func(h, _ uintptr) uintptr {
		if v, _, _ := procIsWindowVisible.Call(h); v == 0 {
			return 1
		}
		if windowTitle(h) == "" {
			return 1
		}
		if cb(h) {
			return 1
		}
		return 0
	})
	procEnumWindows.Call(cbPtr, 0)
}

func focusHWND() uintptr {
	fg := foreground()
	if fg == 0 {
		return 0
	}
	var pid uint32
	thread, _, _ := procGetWindowThreadProc.Call(fg, uintptr(unsafe.Pointer(&pid)))
	if thread == 0 {
		return 0
	}
	info := guiThreadInfo{Size: guithreadinfoLen}
	r, _, _ := procGetGUIThreadInfo.Call(thread, uintptr(unsafe.Pointer(&info)))
	if r == 0 {
		return 0
	}
	return info.Focus
}

// attachSetFocus borrows the target thread's input state long enough to place
// keyboard focus on the edit control (a cross-process SetFocus is otherwise a
// silent no-op), then detaches.
func attachSetFocus(hwnd, focus uintptr) bool {
	var pid uint32
	thread, _, _ := procGetWindowThreadProc.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	myThread, _, _ := procGetCurrentThreadID.Call()
	if thread == 0 {
		return false
	}
	procAttachThreadInput.Call(myThread, thread, 1)
	defer procAttachThreadInput.Call(myThread, thread, 0)
	procSetFocus.Call(focus)
	return focusHWND() == focus
}

func rawSendUnicode(r rune) (int, error) {
	units := utf16.Encode([]rune{r})
	ins := make([]inputUnion, 0, len(units)*2)
	for _, u := range units {
		ins = append(ins,
			inputUnion{Type: inputKeyboard, Data: kbInput{Scan: u, Flags: keyeventfUnicode}},
			inputUnion{Type: inputKeyboard, Data: kbInput{Scan: u, Flags: keyeventfUnicode | keyeventfKeyup}})
	}
	n, _, e := procSendInput.Call(uintptr(len(ins)), uintptr(unsafe.Pointer(&ins[0])), unsafe.Sizeof(ins[0]))
	if int(n) != len(ins) {
		return int(n), fmt.Errorf("SendInput inserted %d/%d (win32=%v)", int(n), len(ins), e)
	}
	return int(n), nil
}

// ---------------------------------------------------------------------------
// fixture readback
// ---------------------------------------------------------------------------

type dumpState struct {
	chars  int
	runes  int
	units  int
	sha    string
	wmChar int
	wmKey  int
}

func dumpOf(dir string) dumpState {
	b, err := os.ReadFile(filepath.Join(dir, "dump.txt"))
	if err != nil {
		return dumpState{}
	}
	return parseDump(string(b))
}

func parseDump(line string) dumpState {
	var d dumpState
	for _, f := range strings.Fields(line) {
		k, v, ok := strings.Cut(f, "=")
		if !ok {
			continue
		}
		n, _ := strconv.Atoi(v)
		switch k {
		case "chars":
			d.chars = n
		case "runes":
			d.runes = n
		case "utf16_units":
			d.units = n
		case "sha256":
			d.sha = v
		case "wm_char":
			d.wmChar = n
		case "wm_keydown":
			d.wmKey = n
		}
	}
	return d
}

func waitDump(dir string, timeout time.Duration, ready func(dumpState) bool) dumpState {
	deadline := time.Now().Add(timeout)
	var d dumpState
	for {
		d = dumpOf(dir)
		if ready(d) || time.Now().After(deadline) {
			return d
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func clearFixture(dir string) {
	_ = os.WriteFile(filepath.Join(dir, "cmd.txt"), []byte("clear\n"), 0o644)
	time.Sleep(150 * time.Millisecond)
}

func waitStatus(dir string, timeout time.Duration) (status, error) {
	deadline := time.Now().Add(timeout)
	for {
		b, err := os.ReadFile(filepath.Join(dir, "status.txt"))
		if err == nil {
			var st status
			for _, f := range strings.Fields(string(b)) {
				k, v, ok := strings.Cut(f, "=")
				if !ok {
					continue
				}
				switch k {
				case "ts":
					st.ts = v
				case "hwnd":
					n, _ := strconv.ParseUint(v, 10, 64)
					st.hwnd = uintptr(n)
				case "edit":
					n, _ := strconv.ParseUint(v, 10, 64)
					st.edit = uintptr(n)
				case "pid":
					n, _ := strconv.ParseUint(v, 10, 64)
					st.pid = uint32(n)
				case "foreground_method":
					st.foreground = v
				case "is_foreground":
					n, _ := strconv.Atoi(v)
					st.isFg = n
				case "focus_hwnd":
					n, _ := strconv.ParseUint(v, 10, 64)
					st.focus = uintptr(n)
				}
			}
			if st.hwnd != 0 && st.foreground != "" {
				return st, nil
			}
		}
		if time.Now().After(deadline) {
			return status{}, fmt.Errorf("status.txt did not report a live window in %s", timeout)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// ---------------------------------------------------------------------------
// texts and small helpers
// ---------------------------------------------------------------------------

// twoHundredNinetyNineEdge300 is exactly 300 characters of pure ASCII: 18 full
// repetitions of a 16-char block plus a 12-char tail.
func twoHundredNinetyNineEdge300() string {
	return strings.Repeat("0123456789abcdef", 18) + "0123456789AB"
}

// cjk300 is exactly 300 BMP CJK characters, one UTF-16 unit each, so a run that
// types 300 of them also proves the off-by-one traps (bytes vs runes vs units)
// are not being reported as a match.
func cjk300() string {
	var b strings.Builder
	for i := 0; i < 300; i++ {
		b.WriteRune(rune(0x4E00 + i))
	}
	return b.String()
}

func sha(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func head(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

func abs(p string) string {
	a, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return a
}

func writeSummary() {
	b, err := json.MarshalIndent(steps, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(outDir, "summary.json"), b, 0o644)
}
