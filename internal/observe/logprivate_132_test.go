package observe

// Ticket 132: the log tree is private data (owner's ruling on Q-31, 2026-09-23:
// "算、要锁") and the product has to seal it. What these cases assert is the
// READING on the objects that actually landed - the descriptor, or the mode, of
// the directory and the file the pipeline made - and never "SealDir was called".
// The distinction is the ticket's own: a seal that runs is not a seal that holds,
// and a call site nobody can observe is the exact shape that shipped unsealed.
//
// The two assertions are deliberately split by the instant they are taken at:
// the directory case creates a directory that did not exist, and the file case
// reads the boot file while it still holds zero bytes. An instrument that reads
// after the first line cannot tell "sealed before the write" from "sealed after
// it" (ticket 132's c3 leg booked that shape as a green test measuring air).

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestLogDirIsSealedAtCreation132r1 - the log directory does not exist when
// InitLog runs, so the pipeline is what creates it. Sealing afterwards is not
// enough (winsec.PrivateDirAll's doc: a child inherits the ACL it had at
// creation), which is why the fixture refuses to run against an existing
// directory: this case is only worth anything on a level the pipeline made.
func TestLogDirIsSealedAtCreation132r1(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "logs")
	if _, err := os.Stat(dir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("fixture is wrong: %s already exists, so this case would read a directory the pipeline did not create", dir)
	}
	p, err := InitLogWithRegistry(LogConfig{Dir: dir}, NewRegistry())
	if err != nil {
		t.Fatalf("InitLogWithRegistry: %v", err)
	}
	t.Cleanup(func() { _ = p.Close() })

	assertSealedDir132(t, dir)
}

// TestLogFileIsSealedBeforeTheFirstContentByte132r1 - the boot-time open happens
// inside newRollingWriter, before any record exists, so the file is empty when
// this reading is taken. That is the instant AC#2 names ("首次写入前"), and the
// size check first is what keeps the reading on that instant: a file that
// already holds a line cannot answer the question the ticket asks.
func TestLogFileIsSealedBeforeTheFirstContentByte132r1(t *testing.T) {
	dir := t.TempDir()
	p, err := InitLogWithRegistry(LogConfig{Dir: dir, Level: "info"}, NewRegistry())
	if err != nil {
		t.Fatalf("InitLogWithRegistry: %v", err)
	}
	t.Cleanup(func() { _ = p.Close() })

	files := logFileNames(t, dir)
	if len(files) != 1 {
		t.Fatalf("want exactly the boot file, got %d: %v", len(files), files)
	}
	path := filepath.Join(dir, files[0])
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() != 0 {
		t.Fatalf("%s already holds %d bytes, so this reading is not the pre-content instant", files[0], info.Size())
	}
	assertSealedFile132(t, path)

	// Then prove this is the file the pipeline really writes to, so the reading
	// above was taken on the object that carries the records rather than on a
	// file that was sealed and then abandoned.
	rec := slog.NewRecord(time.Now(), slog.LevelInfo, "ticket132: sealed before this line", 0)
	if err := p.Handler().Handle(context.Background(), rec); err != nil {
		t.Fatalf("write after the seal: %v", err)
	}
	if err := p.FlushNow(); err != nil {
		t.Fatal(err)
	}
	if info, err = os.Stat(path); err != nil || info.Size() == 0 {
		t.Fatalf("the sealed file took no content (err=%v), so this case reads the wrong object", err)
	}
	assertSealedFile132(t, path)
}

// TestExistingWideLogTreeIsRepaired132r1 - the repair half of
// winsec.PrivateDirAll's contract, and the state every machine that ran
// yesterday's build is in: a <data>\logs tree that already exists and is already
// wide. Sealing only on the start that creates the directory would leave that
// tree wide forever, because the directory is created exactly once. The widening
// is done by the platform helper and the fixture refuses to proceed if it did not
// land, so a green here cannot be a green on a tree that was never wide.
func TestExistingWideLogTreeIsRepaired132r1(t *testing.T) {
	dir := t.TempDir()
	widenDir132(t, dir)

	p, err := InitLogWithRegistry(LogConfig{Dir: dir}, NewRegistry())
	if err != nil {
		t.Fatalf("InitLogWithRegistry over a wide but writable tree: %v", err)
	}
	t.Cleanup(func() { _ = p.Close() })

	assertSealedDir132(t, dir)
}

// TestSealFailureRefusesTheWrite132r1 - AC#2's failure direction: a seal that
// cannot be applied must cost the record, never the guarantee. "Unsealed but
// writable" is the one outcome the ticket forbids, so the reading is "InitLog
// failed AND no line landed", not "an error was logged somewhere".
//
// The instrument replaces the seal at the instant the pipeline calls it, so what
// turns red is this route's ORDER rather than a test that writes files on its
// own. Removal of the seal call leaves the swapped function unreferenced and this
// case green-by-default, which is why the assertion is on the error rather than
// on the swap.
func TestSealFailureRefusesTheWrite132r1(t *testing.T) {
	dir := t.TempDir()
	restore := sealLogFile
	sealLogFile = func(string) error { return errors.New("132r1: seal refused by the instrument") }
	t.Cleanup(func() { sealLogFile = restore })

	p, err := InitLogWithRegistry(LogConfig{Dir: dir}, NewRegistry())
	if err == nil {
		_ = p.Close()
		t.Fatal("InitLogWithRegistry reported success with a seal that always fails: an unsealed log file stayed writable, which is the outcome AC#2 forbids")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "seal") {
		t.Fatalf("the refusal does not name the seal, so a reader cannot tell which act refused it: %v", err)
	}
	for _, name := range logFileNames(t, dir) {
		info, statErr := os.Stat(filepath.Join(dir, name))
		if statErr != nil {
			continue
		}
		if info.Size() != 0 {
			t.Errorf("%s holds %d bytes although its seal failed; the failure direction is 'refuse the bytes', not 'land them wider'", name, info.Size())
		}
	}
}
