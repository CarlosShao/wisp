# 32 — S4 acceptance: scenario ③ dictation + scenario ④ chat companion gates

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 27-punctuation, 29-memory-l1l2, 30-result-routing-d10, 31-reminders
**Parallel slots:** ≤2 sub-agents (A: scenario ③; B: scenario ④)
**Spec refs:** §10 #12 #13, D44(b) rows ③④, P4/P7 conclusions, S4 done

## What to build
The S4 slice gate executing the two human-experience scenarios end-to-end on real hardware:
③ dictation-to-text-in-place and ④ multi-turn voice companion, each with the exact PASS bars
from the verification plan.

## Key constraints
- Scenario ③ dictation: continuous 300-char dictation → text appears AT THE TARGET APP's cursor
  via input.type (saving to a file does NOT count as pass); punctuation subjective ≥8/10;
  plus one elevated-window attempt → explicit UIPI message (§16.11#2).
- Scenario ④ companion: 5 consecutive conversation rounds; round ≥2 first-char ≤1.5s (model
  excluded — Warm alive, evidenced by round1-vs-round2 delta >1s); TTS subjective ≥7/10 (P7
  sheet); remembers one fact from round 1 (L1 profile live). **D47 additions: barge-in gate —
  mid-playback speech stops TTS ≤400ms and never transcribes our own playback; external-audio
  case (music playing in background) false-trigger rate recorded against P15 thresholds;
  Conversation ASR+TTS co-residency memory numbers backfilled into SLO table (or degradation
  mode verified).**
- Warm/Conversation behavior observed: default Warm (mic closed, warm breathing); explicit
  Conversation enable (privacy confirm once, red ring); 30s silence → Warm.
- Result routing during ④: mixed short/long replies hit correct D10 tiers.
- Evidence: screen recordings + logs + tool traces committed under docs/evidence/s4/; human
  (user) sign-off required for both subjective scores (D29 rule: human acceptance is the bar).

## Out of scope
- Panel UI (results show via notify/ball/native primitives until 36); KWS.

## Acceptance criteria
- [ ] Scenario ③ PASS: video + transcript + input.type trace; elevated-window failure message
      captured; score sheet ≥8/10.
- [ ] Scenario ④ PASS: 5-round log; latency table (round1 vs round2 delta >1s proving Warm);
      TTS ≥7/10; fact-recall demonstrated. **D47: barge-in ≤400ms gate + no-self-transcription
      PASS; external-audio false-trigger evidence; Conversation memory row backfilled.**
- [ ] Conversation privacy-confirm-once + red-ring visuals verified by user.
- [ ] D10 tier hits during ④ consistent with 30's matrix (no regressions).
- [ ] Pre-mortem: TTS device yanked mid-speak; mic device yanked mid-dictation — correct error
      classes + no stuck states.

## Progress log (append-only, newest last)
