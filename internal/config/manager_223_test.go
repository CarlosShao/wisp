package config

// Ticket 223 AC#6's second half, measured at the config layer: the D36
// re-confirmation must not be a wait the Manager holds its own mutex through.
//
// Before this ticket manager.go reached the hook from inside apply(), which the
// caller held under mu, while the file's own header comment promised "Callbacks
// run OUTSIDE the lock so they may call Config()". Two things followed from that
// lie, and both are reddened below rather than argued:
//
//   - a hook that read its own config deadlocked (sync.Mutex is not reentrant:
//     Config() takes mu, and mu was already held);
//   - while the hook was parked on a human, every OTHER Config()/Resolved()
//     reader in the process was parked too - up to C18's 300s for one card.
//
// The plan/commit split in manager.go moved the ask between the two lock
// sections, which also makes the second half of D36 rule 1 observable: the
// pending loosening is NOT in cur while the card is unanswered.
//
// A bare `go` appears in this file on purpose. tools/d22scan's bans #1-5 skip
// *_test.go by design (walkGo, and 票 223's census records why: "测试里绿" proves
// nothing about production shape), and these two cases are probes for a
// deadlock - they must start a goroutine without joining a registry that would
// itself be part of what is under test. Production's tick goes through
// observe.Default.Spawn, where the ban bites.

import (
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// completedWithin runs fn on its own goroutine and reports whether it returned
// before budget elapsed. fn must hand back anything it observes through a
// buffered channel (never a shared variable), so an abandoned goroutine cannot
// race the assertions.
func completedWithin(budget time.Duration, fn func()) bool {
	done := make(chan struct{})
	go func() {
		defer close(done)
		fn()
	}()
	select {
	case <-done:
		return true
	case <-time.After(budget):
		return false
	}
}

// TestConfirmLockedRunsOutsideTheManagerLock is AC#6's "卡挂着的时候 Config()
// 仍可读" judgement, planted at the layer that owns the mutex.
func TestConfirmLockedRunsOutsideTheManagerLock(t *testing.T) {
	m, _, mutate := newTestManager(t, matrixBase)
	loosened := strings.Replace(matrixBase, "allowed_dirs = []", `allowed_dirs = ["D:\\data"]`, 1)

	type observation struct {
		selfReadOK   bool
		otherReadOK  bool
		pendingValue []string
	}
	observedCh := make(chan observation, 1)
	// What the hook saw through Config() while its own card was still open.
	var seenViaConfig []string
	selfCh := make(chan []string, 1)

	var once sync.Once
	release := make(chan struct{})
	t.Cleanup(func() { once.Do(func() { close(release) }) })

	m.ConfirmLocked = func(section string, keys []string) bool {
		// (1) the shape that used to self-deadlock: the hook reads the Manager.
		selfOK := completedWithin(3*time.Second, func() {
			selfCh <- m.Config().FS.AllowedDirs
		})
		// (2) an unrelated reader must not be frozen by the parked ask.
		otherOK := completedWithin(3*time.Second, func() {
			_ = m.Resolved()
		})
		select {
		case seenViaConfig = <-selfCh:
		default:
		}
		observedCh <- observation{
			selfReadOK: selfOK, otherReadOK: otherOK,
			pendingValue: seenViaConfig,
		}
		<-release // park on the "human"
		return true
	}

	mutate(loosened)

	type result struct {
		rep *Report
		err error
	}
	done := make(chan result, 1)
	go func() {
		rep, err := m.CheckAndReload()
		done <- result{rep, err}
	}()

	var o observation
	select {
	case o = <-observedCh:
	case <-time.After(10 * time.Second):
		once.Do(func() { close(release) })
		t.Fatal("the D36 hook was never asked - nothing to measure")
	}
	if !o.selfReadOK {
		once.Do(func() { close(release) })
		t.Fatal("ConfirmLocked could not call Config(): the hook still runs inside the Manager lock")
	}
	if !o.otherReadOK {
		once.Do(func() { close(release) })
		t.Fatal("Config()/Resolved() froze while a confirmation was parked")
	}
	// D36 rule 1's "不得静默生效", read off the live object while the card is up.
	if len(o.pendingValue) != 0 {
		once.Do(func() { close(release) })
		t.Fatalf("cur already carried the pending loosening while the card was open: %v", o.pendingValue)
	}

	once.Do(func() { close(release) })
	var res result
	select {
	case res = <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("CheckAndReload did not return after the confirmation was answered")
	}
	if res.err != nil {
		t.Fatalf("CheckAndReload: %v", res.err)
	}
	if res.rep == nil || len(res.rep.Locked) != 1 {
		t.Fatalf("report = %+v, want one locked decision", res.rep)
	}
	if !res.rep.Locked[0].Approved || res.rep.Locked[0].Direction != DirLoosen {
		t.Fatalf("report.Locked[0] = %+v, want an approved loosen", res.rep.Locked[0])
	}
	if got := m.Config().FS.AllowedDirs; len(got) != 1 || got[0] != `D:\data` {
		t.Fatalf("an approved loosening must be committed after the ask: %v", got)
	}
}

// TestCheckAndReloadAsksOncePerFileChange nails the serialization this ticket
// had to add: releasing mu for the length of a confirmation means two pollers
// reading the same new bytes would otherwise each plan, each ask, and each
// commit. One file change costs one card.
func TestCheckAndReloadAsksOncePerFileChange(t *testing.T) {
	m, _, mutate := newTestManager(t, matrixBase)
	loosened := strings.Replace(matrixBase, "allowed_dirs = []", `allowed_dirs = ["D:\\data"]`, 1)

	var asks atomic.Int32
	m.ConfirmLocked = func(string, []string) bool {
		asks.Add(1)
		return false
	}
	mutate(loosened)

	var wg sync.WaitGroup
	for range 3 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := m.CheckAndReload(); err != nil {
				t.Errorf("CheckAndReload: %v", err)
			}
		}()
	}
	wg.Wait()

	if got := asks.Load(); got != 1 {
		t.Fatalf("one file change raised %d confirmations, want exactly 1", got)
	}
	if got := m.Config().FS.AllowedDirs; len(got) != 0 {
		t.Fatalf("a denied loosening must keep the old values: %v", got)
	}
}
