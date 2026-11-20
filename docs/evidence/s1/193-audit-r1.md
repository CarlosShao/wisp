# 193-audit-r1 — 对抗验收（裁决者≠实现者），对象 `f34210795` ＋ `acce6ea13` · 2026-11-20

**腿**：只读审计（审计路由）。**总判：6 节全部「成立」**（0 夸大 / 0 不成立）。
本件是编排者对腿报告的摘要；逐条红句原文与变异细节见腿报告落盘件（`D:\tmp\wisp\20261120\aud193\`）。

## 1 牙是真的吗 —— 成立（四发全红，还原逐字节 SHA256 13/13 全等）
| 牙 | 变异 | rc | 红句要点 |
|---|---|---|---|
| 基线 | 无 | 0 | `render-animation-193: OK, ticket 193 AC#4+AC#5 (off 真停: getAnimations()=0 + zero active rAF handles, …)` |
| 1 | 删 `panel-visibility.ts:102` 的 `if (!motionHeavyAllowedOn(doc)) return;` | 1 | `RED, 3 failed check(s)`：loop factory must gate start() on the tier / a start() under off must register NOTHING |
| 2 | 把 `[data-motion="off"]` 规则挪进 `@media (prefers-reduced-motion: reduce)`（塞进既有块＝3 failed；另起自己的 media 块＝1 failed，与证据件口径逐字对齐） | 1 | `the off rules must sit at the stylesheet's top level … [brace depth = 1]` |
| 3 | `motion.ts` 未知分支改成 `applyMotionTier("off")` | 1 | `RED, 11 failed`：an absent section, an empty tier or an unknown word must withdraw the attribute |
| 4 | 删 `motion.ts` 翻档时的 `for (const a of live) a.cancel();` | 1 | `the flip to off must cancel what is already in flight` / `getAnimations() must be EMPTY after the flip` |
**还原**：先备份 13 文件记 SHA256，逐发还原后复算 **MISMATCH_COUNT=0**；还原后仪器 rc=0、`go test -count=1 ./internal/config/ ./internal/panel/` → ok 2.017s / ok 18.350s。

## 2 有没有放宽既有断言 —— 成立（**0 条**）
- `git diff acce6ea13~1 acce6ea13 -- internal/panel/composer_test.go` 的删除行**只有** `"PanelSnapshot": {"panelAnimation"},` ＋ 其上注释 ⇒ 删的是**台账登记**不是断言。
- `git diff f34210795~1 f34210795 -- "*.go"` 删除行**只有 7 行**，全是 `composerTSNameMap` 的 gofmt 对齐，零语义。
- 台账收口真是"两向比较"：三发探针全红 —— C（只回退 `panel.ts`）⇒ `interface PanelSnapshot does not declare it - register it in composerTSGapLedger ONLY with a judgement, never to silence it`；D（只删 `composerTSNameMap` 那行）⇒ `composerTSNameMap has no page-side interface for it - add the translation and the panel.ts interface in the same move`；E（Go tag 改名）⇒ **双向两行**同时红。**注：这三条探针是"能力在"的正控（探针未改 revert 前全部跑在同一棵未变异树上，已逐发还原；每条红句都指名了要发现的东西）**。

## 3 tag 逐字一致 —— 成立
`panel.ts` 的 `PanelAnimationSection`/`panelAnimation` 与 `internal/panel/composer.go json:"panelAnimation,omitempty"`、`panel_animation_193.go json:"tier"` 逐字对齐（腿报告给对照表）。

## 4 默认值确实是"待 owner 拍" —— 成立
schema 默认词、设置页文案「（默认，待定）」、证据件三处口径一致；无一处写成定案。

## 5 诚实边界 —— 成立
`panelAnimation` 在 `cmd/wisp/*.go` 无发射点（泵源未接），`internal/panel/pump.go` 无假 off 兜底 ⇒ 今天"不发键"而不是"发 off"。

## 6 零 emoji / 零 dist —— 成立
`git show --stat acce6ea13` 无 `dist` 路径；新增文件 emoji 扫描干净。

## 结论
票 193 已勾的 AC#4/AC#5 由**非实现者**复验通过；AC#3（泵源＋写名册＋ApplySetting 接线）仍开，接线腿排队中。
