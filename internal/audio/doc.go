// Package audio owns sound input and output plumbing (SPEC-01 §3):
// AudioSource (C8), WASAPI capture on the pinned audio-capture thread (D38a),
// VAD, and the half-duplex switch (D16: keep_audio is hardwired false).
//
// Responsibilities:
//   - device enumeration/open/close, capture and playback buffers
//   - the pinned thread contract: the audio-capture thread runs no other Go
//     code (runtime.LockOSThread)
//   - bounded channel to ASR (<=200ms audio); overflow drops frames AND
//     counts them (D38d) — the counter goes to observe
//   - device hotplug re-enumeration once (D42#2, S2)
//
// Implemented (ticket 13): the C8 AudioSource seam with two real sources -
// WASAPIMicrophone (shared-mode WASAPI, pinned thread, in-process linear
// resampling to 16k/mono/int16, hotplug reopen-once, occupied/permission
// error mapping) and WavInjector (the test backbone; ALL pipeline tests
// drive data through it, SPEC-04 §9) - plus HalfDuplexGate (D16, Path T/C
// per D47) and the bounded metered frame channel. The capture/render
// endpoint pair stays queryable for the P15 AEC spike (D47).
//
// Non-responsibilities:
//   - no ASR/TTS/KWS model inference (speech)
//   - no audio device policy decisions (risk/permission), no UI
//   - no AEC source (Path C AEC is tickets 26/59; the gate interface carries
//     the Path switch, not the processing)
//
// DEFERRED(playback): TTS output lands with ticket 26 on the same WASAPI
// technique. Audio buffers are never persisted or logged (D16③).
package audio
