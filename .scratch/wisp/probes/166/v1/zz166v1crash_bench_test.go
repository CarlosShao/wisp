package memory

// 166-v1 acceptance bench (ticket 166 AC#3): a REAL process kill in the middle
// of the 2->3 migration step, plus a reopen. Nothing here is production code;
// this file lives only in the out-of-repo exported tree
// (/d/tmp/wisp-166v1/tree-head/internal/memory/) and is archived as a reading
// under .scratch/wisp/probes/166/v1/. It asserts: after the kill the database is
// either fully v2 (nothing of v3 visible) or fully v3 - never a half shape - and
// the rolled-back case replays whole on the next production Open.
//
// The kill is real: the parent runs os.Process.Kill (TerminateProcess on
// Windows) on a separate child process, at a deterministic point announced by
// the child over stdout. The child's transaction is therefore neither committed
// nor rolled back by anyone - the process simply stops existing with an open
// write transaction and uncommitted WAL frames.

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	v166EnvRole = "WISP166V1_ROLE"
	v166EnvDir  = "WISP166V1_DIR"
	v166EnvStmt = "WISP166V1_STMT"
)

func v166Say(msg string) {
	fmt.Println(msg)
	_ = os.Stdout.Sync()
}

// v166BlockForever parks the child until the parent kills it.
func v166BlockForever() {
	for {
		time.Sleep(50 * time.Millisecond)
	}
}

// TestZZ166V1CrashChild is the CHILD half. It runs the production Open path on a
// database the parent already brought to a genuine v2, and stops dying at one of
// three announced points inside the real 2->3 step.
func TestZZ166V1CrashChild(t *testing.T) {
	role := os.Getenv(v166EnvRole)
	if role == "" {
		t.Log("166-v1 bench: not the child role (no " + v166EnvRole + "); nothing to do")
		return
	}
	dir := os.Getenv(v166EnvDir)
	stmtIdx := -1
	if s := os.Getenv(v166EnvStmt); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil {
			t.Fatalf("bad %s=%q: %v", v166EnvStmt, s, err)
		}
		stmtIdx = n
	}

	saved := migrationChain
	defer func() { migrationChain = saved }()

	step3 := migrationStep{from: 2, to: 3, apply: func(tx txExec) error {
		switch {
		case strings.HasPrefix(role, "midstep"):
			stmts := splitSQLStatements(ddlV3)
			v166Say(fmt.Sprintf("V166-PLAN role=%s stmts=%d", role, len(stmts)))
			for i, stmt := range stmts {
				if _, err := tx.ExecContext(context.Background(), stmt); err != nil {
					return err
				}
				if i == stmtIdx {
					// Inside the still-open transaction: the real ddlV3 statements
					// 0..i have run, the watermark row has NOT been written yet.
					v166Say(fmt.Sprintf("V166-READY %s after-statement=%d/%d", role, i+1, len(stmts)))
					v166BlockForever()
				}
			}
			return nil
		case role == "precommit":
			if err := applyV3(tx); err != nil {
				return err
			}
			// Whole real applyV3 ran on the open tx; applyMigrationStep still has
			// to write schema_version and commit. Die here.
			v166Say("V166-READY precommit all-ddl-done")
			v166BlockForever()
			return nil
		case role == "spill" || strings.HasPrefix(role, "spill"):
			// Same as precommit, but the transaction is then made BIG so dirty pages
			// really leave the page cache and land as uncommitted frames in wisp.db-wal
			// before the kill. Without this, a kill inside the step could be waved
			// through as "nothing was ever written"; with it the reopen has to prove
			// WAL recovery throws the frames away.
			if err := applyV3(tx); err != nil {
				return err
			}
			if _, err := tx.ExecContext(context.Background(), `PRAGMA cache_size=-64`); err != nil {
				v166Say("V166-PRAGMA-WARN " + err.Error())
			}
			blob := strings.Repeat("x", 4096)
			limit := 3000
			if s := strings.TrimPrefix(role, "spill="); s != role {
				if n, err := strconv.Atoi(s); err == nil {
					limit = n
				}
			}
			for i := 0; i < limit; i++ {
				if _, err := tx.ExecContext(context.Background(),
					`INSERT INTO session_turn(session_id, seq, role, content, created_at) VALUES(?,?,?,?,?)`,
					"v166-spill", i, "user", blob, i); err != nil {
					return err
				}
				if i%500 == 0 {
					v166Say(fmt.Sprintf("V166-SPILL-PROGRESS %d/%d", i, limit))
				}
			}
			v166Say(fmt.Sprintf("V166-READY spill rows=%d", limit))
			v166BlockForever()
			return nil

		default:
			return applyV3(tx)
		}
	}}
	migrationChain = []migrationStep{
		{from: 0, to: 1, apply: applyV1},
		{from: 1, to: 2, apply: applyV2},
		step3,
	}

	s, err := Open(dir)
	if err != nil {
		v166Say("V166-CHILD-OPEN-ERR " + err.Error())
		t.Fatalf("child Open: %v", err)
	}
	if role == "postcommit" {
		// The 2->3 step is committed (applyMigrationStep returned). Die now, with
		// the store still open: this is the "全成" side of the boundary.
		v166Say("V166-READY postcommit committed-open-store")
		v166BlockForever()
	}
	v166Say("V166-CHILD-CLEAN-EXIT")
	_ = s.Close()
}

// v166Shape is what a raw reopen (no migration code involved) sees on disk.
type v166Shape struct {
	watermark  int
	hasSession bool
	hasTurn    bool
	cols       map[string][]string
	indexes    []string
	integrity  string
}

func v166ReadShape(t *testing.T, dbPath string) v166Shape {
	t.Helper()
	out := v166Shape{cols: map[string][]string{}}
	db, err := sql.Open(DriverName, dsn(dbPath, false))
	if err != nil {
		t.Fatalf("raw open %s: %v", dbPath, err)
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var raw string
	if err := db.QueryRowContext(ctx,
		`SELECT value FROM schema_meta WHERE key='schema_version'`).Scan(&raw); err != nil {
		t.Fatalf("read watermark of %s: %v", dbPath, err)
	}
	out.watermark, err = strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		t.Fatalf("watermark %q: %v", raw, err)
	}
	for _, pair := range []struct {
		name string
		dst  *bool
	}{{"session", &out.hasSession}, {"session_turn", &out.hasTurn}} {
		var n int
		if err := db.QueryRowContext(ctx,
			`SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?`, pair.name).Scan(&n); err != nil {
			t.Fatalf("hasTable %s: %v", pair.name, err)
		}
		*pair.dst = n > 0
	}
	for _, tbl := range []string{"session", "session_turn"} {
		rows, err := db.QueryContext(ctx, `SELECT name FROM pragma_table_info(?)`, tbl)
		if err != nil {
			t.Fatalf("table_info %s: %v", tbl, err)
		}
		var cols []string
		for rows.Next() {
			var c string
			if err := rows.Scan(&c); err != nil {
				t.Fatal(err)
			}
			cols = append(cols, c)
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		rows.Close()
		out.cols[tbl] = cols
	}
	irows, err := db.QueryContext(ctx,
		`SELECT name FROM sqlite_master WHERE type='index' AND name LIKE 'idx_session%' ORDER BY name`)
	if err != nil {
		t.Fatalf("index roster: %v", err)
	}
	for irows.Next() {
		var n string
		if err := irows.Scan(&n); err != nil {
			t.Fatal(err)
		}
		out.indexes = append(out.indexes, n)
	}
	irows.Close()
	var ic string
	if err := db.QueryRowContext(ctx, integrityCheckQuery).Scan(&ic); err != nil {
		t.Fatalf("integrity_check: %v", err)
	}
	out.integrity = ic
	return out
}

func (s v166Shape) String() string {
	return fmt.Sprintf("wm=%d session=%t session_turn=%t cols(session)=%d cols(session_turn)=%d idx=%v integrity=%s",
		s.watermark, s.hasSession, s.hasTurn, len(s.cols["session"]), len(s.cols["session_turn"]), s.indexes, s.integrity)
}

// classify: all-or-nothing, exactly as AC#3 asks.
func (s v166Shape) classify() (class string, legal bool) {
	full := s.watermark == 3 && s.hasSession && s.hasTurn &&
		len(s.cols["session"]) == 7 && len(s.cols["session_turn"]) == 6 &&
		len(s.indexes) == 3 && s.integrity == "ok"
	none := s.watermark == 2 && !s.hasSession && !s.hasTurn &&
		len(s.cols["session"]) == 0 && len(s.cols["session_turn"]) == 0 &&
		len(s.indexes) == 0 && s.integrity == "ok"
	switch {
	case full:
		return "ALL(v3-complete)", true
	case none:
		return "NOTHING(v2-untouched)", true
	default:
		return "HALF-STATE", false
	}
}

// v166RunChildToKill starts the child, waits for its V166-READY marker, then
// hard-kills it. Returns the child's captured stdout and the Wait error text.
func v166RunChildToKill(t *testing.T, dir, role string, stmt int) (string, string) {
	t.Helper()
	roleArg := role
	if role == "midstep" {
		roleArg = fmt.Sprintf("midstep:%d", stmt)
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestZZ166V1CrashChild$", "-test.timeout=180s")
	cmd.Env = append(os.Environ(),
		v166EnvRole+"="+roleArg,
		v166EnvDir+"="+dir,
		v166EnvStmt+"="+strconv.Itoa(stmt),
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	var capture []string
	if err := cmd.Start(); err != nil {
		t.Fatalf("start child: %v", err)
	}
	scanner := bufio.NewScanner(stdout)
	ready := make(chan string, 64)
	done := make(chan struct{})
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			capture = append(capture, line)
			if strings.HasPrefix(line, "V166-READY") || strings.HasPrefix(line, "V166-CHILD-OPEN-ERR") {
				select {
				case ready <- line:
				default:
				}
			}
		}
		close(done)
	}()
	select {
	case marker := <-ready:
		if strings.HasPrefix(marker, "V166-CHILD-OPEN-ERR") {
			t.Fatalf("child failed before the kill point: %s", marker)
		}
	case <-time.After(120 * time.Second):
		cmd.Process.Kill()
		t.Fatalf("child never announced a kill point for role=%s", roleArg)
	}
	if err := cmd.Process.Kill(); err != nil {
		t.Fatalf("Process.Kill: %v", err)
	}
	waitErr := cmd.Wait()
	// drain whatever the pipe still holds (the child is dead; this returns fast)
	select {
	case <-done:
	case <-time.After(5 * time.Second):
	}
	werr := "<nil>"
	if waitErr != nil {
		werr = waitErr.Error()
	}
	return strings.Join(capture, "\n"), werr
}

func v166DirListing(dir string) string {
	var lines []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "readdir error: " + err.Error()
	}
	for _, e := range entries {
		if e.IsDir() {
			sub, _ := os.ReadDir(filepath.Join(dir, e.Name()))
			for _, s := range sub {
				info, _ := s.Info()
				size := int64(-1)
				if info != nil {
					size = info.Size()
				}
				lines = append(lines, fmt.Sprintf("%s/%s size=%d", e.Name(), s.Name(), size))
			}
			continue
		}
		info, _ := e.Info()
		size := int64(-1)
		if info != nil {
			size = info.Size()
		}
		lines = append(lines, fmt.Sprintf("%s size=%d", e.Name(), size))
	}
	return strings.Join(lines, "\n")
}

// TestZZ166V1CrashParent drives the five kill points and prints every reading.
func TestZZ166V1CrashParent(t *testing.T) {
	cases := []struct {
		name, role string
		stmt       int
	}{
		{"midstep-after-stmt0-create-table-session", "midstep", 0},
		{"midstep-after-stmt2-index-workspace", "midstep", 2},
		{"midstep-after-stmt3-create-table-session_turn", "midstep", 3},
		{"precommit-all-ddl-run-watermark-missing", "precommit", -1},
		{"spill-uncommitted-wal-frames-then-kill", "spill", -1},
		{"postcommit-committed-then-killed", "postcommit", -1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := filepath.Join(sealableTempDir124(t), "data")
			// 1. A genuine v2 database, written through the production path, with a
			//    real row in it so "the rollback ate user data too" would be visible.
			s2, err := Open(dir, withSchemaTarget(2))
			if err != nil {
				t.Fatalf("open v2: %v", err)
			}
			if hasTable(t, s2.reader, "session") {
				t.Fatal("premise broken: v2 target already created the session table")
			}
			ctx := context.Background()
			if err := s2.write(ctx, func(ctx context.Context, tx *sql.Tx) error {
				_, err := tx.ExecContext(ctx,
					`INSERT INTO profile(slot, content, source, updated_at) VALUES('v166.slot','seeded-by-166-v1','manual',?)`,
					time.Now().UTC().Unix())
				return err
			}); err != nil {
				t.Fatalf("seed v2: %v", err)
			}
			if err := s2.Close(); err != nil {
				t.Fatalf("close v2: %v", err)
			}
			pre := v166ReadShape(t, filepath.Join(dir, dbFileName))
			fmt.Printf("V166-BEFORE %s | %s\n", tc.name, pre)

			// 2. Kill the migrating process at the announced point.
			out, werr := v166RunChildToKill(t, dir, tc.role, tc.stmt)
			fmt.Printf("V166-KILL %s | wait-error=%s\nchild-log:\n%s\n", tc.name, werr, out)
			fmt.Printf("V166-DIRLIST %s\n%s\n", tc.name, v166DirListing(dir))

			// 3. Reopen WITHOUT any migration code: read the shape off the disk.
			shape := v166ReadShape(t, filepath.Join(dir, dbFileName))
			class, legal := shape.classify()
			bak := filepath.Join(dir, "backup", "wisp.db.bak-2-3")
			_, bakErr := os.Stat(bak)
			fmt.Printf("V166-AFTER %s | class=%s legal=%t pre-backup-2-3=%t | %s\n",
				tc.name, class, legal, bakErr == nil, shape)
			if !legal {
				t.Errorf("%s: HALF-STATE after a real kill: %s", tc.name, shape)
			}
			if tc.role != "postcommit" && class == "ALL(v3-complete)" {
				t.Errorf("%s: the step was supposed to die uncommitted yet the disk shows v3", tc.name)
			}
			if tc.role == "postcommit" && class != "ALL(v3-complete)" {
				t.Errorf("%s: committed step must survive the kill, disk shows %s", tc.name, class)
			}

			// 4. Production reopen must land on a complete v3 and the tables must work.
			s3, err := Open(dir)
			if err != nil {
				t.Fatalf("production reopen: %v", err)
			}
			final := v166ReadShape(t, filepath.Join(dir, dbFileName))
			fClass, fLegal := final.classify()
			fmt.Printf("V166-REOPEN %s | class=%s legal=%t | %s\n", tc.name, fClass, fLegal, final)
			if !fLegal {
				t.Errorf("%s: reopen did not reach a complete v3: %s", tc.name, final)
			}
			var seed string
			if err := s3.reader.QueryRowContext(ctx,
				`SELECT content FROM profile WHERE slot='v166.slot'`).Scan(&seed); err != nil {
				t.Errorf("%s: seeded row lost across the kill+reopen: %v", tc.name, err)
			} else {
				fmt.Printf("V166-SEED %s | profile content=%q\n", tc.name, seed)
			}
			if err := s3.InsertSession(ctx, Session{ID: "v166-" + tc.name, Title: "after kill", Workspace: "wk"}); err != nil {
				t.Errorf("%s: session table not usable after reopen: %v", tc.name, err)
			}
			if got, err := s3.ListSessions(ctx, SessionListOptions{Workspace: "wk"}); err != nil || len(got) != 1 {
				t.Errorf("%s: ListSessions after reopen = %d rows (%v)", tc.name, len(got), err)
			}
			_ = s3.Close()
		})
	}
}

// v166KillAfterDelay spawns the migrating child and hard-kills it after exactly
// d, WITHOUT waiting for any marker from it - the death point is not
// coordinated with the code under test at all.
func v166KillAfterDelay(t *testing.T, dir, role string, d time.Duration) (string, string) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestZZ166V1CrashChild$", "-test.timeout=180s")
	cmd.Env = append(os.Environ(),
		v166EnvRole+"="+role,
		v166EnvDir+"="+dir,
		v166EnvStmt+"=-1",
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	var mu sync.Mutex
	var capture []string
	if err := cmd.Start(); err != nil {
		t.Fatalf("start child: %v", err)
	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			mu.Lock()
			capture = append(capture, line)
			mu.Unlock()
		}
	}()
	time.Sleep(d)
	if err := cmd.Process.Kill(); err != nil {
		t.Fatalf("Process.Kill: %v", err)
	}
	waitErr := cmd.Wait()
	time.Sleep(200 * time.Millisecond)
	mu.Lock()
	joined := strings.Join(capture, "\n")
	mu.Unlock()
	werr := "<nil>"
	if waitErr != nil {
		werr = waitErr.Error()
	}
	return joined, werr
}

// TestZZ166V1CrashSweep kills the migrating process at 12 uncoordinated moments
// (seeded, printed) and asserts every single outcome is all-or-nothing.
func TestZZ166V1CrashSweep(t *testing.T) {
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	fmt.Printf("V166-SWEEP-SEED %d\n", seed)
	half := 0
	for i := 0; i < 12; i++ {
		delay := time.Duration(150+rng.Intn(2850)) * time.Millisecond
		name := fmt.Sprintf("sweep%02d", i)
		dir := filepath.Join(sealableTempDir124(t), "data")
		s2, err := Open(dir, withSchemaTarget(2))
		if err != nil {
			t.Fatalf("open v2: %v", err)
		}
		if err := s2.Close(); err != nil {
			t.Fatal(err)
		}
		out, werr := v166KillAfterDelay(t, dir, "spill=4000", delay)
		ready := strings.Contains(out, "V166-READY")
		progress := strings.Count(out, "V166-SPILL-PROGRESS")
		bakErr := func() error {
			_, e := os.Stat(filepath.Join(dir, "backup", "wisp.db.bak-2-3"))
			return e
		}()
		shape := v166ReadShape(t, filepath.Join(dir, dbFileName))
		class, legal := shape.classify()
		wal := int64(-1)
		if st, e := os.Stat(filepath.Join(dir, dbFileName+"-wal")); e == nil {
			wal = st.Size()
		}
		fmt.Printf("V166-SWEEP %s | delay=%s wait=%s killed-midstep=%t progress-lines=%d reached-ready=%t step-started(backup-2-3)=%t wal=%d | class=%s legal=%t | %s\n",
			name, delay, werr, !ready, progress, ready, bakErr == nil, wal, class, legal, shape)
		if !legal {
			t.Errorf("%s (delay=%s): HALF-STATE after an uncoordinated kill: %s", name, delay, shape)
		}
		if !ready {
			half++
		}
		s3, err := Open(dir)
		if err != nil {
			t.Fatalf("%s: reopen: %v", name, err)
		}
		final, fLegal := v166ReadShape(t, filepath.Join(dir, dbFileName)).classify()
		fmt.Printf("V166-SWEEP-REOPEN %s | class=%s legal=%t\n", name, final, fLegal)
		if !fLegal {
			t.Errorf("%s: reopen did not reach a complete v3 (%s)", name, final)
		}
		_ = s3.Close()
	}
	fmt.Printf("V166-SWEEP-SUMMARY kills-inside-the-step=%d/12\n", half)
}

// TestZZ166V1NewVsOldRoster answers AC#3's second half: the numeric watermark as
// the idempotency key, pushed once for a NEW database (0->1->2->3) and once for

// an OLD one (a v2 file written before ticket 166, then 2->3). Same instrument,
// whole roster (tables + columns + indexes) printed for both and compared.
func TestZZ166V1NewVsOldRoster(t *testing.T) {
	ctx := context.Background()
	dump := func(t *testing.T, dir, label string) string {
		t.Helper()
		db, err := sql.Open(DriverName, dsn(filepath.Join(dir, dbFileName), false))
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		var sb strings.Builder
		var wm string
		if err := db.QueryRowContext(ctx,
			`SELECT value FROM schema_meta WHERE key='schema_version'`).Scan(&wm); err != nil {
			t.Fatal(err)
		}
		sb.WriteString(fmt.Sprintf("[%s] schema_version=%s\n", label, wm))
		rows, err := db.QueryContext(ctx,
			`SELECT type, name, sql FROM sqlite_master WHERE name NOT LIKE 'sqlite_%' ORDER BY type DESC, name`)
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var typ, name string
			var sqlText sql.NullString
			if err := rows.Scan(&typ, &name, &sqlText); err != nil {
				t.Fatal(err)
			}
			var cols string
			if typ == "table" {
				crows, err := db.QueryContext(ctx, `SELECT name FROM pragma_table_info(?)`, name)
				if err != nil {
					t.Fatal(err)
				}
				var parts []string
				for crows.Next() {
					var c string
					if err := crows.Scan(&c); err != nil {
						t.Fatal(err)
					}
					parts = append(parts, c)
				}
				crows.Close()
				cols = " (" + strings.Join(parts, ",") + ")"
			}
			sb.WriteString(fmt.Sprintf("  %s %s%s\n", typ, name, cols))
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		fmt.Print(sb.String())
		return sb.String()
	}

	// NEW: an empty dir migrated to the production target.
	newDir := filepath.Join(sealableTempDir124(t), "data")
	sn, err := Open(newDir)
	if err != nil {
		t.Fatalf("open fresh: %v", err)
	}
	newRoster := dump(t, newDir, "NEW-0to3")
	_ = sn.Close()

	// OLD: a v2 database written by the pre-166 target, then reopened.
	oldDir := filepath.Join(sealableTempDir124(t), "data")
	so, err := Open(oldDir, withSchemaTarget(2))
	if err != nil {
		t.Fatalf("open v2: %v", err)
	}
	before := dump(t, oldDir, "OLD-at-v2")
	_ = so.Close()
	so2, err := Open(oldDir)
	if err != nil {
		t.Fatalf("reopen old db: %v", err)
	}
	oldRoster := dump(t, oldDir, "OLD-2to3")
	_ = so2.Close()

	// Compare the roster BODIES: the first line of each dump carries the label and
	// the watermark, which is exactly what the two paths are supposed to share
	// after the upgrade, so it is dropped before the byte compare.
	body := func(s string) string {
		if i := strings.IndexByte(s, '\n'); i >= 0 {
			return s[i:]
		}
		return s
	}
	if body(newRoster) != body(oldRoster) {
		t.Errorf("new-vs-old rosters differ:\n--- NEW ---\n%s\n--- OLD ---\n%s", newRoster, oldRoster)
	}
	fmt.Printf("V166-ROSTER-CMP body(new)==body(old)=%t new-lines=%d old-lines=%d\n",
		body(newRoster) == body(oldRoster), strings.Count(newRoster, "\n"), strings.Count(oldRoster, "\n"))
	fmt.Printf("V166-ROSTER-OLD-BEFORE\n%s\n", before)
}

var _ = io.Discard
