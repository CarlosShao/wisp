package models

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

func newWalkRig(t *testing.T, srv *httptest.Server) (*DownloadingBridge, *[]statemachine.Effect, *statemachine.Machine) {
	t.Helper()
	m := manifestWithURLs(t, srv.URL)
	mgr, _ := newTestManager(t, m, nil, nil)
	var effects []statemachine.Effect
	machine := statemachine.New(statemachine.Options{
		Initial: statemachine.StateFirstRun,
		Sink:    func(e statemachine.Effect) { effects = append(effects, e) },
	})
	return WireDownloading(mgr, machine), &effects, machine
}

func TestDownloadingWalkSuccess(t *testing.T) {
	srv := newCountingServer(t, map[string][]byte{
		"/vad-fixture/vad.onnx": vadBody,
	})
	bridge, effects, machine := newWalkRig(t, srv.srv)

	final, err := bridge.Run(context.Background(), "vad-fixture")
	if err != nil {
		t.Fatalf("walk failed: %v", err)
	}
	if final != statemachine.StateFirstRun {
		t.Fatalf("walk ended in %s, want FirstRun", final)
	}
	// State log: enter Downloading (row #2), exit to FirstRun (row #37).
	log := bridge.StateLog()
	if len(log) != 2 || log[0] != statemachine.StateDownloading || log[1] != statemachine.StateFirstRun {
		t.Fatalf("state walk %v, want [Downloading FirstRun]", log)
	}
	if bridge.Ticks() == 0 {
		t.Fatal("no progress ticks observed during Downloading")
	}
	if bridge.LastPercent() != 100 {
		t.Fatalf("last percent %v, want 100", bridge.LastPercent())
	}
	// Row #37 side effect fired.
	found := false
	for _, e := range *effects {
		if e.Name == "model.verify-sha256-signature" {
			found = true
		}
	}
	if !found {
		t.Fatalf("row #37 side effect not fired; effects %v", *effects)
	}
	_ = machine
}

func TestDownloadingWalkFailure(t *testing.T) {
	srv := newCountingServer(t, map[string][]byte{}) // everything 404
	bridge, _, machine := newWalkRig(t, srv.srv)

	_, err := bridge.Run(context.Background(), "vad-fixture")
	if err == nil {
		t.Fatal("walk against a dead mirror must fail")
	}
	if machine.State() != statemachine.StateError {
		t.Fatalf("machine in %s, want Error (row #37 failed exit)", machine.State())
	}
	log := bridge.StateLog()
	if len(log) != 2 || log[0] != statemachine.StateDownloading || log[1] != statemachine.StateError {
		t.Fatalf("state walk %v, want [Downloading Error]", log)
	}
}

func TestDownloadingWalkRejectsIllegalEnter(t *testing.T) {
	// D43 row #2 only allows FirstRun + model-missing -> Downloading. From
	// Sleeping the enter is illegal and must be REJECTED, not force-applied
	// (undefined = stop, D22 gate 3).
	srv := newCountingServer(t, map[string][]byte{
		"/vad-fixture/vad.onnx": vadBody,
	})
	m := manifestWithURLs(t, srv.URL())
	mgr, _ := newTestManager(t, m, nil, nil)
	machine := statemachine.New(statemachine.Options{Initial: statemachine.StateSleeping})
	bridge := WireDownloading(mgr, machine)

	if _, err := bridge.Run(context.Background(), "vad-fixture"); err == nil {
		t.Fatal("illegal enter (Sleeping + model-missing) was accepted")
	}
	if machine.State() != statemachine.StateSleeping {
		t.Fatalf("state changed on rejected dispatch: %s", machine.State())
	}
}
