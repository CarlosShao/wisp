package perm

// Ticket 372-r2: the READ-SIDE downgrade.
//
// The revert of ticket 367's direction rule (store_367_test.go) makes
// ask_every_step -> ask_high_risk reachable again on every entry, including ones
// with no confirmation channel. What keeps that harmless is this half: a Store
// with Confirm == nil answers ask_every_step from PermissionMode(), so the bridge
// (which reads the mode once per call through tools.ModeSource) asks about
// everything whatever config.toml says. config.toml itself is never rewritten -
// that is the difference between "downgrade the effective read" and "change the
// user's setting", and the same shape Step-Code's effectiveState uses
// (docs/evidence/s1/m4-survey.md:106).
//
// Teeth on the pre-372 blob: TestPermissionModeDowngradesWithoutAConfirmChannel
// reads auto_approve / ask_high_risk there and fails, and
// TestPermissionModeDowngradeIsRecordedOnce finds no MODE-DOWNGRADE line at all.

import (
	"context"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// TestPermissionModeDowngradesWithoutAConfirmChannel is the effective-read table:
// with no channel every configured value reads ask_every_step, the configured
// value is untouched, and attaching a channel restores the configured reading.
func TestPermissionModeDowngradesWithoutAConfirmChannel(t *testing.T) {
	for _, configured := range []risk.Mode{
		risk.ModeAskEveryStep, risk.ModeAskHighRisk, risk.ModeAutoApprove,
	} {
		s, mc, _ := newTestStore(t, configured.String(), nil)
		if got := s.PermissionMode(); got != risk.ModeAskEveryStep {
			t.Errorf("configured %s with NO channel reads %s, want the downgrade to ask_every_step (372-r2)",
				configured, got)
		}
		// The downgrade must not become a config write: that would be a different
		// feature (and a silent change of the operator's setting).
		if mode, saves := mc.saved(); mode != configured.String() || saves != 0 {
			t.Errorf("configured %s: the downgrade rewrote the stored value (mode=%q saves=%d, want %q/0)",
				configured, mode, saves, configured)
		}

		// Same configured value, a channel attached: no downgrade, because the
		// process can ask for the one card that makes the looser档 legitimate.
		withLeg, _, _ := newTestStore(t, configured.String(), (&confirmSpy{}).Confirm)
		if got := withLeg.PermissionMode(); got != configured {
			t.Errorf("configured %s WITH a channel reads %s, want the configured %s",
				configured, got, configured)
		}
	}
}

// TestPermissionModeDowngradeIsRecordedOnce pins the audit half: a security
// posture change that only lives in a getter is not audited. One MODE-DOWNGRADE
// line at construction names both ends and its own result value; with a channel
// there is no downgrade and therefore no line.
func TestPermissionModeDowngradeIsRecordedOnce(t *testing.T) {
	_, _, rec := newTestStore(t, risk.ModeAutoApproveName, nil)
	if got := rec.count("MODE-DOWNGRADE"); got != 1 {
		t.Fatalf("MODE-DOWNGRADE lines = %d, want exactly 1: %v", got, rec.all())
	}
	var line string
	for _, l := range rec.all() {
		if strings.Contains(l, "MODE-DOWNGRADE") {
			line = l
		}
	}
	for _, needle := range []string{
		"from=auto_approve",
		"to=ask_every_step",
		"origin=\"startup\"",
		"result=" + string(ResultDowngraded),
	} {
		if !strings.Contains(line, needle) {
			t.Errorf("the downgrade line does not carry %q: %s", needle, line)
		}
	}
	if strings.Contains(line, "MODE-SWITCH") {
		t.Errorf("the downgrade calls itself a switch, but config.toml did not move: %s", line)
	}

	withLeg, _, rec2 := newTestStore(t, risk.ModeAutoApproveName, (&confirmSpy{}).Confirm)
	if got := rec2.count("MODE-DOWNGRADE"); got != 0 {
		t.Errorf("a Store with a channel recorded %d downgrade line(s), want 0: %v", got, rec2.all())
	}
	if got := withLeg.PermissionMode(); got != risk.ModeAutoApprove {
		t.Errorf("a Store with a channel reads %s, want auto_approve", got)
	}
}

// TestSetAccountsOnTheConfiguredValueNotTheDowngradedRead is the accounting tooth
// for the read-side change. Configured auto_approve + no channel reads
// ask_every_step; the operator then asks for ask_every_step. If Set compared
// against the EFFECTIVE read it would call that a no-change and persist nothing,
// leaving config.toml on the looser value forever. It must compare against the
// configured value, apply, and record the persisted pair.
func TestSetAccountsOnTheConfiguredValueNotTheDowngradedRead(t *testing.T) {
	s, mc, _ := newTestStore(t, risk.ModeAutoApproveName, nil)
	if got := s.PermissionMode(); got != risk.ModeAskEveryStep {
		t.Fatalf("premise moved: the legless read is %s, want ask_every_step", got)
	}
	if err := s.Set(context.Background(), risk.ModeAskEveryStep, "panel", "owner"); err != nil {
		t.Fatalf("the request to make the file stricter was refused: %v", err)
	}
	last, ok := s.Last()
	if !ok {
		t.Fatal("the switch left no record")
	}
	if last.From != risk.ModeAutoApprove || last.To != risk.ModeAskEveryStep {
		t.Errorf("audit pair = %s -> %s, want the PERSISTED pair auto_approve -> ask_every_step", last.From, last.To)
	}
	if last.Result != ResultApplied {
		t.Errorf("result = %q, want %q: the stricter value must land on disk", last.Result, ResultApplied)
	}
	if mode, saves := mc.saved(); mode != risk.ModeAskEveryStepName || saves != 1 {
		t.Errorf("stored mode=%q saves=%d, want ask_every_step/1", mode, saves)
	}
	if got := s.PermissionMode(); got != risk.ModeAskEveryStep {
		t.Errorf("effective read = %s after the tightening, want ask_every_step", got)
	}
}

// TestAutoApproveWithNoChannelIsRefusedNotPanicked is the panic shot as a case.
// Set's auto_approve branch calls s.confirm with no nil test of its own; the only
// thing that makes that safe is the target gate refusing to reach it when the
// channel is nil. A mutation that drops the nil half of the gate makes this test
// panic (a failure, not a pass), which is the property ticket 372 must not lose.
func TestAutoApproveWithNoChannelIsRefusedNotPanicked(t *testing.T) {
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("Set panicked on auto_approve with no channel: %v (the target gate must refuse first)", rec)
		}
	}()
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, nil)
	err := s.Set(context.Background(), risk.ModeAutoApprove, "panel", "owner")
	if err == nil {
		t.Fatal("auto_approve was reached with no confirmation channel")
	}
	if last, _ := s.Last(); last.Result != ResultRefused {
		t.Errorf("audit result = %q, want %q", last.Result, ResultRefused)
	}
	if mode, saves := mc.saved(); mode != risk.ModeAskEveryStepName || saves != 0 {
		t.Errorf("stored mode=%q saves=%d, want ask_every_step/0", mode, saves)
	}
	if got := rec.count(string(ResultRefused)); got != 1 {
		t.Errorf("refusal lines = %d, want 1: %v", got, rec.all())
	}
	// The edge the panic shot's start mode matters for: the same refusal from the
	// middle档 (no config fallback involved).
	if err := s.Set(context.Background(), risk.ModeAutoApprove, "panel", "owner"); err == nil {
		t.Fatal("auto_approve was reached on a second attempt")
	}
}
