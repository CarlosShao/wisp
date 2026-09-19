// Package secret owns the SecretStore (C28): DPAPI (CryptProtectData,
// CurrentUser) blobs referenced by indirection keys, plaintext first-run
// migration, and last-4-only logging.
//
// Responsibilities:
//   - store/resolve by ref ("dpapi:<blob-id>" | "env:NAME")
//   - plaintext api_key migration with backup + user-visible notice (D33)
//   - portable-mode fallback: undecryptable dpapi refs produce an explicit
//     error pointing to env:, never plaintext (P13)
//
// Non-responsibilities:
//   - no config parsing (config), no key material generation for C29 signing
//     (models/manifest), no env-fork data dir policy (proc)
//
// DEFERRED(store/migration): implemented by ticket 06. This ticket only
// freezes the package boundary.
package secret
