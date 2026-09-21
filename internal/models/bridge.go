package models

import (
	"context"
	"fmt"
	"sync"

	"github.com/CarlosShao/wisp/internal/statemachine"
)

// DownloadingBridge wires download progress to the D43 ball states:
//
//	row #2  FirstRun  + model-missing        -> Downloading   (enter)
//	        Downloading + progress ticks      -> Downloading   (ring %, no state change)
//	row #37 Downloading + download-completed -> FirstRun      (exit, ok)
//	        Downloading + download-failed     -> Error         (exit, failed/cancelled)
//
// Progress itself is ball-local UI state (ring percent), not a D43 event, so
// ticks never dispatch; the walk is observable through StateLog. Cancellation
// exits via download-failed too: row #37 offers exactly two exits and
// inventing a third is a D22 violation.
type DownloadingBridge struct {
	mgr     *Manager
	machine *statemachine.Machine

	mu      sync.Mutex
	ticks   int
	percent float64
	log     []statemachine.State
}

// WireDownloading attaches a Manager to a Machine.
func WireDownloading(mgr *Manager, machine *statemachine.Machine) *DownloadingBridge {
	return &DownloadingBridge{mgr: mgr, machine: machine}
}

// Run performs one full downloading walk for model id. Returns the machine
// state after the walk. StateLog captures [post-enter, post-exit, ...].
func (b *DownloadingBridge) Run(ctx context.Context, id string) (statemachine.State, error) {
	if _, err := b.machine.Dispatch(statemachine.EvModelMissing, nil); err != nil {
		return b.machine.State(), fmt.Errorf("models: enter Downloading: %w", err)
	}
	_ = b.record() // entered Downloading (row #2)
	_, err := b.mgr.Ensure(ctx, id, b.onProgress)
	if err != nil {
		if _, derr := b.machine.Dispatch(statemachine.EvDownloadFailed, nil); derr != nil {
			return b.record(), fmt.Errorf("models: exit Downloading after failure: %w (download err: %v)", derr, err)
		}
		return b.record(), err
	}
	// Ticket 109's hand-off check. Ensure verified the bytes *before* it
	// returned and the reader opens them *after* this call announces the model
	// available; the model directory stays inherit-wide on purpose (ticket 95),
	// so that span is writable by anything the parent grants write to. Re-verify
	// here, at the one point in this package where "available" is decided, and
	// fail the walk instead of announcing a swapped file. The state machine gets
	// no new event: row #37 already has exactly the two exits this needs.
	if verr := b.mgr.VerifyInstalled(id); verr != nil {
		if _, derr := b.machine.Dispatch(statemachine.EvDownloadFailed, nil); derr != nil {
			return b.record(), fmt.Errorf("models: exit Downloading after the hand-off check: %w (re-verify err: %v)", derr, verr)
		}
		return b.record(), verr
	}
	if _, derr := b.machine.Dispatch(statemachine.EvDownloadCompleted, nil); derr != nil {
		return b.record(), fmt.Errorf("models: exit Downloading after completion: %w", derr)
	}
	return b.record(), nil
}

func (b *DownloadingBridge) onProgress(ev ProgressEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.ticks++
	if ev.Phase == PhaseDownloading || ev.Phase == PhaseDone {
		b.percent = ev.Percent
	}
}

// record snapshots the current state into the walk log and returns it.
func (b *DownloadingBridge) record() statemachine.State {
	s := b.machine.State()
	b.mu.Lock()
	defer b.mu.Unlock()
	b.log = append(b.log, s)
	return s
}

// Ticks reports how many progress events were observed (ring activity).
func (b *DownloadingBridge) Ticks() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.ticks
}

// LastPercent reports the last download-phase percent seen (0-100).
func (b *DownloadingBridge) LastPercent() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.percent
}

// StateLog returns the state walk in order (enter, exits across runs).
func (b *DownloadingBridge) StateLog() []statemachine.State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]statemachine.State(nil), b.log...)
}
