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
// Ticket 226 put a condition on that adoption (only content this write produced
// is claimed) and moved the whole persist behind a re-read, so a hand edit of any
// OTHER key survives this write - see writeguard.go.

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
// success rather than duplicating the line. Since ticket 226 "already
// allowlisted" is measured against the FILE too, not only against memory: an
// entry the operator hand-added while this process ran is left where it is and
// not written a second time.
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
	return m.writeAllowedDirs(next, func(base *Config) {
		// Append to what the file holds, never over it: a directory the operator
		// added by hand while this process ran is their entry, and this answer's
		// business is exactly one line.
		for _, have := range base.FS.AllowedDirs {
			if strings.TrimSpace(have) == d {
				return
			}
		}
		base.FS.AllowedDirs = append(append([]string(nil), base.FS.AllowedDirs...), d)
	})
}

// SetAllowedDirs persists a whole list (the撤销 half of AddAllowedDir: dropping
// one line is a rewrite of the same section).
//
// Tightening this section needs no L2 re-confirmation - applyLocked's own
// bookkeeping says a tightening is hot-applied without asking - so a caller that
// only removes entries is not obliged to raise a card. A caller that ADDS one is,
// and this method cannot tell the two apart, which is why it is exported under a
// name that says SET and the additive door above carries the doc.
//
// "SET" is scoped to the one key it names: since ticket 226 this writes
// fs.allowed_dirs and nothing else, so every other key in the file keeps the
// value the file has. It is still authoritative for that key - a directory the
// operator added by hand and did not pass here is dropped, which is what a whole
// list means - and the write reports that key as one it changed.
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
	return m.writeAllowedDirs(cleaned, func(base *Config) {
		base.FS.AllowedDirs = append([]string(nil), cleaned...)
	})
}

// writeAllowedDirs is the shared body: apply to memory, persist ONE key through
// the guarded merge writer (writeguard.go re-reads the file and replaces
// fs.allowed_dirs on what it actually holds), and roll the memory value back if
// the write fails. Caller holds m.mu.
//
// The two behaviours ticket 201's ruling pinned are unchanged and neither is this
// body's to soften: a refusal (bad entry, unreadable file, failed write) leaves
// memory and file exactly as they were, and only a write this process can vouch
// for - byte-for-byte what it just put there - is adopted as our own (that part
// moved into mergeWrite, which is why it is no longer four lines here).
func (m *Manager) writeAllowedDirs(next []string, setOn func(base *Config)) error {
	old := m.cur.FS.AllowedDirs
	m.cur.FS.AllowedDirs = next
	if err := m.mergeWrite(keyFSAllowedDirs, setOn); err != nil {
		m.cur.FS.AllowedDirs = old
		return observe.Wrap(observe.ClassConfig, err,
			"config: allowed_dirs write failed; keeping the previous allowlist in memory")
	}
	return nil
}

// statOwnWrite adopts the watched file's new mtime+size after a write this
// Manager made itself, so CheckAndReload does not read the program's own write as
// an unconfirmed hand edit. It is called from exactly one place - mergeWrite, at
// the end of the branch that proved the bytes it wrote are the bytes memory holds
// (ticket 226 AC#3 narrowed it from "we wrote the file once" to "this write
// produced this content"). A file it cannot stat is left alone: the next poll then
// sees the change, which is the conservative direction.
func (m *Manager) statOwnWrite() {
	st, err := os.Stat(m.path)
	if err != nil {
		return
	}
	m.seenMtime, m.seenSize, m.seenValid = st.ModTime(), st.Size(), true
}
