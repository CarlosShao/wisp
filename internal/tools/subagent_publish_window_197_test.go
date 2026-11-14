package tools

// Ticket 70 step 7 read a real defect out of CI and ticket 197 owns the code it
// lives in: the spawn path used to publish a subagent's roster row and then,
// one statement later, attach its stop handle. Between those two statements the
// row is observable ("state=Thinking", i.e. this child is in flight) while
// TaskRoster.Cancel finds no handle and answers 「没有在跑 … 没有可停的句柄」 -
// two facts in one sentence that cannot both be true.
//
// The two legs below are the two halves of the instrument, and they answer
// different questions on purpose:
//
//	Test197InFlightRowNeverLacksItsStopHandle   the SYMPTOM, through the real
//	    bridge / real roster / real spawn path. It is the shape CI failed on.
//	    Pre-fix it is scheduler-bound (the window closes when the spawner's
//	    goroutine executes the statement after the publish), so it is counted
//	    over a denominator instead of being pretended to be deterministic:
//	    publishWindowRounds rounds, hits reported verbatim.
//	    Post-fix it is deterministic: admission publishes row+handle in ONE
//	    critical section, admission happens-before the child's first model call,
//	    and a round's release point IS that first model call - so the handle is
//	    already there, by happens-before, every time.
//
//	Test197SpawnPathPublishesRowAndStopHandleTogether   the SHAPE, deterministic.
//	    The one deterministic red available for a two-statement window is a
//	    reading of the two statements themselves: the spawn path may file the row
//	    only through the publisher that takes the cancel handle with it, and may
//	    not pair a bare publish with a later AttachCancel.
//
// Neither leg adds a clock: no sleep, no deadline, no elapsed-time threshold
// (D22 bans wall-clock timeouts, and a window does not become measured by
// waiting longer).

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// publishWindowRounds is the denominator the ticket asked for: 40 rounds. It is
// a sample size, not a threshold - a single hit is a red, and the failure text
// carries the hit count so the reading cannot be mistaken for a pass rate.
const publishWindowRounds = 40

// publishWindowPerRound is how many spawners run at once in one round. Four is
// the pool's own ceiling (MaxConcurrentSubagents), so a round is exactly "the
// whole pool spawning at the same instant" - the contention the CI reading was
// taken under (triage step 7 ④: 负载一高，父 goroutine 在 RunAsync 返回后被抢占一下).
// It is also four independent windows per round instead of one.
const publishWindowPerRound = 4

// Test197InFlightRowNeverLacksItsStopHandle holds a full pool of children inside
// their first model call and asks the product's own task.cancel the only
// question the defect makes answerable: the row says this child is in flight -
// can its parent stop it?
func Test197InFlightRowNeverLacksItsStopHandle(t *testing.T) {
	hits := 0
	reds := make([]string, 0, 4)
	for round := 0; round < publishWindowRounds; round++ {
		started := make(chan struct{}, publishWindowPerRound)
		release := make(chan struct{})
		x := build221(t, windowGate221(), heldProvider(started, release), false)

		type held struct {
			label string
			done  chan Result
		}
		kids := make([]held, 0, publishWindowPerRound)
		for i := 0; i < publishWindowPerRound; i++ {
			label := fmt.Sprintf("窗口-%d-%d", round, i)
			done := make(chan Result, 1)
			go func(label string) {
				res, err := x.h.spawnResult(t.Context(), label, "长任务")
				if err != nil {
					res = Result{Text: "dispatch error: " + err.Error(), IsError: true}
				}
				done <- res
			}(label)
			kids = append(kids, held{label: label, done: done})
		}
		// One token per in-flight model call. When this loop is done, every
		// child of this round is past everything that has to happen "before the
		// child does any work" - which is exactly the instant CI's failing
		// assertion read.
		for i := 0; i < publishWindowPerRound; i++ {
			select {
			case <-started:
			case <-t.Context().Done():
				t.Fatalf("第 %d 发：有子代理没有开始它的模型调用", round)
			}
		}
		byLabel := map[string]string{}
		for _, r := range x.h.childRows() {
			byLabel[r.Out.Label] = r.TaskID
		}
		for _, kid := range kids {
			childID, ok := byLabel[kid.label]
			if !ok {
				t.Fatalf("第 %d 发：名册里没有 %q 那一行：%+v", round, kid.label, x.h.childRows())
			}
			row := x.h.mustLook(t, childID)
			if got := row.Out.State; got != statemachine.StateThinking {
				t.Fatalf("第 %d 发：%s 那一行 state = %q, want %q —— 这台仪器只在「在飞」那一行上读数",
					round, kid.label, string(got), string(statemachine.StateThinking))
			}
			if res := x.cancel(t, parent197, childID); res.IsError {
				hits++
				if len(reds) < 4 {
					reds = append(reds, fmt.Sprintf("第 %d 发 %s（%s）：%s", round, kid.label, childID, res.Text))
				}
			}
		}

		// Let the whole pool join before the next round. This is a join, not a
		// wait for a clock.
		close(release)
		for _, kid := range kids {
			select {
			case <-kid.done:
			case <-t.Context().Done():
				t.Fatalf("第 %d 发：被停/被放行的子代理没有让 task.spawn 收口", round)
			}
		}
	}
	if hits != 0 {
		t.Errorf("行已发布、句柄未挂：%d/%d 个在飞子代理被父任务停时回了「没有在跑」"+
			"（CI 20 发里 2 发命中 = 10%%，本地空载一时打不开；这就是那枚窗口）：\n%s",
			hits, publishWindowRounds*publishWindowPerRound, strings.Join(reds, "\n"))
	}
}

// Test197SpawnPathPublishesRowAndStopHandleTogether is the deterministic half:
// it reads the spawn path's own two statements. A row published by a call that
// does not carry the handle, followed by a separate AttachCancel, is the window
// no matter how the scheduler happens to interleave it - so the shape itself is
// what this leg refuses.
func Test197SpawnPathPublishesRowAndStopHandleTogether(t *testing.T) {
	src, err := os.ReadFile("subagent_197.go")
	if err != nil {
		t.Fatalf("读不到 subagent_197.go（这把尺的落点是包内源文件）：%v", err)
	}
	lines := strings.Split(string(src), "\n")
	if !strings.Contains(string(src), "PublishSubagentInFlight(") {
		t.Errorf("派生路径没有调用 PublishSubagentInFlight：行与句柄必须由一次发布落地，" +
			"否则「行可见 ⇒ 句柄已挂」就只是一个约定而不是不变量")
	}
	for i, line := range lines {
		code := strings.TrimSpace(line)
		if strings.HasPrefix(code, "//") {
			continue
		}
		if strings.Contains(code, ".PublishSubagent(") {
			t.Errorf("subagent_197.go:%d 仍在用裸 PublishSubagent 发行（这一调用不带句柄）：%s\n"+
				"这一行就是票 70 step 7 读出来的窗口：行先可见、句柄后挂，中间那一瞬"+
				"state=Thinking 而 Cancel 说「没有在跑」", i+1, code)
		}
		if strings.Contains(code, ".AttachCancel(") {
			t.Errorf("subagent_197.go:%d 仍把 AttachCancel 当第二次发布：%s\n"+
				"句柄必须和行进同一个临界区（PublishSubagentInFlight），"+
				"否则两次发布之间就是一个可观察的中间态", i+1, code)
		}
	}
}
