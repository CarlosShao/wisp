// Package panel owns the WebView2 panel host (SPEC-01 §3): C27 singleton
// window (hidden, not destroyed), PanelBridge (C17) with method whitelist,
// and the embedded web assets (D29: AddWebResourceRequestedFilter, no local
// HTTP server).
//
// Responsibilities:
//   - WebView2 lifecycle on the shared ui-sta STA thread; 3-5 msedgewebview2
//     child processes all join the Job Object (proc.JobScope)
//   - PanelBridge method whitelist + capability annotations + correlationId
//     routing + push backpressure/merge (D38d)
//
// Non-responsibilities:
//   - no React application logic (frontend/, ticket 34), no ball drawing
//     (ball), no business decisions
//
// DEFERRED(host/bridge): implemented by ticket 33 (host), ticket 35 (bridge).
// This ticket only freezes the package boundary.
package panel
