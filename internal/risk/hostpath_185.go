package risk

// Ticket 185 (ledger A381) — the host-minted-path roster, DERIVED and not stored.
//
// What ticket 185 measured: every successful reread of a host pointer stamps one
// more mark, and that mark nobody declared becomes the blocker for the next
// reread of the very same path (真机 183-r1 log: first call kind=success
// rules_hit=[], second call kind=refused rules_hit=[R4] reason="R4: 包含来自
// fs.read C:\Users\…", close dropped=2). Ticket 183's exemption is per-mark and
// dies with the mark that declared it, so it never reached the read-back body.
//
// The shape owner approved on 2026-09-28 16:5x for it — 「那超长输出，我批了，要做的」,
// recorded as A381 — is: the exemption asks "is this path the very file the HOST
// itself minted in this scope", not "which path did the model just name". That
// approval is what overturned ticket 177's ruling that a per-scope path roster
// was not approved (A353; the sentence in taintmatch.go that said so is corrected
// in place, naming A381). 撤销口令：「撤分页」.
//
// What lands here is the narrowest delivery of that shape: NO new state.
//
//   - where the roster lives: nowhere separate. It is the projection of one
//     scope's own marks over the field ticket 183 already added —
//     fragmentIndex.declaredPath — which is set only when the host PROVED it
//     wrote that path into a result it just stamped (MarkWithHostPath's
//     runeIndexOf hit inside the declaring mark's own body).
//   - who writes it: the stamping step that already declares (bridge.mark ->
//     MarkWithHostPath with a non-empty hostPath, today only task.output's
//     backfilled artifact through hostPathBox). Nothing else can add an entry,
//     and no exported method hands out that power.
//   - who reads it: MarkWithHostPath, and only for a mark whose own recorded
//     provenance (origin) is exactly one of those paths — that is, a read-back of
//     a host-minted artifact. The candidate parameter is NEVER compared against
//     this list (ticket 183 AC#5 (1) still holds: no parameter-side by-value
//     pass, no name table lookup on the model's own string).
//   - how long it lives: exactly as long as the mark that carried the
//     declaration. Scope.Close drops the marks, so the roster drops with them;
//     re-opening the same id starts with an empty roster; a different scope id
//     never sees it (internal/risk/pointer_185_test.go L5/L6 and the CLI-seam leg
//     TestPointer185AnotherTaskCannotBorrowTheRoster are the catchers).
//
// Two deliberate bounds keep this from being a different contract:
//
//  1. value rule ONLY, never a positional span. The rostered mark excludes no
//     span when it is indexed — a window is skipped only if its own spelling
//     occurs inside the host path — so it can never drop MORE evidence than the
//     ticket 183 declaration on the same body (pinned quantitatively by
//     TestPointer185CostFaceOfTheRosterIsPinnedAndNarrowerThanTheSpan).
//  2. exact whole-string comparison over normalized paths, never a directory,
//     a prefix or a basename: a sibling artifact in the same artifacts
//     directory, a path that differs in one rune, and the same file name under a
//     different directory all stay evidence (ticket 177 W-3, 183 AC#4, 185 AC#7).
//
// The cost, stated with the benefit and not after it: body text of a
// host-read-back mark that happens to spell >=8 consecutive runes of that host
// path — the user-profile and artifacts-directory prefix is the realistic case,
// and 3 of the 10 bodies of ticket 183's frequency ruler are exactly that shape —
// stops being evidence against a later parameter. One fewer L2 card in that
// shape; the index never gets wider, and the same text arriving through any
// source the host did not mint (web.fetch, or a file the model named) still hits.

// hostMintedPathFor returns the normalized host-minted path that THIS SCOPE's own
// marks declared and that equals the recorded origin of the mark about to be
// built, or "" when this scope declared nothing matching. "" means "no exemption"
// — every miss is fail-closed, so a canonicalizer that hands back a path in
// another spelling only costs the feature, never the wall.
func (p *Provenance) hostMintedPathFor(scopeID, origin string) string {
	on := normalizeTaint(origin)
	if on == "" {
		return ""
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	reg := p.scopes[scopeID]
	if reg == nil {
		return ""
	}
	for _, m := range reg.marks {
		if m.idx == nil || len(m.idx.declaredPath) == 0 {
			continue // an ordinary mark declares nothing, so it minters nothing
		}
		if string(m.idx.declaredPath) == on {
			return on
		}
	}
	return ""
}
