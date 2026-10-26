package speech

import (
	"errors"
	"strings"
	"testing"
)

// recorder26 records engine lifecycle events so the slot's *ordering* can be
// asserted: a lock alone would keep two engines alive without ever using them
// concurrently, and the RSS peak is about being alive, not about being used.
type recorder26 struct{ events []string }

func (r *recorder26) record(ev string) { r.events = append(r.events, ev) }
func (r *recorder26) String() string   { return strings.Join(r.events, ",") }

type fakeAsr26 struct {
	name string
	rec  *recorder26
}

func (f *fakeAsr26) Begin()                 {}
func (f *fakeAsr26) Feed(samples []float32) {}
func (f *fakeAsr26) Partial() string        { return "" }
func (f *fakeAsr26) Finish() string         { return "" }
func (f *fakeAsr26) InFlight() bool         { return false }
func (f *fakeAsr26) Fed() int               { return 0 }
func (f *fakeAsr26) Close() error {
	f.rec.record("close:" + f.name)
	return nil
}

type fakeTts26 struct {
	name string
	rec  *recorder26
}

func (f *fakeTts26) Speak(text string, emit func([]float32) error) error { return nil }
func (f *fakeTts26) SampleRate() int                                     { return 22050 }
func (f *fakeTts26) Close() error {
	f.rec.record("close:" + f.name)
	return nil
}

// Test26TheSlotReleasesBeforeItLoads is the core of AC5: the resident family is
// closed *before* the incoming loader runs, so the two never overlap.
func Test26TheSlotReleasesBeforeItLoads(t *testing.T) {
	rec := &recorder26{}
	slot := NewSlot()

	asr, err := slot.LoadAsr(func() (AsrEngine, error) {
		rec.record("load:asr")
		return &fakeAsr26{name: "asr", rec: rec}, nil
	})
	if err != nil {
		t.Fatalf("LoadAsr: %v", err)
	}
	if slot.Resident() != FamilyAsr || slot.Asr() != asr {
		t.Fatalf("resident = %q, asr match = %v, want the ASR family", slot.Resident(), slot.Asr() == asr)
	}

	tts, err := slot.LoadTts(func() (TtsEngine, error) {
		rec.record("load:tts")
		return &fakeTts26{name: "tts", rec: rec}, nil
	})
	if err != nil {
		t.Fatalf("LoadTts: %v", err)
	}

	const want = "load:asr,close:asr,load:tts"
	if got := rec.String(); got != want {
		t.Fatalf("slot lifecycle = %q, want %q (the ASR engine must be released before TTS loads)", got, want)
	}
	if slot.Resident() != FamilyTts || slot.Tts() != tts {
		t.Fatalf("resident = %q, want the TTS family", slot.Resident())
	}
	if slot.Asr() != nil {
		t.Fatalf("the ASR family is still reachable while TTS is resident")
	}
	if slot.Loads() != 2 || slot.Closes() != 1 {
		t.Fatalf("counters loads=%d closes=%d, want 2 and 1", slot.Loads(), slot.Closes())
	}

	// Asking for the family that already holds the slot must not rebuild it: a
	// second load would pay 170MB and 5.8s for nothing.
	again, err := slot.LoadTts(func() (TtsEngine, error) {
		rec.record("load:tts-again")
		return &fakeTts26{name: "tts-again", rec: rec}, nil
	})
	if err != nil {
		t.Fatalf("second LoadTts: %v", err)
	}
	if again != tts {
		t.Fatalf("the second LoadTts built a new engine instead of returning the resident one")
	}
	if got := rec.String(); got != want {
		t.Fatalf("slot lifecycle after a redundant load = %q, want %q unchanged", got, want)
	}

	if err := slot.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if slot.Resident() != FamilyNone || slot.Tts() != nil {
		t.Fatalf("slot still holds %q after Close", slot.Resident())
	}
	if got, want := rec.String(), "load:asr,close:asr,load:tts,close:tts"; got != want {
		t.Fatalf("slot lifecycle = %q, want %q", got, want)
	}
	if err := slot.Close(); err != nil {
		t.Fatalf("second Close must be a no-op: %v", err)
	}
	if got, want := rec.String(), "load:asr,close:asr,load:tts,close:tts"; got != want {
		t.Fatalf("a second Close released something twice: %q", got)
	}
}

// Test26AFailedLoadLeavesTheSlotEmpty is the fail-closed half: the release has
// already happened, so a loader error must not leave a half-built residency
// behind.
func Test26AFailedLoadLeavesTheSlotEmpty(t *testing.T) {
	rec := &recorder26{}
	slot := NewSlot()
	if _, err := slot.LoadAsr(func() (AsrEngine, error) {
		rec.record("load:asr")
		return &fakeAsr26{name: "asr", rec: rec}, nil
	}); err != nil {
		t.Fatalf("LoadAsr: %v", err)
	}

	boom := errors.New("no model")
	if _, err := slot.LoadTts(func() (TtsEngine, error) { return nil, boom }); !errors.Is(err, boom) {
		t.Fatalf("LoadTts = %v, want the loader's error", err)
	}
	if slot.Resident() != FamilyNone || slot.Asr() != nil || slot.Tts() != nil {
		t.Fatalf("failed load left residency %q behind", slot.Resident())
	}
	const want = "load:asr,close:asr"
	if got := rec.String(); got != want {
		t.Fatalf("slot lifecycle = %q, want %q", got, want)
	}
}

// Test26TheSlotNeedsLoadersAndSurvivesNilEngines covers the two degenerate calls
// a caller will eventually make.
func Test26TheSlotNeedsLoadersAndSurvivesNilEngines(t *testing.T) {
	slot := NewSlot()
	if _, err := slot.LoadAsr(nil); err == nil {
		t.Fatalf("LoadAsr(nil) returned nil error")
	}
	if _, err := slot.LoadTts(nil); err == nil {
		t.Fatalf("LoadTts(nil) returned nil error")
	}
	if err := slot.Close(); err != nil {
		t.Fatalf("Close on an empty slot: %v", err)
	}
	if slot.Loads() != 0 || slot.Closes() != 0 {
		t.Fatalf("empty slot counters loads=%d closes=%d, want 0 and 0", slot.Loads(), slot.Closes())
	}
}
