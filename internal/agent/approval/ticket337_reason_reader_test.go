//go:build windows

package approval

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// Ticket 337 form 甲 (local audit log): the honest fail-closed sentence ticket
// 334 bought had ZERO readers - it never reached a bubble, a panel, the model
// text, or the log (census: docs/evidence/s1/337r1-*.md, ruler R1/R2 = 0). The
// fix arms risk.Logf (the package's own security sink) with the host's audit
// sink, in approval.New. risk.Gate is the only producer of that sentence and it
// emits it through risk.Logf as the `rule=` field of its "risk: A-list DENY ..."
// line, so arming the sink is what makes the reason reach the operator/audit
// log the moment readBlacklist runs risk.Gate.
//
// This instrument proves BOTH halves of that claim on a real Gate:
//   - the wiring: New(opts) turns a nil risk.Logf into the host sink; and
//   - the flow: the fail-closed sentence ("[fail-closed, NOT a match ...]") and
//     a guaranteed A-deny line both land in that sink when Gate runs.
//
// REVERSE: comment out the `if o.Logf != nil { risk.Logf = o.Logf }` block in
// gate.go New and this test goes red at the wiring check - the reason is
// invisible again. (The mutant run is recorded in the evidence file.)
//
// The world is built through PUBLIC API only (os / filepath / t.Setenv), never
// the unexported newWorld334, because this lives in package approval. It mirrors
// that builder: the A-anchor leaves are installed by name or left absent, so the
// over-strict fail-closed shape is decidable on any Windows machine.
const hostile337 = "id_ed25519\u001b\u202e" // control chars: the OS refuses to parse the name -> "present but unproven"

type world337 struct {
	profile string
	work    string
}

func buildWorld337(t *testing.T) world337 {
	t.Helper()
	base := t.TempDir()
	profile := filepath.Join(base, "profile")
	work := filepath.Join(profile, "work")
	for _, d := range []string{profile, work} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatalf("setup mkdir %s: %v", d, err)
		}
	}
	// Install everything EXCEPT the ~/.aws leaves, so the first A anchor that an
	// inside-profile hostile candidate cannot be proven against is ~/.aws.
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
	return world337{profile: profile, work: work}
}

func TestTicket337FailClosedReasonReachesAuditLog(t *testing.T) {
	// Deterministic base: a leaked sink from another test would mask a missing
	// wiring line, so clear it and restore on exit.
	old := risk.Logf
	risk.Logf = nil
	t.Cleanup(func() { risk.Logf = old })

	var lines []string
	capture := func(format string, args ...any) {
		lines = append(lines, fmt.Sprintf(format, args...))
	}

	// --- half 1: the wiring. New(opts) must arm the risk sink with the host sink.
	g := New(Options{Logf: capture})
	if g == nil {
		t.Fatal("New returned nil")
	}
	if risk.Logf == nil {
		t.Fatal("approval.New did not arm risk.Logf: the fail-closed reason has no reader again (ticket 337)")
	}
	risk.Logf("ticket337 wiring probe %s", "ok")
	if !strings.Contains(strings.Join(lines, "\n"), "ticket337 wiring probe ok") {
		t.Fatalf("risk.Logf is non-nil but does not deliver to the host sink; captured:\n%s", strings.Join(lines, "\n"))
	}

	// --- half 2: the flow, on a real Gate run.
	w := buildWorld337(t)

	// (a) A guaranteed A-deny: an installed ~/.ssh leaf. No dependence on the
	// control-character name, so this half cannot be flaky. Its rule= field is
	// the anchor literal, and the DENY line must reach the sink.
	sshLeaf := filepath.Join(w.profile, ".ssh", "id_installed")
	if d := risk.Gate(sshLeaf, nil); d.Class != risk.ClassA {
		t.Fatalf("fixture guard: installed ~/.ssh leaf must be A, got class=%s reason=%q", d.Class, d.Reason)
	}
	lines = nil
	risk.Gate(sshLeaf, nil)
	if joined := strings.Join(lines, "\n"); !strings.Contains(joined, "risk: A-list DENY") || !strings.Contains(joined, "~/.ssh/**") {
		t.Fatalf("a real A deny must reach the audit sink (ticket 337); captured:\n%s", joined)
	}

	// (b) The sentence ticket 334 bought: a hostile candidate inside the profile
	// whose own spelling cannot be proven expanded, with ~/.aws absent, must
	// fail closed AS A and the honest sentence must reach the sink.
	cand := filepath.Join(w.work, hostile337)
	d := risk.Gate(cand, nil)
	if d.Class != risk.ClassA {
		t.Fatalf("fixture guard: the fail-closed shape must still be A here, got class=%s reason=%q", d.Class, d.Reason)
	}
	if !strings.Contains(d.Reason, "fail-closed") || !strings.Contains(d.Reason, "NOT a match") {
		t.Fatalf("fixture guard: not the ticket-334 fail-closed sentence: %q", d.Reason)
	}
	lines = nil
	risk.Gate(cand, nil)
	joined := strings.Join(lines, "\n")
	for _, want := range []string{"risk: A-list DENY", "fail-closed", "NOT a match", "~/.aws/credentials"} {
		if !strings.Contains(joined, want) {
			t.Errorf("the audit log must carry %q so the fail-closed reason is visible to a human; captured:\n%s", want, joined)
		}
	}
	t.Logf("ticket337 audit line: %s", strings.TrimSpace(joined))
}
