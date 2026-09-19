# 61 — Cloud voice providers via C9: cascade ASR/TTS from configured models

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 16-s2-acceptance, 11-llm-adapters-rest (probe/compat infra), 05-config-model (catalog schema)
**Parallel slots:** ≤2 sub-agents (A: cloud ASR provider + streaming; B: cloud TTS provider + voices)
**Spec refs:** C9, D47 Path C cascade, SPEC-03 §3.1 (voice chains), SPEC-12 registry row
（已采纳·提前 2026-09-19）, user requirement: 语音模型下拉

## What to build
The first CLOUD cascade voice providers behind the existing C9 seam (`AsrEngine`/`TtsEngine`
pluggability was reserved by D5 all along): cloud ASR and cloud TTS driven by **configured
models from the catalog** (models with `capabilities.audio_in` / `audio_out` bits) — this is the
implementation layer behind the frontend's default-voice-model dropdowns. First candidates:
**StepAudio 2.5 ASR** (streaming, 5-min audio in ~1s) and **StepAudio 3 TTS** (low-latency
streaming, laughter/hesitation) via StepFun's OpenAI-compatible-style audio endpoints; provider
abstraction keeps any later vendor drop-in.

## Key constraints
- C9 contract: cloud providers are interchangeable with local sherpa engines — same
  `AsrEngine`/`TtsEngine` interfaces, selected via `[voice] cloud_asr_chain[]` /
  `cloud_tts_chain[]` (ordered fallback chains from catalog entries with voice capability bits;
  chains validated against catalog, ticket 05).
- **Local sherpa cascade is the FINAL fallback and never enters the chains** (zero cost, always
  works offline) — chain exhaustion falls to local (visible event, never silent).
- Billing: audio usage metered through C23 (audio_in/audio_out pricing fields from the catalog
  entry); quota pre-dispatch check applies (ticket 44) — over-quota → advance voice chain.
- Privacy: cloud voice is opt-in per chain configuration; when a cloud ASR/TTS is ACTIVE, surface
  the "audio will be uploaded" state (D16 rule extended to cascade; Conversation red-ring rules
  unchanged). Local fallback keeps privacy posture intact.
- Streaming: cloud ASR must support streaming partials feeding the same VAD→ASR pipeline
  (SPEC-04 §4); TTS first-frame budget uses the cascade path budgets (cloud TTS may exceed local
  first-frame budget on cold network — measure, mark "provider-dependent", no invented SLO gate).
- Failure semantics: provider errors classified per D37 (`provider`/`rate_limit`/`network`) →
  advance voice chain (ticket 11 chain-stepping mechanism, shared with text_chain).
- Probes: audio_in/audio_out capability probes defined here feed `provider_health`
  (ticket 11's probe framework).
- Dependency whitelist: no new native deps (HTTP streaming only); SDK-less REST/SSE.

## Out of scope
- Realtime S2S (ticket 60/C32 — separate engine); local model changes; memory pipeline changes
  (extraction is transcript-based, provider-agnostic).

## Acceptance criteria
- [ ] Voice E2E with cloud chain: wav → cloud ASR (streaming partials) → agent → cloud TTS
      playback; ball state walk unchanged.
- [ ] Chain fallback: cloud ASR fail_next → advances to next chain entry → local sherpa as final
      (visible events at each hop); same for TTS.
- [ ] Quota integration: cloud voice model over quota → chain advances BEFORE dispatch (no
      negative), event visible.
- [ ] Cost: audio usage lands in task_log + cost_daily with audio pricing fields.
- [ ] Privacy state: active-cloud-voice indicator; local fallback restores privacy posture.
- [ ] C9 contract test suite passes for both cloud providers (interchangeability with local).

## Progress log (append-only, newest last)
