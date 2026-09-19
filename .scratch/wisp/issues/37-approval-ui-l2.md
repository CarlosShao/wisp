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
