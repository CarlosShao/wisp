package main

// Ticket 49's VERB teeth (2026-11-17 audit docs/evidence/s1/
// 4849-adversarial-20261117-1052.md, 判票 49「动词层空牙」): the operator
// verbs `grants` and `revoke-grant` in cmd/wisp/approval_reply.go had ZERO
// test calls - the assembly file ticket49_assembly_test.go exercises the same
// ledger by calling session.Record / session.Revoke DIRECTLY, so a broken
// dispatch, a broken number parse or a broken list rendering in the verbs
// themselves could ship green. These two teeth drive the OPERATOR VERBS the
// console really runs, through replySurface.handle, against the same
// assembled-run shape:
//
//	revoke-grant	TestTicket49RevokeGrantVerbPopsTheNextCoveredCallBack
//	grants		TestTicket49GrantsVerbReadsTheRealRowsFieldByField
//
// The judgment ceiling ticket49_assembly_test.go states applies here too:
// every cell below is ONE ASSEMBLED RUN in the test process - the real
// wisp.db, the real ledger, the real bridge, the real gate - and the card
// counts are read off rt.windowCount() (the production UI's own shown()
// counter), never off a mock. No source-string assertion anywhere: every
// claim is a behavior the run really performed.

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/memory"
)

// ---------------------------------------------------------------------------
// tooth 1: the revoke-grant verb spends the row, and the NEXT covered call
// really asks again
// ---------------------------------------------------------------------------

// TestTicket49RevokeGrantVerbPopsTheNextCoveredCallBack is the verb-shaped
// redo of the assembly tooth's red/green cell: with a live session grant, a
// covered fs.write call runs silently (NEGATIVE CONTROL: zero cards); then
// `revoke-grant <编号>` is sent through replySurface.handle - the exact door
// the console's answer stream walks - and the SAME shaped call must show a
// real card on the real gate, counted on the production UI counter
// rt.windowCount() (the shape ticket49_assembly_test.go's revoke cell reads).
func TestTicket49RevokeGrantVerbPopsTheNextCoveredCallBack(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	// The console stays attached (a silent reader): an L1 window only opens
	// when somebody could answer, per ticket 375.
	h.reply = strings.NewReader("")
	granted := filepath.Join(h.dir, "verb-revocable.txt")
	ctx := context.Background()

	var (
		rowID    int64
		canonic  string
		verbSaid string
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
		rowID, err = rt.session.Record(ctx, "fs.write", c)
		if err != nil {
			t.Fatalf("session.Record: %v", err)
		}

		revoke := rt.gate.AdmitTextTask("t49-verb-revoke")
		defer revoke()

		// NEGATIVE CONTROL: the grant is live and nothing is revoked, so the
		// covered call must run silently - zero cards on the production UI.
		before := rt.windowCount()
		args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"covered, nothing revoked"}`,
			filepath.ToSlash(granted)))
		if _, err := rt.bridge.Execute(context.Background(), agent.ToolRequest{
			TaskID: "t49-verb-revoke", CorrelationID: "t49-verb-green", CallID: "g1",
			Name: "fs.write", Args: args,
		}); err != nil {
			t.Fatalf("negative-control call: %v", err)
		}
		if cards := rt.windowCount() - before; cards != 0 {
			t.Errorf("with the grant live the covered call showed %d cards, want 0", cards)
		}

		// THE VERB: `revoke-grant <编号>` through replySurface.handle - not a
		// direct session.Revoke call. The number is the row the ledger minted.
		s := rt.reply
		if s == nil {
			t.Fatalf("assembled run has no reply surface")
		}
		said, err := s.handle("revoke-grant", fmt.Sprint(rowID), "")
		if err != nil {
			t.Fatalf("revoke-grant %d through the verb: %v", rowID, err)
		}
		verbSaid = said

		// RED: the SAME shaped call must ask again - a real card on the real
		// gate, counted on the same production counter the assembly cell reads.
		before = rt.windowCount()
		args2 := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"after the verb revoked"}`,
			filepath.ToSlash(granted)))
		if _, err := rt.bridge.Execute(context.Background(), agent.ToolRequest{
			TaskID: "t49-verb-revoke", CorrelationID: "t49-verb-red", CallID: "r1",
			Name: "fs.write", Args: args2,
		}); err != nil {
			t.Fatalf("red call: %v", err)
		}
		redCards = rt.windowCount() - before
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}

	// The verb's own answer names the row it spent (the number an operator
	// reads back against `grants`), and the audit names the revocation.
	if !strings.Contains(verbSaid, fmt.Sprintf("已撤销授权 #%d", rowID)) {
		t.Errorf("the revoke-grant verb answered %q, want it to name row #%d", verbSaid, rowID)
	}
	if !strings.Contains(grantLinesOf(h), "GRANT-REVOKED") {
		t.Errorf("the audit is missing the verb's GRANT-REVOKED line:\n%s", grantLinesOf(h))
	}
	if redCards < 1 {
		t.Errorf("after the revoke-grant verb the same call showed %d cards, want at least 1: "+
			"the verb spent row %d but the next covered call did not re-ask (pattern %q)",
			redCards, rowID, canonic)
	}

	// SQL reading off the real store: the negative-control row carried
	// grant_id; the row the verb revoked answered nothing afterwards.
	st := h.openStore()
	rows, err := st.ListToolCallsByTask(ctx, "t49-verb-revoke")
	if err != nil {
		t.Fatalf("ListToolCallsByTask: %v", err)
	}
	var green, red *memory.ToolCall
	for i := range rows {
		switch rows[i].CorrelationID {
		case "t49-verb-green":
			green = &rows[i]
		case "t49-verb-red":
			red = &rows[i]
		}
	}
	if green == nil || red == nil {
		t.Fatalf("missing rows: green=%v red=%v (rows=%+v)", green != nil, red != nil, rows)
	}
	if green.Decision != agent.DecisionAllowGrant || green.GrantID == nil || *green.GrantID != rowID {
		t.Errorf("negative-control row: decision=%q grant_id=%v, want %q/%d",
			green.Decision, green.GrantID, agent.DecisionAllowGrant, rowID)
	}
	if red.Decision == agent.DecisionAllowGrant {
		t.Errorf("red row booked decision=allow_session_grant: a revoked grant answered a call")
	}
	if red.GrantID != nil {
		t.Errorf("red row booked grant_id=%d, want NULL", *red.GrantID)
	}
}

// ---------------------------------------------------------------------------
// tooth 2: the grants verb reads the REAL rows off the store, field by field
// ---------------------------------------------------------------------------

// TestTicket49GrantsVerbReadsTheRealRowsFieldByField drives `grants` through
// replySurface.handle three times in one assembled run:
//
//	EMPTY LEDGER  - the negative control: the verb says zero rows OUT LOUD
//	                (the explicit 「当前没有生效中的授权」 sentence), it does
//	                not swallow the read into silence or an error;
//	TWO REAL ROWS - one fs.write path row through Record and one input.type
//	                target row through RecordTarget, then the verb's listing
//	                must carry each row field by field (#grant_id tool
//	                「pattern」 created=CreatedAt state=生效中) against the
//	                rows read back from the store;
//	AFTER REVOKE  - the target row is revoked THROUGH THE VERB, and the next
//	                listing must show that row 已撤销 (its revoked_at really
//	                stamped on disk, read back mid-run before the session-end
//	                sweep) while the untouched row stays 生效中.
func TestTicket49GrantsVerbReadsTheRealRowsFieldByField(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	h.reply = strings.NewReader("")
	note := filepath.Join(h.dir, "grants-note.txt")
	ctx := context.Background()

	var (
		minted string
		// the two rows as the store really holds them (read back mid-run)
		idPath, idTarget     int64
		pathTool, pathPat    string
		targetTool, targetPt string
		createdPath          int64
		createdTarget        int64
		targetRevokedOnDisk  bool
		pathLiveOnDisk       bool
		// the verb's three answers
		emptySaid string
		listSaid  string
		afterSaid string
	)

	h.rtHook = func(rt *agentRuntime) {
		if rt.session == nil {
			t.Fatalf("assembled run has no session ledger")
		}
		s := rt.reply
		if s == nil {
			t.Fatalf("assembled run has no reply surface")
		}
		minted = rt.session.SessionID().String()

		// NEGATIVE CONTROL: nothing recorded yet. grants must answer with the
		// explicit zero-row sentence, not an error and not silence.
		said, err := s.handle("grants", "", "")
		if err != nil {
			t.Fatalf("grants on an empty ledger: %v", err)
		}
		emptySaid = said

		// Two REAL rows through the ledger's two write doors.
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
		listSaid, err = s.handle("grants", "", "")
		if err != nil {
			t.Fatalf("grants with two rows: %v", err)
		}

		// Revoke the target row THROUGH THE VERB (the same door tooth 1
		// drives), then read the listing again.
		if _, err := s.handle("revoke-grant", fmt.Sprint(idTarget), ""); err != nil {
			t.Fatalf("revoke-grant %d through the verb: %v", idTarget, err)
		}
		afterSaid, err = s.handle("grants", "", "")
		if err != nil {
			t.Fatalf("grants after the revoke: %v", err)
		}

		// The rows as the store really holds them at this moment (BEFORE the
		// session-end sweep stamps everything): the field-by-field expectation
		// below is built from these read-back rows, and the revoked_at state
		// is proven on disk here, where the sweep cannot have muddied it.
		mid, err := rt.session.ListBySession(ctx)
		if err != nil {
			t.Fatalf("ListBySession: %v", err)
		}
		if len(mid) != 2 {
			t.Fatalf("the ledger holds %d rows mid-run, want exactly 2 (rows=%+v)", len(mid), mid)
		}
		for _, g := range mid {
			switch g.ID {
			case idPath:
				createdPath, pathTool, pathPat = g.CreatedAt, g.Tool, g.Pattern
				pathLiveOnDisk = g.RevokedAt == nil
			case idTarget:
				createdTarget, targetTool, targetPt = g.CreatedAt, g.Tool, g.Pattern
				targetRevokedOnDisk = g.RevokedAt != nil
			}
		}
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}

	// EMPTY LEDGER: zero rows said out loud, naming this session.
	if !strings.Contains(emptySaid, "当前没有生效中的授权") {
		t.Errorf("grants on an empty ledger said %q, want the explicit zero-row sentence", emptySaid)
	}
	if !strings.Contains(emptySaid, minted) {
		t.Errorf("the zero-row answer %q does not name the session %s it read", emptySaid, minted)
	}

	// TWO REAL ROWS: every field of every read-back row must appear in the
	// listing, spelled the way an operator reads it: #grant_id tool「pattern」
	// created=CreatedAt state=生效中.
	if !strings.Contains(listSaid, fmt.Sprintf("本会话（%s）的授权", minted)) {
		t.Errorf("the listing %q does not name the session %s", listSaid, minted)
	}
	for _, want := range []struct {
		row, tool, pat string
		created        int64
	}{
		{fmt.Sprint(idPath), pathTool, pathPat, createdPath},
		{fmt.Sprint(idTarget), targetTool, targetPt, createdTarget},
	} {
		line := fmt.Sprintf("#%s %s「%s」 created=%d state=生效中", want.row, want.tool, want.pat, want.created)
		if !strings.Contains(listSaid, line) {
			t.Errorf("the grants listing is missing the real row verbatim:\n  want: %s\n  got: %s", line, listSaid)
		}
	}

	// AFTER REVOKE: the revoked row's state comes from its real revoked_at
	// stamp (proven on disk mid-run), and the untouched row stays live.
	if !targetRevokedOnDisk {
		t.Errorf("row %d the verb revoked carries no revoked_at on disk (mid-run read)", idTarget)
	}
	if !pathLiveOnDisk {
		t.Errorf("row %d shows a revoked_at stamp before any revoke touched it", idPath)
	}
	revokedLine := fmt.Sprintf("#%d %s「%s」 created=%d state=已撤销", idTarget, targetTool, targetPt, createdTarget)
	if !strings.Contains(afterSaid, revokedLine) {
		t.Errorf("after the revoke the grants listing is missing the revoked row verbatim:\n  want: %s\n  got: %s",
			revokedLine, afterSaid)
	}
	liveLine := fmt.Sprintf("#%d %s「%s」 created=%d state=生效中", idPath, pathTool, pathPat, createdPath)
	if !strings.Contains(afterSaid, liveLine) {
		t.Errorf("after the revoke the grants listing lost the untouched live row:\n  want: %s\n  got: %s",
			liveLine, afterSaid)
	}
	if strings.Contains(afterSaid, fmt.Sprintf("#%d %s「%s」 created=%d state=生效中",
		idTarget, targetTool, targetPt, createdTarget)) {
		t.Errorf("the revoked row %d still lists as 生效中 after the verb spent it:\n%s", idTarget, afterSaid)
	}
}
