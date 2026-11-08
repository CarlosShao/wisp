//go:build windows

package main

// Ticket 291 AC#2 four-shape live-mic probe. CLONE-TREE ONLY: this file lives
// in the throwaway export tree /d/tmp/wisp/20261108/291-mic-tree and is never
// committed - the main repo stays a zero-production-change read-only leg for
// ticket 291 (write surface = .scratch/wisp/probes/291/** + evidence + one
// Progress log line). It reuses the 247 live harness unchanged
// (bootAudioRuntime / writeAudioConfig / assembleCapture /
// requireNoCompetingAudioProcesses / livePhase / phaseRouter) so the readings
// come from the same production assembly family the ticket names.
//
// The four shapes of ticket 291 AC#2:
//   shape 1: room silence, 3s, real microphone.
//   shape 2: a person speaking at normal volume, 25s, real microphone.
//   shape 3: a known-amplitude tone (440 Hz sine, sample amplitude 12000)
//            played out of the render endpoint; its SOURCE level is
//            12000/sqrt(2)/32768 = 0.258926 - what the microphone hears is
//            attenuated by speaker volume, distance and the room, so the mic
//            reading is reported next to the source figure, not claimed equal.
//   shape 4: a synthetic DC(12345) + whole-cycle sine(amp 6000) frame pushed
//            through the production ruler audio.LevelOfSamples and checked
//            against the closed-form RMS sqrt(dc^2 + amp^2/2)/32768. This is
//            the in-repo fixture family (capturelevel_windows_test.go's DC
//            12345); per the ticket it judges nothing about production.

import (
	"encoding/binary"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/config"
)

func TestAC291FourShapeLiveMicReadings(t *testing.T) {
	if os.Getenv("WISP_LIVE_MIC") != "1" {
		t.Skip("291 AC#2 needs a real microphone: run with WISP_LIVE_MIC=1 after both tasklist rulers read 0")
	}
	requireNoCompetingAudioProcesses(t)

	// ---- shape 4 (no device): DC + whole-cycle sine through the production ruler ----
	const dc = 12345.0
	const amp = 6000.0
	const period = 32 // samples; FrameSamples 512 / 32 = 16 whole cycles
	synth := make([]int16, audio.FrameSamples)
	for i := range synth {
		synth[i] = int16(math.Round(dc + amp*math.Sin(2*math.Pi*float64(i%period)/period)))
	}
	gotSynth := audio.LevelOfSamples(synth)
	wantSynth := math.Sqrt(dc*dc+amp*amp/2) / audio.LevelFullScale
	t.Logf("291 shape4 DC+sin frame: LevelOfSamples=%.6f closed-form=%.6f delta=%.3e",
		gotSynth, wantSynth, math.Abs(gotSynth-wantSynth))
	dcOnly := make([]int16, audio.FrameSamples)
	for i := range dcOnly {
		dcOnly[i] = 12345
	}
	t.Logf("291 shape4 ref pure-DC-12345 frame: LevelOfSamples=%.6f (=12345/32768)", audio.LevelOfSamples(dcOnly))

	// ---- live shapes 1/2/3 on the real microphone ----
	rt, shut := bootAudioRuntime(t)
	dir := writeAudioConfig(t, func(c *config.Config) { c.Audio.MicMutedDefault = false })

	quiet1 := &livePhase{name: "shape1-quiet"}
	tone := &livePhase{name: "shape3-tone"}
	speech := &livePhase{name: "shape2-speech"}
	quiet2 := &livePhase{name: "control-quiet-after"}
	router := &phaseRouter{}
	router.set(quiet1)

	ra := assembleCapture(rt, dir, router.route, newRealCaptureSource)
	if !ra.started.Load() {
		t.Fatalf("the live capture leg did not come up on this machine: %s", ra.posture())
	}
	mic := ra.mic.(*audio.WASAPIMicrophone)
	capDev, renderDev, err := mic.Endpoints()
	if err != nil {
		t.Fatalf("endpoint pair: %v", err)
	}
	t.Logf("291 source: capture=%q render=%q", capDev.String(), renderDev.String())

	time.Sleep(3 * time.Second)
	quiet1.span = 3 * time.Second

	router.set(tone)
	stopTone := playToneWav(t)
	time.Sleep(6 * time.Second)
	tone.span = 6 * time.Second
	stopTone()

	router.set(speech)
	time.Sleep(25 * time.Second)
	speech.span = 25 * time.Second

	router.set(quiet2)
	time.Sleep(3 * time.Second)
	quiet2.span = 3 * time.Second

	st := mic.Stats()
	t.Logf("291 counters at the seam: levels=%d frames_sent=%d frames_dropped=%d last_err=%q",
		ra.levels.Load(), st.FramesSent, st.FramesDropped, st.LastError)
	for _, p := range []*livePhase{quiet1, tone, speech, quiet2} {
		t.Logf("291 reading: %s", p.String())
	}
	for _, p := range []*livePhase{quiet1, tone, speech, quiet2} {
		if p.count() == 0 {
			t.Fatalf("a phase recorded nothing: %s", p)
		}
	}
	_ = shut()
}

// playToneWav writes the known 440 Hz sine wav (sample amplitude 12000, source
// level 12000/sqrt(2)/32768 = 0.258926) to a temp file and loops it out of the
// render endpoint; returns the stop function. Same shape as the 247 leg's
// playAlarmThroughSpeakers, but with a known amplitude instead of Alarm01.wav.
func playToneWav(t *testing.T) func() {
	t.Helper()
	const rate = 16000
	const amp = 12000.0
	const secs = 2
	n := rate * secs
	pcm := make([]byte, 44+n*2)
	copy(pcm[0:], "RIFF")
	binary.LittleEndian.PutUint32(pcm[4:], uint32(36+n*2))
	copy(pcm[8:], "WAVEfmt ")
	binary.LittleEndian.PutUint32(pcm[16:], 16)
	binary.LittleEndian.PutUint16(pcm[20:], 1)
	binary.LittleEndian.PutUint16(pcm[22:], 1)
	binary.LittleEndian.PutUint32(pcm[24:], rate)
	binary.LittleEndian.PutUint32(pcm[28:], rate*2)
	binary.LittleEndian.PutUint16(pcm[32:], 2)
	binary.LittleEndian.PutUint16(pcm[34:], 16)
	copy(pcm[36:], "data")
	binary.LittleEndian.PutUint32(pcm[40:], uint32(n*2))
	for i := 0; i < n; i++ {
		v := int16(amp * math.Sin(2*math.Pi*440*float64(i)/rate))
		binary.LittleEndian.PutUint16(pcm[44+i*2:], uint16(v))
	}
	path := filepath.Join(t.TempDir(), "tone440-amp12000.wav")
	if err := os.WriteFile(path, pcm, 0o644); err != nil {
		t.Fatalf("write tone wav: %v", err)
	}
	t.Logf("291 shape3 tone source: 440Hz amp=12000 source_level=%.6f", amp/math.Sqrt2/audio.LevelFullScale)
	script := "$ErrorActionPreference='Stop';" +
		"$p=New-Object System.Media.SoundPlayer '" + path + "';" +
		"while ($true) { $p.PlaySync(); Start-Sleep -Milliseconds 100 }"
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", script)
	if err := cmd.Start(); err != nil {
		t.Fatalf("start the tone source: %v", err)
	}
	var once sync.Once
	stop := func() {
		once.Do(func() {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		})
	}
	t.Cleanup(stop)
	return stop
}
