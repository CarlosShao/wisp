package panel

// Readface wave 4 (2026-11-18), ticket 40 AC#2's wire instruments: the two
// roster names this wave added (session.rules.list / session.rules.revoke).
// The wave-3 file (readface_wave3_test.go) is the shape these follow; nothing
// here repeats its assertions - each check below is the one wave 4 needs and
// wave 3 does not carry.
//
// What is pinned, and why:
//
//  1. The two names are on the whitelist and route to their OWN socket - the
//     closed-roster ruler and the cmd/wisp fuzz table both grew with the same
//     commit, but the router case is the wire itself, so it is asserted here
//     where the router lives.
//  2. The pair declares EXACTLY two methods and the list half stays a read:
//     the list's payload contract (real fields, explicit zero) lives in the
//     cmd/wisp door; what this file pins is the wire shape the page relies on.
//  3. The revoke selector is selector-shaped and the envelope field survives
//     the parser under its own JSON key - the key the page's literal echoes.
//  4. The naming half is pinned where the instrument that forced it lives:
//     both route values are CLEAN under grantRouteWords (they must be, or the
//     frozen grant-boundary instrument would already be red - this check is
//     the witness that the wave-4 spellings were chosen with that vocabulary
//     in hand, not by luck), and the ban words are NOT re-litigated here.

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

// wave4Envelope builds the ONE envelope the parser accepts for method, with
// the optional selector fields.
func wave4Envelope(method, rid string, fields map[string]string) string {
	m := map[string]string{"method": method, "requestId": rid, "source": ComposerRequestSource}
	for k, v := range fields {
		m[k] = v
	}
	b, _ := json.Marshal(m)
	return string(b)
}

type wave4Rules struct{ list, revoke int }

func (h *wave4Rules) HandleSessionRulesList(context.Context, ComposerRequest) (string, error) {
	h.list++
	return "授权行", nil
}

func (h *wave4Rules) HandleSessionRulesRevoke(context.Context, ComposerRequest) (string, error) {
	h.revoke++
	return "已撤销", nil
}

// TestWave4RosterNamesRouteToTheirOwnSockets is the router half: each of the
// two names dispatches to exactly its own method and no other method moves,
// and the running guard accepts each name (a case label the whitelist refuses
// would be a door the page can never ring).
func TestWave4RosterNamesRouteToTheirOwnSockets(t *testing.T) {
	rules := &wave4Rules{}
	d := &ComposerDispatch{SessionRules: rules}
	rows := []struct {
		method string
		own    *int
		other  *int
	}{
		{MethodSessionRulesList, &rules.list, &rules.revoke},
		{MethodSessionRulesRevoke, &rules.revoke, &rules.list},
	}
	for _, r := range rows {
		before, otherBefore := *r.own, *r.other
		if _, err := d.Handle(context.Background(), wave4Envelope(r.method, "w4-"+r.method, nil)); err != nil {
			t.Fatalf("roster method %q was refused with its handler attached: %v", r.method, err)
		}
		if got := *r.own; got != before+1 {
			t.Errorf("method %q: own socket moved by %d, want 1", r.method, got-before)
		}
		if *r.other != otherBefore {
			t.Errorf("method %q moved its sibling method (%d -> %d)", r.method, otherBefore, *r.other)
		}
		if _, err := ParseComposerRequest(wave4Envelope(r.method, "w4-guard", nil)); err != nil {
			t.Errorf("the running whitelist guard refuses %q: %v", r.method, err)
		}
	}
}

// TestWave4SessionRulesHandlerDeclaresExactlyTheTwoVerbs pins the shape the
// pair joined the roster with: two methods, no third - a read list and a
// revoke, nothing else (an interface that grew an approval-shaped verb must go
// red here, where the interface lives).
func TestWave4SessionRulesHandlerDeclaresExactlyTheTwoVerbs(t *testing.T) {
	iface := reflect.TypeOf((*SessionRulesHandler)(nil)).Elem()
	if n := iface.NumMethod(); n != 2 {
		t.Errorf("SessionRulesHandler declares %d methods, want exactly 2 (list + revoke) - a third verb cannot join the roster through this interface", n)
	}
	var _ SessionRulesHandler = &wave4Rules{}
}

// TestWave4NilSocketRefusesByName keeps the fail-closed posture every other
// door carries: a whitelisted name on a chain with no leg attached refuses
// with ErrNoHandlerAttached, never drops and never fakes success.
func TestWave4NilSocketRefusesByName(t *testing.T) {
	d := &ComposerDispatch{}
	for _, method := range []string{MethodSessionRulesList, MethodSessionRulesRevoke} {
		_, err := d.Handle(context.Background(), wave4Envelope(method, "w4-nil", nil))
		if !errors.Is(err, ErrNoHandlerAttached) {
			t.Errorf("method %q with a nil socket answered %v, want an ErrNoHandlerAttached refusal", method, err)
		}
	}
}

// TestWave4RevokeSelectorSurvivesTheParser pins the revoke selector: the JSON
// key the page echoes (sessionRuleId) binds into ComposerRequest.SessionRuleID,
// and the value stays selector-shaped (a row id, never a key path).
func TestWave4RevokeSelectorSurvivesTheParser(t *testing.T) {
	req, err := ParseComposerRequest(wave4Envelope(MethodSessionRulesRevoke, "w4-sel", map[string]string{"sessionRuleId": "42"}))
	if err != nil {
		t.Fatalf("the revoke selector did not survive the parser: %v", err)
	}
	if req.SessionRuleID != "42" {
		t.Errorf("sessionRuleId = %q, want the echoed row id \"42\"", req.SessionRuleID)
	}
	if strings.HasPrefix(req.SessionRuleID, "$.") || strings.Contains(req.SessionRuleID, "..") {
		t.Errorf("selector value %q reads as a key path, not a selector", req.SessionRuleID)
	}
	// The list door takes NO selector: one filled on a list request is refused
	// by the cmd/wisp door (not by the parser, which is selector-blind) - this
	// pins only that the selector does not leak into the read's own contract.
	req, err = ParseComposerRequest(wave4Envelope(MethodSessionRulesList, "w4-sel2", nil))
	if err != nil {
		t.Fatalf("the list envelope did not survive the parser: %v", err)
	}
	if req.SessionRuleID != "" {
		t.Errorf("a selector-less list request carried %q", req.SessionRuleID)
	}
}

// TestWave4RouteSpellingsAreCleanUnderTheGrantVocabulary is the witness that
// the wave-4 spellings were chosen with the frozen instrument's vocabulary in
// hand: both names must be CLEAN under grantRouteWords (the instrument itself
// would already be red otherwise), and the spec's grants.list / grants.revoke
// must stay refused - the documented reason the synonym spellings exist. This
// check re-asks the vocabulary instead of trusting the constants: a rename
// that walks a banned word back onto the wire goes red HERE by name.
func TestWave4RouteSpellingsAreCleanUnderTheGrantVocabulary(t *testing.T) {
	for _, route := range []string{MethodSessionRulesList, MethodSessionRulesRevoke} {
		if carriesGrantWord(route, grantRouteWords) {
			t.Errorf("route %q carries a grantRouteWords word: the frozen grant-boundary instrument bans exactly this shape (AGENTS.md §1.2 ban #6)", route)
		}
	}
	for _, banned := range []string{"grants.list", "grants.revoke"} {
		if knownComposerMethod(banned) {
			t.Errorf("the guard answers %q - the spec spelling the grant-boundary instrument bans must never be admitted under its own name", banned)
		}
	}
}
