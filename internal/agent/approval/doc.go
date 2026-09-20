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
// Wiring still owed by ticket 12 (this package has NO production caller yet,
// so nothing here runs in cmd/wisp today):
//   - compose approval.New(...) as tools.Options.Gate (replacing NoGate);
//   - call Gate.AdmitTextTask(taskID) once per TEXT-loop task and defer its
//     revoke - unadmitted tasks are refused by design (D47), so forgetting it
//     fails closed, it does not fall open;
//   - hand Queue.Native() to the ball click / native card / hotkey handlers
//     and Queue.Panel() to the ticket 37 panel bridge;
//   - have the fs.write family (ticket 20 segment 2) poll Gate.Bus() and fill
//     Result.AppliedSteps;
//   - delete loop.decideRisk's declared-L1/L2 refusal, which currently kills
//     fs.write before the bridge ever assesses it.
package approval
