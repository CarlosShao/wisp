# 15 — Speech engines: VAD + streaming ASR via sherpa, engine slot mutex, CER harness

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 02-s0-spike, 13-audio-capture, 14-model-distribution
**Parallel slots:** ≤2 sub-agents (A: sherpa bindings + VAD/ASR engines + slot mutex; B: CER
harness + baseline sets + latency probes)
**Spec refs:** SPEC-04 §4, §8, D5, D32 16.3.3/16.3.4, S2, D44 CER gates

## What to build
`speech` module: `AsrEngine`/`WakeWordEngine` (slot reserved) C9 implementations on sherpa-onnx —
silero VAD endpointing + streaming Paraformer (int8) feeding transcripts into the agent — with
the engine-slot mutex (serial ASR/TTS residency), `intra_op_num_threads=1` enforcement, and the
CER measurement harness with the two committed baseline wav sets.

## Key constraints
- One shared onnxruntime instance for the whole chain (D5); `intra_op_num_threads=1` pinned
  (D32) — assert via session options in code + measured CPU shape.
- VAD: frame 512@16k; end-of-speech silence ~700ms (configurable), min speech ≥300ms (D43#11),
  max utterance 60s guard (SPEC: auto-finalize + notify).
- Streaming ASR: partials for ball latency, final on VAD stop → punctuation is a LATER stage
  (27) — transcript emitted unpunctuated is acceptable at this ticket; `asr` error class for
  pure-silence/too-short → one retry (D37), then user-visible 「没听清」.
- Engine slot mutex: **Spike backfill (T02, 2026-09-19): measured ASR↔TTS serial switch P50 4726ms vs 1600ms
      budget → D32 16.3.3 preset degradation is ACTIVE: default policy = "TTS resident + ASR
      on-demand" (serial scheme re-evaluated only at S4). Also: settle residual obligation —
      `debug.SetMemoryLimit`/GOGC tuning is part of this ticket's S1 scope (see ticket 12 note).**
 at most one big-model family resident (ASR xor TTS) — **Path T only (D47)**.
  **Path C (Conversation) is FULL DUPLEX: ASR+TTS+AEC co-resident during playback**; serial
  mitigation unavailable there — Conversation memory budget is filled in by S4 measurement
  (ticket 59/32); approved degradation = suspend ASR during playback, keep VAD+AEC for barge-in
  detection. Switch latency (Path T) measured and reported against the 1.6s budget.
- CER harness: `testdata/asr-baseline/near_clean/` (target ≤6%) and `noisy_far/` (≤15%, non-
  blocking but must report + UI hint requirement logged); sets COMMITTED to repo; runner emits
  per-file CER + aggregate; wired to a CI script (nightly/full self-hosted per 08).
- Latency probes: wake→listening ≤300ms stub (via injected wake event), VAD-stop→ASR-final
  ≤800ms streaming / ≤1500ms offline, measured with wav injection at C8 seam.
- Transcript taint marking hook: mark transcript as user-speech source (consumed by 19).

## Out of scope
- KWS (41); TTS (26); punctuation model (27); approval/security; cloud providers (RESERVED).

## Acceptance criteria
- [ ] Wav-injected utterance → VAD stop → transcript into agent (12's text path) end-to-end;
      ball walk Listening→Thinking verified.
- [ ] CER: near_clean ≤6% AND noisy_far ≤15% on the committed sets (CI-emitted report; noisy
      failure blocks merge only with a logged DEFERRED + UI-hint task, per D44).
- [ ] Endpointing: speak-pause-speak fixture splits correctly; <300ms speech ignored; 60s
      utterance auto-finalizes.
- [ ] Slot mutex: forcing ASR+TTS load requests → serial residency, unload frees memory
      (sampler delta), switch latency recorded vs 1.6s budget.
- [ ] Pure-silence input → one retry then `asr` error with 「没听清」 semantics.
- [ ] CPU: ASR inference single-core pinned (no thread-pool fan-out), sustained >3s single-core
      only within work-peak rules.

## Progress log (append-only, newest last)
