//go:build windows

package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/proc"
)

// runResident is the no-args path: boot (env layout, Job Object, single
// instance, goroutine registry), host the floating ball window and its tray
// (resident_ball_windows.go), take the task source AC#7 added
// (resident_task_source_windows.go), run the event loop, then exit through
// the D38(e) shutdown order. A second launch in the same session signals the
// running instance's activation event and exits (D42#7).
//
// This body is verbatim from main.go (ticket 78); it moved here because every
// call it makes is Windows-only - proc.Boot, proc.ErrAlreadyRunning and
// proc.SignalExistingInstance are declared in internal/proc's
// //go:build windows files, so an untagged main.go could not type-check under
// GOOS=linux. resident_other.go carries the refusal for every other platform.
func runResident() {
	printVersions("")

	env, err := buildinfo.ResolveEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "wisp: %v\n", err)
		os.Exit(2)
	}

	rt, err := proc.Boot(env)
	if errors.Is(err, proc.ErrAlreadyRunning) {
		layout, lerr := proc.DefaultLayout(env)
		if lerr == nil && layout.MutexEnabled {
			_ = proc.SignalExistingInstance(layout.ActivateEventName)
		}
		fmt.Println("wisp: another instance is running in this session; activated it; exiting")
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "wisp: boot failed: %v\n", err)
		os.Exit(1)
	}
	// Ticket 117: the resident process is the leg owner actually uses - double
	// click the icon, no terminal attached, stderr going nowhere. Everything
	// this process logs (the seal notices of internal/winsec, a locked config
	// section loosened under L2, the D33 credential migration, the panic sink)
	// needs a listener on disk, or "you will see it when an authorization is
	// cleared" is a sentence that only holds in a test binary.
	//
	// Ordering: proc.Boot is the only thing that ran before this, and it seals
	// nothing - internal/proc has zero winsec imports - so no event can have
	// been missed. The close defer is registered BEFORE the shutdown defer so
	// LIFO runs the D38(e) sequence first and its own log lines still land.
	sink, sinkErr := installLogSink(rt.Layout.DataDir)
	if sinkErr != nil {
		// Loud, and it does not stop the app: a log directory that will not
		// open must not become a way to keep Wisp from starting. The notices
		// fall back to stderr, which is the state before this ticket.
		fmt.Fprintf(os.Stderr, "wisp: 持久日志未启用（%v）：安全告警只会到 stderr，不会落盘\n", sinkErr)
	} else {
		defer sink.close()
	}
	defer func() {
		records := rt.Shutdown(false)
		failed := 0
		for _, rec := range records {
			if rec.Err != nil {
				failed++
				fmt.Printf("wisp: shutdown step %d (%s) FAILED: %v\n", rec.Step, rec.Name, rec.Err)
			}
		}
		fmt.Printf("wisp: exited through the D38(e) shutdown order (10 steps, %d failed)\n", failed)
	}()

	sum := rt.Layout.Summary()
	fmt.Printf("wisp: resident runtime booted (%s, data dir = %s, portable = %v, job object = on, single instance = %v)\n",
		sum.EnvBadge(), sum.DataDir, sum.Portable, rt.Instance != nil)

	// An exit request that arrives DURING boot is now a request, not a crash.
	// The handler that stops this process is installed inside RunEventLoop, and
	// the ball path in front of it - D2D factory, layered window, tray icon,
	// four RegisterHotKey calls - measured 268ms on the bench machine
	// (log stamps: sink installed 16:05:20.7746, ball booked 16:05:21.0424).
	// Without this registration a Ctrl+C inside that window is delivered to
	// nothing, so Go's console handler returns 0 (runtime/os_windows.go
	// ctrlHandler), the default disposition kills the process with
	// 0xc000013a, and the operator loses the whole D38(e) trail, the tray
	// removal and the hot key release. Registering above, while the sink is
	// already open and the shutdown defer is already in place, is what makes
	// the deferred chain below the only way this function ends.
	//
	// It stays registered for the rest of the function: os/signal forwards to
	// every registered channel, so proc's own context still gets the signal
	// once it exists, and the only gap left is the few instructions between the
	// check below and NotifyContext inside RunEventLoop - which costs one more
	// Ctrl+C, not an unclean death.
	bootExit := make(chan os.Signal, 1)
	signal.Notify(bootExit, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(bootExit)

	// Ticket 246 AC#1, form 乙 (orchestrator ruling, ledger A481): the approval
	// gate is BUILT HERE, by the assembly root, and handed down as a function
	// value and a binding. The resident file does not compose a gate of its own -
	// that would be the second truth source 246-a1 §5 reason 1 refuses, and it is
	// the same shape ticket 238's cut-1 blocked ("两枚正向依赖边一律不开、改注入",
	// extended to this ticket by the ruling named above).
	//
	// Zero new package-level dependency edges: cmd/wisp already imports
	// internal/agent/approval, internal/ball, internal/proc and internal/tools in
	// production code (run.go:47, resident_ball_windows.go:41, this file's own
	// import block), and internal/proc and internal/ball stay leaves. The ruler
	// for that claim is `GOOS=windows go list -deps ./cmd/wisp`, run before and
	// after, in this leg's evidence table.
	// Ticket 256 AC#1 (the [risk] half, orchestrator ruling §8.2 form 甲ⓐ): this
	// is the one production call site, and it now passes rt.Layout.DataDir so the
	// gate is built from this host's confirm_timeout_sec / l1_window_sec instead
	// of the compiled 300s / 3s. rt.Layout.DataDir is already resolved here (the
	// same value installLogSink takes at :66 and newResidentPanelManager takes at
	// :151). Construction time only - see newResidentApprovalWithConfig.
	ra := newResidentApprovalWithConfig(rt.Layout.DataDir)

	// Ticket 33 AC#1..AC#4: the panel host, on its own dedicated STA thread
	// (orchestrator ruling P1). It is BUILT here, by the assembly root, and handed
	// to the ball host as a function value - the same shape ticket 246's ruling
	// (ledger A481) set for the cancel executor, so resident_ball_windows.go keeps
	// knowing nothing about WebView2 and this file keeps knowing nothing about the
	// approval UI.
	//
	// Building the manager is the first non-test call of NewPanelManager in this
	// repository: until now the host existed only in a test binary, which is what
	// made AC#1..AC#4 read as "the host works" instead of "the user can open it".
	//
	// A panel that cannot be assembled does not stop the boot (same posture as
	// installLogSink and startResidentBall): it is printed, and the boot report
	// below says which shape this process is in. Its teardown is the defer that
	// runs AFTER rb.stop() above registers - LIFO, so the panel thread ends last
	// and the frozen D38(e) ten-step order (internal/proc/shutdown.go) is not
	// touched: no new step, no new hook name on that closed roster.
	rp, rpErr := newResidentPanelManager(rt.Layout.DataDir)
	var panel *residentPanel
	if rpErr != nil {
		slog.Error("panel host: this process has no panel", "err", rpErr)
		fmt.Printf("wisp: panel host unavailable (%v): the panel hot key and the tray item will be recorded, not executed\n", rpErr)
	} else {
		panel = startResidentPanel(rt.Registry, rp)
		defer panel.stop()
	}

	// Ticket 228 AC#1: this leg is the process the ball belongs in (D2,
	// PLAN.md:74 and :83-88 - layered window, tray, hot keys and the Job Object
	// holder in one resident main process). The defer below is registered last,
	// so LIFO runs it FIRST: D38(e) step 2's work ("hotkey + wake-word
	// listening stops", internal/proc/shutdown.go:18) happens ahead of the
	// sequence it is part of, and both the shutdown trail and the ball's own
	// records still reach the sink installed above.
	//
	// A ball that will not come up does not stop the boot: the window, the tray
	// and the hot keys are reported by what Win32 actually returned, and the
	// sentence printed under this call is built from that same result.
	// Ticket 258 AC#1: the [hotkey] chain, built here (the assembly-root shape
	// tickets 246 and 33 set for everything this leg hosts). The ball lives at
	// this point of the boot; the task pipeline that owns rt.mgr does not exist
	// until startResidentTaskSource below (258-a1 census Q1c), so the chain
	// re-reads config.toml per call instead of holding a Manager - the same
	// per-use read shape panelGeometrySource has been since ticket 255 AC#4.
	// Every branch names its provenance out loud: a missing or unreadable file
	// is the AC#1 fallback to ball.DefaultHotkeys() with the word "defaults"
	// in the verdict, never a silent swap.
	hotCfg258 := func() ball.HotkeyConfig {
		cfgPath := filepath.Join(rt.Layout.DataDir, configFileName)
		c, _, err := config.LoadFile(cfgPath, nil)
		if err != nil || c == nil {
			slog.Warn("ball: [hotkey] source unreadable at construction; falling back to the compiled defaults",
				"path", cfgPath, "err", err, "fallback", "DefaultHotkeys")
			fmt.Printf("wisp: ball [hotkey]: config.toml unreadable (%v); the four hotkeys fall back to the compiled defaults (DefaultHotkeys)\n", err)
			return ball.HotkeyConfig{}
		}
		h := c.Hotkey
		slog.Info("ball: [hotkey] view read from config.toml",
			"summon", h.Summon, "mute", h.Mute, "cancel", h.Cancel, "panel", h.Panel,
			"empty_slots_note", "empty slots are filled from the product defaults by ball.ApplyHotkeyDefaults; "+
				"when the merged set is the compiled defaults the verdict prints the word defaults (ticket 258-r2)")
		return ball.HotkeyConfig{Summon: h.Summon, Mute: h.Mute, Cancel: h.Cancel, Panel: h.Panel}
	}
	// The hot-tier half (form A): the bridge's src re-reads config.toml each
	// tick. The task pipeline's Manager (assembleRuntime's rt.mgr) does not
	// exist until startResidentTaskSource below - the ball is built first, the
	// 258-a1 census measured that order - so the src is a per-tick fresh read of
	// the file, the exact shape panelGeometrySource has been since ticket 255
	// AC#4 and models.go:184 before it. CheckAndReload's diff still runs inside
	// the bridge (it diffs against what the ball holds), so an unchanged file
	// costs zero Win32 calls.
	hotReload258 := func() ball.HotkeyConfig {
		c, _, err := config.LoadFile(filepath.Join(rt.Layout.DataDir, configFileName), nil)
		if err != nil || c == nil {
			// The bridge keeps the current bindings on an empty answer
			// (hotkey_reload.go leaves the applied set untouched when the source
			// errors into all-empty after the first diff), and the refresh Warn
			// it prints names the failure - the AC#1 fallback stays said, per
			// tick, in the log.
			return ball.HotkeyConfig{}
		}
		return ball.HotkeyConfig{Summon: c.Hotkey.Summon, Mute: c.Hotkey.Mute, Cancel: c.Hotkey.Cancel, Panel: c.Hotkey.Panel}
	}
	rb := startResidentBall(rt.Registry, ra.vetoByEsc, hotCfg258, hotReload258, withPanelHost(func(via string) bool {
		return panel.RequestToggle(via)
	}))
	defer rb.stop()

	// The gate is bound to the ball only if the ball exists, and it is detached
	// BEFORE the window goes away: after Ball.Close posts WM_QUIT there is no STA
	// thread left to answer a release, and a card settling at that moment must
	// say so rather than block. Defer order is the mechanism - this line is
	// registered after rb.stop() above, so LIFO runs it first.
	defer ra.detachBall()
	ra.bindBallHost(rb)

	// Ticket 246 AC#4: D38(e) step 3 ("all task root ctxs cancelled -> wait
	// <= 3s") gets the one thing it never had - a producer. Before this,
	// internal/proc/boot_windows.go built a fresh ShutdownHooks{} per call and
	// filled only CloseJob, so a card still waiting when an exit signal arrived
	// was killed by the OS instead of being refused, and the audit trail recorded
	// the step as skipped. The ORDER is untouched: RunShutdownSequence still
	// walks 1..10 (shutdown.go:112, pinned by TestShutdownOrderAudit), and a step
	// nobody registered is still honestly recorded as skipped - which is why a
	// refusal to register below is printed and not fatal: the record stays true
	// either way, it just says the harder thing.
	if err := rt.RegisterShutdownHook(proc.StepCancelTasks, ra.cancelTaskRoots); err != nil {
		slog.Error("proc: 第 3 步（取消任务根）未能注册为钩子，退出时待确认项不会被有序拒绝",
			"step", int(proc.StepCancelTasks), "err", err)
		fmt.Printf("wisp: 退出序列第 3 步未注册（%v）：卡片挂起时收到退出信号，这一步会被记为 skipped\n", err)
	}

	// Ticket 246 AC#7: the hop every earlier leg left open. The gate above can
	// show a card, the ball above can veto one and step 3 above can refuse one on
	// the way out - and AskOnTaskRoot / askConfirmation still had zero product
	// callers, so nothing was going to raise one. This call is the caller: it reads
	// the operator's console through interactiveStdin() (the same gate ticket 201's
	// answer side uses, and the only one this process may test), assembles the task
	// pipeline through cmd/wisp/run.go with ra's gate injected, and books the two
	// D38(e) steps that entry owns (1 and 7). Every degraded branch inside it - no
	// console, no config.toml, no credential - prints and returns nil, because the
	// only condition that may stop this boot is the one ticket 128 settled.
	//
	// It runs after bindBallHost so the veto channel this leg advertises is
	// already the honest one, and before the report below so the step roster the
	// report prints includes whatever this registered.
	src := startResidentTaskSource(rt, ra)

	// The boot report names the steps this process really owns, taken off the
	// registration rather than off a comment (AC#4's "不许撒谎" half).
	registered := rt.RegisteredShutdownSteps()
	names := make([]string, 0, len(registered))
	for _, s := range registered {
		names = append(names, fmt.Sprintf("%d:%s", int(s), s.Name()))
	}
	fmt.Printf("wisp: %s; 任务来源：%s; 面板：%s; D38(e) steps with an owner in this process: %s\n",
		ra.residentStatusLine(), src.taskPosture(), panel.statusLine(), strings.Join(names, ", "))

	// The boot report has to match what happens next: if an exit request already
	// arrived during the ball path, printing "resident event loop running" and then
	// never entering the loop would be exactly the kind of sentence AC#7 exists
	// for. The ball status and the task-source posture are true in either branch,
	// so they print in either.
	var reason string
	select {
	case sig := <-bootExit:
		reason = "exit request (" + sig.String() + ") arrived during boot; the event loop was never entered"
		fmt.Printf("wisp: %s; %s\n", reason, rb.statusLine())
	default:
		fmt.Printf("wisp: resident event loop running (task source: %s); %s (Ctrl+C exits cleanly)\n",
			src.taskPosture(), rb.statusLine())
		reason = rt.RunEventLoop()
	}
	fmt.Printf("wisp: event loop ending (%s); running the D38(e) shutdown order\n", reason)
}
