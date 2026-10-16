package memory

// Ticket 326 AC#3: the `UPDATE session SET updated_at=?` that AppendSessionTurn
// now writes lives in the SAME transaction as the turn INSERT, so a hard kill
// mid-write leaves one of exactly two shapes on disk — either the turn and its
// touch are both there, or neither is. This is the real-process-kill shape
// ticket 166 AC#3 established (TestCrashRecoveryKillMidWrite is its sibling for
// task_log/tool_call); nothing here mocks the crash, the parent hard-kills
// (TerminateProcess) the writer while its next transaction is in flight.

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

const (
	crashChildEnv = "WISP_326_APPEND_CRASH_CHILD"
	crashDirEnv   = "WISP_326_APPEND_CRASH_DIR"
)

// TestAppendSessionTurnCrashChild is the CHILD role of
// TestAppendSessionTurnCrashAtomicAcrossKill: it opens the store the parent
// points it at, creates one conversation stamped an hour in the past, then
// appends turns forever, announcing each one AFTER AppendSessionTurn returned
// (so the announce is a committed point and the transaction in flight at kill
// time is the next one).
//
// The turn stamps are minted one second apart (base+i) instead of taken from the
// wall clock, because the whole loop finishes inside a single wall second: with
// equal stamps "no refresh at all" and "refresh in a separate transaction" both
// satisfy any timestamp comparison and the台件 would grade itself green for the
// wrong reason. A caller minting its own per-turn stamp is the production shape
// too — internal/session.Manager.AppendTurn passes m.now().Unix() in.
func TestAppendSessionTurnCrashChild(t *testing.T) {
	if os.Getenv(crashChildEnv) != "1" {
		t.Skip("crash child for TestAppendSessionTurnCrashAtomicAcrossKill")
	}
	s, err := Open(os.Getenv(crashDirEnv))
	if err != nil {
		t.Fatalf("child open: %v", err)
	}
	defer s.Close()
	ctx := context.Background()
	base := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC).Unix()
	birth := base - 3600
	if err := s.InsertSession(ctx, Session{
		ID: "crash-sess", Title: "the conversation killed mid-write", CreatedAt: birth, UpdatedAt: birth,
	}); err != nil {
		t.Fatalf("child InsertSession: %v", err)
	}
	for i := 0; ; i++ {
		if _, err := s.AppendSessionTurn(ctx, SessionTurn{
			SessionID: "crash-sess", Role: "user",
			Content:   fmt.Sprintf("turn %06d filler %s", i, strings.Repeat("x", 4096)),
			CreatedAt: base + int64(i),
		}); err != nil {
			t.Fatalf("child append %d: %v", i, err)
		}
		fmt.Printf("commit %d\n", i)
	}
}

// TestAppendSessionTurnCrashAtomicAcrossKill kills the child at five fixed
// committed-counts and three derived ones, reopens through the production path
// (WAL recovery + startup checkpoint) and classifies what it finds. Only
// NOTHING and ALL are legal. The shape that must be unreachable is a committed
// turn whose session row still carries an OLDER updated_at — that is what
// "the touch is not in the same transaction as the turn" looks like on disk after
// a kill, and it is also exactly what ticket 326 is about (a turn that does not
// count as a touch).
func TestAppendSessionTurnCrashAtomicAcrossKill(t *testing.T) {
	if os.Getenv(crashChildEnv) == "1" {
		t.Skip("parent role; the child is TestAppendSessionTurnCrashChild")
	}
	for _, killAfter := range []int{2, 4, 9, 20, 45} {
		killAndAudit(t, killAfter)
	}
	for _, seed := range []int64{1792117039198721100, 42, 20261016} {
		killAndAudit(t, int(seed%13)+3)
	}
}

func killAndAudit(t *testing.T, killAfter int) {
	t.Helper()
	dir := filepath.Join(sealableTempDir124(t), "data")

	cmd := exec.Command(os.Args[0], "-test.run=^TestAppendSessionTurnCrashChild$", "-test.timeout=120s")
	cmd.Env = append(os.Environ(), crashChildEnv+"=1", crashDirEnv+"="+dir)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("start child: %v", err)
	}

	var announced atomic.Int64
	scanDone := make(chan struct{})
	go func() {
		defer close(scanDone)
		sc := bufio.NewScanner(stdout)
		sc.Buffer(make([]byte, 64*1024), 1024*1024)
		for sc.Scan() {
			if strings.HasPrefix(sc.Text(), "commit ") {
				announced.Add(1)
			}
		}
	}()

	deadline := observe.NewTimeout(60 * time.Second)
	for int(announced.Load()) < killAfter {
		if deadline.Expired() {
			_ = cmd.Process.Kill()
			<-scanDone
			t.Fatalf("child announced %d commits within 60s, want >= %d (台件故障，不是判据)", announced.Load(), killAfter)
		}
		time.Sleep(2 * time.Millisecond)
	}
	if err := cmd.Process.Kill(); err != nil {
		t.Fatalf("kill child: %v", err)
	}
	_ = cmd.Wait() // killed: a non-zero exit is the point
	<-scanDone
	seen := int(announced.Load())

	s, err := Open(dir)
	if err != nil {
		t.Fatalf("reopen after kill: %v", err)
	}
	defer s.Close()

	var check string
	if err := s.reader.QueryRow(integrityCheckQuery).Scan(&check); err != nil {
		t.Fatalf("integrity_check: %v", err)
	}
	if check != "ok" {
		t.Errorf("integrity_check = %q after kill at %d commits, want ok", check, killAfter)
	}

	var (
		createdAt, updatedAt, turns, maxTurnStamp, maxSeq int64
		verdict                                           string
	)
	row := s.reader.QueryRow(`SELECT s.created_at, s.updated_at,
		(SELECT COUNT(*) FROM session_turn t WHERE t.session_id = s.id),
		(COALESCE((SELECT MAX(t.created_at) FROM session_turn t WHERE t.session_id = s.id), 0)),
		(COALESCE((SELECT MAX(t.seq) FROM session_turn t WHERE t.session_id = s.id), 0))
		FROM session s WHERE s.id='crash-sess'`)
	if err := row.Scan(&createdAt, &updatedAt, &turns, &maxTurnStamp, &maxSeq); err != nil {
		t.Fatalf("read crash-sess after kill: %v", err)
	}
	switch {
	case turns == 0 && updatedAt == createdAt:
		verdict = "NOTHING"
	case turns > 0 && updatedAt >= maxTurnStamp && turns == maxSeq:
		verdict = "ALL"
	default:
		verdict = "HALF-STATE"
	}
	t.Logf("326 AC#3 kill@%d: announced=%d turns=%d maxSeq=%d updated_at=%d created_at=%d max(turn.created_at)=%d lag=%ds => %s",
		killAfter, seen, turns, maxSeq, updatedAt, createdAt, maxTurnStamp, maxTurnStamp-updatedAt, verdict)

	if verdict == "HALF-STATE" {
		t.Errorf("half state after kill at %d commits: session_turn has %d rows (maxSeq %d, newest turn stamp %d) "+
			"while session.updated_at = %d — the touch and the turn did not share one transaction (want NOTHING or ALL, got a %ds lag or a torn seq)",
			killAfter, turns, maxSeq, maxTurnStamp, updatedAt, maxTurnStamp-updatedAt)
	}
	if int(turns) < seen {
		t.Errorf("committed turns lost across the kill: %d rows on disk, %d announced commits", turns, seen)
	}
}
