# 13 — Audio capture stack: C8 sources, WASAPI pinned thread, hotplug, permissions (DONE ✅)

**Status:** done
**Claimed by:** orchestrator -> sub-agent T13-impl
**Last update:** 2026-09-19T15:11:33ZT13:56:41Z
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
- [x] WavInjector → channel: exact frames, backpressure drop counter increments under a slow
      consumer, counter visible in logs.
- [x] Pinned-thread test: capture goroutine OS thread stable across 10s (no M-migration);
      no other work scheduled on it.
- [x] Resampler test: 48k fixture → 16k output, SNR sanity + latency bounded.
- [x] Hotplug: injected default-device-change → re-enum; stale-handle scenario → error with
      device name (not silent).
- [x] Occupied/permission fixtures → `Error(audio_device)` with guidance text.
- [x] Half-duplex gate test: gate closed → zero frames delivered; reopen restores.

## Progress log (append-only, newest last)
- [2026-09-19T13:56:41Z] agent=orchestrator claimed=T13-impl did=dispatched next=sub-agent works through acceptance criteria
- [2026-09-19T14:16:44Z] agent=T13-impl did=c8-seam+resampler+wavinjector (AudioSource C8 seam + 16k/mono/int16 512-sample frame contract + bounded <=192ms channel w/ metered push: drop+count+rate-limited warn log (D38d); streaming fixed-point linear resampler 48k/44.1k->16k, chunk-invariant + 1-sample latency bound + SNR tested; WavInjector: RIFF PCM16/float32 parse -> mono downmix -> 16k, frame-exact replay, realtime/flood pacing, backpressure drop counter visible in logs; race-clean) next=half-duplex-gate
- [2026-09-19T14:27:25Z] agent=T13-impl did=half-duplex-gate-D16 (HalfDuplexGate implements AudioSource: Path T closes capture on SetSpeaking (inner device Stop, proven frozen FramesSent, not discard-only) and reopens on speech end; Path C (D47) full-duplex pass-through with speaking events only (AEC=26/59); SetMuted closes on both paths + Muted/Unmuted events for state 07 hotkey path; WithStartMuted maps [audio] mic_muted_default at boot; residual <=192ms pre-TTS frames left to consumer (not self-echo); tests: zero-frames-while-closed, reopen restores, Path C untouched, muted-wins-over-speaking, stop idempotent; race-clean) next=wasapi-capture+hotplug
- [2026-09-19T14:56:21Z] agent=T13-impl did=wasapi-capture+hotplug+errors (WASAPIMicrophone: shared-mode WASAPI via x/sys COM vtables (ball/spike-proven pattern, no cgo), runtime.LockOSThread pinned audio-capture loop spawned via observe roster, device-native float/int16 mix format -> MonoDownmix -> fixed-point resampler -> 512-sample frames -> metered bounded push; hotplug D42#2: Go-synthesized IMMNotificationClient (global-registry callback identity) -> injected watcher event -> reopen-once w/ aligned-Initialize retry; reopen failure or stale handle = terminal Error(audio_device) NAMING the device (prev-device name on enum failure), never silent; D42#12 mapping: DEVICE_IN_USE->exclusive-mode guidance, E_ACCESSDENIED->Windows privacy-settings text, fixture-tested; Endpoints() exposes capture/render pair for P15/D47; tests: hotplug re-enum-once, stale-handle loud failure, enum-failure naming, occupied/permission fixtures, wav-integrity through real loop, 10s pinned-thread stability (1 OS thread, no foreign goroutine on it), live smoke; LIVE-VERIFIED on Realtek mic: enumeration w/ friendly names + sustained capture 25 frames/820ms == 32ms cadence; note: IAudioClient vtable slot correction (Start=10 Stop=11 SetEventHandle=13) + AUDCLNT_E_BUFFER_SIZE_NOT_ALIGNED aligned-retry found during live run) next=final-gates+handoff
- [2026-09-19T14:58:45Z] agent=T13-impl did=handoff-to-orchestrator (all 6 acceptance criteria implemented+tested: wav frames exact+drop counter in stats/log; 10s pinned-thread stability; 48k/44.1k->16k resampler length/SNR/latency; hotplug inject->re-enum-once + stale-handle named error; occupied/permission fixtures->Error(audio_device)+guidance; half-duplex gate zero-frames+reopen. Gates: go vet ./... green, go build ./... green, go test ./... green (audio 21 tests race-clean), scripts/build.ps1 -Env dev green (doctor PASS, commit 77538f1 baked). LIVE smoke on real Realtek mic PASS (WISP_LIVE_MIC=1 hook retained for ticket 16). Parallel note: T09-impl shares this worktree; all commits used explicit paths, no ticket marked done; handoff notes: wiring points = NewWASAPIMicrophone->HalfDuplexGate->NewBoundedFrames consumer, SetSpeaking/SetMuted + WithStartMuted(mic_muted_default), events to state 07) next=ticket-done-decision-belongs-to-orchestrator
- [2026-09-19T15:11:33Z] agent=orchestrator did=T13-adv VERDICT PASS (0 BLOCKER/0 MAJOR/3 MINOR; report docs/evidence/s2/13-adversarial-acceptance.md; 26/26 race green + live-mic reproduced; MINORs are undercounts in handoff claims) next=ticket DONE
- [2026-09-20T02:25Z] agent=agent-bookkeeping-1 did=AC boxes reconciled against docs/evidence/s2/13-adversarial-acceptance.md (§6 六条验收逐条 table: each AC ruled PASS with named tests + file:line, independently re-run under -race): 6 checked, 0 left open next=none
