package config

// The programmatic section-write entry, and the D36 direction verdict that runs
// BEFORE a single byte of config.toml is replaced (ticket 195, orchestrator
// ruling 2026-10-02 section 1).
//
// WHY THIS FILE EXISTS. The verdict "is this key a loosening?" lived in exactly
// one place: the reload path. manager.go's planLocked asks
// riskDirection/fsDirection/netDirection/pluginsDirection, and CheckAndReload
// raises the L2 re-confirmation with mu released, so a HAND EDIT of
// [risk]/[fs]/[net]/[plugins] cannot loosen anything silently: with no hook
// wired the answer is a denial (Manager.ConfirmLocked, fail-closed). What the
// same verdict did not cover was a PROGRAMMATIC write. settings.go's writeOneKey
// (ticket 248) validates a value and merges one key, and the two older writers
// (SetPermissionMode, AddAllowedDir) state in their own headers that their
// caller already raised the card - none of the three asks whether the key it
// just wrote was a loosening. Ticket 194-c1's census (ledger A384 ruling 2)
// measured the consequence: the only reason fs.delete_enabled and
// fs.reparse_point_exceptions cannot be loosened through a programmatic write
// today is that NO writer addresses them. This file makes the verdict part of
// the write entry itself, so the first writer that addresses a locked key
// carries it by construction.
//
// WHAT IT IS NOT. It is not a second settings page: the setters writeOneKey
// serves (the six [llm] keys plus ticket 313/320's voice and billing fields)
// keep their own names, messages and one-key-per-call receipt (ruling section 5
// forbids rebuilding settings.go), and writeOneKey is now the thin adapter over
// this entry rather than a parallel body. It also adds no capability to the
// panel: there is no "allow" parameter anywhere in this file, so no caller -
// panel or otherwise - can approve its own loosening through it. The panel's
// config.set route is still refused wholesale by internal/panel's
// lockedFieldFamilies (ticket 248) and this ticket does not touch that.

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/CarlosShao/wisp/internal/observe"
)

// WriteSource says who asked for a programmatic write. It is a closed set: an
// empty or invented value is refused rather than defaulted, because the whole
// point of the type is that "who asked" decides what a loosening may do.
type WriteSource string

const (
	// SourceNative is the native side: the operator's own surface, the CLI and
	// the resident approvals path. A native loosening is QUEUED (one pending
	// L2 request), never written in place and never waited on.
	SourceNative WriteSource = "native"
	// SourcePanel is the panel side (C17). SPEC-06:19 and R20 make the panel a
	// permission INPUT - display plus request - and never an authority, so a
	// panel-sourced loosening of a locked section is refused here and there is
	// no parameter on this entry that can approve one.
	SourcePanel WriteSource = "panel"
)

// valid reports whether s is one of the two defined sources.
func (s WriteSource) valid() bool { return s == SourceNative || s == SourcePanel }

// The two refusals this entry can hand back on a loosening. Both are wrapped in
// an observe.Error so the caller sees the same config-class shape every other
// refused write in this package returns.
var (
	// ErrPanelLoosenRefused is the structural refusal AC#3 asks for: a
	// panel-sourced loosening of a locked section. Nothing is written, nothing
	// is queued, and the in-memory config is untouched.
	ErrPanelLoosenRefused = errors.New(
		"config: a panel-sourced loosening of a locked section is refused; the panel is a permission input (display plus request), not an authority")
	// ErrLoosenQueueUnwired is the fail-closed answer when no approval queue is
	// wired: a loosening that cannot be handed to C18 is not written either
	// (the same nil-means-deny rule Manager.ConfirmLocked already carries).
	ErrLoosenQueueUnwired = errors.New(
		"config: a locked-section loosening needs the L2 approval queue, and this Manager has none wired; refusing to write it")
)

// LoosenRequest is one pending loosening this entry did NOT write. It names
// what would change and who asked, which is what a human needs to answer the
// card, and it carries no value and no authority of its own: the write, if it
// ever happens, goes back through the same validation and the same queue.
type LoosenRequest struct {
	// Section is the locked section the change would loosen ("risk", "fs",
	// "net", "plugins").
	Section string
	// Keys are the dotted key paths the reload path itself would call a
	// loosening (the direction auditor's spelling, so one key has one name in
	// a log line, in a card and in a D36 verdict).
	Keys []string
	// Source is who asked.
	Source WriteSource
	// Origin names the route that asked ("panel.settings", "cli.approve", ...),
	// so the audit line and the card can say who to blame without guessing.
	Origin string
}

// LoosenQueue is the C18 approval queue as this package needs it: hand over one
// pending request and get an id back. It must NOT block on a human - the entry
// below returns the moment the request is queued, and waiting for a card here
// would drag C18's 300s deadline into the config write path (ruling N9). The
// production adapter is wired at the assembly root; nil refuses (fail-closed).
type LoosenQueue interface {
	QueueLoosen(req LoosenRequest) (id string, err error)
}

// AuditFunc receives one attribution line per adjudicated write. It is called
// without mu held. nil means "the process log" (see auditf).
type AuditFunc func(format string, args ...any)

// SectionWrite is one programmatic write of one dotted key inside one section.
type SectionWrite struct {
	// Section is the config.toml section the key belongs to. It must be the
	// key's own first segment: a write that names a different section than the
	// key it addresses is refused rather than guessed at.
	Section string
	// Key is the dotted key path this write owns, spelled the way the direction
	// auditor spells it ("fs.delete_enabled"). One call, one key.
	Key string
	// Source is who asked. SourcePanel may not loosen a locked section.
	Source WriteSource
	// Origin names the route that asked, for the audit line (and for the card
	// if the write turns out to be a loosening).
	Origin string
	// Check is the pre-write refusal that names a missing row (settings.go's
	// three refusal reasons). Optional.
	Check func(*Config) error
	// Apply stamps the value onto a config. It must touch Key and nothing else:
	// the entry verifies that by diffing, so a caller cannot smuggle a second
	// key in on the back of one receipt line.
	Apply func(*Config) error
}

// SectionWriteResult is what one call actually did.
type SectionWriteResult struct {
	// Direction is the verdict the entry took before it touched anything
	// (neutral for every unlocked section).
	Direction Direction
	// Written reports whether config.toml changed. False on a queued
	// loosening, a refusal, and a write whose value the file already carried.
	Written bool
	// PendingID is the id the approval queue returned when a loosening was
	// queued instead of written. Empty in every other case.
	PendingID string
	// Keys lists the key paths the file actually changed by (the merged write's
	// own report), empty when nothing changed.
	Keys []string
}

// WriteSection writes one key of one section, and it takes the D36 direction
// verdict BEFORE it writes:
//
//   - neutral or tightening - the guarded merge write runs exactly as ticket
//     248's settings path does (validate first, one key, memory rolled back if
//     the file write fails);
//   - loosening from the native side - NOTHING is written and nothing in
//     memory changes: one LoosenRequest is handed to the approval queue and
//     this returns immediately with PendingID set;
//   - loosening from the panel - refused outright, nothing written, nothing
//     queued (ErrPanelLoosenRefused);
//   - loosening with no queue wired - refused (ErrLoosenQueueUnwired), which is
//     the fail-closed direction the reload path already takes with a nil hook.
//
// The verdict itself is computed with the SAME four helper functions the reload
// path uses, so a key planLocked calls a loosening cannot be a tightening here.
// The verdict, the validation and the file read all happen with mu released
// (ruling N9); only the merge write and the in-memory apply run under it.
func (m *Manager) WriteSection(w SectionWrite) (SectionWriteResult, error) {
	res := SectionWriteResult{Direction: DirNeutral}

	if w.Apply == nil {
		return res, observe.New(observe.ClassConfig,
			"config: section write for "+w.Key+" has no value to apply; nothing was written")
	}
	if !w.Source.valid() {
		return res, observe.New(observe.ClassConfig,
			"config: section write for "+w.Key+" names no defined source ("+string(w.Source)+
				"); nothing was written")
	}
	section := sectionOfKey(w.Key)
	if w.Section == "" || w.Section != section {
		return res, observe.New(observe.ClassConfig,
			fmt.Sprintf("config: section write for %s names section %q, but the key's own section is %q; "+
				"nothing was written", w.Key, w.Section, section))
	}
	if !knownSection(w.Section) {
		return res, observe.New(observe.ClassConfig,
			"config: section write for "+w.Key+" names a section this schema does not have ("+
				w.Section+"); nothing was written")
	}

	// (1) The base this write is judged against is what the FILE holds, not the
	// in-memory snapshot: that is ticket 226's rule and it is the same base the
	// guarded merge writer will re-read below. A missing file falls back to
	// memory exactly as writeOneKey always has (first-run may create one).
	base, exists, err := m.readCurrentFile()
	if err != nil {
		// The file is here and this process cannot read it. Overwriting it is
		// the exact destruction ticket 226 was filed for, so the write gets the
		// refusal, not the file.
		return res, err
	}
	if !exists {
		base = deepCopyConfig(m.Config())
		if base == nil {
			base = NewDefaults()
		}
	}

	candidate := deepCopyConfig(base)
	if w.Check != nil {
		if err := w.Check(candidate); err != nil {
			return res, err
		}
	}
	if err := w.Apply(candidate); err != nil {
		return res, err
	}

	// (2) The write owns w.Key and nothing beside it: every path it changes must
	// be that key or one of its children (ticket 313's chat pair writes
	// llm.roles.chat.model and .provider under the one key llm.roles.chat). A
	// write that reached an UNRELATED key is refused here rather than landing
	// under a receipt that never names it.
	changed := diffKeyPaths(base, candidate)
	for _, path := range changed {
		if path == w.Key || strings.HasPrefix(path, w.Key+".") {
			continue
		}
		return res, observe.New(observe.ClassConfig,
			"config: section write owns "+w.Key+" and would also change "+
				strings.Join(changed, ", ")+"; refusing a write that carries keys outside the one it owns")
	}

	// (3) Ticket 248's AC#7 gate: the SAME validate() the start-up path runs,
	// over the base the write would produce, before a single byte is written.
	if err := validate(candidate); err != nil {
		return res, observe.Wrap(observe.ClassConfig, err,
			"config: refused to write "+w.Key+"; the value would not survive validation, so config.toml is unchanged。"+
				refusalInvalid+"：值本身过不了这份 schema 的校验，这一行要改的是值；文件一个字节都没动。")
	}

	// (4) The verdict, still with mu released (ruling N9).
	dir, keys := lockedDirectionOf(w.Section, base, candidate)
	res.Direction = dir
	res.Keys = nil
	if dir == DirLoosen {
		res.Keys = keys
		switch {
		case w.Source == SourcePanel:
			m.auditf("config: SECTION-WRITE section=%s key=%s direction=loosen source=%s origin=%s outcome=refused-panel",
				w.Section, w.Key, w.Source, w.Origin)
			return res, observe.Wrap(observe.ClassConfig, ErrPanelLoosenRefused,
				"config: refused to write "+w.Key+": "+strings.Join(keys, ", ")+" would loosen "+w.Section+
					", and this request came from the panel")
		case m.LoosenQueue == nil:
			m.auditf("config: SECTION-WRITE section=%s key=%s direction=loosen source=%s origin=%s outcome=refused-no-queue",
				w.Section, w.Key, w.Source, w.Origin)
			return res, observe.Wrap(observe.ClassConfig, ErrLoosenQueueUnwired,
				"config: refused to write "+w.Key+": "+strings.Join(keys, ", ")+" would loosen "+w.Section+
					" and no approval queue is wired, so the request has nowhere to go")
		}
		id, qerr := m.LoosenQueue.QueueLoosen(LoosenRequest{
			Section: w.Section,
			Keys:    keys,
			Source:  w.Source,
			Origin:  w.Origin,
		})
		if qerr != nil {
			m.auditf("config: SECTION-WRITE section=%s key=%s direction=loosen source=%s origin=%s outcome=queue-failed err=%v",
				w.Section, w.Key, w.Source, w.Origin, qerr)
			return res, observe.Wrap(observe.ClassConfig, qerr,
				"config: could not queue the "+strings.Join(keys, ", ")+" loosening ("+w.Section+
					"); nothing was written")
		}
		res.PendingID = id
		m.auditf("config: SECTION-WRITE section=%s key=%s direction=loosen source=%s origin=%s outcome=queued id=%s keys=%s",
			w.Section, w.Key, w.Source, w.Origin, id, strings.Join(keys, ","))
		return res, nil
	}
	// (5) The write itself: memory first, then the file, then memory back if the
	// file refused - the same order and rollback settings.go has always used.
	// The tightening audit line is written AFTER this returns, because
	// "outcome=written" before the merge writer has answered would be the one
	// lie an audit line must never tell.
	res, err = m.commitSectionWrite(w, base, exists, res)
	if err != nil {
		return res, err
	}
	if dir == DirTighten {
		outcome := "written"
		if !res.Written {
			outcome = "unchanged"
		}
		m.auditf("config: SECTION-WRITE section=%s key=%s direction=tighten source=%s origin=%s outcome=%s keys=%s",
			w.Section, w.Key, w.Source, w.Origin, outcome, strings.Join(res.Keys, ","))
	}
	return res, nil
}

// commitSectionWrite is the under-lock half of WriteSection: the in-memory
// apply, the guarded merge write, and the rollback that keeps memory and file
// from diverging on a failed write. The direction verdict, the validation and
// the file pre-read all happened before this was called (ruling N9).
func (m *Manager) commitSectionWrite(w SectionWrite, base *Config, exists bool,
	res SectionWriteResult,
) (SectionWriteResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	memOld := deepCopyConfig(m.cur)
	if memOld == nil {
		return res, observe.New(observe.ClassConfig, "config: no loaded config in this Manager, refusing to write one")
	}
	if err := w.Apply(m.cur); err != nil {
		m.cur = memOld
		return res, err
	}
	var applyErr error
	if err := m.mergeWrite(w.Key, func(b *Config) { applyErr = w.Apply(b) }); err != nil {
		m.cur = memOld
		return res, observe.Wrap(observe.ClassConfig, err,
			"config: section write for "+w.Key+" failed; the values in memory are the ones this process started with")
	}
	if applyErr != nil {
		// The file changed underneath between the verdict and the merge, and the
		// key no longer exists on what the file now holds. mergeWrite wrote
		// nothing (its own re-read saw no change), and memory goes back.
		m.cur = memOld
		return res, applyErr
	}

	res.Keys = m.writtenKeyPaths(base, exists, w.Key)
	res.Written = len(res.Keys) > 0
	return res, nil
}

// lockedDirectionOf is the shared verdict: it returns the direction of the
// change between two configs for one section, plus the key paths on the winning
// side (loosening wins over tightening, the same precedence planLocked uses).
// Unlocked sections are always neutral - their tier is what governs them, not
// D36 rule 1.
func lockedDirectionOf(section string, old, new *Config) (Direction, []string) {
	if old == nil || new == nil {
		return DirNeutral, nil
	}
	var loosen, tighten []string
	switch section {
	case "risk":
		loosen, tighten = riskDirection(&old.Risk, &new.Risk)
	case "fs":
		loosen, tighten = fsDirection(&old.FS, &new.FS)
	case "net":
		loosen, tighten = netDirection(&old.Net, &new.Net)
	case "plugins":
		loosen, tighten = pluginsDirection(&old.Plugins, &new.Plugins)
	default:
		return DirNeutral, nil
	}
	switch {
	case len(loosen) > 0:
		return DirLoosen, loosen
	case len(tighten) > 0:
		return DirTighten, tighten
	}
	return DirNeutral, nil
}

// sectionOfKey returns the first segment of a dotted key path ("" for a key
// that has none).
func sectionOfKey(key string) string {
	if i := strings.IndexByte(key, '.'); i > 0 {
		return key[:i]
	}
	return ""
}

// knownSection reports whether the schema has this top-level section. Whole
// sections are registered by name; [app] and [voice] are registered per key, so
// a section counts as known when the registry holds one of its keys.
func knownSection(name string) bool {
	if name == "" {
		return false
	}
	if _, ok := TierRegistry[name]; ok {
		return true
	}
	for path := range TierRegistry {
		if strings.HasPrefix(path, name+".") {
			return true
		}
	}
	return false
}

// auditf is the one place an adjudicated write says who asked and what was
// decided. Called without mu held: a sink may read the config back.
func (m *Manager) auditf(format string, args ...any) {
	if m.Audit != nil {
		m.Audit(format, args...)
		return
	}
	slog.Info(fmt.Sprintf(format, args...))
}
