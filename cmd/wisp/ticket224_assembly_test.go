package main

// Ticket 224-r2, cells N#1 (assembly half) and N#2 (the case three comments named
// but nobody wrote).
//
// What this file is for, in one line each:
//
//	N#1  - 224-v1 §2 found that deleting Gate.allowSession's Record call leaves
//	       all four packages green while the log still prints GRANT-RECORDED. Its
//	       §6-c row 1 asked for two instruments: one in internal/agent/approval
//	       (delivered by ticket224_reply_grant_test.go) and one HERE, proving the
//	       answer lands as a real approval_grant row in the run cmd/wisp actually
//	       assembles. The first case below is that one: the operator's verb is
//	       typed on runSpec.reply, the seam AGENTS.md §1.3 names for `wisp run`.
//
//	N#1读 - the read cell has never been run through the assembled bridge with the
//	       REAL ledger on one side and the REAL store on the other: tools' eight
//	       cases stub the source (grant_test.go:51), and the case that books
//	       grant_id books it into a stub journal. The second case joins bridge,
//	       ledger, store and the tool_call.grant_id column in one boot.
//
//	N#2  - TestTicket224ProductionSessionDoesNotSurviveRestart. Three comments in
//	       internal/session/grants_test.go (:20/:94/:134 in the ticket's anchor,
//	       now :20/:94/:134 here as well) and one in internal/tools/grant_test.go
//	       point AT THIS NAME for the cross-boot cell, the expensive half of the
//	       test-literal lock, and ⑩'s controls ① and ②. 224-v1 §6-a ran
//	       `go test ./cmd/wisp/ -run TestTicket224...` and got "no tests to run":
//	       the case did not exist. It exists now.
//
// ---------------------------------------------------------------------------
// JUDGEMENT CEILING - read this before quoting anything below as "restart proof"
// ---------------------------------------------------------------------------
//
// cmd/wisp/run_mode101_test.go:72-76 describes this rig's "Restart" verbatim as
// "a second runTextTask over the same dir: a new process surface", and :143-158
// shows what that is: a second call into the composition root INSIDE THE SAME TEST
// PROCESS. The bytes on disk, the config Manager, the session ledger, the bridge
// and the minted identity are all freshly built; the OS process is not.
//
// So the third case below proves "a second assembly mints a different key and
// reads zero rows". It does NOT prove "a second OS process does", and no
// instrument in this repository currently distinguishes those two meanings of
// "restart" - 224-v1 §7 says so of the pre-existing nails, and it says so of this
// one too. A mint that happened to equal itself across processes while differing
// inside one process (the pid-derived form M1 simulates is caught here, but a
// start-second-derived form is not distinguishable by any case on盘) would still
// slip past. Anything that wants the stronger claim needs a real two-process
// harness, and that harness does not exist yet. Do not cite this file as having
// settled AC#3's OS-process wording.

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
	"github.com/CarlosShao/wisp/internal/memory"
)

// ---------------------------------------------------------------------------
// N#1 - the operator's 「本会话内允许」 verb puts a row in the assembled run's db
// ---------------------------------------------------------------------------

// TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun is the answer-to-row
// hop measured where the user's click actually lands: a whole `wisp run`, its real
// gate, its real ledger, its real wisp.db, and the reply arriving on the CLI's own
// answer stream. The row is then read back with the identity this boot minted, and
// its pattern is required to be the path the CARD printed - not the path this test
// typed.
//
// Mutation direction: with Gate.allowSession's Record call removed (224-v1's M4)
// this case fails at "0 rows for the minted session", and the audit line it also
// catches is the one that still claims GRANT-RECORDED.
func TestTicket224SessionVerbPutsARowOnDiskInTheAssembledRun(t *testing.T) {
	h := newReplyHost(t, 20*time.Second)
	target := filepath.Join(h.outside, "remembered-by-a-click.txt")
	pr, pw := io.Pipe()
	h.reply = pr
	t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })

	const task, corr = "t224r2-verb", "t224r2-verb-corr"
	var (
		minted    string
		printed   []string
		cardTool  string
		out       agent.ToolOutcome
		sessionOK bool
		done      = make(chan struct{})
	)
	h.rtHook = func(rt *agentRuntime) {
		if rt.session == nil {
			t.Errorf("the assembled run has no session ledger: nothing here can store a scope")
			return
		}
		minted = rt.session.SessionID().String()
		revoke := rt.gate.AdmitTextTask(task)
		defer revoke()
		args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"written after 本会话内允许"}`,
			filepath.ToSlash(target)))
		go func() {
			defer close(done)
			out, _ = rt.bridge.Execute(context.Background(), agent.ToolRequest{
				TaskID: task, CorrelationID: corr, CallID: "v1",
				Name: "fs.write", Args: args,
			})
		}()
		if !waitForCard187(rt.liveCards, corr, 10*time.Second) {
			t.Errorf("no card reached the native ledger, so there was nothing to answer")
			return
		}
		card, _ := rt.liveCards.look(corr)
		cardTool, printed = card.Tool, append([]string(nil), card.Paths...)
		if _, err := fmt.Fprintln(pw, "session "+corr); err != nil {
			t.Errorf("writing the reply: %v", err)
			return
		}
		<-done
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	<-done
	if out.IsError {
		t.Fatalf("the answered call must execute, got: %s", out.Text)
	}
	if len(printed) == 0 {
		t.Fatalf("the card printed no paths, so the row below has nothing to be matched against")
	}

	ctx := context.Background()
	st := h.openStore()
	rows, err := st.ListGrantsBySession(ctx, minted)
	if err != nil {
		t.Fatalf("ListGrantsBySession(%s): %v", minted, err)
	}
	// THE CELL: one answered card, one row per printed path. Zero rows here means
	// the user's 「本次会话内」 was a single-use allow wearing a remembered answer's
	// clothes, which is the exact bug this ticket was opened for.
	if len(rows) != len(printed) {
		t.Fatalf("%d approval_grant rows for the session this run minted (%s), want %d - the card "+
			"printed %v. The audit still says: %s",
			len(rows), minted, len(printed), printed, grantLinesOf(h))
	}
	got := rows[0]
	if got.Scope != memory.GrantScopeSession {
		t.Errorf("row scope = %q, want %q: SPEC-02 §3 has one scope value and 「长期」 is not it",
			got.Scope, memory.GrantScopeSession)
	}
	if got.Tool != cardTool {
		t.Errorf("row tool = %q, want the card's own %q", got.Tool, cardTool)
	}
	if got.Pattern != printed[0] {
		t.Errorf("row pattern = %q, want the path the card printed %q: a rule narrower or wider "+
			"than what was shown is a rule nobody agreed to", got.Pattern, printed[0])
	}
	if got.SessionID != minted {
		t.Errorf("row session_id = %q, want the id this boot minted %q", got.SessionID, minted)
	}
	if got.ExpiresAt <= got.CreatedAt {
		t.Errorf("row created_at=%d expires_at=%d: born dead, so Covering would never honour it",
			got.CreatedAt, got.ExpiresAt)
	}
	if got.RevokedAt != nil {
		t.Errorf("row revoked_at=%v, want NULL for a fresh answer", *got.RevokedAt)
	}

	// The row must be readable by the ledger that wrote it: a row nothing can find
	// is the other half of "silently grants nothing".
	sessionOK = func() bool {
		lt, err := memory.Open(h.dir)
		if err != nil {
			return false
		}
		defer lt.Close()
		r2, err := lt.ListGrantsBySession(ctx, minted)
		return err == nil && len(r2) == len(printed) && r2[0].ID == got.ID
	}()
	if !sessionOK {
		t.Errorf("a second handle on the same db cannot read the row back by the minted id")
	}

	// And the audit line names it with the row's real id, never with 0.
	if !strings.Contains(h.err.String(),
		fmt.Sprintf("approval: GRANT-RECORDED corr=%s grant_id=%d tool=%s", corr, got.ID, cardTool)) {
		t.Errorf("audit is missing the recording line naming grant_id=%d; got:\n%s",
			got.ID, grantLinesOf(h))
	}
	if strings.Contains(h.err.String(), "GRANT-DROPPED") {
		t.Errorf("audit says the scope was DROPPED on a host that does have a ledger:\n%s", grantLinesOf(h))
	}
	// The console owes the operator the same sentence it prints for any answer.
	if !strings.Contains(h.out.String(), corr) {
		t.Errorf("the console never mentioned this correlation id at all:\n%s", h.out.String())
	}
	if _, err := os.Stat(target); err != nil {
		t.Errorf("the answered call should have written %s: %v", target, err)
	}
}

// grantLinesOf pulls just the grant-shaped audit lines out of the run's stderr, so
// a failure prints the sentences that matter instead of a 400-line dump.
func grantLinesOf(h *replyHost) string {
	var kept []string
	for _, l := range strings.Split(h.err.String(), "\n") {
		if strings.Contains(l, "GRANT") || strings.Contains(l, "SESSION-MINT") ||
			strings.Contains(l, "REPLY ") {
			kept = append(kept, l)
		}
	}
	return strings.Join(kept, "\n")
}

// ---------------------------------------------------------------------------
// N#1读 - the assembled bridge stops asking, and books WHICH row paid for it
// ---------------------------------------------------------------------------

// TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking is AC#2's read cell at
// the one layer that had never run it: the bridge cmd/wisp assembled, reading
// through the ledger cmd/wisp minted, out to the store cmd/wisp opened. 224-v1 §1
// proved this by hand with an -overlay probe and §2 judged it "实现成立、判据缺";
// this is that judgement as a tracked case.
//
// Both halves of 「问与不问」 are measured, and by card counts rather than by
// errors: an L1 window that runs out answers allow (SPEC-06 §2), so "was I
// refused" says nothing about whether anything was asked. That is the polarity trap
// 224-v1 §6-b row 3 names, and the control below is what keeps it honest.
func TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking(t *testing.T) {
	h := newReplyHost(t, 20*time.Second)
	asked := filepath.Join(h.dir, "never-granted.txt")
	granted := filepath.Join(h.dir, "granted-in-this-session.txt")
	ctx := context.Background()

	type booked struct {
		decision string
		grantID  *int64
		cards    int
	}
	var (
		ledgerID string
		control  booked
		measured booked
		rowID    int64
		canonic  string
	)
	h.rtHook = func(rt *agentRuntime) {
		if rt.session == nil {
			t.Errorf("assembled run has no session ledger")
			return
		}
		ledgerID = rt.session.SessionID().String()
		revoke := rt.gate.AdmitTextTask("t224r2-read")
		defer revoke()

		// CONTROL first: with nothing on record, the same shape of call asks once.
		before := rt.windowCount()
		_, _, ctrlErr := t224r2call(rt, "t224r2-ask", asked)
		cardsAfterControl := rt.windowCount() - before
		control = booked{
			t224r2decision(t, rt.store, "t224r2-read", "t224r2-ask"),
			t224r2grantID(t, rt.store, "t224r2-read", "t224r2-ask"), cardsAfterControl,
		}
		if ctrlErr {
			t.Errorf("the control call errored: the L1 window running out means execute")
		}

		// Now answer-as-if-recorded through the PRODUCTION ledger, with the pattern
		// taken from the production canonicalizer (the same route the card's own
		// paths came by in the case above).
		c, err := rt.paths.Canonicalize(filepath.ToSlash(granted))
		if err != nil {
			t.Fatalf("Canonicalize(%q): %v", granted, err)
		}
		canonic = c
		id, err := rt.session.Record(ctx, "fs.write", c)
		if err != nil {
			t.Fatalf("session.Record: %v", err)
		}
		rowID = id

		before = rt.windowCount()
		_, _, measErr := t224r2call(rt, "t224r2-hit", granted)
		measured = booked{
			t224r2decision(t, rt.store, "t224r2-read", "t224r2-hit"),
			t224r2grantID(t, rt.store, "t224r2-read", "t224r2-hit"), rt.windowCount() - before,
		}
		if measErr {
			t.Errorf("the granted call errored: %v", measErr)
		}
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\n%s\n%s", code, h.out.String(), h.err.String())
	}

	// Non-vacuity: the chain does ask when no row covers it, and books no grant_id.
	if control.cards != 1 {
		t.Fatalf("the ungranted call showed %d cards, want 1: the control half is broken, so a "+
			"zero below would prove nothing", control.cards)
	}
	if control.grantID != nil {
		t.Errorf("the ungranted call booked grant_id=%d, want NULL", *control.grantID)
	}

	// THE CELL: the covered call shows no card at all.
	if measured.cards != 0 {
		t.Errorf("a live session grant left %d cards on screen for the covered call: 「命中授权就"+
			"不再弹卡」 is not happening in the assembled run (ledger id=%s, pattern=%q, row=%d)",
			measured.cards, ledgerID, canonic, rowID)
	}
	if measured.decision != agent.DecisionAllowGrant {
		t.Errorf("covered call booked decision=%q, want %q", measured.decision, agent.DecisionAllowGrant)
	}
	if measured.grantID == nil || *measured.grantID != rowID {
		t.Errorf("covered call booked grant_id=%v, want the covering row %d: tool_call.grant_id is "+
			"the column D45-2 promised for forensics", measured.grantID, rowID)
	}
	if !strings.Contains(h.err.String(),
		fmt.Sprintf("tools: GRANT-USE tool=fs.write paths=1 grant_id=%d assessed=L1", rowID)) {
		t.Errorf("audit is missing the GRANT-USE line naming row %d; got:\n%s", rowID, grantLinesOf(h))
	}
	if _, err := os.Stat(granted); err != nil {
		t.Errorf("the covered call should have written %s: %v", granted, err)
	}
	// The row is still the only one, still keyed to this boot's identity.
	rows, err := h.openStore().ListGrantsBySession(ctx, ledgerID)
	if err != nil {
		t.Fatalf("ListGrantsBySession: %v", err)
	}
	if len(rows) != 1 || rows[0].ID != rowID {
		t.Errorf("ledger rows = %+v, want exactly the one row %d", rows, rowID)
	}
}

// t224r2call dispatches one host-side fs.write through the bridge this boot
// assembled and reports whether it came back an error outcome.
func t224r2call(rt *agentRuntime, corr, path string) (agent.ToolOutcome, string, bool) {
	args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"through the assembled chain"}`,
		filepath.ToSlash(path)))
	out, err := rt.bridge.Execute(context.Background(), agent.ToolRequest{
		TaskID: "t224r2-read", CorrelationID: corr, CallID: "call-" + corr,
		Name: "fs.write", Args: args,
	})
	if err != nil {
		return out, err.Error(), true
	}
	return out, out.Text, out.IsError
}

// t224r2row reads one tool_call row off the real store by (task, correlation). The
// task id is a parameter rather than a constant because this file's two hosts book
// under different ones: its own dispatches use t224r2-read, while the restart case
// reuses run_mode101_test.go's t101call, which keys on t101-task-<n>.
func t224r2row(t *testing.T, st *memory.Store, task, corr string) memory.ToolCall {
	t.Helper()
	rows, err := st.ListToolCallsByTask(context.Background(), task)
	if err != nil {
		t.Fatalf("ListToolCallsByTask(%s): %v", task, err)
	}
	for _, r := range rows {
		if r.CorrelationID == corr {
			return r
		}
	}
	t.Fatalf("no tool_call row for task %q carries correlation %q; rows=%+v", task, corr, rows)
	return memory.ToolCall{}
}

func t224r2decision(t *testing.T, st *memory.Store, task, corr string) string {
	t.Helper()
	return t224r2row(t, st, task, corr).Decision
}

func t224r2grantID(t *testing.T, st *memory.Store, task, corr string) *int64 {
	t.Helper()
	return t224r2row(t, st, task, corr).GrantID
}

// ---------------------------------------------------------------------------
// N#2 - the cross-boot case three comments have been pointing at since 224
// ---------------------------------------------------------------------------

// TestTicket224ProductionSessionDoesNotSurviveRestart is the case
// internal/session/grants_test.go:20, :94 and :134 and internal/tools/grant_test.go:7
// all name. It carries the three controls ticket 224's ⑩ lists, and ⑩ says a
// missing one does not count as a judgement:
//
//	① two boots mint different ids          (asserted on both boots' real ledgers)
//	② rows written by boot 1's id read 0    (asserted with boot 2's minted id)
//	③ the same L1 call ASKS again           (asserted through the assembled bridge)
//
// It is also the expensive half of the test-literal lock: the ids here are the ones
// the production minting point produced, and neither equals
// "session-before-restart" or "session-after-restart", which is what makes the
// four pre-224 guards' hand-written fixtures unable to collide with a real session
// by accident.
//
// ⛔ Read the judgement ceiling at the top of this file before describing a green
// run below as "restart-proofed": boot 2 is a second assembly inside this same test
// process, not a second OS process, and nothing in this repository can tell those
// two statements apart yet.
func TestTicket224ProductionSessionDoesNotSurviveRestart(t *testing.T) {
	h := t101boot(t, "") // the default档, ask_every_step: a question definitely exists
	target := filepath.Join(h.dir, "cross-boot-note.txt")
	ctx := context.Background()

	// The literals the pre-224 guards key their fixture rows with. Named here
	// because the lock is about THESE strings, not about "some test string".
	const litBefore, litAfter = "session-before-restart", "session-after-restart"

	var (
		id1, id2  string
		rowID1    int64
		canonic   string
		covered   bool
		rowsIn1   int
		rowsIn2   int
		cardsBoot int
		decision  string
		grant     *int64
	)

	hook1 := func(rt *agentRuntime) {
		if rt.session == nil {
			t.Fatalf("boot 1 assembled without a session ledger, so there is no production id to measure")
		}
		id1 = rt.session.SessionID().String()
		if !rt.session.SessionID().Valid() {
			t.Fatalf("boot 1 minted %q, which the shape lock rejects: the minting point and the "+
				"lock disagree, so no literal can be excluded from either", id1)
		}
		if id1 == litBefore || id1 == litAfter {
			t.Fatalf("the production mint produced the test literal %q: the four pre-224 guards "+
				"would then be measuring a fixture against itself", id1)
		}
		c, err := rt.paths.Canonicalize(filepath.ToSlash(target))
		if err != nil {
			t.Fatalf("Canonicalize: %v", err)
		}
		canonic = c
		id, err := rt.session.Record(ctx, "fs.write", c)
		if err != nil {
			t.Fatalf("session.Record: %v", err)
		}
		rowID1 = id
		if _, ok := rt.session.Covering(ctx, "fs.write", []string{c}); !ok {
			t.Fatalf("boot 1's own session cannot cover its own answer: the fixture is broken, "+
				"so boot 2's refusal below would be vacuous (row %d, pattern %q)", id, c)
		}
		covered = true
		if rows, err := rt.store.ListGrantsBySession(ctx, id1); err != nil {
			t.Fatalf("ListGrantsBySession: %v", err)
		} else {
			rowsIn1 = len(rows)
		}
	}
	code1, log1 := h.start(t, nil, hook1)
	if code1 != 0 {
		t.Fatalf("boot 1 exit %d\n%s", code1, log1)
	}
	if !covered {
		t.Fatal("boot 1's ledger never covered its own row: nothing below can be attributed to a restart")
	}
	if rowsIn1 != 1 {
		t.Fatalf("boot 1's session holds %d rows, want 1", rowsIn1)
	}
	// The identity that keyed the row is the identity the minting point announced:
	// "which session wrote this row" stays answerable from the log, not inferred.
	if !strings.Contains(log1, "SESSION-MINT id="+id1) {
		t.Errorf("boot 1's audit does not announce the id it minted (%s):\n%s",
			id1, grantLinesIn(log1))
	}

	hook2 := func(rt *agentRuntime) {
		if rt.session == nil {
			t.Fatalf("boot 2 assembled without a session ledger")
		}
		id2 = rt.session.SessionID().String()
		// ① two boots, two ids - the control ticket 224's ⑩ says the suite lacked.
		if id2 == id1 {
			t.Fatalf("the second assembly minted boot 1's identity %q: 「本次会话内」 would then be "+
				"a permanent pass, since every row on disk is keyed by this string", id2)
		}
		if !rt.session.SessionID().Valid() {
			t.Fatalf("boot 2 minted %q, which the shape lock rejects", id2)
		}
		if id2 == litBefore || id2 == litAfter {
			t.Fatalf("the production mint produced the test literal %q", id2)
		}
		// ② boot 2's own id finds nothing.
		rows, err := rt.store.ListGrantsBySession(ctx, id2)
		if err != nil {
			t.Fatalf("ListGrantsBySession(%s): %v", id2, err)
		}
		rowsIn2 = len(rows)
		// The row is not deleted: invalidation is a key that stopped existing, not a
		// cleanup step somebody could forget (SPEC-02 §4 keeps it for the audit).
		if old, err := rt.store.ListGrantsBySession(ctx, id1); err != nil {
			t.Fatalf("ListGrantsBySession(%s): %v", id1, err)
		} else if len(old) != 1 || old[0].ID != rowID1 {
			t.Errorf("boot 1's row is not on disk any more (rows=%+v, want the one row %d)", old, rowID1)
		}
		// ③ and the behaviour half, through the real bridge: the same call now ASKS.
		// Card counts, not refusals - an L1 window running out answers allow, so a
		// refusal would say nothing about whether anything was asked.
		cardsBoot = rt.windowCount()
		_, _, _ = t101call(rt, 9, "fs.write", target)
		cardsBoot = rt.windowCount() - cardsBoot
		// t101call keys its rows on task "t101-task-<n>" with correlation
		// "t101-corr-<n>", so both halves name the same dispatch.
		grant = t224r2grantID(t, rt.store, "t101-task-9", "t101-corr-9")
		decision = t224r2decision(t, rt.store, "t101-task-9", "t101-corr-9")
	}
	code2, log2 := h.start(t, nil, hook2)
	if code2 != 0 {
		t.Fatalf("boot 2 exit %d\n%s", code2, log2)
	}

	if rowsIn2 != 0 {
		t.Errorf("the second assembly's session can read %d grant rows left by the first, want 0 "+
			"(PLAN.md:1642 「会话结束后授权必须失效」)", rowsIn2)
	}
	if cardsBoot < 1 {
		t.Errorf("AC#3 FAILS: after the restart the same L1 call showed %d cards, want at least 1 - "+
			"a session grant that survives into the next assembly silences the question without "+
			"anyone being asked again (pattern %q, row %d)", cardsBoot, canonic, rowID1)
	}
	if decision == agent.DecisionAllowGrant {
		t.Errorf("the post-restart call booked decision=allow_session_grant: the second assembly "+
			"acted on the first one's authorization (grant_id=%v)", grant)
	}
	if grant != nil {
		t.Errorf("the post-restart call booked grant_id=%d, want NULL", *grant)
	}

	// And the row itself survived for audit, keyed to the identity that died.
	st, err := memory.Open(h.dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	all, err := st.ListGrants(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 || all[0].SessionID != id1 || all[0].ID != rowID1 {
		t.Errorf("approval_grant rows after restart = %+v, want exactly the dead-session row %d "+
			"kept for the audit window", all, rowID1)
	}
	// Boot 2 announced its OWN id, and the two announcements are the whole case in
	// two lines: two mint points, two keys, zero rows in common.
	if !strings.Contains(log2, "SESSION-MINT id="+id2) {
		t.Errorf("boot 2's audit does not announce the id it minted (%s):\n%s", id2, grantLinesIn(log2))
	}
	if !strings.Contains(log2, "session: GRANT-READ") && strings.Contains(log2, "GRANT-HIT") {
		t.Errorf("boot 2 hit a grant row it should not have been able to see:\n%s", grantLinesIn(log2))
	}
}

// grantLinesIn trims a boot's merged output down to the lines this ticket's cases
// read, so a failure prints the sentences instead of a several-hundred-line dump.
func grantLinesIn(log string) string {
	var kept []string
	for _, l := range strings.Split(log, "\n") {
		if strings.Contains(l, "GRANT") || strings.Contains(l, "SESSION-") {
			kept = append(kept, l)
		}
	}
	return strings.Join(kept, "\n")
}
