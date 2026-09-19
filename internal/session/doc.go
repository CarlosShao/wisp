// Package session owns SessionScope (C31) (SPEC-01 §3): the wake -> explicit
// end / 90s-idle lifecycle, Warm/Conversation timing, model reference
// counting, and D45 authorization registration.
//
// Responsibilities:
//   - build the session-level DisposalScope (C11 hierarchy: session > task >
//     tool > plugin); session end MUST fully dispose, including
//     debug.FreeOSMemory (C11 step 4) so the D32 10s fall-back stays reachable
//   - Warm keepalive (B4) and the 90s -> Settling transition trigger
//   - per-session grants ledger (D45-2)
//
// Non-responsibilities:
//   - no state transitions (statemachine), no model loading (speech/models),
//   - no audio device ownership (audio)
//
// DEFERRED(SessionScope): implemented by ticket 28. This ticket only freezes
// the package boundary.
package session
