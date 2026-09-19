package memory

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
)

// isBusyErr detects SQLITE_BUSY leaking out of any layer (D35 rule 3: there
// must be none — the single writer plus WAL makes retries unnecessary, so a
// busy error is always a bug).
func isBusyErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "busy") || strings.Contains(msg, "database is locked")
}

// TestConcurrentWritersReaders runs 2 writers + 4 readers for 10 seconds
// (SPEC-02 §8) and asserts: zero busy errors, all writes committed, and the
// WAL bounded by the autocheckpoint setting.
func TestConcurrentWritersReaders(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()

	// Give search something to find.
	if _, err := s.AddMemory(ctx, "并发测试记忆", "并发 测试"); err != nil {
		t.Fatal(err)
	}

	const (
		writers = 2
		readers = 4
		runFor  = 10 * time.Second
	)
	var (
		wg     sync.WaitGroup
		errCh  = make(chan error, writers+readers)
		writes atomic.Int64
		reads  atomic.Int64
	)

	stop := make(chan struct{})

	// 2 writers: each task_log insert is a two-statement transaction through
	// the db-writer (the only write path).
	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; ; i++ {
				select {
				case <-stop:
					return
				default:
				}
				id := fmt.Sprintf("conc-%d-%d", w, i)
				err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
					if _, err := tx.ExecContext(ctx,
						`INSERT INTO task_log(id, started_at, state, query_text) VALUES(?,?, 'running', ?)`,
						id, observe.NowWallUTC().Unix(), "q"); err != nil {
						return err
					}
					_, err := tx.ExecContext(ctx,
						`INSERT INTO tool_call(task_id, seq, tool, args_json, risk_level, correlation_id)
						 VALUES(?,1,'fs.read','{}','L0',?)`,
						id, "corr-"+id)
					return err
				})
				if err != nil {
					errCh <- fmt.Errorf("writer %d: %w", w, err)
					return
				}
				writes.Add(1)
			}
		}(w)
	}

	// 4 readers: concurrent SELECT traffic on the reader pool (WAL: reads
	// never block the writer).
	for r := 0; r < readers; r++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				var n int
				err := s.reader.QueryRowContext(ctx, `SELECT count(*) FROM task_log`).Scan(&n)
				if err != nil {
					errCh <- fmt.Errorf("reader %d: %w", r, err)
					return
				}
				if _, err := s.SearchMemory(ctx, "并发", 5); err != nil {
					errCh <- fmt.Errorf("reader %d search: %w", r, err)
					return
				}
				reads.Add(1)
			}
		}(r)
	}

	// Monotonic 10s budget (observe.Timeout — D42#9).
	budget := observe.NewTimeout(runFor)
	for budget.Remaining() > 0 {
		time.Sleep(50 * time.Millisecond)
	}
	close(stop)
	wg.Wait()
	close(errCh)

	var busyCount int
	for err := range errCh {
		if isBusyErr(err) {
			busyCount++
		}
		t.Error(err)
	}
	if busyCount != 0 {
		t.Errorf("%d SQLITE_BUSY errors leaked to callers (single-writer invariant broken)", busyCount)
	}
	if writes.Load() == 0 {
		t.Fatal("writers made no progress")
	}
	if reads.Load() == 0 {
		t.Error("readers made no progress")
	}

	// Every committed write is visible.
	var n int64
	if err := s.reader.QueryRow(`SELECT count(*) FROM task_log`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != writes.Load() {
		t.Errorf("task_log rows = %d, writers committed %d", n, writes.Load())
	}

	// WAL bounded by wal_autocheckpoint=1000 pages (D35 rule 3).
	var pageSize int64
	if err := s.reader.QueryRow(`PRAGMA page_size`).Scan(&pageSize); err != nil {
		t.Fatal(err)
	}
	walPath := s.dbPath + "-wal"
	walSize := int64(0)
	if fi, err := os.Stat(walPath); err == nil {
		walSize = fi.Size()
	}
	maxExpected := pageSize * 1000 * 2 // 1000-page target, generous 2x headroom
	if walSize > maxExpected {
		t.Errorf("WAL size %d > %d (autocheckpoint not bounding growth)", walSize, maxExpected)
	}
	t.Logf("10s run: %d writes, %d read-cycles, wal=%d bytes (page=%d)",
		writes.Load(), reads.Load(), walSize, pageSize)

	// Clean exit leaves the WAL truncated (Close checkpoints TRUNCATE).
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if fi, err := os.Stat(walPath); err == nil && fi.Size() > 0 {
		t.Errorf("WAL not truncated on clean close: %d bytes", fi.Size())
	}
}

// --- Crash recovery with a real subprocess (SPEC-02 §2/§8: the WAL must
// survive kill / power-loss reopen; here: hard process kill mid-write). ---

// TestSubprocessCrashWriter is the CHILD of TestCrashRecoveryKillMidWrite:
// it writes continuously until killed. Runs only when the env marks it.
func TestSubprocessCrashWriter(t *testing.T) {
	if os.Getenv("WISP_CRASH_CHILD") != "1" {
		t.Skip("crash-writer subprocess; runs under TestCrashRecoveryKillMidWrite")
	}
	dir := os.Getenv("WISP_CRASH_DIR")
	s, err := Open(dir)
	if err != nil {
		t.Fatalf("child open: %v", err)
	}
	defer s.Close()
	ctx := context.Background()
	for i := 0; ; i++ {
		id := fmt.Sprintf("crash-%d", i)
		err := s.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO task_log(id, started_at, state, query_text) VALUES(?,?, 'running','crash test')`,
				id, observe.NowWallUTC().Unix()); err != nil {
				return err
			}
			_, err := tx.ExecContext(ctx,
				`INSERT INTO tool_call(task_id, seq, tool, args_json, risk_level, correlation_id)
				 VALUES(?,1,'fs.read','{}','L0',?)`, id, "corr-"+id)
			return err
		})
		if err != nil {
			t.Fatalf("child write: %v", err)
		}
		// Announce the COMMIT (after it happened) so the parent kills us at a
		// known-committed point while the next transaction is in flight.
		fmt.Printf("commit %d\n", i)
	}
}

// TestCrashRecoveryKillMidWrite kills the writer subprocess mid-write, then
// reopens: committed rows survive, integrity_check passes, and the startup
// checkpoint truncates the WAL.
func TestCrashRecoveryKillMidWrite(t *testing.T) {
	if os.Getenv("WISP_CRASH_CHILD") == "1" {
		t.Skip("parent test; the child role is TestSubprocessCrashWriter")
	}
	dir := filepath.Join(t.TempDir(), "data")

	cmd := exec.Command(os.Args[0], "-test.run=^TestSubprocessCrashWriter$", "-test.timeout=60s")
	cmd.Env = append(os.Environ(), "WISP_CRASH_CHILD=1", "WISP_CRASH_DIR="+dir)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("start subprocess: %v", err)
	}

	// Wait for the killAfter-th committed transaction (each line is printed
	// AFTER the commit), then hard-kill: the child is almost certainly inside
	// the next write transaction — the "write half done" scenario.
	const killAfter = 20
	var committed atomic.Int64
	scanDone := make(chan struct{})
	go func() {
		defer close(scanDone)
		sc := bufio.NewScanner(stdout)
		for sc.Scan() {
			if strings.HasPrefix(sc.Text(), "commit ") {
				committed.Add(1)
			}
		}
	}()
	deadline := observe.NewTimeout(30 * time.Second)
	for committed.Load() < killAfter {
		if deadline.Expired() {
			t.Fatalf("subprocess only committed %d rows within budget (want >= %d)", committed.Load(), killAfter)
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := cmd.Process.Kill(); err != nil {
		t.Fatalf("kill subprocess: %v", err)
	}
	_ = cmd.Wait() // exit error expected (killed)
	<-scanDone     // scanner drained after the pipe closed

	// The WAL sidecars exist at kill time.
	if _, err := os.Stat(dir + string(os.PathSeparator) + "wisp.db-wal"); err != nil {
		t.Logf("note: no -wal file at kill time: %v", err)
	}

	// Reopen: migration must recognize the schema, WAL replay must recover
	// every committed transaction, and the startup checkpoint (TRUNCATE)
	// must shrink the WAL.
	s, err := Open(dir)
	if err != nil {
		t.Fatalf("reopen after crash: %v", err)
	}
	defer s.Close()

	var check string
	if err := s.reader.QueryRow(integrityCheckQuery).Scan(&check); err != nil {
		t.Fatalf("integrity_check: %v", err)
	}
	if check != "ok" {
		t.Errorf("integrity_check = %q, want ok", check)
	}

	var rows int
	if err := s.reader.QueryRow(`SELECT count(*) FROM task_log`).Scan(&rows); err != nil {
		t.Fatal(err)
	}
	if rows < killAfter {
		t.Errorf("recovered %d committed rows, want >= %d (lost committed data)", rows, killAfter)
	}
	var toolRows int
	if err := s.reader.QueryRow(`SELECT count(*) FROM tool_call`).Scan(&toolRows); err != nil {
		t.Fatal(err)
	}
	if toolRows != rows {
		t.Errorf("transaction atomicity broken: task_log=%d tool_call=%d (two-statement txns)", rows, toolRows)
	}

	walSize := int64(0)
	if fi, err := os.Stat(s.dbPath + "-wal"); err == nil {
		walSize = fi.Size()
	}
	if walSize != 0 {
		t.Errorf("startup checkpoint did not truncate WAL: %d bytes", walSize)
	}
	t.Logf("crash recovery: %d committed rows recovered, integrity ok, wal=%d bytes", rows, walSize)
}
