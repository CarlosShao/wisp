package main

// Ticket 377: the档 this process books must be the EFFECTIVE one.
//
// Decision.Mode's own doc comment (internal/tools/gate.go) says the field is
// "the permission mode this call was screened under" and that it exists so
// "why was I not asked?" is answerable "from the record, not from the config
// file's current contents". Two writers in this package did the opposite: they
// read rt.mgr.Config().PermissionMode(), i.e. config.toml's value, through
// agentRuntime.permissionMode (approval_always.go's widening card at
// runWidening, and config_reload.go's D36 re-confirmation card at
// confirmLockedLoosening).
//
// The census measured the consequence (docs/evidence/s1/377-a1-census.md
// section 6): on the `wisp run` seam the two faces are equal by construction,
// because confirm is the confirmModeSwitch method value and is never nil, so
// ticket 372's read-side collapse cannot fire. They come apart on the
// `wisp panel-inbound` seam, which passes nil - the exact assembly this case
// builds. Before this ticket the record said auto_approve about a process that
// screens at ask_every_step: the widening card's Mode would have answered the
// field's own question with the opposite of what happened.
//
// This is a nail on the READ, and it is red on the pre-377 tree by running (no
// API the change introduces): the two legs below differ only in the confirm
// leg, and the effective read must follow that difference while the config read
// cannot.

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/perm"
	"github.com/CarlosShao/wisp/internal/risk"
)

// t377ModeRuntime builds the smallest honest owner: a real config.toml holding
// one mode, the real Manager over it, and the real Store the assembly would hand
// the runtime. Nothing here is a mock - the mode read is the seam under test.
func t377ModeRuntime(t *testing.T, confirm perm.ConfirmFunc) *agentRuntime {
	t.Helper()
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, configFileName)
	body := "schema_version = 2\n\n[risk]\npermission_mode = \"auto_approve\"\n"
	if err := os.WriteFile(cfgPath, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	mgr, err := config.NewManager(cfgPath, nil)
	if err != nil {
		t.Fatalf("config.NewManager: %v", err)
	}
	store, err := perm.New(perm.Options{Manager: mgr, Confirm: confirm})
	if err != nil {
		t.Fatalf("perm.New: %v", err)
	}
	return &agentRuntime{mgr: mgr, modes: store}
}

// Test377BookedModeReadsTheEffectivePermissionMode is the ticket's 〔丁〕 item 2:
// with no confirmation leg the configured auto_approve is not what this process
// screens with, so it is not what the card may book.
func Test377BookedModeReadsTheEffectivePermissionMode(t *testing.T) {
	rt := t377ModeRuntime(t, nil)
	if got := rt.mgr.Config().PermissionMode(); got != risk.ModeAutoApprove {
		t.Fatalf("premise broken: config.toml reads %s, want auto_approve", got)
	}
	if got := rt.permissionMode(); got != risk.ModeAskEveryStep {
		t.Errorf("the widening card would book %s for a process that screens at %s: "+
			"Decision.Mode has to answer \"why was I not asked?\" from the record",
			got, risk.ModeAskEveryStep)
	}
}

// Test377BookedModeKeepsTheConfiguredModeWhenAChannelExists is the positive
// control: the `wisp run` seam always has a leg (confirmModeSwitch is a method
// value, never nil), so there the configured and effective faces are the same
// value and the read must stay auto_approve. A fix that simply hard-coded the
// strictest档 would pass the case above and fail this one.
func Test377BookedModeKeepsTheConfiguredModeWhenAChannelExists(t *testing.T) {
	leg := func(context.Context, perm.Switch) error { return nil }
	rt := t377ModeRuntime(t, leg)
	if got := rt.permissionMode(); got != risk.ModeAutoApprove {
		t.Errorf("with a confirmation leg attached the effective档 is config.toml's "+
			"auto_approve; booked %s instead", got)
	}
}

// Test377BookedModeWithoutAStoreFailsClosed keeps the "unreadable means
// strictest" half of the old function: a runtime with no store has no evidence
// that any channel exists, and this field must never guess toward silence.
func Test377BookedModeWithoutAStoreFailsClosed(t *testing.T) {
	var nilRuntime *agentRuntime
	if got := nilRuntime.permissionMode(); got != risk.ModeAskEveryStep {
		t.Errorf("a nil runtime booked %s, want the strictest档", got)
	}
	bare := &agentRuntime{}
	if got := bare.permissionMode(); got != risk.ModeAskEveryStep {
		t.Errorf("a runtime with no mode store booked %s, want the strictest档", got)
	}
}
