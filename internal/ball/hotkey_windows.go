//go:build windows

package ball

// Global hotkeys ([hotkey]: summon / mute / cancel / panel; SPEC-03 §4).
// B1 Esc rule: the cancel binding is temporarily Esc during Confirming and
// MUST be handed back at session end - TakeEscForCancel / ReleaseEscAfterSes
// sion implement exactly that, and the tests pin both directions.

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

// Hotkey ids (wParam of WM_HOTKEY).
const (
	hkSummon = 1
	hkMute   = 2
	hkCancel = 3
	hkPanel  = 4
)

// hkNames is the id -> user-visible name table. The order is the registration
// order, so a report reads summon, mute, cancel, panel.
var hkNames = []struct {
	id   uint32
	name string
}{
	{hkSummon, "summon"},
	{hkMute, "mute"},
	{hkCancel, "cancel"},
	{hkPanel, "panel"},
}

// HotkeyConfig is the parsed [hotkey] section. Empty string = disabled.
type HotkeyConfig struct {
	Summon string // e.g. "Ctrl+Alt+Q"
	Mute   string // e.g. "Ctrl+Alt+M"
	Cancel string // default "Esc" (B1: temporary takeover during Confirming)
	Panel  string // e.g. "Ctrl+Alt+P"
}

// AltSummonSpace is the documented alternative summon binding:
// "Ctrl+Alt+Space" registers fine on this machine and stays user-selectable,
// but it is NOT the default because several Chinese IMEs (Sogou/WeChat-style
// shift-space and ctrl-space commits) grab Ctrl+Alt+Space before the app can.
const AltSummonSpace = "Ctrl+Alt+Space"

// DefaultHotkeys returns the product defaults (SPEC-03 §4: cancel = Esc).
// Summon is Ctrl+Alt+Q by owner ruling R10 (2026-09-20): a sweep of 84
// candidate combinations on the owner's machine found Ctrl+Alt+W occupied by
// a third-party app (the only occupied one), so the previous default could
// never register there. Ctrl+Alt+Space (AltSummonSpace) is the supported
// alternative, deliberately not the default - see that constant's comment.
func DefaultHotkeys() HotkeyConfig {
	return HotkeyConfig{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Esc", Panel: "Ctrl+Alt+P"}
}

// ApplyHotkeyDefaults fills the empty fields of a config-sourced binding set
// with the product defaults. [hotkey] in config.toml distinguishes nothing
// between "unset" and "explicitly disabled" - both are an empty string - so a
// host that maps config.Hotkey straight in would silently end up with NO
// summon key at all on a fresh install (summon/mute/panel default to empty in
// the schema). Hosts that read [hotkey] from the config file pass the mapping
// through here; a config with summon = "none"/"off" is the host's business
// (it arrives as "" and re-enables the default, which is the safe direction:
// a live hotkey the user can see, never a dead silent one).
func ApplyHotkeyDefaults(cfg HotkeyConfig) HotkeyConfig {
	d := DefaultHotkeys()
	if cfg.Summon == "" {
		cfg.Summon = d.Summon
	}
	if cfg.Mute == "" {
		cfg.Mute = d.Mute
	}
	if cfg.Cancel == "" {
		cfg.Cancel = d.Cancel
	}
	if cfg.Panel == "" {
		cfg.Panel = d.Panel
	}
	return cfg
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

// HotkeyStatus is the outcome of ONE binding. The two failure families are
// deliberately distinct (ticket 64 A1b): "somebody else already owns this
// key" is actionable advice for the user, "we never even tried" is a config
// bug, and the old code collapsed both into one `slog.Warn` with a question
// mark in it ("already taken?"), which told the user nothing.
type HotkeyStatus uint8

const (
	// HotkeyLive: RegisterHotKey succeeded, the binding is ours and presses
	// arrive as WM_HOTKEY.
	HotkeyLive HotkeyStatus = iota
	// HotkeyDisabled: NOT attempted - the binding is empty in config.
	HotkeyDisabled
	// HotkeyUnparsable: NOT attempted - the binding string is not a valid
	// accelerator (ParseAccelerator refused it).
	HotkeyUnparsable
	// HotkeyTaken: attempted and refused by Win32 with
	// ERROR_HOTKEY_ALREADY_REGISTERED (1409), i.e. another process (or
	// another window) owns that exact combination.
	HotkeyTaken
	// HotkeyError: attempted and refused for any other Win32 reason.
	HotkeyError
)

// String names the status for logs and for the user-visible problem lines.
func (s HotkeyStatus) String() string {
	switch s {
	case HotkeyLive:
		return "live"
	case HotkeyDisabled:
		return "disabled (unset in config)"
	case HotkeyUnparsable:
		return "not attempted (binding cannot be parsed)"
	case HotkeyTaken:
		return "occupied by another app"
	case HotkeyError:
		return "registration failed"
	}
	return "unknown"
}

// Attempted reports whether a RegisterHotKey call was made for this binding.
// The split the ticket asks for is exactly this bit: HotkeyTaken/HotkeyError
// were attempted and refused; HotkeyDisabled/HotkeyUnparsable were never
// attempted at all.
func (s HotkeyStatus) Attempted() bool {
	return s == HotkeyTaken || s == HotkeyError
}

// HotkeyBinding is the per-binding line of a HotkeyReport.
type HotkeyBinding struct {
	Name    string // summon|mute|cancel|panel
	ID      uint32 // wParam of WM_HOTKEY
	Binding string // the configured string, verbatim ("" when unset)
	Status  HotkeyStatus
	Acc     Accelerator // parsed accelerator (zero when not attempted)
	Err     error       // Win32 error for HotkeyTaken/HotkeyError, parse error for HotkeyUnparsable
}

// HotkeyReport is the outcome of one registration pass over all four ids.
type HotkeyReport struct {
	bindings []HotkeyBinding
}

// Bindings lists the four ids in registration order.
func (r HotkeyReport) Bindings() []HotkeyBinding { return slices.Clone(r.bindings) }

// Live is the id -> accelerator map of the bindings that are OURS right now
// (the "registration set" the ticket says was never asserted anywhere).
func (r HotkeyReport) Live() map[uint32]Accelerator {
	m := map[uint32]Accelerator{}
	for _, b := range r.bindings {
		if b.Status == HotkeyLive {
			m[b.ID] = b.Acc
		}
	}
	return m
}

// Binding returns the recorded outcome for one hotkey id.
func (r HotkeyReport) Binding(id uint32) (HotkeyBinding, bool) {
	for _, b := range r.bindings {
		if b.ID == id {
			return b, true
		}
	}
	return HotkeyBinding{}, false
}

// IsLive reports whether id registered.
func (r HotkeyReport) IsLive(id uint32) bool {
	b, ok := r.Binding(id)
	return ok && b.Status == HotkeyLive
}

// AllLive reports every configured binding registered (disabled ones count
// as satisfied: the user asked for no key).
func (r HotkeyReport) AllLive() bool {
	for _, b := range r.bindings {
		if b.Status != HotkeyLive && b.Status != HotkeyDisabled {
			return false
		}
	}
	return true
}

// Problems renders the user-visible lines: one per binding that is not live,
// naming the key and saying which of the two families it is. No emoji, no
// question marks about which case it is - the case IS known here.
func (r HotkeyReport) Problems() []string {
	out := make([]string, 0, len(r.bindings))
	for _, b := range r.bindings {
		switch b.Status {
		case HotkeyLive, HotkeyDisabled:
			continue
		case HotkeyUnparsable:
			out = append(out, fmt.Sprintf(
				"hotkey %s = %q was not registered: the binding is not a valid key combination (%v)",
				b.Name, b.Binding, b.Err))
		case HotkeyTaken:
			out = append(out, fmt.Sprintf(
				"hotkey %s = %q is occupied by another program and was NOT registered; "+
					"pressing it will do nothing until you pick a free combination in [hotkey]",
				b.Name, b.Binding))
		default:
			out = append(out, fmt.Sprintf(
				"hotkey %s = %q was not registered: %v", b.Name, b.Binding, b.Err))
		}
	}
	return out
}

// registerFn performs one Win32 registration; it is the seam that lets the
// status classification be tested without a desktop (ticket 64 A1b).
type registerFn func(id uint32, acc Accelerator) error

// hotkeyRegisterer is the Win32 RegisterHotKey for a fixed window.
func hotkeyRegisterer(hwnd windows.HWND) registerFn {
	return func(id uint32, acc Accelerator) error {
		r, _, err := pRegisterHotKey.Call(uintptr(hwnd), uintptr(id),
			uintptr(acc.Mods), uintptr(acc.VK))
		if r != 0 {
			return nil
		}
		if errno, ok := err.(syscall.Errno); ok && errno != 0 {
			return errno
		}
		return syscall.EINVAL // RegisterHotKey returned FALSE with no error info
	}
}

// unregisterFn performs one Win32 unregistration (test seam).
type unregisterFn func(id uint32)

// hotkeyUnregisterer is the Win32 UnregisterHotKey for a fixed window.
func hotkeyUnregisterer(hwnd windows.HWND) unregisterFn {
	return func(id uint32) { pUnregisterHotKey.Call(uintptr(hwnd), uintptr(id)) }
}

// registerAll binds every configured hotkey on hwnd and returns the outcome
// per binding. Failure of one binding never blocks the others, and no failure
// is silent: each one carries a status the caller can surface to the user.
func registerAll(hwnd windows.HWND, cfg HotkeyConfig) HotkeyReport {
	return registerAllWith(hwnd, cfg, hotkeyRegisterer(hwnd))
}

// registerAllWith is registerAll over an injectable registrar.
func registerAllWith(_ windows.HWND, cfg HotkeyConfig, reg registerFn) HotkeyReport {
	rep := HotkeyReport{bindings: make([]HotkeyBinding, 0, len(hkNames))}
	bindings := []string{cfg.Summon, cfg.Mute, cfg.Cancel, cfg.Panel}
	for i, p := range hkNames {
		b := HotkeyBinding{Name: p.name, ID: p.id, Binding: bindings[i]}
		switch {
		case b.Binding == "":
			b.Status = HotkeyDisabled // disabled by config: never attempted
			slog.Info("hotkey disabled (unset in config)", "hotkey", b.Name)
		default:
			acc, err := ParseAccelerator(b.Binding)
			if err != nil {
				b.Status = HotkeyUnparsable
				b.Err = err
				slog.Error("hotkey not attempted (bad binding)", "hotkey", b.Name,
					"binding", b.Binding, "err", err)
			} else if rerr := reg(p.id, acc); rerr != nil {
				b.Acc, b.Err = acc, rerr
				if errors.Is(rerr, windows.ERROR_HOTKEY_ALREADY_REGISTERED) {
					b.Status = HotkeyTaken
					slog.Error("hotkey occupied by another program, not registered",
						"hotkey", b.Name, "binding", b.Binding, "err", rerr)
				} else {
					b.Status = HotkeyError
					slog.Error("hotkey registration failed",
						"hotkey", b.Name, "binding", b.Binding, "err", rerr)
				}
			} else {
				b.Status = HotkeyLive
				b.Acc = acc
			}
		}
		rep.bindings = append(rep.bindings, b)
	}
	return rep
}

// unregisterAll drops the four known hotkey slots on this window, not just
// the ones the caller still tracks: a slot we lost track of (a rebind that
// half-failed, a takeover id) would otherwise come back on the NEXT rebind as
// a 1409 that reads "somebody else owns it" while the somebody is us.
// Unregistering an id we do not hold is a no-op Win32 refuses silently.
func unregisterAll(hwnd windows.HWND) {
	unregisterAllWith(hotkeyUnregisterer(hwnd))
}

func unregisterAllWith(unreg unregisterFn) {
	for _, p := range hkNames {
		unreg(p.id)
	}
}

// takeEsc binds VK_ESCAPE as the cancel hotkey (B1 Confirming takeover). The
// configured cancel binding stays remembered in the ball for release.
func takeEsc(hwnd windows.HWND) bool {
	pUnregisterHotKey.Call(uintptr(hwnd), hkCancel)
	err := hotkeyRegisterer(hwnd)(hkCancel, Accelerator{Mods: modNoRepeat, VK: vkEscape})
	if err != nil {
		slog.Error("Esc cancel takeover refused by Win32; Confirming will have no cancel key",
			"hotkey", "cancel", "binding", "Esc", "err", err)
		return false
	}
	return true
}

// releaseEsc re-registers the configured cancel binding (B1: the Esc key
// must be handed back when the session ends).
// releaseEsc hands Esc back (B1): the takeover binding is ALWAYS dropped -
// even when the configured binding cannot be re-registered (taken by
// another app) - and true is returned only when the original binding is
// live again. escTakenOver clears either way: Esc is no longer ours.
func releaseEsc(hwnd windows.HWND, bind string) bool {
	pUnregisterHotKey.Call(uintptr(hwnd), hkCancel)
	if bind == "" {
		return false
	}
	acc, err := ParseAccelerator(bind)
	if err != nil {
		slog.Error("cancel hotkey binding unparsable; Esc returned but unbound", "hotkey", bind, "err", err)
		return false
	}
	if rerr := hotkeyRegisterer(hwnd)(hkCancel, acc); rerr != nil {
		// The same two families as registerAll: this is where "my cancel key
		// stopped working after the session" would otherwise be unexplainable.
		slog.Error("cancel hotkey re-registration failed; Esc returned but unbound",
			"hotkey", "cancel", "binding", bind,
			"occupied", errors.Is(rerr, windows.ERROR_HOTKEY_ALREADY_REGISTERED), "err", rerr)
		return false
	}
	return true
}
