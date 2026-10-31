//go:build windows && cgo

package aecspike

// The cgo surface of the AEC DMO probe. This file IS the ticket 59 AC#2
// cgo smoke: it is compiled and linked by the pinned MSYS2 mingw-w64 GCC
// (BUILD.md) through cgo, and its probe call is exercised at build-verified
// runtime by the gated tests. See aecsmoke.h for the C contract; the
// session calls must stay on ONE pinned OS thread (the rig thread or the
// test's LockOSThread).

/*
#cgo windows LDFLAGS: -lole32
#include "aecsmoke.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// DmoInfo is the presence reading for the OS voice-AEC DMO.
type DmoInfo struct {
	Present    bool   `json:"present"`
	HRESULT    string `json:"hresult"`
	InStreams  int    `json:"in_streams"`
	OutStreams int    `json:"out_streams"`
}

// DmoProbe CoCreateInstances the Voice Capture DSP and reports its stream
// layout. HRESULT is rendered 0x%08X (S_OK = present).
func DmoProbe() DmoInfo {
	var in, out C.int32_t
	hr := C.aecsmoke_probe(&in, &out)
	info := DmoInfo{
		HRESULT:    fmt.Sprintf("0x%08X", uint32(hr)),
		InStreams:  int(in),
		OutStreams: int(out),
	}
	info.Present = hr == 0
	return info
}

// DmoSession is one full-duplex AEC processing session. Process-global C
// state: one session per process, owned by the pinned thread that called
// DmoBegin.
type DmoSession struct{ active bool }

// DmoBegin starts a session (feed mode, AEC system mode, 16k/mono/s16).
func DmoBegin() (*DmoSession, error) {
	hr := C.aecsmoke_begin()
	if hr != 0 {
		return nil, fmt.Errorf("aecspike: aecsmoke_begin hr=0x%08X", uint32(hr))
	}
	return &DmoSession{active: true}, nil
}

// Process feeds one aligned chunk (far reference + near signal, up to 1600
// samples each) and copies the produced output samples into out. Returns
// the produced sample count; 0 with nil error means the engine is still
// buffering (S_FALSE) and the caller should feed the next chunk.
func (s *DmoSession) Process(far, near []int16, out []int16) (int, error) {
	if !s.active {
		return 0, fmt.Errorf("aecspike: AEC session is closed")
	}
	n := len(far)
	if len(near) < n {
		n = len(near)
	}
	if len(out) < n {
		return 0, fmt.Errorf("aecspike: output buffer too small (%d < %d)", len(out), n)
	}
	var produced C.int32_t
	hr := C.aecsmoke_process((*C.int16_t)(unsafe.Pointer(&far[0])),
		(*C.int16_t)(unsafe.Pointer(&near[0])),
		(*C.int16_t)(unsafe.Pointer(&out[0])), C.int32_t(n), &produced)
	switch {
	case hr == 0:
		return int(produced), nil
	case hr == 1:
		return 0, nil
	default:
		return 0, fmt.Errorf("aecspike: aecsmoke_process hr=0x%08X", uint32(hr))
	}
}

// SetSessionVolume pins the rig render stream's session volume through a
// real C call (float arguments ride XMM registers that syscall.SyscallN
// cannot populate; see aecsmoke.c) and reports the read-back value.
func SetSessionVolume(client unsafe.Pointer, v float64) (float64, error) {
	hr := C.aecsmoke_set_client_volume(client, C.float(v))
	if hr != 0 {
		return 0, fmt.Errorf("aecspike: aecsmoke_set_client_volume hr=0x%08X", uint32(hr))
	}
	return float64(C.aecsmoke_last_volume()), nil
}

// Close ends the session (balances COM init). Idempotent.
func (s *DmoSession) Close() {
	if s.active {
		C.aecsmoke_end()
		s.active = false
	}
}
