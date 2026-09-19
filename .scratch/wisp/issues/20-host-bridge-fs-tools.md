# 20 — host bridge + first tool family: fs.*, capability checks, spill rule

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
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
- [ ] Bridge tests: undeclared capability → reject; risk decision → gate branch (L0 pass, L1
      pending-window callback, L2 pending-approval callback); tool_call rows complete.
- [ ] fs matrix: each action's RiskLevel asserted incl. trash=L1, overwrite=L2, cross-volume
      move→L2, out-of-allowlist read→L2.
- [ ] Atomic write test: kill mid-write → no partial file at target (temp+rename verified).
- [ ] Cancel-semantics test: cancelled after N ops → "applied steps" list correct.
- [ ] Reparse/short-name target via junction → denied (18 integration).
- [ ] Artifacts API writes only under data dir; attempting user-dir write via artifacts API →
      rejected; `fs.write` to user dir remains gated (rule separation test).

## Progress log (append-only, newest last)
