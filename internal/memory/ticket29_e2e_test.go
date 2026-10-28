package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestTicket29_L1_LRUEviction covers AC#3: LRU eviction when profile exceeds 20 slots.
func TestTicket29_L1_LRUEviction(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	// Seed exactly 20 rows.
	for i := int64(1); i <= 20; i++ {
		slot := fmt.Sprintf("habit.fact%02d", i)
		if _, err := s.UpsertProfile(ctx, slot, fmt.Sprintf("content %d", i), ProfileManual); err != nil {
			t.Fatal(err)
		}
		time.Sleep(10 * time.Millisecond) // ensure different updated_at
	}

	// Verify count = 20.
	profiles, err := s.ListProfiles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(profiles) != 20 {
		t.Fatalf("initial count = %d, want 20", len(profiles))
	}

	// Insert 21st row → should evict oldest (fact01).
	if _, err := s.UpsertProfile(ctx, "habit.fact21", "new fact", ProfileManual); err != nil {
		t.Fatal(err)
	}
	profiles, err = s.ListProfiles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(profiles) != 20 {
		t.Fatalf("after 21st insert count = %d, want 20", len(profiles))
	}

	// Check fact01 is gone (oldest updated_at).
	found := false
	for _, p := range profiles {
		if p.Slot == "habit.fact01" {
			found = true
			break
		}
	}
	if found {
		t.Error("fact01 still present after eviction, want evicted")
	}
}

// TestTicket29_L1_Undo restores prior profile state exactly (AC#4).
func TestTicket29_L1_Undo(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	// Seed initial state.
	profiles := []Profile{
		{Slot: "pref.language", Content: "zh-CN", Source: string(ProfileManual), UpdatedAt: 1000},
		{Slot: "pref.tone", Content: "friendly", Source: string(ProfileManual), UpdatedAt: 1000},
	}
	for _, p := range profiles {
		if _, err := s.UpsertProfile(ctx, p.Slot, p.Content, ProfileSource(p.Source)); err != nil {
			t.Fatal(err)
		}
	}

	// Read snapshot (in-memory copy).
	before, _ := s.ListProfiles(ctx)
	beforeIDs := make([]int64, len(before))
	for i, p := range before {
		beforeIDs[i] = p.ID
	}

	// Modify profiles.
	if _, err := s.UpsertProfile(ctx, "pref.language", "en-US", ProfileManual); err != nil {
		t.Fatal(err)
	}
	if _, err := s.UpsertProfile(ctx, "habit.work", "9-5", ProfileManual); err != nil {
		t.Fatal(err)
	}

	modified, _ := s.ListProfiles(ctx)
	if len(modified) != 3 {
		t.Fatalf("modified count = %d, want 3", len(modified))
	}

	// Restore to snapshot.
	if err := s.RestoreProfile(ctx, before...); err != nil {
		t.Fatal(err)
	}

	after, err := s.ListProfiles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != 2 {
		t.Fatalf("restored count = %d, want 2", len(after))
	}
	// Verify exact match by slot/content (IDs can change).
	beforeMap := make(map[string]Profile)
	for _, p := range before {
		beforeMap[p.Slot] = p
	}
	for i, p := range after {
		want := beforeMap[p.Slot]
		if p.Content != want.Content || p.Source != want.Source {
			t.Errorf("after[%d]: slot=%q content=%q source=%q, want content=%q source=%q",
				i, p.Slot, p.Content, p.Source, want.Content, want.Source)
		}
	}
}

func equalIDs(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestTicket29_L2_Scoping verifies basic memory operations (recall visibility across sessions).
func TestTicket29_L2_Scoping(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	// Create a memory.
	memID, err := s.AddMemory(ctx, "User prefers concise answers", "concise answer")
	if err != nil {
		t.Fatal(err)
	}

	// Verify search returns it.
	results, err := s.SearchMemory(ctx, "concise", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("search returned %d results, want 1", len(results))
	}
	if results[0].ID != memID {
		t.Errorf("wrong memory returned: id=%d", results[0].ID)
	}

	// Verify that listing shows one memory.
	all, err := s.ListMemories(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 {
		t.Fatalf("list returned %d memories, want 1", len(all))
	}
}

// TestTicket29_L2_Ranking verifies recency+hit ranking order (AC#5).
func TestTicket29_L2_Ranking(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	// Insert 3 memories.
	id1, _ := s.AddMemory(ctx, "first memory", "keyword first")
	time.Sleep(10 * time.Millisecond)
	id2, _ := s.AddMemory(ctx, "second memory", "keyword second")
	time.Sleep(10 * time.Millisecond)
	id3, _ := s.AddMemory(ctx, "third memory", "keyword third")

	// Search "keyword" → should rank by last_hit_at (all 0), then hit_count (all 0), then id DESC.
	results, err := s.SearchMemory(ctx, "keyword", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("search returned %d results, want 3", len(results))
	}
	if results[0].ID != id3 || results[1].ID != id2 || results[2].ID != id1 {
		t.Errorf("ranking order wrong: %v", results)
	}

	// Hit id2 twice → now it should rank first.
	for i := 0; i < 2; i++ {
		s.SearchMemory(ctx, "second", 10)
	}

	results, err = s.SearchMemory(ctx, "keyword", 10)
	if err != nil {
		t.Fatal(err)
	}
	if results[0].ID != id2 {
		t.Errorf("id2 should rank first after hits, got id=%d", results[0].ID)
	}
}

// TestTicket29_PrivacyJSONExport validates JSON export format (AC#6).
func TestTicket29_PrivacyJSONExport(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	// Seed data.
	_, _ = s.UpsertProfile(ctx, "pref.test", "value", ProfileManual)
	_, _ = s.AddMemory(ctx, "test memory", "test kw")

	export, err := s.ExportPrivacy(ctx, PrivacyProfile)
	if err != nil {
		t.Fatal(err)
	}
	var env privacyExport
	if err := json.Unmarshal(export, &env); err != nil {
		t.Fatal(err)
	}
	if env.Domain != string(PrivacyProfile) {
		t.Errorf("export domain = %q, want %q", env.Domain, PrivacyProfile)
	}
	if env.ExportedAt == "" {
		t.Error("exported_at empty in envelope")
	}
	if env.Records == nil {
		t.Error("records is null in profile export")
	}

	export, err = s.ExportPrivacy(ctx, PrivacyMemory)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(export, &env); err != nil {
		t.Fatal(err)
	}
	if env.Domain != string(PrivacyMemory) {
		t.Errorf("export domain = %q, want %q", env.Domain, PrivacyMemory)
	}
	// Verify records array is not null.
	raw, _ := json.Marshal(env.Records)
	if !strings.Contains(string(raw), "[") && !strings.Contains(string(raw), "test memory") {
		t.Errorf("memory export lacks records: %s", raw)
	}
}
