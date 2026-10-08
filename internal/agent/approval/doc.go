// Package approval owns the ApprovalQueue (C18) for L2 (native-side confirmed,
// F2) operations (SPEC-01 §3). Ownership assignment per D31/D37: approval UI
// lives in agent (queue) + panel (rendering).
//
// Responsibilities:
//   - queue of pending L2 grants, 300s timeout with a visible 30s warning,
//   - decision callbacks from the native side / panel
//
// Non-responsibilities:
//   - no risk assessment (risk decides L0/L1/L2), no approval UI rendering
//     (panel), no L1 countdown confirm (statemachine Confirming state)
//
// DEFERRED(queue): minimal gates landed in ticket 21 (this package), the
// multi-task routing half belongs to ticket 48 -> SPEC-12 §5 has no row for it;
// the row that does exist and still binds this package is
// 「DEFERRED | 快捷键路径语音否决（B1）」, which is why the veto word is
// reported as 「语音取消不可用」 until ticket 41 loads KWS.
//
// What ticket 21 segment 1 actually put here (decision + queue layer, UI
// injected): Gate (implements tools.Gate), the 2-3s L1 block window with four
// honestly-reported veto channels, the C18 single-task queue with a 300s
// auto-reject and a 270s warning, D45-1 batch aggregation, D47's task
// admission guard, D31's applied-steps report + Bus, and the native grant that
// makes a panel-sourced allow unrepresentable rather than merely disallowed.
//
// Wiring owed by the ticket 12 list. Landed since this paragraph was written:
// approval.New is composed in cmd/wisp (run.go:612, resident_approval_windows.go:370)
// and handed to the tool bridge as tools.Options.Gate (run.go:748); Gate.AdmitTextTask
// has production call sites (run.go's mode-switch card and admitTask, the config-reload
// tick, the resident card path); and the D31 ledger seam is wired (tools.Options.Cancel
// takes Gate.ToolsCancelBus, and the fs.write family fills Result.AppliedSteps).
// Still owed:
//   - hand Queue.Native() to the ball click / native card / hotkey handlers
//     and Queue.Panel() to the ticket 37 panel bridge;
//   - delete loop.decideRisk's declared-L1/L2 refusal: the branch still exists
//     (loop.go:794), but `wisp run` binds AdmitTask (run.go:1008), so it is no
//     longer reached; the deletion is dead-code cleanup now, not a live fix.
package approval
