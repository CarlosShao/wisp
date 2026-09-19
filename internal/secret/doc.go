// Package secret owns the SecretStore (C28): DPAPI (CryptProtectData,
// CurrentUser) blobs referenced by indirection keys, plaintext first-run
// migration, and last-4-only logging.
//
// Responsibilities:
//   - store/resolve by ref ("dpapi:<blob-id>" | "env:NAME"); one blob file
//     per ref under <data>\secrets\ (SPEC-02 §6)
//   - plaintext api_key migration with backup + user-visible notice (D33)
//   - portable-mode fallback: undecryptable dpapi refs produce an explicit
//     error pointing to env:, never plaintext (P13, SPEC-02 §6)
//   - RedactSecret: the only sanctioned way to put secret material into a
//     log line (last 4 characters at most)
//
// Non-responsibilities:
//   - no config parsing/normalization (internal/config owns that), no C29
//     signing key material (models/manifest), no env-fork data dir policy
//     (proc) - callers pass the data dir in
package secret
