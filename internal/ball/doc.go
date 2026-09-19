// Package ball owns the floating ball UI surface (SPEC-01 §3): the Win32
// layered window, Direct2D/DirectWrite drawing, mouse hit-testing and drag,
// multi-monitor positioning, and the tray icon.
//
// Responsibilities (frozen module table SPEC-01 §1.2/§3):
//   - layered window creation/paint on the single ui-sta STA thread (D38a)
//   - mouse/drag interaction, multi-monitor and DPI repositioning (D42#3)
//   - tray icon and its minimal menu
//
// Non-responsibilities:
//   - no business logic: state transitions live in statemachine, session
//     control in session, results routing in tools/panel
//   - no WebView2 hosting (panel), no global hotkey registration itself
//     (ui-sta thread is provided here; hotkey policy is ball+statemachine)
//
// DEFERRED(window/drag/tray): implemented by ticket 07 (ball core) and
// refined in S7 (D42#3 display-topology handling). This ticket only freezes
// the package boundary.
package ball
