package audio

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// WavInjector is the C8 test backbone (SPEC-04 sec 3/9): it replays a WAV
// file as 16kHz/mono/int16 frames through the exact same bounded-channel
// push, drop counting and frame pacing contract as the real microphone, so
// the whole voice pipeline is testable without audio hardware ("all tests
// drive through C8 WavInjector" is the foundation of the test strategy).
//
// Frame contract: samples are chunked into FrameSamples frames; the final
// frame of a finite file is zero-padded to FrameBytes. Pacing is monotonic
// (frame i due at start + i*FrameDuration), never wall-clock based (D42#9).
type WavInjector struct {
	samples []int16 // 16k mono, the seam rate
	pace    time.Duration
	meter   *meter

	mu     sync.Mutex
	cancel context.CancelFunc
	done   chan struct{}
}

// InjectorOption configures a WavInjector.
type InjectorOption func(*WavInjector)

// WithInjectorPace sets the inter-frame delay. The default is FrameDuration
// (real-time replay). A pace of 0 floods (no pacing), which is how back-
// pressure drops are driven deterministically in tests.
func WithInjectorPace(d time.Duration) InjectorOption {
	return func(w *WavInjector) { w.pace = d }
}

// NewWavInjector parses a RIFF/WAVE file (PCM16 or IEEE float32, any channel
// count and sample rate), converts it to the 16k/mono seam rate in-process
// (same MonoDownmix + linear resampler as the microphone path) and returns a
// reusable injector. Parse errors fail fast at construction.
func NewWavInjector(path string, opts ...InjectorOption) (*WavInjector, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	samples, rate, chans, err := parseWav(data)
	if err != nil {
		return nil, fmt.Errorf("wav %s: %w", path, err)
	}
	samples = MonoDownmix(samples, chans)
	if rate != TargetRate {
		samples = ResampleLinear(samples, rate, TargetRate)
	}
	w := &WavInjector{
		samples: samples,
		pace:    FrameDuration,
		meter:   newMeter("wav-injector"),
		done:    make(chan struct{}),
	}
	for _, o := range opts {
		o(w)
	}
	return w, nil
}

// Start implements AudioSource. It returns after the pump is running; frames
// begin flowing into buf immediately after.
func (w *WavInjector) Start(ctx context.Context, buf chan<- []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel != nil {
		return fmt.Errorf("wav-injector: already started")
	}
	cctx, cancel := context.WithCancel(ctx)
	w.cancel = cancel
	w.done = make(chan struct{})
	observe.Default.Spawn("audio-capture", "audio", nil, func(ctx context.Context) {
		w.pump(cctx, buf)
		close(w.done)
	})
	return nil
}

// Stop implements AudioSource: cancels the pump and waits for it to exit.
// Idempotent; safe without Start.
func (w *WavInjector) Stop() error {
	w.mu.Lock()
	cancel := w.cancel
	done := w.done
	w.cancel = nil
	w.mu.Unlock()
	if cancel != nil {
		cancel()
		<-done
	}
	return nil
}

// Stats implements Stater (D38d telemetry).
func (w *WavInjector) Stats() Stats { return w.meter.snapshot() }

// pump replays the sample stream: pace first, push second (so the first
// frame is immediate and the cadence is start + i*FrameDuration).
func (w *WavInjector) pump(ctx context.Context, buf chan<- []byte) {
	start := time.Now()
	for off := 0; off < len(w.samples); off += FrameSamples {
		if w.pace > 0 {
			target := start.Add(time.Duration(off/FrameSamples+1) * w.pace)
			if d := time.Until(target); d > 0 { // monotonic (D42#9)
				select {
				case <-ctx.Done():
					return
				case <-time.After(d):
				}
			}
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
		w.meter.push(buf, encodeFrame(w.samples[off:min(off+FrameSamples, len(w.samples))]))
	}
}

// encodeFrame little-endian-encodes samples into a full FrameBytes frame,
// zero-padding the tail (finite-source contract, see the type comment).
func encodeFrame(samples []int16) []byte {
	frame := make([]byte, FrameBytes)
	for i, s := range samples {
		frame[2*i] = byte(s)
		frame[2*i+1] = byte(s >> 8)
	}
	return frame
}

// parseWav parses the minimal RIFF/WAVE subset the injector needs: PCM16
// (format 1) and IEEE float32 (format 3, incl. WAVE_FORMAT_EXTENSIBLE with a
// float SubFormat). Anything else is a loud error, never a silent mangle.
func parseWav(data []byte) (samples []int16, rate, chans int, err error) {
	if len(data) < 12 || string(data[0:4]) != "RIFF" || string(data[8:12]) != "WAVE" {
		return nil, 0, 0, fmt.Errorf("not a RIFF/WAVE file")
	}
	var (
		fmtTag    uint16
		bits      uint16
		fmtSeen   bool
		dataChunk []byte
	)
	off := 12
	for off+8 <= len(data) {
		id := string(data[off : off+4])
		size := int(binary.LittleEndian.Uint32(data[off+4 : off+8]))
		body := off + 8
		if body+size > len(data) {
			size = len(data) - body // tolerate a truncated final chunk
		}
		switch id {
		case "fmt ":
			if size < 16 {
				return nil, 0, 0, fmt.Errorf("fmt chunk too short (%d)", size)
			}
			fmtTag = binary.LittleEndian.Uint16(data[body : body+2])
			chans = int(binary.LittleEndian.Uint16(data[body+2 : body+4]))
			rate = int(binary.LittleEndian.Uint32(data[body+4 : body+8]))
			bits = binary.LittleEndian.Uint16(data[body+14 : body+16])
			if fmtTag == 0xFFFE { // WAVE_FORMAT_EXTENSIBLE: real tag is SubFormat[0:2]
				if size < 40 {
					return nil, 0, 0, fmt.Errorf("extensible fmt chunk too short (%d)", size)
				}
				fmtTag = binary.LittleEndian.Uint16(data[body+24 : body+26])
			}
			fmtSeen = true
		case "data":
			dataChunk = data[body : body+size]
		}
		off = body + size + (size & 1) // chunks are word-aligned
	}
	if !fmtSeen {
		return nil, 0, 0, fmt.Errorf("missing fmt chunk")
	}
	if chans <= 0 || rate <= 0 {
		return nil, 0, 0, fmt.Errorf("bad format: rate=%d chans=%d", rate, chans)
	}
	switch {
	case fmtTag == 1 && bits == 16:
		n := len(dataChunk) / 2
		samples = make([]int16, n)
		for i := 0; i < n; i++ {
			samples[i] = int16(binary.LittleEndian.Uint16(dataChunk[2*i : 2*i+2]))
		}
		return samples, rate, chans, nil
	case fmtTag == 3 && bits == 32:
		n := len(dataChunk) / 4
		f := make([]float32, n)
		for i := 0; i < n; i++ {
			f[i] = math.Float32frombits(binary.LittleEndian.Uint32(dataChunk[4*i : 4*i+4]))
		}
		return FloatToPCM16(f), rate, chans, nil
	default:
		return nil, 0, 0, fmt.Errorf("unsupported wav format: tag=%d bits=%d (want PCM16 or float32)", fmtTag, bits)
	}
}
