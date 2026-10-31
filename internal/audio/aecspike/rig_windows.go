//go:build windows

package aecspike

// The loop-back rig: one pinned OS thread opens a shared-mode WASAPI render
// stream and an endpoint-loopback capture stream on the SAME render endpoint
// (AUDCLNT_STREAMFLAGS_LOOPBACK), plays the far-end fixture and records what
// the endpoint mix returned. The technique is the package-internal one
// proven by internal/audio's own WASAPI code (x/sys COM vtable calls, no
// cgo); it is duplicated here because those helpers are unexported and the
// spike must not modify the production capture files.
//
// Threading: RunLoopback books its whole life through audio.SpawnCapture (the
// sanctioned "audio-capture" roster name), pins the thread (D38a) and owns
// every COM object it creates. The bounded wait is the device-event wait the
// production capture loop also uses - never a wall-clock timeout (D42#9).
//
// This is a spike instrument: it holds no audio after the caller consumes
// the capture, and nothing in the resident process references it.

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/observe"
)

var (
	modOle32    = windows.NewLazySystemDLL("ole32.dll")
	modKernel32 = windows.NewLazySystemDLL("kernel32.dll")

	pCoInitializeEx   = modOle32.NewProc("CoInitializeEx")
	pCoUninitialize   = modOle32.NewProc("CoUninitialize")
	pCoCreateInstance = modOle32.NewProc("CoCreateInstance")
	pCoTaskMemFree    = modOle32.NewProc("CoTaskMemFree")
	pWaitForMultiple  = modKernel32.NewProc("WaitForMultipleObjects")
)

var (
	iidIUnknown             = windows.GUID{Data4: [8]byte{0xC0, 0, 0, 0, 0, 0, 0, 0x46}}
	clsidMMDeviceEnumerator = windows.GUID{
		Data1: 0xBCDE0395, Data2: 0xE52F, Data3: 0x467C,
		Data4: [8]byte{0x8E, 0x3D, 0xC4, 0x57, 0x92, 0x91, 0x69, 0x2E},
	}
	iidIMMDeviceEnumerator = windows.GUID{
		Data1: 0xA95664D2, Data2: 0x9614, Data3: 0x4F35,
		Data4: [8]byte{0xA7, 0x46, 0xDE, 0x8D, 0xB6, 0x36, 0x17, 0xE6},
	}
	iidIAudioClient = windows.GUID{
		Data1: 0x1CB9AD4C, Data2: 0xDBFA, Data3: 0x4C32,
		Data4: [8]byte{0xB1, 0x78, 0xC2, 0xF5, 0x68, 0xA7, 0x03, 0xB2},
	}
	iidIAudioCaptureClient = windows.GUID{
		Data1: 0xC8ADBD64, Data2: 0xE71E, Data3: 0x48A0,
		Data4: [8]byte{0xA4, 0xDE, 0x18, 0x5C, 0x39, 0x5C, 0xD3, 0x17},
	}
	iidIAudioRenderClient = windows.GUID{
		Data1: 0xF294ACFC, Data2: 0x3146, Data3: 0x4483,
		Data4: [8]byte{0xA7, 0xBF, 0xAD, 0xDC, 0xA7, 0xC2, 0x60, 0xE2},
	}
	pkeyFriendlyName = propKey{
		fmtid: windows.GUID{
			Data1: 0xA45C254E, Data2: 0xDF1C, Data3: 0x4EFD,
			Data4: [8]byte{0x80, 0x20, 0x67, 0xD1, 0x46, 0xA8, 0x50, 0xE0},
		},
		pid: 14,
	}
)

const (
	clsctxAll       = 0x17
	coinitMTA       = 0x0
	eRender         = 0
	eCapture        = 1
	roleConsole     = 0
	stateActive     = 1
	vtLPWSTR        = 31
	waveFormatPCM   = 1
	waveFormatExt   = 0xFFFE
	waveFormatFloat = 3
	shareShared     = 0
	flagEventCB     = 0x00040000
	flagLoopback    = 0x00020000
	flagSilent      = 0x2
	bufHns          = 2_000_000 // 200ms
	waitObject0     = 0
	waitTimeout     = 0x00000102
	waitFailed      = 0xFFFFFFFF
	slotRelease     = 2
)

type propKey struct {
	fmtid windows.GUID
	pid   uint32
}

type propVariant struct {
	vt uint16
	_  [6]byte
	// ptr is the union member this spike reads: VT_LPWSTR's pointer.
	ptr unsafe.Pointer
}

// waveFormat is the negotiated mix-format subset the rig needs.
type waveFormat struct {
	tag      uint16
	channels uint16
	rate     uint32
	bits     uint16
}

func comCall(this unsafe.Pointer, slot uintptr, args ...uintptr) (uintptr, uintptr, error) {
	vt := *(*unsafe.Pointer)(this)
	fn := *(*unsafe.Pointer)(unsafe.Add(vt, slot*8))
	return syscall.SyscallN(uintptr(fn), append([]uintptr{uintptr(this)}, args...)...)
}

func comRelease(obj *unsafe.Pointer) {
	if *obj != nil {
		_, _, _ = comCall(*obj, slotRelease)
		*obj = nil
	}
}

func hrOK(hr uintptr) bool { return int32(hr) >= 0 }

func hrHex(hr uintptr) string { return fmt.Sprintf("0x%08X", uint32(hr)) }

func comInitMTA() bool {
	hr, _, _ := pCoInitializeEx.Call(0, coinitMTA)
	return hrOK(hr)
}

func comUninitMTA() { pCoUninitialize.Call() }

func createEnumerator() (unsafe.Pointer, func(), error) {
	var enum unsafe.Pointer
	hr, _, _ := pCoCreateInstance.Call(
		uintptr(unsafe.Pointer(&clsidMMDeviceEnumerator)), 0, clsctxAll,
		uintptr(unsafe.Pointer(&iidIMMDeviceEnumerator)),
		uintptr(unsafe.Pointer(&enum)))
	if !hrOK(hr) || enum == nil {
		return nil, nil, fmt.Errorf("aecspike: CoCreateInstance(MMDeviceEnumerator) hr=%s", hrHex(hr))
	}
	return enum, func() { comRelease(&enum) }, nil
}

// EndpointInfo names one active audio endpoint.
type EndpointInfo struct {
	Flow string `json:"flow"` // "render" or "capture"
	ID   string `json:"id"`
	Name string `json:"name"`
}

// EndpointInventory is the device-side premise of every rig reading.
type EndpointInventory struct {
	Endpoints      []EndpointInfo `json:"endpoints"`
	DefaultRender  string         `json:"default_render"`
	DefaultCapture string         `json:"default_capture"`
}

// EnumerateEndpoints lists the active endpoints and the console-role
// defaults. An empty list is a valid answer (a headless runner), not an
// error; the rig tests skip on it. COM is initialized for the call and
// balanced before it returns.
func EnumerateEndpoints() (*EndpointInventory, error) {
	if comInit := comInitMTA(); comInit {
		defer comUninitMTA()
	}
	enum, release, err := createEnumerator()
	if err != nil {
		return nil, err
	}
	defer release()
	inv := &EndpointInventory{}
	for _, flow := range []struct {
		id   uint32
		name string
	}{{eRender, "render"}, {eCapture, "capture"}} {
		var coll unsafe.Pointer
		if hr, _, _ := comCall(enum, 3, uintptr(flow.id), stateActive, uintptr(unsafe.Pointer(&coll))); !hrOK(hr) {
			continue
		}
		var n uint32
		comCall(coll, 3, uintptr(unsafe.Pointer(&n)))
		for i := uint32(0); i < n; i++ {
			var dev unsafe.Pointer
			if hr, _, _ := comCall(coll, 4, uintptr(i), uintptr(unsafe.Pointer(&dev))); !hrOK(hr) {
				continue
			}
			info := EndpointInfo{Flow: flow.name}
			info.ID, info.Name = deviceIdentity(dev)
			comRelease(&dev)
			inv.Endpoints = append(inv.Endpoints, info)
		}
		comRelease(&coll)
	}
	if d, ok := defaultEndpoint(enum, eRender); ok {
		inv.DefaultRender = d
	}
	if d, ok := defaultEndpoint(enum, eCapture); ok {
		inv.DefaultCapture = d
	}
	return inv, nil
}

func defaultEndpoint(enum unsafe.Pointer, flow uint32) (string, bool) {
	var dev unsafe.Pointer
	if hr, _, _ := comCall(enum, 4, uintptr(flow), uintptr(roleConsole), uintptr(unsafe.Pointer(&dev))); !hrOK(hr) || dev == nil {
		return "", false
	}
	defer comRelease(&dev)
	id, _ := deviceIdentity(dev)
	return id, id != ""
}

// deviceIdentity reads IMMDevice::GetId (slot 5) and the friendly name
// property. Best effort: an unnamed endpoint still reports its id.
func deviceIdentity(dev unsafe.Pointer) (id, name string) {
	var idPtr *uint16
	if hr, _, _ := comCall(dev, 5, uintptr(unsafe.Pointer(&idPtr))); hrOK(hr) && idPtr != nil {
		id = windows.UTF16PtrToString(idPtr)
		pCoTaskMemFree.Call(uintptr(unsafe.Pointer(idPtr)))
	}
	var store unsafe.Pointer
	if hr, _, _ := comCall(dev, 4, 0, uintptr(unsafe.Pointer(&store))); hrOK(hr) && store != nil {
		defer comRelease(&store)
		var pv propVariant
		if hr, _, _ := comCall(store, 5, uintptr(unsafe.Pointer(&pkeyFriendlyName)),
			uintptr(unsafe.Pointer(&pv))); hrOK(hr) && pv.vt == vtLPWSTR && pv.ptr != nil {
			name = windows.UTF16PtrToString((*uint16)(pv.ptr))
		}
	}
	return id, name
}

// parseMixFormat reads the WAVEFORMATEX the engine negotiated (see
// internal/audio wasapi_windows.go for the packing rationale; WAVEFORMATEX
// is 18 bytes, the extensible body adds 22 more with SubFormat at offset 24).
func parseMixFormat(p unsafe.Pointer) (waveFormat, error) {
	if p == nil {
		return waveFormat{}, fmt.Errorf("aecspike: GetMixFormat returned no format")
	}
	f := waveFormat{
		tag:      *(*uint16)(p),
		channels: *(*uint16)(unsafe.Add(p, 2)),
		rate:     *(*uint32)(unsafe.Add(p, 4)),
		bits:     *(*uint16)(unsafe.Add(p, 14)),
	}
	if f.tag == waveFormatExt {
		sub := *(*windows.GUID)(unsafe.Add(p, 24))
		f.tag = uint16(sub.Data1)
	}
	if f.rate == 0 || f.channels == 0 {
		return f, fmt.Errorf("aecspike: degenerate mix format")
	}
	return f, nil
}

// RigReading is one loop-back run's telemetry.
type RigReading struct {
	EndpointID   string `json:"endpoint_id"`
	EndpointName string `json:"endpoint_name"`
	MixRate      int    `json:"mix_rate"`
	MixChannels  int    `json:"mix_channels"`
	MixFloat     bool   `json:"mix_float"`
	// SessionVolume is the SOFTWARE amplitude the fixture was scaled by
	// (see the fill comment for why ISimpleAudioVolume is not used).
	SessionVolume  float64 `json:"session_volume"`
	PushedFrames   int     `json:"pushed_frames"`
	CapturedFrames int     `json:"captured_frames"`
	// RatioPpm is (captured-pushed)/pushed in ppm: on one endpoint both
	// clocks are the same hardware clock, so a healthy run is ~0 and any
	// residue is the engine's own resampling, not device drift.
	RatioPpm int `json:"ratio_ppm"`
}

// RunLoopback plays far (mono float32 at farRate) on the named render
// endpoint (empty id = console default) at the given session volume and
// captures the endpoint loopback. It returns the telemetry and the captured
// mix-rate mono audio. Runs entirely on one pinned thread booked through
// audio.SpawnCapture; ctx cancellation stops it between device periods.
func RunLoopback(reg *observe.Registry, ctx context.Context, endpointID string, volume float64,
	far []float32, farRate int,
) (*RigReading, []int16, error) {
	rb := &rigBundle{
		endpointID: endpointID, volume: volume, far: far, farRate: farRate,
		ready: make(chan error, 1),
	}
	done := make(chan struct{})
	_ = audio.SpawnCapture(reg, func(cctx context.Context) {
		defer close(done)
		rb.run(cctx)
	})
	select {
	case err := <-rb.ready:
		if err != nil {
			return nil, nil, err
		}
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	}
	<-done
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return rb.reading, rb.captured, rb.runErr
}

// rigBundle carries one loop-back run. All fields are touched on the pinned
// rig thread only, except the three handed back under mu at the end.
type rigBundle struct {
	endpointID string
	volume     float64
	far        []float32
	farRate    int

	ready    chan error
	mu       sync.Mutex
	answered bool
	reading  *RigReading
	captured []int16
	runErr   error
}

// answer delivers the run's result exactly once (fail/finish/panic guard).
func (rb *rigBundle) answer(err error) {
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

// run is the pinned rig thread: open, play, capture, close. The IAudioClient
// vtable slots (audioclient.h MIDL order; the same mapping the production
// capture code pins): 3 Initialize, 4 GetBufferSize, 6 GetCurrentPadding,
// 8 GetMixFormat, 9 GetDevicePeriod, 10 Start, 11 Stop, 12 Reset,
// 13 SetEventHandle, 14 GetService.
func (rb *rigBundle) run(ctx context.Context) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	// A panic here would otherwise leave RunLoopback waiting on ready
	// forever; the registry's recover boundary closes done, so re-answer.
	defer func() {
		rb.mu.Lock()
		answered := rb.answered
		rb.mu.Unlock()
		if !answered {
			rb.answer(fmt.Errorf("aecspike: rig thread exited without an answer"))
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
	if rb.endpointID == "" {
		if hr, _, _ := comCall(enum, 4, eRender, roleConsole, uintptr(unsafe.Pointer(&dev))); !hrOK(hr) || dev == nil {
			err := fmt.Errorf("aecspike: no default render endpoint (hr=%s)", hrHex(hr))
			rb.fail(err)
			return err
		}
	} else {
		id16, e := windows.UTF16PtrFromString(rb.endpointID)
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
	// Pin the render session volume at 1.0 through the real C float call:
	// the volume axis is applied by scaling the samples, and the explicit
	// reset undoes the integer-path ISimpleAudioVolume damage recorded in
	// the evidence file (a garbage float once persisted a near-zero session
	// volume here, silencing every render run).
	if vol, err := SetSessionVolume(render.client, 1.0); err != nil {
		slog.Warn("aecspike rig could not reset session volume", "error", err.Error())
	} else {
		slog.Info("aecspike rig session volume reset", "readback", vol)
	}
	slog.Debug("aecspike rig loopback format", "tag", loop.fmt.tag, "rate", loop.fmt.rate,
		"channels", loop.fmt.channels, "bits", loop.fmt.bits)

	mix := render.fmt
	var bufferSize uint32
	comCall(render.client, 4, uintptr(unsafe.Pointer(&bufferSize))) // GetBufferSize
	slog.Debug("aecspike rig opened", "mix_tag", mix.tag, "rate", mix.rate,
		"channels", mix.channels, "bits", mix.bits, "render_buf_frames", bufferSize)
	// Resample the far fixture to the mix rate and interleave it. The volume
	// axis is applied by SCALING THE SAMPLES (with the session volume pinned
	// at 1.0 above): ISimpleAudioVolume cannot be driven from Go directly
	// because its float parameter rides the XMM registers (see the rig-open
	// reset and the evidence file).
	mono := ResampleLinear(rb.far, rb.farRate, int(mix.rate))
	frames := len(mono)
	slog.Debug("aecspike rig fixture", "far_samples", len(rb.far), "far_rate", rb.farRate,
		"mono_mix_frames", frames)
	inter := make([]float32, frames*int(mix.channels))
	gain := float32(rb.volume)
	for i, v := range mono {
		v *= gain
		for c := 0; c < int(mix.channels); c++ {
			inter[i*int(mix.channels)+c] = v
		}
	}

	// Loopback capture starts first so the render's first period is caught.
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
	pushed, written := 0, 0
	var capFrames int
	captured := make([]int16, 0, frames+int(mix.rate))
	idle, lastCap := 0, -1
	reason := "done"
	defer func() { slog.Debug("aecspike rig exited", "reason", reason, "pushed", pushed, "captured", capFrames) }()
	for {
		if ctx.Err() != nil {
			rb.fail(ctx.Err())
			return ctx.Err()
		}
		handles := []windows.Handle{render.ev, loop.ev}
		state, _, _ := pWaitForMultiple.Call(2, uintptr(unsafe.Pointer(&handles[0])), 0, 100)
		if state == waitFailed {
			err := fmt.Errorf("aecspike: WaitForMultipleObjects failed")
			rb.fail(err)
			return err
		}
		if state == waitObject0 {
			n, err := render.fill(inter, &written, bufferSize, mix)
			if err != nil {
				rb.fail(err)
				return err
			}
			pushed += n
		}
		if state == waitObject0+1 {
			n, err := loop.drain(&captured, mix)
			if err != nil {
				rb.fail(err)
				return err
			}
			capFrames += n
		}
		if pushed >= frames && capFrames >= frames {
			break
		}
		if pushed >= frames && capFrames == lastCap {
			idle++
			if idle > 50 { // 5s of waits with no capture progress: stop, keep data
				reason = "idle-stall"
				break
			}
		} else {
			idle = 0
		}
		lastCap = capFrames
		if capFrames > 2*frames+int(mix.rate) {
			reason = "runaway-capture"
			break // runaway mix from other sessions; keep what we have
		}
	}
	vol := rb.volume
	rb.finish(mix, pushed, capFrames, captured, vol)
	return nil
}

func (rb *rigBundle) fail(err error) {
	slog.Debug("aecspike rig failed", "error", err.Error())
	rb.answer(err)
}

func (rb *rigBundle) finish(mix waveFormat, pushed, capFrames int, captured []int16, volume float64) {
	rb.mu.Lock()
	if !rb.answered {
		rb.reading = &RigReading{
			EndpointID:     rb.endpointID,
			MixRate:        int(mix.rate),
			MixChannels:    int(mix.channels),
			MixFloat:       mix.tag == waveFormatFloat,
			SessionVolume:  volume,
			PushedFrames:   pushed,
			CapturedFrames: capFrames,
		}
		if pushed > 0 {
			rb.reading.RatioPpm = int(int64(capFrames-pushed) * 1_000_000 / int64(pushed))
		}
		rb.captured = captured
	}
	rb.mu.Unlock()
	rb.answer(nil)
}

// stream is one open IAudioClient on the rig's render endpoint, either the
// render side or the loopback-capture side.
type stream struct {
	dev      unsafe.Pointer
	client   unsafe.Pointer
	service  unsafe.Pointer // IAudioRenderClient or IAudioCaptureClient
	ev       windows.Handle
	fmt      waveFormat
	fillLogs int

	denormChecked bool
	denormScale   bool
}

// initRender opens the shared-mode render stream.
func (s *stream) initRender() error {
	if err := s.activate(); err != nil {
		return err
	}
	if err := s.init(shareShared, flagEventCB); err != nil {
		return err
	}
	if hr, _, _ := comCall(s.client, 14, uintptr(unsafe.Pointer(&iidIAudioRenderClient)),
		uintptr(unsafe.Pointer(&s.service))); !hrOK(hr) || s.service == nil {
		return fmt.Errorf("aecspike: GetService(IAudioRenderClient) hr=%s", hrHex(hr))
	}
	ev, err := windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return fmt.Errorf("aecspike: CreateEvent(render): %w", err)
	}
	s.ev = ev
	if hr, _, _ := comCall(s.client, 13, uintptr(ev)); !hrOK(hr) {
		return fmt.Errorf("aecspike: SetEventHandle(render) hr=%s", hrHex(hr))
	}
	return nil
}

// initLoopback opens the loopback capture stream on the same endpoint.
func (s *stream) initLoopback() error {
	if err := s.activate(); err != nil {
		return err
	}
	if err := s.init(shareShared, flagEventCB|flagLoopback); err != nil {
		return err
	}
	if hr, _, _ := comCall(s.client, 14, uintptr(unsafe.Pointer(&iidIAudioCaptureClient)),
		uintptr(unsafe.Pointer(&s.service))); !hrOK(hr) || s.service == nil {
		return fmt.Errorf("aecspike: GetService(IAudioCaptureClient) hr=%s", hrHex(hr))
	}
	ev, err := windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return fmt.Errorf("aecspike: CreateEvent(loopback): %w", err)
	}
	s.ev = ev
	if hr, _, _ := comCall(s.client, 13, uintptr(ev)); !hrOK(hr) {
		return fmt.Errorf("aecspike: SetEventHandle(loopback) hr=%s", hrHex(hr))
	}
	return nil
}

func (s *stream) activate() error {
	var client unsafe.Pointer
	if hr, _, _ := comCall(s.dev, 3, uintptr(unsafe.Pointer(&iidIAudioClient)), clsctxAll, 0,
		uintptr(unsafe.Pointer(&client))); !hrOK(hr) || client == nil {
		return fmt.Errorf("aecspike: Activate(IAudioClient) hr=%s", hrHex(hr))
	}
	s.client = client
	return nil
}

// init negotiates the mix format and initializes the stream (with the
// period-aligned retry the production capture code documents).
func (s *stream) init(share, flags uint32) error {
	var wfex unsafe.Pointer
	if hr, _, _ := comCall(s.client, 8, uintptr(unsafe.Pointer(&wfex))); !hrOK(hr) || wfex == nil {
		return fmt.Errorf("aecspike: GetMixFormat hr=%s", hrHex(hr))
	}
	f, ferr := parseMixFormat(wfex)
	if ferr != nil {
		pCoTaskMemFree.Call(uintptr(wfex))
		return ferr
	}
	s.fmt = f
	hr, _, _ := comCall(s.client, 3, uintptr(share), uintptr(flags), bufHns, 0, uintptr(wfex), 0)
	if uint32(hr) == 0x88890014 { // AUDCLNT_E_BUFFER_SIZE_NOT_ALIGNED: retry aligned
		var defPeriod int64
		comCall(s.client, 9, uintptr(unsafe.Pointer(&defPeriod)), uintptr(unsafe.Pointer(&defPeriod)))
		if defPeriod > 0 {
			rate := uint64(f.rate)
			aligned := (uint64(bufHns)*rate + 9_999_999) / 10_000_000
			periodFrames := uint64(defPeriod) * rate / 10_000_000
			if periodFrames > 0 {
				aligned = (aligned + periodFrames - 1) / periodFrames * periodFrames
			}
			comCall(s.client, 12) // Reset before re-Initialize
			hr, _, _ = comCall(s.client, 3, uintptr(share), uintptr(flags),
				uintptr(aligned*10_000_000/rate), 0, uintptr(wfex), 0)
		}
	}
	pCoTaskMemFree.Call(uintptr(wfex))
	if !hrOK(hr) {
		return fmt.Errorf("aecspike: Initialize hr=%s", hrHex(hr))
	}
	return nil
}

// fill writes as much of the interleaved far buffer as the render buffer
// accepts (IAudioRenderClient: 3 GetBuffer, 4 ReleaseBuffer).
func (s *stream) fill(inter []float32, written *int, bufferSize uint32, mix waveFormat) (int, error) {
	var padding uint32
	if hr, _, _ := comCall(s.client, 6, uintptr(unsafe.Pointer(&padding))); !hrOK(hr) {
		return 0, fmt.Errorf("aecspike: GetCurrentPadding hr=%s", hrHex(hr))
	}
	if s.fillLogs < 10 {
		s.fillLogs++
		slog.Debug("aecspike rig fill", "n", s.fillLogs, "padding", padding,
			"buffer", bufferSize, "written", *written)
	}
	ch := int(mix.channels)
	avail := int(bufferSize) - int(padding)
	if avail <= 0 {
		return 0, nil
	}
	var data unsafe.Pointer
	var gbHr uintptr
	gbHr, _, _ = comCall(s.service, 3, uintptr(avail), uintptr(unsafe.Pointer(&data)))
	if !hrOK(gbHr) {
		// AUDCLNT_E_BUFFER_TOO_LARGE (0x88890007): right after Start the
		// padding counter can lag the fill state; treat it as "no room this
		// period" and wait for the next event instead of failing the run.
		if uint32(gbHr) == 0x88890007 {
			return 0, nil
		}
		return 0, fmt.Errorf("aecspike: GetBuffer(render) hr=%s", hrHex(gbHr))
	}
	n := avail
	if left := (len(inter) - *written) / ch; n > left {
		n = left
	}
	if n <= 0 {
		comCall(s.service, 4, 0)
		return 0, nil
	}
	if mix.tag == waveFormatFloat {
		out := unsafe.Slice((*float32)(data), n*ch)
		copy(out, inter[*written:*written+n*ch])
	} else {
		out := unsafe.Slice((*int16)(data), n*ch)
		for k := 0; k < n*ch; k++ {
			v := inter[*written+k] * 32767
			if v > 32767 {
				v = 32767
			} else if v < -32768 {
				v = -32768
			}
			out[k] = int16(v)
		}
	}
	if rbHr, _, _ := comCall(s.service, 4, uintptr(n), 0); !hrOK(rbHr) { // ReleaseBuffer(n, 0 flags)
		if s.fillLogs <= 10 {
			slog.Info("aecspike rig ReleaseBuffer failed", "hr", hrHex(rbHr), "n", n)
		}
		return 0, fmt.Errorf("aecspike: ReleaseBuffer(render) hr=%s", hrHex(rbHr))
	}
	*written += n * ch
	return n, nil
}

// drain pulls every pending loopback packet down to mono int16
// (IAudioCaptureClient: 3 GetBuffer, 4 ReleaseBuffer, 5 GetNextPacketSize).
func (s *stream) drain(out *[]int16, mix waveFormat) (int, error) {
	ch := int(mix.channels)
	floating := mix.tag == waveFormatFloat
	total := 0
	for {
		var packet uint32
		if hr, _, _ := comCall(s.service, 5, uintptr(unsafe.Pointer(&packet))); !hrOK(hr) {
			return total, fmt.Errorf("aecspike: GetNextPacketSize hr=%s", hrHex(hr))
		}
		if packet == 0 {
			return total, nil
		}
		var (
			data          unsafe.Pointer
			frames, flags uint32
			devPos, qpc   uint64
		)
		if hr, _, _ := comCall(s.service, 3, uintptr(unsafe.Pointer(&data)),
			uintptr(unsafe.Pointer(&frames)), uintptr(unsafe.Pointer(&flags)),
			uintptr(unsafe.Pointer(&devPos)), uintptr(unsafe.Pointer(&qpc))); !hrOK(hr) {
			return total, fmt.Errorf("aecspike: GetBuffer(capture) hr=%s", hrHex(hr))
		}
		if s.fillLogs < 15 {
			s.fillLogs++
			hex := ""
			if data != nil && frames > 0 {
				raw := unsafe.Slice((*byte)(data), 8)
				hex = fmt.Sprintf("%02x %02x %02x %02x %02x %02x %02x %02x",
					raw[0], raw[1], raw[2], raw[3], raw[4], raw[5], raw[6], raw[7])
			}
			slog.Debug("aecspike rig capture packet", "n", s.fillLogs, "frames", frames,
				"flags", flags, "first8bytes", hex, "qpc", qpc)
		}
		if frames > 0 && flags&flagSilent == 0 && data != nil {
			if floating {
				in := unsafe.Slice((*float32)(data), int(frames)*ch)
				// Observed on this host (Windows 11 24H2, Realtek APO chain):
				// the loopback tap hands out float32 samples whose VALUES are
				// the int16 mix samples' BIT PATTERNS, i.e. scaled by 2^-149
				// instead of 2^-15 - every sample a denormal (evidence: the
				// first-8-byte hex dumps in the spike logs). Renormalize by
				// 2^134 (= 2^-15 / 2^-149) when the denormal signature shows.
				if !s.denormChecked {
					for _, v := range in {
						a := v
						if a < 0 {
							a = -a
						}
						if a > 1e-30 {
							s.denormChecked = true
							break
						}
						if a > 0 {
							s.denormChecked = true
							s.denormScale = true
							slog.Debug("aecspike rig loopback denormal signature detected; rescaling by 2^134")
							break
						}
					}
				}
				for i := 0; i < int(frames); i++ {
					var v float32
					for c := 0; c < ch; c++ {
						v += in[i*ch+c]
					}
					v /= float32(ch)
					if s.denormScale {
						v = float32(math.Ldexp(float64(v), 134))
					}
					*out = append(*out, pcmClip(v*32767))
				}
			} else {
				in := unsafe.Slice((*int16)(data), int(frames)*ch)
				for i := 0; i < int(frames); i++ {
					var sum int32
					for c := 0; c < ch; c++ {
						sum += int32(in[i*ch+c])
					}
					*out = append(*out, int16(sum/int32(ch)))
				}
			}
			total += int(frames)
		} else {
			for i := 0; i < int(frames); i++ {
				*out = append(*out, 0)
			}
			total += int(frames)
		}
		comCall(s.service, 4, uintptr(frames))
	}
}

func pcmClip(v float32) int16 {
	if v > 32767 {
		return 32767
	}
	if v < -32768 {
		return -32768
	}
	return int16(v)
}

// close stops and releases the stream. Idempotent, best effort.
func (s *stream) close() {
	comRelease(&s.service)
	if s.client != nil {
		comCall(s.client, 11) // Stop
		comRelease(&s.client)
	}
	if s.ev != 0 {
		_ = windows.CloseHandle(s.ev)
		s.ev = 0
	}
}
