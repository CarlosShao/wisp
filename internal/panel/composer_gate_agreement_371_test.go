package panel

// Ticket 371 AC#2'-b: the panel gate and the store gate must judge the SAME rule,
// and the layer that asks must not be the only one that can name the verdict.
//
// WHY A SECOND AGREEMENT FILE. composer_gate_agreement_367_test.go pinned the
// TARGET rule for the no-channel matrix. Ticket 371 adds the ANSWER to the rule:
// auto_approve + a non-nil limb that never raises a card used to apply silently
// (docs/evidence/s1/371-a1p-impl.md section 5). The store now refuses that with
// its own verdict (perm.ResultCardNotAnswered). The panel cannot judge that answer
// - it raises no card, on purpose (composer_handlers.go's header) - but it must
// not describe it differently: it books the store's refusal with the SAME
// classification, on its own line, with its requestId. Otherwise "one rule, two
// gates" is a claim about a function call, and the operator's log shows two
// vocabularies for one refusal.
//
// THE ROWS BELOW ARE THE INVARIANT, not either implementation:
//
//	no channel     -> the PANEL refuses before the write leg (one panel line)
//	never-asking   -> the panel passes, the STORE refuses, BOTH layers book
//	                  result=card-not-answered, and nothing lands
//	answered allow -> the switch applies, one card, no panel line
//	target M4 does not card -> the switch applies, the limb is never entered
//
// Each row runs a REAL perm.Store behind the handler (no fake verdict), and the
// store's verdict is read from its own audit row rather than from the error the
// handler returned - so "both agree" is measured against two independent records.
//
// THE REVERSE TOOTH is specified where it belongs: the two cases at the bottom pin
// the shape of the panel's half (it calls perm.JudgeCard, and its ModeConfirm
// carries the credential) so that updating perm alone cannot leave this gate
// spelling the old rule. docs/evidence/s1/371-r1p-impl.md measures the overlay run
// that puts the panel half back and watches the never-asking row go red.

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/CarlosShao/wisp/internal/perm"
	"github.com/CarlosShao/wisp/internal/risk"
)

// t371Limb is the confirmation channel BOTH layers are handed, in the two shapes
// the assembly really uses: perm's ConfirmFunc (the store) and ModeConfirm (the
// handler). One limb, two adapters - the same relationship cmd/wisp's
// panelModeConfirm / storeConfirmFromPanel have, which their own package pins.
type t371Limb struct {
	mu     sync.Mutex
	asks   int
	answer perm.Card
}

func (l *t371Limb) storeLeg() perm.ConfirmFunc {
	return func(context.Context, perm.Switch) (perm.Card, error) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.asks++
		return l.answer, nil
	}
}

func (l *t371Limb) panelLeg() ModeConfirm {
	return func(context.Context, risk.Mode, risk.Mode, string, string) (perm.Card, error) {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.asks++
		return l.answer, nil
	}
}

func (l *t371Limb) asked() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.asks
}

// t371Writer is the write leg the handler crosses: the REAL store, with a call
// count so "the panel refused before the write leg" is a reading and not a claim.
type t371Writer struct {
	store *perm.Store
	calls int
}

func (w *t371Writer) PermissionMode() risk.Mode { return w.store.PermissionMode() }

func (w *t371Writer) Set(ctx context.Context, to risk.Mode, origin, actor string) error {
	w.calls++
	return w.store.Set(ctx, to, origin, actor)
}

// t371StoreRow returns the store's own MODE-SWITCH row for the target.
func t371StoreRow(t *testing.T, log *auditLog, to string) string {
	t.Helper()
	var hit string
	for _, l := range log.all() {
		if strings.Contains(l, "MODE-SWITCH") && strings.Contains(l, "to="+to) {
			hit = l
		}
	}
	if hit == "" {
		t.Fatalf("no MODE-SWITCH row for to=%s in:\n%s", to, log.joined())
	}
	return hit
}

// all is the slice form the joined helper does not expose.
func (a *auditLog) all() []string { return append([]string(nil), (*a)...) }

// TestPanelGateAndStoreGateAgreeOnTheAnsweredCard371 is the matrix.
func TestPanelGateAndStoreGateAgreeOnTheAnsweredCard371(t *testing.T) {
	ladder := modeLadder(t)
	names := risk.ModeNames()
	ctx := context.Background()
	pairs, kinds := 0, 0
	refusedByPanel, refusedByStore, applied := 0, 0, 0

	for fromIdx, from := range ladder {
		for toIdx, to := range ladder {
			if fromIdx == toIdx {
				// A no-change row never reaches a gate in either layer; ticket
				// 367's matrix and 114's handler cases own it.
				continue
			}
			pairs++
			for _, kind := range []string{"no channel", "never-asking", "answered allow"} {
				kinds++
				limb := &t371Limb{}
				var storeLeg perm.ConfirmFunc
				var panelLeg ModeConfirm
				if kind != "no channel" {
					storeLeg, panelLeg = limb.storeLeg(), limb.panelLeg()
				}
				if kind == "answered allow" {
					limb.answer = perm.AnsweredCard("t371/panel")
				}

				mgr := &t367Manager{mode: from}
				var storeLog, panelLog auditLog
				store, err := perm.New(perm.Options{Manager: mgr, Confirm: storeLeg, Logf: storeLog.f})
				if err != nil {
					t.Fatalf("%s -> %s / %s: perm.New: %v", from, to, kind, err)
				}
				writer := &t371Writer{store: store}
				h := &ModeWriteHandler{Modes: writer, Confirm: panelLeg, Audit: panelLog.f, Actor: "t371"}
				id := fmt.Sprintf("t371-%s-%s-%s", names[fromIdx], names[toIdx], strings.ReplaceAll(kind, " ", "-"))
				handlerErr := h.HandleModeRequest(ctx, modeReq(names[toIdx], id))

				carded := perm.RequiresCard(to)
				wantRefused := carded && kind != "answered allow"

				// (1) The handler's own verdict, before either record is read.
				if (handlerErr != nil) != wantRefused {
					t.Errorf("%s -> %s / %s: panel refused=%v (%v), want refused=%v",
						from, to, kind, handlerErr != nil, handlerErr, wantRefused)
				}
				// (2) The persisted档, which is the thing a user feels.
				if wantRefused {
					if mgr.mode != from {
						t.Errorf("%s -> %s / %s: the refusal still moved the persisted档 to %s",
							from, to, kind, mgr.mode)
					}
				} else if mgr.mode != to {
					t.Errorf("%s -> %s / %s: the allowed switch did not persist %s (file says %s)",
						from, to, kind, to, mgr.mode)
				}
				// (3) Who refused. Only the no-channel row may be refused WITHOUT
				// the store ever running; every other verdict must come from the
				// store, which is the layer that asks.
				if writer.calls == 0 {
					refusedByPanel++
					if !(carded && kind == "no channel") {
						t.Fatalf("%s -> %s / %s: the write leg was never reached", from, to, kind)
					}
					if handlerErr == nil || !strings.Contains(handlerErr.Error(), "没有接入 L2 确认腿") {
						t.Errorf("%s -> %s / %s: the refusal is not the panel's own: %v", from, to, kind, handlerErr)
					}
					if limb.asked() != 0 {
						t.Errorf("%s -> %s / %s: the limb was entered with no channel attached", from, to, kind)
					}
					if len(panelLog) != 1 || !strings.Contains(panelLog.joined(), id) {
						t.Errorf("%s -> %s / %s: want exactly one panel line naming %q, got %q",
							from, to, kind, id, panelLog.joined())
					}
					continue
				}

				row := t371StoreRow(t, &storeLog, names[toIdx])
				storeRefused := !strings.Contains(row, "result="+string(perm.ResultApplied))
				if storeRefused != wantRefused {
					t.Errorf("%s -> %s / %s: the store's own row says refused=%v, want %v: %s",
						from, to, kind, storeRefused, wantRefused, row)
				}
				if (handlerErr != nil) != storeRefused {
					t.Errorf("%s -> %s / %s: THE TWO GATES DISAGREE - panel error=%v, store row refused=%v (%s)",
						from, to, kind, handlerErr, storeRefused, row)
				}

				switch {
				case carded && kind == "never-asking":
					refusedByStore++
					// The panel passed it, the store asked and got nothing back:
					// the two layers must name ONE rule for it.
					if writer.calls != 1 || limb.asked() != 1 {
						t.Errorf("%s -> %s / %s: writeCalls=%d limbEnters=%d, want 1/1 (the store asks once)",
							from, to, kind, writer.calls, limb.asked())
					}
					// handlerErr may be nil on the should-refuse-but-did-not path;
					// fail the assertion with a clear message instead of panicking
					// while trying to read the missing error.
					if handlerErr == nil {
						t.Errorf("%s -> %s / %s: the store refused the unanswered card but the panel "+
							"surfaced no error; the refusal never reached the page", from, to, kind)
					} else if !strings.Contains(handlerErr.Error(), "card was not answered") {
						t.Errorf("%s -> %s / %s: the surfaced refusal does not name the unanswered card: %v",
							from, to, kind, handlerErr)
					}
					if !strings.Contains(row, "result="+string(perm.ResultCardNotAnswered)) ||
						!strings.Contains(row, "asked=0") {
						t.Errorf("%s -> %s / %s: the store row does not book the unanswered card: %s",
							from, to, kind, row)
					}
					if len(panelLog) != 1 {
						t.Fatalf("%s -> %s / %s: the panel wrote %d refusal line(s), want 1: %q",
							from, to, kind, len(panelLog), panelLog.joined())
					}
					panelLine := panelLog.joined()
					if !strings.Contains(panelLine, "result="+string(perm.ResultCardNotAnswered)) {
						t.Errorf("%s -> %s / %s: the panel booked the same refusal as a different thing (%q), "+
							"while the store booked %q - that is the drift AC#2'-b exists to stop",
							from, to, kind, panelLine, perm.ResultCardNotAnswered)
					}
					if !strings.Contains(panelLine, id) {
						t.Errorf("%s -> %s / %s: the panel's line does not name the requestId %q: %q",
							from, to, kind, id, panelLine)
					}
				case carded:
					applied++
					if writer.calls != 1 || limb.asked() != 1 {
						t.Errorf("%s -> %s / %s: writeCalls=%d asks=%d, want exactly 1/1",
							from, to, kind, writer.calls, limb.asked())
					}
					if len(panelLog) != 0 {
						t.Errorf("%s -> %s / %s: a request that reached the write leg added a verdict: %q",
							from, to, kind, panelLog.joined())
					}
					if !strings.Contains(row, "asked=1") || !strings.Contains(row, `card="t371/panel"`) {
						t.Errorf("%s -> %s / %s: the applied row does not name the answered card: %s",
							from, to, kind, row)
					}
				default:
					// !carded: M4 says this target costs no card in ANY assembly.
					applied++
					if writer.calls != 1 {
						t.Errorf("%s -> %s / %s: writeCalls=%d, want 1", from, to, kind, writer.calls)
					}
					if limb.asked() != 0 {
						t.Errorf("%s -> %s / %s: R20/M4 was widened - the limb was entered for a target "+
							"that costs no card", from, to, kind)
					}
					if len(panelLog) != 0 {
						t.Errorf("%s -> %s / %s: a passing request added a panel verdict: %q",
							from, to, kind, panelLog.joined())
					}
				}
			}
		}
	}

	// The counts are the matrix's own statement: 6 ordered pairs x 3 limb kinds.
	if pairs != 6 || kinds != 18 {
		t.Fatalf("matrix covered %d pairs / %d rows, want 6/18: the enumeration was shortened", pairs, kinds)
	}
	// Two auto_approve pairs are refused by the panel (no channel) and two by the
	// store (never-asking); the rest apply.
	if refusedByPanel != 2 || refusedByStore != 2 {
		t.Errorf("refused by panel=%d by store=%d, want 2/2 (auto_approve from the two stricter档)",
			refusedByPanel, refusedByStore)
	}
	if want := kinds - refusedByPanel - refusedByStore; applied != want {
		t.Errorf("%d rows applied, want %d", applied, want)
	}
}

// Test371R1PModeConfirmCarriesTheCredential is the compile-and-run half of the
// same invariant: the panel's leg reports perm.Card, i.e. the answer the store
// judges, so no adapter on this path can flatten an answered allow to "no error".
// If a future change put ModeConfirm back on `error`, this file would not compile -
// which is the strongest red a cross-package contract change can have, and it is
// named here so it is a decision rather than an accident.
func Test371R1PModeConfirmCarriesTheCredential(t *testing.T) {
	var leg ModeConfirm = func(context.Context, risk.Mode, risk.Mode, string, string) (perm.Card, error) {
		return perm.AnsweredCard("t371/direct"), nil
	}
	card, err := leg(context.Background(), risk.ModeAskEveryStep, risk.ModeAutoApprove, "panel", "owner")
	if err != nil {
		t.Fatalf("the limb reported an error: %v", err)
	}
	if !card.Answered() || card.Corr() != "t371/direct" {
		t.Errorf("the credential did not survive the panel's own contract: %+v", card)
	}
	// And the zero value is still "not answered", which is what a limb that never
	// raised a card reports.
	var zero perm.Card
	if zero.Answered() {
		t.Error("the zero perm.Card reads as answered; a never-asking limb would pass the gate again")
	}
}

// Test371R1PThePanelGateCallsTheSameCardPredicate is the source pin for the CALL
// SITE. The agreement matrix cannot see which function a gate called, and that is
// exactly the edit AC#2'-b is about: update perm and leave this file spelling
// `RequiresCard(to) && h.Confirm == nil`, and every behavioural row still agrees
// because both terms are equivalent today. The rule drifts the day one of them
// changes; this pin makes that day red.
func Test371R1PThePanelGateCallsTheSameCardPredicate(t *testing.T) {
	const path = "composer_handlers.go"
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read %s: %v", path, err)
	}
	f, err := parser.ParseFile(token.NewFileSet(), path, src, 0)
	if err != nil {
		t.Fatalf("cannot parse %s: %v", path, err)
	}

	var calls []*ast.CallExpr
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel == nil || sel.Sel.Name != "JudgeCard" {
			return true
		}
		if id, ok := sel.X.(*ast.Ident); ok && id.Name == "perm" {
			calls = append(calls, call)
		}
		return true
	})
	if len(calls) != 1 {
		t.Fatalf("%s calls perm.JudgeCard %d time(s), want exactly 1: the panel's half of the mode gate "+
			"must be the SAME judgement perm.Store.Set applies (ticket 371-b), or the two gates are two rules",
			path, len(calls))
	}
	call := calls[0]
	if len(call.Args) != 2 {
		t.Fatalf("perm.JudgeCard takes (target, channel); the gate passes %d argument(s)", len(call.Args))
	}
	channel := exprText371(call.Args[1])
	if channel != "h.Confirm != nil" {
		t.Errorf("the panel's channel term is %q, want \"h.Confirm != nil\": that is the term "+
			"perm.Store.Set passes its own channel as", channel)
	}
	// The old spelling must be gone, not merely shadowed by the new call.
	for i, line := range strings.Split(string(src), "\n") {
		if strings.Contains(line, "perm.RequiresCard(to) &&") {
			t.Errorf("%s:%d still spells the (target, channel) rule by hand: %q", path, i+1, strings.TrimSpace(line))
		}
	}
	// The panel's refusal line carries perm's classification, which is the other
	// half of "one rule, two gates".
	if !strings.Contains(string(src), "perm.RefusalResult(err)") {
		t.Errorf("%s does not classify a refusal with perm.RefusalResult, so the panel's record and the "+
			"store's MODE-SWITCH row can name one refusal differently", path)
	}
}

// exprText371 renders the one expression the pin above compares, so a failure
// prints what the gate actually passes.
func exprText371(e ast.Expr) string {
	switch v := e.(type) {
	case *ast.BinaryExpr:
		return exprText371(v.X) + " " + v.Op.String() + " " + exprText371(v.Y)
	case *ast.SelectorExpr:
		return exprText371(v.X) + "." + v.Sel.Name
	case *ast.Ident:
		return v.Name
	case *ast.CallExpr:
		// A nil literal in this AST shows up as the predeclared identifier "nil",
		// which the Ident case above already renders.
		return exprText371(v.Fun) + "(...)"
	default:
		return fmt.Sprintf("%T", e)
	}
}
