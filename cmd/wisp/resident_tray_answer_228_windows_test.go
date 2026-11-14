//go:build windows

package main

// Ticket 228 AC#2 / AC#3 - the four rulers this slice owes, and what each one
// would go red on.
//
//	TestAC2TrayAllowOnceSpendsTheNativeGrant        AC#2: 允许一次 → Replies.Allow,
//	    end to end: the card is really answered on the native route and the tool
//	    call it was blocking really executes (`ANSWER-ALLOW ... route=native`).
//	TestAC2TrayRejectRefusesOnTheNativeRoute        AC#2: 拒绝 → Replies.Reject,
//	    end to end: the call is refused, the file never lands, the ledger line
//	    names route="native/reject".
//	TestAC2TrayAllowAlwaysRoutesIntoTheWideningVerb AC#2: 长期允许 →
//	    replySurface.always (it refuses here for the one reason a headless case can
//	    state honestly - no config write door - and the ROUTE it booked proves it
//	    reached that verb, not Allow).
//	TestAC2TrayAnswerWithNoPendingCardIsLoud         F4's 落空 shape: no card ⇒ the
//	    named "no pending card" line + no fake success, never a silent no-op.
//	TestAC2TrayAnswerWithoutAnExecutorSaysSo         the boot-window shape: items
//	    mounted, no executor attached yet ⇒ the existing unhosted-gesture record.
//	TestAC2TrayItemsAreMountedInTheBallPackage       internal/ball really builds
//	    the three items and dispatches them (the headless half of a tray click: the
//	    click itself needs a real notification area, ticket 228's own note).
//	TestAC3BookWaitingStateDrivesTheBallWaitingState AC#3: the D43 「等人」 name
//	    the runtime already computed now DRIVES the ball projection, on the
//	    production path (consoleApprovalUI.Prompt → bookWaitingState → sink).
//	TestAC3TheWaitingStateDriveIsWiredIntoTheShippedPath  the two wiring hops the
//	    resident leg owns (bindWaitingStateSink / attachTrayAnswers) are in the
//	    shipped assembly, as a token census that reds when either line is deleted.
//
// WHAT IS NOT MEASURED HERE, named instead of implied: a mouse click on a real
// tray menu. That needs a window in the notification area and a human hand
// (internal/ball's winlive family), so the click's own dispatch is pinned by the
// source census below and the wiring by the runtime cases; the live click stays
// ticket 228's owner-present pass.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent"
	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// cardExec228 is one in-flight fs.write a case uses to put a real L2 card on the
// ledger. The call is started on its own goroutine exactly the way
// approval_seam_201_test.go does: the gate blocks inside bridge.Execute until an
// answer arrives, so the answer has to come from the test's own goroutine.
type cardExec228 struct {
	out  agent.ToolOutcome
	err  error
	done chan struct{}
}

func startWriteCard228(rt *agentRuntime, corr, target string) *cardExec228 {
	c := &cardExec228{done: make(chan struct{})}
	args := json.RawMessage(fmt.Sprintf(`{"path":%q,"content":"written from the tray"}`, filepath.ToSlash(target)))
	go func() {
		defer close(c.done)
		c.out, c.err = rt.bridge.Execute(context.Background(), agent.ToolRequest{
			TaskID: "t228-tray", CorrelationID: corr, CallID: "a1", Name: "fs.write", Args: args,
		})
	}()
	return c
}

// trayExecutor228 builds the executor the way the shipped assembly does: a
// residentApproval over THIS run's ledger and gate, and a replySurface built by
// the runtime's own constructor. Nothing here re-implements a route.
func trayExecutor228(rt *agentRuntime) trayAnswerFunc {
	ra := &residentApproval{cards: rt.liveCards.h, gate: rt.gate}
	return ra.trayAnswerExecutor(rt.newReplySurface(nil, "", false))
}

// TestAC2TrayAllowOnceSpendsTheNativeGrant is AC#2's "有一枚真执行者", on the
// strongest reading available headless: the tray item's executor is called with
// the allow-once verb and the write the card was blocking really lands.
func TestAC2TrayAllowOnceSpendsTheNativeGrant(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	target := filepath.Join(h.outside, "landed-after-a-tray-allow-once.txt")
	var seamErr error

	h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask("t228-tray")
		defer revoke()
		exec := trayExecutor228(rt)
		card := startWriteCard228(rt, "t228-tray-allow", target)
		if !waitForCard187(rt.liveCards, "t228-tray-allow", 5*time.Second) {
			seamErr = fmt.Errorf("the L2 card never reached the host ledger, so no tray item could answer it")
			return
		}
		sentence := exec(trayAnswerAllowOnce)
		if !strings.Contains(sentence, "已允许") || !strings.Contains(sentence, "t228-tray-allow") {
			seamErr = fmt.Errorf("tray allow-once sentence = %q, want the surface's own allow sentence naming the card", sentence)
			return
		}
		<-card.done
		if card.err != nil {
			seamErr = fmt.Errorf("bridge.Execute: %w", card.err)
			return
		}
		if card.out.IsError {
			seamErr = fmt.Errorf("a call the tray allowed must execute, got: %s", card.out.Text)
		}
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	if seamErr != nil {
		t.Fatal(seamErr)
	}
	body, err := os.ReadFile(target)
	if err != nil || !strings.Contains(string(body), "written from the tray") {
		t.Fatalf("the tray's allow never reached the write: %v %q", err, body)
	}
	audit := h.err.String()
	// The gate's own native-route line: the answer travelled Replies.Allow →
	// DecideFromNative, not a panel route and not a re-implementation.
	const want = "approval: ANSWER-ALLOW corr=t228-tray-allow tool=fs.write route=native decision=allow"
	if !strings.Contains(audit, want) {
		t.Errorf("audit is missing %q; got:\n%s", want, audit)
	}
}

// TestAC2TrayRejectRefusesOnTheNativeRoute is the second of AC#2's three answer
// verbs: 拒绝 → Replies.Reject, and the call it was blocking never lands.
func TestAC2TrayRejectRefusesOnTheNativeRoute(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	target := filepath.Join(h.outside, "must-never-exist-after-a-tray-reject.txt")
	var seamErr error

	h.rtHook = func(rt *agentRuntime) {
		revoke := rt.gate.AdmitTextTask("t228-tray")
		defer revoke()
		exec := trayExecutor228(rt)
		card := startWriteCard228(rt, "t228-tray-reject", target)
		if !waitForCard187(rt.liveCards, "t228-tray-reject", 5*time.Second) {
			seamErr = fmt.Errorf("the L2 card never reached the host ledger")
			return
		}
		sentence := exec(trayAnswerReject)
		if !strings.Contains(sentence, "已拒绝") {
			seamErr = fmt.Errorf("tray reject sentence = %q, want the surface's own refusal sentence", sentence)
			return
		}
		<-card.done
		if card.err != nil {
			seamErr = fmt.Errorf("bridge.Execute: %w", card.err)
			return
		}
		if !card.out.IsError {
			seamErr = fmt.Errorf("a call the tray refused must not execute, got: %s", card.out.Text)
		}
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	if seamErr != nil {
		t.Fatal(seamErr)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("the refused call wrote its file anyway (stat err = %v)", err)
	}
	// The route fingerprint: a reject that reached Replies.Reject books
	// route="native/reject"; an allow-once would have booked the allow line.
	audit := h.err.String()
	const want = `approval: REPLY ANSWERED corr="t228-tray-reject" tool=fs.write route="native/reject"`
	if !strings.Contains(audit, want) {
		t.Errorf("audit is missing %q; got:\n%s", want, audit)
	}
}

// TestAC2TrayAllowAlwaysRoutesIntoTheWideningVerb is AC#2's third answer verb. The
// card is a real recorded L2 card; the surface is built WITHOUT a runtime on
// purpose, because that is the one shape in which `always` refuses for a reason a
// headless case can state honestly (no config write door) instead of spawning the
// second confirmation card and waiting 300s for an answer nobody can type. What
// this pins is the ROUTE: only always() books "native/always", so seeing it means
// the tray item reached that verb and not Replies.Allow.
func TestAC2TrayAllowAlwaysRoutesIntoTheWideningVerb(t *testing.T) {
	ledger := approval.NewReplies()
	ledger.Record(approval.Prompt{
		CorrelationID: "t228-tray-always", Level: "L2", Tool: "fs.write", TaskID: "t228-tray",
		Paths: []string{filepath.ToSlash(filepath.Join(t.TempDir(), "elsewhere"))},
	})
	var booked []string
	var mu sync.Mutex
	surface := &replySurface{
		live: &nativeCards{h: ledger},
		ctx:  context.Background(),
		audit: func(format string, args ...any) {
			mu.Lock()
			booked = append(booked, fmt.Sprintf(format, args...))
			mu.Unlock()
		},
	}
	ra := &residentApproval{cards: ledger}
	exec := ra.trayAnswerExecutor(surface)

	sentence := exec(trayAnswerAllowAlways)
	if !strings.Contains(sentence, "长期允许") {
		t.Errorf("the no-config-write-door sentence must name the item that was pressed, got %q", sentence)
	}
	mu.Lock()
	defer mu.Unlock()
	joined := strings.Join(booked, "\n")
	if !strings.Contains(joined, `route="native/always"`) {
		t.Fatalf("the 长期允许 item did not reach replySurface.always; booked lines:\n%s", joined)
	}
	if strings.Contains(joined, `route="native/allow"`) {
		t.Errorf("the 长期允许 item booked the plain allow route; booked lines:\n%s", joined)
	}
}

// TestAC2TrayAnswerWithNoPendingCardIsLoud is F4's ruling made executable: the
// three items are always visible, so a press with nothing to answer must take a
// NAMED path - one log line carrying the fact, and a sentence that cannot be read
// as success. A silent no-op or a fabricated "已允许" is the defect class.
func TestAC2TrayAnswerWithNoPendingCardIsLoud(t *testing.T) {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	defer slog.SetDefault(old)

	ra := &residentApproval{cards: approval.NewReplies()}
	exec := ra.trayAnswerExecutor(&replySurface{})
	sentence := exec(trayAnswerAllowOnce)

	if !strings.Contains(sentence, "no pending card") {
		t.Errorf("the no-card sentence must carry the literal words the ruling names, got %q", sentence)
	}
	if !strings.Contains(sentence, "nothing was allowed, rejected or remembered") {
		t.Errorf("the no-card sentence must state that nothing happened, got %q", sentence)
	}
	for _, forbidden := range []string{"已允许", "已拒绝", "ANSWERED"} {
		if strings.Contains(sentence, forbidden) {
			t.Errorf("the no-card sentence reads like a success (%q in %q)", forbidden, sentence)
		}
	}
	logged := buf.String()
	for _, want := range []string{"tray answer: no pending card", "verb=allow-once"} {
		if !strings.Contains(logged, want) {
			t.Errorf("the no-card press booked no record containing %q; log:\n%s", want, logged)
		}
	}
}

// TestAC2TrayAnswerWithoutAnExecutorSaysSo is the boot-window shape: the three
// items are mounted from the first frame, but the executor is attached only once
// the task pipeline exists. A press in between must be recorded as the unhosted
// gesture it is, by name - not dropped and not answered by this file.
func TestAC2TrayAnswerWithoutAnExecutorSaysSo(t *testing.T) {
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
	defer slog.SetDefault(old)

	rb := &residentBall{} // no window, no executor: the pre-attach shape
	rb.trayAnswerGesture(trayAnswerAllowOnce)

	logged := buf.String()
	for _, want := range []string{"ball gesture arrived with no executor", "tray-allow-once"} {
		if !strings.Contains(logged, want) {
			t.Errorf("an unattached tray answer press booked no record containing %q; log:\n%s", want, logged)
		}
	}
}

// TestAC2TrayItemsAreMountedInTheBallPackage is the headless half of "the tray
// really carries the three items": internal/ball builds them and its wndproc
// dispatches each id to its own callback. The click itself needs the notification
// area (and a human), so this reads the source instead of pressing it - deleting
// an appendItem line, a case line, or the id constant turns it red by name.
func TestAC2TrayItemsAreMountedInTheBallPackage(t *testing.T) {
	fset := token.NewFileSet()
	trayPath := filepath.Join("..", "..", "internal", "ball", "tray_windows.go")
	trayFile, err := parser.ParseFile(fset, trayPath, nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", trayPath, err)
	}

	// 1. ids are appended, never renumbered: 1..4 stay exactly as shipped.
	ids := map[string]int{}
	ast.Inspect(trayFile, func(n ast.Node) bool {
		spec, ok := n.(*ast.ValueSpec)
		if !ok || len(spec.Names) != 1 || len(spec.Values) != 1 {
			return true
		}
		lit, ok := spec.Values[0].(*ast.BasicLit)
		if !ok || lit.Kind != token.INT {
			return true
		}
		var v int
		if _, err := fmt.Sscanf(lit.Value, "%d", &v); err == nil {
			ids[spec.Names[0].Name] = v
		}
		return true
	})
	for name, want := range map[string]int{
		"menuOpenPanel": 1, "menuMute": 2, "menuPauseWake": 3, "menuExit": 4,
		"menuAllowOnce": 5, "menuReject": 6, "menuAllowAlways": 7,
	} {
		if got, ok := ids[name]; !ok || got != want {
			t.Errorf("internal/ball id %s = %d (present=%v), want %d - the answer items are APPENDED to 1..4, never renumbered", name, got, ok, want)
		}
	}

	// 2. the three labels exist, in the menu builder.
	labels := map[string]string{}
	ast.Inspect(trayFile, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || len(call.Args) != 3 {
			return true
		}
		if fn, ok := call.Fun.(*ast.Ident); !ok || fn.Name != "appendItem" {
			return true
		}
		id, ok := call.Args[0].(*ast.Ident)
		if !ok {
			return true
		}
		lit, ok := call.Args[1].(*ast.BasicLit)
		if !ok {
			return true
		}
		labels[id.Name] = strings.Trim(lit.Value, `"`)
		return true
	})
	for id, want := range map[string]string{
		"menuAllowOnce": "允许一次", "menuReject": "拒绝", "menuAllowAlways": "长期允许",
	} {
		if got := labels[id]; got != want {
			t.Errorf("the tray's %s item is labelled %q, want %q", id, got, want)
		}
	}

	// 3. the wndproc dispatches each id to its OWN callback - the whole mapping.
	ballPath := filepath.Join("..", "..", "internal", "ball", "ball_windows.go")
	ballFile, err := parser.ParseFile(fset, ballPath, nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", ballPath, err)
	}
	dispatched := map[string][]string{}
	ast.Inspect(ballFile, func(n ast.Node) bool {
		cc, ok := n.(*ast.CaseClause)
		if !ok || len(cc.List) != 1 {
			return true
		}
		id, ok := cc.List[0].(*ast.Ident)
		if !ok {
			return true
		}
		for _, stmt := range cc.Body {
			ast.Inspect(stmt, func(m ast.Node) bool {
				call, ok := m.(*ast.CallExpr)
				if !ok || len(call.Args) != 1 {
					return true
				}
				if sel, ok := call.Fun.(*ast.SelectorExpr); !ok || sel.Sel.Name != "fire" {
					return true
				}
				arg, ok := call.Args[0].(*ast.SelectorExpr)
				if !ok {
					return true
				}
				if strings.HasPrefix(arg.Sel.Name, "OnTray") {
					dispatched[id.Name] = append(dispatched[id.Name], arg.Sel.Name)
				}
				return true
			})
		}
		return true
	})
	for id, want := range map[string]string{
		"menuAllowOnce": "OnTrayAllowOnce", "menuReject": "OnTrayReject", "menuAllowAlways": "OnTrayAllowAlways",
	} {
		got := dispatched[id]
		if len(got) != 1 || got[0] != want {
			t.Errorf("wndproc dispatches %s to %v, want exactly [%s]", id, got, want)
		}
	}
}

// TestAC3BookWaitingStateDrivesTheBallWaitingState is AC#3: the D43 name the
// runtime already computed (AC#6a) now reaches the ball, and it does so on the
// PRODUCTION path - this case never calls bookWaitingState itself. It binds the
// projection the resident leg binds, lets consoleApprovalUI.Prompt raise a real
// card, and reads what arrived.
//
// The mutation it is built for: delete the projectWaitingState call from
// bookWaitingState and this case reports the empty reading instead of green -
// which is exactly the control the ticket says the cell did not have.
func TestAC3BookWaitingStateDrivesTheBallWaitingState(t *testing.T) {
	h := newReplyHost(t, 40*time.Second)
	target := filepath.Join(h.outside, "ac3-card-target.txt")
	// The ball host stands outside the hook so the case can read what IT recorded,
	// not only what the sink saw: both halves of the hop are this test's subject.
	rb := &residentBall{}
	var (
		mu     sync.Mutex
		sunk   []statemachine.State
		errBox error
	)

	h.rtHook = func(rt *agentRuntime) {
		// The admitted task id has to be the one the bridge call carries
		// (startWriteCard228's "t228-tray"): an unadmitted task is refused before
		// any card exists, which reads as "the card never reached the ledger".
		revoke := rt.gate.AdmitTextTask("t228-tray")
		defer revoke()
		// The ball host's half, wired the way resident_windows.go wires it
		// (bindWaitingStateSink(rb.projectWaitingState)) - the recording wrapper
		// IS that binding, with this case's instrument wrapped around the same
		// call, so what is measured is the shipped hop and not a stand-in for it.
		rt.bindWaitingStateSink(func(s statemachine.State) {
			mu.Lock()
			sunk = append(sunk, s)
			mu.Unlock()
			rb.projectWaitingState(s)
		})
		card := startWriteCard228(rt, "t228-ac3-corr", target)
		if !waitForCard187(rt.liveCards, "t228-ac3-corr", 5*time.Second) {
			errBox = fmt.Errorf("the L2 card never reached the host ledger")
			return
		}
		// Answer it so the run can finish; the reading below was taken already.
		_ = rt.liveCards.h.Reject(context.Background(), "t228-ac3-corr", "ac3 case: refused after the reading")
		<-card.done
	}
	if code := h.run("总结一下 这份笔记"); code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, h.out.String(), h.err.String())
	}
	if errBox != nil {
		t.Fatal(errBox)
	}
	mu.Lock()
	got := append([]statemachine.State(nil), sunk...)
	mu.Unlock()
	if len(got) == 0 {
		t.Fatalf("bookWaitingState drove nothing onto the ball projection; the producer only wrote its audit line (AC#3's defect)")
	}
	if got[0] != statemachine.StateAwaitingApproval {
		t.Errorf("the first projection was %q, want %q (an L2 card is D43's 等人 row)", got[0], statemachine.StateAwaitingApproval)
	}
	// And the BALL recorded the same walk: the sink above is this case's
	// instrument, residentBall.projectWaitingState is the shipped projection, and
	// the two must agree or "the orb entered it" would be a claim about the
	// instrument rather than about the host.
	ballWalk := rb.waitingStatesSnapshot()
	if len(ballWalk) == 0 || ballWalk[0] != statemachine.StateAwaitingApproval {
		t.Errorf("the ball host recorded %v, want [%q, …] (its own projection never ran)", ballWalk, statemachine.StateAwaitingApproval)
	}
	const wantLine = `approval: WAITING-STATE state="AwaitingApproval" awaiting=true corr="t228-ac3-corr"`
	if audit := h.err.String(); !strings.Contains(audit, wantLine) {
		t.Errorf("audit is missing %q; got:\n%s", wantLine, audit)
	}
}

// TestAC3TheWaitingStateDriveIsWiredIntoTheShippedPath is the wiring census: the
// two hops have to exist in the shipped assembly and on the real card surface, not
// merely in a helper a case calls. Deleting either line turns this red by name.
func TestAC3TheWaitingStateDriveIsWiredIntoTheShippedPath(t *testing.T) {
	// 1. The card surface asks the runtime's producer, exactly once, on every card
	// it shows (the established census instrument, see approval_badge_21s2).
	census := census21S2(t, "resident_approval_windows.go")
	if n := census["Prompt"]["bookWaitingState"]; n != 1 {
		t.Errorf("ballCardUI.Prompt calls bookWaitingState %d times, want exactly 1 - without it the driver has no production caller", n)
	}

	// 2. The resident assembly binds both hops, on the line where the runtime (and
	// therefore its ledger and answer surface) first exists.
	body, err := os.ReadFile("resident_windows.go")
	if err != nil {
		t.Fatalf("read resident_windows.go: %v", err)
	}
	for _, want := range []string{
		"src.run.bindWaitingStateSink(rb.projectWaitingState)",
		"rb.attachTrayAnswers(ra.trayAnswerExecutor(src.answers))",
		"ra.bindRuntime(src.run)",
	} {
		if !bytes.Contains(body, []byte(want)) {
			t.Errorf("the shipped resident assembly is missing the line %q", want)
		}
	}
}
