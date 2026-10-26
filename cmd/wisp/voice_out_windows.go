//go:build windows

package main

// The announce path's last hop: samples to the loudspeaker.
//
// WHY THIS FILE EXISTS AT ALL. internal/audio owns the C8 capture seam and its
// doc.go carries the marker this file answers:
//
//	// DEFERRED(playback): TTS output lands with ticket 26 on the same WASAPI
//	//   seam...
//
// Ticket 26 landed the synthesis (internal/speech's MatchaTts) and 26-r1 said in
// as many words that it could not ship the render side: `grep '\.Speak(' cmd
// internal` had no production caller. Ticket wire-r1 is that caller, and a
// caller that hands samples to a callback nobody renders is exactly the
// decoration 26-r1 refused to add. So the render side lands here, in the package
// that owns the composition root, and it lands on winmm's waveOut rather than
// WASAPI:
//
//   - waveOut is the OS's own PCM queue. It needs no COM apartment, no
//     IMMDevice activation, no MTA thread and no callback pointer from Go into C
//     (which would be an illegal cgo-visible function pointer). Every one of
//     those is a failure mode this leg would have to test instead of the
//     announcement.
//   - the chain is 16 kHz mono int16 on the capture side and 22.05 kHz mono
//     float32 on the synthesis side, which is well inside what waveOut has
//     supported since Windows 3.1. There is no format negotiation to lose.
//   - it is the same "one seam, one implementation" shape the rest of this tree
//     uses: voiceSpeaker is two methods, so the tests can hand the announce path
//     a recorder and read the samples back without a sound card.
//
// WHAT IT DOES NOT DO. It does not mix, does not resample on the fly (it picks a
// device format up front from the engine's own rate and falls back through the
// rates every Windows box supports), does not keep a thread of its own (the
// owner is the voice chain's single goroutine, D38b), and does not measure
// anything: the bookkeeping the announce path reports is sample counts, never a
// wall clock (D42#9). Barge-in (D47 Path C) stops a sentence through the SYNTH
// side's cooperative abort - the seam's emit returns an error - and the same
// error path calls Abort here to drop the queued buffers, which is the
// `tts.stop` / `audio.release-output` side effect rows #27/#41 name.

import (
	"errors"
	"fmt"
	"sync"
	"syscall"
	"unsafe"

	"github.com/CarlosShao/wisp/internal/audio"
)

// voiceSpeaker is the render seam the announce path is written against.
//
// Play is called once per chunk the TTS engine emits, in playback order, and it
// returns as soon as the chunk is QUEUED (waveOut plays it asynchronously). Drain
// is what makes the announcement's end observable: it returns after the last
// queued chunk has been played, which is the moment D43 #26 (EvSpeakDone) is
// legal to dispatch. Abort drops whatever is queued, for the interrupt rows.
type voiceSpeaker interface {
	// SampleRate is the device rate this speaker was opened at. A chunk from an
	// engine running at a different rate is resampled to it.
	SampleRate() int
	// Play queues one chunk of float32 mono samples in -1..+1.
	Play(samples []float32, engineRate int) error
	// Drain waits until everything queued has been played.
	Drain() error
	// Abort drops the queue and returns the number of chunks dropped.
	Abort() int
	// Close releases the device (C11: idempotent).
	Close() error
}

// The winmm entry points. They are resolved lazily so a machine whose winmm.dll
// refuses to load produces a named error at the announce, not a boot failure:
// the voice chain runs without a speaker, it just cannot be heard.
var (
	winmmDLL            = syscall.NewLazyDLL("winmm.dll")
	procWaveOutOpen     = winmmDLL.NewProc("waveOutOpen")
	procWaveOutClose    = winmmDLL.NewProc("waveOutClose")
	procWaveOutPrepare  = winmmDLL.NewProc("waveOutPrepareHeader")
	procWaveOutUnprep   = winmmDLL.NewProc("waveOutUnprepareHeader")
	procWaveOutWrite    = winmmDLL.NewProc("waveOutWrite")
	procWaveOutReset    = winmmDLL.NewProc("waveOutReset")
	procWaveOutGetNumDe = winmmDLL.NewProc("waveOutGetNumDevs")
	kernel32DLL         = syscall.NewLazyDLL("kernel32.dll")
	procCreateEventW    = kernel32DLL.NewProc("CreateEventW")
	procWaitForSingle   = kernel32DLL.NewProc("WaitForSingleObject")
	procCloseHandle     = kernel32DLL.NewProc("CloseHandle")
)

const (
	waveMapper     = 0xFFFFFFFF
	waveFormatPCM  = 1
	callbackEvent  = 0x00050000 // CALLBACK_EVENT
	whdrDone       = 0x00000001
	mmSysErrNoErr  = 0
	waveOutBuffers = 2 // one playing, one being filled: the smallest queue that
	//                            keeps the device fed between TTS chunks
	// voiceOutWaitGrace is added to a chunk's own duration when the wait for the
	// device is bounded. The bound is derived from the AUDIO clock (samples to
	// duration) and never from "how long has it been" (D42#9); the grace is what
	// covers the driver's own start-up, not a deadline anybody is racing.
	voiceOutWaitGrace = 2000 // milliseconds
)

// waveFormatEx is WAVEFORMATEX (mmreg.h). Size 18 on the wire; Go pads the
// struct to 20 and only the first 18 bytes are read.
type waveFormatEx struct {
	FormatTag      uint16
	Channels       uint16
	SamplesPerSec  uint32
	AvgBytesPerSec uint32
	BlockAlign     uint16
	BitsPerSample  uint16
	CBytes         uint16
}

// waveHdr is WAVEHDR (mmsystem.h), 48 bytes on amd64.
type waveHdr struct {
	Data          uintptr
	BufferLength  uint32
	BytesRecorded uint32
	User          uintptr
	Flags         uint32
	Loops         uint32
	Next          uintptr
	Reserved      uintptr
}

// waveOutSpeaker is the winmm implementation of voiceSpeaker.
//
// The buffers are a fixed ring of waveOutBuffers: a chunk reuses the slot that is
// oldest, waiting for it first. That is the whole flow control - no goroutine, no
// channel, no timer - which is what keeps this file inside D38b's "one owner
// goroutine" rule: the caller is the owner and Play blocks until there is room.
type waveOutSpeaker struct {
	mu     sync.Mutex
	handle uintptr
	event  uintptr
	rate   int
	closed bool

	hdrs    []waveHdr
	bufs    [][]int16
	pending []bool // buf i was written and has not been waited for yet
	next    int
	playing int
	dropped int

	failed error // a device that failed once is not retried: the announce says so
}

// waveOutError renders one MMRESULT the way the announce line needs it: the code
// and the constant that names it, so a log line can be looked up without a copy
// of mmsystem.h.
func waveOutError(rc uintptr) string {
	name := "MMSYSERR_UNKNOWN"
	switch rc {
	case 0:
		name = "MMSYSERR_NOERROR"
	case 2:
		name = "MMSYSERR_BADDEVICEID"
	case 4:
		name = "MMSYSERR_ALLOCATED"
	case 5:
		name = "MMSYSERR_INVALHANDLE"
	case 6:
		name = "MMSYSERR_NODRIVER"
	case 7:
		name = "MMSYSERR_NOMEM"
	case 8:
		name = "MMSYSERR_NOTSUPPORTED"
	case 11:
		name = "MMSYSERR_NOTENABLED"
	case 32:
		name = "WAVERR_BADFORMAT"
	case 33:
		name = "WAVERR_STILLPLAYING"
	}
	return fmt.Sprintf("%s (%d)", name, rc)
}

// voiceOutRates is the fallback ladder. The TTS engine's own rate comes first
// because it needs no resampling at all; the rest are the rates every Windows
// audio device supports, so a device that refuses 22.05 kHz still gets the
// announcement instead of a refusal the user cannot act on.
var voiceOutRates = []int{48000, 44100, 16000}

// openWaveOutSpeaker opens the default output device.
//
// engineRate is the synthesis rate; it is tried first. A machine with no output
// device at all answers a named error, and the announce path degrades to the
// text notification rather than pretending it spoke (ticket 26's "text-only
// fallback still delivers the result").
func openWaveOutSpeaker(engineRate int) (voiceSpeaker, error) {
	if err := procWaveOutGetNumDe.Find(); err != nil {
		return nil, fmt.Errorf("waveOut is not available in this process: %w", err)
	}
	n, _, _ := procWaveOutGetNumDe.Call()
	if n == 0 {
		return nil, errors.New("this machine reports no waveOut output device: the announcement cannot be heard")
	}
	var firstErr error
	rates := append([]int{engineRate}, voiceOutRates...)
	for _, rate := range rates {
		if rate <= 0 {
			continue
		}
		sp, err := newWaveOutSpeaker(rate)
		if err == nil {
			return sp, nil
		}
		if firstErr == nil {
			firstErr = err
		}
	}
	return nil, firstErr
}

func newWaveOutSpeaker(rate int) (*waveOutSpeaker, error) {
	ev, _, evErr := procCreateEventW.Call(0, 0, 0, 0)
	if ev == 0 {
		return nil, fmt.Errorf("CreateEvent for the waveOut callback failed: %v", evErr)
	}
	wfx := waveFormatEx{
		FormatTag:      waveFormatPCM,
		Channels:       1,
		SamplesPerSec:  uint32(rate),
		AvgBytesPerSec: uint32(rate * 2),
		BlockAlign:     2,
		BitsPerSample:  16,
		CBytes:         0,
	}
	var handle uintptr
	rc, _, _ := procWaveOutOpen.Call(
		uintptr(unsafe.Pointer(&handle)),
		uintptr(waveMapper),
		uintptr(unsafe.Pointer(&wfx)),
		ev,
		0,
		uintptr(callbackEvent),
	)
	if uintptr(rc) != mmSysErrNoErr {
		_, _, _ = procCloseHandle.Call(ev)
		return nil, fmt.Errorf("waveOutOpen at %d Hz mono 16-bit refused: %s", rate, waveOutError(uintptr(rc)))
	}
	sp := &waveOutSpeaker{handle: handle, event: ev, rate: rate}
	sp.hdrs = make([]waveHdr, waveOutBuffers)
	sp.bufs = make([][]int16, waveOutBuffers)
	sp.pending = make([]bool, waveOutBuffers)
	return sp, nil
}

func (s *waveOutSpeaker) SampleRate() int { return s.rate }

// Play converts one chunk and queues it, waiting for the slot it reuses first.
func (s *waveOutSpeaker) Play(samples []float32, engineRate int) error {
	if s == nil {
		return errors.New("voice: no speaker")
	}
	if len(samples) == 0 {
		return nil
	}
	pcm := audio.FloatToPCM16(samples)
	if engineRate > 0 && engineRate != s.rate {
		pcm = audio.ResampleLinear(pcm, engineRate, s.rate)
	}
	if len(pcm) == 0 {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return errors.New("voice: the speaker is already closed")
	}
	if s.failed != nil {
		return s.failed
	}

	i := s.next
	s.next = (s.next + 1) % waveOutBuffers
	if err := s.releaseLocked(i); err != nil {
		s.failed = err
		return err
	}

	// The buffer must outlive the header, and the header must not move: both are
	// fields of the speaker, and the header is passed by address, so the slice
	// element is the stable home Windows writes its flags into.
	s.bufs[i] = pcm
	hdr := &s.hdrs[i]
	*hdr = waveHdr{
		Data:         uintptr(unsafe.Pointer(&s.bufs[i][0])),
		BufferLength: uint32(len(pcm) * 2),
	}
	if rc, _, _ := procWaveOutPrepare.Call(s.handle, uintptr(unsafe.Pointer(hdr)), unsafe.Sizeof(*hdr)); uintptr(rc) != mmSysErrNoErr {
		return fmt.Errorf("waveOutPrepareHeader refused chunk %d: %s", i, waveOutError(uintptr(rc)))
	}
	if rc, _, _ := procWaveOutWrite.Call(s.handle, uintptr(unsafe.Pointer(hdr)), unsafe.Sizeof(*hdr)); uintptr(rc) != mmSysErrNoErr {
		_, _, _ = procWaveOutUnprep.Call(s.handle, uintptr(unsafe.Pointer(hdr)), unsafe.Sizeof(*hdr))
		return fmt.Errorf("waveOutWrite refused chunk %d: %s", i, waveOutError(uintptr(rc)))
	}
	s.pending[i] = true
	s.playing++
	return nil
}

// releaseLocked waits out one buffer and unprepares it. The wait is bounded by
// the buffer's OWN duration on the audio clock plus one fixed grace: a device
// that stops signalling is reported, never waited on forever.
func (s *waveOutSpeaker) releaseLocked(i int) error {
	if !s.pending[i] {
		return nil
	}
	hdr := &s.hdrs[i]
	samples := len(s.bufs[i])
	budget := samples*1000/s.rate + voiceOutWaitGrace // samples to milliseconds, audio clock
	rc, _, _ := procWaitForSingle.Call(s.event, uintptr(budget))
	switch uintptr(rc) {
	case 0: // WAIT_OBJECT_0
	case 0x102: // WAIT_TIMEOUT
		return fmt.Errorf("waveOut did not finish chunk %d within its own %d ms of audio plus %d ms grace: the device stopped signalling",
			i, samples*1000/s.rate, voiceOutWaitGrace)
	default:
		return fmt.Errorf("waiting for waveOut chunk %d returned %d", i, uintptr(rc))
	}
	if rc, _, _ := procWaveOutUnprep.Call(s.handle, uintptr(unsafe.Pointer(hdr)), unsafe.Sizeof(*hdr)); uintptr(rc) != mmSysErrNoErr && uintptr(rc) != 33 {
		return fmt.Errorf("waveOutUnprepareHeader refused chunk %d: %s", i, waveOutError(uintptr(rc)))
	}
	s.pending[i] = false
	s.bufs[i] = nil
	s.playing--
	return nil
}

// Drain waits for every queued chunk. This is the "the sentence has been
// spoken" signal the D43 walk hangs EvSpeakDone on.
func (s *waveOutSpeaker) Drain() error {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	if s.failed != nil {
		return s.failed
	}
	for i := range s.hdrs {
		if err := s.releaseLocked(i); err != nil {
			s.failed = err
			return err
		}
	}
	return nil
}

// Abort drops what is queued (rows #27/#41: tts.stop / audio.release-output).
func (s *waveOutSpeaker) Abort() int {
	if s == nil {
		return 0
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return 0
	}
	n := s.playing
	if s.handle != 0 {
		_, _, _ = procWaveOutReset.Call(s.handle)
	}
	for i := range s.hdrs {
		if s.pending[i] {
			_, _, _ = procWaveOutUnprep.Call(s.handle, uintptr(unsafe.Pointer(&s.hdrs[i])), unsafe.Sizeof(s.hdrs[i]))
			s.pending[i] = false
			s.bufs[i] = nil
		}
	}
	s.playing = 0
	s.dropped += n
	return n
}

// Close releases the device. Idempotent (C11).
func (s *waveOutSpeaker) Close() error {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	if s.handle != 0 {
		_, _, _ = procWaveOutReset.Call(s.handle)
		for i := range s.hdrs {
			if s.pending[i] {
				_, _, _ = procWaveOutUnprep.Call(s.handle, uintptr(unsafe.Pointer(&s.hdrs[i])), unsafe.Sizeof(s.hdrs[i]))
				s.pending[i] = false
			}
		}
		_, _, _ = procWaveOutClose.Call(s.handle)
		s.handle = 0
	}
	if s.event != 0 {
		_, _, _ = procCloseHandle.Call(s.event)
		s.event = 0
	}
	for i := range s.bufs {
		s.bufs[i] = nil
	}
	return nil
}
