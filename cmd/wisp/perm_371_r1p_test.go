package main

// Ticket 371-b at the PRODUCTION seam: the mode-switch audit row has to carry an
// "asked" reading that joins the approval queue's own ANSWER-ALLOW line.
//
// WHY THIS FILE EXISTS, SEPARATELY FROM internal/perm's cases. perm can only test
// what it is handed: it cannot know whether the limb it called ever reached a
// card. The credential is minted in exactly one production function -
// confirmModeSwitchVia - and this file drives THAT function against a REAL
// approval.Gate, a real queue, and a real native allow, then reads the two audit
// rows and joins them by correlation id:
//
//	perm: MODE-SWITCH ... asked=1 card="host:mode-switch#N" result=applied
//	approval: ANSWER-ALLOW corr=host:mode-switch#N tool=permission.mode route=native decision=allow
//
// The id is chosen by the limb (the queue does not return it), which is why the
// join is possible at all: a row that could only print asked=1 would answer "did
// anyone answer?" but not "which card was it?".
//
// The second case is the counter-example, and it is the hole 371-a1p measured
// moved one layer down: with no click, the same production limb must report NO
// answered card, and the store must refuse while the row says asked=0.

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/perm"
	"github.com/CarlosShao/wisp/internal/risk"
)

// log371 is the audit sink both the gate and the store write into.
type log371 struct {
	mu    sync.Mutex
	lines []string
}

func (l *log371) logf(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lines = append(l.lines, fmt.Sprintf(format, args...))
}

func (l *log371) all() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.lines...)
}

func (l *log371) line(needles ...string) string {
	var hit string
	for _, line := range l.all() {
		ok := true
		for _, n := range needles {
			if !strings.Contains(line, n) {
				ok = false
				break
			}
		}
		if ok {
			hit = line
		}
	}
	return hit
}

// nativeCard371 is the native surface the gate is handed: it records the prompt it
// was asked to show (which is where the grant lives, by design) and closes ready
// once. It never answers by itself - the test does that through Gate.Native, so
// "somebody clicked allow" stays the test's action and not the fixture's opinion.
type nativeCard371 struct {
	mu    sync.Mutex
	p     approval.Prompt
	ready chan struct{}
}

func newNativeCard371() *nativeCard371 {
	return &nativeCard371{ready: make(chan struct{})}
}

func (n *nativeCard371) Prompt(_ context.Context, p approval.Prompt) error {
	n.mu.Lock()
	n.p = p
	n.mu.Unlock()
	close(n.ready)
	return nil
}

func (n *nativeCard371) Update(context.Context, approval.Event) error { return nil }

func (n *nativeCard371) prompt(t *testing.T) approval.Prompt {
	t.Helper()
	select {
	case <-n.ready:
	case <-time.After(5 * time.Second):
		t.Fatal("the L2 card never reached the native surface")
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.p
}

// store371 builds the smallest honest assembly: a real config.toml on a real
// Manager, and a real perm.Store whose limb is the PRODUCTION body.
func store371(t *testing.T, gate *approval.Gate, logf func(string, ...any)) (*perm.Store, string) {
	t.Helper()
	dir := newInboundDataDir33(t, risk.ModeAskEveryStepName)
	mgr, err := config.NewManager(filepath.Join(dir, configFileName), nil)
	if err != nil {
		t.Fatalf("config.NewManager: %v", err)
	}
	store, err := perm.New(perm.Options{
		Manager: mgr,
		Confirm: func(ctx context.Context, sw perm.Switch) (perm.Card, error) {
			return confirmModeSwitchVia(gate, ctx, sw)
		},
		Logf: logf,
	})
	if err != nil {
		t.Fatalf("perm.New: %v", err)
	}
	return store, dir
}

// Test371R1PTheAnsweredCardJoinsTheTwoAuditRows is the accepted reading of
// ticket 371-b: one switch into auto_approve, one native allow, and both records
// name the same card.
func Test371R1PTheAnsweredCardJoinsTheTwoAuditRows(t *testing.T) {
	sink := &log371{}
	ui := newNativeCard371()
	gate := approval.New(approval.Options{UI: ui, Logf: sink.logf})
	store, dir := store371(t, gate, sink.logf)

	var (
		mu     sync.Mutex
		setErr error
		done   = make(chan struct{})
	)
	go func() {
		defer close(done)
		mu.Lock()
		defer mu.Unlock()
		setErr = store.Set(context.Background(), risk.ModeAutoApprove, "panel", "owner")
	}()

	p := ui.prompt(t)
	// The card the limb raised must be identifiable and it must not be the queue's
	// own "approval-N": the credential has to carry a value chosen before the gate
	// was called, or the store's row could not name it.
	if !strings.HasPrefix(p.CorrelationID, modeSwitchTaskID+"#") {
		t.Errorf("the mode-switch card id = %q, want %q + a sequence number (the limb chooses it so "+
			"the store's audit row can name it)", p.CorrelationID, modeSwitchTaskID)
	}
	if p.Tool != modeSwitchToolName || p.Level != "L2" {
		t.Errorf("the card is tool=%q level=%q, want %q/L2", p.Tool, p.Level, modeSwitchToolName)
	}
	if p.Grant == "" {
		t.Fatal("the native prompt carries no grant; a native allow needs one")
	}
	if err := gate.Native().Allow(context.Background(), p.CorrelationID, p.Grant); err != nil {
		t.Fatalf("the native allow was refused: %v", err)
	}
	<-done

	mu.Lock()
	err := setErr
	mu.Unlock()
	if err != nil {
		t.Fatalf("an answered card was refused by the store: %v", err)
	}
	if got := readInboundCfg33(t, dir); got != risk.ModeAutoApproveName {
		t.Errorf("the answered switch did not land on disk: %q", got)
	}

	row := sink.line("perm: MODE-SWITCH", "to="+risk.ModeAutoApproveName)
	if row == "" {
		t.Fatalf("no MODE-SWITCH row for auto_approve in:\n%s", strings.Join(sink.all(), "\n"))
	}
	for _, want := range []string{"asked=1", `card="` + p.CorrelationID + `"`, "result=" + string(perm.ResultApplied)} {
		if !strings.Contains(row, want) {
			t.Errorf("the MODE-SWITCH row is missing %q: %s", want, row)
		}
	}
	allow := sink.line("approval: ANSWER-ALLOW corr=" + p.CorrelationID)
	if allow == "" {
		t.Fatalf("no ANSWER-ALLOW line for corr=%s in:\n%s", p.CorrelationID, strings.Join(sink.all(), "\n"))
	}
	// The join itself: the id on the switch row is the id on the answer line.
	if !strings.Contains(row, p.CorrelationID) {
		t.Errorf("the two rows cannot be joined: switch=%q answer=%q", row, allow)
	}
	t.Logf("371R1P|join|%s", row)
	t.Logf("371R1P|join|%s", allow)
}

// Test371R1PNoClickMeansNoAnsweredCard is the counter-example on the same limb: a
// card is raised, nobody answers, the context ends, and the switch is refused with
// asked=0. Without this, "asked=1" could be a constant.
func Test371R1PNoClickMeansNoAnsweredCard(t *testing.T) {
	sink := &log371{}
	ui := newNativeCard371()
	gate := approval.New(approval.Options{UI: ui, Logf: sink.logf})
	store, dir := store371(t, gate, sink.logf)

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	err := store.Set(ctx, risk.ModeAutoApprove, "panel", "owner")
	if err == nil {
		t.Fatal("an unanswered card still applied the switch")
	}
	if got := readInboundCfg33(t, dir); got != risk.ModeAskEveryStepName {
		t.Errorf("the unanswered switch moved the档 on disk: %q", got)
	}
	row := sink.line("perm: MODE-SWITCH", "to="+risk.ModeAutoApproveName)
	if row == "" {
		t.Fatalf("no MODE-SWITCH row in:\n%s", strings.Join(sink.all(), "\n"))
	}
	if !strings.Contains(row, "asked=0") || !strings.Contains(row, `card=""`) {
		t.Errorf("the unanswered attempt claims a card was answered: %s", row)
	}
	if strings.Contains(row, "result="+string(perm.ResultCardNotAnswered)) {
		t.Errorf("a card that was shown and timed out is booked as 'not answered'; the two stories "+
			"differ and the row must keep them apart: %s", row)
	}
	if allow := sink.line("approval: ANSWER-ALLOW"); allow != "" {
		t.Errorf("an allow line exists although nothing clicked allow: %s", allow)
	}
}

// Test371R1PThePanelAdaptersPreserveTheAnsweredCard pins the two adapters that
// translate one limb between the layers: they must carry the credential, not
// flatten it to an error. A flattening adapter would put the panel seam back on
// "no error means allowed" while the run seam judged the credential - two gates,
// two rules - which is exactly what AC#2'-b forbids.
func Test371R1PThePanelAdaptersPreserveTheAnsweredCard(t *testing.T) {
	ctx := context.Background()

	// A limb that can ask nothing reports the zero credential.
	silent := panelModeConfirm(func(context.Context, perm.Switch) (perm.Card, error) {
		return perm.Card{}, nil
	})
	if silent == nil {
		t.Fatal("a non-nil perm limb became a nil panel limb")
	}
	card, err := silent(ctx, risk.ModeAskEveryStep, risk.ModeAutoApprove, "panel", "owner")
	if err != nil || card.Answered() {
		t.Fatalf("the panel adapter invented an answer: card=%+v err=%v", card, err)
	}
	back := storeConfirmFromPanel(silent)
	card, err = back(ctx, perm.Switch{From: risk.ModeAskEveryStep, To: risk.ModeAutoApprove})
	if err != nil || card.Answered() {
		t.Fatalf("the round trip invented an answer: card=%+v err=%v", card, err)
	}

	// A limb that answered keeps its card id through both hops.
	const corr = "host:mode-switch#99"
	answered := panelModeConfirm(func(context.Context, perm.Switch) (perm.Card, error) {
		return perm.AnsweredCard(corr), nil
	})
	card, err = answered(ctx, risk.ModeAskEveryStep, risk.ModeAutoApprove, "panel", "owner")
	if err != nil || !card.Answered() || card.Corr() != corr {
		t.Fatalf("the panel adapter dropped the credential: card=%+v err=%v", card, err)
	}
	card, err = storeConfirmFromPanel(answered)(ctx, perm.Switch{From: risk.ModeAskEveryStep, To: risk.ModeAutoApprove})
	if err != nil || !card.Answered() || card.Corr() != corr {
		t.Fatalf("the round trip dropped the credential: card=%+v err=%v", card, err)
	}

	// And a nil leg stays nil in both directions: a closure that only refused when
	// called would switch both gates off.
	if panelModeConfirm(nil) != nil || storeConfirmFromPanel(nil) != nil {
		t.Error("a nil leg became a non-nil closure; both gates key on its presence")
	}
}
