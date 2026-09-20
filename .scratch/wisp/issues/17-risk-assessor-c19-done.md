risk module: C19 RiskAssessor, R1–R9 rule set, fail-closed fusion (DONE ✅)

**Status:** done
**Claimed by:** orchestrator -> sub-agent T17-impl
**Last update:** 2026-09-20T00:32:00Z
**Blocked by:** 03-skeleton-runtime-rules
**Parallel slots:** ≤2 sub-agents (A: assessor core + fusion; B: R2–R9 rule implementations +
test matrix)
**Spec refs:** SPEC-06 §2–§3, D4, D31②, C19, D39 C19 table, 16.9#5

## What to build
`risk` module's central risk assessor: the ONLY authority for risk decisions. Built-in rule set
R1–R9 is the decider; plugin/manifest declarations are inputs (R1 = lower bound, never
conclusion); multi-assessor fusion takes max severity; any assessor panic/absence → fail-closed
L2 (R9).

## Key constraints
- Rules exactly per SPEC-06 §3 table: R1 declared level (input-only); R2 path-in-allowlist
  (needs 18's PathResolver — stub interface now, integrate later); R3 sensitive-path A/B
  (interface now, real lists in 18); R4 Provenance taint-hit (interface now, real in 19);
  R5 network target (allowlist/private-ranges/protocols); R6 shell argv meta-chars; R7 batch
  scale ≥50 files → L2; R8 irreversibility (delete/overwrite/power/close-window/send-class);
  R9 fail-closed L2.
- Fusion = max severity; single rejection beats any allow. Rule-set changes = contract change
  (human approval) — encode rule IDs as data, stable and reviewable.
- Output type: `{level, rules_hit[], reason}` — the confirmation card (21/37) renders rules_hit
  verbatim ("R4: 包含来自 web.fetch 的内容").
- Pluggable assessor interface (RESERVED guardian scoring attaches here later; do NOT implement).
- Agent-loop integration point defined: ToolProvider calls `risk.Assess(tool, params, context)`
  before execution; loop consumes decision (loop-side wiring lands 20/21).
- Unit tests: per rule one positive + one negative + panic-injection (R9) + fusion max.

## Out of scope
- PathResolver implementation (18); taint engine (19); approval UI/queue (21); plugin manifests.

## Acceptance criteria
- [ ] Rule matrix: each R1–R9 positive/negative/edge case passes (table-driven).
- [ ] Panic in any sub-assessor → R9 L2 fail-closed (no decision path returns L0/L1 on error).
- [ ] Fusion: conflicting severities → max wins; send-class via R8 overrides a declared L0.
- [ ] Decision object includes rules_hit + human-readable reason (golden snapshots).
- [ ] Contract-freeze note: rule IDs + semantics documented as frozen; doc committed.

## Progress log (append-only, newest last)
- [2026-09-20T00:10:30Z] agent=orchestrator claimed=T17-impl did=dispatched (50-min deadline window; hard stop 08:40 local, clean-unit boundary only) next=work
- [2026-09-20T00:32:00Z] agent=T17-impl did=R1-R9 implemented (assessor.go contract+fusion+R9 recover; rules_gateway R1-R4 interface-injected dormant-until-wired; rules_network R5 full; rules_shell R6 full; rules_scale R7 full; rules_irreversible R8 full; 14 tests PASS: per-rule pos/neg/edge, panic-injection R9 x2, fusion max, Deny-beats, R4 session-override-block, 10 golden snapshots) next=integrate 18 resolver/classifier + 19 taint via With* setters; note: T18 WIP landing in parallel (pathresolver files untouched per boundary)
- [2026-09-20T00:35:00Z] agent=T17-impl did=commit 67ffbd8 pushed to origin+cnb (TLS retry 2 attempts); acceptance criteria all met incl. contract-freeze block in assessor.go; did=handoff-to-orchestrator
- [2026-09-20T00:27:30Z] agent=orchestrator did=adversarial PASS (orchestrator-executed; reports docs/evidence/s1/.scratch/wisp/issues/17-*.md; full-package tests + race green after both agents merged) next=ticket DONE
