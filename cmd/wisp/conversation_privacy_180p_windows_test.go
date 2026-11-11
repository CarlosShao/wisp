//go:build windows

package main

// The deferred-disposition 20261110 contract item C2 tests: voice
// .conversation_mode=true's first-enable L2 privacy confirm.
//
// THE REVERSE CONTRACT (摘喂线 -> 红): the gate's four outcomes are each
// asserted against the thing that produces them, so gutting any half reddens
// its own case - remove the marker write and the granted case's
// second-build-no-re-ask reading reds; remove the deny branch and the denied
// case's plan-flip red; remove the nil fail-closed and the no-surface case
// reds; remove the scopeOptions assignment and the EnableConversation case
// reds (session's own nil hook refuses first, loudly).

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/session"
)

// conversation180PStubConfirm is a card surface the cases control: it records
// that it was asked (the "asked once, not twice" reading) and answers from a
// table.
func conversation180PStubConfirm(answers *[]bool) func() bool {
	return func() bool {
		if len(*answers) == 0 {
			return false
		}
		a := (*answers)[0]
		*answers = (*answers)[1:]
		return a
	}
}

// Test180PConversationPrivacyGateOutcomes is the gate table, each row against
// its own data dir so no row inherits another row's marker:
//
//	already-confirmed  marker present, confirm NOT called, granted
//	no-surface         nil confirm, no marker, refused (session's own direction)
//	denied             confirm says no, refused, marker NOT written
//	granted            confirm says yes, granted, marker written with a receipt
func Test180PConversationPrivacyGateOutcomes(t *testing.T) {
	mkDir := func(t *testing.T) string {
		return t.TempDir()
	}
	t.Run("already-confirmed", func(t *testing.T) {
		dir := mkDir(t)
		if err := writeConversationPrivacyMarker(dir); err != nil {
			t.Fatalf("seed marker: %v", err)
		}
		asked := false
		granted, outcome := conversationPrivacyGate180P(dir, func() bool { asked = true; return false })
		if !granted || outcome != convPrivacyAlready {
			t.Fatalf("a recorded grant re-asked the card (granted=%v outcome=%s): the marker is the「一次」", granted, outcome)
		}
		if asked {
			t.Error("the card surface was invoked on an already-confirmed dir")
		}
	})
	t.Run("no-surface-fail-closed", func(t *testing.T) {
		dir := mkDir(t)
		granted, outcome := conversationPrivacyGate180P(dir, nil)
		if granted || outcome != convPrivacyNoSurface {
			t.Fatalf("a nil card surface did not fail closed (granted=%v outcome=%s): the same direction session.EnableConversation takes for a nil hook", granted, outcome)
		}
		if conversationPrivacyConfirmed(dir) {
			t.Error("a refusal wrote a consent marker")
		}
	})
	t.Run("denied", func(t *testing.T) {
		dir := mkDir(t)
		granted, outcome := conversationPrivacyGate180P(dir, func() bool { return false })
		if granted || outcome != convPrivacyDenied {
			t.Fatalf("a denied card granted the mode (granted=%v outcome=%s)", granted, outcome)
		}
		if conversationPrivacyConfirmed(dir) {
			t.Error("a denial wrote a consent marker: the next boot would never ask again")
		}
	})
	t.Run("granted", func(t *testing.T) {
		dir := mkDir(t)
		granted, outcome := conversationPrivacyGate180P(dir, func() bool { return true })
		if !granted || outcome != convPrivacyGranted {
			t.Fatalf("an explicit allow did not grant (granted=%v outcome=%s)", granted, outcome)
		}
		body, err := os.ReadFile(conversationPrivacyMarkerPath(dir))
		if err != nil {
			t.Fatalf("a grant wrote no marker: %v", err)
		}
		for _, want := range []string{"key=voice.conversation_mode", "level=L2", "granted_at="} {
			if !strings.Contains(string(body), want) {
				t.Errorf("marker receipt lacks %q (body=%q): the file must say what was confirmed and when", want, string(body))
			}
		}
	})
}

// Test180PConversationPrivacyGateFlipsTheMachineSelector is the machine-side
// face: a residency whose plan carries the mode true leaves applyConversation
// PrivacyGate180P with the POST-gate value, and the residency carries the two
// values the scope's own hook will need. The denied row is the reverse half:
// gut the deny branch and the mode would ride row #28 unconfirmed, which is
// exactly the shape PLAN.md:2211 forbids.
func Test180PConversationPrivacyGateFlipsTheMachineSelector(t *testing.T) {
	for _, tc := range []struct {
		name        string
		answers     []bool
		haveCard    bool
		wantMode    bool
		wantOutcome conversationPrivacyOutcome
	}{
		{"granted", []bool{true}, true, true, convPrivacyGranted},
		{"denied", []bool{false}, true, false, convPrivacyDenied},
		{"no-surface", nil, false, false, convPrivacyNoSurface},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			rv := &residentVoice{plan: voicePlan{ConversationMode: true}, residency: newVoiceResidency(voicePlan{ConversationMode: true}, newVoiceOwnerQueue(), nil)}
			answers := tc.answers
			var confirm func() bool
			if tc.haveCard {
				confirm = conversation180PStubConfirm(&answers)
			}
			rv.applyConversationPrivacyGate180P(dir, confirm)
			if rv.plan.ConversationMode != tc.wantMode {
				t.Fatalf("%s: plan.ConversationMode = %v, want %v (the selector must carry the post-gate value)", tc.name, rv.plan.ConversationMode, tc.wantMode)
			}
			if got := rv.residency.conversationPrivacyDataDir; got != dir {
				t.Errorf("%s: the residency's data dir was not set (got %q): the scope face would read a different marker", tc.name, got)
			}
			if (rv.residency.conversationPrivacyConfirm == nil) != !tc.haveCard {
				t.Errorf("%s: the residency's card surface was not carried for the scope face", tc.name)
			}
			if conversationPrivacyConfirmed(dir) != (tc.wantOutcome == convPrivacyGranted) {
				t.Errorf("%s: marker presence disagrees with the outcome %s", tc.name, tc.wantOutcome)
			}
		})
	}
}

// Test180PConversationPrivacyGateAsksOncePerDataDir is the「一次」contract: the
// first build asks and grants; the second build on the same dir needs NO card
// surface at all (nil confirm) and still runs the mode. A process-local
// confirmed flag would re-ask every boot - that red is the point of the case.
func Test180PConversationPrivacyGateAsksOncePerDataDir(t *testing.T) {
	dir := t.TempDir()
	rv := &residentVoice{plan: voicePlan{ConversationMode: true}, residency: newVoiceResidency(voicePlan{ConversationMode: true}, newVoiceOwnerQueue(), nil)}
	answers := []bool{true}
	rv.applyConversationPrivacyGate180P(dir, conversation180PStubConfirm(&answers))
	if len(answers) != 0 || !rv.plan.ConversationMode {
		t.Fatalf("the first enable did not ask-and-grant (answers left=%d mode=%v)", len(answers), rv.plan.ConversationMode)
	}
	rv2 := &residentVoice{plan: voicePlan{ConversationMode: true}, residency: newVoiceResidency(voicePlan{ConversationMode: true}, newVoiceOwnerQueue(), nil)}
	rv2.applyConversationPrivacyGate180P(dir, nil)
	if !rv2.plan.ConversationMode {
		t.Fatal("the second boot re-asked (nil surface refused a recorded grant): the marker is not doing its「一次」job")
	}
	if got := rv2.residency.conversationPrivacyDataDir; got != dir {
		t.Errorf("the second boot did not carry the dir into the residency (got %q)", got)
	}
}

// Test180PScopeOptionsPrivacyConfirmRidesTheSharedGate is the scope face: the
// residency's scopeOptions must assign Hooks.PrivacyConfirm from the same gate
// the machine face ran. Three reads, one case:
//
//	with a card and no marker  -> the first EnableConversation asks and grants;
//	                              the marker lands (the scope's「一次」persisted);
//	with the marker present    -> a nil surface still enables (already confirmed);
//	with no card and no marker -> EnableConversation refuses (fail-closed, the
//	                              session package's own nil-hook direction).
func Test180PScopeOptionsPrivacyConfirmRidesTheSharedGate(t *testing.T) {
	newScope := func(t *testing.T, r *voiceResidency) *session.Scope {
		t.Helper()
		opts := r.scopeOptions()
		opts.Clock = &voice423InstantClock{}
		opts.Loader = func(session.EngineKind) (session.Engine, error) { return nil, errors.New("180P case arms no engines") }
		sc, err := session.NewScope(opts)
		if err != nil {
			t.Fatalf("NewScope: %v", err)
		}
		return sc
	}
	t.Run("first-enable-asks-and-persists", func(t *testing.T) {
		dir := t.TempDir()
		r := newVoiceResidency(voicePlan{}, newVoiceOwnerQueue(), nil)
		r.conversationPrivacyDataDir = dir
		answers := []bool{true}
		r.conversationPrivacyConfirm = conversation180PStubConfirm(&answers)
		sc := newScope(t, r)
		ok, err := sc.EnableConversation()
		if err != nil || !ok {
			t.Fatalf("the scope's first enable was refused through the shared gate (ok=%v): %v", ok, err)
		}
		if len(answers) != 0 {
			t.Error("the card was never asked: the hook is not riding the shared gate")
		}
		if !conversationPrivacyConfirmed(dir) {
			t.Error("the grant left no marker: every later boot would re-ask")
		}
	})
	t.Run("recorded-grant-needs-no-card", func(t *testing.T) {
		dir := t.TempDir()
		if err := writeConversationPrivacyMarker(dir); err != nil {
			t.Fatalf("seed marker: %v", err)
		}
		r := newVoiceResidency(voicePlan{}, newVoiceOwnerQueue(), nil)
		r.conversationPrivacyDataDir = dir
		r.conversationPrivacyConfirm = nil
		sc := newScope(t, r)
		if ok, err := sc.EnableConversation(); err != nil || !ok {
			t.Fatalf("a recorded grant was refused on a nil surface (ok=%v): %v", ok, err)
		}
	})
	t.Run("no-card-no-marker-fails-closed", func(t *testing.T) {
		dir := t.TempDir()
		r := newVoiceResidency(voicePlan{}, newVoiceOwnerQueue(), nil)
		r.conversationPrivacyDataDir = dir
		r.conversationPrivacyConfirm = nil
		sc := newScope(t, r)
		ok, _ := sc.EnableConversation()
		// The refused direction is session's own (scope.go:545-548): a card the
		// hook does not grant answers (false, nil) - "仍是 OFF" is a verdict,
		// not an error. The fail-closed reading is ok=false with no marker.
		if ok {
			t.Fatal("a mode that keeps a mic open was enabled with no card surface and no marker")
		}
		if conversationPrivacyConfirmed(dir) {
			t.Error("the refusal wrote a marker")
		}
	})
}

// conversation180PMarkerDeleteReArms is the reverse of「一次」, pinned so the
// marker stays deletable: removing the file puts the gate back into ask mode.
// If somebody swaps the marker for an un-deletable store, this case reds.
func Test180PConversationPrivacyMarkerDeleteReArms(t *testing.T) {
	dir := t.TempDir()
	if err := writeConversationPrivacyMarker(dir); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := os.Remove(conversationPrivacyMarkerPath(dir)); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if conversationPrivacyConfirmed(dir) {
		t.Fatal("a deleted marker still reads as confirmed")
	}
	granted, outcome := conversationPrivacyGate180P(dir, func() bool { return false })
	if granted || outcome != convPrivacyDenied {
		t.Fatalf("a re-armed dir did not ask again (granted=%v outcome=%s)", granted, outcome)
	}
}

// compile-time: the instant clock type the scope cases borrow lives in the 423
// file; this keeps the borrow honest if that helper ever moves.
var _ = voice423InstantClock{}
