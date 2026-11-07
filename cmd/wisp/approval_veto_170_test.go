package main

// Ticket 170 - the CLI leg's own veto channel (form 甲, owner-approved): when
// stdin is a TTY the console declares a cancel transport, and any input that
// reaches the reply listener while an L1 block window is live cancels it.
//
// What each case pins, and why each one is a case rather than a comment:
//
//	Test170AnyInputStopsTheL1Window
//	    the new shape: input that is NOT a reply verb ("停一下") cancels the
//	    live window - the operator does not need a grammar in a 2-second
//	    window, they need to say stop.
//	Test170BareEnterStopsTheL1Window
//	    the key the card actually names on this leg (回车): an empty line is
//	    still input while the window is live, and it must land the same veto.
//	Test170WithoutTheDeclaredChannelTheSameInputExecutes
//	    AC#5's reverse reading: the SAME input, with the host declaration
//	    pulled out, is refused by the gate's own registry and the window
//	    expires into EXECUTION (SPEC-06 §2 row L1) - the channel is a real
//	    transport, not a decoration.
//	Test170TTYEntryDeclaresTheCancelKeySlotAndNamesTheRealKey
//	    the entry's two statements: cmdRun's TTY declaration yields the esc
//	    slot, and the card names the key this host really holds (回车), never
//	    a phantom global Esc hook (B1, ticket 260's spelling seam).
//	Test170NoDeclaredChannelBooksTheLoudAuditLine
//	    the non-TTY posture is SAID: run.go books CLI-VETO-NONE instead of
//	    leaving the silence as an implied decision.
//
// The timeout POLARITY of the L1 route is asserted, never edited here: an
// unopposed window still executes (SPEC-06 §2 row L1), which is exactly what
// Test170WithoutTheDeclaredChannelTheSameInputExecutes re-proves after the
// veto attempt is refused.
//
// PATH warning (ticket 98): this package's test binary links sherpa-onnx and
// dies at load unless third_party/sherpa-onnx is on PATH.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
)

func Test170AnyInputStopsTheL1Window(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	h.replyVeto = approval.ChannelEsc
	target := filepath.Join(h.dir, "t170-any-input-vetoed.txt")
	pr, pw := io.Pipe()
	h.reply = pr
	t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })

	var (
		out  agent.ToolOutcome
		done = make(chan struct{})
		corr = "t170-any-input-corr"
	)
	h.rtHook = func(rt *agentRuntime) {
		if rt.reply == nil {
			t.Fatal("no reply surface, so nothing could carry the veto")
		}
		if got := rt.reply.vetoChannel; got != approval.ChannelEsc {
			t.Fatalf("the host declared esc, the assembled surface carries %q", got)
		}
		revoke := rt.gate.AdmitTextTask("t170-any-input")
		defer revoke()
		args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"must never land"}`,
			filepath.ToSlash(target)))
		go func() {
			defer close(done)
			out, _ = rt.bridge.Execute(context.Background(), agent.ToolRequest{
				TaskID: "t170-any-input", CorrelationID: corr, CallID: "a1",
				Name: "fs.write", Args: args,
			})
		}()
		if !waitForCard187(rt.liveCards, corr, 5*time.Second) {
			t.Error("the L1 window never reached the native ledger")
			return
		}
		// NOT a reply verb: in the old grammar this line answered 未知答复指令
		// and the window kept counting. Under ticket 170 the window is live, so
		// the line IS the 停.
		if _, err := fmt.Fprintln(pw, "停一下"); err != nil {
			t.Errorf("writing the stop line: %v", err)
			return
		}
		<-done
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	<-done
	if !out.IsError {
		t.Fatalf("any input during a live L1 window must cancel the call: %+v", out)
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
	if !strings.Contains(h.out.String(), "已否决 "+corr) {
		t.Errorf("the console never told the operator the stop landed:\n%s", h.out.String())
	}
}

func Test170BareEnterStopsTheL1Window(t *testing.T) {
	// The card on this leg names 回车 as the cancel key (main.go's spelling
	// seam), so a bare Enter - an EMPTY line - has to land the same veto, or
	// the card's sentence would be the only thing about it that works.
	h := newReplyHost(t, 40*time.Second)
	h.replyVeto = approval.ChannelEsc
	target := filepath.Join(h.dir, "t170-bare-enter-vetoed.txt")
	pr, pw := io.Pipe()
	h.reply = pr
	t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })

	var (
		out  agent.ToolOutcome
		done = make(chan struct{})
		corr = "t170-bare-enter-corr"
	)
	h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask("t170-bare-enter")
		defer revoke()
		args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"must never land"}`,
			filepath.ToSlash(target)))
		go func() {
			defer close(done)
			out, _ = rt.bridge.Execute(context.Background(), agent.ToolRequest{
				TaskID: "t170-bare-enter", CorrelationID: corr, CallID: "e1",
				Name: "fs.write", Args: args,
			})
		}()
		if !waitForCard187(rt.liveCards, corr, 5*time.Second) {
			t.Error("the L1 window never reached the native ledger")
			return
		}
		if _, err := fmt.Fprintln(pw, ""); err != nil {
			t.Errorf("writing the bare Enter: %v", err)
			return
		}
		<-done
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	<-done
	if !out.IsError {
		t.Fatalf("a bare Enter during a live L1 window must cancel the call: %+v", out)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Errorf("the vetoed L1 window wrote the file anyway: %v", err)
	}
	if !strings.Contains(h.err.String(),
		"approval: ANSWER-VETO corr="+corr+" tool=fs.write channel=esc decision=veto") {
		t.Errorf("audit is missing the veto line; got:\n%s", h.err.String())
	}
}

func Test170WithoutTheDeclaredChannelTheSameInputExecutes(t *testing.T) {
	// AC#5's reverse reading, run against the assembled stack: the SAME stop
	// line, with the host declaration pulled out (replyVeto empty - the exact
	// value cmdRun fills for a non-TTY stdin), is refused by the gate's own
	// registry and the window keeps its frozen polarity: it expires into
	// EXECUTION and the write lands. If pulling the transport out changed
	// nothing here, the transport would be a decoration.
	h := newReplyHost(t, 40*time.Second)
	target := filepath.Join(h.dir, "t170-undeclared-executes.txt")
	pr, pw := io.Pipe()
	h.reply = pr
	t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })

	var (
		out  agent.ToolOutcome
		done = make(chan struct{})
		corr = "t170-undeclared-corr"
	)
	h.rtHook = func(rt *agentRuntime) {
		if got := rt.reply.vetoChannel; got != "" {
			t.Fatalf("this case is the undeclared control, but the surface carries %q", got)
		}
		if rt.gate.Channels().Loaded(approval.ChannelEsc) {
			t.Error("no host declared a cancel transport, yet esc reads as loaded")
		}
		revoke := rt.gate.AdmitTextTask("t170-undeclared")
		defer revoke()
		args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"written because nobody said stop"}`,
			filepath.ToSlash(target)))
		go func() {
			defer close(done)
			out, _ = rt.bridge.Execute(context.Background(), agent.ToolRequest{
				TaskID: "t170-undeclared", CorrelationID: corr, CallID: "u1",
				Name: "fs.write", Args: args,
			})
		}()
		if !waitForCard187(rt.liveCards, corr, 5*time.Second) {
			t.Error("the L1 window never reached the native ledger")
			return
		}
		if _, err := fmt.Fprintln(pw, "停一下"); err != nil {
			t.Errorf("writing the stop line: %v", err)
			return
		}
		<-done
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	<-done
	if out.IsError {
		t.Fatalf("with no declared channel the window's polarity is unchanged: %s", out.Text)
	}
	body, err := os.ReadFile(target)
	if err != nil || !strings.Contains(string(body), "nobody said stop") {
		t.Fatalf("the unopposed L1 window did not execute as SPEC-06 §2 says: %v %q", err, body)
	}
	audit := h.err.String()
	for _, want := range []string{
		"approval: ANSWER-EXPIRED corr=" + corr + " tool=fs.write decision=timeout->execute",
		`approval: REPLY REFUSED corr="` + corr + `" tool=fs.write route="native/veto"`,
		// the loud non-declared posture line, with a reply source attached
		"approval: CLI-VETO-NONE reply_attached=true",
	} {
		if !strings.Contains(audit, want) {
			t.Errorf("audit is missing %q; got:\n%s", want, audit)
		}
	}
	if strings.Contains(audit, "approval: ANSWER-VETO corr="+corr) {
		t.Errorf("a veto that never landed was booked as one that did:\n%s", audit)
	}
	if !strings.Contains(h.out.String(), "答复未被接受（L1 窗口否决）") {
		t.Errorf("the operator was never told the stop could not be delivered:\n%s", h.out.String())
	}
}

func Test170TTYEntryDeclaresTheCancelKeySlotAndNamesTheRealKey(t *testing.T) {
	// The entry's two statements together: the slot is one of the closed four
	// (esc = the cancel-KEY channel - the set is frozen, a fifth selector is
	// the M-7/C-3 shape), and the card names the key this host really holds.
	prev := approval.SetCancelKeySpelling(nil)
	t.Cleanup(func() { approval.SetCancelKeySpelling(prev) })
	if ch := declareConsoleCancelTransport(); ch != approval.ChannelEsc {
		t.Fatalf("the TTY entry declared %q, want the cancel-key slot esc", ch)
	}
	reg := approval.NewChannels(approval.ChannelEsc)
	p := approval.Prompt{
		CorrelationID: "t170-label-corr",
		Level:         "L1",
		Window:        2 * time.Second,
		Channels:      reg.Statuses(),
	}
	card := strings.Join(approval.RenderCard(p), "\n")
	if !strings.Contains(card, "回车") {
		t.Errorf("the card does not name the key this host really holds:\n%s", card)
	}
	if strings.Contains(card, "按 Esc 键") {
		t.Errorf("the card names a global Esc hook this console never registered:\n%s", card)
	}
}

func Test170NoDeclaredChannelBooksTheLoudAuditLine(t *testing.T) {
	// reply stays nil - the non-TTY posture (pipe, script, CI). Nothing may
	// pretend a veto exists, and the assembly must SAY that posture rather
	// than leave it as an implied decision.
	h := newReplyHost(t, 40*time.Second)
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	audit := h.err.String()
	if !strings.Contains(audit, "approval: CLI-VETO-NONE reply_attached=false") {
		t.Errorf("audit is missing the loud no-veto line; got:\n%s", audit)
	}
	if strings.Contains(audit, "approval: CLI-VETO-ARMED") {
		t.Errorf("a run with no declared channel must not book an armed line:\n%s", audit)
	}
}
