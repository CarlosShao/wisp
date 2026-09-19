// Package statemachine owns the BallState state vocabulary and the D43
// transition table (SPEC-01 §3; C12 freezes the table in D43).
//
// Responsibilities:
//   - the 20 BallState values (see states.go)
//   - the authoritative D43 transition table (40 edges) and per-state timeouts
//   - "transitions not listed are illegal" enforcement (D22 gate 3)
//
// Non-responsibilities:
//   - no capabilities of its own: it never touches audio, LLM, tools or UI;
//     side effects listed in D43 guards are executed by callers
//   - no error-class knowledge beyond what observe defines (mapping lives in
//     observe per D37)
//
// Ticket 03 only pins the State vocabulary (states.go). The Idle-only
// placeholder is intentional: transitions and the per-state timeout table are
// implemented by ticket 07 (D43).
package statemachine
