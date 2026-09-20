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
		s.dataDir = resolveDataDir(buildEnvString())
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
	spec     runSpec
	cfg      *config.Config
	paths    *tools.PathCanonicalizer
	store    *memory.Store
	provs    []llm.LlmProvider // built chain, kept for the probe path
	names    []string
	endpoint llm.Endpoint
	gate     *approval.Gate
	ui       *consoleApprovalUI
	bridge   *tools.Bridge
	notify   notifyPoster
	stdout   io.Writer
	stderr   io.Writer
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
	cfg, _, err := config.LoadFile(cfgPath, nil)
	if err != nil {
		fmt.Fprintf(s.stderr, "wisp run: 配置未就绪（Unconfigured）：%v\n", err)
		return rt, 2
	}
	rt.cfg = cfg

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

	// The approval layer (A13). Channels loaded here are the honest set for a
	// console run: no floating ball, no global Esc hook, so an L1 window can
	// only expire unvetoed and an L2 card can only time out into a reject.
	rt.ui = &consoleApprovalUI{out: s.stdout}
	rt.gate = approval.New(approval.Options{
		UI:              rt.ui,
		Channels:        approval.NewChannels(),
		Window:          time.Duration(cfg.Risk.L1WindowSec) * time.Second,
		ApprovalTimeout: time.Duration(cfg.Risk.ConfirmTimeoutSec) * time.Second,
		Logf:            rt.auditf,
	})
	rt.bridge = tools.New(tools.Options{
		Registry: reg,
		Paths:    rt.paths,
		Gate:     rt.gate,
		Cancel:   rt.gate.ToolsCancelBus(),
		Journal:  mem,
		// R4 stays dormant: probing every result for sensitive sources needs
		// the C25 detector's own wiring (ticket 25), and a dormant R4 is the
		// honest state, not a weakened one.
		Provenance:     risk.NewProvenance(risk.ProvOptions{NoProbe: true}),
		DefaultTimeout: time.Duration(cfg.Agent.PerToolTimeoutMS) * time.Millisecond,
		OnDecision: func(d tools.Decision) {
			rt.auditf("wisp run: 风险判定 tool=%s level=%s rules=%v reason=%q",
				d.Tool, d.LevelString(), d.RulesHit, d.Reason)
		},
		Logf: rt.auditf,
	})
	return rt, 0
}

// windowCount reports how many confirmation cards the composed gate displayed.
func (rt *agentRuntime) windowCount() int {
	if rt.ui == nil {
		return 0
	}
	return rt.ui.shown()
}

// close releases the store.
func (rt *agentRuntime) close() {
	if rt.store != nil {
		_ = rt.store.Close()
	}
}

// logf is the audit sink: the bridge's and the gate's lines.
func (rt *agentRuntime) auditf(format string, args ...any) {
	fmt.Fprintf(rt.stderr, "[audit] "+format+"\n", args...)
}

// execute runs one task through the loop and presents the result.
func (rt *agentRuntime) execute(task string) int {
	cfg := rt.cfg
	prov := rt.provs[0]
	info := prov.Info()
	ep := rt.endpoint
	loop, err := agent.New(agent.Options{
		Provider: prov,
		Tools:    rt.bridge,
		Sink:     consoleSink{out: rt.stdout},
		// The loop's own tool_call rows are deliberately NOT booked: the
		// bridge already writes the authoritative one (assessed level, gate
		// decision, outcome, correlation id). Two rows per call would make
		// tool_call a guess instead of a record.
		Journal: nil,
		// D47: only a task the TEXT loop registered may reach a gate, and the
		// loop owns its task id, so registration travels through this hook.
		AdmitTask: rt.gate.AdmitTextTask,
		Registry:  observe.NewRegistry(),
		Config: agent.Config{
			Model:         info.Model,
			ContextWindow: ep.ContextWindow,
			ArtifactsDir:  filepath.Join(rt.spec.dataDir, "artifacts"),
			PerToolTimeout: time.Duration(cfg.Agent.PerToolTimeoutMS) *
				time.Millisecond,
			SteeringEnabled: cfg.Agent.SteeringEnabled,
		},
	})
	if err != nil {
		fmt.Fprintf(rt.stderr, "wisp run: agent 环路装配失败：%v\n", err)
		return 1
	}
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(cfg.LLM.TimeoutMS)*time.Millisecond)
	defer cancel()
	res := loop.Run(ctx, task)

	class := ""
	if res.Err != nil {
		class = string(res.Err.Class)
	}
	state := taskLogState(res.Status)
	if rt.store != nil {
		c, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel2()
		tl := memory.TaskLog{ID: res.TaskID, State: state, QueryText: task,
			CostTokensIn: int64(res.Usage.InputTokens), CostTokensOut: int64(res.Usage.OutputTokens),
			CostAmountMicro: res.CostMicros, Currency: res.Currency, ErrorClass: class}
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
	flags llm.MeasuredFlags, ok bool, latencyMS int64, at time.Time) error {
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
type consoleSink struct{ out io.Writer }

func (c consoleSink) Publish(e agent.Event) {
	switch e.Kind {
	case agent.EvTextDelta:
		fmt.Fprint(c.out, e.Text)
	case agent.EvReasoningDelta:
		// Reasoning is shown as it is what the thinking probe measures; it is
		// never mixed into the reply text.
		fmt.Fprint(c.out, e.Text)
	case agent.EvToolStart:
		fmt.Fprintf(c.out, "\n[工具 %s]\n", e.ToolName)
	case agent.EvToolEnd:
		fmt.Fprintf(c.out, "[工具 %s -> %s]\n", e.ToolName, e.Outcome)
	case agent.EvStuck:
		fmt.Fprintf(c.out, "\n[停滞] %s\n", e.Text)
	case agent.EvError:
		if e.Err != nil {
			fmt.Fprintf(c.out, "\n[错误 %s] %s\n", e.Err.Class, e.Err.Detail)
		}
	}
}

// consoleApprovalUI is the injected presentation surface for a console run: it
// prints the card and the countdown events. It never answers on the user's
// behalf, which is exactly why an L2 card here ends in an auto-reject.
type consoleApprovalUI struct {
	out io.Writer
	// cards counts how many confirmations this surface actually displayed, so a
	// test can tell "the gate opened a window" from "the call never reached the
	// gate".
	mu    sync.Mutex
	cards int
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
	return nil
}

// Update prints the transient states (countdown, warning, dismissal, handoff).
func (u *consoleApprovalUI) Update(_ context.Context, e approval.Event) error {
	fmt.Fprintf(u.out, "[%s] %s\n", e.Kind, e.Text)
	return nil
}

func channelState(c approval.ChannelStatus) string {
	if c.Loaded {
		return "已加载"
	}
	return "未加载"
}
