package memory

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// Ticket 40 AC#6 teeth: a privacy row carrying risk_level=="L0" and
// decision=="allow" must NOT read better than reality — the human detail
// must say it was auto-allowed by tier and nobody was asked, not a bare
// "allow"; only L0 may say 自动放行 (an L2+allow row has a human behind it);
// and the exported bundle's human-readable field must carry the SAME
// sentence the page renders. The recorded column values themselves stay
// verbatim (ticket 184 AC#5(a): 取证面不许改字) — rewording them reddens
// these tests.
//
// Mutations live OUTSIDE the repo (each tooth is proven by mutating its own
// subject in a throwaway copy, never by editing tracked files here).

// ac6ToolCallRow books one tool_call row with the given level and decision.
func ac6ToolCallRow(t *testing.T, s *Store, taskID string, seq int64, level, decision string) int64 {
	t.Helper()
	id, err := s.InsertToolCall(ctx40(), ToolCall{
		TaskID: taskID, Seq: seq, Tool: "fs.read", ArgsJSON: "{}",
		RiskLevel: level, Decision: decision, CorrelationID: "corr-ac6",
	})
	if err != nil {
		t.Fatalf("insert tool_call %s/%s: %v", level, decision, err)
	}
	return id
}

func ctx40() context.Context { return context.Background() }

// detailSegments pins the derived detail's exact shape: the recorded column
// values verbatim as the first segments (取证面 echoed, not reworded), and —
// for an L0+allow row only — the honest-by-tier sentence as the last one.
func detailSegments(t *testing.T, detail string, want []string) {
	t.Helper()
	got := strings.Split(detail, " · ")
	if len(got) != len(want) {
		t.Fatalf("detail %q has %d segments, want %d (%v)", detail, len(got), len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("detail %q segment %d = %q, want %q", detail, i, got[i], want[i])
		}
	}
}

// Tooth ①: an L0+allow row rendered as a bare "allow" is red — the human
// detail must carry the honest-by-tier sentence.
func TestPrivacyL0AutoAllowDetailSaysNobodyWasAsked(t *testing.T) {
	s := openTestStore(t)
	ctx := ctx40()

	id := ac6ToolCallRow(t, s, "task-ac6", 1, RiskL0, "allow")

	items, err := s.ListPrivacy(ctx, PrivacyToolCall)
	if err != nil {
		t.Fatal(err)
	}
	var row *PrivacyItem
	for i := range items {
		if items[i].ID == itoa64(id) {
			row = &items[i]
		}
	}
	if row == nil {
		t.Fatalf("L0+allow row id %d missing from the privacy list", id)
	}
	// The bare recorded decision is still echoed verbatim…
	if !strings.Contains(row.Detail, "allow") {
		t.Errorf("detail lost the recorded decision word: %q", row.Detail)
	}
	// …but a bare "allow" alone is exactly what this AC forbids.
	if strings.TrimSpace(row.Detail) == "allow" || strings.HasSuffix(row.Detail, " · allow") {
		t.Errorf("L0+allow row rendered as a bare allow (read better than reality): %q", row.Detail)
	}
	for _, want := range []string{"按档自动放行（L0）", "未询问任何人"} {
		if !strings.Contains(row.Detail, want) {
			t.Errorf("L0+allow detail %q lacks the honest-by-tier sentence %q", row.Detail, want)
		}
	}
	// Exact shape: recorded columns verbatim, then the honest sentence —
	// no rewording, no inserted gloss, in this exact order.
	detailSegments(t, row.Detail, []string{"L0", "fs.read", "allow", "按档自动放行（L0），未询问任何人"})
}

// Tooth ② (reverse): only L0 may say 自动放行. An L2+allow row has a human
// behind it — writing "自动放行" there reads better than reality and is red.
func TestPrivacyL2AllowNeverSaysAutoAllow(t *testing.T) {
	s := openTestStore(t)
	ctx := ctx40()

	id := ac6ToolCallRow(t, s, "task-ac6", 2, RiskL2, "allow")

	items, err := s.ListPrivacy(ctx, PrivacyToolCall)
	if err != nil {
		t.Fatal(err)
	}
	var row *PrivacyItem
	for i := range items {
		if items[i].ID == itoa64(id) {
			row = &items[i]
		}
	}
	if row == nil {
		t.Fatalf("L2+allow row id %d missing from the privacy list", id)
	}
	if strings.Contains(row.Detail, "自动放行") || strings.Contains(row.Detail, l0AutoAllowNotice) {
		t.Errorf("L2+allow row wrongly claimed a tier auto-allow: %q", row.Detail)
	}
	if !strings.Contains(row.Detail, "allow") {
		t.Errorf("L2+allow detail lost the recorded decision word: %q", row.Detail)
	}
	// Exact shape: the recorded columns verbatim and NOTHING after the
	// decision — an L2 allow row's human story is not "automatic".
	detailSegments(t, row.Detail, []string{"L2", "fs.read", "allow"})
}

// Tooth ③: the exported bundle's human-readable field (decision_notes,
// keyed by row id) must carry the SAME sentence the page renders — two
// wordings across the two faces is red. The records themselves stay the
// verbatim column values: this pin doubles as the 取证面不许改字 tooth
// (ticket 184 AC#5(a)) — rewording decision/risk_level/tool in the export
// reddens it.
func TestPrivacyExportDecisionNotesMatchListDetail(t *testing.T) {
	s := openTestStore(t)
	ctx := ctx40()

	l0 := ac6ToolCallRow(t, s, "task-ac6", 3, RiskL0, "allow")
	l2 := ac6ToolCallRow(t, s, "task-ac6", 4, RiskL2, "allow")

	items, err := s.ListPrivacy(ctx, PrivacyToolCall)
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]PrivacyItem{}
	for _, it := range items {
		byID[it.ID] = it
	}

	data, err := s.ExportPrivacy(ctx, PrivacyToolCall)
	if err != nil {
		t.Fatal(err)
	}
	var env struct {
		Records       []json.RawMessage `json:"records"`
		DecisionNotes map[string]string `json:"decision_notes"`
	}
	if err := json.Unmarshal(data, &env); err != nil {
		t.Fatalf("export not valid JSON: %v\n%s", err, data)
	}
	if env.DecisionNotes == nil {
		t.Fatalf("export envelope lacks decision_notes: %s", data)
	}

	// The page face's strings…
	wantL0 := byID[itoa64(l0)].Detail
	wantL2 := byID[itoa64(l2)].Detail
	if !strings.Contains(wantL0, l0AutoAllowNotice) || strings.Contains(wantL2, l0AutoAllowNotice) {
		t.Fatalf("page face itself is wrong — L0 detail %q, L2 detail %q", wantL0, wantL2)
	}
	// …must be byte-identical to the export face's strings.
	gotL0, ok := env.DecisionNotes[itoa64(l0)]
	if !ok {
		t.Fatalf("export decision_notes missing row %d: %v", l0, env.DecisionNotes)
	}
	if gotL0 != wantL0 {
		t.Errorf("export/page drift on L0 row:\n export: %q\n page:   %q", gotL0, wantL0)
	}
	gotL2, ok := env.DecisionNotes[itoa64(l2)]
	if !ok {
		t.Fatalf("export decision_notes missing row %d: %v", l2, env.DecisionNotes)
	}
	if gotL2 != wantL2 {
		t.Errorf("export/page drift on L2 row:\n export: %q\n page:   %q", gotL2, wantL2)
	}

	// The verbatim columns (ticket 184 AC#5(a)): each exported record echoes
	// the recorded column values untouched — the rewording tooth. ToolCall
	// has no json tags, so the record keys are the Go field names.
	for i, raw := range env.Records {
		var rec struct {
			ID        int64  `json:"ID"`
			RiskLevel string `json:"RiskLevel"`
			Decision  string `json:"Decision"`
			Tool      string `json:"Tool"`
		}
		if err := json.Unmarshal(raw, &rec); err != nil {
			t.Fatalf("record %d not valid JSON: %v\n%s", i, err, raw)
		}
		if rec.RiskLevel != RiskL0 && rec.RiskLevel != RiskL2 {
			t.Errorf("record %d risk_level reworded: %q", rec.ID, rec.RiskLevel)
		}
		if rec.Decision != "allow" {
			t.Errorf("record %d decision reworded: %q", rec.ID, rec.Decision)
		}
		if rec.Tool != "fs.read" {
			t.Errorf("record %d tool reworded: %q", rec.ID, rec.Tool)
		}
	}
}
