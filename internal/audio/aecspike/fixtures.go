package aecspike

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"
)

// Synthetic false-trigger fixtures and minimal 16-bit PCM wav IO. The
// fixtures are deterministic pure math (no PRNG): the same run re-reads the
// same waveforms, which is what makes the false-trigger table comparable
// across volumes and across machines.

// wavHeaderSize is the canonical 44-byte PCM wav header.
const wavHeaderSize = 44

// WriteWav16 writes mono 16-bit PCM at rate. Called only under the
// WISP_AEC_SPIKE gate, only into the workspace temp root (doc.go).
func WriteWav16(path string, rate int, samples []int16) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	dataLen := len(samples) * 2
	h := make([]byte, wavHeaderSize)
	copy(h[0:], "RIFF")
	binary.LittleEndian.PutUint32(h[4:], uint32(wavHeaderSize-8+dataLen))
	copy(h[8:], "WAVE")
	copy(h[12:], "fmt ")
	binary.LittleEndian.PutUint32(h[16:], 16)
	binary.LittleEndian.PutUint16(h[20:], 1)
	binary.LittleEndian.PutUint16(h[22:], 1)
	binary.LittleEndian.PutUint32(h[24:], uint32(rate))
	binary.LittleEndian.PutUint32(h[28:], uint32(rate*2))
	binary.LittleEndian.PutUint16(h[32:], 2)
	binary.LittleEndian.PutUint16(h[34:], 16)
	copy(h[36:], "data")
	binary.LittleEndian.PutUint32(h[40:], uint32(dataLen))
	body := make([]byte, dataLen)
	for i, s := range samples {
		binary.LittleEndian.PutUint16(body[i*2:], uint16(s))
	}
	return os.WriteFile(path, append(h, body...), 0o644)
}

// ReadWav16 reads a 16-bit PCM wav (mono or stereo, first channel kept).
func ReadWav16(path string) (rate int, samples []int16, err error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return 0, nil, err
	}
	if len(raw) < wavHeaderSize || string(raw[0:4]) != "RIFF" || string(raw[8:12]) != "WAVE" {
		return 0, nil, fmt.Errorf("aecspike: %s is not a RIFF/WAVE file", path)
	}
	if binary.LittleEndian.Uint16(raw[20:]) != 1 || binary.LittleEndian.Uint16(raw[34:]) != 16 {
		return 0, nil, fmt.Errorf("aecspike: %s is not 16-bit PCM", path)
	}
	rate = int(binary.LittleEndian.Uint32(raw[24:]))
	channels := int(binary.LittleEndian.Uint16(raw[22:]))
	body := raw[wavHeaderSize:]
	n := len(body) / 2 / channels
	samples = make([]int16, n)
	for i := 0; i < n; i++ {
		samples[i] = int16(binary.LittleEndian.Uint16(body[i*2*channels:]))
	}
	return rate, samples, nil
}

// SynthMusic is a deterministic music-like fixture: a cycling four-chord
// pad (harmonic stacks with a 2Hz tremolo) over a 120BPM low-sine kick.
// It is NOT real music; the table rows built from it carry that label.
func SynthMusic(rate, seconds int) []float32 {
	n := rate * seconds
	out := make([]float32, n)
	chords := [4][3]float64{
		{261.63, 329.63, 392.00}, // C
		{196.00, 246.94, 392.00}, // G
		{220.00, 261.63, 329.63}, // Am
		{174.61, 220.00, 261.63}, // F
	}
	for i := 0; i < n; i++ {
		t := float64(i) / float64(rate)
		chord := chords[int(t/0.5)%len(chords)]
		v := 0.0
		for k, f := range chord {
			v += math.Sin(2*math.Pi*f*t) / float64(k+2)
			v += math.Sin(2*math.Pi*2*f*t) / float64(4*(k+2))
		}
		v *= 0.55 + 0.45*math.Sin(2*math.Pi*2*t) // tremolo
		beat := math.Mod(t, 0.5)
		if beat < 0.06 {
			v += 0.8 * math.Sin(2*math.Pi*55*t) * math.Exp(-beat*60)
		}
		out[i] = clamp(float32(v * 0.4))
	}
	return out
}

// SynthVideo mixes the music bed with deterministic wideband bursts (the
// "effects/action" signature of video playback) and, when speech is given,
// a ducked speech excerpt - the same composition the ticket's "user watches
// a video" row needs.
func SynthVideo(rate, seconds int, speech []float32) []float32 {
	music := SynthMusic(rate, seconds)
	out := make([]float32, len(music))
	var speechOn bool
	for i := range out {
		t := float64(i) / float64(rate)
		burst := 0.0
		if frac := math.Mod(t, 1.5); frac < 0.15 {
			burst = 0.5 * math.Sin(2*math.Pi*1800*t) * (1 - frac/0.15)
		}
		v := float64(music[i])
		if speechOn {
			k := i % len(speech)
			v += 0.7 * float64(speech[k])
		}
		if int(t)%4 == 3 {
			speechOn = true // a 1-in-4 speech excerpt window, deterministic
		} else {
			speechOn = false
		}
		out[i] = clamp(float32(v) + float32(burst)*0.3)
	}
	return out
}

// SynthSpeechLike is the fallback speech fixture when the pinned matcha TTS
// is unavailable: syllable-rate modulated formant tones. Fixture rows built
// from it are labelled as such in the table; the primary rows use real TTS.
func SynthSpeechLike(rate, seconds int) []float32 {
	n := rate * seconds
	out := make([]float32, n)
	formants := [3]float64{500, 1500, 2500}
	for i := 0; i < n; i++ {
		t := float64(i) / float64(rate)
		syl := math.Mod(t, 0.25)
		if syl > 0.16 {
			out[i] = 0
			continue
		}
		f0 := 130 + 40*math.Sin(2*math.Pi*0.7*t)
		v := 0.0
		for k, f := range formants {
			v += math.Sin(2*math.Pi*f*t) / float64(k+1)
		}
		out[i] = clamp(float32(v * 0.3 * math.Exp(-syl*12) * (0.8 + 0.2*math.Sin(2*math.Pi*f0*t))))
	}
	return out
}

func clamp(v float32) float32 {
	if v > 1 {
		return 1
	}
	if v < -1 {
		return -1
	}
	return v
}

// FloatToPCM16 converts normalized float32 samples to int16 PCM.
func FloatToPCM16(in []float32) []int16 {
	out := make([]int16, len(in))
	for i, v := range in {
		x := v * 32767
		if x > 32767 {
			x = 32767
		}
		if x < -32768 {
			x = -32768
		}
		out[i] = int16(x)
	}
	return out
}

// PCM16ToFloat converts int16 PCM to normalized float32 samples.
func PCM16ToFloat(in []int16) []float32 {
	out := make([]float32, len(in))
	for i, v := range in {
		out[i] = float32(v) / 32768
	}
	return out
}

// ResampleLinear is a linear-interpolation resampler for float32 mono. The
// rig and fixtures use it between the TTS rate, the seam rate and the mix
// rate; quality is adequate for spike readings (labelled as such).
func ResampleLinear(in []float32, from, to int) []float32 {
	if from == to || len(in) == 0 {
		return append([]float32(nil), in...)
	}
	ratio := float64(from) / float64(to)
	n := int(float64(len(in)) / ratio)
	out := make([]float32, n)
	for i := range out {
		x := float64(i) * ratio
		lo := int(x)
		if lo >= len(in)-1 {
			out[i] = in[len(in)-1]
			continue
		}
		frac := x - float64(lo)
		out[i] = in[lo] + float32(frac)*(in[lo+1]-in[lo])
	}
	return out
}
