// Package watchdog owns the resource watchdog (SPEC-01 §3): per-state lookup
// table of limits (D32 revised 3) and the enforcement actions compatible with
// each state (e.g. never unload KWS in Armed), plus config-file mtime polling
// reusing its loop (SPEC-03 §4.3).
//
// Responsibilities:
//   - sample the D32 SLO metrics: process-tree private memory via
//     proc.TreePrivateBytes (Job Object), CPU, plus GDI/User objects, handles,
//     goroutine count and thread count (D42#10)
//   - threshold breaches -> WatchdogAlert state / resource error class (D37),
//     never silent auto-unload of user-visible capabilities
//   - goroutine roster verification against the observe registry (leak check)
//
// Non-responsibilities:
//   - no sampling backend itself (observe provides SLO plumbing), no Job
//     Object ownership (proc), no state transitions (statemachine)
//
// DEFERRED(watchdog loop/thresholds): implemented by ticket 42. This ticket
// only freezes the package boundary.
package watchdog
