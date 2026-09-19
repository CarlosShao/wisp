// Package risk owns the security gate (SPEC-01 §3): RiskAssessor R1-R9 (C19),
// PathResolver (C26), Provenance/taint (C25), and the A/B tool blacklists.
//
// Responsibilities:
//   - classify every tool call into L0 (direct) / L1 (quick confirm) /
//     L2 (ApprovalQueue) / deny; declared RiskLevel is a lower bound, never
//     the conclusion (R1)
//   - path canonicalization robust against junction/8.3/UNC/\\?\ bypasses
//   - injection-suspect detection produces class "injection" errors (D37)
//
// Non-responsibilities:
//   - no tool execution (tools), no approval UI or queue (agent/approval),
//   - no capability declaration storage (plugin manifests feed it inputs)
//
// DEFERRED(R1-R9/PathResolver/Provenance): implemented by tickets 17, 18 and
// 19. This ticket only freezes the package boundary.
package risk
