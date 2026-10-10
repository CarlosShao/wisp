//go:build windows

package audio

// Ticket 300 AC#6, landing shape C2 (orchestrator ruling, 2026-10-10 17:5x).
//
// What this file is for: AC#6 names a defect in the evidence, not in the
// product -- waveFormatFloat had NO ruler in this repository that could be
// broken by changing it. 300-v3 proved that by mutation (N2: waveFormatFloat
// 3 -> 1 left the whole package green, all 41 top-level cases passing by the
// runtime "--- PASS" ruler). The reason is recorded in the ticket: the only
// existing case that touches the constant,
// parse_wave_format_300_windows_test.go:152-153, puts that same constant on
// BOTH sides of the comparison (got.tag == waveFormatFloat versus
// tc.wantTag == waveFormatFloat), so moving the constant moves both sides and
// the identity keeps holding. A ruler whose expected side reads the thing it is
// supposed to pin pins nothing.
//
// So the shape below is fixed by AC#6's own wording: the EXPECTED side carries
// literals that come from outside this package, and the OBSERVED side is what
// this package actually declares. Nothing on the expected side below reads
// waveFormatFloat, waveFormatExt or waveFormatPCM.
//
// Where the expected literals come from (named as required): they were copied
// out of the Windows SDK header shared/mmreg.h -- an out-of-repo file, not
// anything in this tree -- from the three lines quoted verbatim next to each
// literal below. The readings the orchestrator took are archived in
// .scratch/wisp/probes/300/orch/2026-10-10-mmreg-authority.txt. This case does
// NOT read that header at run time, and it does not t.Skip: the path is
// machine- and build-version-specific, so on a hosted runner the file may
// simply not exist, and an unregistered "--- SKIP" is scored red by
// scripts/portable-tests.sh while registering one would mean editing scripts,
// which is outside AC#4's roster. The authority therefore lives in these copied
// literals plus the citations; the disk reading is evidence about the citation,
// not a precondition of the test.
//
// Terminology, because AC#6 asked that two things not be merged into one
// citation: the WORD values WAVE_FORMAT_PCM, WAVE_FORMAT_IEEE_FLOAT and
// WAVE_FORMAT_EXTENSIBLE are defined in mmreg.h. The GUID
// KSDATAFORMAT_SUBTYPE_IEEE_FLOAT is a different object and lives in ksmedia.h,
// not in mmreg.h; mmreg.h is also what fixed the field offsets back in AC#0.
// Same header family, different things -- AC#0 pinned the OFFSET, this case
// pins the VALUE.
//
// What this case can NOT observe, stated plainly instead of faked: the only
// production consumer of waveFormatFloat is the line
// "floating := s.format.tag == waveFormatFloat" inside (*wasapiStream).Drain,
// so the observed side here never reaches that line. Drain needs a live stream
// from a real device, which ticket 300 bans for a leg (its own line 31), and
// convertPacket takes floating as a bool parameter, so driving convertPacket
// would exercise the parameter rather than the constant. That residual is
// named R2 in the orchestrator's ruling (extracting that line into a package
// predicate would be shape C5, a production-semantics change this ticket does
// not authorize).
//
// No device, no microphone, no .wav file, no Initialize, no Start, no window.

import (
	"encoding/binary"
	"runtime"
	"testing"
	"unsafe"
)

// The three expected-side literals, written in the form the header itself uses
// (0x0003, not a simplified 3 -- simplifying it would recreate the same-source
// problem AC#6 exists to close). Each carries its own citation.
const (
	// mmreg.h:2418, verbatim: "#define WAVE_FORMAT_PCM         1"
	mmregWaveFormatPCM uint16 = 1

	// mmreg.h:2110, verbatim:
	// "#define  WAVE_FORMAT_IEEE_FLOAT                 0x0003 /* Microsoft Corporation */"
	mmregWaveFormatIEEEFloat uint16 = 0x0003

	// mmreg.h:2376, verbatim:
	// "#define  WAVE_FORMAT_EXTENSIBLE                 0xFFFE /* Microsoft */"
	mmregWaveFormatExtensible uint16 = 0xFFFE
)

// r2Face builds a 40-byte WAVEFORMATEXTENSIBLE-shaped face from arguments that
// are all test-carried values. It is a separate helper under a separate name so
// that this file touches no declaration owned by
// parse_wave_format_300_windows_test.go, which is the AC#1/AC#2 evidence body
// and is off-limits to AC#6. A data1At below zero means "do not write
// SubFormat.Data1 at all".
func r2Face(tag uint16, channels uint16, rate uint32, bits uint16, cbSize uint16, data1 uint32, data1At int) []byte {
	b := make([]byte, 40)
	le := binary.LittleEndian
	le.PutUint16(b[0:], tag)
	le.PutUint16(b[2:], channels)
	le.PutUint32(b[4:], rate)
	le.PutUint32(b[8:], rate*uint32(channels)*(uint32(bits)/8))
	le.PutUint16(b[12:], channels*(bits/8))
	le.PutUint16(b[14:], bits)
	le.PutUint16(b[16:], cbSize)
	le.PutUint16(b[18:], bits) // Samples.wValidBitsPerSample
	le.PutUint32(b[20:], uint32((1<<channels)-1))
	if data1At >= 0 {
		le.PutUint32(b[data1At:], data1)
	}
	return b
}

// r2Parse calls parseWaveFormat on a face and keeps the buffer alive across the
// unsafe read, the same way the sibling case does.
func r2Parse(t *testing.T, face []byte) waveFormat {
	t.Helper()
	if len(face) != 40 {
		t.Fatalf("face length = %d, want 40", len(face))
	}
	got, err := parseWaveFormat(unsafe.Pointer(&face[0]))
	runtime.KeepAlive(face)
	if err != nil {
		t.Fatalf("parseWaveFormat error = %v, want nil", err)
	}
	return got
}

// r2AuthorityEntry is one row of the table copied out of mmreg.h. The table is
// keyed by nothing: it is scanned at run time, so mutating a production
// constant changes a value in this package and never turns this file into a
// compile error (a map literal keyed on those constants would).
type r2AuthorityEntry struct {
	value    uint16
	macro    string
	cite     string
	constant string // the name this package declares for that value
}

var r2Authority = []r2AuthorityEntry{
	{mmregWaveFormatPCM, "WAVE_FORMAT_PCM", "mmreg.h:2418", "waveFormatPCM"},
	{mmregWaveFormatIEEEFloat, "WAVE_FORMAT_IEEE_FLOAT", "mmreg.h:2110", "waveFormatFloat"},
	{mmregWaveFormatExtensible, "WAVE_FORMAT_EXTENSIBLE", "mmreg.h:2376", "waveFormatExt"},
}

// TestWaveFormatConstantsMatchMmregAuthority300 is the AC#6 nail.
//
// Every assertion compares a value OBSERVED from this package against a literal
// the test carries from mmreg.h. Changing waveFormatFloat therefore moves the
// observed side only, which is what makes the mutation red: 3 -> 1 fires both
// the float comparison and the one-to-one mapping check (it then collides with
// waveFormatPCM), and 3 -> 0 fires both as well (0 is none of the three
// authority values). The failure text names the constant that moved, not a
// parsed tag, which is the half of AC#6 that the existing identity comparison
// cannot express.
func TestWaveFormatConstantsMatchMmregAuthority300(t *testing.T) {
	cases := []struct {
		name        string
		constName   string
		got         uint16
		want        uint16
		authority   string
		consequence string
	}{
		{
			name:      "ieee_float_low_word",
			constName: "waveFormatFloat",
			got:       waveFormatFloat,
			want:      mmregWaveFormatIEEEFloat,
			authority: "mmreg.h:2110 WAVE_FORMAT_IEEE_FLOAT 0x0003, copied out of the SDK header, not derived from this package",
			consequence: "waveFormatFloat is the constant wasapi_windows.go uses for floating := s.format.tag == " +
				"waveFormatFloat, so its value decides whether 32-bit float frames are sliced as float32 or as " +
				"int16; at 1 it would fire for PCM mix formats and never for float ones, at 0 for nothing at all",
		},
		{
			name:      "extensible_tag",
			constName: "waveFormatExt",
			got:       waveFormatExt,
			want:      mmregWaveFormatExtensible,
			authority: "mmreg.h:2376 WAVE_FORMAT_EXTENSIBLE 0xFFFE, copied out of the SDK header",
			consequence: "waveFormatExt is the tag parseWaveFormat tests before it expands SubFormat; off 0xFFFE every " +
				"extensible mix format keeps 0xFFFE as its tag and never resolves to a subtype at all",
		},
		{
			name:      "pcm_tag",
			constName: "waveFormatPCM",
			got:       waveFormatPCM,
			want:      mmregWaveFormatPCM,
			authority: "mmreg.h:2418 WAVE_FORMAT_PCM 1, copied out of the SDK header",
			consequence: "waveFormatPCM is declared and unused in this package (300-a2 F2, zero uses), so before this " +
				"case no ruler of any kind read it; this cell pins its value and does NOT delete it -- deleting an " +
				"unused declaration is a different cell under a different authorization",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Errorf("%s = %#04x, want %#04x; %s. %s",
					tc.constName, tc.got, tc.want, tc.authority, tc.consequence)
			}
			t.Logf("%s observed %#04x against test-carried %#04x; %s",
				tc.constName, tc.got, tc.want, tc.consequence)
		})
	}

	// Second, independent face for waveFormatFloat: the three constants must map
	// one-to-one onto the three authority literals. Kept as its own subtest so a
	// future edit cannot delete the pin by trimming one comparison.
	t.Run("the_three_constants_map_one_to_one_onto_the_authority_table", func(t *testing.T) {
		observed := []struct {
			name  string
			value uint16
		}{
			{"waveFormatPCM", waveFormatPCM},
			{"waveFormatFloat", waveFormatFloat},
			{"waveFormatExt", waveFormatExt},
		}
		for i := 0; i < len(observed); i++ {
			for j := i + 1; j < len(observed); j++ {
				if observed[i].value == observed[j].value {
					t.Errorf("%s and %s both = %#04x, but mmreg.h gives WAVE_FORMAT_PCM 1, "+
						"WAVE_FORMAT_IEEE_FLOAT 0x0003 and WAVE_FORMAT_EXTENSIBLE 0xFFFE as three distinct values, "+
						"so one of those two constants has moved off the value copied from mmreg.h:2418/2110/2376",
						observed[i].name, observed[j].name, observed[i].value)
				}
			}
		}
		for _, o := range observed {
			matched := false
			for _, a := range r2Authority {
				if o.value != a.value {
					continue
				}
				matched = true
				if o.name != a.constant {
					t.Errorf("%s = %#04x, but mmreg.h reserves %#04x for %s (%s), which this package declares as "+
						"%s; the constant names and the authority values no longer map one to one",
						o.name, o.value, a.value, a.macro, a.cite, a.constant)
				}
			}
			if !matched {
				t.Errorf("%s = %#04x, which is none of WAVE_FORMAT_PCM 1, WAVE_FORMAT_IEEE_FLOAT 0x0003 or "+
					"WAVE_FORMAT_EXTENSIBLE 0xFFFE as copied from mmreg.h:2418/2110/2376", o.name, o.value)
			}
		}
	})

	// Behavioural face for waveFormatExt: the tag gate inside parseWaveFormat is
	// exercised through the parser instead of being read off the declaration, and
	// the expected side is still a test-carried literal. This subtest cannot see
	// waveFormatFloat at all (parseWaveFormat does not use it), so it is named for
	// waveFormatExt and is not to be read as the float nail.
	t.Run("parseWaveFormat_expands_a_face_tagged_with_the_authority_extensible_value", func(t *testing.T) {
		face := r2Face(mmregWaveFormatExtensible, 2, 48000, 32, 22, uint32(mmregWaveFormatIEEEFloat), 24)
		got := r2Parse(t, face)
		if got.tag != mmregWaveFormatIEEEFloat {
			t.Errorf("parseWaveFormat on a face tagged %#04x (mmreg.h:2376 WAVE_FORMAT_EXTENSIBLE) returned tag = "+
				"%#04x, want %#04x (mmreg.h:2110 WAVE_FORMAT_IEEE_FLOAT); waveFormatExt is the constant the parser "+
				"compares that tag against, so waveFormatExt off 0xFFFE switches the extensible branch off and the "+
				"tag field is left holding 0xFFFE",
				mmregWaveFormatExtensible, got.tag, mmregWaveFormatIEEEFloat)
		}
	})

	// Control for the face above: on a plain PCM tag the extensible branch must not
	// fire, so the returned tag stays what the header bytes say. Without this the
	// subtest above could go green for the wrong reason (a parser that resolved
	// SubFormat for every face would pass it and fail here).
	t.Run("control_plain_pcm_tag_is_not_expanded", func(t *testing.T) {
		face := r2Face(mmregWaveFormatPCM, 1, 48000, 16, 0, uint32(mmregWaveFormatIEEEFloat), 24)
		got := r2Parse(t, face)
		if got.tag != mmregWaveFormatPCM {
			t.Errorf("parseWaveFormat on a face tagged %#04x (mmreg.h:2418 WAVE_FORMAT_PCM, cbSize 0) returned tag = "+
				"%#04x, want %#04x; the bytes at offset 24 must be ignored for a non-extensible face",
				mmregWaveFormatPCM, got.tag, mmregWaveFormatPCM)
		}
	})
}
