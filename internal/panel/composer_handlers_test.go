package panel

// Ticket 114 AC#2 - the native gate on a panel mode write (R-92-2).
//
// The judgement this file has to be able to fail is not "the handler returned an
// error". It is "the injected write leg was never reached": a handler that logs
// a refusal and then calls Set anyway is green under a string check and red
// under these cases, which is exactly the mutation AC#2 names. internal/panel is
// in scripts/portable-tests.sh's core scope, so all of this runs on the ubuntu
// leg as well as here; nothing about the gate is delegated to cmd/wisp, whose
// tests only the windows leg runs.

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/risk"
)

// fakeModeWriter records every call the handler makes. setCalls is the number
// AC#2 turns on: a refusal that reached Set is a refusal that wrote the档.
type fakeModeWriter struct {
	current   risk.Mode
	currentOK bool // false = PermissionMode() answers something off the ladder
	setCalls  int
	setMode   risk.Mode
	setOrigin string
	setActor  string
	setErr    error
	readCalls int
}

func (f *fakeModeWriter) PermissionMode() risk.Mode {
	f.readCalls++
	if !f.currentOK {
		// A writer that invents a档 outside the three: not the same as asking
		// "which档", it is the shape that would let widening look like
		// narrowing and close this gate from the outside.
		return risk.Mode(97)
	}
	return f.current
}

func (f *fakeModeWriter) Set(_ context.Context, to risk.Mode, origin, actor string) error {
	f.setCalls++
	f.setMode, f.setOrigin, f.setActor = to, origin, actor
	return f.setErr
}

// countingConfirm is the attached leg. It counts, because the single-card rule
// says the handler must never call it.
type countingConfirm struct{ calls int }

func (c *countingConfirm) fn() ModeConfirm {
	return func(context.Context, risk.Mode, risk.Mode, string, string) error {
		c.calls++
		return nil
	}
}

func modeReq(to, id string) ComposerRequest {
	return ComposerRequest{Method: MethodModeRequest, RequestID: id, Source: ComposerRequestSource, To: to}
}

// modeLadder is the order the gate's one comparison depends on.
func modeLadder(t *testing.T) []risk.Mode {
	t.Helper()
	names := risk.ModeNames()
	ladder := make([]risk.Mode, 0, len(names))
	for _, n := range names {
		m, err := risk.ParseMode(n)
		if err != nil {
			t.Fatalf("risk.ParseMode(%q): %v", n, err)
		}
		ladder = append(ladder, m)
	}
	return ladder
}

// TestTheModeLadderThisGateJudgesAgainstIsTheDocumentedThree nails the premise
// of modeIsWidening, which is a raw `to > from` over risk.Mode's constants. The
// vocabulary is a list, so it is nailed three ways: exact equality (a rename or
// a reorder goes red), a length floor and a strict-prefix negative assertion
// (a SHRUNKEN list goes red on its own line, not only inside the equality -
// A109(2) measured that a 6-to-5 trim leaves per-entry assertions green).
func TestTheModeLadderThisGateJudgesAgainstIsTheDocumentedThree(t *testing.T) {
	names := risk.ModeNames()
	want := []string{risk.ModeAskEveryStepName, risk.ModeAskHighRiskName, risk.ModeAutoApproveName}
	if len(names) < len(want) {
		t.Fatalf("the mode vocabulary holds %d names, this gate needs at least %d: %v",
			len(names), len(want), names)
	}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("the mode ladder this gate compares against moved:\n got %v\nwant %v", names, want)
	}
	for n := 1; n < len(want); n++ {
		prefix := want[:n]
		if reflect.DeepEqual(names, prefix) {
			t.Fatalf("the vocabulary collapsed to the strict prefix %v: widening can no longer be recognised", prefix)
		}
	}
	ladder := modeLadder(t)
	for i := 1; i < len(ladder); i++ {
		if !modeIsWidening(ladder[i-1], ladder[i]) {
			t.Fatalf("%s -> %s must read as widening; risk.Mode's constants are ordered strict -> loose on purpose",
				names[i-1], names[i])
		}
		if modeIsWidening(ladder[i], ladder[i-1]) {
			t.Fatalf("%s -> %s must read as narrowing (R20: getting stricter is never blocked)",
				names[i], names[i-1])
		}
	}
}

// TestAWideningModeRequestIsRefusedWhenNoL2ConfirmLegIsAttached is AC#2.
func TestAWideningModeRequestIsRefusedWhenNoL2ConfirmLegIsAttached(t *testing.T) {
	ladder := modeLadder(t)
	names := risk.ModeNames()
	cases := 0
	refused := 0
	for fromIdx, from := range ladder {
		for toIdx := range ladder {
			cases++
			if toIdx <= fromIdx {
				continue // narrowing or no-change: this case is about widening
			}
			w := &fakeModeWriter{current: from, currentOK: true}
			var audit auditLog
			h := &ModeWriteHandler{Modes: w, Audit: audit.f}
			id := fmt.Sprintf("rid-%s-%s", names[fromIdx], names[toIdx])
			err := h.HandleModeRequest(context.Background(), modeReq(names[toIdx], id))
			if err == nil {
				t.Fatalf("%s -> %s with no confirm leg: got nil error, the档 would have been written",
					names[fromIdx], names[toIdx])
			}
			if !errors.Is(err, ErrNoL2Confirm) {
				t.Fatalf("%s -> %s: error %v does not say the confirmation leg is what is missing",
					names[fromIdx], names[toIdx], err)
			}
			if w.setCalls != 0 {
				t.Fatalf("%s -> %s: the injected ModeWriter was called %d time(s) despite the refusal - the gate did not stop the write",
					names[fromIdx], names[toIdx], w.setCalls)
			}
			if len(audit) != 1 {
				t.Fatalf("%s -> %s: refusal wrote %d audit line(s), want exactly 1 (a refused switch must not be silent)",
					names[fromIdx], names[toIdx], len(audit))
			}
			if !strings.Contains(audit.joined(), id) {
				t.Fatalf("audit line %q does not name the requestId it refused", audit.joined())
			}
			refused++
		}
	}
	if cases != len(ladder)*len(ladder) {
		t.Fatalf("the switch matrix covered %d pairs of %d modes (%d): the enumeration was shortened",
			cases, len(ladder), len(ladder)*len(ladder))
	}
	if want := len(ladder) * (len(ladder) - 1) / 2; refused != want {
		t.Fatalf("%d widening pairs asserted, want %d: some direction of the ladder went untested", refused, want)
	}
}

// TestAWideningRequestWithAnAttachedLegReachesTheOneWriterAndRaisesNoSecondCard
// is the other half of the same property: the leg being there is what opens the
// gate, and opening it means ONE write through ONE leg with no second card.
func TestAWideningRequestWithAnAttachedLegReachesTheOneWriterAndRaisesNoSecondCard(t *testing.T) {
	ladder := modeLadder(t)
	names := risk.ModeNames()
	for fromIdx, from := range ladder {
		for toIdx, to := range ladder {
			if toIdx <= fromIdx {
				continue
			}
			w := &fakeModeWriter{current: from, currentOK: true}
			card := &countingConfirm{}
			var audit auditLog
			h := &ModeWriteHandler{Modes: w, Confirm: card.fn(), Audit: audit.f, Actor: "sess-1"}
			if err := h.HandleModeRequest(context.Background(), modeReq(names[toIdx], "rid-ok")); err != nil {
				t.Fatalf("%s -> %s with a leg attached: %v", names[fromIdx], names[toIdx], err)
			}
			if w.setCalls != 1 {
				t.Fatalf("%s -> %s: ModeWriter called %d time(s), want 1", names[fromIdx], names[toIdx], w.setCalls)
			}
			if w.setMode != to {
				t.Fatalf("%s -> %s: the write leg was asked for %s", names[fromIdx], names[toIdx], w.setMode)
			}
			if w.setOrigin != PanelModeOrigin {
				t.Fatalf("origin = %q, want %q (perm's audit vocabulary)", w.setOrigin, PanelModeOrigin)
			}
			if w.setActor != "sess-1" {
				t.Fatalf("actor = %q, want the injected %q", w.setActor, "sess-1")
			}
			if card.calls != 0 {
				t.Fatalf("%s -> %s: the handler raised its own confirmation (%d call(s)) on top of the one the write leg owns - two cards for one request",
					names[fromIdx], names[toIdx], card.calls)
			}
			if len(audit) != 0 {
				t.Fatalf("a request that reached the write leg must not add a second verdict, got %q", audit.joined())
			}
		}
	}
}

// TestAStricterOrEqualModeRequestNeedsNoConfirmLeg keeps the gate from being a
// lock: R20 makes getting stricter free, and perm.Set is where a no-change
// request is audited. Refusing these would be a gate that protects nothing.
func TestAStricterOrEqualModeRequestNeedsNoConfirmLeg(t *testing.T) {
	ladder := modeLadder(t)
	names := risk.ModeNames()
	passed := 0
	for fromIdx, from := range ladder {
		for toIdx := range ladder {
			if toIdx >= fromIdx {
				continue
			}
			w := &fakeModeWriter{current: from, currentOK: true}
			h := &ModeWriteHandler{Modes: w}
			if err := h.HandleModeRequest(context.Background(), modeReq(names[toIdx], "rid-narrow")); err != nil {
				t.Fatalf("%s -> %s (stricter) was refused: %v", names[fromIdx], names[toIdx], err)
			}
			if w.setCalls != 1 {
				t.Fatalf("%s -> %s (stricter): ModeWriter called %d time(s), want 1 - a legless host must still be able to get safer",
					names[fromIdx], names[toIdx], w.setCalls)
			}
			passed++
		}
	}
	if want := len(ladder) * (len(ladder) - 1) / 2; passed != want {
		t.Fatalf("%d narrowing pairs reached the write leg, want %d", passed, want)
	}
	// The no-change row belongs to the write leg too: perm audits it as
	// no-change, and this gate must not eat that line either.
	for idx, m := range ladder {
		w := &fakeModeWriter{current: m, currentOK: true}
		h := &ModeWriteHandler{Modes: w}
		if err := h.HandleModeRequest(context.Background(), modeReq(names[idx], "rid-same")); err != nil {
			t.Fatalf("%s -> %s (unchanged) was refused: %v", names[idx], names[idx], err)
		}
		if w.setCalls != 1 {
			t.Fatalf("%s unchanged: ModeWriter called %d time(s), want 1 (R20 audits all three档)", names[idx], w.setCalls)
		}
	}
}

// TestAModeRequestIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete covers the
// three ways the handler can be handed nothing to write with, none of which may
// reach a write leg.
func TestAModeRequestIsRefusedBeforeAnyWriteWhenTheAssemblyIsIncomplete(t *testing.T) {
	t.Run("no write leg at all", func(t *testing.T) {
		var audit auditLog
		h := &ModeWriteHandler{Audit: audit.f}
		err := h.HandleModeRequest(context.Background(), modeReq(risk.ModeAutoApproveName, "rid-nowriter"))
		if !errors.Is(err, ErrNoModeWriter) {
			t.Fatalf("got %v, want %v", err, ErrNoModeWriter)
		}
		if len(audit) != 1 {
			t.Fatalf("audit lines = %d, want 1", len(audit))
		}
	})
	t.Run("current档 unreadable", func(t *testing.T) {
		w := &fakeModeWriter{currentOK: false}
		var audit auditLog
		h := &ModeWriteHandler{Modes: w, Audit: audit.f}
		err := h.HandleModeRequest(context.Background(), modeReq(risk.ModeAskHighRiskName, "rid-unreadable"))
		if !errors.Is(err, ErrNoCurrentMode) {
			t.Fatalf("got %v, want %v", err, ErrNoCurrentMode)
		}
		if w.setCalls != 0 {
			t.Fatalf("ModeWriter called %d time(s) on an unreadable档, want 0", w.setCalls)
		}
		if len(audit) != 1 {
			t.Fatalf("audit lines = %d, want 1", len(audit))
		}
	})
	t.Run("unknown spelling", func(t *testing.T) {
		w := &fakeModeWriter{current: risk.ModeAskEveryStep, currentOK: true}
		h := &ModeWriteHandler{Modes: w}
		if err := h.HandleModeRequest(context.Background(), modeReq("off", "rid-bad")); err == nil {
			t.Fatal("an unknown档 name was accepted")
		}
		if w.setCalls != 0 {
			t.Fatalf("ModeWriter called %d time(s) for an unknown档, want 0", w.setCalls)
		}
	})
	t.Run("empty target", func(t *testing.T) {
		w := &fakeModeWriter{current: risk.ModeAskEveryStep, currentOK: true}
		h := &ModeWriteHandler{Modes: w}
		if err := h.HandleModeRequest(context.Background(), modeReq("", "rid-empty")); err == nil {
			t.Fatal("a request with no target was accepted")
		}
		if w.setCalls != 0 {
			t.Fatalf("ModeWriter called %d time(s) for an empty target, want 0", w.setCalls)
		}
	})
	t.Run("another method", func(t *testing.T) {
		w := &fakeModeWriter{current: risk.ModeAskEveryStep, currentOK: true}
		h := &ModeWriteHandler{Modes: w}
		err := h.HandleModeRequest(context.Background(),
			ComposerRequest{Method: MethodMessageSend, RequestID: "rid-route", To: risk.ModeAutoApproveName})
		if err == nil {
			t.Fatal("the mode handler answered a message.send request")
		}
		if w.setCalls != 0 {
			t.Fatalf("ModeWriter called %d time(s) for a misrouted request, want 0", w.setCalls)
		}
	})
}

// TestTheWriteLegsOwnRefusalIsReturnedUntouched: once the gate opens, perm's
// verdict is the answer - it is the record R20 asks for, and rewriting it here
// would be a second, softer account of the same switch.
func TestTheWriteLegsOwnRefusalIsReturnedUntouched(t *testing.T) {
	sentinel := errors.New("perm: 切到 auto_approve 被拒绝: 审批超时")
	w := &fakeModeWriter{current: risk.ModeAskEveryStep, currentOK: true, setErr: sentinel}
	h := &ModeWriteHandler{Modes: w, Confirm: (&countingConfirm{}).fn()}
	err := h.HandleModeRequest(context.Background(), modeReq(risk.ModeAutoApproveName, "rid-timeout"))
	if !errors.Is(err, sentinel) {
		t.Fatalf("got %v, want the write leg's own error %v", err, sentinel)
	}
	if w.setCalls != 1 {
		t.Fatalf("ModeWriter called %d time(s), want 1", w.setCalls)
	}
}
