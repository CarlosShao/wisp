package observe

import (
	"testing"
	"time"
)

// TestTimeoutIsMonotonic checks the Timeout budget arithmetic. (The
// monotonic property itself comes from time.Since on a start value carrying
// the monotonic reading; faking a clock jump in-process is not possible, so
// this test verifies the wrapper never touches wall-clock fields.)
func TestTimeoutIsMonotonic(t *testing.T) {
	tm := NewTimeout(60 * time.Millisecond)
	if tm.Expired() {
		t.Fatal("fresh 60ms budget must not be expired")
	}
	if got := tm.Remaining(); got <= 0 || got > 60*time.Millisecond {
		t.Fatalf("Remaining() = %v, want (0,60ms]", got)
	}
	if got := tm.Elapsed(); got < 0 {
		t.Fatalf("Elapsed() = %v, want >= 0", got)
	}
	time.Sleep(90 * time.Millisecond)
	if !tm.Expired() {
		t.Fatal("budget must be expired after 90ms")
	}
	if got := tm.Remaining(); got != 0 {
		t.Fatalf("Remaining() after expiry = %v, want 0", got)
	}
	if got := tm.Elapsed(); got < 90*time.Millisecond {
		t.Fatalf("Elapsed() = %v, want >= 90ms", got)
	}
	if tm.Budget() != 60*time.Millisecond {
		t.Fatalf("Budget() = %v", tm.Budget())
	}

	// Deadline keeps the monotonic reading: time.Until on it must be usable
	// and non-positive after expiry.
	d := tm.Deadline()
	if time.Until(d) > 0 {
		t.Fatal("deadline of expired budget must be in the past (monotonically)")
	}

	// Zero and negative budgets start expired.
	if !NewTimeout(0).Expired() || !NewTimeout(-time.Second).Expired() {
		t.Fatal("non-positive budgets must start expired")
	}
}

// TestWallTimestampUTC verifies the persisted-timestamp helper: RFC3339 UTC,
// second precision, monotonic reading stripped.
func TestWallTimestampUTC(t *testing.T) {
	now := NowWallUTC()
	if now.Location() != time.UTC {
		t.Fatal("NowWallUTC must return UTC")
	}
	ts := WallTimestampUTC(now)
	parsed, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		t.Fatalf("WallTimestampUTC produced unparseable value %q: %v", ts, err)
	}
	if parsed.Unix() != now.Round(0).Unix() {
		t.Fatalf("round trip mismatch: %q vs %v", ts, now)
	}
	if ts != WallTimestampUTC(parsed) {
		t.Fatal("timestamp not idempotent")
	}
	// Non-UTC input is normalized to UTC with the same instant.
	local := now.In(time.FixedZone("TEST", 8*3600))
	if got := WallTimestampUTC(local); got != ts {
		t.Fatalf("WallTimestampUTC(non-UTC) = %q, want %q", got, ts)
	}
}
