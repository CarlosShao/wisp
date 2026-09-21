# 37 — Approval UI: queue view, L2 card panel-grade, native fallback card verification

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 21-approval-gates-minimal, 35-panel-bridge-c17
**Parallel slots:** ≤2 sub-agents (A: queue/task list UI + depth badge sync; B: panel L2 card +
native fallback polish + batch view)
**Spec refs:** SPEC-08 §5.4 (approval rows), §17.5, C18, D31 visibility mandates, D45-1 batch UI,
C27 fallback

## What to build
The approval experience at panel grade: queue head + depth indicator, per-request full context
(task, tool, params, reason), batch confirmation view, and the verified native fallback card
path when the panel is unavailable — with the ball depth-badge as the always-visible channel.

## Key constraints
- Visibility mandates (D31): queue depth ALWAYS on ball badge (sync <100ms); queued/blocked
  tasks visible in panel task list WITH "waiting on what"; every request self-describes
  (task, tool, FULL params, risk, reason) — never a bare "需要确认".
- Panel L2 card: top 2px danger bar, tool + L2 pill, complete params, rules_hit reasons,
  batch expandable list (first 5 + expand), persistent "点击悬浮球以批准" guide + pulsing ball
  glyph, buttons = reject + view-full-params ONLY (no allow — F2/§15#6; redundant assert here).
- Native fallback (panel unavailable): 320×140 native card (allow + reject + open-panel ghost
  link) — verify C27 rule "L2 capability never disappears"; visuals per §17.5 native-row.
- Timeout UX: 270s warning state on card (30s before 300s auto-reject); post-reject one-click
  replay button (C18).
- Queue FIFO order display; correlationId shown for forensics; decisions (allow/reject/timeout)
  land in tool_call rows with grant linkage hooks.
- Batch aggregated confirmation renders total + affected dirs + expandable detail; L2 batches
  never aggregate (guard re-asserted at UI level).

## Out of scope
- Session-grant options in card (49 adds the third option); multi-task scheduler itself (47).

## Acceptance criteria
- [ ] E2E: two sequential L2 requests → depth badge 1→2 (simulated pre-47 via queue test mode),
      head-only display, decision routing by correlationId.
- [ ] Timeout: card warning at 270s → auto-reject at 300s → replay works; ball badge clears.
- [ ] Fallback: with panel runtime masked, L2 still approvable via native card (33 fixture).
- [ ] No-allow invariant: panel card markup contains no allow action wired to decision API
      (DOM + bridge double test).
- [ ] Batch view: 10-op batch renders aggregated card; 50+ upgrades to L2 card (R7).
- [ ] Human visual sign-off of card + fallback (screenshots vs design/screens/approval.html).

## Progress log (append-only, newest last)

- 2026-09-21 15:4x（编排者，**接线时必须做的一次键对账 = 票 87 的残留 R-3**）：
  `acceptor-ticket87` 实测（其报告 §6.2 攻击 2）：**卡片 B 的"入键"与卡片 A 的队列签发键撞车时，
  宿主按入键取消 B 会拒到 A**。方向是拒绝 ⇒ **不触"查不到就放行"那条底线**，而且 `reject()` 这条路线
  **先于票 87 存在**；票 87 新增的只是"`Veto` 这条以前什么都不拒的路线现在也会拒错"。
  ⇒ **判据（票 37 接线时算你的验收项）**：`internal/panel/approval.go` 带来的 C17 入键（`CorrelationID`
  取自 `ApprovalSubject.CorrelationID`，见 `:28/:40/:78`）与队列自己的签发键（`queue.go:133-139` 的 `it.Corr`）
  必须做**一次显式对账**，并且要有一条用例：**两张卡共存时，取消其中一张绝不能命中另一张**
  （今天两者相同只因为 `internal/agent/loop.go:644` 把 corr 写成 taskID 且队列没重签发——**那是巧合，不是设计**）。
  相关：票 87 已 `-done`（库层就绪、用户层未接线的判词在 `docs/evidence/s1/87-adversarial-acceptance.md` §12），
  R-1/R-2 在**票 97**。
