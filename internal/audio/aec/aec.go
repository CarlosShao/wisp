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

// Apm is one mono AEC3 pipeline at a fixed sample rate. A zero Apm is not
// usable; obtain one from NewApm.
type Apm struct {
	rate  int
	frame int
	sess  unsafe.Pointer
}

// FrameSamples is the 10ms frame length at the given rate, the chunk size
// ProcessRender and ProcessCapture accept.
func FrameSamples(rate int) int { return rate / 100 }
