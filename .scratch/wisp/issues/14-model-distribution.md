# 14 — Model distribution: C29 signed manifest, mirrors, resume, Downloading state

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 03-skeleton-runtime-rules
**Parallel slots:** ≤2 sub-agents (A: manifest + minisign verify + downloader; B: state wiring +
local_override + mirror tooling/compose)
**Spec refs:** SPEC-04 §7, D26, C29, D33/F3, D42 model rows, P3/P5, S2 (P5 blocks S2)

## What to build
`models` module: signed model manifest, verified multi-mirror downloads with resume/retry/
progress/cancel, `Downloading` state wiring, local-dir escape hatch, and the compose
`model-mirror` service (normal/corrupt/bad-signature fixtures) used by tests and dev env.

## Key constraints
- `models/manifest.json` in-repo: per model `{id, purpose, urls[], sha256, size_bytes, license,
  quant}`; signed with **minisign**; public key hardcoded in `buildinfo` (C29). Verification is
  offline (manifest signature checked before any network use).
- Download flow: verify manifest signature → fetch bytes from mirror list (domestic mirrors
  first, official fallback; `[models] mirror` config) → verify file sha256 → install to
  `models\<id>\`. Mirror provides BYTES ONLY; hashes come exclusively from the signed in-repo
  manifest (F3: same-source hash = zero auth value).
- HTTP Range resume; retry with backoff; progress events (ball `Downloading` ring + percent);
  cancel cleans partial files to `staging\`.
- `[models] verify_signature=false` must ERROR, never disable (config already enforces; runtime
  double-checks).
- `[models] local_override{id:path}`: skip download, run integrity sha256 anyway (cheap), then use.
- P3: fill `license` field per model from license verification; any non-commercial finding = a
  blocked-decision log, model not shipped.
- Model set for this ticket: KWS zipformer 3.3M (bundled candidate — note: bundling decision
  recorded; default = download-on-first-use unless spike says bundle), silero VAD, streaming
  Paraformer int8, offline SenseVoice int8, CT-Punc int8, matcha-zh int8. Actual binaries pulled
  by 15; this ticket ships the pipeline + registry entries.
- Compose `model-mirror` (port 18081): serves good files, a corrupted file, a tampered manifest.

## Out of scope
- Engine loading/inference (15); GUI model management page (40); CPU-heavy CER (15/16).

## Acceptance criteria
- [ ] Tamper tests: flip one model byte → sha256 fail → reject+delete; tampered manifest →
      signature fail before download; `verify_signature=false` in config → hard error.
- [ ] Resume test: kill downloader mid-transfer → restart continues (byte-offset verified),
      completes, hash passes.
- [ ] Mirror failover: primary mirror 404 → fallback used; progress events emitted throughout;
      cancel leaves no partials outside staging.
- [ ] `Downloading` state walk on ball (enter/progress/exit) asserted via state log.
- [ ] local_override: model loaded from local path with integrity check, no network calls.
- [ ] manifest.json committed with license fields filled; P5 conclusion recorded in PRECHECK.

## Progress log (append-only, newest last)
