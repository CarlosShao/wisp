package main

// Ticket 230 AC#4 (owner ruling 2026-11-05, ledger A1549) - the composition
// root's half of the standing B-tier exemption.
//
// WHAT WAS MISSING BETWEEN THE TWO LANDED HALVES. internal/config/overrides.go
// made risk.blacklist_overrides a real setting: every entry is canonicalized
// through the C26 resolver, refused if it classifies A (the A tier is absolutely
// denied and never overridable) and refused if it names no B-tier shape.
// internal/tools/bridge.go added Options.PersistedExemptions and
// internal/tools/mode.go's readBlacklist passes it to risk.GateStanding as the
// second of its two override sources, kept apart from the runtime L2 set so the
// security line can say which claim allowed a use. Between the two sat this file:
// with nobody arming the field, the bridge's closure answered nil ("the operator
// exempted nothing"), so a configured B-tier exemption still came back from the
// gate as "ask first" and the setting was true but inert - the shape the ticket
// exists to close (docs/evidence/s1/230-ac4-config-half.md section 4c).
//
// WHAT THIS FILE IS NOT. It is not a second classifier and not a second
// comparison fold. The tier verdict stays inside the frozen tables
// (internal/risk's A/B rules), and the comparison key is risk.OverrideKey - the
// exported form of the very fold risk.GateStanding applies to the path it is
// asked about, so the map built here and the lookup there cannot drift. It also
// cannot widen anything: risk.GateStanding checks the tier BEFORE it reads
// either override map, and no entry can reach this fold without having passed
// the load-time A-tier refusal, so an A path is denied whether or not a
// hand-built set ever names it.
//
// WHY THE SNAPSHOT IS ATOMIC RATHER THAN A LIVE READ. A live read would have to
// call config.Manager.Config per tool call, and that accessor returns a DEEP
// COPY of the whole config (internal/config/manager.go) - a whole-config copy on
// every call, to read one list. The list cannot change under this process
// without a reload (the reload path's plan/commit is the only runtime mutator of
// cur.Risk.BlacklistOverrides), so the root folds it once at assembly and again
// after every reload that committed something (cmd/wisp/config_reload.go's
// reloadOnce), and the bridge reads a snapshot. The maps built here are never
// mutated afterwards, so handing the same map to concurrent callers is safe.

import "github.com/CarlosShao/wisp/internal/risk"

// exemptionKeySet folds a C26-canonical exemption list into the key space
// risk.GateStanding compares in.
//
// Each entry is already canonical (the load pipeline published the resolver's
// answer, not the operator's spelling - internal/config/overrides.go), so this
// is a pure fold: OverrideKey lowercases, unifies separators and trims a
// trailing one. Keying the set with the canonical string instead would unlock
// nothing on this platform, because the key GateStanding looks up is the FOLDED
// form; that asymmetry is deliberate on the deny side and the reason a second
// copy of this fold must not appear here.
//
// An empty list folds to a nil map, which is also what the bridge's closure
// answers when no source is armed at all: "no path was ever exempted by config"
// is the fail-closed reading, and it keeps the default (every B-tier path asks)
// identical to what this process did before the ticket.
func exemptionKeySet(entries []string) map[string]bool {
	if len(entries) == 0 {
		return nil
	}
	set := make(map[string]bool, len(entries))
	for _, entry := range entries {
		set[risk.OverrideKey(entry)] = true
	}
	return set
}

// armPersistedExemptions rebuilds the standing exemption snapshot from the
// config this process currently considers effective, and swaps it in.
//
// The source is the Manager's live config rather than the boot snapshot on
// purpose: a reload commits into the Manager's cur (D36's plan/commit split),
// so folding here after CheckAndReload returns is what makes an added or removed
// exemption take effect without a restart. A denied loosening never reaches cur,
// so it also never reaches this fold: the old values (and the old set) stay.
//
// A nil Manager (a runtime assembled without one) leaves the snapshot nil, which
// is the same fail-closed default as an empty list.
func (rt *agentRuntime) armPersistedExemptions() {
	if rt == nil || rt.mgr == nil {
		return
	}
	set := exemptionKeySet(rt.mgr.Config().Risk.BlacklistOverrides)
	rt.exemptions.Store(&set)
}

// persistedExemptions is the closure tools.Options.PersistedExemptions receives:
// the ARMED set, read once per tool call from whatever the last arm installed.
// It returns a shared, read-only map (the bridge never writes it), and nil means
// nothing was ever armed.
func (rt *agentRuntime) persistedExemptions() map[string]bool {
	if rt == nil {
		return nil
	}
	if p := rt.exemptions.Load(); p != nil {
		return *p
	}
	return nil
}
