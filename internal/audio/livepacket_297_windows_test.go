//go:build windows

package audio

// Ticket 297 AC#0, the real-device half: one in-process run against THIS
// machine's default capture endpoint that prints the four things the ticket's
// discriminating rule needs, before any production change is even discussed:
//
//	[A] the parsed mix-format quad (tag/channels/rate/bits), verbatim from
//	    the waveFormat struct wasapi_windows.go fills in;
//	[B] the first real packet's first 16 samples read as int16;
//	[C] the same bytes re-read as float32 (first 8);
//	[D] the FrameLevel / window levels of the converted stream.
//
// The rule (ticket 297 line 23, restated): if [B] shows large/small
// alternating values around the low tens of thousands, the bytes are being
// sliced wrong (float32 packet parsed as int16); if [B] looks like a smooth
// small signal while [C] looks like garbage, the shape is reversed; if both
// are smooth and the RMS really is ~0.37, the device floor (AC#2 form 3) is
// the candidate. This test PRINTS the verdict heuristic; it does not assert
// a disposition -- AC#2 belongs to the orchestrator, not to this leg.
//
// Gated behind WISP_LIVE_MIC=1 exactly like TestLiveWasapiSmoke: CI has no
// microphone, and the ticket requires the one reading it wants to be a real
// run, not a skipped box. It talks to the device through the production
// opener (wasapiOpener) and reads raw packets with the same COM vtable slots
// Drain uses (GetNextPacketSize=5, GetBuffer=3, ReleaseBuffer=4), because
// Drain only hands out converted mono and the ticket asks for the RAW bytes.
// Everything runs on the single test goroutine that opened the stream; no
// goroutine is created here, and convertPacket itself is the production
// converter (same call shape as wasapi_windows.go Drain).

import (
	"math"
	"os"
	"testing"
	"time"
	"unsafe"
)

func TestAC297LivePacketByteFace(t *testing.T) {
	if os.Getenv("WISP_LIVE_MIC") != "1" {
		t.Skip("ticket 297 AC#0 live reading: set WISP_LIVE_MIC=1, keep the room quiet, no other mic consumer (真机那一发)")
	}

	mic := NewWASAPIMicrophone()
	capDev, renderDev, err := mic.Endpoints()
	if err != nil {
		t.Skipf("endpoint enumeration failed: %v", err)
	}
	t.Logf("[A0] capture endpoint id=%q name=%q | render endpoint name=%q",
		capDev.ID, capDev.Name, renderDev.Name)

	// COM on this thread: newMMDeviceWatcher() already initialized MTA here
	// (mmDeviceWatcher leaves it balanced-by-design unbalanced, process
	// resident), so the opener's COM calls below run inside that apartment.
	stream, err := wasapiOpener{}.Open(capDev)
	if err != nil {
		t.Skipf("device open failed (occupied? privacy settings?): %v", err)
	}
	s, ok := stream.(*wasapiStream)
	if !ok {
		t.Fatalf("opener returned %T, want *wasapiStream", stream)
	}
	defer s.Close()

	// [A] the quad the capture stack actually negotiated.
	t.Logf("[A] mix format quad: tag=%d channels=%d rate=%d bits=%d",
		s.format.tag, s.format.channels, s.format.rate, s.format.bits)
	t.Logf("[A2] convertPacket path: floating=%v (tag==waveFormatFloat)",
		s.format.tag == waveFormatFloat)

	ch := int(s.format.channels)
	floating := s.format.tag == waveFormatFloat

	var (
		mono       []int16
		windows    []float64
		printed    bool
		silentPkts int
		first16    []int16
		first8f    []float32
		firstFlags uint32
	)
	deadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) && len(windows) < 64 {
		ok, err := s.WaitEvent(1000)
		if err != nil {
			t.Fatalf("WaitEvent: %v", err)
		}
		if !ok {
			continue
		}
		for {
			var packet uint32
			if hr, _, _ := comCall(s.capture, 5, uintptr(unsafe.Pointer(&packet))); !hrOK(hr) {
				t.Fatalf("GetNextPacketSize hr=%s", hresultString(hr))
			}
			if packet == 0 {
				break
			}
			var (
				data           unsafe.Pointer
				frames, flags  uint32
				devPos, qpcPos uint64
			)
			if hr, _, _ := comCall(s.capture, 3,
				uintptr(unsafe.Pointer(&data)), uintptr(unsafe.Pointer(&frames)),
				uintptr(unsafe.Pointer(&flags)),
				uintptr(unsafe.Pointer(&devPos)), uintptr(unsafe.Pointer(&qpcPos))); !hrOK(hr) {
				t.Fatalf("GetBuffer hr=%s", hresultString(hr))
			}
			if frames > 0 {
				bytesAvail := int(frames) * ch * int(s.format.bits) / 8
				silent := flags&audclntBufferFlagsSilent != 0
				if !printed {
					printed = true
					firstFlags = flags
					if silent || data == nil {
						t.Logf("[B/C] first packet: frames=%d flags=0x%x SILENT-flagged (raw bytes unavailable; zeros contract applies)", frames, flags)
					} else if bytesAvail >= 32 {
						i16 := unsafe.Slice((*int16)(data), bytesAvail/2)
						n16 := 16
						if len(i16) < n16 {
							n16 = len(i16)
						}
						first16 = append([]int16(nil), i16[:n16]...)
						f32 := unsafe.Slice((*float32)(data), bytesAvail/4)
						n8 := 8
						if len(f32) < n8 {
							n8 = len(f32)
						}
						first8f = append([]float32(nil), f32[:n8]...)
						t.Logf("[B/C] first packet: frames=%d flags=0x%x bytes=%d (bits=%d ch=%d)",
							frames, flags, bytesAvail, s.format.bits, ch)
					}
				}
				if !silent && data == nil {
					t.Fatalf("GetBuffer returned nil data with flags=0x%x, no SILENT flag", flags)
				}
				// [D] production converter, same call shape as Drain.
				mono = append(mono, convertPacket(data, int(frames), ch, floating, silent)...)
				for len(mono) >= FrameSamples {
					win := mono[:FrameSamples]
					lvl := LevelOfSamples(win)
					windows = append(windows, lvl)
					if len(windows) == 1 {
						fl, ferr := FrameLevel(EncodeFrame(win))
						if ferr != nil {
							t.Fatalf("FrameLevel on own window: %v", ferr)
						}
						t.Logf("[D] window 1 FrameLevel = %.6f", fl)
					}
					mono = mono[FrameSamples:]
				}
			}
			if hr, _, _ := comCall(s.capture, 4, uintptr(frames)); !hrOK(hr) {
				t.Fatalf("ReleaseBuffer hr=%s", hresultString(hr))
			}
			if flags&audclntBufferFlagsSilent != 0 {
				silentPkts++
			}
		}
	}

	if !printed {
		t.Fatalf("no device packet arrived within the deadline (device dead or stream stalled)")
	}
	if len(windows) == 0 {
		t.Fatalf("packets arrived but no full seam window accumulated")
	}

	// [B]/[C] verbatim sample lines.
	if len(first16) > 0 {
		t.Logf("[B] first packet, int16 view, first %d: %v", len(first16), first16)
	}
	if len(first8f) > 0 {
		t.Logf("[C] same bytes, float32 view, first %d: %v", len(first8f), first8f)
	}

	// [D] summary over the collected windows (device rate; RMS is rate-blind
	// for identical waveforms, and the quad above says which rate this is).
	var sum, wmax, wmin float64
	wmin = math.MaxFloat64
	for _, l := range windows {
		sum += l
		if l > wmax {
			wmax = l
		}
		if l < wmin {
			wmin = l
		}
	}
	mean := sum / float64(len(windows))
	t.Logf("[D] %d windows x %d samples (rate %d): mean=%.6f max=%.6f min=%.6f silentPackets=%d firstFlags=0x%x -- ticket 247 floor was mean=0.367691 max=0.409285",
		len(windows), FrameSamples, s.format.rate, mean, wmax, wmin, silentPkts, firstFlags)

	// Verdict heuristic: a PRINT for the readings file, never an assertion.
	alt := false
	if len(first16) >= 16 {
		var maxMag int
		minMag := math.MaxInt
		altPairs := 0
		for _, v := range first16 {
			m := int(v)
			if m < 0 {
				m = -m
			}
			if m > maxMag {
				maxMag = m
			}
			if m > 0 && m < minMag {
				minMag = m
			}
		}
		for i := 0; i+1 < len(first16); i++ {
			a, b := int(first16[i]), int(first16[i+1])
			if a < 0 {
				a = -a
			}
			if b < 0 {
				b = -b
			}
			if a > 4*b || b > 4*a {
				altPairs++
			}
		}
		alt = altPairs >= 8 && maxMag >= 6000 && minMag <= maxMag/4
	}
	fclean := false
	if len(first8f) > 0 {
		fclean = true
		for _, v := range first8f {
			if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) || math.Abs(float64(v)) > 1.5 {
				fclean = false
			}
		}
	}
	switch {
	case alt && fclean:
		t.Logf("[verdict] int16 view alternates large/small while float32 view is clean: byte-slicing signature (float32 packet parsed as int16)")
	case !alt && len(first8f) > 0 && !fclean:
		t.Logf("[verdict] int16 view smooth, float32 view garbage: reversed shape (packet really is int16)")
	case !alt && fclean && mean >= 0.30 && mean <= 0.45:
		t.Logf("[verdict] both views smooth and window mean lands near the 0.37 floor: device-floor candidate (AC#2 form 3)")
	default:
		t.Logf("[verdict] no signature matched; the printed faces above are the evidence for AC#2")
	}
}
