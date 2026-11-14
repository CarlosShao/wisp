package panel

// Ticket 216 AC#5's instruments. The branch point errorSentence216 / the
// dispatch refusal sentence now ride on observe.ClassOf (the D37 class pulled
// by errors.As over the chain); these tests pin the three things the ticket
// names:
//
//  1. THE CODE→SENTENCE TABLE (对照测试): one error CODE selects one leading
//     sentence - 17 classes, 17 distinct sentences, picked by the class
//     constant, never by anything found inside err.Error().
//  1b. THE WIRE: the dispatch refusal sentence (RefusedEnvelopeForUser, the
//     only page-visible error sentence this package renders for a refused
//     composer request) carries the class's own leading sentence when the
//     chain carries a D37 class, through the real config.set write route
//     whose ApplySetting errors are observe.ClassConfig wrapped (%w) upstream
//     (internal/config/sectionwrite.go), so the branch is reached through
//     production wiring, not a decoration.
//  2. THE SUBSTRING FORBIDDEN SHAPE (正控): the same error text fed under 17
//     different classes must yield 17 DIFFERENT leading sentences - a
//     substring-of-err.Error() implementation collapses those 17 to one
//     identical branch (the openchamber ChatInput.tsx:2205 "413" shape) and
//     goes red here. This is the discriminating half: an implementation that
//     branches on the error string cannot satisfy both this and test 1.
//
// The negative half AC#5 forbids is already scanned by the a2/v2 re-runs
// (zero substring-branch sites in panel production code); this file adds the
// positive branch and the instrument that keeps it that way.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/observe"
)

// TestErrorCodeSelectsTheSentence216 is the code→sentence table test. Each of
// the 17 D37 classes feeds the SAME error text; the branch must pick the
// class's own leading sentence.
func TestErrorCodeSelectsTheSentence216(t *testing.T) {
	// The same underlying text for every class - the branch key must be the
	// class, never this text.
	const sameDetail = "upstream 413 payload too large"
	for _, c := range observe.AllClasses() {
		err := observe.New(c, sameDetail)
		got := errorSentence216(err)
		want := errorCopy216[c] + err.Error()
		if got != want {
			t.Errorf("class %s: sentence = %q, want %q", c, got, want)
		}
	}
	// The table must be exhaustive over the taxonomy: an 18th class added
	// later without a copy row would fall back silently - this ruler makes
	// that red, so a class cannot drift out of the table.
	if complete := errorCopy216Complete(); !complete {
		t.Errorf("errorCopy216 is not exhaustive over observe.AllClasses()")
	}
}

// TestUnclassifiedKeepsThePlainSentence216: ClassOf ok=false keeps the
// pre-AC#5 plain sentence byte-for-byte, so the unclassified path (the shape
// every existing panel fixture feeds) does not move.
func TestUnclassifiedKeepsThePlainSentence216(t *testing.T) {
	plain := errors.New("配置读不了")
	if got := errorSentence216(plain); got != plain.Error() {
		t.Errorf("unclassified error = %q, want the plain text verbatim %q", got, plain.Error())
	}
	// A wrapped unclassified chain stays unclassified: errors.As finds no
	// observe.Error anywhere in it.
	wrapped := fmt.Errorf("panel: 配置读取失败（requestId=%q）：%w", "r-1", plain)
	if _, ok := observe.ClassOf(wrapped); ok {
		t.Fatalf("the wrapped chain suddenly classifies; the test premise broke")
	}
	if got := errorSentence216(wrapped); got != wrapped.Error() {
		t.Errorf("wrapped unclassified = %q, want verbatim %q", got, wrapped.Error())
	}
}

// TestErrorCopySentencesAreAllDistinct216 is the discriminating half: the same
// error text under 17 different classes must yield 17 DIFFERENT sentences. A
// substring-of-err.Error() implementation (the forbidden ChatInput.tsx:2205
// "413" shape) reads the same text under every class and collapses to one
// branch - this test is red the moment the branch key becomes the string.
func TestErrorCopySentencesAreAllDistinct216(t *testing.T) {
	const sameDetail = "upstream 413 payload too large"
	seen := map[observe.ErrorClass]string{}
	for _, c := range observe.AllClasses() {
		seen[c] = errorSentence216(observe.New(c, sameDetail))
	}
	for i, a := range observe.AllClasses() {
		for _, b := range observe.AllClasses()[i+1:] {
			if seen[a] == seen[b] {
				t.Errorf("classes %s and %s produce the same sentence %q - the branch no longer separates codes", a, b, seen[a])
			}
		}
	}
}

// TestRefusedEnvelopeCarriesTheClassSentence216 is the wire half: a real
// config.set write through the production dispatch, with the leg returning a
// classified error wrapped exactly as internal/config's ApplySetting wraps it
// (observe.ClassConfig + %w through fmt.Errorf), must come back to the page
// with the class's own leading sentence in the refusal.
func TestRefusedEnvelopeCarriesTheClassSentence216(t *testing.T) {
	legErr := fmt.Errorf("%w", observe.New(observe.ClassConfig, "role_chat_thinking 传入了名册外的档位"))
	leg := &cfg248Leg{applyEr: legErr}
	dispatch := cfg248Dispatch(leg)
	raw := cfg248Raw(t, map[string]any{
		"method": MethodConfigSet, "requestId": "r-216", "source": ComposerRequestSource,
		"configField": FieldRoleChatThinking, "configValue": "high",
	})
	reply, err := dispatch.Handle(context.Background(), raw)
	if err == nil {
		t.Fatalf("a refused write came back nil-error; the door did not refuse")
	}
	if leg.applyCalls != 1 {
		t.Errorf("leg apply calls = %d, want 1 - the refusal was made before the leg", leg.applyCalls)
	}
	want := errorCopy216[observe.ClassConfig]
	if !strings.Contains(reply, want) {
		t.Errorf("refusal reply %q does not carry the config-class sentence %q", reply, want)
	}
	// And the raw cause still follows - the branch adds a head, it never
	// replaces the cause (ticket 147's "读不到 without a name" rule).
	if !strings.Contains(reply, "名册外的档位") {
		t.Errorf("refusal reply %q lost the underlying cause", reply)
	}
}
