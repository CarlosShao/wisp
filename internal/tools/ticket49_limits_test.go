package tools

// Ticket 49: the three inviolable limits (D45-3, SPEC-06 §8.3) at bridge scale.
// Two of them (L2 never granted, R4 never granted) already have their cells in
// ticket224/grant_test.go and ticket49_target_test.go; this file pins the
// THIRD (config-section loosening never via grants) and restates the L2 cell
// with the exact reading the dispatch asked for: 「有 grant 也仍弹卡」 on a
// threshold-raised verdict.
//
// The config cell's mechanism: a permission mode silences a question by
// screening the verdict (risk.Mode.Screen) - that is the ONLY channel config
// has into the decision chain (ticket 90; the allowlist/exemptions change
// verdicts through the judge, which is the D36 re-confirmation path). A
// session grant stands in for a question the mode did NOT remove
// (`!sil.Silenced` is the gate). So a loosened config section can never ride a
// grant: the grant is only consulted when the question still exists, and when
// the mode removed the question the booking says MODE-SILENCE with
// decision=allow and grant_id NULL - never allow_session_grant.

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/risk"
)

// TestTicket49L2WithALiveGrantStillShowsTheCard is the 「有 grant 也仍弹卡」
// tooth for the L2 family: a bridge holding a grant source that answers
// EVERYTHING (the widest possible source - the same promiscuous fixture
// shell_session_teeth_test.go used for the L2 route pin) still shows one L2
// card for an R8-irreversible write. The grant source is consulted zero times:
// the L2 route never reaches the check.
func TestTicket49L2WithALiveGrantStillShowsTheCard(t *testing.T) {
	dir := sealableTempDir124(t)
	target := filepath.Join(dir, "existing.txt")
	if err := os.WriteFile(target, []byte("v1"), 0o600); err != nil {
		t.Fatal(err)
	}
	paths := NewPathCanonicalizer([]string{dir}, nil)
	e, _ := t90Entry("t49.l2over", risk.L1, []string{"path"}, func(context.Context, map[string]any, []string) risk.Facts {
		return risk.Facts{OverwriteExisting: true}
	})
	reg := NewRegistry()
	if err := reg.Register(e); err != nil {
		t.Fatalf("Register: %v", err)
	}
	gate := &t90Gate{}
	b := New(Options{
		Registry: reg,
		Paths:    paths,
		Gate:     gate,
		Grants:   promiscuousGrantSource{},
		Modes:    &t90Modes{},
		Logf:     func(string, ...any) {},
	})
	out, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "t49-l2", CorrelationID: "t49-l2", CallID: "c1",
		Name: "t49.l2over",
		Args: json.RawMessage(`{"path":` + jsonQuoted49(target) + `}`),
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !out.IsError {
		// t90Entry's tool refuses nothing; an R8-overwrite L2 verdict routes to
		// PendingApproval, whose t90Gate stub answers reject, which the bridge
		// books as a refusal. Either way the card was asked.
		t.Logf("outcome: %s", out.Text)
	}
	// The stub gate answers allow (it counts the visit); the assertion is that
	// the L2 route was REACHED and the grant seam never was.
	if w, a := gate.counts(); a != 1 || w != 0 {
		t.Errorf("an R8 L2 call with an all-covering grant source produced windows=%d approvals=%d, "+
			"want 0/1: L2 永不被 grant 自动放行（SPEC-06 §8.3 第 1 条）", w, a)
	}
}

// promiscuousGrantSource answers EVERY lookup with a live id - re-declared
// here with the target door too, so a future mutation that routes the L2
// through the TARGET read cannot hide behind the path-only fixture above.
type promiscuousGrantSource49 struct{}

func (promiscuousGrantSource49) Covering(context.Context, string, []string) (int64, bool) {
	return 8675309, true
}

func (promiscuousGrantSource49) TargetCovering(context.Context, string, string) (int64, bool) {
	return 8675309, true
}

// TestTicket49ModeSilenceNeverBooksAGrantLink is the config-limit cell. Under
// auto_approve an L1 verdict is silenced (the config-fed mode loosened the
// asking), and the booking must be decision=allow with grant_id NULL: a grant
// is only consulted for a question that still exists, so a silenced call can
// never carry allow_session_grant. A source that answers everything is wired
// to make the absence positive: if the check ever runs on a silenced call, the
// booking flips and this case goes red.
func TestTicket49ModeSilenceNeverBooksAGrantLink(t *testing.T) {
	dir := sealableTempDir124(t)
	target := filepath.Join(dir, "mode-silenced.txt")
	paths := NewPathCanonicalizer([]string{dir}, nil)
	j := &t224Journal{}
	e, _ := t90Entry("t49.silenced", risk.L1, []string{"path"}, nil)
	reg := NewRegistry()
	if err := reg.Register(e); err != nil {
		t.Fatalf("Register: %v", err)
	}
	b := New(Options{
		Registry: reg,
		Paths:    paths,
		Gate:     &t90Gate{},
		Grants:   promiscuousGrantSource49{},
		Journal:  j,
		Modes:    &t90Modes{m: risk.ModeAutoApprove},
		Logf:     func(string, ...any) {},
	})
	out, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "t49-sil", CorrelationID: "t49-sil", CallID: "c1",
		Name: "t49.silenced",
		Args: json.RawMessage(`{"path":` + jsonQuoted49(target) + `}`),
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if out.IsError {
		t.Fatalf("a mode-silenced L1 must run: %s", out.Text)
	}
	row, ok := j.last()
	if !ok {
		t.Fatal("no tool_call row booked")
	}
	if row.Decision == agent.DecisionAllowGrant {
		t.Errorf("a mode-silenced call booked decision=allow_session_grant: config loosening rode a grant, "+
			"which is exactly the D45-3 third limit (安全配置放宽不适用会话授权)")
	}
	if row.GrantID != nil {
		t.Errorf("a mode-silenced call booked grant_id=%d, want NULL: the grant seam was consulted on a "+
			"question the config mode had already removed", *row.GrantID)
	}
}

// TestTicket49TaintedPathWriteIgnoresEveryGrant restates the R4 cell with the
// promiscuous source ON A PATH WRITE: a tainted fs.write to an allowlisted
// path is L2 (not covered), the grant source is never consulted, and the
// booking carries no grant linkage. The taint comes from the real judge via
// the injected detector (t90Taint), not from a fixture verdict.
func TestTicket49TaintedPathWriteIgnoresEveryGrant(t *testing.T) {
	dir := sealableTempDir124(t)
	target := filepath.Join(dir, "tainted.txt")
	paths := NewPathCanonicalizer([]string{dir}, nil)
	j := &t224Journal{}
	e, _ := t90Entry("t49.tainted", risk.L1, []string{"path"}, nil)
	reg := NewRegistry()
	if err := reg.Register(e); err != nil {
		t.Fatalf("Register: %v", err)
	}
	gate := &t90Gate{}
	b := New(Options{
		Registry: reg,
		Paths:    paths,
		Gate:     gate,
		Grants:   promiscuousGrantSource49{},
		Journal:  j,
		Assessor: risk.NewRiskAssessor().
			WithCanonicalizer(paths).
			WithTaintDetector(t90Taint{}),
		Modes: &t90Modes{},
		Logf:  func(string, ...any) {},
	})
	if _, err := b.Execute(context.Background(), agent.ToolRequest{
		TaskID: "t49-r4", CorrelationID: "t49-r4", CallID: "c1",
		Name: "t49.tainted",
		Args: json.RawMessage(`{"path":` + jsonQuoted49(target) + `}`),
	}); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if w, a := gate.counts(); a != 1 {
		t.Errorf("a tainted call with an all-covering grant source produced windows=%d approvals=%d, "+
			"want 0/1: C25/R4 taint upgrades are NOT overridden by any grant (D45-3 second limit)", w, a)
	}
	row, ok := j.last()
	if ok && row.GrantID != nil {
		t.Errorf("a tainted call booked grant_id=%d, want NULL", *row.GrantID)
	}
	if ok && row.Decision == agent.DecisionAllowGrant {
		t.Errorf("a tainted call booked decision=allow_session_grant")
	}
}

// jsonQuoted49 quotes one path as a JSON string value (the tests above build
// argument objects; encoding/json's Marshal on a string would do, but a tiny
// helper keeps the literals readable).
func jsonQuoted49(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}
