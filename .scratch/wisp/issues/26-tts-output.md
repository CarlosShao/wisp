# 26 — TTS output: sherpa matcha engine, serial slot, Speaking pipeline, P7 gate

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 15-speech-engines-cer-harness
**Parallel slots:** ≤2 sub-agents (A: TTS engine + playback; B: half-duplex integration +
latency/P7 harness)
**Spec refs:** SPEC-04 §4, §8, D5, D16 half-duplex, D32 16.3.3/16.3.4, P7, S4

## What to build
`TtsEngine` (C9) on sherpa matcha/vits zh (int8): synthesize-and-play for short results, serial
engine-slot residency with ASR, half-duplex mic gating during playback, interrupt handling, and
the P7 subjective-quality gate harness.

## Key constraints
- TTS slot via the 15 engine-slot mutex: **Spike backfill (T02): serial switch measured P50 4726ms > 1.6s budget → D32 16.3.3 preset
      degradation ACTIVE: TTS RESIDENT + ASR on-demand is the default Path T policy (S4 may
      re-evaluate serial after tuning).** peak = max(ASR,TTS) not sum; unload via DisposalScope;
  measured switch latency vs "first-token→TTS-first-frame ≤800ms (loaded) / ≤1600ms (load)".
- Half-duplex (D16, **Path T only per D47**): during playback, mic capture closed AND KWS
  inference suspended (except veto-word channel rule in 41); playback end → Speaking→Warm
  transition (mic stays closed, default policy). **Path C (Conversation) is the D47 exception:
  mic STAYS OPEN with AEC (ticket 59/P15-gated), user speech mid-playback = barge-in → stop TTS
  within ≤400ms and never transcribe our own playback audio.**
- Interrupt channels during Speaking: hotkey / ball click / veto word (when KWS running) →
  stop synthesis+playback, release audio out, → Listening (D43 #27).
- Streaming: speak first sentence as soon as available (sentence-chunk synthesis) to meet the
  800ms budget; long text speaks summary only (routing at 30).
- Audio out: WASAPI render; volume follows system; no persistence of synthesized audio.
- P7 gate harness: fixed 20-sentence Chinese set; double-blind human scoring ≥7/10 recorded in
  evidence; below → fallback decision point (cloud TTS provider slot stays RESERVED, not
  implemented).
- Failure: TTS engine load fail → `model`/`resource` error class; playback device fail →
  `audio_device`; text-only fallback path (notify) still delivers result.

## Out of scope
- Result routing tiers (30); Conversation-mode mic policy internals (28, Path C owner); punctuation
  (27); AEC implementation (13 + spike 59 — this ticket consumes its gate for Path C barge-in).

## Acceptance criteria
- [ ] Speak path: short result → TTS audible, ball Speaking with success breathing; interrupt
      mid-playback → instant stop + Listening.
- [ ] Half-duplex (Path T): playback injects no ASR transcripts (self-echo fixture); capture gate
      closed during playback (13 counter asserts zero frames).
- [ ] **Path C barge-in (D47, blocked by 59): with AEC on, speak mid-playback → TTS stops ≤400ms →
      Listening; loop back our own TTS audio into the mic input → zero ASR transcripts (no
      self-excitation); P15 numbers recorded (CPU sustained, false-trigger incl. external audio).**
- [ ] Latency: first-audio ≤800ms warm / ≤1600ms cold-load measured (10 runs, P50/P95 recorded).
- [ ] Serial slot: ASR and TTS never co-resident (sampler peak = max, asserted); unload frees.
- [ ] P7: score sheet committed ≥7/10 or blocked-decision logged for cloud-TTS fallback.
- [ ] C25 wiring (registered by ticket 19's DEFERRED(C25-loop-wiring), adversarial report
      N-6): every announced string goes through `risk.Provenance.CheckText(scope,
      risk.ChTTS, text)` before playback; a hit means the announce is replaced by the L2
      confirmation (source named) and never spoken aloud. Test: the TTS channel assertion in
      `internal/risk/provenance_test.go` (TestFourChannelExfilSuite) plus a 26-side gate test.

## Progress log (append-only, newest last)
