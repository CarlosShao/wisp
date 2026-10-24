package tools_test

// Ticket 377: the `tools: call` line must state the EFFECTIVE permission档 -
// the one this process screens with - and not config.toml's value.
//
// Why the line and not Decision.Mode: the census (docs/evidence/s1/377-a1-census.md
// section 4) enumerated Decision.Mode's consumers and found zero production
// readers - the approval queue stores it, panel_pump's LiveApprovals projection
// drops it, Gate.Replay has no production caller. Ticket 377's AC#2 forbids
// asserting "a field nobody reads exists", so the recorded value needed a
// surface a human (or a grep) actually reads. book()'s line is that surface: one
// per call, in the same sink as MODE-READ/MODE-SILENCE, and the biggest forensic
// face this repository has.
//
// The difference is REAL, not cosmetic, and these cases are built so that a line
// reading the wrong face goes red with the wrong value rather than merely with a
// missing token:
//
//	panel-inbound shape: config auto_approve + Confirm nil + Set(ask_high_risk)
//	  the file ends at ask_high_risk, the process screens at ask_every_step
//	  (ticket 372's read-side collapse) -> the line must say ask_every_step,
//	  i.e. NEITHER of the two values config.toml ever held.
//	`wisp run` shape: confirm is never nil, so the two faces are equal
//	  -> the line must still say auto_approve (the positive control: the fix
//	  must not turn every booking into the strictest档).

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/perm"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// --- fixtures ---------------------------------------------------------------

// t377Gate counts the two card branches and answers "the user let it through",
// so "one window was raised" is never also an assertion about refusal.
type t377Gate struct {
	mu        sync.Mutex
	windows   int
	approvals int
}

func (g *t377Gate) PendingWindow(context.Context, tools.Decision) (tools.Answer, string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.windows++
	return tools.AnswerAllow, ""
}

func (g *t377Gate) PendingApproval(context.Context, tools.Decision) (tools.Answer, string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.approvals++
	return tools.AnswerAllow, ""
}

func (g *t377Gate) counts() (int, int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.windows, g.approvals
}

// t377Log is the audit sink (production: agentRuntime.auditf), verbatim.
type t377Log struct {
	mu    sync.Mutex
	lines []string
}

func (l *t377Log) logf(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lines = append(l.lines, fmt.Sprintf(format, args...))
}

func (l *t377Log) all() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.lines...)
}

// t377Line returns the LAST line containing marker.
func (l *t377Log) line(marker string) string {
	var hit string
	for _, s := range l.all() {
		if strings.Contains(s, marker) {
			hit = s
		}
	}
	return hit
}

// t377Rig is one assembly of the seam, built out of the REAL objects: a
// config.toml on disk, config.NewManager, perm.New with the confirm leg the
// caller names, the D34 fs roster, the C26 canonicalizer over one allowed root
// and the bridge itself. The only substitution is the audit sink and the gate.
type t377Rig struct {
	b     *tools.Bridge
	store *perm.Store
	mgr   *config.Manager
	gate  *t377Gate
	log   *t377Log
	dir   string
}

// t377New composes the shape. mode is written to config.toml; confirm nil is the
// `wisp panel-inbound` leg (panel_inbound.go's own nil), a function is the
// `wisp run` leg (rt.confirmModeSwitch, never nil).
func t377New(t *testing.T, mode string, confirm perm.ConfirmFunc) *t377Rig {
	t.Helper()
	dir := t.TempDir()
	cfg := fmt.Sprintf("schema_version = 2\n\n[fs]\nallowed_dirs = [%q]\n\n[risk]\npermission_mode = %q\n",
		filepath.ToSlash(dir), mode)
	cfgPath := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o600); err != nil {
		t.Fatal(err)
	}
	mgr, err := config.NewManager(cfgPath, nil)
	if err != nil {
		t.Fatalf("config.NewManager: %v", err)
	}
	rig := &t377Rig{mgr: mgr, gate: &t377Gate{}, log: &t377Log{}, dir: dir}
	store, err := perm.New(perm.Options{Manager: mgr, Confirm: confirm, Logf: rig.log.logf})
	if err != nil {
		t.Fatalf("perm.New: %v", err)
	}
	rig.store = store

	paths := tools.NewPathCanonicalizer([]string{dir}, nil)
	reg := tools.NewRegistry()
	for _, e := range tools.BuiltinFSEntries(tools.FSDeps{Paths: paths}) {
		if err := reg.Register(e); err != nil {
			t.Fatalf("Register(%s): %v", e.Tool.Name(), err)
		}
	}
	rig.b = tools.New(tools.Options{
		Registry: reg,
		Paths:    paths,
		Gate:     rig.gate,
		Modes:    store, // the same object the assembly hands the bridge
		Logf:     rig.log.logf,
	})
	return rig
}

// t377Write is one L1 write inside the allowed root - the level the L1 window
// owns, so the mode can silence it (auto_approve) or leave the question standing
// (ask_every_step).
func (r *t377Rig) write(t *testing.T, name string) agent.ToolOutcome {
	t.Helper()
	args := fmt.Sprintf(`{"path":%q,"content":"377-r1 byte"}`,
		filepath.ToSlash(filepath.Join(r.dir, name)))
	out, err := r.b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "task-377", CorrelationID: "corr-377", CallID: "call-377",
		Name: "fs.write", Args: json.RawMessage(args),
	})
	if err != nil {
		t.Fatalf("a refusal is not a host fault: %v", err)
	}
	return out
}

// t377CallLine is the audit line this ticket extends. Its absence is its own
// failure message: the assertion below must never silently pass on "".
func (r *t377Rig) callLine(t *testing.T) string {
	t.Helper()
	line := r.log.line("tools: call ")
	if line == "" {
		t.Fatalf("no `tools: call` line was booked at all; sink:\n%s",
			strings.Join(r.log.all(), "\n"))
	}
	return line
}

// --- AC#2: the panel-inbound shape -----------------------------------------

// Test377PanelInboundBookingNamesTheEffectiveMode is AC#2's criterion, and it is
// the whole ticket in one call:
//
//	config.toml           auto_approve       (what the operator wrote)
//	Confirm               nil                (this seam has no leg)
//	perm.Set(ask_high_risk)                  (a panel request; the FILE moves)
//	=> config.toml        ask_high_risk      (after)
//	=> effective档        ask_every_step     (before AND after: 372 collapse)
//	=> tools: call must   mode=ask_every_step
//
// A line that printed either config value would say auto_approve or
// ask_high_risk; only a line reading the effective face can say ask_every_step.
// That is why this case distinguishes "the field exists" from "the field is
// right", which AC#2 forbids confusing.
func Test377PanelInboundBookingNamesTheEffectiveMode(t *testing.T) {
	rig := t377New(t, risk.ModeAutoApproveName, nil)

	// The premise: this assembly really is the one where the two faces differ.
	if got := rig.store.PermissionMode(); got != risk.ModeAskEveryStep {
		t.Fatalf("premise broken: effective档 = %s, want %s (372's collapse did not fire)",
			got, risk.ModeAskEveryStep)
	}
	if !strings.Contains(strings.Join(rig.log.all(), "\n"), "MODE-DOWNGRADE") {
		t.Errorf("premise broken: no MODE-DOWNGRADE line was booked by perm.New")
	}

	// The switch the panel would make. The FILE moves; the effective档 cannot.
	if err := rig.store.Set(context.Background(), risk.ModeAskHighRisk, "panel", "t377-ac2"); err != nil {
		t.Fatalf("Set(ask_high_risk): %v", err)
	}
	if got := rig.mgr.Config().PermissionMode(); got != risk.ModeAskHighRisk {
		t.Fatalf("premise broken: config.toml still says %s after Set, want %s", got, risk.ModeAskHighRisk)
	}
	if got := rig.store.PermissionMode(); got != risk.ModeAskEveryStep {
		t.Fatalf("premise broken: effective档 moved to %s; this seam must not", got)
	}

	out := rig.write(t, "panel-inbound.txt")
	if out.IsError {
		t.Fatalf("the L1 write did not execute in the downgraded assembly: %q", out.Text)
	}
	if w, a := rig.gate.counts(); w != 1 || a != 0 {
		t.Errorf("the downgraded assembly raised %d windows / %d approvals, want 1/0: "+
			"the effective档 must still ASK about an L1", w, a)
	}

	line := rig.callLine(t)
	t.Logf("AC#2 reading: %s", line)
	if !strings.Contains(line, "mode="+risk.ModeAskEveryStepName) {
		t.Errorf("the booked line does not name the effective档 %q: %s",
			risk.ModeAskEveryStepName, line)
	}
	// The two wrong answers, named so a mutation that reads either face fails
	// with the value it printed rather than only with "token missing".
	for _, wrong := range []string{risk.ModeAutoApproveName, risk.ModeAskHighRiskName} {
		if strings.Contains(line, "mode="+wrong) {
			t.Errorf("the booked line reports config.toml's %q, not the effective档: %s", wrong, line)
		}
	}
}

// --- the two-shot comparison -------------------------------------------------

// Test377ConfiguredVersusEffectiveTwoShots is the census's section-2 reading
// re-run through the booking line: ONE config.toml (auto_approve) and one L1
// call, and the only variable is whether the process was handed a leg.
//
//	confirm=nil  -> effective ask_every_step, 1 L1 window, line mode=ask_every_step
//	confirm=leg  -> effective auto_approve,   0 windows,    line mode=auto_approve
//
// The leg case is also the positive control the ruling asked for: on the
// `wisp run` seam the configured and effective faces are equal by construction
// (confirm is never nil there), and the line must keep saying auto_approve.
func Test377ConfiguredVersusEffectiveTwoShots(t *testing.T) {
	leg := func(context.Context, perm.Switch) error { return nil }
	legs := []struct {
		name        string
		confirm     perm.ConfirmFunc
		wantMode    string
		wantWindows int
		wantSilence bool
	}{
		{"wisp panel-inbound shape: confirm=nil", nil, risk.ModeAskEveryStepName, 1, false},
		{"wisp run shape: confirm=leg", leg, risk.ModeAutoApproveName, 0, true},
	}
	for _, c := range legs {
		t.Run(c.name, func(t *testing.T) {
			rig := t377New(t, risk.ModeAutoApproveName, c.confirm)
			out := rig.write(t, "two-shot.txt")
			if out.IsError {
				t.Fatalf("the L1 write was refused in this leg: %q", out.Text)
			}
			if w, _ := rig.gate.counts(); w != c.wantWindows {
				t.Errorf("windows = %d, want %d: the mode's effect is the half the line must "+
					"agree with", w, c.wantWindows)
			}
			line := rig.callLine(t)
			t.Logf("%s -> %s", c.name, line)
			if !strings.Contains(line, "mode="+c.wantMode) {
				t.Errorf("mode= is not %q: %s", c.wantMode, line)
			}
			// Cross-check the printed token against the mode's own effect: a
			// line hard-coded to one name would agree with one leg and not the
			// other, and would also disagree with MODE-SILENCE here.
			silenced := rig.log.line("tools: MODE-SILENCE ") != ""
			if silenced != c.wantSilence {
				t.Errorf("MODE-SILENCE present = %v, want %v: the ledger and the booking "+
					"line must tell one story", silenced, c.wantSilence)
			}
			if c.wantSilence && !strings.Contains(rig.log.line("tools: MODE-SILENCE "), "mode="+c.wantMode) {
				t.Errorf("MODE-SILENCE names a different mode than the booking line: %s",
					rig.log.line("tools: MODE-SILENCE "))
			}
		})
	}
}

// Test377UnknownToolBookingStillNamesAMode pins the second half of the hoist: a
// call refused BEFORE C3 ever ran (the unknown-tool branch) has no screening to
// inherit a mode from, and it used to have no Mode at all. The line is emitted
// for every call, so it must not print a fabricated zero on this branch - it
// prints the effective档 the process is in, which here is the collapsed one.
func Test377UnknownToolBookingStillNamesAMode(t *testing.T) {
	rig := t377New(t, risk.ModeAutoApproveName, nil)
	out, err := rig.b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "task-377", CorrelationID: "corr-377-unknown", CallID: "call-377-unknown",
		Name: "no.such.tool", Args: json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("an unknown tool is a refusal, not a host fault: %v", err)
	}
	if !out.IsError {
		t.Fatalf("an unknown tool must be refused: %+v", out)
	}
	line := rig.callLine(t)
	if !strings.Contains(line, "mode="+risk.ModeAskEveryStepName) {
		t.Errorf("the pre-C3 refusal booked no effective档: %s", line)
	}
}
