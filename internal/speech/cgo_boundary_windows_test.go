//go:build windows

package speech

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unsafe"

	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
	"golang.org/x/sys/windows"
)

// The cgo boundary evidence lives in tests, not in the product API: the claims
// ("one onnxruntime instance for the whole chain", "intra_op pinned to 1") are
// properties of the process, so they are observed rather than asserted in code.
//
// Why these observables are effective, and what they cannot see, is written up
// in docs/evidence/s1/15b-r1-impl.md. In short: a second runtime image is a
// second loaded module (observable below), and an intra-op pool is a thread
// count (observable below); the Ort::Env refcount inside C is NOT observable
// from Go, and nothing here claims to see it.

// loadedModule is one image mapped into this process.
type loadedModule struct {
	base uintptr
	path string
}

// loadedModules enumerates the process's loaded images.
func loadedModules(t *testing.T) []loadedModule {
	t.Helper()
	proc, err := windows.GetCurrentProcess()
	if err != nil {
		t.Fatalf("GetCurrentProcess: %v", err)
	}
	buf := make([]windows.Handle, 2048)
	var needed uint32
	cb := uint32(len(buf)) * uint32(unsafe.Sizeof(buf[0]))
	if err := windows.EnumProcessModules(proc, &buf[0], cb, &needed); err != nil {
		t.Fatalf("EnumProcessModules: %v", err)
	}
	n := int(needed / uint32(unsafe.Sizeof(buf[0])))
	if n > len(buf) {
		n = len(buf)
	}
	name := make([]uint16, windows.MAX_PATH)
	out := make([]loadedModule, 0, n)
	for _, m := range buf[:n] {
		if err := windows.GetModuleFileNameEx(proc, m, &name[0], uint32(len(name))); err != nil {
			continue
		}
		out = append(out, loadedModule{base: uintptr(m), path: windows.UTF16ToString(name)})
	}
	return out
}

// imagesMatching returns the loaded images whose base name contains want.
func imagesMatching(mods []loadedModule, want string) []loadedModule {
	var out []loadedModule
	for _, m := range mods {
		if strings.Contains(strings.ToLower(filepath.Base(m.path)), want) {
			out = append(out, m)
		}
	}
	return out
}

// processThreads counts this process's threads.
func processThreads(t *testing.T) int {
	t.Helper()
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err != nil {
		t.Fatalf("snapshot threads: %v", err)
	}
	defer func() { _ = windows.CloseHandle(snap) }()
	pid := windows.GetCurrentProcessId()
	var e windows.ThreadEntry32
	e.Size = uint32(unsafe.Sizeof(e))
	n := 0
	for err := windows.Thread32First(snap, &e); err == nil; err = windows.Thread32Next(snap, &e) {
		if e.OwnerProcessID == pid {
			n++
		}
	}
	return n
}

// processCPUTime is kernel+user CPU consumed by the whole process.
func processCPUTime(t *testing.T) time.Duration {
	t.Helper()
	proc, err := windows.GetCurrentProcess()
	if err != nil {
		t.Fatalf("GetCurrentProcess: %v", err)
	}
	var created, exited, kernel, user windows.Filetime
	if err := windows.GetProcessTimes(proc, &created, &exited, &kernel, &user); err != nil {
		t.Fatalf("GetProcessTimes: %v", err)
	}
	return time.Duration(kernel.Nanoseconds() + user.Nanoseconds())
}

// rawVAD opens a silero session directly, bypassing Detector, so the thread
// count can be explored outside the pinned value. Production cannot reach this
// path: NewDetector always passes IntraOpThreads.
func rawVAD(t *testing.T, numThreads int) *sherpa.VoiceActivityDetector {
	t.Helper()
	v := sherpa.NewVoiceActivityDetector(&sherpa.VadModelConfig{
		SileroVad: sherpa.SileroVadModelConfig{
			Model:              sileroModelPath(t),
			Threshold:          DefaultVadThreshold,
			MinSilenceDuration: 0,
			MinSpeechDuration:  0,
			WindowSize:         WindowSamples,
			MaxSpeechDuration:  float32(MaxUtterance.Seconds()),
		},
		SampleRate: 16000,
		NumThreads: numThreads,
		Provider:   "cpu",
	}, vadBufferSeconds)
	if v == nil {
		t.Skipf("sherpa refused a %d-thread VAD session", numThreads)
	}
	return v
}

// ortfImageName is the module base name the whole chain must map exactly once.
const ortfImageName = "onnxruntime"

// TestSingleOnnxruntimeInstance proves the chain runs on ONE onnxruntime image.
//
// Leg 1 counts the images before any session exists (the sherpa import already
// pulled the runtime in at load time).
// Leg 2 creates and disposes 8 detectors and counts again: a second runtime
// would show up as a second image with a second path.
// Leg 3 is the control that keeps leg 1-2 from being self-certifying: it maps a
// copy of the same DLL under another name and shows the counter DOES see two.
func TestSingleOnnxruntimeInstance(t *testing.T) {
	mods := loadedModules(t)
	ort := imagesMatching(mods, ortfImageName)
	sherpaMods := imagesMatching(mods, "sherpa-onnx")
	t.Logf("leg1: %d image(s) matching %q, %d sherpa image(s)", len(ort), ortfImageName, len(sherpaMods))
	for _, m := range ort {
		t.Logf("leg1: onnxruntime base=%#x path=%s", m.base, m.path)
	}
	if len(ort) != 1 {
		t.Fatalf("leg1: %d onnxruntime image(s) mapped, want exactly 1: %v", len(ort), ort)
	}
	firstPath := ort[0].path
	if len(sherpaMods) == 0 {
		t.Fatalf("leg1: no sherpa image mapped, so the C API is not even linked")
	}

	if err := checkOnnxruntime(); err != nil {
		t.Skipf("cannot create sessions here (%v); legs 2-3 need a working runtime", err)
	}

	// Leg 2: eight create/dispose cycles must not map a second runtime.
	for i := 0; i < 8; i++ {
		d, err := NewDetector(DetectorConfig{ModelPath: sileroModelPath(t)})
		if err != nil {
			t.Fatalf("leg2: detector %d: %v", i, err)
		}
		if err := d.Close(); err != nil {
			t.Fatalf("leg2: close %d: %v", i, err)
		}
	}
	after := imagesMatching(loadedModules(t), ortfImageName)
	t.Logf("leg2: after 8 create/dispose cycles: %d onnxruntime image(s)", len(after))
	if len(after) != 1 {
		t.Fatalf("leg2: %d onnxruntime image(s) after 8 sessions, want 1: %v", len(after), after)
	}
	if !strings.EqualFold(after[0].path, firstPath) {
		t.Fatalf("leg2: runtime image changed from %s to %s", firstPath, after[0].path)
	}

	// Leg 3: instrument control. Map a byte-identical copy under another name;
	// the counter must report two distinct images, or legs 1-2 prove nothing.
	src := ""
	for _, m := range imagesMatching(loadedModules(t), ortfImageName) {
		if strings.EqualFold(filepath.Base(m.path), "onnxruntime.dll") {
			src = m.path
		}
	}
	if src == "" {
		t.Fatalf("leg3: cannot find the onnxruntime.dll path to copy")
	}
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("leg3: read %s: %v", src, err)
	}
	probe := filepath.Join(t.TempDir(), "onnxruntime_probe_copy.dll")
	if err := os.WriteFile(probe, data, 0o600); err != nil {
		t.Fatalf("leg3: write %s: %v", probe, err)
	}
	h, err := windows.LoadLibraryEx(probe, 0, 0)
	if err != nil {
		t.Fatalf("leg3: LoadLibraryEx(%s): %v", probe, err)
	}
	two := imagesMatching(loadedModules(t), ortfImageName)
	t.Logf("leg3: with the probe copy mapped: %d onnxruntime image(s)", len(two))
	if len(two) != 2 {
		t.Fatalf("leg3: counter saw %d image(s) with a second copy mapped, want 2 - "+
			"the instrument is blind and legs 1-2 are worthless", len(two))
	}
	if err := windows.FreeLibrary(h); err != nil {
		t.Logf("leg3: FreeLibrary: %v", err)
	}
}

// TestIntraOpThreadShape measures whether the pinned intra-op count is visible.
//
// Two independent observables, both taken around the same fixed workload:
//   - thread delta while the session exists (an intra-op pool is threads);
//   - CPU seconds burned per wall second (with one intra-op thread at most one
//     core can be busy, so the ratio stays near 1; a pool pushes it up).
//
// The 4-thread session is the positive control: if neither observable separates
// 1 from 4 on this host, the measurement has no teeth and the receipt says so
// instead of claiming the pin is proven.
func TestIntraOpThreadShape(t *testing.T) {
	if err := checkOnnxruntime(); err != nil {
		t.Skipf("cannot create sessions here (%v)", err)
	}
	const windowsToFeed = 1000
	audioBuf := pcmToFloat(synthSpeech(msSamples(32)))

	type shape struct {
		threads int
		cpu     time.Duration
		wall    time.Duration
	}
	measure := func(numThreads int) shape {
		base := processThreads(t)
		v := rawVAD(t, numThreads)
		defer sherpa.DeleteVoiceActivityDetector(v)
		// Touch the session so any lazily created pool exists before counting.
		v.AcceptWaveform(audioBuf)
		_ = v.IsSpeech()
		peak := processThreads(t)
		cpu0, wall0 := processCPUTime(t), time.Now()
		for i := 0; i < windowsToFeed; i++ {
			v.AcceptWaveform(audioBuf)
		}
		return shape{
			threads: peak - base,
			cpu:     processCPUTime(t) - cpu0,
			wall:    time.Since(wall0),
		}
	}

	one := measure(IntraOpThreads)
	four := measure(4)
	ratio := func(s shape) float64 { return float64(s.cpu) / float64(s.wall) }
	t.Logf("intra_op=1: threads+%d cpu=%v wall=%v cpu/wall=%.2f",
		one.threads, one.cpu.Round(time.Millisecond), one.wall.Round(time.Millisecond), ratio(one))
	t.Logf("intra_op=4: threads+%d cpu=%v wall=%v cpu/wall=%.2f",
		four.threads, four.cpu.Round(time.Millisecond), four.wall.Round(time.Millisecond), ratio(four))

	// The pinned shape must not look like a fan-out: one core's worth of CPU.
	if r := ratio(one); r > 1.50 {
		t.Fatalf("intra_op=1 burned %.2f cores (cpu/wall), so the pin is not holding", r)
	}
	// The control must separate, or the check above is vacuous.
	if four.threads <= one.threads && ratio(four) <= ratio(one) {
		t.Logf("WARNING: neither thread delta nor cpu/wall separates 1 from 4 on this host "+
			"(threads %d vs %d, ratio %.2f vs %.2f); the intra-op pin is NOT independently "+
			"confirmed by this instrument here", one.threads, four.threads, ratio(one), ratio(four))
	}
	if one.threads > 4 {
		t.Fatalf("intra_op=1 created %d threads, that is a pool", one.threads)
	}
}
