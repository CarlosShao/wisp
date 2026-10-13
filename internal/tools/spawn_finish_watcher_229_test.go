package tools

// Ticket 229, AC#2/#3/#5: the subagent finish watcher (the goroutine that
// joins a detached child and files its roster row) used to spawn under
// "subagent-finish-<id>", a name in none of the D38b tables - so every spawn
// rang the leak alarm ("goroutine outside the D38 roster (leak symptom)") and
// the one signal that is supposed to mean "a goroutine ran away" drowned in
// routine noise.
//
// This file pins the 甲形 fix from the outside: while a really spawned child is
// held mid-model-call, observe.Default must show the watcher under a per-task
// roster name (agent-task-<childID>) and RosterReport.Unknown must be free of
// the old name. The positive control below keeps the door shut the other way:
// a genuinely unregistered name must still be classified CategoryUnknown AND
// still make the WARN fire - the fix may not, and does not, weaken the alarm.

import (
	"context"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/llm"
	"github.com/CarlosShao/wisp/internal/observe"
)

// Test229FinishWatcherIsRosterNamedWhileChildRuns is the AC#2 judgement on the
// production path: harness -> real task.spawn -> observe.Default. The child's
// first model call is held open by the C5 fake (the frozen seam), so the
// watcher is deterministically alive while the registry is read.
func Test229FinishWatcherIsRosterNamedWhileChildRuns(t *testing.T) {
	h := newSub197Harness(t)
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	prov := &fake197Provider{
		name: "child-229", answer: "孩子把日志读完了 229", roster: h.roster,
		started: started, release: release,
	}
	h.build(t, func() llm.LlmProvider { return prov }, true)

	var releaseOnce sync.Once
	releaseChild := func() { releaseOnce.Do(func() { close(release) }) }
	t.Cleanup(releaseChild)

	type spawnOut struct {
		res Result
		err error
	}
	outc := make(chan spawnOut, 1)
	go func() {
		res, err := h.spawnResult(t.Context(), label197, "把日志读完并总结")
		outc <- spawnOut{res: res, err: err}
	}()

	select {
	case <-started:
	case <-time.After(10 * time.Second):
		t.Fatal("孩子的第一发模型调用没起来（孩子在按住的缝上），这一发的读数无从谈起")
	}
	id := h.childRow(t).TaskID
	want := "agent-task-" + id     // 甲形：每任务前缀族里属于孩子任务的那一行
	old := "subagent-finish-" + id // 修前的名字，不在名册任何一表

	deadline := time.Now().Add(5 * time.Second)
	for observe.Default.CountByName(want) == 0 {
		if time.Now().After(deadline) {
			rep := observe.Default.RosterReport()
			t.Fatalf("收尾协程不在册：CountByName(%q)=0；此刻 observe.Default：Total=%d PerTask=%d Unknown=%v；"+
				"旧名 %q 的 CountByName=%d（AC#2 没落地，或名册判定被改）",
				want, rep.Total, rep.PerTask, rep.Unknown, old, observe.Default.CountByName(old))
		}
		time.Sleep(5 * time.Millisecond)
	}

	rep := observe.Default.RosterReport()
	t.Logf("AC#5 安静期读数（孩子被按住、父任务还停在等待上）：observe.Default.Count()=%d —— 与 SLO 同一口径"+
		"（sampler.go:365 读 s.reg.Count()，生产里 reg=observe.Default）；RosterReport：Resident=%d OnDemand=%d "+
		"PerTask=%d Temporary=%d Unknown=%v；CountByName(%q)=%d；CountByName(%q)=%d",
		observe.Default.Count(), rep.Resident, rep.OnDemand, rep.PerTask, rep.Temporary, rep.Unknown,
		want, observe.Default.CountByName(want), old, observe.Default.CountByName(old))
	if cat := observe.ClassifyGoroutine(want); cat != observe.CategoryPerTask {
		t.Errorf("ClassifyGoroutine(%q) = %q, want %q", want, cat, observe.CategoryPerTask)
	}
	for _, u := range rep.Unknown {
		if strings.HasPrefix(u, "subagent-finish-") {
			t.Errorf("收尾协程仍以名册外名字躺在 RosterReport.Unknown 里：%q", u)
		}
	}

	// 放手：孩子正常收尾，父侧的 task.spawn 拿回结论，名册里不留尸位。
	releaseChild()
	select {
	case o := <-outc:
		if o.err != nil {
			t.Fatalf("task.spawn 返回错误：%v", o.err)
		}
		if o.res.IsError {
			t.Fatalf("task.spawn 答了错误：%s", o.res.Text)
		}
		if !strings.Contains(o.res.Text, "孩子把日志读完了 229") {
			t.Errorf("孩子的结论没回到父任务：%q", o.res.Text)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("孩子收尾后 task.spawn 没有在等待窗口内返回")
	}

	deadline = time.Now().Add(15 * time.Second)
	for observe.Default.CountByName(want) != 0 {
		if time.Now().After(deadline) {
			t.Fatalf("孩子收尾后收尾协程仍在册：CountByName(%q)=%d（留了尸位）", want, observe.Default.CountByName(want))
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// Test229UnknownNameStillRingsLeakAlarm is the AC#3 positive control: plant a
// name that is in no table and assert the alarm still rings, both as a live
// WARN record and as CategoryUnknown. The fix must only change WHICH names are
// accounted for, never how loudly an unaccounted one screams.
func Test229UnknownNameStillRingsLeakAlarm(t *testing.T) {
	const fake = "229-probe-not-in-any-roster-table"

	prev := slog.Default()
	rec := &slogRec229{}
	slog.SetDefault(slog.New(slog.NewTextHandler(rec, &slog.HandlerOptions{Level: slog.LevelWarn})))
	t.Cleanup(func() { slog.SetDefault(prev) })

	reg := observe.NewRegistry()
	release := make(chan struct{})
	h := reg.Spawn(fake, "test-229", nil, func(context.Context) { <-release })

	if rep := reg.RosterReport(); len(rep.Unknown) != 1 || rep.Unknown[0] != fake {
		t.Fatalf("RosterReport.Unknown = %v, want [%q]：分类门被拆了？", rep.Unknown, fake)
	}
	if cat := observe.ClassifyGoroutine(fake); cat != observe.CategoryUnknown {
		t.Fatalf("ClassifyGoroutine(%q) = %q, want %q：分类门被拆了？", fake, cat, observe.CategoryUnknown)
	}
	logged := rec.String()
	if !strings.Contains(logged, "goroutine outside the D38 roster (leak symptom)") {
		t.Fatalf("名册外的名字没让漏哨响，日志=%q：WARN 被调轻或删掉了？", logged)
	}
	if !strings.Contains(logged, fake) {
		t.Fatalf("漏哨响了但没带上那枚名字，日志=%q", logged)
	}

	close(release)
	<-h.Done()
}

// slogRec229 is a concurrency-safe sink for the captured WARN lines.
type slogRec229 struct {
	mu  sync.Mutex
	buf strings.Builder
}

func (w *slogRec229) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.Write(p)
}

func (w *slogRec229) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.String()
}
