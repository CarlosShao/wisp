# 27 — Punctuation restoration (P4): CT-Punc model, dictation-grade output

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 15-speech-engines-cer-harness
**Parallel slots:** 1
**Spec refs:** SPEC-04 §4, D32 16.3.4 (≤300ms), P4, S4

## What to build
Punctuation restoration stage in the speech pipeline (CT-Punc-class int8 model via sherpa),
feeding punctuated transcripts to the agent — the make-or-break for dictation (scenario ③).

## Key constraints
- Latency: ASR-final → punctuated text ≤300ms (D32); measured on committed fixture set.
- Model residency: joins the Warm-resident set; participates in slot mutex accounting; unload
  at session end.
- Streaming ASR partials remain unpunctuated (display-only); finalization always punctuates.
- Toggle `[voice] punctuation=true`; disabled → raw transcript marked as such downstream.
- Quality: subjective punctuation correctness ≥8/10 on the dictation fixture (300-char
  continuous speech) — this is the S4 scenario-③ dependency; below → blocked-decision (try
  alternative model size before deferring).
- Failure: model load fail → degrade to unpunctuated + `model` error log + UI hint, never fail
  the utterance.

## Out of scope
- Dictation E2E acceptance (32); input.type (23).

## Acceptance criteria
- [ ] Punctuated E2E: wav → transcript with punctuation reaching agent (mock LLM asserts
      punctuation present).
- [ ] Latency ≤300ms P50/P95 on fixtures, recorded in SLO appendix.
- [ ] 300-char dictation fixture subjective score ≥8/10 (evidence sheet) or blocked-decision
      logged with next-step.
- [ ] Toggle off path works; load-failure degrade path works.
- [ ] Memory accounting: punct model included in Warm residency numbers (sampler).

## Progress log (append-only, newest last)
