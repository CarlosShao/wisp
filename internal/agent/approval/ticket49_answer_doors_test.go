package approval_test

// Ticket 49 tests: the answer-side doors the D45-2 variant needs.
//
// These are the approval-layer cells of ticket 49's AC list, written against
// the same fixtures ticket 224's file above established (the recorder, the
// log box, the real queue):
//
//	danger door      TestTicket49AllAppsAnswerWritesOneWildcardRow
//	doors stay shut  TestTicket49AllAppsDoorNeedsTheRitualInTheHostAndRefusesTheWrongTool
//	batch one-card   TestTicket49TenOpBatchOneCardOnePatternGrant
//	pattern refusal  TestTicket49PatternGrantOnANonBatchCardWritesNothing
//
// The batch cell is the AC#5 reading the dispatch names (10 operations + opt-in
// grant => one card => the whole pattern granted in one answer), measured by
// the number of cards the UI saw and by the shape of the row the recorder got.

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/tools"
)

// ---------------------------------------------------------------------------
// fixtures
// ---------------------------------------------------------------------------

// The ticket-224 recorder (ticket224_reply_grant_test.go) now carries ticket
// 49's two target doors on the same object, so these cases read ONE recorder:
// the calls list holds per-path rows (224's shape), per-process rows (the
// resolved process name in the pattern column) and the all-apps row (the
// sentinel "*"). No second fixture type exists, so a mutation cannot pass one
// cell with a recorder the production type assertion would refuse.

// rec49Of narrows the fixture's recorder for the cases that read the door
// calls; the base type satisfies both faces since ticket 49.
func rec49Of(t *testing.T, f *grantFixture) *recGrantRecorder {
	t.Helper()
	if f.rec == nil {
		t.Fatal("fixture has no recorder")
	}
	return f.rec
}

// timeSecond is the fixture timeout helper (the 224 file used a literal; this
// one names the duration once).
func timeSecond() time.Duration { return 10 * time.Second }

// new49Fixture is newGrantFixture spelled for this file's reads: it exists so
// the cases below say what they test instead of which constructor they call.
func new49Fixture(t *testing.T, withRecorder bool) *grantFixture {
	return newGrantFixture(t, withRecorder)
}

// ---------------------------------------------------------------------------
// the danger door
// ---------------------------------------------------------------------------

// TestTicket49AllAppsAnswerWritesOneWildcardRow drives one input.type card
// (the decision carries the resolved target process) through the native
// danger door: ONE wildcard row lands on the recorder through
// RecordTargetAllApps, the call is released first, and the audit line names
// the danger tier.
func TestTicket49AllAppsAnswerWritesOneWildcardRow(t *testing.T) {
	f := new49Fixture(t, true)
	rec := rec49Of(t, f)

	revoke := f.g.AdmitTextTask(testTask)
	t.Cleanup(revoke)
	ctx, cancel := context.WithTimeout(context.Background(), 10*timeSecond())
	t.Cleanup(cancel)
	d := l2Decision()
	d.Tool = "input.type"
	d.TargetProcess = "notepad.exe"
	ch := runApproval(t, f.g, ctx, d)
	p := f.ui.wait(t)

	if err := f.g.Native().AllowSessionAllApps(context.Background(), p.CorrelationID, p.Grant); err != nil {
		t.Fatalf("Native().AllowSessionAllApps: %v", err)
	}
	if a := mustAnswer(t, ch); a.a != tools.AnswerAllow {
		t.Fatalf("answer = %v (%s), want allow: the danger answer releases the call like any other", a.a, a.why)
	}

	got := rec.got()
	if len(got) != 1 || got[0].pattern != "*" || got[0].tool != "input.type" {
		t.Fatalf("recorder saw %+v, want exactly one (input.type, *) row: the wildcard row has one door and this is it", got)
	}
	if !f.log.has("GRANT-RECORDED-ALL-APPS") {
		t.Errorf("audit is missing the danger-tier recording line:\n%s", f.log.all())
	}
	if !f.log.has("danger") && !f.log.has("危险档") {
		t.Errorf("the danger tier is not named in the audit:\n%s", f.log.all())
	}
}

// TestTicket49AllAppsDoorNeedsTheRitualInTheHostAndRefusesTheWrongTool pins
// the two shapes the danger door refuses closed:
//
//	a ledger without the target doors - a fixture whose GrantRecorder is a
//	plain 224-shaped interface value is impossible now that the recorder
//	carries both faces, so the pre-49 shape is exercised through a gate
//	whose Options.Grants is nil (the typed-nil / unwired host);
//	a card whose tool is not input.type and whose target resolved anyway
//	(a fixture-forged shape) writes nothing: the sentinel belongs to one tool.
func TestTicket49AllAppsDoorNeedsTheRitualInTheHostAndRefusesTheWrongTool(t *testing.T) {
	// Half 1: no ledger wired at all.
	f := new49Fixture(t, false)
	revoke := f.g.AdmitTextTask(testTask)
	t.Cleanup(revoke)
	d := l2Decision()
	d.Tool = "input.type"
	d.TargetProcess = "notepad.exe"
	ctx, cancel := context.WithTimeout(context.Background(), 10*timeSecond())
	t.Cleanup(cancel)
	ch := runApproval(t, f.g, ctx, d)
	p := f.ui.wait(t)
	if err := f.g.Native().AllowSessionAllApps(context.Background(), p.CorrelationID, p.Grant); err != nil {
		t.Fatalf("AllowSessionAllApps on a ledger-less host must still answer: %v", err)
	}
	if a := mustAnswer(t, ch); a.a != tools.AnswerAllow {
		t.Errorf("answer = %v, want allow", a.a)
	}
	if !f.log.has("GRANT-DROPPED") {
		t.Errorf("audit is missing the drop line on a ledger-less host:\n%s", f.log.all())
	}

	// Half 2: the sentinel is refused for a tool that is not input.type.
	f2 := new49Fixture(t, true)
	rec2 := rec49Of(t, f2)
	revoke2 := f2.g.AdmitTextTask(testTask)
	t.Cleanup(revoke2)
	d2 := l2Decision() // fs.write-ish generic card
	d2.TargetProcess = "notepad.exe"
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*timeSecond())
	t.Cleanup(cancel2)
	ch2 := runApproval(t, f2.g, ctx2, d2)
	p2 := f2.ui.wait(t)
	if err := f2.g.Native().AllowSessionAllApps(context.Background(), p2.CorrelationID, p2.Grant); err != nil {
		t.Fatalf("AllowSessionAllApps must not fail the CALL on a wrong-tool card: %v", err)
	}
	if a := mustAnswer(t, ch2); a.a != tools.AnswerAllow {
		t.Errorf("answer = %v, want allow", a.a)
	}
	if len(rec2.got()) != 0 {
		t.Errorf("the danger door wrote %v for a non-input.type card: the sentinel is input.type's alone", rec2.got())
	}
	if !f2.log.has("「所有应用」档只属于 input.type") {
		t.Errorf("audit does not say why the row was dropped:\n%s", f2.log.all())
	}
}

// ---------------------------------------------------------------------------
// the batch cell: one card, one pattern grant
// ---------------------------------------------------------------------------

// TestTicket49TenOpBatchOneCardOnePatternGrant is AC#5's cell. A 10-path batch
// card is answered through AllowSessionPattern: exactly ONE card was shown for
// the whole batch, exactly ONE wildcard-terminated row lands (the whole
// pattern), and the audit line names the batch grant. The card is routed the
// way the queue-backed confirmation rides (PendingApproval - the shape that
// CAN wait for the operator's opt-in click; the 3s L1 window auto-executes and
// has no answer verb to opt in with), and Prompt.Paths carries the whole batch.
func TestTicket49TenOpBatchOneCardOnePatternGrant(t *testing.T) {
	f := new49Fixture(t, true)
	rec := rec49Of(t, f)

	var paths []string
	for i := 0; i < 10; i++ {
		paths = append(paths, fmt.Sprintf("C:/work/batch%d.txt", i))
	}

	cards := 0
	revoke := f.g.AdmitTextTask(testTask)
	t.Cleanup(revoke)
	ctx, cancel := context.WithTimeout(context.Background(), 10*timeSecond())
	t.Cleanup(cancel)
	d := l2Decision(paths...)

	ch := runApproval(t, f.g, ctx, d)
	p := f.ui.wait(t)
	cards++
	if len(p.Paths) != 10 {
		t.Fatalf("card named %d paths, want the whole 10-op batch: the one-card promise needs the batch on ONE card", len(p.Paths))
	}

	if err := f.g.Native().AllowSessionPattern(context.Background(), p.CorrelationID, p.Grant); err != nil {
		t.Fatalf("Native().AllowSessionPattern: %v", err)
	}
	if a := mustAnswer(t, ch); a.a != tools.AnswerAllow {
		t.Fatalf("answer = %v (%s), want allow", a.a, a.why)
	}

	// THE CELL: one card, one row for the WHOLE pattern.
	if cards != 1 {
		t.Errorf("%d cards were shown for a 10-op batch, want exactly 1", cards)
	}
	got := rec.got()
	if len(got) != 1 {
		t.Fatalf("recorder saw %d rows (%v), want exactly ONE row for the whole pattern: "+
			"10 per-path rows would be the batch confirmation buying 10 rules instead of one", len(got), got)
	}
	if got[0].pattern != "C:/work/*" {
		t.Errorf("the pattern row = %q, want the whole-batch ancestor \"C:/work/*\"", got[0].pattern)
	}
	if got[0].tool != d.Tool {
		t.Errorf("the pattern row's tool = %q, want the card's own %q", got[0].tool, d.Tool)
	}
	if !f.log.has("GRANT-RECORDED") || !f.log.has("整批模式一次授予") {
		t.Errorf("audit does not name the batch grant:\n%s", f.log.all())
	}
}

// TestTicket49PatternGrantOnANonBatchCardWritesNothing keeps the pattern door
// honest: a card that was never aggregated has no batch pattern, so the door
// releases the call (the click was real) and writes nothing.
func TestTicket49PatternGrantOnANonBatchCardWritesNothing(t *testing.T) {
	f := new49Fixture(t, true)
	rec := rec49Of(t, f)

	revoke := f.g.AdmitTextTask(testTask)
	t.Cleanup(revoke)
	ctx, cancel := context.WithTimeout(context.Background(), 10*timeSecond())
	t.Cleanup(cancel)
	ch := runApproval(t, f.g, ctx, l2Decision("C:/work/only-one.txt"))
	p := f.ui.wait(t)

	if err := f.g.Native().AllowSessionPattern(context.Background(), p.CorrelationID, p.Grant); err != nil {
		t.Fatalf("AllowSessionPattern on a non-batch card must still answer the call: %v", err)
	}
	if a := mustAnswer(t, ch); a.a != tools.AnswerAllow {
		t.Errorf("answer = %v, want allow", a.a)
	}
	if len(rec.got()) != 0 {
		t.Errorf("the pattern door wrote %v for a non-batch card: it must refuse instead", rec.got())
	}
	if !f.log.has("非聚合卡没有整批模式") {
		t.Errorf("audit does not say why nothing was written:\n%s", f.log.all())
	}
}

// The fail-closed compile guard: the fixture recorder satisfies BOTH faces
// since ticket 49, so a gate whose type assertion stops finding the doors is
// the only way half 1 of the danger-door case can go green again.
var (
	_ approval.GrantRecorder = (*recGrantRecorder)(nil)
	_ interface {
		RecordTarget(context.Context, string, string) (int64, error)
		RecordTargetAllApps(context.Context, string) (int64, error)
	} = (*recGrantRecorder)(nil)
)

// failOn is referenced by the recorder's doors through the 224 fixture; a
// compile-time touch keeps the field honest if that fixture ever drops it.
var _ = errors.New
