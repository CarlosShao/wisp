//go:build windows

package audio

// Ticket 300 AC#2 tracked test, landed with the offset fix.
//
// This is the rig that 300-a1r ran in an out-of-repo export tree and that
// 300-v1 re-ran and scanned across five offsets (20/22/24/26/28). It is
// committed here because parseWaveFormat now reads SubFormat at the
// authoritative offset: the four faces below assert the values that the
// authoritative layout actually holds there (3 / 1 / 7 / 0). 300-v1's and the
// orchestrator's own scans put this suite green at 24 and red at each of the
// other offsets they swept (20, 22, 26, 28). The "must be red before the fix"
// half of the evidence lives in .scratch/wisp/probes/300/ (a1r, v1, orch),
// never in this tree -- an always-red tracked test is banned here.
//
// Fixture byte layout comes from the AC#0 authority (on-disk Windows SDK
// header C:\Program Files (x86)\Windows Kits\10\Include\10.0.26100.0\shared\mmreg.h,
// compiled with 1-byte packing via pshpack1.h, so there is no padding
// anywhere), NOT from either implementation in this package:
//
//	WAVEFORMATEX      : tag@0 channels@2 rate@4 avgBytes@8 blockAlign@12
//	                  bits@14 cbSize@16                       -> 18 bytes
//	extension (cbSize = 22, mmreg.h:2540 and :2550):
//	                  Samples@18 (union of three WORDs, 2 bytes)
//	                  dwChannelMask@20 (DWORD, 4 bytes)
//	                  SubFormat@24 (GUID, 16 bytes)
//	total 40 bytes; SubFormat.Data1 therefore starts at byte 24, NOT 26.
//
// Nothing here uses a device, a microphone, or a .wav file: parseWaveFormat
// takes an unsafe.Pointer, so the face is built by hand in a []byte. No
// Initialize, no Start.

import (
	"encoding/binary"
	"runtime"
	"testing"
	"unsafe"
)

// face300 describes one 40-byte WAVEFORMATEXTENSIBLE byte face.
type face300 struct {
	channels uint16
	rate     uint32
	bits     uint16
	data1    uint32 // SubFormat.Data1 value to write
	data1At  int    // byte offset the face writes Data1 at (24 = authoritative)
}

// build returns the face. The slice is created with make and returned, so it
// escapes to the heap by construction (a make result that flows out of a
// function cannot be stack-allocated); the caller pins it with
// runtime.KeepAlive across the parseWaveFormat call, and Go's GC does not move
// live objects, so the address handed to unsafe.Pointer stays valid.
func (f face300) build() []byte {
	b := make([]byte, 40)
	le := binary.LittleEndian
	le.PutUint16(b[0:], 0xFFFE)                                       // WAVE_FORMAT_EXTENSIBLE (mmreg.h:2376)
	le.PutUint16(b[2:], f.channels)                                   // nChannels
	le.PutUint32(b[4:], f.rate)                                       // nSamplesPerSec
	le.PutUint32(b[8:], f.rate*uint32(f.channels)*(uint32(f.bits)/8)) // nAvgBytesPerSec
	le.PutUint16(b[12:], f.channels*(f.bits/8))                       // nBlockAlign
	le.PutUint16(b[14:], f.bits)                                      // wBitsPerSample
	le.PutUint16(b[16:], 22)                                          // cbSize = 22 (mmreg.h:2540/2550)
	le.PutUint16(b[18:], f.bits)                                      // Samples.wValidBitsPerSample
	le.PutUint32(b[20:], uint32((1<<f.channels)-1))                   // dwChannelMask
	le.PutUint32(b[f.data1At:], f.data1)                              // SubFormat.Data1
	return b
}

// TestParseWaveFormatSubFormatOffset300 asserts the value parseWaveFormat
// resolves out of an extensible mix format -- not that it was called.
//
// The four faces are chosen so that the assertion cannot be a tautology:
// S1/S2/S3 put SubFormat.Data1 at the authoritative offset 24 and demand the
// value that lives there; S4 is a deliberately malformed face that puts Data1
// at 26 and therefore demands 0, because at offset 24 of that face the
// authoritative layout has the tail of dwChannelMask plus zero bytes. A ruler
// that let both S1 and S4 through would be a ruler that distinguishes nothing.
//
// 300-v1's offset scan is what says where the discriminating power actually
// sits, and the naming here follows it: S2 (Data1 = 1 collides with a
// dwChannelMask of 3 at offset 20) and S3 (7 collides with nothing) carry it;
// S4 is the weakest face in the set -- on its own it stays green at offsets
// 22, 24 and 28 -- and S1 on its own stays green at 20 and 24. No single face
// here is the decisive one; the four together only accept 24.
func TestParseWaveFormatSubFormatOffset300(t *testing.T) {
	cases := []struct {
		name    string
		face    face300
		wantTag uint16
		why     string
	}{
		{
			name:    "S1_float_at_24",
			face:    face300{channels: 2, rate: 48000, bits: 32, data1: 3, data1At: 24},
			wantTag: 3,
			why:     "KSDATAFORMAT_SUBTYPE_IEEE_FLOAT (mmreg.h:2483, Data1=0x00000003) at byte 24",
		},
		{
			name:    "S2_pcm_at_24",
			face:    face300{channels: 2, rate: 48000, bits: 16, data1: 1, data1At: 24},
			wantTag: 1,
			why:     "KSDATAFORMAT_SUBTYPE_PCM (mmreg.h:2474, Data1=0x00000001) at byte 24",
		},
		{
			name:    "S3_bogus_subtype_7_at_24",
			face:    face300{channels: 2, rate: 48000, bits: 32, data1: 7, data1At: 24},
			wantTag: 7,
			why:     "positive control: an obviously wrong Data1 must NOT be swallowed into PCM/FLOAT; the parser has to report the value it actually read, so this face may not go green on a constant",
		},
		{
			name:    "S4_data1_shifted_to_26",
			face:    face300{channels: 2, rate: 48000, bits: 32, data1: 3, data1At: 26},
			wantTag: 0,
			why:     "positive control for offset sensitivity (the weakest face in this set, see the scan note above): malformed face with Data1 two bytes late; per the authoritative layout the value at byte 24 is 0, so an assertion that passes here would read nothing meaningful",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			face := tc.face.build()
			if len(face) != 40 {
				t.Fatalf("face length = %d, want 40", len(face))
			}
			p := unsafe.Pointer(&face[0])
			got, err := parseWaveFormat(p)
			runtime.KeepAlive(face)

			if err != nil {
				t.Fatalf("parseWaveFormat error = %v, want nil; face: %s", err, tc.why)
			}
			if got.tag != tc.wantTag {
				t.Errorf("tag = %d, want %d (%s)", got.tag, tc.wantTag, tc.why)
			}
			// Header sanity: these three offsets are read correctly whatever
			// happens to SubFormat, so a pass line here proves the face is a
			// well-formed WAVEFORMATEXTENSIBLE and that only the SubFormat
			// position is in play.
			if got.channels != tc.face.channels || got.rate != tc.face.rate || got.bits != tc.face.bits {
				t.Errorf("header = {channels %d rate %d bits %d}, want {channels %d rate %d bits %d}",
					got.channels, got.rate, got.bits,
					tc.face.channels, tc.face.rate, tc.face.bits)
			}
			// Real-world consequence named in ticket 300: convertPacket picks
			// the sample path from floating alone, so a wrong tag silently
			// reparses 32-bit float frames as int16.
			gotFloat := got.tag == waveFormatFloat
			wantFloat := tc.wantTag == waveFormatFloat
			if gotFloat != wantFloat {
				t.Errorf("floating = %v, want %v (convertPacket would take the %s path)",
					gotFloat, wantFloat, map[bool]string{true: "float32", false: "int16"}[wantFloat])
			}
		})
	}
}
