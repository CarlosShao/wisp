//go:build windows

package audio

// WASAPIMicrophone: the real C8 source (SPEC-04 sec 3). Shared-mode WASAPI
// capture on a PINNED OS thread (D38a: the audio-capture thread runs no
// other Go code - runtime.LockOSThread for the whole loop), device-native
// rate resampled to the 16k/mono/int16 seam in-process, frames pushed into
// the bounded channel with D38d drop counting, hotplug reopen-once
// supervision (D42#2) and occupied/permission error mapping (D42#12).
//
// Testability: the enumeration seam (DeviceWatcher) and the open seam
// (streamOpener) are injected. Tests drive the REAL capture loop, resampler
// and bounded push with fake watcher/opener whose streams carry wav data -
// the audio-data seam is never mocked (SPEC-04 sec 9).

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// reopenWait is the WaitEvent timeout: how often the pinned loop re-checks
// stop/hotplug signals while idle. Bounded loop exit keeps D38e step 4
// (stop capture before ASR release) fast.
const reopenWaitMs = 200

// WASAPIMicrophone implements AudioSource (C8) over WASAPI shared mode.
type WASAPIMicrophone struct {
	watcher DeviceWatcher
	opener  streamOpener
	meter   *meter

	threadID atomic.Uint32 // OS thread of the pinned capture loop (diagnostics/pinned test)

	mu      sync.Mutex
	running bool
	cancel  context.CancelFunc
	done    chan struct{}
	errC    chan error
	hotplug chan struct{}
}

// NewWASAPIMicrophone builds the real device stack (MMDeviceEnumerator
// watcher + WASAPI shared-mode opener).
func NewWASAPIMicrophone() *WASAPIMicrophone {
	return newWASAPIMicrophoneWith(newMMDeviceWatcher(), wasapiOpener{})
}

// newWASAPIMicrophoneWith injects the seams (tests).
func newWASAPIMicrophoneWith(w DeviceWatcher, o streamOpener) *WASAPIMicrophone {
	return &WASAPIMicrophone{
		watcher: w,
		opener:  o,
		meter:   newMeter("wasapi-mic"),
		hotplug: make(chan struct{}, 1),
	}
}

// Start implements AudioSource. It returns after the pinned capture thread
// has opened the device (or with the D42-classified open error, so occupied
// /permission failures surface at Start - D42#12, never silent).
func (m *WASAPIMicrophone) Start(ctx context.Context, buf chan<- []byte) error {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return fmt.Errorf("wasapi-mic: already started")
	}
	loopCtx, cancel := context.WithCancel(ctx)
	m.cancel = cancel
	m.done = make(chan struct{})
	m.errC = make(chan error, 1)
	m.running = true
	m.mu.Unlock()

	started := make(chan error, 1)
	observe.Default.Spawn("audio-capture", "audio", nil, func(c context.Context) {
		// c is the registry's ctx (nil root -> background); the loop must
		// honor the CONSUMER's ctx, so run the derived one explicitly.
		m.run(loopCtx, buf, started)
	})
	return <-started
}

// Stop implements AudioSource: cancels the loop (exits within reopenWaitMs)
// and waits for the pinned thread to release the device (D38e step 4).
// Idempotent; safe without Start.
func (m *WASAPIMicrophone) Stop() error {
	m.mu.Lock()
	cancel := m.cancel
	done := m.done
	running := m.running
	m.cancel = nil
	m.running = false
	m.mu.Unlock()
	if !running {
		return nil
	}
	cancel()
	tm := observe.NewTimeout(2 * time.Second)
	select {
	case <-done:
		return nil
	case <-time.After(tm.Remaining()):
		// Abandon the wait like D38e step 3 does; keep it visible.
		slog.Error("audio capture thread did not exit within 2s (abandoned wait)", "source", "wasapi-mic")
		return nil
	}
}

// Err returns the terminal error of the last capture session, if any
// (non-blocking). Hotplug reopen failures (D42#2) land here.
func (m *WASAPIMicrophone) Err() error {
	select {
	case err := <-m.errC:
		return err
	default:
		return nil
	}
}

// ThreadID returns the OS thread id of the pinned capture loop (0 when not
// running). Diagnostics + the pinned-thread stability test use it.
func (m *WASAPIMicrophone) ThreadID() uint32 { return m.threadID.Load() }

// Stats implements Stater.
func (m *WASAPIMicrophone) Stats() Stats { return m.meter.snapshot() }

// Endpoints returns the current capture/render endpoint pair (D47: kept
// queryable for the P15 AEC clock-drift spike).
func (m *WASAPIMicrophone) Endpoints() (capture, render DeviceDescriptor, err error) {
	capture, err = m.watcher.DefaultCaptureDevice()
	if err != nil {
		return capture, render, err
	}
	render, err = m.watcher.DefaultRenderDevice()
	return capture, render, err
}

// run is the pinned capture loop. It pins the goroutine to ONE OS thread for
// its whole life (D38a) and runs nothing but capture work: open, wait, drain,
// resample, bounded push, hotplug reopen.
func (m *WASAPIMicrophone) run(ctx context.Context, buf chan<- []byte, started chan<- error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	m.threadID.Store(currentThreadID())
	defer func() { m.threadID.Store(0); close(m.done) }()

	comInit := comInitMTA()
	if comInit {
		defer comUninitMTA()
	}

	unsub, err := m.watcher.OnDefaultCaptureChanged(func() {
		// Runs on a Windows notification thread: non-blocking, no COM.
		select {
		case m.hotplug <- struct{}{}:
		default:
		}
	})
	if err != nil {
		err := observe.Wrap(observe.ClassAudioDevice, err, "hotplug notification registration failed")
		m.meter.fail(err)
		started <- err
		return
	}
	defer unsub()

	stream, err := m.openCurrent()
	if err != nil {
		m.meter.fail(err)
		started <- err
		m.postErr(err)
		return
	}
	started <- nil
	defer func() { _ = stream.Close() }() // closure: closes whatever stream is CURRENT (swap-safe)

	res := NewResampler(stream.Rate(), TargetRate)
	var pending []int16

	for {
		if ctx.Err() != nil {
			return
		}
		// Hotplug (D42#2): default device changed -> re-enumerate once.
		select {
		case <-m.hotplug:
			next, err := m.reopenOnce(stream.Descriptor())
			if err != nil {
				// Never silently keep the stale handle: fail loudly, stop.
				m.meter.fail(err)
				m.postErr(err)
				return
			}
			_ = stream.Close()
			stream = next
			res = NewResampler(stream.Rate(), TargetRate)
			pending = pending[:0]
			continue
		default:
		}

		ok, err := stream.WaitEvent(reopenWaitMs)
		if err != nil || !ok {
			// Timeout just means no new period; an error is a stale handle.
			if err == nil {
				continue
			}
			if !m.handleStreamError(ctx, &stream, &res, &pending) {
				return
			}
			continue
		}

		samples, err := stream.Drain()
		if err != nil {
			if !m.handleStreamError(ctx, &stream, &res, &pending) {
				return
			}
			continue
		}
		pending = append(pending, res.Process(samples)...)
		for len(pending) >= FrameSamples {
			m.meter.push(buf, encodeFrame(pending[:FrameSamples]))
			pending = append(pending[:0], pending[FrameSamples:]...)
		}
	}
}

// handleStreamError applies the D42#2 recovery: re-enumerate ONCE and swap
// the stream; a failed reopen is terminal (no silent stale handle). Returns
// false when the loop must stop.
func (m *WASAPIMicrophone) handleStreamError(ctx context.Context, stream *deviceStream, res **Resampler, pending *[]int16) bool {
	if ctx.Err() != nil {
		return false
	}
	next, err := m.reopenOnce((*stream).Descriptor())
	if err != nil {
		m.meter.fail(err)
		m.postErr(err)
		return false
	}
	_ = (*stream).Close()
	*stream = next
	*res = NewResampler((*stream).Rate(), TargetRate)
	*pending = (*pending)[:0]
	return true
}

// openCurrent opens the current default capture endpoint (initial open).
func (m *WASAPIMicrophone) openCurrent() (deviceStream, error) {
	dev, err := m.watcher.DefaultCaptureDevice()
	if err != nil {
		return nil, observe.Wrap(observe.ClassAudioDevice, err, "capture device enumeration failed")
	}
	return m.opener.Open(dev)
}

// reopenOnce is the D42#2 recovery step: re-enumerate once, open once
// (observe.ClassAudioDevice.RetryPolicy() == RetryOnce, D37). Success bumps
// the Reopens telemetry counter; failure names the device we refuse to keep
// silently (prev = the stale handle's descriptor).
func (m *WASAPIMicrophone) reopenOnce(prev DeviceDescriptor) (deviceStream, error) {
	s, err := m.openCurrent()
	if err != nil {
		return nil, observe.Wrap(observe.ClassAudioDevice, err,
			"hotplug reopen failed; keeping no stale handle (previous device "+prev.String()+")")
	}
	m.meter.reopened()
	slog.Info("audio capture reopened after device change",
		"source", "wasapi-mic", "device", s.Descriptor().String())
	return s, nil
}

// postErr delivers a terminal error to Err() without ever blocking.
func (m *WASAPIMicrophone) postErr(err error) {
	select {
	case m.errC <- err:
	default:
	}
}
