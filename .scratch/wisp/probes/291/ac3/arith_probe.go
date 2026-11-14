// Probe (not production code, not in the module graph: .scratch is skipped by
// go build ./...). Ticket 291 AC#3: measure the arithmetic of the PROPOSED
// minimal change (per-frame mean removal) against the CURRENT one, on every
// fixture the existing nails use, so the landing table carries readings
// instead of opinions.
//
// Baseline = internal/audio/level.go:96-107 as it stands at HEAD (int64 sum of
// squares, no mean removal). Proposed = the same function with a per-frame
// mean subtracted (two passes, float64), i.e. what the landing table proposes.
//
// Run: go run arith_probe.go   (stdlib only, no module needed)
package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

const (
	frameSamples   = 512
	frameBytes     = frameSamples * 2
	levelFullScale = 32768.0
	levelLSB       = 1.0 / levelFullScale
	sineTolerance  = 5e-5
	targetRate     = 16000 // internal/audio TargetRate
)

// baseline is internal/audio/level.go LevelOfSamples verbatim.
func baseline(samples []int16) float64 {
	if len(samples) == 0 {
		return 0
	}
	var sumSquares int64
	for _, s := range samples {
		v := int64(s)
		sumSquares += v * v
	}
	rms := math.Sqrt(float64(sumSquares) / float64(len(samples)))
	return rms / levelFullScale
}

// proposed is the minimal change: subtract this window's own mean, then RMS.
// Two passes, float64, no allocation, no state (same signature and contract).
func proposed(samples []int16) float64 {
	n := len(samples)
	if n == 0 {
		return 0
	}
	var sum int64
	for _, s := range samples {
		sum += int64(s)
	}
	mean := float64(sum) / float64(n)
	var sumSquares float64
	for _, s := range samples {
		d := float64(s) - mean
		sumSquares += d * d
	}
	rms := math.Sqrt(sumSquares / float64(n))
	return rms / levelFullScale
}

// exactIntVariant is the integer-exact variance identity
//
//	var = (n*sum(s^2) - sum(s)^2) / n^2
//
// kept here only to measure its overflow boundary (it is NOT the recommendation).
func exactIntVariant(samples []int16) (float64, bool) {
	n := int64(len(samples))
	if n == 0 {
		return 0, true
	}
	var sum, sumSq int64
	for _, s := range samples {
		sum += int64(s)
		sumSq += int64(s) * int64(s)
	}
	// n*sumSq overflows int64 when n^2*32768^2 > 2^63.
	if n > 0 && sumSq > math.MaxInt64/n {
		return 0, false
	}
	num := n*sumSq - sum*sum
	if num < 0 {
		num = 0
	}
	return math.Sqrt(float64(num)) / float64(n) / levelFullScale, true
}

func squareFrame(amp int16) []int16 {
	s := make([]int16, frameSamples)
	for i := range s {
		if i%2 == 0 {
			s[i] = amp
		} else {
			s[i] = -amp
		}
	}
	return s
}

func wholeCycleSine(n int, amp int64, cycles int) []int16 {
	out := make([]int16, n)
	for i := range out {
		out[i] = int16(math.Round(float64(amp) * math.Sin(2*math.Pi*float64(cycles)*float64(i)/float64(n))))
	}
	return out
}

func constFrame(v int16) []int16 {
	s := make([]int16, frameSamples)
	for i := range s {
		s[i] = v
	}
	return s
}

func meanOf(s []int16) float64 {
	var sum int64
	for _, v := range s {
		sum += int64(v)
	}
	return float64(sum) / float64(len(s))
}

func show(name string, s []int16) {
	b := baseline(s)
	p := proposed(s)
	fmt.Printf("%-46s mean=%12.4f  baseline=%.9f  proposed=%.9f  delta=%+.3e\n",
		name, meanOf(s), b, p, p-b)
}

func main() {
	fmt.Println("== fixture readings: baseline (no DC removal) vs proposed (per-frame mean removal) ==")
	show("all-zero frame", constFrame(0))
	show("constant -32768 frame", constFrame(-32768))
	show("constant +32767 frame", constFrame(32767))
	show("constant 12345 frame (capturelevel fixture)", constFrame(12345))
	show("square +-32767", squareFrame(32767))
	show("square +-16384 (half scale)", squareFrame(16384))
	show("square +-8192 (quarter)", squareFrame(8192))
	show("square +-1 (LevelLSB)", squareFrame(1))
	show("hottest legal bipolar square +32767/-32768", func() []int16 {
		s := make([]int16, frameSamples)
		for i := range s {
			if i%2 == 0 {
				s[i] = 32767
			} else {
				s[i] = -32768
			}
		}
		return s
	}())
	show("whole-cycle sine amp 32767 (8 cycles)", wholeCycleSine(frameSamples, 32767, 8))
	show("whole-cycle sine amp 16384", wholeCycleSine(frameSamples, 16384, 8))
	show("whole-cycle sine amp 8192", wholeCycleSine(frameSamples, 8192, 8))
	// ticket 291 mic probe shape: DC 12345 + whole-cycle sine 6000
	dcSin := make([]int16, frameSamples)
	sinPart := wholeCycleSine(frameSamples, 6000, 8)
	for i := range dcSin {
		dcSin[i] = 12345 + sinPart[i]
	}
	show("DC 12345 + sine 6000 (ticket 291 shape)", dcSin)

	fmt.Println("\n== amplitude-pair fixture (internal/audio/amplitudepair_297_windows_test.go) ==")
	lowCounts := int16(math.Round(0.05 * 32767))
	highCounts := int16(math.Round(0.4 * 32767))
	lowWin := constFrame(lowCounts)
	highWin := constFrame(highCounts)
	fmt.Printf("lowCounts=%d highCounts=%d amplitude ratio=%.6f\n", lowCounts, highCounts,
		float64(highCounts)/float64(lowCounts))
	fmt.Printf("window A: baseline=%.9f proposed=%.9f (want |%d|/FS=%.9f)\n",
		baseline(lowWin), proposed(lowWin), lowCounts, math.Abs(float64(lowCounts))/levelFullScale)
	fmt.Printf("window B: baseline=%.9f proposed=%.9f\n", baseline(highWin), proposed(highWin))
	fmt.Printf("mean ratio: baseline=%.6f proposed=%.6f (0/0 => %v)\n",
		baseline(highWin)/baseline(lowWin), proposed(highWin)/proposed(lowWin),
		math.IsNaN(proposed(highWin)/proposed(lowWin)))

	fmt.Println("\n== capturelevel fixture (internal/audio/capturelevel_windows_test.go:21-24) ==")
	fmt.Printf("constant 12345: baseline=%.9f proposed=%.9f; the test's want is derived with the same function\n",
		baseline(constFrame(12345)), proposed(constFrame(12345)))
	fmt.Printf("so the fixture stays GREEN but witnesses 0 instead of 0.3767, and its comment 'RMS is exactly this magnitude' becomes false\n")

	fmt.Println("\n== convertpacket S2 (float32 bytes sliced as int16, pair-mean) ==")
	{
		// Replicates convertpacket_297_windows_test.go S1/S2/S4 without unsafe.
		const frames = 512
		const rate = 48000.0
		const freq = 375.0
		const amp = 0.05
		pkt := make([]byte, frames*2*4)
		for i := 0; i < frames*2; i++ {
			v := amp * math.Sin(2*math.Pi*freq*float64(i/2)/rate)
			binary.LittleEndian.PutUint32(pkt[4*i:], math.Float32bits(float32(v)))
		}
		raw := make([]int16, frames*2)
		for i := range raw {
			raw[i] = int16(binary.LittleEndian.Uint16(pkt[2*i:]))
		}
		wrong := make([]int16, frames)
		for i := range wrong {
			v := (int32(raw[2*i]) + int32(raw[2*i+1])) / 2
			if v > 32767 {
				v = 32767
			}
			if v < -32768 {
				v = -32768
			}
			wrong[i] = int16(v)
		}
		fmt.Printf("S2 wrong-byte window: mean=%.3f counts  baseline=%.6f  proposed=%.6f (assertion is > 0.1)\n",
			meanOf(wrong), baseline(wrong), proposed(wrong))
		// S1: float path mono = the float samples scaled to int16 by convertPacket
		mono := make([]int16, frames)
		for i := 0; i < frames; i++ {
			mono[i] = int16(math.Round(amp * math.Sin(2*math.Pi*freq*float64(i)/rate) * 32768))
		}
		fmt.Printf("S1 float-path window: baseline=%.9f proposed=%.9f (want amp/sqrt2=%.9f)\n",
			baseline(mono), proposed(mono), amp/math.Sqrt2)
		// S4: +-1638 square in pairs of two samples
		s4 := make([]int16, frames)
		for i := 0; i < frames; i++ {
			if (i/2)%2 == 1 {
				s4[i] = -1638
			} else {
				s4[i] = 1638
			}
		}
		fmt.Printf("S4 pair-square window: mean=%.3f baseline=%.9f proposed=%.9f (want 1638/FS=%.9f)\n",
			meanOf(s4), baseline(s4), proposed(s4), 1638.0/levelFullScale)
	}

	fmt.Println("\n== sine tolerance after the change (level_test.go TestLevelSineAgainstRootTwo) ==")
	for _, amp := range []int64{32767, 16384, 8192} {
		s := wholeCycleSine(frameSamples, amp, 8)
		ideal := float64(amp) / levelFullScale / math.Sqrt2
		b := baseline(s)
		p := proposed(s)
		fmt.Printf("amp %5d: ideal=%.9f baseline err=%+.3e proposed err=%+.3e (tolerance %.1e) %s\n",
			amp, ideal, b-ideal, p-ideal, sineTolerance,
			map[bool]string{true: "PASS", false: "RED"}[math.Abs(p-ideal) <= sineTolerance])
	}

	fmt.Println("\n== 'exactly zero' / bit checks for a constant frame ==")
	for _, v := range []int16{0, 1, -1, 12345, 32767, -32768} {
		f := constFrame(v)
		p := proposed(f)
		pf := float32(p)
		fmt.Printf("constant %6d: proposed float64 bits=%016x value=%v ; float32 bits=%08x value=%v\n",
			v, math.Float64bits(p), p, math.Float32bits(pf), pf)
	}

	fmt.Println("\n== translation invariance (the judgeable 'DC is gone' fact) ==")
	{
		base := wholeCycleSine(frameSamples, 3000, 5)
		lv0 := proposed(base)
		ok := true
		for _, c := range []int16{-12345, -3000, -1, 1, 3000, 12345, 20000} {
			shifted := make([]int16, len(base))
			for i, v := range base {
				shifted[i] = v + c // may clip at int16 bounds; only report in-range shifts
			}
			lv := proposed(shifted)
			same := math.Float64bits(lv) == math.Float64bits(lv0)
			if !same {
				ok = false
			}
			fmt.Printf("shift %+6d: level=%.12f bits-equal-to-unshifted=%v\n", c, lv, same)
		}
		fmt.Printf("translation invariant (bit-exact) over all shifts that did not clip: %v\n", ok)
	}

	fmt.Println("\n== frame-mean removal frequency response (gain in dB, best/worst phase, 512 samples @16kHz) ==")
	{
		sweep := []float64{5, 10, 15, 20, 23, 25, 30, 31.25, 40, 50, 60, 70, 85, 100, 120, 150, 200, 300, 440, 800, 1000}
		for _, f := range sweep {
			minG, maxG := math.Inf(1), math.Inf(-1)
			for k := 0; k < 128; k++ {
				phase := 2 * math.Pi * float64(k) / 128
				s := make([]int16, frameSamples)
				for i := range s {
					s[i] = int16(math.Round(3000 * math.Sin(2*math.Pi*f*float64(i)/targetRate+phase)))
				}
				b := baseline(s)
				p := proposed(s)
				var g float64
				if b > 0 {
					g = 20 * math.Log10(p/b)
				}
				if g < minG {
					minG = g
				}
				if g > maxG {
					maxG = g
				}
			}
			fmt.Printf("f=%7.2f Hz  gain min=%+.3f dB  max=%+.3f dB  ripple=%.3f dB\n", f, minG, maxG, maxG-minG)
		}
		// -3 dB corner: first frequency (worst phase) whose gain is <= -3 dB.
		corner := 0.0
		for f := 0.5; f <= 60; f += 0.25 {
			minG := math.Inf(1)
			for k := 0; k < 256; k++ {
				phase := 2 * math.Pi * float64(k) / 256
				s := make([]int16, frameSamples)
				for i := range s {
					s[i] = int16(math.Round(3000 * math.Sin(2*math.Pi*f*float64(i)/targetRate+phase)))
				}
				b := baseline(s)
				p := proposed(s)
				if b <= 0 {
					continue
				}
				if g := 20 * math.Log10(p/b); g < minG {
					minG = g
				}
			}
			if minG <= -3 {
				corner = f
				break
			}
		}
		fmt.Printf("worst-phase -3 dB corner = %.2f Hz\n", corner)
	}

	fmt.Println("\n== endpoint reachability after the change ==")
	{
		hottest := 0.0
		for _, f := range [][]int16{
			squareFrame(32767),
			func() []int16 {
				s := make([]int16, frameSamples)
				for i := range s {
					if i%2 == 0 {
						s[i] = 32767
					} else {
						s[i] = -32768
					}
				}
				return s
			}(),
			func() []int16 {
				s := make([]int16, frameSamples)
				for i := range s {
					if i < frameSamples/2 {
						s[i] = 32767
					} else {
						s[i] = -32768
					}
				}
				return s
			}(),
		} {
			if p := proposed(f); p > hottest {
				hottest = p
			}
		}
		fmt.Printf("max proposed level found = %.12f  (MaxLevel=1.0 attained? %v)\n", hottest, hottest == 1.0)
		fmt.Printf("named bipolar-square bound 32767.5/32768 = %.12f\n", 32767.5/levelFullScale)
		fmt.Printf("closed-scale check: 0 <= max <= 1 -> %v\n", hottest >= 0 && hottest <= 1)
	}

	fmt.Println("\n== exact-integer variant: int64 overflow boundary ==")
	{
		// Largest n with n*sumSq <= MaxInt64 when every sample is +-32768:
		// n^2 * 32768^2 <= 2^63  =>  n <= sqrt(2^63)/32768 = 92681.9.
		// The check uses float64 on purpose: the int64 product WRAPS past the
		// boundary (that wrap is the finding, not a bug in the probe).
		const maxI64f = 9223372036854775808.0 // 2^63
		nMax := int(math.Floor(math.Sqrt(maxI64f) / levelFullScale))
		for float64(nMax+1)*float64(nMax+1)*levelFullScale*levelFullScale <= maxI64f {
			nMax++
		}
		for float64(nMax)*float64(nMax)*levelFullScale*levelFullScale > maxI64f {
			nMax--
		}
		fmt.Printf("n*sumSq is safe up to n = %d samples (%.2f s at 16 kHz); LevelOfSamples accepts any n\n",
			nMax, float64(nMax)/targetRate)
		// sanity: same answer as proposed on a DC+AC window
		w, ok := exactIntVariant(dcSin)
		fmt.Printf("exact variant on DC+sine window: %.12f ok=%v (proposed %.12f)\n", w, ok, proposed(dcSin))
	}

	fmt.Println("\n== cost: ns per 512-sample window ==")
	{
		w := wholeCycleSine(frameSamples, 3000, 5)
		const iters = 200000
		start := time.Now()
		var sink float64
		for i := 0; i < iters; i++ {
			sink += baseline(w)
		}
		d1 := time.Since(start)
		start = time.Now()
		for i := 0; i < iters; i++ {
			sink += proposed(w)
		}
		d2 := time.Since(start)
		fmt.Printf("baseline %.1f ns/window ; proposed %.1f ns/window ; x%.2f ; at 31.25 windows/s the added cost is %.3f ms/s of one thread (sink %v)\n",
			float64(d1.Nanoseconds())/iters, float64(d2.Nanoseconds())/iters,
			float64(d2)/float64(d1), (float64(d2-d1)/iters)*31.25/1e6, sink != 0)
	}

	fmt.Println("\n== SilenceLevelGate direction (monotonicity) ==")
	{
		// proposed <= baseline for every frame: mean removal subtracts the mean's
		// own energy, so variance <= E[s^2]. Verified on a spread of fixtures.
		worst := 0.0
		for _, s := range [][]int16{
			constFrame(12345), squareFrame(32767), wholeCycleSine(frameSamples, 32767, 8), dcSin,
			func() []int16 {
				s := make([]int16, frameSamples)
				for i := range s {
					s[i] = int16((i*37)%2000 - 1000 + 12345)
				}
				return s
			}(),
		} {
			if d := proposed(s) - baseline(s); d > worst {
				worst = d
			}
		}
		fmt.Printf("largest increase of proposed over baseline across fixtures = %+.3e (must be <= 0)\n", worst)
	}

	fmt.Println("\n== 297 live heuristic band (printed, not asserted) ==")
	fmt.Printf("livepacket_297_windows_test.go:244 prints the 0.30..0.45 band as a device-floor signature; after the change a DC floor reads ~0 so that branch can never print -- a note, not a red (WISP_LIVE_MIC-gated).\n")

	_ = frameBytes
	_ = levelLSB
}
