//go:build windows

package aecspike

// The long-window loop-back rig (59-w-r1): the aec3-r1 rig's WASAPI plumbing
// run for a caller-named number of SECONDS of audio instead of one fixture
// pass. The far fixture loops; the loopback capture streams to disk in
// bounded segments so run length costs disk, not memory. The stop condition
// is the audio clock alone - a pushed/captured frame target derived from the
// negotiated mix rate - never a wall clock (D42#9), and the device-event
// wait stays the bounded wait the production capture loop uses.
//
// Everything else (pinned thread via audio.SpawnCapture, the session-volume
// reset, the denormal rescale, the runaway/idle guards) is the proven
// aec3-r1 shape; see rig_windows.go and docs/evidence/s1/aec3r1-remeasure.md.

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/observe"
)

// LongRunConfig names one long-window capture.
type LongRunConfig struct {
	EndpointID string // empty = console default render endpoint
	Volume     float64
	Far        []float32 // mono far fixture at FarRate; looped for the run
	FarRate    int
	Seconds    float64 // the audio-clock target (frame count derived from this)
	// SegmentSeconds bounds each capture WAV written to OutDir (5s keeps a
	// 60s run's memory at one segment, ~0.5MB at the mix rate).
	SegmentSeconds float64
	OutDir         string
	Name           string // file prefix; segments land as <name>-seg-NNN.wav
}

// LongRunSegment is one flushed capture segment.
type LongRunSegment struct {
	File       string `json:"file"`
	StartFrame int    `json:"start_frame"` // mix-rate frame index of segment start
	Frames     int    `json:"frames"`
}

// LongRunResult is one long-window capture's telemetry. MixFar is the mono
// far AT THE MIX RATE exactly as the rig looped it - the phase-continuous
// far reference the offline AEC pass reconstructs with FarClock.
type LongRunResult struct {
	EndpointID     string  `json:"endpoint_id"`
	EndpointName   string  `json:"endpoint_name"`
	MixRate        int     `json:"mix_rate"`
	MixChannels    int     `json:"mix_channels"`
	MixFloat       bool    `json:"mix_float"`
	SessionVolume  float64 `json:"session_volume"`
	TargetFrames   int     `json:"target_frames"`
	PushedFrames   int     `json:"pushed_frames"`
	CapturedFrames int     `json:"captured_frames"`
	// RatioPpm keeps the baseline rig's field but NOT its baseline
	// semantics: the long rig's render cursor keeps looping while the
	// capture drains the echo tail to the target, so captured-pushed is
	// about minus one render buffer by construction (observed -3156ppm at
	// 60s, -52924ppm at 3.4s). Clock health on a long run reads from the
	// alignment strength and the delay-track instruments, not from here.
	RatioPpm int              `json:"ratio_ppm"`
	Loops    int              `json:"far_loops"`
	MixFar   []float32        `json:"-"`
	Segments []LongRunSegment `json:"segments"`
}

// RunLoopbackLong is RunLoopback's long-window twin: same endpoint, same
// session volume discipline, same event-driven loop, but the far loops and
// the capture is flushed to disk every SegmentSeconds. It runs entirely on
// one pinned thread booked through audio.SpawnCapture and answers through
// the same fail/finish discipline as RunLoopback.
func RunLoopbackLong(reg *observe.Registry, ctx context.Context, cfg LongRunConfig) (*LongRunResult, error) {
	rb := &longRunBundle{cfg: cfg, ready: make(chan error, 1)}
	done := make(chan struct{})
	_ = audio.SpawnCapture(reg, func(cctx context.Context) {
		defer close(done)
		rb.run(cctx)
	})
	select {
	case err := <-rb.ready:
		if err != nil {
			return nil, err
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	<-done
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return rb.result, rb.runErr
}

// longRunBundle carries one long-window run. All fields live on the pinned
// rig thread; the handback fields go back under mu at the end.
type longRunBundle struct {
	cfg      LongRunConfig
	ready    chan error
	mu       sync.Mutex
	answered bool
	result   *LongRunResult
	runErr   error
}

// answer delivers the run's result exactly once (fail/finish/panic guard).
func (rb *longRunBundle) answer(err error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	if rb.answered {
		return
	}
	rb.answered = true
	rb.runErr = err
	select {
	case rb.ready <- err:
	default:
	}
}

func (rb *longRunBundle) fail(err error) {
	slog.Debug("aecspike long rig failed", "error", err.Error())
	rb.answer(err)
}

// run is the pinned rig thread for the long window. The stream setup is the
// aec3-r1 rig's; only the loop body differs (looping far cursor, segment
// flush, audio-clock target).
func (rb *longRunBundle) run(ctx context.Context) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer func() {
		rb.mu.Lock()
		answered := rb.answered
		rb.mu.Unlock()
		if !answered {
			rb.answer(fmt.Errorf("aecspike: long rig thread exited without an answer"))
		}
	}()
	if comInit := comInitMTA(); comInit {
		defer comUninitMTA()
	}

	enum, release, err := createEnumerator()
	if err != nil {
		rb.fail(err)
		return err
	}
	defer release()

	var dev unsafe.Pointer
	if rb.cfg.EndpointID == "" {
		if hr, _, _ := comCall(enum, 4, eRender, roleConsole, uintptr(unsafe.Pointer(&dev))); !hrOK(hr) || dev == nil {
			err := fmt.Errorf("aecspike: no default render endpoint (hr=%s)", hrHex(hr))
			rb.fail(err)
			return err
		}
	} else {
		id16, e := windows.UTF16PtrFromString(rb.cfg.EndpointID)
		if e != nil {
			rb.fail(e)
			return e
		}
		if hr, _, _ := comCall(enum, 5, uintptr(unsafe.Pointer(id16)), uintptr(unsafe.Pointer(&dev))); !hrOK(hr) || dev == nil {
			err := fmt.Errorf("aecspike: render endpoint not found (hr=%s)", hrHex(hr))
			rb.fail(err)
			return err
		}
	}
	defer comRelease(&dev)

	render := &stream{dev: dev}
	if err := render.initRender(); err != nil {
		rb.fail(err)
		return err
	}
	defer render.close()
	loop := &stream{dev: dev}
	if err := loop.initLoopback(); err != nil {
		rb.fail(err)
		return err
	}
	defer loop.close()
	if vol, err := SetSessionVolume(render.client, 1.0); err != nil {
		slog.Warn("aecspike long rig could not reset session volume", "error", err.Error())
	} else {
		slog.Info("aecspike long rig session volume reset", "readback", vol)
	}

	mix := render.fmt
	var bufferSize uint32
	comCall(render.client, 4, uintptr(unsafe.Pointer(&bufferSize)))
	slog.Info("aecspike long rig opened", "name", rb.cfg.Name, "rate", mix.rate,
		"channels", mix.channels, "bits", mix.bits, "render_buf_frames", bufferSize)

	mono := ResampleLinear(rb.cfg.Far, rb.cfg.FarRate, int(mix.rate))
	loopSamples := len(mono)
	if loopSamples == 0 {
		err := fmt.Errorf("aecspike: long rig far fixture is empty")
		rb.fail(err)
		return err
	}
	gain := float32(rb.cfg.Volume)
	inter := make([]float32, loopSamples*int(mix.channels))
	for i, v := range mono {
		v *= gain
		for c := 0; c < int(mix.channels); c++ {
			inter[i*int(mix.channels)+c] = v
		}
	}

	target := int(float64(rb.cfg.Seconds) * float64(mix.rate))
	segFrames := int(rb.cfg.SegmentSeconds * float64(mix.rate))
	if segFrames <= 0 || segFrames > target {
		segFrames = target // degenerate window: one segment
	}

	// Loopback capture starts first so the render's first period is caught
	// (the aec3-r1 rig's ordering).
	if hr, _, _ := comCall(loop.client, 10); !hrOK(hr) {
		err := fmt.Errorf("aecspike: loopback Start hr=%s", hrHex(hr))
		rb.fail(err)
		return err
	}
	if hr, _, _ := comCall(render.client, 10); !hrOK(hr) {
		err := fmt.Errorf("aecspike: render Start hr=%s", hrHex(hr))
		rb.fail(err)
		return err
	}

	rb.mu.Lock()
	rb.result = &LongRunResult{
		EndpointID:    rb.cfg.EndpointID,
		EndpointName:  deviceNameOf(dev),
		MixRate:       int(mix.rate),
		MixChannels:   int(mix.channels),
		MixFloat:      mix.tag == waveFormatFloat,
		SessionVolume: rb.cfg.Volume,
		TargetFrames:  target,
		MixFar:        mono,
	}
	rb.mu.Unlock()

	pushed, written, capFrames := 0, 0, 0
	var seg []int16
	idle, lastCap := 0, -1
	reason := "target-reached"
	defer func() {
		slog.Debug("aecspike long rig exited", "reason", reason, "pushed", pushed, "captured", capFrames)
	}()
	for {
		if ctx.Err() != nil {
			rb.flushSegment(&seg)
			rb.fail(ctx.Err())
			return ctx.Err()
		}
		handles := []windows.Handle{render.ev, loop.ev}
		state, _, _ := pWaitForMultiple.Call(2, uintptr(unsafe.Pointer(&handles[0])), 0, 100)
		if state == waitFailed {
			rb.flushSegment(&seg)
			err := fmt.Errorf("aecspike: WaitForMultipleObjects failed")
			rb.fail(err)
			return err
		}
		if state == waitObject0 {
			n, err := fillLoop(render, inter, &written, bufferSize, mix)
			if err != nil {
				rb.flushSegment(&seg)
				rb.fail(err)
				return err
			}
			pushed += n
		}
		if state == waitObject0+1 {
			n, err := loop.drain(&seg, mix)
			if err != nil {
				rb.flushSegment(&seg)
				rb.fail(err)
				return err
			}
			capFrames += n
			for len(seg) >= segFrames {
				rb.writeSegment(&seg, segFrames)
			}
		}
		if pushed >= target && capFrames >= target {
			break
		}
		if pushed >= target && capFrames == lastCap {
			idle++
			if idle > 50 { // 5s of waits with no capture progress: stop, keep data
				reason = "idle-stall"
				break
			}
		} else {
			idle = 0
		}
		lastCap = capFrames
		if capFrames > target+loopSamples+int(mix.rate) {
			reason = "runaway-capture"
			break // runaway mix from other sessions; keep what we have
		}
	}
	rb.flushSegment(&seg)

	rb.mu.Lock()
	if res := rb.result; res != nil {
		res.PushedFrames = pushed
		res.CapturedFrames = capFrames
		res.Loops = written / (loopSamples * int(mix.channels))
		if pushed > 0 {
			res.RatioPpm = int(int64(capFrames-pushed) * 1_000_000 / int64(pushed))
		}
	}
	rb.mu.Unlock()
	rb.answer(nil)
	return nil
}

// writeSegment flushes a full segment to disk and keeps the remainder in
// place (the capture stream continues at segFrames).
func (rb *longRunBundle) writeSegment(seg *[]int16, segFrames int) {
	if len(*seg) < segFrames {
		return
	}
	rb.appendSegment((*seg)[:segFrames])
	rest := copy(*seg, (*seg)[segFrames:])
	*seg = (*seg)[:rest]
}

// flushSegment writes whatever is left as the run's final (partial) segment.
func (rb *longRunBundle) flushSegment(seg *[]int16) {
	if len(*seg) == 0 {
		return
	}
	rb.appendSegment(*seg)
	*seg = (*seg)[:0]
}

func (rb *longRunBundle) appendSegment(samples []int16) {
	rb.mu.Lock()
	res := rb.result
	rb.mu.Unlock()
	if res == nil {
		return
	}
	start := 0
	if len(res.Segments) > 0 {
		last := res.Segments[len(res.Segments)-1]
		start = last.StartFrame + last.Frames
	}
	name := fmt.Sprintf("%s-seg-%03d.wav", rb.cfg.Name, len(res.Segments))
	path := filepath.Join(rb.cfg.OutDir, name)
	if err := WriteWav16(path, res.MixRate, samples); err != nil {
		slog.Warn("aecspike long rig segment write failed", "file", name, "error", err.Error())
		return
	}
	res.Segments = append(res.Segments, LongRunSegment{
		File: path, StartFrame: start, Frames: len(samples),
	})
}

// fillLoop is stream.fill with a looping source: the render cursor wraps the
// one-loop interleave buffer, so a run of any length plays the fixture
// seamlessly. Same GetBuffer/ReleaseBuffer discipline and the same
// BUFFER_TOO_LARGE grace as the baseline rig.
func fillLoop(s *stream, inter []float32, written *int, bufferSize uint32, mix waveFormat) (int, error) {
	var padding uint32
	if hr, _, _ := comCall(s.client, 6, uintptr(unsafe.Pointer(&padding))); !hrOK(hr) {
		return 0, fmt.Errorf("aecspike: GetCurrentPadding hr=%s", hrHex(hr))
	}
	ch := int(mix.channels)
	avail := int(bufferSize) - int(padding)
	if avail <= 0 {
		return 0, nil
	}
	var data unsafe.Pointer
	gbHr, _, _ := comCall(s.service, 3, uintptr(avail), uintptr(unsafe.Pointer(&data)))
	if !hrOK(gbHr) {
		if uint32(gbHr) == 0x88890007 { // right-after-Start padding lag: wait a period
			return 0, nil
		}
		return 0, fmt.Errorf("aecspike: GetBuffer(render) hr=%s", hrHex(gbHr))
	}
	n := avail
	if n <= 0 {
		comCall(s.service, 4, 0)
		return 0, nil
	}
	loopSamples := len(inter)
	if mix.tag == waveFormatFloat {
		out := unsafe.Slice((*float32)(data), n*ch)
		for k := 0; k < n*ch; k++ {
			out[k] = inter[(*written+k)%loopSamples]
		}
	} else {
		out := unsafe.Slice((*int16)(data), n*ch)
		for k := 0; k < n*ch; k++ {
			v := inter[(*written+k)%loopSamples] * 32767
			if v > 32767 {
				v = 32767
			} else if v < -32768 {
				v = -32768
			}
			out[k] = int16(v)
		}
	}
	if rbHr, _, _ := comCall(s.service, 4, uintptr(n), 0); !hrOK(rbHr) {
		return 0, fmt.Errorf("aecspike: ReleaseBuffer(render) hr=%s", hrHex(rbHr))
	}
	*written += n * ch
	return n, nil
}

// deviceNameOf re-reads the endpoint's friendly name for the report.
func deviceNameOf(dev unsafe.Pointer) string {
	_, name := deviceIdentity(dev)
	return name
}
