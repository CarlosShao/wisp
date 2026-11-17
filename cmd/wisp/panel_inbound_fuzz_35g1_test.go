//go:build windows

package main

// Ticket 35 :42 (Acceptance: "Whitelist fuzz: 50 random/forbidden method names ->
// all rejected + logged; each listed method dispatches") - G1's ruler for the
// inbound whitelist half.
//
// WHY THIS FILE LIVES IN cmd/wisp AND NOT internal/panel. The recon leg
// (docs/evidence/s1/35-recon-20261115.md, section (1)-c) proposed landing this in
// internal/panel/composer_dispatch_test.go. That package is READ-ONLY for this
// wave (another leg is using it), so the same assertions are made here against the
// SAME production type through its exported surface: panel.ComposerDispatch.Handle,
// panel.ParseComposerRequest's exported sentinels, and the exported handler
// interfaces. Nothing about the object under test changes; only the package the
// test binary is built in does.
//
// WHAT WAS MISSING BEFORE THIS FILE (the recon's reading, kept so a later reader
// can tell what this ruler actually added):
//   - The largest existing ruler enumerates 11 approval-shaped candidate names and
//     calls knownComposerMethod DIRECTLY - it never goes through Handle, so it
//     never reaches the audit line at composer_dispatch.go:298.
//   - "Each listed method dispatches" was only nailed for 3 of the roster names
//     (workspace/attachment/message); the other names (mode, the two settings
//     routes, the fetch route, the cost door, the three privacy doors - wave 1's
//     b1d2162dd - and readface wave 2's history/cancel doors) had no "the right
//     handler ran exactly once" nail until the 2026-11-17 re-adjudication grew
//     this table to the full thirteen. This file drives all 13.
//
// TWO TESTS, because the AC has two halves and they fail for different reasons:
//   1. TestWhitelistFuzz50NamesAllRefusedAndAudited35g1 - 50 forbidden names through
//      Handle: each is refused (error wraps ErrComposerRequest), the user gets a
//      non-empty sentence naming their requestId, NO handler runs, and the audit
//      sink carries EXACTLY 50 lines, each naming the method. Deleting the guard at
//      bridge.go:172-174 turns this red (the name would then reach a handler or the
//      default branch, and the audit count changes).
//   2. TestEveryRosterMethodDispatchesToItsOwnHandler35g1 - a table over all 13
//      roster names: each lands on ITS OWN handler exactly once and on no other.
//      Deleting a `case` from composer_dispatch.go's switch, or dropping a name
//      from bridge.go's knownComposerMethod, turns this red.

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/panel"
)

// ----------------------------------------------------------------------------
// the audit sink and the handler sockets, as spies
// ----------------------------------------------------------------------------

// g1Audit is the project's existing "[audit]" sink shape (workspace.go's
// AuditFunc): the dispatch chain writes one line per refusal here, and the ruler
// reads them back verbatim. A nil-safe slice, not a counter, because the AC says
// the 50 lines must be READABLE, not merely counted.
type g1Audit struct {
	lines []string
}

func (a *g1Audit) fn(format string, args ...any) {
	a.lines = append(a.lines, fmt.Sprintf(format, args...))
}

// g1Call records how often one socket was reached and with what. Every socket
// below increments its own counter, which is what lets the 8-name test assert
// "this door ran once and the other six ran zero times".
type g1Call struct {
	calls int
	last  panel.ComposerRequest
}

type g1Mode struct{ c g1Call }

func (s *g1Mode) HandleModeRequest(_ context.Context, req panel.ComposerRequest) error {
	s.c.calls++
	s.c.last = req
	return nil
}

type g1Workspace struct{ c g1Call }

func (s *g1Workspace) HandleWorkspaceRequest(_ context.Context, req panel.ComposerRequest) error {
	s.c.calls++
	s.c.last = req
	return nil
}

type g1Attachment struct{ c g1Call }

func (s *g1Attachment) HandleAttachmentRequest(_ context.Context, req panel.ComposerRequest) error {
	s.c.calls++
	s.c.last = req
	return nil
}

type g1Message struct{ c g1Call }

func (s *g1Message) HandleMessageRequest(_ context.Context, req panel.ComposerRequest) error {
	s.c.calls++
	s.c.last = req
	return nil
}

type g1Config struct{ c g1Call }

func (s *g1Config) HandleConfigRequest(_ context.Context, req panel.ComposerRequest, _ string) (string, error) {
	s.c.calls++
	s.c.last = req
	return "settings-applied", nil
}

type g1Models struct{ c g1Call }

func (s *g1Models) HandleModelsListRequest(_ context.Context, req panel.ComposerRequest) (string, error) {
	s.c.calls++
	s.c.last = req
	return "models-listed", nil
}

type g1Costs struct{ c g1Call }

func (s *g1Costs) HandleCostSummary(_ context.Context, req panel.ComposerRequest) (string, error) {
	s.c.calls++
	s.c.last = req
	return "cost-summarised", nil
}

// g1Privacy / g1History / g1RosterCancel are the wave-1/wave-2 sockets' spies
// (2026-11-17 re-adjudication): the dispatch ruler's table now drives all
// THIRTEEN roster names, so the three privacy verbs and the two readface
// wave-2 doors get the same "its own socket moved by exactly one" nail the
// eight 35g1-era names have. One call counter per privacy verb: a purge that
// reached delete-one's handler would be the wrong answer to a destructive
// request, which is exactly what this ruler exists to catch.
type g1Privacy struct {
	deleteOne g1Call
	purge     g1Call
	export    g1Call
}

func (s *g1Privacy) HandlePrivacyDeleteOne(_ context.Context, req panel.ComposerRequest) (string, error) {
	s.deleteOne.calls++
	s.deleteOne.last = req
	return "privacy-deleted-one", nil
}

func (s *g1Privacy) HandlePrivacyPurge(_ context.Context, req panel.ComposerRequest) (string, error) {
	s.purge.calls++
	s.purge.last = req
	return "privacy-purged", nil
}

func (s *g1Privacy) HandlePrivacyExport(_ context.Context, req panel.ComposerRequest) (string, error) {
	s.export.calls++
	s.export.last = req
	return "privacy-exported", nil
}

type g1History struct{ c g1Call }

func (s *g1History) HandleHistoryPage(_ context.Context, req panel.ComposerRequest) (string, error) {
	s.c.calls++
	s.c.last = req
	return "history-paged", nil
}

type g1RosterCancel struct{ c g1Call }

func (s *g1RosterCancel) HandleRosterCancel(_ context.Context, req panel.ComposerRequest) (string, error) {
	s.c.calls++
	s.c.last = req
	return "roster-cancelled", nil
}

// The wave-3 sockets (2026-11-17, ticket 38 AC#3/AC#5): same spy shape, one
// handler each, so the dispatch ruler covers the grown roster without
// loosening anything.
type g1ToolsList struct{ c g1Call }

func (s *g1ToolsList) HandleToolsList(_ context.Context, req panel.ComposerRequest) (string, error) {
	s.c.calls++
	s.c.last = req
	return "tools-listed", nil
}

type g1ToolLaunch struct{ c g1Call }

func (s *g1ToolLaunch) HandleToolLaunch(_ context.Context, req panel.ComposerRequest) (string, error) {
	s.c.calls++
	s.c.last = req
	return "tool-launched", nil
}

type g1RulesForce struct{ c g1Call }

func (s *g1RulesForce) HandleRulesForce(_ context.Context, req panel.ComposerRequest) (string, error) {
	s.c.calls++
	s.c.last = req
	return "rules-forced", nil
}

// g1Envelope marshals one composer envelope with the attribution fields the two
// later guards require, so the ONLY refusal any name in this file can draw is the
// method-whitelist guard - never source and never requestId.
func g1Envelope(method, requestID string) string {
	m := map[string]any{
		"method":    method,
		"requestId": requestID,
		"source":    panel.ComposerRequestSource,
	}
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// ----------------------------------------------------------------------------
// half 1: 50 forbidden names, all refused and all logged
// ----------------------------------------------------------------------------

// g1ForbiddenNames is EXACTLY 50 names, none of which is one of the 8 roster
// names, and none of which TrimSpace's into a roster name (the guard trims before
// comparing, so a trailing space would be ACCEPTED and would make this ruler lie).
//
// The list deliberately mixes five shapes, so a guard that only knows one of them
// cannot pass:
//   - names a real page or a neighbouring ticket sends (approval.decide, the
//     panel.approval.request envelope, panel.resync, tasks.list, ...);
//   - near-misses of roster names (a suffix, a case change, a segment swap);
//   - namespaced names in a different namespace (models.delete, grants.revoke);
//   - shape-only noise that still passes the inbound dot-shape test, so it reaches
//     the guard rather than being lost at JSON parse;
//   - the empty method, which must be refused by the same door ("" is not on the
//     roster; TrimSpace("") stays "").
func g1ForbiddenNames() []string {
	names := []string{
		// real / neighbouring names the tree itself names
		"approval.decide",
		"panel.approval.request",
		"panel.resync",
		"tasks.list",
		"tasks.detail",
		"history.query",
		"transcript.get",
		"approval.current",
		"approval.queue",
		"grants.list",
		"grants.revoke",
		// 2026-11-17 re-adjudication (readface wave 2's full-package carrier
		// run caught it): privacy.purge / privacy.export joined the whitelist
		// with readface wave 1 (b1d2162dd, ticket 40 AC#1) and this table was
		// not moved with it - a name the whitelist answers is by definition
		// not a fuzz-table member. The two join as fresh near-miss shapes so
		// the table stays exactly 50; history.page / roster.cancel (wave 2)
		// were minted off-table, so nothing else moves.
		"privacy.delete-all",
		"privacy.export.all",
		"models.delete",
		"diagnostics.export",
		// near-misses of the eight roster names
		"panel.mode.request2",
		"panel.mode.Request",
		"panel.mode.request.extra",
		"mode.request",
		"config.get2",
		"config.Get",
		"config.get.extra",
		"config.set2",
		"config.del",
		"config.models.list2",
		"config.models.delete",
		"cost.summary2",
		"cost.summary.all",
		"panel.workspace.requestx",
		"panel.attachment.add2",
		"panel.message.send2",
		// other namespaces / random-looking method names
		"auth.login",
		"auth.logout",
		"session.reset",
		"fs.read",
		"fs.write",
		"exec.shell",
		"net.fetch",
		"db.query",
		"admin.debug",
		"internal.dump",
		"tool.invoke",
		"plugin.load",
		"plugin.unload",
		"eval.script",
		// shape-only noise
		"panel.",
		".hidden",
		"a.b",
		"x.y.z",
		"fuzz.method.1",
	}
	return names
}

func TestWhitelistFuzz50NamesAllRefusedAndAudited35g1(t *testing.T) {
	names := g1ForbiddenNames()
	if len(names) != 50 {
		t.Fatalf("the fuzz table holds %d names, want exactly 50: ticket 35 :42 asks for 50 random/forbidden names, and a smaller table would silently lower the bar", len(names))
	}

	// A roster must exist and be the closed sixteen, or "all rejected" is
	// trivially true against an empty whitelist. This is the anti-empty-universe
	// guard the closed-roster ruler in internal/panel also carries. 2026-11-17
	// re-adjudication (readface wave 2's carrier run): the slice was pinned at
	// the 35g1-era eight while readface wave 1 (b1d2162dd) had already grown
	// the whitelist to eleven with the three privacy doors; wave 2 added the
	// two doors this package wires (thirteen), and readface wave 3 adds
	// tools.list / tool.launch / rules.force (sixteen).
	roster := []string{
		panel.MethodModeRequest, panel.MethodWorkspaceRequest, panel.MethodAttachmentAdd,
		panel.MethodMessageSend, panel.MethodConfigGet, panel.MethodConfigSet,
		panel.MethodModelsList, panel.MethodCostSummary,
		panel.MethodPrivacyDeleteOne, panel.MethodPrivacyPurge, panel.MethodPrivacyExport,
		panel.MethodHistoryPage, panel.MethodRosterCancel,
		panel.MethodToolsList, panel.MethodToolLaunch, panel.MethodRulesForce,
	}
	if len(roster) != 16 {
		t.Fatalf("the roster slice here has %d names, want 16 - the count this ruler reports against", len(roster))
	}
	onRoster := map[string]bool{}
	for _, r := range roster {
		onRoster[r] = true
	}
	for i, name := range names {
		if onRoster[name] {
			t.Fatalf("fuzz name #%d %q is a ROSTER name - the table must hold forbidden names only, or this test is asserting the whitelist refuses its own entries", i, name)
		}
	}

	mode := &g1Mode{}
	audit := &g1Audit{}
	d := &panel.ComposerDispatch{Mode: mode, Audit: audit.fn}

	for i, name := range names {
		rid := fmt.Sprintf("fuzz-rid-%02d", i)
		raw := g1Envelope(name, rid)
		reply, err := d.Handle(context.Background(), raw)
		if err == nil {
			t.Fatalf("forbidden method %q (#%d) was ACCEPTED (no error) - a name off the roster must never reach a handler", name, i)
		}
		if !isComposerRequestRefusal(err) {
			t.Fatalf("forbidden method %q (#%d) was refused by the WRONG door: %v - want an error wrapping panel.ErrComposerRequest (the whitelist door), so a source/requestId refusal cannot masquerade as a whitelist refusal here", name, i, err)
		}
		if reply == "" {
			t.Fatalf("forbidden method %q (#%d) produced the empty sentence: a dropped request that looks handled", name, i)
		}
		if !strings.Contains(reply, rid) {
			t.Errorf("forbidden method %q (#%d): the user-facing sentence %q does not name requestId %q, so the user cannot trace their own click", name, i, reply, rid)
		}
		if mode.c.calls != 0 {
			t.Fatalf("forbidden method %q (#%d) reached the mode handler %d time(s), want 0", name, i, mode.c.calls)
		}
	}

	// The audit half: exactly one line per refusal, 50 total, and EACH line names
	// its own method. Counting alone would pass a ruler that logged a single line
	// fifty times; naming the method is what makes the log readable per request.
	if len(audit.lines) != 50 {
		t.Fatalf("audit lines = %d, want exactly 50 (one per refused name): a refusal with no line is a silent drop wearing a return value, and a wrong count means a name was not logged or was logged twice", len(audit.lines))
	}
	for i, name := range names {
		want := "method=" + strconv.Quote(name)
		if !strings.Contains(audit.lines[i], want) {
			t.Errorf("audit line #%d does not name its method %s:\n  got: %s", i, want, audit.lines[i])
		}
		if !strings.Contains(audit.lines[i], "INBOUND-DISPATCH") {
			t.Errorf("audit line #%d is not the inbound-dispatch record: %s", i, audit.lines[i])
		}
	}

	// The positive control: a name that IS on the roster must draw no refusal from
	// the whitelist guard, so this ruler is not simply asserting "Handle always
	// errors". A roster name with its handler attached must be accepted.
	mode.c.calls = 0
	if _, err := d.Handle(context.Background(), g1Envelope(panel.MethodModeRequest, "positive-control-rid")); err != nil {
		t.Fatalf("positive control: a ROSTER name (%q) with its handler attached was refused: %v - a ruler that cannot tell a valid call from a forged one has no discriminating power", panel.MethodModeRequest, err)
	}
	if mode.c.calls != 1 {
		t.Fatalf("positive control: the roster name reached the mode handler %d time(s), want 1", mode.c.calls)
	}
	t.Logf("whitelist fuzz: %d/%d forbidden names refused by the ErrComposerRequest door, %d audit line(s), each naming its method; positive control accepted", len(names), len(names), len(audit.lines))
}

// isComposerRequestRefusal is errors.Is(err, panel.ErrComposerRequest) without
// importing errors at the call site; kept as a named helper so the failure text
// above can say which door (rather than "some error") refused.
func isComposerRequestRefusal(err error) bool {
	for e := err; e != nil; {
		if e == panel.ErrComposerRequest {
			return true
		}
		type unwrapper interface{ Unwrap() error }
		u, ok := e.(unwrapper)
		if !ok {
			return false
		}
		e = u.Unwrap()
	}
	return false
}

// ----------------------------------------------------------------------------
// half 2: each of the 8 roster names dispatches to ITS OWN handler
// ----------------------------------------------------------------------------

func TestEveryRosterMethodDispatchesToItsOwnHandler35g1(t *testing.T) {
	mode := &g1Mode{}
	workspace := &g1Workspace{}
	attachment := &g1Attachment{}
	message := &g1Message{}
	config := &g1Config{}
	models := &g1Models{}
	costs := &g1Costs{}
	privacy := &g1Privacy{}
	history := &g1History{}
	rosterCancel := &g1RosterCancel{}
	tools := &g1ToolsList{}
	launch := &g1ToolLaunch{}
	rules := &g1RulesForce{}

	d := &panel.ComposerDispatch{
		Mode:         mode,
		Workspace:    workspace,
		Attachment:   attachment,
		Message:      message,
		Config:       config,
		Models:       models,
		Costs:        costs,
		Privacy:      privacy,
		History:      history,
		RosterCancel: rosterCancel,
		Tools:        tools,
		Launch:       launch,
		Rules:        rules,
		Audit:        (&g1Audit{}).fn,
	}

	// each row names the ONE socket that must fire, by key into the all map below.
	type row struct {
		method string
	}
	rows := []row{
		{panel.MethodModeRequest},
		{panel.MethodWorkspaceRequest},
		{panel.MethodAttachmentAdd},
		{panel.MethodMessageSend},
		{panel.MethodConfigGet},
		{panel.MethodConfigSet},
		{panel.MethodModelsList},
		{panel.MethodCostSummary},
		{panel.MethodPrivacyDeleteOne},
		{panel.MethodPrivacyPurge},
		{panel.MethodPrivacyExport},
		{panel.MethodHistoryPage},
		{panel.MethodRosterCancel},
		{panel.MethodToolsList},
		{panel.MethodToolLaunch},
		{panel.MethodRulesForce},
	}
	if len(rows) != 16 {
		t.Fatalf("the dispatch table covers %d roster names, want all 16", len(rows))
	}

	all := map[string]func() int{
		"mode":               func() int { return mode.c.calls },
		"workspace":          func() int { return workspace.c.calls },
		"attachment":         func() int { return attachment.c.calls },
		"message":            func() int { return message.c.calls },
		"config":             func() int { return config.c.calls },
		"models":             func() int { return models.c.calls },
		"costs":              func() int { return costs.c.calls },
		"privacy.delete-one": func() int { return privacy.deleteOne.calls },
		"privacy.purge":      func() int { return privacy.purge.calls },
		"privacy.export":     func() int { return privacy.export.calls },
		"history":            func() int { return history.c.calls },
		"cancel":             func() int { return rosterCancel.c.calls },
		"tools.list":         func() int { return tools.c.calls },
		"tool.launch":        func() int { return launch.c.calls },
		"rules.force":        func() int { return rules.c.calls },
	}

	dispatched := 0
	for i, r := range rows {
		rid := fmt.Sprintf("roster-rid-%02d", i)
		before := map[string]int{}
		for k, get := range all {
			before[k] = get()
		}
		_, err := d.Handle(context.Background(), g1Envelope(r.method, rid))
		if err != nil {
			t.Fatalf("roster method %q (#%d) was refused with every handler attached: %v - a listed method that cannot dispatch is a page-visible dead door", r.method, i, err)
		}
		// the owning socket moved by exactly one.
		var ownKey string
		switch r.method {
		case panel.MethodModeRequest:
			ownKey = "mode"
		case panel.MethodWorkspaceRequest:
			ownKey = "workspace"
		case panel.MethodAttachmentAdd:
			ownKey = "attachment"
		case panel.MethodMessageSend:
			ownKey = "message"
		case panel.MethodConfigGet, panel.MethodConfigSet:
			ownKey = "config"
		case panel.MethodModelsList:
			ownKey = "models"
		case panel.MethodCostSummary:
			ownKey = "costs"
		case panel.MethodPrivacyDeleteOne:
			ownKey = "privacy.delete-one"
		case panel.MethodPrivacyPurge:
			ownKey = "privacy.purge"
		case panel.MethodPrivacyExport:
			ownKey = "privacy.export"
		case panel.MethodHistoryPage:
			ownKey = "history"
		case panel.MethodRosterCancel:
			ownKey = "cancel"
		case panel.MethodToolsList:
			ownKey = "tools.list"
		case panel.MethodToolLaunch:
			ownKey = "tool.launch"
		case panel.MethodRulesForce:
			ownKey = "rules.force"
		}
		for k, get := range all {
			delta := get() - before[k]
			want := 0
			if k == ownKey {
				want = 1
			}
			if delta != want {
				t.Errorf("roster method %q (#%d): socket %q moved by %d, want %d (a listed method must reach exactly its own handler and no other - a misroute is a wrong answer to a real request)", r.method, i, k, delta, want)
			}
		}
		dispatched++
	}

	if dispatched != 16 {
		t.Fatalf("dispatched %d roster names, want 16", dispatched)
	}

	// Every roster name is also answered by the RUNNING guard. This is the other
	// direction of the closed roster: a name the dispatch table handles but the
	// guard would refuse is a name that can never be reached from the page. The
	// guard is asked through the production parser (its first gate), never
	// re-implemented here.
	for _, r := range rows {
		if _, err := panel.ParseComposerRequest(g1Envelope(r.method, "guard-check")); err != nil {
			t.Errorf("roster method %q is handled by the dispatch table but the running whitelist guard refuses it: %v", r.method, err)
		}
	}
	t.Logf("roster dispatch: %d/%d roster names each reached their own handler exactly once and no other socket moved", dispatched, 16)
}
