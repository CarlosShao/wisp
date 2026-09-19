# 41 — KWS wake word: Armed state, keywords+veto words, mute, opt-in privacy

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 15-speech-engines-cer-harness, 28-session-scope-warm
**Parallel slots:** ≤2 sub-agents (A: KWS engine + Armed lifecycle; B: veto words + L1 veto
integration + thresholds P14)
**Spec refs:** SPEC-04 §5, D2/D5/D16, B1, D43 #5,7–10, P14, S6

## What to build
`kws-zipformer-wenetspeech-3.3M` wake word engine (opt-in): Armed state lifecycle with ≤90/110MB
budget, keywords-file customization (no training), veto words sharing the model for the L1
block-window channel, one-key mute, and the P14 multi-keyword/threshold verification.

## Key constraints
- Opt-in only (`[voice] wake_word.enabled=false` default); enabling loads model (Armed);
  disabling unloads (Sleeping cap); mic stays OFF until wake or explicit listening (D16①).
- Armed budget: tree RSS ≤90MB (path Y) / ≤110MB (path X) per spike verdict; goroutine ≤7;
  audio buffers never persisted. `kws-infer` named goroutine.
- Keywords file: user-editable text list + per-word thresholds (`[voice] wake_word`); wake hit
  → Listening (KWS paused, resumes per state rules).
- **Veto words** (取消/停下/别) registered in the SAME model with per-word thresholds: hit during
  Confirming → veto (the B1 voice-cancel channel, ~100ms-class, no ASR load). P14: measure
  false-trigger rate of veto words under TTS/environment audio; too high → raise thresholds or
  disable channel with explicit UI note (decision point recorded).
- KWS-not-loaded honesty: Confirming UI shows 「语音取消不可用」 whenever KWS isn't running
  (short-cut-key sessions) — asserted e2e.
- Mute: global mute key → Muted (KWS inference stopped, model retained per #8); unmute
  restores; Armed visual per SPEC-08 §2.1.
- Wake latency: wake-word-end → Listening ≤300ms (D32) measured with fixture audio.
- Watchdog contract honored: Armed over-limit NEVER unloads KWS (42 rule; integration test).

## Out of scope
- AEC/barge-in (deferred); fast-path voice veto without KWS (deferred B1).

## Acceptance criteria
- [ ] Armed RSS/CPU gate measured (sampler) within D32 row; wake latency ≤300ms.
- [ ] Keywords: custom keyword list change → reload without restart; hit on fixture audio.
- [ ] Veto: KWS running + Confirming + veto-word wav → L1 call cancelled; KWS absent → UI
      message asserted.
- [ ] P14 false-trigger report (thresholds table + recommended values) committed.
- [ ] Mute cycle Armed→Muted→Armed with inference stop/start asserted (no model reload).
- [ ] Opt-in off → zero mic opens (device-handle probe) + Sleeping budget.

## Progress log (append-only, newest last)
