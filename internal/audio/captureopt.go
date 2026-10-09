package audio

import (
	"log/slog"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Capture-side injection points (ticket 247, orchestrator rulings P4 form 甲
// and P6 form 甲, ledger A485).
//
// Before this file the two real C8 sources booked their capture thread into
// the process registry by name (wasapimic_windows.go:83, wavinjector.go:84),
// which meant the assembly root could not say which registry owns the one
// goroutine D38b put on the resident roster, and there was no way to hand the
// capture leg's per-frame loudness to anybody: the package only ever produced
// PCM frames (AudioSource.Start), and level.go's FrameLevel was a pure
// function with no caller. 247-a4 measured exactly that gap ("no current-level
// accessor; crossing the seam needs a new func(float32) on the capture side").
//
// Both problems are injection points, not new goroutines and not new
// dependency edges: the level is computed INSIDE the existing audio-capture
// goroutine (P8 form 甲, so the frozen six-name resident roster of D38b gains
// nothing) and handed out as one scalar to a callback the caller owns. The
// caller in production is cmd/wisp, which is what keeps the edge count at the
// single new cmd/wisp -> internal/audio edge AC#0 ruled; internal/audio still
// imports neither internal/ball nor internal/statemachine.

// CaptureConfig carries the two capture-side seams one source is built with.
//
// Exactly one registry is ever in play per source: newCaptureConfig seeds the
// field with the process registry (observe.Default - what both sources used
// before this seam existed) and an injected registry REPLACES it. There is no
// path where the same capture thread is booked into two registries, which is
// the false-green shape 247-a1 census E6 refused (P4 form 甲: "no two
// registries coexisting").
type CaptureConfig struct {
	// Registry owns the "audio-capture" goroutine this source spawns.
	Registry *observe.Registry
	// OnLevel receives one level scalar per seam frame, on the capture thread.
	OnLevel func(float32)
}

// CaptureOption configures a CaptureConfig.
type CaptureOption func(*CaptureConfig)

// WithCaptureRegistry books this source's audio-capture goroutine into r
// instead of the process registry. A nil r is ignored and the process
// registry stays: a capture thread that belongs to no roster is a leak the
// watchdog cannot see, so "no registry" is not an available answer.
func WithCaptureRegistry(r *observe.Registry) CaptureOption {
	return func(c *CaptureConfig) {
		if r != nil {
			c.Registry = r
		}
	}
}

// WithLevelSink registers the per-frame level callback.
//
// CONTRACT (AC#3, capability-shaped): fn is called with ONE float32, from the
// audio-capture goroutine, once per C8 seam frame (every FrameDuration, about
// 31 times a second at 16k). It must not block: the capture thread is pinned
// to one OS thread (D38a) and every microsecond it spends inside the sink is a
// microsecond the device buffer is not being drained.
//
// No samples cross this boundary. The frame is turned into the scalar here
// with FrameLevel and the frame itself never leaves the capture loop; the
// render side (internal/ball) is only ever handed a number, which is the
// "no samples and no transcript cross into the ball" rule that file states at
// liquid_windows.go:40-41.
func WithLevelSink(fn func(float32)) CaptureOption {
	return func(c *CaptureConfig) { c.OnLevel = fn }
}

// newCaptureConfig applies opts over the one legal default (the process
// registry) and normalises a cleared registry back to it.
func newCaptureConfig(opts ...CaptureOption) CaptureConfig {
	c := CaptureConfig{Registry: observe.Default}
	for _, o := range opts {
		o(&c)
	}
	if c.Registry == nil {
		c.Registry = observe.Default
	}
	return c
}

// emitLevel turns one seam frame into the scalar the sink wants and hands it
// over. Called from the capture loop only.
//
// A frame that FrameLevel refuses is said out loud and emits nothing: a
// silently missing level would look like silence on the consumer's side, which
// is the same masquerade this package forbids for dropped frames (D38d).
//
// The sink runs even when the frame did not reach the consumer's channel (the
// meter pushed into a full bounded window and counted the drop): the level is
// a fact about the microphone, not about how fast the ASR side reads, and a
// process that has no frame consumer yet (internal/speech is still unwritten)
// must still get its loudness.
func (c CaptureConfig) emitLevel(frame []byte) {
	if c.OnLevel == nil {
		return
	}
	level, err := FrameLevel(frame)
	if err != nil {
		slog.Error("audio level: seam frame rejected, no level emitted", "err", err.Error())
		return
	}
	c.OnLevel(level)
}
