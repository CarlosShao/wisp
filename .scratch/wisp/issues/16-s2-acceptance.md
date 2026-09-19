# 16 — S2 acceptance: voice end-to-end gate, latency segments, hotplug pre-mortems

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 12-cli-text-path-s1-gate, 15-speech-engines-cer-harness
**Parallel slots:** 1
**Spec refs:** S2 done (D44), D32 16.3.4, §10 verification, SPEC-10 §3–§5

## What to build
The S2 slice gate: full voice-input chain (mic → VAD → ASR → agent → notify result → settle)
verified with wav injection and real-device smoke, all D32 latency segments measured against
budget, and the S2 pre-mortem items executed for real.

## Key constraints (must all be evidenced, not asserted)
- Chain demo: say command (wav) → ball Listening→Thinking→Acting→notify reply → Warm/Settling;
  transcript + tool evidence rows present.
- Latency segments vs D32 16.3.4: VAD-stop→ASR-final ≤800ms(stream)/≤1500ms(offline);
  context-assembly ≤50ms; ASR→punct placeholder ≤300ms (punct real in 27 — measure pipeline
  overhead only); cold-session first-token P50 ≤4s / P95 ≤6s (golden LLM), Warm P50 ≤2.5s
  (from 12's Warm data if present, else deferred note).
- Idle/Armed SLO re-run (sampler) after speech stack exists: Armed only when KWS lands (41) —
  here assert Sleeping unaffected and post-settle return ≤10s.
- Pre-mortems executed for real (SPEC-10 §7): mic occupied (test process holds device), mic
  permission denied, device hotplug mid-utterance, disk-full during transcript write, model
  file corrupted post-install (sha256 fail on load), process kill during download + resume.
- Every failure shows correct D37 class + user-visible message; nothing silent.

## Out of scope
- TTS/punctuation (26/27); KWS (41); security tools (17+).

## Acceptance criteria
- [ ] E2E voice demo recorded (log + state trace + tool rows), run on real hardware.
- [ ] Latency table filled with measured numbers vs budget, committed to docs/SLO.md appendix.
- [ ] All six pre-mortem scenarios executed with evidence and correct error classes.
- [ ] CER gates from 15 green in CI artifacts (near_clean mandatory; noisy_far status recorded).
- [ ] Adversarial acceptance pass by a different agent (stub scan + contract conformance).

## Progress log (append-only, newest last)
