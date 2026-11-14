//go:build windows

package main

// Ticket 216 AC#5's instrument for the announcement path. The criterion is
// "the interrupted/failed split is decided by the error CODE, never by a
// substring of the English error text" - the negative example is openchamber
// ChatInput.tsx:2205-2206 (`includes('payload too large') ||
// includes('413') || includes('entity too large')`, transcribed verbatim in
// .scratch/wisp/probes/216/ac5/control/openchamber-chatinput-2205.txt).
//
// The arm that gives this ruler teeth is arm 3: an error whose TEXT is exactly
// the sentinel's text, but which is a DIFFERENT error, must answer false. The
// pre-216 shape (`strings.Contains(speakErr.Error(), "interrupted")`) answers
// true there and this test goes red. Restoring that shape is the mutation the
// evidence file records.
//
// What this does NOT cover: whether the real engines preserve the sentinel.
// They do by construction (MatchaTts.Speak returns the emit error verbatim,
// CloudTTS.Speak documents the same contract, CloudTTSChain wraps with %w), and
// arms 1-2 pin the two shapes that reach speak(); a live announcement would
// need the real models and device, which this ticket's write face does not
// reach.

import (
	"errors"
	"fmt"
	"testing"
)

func TestSpeakInterruptIsDecidedByCode216(t *testing.T) {
	// Arm 1: the code itself.
	if !speakFailureIsInterrupt(errSpeakInterrupted) {
		t.Errorf("the interrupt sentinel does not read as an interrupt: the code branch lost its own code")
	}

	// Arm 2: the sentinel through the wrapping shapes the speakers use. The
	// cloud chain wraps with %w (internal/speech/cloudchain.go:655), so a chain
	// failure around our callback error must still classify as an interrupt.
	wrapped := fmt.Errorf("speech: cloud TTS chain step: %w", errSpeakInterrupted)
	if !speakFailureIsInterrupt(wrapped) {
		t.Errorf("a %%w-wrapped interrupt sentinel stopped reading as an interrupt: %v", wrapped)
	}

	// Arm 3 (teeth): THE SAME ENGLISH TEXT, a different error. A substring
	// implementation cannot tell these apart; a code implementation must.
	sameText := errors.New("interrupted mid-speak")
	if speakFailureIsInterrupt(sameText) {
		t.Errorf("an unrelated error with the text %q classified as a user interrupt - the branch is reading the English text again (openchamber ChatInput.tsx:2205 shape)", sameText.Error())
	}

	// Arm 4: no error is not an interrupt.
	if speakFailureIsInterrupt(nil) {
		t.Errorf("nil classified as an interrupt")
	}

	// Arm 5: the neighbouring code (the speaker abort, D43#27) is a different
	// fact and must not be folded into this branch.
	if speakFailureIsInterrupt(errSpeakerAborted) {
		t.Errorf("the speaker-abort sentinel classified as a mid-synthesis interrupt; the two codes are distinct facts")
	}
}
