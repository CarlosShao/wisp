package speech

import "testing"

// mb26 renders a byte count the way the evidence quotes it.
func mb26(v int64) float64 { return float64(v) / (1024 * 1024) }

// Test26TheRealSwitchDoesNotAddTheTwoFamilies is AC5's RSS half on the real
// models: the switch from one family to the other must not grow the process by a
// whole second engine, because that growth is the sum the resource budget forbids
// (SPEC-04 section 4: peak = max(ASR, TTS)).
//
// The discriminating quantity is the growth *across* the switch, compared with the
// cost of loading the incoming family from empty. With the release in place the
// switch nets out (the outgoing engine's memory goes away while the incoming one
// arrives); without it the growth is a full engine.
func Test26TheRealSwitchDoesNotAddTheTwoFamilies(t *testing.T) {
	base, ok := workingSetBytes()
	if !ok {
		t.Skipf("process working set is unavailable on this host; the residency reading is a Windows measurement")
	}

	slot := NewSlot()
	t.Cleanup(func() { _ = slot.Close() })

	// 1. TTS resident alone.
	if _, err := slot.LoadTts(func() (TtsEngine, error) { return NewMatchaTts(matchaConfigForTest(t)) }); err != nil {
		t.Skipf("matcha tts unavailable on this host: %v", err)
	}
	wsTts, ok := workingSetBytes()
	if !ok {
		t.Skipf("process working set became unavailable")
	}
	dTts := wsTts - base

	// 2. Switch the slot to ASR. If the release holds, this does not add a second
	//    engine's worth of memory on top of the first.
	if _, err := slot.LoadAsr(func() (AsrEngine, error) { return NewStreamingAsr(asrConfigForTest(t)) }); err != nil {
		t.Skipf("streaming ASR unavailable on this host: %v", err)
	}
	wsAsr, ok := workingSetBytes()
	if !ok {
		t.Skipf("process working set became unavailable")
	}
	switchGrowth := wsAsr - wsTts

	// 3. The cost of loading ASR from empty, on the same process, as the yardstick
	//    for step 2.
	if err := slot.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	wsEmpty, ok := workingSetBytes()
	if !ok {
		t.Skipf("process working set became unavailable")
	}
	if _, err := slot.LoadAsr(func() (AsrEngine, error) { return NewStreamingAsr(asrConfigForTest(t)) }); err != nil {
		t.Skipf("streaming ASR unavailable on this host: %v", err)
	}
	wsAsrAlone, ok := workingSetBytes()
	if !ok {
		t.Skipf("process working set became unavailable")
	}
	dAsrAlone := wsAsrAlone - wsEmpty

	if err := slot.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	wsReleased, ok := workingSetBytes()
	released := wsReleased - base

	t.Logf("rss: baseline=%.1fMB ttsResident=%.1fMB (+%.1fMB) asrResidentAfterSwitch=%.1fMB (+%.1fMB over tts) "+
		"switchGrowth=%.1fMB asrFromEmpty=%.1fMB releasedResidual=%.1fMB",
		mb26(base), mb26(wsTts), mb26(dTts), mb26(wsAsr), mb26(switchGrowth),
		mb26(switchGrowth), mb26(dAsrAlone), mb26(released))

	if dAsrAlone <= 0 {
		t.Fatalf("loading ASR from empty grew the process by %.1fMB; the yardstick is unusable", mb26(dAsrAlone))
	}
	if switchGrowth >= dAsrAlone {
		t.Fatalf("switching families grew the process by %.1fMB, as much as loading the incoming engine from "+
			"empty (%.1fMB): the outgoing family was still resident", mb26(switchGrowth), mb26(dAsrAlone))
	}
	if switchGrowth > dAsrAlone*3/4 {
		t.Fatalf("switching families grew the process by %.1fMB against a %.1fMB incoming engine: too close to a sum",
			mb26(switchGrowth), mb26(dAsrAlone))
	}
	if ok && released > maxInt64(dTts, dAsrAlone)/2 {
		t.Fatalf("releasing the slot left %.1fMB over baseline against the larger family's %.1fMB: unload did not free",
			mb26(released), mb26(maxInt64(dTts, dAsrAlone)))
	}
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
