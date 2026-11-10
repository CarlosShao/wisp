//go:build windows

package audio

// Ticket 297 AC#0, the package-internal half: an instrument that feeds
// Go-owned byte buffers straight into convertPacket (and parseWaveFormat),
// because the ticket's ruling note (a) established that "convertPacket cannot
// be tested in-repo" was wrong -- it is a package function whose data
// argument is a byte address, so a []byte this test owns is a legal device
// packet substitute. No fake device, no microphone, no COM.
//
// What this instrument is FOR (ticket 297 line 23): print the three faces the
// discriminating rule needs --
//
//	[B] the same bytes read as int16 (what convertPacket does when
//	    floating=false),
//	[C] the same bytes read as float32 (what convertPacket does when
//	    floating=true),
//	[D] the level of the converted window,
//
// and pin the exact arithmetic on synthetic bytes so the real-device run
// (livepacket_297_windows_test.go) has a computed counterpart to compare
// against. The discriminating rule from the ticket, restated on synthetic
// signal: a float32 packet sliced as int16 shows large/small alternating
// values with |value| in the low tens of thousands, and its RMS lands near
// 0.37 -- the exact floor ticket 247 measured on the real device. This file
// does NOT decide the disposition (that is AC#2, the orchestrator's call); it
// only makes the shape computable in-repo.
//
// Zero production semantics touched: convertPacket, parseWaveFormat,
// MonoDownmix, FloatToPCM16 are called, never changed.

import (
	"encoding/binary"
	"math"
	"runtime"
	"testing"
	"unsafe"
)

// sinePacketF32Stereo builds one interleaved stereo float32 device packet:
// both channels carry the same sine, frames frames of it. Deterministic:
// math.Sin is a pure function, so the bytes are identical on every run.
func sinePacketF32Stereo(t *testing.T, frames int, amp float64) []byte {
	t.Helper()
	const (
		rate = 48000 // nominal; convertPacket is rate-blind
		freq = 375.0 // 48000/512*4: exactly 4 whole cycles per 512-sample window, so mean(sin^2) is exactly 1/2 and the RMS check cannot hide a partial-cycle bias
	)
	pkt := make([]byte, frames*2*4)
	for i := 0; i < frames*2; i++ {
		v := amp * math.Sin(2*math.Pi*freq*float64(i/2)/rate)
		binary.LittleEndian.PutUint32(pkt[4*i:], math.Float32bits(float32(v)))
	}
	return pkt
}

// int16View reads the packet as convertPacket's int16 path would.
func int16View(pkt []byte, n int) []int16 {
	out := make([]int16, n)
	for i := 0; i < n; i++ {
		out[i] = int16(binary.LittleEndian.Uint16(pkt[2*i:]))
	}
	return out
}

// float32View reads the packet as convertPacket's float path would.
func float32View(pkt []byte, n int) []float32 {
	out := make([]float32, n)
	for i := 0; i < n; i++ {
		out[i] = math.Float32frombits(binary.LittleEndian.Uint32(pkt[4*i:]))
	}
	return out
}

// TestAC297ConvertPacketSyntheticByteFace is the synthetic byte-face
// instrument. Faces:
//
//	S1 float32 stereo packet through the CORRECT path (floating=true):
//	   mono RMS must read amp/sqrt(2) within SineLevelTolerance.
//	S2 the SAME bytes through the WRONG path (floating=false, i.e. the
//	   shape a tag resolved to 0 would produce): output must equal the
//	   plain int16 pair-mean of the byte stream, computed here with
//	   encoding/binary only -- and its level is printed next to the 0.37
//	   floor from ticket 247 for the readings file.
//	S3 silent flag: the packet must come back all zeros at the same
//	   length (the wasapi_windows.go silent contract).
//	S4 int16 stereo packet through the int16 path: pair-mean downmix,
//	   level reads the DC amplitude exactly.
//	S5 [A] face: parseWaveFormat resolves the extensible quad -- prints
//	   the tag/channels/rate/bits tuple for both subtypes (the quad the
//	   live run prints for the real mix format).
func TestAC297ConvertPacketSyntheticByteFace(t *testing.T) {
	const (
		frames = 512 // one seam window of converted mono
		amp    = 0.05
	)

	// --- S1: correct float path -------------------------------------------
	pkt := sinePacketF32Stereo(t, frames, amp)
	f32 := float32View(pkt, frames*2)
	t.Logf("[C] same bytes as float32, first 8: %v", f32[:8])
	t.Logf("[B] same bytes as int16, first 16: %v", int16View(pkt, 16))

	p := unsafe.Pointer(&pkt[0])
	mono := convertPacket(p, frames, 2, true, false)
	runtime.KeepAlive(pkt)
	if len(mono) != frames {
		t.Fatalf("S1 float path: mono length = %d, want %d", len(mono), frames)
	}
	wantLevel := amp / math.Sqrt2
	gotLevel := LevelOfSamples(mono)
	if diff := math.Abs(gotLevel - wantLevel); diff > SineLevelTolerance {
		t.Fatalf("S1 float path: level = %.8f, want %.8f (diff %.2e > SineLevelTolerance)",
			gotLevel, wantLevel, diff)
	}
	t.Logf("[D] float path level = %.6f (amp %v / sqrt(2) = %.6f, within SineLevelTolerance)",
		gotLevel, amp, wantLevel)

	// --- S2: the same bytes through the int16 path (wrong-byte shape) -----
	wrong := convertPacket(p, frames, 2, false, false)
	runtime.KeepAlive(pkt)
	if len(wrong) != frames {
		t.Fatalf("S2 int16-on-float path: mono length = %d, want %d", len(wrong), frames)
	}
	// Expectation computed WITHOUT unsafe: read the byte buffer as int16 and
	// pair-average exactly like MonoDownmix (truncating int32 division).
	raw := int16View(pkt, frames*2)
	want := make([]int16, frames)
	for i := range want {
		want[i] = clampI16((int32(raw[2*i]) + int32(raw[2*i+1])) / 2)
	}
	for i := range want {
		if wrong[i] != want[i] {
			t.Fatalf("S2 int16-on-float path: sample %d = %d, want %d (byte-level downmix)",
				i, wrong[i], want[i])
		}
	}
	wrongLevel := LevelOfSamples(wrong)
	t.Logf("[D2] SAME float bytes sliced as int16 (tag=0 shape): level = %.6f (RMS %.0f counts) -- ticket 247 floor was 0.367691 (12050 counts)",
		wrongLevel, wrongLevel*32768)
	if wrongLevel < 0.1 {
		t.Fatalf("S2 sanity: int16 slicing of float bits should land far from the float level; got %.6f", wrongLevel)
	}

	// --- S3: silent flag ---------------------------------------------------
	silent := convertPacket(p, frames, 2, true, true)
	runtime.KeepAlive(pkt)
	if len(silent) != frames {
		t.Fatalf("S3 silent: length = %d, want %d", len(silent), frames)
	}
	for i, v := range silent {
		if v != 0 {
			t.Fatalf("S3 silent: sample %d = %d, want 0", i, v)
		}
	}
	t.Logf("[S3] silent-flagged packet: %d samples, all zero", frames)

	// --- S4: int16 stereo packet ------------------------------------------
	const ampCounts = 1638 // ~0.05 on the raw scale
	ipkt := make([]byte, frames*2*2)
	for i := 0; i < frames*2; i++ {
		v := int16(ampCounts)
		if (i/2)%2 == 1 {
			v = -int16(ampCounts) // square wave: RMS is exactly the amplitude
		}
		binary.LittleEndian.PutUint16(ipkt[2*i:], uint16(v))
	}
	ip := unsafe.Pointer(&ipkt[0])
	mono16 := convertPacket(ip, frames, 2, false, false)
	runtime.KeepAlive(ipkt)
	if len(mono16) != frames {
		t.Fatalf("S4 int16 path: mono length = %d, want %d", len(mono16), frames)
	}
	// Both channels carry the same square wave, so the pair-mean returns the
	// exact sample; RMS of +-ampCounts is exactly ampCounts.
	got16 := LevelOfSamples(mono16)
	want16 := float64(ampCounts) / LevelFullScale
	if diff := math.Abs(got16 - want16); diff > SineLevelTolerance {
		t.Fatalf("S4 int16 path: level = %.8f, want %.8f (diff %.2e > SineLevelTolerance)",
			got16, want16, diff)
	}
	t.Logf("[S4] int16 stereo square wave: level = %.8f (exact amplitude/FullScale)", got16)

	// --- S5: the [A] quad from parseWaveFormat ----------------------------
	for _, f := range []face300{
		{channels: 2, rate: 48000, bits: 32, data1: 3, data1At: 24}, // IEEE float
		{channels: 2, rate: 48000, bits: 16, data1: 1, data1At: 24}, // PCM
	} {
		face := f.build()
		wf, err := parseWaveFormat(unsafe.Pointer(&face[0]))
		runtime.KeepAlive(face)
		if err != nil {
			t.Fatalf("S5 parseWaveFormat error = %v", err)
		}
		t.Logf("[A] extensible face (bits=%d) resolved quad: tag=%d channels=%d rate=%d bits=%d",
			f.bits, wf.tag, wf.channels, wf.rate, wf.bits)
		if wf.tag != uint16(f.data1) || wf.channels != f.channels || wf.rate != f.rate || wf.bits != f.bits {
			t.Fatalf("S5 quad = {%d %d %d %d}, want {%d %d %d %d}",
				wf.tag, wf.channels, wf.rate, wf.bits, f.data1, f.channels, f.rate, f.bits)
		}
	}
}
