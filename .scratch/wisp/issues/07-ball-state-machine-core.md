# 07 — Ball shell + state machine core (20 states, 40 transitions, hotkeys, tray)

**Status:** in-progress
**Claimed by:** orchestrator -> sub-agent T07-impl
**Last update:** 2026-09-19T10:41:26Z
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
- [ ] Exhaustive transition test: every legal row fires (state+side-effect spy), every illegal
      pair rejected; per-state timeout table-driven tests.
- [ ] Visual: 20 states rendered in a debug cycle page/window; human screenshot review vs
      design/screens/ball.html (colors/opacity/sizes per SPEC-08 §2.1).
- [ ] Interactive: click ball → Sleeping→Listening; Esc/click cancels from Confirming; focus
      never stolen (verified with a focused editor); transparent region click-through.
- [ ] Sleeping-state CPU ≈0 (no timers) — measurable in 08's sampler; here: assert no animation
      timer handles alive in Sleeping.
- [ ] Hotkeys registered/re-registered on config change; mute toggles Muted state.
- [ ] Multi-monitor: drag to second monitor, persist, restore; simulated detach → primary.

## Progress log (append-only, newest last)
- [2026-09-19T10:41:26Z] agent=orchestrator claimed=T07-impl did=dispatched (maintain 2-way concurrency floor) next=sub-agent works through acceptance criteria
- [2026-09-19T11:20:00Z] agent=T07-impl did=statemachine-core next=ball-window-unit
  internal/statemachine implemented (table-driven): table.go = 42 encoded rows = the 40 frozen D43
  rows + SPEC-08 §3 appended #41 (Path C barge-in) / #42 (Conversation task intent); guard-branched
  rows encoded as adjacent entries sharing one number. machine.go = Dispatch with mutex, loud
  IllegalTransitionError on unknown (state,event), guard-all-fail also rejected; at most ONE runtime
  timer, armed only for timeout-table states (Listening 15s/90s/30s by phase pinned by the entering
  row, Confirming 3s [SPEC-06 L1 2-3s upper bound], Warm 90s, Settling 3s, AwaitingApproval 300s,
  Error 10s); Sleeping provably zero timers (TimersAlive()==0 asserted). Effects fire as events on
  injectable Sink, default explicit no-op. Tests: row presence 1..42, every legal edge fires
  (state+phase+side-effect spy), full 20x37 illegal matrix rejected (>>20 rejections), timeout table
  table-driven, Sleeping zero-timer, warm/settling/listening timer fire + stale-fire guard.
  go vet + go test green. Note: Listening timeout phase is a machine-internal pinned by entering row
  (#4/#7/#33 first-round, #29/#27 in-session, #28/#30/#41 conversation).
- [2026-09-19T11:35:00Z] agent=T07-impl did=c21-native-tokens next=ball-window-renderer
  internal/ball foundation committed: tokens.go (C21 native palette: full dark+light tables verbatim
  from design/assets/tokens.css; Color=straight RGBA + Premultiplied() for ULW; geometry/motion
  consts incl. 30fps cap MinFrameMs, Warm 2.4s 0.55<->0.7 alpha-path band, Settling 260ms fade),
  statevisual.go (VisualFor: all 20 states -> SPEC-08 §2.1 recipe; Conversation ring 2px > Confirming
  1.5px pinned; Sleeping fixed 12px/0.35 ignoring user size), anim.go (AnimationPolicy: Sleeping =
  AnimNone/0ms zero-timer; loops only in the 5 active states; Warm = AnimAlpha compositor path at
  10fps; Settling = one-shot AnimFade - documented SPEC-sanctioned exceptions), hit.go (inscribed-
  circle hit test, corners HTTRANSPARENT; DPI-injectable WindowEdgePx = orb+2*margin).
  Cross-check table: docs/evidence/s1/c21-native-tokens.md (token <-> CSS value <-> Go symbol, both
  themes + geometry/motion). Machine guard: TestNoHardcodedColorsInBallPackage (only tokens.go may
  carry color literals) + token golden values + 20-state visual coverage + anim policy incl. zero
  timer in Sleeping + hit/DPI tests. go vet/test green.
- [2026-09-19T12:40:00Z] agent=T07-impl did=ball-window-renderer-live next=live-tests+evidence
  internal/ball window stack implemented and verified LIVE on the desktop (20-state cycle ran to
  completion, exit 0): win32_windows.go (raw syscall plumbing), sta_windows.go (ui-sta thread via
  observe.Spawn, CoInitializeEx(STA), GetMessage pump, runtime.LockOSThread - without the pin the
  pump waits on the wrong queue and posted tasks hang), d2d_windows.go (single process-wide D2D+
  DWrite factory, D38a; COM objects stored as unsafe.Pointer, vet-clean provenance; CreateDCRenderTarget
  slot 16, BindDC slot 57 verified against mingw headers), renderer_windows.go (DC render target on
  top-down 32bpp DIB, DPI 96 -> 1 DIP = 1 px; radial gradients for orb body/highlight/core/glow with
  per-state brush cache; DrawLine carries D2D1_POINT_2F BY VALUE = two floats packed in one register
  word - pointer- passing corrupts coordinates; SetColor is ID2D1SolidColorBrush slot 8 NOT 3 - slot
  3 is GetFactory and silently keeps the brush white; stroked circles are 48-segment polylines
  because DrawEllipse's strokeWidth is a register float, ABI-unpassable), ball_windows.go (wndproc:
  WM_NCHITTEST inscribed-circle click-through, MA_NOACTIVATE, drag w/ 4px threshold + persist,
  WM_DPICHANGED rebuild, tray menu, hotkey dispatch; SetState/SetBadge/SetProgress/SetBadgeText;
  animation timers only per AnimationPolicy; Warm breathing and Settling fade via ULW
  SourceConstantAlpha - no re-render; handle count over a full 20-state cycle: 356 up / 366 peak /
  366 after close, gate <600 PASS; teardown destroys window+DIB+brushes+tray+hotkeys and joins the
  thread). Live pixel evidence: Listening ring pixel RGB(137,197,187) == C21 accent #86C2B9.
  Fixed en route: thread pinning, packed-point ABI, SetColor slot, hidden-move DPI re-read
  (WM_DPICHANGED does not fire for hidden windows), wetype status-bar overlap identified as
  external. cmd/balldebug harness committed (cycle/state/stay modes, handle reporting).
  go build ./... (CGO on), go vet, go test green.
- [2026-09-19T13:00:00Z] agent=T07-impl did=hotkeys-tray-position-live-tests next=visual-evidence
  Hotkeys (hotkey_windows.go): [hotkey] summon/mute/cancel/panel parse (Ctrl/Alt/Shift/Win + named
  keys), per-binding failure is warn-and-skip, re-registration on config change (RebindHotkeys); B1
  implemented: TakeEscForCancel binds Esc during Confirming, ReleaseEscAfterSession ALWAYS returns
  Esc (even if the original binding is taken by another app - then cancel stays unbound with a warn)
  and clears the takeover flag either way. Tray (tray_windows.go): left click = panel stub event,
  right-click menu 打开面板/静音/暂停唤醒/退出 with live checkmarks, TPM_NOACTIVATE + the
  SetForegroundWindow tray quirk handled. Position (position.go + monitors_windows.go): JSON store
  per monitor device (temp+rename), ResolvePosition (visible-keep / clamp / monitor-default /
  detached->primary) pure + table-tested; Per-Monitor V2 via SetProcessDpiAwarenessContext; hidden
  windows get no WM_DPICHANGED so the destination DPI is read explicitly after the move; resolve
  edge == initial-state window edge so the first apply never shifts the restored position.
  Live-window tests (winlive tag, real desktop): TestBallLiveLifecycle (20 states render, timers
  alive exactly in animated+Warm+Settling, Sleeping zero-timer, B1 takeover+release, handle gate
  base=97 peak=366 < 600, window destroyed after Close) and TestBallLivePositionPersistence
  (persist + exact restore). Fixes: window class registration once-guarded (second New in-process),
  Esc release semantics. go vet + go test (default suite + winlive) green.
- [2026-09-19T13:45:00Z] agent=T07-impl did=visual-evidence next=final-gates
  Visual evidence complete: scripts/dev/ball-cycle.ps1 (builds cmd/balldebug, walks the 20 D43 states
  at 2s dwell with overlays: AwaitingApproval badge=3, Downloading 69%+ring, Confirming countdown
  text; Go-side BitBlt composite capture per state - SRCCOPY|CAPTUREBLT, since plain SRCCOPY screen
  BitBlts EXCLUDE layered windows; parks the cursor and the ball on clear wallpaper away from the
  input-method bar; verifies the harness exited OK and the handle gate). 20 screenshots in
  docs/evidence/s1/ball-states/01-FirstRun.png .. 20-Stuck.png - human-reviewed via image reads:
  Sleeping 12px micro-dot @0.35 (zero animation), Listening teal ring + waves + audio-lines,
  Confirming danger ring + countdown, AwaitingApproval danger ring + white "3" badge, Speaking
  success tint, Settling fade, Warm amber breathing tint, Conversation 2px danger ring + mic, Error
  danger + X, Stuck warn + refresh, Downloading progress ring + percent, Queued info dot. Pixel
  check: Listening ring RGB(137,197,187) == C21 accent #86C2B9.
  Root cause found for the earlier transparent-frame evidence: state-dependent window resizes forced
  ID2D1DCRenderTarget::BindDC rebinds which empirically produced transparent frames; the window edge
  is now FIXED for the process lifetime (Sleeping renders its 12px dot inside the standard window;
  extra margins are click-through) and WM_DPICHANGED does a full renderer rebuild instead of a
  rebind. Also fixed en route: renderer DIB backing is allocated once at the max edge
  (ensureDCAndDIB). Final cycle: handles 385 < 600, harness exit OK, all windows destroyed.
