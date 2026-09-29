package config

// The [fs] allowlist write path (ticket 201 AC#4/AC#5's landing spot).
//
// WHY THIS FILE EXISTS. PLAN.md:1645-1646 freezes one sentence about a loosening
// of [risk]/[fs]/[net]/[plugins]: 热加载放宽「必须触发 L2 级重新确认，不得静默
// 生效」. The reload side of that rule has existed since ticket 06 -
// Manager.ConfirmLocked is the hook, and a nil hook denies (fail closed), which
// is what `cmd/wisp run` passes today. What did NOT exist was the other half: no
// in-process caller could add one [fs] allowed_dirs entry at all, so a 「一直允
// 许」 answer had nowhere to go. Permmode.go's own header names the asymmetry -
// the permission MODE is programmatic and persisted, everything else is a hand
// edit - and ticket 201 AC#5's rewritten form says the long-lived branch's real
// landing spot is exactly one line in allowed_dirs.
//
// THE ONE RULE THIS FILE MUST NOT SOFTEN. A write here widens a locked section,
// so its caller owes the L2 re-confirmation first. This is the same contract
// SetPermissionMode documents above its own name: the programmatic path is
// confirmed by ITS caller (ticket 90's mode switcher raises exactly one L2 card),
// and re-running the D36 reload hook here would ask the same question twice.
// cmd/wisp's reply listener is therefore the confirmation, and the pair is pinned
// by a test (cmd/wisp/approval_always_201_test.go) that fails if the widening is
// persisted without a settled L2 answer in front of it. Nothing in this file can
// raise, fake or skip that card - it has no Gate, no UI and no route to either.
//
// Adopting the file's new mtime+size is what makes the write take effect instead
// of being read back as an unconfirmed hand edit; a hand edit that loosens [fs]
// still goes through ConfirmLocked and is still denied when nothing is wired.

import (
	"os"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
)

// AllowedDirs returns the readable/writable roots this config carries, as one
// copy so a caller cannot mutate the manager's own slice by appending.
func (c *Config) AllowedDirs() []string {
	if c == nil {
		return nil
	}
	return append([]string(nil), c.FS.AllowedDirs...)
}

// AddAllowedDir persists one more [fs] allowed_dirs entry.
//
// It is additive by construction: an existing entry is left where it is, so the
// order of what the operator already approved is never reshaped by one answer,
// and adding a directory that is already allowlisted is a no-op that reports
// success rather than duplicating the line.
//
// A refusal is loud and specific, and every refusal path leaves both the file and
// the in-memory config exactly as they were: the value the runtime acts on and
// the value the next start reads cannot diverge.
func (m *Manager) AddAllowedDir(dir string) error {
	d := strings.TrimSpace(dir)
	if d == "" {
		return observe.New(observe.ClassConfig, "config: refused to persist an empty allowed_dirs entry")
	}
	if strings.Contains(d, "..") {
		// A stored allowlist entry with a relative hop in its text is not a root,
		// and writing it would hand the next start something C26 has to resolve
		// before any verdict exists. Fail closed instead.
		return observe.New(observe.ClassConfig,
			"config: refused to persist an allowed_dirs entry containing \"..\": "+d)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, have := range m.cur.FS.AllowedDirs {
		if strings.TrimSpace(have) == d {
			return nil
		}
	}
	next := append(append([]string(nil), m.cur.FS.AllowedDirs...), d)
	return m.writeAllowedDirs(next)
}

// SetAllowedDirs persists a whole list (the撤销 half of AddAllowedDir: dropping
// one line is a rewrite of the same section).
//
// Tightening this section needs no L2 re-confirmation - applyLocked's own
// bookkeeping says a tightening is hot-applied without asking - so a caller that
// only removes entries is not obliged to raise a card. A caller that ADDS one is,
// and this method cannot tell the two apart, which is why it is exported under a
// name that says SET and the additive door above carries the doc.
func (m *Manager) SetAllowedDirs(dirs []string) error {
	cleaned := make([]string, 0, len(dirs))
	for _, d := range dirs {
		s := strings.TrimSpace(d)
		if s == "" {
			continue
		}
		if strings.Contains(s, "..") {
			return observe.New(observe.ClassConfig,
				"config: refused to persist an allowed_dirs entry containing \"..\": "+s)
		}
		cleaned = append(cleaned, s)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.writeAllowedDirs(cleaned)
}

// writeAllowedDirs is the shared body: apply, persist, roll back on a failed
// write, and adopt the file's new stat so the program's own write is not read
// back as an unconfirmed hand edit. Caller holds m.mu.
func (m *Manager) writeAllowedDirs(next []string) error {
	old := m.cur.FS.AllowedDirs
	m.cur.FS.AllowedDirs = next
	if err := SaveFile(m.path, m.cur); err != nil {
		m.cur.FS.AllowedDirs = old
		return observe.Wrap(observe.ClassConfig, err,
			"config: allowed_dirs write failed; keeping the previous allowlist in memory")
	}
	m.statOwnWrite()
	return nil
}

// statOwnWrite adopts the watched file's new mtime+size after a write this
// Manager made itself, so CheckAndReload does not read the program's own write as
// an unconfirmed hand edit (the same four lines SetPermissionMode runs, shared
// here rather than copied a third time). A file it cannot stat is left alone: the
// next poll then sees the change, which is the conservative direction.
func (m *Manager) statOwnWrite() {
	st, err := os.Stat(m.path)
	if err != nil {
		return
	}
	m.seenMtime, m.seenSize, m.seenValid = st.ModTime(), st.Size(), true
}
