//go:build windows

package audio

// [audio] input_device (the 180-census follow-up's wire): the picker and the
// open path. The census measured the key as zero production readers - the
// capture stack always opened watcher.DefaultCaptureDevice() - so the reader
// is device.go's selectCaptureDevice, consulted on the initial open and on
// every hotplug reopen. The tests run at the same seams the hotplug family
// uses (fake watcher/opener, real loop), never at the audio-data seam.

import (
	"errors"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
)

// listerWatcher wraps the hotplug family's fakeWatcher with the optional
// enumeration face the production MMDevice watcher carries.
type listerWatcher struct {
	fakeWatcher
	list []DeviceDescriptor
}

func (w *listerWatcher) CaptureDevices() ([]DeviceDescriptor, error) { return w.list, nil }

var (
	dev180kA = DeviceDescriptor{ID: "end-a", Name: "Blue Yeti"}
	dev180kB = DeviceDescriptor{ID: "end-b", Name: "Realtek Mic Array"}
	dev180kC = DeviceDescriptor{ID: "end-c", Name: "usb yeti mini"}
)

func dev180kLister() *listerWatcher {
	return &listerWatcher{
		fakeWatcher: fakeWatcher{devs: []DeviceDescriptor{dev180kA}},
		list:        []DeviceDescriptor{dev180kA, dev180kB, dev180kC},
	}
}

func Test180KInputDeviceDefaultKeepsTheDefaultEndpoint(t *testing.T) {
	w := dev180kLister()
	for _, name := range []string{"", "default", "  "} {
		dev, err := selectCaptureDevice(w, name)
		if err != nil {
			t.Fatalf("name %q: %v", name, err)
		}
		if dev.ID != dev180kA.ID {
			t.Fatalf("name %q resolved %q, want the watcher's default %q", name, dev.ID, dev180kA.ID)
		}
	}
}

func Test180KInputDeviceResolvesByNameAndByID(t *testing.T) {
	w := dev180kLister()
	cases := map[string]string{
		"Blue Yeti":            dev180kA.ID, // exact friendly name
		"  Realtek Mic Array ": dev180kB.ID, // trimmed
		"USB YETI MINI":        dev180kC.ID, // case-insensitive
		"end-b":                dev180kB.ID, // endpoint id wins even when it is not a name
	}
	for name, want := range cases {
		dev, err := selectCaptureDevice(w, name)
		if err != nil {
			t.Fatalf("name %q: %v", name, err)
		}
		if dev.ID != want {
			t.Fatalf("name %q resolved %q, want %q", name, dev.ID, want)
		}
	}
}

func Test180KInputDeviceUnknownNameRefusesAndLists(t *testing.T) {
	w := dev180kLister()
	_, err := selectCaptureDevice(w, "No Such Mic")
	if err == nil {
		t.Fatal("an unmatched input_device must be refused, never silently opened as the default")
	}
	var oe *observe.Error
	if !errors.As(err, &oe) || string(oe.Class) != string(observe.ClassAudioDevice) {
		t.Fatalf("the refusal must be a ClassAudioDevice error, got %T %v", err, err)
	}
	for _, dev := range w.list {
		if !strings.Contains(err.Error(), dev.String()) {
			t.Fatalf("the refusal must name the present endpoints so the user can copy one; %q missing from %q", dev.String(), err.Error())
		}
	}
}

func Test180KInputDeviceWatcherWithoutListerRefuses(t *testing.T) {
	// A watcher that cannot enumerate must not fall back to the default: that
	// is the "it fell back" and "it was never written" lookalike the config
	// layer refuses for every other read-only key.
	_, err := selectCaptureDevice(&fakeWatcher{devs: []DeviceDescriptor{dev180kA}}, "Blue Yeti")
	if err == nil {
		t.Fatal("a named device on a lister-less watcher must be refused")
	}
	if !strings.Contains(err.Error(), "cannot enumerate") {
		t.Fatalf("the refusal must say the watcher cannot enumerate, got: %v", err)
	}
}

func Test180KOpenPathOpensTheNamedDevice(t *testing.T) {
	// The production reader: the mic's open path (initial open, and the same
	// call every hotplug reopen makes) resolves the configured name, so the
	// opener receives the named endpoint, not the default one.
	op := &fakeOpener{}
	mic := newWASAPIMicrophoneWith(dev180kLister(), op, WithInputDevice("Realtek Mic Array"))
	if _, err := mic.openCurrent(); err != nil {
		t.Fatalf("openCurrent with a named device: %v", err)
	}
	if op.openCount() != 1 || op.opens[0].ID != dev180kB.ID {
		t.Fatalf("the open path opened %+v, want the named endpoint %q", op.opens, dev180kB.ID)
	}

	op2 := &fakeOpener{}
	mic2 := newWASAPIMicrophoneWith(dev180kLister(), op2) // no option: the pre-wire shape
	if _, err := mic2.openCurrent(); err != nil {
		t.Fatalf("openCurrent without input_device: %v", err)
	}
	if op2.opens[0].ID != dev180kA.ID {
		t.Fatalf("the default shape moved: the open hit %q, want the watcher default %q", op2.opens[0].ID, dev180kA.ID)
	}
}
