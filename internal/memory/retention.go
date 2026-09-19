package memory

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/plugin"
)

// RetentionJob (D35 rule 1, SPEC-02 §4): retention windows need an owner or
// they silently never run.
//
// Schedule (contractual): once 5 minutes after process start, then every 24h
// — a resident process may only live 10 minutes before shutdown, so the
// periodic timer alone would NEVER fire; the short first delay is the
// guarantee. The retention-job goroutine is temporary (D38b roster) and
// governed by a DisposalScope (C11): disposal cancels its ctx and joins it.
//
// Clock discipline (D42#9): the schedule and every in-process elapsed
// comparison use the monotonic clock (time.Timer/Ticker/observe.Timeout —
// wall-clock deltas are banned). Row age itself is necessarily a comparison
// against PERSISTED wall-clock seconds (there is no other frame for "30 days
// ago" once the process restarted); that is a semantic of the stored data,
// not elapsed-time measurement.

// Retention windows (SPEC-02 §4 table).
const (
	// TaskLogTTL / ToolCallTTL: L3 task logs and tool forensics live 30 days.
	TaskLogTTL = 30 * 24 * time.Hour
	// CostTTL: the C23 daily aggregation lives 400 days.
	CostTTL = 400 * 24 * time.Hour
	// GrantAuditTTL: expired/revoked grants stay 30 days for audit, then go.
	GrantAuditTTL = 30 * 24 * time.Hour
	// ArtifactsQuotaBytes: the artifacts directory LRU quota (500MB).
	ArtifactsQuotaBytes = 500 * 1024 * 1024

	defaultFirstDelay = 5 * time.Minute
	defaultPeriod     = 24 * time.Hour
)

// RetentionConfig carries the contractual windows; zero fields fall back to
// the defaults above. Tests shorten the schedule and pin the cutoffs.
type RetentionConfig struct {
	FirstDelay time.Duration // default 5m after boot
	Period     time.Duration // default 24h
	TaskLogTTL time.Duration // default TaskLogTTL
	CostTTL    time.Duration // default CostTTL
	GrantTTL   time.Duration // audit window after expiry; default GrantAuditTTL
	Artifacts  int64         // artifacts quota bytes; default ArtifactsQuotaBytes

	// Now supplies the wall-clock reading for cutoff computation. Nil =
	// observe.NowWallUTC (production). Tests pin it so window boundaries are
	// exact despite second-granularity unix timestamps.
	Now func() time.Time
}

func (c RetentionConfig) withDefaults() RetentionConfig {
	if c.FirstDelay <= 0 {
		c.FirstDelay = defaultFirstDelay
	}
	if c.Period <= 0 {
		c.Period = defaultPeriod
	}
	if c.TaskLogTTL <= 0 {
		c.TaskLogTTL = TaskLogTTL
	}
	if c.CostTTL <= 0 {
		c.CostTTL = CostTTL
	}
	if c.GrantTTL <= 0 {
		c.GrantTTL = GrantAuditTTL
	}
	if c.Artifacts <= 0 {
		c.Artifacts = ArtifactsQuotaBytes
	}
	if c.Now == nil {
		c.Now = observe.NowWallUTC
	}
	return c
}

// RetentionResult is one pass's outcome, logged and returned for tests.
type RetentionResult struct {
	TaskLogDeleted  int64
	ToolCallDeleted int64
	CostDayDeleted  int64
	GrantDeleted    int64
	ArtifactFiles   int
	ArtifactBytes   int64
}

// StartRetentionJob launches the retention-job goroutine inside the given
// DisposalScope (C11: the scope's Dispose cancels and joins it). It returns
// immediately.
func (s *Store) StartRetentionJob(scope *plugin.DisposalScope, cfg RetentionConfig) {
	if scope == nil {
		panic("memory: StartRetentionJob requires a DisposalScope")
	}
	c := cfg.withDefaults()
	scope.Go("retention-job", func(ctx context.Context) {
		s.retentionLoop(ctx, c)
	})
}

func (s *Store) retentionLoop(ctx context.Context, c RetentionConfig) {
	s.logger.Info("retention job scheduled",
		"component", "memory",
		"first_delay", c.FirstDelay.String(), "period", c.Period.String())

	// Monotonic schedule: Timer + Ticker (Go timers run on the monotonic
	// clock; wall-clock differences are banned by D42#9).
	timer := time.NewTimer(c.FirstDelay)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("retention job stopped", "component", "memory")
			return
		case <-timer.C:
		}

		s.runRetention(ctx, c)

		// Re-arm the 24h period unless we are shutting down. After the first
		// run the loop uses a plain ticker cadence via timer reset.
		if ctx.Err() != nil {
			return
		}
		timer.Reset(c.Period)
	}
}

// runRetention executes one cleanup pass. Windows: task_log/tool_call older
// than 30 days, cost_daily older than 400 days, grants invalid for more than
// the audit window, and the artifacts LRU over its quota.
func (s *Store) runRetention(ctx context.Context, c RetentionConfig) (res RetentionResult) {
	start := observe.NewTimeout(0)
	s.logger.Info("retention pass started", "component", "memory")

	err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		now := c.Now()
		res.TaskLogDeleted, err = deleteOlderThan(ctx, tx, "task_log", "started_at", c.TaskLogTTL, now)
		if err != nil {
			return err
		}
		// tool_call ages by its start, falling back to the decision time for
		// rows captured mid-approval.
		res.ToolCallDeleted, err = deleteOlderThan(ctx, tx, "tool_call",
			"COALESCE(started_at, decided_at)", c.TaskLogTTL, now)
		if err != nil {
			return err
		}
		res.CostDayDeleted, err = deleteCostOlderThan(ctx, tx, c.CostTTL, now)
		if err != nil {
			return err
		}
		res.GrantDeleted, err = deleteExpiredGrants(ctx, tx, c.GrantTTL, now)
		return err
	})
	if err != nil {
		s.logger.Error("retention pass db cleanup failed",
			"component", "memory", "err", err,
			"error_class", string(observe.ClassOfOrInternal(err)))
	}

	af, ab, aerr := s.enforceArtifactsQuota(ctx, c.Artifacts)
	if aerr != nil {
		s.logger.Error("retention pass artifacts cleanup failed",
			"component", "memory", "err", aerr,
			"error_class", string(observe.ClassOfOrInternal(aerr)))
	}
	res.ArtifactFiles, res.ArtifactBytes = af, ab

	s.logger.Info("retention pass done",
		"component", "memory",
		"task_log_deleted", res.TaskLogDeleted,
		"tool_call_deleted", res.ToolCallDeleted,
		"cost_day_deleted", res.CostDayDeleted,
		"grant_deleted", res.GrantDeleted,
		"artifact_files_removed", res.ArtifactFiles,
		"artifact_bytes_removed", res.ArtifactBytes,
		"duration_ms", start.Elapsed().Milliseconds())
	return res
}

// deleteOlderThan removes rows whose timestamp column (possibly an
// expression) is strictly older than now-ttl. A row at exactly the window
// boundary (age == ttl) is KEPT: ">30 days" is strictly greater.
func deleteOlderThan(ctx context.Context, tx *sql.Tx, table, tsExpr string, ttl time.Duration, now time.Time) (int64, error) {
	cutoff := now.Add(-ttl).Unix()
	q := fmt.Sprintf(`DELETE FROM %s WHERE %s IS NOT NULL AND %s < ?`, table, tsExpr, tsExpr)
	res, err := tx.ExecContext(ctx, q, cutoff)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// deleteCostOlderThan removes cost_daily rows whose 'YYYY-MM-DD' day is
// strictly older than the 400-day cutoff (string comparison is correct for
// zero-padded ISO days).
func deleteCostOlderThan(ctx context.Context, tx *sql.Tx, ttl time.Duration, now time.Time) (int64, error) {
	cutoffDay := now.Add(-ttl).Format("2006-01-02")
	res, err := tx.ExecContext(ctx, `DELETE FROM cost_daily WHERE day < ?`, cutoffDay)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// deleteExpiredGrants removes grants that have been invalid (expired or
// revoked) for longer than the audit window: an invalid grant stays 30 days
// for audit, then goes (SPEC-02 §4).
func deleteExpiredGrants(ctx context.Context, tx *sql.Tx, ttl time.Duration, now time.Time) (int64, error) {
	cutoff := now.Add(-ttl).Unix()
	// invalid_at = MAX(expires_at, revoked_at if present); MAX() ignores
	// NULLs, so COALESCE keeps revoked rows judged by the later of the two.
	res, err := tx.ExecContext(ctx,
		`DELETE FROM approval_grant
		 WHERE MAX(expires_at, COALESCE(revoked_at, 0)) < ?`, cutoff)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// RunRetentionOnce executes one pass synchronously (tests and the doctor
// command); production scheduling goes through StartRetentionJob.
func (s *Store) RunRetentionOnce(ctx context.Context, cfg RetentionConfig) RetentionResult {
	return s.runRetention(ctx, cfg.withDefaults())
}
