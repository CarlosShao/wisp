package memory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/plugin"
)

// seedTaskLogAt inserts a task_log row at an absolute wall timestamp.
func seedTaskLogAt(t *testing.T, s *Store, id string, ts time.Time) {
	t.Helper()
	err := s.write(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO task_log(id, started_at, state, query_text) VALUES(?,?,'succeeded','q')`,
			id, ts.Unix())
		return err
	})
	if err != nil {
		t.Fatalf("seed task_log %s: %v", id, err)
	}
}

// seedToolCallAt inserts a tool_call row at an absolute wall timestamp.
func seedToolCallAt(t *testing.T, s *Store, taskID string, seq int64, ts time.Time) {
	t.Helper()
	err := s.write(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO tool_call(task_id, seq, tool, args_json, risk_level, correlation_id, started_at)
			 VALUES(?,?,?,?,?,?,?)`,
			taskID, seq, "fs.read", "{}", RiskL0, fmt.Sprintf("corr-%d", seq), ts.Unix())
		return err
	})
	if err != nil {
		t.Fatalf("seed tool_call: %v", err)
	}
}

// seedToolCallFull inserts a tool_call row with explicit started/ended
// instants (hasEnded=false leaves ended_at NULL, decided_at always NULL).
func seedToolCallFull(t *testing.T, s *Store, taskID string, seq int64, started time.Time, ended time.Time, hasEnded bool) {
	t.Helper()
	var endedAt any
	if hasEnded {
		endedAt = ended.Unix()
	}
	err := s.write(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO tool_call(task_id, seq, tool, args_json, risk_level, correlation_id, started_at, ended_at)
			 VALUES(?,?,?,?,?,?,?,?)`,
			taskID, seq, "fs.read", "{}", RiskL0, fmt.Sprintf("corr-full-%d", seq), started.Unix(), endedAt)
		return err
	})
	if err != nil {
		t.Fatalf("seed tool_call full: %v", err)
	}
}

// seedCostDayAt inserts a cost_daily row for an absolute day.
func seedCostDayAt(t *testing.T, s *Store, day string) {
	t.Helper()
	err := s.write(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO cost_daily(day, tokens_in, tokens_out, cost_micros, tasks) VALUES(?,1,1,1,1)`, day)
		return err
	})
	if err != nil {
		t.Fatalf("seed cost_daily %s: %v", day, err)
	}
}

// seedGrantAt inserts a grant with absolute expires/revoked instants
// (nilPtr => NULL revoked_at).
func seedGrantAt(t *testing.T, s *Store, session string, expires, revoked time.Time, hasRevoked bool) {
	t.Helper()
	var revokedAt any
	if hasRevoked {
		revokedAt = revoked.Unix()
	}
	err := s.write(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO approval_grant(scope, tool, pattern, session_id, created_at, expires_at, revoked_at)
			 VALUES('session','fs.read','p',?,0,?,?)`,
			session, expires.Unix(), revokedAt)
		return err
	})
	if err != nil {
		t.Fatalf("seed grant: %v", err)
	}
}

func countRows(t *testing.T, s *Store, table string) int64 {
	t.Helper()
	var n int64
	if err := s.reader.QueryRow(`SELECT count(*) FROM ` + table).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// TestRetentionBoundaries covers the SPEC-02 §8 boundary matrix with a pinned
// clock: rows exactly AT the window stay, strictly older rows go
// (29/30/31 days; 399/400/401 days; grant audit window).
func TestRetentionBoundaries(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	day := 24 * time.Hour

	// Pin the clock AFTER seeding so second-granularity unix timestamps stay
	// exact: every seed is computed against this same instant, and the
	// retention pass reads it via RetentionConfig.Now.
	now := observe.NowWallUTC()
	ago := func(d time.Duration) time.Time { return now.Add(-d) }

	seedTaskLogAt(t, s, "t-29", ago(29*day))
	seedTaskLogAt(t, s, "t-30", ago(30*day))
	seedTaskLogAt(t, s, "t-31", ago(31*day))

	seedToolCallAt(t, s, "t-29", 1, ago(29*day))
	seedToolCallAt(t, s, "t-30", 2, ago(30*day))
	seedToolCallAt(t, s, "t-31", 3, ago(31*day))
	// NULL ended_at rows must still age (MINOR-4: no immortal NULLs) and
	// finished rows age by ended_at, not the older started_at.
	seedToolCallFull(t, s, "tc-nullended-31", 4, ago(31*day), time.Time{}, false)
	seedToolCallFull(t, s, "tc-ended-31", 5, ago(32*day), ago(31*day), true)
	seedToolCallFull(t, s, "tc-ended-29", 6, ago(30*day), ago(29*day), true)

	seedCostDayAt(t, s, now.Add(-399*day).Format("2006-01-02"))
	seedCostDayAt(t, s, now.Add(-400*day).Format("2006-01-02"))
	seedCostDayAt(t, s, now.Add(-401*day).Format("2006-01-02"))

	// Grants: created at now-60d, expired at now-60d; revoked (or expired
	// alone) the given ages ago. Audit window = 30d after invalidation.
	seedGrantAt(t, s, "g-exp-29", ago(60*day), ago(29*day), true)
	seedGrantAt(t, s, "g-exp-30", ago(60*day), ago(30*day), true)
	seedGrantAt(t, s, "g-exp-31", ago(60*day), ago(31*day), true)
	// Expired-but-never-revoked: invalid_at = expires_at.
	seedGrantAt(t, s, "g-norevoke-31", ago(31*day), time.Time{}, false)
	// Still-valid future grant survives everything.
	seedGrantAt(t, s, "g-future", now.Add(24*time.Hour), time.Time{}, false)

	res := s.RunRetentionOnce(ctx, RetentionConfig{Now: func() time.Time { return now }})

	if res.TaskLogDeleted != 1 {
		t.Errorf("task_log deleted = %d, want 1 (the 31d row)", res.TaskLogDeleted)
	}
	if res.ToolCallDeleted != 3 {
		t.Errorf("tool_call deleted = %d, want 3 (31d row + NULL-ended 31d + ended-31d)", res.ToolCallDeleted)
	}
	if res.CostDayDeleted != 1 {
		t.Errorf("cost_daily deleted = %d, want 1 (the 401d day)", res.CostDayDeleted)
	}
	if res.GrantDeleted != 2 {
		t.Errorf("grant deleted = %d, want 2 (revoked-31d + expired-31d)", res.GrantDeleted)
	}

	// 29d and 30d rows are kept; the 31d row is gone.
	for _, id := range []string{"t-29", "t-30"} {
		if _, err := s.TaskLogByID(ctx, id); err != nil {
			t.Errorf("boundary row %s must be kept: %v", id, err)
		}
	}
	if _, err := s.TaskLogByID(ctx, "t-31"); !errors.Is(err, ErrNotFound) {
		t.Errorf("31d row must be deleted, got %v", err)
	}
	if got := countRows(t, s, "tool_call"); got != 3 {
		t.Errorf("tool_call rows = %d, want 3 (t-29 + t-30 + tc-ended-29)", got)
	}
	if got := countRows(t, s, "cost_daily"); got != 2 {
		t.Errorf("cost_daily rows = %d, want 2 (399d+400d)", got)
	}
	if got := countRows(t, s, "approval_grant"); got != 3 {
		t.Errorf("grant rows = %d, want 3 (29d+30d+future)", got)
	}
}

// TestRetentionSchedule pins the contractual schedule: one run after
// FirstDelay, then every Period (shortened for the test, monotonic timers),
// and DisposalScope governance: after Dispose the goroutine is joined (C11
// step 2), so no further pass can start.
func TestRetentionSchedule(t *testing.T) {
	logger, buf := bufLogger(t)
	s := openTestStore(t, WithLogger(logger))
	reg := observe.NewRegistry()
	scope := plugin.NewDisposalScope("test-retention", nil,
		plugin.WithRegistry(reg), plugin.WithTimeout(2*time.Second))
	defer scope.Dispose()

	cfg := RetentionConfig{FirstDelay: 30 * time.Millisecond, Period: 80 * time.Millisecond}
	s.StartRetentionJob(scope, cfg)

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) && strings.Count(buf.String(), "retention pass done") < 2 {
		time.Sleep(5 * time.Millisecond)
	}
	n := strings.Count(buf.String(), "retention pass done")
	if n < 2 {
		t.Fatalf("retention passes = %d, want >= 2 (first delay + period)", n)
	}

	scope.Dispose() // C11: cancel + join; returns only after the goroutine exits
	time.Sleep(3 * scopeCfgPeriod(cfg))
	after := strings.Count(buf.String(), "retention pass done")
	if after != n {
		t.Errorf("retention pass ran after Dispose (%d -> %d): goroutine not governed", n, after)
	}
}

// scopeCfgPeriod guards against a zero period in the test above.
func scopeCfgPeriod(c RetentionConfig) time.Duration {
	if c.Period <= 0 {
		return time.Second
	}
	return c.Period
}

// TestArtifactsLRUQuota pins the LRU behavior with a tiny quota (the
// production 500MB constant is ArtifactsQuotaBytes).
func TestArtifactsLRUQuota(t *testing.T) {
	logger, buf := bufLogger(t)
	s := openTestStore(t, WithLogger(logger))
	ctx := context.Background()

	mk := func(name string, size int, age time.Duration) {
		t.Helper()
		p := filepath.Join(s.ArtifactsDir(), name)
		if err := os.WriteFile(p, make([]byte, size), 0o644); err != nil {
			t.Fatal(err)
		}
		past := time.Now().Add(-age)
		if err := os.Chtimes(p, past, past); err != nil {
			t.Fatal(err)
		}
	}
	mk("old.txt", 100, 3*time.Hour)
	mk("mid.txt", 100, 2*time.Hour)
	mk("new.txt", 100, 1*time.Hour)

	files, freed, err := s.enforceArtifactsQuota(ctx, 250) // 300 bytes on disk, 50 over
	if err != nil {
		t.Fatal(err)
	}
	if files != 1 || freed != 100 {
		t.Fatalf("LRU removed (%d, %d), want (1, 100): oldest file only", files, freed)
	}
	if _, err := os.Stat(filepath.Join(s.ArtifactsDir(), "old.txt")); !os.IsNotExist(err) {
		t.Errorf("oldest artifact survived LRU: %v", err)
	}
	if _, err := os.Stat(filepath.Join(s.ArtifactsDir(), "mid.txt")); err != nil {
		t.Errorf("mid artifact wrongly removed: %v", err)
	}
	if !strings.Contains(buf.String(), "artifacts LRU cleanup") {
		t.Error("LRU cleanup not logged")
	}

	// Quota >= size: no-op.
	files, freed, err = s.enforceArtifactsQuota(ctx, 1<<30)
	if err != nil || files != 0 || freed != 0 {
		t.Errorf("no-op LRU removed (%d, %d, %v), want zeros", files, freed, err)
	}
}

func TestArtifactsDeleteAndPurgeAndGuards(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	write := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(s.ArtifactsDir(), name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("a.txt", "hello")
	write("b.txt", "world")

	arts, err := s.ListArtifacts(ctx)
	if err != nil || len(arts) != 2 {
		t.Fatalf("list artifacts = %v, %v", arts, err)
	}
	if err := s.DeleteArtifact(ctx, "a.txt"); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteArtifact(ctx, "a.txt"); !errors.Is(err, ErrNotFound) {
		t.Errorf("delete missing artifact: %v", err)
	}
	// Path traversal guards: only bare file names are valid.
	for _, bad := range []string{"../evil.txt", `sub\dir.txt`, "..", "", "a:b.txt"} {
		if err := s.DeleteArtifact(ctx, bad); err == nil {
			t.Errorf("traversal name %q not rejected", bad)
		}
	}
	n, err := s.PurgeArtifacts(ctx)
	if err != nil || n != 1 {
		t.Fatalf("purge = %d, %v; want 1", n, err)
	}
	arts, _ = s.ListArtifacts(ctx)
	if len(arts) != 0 {
		t.Errorf("purge left %d artifacts", len(arts))
	}
}
