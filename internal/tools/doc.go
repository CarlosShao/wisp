// Package tools owns the ToolProvider surface (C4) and the host bridge that
// funnels every capability through enforcement (SPEC-01 §3, SPEC-07 §2).
//
// Responsibilities:
//   - the host bridge (bridge.go): the ONLY path from agent/LLM to effects. Per
//     call it runs the C3 capability check, one risk.Assess verdict over
//     C26-canonical paths and C25 taint, the L0/L1/L2/Deny routing, the D38d
//     concurrency ceiling and the C22 per-tool deadline, then books the
//     tool_call row. It implements agent.ToolProvider, so the loop's existing
//     dispatch is this code - there is no second path to forget to secure.
//   - the C1 Tool contract (tool.go) and the C4 registry (registry.go)
//   - the fs family, landing slice by slice per D34: fs.read/fs.list are here
//     (ticket 20 segment 1); fs.write/trash/move and the delete_enabled-gated
//     fs.delete are segment 2; web/system/doc/search tools are 22-24
//   - tool errors: internal failures go back to the LLM as class "tool"
//     (self-correction), never surface as user-facing failures (D37)
//
// Non-responsibilities:
//   - no risk rules (risk), no approval mechanics (agent/approval, ticket 21),
//     no plugin sandbox (plugin), no artifact spill (agent/spill.go: that is a
//     host-internal write and D34 note 2 forbids modelling it as a tool)
//
// Composition status: nothing in cmd/wisp builds a Bridge yet. The wiring is
// ticket 12's (see docs/reports/pending-and-issues.md R12, which turns 12 into
// the wiring ticket after rulings A8/A11 found capabilities proven only in
// test harnesses). Until then the loop runs on agent.EchoProvider
// (internal/agent/loop.go:190).
//
// DEFERRED(tools): implemented from ticket 20 (host bridge + fs tools) onward
// (22/23/24 add web/system/doc tools). This ticket only freezes the package
// boundary.
package tools
