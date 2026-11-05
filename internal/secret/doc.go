// Package secret owns the SecretStore (C28): the platform-sealed blobs
// (Windows: DPAPI, CryptProtectData, CurrentUser - see protect.go) referenced by
// indirection keys, plaintext first-run migration, and last-4-only logging.
//
// Responsibilities:
//   - store/resolve by ref ("dpapi:<blob-id>" | "env:NAME"); one blob file
//     per ref under <data>\secrets\ (SPEC-02 §6)
//   - plaintext api_key migration with backup + user-visible notice (D33)
//   - portable-mode fallback: undecryptable dpapi refs produce an explicit
//     error pointing to env:, never plaintext (P13, SPEC-02 §6)
//   - RedactSecret: the only sanctioned way to put secret material into a
//     log line (last 4 characters at most)
//   - the platform protector seam (protect.go): C28's Windows mechanism is
//     DPAPI; non-Windows builds seal with protect_other.go so the store's
//     platform-neutral contract is exercisable outside Windows (Linux is not in
//     the support matrix, D7) without touching the Windows path
//
// Non-responsibilities:
//   - no config parsing/normalization (internal/config owns that), no C29
//     signing key material (models/manifest), no env-fork data dir policy
//     (proc) - callers pass the data dir in
package secret
