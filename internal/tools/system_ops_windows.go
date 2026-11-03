//go:build windows

package tools

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sys/windows"
)

// realVolume is the Windows mixer control behind system.set's 音量/静音 and
// media.control's transport verbs. Windows' Core Audio mixer (IAudioEndpointVolume)
// is a COM interface; reaching it from Go without a new module means calling the
// same MMDevice/EndpointVolume vtable entry points ole32 exposes, which this file
// does not do - instead it uses the two things the OS itself offers without COM:
//
//	the volume keys (VK_VOLUME_UP / VK_VOLUME_DOWN / VK_VOLUME_MUTE), which the
//	shell turns into one mixer step each.
//
// Why keys and not a "set to exactly N%" call: a step count is a delta, and the
// set verb wants an absolute level. Emitting (current-level) keystrokes requires
// reading the level first, which is exactly the COM call above. So the native
// backend answers "set to N%" with the honest limitation below, and an injected
// helper (realVolumeGet / SetVolumeHook) - which the composition root can supply
// from the panel's existing audio path - is what makes the absolute form work.
// Nothing here reports a level it did not read.

var (
	volumeHookMu sync.RWMutex
	volumeSet    func(action string, level int) (string, error)
)

// SetVolumeHook installs the absolute mixer writer (nil removes it). It is the
// seam a composition root uses when it has a Core Audio helper; the audit line
// names it on every call so a run that used it is distinguishable.
func SetVolumeHook(fn func(action string, level int) (string, error)) {
	volumeHookMu.Lock()
	defer volumeHookMu.Unlock()
	volumeSet = fn
}

// realVolume performs one mixer change.
func realVolume(action string, level int) (string, error) {
	volumeHookMu.RLock()
	hook := volumeSet
	volumeHookMu.RUnlock()
	if hook != nil {
		return hook(action, level)
	}
	switch action {
	case "mute", "unmute":
		if err := sendVolumeKey(vkVolumeMute); err != nil {
			return "", err
		}
		if action == "mute" {
			return "已按下静音键（切换语义：本机没有接入绝对电平读取，无法确认当前是否静音）", nil
		}
		return "已按下静音键（切换语义：本机没有接入绝对电平读取，无法确认当前是否取消静音）", nil
	case "up", "down":
		var key uint16 = vkVolumeUp
		if action == "down" {
			key = vkVolumeDown
		}
		if err := sendVolumeKey(key); err != nil {
			return "", err
		}
		return fmt.Sprintf("已按音量%s键一格", map[string]string{"up": "增", "down": "减"}[action]), nil
	default:
		return "", fmt.Errorf("%w: 绝对音量设置需要 Core Audio 电平接口，本后端未接入"+
			"（可用 SetVolumeHook 注入；按键通道只有一步一格）", ErrNotImplementedHere)
	}
}

func sendVolumeKey(vk uint16) error {
	return sendInputs([]inputUnion{
		{Type: inputKeyboard, Data: kbInput{Vk: vk}},
		{Type: inputKeyboard, Data: kbInput{Vk: vk, Flags: keyeventfKeyup}},
	})
}

// platformSystemBackend answers the SystemBridge seam's platform default.
func platformSystemBackend() (systemFuncs, error) {
	if err := ensureDPIAwareness(); err != nil {
		// A failed DPI declaration is not fatal: coordinates are then in
		// physical pixels, which the capture path handles (it copies a
		// rectangle, it does not compute from a DPI scale).
		_ = err
	}
	return windowsSystemBackend(), nil
}

var dpiOnce sync.Once

// ensureDPIAwareness declares per-monitor DPI awareness once per process so a
// captured rectangle matches the coordinates the window API reported. It is
// idempotent and never fails the call: the shcore entry point is absent on
// Windows 8.
func ensureDPIAwareness() error {
	var firstErr error
	dpiOnce.Do(func() {
		if procSetProcessDpiAwareness.Find() != nil {
			firstErr = ErrNotImplementedHere
			return
		}
		// PROCESS_PER_MONITOR_DPI_AWARE = 2
		r, _, e := procSetProcessDpiAwareness.Call(2)
		if r != 0 && e != nil {
			firstErr = e
		}
	})
	return firstErr
}

// keep the syscall package honest: ExitWindowsEx needs a handle-shaped second
// argument, which the backend passes as 0 via callBool.
var _ = windows.CurrentProcess

var _ = context.Background
