//go:build windows && cgo && webrtcap

package aec

// The webrtc AEC3 binding proper. wap_shim.cpp (same directory, same build
// constraint) is compiled by cgo and owns every webrtc type; the Go side
// holds only the opaque session pointer, so the webrtc C++ API never leaks
// into Go. CGO_CXXFLAGS must carry the include dirs of the fetched
// webrtc-audio-processing tree and its pinned abseil (deps.toml), and
// CGO_LDFLAGS the -L/-l pair for the built static archive.

/*
#cgo CXXFLAGS: -std=c++17
#cgo LDFLAGS: -lwebrtc_audio_processing -lwinmm -lstdc++
#include "wap_shim.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// NewApm creates the mono AEC3 pipeline at rate (10ms frames). Only the
// echo canceller is enabled: AGC/NS/HPF stay off so the rig measures the
// engine's own residual, not a gain-controlled rescale of it.
func NewApm(rate int) (*Apm, error) {
	if rate <= 0 || rate%100 != 0 {
		return nil, fmt.Errorf("aec: unsupported rate %d", rate)
	}
	var cerr C.int
	sess := C.wap_new(C.int(rate), &cerr)
	if sess == nil {
		return nil, fmt.Errorf("aec: wap_new failed (code=%d)", int(cerr))
	}
	return &Apm{rate: rate, frame: rate / 100, sess: sess}, nil
}

// NewApmTuned is NewApm with the ticket-392 parameter face: the overrides
// in t ride the public EchoControlFactory seam into an EchoCanceller3Config
// at construction. A zero Tuning builds the byte-identical default engine
// (it answers NewApm exactly); a tuned build is reported by the readings
// that carry the tuning (392 AC#2's one-parameter-one-firing discipline).
func NewApmTuned(rate int, t Tuning) (*Apm, error) {
	if !t.Overrides() {
		return NewApm(rate)
	}
	if rate <= 0 || rate%100 != 0 {
		return nil, fmt.Errorf("aec: unsupported rate %d", rate)
	}
	var ct C.wap_aec3_tuning
	ct.dominant_nearend_enr_threshold = C.double(t.DominantNearendEnrThreshold)
	ct.dominant_nearend_trigger_threshold = C.int(t.DominantNearendTriggerThreshold)
	ct.nearend_max_dec_factor_lf = C.double(t.NearendMaxDecFactorLf)
	var cerr C.int
	sess := C.wap_new_tuned(C.int(rate), &ct, &cerr)
	if sess == nil {
		return nil, fmt.Errorf("aec: wap_new_tuned failed (code=%d)", int(cerr))
	}
	return &Apm{rate: rate, frame: rate / 100, sess: sess}, nil
}

// ProcessRender feeds one 10ms far (render) frame. The echo canceller
// buffers the reference internally; call it before the matching capture
// frame of the same tick.
func (a *Apm) ProcessRender(far []float32) error {
	if len(far) != a.frame {
		return fmt.Errorf("aec: render frame len %d, want %d", len(far), a.frame)
	}
	if rc := C.wap_process_reverse(a.sess, (*C.float)(unsafe.Pointer(&far[0])), C.int(a.frame)); rc != 0 {
		return fmt.Errorf("aec: ProcessReverseStream rc=%d: %s", int(rc), a.errString())
	}
	return nil
}

// ProcessCapture feeds one 10ms near (capture) frame with the stream delay
// the rig measured between render and capture, and returns the
// echo-cancelled frame.
func (a *Apm) ProcessCapture(near []float32, delayMs int) ([]float32, error) {
	if len(near) != a.frame {
		return nil, fmt.Errorf("aec: capture frame len %d, want %d", len(near), a.frame)
	}
	out := make([]float32, a.frame)
	rc := C.wap_process_capture(a.sess, (*C.float)(unsafe.Pointer(&near[0])),
		C.int(a.frame), C.int(delayMs), (*C.float)(unsafe.Pointer(&out[0])))
	if rc != 0 {
		return nil, fmt.Errorf("aec: ProcessStream rc=%d: %s", int(rc), a.errString())
	}
	return out, nil
}

// Reset restores the session to its just-created state. This is the
// ticket-393 flush/reset face's reset half, ridden on AudioProcessing's
// public Initialize() face: every internal buffer and the echo canceller
// itself are rebuilt (through the injected EchoControlFactory, so a tuned
// session re-creates the same tuned engine) while the user settings are
// retained - the between-windows lever for the cross-window carry the 392
// findings pinned. A keep-filter flush (drain the latency buffers, keep the
// adapted filter) has no public face in the pinned webrtc tree; that half is
// named absent in docs/evidence/s1/393r2-flush-reset.md.
func (a *Apm) Reset() error {
	if rc := C.wap_reset(a.sess); rc != 0 {
		return fmt.Errorf("aec: Initialize rc=%d: %s", int(rc), a.errString())
	}
	return nil
}

// Close releases the native session. Idempotent; a used-after-Close Apm
// reports an error instead of crashing.
func (a *Apm) Close() error {
	if a.sess == nil {
		return nil
	}
	C.wap_free(a.sess)
	a.sess = nil
	return nil
}

func (a *Apm) errString() string {
	return C.GoString(C.wap_last_error(a.sess))
}
