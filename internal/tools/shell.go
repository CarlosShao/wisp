package tools

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/risk"
)

// shell family — the D34 command-execution rows ticket 163 lands:
//
//	shell.exec     L2  one-shot: start one process, Wait it, take the exit
//	                    code from the Wait state field (never from output text)
//	shell.session  L2  one session process per task; each call feeds the next
//	                    command into the SAME live shell, so working directory
//	                    and environment carry over (ticket 163 AC#2/AC#3)
//
// Both names entered the frozen D34 table by owner approval (ledger A322,
// commit ced72f8: PLAN.md D34 rows and the SPEC-07 §3 mirror), and both
// registered nowhere until this family: the roster-diff adjudication table
// (docs/evidence/s1/163-ac1b-recount.md §2) ruled "在册＋零实现" before this
// leg started, and this file is the landing of that ruling. Like task.* in
// ticket 164, the two names register together from ONE family constructor
// (BuiltinShellEntries below) - there is no second registration path.
//
// Shapes the frozen text pins and this family keeps:
//
//   - 默认禁用 (D14 conclusion kept on the shell.exec row): the family
//     registers ONLY when ShellDeps.Enabled mirrors [risk] shell_enabled
//     = true - the same shape FSDeps.DeleteEnabled gives fs.delete. The
//     three unwired [risk] shell_* config keys still fail load (the
//     internal/config/unwired.go guard rows are not this ticket's to move;
//     owner wording approval is registered as the one open half in the
//     ticket progress log).
//   - argv 向量强制 (SPEC-06 §10 / D33): string mode is NOT implemented.
//     Every call carries argv only; Facts.ShellString stays false, so R6's
//     string-mode verdict never runs and the [risk] allow_shell_string row
//     keeps saying exactly what is still true.
//   - 每发命令仍按 R6 单独判 (shell.session row, PLAN.md D34): a session is
//     NOT an authorization. Session existence carries zero authority - the
//     bridge re-runs the full C3 + C19 + L2-routing path on EVERY call, so
//     "approved once, allowed forever" is structurally impossible here
//     (ticket 163 AC#4; the teeth test pins two gate visits for two calls).
//   - 会话寿命绑任务作用域 (shell.session row): sessions are keyed by the
//     host's task id and the bridge's task-end boundary collects them
//     (Bridge.CloseTask -> Options.TaskCleanup, one janitor verb, no second
//     closing path). AC#3's teeth test kills a session through that boundary
//     and proves a foreign task's end cannot reach it.
//   - 超时＝真 deadline (ban on wall-clock differences): every wait here is
//     a select on the call's context (the bridge's C22 deadline), never a
//     time.Now() subtraction. AC#2's marker question is answered at the
//     marker code below.
//   - 不引 PTY (shell.session row; the dsh ConPTY half is the "don't copy"
//     list): the session carrier is one cmd.exe driven through stdin/stdout
//     pipes, and every process tree lives in its own kill-on-close Job
//     (proc.OpenJobScope's exact shape through the SessionTree seam), which
//     is what makes "cancel => the whole tree is gone" an OS fact instead of
//     a best-effort kill.

// ShellDeps is what the shell family needs from the host.
type ShellDeps struct {
	// Enabled mirrors [risk] shell_enabled (D14 默认禁用): BuiltinShellEntries
	// registers NOTHING while it is false - the family does not exist for the
	// model until the operator turned it on, exactly like fs.delete.
	Enabled bool
	// Allowlist mirrors [risk] shell_allowlist and is fed to R6 through the
	// Facts hooks below: clean argv whose program hits the allowlist judges
	// L1, everything else L2 (internal/risk/rules_shell.go). Empty = every
	// call judges L2, which is the fail-closed direction.
	Allowlist []string
	// Paths is the C26 canonicalizer, and it is the ONLY way an optional cwd
	// parameter reaches cmd.Dir - a filepath decision made anywhere else is
	// the D22 ban AGENTS §1.2 states verbatim. nil = a call that passes cwd
	// refuses fail-closed instead of trusting the raw string.
	Paths *PathCanonicalizer
	// NewTree hands out one kill-on-close process tree per one-shot call and
	// per session (Windows production: a closure over proc.OpenJobScope,
	// whose StartInJob/Close method set is exactly SessionTree). Closing the
	// tree is what terminates the whole descendant tree at cancel/end time -
	// AC#3's "会话必须跟着收" is this OS fact, not a polite Kill on one pid.
	// nil = the family refuses to start any process (fail-closed): a command
	// run outside any Job containment is the exact posture ticket 100/R21
	// recorded as not landing.
	NewTree func() (SessionTree, error)
	// NewSessionShell overrides the session carrier itself (tests inject a
	// fake shell with pipes they drive). nil = the production carrier:
	// cmd.exe /q /d with a stdin pipe, stdout read end, a guarded stderr
	// buffer and its own tree from NewTree.
	NewSessionShell func() (SessionShell, error)
	// Sessions is the task-scoped session roster, shared with the bridge's
	// task-end janitor (Options.TaskCleanup = the roster's TerminateForTask),
	// so the host and the tool always look at the SAME table - the TaskDeps.
	// Roster shape. nil = shell.session refuses fail-closed ("会话名册未接线"),
	// never silently keeps a session that nothing can collect.
	Sessions *ShellSessions
	// Logf is the family's own audit sink (session create/reset events). The
	// bridge logs every call on its own; nil = quiet.
	Logf func(format string, args ...any)
}

// SessionTree is one kill-on-close process tree (C30's per-scope shape:
// proc.JobScope on Windows satisfies it by method set). StartInJob must
// start cmd and land it in the tree before it does work, killing it when the
// assignment fails - the proc.JobScope.StartInJob contract, restated here so
// this package stays consume-only for internal/proc.
type SessionTree interface {
	StartInJob(cmd *exec.Cmd) (*os.Process, error)
	Close() error
}

// SessionShell is the live shell process a session drives. Production
// carries cmd.exe through pipes (no PTY - the frozen row forbids the
// dependency); tests drive a fake with their own pipes.
type SessionShell interface {
	// Stdin is the shell's command input. Lines end with CRLF; the session
	// writes one command line plus the marker echo line per call.
	Stdin() io.Writer
	// Stdout is the shell's output stream: every command's bytes and every
	// marker line arrive here, and the session's one reader drains it.
	Stdout() io.Reader
	// TakeStderr returns and clears what the shell itself wrote to its
	// diagnostics stream since the last take (guarded inside the carrier -
	// the drain runs on another goroutine). Command OUTPUT never comes from
	// here; this is cmd.exe's own error text only.
	TakeStderr() string
	// Kill terminates the shell process (belt; the tree Close is the brace).
	Kill() error
	// Wait reaps the shell process once it is gone. The returned state is
	// the ONLY source of the shell's own exit code (a1 ruling: the Wait
	// state field, never a guess from output text).
	Wait() (*os.ProcessState, error)
	// Close releases the kill-on-close tree: every process still inside it,
	// the shell and all of its descendants, dies with the handle.
	Close() error
}

// shellTimeout is the C22 budget the family declares for one call (both
// entries). Real shell work (install, build, test) does not fit the 30s
// bridge default; the house number for that class of work is the 30m the
// owner set for cmd/wisp test runs, and a foreground shell command is the
// same class. Longer or detached work belongs to task.spawn. The budget is
// honored as a context deadline - never a wall-clock comparison - so the
// d22scan ban #4 shape has no occurrence here.
const shellTimeout = 30 * time.Minute

// shellExecDecl is the host-side declaration for shell.exec (D34 row:
// L2, 白名单命中 → L1 comes from R6 at C19 time, not from the declaration;
// the declared floor is the row's own L2).
//
// Capabilities/Needs are the frozen C3 token `shell` (PLAN.md:1359): the
// family's whole capability face is one token that already existed. cwd is
// the one model-supplied path, so it is a declared PathParams entry and R2/R3
// judge it through the bridge exactly like an fs path - an out-of-allowlist
// cwd judges L2 fail-closed, never "trusted because it is a directory".
func shellExecDecl(d ShellDeps) Decl {
	return Decl{
		Capabilities: []Capability{CapShell},
		Needs:        []Capability{CapShell},
		Declared:     risk.L2,
		PathParams:   []string{"cwd"},
		Facts:        shellFacts(d),
		Timeout:      shellTimeout,
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// shellSessionDecl declares shell.session (D34 row, PLAN.md:2575 / SPEC-07:75):
// same capability token, same declared floor, same cwd path parameter, and
// the SAME per-call C19 judging - which is the whole AC#4 shape. The session
// row's lifetime wording lands in ShellSessions + Options.TaskCleanup; the
// no-PTY wording lands in the carrier below.
func shellSessionDecl(d ShellDeps) Decl {
	return Decl{
		Capabilities: []Capability{CapShell},
		Needs:        []Capability{CapShell},
		Declared:     risk.L2,
		PathParams:   []string{"cwd"},
		Facts:        shellFacts(d),
		Timeout:      shellTimeout,
		Resident:     true,
		Provider:     KindBuiltin,
	}
}

// shellFacts is the C19 input hook both entries declare. It fills ONLY what
// the tool itself knows: the argv vector (R6's forced shape) and the host's
// allowlist. Facts.Paths belongs to the bridge (it stamps the judged paths),
// and Facts.ShellString stays false - string mode is not implemented here,
// so R6's string-mode verdict never runs and the unwired row for
// allow_shell_string keeps telling the truth.
func shellFacts(d ShellDeps) func(context.Context, map[string]any, []string) risk.Facts {
	return func(_ context.Context, params map[string]any, _ []string) risk.Facts {
		return risk.Facts{
			ShellArgv:      argvFromParams(params),
			ShellAllowlist: d.Allowlist,
		}
	}
}

// argvFromParams extracts the argv vector from the decoded call params.
// A malformed vector (not an array of strings) yields nil, which makes R6
// stand down and leaves the declared L2 floor as the verdict - and Execute
// refuses the same malformed input anyway, so a judgment built on an empty
// vector never reaches a process.
func argvFromParams(params map[string]any) []string {
	raw, ok := params["argv"].([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		s, ok := v.(string)
		if !ok {
			return nil
		}
		out = append(out, s)
	}
	return out
}

// BuiltinShellEntries returns the shell family. Registered ONLY when
// ShellDeps.Enabled mirrors [risk] shell_enabled = true (D14 默认禁用,
// FSDeps.DeleteEnabled shape): a disabled family is not a roster of tools
// that refuse - it is no roster at all. Both names land together from this
// one constructor, per the AC#1b ruling ("两枚名字同批由同一枚工具家族注册").
func BuiltinShellEntries(d ShellDeps) []Entry {
	if !d.Enabled {
		return nil
	}
	return []Entry{
		{Tool: shellExec{d: d}, Decl: shellExecDecl(d)},
		{Tool: shellSession{d: d}, Decl: shellSessionDecl(d)},
	}
}

// sessionLineLimit is the practical cmd.exe command-line ceiling (8191
// chars); a joined line above it is refused instead of silently truncated.
const sessionLineLimit = 8000

// joinSessionArgv renders the session form's ONE input line. The session has
// no process boundary per command - the line is re-parsed by cmd.exe - so
// this join is stricter than a one-shot spawn, and it refuses rather than
// improvising shell quoting:
//
//   - empty elements and pure-space elements: unrepresentable, refuse;
//   - an element containing `"` (or CR/LF/NUL): the line integrity the marker
//     protocol rests on, refuse fail-closed (SPEC-06 §10 keeps string mode
//     off; the metacharacters R6 already tiers to L2 stay legal here because
//     the user approves the exact argv C19 judged - the tool does not
//     second-guess an approved verdict);
//   - elements with spaces are wrapped in double quotes; everything else
//     joins with single spaces.
//
// The trailing " 2>&1" and the marker echo line are HOST machinery appended
// by runCommand, not part of the judged argv and never fed through here.
func joinSessionArgv(argv []string) (string, error) {
	parts := make([]string, 0, len(argv))
	for i, arg := range argv {
		switch {
		case strings.TrimSpace(arg) == "":
			return "", fmt.Errorf("argv[%d] 是空串或纯空白（会话形无法表达，fail-closed 拒绝）", i)
		case strings.ContainsAny(arg, "\"\r\n\x00"):
			return "", fmt.Errorf("argv[%d] 含引号或换行/NUL（会话按行喂数，这一味会破坏行完整性，fail-closed 拒绝；单发请用 shell.exec）", i)
		}
		if strings.ContainsFunc(arg, func(r rune) bool { return r == ' ' || r == '\t' }) {
			parts = append(parts, `"`+arg+`"`)
			continue
		}
		parts = append(parts, arg)
	}
	line := strings.Join(parts, " ")
	if len(line) > sessionLineLimit {
		return "", fmt.Errorf("拼出的命令行 %d 字符超过会话上限 %d（拆成多发，或改用 task.spawn）",
			len(line), sessionLineLimit)
	}
	return line, nil
}

// newSessionMarker mints one random marker: 12 crypto-random bytes, hex
// spelled. It is per-call, unguessable, and free of shell metacharacters by
// construction - the property the exit-code protocol below rests on. A
// session call without randomness cannot bound its own result, so the error
// fails the call instead of degrading to a guessable marker.
func newSessionMarker() (string, error) {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("shell: crypto/rand: %w", err)
	}
	return "w163" + hex.EncodeToString(b), nil
}

// ErrShellTreeUnwired is the fail-closed refusal for a family asked to start
// a process with no kill-on-close tree wired.
var ErrShellTreeUnwired = errors.New("shell: 进程树管控未接线（Job）")
