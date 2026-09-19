# 55 — S8: macOS port (second-class platform) — DEFERRED until open-source prep

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 54-s7-acceptance + self-use week + explicit user go-ahead (S8 starts only
pre-open-source, D23)
**Parallel slots:** ≤2 sub-agents
**Spec refs:** D7, D23, DEFERRED registry macOS row, SPEC-12 §5

## What to build
The macOS platform layer: borderless NSWindow + `level = .floating` ball, Keychain-backed
SecretStore (S8 row of C28), symlink-realpath path rules (18 stub), and CI cross-builds with
prebuilt artifacts.

## Key constraints
- D7 stance unchanged: macOS ships as "second-class, unverified on real hardware" unless a
  maintainer with hardware joins; README must say so.
- Ball parity: 20 states via NSWindow; layer-shell Linux still RESERVED (do not start).
- CI: darwin/amd64 + arm64 artifacts (cgo strategy per BUILD.md extension); notarization
  DEFERRED (Gatekeeper docs `xattr` workaround, D17).
- All DEFERRED rows' completion criteria from SPEC-12 §5 must be met before S8 closes.

## Acceptance criteria
- [ ] Both ball implementations pass interaction tests; CI produces signed-optional macOS
      artifacts; CER-equivalent English models out of scope (i18n is 58).
- [ ] Keychain SecretStore passes 06's test matrix ported.
- [ ] PathResolver macOS branch (realpath + lstat) passes red-team port.
- [ ] README platform matrix updated honestly.

## Progress log (append-only, newest last)
