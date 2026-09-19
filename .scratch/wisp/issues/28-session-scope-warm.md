# 28 — SessionScope (C31): Warm/Conversation keep-alive, timers, privacy gating

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 07-ball-state-machine-core, 15-speech-engines-cer-harness
**Parallel slots:** ≤2 sub-agents (A: SessionScope + model refcount + timers; B: Conversation
mode + privacy confirm + battery shrink)
**Spec refs:** SPEC-05 §5, B4, C31, D43 #26,28–31, §15#5 decision, D42#6, S4

## What to build
The `session` module: the single owner of keep-alive — session definition, Warm/Conversation
lifecycles, model reference counting, panel keep-alive reservation, session-scoped grant
registry (rows now, lifecycle at 49), and the settle/dispose cascade.

## Key constraints
- Session = wake → explicit end OR 90s idle (`[session] warm_timeout_sec`); `Settling` 3s static
  period; dispose cascade per D43 #31 (models unload, panel hide-or-destroy, FreeOSMemory,
  grants expired).
- `Warm`: models hot + mic CLOSED + faint warm breathing; next wake via click/hotkey → Listening
  with ZERO model load (assert no engine loads on re-wake — B4's whole point).
- `Conversation` (default OFF via `[voice] conversation_mode`): **Path C full duplex (D47) —
  mic stays OPEN through playback with AEC (ticket 59/P15-gated); barge-in mid-playback →
  Listening (D43 #41)**; red solid ring never fades; Speaking-end → Listening directly (#28);
  30s no-speech → Warm (close mic FIRST); first-ever enable → L2-grade privacy confirm
  (native card) + logged. **Memory note: ASR+TTS co-residency breaks the Path-T serial
  mitigation — Conversation budget measured in S4; approved degradation = suspend ASR during
  playback, keep VAD+AEC for barge-in (ticket 59/15).**
- Battery saver (D42#6): Warm 30s / Conversation 15s; KWS resident OFF.
- Model refcount: ASR/TTS/VAD/punct/KWS acquire-release through SessionScope only (no other
  module may keep models alive); leak test: session end → all refcounts zero.
- Async L1 extraction + history compression run inside the Warm window (hooks for 29/10).
- Interrupt/settle interplay: new wake during Settling cancels settle (#33); explicit end
  (hotkey long-press/menu) → immediate dispose.

## Out of scope
- Panel keep-alive UX (33 wires its window into scope); grant semantics (49); extraction logic (29).

## Acceptance criteria
- [ ] Cold vs Warm wake: second-round first-token latency measured; Warm wake shows zero engine
      loads (instrumentation counter) — B4 evidence.
- [ ] Timers: 90s Warm → Settling → dispose; 30s Conversation → Warm with mic closed FIRST
      (order asserted); battery-saver values applied.
- [ ] **Path C full duplex (D47, blocked by 59): during Conversation playback mic stays open
      (AEC), barge-in works (≤400ms), Conversation memory measured and budget row backfilled
      (or degradation mode verified if over budget).**
- [ ] Conversation first-enable → privacy confirm + log; subsequent enables no prompt.
- [ ] Dispose cascade order test (models → panel-hide → FreeOSMemory → grants expiry) and RSS
      return ≤10s (sampler).
- [ ] Refcount leak test: 100 wake/sleep cycles → zero residual refs/goroutines.

## Progress log (append-only, newest last)
