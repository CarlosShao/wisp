//go:build windows

package audio

// MMDevice enumeration watcher (D42#2): default endpoint discovery, friendly
// names, and IMMNotificationClient::OnDefaultDeviceChanged bridged into a
// plain Go callback. Pure x/sys COM like the rest of the package.
//
// The IMMNotificationClient is a Go-synthesized COM object: a Go struct
// whose first word is a manually built vtable of windows.NewCallback
// trampolines. NewCallback accepts only top-level functions (no closures),
// so the callbacks recover their watcher through a registry keyed by the
// `this` pointer value; the callbacks run on a Windows audio-service thread
// and only record the event (non-blocking, no COM) before returning S_OK.

import (
	"log/slog"
	"strconv"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/observe"
)

// notifRegistry maps a synthesized client's `this` pointer to its watcher.
var notifRegistry sync.Map // uintptr -> *mmDeviceWatcher

// mmDeviceWatcher is the real DeviceWatcher.
type mmDeviceWatcher struct {
	mu     sync.Mutex
	nc     *notifClient // keeps the synthesized COM object alive
	this   uintptr      // the pointer handed to COM (registry key)
	reg    unsafe.Pointer
	cbs    []*cbEntry
	closed bool
}

// cbEntry is one registered callback (pointer identity = unsubscribe key).
type cbEntry struct {
	fn func()
}

// newMMDeviceWatcher creates the watcher and registers the notification
// client. COM is initialized MTA on the calling thread; the balancing
// CoUninitialize is deliberately skipped (process-resident watcher, like
// the ball STA thread).
func newMMDeviceWatcher() *mmDeviceWatcher {
	w := &mmDeviceWatcher{}
	comInitMTA()

	nc := newNotifClient(w)
	w.nc = nc
	w.this = uintptr(unsafe.Pointer(nc))
	notifRegistry.Store(w.this, w)

	enum, release, err := createEnumerator()
	if err != nil {
		slog.Error("audio watcher: enumerator unavailable (hotplug disabled)", "error", err.Error())
		return w
	}
	// IMMDeviceEnumerator::RegisterEndpointNotificationCallback (slot 6).
	if hr, _, _ := comCall(enum, 6, w.this); !hrOK(hr) {
		slog.Error("audio watcher: notification registration failed (hotplug disabled)",
			"hr", hresultString(hr))
		release()
		return w
	}
	w.reg = enum
	return w
}

// defaultEndpoint resolves the default endpoint for a data flow, with a
// fresh enumerator per call (no cross-thread object sharing).
func (w *mmDeviceWatcher) defaultEndpoint(flow uint32) (DeviceDescriptor, error) {
	enum, release, err := createEnumerator()
	if err != nil {
		return DeviceDescriptor{}, err
	}
	defer release()

	// GetDefaultAudioEndpoint (slot 4).
	var dev unsafe.Pointer
	if hr, _, _ := comCall(enum, 4, uintptr(flow), uintptr(eConsole), uintptr(unsafe.Pointer(&dev))); !hrOK(hr) || dev == nil {
		return DeviceDescriptor{}, observe.New(observe.ClassAudioDevice,
			"GetDefaultAudioEndpoint(flow="+strconv.Itoa(int(flow))+") failed ("+hresultString(hr)+")")
	}
	defer comRelease(&dev)

	// IMMDevice::GetId (slot 5) -> CoTaskMem LPWSTR.
	var pid unsafe.Pointer
	if hr, _, _ := comCall(dev, 5, uintptr(unsafe.Pointer(&pid))); !hrOK(hr) || pid == nil {
		return DeviceDescriptor{}, observe.New(observe.ClassAudioDevice,
			"GetId failed ("+hresultString(hr)+")")
	}
	id := windows.UTF16PtrToString((*uint16)(pid))
	pCoTaskMemFree.Call(uintptr(pid))

	return DeviceDescriptor{ID: id, Name: friendlyName(dev)}, nil
}

// DefaultCaptureDevice implements DeviceWatcher.
func (w *mmDeviceWatcher) DefaultCaptureDevice() (DeviceDescriptor, error) {
	return w.defaultEndpoint(eCapture)
}

// DefaultRenderDevice implements DeviceWatcher (D47 capture/render pair,
// queryable for the P15 AEC spike).
func (w *mmDeviceWatcher) DefaultRenderDevice() (DeviceDescriptor, error) {
	return w.defaultEndpoint(eRender)
}

// OnDefaultCaptureChanged implements DeviceWatcher.
func (w *mmDeviceWatcher) OnDefaultCaptureChanged(cb func()) (func(), error) {
	if cb == nil {
		return nil, observe.New(observe.ClassAudioDevice, "nil hotplug callback")
	}
	e := &cbEntry{fn: cb}
	w.mu.Lock()
	w.cbs = append(w.cbs, e)
	w.mu.Unlock()
	return func() {
		w.mu.Lock()
		for i, c := range w.cbs {
			if c == e {
				w.cbs = append(w.cbs[:i], w.cbs[i+1:]...)
				break
			}
		}
		w.mu.Unlock()
	}, nil
}

// Close implements DeviceWatcher: unregister + release. Idempotent.
func (w *mmDeviceWatcher) Close() error {
	w.mu.Lock()
	closed := w.closed
	w.closed = true
	reg := w.reg
	w.reg = nil
	w.cbs = nil
	w.mu.Unlock()
	if closed {
		return nil
	}
	if reg != nil {
		// UnregisterEndpointNotificationCallback (slot 7).
		comCall(reg, 7, w.this)
		comRelease(&reg)
	}
	if w.this != 0 {
		notifRegistry.Delete(w.this)
		w.this = 0
	}
	return nil
}

// notify fans an event out to the registered callbacks. Runs on a Windows
// audio-service thread: no COM, no blocking work beyond the copy.
func (w *mmDeviceWatcher) notify() {
	w.mu.Lock()
	cbs := append([]*cbEntry{}, w.cbs...)
	w.mu.Unlock()
	for _, e := range cbs {
		e.fn()
	}
}

// ---------------------------------------------------------------- notifClient

// notifVtbl is the IMMNotificationClient vtable (IUnknown + 5 methods).
type notifVtbl struct {
	queryInterface         uintptr
	addRef                 uintptr
	release                uintptr
	onDeviceStateChanged   uintptr
	onDeviceAdded          uintptr
	onDeviceRemoved        uintptr
	onDefaultDeviceChanged uintptr
	onPropertyValueChanged uintptr
}

// notifClient synthesizes the COM object: its FIRST word is the vtable
// pointer, which is the layout COM expects at `this`.
type notifClient struct {
	vt *notifVtbl
	w  *mmDeviceWatcher
}

func newNotifClient(w *mmDeviceWatcher) *notifClient {
	return &notifClient{vt: &notifVtbl{
		queryInterface:         windows.NewCallback(nClientQueryInterface),
		addRef:                 windows.NewCallback(nClientAddRef),
		release:                windows.NewCallback(nClientRelease),
		onDeviceStateChanged:   windows.NewCallback(nClientOnDeviceStateChanged),
		onDeviceAdded:          windows.NewCallback(nClientOnDeviceAdded),
		onDeviceRemoved:        windows.NewCallback(nClientOnDeviceRemoved),
		onDefaultDeviceChanged: windows.NewCallback(nClientOnDefaultDeviceChanged),
		onPropertyValueChanged: windows.NewCallback(nClientOnPropertyValueChanged),
	}, w: w}
}

// HRESULTs used by the callbacks.
const (
	hrEPointer     = 0x80004003
	hrENoInterface = 0x80004002
)

// nClientQueryInterface: the audio service never type-negotiates a
// registered notification client (registration passes the typed pointer and
// drives AddRef/Release directly), and a uintptr out-slot cannot be written
// vet-cleanly from Go. So QI declines everything - a normal COM outcome for
// any consumer that does QI first - and lifetime stays Go-managed.
func nClientQueryInterface(this, riid, out uintptr) uintptr {
	return hrENoInterface
}

func nClientAddRef(this uintptr) uintptr  { return 1 }
func nClientRelease(this uintptr) uintptr { return 1 }

// nClientOnDefaultDeviceChanged(dataFlow, role, pwstrDefaultDevice): the
// D42#2 trigger. Only capture-flow changes wake the supervisor.
func nClientOnDefaultDeviceChanged(this, flow, role, device uintptr) uintptr {
	if w := notifWatcher(this); w != nil && flow == eCapture {
		w.notify()
	}
	return hrSOK
}

// nClientOnDeviceStateChanged(pwstrDeviceId, newState): also fire when a
// device leaves the active state (unplugged/disabled mic) so the stale-
// handle path runs even without a default-change event.
const deviceStateActive = 0x1 // DEVICE_STATE_ACTIVE

func nClientOnDeviceStateChanged(this, deviceId, newState uintptr) uintptr {
	if w := notifWatcher(this); w != nil && newState != deviceStateActive {
		w.notify()
	}
	return hrSOK
}

func nClientOnDeviceAdded(this, deviceId uintptr) uintptr { return hrSOK }

func nClientOnDeviceRemoved(this, deviceId uintptr) uintptr { return hrSOK }

func nClientOnPropertyValueChanged(this, deviceId, key uintptr) uintptr { return hrSOK }

// notifWatcher recovers the watcher from the raw `this` value via the
// registry (no uintptr->pointer conversion; vet-clean).
func notifWatcher(this uintptr) *mmDeviceWatcher {
	if v, ok := notifRegistry.Load(this); ok {
		return v.(*mmDeviceWatcher)
	}
	return nil
}

// friendlyName reads PKEY_Device_FriendlyName via IPropertyStore
// (IMMDevice vtable: 4 OpenPropertyStore; IPropertyStore: 5 GetValue).
// Best-effort: failures return "" and the endpoint id is used instead.
func friendlyName(dev unsafe.Pointer) string {
	var store unsafe.Pointer
	if hr, _, _ := comCall(dev, 4, stgmRead, uintptr(unsafe.Pointer(&store))); !hrOK(hr) || store == nil {
		return ""
	}
	defer comRelease(&store)

	var pv propVariant
	if hr, _, _ := comCall(store, 5,
		uintptr(unsafe.Pointer(&pkeyDeviceFriendlyName)),
		uintptr(unsafe.Pointer(&pv))); !hrOK(hr) {
		return ""
	}
	defer pPropVariantClear.Call(uintptr(unsafe.Pointer(&pv)))
	if pv.vt != vtLPWSTR {
		return ""
	}
	p := *(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&pv), 8))
	if p == nil {
		return ""
	}
	return windows.UTF16PtrToString((*uint16)(p))
}
