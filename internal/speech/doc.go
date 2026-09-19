// Package speech owns the voice model engines (SPEC-01 §3): AsrEngine,
// TtsEngine and WakeWordEngine (C9) over the sherpa-onnx runtime, plus the
// interface seams for cloud providers (C9, ticket 61).
//
// Responsibilities:
//   - engine/session lifecycle: load, infer, release (inside a DisposalScope,
//     C11 — sessions are released, not unloaded, since FreeLibrary of
//     onnxruntime is not reliable)
//   - cgo boundary discipline (D37b): C return codes and char* error strings
//     are copied to Go strings immediately, never held
//   - inference runs in disposable goroutines or (path X) a speech subprocess;
//     cancellation means abandoning the result, not interrupting onnxruntime
//     Run (D37c)
//
// Non-responsibilities:
//   - no audio capture/playback (audio), no model download/verification
//     (models), no wake-word arming policy (statemachine/watchdog)
//
// DEFERRED(engines): implemented by ticket 15 (ASR + CER harness), ticket 26
// (TTS), ticket 41 (KWS). This ticket only freezes the package boundary.
package speech
