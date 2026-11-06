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
// Wiring owed by the ticket 12 list. Landed since that paragraph was written:
// approval.New is composed in cmd/wisp (run.go's assembleRuntime and
// resident_approval_windows.go's newResidentApproval) and handed to the tool
// bridge as tools.Options.Gate; Gate.AdmitTextTask has production call sites
// (run.go's mode-switch card and admitTask, the config-reload tick, the
// resident card path); the D31 ledger seam is wired (tools.Options.Cancel
// takes Gate.ToolsCancelBus, and the fs.write family fills
// Result.AppliedSteps); and the answer side is approval.Replies (replies.go),
// which both cmd/wisp surfaces bind - the console listener and the resident
// ball host - so Native() and Panel() are consumed, not just published.
//
// Landed with ticket 21 segment 2 (cardrender.go): approval.RenderCard is the
// ONE renderer both native surfaces print cards through - rules verbatim, the
// owning task, the L2 queue depth, the D45-1 batch view, and the FULL
// parameter object - and the resident orb posts the queue depth to
// ball.SetBadge on every card (SPEC-06 §7's 「深度角标：悬浮球必显」).
// Still owed:
//   - delete loop.decideRisk's declared-L1/L2 refusal branch: it is dead code
//     today (the loop passes a declared L1/L2 through whenever the host bound
//     Options.AdmitTask, and every real host does), so the deletion is
//     dead-code cleanup, not a live fix;
//   - the panel-side card half belongs to ticket 37 (PanelItem/PanelAPI are
//     its whole surface; the frontend page is outside this tree).
package approval
