# 05 — config.toml full model: D36 sections, three effect tiers, hot reload, migration (DONE ✅)

**Status:** done
**Claimed by:** orchestrator -> sub-agent T05-impl
**Last update:** 2026-09-19T13:56:25ZT09:20:36Z
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
- [x] Per-section valid/unknown-key/wrong-type test sets; unknown key error contains line number.
- [x] Hot/reload/restart matrix test: writing a file mutation applies per tier (reload emits event).
- [x] 🔒 loosening test: hook invoked; reject keeps old values; both directions logged.
- [x] Migration: v-1 fixture migrates, backup exists; corrupt fixture → Unconfigured, file untouched.
- [ ] Hard-coded read-only fields reject writes with explicit errors.
- [ ] Round-trip: load → marshal → load yields identical structs (GUI double-source guard).
- [x] Catalog tests: chain referencing unknown provider/model → error naming the element;
      capabilities/billing/quota round-trip; voice chains validated; thinking_intensity enum
      enforced; compat flag defaults.
- [x] Storage-boundary test: probe/health fields have no representation in config structs;
      `provider_health` table exists per SPEC-02 schema v2.
- note: AC#5 left open — no ruling found in docs/evidence/s1/05-adversarial-acceptance.md (its 5-row
      audit table searched for read-only/half_duplex/keep_transcript/verify_signature; no row names
      them). Candidate evidence exists unadjudicated at internal/config/validate_test.go:129
      TestValidateHardcodedReadOnly — needs an addendum ruling, not new work
- note: AC#6 left open — no ruling found in docs/evidence/s1/05-adversarial-acceptance.md (searched
      for round-trip/marshal/往返/TestRoundTrip*; row 1 only says "go test ./internal/config ok
      47 顶层/115 RUN"). Candidate evidence unadjudicated at
      internal/config/boundary_test.go:188 TestRoundTripLoadMarshalLoad

## Progress log (append-only, newest last)
- [2026-09-19T09:20:36Z] agent=orchestrator claimed=T05-impl did=dispatched (3rd concurrent slot trial) next=sub-agent works through acceptance criteria
- [2026-09-19T10:39:08Z] agent=T05-resume3 did=startup: read ticket+SPEC-03; verified baseline (CGO0 go test green except pre-existing cmd/wisp sherpa cgo gate, same exclusion as T06 closeout; no gcc in this sandbox - full CGO build deferred, documented at closeout); deleted untracked probe junk tmp_probe/ + cmd/tomlprobe_main.go.bak next=reconcile WIP schema/defaults/parse vs SPEC-03 §3 key-by-key
- [2026-09-19T10:53:05Z] agent=T05-resume4 did=startup+reconcile: read ticket+SPEC-03; CGO0 baseline re-verified (go vet config/secret clean; go test green except pre-existing cmd/wisp sherpa gate); key-by-key reconcile of WIP schema.go/defaults.go/parse.go vs SPEC-03 sec 3 + 3.1 PASS (all 17 sections + defaults match table; 12 presets incl minimax/mimo/stepfun; catalog v2 keys complete; SchemaVersionCurrent=2; Tier/Direction types in place; plugins raw-map absorber documented); schema.go comments already ASCII (only non-ASCII is VetoWords default = data not comment); committing reconciled WIP as baseline, then TDD: validate/loader/manager/CheckAndReload/migrate/tests next=validate.go
- [2026-09-19T11:00:00Z] agent=T05-resume4 did=catalog-validation: TDD catalog.go (catalog_test.go first) - text_chain/cloud_asr_chain/cloud_tts_chain elements must reference existing provider/model pairs (split at first '/', openrouter-style nested model ids kept), unknown element named verbatim in error per SPEC-03 sec 3.1, malformed "noslash" elements rejected, local-sherpa rejected from cloud chains (terminal fallback never a member), roles set-both-or-neither + pair existence, realtime pair checked only when enabled; go test ./internal/config ok next=loader (Load/Save/defaults merge/resolve hook)
- [2026-09-19T11:02Z] agent=T05-resume4 did=validate-core: TDD validate.go (validate_test.go first, 9 test funcs) - enums (theme/observe.level/proxy.mode/over_budget/protocol/billing provider+model/thinking_intensity/roles thinking_intensity), ball.size 44-72 + opacity_idle + alert_threshold ranges, hard-coded read-onlys (audio.half_duplex=false, privacy.keep_transcript/keep_audio=true, models.verify_signature=false all rejected), api_key_ref format via secret.ParseRef (providers + voice.realtime; plaintext/traversal rejected); all key-path errors; go test ./internal/config ok next=catalog chain validation
- [2026-09-19T11:04:15Z] agent=T05-resume4 did=loader: TDD loader.go (loader_test.go first, 13 test funcs) - LoadFile pipeline per SPEC-03 sec 4.1 (peek schema_version -> newer-version reject -> migrate dispatch [stub until migrate unit] -> strict decode onto NewDefaults -> normalize -> preset inheritance -> validate -> resolve), unknown-key errors carry 1-based line numbers (asserted "line 6"/"line 2"), type errors carry position, preset inheritance fills protocol/base_url for preset-named providers (never overwrites explicit), non-preset provider without protocol rejected naming key, SecretResolver interface + Resolved{ProviderKeys,RealtimeKey} (plaintext never enters Config; *secret.Store satisfies the interface directly), resolve failure names failing ref (Unconfigured mapping left to caller), SaveFile stamps current version + canonical marshal + atomicWrite; go test ./internal/config ok next=manager three-tier reload
- [2026-09-19T11:10:48Z] agent=T05-resume4 did=manager+reload-idempotence: TDD manager.go (manager_test.go first, 12 test funcs) - Manager{ConfirmLocked/OnReload/OnRestartPending callbacks, nil hook = deny fail-closed}, CheckAndReload() idempotent via mtime+size stat (no goroutine, watchdog tick calls it), hot/reload/restart matrix verified (ball hot applies; voice asr.model swap applies+emits OnReload; app.language keeps old value + OnRestartPending; app.theme hot inside restart section; voice punctuation/veto/thresholds/speed hot without event), locked sections risk/fs/net/plugins direction detection (additions loosen, removals tighten, delete_enabled/tier2/enabled flips, proxy none->open loosen, l1_window_sec raise loosen), loosening gates through ConfirmLocked - reject keeps old values (atomic per section), approved applies; BOTH directions logged via slog (captured in test); failed load keeps current config + adopts stat (no per-tick re-error), missing file errors but keeps current, Resolved() accessor exposes re-resolved secrets (llm api_key_ref change re-resolves on next reload); callbacks fire outside lock next=migration registry
- [2026-09-19T11:14:57Z] agent=T05-resume4 did=migration: TDD migrate.go registry (migrate_test.go first, 5 test funcs) - migrations map[int]migrationStep {1->2 migrateV1toV2}, ver 0 (no schema_version key) treated as v1, chain walks to current; v1->v2 mapping: llm.timeout->timeout_ms, llm.temperature->roles.chat.temperature, default_provider(+model)->text_chain[0]+roles.chat, fallback_provider->text_chain[1], flat providers.<n>.model->models.<model>{enabled=true}, everything else passes through; byte-identical backup config.toml.bak-<ver> written before atomic rewrite; SAFETY: migrated bytes decode+validate dry-checked BEFORE file replaced (never writes a non-loading config); unmigratable -> explicit error naming stuck version, file untouched, no backup (test pins corrupt fixture + unsupported version); tests verify backup identical, values mapped, non-llm sections preserved, migrated file reloads clean next=boundary/round-trip tests + closeout
- [2026-09-19T11:19:08Z] agent=T05-resume4 did=boundary+round-trip: TDD boundary_test.go - TestBoundaryNoPlaintextSecretFields (reflect walk over Config type graph: only APIKeyRef may carry key material, D36 rule 5 pinned structurally), TestBoundaryNoRuntimeObservationFields (no health/probe/latency/lasterror/verified fields - storage split; provider_health DDL deliberately NOT asserted here: ticket 09 owns it, memory schema.go documents absence), TestRoundTripLoadMarshalLoad (full every-section config incl plugins/compat/price/quotas/chains: save->load->save->load DeepEqual-identical, GUI double-source guard), TestRoundTripBytes (MarshalCanonical->decodeStrict identity), TestResolvedNeverPersists (refs verbatim in file, plaintext never); note: go-toml v2 marshals refs as literal strings api_key_ref = 'env:ACME' (single-quote form pinned); go test ./internal/config ok next=closeout (full vet/build/test matrix + handoff)
- [2026-09-19T11:21:30Z] agent=T05-resume4 did=handoff-to-orchestrator: ticket 05 implementation COMPLETE except items explicitly deferred (see below). Closeout matrix: gofmt -l internal/ clean; go vet ./internal/... clean; non-ASCII grep over internal/config clean (only remaining non-ASCII is the VetoWords default value - data, not comments; schema.go comments were already ASCII); CGO_ENABLED=0 go test -count=1 all packages green (buildinfo config memory observe plugin proc secret statemachine) except pre-existing cmd/wisp sherpa cgo gate (same exclusion as T05-resume3/T06 closeout baselines); internal/... builds under CGO0. Coverage: 47 top-level test funcs / 115 RUN entries in internal/config - per-section unknown-key with LINE NUMBERS (loader_test asserts "line 6"/"line 2" evidence), wrong-type with position, all enums, ball.size 44-72, hard-coded read-onlys (audio.half_duplex=false / privacy.keep_transcript|keep_audio=true / models.verify_signature=false all rejected), api_key_ref grammar via secret.ParseRef (providers+realtime, plaintext/traversal rejected), catalog chains name failing element verbatim + local-sherpa ban, roles set-both-or-neither, realtime checked only when enabled, thinking_intensity/thinking_levels enum, preset inheritance, resolve hook (SecretResolver interface; *secret.Store satisfies directly; Resolved{ProviderKeys,RealtimeKey}; failure names ref), three-tier matrix (hot applies / reload applies+OnReload event / restart keeps old+OnRestartPending; app.theme hot inside restart section; voice tuning keys hot), locked sections direction engine (risk/fs/net/plugins: additions loosen, removals tighten, flips classified; loosen -> ConfirmLocked hook, nil=deny fail-closed, reject keeps old values atomically per section; BOTH directions slog-logged, captured in test), CheckAndReload mtime+size idempotent no-op (asserted), failed load keeps current config + adopts stat (no per-tick spam), missing file errors but keeps current, migration v1->v2 registry with byte-identical .bak-1 backup + dry-check-before-write + corrupt/unsupported file untouched, full round-trip (save->load->save->load DeepEqual; marshal->decode identity on every-section config incl plugins/compat/price/quota), storage boundary (no plaintext api_key field possible; no health/probe/latency fields; provider_health DDL deliberately left to ticket 09 per memory/schema.go). Commits: fe126e2 c352769 18ab82c 1fd19d4 ab638dd ae838ef 02b9a72 (all pushed origin+cnb). Deferred/notes: (1) watchdog-loop wiring of CheckAndReload into its 1s tick = ticket 42 per watchdog/doc.go DEFERRED; (2) actual L2 confirmation UI = S3/S5, hook seam ready; (3) provider_health table = ticket 09; (4) -race skipped (requires gcc/CGO on windows, sandbox has no gcc - documented at T05-resume3); (5) full CGO build not run (no gcc) - CGO0 matrix green. Ticket NOT marked done (orchestrator decides).
- [2026-09-19T13:56:25Z] agent=orchestrator did=T05-adv PASS (orchestrator-executed; report docs/evidence/s1/05-adversarial-acceptance.md) next=ticket DONE
- [2026-09-20T02:25Z] agent=agent-bookkeeping-1 did=AC boxes reconciled against docs/evidence/s1/05-adversarial-acceptance.md (5-row audit table; AC#8's provider_health half cross-ruled in docs/evidence/s1/09-adversarial-acceptance.md §5 schema v2 审计 PASS): 6 checked, 2 left open (reasons above) next=none
