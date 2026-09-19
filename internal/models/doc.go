// Package models owns local model distribution (SPEC-01 §3): download,
// verification, resume, local path pointing (D26) and the C29 manifest /
// minisign verification.
//
// Responsibilities:
//   - download with resume against the mirror (env-forked endpoint, SPEC-03
//     §5.2; values land in ticket 06)
//   - sha256 + manifest signature verification (C29) before activation
//   - exposing model file locations to speech
//
// Non-responsibilities:
//   - no inference (speech), no config of provider API keys (secret/config),
//   - no mirror/mock server (tools/mockllm and compose own those, SPEC-11)
//
// DEFERRED(download/verify): implemented by ticket 14. This ticket only
// freezes the package boundary.
package models
