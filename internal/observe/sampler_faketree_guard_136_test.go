package observe

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

// Ticket 136 AC#13 - the second instrument that used to swallow readings.
//
// AC#8 nailed the `rep.Samples[0]` direct read in
// TestSamplerGoroutineAccountingFollowsRegistry. R-136-9 is the same family in
// the seam itself: fakeTree.ReadTree answered "script exhausted AND nothing
// was ever scripted" with f.mu[len(f.mu)-1], i.e. f.mu[-1] on a bare
// &fakeTree{} - a panic inside the test fixture. A panic like that does not
// just fail one case: it takes the test binary down, so every case scheduled
// after it reports neither red nor green (evidence §1.1: baseline 58 RUN, with
// one ordinary bare-&fakeTree{} sampling case added 43 RUN, 16 named cases
// swallowed).
//
// The ruling that kept this ticket honest is that nothing in production
// reachable today hits the branch (the two bare &fakeTree{} constructions stop
// before the first read: sampler.go:253-255 rejects an unknown state, and
// TestMarkTransitionTimestamps never reads a tree). So this is the landmine
// the *next* person adds a case over. The fix therefore only guards the seam
// (sampler_test.go:31) - the sampler's semantics are untouched.
//
// Legs, and what mutation turns each one red (evidence §1.4):
//   - the seam itself must return an error, never panic            (removing
//     the guard makes this leg report the panic value it recovered);
//   - a scripted seam keeps repeating its last reading unchanged   (guard
//     written too wide - e.g. failing on every read - turns this red);
//   - an ordinary sampling case over the misuse goes red through the
//     error channel and the package still finishes.

// recoverRead calls the seam once and reports a panic as this case's own
// failure instead of letting it kill the binary. In the fixed state nothing
// panics, so this never hides anything: if the guard is removed the caller
// still gets the panic value and the leg that asked for it fails red.
func recoverRead(t *testing.T, f *fakeTree) (m TreeMetrics, err error) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("fakeTree.ReadTree panicked instead of failing closed: %v "+
				"(a panic here takes the whole package's readings with it)", r)
		}
	}()
	return f.ReadTree()
}

// TestFakeTreeEmptyScriptFailsClosedAndNotPanics is the seam leg.
func TestFakeTreeEmptyScriptFailsClosedAndNotPanics(t *testing.T) {
	m, err := recoverRead(t, &fakeTree{})
	if err == nil {
		t.Fatalf("bare &fakeTree{} returned a reading it never measured (metrics=%+v): the seam must fail closed", m)
	}
	if !errors.Is(err, errFakeTreeNoScript) {
		t.Fatalf("bare &fakeTree{} must fail with the seam's own sentence, got %v", err)
	}
	if !strings.Contains(err.Error(), "no scripted reads") {
		t.Fatalf("the failure must say what is missing, got %q", err.Error())
	}
	if m != (TreeMetrics{}) {
		t.Fatalf("a refused read must not carry invented metrics: %+v", m)
	}
}

// TestFakeTreeScriptedReadsStillRepeatLastReading is the control leg that keeps
// the guard from being written too wide: the documented behaviour of a scripted
// fake (run past the end of the script and you keep getting the last reading)
// must survive. sampler_test.go:238 and the mu-driven cases depend on it.
func TestFakeTreeScriptedReadsStillRepeatLastReading(t *testing.T) {
	f := &fakeTree{mu: []TreeMetrics{{PIDs: 1, PrivateWorkingSetBytes: 4 << 20}}}
	first, err := recoverRead(t, f)
	if err != nil {
		t.Fatalf("scripted read must not fail: %v", err)
	}
	if first.PrivateWorkingSetBytes != 4<<20 {
		t.Fatalf("scripted read wrong: %+v", first)
	}
	for i := 0; i < 3; i++ {
		again, err := recoverRead(t, f)
		if err != nil {
			t.Fatalf("exhausted script %d: the last reading must still be repeated, got err=%v", i, err)
		}
		if again != first {
			t.Fatalf("exhausted script must repeat the last reading, got %+v want %+v", again, first)
		}
	}
}

// TestBareFakeTreeNormalSamplingCaseGoesRedThroughErrorChannel is the shape the
// next person writes: an ordinary sampling case over a bare &fakeTree{}. It
// must come back as an error the case can fail on - not as a panic - and the
// sampler's semantics must be what carries it (sampler.go:269-272 rejects the
// failed baseline read; nothing here loosens that).
func TestBareFakeTreeNormalSamplingCaseGoesRedThroughErrorChannel(t *testing.T) {
	s := NewSampler(&fakeTree{}, NewRegistry())
	rep, err := s.SampleState(context.Background(), SLOSleeping, 10*time.Millisecond, 40*time.Millisecond)
	if err == nil {
		t.Fatalf("a sampling window over an unscripted tree must error, got %+v", rep)
	}
	if rep != nil {
		t.Fatalf("the failed window must not hand back a report: %+v", rep)
	}
	if !errors.Is(err, errFakeTreeNoScript) {
		t.Fatalf("the seam's sentence must reach the caller, got %v", err)
	}
}
