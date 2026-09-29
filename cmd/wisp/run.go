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
	// displays (ticket 211). nil means nobody can answer, which is exactly what
	// every run before this field was: an L2 card waits out the C18 deadline and
	// auto-rejects, an L1 window cannot be opposed. cmdRun fills it with the
	// console's input handle when one is really interactive; the CLI tests fill it
	// with a scripted reader, which is the injection seam AGENTS.md §1.3 names
	// for `wisp run` - it is not a mock standing in for a missing subsystem,
	// because the subsystem here IS a human typing at this terminal.
	reply io.Reader
	// replyVeto is the veto channel this host's cancel transport really is
	// (ticket 211). Production leaves it empty because a console run wires none
	// of SPEC-06 §2's four channels, and an empty value makes Gate.Veto say that
	// back instead of letting this process claim a cancel path it does not have.
	// A host that DOES own one - the ball click, the global Esc hook - names it
	// here and marks it loaded in the same step.
	replyVeto approval.Channel
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
	return rt.execute(task)
}

// runtime is the assembled S1 stack.
type agentRuntime struct {
	spec  runSpec
	cfg   *config.Config
	mgr   *config.Manager
	paths *tools.PathCanonicalizer
	store *memory.Store
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
	// one thing that makes an allow possible at all (ticket 211). It is filled
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
	// task.output (ticket 164 AC#3, D34 row PLAN.md:2564, L0) reads the
	// process-local task table above. ⚠ WHO FILLS IT IS NOT THIS TICKET: the
	// background spawn and its cancel are AC#4's and ticket 163's land, and
	// RunAsync now has one production call site (measured: the line below in
	// this same file, :633 `bg := loop.RunAsync(ctx, task)`, landed with
	// ticket 176 AC#1; the count was zero when evidence file
	// 164-task-output-impl-r1-ac2-ac3.md §4 recorded it). Until a spawner
	// records into it,
	// every call here answers "查不到这个任务" - loud, per 票 164 定案②, and
	// never an empty success that would read as "the task printed nothing".
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
	// Ticket 211 changed one half of that sentence and left the other half
	// standing: the console can now ANSWER a card (a reply source, attached at
	// the end of this function), so an L2 card no longer has to be a question
	// nobody hears. What it still cannot do is veto an L1 window, because none
	// of SPEC-06 §2's four veto channels (ball / Esc / KWS / panel) exists in a
	// terminal - so NewChannels() stays empty and runSpec.replyVeto stays unset,
	// and approval_reply.go says what both of those choices protect.
	rt.liveCards = newNativeCards()
	rt.ui = &consoleApprovalUI{out: s.stdout, live: rt.liveCards}
	rt.gate = approval.New(approval.Options{
		UI:              rt.ui,
		Channels:        approval.NewChannels(),
		Window:          time.Duration(cfg.Risk.L1WindowSec) * time.Second,
		ApprovalTimeout: time.Duration(cfg.Risk.ConfirmTimeoutSec) * time.Second,
		Logf:            rt.auditf,
	})

	// The permission mode (ticket 90's storage, wired here by ticket 101). This
	// is the line that makes R20/M3 real: without it the bridge reads a nil
	// ModeSource, which answers the strictest档 for every call, and a manually
	// chosen档 would be silently forgotten at the next start.
	//
	// What is deliberately NOT here: any grant source. tools.Options.Confirmations
	// stays nil, because R20/M3 persists the MODE and nothing else - a D45
	// session grant must not come back from disk just because something on this
	// boot learned to read config.toml (PLAN.md:1640, pinned by AC#2(c)).
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
	rt.ui.publish = rt.publishPanelSnapshot

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
	// runtime (ticket 211). With no reply source this is a no-op and the run
	// keeps the posture it had before the field existed: cards get shown, nobody
	// answers them, and each route resolves on its own clock.
	if s.reply != nil {
		rt.attachReplyListener(s.reply, s.replyVeto)
	}

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

// windowCount reports how many confirmation cards the composed gate displayed.
func (rt *agentRuntime) windowCount() int {
	if rt.ui == nil {
		return 0
	}
	return rt.ui.shown()
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
func (rt *agentRuntime) execute(task string) int {
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

	ctx, cancel := context.WithTimeout(context.Background(),
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
// gate's own comment on that rule is unchanged. What ticket 211 added is the
// other half of a native surface: it HOLDS what an answer needs. The card's
// single-use grant arrives in Prompt and is booked into live, the ledger the
// reply surface spends from, and the card's id is printed because an operator
// cannot answer a question they cannot address. An L2 card on a run with no
// reply source attached still ends in an auto-reject, and that is the route's
// own behaviour rather than this surface answering for anyone.
type consoleApprovalUI struct {
	out io.Writer
	// cards counts how many confirmations this surface actually displayed, so a
	// test can tell "the gate opened a window" from "the call never reached the
	// gate".
	mu    sync.Mutex
	cards int
	// live is the native-side ledger of displayed cards (ticket 211). nil means
	// this surface stands alone, exactly as every pre-211 console run did: it
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
	// and ticket 211's 「允许只长在原生侧」).
	if u.live != nil {
		u.live.record(liveCard{
			CorrelationID: p.CorrelationID,
			Tool:          p.Tool,
			Level:         p.Level,
			Grant:         p.Grant,
			Window:        p.Window,
		})
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
