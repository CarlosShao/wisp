# 05 — config.toml full model: D36 sections, three effect tiers, hot reload, migration

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 03-skeleton-runtime-rules
**Parallel slots:** ≤2 sub-agents (A: schema structs + load/validate + migration; B: hot-reload
watch + tier semantics + security-section re-confirm gate)
**Spec refs:** SPEC-03 §2–§4, D36, §14.5, C15, D33 config-escalation

## What to build
`config` module: the complete `config.toml` section tree as Go structs with default tags (single
source of defaults), strict validation, three-tier effect semantics (hot/reload/restart),
file-watch hot reload, schema migration with backup, and the security-section loosening gate
(hook that fires an L2 re-confirmation callback; the actual confirm UI arrives later, S3/S5).

## Key constraints
- All sections/keys exactly per SPEC-03 §3 table (`app ball hotkey session voice audio llm agent
  risk fs net privacy memory panel cost plugins models observe`), types + defaults as listed.
  Provider presets: openai, anthropic, deepseek, qwen, zhipu, moonshot, siliconflow, openrouter,
  ollama (base_url/model defaults only; keys always via `api_key_ref`, never plaintext).
- Unknown key → error naming key + line number (never silently ignored).
- Hard-coded read-onlys: `audio.half_duplex=true` (writing false = error), `privacy.keep_transcript/
  keep_audio=false` (same), `models.verify_signature=true` (false = error, C29).
- Effect tiers: restart (`app`), reload (`voice` model swaps → emits reload event), hot (rest).
  🔒 sections `risk/fs/net/plugins`: any LOOSENING change → invoke re-confirm hook; rejected →
  keep old values; tightened → hot-apply. Both directions logged.
- Change detection: watchdog-loop mtime+size poll every 1s (no fsnotify dep, no new goroutine);
  GUI writes TOML → same path (single source of truth).
- Migration: `schema_version` key; auto-migrate + backup `config.toml.bak-<ver>`; unmigratable →
  Unconfigured + explicit guidance, never silently reset (would wipe allowlists = security bug).
- Defaults live ONLY in struct tags; TOML is serialization (D36 rule 3).
- **LLM catalog schema v2 (2026-09-19 user-approved supplement, SPEC-03 §3.1)**: `[llm]` carries
  `text_chain[]` (ordered fallback), `roles.{chat,memory_extract,handoff,summarize}` (per-role
  model+temperature+thinking_intensity), `providers.<name>` with `protocol/billing(plan|
  pay-per-token)/plan_credit_total_micro/compat{loose,allow_missing_usage,extra_headers}/rpm/tpm`
  and `models.<id>` with `capabilities{...}/context_window/max_output_tokens/price{...}/
  quota_daily_micro/quota_monthly_micro`; `[voice]` carries `realtime{...}` + `cloud_asr_chain[]` +
  `cloud_tts_chain[]` (local sherpa cascade = final fallback, never in chain). Chain elements must
  reference existing provider/model pairs (unknown → error naming the element). Storage split:
  catalog/quota = config.toml; probe/health = SQLite `provider_health` (SPEC-02 schema v2) —
  enforce the boundary in the structs.

## Out of scope
- SecretStore resolution of `api_key_ref` (06); GUI editor (39); risk enforcement (17+).

## Acceptance criteria
- [ ] Per-section valid/unknown-key/wrong-type test sets; unknown key error contains line number.
- [ ] Hot/reload/restart matrix test: writing a file mutation applies per tier (reload emits event).
- [ ] 🔒 loosening test: hook invoked; reject keeps old values; both directions logged.
- [ ] Migration: v-1 fixture migrates, backup exists; corrupt fixture → Unconfigured, file untouched.
- [ ] Hard-coded read-only fields reject writes with explicit errors.
- [ ] Round-trip: load → marshal → load yields identical structs (GUI double-source guard).
- [ ] Catalog tests: chain referencing unknown provider/model → error naming the element;
      capabilities/billing/quota round-trip; voice chains validated; thinking_intensity enum
      enforced; compat flag defaults.
- [ ] Storage-boundary test: probe/health fields have no representation in config structs;
      `provider_health` table exists per SPEC-02 schema v2.

## Progress log (append-only, newest last)
