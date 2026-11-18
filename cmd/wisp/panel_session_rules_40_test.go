package main

// Readface wave 4 (2026-11-18), ticket 40 AC#2's teeth: the security page's
// grants view reads 49's grants LIVE and the panel revoke bites on the very
// NEXT covered call. Three tests, each measuring a behavior the run really
// performed - never a source string:
//
//	Test40SessionRulesWiringNailInProductionChain   the production chain
//	                                                carries the real door; an
//	                                                unbound seam and a bad
//	                                                selector refuse by name
//	                                                (zero effect refusals).
//	Test40SessionRulesReadFaceReflectsTicket49GrantsLive
//	                                                one assembled run: the
//	                                                empty ledger is said OUT
//	                                                LOUD, two real rows come
//	                                                back field by field, and
//	                                                after a panel revoke the
//	                                                spent row LEAVES the live
//	                                                list while its audited
//	                                                row stays on disk with
//	                                                revoked_at stamped.
//	Test40SessionRulesPanelRevokePopsTheNextCoveredCall
//	                                                the AC#2 integration
//	                                                fixture: grant live =>
//	                                                covered call silent
//	                                                (0 cards), the panel
//	                                                revoke receipt names the
//	                                                row, the SAME shaped call
//	                                                really re-asks (>= 1 card
//	                                                on the production UI
//	                                                counter) - the two
//	                                                readings the dispatch
//	                                                asks for.
//
// The judgement ceiling ticket49_verbs_test.go states applies here too: every
// cell is ONE ASSEMBLED RUN in the test process - the real wisp.db, the real
// ledger, the real bridge, the real gate - and the card counts are read off
// rt.windowCount() (the production UI's own shown() counter), never off a
// mock. The panel route is driven through panel.ComposerDispatch.Handle with
// the production chain assembled by newComposerDispatchChain and its seam
// bound the way the boot binds it (the runtime's own ledger) - the same
// objects production uses, no parallel leg.

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/panel"
)

func sessionRulesEnvelope40(requestID, method, ruleID string) string {
	env := `{"method":"` + method + `","requestId":"` + requestID + `","source":"` + panel.ComposerRequestSource + `"`
	if ruleID != "" {
		env += `,"sessionRuleId":"` + ruleID + `"`
	}
	return env + `}`
}

// sessionRulesPayloadView40 mirrors the door's wire shape
// (panel_session_rules_40.go sessionRulesPayload40) for the teeth below.
type sessionRulesPayloadView40 struct {
	Session string `json:"session"`
	Rows    []struct {
		ID        int64  `json:"id"`
		Scope     string `json:"scope"`
		Tool      string `json:"tool"`
		Pattern   string `json:"pattern"`
		SessionID string `json:"session"`
		CreatedAt int64  `json:"createdAt"`
		ExpiresAt int64  `json:"expiresAt"`
	} `json:"rows"`
	Zero string `json:"zero"`
}

// parseSessionRulesPayload40 pulls the ONE payload line out of a list reply
// the way the page does (find the marker, parse the rest of the line), so
// every claim below is read off the wire the page reads - never off the prose.
func parseSessionRulesPayload40(t *testing.T, reply string) sessionRulesPayloadView40 {
	t.Helper()
	var payloadLine string
	for _, line := range strings.Split(reply, "\n") {
		if strings.HasPrefix(line, SessionRulesPayloadMarker) {
			payloadLine = strings.TrimPrefix(line, SessionRulesPayloadMarker)
		}
	}
	if payloadLine == "" {
		t.Fatalf("list reply carries no %s line: %q", SessionRulesPayloadMarker, reply)
	}
	var p sessionRulesPayloadView40
	if err := json.Unmarshal([]byte(payloadLine), &p); err != nil {
		t.Fatalf("list payload is not JSON: %v\n%s", err, payloadLine)
	}
	return p
}

// ---------------------------------------------------------------------------
// wiring nail + selector refusals (zero effect on every refusal)
// ---------------------------------------------------------------------------

func Test40SessionRulesWiringNailInProductionChain(t *testing.T) {
	dir := newInboundDataDir33(t, "ask_every_step")
	disp, _, err := newComposerDispatchChain(dir, auditQuiet194, panelInboundActor, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("the production chain did not assemble on %s: %v", dir, err)
	}
	if disp.SessionRules == nil {
		t.Fatal("the production chain carries no session-rules leg: the two doors are back on the nil socket that refuses by name")
	}
	door, ok := disp.SessionRules.(*sessionRulesDoor40)
	if !ok {
		t.Fatalf("session-rules leg is %T, want the real door sessionRulesDoor40 (a stub here would be the non-真 answer)", disp.SessionRules)
	}
	if door.seamOf().list != nil {
		t.Fatal("the CLI chain's seam is bound before any runtime exists - the late-bound pattern is broken")
	}
	ctx := context.Background()
	// Unbound seam: both doors refuse BY NAME, with zero effect.
	_, err = disp.Handle(ctx, sessionRulesEnvelope40("rid-sr-nail-list", panel.MethodSessionRulesList, ""))
	if err == nil || !strings.Contains(err.Error(), "会话授权读面未接入") {
		t.Errorf("unbound list = %v, want the named 会话授权读面未接入 refusal", err)
	}
	_, err = disp.Handle(ctx, sessionRulesEnvelope40("rid-sr-nail-rev", panel.MethodSessionRulesRevoke, "42"))
	if err == nil || !strings.Contains(err.Error(), "会话授权撤销通路未接入") {
		t.Errorf("unbound revoke = %v, want the named 会话授权撤销通路未接入 refusal", err)
	}
	// Bad selectors refuse BEFORE the seam is reached (the door is nil here,
	// so reaching the seam would panic - the refusals below prove the
	// selector gates run first).
	_, err = disp.Handle(ctx, sessionRulesEnvelope40("rid-sr-noid", panel.MethodSessionRulesRevoke, ""))
	if err == nil || !strings.Contains(err.Error(), "sessionRuleId 是空的") {
		t.Errorf("empty-selector refusal = %v, want one that names the missing 编号", err)
	}
	_, err = disp.Handle(ctx, sessionRulesEnvelope40("rid-sr-badid", panel.MethodSessionRulesRevoke, "abc"))
	if err == nil || !strings.Contains(err.Error(), "不是读面列出的授权编号") {
		t.Errorf("non-numeric-selector refusal = %v, want one that names the 编号 rule", err)
	}
	_, err = disp.Handle(ctx, sessionRulesEnvelope40("rid-sr-zero", panel.MethodSessionRulesRevoke, "0"))
	if err == nil {
		t.Errorf("sessionRuleId=0 was admitted; want the 编号 rule's refusal")
	}
}

// ---------------------------------------------------------------------------
// the read face reflects 49's grants live (one assembled run)
// ---------------------------------------------------------------------------

func Test40SessionRulesReadFaceReflectsTicket49GrantsLive(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	h.reply = strings.NewReader("")
	note := filepath.Join(h.dir, "rules-note.txt")
	ctx := context.Background()

	var (
		minted              string
		idPath              int64
		idTarget            int64
		listPayload         sessionRulesPayloadView40
		emptyZero           string
		afterRows           int
		targetRevokedOnDisk bool
		pathLiveOnDisk      bool
		spentGoneFromList   bool
	)

	h.rtHook = func(rt *agentRuntime) {
		if rt.session == nil {
			t.Fatalf("assembled run has no session ledger")
		}
		minted = rt.session.SessionID().String()
		// The production chain, bound the way the boot binds it: the seam is
		// the runtime's own ledger - the same rows the gate asks.
		disp, _, err := newComposerDispatchChain(h.dir, auditQuiet194, "test-panel", nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("the production chain did not assemble on %s: %v", h.dir, err)
		}
		door, ok := disp.SessionRules.(*sessionRulesDoor40)
		if !ok {
			t.Fatalf("session-rules leg is %T, want the real door", disp.SessionRules)
		}
		door.bindSeam(rt.sessionRulesSeam40())

		// EMPTY LEDGER: the negative control. The read face must say zero
		// rows OUT LOUD in the payload's own zero field - not silence, not
		// an error, and not an empty table dressed up as a reading.
		reply, err := disp.Handle(ctx, sessionRulesEnvelope40("rid-sr-empty", panel.MethodSessionRulesList, ""))
		if err != nil {
			t.Fatalf("list on an empty ledger: %v", err)
		}
		empty := parseSessionRulesPayload40(t, reply)
		emptyZero = empty.Zero
		if len(empty.Rows) != 0 {
			t.Fatalf("an empty ledger listed %d rows", len(empty.Rows))
		}
		if empty.Session != minted {
			t.Fatalf("the empty payload names session %q, want %s", empty.Session, minted)
		}

		// Two REAL rows through the ledger's two write doors - the same pair
		// ticket 49's tooth mints, so the view is measured against the same
		// shape the console verb lists.
		c, err := rt.paths.Canonicalize(filepath.ToSlash(note))
		if err != nil {
			t.Fatalf("Canonicalize: %v", err)
		}
		idPath, err = rt.session.Record(ctx, "fs.write", c)
		if err != nil {
			t.Fatalf("session.Record: %v", err)
		}
		idTarget, err = rt.session.RecordTarget(ctx, "input.type", "notepad.exe")
		if err != nil {
			t.Fatalf("session.RecordTarget: %v", err)
		}

		// THE LIVE LIST: both rows, field by field against the rows the
		// store really holds at this moment.
		reply, err = disp.Handle(ctx, sessionRulesEnvelope40("rid-sr-live", panel.MethodSessionRulesList, ""))
		if err != nil {
			t.Fatalf("list with two rows: %v", err)
		}
		listPayload = parseSessionRulesPayload40(t, reply)
		if len(listPayload.Rows) != 2 {
			t.Fatalf("the live list holds %d rows, want exactly 2 (payload=%+v)", len(listPayload.Rows), listPayload)
		}
		if listPayload.Zero != "" {
			t.Fatalf("a two-row list carries a zero sentence %q", listPayload.Zero)
		}
		mid, err := rt.session.ListBySession(ctx)
		if err != nil {
			t.Fatalf("ListBySession: %v", err)
		}
		if len(mid) != 2 {
			t.Fatalf("the ledger holds %d rows mid-run, want exactly 2 (rows=%+v)", len(mid), mid)
		}
		for _, got := range listPayload.Rows {
			var store *struct {
				id        int64
				scope     string
				tool      string
				pattern   string
				sessionID string
				createdAt int64
				expiresAt int64
			}
			for i := range mid {
				if mid[i].ID == got.ID {
					g := mid[i]
					store = &struct {
						id        int64
						scope     string
						tool      string
						pattern   string
						sessionID string
						createdAt int64
						expiresAt int64
					}{g.ID, g.Scope, g.Tool, g.Pattern, g.SessionID, g.CreatedAt, g.ExpiresAt}
				}
			}
			if store == nil {
				t.Errorf("the live list carries row %d the store does not hold", got.ID)
				continue
			}
			if got.Scope != store.scope || got.Tool != store.tool || got.Pattern != store.pattern ||
				got.SessionID != store.sessionID || got.CreatedAt != store.createdAt || got.ExpiresAt != store.expiresAt {
				t.Errorf("row %d wire fields = %+v, store row = %+v: the payload is not the store's own fields",
					got.ID, got, store)
			}
			if got.Scope != "session" {
				t.Errorf("row %d scope = %q, want the store's session scope word", got.ID, got.Scope)
			}
		}

		// THE PANEL REVOKE: the target row spent through the production
		// dispatch - the same door the security page's one-click revoke
		// rings (the seam is the ledger's own Revoke entrance, the door
		// ticket 49's console verb spends).
		reply, err = disp.Handle(ctx, sessionRulesEnvelope40("rid-sr-revoke", panel.MethodSessionRulesRevoke, fmt.Sprint(idTarget)))
		if err != nil {
			t.Fatalf("revoke %d through the panel route: %v", idTarget, err)
		}
		if !strings.Contains(reply, fmt.Sprintf("已撤销会话授权 #%d", idTarget)) {
			t.Errorf("the revoke receipt %q does not name row #%d", reply, idTarget)
		}

		// THE LIVE LIST AFTER THE REVOKE: the spent row LEAVES the live list
		// (tooth ③'s green half), the untouched row stays, and the payload's
		// zero field stays empty - one live row is not an empty ledger.
		reply, err = disp.Handle(ctx, sessionRulesEnvelope40("rid-sr-after", panel.MethodSessionRulesList, ""))
		if err != nil {
			t.Fatalf("list after the revoke: %v", err)
		}
		after := parseSessionRulesPayload40(t, reply)
		afterRows = len(after.Rows)
		spentGoneFromList = true
		for _, got := range after.Rows {
			if got.ID == idTarget {
				spentGoneFromList = false
				t.Errorf("the revoked row %d is still listed as live after the panel spent it (payload=%+v)", idTarget, after)
			}
		}
		if afterRows != 1 {
			t.Errorf("the live list after the revoke holds %d rows, want exactly 1 (the untouched path row)", afterRows)
		}
		if after.Zero != "" {
			t.Errorf("a one-row list carries a zero sentence %q", after.Zero)
		}

		// The rows as the store really holds them at this moment (BEFORE the
		// session-end sweep stamps everything): the revoked row is retained
		// on disk with revoked_at stamped - the audit window's own row - and
		// the untouched row shows no stamp.
		mid2, err := rt.session.ListBySession(ctx)
		if err != nil {
			t.Fatalf("ListBySession after revoke: %v", err)
		}
		for _, g := range mid2 {
			switch g.ID {
			case idTarget:
				targetRevokedOnDisk = g.RevokedAt != nil
			case idPath:
				pathLiveOnDisk = g.RevokedAt == nil
			}
		}
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}

	// EMPTY LEDGER: the zero sentence names the session it read.
	if !strings.Contains(emptyZero, "当前没有生效中的授权") {
		t.Errorf("the empty-ledger zero sentence %q is not the explicit zero-row sentence", emptyZero)
	}
	if !strings.Contains(emptyZero, minted) {
		t.Errorf("the zero sentence %q does not name the session %s it read", emptyZero, minted)
	}
	if !targetRevokedOnDisk {
		t.Errorf("row %d the panel route revoked carries no revoked_at on disk (mid-run read)", idTarget)
	}
	if !pathLiveOnDisk {
		t.Errorf("row %d shows a revoked_at stamp before any revoke touched it", idPath)
	}
	if !spentGoneFromList {
		t.Errorf("the revoked row %d never left the live list", idTarget)
	}
}

// ---------------------------------------------------------------------------
// the AC#2 integration fixture: panel revoke bites on the NEXT covered call
// ---------------------------------------------------------------------------

func Test40SessionRulesPanelRevokePopsTheNextCoveredCall(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	// The console stays attached (a silent reader): an L1 window only opens
	// when somebody could answer, per ticket 375.
	h.reply = strings.NewReader("")
	granted := filepath.Join(h.dir, "rules-revocable.txt")
	ctx := context.Background()

	var (
		rowID    int64
		canonic  string
		receipt  string
		redCards int
	)

	h.rtHook = func(rt *agentRuntime) {
		if rt.session == nil {
			t.Fatalf("assembled run has no session ledger")
		}
		c, err := rt.paths.Canonicalize(filepath.ToSlash(granted))
		if err != nil {
			t.Fatalf("Canonicalize: %v", err)
		}
		canonic = c
		// 「先本会话允许」: the grant is minted through the ledger's Record
		// door - the SAME write the 本会话内允许 answer route spends
		// (internal/session/grants.go:190, "called from the native answer
		// route ONLY"). One answered card, one stored rule.
		rowID, err = rt.session.Record(ctx, "fs.write", c)
		if err != nil {
			t.Fatalf("session.Record: %v", err)
		}

		// The production chain, bound the way the boot binds it.
		disp, _, err := newComposerDispatchChain(h.dir, auditQuiet194, "test-panel", nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("the production chain did not assemble on %s: %v", h.dir, err)
		}
		door, ok := disp.SessionRules.(*sessionRulesDoor40)
		if !ok {
			t.Fatalf("session-rules leg is %T, want the real door", disp.SessionRules)
		}
		door.bindSeam(rt.sessionRulesSeam40())

		revoke := rt.gate.AdmitTextTask("t40-panel-revoke")
		defer revoke()

		// READING 1 (negative control): the grant is live and nothing is
		// revoked, so the covered call must run silently - zero cards on the
		// production UI.
		// The file is removed first so this call keeps the CREATE shape:
		// probed 2026-11-18 (mutation lab, three-covered-calls reading), a
		// second fs.write OVER an existing file elevates the verdict beyond
		// L1 - and L2 is never grant-covered (SPEC-06 §8.3/D45-3) - so an
		// overwrite-shaped "red" call would ask even with the grant live and
		// the reading would be confounded. Same shape in, same shape out;
		// the ONLY variable between the two readings is the revoke.
		if err := os.Remove(granted); err != nil && !os.IsNotExist(err) {
			t.Fatalf("pre-green remove: %v", err)
		}
		before := rt.windowCount()
		args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"covered, nothing revoked"}`,
			filepath.ToSlash(granted)))
		if _, err := rt.bridge.Execute(context.Background(), agent.ToolRequest{
			TaskID: "t40-panel-revoke", CorrelationID: "t40-panel-green", CallID: "g1",
			Name: "fs.write", Args: args,
		}); err != nil {
			t.Fatalf("negative-control call: %v", err)
		}
		if cards := rt.windowCount() - before; cards != 0 {
			t.Errorf("with the grant live the covered call showed %d cards, want 0", cards)
		}

		// THE PANEL REVOKE: session.rules.revoke through the production
		// dispatch - not a direct session.Revoke call. The selector is the
		// row id the read face carried.
		receipt, err = disp.Handle(ctx, sessionRulesEnvelope40("rid-sr-pop", panel.MethodSessionRulesRevoke, fmt.Sprint(rowID)))
		if err != nil {
			t.Fatalf("session.rules.revoke %d through the panel route: %v", rowID, err)
		}

		// READING 2 (red): the SAME shaped call must ask again - a real card
		// on the real gate, counted on the same production counter the
		// ticket-49 cells read. The remove keeps the create shape (see the
		// negative control's note: without it the overwrite elevation, not
		// the revoke, would be what asks - the mutation lab proved that
		// reading vacuous, 2026-11-18).
		if err := os.Remove(granted); err != nil && !os.IsNotExist(err) {
			t.Fatalf("pre-red remove: %v", err)
		}
		before = rt.windowCount()
		args2 := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"after the panel revoked"}`,
			filepath.ToSlash(granted)))
		if _, err := rt.bridge.Execute(context.Background(), agent.ToolRequest{
			TaskID: "t40-panel-revoke", CorrelationID: "t40-panel-red", CallID: "r1",
			Name: "fs.write", Args: args2,
		}); err != nil {
			t.Fatalf("red call: %v", err)
		}
		redCards = rt.windowCount() - before
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}

	// The receipt names the row it spent.
	if !strings.Contains(receipt, fmt.Sprintf("已撤销会话授权 #%d", rowID)) {
		t.Errorf("the revoke receipt %q does not name row #%d", receipt, rowID)
	}
	if !strings.Contains(grantLinesOf(h), "GRANT-REVOKE") {
		t.Errorf("the audit is missing the ledger's GRANT-REVOKE line:\n%s", grantLinesOf(h))
	}
	if redCards < 1 {
		t.Errorf("after the panel revoke the same call showed %d cards, want at least 1: "+
			"the route spent row %d but the next covered call did not re-ask (pattern %q)",
			redCards, rowID, canonic)
	}

	// SQL reading off the real store: the negative-control row carried
	// grant_id; the row the panel revoked answered nothing afterwards.
	st := h.openStore()
	rows, err := st.ListToolCallsByTask(ctx, "t40-panel-revoke")
	if err != nil {
		t.Fatalf("ListToolCallsByTask: %v", err)
	}
	var green, red *memoryToolCallView
	for i := range rows {
		switch rows[i].CorrelationID {
		case "t40-panel-green":
			green = &memoryToolCallView{rows[i].Decision, rows[i].GrantID}
		case "t40-panel-red":
			red = &memoryToolCallView{rows[i].Decision, rows[i].GrantID}
		}
	}
	if green == nil || red == nil {
		t.Fatalf("missing rows: green=%v red=%v (rows=%+v)", green != nil, red != nil, rows)
	}
	if green.decision != agent.DecisionAllowGrant || green.grantID == nil || *green.grantID != rowID {
		t.Errorf("negative-control row: decision=%q grant_id=%v, want %q/%d",
			green.decision, green.grantID, agent.DecisionAllowGrant, rowID)
	}
	if red.decision == agent.DecisionAllowGrant {
		t.Errorf("red row booked decision=allow_session_grant: a revoked grant answered a call")
	}
	if red.grantID != nil {
		t.Errorf("red row booked grant_id=%d, want NULL", *red.grantID)
	}
}

// memoryToolCallView keeps the two SQL facts this fixture reads off the store,
// so the switch above needs no memory import beyond what the rows carry.
type memoryToolCallView struct {
	decision string
	grantID  *int64
}
