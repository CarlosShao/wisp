package agent

import (
	"fmt"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/observe"
)

// 180 orphan-red repair for Test31DisarmStopsTheWait.
//
// That case is a coin flip on the package gate, not a slow machine: Arm spawns
// the wait, Disarm closes its stop channel, and Advance can release the timer
// before the spawned goroutine has ever reached its select. With both channels
// ready Go's select picks uniformly at random, so ~half of the runs went on to
// fire - MarkReminderFired answered applied=true (the row was still armed),
// a toast was delivered, and the "cancel removes" contract was broken. The
// single-shot original case reports it only when the coin lands that way (3 of
// 4 package runs measured, docs/evidence/s1/orphan-reds-internal-180.md).
//
// This case runs the same sequence many times so the verdict is deterministic:
// a delivered toast is not allowed in ANY round.
func Test180DisarmBeatsAnAlreadyReleasedWait(t *testing.T) {
	const rounds = 48 // one delivery is enough to fail; 48 x ~50% makes a miss ~1e-15
	for i := 0; i < rounds; i++ {
		id := fmt.Sprintf("rem_180_%d", i)
		clock := newFakeReminderClock()
		rows := newFakeRows(map[string]string{id: memory.ReminderArmed})
		ntf := &fakeNotifier{}
		j := &fakeJournal{}
		tm, err := NewReminderTimers(ReminderTimerOptions{
			Rows: rows, Notify: ntf, Journal: j, Clock: clock,
			Registry: observe.NewRegistry(), Root: observe.NewRoot("t180-disarm-race"),
		})
		if err != nil {
			t.Fatalf("round %d: NewReminderTimers: %v", i, err)
		}

		target := clock.WallNow().Add(time.Minute).Unix()
		if err := tm.Arm(memory.Reminder{ID: id, FireAt: target, Note: "已取消", Status: memory.ReminderArmed}); err != nil {
			t.Fatalf("round %d: Arm: %v", i, err)
		}
		if !tm.Disarm(id) {
			t.Fatalf("round %d: Disarm = false, want true", i)
		}

		// Release the wait AFTER the cancel: this is the order that used to
		// lose, because the release and the stop are both ready by the time the
		// goroutine wakes up.
		clock.Advance(2 * time.Minute)
		waitFor(t, 2*time.Second, func() bool { return tm.Pending() == 0 }, "the wait never exited")
		_ = tm.Close(2 * time.Second)

		if got := ntf.count(); got != 0 {
			t.Fatalf("round %d: 取消后仍发了 %d 条通知", i, got)
		}
		if got := rows.statusOf(id); got != memory.ReminderArmed {
			t.Fatalf("round %d: row status = %q, want armed (the cancel's own write is tools-side)", i, got)
		}
		if n := len(j.all()); n != 0 {
			t.Fatalf("round %d: audit rows = %d, want 0 (nothing happened, so nothing is recorded)", i, n)
		}
	}
}
