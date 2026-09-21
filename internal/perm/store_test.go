package perm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/risk"
)

// Ticket 90, part 1: the switch itself (R20/M1 三档, M4 只有全自动吃一次 L2,
// 审计三档全写) and the screening semantics every mode must obey.

type memConfig struct {
	mu   sync.Mutex
	mode string
	save int
	fail error
}

func (m *memConfig) Config() *config.Config {
	m.mu.Lock()
	defer m.mu.Unlock()
	c := config.NewDefaults()
	c.Risk.PermissionMode = m.mode
	return c
}

func (m *memConfig) SetPermissionMode(mode risk.Mode) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.fail != nil {
		return m.fail
	}
	m.mode = mode.String()
	m.save++
	return nil
}

func (m *memConfig) saved() (string, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.mode, m.save
}

// recorder is the test's stand-in for the "[audit]" sink the composition root
// hands every security component (cmd/wisp's auditf).
type recorder struct {
	mu    sync.Mutex
	lines []string
}

func (r *recorder) logf(format string, args ...any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lines = append(r.lines, fmt.Sprintf(format, args...))
}

func (r *recorder) all() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.lines...)
}

func (r *recorder) count(needle string) int {
	var n int
	for _, l := range r.all() {
		if strings.Contains(l, needle) {
			n++
		}
	}
	return n
}

type confirmSpy struct {
	mu     sync.Mutex
	calls  []Switch
	answer error
}

func (c *confirmSpy) Confirm(_ context.Context, sw Switch) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.calls = append(c.calls, sw)
	return c.answer
}

func (c *confirmSpy) n() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.calls)
}

func newTestStore(t *testing.T, mode string, cf ConfirmFunc) (*Store, *memConfig, *recorder) {
	t.Helper()
	mc := &memConfig{mode: mode}
	rec := &recorder{}
	s, err := New(Options{Manager: mc, Confirm: cf, Logf: rec.logf, Now: func() time.Time {
		return time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)
	}})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s, mc, rec
}

// TestTicket90ThreeModesExistAndNoFourth is R20/M1: exactly three档, and
// "一键关掉一切" is not reachable through the switch.
func TestTicket90ThreeModesAndNoFourth(t *testing.T) {
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, nil)
	if err := s.Set(context.Background(), risk.Mode(9), "test", "nobody"); err == nil {
		t.Fatal("an undefined fourth mode was accepted; R20/M1 defines exactly three")
	}
	mode, saves := mc.saved()
	if mode != risk.ModeAskEveryStepName || saves != 0 {
		t.Errorf("after a rejected switch the persisted state moved: mode=%q saves=%d", mode, saves)
	}
	if rec.count("invalid-mode") != 1 {
		t.Errorf("the refused fourth-mode switch was not audited: %v", rec.all())
	}
}

// TestTicket90OnlyAutoCostsAConfirmAndAllThreeAreAudited is R20/M4, both halves:
// the confirmation is asymmetric (only the switch INTO auto_approve), the audit
// line is not (every档 change writes one, naming from/to/time/origin/actor).
func TestTicket90OnlyAutoCostsAConfirmAndAllThreeAreAudited(t *testing.T) {
	spy := &confirmSpy{}
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, spy.Confirm)
	ctx := context.Background()

	if err := s.Set(ctx, risk.ModeAskHighRisk, "panel", "owner"); err != nil {
		t.Fatalf("Set(ask_high_risk): %v", err)
	}
	if spy.n() != 0 {
		t.Errorf("switching to ask_high_risk cost %d confirmations, want 0 (R20/M4)", spy.n())
	}
	if err := s.Set(ctx, risk.ModeAutoApprove, "panel", "owner"); err != nil {
		t.Fatalf("Set(auto_approve): %v", err)
	}
	if spy.n() != 1 {
		t.Errorf("switching to auto_approve cost %d confirmations, want 1", spy.n())
	}
	if spy.calls[0].To != risk.ModeAutoApprove || spy.calls[0].From != risk.ModeAskHighRisk {
		t.Errorf("the confirmation was asked about %v->%v, want ask_high_risk->auto_approve",
			spy.calls[0].From, spy.calls[0].To)
	}
	if err := s.Set(ctx, risk.ModeAskEveryStep, "ball", "owner"); err != nil {
		t.Fatalf("Set(ask_every_step): %v", err)
	}
	if spy.n() != 1 {
		t.Errorf("tightening back cost an extra confirmation (%d), want the same 1", spy.n())
	}

	if mode, saves := mc.saved(); mode != risk.ModeAskEveryStepName || saves != 3 {
		t.Errorf("persisted=%q saves=%d, want ask_every_step/3", mode, saves)
	}
	// Three real switches, three audit lines, each carrying the five fields M4
	// asks for: from, to, time, origin, actor.
	switches := rec.count("MODE-SWITCH")
	if switches != 3 {
		t.Fatalf("MODE-SWITCH audit lines = %d, want 3 (one per档 change): %v", switches, rec.all())
	}
	for _, l := range rec.all() {
		if !strings.Contains(l, "from=") || !strings.Contains(l, "to=") ||
			!strings.Contains(l, "at=") || !strings.Contains(l, "origin=") ||
			!strings.Contains(l, "actor=") || !strings.Contains(l, "result=") {
			t.Errorf("audit line is missing a mandatory field: %q", l)
		}
	}
	if rec.count("origin=\"panel\"") != 2 || rec.count("origin=\"ball\"") != 1 {
		t.Errorf("origin not recorded per switch: %v", rec.all())
	}
}

// TestTicket90NoConfirmChannelMeansNoAuto is the fail-closed half of M4: a
// composition that never wired an L2 route cannot reach auto_approve at all.
// "We could not ask" is a refusal, never an assumed yes.
func TestTicket90NoConfirmChannelMeansNoAuto(t *testing.T) {
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, nil)
	err := s.Set(context.Background(), risk.ModeAutoApprove, "panel", "owner")
	if err == nil {
		t.Fatal("auto_approve was reached with no confirmation channel wired")
	}
	if mode, saves := mc.saved(); mode != risk.ModeAskEveryStepName || saves != 0 {
		t.Errorf("state moved anyway: mode=%q saves=%d", mode, saves)
	}
	if rec.count(string(ResultRefused)) != 1 {
		t.Errorf("refused switch was not audited as %s: %v", ResultRefused, rec.all())
	}
	if s.PermissionMode() != risk.ModeAskEveryStep {
		t.Errorf("mode = %v after a refused switch, want ask_every_step", s.PermissionMode())
	}
}

// TestTicket90ConfirmRefusalLeavesModeAlone covers "the operator said no".
func TestTicket90ConfirmRefusalLeavesModeAlone(t *testing.T) {
	spy := &confirmSpy{answer: errors.New("用户拒绝")}
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, spy.Confirm)
	if err := s.Set(context.Background(), risk.ModeAutoApprove, "panel", "owner"); err == nil {
		t.Fatal("Set succeeded although the confirmation was refused")
	}
	if _, saves := mc.saved(); saves != 0 {
		t.Errorf("a refused confirmation still persisted (%d writes)", saves)
	}
	if rec.count(string(ResultRefused)) != 1 {
		t.Errorf("refusal not audited: %v", rec.all())
	}
}

// TestTicket90PersistFailureKeepsMemory is the anti-split-state half of M3: if
// the write to config.toml fails, the runtime must not run a档 the next start
// will not read back.
func TestTicket90PersistFailureKeepsMemory(t *testing.T) {
	mc := &memConfig{mode: risk.ModeAskEveryStepName, fail: errors.New("disk gone")}
	rec := &recorder{}
	s, err := New(Options{Manager: mc, Logf: rec.logf})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := s.Set(context.Background(), risk.ModeAskHighRisk, "cli", "owner"); err == nil {
		t.Fatal("Set reported success while persistence failed")
	}
	if got := s.PermissionMode(); got != risk.ModeAskEveryStep {
		t.Errorf("mode = %v after a failed write, want the previous %v", got, risk.ModeAskEveryStep)
	}
	if rec.count(string(ResultPersistError)) != 1 {
		t.Errorf("persist failure not audited: %v", rec.all())
	}
}

// TestTicket90SnapshotIsReadOnly documents the interface ticket 92 renders: the
// 档 and the switch history, with no way back into the store.
func TestTicket90SnapshotIsReadOnly(t *testing.T) {
	s, _, _ := newTestStore(t, risk.ModeAskEveryStepName, nil)
	snap := s.Snapshot()
	if snap.Mode != risk.ModeAskEveryStep {
		t.Errorf("snapshot mode = %v", snap.Mode)
	}
	last, ok := s.Last()
	if !ok || last.Origin != "startup" {
		t.Fatalf("a fresh store must still report what档 it started on: %+v ok=%v", last, ok)
	}
	// Mutating the returned slice must not touch the store's record.
	snap.History[0].Result = ResultApplied
	again, _ := s.Last()
	if again.Result == ResultApplied {
		t.Error("Snapshot().History aliases the store's own records")
	}
}

// TestTicket90ScreenTable is the semantics AC#2 and AC#4 hang on: for every
// mode, the table below is what the enforcement chain must do. The three
// red-line rows are the ones a mutation of risk.redLine has to turn red.
func TestTicket90ScreenTable(t *testing.T) {
	l2 := func(rules ...risk.RuleID) risk.Decision {
		return risk.Decision{Level: risk.L2, RulesHit: rules}
	}
	cases := []struct {
		name    string
		d       risk.Decision
		askAll  risk.Level
		high    risk.Level
		auto    risk.Level
		highSil bool
		autoSil bool
	}{
		{
			"L1 reversible write",
			risk.Decision{Level: risk.L1, RulesHit: []risk.RuleID{risk.R1}},
			risk.L1, risk.L0, risk.L0, true, true,
		},
		{"R8 不可逆", l2(risk.R8), risk.L2, risk.L2, risk.L2, false, false},
		{"R1 declared L2", l2(risk.R1), risk.L2, risk.L2, risk.L2, false, false},
		{"R2 工作区外", l2(risk.R2), risk.L2, risk.L2, risk.L2, false, false},
		{"R3 敏感路径", l2(risk.R3), risk.L2, risk.L2, risk.L2, false, false},
		{"R4 污染", l2(risk.R4), risk.L2, risk.L2, risk.L2, false, false},
		{"R5 外网", l2(risk.R5), risk.L2, risk.L2, risk.L2, false, false},
		{"R9 fail-closed", l2(risk.R9), risk.L2, risk.L2, risk.L2, false, false},
		{"R6 shell 元字符", l2(risk.R6), risk.L2, risk.L2, risk.L0, false, true},
		{"R7 批量规模", l2(risk.R7), risk.L2, risk.L2, risk.L0, false, true},
	}
	for _, tc := range cases {
		for _, c := range []struct {
			m       risk.Mode
			want    risk.Level
			silence bool
		}{
			{risk.ModeAskEveryStep, tc.askAll, false},
			{risk.ModeAskHighRisk, tc.high, tc.highSil},
			{risk.ModeAutoApprove, tc.auto, tc.autoSil},
		} {
			got := c.m.Screen(tc.d)
			if got.Level != c.want {
				t.Errorf("%s under %s: level = %v, want %v", tc.name, c.m, got.Level, c.want)
			}
			if got.Silenced != c.silence {
				t.Errorf("%s under %s: silenced = %v, want %v", tc.name, c.m, got.Silenced, c.silence)
			}
			if c.m == risk.ModeAskHighRisk && tc.d.Level == risk.L2 && got.Kept == "" {
				t.Errorf("%s: ask_high_risk must be able to say why it kept the question", tc.name)
			}
		}
	}
}

// TestTicket90TaintFlagAndDenyAreNeverSilenced is AC#2's red lines 1 and 2 in
// pure-function form, so a mutation of redLine() has two independent witnesses.
func TestTicket90TaintFlagAndDenyAreNeverSilenced(t *testing.T) {
	tainted := risk.Decision{
		Level: risk.L2, RulesHit: []risk.RuleID{risk.R6}, SessionOverrideBlocked: true,
	}
	for _, m := range []risk.Mode{risk.ModeAskEveryStep, risk.ModeAskHighRisk, risk.ModeAutoApprove} {
		if got := m.Screen(tainted); got.Level != risk.L2 || got.Silenced {
			t.Errorf("mode=%s: a tainted verdict came out %v (silenced=%v); PLAN.md:1640 says no "+
				"authorization - including a permission mode - covers a C25 escalation",
				m, got.Level, got.Silenced)
		}
	}
	deny := risk.Decision{Level: risk.Deny, RulesHit: []risk.RuleID{risk.R3}}
	for _, m := range []risk.Mode{risk.ModeAskEveryStep, risk.ModeAskHighRisk, risk.ModeAutoApprove} {
		got := m.Screen(deny)
		if got.Level != risk.Deny {
			t.Errorf("mode=%s: Deny came out as %v", m, got.Level)
		}
		if got.Silenced {
			t.Errorf("mode=%s: a Deny was marked silenced", m)
		}
	}
	// An impossible level fails closed rather than falling through the switch.
	odd := risk.Decision{Level: risk.Level(42), RulesHit: []risk.RuleID{risk.R7}}
	if got := risk.ModeAutoApprove.Screen(odd); got.Level != risk.Level(42) || got.Silenced {
		t.Errorf("unknown level screened to %v (silenced=%v); want untouched and asking",
			got.Level, got.Silenced)
	}
	// And an out-of-range MODE reads as the strictest one, never as silence.
	bad := risk.Mode(7)
	if got := bad.Screen(risk.Decision{Level: risk.L1}); got.Level != risk.L1 || got.Silenced {
		t.Errorf("undefined mode %d silenced an L1: %v", int(bad), got)
	}
}

// TestTicket90ParseModeIsLoud pins the ticket-83 rule on the read side: an
// unparseable档 is an error that names the vocabulary, and the default is the
// strictest档.
func TestTicket90ParseModeIsLoud(t *testing.T) {
	if _, err := risk.ParseMode("yolo"); err == nil {
		t.Error("ParseMode(yolo) returned no error")
	} else if !strings.Contains(err.Error(), "ask_every_step|ask_high_risk|auto_approve") {
		t.Errorf("error %q must name the whole vocabulary", err)
	}
	if got, err := risk.ParseMode(""); err != nil || got != risk.ModeAskEveryStep {
		t.Errorf("ParseMode(\"\") = %v, %v; want ask_every_step (R20/M2), nil", got, err)
	}
	if risk.DefaultMode() != risk.ModeAskEveryStep {
		t.Error("DefaultMode() is not the strictest档; R20/M2 requires it")
	}
	for _, name := range risk.ModeNames() {
		m, err := risk.ParseMode(name)
		if err != nil || m.String() != name {
			t.Errorf("round trip %q -> %v (%v)", name, m, err)
		}
	}
}
