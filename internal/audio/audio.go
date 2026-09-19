package audio

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// C8 AudioSource seam (SPEC-04 sec 3) plus the shared frame data path:
// 16kHz/mono/int16 frames, the <=200ms bounded transfer channel discipline
// (D38d: overflow drops frames AND counts them - silent drops are forbidden
// because they masquerade as ASR accuracy loss), and the hotplug/diagnostics
// stats both real implementations share.
//
// D16 privacy rule: audio buffers are NEVER persisted and NEVER logged. The
// logging in this package records counts, durations and device names only -
// never frame payloads. keep_audio is hardwired false (SPEC-03 sec 3).

// AudioSource is the C8 seam: a producer of 16kHz/mono/int16 frames.
type AudioSource interface {
	// Start begins pushing frames into buf after Start returns. Frames are
	// 16kHz/mono/int16, FrameSamples samples (FrameBytes bytes) each, except
	// a possibly zero-padded final frame of a finite source. ctx cancellation
	// stops the source; buf is owned by the consumer and is NOT closed by
	// the source. Start must not block on the consumer.
	Start(ctx context.Context, buf chan<- []byte) error
	// Stop stops the source and releases device resources. Stop is
	// idempotent and safe to call without Start.
	Stop() error
}

// Stater is implemented by sources that expose D38d telemetry.
type Stater interface {
	// Stats returns a snapshot of the source's counters and last error.
	Stats() Stats
}

// Target sample rate of the C8 seam: everything downstream (VAD/ASR/KWS)
// consumes 16kHz mono (SPEC-04 sec 3).
const TargetRate = 16000

// FrameSamples is the frame length the seam guarantees, aligned to the VAD
// input (SPEC-04 sec 3: 512 samples @16k = 32ms).
const FrameSamples = 512

// FrameBytes is the byte length of one full frame (int16 = 2 bytes/sample).
const FrameBytes = FrameSamples * 2

// FrameDuration is the wall time of one full frame.
const FrameDuration = FrameSamples * time.Second / TargetRate // 32ms

// BoundedWindow is the maximum audio the transfer channel may hold (D38d:
// <=200ms). Six frames of 32ms = 192ms is the largest whole-frame window
// under the cap.
const BoundedWindow = 6 * FrameDuration // 192ms

// BoundedFrameCapacity is the channel capacity implementing BoundedWindow.
func BoundedFrameCapacity() int { return int(BoundedWindow / FrameDuration) }

// NewBoundedFrames creates the bounded audio->ASR transfer channel (D38d).
// Sources drop (never block, never grow it) and count on overflow.
func NewBoundedFrames() chan []byte { return make(chan []byte, BoundedFrameCapacity()) }

// Path is the voice path a gate serves (D47): T = task/work path, kept
// half-duplex by design (D16); C = conversation path, full duplex with AEC
// (AEC itself is ticket 26/59 - the gate only stops policing it).
type Path string

const (
	// PathT is the half-duplex work path (D16): speaking closes capture.
	PathT Path = "T"
	// PathC is the conversation path (D47): capture stays open while
	// speaking; echo is the AEC layer's problem, not the gate's.
	PathC Path = "C"
)

// Stats is a diagnostics snapshot (D38d: drop counters must reach logs and
// the diagnostics package; D42#2/#12: device errors keep the device name).
type Stats struct {
	FramesSent    uint64 `json:"frames_sent"`
	FramesDropped uint64 `json:"frames_dropped"`
	BytesSent     uint64 `json:"bytes_sent"`
	BytesDropped  uint64 `json:"bytes_dropped"`
	Reopens       uint64 `json:"reopens"`  // hotplug re-enumerations performed
	LastError     string `json:"last_error,omitempty"`
}

// Dropped reports whether any frame was ever dropped.
func (s Stats) Dropped() bool { return s.FramesDropped > 0 }

// meter carries the D38d counters of one source and does the metered push
// into the consumer's bounded channel: non-blocking send; a full buffer
// drops the frame, increments the counters and rate-limited-logs. The log
// carries counts only (D16: frame payloads never reach logs).
type meter struct {
	name string // source name for log lines

	mu       sync.Mutex
	sent     uint64
	dropped  uint64
	bytesSent uint64
	bytesDropped uint64
	reopens  uint64
	lastErr  string
	lastWarn time.Time // monotonic reading for the 1s log throttle (D42#9)
}

func newMeter(name string) *meter { return &meter{name: name} }

// push sends one frame without ever blocking on the consumer. Returns true
// when delivered, false when dropped (D38d: drop + count + visible log).
func (m *meter) push(buf chan<- []byte, frame []byte) bool {
	n := len(frame)
	select {
	case buf <- frame:
		m.mu.Lock()
		m.sent++
		m.bytesSent += uint64(n)
		m.mu.Unlock()
		return true
	default:
		m.mu.Lock()
		m.dropped++
		m.bytesDropped += uint64(n)
		total := m.dropped
		droppedBytes := m.bytesDropped
		logIt := total == 1 || time.Since(m.lastWarn) >= time.Second
		if logIt {
			m.lastWarn = time.Now()
		}
		m.mu.Unlock()
		if logIt {
			slog.Warn("audio frames dropped: bounded channel full (slow consumer)",
				"source", m.name,
				"dropped_frames_total", total,
				"dropped_bytes_total", droppedBytes,
				"bounded_window", BoundedWindow.String(),
				"frame_ms", FrameDuration.Milliseconds())
		}
		return false
	}
}

// fail records a terminal device error (D42#2/#12: device name inside).
func (m *meter) fail(err error) {
	m.mu.Lock()
	m.lastErr = err.Error()
	m.mu.Unlock()
	slog.Error("audio source failed", "source", m.name, "error", err.Error())
}

// reopened records a hotplug re-enumeration (D42#2).
func (m *meter) reopened() {
	m.mu.Lock()
	m.reopens++
	m.mu.Unlock()
}

// snapshot builds the Stats view.
func (m *meter) snapshot() Stats {
	m.mu.Lock()
	defer m.mu.Unlock()
	return Stats{
		FramesSent:    m.sent,
		FramesDropped: m.dropped,
		BytesSent:     m.bytesSent,
		BytesDropped:  m.bytesDropped,
		Reopens:       m.reopens,
		LastError:     m.lastErr,
	}
}

// SpawnCapture spawns the audio-capture goroutine (the D38b roster name) as
// the single sanctioned entry point for capture threads. fn runs pinned to
// one OS thread (D38a) - see runPinned in the Windows capture file. Here it
// only exists so tests on any platform can exercise the same spawn seam.
func SpawnCapture(registry *observe.Registry, run func(ctx context.Context)) *observe.Handle {
	return registry.Spawn("audio-capture", "audio", nil, run)
}
