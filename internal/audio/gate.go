package audio

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// Half-duplex gate (D16): the seam the TTS side and the mute hotkey drive.
//
//   - Path T (work path): while TTS is speaking the gate CLOSES capture -
//     the inner source is Stop()ed (device off, D16 "关闭麦克风", not a
//     discard-only filter); on speech end it reopens with zero model reload.
//   - Path C (conversation path, D47): full duplex - speaking does NOT close
//     capture; echo cancellation is the AEC layer's job (tickets 26/59). The
//     gate still reports speaking events so downstream layers can adapt.
//   - Muted (both paths, user intent wins): capture closed; Muted/Unmuted
//     events hook the mute hotkey path into the Muted state (ticket 07).
//
// Frames already sitting in the bounded channel when the gate closes were
// captured BEFORE the TTS started (<= 192ms, D38d) and are not self-echo;
// the gate leaves the consumer's channel alone. Audio buffers are never
// persisted or logged here (D16③).
//
// The gate itself implements AudioSource, so callers wire
// [WASAPIMicrophone | WavInjector] -> Gate -> consumer unchanged.

// GateEventKind enumerates the gate's state transitions.
type GateEventKind string

const (
	EventCaptureOpened   GateEventKind = "capture_opened"
	EventCaptureClosed   GateEventKind = "capture_closed"
	EventMuted           GateEventKind = "muted"
	EventUnmuted         GateEventKind = "unmuted"
	EventSpeakingStarted GateEventKind = "speaking_started"
	EventSpeakingStopped GateEventKind = "speaking_stopped"
)

// GateEvent is delivered (synchronously, after the state change settles) to
// the registered sink. At is a wall-clock UTC stamp for persisted event
// records only (D42#9: never use it for interval math).
type GateEvent struct {
	Kind   GateEventKind `json:"kind"`
	Reason string        `json:"reason,omitempty"`
	At     time.Time     `json:"at"`
}

// GateOption configures a HalfDuplexGate.
type GateOption func(*HalfDuplexGate)

// WithStartMuted starts the gate muted - map `[audio] mic_muted_default`
// (config AudioSection.MicMutedDefault) here at boot wiring.
func WithStartMuted(muted bool) GateOption {
	return func(g *HalfDuplexGate) { g.muted = muted }
}

// WithGateEvents registers the event sink (mute hotkey path -> Muted state;
// the ball/state machine subscribes at wiring time, ticket 07/16).
func WithGateEvents(fn func(GateEvent)) GateOption {
	return func(g *HalfDuplexGate) { g.onEvent = fn }
}

// HalfDuplexGate wraps an AudioSource with the D16 duplex policy.
type HalfDuplexGate struct {
	inner   AudioSource
	path    Path
	onEvent func(GateEvent)
	meter   *meter

	mu       sync.Mutex
	ctx      context.Context
	buf      chan<- []byte
	running  bool
	open     bool // capture currently handed to inner
	muted    bool
	speaking bool
}

// NewHalfDuplexGate wraps inner with the duplex policy of path.
func NewHalfDuplexGate(inner AudioSource, path Path, opts ...GateOption) *HalfDuplexGate {
	g := &HalfDuplexGate{
		inner: inner,
		path:  path,
		meter: newMeter("gate-" + string(path)),
	}
	for _, o := range opts {
		o(g)
	}
	return g
}

// Start implements AudioSource. When the gate is (currently) closed the inner
// source is NOT started; SetMuted(false)/SetSpeaking(false) will start it
// later with the stored ctx/buf, so Start still returns nil.
func (g *HalfDuplexGate) Start(ctx context.Context, buf chan<- []byte) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.running {
		return fmt.Errorf("gate: already started")
	}
	g.ctx = ctx
	g.buf = buf
	g.running = true
	if g.effectiveOpen() {
		g.open = g.openInnerLocked()
	}
	return nil
}

// Stop implements AudioSource. Idempotent; closes the inner capture.
func (g *HalfDuplexGate) Stop() error {
	g.mu.Lock()
	running := g.running
	g.running = false
	g.open = false
	g.ctx, g.buf = nil, nil
	g.mu.Unlock()
	if !running {
		return nil
	}
	return g.inner.Stop()
}

// SetSpeaking is driven by the TTS side. On Path T, speaking closes capture
// (half duplex); on Path C it only reports (full duplex, D47).
func (g *HalfDuplexGate) SetSpeaking(speaking bool) {
	g.mu.Lock()
	if g.speaking == speaking {
		g.mu.Unlock()
		return
	}
	g.speaking = speaking
	reason := "speaking"
	var evs []GateEvent
	if speaking {
		evs = append(evs, GateEvent{Kind: EventSpeakingStarted, At: observe.NowWallUTC()})
	} else {
		reason = "speech-ended"
		evs = append(evs, GateEvent{Kind: EventSpeakingStopped, At: observe.NowWallUTC()})
	}
	shouldOpen := g.running && g.effectiveOpen()
	if shouldOpen && !g.open {
		g.open = g.openInnerLocked()
		if g.open {
			evs = append(evs, GateEvent{Kind: EventCaptureOpened, Reason: reason, At: observe.NowWallUTC()})
		}
	} else if !shouldOpen && g.open {
		g.open = false
		g.closeInnerLocked()
		evs = append(evs, GateEvent{Kind: EventCaptureClosed, Reason: reason, At: observe.NowWallUTC()})
	}
	fn := g.onEvent
	g.mu.Unlock()
	for _, e := range evs {
		if fn != nil {
			fn(e)
		}
	}
}

// SetMuted is driven by the mute hotkey. Muted closes capture on both paths
// (user intent outranks the duplex policy); the Muted/Unmuted events are the
// hook into the Muted ball state (07).
func (g *HalfDuplexGate) SetMuted(muted bool) {
	g.mu.Lock()
	if g.muted == muted {
		g.mu.Unlock()
		return
	}
	g.muted = muted
	var evs []GateEvent
	kind := EventMuted
	if !muted {
		kind = EventUnmuted
	}
	evs = append(evs, GateEvent{Kind: kind, At: observe.NowWallUTC()})
	shouldOpen := g.running && g.effectiveOpen()
	if shouldOpen && !g.open {
		g.open = g.openInnerLocked()
		if g.open {
			evs = append(evs, GateEvent{Kind: EventCaptureOpened, Reason: "unmuted", At: observe.NowWallUTC()})
		}
	} else if !shouldOpen && g.open {
		g.open = false
		g.closeInnerLocked()
		evs = append(evs, GateEvent{Kind: EventCaptureClosed, Reason: "muted", At: observe.NowWallUTC()})
	}
	fn := g.onEvent
	g.mu.Unlock()
	for _, e := range evs {
		if fn != nil {
			fn(e)
		}
	}
}

// Muted reports the current mute state.
func (g *HalfDuplexGate) Muted() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.muted
}

// Open reports whether capture is currently handed to the inner source.
func (g *HalfDuplexGate) Open() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.open
}

// Path reports the configured voice path.
func (g *HalfDuplexGate) Path() Path { return g.path }

// Stats implements Stater: it surfaces the inner source's counters (drops,
// reopens, device errors) plus the gate's own close counters in LastError-
// style fields. When the inner source is not a Stater the snapshot is empty.
func (g *HalfDuplexGate) Stats() Stats {
	g.mu.Lock()
	inner := g.inner
	g.mu.Unlock()
	if s, ok := inner.(Stater); ok {
		return s.Stats()
	}
	return Stats{}
}

// effectiveOpen reports whether the CURRENT state allows capture. Callers
// hold mu. muted outranks everything (user intent); on Path T speaking
// closes capture; on Path C (D47) speaking does not.
func (g *HalfDuplexGate) effectiveOpen() bool {
	if g.muted {
		return false
	}
	if g.speaking && g.path == PathT {
		return false
	}
	return true
}

// openInnerLocked starts the inner source and reports whether it is now
// delivering. A failed reopen is surfaced (meter.fail) - never a silent
// dead seam (D42#2/D42#12).
func (g *HalfDuplexGate) openInnerLocked() bool {
	if err := g.inner.Start(g.ctx, g.buf); err != nil {
		g.meter.fail(err)
		return false
	}
	return true
}

func (g *HalfDuplexGate) closeInnerLocked() {
	if err := g.inner.Stop(); err != nil {
		g.meter.fail(err)
	}
}
