package observe

import "testing"

// TestFrozenCapAccessors pins the watchdog-facing accessors to the frozen
// thresholds.go values (ticket 42). The watchdog gates with the same digits
// the sampler reports against; any edit to the frozen numbers surfaces here
// as a conscious, human-approved SLO.md change, never as a silent drift.
func TestFrozenCapAccessors(t *testing.T) {
	caps := map[SLOState]int64{
		SLOSleeping:     25 << 20,
		SLOArmed:        90 << 20,
		SLOWarm:         350 << 20,
		SLOConversation: 350 << 20,
		SLOPanelOpen:    600 << 20,
		SLOWorkPeak:     700 << 20, // target, not acceptance (SLO.md §2)
	}
	for _, st := range SLOStates {
		if got := MemCapBytes(st); got != caps[st] {
			t.Errorf("MemCapBytes(%s) = %d, want %d (frozen D32 row)", st, got, caps[st])
		}
	}

	cpu := map[SLOState]float64{
		SLOSleeping:     0.5,
		SLOArmed:        2.0,
		SLOWarm:         1.0,
		SLOConversation: 5.0,
		SLOPanelOpen:    10.0,
		SLOWorkPeak:     30.0,
	}
	for _, st := range SLOStates {
		if got := CPULimitPercent(st); got != cpu[st] {
			t.Errorf("CPULimitPercent(%s) = %g, want %g (frozen D32 row)", st, got, cpu[st])
		}
	}

	if got := GDILimit(); got != 200 {
		t.Errorf("GDILimit() = %d, want 200 (frozen guardrail)", got)
	}
	if got := HandleLimit(); got != 600 {
		t.Errorf("HandleLimit() = %d, want 600 (SLO.md ruling 1)", got)
	}

	limit, gated := GoroutineLimit(SLOSleeping)
	if !gated || limit != 6 {
		t.Errorf("GoroutineLimit(Sleeping) = %d/%v, want 6/true (D38b resident baseline)", limit, gated)
	}
	limit, gated = GoroutineLimit(SLOArmed)
	if !gated || limit != 7 {
		t.Errorf("GoroutineLimit(Armed) = %d/%v, want 7/true (D38b + kws-infer)", limit, gated)
	}
	for _, st := range []SLOState{SLOWarm, SLOConversation, SLOPanelOpen, SLOWorkPeak} {
		if _, gated := GoroutineLimit(st); gated {
			t.Errorf("GoroutineLimit(%s) claims a frozen cap; the frozen table records it without one", st)
		}
	}
}
