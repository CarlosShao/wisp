# 07 — Ball shell + state machine core (20 states, 40 transitions, hotkeys, tray) (DONE ✅)

**Status:** done
**Claimed by:** orchestrator -> sub-agent T07-impl
**Last update:** 2026-09-19T13:56:25ZT10:41:26Z
**Blocked by:** 03-skeleton-runtime-rules (soft: 02 baseline ④ for draw-path memory budget —
build against the budget; verify in 12)
**Parallel slots:** ≤2 sub-agents (A: Win32 layered window + D2D renderer + tokens; B:
statemachine table + hotkeys + tray + focus rules)
**Spec refs:** SPEC-08 §2–§3, D43, §14.6, D29, C12, C21, B1 Esc rule

## What to build
The native floating ball: Win32 layered window rendering via Direct2D/DirectWrite with C21
DesignTokens, plus the `statemachine` module implementing the 20-state / 40-transition table
(D43) as data (table-driven, exhaustive, illegal transitions rejected), global hotkeys, tray
icon+menu, focus rules, click-through, and per-state visuals sufficient for S1 (static states +
basic animations; full visual polish gate is human acceptance at 12).

## Key constraints
- Window: `WS_EX_LAYERED|WS_EX_NOACTIVATE|WS_EX_TOOLWINDOW|WS_EX_TOPMOST`; `UpdateLayeredWindow`
  with 32bpp premultiplied ARGB; D2D + DirectWrite (never GDI text); shared `ui-sta` STA thread
  with future panel (D38a — one D2D factory only).
- Animation discipline: `Sleeping` fully static (no timer); animated states only
  Listening/Thinking/Acting/Speaking/Confirming; ≤30fps cap; loops only in non-idle states
  (`Warm` = ≥2s opacity breathing via compositor-friendly path).
- All colors/radii/shadows/easing from C21 tokens (`design/assets/tokens.css` is the reference
  implementation) — hardcoded hex in ball code = review-rejected. Implement the native-side token
  table + a doc listing token↔value pairs for cross-checking.
- State machine: all 20 states and 40 transitions from SPEC-08 §3 as an exhaustive table; each
  state's timeout per table (Listening 15s first-round / 90s in-session / 30s Conversation;
  Confirming countdown; Warm 90s; Settling 3s; Error 10s auto...). Unknown (from,event) pairs →
  rejected loudly (undefined = stop).
- Transition side-effects hooks fired as events (SessionScope creation etc. are no-ops here,
  consumed by later tickets).
- Hotkeys: summon/mute/cancel/panel from `[hotkey]`; `Esc` takeover semantics reserved — during
  Confirming, cancel key temporarily takes Esc and MUST be returned after session (B1).
- Tray: left-click = open panel (no-op stub), right-click menu (open panel / mute / pause wake /
  exit). Click-through transparent regions; ball body clickable; never steals focus.
- Multi-monitor: Per-Monitor V2 DPI; position saved per monitor; off-screen → back to primary **Spike backfill (T02): Path Y confirmed (idle 16.5/16.6MB ≤ 25MB with cgo resident); ball
  window stack (layered+D2D+DWrite) ≈420 handles → SLO handle gate for window-bearing states is
  <600 (orchestrator ruling, docs/SLO.md); release path must bound non-Sleeping handle growth.**
  visible (DPI resource rebuild stubbed to full for S1).

## Out of scope
- KWS/Armed loading (41), approval depth badge visuals (37), panel WebView (33), Conversation
  privacy confirm dialog content (28 wires the transition), watchdog interplay (42).

## Acceptance criteria
- [x] Exhaustive transition test: every legal row fires (state+side-effect spy), every illegal
      pair rejected; per-state timeout table-driven tests.
- [ ] Visual: 20 states rendered in a debug cycle page/window; human screenshot review vs
      design/screens/ball.html (colors/opacity/sizes per SPEC-08 §2.1).
- [ ] Interactive: click ball → Sleeping→Listening; Esc/click cancels from Confirming; focus
      never stolen (verified with a focused editor); transparent region click-through.
- [ ] Sleeping-state CPU ≈0 (no timers) — measurable in 08's sampler; here: assert no animation
      timer handles alive in Sleeping.
- [ ] Hotkeys registered/re-registered on config change; mute toggles Muted state.
- [ ] Multi-monitor: drag to second monitor, persist, restore; simulated detach → primary.
- note: AC#2 left open — docs/evidence/s1/07-adversarial-acceptance.md row 3 rules 视觉证据
      "PASS（人工签收挂起）"; the AC's second clause (human review vs ball.html) is still the open
      待人项 H1 in docs/reports/pending-and-issues.md, so no completion ruling exists
- note: AC#3 left open — no ruling found in docs/evidence/s1/07-adversarial-acceptance.md (5-row
      table searched for click/Esc/focus/click-through; row 2 only says "ball ok / statemachine ok"
      without naming what those tests assert). Candidate evidence unadjudicated at
      internal/ball/tokens_test.go:239 TestHitTestAndDPIInjection
- note: AC#4 left open — no ruling found in docs/evidence/s1/07-adversarial-acceptance.md (searched
      for Sleeping/CPU/timer/animation); the AC defers measurement to ticket 08's sampler, whose own
      report (docs/evidence/s1/08-adversarial-acceptance.md row 2) rules the Sleeping state gates
      PASS at 8.7MB but says nothing about animation timer handles. Candidate evidence unadjudicated
      at internal/ball/tokens_test.go:202 TestAnimationPolicyZeroTimerInSleeping
- note: AC#5 left open — no ruling found in docs/evidence/s1/07-adversarial-acceptance.md (searched
      for 热键/hotkey/re-register/mute); row 5 only lists "热键默认值待 39 接线" as a registered
      deviation, not a pass on re-registration on config change. Candidate evidence unadjudicated at
      internal/ball/hotkey_test.go:51 TestDefaultHotkeys
- note: AC#6 left open — no ruling found in docs/evidence/s1/07-adversarial-acceptance.md (searched
      for 多显示器/DPI/monitor/drag/persist/restore). Candidate evidence unadjudicated at
      internal/ball/live_windows_test.go:125 TestBallLivePositionPersistence

## Progress log (append-only, newest last)
- [2026-09-19T10:41:26Z] agent=orchestrator claimed=T07-impl did=dispatched (maintain 2-way concurrency floor) next=sub-agent works through acceptance criteria
- [2026-09-19T11:20:00Z] agent=T07-impl did=statemachine-core next=ball-window-unit
- [2026-09-19T11:35:00Z] agent=T07-impl did=c21-native-tokens next=ball-window-renderer
- [2026-09-19T12:40:00Z] agent=T07-impl did=ball-window-renderer-live next=live-tests+evidence
- [2026-09-19T13:00:00Z] agent=T07-impl did=hotkeys-tray-position-live-tests next=visual-evidence
- [2026-09-19T13:45:00Z] agent=T07-impl did=visual-evidence next=final-gates
- [2026-09-19T13:55:00Z] agent=T07-impl did=handoff-to-orchestrator status=implementation-complete
- [2026-09-19T13:56:25Z] agent=orchestrator did=T07-adv PASS (orchestrator-executed; D43 #1-#40+#41/#42 table audit; report docs/evidence/s1/07-adversarial-acceptance.md; human visual sign-off pending -> docs/reports) next=ticket DONE
- [2026-09-20T02:25Z] agent=agent-bookkeeping-1 did=AC boxes reconciled against docs/evidence/s1/07-adversarial-acceptance.md (5-row audit-item table, not a per-AC table): 1 checked (AC#1 via rows 1+2), 5 left open (reasons above) next=none
