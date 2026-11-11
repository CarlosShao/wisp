package audio

// Device enumeration seam (SPEC-04 sec 9: hotplug tests inject fake device
// changes at the enumeration interface - NEVER at the audio-data seam) and
// the D42#2/#12 error mapping from open failures to `Error(audio_device)`
// with the device name and user guidance text.

import (
	"strconv"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
)

// DeviceDescriptor identifies one audio endpoint. ID is the WASAPI endpoint
// id, Name the friendly name (best effort).
type DeviceDescriptor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (d DeviceDescriptor) String() string {
	if d.Name != "" {
		return d.Name
	}
	return d.ID
}

// DeviceWatcher is the enumeration seam: which capture endpoint is default,
// the paired render endpoint (D47: capture and render may be DIFFERENT
// endpoints - e.g. Bluetooth speaker + wired mic; the pair stays queryable
// for the P15 AEC spike), and default-device-change notifications.
type DeviceWatcher interface {
	// DefaultCaptureDevice returns the current default capture endpoint.
	DefaultCaptureDevice() (DeviceDescriptor, error)
	// DefaultRenderDevice returns the current default render endpoint.
	DefaultRenderDevice() (DeviceDescriptor, error)
	// OnDefaultCaptureChanged registers cb, invoked (from a Windows
	// notification thread or a test) when the default capture device
	// changes. The callback MUST be non-blocking and must not call COM.
	// Returns an unsubscribe function.
	OnDefaultCaptureChanged(cb func()) (unsub func(), err error)
	// Close releases watcher resources. Idempotent.
	Close() error
}

// streamOpener opens a capture stream on a specific endpoint (the seam the
// capture thread drives on open and on every hotplug reopen; tests inject a
// fake opener whose streams carry real wav data through the real loop).
type streamOpener interface {
	Open(dev DeviceDescriptor) (deviceStream, error)
}

// captureDeviceLister is the OPTIONAL enumeration face: a DeviceWatcher that
// can list the present capture endpoints carries it, and [audio] input_device
// resolves through it. The production watcher implements it over
// IMMDeviceEnumerator::EnumAudioEndpoints; a watcher that cannot enumerate
// keeps the default-endpoint contract, and a config that names a device on
// such a watcher fails loudly instead of silently opening the default.
type captureDeviceLister interface {
	CaptureDevices() ([]DeviceDescriptor, error)
}

// inputDeviceDefault is the [audio] input_device value that means "the OS
// default capture endpoint" (the schema default and the pre-wire shape).
const inputDeviceDefault = "default"

// selectCaptureDevice resolves [audio] input_device into one endpoint.
// name "" or "default" asks the watcher for its default endpoint (the exact
// pre-wire read); any other name goes through the watcher's enumeration face
// and matches by endpoint id first, friendly name second (both compared after
// whitespace trim, the name case-insensitively). A name that matches nothing
// is a ClassAudioDevice error that lists the present endpoints - "it fell
// back to the default" and "it was never written" must not look identical
// (the ticket-83 rule the config layer enforces for the other read-only
// keys), and a missing enumeration face is refused the same way.
func selectCaptureDevice(w DeviceWatcher, name string) (DeviceDescriptor, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" || trimmed == inputDeviceDefault {
		dev, err := w.DefaultCaptureDevice()
		if err != nil {
			return DeviceDescriptor{}, observe.Wrap(observe.ClassAudioDevice, err, "default capture endpoint")
		}
		return dev, nil
	}
	lister, ok := w.(captureDeviceLister)
	if !ok {
		return DeviceDescriptor{}, observe.New(observe.ClassAudioDevice,
			"capture device "+strconv.Quote(trimmed)+" was requested but this watcher cannot enumerate endpoints; remove audio.input_device or build against the MMDevice watcher")
	}
	devs, err := lister.CaptureDevices()
	if err != nil {
		return DeviceDescriptor{}, observe.Wrap(observe.ClassAudioDevice, err,
			"capture endpoint enumeration failed while resolving audio.input_device")
	}
	for _, d := range devs {
		if strings.TrimSpace(d.ID) == trimmed {
			return d, nil
		}
	}
	for _, d := range devs {
		if strings.EqualFold(strings.TrimSpace(d.Name), trimmed) {
			return d, nil
		}
	}
	names := make([]string, 0, len(devs))
	for _, d := range devs {
		names = append(names, d.String())
	}
	return DeviceDescriptor{}, observe.New(observe.ClassAudioDevice,
		"audio.input_device "+strconv.Quote(trimmed)+" matches no capture endpoint; present endpoints: "+strings.Join(names, ", "))
}

// deviceStream is the capture-thread view of one open endpoint: mono int16
// samples at the device's native rate, drained per device period.
type deviceStream interface {
	// Descriptor names the endpoint the stream was opened on.
	Descriptor() DeviceDescriptor
	// Rate is the negotiated device-native sample rate (per stream constant;
	// known at open time from the mix format).
	Rate() int
	// WaitEvent waits up to timeoutMs for the next device period.
	// ok=false means timeout (no data yet), not failure.
	WaitEvent(timeoutMs uint32) (ok bool, err error)
	// Drain returns all samples buffered since the last event, downmixed to
	// mono at Rate(). A device SILENT period comes back as zeros (silence
	// is data, not an error).
	Drain() (samples []int16, err error)
	// Close stops capture and releases the device. Idempotent.
	Close() error
}

// Both interfaces are the capture seam the windows WASAPI implementation and
// its hotplug tests drive; the linux build has no user for them, so the seam
// types are pinned here instead of being deleted out from under that side.
var (
	_ = streamOpener(nil)
	_ = deviceStream(nil)
)

// Windows HRESULTs the audio stack maps (D42#2/#12). Values from
// audioclient.h / winerror.h.
const (
	hrSOK                  = 0x00000000
	hrEAccessDenied        = 0x80070005
	hrEOutOfBounds         = 0x8000FFFF
	hrAUDCLNTNotInit       = 0x88890001
	hrAUDCLNTDeviceInv     = 0x88890004 // endpoint no longer present (hotplug)
	hrAUDCLNTUnsupported   = 0x88890008 // mix format unusable
	hrAUDCLNTDeviceInUse   = 0x8889000A // held exclusively by another app
	hrAUDCLNTEndpointFail  = 0x8889000E // endpoint creation failed
	hrAUDCLNTBufNotAligned = 0x88890014 // shared-mode buffer must align to device period
	hrCOENotInitialized    = 0x800401F0
)

// privacyGuidance is the D42#12 copy pointer: access-denied open failures
// must tell the user where to fix microphone permission. (Log-safe: device
// names only, never audio content - D16.)
const privacyGuidance = "microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone"

// inUseGuidance covers the exclusive-mode occupation case (D42#12).
const inUseGuidance = "capture device is held in exclusive mode by another application; close that application or pick another device"

// DeviceError builds the canonical D42 audio_device error for an open or
// capture failure on dev, mapping known HRESULTs to explicit guidance text.
// The device name is always included (silent failures are forbidden).
func DeviceError(hr uintptr, dev DeviceDescriptor, what string) *observe.Error {
	code := hresultString(hr)
	detail := what + " on device " + dev.String()
	switch hr {
	case hrEAccessDenied:
		return observe.New(observe.ClassAudioDevice, detail+" ("+code+"): "+privacyGuidance)
	case hrAUDCLNTDeviceInUse:
		return observe.New(observe.ClassAudioDevice, detail+" ("+code+"): "+inUseGuidance)
	case hrAUDCLNTDeviceInv:
		return observe.New(observe.ClassAudioDevice, detail+" ("+code+"): device invalidated (unplugged or disabled)")
	case hrAUDCLNTUnsupported:
		return observe.New(observe.ClassAudioDevice, detail+" ("+code+"): unsupported device mix format")
	case hrAUDCLNTEndpointFail:
		return observe.New(observe.ClassAudioDevice, detail+" ("+code+"): endpoint creation failed")
	default:
		return observe.New(observe.ClassAudioDevice, detail+" ("+code+")")
	}
}

// hresultString renders an HRESULT as 0x%08X.
func hresultString(hr uintptr) string {
	const hexdigits = "0123456789ABCDEF"
	buf := []byte{'0', 'x', 0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < 8; i++ {
		buf[9-i] = hexdigits[hr&0xF]
		hr >>= 4
	}
	return string(buf)
}
