package main

// Ticket 49 tests at the assembled-run layer - the cells the dispatch asked
// for as SQL/behavior READINGS from the run cmd/wisp actually builds, in the
// same host shape ticket224_assembly_test.go established (the real wisp.db,
// the real ledger, the real bridge, the operator's reply stream).
//
// The cells and where their halves live:
//
//	lifecycle (audit rows per call, session-end sweep)
//	                      TestTicket49SessionEndSweepRevokesEveryLiveGrant
//	revoke (immediate)    TestTicket49RevokeVerbStopsTheNextCoveredCall
//	audit (SQL)           TestTicket49EveryGrantBackedExecutionCarriesGrantID
//	                      (reads the real tool_call rows off disk)
//	batch one-card        covered at the approval layer
//	                      (ticket49_answer_doors_test.go); the assembled-run
//	                      half of the danger tier is
//	                      TestTicket49SessionAllVerbRunsTheRitualEndToEnd
//
// The judgment ceiling ticket224_assembly_test.go states applies here too:
// every "restart"/"session end" below is a SECOND ASSEMBLY in the same test
// process, not a second OS process; what the cells prove is the STAMP (the
// rows' revoked_at is set and the same ledger's next read answers nothing),
// with the minted-identity invalidation already covered by ticket 224's own
// cross-boot case.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// ---------------------------------------------------------------------------
// session end: the close-time sweep stamps every live row, and the audit
// per-call rows the dispatch asked for
// ---------------------------------------------------------------------------

// TestTicket49SessionEndSweepRevokesEveryLiveGrant drives one assembled run
// whose hook records a session grant for one path (the production shape of
// 「本会话内允许」: the row exists because a card was answered), asserts the
// per-call audit READINGS the dispatch named, and then closes the runtime -
// the sweep's stamp must leave every row revoked while the rows stay on disk
// for the audit window.
func TestTicket49SessionEndSweepRevokesEveryLiveGrant(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	// The console stays attached (a silent reader): the L1-window control
	// needs somebody who could answer, per ticket 375.
	h.reply = strings.NewReader("")
	granted := filepath.Join(h.dir, "sweep-note.txt")
	ctx := context.Background()

	var (
		ledgerID string
		rowID    int64
	)

	h.rtHook = func(rt *agentRuntime) {
		if rt.session == nil {
			t.Fatalf("assembled run has no session ledger")
		}
		ledgerID = rt.session.SessionID().String()
		c, err := rt.paths.Canonicalize(filepath.ToSlash(granted))
		if err != nil {
			t.Fatalf("Canonicalize: %v", err)
		}
		id, err := rt.session.Record(ctx, "fs.write", c)
		if err != nil {
			t.Fatalf("session.Record: %v", err)
		}
		rowID = id

		// Per-call audit reading 1 (the dispatch's 逐调用审计行): one covered
		// call books its own tool_call row with decision=allow_session_grant
		// and grant_id pointing at the covering row. Read back OFF DISK below.
		revoke := rt.gate.AdmitTextTask("t49-sweep")
		defer revoke()
		args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"covered by the session grant"}`,
			filepath.ToSlash(granted)))
		if _, err := rt.bridge.Execute(context.Background(), agent.ToolRequest{
			TaskID: "t49-sweep", CorrelationID: "t49-sweep-c1", CallID: "c1",
			Name: "fs.write", Args: args,
		}); err != nil {
			t.Fatalf("covered call: %v", err)
		}
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}

	// The tool_call row the covered call booked: grant_id on every
	// grant-backed execution, decision=allow_session_grant.
	st := h.openStore()
	rows, err := st.ListToolCallsByTask(ctx, "t49-sweep")
	if err != nil {
		t.Fatalf("ListToolCallsByTask: %v", err)
	}
	var covered *memory.ToolCall
	for i := range rows {
		if rows[i].CorrelationID == "t49-sweep-c1" {
			covered = &rows[i]
		}
	}
	if covered == nil {
		t.Fatalf("no tool_call row for corr=t49-sweep-c1: rows=%+v", rows)
	}
	if covered.Decision != agent.DecisionAllowGrant {
		t.Errorf("covered call decision=%q, want %q", covered.Decision, agent.DecisionAllowGrant)
	}
	if covered.GrantID == nil || *covered.GrantID != rowID {
		t.Errorf("covered call grant_id=%v, want the covering row %d", covered.GrantID, rowID)
	}

	// The grant row itself, keyed to the minted identity.
	grants, err := st.ListGrantsBySession(ctx, ledgerID)
	if err != nil {
		t.Fatalf("ListGrantsBySession: %v", err)
	}
	if len(grants) != 1 || grants[0].ID != rowID {
		t.Fatalf("grant rows=%+v, want exactly %d", grants, rowID)
	}

	// SESSION END: the run's own shutdown sequence ALREADY swept - rt.close()
	// is called by runTextTask before it returns, and its SESSION-GRANT-SWEEP
	// stamped every live row revoked (the production lifecycle step). Read the
	// two facts back: every row stamped ON DISK, the stamp visible to the very
	// next read of the SAME ledger identity, and the rows retained.
	after, err := st.ListGrantsBySession(ctx, ledgerID)
	if err != nil {
		t.Fatalf("ListGrantsBySession(after): %v", err)
	}
	if len(after) != 1 {
		t.Fatalf("grant rows after session end = %d, want 1 (retained for audit, not deleted)", len(after))
	}
	if after[0].RevokedAt == nil {
		t.Errorf("after session end the row must be revoked ON DISK: rows=%+v", after)
	}

	// And the audit line named the sweep with the row count.
	if !strings.Contains(h.err.String(), "SESSION-GRANT-SWEEP") ||
		!strings.Contains(h.err.String(), "revoked=1") {
		t.Errorf("audit is missing the session-end sweep line naming the count:\n%s", grantLinesOf(h))
	}
}

// argIntAny extracts one integer argument from the call's parameter map (the
// fixture resolver's hwnd selector). tools.argInt is package-internal to
// internal/tools; this 10-line local copy keeps the resolver testable without
// exporting anything.
func argIntAny(params map[string]any, key string) (int64, bool) {
	v, ok := params[key]
	if !ok {
		return 0, false
	}
	switch t := v.(type) {
	case int:
		return int64(t), true
	case int64:
		return t, true
	case float64:
		return int64(t), true
	default:
		return 0, false
	}
}

// ---------------------------------------------------------------------------
// revoke: the operator verb stops the NEXT covered call
// ---------------------------------------------------------------------------

// TestTicket49RevokeVerbStopsTheNextCoveredCall is AC#2's red/green cell at
// assembled-run scale: with the grant live, a covered call runs silently with
// grant_id booked; after the ledger's Revoke (the door `revoke-grant` spends),
// the SAME call asks its L1 window again and books no grant linkage.
func TestTicket49RevokeVerbStopsTheNextCoveredCall(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	h.reply = strings.NewReader("")
	granted := filepath.Join(h.dir, "revocable.txt")
	ctx := context.Background()

	var (
		rowID    int64
		canonic  string
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

		revoke := rt.gate.AdmitTextTask("t49-revoke")
		defer revoke()

		// GREEN: covered, silent, booked.
		before := rt.windowCount()
		args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"covered"}`, filepath.ToSlash(granted)))
		if _, err := rt.bridge.Execute(context.Background(), agent.ToolRequest{
			TaskID: "t49-revoke", CorrelationID: "t49-revoke-green", CallID: "g1",
			Name: "fs.write", Args: args,
		}); err != nil {
			t.Fatalf("green call: %v", err)
		}
		if cards := rt.windowCount() - before; cards != 0 {
			t.Errorf("the covered call showed %d cards, want 0", cards)
		}

		// RED: revoke, then the SAME call asks.
		if err := rt.session.Revoke(ctx, rowID); err != nil {
			t.Fatalf("Revoke: %v", err)
		}
		before = rt.windowCount()
		args2 := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"after revoke"}`, filepath.ToSlash(granted)))
		if _, err := rt.bridge.Execute(context.Background(), agent.ToolRequest{
			TaskID: "t49-revoke", CorrelationID: "t49-revoke-red", CallID: "r1",
			Name: "fs.write", Args: args2,
		}); err != nil {
			t.Fatalf("red call: %v", err)
		}
		redCards = rt.windowCount() - before
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}

	if redCards < 1 {
		t.Errorf("after revoke the same call showed %d cards, want at least 1: "+
			"「撤销后下一次受门控调用重新弹卡」 is not happening in the assembled run (pattern %q, row %d)",
			redCards, canonic, rowID)
	}

	// SQL reading: the green row carries grant_id, the red row does not.
	st := h.openStore()
	rows, err := st.ListToolCallsByTask(ctx, "t49-revoke")
	if err != nil {
		t.Fatalf("ListToolCallsByTask: %v", err)
	}
	var green, red *memory.ToolCall
	for i := range rows {
		switch rows[i].CorrelationID {
		case "t49-revoke-green":
			green = &rows[i]
		case "t49-revoke-red":
			red = &rows[i]
		}
	}
	if green == nil || red == nil {
		t.Fatalf("missing rows: green=%v red=%v (rows=%+v)", green != nil, red != nil, rows)
	}
	if green.Decision != agent.DecisionAllowGrant || green.GrantID == nil || *green.GrantID != rowID {
		t.Errorf("green row: decision=%q grant_id=%v, want %q/%d", green.Decision, green.GrantID, agent.DecisionAllowGrant, rowID)
	}
	if red.Decision == agent.DecisionAllowGrant {
		t.Errorf("red row booked decision=allow_session_grant: a revoked grant answered a call")
	}
	if red.GrantID != nil {
		t.Errorf("red row booked grant_id=%d, want NULL", *red.GrantID)
	}
}

// ---------------------------------------------------------------------------
// the danger tier through the operator's own verb: expand → all-apps
// ---------------------------------------------------------------------------

// TestTicket49SessionAllVerbRunsTheRitualEndToEnd is the input.type danger
// cell where the user actually is: an assembled run, an input.type card (the
// decision carries the resolved target process), the operator types
// session-all WITHOUT the expand step (refused - the ritual is a sequence),
// then expands and confirms within the window - and the wildcard row lands
// keyed to the run's own session.
func TestTicket49SessionAllVerbRunsTheRitualEndToEnd(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	pr, pw := io.Pipe()
	h.reply = pr
	t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })

	const task, corr = "t49-all", "t49-all-corr"
	var (
		minted string
		done   = make(chan struct{})
		out    agent.ToolOutcome
	)
	h.rtHook = func(rt *agentRuntime) {
		if rt.session == nil {
			t.Fatalf("assembled run has no session ledger")
		}
		minted = rt.session.SessionID().String()
		revoke := rt.gate.AdmitTextTask(task)
		defer revoke()
		args := json.RawMessage(`{"text":"hello","hwnd":7}`)
		go func() {
			defer close(done)
			out, _ = rt.bridge.Execute(context.Background(), agent.ToolRequest{
				TaskID: task, CorrelationID: corr, CallID: "v1",
				Name: "input.type", Args: args,
			})
		}()
		if !waitForCard187(rt.liveCards, corr, 10*time.Second) {
			t.Errorf("the input.type card never reached the native ledger")
			return
		}

		// STEP 1: session-all without the expand step must be refused (the
		// ritual is a sequence), and the card stays pending. The refusal names
		// the expand step.
		s := rt.reply
		_, err := s.handle("session-all", corr, "")
		if err == nil || !strings.Contains(err.Error(), "session-expand") {
			t.Errorf("the un-expanded session-all must be refused naming the expand step: err=%v", err)
		}

		// STEP 2: expand arms the window, session-all spends it. The card is
		// answered, the call released, and the wildcard row lands.
		if _, err := s.handle("session-expand", corr, ""); err != nil {
			t.Fatalf("session-expand: %v", err)
		}
		if _, err := s.handle("session-all", corr, ""); err != nil {
			t.Fatalf("session-all after expand: %v", err)
		}
		<-done
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	<-done
	if out.IsError {
		t.Fatalf("the danger-answered call must execute, got: %s", out.Text)
	}

	st := h.openStore()
	rows, err := st.ListGrantsBySession(context.Background(), minted)
	if err != nil {
		t.Fatalf("ListGrantsBySession: %v", err)
	}
	if len(rows) != 1 || rows[0].Pattern != "*" || rows[0].Tool != "input.type" {
		t.Fatalf("rows=%+v, want exactly one (input.type, *) row keyed to %s", rows, minted)
	}
	// The audit named the danger tier.
	if !strings.Contains(h.err.String(), "GRANT-RECORDED-ALL-APPS") {
		t.Errorf("audit is missing the danger-tier recording line:\n%s", grantLinesOf(h))
	}
}

// ---------------------------------------------------------------------------
// input.type binding through the assembled bridge: per-process + downgrade
// ---------------------------------------------------------------------------

// TestTicket49AssembledTargetBindingDowngradesAndIsolates joins the assembled
// bridge, the real ledger and the real store for the §15 item-7 cells:
//
//	RED   - the first injection into a process (no covering row) asks its L2
//	        card and books no grant linkage;
//	GREEN - after the per-process row is recorded through the ledger's TARGET
//	        door (the production answer route), the SAME call runs silently
//	        and books decision=allow_session_grant with the row's id;
//	ISOLATION - a different target process still asks (Word ≠ 终端).
//
// The target resolver rides the run's OWN SystemBridge (the same object the
// tools execute through). A runner cannot guarantee a desktop window, so the
// hook pins the resolution the Windows backend would have produced for the
// hwnd the calls name - by overriding the bridge's resolver field DIRECTLY
// (same package), with each call's hwnd mapped to a stable process name the
// way realProcessName would. That keeps the test's resolution one field wide
// and the production resolution path untouched for real hwnds.
func TestTicket49AssembledTargetBindingDowngradesAndIsolates(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	h.reply = strings.NewReader("")
	// The fixture resolver: hwnd 7 → notepad.exe, hwnd 9 → cmd.exe. A runner
	// cannot guarantee a desktop window to resolve a real hwnd against, and
	// the pid-fallback refusal is a shape rule (not overridden).
	h.resolveTarget = func(params map[string]any) (string, bool) {
		if n, ok := argIntAny(params, "hwnd"); ok && n == 9 {
			return "cmd.exe", true
		}
		return "notepad.exe", true
	}
	ctx := context.Background()

	var (
		ledgerID string
		rowID    int64
		redCards int
	)

	h.rtHook = func(rt *agentRuntime) {
		if rt.session == nil {
			t.Fatalf("assembled run has no session ledger")
		}
		ledgerID = rt.session.SessionID().String()
		revoke := rt.gate.AdmitTextTask("t49-target")
		defer revoke()

		targetArgs := json.RawMessage(`{"text":"hello","hwnd":7}`)

		// RED: the first injection asks its L2 card.
		before := rt.windowCount()
		if _, err := rt.bridge.Execute(context.Background(), agent.ToolRequest{
			TaskID: "t49-target", CorrelationID: "t49-target-red", CallID: "r1",
			Name: "input.type", Args: targetArgs,
		}); err != nil {
			t.Fatalf("red call: %v", err)
		}
		redCards = rt.windowCount() - before

		// GREEN: record the per-process row through the production TARGET
		// door, then the same call runs silently and books the linkage.
		id, err := rt.session.RecordTarget(ctx, "input.type", "notepad.exe")
		if err != nil {
			t.Fatalf("RecordTarget: %v", err)
		}
		rowID = id
		before = rt.windowCount()
		if _, err := rt.bridge.Execute(context.Background(), agent.ToolRequest{
			TaskID: "t49-target", CorrelationID: "t49-target-green", CallID: "g1",
			Name: "input.type", Args: targetArgs,
		}); err != nil {
			t.Fatalf("green call: %v", err)
		}
		if cards := rt.windowCount() - before; cards != 0 {
			t.Errorf("the covered input.type call showed %d cards, want 0: the downgrade did not happen", cards)
		}

		// ISOLATION: a different process asks again.
		before = rt.windowCount()
		if _, err := rt.bridge.Execute(context.Background(), agent.ToolRequest{
			TaskID: "t49-target", CorrelationID: "t49-target-isolated", CallID: "i1",
			Name: "input.type",
			Args: json.RawMessage(`{"text":"hello","hwnd":9}`),
		}); err != nil {
			t.Fatalf("isolated call: %v", err)
		}
		if cards := rt.windowCount() - before; cards < 1 {
			t.Errorf("a different target showed %d cards, want at least 1: per-process isolation (Word ≠ 终端)", cards)
		}
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}

	// RED must have asked.
	if redCards < 1 {
		t.Errorf("the first injection showed %d cards, want at least 1 (L2 without a grant)", redCards)
	}

	// SQL reading off the real store.
	st := h.openStore()
	rows, err := st.ListToolCallsByTask(ctx, "t49-target")
	if err != nil {
		t.Fatalf("ListToolCallsByTask: %v", err)
	}
	var red, green *memory.ToolCall
	for i := range rows {
		switch rows[i].CorrelationID {
		case "t49-target-red":
			red = &rows[i]
		case "t49-target-green":
			green = &rows[i]
		}
	}
	if green == nil {
		t.Fatalf("no row for the green call: rows=%+v", rows)
	}
	if green.Decision != agent.DecisionAllowGrant {
		t.Errorf("green input.type decision=%q, want %q", green.Decision, agent.DecisionAllowGrant)
	}
	if green.GrantID == nil || *green.GrantID != rowID {
		t.Errorf("green input.type grant_id=%v, want %d", green.GrantID, rowID)
	}
	if red != nil && red.GrantID != nil {
		t.Errorf("red input.type row booked grant_id=%d, want NULL", *red.GrantID)
	}
	// The binding row exists, keyed to the process the resolver named.
	grants, err := st.ListGrantsBySession(ctx, ledgerID)
	if err != nil {
		t.Fatalf("ListGrantsBySession: %v", err)
	}
	found := false
	for _, g := range grants {
		if g.Tool == "input.type" && g.Pattern == "notepad.exe" && g.ID == rowID {
			found = true
		}
	}
	if !found {
		t.Errorf("the (input.type, notepad.exe) binding row is not on disk: rows=%+v", grants)
	}
}

// ---------------------------------------------------------------------------
// the batch cell at assembled-run scale: one card, one pattern row
// ---------------------------------------------------------------------------

// TestTicket49AssembledBatchOneCardOnePatternRow drives the 10-path batch
// through the ASSEMBLED gate, ledger and store: the decision is the D45-1
// batch shape (10 homogeneous paths on one decision), routed through the
// queue-backed confirmation the opt-in click needs (the 3s window
// auto-executes and has no answer verb), answered by the OPERATOR VERB the
// console really runs (`session-pattern`), and read back as exactly one
// whole-pattern row keyed to the run's own session. The card-count witness is
// the one Prompt the answer-side ledger recorded for the correlation.
func TestTicket49AssembledBatchOneCardOnePatternRow(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	h.reply = strings.NewReader("")
	ctx := context.Background()

	var (
		minted  string
		oneCard bool
	)

	h.rtHook = func(rt *agentRuntime) {
		if rt.session == nil {
			t.Fatalf("assembled run has no session ledger")
		}
		minted = rt.session.SessionID().String()
		revoke := rt.gate.AdmitTextTask("t49-batch")
		t.Cleanup(revoke)

		var paths []string
		for i := 0; i < 10; i++ {
			paths = append(paths, filepath.ToSlash(filepath.Join(h.dir, fmt.Sprintf("batch%d.txt", i))))
		}
		d := tools.Decision{
			Tool: "fs.write", Provider: tools.KindBuiltin,
			Level: risk.L1,
			Paths: paths, CorrelationID: "t49-batch-corr",
			TaskID: "t49-batch", CallID: "v1",
		}
		// The D45-1 aggregation witness: the gate's own view of the batch.
		bv := approval.Aggregate(d)
		if bv == nil || bv.Total != 10 {
			t.Fatalf("Aggregate = %+v, want a 10-item view: the batch shape is broken", bv)
		}

		// ONE card (queue-backed) + the operator's own verb. The answer side
		// records exactly one Prompt for this correlation - the one-card
		// witness the AC reads.
		go func() {
			_, _ = rt.gate.PendingApproval(ctx, d)
		}()
		if !waitForCard187(rt.liveCards, "t49-batch-corr", 10*time.Second) {
			t.Fatalf("the batch card never reached the native ledger")
		}
		n := 0
		for _, c := range rt.liveCards.pending() {
			if c.CorrelationID == "t49-batch-corr" {
				n++
			}
		}
		oneCard = n == 1

		s := rt.reply
		if _, err := s.handle("session-pattern", "t49-batch-corr", ""); err != nil {
			t.Fatalf("session-pattern: %v", err)
		}
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}

	// THE CELL, in two readings: one card, one whole-pattern row.
	if !oneCard {
		t.Errorf("the batch confirmation did not show as exactly one tracked card")
	}
	st := h.openStore()
	rows, err := st.ListGrantsBySession(ctx, minted)
	if err != nil {
		t.Fatalf("ListGrantsBySession: %v", err)
	}
	if len(rows) != 1 || rows[0].Pattern != filepath.ToSlash(h.dir)+"/*" {
		t.Fatalf("rows=%+v, want exactly one whole-pattern row %q", rows, filepath.ToSlash(h.dir)+"/*")
	}
	if rows[0].Tool != "fs.write" {
		t.Errorf("the pattern row's tool = %q, want fs.write", rows[0].Tool)
	}
	if !strings.Contains(h.err.String(), "整批模式一次授予") {
		t.Errorf("audit does not name the batch grant:\n%s", grantLinesOf(h))
	}
}

func joinQuoted49(items []string) string {
	quoted := make([]string, len(items))
	for i, s := range items {
		quoted[i] = fmt.Sprintf("%q", s)
	}
	return strings.Join(quoted, ",")
}

// ---------------------------------------------------------------------------
// compile-time keepers for the anchors this file leans on
// ---------------------------------------------------------------------------

var (
	_ = os.Stat
	_ = agent.DecisionAllowGrant
)
