package config

// Ticket 42-r1's finding (R-42r1-1), fixed in internal/config/manager.go:
// CheckAndReload's failure path adopted the seen stat for EVERY LoadFile
// failure, so one transient sharing violation - the Windows shape of racing
// the atomic-rename writer (temp file + rename, the C19 shape) - made the poll
// "see" the new mtime+size, fail the read, adopt the failed stat, and hide
// that write until the NEXT successful write. 42-r1's concurrent-poll evidence
// (internal/watchdog/config_poll_test.go) had to replant edits to work around
// exactly this ("a real writer would simply write again"); this file is the
// direct pin that the poll itself now recovers on the next tick.
//
// The two classes the fix splits (manager.go classifies on loader.go's
// readError marker):
//
//   - READ class (first case): the file could not be read off the disk. The
//     seen stat is NOT adopted, so the next tick retries the read and a real
//     writer's write becomes visible without being rewritten.
//   - CONTENT class (second case): the bytes are bad. Adoption stays - it is
//     the anti-loop that keeps a file that will not fix itself from
//     re-erroring every tick - and this case is the regression pin that the
//     read rule was not broadened into "never adopt".
//
// HOW THE SHARING VIOLATION IS MADE (deterministically, no sleeping): the
// test holds a CreateFile handle with dwShareMode=0 for the exact window it
// needs. Go's os.Stat survives that handle - GetFileAttributesEx answers
// ERROR_SHARING_VIOLATION and Go's stat falls back to FindFirstFile
// (go/src/os/stat_windows.go, the ERROR_SHARING_VIOLATION branch) - while
// os.ReadFile's CreateFile(GENERIC_READ, FILE_SHARE_READ|FILE_SHARE_WRITE)
// conflicts with the 0-share handle and fails. That is precisely the
// production shape: the stat succeeds, the read fails. It was measured by a
// standalone probe before this file was written, and it is re-proven here by
// the first case's own assertions: if the conflict ever stops producing
// ERROR_SHARING_VIOLATION on the read, the case fails loudly instead of
// passing on a different code path (the stat-failure branch adopts nothing
// and would make the recovery assertion pass vacuously).

import (
	"errors"
	"os"
	"strings"
	"testing"

	"golang.org/x/sys/windows"
)

// holdExclusive opens path with dwShareMode=0 (share nothing) and returns an
// idempotent release func. It is the credible Windows shape of the rename
// race: the holder needs no write access and no timing, and the conflict it
// creates is the same errno the production poll loses to.
func holdExclusive(t *testing.T, path string) (release func()) {
	t.Helper()
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		t.Fatal(err)
	}
	h, err := windows.CreateFile(p, windows.GENERIC_READ, 0, nil,
		windows.OPEN_EXISTING, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		t.Fatalf("exclusive open (dwShareMode=0) of %s: %v", path, err)
	}
	var done bool
	return func() {
		if done {
			return
		}
		done = true
		_ = windows.CloseHandle(h)
	}
}

// TestReadFailureDoesNotAdoptStatTheNextTickSeesTheWrite is the fix's direct
// pin: a read that fails to a transient sharing violation must leave the seen
// stat alone, so the very next tick retries the read and applies the write.
// Under the pre-fix always-adopt behavior this case is red at the marked
// assertion: the failed call adopts the new stat, the next tick answers
// (nil, nil), and the write stays hidden until the next successful write.
func TestReadFailureDoesNotAdoptStatTheNextTickSeesTheWrite(t *testing.T) {
	m, path, mutate := newTestManager(t, matrixBase)
	beforeMtime, beforeSize, beforeValid := m.seenMtime, m.seenSize, m.seenValid

	// A real write lands on disk while the poll is not looking (the production
	// shape: writer done, fingerprint changed, no poll ran in between).
	mutate(strings.Replace(matrixBase, "size = 56", "size = 60", 1))

	// One transient sharing violation, held across exactly one poll.
	release := holdExclusive(t, path)
	defer release()

	st, err := os.Stat(path)
	if err != nil {
		t.Fatalf("os.Stat must survive the sharing conflict (Go stat falls back to FindFirstFile): %v", err)
	}
	if st.ModTime().Equal(beforeMtime) {
		t.Fatal("the planted write did not move the mtime fingerprint; this case would measure nothing")
	}

	rep, err := m.CheckAndReload()
	if err == nil {
		t.Fatalf("the read under a sharing conflict must surface an error, got report %+v", rep)
	}
	if rep != nil {
		t.Fatalf("a failed read must not produce a report, got %+v", rep)
	}
	// The failure must be the READ class, not some other branch: the
	// stat-failure branch also adopts nothing, so a case that ended there
	// would pass the recovery assertion below vacuously.
	var re *readError
	if !errors.As(err, &re) {
		t.Fatalf("the failure must be the read class (loader.go's readError), got: %v", err)
	}
	if !errors.Is(err, windows.ERROR_SHARING_VIOLATION) {
		t.Fatalf("the read failure must carry ERROR_SHARING_VIOLATION off the chain, got: %v", err)
	}

	// THE FIX, read directly off the watch state: the failed read adopted
	// nothing - the seen triple is still the pre-write one.
	if !m.seenMtime.Equal(beforeMtime) || m.seenSize != beforeSize || !m.seenValid || !beforeValid {
		t.Fatalf("the failed read adopted the stat it saw: seen=(mtime=%v size=%d valid=%v), "+
			"want the pre-write (mtime=%v size=%d valid=%v) - adopting a failed read's stat is the R-42r1-1 loss",
			m.seenMtime, m.seenSize, m.seenValid, beforeMtime, beforeSize, beforeValid)
	}

	// The writer's handle goes away (the violation was transient).
	release()

	// Next tick: the poll retries the read and the write becomes visible.
	rep2, err2 := m.CheckAndReload()
	if err2 != nil {
		t.Fatalf("the tick after the conflict cleared: %v", err2)
	}
	if rep2 == nil {
		t.Fatal("THE WRITE WAS HIDDEN: the tick after a transient sharing violation reported no change - " +
			"the failed read adopted the stat, the pre-fix behavior this case exists to bite")
	}
	if got := m.Config().Ball.Size; got != 60 {
		t.Fatalf("ball.size = %d, want 60: the write that raced the poll must apply on the next tick "+
			"(hot tier, without the writer having to write again)", got)
	}
}

// TestContentFailureStillAdoptsStatTheNextTickStaysQuiet is the regression pin
// for the other half of the classification: a CONTENT failure (the bytes are
// bad - the ticket 223 shape) must keep adopting the seen stat, so the next
// tick over the same bytes is a silent no-op instead of the same error again.
// This is the assertion that bites if the read-class rule is ever broadened
// into "never adopt" (the anti-loop would re-error every tick here).
func TestContentFailureStillAdoptsStatTheNextTickStaysQuiet(t *testing.T) {
	m, path, mutate := newTestManager(t, matrixBase)

	// Content class: the bytes themselves are broken.
	broken := "schema_version = 2\n\nthis is not toml [[["
	mutate(broken)
	st, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	rep, err := m.CheckAndReload()
	if err == nil {
		t.Fatal("corrupt content must surface an error")
	}
	if rep != nil {
		t.Fatalf("a failed load must not produce a report, got %+v", rep)
	}
	var re *readError
	if errors.As(err, &re) {
		t.Fatalf("broken bytes must be content class, not read class: %v", err)
	}

	// Adoption (pre-existing behavior this fix must keep): the failed load
	// adopted the stat it saw...
	if !m.seenMtime.Equal(st.ModTime()) || m.seenSize != st.Size() || !m.seenValid {
		t.Fatalf("the content failure did not adopt the stat: seen=(mtime=%v size=%d valid=%v), "+
			"want the planted (mtime=%v size=%d valid=true) - the anti-loop depends on this adoption",
			m.seenMtime, m.seenSize, m.seenValid, st.ModTime(), st.Size())
	}
	// ...so the next tick over the SAME bytes is a silent no-op, not the same
	// error again.
	rep2, err2 := m.CheckAndReload()
	if err2 != nil || rep2 != nil {
		t.Fatalf("a content failure must adopt the stat so the next tick stays quiet (the anti-loop), "+
			"got rep=%v err=%v", rep2, err2)
	}
}
