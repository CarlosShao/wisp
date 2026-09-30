package main

// Ticket 197 AC#5 - 「子代理永不自批」 turned from a sentence into a refusal a
// machine can see.
//
// WHY THIS FILE EXISTS. The 09-29 现读普查 (只读腿 197-c2, evidence
// .scratch/wisp/probes/197/r5c/census.md, transcribed onto the ticket face) judged
// AC#5 〔仅文档〕 and named the reason precisely: the child structurally cannot
// reach an answer - the interface handed down to every tool is tools.Gate, whose
// whole method set is PendingWindow and PendingApproval (internal/tools/gate.go),
// and the single-use native proof is handed to exactly one recipient, the injected
// UI (internal/agent/approval/approval.go's grant comment, gate.go's promptFor) -
// BUT 「造一枚自批假腿 => 要红」 had no ruler: grant binds the EVENT, not the
// PERSON (bindDigest(corr, taskID, tool, level, seq, args), approval.go:253), and
// Request.Source is documented as 「logged and never consulted」 (approval.go:231).
// The census's own conclusion was that AC#5 stays unchecked until a positive
// control exists, and it named the cheapest legal seam: `package main` in cmd/wisp,
// the one place that holds both the concrete *approval.Gate and the spawn path, and
// the assertion 「记下 grant 后 Native().Allow(childCorr, capturedGrant) 必须失败」.
//
// WHAT THIS FILE ACTUALLY MEASURES, in the order the legs below run.
//
//	(1) 一枚真子代理 - 名册里 Kind=subagent 的那一行 - 用自己的 task id 走真桥
//	    撞出一张真 L2 卡；这张卡能被宿主那一侧的令牌允许（文件真的落了盘），
//	    所以后面每一发「被拒」都不是一句空话；
//	(2) 同一个孩子，用它在场时能读到的每一枚载具去凑一枚令牌，每一发都是
//	    ErrBadGrant：空令牌、假令牌、自称 Source="subagent-self" 的原生路由、
//	    别张卡的令牌、已结卡的重放；
//	(3) 一张真令牌一旦被面板路线递出去，它就被按泄露烧掉，连原生侧也花不出去了
//	    （ErrPanelAllow -> ErrBadGrant），收口是宿主的拒绝，文件始终不存在；
//	(4) 那枚令牌从来没出现在孩子读得到的任何一处：名册行（工具侧与面板侧两份）、
//	    它自己那一页的流式正文、这一程自己发布的快照字节、stdout/stderr、
//	    以及落了盘的持久日志；
//	(5) 允许出口连"可达"这件事都不成立：从子代理装配能静态走到的每一个类型
//	    （tools.SubagentDeps / tools.TaskDeps / tools.Options / agent.Options）里
//	    没有一个声明 allow 类方法，tools.Gate 的方法集仍然只有那两枚「问」的动作。
//
// WHAT IT DELIBERATELY DOES NOT CLAIM.
//   - It does not claim the child's MODEL emitted the tool call: mockllm echoes and
//     cannot script a tool_call, so the call is dispatched through the assembled
//     bridge in the child's own identity (TaskID == the roster's subagent row) -
//     the same seam internal/tools/subagent_197.go's own depth check keys on. A
//     model-emitted variant is a mockllm capability question, not a 197 AC#5 one.
//   - Leg (5) walks STATIC types. The one field whose dynamic target it cannot see
//     is tools.Options.Gate, whose static type is the two-method ask-only seam -
//     which is exactly why the same leg pins that interface's method set by name.
//     Adding an answer verb to tools.Gate, or a new approval-shaped field to any of
//     the four roots, is red here; that is the reachability half, and (1)-(4) are
//     the behavioural half.
//   - Nothing here loosens or re-judges AC#6's two cells, which the 09-29 census
//     moved out of this ticket to 票 220 and 票 221.
//
// NO PRODUCTION CODE CHANGED. Same posture as 197-r4's blocked-carrier leg: the
// property was already structural, this file only makes it observable. The
// mutations that should turn these legs red are listed at the bottom.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/tools"
)

// selfApp197ProbeContent is the body a host answer is supposed to make appear. It
// names the case so a reader of the file on disk knows who wrote it.
const selfApp197ProbeContent = "written under one host-side native answer (ticket 197 AC#5)\n"

// childWrite197 is one tool call dispatched through the assembled bridge in a
// CHILD task's own identity, and the channel its outcome lands on.
type childWrite197 struct {
	done    chan struct{}
	text    string
	isError bool
	execErr error
}

// startChildWrite197 fires the call and returns immediately. The context is bounded
// on purpose: an unanswered L2 would otherwise sit on the C18 300s default, and a
// bounded context makes the abandonment (which resolves as a reject) the test's own
// act rather than the queue's. 30s is far above this case's own polling granularity
// (2ms) and far below the package's timeout, so a stuck leg fails loudly instead of
// eating the whole run.
func startChildWrite197(rt *agentRuntime, taskID, path string) *childWrite197 {
	w := &childWrite197{done: make(chan struct{})}
	args, err := json.Marshal(map[string]string{
		"path": filepath.ToSlash(path), "content": selfApp197ProbeContent,
	})
	if err != nil {
		w.execErr = err
		close(w.done)
		return w
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	go func() {
		defer cancel()
		defer close(w.done)
		out, e := rt.bridge.Execute(ctx, agent.ToolRequest{
			TaskID: taskID, CorrelationID: taskID,
			CallID: "self197-" + filepath.Base(path),
			Name:   "fs.write", Args: args,
		})
		w.text, w.isError, w.execErr = out.Text, out.IsError, e
	}()
	return w
}

// waitForChildCard197 polls the run's OWN native ledger until a live L2 card belongs
// to the given task, and hands back that card - which is how this file gets the
// proof the host was handed without any of it being typed here. The ledger is the
// only place a grant exists at all (consoleApprovalUI.Prompt books it from the
// Prompt the gate displayed), so reading it is the honest native-side vantage.
//
// The level test is load-bearing: an L1 window is booked into the same ledger with
// an empty Grant (gate.go's promptFor passes "" on that route), so matching a card
// by task id alone could hand this leg a countdown it never asked for.
func waitForChildCard197(live *nativeCards, taskID string, limit time.Duration) (approval.ReplyCard, bool) {
	deadline := time.Now().Add(limit)
	for time.Now().Before(deadline) {
		for _, c := range live.pending() {
			if c.TaskID == taskID && c.Level == "L2" {
				return c, true
			}
		}
		time.Sleep(2 * time.Millisecond)
	}
	return approval.ReplyCard{}, false
}

// selfApp197Readings carries every measurement the probe makes, so no t.Fatal ever
// runs on a non-test goroutine (the shape subagent_carrier_197_test.go's spawn197
// already uses for the same reason).
type selfApp197Readings struct {
	mu sync.Mutex

	err error // a setup failure: named, and fatal in the test goroutine

	childID   string
	spawnText string

	corr1, grant1 string
	corr2, grant2 string

	// carriers sampled while card 1 was still on screen
	pendingPacket []byte
	rosterDump    string
	chatChunkDump string
	packetSeen    bool
	pendingCount  int

	// route answers (each is an error the gate returned; nil means it should not be)
	emptyProofErr error
	bogusProofErr error
	selfClaimErr  error
	foreignErr    error
	replayErr     error
	hostAllowErr  error
	panelAllowErr error
	burnedErr     error
	hostRejectErr error

	write1Text     string
	write1IsError  bool
	write1ExecErr  error
	write2IsError  bool
	write2ExecErr  error
	allowedLanded  bool
	refusedLanded  bool
	snapClearedCar bool
}

func (r *selfApp197Readings) setErr(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err == nil {
		r.err = err
	}
}

// Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands is AC#5's positive
// control on a production assembly: one real subagent, two real L2 cards in its own
// task identity, every assemblable self-approval route refused, and the host's own
// answer proven to land - in the same run, in that order.
func Test197SubagentSelfApprovalIsRefusedAndTheHostAnswerLands(t *testing.T) {
	f := newRunFixture(t, "openai-chat")
	// The window this case needs is "the root loop is still running while a child's
	// card is being answered", because runTextTask closes the store the moment the
	// root returns (run.go:207-220). 2.5s of mockllm latency is that window; the
	// probe below finishes its two exchanges in milliseconds.
	f.srv.Control(t, "/__control/latency", `{"ms":2500}`)
	outside := t.TempDir() // nothing allowlists it: R2 => L2, the shape 票 201's seam cases use
	file1 := filepath.Join(outside, "child197-under-host-allow.txt")
	file2 := filepath.Join(outside, "child197-under-self-approval.txt")

	rec := &selfApp197Readings{}
	sp := newSpawn197()
	probeDone := make(chan struct{})

	f.rtHook = func(rt *agentRuntime) {
		go func() {
			defer close(probeDone)
			go sp.launch(rt, carrier197ChildPrompt, 1)
			probeSelfApproval197(t, rt, rec, file1, file2)
		}()
	}

	code := f.run(carrier197RootTask)
	<-probeDone
	<-sp.done

	rec.mu.Lock()
	defer rec.mu.Unlock()

	if code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, f.out.String(), f.err.String())
	}
	if rec.err != nil {
		t.Fatalf("the self-approval probe never completed its setup: %v\nstdout:\n%s\nstderr:\n%s",
			rec.err, f.out.String(), f.err.String())
	}
	if rec.childID == "" {
		t.Fatal("no subagent roster row was ever read, so nothing below measured a child")
	}

	// ---------------------------------------------------------------------
	// (2) 自批的每一发都被拒 - 而每一发的理由都是「令牌不对」，不是「这条路由没在跑」
	// ---------------------------------------------------------------------
	for _, tc := range []struct {
		name string
		err  error
	}{
		{"空令牌（孩子在任何载具上都读不到的那一枚）", rec.emptyProofErr},
		{"假令牌（凑出来的形状）", rec.bogusProofErr},
		{"自称 Source=subagent-self 的原生路由", rec.selfClaimErr},
		{"别张卡的令牌（跨卡借证）", rec.foreignErr},
		{"面板路线递出真令牌", rec.panelAllowErr},
		{"被烧掉的令牌再走原生侧", rec.burnedErr},
	} {
		if tc.err == nil {
			t.Errorf("%s：这一发竟然没有被拒。AC#5 要的正是「子代理拿不到允许」，"+
				"任何一发放行都等于给孩子开了一道自批的出口", tc.name)
			continue
		}
		if !errors.Is(tc.err, approval.ErrBadGrant) && !errors.Is(tc.err, approval.ErrPanelAllow) {
			t.Errorf("%s：被拒的理由是 %v，想要 ErrBadGrant 或 ErrPanelAllow。"+
				"「因为别的原因失败」不构成本票的证据", tc.name, tc.err)
		}
	}
	// 一张已经结算的卡不能被再花一次（重放）。ErrNotPending 也算被拒，但不是这里要的：
	// 卡片离开队列后 byID 里就没有它了，正确读数只能是 ErrUnknownCorrelation。
	if rec.replayErr == nil || !errors.Is(rec.replayErr, approval.ErrUnknownCorrelation) {
		t.Errorf("已结卡的重放回了 %v，想要 ErrUnknownCorrelation：一张答过的卡留着可用令牌，"+
			"就是孩子（或任何别的东西）第二次拿到允许的那条路", rec.replayErr)
	}

	// ---------------------------------------------------------------------
	// (1) 正控：同一枚孩子的卡，宿主那一侧的真令牌确实能允许，文件真的落了盘
	// ---------------------------------------------------------------------
	if rec.hostAllowErr != nil {
		t.Errorf("宿主侧允许这张子代理的卡竟然失败了：%v（正控不成立，上面每一发「被拒」就都是空气）",
			rec.hostAllowErr)
	}
	if rec.write1ExecErr != nil {
		t.Errorf("被允许的调用本身报错：%v", rec.write1ExecErr)
	}
	if rec.write1IsError {
		t.Errorf("被允许的写回了错误结果：%s", rec.write1Text)
	}
	if !rec.allowedLanded {
		t.Errorf("正控失败：%s 没有被写出来。允许的那一发没有真的执行，"+
			"这一格的判据就退回到「填了字段但没人读」那一族（票 181 AC#7）", file1)
	}
	if rec.write2IsError == false {
		t.Errorf("一张允许只被面板路线递出去过的卡竟然执行了：%+v", rec.write2IsError)
	}
	if rec.write2ExecErr != nil {
		t.Errorf("第二发（被烧掉令牌那一张）在桥这一层报了 Go 错误：%v，想要的是「结果里带拒绝」而不是「调用炸了」",
			rec.write2ExecErr)
	}
	if rec.refusedLanded {
		t.Errorf("%s 被写出来了，可它那张卡的令牌是在不受信路线上递出去的：面板来源的允许（AGENTS.md §1.2 铁律）"+
			"一旦能落地，本文件上面每一发被拒都没有意义了", file2)
	}
	if rec.hostRejectErr != nil {
		t.Errorf("收口这一发（宿主拒绝第二张卡）失败：%v", rec.hostRejectErr)
	}

	// ---------------------------------------------------------------------
	// (4) 那枚令牌从来没出现在孩子读得到的任何一处
	// ---------------------------------------------------------------------
	carriers := []struct {
		name  string
		bytes string
	}{
		{"这一程自己发布的快照字节（卡还挂在队上时取的）", string(rec.pendingPacket)},
		{"名册行（工具侧 TaskOutput 全字段）", rec.rosterDump},
		{"它自己那一页的流式正文（StreamLog 全字段）", rec.chatChunkDump},
		{"task.spawn 回给父任务模型的正文", rec.spawnText},
		{"stdout", f.out.String()},
		{"stderr（含 gate 自己的审计行）", f.err.String()},
		{"持久日志目录", readRunLogDir197(t, f.dir)},
	}
	for _, c := range carriers {
		for _, proof := range []string{rec.grant1, rec.grant2} {
			if proof == "" {
				t.Fatalf("本用例没能从宿主账上读到一枚真令牌，载具扫描就是空转：%+v", c.name)
			}
			if strings.Contains(c.bytes, proof) {
				t.Errorf("%s 里出现了原生一次性令牌本身。孩子读得到的每一处都不该有它："+
					"bindDigest 绑的是事不是人（approval.go:253），所以「拿不到」是这一层的唯一防线", c.name)
			}
		}
	}
	if rec.pendingCount != 1 {
		if !rec.packetSeen {
			t.Error("这一程从来没有留下可读的快照（lastPanelSnapshot 回了 false），" +
				"上面那一整格载具扫描因此没有读到过「卡还挂着」的字节")
		} else {
			t.Errorf("卡还挂着时快照里有 %d 枚 pending，想要 1 枚（本用例只给孩子挂了一张）", rec.pendingCount)
		}
	}

	// 审计侧的三条实话：被拒的每一发都留了痕，允许的那一发也留了。
	audit := f.err.String()
	for _, want := range []string{
		"approval: FORGED-OR-STALE allow rejected corr=" + rec.corr1,
		"approval: PANEL-ALLOW-REJECTED corr=" + rec.corr2,
		"approval: ANSWER-ALLOW corr=" + rec.corr1 + " tool=fs.write route=native decision=allow",
	} {
		if !strings.Contains(audit, want) {
			t.Errorf("审计里找不到 %q。被拒的与放行的必须是两种读数，缺一种这一格就只剩推断：\n%s", want, audit)
		}
	}
	// This file never calls publishPanelSnapshot, so the packet whose bytes were
	// scanned above was built by the run's own publisher. What this statement is
	// NOT: the sha-to-ledger tie 197-r3b used. That tie is only reachable on a
	// packet that is still the pump's newest when the run books it, and this case
	// samples its packet while the card is up - the run's exit publish overwrites
	// rt.lastSnapBytes with a wider one, so demanding the equality here would be a
	// test that samples a state no production publisher ever books under that name.
	// The booking half is covered by that family in
	// subagent_carrier_197_test.go; naming the gap is this leg's job.
	if !rec.snapClearedCar {
		t.Errorf("没在快照里读到过这张卡的名字 %q，pending 那一发就无从判断", rec.corr1)
	}
	if len(ledgerSummaries145(t, f.dir)) == 0 {
		t.Error("这一程没有把任何一枚快照落进持久台账，本腿读的载具因此不能算生产发布的东西")
	}
	t.Logf("child=%s corr1=%s corr2=%s | 空令牌=%v | 假令牌=%v | 自称来源=%v | 借证=%v | 重放=%v | "+
		"面板递证=%v | 烧后再试=%v | 宿主允许=%v | 落盘=%v/%v",
		rec.childID, rec.corr1, rec.corr2, rec.emptyProofErr, rec.bogusProofErr, rec.selfClaimErr,
		rec.foreignErr, rec.replayErr, rec.panelAllowErr, rec.burnedErr, rec.hostAllowErr,
		rec.allowedLanded, rec.refusedLanded)
}

// probeSelfApproval197 runs the whole sequence above on the hook's goroutine. It
// records, it never fails: a t.Fatal off the test goroutine would abort the wrong
// thing (the same rule spawn197 follows).
func probeSelfApproval197(t *testing.T, rt *agentRuntime, rec *selfApp197Readings, file1, file2 string) {
	t.Helper()
	rootID, err := waitRootRow197(rt)
	if err != nil {
		rec.setErr(err)
		return
	}
	childID, err := waitChildRow197(rt, rootID)
	if err != nil {
		rec.setErr(err)
		return
	}
	// The child's own admission hook already registered it and will revoke it when
	// it joins; this registration is what makes the card reachable independent of
	// when that happens, exactly as subagent_blocked_197_test.go does for the same
	// reason. It is revoked on the way out of this probe.
	revoke := rt.gate.AdmitTextTask(childID)
	defer revoke()

	// ---- card 1: 允许出口被逐头发拒，然后宿主那一发落地 ---------------------
	w1 := startChildWrite197(rt, childID, file1)
	card1, ok := waitForChildCard197(rt.liveCards, childID, 20*time.Second)
	if !ok {
		<-w1.done
		rec.setErr(fmt.Errorf("孩子在 %s 名下始终没有挂出一张 L2 卡（本用例的全部判据都挂在它上面）；"+
			"这一程已经显示了 %d 枚确认卡", childID, rt.windowCount()))
		return
	}
	rec.mu.Lock()
	rec.childID, rec.corr1, rec.grant1 = childID, card1.CorrelationID, card1.Grant
	rec.mu.Unlock()

	// Sample the carriers WHILE the card is up: a leak would show there, and 197-r4's
	// lesson is that a field read after the fact cannot tell a report from a constant.
	var (
		snapBytes []byte
		pending   int
		seen      bool
		named     bool
	)
	for i := 0; i < 600; i++ {
		snap, data, ok2 := rt.lastPanelSnapshot()
		pending, seen = len(snap.Pending), ok2
		if ok2 && len(snap.Pending) == 1 && snap.Pending[0].CorrelationID == card1.CorrelationID {
			named = true
			snapBytes = append([]byte(nil), data...)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	rowOut, rowOK := rt.tasks.Look(childID)
	rec.mu.Lock()
	rec.pendingPacket, rec.pendingCount, rec.snapClearedCar = snapBytes, pending, named
	rec.packetSeen = seen
	rec.rosterDump = fmt.Sprintf("%+v ok=%v", rowOut, rowOK)
	rec.chatChunkDump = fmt.Sprintf("%+v", rt.stream.Chunks())
	rec.mu.Unlock()

	rec.mu.Lock()
	rec.emptyProofErr = rt.gate.Native().Allow(context.Background(), card1.CorrelationID, "")
	rec.bogusProofErr = rt.gate.Native().Allow(context.Background(), card1.CorrelationID,
		bogusProof197(card1.Grant))
	rec.selfClaimErr = rt.gate.DecideFromNative(context.Background(), approval.Request{
		CorrelationID: card1.CorrelationID, Allow: true, Grant: "", Source: "subagent-self",
	})
	rec.mu.Unlock()

	// The host's own answer: this is the reading that makes every refusal above a
	// refusal rather than an artifact of a queue nobody can reach.
	hostErr := rt.liveCards.h.Allow(context.Background(), card1.CorrelationID)
	<-w1.done
	_, statErr := os.Stat(file1)

	rec.mu.Lock()
	rec.hostAllowErr, rec.write1Text, rec.write1IsError, rec.write1ExecErr =
		hostErr, w1.text, w1.isError, w1.execErr
	rec.allowedLanded = statErr == nil
	rec.mu.Unlock()

	// A settled card has no live item under its name any more: the only honest
	// reading is "unknown correlation", never "still spendable".
	replay := rt.gate.Native().Allow(context.Background(), card1.CorrelationID, card1.Grant)
	rec.mu.Lock()
	rec.replayErr = replay
	rec.mu.Unlock()

	// ---- card 2: 真令牌在不受信路线上露一次面，就再也花不出去了 --------------
	w2 := startChildWrite197(rt, childID, file2)
	card2, ok := waitForChildCard197(rt.liveCards, childID, 20*time.Second)
	if !ok {
		<-w2.done
		rec.setErr(fmt.Errorf("孩子的第二张卡没有挂出来（%s）", childID))
		return
	}
	rec.mu.Lock()
	rec.corr2, rec.grant2 = card2.CorrelationID, card2.Grant
	// 借别张卡的证：card1 的令牌此刻已经花掉/不存在了，但它连"形状"都不属于这张卡。
	rec.foreignErr = rt.gate.Native().Allow(context.Background(), card2.CorrelationID, card1.Grant)
	rec.mu.Unlock()

	panelErr := rt.gate.DecideFromPanel(context.Background(), approval.Request{
		CorrelationID: card2.CorrelationID, Allow: true, Grant: card2.Grant, Source: "subagent-via-panel",
	})
	burned := rt.gate.Native().Allow(context.Background(), card2.CorrelationID, card2.Grant)
	rejectErr := rt.liveCards.h.Reject(context.Background(), card2.CorrelationID, "本用例不收这张卡")
	<-w2.done
	_, stat2 := os.Stat(file2)

	rec.mu.Lock()
	rec.panelAllowErr, rec.burnedErr, rec.hostRejectErr = panelErr, burned, rejectErr
	rec.write2IsError, rec.write2ExecErr, rec.refusedLanded = w2.isError, w2.execErr, stat2 == nil
	rec.mu.Unlock()
}

// bogusProof197 builds one grant-shaped string that is NOT this card's proof, at
// run time: a value typed as a literal in a test source is exactly the shape
// d22scan's ban #3 exists to keep out of a repository.
func bogusProof197(real string) string {
	if len(real) < 8 {
		return "grant_bogus_value"
	}
	cut := real[:len(real)-4]
	return cut + "-xxx"
}

// readRunLogDir197 returns the bytes of this run's durable audit files. The grant
// must not outlive the process: rt.auditf writes every gate line to stderr AND to
// the rolling JSONL under <data>/logs (run.go:732-736), so the on-disk copy is a
// carrier of its own. Missing directory is not an error - the sink is explicitly
// allowed to degrade (run.go:199-205) - and an unreadable file is.
func readRunLogDir197(t *testing.T, dataDir string) string {
	t.Helper()
	dir := logSinkDir(dataDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	var b strings.Builder
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		src, rerr := os.ReadFile(filepath.Join(dir, e.Name()))
		if rerr != nil {
			t.Errorf("持久日志 %s 读不出来：%v（载具扫描不能跳过一格）", e.Name(), rerr)
			continue
		}
		b.Write(src)
	}
	return b.String()
}

// ---------------------------------------------------------------------------
// Leg (5): the reachability half - no type a subagent is assembled from declares
// an allow door, and the one approval interface a tool can hold asks and never
// answers.
//
// This is the structural answer to 「子代理不许自带允许出口」: (1)-(4) above measure
// that a child cannot get today's proof, this measures that there is nothing to
// get. It is a name-and-signature reading over STATIC types, so the day someone
// hands the tool layer a *approval.Gate / *approval.Replies, or widens tools.Gate
// with an answer verb, the leg goes red on its own - which is the failure the
// 09-29 census described as 「今天没有任何尺会红」.
// ---------------------------------------------------------------------------

// allowDoorMethods197 is deliberately narrow: the five names that can put an ALLOW
// into motion (approval.Gate.Native / its two routers / Queue.grantNonce, and
// approval.Replies.Allow). 「Reject」 and 「Veto」 are not here on purpose - a
// refusal a child could reach would be its own finding (票 220's neighbourhood),
// and mixing the two directions would make this leg's red ambiguous. Measured at
// this anchor, no type outside internal/agent/approval declares any of these five
// names; the planted control below is what keeps that claim from being a sentence.
var allowDoorMethods197 = []string{"Allow", "Native", "DecideFromNative", "DecideFromPanel", "GrantNonce"}

// toolAssembledRoots197 are the four production types a subagent's runtime is put
// together from. A child's loop is a copy of agent.Options (SubagentDeps.BaseOptions),
// its tool surface is the bridge built from tools.Options, and the spawn path itself
// holds tools.SubagentDeps / tools.TaskDeps.
var toolAssembledRoots197 = []struct {
	name string
	typ  reflect.Type
}{
	{"tools.SubagentDeps", reflect.TypeOf(tools.SubagentDeps{})},
	{"tools.TaskDeps", reflect.TypeOf(tools.TaskDeps{})},
	{"tools.Options", reflect.TypeOf(tools.Options{})},
	{"agent.Options", reflect.TypeOf(agent.Options{})},
}

func Test197NoAllowDoorIsReachableFromASubagentsAssembly(t *testing.T) {
	var findings []string
	for _, root := range toolAssembledRoots197 {
		findings = append(findings, findAllowDoors197(root.name, root.typ, 0, map[string]bool{})...)
	}
	if len(findings) != 0 {
		t.Errorf("从子代理的装配里能静态走到声明了允许出口的类型：%v。"+
			"孩子只要拿得到这些名字里的任何一个，「子代理永不自批」就退化成一句约定", findings)
	}

	// The one approval surface a tool can hold at all, and its whole method set.
	gateNames := methodNames197(reflect.TypeOf((*tools.Gate)(nil)).Elem())
	if len(gateNames) != 2 || !hasMethod197(gateNames, "PendingWindow") || !hasMethod197(gateNames, "PendingApproval") {
		t.Errorf("tools.Gate 的方法集是 %v，想要 {PendingWindow, PendingApproval}。"+
			"这是每个工具都摸得到的一枚接口：往上添任何一枚「答」的动作，就是给孩子开出口", gateNames)
	}

	// Positive control for the walker itself (本仓纪律：负向尺必配"种 X 必响"的正控).
	planted := findAllowDoors197("planted", reflect.TypeOf(selfApprovalDoor197{}), 0, map[string]bool{})
	if len(planted) < 2 {
		t.Fatalf("这架扫描器在种了一扇门的载体上只读出 %v：负向判据不能由一架照不见门的仪器来签发", planted)
	}
	var verbs []string
	for _, p := range planted {
		verbs = append(verbs, p)
	}
	if !strings.Contains(strings.Join(verbs, "|"), "Allow") ||
		!strings.Contains(strings.Join(verbs, "|"), "Native") {
		t.Errorf("种下的两枚出口名字里读出了 %v，想要 Allow 与 Native 都在", verbs)
	}
	// The planted door's own calls exist so a linter cannot call them dead code, and
	// so the reading above is about a method that really answers something.
	var door nativeDoor197
	if err := door.Allow(context.Background(), "c", "g"); err == nil {
		t.Error("种下的门回了一个 nil 错误 - 正控的读数不该依赖它的行为，但也不该是一枚空方法")
	}
	_ = door.Native()
}

// selfApprovalDoor197 is the planted shape: a deps-style struct one field away from
// handing a tool layer a native answer surface.
type selfApprovalDoor197 struct{ holder carrierWithDoor197 }

type carrierWithDoor197 struct{ Inner nativeDoor197 }

type nativeDoor197 struct{}

func (nativeDoor197) Allow(_ context.Context, corr, proof string) error {
	return fmt.Errorf("planted door, not a route: %s/%s", corr, proof)
}

func (nativeDoor197) Native() string { return "planted" }

// findAllowDoors197 walks a type graph and names every type that declares one of
// allowDoorMethods197. Interfaces are read by their declared methods (that is how a
// future `Approvals answerer` field would be caught) and are not followed further,
// because a static type has no dynamic target to walk.
func findAllowDoors197(path string, typ reflect.Type, depth int, seen map[string]bool) []string {
	if typ == nil || depth > 4 {
		return nil
	}
	switch typ.Kind() {
	case reflect.Pointer:
		return findAllowDoors197(path+"*", typ.Elem(), depth, seen)
	case reflect.Slice, reflect.Array:
		return findAllowDoors197(path+"[]", typ.Elem(), depth+1, seen)
	case reflect.Map:
		return findAllowDoors197(path+"{}", typ.Elem(), depth+1, seen)
	case reflect.Func, reflect.Chan:
		return nil
	}
	if key := typ.PkgPath() + "/" + typ.Name(); key != "/" {
		if seen[key] {
			return nil
		}
		seen[key] = true
	}
	var hits []string
	for _, name := range methodNames197(reflect.PointerTo(typ)) {
		if hasMethod197(allowDoorMethods197, name) {
			hits = append(hits, fmt.Sprintf("%s: %s declares %s()", path, typ.String(), name))
		}
	}
	if typ.Kind() == reflect.Interface {
		for _, name := range methodNames197(typ) {
			if hasMethod197(allowDoorMethods197, name) {
				hits = append(hits, fmt.Sprintf("%s: interface %s declares %s()", path, typ.String(), name))
			}
		}
		return dedupe197(hits)
	}
	if typ.Kind() != reflect.Struct {
		return dedupe197(hits)
	}
	for i := 0; i < typ.NumField(); i++ {
		hits = append(hits, findAllowDoors197(path+"."+typ.Field(i).Name, typ.Field(i).Type, depth+1, seen)...)
	}
	return dedupe197(hits)
}

// methodNames197 reads a method set as names, sorted by declaration order which
// reflect already normalises for non-interface types.
func methodNames197(typ reflect.Type) []string {
	if typ == nil {
		return nil
	}
	out := make([]string, 0, typ.NumMethod())
	for i := 0; i < typ.NumMethod(); i++ {
		out = append(out, typ.Method(i).Name)
	}
	return out
}

func hasMethod197(names []string, want string) bool {
	for _, n := range names {
		if n == want {
			return true
		}
	}
	return false
}

func dedupe197(in []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(in))
	for _, v := range in {
		if seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	return out
}

// ---------------------------------------------------------------------------
// MUTATION CONTROLS (what the acceptance leg should plant, and what must go red).
// None of these were run by this leg: the gate window belonged to another agent.
//
//	M1 internal/agent/approval/ui.go: add Grant to PanelItem -> (4) red on the
//	    snapshot-bytes carrier, and internal/panel's own boundary family too.
//	M2 cmd/wisp/run.go: print p.Grant in consoleApprovalUI.Prompt -> (4) red on
//	    stdout. This is the "convenience" shape the leg exists for.
//	M3 internal/tools: add a sixth field to SubagentDeps carrying the gate or the
//	    reply seam -> leg (5) red, plus 197's own field-count nail in
//	    internal/tools/subagent_197_test.go.
//	M4 internal/tools/gate.go: add an answer verb to the Gate interface -> leg (5)
//	    red on the method-set pin.
//	M5 internal/agent/approval/gate.go: consult Request.Source for authority (e.g.
//	    honour Source=="subagent-self") -> (2) red: the self-claimed source stops
//	    being decoration and starts buying a refusal or an allow.
//	M6 internal/agent/approval/replies.go: let Replies.Allow spend without the
//	    card's own proof -> (1)/(2) red.
// ---------------------------------------------------------------------------
