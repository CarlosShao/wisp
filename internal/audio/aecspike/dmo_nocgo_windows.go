//go:build windows && !cgo

package aecspike

// No-cgo shape of the DMO probe: the package (and the repo's non-cgo build
// of it) compiles, and every DMO reading reports itself as unavailable
// instead of silently vanishing.

import "fmt"

// DmoInfo is the presence reading for the OS voice-AEC DMO.
type DmoInfo struct {
	Present    bool   `json:"present"`
	HRESULT    string `json:"hresult"`
	InStreams  int    `json:"in_streams"`
	OutStreams int    `json:"out_streams"`
}

// DmoProbe without cgo: unavailable is a reading, not an error.
func DmoProbe() DmoInfo {
	return DmoInfo{Present: false, HRESULT: "cgo-disabled"}
}

// DmoSession is one full-duplex AEC processing session (unavailable here).
type DmoSession struct{}

// DmoBegin without cgo refuses with the reason named.
func DmoBegin() (*DmoSession, error) {
	return nil, fmt.Errorf("aecspike: DMO probe requires cgo (CGO_ENABLED=0 build)")
}

// Process is unreachable on the no-cgo shape.
func (s *DmoSession) Process(far, near []int16, out []int16) (int, error) {
	return 0, fmt.Errorf("aecspike: DMO probe requires cgo")
}

// Close is a no-op on the no-cgo shape.
func (s *DmoSession) Close() {}

// SetSessionVolume without cgo refuses with the reason named.
func SetSessionVolume(client unsafe.Pointer, v float64) (float64, error) {
	return 0, fmt.Errorf("aecspike: session volume requires cgo")
}
