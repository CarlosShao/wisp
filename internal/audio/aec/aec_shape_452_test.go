package aec

// Ticket 452 AC#1(a): what NewApm really answers in each build shape. The
// shipped binary is the shape scripts/build.ps1 produces - no `-tags webrtcap`
// - so this file reads the stub's answer and pins the capability constant the
// product surface states. The tagged branch is the dedicated staging shape and
// pins the other direction, so neither answer can silently become the other.

import (
	"errors"
	"testing"
)

func TestAC452NewApmAnswersMatchTheBuildShape(t *testing.T) {
	apm, err := NewApm(16000)

	if BindingBuilt() {
		// The native binding's shape: an engine, never the stub error. This
		// branch needs the pinned webrtc-audio-processing tree on the CGO
		// flags (deps.toml); on a host without it the tagged build does not
		// compile at all, which is a louder answer than a red here.
		if err != nil {
			t.Fatalf("452 AC#1(a): BindingBuilt=true but NewApm(16000) answered %v, so the constant and the binding disagree", err)
		}
		if apm == nil {
			t.Fatal("452 AC#1(a): BindingBuilt=true and NewApm answered no error but no engine either")
		}
		if cerr := apm.Close(); cerr != nil {
			t.Fatalf("452 AC#1(a): Close() on the engine NewApm built: %v", cerr)
		}
		return
	}

	// The shipped shape: the stub, always ErrUnavailable, never an engine.
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("452 AC#1(a): in this build (BindingBuilt=false) NewApm(16000) answered %v, want ErrUnavailable: the stub is not the answer the shipped shape promises", err)
	}
	if apm != nil {
		t.Fatalf("452 AC#1(a): the stub returned an engine (%+v) alongside ErrUnavailable: a caller that checked only the object would think echo cancellation exists", apm)
	}
}
