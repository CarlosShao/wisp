package scheduler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// Ticket 47 acceptance legs. The fixture discipline mirrors the repo's
// existing concurrency tests: every wait is bounded by observe.NewTimeout
// (monotonic; D42#9), no wall-clock difference anywhere, and the goroutine
// assertions read the D38b registry - the same source of truth the watchdog
// reads.

// tc is a no-op Exec body returning a completed result.
func tc(text string) func(context.Context, *TaskRun) TaskResult {
	return func(context.Context, *TaskRun) TaskResult {
		return TaskResult{Status: "completed", Text: text}
	}
}

// waitFor polls cond with a deadline; the failure names what was waited for.
func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	tm := observe.NewTimeout(5 * time.Second)
	for !cond() {
		if tm.Expired() {
			t.Fatalf("等不到 %s（5s 上限）", what)
		}
		time.Sleep(2 * time.Millisecond)
	}
}

// ---------------------------------------------------------------------------
// AC#1 - concurrency: FIFO accept, concurrent run, trio names, drain
// ---------------------------------------------------------------------------

// Test47AC1FourTasksRunConcurrentlyAndRosterDrains is AC#1: 4 tasks with
// mixed tool loads all complete, and the goroutine count returns to the
// baseline (zero leaks) - read through the registry's own RosterReport.
func Test47AC1FourTasksRunConcurrentlyAndRosterDrains(t *testing.T) {
	reg := observe.NewRegistry()
	locks := NewPathLock(PathLockOptions{})
	s := New(Options{Registry: reg, Locks: locks, MaxConcurrent: 4})

	// Each task spawns its trio helpers (tool-exec-<id> + approval-waiter)
	// under its own root, holds until released, then finishes - the shape
	// the loop's per-task roster produces (agent-task-<id> is the root's own
	// Spawn name).
	var rel []chan struct{}
	for i := 0; i < 4; i++ {
		rel = append(rel, make(chan struct{}))
	}
	ids := make([]string, 4)
	for i := 0; i < 4; i++ {
		i, r := i, rel[i]
		id, err := s.Submit(TaskSpec{
			Label: fmt.Sprintf("任务%d", i),
			Exec: func(ctx context.Context, run *TaskRun) TaskResult {
				// Trio helpers: both named per the D38b per-task family and
				// spawned under the task root.
				done1 := run.Reg.Spawn("tool-exec-"+run.ID, "scheduler-test", run.Root, func(context.Context) {
					<-r
				})
				done2 := run.Reg.Spawn("approval-waiter", "scheduler-test", run.Root, func(context.Context) {
					<-r
				})
				<-done1.Done()
				<-done2.Done()
				return TaskResult{Status: "completed", Text: "完事"}
			},
		})
		if err != nil {
			t.Fatalf("Submit %d: %v", i, err)
		}
		ids[i] = id
	}

	// All four run concurrently (each spawned its helpers, none finished).
	waitFor(t, "四任务的三件套并发到场", func() bool {
		return reg.CountByName("approval-waiter") == 4
	})

	// Mixed loads complete in FIFO acceptance order of release.
	close(rel[2])
	close(rel[0])
	close(rel[3])
	close(rel[1])
	for _, id := range ids {
		res, err := s.Wait(id)
		if err != nil {
			t.Fatalf("Wait(%s): %v", id, err)
		}
		if res.Status != "completed" {
			t.Errorf("任务 %s 落到 %s, want completed（%+v）", id, res.Status, res)
		}
	}

	// Zero leaks: every per-task goroutine drained. The registry reports
	// exactly zero live goroutines (this test's registry sees only what the
	// scheduler and the exec bodies spawned).
	waitFor(t, "名册排干（零泄漏）", func() bool { return reg.Count() == 0 })
	rep := reg.RosterReport()
	if rep.Total != 0 || rep.ResidentOverBaseline || len(rep.Unknown) != 0 {
		t.Errorf("收尾名册读数 = %+v, want 全零（AC#1 的零泄漏判据）", rep)
	}
}

// Test47AC1FIFOAcceptOrderIsTheQueueOrder pins the FIFO accept half: four
// tasks, one slot - they run strictly in acceptance order.
func Test47AC1FIFOAcceptOrderIsTheQueueOrder(t *testing.T) {
	reg := observe.NewRegistry()
	s := New(Options{Registry: reg, MaxConcurrent: 1})

	var mu sync.Mutex
	var order []string
	gates := []chan struct{}{make(chan struct{}), make(chan struct{}), make(chan struct{})}
	var ids [3]string
	for i := range gates {
		i := i
		id, err := s.Submit(TaskSpec{
			Label: fmt.Sprintf("FIFO%d", i),
			Exec: func(context.Context, *TaskRun) TaskResult {
				<-gates[i]
				mu.Lock()
				order = append(order, fmt.Sprint(i))
				mu.Unlock()
				return TaskResult{Status: "completed"}
			},
		})
		if err != nil {
			t.Fatalf("Submit %d: %v", i, err)
		}
		ids[i] = id
	}
	// The first task holds the only slot; the other two are Queued.
	waitFor(t, "两枚排队行可见", func() bool {
		rows := s.Rows()
		return len(rows) == 3 && rows[1].State == statemachine.StateQueued && rows[2].State == statemachine.StateQueued
	})
	close(gates[0])
	if _, err := s.Wait(ids[0]); err != nil {
		t.Fatalf("Wait 0: %v", err)
	}
	waitFor(t, "第二枚起跑", func() bool {
		rows := s.Rows()
		return len(rows) == 3 && rows[1].State == statemachine.StateThinking
	})
	close(gates[1])
	if _, err := s.Wait(ids[1]); err != nil {
		t.Fatalf("Wait 1: %v", err)
	}
	waitFor(t, "第三枚起跑", func() bool {
		rows := s.Rows()
		return len(rows) == 3 && rows[2].State == statemachine.StateThinking
	})
	close(gates[2])
	if _, err := s.Wait(ids[2]); err != nil {
		t.Fatalf("Wait 2: %v", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if strings.Join(order, ",") != "0,1,2" {
		t.Errorf("FIFO 完成序 = %v, want [0 1 2]", order)
	}
}

// Test47AC1DrainIsRequiredBeforeCompletion pins the D38c completion law: a
// task whose helpers have not drained is NOT complete even after its Exec
// returned; completion lands only after the drain.
func Test47AC1DrainIsRequiredBeforeCompletion(t *testing.T) {
	reg := observe.NewRegistry()
	s := New(Options{Registry: reg, MaxConcurrent: 1})

	release := make(chan struct{})
	id, err := s.Submit(TaskSpec{
		Label: "排干",
		Exec: func(ctx context.Context, run *TaskRun) TaskResult {
			// A helper the task spawned that outlives the Exec's return
			// until released - exactly the shape the drain requirement
			// exists to make observable.
			h := run.Reg.Spawn("approval-waiter", "scheduler-test", run.Root, func(context.Context) {
				<-release
			})
			go func() { <-h.Done() }() // the task does not join it before returning
			//nolint: the spawned reader joins the helper's done on a
			// registry-spawned goroutine below.
			return TaskResult{Status: "completed"}
		},
	})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	// Give the helper a beat to spawn, then wait for the task.
	time.Sleep(50 * time.Millisecond)
	done := make(chan TaskResult, 1)
	go func() {
		res, _ := s.Wait(id)
		done <- res
	}()
	select {
	case res := <-done:
		t.Fatalf("helper 未排干任务就报告完成：%+v（D38c 的「完成＝WaitGroup 排干」被跳过）", res)
	case <-time.After(200 * time.Millisecond):
		// Expected: Wait is still blocked because the helper holds the root.
	}
	close(release)
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("释放 helper 后 Wait 仍未返回：排干判据把完成卡死了")
	}
}

// ---------------------------------------------------------------------------
// AC#2 - PathLock: overlap queues the second task with the visible reason
// ---------------------------------------------------------------------------

// Test47AC2OverlappingPathQueuesSecondWithReasonAndResumes is AC#2: the
// later task queues with 「在等 <task> 释放 <path>」, shows Queued, and
// resumes when the holder releases (D43 #20 -> #39).
func Test47AC2OverlappingPathQueuesSecondWithReasonAndResumes(t *testing.T) {
	reg := observe.NewRegistry()
	locks := NewPathLock(PathLockOptions{})
	s := New(Options{Registry: reg, Locks: locks, MaxConcurrent: 2})

	const p = `C:/wisp47/工件.txt`
	holderStarted := make(chan struct{})
	releaseHolder := make(chan struct{})

	held := make(chan string, 1)
	idA, err := s.Submit(TaskSpec{
		Label: "A 占路",
		Exec: func(ctx context.Context, run *TaskRun) TaskResult {
			if err := run.Locks.Acquire(ctx, run.ID, []string{p}); err != nil {
				return TaskResult{Status: "failed", Err: err.Error()}
			}
			held <- run.ID
			close(holderStarted)
			<-releaseHolder
			run.Locks.Release(run.ID, p)
			return TaskResult{Status: "completed"}
		},
	})
	if err != nil {
		t.Fatalf("Submit A: %v", err)
	}
	<-holderStarted
	aID := <-held

	var idB string
	bQueued := make(chan struct{})
	go func() {
		var err error
		idB, err = s.Submit(TaskSpec{
			Label: "B 撞路",
			Exec: func(ctx context.Context, run *TaskRun) TaskResult {
				if err := run.Locks.Acquire(ctx, run.ID, []string{p}); err != nil {
					return TaskResult{Status: "failed", Err: err.Error()}
				}
				run.Locks.Release(run.ID, p)
				return TaskResult{Status: "completed"}
			},
		})
		if err != nil {
			t.Errorf("Submit B: %v", err)
			close(bQueued)
			return
		}
		// B cannot take the path: it is visible as Queued with the reason.
		waitFor(t, "B 的排队行与原因", func() bool {
			row, ok := s.Row(idB)
			return ok && row.State == statemachine.StateQueued &&
				strings.Contains(row.Note, "在等 "+aID+" 释放")
		})
		close(bQueued)
	}()
	<-bQueued

	// The state trace the ticket names: #20 (conflict → Queued) is the row
	// state above; release drives #39 (Queued → Acting/Thinking here).
	close(releaseHolder)
	_, errA := s.Wait(idA)
	if errA != nil {
		t.Fatalf("A: %v", errA)
	}
	resB, errB := s.Wait(idB)
	if errB != nil {
		t.Fatalf("B: %v", errB)
	}
	if resB.Status != "completed" {
		t.Errorf("B 释放后落到 %s, want completed", resB.Status)
	}
	if owner, ok := locks.HolderOf(p); ok {
		t.Errorf("收尾后路径 %s 仍被 %s 占着, want 已释放", p, owner)
	}
}

// Test47AC2SameTaskReacquireIsNotAConflict pins the idempotent half: a task
// re-touching its own file across steps does not queue behind itself.
func Test47AC2SameTaskReacquireIsNotAConflict(t *testing.T) {
	locks := NewPathLock(PathLockOptions{})
	s := New(Options{Registry: observe.NewRegistry(), Locks: locks, MaxConcurrent: 1})
	const p = `C:/wisp47/复用.txt`
	id, err := s.Submit(TaskSpec{
		Label: "自复用",
		Exec: func(ctx context.Context, run *TaskRun) TaskResult {
			for i := 0; i < 3; i++ {
				if err := run.Locks.Acquire(ctx, run.ID, []string{p}); err != nil {
					return TaskResult{Status: "failed", Err: err.Error()}
				}
			}
			run.Locks.ReleaseAll(run.ID)
			return TaskResult{Status: "completed"}
		},
	})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	res, err := s.Wait(id)
	if err != nil {
		t.Fatalf("Wait: %v", err)
	}
	if res.Status != "completed" {
		t.Errorf("自复用三次落到 %+v, want completed", res)
	}
}

// ---------------------------------------------------------------------------
// AC#3 - deadlock fixture: timeout + explicit error, no silent hang
// ---------------------------------------------------------------------------

// Test47AC3MutualWaitTimesOutWithExplicitError is the ticket's exact fixture:
// A holds p1 wants p2; B holds p2 wants p1. The guard gives up with
// ErrDeadlock instead of hanging.
func Test47AC3MutualWaitTimesOutWithExplicitError(t *testing.T) {
	reg := observe.NewRegistry()
	locks := NewPathLock(PathLockOptions{})
	s := New(Options{Registry: reg, Locks: locks, MaxConcurrent: 2})

	p1, p2 := `C:/wisp47/甲.txt`, `C:/wisp47/乙.txt`
	shrink := func(id string) { locks.SetWaitCeiling(id, 150*time.Millisecond) }

	aReady := make(chan struct{}) // A holds p1
	bReady := make(chan struct{}) // B holds p2
	goOn := make(chan struct{})   // both hold: now the mutual grab begins

	idA, _ := s.Submit(TaskSpec{
		Label: "A 锁甲要乙",
		Exec: func(ctx context.Context, run *TaskRun) TaskResult {
			shrink(run.ID)
			if err := run.Locks.Acquire(ctx, run.ID, []string{p1}); err != nil {
				return TaskResult{Status: "failed", Err: err.Error()}
			}
			aReady <- struct{}{}
			<-goOn
			// The conflict: A wants what B holds.
			err := run.Locks.Acquire(ctx, run.ID, []string{p2})
			run.Locks.ReleaseAll(run.ID)
			if err == nil {
				return TaskResult{Status: "completed", Text: "A 拿到了乙（不该发生）"}
			}
			if !errors.Is(err, ErrDeadlock) {
				return TaskResult{Status: "failed", Err: "A 的等待错了类：" + err.Error()}
			}
			return TaskResult{Status: "failed", Err: err.Error()}
		},
	})
	idB, _ := s.Submit(TaskSpec{
		Label: "B 锁乙要甲",
		Exec: func(ctx context.Context, run *TaskRun) TaskResult {
			shrink(run.ID)
			if err := run.Locks.Acquire(ctx, run.ID, []string{p2}); err != nil {
				return TaskResult{Status: "failed", Err: err.Error()}
			}
			bReady <- struct{}{}
			<-goOn
			err := run.Locks.Acquire(ctx, run.ID, []string{p1})
			run.Locks.ReleaseAll(run.ID)
			if err == nil {
				return TaskResult{Status: "completed", Text: "B 拿到了甲（不该发生）"}
			}
			if !errors.Is(err, ErrDeadlock) {
				return TaskResult{Status: "failed", Err: "B 的等待错了类：" + err.Error()}
			}
			return TaskResult{Status: "failed", Err: err.Error()}
		},
	})

	<-aReady
	<-bReady
	close(goOn)

	resA, errA := s.Wait(idA)
	resB, errB := s.Wait(idB)
	for name, res := range map[string]TaskResult{"A": resA, "B": resB} {
		if !strings.Contains(res.Err, ErrDeadlock.Error()) {
			t.Errorf("%s 的收尾没有死锁措辞：%+v", name, res)
		}
	}
	if errA != nil || errB != nil {
		t.Fatalf("Wait 错误：A=%v B=%v", errA, errB)
	}
}

// Test47AC3CancelUnblocksAPathWaiter pins the stoppable half: a task parked
// on Acquire leaves the wait when its root ctx is cancelled - the wait is
// stoppable, not a timer that owns the task.
func Test47AC3CancelUnblocksAPathWaiter(t *testing.T) {
	reg := observe.NewRegistry()
	locks := NewPathLock(PathLockOptions{})
	s := New(Options{Registry: reg, Locks: locks, MaxConcurrent: 2})

	const p = `C:/wisp47/被停.txt`
	holderDone := make(chan struct{})
	aHolds := make(chan struct{}) // A holds p before B starts acquiring
	idA, _ := s.Submit(TaskSpec{
		Label: "持有人",
		Exec: func(ctx context.Context, run *TaskRun) TaskResult {
			if err := run.Locks.Acquire(ctx, run.ID, []string{p}); err != nil {
				return TaskResult{Status: "failed", Err: err.Error()}
			}
			aHolds <- struct{}{}
			<-holderDone
			run.Locks.ReleaseAll(run.ID)
			return TaskResult{Status: "completed"}
		},
	})
	waiterStarted := make(chan struct{}, 1)
	idB, err := s.Submit(TaskSpec{
		Label: "等待者",
		Exec: func(ctx context.Context, run *TaskRun) TaskResult {
			<-aHolds // deterministic: B only ever SEES the contended wait
			acqErr := run.Locks.Acquire(ctx, run.ID, []string{p})
			waiterStarted <- struct{}{} // buffered; unblocks even on cancel
			if acqErr == nil {
				return TaskResult{Status: "failed", Err: "B 拿到了 A 占着的路径（不该发生）"}
			}
			return TaskResult{Status: "failed", Err: acqErr.Error()}
		},
	})
	if err != nil {
		t.Fatalf("Submit B: %v", err)
	}
	waiterDone := make(chan TaskResult, 1)
	go func() {
		res, _ := s.Wait(idB)
		waiterDone <- res
	}()
	// B's entry into the contended wait is directly observable through the
	// lock (before it starts: ok=false; once parked: true - and it cannot
	// finish the Acquire while A holds, so a true is the parked reading).
	waitFor(t, "等待者进入等路", func() bool {
		_, _, ok := locks.Waiting(idB)
		return ok
	})
	if err := s.Cancel(idB); err != nil {
		t.Fatalf("Cancel waiter: %v", err)
	}
	select {
	case res := <-waiterDone:
		if !strings.Contains(res.Err, ErrWaitCancelled.Error()) {
			t.Errorf("取消后等待者的错误形状 = %+v, want ErrWaitCancelled 族", res)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("取消等待者后它没有解除阻塞：等待不可停（D32 违）")
	}
	close(holderDone)
	if _, err := s.Wait(idA); err != nil {
		t.Fatalf("Wait A: %v", err)
	}
}

// ---------------------------------------------------------------------------
// AC#4 - approval independence
// ---------------------------------------------------------------------------

// Test47AC4ApprovalWaitDoesNotBlockOtherTasks is AC#4's core D31 property:
// task A parked on an approval-waiter while task B executes tooling to
// completion.
func Test47AC4ApprovalWaitDoesNotBlockOtherTasks(t *testing.T) {
	reg := observe.NewRegistry()
	s := New(Options{Registry: reg, MaxConcurrent: 2})

	approve := make(chan struct{})
	var idA string
	aReady := make(chan struct{})
	go func() {
		id, err := s.Submit(TaskSpec{
			Label: "A 等批准",
			Exec: func(ctx context.Context, run *TaskRun) TaskResult {
				// The host's per-task approval-waiter shape (D38b): a named
				// goroutine parked until the card is answered.
				run.Reg.Spawn("approval-waiter", "scheduler-test", run.Root, func(c context.Context) {
					close(aReady)
					<-approve
					_ = c
				})
				return TaskResult{Status: "completed", Text: "批准后收尾"}
			},
		})
		if err != nil {
			t.Errorf("Submit A: %v", err)
			return
		}
		idA = id
	}()
	<-aReady

	idB, err := s.Submit(TaskSpec{
		Label: "B 照常干活",
		Exec: func(ctx context.Context, run *TaskRun) TaskResult {
			done := run.Reg.Spawn("tool-exec-"+run.ID, "scheduler-test", run.Root, func(context.Context) {
			})
			<-done.Done()
			return TaskResult{Status: "completed", Text: "B 的工具跑完了"}
		},
	})
	if err != nil {
		t.Fatalf("Submit B: %v", err)
	}

	// B completes while A's card is still open.
	resB, err := s.Wait(idB)
	if err != nil {
		t.Fatalf("Wait B: %v", err)
	}
	if resB.Status != "completed" {
		t.Errorf("B 在 A 等批时落到 %+v, want completed（D31 核心性质被破）", resB)
	}
	// Then A's approval lands and A completes too.
	close(approve)
	resA, err := s.Wait(idA)
	if err != nil {
		t.Fatalf("Wait A: %v", err)
	}
	if resA.Status != "completed" {
		t.Errorf("A 批准后落到 %+v, want completed", resA)
	}
	waitFor(t, "A 的 waiter 退场", func() bool { return reg.CountByName("approval-waiter") == 0 })
}

// ---------------------------------------------------------------------------
// AC#5 - cancel under concurrency
// ---------------------------------------------------------------------------

// Test47AC5MidFlightCancelCarriesAppliedStepsSiblingsUnaffected is AC#5: a
// mid-flight cancel settles the cancelled task with its applied-steps
// report; the siblings finish unaffected.
func Test47AC5MidFlightCancelCarriesAppliedStepsSiblingsUnaffected(t *testing.T) {
	reg := observe.NewRegistry()
	s := New(Options{Registry: reg, MaxConcurrent: 2})

	cancelMe := make(chan struct{})
	aStarted := make(chan struct{})
	idA, _ := s.Submit(TaskSpec{
		Label: "A 被停",
		Exec: func(ctx context.Context, run *TaskRun) TaskResult {
			close(aStarted)
			<-cancelMe
			if ctx.Err() != nil {
				// The D31 shape: report what already landed, never a clean undo.
				return TaskResult{
					Status:       "cancelled",
					AppliedSteps: []string{"已写入 C:/wisp47/半成品.txt"},
				}
			}
			return TaskResult{Status: "completed"}
		},
	})
	<-aStarted

	var idB string
	bDone := make(chan TaskResult, 1)
	go func() {
		id, err := s.Submit(TaskSpec{
			Label: "B 兄弟",
			Exec: func(ctx context.Context, run *TaskRun) TaskResult {
				done := run.Reg.Spawn("tool-exec-"+run.ID, "scheduler-test", run.Root, func(context.Context) {})
				<-done.Done()
				return TaskResult{Status: "completed", Text: "B 完事"}
			},
		})
		if err != nil {
			bDone <- TaskResult{Status: "failed", Err: err.Error()}
			return
		}
		idB = id
		res, _ := s.Wait(id)
		bDone <- res
	}()
	// B must reach completion even though A is still running.
	var resB TaskResult
	select {
	case resB = <-bDone:
	case <-time.After(5 * time.Second):
		t.Fatal("B 没能独立完成（A 占着并发窗不是理由）")
	}
	if resB.Status != "completed" {
		t.Errorf("B 落到 %+v, want completed", resB)
	}

	// Cancel A mid-flight; the row settles with the applied steps verbatim.
	if err := s.Cancel(idA); err != nil {
		t.Fatalf("Cancel A: %v", err)
	}
	close(cancelMe)
	resA, err := s.Wait(idA)
	if err != nil {
		t.Fatalf("Wait A: %v", err)
	}
	if resA.Status != "cancelled" {
		t.Errorf("A 落到 %s, want cancelled", resA.Status)
	}
	row, ok := s.Row(idA)
	if !ok {
		t.Fatal("A 的行不见了")
	}
	if !strings.Contains(row.Note, "已写入 C:/wisp47/半成品.txt") ||
		!strings.Contains(row.Note, "未自动回退") {
		t.Errorf("A 的行注没有带 applied-steps 台账（D31 措辞）：%q", row.Note)
	}
	_ = idB
}

// Test47AC5CancelAQueuedTaskSettlesItWithoutRunning pins the queued-cancel
// half: a task cancelled before its first slot never runs its Exec.
func Test47AC5CancelAQueuedTaskSettlesItWithoutRunning(t *testing.T) {
	reg := observe.NewRegistry()
	s := New(Options{Registry: reg, MaxConcurrent: 1})

	gate := make(chan struct{})
	idA, _ := s.Submit(TaskSpec{Label: "占窗", Exec: func(context.Context, *TaskRun) TaskResult {
		<-gate
		return TaskResult{Status: "completed"}
	}})
	ran := make(chan struct{})
	idB, err := s.Submit(TaskSpec{Label: "在队", Exec: func(context.Context, *TaskRun) TaskResult {
		close(ran)
		return TaskResult{Status: "completed"}
	}})
	if err != nil {
		t.Fatalf("Submit B: %v", err)
	}
	waitFor(t, "B 进入队列", func() bool {
		row, ok := s.Row(idB)
		return ok && row.State == statemachine.StateQueued
	})
	if err := s.Cancel(idB); err != nil {
		t.Fatalf("Cancel B: %v", err)
	}
	res, err := s.Wait(idB)
	if err != nil {
		t.Fatalf("Wait B: %v", err)
	}
	if res.Status != "" {
		t.Errorf("被取消的排队任务不该有执行结果：%+v", res)
	}
	select {
	case <-ran:
		t.Fatal("排队任务被取消后仍然执行了 Exec")
	default:
	}
	close(gate)
	if _, err := s.Wait(idA); err != nil {
		t.Fatalf("Wait A: %v", err)
	}
}

// ---------------------------------------------------------------------------
// D38e - door + CancelAll
// ---------------------------------------------------------------------------

// Test47D38eDoorShutsAndCancelAllDrainsOrAbandons pins the door (step 1) and
// the step-3 shape: no new tasks after Close; CancelAll cancels roots and
// drains the survivors.
func Test47D38eDoorShutsAndCancelAllDrainsOrAbandons(t *testing.T) {
	reg := observe.NewRegistry()
	s := New(Options{Registry: reg, MaxConcurrent: 2})

	stuck := make(chan struct{})
	_, _ = s.Submit(TaskSpec{Label: "粘住", Exec: func(ctx context.Context, _ *TaskRun) TaskResult {
		<-stuck
		// Honour the cancellation: a cancelled root unblocks the body.
		return TaskResult{Status: "cancelled"}
	}})
	// The stuck body above only observes the local channel, so the drain
	// wait below times out against it on purpose: CancelAll must answer the
	// abandoned list (D37c) instead of hanging.
	_, err := s.Submit(TaskSpec{Label: "陪跑", Exec: tc("收工")})
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	s.Close()
	if _, err := s.Submit(TaskSpec{Label: "迟到", Exec: tc("不该跑")}); !errors.Is(err, ErrClosed) {
		t.Errorf("关门后 Submit 回了 %v, want ErrClosed", err)
	}
	abandoned, err := s.CancelAll(context.Background(), 150*time.Millisecond)
	if err != nil {
		t.Fatalf("CancelAll: %v", err)
	}
	if len(abandoned) != 1 {
		t.Errorf("CancelAll 报弃 %v, want 恰 1（只有粘住那枚不响应）", abandoned)
	}
	close(stuck)
	waitFor(t, "名册排干", func() bool { return reg.Count() == 0 })
}
