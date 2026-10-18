package panel

// Ticket 328 AC#2's shape readings on the join side: the four correlation shapes
// production actually mints must LIGHT the roster row of the task they were built
// on, and the two shapes that carry no task id must stay dark.
//
// WHY THE BYTES BELOW ARE WRITTEN HERE AT ALL. The shape's provenance is not this
// file: internal/agent/corr_roster_join_328_test.go drives a REAL loop
// (run -> executeCalls -> loop.go:676, the only call site of callCorr) and reads
// the correlation id off the ToolRequest the loop dispatched, then carries it into
// this package through the production card constructor. callCorr is not exported,
// so no test in this package can mint one; what this package can and does measure
// is which shapes the roster accepts once those bytes exist. Each row therefore
// names the mint point that produces it, so a reader can check the pair:
//
//   "child-1#call_e1"    loop.go:607 - the provider gave the call an id (arm A)
//   "child-1#call-0"     loop.go:605 - the call came with no id (arm B, index 0)
//   "child-1#call-0#7"   queue.go:168 pushed arm B's corr a second time and
//                        rewrote it as "<original>#<seq>" (stacked rewrite)
//   "child-1#9"          queue.go:168 rewriting a corr that was the bare task id
//   "approval-7"         queue.go:165 - the empty-correlation fallback, no '#'
//   "host:mode-switch"   cmd/wisp's host-initiated cards (gate.go:281 orDefaultText)
//
// The last two are the shapes the ruling says must stay unlit: neither contains a
// rostered task id followed by '#', so no row may light for them and no special
// case is opened for either.

import "testing"

// TestProductionCorrelationShapesLightTheirRows is the four-shape reading. It uses
// the rows every ticket-220 case shares (root-1 / child-1) so the widened arm is
// measured against the same roster the guard rails are.
func TestProductionCorrelationShapesLightTheirRows(t *testing.T) {
	cases := []struct {
		name string
		corr string
		want bool
	}{
		{"arm A provider call id lights its task", "child-1#call_e1", true},
		{"arm B id-less call lights its task", "child-1#call-0", true},
		{"arm B with the queue clash rewrite stacked lights its task", "child-1#call-0#7", true},
		{"the queue rewrite of a bare task id lights its task", "child-1#9", true},
		{"arm A rewrite of a bare task id lights its task", "child-1#call_e1#7", true},
		{"the empty-corr fallback names no row", "approval-7", false},
		{"a host-initiated correlation names no row", "host:mode-switch", false},
		{"a correlation built on a task nobody filed names no row", "ghost-9#call-0", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data := blockedPacket(t, []NativeVerdict{card197(c.corr)}, nil)
			assertBlocked(t, data, "child-1", c.want)
			// One card never blocks the whole roster: the other row must stay false
			// whatever the widened arm did to the row the correlation names.
			assertBlocked(t, data, "root-1", false)
		})
	}
}

// TestEachCorrelationLightsOnlyTheRowItIsBuiltOn is the per-row half of the join:
// the widened arm must move exactly the row whose id the correlation was built on,
// not every row that merely resembles it. This is the shape 328's boundary exists
// to keep - two tasks, one card, one row lit.
func TestEachCorrelationLightsOnlyTheRowItIsBuiltOn(t *testing.T) {
	data := blockedPacket(t, []NativeVerdict{
		card197("root-1#call-0"),
		card197("child-1#call_e1"),
	}, nil)
	assertBlocked(t, data, "root-1", true)
	assertBlocked(t, data, "child-1", true)

	// Now the asymmetric world: one card, built on root-1. child-1 shares the
	// prefix-free shape and must not be reached through it.
	data = blockedPacket(t, []NativeVerdict{card197("root-1#call-0")}, nil)
	assertBlocked(t, data, "root-1", true)
	assertBlocked(t, data, "child-1", false)
}

// TestAnL1WindowCorrelationIsMatchedByTheSameArm records that the widened arm is
// not card-only: an L1 window's correlation is filed by the same
// indexWaitingKey, so a window the host reports as "<taskID>#<callID>" blocks that
// row too - and a window with an id nobody filed still blocks nothing.
func TestAnL1WindowCorrelationIsMatchedByTheSameArm(t *testing.T) {
	data := blockedPacket(t, nil, []L1WindowWait{{CorrelationID: "child-1#call-0", TaskID: ""}})
	assertBlocked(t, data, "child-1", true)
	assertBlocked(t, data, "root-1", false)

	data = blockedPacket(t, nil, []L1WindowWait{{CorrelationID: "child-1x", TaskID: ""}})
	assertBlocked(t, data, "child-1", false)
}
