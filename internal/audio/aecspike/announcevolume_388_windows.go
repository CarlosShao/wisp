//go:build windows

package aecspike

// Ticket-388 attenuation guard: the voice tests render REAL matcha speech on
// the REAL device (waveOut - the <=400ms interrupt measurement needs the true
// abort path), so the test sentence reaches the owner's speakers at whatever
// volume the process session currently has. The guard attenuates that
// session volume before such a test and puts the owner's original value back
// afterwards.
//
// Why the session volume, and why ONLY the cgo float path:
//   - ISimpleAudioVolume::SetMasterVolume takes a FLOAT parameter. On the
//     x64 Windows ABI float arguments ride the XMM registers, which
//     syscall.SyscallN cannot populate: the integer path once wrote a
//     garbage float that Windows PERSISTED as a near-zero session volume,
//     silencing every later render run (59-r1 instrument finding #2). The
//     guard therefore only drives the proven C shim (SetSessionVolume /
//     GetSessionVolume over aecsmoke.c); there is no integer path here.
//   - The guard opens its own shared-mode IAudioClient on the console
//     default render endpoint and reaches the stream's session through
//     IAudioClient::GetService. A stream initialized without a session GUID
//     joins the process default session - the same session the waveOut
//     announce path renders through - so SetMasterVolume on the guard's
//     client attenuates the announce output WITHOUT touching the waveOut
//     stream itself: the measured object (the abort-latency path) keeps its
//     exact semantics; only loudness changes.
//   - Restore axis: the guard saves the owner's volume BEFORE touching it
//     and verifies the read-back after every write (fail-closed: a failed
//     apply puts the original back before returning). The rig's shim also
//     unmutes the session on every SetSessionVolume call - that is
//     established instrument behavior (every RunLoopback run does it), and
//     the guard inherits it rather than widening the C surface.

import (
	"fmt"
	"math"
	"unsafe"
)

// DefaultAnnounceAttenuatedVolume is the session volume the voice tests
// apply unless the caller names another value: 15% of full scale (about
// -16.5 dB below a full-scale session) - plainly audible to an operator
// listening for it, far below startle level. The guard's value axis is
// parameterized: tests pass their own value when a reading needs one.
const DefaultAnnounceAttenuatedVolume = 0.15

// volumeReadbackEpsilon bounds the read-back agreement of every volume
// write. ISimpleAudioVolume stores float32, so 1e-6 is comfortably above
// the representation error and far below anything audible.
const volumeReadbackEpsilon = 1e-6

// VolumeReadings is the AC#1/AC#3 triple: the owner's volume as found, the
// attenuated value the session actually read back, and the restored value
// the session actually read back.
type VolumeReadings struct {
	Original   float64 `json:"original"`
	Attenuated float64 `json:"attenuated"`
	Restored   float64 `json:"restored"`
}

// AnnounceAttenuation is one applied attenuation. Its zero value is not
// usable; construct it only through ApplyAnnounceAttenuation, and always
// pair it with exactly one Restore (a test's t.Cleanup is the right owner).
type AnnounceAttenuation struct {
	client   unsafe.Pointer
	release  func()
	original float64
	applied  float64
}

// ApplyAnnounceAttenuation saves the process render session's current
// master volume and sets attenuated through the real C float call,
// verifying the read-back. On any failure it puts the saved volume back
// before returning, so a failed apply leaves the session as it found it.
func ApplyAnnounceAttenuation(attenuated float64) (*AnnounceAttenuation, error) {
	if attenuated <= 0 || attenuated > 1 {
		return nil, fmt.Errorf("aecspike: announce attenuation volume %g outside (0,1]", attenuated)
	}
	client, release, err := openSessionVolumeClient()
	if err != nil {
		return nil, err
	}
	original, err := GetSessionVolume(client)
	if err != nil {
		release()
		return nil, fmt.Errorf("aecspike: reading the original session volume: %w", err)
	}
	back, err := SetSessionVolume(client, attenuated)
	if err != nil {
		release()
		if _, rerr := SetSessionVolume(client, original); rerr != nil {
			return nil, fmt.Errorf("aecspike: apply to %g failed (%v) AND the restore to %g failed (%v); "+
				"the session volume may be left at the attenuated value - set it back manually",
				attenuated, err, original, rerr)
		}
		return nil, fmt.Errorf("aecspike: setting the attenuated session volume: %w", err)
	}
	if math.Abs(back-attenuated) > volumeReadbackEpsilon {
		release()
		if _, rerr := SetSessionVolume(client, original); rerr != nil {
			return nil, fmt.Errorf("aecspike: apply read back %g, wanted %g, AND the restore to %g failed (%v); "+
				"the session volume may be left at the attenuated value - set it back manually",
				back, attenuated, original, rerr)
		}
		return nil, fmt.Errorf("aecspike: attenuated volume read back %g, wanted %g", back, attenuated)
	}
	return &AnnounceAttenuation{
		client: client, release: release, original: original, applied: back,
	}, nil
}

// Restore puts the owner's original session volume back and verifies the
// read-back, then releases the guard's client (Restore is single-use; the
// readings carry the full triple either way).
func (a *AnnounceAttenuation) Restore() (VolumeReadings, error) {
	if a == nil || a.client == nil {
		return VolumeReadings{}, fmt.Errorf("aecspike: the announce attenuation guard is not applied")
	}
	readings := VolumeReadings{Original: a.original, Attenuated: a.applied}
	back, err := SetSessionVolume(a.client, a.original)
	readings.Restored = back
	a.release()
	a.client = nil
	if err != nil {
		return readings, fmt.Errorf("aecspike: restoring the session volume %g: %w", a.original, err)
	}
	if math.Abs(back-a.original) > volumeReadbackEpsilon {
		return readings, fmt.Errorf("aecspike: restored volume read back %g, wanted the original %g", back, a.original)
	}
	return readings, nil
}

// openSessionVolumeClient opens one shared-mode IAudioClient on the console
// default render endpoint (initialized, never started) so its session can be
// driven through IAudioClient::GetService(ISimpleAudioVolume). The returned
// release balances COM init, the device, the client and the enumerator, and
// is safe to call more than once.
func openSessionVolumeClient() (unsafe.Pointer, func(), error) {
	comMine := comInitMTA()
	released := false
	var (
		enum        unsafe.Pointer
		enumRelease func()
		dev         unsafe.Pointer
		client      unsafe.Pointer
	)
	// closeAll closes whatever was opened so far, so it is safe to run
	// after any prefix of the open steps below.
	closeAll := func() {
		comRelease(&client)
		comRelease(&dev)
		if enumRelease != nil {
			enumRelease()
		}
		if comMine {
			comUninitMTA()
		}
	}
	releaseAll := func() {
		if released {
			return
		}
		released = true
		closeAll()
	}
	enum, enumRelease, err := createEnumerator()
	if err != nil {
		releaseAll()
		return nil, nil, err
	}
	if hr, _, _ := comCall(enum, 4, eRender, roleConsole, uintptr(unsafe.Pointer(&dev))); !hrOK(hr) || dev == nil {
		releaseAll()
		return nil, nil, fmt.Errorf("aecspike: announce guard found no default render endpoint (hr=%s)", hrHex(hr))
	}
	s := &stream{dev: dev}
	if err := s.activate(); err != nil {
		releaseAll()
		return nil, nil, err
	}
	client = s.client
	if err := s.init(shareShared, 0); err != nil {
		releaseAll()
		return nil, nil, err
	}
	return client, releaseAll, nil
}
