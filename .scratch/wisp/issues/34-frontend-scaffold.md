# 34 — Frontend scaffold: React+TS+Tailwind+shadcn, tokens port, docker build, embed pipeline

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 01-build-chain (assets/web embed target exists)
**Parallel slots:** ≤2 sub-agents (A: scaffold + toolchain + build pipeline; B: token port +
base components + icon set)
**Spec refs:** SPEC-08 §5.3, §6, §17.3/17.4/17.6, D29, C21, SPEC-11 §3.2

## What to build
The `frontend/` workspace producing `assets/web/dist` via the Node docker builder: React+TS+
Tailwind with shadcn/ui source-copied components, the C21 token port (`tokens.css` → CSS vars),
base component set, the Lucide icon subset (zero emoji), and the CI checks (no hardcoded hex
outside tokens, zero emoji scan, tsc/lint/build).

## Key constraints
- Stack exactly per D29/17.6: React+TS, Tailwind, shadcn/ui COPIED source (not npm dep), Radix
  primitives, CVA+clsx+tailwind-merge, cmdk, Sonner, motion (panel-only, NO infinite loops),
  react-markdown+remark-gfm, rehype-sanitize (whitelist config!), Shiki, @tanstack/react-virtual,
  Zod (input-feedback only; Go structs own defaults), date-fns (display only).
- Banned: emoji libs, Font Awesome, Material/MUI/antd, Lottie/Rive, styled-components, chart
  libs (17.6 exclusion list — CI scans).
- Tokens: `frontend/src/tokens.css` derived from `design/assets/tokens.css` (single source
  referenced, v2.3 cold ambient + glass system); every screen uses vars only — CI grep
  `#[0-9A-Fa-f]{6}` outside tokens files must be zero.
- Icons: `lucide-react` subset (~55 names per §17.4 list), stroke 1.5, sizes 14/16/18/20 only;
  zero emoji including text symbols (U+2190–U+2BFF etc.) — CI scan.
- Type scale/spacing/radius/motion per tokens (30/18/15/14/12.5/11 px tiers; 4px grid; 120/180/
  260ms; cubic-bezier(0.32,0.72,0,1)); Chinese rules: body ≥14px, no weight <400.
- Build: `docker/frontend.Dockerfile` (node:20-alpine, npm ci, build, export dist) →
  `assets/web/dist`; host needs no Node. CI: tsc+lint+build job.
- App shell: router skeleton with placeholder routes matching the 11 design screens; dark
  default + light toggle; 640px panel width; keyboard-first focus states.
- Frontend stays STATELESS (no cross-show caching) — scaffold enforces a single store hydrated
  by `panel.resync` only (35).

## Out of scope
- Bridge implementation (35); actual page features (36–40); ball visuals (native side).

## Acceptance criteria
- [ ] `docker build` → dist → `go build` serves the scaffold in the real panel window (33).
- [ ] CI job: tsc, lint, build green; seeded violations (hardcoded hex, an emoji char, a banned
      dep import) each fail.
- [ ] Token parity check: scripted diff of token names design/assets/tokens.css ↔
      frontend/src/tokens.css (no drift).
- [ ] Dark+light render of the shell + one sample screen; contrast spot-check ≥AA on both.
- [ ] Bundle budget: initial JS ≤400KB gz (React cost acceptable per D29; motion/Shiki lazy).

## Progress log (append-only, newest last)
