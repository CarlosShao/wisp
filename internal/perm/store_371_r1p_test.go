package perm

// Ticket 371, AC#2' and AC#2'-b: the gate must require an ANSWERED allow, not a
// non-nil channel.
//
// WHAT WAS WRONG (measured by 371-a1p on 786a13b22, docs/evidence/s1/371-a1p-impl.md
// section 5): auto_approve + a NON-NIL limb that never raises a card read
// invoked=1 / asked=0 / applied, i.e. the switch applied with nobody asked. The
// gate's channel term was nil-ness only; nothing judged "was the card answered",
// even though that evidence existed further down (run.go's ans==AnswerAllow,
// approval/queue.go's ANSWER-ALLOW line) and was flattened to `error == nil` here.
//
// WHAT THIS FILE PINS, in the ticket's own acceptance order:
//
//	1. AA + non-nil, never-asking limb  => applied=0 AND the refusal NAMES
//	   "card was not answered" (a bare "refused" would be indistinguishable from
//	   the nil-channel branch, which is the whole reason for a separate verdict).
//	2. POSITIVE CONTROL: AA + a limb that really answered allow => applied=1,
//	   asked=1, and the audit row carries asked=1 with the card's correlation id.
//	3. The two neighbouring verdicts did NOT move: a nil channel is still the
//	   older refused-confirm story, and an operator's "no" is still a refusal.
//
// The reverse tooth (remove the answered hop => case 1 reads applied=1 and silent)
// is measured in docs/evidence/s1/371-r1p-impl.md, not here: a test cannot delete
// the line it is judging.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// leg371 is one confirmation limb with the three things these cases distinguish:
// whether it was entered, whether it really raised a card, and what it reported.
// The limb literals below each set those three facts explicitly, so a case can
// never pass by borrowing another case's fixture.
type leg371 struct {
	invoked int
	asked   int
}

// neverAsks371 is AC#0'②'s shape with 371's answer contract: non-nil, entered,
// and it returns the ZERO credential - no card was ever shown.
func neverAsks371(l *leg371) ConfirmFunc {
	return func(context.Context, Switch) (Card, error) {
		l.invoked++
		return Card{}, nil
	}
}

// answers371 is the positive control: one card is raised and the operator answered
// allow, which is exactly what the production limb reports after
// approval.Gate answered tools.AnswerAllow (cmd/wisp/run.go's
// confirmModeSwitchVia).
func answers371(l *leg371, corr string) ConfirmFunc {
	return func(context.Context, Switch) (Card, error) {
		l.invoked++
		l.asked++
		return AnsweredCard(corr), nil
	}
}

// cardLine371 returns the MODE-SWITCH row a recorder holds for to=auto_approve.
func cardLine371(t *testing.T, rec *recorder, to string) string {
	t.Helper()
	var hit string
	for _, l := range rec.all() {
		if strings.Contains(l, "MODE-SWITCH") && strings.Contains(l, "to="+to) {
			hit = l
		}
	}
	if hit == "" {
		t.Fatalf("no MODE-SWITCH row for to=%s in:\n%s", to, strings.Join(rec.all(), "\n"))
	}
	return hit
}

// Test371R1PTheNeverAskingLimbIsRefusedAndTheReasonNamesIt is the ticket's single
// acceptance reading: applied=0 AND the refusal says the card was not answered.
func Test371R1PTheNeverAskingLimbIsRefusedAndTheReasonNamesIt(t *testing.T) {
	ctx := context.Background()
	limb := &leg371{}
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, neverAsks371(limb))

	err := s.Set(ctx, risk.ModeAutoApprove, "t371", "owner")

	// The limb WAS entered: this is "asked and not answered", not "never asked".
	if limb.invoked != 1 {
		t.Fatalf("the limb was entered %d time(s), want 1 - the gate still calls it for auto_approve", limb.invoked)
	}
	// 1. applied=0.
	if err == nil {
		t.Fatal("the switch applied although no card was ever answered: ticket 371's hole is open")
	}
	if mode, saves := mc.saved(); mode != risk.ModeAskEveryStepName || saves != 0 {
		t.Errorf("persisted mode=%q saves=%d, want ask_every_step/0 (nothing may land)", mode, saves)
	}
	// 2. The refusal NAMES the unanswered card, in its own vocabulary...
	if !errors.Is(err, ErrCardNotAnswered) {
		t.Errorf("refusal %v does not wrap ErrCardNotAnswered", err)
	}
	if !strings.Contains(err.Error(), "card was not answered") {
		t.Errorf("refusal %q does not name \"card was not answered\"", err.Error())
	}
	// ... and it is NOT the nil-channel refusal, which is the distinction the
	// ticket calls out as the one that would otherwise be lost.
	if strings.Contains(err.Error(), "没有接入确认通道") {
		t.Errorf("the never-asking limb was refused with the nil-channel sentence, so the two stories "+
			"are indistinguishable again: %q", err.Error())
	}
	// 3. The audit row says the same thing, and says nobody was asked.
	last, ok := s.Last()
	if !ok {
		t.Fatal("the refusal left no record")
	}
	if last.Result != ResultCardNotAnswered {
		t.Errorf("audit result = %q, want %q", last.Result, ResultCardNotAnswered)
	}
	if last.Asked || last.Card != "" {
		t.Errorf("audit asked=%v card=%q, want false/\"\"", last.Asked, last.Card)
	}
	line := cardLine371(t, rec, risk.ModeAutoApproveName)
	for _, needle := range []string{"result=" + string(ResultCardNotAnswered), "asked=0", `card=""`} {
		if !strings.Contains(line, needle) {
			t.Errorf("the MODE-SWITCH row does not carry %q: %s", needle, line)
		}
	}
}

// Test371R1PTheAnsweredLimbAppliesAndTheRowCarriesTheCard is the positive
// control, and it is also AC#2'-b's 〔丙-b〕 half: the row carries asked=1 and the
// correlation id of the card that was answered, which is the value the approval
// layer prints on its own ANSWER-ALLOW line.
func Test371R1PTheAnsweredLimbAppliesAndTheRowCarriesTheCard(t *testing.T) {
	ctx := context.Background()
	limb := &leg371{}
	const corr = "host:mode-switch#7"
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, answers371(limb, corr))

	if err := s.Set(ctx, risk.ModeAutoApprove, "panel", "owner"); err != nil {
		t.Fatalf("an answered allow was refused: %v", err)
	}
	if limb.invoked != 1 || limb.asked != 1 {
		t.Errorf("invoked=%d asked=%d, want 1/1: one card, answered once (R20/M4)", limb.invoked, limb.asked)
	}
	if mode, saves := mc.saved(); mode != risk.ModeAutoApproveName || saves != 1 {
		t.Errorf("persisted mode=%q saves=%d, want auto_approve/1", mode, saves)
	}
	last, ok := s.Last()
	if !ok {
		t.Fatal("the applied switch left no record")
	}
	if last.Result != ResultApplied {
		t.Errorf("audit result = %q, want %q", last.Result, ResultApplied)
	}
	if !last.Asked || last.Card != corr {
		t.Errorf("audit asked=%v card=%q, want true/%q: the row must name the card that was answered",
			last.Asked, last.Card, corr)
	}
	line := cardLine371(t, rec, risk.ModeAutoApproveName)
	for _, needle := range []string{"asked=1", `card="` + corr + `"`, "result=" + string(ResultApplied)} {
		if !strings.Contains(line, needle) {
			t.Errorf("the MODE-SWITCH row does not carry %q: %s", needle, line)
		}
	}
}

// Test371R1PTheNilChannelRefusalIsStillItsOwnVerdict keeps the neighbouring story
// where ticket 372 left it: no channel at all is refused-confirm, names the
// channel, and is NOT the unanswered-card verdict.
func Test371R1PTheNilChannelRefusalIsStillItsOwnVerdict(t *testing.T) {
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, nil)
	err := s.Set(context.Background(), risk.ModeAutoApprove, "panel", "owner")
	if err == nil {
		t.Fatal("auto_approve was reached with no confirmation channel")
	}
	if errors.Is(err, ErrCardNotAnswered) {
		t.Errorf("the nil-channel refusal claims a card was not answered: %v", err)
	}
	if !strings.Contains(err.Error(), "没有接入确认通道") {
		t.Errorf("the nil-channel refusal does not name the channel: %v", err)
	}
	if last, _ := s.Last(); last.Result != ResultRefused {
		t.Errorf("audit result = %q, want %q", last.Result, ResultRefused)
	}
	if _, saves := mc.saved(); saves != 0 {
		t.Errorf("state moved under the refusal: saves=%d", saves)
	}
	line := cardLine371(t, rec, risk.ModeAutoApproveName)
	if !strings.Contains(line, "asked=0") {
		t.Errorf("the refusal row does not say nobody was asked: %s", line)
	}
}

// Test371R1PAnOperatorsNoIsStillARefusalNotAnUnansweredCard keeps the third
// verdict apart: a limb that reports an error said something, and that is not the
// "no card was answered" story (the operator DID answer - no).
func Test371R1PAnOperatorsNoIsStillARefusalNotAnUnansweredCard(t *testing.T) {
	denied := errors.New("用户拒绝")
	limb := &leg371{}
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, func(context.Context, Switch) (Card, error) {
		limb.invoked++
		return Card{}, denied
	})
	err := s.Set(context.Background(), risk.ModeAutoApprove, "ball", "owner")
	if err == nil {
		t.Fatal("a denied card still applied the switch")
	}
	if !errors.Is(err, denied) {
		t.Errorf("refusal %v lost the operator's own reason", err)
	}
	if errors.Is(err, ErrCardNotAnswered) {
		t.Errorf("a denial is booked as an unanswered card: %v", err)
	}
	if last, _ := s.Last(); last.Result != ResultRefused {
		t.Errorf("audit result = %q, want %q", last.Result, ResultRefused)
	}
	if _, saves := mc.saved(); saves != 0 {
		t.Errorf("a denied card persisted (%d writes)", saves)
	}
	if got := rec.count(string(ResultCardNotAnswered)); got != 0 {
		t.Errorf("a denial wrote %d unanswered-card rows, want 0: %v", got, rec.all())
	}
}

// Test371R1PTheMiddleTargetIsStillFreeOfCards is R20/M4's asymmetry, unmoved by
// this ticket: ask_every_step -> ask_high_risk costs no card at all, so the new
// answered rule cannot have leaked onto it.
func Test371R1PTheMiddleTargetIsStillFreeOfCards(t *testing.T) {
	limb := &leg371{}
	s, mc, rec := newTestStore(t, risk.ModeAskEveryStepName, neverAsks371(limb))
	if err := s.Set(context.Background(), risk.ModeAskHighRisk, "panel", "owner"); err != nil {
		t.Fatalf("the middle target was refused: %v", err)
	}
	if limb.invoked != 0 {
		t.Errorf("ask_high_risk entered the confirmation limb %d time(s); R20/M4 cards only auto_approve", limb.invoked)
	}
	if mode, saves := mc.saved(); mode != risk.ModeAskHighRiskName || saves != 1 {
		t.Errorf("persisted mode=%q saves=%d, want ask_high_risk/1", mode, saves)
	}
	line := cardLine371(t, rec, risk.ModeAskHighRiskName)
	if !strings.Contains(line, "asked=0") {
		t.Errorf("a target that costs no card claims somebody was asked: %s", line)
	}
}

// Test371R1PRefusalResultIsTheOneVocabulary pins the mapping both transition
// points share (AC#2'-b): the panel books the store's refusal with these values,
// so a change here moves both layers at once instead of only one.
func Test371R1PRefusalResultIsTheOneVocabulary(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want Result
	}{
		{"no error", nil, ResultApplied},
		{"unanswered card", ErrCardNotAnswered, ResultCardNotAnswered},
		{
			"wrapped unanswered card",
			fmt.Errorf("perm: 切到 auto_approve 被拒绝: %w", ErrCardNotAnswered), ResultCardNotAnswered,
		},
		{"operator said no", errors.New("用户拒绝"), ResultRefused},
		{"deadline", context.DeadlineExceeded, ResultConfirmError},
		{"cancelled", context.Canceled, ResultConfirmError},
	}
	for _, tc := range cases {
		if got := RefusalResult(tc.err); got != tc.want {
			t.Errorf("%s: RefusalResult(%v) = %q, want %q", tc.name, tc.err, got, tc.want)
		}
	}
}
