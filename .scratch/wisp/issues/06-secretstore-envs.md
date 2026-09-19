# 06 — SecretStore (C28, DPAPI) + WISP_ENV environment isolation

**Status:** in-progress
**Claimed by:** orchestrator -> sub-agent T06-impl
**Last update:** 2026-09-19T09:31:25Z
**Blocked by:** 03-skeleton-runtime-rules
**Parallel slots:** ≤2 sub-agents (A: DPAPI store + plaintext migration; B: env forks + test
env injection)
**Spec refs:** SPEC-03 §5–§6, SPEC-02 §6, C28, D33 config-security, P13

## What to build
`secret` module: DPAPI-backed secret store with reference indirection (`api_key_ref =
"dpapi:<blob-id>" | "env:NAME"`), first-run plaintext migration, last-4-only logging; plus the
`WISP_ENV = dev|test|prod` fork: per-env data dirs, mutex names, default endpoints, and the
visible env badge contract consumed by ball/panel.

## Key constraints
- Store: `CryptProtectData` (CurrentUser) → one blob file per ref under `<data>\secrets\`;
  resolve API takes ref string, returns secret; logs/diagnostics show only last 4 chars.
- Plaintext migration: on first run, any plaintext `api_key` in config.toml → migrate to DPAPI
  blob, remove field, back up original file, log a user-visible notice (D33).
- Env forks (SPEC-03 §5.2): prod `%APPDATA%\wisp\` + mutex `Local\wisp-single-instance` +
  auto-update checks on; dev `%APPDATA%\wisp-dev\` + `Local\wisp-dev-single-instance` + log debug
  + update checks OFF + default LLM base_url `http://127.0.0.1:18080/v1` + default mirror
  `http://127.0.0.1:18081`; test = `WISP_TEST_DATA_DIR` or `%TEMP%\wisp-test-<pid>\`, NO mutex
  registration, injectable endpoints. `go run` defaults dev; release builds default prod.
- Portable mode (`portable.txt`) overrides data dir to exe-relative `data\` (dev → `data-dev\`),
  taking precedence over env dir choice (SPEC-02 §6).
- P13 (portable × DPAPI): implement the recommended fallback — `dpapi:` refs that fail to decrypt
  under portable mode produce an explicit error guiding to `env:`; never plaintext. Record decision
  in PRECHECK.
- Env identity must be exposed for UI badge ("Wisp · dev") — provide `buildinfo.Env()` +
  data-dir summary API; badge rendering lands in 07/35.

## Out of scope
- Mock-llm server itself (09); panel badge UI (35); update channel logic (44/56).

## Acceptance criteria
- [ ] Round-trip: store → resolve; blob file on disk not readable as plaintext; log shows last-4 only.
- [ ] Plaintext migration test: seeded plaintext config → migrated, field removed, backup exists,
      notice logged; second run idempotent.
- [ ] Env fork matrix test: three envs → distinct data dirs/mutex names/endpoints; dev and prod
      dirs mutually invisible.
- [ ] Mutex names per env (integration test on Windows runner).
- [ ] Portable mode: data dir relocation + dpapi-decrypt-failure → explicit error (not plaintext).
- [ ] `env:` refs work with arbitrary dummy values (CI-friendly).

## Progress log (append-only, newest last)
- [2026-09-19T09:31:25Z] agent=orchestrator claimed=T06-impl did=dispatched (3-way concurrency) next=sub-agent works through acceptance criteria
- [2026-09-19T09:50:45Z] agent=T06-impl did=env fork complete (test data dir WISP_TEST_DATA_DIR/TEMP per pid, no mutex; dev endpoints+update-off field defaults; portable.txt override; Summary/EnvBadge API; boot LayoutErr removed; mutex-per-env integration test) next=secret DPAPI store
- [2026-09-19T10:00:33Z] agent=T06-impl did=SecretStore C28 (dpapi:/env: refs, CryptProtectData CurrentUser blobs 0600 under <data>/secrets, RedactSecret last-4, no-plaintext-in-log tests, P13 ErrPortableDecrypt explicit env: guidance) next=plaintext migration
