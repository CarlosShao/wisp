//go:build windows

package ball

// Global hotkeys ([hotkey]: summon / mute / cancel / panel; SPEC-03 §4).
// B1 Esc rule: the cancel binding is temporarily Esc during Confirming and
// MUST be handed back at session end - TakeEscForCancel / ReleaseEscAfterSes
// sion implement exactly that, and the tests pin both directions.

import (
	"fmt"
	"log/slog"
	"strings"

	"golang.org/x/sys/windows"
)

// Hotkey ids (wParam of WM_HOTKEY).
const (
	hkSummon = 1
	hkMute   = 2
	hkCancel = 3
	hkPanel  = 4
)

// HotkeyConfig is the parsed [hotkey] section. Empty string = disabled.
type HotkeyConfig struct {
	Summon string // e.g. "Ctrl+Alt+W"
	Mute   string // e.g. "Ctrl+Alt+M"
	Cancel string // default "Esc" (B1: temporary takeover during Confirming)
	Panel  string // e.g. "Ctrl+Alt+P"
}

// DefaultHotkeys returns the product defaults (SPEC-03 §4: cancel = Esc).
func DefaultHotkeys() HotkeyConfig {
	return HotkeyConfig{Summon: "Ctrl+Alt+W", Mute: "Ctrl+Alt+M", Cancel: "Esc", Panel: "Ctrl+Alt+P"}
}

// Accelerator is a parsed hotkey binding.
type Accelerator struct {
	Mods uint32
	VK   uint32
}

// Win32 modifier + key constants.
const (
	modAlt      = 0x0001
	modControl  = 0x0002
	modShift    = 0x0004
	modWin      = 0x0008
	modNoRepeat = 0x4000

	vkEscape = 0x1B
	vkSpace  = 0x20
)

// ParseAccelerator parses "Ctrl+Alt+W" / "Esc" / "F9" style bindings.
// Parse errors name the offending part; the caller logs and disables that
// binding (a bad hotkey string must never prevent boot).
func ParseAccelerator(s string) (Accelerator, error) {
	var acc Accelerator
	s = strings.TrimSpace(s)
	if s == "" {
		return acc, fmt.Errorf("empty hotkey")
	}
	for _, part := range strings.Split(s, "+") {
		switch strings.ToLower(strings.TrimSpace(part)) {
		case "ctrl", "control":
			acc.Mods |= modControl
		case "alt":
			acc.Mods |= modAlt
		case "shift":
			acc.Mods |= modShift
		case "win", "super":
			acc.Mods |= modWin
		case "esc", "escape":
			acc.VK = vkEscape
		case "space":
			acc.VK = vkSpace
		default:
			if vk, ok := parseVK(part); ok {
				acc.VK = vk
			} else {
				return acc, fmt.Errorf("unknown key %q in hotkey %q", part, s)
			}
		}
	}
	if acc.VK == 0 {
		return acc, fmt.Errorf("hotkey %q has no key", s)
	}
	acc.Mods |= modNoRepeat // no auto-repeat storms while held
	return acc, nil
}

// parseVK resolves named keys and single characters (letters/digits).
func parseVK(name string) (uint32, bool) {
	n := strings.TrimSpace(name)
	if len(n) == 1 {
		c := n[0]
		switch {
		case c >= 'A' && c <= 'Z':
			return uint32(c), true
		case c >= 'a' && c <= 'z':
			return uint32(c) - 'a' + 'A', true
		case c >= '0' && c <= '9':
			return uint32(c), true
		}
	}
	upper := strings.ToUpper(n)
	switch upper {
	case "F1":
		return 0x70, true
	case "F2":
		return 0x71, true
	case "F3":
		return 0x72, true
	case "F4":
		return 0x73, true
	case "F5":
		return 0x74, true
	case "F6":
		return 0x75, true
	case "F7":
		return 0x76, true
	case "F8":
		return 0x77, true
	case "F9":
		return 0x78, true
	case "F10":
		return 0x79, true
	case "F11":
		return 0x7A, true
	case "F12":
		return 0x7B, true
	case "ENTER":
		return 0x0D, true
	case "TAB":
		return 0x09, true
	case "HOME":
		return 0x24, true
	case "END":
		return 0x23, true
	case "PGUP", "PAGEUP":
		return 0x21, true
	case "PGDN", "PAGEDOWN":
		return 0x22, true
	case "UP":
		return 0x26, true
	case "DOWN":
		return 0x28, true
	case "LEFT":
		return 0x25, true
	case "RIGHT":
		return 0x27, true
	case "INSERT":
		return 0x2D, true
	case "DELETE", "DEL":
		return 0x2E, true
	}
	return 0, false
}

// registerAll binds every configured hotkey on hwnd, returning the ids that
// registered. Failure of one binding never blocks the others.
func registerAll(hwnd windows.HWND, cfg HotkeyConfig) map[uint32]Accelerator {
	ok := map[uint32]Accelerator{}
	pairs := []struct {
		id   uint32
		bind string
	}{
		{hkSummon, cfg.Summon},
		{hkMute, cfg.Mute},
		{hkCancel, cfg.Cancel},
		{hkPanel, cfg.Panel},
	}
	for _, p := range pairs {
		if p.bind == "" {
			continue // disabled by config
		}
		acc, err := ParseAccelerator(p.bind)
		if err != nil {
			slog.Warn("hotkey disabled (bad binding)", "hotkey", p.bind, "err", err)
			continue
		}
		if r, _, _ := pRegisterHotKey.Call(uintptr(hwnd), uintptr(p.id),
			uintptr(acc.Mods), uintptr(acc.VK)); r != 0 {
			ok[p.id] = acc
		} else {
			slog.Warn("hotkey registration failed (already taken?)", "hotkey", p.bind)
		}
	}
	return ok
}

func unregisterAll(hwnd windows.HWND, registered map[uint32]Accelerator) {
	for id := range registered {
		pUnregisterHotKey.Call(uintptr(hwnd), uintptr(id))
	}
}

// takeEsc binds VK_ESCAPE as the cancel hotkey (B1 Confirming takeover). The
// configured cancel binding stays remembered in the ball for release.
func takeEsc(hwnd windows.HWND) bool {
	pUnregisterHotKey.Call(uintptr(hwnd), hkCancel)
	r, _, _ := pRegisterHotKey.Call(uintptr(hwnd), hkCancel, modNoRepeat, vkEscape)
	return r != 0
}

// releaseEsc re-registers the configured cancel binding (B1: the Esc key
// must be handed back when the session ends).
func releaseEsc(hwnd windows.HWND, bind string) bool {
	pUnregisterHotKey.Call(uintptr(hwnd), hkCancel)
	if bind == "" {
		return false
	}
	acc, err := ParseAccelerator(bind)
	if err != nil {
		return false
	}
	r, _, _ := pRegisterHotKey.Call(uintptr(hwnd), hkCancel, uintptr(acc.Mods), uintptr(acc.VK))
	return r != 0
}
