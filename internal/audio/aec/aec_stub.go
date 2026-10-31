//go:build !windows || !cgo || !webrtcap

package aec

// The no-native-build shape: the spike rig skips its AEC3 reading with this
// error, every other caller sees the same honest answer, and nothing tries
// to link a library that is not there. The method set mirrors the real
// binding so callers compile everywhere; NewApm fails first, so the methods
// are unreachable in this configuration.

// NewApm always fails with ErrUnavailable in this build configuration.
func NewApm(rate int) (*Apm, error) { return nil, ErrUnavailable }

// ProcessRender always fails with ErrUnavailable in this build configuration.
func (a *Apm) ProcessRender(far []float32) error { return ErrUnavailable }

// ProcessCapture always fails with ErrUnavailable in this build configuration.
func (a *Apm) ProcessCapture(near []float32, delayMs int) ([]float32, error) {
	return nil, ErrUnavailable
}

// Close always fails with ErrUnavailable in this build configuration.
func (a *Apm) Close() error { return ErrUnavailable }
