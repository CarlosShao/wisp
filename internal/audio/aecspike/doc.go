// Package aecspike is the ticket 59 P15 AEC feasibility spike: the decision
// readings that gate Path C full duplex (D47), not a production AEC.
//
// What this package is:
//
//   - a loop-back rig (WASAPI render + endpoint loopback capture) with the
//     reference-subtraction math that turns captured audio into ERLE-class
//     readings, same-endpoint drift data and a synthetic clock-drift
//     sensitivity curve;
//   - a probe of the zero-download AEC candidate on this machine (the
//     Windows Voice Capture DSP, CLSID_CWMAudioAEC, driven full duplex from
//     cgo) - the candidate the license verdict cell (ticket 59 AC#1) rules
//     on, because the ticket-named webrtc-audio-processing build cannot be
//     evaluated without fetching it;
//   - the false-trigger table driver: the product's own VAD (internal/speech
//     Detector, silero + Endpointer) run over raw vs echo-cancelled captures
//     of music/speech/video fixtures at several playback volumes;
//   - a self-transcription suppression check (the product's streaming ASR on
//     the raw playback capture vs the residual).
//
// What this package is NOT: production wiring. Nothing here is referenced by
// the resident process; the AEC seam integration belongs to ticket 13/26/28.
// Nothing here changes HalfDuplexGate semantics, D43, SPEC text, config keys
// or cmd/wisp.
//
// Privacy (D16): the rig only ever renders and captures SYNTHETIC fixtures
// and self-synthesised TTS from the pinned third_party models. Audio
// artifacts are written ONLY when the WISP_AEC_SPIKE=1 gate is set, ONLY to
// the workspace temp root the dispatcher named (never into the repository,
// never into user data directories), and the committed evidence carries
// numbers, never payloads.
package aecspike
