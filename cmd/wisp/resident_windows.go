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
	"sync"
	"syscall"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/proc"
)

// hotkeyReloadSource296 is ticket 258's hot-tier [hotkey] source - the reload
// bridge's src, hoisted out of runResident so the ticket-296 case in
// cmd/wisp/resident_hotkey_296_windows_test.go drives the body production
// actually runs. Ticket 258 hoisted the construction half the same way
// (residentBallHotkeyChain258); the reason is the same: a re-typed copy of a
// closure inside a test is not an instrument for a bug that lives in the
// closure's own body, it reads green whichever way production changed.
//
// It answers what config.toml says, one fresh read per call, and nothing else -
// the per-tick shape panelGeometrySource has used since ticket 255 AC#4. The
// bridge diffs against what the ball holds, so an unchanged file costs zero
// Win32 calls.
//
// NO SLOT IS EVER ANSWERED EMPTY (ticket 296). Both branches merge through
// ball.ApplyHotkeyDefaults, which is the precondition ball.HotkeySource states
// in its own doc comment and the shape the construction half gets from
// residentBallHotkeyChain258. The bridge does NOT keep the current bindings on
// an empty answer: Check()'s only skip is the DeepEqual diff against what is
// applied, and internal/ball has no all-empty guard, so four empty strings from
// here would re-register the set and leave summon, mute and panel dead until
// the next read - which is exactly the regression this closure caused on the
// owner's machine (an empty [hotkey] slot means "unset", not "turned off").
func hotkeyReloadSource296(dataDir string) ball.HotkeyConfig {
	c, _, err := config.LoadFile(filepath.Join(dataDir, configFileName), nil)
	if err != nil || c == nil {
		// Ticket 296: the unreadable branch answers the compiled defaults, not
		// four empty strings - the same fallback hotCfg258 prints a Warn for at
		// construction, said per tick by the value it hands the bridge.
		return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{})
	}
	// Ticket 296: the three slots the owner's file leaves empty are "unset", and
	// the bridge would read them as "disabled by config". Merge before the hop.
	return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{Summon: c.Hotkey.Summon, Mute: c.Hotkey.Mute, Cancel: c.Hotkey.Cancel, Panel: c.Hotkey.Panel})
}

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
	// Ticket 186: the second return value is the holder the panel's
	// panel.workspace.request narrows once the task pipeline exists. It is handed
	// to the task source below, which is the only object that owns the runtime's
	// path scope; a boot whose panel did not come up passes nil, and the task
	// source then binds nothing (its switches are refused by name rather than
	// narrowing a scope no panel can reach).
	rp, panelScope, rpErr := newResidentPanelManager(rt.Layout.DataDir)
	// Ticket 319 AC#1: the assembly root's handle on the panel thread is a HOLDER
	// now, not the thread itself. The one thing that changes in a visible way:
	// after the thread has ended (any way at all - the owner's X, a system
	// SC_CLOSE, a named refusal), the next gesture starts a fresh thread on the
	// same manager instead of being refused forever. See panelHolder above for
	// the defect, the field, and the three frozen things it leaves alone.
	holder := newPanelHolder(rt.Registry, rp)
	if rpErr != nil {
		slog.Error("panel host: this process has no panel", "err", rpErr)
		fmt.Printf("wisp: panel host unavailable (%v): the panel hot key and the tray item will be recorded, not executed\n", rpErr)
	} else {
		// The eager first thread, the same step the boot has always taken (the
		// thread exists from boot; the window still appears only when asked).
		// What moved is who holds the name: the holder, which is also what
		// re-starts a thread later.
		holder.start()
	}
	defer holder.stop()

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
	// Ticket 296: hotReload258's body is the package-scope hotkeyReloadSource296
	// above, hoisted so the ticket-296 case drives the SAME code the assembly
	// root runs (ticket 258 hoisted the construction half the same way, as
	// residentBallHotkeyChain258). The two halves stay two separate closures on
	// purpose: the construction view is read once and merged by its consumer
	// (residentBallHotkeyChain258), the reload view is merged by this source
	// itself because its consumer - the bridge - performs no merge at all
	// (ticket 296). Hoisting is what makes the merge measurable: a re-typed copy
	// of the closure inside a test reads green no matter what production does.
	hotReload258 := func() ball.HotkeyConfig { return hotkeyReloadSource296(rt.Layout.DataDir) }
	rb := startResidentBall(rt.Registry, ra.vetoByEsc, hotCfg258, hotReload258, withPanelHost(func(via string) bool {
		return holder.toggle(via)
	}))
	defer rb.stop()

	// The gate is bound to the ball only if the ball exists, and it is detached
	// BEFORE the window goes away: after Ball.Close posts WM_QUIT there is no STA
	// thread left to answer a release, and a card settling at that moment must
	// say so rather than block. Defer order is the mechanism - this line is
	// registered after rb.stop() above, so LIFO runs it first.
	defer ra.detachBall()
	ra.bindBallHost(rb)

	// Ticket 319 AC#2: the hop every earlier leg left open, and the ONE line that
	// turns a second launch into something the operator can see. Before this,
	// RunEventLoop's activation branch booked and cleared the signal and said out
	// loud why it could do nothing else ("no path from this loop to a window") -
	// true as a fact about internal/proc, false as a fact about this process,
	// which has held a window-raising executor since ticket 33 (the ball's panel
	// gestures; :257 above drives the same one). The owner's 2026-10-14 log has 12
	// "activation requested by second launch" entries, one of them (21:40:38)
	// while the panel thread and its window were both alive: the double click was
	// silence because nobody had connected these two ends, not because the panel
	// was dead.
	//
	// The value handed over is closed over the holder and the ball host and does
	// the work itself - SetActivationHook only stores it, so internal/proc gains
	// no dependency edge (Runtime.SetActivationHook takes a func() string, the
	// same shape RegisterShutdownHook set). It is registered here, after the ball
	// exists (the fallback tip needs it) and before the event loop below (the only
	// consumer). Ticket 319-r3 corrected which holder executor it drives: show,
	// not toggle - a toggle hides a panel that is already on screen, and AC#2 is
	// about bringing it to the front (see secondLaunchActivation319's comment).
	rt.SetActivationHook(secondLaunchActivation319(holder, rb))

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
	//
	// panelScope (ticket 186) rides in with it: the pipeline this call assembles
	// owns the ONE path scope the panel's workspace request may narrow, and the
	// bind happens inside, on the first line where that object exists.
	src := startResidentTaskSource(rt, ra, panelScope)

	// Ticket 329 form 甲: the one panel host this process owns becomes the
	// destination of the task pipeline's booked snapshot packets. This is the
	// attach, and it is the whole of the wiring - publishPanelSnapshot pushes
	// through whatever host is set here, and PanelManager.evalPanelJS routes each
	// push to the panel thread (no window open, nothing pushed; the packet stays
	// booked as before). It sits on this line because src.run is the first moment
	// the runtime exists here, and every later line is a process whose pumps could
	// already publish. A refused assembly (src nil or src.run nil) and a boot whose
	// panel host never built (rp nil) each leave the host unset, which is the
	// honest statement that this process has no page to push to.
	if src != nil && src.run != nil && rp != nil {
		src.run.setPanelSnapshotHost(rp)
	}

	// Ticket 247 AC#1/AC#5: the capture leg. This is the line that makes
	// internal/audio have a non-test importer - before it the whole stack
	// (real WASAPI shared mode, the D38d bounded frame path, the gate, ticket
	// 241's level scale) had zero production callers, and the only number ever
	// handed to Ball.SetAudioLevel came from cmd/balldebug's command line.
	//
	// It runs AFTER startResidentBall because the level's destination is that
	// window, and BEFORE the report below so the step roster the report prints
	// includes D38(e) step 4 if this leg really owns it. The ball host is
	// handed a nil-safe method value taking one float32: no sample, no
	// transcript and no second dependency edge cross into internal/ball (AC#3),
	// and internal/audio still does not import the ball (AC#0 甲 + 乙-2 exit).
	//
	// No defer here on purpose: the teardown is the StepStopAudio hook this call
	// registers, so the frozen ten-step order closes the device - not a LIFO
	// accident that would leave step 4 reporting a stop it did not perform.
	raudio := startResidentAudio(rt, rb.setAudioLevel)

	// Ticket 290 AC#2 (form 甲-1, orchestrator ruling 2026-10-09 14:2x): the hop
	// that makes the default posture something a user can leave. The capture leg
	// exists from the line above, so the two mute gestures - the orb's global mute
	// hot key and the tray's mute item - are now pointed at the gate THIS process
	// assembled, and turning it is their only effect.
	//
	// Why a setter and not a parameter: the ball is built first (its window is
	// where the level is going, so startResidentBall has to precede the capture
	// leg), while the gate is built here. residentBall's own lock orders this
	// handoff against a key press that arrives inside that boot window, and the
	// gesture says so out loud when it does.
	//
	// What this deliberately does NOT do: no state machine event is dispatched
	// (D43's table has no Sleeping -> Muted edge, and inventing one is a
	// human-approved contract change, not a wiring detail - ticket 290's ruling
	// names that follow-up as outside this installment), no second "who is muted"
	// flag anywhere
	// (the gate's own flag is the truth source, internal/audio/gate.go:20-21), no
	// new goroutine (the gestures fire on the ui-sta thread and the device open is
	// synchronous inside internal/audio, on the roster name ticket 247 already
	// booked), and - as of ticket 293 - a one-way MIRROR of the gate into the
	// tray's own checkmark. The relation is master and projection, not two
	// authorities: gate.Muted() (internal/audio/gate.go) is the only truth about
	// mute, the ball-side trayMuted flag is display state written FROM it and read
	// back by nothing, and the hop below is the only production writer of
	// Ball.SetTrayChecks. Undriven stopped being the honest choice the moment
	// ticket 290 made the tray item turn the real gate: an unchecked "静音" sitting
	// next to a microphone the user just opened is a privacy lie, so the checkmark
	// now follows the gate instead of staying a stranger to it.
	rb.attachMuteGate(raudio.toggleMute)

	// Ticket 293 AC#2 (form 甲, orchestrator ruling 2026-10-09 20:2x): injected
	// from THIS side, so internal/ball keeps zero changes, zero new package-level
	// dependency edges and zero new goroutines.
	//
	// Why this one shape covers all three ways the shown state can change:
	//   - the tray's mute item and the orb's global mute hot key both land in
	//     residentBall.muteGesture, and that function re-projects after every flip
	//     the gate accepted;
	//   - the factory posture (audio.mic_muted_default, shipped true) never passes
	//     through a gesture, so the line below projects it once, immediately after
	//     the leg that owns the gate exists.
	// The menu is rebuilt from b.trayMuted on every right-click
	// (internal/ball/ball_windows.go's showMenu hop), so pushing the projection on
	// each change is exactly what makes the next right-click show the state as of
	// that moment - and a projection nobody pushed would stay at the zero value
	// forever, which is the bug this ticket is.
	//
	// The nil guard comes first because it has to precede the method value: taking
	// rb.b.SetTrayChecks on a host whose window never came up would panic at the
	// expression itself (startResidentBall returns with rb.b == nil on that
	// branch), and every other ball-facing hop in this process keeps the guard.
	if rb.b != nil {
		rb.attachTrayMuteProjection(raudio.trayMuteState, rb.b.SetTrayChecks)
		rb.mirrorTrayMute()
	}

	// The boot report names the steps this process really owns, taken off the
	// registration rather than off a comment (AC#4's "不许撒谎" half).
	registered := rt.RegisteredShutdownSteps()
	names := make([]string, 0, len(registered))
	for _, s := range registered {
		names = append(names, fmt.Sprintf("%d:%s", int(s), s.Name()))
	}
	fmt.Printf("wisp: %s; 任务来源：%s; 面板：%s; D38(e) steps with an owner in this process: %s\n",
		ra.residentStatusLine(), src.taskPosture(), holder.statusLine(), strings.Join(names, ", "))
	// Ticket 247 AC#6/AC#10: the capture posture gets its own line, taken off
	// the state the assembly recorded. It says which of the three shapes this
	// process is in (no collector / armed and muted / running) and it does NOT
	// claim the ball is visibly breathing: with prototypeVisuals off, which is
	// the default and this leg's ruling forbids flipping, Ball.SetAudioLevel
	// returns at once (internal/ball/liquid_windows.go:52).
	fmt.Printf("wisp: 采集腿：%s\n", raudio.posture())

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

// panelHolder is the assembly root's handle on the panel thread, and it exists
// because a dead panel thread must not be the end of this process's panel
// (ticket 319 AC#1).
//
// THE DEFECT IT KILLS, verbatim from the field: the owner closed the panel
// window at 15:41:45 on 2026-10-14; the system's close chain
// (WM_SYSCOMMAND/SC_CLOSE -> WM_CLOSE -> DestroyWindow -> WM_DESTROY ->
// PostQuitMessage(0)) made the library pump return, the panel thread's loop()
// returned, rp.finished closed - and from then on every tray/ball gesture was
// refused by post() ("panel gesture could not be handed over", 9 entries that
// afternoon) with no path anywhere that could start a thread again:
// startResidentPanel had exactly one production caller, the boot line below.
// Root cause fixed by 319-a1's census (.scratch/wisp/probes/319/a1/census.md).
//
// WHAT THIS HOLDER IS. One field of state (the current thread) behind one lock,
// and four answers the assembly root needs: the gesture executors - toggle for
// the ball's panel gestures, the panel hot key and the tray item, show for the
// second-launch activation (ticket 319 AC#2's D1 correction, leg 319-r3) - plus
// the boot-report line (statusLine) and the shutdown stop (stop). Both executors
// notice a thread whose goroutine has returned and start a fresh one on the SAME
// manager, which is what makes the window really come back.
//
// WHAT IT DOES NOT CHANGE, each named because this is the kind of edit that can
// quietly move a frozen contract:
//   - the fresh thread goes through Registry.Spawn with the same panel-sta
//     name and the same owner (startResidentPanel, untouched), so the D38b
//     resident roster - six frozen names, the panel thread deliberately not
//     among them - is not touched, and no new name appears;
//   - the D38(e) ten-step order is not touched (the holder's stop is still a
//     plain defer in runResident, in the same position as before);
//   - no new dependency edge: this file already imports internal/observe (the
//     registry it hands every leg) and internal/panel types stay out of it;
//   - the fresh goroutine's runtime.LockOSThread gets a BRAND-NEW OS thread
//     (Go destroys a locked thread when its goroutine returns), which is what
//     keeps bringUp's precondition honest: the thread about to create a window
//     has no used message queue behind it. A fresh thread on a REUSED OS thread
//     is exactly the shape 33-r7 / 33-r9 measured panicking.
type panelHolder struct {
	mu  sync.Mutex
	reg *observe.Registry
	mgr *PanelManager
	// spawn is startResidentPanel in production, and a seam in the headless nail
	// (panel_holder_319_windows_test.go) so "a second thread was really asked
	// for" is an observable fact instead of an inference from behaviour.
	spawn func(*observe.Registry, *PanelManager) *residentPanel
	rp    *residentPanel
}

// newPanelHolder builds the holder around the manager the boot assembled. A nil
// mgr is legal and means "the panel host could not be assembled at all": the
// holder then refuses by name (toggle and show false, statusLine says so out
// loud), which is the pre-319 behaviour for that boot shape, kept.
func newPanelHolder(reg *observe.Registry, mgr *PanelManager) *panelHolder {
	return &panelHolder{reg: reg, mgr: mgr, spawn: startResidentPanel}
}

// start brings up the FIRST thread - the eager step the resident boot has always
// taken, so the boot report's statusLine still describes a thread that exists.
// The window is still created only when someone asks; nothing here changes the
// "no gesture, no window" posture.
func (h *panelHolder) start() {
	if h == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.mgr == nil || h.rp != nil {
		return
	}
	h.rp = h.spawn(h.reg, h.mgr)
}

// toggle is the one executor the ball's two panel gestures, the panel hot key
// and the tray item drive. It answers exactly what residentPanel.RequestToggle
// answers, and it adds the one thing that answer never had: when the thread
// behind it is gone, a fresh thread is started for this request instead of
// refusing forever.
//
// It is deliberately NOT the executor the second-launch activation drives: a
// toggle HIDES a panel that is already on screen, and AC#2 asks for the window
// to be brought to the front. That trigger drives show below instead (ticket
// 319-r3; the reading that split the two is probes/319/v/03-ac2-live-adjudication.txt,
// shape A: "THE SECOND LAUNCH HID THE PANEL").
//
// The lock is held across the whole call on purpose: two gestures racing here
// must not start two threads, and RequestToggle never calls back into the
// holder, so the lock is a serialization point and not a deadlock.
func (h *panelHolder) toggle(via string) bool {
	if h == nil {
		return false
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.mgr == nil {
		return false
	}
	if h.rp == nil || h.rp.isFinished() {
		why := "the previous panel thread had exited"
		if h.rp == nil {
			why = "no panel thread had been started"
		}
		h.rp = h.spawn(h.reg, h.mgr)
		slog.Info("panel host: starting a fresh panel thread for this gesture",
			"via", via, "why", why)
	}
	return h.rp.RequestToggle(via)
}

// show is the executor the second-launch activation drives (ticket 319 AC#2,
// leg 319-r3; extended with the raise hop by leg 319-r4). It answers exactly
// what residentPanel.RequestShow answers - ask for the window to appear, never
// for a state flip - and it carries the same dead-thread rebuild toggle has, for
// the same reason (AC#1: a thread that ended must not make the request
// impossible). Since 319-r4 it asks for TWO requests, in this order and under
// this one lock: RequestShow (the window on screen) and then RequestRaise (the
// window in front, panel_resident_windows.go). The order is by construction -
// both posts are FIFO behind each other on the panel thread's queue - and it is
// the order AC#2's field shape needs (a window cannot be brought to the front
// before it is shown). The method answers the conjunction: if either half was
// refused by the thread, the caller sees false and the activation takes its
// tray-tip branch instead of reporting a window move that never happened.
//
// WHY THE RAISE IS ITS OWN REQUEST. RequestShow's meaning is "this window should
// be on screen"; for an ALREADY-VISIBLE panel mgr.Show re-asserts visibility
// without touching the foreground (its SetForegroundWindow runs only when the
// window was not already shown), and 319-r3 measured that this leaves the prior
// window in front - the "panel open but pushed behind, double click does
// nothing" shape of 2026-10-14 21:40:38. RequestRaise is the named second hop
// that covers exactly that shape; it does not edit mgr.Show's guard and does not
// touch any focus bookkeeping.
//
// WHY THIS IS A SIBLING COPY OF toggle AND NOT A MODE OF IT. The two route to
// two different requests on purpose: toggle's contract is a toggle, so a panel
// already on screen gets hidden (RequestToggle, panel_resident_windows.go), and
// that is correct for the ball's gestures, the hot key and the tray item. The
// second launch is not a toggle: AC#2's letter is "bring the existing instance's
// UI to the front", and the live reading behind D1 (319-v's shape A) is that
// the toggle road made the window DISAPPEAR on a double click. What the show
// road delivers per shape, measured: the cold and hidden shapes really come
// forward (mgr.Show runs ShowWindow(SW_SHOW) + SetForegroundWindow when the
// window was not already shown), the already-visible shape is raised by the
// RequestRaise hop added in 319-r4 (before it, 319-r3 measured the foreground
// staying on the prior window: probes/319/r3/02-probe-foreground-fixed.txt). The
// rebuild block is a second copy rather than a shared helper because the AC#1
// mutation teeth sit on toggle's own isFinished() term: sharing would move an
// existing pin's subject without changing anything it asserts.
func (h *panelHolder) show(via string) bool {
	if h == nil {
		return false
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.mgr == nil {
		return false
	}
	if h.rp == nil || h.rp.isFinished() {
		why := "the previous panel thread had exited"
		if h.rp == nil {
			why = "no panel thread had been started"
		}
		h.rp = h.spawn(h.reg, h.mgr)
		slog.Info("panel host: starting a fresh panel thread for this request",
			"via", via, "why", why)
	}
	tookShow := h.rp.RequestShow(via)
	tookRaise := h.rp.RequestRaise(via)
	return tookShow && tookRaise
}

// statusLine is the boot report's "does this process have a panel" answer. The
// no-manager branch returns the same sentence residentPanel's nil receiver
// returns (panelNoThreadLine, panel_resident_windows.go): the boot shape where
// the panel host never assembled must not read as though a thread were there.
func (h *panelHolder) statusLine() string {
	if h == nil {
		return panelNoThreadLine
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.mgr == nil || h.rp == nil {
		return panelNoThreadLine
	}
	return h.rp.statusLine()
}

// stop ends whichever thread is current. It is the assembly root's defer, in the
// same position the single panel.stop() defer had before ticket 319: a thread
// that was rebuilt during this process's life is the one that must be stopped
// (the earlier ones ended by themselves - that is why a rebuild happened).
func (h *panelHolder) stop() {
	if h == nil {
		return
	}
	h.mu.Lock()
	rp := h.rp
	h.mu.Unlock()
	if rp != nil {
		rp.stop()
	}
}

// secondLaunchTrayTip319 is the fallback the operator can SEE when the panel
// cannot be raised (ticket 319 AC#2's "或给机主可见的提示"). It goes through
// internal/ball's existing SetTrayTip (ball_windows.go), which posts to the ui-sta
// thread: zero new API, zero new dependency edge. The stronger shape - a real
// tray balloon (Shell_NotifyIcon + NIF_INFO) or a flashing orb - needs a new
// internal/ball method and is another ticket; what this constant buys is that a
// second launch is never silence again.
const secondLaunchTrayTip319 = "Wisp: another launch asked for the window, but this process found no panel to raise (see the wisp-dev log)"

// secondLaunchActivation319 is the function value the resident runtime runs when
// the second-instance activation event fires (ticket 319 AC#2). The runtime's
// event loop hands it over through Runtime.SetActivationHook - the same
// injection shape RegisterShutdownHook uses, so internal/proc stays a leaf and
// knows nothing about panels, balls or windows.
//
// WHY THE SHOW EXECUTOR AND NOT A TOGGLE. Bringing the window forward is not a
// capability this process lacks - Show already runs ShowWindow(SW_SHOW) +
// SetForegroundWindow (panel_host_windows.go) - it was a hop nobody had
// connected. The day's reading shows the two defects are separable: the
// 21:40:38 second launch arrived while the panel thread AND its window were
// alive and still did nothing. What the same day did NOT settle, and 319-v's
// live probes did (2026-10-15): a toggle on an already-visible panel HIDES it
// (shape A in .scratch/wisp/probes/319/v/03-ac2-live-adjudication.txt), which is
// the opposite of AC#2's letter. So the executor here is panelHolder.show: it
// ensures a live thread and asks RequestShow, the request whose meaning is
// "this window should be on screen". The tray item, the hot key and the ball's
// gestures keep driving toggle, unchanged. The in-front half for an
// already-visible panel - the one hop 319-r3 measured missing and stopped on -
// is delivered since 319-r4 by panelHolder.show's own second request,
// RequestRaise (see show's comment; mgr.Show's `if !shown` guard was left
// untouched).
//
// WHAT THE RECEIPT CAN AND CANNOT SAY. The request is asynchronous by contract
// (post() returns once the request is queued; the window verdict lands on the
// panel thread), so this sentence reports the hand-off, never the window's fate.
// When the hand-off itself was refused - no panel host, or a thread that will
// not take requests - the tray tip above is written instead, because that is the
// one shape where the double click would otherwise be silence again.
func secondLaunchActivation319(h *panelHolder, rb *residentBall) func() string {
	return func() string {
		if h.show("second-launch") {
			return "the panel executor took the request (second-launch); the window verdict lands on the panel thread's own log lines"
		}
		if rb != nil && rb.b != nil {
			rb.b.SetTrayTip(secondLaunchTrayTip319)
			return "no panel host would take the request; the tray was told to say so: " + secondLaunchTrayTip319
		}
		return "no panel host and no tray in this process; the activation was booked and cleared only"
	}
}
