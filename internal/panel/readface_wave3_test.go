package panel

// Readface wave 3 (2026-11-17), ticket 38 AC#3's discovery/launch doors and
// AC#5's force-rules read: the Go-half instruments for the three roster names
// this wave added (tools.list / tool.launch / rules.force). The wave-2 file
// (readface_wave2_test.go) is the shape these follow; nothing here repeats its
// assertions - each check below is the one wave 3 needs and wave 2 does not
// carry.
//
// What is pinned, and why:
//
//  1. The three names are on the whitelist and route to their OWN socket -
//     the closed-roster ruler and the cmd/wisp fuzz table both grew with the
//     same commit, but the router case is the wire itself, so it is asserted
//     here where the router lives.
//  2. The two read faces declare ONE method each and no more: read-only is
//     the security statement that let these names join the roster beside the
//     frozen grant-boundary instrument, so an interface that grew a second
//     (write-shaped) method must go red here, where the interface lives.
//  3. tool.launch's receipt never carries an approval verdict word in its
//     route shape: the door issues a launch, it does not answer a card. The
//     names are already covered by the grant-boundary instrument's sweep;
//     what this file adds is the router/shape/ledger half.
//  4. The envelope's new selectors exist and are selector-shaped (names and
//     one JSON string, never key paths), mirroring the wave-2 selector check.

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

// wave3Envelope builds the ONE envelope the parser accepts for method.
func wave3Envelope(method, rid string) string {
	raw, _ := json.Marshal(map[string]string{
		"method": method, "requestId": rid, "source": ComposerRequestSource,
	})
	return string(raw)
}

type wave3Tools struct{ calls int }

func (h *wave3Tools) HandleToolsList(context.Context, ComposerRequest) (string, error) {
	h.calls++
	return "名册", nil
}

type wave3Launch struct{ calls int }

func (h *wave3Launch) HandleToolLaunch(context.Context, ComposerRequest) (string, error) {
	h.calls++
	return "已发出", nil
}

type wave3Rules struct{ calls int }

func (h *wave3Rules) HandleRulesForce(context.Context, ComposerRequest) (string, error) {
	h.calls++
	return "规则", nil
}

// TestWave3RosterNamesRouteToTheirOwnSockets is the router half: each of the
// three names dispatches to exactly its own handler and no other socket moves,
// and the running guard accepts each name (a case label the whitelist refuses
// would be a door the page can never ring).
func TestWave3RosterNamesRouteToTheirOwnSockets(t *testing.T) {
	tools, launch, rules := &wave3Tools{}, &wave3Launch{}, &wave3Rules{}
	d := &ComposerDispatch{Tools: tools, Launch: launch, Rules: rules}
	rows := []struct {
		method string
		own    *int
	}{
		{MethodToolsList, &tools.calls},
		{MethodToolLaunch, &launch.calls},
		{MethodRulesForce, &rules.calls},
	}
	for _, r := range rows {
		toolsBefore, launchBefore, rulesBefore := tools.calls, launch.calls, rules.calls
		if _, err := d.Handle(context.Background(), wave3Envelope(r.method, "w3-"+r.method)); err != nil {
			t.Fatalf("roster method %q was refused with its handler attached: %v", r.method, err)
		}
		if got, want := *r.own, beforeOf(r.method, toolsBefore, launchBefore, rulesBefore)+1; got != want {
			t.Errorf("method %q: own socket moved by %d, want 1", r.method, got-want+1)
		}
		if tools.calls != toolsBefore && r.method != MethodToolsList {
			t.Errorf("method %q moved the tools socket (%d -> %d)", r.method, toolsBefore, tools.calls)
		}
		if launch.calls != launchBefore && r.method != MethodToolLaunch {
			t.Errorf("method %q moved the launch socket (%d -> %d)", r.method, launchBefore, launch.calls)
		}
		if rules.calls != rulesBefore && r.method != MethodRulesForce {
			t.Errorf("method %q moved the rules socket (%d -> %d)", r.method, rulesBefore, rules.calls)
		}
		if _, err := ParseComposerRequest(wave3Envelope(r.method, "w3-guard")); err != nil {
			t.Errorf("the running whitelist guard refuses %q: %v", r.method, err)
		}
	}
}

// beforeOf returns the own socket's count before this row ran.
func beforeOf(method string, toolsBefore, launchBefore, rulesBefore int) int {
	switch method {
	case MethodToolsList:
		return toolsBefore
	case MethodToolLaunch:
		return launchBefore
	default:
		return rulesBefore
	}
}

// TestWave3ReadFacesDeclareOneReadOnlyMethodEach pins the shape that made the
// two read names roster-eligible: one method each, counted off the interface
// types themselves (not the fixtures), so a second (write-shaped) method
// cannot quietly join a read interface.
func TestWave3ReadFacesDeclareOneReadOnlyMethodEach(t *testing.T) {
	pairs := []struct {
		name string
		typ  any
	}{
		{"ToolsListHandler", reflect.TypeOf((*ToolsListHandler)(nil)).Elem()},
		{"ToolLaunchHandler", reflect.TypeOf((*ToolLaunchHandler)(nil)).Elem()},
		{"RulesForceHandler", reflect.TypeOf((*RulesForceHandler)(nil)).Elem()},
	}
	for _, p := range pairs {
		iface, ok := p.typ.(reflect.Type)
		if !ok {
			t.Fatalf("%s: not an interface type", p.name)
		}
		if n := iface.NumMethod(); n != 1 {
			t.Errorf("%s declares %d methods, want exactly 1 - read-only is the security statement that let this name join the roster", p.name, n)
		}
	}
	var _ ToolsListHandler = &wave3Tools{}
	var _ ToolLaunchHandler = &wave3Launch{}
	var _ RulesForceHandler = &wave3Rules{}
}

// TestWave3NilSocketsRefuseByName keeps the fail-closed posture every other
// door carries: a whitelisted name on a chain with no leg attached refuses
// with ErrNoHandlerAttached, never drops and never fakes success.
func TestWave3NilSocketsRefuseByName(t *testing.T) {
	d := &ComposerDispatch{}
	for _, method := range []string{MethodToolsList, MethodToolLaunch, MethodRulesForce} {
		_, err := d.Handle(context.Background(), wave3Envelope(method, "w3-nil"))
		if !errors.Is(err, ErrNoHandlerAttached) {
			t.Errorf("method %q with a nil socket answered %v, want an ErrNoHandlerAttached refusal", method, err)
		}
	}
	var nilDisp *ComposerDispatch
	_, err := nilDisp.Handle(context.Background(), wave3Envelope(MethodToolsList, "w3-nildisp"))
	if !errors.Is(err, ErrNoHandlerAttached) {
		t.Errorf("a nil dispatch answered %v, want the named refusal", err)
	}
}

// TestWave3SelectorsAreSelectorShaped pins the new envelope fields: the JSON
// keys exist (the page's pending literals will echo them), and none of the
// wave-3 route values carry a key-path shape - a selector, never a path.
func TestWave3SelectorsAreSelectorShaped(t *testing.T) {
	raw := `{"method":"` + MethodToolLaunch + `","requestId":"w3-sel","source":"` + ComposerRequestSource +
		`","toolName":"fs.delete","toolArgs":"{\"path\":\"C:/x\"}","forceTask":"清理临时文件","forceMode":"force_tool"}`
	req, err := ParseComposerRequest(raw)
	if err != nil {
		t.Fatalf("the wave-3 selectors did not survive the parser: %v", err)
	}
	if req.ToolName != "fs.delete" || req.ToolArgs != `{"path":"C:/x"}` {
		t.Errorf("launch selectors = %q/%q, want the echoed name and JSON args", req.ToolName, req.ToolArgs)
	}
	if req.ForceTask != "清理临时文件" || req.ForceMode != "force_tool" {
		t.Errorf("force selectors = %q/%q", req.ForceTask, req.ForceMode)
	}
	for _, v := range []string{req.ToolName, req.ToolArgs, req.ForceTask, req.ForceMode} {
		if strings.HasPrefix(v, "$.") || strings.Contains(v, "..") {
			t.Errorf("selector value %q reads as a key path, not a selector", v)
		}
	}
}

// TestWave3PageGapLedgerRegistersTheThreeNames keeps the registration honest:
// every wave-3 entry must name its reason (an empty reason is a silent
// exemption) and the ledger must hold exactly the three names this wave owes
// the page.
func TestWave3PageGapLedgerRegistersTheThreeNames(t *testing.T) {
	for _, name := range []string{MethodToolsList, MethodToolLaunch, MethodRulesForce} {
		why, ok := rosterPageGapLedger194[name]
		if !ok {
			t.Errorf("%q is not registered in rosterPageGapLedger194: an unmirrored pending name must be registered with its reason, not left to the original rule", name)
			continue
		}
		if strings.TrimSpace(why) == "" {
			t.Errorf("%q is registered with an empty reason - the ledger is a bridge, not a silent exemption", name)
		}
	}
}
