//go:build windows

package audio

// Raw WASAPI plumbing via x/sys COM vtable calls (no cgo): the technique is
// the one proven in the S0 spike (scripts/spike/common/winshell.go, ticket
// 02) and reused by internal/ball. Vtable slot indices follow the MIDL
// declaration order of mmdeviceapi.h / audioclient.h and are named per call.
//
// unsafe discipline (vet-clean, same as internal/ball): COM object pointers
// stay unsafe.Pointer throughout; uintptr only ever appears as syscall
// arguments, and pointer arithmetic uses unsafe.Add.
//
// Threading contract: wasapiOpener.Open runs on the pinned capture thread,
// which owns every COM object it creates (IAudioClient/IAudioCaptureClient
// are used from that thread only).

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/observe"
)

var (
	modOle32    = windows.NewLazySystemDLL("ole32.dll")
	modKernel32 = windows.NewLazySystemDLL("kernel32.dll")

	pCoCreateInstance   = modOle32.NewProc("CoCreateInstance")
	pCoInitializeEx     = modOle32.NewProc("CoInitializeEx")
	pCoUninitialize     = modOle32.NewProc("CoUninitialize")
	pCoTaskMemFree      = modOle32.NewProc("CoTaskMemFree")
	pPropVariantClear   = modOle32.NewProc("PropVariantClear")
	pGetCurrentThreadId = modKernel32.NewProc("GetCurrentThreadId")
)

// currentThreadID returns the Win32 thread id of the calling thread (the
// pinned-loop identity probe for D38a diagnostics and tests).
func currentThreadID() uint32 {
	id, _, _ := pGetCurrentThreadId.Call()
	return uint32(id)
}

// COM identifiers (from mmdeviceapi.h / audioclient.h).
var (
	// IID_IUnknown {00000000-0000-0000-C000-000000000046}
	iidIUnknown = windows.GUID{
		Data1: 0x00000000, Data2: 0x0000, Data3: 0x0000,
		Data4: [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46},
	}
	// CLSID_MMDeviceEnumerator {BCDE0395-E52F-467C-8E3D-C4579291692E}
	clsidMMDeviceEnumerator = windows.GUID{
		Data1: 0xBCDE0395, Data2: 0xE52F, Data3: 0x467C,
		Data4: [8]byte{0x8E, 0x3D, 0xC4, 0x57, 0x92, 0x91, 0x69, 0x2E},
	}
	// IID_IMMDeviceEnumerator {A95664D2-9614-4F35-A746-DE8DB63617E6}
	iidIMMDeviceEnumerator = windows.GUID{
		Data1: 0xA95664D2, Data2: 0x9614, Data3: 0x4F35,
		Data4: [8]byte{0xA7, 0x46, 0xDE, 0x8D, 0xB6, 0x36, 0x17, 0xE6},
	}
	// IID_IAudioClient {1CB9AD4C-DBFA-4c32-B178-C2F568A703B2}
	iidIAudioClient = windows.GUID{
		Data1: 0x1CB9AD4C, Data2: 0xDBFA, Data3: 0x4C32,
		Data4: [8]byte{0xB1, 0x78, 0xC2, 0xF5, 0x68, 0xA7, 0x03, 0xB2},
	}
	// IID_IAudioCaptureClient {C8ADBD64-E71E-48a0-A4DE-185C395CD317}
	iidIAudioCaptureClient = windows.GUID{
		Data1: 0xC8ADBD64, Data2: 0xE71E, Data3: 0x48A0,
		Data4: [8]byte{0xA4, 0xDE, 0x18, 0x5C, 0x39, 0x5C, 0xD3, 0x17},
	}
	// PKEY_Device_FriendlyName: fmtid {A45C254E-DF1C-4EFD-8020-67D146A850E0} pid 14
	pkeyDeviceFriendlyName = propKey{
		fmtid: windows.GUID{
			Data1: 0xA45C254E, Data2: 0xDF1C, Data3: 0x4EFD,
			Data4: [8]byte{0x80, 0x20, 0x67, 0xD1, 0x46, 0xA8, 0x50, 0xE0},
		},
		pid: 14,
	}
)

// COM constants.
const (
	clsctxAll = 0x17 // CLSCTX_ALL

	coinitMTA = 0x0 // COINIT_MULTITHREADED

	eRender  = 0
	eCapture = 1
	eConsole = 0

	stgmRead = 0

	vtLPWSTR = 31

	waveFormatPCM   = 1
	waveFormatFloat = 3
	waveFormatExt   = 0xFFFE

	audclntShareModeShared    = 0
	audclntStreamFlagsEventCB = 0x00040000
	audclntBufferFlagsSilent  = 0x2
	audclntBufferHns          = 1_000_000 // 100ms in 100ns units

	waitObject0 = 0
	waitTimeout = 0x00000102
	waitFailed  = 0xFFFFFFFF
)

// propKey is a PROPERTYKEY.
type propKey struct {
	fmtid windows.GUID
	pid   uint32
}

// propVariant is the minimal PROPVARIANT layout (x64: vt + reserved + union
// at offset 8, 24 bytes total) needed to read VT_LPWSTR.
type propVariant struct {
	vt    uint16
	_     [6]byte
	union [16]byte
}

// comCall invokes vtable slot `slot` (0-based; IUnknown::QueryInterface = 0)
// of the COM object at `this`. Same shape as internal/ball's comCall: object
// pointers stay unsafe.Pointer, so no uintptr->pointer conversion happens
// here (vet-clean; unsafeptr check).
func comCall(this unsafe.Pointer, slot uintptr, args ...uintptr) (uintptr, uintptr, error) {
	vt := *(*unsafe.Pointer)(this)
	fn := *(*unsafe.Pointer)(unsafe.Add(vt, slot*8))
	return syscall.SyscallN(uintptr(fn), append([]uintptr{uintptr(this)}, args...)...)
}

// comRelease calls IUnknown::Release (slot 2) and nils the pointer.
func comRelease(obj *unsafe.Pointer) {
	if *obj != nil {
		_, _, _ = comCall(*obj, slotIUnknownRelease)
		*obj = nil
	}
}

const slotIUnknownRelease = 2

func hrOK(hr uintptr) bool { return int32(hr) >= 0 }

// comInitMTA initializes COM (MTA) on the calling thread; returns true when
// a balancing CoUninitialize is owed.
func comInitMTA() bool {
	hr, _, _ := pCoInitializeEx.Call(0, coinitMTA)
	return hrOK(hr) // S_OK or S_FALSE both owe the balancing uninit
}

func comUninitMTA() { pCoUninitialize.Call() }

// createEnumerator CoCreates a fresh IMMDeviceEnumerator on the calling
// thread (fresh-per-call avoids any cross-thread object sharing). The
// release func is idempotent.
func createEnumerator() (unsafe.Pointer, func(), error) {
	var enum unsafe.Pointer
	hr, _, _ := pCoCreateInstance.Call(
		uintptr(unsafe.Pointer(&clsidMMDeviceEnumerator)),
		0,
		clsctxAll,
		uintptr(unsafe.Pointer(&iidIMMDeviceEnumerator)),
		uintptr(unsafe.Pointer(&enum)))
	if !hrOK(hr) || enum == nil {
		return nil, nil, observe.New(observe.ClassAudioDevice,
			"CoCreateInstance(MMDeviceEnumerator) failed ("+hresultString(hr)+")")
	}
	return enum, func() { comRelease(&enum) }, nil
}

// waveFormat is the negotiated mix-format subset the loop needs.
type waveFormat struct {
	tag      uint16 // resolved to PCM(1) or FLOAT(3) for extensible formats
	channels uint16
	rate     uint32
	bits     uint16
}

// parseWaveFormat reads WAVEFORMATEX (+ WAVEFORMATEXTENSIBLE SubFormat) from
// a CoTaskMem-allocated format pointer. Layout: tag@0, channels@2, rate@4,
// bits@14, cbSize@16, [validBits@18, channelMask@22, SubFormat@26].
func parseWaveFormat(p unsafe.Pointer) (waveFormat, error) {
	if p == nil {
		return waveFormat{}, observe.New(observe.ClassAudioDevice, "GetMixFormat returned no format")
	}
	f := waveFormat{
		tag:      *(*uint16)(p),
		channels: *(*uint16)(unsafe.Add(p, 2)),
		rate:     *(*uint32)(unsafe.Add(p, 4)),
		bits:     *(*uint16)(unsafe.Add(p, 14)),
	}
	if f.tag == waveFormatExt {
		sub := *(*windows.GUID)(unsafe.Add(p, 26))
		f.tag = uint16(sub.Data1) // KSDATAFORMAT_SUBTYPE_* low word: 1=PCM 3=FLOAT
	}
	if f.rate == 0 || f.channels == 0 {
		return f, observe.New(observe.ClassAudioDevice, "degenerate mix format")
	}
	return f, nil
}

// wasapiOpener opens shared-mode capture streams (the real streamOpener).
type wasapiOpener struct{}

// Open implements streamOpener. It runs on the pinned capture thread, which
// has COM initialized (comInitMTA in run).
func (wasapiOpener) Open(dev DeviceDescriptor) (deviceStream, error) {
	enum, release, err := createEnumerator()
	if err != nil {
		return nil, err
	}
	defer release()

	// IMMDeviceEnumerator::GetDevice (slot 5) -> IMMDevice for the endpoint id.
	id16, err := windows.UTF16PtrFromString(dev.ID)
	if err != nil {
		return nil, DeviceError(hrEOutOfBounds, dev, "endpoint id encoding failed")
	}
	var mmdev unsafe.Pointer
	if hr, _, _ := comCall(enum, 5, uintptr(unsafe.Pointer(id16)), uintptr(unsafe.Pointer(&mmdev))); !hrOK(hr) || mmdev == nil {
		return nil, DeviceError(hr, dev, "endpoint not found (GetDevice)")
	}
	defer comRelease(&mmdev)

	// IMMDevice::Activate (slot 3) -> IAudioClient.
	var client unsafe.Pointer
	if hr, _, _ := comCall(mmdev, 3,
		uintptr(unsafe.Pointer(&iidIAudioClient)), clsctxAll, 0,
		uintptr(unsafe.Pointer(&client))); !hrOK(hr) || client == nil {
		return nil, DeviceError(hr, dev, "Activate(IAudioClient)")
	}

	s := &wasapiStream{dev: dev, client: client}
	if err := s.init(); err != nil {
		_ = s.Close() // release the half-initialized client (idempotent)
		return nil, err
	}
	return s, nil
}

// wasapiStream is one open shared-mode capture endpoint.
type wasapiStream struct {
	dev     DeviceDescriptor
	client  unsafe.Pointer // IAudioClient
	capture unsafe.Pointer // IAudioCaptureClient
	ev      windows.Handle
	format  waveFormat
}

// init negotiates the mix format and starts event-driven shared capture.
// IAudioClient vtable slots: 3 Initialize, 8 GetMixFormat, 9 GetDevicePeriod,
// 10 Start, 11 Stop, 12 Reset, 13 SetEventHandle, 14 GetService.
func (s *wasapiStream) init() error {
	// GetMixFormat (slot 8) - allocated by COM, freed with CoTaskMemFree
	// after Initialize consumed it.
	var wfex unsafe.Pointer
	if hr, _, _ := comCall(s.client, slotIAudioClientGetMixFormat, uintptr(unsafe.Pointer(&wfex))); !hrOK(hr) || wfex == nil {
		return DeviceError(hr, s.dev, "GetMixFormat")
	}
	f, ferr := parseWaveFormat(wfex)
	if ferr != nil {
		pCoTaskMemFree.Call(uintptr(wfex))
		return ferr
	}
	s.format = f

	// Initialize (slot 3): shared mode, event callback, ~100ms buffer. In
	// shared mode the buffer must align to the device period; on
	// AUDCLNT_E_BUFFER_SIZE_NOT_ALIGNED retry with the aligned duration
	// (Microsoft's documented recipe), falling back to the system default.
	hr, _, _ := comCall(s.client, slotIAudioClientInitialize,
		audclntShareModeShared, audclntStreamFlagsEventCB,
		audclntBufferHns, 0, uintptr(wfex), 0)
	if uint32(hr) == hrAUDCLNTBufNotAligned {
		aligned := alignedBufferHns(s.client, uint32(f.rate))
		hr, _, _ = comCall(s.client, slotIAudioClientInitialize,
			audclntShareModeShared, audclntStreamFlagsEventCB,
			aligned, 0, uintptr(wfex), 0)
	}
	if !hrOK(hr) {
		pCoTaskMemFree.Call(uintptr(wfex))
		return DeviceError(hr, s.dev, "Initialize(shared mode)")
	}
	pCoTaskMemFree.Call(uintptr(wfex))

	// Event handle (slot 13 SetEventHandle) for period-driven draining.
	ev, err := windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return observe.Wrap(observe.ClassAudioDevice, err, "CreateEvent for capture")
	}
	s.ev = ev
	if hr, _, _ := comCall(s.client, slotIAudioClientSetEventHandle, uintptr(ev)); !hrOK(hr) {
		return DeviceError(hr, s.dev, "SetEventHandle")
	}

	// GetService (slot 14) -> IAudioCaptureClient.
	if hr, _, _ := comCall(s.client, slotIAudioClientGetService,
		uintptr(unsafe.Pointer(&iidIAudioCaptureClient)),
		uintptr(unsafe.Pointer(&s.capture))); !hrOK(hr) || s.capture == nil {
		return DeviceError(hr, s.dev, "GetService(IAudioCaptureClient)")
	}

	// Start (slot 10).
	if hr, _, _ := comCall(s.client, slotIAudioClientStart); !hrOK(hr) {
		return DeviceError(hr, s.dev, "Start capture")
	}
	return nil
}

// alignedBufferHns computes the period-aligned buffer duration for the
// AUDCLNT_E_BUFFER_SIZE_NOT_ALIGNED retry: ceil(target/period) periods in
// 100ns units. Returns 0 (system default) when the period query fails.
func alignedBufferHns(client unsafe.Pointer, rate uint32) uintptr {
	// IAudioClient::GetDevicePeriod (slot 9): default + minimum, 100ns units.
	var defPeriod, minPeriod int64
	if hr, _, _ := comCall(client, slotIAudioClientGetDevicePeriod,
		uintptr(unsafe.Pointer(&defPeriod)), uintptr(unsafe.Pointer(&minPeriod))); !hrOK(hr) ||
		defPeriod <= 0 || rate == 0 {
		return 0
	}
	targetFrames := uint64(audclntBufferHns) * uint64(rate) / 10_000_000
	periodFrames := uint64(defPeriod) * uint64(rate) / 10_000_000
	if periodFrames == 0 {
		return 0
	}
	alignedFrames := (targetFrames + periodFrames - 1) / periodFrames * periodFrames
	return uintptr(alignedFrames * 10_000_000 / uint64(rate))
}

// Descriptor implements deviceStream.
func (s *wasapiStream) Descriptor() DeviceDescriptor { return s.dev }

// Rate implements deviceStream (negotiated mix-format rate).
func (s *wasapiStream) Rate() int { return int(s.format.rate) }

// WaitEvent implements deviceStream: block up to timeoutMs for the next
// device period.
func (s *wasapiStream) WaitEvent(timeoutMs uint32) (bool, error) {
	if s.ev == 0 {
		return false, DeviceError(hrAUDCLNTNotInit, s.dev, "WaitEvent on closed stream")
	}
	state, err := windows.WaitForSingleObject(s.ev, timeoutMs)
	if err != nil {
		return false, DeviceError(waitFailed, s.dev, "WaitForSingleObject")
	}
	if state == waitTimeout {
		return false, nil
	}
	if state != waitObject0 {
		return false, DeviceError(uintptr(state), s.dev, "WaitForSingleObject")
	}
	return true, nil
}

// Drain implements deviceStream: pull every pending packet from the capture
// client (vtable: 3 GetBuffer, 4 ReleaseBuffer, 5 GetNextPacketSize),
// converting to mono int16. A SILENT-flagged period comes back as zeros.
func (s *wasapiStream) Drain() ([]int16, error) {
	ch := int(s.format.channels)
	floating := s.format.tag == waveFormatFloat
	var out []int16
	for {
		var packet uint32
		if hr, _, _ := comCall(s.capture, 5, uintptr(unsafe.Pointer(&packet))); !hrOK(hr) {
			return out, DeviceError(hr, s.dev, "GetNextPacketSize")
		}
		if packet == 0 {
			return out, nil
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
			return out, DeviceError(hr, s.dev, "GetBuffer")
		}
		if frames > 0 {
			out = append(out, convertPacket(data, int(frames), ch, floating, flags&audclntBufferFlagsSilent != 0)...)
		}
		if hr, _, _ := comCall(s.capture, 4, uintptr(frames)); !hrOK(hr) {
			return out, DeviceError(hr, s.dev, "ReleaseBuffer")
		}
	}
}

// convertPacket turns one interleaved device packet into mono int16.
func convertPacket(data unsafe.Pointer, frames, channels int, floating, silent bool) []int16 {
	if silent || data == nil {
		return make([]int16, frames)
	}
	if floating {
		f := unsafe.Slice((*float32)(data), frames*channels)
		return MonoDownmix(FloatToPCM16(f), channels)
	}
	i := unsafe.Slice((*int16)(data), frames*channels)
	return MonoDownmix(i, channels)
}

// Close implements deviceStream: stop capture, release COM objects, close
// the event. Idempotent and best-effort (D38e step 4 needs this fast).
func (s *wasapiStream) Close() error {
	comRelease(&s.capture)
	if s.client != nil {
		_, _, _ = comCall(s.client, slotIAudioClientStop) // IAudioClient::Stop
		comRelease(&s.client)
	}
	if s.ev != 0 {
		_ = windows.CloseHandle(s.ev)
		s.ev = 0
	}
	return nil
}

// IAudioClient vtable slots (IUnknown 0-2).
const (
	slotIAudioClientInitialize      = 3
	slotIAudioClientGetMixFormat    = 8
	slotIAudioClientGetDevicePeriod = 9
	slotIAudioClientStart           = 10
	slotIAudioClientStop            = 11
	slotIAudioClientSetEventHandle  = 13
	slotIAudioClientGetService      = 14
)
