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
// DEFERRED(queue): minimal gates land in ticket 21, full queue in ticket 48.
// This ticket only freezes the package boundary.
package approval
