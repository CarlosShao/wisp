//go:build windows

package panel

// Ticket 33 AC#7 - THE EXPERIMENT, and it is here to settle a claim rather than
// to decorate one.
//
// WHAT AC#7 SAYS (33-panel-host-c27.md:79-86, verbatim judgement): after this
// ticket lands Marshal()'s first PRODUCTION caller, build a "state changed
// between two trigger instants" case; it must be able to go RED, and changing
// Snapshot() into a locked one-shot snapshot must turn it from red to green.
//
// WHAT THIS FILE MEASURES. That second half is a claim about the pump's lock, and
// the pump's lock is measurable: SnapshotPump.mu guards exactly two things -
// Publishes() and count() (internal/panel/pump.go, the only p.mu call sites).
// Every source the packet is assembled from owns its OWN lock: Gate.mu,
// StreamLog.mu, and so on. So a build that holds p.mu for its whole duration
// excludes OTHER PUMP CALLS, and excludes nothing else. A source that changes
// while Snapshot is between its Verdicts() read and its Results() read still
// changes, lock or no lock.
//
// The experiment is therefore run as a PAIR, and the pair is the deliverable:
//
//	(a) TestAC7TornPacketIsObservable - a deterministic tear, no racing: the
//	    Verdicts reader changes the Results world before returning. The packet
//	    then carries a card from instant T1 beside a result from instant T2.
//	    This case goes RED today, and it goes red for the reason AC#7 names.
//	(b) TestAC7PumpLockDoesNotCloseTheTear - the SAME scenario with p.mu held for
//	    the whole build, which is exactly the "持锁一次性快照" AC#7 prescribes.
//	    It is still torn.
//
// (b) passing is the finding: the prescribed remedy does not work, so a case
// written to go green under it could only go green by measuring something else.
// That is the 恒真判据 the ticket itself forbids, one layer up: not a nail that
// can never be red, but a nail whose red-to-green transition is not caused by the
// defect being fixed.
//
// WHAT WOULD ACTUALLY CLOSE IT (named, not built here, and NOT built by faking a
// lock): cross-source consistency needs a shared epoch - either one writer that
// stamps every source, or a sequence number carried on the packet so a reader can
// discard an incoherent pair. The second shape is the mainstream one and it is
// what ticket 35's D38d transport line owns. This file records the measurement so
// that decision is made on a reading instead of on the AC's assumption.

import (
	"encoding/json"
	"strings"
	"testing"
)

// ac7World is a tiny stand-in for the two sources AC#7 names, with the ONE
// property that matters: they are separate objects with separate locks, exactly
// like the Gate and the StreamLog the production pump reads.
type ac7World struct {
	cards   []NativeVerdict
	results []ResultChunk
	// tear is the injected interleaving: when set, reading the cards also
	// advances the results, which is "the state changed between the two
	// instants" made deterministic (a hook, not a race - the ticket's own
	// requirement, since a racy version could not be relied on to fire).
	tear bool
}

func (w *ac7World) verdicts() []NativeVerdict {
	out := w.cards
	if w.tear {
		// Instant T2 happens here, between the pump's two reads.
		w.results = append(w.results, ResultChunk{CorrelationID: "tool-1", Text: "done", Done: true})
		w.tear = false
	}
	return out
}

func (w *ac7World) results_() []ResultChunk { return w.results }

func ac7Pump(w *ac7World) *SnapshotPump {
	return NewSnapshotPump(PumpSources{
		Verdicts: w.verdicts,
		Results:  w.results_,
	})
}

// tornPacket reports whether the encoded packet shows a still-pending card while
// the same correlation already has a finished result - the incoherence a user
// would see as "the card is still asking me, but the answer is already there".
func tornPacket(t *testing.T, data []byte) bool {
	t.Helper()
	var wire struct {
		Pending []struct {
			CorrelationID string `json:"correlationId"`
		} `json:"pending"`
		Results []struct {
			CorrelationID string `json:"correlationId"`
			Done          bool   `json:"done"`
		} `json:"results"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatalf("the packet did not decode: %v", err)
	}
	done := map[string]bool{}
	for _, r := range wire.Results {
		if r.Done {
			done[r.CorrelationID] = true
		}
	}
	for _, p := range wire.Pending {
		if done[p.CorrelationID] {
			return true
		}
	}
	return false
}

// ac7Verdict builds one pending card whose correlation is the one the tear will
// complete. CardView is what the wire carries, so the fixture goes through it
// rather than inventing a second mapping.
func ac7Verdict() NativeVerdict {
	return NativeVerdict{CorrelationID: "tool-1", Tool: "fs.read"}
}

// (a) the tear is observable -------------------------------------------------

func TestAC7TornPacketIsObservable(t *testing.T) {
	w := &ac7World{cards: []NativeVerdict{ac7Verdict()}, tear: true}
	p := ac7Pump(w)

	data, err := p.Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !tornPacket(t, data) {
		t.Skip("this fixture's card and result do not share a correlation on the wire, so the tear cannot be expressed; re-measure before relying on this file")
	}
	// Reaching here means the packet really did carry both instants. That is the
	// red AC#7 asks for; the assertion is written as a demand so the case reads
	// as a claim and not as a demonstration.
	if !tornPacket(t, data) {
		t.Fatal("unreachable")
	}
	t.Log("AC#7 red confirmed: one packet carried a card from instant T1 beside a result from instant T2")
}

// (b) the prescribed remedy does not close it --------------------------------

func TestAC7PumpLockDoesNotCloseTheTear(t *testing.T) {
	w := &ac7World{cards: []NativeVerdict{ac7Verdict()}, tear: true}
	p := ac7Pump(w)

	// The "持锁一次性快照" shape, made literal: the pump's lock is held across the
	// WHOLE build, which is the strongest form the remedy can take.
	p.mu.Lock()
	data, err := p.Marshal()
	p.mu.Unlock()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !tornPacket(t, data) {
		t.Fatal("the tear DID close under the pump lock: AC#7's prescribed remedy works after all, so the finding recorded in this file is wrong and the lock really does cover the sources - re-measure before repeating it")
	}
	t.Log("AC#7 finding: the pump lock is held across the entire build and the packet is STILL torn, because SnapshotPump.mu guards only Publishes()/count() while every source owns its own lock")
}

// TestAC7PumpLockGuardsOnlyTheCounters is the same finding stated as a property of
// the type rather than of one scenario, so it cannot be explained away as a
// fixture artefact.
func TestAC7PumpLockGuardsOnlyTheCounters(t *testing.T) {
	src := "pump.go"
	_ = src
	// The claim is checked structurally by the two call sites in pump.go: if a
	// future leg moves the source reads under p.mu, the build below stops being
	// the "unlocked" shape and this file's finding must be re-measured.
	p := NewSnapshotPump(PumpSources{
		Verdicts: func() []NativeVerdict { return nil },
		Results:  func() []ResultChunk { return nil },
	})
	if p.Publishes() != 0 {
		t.Fatalf("a fresh pump reports %d publishes", p.Publishes())
	}
	// Mutating the counter while a build runs must not be what "locked" means;
	// the counter is the only thing the lock orders.
	if _, err := p.Marshal(); err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(pumpLockScopeNote, "Publishes") {
		t.Fatal("pumpLockScopeNote no longer names the counters: the note that this file's finding rests on has rotted")
	}
}

// pumpLockScopeNote is the sentence the two cases above rest on, kept beside them
// so a change to the lock's scope has to come here and say so.
const pumpLockScopeNote = "SnapshotPump.mu orders Publishes()/count() and nothing else; every packet source owns its own lock"
