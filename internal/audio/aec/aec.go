// Package aec is the cgo seam over the webrtc AEC3 engine
// (webrtc-audio-processing, pinned in deps.toml) used by the aecspike
// re-measurement rig: far (render) reference and near (capture) 10ms frames
// in, echo-cancelled capture out.
//
// The native build is NOT part of the default repo build. The binding
// compiles only under `-tags webrtcap` on Windows with cgo, with
// CGO_CXXFLAGS/CGO_LDFLAGS pointing at the fetched-and-built tree recorded
// in deps.toml (the reproduction recipe lives in
// docs/evidence/s1/aec3r1-remeasure.md). Everywhere else - other GOOSes, no
// cgo, no tag - NewApm returns ErrUnavailable, so nothing outside the gated
// spike test ever reaches for the native library.
//
// An Apm instance processes one mono stream pair and is not safe for
// concurrent use: the webrtc APM pins its stream state to the calling
// sequence render-frame then capture-frame, exactly the offline lockstep the
// rig drives. This is the spike-grade binding; the production AEC seam is
// ticket 13/26/28 work, not this package.
package aec

import (
	"errors"
	"unsafe"
)

// ErrUnavailable is returned by NewApm when the native webrtc
// audio-processing library is not part of the build (no cgo, no webrtcap
// tag, or the library is absent).
var ErrUnavailable = errors.New("aec: webrtc-audio-processing binding not built (requires cgo + -tags webrtcap; see deps.toml [webrtc-audio-processing])")

// Tuning is the AEC3 initialization-parameter face (ticket 392 AC#2): the
// double-talk suppressor knobs the binding exposes at construction time.
// Every field is an override - the zero value keeps the engine's own default
// (the exact pre-392 posture), so a zero Tuning builds the identical engine
// NewApm builds. The parameters ride the public EchoControlFactory seam
// (AudioProcessingBuilder.SetEchoControlFactory) with an EchoCanceller3Config
// - no webrtc source change, the C++ tree stays pinned and untouched.
type Tuning struct {
	// DominantNearendEnrThreshold overrides suppressor.
	// dominant_nearend_detection.enr_threshold (engine default 0.25): the
	// echo-to-nearend ratio under which the nearend-dominant state may
	// engage. The nearend state is the transparent suppressor tuning - the
	// state where the engine stops eating the near talker - so raising it is
	// the double-talk damage lever: the state engages while the talker and
	// the echo are still comparable instead of only when the talker
	// dominates outright.
	DominantNearendEnrThreshold float64
	// DominantNearendTriggerThreshold overrides suppressor.
	// dominant_nearend_detection.trigger_threshold (engine default 12, i.e.
	// 12 consecutive 64-sample blocks ~= 48ms of evidence): how long the
	// nearend evidence must hold before the state engages. Lower engages
	// faster inside a barge-in.
	DominantNearendTriggerThreshold int
	// NearendMaxDecFactorLf overrides suppressor.nearend_tuning.
	// max_dec_factor_lf (engine default 0.25): the per-frame drop floor of
	// the low-band suppression gain while the nearend state holds. Higher
	// makes the dive into deep suppression slower after strong nearend, so
	// the talker's first syllables survive more of the engine's own
	// suppression.
	NearendMaxDecFactorLf float64
}

// Overrides reports whether any field carries a non-default instruction
// (all three knobs are positive quantities, so zero means "unset").
func (t Tuning) Overrides() bool {
	return t.DominantNearendEnrThreshold > 0 ||
		t.DominantNearendTriggerThreshold > 0 ||
		t.NearendMaxDecFactorLf > 0
}

// Apm is one mono AEC3 pipeline at a fixed sample rate. A zero Apm is not
// usable; obtain one from NewApm.
type Apm struct {
	rate  int
	frame int
	sess  unsafe.Pointer
}

// The fields are the binding's private state: every reader lives in the
// windows && cgo && webrtcap binding (aec_webrtc_windows.go), which the linux
// build cannot see. This reference keeps the seam struct whole per platform.
var _ = [3]any{Apm{}.rate, Apm{}.frame, Apm{}.sess}

// FrameSamples is the 10ms frame length at the given rate, the chunk size
// ProcessRender and ProcessCapture accept.
func FrameSamples(rate int) int { return rate / 100 }
