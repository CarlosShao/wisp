// Package agent owns the ReAct agent loop (SPEC-01 §3): the LLM/tool cycle,
// context budgeting, cancellation, and truncated tool-call handling (D21:
// max_tokens stop reason fails all unclosed tool calls).
//
// Responsibilities:
//   - ReAct loop driving llm adapters and tools, one root ctx per task (D38c)
//   - context budget, truncation, retry/backoff policy consumption (D37)
//   - failure classification into observe error classes
//
// Non-responsibilities:
//   - no provider protocol code (llm), no tool implementations (tools),
//     no approval decisions themselves (agent/approval), no UI
//
// DEFERRED(loop): implemented by ticket 10. This ticket only freezes the
// package boundary and the two sub-packages below.
package agent
