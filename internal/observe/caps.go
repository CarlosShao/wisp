package observe

// Frozen-threshold read accessors for the resource watchdog (ticket 42).
//
// The numbers live in thresholds.go and are FROZEN (ticket 08; D32 16.3.2
// as revised by docs/SLO.md and the orchestrator ruling of 2026-09-19).
// The watchdog must gate against the same digits the sampler gates with,
// so it reads them through these accessors instead of keeping a second
// copy of the table: one table, two consumers - referencing the frozen
// constants is the sanctioned shape, copying them into a second package is
// not (the duplicate-truth shape the 188-task-state ruling named).
//
// These are read-only views and nothing else. No accessor may loosen, skip
// or re-derive a frozen value (D22 run-away mode 6); changing a number is
// an SLO.md change with an SLO.md ruling behind it, never a watchdog edit.

// MemCapBytes returns the frozen private-working-set cap of a D32 state.
// WorkPeak returns its target value, which stays a recorded target, not an
// acceptance gate (thresholds.go / docs/SLO.md §2).
func MemCapBytes(st SLOState) int64 { return stateMemCap(st) }

// CPULimitPercent returns the frozen all-core CPU mean limit of a state.
// The watchdog records it with its tick metrics; unload passes are driven
// by the memory sustain rule (D32 16.3.2 revision 3), not by CPU.
func CPULimitPercent(st SLOState) float64 { return stateCPULimit(st) }

// GDILimit returns the frozen GDI object guardrail (all states).
func GDILimit() int { return gdiLimitAll }

// HandleLimit returns the frozen handle guardrail (all states; docs/SLO.md
// ruling 1: the layered-window + D2D + DWrite stack measured ~420 handles,
// gated at <600 with the leak trend as the primary verdict).
func HandleLimit() int { return handleLimitAll }

// GoroutineLimit returns the frozen per-state goroutine cap and whether a
// frozen cap exists at all. Only Sleeping (6) and Armed (7) carry frozen
// caps in the table; every other state is a record-only row, where the
// second return is false and the caller records without gating.
func GoroutineLimit(st SLOState) (limit int, gated bool) {
	switch st {
	case SLOSleeping:
		return goroutineLimitSleeping, true
	case SLOArmed:
		return goroutineLimitArmed, true
	default:
		return 0, false
	}
}
