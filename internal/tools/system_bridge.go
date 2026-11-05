package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"sync"
)

// ---------------------------------------------------------------------------
// The host bridge seam for the machine-control family (ticket 23)
//
// D34's system.* / window.* / clipboard.* / media.control / screen.capture /
// input.type rows all reach the operating system rather than the filesystem, and
// the ticket's hard constraint is "All system calls through the host bridge"
// (SPEC-07 §3, the same wording the D34 note carries). So nothing in this
// family calls a Win32 entry point directly: every action arrives here as an
// Action, and exactly one implementation - SystemBridge - answers.
//
// Why an interface rather than a package-level proc table:
//
//   - The risk verdict is decided BEFORE the action runs (the bridge assesses,
//     routes, and only then calls Tool.Execute), so the thing being executed and
//     the thing that was judged must be the same value. An Action is that value:
//     it is what the Decl.Facts hook read (the R8 irreversibility classes) and
//     what Execute hands over, byte for byte.
//   - Every acceptance criterion in ticket 23 that is about a REFUSAL - the UIPI
//     answer for an elevated window, the typed confirmation word for a power
//     action - is only testable if the refusal can be provoked without an
//     elevated window and without a shutdown. The seam is where a test injects
//     "this target is elevated" and "record the shutdown, do not run it".
//
// The production implementation is split by platform: system_backend_windows.go
// speaks to the OS (with a per-surface fake in this file selected by
// system_ops_windows.go only when told to), and system_ops_other.go answers
// every action with ErrNoSystemBackend. A build for another platform therefore
// REFUSES (fail-closed, IsError text the model can act on) instead of pretending
// the machine was touched.
// ---------------------------------------------------------------------------

// Action names one machine-control verb. The vocabulary mirrors the D34 rows
// exactly (system.get, system.set, system.power, window.*, clipboard.*,
// media.control, screen.capture, input.type); a verb outside it is not an
// action and is refused rather than guessed at.
type Action string

// The frozen action vocabulary of ticket 23.
const (
	// system.get / system.set
	ActionSystemGet Action = "system.get"
	ActionSetVolume Action = "system.set.volume"
	ActionSetMute   Action = "system.set.mute"
	ActionSetBright Action = "system.set.brightness"
	ActionMediaKey  Action = "system.set.media_key"
	ActionSetIME    Action = "system.set.ime"
	// system.power
	ActionPowerLock      Action = "system.power.lock"
	ActionPowerSleep     Action = "system.power.sleep"
	ActionPowerRestart   Action = "system.power.restart"
	ActionPowerShutdown  Action = "system.power.shutdown"
	ActionPowerLogoff    Action = "system.power.logoff"
	actionPowerUnknown   Action = "system.power.unknown"
	actionMediaPlayPause Action = "media.control.play_pause"
	actionMediaNext      Action = "media.control.next"
	actionMediaPrev      Action = "media.control.previous"
	actionMediaStop      Action = "media.control.stop"
	// window.*
	ActionWindowList    Action = "window.list"
	ActionWindowFocus   Action = "window.focus"
	ActionWindowMin     Action = "window.manage.minimize"
	ActionWindowMax     Action = "window.manage.maximize"
	ActionWindowRestore Action = "window.manage.restore"
	ActionWindowMove    Action = "window.manage.move"
	ActionWindowResize  Action = "window.manage.resize"
	ActionWindowSnap    Action = "window.manage.snap"
	ActionWindowClose   Action = "window.manage.close"
	// clipboard.*
	ActionClipboardRead  Action = "clipboard.read"
	ActionClipboardWrite Action = "clipboard.write"
	// screen.capture
	ActionScreenCapture Action = "screen.capture"
	// input.type
	ActionInputType Action = "input.type"
)

// systemFuncs is the whole surface the family needs from an OS backend. A
// backend returns (payload, err): payload is whatever Facts it collected (an
// ActionResult for most verbs, a richer report for the read verbs), err is a
// refusal or a platform limit. Nothing here retries, sleeps or owns a goroutine:
// a backend must honor ctx and return (the bridge already holds the C22
// deadline on that ctx).
type systemFuncs interface {
	// act is the single entry point. It is given the action plus the argument
	// map the tool already parsed, so a new verb is a new case in the backend
	// rather than a new method on a host-facing interface.
	act(ctx context.Context, a Action, args map[string]any) (ActionResult, error)
}

// SystemBridge is the host bridge: the SystemBridge interface the tools are
// constructed with, and the ONE implementation production uses.
//
// Zero value: no backend wired for any surface, so every action answers
// ErrNoSystemBackend. That is the fail-closed shape, and it is why the family's
// constructors accept a nil SystemBridge without nil-panicking later.
type SystemBridge struct {
	mu sync.RWMutex
	// Backend is the OS implementation. nil = selected by the platform build
	// (system_ops_windows.go on Windows, system_ops_other.go elsewhere).
	Backend systemFuncs
	// Ops is the observable side for tests and for the audit sink: every action
	// passes through it even when Backend is nil, so "was this refused before
	// the OS was touched" is answerable.
	Ops *systemOps
}

// systemFuncsOrNil returns the backend to use, honoring the platform default.
func (b *SystemBridge) funcs() (systemFuncs, error) {
	if b == nil {
		return nil, ErrNoSystemBackend
	}
	b.mu.RLock()
	backend := b.Backend
	b.mu.RUnlock()
	if backend != nil {
		return backend, nil
	}
	return platformSystemBackend()
}

// Do runs one action. It is the only method the tools call.
func (b *SystemBridge) Do(ctx context.Context, a Action, args map[string]any) (ActionResult, error) {
	if err := ctx.Err(); err != nil {
		return ActionResult{}, err
	}
	if b != nil && b.Ops != nil {
		if err := b.Ops.guard(a, args); err != nil {
			return ActionResult{}, err
		}
	}
	f, err := b.funcs()
	if err != nil {
		return ActionResult{}, err
	}
	res, err := f.act(ctx, a, args)
	if b != nil && b.Ops != nil {
		b.Ops.note(a, args, res, err)
	}
	return res, err
}

// ErrNoSystemBackend is the honest answer of a build (or a composition) with no
// OS backend. It is a REFUSAL, not a fault: the tools turn it into IsError text
// so the model learns the machine cannot be reached from here, and the task
// continues (D37).
var ErrNoSystemBackend = errors.New(
	"tools: 本平台/本组合没有接入系统能力后端（票 23 的 host bridge 未接线），已拒绝该动作")

// ErrElevatedTarget is the UIPI answer, and it is a separate value rather than
// a formatted string because ticket 23's AC#3 asserts the exact wording that
// reaches the model: PLAN.md:3234 requires 「目标程序以管理员身份运行，Wisp 无法
// 向其输入」 and forbids every attempt at an elevation bypass.
var ErrElevatedTarget = errors.New(
	"目标程序以管理员身份运行，Wisp 无法向其输入（Windows UIPI 平台硬限制：普通权限进程不能" +
		"向更高完整性级别的窗口发送输入，本工具不会尝试提权绕过）。请把目标应用以普通权限运行后重试")

// ErrNotImplementedHere names an action the current platform backend cannot
// perform. It is distinct from ErrNoSystemBackend (no bridge at all) so the
// model is told which side is missing.
var ErrNotImplementedHere = errors.New("tools: 该动作在本平台后端未实现")

// ActionResult is what one action reports back. It is data, never prose, so a
// tool can render it in the shape its row promises (the read verbs render a
// report, the write verbs render one line).
type ActionResult struct {
	// Detail is the human-readable one-liner for a write verb ("音量已设为 40%").
	Detail string
	// Report is the structured payload of a read verb (system.get, window.list,
	// clipboard.read, screen.capture). Keys are stable strings so a test can
	// assert them instead of matching prose.
	Report map[string]any
	// Bytes is the memory-only payload of screen.capture. It NEVER leaves this
	// process: see screen_capture.go for the four places bytes_ref must not go.
	Bytes []byte
	// Mime is Bytes' media type.
	Mime string
	// DRM marks a capture that came back black because the window is
	// protected-content (PLAN.md:3236, #4 of §16.11).
	DRM bool
	// Target names what was acted on (a window title, a process name, a
	// clipboard origin) for the C25 origin and the audit trail.
	Target string
}

// ---------------------------------------------------------------------------
// systemOps: the observable side (audit + test spy)
// ---------------------------------------------------------------------------

// systemOps records what the family asked the OS for. It is deliberately
// bounded (a fixed ring of the last N lines) because an unbounded audit buffer
// inside a 24/7 background process is a leak with a friendly name.
type systemOps struct {
	mu sync.Mutex
	// Hooks lets a test/host answer "this target is elevated" or "this window is
	// DRM-protected" without touching a real machine. Returning an error from
	// Guard refuses the action before the backend sees it.
	Guard func(a Action, args map[string]any) error
	// Sink receives one line per action. nil discards.
	Sink func(format string, args ...any)
	// lines keeps the last maxOpsLines actions for assertions.
	lines []string
}

// maxOpsLines bounds systemOps.lines.
const maxOpsLines = 64

func (o *systemOps) guard(a Action, args map[string]any) error {
	o.mu.Lock()
	hook := o.Guard
	o.mu.Unlock()
	if hook == nil {
		return nil
	}
	return hook(a, args)
}

func (o *systemOps) note(a Action, args map[string]any, res ActionResult, err error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	line := fmt.Sprintf("action=%s target=%q err=%v", a, res.Target, err)
	o.lines = append(o.lines, line)
	if len(o.lines) > maxOpsLines {
		o.lines = o.lines[len(o.lines)-maxOpsLines:]
	}
	if o.Sink != nil {
		// The argument VALUES are deliberately absent from this line: a
		// clipboard payload or a typed text would otherwise land in the audit
		// sink, which is exactly the leak screen.capture's row forbids for its
		// own bytes and which no other verb promises.
		o.Sink("tools: system %s keys=%d target=%q err=%v",
			a, len(args), res.Target, err)
	}
}

// Lines returns a copy of the recorded action lines (test/audit surface).
func (o *systemOps) Lines() []string {
	if o == nil {
		return nil
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	return append([]string(nil), o.lines...)
}

// ---------------------------------------------------------------------------
// small helpers shared by the family
// ---------------------------------------------------------------------------

// isWindowsGate is the one place a non-test file asks "am I on Windows". It is
// used only for prose (which limit applies), never for a security decision: a
// security decision lives in a build-tagged backend, where the compiler is the
// arbiter.
func isWindowsGate() bool { return runtime.GOOS == "windows" }

// argString reads one string argument.
func argString(args map[string]any, key string) string {
	v, ok := args[key]
	if !ok || v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

// argInt reads one integer argument. json.Unmarshal into map[string]any yields
// float64, so both spellings are accepted and a non-numeric value is a zero.
func argInt(args map[string]any, key string) (int, bool) {
	v, ok := args[key]
	if !ok || v == nil {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	case json.Number:
		f, err := n.Float64()
		if err != nil {
			return 0, false
		}
		return int(f), true
	}
	return 0, false
}

// argBool reads one boolean argument.
func argBool(args map[string]any, key string) bool {
	v, ok := args[key]
	if !ok || v == nil {
		return false
	}
	b, ok := v.(bool)
	return ok && b
}

// systemErrText renders a backend refusal for the model. The distinction it
// keeps is the one D37 draws: an unreachable platform is a tool-class failure
// the model can act on, never a host fault.
func systemErrText(what string, err error) string {
	switch {
	case errors.Is(err, ErrElevatedTarget):
		return err.Error()
	case errors.Is(err, ErrNoSystemBackend), errors.Is(err, ErrNotImplementedHere):
		return what + "：" + err.Error()
	default:
		return what + "失败：" + err.Error()
	}
}
