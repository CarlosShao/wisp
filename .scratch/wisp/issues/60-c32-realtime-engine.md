# 60 — C32 RealtimeEngine: cloud S2S realtime as Path C enhancement (gated)

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 59-p15-aec-spike (gate condition #1), 16-s2-acceptance, 28-session-scope-warm,
44-costmeter-c23
**Gate condition #2:** user provides a realtime provider Key — **if either condition is missing,
this ticket AUTO-DEFERS with a registry entry (never silently disappears)**
**Parallel slots:** ≤2 sub-agents
**Spec refs:** D47 (§16.12), C32 (PLAN contract row incl. D47 细化 ①②③), D23, F4/TTS-channel,
C23, D43 #41/#42

## What to build
The `RealtimeEngine` (C32): cloud speech-to-speech realtime as an OPT-IN enhancement layer of
Path C (Conversation) — audio-in/audio-out streaming with server-side endpoint detection and
barge-in events, ZERO tool permissions, full cost transparency, transcript persistence, and a
clean handoff back to the text loop when the user asks the assistant to DO something.

## Key constraints
- **Zero tool permissions (the entire security model)**: the realtime brain can only talk. It
  never sees tools, never produces tool_calls, never reads files. Approval gates (D4/D31/C19)
  are structurally unreachable from this path (ticket 21's scope note).
- **Context hygiene rule (C32 细化③ — hard rule)**: local file contents and taint-marked content
  MUST NOT be injected into the realtime context. What it can speak IS an exfiltration channel
  (F4's TTS channel equivalent); starve it of sensitive raw material. Its context = user speech
  transcript + L1 profile + conversation history only.
- **Handoff protocol (C32 细化①)**: intent detection = realtime model's structured handoff
  signal (a HOST-side trusted primitive, NOT a tool call); payload = current transcript +
  conversation context reference; realtime session SUSPENDS (not destroyed); the text loop
  (Path T) executes the task under FULL gating; result returns via notify/TTS; realtime can
  resume. D43 #42 transition.
- **Provider abstraction (C32 细化②)**: preset table + custom `base_url` (D8 pattern) — domestic
  candidates: **StepAudio 3 Realtime (FIRST candidate, user has StepFun account; note: it supports
  Voice-Agent tool calls on the provider side — our ZERO-tool rule still forces it off), 豆包 /
  Qwen-Omni realtime**, and OpenAI Realtime (proxy) evaluated at build time;
  provider differences hidden behind C32; `[voice] realtime{provider, model, api_key_ref,
  base_url}` config keys (hot/reload tiers per SPEC-03 conventions; add section row).
- **Cost (C23)**: audio tokens metered in real time (in/out), same per-task/daily aggregation;
  realtime tasks carry their cost line like every task.
- **Transcripts → task_log** (with §14.4 masking rules); L1 profile extraction unaffected
  (works off transcripts, provider-agnostic).
- **Privacy**: opt-in + explicit "audio will be uploaded" notice + first-enable L2 privacy
  confirm (distinct from Conversation's confirm if both apply — one combined confirm flow is
  acceptable, log both consents); red ring unchanged.
- **Barge-in**: uses AEC (59) + server endpoint events AND local AEC-cleaned VAD — whichever
  fires first; ≤400ms stop-or-overlap rule per Path C gate.
- Latency expectation: realtime turns target ≤1s first-audio (provider-dependent) — record
  measured numbers, do NOT invent an SLO gate until measured (target ≠ acceptance, D32 rule).

## Out of scope
- Making realtime the default chat brain (default stays local cascade ASR→LLM→TTS; realtime is
  an opt-in switch `[voice] realtime.enabled=false` default); work-mode changes; tool access of
  any kind; token-level memory of audio tone (future).

## Acceptance criteria
- [ ] E2E: Conversation with realtime provider (or recorded-real fixture if key absent at test
      time — provider-abstracted) → natural multi-turn: speak → ≤~1s first audio → barge-in
      mid-reply → resume.
- [ ] Zero-tool proof: adversarial prompts trying to induce tool behavior produce NO tool_calls
      and NO host-bridge invocations (bridge audit log empty for the session).
- [ ] Context hygiene: seeded taint/local-file injection attempt into realtime context →
      rejected by construction (test at the context assembly boundary).
- [ ] Handoff: "帮我把下载里的 PDF 归类" spoken in realtime mode → handoff signal → realtime
      suspends → text loop executes under full gating (L1/L2 as applicable) → result announced →
      realtime resumes; D43 #42 trace.
- [ ] Cost: audio token usage visible in task_log + cost_daily within one session.
- [ ] Gate behavior: without key or P15-fail → ticket auto-defers, registry entry written,
      Conversation continues on cascade path (no UX regression).
- [ ] Privacy: upload notice + consent flow logged; `keep_audio` still impossible (D16).

## Progress log (append-only, newest last)
