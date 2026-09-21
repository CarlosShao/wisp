package config

import (
	"os"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/risk"
)

// Permission-mode read/write helpers (ticket 90, ruling R20/M3).
//
// R20/M3 makes the mode the ONE user preference in this system that persists
// across a new session and a restart: "手动选过哪档就一直按那档". That is
// deliberately narrow, and the narrowness is the boundary the orchestrator
// pinned when the ruling was recorded:
//
//   - persisted here: the MODE, as a single string in config.toml, which is
//     already the project's single source of truth for configuration (D6);
//   - NOT persisted anywhere: a D45 session grant. GrantScopeSession stays
//     session-scoped and dies with the process (PLAN.md:1640, ticket 49), and
//     PLAN.md:1640's "会话授权不得覆盖 C25 污染升级" stands untouched. There is
//     no code path in this file that can write a grant, and the two behaviors
//     are pinned by separate tests (AC#3a/3b vs AC#3b) precisely so nobody can
//     merge them into one "persistence" case later.
//
// Persistence rides the reload manager rather than a sidecar file: a mode that
// lived outside config.toml would give the same setting two truths, which is
// the split-state bug SPEC-03 §3.1 exists to prevent.

// PermissionMode returns the effective mode of this config.
//
// An unparseable value yields the strictest mode rather than an error, because
// this accessor is on the hot decision path (the bridge reads it once per tool
// call) and "could not read the mode" must fail toward asking, never toward
// silence. A load-time config cannot normally reach that state - validateRisk
// rejects an unknown mode before the Config is ever handed out - so this is a
// defense against a hand-built Config struct, not a second parser.
func (c *Config) PermissionMode() risk.Mode {
	if c == nil {
		return risk.DefaultMode()
	}
	m, err := risk.ParseMode(c.Risk.PermissionMode)
	if err != nil {
		return risk.DefaultMode()
	}
	return m
}

// SetPermissionMode applies a mode to the effective config and persists it to
// the watched file, atomically (atomicWrite) and in canonical form.
//
// It is the programmatic counterpart of the operator editing config.toml by
// hand, and the difference matters for D36: a hand-edit that LOOSENS a locked
// section goes through Manager.ConfirmLocked (an L2 re-confirmation), while a
// switch made here has already been confirmed by its own caller - ticket 90's
// mode switcher raises exactly one L2 card, and only for auto_approve (R20/M4).
// Re-running the D36 hook here would ask the same question twice, so instead
// this method adopts the file's new mtime+size: the change is already applied
// in memory, and CheckAndReload must not read the program's own write as an
// unconfirmed edit.
//
// A failed write rolls the in-memory value back, so the mode the runtime acts
// on and the mode the next start reads cannot diverge.
func (m *Manager) SetPermissionMode(mode risk.Mode) error {
	if !mode.Valid() {
		return observe.New(observe.ClassConfig, "config: refusing to persist an undefined permission mode")
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	old := m.cur.Risk.PermissionMode
	m.cur.Risk.PermissionMode = mode.String()
	if err := SaveFile(m.path, m.cur); err != nil {
		m.cur.Risk.PermissionMode = old
		return observe.Wrap(observe.ClassConfig, err,
			"config: permission_mode write failed; keeping the previous mode in memory")
	}
	if st, statErr := os.Stat(m.path); statErr == nil {
		m.seenMtime, m.seenSize, m.seenValid = st.ModTime(), st.Size(), true
	}
	return nil
}
