package approval_test

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// The L1 pre-execution BLOCK window (B1-corrected: a block, not an "undo").

func TestL1WindowTimeoutMeansExecute(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui, approval.Options{Window: 3 * time.Second})
	g.AdmitTextTask(testTask)

	res := runWindow(t, g, context.Background(), l1Decision("C:/data/a.txt"))
	p := ui.wait(t)
	if p.Level != "L1" || p.CorrelationID != testCorr {
		t.Fatalf("prompt = %+v", p)
	}
	if p.Window != 3*time.Second {
		t.Errorf("Window=%v，期望 3s（2–3s 契约）", p.Window)
	}
	// B1: the wording must never promise an undo.
	if strings.Contains(p.Reason, "撤销成功") || strings.Contains(p.Reason, "可撤销") {
		t.Errorf("L1 提示不得承诺可撤销：%q", p.Reason)
	}

	clk.Advance(3*time.Second - time.Millisecond)
	select {
	case a := <-res:
		t.Fatalf("窗口尚未结束就返回了 %v（提前放行）", a)
	default:
	}
	clk.Advance(time.Millisecond)

	a := mustAnswer(t, res)
	if a.a != tools.AnswerTimeout {
		t.Fatalf("answer=%q，期望 timeout（超时即执行）：%s", a.a, a.why)
	}
	if len(ui.ofKind(approval.EventStarted)) != 1 {
		t.Errorf("放行后必须发出 started 事件，得 %v", ui.ofKind(approval.EventStarted))
	}
}

func TestL1WindowVetoedByEachLoadedChannel(t *testing.T) {
	for _, ch := range []approval.Channel{approval.ChannelBall, approval.ChannelEsc} {
		t.Run(string(ch), func(t *testing.T) {
			ui := newFakeUI()
			g, clk, _ := newGate(t, ui, approval.Options{Window: 3 * time.Second})
			g.AdmitTextTask(testTask)
			res := runWindow(t, g, context.Background(), l1Decision("C:/data/a.txt"))
			ui.wait(t)

			if err := g.Veto(approval.Veto{CorrelationID: testCorr, Channel: ch}); err != nil {
				t.Fatalf("Veto(%s): %v", ch, err)
			}
			a := mustAnswer(t, res)
			if a.a != tools.AnswerVeto {
				t.Fatalf("answer=%q，期望 veto：%s", a.a, a.why)
			}
			if !strings.Contains(a.why, "否决") {
				t.Errorf("否决理由必须回到模型可读的文案：%q", a.why)
			}
			clk.Advance(10 * time.Second)
			if len(ui.ofKind(approval.EventStarted)) != 0 {
				t.Errorf("已否决的调用不得进入执行：%+v", ui.all())
			}
		})
	}
}

// The ticket's demanded assertion: a channel that is not loaded must be
// REPORTED as unavailable - not silently skipped, and never pretended.
func TestL1WindowUnavailableVoiceChannelIsAnnounced(t *testing.T) {
	ui := newFakeUI()
	// Ball + Esc loaded, KWS (ticket 41) and panel (ticket 37) not.
	g, clk, reg := newGate(t, ui, approval.Options{Window: 3 * time.Second})
	g.AdmitTextTask(testTask)
	res := runWindow(t, g, context.Background(), l1Decision("C:/data/a.txt"))
	p := ui.wait(t)

	st, ok := channelText(p, approval.ChannelKWS)
	if !ok {
		t.Fatal("提示必须逐条列出四个否决通道")
	}
	if st.Loaded {
		t.Fatal("KWS 未加载却报为可用（这是在假装通道存在）")
	}
	if !strings.Contains(st.Text, "语音取消不可用") {
		t.Errorf("KWS 通道文案必须原样包含契约要求的「语音取消不可用」，得 %q", st.Text)
	}
	if len(p.Channels) != 4 {
		t.Errorf("四个否决通道都得列出，得 %d", len(p.Channels))
	}
	if ps, _ := channelText(p, approval.ChannelPanel); ps.Loaded {
		t.Error("面板（票 37）未接入却报为可用")
	}

	// An attempt to use the unloaded channel is refused with the same honest
	// message, and it does NOT cancel the window.
	err := g.Veto(approval.Veto{CorrelationID: testCorr, Channel: approval.ChannelKWS})
	if err == nil {
		t.Fatal("未加载的语音通道竟然否决成功")
	}
	if !errors.Is(err, approval.ErrChannelUnavailable) {
		t.Fatalf("err=%v，期望 ErrChannelUnavailable", err)
	}
	var ce *approval.ChannelError
	if !errors.As(err, &ce) || !strings.Contains(ce.Msg, "语音取消不可用") {
		t.Fatalf("错误必须携带用户可见的不可用文案，得 %v", err)
	}
	// The strip is told out loud too (a failed attempt must not be silent).
	if ev := ui.ofKind(approval.EventWarning); len(ev) == 0 ||
		!strings.Contains(ev[0].Text, "语音取消不可用") {
		t.Errorf("界面必须收到「语音取消不可用」提示事件，得 %+v", ev)
	}
	select {
	case a := <-res:
		t.Fatalf("未加载通道的否决竟然结束了窗口：%+v", a)
	default:
	}

	// Once ticket 41 loads the model, the same channel really does veto.
	reg.SetLoaded(approval.ChannelKWS, true)
	if err := g.Veto(approval.Veto{CorrelationID: testCorr, Channel: approval.ChannelKWS}); err != nil {
		t.Fatalf("KWS 加载后否决失败：%v", err)
	}
	if a := mustAnswer(t, res); a.a != tools.AnswerVeto {
		t.Fatalf("answer=%q，期望 veto：%s", a.a, a.why)
	}
	clk.Advance(10 * time.Second)
}

func TestL1WindowVetoOnForeignCorrelationCannotCancel(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui, approval.Options{Window: 3 * time.Second})
	g.AdmitTextTask(testTask)
	res := runWindow(t, g, context.Background(), l1Decision("C:/data/a.txt"))
	ui.wait(t)

	// C18: a reply is routed by correlation id, so a stray click cannot land
	// on someone else's request.
	err := g.Veto(approval.Veto{CorrelationID: "corr-someone-else", Channel: approval.ChannelBall})
	if !errors.Is(err, approval.ErrUnknownCorrelation) {
		t.Fatalf("err=%v，期望 ErrUnknownCorrelation", err)
	}
	clk.Advance(3 * time.Second)
	if a := mustAnswer(t, res); a.a != tools.AnswerTimeout {
		t.Fatalf("别人的 correlation_id 不该取消本窗口，得 %+v", a)
	}
}

func TestL1WindowUnreachableUIFailsClosed(t *testing.T) {
	ui := newFakeUI()
	ui.setErr(errors.New("球窗口创建失败"))
	g, _, _ := newGate(t, ui, approval.Options{Window: 3 * time.Second})
	g.AdmitTextTask(testTask)

	a, why := g.PendingWindow(context.Background(), l1Decision("C:/data/a.txt"))
	if a != tools.AnswerReject {
		t.Fatalf("answer=%q，宿主不可达必须 reject：%s", a, why)
	}
	if !strings.Contains(why, "fail-closed") {
		t.Errorf("文案要说明是 fail-closed：%q", why)
	}
}

// D47: the Path C realtime brain has zero tool permissions, so nothing from
// that path may reach the gates - and the gate cannot be talked past by a
// request that merely claims a task id.
func TestUnadmittedTaskIsRefusedBeforeQueueing(t *testing.T) {
	ui := newFakeUI()
	g, _, _ := newGate(t, ui, approval.Options{Window: 3 * time.Second})

	for name, call := range map[string]func(context.Context, tools.Decision) (tools.Answer, string){
		"L1": g.PendingWindow,
		"L2": g.PendingApproval,
	} {
		a, why := call(context.Background(), l2Decision("C:/data/a.txt"))
		if a != tools.AnswerReject {
			t.Fatalf("%s: answer=%q，未登记任务必须拒绝：%s", name, a, why)
		}
		if !strings.Contains(why, "D47") && !strings.Contains(why, "Path C") {
			t.Errorf("%s: 文案要指出 D47：%q", name, why)
		}
	}
	if n := g.Queue().Depth(); n != 0 {
		t.Errorf("未经登记的请求不得进入审批队列，队列深度=%d", n)
	}
	if n := ui.count(); n != 0 {
		t.Errorf("不该弹出任何确认界面，得 %d", n)
	}
}

// D31: cancellation is not atomic. A veto after hand-off must produce the
// applied-steps report shape, and must never read like a clean undo.
func TestVetoAfterStartYieldsAppliedStepsReport(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui, approval.Options{Window: 3 * time.Second})
	g.AdmitTextTask(testTask)
	d := l1Decision("C:/data/a.txt")

	res := runWindow(t, g, context.Background(), d)
	ui.wait(t)
	clk.Advance(3 * time.Second)
	if a := mustAnswer(t, res); a.a != tools.AnswerTimeout {
		t.Fatalf("窗口应超时放行，得 %+v", a)
	}

	// The user hits the ball only after execution began.
	err := g.Veto(approval.Veto{CorrelationID: testCorr, Channel: approval.ChannelBall})
	if !errors.Is(err, approval.ErrAlreadyStarted) {
		t.Fatalf("执行后否决必须回报已开始执行，得 %v", err)
	}
	bus := g.Bus()
	v, ok := bus.Vetoed(testCorr)
	if !ok || v.Channel != approval.ChannelBall {
		t.Fatalf("运行中的工具必须能看到迟到的否决，得 %+v ok=%v", v, ok)
	}

	rep := bus.Report(d, tools.Result{
		Text:         "已中止",
		IsError:      true,
		AppliedSteps: []string{"写入临时文件 C:/data/.a.txt.tmp", "重命名覆盖 C:/data/a.txt"},
	})
	if !rep.StartedBeforeVeto || !rep.Vetoed {
		t.Fatalf("报告形状不对：%+v", rep)
	}
	if len(rep.AppliedSteps) != 2 {
		t.Errorf("AppliedSteps 必须原样带出：%+v", rep.AppliedSteps)
	}
	for _, want := range []string{"取消不是原子的", "已执行：写入临时文件", "已执行：重命名覆盖"} {
		if !strings.Contains(rep.Text, want) {
			t.Errorf("报告文案缺少 %q：%q", want, rep.Text)
		}
	}
	if strings.Contains(rep.Text, "已撤销") || strings.Contains(rep.Text, "已回滚") {
		t.Errorf("报告不得假装撤销：%q", rep.Text)
	}

	// A tool that reported nothing is the honest worst case, not a blank.
	rep2 := bus.Report(d, tools.Result{Text: "已中止", IsError: true})
	if !rep2.UnreportedSteps || !strings.Contains(rep2.Text, "无上报") {
		t.Errorf("工具未上报步骤时必须按最坏情况措辞：%+v %q", rep2.UnreportedSteps, rep2.Text)
	}
	bus.Complete(testCorr)
	if bus.Started(testCorr) {
		t.Error("Complete 之后应释放跟踪记录")
	}
}

// The veto has to reach the model as a failure the loop can continue from, not
// as a provider fault (bridge.Execute's rejection contract).
func TestVetoTravelsBackToTheLoopAsAFailure(t *testing.T) {
	ui := newFakeUI()
	g, clk, _ := newGate(t, ui, approval.Options{Window: 3 * time.Second})
	g.AdmitTextTask(testTask)
	dir := t.TempDir()
	ft := newFakeTool("fs.write", tools.CapFSWrite, risk.L1, "paths")
	b := bridgeFor(t, g, ft, []string{dir})

	var gerr error
	outCh := make(chan agent.ToolOutcome, 1)
	spawn(t, func() {
		out, err := b.Execute(context.Background(),
			agentToolRequest("corr-veto",
				argsJSON(t, slashPath(filepath.Join(dir, "a.txt")))))
		gerr = err
		outCh <- out
	})
	ui.wait(t)
	if err := g.Veto(approval.Veto{CorrelationID: "corr-veto", Channel: approval.ChannelBall}); err != nil {
		t.Fatalf("Veto: %v", err)
	}
	out := <-outCh
	clk.Advance(10 * time.Second)

	if gerr != nil {
		t.Fatalf("拒绝不是宿主故障（D37）：err=%v", gerr)
	}
	if !out.IsError {
		t.Fatal("否决必须以失败结果回给模型，否则循环不知道调用被拒")
	}
	if out.ErrorClass != string(observe.ClassUserRejected) {
		t.Errorf("error_class=%q，期望 %s", out.ErrorClass, observe.ClassUserRejected)
	}
	if !strings.Contains(out.Text, "否决") {
		t.Errorf("模型要看到否决原因：%q", out.Text)
	}
	if n := ft.count(); n != 0 {
		t.Errorf("被否决的调用绝不能执行，工具被调用 %d 次", n)
	}
}
