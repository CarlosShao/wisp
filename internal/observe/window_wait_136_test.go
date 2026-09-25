package observe

import (
	"testing"
	"time"
)

// Ticket 136 AC#15 - the "count without waiting" precondition, fixed once.
//
// Twelve legs in this package state part of their precondition as "this window
// took at least N reads" or "recorded at least N samples" and then read that
// count the instant the window closes. The count is not theirs to have yet:
// CheckSettle only ever reads on a ticker tick (sampler.go:493-525 - there is
// no read before the first tick, unlike SampleState, which reads once up front
// at sampler.go:269 and once to close at :306), so the number of reads a window
// returns is a property of how the Go scheduler delivered ticks during those
// 100-200ms, not of the fixture. Reading it without waiting for it is the same
// defect ticket 136 AC#11 measured and fixed on the registry side
// (goroutine_test.go:44-66): a guard that races an event instead of waiting for
// it, green by luck and red by scheduling.
//
// Measured, in order of who measured it:
//   - docs/evidence/s1/136-ac15-rate-r1.md §4.2 (snapshot beac693): 4 reds in
//     8200 same-process shots of the half-covered leg = 0.049%, all four printed
//     the read-count guard ("only 2 reads taken", "only 3 reads taken"), and the
//     two shadow guards underneath it never got a turn. That is the signature of
//     a missing wait, not of a wrong product computation.
//   - this change, re-measured on the HEAD it lands on (ee5a25e): 0 reds in
//     2000 shots of that leg and 0 in 1800 shots of all six Group A legs, on a
//     quiet machine. The archived rate is small enough that a quiet 3800-shot
//     sample cannot be expected to reproduce it, so the shape is what is being
//     fixed here and the rate claim is the archived one; what this file does
//     prove is stated in docs/evidence/s1/observe-ac15-poll-r1.md.
//
// The discipline the fix is held to, same as AC#11's: the wait is bounded and
// monotonic, the threshold of every guard below is unchanged, nothing is
// skipped, and the criterion is the FIXTURE's own read counter wherever one
// exists - so a product-side breakage (a sweeper that stopped counting dropped
// reads, a disclosure that stopped naming sample_errors) can never satisfy the
// wait and turn a leg green by not testing the thing it exists to test.

// settleWindowWaitBound is the monotonic budget awaitWindow grants a leg before
// it reports its own precondition as unmet.
//
// It is derived from the durations this change measured on the HEAD it lands
// on, not from a handed-down number:
//
//	window                                   n      p50     p99     max
//	half-covered leg, same-process shots    2000     100ms   100ms   110ms
//	all six Group A legs, same-process      1800     100ms   200ms   200ms
//
// 2s is 10x the p99 of the slowest leg in the family (200ms) and about 20x the
// p99 of the leg that actually went red (100ms), so the budget buys roughly 20
// median-cost attempts and never expires while one ordinary window is still
// open. The only tail witness on record - the archived run's single 880ms pass
// (136-ac15-rate-r1.md §1.2, stale by this file's standards but the only one
// that exists) - still fits inside it twice over, so one stalled window cannot
// eat the budget by itself.
//
// What the bound has to outlast is a run of starved windows. At the archived
// whole-window shortfall rate (4/8200 = 0.049%) the chance of 20 independent
// attempts all falling short is 1e-56; at the pessimistic busy-machine stratum
// the same run reported (3/1200 = 0.25%) it is 1e-52. Both are far below any
// rate this package could ever confuse with "fixed". And expiry is a loud red
// naming the bound and every count seen, never a pass - so the residual risk
// this number carries is "a red on a machine stalled for two seconds straight",
// which is a red worth having.
const settleWindowWaitBound = 2 * time.Second

// settleWindow is one settle window: the report it handed back and how many
// reads its own seam took inside it. They are carried together on purpose -
// a leg that reopened its window must assert against the counts of the window
// it kept, not against a total accumulated over the attempts it threw away.
type settleWindow struct {
	rep   *SettleReport
	reads int
}

// stateWindow is the same pairing for a SampleState window.
type stateWindow struct {
	rep   *StateReport
	reads int
}

// awaitWindow reopens a window until short(w) reaches want, or the monotonic
// budget expires.
//
// The bound is observe.Timeout (clock.go:25-56), the instrument this package
// already uses for every bounded wait - goroutine.go:139 in production,
// goroutine_test.go:57-66 for AC#11's fix of this same defect. There is no
// wall-clock subtraction anywhere in here (AGENTS.md §1.2 bans a timeout built
// from wall-clock differences; tools/d22scan enforces it), no time.Sleep
// papering over a window, and no t.Skip: ticket 136 AC#15's judgement 3 allows
// only "waiting with a criterion", and its ⛔ list names the three cheap ways
// out that this shape closes off.
func awaitWindow[R any](t *testing.T, want int, open func() R, short func(R) int) R {
	t.Helper()
	tm := NewTimeout(settleWindowWaitBound)
	var seen []int
	for tries := 1; ; tries++ {
		w := open()
		n := short(w)
		if n >= want {
			return w
		}
		seen = append(seen, n)
		if tm.Expired() {
			t.Fatalf("precondition broken: %d windows opened inside the %v monotonic bound (clock.go Timeout) all fell short of %d; counts seen: %v",
				tries, settleWindowWaitBound, want, seen)
		}
	}
}

// awaitSettleReads opens CheckSettle windows until the seam has taken want
// reads, and hands back that window's report and count.
//
// open() is called again for every attempt, so each attempt is a complete
// window with a fresh fixture and a fresh report.
func awaitSettleReads(t *testing.T, want int, open func() (*SettleReport, int)) settleWindow {
	t.Helper()
	return awaitWindow(t, want,
		func() settleWindow {
			rep, reads := open()
			return settleWindow{rep: rep, reads: reads}
		},
		func(w settleWindow) int { return w.reads })
}

// awaitStateReads is the SampleState twin. want is a READ count, and a
// SampleState window takes 2 reads outside the ticker loop (sampler.go:269 and
// :306), so a leg that needs N samples asks for N+2 reads - spelled out at
// each call site rather than hidden here, because the arithmetic belongs to
// the window shape and not to the wait.
func awaitStateReads(t *testing.T, want int, open func() (*StateReport, int)) stateWindow {
	t.Helper()
	return awaitWindow(t, want,
		func() stateWindow {
			rep, reads := open()
			return stateWindow{rep: rep, reads: reads}
		},
		func(w stateWindow) int { return w.reads })
}
