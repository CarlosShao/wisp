// Package config_test holds the ticket-267 checks that have to see BOTH sides
// of the contract they guard: the [risk].confirm_timeout_sec band in this
// layer, and the C18 warning lead it is derived from, which lives in
// internal/agent/approval.
//
// It is the EXTERNAL test package deliberately. internal/agent/approval
// imports internal/agent, which imports internal/config, so a test file in
// `package config` importing approval would close an import cycle. The external
// package may import approval precisely because config itself does not - and
// that asymmetry is the invariant this file pins: the config layer never
// depends on the thing it is protecting, it just refuses a value the consumer
// cannot honour.
package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/config"
)

// warningLeadSec is the C18 pre-timeout warning lead, read off the symbol it
// is defined by (never copied as a literal - a copy is what let this drift to
// "no bounds at all" in the first place).
func warningLeadSec(t *testing.T) int {
	t.Helper()
	lead := approval.DefaultApprovalWarning
	if lead <= 0 || lead%time.Second != 0 {
		t.Fatalf("approval.DefaultApprovalWarning = %v, want a whole positive number of seconds", lead)
	}
	return int(lead / time.Second)
}

// loadWithTimeout writes a minimal config carrying only [risk].confirm_timeout_sec
// and runs it through the real load path (LoadFile = strict decode + semantic
// validation), so what is being asserted is refusal-to-load, not a value the
// program later ignores.
func loadWithTimeout(t *testing.T, sec int) (*config.Config, error) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	body := fmt.Sprintf("schema_version = 2\n\n[risk]\nconfirm_timeout_sec = %d\n", sec)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	c, _, err := config.LoadFile(path, nil)
	return c, err
}

// TestConfirmTimeoutFloorSitsAboveApprovalWarningLead is ticket 267's
// same-source guard. The floor of the config band is not a number this package
// gets to invent: internal/agent/approval/gate.go arms the C18 warning only
// while Timeout()-WarningLead() > 0, and the queue falls back to
// approval.DefaultApprovalWarning whenever no producer passed a lead. So a
// timeout equal to that lead silently drops the warning, and the band must
// start one second above it. This test derives that boundary from the symbol,
// so moving DefaultApprovalWarning without moving the band turns this red.
func TestConfirmTimeoutFloorSitsAboveApprovalWarningLead(t *testing.T) {
	lead := warningLeadSec(t)

	// At the lead: the config must be refused, not loaded-and-degraded.
	if _, err := loadWithTimeout(t, lead); err == nil {
		t.Fatalf("confirm_timeout_sec = %d must be refused: it equals approval.DefaultApprovalWarning (%v), so gate.go's `lead := Timeout() - WarningLead()` is <= 0 and the C18 warning never arms", lead, approval.DefaultApprovalWarning)
	} else if !strings.Contains(err.Error(), "risk.confirm_timeout_sec") {
		t.Fatalf("error %q must name risk.confirm_timeout_sec", err)
	}

	// One above the lead: the smallest value that keeps the warning alive, and
	// it must load clean.
	c, err := loadWithTimeout(t, lead+1)
	if err != nil {
		t.Fatalf("confirm_timeout_sec = %d (lead+1) must load: %v", lead+1, err)
	}
	if c.Risk.ConfirmTimeoutSec != lead+1 {
		t.Fatalf("loaded confirm_timeout_sec = %d, want %d", c.Risk.ConfirmTimeoutSec, lead+1)
	}
}

// TestConfirmTimeoutIsBandedOnTheLoadPath is AC#1 read through the loader
// rather than through validate() alone: the seeds ticket 267 names (30 / 10 /
// 3601 / 99999) must fail to load at all, and the message must say which key
// and which range.
func TestConfirmTimeoutIsBandedOnTheLoadPath(t *testing.T) {
	lead := warningLeadSec(t)
	for _, bad := range []int{lead, 10, 3601, 99999, 0, -1} {
		_, err := loadWithTimeout(t, bad)
		if err == nil {
			t.Fatalf("confirm_timeout_sec = %d must be refused at load", bad)
		}
		if !strings.Contains(err.Error(), "risk.confirm_timeout_sec") || !strings.Contains(err.Error(), "out of range") {
			t.Fatalf("error %q must name risk.confirm_timeout_sec and state the range it broke", err)
		}
	}
	for _, good := range []int{lead + 1, 300, 3600} {
		if _, err := loadWithTimeout(t, good); err != nil {
			t.Fatalf("confirm_timeout_sec = %d must load: %v", good, err)
		}
	}
}

// TestMinimumLegalConfirmTimeoutKeepsC18Warning is AC#2, and it is a real
// reading rather than a stand-in: it loads the smallest legal config, builds
// the approval queue the way the assembly root does
// (cmd/wisp/run.go:616 and cmd/wisp/resident_approval_windows.go:460 both pass
// time.Duration(ConfirmTimeoutSec)*time.Second and no warning lead), and then
// evaluates the very expression internal/agent/approval/gate.go:528 branches
// on. If that expression ever stops being positive, the C18 warning is gone
// with no error and no log line.
func TestMinimumLegalConfirmTimeoutKeepsC18Warning(t *testing.T) {
	lead := warningLeadSec(t)
	c, err := loadWithTimeout(t, lead+1)
	if err != nil {
		t.Fatalf("the smallest legal confirm_timeout_sec must load: %v", err)
	}

	q := approval.NewQueue(time.Duration(c.Risk.ConfirmTimeoutSec)*time.Second, 0, 0, nil)
	gateLead := q.Timeout() - q.WarningLead()
	if gateLead <= 0 {
		t.Fatalf("C18 warning would not arm: Queue.Timeout()=%v - WarningLead()=%v = %v at the smallest legal confirm_timeout_sec=%d",
			q.Timeout(), q.WarningLead(), gateLead, c.Risk.ConfirmTimeoutSec)
	}
	t.Logf("smallest legal timeout %v, warning lead %v, arming lead %v > 0", q.Timeout(), q.WarningLead(), gateLead)

	// The other half of the same branch: the ceiling must not turn the card
	// into an all-day wait either. 3600s is the band top, and it is still a
	// timeout the user can outlive by answering, not by walking away.
	if q.Timeout() > time.Duration(3600)*time.Second {
		t.Fatalf("queue timeout %v exceeds the configured band ceiling", q.Timeout())
	}
}
