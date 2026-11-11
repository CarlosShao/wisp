//go:build windows

package aecspike

// Ticket 297 AC#2 form ① discriminator: is "the loud phase is not
// distinguishable from silence" (ticket 247's surviving red sentence) the
// shape where the sound never reaches the microphone?
//
// What ticket 297-l1 established (probes/297/l1/): the old 0.37 floor was
// form ② (float32 packets sliced as int16) and ticket 300's @24 fix cured
// it; form ③ (a ~0.37 device floor) is excluded by the post-fix live
// readings (quiet 0.0013-0.0035); what remains red is form ① -- the
// six-window sound window read mean=0.001358, BELOW its quiet-A, i.e. the
// microphone heard silence while the announce render was playing. The ticket
// names the candidate facts: speakers off, or routed to an endpoint this
// microphone cannot hear. Those are facts about THIS machine at THIS
// moment -- most of them are SOFTWARE-READABLE, and the reach question
// itself is answerable in-process with one controlled experiment:
//
//	phase A   the production microphone (real WASAPIMicrophone, real capture
//	          loop, per-frame LevelSink) records ~6 quiet seconds;
//	phase B   the rig renders FIVE one-second known-amplitude speech-like
//	          fixtures on the console default render endpoint; each render's
//	          own endpoint-loopback tap (RunLoopback's captured mix) is the
//	          amplitude TRUTH of what the endpoint feed carried;
//	phase C   ~2 more quiet seconds as the recovery face.
//
// and prints, next to the microphone's quiet/render/after means:
//
//	[volume] the render session's ISimpleAudioVolume master volume and mute,
//	         read through a client opened the way the ticket-388 guard opens
//	         its (GetService on an initialized, never-started client);
//	[hardware] the default render endpoint's IAudioEndpointVolume master
//	         volume and mute.
//
// Both reads are OUT-PARAMETER reads: GetMasterVolume(float*), GetMute(BOOL*)
// on both interfaces. A float passed BY POINTER through syscall.SyscallN is
// safe -- the poisoned-persisted-volume incident the 388 evidence file
// records happened only on the WRITE path, where a float VALUE must ride
// XMM; no volume is written here. This is why the instrument works under
// CGO_ENABLED=0, where the 388 C-shim write path refuses: reading needs no
// shim.
//
// Sampling order, stated because it bounds the reading: the microphone's
// render-window face is collected BEFORE its paired render (a one-second
// window immediately following the previous second of render) and the rig
// run happens AFTER the window closes, so a window can never sample its OWN
// render; consecutive runs back-to-back make window N+1 cover render N
// almost exactly. The tap is the render-side truth; the mic window is the
// acoustic answer one render behind. A future run that wants exact overlap
// needs a sink-side window mark, which is a production-seam change this
// wave refuses (internal/audio/** is l1's written face; this file is
// aecspike-only).
//
// Verdict sentences are PRINTED, never asserted: whether the microphone
// heard the render is a fact question whose answer is evidence for the
// orchestrator's AC#2 ruling, not a pass/fail of this file. What IS asserted
// is the instrument's own premise (the 388 precedent): the endpoint tap must
// have carried the render, the microphone must stream at the seam's ~31
// levels/s, and rig/device failures skip the way 388 skips.
//
// Gate: WISP_AEC_SPIKE=1 exactly like the rest of the gated package. The
// test lives in aecspike -- NOT internal/audio -- deliberately: the
// portable-tests windows tier runs ./internal/audio/ with SKIP-is-fatal and
// a ledger this wave must not touch, so an env-gated skip there would be a
// CI red; aecspike is outside every CI tier's scope (census answer, ticket
// 297 AC#5 reading), so this instrument is a manual run by construction.

import (
	"context"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/observe"
)

// COM identifiers this instrument reads through (audioclient.h /
// mmdeviceapi.h / endpointvolume.h). Reads only: every call below takes an
// out pointer, none passes a float value.
var (
	// IID_ISimpleAudioVolume {87CE5498-68D6-44E5-9215-6DA47EF883D8}
	iidISimpleAudioVolume = windows.GUID{
		Data1: 0x87CE5498, Data2: 0x68D6, Data3: 0x44E5,
		Data4: [8]byte{0x92, 0x15, 0x6D, 0xA4, 0x7E, 0xF8, 0x83, 0xD8},
	}
	// IID_IAudioEndpointVolume {5CDF2C82-841E-4546-9722-0CF74078229A}
	iidIAudioEndpointVolume = windows.GUID{
		Data1: 0x5CDF2C82, Data2: 0x841E, Data3: 0x4546,
		Data4: [8]byte{0x97, 0x22, 0x0C, 0xF7, 0x40, 0x78, 0x22, 0x9A},
	}
)

// ISimpleAudioVolume vtable (after IUnknown 0-2): 3
// SetMasterVolume(float, GUID*) -- the XMM-value call this instrument must
// never drive from Go (388's finding), 4 GetMasterVolume(float*),
// 5 SetMute(BOOL, GUID*), 6 GetMute(BOOL*).
const (
	slotSimpleGetMasterVolume = 4
	slotSimpleGetMute         = 6
)

// IAudioEndpointVolume vtable (after IUnknown 0-2), endpointvolume.h MIDL
// order: 4 RegisterControlChangeNotify, 5 UnregisterControlChangeNotify,
// 6 GetChannelCount, 7 SetMasterVolumeLevel, 8 SetMasterVolumeLevelScalar,
// 9 GetMasterVolumeLevelScalar, 10 SetMasterVolumeLevel, 11
// GetMasterVolumeLevel, 12 SetVolumeStepInfo(ish pair), 13
// GetVolumeStepInfo, 14 SetMute, 15 GetMute.
const (
	slotEndpointGetMasterVolumeScalar = 9
	slotEndpointGetMute               = 15
)

// Vtable slots shared with the existing package instruments:
// IMMDevice::Activate = 3, IMMDeviceEnumerator::GetDefaultAudioEndpoint = 4,
// IAudioClient::GetService = 14 (rig_windows.go's mapping).
const (
	slotIMMDeviceActivate           = 3
	slotEnumGetDefaultAudioEndpoint = 4
	slotIAudioClientGetService      = 14
)

// formOneMicLevels collects per-frame levels from the production microphone
// under a mutex: the sink runs on the pinned capture thread while the test
// goroutine reads snapshots.
type formOneMicLevels struct {
	mu     sync.Mutex
	levels []float32
}

func (f *formOneMicLevels) call(v float32) {
	f.mu.Lock()
	f.levels = append(f.levels, v)
	f.mu.Unlock()
}

func (f *formOneMicLevels) snapshot() []float32 {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]float32(nil), f.levels...)
}

func (f *formOneMicLevels) snapshotLen() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.levels)
}

// readSimpleAudioVolume opens a session-volume view of the render session
// (the 388 guard's open shape) and reads master volume + mute. Read-only.
func readSimpleAudioVolume() (float64, bool, error) {
	client, release, cerr := openSessionVolumeClient()
	if cerr != nil {
		return 0, false, cerr
	}
	defer release()
	var simple unsafe.Pointer
	if hr, _, _ := comCall(client, slotIAudioClientGetService,
		uintptr(unsafe.Pointer(&iidISimpleAudioVolume)),
		uintptr(unsafe.Pointer(&simple))); !hrOK(hr) || simple == nil {
		return 0, false, fmt.Errorf("GetService(ISimpleAudioVolume) hr=%s", hrHex(hr))
	}
	defer comRelease(&simple)
	var vol float32
	if hr, _, _ := comCall(simple, slotSimpleGetMasterVolume,
		uintptr(unsafe.Pointer(&vol)), 0); !hrOK(hr) {
		return 0, false, fmt.Errorf("GetMasterVolume hr=%s", hrHex(hr))
	}
	var mute int32
	if hr, _, _ := comCall(simple, slotSimpleGetMute,
		uintptr(unsafe.Pointer(&mute)), 0); !hrOK(hr) {
		return 0, false, fmt.Errorf("GetMute hr=%s", hrHex(hr))
	}
	return float64(vol), mute != 0, nil
}

// readEndpointVolume reads the default render endpoint's hardware master
// volume + mute through IAudioEndpointVolume. Read-only.
func readEndpointVolume(enum unsafe.Pointer) (float64, bool, error) {
	var dev unsafe.Pointer
	if hr, _, _ := comCall(enum, slotEnumGetDefaultAudioEndpoint, eRender, roleConsole,
		uintptr(unsafe.Pointer(&dev))); !hrOK(hr) || dev == nil {
		return 0, false, fmt.Errorf("no default render endpoint (hr=%s)", hrHex(hr))
	}
	defer comRelease(&dev)
	var epv unsafe.Pointer
	if hr, _, _ := comCall(dev, slotIMMDeviceActivate,
		uintptr(unsafe.Pointer(&iidIAudioEndpointVolume)), clsctxAll, 0,
		uintptr(unsafe.Pointer(&epv))); !hrOK(hr) || epv == nil {
		return 0, false, fmt.Errorf("Activate(IAudioEndpointVolume) hr=%s", hrHex(hr))
	}
	defer comRelease(&epv)
	var vol float32
	if hr, _, _ := comCall(epv, slotEndpointGetMasterVolumeScalar,
		uintptr(unsafe.Pointer(&vol))); !hrOK(hr) {
		return 0, false, fmt.Errorf("GetMasterVolumeLevelScalar hr=%s", hrHex(hr))
	}
	var mute int32
	if hr, _, _ := comCall(epv, slotEndpointGetMute,
		uintptr(unsafe.Pointer(&mute))); !hrOK(hr) {
		return 0, false, fmt.Errorf("GetMute hr=%s", hrHex(hr))
	}
	return float64(vol), mute != 0, nil
}

// waitForLevelsUntil returns the sink snapshot once want more levels arrived
// or the deadline passed (the l1 livepacket deadline-loop shape: a test-side
// bound, never a wall-clock production timeout).
func waitForLevelsUntil(rec *formOneMicLevels, have, want int, deadline time.Time) []float32 {
	for {
		got := rec.snapshot()
		if len(got)-have >= want || !time.Now().Before(deadline) {
			return got
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// meanRangeStats is one phase's face: count, mean, max, min.
type meanRangeStats struct {
	n    int
	mean float64
	max  float64
	min  float64
}

// renderFace pairs one render's endpoint-tap truth with the microphone
// window that covers it (see the sampling-order note at the top).
type renderFace struct {
	run      int
	tapMean  float64
	tapPeak  float64
	micN     int
	micMean  float64
	micMax   float64
	pushed   int
	captured int
}

func statsOf(levels []float32) meanRangeStats {
	var s meanRangeStats
	if len(levels) == 0 {
		return s
	}
	mn, mx := math.MaxFloat64, -math.MaxFloat64
	var sum float64
	for _, v := range levels {
		f := float64(v)
		sum += f
		if f > mx {
			mx = f
		}
		if f < mn {
			mn = f
		}
	}
	return meanRangeStats{n: len(levels), mean: sum / float64(len(levels)), max: mx, min: mn}
}

// TestAC297FormOneRenderReachProbe runs the quiet/render/after microphone
// window against five known-amplitude endpoint renders, prints the reach
// verdict with the volume/mute snapshot beside it. Gated manual instrument;
// it never writes a volume and never asserts an AC#2 disposition.
func TestAC297FormOneRenderReachProbe(t *testing.T) {
	spikeGate(t)
	reg := observe.NewRegistry()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const (
		quietSeconds = 6
		renderRuns   = 5
		afterSeconds = 2
		gain         = 0.35 // fixture sample gain; loudness axis = gain x session volume
	)

	// [read] the volume snapshot BEFORE anything is touched (a fact about
	// the machine as the owner left it). COM init lives HERE (the test
	// goroutine): the rig's own thread initializes its own apartment, and
	// EnumerateEndpoints is the package precedent for a balanced init on a
	// plain goroutine.
	if comInit := comInitMTA(); comInit {
		defer comUninitMTA()
	}
	enum, releaseEnum, err := createEnumerator()
	if err != nil {
		t.Skipf("enumerator unavailable: %v", err)
	}
	defer releaseEnum()
	sessVol, sessMute, err := readSimpleAudioVolume()
	if err != nil {
		t.Skipf("render session volume unavailable: %v", err)
	}
	hwVol, hwMute, err := readEndpointVolume(enum)
	if err != nil {
		t.Skipf("endpoint volume unavailable: %v", err)
	}
	t.Logf("[volume] render session ISimpleAudioVolume: volume=%.3f muted=%v", sessVol, sessMute)
	t.Logf("[volume] render endpoint IAudioEndpointVolume: volume=%.3f muted=%v", hwVol, hwMute)

	// [mic] phase A: the production microphone, quiet room.
	rec := &formOneMicLevels{}
	mic := audio.NewWASAPIMicrophone(
		audio.WithCaptureRegistry(reg),
		audio.WithLevelSink(rec.call),
	)
	buf := make(chan []byte, 1024) // absorbs the whole run; no consumer prints drop warnings
	if err := mic.Start(ctx, buf); err != nil {
		t.Skipf("production microphone unavailable on this host: %v", err)
	}
	quietFace := func() meanRangeStats {
		deadline := time.Now().Add(time.Duration(quietSeconds)*time.Second + 500*time.Millisecond)
		got := waitForLevelsUntil(rec, 0, 31*quietSeconds, deadline)
		return statsOf(got)
	}()
	if quietFace.n == 0 {
		_ = mic.Stop()
		t.Fatalf("phase A collected no levels (capture loop produced nothing in %ds)", quietSeconds)
	}
	t.Logf("[mic A] quiet %ds: n=%d mean=%.6f max=%.6f min=%.6f (~%d levels/s)",
		quietSeconds, quietFace.n, quietFace.mean, quietFace.max, quietFace.min, quietFace.n/quietSeconds)

	// [render] phase B: five known-amplitude fixtures; the endpoint loopback
	// tap is the truth of what the endpoint feed carried, run by run. The
	// mic window is sampled FIRST (see the sampling-order note above): it
	// covers the PREVIOUS run's render; run 1's window is a second quiet
	// control inside phase B. A rig failure mid-run is recorded and stops
	// the loop so the verdict below still prints on what ran.
	far := SynthSpeechLike(audio.TargetRate, 1)
	var renders []renderFace
	rigFailed := false
	for run := 1; run <= renderRuns && !rigFailed; run++ {
		// One-second microphone window (covers the previous run's render).
		have := rec.snapshotLen()
		got := waitForLevelsUntil(rec, have, 31, time.Now().Add(1500*time.Millisecond))
		face := statsOf(got[have:])

		// The render, after the window closed.
		r, mix, rerr := RunLoopback(reg, ctx, "", gain, far, audio.TargetRate)
		if rerr != nil {
			t.Logf("[render %d] rig failed: %v (verdict below is computed on runs 1..%d)", run, rerr, run-1)
			rigFailed = true
			break
		}
		prms := peakWindowRMS(mix, 1024)
		renders = append(renders, renderFace{
			run: run, tapMean: rmsInt16(mix), tapPeak: prms,
			micN: face.n, micMean: face.mean, micMax: face.max,
			pushed: r.PushedFrames, captured: r.CapturedFrames,
		})
		t.Logf("[render %d] endpoint tap: rms=%.6f peakrms=%.6f pushed=%d captured=%d | mic window (covers render %d): n=%d mean=%.6f max=%.6f",
			run, rmsInt16(mix), prms, r.PushedFrames, r.CapturedFrames, run-1, face.n, face.mean, face.max)
	}

	// [mic] phase C: recovery face (quiet again; also covers render 5).
	have := rec.snapshotLen()
	got := waitForLevelsUntil(rec, have, 31*afterSeconds, time.Now().Add(time.Duration(afterSeconds)*time.Second+500*time.Millisecond))
	afterFace := statsOf(got[have:])
	_ = mic.Stop()

	// The reach verdict is computed BEFORE the premise checks so a rig
	// failure mid-run (the 297-l2 run-4 shape: the rig's render fill panicked
	// after three clean runs) still leaves the discriminating sentences on
	// the record; the premise checks then decide pass, skip or fail.
	quietRef := quietFace.mean
	if quietRef <= 0 {
		quietRef = 1e-6
	}
	best := 0.0
	for _, f := range renders {
		if f.micMean > best {
			best = f.micMean
		}
	}
	ratio := best / quietRef
	verdict := func() {
		switch {
		case ratio >= 2.0:
			t.Logf("[verdict] the microphone DID hear the render (best render-covering mean %.6f = %.2fx quiet %.6f): form ① is NOT the whole story this run -- loudness/routing detail is next",
				best, ratio, quietFace.mean)
		case ratio >= 1.3:
			t.Logf("[verdict] marginal: the microphone may have heard the render (best %.6f = %.2fx quiet %.6f); repeat the run before ruling",
				best, ratio, quietFace.mean)
		default:
			t.Logf("[verdict] the microphone did NOT hear the render (best %.6f = %.2fx quiet %.6f, after-face mean %.6f): form ① SUPPORTED this run -- the sound did not reach the capture device while the endpoint feed carried it (taps above); the fact question is WHY: session volume=%.3f muted=%v, endpoint volume=%.3f muted=%v, or physical speakers/headphones off -- the owner's ears settle it",
				best, ratio, quietFace.mean, afterFace.mean, sessVol, sessMute, hwVol, hwMute)
		}
		t.Logf("[summary] quiet mean=%.6f | best render-covering mean=%.6f (ratio %.2f) | after mean=%.6f | endpoint taps carried %d/%d at fixture gain=%g",
			quietFace.mean, best, ratio, afterFace.mean, tapOK(renders), len(renders), gain)
	}

	// The instrument's own premise (asserted; the verdict is not): the
	// endpoint tap must have carried the render -- every run's peak window
	// must sit far above a silent tap (~0 counts).
	if tapOK(renders) != len(renders) {
		verdict()
		// The tap is fed through the session/hardware volume axis: a silent
		// tap with a muted/zero volume axis is itself the form-① answer, not
		// a broken instrument -- report and skip rather than fail.
		if hwMute || hwVol <= 0.001 || sessMute || sessVol <= 0.001 {
			t.Skipf("the endpoint feed itself was silent under the current volume axis (session %.3f muted=%v, endpoint %.3f muted=%v): that IS a form-① finding; raise the volume and re-run for the reach half",
				sessVol, sessMute, hwVol, hwMute)
		}
		t.Fatalf("instrument premise: endpoint taps did not carry the render (fixture gain=%g, session %.3f, endpoint %.3f); see the verdict above",
			gain, sessVol, hwVol)
	}
	if len(renders) == 0 || renders[0].captured == 0 {
		t.Fatalf("instrument premise: the tap captured no frames; the loopback reading is invalid")
	}

	// [mic] premise: the production microphone streamed at the seam rate.
	if quietFace.n < 31*quietSeconds/2 {
		t.Fatalf("instrument premise: the microphone produced only %d levels in %ds (want ~%d); the reach question is unanswerable this run",
			quietFace.n, quietSeconds, 31*quietSeconds)
	}

	verdict()
	if rigFailed {
		t.Skipf("the rig stopped early (see the [render] line above); the verdict above is the evidence this run produced")
	}
}

// tapOK counts the renders whose endpoint tap provably carried audio
// (peak window far above a silent tap's ~0 counts).
func tapOK(renders []renderFace) int {
	var n int
	for _, f := range renders {
		if f.tapPeak > 200 { // counts; a silent tap reads ~0
			n++
		}
	}
	return n
}
