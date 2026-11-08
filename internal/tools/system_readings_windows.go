//go:build windows

package tools

import (
	"sync"
	"time"
	"unsafe"
)

// The readings system.get reports that have no home in the verb table: the wall
// clock (display only, never a timeout source) and the platform readings the
// host can supply without a new module. A function field, not a package var, so
// a test can pin the rendered text.

var (
	systemNowFn  = func() (string, error) { return time.Now().Format(time.RFC3339), nil }
	platformRead = func() map[string]any {
		return map[string]any{"battery": batteryState(), "volume": volumeReading()}
	}
)

// systemNow returns the local wall clock for the report (display only).
func systemNow() (string, error) { return systemNowFn() }

// platformReadings returns the platform readings present in this build.
func platformReadings() map[string]any { return platformRead() }

// batteryState reports AC line status and battery life percent when the machine
// has a battery. A desktop answers {"present": false} - an absent key would be
// indistinguishable from "we did not look".
type systemPowerStatus struct {
	ACLineStatus        uint8
	BatteryFlag         uint8
	BatteryLifePercent  uint8
	SystemStatusFlag    uint8
	BatteryLifeTime     uint32
	BatteryFullLifeTime uint32
}

var procGetSystemPowerStatus = kernel32.NewProc("GetSystemPowerStatus")

func batteryState() map[string]any {
	var status systemPowerStatus
	if r, _, _ := procGetSystemPowerStatus.Call(uintptr(unsafe.Pointer(&status))); r == 0 {
		return map[string]any{"present": false, "error": "GetSystemPowerStatus 调用失败"}
	}
	out := map[string]any{"present": status.BatteryFlag != 128}
	if status.BatteryLifePercent != 255 {
		out["percent"] = int(status.BatteryLifePercent)
	}
	switch status.ACLineStatus {
	case 0:
		out["ac"] = false
	case 1:
		out["ac"] = true
	default:
		out["ac"] = "unknown"
	}
	return out
}

// volumeReading is the mixer reading behind system.get's 音量. The Windows mixer
// has no plain Win32 getter; the composition root injects one through
// SystemBridge.Ops when it has a helper, and until then the key is absent rather
// than a made-up zero.
func volumeReading() any {
	realVolumeGetMu.RLock()
	fn := realVolumeGet
	realVolumeGetMu.RUnlock()
	if fn == nil {
		return nil
	}
	level, muted, err := fn()
	if err != nil {
		return map[string]any{"error": err.Error()}
	}
	return map[string]any{"level": level, "muted": muted}
}

var (
	realVolumeGetMu sync.RWMutex
	realVolumeGet   func() (int, bool, error)
)
