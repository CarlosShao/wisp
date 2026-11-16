# 34 — Frontend scaffold: React+TS+Tailwind+shadcn, tokens port, docker build, embed pipeline

**Status:** **superseded（2026-11-16 编排者裁）** —— 与票 77 是同一件事（前端脚手架），77 是活的那张（3 格已交）⇒ 本票**不再派**，余下 5 格已**具名并入票 77 的 AC#8–AC#12**（见 77 票尾）。⛔ 悄悄消失：接手名字＝票 77。
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

---

## 编排者裁定（2026-11-16 18:4x）：本票并入票 77，不再派

理由：34 与 77 都写「前端脚手架（React+TS+Tailwind+shadcn）」，77 已有 3 格交付并是活票；两张并行会双写 rontend/（写面互斥）。
**具名接手**（⛔ 不许写"暂缓"了事）：
| 34 的格 | 并入 |
|---|---|
| docker build→dist→go build 在真面板窗服务脚手架 | 77 AC#8 |
| CI job: tsc/lint/build 绿＋种子违规 | 77 AC#9 |
| token parity（design/assets/tokens.css ↔ 前端 tokens 脚本 diff） | 77 AC#10 |
| 深/浅色渲染＋对比度 ≥AA | 77 AC#11 |
| bundle 预算 initial JS ≤400KB gz | 77 AC#12 |