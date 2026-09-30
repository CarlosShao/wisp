//go:build windows

package main

// The resident leg's source of cards (ticket 246 AC#7).
//
// WHY THIS FILE EXISTS. Ticket 246's title sentence is "the process that runs
// tasks is not the process that holds the ball, so no veto channel has an
// executor". The first four legs moved the door into the resident process:
// resident_windows.go:122 builds the gate, :135 starts the ball with
// ra.vetoByEsc, :144 binds the two, :156 registers D38(e) step 3's producer. What
// none of them did was produce a card: the orchestrator's own re-measurement at
// 22:5x found AskOnTaskRoot (resident_approval_windows.go:261) and
// askConfirmation (:205) with ZERO product callers - only this package's own
// cases call them. So the gate was real, the card route was real, Esc really
// could veto, the exit really did refuse - and nothing ever asked. AC#7 is that
// missing hop, and this file is it.
//
// WHAT A TASK SOURCE IS HERE: the operator's keyboard, and nothing else.
// interactiveStdin() (approval_reply_stdin_windows.go:41) is the ONE gate on it -
// ruling 2.1 forbids a second "is this a console" test, and the reason is in that
// file's header: a stream a script can pre-fill is not "the user confirmed", it is
// a way to make confirmation happen without a user. A task text is MORE dangerous
// than an answer, because it drives a real tool call, so it gets the same gate and
// therefore the same verdicts: a pipe, a redirected file, Explorer's detached
// start and CI all get nil, and this process then says out loud that its task
// entry is not enabled and carries on hosting the ball.
//
// A MISSING CONSOLE IS NEVER A REASON TO STOP. Ticket 128 settled that the only
// fatal condition in this leg is "no data root", and SPEC-05 §3.4 asks for a
// classified and VISIBLE downgrade, not a hard stop. So every failure below - no
// console, no config.toml, no credential, a store that will not open - prints a
// named line and leaves the ball running.
//
// ONE GATE PER PROCESS (ruling 2.2). The pipeline this file starts is assembled by
// cmd/wisp/run.go's assembleRuntime, the only assembly root in the repository,
// with ra's gate, ra's presentation surface and ra's Replies ledger handed IN
// through runSpec. Nothing here calls approval.New: a second gate in this process
// would be a second "who is waiting on a decision" truth source, which is exactly
// the shape 246-a1 §5 reason 1 and ledger A481 refused. The consequence of that
// ordering is printed at boot rather than smoothed over - a gate built before the
// session ledger exists has no Grants writer, so 「本会话内允许」 releases the call
// and books GRANT-DROPPED, which is approval.Gate's own documented posture.
//
// THE TASK ROOT (ruling 2.3). runSpec.taskCtx is ra.root, the context D38(e) step
// 3 cancels. Without that line "leaving cancels the running task" would still be a
// sentence with no mechanism behind it: the loop's ctx would hang off
// context.Background() the way the CLI leg's does, and an exit signal would refuse
// the cards while the task kept running to its own LLM timeout.
//
// WHAT THIS FILE DOES NOT CLAIM.
//   - Not that a card is visible. This process links no WebView2 host
//     (approval_always.go:165), so the reading is the orb's state, the ledger, the
//     audit line and what the desktop does with the cancel key - the same limit
//     ticket 246 AC#2 wrote down.
//   - Not a panel, not a microphone, not the tray menu (ruling 2.5). The panel's
//     Message slot (panel_inbound.go:238) is still nil and stays nil; AC#6 settled
//     that of D43's four veto channels only Esc lands in this ticket.
//   - Not "the model answered": the provider this pipeline talks to is whatever
//     config.toml names. The credential assertion is the owner's to make, so the
//     four-step reading this leg can hand over is split in the evidence table
//     between the injected seam and the real endpoint.

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/proc"
)

// The four postures of the task entry, as constants because a case in
// cmd/wisp/resident_task_source_246_windows_test.go reads each of them and the
// honest/lying split of the boot report lives in their wording.
const (
	// taskEntryDisabledClaim is printed when interactiveStdin() said "no console"
	// and no test injection was accepted. It names the consequence, not just the
	// condition: a leg that cannot be handed a task also has nothing that asks.
	taskEntryDisabledClaim = "任务入口未启用（本机没有可交互控制台）"
	// taskEntryConsoleClaim is the other branch: a real terminal this process owns.
	taskEntryConsoleClaim = "任务入口已启用（控制台键盘）"
	// taskEntryInjectedClaim names the test-only entry, and says in the same
	// breath that it is not the console path.
	taskEntryInjectedClaim = "任务入口经测试注入位打开"
	// taskEntryRefusedClaim is the sentence the injection gets when the predicate
	// is not met. It is printed even when the console path is up, because a
	// refused injection must never be a quiet drop.
	taskEntryRefusedClaim = "这条注入位在当前环境不被受理"
	// pipelineAbsentClaim is the loud half of a degraded assembly: the gate is
	// there, the entry is there, and the stack that would raise a card is not.
	pipelineAbsentClaim = "任务管线未装配"
	// taskPostureConsole / taskPostureInjected / taskPostureAbsent /
	// taskPostureRefused are the four forms the boot report's 任务来源 field can
	// take. They are one per branch of startResidentTaskSource's return, so the
	// report can tell "no entry" apart from "entry refused to assemble" apart from
	// "entry is a test injection" - the three of them used to print the same
	// nothing, which is the sentence AC#7 exists to stop.
	taskPostureConsole  = "控制台键盘（已接入）"
	taskPostureInjected = "仅本发的测试注入位"
	taskPostureAbsent   = "无（任务入口未启用，理由见上面那行）"
	taskPostureRefused  = "无（任务管线未装配，理由见上面那行）"
)

// testTaskTextEnv is ticket 246 AC#7's narrow test escape hatch, and the whole of
// its remit: it supplies ONE task text, in ONE environment, when the harness
// itself named the data root. It does not open the reply stream (answers still
// need a real console), it is not read in dev or prod, and it is not read when the
// data root came from anywhere but WISP_TEST_DATA_DIR.
//
// The reason it exists at all is measured, not assumed: a subprocess started by a
// test has no console (its stdin is a pipe or nil), so the four steps AC#7 asks
// for could otherwise never be driven from a shipped wisp.exe. AGENTS.md §1.3
// allows the C5 golden-SSE seam and the CLI seam; this env var is the CLI side of
// that pair, and the evidence table states which readings it therefore licenses.
const testTaskTextEnv = "WISP_TEST_TASK_TEXT"

// taskSourceGoroutinePrefix is the D38b per-task roster prefix (observe's
// "agent-task-<id>") the task goroutines are booked under. The suffix is this
// leg's own sequence number, so two tasks in one process are two roster entries.
const taskSourceGoroutinePrefix = "agent-task-resident-"

// taskSourceConsoleGoroutine is the console reader's roster name: the per-task
// "approval-waiter" entry, which is what every other answer side in this package
// already uses (approval_reply.go:454, approval_always.go:99).
const taskSourceConsoleGoroutine = "approval-waiter"

// residentTaskSource is the resident leg's task entry plus the pipeline it feeds.
// nil is a legal return from startResidentTaskSource and means "this process has
// no source of cards", which every sentence built from it says out loud.
type residentTaskSource struct {
	ra  *residentApproval
	rt  *proc.Runtime
	run *agentRuntime

	// root owns this leg's console reader and every task goroutine. It derives
	// from ra.root, so D38(e) step 3's cancel reaches both (ruling 2.3).
	root   *observe.Root
	handle *observe.Handle

	mu sync.Mutex
	// closed is D38(e) step 1's door: "no new tasks". It is set before the
	// sequence cancels anything, so an exit can never be followed by a fresh
	// question.
	closed  bool
	running bool
	// taskDone closes when the task that set running leaves execute(). drainAnd-
	// Close waits on it rather than guessing how long a task takes.
	taskDone chan struct{}
	// seq numbers the per-task roster entries.
	seq int
	// surface is the answer side, built from approval_reply.go's own verb table.
	// nil means nobody can type an answer, which is the no-console branch.
	surface *replySurface
	// consoleStop ends the console reader's verb handling without killing the
	// process (the `quit` verb's route).
	consoleStop func()
}

// residentTestTaskText reads the narrow escape hatch. It returns the text ONLY
// when all three of these hold, and a reason string whenever the variable was set
// and did not hold:
//
//   - the environment is test (buildinfo.ResolveEnv's answer, the same one the
//     data root and the mutex rules were picked with);
//   - WISP_TEST_DATA_DIR is set, i.e. the harness named this process's data root;
//   - the layout's data dir IS that string (internal/proc/envfork.go makes it an
//     identity contract: it comes back verbatim, so this is a plain comparison and
//     never a filepath decision - no Clean, no Abs, no risk.PathResolver bypass).
func residentTestTaskText(env buildinfo.Env, layout proc.Layout) (text, refusal string) {
	raw := strings.TrimSpace(os.Getenv(testTaskTextEnv))
	if raw == "" {
		return "", ""
	}
	injected := strings.TrimSpace(os.Getenv(proc.TestDataDirEnv))
	switch {
	case env != buildinfo.EnvTest:
		return "", fmt.Sprintf("%s=%q 已设置，但本进程的环境是 %s：%s（只有 test 受理任务文本注入）",
			testTaskTextEnv, redactTaskTextForLine(raw), env, taskEntryRefusedClaim)
	case injected == "":
		return "", fmt.Sprintf("%s=%q 已设置，但本进程的数据根不是台件指定的那一枚（没有 %s）：%s",
			testTaskTextEnv, redactTaskTextForLine(raw), proc.TestDataDirEnv, taskEntryRefusedClaim)
	case !strings.EqualFold(layout.DataDir, injected):
		return "", fmt.Sprintf("%s=%q 已设置，但本进程的数据根 %q 与 %s=%q 不是同一枚目录：%s",
			testTaskTextEnv, redactTaskTextForLine(raw), layout.DataDir,
			proc.TestDataDirEnv, injected, taskEntryRefusedClaim)
	}
	return raw, ""
}

// redactTaskTextForLine keeps a refused injection readable without putting the
// whole task text into the log: the record says that a value of that length was
// rejected, not what it asked for. (Accepted values are handed to the loop, which
// books the query text into task_log on its own.)
func redactTaskTextForLine(s string) string {
	r := []rune(s)
	if len(r) <= 12 {
		return s
	}
	return string(r[:12]) + "…"
}

// startResidentTaskSource is the hop AC#7 measured as missing: it decides whether
// this process has a task entry at all, assembles the pipeline through the ONE
// assembly root with ra's gate injected, and puts a reader on that entry. Every
// return of nil is printed by the path that produced it.
func startResidentTaskSource(rt *proc.Runtime, ra *residentApproval) *residentTaskSource {
	// The single console gate (ruling 2.1). No second test of "is this a
	// terminal" may appear anywhere in this file.
	console := interactiveStdin()
	injectedText, injectedRefusal := residentTestTaskText(rt.Env, rt.Layout)
	if injectedRefusal != "" {
		slog.Warn("task source: " + injectedRefusal)
		fmt.Printf("wisp: %s\n", injectedRefusal)
	}
	if console == nil && injectedText == "" {
		slog.Warn("task source: 任务入口未启用",
			"why", "interactiveStdin() 返回 nil：标准输入不是本进程拥有的控制台输入缓冲",
			"effect", "没有东西会去举一张确认卡片；球、托盘与退出序列照旧")
		fmt.Printf("wisp: %s：本进程仍然带球常驻、仍然会在退出时拒绝挂起的卡片，但没有任何东西会去举一张卡\n",
			taskEntryDisabledClaim)
		return nil
	}

	src := &residentTaskSource{ra: ra, rt: rt}
	// The injection is the whole of ruling 2.2: the gate, the surface it shows
	// cards on and the ledger those cards were booked into all come from ra, and
	// assembleRuntime builds everything else around them. reply stays nil - this
	// file owns the console stream and routes the answer verbs itself, because a
	// second reader on one keyboard is not two listeners, it is a race.
	spec := runSpec{
		stdout:  os.Stdout,
		stderr:  os.Stderr,
		dataDir: rt.Layout.DataDir,
		gate:    ra.gate,
		ui:      ra.ui,
		cards:   ra.cards,
		taskCtx: ra.root,
	}
	run, code := assembleRuntime(spec)
	if code != 0 {
		// Classified and visible, and the boot goes on: SPEC-05 §3.4 forbids the
		// silent downgrade, not the degraded one. Nothing was opened, so there is
		// nothing to close except a half-built runtime that did reach a store.
		slog.Error("task source: 任务管线未装配", "exit_code", code,
			"why", "装配根在这一发返回了非零分类；审批门仍在，但没有会举卡的任务环路")
		fmt.Printf("wisp: %s（退出码 %d，原因见上面这几行）：审批门已装配而任务管线没有，本进程仍然带球常驻\n",
			pipelineAbsentClaim, code)
		if run != nil {
			run.close()
		}
		// Non-nil with run == nil: the entry WAS offered and the pipeline refused,
		// which the boot report then says differently from "no console at all".
		return src
	}
	src.run = run
	src.root = observe.NewRootFrom(ra.root, "resident-task-source")
	// The reply side, on THIS process's root: rt.close() cancels rt.replyRoot, and
	// rebindLedger is false because ra attached its ledger to ra.gate at boot with
	// the resident transport names. Re-attaching would restamp a resident card as
	// a console card in the audit line.
	if console != nil {
		src.surface = run.newReplySurface(src.root, ra.consoleVetoChannel(), false)
		src.consoleStop = func() { src.root.Cancel() }
	}

	// Two D38(e) steps that had no producer in this process before the entry
	// existed. Neither touches the frozen ORDER - RunShutdownSequence still walks
	// 1..10 and a step nobody handed is still honestly skipped:
	//
	//   - step 1 "TaskScheduler closes its door (no new tasks)": that door is
	//     src.closed. Without it a Ctrl+C would refuse the hanging cards in step 3
	//     while the console kept starting new tasks behind them.
	//   - step 7 "flush logs -> db-writer drained -> wal_checkpoint -> SQLite
	//     closed": the store this assembly opened belongs there, and closing it
	//     while a task is still writing rows would be the out-of-order incident
	//     shutdown.go's header names as an incident rather than a style choice.
	for _, h := range []struct {
		step proc.ShutdownStep
		fn   func(context.Context) error
		name string
	}{
		{proc.StepSchedulerClose, src.closeDoor, "第 1 步（关闭任务门）"},
		{proc.StepFlushAndCloseDB, src.drainAndClose, "第 7 步（冲刷并关闭存储）"},
	} {
		if err := rt.RegisterShutdownHook(h.step, h.fn); err != nil {
			slog.Error("task source: 退出序列步骤未能注册为钩子", "step", int(h.step), "err", err)
			fmt.Printf("wisp: %s未注册（%v）：这一步会被记为 skipped，任务源不会被有序关闭\n", h.name, err)
		}
	}

	if console != nil {
		src.handle = observe.Default.Spawn(taskSourceConsoleGoroutine, "agent", src.root,
			func(ctx context.Context) { src.runConsoleLoop(ctx, console, os.Stdout) })
		fmt.Printf("wisp: %s：%s\n", taskEntryConsoleClaim, consoleEntryHelpLine)
	}
	if injectedText != "" {
		fmt.Printf("wisp: %s（%s）：这是只在本机测试台件里受理的一发任务文本，本机没有可交互控制台，"+
			"挂起的卡片只能由取消键否决\n", taskEntryInjectedClaim, testTaskTextEnv)
		if err := src.submitTask(injectedText); err != nil {
			slog.Error("task source: 注入的任务没有起来", "err", err)
			fmt.Printf("wisp: 任务没有起来（%v）\n", err)
		}
	}
	return src
}

// consoleEntryHelpLine is the one place the resident leg's console grammar is
// written down. It is printed at boot and repeated by an unknown verb, so what a
// user is told and what the reader accepts cannot drift apart.
const consoleEntryHelpLine = "输入 task <文本> 起一发任务；卡片上给了编号，" +
	"yes / session / no / veto / always / head / view 是答复那一张的路，quit 关掉这一路"

// runConsoleLoop reads the operator's lines for as long as this process owns the
// console. One reader carries BOTH halves of the stream on purpose: a second
// reader on one keyboard is a race, and the task entry and the answer entry are
// the same physical channel.
//
// The verb table is not copied here - every line that is not this loop's own two
// verbs goes through replySurface.handle, the same funnel approval_reply.go's
// console listener uses, so a resident answer and a `wisp run` answer cannot grow
// different leniencies.
func (src *residentTaskSource) runConsoleLoop(ctx context.Context, in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Buffer(make([]byte, 0, 4096), 1<<16)
	for sc.Scan() {
		if ctx.Err() != nil {
			return
		}
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		verb, rest, _ := strings.Cut(line, " ")
		switch verb {
		case "task":
			// The whole remainder is the task text - a sentence with spaces in it
			// is the normal case here, which is why this branch does not go
			// through splitReplyLine (that grammar has one word for the card id
			// and one for its argument, and a task is neither).
			if err := src.submitTask(strings.TrimSpace(rest)); err != nil {
				fmt.Fprintf(out, "wisp: 任务没有起来（%v）\n", err)
			}
			continue
		case "quit":
			// Ending THIS route is all quit means here: the cards stay pending and
			// resolve on their own clocks, exactly like the console listener's quit,
			// and the process still leaves through a signal, which is the state
			// ticket 228 recorded for the tray's Exit item too.
			fmt.Fprintf(out, "wisp: 任务入口已关闭（卡片仍在等待，未答复的那一张按各自的超时与退出序列处理）\n")
			src.closeDoorOnly()
			return
		default:
			if src.surface == nil {
				fmt.Fprintf(out, "wisp: 这一路没有答复通道（装配时读不到可交互控制台）：%s\n", line)
				continue
			}
			_, corr, arg := splitReplyLine(line)
			text, err := src.surface.handle(verb, corr, arg)
			if err != nil {
				fmt.Fprintf(out, "wisp: 答复未被接受（%s）：%v\n", line, err)
				fmt.Fprintf(out, "wisp: %s\n", consoleEntryHelpLine)
				continue
			}
			fmt.Fprintf(out, "wisp: %s\n", text)
		}
	}
	if err := sc.Err(); err != nil && !errors.Is(err, io.EOF) {
		slog.Warn("task source: 控制台输入流断了", "err", err)
		fmt.Fprintf(os.Stdout, "wisp: 任务输入流已断（%v），这一发之后无人能再提交\n", err)
	}
}

// submitTask starts ONE task on the pipeline this file assembled. The single-slot
// rule is a decision, not a limitation to hide: the loop's per-task budgets, the
// panel roster and the L1 window were all designed around one operator waiting on
// one orb, and a queue of tasks in a leg that cannot show a page would be a queue
// nobody can read. A second line while a task runs is refused out loud.
func (src *residentTaskSource) submitTask(text string) error {
	text = strings.TrimSpace(text)
	if text == "" {
		return errors.New("task 后面要跟任务文本")
	}
	src.mu.Lock()
	if src.closed {
		src.mu.Unlock()
		return errors.New("任务入口已随退出序列关闭（D38e 第 1 步），本进程不再接新任务")
	}
	if src.running {
		src.mu.Unlock()
		return errors.New("本进程此刻已有一发任务在跑：常驻腿一次只接一发，等它结束再提交")
	}
	src.running = true
	src.seq++
	seq := src.seq
	src.taskDone = make(chan struct{})
	done := src.taskDone
	src.mu.Unlock()

	root := observe.NewRootFrom(src.root.Ctx, fmt.Sprintf("resident-task-%d", seq))
	// The roster name is the per-task "agent-task-<id>" shape (D38b), and the
	// recover boundary is Registry.run's, never a local recover (D22 ban #1).
	observe.Default.Spawn(taskSourceGoroutinePrefix+fmt.Sprint(seq), "agent", root,
		func(ctx context.Context) {
			defer close(done)
			// root.Cancel on the way out: this task's own ctx is what execute
			// derives its deadline from (ruling 2.3), and cancelling it here is
			// what makes a finished task stop holding a live subtree while the
			// entry that produced it stays up.
			defer root.Cancel()
			code := src.run.execute(root.Ctx, text)
			src.noteTaskFinished(code)
		})
	slog.Info("task source: 常驻腿提交一发任务",
		"seq", seq, "task_chars", len([]rune(text)), "console_entry", src.surface != nil)
	return nil
}

// noteTaskFinished frees the single slot. The exit code is booked, not swallowed:
// a resident leg has no exit code of its own to carry it in.
func (src *residentTaskSource) noteTaskFinished(code int) {
	src.mu.Lock()
	src.running = false
	src.mu.Unlock()
	if code != 0 {
		slog.Warn("task source: 任务没有正常收口", "exit_code", code)
		fmt.Printf("wisp: 这一发任务以退出码 %d 收口（分类见上面那行 wisp run: 状态）\n", code)
		return
	}
	fmt.Printf("wisp: 这一发任务已完成\n")
}

// closeDoor is D38(e) step 1: no new tasks. It does NOT cancel the running task -
// that is step 3's job through ra.root, and doing it here would put the
// cancellation ahead of the refusal the sequence owes the hanging cards.
func (src *residentTaskSource) closeDoor(ctx context.Context) error {
	src.closeDoorOnly()
	slog.Info("task source: 退出第 1 步完成：任务门已关闭，不再接新任务")
	return nil
}

func (src *residentTaskSource) closeDoorOnly() {
	src.mu.Lock()
	src.closed = true
	src.mu.Unlock()
	if src.consoleStop != nil {
		src.consoleStop()
	}
}

// drainAndClose is D38(e) step 7's producer for this process: wait, bounded, for
// the task that step 3 already cancelled, then close the store it writes rows
// into. The bound is a context deadline (D42#9 - no wall-clock difference decides
// anything), and exceeding it means the store is NOT closed: a closed SQLite file
// under a still-writing task is a worse record than one the OS releases.
func (src *residentTaskSource) drainAndClose(ctx context.Context) error {
	src.mu.Lock()
	waiting := src.running
	done := src.taskDone
	src.mu.Unlock()
	if waiting && done != nil {
		waitCtx, cancel := context.WithTimeout(ctx, proc.TaskWaitTimeout)
		defer cancel()
		select {
		case <-done:
		case <-waitCtx.Done():
			slog.Warn("task source: 退出第 7 步到点仍有任务在写，本轮不关闭存储（交给进程退出）")
			fmt.Printf("wisp: 任务在退出预算内没有收口，本进程不关闭存储句柄\n")
			return nil
		}
	}
	src.run.close()
	slog.Info("task source: 退出第 7 步完成：日志与存储已冲刷关闭")
	return nil
}

// taskPosture is the boot report's 任务来源 field - the half that AC#7 exists to
// stop this process from lying about. It reads the same three facts the entry
// function used, so the sentence and the wiring are one statement.
func (src *residentTaskSource) taskPosture() string {
	switch {
	case src == nil:
		return taskPostureAbsent
	case src.run == nil:
		return taskPostureRefused
	case src.surface != nil:
		return taskPostureConsole
	default:
		return taskPostureInjected
	}
}

// consoleVetoChannel is the answer side's own statement of which of SPEC-06 §2's
// four channels this process really owns. It is conditional on bindBallHost
// having loaded Esc behind a real window, because approval_reply.go's rule is that
// declaring a transport IS asserting it is up: a resident leg with no ball may not
// hand its console a veto route it cannot fire.
func (ra *residentApproval) consoleVetoChannel() approval.Channel {
	ra.mu.Lock()
	loaded := ra.escLoad
	ra.mu.Unlock()
	if !loaded {
		return ""
	}
	return approval.ChannelEsc
}
