package speech

import (
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/CarlosShao/wisp/internal/audio"
)

// Speech surrogate for the endpointing tests. These tests must not need a
// corpus (ticket 15b-r1 refused to block AC#3 on recorded speech), so they
// synthesise one. This shape was chosen by measurement, not by taste: a
// harmonic stack at 120Hz amplitude-modulated at 3.5Hz (syllable rate) is the
// one synthetic signal out of five candidates that silero classifies as
// speech AND keeps classifying as speech for 65s straight, which the 60s
// guard needs. A 440Hz sine, white noise and noise-AM are all rejected by the
// model, and a 700Hz resonator surrogate flickers over long runs.
const (
	surrogateFundamental = 120.0
	surrogateHarmonics   = 8
	surrogateSyllableHz  = 3.5
	surrogateAmplitude   = 0.2
)

// synthSpeech returns n samples of the surrogate at the seam rate.
func synthSpeech(n int) []int16 {
	out := make([]int16, n)
	for i := range out {
		tt := float64(i) / float64(audio.TargetRate)
		env := 0.5 + 0.5*math.Sin(2*math.Pi*surrogateSyllableHz*tt)
		var acc float64
		for h := 1; h <= surrogateHarmonics; h++ {
			acc += math.Sin(2*math.Pi*surrogateFundamental*float64(h)*tt) / float64(h)
		}
		out[i] = int16(math.Max(-1, math.Min(1, surrogateAmplitude*env*acc)) * 32767)
	}
	return out
}

// silence returns n zero samples.
func silence(n int) []int16 { return make([]int16, n) }

// msSamples converts milliseconds to samples at the seam rate.
func msSamples(ms int) int { return ms * audio.TargetRate / 1000 }

// writePCM16Wav writes a canonical 44-byte-header 16k/mono PCM16 RIFF file so
// the fixtures travel through the real C8 seam (audio.WavInjector) rather than
// through a test-only shortcut.
func writePCM16Wav(t *testing.T, path string, samples []int16) {
	t.Helper()
	data := make([]byte, 44+len(samples)*2)
	copy(data[0:4], "RIFF")
	binary.LittleEndian.PutUint32(data[4:8], uint32(36+len(samples)*2))
	copy(data[8:12], "WAVE")
	copy(data[12:16], "fmt ")
	binary.LittleEndian.PutUint32(data[16:20], 16)
	binary.LittleEndian.PutUint16(data[20:22], 1) // PCM
	binary.LittleEndian.PutUint16(data[22:24], 1) // mono
	binary.LittleEndian.PutUint32(data[24:28], uint32(audio.TargetRate))
	binary.LittleEndian.PutUint32(data[28:32], uint32(audio.TargetRate*2))
	binary.LittleEndian.PutUint16(data[32:34], 2)
	binary.LittleEndian.PutUint16(data[34:36], 16)
	copy(data[36:40], "data")
	binary.LittleEndian.PutUint32(data[40:44], uint32(len(samples)*2))
	for i, s := range samples {
		binary.LittleEndian.PutUint16(data[44+2*i:], uint16(s))
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write fixture wav: %v", err)
	}
}

// seamFrameCount is how many C8 frames a finite source of n samples produces
// (the last one zero-padded, per the WavInjector contract).
func seamFrameCount(n int) int {
	return (n + audio.FrameSamples - 1) / audio.FrameSamples
}

// sileroModelPath locates the pinned silero model. It honours
// WISP_SILERO_MODEL and otherwise walks up from the working directory, so the
// same test binary works both under "go test" (cwd = the package directory)
// and when run from a staging directory that holds the sherpa DLLs.
func sileroModelPath(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("WISP_SILERO_MODEL"); p != "" {
		return p
	}
	rel := filepath.Join("third_party", "spike-models", "silero_vad.onnx")
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, rel)); err == nil {
			return filepath.Join(dir, rel)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatalf("silero model not found: set WISP_SILERO_MODEL, or run from inside the repo "+
		"(looked for %s; runtime.GOOS=%s)", rel, runtime.GOOS)
	return ""
}
