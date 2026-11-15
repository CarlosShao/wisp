//go:build windows

package main

// Ticket 46's D32 settle row, in production: the one consumer of
// watchdog.Loop.VerifySettle.
//
// WHAT WAS WRONG BEFORE THIS FILE. The watchdog loop was constructed at boot
// (startResidentWatchdog) and then thrown away - resident_windows.go carried
// `_ = wd` with a comment saying the dispose-path consumer "rides the
// session-scope wave (ticket 28)". Ticket 28 closed weeks ago, so the sentence
// had lost its premise and the layer it promised never arrived: VerifySettle's
// only caller in the whole tree was a test, which means the D32 settle row -
// "after a dispose the tree is back under the state cap within 10s, with the
// FreeOSMemory counter > 0" - had ZERO executions in the product. A gate nobody
// runs is not a gate; it is a comment with a function body.
//
// WHERE IT RUNS, AND WHY THERE. The trigger is the session scope's dispose
// cascade, step 3 (Hooks.FreeOSMemory, resident_voice_engines_windows.go's
// scopeOptions). That is the only instant in this process where all three
// halves of the row are simultaneously true: the engines are released (step 1),
// the panel has been hidden or destroyed (step 2), and the heap has just been
// handed back to the OS (step 3). Sampling must therefore START there and not
// earlier - a check that begins while the engines are still hot measures the
// session, not the settle. And it must not run ON that thread: VerifySettle
// blocks for the whole window by construction (observe.CheckSettle samples until
// its deadline and never returns early - sampler.go:495-526), so it runs on its
// own goroutine, spawned through observe.Registry.Spawn under the D38b roster
// name the watchdog family already uses (the same booking cmd/wisp/config_reload.go
// makes for the config tick: name "watchdog", a different owner).
//
// WHAT KEEPS IT HONEST RATHER THAN MERELY RUNNING:
//
//   - released is the DELTA of observe.FreeOSMemoryCount() across the release
//     call, not a constant. If a future change hands the release back to
//     debug.FreeOSMemory directly, every verification fails by contract instead
//     of silently passing a check that never saw a release.
//   - the target state is read from the live D43 machine (the same mapping the
//     loop gates its resource duty with, sloStateFromD43Machine). Post-dispose is
//     Sleeping when nothing is armed and Armed when KWS is still resident, and
//     the two caps differ by 65MB (25 vs 90) - a hardcoded Sleeping would report
//     a correct Armed process as a settle failure.
//   - peak is one real tree read taken before the release; a report whose
//     peak_bytes is a guess is not evidence.
//   - a wake cancels an in-flight verification (voiceResidency.openScope). The
//     question is "did the disposed tree fall back", never "is the tree small
//     while the next session boots"; without the cancel three ordinary wakes
//     would walk D43 #35 (WatchdogAlert + one-click restart) on a healthy
//     process.
//   - the failure path is the existing one and is not re-invented here: three
//     CONSECUTIVE failures raise AlertFatal from inside the loop
//     (watchdog/settle.go:67-76), the host maps it to D43 #35
//     (dispatchWatchdogAlert), and the numbers land in the loop's own
//     "watchdog: settle verification" record - which is the JSONL log file.

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/watchdog"
)

// settleRequest is one dispose's evidence for the gate: the tree footprint read
// immediately before the release, and whether that release really went through
// observe.ReleaseMemory (the instrumented wrapper the D32 row's counter is).
type settleRequest struct {
	peak     int64
	released bool
}

// residentSettleGate is the host-side driver: it owns the goroutine, the
// coalescing and the readings, while the verdict itself stays the loop's
// (VerifySettle, internal/watchdog - one settle rule in the tree, not two).
type residentSettleGate struct {
	loop    *watchdog.Loop
	tree    observe.TreeReader
	reg     *observe.Registry
	root    *observe.Root
	ctx     context.Context
	stateFn func() observe.SLOState

	mu        sync.Mutex
	cancel    context.CancelFunc // the in-flight verification, if any
	requests  int
	coalesced int
	runs      int
	cancelled int
	failures  int

	lastTarget observe.SLOState
	lastPass   bool
	lastPeak   int64
	lastBackMS int64
	lastFinal  int64

	// attached / unattachedReason are the mount half of the boot row (the 46b
	// audit's "另记 2 条" #1): an assembled-but-unmounted gate must say so with
	// a named reason instead of the assembled sentence, because a boot whose
	// voice leg never materialized (first install without `wisp models ensure`
	// is the common shape) has NO trigger in the dispose cascade - claiming
	// "装配" there was a false statement measured by the audit.
	attached         bool
	unattachedReason string
}

// newResidentSettleGate arms the gate. loop and tree are required - the verdict
// needs a measuring instrument, and the peak read needs the same reader the loop
// samples with (one TreeSampler over the C30 Job, so the two can never disagree
// about what the tree held). reg nil falls back to the process registry, the
// same default watchdog.New uses. root's ctx bounds every verification: a
// process that is leaving stops measuring.
func newResidentSettleGate(loop *watchdog.Loop, tree observe.TreeReader, reg *observe.Registry, root *observe.Root, stateFn func() observe.SLOState) *residentSettleGate {
	if loop == nil || tree == nil {
		return nil
	}
	if reg == nil {
		reg = observe.Default
	}
	ctx := context.Background()
	if root != nil {
		ctx = root.Ctx
	}
	return &residentSettleGate{
		loop: loop, tree: tree, reg: reg, root: root, ctx: ctx, stateFn: stateFn,
		lastTarget: observe.SLOSleeping,
	}
}

// readPeak is the pre-release tree read. Called on the dispose thread (the
// cascade's third step) and therefore cheap by construction: one ReadTree, no
// allocation of note, and a failure that is named in the report rather than
// papered over with a zero.
func (g *residentSettleGate) readPeak() int64 {
	if g == nil || g.tree == nil {
		return 0
	}
	m, err := g.tree.ReadTree()
	if err != nil {
		slog.Warn("settle gate: 释放前的树读数失败（peak_bytes 会是 0）",
			"err", err.Error(), "error_class", string(observe.ClassResource))
		return 0
	}
	if m.PrivateWorkingSetBytes <= 0 {
		slog.Warn("settle gate: 释放前的树读数为 0，不当读数用（peak_bytes 会是 0）",
			"error_class", string(observe.ClassResource))
		return 0
	}
	return m.PrivateWorkingSetBytes
}

// settleTarget reads the state the tree is being judged against, off the live
// D43 machine. A machine in a state with no D32 sampling row of its own answers
// the zero SLOState (sloStateFromD43Machine's documented shape), and the row
// this gate serves is the at-rest one, so Sleeping is the fallback - said out
// loud in the record rather than silently substituted.
func (g *residentSettleGate) settleTarget() (observe.SLOState, bool) {
	if g != nil && g.stateFn != nil {
		if st := g.stateFn(); st.Valid() {
			return st, false
		}
	}
	return observe.SLOSleeping, true
}

// request starts one verification, or coalesces into the one already running.
// Non-blocking by contract: it is called from the dispose cascade, whose thread
// must not be held for the settle window.
func (g *residentSettleGate) request(req settleRequest) {
	if g == nil || g.loop == nil {
		return
	}
	g.mu.Lock()
	g.requests++
	if g.cancel != nil {
		g.coalesced++
		n := g.coalesced
		g.mu.Unlock()
		// Not a silent merge: a dispose whose settle went unverified is exactly
		// the kind of absence this ticket exists to end, so it is counted and
		// said. The verification already in flight is still measuring THIS tree.
		slog.Info("settle gate: 上一次校验还在窗口内，本次 dispose 的校验并入（不并发测同一棵树）",
			"coalesced", n, "peak_bytes", req.peak, "released", req.released)
		return
	}
	vctx, cancel := context.WithCancel(g.ctx)
	g.cancel = cancel
	g.mu.Unlock()

	// D38b roster name, watchdog family, its own owner: the same booking shape
	// cmd/wisp/config_reload.go already makes for a watchdog-side background
	// loop ("watchdog" + owner "config"), so the frozen roster is reused rather
	// than extended.
	g.reg.Spawn("watchdog", "settle-gate", g.root, func(context.Context) {
		// The ctx Spawn handed this body is the root ctx; vctx is that same
		// root plus this run's cancel, which is what abort() flips.
		target, fellBack := g.settleTarget()
		pass, rep, err := g.loop.VerifySettle(vctx, target, req.released, req.peak)
		g.finish(target, fellBack, pass, rep, err)
	})
}

// finish books one verification's outcome. A cancelled run is NOT a failure
// (watchdog/settle.go returns the ClassCancelled error before it touches its
// consecutive-failure counter, so the loop's own count is authoritative too).
func (g *residentSettleGate) finish(target observe.SLOState, fellBack bool, pass bool, rep *observe.SettleReport, err error) {
	g.mu.Lock()
	g.cancel = nil
	g.runs++
	g.lastTarget = target
	g.lastPass = pass
	if err != nil {
		g.cancelled++
	}
	if rep != nil {
		g.lastPeak = rep.PeakBytes
		g.lastBackMS = rep.BackWithinCapMS
		g.lastFinal = rep.FinalBytes
	}
	if err == nil && !pass {
		g.failures++
	}
	failures := g.failures
	final, back := g.lastFinal, g.lastBackMS
	g.mu.Unlock()

	if err != nil {
		// The loop returns the cancellation as a classified error and counts
		// nothing; saying so here is what keeps "no verification ran" apart from
		// "a verification ran and was ended by a wake".
		slog.Info("settle gate: 校验没有走完（被取消或被截断），不计失败",
			"target_state", string(target), "err", err.Error())
		return
	}
	if !pass {
		slog.Warn("settle gate: D32 settle 行本发不过（连续失败数由 watchdog 侧升 D43#35）",
			"target_state", string(target),
			"final_bytes", final, "back_within_cap_ms", back,
			"gate_failures", failures,
			"error_class", string(observe.ClassResource))
		return
	}
	slog.Info("settle gate: D32 settle 行本发通过",
		"target_state", string(target),
		"state_from_machine", !fellBack,
		"final_bytes", final, "back_within_cap_ms", back)
}

// abort cancels the in-flight verification. Called when a wake reopens the
// session: the disposed tree is no longer the thing being measured.
func (g *residentSettleGate) abort(reason string) {
	if g == nil {
		return
	}
	g.mu.Lock()
	cancel := g.cancel
	g.mu.Unlock()
	if cancel == nil {
		return
	}
	cancel()
	slog.Info("settle gate: 校验被唤醒取消（树已不是 dispose 后的那棵）", "reason", reason)
}

// line is the boot report's settle-gate row: what is armed, where it fires and
// what its evidence is. It states no window number on purpose - the window is
// the loop's own SettleWithin (D32's 10s row, watchdog.go:225-228) and this
// assembly does not restate a frozen default it does not own.
//
// Three states, per the 46b audit's ruling (the old two-state line was a false
// claim on every boot whose voice leg never materialized: first install
// without `wisp models ensure` is the COMMON shape, the gate sat assembled but
// unmounted, and the row still said the dispose cascade would trigger it):
//
//   - 未装配 (g == nil): no watchdog loop, no verifier at all;
//   - 已装配但未挂载 (assembled, attachSettleGate found nothing to mount on):
//     the loop and the tree exist, but this boot built no voice residency, so
//     the dispose cascade's third step will never reach the gate - the reason
//     is named, taken from the boot's own voice-leg verdict;
//   - 装配 (mounted): the cascade's third step really calls this gate.
func (g *residentSettleGate) line() string {
	if g == nil {
		return "收口闸：未装配（本进程没有 watchdog 循环，D32 settle 行在这次启动里没有验证者）"
	}
	g.mu.Lock()
	attached, reason := g.attached, g.unattachedReason
	g.mu.Unlock()
	if !attached {
		if reason == "" {
			reason = "startResidentWatchdog 尚未执行挂载（装配与挂载之间）"
		}
		return "收口闸：已装配但未挂载（" + reason +
			"）：dispose 级联第 3 步摸不到这枚闸，D32 settle 行在这次启动里不会被触发验证"
	}
	return "收口闸：装配（触发＝会话 dispose 级联第 3 步 FreeOSMemory；目标态＝D43 机现态，非采样态回落 Sleeping；" +
		"窗口＝watchdog 的 SettleWithin；证据＝watchdog settle verification 日志行＋三次连续失败升 D43#35）"
}

// markAttached / markUnattached record the mount half of line(). attachSettleGate
// (resident_windows-side wrapper) is the only caller, so the row can only move
// from unmounted to mounted, never the other way.
func (g *residentSettleGate) markAttached() {
	if g == nil {
		return
	}
	g.mu.Lock()
	g.attached = true
	g.unattachedReason = ""
	g.mu.Unlock()
}

func (g *residentSettleGate) markUnattached(reason string) {
	if g == nil {
		return
	}
	g.mu.Lock()
	g.attached = false
	g.unattachedReason = reason
	g.mu.Unlock()
}

// attachSettleGate hands the gate to the voice leg's residency, which is what
// the dispose cascade's third step reads. On a boot with no voice leg (or no
// residency) there is nothing to arm and the boot row says the gate is
// assembled-but-unmounted rather than pretending a trigger exists - with the
// named reason the boot itself produced (buildResidentVoice's verdict), so the
// row points at the actual cause instead of a generic "no residency".
//
// The returned bool is the mount result the boot row reads; true means the
// residency now holds the gate.
func (rv *residentVoice) attachSettleGate(g *residentSettleGate) bool {
	if rv == nil {
		if g != nil {
			g.markUnattached("本进程没有语音腿（residentVoice 为 nil）")
		}
		return false
	}
	if rv.residency == nil {
		if g != nil {
			g.markUnattached("语音腿未构造无常驻（" + rv.verdict + "）")
		}
		return false
	}
	return rv.residency.attachSettleGate(g)
}

// settleGateSnapshot is the readable state (diagnostics and the evidence
// readings; nothing in the boot path depends on it).
type settleGateSnapshot struct {
	Requests  int
	Coalesced int
	Runs      int
	Cancelled int
	Failures  int
	Target    observe.SLOState
	Pass      bool
	Peak      int64
	BackMS    int64
	Final     int64
}

func (g *residentSettleGate) snapshot() settleGateSnapshot {
	if g == nil {
		return settleGateSnapshot{}
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	return settleGateSnapshot{
		Requests: g.requests, Coalesced: g.coalesced, Runs: g.runs,
		Cancelled: g.cancelled, Failures: g.failures,
		Target: g.lastTarget, Pass: g.lastPass,
		Peak: g.lastPeak, BackMS: g.lastBackMS, Final: g.lastFinal,
	}
}

// String keeps the gate readable in a log line an operator greps.
func (g *residentSettleGate) String() string {
	s := g.snapshot()
	return fmt.Sprintf("settle-gate requests=%d coalesced=%d runs=%d cancelled=%d failures=%d last_target=%s last_pass=%v back_ms=%d",
		s.Requests, s.Coalesced, s.Runs, s.Cancelled, s.Failures, string(s.Target), s.Pass, s.BackMS)
}
