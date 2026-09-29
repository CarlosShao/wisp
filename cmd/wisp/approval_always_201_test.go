package main

// Ticket 201 AC#4 - the 长期 branch shows the rule it would store, and the
// widening itself is a second, L2-level confirmation.
//
// PLAN.md:1645-1646 freezes the sentence these two cases are built on: 热加载放宽
// [risk]/[fs]/[net]/[plugins]「必须触发 L2 级重新确认，不得静默生效」. So an
// 「一直」 answer is not one card doing two jobs - the card that lets THIS call run
// is about the call, and the line that lands in [fs] allowed_dirs governs every
// future one. The two cases differ on exactly one thing, the second card's answer:
//
//	TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card   the rule text is printed
//	    before anything happens, a second L2 card is raised about that exact line,
//	    answering it 'yes' persists the line into config.toml, and the original
//	    card's own answer stays independent of it.
//	TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused   the fail-closed
//	    control: the second card is answered 'no' (and, in the same case, the
//	    original card too) - config.toml gains no line, and the widening is booked
//	    as a refusal rather than silence.
//
// Both drive the assembled `wisp run` with a scripted reply stream, which is the
// injection seam AGENTS.md §1.3 names for the CLI - no mock stands in for the gate,
// the queue, the config writer or the file.
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and dies
// at load (0xc0000135) unless third_party/sherpa-onnx is on PATH.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
)

// widenRun is one always-branch run: the host, the in-flight tool call it drives,
// and the state the reply driver and the assertions share.
type widenRun struct {
	h      *replyHost
	done   chan struct{}
	reply  chan struct{}
	corr   string
	target string
	mu     sync.Mutex
	out    agent.ToolOutcome
	exec   error
	fail   error
	rtRef  *agentRuntime
}

func (w *widenRun) setFail(err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.fail == nil && err != nil {
		w.fail = err
	}
}

func (w *widenRun) failure() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.fail
}

func (w *widenRun) runtime() *agentRuntime {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.rtRef
}

// newWidenRun builds the host plus one outside-the-allowlist fs.write, whose L2
// card is what the 长期 verb is typed against.
func newWidenRun(t *testing.T, l2Wait time.Duration, taskID, corr string) (*widenRun, *io.PipeWriter) {
	t.Helper()
	w := &widenRun{
		h:     newReplyHost(t, l2Wait),
		done:  make(chan struct{}),
		reply: make(chan struct{}),
		corr:  corr,
	}
	pr, pw := io.Pipe()
	w.h.reply = pr
	t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })

	target := filepath.Join(w.h.outside, strings.TrimSuffix(corr, "-corr")+".txt")
	w.h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask(taskID)
		defer revoke()
		args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"long-lived answer"}`,
			filepath.ToSlash(target)))
		go func() {
			defer close(w.done)
			out, err := rt.bridge.Execute(context.Background(), agent.ToolRequest{
				TaskID: taskID, CorrelationID: corr, CallID: "a1",
				Name: "fs.write", Args: args,
			})
			w.mu.Lock()
			w.out, w.exec = out, err
			w.mu.Unlock()
		}()
		if !waitForCard187(rt.liveCards, corr, 5*time.Second) {
			w.setFail(fmt.Errorf("the L2 card for %s never reached the ledger", corr))
			return
		}
		w.mu.Lock()
		w.rtRef = rt
		w.mu.Unlock()
	}
	w.target = target
	return w, pw
}

func TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card(t *testing.T) {
	w, pw := newWidenRun(t, 40*time.Second, "t201-always", "t201-always-corr")

	go driveAlways(t, w, pw, true)

	if code := w.h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, w.h.out.String(), w.h.err.String())
	}
	<-w.done
	<-w.reply
	if err := w.failure(); err != nil {
		t.Fatal(err)
	}

	text, err := os.ReadFile(filepath.Join(w.h.dir, configFileName))
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.ToSlash(w.h.outside)
	if !strings.Contains(string(text), want) {
		t.Fatalf("config.toml did not gain the stored rule %q; it holds:\n%s", want, text)
	}
	audit := w.h.err.String()
	for _, needle := range []string{
		`approval: REPLY STAGED corr="t201-always-corr"`,
		"approval: REPLY WIDEN-APPLIED",
	} {
		if !strings.Contains(audit, needle) {
			t.Errorf("audit is missing %q; got:\n%s", needle, audit)
		}
	}
	w.mu.Lock()
	out, execErr := w.out, w.exec
	w.mu.Unlock()
	if execErr != nil {
		t.Fatalf("bridge.Execute: %v", execErr)
	}
	if out.IsError {
		t.Fatalf("the first card was answered yes, so the write must have run: %s", out.Text)
	}
}

func TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused(t *testing.T) {
	w, pw := newWidenRun(t, 40*time.Second, "t201-always2", "t201-always2-corr")

	go driveAlways(t, w, pw, false)

	if code := w.h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, w.h.out.String(), w.h.err.String())
	}
	<-w.done
	<-w.reply
	if err := w.failure(); err != nil {
		t.Fatal(err)
	}

	text, err := os.ReadFile(filepath.Join(w.h.dir, configFileName))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(text), filepath.ToSlash(w.h.outside)) {
		t.Fatalf("config.toml gained an allowlist entry the second card refused (PLAN.md:1645-1646: 不得静默生效):\n%s", text)
	}
	if _, err := os.Stat(w.target); err == nil {
		t.Fatalf("the refused call still wrote: %s", w.target)
	}
	if !strings.Contains(w.h.err.String(), "approval: REPLY WIDEN-REFUSED") {
		t.Errorf("a refused widening was not booked as a refusal:\n%s", w.h.err.String())
	}
}

// driveAlways types the 长期 verb, waits for the rule text and the second card,
// answers that card (yes when persist is true, no otherwise), and only then
// answers the original card - so the run cannot finish out from under the
// widening worker.
func driveAlways(t *testing.T, w *widenRun, pw *io.PipeWriter, persist bool) {
	t.Helper()
	defer func() { _ = pw.Close() }()
	defer close(w.reply)

	waitForSentence187(t, w.h, "卡片编号：", "the first card must show the address the operator answers against")
	if _, err := fmt.Fprintln(pw, "always "+w.corr); err != nil {
		w.setFail(err)
		return
	}
	// AC#4's first demand: the rule text is on the screen before anything is
	// stored, and it is the line itself, not a paraphrase of it.
	waitForSentence187(t, w.h, "要存的规则：", "the 长期 answer must print the rule it stores")
	waitForSentence187(t, w.h, `allowed_dirs += "`, "the printed rule must name the key and the directory")

	corr := waitWidenCard187(t, w, 15*time.Second)
	if corr == "" {
		w.setFail(errors.New("the second, L2-level re-confirmation card never appeared"))
		return
	}
	// The widening card must itself say what it stores (the reason the model and
	// the operator read is the same line the file will hold).
	waitForSentence187(t, w.h, "这是「一直」要求的 L2 级重新确认", "the second card must say it is the re-confirmation")

	answer := "no "
	if persist {
		answer = "yes "
	}
	if _, err := fmt.Fprintln(pw, answer+corr); err != nil {
		w.setFail(err)
		return
	}
	needle := "approval: REPLY WIDEN-APPLIED"
	if !persist {
		needle = "approval: REPLY WIDEN-REFUSED"
	}
	if !waitForAudit187(w.h, needle, 25*time.Second) {
		w.setFail(fmt.Errorf("the widening never booked %q", needle))
		return
	}
	final := "no " + w.corr + " 这一次不干活"
	if persist {
		// The two answers stay independent: the second card decided the RULE,
		// this one decides only this call.
		final = "yes " + w.corr
	}
	if _, err := fmt.Fprintln(pw, final); err != nil {
		w.setFail(err)
		return
	}
}

// waitWidenCard187 returns the correlation id of the widening card this run
// raised (the one whose tool is the config door), or "" after the limit.
func waitWidenCard187(t *testing.T, w *widenRun, limit time.Duration) string {
	t.Helper()
	deadline := time.Now().Add(limit)
	for time.Now().Before(deadline) {
		if rt := w.runtime(); rt != nil {
			for _, c := range rt.liveCards.pending() {
				if c.Tool == allowWidenTool {
					return c.CorrelationID
				}
			}
		}
		time.Sleep(2 * time.Millisecond)
	}
	return ""
}

// waitForAudit187 waits for one line to reach the run's audit stream.
func waitForAudit187(h *replyHost, needle string, limit time.Duration) bool {
	deadline := time.Now().Add(limit)
	for time.Now().Before(deadline) {
		if strings.Contains(h.err.String(), needle) {
			return true
		}
		time.Sleep(2 * time.Millisecond)
	}
	return false
}
