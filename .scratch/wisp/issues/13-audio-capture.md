# 13 — Audio capture stack: C8 sources, WASAPI pinned thread, hotplug, permissions

**Status:** ready-for-agent
**Claimed by:** —
**Last update:** 2026-09-19
**Blocked by:** 03-skeleton-runtime-rules
**Parallel slots:** ≤2 sub-agents (A: WASAPI source + resampler + pinned thread; B: WavInjector
+ hotplug + error mapping + tests)
**Spec refs:** SPEC-04 §2–§3, C8, D16, D38(a)(d), D42#2#12, S2

## What to build
`audio` module: the C8 `AudioSource` seam with two real implementations — WASAPI microphone
(share mode, pinned thread) and WavInjector (test backbone) — plus in-process resampling to
16k/mono, bounded transfer channel with drop+count telemetry, half-duplex gate, device hotplug
handling, and mic-occupied/permission-denied error mapping.

## Key constraints
- Interface per C8: `Start(ctx, chan<- []byte)/Stop`; frames 16kHz/mono/int16; implementations:
  WASAPIMicrophone, WavInjector, and **AecSource — upgraded from interface-slot to a real S4
  implementation (D47): wraps webrtc-audio-processing (AEC3, P15-gated) using self-rendered TTS
  PCM as reference + WASAPI render position; used ONLY on Path C (Conversation)**. Path T
  remains half-duplex (design, not debt).
- WASAPI: shared mode; `runtime.LockOSThread()` pinned capture thread that runs NO other Go code
  (D38a); resample device-native rate → 16k via in-module linear interpolator (no new deps).
- Bounded channel ≤200ms of audio; overflow drops frames AND increments a counter exposed to
  logs/diagnostics (silent drops are forbidden — they masquerade as accuracy loss).
- Half-duplex gate (D16): a mute gate the TTS side will drive (Speaking closes capture); audio
  buffers NEVER persisted/logged (D16③).
- Hotplug (D42#2): `IMMNotificationClient::OnDefaultDeviceChanged` → re-enumerate once; failure →
  `Error(audio_device)` naming the device; never silently keep a stale handle. **D47 note: capture
  and render may be DIFFERENT endpoints (Bluetooth speaker + wired mic) — clock drift between them
  is an AEC quality input; P15 (ticket 59) measures this case; keep capture/render endpoint pair
  queryable for the spike.**
- Occupied/permission-denied (D42#12): open-device failures → `Error(audio_device)` with explicit
  guidance text pointing to Windows privacy settings; silent failure forbidden.
- Mic-muted-default respected on boot (`[audio] mic_muted_default`); mute hotkey path hooks into
  Muted state (07) via events.
- Tests drive ONLY through WavInjector for data; device enumeration may be faked at the
  enumeration interface (not at the audio-data seam).

## Out of scope
- VAD/ASR consumption (15); KWS ring-buffer (41); Conversation-mode mic policy (28).

## Acceptance criteria
- [ ] WavInjector → channel: exact frames, backpressure drop counter increments under a slow
      consumer, counter visible in logs.
- [ ] Pinned-thread test: capture goroutine OS thread stable across 10s (no M-migration);
      no other work scheduled on it.
- [ ] Resampler test: 48k fixture → 16k output, SNR sanity + latency bounded.
- [ ] Hotplug: injected default-device-change → re-enum; stale-handle scenario → error with
      device name (not silent).
- [ ] Occupied/permission fixtures → `Error(audio_device)` with guidance text.
- [ ] Half-duplex gate test: gate closed → zero frames delivered; reopen restores.

## Progress log (append-only, newest last)
