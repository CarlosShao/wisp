package main

// Ticket 201 AC#1/AC#3/AC#6 - the host-facing seam, judged WITHOUT a console.
//
// These two cases are what the seam exists for. The five cases in
// approval_reply_201_test.go all answer a card by typing a line into the reply
// listener, which is the one shape of "someone answered" a terminal can produce;
// a ball click handler or a tray menu item is not a terminal. So both cases below
// run with h.reply = nil - attachReplyListener is never called, no listener
// goroutine exists, no verb is ever typed - and answer the card through
// approval.Replies, the exported entry point a host is handed
// (internal/agent/approval/replies.go). If that file's methods were only
// reachable from the console grammar, these two cases could not be written, and
// "有入口" would again mean "有一个测试能碰到它" - the reading ticket 219's AC#7
// lesson refuses.
//
//	TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole   AC#1 + AC#3's allow side
//	TestNativeHostSeamRefusesAPanelSourcedAllow         AC#3's negative control,
//	    through the seam rather than through the console: the panel route cannot
//	    allow, the nonce that surfaced there is burned, and the route's real power
//	    (拒绝) still works.
//
// Both also read AC#6 off the same ledger: the waiting state has to be a fact
// produced by the running process at the instant a card exists, not a name a
// display file invented. WaitingState returns only D43 names, and the audit line
// booking it comes from cmd/wisp/run.go's own Prompt path.
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and dies
// at load (0xc0000135) unless third_party/sherpa-onnx is on PATH.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole is AC#1's criterion: an
// allow that arrives through the exported seam, with no reply listener attached.
func TestNativeHostSeamAnswersAnL2CardWithoutAnyConsole(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	target := filepath.Join(h.outside, "landed-through-the-seam.txt")
	// h.reply stays nil: no console, no verbs, no listener goroutine.

	var (
		out     agent.ToolOutcome
		execErr error
		done    = make(chan struct{})
		seamErr error
	)
	h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask("t201-seam")
		defer revoke()
		args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"written by a native click, not a keystroke"}`,
			filepath.ToSlash(target)))
		go func() {
			defer close(done)
			out, execErr = rt.bridge.Execute(context.Background(), agent.ToolRequest{
				TaskID: "t201-seam", CorrelationID: "t201-seam-corr", CallID: "a1",
				Name: "fs.write", Args: args,
			})
		}()
		if !waitForCard187(rt.liveCards, "t201-seam-corr", 5*time.Second) {
			seamErr = errors.New("the card never reached the host ledger")
			return
		}

		// AC#6, at the instant the card is really pending.
		if st, awaiting := rt.waitingStateName(); !awaiting || st != statemachine.StateAwaitingApproval {
			seamErr = fmt.Errorf("waitingStateName() = (%q,%v), want (AwaitingApproval,true) while the card is pending", st, awaiting)
			return
		}
		if card, ok := rt.liveCards.h.AwaitingHuman(); !ok || card.CorrelationID != "t201-seam-corr" {
			seamErr = fmt.Errorf("AwaitingHuman() = (%+v,%v), want the pending card", card, ok)
			return
		}

		// AC#1: the host answers through the exported seam. Nothing in this line
		// goes through the console grammar.
		if err := rt.liveCards.h.Allow(context.Background(), "t201-seam-corr"); err != nil {
			seamErr = fmt.Errorf("the seam refused a native allow it must accept: %w", err)
			return
		}
		// The wait is over: the same reading must flip.
		if st, awaiting := rt.waitingStateName(); awaiting {
			seamErr = fmt.Errorf("the waiting state never cleared; got %q still awaiting", st)
			return
		}
		<-done
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	<-done
	if seamErr != nil {
		t.Fatal(seamErr)
	}
	if execErr != nil {
		t.Fatalf("bridge.Execute: %v", execErr)
	}
	if out.IsError {
		t.Fatalf("a call the host seam allowed must execute, got: %s", out.Text)
	}
	body, err := os.ReadFile(target)
	if err != nil || !strings.Contains(string(body), "written by a native click") {
		t.Fatalf("the seam's allow never reached the write: %v %q", err, body)
	}
	audit := h.err.String()
	for _, want := range []string{
		"approval: ANSWER-ALLOW corr=t201-seam-corr tool=fs.write route=native decision=allow",
		`approval: WAITING-STATE state="AwaitingApproval" awaiting=true corr="t201-seam-corr"`,
	} {
		if !strings.Contains(audit, want) {
			t.Errorf("audit is missing %q; got:\n%s", want, audit)
		}
	}
	// And it stayed a two-event story: the console listener was never attached,
	// so no REPLY line of the host grammar can exist.
	if strings.Contains(audit, "approval: REPLY ") {
		t.Errorf("a host answer arrived without any console, yet the console's own reply ledger booked a line:\n%s", audit)
	}
}

// TestNativeHostSeamRefusesAPanelSourcedAllow is AC#3's negative control on the
// seam itself: AGENTS.md §1.2's ban (#6, no panel-sourced L2 allow) has to hold
// for the entry point a WebView host would be handed, not only for the terminal's
// panel-yes verb.
func TestNativeHostSeamRefusesAPanelSourcedAllow(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	target := filepath.Join(h.outside, "must-never-exist.txt")

	var (
		out     agent.ToolOutcome
		execErr error
		done    = make(chan struct{})
		fail    error
		seen    *agentRuntime
	)
	h.rtHook = func(rt *agentRuntime) {
		seen = rt
		revoke := rt.gate.AdmitTextTask("t201-seam-panel")
		defer revoke()
		args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"only a page could have asked for this"}`,
			filepath.ToSlash(target)))
		go func() {
			defer close(done)
			out, execErr = rt.bridge.Execute(context.Background(), agent.ToolRequest{
				TaskID: "t201-seam-panel", CorrelationID: "t201-seam-panel-corr", CallID: "a1",
				Name: "fs.write", Args: args,
			})
		}()
		if !waitForCard187(rt.liveCards, "t201-seam-panel-corr", 5*time.Second) {
			fail = errors.New("the card never reached the host ledger")
			return
		}

		// (1) The panel route's allow is refused on the ROUTE, and the seam says
		//     out loud that a native grant was offered on it (burned as leaked).
		offered, err := rt.liveCards.h.PanelAllow(context.Background(), "t201-seam-panel-corr")
		if !errors.Is(err, approval.ErrPanelAllow) {
			fail = fmt.Errorf("PanelAllow must be refused on the route, got (%v,%v)", offered, err)
			return
		}
		if !offered {
			fail = fmt.Errorf("PanelAllow must report that a native grant was offered on the untrusted route: %v", err)
			return
		}
		// (2) The burned nonce cannot support a native allow either.
		if err := rt.liveCards.h.Allow(context.Background(), "t201-seam-panel-corr"); !errors.Is(err, approval.ErrBadGrant) {
			fail = fmt.Errorf("a grant that surfaced on the panel route must be unspendable natively, got %v", err)
			return
		}
		// (3) The route's real power still works: a refusal.
		if err := rt.liveCards.h.PanelReject(context.Background(), "t201-seam-panel-corr", "页面只能拒绝"); err != nil {
			fail = fmt.Errorf("the panel route's refusal must still land: %w", err)
			return
		}
		<-done
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	<-done
	if fail != nil {
		t.Fatal(fail)
	}
	if execErr != nil {
		t.Fatalf("bridge.Execute: %v", execErr)
	}
	if !out.IsError {
		t.Fatalf("a card whose allow arrived on the panel route must never execute: %+v", out)
	}
	if _, err := os.Stat(target); err == nil {
		t.Fatalf("the panel route wrote a file: %s", target)
	}
	if seen != nil {
		if _, awaiting := seen.waitingStateName(); awaiting {
			t.Errorf("the waiting state outlived the settled card")
		}
	}
}
