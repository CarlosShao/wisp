//go:build windows

package main

// Ticket 337 residual, leg 337r2 (cmd/wisp side).
//
// What ticket 337-r1 landed (internal/agent/approval/gate.go:180-182) arms the
// risk package's OWN security sink - risk.Logf - from Options.Logf inside
// approval.New. That covers every assembly that BUILDS its gate through that
// constructor. gate.go:175-179 states the residual it cannot cover, verbatim:
//
//	"a host that injects its own gate (run.go's s.gate != nil branch) never
//	 calls this constructor, so risk.Logf stays unwired there"
//
// risk.Gate is the sole producer of the honest fail-closed sentence (the
// aRuleReason text "... [fail-closed, NOT a match on this rule] ...",
// internal/risk/blacklist.go:346), and it emits it ONLY through risk.Logf, as
// the `rule=` field of "risk: A-list DENY path=... rule=..." (blacklist.go:83).
// Unwired means dropped (blacklist.go:53-57), which is the "消费者 0 枚" the
// ticket was filed over - on the injected path, since 337-r1.
//
// These two cases measure the fix at the injection point (run.go, the
// s.gate != nil branch): the fail-closed reason a real risk.Gate run produces
// lands in THIS runtime's audit sink (agentRuntime.auditf = "[audit] ..." on
// the run's stderr plus the same record through the run's log sink), and it
// lands there whether or not the injecting host had armed risk.Logf with its
// own line first.
//
// Nothing here mocks the producer or the sink: the gate is the resident leg's
// real approval.Gate (newResidentApproval), the assembly is the production
// assembleRuntime, the paths are real files under a relocated USERPROFILE, and
// the verdict lines are risk.Gate's own - the same shape
// internal/agent/approval/ticket337_reason_reader_test.go reads on the console
// leg. AGENTS.md §1.3's injection surface for the CLI is exactly this seam.
//
// REVERSE (AC#2): delete the `risk.Logf = rt.auditf` line in run.go's injected
// branch and case 1 goes red at "the injected assembly did not re-arm" while
// case 2 goes red at "no audit line reached the runtime sink" - the reason is
// invisible on the injected path again. The mutant run is booked in
// docs/evidence/s1/337r2-failclosed-sink.md.

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// safeErr is a stderr with a mutex: assembleRuntime hands this writer to the
// pipeline, and the config-reload watchdog may book a line on it while the case
// reads the tail. A plain bytes.Buffer here is a data race under -race, not a
// smaller test.
type safeErr struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (w *safeErr) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Write(p)
}

func (w *safeErr) Len() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Len()
}

func (w *safeErr) Tail(from int) string {
	w.mu.Lock()
	defer w.mu.Unlock()
	s := w.b.String()
	if from > len(s) {
		return ""
	}
	return s[from:]
}

func (w *safeErr) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.String()
}

// hostile337r2 is the same control-character name the 337-v1 instrument uses:
// the OS refuses to parse it, so the candidate's own expansion stops short and
// the anchor it shares ancestry with is not installed -> uncertainAnchorMiss
// fires and aRuleReason says so out loud. It is never created on disk.
const hostile337r2 = "id_ed25519\u001b\u202e"

// world337r2 is the fail-closed world: an installed ~/.ssh leaf (a guaranteed
// A deny whose rule= is the bare literal), ~/.aws deliberately absent (the
// anchor the over-strict clause names), and a work dir under the profile so an
// unprovable candidate shares ancestry with it.
type world337r2 struct {
	profile string
	work    string
}

func buildWorld337r2(t *testing.T) world337r2 {
	t.Helper()
	base := t.TempDir()
	profile := filepath.Join(base, "profile")
	work := filepath.Join(profile, "work")
	for _, d := range []string{profile, work} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatalf("setup mkdir %s: %v", d, err)
		}
	}
	// Everything EXCEPT ~/.aws, so the first A anchor an inside-profile hostile
	// candidate cannot be proven against is the credentials file.
	for _, leaf := range []string{".git-credentials", ".ssh/id_installed", `AppData/Roaming/wisp/config.toml`} {
		full := filepath.Join(profile, filepath.FromSlash(leaf))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("setup mkdir parent of %s: %v", full, err)
		}
		if err := os.WriteFile(full, []byte("installed leaf\n"), 0o600); err != nil {
			t.Fatalf("setup write %s: %v", full, err)
		}
	}
	t.Setenv("USERPROFILE", profile)
	t.Setenv("HOME", profile)
	t.Setenv("APPDATA", filepath.Join(profile, "AppData", "Roaming"))
	t.Setenv("LOCALAPPDATA", filepath.Join(profile, "AppData", "Local"))
	return world337r2{profile: profile, work: work}
}

// TestTicket337R2InjectedGateArmsRiskSinkOnRuntimeAudit is AC#1's strong half:
// an injected host whose gate ALREADY armed risk.Logf with its own line does not
// keep that wiring - the assembly root re-arms the security sink with the
// runtime's auditf, and a real risk.Gate run's fail-closed sentence lands in
// this run's stderr audit trail.
func TestTicket337R2InjectedGateArmsRiskSinkOnRuntimeAudit(t *testing.T) {
	old := risk.Logf
	risk.Logf = nil
	t.Cleanup(func() { risk.Logf = old })

	w := buildWorld337r2(t)
	f := newRunFixture(t, "openai-chat")
	ra := newResidentApproval()

	// A sentinel standing for "the host armed the sink with something else".
	// run.go cannot read Options.Logf back off the gate (gate.go keeps it
	// write-only), so the honest shape is: whatever a foreign sink was holding,
	// the audit line ends up in the runtime sink after the assembly runs.
	var foreign bytes.Buffer
	risk.Logf = func(format string, args ...any) {
		fmt.Fprintf(&foreign, format+"\n", args...)
	}

	errs := &safeErr{}
	ar, code := assembleRuntime(runSpec{
		stdout:  f.out,
		stderr:  errs,
		dataDir: f.dir,
		gate:    ra.gate,
		ui:      ra.ui,
		cards:   ra.cards,
		answer:  ra.ui.hasBall,
		taskCtx: ra.root,
		notify:  swallowPoster,
	})
	if code != 0 {
		t.Fatalf("assemble with the resident gate injected: exit %d, stderr:\n%s", code, errs.String())
	}
	defer ar.close()

	if risk.Logf == nil {
		t.Fatal("the injected assembly left risk.Logf unwired: the fail-closed reason has no reader on this path (ticket 337 residual)")
	}
	mark := errs.Len()
	sshLeaf := filepath.Join(w.profile, ".ssh", "id_installed")
	if d := risk.Gate(sshLeaf, nil); d.Class != risk.ClassA {
		t.Fatalf("fixture guard: installed ~/.ssh leaf must be A, got class=%s reason=%q", d.Class, d.Reason)
	}
	line := errs.Tail(mark)
	if !strings.Contains(line, "[audit] risk: A-list DENY") || !strings.Contains(line, "~/.ssh/**") {
		t.Fatalf("an A deny on the injected-host assembly did not reach the runtime audit sink; got:\n%q\n"+
			"(that is ticket 337's residual: risk.Gate produced the sentence and nobody was listening)", line)
	}
	if foreign.Len() != 0 {
		t.Errorf("the pre-injection sink still received a risk line (%q): the assembly root did not re-arm risk.Logf with rt.auditf", foreign.String())
	}

	// The sentence ticket 334 bought, on the same injected assembly.
	cand := filepath.Join(w.work, hostile337r2)
	d := risk.Gate(cand, nil)
	if d.Class != risk.ClassA {
		t.Fatalf("fixture guard: the fail-closed shape must still be A here, got class=%s reason=%q", d.Class, d.Reason)
	}
	if !strings.Contains(d.Reason, "fail-closed") || !strings.Contains(d.Reason, "NOT a match") {
		t.Fatalf("fixture guard: not the ticket-334 fail-closed sentence: %q", d.Reason)
	}
	mark = errs.Len()
	risk.Gate(cand, nil)
	line = errs.Tail(mark)
	for _, want := range []string{"risk: A-list DENY", "fail-closed", "NOT a match", "~/.aws/credentials"} {
		if !strings.Contains(line, want) {
			t.Errorf("the runtime audit sink must carry %q so the fail-closed reason is visible to an auditor; got:\n%s", want, line)
		}
	}
	t.Logf("ticket337r2 injected-host audit line: %s", strings.TrimSpace(line))
}

// TestTicket337R2InjectedGateArmsRiskSinkWhenHostNeverArmed covers the plain
// residual: a host that injects a gate without ever arming risk.Logf. Before
// this leg the sink stayed nil on this path and every DENY line dropped.
func TestTicket337R2InjectedGateArmsRiskSinkWhenHostNeverArmed(t *testing.T) {
	old := risk.Logf
	risk.Logf = nil
	t.Cleanup(func() { risk.Logf = old })

	w := buildWorld337r2(t)
	f := newRunFixture(t, "openai-chat")
	ra := newResidentApproval()
	// Deterministic base: the resident leg's own approval.New arms the sink with
	// ra.residentAuditf, so clear it and let the assembly prove its own wiring.
	risk.Logf = nil

	errs := &safeErr{}
	ar, code := assembleRuntime(runSpec{
		stdout:  f.out,
		stderr:  errs,
		dataDir: f.dir,
		gate:    ra.gate,
		ui:      ra.ui,
		cards:   ra.cards,
		answer:  ra.ui.hasBall,
		taskCtx: ra.root,
		notify:  swallowPoster,
	})
	if code != 0 {
		t.Fatalf("assemble with the resident gate injected: exit %d, stderr:\n%s", code, errs.String())
	}
	defer ar.close()

	if risk.Logf == nil {
		t.Fatal("the injected assembly did not arm risk.Logf: the fail-closed reason is still dropped on this path (ticket 337 residual)")
	}
	mark := errs.Len()
	sshLeaf := filepath.Join(w.profile, ".ssh", "id_installed")
	risk.Gate(sshLeaf, nil)
	line := errs.Tail(mark)
	if !strings.Contains(line, "[audit] risk: A-list DENY") {
		t.Fatalf("with an unwired host the DENY line must still land in the runtime audit sink; got:\n%q", line)
	}
}
