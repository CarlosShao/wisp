# 59 — P15 AEC spike: license, binding, echo quality, CPU, false-trigger thresholds

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 13-audio-capture
**Blocks:** 26 (Path C barge-in part), 28 (Path C full duplex), 32 (scenario ④ gate), 60 (C32)
**Parallel slots:** ≤2 sub-agents (A: license/binding/build integration; B: measurement rig +
threshold report)
**Spec refs:** D47 (§16.12), P15, C32 gate condition, SPEC-01 whitelist entry, D32 16.3.3 D47 note

## What to build
The P15 feasibility spike that gates ALL Path C full-duplex work: verify webrtc-audio-processing
(AEC3) license + Go binding viability, measure echo-cancellation quality with our self-rendered
TTS PCM as reference, measure sustained CPU, and set the barge-in false-trigger thresholds —
including the external-audio-source case.

## Key constraints
- **License first**: webrtc-audio-processing is BSD-3 — verify the exact build/distribution we
  fetch (deps.toml entry + fetch-deps, SPEC-11 §2.1); any license surprise → Path C barge-in
  degrades to hotkey/click interruption, whitelist entry retracted (D47 rule; do NOT revive
  Path T full duplex).
- Reference signal = OUR OWN rendered TTS PCM (known exactly) + WASAPI render position —
  measure alignment quality and residual echo; clock drift concern is "an order of magnitude
  smaller" than generic AEC (D47 argument) — VERIFY, don't assume.
- **Cross-endpoint case**: capture and render on different devices (Bluetooth speaker + wired
  mic) — drift returns; measure and set expectations/thresholds.
- **External-audio case (missed by everyone, mandatory)**: user plays music/video while chatting
  — third-party sound is NOT in our reference; measure barge-in false-trigger rate (VAD on
  AEC-cleaned signal) with music/speech/video fixtures; threshold = external audio must not
  trigger barge-in (target: ≤1 false barge-in per 10min playback under normal room volume);
  if unachievable, Path C keeps barge-in but requires voice-activity evidence (e.g., direction
  of AEC residual), record honest limits.
- Sustained CPU: AEC+VAD processing during playback ≤ single-core 5% PROVISIONAL (SLO table
  row) — measure and backfill the acceptance value (it's a continuous load, not transient).
- 400ms gate sanity: from speech onset detection (AEC-cleaned VAD) to TTS-stop call — the
  ≤400ms budget includes detection (~100ms-class) + stop latency; measure the chain, not just VAD.
- Deliverable: `docs/PRECHECK.md` P15 conclusions + threshold table + go/no-go for Path C
  barge-in and C32 gate condition #1.

## Out of scope
- Product AEC integration wiring (13 does the module, 26/28 consume); C32 provider work (60);
  conversation UX.

## Acceptance criteria
- [ ] License verdict recorded (exact source/build + license file committed to evidence).
- [ ] Go binding builds + links on the pinned toolchain (mingw-w64 per BUILD.md); cgo smoke for
      AEC3 passes; cross-compile status noted (mingw Linux-cross likely NOT viable for this lib —
      record; native Windows build path unaffected).
- [ ] Echo-quality rig: loop-back fixture (render TTS wav via WASAPI + capture via mic/loopback)
      → residual echo measured (ERLE-class metric or perceptual check), self-transcription
      suppression proven (zero ASR text from playback).
- [ ] Cross-endpoint drift data recorded; thresholds or restrictions documented.
- [ ] External-audio false-trigger table (music/speech/video × volumes) with recommended
      barge-in threshold config.
- [ ] Sustained CPU number backfilled into SPEC-10 §3.1 Path C row (replacing the 5% provisional).
- [ ] PRECHECK P15 = PASS/NO-GO recorded; if NO-GO → tickets 26/28 Path-C parts + 60
      auto-defer with DEFERRED registry entries (no silent disappearance).

## Progress log (append-only, newest last)
