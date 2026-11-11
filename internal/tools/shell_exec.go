package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/CarlosShao/wisp/internal/observe"
)

// shell.exec (C1 Tool, D34 row: PLAN.md:2574 / SPEC-07:74) - the one-shot
// half of ticket 163. It starts ONE process from the forced argv vector,
// Waits it, and reports the exit code taken from the Wait state field.
//
// AC#2's shape ruling (163-a1 §3.2) is implemented here literally: the exit
// code's channel is ProcessState.ExitCode(), retrieved where cmd.Wait
// returns - never a parse of the command's own output. The live sample a1
// measured (a process exiting 0 while PRINTING "exit 1") is pinned in
// shell_exec_live_windows_test.go as the anti-guess control: the text says
// "exit 1" and the verdict says 0, because the verdict comes from the field.
//
// Containment: the child enters its own kill-on-close tree (ShellDeps.
// NewTree, proc.OpenJobScope's shape on Windows). Closing the tree after
// Wait reaps any descendant that outlived the direct child, which is what
// keeps "cancel / end => nothing lingers" an OS fact. A call whose context
// dies first closes the tree INSTEAD of reporting an exit code: a killed
// process has no exit code worth reporting, and inventing one would be the
// exact misjudgment AC#2 forbids.

type shellExec struct{ d ShellDeps }

type shellExecArgs struct {
	Argv []string `json:"argv"`
	Cwd  string   `json:"cwd"`
}

// shellExecSchema is the C1 Parameters document: forced argv (SPEC-06 §10),
// optional cwd. There is deliberately NO string-command parameter: string
// mode is [risk] allow_shell_string's row and it is not implemented here.
var shellExecSchema = JSONSchema(`{"type":"object","properties":{` +
	`"argv":{"type":"array","items":{"type":"string"},"minItems":1,` +
	`"description":"命令与参数向量（强制 argv 形，SPEC-06 §10；不接受字符串拼接）"},` +
	`"cwd":{"type":"string","description":"可选：工作目录（经 C26 规范化并按 R2/R3 判级；省略则沿用宿主进程目录）"}}` +
	`,"required":["argv"],"additionalProperties":false}`)

func (shellExec) Name() string { return "shell.exec" }

func (shellExec) Description() string {
	return "跑一条一次性命令（argv 向量）：返回命令的 stdout/stderr 与退出码；" +
		"连续多步（先 cd 再跑、装环境）请用 shell.session，它保住工作目录和环境"
}

func (shellExec) Parameters() JSONSchema { return shellExecSchema }

// waitOutcome is what the Wait watcher goroutine reports back: the raw Wait
// result plus the captured output. The exit code is read from State at the
// one place Wait returned - AC#2's channel, stated in the type so no second
// channel can grow silently.
type waitOutcome struct {
	state *os.ProcessState
	err   error
}

func (t shellExec) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a shellExecArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	if len(a.Argv) == 0 {
		return Result{Text: "缺少 argv 参数（shell.exec 只收 argv 向量，SPEC-06 §10 强制）", IsError: true}, nil
	}
	for i, arg := range a.Argv {
		if arg == "" {
			return Result{Text: fmt.Sprintf("argv[%d] 是空串（fail-closed 拒绝）", i), IsError: true}, nil
		}
	}
	// cwd goes through the SAME C26 pipeline an fs path would: the judged
	// spelling is the opened spelling. No resolver wired = refuse that arm,
	// never "trust the raw string".
	dir := ""
	if strings.TrimSpace(a.Cwd) != "" {
		if t.d.Paths == nil {
			return Result{Text: "cwd 传了，但 C26 规范化器未接线（fail-closed：这条目录没法判，拒绝执行）", IsError: true}, nil
		}
		c, err := t.d.Paths.Canonicalize(a.Cwd)
		if err != nil {
			return Result{Text: "cwd 无法规范化（C26 拒绝）：" + cwdFailText(err), IsError: true}, nil
		}
		dir = c
	}
	if t.d.NewTree == nil {
		return Result{Text: ErrShellTreeUnwired.Error() +
			"（fail-closed：没有 kill-on-close 的进程树，拒绝起任何命令进程——取消时整棵树必须能被回收）", IsError: true}, nil
	}
	tree, err := t.d.NewTree()
	if err != nil {
		return Result{Text: "进程树建立失败：" + err.Error(), IsError: true}, nil
	}
	// The tree is closed on EVERY return path: after a clean Wait it reaps
	// descendants that outlived the child; on cancel it IS the kill.
	defer func() { _ = tree.Close() }()

	cmd := exec.Command(a.Argv[0], a.Argv[1:]...)
	cmd.Dir = dir
	var outBuf, errBuf guardedBytes
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if _, err := tree.StartInJob(cmd); err != nil {
		return Result{Text: "命令进程未能启动：" + err.Error(), IsError: true}, nil
	}
	if onUpdate != nil {
		onUpdate(fmt.Sprintf("shell.exec 已启动 %s（独立 kill-on-close 进程树）", a.Argv[0]))
	}

	done := make(chan waitOutcome, 1)
	observe.Default.Spawn(shellExecGoroutine(TaskID(ctx)), "tools", nil, func(context.Context) {
		werr := cmd.Wait()
		done <- waitOutcome{state: cmd.ProcessState, err: werr}
	})

	select {
	case w := <-done:
		code := exitCodeOf(w)
		sections := shellResultSections(outBuf.String(), errBuf.String())
		steps := []string{fmt.Sprintf("启动 %s", a.Argv[0])}
		var text string
		switch {
		case code >= 0:
			text = fmt.Sprintf("%s退出码：%d（来自进程 Wait 状态字段）", sections, code)
			steps = append(steps, fmt.Sprintf("Wait 收回退出码 %d", code))
		default:
			// Wait failed with something that is not a process exit: there is
			// no command verdict at all, so none is claimed.
			text = fmt.Sprintf("%s退出码读不回来（进程 Wait 未给出可认的退出码：%v）", sections, w.err)
			steps = append(steps, "Wait 未给出可认的退出码")
		}
		return Result{Text: text, IsError: code != 0, AppliedSteps: steps}, nil
	case <-ctx.Done():
		// The tree close terminates the direct child and every descendant;
		// the Wait watcher drains on its own. NO exit code is reported: the
		// process was killed by us, and booking kill's "1" as the command's
		// answer would be the misjudgment AC#2 forbids.
		return Result{Text: fmt.Sprintf(
			"命令未在 C22 预算内回来（预算 %s），整棵进程树已被 kill-on-close Job 回收；"+
				"不发退出码（被杀进程没有可认的退出码），也没有输出可归属——"+
				"需要更长时间或后台跑的活请改用 task.spawn。",
			shellTimeout), IsError: true, AppliedSteps: []string{
			fmt.Sprintf("启动 %s", a.Argv[0]),
			"C22 预算到点，关闭进程树回收整棵后代",
		}}, nil
	}
}

// exitCodeOf reads the code from the Wait outcome. A clean exit answers 0;
// an ExitError answers ExitCode(); anything else (a start fault reaching
// Wait, an internal drain error) has no code and answers -1, which the
// caller renders as an internal refusal rather than a command result.
func exitCodeOf(w waitOutcome) int {
	if w.err == nil {
		return 0
	}
	var ee *exec.ExitError
	if errors.As(w.err, &ee) && ee.ProcessState != nil {
		return ee.ExitCode()
	}
	return -1
}

// shellResultSections renders captured stdout and stderr. Both are
// byte-copies of what the child wrote (os/exec drains them on its own
// goroutines); ordering between the two streams is not promised here.
func shellResultSections(stdout, stderr string) string {
	var b strings.Builder
	b.WriteString(strings.TrimRight(stdout, "\r\n"))
	b.WriteString("\n")
	if stderr != "" {
		b.WriteString("[stderr]\n")
		b.WriteString(strings.TrimRight(stderr, "\r\n"))
		b.WriteString("\n")
	}
	return b.String()
}

func cwdFailText(err error) string {
	msg := err.Error()
	if len(msg) > 200 {
		msg = msg[:200]
	}
	return msg
}

// shellExecGoroutine names the Wait watcher on the D38b per-task roster
// (tool-exec- prefix): the watchdog sees one picture.
func shellExecGoroutine(taskID string) string {
	if taskID == "" {
		return "tool-exec-shell-exec"
	}
	return "tool-exec-shell-" + taskID
}

// guardedBytes is a concurrency-safe bytes buffer: os/exec's own drain
// goroutines write into it while the calling goroutine reads it after Wait.
type guardedBytes struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (g *guardedBytes) Write(p []byte) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.buf.Write(p)
}

func (g *guardedBytes) String() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.buf.String()
}

// takeAll returns and clears the buffer (the session carrier drains cmd.exe's
// diagnostics stream after every command through this method).
func (g *guardedBytes) takeAll() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := g.buf.String()
	g.buf.Reset()
	return out
}
