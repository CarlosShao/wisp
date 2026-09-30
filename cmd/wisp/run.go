package main

// wisp run - the S1 text vertical slice (ticket 12 AC#1) and, more importantly,
// the ASSEMBLY ROOT rulings A8/A11/A13/R12 say was missing.
//
// Before this file the three capabilities that landed this week were each green
// inside their own package and unreachable from production:
//
//	A8  config api_key_ref -> DPAPI store -> the Authorization header a
//	    PROVIDER sets (ticket 63 proved it with a test-written http client)
//	A11 llm.RunProbeSuite had no caller at all (ticket 11 AC#6)
//	A13 approval.Gate and tools.Bridge had zero importers outside their own
//	    packages, so no L1 window and no native card could ever open
//
// What is composed here, in one place, is therefore:
//
//	config.LoadFile -> llm.NewResolver(cfg, secretStore) -> provider      (A8)
//	                -> llm.RunProbeSuite on the save/discovery path       (A11)
//	                -> approval.Gate(console UI)                          (A13)
//	                -> tools.Bridge(gate, C26 paths, journal, provenance)
//	                -> agent.Loop(provider, bridge, AdmitTask=D47 hook)
//	                -> streamed text + notify + task_log/tool_call rows
//
// Ticket 101 added one line to that chain, and it is the whole point of this
// file's current shape: the permission mode is READ from config.toml here and
// INJECTED into the decision chain (config.NewManager -> perm.Store ->
// tools.Options.Modes). Ticket 90 built the storage and proved it works in its
// own tests; an unwired capability is a capability nobody has (R20/M3 did not
// take effect until this line existed).
//
// The dependency direction matters: cmd/wisp is the only package allowed to
// know about all of these at once (internal/llm must not import storage,
// internal/agent/approval imports internal/tools, never the reverse).

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/perm"
	"github.com/CarlosShao/wisp/internal/projctx"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/secret"
	"github.com/CarlosShao/wisp/internal/session"
	"github.com/CarlosShao/wisp/internal/tools"
	// The C5 adapters: llm.NewProvider resolves a config protocol enum through
	// this registry, so a build that forgets one of these lines fails at the
	// first request for that protocol rather than at startup.
	_ "github.com/CarlosShao/wisp/internal/llm/anthropic"
	_ "github.com/CarlosShao/wisp/internal/llm/openaichat"
	_ "github.com/CarlosShao/wisp/internal/llm/openairesponses"
)

// runExit* are the process exit codes of `wisp run`, which the AC pins: "exit
// code reflects error_class" (ticket 12 Key constraints). They are stable
// numbers, not per-run accidents:
//
//	0 completed / control        2 Unconfigured (config or auth, incl. a
//	                              missing credential blob)
//	1 any other classified       3 a loop brake (budget / rounds / repeat)
//	  failure                    4 the user vetoed or rejected the call
func runExitCode(status agent.Status, errClass string) int {
	switch {
	case status == agent.StatusCompleted || status == agent.StatusControl:
		return 0
	case status == agent.StatusStuck || status == agent.StatusTruncated:
		return 3
	}
	switch observe.ErrorClass(errClass) {
	case observe.ClassConfig, observe.ClassAuth:
		return 2
	case observe.ClassUserRejected:
		return 4
	default:
		return 1
	}
}

// notifyPoster posts one user-visible notification (D10 result presentation).
// It is a seam so the end-to-end tests can observe the post; production fills
// it with the real Windows implementation in notify_windows.go.
type notifyPoster func(title, body string) error

// runSpec is one `wisp run` invocation with its process surface injected. The
// zero-value fields are filled by runDefaults, so main() passes only argv and
// the two streams: the data dir, the config file, the credential store, the
// SQLite store and the notification poster are all the REAL ones.
type runSpec struct {
	argv   []string
	stdout io.Writer
	stderr io.Writer
	// dataDir is the env's data root (config.toml, secrets\, memory.db).
	dataDir string
	// notify posts the result notification; nil = postSystemNotification.
	notify notifyPoster
	// probeSink, when non-nil, also receives the capability probe's verdicts.
	// Production passes the provider_health DAO adapter.
	probeSink llm.HealthSink
	// now is the clock source for persisted records (never for deadlines).
	now func() time.Time
	// onRuntime hands the assembled stack to the caller before the task runs.
	// Production leaves it nil; the composition tests use it to dispatch a
	// host-side tool call through the same bridge the loop uses, instead of
	// rebuilding the wiring and proving nothing.
	onRuntime func(*agentRuntime)
	// modeConfirm is the L2 strong confirmation a switch INTO auto_approve
	// costs (R20/M4). Production leaves it nil, which means confirmModeSwitch:
	// one card through the same C18 gate a tool call goes through, answerable
	// only from the native side. The end-to-end tests fill it to stand for "the
	// operator clicked allow", because a console run has no native channel and
	// its L2 card therefore always resolves to a reject.
	modeConfirm perm.ConfirmFunc
	// sink is the persistent log pipeline this run installed. Every caller
	// leaves it nil: runTextTask fills it in after installing the sink, and
	// agentRuntime.auditf books the audit trail through it so the ledger lands
	// on disk without the same sentence reaching the console twice.
	sink *logSink
	// reply is the operator's ANSWER stream for the confirmations this run
	// displays (ticket 201). nil means nobody can answer, which is exactly what
	// every run before this field was: an L2 card waits out the C18 deadline and
	// auto-rejects, an L1 window cannot be opposed. cmdRun fills it with the
	// console's input handle when one is really interactive; the CLI tests fill it
	// with a scripted reader, which is the injection seam AGENTS.md §1.3 names
	// for `wisp run` - it is not a mock standing in for a missing subsystem,
	// because the subsystem here IS a human typing at this terminal.
	reply io.Reader
	// replyVeto is the veto channel this host's cancel transport really is
	// (ticket 201). Production leaves it empty because a console run wires none
	// of SPEC-06 §2's four channels, and an empty value makes Gate.Veto say that
	// back instead of letting this process claim a cancel path it does not have.
	// A host that DOES own one - the ball click, the global Esc hook - names it
	// here and marks it loaded in the same step.
	replyVeto approval.Channel
	// gate, ui and cards are ticket 246 AC#7's injection points, and they exist
	// for exactly one reason: a process may hold ONE approval gate. The resident
	// leg builds its gate at boot (resident_approval_windows.go, form 乙 per
	// orchestrator ruling A481) because the ball and the Esc channel have to be
	// bound before anything can answer a card, and that gate is what the task
	// pipeline this file assembles must use.
	//
	// Leaving all three nil is what every caller before this field was, and what
	// every caller still leaves them for `wisp run`: assembleRuntime composes the
	// console gate, the console UI and a fresh Replies ledger exactly as before.
	// Handing a gate in WITHOUT the ledger and the UI it was built with is
	// refused below rather than patched over, because a second ledger next to a
	// first gate is the second "who is waiting on whom" truth source that 246-a1
	// §5 reason 1 and ledger A481 exist to keep out.
	gate *approval.Gate
	ui   approval.UI
	// cards is the Replies ledger the injected UI books its displayed cards
	// into, and therefore the ledger the answer side spends grants from.
	cards *approval.Replies
	// taskCtx is the parent of the context this run's task loop runs on. nil
	// means context.Background(), which is what `wisp run` has always meant: the
	// process is the task's lifetime. A host that outlives one command line names
	// its OWN task root here (ticket 246 AC#7's ruling 2.3: the resident leg's
	// loop must hang under the ctx that D38(e) step 3 cancels, or "leaving
	// cancels the running task" stays a sentence with no mechanism behind it).
	taskCtx context.Context
}

// runTextTask executes one CLI text task and returns the process exit code.
//
// Every failure here is a CLASSIFIED failure with a user-visible line: SPEC-05
// §3.4 forbids the silent downgrade, so nothing returns 0 unless the task
// really completed.
func runTextTask(s runSpec) int {
	if s.stdout == nil {
		s.stdout = os.Stdout
	}
	if s.stderr == nil {
		s.stderr = os.Stderr
	}
	if s.notify == nil {
		s.notify = postSystemNotification
	}
	if s.now == nil {
		s.now = observe.NowWallUTC
	}
	task := strings.TrimSpace(firstArg(s.argv))
	if task == "" {
		fmt.Fprintln(s.stderr, `wisp run: 没有任务文本（用法：wisp run "总结一下…"）`)
		return 2
	}
	if s.dataDir == "" {
		dir, dirErr := resolveDataDir(buildEnvString())
		if dirErr != nil {
			// Ticket 128 AC#2: no data root, no run. This is the leg that used
			// to answer "." and write its logs, config.toml and DPAPI store
			// into the start-up directory; the resident and `wisp secret` legs
			// already refused the same shape, so refusing here is what makes
			// the three agree. 2 is SPEC-05 §3.4's setup-problem class, which
			// is also what the other two legs return. Booked cost (A105 ⑦): on
			// such a machine this leg is now unusable, which is why the line
			// below carries the fix and not just the verdict.
			fmt.Fprintf(s.stderr, "wisp run: %v\n", dirErr)
			return 2
		}
		s.dataDir = dir
	}

	// Ticket 117: the persistent listener goes on BEFORE the assembly, not
	// after it. assembleRuntime's first statement is secret.NewStore ->
	// winsec.PrivateDirAll, which is already a sealing site, so a sink
	// installed after that point would miss exactly the record this exists to
	// keep. Registered before rt.close() so the shutdown defer below runs
	// first and its own log lines land in the file (LIFO).
	sink, sinkErr := installLogSink(s.dataDir)
	if sinkErr != nil {
		// Loud and named, and it does not stop the run (SPEC-05 §3.4 forbids
		// the silent downgrade, not the degraded boot): refusing to start over
		// a log directory would turn a full disk into an app that will not
		// open, while the notices keep going to stderr as before.
		fmt.Fprintf(s.stderr, "wisp run: 持久日志未启用（%v）：本轮安全告警只会到 stderr，不会落盘\n", sinkErr)
	} else {
		s.sink = sink
		defer sink.close()
	}

	rt, code := assembleRuntime(s)
	if rt != nil {
		// The store opens before the registry is filled, so a partially
		// assembled runtime still has to be closed: a failed run that leaks
		// wisp.db holds the per-session file open after the process is gone.
		defer rt.close()
	}
	if code != 0 {
		return code
	}
	if s.onRuntime != nil {
		s.onRuntime(rt)
	}
	return rt.execute(s.taskCtx, task)
}

// runtime is the assembled S1 stack.
type agentRuntime struct {
	spec  runSpec
	cfg   *config.Config
	mgr   *config.Manager
	paths *tools.PathCanonicalizer
	store *memory.Store
	// session is ticket 224's D45 grants ledger for THIS process: it holds the
	// minted session identity, writes the rows a native 「本会话内允许」 answer
	// produces, and is the read side the bridge consults. nil means this boot
	// could not mint (or could not open the ledger), which leaves the session
	// scope unusable and every call asking - the fail-closed direction.
	session *session.Ledger
	// tasks is ticket 164's process-local background-task table: the only
	// thing task.output reads. v1 keeps no roster across restarts (票 164
	// 定案②), so a restart answers "查不到这个任务" rather than empty.
	tasks *tools.TaskRoster
	// loopOpt / loopOptSet are ticket 197's derivation seam: the agent.Options
	// the CURRENT root loop was assembled with, kept so task.spawn builds its
	// child from the SAME host wiring (same C5 provider, same admission hook,
	// same budgets) instead of inventing a second runtime next to it. Set in
	// execute(), read through SubagentDeps.BaseOptions per spawn.
	loopOpt    agent.Options
	loopOptSet bool
	// c25 is the provenance engine the bridge was built with, kept so the
	// spawner can stamp a subagent's conclusion into the PARENT task's scope
	// under the existing risk.SrcTaskOutput source name (197 §0: no invented
	// source names, or a child's outside text travels a channel nobody owns).
	c25 *risk.Provenance
	// provs is the built provider chain, kept for the probe path.
	provs    []llm.LlmProvider // built chain, kept for the probe path
	names    []string
	endpoint llm.Endpoint
	gate     *approval.Gate
	ui       *consoleApprovalUI
	bridge   *tools.Bridge
	// liveCards is the native-side ledger of the cards this run displayed: the
	// one place a displayed card's single-use grant is kept, and therefore the
	// one thing that makes an allow possible at all (ticket 201). It is filled
	// by the injected UI, which is the only recipient of the grant, and read by
	// the reply surface. A card leaves it when it is answered, dismissed, or
	// handed off to execution.
	liveCards *nativeCards
	// reply is this run's answer side, and replyRoot/replyHandle its owner and
	// goroutine. All three stay nil when no reply source was handed in, which is
	// the honest statement that nobody can answer this run's cards.
	reply       *replySurface
	replyRoot   *observe.Root
	replyHandle *observe.Handle
	// reloadRoot / reloadHandle own this run's config-reload tick (ticket 223:
	// the first production caller of Manager.CheckAndReload in `wisp run`). The
	// root is what the D36 confirmation card is parked on - cancelling it at
	// close abandons an unanswered card, which resolves as a deny - and it is
	// also the ctx the hook reads, because ConfirmLocked's own signature carries
	// no context.
	reloadRoot   *observe.Root
	reloadHandle *observe.Handle
	// modes is the assembled owner of the permission mode (ticket 90's storage,
	// ticket 101's wire): the read the bridge consults once per call, and the
	// only object in this process that may change the档.
	modes *perm.Store
	// modeWrites is ticket 114 AC#2's gate: the one handler a panel mode request
	// is answered by, assembled with the SAME L2 leg the Store above was handed.
	// Nothing calls it yet - the WebView2 "event -> ParseComposerRequest" hop does
	// not exist in this tree (tickets 33/35) - so this is the gate's other half,
	// not a path: the host that gets wired later cannot assemble a wider档 without
	// a confirmation leg.
	modeWrites *panel.ModeWriteHandler
	// pump assembles the packet the panel is rendered from, out of the live
	// objects this struct already holds (ticket 35's data half). It is the first
	// thing in this repository that calls panel.NewSnapshot /
	// panel.NewComposerState from a running process rather than from a test -
	// see panel_pump.go for what its exit is and what it deliberately is not.
	pump *panel.SnapshotPump
	// stream is the results section's producer: the deltas this run streams,
	// accumulated so a snapshot can carry them.
	stream *panel.StreamLog
	// seenTasks are the task ids this process admitted - every root loop and
	// every child loop, because agent.Options.AdmitTask is the one hook both of
	// them pass through (internal/agent/loop.go:360, and task.spawn wraps the
	// same hook at internal/tools/subagent_197.go's opt.AdmitTask). It is the
	// enumeration set for the panel's roster reader: tools.TaskRoster has no
	// "list everything" port, and minting one is ticket 194's roster work, not
	// this leg's. taskMu guards it because admitTask runs on each task's own
	// goroutine while the pump reads it from the approval UI's and the sink's.
	taskMu    sync.Mutex
	seenTasks []string
	// instrLoader is ticket 200's project-instruction loader for the CURRENT run,
	// kept so the panel pump has a live reader for the instructions carrier
	// (200-r2, AC#7's second hop). instrMu guards it because execute() writes it
	// while publishes read it from the approval UI's and the sink's goroutines.
	instrMu     sync.Mutex
	instrLoader *projctx.Loader
	// lastSnap / lastSnapBytes are the most recent packet this run built, kept
	// for the same reason consoleApprovalUI.cards is kept: the production path
	// records what it did, and a reader asks afterwards. The ledger cannot hold
	// the packet (internal/observe bounds one logged string at 512 chars), so
	// this is where the bytes are while the transport is still missing.
	snapMu        sync.Mutex
	lastSnap      panel.Snapshot
	lastSnapBytes []byte
	snapSeen      bool
	notify        notifyPoster
	stdout        io.Writer
	stderr        io.Writer
	// logf receives the bridge's audit lines.
	logf func(string, ...any)
}

// assembleRuntime builds the whole stack from the data dir. A failure here is
// the Unconfigured state (SPEC-03 §4.1: never fall back to a half-configured
// runtime silently), which exits 2 with a named reason.
func assembleRuntime(s runSpec) (*agentRuntime, int) {
	// The result poster is defaulted HERE, not only in runTextTask, because
	// execute() calls it on every task's way out and a caller that assembled
	// through this function without going through runTextTask (ticket 246 AC#7's
	// resident leg does exactly that) would otherwise be calling a nil func.
	if s.notify == nil {
		s.notify = postSystemNotification
	}
	rt := &agentRuntime{spec: s, stdout: s.stdout, stderr: s.stderr, notify: s.notify}
	cfgPath := filepath.Join(s.dataDir, configFileName)
	st, err := secret.NewStore(s.dataDir)
	if err != nil {
		fmt.Fprintf(s.stderr, "wisp run: 凭据存储不可用：%v\n", err)
		return nil, 2
	}
	// res=nil on purpose: config keeps the api_key_ref and the ONLY thing that
	// turns a ref into a key is llm.Resolver -> secret.Store (A8). Had we
	// resolved through LoadFile here, the provider's key would come from a
	// map the loader filled and this ticket's credential assertion would be
	// proving the loader instead of the provider.
	//
	// Manager rather than LoadFile (ticket 101): the permission mode is the one
	// preference R20/M3 lets outlive a session, and it persists THROUGH the
	// config file, so the assembly root needs the object that can write it back
	// atomically (Manager.SetPermissionMode). Loading the file twice over - once
	// read-only here, once writable elsewhere - is exactly the split state
	// SPEC-03 §3.1 exists to prevent, so there is one Manager and one truth.
	mgr, err := config.NewManager(cfgPath, nil)
	if err != nil {
		// AC#3 (ticket 101): a mode that cannot be read is a CLASSIFIED failure,
		// never "keep the last value we saw". This process caches no mode of its
		// own - the档's only source is the file that just failed - so the honest
		// answer is: audit which档 a fail-closed read answers with (the
		// strictest), refuse to assemble a decision chain at all, exit 2.
		rt.auditModeUnreadable(cfgPath, err)
		fmt.Fprintf(s.stderr, "wisp run: 配置未就绪（Unconfigured）：%v\n", err)
		return rt, 2
	}
	// A snapshot copy: the runtime's non-mode settings stay frozen at what this
	// boot read, while perm.Store keeps reading the live section (ticket 90's
	// rule that the running mode always equals what the config layer concluded).
	cfg := mgr.Config()
	rt.cfg = cfg
	rt.mgr = mgr

	// Role -> endpoint -> provider, with the key resolved through the DPAPI
	// store by the resolver (A8's whole point).
	res := llm.NewResolver(cfg, st)
	ep, _, err := res.ResolveRole(llm.RoleChat)
	if err != nil {
		fmt.Fprintf(s.stderr, "wisp run: 文本角色未配置（Unconfigured）：%v\n", err)
		return rt, 2
	}
	rt.endpoint = ep
	prov, err := llm.BuildEndpointProvider(ep, chainOptions(cfg))
	if err != nil {
		fmt.Fprintf(s.stderr, "wisp run: 无法构造 provider：%v\n", err)
		return rt, 2
	}
	rt.provs = []llm.LlmProvider{prov}
	rt.names = []string{ep.String()}

	mem, err := memory.Open(s.dataDir)
	if err != nil {
		fmt.Fprintf(s.stderr, "wisp run: 存储不可用：%v\n", err)
		return rt, 2
	}
	rt.store = mem

	// The host session identity (ticket 224, approval record A435).
	//
	// One mint per process, from crypto/rand, and nothing derived: A435 clause 1
	// forbids recomputing it from the pid / the clock / the working directory,
	// because every approval_grant row on disk is keyed by this exact string. A
	// recomputable id would let the next boot read the last boot's D45 grants
	// back, which is the permanent免审通行证 internal/perm/store.go:19-27 and
	// PLAN.md:1642 name as the thing to keep out. Clause 2 - the session ends
	// when this process does - needs no code at all: this value dies here, so no
	// later Covering call can be keyed with it, and AC#3 (重启必失效) holds by
	// construction rather than by a cleanup step somebody might forget.
	//
	// Failing to mint is fatal for THIS feature and for nothing else: the ledger
	// stays nil, so 「本会话内允许」 answers as a single-use allow and says so, and
	// the bridge reads no grant source at all. A run that cannot get a session
	// must still be a run, and it must be a run that asks.
	sessID, err := session.Mint()
	if err != nil {
		fmt.Fprintf(s.stderr, "wisp run: 本次没有会话身份（%v），「本会话内允许」将按「仅本次」答复处理\n", err)
		rt.auditf("wisp run: SESSION-MINT-FAILED err=%v (会话授权档今天不可用)", err)
	} else {
		ledger, lerr := session.NewLedger(session.LedgerOptions{
			ID:    sessID,
			Store: mem,
			Logf:  rt.auditf,
		})
		if lerr != nil {
			fmt.Fprintf(s.stderr, "wisp run: 会话授权记账不可用：%v\n", lerr)
			rt.auditf("wisp run: SESSION-LEDGER-FAILED err=%v", lerr)
		} else {
			rt.session = ledger
			// The one line that makes "which session wrote this row" answerable
			// from the log instead of from the database afterwards.
			rt.auditf("wisp run: SESSION-MINT id=%s (结束点＝本进程退出，A435 第 2 条)", sessID)
		}
	}

	// C26 canonicalizer over the [fs] allowlist; the same instance is shared
	// by the bridge and the fs tools, so the path the risk verdict judged is
	// the path that opens.
	allowed := make([]string, 0, len(cfg.FS.AllowedDirs))
	for _, d := range cfg.FS.AllowedDirs {
		allowed = append(allowed, d)
	}
	rt.paths = tools.NewPathCanonicalizer(allowed, cfg.FS.ReparsePointExceptions)

	reg := tools.NewRegistry()
	// WHY notify IS NOT A TOOL HERE YET (a blocker for the owner, not a
	// shortcut). SPEC-07 §3's frozen S1 roster names the first two builtins
	// `notify` and `list_tools`, and internal/tools' C1 registration validates
	// the 'namespace.action' name shape - which both of those names violate
	// (`Registry.Register` refuses "notify": "name must be 'namespace.action'").
	// Renaming it to `system.notify`/`wisp.notify.send` departs from the frozen
	// table just as badly, and internal/tools is consume-only for this ticket,
	// so neither side gets edited here. The notification therefore stays on the
	// CLI's own D10 result path (which is where the S1 demo needs it anyway)
	// and list_tools keeps its D15(2) in-loop answer. Registering notify as a
	// C4 tool needs a ruling on which frozen name wins.
	for _, e := range tools.BuiltinFSEntries(tools.FSDeps{
		Paths:         rt.paths,
		DeleteEnabled: cfg.FS.DeleteEnabled,
	}) {
		if err := reg.Register(e); err != nil {
			fmt.Fprintf(s.stderr, "wisp run: 工具注册失败（%s）：%v\n", e.Tool.Name(), err)
			return rt, 2
		}
	}
	// The task family (ticket 164 AC#3, plus ticket 221 甲形 per ruling A434) is
	// what reads and stops the process-local task table above:
	//
	//	task.output  reads back what one row's task printed
	//	task.cancel  stops ONE subagent the CALLING task itself derived
	//
	// Both are registered from tools.BuiltinTaskEntries below, so this line is
	// task.cancel's production landing point (ticket 221 AC#2: TaskRoster's Cancel
	// had zero production callers until the tool that reads it got registered -
	// count it with a grep for calls of that method outside _test, which answers
	// exactly one hit now, internal/tools/task.go's taskCancel.Execute; the
	// description of the measure is spelled out here without the call syntax on
	// purpose, so this comment can never be miscounted as a caller).
	// Who FILLS the roster is ticket 197's spawner, which now
	// does (PublishSubagent on both the spawn goroutine and the child's admission
	// hook), and the row a cancelled child lands on is finalised by that spawner's
	// finish watcher, not by this tool. A call for an id nobody filed still answers
	// "查不到这个任务" - loud, per 票 164 定案②, and never an empty success that
	// would read as "the task printed nothing".
	rt.tasks = tools.NewTaskRoster()
	for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: rt.tasks, Paths: rt.paths}) {
		if err := reg.Register(e); err != nil {
			fmt.Fprintf(s.stderr, "wisp run: 工具注册失败（%s）：%v\n", e.Tool.Name(), err)
			return rt, 2
		}
	}

	// The approval layer (A13). Channels loaded here are the honest set for a
	// console run: no floating ball, no global Esc hook, so an L1 window can
	// only expire unvetoed and an L2 card can only time out into a reject.
	//
	// Ticket 201 changed one half of that sentence and left the other half
	// standing: the console can now ANSWER a card (a reply source, attached at
	// the end of this function), so an L2 card no longer has to be a question
	// nobody hears. What it still cannot do is veto an L1 window, because none
	// of SPEC-06 §2's four veto channels (ball / Esc / KWS / panel) exists in a
	// terminal - so NewChannels() stays empty and runSpec.replyVeto stays unset,
	// and approval_reply.go says what both of those choices protect.
	// Ticket 224's two injections of the SAME ledger, split by which side of the
	// boundary they serve: the approval layer may only WRITE a session rule (a
	// native click produced it), the bridge may only READ one (it decides whether
	// a question still has to be asked). Neither gets the other's half, so no
	// component in this process can both grant itself a rule and wave a call
	// through on it. See the comment at the permission-mode block below for why
	// this is not the same thing as Options.Confirmations.
	//
	// The typed-nil guard is load-bearing: a nil *session.Ledger assigned into an
	// interface field yields a NON-nil interface value, and both consumers fail
	// closed by testing that field against nil.
	var grantRead tools.GrantSource
	var grantWrite approval.GrantRecorder
	if rt.session != nil {
		grantRead = rt.session
		grantWrite = rt.session
	}

	// Two postures, and the second one is ticket 246 AC#7 (orchestrator ruling
	// 2.2: 「一个进程只许一枚 approval.Gate」):
	//
	//   - no injected gate: this assembly owns the gate, the console surface and
	//     a fresh ledger, exactly as every run since ticket 12 has done;
	//   - injected gate: the host built all three BEFORE it called here (the
	//     resident leg needs the ball bound to that gate's UI to load its Esc
	//     channel, and the ball exists before a config file is ever read), and
	//     this function owes that host the one thing it never had - a task
	//     pipeline that asks through THAT gate instead of standing a second one
	//     next to it.
	//
	// The cost of the second shape is named, not hidden: a gate built before this
	// function ran cannot have been handed the session ledger minted above, so
	// 「本会话内允许」 has no writer here. That is the state approval.Gate answers
	// with when Options.Grants is nil - the call is released, nothing is stored,
	// and the gate books `approval: GRANT-DROPPED ... 本机没有接入会话授权记账`
	// on the route (gate.go's allowSession). resident_windows.go prints the same
	// limit at boot; widening it needs the gate itself moved into this function,
	// which is a ruling about AC#1's form 乙, not a line this leg may write.
	if s.gate != nil {
		if s.cards == nil || s.ui == nil {
			fmt.Fprintf(s.stderr, "wisp run: 注入的审批门必须连同它自己的界面与会话账本一起递进来"+
				"（gate=%v ui=%v cards=%v），已拒绝装配\n", s.gate, s.ui, s.cards)
			return rt, 2
		}
		rt.gate = s.gate
		rt.ui = nil // no console surface in a process whose cards go to another host
		rt.liveCards = &nativeCards{h: s.cards}
	} else {
		rt.liveCards = newNativeCards()
		rt.ui = &consoleApprovalUI{out: s.stdout, run: rt, live: rt.liveCards}
		rt.gate = approval.New(approval.Options{
			UI:              rt.ui,
			Channels:        approval.NewChannels(),
			Window:          time.Duration(cfg.Risk.L1WindowSec) * time.Second,
			ApprovalTimeout: time.Duration(cfg.Risk.ConfirmTimeoutSec) * time.Second,
			Logf:            rt.auditf,
			Grants:          grantWrite,
		})
		// Bind the reply ledger to the gate it answers to (ticket 201 AC#1). The
		// ledger had to exist before the gate, because the UI that fills it is one of
		// the gate's options; from this line on, any host surface holding this ledger
		// can route an answer. No veto channel and no reply source yet means exactly
		// what it says: cards are recorded, an undeliverable veto is refused by the
		// gate's own registry, and attachReplyListener re-binds with the channel this
		// host really wired when (and if) it has one.
		rt.liveCards.bind(rt.gate, "", "", "")
	}

	// The permission mode (ticket 90's storage, wired here by ticket 101). This
	// is the line that makes R20/M3 real: without it the bridge reads a nil
	// ModeSource, which answers the strictest档 for every call, and a manually
	// chosen档 would be silently forgotten at the next start.
	//
	// Two other injections sit next to it, and they are NOT the same kind of
	// thing, which is worth spelling out because this comment used to conflate
	// two of them:
	//
	//   - tools.Options.Confirmations stays nil, and stays nil after ticket 224.
	//     That field is the B-tier SINGLE-FILE override set risk.Gate reads, and
	//     populating it is ticket 21's confirmation flow. Nothing in this
	//     assembly touches it.
	//   - tools.Options.Grants is new here: the read side of a D45 session grant
	//     (internal/tools/grant.go). R20/M3 still persists the MODE and nothing
	//     else - a session grant does not come back from disk "because this boot
	//     learned to read config.toml". It comes back only ever keyed by the
	//     identity minted above, which no later boot can recompute (A435 第 1
	//     条), so PLAN.md:1642's "会话结束后授权必须失效" is the property the read
	//     depends on rather than a rule this file has to remember to enforce.
	confirm := s.modeConfirm
	if confirm == nil {
		confirm = rt.confirmModeSwitch
	}
	modeStore, err := perm.New(perm.Options{
		Manager: mgr,
		Confirm: confirm,
		Logf:    rt.auditf,
	})
	if err != nil {
		rt.auditModeUnreadable(cfgPath, err)
		fmt.Fprintf(s.stderr, "wisp run: 权限档位不可用（Unconfigured）：%v\n", err)
		return rt, 2
	}
	rt.modes = modeStore

	// Ticket 114 AC#2 (R-92-2): the panel's档 request gets ONE handler, and it is
	// handed the same `confirm` value the Store got above - so "no L2 leg attached"
	// means the same thing on both sides of the boundary and the handler refuses a
	// widening request before ever reaching Set. The handler never calls this leg
	// itself (Set raises the single card; a second ask would be two 300s windows);
	// the forward is here so that the day the ask moves, it moves to the same C18
	// route. Switch.At is left alone: perm.Set stamps its own record and
	// confirmModeSwitch reads only From/To.
	rt.modeWrites = &panel.ModeWriteHandler{
		Modes: modeStore,
		Confirm: func(ctx context.Context, from, to risk.Mode, origin, actor string) error {
			return confirm(ctx, perm.Switch{From: from, To: to, Origin: origin, Actor: actor})
		},
		Audit: rt.auditf,
	}

	// Which posture this boot came up in has to be readable in the log, not
	// inferred from the file afterwards. The Store keeps its startup record in
	// its own history; the line is the assembly root's to write, because "this
	// host started in auto_approve" is the first fact an operator reading a
	// long-running host's audit wants, and it is also what makes a restart
	// auditable against the mode read before it (AC#3's audit half).
	rt.auditf("perm: MODE-READ origin=startup mode=%s source=%q",
		modeStore.PermissionMode(), cfgPath)

	// The panel's snapshot pump (ticket 35's data half). Every reader below is a
	// live object this process is actually running - the approval queue, the perm
	// store, the C26 canonicalizer, the stream this run is writing - so the
	// packet is built from state, not from a test's arguments. What it does NOT
	// have is a page to go to: this tree carries no WebView2 host (tickets 33/35),
	// so the bytes are booked to the run's persistent ledger and the last mile is
	// still open. panel_pump.go states the whole rule.
	rt.stream = panel.NewStreamLog(panel.DefaultStreamKeys)
	rt.pump = panel.NewSnapshotPump(panel.PumpSources{
		Verdicts:  rt.liveVerdicts,
		Mode:      rt.modes.PermissionMode,
		Workspace: rt.workspaceView,
		Git:       rt.gitView,
		Model:     rt.currentModel,
		Results:   rt.stream.Chunks,
		// ticket 200 AC#7's second hop, closed by 200-r2: the carrier existed at
		// r1 and NOTHING read it, which is the same shape as ledger A408's
		// "文案在、控件不在" moved from the widget side to the wire side. This
		// reader hands the pump the loader's own last bundle, so the packet
		// reports what was loaded (or which of the five states stopped it)
		// instead of a key that can never carry a value.
		Instructions: rt.instructionBundle,
		// ticket 197 载体层's reader, added in the shape the line above it
		// established (票 200's hop): the carrier exists on the pump, so the
		// assembly root owes it a reader of a live object, or the wire key can
		// never arrive. This one reads the roster this process is running plus
		// the stream log this process is writing - see panel_pump.go.
		Tasks: rt.taskRosterState,
		Out:   rt.bookPanelSnapshot,
	})
	// The console surface publishes a packet when the approval state moves. A
	// host whose cards go somewhere else (ticket 246 AC#7's injected gate) has no
	// console surface here, and its own UI owns its publishing - which this leg
	// does NOT give it: the panel route is ticket 33/35's, not AC#7's, and
	// inventing a second publish path from the resident side would be a claim
	// about a page this process cannot show. The pump stays assembled and its
	// bytes keep landing in this run's ledger, exactly as panel_pump.go says.
	if rt.ui != nil {
		rt.ui.publish = rt.publishPanelSnapshot
	}

	// c25 is the ONE C25 engine this process has: the bridge marks sensitive
	// results with it, and task.spawn stamps a subagent's conclusion into its
	// parent's scope with the same engine under the same registered source name
	// (risk.SrcTaskOutput). Handing the spawner a second engine would split the
	// taint index the R4 checks read.
	c25 := risk.NewProvenance(risk.ProvOptions{NoProbe: true})
	rt.c25 = c25
	rt.bridge = tools.New(tools.Options{
		Registry: reg,
		Paths:    rt.paths,
		Gate:     rt.gate,
		Cancel:   rt.gate.ToolsCancelBus(),
		Journal:  mem,
		// AC#1/AC#4's anchor line: the档 this boot read enters the decision
		// chain here and nowhere else.
		Modes: rt.modes,
		// Ticket 224 AC#1/AC#2's anchor line: this is the only place a D45
		// session authorization enters the decision chain. Nil when this boot
		// could not mint a session, which is the state every pre-224 assembly was
		// in and which leaves every question asked.
		Grants: grantRead,
		// R4 stays dormant: probing every result for sensitive sources needs
		// the C25 detector's own wiring (ticket 25), and a dormant R4 is the
		// honest state, not a weakened one.
		Provenance:     c25,
		DefaultTimeout: time.Duration(cfg.Agent.PerToolTimeoutMS) * time.Millisecond,
		OnDecision: func(d tools.Decision) {
			rt.auditf("wisp run: 风险判定 tool=%s level=%s rules=%v reason=%q",
				d.Tool, d.LevelString(), d.RulesHit, d.Reason)
		},
		Logf: rt.auditf,
	})

	// task.spawn (ticket 197 leg A) is the only tool this run adds beyond D34's
	// rows, and it is wired to the objects this runtime is actually running: the
	// same roster task.output reads, the same bridge as the child's tool surface
	// (the spawner filters task.spawn out of it, so a child never even sees the
	// name - depth 1 by structure, not by counting), the same C25 engine, the
	// panel's stream log for the subagent:<taskID> keys, and the agent.Options the
	// current root loop was assembled with. Before execute() has assembled one the
	// snapshot is unset and the tool answers a fail-closed refusal instead of
	// starting a runtime of its own.
	for _, e := range tools.BuiltinSubagentEntries(tools.SubagentDeps{
		Roster:      rt.tasks,
		ParentTools: rt.bridge,
		Provenance:  rt.c25,
		// The panel's stream log itself, handed over as the two methods the
		// spawner needs (tools.SubagentStreamSink). It is the log rather than a
		// closure around it because Close is half of the job: a finished
		// subagent whose row never says done streams forever on the page.
		Stream:      rt.stream,
		BaseOptions: func() (agent.Options, bool) { return rt.loopOpt, rt.loopOptSet },
	}) {
		if err := reg.Register(e); err != nil {
			fmt.Fprintf(s.stderr, "wisp run: 工具注册失败（%s）：%v\n", e.Tool.Name(), err)
			return rt, 2
		}
	}

	// The answer side, last - so a reply can never arrive at a half-built
	// runtime (ticket 201). With no reply source this is a no-op and the run
	// keeps the posture it had before the field existed: cards get shown, nobody
	// answers them, and each route resolves on its own clock.
	if s.reply != nil {
		rt.attachReplyListener(s.reply, s.replyVeto)
	}

	// The config-reload tick, after the answer side (ticket 223). Order matters
	// and is not cosmetic: the tick can raise an L2 card the moment it sees a
	// hand-edited loosening, and a card raised before the reply listener exists
	// is a card nobody in this process can answer - which would make D33's
	// re-confirmation resolve as a deny for a reason that has nothing to do with
	// what the operator wrote. config_reload.go owns the three lines this call
	// wires (CheckAndReload's caller, ConfirmLocked, OnRestartPending) and says
	// what still does not move mid-run.
	rt.startConfigReload()

	return rt, 0
}

// modeSwitchToolName is the card's tool name for a mode switch. It is not a
// registered C4 tool: no call route uses it, it exists so the card and the
// audit line name what is being confirmed.
const modeSwitchToolName = "permission.mode"

// confirmModeSwitch is the production R20/M4 confirmation: exactly ONE L2 card,
// raised through the same C18 approval gate a tool call goes through, and it can
// only be answered from the native side (SPEC-06:1588, machine-gated by ban #6).
//
// A console run has no native channel, so this refuses - the card times out into
// a reject, or the caller's context ends first and the gate abandons it. That is
// fail-closed by construction and it is the point: the档 that removes questions
// cannot be reached from a surface that cannot ask them. The floating ball /
// panel hosts (tickets 77/92) hand the same card a real click.
//
// D47 applies to a host-initiated card as much as to a loop-initiated one: the
// gate refuses any request whose task was never admitted, so this registers its
// own synthetic task id for the duration of the confirmation and revokes it on
// the way out.
func (rt *agentRuntime) confirmModeSwitch(ctx context.Context, sw perm.Switch) error {
	if rt.gate == nil {
		return fmt.Errorf("审批 gate 未装配，无法为切到 %s 签发 L2 卡片", sw.To)
	}
	const taskID = "host:mode-switch"
	revoke := rt.gate.AdmitTextTask(taskID)
	defer revoke()
	ans, why := rt.gate.PendingApproval(ctx, tools.Decision{
		TaskID: taskID,
		Tool:   modeSwitchToolName,
		Level:  risk.L2,
		Reason: fmt.Sprintf("切换到 %s 会移除本应询问的确认，需要一次 L2 强确认（R20/M4，当前档 %s）",
			sw.To, sw.From),
		Mode: sw.From,
	})
	if ans != tools.AnswerAllow {
		if why == "" {
			why = string(ans)
		}
		return fmt.Errorf("L2 强确认未通过（%s）：%s", ans, why)
	}
	return nil
}

// auditModeUnreadable is AC#3's loud half. It writes the same "[audit] perm:"
// family the Store itself writes, so every mode event in the log has one shape
// and one grep, and it names the档 a failed read answers with instead of
// leaving that to the reader's imagination.
func (rt *agentRuntime) auditModeUnreadable(path string, err error) {
	rt.auditf("perm: MODE-READ-FAILED path=%q err=%v mode=%s origin=startup result=fail-closed "+
		"detail=%q", path, err, risk.DefaultMode(),
		"档位读不到：本进程不缓存任何上一次的宽松值，决策链不会被装配（退出码 2）")
}

// displayedCounter is the one question a host-side presentation surface answers
// about its own output: how many cards did it actually put in front of somebody.
// It is an interface rather than the concrete resident type because run.go is
// not build-tagged and the ball surface is (resident_approval_windows.go), and
// because the count is the surface's own fact, never this file's inference.
type displayedCounter interface {
	displayedCards() int
}

// windowCount reports how many confirmation cards the composed gate displayed.
func (rt *agentRuntime) windowCount() int {
	if rt.ui != nil {
		return rt.ui.shown()
	}
	// A run assembled with an injected gate (ticket 246 AC#7) shows nothing on a
	// console, so the honest count is the one the host's own surface keeps. With
	// neither, this run displayed nothing and 0 is the reading.
	if d, ok := rt.spec.ui.(displayedCounter); ok {
		return d.displayedCards()
	}
	return 0
}

// close releases the store and stands the reply listener down.
//
// The cancel is not a join: a listener parked in a blocking Read on the
// operator's stream does not observe a context, so the root's Cancel is what
// stops it from taking any further reply (runReplyLoop checks ctx.Err() per line
// and refuses to act on one after that). For a CLI process that is the whole
// requirement - the goroutine dies with the process, and a pipe an end-to-end
// test closed returns EOF on its own.
func (rt *agentRuntime) close() {
	if rt.replyRoot != nil {
		rt.replyRoot.Cancel()
	}
	// Stand the reload tick down (ticket 223): this is also what abandons an L2
	// re-confirmation still waiting on a human, and an abandoned card resolves as
	// a deny, so shutting down can never be the moment a loosening slips in.
	if rt.reloadRoot != nil {
		rt.reloadRoot.Cancel()
	}
	if rt.store != nil {
		_ = rt.store.Close()
	}
}

// logf is the audit sink: the bridge's and the gate's lines.
//
// Both halves are load-bearing and they answer two different complaints:
//
//   - stderr is the operator standing at a terminal, reading this run.
//   - slog is the record that outlives the process. Ticket 105 wired the D31
//     path-rewrite account ("tools: d31 report ...") and the mode reads
//     ("perm: MODE-READ ...") into this function, and R-105-1 was right that
//     calling that "persisted" was false: an fmt.Fprintf to a stream is not a
//     ledger. Now the same line also goes through the sink this run installed,
//     the rolling JSONL file under <data>\logs (see logsink.go) - so the
//     rewrite ledger is on disk on the production path, not only in a test that
//     built its own sink. It rides the sink's own logger rather than the
//     process default precisely so the console keeps ONE copy of each line.
func (rt *agentRuntime) auditf(format string, args ...any) {
	line := fmt.Sprintf(format, args...)
	fmt.Fprint(rt.stderr, "[audit] "+line+"\n")
	rt.spec.sink.logger().Info("audit: " + line)
}

// admitTask is this composition root's per-task boundary: it registers the
// TEXT-loop task with the D47 gate and hands back the revocation the loop
// defers (internal/agent/loop.go:366). The revocation now also closes the
// task's C25 taint scope on the bridge (ticket 151).
//
// WHY THIS HOOK OWNS THE CLOSE. CloseTask's own comment says the composition
// root defers it "on the task's ..." boundary, and this is the only boundary
// this root has that (a) is handed the loop's task id - the loop generates it
// inside Run, so nothing above this point knows it - and (b) already runs on a
// defer, so it covers the brake/cancel/error exits too. The alternative
// readings were rejected: closing after loop.Run returns would skip every
// panic exit, and a task-lifecycle event does not exist in this tree yet.
//
// What this does NOT close: a task id a host-internal caller invents when it
// dispatches on the bridge directly (no admission, no boundary). That shape is
// recorded as an open end, not silently folded into this line.
func (rt *agentRuntime) admitTask(taskID string) func() {
	rt.noteTask(taskID)
	revoke := rt.gate.AdmitTextTask(taskID)
	return func() {
		if revoke != nil {
			revoke()
		}
		rt.bridge.CloseTask(taskID)
	}
}

// noteTask records that a task of this process was admitted, which is what makes
// it nameable by the panel's roster reader (ticket 197 载体层). It adds no state
// of its own: the row behind an id is the roster's, and a task that was admitted
// but never filed a row is simply not in the packet.
func (rt *agentRuntime) noteTask(taskID string) {
	if rt == nil || taskID == "" {
		return
	}
	rt.taskMu.Lock()
	defer rt.taskMu.Unlock()
	for _, seen := range rt.seenTasks {
		if seen == taskID {
			return
		}
	}
	rt.seenTasks = append(rt.seenTasks, taskID)
}

// execute runs one task through the loop and presents the result.
//
// parent is the context the task derives from, and the three-way default is the
// whole of ticket 246 AC#7's ruling 2.3: a caller that names none gets
// context.Background() (what `wisp run` has always meant - the process IS the
// task's lifetime), a run assembled with runSpec.taskCtx gets that host's task
// root, and a host that submits one task at a time passes a root derived from it
// so a single task can be cancelled without touching the entry that produced it.
func (rt *agentRuntime) execute(parent context.Context, task string) int {
	cfg := rt.cfg
	prov := rt.provs[0]
	info := prov.Info()
	ep := rt.endpoint
	opts := agent.Options{
		Provider: prov,
		Tools:    rt.bridge,
		Sink:     consoleSink{out: rt.stdout, stream: rt.stream, publish: rt.publishPanelSnapshot},
		// The loop's own tool_call rows are deliberately NOT booked: the
		// bridge already writes the authoritative one (assessed level, gate
		// decision, outcome, correlation id). Two rows per call would make
		// tool_call a guess instead of a record.
		Journal: nil,
		// D47: only a task the TEXT loop registered may reach a gate, and the
		// loop owns its task id, so registration travels through this hook.
		// The same hook is this task's only per-task boundary in this
		// composition root, so it also owns the C25 teardown (admitTask).
		AdmitTask: rt.admitTask,
		Registry:  observe.NewRegistry(),
		Config: agent.Config{
			Model:         info.Model,
			ContextWindow: ep.ContextWindow,
			ArtifactsDir:  filepath.Join(rt.spec.dataDir, "artifacts"),
			PerToolTimeout: time.Duration(cfg.Agent.PerToolTimeoutMS) *
				time.Millisecond,
			SteeringEnabled: cfg.Agent.SteeringEnabled,
		},
	}
	// ticket 197: task.spawn reads this snapshot to build its child loop, so a
	// subagent is assembled by the SAME host wiring as its parent (same provider,
	// same admission hook, same budgets) and not by a runtime invented here.
	rt.loopOpt, rt.loopOptSet = opts, true
	loop, err := agent.New(opts)
	if err != nil {
		fmt.Fprintf(rt.stderr, "wisp run: agent 环路装配失败：%v\n", err)
		return 1
	}
	// ticket 200: read the project's own instruction files (AGENTS.md and its
	// siblings) and attach them to the loop assembled just above, so what the
	// repository says about itself reaches the request instead of living in a
	// struct nobody fills. Every input comes from an owner that already exists:
	// the workspace from the view the path resolver feeds (through the account
	// below, not past it), the global tier's directory from the layout decider
	// (the same dir config.toml lives in), the budget from this loop's own scaled
	// D39 prompt total, the token estimator from the loop's package, and the C25
	// engine from the one this process has.
	// The switch is agent.project_instructions_enabled, default on; off means the
	// loader reads nothing and prints why instead of going silent.
	//
	// 200-r2's two changes, both load-bearing:
	//   - the ticket 102 rewrite account is consumed BEFORE a tree is named. This
	//     file hands the whole WorkspaceView to panel and never reaches past the
	//     account for its path, so a %VAR% or ~ rewrite cannot decide which
	//     repository writes the model's system prompt;
	//   - the loader is kept on the runtime so the panel pump has a READER for
	//     the carrier ticket 200 AC#7 asked for. At r1 the carrier existed and
	//     nothing filled it - the wire's version of 文案在、控件不在 (ledger A408).
	instrReq := panel.ProjectInstructionLoadRequestFor(rt.workspaceView(), rt.spec.dataDir)
	if instrReq.Refused != "" {
		rt.auditf("wisp run: %s", instrReq.Refused)
	}
	instrScope := rt.c25.OpenScope("run:project-instructions")
	defer instrScope.Close()
	rt.auditf("wisp run: 项目说明加载器 workspace=%q globalDir=%q enabled=%v budget=%d",
		instrReq.WorkspaceDir, instrReq.GlobalDir, cfg.Agent.ProjectInstructionsEnabled,
		loop.Budgets().PromptTotal)
	instrLoader := projctx.New(projctx.Options{
		WorkspaceDir: instrReq.WorkspaceDir,
		GlobalDir:    instrReq.GlobalDir,
		Enabled:      cfg.Agent.ProjectInstructionsEnabled,
		Refused:      instrReq.Refused,
		BudgetTokens: loop.Budgets().PromptTotal,
		Tokenize:     agent.ApproxTokens,
		Prov:         rt.c25,
		ScopeID:      instrScope.ID(),
		// A format verb, never the line itself: auditf takes a format, and a
		// file path is data, not a directive.
		Log: func(s string) { rt.auditf("projctx: %s", s) },
	})
	rt.setInstructionLoader(instrLoader)
	loop.AttachProjectInstructions(instrLoader)

	// The task's own deadline rides on the parent named above (ticket 246 AC#7's
	// ruling 2.3). For `wisp run` that is context.Background(), which is what this
	// line always produced. For the resident leg it is a root under the ctx that
	// D38(e) step 3 cancels, so "leaving cancels the running task" is a mechanism
	// and not a sentence.
	baseCtx := parent
	if baseCtx == nil {
		baseCtx = rt.spec.taskCtx
	}
	if baseCtx == nil {
		baseCtx = context.Background()
	}
	ctx, cancel := context.WithTimeout(baseCtx,
		time.Duration(cfg.LLM.TimeoutMS)*time.Millisecond)
	defer cancel()
	// 起跑口（票 176-r1）：这一程的任务走**现成**的异步入口
	// internal/agent/loop.go:321，而不是新造一枚同类函数。三件事是这条线的全部要求：
	//   ① id 仍由环路铸（RunningTask.ID / res.TaskID）。组合根自造 id 直接上桥
	//      是上面那段说明里已经登记过的未收口开放端，不走那条路。
	//   ② 名册回填只在 Wait() 之后：那枚 Wait 同时 join 结果信号与登记册句柄，
	//      所以任务真的收口（D38e 的 join 计数归零）以后才有人写它的输出。
	//   ③ 写之前看一眼取消态：被取消掉的任务不再回填（票 176 AC#3 的正方；
	//      判据本体与反向读数在 internal/tools/ticket176r1_start_port_test.go）。
	// 回填走真落盘那套：spills 用的就是环路自己那枚 D15(3) 写入器（同一
	// artifacts 目录、同一份缩放后的预算），产物 key 带 agent-task- 前缀，
	// 不与模型 supplied 的 call id 共用命名空间（internal/tools/task_backfill.go）。
	bg := loop.RunAsync(ctx, task)
	// The root's own identity goes into the roster the moment its id exists, so
	// the tree task.spawn builds has something to hang under (kind=root, no
	// parent). State is deliberately NOT filed here: that producer rewiring is
	// ticket 196's per ruling A394, so StateAnswer keeps saying 「宿主没有登记这一维」
	// for a root row and only subagent rows carry a D43 name.
	rt.tasks.MarkRoot(bg.ID, task)
	res := bg.Wait()
	_, why := (tools.TaskBackfill{
		Roster: rt.tasks,
		Spills: agent.NewSpiller(filepath.Join(rt.spec.dataDir, "artifacts"), loop.Budgets()),
	}).Backfill(bg.Root().Ctx, res.TaskID, res.Text)
	if why != "" {
		rt.auditf("wisp run: 后台任务的输出没有进名册：%s", why)
	}

	class := ""
	if res.Err != nil {
		class = string(res.Err.Class)
	}
	state := taskLogState(res.Status)
	if rt.store != nil {
		c, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel2()
		tl := memory.TaskLog{
			ID: res.TaskID, State: state, QueryText: task,
			CostTokensIn: int64(res.Usage.InputTokens), CostTokensOut: int64(res.Usage.OutputTokens),
			CostAmountMicro: res.CostMicros, Currency: res.Currency, ErrorClass: class,
		}
		if err := rt.store.StartTaskLog(c, tl); err != nil {
			rt.auditf("wisp run: task_log open failed: %v", err)
		}
		summary := res.Text
		if len(summary) > 400 {
			summary = summary[:400]
		}
		if err := rt.store.FinishTaskLog(c, res.TaskID, state, summary,
			int64(res.Usage.InputTokens), int64(res.Usage.OutputTokens),
			res.CostMicros, class); err != nil {
			rt.auditf("wisp run: task_log close failed: %v", err)
		}
	}

	code := runExitCode(res.Status, class)
	fmt.Fprintln(rt.stdout)
	fmt.Fprintf(rt.stdout, "wisp run: 任务 %s 结束（%s，%d 轮，%d 次工具调用，成本 %s）\n",
		res.TaskID, res.Status, res.Rounds, res.ToolCalls, costLine(res))
	if res.Message != "" {
		fmt.Fprintf(rt.stdout, "wisp run: %s\n", res.Message)
	}
	if res.Text != "" && res.Status == agent.StatusCompleted {
		// The reply itself already streamed; this line is the closure.
		fmt.Fprintf(rt.stdout, "wisp run: 回复已完成\n")
	}

	// D10: the result arrives as a notification too, because the CLI has no
	// ball to pop. A poster failure is reported and never swallowed.
	if err := rt.notify("Wisp", notifyBody(res)); err != nil {
		fmt.Fprintf(rt.stderr, "wisp run: 通知未能发出：%v\n", err)
		if code == 0 {
			code = 1
		}
	}
	if code != 0 {
		fmt.Fprintf(rt.stderr, "wisp run: 退出码 %d（error_class=%s）\n", code, orDefaultS(class, "无"))
	}
	return code
}

// storeHealthSink is the composition-side adapter llm.HealthSink ->
// memory.Store (the exact mapping internal/llm documents but must not own).
type storeHealthSink struct{ store *memory.Store }

func (m storeHealthSink) RecordProbe(ctx context.Context, provider, model string,
	flags llm.MeasuredFlags, ok bool, latencyMS int64, at time.Time,
) error {
	return m.store.UpsertProviderProbe(ctx, provider, model, memory.ProbeFlags{
		Text: flags.Text, Vision: flags.Vision, AudioIn: flags.AudioIn,
		AudioOut: flags.AudioOut, Thinking: flags.Thinking, FC: flags.FC,
	}, ok, latencyMS, at)
}

// notifyBody is the notification text: the reply when there is one, the
// failure line when there is not (SPEC-05 §3.4: a failure the user cannot see
// is a failure that did not happen).
func notifyBody(res agent.Result) string {
	if res.Text != "" {
		return res.Text
	}
	if res.Message != "" {
		return res.Message
	}
	return "任务结束：" + string(res.Status)
}

// costLine renders the C23 cost line without floats.
func costLine(res agent.Result) string {
	if res.CostMicros == 0 {
		return strconv.FormatInt(0, 10) + " " + orDefaultS(res.Currency, "CNY") + "（未计价）"
	}
	return fmt.Sprintf("%d micro-%s（in %d / out %d tokens）",
		res.CostMicros, orDefaultS(res.Currency, "CNY"),
		res.Usage.InputTokens, res.Usage.OutputTokens)
}

// taskLogState maps the loop's terminal state onto the memory.task_log.state
// vocabulary.
func taskLogState(s agent.Status) string {
	switch s {
	case agent.StatusCompleted:
		return "done"
	case agent.StatusCancelled:
		return "cancelled"
	case agent.StatusControl:
		return "done"
	default:
		return "error"
	}
}

// chainOptions carries the [net] and [llm.retry] settings into the provider
// stack; tests reach the same function through the config file, so there is
// no second construction path to keep honest.
func chainOptions(cfg *config.Config) llm.ChainBuildOptions {
	return llm.ChainBuildOptions{
		ProxyMode: cfg.Net.Proxy.Mode,
		ProxyURL:  cfg.Net.Proxy.URL,
		RetryBase: time.Duration(cfg.LLM.Retry.BackoffMS) * time.Millisecond,
		RetryMax:  cfg.LLM.Retry.Max,
	}
}

// orDefaultS renders v, or def when v is empty.
func orDefaultS(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

func firstArg(argv []string) string {
	if len(argv) == 0 {
		return ""
	}
	return argv[0]
}

// buildEnvString is the environment badge used to pick the data dir.
func buildEnvString() string { return buildinfo.EnvString() }

// consoleSink streams the reply. It must never block (the loop publishes
// synchronously on the task goroutine).
type consoleSink struct {
	out io.Writer
	// stream collects the reply's deltas into the snapshot's results section, and
	// publish puts a packet on the wire where a visible state moved. Both are nil
	// on a sink built without the pump (see assembleRuntime).
	stream  *panel.StreamLog
	publish func()
}

func (c consoleSink) Publish(e agent.Event) {
	changed := false
	switch e.Kind {
	case agent.EvTextDelta:
		fmt.Fprint(c.out, e.Text)
		if c.stream != nil {
			c.stream.Append(e.TaskID, e.Text)
		}
	case agent.EvReasoningDelta:
		// Reasoning is shown as it is what the thinking probe measures; it is
		// never mixed into the reply text. It is equally kept out of the
		// snapshot's results section, whose contract is an assistant RESULT
		// (panel.ResultChunk, composer.go:51-56); the reasoning field is
		// ticket 145 row 3, which is a key the snapshot does not have.
		fmt.Fprint(c.out, e.Text)
	case agent.EvToolStart:
		fmt.Fprintf(c.out, "\n[工具 %s]\n", e.ToolName)
		changed = true
	case agent.EvToolEnd:
		fmt.Fprintf(c.out, "[工具 %s -> %s]\n", e.ToolName, e.Outcome)
		changed = true
	case agent.EvStuck:
		fmt.Fprintf(c.out, "\n[停滞] %s\n", e.Text)
		changed = true
	case agent.EvError:
		if e.Err != nil {
			fmt.Fprintf(c.out, "\n[错误 %s] %s\n", e.Err.Class, e.Err.Detail)
		}
		changed = true
	case agent.EvDone:
		// The reply is over; the packet the panel renders has to carry that, or
		// the result region would stream forever on a finished task.
		if c.stream != nil {
			c.stream.Close(e.TaskID)
		}
		changed = true
	}
	if changed && c.publish != nil {
		c.publish()
	}
}

// consoleApprovalUI is the injected presentation surface for a console run: it
// prints the card and the countdown events.
//
// It still never answers on the user's behalf - it decides nothing, and the
// gate's own comment on that rule is unchanged. What ticket 201 added is the
// other half of a native surface: it HOLDS what an answer needs. The card's
// single-use grant arrives in Prompt and is booked into live, the ledger the
// reply surface spends from, and the card's id is printed because an operator
// cannot answer a question they cannot address. An L2 card on a run with no
// reply source attached still ends in an auto-reject, and that is the route's
// own behaviour rather than this surface answering for anyone.
type consoleApprovalUI struct {
	out io.Writer
	// run is the assembly this surface belongs to, for the one thing a surface
	// cannot do alone: book the waiting state a card puts this process into
	// (ticket 201 AC#6). nil means a standalone surface, which prints cards and
	// names nothing.
	run *agentRuntime
	// cards counts how many confirmations this surface actually displayed, so a
	// test can tell "the gate opened a window" from "the call never reached the
	// gate".
	mu    sync.Mutex
	cards int
	// live is the native-side ledger of displayed cards (ticket 201). nil means
	// this surface stands alone, exactly as every pre-201 console run did: it
	// shows cards and holds nothing that could answer one.
	live *nativeCards
	// publish puts one snapshot on the wire after the approval state moved. The
	// assembly root attaches it (assembleRuntime); nil means this surface is
	// standing alone, which is what every pre-ticket-35 console run was.
	publish func()
}

// shown reports how many prompts have been displayed.
func (u *consoleApprovalUI) shown() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.cards
}

// Prompt prints one confirmation. Returning nil means "the card was shown",
// not "the user agreed": the gate decides that, and only the gate.
func (u *consoleApprovalUI) Prompt(_ context.Context, p approval.Prompt) error {
	u.mu.Lock()
	u.cards++
	u.mu.Unlock()
	// The grant is booked before anything is printed. The order is the point:
	// from the moment a card is on screen someone may answer it, and the answer
	// needs the proof this surface was handed - which is also why the ledger
	// lives here and nowhere the page can reach (approval.go's grant comment,
	// and ticket 201's 「允许只长在原生侧」).
	if u.live != nil {
		u.live.record(p)
		// AC#6 (ticket 201): the instant a card exists is the instant someone is
		// being waited on, so that fact is booked HERE, on the same statement
		// that made the card answerable, and not inferred downstream. Printed
		// before the card's own lines, because the reading has to exist even if
		// the printer is the thing that fails next.
		if u.run != nil {
			u.run.bookWaitingState("ui-prompt")
		}
	}
	fmt.Fprintf(u.out, "\n[确认 %s %s] %s\n", p.Level, p.Tool, p.Reason)
	for _, r := range p.RulesHit {
		fmt.Fprintf(u.out, "  规则 %s\n", r)
	}
	for _, c := range p.Channels {
		fmt.Fprintf(u.out, "  取消方式：%s（%s）\n", c.Text, channelState(c))
	}
	if len(p.Paths) > 0 {
		fmt.Fprintf(u.out, "  影响路径：%s\n", strings.Join(p.Paths, ", "))
	}
	fmt.Fprintf(u.out, "  窗口 %.1fs\n", p.Window.Seconds())
	// The address the operator answers against. Without it C18's 「回复按
	// correlationId 路由」 is a rule a terminal user cannot obey, and a reply
	// that names nothing is a reply the gate has to refuse as an unknown
	// correlation. Printed last so the rule lines above it keep their shape for
	// the cross-surface card comparison in panel_pump_test.go.
	if u.live != nil {
		fmt.Fprintf(u.out, "  卡片编号：%s\n", p.CorrelationID)
	}
	// The card is on screen, so the queue now has one more pending item than it
	// did a microsecond ago: this is the first moment a snapshot about this task
	// is worth sending.
	if u.publish != nil {
		u.publish()
	}
	return nil
}

// Update prints the transient states (countdown, warning, dismissal, handoff).
func (u *consoleApprovalUI) Update(_ context.Context, e approval.Event) error {
	fmt.Fprintf(u.out, "[%s] %s\n", e.Kind, e.Text)
	// Both kinds that end a card also end its answerability: a dismissed card
	// left the queue, and a started one handed off to execution, so a reply
	// arriving afterwards is a lost vote either way and the ledger must not keep
	// a spendable grant for it (queue.revokeGrants would burn it regardless).
	switch e.Kind {
	case approval.EventDismissed, approval.EventStarted:
		if u.live != nil {
			u.live.forget(e.CorrelationID)
		}
		// The other half of AC#6: the wait is over, and the reading has to say
		// so rather than leaving a stale 「等人」 standing in the log.
		if u.run != nil {
			u.run.bookWaitingState("ui-" + string(e.Kind))
		}
	}
	// Only the two kinds that move the queue publish a packet. A tick or a
	// warning changes a countdown the four-key snapshot has no field for, so
	// publishing on it would book the same packet over and over while claiming
	// nothing new - and inventing that field is Q-51, not this ticket.
	if e.Kind == approval.EventDismissed || e.Kind == approval.EventStarted {
		if u.publish != nil {
			u.publish()
		}
	}
	return nil
}

func channelState(c approval.ChannelStatus) string {
	if c.Loaded {
		return "已加载"
	}
	return "未加载"
}
