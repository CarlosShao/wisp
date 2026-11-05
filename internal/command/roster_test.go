package command

// Ticket 213 AC#1, AC#2, AC#3 and AC#5, plus the one premise AC#4 turns on.
//
// The rows here are written so that each acceptance criterion has a case that can
// fail for its OWN reason:
//
//   - AC#1: the production roster is enumerable, slash-less, and EVERY row's door
//     is driven for real - the mode row through panel.ModeWriteHandler +
//     perm.Store, the control rows through a real agent.Loop that is cancelled
//     while a task is running (so "there is an executor" is a measurement, not a
//     comment).
//   - AC#2: a planted command with no door is refused at construction, and so is a
//     renamed copy of a door that already exists. The planted rows are the negative
//     control for the gate; remove the gate and these cases go red.
//   - AC#3: for every ordered pair of档 and both channel states, the roster's
//     verdict is compared against what the PRODUCTION mode door actually does:
//     Refused must be the door's error, NeedsApproval must be exactly one confirm
//     call that the panel does not make, Allowed must be zero asks. A guessed
//     policy in Judge disagrees and reds.
//   - AC#5: every one of those rows must also leave an audit line in every档 (the
//     "三档都写审计" half), and this package is scanned to prove it holds no
//     execution of its own (there is no second path to audit).
//   - AC#4 is NOT asserted green here: the roster is enumerable and the doors are
//     command-line executable, but the page has no commands key, so the ticket's
//     own cell stays for the orchestrator/verifier.

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/perm"
	"github.com/CarlosShao/wisp/internal/risk"
)

// ---------------------------------------------------------------------- AC#1

// TestTheRosterIsEnumerableAndSlashLess213 pins the enumeration's shape and the
// census. The literal name set is deliberate: a hand-added row in
// productionEntries() that is not also declared here reds, and a row that is
// declared here but has no door reds inside New.
func TestTheRosterIsEnumerableAndSlashLess213(t *testing.T) {
	r, err := Production()
	if err != nil {
		t.Fatalf("Production(): %v", err)
	}
	got := r.Names()
	want := []string{"cancel", "mode", "stop"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Names() = %v, want %v (the census: only commands with a door today)", got, want)
	}
	for _, name := range got {
		e, ok := r.Lookup(name)
		if !ok {
			t.Fatalf("Lookup(%q) missed a name Names() returned", name)
		}
		if e.Name != name {
			t.Errorf("Lookup(%q).Name = %q", name, e.Name)
		}
		if strings.HasPrefix(e.Name, "/") {
			t.Errorf("%q carries its slash: the menu draws the slash, the roster's name is slash-less (ticket 213 交付 1)", e.Name)
		}
		if strings.TrimSpace(e.Summary) == "" {
			t.Errorf("%q has no summary: a menu row that cannot say what it does is the 'clicked and nothing happened' shape", e.Name)
		}
		if e.Channel == ChannelPanelMode && strings.TrimSpace(e.Args) == "" {
			t.Errorf("%q is a panel command with no argument shape: the menu cannot tell the user it wants a target", e.Name)
		}
		if e.Channel == ChannelControlWord && e.Args != "" {
			t.Errorf("%q is a control word but declares args %q: D11's words carry no arguments", e.Name, e.Args)
		}
		if e.Channel == ChannelPanelMode && e.Method == "" {
			t.Errorf("%q has no panel method", e.Name)
		}
		if e.Channel == ChannelControlWord && e.Utterance == "" {
			t.Errorf("%q has no control utterance", e.Name)
		}
	}
	// Names are the roster's own spelling: neither the menu's slash nor a
	// capability that is not implemented may answer.
	for _, notThere := range []string{"/mode", "compact", "goal", "plan", "model", "export", "feedback"} {
		if _, ok := r.Lookup(notThere); ok {
			t.Errorf("Lookup(%q) answered: AC#2 says a capability that does not exist must not be in the roster (the ticket's 8 minus the two real ones)", notThere)
		}
	}
	if n := len(r.Entries()); n != len(want) {
		t.Errorf("Entries() = %d rows, Names() = %d", n, len(want))
	}
}

// TestEveryRosterDoorHasAProductionDoor213 is AC#1's "生产路径上有调用方" half, read
// off the disk: both doors must be attached by production code, not only by tests.
func TestEveryRosterDoorHasAProductionDoor213(t *testing.T) {
	root := repoRoot213(t)
	cases := []struct {
		file string
		what string
	}{
		{"cmd/wisp/run.go", "panel.ModeWriteHandler{"},           // the resident panel's mode leg
		{"cmd/wisp/panel_inbound.go", "panel.ModeWriteHandler{"}, // the console leg
		{"cmd/wisp/run.go", "agent.New("},                        // the loop production builds
		{"internal/agent/loop.go", "MatchControl(input)"},        // D11(1): control words before any request
		{"internal/agent/loop.go", "defaultControl"},             // the executor that cancels the running task
	}
	for _, c := range cases {
		text := mustRead213(t, filepath.Join(root, c.file))
		if !strings.Contains(text, c.what) {
			t.Errorf("%s no longer contains %q: a roster row would name a door with no production caller (AC#1)", c.file, c.what)
		}
	}
}

// ---------------------------------------------------------------------- AC#2

// TestACommandWithoutADoorCannotEnterTheRoster213 is the gate's positive control:
// each planted row lacks a real door for a different reason, and New must refuse
// all of them loudly.
func TestACommandWithoutADoorCannotEnterTheRoster213(t *testing.T) {
	modeRow := Command{
		Name: "mode", Summary: "切换权限预设", Args: "目标档位",
		Approval: ApprovalPerRequest, Channel: ChannelPanelMode, Method: panel.MethodModeRequest,
	}
	cases := []struct {
		what    string
		entries []Command
		want    string
	}{
		{
			"a control word nobody parses",
			[]Command{{Name: "compress", Summary: "压缩当前对话", Approval: ApprovalNever, Channel: ChannelControlWord, Utterance: "压缩"}},
			"词表",
		},
		{
			"a panel method the wall does not answer",
			[]Command{{Name: "compact", Summary: "压缩当前对话", Approval: ApprovalNever, Channel: ChannelPanelMode, Method: "panel.compact.request"}},
			"panel.compact.request",
		},
		{
			"no channel at all",
			[]Command{{Name: "goal", Summary: "设定长期目标", Approval: ApprovalNever}},
			"没有执行门",
		},
		{
			"a second command claiming a door that already exists",
			[]Command{modeRow, {
				Name: "preset", Summary: "换个名字叫档位", Args: "目标档位",
				Approval: ApprovalPerRequest, Channel: ChannelPanelMode, Method: panel.MethodModeRequest,
			}},
			"门",
		},
		{
			"a row that cannot say what it does",
			[]Command{{Name: "export", Summary: "   ", Approval: ApprovalNever, Channel: ChannelControlWord, Utterance: "确认"}},
			"说明",
		},
		{
			"a route-shaped name",
			[]Command{{Name: "panel.mode.request", Summary: "伪装成路由", Approval: ApprovalNever, Channel: ChannelControlWord, Utterance: "确认"}},
			"不合法",
		},
		{
			"a slash in the name",
			[]Command{{Name: "/mode", Summary: "带斜杠", Approval: ApprovalNever, Channel: ChannelControlWord, Utterance: "确认"}},
			"不合法",
		},
	}
	for _, tc := range cases {
		r, err := New(tc.entries...)
		if err == nil {
			t.Errorf("%s: New accepted it (rows: %v) - AC#2 requires a command with no capability to be unrepresentable", tc.what, r.Names())
			continue
		}
		if r != nil {
			t.Errorf("%s: New returned a roster AND an error", tc.what)
		}
		if !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: refusal %q does not name %q", tc.what, err, tc.want)
		}
	}
}

// TestTheRosterRefusesLoudOnDuplicateNamesAndDoors213 is the "撞名要响亮报" half
// (the ticket picked Pi's form over openchamber's silent ignore).
func TestTheRosterRefusesLoudOnDuplicateNamesAndDoors213(t *testing.T) {
	row := Command{Name: "stop", Summary: "停", Approval: ApprovalNever, Channel: ChannelControlWord, Utterance: "停"}
	_, err := New(row, row)
	if err == nil || !strings.Contains(err.Error(), "重名") {
		t.Fatalf("New(row, row) = %v, want a loud duplicate-name refusal", err)
	}
	if !strings.Contains(err.Error(), `"stop"`) {
		t.Errorf("the duplicate refusal does not name the name that clashed: %v", err)
	}

	// Two different names on ONE utterance is the other half: the roster may not
	// pretend the same D11 word is two commands.
	alias := Command{Name: "halt", Summary: "停", Approval: ApprovalNever, Channel: ChannelControlWord, Utterance: "停"}
	if _, err := New(row, alias); err == nil || !strings.Contains(err.Error(), "门") {
		t.Fatalf("New(one word, two names) = %v, want a door-collision refusal", err)
	}
}

// ---------------------------------------------------------------------- AC<3

// rosterManager is the smallest perm.ConfigManager a real Store can be handed; it
// is the same stand-in the panel's own gate tests use, so no file is touched.
type rosterManager struct{ mode risk.Mode }

func (m *rosterManager) Config() *config.Config {
	c := config.NewDefaults()
	c.Risk.PermissionMode = m.mode.String()
	return c
}

func (m *rosterManager) SetPermissionMode(mode risk.Mode) error {
	m.mode = mode
	return nil
}

// modeLadder213 is the vocabulary these rows iterate: risk.ModeNames() parsed back,
// in the order risk.Mode declares them.
func modeLadder213(t *testing.T) []risk.Mode {
	t.Helper()
	names := risk.ModeNames()
	out := make([]risk.Mode, 0, len(names))
	for _, n := range names {
		m, err := risk.ParseMode(n)
		if err != nil {
			t.Fatalf("risk.ParseMode(%q): %v", n, err)
		}
		out = append(out, m)
	}
	return out
}

// TestAvailabilityAgreesWithTheProductionModeDoor213 is AC#3's teeth. It compares
// the roster's verdict with what the REAL panel handler in front of a REAL
// perm.Store does, for every ordered pair of档 and both channel states. It also
// carries AC#5's audit half: each row must leave an audit line.
func TestAvailabilityAgreesWithTheProductionModeDoor213(t *testing.T) {
	r, err := Production()
	if err != nil {
		t.Fatalf("Production(): %v", err)
	}
	ctx := context.Background()
	ladder := modeLadder213(t)
	rows, auditLines := 0, 0
	verdicts := map[Verdict]int{}

	for _, configured := range ladder {
		for _, target := range ladder {
			for _, card := range []bool{false, true} {
				rows++
				name := fmt.Sprintf("%s-to-%s-card-%v", configured, target, card)

				mgr := &rosterManager{mode: configured}
				var storeAudits []string
				var storeAsks int32
				var storeConfirm perm.ConfirmFunc
				var panelAsks int32
				var panelConfirm panel.ModeConfirm
				if card {
					storeConfirm = func(_ context.Context, sw perm.Switch) (perm.Card, error) {
						atomic.AddInt32(&storeAsks, 1)
						return perm.AnsweredCard("card-" + name), nil
					}
					panelConfirm = func(_ context.Context, _, _ risk.Mode, _, _ string) (perm.Card, error) {
						atomic.AddInt32(&panelAsks, 1)
						return perm.AnsweredCard("card-" + name), nil
					}
				}
				store, err := perm.New(perm.Options{
					Manager: mgr,
					Confirm: storeConfirm,
					Logf: func(format string, args ...any) {
						storeAudits = append(storeAudits, fmt.Sprintf(format, args...))
					},
				})
				if err != nil {
					t.Fatalf("%s: perm.New: %v", name, err)
				}
				var handlerAudits []string
				h := &panel.ModeWriteHandler{
					Modes:   store,
					Confirm: panelConfirm,
					Audit:   func(format string, args ...any) { handlerAudits = append(handlerAudits, fmt.Sprintf(format, args...)) },
					Actor:   "roster-213",
				}

				// The ENV carries the EFFECTIVE档, which is what the handler's own
				// gate reads (a legless store answers the strictest档 even when the
				// file says otherwise: ticket 372-r2).
				env := Env{Mode: store.PermissionMode(), CardChannel: card}
				av := r.Judge(env, Request{Name: "mode", Target: target.String()})
				verdicts[av.Verdict]++

				if strings.TrimSpace(av.Reason) == "" {
					t.Errorf("%s: verdict %v carries no reason (AC#3: 拒了要说为什么)", name, av.Verdict)
				}

				doorErr := h.HandleModeRequest(ctx, panel.ComposerRequest{
					Method: panel.MethodModeRequest, RequestID: name,
					Source: panel.ComposerRequestSource, To: target.String(),
				})

				switch av.Verdict {
				case Refused:
					if doorErr == nil {
						t.Errorf("%s: roster says Refused (%s) but the door accepted the switch", name, av.Reason)
					}
					if mgr.mode != configured {
						t.Errorf("%s: roster says Refused but the persisted档 moved to %s", name, mgr.mode)
					}
				case NeedsApproval:
					if doorErr != nil {
						t.Errorf("%s: roster says NeedsApproval but the door refused: %v", name, doorErr)
					}
					if got := atomic.LoadInt32(&storeAsks); got != 1 {
						t.Errorf("%s: roster says NeedsApproval but perm.Store asked %d time(s), want exactly 1", name, got)
					}
					if got := atomic.LoadInt32(&panelAsks); got != 0 {
						t.Errorf("%s: the panel raised %d card(s) of its own; the panel deliberately raises no second card", name, got)
					}
					if mgr.mode != target {
						t.Errorf("%s: NeedsApproval was answered allow but the persisted档 is %s, want %s", name, mgr.mode, target)
					}
				case Allowed:
					if doorErr != nil {
						t.Errorf("%s: roster says Allowed (%s) but the door refused: %v", name, av.Reason, doorErr)
					}
					if got := atomic.LoadInt32(&storeAsks); got != 0 {
						t.Errorf("%s: roster says Allowed but perm.Store asked for a card %d time(s)", name, got)
					}
					if mgr.mode != target {
						t.Errorf("%s: Allowed but the persisted档 is %s, want %s", name, mgr.mode, target)
					}
				default:
					t.Fatalf("%s: unknown verdict %v", name, av.Verdict)
				}

				// AC#5: whichever way it went, an audit line exists. perm.Store
				// books applied/refused/no-change; the panel books the refusals it
				// makes before the store is reached.
				lines := len(storeAudits) + len(handlerAudits)
				auditLines += lines
				if lines == 0 {
					t.Errorf("%s: verdict %v left no audit line (AC#5: 三档都写审计，不留旁路)", name, av.Verdict)
				}
			}
		}
	}
	if rows != len(ladder)*len(ladder)*2 {
		t.Fatalf("measured %d rows, want %d", rows, len(ladder)*len(ladder)*2)
	}
	// All three verdicts must actually appear, or the agreement is vacuous.
	for _, v := range []Verdict{Allowed, NeedsApproval, Refused} {
		if verdicts[v] == 0 {
			t.Errorf("verdict %v never appeared in %d rows: this matrix is not exercising it", v, rows)
		}
	}
	if auditLines < rows {
		t.Errorf("audit lines = %d for %d rows, want at least one per row", auditLines, rows)
	}
	t.Logf("rows=%d verdicts=allowed:%d needs-approval:%d refused:%d audit_lines=%d",
		rows, verdicts[Allowed], verdicts[NeedsApproval], verdicts[Refused], auditLines)
}

// TestTheRosterRefusesWhatItDoesNotKnow213 pins the "never silent" half for the
// two inputs that never reach a door: an unknown name and an unparseable target.
func TestTheRosterRefusesWhatItDoesNotKnow213(t *testing.T) {
	r, err := Production()
	if err != nil {
		t.Fatalf("Production(): %v", err)
	}
	cases := []struct {
		req  Request
		want string
	}{
		{Request{Name: "compact"}, "名册里没有"},
		{Request{Name: "mode"}, "目标档位"},
		{Request{Name: "mode", Target: "god_mode"}, "不可解析"},
	}
	for _, tc := range cases {
		av := r.Judge(Env{Mode: risk.ModeAskEveryStep}, tc.req)
		if av.Verdict != Refused {
			t.Errorf("Judge(%+v) = %v, want Refused", tc.req, av.Verdict)
		}
		if !strings.Contains(av.Reason, tc.want) || strings.TrimSpace(av.Reason) == "" {
			t.Errorf("Judge(%+v) reason = %q, want it to name %q", tc.req, av.Reason, tc.want)
		}
	}
	// The control rows are runnable under every档: stopping is not a privilege.
	ladder := modeLadder213(t)
	for _, m := range ladder {
		for _, name := range []string{"cancel", "stop"} {
			if av := r.Judge(Env{Mode: m}, Request{Name: name}); av.Verdict != Allowed {
				t.Errorf("Judge(%s, %s) = %v (%s), want Allowed in every档", m, name, av.Verdict, av.Reason)
			}
		}
	}
}

// ---------------------------------------------------------------------- AC#1

// blockingProvider is a C5 stub whose ONLY job is to be cancelled: the layer under
// test here is D11(1)'s local control path, not the provider. It records how many
// streams were opened, because the control word itself must cost none.
type blockingProvider struct {
	started chan struct{}
	once    sync.Once
	calls   int32
}

func newBlockingProvider() *blockingProvider {
	return &blockingProvider{started: make(chan struct{})}
}

func (p *blockingProvider) Stream(ctx context.Context, _ *llm.Request, _ func(llm.StreamEvent) error) error {
	atomic.AddInt32(&p.calls, 1)
	p.once.Do(func() { close(p.started) })
	<-ctx.Done()
	return ctx.Err()
}

func (p *blockingProvider) Info() llm.ProviderInfo {
	return llm.ProviderInfo{Protocol: "test", Model: "blocking"}
}

// TestTheControlRowsStopARealRunningTask213 is AC#1's executor half for the two
// control rows: a real agent.Loop, a real running task, and the real control word
// cancelling it. It also pins that the control word costs zero provider calls,
// which is what makes this a command (local execution) rather than a prompt.
func TestTheControlRowsStopARealRunningTask213(t *testing.T) {
	p := newBlockingProvider()
	loop, err := agent.New(agent.Options{Provider: p})
	if err != nil {
		t.Fatalf("agent.New: %v", err)
	}
	task := loop.RunAsync(context.Background(), "统计仓库里每个文件的字节数")
	if !waitClosed213(t, p.started, 5*time.Second) {
		task.Cancel()
		task.Wait()
		t.Fatalf("the task never reached the provider: the control word would have nothing to cancel")
	}

	res := loop.Run(context.Background(), "取消")
	if res.Status != agent.StatusControl {
		t.Errorf("status = %s, want control (D11(1): the control word never reaches the provider)", res.Status)
	}
	if !strings.Contains(res.Message, "已停止当前任务") {
		t.Errorf("message = %q, want the loop's own '已停止当前任务' (the default handler ran, not a no-op)", res.Message)
	}
	done := task.Wait()
	if done.Status != agent.StatusCancelled {
		t.Errorf("the running task ended as %s, want cancelled: the control word did not reach it", done.Status)
	}
	if got := atomic.LoadInt32(&p.calls); got != 1 {
		t.Errorf("provider streams = %d, want 1 (the cancelled task only): the control word itself must cost no LLM call", got)
	}
}

// waitClosed213 reports whether ch closed within d (a bounded wait, not a sleep).
func waitClosed213(t *testing.T, ch <-chan struct{}, d time.Duration) bool {
	t.Helper()
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ch:
		return true
	case <-timer.C:
		return false
	}
}

// ---------------------------------------------------------------------- AC#5

// TestTheRosterPackageExecutesNothing213 is AC#5's structural half: the package
// holds no store, no loop and no handler, so it cannot be a second execution path.
// If a later leg decides to execute from here, this case reds and forces the
// audit question into the same review.
//
// The scan is over the CALL GRAPH, not the text: this package's own prose and
// refusal sentences name the doors it does not call (that is how they are
// documented), so a text match would red on the explanation instead of on a call.
func TestTheRosterPackageExecutesNothing213(t *testing.T) {
	dir := packageDir213(t)
	// Any of these calls would make this package an execution path of its own.
	forbidden := map[string]string{
		"Set":               "calling a write leg (perm.Store.Set is the only one)",
		"New":               "building a store or a loop (perm.New / agent.New)",
		"Run":               "running a task (agent.Loop.Run)",
		"RunAsync":          "starting a task (agent.Loop.RunAsync)",
		"HandleModeRequest": "calling the panel handler directly",
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	scanned := 0
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		scanned++
		if bad, why := forbiddenCall213(t, filepath.Join(dir, name), forbidden); bad != "" {
			t.Errorf("%s calls %s (%s): the roster must not grow an execution path (AC#5: every execution goes through the existing audited door)", name, bad, why)
		}
	}
	if scanned == 0 {
		t.Fatal("scanned 0 production files: this case would be decoration")
	}
}

// forbiddenCall213 parses one file and reports the first call whose callee's
// selector is in the forbidden set, e.g. "perm.New" for `perm.New(...)`.
func forbiddenCall213(t *testing.T, path string, forbidden map[string]string) (string, string) {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	found, why := "", ""
	ast.Inspect(file, func(n ast.Node) bool {
		if found != "" {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if reason, banned := forbidden[sel.Sel.Name]; banned {
			qual := sel.Sel.Name
			if id, ok := sel.X.(*ast.Ident); ok {
				qual = id.Name + "." + sel.Sel.Name
			}
			found, why = qual, reason
			return false
		}
		return true
	})
	return found, why
}

// --------------------------------------------------------------------- helpers

func packageDir213(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Dir(thisFile)
}

func repoRoot213(t *testing.T) string {
	t.Helper()
	return filepath.Join(packageDir213(t), "..", "..")
}

func mustRead213(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
