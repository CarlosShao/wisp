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
//   - the fs family, per D34: fs.read/fs.list (L0 pair, ticket 20 segment 1)
//     and fs.write/fs.trash/fs.move (segment 2) are here, with fs.delete
//     registered ONLY when [fs] delete_enabled=true. web/system/doc/search
//     tools are 22-24. The write half brings the D31 temp+atomic-rename
//     writer, the applied-steps ledger (Result.AppliedSteps + CancelBus) and
//     the real Shell recycle-bin call behind fs.trash.
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
// THE BLOCKER ticket 12/21 must clear before fs.write is reachable (verified by
// reading the code, not by trusting this paragraph): internal/agent/loop.go
// decideRisk (internal/agent/loop.go:719-738 in the file as it stands today)
// refuses any call whose DIRECTORY ENTRY declares L1/L2 before Execute runs -
// it can only see the declared floor, so it cannot tell an L1 write from an L2
// overwrite. That guard was the stopgap while no gate existed; the gate now
// exists (internal/agent/approval) and is proven end-to-end against this
// bridge by wiring_test.go, but nothing in cmd/wisp constructs it, so deleting
// the guard would drop the only thing standing between a declared-L1 tool and
// an ungated write. It must be replaced IN THE SAME COMMIT that injects the
// gate, by handing the loop the bridge's verdict instead of the roster's floor
// (single-writer contract: the bridge owns tool_call, the loop runs with a nil
// Journal, or the rows double).
//
// DEFERRED(tools): implemented from ticket 20 (host bridge + fs tools) onward
// (22/23/24 add web/system/doc tools). This ticket only freezes the package
// boundary.
package tools
