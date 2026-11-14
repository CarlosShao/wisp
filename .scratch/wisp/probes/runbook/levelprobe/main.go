//go:build windows

// Command levelprobe measures this machine's microphone level with the SAME
// ruler the product uses, without starting, stopping or otherwise touching the
// running Wisp process (ticket 290 / 291 / 449 / 247 "one real-machine session"
// runbook).
//
// Why this program exists at all. The resident leg never prints a level value:
// resident_audio_windows.go's only level-shaped output is the exit line
// "levels_delivered=%d" (there:429), which needs the product process to stop -
// and the runbook forbids stopping it. The product side is therefore read from
// the persistent log sink (records, with line numbers), and the NUMBERS come
// from here, off the same default capture endpoint and through the same
// internal/audio functions the product calls:
//
//	audio.NewWASAPIMicrophone  - the production capture source (wasapimic)
//	audio.DecodeFrame          - the C8 seam decoder
//	audio.LevelOfSamples       - the production level scale (ticket 449's
//	                             DC-stripped definition; level.go:119)
//
// Two numbers are computed per frame, on the SAME bytes, so that ticket 449's
// "it must only ever read lower" claim is measured rather than asserted:
//
//	new = LevelOfSamples(samples)          (shipping ruler, DC removed)
//	old = sqrt(mean(s^2)) / 32768          (pre-449 ruler, DC included)
//
// The proof that new <= old for every frame is in internal/audio/level.go: the
// mean is the minimiser of the squared deviation, so removing it can only
// shrink the RMS. This program prints the measured delta; it does not assume it.
//
// What this program is NOT: it is not the product, it is not a test of the
// product, and it never asserts that "the ball breathes" (prototypeVisuals is
// off; ticket 247 AC#10). It is a measuring instrument that happens to be
// written in Go because the ruler lives in Go.
//
// Output, all under -out:
//
//	levelprobe.json        summary (one object; absent values are null, never 0)
//	levelprobe.frames.jsonl one record per C8 frame: phase, t_ms, new, old, dc
//	levelprobe.stdout.log  this program's own narration (the script tees it too)
//	tone-1k-halffs.wav     the played tone of window 3 (0.5 FS, 1 kHz, 16 kHz mono)
//
// Exit codes: 0 = every window recorded frames; 2 = the microphone never
// opened (nothing measured, and the JSON says so); 3 = opened but at least one
// window recorded zero frames (that window's numbers stay null).
package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/audio"
)

// ---------------------------------------------------------------------------
// measurement shapes
// ---------------------------------------------------------------------------

// windowSpec is one timed window of the session. tone=true makes this program
// play its own WAV through the default render endpoint for the window's length.
type windowSpec struct {
	Name string
	Secs float64
	Say  string // what the operator is asked to do, printed before the window
	Tone bool
}

// defaultWindows is the 17-second measuring part of the runbook: quiet first
// (ticket 449 AC#8 wants the quiet reading), then the voice (291 AC#2 shape 2 /
// 247 AC#2), then quiet again (the control that catches "the room got louder"),
// then the played tone (291 AC#2 shape 3). Shape 1 (completely quiet) and shape
// 3 are the same instrument; shape 4 is synthetic and needs no microphone.
func defaultWindows() []windowSpec {
	return []windowSpec{
		{Name: "quiet-1", Secs: 4, Say: "完全安静：别说话、别碰键盘（4 秒）"},
		{Name: "speak", Secs: 5, Say: "用正常音量说一句话，一直说到窗口结束（5 秒）"},
		{Name: "quiet-2", Secs: 4, Say: "再安静一次，别出声（4 秒，对照组）"},
		{Name: "tone", Secs: 4, Say: "什么都别说：本程序现在自己放一段 1 kHz 音（4 秒）", Tone: true},
	}
}

// ---------------------------------------------------------------------------
// per-frame record and per-window statistics
// ---------------------------------------------------------------------------

type frameRecord struct {
	Phase string  `json:"phase"`
	TMS   int64   `json:"t_ms"`
	New   float64 `json:"new"`
	Old   float64 `json:"old"`
	DC    float64 `json:"dc_counts"`
}

type nums struct {
	Count int      `json:"count"`
	Mean  *float64 `json:"mean"` // nil => nothing measured (rendered as null, never 0)
	Max   *float64 `json:"max"`
	Min   *float64 `json:"min"`
}

type windowResult struct {
	Name       string   `json:"name"`
	SpanS      float64  `json:"span_s"`
	Samples    int      `json:"samples"`
	LevelsPerS *float64 `json:"levels_per_s"`
	Requested  string   `json:"requested"`
	NewLevel   nums     `json:"new_level_dc_stripped"`
	OldLevel   nums     `json:"old_level_dc_included"`
	DCCounts   nums     `json:"dc_counts"`
	Note       string   `json:"note,omitempty"`
}

type accumulator struct {
	spec windowSpec
	recs []frameRecord
}

func (a *accumulator) add(r frameRecord) { a.recs = append(a.recs, r) }

func statsOf(values []float64) nums {
	if len(values) == 0 {
		return nums{}
	}
	sum, mx, mn := 0.0, values[0], values[0]
	for _, v := range values {
		sum += v
		if v > mx {
			mx = v
		}
		if v < mn {
			mn = v
		}
	}
	mean := sum / float64(len(values))
	return nums{Count: len(values), Mean: &mean, Max: &mx, Min: &mn}
}

func (a *accumulator) result(spanS float64) windowResult {
	newVals := make([]float64, 0, len(a.recs))
	oldVals := make([]float64, 0, len(a.recs))
	dcVals := make([]float64, 0, len(a.recs))
	for _, r := range a.recs {
		newVals = append(newVals, r.New)
		oldVals = append(oldVals, r.Old)
		dcVals = append(dcVals, r.DC)
	}
	res := windowResult{
		Name:      a.spec.Name,
		SpanS:     spanS,
		Samples:   len(a.recs),
		Requested: a.spec.Say,
		NewLevel:  statsOf(newVals),
		OldLevel:  statsOf(oldVals),
		DCCounts:  statsOf(dcVals),
	}
	if spanS > 0 && len(a.recs) > 0 {
		ps := float64(len(a.recs)) / spanS
		res.LevelsPerS = &ps
	}
	return res
}

// ---------------------------------------------------------------------------
// synthetic shape (ticket 291 AC#2 shape 4): no microphone involved
// ---------------------------------------------------------------------------

type synthCase struct {
	Name       string  `json:"name"`
	Frame      string  `json:"frame"`
	Measured   float64 `json:"measured"`
	Expected   float64 `json:"expected"`
	Delta      float64 `json:"delta"`
	Tolerance  float64 `json:"tolerance"`
	Within     bool    `json:"within_tolerance"`
	Derivation string  `json:"derivation"`
}

func squareFrame(lo, hi int16, n int) []int16 {
	s := make([]int16, n)
	for i := range s {
		if i%2 == 0 {
			s[i] = hi
		} else {
			s[i] = lo
		}
	}
	return s
}

// sineFrame builds n samples of round(dc + amp*sin(2*pi*k/period)). A period of
// 64 over 512 samples is exactly 8 whole cycles, so mean(sin^2) is 1/2 and the
// closed form below is exact up to int16 rounding (level.go:79-89).
func sineFrame(dc, amp float64, period, n int) []int16 {
	s := make([]int16, n)
	for i := range s {
		s[i] = int16(math.Round(dc + amp*math.Sin(2*math.Pi*float64(i)/float64(period))))
	}
	return s
}

func synthResults() []synthCase {
	const n = audio.FrameSamples
	level := audio.LevelOfSamples
	out := []synthCase{
		{
			Name: "constant-zero", Frame: "512 x 0",
			Measured: level(make([]int16, n)), Expected: audio.MinLevel,
			Derivation: "level.go:119-135: a constant frame has zero deviation about its own mean",
		},
		{
			Name: "constant-dc-12345", Frame: "512 x 12345",
			Measured: level(func() []int16 {
				s := make([]int16, n)
				for i := range s {
					s[i] = 12345
				}
				return s
			}()),
			Expected:   audio.MinLevel,
			Derivation: "ticket 449 AC#1: pure DC now reads exactly 0; it used to read 12345/32768 = 0.376740",
		},
		{
			Name: "square-1lsb", Frame: "+1/-1 alternating",
			Measured: level(squareFrame(-1, 1, n)), Expected: audio.LevelLSB,
			Derivation: "level.go:40-41: the smallest non-zero reading, 1/32768",
		},
		{
			Name: "square-fullscale", Frame: "+32767/-32767 alternating",
			Measured: level(squareFrame(-32767, 32767, n)), Expected: audio.FullScaleSquareLevel,
			Derivation: "level.go:74-77: FullScaleSquareLevel = 32767/32768",
		},
		{
			Name: "hottest-legal", Frame: "+32767/-32768 alternating",
			Measured: level(squareFrame(-32768, 32767, n)), Expected: audio.MaxLevel - audio.LevelLSB/2,
			Derivation: "level.go:28-34: MaxLevel - LevelLSB/2, the closed upper bound no symmetric frame attains",
		},
		{
			Name: "dc12345-sine6000", Frame: "12345 + 6000*sin, 8 whole cycles",
			Measured: level(sineFrame(12345, 6000, 64, n)), Expected: 6000.0 / audio.LevelFullScale / math.Sqrt2,
			Derivation: "ticket 449 AC#3: DC is removed, so the reading is the sine's own RMS = 6000/32768/sqrt(2)",
		},
	}
	for i := range out {
		out[i].Delta = out[i].Measured - out[i].Expected
		out[i].Tolerance = audio.SineLevelTolerance
		out[i].Within = math.Abs(out[i].Delta) <= out[i].Tolerance
	}
	return out
}

// ---------------------------------------------------------------------------
// tone playback (window 3): the played signal's RMS is known exactly
// ---------------------------------------------------------------------------

const (
	sndAsync     = 0x0001
	sndNoDefault = 0x0002
	sndFilename  = 0x00020000
)

var (
	winmm     = syscall.NewLazyDLL("winmm.dll")
	playSound = winmm.NewProc("PlaySoundW")
)

func playAsync(path string) error {
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	r, _, callErr := playSound.Call(uintptr(unsafe.Pointer(p)), 0, sndAsync|sndFilename|sndNoDefault)
	if r == 0 {
		return callErr
	}
	return nil
}

func stopPlay() { playSound.Call(0, 0, 0) }

// writeToneWAV writes a 16 kHz mono 16-bit PCM sine of the given frequency and
// digital amplitude. Amplitude 16384 is 0.5 of full scale, so the played
// signal's RMS is 16384/32768/sqrt(2) = 0.35355339059327373 of full scale -
// the "known amplitude tone" ticket 291 AC#2 shape 3 asks the operator to play,
// generated here so the amplitude is a fact rather than a memory.
func writeToneWAV(path string, freq, amp float64, secs float64) error {
	const rate = 16000
	n := int(rate * secs)
	data := make([]byte, 0, n*2)
	for i := 0; i < n; i++ {
		v := int16(math.Round(amp * math.Sin(2*math.Pi*freq*float64(i)/rate)))
		var b [2]byte
		binary.LittleEndian.PutUint16(b[:], uint16(v))
		data = append(data, b[0], b[1])
	}
	hdr := make([]byte, 44)
	copy(hdr[0:4], "RIFF")
	binary.LittleEndian.PutUint32(hdr[4:8], uint32(36+len(data)))
	copy(hdr[8:12], "WAVE")
	copy(hdr[12:16], "fmt ")
	binary.LittleEndian.PutUint32(hdr[16:20], 16)     // PCM chunk size
	binary.LittleEndian.PutUint16(hdr[20:22], 1)      // PCM
	binary.LittleEndian.PutUint16(hdr[22:24], 1)      // mono
	binary.LittleEndian.PutUint32(hdr[24:28], rate)   // sample rate
	binary.LittleEndian.PutUint32(hdr[28:32], rate*2) // byte rate
	binary.LittleEndian.PutUint16(hdr[32:34], 2)      // block align
	binary.LittleEndian.PutUint16(hdr[34:36], 16)     // bits
	copy(hdr[36:40], "data")
	binary.LittleEndian.PutUint32(hdr[40:44], uint32(len(data)))
	return os.WriteFile(path, append(hdr, data...), 0o644)
}

// ---------------------------------------------------------------------------
// main
// ---------------------------------------------------------------------------

func main() {
	os.Exit(run())
}

func run() int {
	outDir := flag.String("out", "", "directory to write levelprobe.json / .frames.jsonl / the WAV into (required)")
	deviceName := flag.String("device", "", "[audio] input_device to open; empty = OS default endpoint")
	quietSecs := flag.Float64("quiet", 4, "seconds of the quiet windows")
	speakSecs := flag.Float64("speak", 5, "seconds of the speaking window")
	toneSecs := flag.Float64("tone", 4, "seconds of the played-tone window")
	flag.Parse()

	if *outDir == "" {
		fmt.Fprintln(os.Stderr, "levelprobe: -out is required")
		return 2
	}
	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "levelprobe: out dir: %v\n", err)
		return 2
	}

	sum := map[string]any{
		"probe":      "wisp-levelprobe",
		"started_at": time.Now().Format(time.RFC3339Nano),
		"ruler": map[string]any{
			"new":               "internal/audio.LevelOfSamples (level.go:119) - the DC-stripped ruler ticket 449 landed",
			"old":               "sqrt(mean(s^2))/32768 computed IN THIS PROBE from the same frame - the pre-449 ruler, for the comparison only",
			"fullscale":         audio.LevelFullScale,
			"framesamples":      audio.FrameSamples,
			"frame_duration_ms": audio.FrameDuration.Milliseconds(),
		},
	}

	windows := defaultWindows()
	windows[0].Secs = *quietSecs
	windows[2].Secs = *quietSecs
	windows[1].Secs = *speakSecs
	windows[3].Secs = *toneSecs

	// The tone file is written before the microphone opens so a failure to
	// write it cannot leave an open device behind.
	tonePath := filepath.Join(*outDir, "tone-1k-halffs.wav")
	toneSpec := windows[3]
	if err := writeToneWAV(tonePath, 1000, 16384, toneSpec.Secs); err != nil {
		sum["tone_error"] = err.Error()
	} else {
		sum["tone"] = map[string]any{
			"path": tonePath, "frequency_hz": 1000, "amplitude_counts": 16384,
			"digital_rms_fullscale": 16384.0 / audio.LevelFullScale / math.Sqrt2,
			"note":                  "the played signal's RMS is exact; what the microphone HEARS is empirical (speaker volume, room)",
		}
	}

	framesFile, err := os.OpenFile(filepath.Join(*outDir, "levelprobe.frames.jsonl"),
		os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "levelprobe: frames file (refusing to overwrite): %v\n", err)
		return 2
	}
	defer framesFile.Close()

	mic := audio.NewWASAPIMicrophone(audio.WithInputDevice(*deviceName))
	buf := audio.NewBoundedFrames()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Printf("levelprobe: opening the default capture endpoint (device=%q)\n", *deviceName)
	if err := mic.Start(ctx, buf); err != nil {
		sum["mic_error"] = err.Error()
		sum["windows"] = []windowResult{}
		writeJSON(*outDir, sum)
		fmt.Printf("levelprobe: MICROPHONE DID NOT OPEN: %v\n", err)
		return 2
	}
	defer func() { _ = mic.Stop() }()
	if capDev, renderDev, e := mic.Endpoints(); e == nil {
		sum["endpoints"] = map[string]any{
			"capture": capDev.String(), "render": renderDev.String(),
			"capture_thread": mic.ThreadID(),
		}
	} else {
		sum["endpoints_error"] = e.Error()
	}
	fmt.Printf("levelprobe: device open; ruler FrameDuration=%s\n", audio.FrameDuration)

	var results []windowResult
	for _, spec := range windows {
		acc := &accumulator{spec: spec}
		fmt.Printf("\n=== %s: %s ===\n", spec.Name, spec.Say)
		if spec.Tone {
			if err := playAsync(tonePath); err != nil {
				sum["tone_play_error"] = err.Error()
				fmt.Printf("levelprobe: tone playback failed: %v\n", err)
			}
		}
		span := runWindow(ctx, buf, acc, time.Duration(spec.Secs*float64(time.Second)), framesFile)
		if spec.Tone {
			stopPlay()
		}
		res := acc.result(span.Seconds())
		if res.Samples == 0 {
			res.Note = "ABSENT: this window recorded zero frames - nothing was measured"
		}
		results = append(results, res)
		fmt.Printf("levelprobe: %-8s samples=%d new(mean=%s max=%s) old(mean=%s max=%s) dc(mean=%s)\n",
			res.Name, res.Samples,
			fmtPtr(res.NewLevel.Mean), fmtPtr(res.NewLevel.Max),
			fmtPtr(res.OldLevel.Mean), fmtPtr(res.OldLevel.Max),
			fmtPtr(res.DCCounts.Mean))
	}

	st := mic.Stats()
	sum["windows"] = results
	sum["capture_stats"] = map[string]any{
		"frames_sent": st.FramesSent, "frames_dropped": st.FramesDropped,
		"reopens": st.Reopens, "last_error": st.LastError,
	}
	sum["synthetic"] = synthResults()
	sum["finished_at"] = time.Now().Format(time.RFC3339Nano)
	writeJSON(*outDir, sum)

	missing := 0
	for _, r := range results {
		if r.Samples == 0 {
			missing++
		}
	}
	if missing > 0 {
		fmt.Printf("\nlevelprobe: %d window(s) recorded nothing (see the null entries)\n", missing)
		return 3
	}
	return 0
}

func fmtPtr(p *float64) string {
	if p == nil {
		return "ABSENT"
	}
	return fmt.Sprintf("%.6f", *p)
}

// runWindow consumes frames for the window's duration. A single goroutine does
// the channel, the clock and the playback stop, so the probe never needs a bare
// go func of its own (the capture thread it creates is the package's own
// SpawnCapture, inside internal/audio).
func runWindow(ctx context.Context, buf <-chan []byte, acc *accumulator,
	d time.Duration, out *os.File,
) time.Duration {
	start := time.Now()
	deadline := time.NewTimer(d)
	defer deadline.Stop()
	for {
		select {
		case <-ctx.Done():
			return time.Since(start)
		case <-deadline.C:
			return time.Since(start)
		case frame, ok := <-buf:
			if !ok {
				return time.Since(start)
			}
			samples, err := audio.DecodeFrame(frame)
			if err != nil {
				continue
			}
			var sum, sumSq float64
			for _, s := range samples {
				v := float64(s)
				sum += v
				sumSq += v * v
			}
			rec := frameRecord{
				Phase: acc.spec.Name,
				TMS:   time.Since(start).Milliseconds(),
				New:   audio.LevelOfSamples(samples),
				Old:   math.Sqrt(sumSq/float64(len(samples))) / audio.LevelFullScale,
				DC:    sum / float64(len(samples)),
			}
			acc.add(rec)
			line, _ := json.Marshal(rec)
			_, _ = out.Write(append(line, '\n'))
		}
	}
}

func writeJSON(dir string, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "levelprobe: marshal: %v\n", err)
		return
	}
	b = append(b, '\n')
	p := filepath.Join(dir, "levelprobe.json")
	if _, statErr := os.Stat(p); statErr == nil {
		fmt.Fprintf(os.Stderr, "levelprobe: refusing to overwrite %s\n", p)
		return
	}
	if err := os.WriteFile(p, b, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "levelprobe: write %s: %v\n", p, err)
	}
}
