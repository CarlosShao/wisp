package tools

// Ticket 176-r1 - the standing criteria for the WRITE side of ticket 164's
// task roster: the host now starts a background task through the loop's
// existing async entry point and files its answer afterwards, so task.output
// has an object to read on a real run.
//
// The four legs, and the mutation each one turns red on (readings recorded in
// docs/evidence/s1/176-background-start-port-r1.md):
//
//	L-1 TestBackfillFilesTheRealSpilledArtifact176r1
//	    -> drop the ArtifactPath assignment (record only the text): a long
//	       answer then gets no pointer, which is PLAN.md:2564's 「只截不指」.
//	L-2 TestTaskArtifactKeyStaysOutOfTheCallIDNamespace176r1
//	    -> hand Prepare the BARE task id as the artifact key (the line
//	       taskArtifactPrefix + taskID -> taskID): a model-supplied call id
//	       that repeats the id now lands on the same file name, and
//	       spill.go:117-131 resolves that as last-writer-wins - a silent byte
//	       swap under a pointer the host itself printed.
//	L-3 TestNoRosterWriteAfterCancel176r1
//	    -> remove the Stopped(ctx) guard in Backfill, or lift Record above the
//	       join: the cancelled task's output then lands in the table.
//	       This leg carries its own control half (no cancel -> the same port
//	       DOES write), so it cannot be bought by "never write at all".
//	L-4 TestTerminalWriteCarrierStaysOutOfScope176r1
//	    -> widen the guard to "refuse whenever the PARENT is done": the loop's
//	       own terminal write carrier (context.WithoutCancel,
//	       internal/agent/loop.go:956-958 - the one legitimate still-writing
//	       shape after a cancel) would then be read as a violation, and a
//	       green here would be the host flagging its own teardown.
//
// What these legs do NOT pin: that the composition root actually calls
// RunAsync. That claim is a grep, not a test in this package (cmd/wisp's own
// tests are not this ticket's write surface), and its reading is in the
// evidence file next to the start-port line in cmd/wisp/run.go.

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/observe"
)

// joinBudget176r1 is the join budget for this leg's own harness; it touches no
// product constant and no threshold.
const joinBudget176r1 = 2 * time.Second

// backfill176r1 builds the port's real dependency set: a fresh roster and the
// loop's own D15(3) writer (agent.Spiller, reference budgets) over a temporary
// artifacts directory. No stand-in for the writer - the whole point of these
// legs is that the bytes land on disk.
func backfill176r1(t *testing.T) (*TaskRoster, *agent.Spiller, string) {
	t.Helper()
	dir := t.TempDir()
	roster := NewTaskRoster()
	spills := agent.NewSpiller(filepath.Join(dir, "artifacts"), agent.BudgetsFor(0))
	return roster, spills, dir
}

// ---------------------------------------------------------------------------
// L-1 the roster gets the text AND the path of the copy that really landed.
// ---------------------------------------------------------------------------

func TestBackfillFilesTheRealSpilledArtifact176r1(t *testing.T) {
	roster, spills, _ := backfill176r1(t)
	port := TaskBackfill{Roster: roster, Spills: spills}
	ctx := context.Background()

	// (a) an oversized answer: text + the artifact that holds those same bytes.
	const longID = "0b7d1f3e9c2a4571"
	long := asciiRun(20000) // > the 4000-token / 16000-byte spill threshold
	rec, why := port.Backfill(ctx, longID, long)
	if why != "" {
		t.Fatalf("回填报告了理由，记录却没进名册：%s", why)
	}
	got, ok := roster.Look(longID)
	if !ok {
		t.Fatalf("名册里没有 %s：Count()=%d，task.output 这一发还是会答「查不到这个任务」",
			longID, roster.Count())
	}
	if got.Text != long {
		t.Errorf("名册里的正文与任务打印的不一致：got %d 字节, want %d", len(got.Text), len(long))
	}
	if rec.ArtifactPath != got.ArtifactPath {
		t.Errorf("返回值与记录不同一份路径：rec=%q got=%q", rec.ArtifactPath, got.ArtifactPath)
	}
	if got.ArtifactPath == "" {
		t.Fatal("超长输出没有产物路径：这就是 PLAN.md:2564 的「只截不指」")
	}
	body, err := os.ReadFile(got.ArtifactPath)
	if err != nil {
		t.Fatalf("登记的路径读不回来（票 174 的「说实话」规矩）：%v", err)
	}
	if string(body) != long {
		t.Errorf("产物里是别人的字节：%d 字节 vs 原文 %d 字节", len(body), len(long))
	}

	// (b) a short answer: recorded, and it promises NO path at all.
	const shortID = "5a19c2e0d6b3487f"
	if _, why := port.Backfill(ctx, shortID, "就一句话"); why != "" {
		t.Fatalf("短文回填失败：%s", why)
	}
	short, ok := roster.Look(shortID)
	if !ok {
		t.Fatalf("短文没进名册（Count()=%d）", roster.Count())
	}
	if short.ArtifactPath != "" {
		t.Errorf("没落盘却许了路径：%q（读不回来的许诺就是假指针）", short.ArtifactPath)
	}
	if short.Text != "就一句话" {
		t.Errorf("短文正文 = %q, want 逐字", short.Text)
	}
	if n := roster.Count(); n != 2 {
		t.Errorf("Count() = %d, want 2", n)
	}
}

// ---------------------------------------------------------------------------
// L-2 namespace discipline: a background record's artifact key must not be
// reachable by a model-supplied tool-call id.
// ---------------------------------------------------------------------------

func TestTaskArtifactKeyStaysOutOfTheCallIDNamespace176r1(t *testing.T) {
	roster, spills, _ := backfill176r1(t)

	const taskID = "e3f1c0a97b2d4865"
	mine := asciiRun(20000)
	if _, why := (TaskBackfill{Roster: roster, Spills: spills}).
		Backfill(context.Background(), taskID, mine); why != "" {
		t.Fatalf("回填没有成功：%s", why)
	}
	pointer, _ := roster.Look(taskID)
	if pointer.ArtifactPath == "" {
		t.Fatal("precondition: 这条记录没有产物路径，命名空间的对照就散了")
	}
	// The precedent name (loop.go:324 registers the same task as
	// "agent-task-<id>" in the D38 roster) must be readable off the file name,
	// because that prefix is the only thing separating the two key spaces.
	if name := filepath.Base(pointer.ArtifactPath); !strings.Contains(name, "agent-task-"+taskID) {
		t.Errorf("产物名 %q 不带 agent-task- 前缀：后台记录与模型 supplied 的 call id 又回到同一命名空间", name)
	}

	// The collision attempt, played the way the model can actually play it: a
	// tool-call id that repeats a task id it has seen. The loop's own spill
	// layer addresses artifacts by call id, so this is the same writer, the same
	// directory, the same key space.
	foreign := asciiRun(19000)
	sp, err := spills.Prepare(taskID, foreign)
	if err != nil {
		t.Fatalf("Prepare(模型 call id): %v", err)
	}
	if sp.Path == pointer.ArtifactPath {
		t.Fatalf("模型 supplied 的 call id 与任务 id 落在同一个产物名上：%q", sp.Path)
	}
	body, err := os.ReadFile(pointer.ArtifactPath)
	if err != nil {
		t.Fatalf("任务那份产物读不回来了：%v", err)
	}
	if string(body) != mine {
		t.Errorf("任务产物被静默换成了别的内容：%d 字节 vs 原 %d 字节（spill.go:117-131 的 last-writer-wins 不会报错）",
			len(body), len(mine))
	}
}

// ---------------------------------------------------------------------------
// L-3 AC#3: once a task is cancelled, nobody may still be writing its output.
// ---------------------------------------------------------------------------

// runLikeThePort176r1 reproduces the composition root's shape with the real
// registry and the real root: spawn the body under the D38 name the loop uses,
// let the caller decide when to cancel, join the handle, and only THEN reach
// Backfill with the task's own root context - which is exactly the order
// cmd/wisp/run.go uses (RunAsync -> Wait -> Backfill(bg.Root().Ctx, ...)).
// release is how the test lets an in-flight body finish after the cancel, so
// "cancel arrives while the task is still working" is a real interleaving and
// not a ctx that was already dead before the body started.
func runLikeThePort176r1(t *testing.T, port TaskBackfill, taskID, text string, cancelFirst bool) (string, int, bool) {
	t.Helper()
	reg := observe.NewRegistry()
	root := observe.NewRoot("176r1-leg")
	started := make(chan struct{})
	release := make(chan struct{})
	h := reg.Spawn("agent-task-"+taskID, "test", root, func(c context.Context) {
		close(started)
		select {
		case <-release:
		case <-c.Done():
			<-release // the body keeps working past the veto, like a real tool call
		}
	})
	<-started

	if cancelFirst {
		root.Cancel()
	}
	// The caller joins before it writes - Pending()==0 is D38e's "the task's
	// work is fully drained", and it is what makes the order observable.
	close(release)
	<-h.Done()
	if left := root.Wait(joinBudget176r1); left != 0 {
		t.Fatalf("join 没归零：Pending()=%d", left)
	}

	if port.Roster.Count() != 0 {
		t.Fatalf("precondition: 名册起手不空（%d 条），Count() 的读数就没有意义", port.Roster.Count())
	}
	_, why := port.Backfill(root.Ctx, taskID, text)
	_, present := port.Roster.Look(taskID)
	return why, port.Roster.Count(), present
}

func TestNoRosterWriteAfterCancel176r1(t *testing.T) {
	// (a) the cancelled direction: the task is cancelled mid-flight, the body
	// still finishes, the join drains - and the table is untouched.
	roster, spills, _ := backfill176r1(t)
	port := TaskBackfill{Roster: roster, Spills: spills}
	why, count, present := runLikeThePort176r1(t, port, "c1f0a4d27b6e4803", asciiRun(20000), true)
	if present {
		t.Errorf("任务被取消之后它的输出还是被写进名册了：Look() 命中（why=%q）", why)
	}
	if count != 0 {
		t.Errorf("取消之后 Count() = %d, want 0（取消不是原子的，D31：收尾不许再动名册）", count)
	}
	if why == "" {
		t.Error("跳过回填却没交代理由：「查不到这个任务」必须能被区分成「没这条记录」")
	}

	// (b) the control half, same harness, same port: without a cancel the very
	// same call must write, or this leg would be passing on silence.
	roster2, spills2, _ := backfill176r1(t)
	port2 := TaskBackfill{Roster: roster2, Spills: spills2}
	why2, count2, present2 := runLikeThePort176r1(t, port2, "9d3b7e10c58a4f26", "没被取消的正文", false)
	if !present2 || count2 != 1 {
		t.Fatalf("对照组没能写进名册（present=%v count=%d why=%q）：这一发的绿色就是假的",
			present2, count2, why2)
	}
}

// ---------------------------------------------------------------------------
// L-4 the criterion's range explicitly excludes the loop's terminal write
// carrier - the one legitimate "still writing after a cancel" in this tree.
// ---------------------------------------------------------------------------

func TestTerminalWriteCarrierStaysOutOfScope176r1(t *testing.T) {
	root := observe.NewRoot("176r1-terminal")
	root.Cancel()
	if root.Ctx.Err() == nil {
		t.Fatal("precondition: 这根 ctx 没被取消，终态写入的对照就散了")
	}
	// The exact carrier internal/agent uses for its terminal journal write
	// (loop.go:956-958: context.WithTimeout(context.WithoutCancel(ctx), ...)).
	// A criterion that read this as "still writing after cancel" would flag the
	// host's own teardown, so the guard's range is the task's ctx, not its
	// parent's history.
	terminal := context.WithoutCancel(root.Ctx)
	if stopped, why := Stopped(terminal); stopped {
		t.Errorf("环路那枚终态写入载体被判成「已停止」（%q）：取消判据的射程把自家收尾也算进去了", why)
	}

	roster, spills, _ := backfill176r1(t)
	rec, why := (TaskBackfill{Roster: roster, Spills: spills}).
		Backfill(terminal, "7e40b18d3c9a4f52", "终态窗口里仍要写下的正文")
	if why != "" {
		t.Fatalf("终态载体上回填被拒了：%s（got=%+v）", why, rec)
	}
	if _, ok := roster.Look("7e40b18d3c9a4f52"); !ok {
		t.Error("终态载体上没写进名册：这条腿的排除项没落地")
	}
}
