//go:build windows

package main

// Ticket 43, power & lifecycle events: the window and the notification channel
// through which the OS reaches this process, and the three sequences the ticket
// freezes on top of them.
//
// HOW THE OS REACHES THIS PROCESS, two independent doors because neither one
// alone reaches every shape:
//
//   - WM_QUERYENDSESSION / WM_ENDSESSION arrive only through a top-level
//     window, so this leg owns an INVISIBLE one (WS_OVERLAPPED, never shown).
//     A message-only window (HWND_MESSAGE) would not receive broadcasts; the
//     ball and panel windows are theirs, not this leg's, so this leg brings
//     its own.
//   - WM_POWERBROADCAST for PBT_APMSUSPEND is, since Vista, sent only to the
//     foreground application - a background resident would never see it. The
//     supplementary door is PowerRegisterSuspendResumeNotification
//     (DEVICE_NOTIFY_CALLBACK), which delivers suspend/resume regardless of
//     focus. Both doors funnel into ONE event channel drained by ONE thread,
//     so the sequences below can never run concurrently with each other.
//
// WHAT RUNS WHERE. The events channel is drained by this leg's own thread
// ("lifecycle-sta", the same deliberate-outside-the-roster shape panel-sta
// books - observe's six frozen resident names are full, and adding to them is a
// human-approved contract change, not a wiring detail). The suspend sequence
// itself is the part that has to be fast and ordered:
//
//	1. pending L2 cards -> voided + auto-rejected + counted notice (the
//	   §16.11#10 lean the ticket lands: a stale approval hours later is a
//	   safety hole, and D43 row 24's refusal path is the legal row)
//	2. audio capture stops (device off before sleep - D40 row 1)
//	3. the session scope is disposed ON THE VOICE OWNER (closing a sherpa
//	   session under a decoding thread is the crash the owner-marshalling
//	   exists to prevent), waited for on a bounded monotonic budget
//
// Resume reopens the capture device and nothing else: the disposed scope
// reloads its engines on the next real wake (that IS D43 #31's dispose
// semantics - nobody re-arms models nobody asked for), the state machine has no
// suspend edge in the frozen table so it is left exactly where it was, and the
// wall-derived schedules the ticket mentions have no device in this build yet
// (named residual, not a silent skip).

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/proc"
	"golang.org/x/sys/windows"
)

// The frozen Windows constants this file needs. Kept local: internal/proc and
// internal/ball own their own copies, and this leg adds no dependency edge.
const (
	wmQueryEndSession    = 0x0011
	wmEndSession         = 0x0016
	wmPowerBroadcast     = 0x0218
	pbtApmSuspend        = 4
	pbtApmResumeAuto     = 7
	pbtApmResumeSuspend  = 5
	deviceNotifyCallback = 2
)

// suspendVoiceBudget is the monotonic ceiling the suspend sequence waits for
// the voice owner's dispose. 28-r1 measured the cascade at ~68 ms hot; the
// budget is generous so a cold-family swap can still land inside it, and a
// timeout is a recorded problem, never a hang: the OS is going to sleep either
// way.
const suspendVoiceBudget = 1500 * time.Millisecond

// lifecycleSink is what the assembly root wires. Three function values, the
// same injection shape RegisterShutdownHook and stopRequestFunc use: the module
// that owns each capability hands its body in, and this file knows nothing
// about approvals, audio or speech engines.
type lifecycleSink struct {
	suspend    func(why string)
	resume     func(why string)
	sessionEnd func(why string)
}

// lifecycleEvents is this leg's hold on the two doors. Nil-safe methods, like
// every other resident handle.
type lifecycleEvents struct {
	sink     lifecycleSink
	fastExit *atomic.Bool

	// events is the funnel both doors post into; the pump thread drains it and
	// runs the sink, so sequences serialize and a callback thread never blocks
	// (the OS notification callback must return quickly).
	events chan lifecycleEvent

	// hwnd is the invisible top-level window; powerReg is the notification
	// registration. cb keeps the NewCallback token alive for the process
	// lifetime.
	hwnd     uintptr
	powerReg uintptr
	cb       uintptr

	// suspendDepth counts the suspend sequences that actually ran, and
	// resumeDepth the resume ones. They are the reading the evidence table
	// quotes for "the events reached this process".
	suspendDepth atomic.Int64
	resumeDepth  atomic.Int64
}

type lifecycleEvent struct {
	kind string
	why  string
}

// lifecycleWindowClass is this window's registration class. A literal, named,
// because the live tier finds the child's window by it.
const lifecycleWindowClass = "WispLifecycleWindow"

var (
	lifecycleUser32             = syscall.NewLazyDLL("user32.dll")
	pLifecycleRegisterClassExW  = lifecycleUser32.NewProc("RegisterClassExW")
	pLifecycleCreateWindowExW   = lifecycleUser32.NewProc("CreateWindowExW")
	pLifecycleDefWindowProcW    = lifecycleUser32.NewProc("DefWindowProcW")
	pLifecyclePeekMessageW      = lifecycleUser32.NewProc("PeekMessageW")
	pLifecycleTranslateMessageW = lifecycleUser32.NewProc("TranslateMessage")
	pLifecycleDispatchMessageW  = lifecycleUser32.NewProc("DispatchMessageW")
	lifecycleKernel32           = syscall.NewLazyDLL("kernel32.dll")
	pLifecycleGetModuleHandleW  = lifecycleKernel32.NewProc("GetModuleHandleW")
	// PowerRegisterSuspendResumeNotification is declared in winuser.h but
	// implemented in powrprof.dll - this machine's user32 export table has no
	// such entry point, and a LazyProc that guesses wrong panics at first use.
	lifecyclePowrprof       = syscall.NewLazyDLL("powrprof.dll")
	pLifecyclePowerRegister = lifecyclePowrprof.NewProc("PowerRegisterSuspendResumeNotification")
)

// lifecycleSelf is the ONE lifecycle handle of this process. The resident boot
// builds exactly one; the registration-time wndproc reads it back here rather
// than through GWLP_USERDATA, so no Go pointer is ever handed to Win32 storage
// (a uintptr round-trip of a Go pointer is the shape vet refuses, and a package
// map would be a second truth about whose window this is).
var lifecycleSelf atomic.Pointer[lifecycleEvents]

// wndClassEx mirrors the Win32 WNDCLASSEX layout (amd64 packing).
type wndClassEx struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   windows.Handle
	Icon       windows.Handle
	Cursor     windows.Handle
	Background windows.Handle
	MenuName   *uint16
	ClassName  *uint16
	IconSm     windows.Handle
}

// msg mirrors the Win32 MSG layout (amd64 packing).
type msg struct {
	Hwnd    windows.HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

type windowsHandle = windows.Handle

// startLifecycleEvents registers the power callback and starts the pump on
// this leg's own thread. The window is created ON THAT THREAD (see pump): a
// window's messages queue up on the thread that created it, so a window built
// on the boot thread would have a queue nobody pumps - SendMessageTimeoutW
// would time out on every QUERYENDSESSION. It never stops the boot: a machine
// where either door fails still gets the other one, and the receipt says which.
func startLifecycleEvents(rt *proc.Runtime, sink lifecycleSink, fastExit *atomic.Bool) *lifecycleEvents {
	le := &lifecycleEvents{
		sink:     sink,
		fastExit: fastExit,
		events:   make(chan lifecycleEvent, 8),
	}
	le.registerPowerNotifications()
	// The pump thread. observe.Default.Spawn is the only sanctioned goroutine
	// entry (D22 ban #1); the name is deliberately outside the six-name frozen
	// resident roster, the same shape panel-sta books, and the roster warning
	// that costs is the same one the panel leg accepted in 33.
	observe.Default.Spawn("lifecycle-sta", "lifecycle", nil, func(ctx context.Context) {
		le.pump(ctx)
	})
	fmt.Printf("wisp: 生命周期腿已启动：WM_QUERYENDSESSION 窗口由 lifecycle-sta 线程装配；电源通知=%s\n",
		map[bool]string{true: "已注册", false: "未注册"}[le.powerReg != 0])
	return le
}

// createWindow registers the class and creates the invisible window. False
// means "the door is not there" - recorded by the caller, never fatal.
func (le *lifecycleEvents) createWindow() bool {
	clsName, err := syscall.UTF16PtrFromString(lifecycleWindowClass)
	if err != nil {
		return false
	}
	mod, _, _ := pLifecycleGetModuleHandleW.Call(0)
	if mod == 0 {
		return false
	}
	wc := wndClassEx{
		Size:      uint32(unsafe.Sizeof(wndClassEx{})),
		WndProc:   syscall.NewCallback(lifecycleWndProcTrampoline),
		Instance:  windowsHandle(mod),
		ClassName: clsName,
	}
	if r, _, _ := pLifecycleRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); r == 0 {
		return false
	}
	// WS_OVERLAPPED, never WS_VISIBLE: the window exists to receive broadcasts,
	// not to be seen.
	h, _, err := pLifecycleCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(clsName)),
		uintptr(unsafe.Pointer(clsName)),
		0, // style: WS_OVERLAPPED (0)
		0, 0, 0, 0,
		0, 0, uintptr(mod), 0)
	if h == 0 {
		slog.Warn("lifecycle: 窗口创建失败", "err", err.Error())
		return false
	}
	le.hwnd = h
	lifecycleSelf.Store(le)
	return true
}

// lifecycleWndProcTrampoline is the registration-time wndproc. It reads the one
// process handle from lifecycleSelf (see that var's comment for why not
// GWLP_USERDATA).
func lifecycleWndProcTrampoline(h uintptr, m uint32, w, l uintptr) uintptr {
	le := lifecycleSelf.Load()
	if le == nil {
		r, _, _ := pLifecycleDefWindowProcW.Call(h, uintptr(m), w, l)
		return r
	}
	return le.wndProc(h, m, w, l)
}

// wndProc answers the session-end question and funnels power events. The
// session-end answer returns TRUE (1) IMMEDIATELY - the system is waiting on
// this window's reply with its own timeout, and the shutdown sequence runs
// after the reply, not inside it.
func (le *lifecycleEvents) wndProc(h uintptr, m uint32, w, l uintptr) uintptr {
	switch m {
	case wmQueryEndSession, wmEndSession:
		le.post("session-end", "WM_QUERYENDSESSION/WM_ENDSESSION")
		return 1 // TRUE: this process can be ended; the sequence is already starting
	case wmPowerBroadcast:
		switch w {
		case pbtApmSuspend:
			le.post("suspend", "PBT_APMSUSPEND")
		case pbtApmResumeAuto, pbtApmResumeSuspend:
			le.post("resume", "PBT_APMRESUME")
		}
		return 1 // TRUE: processed (DefWindowProc's BROADCAST query result)
	case wmDestroy:
		le.hwnd = 0
	}
	r, _, _ := pLifecycleDefWindowProcW.Call(h, uintptr(m), w, l)
	return r
}

// post funnels one event to the pump. Non-blocking: a callback thread that
// cannot deliver drops the event with a log line rather than stalling the OS's
// power thread.
func (le *lifecycleEvents) post(kind, why string) {
	if le == nil {
		return
	}
	select {
	case le.events <- lifecycleEvent{kind: kind, why: why}:
	default:
		slog.Warn("lifecycle: 事件通道满，事件被丢弃", "kind", kind, "why", why)
	}
}

// deliverPowerEvent is the seam the notification callback and the tests both
// drive. Production: the PowerRegisterSuspendResumeNotification callback calls
// it; tests call it directly to drive the SAME body (a re-typed copy of the
// sequence in a test would read green whichever way production changed).
func (le *lifecycleEvents) deliverPowerEvent(typ uint32) {
	switch typ {
	case pbtApmSuspend:
		le.post("suspend", "DeviceNotifyCallback(PBT_APMSUSPEND)")
	case pbtApmResumeAuto, pbtApmResumeSuspend:
		le.post("resume", "DeviceNotifyCallback(PBT_APMRESUME)")
	}
}

// registerPowerNotifications subscribes to suspend/resume regardless of
// foreground focus. Failure is logged, not fatal: the WM_POWERBROADCAST door
// still exists, and a machine without either door boots fine.
func (le *lifecycleEvents) registerPowerNotifications() {
	// DEVICE_NOTIFY_SUBSCRIBE_PARAMETERS { Callback; Context }.
	le.cb = syscall.NewCallback(func(ctx, typ, setting uintptr) uintptr {
		le.deliverPowerEvent(uint32(typ))
		return 0
	})
	// Find() first: a missing entry point (an old Windows) is a logged
	// degradation, never the panic a bare LazyProc.Call would raise on the
	// boot thread.
	if err := pLifecyclePowerRegister.Find(); err != nil {
		slog.Warn("lifecycle: PowerRegisterSuspendResumeNotification 不存在于本机（后台挂起通知不可用）",
			"err", err.Error())
		return
	}
	params := struct {
		Callback uintptr
		Context  uintptr
	}{Callback: le.cb}
	var reg uintptr
	// The return is a DWORD error code, not a BOOL: 0 (ERROR_SUCCESS) is the
	// SUCCESS shape here, and the handle lands in reg.
	r, _, err := pLifecyclePowerRegister.Call(
		uintptr(deviceNotifyCallback),
		uintptr(unsafe.Pointer(&params)),
		uintptr(unsafe.Pointer(&reg)))
	if r == 0 {
		le.powerReg = reg
		return
	}
	slog.Warn("lifecycle: PowerRegisterSuspendResumeNotification 未注册（后台挂起通知不可用）",
		"err", err.Error(), "code", r)
}

// pump drains BOTH funnels on this leg's own thread: the message queue (the
// window's wndproc answers the session-end question here) and the events
// channel (the power callback's posts land here). The sequences run on this
// thread, so two doors firing at once serialize instead of racing the engines.
// The queue is pumped non-blocking (PeekMessage/PM_REMOVE) so a quiet queue
// still leaves the events channel responsive - the same tick shape
// internal/proc's RunEventLoop books.
func (le *lifecycleEvents) pump(ctx context.Context) {
	// The window is created HERE, on the thread that will pump it, so its
	// message queue always has a reader (the header comment on
	// startLifecycleEvents names the failure this ordering prevents).
	if !le.createWindow() {
		slog.Warn("lifecycle: WM_QUERYENDSESSION 门未建（窗口创建失败）；电源通知门照装",
			"class", lifecycleWindowClass)
	} else {
		slog.Info("lifecycle: WM_QUERYENDSESSION 门已建（lifecycle-sta 线程自泵）",
			"class", lifecycleWindowClass)
	}
	var m msg
	for {
		for {
			r, _, _ := pLifecyclePeekMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0, 1 /* PM_REMOVE */)
			if r == 0 {
				break
			}
			pLifecycleTranslateMessageW.Call(uintptr(unsafe.Pointer(&m)))
			pLifecycleDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
		}
		select {
		case <-ctx.Done():
			return
		case ev := <-le.events:
			le.run(ev)
		case <-time.After(50 * time.Millisecond):
		}
	}
}

// run executes one event. The suspend path runs inline (bounded); everything
// else is a receipt line.
func (le *lifecycleEvents) run(ev lifecycleEvent) {
	switch ev.kind {
	case "suspend":
		le.suspendDepth.Add(1)
		if le.sink.suspend != nil {
			le.sink.suspend(ev.why)
		}
	case "resume":
		le.resumeDepth.Add(1)
		if le.sink.resume != nil {
			le.sink.resume(ev.why)
		}
	case "session-end":
		if le.sink.sessionEnd != nil {
			le.sink.sessionEnd(ev.why)
		}
	}
}

// ---- the sequences ---------------------------------------------------------

// runSuspendSequence is the ticket's suspend half, in the order D40 row 1 and
// §16.11#10 lean freeze: approvals voided first (so nothing is approved while
// the world goes away), then capture off, then the session's engines released
// on the voice owner, bounded.
//
// Every step is counted and logged, so the sampler reading ("no models resident
// during sleep") has a paper trail next to the counters. The returned count is
// the voided-approvals N the evidence table quotes.
func runSuspendSequence(ra *residentApproval, raudio *residentAudio, rv *residentVoice, why string) int {
	voided := 0
	if ra != nil {
		voided = ra.voidForSuspend(why)
	}
	audioStopped := false
	if raudio != nil {
		audioStopped = suspendAudioCapture(raudio)
	}
	enginesReleased := false
	if rv != nil {
		enginesReleased = suspendVoiceEngines(rv, why, suspendVoiceBudget)
	}
	slog.Info("lifecycle: 挂起序列完成",
		"why", why, "voided_l2", voided, "audio_stopped", audioStopped,
		"engines_released", enginesReleased)
	if voided > 0 {
		fmt.Printf("wisp: 错过 %d 个待确认操作（系统挂起，已全部作废并按拒绝处理；恢复后不会重放）\n", voided)
	}
	return voided
}

// runResumeSequence is the resume half: capture back on, engines left disposed
// (the next real wake reloads them), the state machine untouched (no suspend
// edge exists in the frozen D43 table), wall schedules named as absent.
func runResumeSequence(raudio *residentAudio, why string) bool {
	reopened := false
	if raudio != nil {
		reopened = resumeAudioCapture(raudio)
	}
	slog.Info("lifecycle: 恢复序列完成",
		"why", why, "audio_reopened", reopened,
		"engines", "已随挂起释放，下一次真实唤醒按 D43#29 冷装载",
		"state_machine", "D43 表无挂起边，保持原状态",
		"wall_schedules", "本构建无墙钟派生日程装置（提醒实现未落地），无可重算项")
	return reopened
}

// suspendAudioCapture stops the device for the sleep (D40 row 1's first half).
// It deliberately does NOT walk ra.stop(): that is D38(e) step 4's process-exit
// body (it would also tear the voice chain down), and a sleep that ends must
// find the chain alive. False means there was no running gate to stop.
func suspendAudioCapture(ra *residentAudio) bool {
	if ra == nil || ra.gate == nil {
		return false
	}
	if !ra.started.Load() {
		return false
	}
	if err := ra.gate.Stop(); err != nil {
		slog.Warn("lifecycle: 挂起停采集返回错误", "err", err.Error())
		return false
	}
	return true
}

// resumeAudioCapture reopens the device after the wake, reusing the same gate,
// the same frame channel and a fresh consumer context. The gate's own
// Start/Stop pair is restart-shaped (running=false after Stop, ctx/buf re-stored
// by Start), so this is a reopen and not a rebuild. False is recorded by the
// caller: a device that will not reopen is the pre-existing degraded posture,
// not a new failure.
func resumeAudioCapture(ra *residentAudio) bool {
	if ra == nil || ra.gate == nil || ra.frames == nil {
		return false
	}
	ctx, cancel := context.WithCancel(context.Background())
	ra.cancel = cancel
	if err := ra.gate.Start(ctx, ra.frames); err != nil {
		cancel()
		slog.Warn("lifecycle: 恢复重开采集失败", "err", err.Error())
		return false
	}
	ra.started.Store(true)
	return true
}

// suspendVoiceEngines posts the epoch dispose onto the voice owner and waits a
// bounded monotonic budget. The dispose MUST run on the owner (a sherpa session
// closed under a decoding thread is the crash the owner-marshalling exists to
// prevent), and the wait must be bounded (the OS is going to sleep either way);
// a timeout is a recorded problem, not a hang.
func suspendVoiceEngines(rv *residentVoice, why string, budget time.Duration) bool {
	if rv == nil || rv.residency == nil || rv.queue == nil {
		return false
	}
	done := make(chan error, 1)
	if !rv.queue.post(func() { done <- rv.residency.closeEpoch("suspend:" + why) }) {
		slog.Warn("lifecycle: 挂起 dispose 未送达（voice owner 已停）")
		return false
	}
	tm := observe.NewTimeout(budget)
	select {
	case err := <-done:
		if err != nil {
			slog.Warn("lifecycle: 挂起 dispose 带错误返回", "err", err.Error())
		}
		return true
	case <-time.After(tm.Remaining()):
		slog.Warn("lifecycle: 挂起 dispose 未在预算内完成（系统照睡，唤醒后引擎按 dispose 语义已清）",
			"budget_ms", budget.Milliseconds())
		return false
	}
}

// ---- the sink wiring -------------------------------------------------------

// wireLifecycleSink builds the production sink from the legs the resident boot
// already assembled. It closes over the handles, so the sequences above read
// exactly like the boot report says they do.
func wireLifecycleSink(rt *proc.Runtime, ra *residentApproval, raudio *residentAudio, rv *residentVoice, fastExit *atomic.Bool) lifecycleSink {
	return lifecycleSink{
		suspend: func(why string) { runSuspendSequence(ra, raudio, rv, why) },
		resume:  func(why string) { runResumeSequence(raudio, why) },
		sessionEnd: func(why string) {
			// The fast path (D42#7 / shutdown.go's own comment): steps 4/5/9
			// never skip, only step 7's non-critical flushes do. The flag is
			// what the resident boot's deferred rt.Shutdown reads, and the stop
			// request is what ends the event loop the normal way - the D38(e)
			// chain still hangs off the boot's own defer, not off this callback.
			if fastExit != nil {
				fastExit.Store(true)
			}
			if rt != nil {
				rt.RequestStop("session-end (" + why + ")")
			}
			slog.Warn("lifecycle: 会话结束事件已接（快速关停已请求）", "why", why)
		},
	}
}

// voidForSuspend voids every pending card because the system is going to
// sleep. It walks the same funnel D38(e) step 3's cancelTaskRoots walks - the
// queue's own Reject for L2 (D43 row 24's refusal path, the legal row this
// lands on), the audit line + Forget for an L1 window (SPEC-06 §2 B1: a
// countdown can be opposed, not answered) - with the suspend wording, and it
// does NOT cancel the task root: the process is not exiting, the task continues
// after the wake minus the operation nobody confirmed.
//
// The returned count is the "错过 N 个待确认操作" notice's N.
func (ra *residentApproval) voidForSuspend(why string) int {
	refused, abandoned := 0, 0
	for _, card := range ra.cards.Pending() {
		switch card.Level {
		case "L2":
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			err := ra.cards.Reject(ctx, card.CorrelationID,
				"系统挂起：待确认操作已作废，按拒绝处理（未执行）")
			cancel()
			if err != nil {
				ra.residentAuditf("lifecycle: 挂起作废待批卡片 %s 失败：%v", card.CorrelationID, err)
				continue
			}
			refused++
		default:
			abandoned++
			ra.residentAuditf("approval: SUSPEND-WINDOW-VOIDED corr=%s tool=%s decision=reject "+
				"reason=%q channel=none", card.CorrelationID, card.Tool,
				"系统挂起，L1 确认窗口作废，未执行")
			ra.cards.Forget(card.CorrelationID)
		}
	}
	ra.residentAuditf("lifecycle: SUSPEND-L2-VOIDED why=%q refused_l2=%d voided_l1=%d", why, refused, abandoned)
	return refused + abandoned
}
