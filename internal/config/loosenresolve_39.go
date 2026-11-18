package config

// Ticket 39 AC#2's resolve half of the section-write gate (write leg
// 2026-11-18): the accepted answer of a queued loosening becomes a REAL write,
// and the refused/expired answer keeps the old value - both audited with the
// card id and the old->new values.
//
// WHY THIS FILE EXISTS. Ticket 195's WriteSection hands a native loosening to
// the LoosenQueue seam and returns immediately (ruling N9: C18's 300s deadline
// has no business inside a config write). Until this file the seam stopped
// there: whatever the adapter did with the card, the Manager had no way to
// turn an ACCEPTED card into the write it was raised for, and no record that
// named the card that authorized it. The adapter (cmd/wisp, ticket 39's
// assembly wiring) now calls ResolveQueuedLoosen with the card's answer:
//
//   - accepted  - the write is recomputed against the file AS IT NOW STANDS
//     (the card may have hung for a while; the base the original verdict
//     judged is stale by definition), validated, committed through the same
//     under-lock commit path every section write uses, and booked with ONE
//     audit line that names the card id and the old->new value of every key
//     the write changed;
//   - refused / timed out - nothing is written; the audit line says so, and
//     the old value is what any read-back will find (an accepted-card write
//     that never happened must not be inferable from the file).
//
// THE CARD IS THE ONLY AUTHORITY, AND IT IS SPENT ONCE. A pending entry is
// removed from the registry the moment it is resolved; a second resolve of the
// same card id is refused by name, so one answered card cannot be replayed
// into two writes. The registration happens BEFORE QueueLoosen is called and
// under the id the adapter is told to use (LoosenRequest.RequestID), so an
// instant answer can never arrive before the pending entry exists.
//
// DRIFT IS FAIL-CLOSED. If the file moved while the card was up and the
// recomputed loosening names keys the card never showed (a second key flipped
// in the same section), the write is refused: the operator approved what the
// card named, not what the file grew into meanwhile. A file that already
// carries the approved value resolves as "unchanged" (no write, no lie), and a
// file that moved the OTHER way (the candidate is now a tightening against it)
// commits as the tighten it now is - tightening needs no card (D36 rule 1).

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/CarlosShao/wisp/internal/observe"
)

// pendingLoosen is one queued loosening waiting on its card, held so the
// accepted answer can complete the write the card was raised for.
type pendingLoosen struct {
	w    SectionWrite
	keys []string // the loosening keys the card names
}

// loosenRequestSeq mints the request ids WriteSection hands to the queue. The
// id is the card's correlation id, the registry key, and the audit join all at
// once, which is why the Manager mints it BEFORE the queue is called.
var loosenRequestSeq atomic.Uint64

func mintLoosenRequestID() string {
	return fmt.Sprintf("host:config-loosen#%d", loosenRequestSeq.Add(1))
}

// registerPendingLoosen stores the pending write under id. Called before
// QueueLoosen hands the request out, so the id is already resolvable the
// instant the card exists.
func (m *Manager) registerPendingLoosen(id string, w SectionWrite, keys []string) {
	m.pendingLoosenMu.Lock()
	defer m.pendingLoosenMu.Unlock()
	if m.pendingLoosen == nil {
		m.pendingLoosen = map[string]*pendingLoosen{}
	}
	m.pendingLoosen[id] = &pendingLoosen{w: w, keys: append([]string(nil), keys...)}
}

// takePendingLoosen removes and returns the pending entry, nil when the id is
// unknown or already spent (one card, one write).
func (m *Manager) takePendingLoosen(id string) *pendingLoosen {
	m.pendingLoosenMu.Lock()
	defer m.pendingLoosenMu.Unlock()
	p := m.pendingLoosen[id]
	delete(m.pendingLoosen, id)
	return p
}

// ResolveQueuedLoosen completes one queued loosening with its card's answer.
// id is the request id the queue was handed (LoosenRequest.RequestID) and the
// card answered under; approved is the card's verdict - a timeout arrives as
// false because the gate itself resolves an unanswered card to a reject.
//
// Nothing here trusts the answer beyond what the card covers: the write is
// recomputed against the current file, re-validated, re-checked for ownership,
// and its loosening keys must still be covered by the keys the card named.
func (m *Manager) ResolveQueuedLoosen(id string, approved bool) (SectionWriteResult, error) {
	res := SectionWriteResult{Direction: DirLoosen}
	p := m.takePendingLoosen(id)
	if p == nil {
		return res, observe.New(observe.ClassConfig,
			"config: no pending loosening is registered under card id "+id+
				"; the id is unknown or its card was already resolved - one card authorizes one write")
	}
	w := p.w
	res.Keys = append([]string(nil), p.keys...)

	if !approved {
		// The refused/timeout arm: the file is left exactly as it stands, and
		// the audit line is the record that a card was asked and did not say
		// allow. Read-backs will find the old values.
		m.auditf("config: SECTION-WRITE section=%s key=%s direction=loosen source=%s origin=%s "+
			"outcome=denied card=%s keys=%s detail=%q",
			w.Section, w.Key, w.Source, w.Origin, id, strings.Join(p.keys, ","),
			"这张 L2 卡没有拿到「允许」（拒绝或超时）：一个字节都没写，旧值原样保留")
		return res, nil
	}

	// The accepted arm. The base is re-read because the card may have hung for
	// a while: the verdict the card approved was computed against the file as
	// it stood THEN.
	base, exists, err := m.readCurrentFile()
	if err != nil {
		m.auditf("config: SECTION-WRITE section=%s key=%s direction=loosen source=%s origin=%s "+
			"outcome=refused-unreadable card=%s err=%v", w.Section, w.Key, w.Source, w.Origin, id, err)
		return res, observe.Wrap(observe.ClassConfig, err,
			"config: the card for "+w.Key+" was accepted, but config.toml cannot be re-read; nothing was written")
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
			m.auditf("config: SECTION-WRITE section=%s key=%s direction=loosen source=%s origin=%s "+
				"outcome=refused-check card=%s err=%v", w.Section, w.Key, w.Source, w.Origin, id, err)
			return res, err
		}
	}
	if err := w.Apply(candidate); err != nil {
		m.auditf("config: SECTION-WRITE section=%s key=%s direction=loosen source=%s origin=%s "+
			"outcome=refused-apply card=%s err=%v", w.Section, w.Key, w.Source, w.Origin, id, err)
		return res, err
	}
	// The same one-key ownership rule the original write ran: the card
	// approved one key, not whatever the Apply closure reaches.
	changed := diffKeyPaths(base, candidate)
	for _, path := range changed {
		if path == w.Key || strings.HasPrefix(path, w.Key+".") {
			continue
		}
		m.auditf("config: SECTION-WRITE section=%s key=%s direction=loosen source=%s origin=%s "+
			"outcome=refused-ownership card=%s keys=%s", w.Section, w.Key, w.Source, w.Origin, id,
			strings.Join(changed, ","))
		return res, observe.New(observe.ClassConfig,
			"config: the accepted card for "+w.Key+" would also change "+strings.Join(changed, ", ")+
				"; refusing a completion that carries keys outside the one the card named")
	}
	if err := validate(candidate); err != nil {
		m.auditf("config: SECTION-WRITE section=%s key=%s direction=loosen source=%s origin=%s "+
			"outcome=refused-validate card=%s err=%v", w.Section, w.Key, w.Source, w.Origin, id, err)
		return res, observe.Wrap(observe.ClassConfig, err,
			"config: the accepted card for "+w.Key+" no longer survives validation against the current file; nothing was written")
	}

	dir, keys := lockedDirectionOf(w.Section, base, candidate)
	switch dir {
	case DirNeutral:
		// The file already carries the approved value (someone wrote it while
		// the card was up). Nothing to write, and saying "written" would be
		// the one lie an audit line must not tell.
		m.auditf("config: SECTION-WRITE section=%s key=%s direction=loosen source=%s origin=%s "+
			"outcome=unchanged card=%s keys=%s detail=%q",
			w.Section, w.Key, w.Source, w.Origin, id, strings.Join(keys, ","),
			"卡通过时盘上已经是这个值了；没有写任何字节")
		return res, nil
	case DirTighten:
		// The file moved the other way while the card was up: against what the
		// file holds NOW the candidate is a tightening, and tightening applies
		// quietly (D36 rule 1). Commit as what it is, booked as tighten.
		res.Direction = DirTighten
		res.Keys = keys
	}
	if dir == DirLoosen {
		// Every loosening key the recompute finds must be covered by the keys
		// the card named - exact or within the same subtree (the dynamic
		// plugin families emit parent/child spellings). Anything else is an
		// approval the operator never gave.
		if stray := uncoveredLoosenKeys(keys, p.keys); len(stray) > 0 {
			m.auditf("config: SECTION-WRITE section=%s key=%s direction=loosen source=%s origin=%s "+
				"outcome=refused-drift card=%s approved=%s found=%s detail=%q",
				w.Section, w.Key, w.Source, w.Origin, id,
				strings.Join(p.keys, ","), strings.Join(stray, ","),
				"卡挂着的时候文件又变了：这次写会放宽卡上没有的键；一个字节都没写（fail-closed）")
			return res, observe.New(observe.ClassConfig,
				"config: the accepted card for "+w.Key+" no longer covers "+strings.Join(stray, ", ")+
					" (the file changed while the card was up); refusing the completion, nothing was written")
		}
		res.Keys = keys
	}

	// The commit itself is the same under-lock path every section write uses:
	// memory first, guarded merge write, rollback on failure. restartTier is
	// false here on purpose (ticket 39 AC#1): these keys live in the four
	// LOCKED sections, whose TierRegistry marker is "locked" - not a reload
	// tier - so a completed card applies to memory exactly as before.
	res, err = m.commitSectionWrite(w, base, exists, res, false)
	if err != nil {
		return res, err
	}
	outcome := "written"
	if !res.Written {
		outcome = "unchanged"
	}
	m.auditf("config: SECTION-WRITE section=%s key=%s direction=%s source=%s origin=%s "+
		"outcome=%s card=%s keys=%s changes=%s",
		w.Section, w.Key, res.Direction, w.Source, w.Origin, outcome, id,
		strings.Join(res.Keys, ","), renderLoosenChanges(base, candidate, changed))
	return res, nil
}

// uncoveredLoosenKeys returns the keys in found that the approved card does
// not cover: a key is covered when it equals an approved key, is that key's
// child, or is its parent (pluginsDirection emits a parent spelling for a new
// entry and child spellings for its leaves).
func uncoveredLoosenKeys(found, approved []string) []string {
	var out []string
	for _, f := range found {
		ok := false
		for _, a := range approved {
			if f == a || strings.HasPrefix(f, a+".") || strings.HasPrefix(a, f+".") {
				ok = true
				break
			}
		}
		if !ok {
			out = append(out, f)
		}
	}
	return out
}

// unrosteredLoosenKeys returns the loosening keys the named roster does not
// carry (the gate's fail-closed pin on auditor/roster drift).
func unrosteredLoosenKeys(keys []string) []string {
	var out []string
	for _, k := range keys {
		if !IsSecurityLoosenKey(k) {
			out = append(out, k)
		}
	}
	return out
}

// renderLoosenChanges renders the old->new value of every changed key, the
// factual half the accepted write's audit line carries ("老值→新值 + 卡 id").
// Values are read off the two configs the write was judged between, spelled
// the way the file spells them (the same walk SectionRows uses).
func renderLoosenChanges(old, new *Config, keys []string) string {
	oldMap := sectionValueMap(old)
	newMap := sectionValueMap(new)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s:%s->%s", k, oldMap[k], newMap[k]))
	}
	out := strings.Join(parts, ", ")
	if len(out) > 512 {
		out = out[:512] + "...(truncated)"
	}
	return out
}

// sectionValueMap flattens one config into dotted key -> rendered value.
func sectionValueMap(c *Config) map[string]string {
	m := map[string]string{}
	for _, row := range SectionRows(c) {
		m[row.Key] = row.Value
	}
	return m
}
