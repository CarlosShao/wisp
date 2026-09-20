# 20 — host bridge + first tool family: fs.*, capability checks, spill rule

**Status:** in-progress (both fs segments landed; artifacts-spill rule separation, the
`[fs] allowed_dirs` first-use ask flow and the loop-side gate composition remain)
**Claimed by:** T20-seg1-agent (host bridge + C4 registry), T20-seg2-agent (fs write family)
**Last update:** 2026-09-20 (segment 2)
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
- [x] fs matrix: each action's RiskLevel asserted incl. trash=L1, overwrite=L2, cross-volume
      move→L2, out-of-allowlist read→L2. Names: `TestD34WriteMatrix` (11 subtests),
      `TestCrossVolumeRuleIsR8` (the hardware-free layer),
      `TestOverwriteDetectionFollowsTheCanonicalPath`.
- [x] Atomic write test: kill mid-write → no partial file at target (temp+rename verified).
      Name: `TestAtomicWriteKillsMidWrite` (overwrite branch: old bytes survive byte-for-byte;
      new-file branch: nothing appears at the path; both assert no staging file is left behind).
- [x] Cancel-semantics test: cancelled after N ops → "applied steps" list correct.
      Names: `TestVetoAfterWorkStartedProducesAppliedSteps`,
      `TestLateVetoRendersTheApprovalLayersAppliedStepsReport`,
      `TestCrossVolumeMoveStopsWithTwoCopiesOnLateStop`,
      `TestApprovedL2ThatStopsMidWriteIsNotRenderedAsNeverStarted`.
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
- [2026-09-20T12:10Z] agent=T20-seg2 did=**segment 2 (the fs write family)** — `fs.write`
  (D31 `stageAndRename`: the staging file is created in the TARGET's directory so the final
  `os.Rename` never crosses a volume, and the target is touched by exactly one syscall, so a death
  before it cannot produce a partial file), `fs.trash`, `fs.move` (L1 same-volume rename / L2
  cross-volume or occupied destination), `fs.delete` (registered ONLY when
  `FSDeps.DeleteEnabled`, re-checked inside Execute so a roster that outlived the setting still
  refuses). The R8 escalations reach the FROZEN rules through a new host-side `Decl.Facts` hook:
  the bridge keeps owning `Declared` and `Paths` (a tool can neither lower its own floor nor hide a
  target), the hook contributes the irreversibility classes, and a panicking hook contributes an
  unknown class which R8 fail-closes to L2. Every existence probe canonicalizes through C26 first.
  **trash is the real shell API, and how that is proven:** `SHFileOperationW` with
  `FOF_ALLOWUNDO` + `FOF_SILENT|NOCONFIRMATION|NOERRORUI|NOCANCELLATION`, issued on a thread
  pinned with `runtime.LockOSThread` + `CoInitializeEx` (the COM apartment is per OS thread;
  without it shell32 deadlocks — measured, see below), and success is then CONDITIONED on finding
  the `$I*` restore record shell32 wrote for that exact path under `<drive>\$Recycle.Bin\<SID>`
  (`TestShellTrashLeavesARecycleBinEntryTheOSCanEnumerate` re-derives it independently of the
  implementation's own check). An unlink cannot produce that record, and `os.Remove` appears in
  this package only in `fs.delete` and in the cross-volume move's source removal, both of which
  declare themselves permanent (R8 → L2). Two findings the next agent should not rediscover the
  hard way: (1) **`SHQueryRecycleBinW` blocks forever from a Go thread** — the first version used
  it for the item-count proof and the test binary sat in that syscall until the go-test timeout
  killed it, i.e. a tool call that can hang a task in silence; (2) **`FOF_WANTNUKEWARNING` asks
  shell32 for a MODAL dialog**, which is the same hang wearing a face on the user's desktop. Both
  dropped in favour of the restore-record check. The trash tests DO place one small temp item in
  the real Recycle Bin each run; that is the price of an honest assertion.
  **Two D31 bugs the composition exposed (both in the fake-clean-cancel direction, both fixed
  here):** `approval.Gate.PendingApproval` never recorded a handoff, so an approved L2 write that
  then stopped mid-flight was rendered 「本次调用未进入执行阶段」 over a list of steps that HAD
  landed; and `CancellationReport.TextFor` had no branch for "steps landed, no veto on record"
  (task cancellation / tool fault), so it fell through to the 「未上报其步骤」 sentence. The
  `markStarted` key is the INCOMING correlation id, because the queue re-stamps
  `d.CorrelationID` with its own item id — the card's address, not the call's. The only
  applied-steps shape in the repo is still `approval.CancellationReport`: `tools.CancelBus.Report`
  returns an `fmt.Stringer` precisely so no second shape could be born on the way.
  **AC#2/3/4 GREEN (test names in the AC text). AC#5 and AC#6 NOT ticked:** neither is in this
  segment's scope — AC#5's junction-denial coverage lives in `internal/risk`
  (`pathresolver_junction_windows_test.go`), and AC#6's artifacts-spill rule-separation test is
  still unwritten.
  **THE R12 ANSWER, restated because three ACs this week were left open on exactly this question:
  who calls this in production? NOBODY YET, and it is NOT reachable from `cmd/wisp`.** Re-run the
  判据: `grep -rln '"github.com/CarlosShao/wisp/internal/tools"' --include=*.go cmd/ internal/ |
  grep -v _test.go` → only the four `internal/agent/approval` files;
  `grep -rn '"github.com/CarlosShao/wisp/internal/agent"' --include=*.go cmd/ internal/ | grep -v
  _test.go` → only `internal/tools/bridge.go`; and
  `grep -rn "tools.New\|approval.New\|agent.New(" --include=*.go cmd/` → **zero hits**. Ticket 12
  owns that wiring; no AC was ticked on the strength of a call path that does not exist.
  **`decideRisk` WAS NOT TOUCHED** (`git diff HEAD -- internal/agent/loop.go` empty; the guard sits
  at `internal/agent/loop.go:719-738`). Deliberate, per the hazard register: it refuses any call
  whose roster entry declares L1/L2 before Execute runs, so an `fs.write` dispatched by a Loop is
  still refused rather than gated. Replacing it means handing the loop the bridge's verdict, and
  that replacement cannot be composed honestly inside this segment: `cmd/wisp` builds no loop and
  no gate, so the only test proving "guard gone, gate in its place" would have to author ticket
  12's production wiring. What IS proven here is that nothing else is missing —
  `internal/tools/wiring_test.go` runs the real `approval.Gate` over this bridge: an L1 write
  blocks its full 2s window and then writes, an in-window veto writes nothing, a late veto renders
  approval's own report wording, and an L0 read opens no card.
  **Gates (scoped to ./internal/tools/ + ./internal/agent/..., all run):** `gofmt -l` clean on
  every file this agent touched (`internal/agent/loop_golden_test.go` and `truncation_test.go`
  carry pre-existing drift from d6f6406, untouched here); `go vet` CLEAN;
  `go test -count=2` ok tools 14.386s / agent 2.463s / approval 0.145s;
  `go test -race -count=1` ok 9.982s / 2.801s / 1.131s; `tools/d22scan` → **clean** (segment 1's
  pre-existing mockllm finding is gone). `internal/ball` + `cmd/balldebug` untouched, no GUI
  process launched, no `-tags winlive` run, frozen files (`docs/PLAN.md`, `docs/specs/*`,
  `internal/risk/*`) zero-diff. Commits 0986d63 + 94827e4 + e669567.
  **MUTATION EVIDENCE (so AC#3/AC#4 are not read as self-certifying):** with
  `os.CreateTemp` swapped for a direct `os.OpenFile(target, O_CREATE|O_WRONLY|O_TRUNC)` the two
  kill tests go red (`D31 violated: after a mid-write kill the target holds N bytes…` and
  `AC#3 violated: a partial file exists at the target`); with `fofALLOWUNDO` zeroed out (a
  permanent shell delete, i.e. the unlink-shaped lie) both trash tests go red with
  `回收站中没有出现任何新的还原记录（C:\$Recycle.Bin）：拒绝按已回收处理`. Both mutations were reverted;
  `git status` is clean and every gate above was re-run against the committed tree.
  next=**ticket 12 (composition) or ticket 21 segment 2**: (a) build the Bridge + `approval.Gate`
  in `cmd/wisp`, and in the SAME commit replace `decideRisk`'s roster-floor refusal with the
  bridge's verdict while running the loop with a nil Journal so `tool_call` rows do not double;
  (b) AC#5 bridge-level junction/reparse denial test; (c) AC#6 artifacts-spill rule-separation
  test (`internal/agent/spill.go` is the host-internal precedent; D22 ban 7 forbids naming a gated
  tool after it); (d) the `[fs] allowed_dirs` first-use ask flow (minimal native prompt now, full
  GUI at 39) — note it has NO SPEC-12 §5 row, so marking it `DEFERRED` in code would break the
  1:1 cross-check; (e) still owed to ticket 19 (N-11) the `rules_gateway.go` R4 comment, and the
  missing `tool_call.rules_hit` column (a SPEC-02 §3 contract change, human approval required).
