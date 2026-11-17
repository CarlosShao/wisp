# 49 — Session grants (D45-2): (tool, pattern, session) authorization + input.type binding

**Status:** ready-for-agent
**Claimed by:** agent-97 (2026-11-17, impl leg, commit 0172714d4 — 非实现者验收未做)
**Last update:** 2026-11-17
**Blocked by:** 04-sqlite-core, 48-approval-queue-full
**Parallel slots:** ≤2 sub-agents (A: grant store + evaluation in risk path; B: card third
option + revoke UX + audit logging)
**Spec refs:** SPEC-06 §8, D45-2, §15#7 input binding, C18/C20, S7

## What to build
Scoped session authorization: confirmation cards gain "本会话内允许 <tool> 于 <pattern>"
creating (tool, pattern, session_id) grants honored by the risk path, expiring at session end,
revocable live, fully audited — including the input.type target-process binding variant.

## Key constraints
- Grant = (tool, pattern, session_id) row (approval_grant table from 04); evaluation BEFORE
  re-prompting: matching grant → execute at L1-grade without new confirm; every grant use
  writes a tool_call audit row (grant_id linked).
- input.type variant: pattern = target process name; per-process isolation (Word ≠ 终端);
  "本会话内允许向所有应用注入" is danger-styled + expand-confirm (two clicks); swap of target
  app → re-confirm.
- Three inviolable limits (D45-2): L2 NEVER enters grants (an L2 decision can create only an
  L1-grade grant for future L1-class uses of the same tool+pattern — L2 always re-confirms);
  C25/R4 taint upgrades NOT overridden by any grant; 🔒 config-section loosening never via
  grants.
- Session end (31 dispose) → grants expired (rows retained 30d for audit); mid-session revoke
  (40 page + card) effective immediately on next evaluation.
- Confirm-card UX: three options (仅本次 / 本会话内允许此工具+模式【visually primary】 /
  全部应用 variant for input.type only, danger-styled expand).
- Batch interplay: aggregated L1 confirm may carry one grant creation for the whole batch
  pattern (if user opts in).

## Out of scope
- Cross-session persistent scopes (workspace/user — RESERVED, do not build).

## Acceptance criteria
- [x] Grant lifecycle: create → subsequent calls silent (audit rows written) → session end
      expires → new session re-prompts.
- [x] Revoke: immediate effect test (next gated call re-prompts).
- [x] Limits: L2 never auto-allowed by grant; tainted exfil upgrades despite grant; config
      loosening unaffected by grants (three explicit tests).
- [x] input.type: per-process isolation; all-apps option danger flow; target-swap re-confirm.
- [x] Batch + grant: 10-op batch with opt-in grant → one card → whole pattern granted.
- [x] Audit: grant_id linkage present in tool_call rows for every grant-backed execution.

## Progress log (append-only, newest last)

- 2026-11-17 agent-97（实现腿，commit `0172714d4`，基锚 `26ee1c081`）：六格实现全落（见
  docs/evidence/s1/49-session-grants-impl.md 的逐格读数）。要点：生命周期＝身份失效（224 已有）+
  收口 stamp（rt.close 的 SESSION-GRANT-SWEEP，行保留）；revoke＝Ledger.Revoke + console 动词
  `revoke-grant <编号>`，下一调用重弹；三限全钉（L2 有全覆盖 grant 源仍弹卡、R4 升级无视 grant、
  mode-silenced 永不记 allow_session_grant）；input.type＝解析进程名绑定（resolveInputTarget 与
  执行同管线）+ per-process 隔离 + all-apps 危险档（session-expand→session-all 两步）+ byte 界
  swap guard；batch＝PatternForPaths（最长公共祖先 + /*）+ session-pattern 一卡一行；
  audit＝tool_call.grant_id 逐调用落库（allow_session_grant 词表）。AC 框全部留给非实现者勾。
  自测：internal/session+tools+agent/approval+memory 全 ok；cmd/wisp 载具单跑五枚全 PASS
  （一次一发、跑前 wisp.test 现查、Git bash 门）。整包死法按 A1893/A1882 处置记录。
  未做（登记非省略）：面板侧 grants.list/revoke 入站方法（banned-root 表禁 Go 入站 grant 路由，
  出口今天为 console 动词）。零契约文字改动。

---

## 编排者收件（2026-11-17 09:3x，`0172714d4`＋`8a429abdb`，件 `docs/evidence/s1/49-session-grants-impl.md`）

**尺1** 两枚提交在（24 件）。**尺2** 件在。**尺3 我自己复跑**：`go test ./internal/session ./internal/tools ./internal/agent/approval -run TestTicket49` ⇒ **三包全 ok**（1.779s／0.163s／0.088s）；票内共 27 枚 `TestTicket49*`。
**六格最硬读数（腿交，我认可为行为型）**：①生命周期＝run 收口后 `approval_grant` 1 行、`revoked_at` 非空＋审计 `SESSION-GRANT-SWEEP revoked=1`；②revoke＝撤前 0 卡＋`grant_id=<行id>`、撤后 ≥1 卡＋`grant_id=NULL`；③limits＝**grant 源全覆盖（Covering/TargetCovering 全 true）下 L2 仍 1 卡、grant 源 0 次被查**；R4 全命中升级无视 grant；mode-silenced `decision=allow、grant_id=NULL` 永不 `allow_session_grant`；④input.type＝`RecordTarget("notepad.exe")` 后 0 卡、换目标重弹、swap guard 执行前重解析不符即拒（原文「目标应用已切换…重新确认」）、all-apps 两步仪式、未展开先答被拒；⑤batch＝10 路径 **1 张卡**→**恰 1 行 `<dir>/*`**；⑥audit SQL＝真库两笔被授权调用**两行 tool_call 各带 `grant_id=31`**，未授权行 NULL。
**边界（登记入票，非省略）**：面板侧 `grants.list`/`revoke` **入站方法未建** —— banned-root 表禁 Go 入站 grant 路由，操作者出口＝console 动词 `grants`/`revoke-grant <编号>`。⇒ **票 40 AC#2（grants 视图）因此判不了**，已在该票留白节写明。
⇒ 6/6，**本票结**；零契约改动、C17 白名单与面板入站路由零触碰。