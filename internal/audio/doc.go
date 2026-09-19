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
// Non-responsibilities:
//   - no ASR/TTS/KWS model inference (speech)
//   - no audio device policy decisions (risk/permission), no UI
//
// DEFERRED(capture/vad/playback): implemented by ticket 13 (capture) with
// TTS output in ticket 26. This ticket only freezes the package boundary.
package audio
