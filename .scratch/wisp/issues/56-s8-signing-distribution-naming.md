# 56 — S8: code signing, package-manager distribution, P10 residual naming — DEFERRED

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 54-s7-acceptance + SignPath qualification (P8) + user go-ahead
**Parallel slots:** ≤2 sub-agents
**Spec refs:** D17, D41, D33/F3, P8/P10, SPEC-11 §7, DEFERRED rows

## What to build
Zero-cash signing + trusted distribution: SignPath Foundation application, CA code signing wired
as secrets-gated CI step, winget/Scoop/Homebrew manifests, updater keying finished, and the P10
residual naming checks (npm/PyPI crates/domains/trademark) closed out.

## Key constraints
- Two signature mechanisms BOTH present (D41e): CA code signing (SmartScreen/Gatekeeper) +
  minisign (updates + models, already live from 14/56-prep) — neither substitutes the other.
- Secrets-gated: no secrets in CI → unsigned artifacts + docs, never a failed pipeline (D17).
- winget/Scoop manifests + SHA256; GitHub Releases layout per SPEC-11 §7.1; uninstaller size
  disclosures for models (D41c).
- P10 residual: npm/PyPI/crates name collisions, domain, trademark, 中文名「一缕」 — conflicts
  → fallback candidates Mote/唤/聆 ONLY via user decision (never invent names).
- Update channel e2e: signed update → staged → verified → swapped → rollback proven (54 drill
  made production-real).

## Acceptance criteria
- [ ] Signed Windows build passes SmartScreen (evidence) or documented fallback live.
- [ ] winget/Scoop install e2e on clean VMs.
- [ ] minisign chain re-verified end-to-end post-signing changes.
- [ ] P10 checklist filed; any conflict escalated to user before any rename.
- [ ] docs/BUILD.md + RELEASE.md updated to final state.

## Progress log (append-only, newest last)
