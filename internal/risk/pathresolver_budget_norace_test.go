//go:build !race

package risk

import (
	"testing"
	"time"
)

// resolveBudget is ticket 18 AC#6: one Resolve call must stay at or under 1ms
// on a warm handle cache.
const resolveBudget = time.Millisecond

// TestResolvePerCallBudget pins the C26 performance contract. It reuses the
// benchmark machinery (fixed b.N scaling, no hand-rolled timing loops) so the
// number it asserts is the same number `go test -bench` reports.
//
// Not compiled under -race: the race runtime inflates per-call cost by an order
// of magnitude, which would make a 1ms wall budget meaningless. Run the budget
// check as its own race-free pass (as scripts/ci does).
func TestResolvePerCallBudget(t *testing.T) {
	target := seedResolveTarget(t)
	warmResolve(t, target)

	res := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := Resolve(target, nil); err != nil {
				b.Fatalf("Resolve(%s): %v", target, err)
			}
		}
	})

	ns := res.NsPerOp()
	t.Logf("C26 Resolve: %d ns/op = %.3f ms/op (budget %.3f ms, %d samples)",
		ns, float64(ns)/float64(time.Millisecond), float64(resolveBudget)/float64(time.Millisecond), res.N)
	if time.Duration(ns) > resolveBudget {
		t.Fatalf("C26 budget breach: Resolve averages %v per call, budget %v", time.Duration(ns), resolveBudget)
	}
	// 50 calls/task must fit the same envelope the budget was written for.
	if cost := time.Duration(ns) * 50; cost > 50*resolveBudget {
		t.Fatalf("C26 50-call task envelope exceeded: %v", cost)
	}
}
