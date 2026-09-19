# 57 — S8: plugin SDK, developer docs, example plugins, community registry — DEFERRED

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 54-s7-acceptance + 53 lark-cli example + user go-ahead
**Parallel slots:** ≤2 sub-agents
**Spec refs:** D23, D19, §14.7, DEFERRED registry rows, C14/C16/C24

## What to build
The third-party developer path: SDK docs for Tier1 (incl. command-kind) and Tier2 plugin
authoring, contract-test kit, curated examples, and a minimal searchable registry with
version-compatibility checks.

## Key constraints
- Docs teach the frozen contracts (C14 schema, C24 surface, risk-input-only rule, net
  allowlists) — no API beyond contracts; "risk declarations are inputs, R1 only" stated
  prominently.
- Test kit: manifest validator + capability/negative-case harness a plugin author can run
  locally (reuses CI fixtures).
- Examples: one Tier1 HTTP plugin, one command-kind plugin (non-credential), one Tier2 state
  plugin; all pass the kit + contract tests.
- Registry (minimal): static index (JSON) + install-by-URL with the SAME install authorization
  flow (no bypass); compatibility check via host_api.
- Security stance from D19 intact: no forced signing, no central gatekeeping.

## Acceptance criteria
- [ ] Third-party author (roleplay agent without repo context) ships a Tier1 + Tier2 plugin
      using only docs + kit — acceptance bar from DEFERRED row.
- [ ] Registry install e2e incl. malicious-manifest rejection cases.
- [ ] Docs reviewed against contracts (no drift); examples pass CI.

## Progress log (append-only, newest last)
