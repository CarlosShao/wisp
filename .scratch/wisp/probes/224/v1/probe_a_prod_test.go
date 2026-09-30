package main

// 224-v1 ACCEPTANCE PROBE (adjudicator's own instrument, not part of the delivery).
//
// Why this file exists: both new test files cite a cmd/wisp case
// "TestTicket224ProductionSessionDoesNotSurviveRestart" three times
// (internal/session/grants_test.go:20/:94/:134 and internal/tools/grant_test.go:7),
// and that case does not exist:
//
//	go test ./cmd/wisp/ -count=1 -run Ticket224 -v
//	=> testing: warning: no tests to run ... [no tests to run]
//
// So the ticket's 必答 3 ("那 271 行有没有真的到达 cmd/wisp 装配出来的那条路，
// 不是只到测试") has no executor in the delivered tests, and the 派单增量第 2 条
// lock's "expensive direction" has no executor either. This probe measures the
// production path directly, over the same harness ticket 101's guard uses
// (t101boot / t101host.start / t101call), i.e. a real assembleRuntime().
//
// HOW an unsilenced L1 is obtained here (learned the hard way, two failed drafts):
//
//   - fs.read inside the allowlist is **L0** ("无规则命中（L0 直接执行）"), and L0
//     never reaches the grant check - so a read cannot measure AC#2's read cell.
//   - fs.write over an EXISTING file is **L2** ("rules_hit=[R1 R8] ... R8: 不可逆
//     操作（覆盖已有内容）"), and L2 is exactly what a session grant may never
//     cover (SPEC-06 §8.3 bullet 1 / PLAN.md:2148 D45-3) - so it cannot either.
//   - fs.write to a path that does NOT exist yet is **L1** (only R1's lower bound
//     fires), which is the one shape the seam is consulted for. The file is
//     therefore removed between dispatches so every measured call is the same L1.
//
// The last control in each probe is the scope line itself: with a live grant for
// that exact tool+path in place, the call that turns into an L2 (now the file
// exists) MUST still produce a card. That assertion is about "was a question
// asked", not about "was it refused", which is what keeps it out of 必答 2's trap:
// the L2 card here dies by "ANSWER-EXPIRED ... decision=timeout->reject", i.e. the
// refusal is the GATE's, and no test may read it as the permission layer's.
//
// Run:
//	export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
//	go test ./cmd/wisp/ -count=1 -run Probe224v1 -v \
//	  -overlay=.scratch/wisp/probes/224/v1/overlay.json

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/session"
)

// probe224L1Write dispatches one L1 write (path absent -> only R1 fires) and
// reports how many cards the composed gate had shown by then.
func probe224L1Write(rt *agentRuntime, n int, path string) int {
	probe224Drop(path)
	w, _, _ := t101call(rt, n, "fs.write", path)
	return w
}

func probe224Drop(path string) {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		panic(err)
	}
}

func probe224AuditLines(t *testing.T, tag, log string) {
	t.Helper()
	for _, ln := range strings.Split(log, "\n") {
		low := strings.ToLower(ln)
		if strings.Contains(low, "tools: call") || strings.Contains(low, "grant") ||
			strings.Contains(low, "session-mint") || strings.Contains(low, "answer-expired") {
			t.Logf("%s %s", tag, ln)
		}
	}
}

// TestProbe224v1ProductionGrantReachesAssembledBridge answers 必答 3's capability
// claim: does a row written through the assembly-minted ledger stop a card that
// the ASSEMBLED bridge would otherwise show?
func TestProbe224v1ProductionGrantReachesAssembledBridge(t *testing.T) {
	h := t101boot(t, "") // fresh install -> ask_every_step, the strictest档
	target := filepath.Join(h.dir, "note224.txt")
	probe224Drop(target) // must NOT exist: an existing file makes the call an L2

	var (
		ledgerNil      bool
		mintedID       string
		askedNoGrant   int
		askedWthGrant  int
		recordErr      error
		canonPattern   string
		rowsOwnSession int
		covID          int64
		covOK          bool
		rowPattern     string
		l2Cards        int
	)

	code, log := h.start(t, nil, func(rt *agentRuntime) {
		ledgerNil = rt.session == nil
		if rt.session == nil {
			return
		}
		mintedID = rt.session.SessionID().String()

		// (1) baseline: an unsilenced L1 write shows exactly one window.
		askedNoGrant = probe224L1Write(rt, 1, target)
		probe224Drop(target)

		c, err := rt.paths.Canonicalize(filepath.ToSlash(target))
		if err != nil {
			t.Errorf("Canonicalize(%q): %v", target, err)
			return
		}
		canonPattern = c

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := rt.session.Record(ctx, "fs.write", c); err != nil {
			recordErr = err
			return
		}
		rows, err := rt.store.ListGrantsBySession(ctx, mintedID)
		if err != nil {
			t.Errorf("ListGrantsBySession: %v", err)
			return
		}
		rowsOwnSession = len(rows)
		if len(rows) > 0 {
			rowPattern = rows[0].Pattern
			covID, covOK = rt.session.Covering(ctx, "fs.write", []string{c})
		}

		// (2) the same L1 call again: does the assembled bridge stop asking?
		askedWthGrant = probe224L1Write(rt, 2, target)

		// (3) SCOPE CONTROL, same boot, same live grant: now the file exists, so
		// the identical call is an L2 (R8) and must still cost a card.
		before := rt.windowCount()
		_, _, _ = t101call(rt, 3, "fs.write", target)
		l2Cards = rt.windowCount() - before
	})
	if code != 0 {
		t.Fatalf("boot 1 exit %d\n%s", code, log)
	}
	probe224AuditLines(t, "BOOT1", log)
	if ledgerNil {
		t.Fatalf("AC#1/必答3 FAILS: the production assembly minted no session ledger "+
			"(rt.session == nil); internal/session's 271 lines reach nothing in cmd/wisp.\n%s", log)
	}
	if !session.ID(mintedID).Valid() {
		t.Errorf("the id the assembly minted is %q, which session.ID.Valid rejects: the "+
			"composition root is minting something the ledger calls illegal", mintedID)
	}
	// 派单增量第 2 条, expensive direction: the PRODUCTION id is never a fixture literal.
	for _, lit := range []string{"session-before-restart", "session-after-restart"} {
		if mintedID == lit {
			t.Errorf("the production mint equals the test literal %q: AC#4's control would "+
				"go green while granting nothing", lit)
		}
	}
	if recordErr != nil {
		t.Fatalf("production ledger refused the answer it is there to record: %v", recordErr)
	}
	if rowsOwnSession != 1 {
		t.Errorf("rows keyed to the minted id = %d, want 1", rowsOwnSession)
	}
	if rowPattern != canonPattern {
		t.Errorf("stored pattern %q != canonical %q: the matcher compares strings, so a write "+
			"side storing a different spelling than the read side asks with can NEVER hit",
			rowPattern, canonPattern)
	}
	if !covOK || covID == 0 {
		t.Errorf("the production ledger does not cover its own row: Covering=(%d,%v), want (>0,true)", covID, covOK)
	}
	if askedNoGrant < 1 {
		t.Errorf("non-vacuity control broken: the un-granted L1 call showed %d windows, so the "+
			"delta below would prove nothing", askedNoGrant)
	}
	if delta := askedWthGrant - askedNoGrant; delta != 0 {
		t.Errorf("AC#2 读 FAILS at the assembly root: after a live grant for %q the assembled "+
			"bridge still showed %d more window(s), want 0 (「命中授权就不再弹卡」)", canonPattern, delta)
	}
	if l2Cards != 1 {
		t.Errorf("scope line FAILS at the assembly root: an L2 (R8 overwrite) call with a live "+
			"session grant for the same tool+path showed %d cards, want 1 (SPEC-06 §8.3 bullet 1)", l2Cards)
	}
	t.Logf("PROBE-A minted=%s canon=%s windows_before=%d windows_after_total=%d rows=%d covering=(%d,%v) l2_cards=%d",
		mintedID, canonPattern, askedNoGrant, askedWthGrant, rowsOwnSession, covID, covOK, l2Cards)
}

// TestProbe224v1ProductionSessionDiesWithTheProcess answers AC#3 at the one level
// AC#3 names: a second boot over the same bytes on disk.
func TestProbe224v1ProductionSessionDiesWithTheProcess(t *testing.T) {
	h := t101boot(t, "")
	target := filepath.Join(h.dir, "note224.txt")
	probe224Drop(target)
	ctx := context.Background()

	var (
		id1, id2       string
		boot1Pre       int
		boot1Post      int
		boot1L2        int
		askedBoot2     int
		rowsForID1     int
		rowsForID2     int
		grantRowsBoot1 int
		boot1NoLedger  bool
		boot2NoLedger  bool
	)

	code, log := h.start(t, nil, func(rt *agentRuntime) {
		if rt.session == nil {
			boot1NoLedger = true
			return
		}
		id1 = rt.session.SessionID().String()
		boot1Pre = probe224L1Write(rt, 1, target)
		probe224Drop(target)
		c, err := rt.paths.Canonicalize(filepath.ToSlash(target))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := rt.session.Record(ctx, "fs.write", c); err != nil {
			t.Fatalf("Record: %v", err)
		}
		boot1Post = probe224L1Write(rt, 2, target)
		before := rt.windowCount()
		_, _, _ = t101call(rt, 3, "fs.write", target) // now an L2 -> card must appear
		boot1L2 = rt.windowCount() - before
		rows, err := rt.store.ListGrantsBySession(ctx, id1)
		if err != nil {
			t.Fatal(err)
		}
		grantRowsBoot1 = len(rows)
	})
	if code != 0 {
		t.Fatalf("boot 1 exit %d, no-ledger=%v\n%s", code, boot1NoLedger, log)
	}
	probe224AuditLines(t, "P-B1", log)
	if boot1NoLedger || grantRowsBoot1 != 1 {
		t.Fatalf("boot 1 setup failed (exit %d, rows %d, no-ledger=%v)\n%s", code, grantRowsBoot1, boot1NoLedger, log)
	}
	if boot1Pre < 1 {
		t.Fatalf("boot 1's un-granted L1 call showed %d windows: the fixture is broken", boot1Pre)
	}
	if boot1Post-boot1Pre != 0 {
		t.Fatalf("boot 1's own grant did not silence the L1 window (delta %d): the restart "+
			"half below would be measuring a feature that never worked", boot1Post-boot1Pre)
	}
	if boot1L2 != 1 {
		t.Fatalf("boot 1's L2 control showed %d cards, want 1", boot1L2)
	}

	// ---- restart: a new process surface over the same bytes ----
	probe224Drop(target) // keep boot 2's call an L1 too, so the card is about the grant
	code, log = h.start(t, nil, func(rt *agentRuntime) {
		if rt.session == nil {
			boot2NoLedger = true
			return
		}
		id2 = rt.session.SessionID().String()
		rows, err := rt.store.ListGrantsBySession(ctx, id1)
		if err != nil {
			t.Fatal(err)
		}
		rowsForID1 = len(rows)
		rows, err = rt.store.ListGrantsBySession(ctx, id2)
		if err != nil {
			t.Fatal(err)
		}
		rowsForID2 = len(rows)
		askedBoot2 = probe224L1Write(rt, 4, target)
	})
	if code != 0 {
		t.Fatalf("boot 2 exit %d\n%s", code, log)
	}
	probe224AuditLines(t, "P-B2", log)
	if boot2NoLedger {
		t.Fatalf("AC#1 FAILS at the assembly root: boot 2 minted no session ledger, so "+
			"「本会话内允许」 has no executor in production.\n%s", log)
	}

	if id1 == id2 {
		t.Errorf("AC#3 FAILS: two boots minted the SAME id %q - that is a derived "+
			"identity (A435 第 1 条), and it makes 「本会话内允许」 a permanent pass", id1)
	}
	if rowsForID1 != 1 {
		t.Errorf("the dead session's row is gone (%d rows): invalidation must not mean delete "+
			"(SPEC-02 §4 keeps it 30 days for audit)", rowsForID1)
	}
	if rowsForID2 != 0 {
		t.Errorf("AC#3 FAILS: boot 2's session can read %d grant rows left by boot 1, want 0", rowsForID2)
	}
	if askedBoot2 < 1 {
		t.Errorf("AC#3 FAILS: after a restart the same L1 call showed %d windows, want >=1 "+
			"(「重新弹卡」 is the whole point of this AC)", askedBoot2)
	}
	t.Logf("PROBE-B id1=%s id2=%s rows_id1=%d rows_id2=%d boot1_delta=%d boot1_l2=%d boot2_windows=%d",
		id1, id2, rowsForID1, rowsForID2, boot1Post-boot1Pre, boot1L2, askedBoot2)
}
