// Package models owns local model distribution (SPEC-01 §3, SPEC-04 §7):
// the C29 signed model manifest, minisign signature verification, verified
// multi-mirror downloads with HTTP Range resume, and the Downloading-state
// wiring (D26/D33-F3).
//
// Security posture (C29, CRITICAL):
//   - models/manifest.json is committed signed (models/manifest.minisig);
//     verification is OFFLINE against the hardcoded public key in buildinfo
//     and happens BEFORE any network use. A manifest that does not verify is
//     never parsed into actionable metadata.
//   - hashes come exclusively from the signed manifest, never from a mirror
//     (D33/F3: same-source hash = zero auth value). Mirrors serve bytes only.
//   - verify_signature is not disableable: config hard-rejects false, and
//     NewManager double-checks at runtime.
//   - every downloaded byte is sha256-checked against the manifest before it
//     leaves staging; integrity failure = reject+delete. Install only happens
//     via atomic move out of staging.
//
// Package layout:
//   - manifest.go    C29 manifest types, parsing, signed loading (LoadSignedManifest)
//   - minisign.go    hand-written minisign format verify (Ed25519 + Blake2b pre-hash)
//   - downloader.go  Manager: mirror chain, Range resume, retry/backoff,
//     progress events, cancel, sha256 gate, install
//   - archive.go     tar.bz2 safe extraction (stdlib, no new dependencies)
//   - bridge.go      DownloadingBridge: D43 row #2/#37 state walk + ring ticks
//
// Non-responsibilities:
//   - no engine loading/inference (speech, ticket 15), no GUI model page (40),
//   - no mirror server (compose model-mirror owns that, SPEC-11).
package models
