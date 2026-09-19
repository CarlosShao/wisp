// Package ball owns the floating ball UI surface (SPEC-01 §3): the Win32
// layered window, Direct2D/DirectWrite drawing, mouse hit-testing and drag,
// multi-monitor positioning, and the tray icon.
//
// Responsibilities (frozen module table SPEC-01 §1.2/§3):
//   - layered window creation/paint on the single ui-sta STA thread (D38a):
//     WS_EX_LAYERED|NOACTIVATE|TOOLWINDOW|TOPMOST, UpdateLayeredWindow with
//     32bpp premultiplied ARGB, D2D + DirectWrite (no GDI text) (SPEC-08 §2)
//   - the C21 native token tables (tokens.go; CSS twin + evidence doc)
//   - per-state visuals (statevisual.go: SPEC-08 §2.1 recipes) and the
//     animation policy (anim.go: zero timers in Sleeping, loops only in the
//     five active states, 30fps cap)
//   - mouse/drag interaction, click-through via WM_NCHITTEST, multi-monitor
//     and DPI repositioning (D42#3), per-monitor position persistence
//   - tray icon and its minimal menu; global hotkeys incl. the B1 Esc
//     takeover/release cycle
//
// Non-responsibilities:
//   - no business logic: state transitions live in statemachine; the ball
//     surfaces user gestures as Events and renders the state it is told
//   - no WebView2 hosting (panel); the shared D2D factory is exposed via
//     FactoriesForPanel for the future panel renderer (one factory, D38a)
//
// Implemented by ticket 07 (S1: static states + basic animations; full
// visual polish gate is the human acceptance at ticket 12).
package ball
