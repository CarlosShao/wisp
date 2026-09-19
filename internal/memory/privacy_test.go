package memory

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// TestPrivacyOpsAllDomains covers the SPEC-02 §5 matrix: for profile, memory,
// task_log, tool_call and artifacts — list / delete-one / purge-all /
// export-JSON all work and the export carries full records.
func TestPrivacyOpsAllDomains(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	// Seed every domain.
	if _, err := s.UpsertProfile(ctx, "pref.language", "zh-CN", ProfileManual); err != nil {
		t.Fatal(err)
	}
	memID, err := s.AddMemory(ctx, "喜欢简洁回复", "偏好 简洁")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.StartTaskLog(ctx, TaskLog{ID: "task-p1", State: "running", QueryText: "打开设置"}); err != nil {
		t.Fatal(err)
	}
	if err := s.FinishTaskLog(ctx, "task-p1", "succeeded", "done", 10, 5, 100, ""); err != nil {
		t.Fatal(err)
	}
	if err := s.StartTaskLog(ctx, TaskLog{ID: "task-p2", State: "running", QueryText: "第二条"}); err != nil {
		t.Fatal(err)
	}
	tcID, err := s.InsertToolCall(ctx, ToolCall{
		TaskID: "task-p1", Seq: 1, Tool: "app.launch", ArgsJSON: "{}",
		RiskLevel: RiskL1, CorrelationID: "corr-p1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(s.ArtifactsDir(), "tool-output-1.txt"), []byte("artifact payload"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(s.ArtifactsDir(), "tool-output-2.txt"), []byte("second"), 0o644); err != nil {
		t.Fatal(err)
	}

	// --- list: every domain lists seeded items with IDs.
	lists := map[PrivacyDomain][]PrivacyItem{}
	for _, d := range PrivacyDomains() {
		items, err := s.ListPrivacy(ctx, d)
		if err != nil {
			t.Fatalf("list %s: %v", d, err)
		}
		if len(items) == 0 {
			t.Fatalf("list %s: no items after seeding", d)
		}
		lists[d] = items
	}
	if lists[PrivacyProfile][0].Label != "pref.language" {
		t.Errorf("profile item label = %q", lists[PrivacyProfile][0].Label)
	}
	if lists[PrivacyTaskLog][0].ID != "task-p2" { // newest first
		t.Errorf("task_log item id = %q, want task-p2 (newest first)", lists[PrivacyTaskLog][0].ID)
	}

	// --- delete-one: by row id, task UUID, and artifact file name.
	if err := s.DeletePrivacyItem(ctx, PrivacyMemory, itoa64(memID)); err != nil {
		t.Fatalf("delete memory item: %v", err)
	}
	if items, _ := s.ListPrivacy(ctx, PrivacyMemory); len(items) != 0 {
		t.Errorf("memory still has %d items after delete-one", len(items))
	}
	if err := s.DeletePrivacyItem(ctx, PrivacyTaskLog, "task-p2"); err != nil {
		t.Fatalf("delete task item: %v", err)
	}
	if err := s.DeletePrivacyItem(ctx, PrivacyToolCall, itoa64(tcID)); err != nil {
		t.Fatalf("delete tool_call item: %v", err)
	}
	if err := s.DeletePrivacyItem(ctx, PrivacyArtifacts, "tool-output-1.txt"); err != nil {
		t.Fatalf("delete artifact item: %v", err)
	}
	// Missing id surfaces ErrNotFound through every domain.
	if err := s.DeletePrivacyItem(ctx, PrivacyProfile, "99999"); err == nil {
		t.Error("delete missing profile item accepted")
	}

	// --- export-JSON: full records, valid envelope, per domain.
	for _, d := range PrivacyDomains() {
		data, err := s.ExportPrivacy(ctx, d)
		if err != nil {
			t.Fatalf("export %s: %v", d, err)
		}
		var env privacyExport
		if err := json.Unmarshal(data, &env); err != nil {
			t.Fatalf("export %s not valid JSON: %v\n%s", d, err, data)
		}
		if env.Domain != string(d) {
			t.Errorf("export domain = %q, want %q", env.Domain, d)
		}
		if env.ExportedAt == "" || env.Records == nil {
			t.Errorf("export %s envelope incomplete: %+v", d, env)
		}
		raw, _ := json.Marshal(env.Records)
		switch d {
		case PrivacyProfile:
			if !strings.Contains(string(raw), "pref.language") {
				t.Errorf("profile export lost records: %s", raw)
			}
		case PrivacyArtifacts:
			if !strings.Contains(string(raw), "tool-output-2.txt") ||
				!strings.Contains(string(raw), "size_bytes") {
				t.Errorf("artifacts export lacks manifest details: %s", raw)
			}
		case PrivacyTaskLog:
			if !strings.Contains(string(raw), "task-p1") {
				t.Errorf("task_log export lost records: %s", raw)
			}
		}
	}

	// --- purge-all: countable, and empties the domain.
	n, err := s.PurgePrivacy(ctx, PrivacyProfile)
	if err != nil || n != 1 {
		t.Fatalf("purge profile = %d, %v; want 1", n, err)
	}
	n, err = s.PurgePrivacy(ctx, PrivacyTaskLog)
	if err != nil || n != 1 {
		t.Fatalf("purge task_log = %d, %v; want 1", n, err)
	}
	n, err = s.PurgePrivacy(ctx, PrivacyArtifacts)
	if err != nil || n != 1 {
		t.Fatalf("purge artifacts = %d, %v; want 1", n, err)
	}
	for _, d := range PrivacyDomains() {
		items, err := s.ListPrivacy(ctx, d)
		if err != nil {
			t.Fatalf("list %s after purge: %v", d, err)
		}
		if len(items) != 0 {
			t.Errorf("domain %s still has %d items after purge", d, len(items))
		}
	}

	// Unknown domain rejected everywhere with the valid list.
	if _, err := s.ListPrivacy(ctx, "secrets"); err == nil ||
		!strings.Contains(err.Error(), "unknown privacy domain") {
		t.Errorf("unknown domain accepted: %v", err)
	}
	if err := s.DeletePrivacyItem(ctx, "secrets", "x"); err == nil {
		t.Error("unknown domain delete accepted")
	}
	if _, err := s.PurgePrivacy(ctx, "secrets"); err == nil {
		t.Error("unknown domain purge accepted")
	}
	if _, err := s.ExportPrivacy(ctx, "secrets"); err == nil {
		t.Error("unknown domain export accepted")
	}
}

func itoa64(v int64) string { return strconv.FormatInt(v, 10) }
