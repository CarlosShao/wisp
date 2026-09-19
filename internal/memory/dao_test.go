package memory

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// bufLogger returns a logger plus the buffer it writes to (asserting log
// lines, e.g. the D35 profile-eviction visibility rule).
func bufLogger(t *testing.T) (*slog.Logger, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	h := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	return slog.New(h), &buf
}

func TestProfileUpsertAndLRUEvictionLogs(t *testing.T) {
	logger, buf := bufLogger(t)
	s := openTestStore(t, WithLogger(logger))
	ctx := context.Background()

	// Fill to the cap.
	for i := 0; i < L1MaxProfiles; i++ {
		if _, err := s.UpsertProfile(ctx, fmt.Sprintf("pref.k%02d", i), fmt.Sprintf("value %d", i), ProfileExtracted); err != nil {
			t.Fatalf("upsert %d: %v", i, err)
		}
	}
	// Distinguish updated_at values so LRU order is deterministic even at
	// second granularity: age row 0 the most by rewriting others.
	for i := 1; i < L1MaxProfiles; i++ {
		bumpProfileUpdatedAt(t, s, fmt.Sprintf("pref.k%02d", i), int64(1000+i))
	}
	bumpProfileUpdatedAt(t, s, "pref.k00", 1) // oldest -> eviction candidate
	buf.Reset()

	// One more row crosses the cap: the oldest (k00) must be evicted, with a
	// log line (D35: eviction must be visible).
	if _, err := s.UpsertProfile(ctx, "pref.k21", "new", ProfileManual); err != nil {
		t.Fatalf("upsert over cap: %v", err)
	}

	profs, err := s.ListProfiles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(profs) != L1MaxProfiles {
		t.Fatalf("profile rows = %d, want %d", len(profs), L1MaxProfiles)
	}
	for _, p := range profs {
		if p.Slot == "pref.k00" {
			t.Error("oldest profile was not evicted")
		}
	}
	if _, err := s.ProfileBySlot(ctx, "pref.k21"); err != nil {
		t.Errorf("newest row missing after eviction: %v", err)
	}

	logText := buf.String()
	if !strings.Contains(logText, "profile LRU eviction") {
		t.Errorf("eviction log line missing; log:\n%s", logText)
	}
	if !strings.Contains(logText, "evicted_slot=pref.k00") {
		t.Errorf("eviction log line does not name the slot; log:\n%s", logText)
	}

	// Same-slot upsert refreshes instead of adding a row.
	if _, err := s.UpsertProfile(ctx, "pref.k21", "updated content", ProfileExtracted); err != nil {
		t.Fatal(err)
	}
	profs, _ = s.ListProfiles(ctx)
	if len(profs) != L1MaxProfiles {
		t.Errorf("upsert duplicated rows: %d", len(profs))
	}
	p, err := s.ProfileBySlot(ctx, "pref.k21")
	if err != nil {
		t.Fatal(err)
	}
	if p.Content != "updated content" || p.Source != string(ProfileExtracted) {
		t.Errorf("upsert did not refresh: %+v", p)
	}

	// Invalid source rejected.
	if _, err := s.UpsertProfile(ctx, "pref.bad", "x", ProfileSource("hacked")); err == nil {
		t.Error("invalid profile source accepted")
	}
}

func bumpProfileUpdatedAt(t *testing.T, s *Store, slot string, ts int64) {
	t.Helper()
	err := s.write(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE profile SET updated_at=? WHERE slot=?`, ts, slot)
		return err
	})
	if err != nil {
		t.Fatalf("bump %s: %v", slot, err)
	}
}

func TestProfileDeleteAndPurge(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	if _, err := s.UpsertProfile(ctx, "pref.a", "a", ProfileManual); err != nil {
		t.Fatal(err)
	}
	p, _ := s.ProfileBySlot(ctx, "pref.a")
	if err := s.DeleteProfile(ctx, p.ID); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteProfile(ctx, p.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("delete missing profile: %v, want ErrNotFound", err)
	}
	for i := 0; i < 3; i++ {
		if _, err := s.UpsertProfile(ctx, fmt.Sprintf("pref.p%d", i), "x", ProfileExtracted); err != nil {
			t.Fatal(err)
		}
	}
	n, err := s.PurgeProfiles(ctx)
	if err != nil || n != 3 {
		t.Fatalf("purge = %d, %v; want 3, nil", n, err)
	}
	profs, _ := s.ListProfiles(ctx)
	if len(profs) != 0 {
		t.Errorf("purge left %d rows", len(profs))
	}
}

func TestMemoryAddSearchDeletePurge(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	id1, err := s.AddMemory(ctx, "用户偏好简短回答", "偏好 简短 回答")
	if err != nil {
		t.Fatal(err)
	}
	id2, err := s.AddMemory(ctx, "工作时间为 9:00-18:00", "HABIT 工作时间")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.AddMemory(ctx, "no keywords", "   "); err == nil {
		t.Error("memory with empty keywords accepted")
	}

	hits, err := s.SearchMemory(ctx, "简短", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 || hits[0].ID != id1 {
		t.Fatalf("search hits = %+v, want [id %d]", hits, id1)
	}

	// Hit stats update (last_hit_at/hit_count) via the writer.
	m, err := s.MemoryByID(ctx, id1)
	if err != nil {
		t.Fatal(err)
	}
	if m.HitCount != 1 || m.LastHitAt == 0 {
		t.Errorf("hit stats not updated: %+v", m)
	}

	// Multi-term OR search + case normalization.
	hits, err = s.SearchMemory(ctx, "HABIT 工作时间", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 || hits[0].ID != id2 {
		t.Fatalf("second search hits = %+v, want [id %d]", hits, id2)
	}

	// Recency-first ordering: hit id2 (older created) -> it now sorts first.
	if _, err := s.SearchMemory(ctx, "偏好 OR 工作", 10); err != nil {
		t.Fatal(err)
	}
	hits, err = s.SearchMemory(ctx, "偏好 工作", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 2 {
		t.Fatalf("want 2 hits, got %+v", hits)
	}
	if hits[0].LastHitAt < hits[1].LastHitAt {
		t.Errorf("not ordered most-recent-hit-first: %+v", hits)
	}

	// LIKE wildcards in the query must not leak (% matches everything).
	hits, _ = s.SearchMemory(ctx, "%", 10)
	if len(hits) != 0 {
		t.Errorf("unescaped LIKE wildcard matched %d rows", len(hits))
	}

	// kind is pinned to 'explicit' by the schema comment.
	rows, _ := s.ListMemories(ctx)
	if len(rows) != 2 {
		t.Fatalf("rows = %d", len(rows))
	}
	for _, r := range rows {
		if r.Kind != MemoryKindExplicit {
			t.Errorf("kind = %q, want explicit", r.Kind)
		}
	}

	if err := s.DeleteMemory(ctx, id1); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteMemory(ctx, id1); !errors.Is(err, ErrNotFound) {
		t.Errorf("delete missing memory: %v", err)
	}
	n, err := s.PurgeMemories(ctx)
	if err != nil || n != 1 {
		t.Fatalf("purge = %d, %v; want 1", n, err)
	}
}

func TestTaskLogLifecycleAndErrorClassValidation(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	tl := TaskLog{ID: "task-1", State: "running", QueryText: "帮我查天气 [已脱敏]"}
	if err := s.StartTaskLog(ctx, tl); err != nil {
		t.Fatal(err)
	}
	// Duplicate id rejected by PK.
	if err := s.StartTaskLog(ctx, tl); err == nil {
		t.Error("duplicate task id accepted")
	}
	// D37 enum enforced via observe.ValidateErrorClass (ticket 03 hook).
	if err := s.FinishTaskLog(ctx, "task-1", "succeeded", "done", 120, 45, 3100, "not_a_class"); err == nil {
		t.Error("invalid error_class accepted")
	}
	if err := s.FinishTaskLog(ctx, "task-1", "succeeded", "done", 120, 45, 3100, string(observe.ClassNetwork)); err != nil {
		t.Fatal(err)
	}
	got, err := s.TaskLogByID(ctx, "task-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.State != "succeeded" || got.EndedAt == nil || got.ErrorClass != string(observe.ClassNetwork) {
		t.Errorf("finished row wrong: %+v", got)
	}
	if got.Currency != "CNY" {
		t.Errorf("currency default = %q, want CNY", got.Currency)
	}
	if got.CostAmountMicro != 3100 {
		t.Errorf("cost micros = %d", got.CostAmountMicro)
	}

	if err := s.StartTaskLog(ctx, TaskLog{ID: "task-2", State: "interrupted", QueryText: "q"}); err != nil {
		t.Fatal(err)
	}
	logs, err := s.ListTaskLogs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(logs) != 2 {
		t.Fatalf("logs = %d", len(logs))
	}
	if err := s.DeleteTaskLog(ctx, "task-1"); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteTaskLog(ctx, "task-1"); !errors.Is(err, ErrNotFound) {
		t.Errorf("delete missing task: %v", err)
	}
	n, err := s.PurgeTaskLogs(ctx)
	if err != nil || n != 1 {
		t.Fatalf("purge = %d, %v; want 1", n, err)
	}
}

func TestToolCallLifecycleAndValidation(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	if err := s.StartTaskLog(ctx, TaskLog{ID: "task-a", State: "running", QueryText: "q"}); err != nil {
		t.Fatal(err)
	}

	tc := ToolCall{
		TaskID: "task-a", Seq: 1, Tool: "fs.read",
		ArgsJSON: `{"path":"C:/x.txt"}`, RiskLevel: RiskL0,
		CorrelationID: "corr-1",
	}
	id, err := s.InsertToolCall(ctx, tc)
	if err != nil {
		t.Fatal(err)
	}

	bad := tc
	bad.RiskLevel = "L9"
	if _, err := s.InsertToolCall(ctx, bad); err == nil {
		t.Error("invalid risk_level accepted")
	}
	bad = tc
	bad.Decision = "maybe"
	if _, err := s.InsertToolCall(ctx, bad); err == nil {
		t.Error("invalid decision accepted")
	}
	bad = tc
	bad.ErrorClass = "not_a_class"
	if _, err := s.InsertToolCall(ctx, bad); err == nil {
		t.Error("invalid error_class accepted")
	}

	grantID := int64(7)
	if err := s.DecideToolCall(ctx, id, "allow_session_grant", &grantID); err != nil {
		t.Fatal(err)
	}
	if err := s.FinishToolCall(ctx, id, "success", ""); err != nil {
		t.Fatal(err)
	}
	if err := s.FinishToolCall(ctx, id, "exploded", ""); err == nil {
		t.Error("invalid outcome accepted")
	}

	got, err := s.ToolCallByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if got.Decision != "allow_session_grant" || got.Outcome != "success" {
		t.Errorf("row state wrong: %+v", got)
	}
	if got.GrantID == nil || *got.GrantID != grantID {
		t.Errorf("grant linkage missing: %+v", got.GrantID)
	}
	if got.DecidedAt == nil || got.EndedAt == nil {
		t.Errorf("timestamps missing: %+v", got)
	}

	byTask, err := s.ListToolCallsByTask(ctx, "task-a")
	if err != nil || len(byTask) != 1 {
		t.Fatalf("by-task list = %d, %v", len(byTask), err)
	}
	all, err := s.ListToolCalls(ctx)
	if err != nil || len(all) != 1 {
		t.Fatalf("list = %d, %v", len(all), err)
	}
	if err := s.DeleteToolCall(ctx, id); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteToolCall(ctx, id); !errors.Is(err, ErrNotFound) {
		t.Errorf("delete missing tool_call: %v", err)
	}
}

func TestGrantCostPluginStateDAO(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	now := time.Now().UTC().Unix()

	gid, err := s.InsertGrant(ctx, ApprovalGrant{
		Scope: GrantScopeSession, Tool: "fs.write",
		Pattern: `C:\Users\swq\docs\*`, SessionID: "sess-1",
		CreatedAt: now, ExpiresAt: now + 3600,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.InsertGrant(ctx, ApprovalGrant{Scope: "global", Tool: "x", Pattern: "y", SessionID: "s", ExpiresAt: now}); err == nil {
		t.Error("non-session scope accepted")
	}
	grants, err := s.ListGrantsBySession(ctx, "sess-1")
	if err != nil || len(grants) != 1 {
		t.Fatalf("grants = %d, %v", len(grants), err)
	}
	if err := s.RevokeGrant(ctx, gid); err != nil {
		t.Fatal(err)
	}
	if err := s.RevokeGrant(ctx, gid); err != nil { // idempotent
		t.Errorf("re-revoke: %v", err)
	}
	got2, _ := s.ListGrants(ctx)
	if got2[0].RevokedAt == nil {
		t.Error("revoked_at not stamped")
	}
	if err := s.RevokeGrant(ctx, 999); !errors.Is(err, ErrNotFound) {
		t.Errorf("revoke missing grant: %v", err)
	}

	// cost_daily: upsert-add semantics.
	if err := s.BumpCostDay(ctx, "2026-09-19", 100, 50, 300); err != nil {
		t.Fatal(err)
	}
	if err := s.BumpCostDay(ctx, "2026-09-19", 10, 5, 30); err != nil {
		t.Fatal(err)
	}
	if err := s.BumpCostDay(ctx, "2026-9-9", 1, 1, 1); err == nil {
		t.Error("malformed day accepted")
	}
	c, err := s.CostDay(ctx, "2026-09-19")
	if err != nil {
		t.Fatal(err)
	}
	if c.TokensIn != 110 || c.TokensOut != 55 || c.CostMicros != 330 || c.Tasks != 2 {
		t.Errorf("aggregate wrong: %+v", c)
	}
	if c2, _ := s.CostDay(ctx, "2026-01-01"); c2.Tasks != 0 {
		t.Errorf("absent day should read zero-value, got %+v", c2)
	}

	// plugin_state incl. exe_hash.
	if err := s.UpsertPluginState(ctx, PluginState{
		ID: "plugin-a", Version: "1.2.3", Enabled: true,
		CapabilitiesJSON: `["fs.read"]`, InstalledAt: now,
		Hash: "abc123", ExeHash: "def456",
	}); err != nil {
		t.Fatal(err)
	}
	ps, err := s.PluginStateByID(ctx, "plugin-a")
	if err != nil {
		t.Fatal(err)
	}
	if !ps.Enabled || ps.ExeHash != "def456" {
		t.Errorf("plugin_state roundtrip wrong: %+v", ps)
	}
	if err := s.UpsertPluginState(ctx, PluginState{ID: "plugin-a", Version: "1.2.4", Enabled: false, CapabilitiesJSON: `["fs.read"]`, Hash: "h2"}); err != nil {
		t.Fatal(err)
	}
	plist, _ := s.ListPluginStates(ctx)
	if len(plist) != 1 || plist[0].Version != "1.2.4" {
		t.Errorf("upsert did not replace: %+v", plist)
	}
	if err := s.DeletePluginState(ctx, "plugin-a"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.PluginStateByID(ctx, "plugin-a"); !errors.Is(err, ErrNotFound) {
		t.Errorf("get missing plugin_state: %v", err)
	}
}

// TestWriteQueueLazyLifecycle pins the db-writer lifecycle decision: the
// goroutine exists only while there is work (D20 on-demand + roster baseline).
func TestWriteQueueLazyLifecycle(t *testing.T) {
	logger, _ := bufLogger(t)
	s := openTestStore(t, WithLogger(logger), WithRegistry(observe.NewRegistry()))
	ctx := context.Background()

	if got := s.reg.CountByName("db-writer"); got != 0 {
		t.Fatalf("db-writer alive before any write: %d", got)
	}
	if _, err := s.UpsertProfile(ctx, "pref.x", "x", ProfileManual); err != nil {
		t.Fatal(err)
	}
	// The writer exits as soon as the queue drains; poll briefly.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if s.queue.idle() && s.reg.CountByName("db-writer") == 0 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if !s.queue.idle() || s.reg.CountByName("db-writer") != 0 {
		t.Error("db-writer did not go idle after drain (resident leak)")
	}

	// Burst of writes still serializes correctly.
	errs := make(chan error, 16)
	for i := 0; i < 16; i++ {
		go func(i int) {
			_, err := s.UpsertProfile(ctx, fmt.Sprintf("pref.b%02d", i), "x", ProfileExtracted)
			errs <- err
		}(i)
	}
	for i := 0; i < 16; i++ {
		if err := <-errs; err != nil {
			t.Fatalf("burst write: %v", err)
		}
	}
	n, err := s.PurgeProfiles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if n != 17 { // pref.x + 16 burst rows, still under the cap
		t.Errorf("purged %d, want 17", n)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if err := s.write(ctx, func(context.Context, *sql.Tx) error { return nil }); !errors.Is(err, ErrStoreClosed) {
		t.Errorf("write after close = %v, want ErrStoreClosed", err)
	}
	// Idempotent close.
	if err := s.Close(); err != nil {
		t.Errorf("second close: %v", err)
	}
}
