package audio

// Ticket 306 AC#1/AC#2/AC#2b: the tooth for the three in-line authority values
// inside parseWav (internal/audio/wavinjector.go).
//
// WHY THESE VALUES NEED THEIR OWN RULER. parseWav's fmt-chunk branch reads two
// authority VALUES out of the Windows SDK header family and one authority
// OFFSET, all three written inline as literals instead of referencing the
// package constants:
//
//	wavinjector.go:192  if fmtTag == 0xFFFE {          WAVE_FORMAT_EXTENSIBLE
//	wavinjector.go:196  data[body+24 : body+26]        SubFormat.Data1 start
//	wavinjector.go:218  case fmtTag == 3 && bits == 32 WAVE_FORMAT_IEEE_FLOAT
//
// They cannot reference waveFormatExt / waveFormatFloat / waveFormatPCM because
// those three live in wasapi_windows.go, which starts with "//go:build
// windows", while wavinjector.go is platform-neutral. That is measured, not
// assumed: the orchestrator and the 306-a1 leg each ran the replacement in an
// out-of-repo export tree and got GOOS=linux go build ./internal/audio/ rc=1
// ("undefined: waveFormatExt", "undefined: waveFormatFloat") with rc=0 on the
// windows side of the same tree (readings archived in
// .scratch/wisp/probes/306/orch/logs/r1-orch-verify-306a1-20261010-221240.txt
// section 6 and .scratch/wisp/probes/306/a1/30-table3-candidate-points-and-cost.md).
// So the literals STAY. What was missing was a ruler: before this file the only
// wav fixture encoder in the repo, writeWav in internal/audio/wavinjector_test.go,
// hardcodes the tag and the bit depth (u16le(1) at line 27, u16le(16) at line 32,
// and a fmt-chunk size of 16 at line 26), so every one of its 10 call sites
// drives the PCM16 branch only. Changing any of the three values above to a
// wrong one left the whole package green (both rc=0 attacks are recorded in the
// ticket's "measured" section). This fixture is the missing executor: it feeds
// tag 0xFFFE plus bits 32 with a float SubFormat, so all three sites are on the
// path this test asserts through: :192 rewrites fmtTag via :196, :218 then
// matches and the samples come out of the float32 branch at :224.
//
// WHAT THE ASSERTIONS ARE. Parsed values only: rate, chans, the number of
// decoded samples and each decoded int16 sample. Nothing here asserts that a
// branch "was called" -- there is no counter, no coverage hook and no mock.
//
// THE SAMPLE EXPECTATIONS ARE DERIVED BY HAND, not by calling FloatToPCM16, so
// they are independent of this package's implementation. Each wanted int16 is
// the documented mapping of the float path (resample.go:119-129, pinned
// separately by TestMonoDownmixAndFloatConvert for 0 / 1.0 / -1.0 / 0.25):
//
//	 1.0 -> 32767    (int64(1.0*32767 + 0.5) = 32767)
//	-1.0 -> -32768   (negative: int64(-1.0*32768 - 0.5) = -32768)
//	 0.5 -> 16384    (int64(16383.5 + 0.5) = 16384)
//	-0.5 -> -16384   (int64(-16384 - 0.5) = -16384)
//	 0.25 -> 8192    (int64(8191.75 + 0.5) = 8192)
//	-0.25 -> -8192   (int64(-8192 - 0.5) = -8192)
//	 0.125 -> 4096   (int64(4095.875 + 0.5) = 4096)
//	-0.125 -> -4096  (int64(-4096 - 0.5) = -4096)
//
// The sample COUNT is its own pin: 8 float32 samples are 32 data bytes, and
// reading those bytes as PCM16 would yield 16 samples, so a length equality
// against 8 is what separates "took the float32 branch" from "reinterpreted the
// same bytes as int16".
//
// WHAT THIS FIXTURE DELIBERATELY DOES NOT COVER (ticket 306 is one instrument,
// not a coverage sweep; the injector's other branches are a different family of
// gaps and the ticket bans closing them here):
//   - the "extensible fmt chunk too short" error branch (wavinjector.go:193-195,
//     size < 40) -- this face has size 40, so :194 stays without an executor;
//   - a bare (non-extensible) tag 3 / bits 32 face, which would put an executor
//     on :218 while leaving :192 and :196 toothless;
//   - PCM16 and the unsupported-format default, both already executed by the
//     existing fixtures.
//
// NO DEVICE, NO MICROPHONE, NO REAL .wav FILE IN THE REPO: the bytes are built
// in code (tracked .wav count in this repo is 0 and internal/audio has no
// testdata directory -- 306-a1 table 3 measured both, which is why ticket 306
// picked this shape). The file a test writes goes to t.TempDir(); the bypass
// entry that calls parseWav directly follows hotplug_test.go:581.

import (
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
	"testing"
)

// The expected-side authority values this fixture writes, spelled out here on
// purpose: a fixture that took them from the production code could not reject a
// wrong production value (the same rule ticket 300 AC#6 set for the constants).
// Source: Windows SDK shared/mmreg.h -- WAVE_FORMAT_EXTENSIBLE 0xFFFE (:2376),
// WAVE_FORMAT_IEEE_FLOAT 0x0003 (:2110), KSDATAFORMAT_SUBTYPE_IEEE_FLOAT
// DEFINE_GUIDSTRUCT("00000003-0000-0010-8000-00aa00389b71", ...) (:2480-2484),
// and the WAVEFORMATEXTENSIBLE extension whose cbSize is 22 for a 40-byte fmt
// body. Header line numbers are the installed 10.0.26100.0 tree, quoted from
// the readings in .scratch/wisp/probes/306/ and ticket 300; this test does NOT
// read the header at run time.
const (
	w306TagExtensible  uint16 = 0xFFFE
	w306BitsFloat32    uint16 = 32
	w306CbSize         uint16 = 22 // extension bytes after the 18-byte WAVEFORMATEX
	w306FmtBodySize           = 40 // 18 + cbSize
	w306SubFormatData1 uint32 = 0x00000003
)

// w306SubFormatTail is the rest of KSDATAFORMAT_SUBTYPE_IEEE_FLOAT after Data1:
// Data2 0x0000, Data3 0x0010, Data4 {80 00 00 AA 00 38 9B 71}. It is only
// padding for the offset arithmetic below -- parseWav reads the two low bytes
// of Data1 and nothing else of the GUID -- but writing the real GUID keeps the
// face a byte-for-byte honest WAVEFORMATEXTENSIBLE instead of a made-up one.
var w306SubFormatTail = []byte{0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0xaa, 0x00, 0x38, 0x9b, 0x71}

// extFloatWav306 builds a complete RIFF/WAVE file whose fmt chunk is a 40-byte
// WAVEFORMATEXTENSIBLE with a float32 SubFormat, and whose data chunk holds the
// samples as little-endian float32. Offsets below are relative to the fmt
// chunk body (the first byte after its 8-byte header), which is exactly the
// "body" index parseWav uses at wavinjector.go:188-196.
func extFloatWav306(chans int, rate int, samples []float32) []byte {
	dataLen := len(samples) * 4
	b := make([]byte, 0, 12+8+w306FmtBodySize+8+dataLen)

	// RIFF header. parseWav only checks the two magics and never reads the
	// size field (wavinjector.go:166), but it is written correctly anyway.
	b = append(b, []byte("RIFF")...)
	b = append(b, u32le(uint32(4+8+w306FmtBodySize+8+dataLen))...)
	b = append(b, []byte("WAVE")...)

	// fmt chunk: id + size 40 + the 40-byte WAVEFORMATEXTENSIBLE body.
	b = append(b, []byte("fmt ")...)
	b = append(b, u32le(uint32(w306FmtBodySize))...)
	fm := make([]byte, w306FmtBodySize)
	le := binary.LittleEndian
	le.PutUint16(fm[0:], w306TagExtensible) // +0  wFormatTag = 0xFFFE
	le.PutUint16(fm[2:], uint16(chans))     // +2  nChannels
	le.PutUint32(fm[4:], uint32(rate))      // +4  nSamplesPerSec
	bl := uint16(chans) * (w306BitsFloat32 / 8)
	le.PutUint32(fm[8:], uint32(rate)*uint32(bl)) // +8  nAvgBytesPerSec
	le.PutUint16(fm[12:], bl)                     // +12 nBlockAlign
	le.PutUint16(fm[14:], w306BitsFloat32)        // +14 wBitsPerSample = 32
	le.PutUint16(fm[16:], w306CbSize)             // +16 cbSize = 22
	le.PutUint16(fm[18:], w306BitsFloat32)        // +18 Samples.wValidBitsPerSample
	mask := uint32(1)<<uint32(chans) - 1
	if chans == 1 {
		mask = 0x4 // SPEAKER_FRONT_CENTER. parseWav never reads either this or
		// dwChannelMask at +20; both are here so the face matches the layout.
	}
	le.PutUint32(fm[20:], mask)               // +20 dwChannelMask
	le.PutUint32(fm[24:], w306SubFormatData1) // +24 SubFormat.Data1  <- :196 reads +24..+26
	copy(fm[28:], w306SubFormatTail)          // +28..39 Data2/Data3/Data4
	b = append(b, fm...)

	// data chunk: the samples, little-endian float32.
	b = append(b, []byte("data")...)
	b = append(b, u32le(uint32(dataLen))...)
	for _, f := range samples {
		var tmp [4]byte
		le.PutUint32(tmp[:], math.Float32bits(f))
		b = append(b, tmp[:]...)
	}
	return b
}

// w306Floats / w306WantI16 are the input stream and its hand-derived int16
// image (derivation in the file comment above).
var (
	w306Floats  = []float32{1.0, -1.0, 0.5, -0.5, 0.25, -0.25, 0.125, -0.125}
	w306WantI16 = []int16{32767, -32768, 16384, -16384, 8192, -8192, 4096, -4096}
)

// TestParseWavExtensibleFloat32306 drives parseWav directly with a 40-byte
// extensible float32 face (stereo 48kHz) and asserts the values it resolves.
// Named case for ticket 306 AC#1/AC#2/AC#2b: it must go red, on the spot, for
// each of the three mutations (:192 0xFFFE -> 0xFFFD, :196 body+24 -> body+26,
// :218 3 -> 2).
func TestParseWavExtensibleFloat32306(t *testing.T) {
	raw := extFloatWav306(2, 48000, w306Floats)
	if len(raw) != 12+8+40+8+len(w306Floats)*4 {
		t.Fatalf("fixture length = %d, want %d", len(raw), 12+8+40+8+len(w306Floats)*4)
	}

	samples, rate, chans, err := parseWav(raw)
	if err != nil {
		t.Fatalf("parseWav error = %v, want nil: an extensible fmt tag with a float SubFormat "+
			"has to resolve to the float32 path (tag read at wavinjector.go:192 and :196, "+
			"branch chosen at :218); rate=%d chans=%d", err, rate, chans)
	}
	if rate != 48000 {
		t.Errorf("rate = %d, want 48000 (read at fmt body +4)", rate)
	}
	if chans != 2 {
		t.Errorf("chans = %d, want 2 (read at fmt body +2)", chans)
	}
	if len(samples) != len(w306WantI16) {
		t.Fatalf("len(samples) = %d, want %d: 8 float32 samples are 32 data bytes, and the only "+
			"other reading of those bytes is 16 PCM16 samples", len(samples), len(w306WantI16))
	}
	for i, want := range w306WantI16 {
		if samples[i] != want {
			t.Errorf("sample %d = %d, want %d (float32 %v converted through the format-3 branch)",
				i, samples[i], want, w306Floats[i])
		}
	}
}

// TestWavInjectorExtensibleFloat32306 is the same face through the real C8
// entry point: NewWavInjector over a file in t.TempDir(). It is 16kHz mono so
// construction applies neither MonoDownmix nor ResampleLinear
// (wavinjector.go:65-68) and the seam sample stream is exactly the parsed
// stream, which keeps this assertion about the parse too.
func TestWavInjectorExtensibleFloat32306(t *testing.T) {
	raw := extFloatWav306(1, TargetRate, w306Floats)
	p := filepath.Join(t.TempDir(), "extensible-float32-306.wav")
	if err := os.WriteFile(p, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	inj, err := NewWavInjector(p)
	if err != nil {
		t.Fatalf("NewWavInjector error = %v, want nil: the injector must accept a 40-byte "+
			"extensible fmt chunk with a float32 SubFormat, not reject it as an unsupported "+
			"wav format", err)
	}
	if len(inj.samples) != len(w306WantI16) {
		t.Fatalf("injector holds %d samples, want %d", len(inj.samples), len(w306WantI16))
	}
	for i, want := range w306WantI16 {
		if inj.samples[i] != want {
			t.Errorf("seam sample %d = %d, want %d", i, inj.samples[i], want)
		}
	}
}
