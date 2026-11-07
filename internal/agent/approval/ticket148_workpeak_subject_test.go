package approval_test

// Ticket 148 AC#2's subject proof: the WorkPeak posture's subject must REALLY
// run tool calls before any "cost of the extra copies" reading means anything.
//
// THE REGISTERED GAP THIS CASE CARRIES. Ticket 146 AC#3 measured the clone's
// cost with `wisp slo -state WorkPeak -seconds 5` and honestly recorded that
// the instrument cannot see the change: the WorkPeak subject is the resident
// skeleton ("sit still", cmd/wisp/slo_windows.go posture:"skeleton"), while
// LiveApprovals' only production wiring is `wisp run`'s panel pump
// (cmd/wisp/run.go:422 -> agentRuntime.liveVerdicts) - so during the sampled
// five seconds the function was never called once. The literal re-measurement
// arm ("whoever makes the WorkPeak subject really run tool calls must come
// back and re-take this cell") became ticket 148 AC#2.
//
// WHAT THIS CASE IS, AND IS NOT. This leg's write surface is
// internal/agent/approval only - cmd/wisp is off-limits, so the production
// `wisp run` posture cannot be rebuilt here. What this case proves instead is
// the precondition the ticket asks to prove FIRST: a subject that really runs
// tool calls end to end - real bridge dispatch (tools.Bridge.Execute), real
// gate assessment raising the declared-L1 tool to L2, a real pending card, a
// real answer routed through DecideFromNative - calls LiveApprovals at least
// once inside a five-second window, and the row it returns names the card the
// tool call is actually waiting on. That converts "the pump path never
// exercises this function" from an assumption into a measured fact for a
// tool-call-bearing subject, and it is the precondition the before/after
// readings in docs/evidence/s1/148-r3.md §AC#2 are taken under. It is NOT the
// production resident: the band-level `wisp slo` reading still needs the
// cmd/wisp instrument ("a run posture with a pending approval"), which stays a
// registered gap, not a claim.
//
// The answers ride the fail-closed direction (reject): the loop must never
// widen what executes, and every card it opens it settles, so the gate's
// depth-zero cleanup holds.

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// workPeakWindow148 is the five-second window the ticket names.
const workPeakWindow148 = 5 * time.Second

type toolOutcome148 struct {
	out agent.ToolOutcome
	err error
}

func Test148WorkPeakSubjectRunsToolCallsLiveApprovalsFires(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui)
	g.AdmitTextTask(testTask)
	dir := t.TempDir()
	// Declared L1, raised to L2 by R2 (path outside the allowlist): the same
	// "declared level is a lower bound, never the conclusion" shape the bridge
	// tests use, so every card here is a REAL assessed verdict, not a planted
	// queue row.
	ft := newFakeTool("fs.write", tools.CapFSWrite, risk.L1, "paths")
	b := bridgeFor(t, g, ft, []string{dir})

	deadline := time.Now().Add(workPeakWindow148)
	toolCalls, liveCalls := 0, 0
	for i := 0; time.Now().Before(deadline); i++ {
		corr := fmt.Sprintf("corr-wp-%d", i)
		args := argsJSON(t, "C:/Windows/not-mine.txt")
		res := make(chan toolOutcome148, 1)
		spawn(t, func() {
			out, err := b.Execute(context.Background(), agentToolRequest(corr, args))
			res <- toolOutcome148{out: out, err: err}
		})

		// The card is up: the tool call really reached the approval flow.
		p := ui.wait(t)
		if p.CorrelationID == "" {
			t.Fatalf("第 %d 次工具调用的卡没有 correlation id", i)
		}

		// The read the pump makes, on the same subject, inside the window.
		live := g.Queue().LiveApprovals()
		liveCalls++
		if len(live) == 0 {
			t.Fatalf("第 %d 次工具调用已出卡（corr=%s）但 LiveApprovals() 返回空行：%s 秒窗口的读数前提不成立",
				i, corr, workPeakWindow148)
		}
		found := false
		for _, row := range live {
			if row.CorrelationID == p.CorrelationID {
				found = true
				if row.Decision.Tool != "fs.write" {
					t.Errorf("窗口内读到 corr=%s 的行 tool=%q，期望 %q", row.CorrelationID, row.Decision.Tool, "fs.write")
				}
			}
		}
		if !found {
			t.Fatalf("LiveApprovals() 的行里没有这次工具调用的 corr=%s（行数 %d）", p.CorrelationID, len(live))
		}

		// Settle through the real native route, fail-closed direction, so the
		// loop can run the next tool call and the gate exits at depth zero.
		if err := g.DecideFromNative(context.Background(), approval.Request{
			CorrelationID: p.CorrelationID,
			Allow:         false,
			Reason:        "148 AC#2 取数循环：一律拒绝",
		}); err != nil {
			t.Fatalf("拒绝路由失败（corr=%s）：%v", p.CorrelationID, err)
		}
		got := <-res
		if got.err != nil {
			t.Fatalf("被拒绝的工具调用不得表现为宿主故障：%v", got.err)
		}
		toolCalls++
	}

	t.Logf("148 AC#2 读数：5s 窗口内 工具调用=%d 次，LiveApprovals=%d 次（subject=真桥接真闸门真待决卡）",
		toolCalls, liveCalls)
	if toolCalls < 1 {
		t.Fatalf("5 秒窗口里工具调用 0 次 —— subject 没有真的跑起工具调用，这枚证明不成立")
	}
	if liveCalls < 1 {
		t.Fatalf("5 秒窗口里 LiveApprovals 0 次被调用 —— 票面要求的至少 1 次不成立")
	}
	// The clock never ran in this loop, but the gate shares the fake clock with
	// nothing else; advance it so any lingering timer is observably spent
	// before cleanup checks the depth.
	clk.Advance(time.Second)
}
