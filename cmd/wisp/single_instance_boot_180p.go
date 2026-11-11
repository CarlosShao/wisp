package main

// The deferred-disposition 20261110 #1 wiring: [app] single_instance, the one
// config key the composition root owes the boot. The seam is the one
// proc.WithLayout was left open for (internal/proc/boot_windows.go:81-85 names
// it verbatim: the portable-mode override was ticket 06, and "configuration
// will use the same seam") - Boot decides the mutex off rt.Layout.MutexEnabled,
// so a composition root that reads config BEFORE proc.Boot can turn the lock
// off by handing in a layout whose mutex names are cleared.
//
// THE ORDER IS THE WHOLE POINT. proc.Boot both resolves the layout and takes
// the mutex inside one call, so the config read cannot happen after it: a
// second launch would already have been refused by the time anything read the
// key. This helper therefore does what D36's restart tier expects - it reads
// [app] once, at the boot that owns the next process start (the same tier the
// autostart seed and the interface language ride) - and its layout answer is
// what proc.Boot is handed. single_instance is registered restart-tier in
// config.TierRegistry (internal/config/tiers.go:62), so nothing here re-reads
// it mid-run or promises a live flip: the key decides the NEXT boot, which is
// the tier's own contract, not a limitation invented here.
//
// WHAT false BUYS, AND WHAT IT COSTS (both said, neither hidden). With
// MutexEnabled=false, proc.Boot builds no mutex and rt.Instance stays nil: a
// second launch is NOT refused and NOT activated - it boots a full second
// resident. D42#7's "activate the running instance" semantics are off because
// the user's explicit key turned the feature off, which is the same shape the
// test env has lived in since ticket 06 (no mutex, parallel boots allowed).
// The boot report carries both facts: the summary line prints
// single-instance=false off rt.Instance==nil, and runResident prints a second,
// config-sourced line naming the key and its file so "off" never reads as
// "broken".
//
// THE FAIL DIRECTION IS THE DEFAULT. Any inability to read the key - no
// config.toml yet (first run), an unreadable file, a parse error - returns the
// untouched DefaultLayout: the mutex stays ON, which is the schema default
// (single_instance default:"true", internal/config/schema.go:158) and the
// fail-closed side (a broken config file must not silently dissolve the
// single-instance guarantee it cannot even name).

import (
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/proc"
)

// singleInstanceBootLayout180P is the production form: resolve the env's own
// DefaultLayout, read [app] off the config.toml that data dir owns, answer the
// boot layout.
func singleInstanceBootLayout180P(env buildinfo.Env) (proc.Layout, bool, string) {
	layout, err := proc.DefaultLayout(env)
	if err != nil {
		// Boot will refuse the same env with its own sentence; nothing here
		// can name a data dir to read config from anyway.
		return proc.Layout{}, false, ""
	}
	return singleInstanceBootLayoutRead180P(layout, filepath.Join(layout.DataDir, configFileName))
}

// singleInstanceBootLayoutRead180P is the pure form the tests drive: a layout
// (any env's, so a mutex-carrying prod layout can prove the flip without
// touching a real user data dir) plus the config file to consult. It answers:
// the layout to boot with, whether the KEY (not the env fork) turned the mutex
// off (the boot line prints only on this), and the config file the decision
// was read from ("" when the key was not consulted, so no file is implied).
func singleInstanceBootLayoutRead180P(layout proc.Layout, cfgPath string) (proc.Layout, bool, string) {
	c, _, lerr := config.LoadFile(cfgPath, nil)
	if lerr != nil || c == nil || c.App.SingleInstance {
		// Missing / unreadable / key true (the schema default): the mutex
		// stays exactly as the env fork built it.
		return layout, false, ""
	}
	off := layout
	off.MutexEnabled = false
	// The names go with the flag (proc's own invariant: MutexName/ActivateEventName
	// are "" when MutexEnabled is false, internal/proc/envfork.go:20-21), so a
	// layout that says "no mutex" cannot still carry a mutex NAME somebody
	// else might acquire with.
	off.MutexName = ""
	off.ActivateEventName = ""
	return off, true, cfgPath
}
