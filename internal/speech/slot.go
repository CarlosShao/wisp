package speech

import (
	"errors"
	"sync"
)

// This file is the engine-slot residency policy ticket 26 asks for (SPEC-04
// section 4, D32 16.3.3): at most one engine family resident at a time, so the
// peak is max(ASR, TTS) and not their sum.
//
// The reason it exists is a measurement, not a preference. Ticket 02's spike
// measured the serial switch at P50 4726ms against the 1600ms budget and the
// fallback in D32 16.3.3 became the active policy ("TTS resident, ASR on
// demand"). 26-r1 re-measured the two halves on this host: ASR int8 resident
// costs 281.4MB, the matcha TTS costs 170.2MB, and holding both costs 436.6MB --
// the sum, to the megabyte. 700MB of budget does not survive a policy that adds.
//
// Scope, stated so nobody reads more into it: this is the Path-T rule only.
// D47's Path C needs ASR and TTS resident together (ticket 59's AEC gate), so the
// slot must NOT be wired into the conversation path -- a slot on that path would
// make barge-in impossible by construction. The slot is a policy object with no
// opinion about who calls it; today that is this package's tests, and the
// composition root (cmd/wisp) plus the cloud chains of ticket 61 are where a
// production caller would land.

// Family names which engine family occupies the slot.
type Family string

const (
	// FamilyNone is an empty slot.
	FamilyNone Family = ""
	// FamilyAsr is the ASR family (one AsrEngine).
	FamilyAsr Family = "asr"
	// FamilyTts is the TTS family (one TtsEngine).
	FamilyTts Family = "tts"
	// FamilyKws is the wake-word family (one WakeWordEngine). Ticket 41.
	FamilyKws Family = "kws"
)

// Slot is the serial engine slot: one family resident, the other released first.
//
// The invariant it buys is an ordering, not a lock: the resident engine is closed
// before the incoming loader runs, so there is no instant at which both exist and
// the RSS peak cannot be a sum. A lock alone would not buy that -- two engines
// can be resident without ever being used concurrently.
type Slot struct {
	mu     sync.Mutex
	family Family
	asr    AsrEngine
	tts    TtsEngine
	kws    WakeWordEngine
	loads  int
	closes int
}

// NewSlot opens an empty slot.
func NewSlot() *Slot { return &Slot{} }

// Resident names the family currently holding the slot.
func (s *Slot) Resident() Family {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.family
}

// Asr is the resident ASR engine, or nil when the TTS family holds the slot.
func (s *Slot) Asr() AsrEngine {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.asr
}

// Tts is the resident TTS engine, or nil when the ASR family holds the slot.
func (s *Slot) Tts() TtsEngine {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tts
}

// Kws is the resident wake-word engine, or nil when another family holds the
// slot (ticket 41, same shape as Asr/Tts).
func (s *Slot) Kws() WakeWordEngine {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.kws
}

// Loads and Closes are the counters the readings are made of: a session that
// switches families must show one load per residency and one close per release,
// which is what makes "the old family really went away" observable instead of
// asserted.
func (s *Slot) Loads() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loads
}

// Closes is how many engines this slot has released.
func (s *Slot) Closes() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closes
}

// LoadAsr makes the ASR family resident, returning the engine already resident
// when it is the same family (a second load would trade a live session for a
// second copy of the same 281MB).
//
// The release happens before the load, and a failed load leaves the slot empty
// rather than holding a half-built engine: fail closed.
func (s *Slot) LoadAsr(load func() (AsrEngine, error)) (AsrEngine, error) {
	if load == nil {
		return nil, errors.New("speech: slot LoadAsr needs a loader")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.family == FamilyAsr && s.asr != nil {
		return s.asr, nil
	}
	if err := s.releaseLocked(); err != nil {
		return nil, err
	}
	eng, err := load()
	if err != nil {
		return nil, err
	}
	s.asr, s.family = eng, FamilyAsr
	s.loads++
	return eng, nil
}

// LoadTts makes the TTS family resident, with the same rules as LoadAsr.
func (s *Slot) LoadTts(load func() (TtsEngine, error)) (TtsEngine, error) {
	if load == nil {
		return nil, errors.New("speech: slot LoadTts needs a loader")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.family == FamilyTts && s.tts != nil {
		return s.tts, nil
	}
	if err := s.releaseLocked(); err != nil {
		return nil, err
	}
	eng, err := load()
	if err != nil {
		return nil, err
	}
	s.tts, s.family = eng, FamilyTts
	s.loads++
	return eng, nil
}

// LoadKws makes the wake-word family resident, with the same rules as
// LoadAsr/LoadTts: release first, then load, and a failed load leaves the slot
// empty rather than holding a half-built engine.
func (s *Slot) LoadKws(load func() (WakeWordEngine, error)) (WakeWordEngine, error) {
	if load == nil {
		return nil, errors.New("speech: slot LoadKws needs a loader")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.family == FamilyKws && s.kws != nil {
		return s.kws, nil
	}
	if err := s.releaseLocked(); err != nil {
		return nil, err
	}
	eng, err := load()
	if err != nil {
		return nil, err
	}
	s.kws, s.family = eng, FamilyKws
	s.loads++
	return eng, nil
}

// Close releases whatever is resident (C11). Idempotent.
func (s *Slot) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.releaseLocked()
}

// releaseLocked closes the resident engine and empties the slot. The caller holds
// the mutex. Both engine interfaces have idempotent Close, so this is safe on an
// empty slot; the counters only move when there was something to release.
func (s *Slot) releaseLocked() error {
	var err error
	switch s.family {
	case FamilyAsr:
		if s.asr != nil {
			err = s.asr.Close()
			s.closes++
		}
		s.asr = nil
	case FamilyTts:
		if s.tts != nil {
			err = s.tts.Close()
			s.closes++
		}
		s.tts = nil
	case FamilyKws:
		if s.kws != nil {
			err = s.kws.Close()
			s.closes++
		}
		s.kws = nil
	}
	s.family = FamilyNone
	return err
}
