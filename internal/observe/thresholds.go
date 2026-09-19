package observe

import "fmt"

// Frozen SLO thresholds (ticket 08; D32 16.3.2 as backfilled by docs/SLO.md
// plus the orchestrator ruling of 2026-09-19). DO NOT EDIT without an SLO.md
// change: loosening or skipping these gates is the D22 run-away mode 6, and
// the adversarial acceptance re-checks every row against the table.
//
// Ruling 1 (handles): D32 wrote "handles <300" pre-measurement; the measured
// layered-window + D2D + DWrite stack alone holds ~420 handles, so every
// state that carries the ball window stack gates at <600 (measured 420 +
// margin). The GDI <200 gate and the LEAK-TREND assertion are unchanged; the
// trend stays the primary verdict, the absolute handle limit is the guardrail.
//
// Work-peak memory (700MB) stays a TARGET, not an acceptance gate, until
// S3/S5 tree measurements backfill it (docs/SLO.md §2 keeps the distinction).
const (
	memCapSleeping     int64 = 25 << 20
	memCapArmed        int64 = 90 << 20
	memCapWarm         int64 = 350 << 20
	memCapConversation int64 = 350 << 20
	memCapPanel        int64 = 600 << 20
	memCapWorkPeak     int64 = 700 << 20 // TARGET (gate=false) until S3/S5

	cpuLimitSleeping     = 0.5 // % of all-core mean, 1min window
	cpuLimitArmed        = 2.0
	cpuLimitWarm         = 1.0 // no inference running
	cpuLimitConversation = 5.0 // VAD always-on
	cpuLimitPanel        = 10.0
	cpuLimitWorkPeak     = 30.0

	gdiLimitAll    = 200
	handleLimitAll = 600 // ruling 1, all six states carry the window stack

	goroutineLimitSleeping = 6 // D38b resident baseline
	goroutineLimitArmed    = 7 // + kws-infer

	// captureBufferLimitConversation is the D32 Conversation capture-buffer
	// bound; not measurable until the speech pipeline (tickets 15/26) runs,
	// so it is recorded as a note, never as a fake pass.
	captureBufferLimitConversation = 20 << 20
)

// stateMemCap returns the memory cap of a state (the settle-row target uses
// the same numbers).
func stateMemCap(st SLOState) int64 {
	switch st {
	case SLOArmed:
		return memCapArmed
	case SLOWarm:
		return memCapWarm
	case SLOConversation:
		return memCapConversation
	case SLOPanelOpen:
		return memCapPanel
	case SLOWorkPeak:
		return memCapWorkPeak
	default:
		return memCapSleeping
	}
}

func stateCPULimit(st SLOState) float64 {
	switch st {
	case SLOArmed:
		return cpuLimitArmed
	case SLOWarm:
		return cpuLimitWarm
	case SLOConversation:
		return cpuLimitConversation
	case SLOPanelOpen:
		return cpuLimitPanel
	case SLOWorkPeak:
		return cpuLimitWorkPeak
	default:
		return cpuLimitSleeping
	}
}

// buildVerdicts evaluates the frozen table for one state window. A Gate
// verdict that fails makes the report fail; non-gate rows are recorded
// evidence only.
func buildVerdicts(st SLOState, rep StateReport) []Verdict {
	vs := []Verdict{
		memVerdict(st, rep),
		cpuVerdict(st, rep),
		{
			Metric: "gdi_objects", Measured: fmt.Sprint(rep.GDIMax),
			Limit: fmt.Sprintf("<%d", gdiLimitAll),
			Pass:  rep.GDIMax < gdiLimitAll, Gate: true,
			Note: "D42#10; absolute guardrail, leak trend is the primary check",
		},
		{
			Metric: "user_objects", Measured: fmt.Sprint(rep.USERMax),
			Limit: "record", Pass: true, Gate: false,
			Note: "D42#10 sampling set; no frozen absolute limit",
		},
		{
			Metric: "handles", Measured: fmt.Sprint(rep.HandlesMax),
			Limit: fmt.Sprintf("<%d", handleLimitAll),
			Pass:  rep.HandlesMax < handleLimitAll, Gate: true,
			Note: "SLO.md ruling 1 (2026-09-19): window stack measured at ~420; <600; trend primary",
		},
		goroutineVerdict(st, rep),
		{
			Metric: "threads", Measured: fmt.Sprint(rep.ThreadsMax),
			Limit: "record", Pass: true, Gate: false,
			Note: "D42#10 sampling set; D2D/ORT thread pools vary by machine",
		},
	}
	vs = append(vs, sleepingConstraints(st, rep)...)
	return vs
}

func memVerdict(st SLOState, rep StateReport) Verdict {
	capBytes := stateMemCap(st)
	pass := rep.MemMedianBytes <= capBytes
	v := Verdict{
		Metric: "tree_private_bytes", Measured: fmt.Sprintf("%.1fMB", mb(rep.MemMedianBytes)),
		Limit: fmt.Sprintf("<=%.0fMB", mb(capBytes)),
		Pass:  pass, Gate: true,
	}
	if st == SLOWorkPeak {
		// docs/SLO.md §2: 700MB is a target until S3/S5 backfill - sampled
		// and recorded, never counted as an acceptance pass.
		v.Gate = false
		v.Note = "D32/SLO.md: target, not acceptance, until S3/S5 tree measurements"
	}
	return v
}

func cpuVerdict(st SLOState, rep StateReport) Verdict {
	limit := stateCPULimit(st)
	return Verdict{
		Metric: "cpu_percent_all_core", Measured: fmt.Sprintf("%.3f%%", rep.CPUMeanPercent),
		Limit: fmt.Sprintf("<=%.1f%%", limit),
		Pass:  rep.CPUMeanPercent <= limit, Gate: true,
	}
}

func goroutineVerdict(st SLOState, rep StateReport) Verdict {
	switch st {
	case SLOSleeping:
		return Verdict{
			Metric: "goroutines", Measured: fmt.Sprint(rep.GoroutinesMax),
			Limit: fmt.Sprintf("<=%d", goroutineLimitSleeping),
			Pass:  rep.GoroutinesMax <= goroutineLimitSleeping, Gate: true,
			Note: "D38b resident baseline",
		}
	case SLOArmed:
		return Verdict{
			Metric: "goroutines", Measured: fmt.Sprint(rep.GoroutinesMax),
			Limit: fmt.Sprintf("<=%d", goroutineLimitArmed),
			Pass:  rep.GoroutinesMax <= goroutineLimitArmed, Gate: true,
			Note: "D38b resident baseline + kws-infer",
		}
	default:
		return Verdict{
			Metric: "goroutines", Measured: fmt.Sprint(rep.GoroutinesMax),
			Limit: "record", Pass: true, Gate: false,
			Note: "per-task counts vary; D38b roster report governs",
		}
	}
}

// sleepingConstraints are the D32 Sleeping-row hard constraints: zero
// periodic disk writes and zero long-lived network connections, measured
// from process telemetry over the sampling window.
func sleepingConstraints(st SLOState, rep StateReport) []Verdict {
	if st != SLOSleeping {
		return nil
	}
	return []Verdict{
		{
			Metric: "disk_write_ops", Measured: fmt.Sprint(rep.WriteOpsTotal),
			Limit: "==0", Pass: rep.WriteOpsTotal == 0, Gate: true,
			Note: "tree write ops excluding the sampling process (observer cannot fail its own window); self delta recorded separately",
		},
		{
			Metric: "tcp_connections", Measured: fmt.Sprint(rep.TCPMax),
			Limit: "==0", Pass: rep.TCPMax == 0, Gate: true,
			Note: "zero long-lived network connections (D32 Sleeping row)",
		},
	}
}

func mb(b int64) float64 { return float64(b) / (1 << 20) }
