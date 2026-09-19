package observe

import (
	"time"
)

// Monotonic-clock discipline (D42#9; D22 ban list).
//
// Every timeout, deadline, retention window and interval comparison MUST be
// computed on the monotonic clock. In Go this means: keep time.Time values
// obtained from time.Now (they carry a monotonic reading), use time.Since /
// time.Until / Timer / Ticker / context deadlines, and this package's
// Timeout. Computing intervals from wall-clock readings (Unix seconds,
// formatted timestamps, values that went through encoding/decoding or
// .Round(0)) is BANNED: NTP corrections, timezone changes and sleep/resume
// (D42#1) make wall-clock deltas arbitrarily wrong.
//
// Wall-clock values are for PERSISTED RECORDS ONLY (log timestamps,
// task_log, cost ledger, artifacts): use NowWallUTC / WallTimestampUTC.

// Timeout is a monotonic budget for disposal scopes, task steps and shutdown
// steps. All comparisons are monotonic (time.Since/time.Until on a start
// value that carries the monotonic reading), so system clock jumps cannot
// shorten or extend the budget.
type Timeout struct {
	start  time.Time // carries the monotonic reading
	budget time.Duration
}

// NewTimeout starts a monotonic budget of d. d <= 0 yields an expired budget
// (useful for "already past deadline" branches).
func NewTimeout(d time.Duration) Timeout {
	return Timeout{start: time.Now(), budget: d}
}

// Budget returns the configured total budget.
func (t Timeout) Budget() time.Duration { return t.budget }

// Elapsed returns the monotonic time since the budget started.
func (t Timeout) Elapsed() time.Duration { return time.Since(t.start) }

// Remaining returns the monotonic remaining budget, clamped at 0.
func (t Timeout) Remaining() time.Duration {
	r := t.budget - t.Elapsed()
	if r < 0 {
		return 0
	}
	return r
}

// Expired reports whether the budget is exhausted.
func (t Timeout) Expired() bool { return t.Remaining() <= 0 }

// Deadline returns the deadline as a time.Time that still carries the
// monotonic reading; it is safe for time.Until and context.WithDeadline.
func (t Timeout) Deadline() time.Time { return t.start.Add(t.budget) }

// NowWallUTC returns the current wall-clock UTC time for persistence only.
// It MUST NOT be used to compute intervals (see the package clock rules).
func NowWallUTC() time.Time { return time.Now().UTC() }

// WallTimestampUTC renders an instant for persisted records: RFC3339, UTC,
// second precision, monotonic reading stripped. Only for storage/display,
// never for interval math.
func WallTimestampUTC(t time.Time) string {
	return t.Round(0).UTC().Format(time.RFC3339)
}
