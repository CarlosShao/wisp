// Package statemachine owns the BallState state vocabulary and the D43
// transition table (SPEC-01 §3; C12 freezes the table in D43).
//
// Responsibilities:
//   - the 20 BallState values (see states.go)
//   - the authoritative D43 transition table (table.go: the 40 frozen rows
//     plus SPEC-08 §3's appended #41/#42) and the per-state timeout table
//     (timeouts.go)
//   - "transitions not listed are illegal" enforcement (D22 gate 3):
//     Machine.Dispatch rejects unknown (state, event) pairs loudly
//   - side effects fire as Effect events on an injectable Sink; the default
//     sink is an explicit no-op - execution belongs to later tickets
//
// Non-responsibilities:
//   - no capabilities of its own: it never touches audio, LLM, tools or UI
//   - no error-class knowledge beyond what observe defines (mapping lives in
//     observe per D37)
//
// Implemented by ticket 07; ticket 03 pinned only the State vocabulary.
package statemachine
