package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

// Ticket 90 tests: the user-facing permission mode inside the real decision
// chain. Naming convention follows the package's existing ticketNN_test.go
// files (84, 87) so this file cannot collide with another ticket's helpers.
//
// Coverage map (the AC table on the ticket):
//
//	AC#1  TestTicket90ModeReachesRoutingAsParameter,
//	      TestTicket90TwoBridgesDoNotShareOneMode        (the "no global" proof)
//	AC#2  TestTicket90IrreversibleStillAsksInEveryMode   (PLAN.md:1629)
//	      TestTicket90TaintEscalationNeverSilenced       (PLAN.md:1640)
//	      TestTicket90TierADenySurvivesEveryMode         (SPEC-06 §4.1)
//	      TestTicket90ModeCarriesNoAllowAuthority        (PLAN.md:1588, ban #6)
//	AC#4  the anchors these tests bite are named in each test's comment; the
//	      mutations themselves are run by hand and recorded on the ticket.
//	AC#5  TestTicket90BlacklistGateIsCalledOnLivePath

// ---------------------------------------------------------------------------
// fixtures
// ---------------------------------------------------------------------------

// t90Modes is a ModeSource the test controls. It is a POINTER receiver on
// purpose: the bridge holds the interface, so a mutation that turned the mode
// into a package-level global would show up here as the two bridges below
// agreeing when they were told to disagree.
type t90Modes struct {
	mu sync.Mutex
	m  risk.Mode
}

func (s *t90Modes) PermissionMode() risk.Mode {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.m
}

func (s *t90Modes) set(m risk.Mode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m = m
}

// t90Tool is a minimal Tool: it records that it ran and answers ok.
type t90Tool struct {
	name string

	mu   sync.Mutex
	runs int
}

func (t *t90Tool) Name() string        { return t.name }
func (t *t90Tool) Description() string { return "ticket 90 fixture" }
func (t *t90Tool) Parameters() JSONSchema {
	return JSONSchema(`{"type":"object","properties":{"path":{"type":"string"}}}`)
}

func (t *t90Tool) Execute(context.Context, json.RawMessage, func(string)) (Result, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.runs++
	return Result{Text: "ran"}, nil
}

// t90Entry builds one declared-floor + facts-hook entry.
func t90Entry(name string, lvl risk.Level, pathParams []string,
	facts func(context.Context, map[string]any, []string) risk.Facts,
) (Entry, *t90Tool) {
	tool := &t90Tool{name: name}
	return Entry{
		Tool: tool,
		Decl: Decl{
			Capabilities: []Capability{CapFSRead},
			Needs:        []Capability{CapFSRead},
			Declared:     lvl,
			PathParams:   pathParams,
			Facts:        facts,
			Resident:     true,
			Provider:     KindBuiltin,
		},
	}, tool
}

// t90Gate counts which branch the routing reached and what it answered. The
// answers default to "the user let it through", so a test that asserts "asked"
// is asserting a card existed, not a refusal.
type t90Gate struct {
	mu             sync.Mutex
	windows        int
	approvals      int
	seen           []Decision
	windowAnswer   Answer
	approvalAnswer Answer
	windowReason   string
	approvalReason string
}

func (g *t90Gate) PendingWindow(_ context.Context, d Decision) (Answer, string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.windows++
	g.seen = append(g.seen, d)
	if g.windowAnswer != "" {
		return g.windowAnswer, g.windowReason
	}
	return AnswerTimeout, ""
}

func (g *t90Gate) PendingApproval(_ context.Context, d Decision) (Answer, string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.approvals++
	g.seen = append(g.seen, d)
	if g.approvalAnswer != "" {
		return g.approvalAnswer, g.approvalReason
	}
	return AnswerAllow, ""
}

func (g *t90Gate) counts() (int, int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.windows, g.approvals
}

func (g *t90Gate) last() Decision {
	g.mu.Lock()
	defer g.mu.Unlock()
	if len(g.seen) == 0 {
		return Decision{}
	}
	return g.seen[len(g.seen)-1]
}

func t90Req(tool, path string) agent.ToolRequest {
	args := `{}`
	if path != "" {
		args = fmt.Sprintf(`{"path":%q}`, path)
	}
	return agent.ToolRequest{
		TaskID: "t90", CorrelationID: "t90-corr", CallID: "t90-call",
		Name: tool, Args: json.RawMessage(args),
	}
}

func t90Irreversible(_ context.Context, _ map[string]any, _ []string) risk.Facts {
	return risk.Facts{Irreversible: []string{"delete"}}
}

// t90Taint is a fake C25 detector: every call carries a sensitive source.
type t90Taint struct{}

func (t90Taint) TaintHit(map[string]any) (string, bool) { return "ticket90-fake-source", true }

// ---------------------------------------------------------------------------
// AC#1
// ---------------------------------------------------------------------------

// TestTicket90ModeReachesRoutingAsParameter is AC#1's behavioral half: the mode
// decides whether a question is asked, per call, through the injected source.
//
// AC#4 anchor (mutation target for the orchestrator): the switch arms in
// Bridge.route that read sil.Level, and Mode.Screen's L1 branch
// `if d.Level == L1 { return Silenced{Level: L0, Silenced: true, Mode: m} }`
// in internal/risk/mode.go.
func TestTicket90ModeReachesRoutingAsParameter(t *testing.T) {
	modes := &t90Modes{}
	gate := &t90Gate{}
	e, tool := t90Entry("t90.write", risk.L1, nil, nil)
	reg := NewRegistry()
	if err := reg.Register(e); err != nil {
		t.Fatalf("Register: %v", err)
	}
	// The verdict the mode was applied to must be readable even when the mode
	// removed the question, so the assertions below read the emitted stream
	// (OnDecision), not the gate's memory: a silenced call never reaches a gate,
	// and a test that only looked at gates could not tell "silenced" from
	// "never judged".
	var (
		emMu    sync.Mutex
		emitted []Decision
	)
	b := New(Options{
		Registry: reg, Gate: gate, Modes: modes,
		OnDecision: func(d Decision) {
			emMu.Lock()
			defer emMu.Unlock()
			emitted = append(emitted, d)
		},
	})
	last := func() Decision {
		emMu.Lock()
		defer emMu.Unlock()
		if len(emitted) == 0 {
			return Decision{}
		}
		return emitted[len(emitted)-1]
	}

	cases := []struct {
		mode          risk.Mode
		wantWindows   int
		wantExec      bool
		wantSilencing bool
	}{
		// R20/M2: the default档 asks about the L1 window.
		{risk.ModeAskEveryStep, 1, true, false},
		// R20/M1: 只问高危 silences L1 and keeps L2.
		{risk.ModeAskHighRisk, 0, true, true},
		// R20/M1: 全自动 silences L1 too.
		{risk.ModeAutoApprove, 0, true, true},
	}
	for _, tc := range cases {
		modes.set(tc.mode)
		before := gate.windows
		if _, err := b.Execute(context.Background(), t90Req("t90.write", "")); err != nil {
			t.Fatalf("Execute(%s): %v", tc.mode, err)
		}
		got := gate.windows - before
		if got != tc.wantWindows {
			t.Errorf("mode=%s: L1 windows opened = %d, want %d", tc.mode, got, tc.wantWindows)
		}
		d := last()
		if d.Mode != tc.mode {
			t.Errorf("mode=%s: the decision records mode %v", tc.mode, d.Mode)
		}
		if tc.wantSilencing && !d.ModeSilenced {
			t.Errorf("mode=%s: decision carries no ModeSilenced mark; the audit trail would "+
				"not be able to tell a silenced call from an unjudged one", tc.mode)
		}
		if !tc.wantSilencing && d.ModeSilenced {
			t.Errorf("mode=%s: ModeSilenced=true on a mode that asks about everything", tc.mode)
		}
		if d.Level != risk.L1 {
			t.Errorf("mode=%s: booked level = %v, want the ASSESSED L1 (a mode may not rewrite "+
				"what the judge concluded; that is what tool_call.risk_level means)", tc.mode, d.Level)
		}
	}
	if tool.runs != len(cases) {
		t.Errorf("tool runs = %d, want %d (one per mode)", tool.runs, len(cases))
	}
}

// TestTicket90TwoBridgesDoNotShareOneMode is AC#1's "不许做成包级全局变量" half.
// A package-level mode would be process-wide, so two bridges built in the same
// process with different sources would behave identically. They must not.
func TestTicket90TwoBridgesDoNotShareOneMode(t *testing.T) {
	build := func(m risk.Mode) (*Bridge, *t90Gate) {
		gate := &t90Gate{}
		e, _ := t90Entry("t90.l1", risk.L1, nil, nil)
		reg := NewRegistry()
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register: %v", err)
		}
		return New(Options{Registry: reg, Gate: gate, Modes: &t90Modes{m: m}}), gate
	}
	strict, strictGate := build(risk.ModeAskEveryStep)
	relaxed, relaxedGate := build(risk.ModeAutoApprove)

	// Interleave so a leaked global would have to show up in the second call.
	if _, err := strict.Execute(context.Background(), t90Req("t90.l1", "")); err != nil {
		t.Fatalf("strict Execute: %v", err)
	}
	if _, err := relaxed.Execute(context.Background(), t90Req("t90.l1", "")); err != nil {
		t.Fatalf("relaxed Execute: %v", err)
	}
	if sw, _ := strictGate.counts(); sw != 1 {
		t.Errorf("strict bridge asked %d times, want 1 (mode value leaked across bridges)", sw)
	}
	if rw, _ := relaxedGate.counts(); rw != 0 {
		t.Errorf("relaxed bridge asked %d times, want 0", rw)
	}
}

// TestTicket90UnwiredModeSourceIsTheStrictestMode pins the fail-closed default:
// a composition that never injects a source (which is today's cmd/wisp state,
// see the ticket's handoff note) must ask about everything.
func TestTicket90UnwiredModeSourceIsTheStrictestMode(t *testing.T) {
	gate := &t90Gate{}
	e, _ := t90Entry("t90.l1", risk.L1, nil, nil)
	reg := NewRegistry()
	if err := reg.Register(e); err != nil {
		t.Fatalf("Register: %v", err)
	}
	b := New(Options{Registry: reg, Gate: gate}) // no Modes
	if _, err := b.Execute(context.Background(), t90Req("t90.l1", "")); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if w, _ := gate.counts(); w != 1 {
		t.Errorf("windowless composition opened %d L1 windows, want 1 (fail-closed = ask)", w)
	}
	if got := gate.last().Mode; got != risk.ModeAskEveryStep {
		t.Errorf("Mode = %v, want ask_every_step", got)
	}
}

// ---------------------------------------------------------------------------
// AC#2: the three red lines, one test each
// ---------------------------------------------------------------------------

// TestTicket90IrreversibleStillAsksInEveryMode is red line 1 (PLAN.md:1629:
// 不可逆操作必须拒绝，且任何授权与会话授权都无法覆盖). R8 is silence-immune in
// every档, including 全自动.
//
// AC#4 anchor: `case R8:` in risk.redLine (internal/risk/mode.go). Deleting
// that case is exactly the "去掉红线判定" mutation, and this test must go red.
func TestTicket90IrreversibleStillAsksInEveryMode(t *testing.T) {
	for _, m := range []risk.Mode{risk.ModeAskEveryStep, risk.ModeAskHighRisk, risk.ModeAutoApprove} {
		gate := &t90Gate{}
		e, tool := t90Entry("t90.rm", risk.L0, nil, t90Irreversible)
		reg := NewRegistry()
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register: %v", err)
		}
		b := New(Options{Registry: reg, Gate: gate, Modes: &t90Modes{m: m}})
		if _, err := b.Execute(context.Background(), t90Req("t90.rm", "")); err != nil {
			t.Fatalf("mode=%s Execute: %v", m, err)
		}
		if _, a := gate.counts(); a != 1 {
			t.Errorf("mode=%s: irreversible call raised %d L2 cards, want 1", m, a)
		}
		if tool.runs != 1 {
			t.Errorf("mode=%s: tool ran %d times after an approved card, want 1", m, tool.runs)
		}
		if m == risk.ModeAutoApprove {
			if kept := gate.last().ModeKept; !strings.Contains(kept, "R8") {
				t.Errorf("mode=%s: ModeKept = %q, want it to name R8 (audit must say why "+
					"the mode did not apply)", m, kept)
			}
			if gate.last().ModeSilenced {
				t.Errorf("mode=%s: an R8 verdict was marked silenced", m)
			}
		}
	}
}

// TestTicket90TaintEscalationNeverSilenced is red line 2 (PLAN.md:1640: 会话授权
// 不得覆盖 C25 污染升级). The verdict that carries SessionOverrideBlocked stays
// a question in every档, and no mode value can clear the flag.
//
// AC#4 anchor: `if d.SessionOverrideBlocked {` at the head of risk.redLine.
func TestTicket90TaintEscalationNeverSilenced(t *testing.T) {
	for _, m := range []risk.Mode{risk.ModeAskEveryStep, risk.ModeAskHighRisk, risk.ModeAutoApprove} {
		gate := &t90Gate{}
		e, _ := t90Entry("t90.taint", risk.L0, nil, nil)
		reg := NewRegistry()
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register: %v", err)
		}
		assessor := risk.NewRiskAssessor().WithTaintDetector(t90Taint{})
		b := New(Options{Registry: reg, Gate: gate, Modes: &t90Modes{m: m}, Assessor: assessor})
		if _, err := b.Execute(context.Background(), t90Req("t90.taint", "")); err != nil {
			t.Fatalf("mode=%s Execute: %v", m, err)
		}
		if _, a := gate.counts(); a != 1 {
			t.Errorf("mode=%s: tainted call raised %d L2 cards, want 1", m, a)
		}
		if !gate.last().SessionOverrideBlocked {
			t.Errorf("mode=%s: the decision lost R4's no-session-override flag", m)
		}
		if m == risk.ModeAutoApprove && !strings.Contains(gate.last().ModeKept, "R4") {
			t.Errorf("mode=%s: ModeKept = %q, want it to name R4", m, gate.last().ModeKept)
		}
	}
}

// TestTicket90TierADenySurvivesEveryMode is red line 3's deny half: an A-tier
// path (SPEC-06 §4.1, absolute) must never run and never get a card, in any档.
// The profile is relocated so the frozen anchors name a path this test owns.
func TestTicket90TierADenySurvivesEveryMode(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	key := filepath.Join(home, ".ssh", "id_ed25519")
	if err := os.MkdirAll(filepath.Dir(key), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(key, []byte("private"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	canon := t90Canonical(t, key)

	for _, m := range []risk.Mode{risk.ModeAskEveryStep, risk.ModeAskHighRisk, risk.ModeAutoApprove} {
		gate := &t90Gate{approvalAnswer: AnswerAllow, windowAnswer: AnswerTimeout}
		e, tool := t90Entry("t90.a", risk.L1, []string{"path"}, nil)
		reg := NewRegistry()
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register: %v", err)
		}
		b := New(Options{
			Registry: reg, Gate: gate, Modes: &t90Modes{m: m},
			Paths: NewPathCanonicalizer([]string{filepath.Dir(canon)}, nil),
		})
		if _, err := b.Execute(context.Background(), t90Req("t90.a", key)); err != nil {
			t.Fatalf("mode=%s Execute: %v", m, err)
		}
		if w, a := gate.counts(); w != 0 || a != 0 {
			t.Errorf("mode=%s: A-tier path raised windows=%d approvals=%d, want 0/0 "+
				"(a Deny must not become a question)", m, w, a)
		}
		if tool.runs != 0 {
			t.Errorf("mode=%s: A-tier path ran the tool %d times, want 0", m, tool.runs)
		}
	}
}

// TestTicket90ModeCarriesNoAllowAuthority is red line 3's other half
// (PLAN.md:1588: 「允许」只接受原生侧点击). A permission mode can only ever
// REMOVE a question; it can never express an answer. That is a shape property,
// and the machine gate for the panel half is ban #6 (armed by ticket 88), whose
// rule is that no frontend text can stand in for a native click.
//
// This test fails if anyone later hangs an authority field off the decision or
// gives the mode source a write method - the two shapes that would let a mode
// (or whoever can call it) approve a call.
func TestTicket90ModeCarriesNoAllowAuthority(t *testing.T) {
	// 1. ModeSource is read-only.
	src := reflect.TypeOf((*ModeSource)(nil)).Elem()
	if src.NumMethod() != 1 {
		t.Errorf("ModeSource has %d methods, want exactly 1 (read-only)", src.NumMethod())
	}
	if name := src.Method(0).Name; name != "PermissionMode" {
		t.Errorf("ModeSource's only method is %s, want PermissionMode", name)
	}

	// 2. Nothing on Decision can carry an approval, and nothing on
	//    agent.ToolRequest can carry a mode or an answer: the caller of
	//    Execute supplies arguments, never a verdict.
	dTyp := reflect.TypeOf(Decision{})
	for i := 0; i < dTyp.NumField(); i++ {
		f := dTyp.Field(i)
		low := strings.ToLower(f.Name)
		if strings.Contains(low, "allow") || strings.Contains(low, "approve") ||
			strings.Contains(low, "grant") || f.Type == reflect.TypeOf(Answer("")) {
			t.Errorf("tools.Decision field %s looks like an authority the caller could set; "+
				"an allow must only ever come from the gate's own answer (SPEC-06 §9)", f.Name)
		}
	}
	rTyp := reflect.TypeOf(agent.ToolRequest{})
	for i := 0; i < rTyp.NumField(); i++ {
		low := strings.ToLower(rTyp.Field(i).Name)
		if strings.Contains(low, "mode") || strings.Contains(low, "allow") ||
			strings.Contains(low, "approve") || strings.Contains(low, "auto") {
			t.Errorf("agent.ToolRequest field %s lets the CALLER supply a permission mode or an "+
				"answer; the mode has exactly one owner (internal/perm) and an answer has one "+
				"(the native side)", rTyp.Field(i).Name)
		}
	}

	// 3. Behaviorally: 全自动 executing a call is not an approval record. The
	//    gate is never consulted, so it never saw - and never signed - anything.
	gate := &t90Gate{}
	e, tool := t90Entry("t90.l1b", risk.L1, nil, nil)
	reg := NewRegistry()
	if err := reg.Register(e); err != nil {
		t.Fatalf("Register: %v", err)
	}
	b := New(Options{Registry: reg, Gate: gate, Modes: &t90Modes{m: risk.ModeAutoApprove}})
	if _, err := b.Execute(context.Background(), t90Req("t90.l1b", "")); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	w, a := gate.counts()
	if w != 0 || a != 0 {
		t.Errorf("silenced call consulted the gate (%d windows, %d approvals); want 0/0", w, a)
	}
	if tool.runs != 1 {
		t.Errorf("tool runs = %d, want 1", tool.runs)
	}
}

// ---------------------------------------------------------------------------
// AC#5: risk.Gate gets a real production caller
// ---------------------------------------------------------------------------

// TestTicket90BlacklistGateIsCalledOnLivePath proves the call site is on the
// live routing path (not a test-only helper): a B-tier path on a real call both
// fills BlacklistNote and leaves risk's own security log line behind, because
// risk.Gate is the only code that knows the single-file override rule.
//
// It also pins the honesty half: an override record on file changes the NOTE,
// never the routing. R3 stays silence-immune, so even a "confirmed" B-tier file
// still gets its card until ticket 21's flow exists.
func TestTicket90BlacklistGateIsCalledOnLivePath(t *testing.T) {
	dir := t.TempDir()
	pem := filepath.Join(dir, "server.pem")
	if err := os.WriteFile(pem, []byte("material"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	canon := t90Canonical(t, pem)
	root := t90Canonical(t, dir)

	var logs []string
	prev := risk.Logf
	risk.Logf = func(format string, args ...any) {
		logs = append(logs, fmt.Sprintf(format, args...))
	}
	t.Cleanup(func() { risk.Logf = prev })

	run := func(overrides map[string]bool) (BlacklistNote, int, int) {
		gate := &t90Gate{}
		e, _ := t90Entry("t90.b", risk.L0, []string{"path"}, nil)
		reg := NewRegistry()
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register: %v", err)
		}
		b := New(Options{
			Registry: reg, Gate: gate,
			Paths:         NewPathCanonicalizer([]string{root}, nil),
			Confirmations: func() map[string]bool { return overrides },
		})
		if _, err := b.Execute(context.Background(), t90Req("t90.b", pem)); err != nil {
			t.Fatalf("Execute: %v", err)
		}
		w, a := gate.counts()
		return gate.last().Blacklist, w, a
	}

	note, w, a := run(nil)
	if note.Class != risk.ClassB {
		t.Fatalf("Blacklist.Class = %v, want B (risk.Gate was not consulted at all)", note.Class)
	}
	if len(note.Unlockable) != 1 {
		t.Errorf("Unlockable = %v, want the one B-tier path", note.Unlockable)
	}
	if len(note.AlreadyUnlocked) != 0 {
		t.Errorf("AlreadyUnlocked = %v with no confirmations on file, want empty", note.AlreadyUnlocked)
	}
	if w+a != 1 {
		t.Errorf("B-tier path asked %d windows + %d approvals, want exactly one question", w, a)
	}
	var sawGateLog bool
	for _, l := range logs {
		if strings.Contains(l, "B-list") {
			sawGateLog = true
		}
	}
	if !sawGateLog {
		t.Error("risk.Gate produced no security log line, so the call site is not the " +
			"one that logs A denies and B decisions (risk.Logf sink)")
	}

	// The honesty half: even with a confirmation on record, the routing does
	// not change - only the note does.
	//
	// The key is Gate's own comparison fold (normPath: separators unified,
	// lowercased, trailing sep trimmed), NOT the raw canonical string - that
	// asymmetry is deliberate in blacklist.go (overrideApplies: an allow must
	// rest on the one form the OS vouches for, in the form the gate compares
	// in). Whoever populates this map in production (ticket 21) has to use the
	// same fold, which is exactly why the map is injected and never guessed.
	logs = nil
	note2, w2, a2 := run(map[string]bool{strings.ToLower(canon): true})
	if len(note2.AlreadyUnlocked) != 1 {
		t.Errorf("AlreadyUnlocked = %v with one confirmation on file, want the path listed "+
			"(this is what risk.Gate's bOverrides argument exists for)", note2.AlreadyUnlocked)
	}
	if w2+a2 != 1 {
		t.Errorf("an override record changed the routing (%d/%d): only ticket 21's own flow may "+
			"decide that a B-tier file no longer needs to be asked about", w2, a2)
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func t90Canonical(t *testing.T, p string) string {
	t.Helper()
	c, err := risk.Resolve(p, nil)
	if err != nil {
		t.Fatalf("risk.Resolve(%s): %v", p, err)
	}
	return c.Canonical
}
