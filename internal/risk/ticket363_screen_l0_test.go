package risk

import (
	"strings"
	"testing"
)

// Ticket 363: L0 is a level this package DEFINES, not an unidentified one.
//
// Two sentences in this package disagree about `Screen`'s entry test, and the
// disagreement is the whole bug:
//
//	assessor.go:41  L0 < L1 < L2 < Deny: fusion takes the max, and Deny is terminal.
//	assessor.go:45  L0: read-only, no side effects - execute directly.
//	assessor.go:46  L0 Level = iota            (the enum's ZERO value)
//	mode.go:175+188 return Silenced{Level: L0, Silenced: true, Mode: m}
//
// The function already emits L0 at its own exits, and it already treats an L1
// verdict silenced by the mode as "have the caller skip the question". So an L0
// verdict ARRIVING here must not fall into the `d.Level != L1 && d.Level != L2`
// branch, which was written for "the judge returned something impossible" and
// which stamps Kept="无法识别的风险级，fail-closed 保留询问" onto every read-only
// call. That string is not just ugly: Kept's contract (mode.go:143-146) says it
// names the red line that REFUSED TO SILENCE a question, and L0 never had a
// question, so the audit row named a red line that does not exist.
//
// The test is written so that a revert of the L0 branch in Screen turns it red:
// Kept becomes the fail-closed sentence again.
func TestTicket363ScreenRecognizesL0(t *testing.T) {
	// An out-of-range mode is included on purpose: L0 needs no question in ANY
	// mode, so the level has to be recognised before the mode is consulted.
	for _, m := range []Mode{ModeAskEveryStep, ModeAskHighRisk, ModeAutoApprove, Mode(7)} {
		got := m.Screen(Decision{Level: L0})
		if got.Level != L0 {
			t.Errorf("mode=%s: Screen(L0) returned level %v, want L0 untouched (a mode may "+
				"lower the QUESTION, never rewrite the level the judge concluded)", m, got.Level)
		}
		if got.Kept != "" {
			t.Errorf("mode=%s: Screen(L0).Kept = %q, want empty: L0 is the judge's own zero "+
				"value (assessor.go:46) and its 'read-only, no side effects - execute "+
				"directly' verdict (assessor.go:45); no red line kept a question here, so the "+
				"audit text must not name one (mode.go:143-146)", m, got.Kept)
		}
		if got.Silenced {
			// Two sentences describe this field and they pull in opposite
			// directions for L0. mode.go:141 says "true when the caller must NOT
			// ask"; gate.go:43-47 (the field of the same name on tools.Decision)
			// says "true when the mode removed a question this verdict would
			// otherwise have asked". At L0 there WAS no question, so the first
			// reading would make every read-only tool_call row claim a
			// mode-removed question - the same category error in the other
			// direction - and would also put the bridge on its MODE-SILENCE
			// line, whose text says "a silenced question". The "do not ask" step
			// for L0 is taken by the caller's own `case risk.L0`
			// (internal/tools/bridge.go:433-435), not here.
			t.Errorf("mode=%s: Screen(L0).Silenced = true; nothing was asked, so nothing was "+
				"silenced (ticket 363)", m)
		}
	}
}

// TestTicket363UnknownLevelStillKeepsItsQuestion is the control that keeps the
// fix from being a deletion: the fail-closed branch must still fire for a level
// that really is impossible. Narrowing "not L1/L2" to "not L0/L1/L2" is the
// fix; dropping the branch is not.
func TestTicket363UnknownLevelStillKeepsItsQuestion(t *testing.T) {
	impossible := Level(42)
	for _, m := range []Mode{ModeAskEveryStep, ModeAskHighRisk, ModeAutoApprove} {
		got := m.Screen(Decision{Level: impossible})
		if got.Level != impossible {
			t.Errorf("mode=%s: an impossible level %v came back as %v; it must travel "+
				"unchanged so the caller's default branch can refuse it", m, impossible, got.Level)
		}
		if got.Silenced {
			t.Errorf("mode=%s: an impossible level %v was marked silenced", m, impossible)
		}
		if !strings.Contains(got.Kept, "无法识别") {
			t.Errorf("mode=%s: an impossible level %v came back with Kept=%q, want the "+
				"fail-closed sentence naming the unrecognised level", m, impossible, got.Kept)
		}
	}
}

// TestTicket363FourLevelsKeepTheirRoute pins the no-regression table at the
// unit level, one row per level the routing owns (the behavioural half - how
// many cards the Gate interface is asked for - lives in internal/tools, where
// the real bridge is composed):
//
//	L0    never asks (execute directly)
//	L1    asks under ask_every_step, silenced by the two relaxed modes
//	L2    asks everywhere a red line keeps it
//	Deny  is refused and never treated as a question
func TestTicket363FourLevelsKeepTheirRoute(t *testing.T) {
	cases := []struct {
		name     string
		d        Decision
		m        Mode
		level    Level
		silenced bool
		kept     bool
	}{
		{"L0/ask_every_step", Decision{Level: L0}, ModeAskEveryStep, L0, false, false},
		{"L0/ask_high_risk", Decision{Level: L0}, ModeAskHighRisk, L0, false, false},
		{"L0/auto_approve", Decision{Level: L0}, ModeAutoApprove, L0, false, false},
		{"L1/ask_every_step", Decision{Level: L1}, ModeAskEveryStep, L1, false, false},
		{"L1/ask_high_risk", Decision{Level: L1}, ModeAskHighRisk, L0, true, false},
		{"L1/auto_approve", Decision{Level: L1}, ModeAutoApprove, L0, true, false},
		{"L2 ask_high_risk", Decision{Level: L2}, ModeAskHighRisk, L2, false, true},
		{"L2 R8 auto_approve", Decision{Level: L2, RulesHit: []RuleID{R8}}, ModeAutoApprove, L2, false, true},
		{"L2 R6 auto_approve (silenceable)", Decision{Level: L2, RulesHit: []RuleID{R6}}, ModeAutoApprove, L0, true, false},
		{"Deny ask_every_step", Decision{Level: Deny, RulesHit: []RuleID{R3}}, ModeAskEveryStep, Deny, false, true},
		{"Deny auto_approve", Decision{Level: Deny, RulesHit: []RuleID{R3}}, ModeAutoApprove, Deny, false, true},
	}
	for _, tc := range cases {
		got := tc.m.Screen(tc.d)
		if got.Level != tc.level {
			t.Errorf("%s: effective level = %v, want %v", tc.name, got.Level, tc.level)
		}
		if got.Silenced != tc.silenced {
			t.Errorf("%s: silenced = %v, want %v", tc.name, got.Silenced, tc.silenced)
		}
		if (got.Kept != "") != tc.kept {
			t.Errorf("%s: Kept = %q, want kept-question = %v", tc.name, got.Kept, tc.kept)
		}
		if got.Mode != tc.m {
			t.Errorf("%s: Mode = %v, want the mode it was screened under (%v)", tc.name, got.Mode, tc.m)
		}
	}
}
