# 20 — host bridge + first tool family: fs.*, capability checks, spill rule

**Status:** in-progress (segment 1 of 2 landed; fs.write/trash/move + artifacts spill remain)
**Claimed by:** T20-seg1-agent (host bridge + C4 registry)
**Last update:** 2026-09-20
**Blocked by:** 17-risk-assessor-c19, 18-path-resolver-c26
**Parallel slots:** ≤2 sub-agents (A: host bridge + ToolProvider registry; B: fs tool family +
artifacts spill)
**Spec refs:** SPEC-07 §2–§3, D3, D14, D34, C1/C3/C4, D34 note②, D31 atomic-rename

## What to build
The `tools` module host bridge — the single choke point every capability flows through — with
C1 Tool contract, C3 capability enforcement (11 caps), C4 ToolProvider registry (Builtin now;
Manifest/Goja slots), and the first builtin family: fs.read/list/write/trash/move (+delete
gated). Includes the host-internal artifacts spill path and `list_tools`-resident registration
from 10.

## Key constraints
- Host bridge: every execution passes `risk.Assess` + capability check (undeclared capability →
  hard reject, not error), decision + rules_hit recorded into `tool_call` rows with
  correlation_id; plugin code can never bypass (no direct OS access from plugin runtimes —
  structurally guaranteed later by 51).
- C1 Tool: Name/Description/Parameters(JSON Schema)/Execute(ctx, params, onUpdate). Tool concurrency
  ≤4 per task; per-tool timeoutMs honored via ctx.
- fs family RiskLevels (D34): read/list L0 (in-allowlist; out-of-scope → L2 via R2);
  write-new L1; overwrite L2 (R8); move L1 (cross-volume → L2); **trash L1** (shell API —
  recycle bin); delete **not registered** unless `[fs] delete_enabled=true` → L2.
- All writes: temp file + atomic rename (D31 cancel semantics); cancel after partial work →
  report "steps possibly already applied" list (never fake clean cancel).
- Paths ALWAYS via PathResolver (18); results taint-marked (19).
- Host-internal artifact writes (`artifacts\`): separate internal API, NOT a Tool, no gating,
  data-dir-only, quota-managed (04), visible/deletable later via panel. Distinguishing rule
  documented + tested.
- `[fs] allowed_dirs` first-use ask flow: minimal native prompt now (full GUI at 39).

## Out of scope
- Approval UI/queue mechanics (21 wires the L1/L2 UX); web/system/doc/search tools (22–24);
  plugin manifests (50).

## Acceptance criteria
- [x] Bridge tests: undeclared capability → reject; risk decision → gate branch (L0 pass, L1
      pending-window callback, L2 pending-approval callback); tool_call rows complete.
- [ ] fs matrix: each action's RiskLevel asserted incl. trash=L1, overwrite=L2, cross-volume
      move→L2, out-of-allowlist read→L2.
- [ ] Atomic write test: kill mid-write → no partial file at target (temp+rename verified).
- [ ] Cancel-semantics test: cancelled after N ops → "applied steps" list correct.
- [ ] Reparse/short-name target via junction → denied (18 integration).
- [ ] Artifacts API writes only under data dir; attempting user-dir write via artifacts API →
      rejected; `fs.write` to user dir remains gated (rule separation test).

## Progress log (append-only, newest last)
- [2026-09-20T09:55Z] agent=T20-seg1 did=**segment 1 (scope items 1-6)** — `internal/tools/` now holds
  the host bridge + the C1/C3/C4 layer. `bridge.go` = `agent.ToolProvider` (the loop's existing
  Tools()/Execute() seam, `internal/agent/loop.go:563,696`; `var _ agent.ToolProvider = (*Bridge)(nil)`
  proves it), so no parallel dispatch path exists. Per call: C3 caps → one `risk.Assess` over
  C26-canonical paths + C25 taint bound per task scope → L0 pass / L1 `Gate.PendingWindow` /
  L2 `Gate.PendingApproval` / Deny refuse-without-asking → D38d ceiling (clamped ≤4, enforced in
  the bridge so a host-internal caller cannot bypass it) → C22 per-tool `context.WithTimeout` →
  C25 `Mark` on the result → complete `tool_call` row (assessed risk_level + decision + outcome +
  error_class + correlation_id). A rejection is NEVER a Go error (SPEC-07 §2 未声明即拒绝调用（不是报错，
  是拒绝）: `loop.go:654-657` turns a provider error into class internal, which is exactly the
  "Wisp broke" misreading C3 forbids) — the tests assert `err == nil` AND the tool's run counter
  stayed 0. `paths.go` closed a real gap: **no type in this repo implemented the frozen C19
  `PathCanonicalizer`/`SensitiveClassifier` seams before this**, so tickets 17/18 left R2/R3
  DORMANT in every composition; wiring them is what makes out-of-allowlist read → L2 a verdict
  (`TestDeclaredRiskIsOnlyAFloor`). `registry.go` declares all four C4 slots and answers
  `ErrSlotNotLanded` for manifest/goja/mcp (`TestC4SlotsWithoutAnImplementationAreNotSilent`)
  rather than pretending. Registered **fs.read + fs.list only** (`TestFSRegistrationIsTheL0Pair`
  pins the pair at 2 entries).
  AC#1 GREEN, test names: `TestUndeclaredCapabilityHardRejects`,
  `TestRiskDecisionRoutesEachLevelToItsBranch` (7 subtests: L0-passes-no-gate / L1-window /
  L1-timeout-executes / L1-veto / L2-approval / L2-timeout-auto-rejects / NoGate-fails-closed),
  `TestToolCallRowsAreComplete` (real SQLite, 5 rows: allow+success / L2 reject+user_rejected /
  unknown-tool / bad-args / L1 veto), plus `TestCapabilitySetIsFrozen` (11 tokens),
  `TestToolConcurrencyCeilingIsFour`, `TestPerToolTimeoutHonoredViaContext`, `TestC1ToolContract`,
  `TestC4RegistryRejectsBadDeclarations`, `TestFSReadTaintsIt`→`TestFSReadTaintFeedsR4`
  (R4 + SessionOverrideBlocked through the real C25 engine), `TestSensitiveFileIsDeniedNotEscalated`
  (R3 tier A never reaches a gate), `TestEmptyAllowlistAuthorizesNothing`.
  **Gates (scoped to ./internal/tools/, all run):** `gofmt -l` empty; `go vet` CLEAN;
  `go test -count=2` ok 0.480s; `go test -race -count=1` ok 1.394s; `tools/d22scan -root .`
  → 1 finding, `internal/llm/adaptertest/mockllm.go:68` bare-goroutine, **pre-existing at
  78b1466 (ticket 11), zero findings in internal/tools**; `internal/ball`+`cmd/balldebug`
  (live ticket-62 agent) untouched, frozen files zero-diff. Commits c9c3a6f + 64c5fea.
  **THE COMPOSITION ANSWER (rulings A8/A11/R12 shape — read this before believing "green"):**
  `who constructs the bridge?` — **nobody in production yet.** No non-test code calls `tools.New`
  or `agent.New`; the **only non-test importer of `internal/agent` in the entire repo is
  `internal/tools/bridge.go` itself** (re-run the判据:
  `grep -rn '"github.com/CarlosShao/wisp/internal/agent"' --include=*.go cmd/ internal/ | grep -v _test.go`),
  and `agent/loop.go:190` still defaults a nil provider to `EchoProvider`;
  `cmd/wisp` constructs no loop at all. **Ticket 12 wires it** (R12 retitled 12 the 接线票),
  which is why `internal/tools/doc.go` states this out loud instead of leaving an unwired-but-green
  module. AC#1 was ticked because its text is a *bridge-test* criterion and does not depend on the
  wiring; no other AC was touched, and none of the remaining five can be closed by segment 1 alone.
  next=**segment 2**: `fs.write` (D31 temp+atomic rename + `TestAtomicWriteKillsMidWrite`),
  `fs.trash` (shell API, L1), `fs.move` (L1 / cross-volume L2), `fs.delete` behind
  `[fs] delete_enabled`, applied-steps report on cancel, the artifacts-spill rule-separation test
  (AC#5/AC#6; `internal/agent/spill.go` is the host-internal precedent — the D22 ban-7 regex
  `"spill"`/`"internal.x"` forbids naming a gated tool after it, keep it that way), and the
  `[fs] allowed_dirs` first-use ask flow. **Two hazards segment 2/21 must not trip:**
  (a) `loop.decideRisk` (`loop.go:719-738`) rejects anything declaring L1/L2 BEFORE Execute runs,
  so `fs.write` will be dead on arrival under the loop until ticket 21 replaces that function with
  this bridge's verdict — register it only with the single-writer contract (bridge owns
  `tool_call`, loop runs with a nil Journal) or the row doubles;
  (b) the missing `tool_call.rules_hit` column is a SPEC-02 §3 contract change, not a segment-2
  task — rules_hit currently travels in `Decision`/`OnDecision` + the audit line only.
  Also still owed to ticket 19 (N-11): `rules_gateway.go`'s R4 comment says "Dormant until ticket 19"
  while the loop wiring is what landed here; that file is FROZEN for this agent, so it was NOT
  edited — ticket 21 owns that comment.
