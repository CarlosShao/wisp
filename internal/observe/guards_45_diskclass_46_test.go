package observe

// Ticket 46's post-audit leg, item 3's tooth: the trip is classified.
//
// The 46b audit (docs/evidence/s1/46b-audit-20261115.md "另记 2 条" #2) measured
// guards_45.go's tripped branch returning a BARE fmt.Errorf wrapping
// ErrDiskSpaceLow - readable via errors.Is, but carrying no D37 class at any
// boundary, so the S6 pre-mortem disk-full row could not hand over "correct D37
// classes". The fix wraps that error with Wrap(ClassResource, ...); this case
// pins BOTH halves of that contract as behaviour:
//
//   - the class reads off the refusal (Class() == ClassResource): a reversion
//     to the bare fmt.Errorf deletes the class and goes red here;
//   - errors.Is(err, ErrDiskSpaceLow) still hits: the wrap must never buy the
//     class by spending the sentinel's identity (errors.go's Error.Unwrap
//     exposes the cause chain; losing it would break every existing caller).
//
// Red sentences name the missing reading, not a source string.

import (
	"errors"
	"testing"
	"time"
)

// Test46DiskGuardTripIsClassifiedResource drives the guard's own admission
// rule (the production path a gated write takes) and reads the refusal the
// way the D37 boundary reads it.
func Test46DiskGuardTripIsClassifiedResource(t *testing.T) {
	probe := &varProbe{free: 199 << 20} // 1MiB below the frozen 200MiB threshold
	clock := &fixedClock{at: time.Unix(1780000000, 0)}
	g := newGuard45(t, probe, clock, DiskGuardThresholdDefault)

	err := g.Check(WriteLog)
	if err == nil {
		t.Fatal("46: a write at 199MiB free was admitted: the guard did not trip, so there is nothing to classify")
	}
	cls, ok := ClassOf(err)
	if !ok || cls != ClassResource {
		t.Fatalf("46: the tripped refusal reads as class %q (ok=%t), want %q: the disk-full trip carries no D37 class at the boundary", string(cls), ok, string(ClassResource))
	}
	if !errors.Is(err, ErrDiskSpaceLow) {
		t.Fatalf("46: errors.Is(refusal, ErrDiskSpaceLow) = false (refusal: %v): the classification was bought by losing the sentinel's identity", err)
	}
}
