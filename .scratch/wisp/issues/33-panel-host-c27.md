# 33 — Panel host: C27 PanelManager singleton, WebView2 window, embed.FS serving

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 07-ball-state-machine-core, 12-cli-text-path-s1-gate
**Parallel slots:** 1
**Spec refs:** SPEC-08 §5.1, D29, C27, D32 panel rows, D42#11, S5

## What to build
The `panel` module host side: singleton PanelManager owning at most ONE WebView2 window per
session (hide-don't-destroy), embed.FS resource serving via `AddWebResourceRequestedFilter`
(no localhost server), cold/hot show paths, focus return, and the no-runtime fallback
declaration consumed by 37's native card.

## Key constraints
- `jchv/go-webview2` (pure Go); one window per session; hide/show instead of destroy within a
  session; destroy at session dispose (31 hook). Cold show ≤1500ms / hot show ≤200ms — measured
  and recorded (D32).
- Resources: `//go:embed assets/web/dist` + `AddWebResourceRequestedFilter` mapping to embedded
  FS; NO localhost HTTP listener (D29; D21 anti-pattern).
- WebView2 created on the shared `ui-sta` STA thread (D38a); environment notes from spike (02)
  recorded — single-window strategy sidesteps Environment sharing entirely.
- Focus: only the panel takes focus when shown; on hide/close → focus returns to the recorded
  previous foreground window (D29).
- Multi-task partition hook: render context keyed by correlationId (surface built at 36/38).
- Runtime-missing detection: creation failure → `panel.unavailable` event + no-crash; L2 native
  fallback flag set (37 consumes); guide user to Evergreen install, never auto-download.
- CSP header/meta injection point for the embedded page (actual CSP value enforced at 35/36 —
  provide the injection hook + default strict CSP now: connect-src 'none' etc., SPEC-06 §9).

## Out of scope
- Frontend bundle contents (34); bridge protocol (35); panel pages (36–40).

## Acceptance criteria
- [ ] Show/hide/destroy lifecycle: second show within session reuses window (process-tree
      child-count stable); session dispose destroys + WebView children exit ≤2s.
- [ ] Latency: cold ≤1500ms / hot ≤200ms (10-run P50/P95 into SLO appendix).
- [ ] Embedded resource serving works offline; no listening socket exists (netstat assertion).
- [ ] Focus round-trip: editor → panel → hide → focus back in editor (automated + manual).
- [ ] Runtime-missing fixture (rename/mask loader) → fallback event, app alive, native L2 flag on.
- [ ] CSP present on served documents (response inspection) with connect-src 'none'.

## Progress log (append-only, newest last)
