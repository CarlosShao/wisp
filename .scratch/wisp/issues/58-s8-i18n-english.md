# 58 — S8: i18n / English locale + English voice models — DEFERRED (voice-chain caveat)

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 54-s7-acceptance + user go-ahead + English voice model selection (P-new)
**Parallel slots:** ≤2 sub-agents
**Spec refs:** D23 corrected i18n row, SPEC-12 §5 i18n caveat, D9, S8

## What to build
Internationalization done right per the corrected DEFERRED row: this is NOT just UI strings —
the whole voice chain is Chinese models, so non-Chinese users are voice-dead. Deliver UI string
externalization + full English locale + English ASR/TTS/KWS model selections passing
English baseline sets.

## Key constraints
- UI strings externalized from S1 onward (ticket 07/34 kept them out of code); locale files +
  `app.language`; RTL not in scope.
- English voice: sherpa-compatible English KWS/ASR/TTS candidates evaluated (license per P3
  discipline); baseline wav sets (en near_clean/noisy_far) + CER gates re-established (6%/15%
  analogs); punctuation model for English.
- Locale switch = hot; voice-locale = reload tier; mixed-language utterances documented
  (best-effort).
- Docs/README bilingual minimum (EN + ZH).

## Acceptance criteria
- [ ] All UI via locale files; zero hardcoded Chinese in frontend/native render paths (scan).
- [ ] English voice chain passes its committed baseline gates; latency budgets re-measured.
- [ ] Locale switch e2e (UI + voice reload) without restart.
- [ ] DEFERRED row closure evidence filed.

## Progress log (append-only, newest last)
