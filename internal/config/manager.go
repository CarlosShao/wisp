package config

import (
	"log/slog"
	"os"
	"reflect"
	"slices"
	"sync"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/risk"
)

// Reload manager (SPEC-03 sec 4.2/4.3): owns the in-memory config, polls
// mtime+size (no fsnotify, no extra goroutine - the watchdog tick calls
// CheckAndReload), and applies changes per the D36 three-tier semantics.
//
// Lock analysis: one mutex guards cur + seen stat. Callbacks run OUTSIDE the
// lock so they may call Config() / trigger subsystem reloads. cur is never
// handed out directly - Config() returns a deep copy (the sanctioned
// hand-out, see deepCopyConfig).

// Manager watches config.toml and applies changes per tier.
type Manager struct {
	// ConfirmLocked is the L2 re-confirmation hook for locked-section
	// loosening (D36 rule 1). It receives the section name and the dotted
	// key paths that loosen it. nil = deny (fail-closed).
	ConfirmLocked func(section string, keys []string) bool
	// OnReload fires when reload-tier sections changed (e.g. voice model
	// swap: the owning subsystem reloads models). Values are already applied.
	OnReload func(sections []string)
	// OnRestartPending fires when restart-tier sections changed: the old
	// values stay active until process restart; the app surfaces
	// "restart required" with these section names.
	OnRestartPending func(sections []string)

	mu        sync.Mutex
	path      string
	res       SecretResolver
	cur       *Config
	resolved  *Resolved
	seenMtime time.Time
	seenSize  int64
	seenValid bool
}

// LockedDecision records the outcome for one locked section in a reload.
type LockedDecision struct {
	// Section is risk|fs|net|plugins.
	Section string
	// Direction is the strongest detected posture (loosen wins over tighten
	// over neutral).
	Direction Direction
	// Approved is the hook verdict for loosening (true for tighten/neutral,
	// which need no confirmation).
	Approved bool
	// Keys are the changed key paths (loosening plus tightening).
	Keys []string
}

// Report describes one applied reload. A nil Report from CheckAndReload
// means "nothing changed".
type Report struct {
	// Hot lists sections hot-applied.
	Hot []string
	// Reload lists reload-tier sections (applied + event emitted).
	Reload []string
	// Restart lists restart-tier sections whose new values are deferred to
	// the next process start.
	Restart []string
	// Locked carries one decision per changed locked section.
	Locked []LockedDecision
}

// NewManager loads the initial config from path (first-run file creation is
// the caller's concern) and starts watching it.
func NewManager(path string, res SecretResolver) (*Manager, error) {
	c, resolved, err := LoadFile(path, res)
	if err != nil {
		return nil, err
	}
	m := &Manager{path: path, res: res, cur: c, resolved: resolved}
	if st, err := os.Stat(path); err == nil {
		m.seenMtime, m.seenSize, m.seenValid = st.ModTime(), st.Size(), true
	}
	return m, nil
}

// Config returns a deep copy of the current effective config (restart-tier
// pending changes are not reflected - they are not active yet).
func (m *Manager) Config() *Config {
	m.mu.Lock()
	defer m.mu.Unlock()
	return deepCopyConfig(m.cur)
}

// Resolved returns a copy of the secrets resolved for the current config
// (plaintext lives here only, never in Config; SPEC-03 sec 3.1).
func (m *Manager) Resolved() *Resolved {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.resolved == nil {
		return &Resolved{ProviderKeys: map[string]string{}}
	}
	pk := make(map[string]string, len(m.resolved.ProviderKeys))
	for k, v := range m.resolved.ProviderKeys {
		pk[k] = v
	}
	return &Resolved{ProviderKeys: pk, RealtimeKey: m.resolved.RealtimeKey}
}

// CheckAndReload is the idempotent poll (SPEC-03 sec 4.3, watchdog 1s tick):
// it stats the file and does nothing when mtime+size are unchanged. On
// change it re-runs the full load pipeline; a failed load keeps the current
// config and returns the error. Callbacks fire outside the lock.
func (m *Manager) CheckAndReload() (*Report, error) {
	st, err := os.Stat(m.path)
	if err != nil {
		return nil, observe.Wrap(observe.ClassConfig, err,
			"config.toml stat failed; keeping current config")
	}

	m.mu.Lock()
	if m.seenValid && st.ModTime().Equal(m.seenMtime) && st.Size() == m.seenSize {
		m.mu.Unlock()
		return nil, nil
	}
	fresh, resolved, err := LoadFile(m.path, m.res)
	if err != nil {
		// Adopt the stat even on failure: a newer mtime+size only appears
		// after the next write, so this avoids re-erroring every tick on a
		// file that is not going to fix itself.
		m.seenMtime, m.seenSize, m.seenValid = st.ModTime(), st.Size(), true
		m.mu.Unlock()
		return nil, err
	}
	rep := m.apply(fresh)
	m.resolved = resolved
	m.seenMtime, m.seenSize, m.seenValid = st.ModTime(), st.Size(), true
	m.mu.Unlock()

	if len(rep.Reload) > 0 && m.OnReload != nil {
		m.OnReload(rep.Reload)
	}
	if len(rep.Restart) > 0 && m.OnRestartPending != nil {
		m.OnRestartPending(rep.Restart)
	}
	return rep, nil
}

// apply merges fresh into cur per tier, mutating cur into the new effective
// config. Caller holds m.mu.
func (m *Manager) apply(fresh *Config) *Report {
	cur := m.cur
	rep := &Report{}

	m.applyLocked("risk", fresh, rep,
		func() bool { return reflect.DeepEqual(cur.Risk, fresh.Risk) },
		func() (loosen, tighten []string) { return riskDirection(&cur.Risk, &fresh.Risk) },
		func(apply bool) {
			if apply {
				cur.Risk = fresh.Risk
			}
		})

	m.applyLocked("fs", fresh, rep,
		func() bool { return reflect.DeepEqual(cur.FS, fresh.FS) },
		func() (loosen, tighten []string) { return fsDirection(&cur.FS, &fresh.FS) },
		func(apply bool) {
			if apply {
				cur.FS = fresh.FS
			}
		})

	m.applyLocked("net", fresh, rep,
		func() bool { return reflect.DeepEqual(cur.Net, fresh.Net) },
		func() (loosen, tighten []string) { return netDirection(&cur.Net, &fresh.Net) },
		func(apply bool) {
			if apply {
				cur.Net = fresh.Net
			}
		})

	m.applyLocked("plugins", fresh, rep,
		func() bool { return reflect.DeepEqual(cur.Plugins, fresh.Plugins) },
		func() (loosen, tighten []string) { return pluginsDirection(&cur.Plugins, &fresh.Plugins) },
		func(apply bool) {
			if apply {
				cur.Plugins = fresh.Plugins
			}
		})

	m.applyApp(fresh, rep)
	m.applyVoice(fresh, rep)

	// Everything else is hot-tier: apply wholesale on change.
	rest := []struct {
		name       string
		oldr, newr any
		set        func()
	}{
		{"ball", &cur.Ball, &fresh.Ball, func() { cur.Ball = fresh.Ball }},
		{"hotkey", &cur.Hotkey, &fresh.Hotkey, func() { cur.Hotkey = fresh.Hotkey }},
		{"session", &cur.Session, &fresh.Session, func() { cur.Session = fresh.Session }},
		{"audio", &cur.Audio, &fresh.Audio, func() { cur.Audio = fresh.Audio }},
		{"llm", &cur.LLM, &fresh.LLM, func() { cur.LLM = fresh.LLM }},
		{"agent", &cur.Agent, &fresh.Agent, func() { cur.Agent = fresh.Agent }},
		{"privacy", &cur.Privacy, &fresh.Privacy, func() { cur.Privacy = fresh.Privacy }},
		{"memory", &cur.Memory, &fresh.Memory, func() { cur.Memory = fresh.Memory }},
		{"panel", &cur.Panel, &fresh.Panel, func() { cur.Panel = fresh.Panel }},
		{"cost", &cur.Cost, &fresh.Cost, func() { cur.Cost = fresh.Cost }},
		{"models", &cur.Models, &fresh.Models, func() { cur.Models = fresh.Models }},
		{"observe", &cur.Observe, &fresh.Observe, func() { cur.Observe = fresh.Observe }},
	}
	for _, s := range rest {
		if !reflect.DeepEqual(s.oldr, s.newr) {
			s.set()
			rep.Hot = append(rep.Hot, s.name)
		}
	}
	return rep
}

// applyLocked handles one locked section: unchanged -> skip; any loosening
// -> L2 hook (nil denies); rejected -> keep old values; approved or
// tighten/neutral -> apply. All directions are logged.
func (m *Manager) applyLocked(section string, fresh *Config, rep *Report,
	unchanged func() bool, direction func() (loosen, tighten []string), applyTo func(apply bool),
) {
	if unchanged() {
		return
	}
	loosen, tighten := direction()
	keys := append(slices.Clone(loosen), tighten...)
	dec := LockedDecision{Section: section, Keys: keys}
	switch {
	case len(loosen) > 0:
		dec.Direction = DirLoosen
		approved := m.ConfirmLocked != nil && m.ConfirmLocked(section, loosen)
		dec.Approved = approved
		applyTo(approved)
		if approved {
			slog.Warn("config: locked loosening approved via L2 re-confirmation (D36 rule 1)",
				"section", section, "keys", loosen)
		} else {
			slog.Warn("config: locked loosening rejected; keeping previous values",
				"section", section, "keys", loosen)
		}
	case len(tighten) > 0:
		dec.Direction = DirTighten
		dec.Approved = true
		applyTo(true)
		slog.Info("config: locked section tightened, hot-applied",
			"section", section, "keys", tighten)
	default:
		dec.Direction = DirNeutral
		dec.Approved = true
		applyTo(true)
		slog.Info("config: locked section updated (neutral), hot-applied", "section", section)
	}
	rep.Locked = append(rep.Locked, dec)
}

// applyApp implements the [app] split: theme is the one hot key, the rest is
// restart-tier (kept at old values until process restart, D36).
func (m *Manager) applyApp(fresh *Config, rep *Report) {
	cur := m.cur
	if reflect.DeepEqual(cur.App, fresh.App) {
		return
	}
	themeChanged := cur.App.Theme != fresh.App.Theme
	restartChanged := cur.App.Language != fresh.App.Language ||
		cur.App.Autostart != fresh.App.Autostart ||
		cur.App.SingleInstance != fresh.App.SingleInstance
	if themeChanged {
		cur.App.Theme = fresh.App.Theme
		rep.Hot = append(rep.Hot, "app")
	}
	if restartChanged {
		// keep cur.App.Language/Autostart/SingleInstance at old values
		rep.Restart = append(rep.Restart, "app")
	}
}

// applyVoice implements the [voice] split (SPEC-03 sec 3): model-pipeline
// keys are reload-tier (apply + event), tuning keys are hot.
func (m *Manager) applyVoice(fresh *Config, rep *Report) {
	cur := m.cur
	if reflect.DeepEqual(cur.Voice, fresh.Voice) {
		return
	}
	reload, hot := false, false
	fv, cv := &fresh.Voice, &cur.Voice

	// reload-tier: anything reshaping the audio/model pipeline.
	if cv.Enabled != fv.Enabled || cv.WakeWord.Enabled != fv.WakeWord.Enabled ||
		!slices.Equal(cv.WakeWord.Keywords, fv.WakeWord.Keywords) ||
		cv.ASR != fv.ASR ||
		cv.TTS.Provider != fv.TTS.Provider || cv.TTS.Voice != fv.TTS.Voice ||
		cv.ConversationMode != fv.ConversationMode || cv.AEC != fv.AEC ||
		cv.Realtime != fv.Realtime ||
		!slices.Equal(cv.CloudASRChain, fv.CloudASRChain) ||
		!slices.Equal(cv.CloudTTSChain, fv.CloudTTSChain) {
		reload = true
	}
	cv.Enabled = fv.Enabled
	cv.WakeWord.Enabled = fv.WakeWord.Enabled
	cv.WakeWord.Keywords = slices.Clone(fv.WakeWord.Keywords)
	cv.ASR = fv.ASR
	cv.TTS.Provider = fv.TTS.Provider
	cv.TTS.Voice = fv.TTS.Voice
	cv.ConversationMode = fv.ConversationMode
	cv.AEC = fv.AEC
	cv.Realtime = fv.Realtime
	cv.CloudASRChain = slices.Clone(fv.CloudASRChain)
	cv.CloudTTSChain = slices.Clone(fv.CloudTTSChain)

	// hot-tier: thresholds, veto words, speed, punctuation (D36).
	if !slices.Equal(cv.WakeWord.Thresholds, fv.WakeWord.Thresholds) ||
		!slices.Equal(cv.WakeWord.VetoWords, fv.WakeWord.VetoWords) ||
		cv.TTS.Speed != fv.TTS.Speed || cv.Punctuation != fv.Punctuation {
		hot = true
	}
	cv.WakeWord.Thresholds = slices.Clone(fv.WakeWord.Thresholds)
	cv.WakeWord.VetoWords = slices.Clone(fv.WakeWord.VetoWords)
	cv.TTS.Speed = fv.TTS.Speed
	cv.Punctuation = fv.Punctuation

	if reload {
		rep.Reload = append(rep.Reload, "voice")
	}
	if hot {
		rep.Hot = append(rep.Hot, "voice")
	}
}

// riskDirection returns the loosening and tightening key paths between old
// and new [risk] (raising l1_window_sec auto-approves more = loosening).
func riskDirection(old, new *RiskSection) (loosen, tighten []string) {
	if !old.ShellEnabled && new.ShellEnabled {
		loosen = append(loosen, "risk.shell_enabled")
	} else if old.ShellEnabled && !new.ShellEnabled {
		tighten = append(tighten, "risk.shell_enabled")
	}
	if !old.AllowShellString && new.AllowShellString {
		loosen = append(loosen, "risk.allow_shell_string")
	} else if old.AllowShellString && !new.AllowShellString {
		tighten = append(tighten, "risk.allow_shell_string")
	}
	if new.L1WindowSec > old.L1WindowSec {
		loosen = append(loosen, "risk.l1_window_sec")
	} else if new.L1WindowSec < old.L1WindowSec {
		tighten = append(tighten, "risk.l1_window_sec")
	}
	// risk.permission_mode (ticket 90, R20/M1): the mode ordering IS the
	// ask-nothing ordering, so a higher rank = asks about less = loosening, and
	// D36 rule 1 puts it behind the same L2 re-confirmation as the other
	// locked keys. An unparseable value parses to the strictest mode on both
	// sides (risk.ParseMode never guesses), so a corrupt mode can only ever
	// register as a tightening or as neutral here, never as a loosening.
	oldMode, _ := risk.ParseMode(old.PermissionMode)
	newMode, _ := risk.ParseMode(new.PermissionMode)
	switch {
	case newMode > oldMode:
		loosen = append(loosen, "risk.permission_mode")
	case newMode < oldMode:
		tighten = append(tighten, "risk.permission_mode")
	}
	loosen, tighten = setDirection(loosen, tighten, "risk.shell_allowlist", old.ShellAllowlist, new.ShellAllowlist)
	return setDirection(loosen, tighten, "risk.blacklist_overrides", old.BlacklistOverrides, new.BlacklistOverrides)
}

// fsDirection for [fs]: any addition to the writable roots or reparse
// exceptions is loosening.
func fsDirection(old, new *FSSection) (loosen, tighten []string) {
	loosen, tighten = setDirection(loosen, tighten, "fs.allowed_dirs", old.AllowedDirs, new.AllowedDirs)
	loosen, tighten = setDirection(loosen, tighten, "fs.reparse_point_exceptions", old.ReparsePointExceptions, new.ReparsePointExceptions)
	if !old.DeleteEnabled && new.DeleteEnabled {
		loosen = append(loosen, "fs.delete_enabled")
	} else if old.DeleteEnabled && !new.DeleteEnabled {
		tighten = append(tighten, "fs.delete_enabled")
	}
	return loosen, tighten
}

// netDirection for [net]: moving away from proxy "none" opens a new egress
// path = loosening; disabling private-range blocking = loosening.
func netDirection(old, new *NetSection) (loosen, tighten []string) {
	loosen, tighten = setDirection(loosen, tighten, "net.allowlist", old.Allowlist, new.Allowlist)
	if old.BlockPrivateRanges && !new.BlockPrivateRanges {
		loosen = append(loosen, "net.block_private_ranges")
	} else if !old.BlockPrivateRanges && new.BlockPrivateRanges {
		tighten = append(tighten, "net.block_private_ranges")
	}
	oldOpen := old.Proxy.Mode != ProxyNone
	newOpen := new.Proxy.Mode != ProxyNone
	switch {
	case !oldOpen && newOpen:
		loosen = append(loosen, "net.proxy.mode")
	case oldOpen && !newOpen:
		tighten = append(tighten, "net.proxy.mode")
	}
	return loosen, tighten
}

// pluginsDirection for [plugins]: enabling tier-2, adding plugins or grants
// is loosening; removing any of them is tightening.
func pluginsDirection(old, new *PluginsSection) (loosen, tighten []string) {
	if !old.Tier2Enabled && new.Tier2Enabled {
		loosen = append(loosen, "plugins.tier2_enabled")
	} else if old.Tier2Enabled && !new.Tier2Enabled {
		tighten = append(tighten, "plugins.tier2_enabled")
	}
	for id := range new.Entries {
		if _, ok := old.Entries[id]; !ok {
			loosen = append(loosen, "plugins."+id)
			continue
		}
		o, n := old.Entries[id], new.Entries[id]
		if !o.Enabled && n.Enabled {
			loosen = append(loosen, "plugins."+id+".enabled")
		} else if o.Enabled && !n.Enabled {
			tighten = append(tighten, "plugins."+id+".enabled")
		}
		loosen, tighten = setDirection(loosen, tighten, "plugins."+id+".capabilities", o.Capabilities, n.Capabilities)
		loosen, tighten = setDirection(loosen, tighten, "plugins."+id+".net_allowlist", o.NetAllowlist, n.NetAllowlist)
		loosen, tighten = setDirection(loosen, tighten, "plugins."+id+".host_api", o.HostAPI, n.HostAPI)
	}
	for id := range old.Entries {
		if _, ok := new.Entries[id]; !ok {
			tighten = append(tighten, "plugins."+id)
		}
	}
	return loosen, tighten
}

// setDirection classifies list changes: additions are loosening, removals
// tightening (key is the section-scoped key path prefix).
func setDirection(loosen, tighten []string, key string, old, new []string) ([]string, []string) {
	for _, el := range new {
		if !slices.Contains(old, el) {
			loosen = append(loosen, key)
			break
		}
	}
	for _, el := range old {
		if !slices.Contains(new, el) {
			tighten = append(tighten, key)
			break
		}
	}
	return loosen, tighten
}
