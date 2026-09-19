// Package plugin owns third-party extension execution (SPEC-01 §3): manifest
// loading and validation (C14/C16), the goja JS runtime sandbox, and the
// DisposalScope (C11).
//
// Responsibilities:
//   - Tier-1 manifest plugins (argv templating, exe_hash pinning) and Tier-2
//     goja plugins; goja JS throws become class "tool" errors for the LLM
//     (D37b), never cross into the Go stack
//   - DisposalScope (C11): the deterministic-disposal primitive shared by the
//     scope hierarchy session > task > tool > plugin (see disposal.go;
//     implemented in ticket 03 as a runtime rule, not a product feature)
//
// Non-responsibilities:
//   - no capability/risk decisions (risk), no plugin UI (panel), no download
//     of plugins (models is for models; plugin distribution is S7)
//
// DEFERRED(manifest/goja): implemented by tickets 50 (Tier-1) and 51 (Tier-2
// goja). This ticket implements only DisposalScope.
package plugin
