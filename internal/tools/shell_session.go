package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/CarlosShao/wisp/internal/observe"
)

// shell.session (C1 Tool, D34 row: PLAN.md:2575 / SPEC-07:75) - the
// persistent half of ticket 163: ONE live shell process per task, driven
// one command per call, so the second command sees the first one's working
// directory and environment (the gap AC#1 measured: two one-shot spawns
// print different directories).
//
// The shape the frozen row pins:
//
//   - 一次性外观、常驻内核: from the model's side every call is an ordinary
//     tool call with its own C19 verdict and its own C22 budget; underneath,
//     one cmd.exe stays alive for the task's whole scope.
//   - 不引 PTY: the carrier is pipes (SessionShell), not a terminal. The dsh
//     ConPTY half - "Windows ConPTY remains outside Job containment", the
//     node-pty kill() throw, the \x03 SIGINT poke, the taskkill sweep - is
//     the "don't copy" list (ticket :22), and the in-repo answer to all of
//     it is the kill-on-close Job the carrier already sits in.
//   - 会话寿命绑任务作用域: sessions are keyed by the host's task id; the
//     task-end boundary (Bridge.CloseTask -> Options.TaskCleanup) is the
//     ONLY collector, so a task that ends or is cancelled takes its shell
//     with it (AC#3; the teeth test also proves a foreign task's end cannot
//     reach someone else's session).
//
// AC#2's exit-code protocol (the marker, dsh form): each call writes
//
//	<command line>
//	echo <random marker> %errorlevel%
//
// and reads the stream until a line carrying THAT random marker appears;
// the number after it is the command's exit code. The channel is NOT a
// wall-clock guess (no time arithmetic anywhere below) and NOT a parse of
// the command's own content (the a1 live sample - a process that prints
// "exit 1" while exiting 0 - cannot mislead it, because the marker is
// per-call random and only the host writes it). The answer to AC#2's
// question "摘掉标记这一味，是否存在一发超时从此判错" is YES, and the code is
// built so that world is unreachable: with the marker gone there is no
// boundary in the stream, so a timed-out call's output and the NEXT call's
// output are indistinguishable - that is the misjudgment. With the marker,
// a timeout is a DETECTED desync: the timeout path closes the kill-on-close
// tree, reaps the shell, drops everything buffered, and reports NO exit
// code - the next shell.session call starts a fresh session. The mutation
// arm in shell_session_test.go pins the negative: a desynced stream without
// its marker never becomes a result, least of all an exit code.

// sessionPendingCap bounds ONE command's undelivered output. A command that
// streams faster than it finishes (think `yes`) would otherwise grow the
// buffer for the whole C22 budget; past the cap the session is reset and the
// call refuses - detached high-volume work belongs to task.spawn.
const sessionPendingCap = 64 << 20

// ---------------------------------------------------------------------------
// the task-scoped roster
// ---------------------------------------------------------------------------

// ShellSessions is the task-scoped session table. One live session per task
// id: a task's commands run in order inside one shell, which is the entire
// value the D34 row buys ("在一条活着的会话里跑下一条命令"). The table is shared
// with the bridge's task-end janitor (Options.TaskCleanup), so host and tool
// always look at the SAME table - the TaskDeps.Roster shape.
type ShellSessions struct {
	mu     sync.Mutex
	byTask map[string]*shellSessionState
}

// NewShellSessions returns an empty session roster.
func NewShellSessions() *ShellSessions {
	return &ShellSessions{byTask: map[string]*shellSessionState{}}
}

// sessionFor returns the task's session row, creating the empty row if the
// task has none (the carrier itself is created lazily on first command).
func (r *ShellSessions) sessionFor(taskID string) *shellSessionState {
	if r == nil || taskID == "" {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.byTask == nil {
		r.byTask = map[string]*shellSessionState{}
	}
	if s, ok := r.byTask[taskID]; ok {
		return s
	}
	s := &shellSessionState{taskID: taskID}
	r.byTask[taskID] = s
	return s
}

// TerminateForTask collects the task's session and kills its process tree
// (kill-on-close Job close: the shell and every descendant die with the
// handle). It returns how many sessions were collected - the bridge's
// task-end janitor calls this, and it is the ONLY collector: there is no
// by-name or by-age sweep, because a session outliving its task is a
// defect, not a lifecycle option (AC#3).
func (r *ShellSessions) TerminateForTask(taskID string) int {
	if r == nil || taskID == "" {
		return 0
	}
	r.mu.Lock()
	s := r.byTask[taskID]
	delete(r.byTask, taskID)
	r.mu.Unlock()
	if s == nil {
		return 0
	}
	s.kill()
	return 1
}

// Count reports how many tasks currently hold a session row (host
// diagnostics and the AC#3 teeth test's readout).
func (r *ShellSessions) Count() int {
	if r == nil {
		return 0
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.byTask)
}

// ---------------------------------------------------------------------------
// the session row
// ---------------------------------------------------------------------------

type shellSessionState struct {
	taskID string

	// mu serializes commands (the D38d ceiling still lets two calls of the
	// same task race; the session defines what serialization MEANS here:
	// they queue, in order, like commands typed at one prompt).
	mu sync.Mutex
	// carrier is the live shell. nil = no session (fresh task, or after a
	// reset); the next command builds a new one.
	carrier SessionShell
	// pending is the raw output stream accumulated by the reader goroutine.
	pending guardedText
	// streamOver is closed by the reader when the shell's stdout ended
	// (process exit or killed pipe).
	streamOver chan struct{}
	logf       func(format string, args ...any)
}

func (s *shellSessionState) log(format string, args ...any) {
	if s.logf != nil {
		s.logf("shell.session task=%s "+format, args...)
	}
}

// kill collects one session: close the kill-on-close tree (the OS
// terminates the shell and every descendant), reap the shell, drop the
// carrier. The buffered output dies with the session - nothing of a dead
// session can be misattributed to a later command (AC#2's answer,
// operationalized).
func (s *shellSessionState) kill() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dropLocked()
	s.log("session collected (tree closed, shell reaped)")
}

// dropLocked releases the carrier under the row lock. The caller holds mu.
func (s *shellSessionState) dropLocked() {
	if s.carrier == nil {
		return
	}
	_ = s.carrier.Close()
	_, _ = s.carrier.Wait()
	s.carrier = nil
}

// ensureLocked builds the carrier if this row has none. cwd is honored only
// at creation - afterwards the shell's own `cd` is the working directory,
// and a cwd parameter on a live session is a refusal, never a silent no-op.
func (s *shellSessionState) ensureLocked(d ShellDeps, cwd string) (SessionShell, error) {
	if s.carrier != nil {
		if cwd != "" {
			return nil, errSessionCwdOnLive
		}
		return s.carrier, nil
	}
	var c SessionShell
	var err error
	if d.NewSessionShell != nil {
		c, err = d.NewSessionShell()
	} else {
		c, err = newCmdSessionShell(d, cwd)
	}
	if err != nil {
		return nil, err
	}
	s.carrier = c
	s.pending = newGuardedText()
	s.streamOver = make(chan struct{})
	s.logf = d.Logf
	s.log("session created (cwd=%q, kill-on-close tree assigned)", cwd)
	s.spawnReader(c)
	return c, nil
}

// spawnReader starts the session's ONE output reader on the D38b per-task
// roster ("tool-exec-shell-<task id>"): it appends raw bytes to the guarded
// buffer and ends when the stream ends. The watchdog sees one picture; no
// bare goroutine exists anywhere in this family (D22 ban #1).
func (s *shellSessionState) spawnReader(c SessionShell) {
	done := s.streamOver
	pending := &s.pending
	observe.Default.Spawn(shellExecGoroutine(s.taskID), "tools", nil, func(context.Context) {
		defer close(done)
		buf := make([]byte, 8192)
		r := c.Stdout()
		for {
			n, err := r.Read(buf)
			if n > 0 {
				pending.append(buf[:n])
			}
			if err != nil {
				return
			}
		}
	})
}

// sessionRun is one command's outcome inside the session.
type sessionRun struct {
	// output is the command's own stream, with the marker line stripped.
	output string
	// code is the exit code: from the marker line (a live shell) or from the
	// shell process's Wait state field (a shell that exited). Negative only
	// on paths that report no code at all.
	code int
	// fromWait records which channel produced code (audit text only).
	fromWait bool
	// ended: the shell process exited (the session is over either way).
	ended bool
	// timeout: the call's budget ran out without a boundary in the stream.
	timeout bool
	// reset: the carrier was dropped and the next call starts fresh.
	reset bool
	// errText: a refusal that is not a command result (unwired, bad line...).
	errText string
}

// run drives ONE command through the live shell. The caller holds s.mu.
func (s *shellSessionState) run(ctx context.Context, d ShellDeps, argv []string, cwd string) sessionRun {
	carrier, err := s.ensureLocked(d, cwd)
	if err != nil {
		if err == errSessionCwdOnLive {
			return sessionRun{errText: "会话已在运行：cwd 参数只在会话建立的那一发生效；" +
				"换工作目录请发一条 cd /d 命令（那也是会话存在的意义——目录和环境都是延续的）"}
		}
		return sessionRun{errText: "会话进程建立失败：" + err.Error()}
	}
	line, err := joinSessionArgv(argv)
	if err != nil {
		return sessionRun{errText: err.Error()}
	}
	marker, err := newSessionMarker()
	if err != nil {
		return sessionRun{errText: err.Error()}
	}

	// The command line plus the marker echo line: the boundary in the stream
	// every read below is looking for. This is host machinery - the judged
	// argv is `line` alone; the appended redirect and echo never pass C19,
	// and nothing the MODEL sent is modified by them.
	stdin := carrier.Stdin()
	if _, err := fmt.Fprintf(stdin, "%s 2>&1\r\necho %s %%errorlevel%%\r\n", line, marker); err != nil {
		// The shell died between commands (its own `exit`, a crash). Reap it
		// and say so; the next call starts fresh.
		_, werr := carrier.Wait()
		s.dropLocked()
		return sessionRun{reset: true, errText: fmt.Sprintf(
			"会话进程已失联（写入命令失败：%v；收尸时的 Wait：%v）。会话已回收，"+
				"下一发 shell.session 会开全新会话。", err, werr)}
	}

	// have accumulates every byte this command produced: a marker line can
	// split across two read chunks, so the search runs over the whole
	// accumulated stream, never over one take.
	var have strings.Builder
	for {
		have.WriteString(s.pending.takeAll())
		pending := have.String()
		if code, idx, ok := findMarkerCode(pending, marker); ok {
			out := strings.TrimRight(pending[:idx], "\r\n")
			if diag := carrier.TakeStderr(); diag != "" {
				out += "\n[cmd 诊断]\n" + strings.TrimRight(diag, "\r\n")
			}
			return sessionRun{output: out, code: code}
		}
		if len(pending) > sessionPendingCap {
			s.dropLocked()
			return sessionRun{reset: true, errText: fmt.Sprintf(
				"这条命令的输出超过单发上限（%d 字节）还没结束：会话已回收，不发退出码。"+
					"把输出落盘再读，或改用 task.spawn 跑后台任务。", sessionPendingCap)}
		}
		select {
		case <-s.pending.notify:
			// New bytes arrived: take another look. The token is coalescing,
			// so a burst of output costs one wake-up, not one per chunk.
		case <-s.streamOver:
			// The stream just ended. Bytes can land in the buffer between the
			// reader's final append and its close, so drain once more and
			// search that tail before treating the exit as the answer.
			tail := s.pending.takeAll()
			have.WriteString(tail)
			pending := have.String()
			if code, idx, ok := findMarkerCode(pending, marker); ok {
				out := strings.TrimRight(pending[:idx], "\r\n")
				return sessionRun{output: out, code: code}
			}
			// The shell process itself exited (for example the command was
			// `exit N`). Reap it; the exit code comes from the Wait state
			// field - the a1 ruling's primary channel, used exactly where a
			// process boundary exists.
			st, werr := carrier.Wait()
			s.dropLocked()
			code := -1
			if st != nil {
				code = st.ExitCode()
			}
			return sessionRun{
				output:   strings.TrimRight(pending, "\r\n"),
				code:     code,
				fromWait: true,
				ended:    true,
				errText:  waitNote(werr, code),
			}
		case <-ctx.Done():
			// Budget gone, no boundary in the stream: this is the desync the
			// marker exists to detect (AC#2). Kill the tree, reap, report NO
			// exit code and NO attributable output; the next call starts a
			// fresh session.
			s.dropLocked()
			return sessionRun{timeout: true, reset: true}
		}
	}
}

func waitNote(werr error, code int) string {
	if werr == nil || code >= 0 {
		return ""
	}
	return fmt.Sprintf("；shell 进程的 Wait 没给出可认的退出码（%v）", werr)
}

// errSessionCwdOnLive is the refusal for a cwd parameter aimed at a session
// that already exists.
var errSessionCwdOnLive = fmt.Errorf("cwd on a live session")

// findMarkerCode scans for the LAST occurrence of the marker that carries a
// payload of the shape "<marker> <digits>". Earlier occurrences (cmd echoing
// our echo line if a command re-enabled echo) hold the literal
// %errorlevel%, not digits, and are skipped by the digit check - the
// unguessable per-call marker is what makes a collision impossible.
func findMarkerCode(pending, marker string) (code int, idx int, ok bool) {
	search := len(pending)
	for search > 0 {
		i := strings.LastIndex(pending[:search], marker)
		if i < 0 {
			return 0, 0, false
		}
		rest := strings.TrimLeft(pending[i+len(marker):], " \t")
		j := 0
		for j < len(rest) && rest[j] >= '0' && rest[j] <= '9' {
			j++
		}
		if j > 0 {
			return atoiDigits(rest[:j]), i, true
		}
		search = i
	}
	return 0, 0, false
}

func atoiDigits(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		n = n*10 + int(s[i]-'0')
	}
	return n
}

// ---------------------------------------------------------------------------
// shell.session (C1 Tool)
// ---------------------------------------------------------------------------

type shellSession struct{ d ShellDeps }

type shellSessionArgs struct {
	Argv []string `json:"argv"`
	Cwd  string   `json:"cwd"`
}

var shellSessionSchema = JSONSchema(`{"type":"object","properties":{` +
	`"argv":{"type":"array","items":{"type":"string"},"minItems":1,` +
	`"description":"命令与参数向量（强制 argv 形，SPEC-06 §10）；每发都单独过风险判定与审批"},` +
	`"cwd":{"type":"string","description":"可选：只在会话建立的那一发生效的初始工作目录（经 C26 规范化并判级）；` +
	`之后的目录切换请直接发 cd /d 命令"}}` +
	`,"required":["argv"],"additionalProperties":false}`)

func (shellSession) Name() string { return "shell.session" }

func (shellSession) Description() string {
	return "在你这个任务自己的常驻命令会话里跑下一条命令：工作目录和环境变量在两次调用之间延续" +
		"（装环境、跑测试这类连续几步的活用它）。每发命令仍单独过风险判定与审批；" +
		"一发超时整条会话会被回收重来；任务结束（含被取消）时会话跟着收掉"
}

func (shellSession) Parameters() JSONSchema { return shellSessionSchema }

// Execute implements Tool. The task id is the session's identity and it
// comes from the host's context, never from a parameter - same shape as
// task.cancel's caller identity. A call with no host task id has no scope
// to bind a session to, so it refuses (fail-closed).
func (t shellSession) Execute(ctx context.Context, params json.RawMessage, onUpdate func(string)) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	var a shellSessionArgs
	if err := json.Unmarshal(params, &a); err != nil {
		return Result{Text: "参数解析失败：" + err.Error(), IsError: true}, nil
	}
	if len(a.Argv) == 0 {
		return Result{Text: "缺少 argv 参数（shell.session 与 shell.exec 一样只收 argv 向量）", IsError: true}, nil
	}
	for i, arg := range a.Argv {
		if arg == "" {
			return Result{Text: fmt.Sprintf("argv[%d] 是空串（fail-closed 拒绝）", i), IsError: true}, nil
		}
	}
	taskID := TaskID(ctx)
	if taskID == "" {
		return Result{Text: "这条调用没有宿主给的任务 id（fail-closed：会话寿命必须绑任务作用域，" +
			"没有作用域就不开会话）", IsError: true}, nil
	}
	if t.d.Sessions == nil {
		return Result{Text: "会话名册未接线（fail-closed：没有任务结束边界能收的会话，拒绝开）", IsError: true}, nil
	}
	cwd := ""
	if strings.TrimSpace(a.Cwd) != "" {
		if t.d.Paths == nil {
			return Result{Text: "cwd 传了，但 C26 规范化器未接线（fail-closed：这条目录没法判，拒绝执行）", IsError: true}, nil
		}
		c, err := t.d.Paths.Canonicalize(a.Cwd)
		if err != nil {
			return Result{Text: "cwd 无法规范化（C26 拒绝）：" + cwdFailText(err), IsError: true}, nil
		}
		cwd = c
	}

	sess := t.d.Sessions.sessionFor(taskID)
	sess.mu.Lock()
	if ctx.Err() != nil {
		// The call's own budget/cancel arrived while queued behind another
		// command of the same session: leave the session untouched - it is
		// healthy, THIS call just never ran.
		sess.mu.Unlock()
		return Result{}, ctx.Err()
	}
	run := sess.run(ctx, t.d, a.Argv, cwd)
	sess.mu.Unlock()

	return sessionResult(run), nil
}

// sessionResult renders one sessionRun as the C1 Result the model reads.
func sessionResult(run sessionRun) Result {
	switch {
	case run.errText != "":
		return Result{Text: run.errText, IsError: true}
	case run.timeout:
		return Result{
			Text: "这条命令没在 C22 预算内回来，而且会话流里没有出现本发的结束标记" +
				"——这是失去同步，不是「命令失败」：整条会话进程树已回收，" +
				"本发不报退出码（拿不到真值），本发已积的输出也不归属给任何结果。" +
				"下一发 shell.session 会开一条全新会话（工作目录与环境从头开始）。",
			IsError: true,
		}
	}

	text := run.output + "\n"
	if run.ended {
		if run.code >= 0 {
			text += fmt.Sprintf("会话进程已退出，退出码：%d（来自 cmd 进程的 Wait 状态字段）。", run.code)
		} else {
			text += "会话进程已退出。" + run.errText
		}
		text += "会话已回收；下一发 shell.session 会开一条全新会话（工作目录与环境从头开始）。"
		return Result{Text: text, IsError: run.code != 0}
	}
	text += fmt.Sprintf("退出码：%d（会话内埋标记回读，每发单独判级）", run.code)
	return Result{Text: text, IsError: run.code != 0}
}

// ---------------------------------------------------------------------------
// the production carrier: cmd.exe through pipes, inside its own job
// ---------------------------------------------------------------------------

// cmdSessionShell is the production SessionShell: one `cmd.exe /q /d`
// process with a stdin pipe, a stdout pipe the session's reader drains, a
// guarded stderr buffer (cmd.exe's own diagnostics, drained by os/exec's
// own goroutine), and its own kill-on-close tree (ShellDeps.NewTree;
// proc.OpenJobScope's shape on Windows).
type cmdSessionShell struct {
	tree   SessionTree
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.Reader
	stderr guardedBytes

	waitOnce   sync.Once
	waitState  *os.ProcessState
	waitErr    error
	waitClosed chan struct{}
}

// newCmdSessionShell builds the carrier. cwd is the canonical spelling C26
// produced ("" = inherit the host process's directory).
func newCmdSessionShell(d ShellDeps, cwd string) (SessionShell, error) {
	if d.NewTree == nil {
		return nil, ErrShellTreeUnwired
	}
	tree, err := d.NewTree()
	if err != nil {
		return nil, fmt.Errorf("shell: 进程树建立失败: %w", err)
	}
	cmd := exec.Command("cmd.exe", "/q", "/d")
	cmd.Dir = cwd
	stdin, err := cmd.StdinPipe()
	if err != nil {
		_ = tree.Close()
		return nil, fmt.Errorf("shell: stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		_ = tree.Close()
		return nil, fmt.Errorf("shell: stdout pipe: %w", err)
	}
	s := &cmdSessionShell{tree: tree, cmd: cmd, stdin: stdin, stdout: stdout, waitClosed: make(chan struct{})}
	// cmd.exe's OWN diagnostics (unrecognized commands, syntax errors) ride
	// its stderr; os/exec drains it on its own goroutine into the guarded
	// buffer the session takes after each command. Command OUTPUT never
	// comes from here (commands append 2>&1 themselves at the host level).
	cmd.Stderr = &s.stderr
	if _, err := tree.StartInJob(cmd); err != nil {
		_ = tree.Close()
		return nil, fmt.Errorf("shell: 会话进程未能入树启动: %w", err)
	}
	return s, nil
}

func (s *cmdSessionShell) Stdin() io.Writer { return s.stdin }

// Stdout returns the shell's output stream (the pipe's read end).
func (s *cmdSessionShell) Stdout() io.Reader { return s.stdout }

func (s *cmdSessionShell) TakeStderr() string { return s.stderr.takeAll() }

func (s *cmdSessionShell) Kill() error {
	if s.cmd.Process == nil {
		return nil
	}
	return s.cmd.Process.Kill()
}

// Wait reaps the shell process exactly once; concurrent and later callers
// join the same result (the a1 ruling's Wait-state channel).
func (s *cmdSessionShell) Wait() (*os.ProcessState, error) {
	s.waitOnce.Do(func() {
		s.waitErr = s.cmd.Wait()
		s.waitState = s.cmd.ProcessState
		close(s.waitClosed)
	})
	<-s.waitClosed
	return s.waitState, s.waitErr
}

// Close releases the kill-on-close tree: the shell and every descendant die
// with the job handle.
func (s *cmdSessionShell) Close() error {
	return s.tree.Close()
}

// guardedText is the session's raw output buffer, written by the reader
// goroutine and consumed (taken whole) by whichever command is waiting.
// Every append coalesces one wake-up token into notify (capacity 1): the
// waiting command's select sleeps on that token, on the stream end and on
// its own deadline - never on a poll - so a marker that arrives mid-wait
// wakes the search instead of being discovered only at the deadline (the
// misjudgment AC#2 forbids).
type guardedText struct {
	mu     sync.Mutex
	buf    bytes.Buffer
	notify chan struct{}
}

func newGuardedText() guardedText {
	return guardedText{notify: make(chan struct{}, 1)}
}

func (g *guardedText) append(p []byte) {
	g.mu.Lock()
	g.buf.Write(p)
	g.mu.Unlock()
	select {
	case g.notify <- struct{}{}:
	default:
	}
}

// takeAll returns and clears the accumulated raw stream, atomically with
// respect to appends.
func (g *guardedText) takeAll() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := g.buf.String()
	g.buf.Reset()
	return out
}
