// Package tools owns the ToolProvider surface (C4) and the host bridge that
// funnels every tool call through capability and risk enforcement (SPEC-01 §3).
//
// Responsibilities:
//   - tool registry and schemas (D34 list lands slice by slice)
//   - host bridge: the ONLY path from agent/LLM to effects; it consults risk
//     (C19), path resolution (C26) and provenance (C25) before executing
//   - tool concurrency ceiling 4 (D38d)
//   - tool errors: internal failures go back to the LLM as class "tool"
//     (self-correction), never surface as user-facing failures (D37)
//
// Non-responsibilities:
//   - no risk rules (risk), no approval queue (agent/approval), no plugin
//     sandbox (plugin)
//
// DEFERRED(tools): implemented from ticket 20 (host bridge + fs tools) onward
// (22/23/24 add web/system/doc tools). This ticket only freezes the package
// boundary.
package tools
