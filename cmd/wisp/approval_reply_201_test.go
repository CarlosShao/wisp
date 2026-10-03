package main

// Ticket 201 - the reply listener's production criteria, run against the stack
// `wisp run` assembles.
//
// What each case pins, and why each one is a case rather than a comment:
//
//	TestReplyListenerAllowsAnL2CardFromTheNativeSide
//	    交付#1 + 交付#3. An L2 card raised on the assembled bridge is answered
//	    'allow' by the operator's reply stream, the call then RUNS, and the audit
//	    holds both the approval layer's own line (approval: ANSWER-ALLOW) and the
//	    host's (approval: REPLY ANSWERED route="native/allow"). Before this leg
//	    the allow branch of this route had zero production callers: the card could
//	    only die at the C18 deadline.
//	TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel
//	    交付#3's second half. The reason typed next to 'no' is what the model
//	    reads back and what the ledger books - not the fixed
//	    「用户拒绝了本次操作」 fallback.
//	TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject
//	    交付#1's web-side limit, and ticket 201 AC#3's negative control: the very
//	    same decision, arriving on the panel route, is refused on the ROUTE; the
//	    nonce that surfaced there is then unspendable on the native side; and the
//	    route's two real powers (拒绝 / 看卡片) still work.
//	TestUnansweredL2CardTimesOutIntoRejectNeverExecution
//	    交付#2 - the listener deliberately NOT attached, so this is also the
//	    'pull the listener out' control. The terminal state of an unanswered L2 is
//	    a refusal that booked a timeout, and nothing was written.
//	TestL1VetoNeedsAChannelTheHostReallyWired
//	    交付#1's L1 half and 交付#2's boundary. With no channel wired the veto is
//	    refused out loud and the window still expires into execution; with a host
//	    that declares and loads one, the same call stops the write.
//
// The timeout POLARITY of both routes is asserted, never edited here: L2 timeout
// rejects, L1 timeout executes (SPEC-06 §2 row L1). Changing the second one is
// what ticket 201 AC#2 asks for, it is a D4 contract act, and it is reported in
// docs/evidence/s1/201-reply-listener-r1.md §⑥ instead of being done quietly.
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and dies
// at load (0xc0000135) unless third_party/sherpa-onnx is on PATH.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/llm/adaptertest"
	"github.com/CarlosShao/wisp/internal/memory"
	"github.com/CarlosShao/wisp/internal/secret"
)

// syncWriter serialises one run's console writes. The reply listener prints from
// its own goroutine while the loop's sink prints from the task goroutine, and a
// plain bytes.Buffer under both is a data race, not just interleaved text.
type syncWriter struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (w *syncWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.Write(p)
}

func (w *syncWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.String()
}

// replyHost is one assembled `wisp run` with the two directories this ticket
// needs: dir, which [fs] allowlists (an in-allowlist write is the L1 route), and
// outside, which nothing allowlists (R2 => L2).
type replyHost struct {
	t       *testing.T
	srv     *adaptertest.Mockllm
	dir     string
	outside string
	out     *syncWriter
	err     *syncWriter
	// reply is runSpec.reply: nil means this host wired no answer source, which
	// is the pre-201 posture and this suite's own detached-listener control.
	reply io.Reader
	// replyVeto is runSpec.replyVeto: the channel this host's cancel transport
	// claims to be. Production leaves it empty (a console has none); one case
	// below sets it to stand for the ball/Esc-hook leg's host declaration.
	replyVeto approval.Channel
	l2Wait    time.Duration
	rtHook    func(*agentRuntime)
	cleaner   []func()
}

func newReplyHost(t *testing.T, l2Wait time.Duration) *replyHost {
	t.Helper()
	if l2Wait < time.Second {
		// config.toml carries the C18 deadline in whole seconds, and a zero would
		// land on the contract default of 300 - far too long for a test that needs
		// one to expire.
		l2Wait = time.Second
	}
	h := &replyHost{
		t: t, srv: adaptertest.StartMockllm(t),
		out: &syncWriter{}, err: &syncWriter{},
		l2Wait: l2Wait,
	}
	h.dir = t.TempDir()
	h.outside = t.TempDir()
	cfgPath := filepath.Join(h.dir, configFileName)
	body := fmt.Sprintf(`schema_version = 2

[llm]
text_chain = ["acme/m1"]

[llm.retry]
max = 1
backoff_ms = 1

[llm.providers.acme]
protocol = "openai-chat"
base_url = %q
api_key_ref = "dpapi:acme"

[llm.providers.acme.models.m1]
# 261-r2: without an enabled key this entry decodes to enabled=false and the
# llm resolver's gate refuses it, so the assembled run dies before the card
# these cases measure even exists. The cases measure the 201/223/226/255 routes,
# not the gate, so the entry says enabled explicitly.
enabled = true
context_window = 128000

[fs]
allowed_dirs = [%q]

[risk]
confirm_timeout_sec = %d
`, h.srv.Base+"/v1", filepath.ToSlash(h.dir),
		int(h.l2Wait.Seconds()))
	if err := os.WriteFile(cfgPath, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := secret.NewStore(h.dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.Store("dpapi:acme", "wire-fake-key-0123456789abcdef"); err != nil {
		t.Fatal(err)
	}
	// secret.Store holds no handle to release (its reads and writes open the blob
	// per call), so the store below needs no close - unlike memory.Store, which
	// owns the SQLite file and does.
	t.Cleanup(func() {
		for _, fn := range h.cleaner {
			fn()
		}
	})
	return h
}

// run drives the composition root with only the notification poster swapped.
func (h *replyHost) run(task string) int {
	h.t.Helper()
	return runTextTask(runSpec{
		argv:    []string{task},
		stdout:  h.out,
		stderr:  h.err,
		dataDir: h.dir,
		reply:   h.reply,
		// replyVeto is the host's declaration; only the esc case sets one, and the
		// console cases leave it empty exactly as cmdRun does.
		replyVeto: h.replyVeto,
		notify:    func(string, string) error { return nil },
		onRuntime: func(rt *agentRuntime) {
			if h.rtHook != nil {
				h.rtHook(rt)
			}
		},
	})
}

// openStore reads the rows this run booked.
func (h *replyHost) openStore() *memory.Store {
	h.t.Helper()
	st, err := memory.Open(h.dir)
	if err != nil {
		h.t.Fatal(err)
	}
	h.cleaner = append(h.cleaner, func() { _ = st.Close() })
	return st
}

func TestReplyListenerAllowsAnL2CardFromTheNativeSide(t *testing.T) {
	h := newReplyHost(t, 20*time.Second)
	target := filepath.Join(h.outside, "landed-after-a-native-allow.txt")
	pr, pw := io.Pipe()
	h.reply = pr
	t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })

	var (
		out     agent.ToolOutcome
		execErr error
		done    = make(chan struct{})
	)
	h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask("t201-allow")
		defer revoke()
		args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"written because a human said yes"}`,
			filepath.ToSlash(target)))
		go func() {
			defer close(done)
			out, execErr = rt.bridge.Execute(context.Background(), agent.ToolRequest{
				TaskID: "t201-allow", CorrelationID: "t201-allow-corr", CallID: "a1",
				Name: "fs.write", Args: args,
			})
		}()
		if !waitForCard187(rt.liveCards, "t201-allow-corr", 5*time.Second) {
			t.Errorf("the L2 card never reached the native ledger, so the listener had nothing to answer")
			return
		}
		if _, err := fmt.Fprintln(pw, "yes t201-allow-corr"); err != nil {
			t.Errorf("writing the reply: %v", err)
			return
		}
		<-done
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	<-done
	if execErr != nil {
		t.Fatalf("bridge.Execute: %v", execErr)
	}
	// The answer landed, so this one RUNS. That is the point of the ticket: an
	// L2 card批得动的 route, from the native side only.
	if out.IsError {
		t.Fatalf("a natively allowed L2 call must execute, got: %s", out.Text)
	}
	body, err := os.ReadFile(target)
	if err != nil || !strings.Contains(string(body), "written because a human said yes") {
		t.Fatalf("the approved write never landed: %v %q", err, body)
	}
	sentences := h.out.String()
	if !strings.Contains(sentences, "已允许 t201-allow-corr") {
		t.Errorf("the console never told the operator the reply was accepted:\n%s", sentences)
	}
	audit := h.err.String()
	for _, want := range []string{
		"approval: ANSWER-ALLOW corr=t201-allow-corr tool=fs.write route=native decision=allow",
		`approval: REPLY ANSWERED corr="t201-allow-corr" tool=fs.write route="native/allow"`,
	} {
		if !strings.Contains(audit, want) {
			t.Errorf("audit is missing the allow line %q; got:\n%s", want, audit)
		}
	}
	row := toolCallByCorr187(t, h, "t201-allow-corr")
	if row.RiskLevel != memory.RiskL2 || row.Decision != agent.DecisionAllow ||
		row.Outcome != agent.OutcomeSuccess {
		t.Errorf("tool_call row = %+v, want L2/allow/success (the answered-card booking)", row)
	}
}

func TestReplyListenerRejectCarriesTheOperatorsReasonToTheModel(t *testing.T) {
	h := newReplyHost(t, 20*time.Second)
	target := filepath.Join(h.outside, "never-written-after-a-reasoned-reject.txt")
	pr, pw := io.Pipe()
	h.reply = pr
	t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })

	var (
		out     agent.ToolOutcome
		done    = make(chan struct{})
		reason  = "这棵树不在授权目录里，先别写"
		comment = "  它是临时目录"
	)
	h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask("t201-reject")
		defer revoke()
		args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"must not land"}`,
			filepath.ToSlash(target)))
		go func() {
			defer close(done)
			out, _ = rt.bridge.Execute(context.Background(), agent.ToolRequest{
				TaskID: "t201-reject", CorrelationID: "t201-reject-corr", CallID: "r1",
				Name: "fs.write", Args: args,
			})
		}()
		if !waitForCard187(rt.liveCards, "t201-reject-corr", 5*time.Second) {
			t.Error("the L2 card never reached the native ledger")
			return
		}
		// A reply with control characters and a run of spaces in it: the reason
		// is relayed into the model's tool result, so its shape is washed here
		// while its words are not (sanitizeReplyText).
		if _, err := fmt.Fprintln(pw, "no t201-reject-corr "+strings.ReplaceAll(reason+comment, " ", " \t")); err != nil {
			t.Errorf("writing the reply: %v", err)
			return
		}
		<-done
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\n%s\n%s", code, h.out.String(), h.err.String())
	}
	<-done
	if !out.IsError {
		t.Fatalf("a rejected L2 call must not execute, got outcome: %+v", out)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Errorf("the refused write still landed: %v", err)
	}
	// The reason reaches the model: this text IS the tool result the loop feeds
	// back (internal/tools/bridge.go's orDefault(why, ...) on the reject branch).
	if !strings.Contains(out.Text, reason) {
		t.Errorf("the model-visible rejection text lost the operator's reason:\n%s", out.Text)
	}
	if strings.Contains(out.Text, "\t") || strings.Contains(out.Text, "  ") {
		t.Errorf("the reason reached the model unwashed: %q", out.Text)
	}
	audit := h.err.String()
	for _, want := range []string{
		"approval: ANSWER-REJECT corr=t201-reject-corr tool=fs.write decision=reject",
		`approval: REPLY ANSWERED corr="t201-reject-corr" tool=fs.write route="native/reject"`,
	} {
		if !strings.Contains(audit, want) {
			t.Errorf("audit is missing the reject line %q; got:\n%s", want, audit)
		}
	}
	if !strings.Contains(audit, reason) {
		t.Errorf("the audit never booked the reason the operator typed:\n%s", audit)
	}
	row := toolCallByCorr187(t, h, "t201-reject-corr")
	if row.Decision != agent.DecisionReject || row.ErrorClass != "user_rejected" {
		t.Errorf("tool_call row = %+v, want reject/user_rejected", row)
	}
}

func TestPanelRouteRefusesAnAllowBurnsTheGrantAndCanStillReject(t *testing.T) {
	// A longer C18 deadline here: this case answers the same card four times over
	// (the panel's projection, the panel-route allow attempt, the native attempt
	// the burned nonce can no longer support, and the panel route's refusal), and
	// pinning that on a 2s deadline would measure the clock instead of the route.
	h := newReplyHost(t, 20*time.Second)
	target := filepath.Join(h.outside, "panel-route-allow-would-have-written.txt")
	pr, pw := io.Pipe()
	h.reply = pr
	t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })

	var (
		out  agent.ToolOutcome
		done = make(chan struct{})
		corr = "t201-panel-corr"
	)
	h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask("t201-panel")
		defer revoke()
		args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"must never land"}`,
			filepath.ToSlash(target)))
		go func() {
			defer close(done)
			out, _ = rt.bridge.Execute(context.Background(), agent.ToolRequest{
				TaskID: "t201-panel", CorrelationID: corr, CallID: "p1",
				Name: "fs.write", Args: args,
			})
		}()
		if !waitForCard187(rt.liveCards, corr, 8*time.Second) {
			t.Error("the L2 card never reached the native ledger")
			return
		}
		// Each reply is followed by the sentence it has to produce, so the next
		// line can never overtake the one before it in the ledger or on screen.
		// (1) the card is showable through the panel's own projection, which
		// carries no grant and no allow; (2) an ALLOW sent in on the panel route
		// is refused on the route; (3) the native allow that the burned grant can
		// no longer support; (4) the panel route's real power: a refusal that
		// stops the call now.
		for step, pair := range [][2]string{
			{"view " + corr, "此投影不含令牌、不含允许"},
			{"panel-yes " + corr, "面板路线不得允许"},
			{"yes " + corr, "原生令牌无效"},
			{"panel-no " + corr + " 面板侧只能拒绝，这一发我拒了", "已拒绝 " + corr},
		} {
			if _, err := fmt.Fprintln(pw, pair[0]); err != nil {
				t.Errorf("writing %q: %v", pair[0], err)
				return
			}
			waitForSentence187(t, h, pair[1], fmt.Sprintf("step %d (%s)", step, pair[0]))
		}
		<-done
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\n%s\n%s", code, h.out.String(), h.err.String())
	}
	<-done
	if !out.IsError {
		t.Fatalf("a card whose allow arrived on the panel route must never execute: %+v", out)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Errorf("the panel-route allow wrote the file: %v", err)
	}
	sentences, audit := h.out.String(), h.err.String()
	for _, want := range []string{
		// (1) the projection, and what it says it cannot contain.
		"此投影不含令牌、不含允许",
		// (2) the route-level refusal, twice: the approval layer's own line and
		// the host's booking of what it was asked to do.
		"approval: PANEL-ALLOW-REJECTED corr=t201-panel-corr",
		`approval: REPLY PANEL-ALLOW-REFUSED corr="t201-panel-corr"`,
		"面板路线不得允许",
		// (3) the grant that surfaced on the untrusted route cannot be spent here.
		"approval: FORGED-OR-STALE allow rejected corr=t201-panel-corr",
		// (4) the panel route's refusal works, with the reason attached.
		"approval: ANSWER-REJECT corr=t201-panel-corr",
		`approval: REPLY ANSWERED corr="t201-panel-corr" tool=fs.write route="panel/reject"`,
	} {
		if !strings.Contains(audit, want) && !strings.Contains(sentences, want) {
			t.Errorf("neither the console nor the audit says %q\nstdout:\n%s\nstderr:\n%s",
				want, sentences, audit)
		}
	}
	row := toolCallByCorr187(t, h, corr)
	if row.Decision != agent.DecisionReject {
		t.Errorf("tool_call row = %+v, want the panel refusal booked as reject", row)
	}
}

func TestUnansweredL2CardTimesOutIntoRejectNeverExecution(t *testing.T) {
	// 交付#2's criterion, and the listener-detached control in the same breath:
	// h.reply stays nil, so this run has the answer side assembled but UNREACHABLE.
	// The card's terminal state must be a refusal that booked a timeout - never
	// execution - and the audit must say which of the two clocks ran out.
	h := newReplyHost(t, 2*time.Second)
	target := filepath.Join(h.outside, "unanswered-card-must-not-write.txt")
	var (
		out  agent.ToolOutcome
		done = make(chan struct{})
	)
	h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask("t201-timeout")
		defer revoke()
		args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"must not land on a timeout"}`,
			filepath.ToSlash(target)))
		go func() {
			defer close(done)
			out, _ = rt.bridge.Execute(context.Background(), agent.ToolRequest{
				TaskID: "t201-timeout", CorrelationID: "t201-timeout-corr", CallID: "x1",
				Name: "fs.write", Args: args,
			})
		}()
		if !waitForCard187(rt.liveCards, "t201-timeout-corr", 5*time.Second) {
			t.Error("the L2 card never reached the native ledger")
			return
		}
		if rt.reply != nil {
			t.Errorf("this case is the detached control, but a reply surface was attached anyway")
		}
		// Wait for the call to resolve INSIDE the hook: the run closes its store
		// when the task ends, and a row booked after that would be written to a
		// closed database. This is the timeout path, so the wait is the C18 clock.
		<-done
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\n%s\n%s", code, h.out.String(), h.err.String())
	}
	<-done
	if !out.IsError {
		t.Fatalf("an unanswered L2 card executed: %+v", out)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Errorf("the unanswered card wrote the file: %v", err)
	}
	if strings.Contains(out.Text, "用户拒绝了本次操作") {
		t.Errorf("a timeout must not read like a human's refusal: %q", out.Text)
	}
	audit := h.err.String()
	for _, want := range []string{
		"approval: ANSWER-EXPIRED corr=t201-timeout-corr tool=fs.write decision=timeout->reject",
	} {
		if !strings.Contains(audit, want) {
			t.Errorf("audit is missing the expiry line %q; got:\n%s", want, audit)
		}
	}
	row := toolCallByCorr187(t, h, "t201-timeout-corr")
	if row.Decision != agent.DecisionTimeout {
		t.Errorf("tool_call row = %+v, want decision=%s (the C18 auto-reject, not an allow)",
			row, agent.DecisionTimeout)
	}
}

func TestL1VetoNeedsAChannelTheHostReallyWired(t *testing.T) {
	t.Run("console posture: veto refused, window still executes", func(t *testing.T) {
		// The console leg wires none of SPEC-06 §2's four veto channels, so the
		// veto arrives at Gate.Veto and is refused there - loudly, audited - and
		// the L1 window keeps its frozen polarity: it expires into EXECUTION
		// (SPEC-06 §2 row L1). This half is the evidence that ticket 201 did NOT
		// quietly change that, and that the run says so instead of going silent.
		h := newReplyHost(t, 20*time.Second)
		target := filepath.Join(h.dir, "l1-window-console-veto.txt")
		pr, pw := io.Pipe()
		h.reply = pr
		t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })

		var (
			out  agent.ToolOutcome
			done = make(chan struct{})
			corr = "t201-l1-console-corr"
		)
		h.rtHook = func(rt *agentRuntime) {
			revoke := rt.gate.AdmitTextTask("t201-l1-console")
			defer revoke()
			args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"landed through the window, unanswered"}`,
				filepath.ToSlash(target)))
			go func() {
				defer close(done)
				out, _ = rt.bridge.Execute(context.Background(), agent.ToolRequest{
					TaskID: "t201-l1-console", CorrelationID: corr, CallID: "v1",
					Name: "fs.write", Args: args,
				})
			}()
			if !waitForCard187(rt.liveCards, corr, 5*time.Second) {
				t.Error("the L1 window never reached the native ledger")
				return
			}
			if ch := rt.reply.vetoChannel; ch != "" {
				t.Errorf("a console run must not declare a veto channel, got %q", ch)
			}
			if _, err := fmt.Fprintln(pw, "veto "+corr); err != nil {
				t.Errorf("writing the veto: %v", err)
				return
			}
			<-done
		}
		if code := h.run("总结一下 这份笔记"); code != 0 {
			t.Fatalf("exit %d\n%s\n%s", code, h.out.String(), h.err.String())
		}
		<-done
		if out.IsError {
			t.Fatalf("the L1 timeout polarity changed (this is the D4 contract, not this leg): %s", out.Text)
		}
		if body, err := os.ReadFile(target); err != nil || !strings.Contains(string(body), "unanswered") {
			t.Fatalf("the unopposed L1 window did not execute as SPEC-06 §2 says: %v %q", err, body)
		}
		audit := h.err.String()
		// The veto that never landed must NOT be booked as one that did: this is
		// the assertion that keeps the listener from writing a happy answer line
		// when the gate refused to take the cancel.
		if strings.Contains(audit, "approval: ANSWER-VETO corr="+corr) {
			t.Errorf("a veto that never landed was booked as one that did:\n%s", audit)
		}
		for _, want := range []string{
			`approval: REPLY REFUSED corr="` + corr + `" tool=fs.write route="native/veto"`,
			"approval: ANSWER-EXPIRED corr=" + corr + " tool=fs.write decision=timeout->execute",
		} {
			if !strings.Contains(audit, want) {
				t.Errorf("audit is missing %q; got:\n%s", want, audit)
			}
		}
		if !strings.Contains(h.out.String(), "否决未能送达") {
			t.Errorf("the operator was never told the veto could not be delivered:\n%s", h.out.String())
		}
	})

	t.Run("host declares and loads esc, veto lands", func(t *testing.T) {
		// The same production code path, with the two halves a GUI host owns
		// arriving through the ONE door a host has: runSpec.replyVeto. The
		// assembly then marks that channel loaded (approval_reply.go), because a
		// host that names a transport is asserting it is up - and this is the
		// seam standing in for the ball/Esc-hook leg (tickets 07/77/92), stated as
		// such in the evidence file. The reading it produces is the one the ticket
		// wants: a veto that arrives in time STOPS the write.
		h := newReplyHost(t, 20*time.Second)
		h.replyVeto = approval.ChannelEsc
		target := filepath.Join(h.dir, "l1-window-vetoed-by-a-loaded-channel.txt")
		pr, pw := io.Pipe()
		h.reply = pr
		t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })

		var (
			out  agent.ToolOutcome
			done = make(chan struct{})
			corr = "t201-l1-veto-corr"
		)
		h.rtHook = func(rt *agentRuntime) {
			revoke := rt.gate.AdmitTextTask("t201-l1-veto")
			defer revoke()
			if rt.reply == nil {
				t.Fatal("no reply surface, so nothing could carry the veto")
			}
			if got := rt.reply.vetoChannel; got != approval.ChannelEsc {
				t.Errorf("the host declared %q, the assembled surface carries %q",
					approval.ChannelEsc, got)
			}
			if !rt.gate.Channels().Loaded(approval.ChannelEsc) {
				t.Error("declaring a cancel transport failed to mark that veto channel loaded")
			}
			args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"must never land"}`,
				filepath.ToSlash(target)))
			go func() {
				defer close(done)
				out, _ = rt.bridge.Execute(context.Background(), agent.ToolRequest{
					TaskID: "t201-l1-veto", CorrelationID: corr, CallID: "v2",
					Name: "fs.write", Args: args,
				})
			}()
			if !waitForCard187(rt.liveCards, corr, 5*time.Second) {
				t.Error("the L1 window never reached the native ledger")
				return
			}
			if _, err := fmt.Fprintln(pw, "veto "+corr); err != nil {
				t.Errorf("writing the veto: %v", err)
				return
			}
			<-done
		}
		if code := h.run("总结一下 这份笔记"); code != 0 {
			t.Fatalf("exit %d\n%s\n%s", code, h.out.String(), h.err.String())
		}
		<-done
		if !out.IsError {
			t.Fatalf("a veto that landed must stop the call: %+v", out)
		}
		if _, err := os.Stat(target); !os.IsNotExist(err) {
			t.Errorf("the vetoed L1 window wrote the file anyway: %v", err)
		}
		if !strings.Contains(out.Text, "否决") {
			t.Errorf("the model-visible text lost the veto: %q", out.Text)
		}
		audit := h.err.String()
		if !strings.Contains(audit, "approval: ANSWER-VETO corr="+corr+" tool=fs.write channel=esc decision=veto") {
			t.Errorf("audit is missing the veto line; got:\n%s", audit)
		}
		if strings.Contains(audit, "approval: ANSWER-EXPIRED corr="+corr) {
			t.Errorf("a vetoed window must not also book an expiry:\n%s", audit)
		}
		row := toolCallByCorr187(t, h, corr)
		if row.Decision != agent.DecisionReject {
			t.Errorf("tool_call row = %+v, want the veto booked as reject", row)
		}
	})
}

// waitForCard187 polls the native ledger until the named card is displayed and
// answerable. Polling, not sleeping: the reply has to arrive after the gate
// armed the item and before its own clock runs out, and the ledger write is the
// instant both of those are true.
func waitForCard187(live *nativeCards, corr string, limit time.Duration) bool {
	deadline := time.Now().Add(limit)
	for time.Now().Before(deadline) {
		if _, ok := live.look(corr); ok {
			return true
		}
		time.Sleep(2 * time.Millisecond)
	}
	return false
}

// waitForSentence187 parks the writing test until the reply it just fed has
// produced its sentence, on either stream. Without this the four answers of the
// panel-route case would race into the listener's scanner and the ledger would
// hold them in an order the case cannot describe.
func waitForSentence187(t *testing.T, h *replyHost, needle, why string) {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		if strings.Contains(h.out.String(), needle) || strings.Contains(h.err.String(), needle) {
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Errorf("no stream ever said %q after %q\nstdout:\n%s\nstderr:\n%s",
		needle, why, h.out.String(), h.err.String())
}

// toolCallByCorr187 reads back the one tool_call row this case dispatched.
func toolCallByCorr187(t *testing.T, h *replyHost, corr string) memory.ToolCall {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := h.openStore().ListToolCallsByTask(ctx, strings.TrimSuffix(corr, "-corr"))
	if err != nil {
		t.Fatal(err)
	}
	// task ids here end in the same stem as the correlation id minus its suffix,
	// which the callers above keep aligned by hand; walk the rows instead of
	// trusting one id.
	for _, r := range rows {
		if r.CorrelationID == corr {
			return r
		}
	}
	all := make([]string, 0, len(rows))
	for _, r := range rows {
		all = append(all, fmt.Sprintf("%+v", r))
	}
	t.Fatalf("no tool_call row carries correlation id %q; rows for that task:\n%s",
		corr, strings.Join(all, "\n"))
	return memory.ToolCall{}
}
