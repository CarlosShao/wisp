//go:build windows

package main

// Ticket 247 AC#2's real-machine reading, plus P8 form 甲's cost clause.
//
// This is the file that answers "does the number on the ball come from sound or
// from a command line". It runs the production assembly (assembleCapture with
// the real WASAPI source) against this machine's actual default capture endpoint
// and records the scalars the leg would hand to Ball.SetAudioLevel.
//
// What it is NOT, stated because the ticket names each of these as a cheating
// shape: not cmd/balldebug's synthetic envelope, not internal/audio's wav
// injector (that is the C8 test seam, not the thing under test), and not a skip
// dressed up as a pass. Without WISP_LIVE_MIC=1 it skips, AC#2's reading does
// not exist, and the box stays closed to a person.
//
// How "loud" is produced here: a real Windows alarm sample played out of the
// render endpoint while the real microphone listens to the room. The bytes never
// enter the pipeline - the only thing that reaches the seam is what the capsule
// heard. If the speakers are off, routed to an endpoint this microphone cannot
// hear, or nobody is speaking into it, the loud phase reads like the quiet ones
// and this test goes RED with all three phases printed. That is the honest
// answer; it is not a pass and it is not faked with a synthetic source.

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// livePhase is one window of recorded scalars plus how long it stayed open.
type livePhase struct {
	name   string
	mu     sync.Mutex
	levels []float32
	span   time.Duration
}

func (p *livePhase) take(level float32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.levels = append(p.levels, level)
}

func (p *livePhase) count() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.levels)
}

func (p *livePhase) max() float32 {
	p.mu.Lock()
	defer p.mu.Unlock()
	var m float32
	for _, l := range p.levels {
		if l > m {
			m = l
		}
	}
	return m
}

func (p *livePhase) mean() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.levels) == 0 {
		return 0
	}
	var sum float64
	for _, l := range p.levels {
		sum += float64(l)
	}
	return sum / float64(len(p.levels))
}

// perSecond is the cadence P8 form 甲 demands be re-measured now that delivery
// happens on the pinned capture thread: one level per seam frame, about 31.25/s
// at FrameDuration = 32ms.
func (p *livePhase) perSecond() float64 {
	if p.span <= 0 {
		return 0
	}
	return float64(p.count()) / p.span.Seconds()
}

func (p *livePhase) String() string {
	return p.name + ": samples=" + strconv.Itoa(p.count()) +
		" mean=" + strconv.FormatFloat(p.mean(), 'f', 6, 32) +
		" max=" + strconv.FormatFloat(float64(p.max()), 'f', 6, 32) +
		" levels_per_s=" + strconv.FormatFloat(p.perSecond(), 'f', 2, 32)
}

// phaseRouter is the window the capture thread is currently writing into. The
// sink holds one pointer swap instead of three flags, because a phase that is
// still open after its window would silently steal samples from the next one.
type phaseRouter struct {
	mu  sync.RWMutex
	cur *livePhase
}

func (r *phaseRouter) set(p *livePhase) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cur = p
}

func (r *phaseRouter) route(level float32) {
	r.mu.RLock()
	cur := r.cur
	r.mu.RUnlock()
	if cur != nil {
		cur.take(level)
	}
}

// TestAC247LiveMicrophoneLevelsReachTheBallSeam reads the live level path in
// three windows: quiet, something loud played into the room, quiet again.
func TestAC247LiveMicrophoneLevelsReachTheBallSeam(t *testing.T) {
	if os.Getenv("WISP_LIVE_MIC") != "1" {
		t.Skip("AC#2 needs a real microphone: run with WISP_LIVE_MIC=1 after both tasklist rulers read 0")
	}
	requireNoCompetingAudioProcesses(t)

	rt, shut := bootAudioRuntime(t)
	// The unmuted value is a choice written into a temp data dir for this one
	// run, not a change to a shipped default (that one is AC#4's subject, pinned
	// by the sibling file): assembleCapture reads it through config.LoadFile
	// exactly as the resident boot does.
	dir := writeAudioConfig(t, func(c *config.Config) { c.Audio.MicMutedDefault = false })

	quietA := &livePhase{name: "quiet-A"}
	loud := &livePhase{name: "sound"}
	quietB := &livePhase{name: "quiet-B"}
	router := &phaseRouter{}
	router.set(quietA)

	ra := assembleCapture(rt, dir, router.route, newRealCaptureSource)
	if !ra.started {
		t.Fatalf("the live capture leg did not come up on this machine: %s", ra.posture())
	}
	mic := ra.mic.(*audio.WASAPIMicrophone)
	capDev, renderDev, err := mic.Endpoints()
	if err != nil {
		t.Fatalf("endpoint pair: %v", err)
	}
	t.Logf("AC#2 出处: capture=%q render=%q pinned_thread=%d", capDev.String(), renderDev.String(), mic.ThreadID())

	// P4 form 甲 on the live thread: booked in the runtime's registry, and not
	// double-booked into the process one.
	if n := rt.Registry.CountByName("audio-capture"); n != 1 {
		t.Fatalf("live capture thread in rt.Registry = %d, want 1", n)
	}
	if n := observe.Default.CountByName("audio-capture"); n != 0 {
		t.Fatalf("the live capture thread was ALSO booked into the process registry (%d)", n)
	}

	time.Sleep(3 * time.Second)
	quietA.span = 3 * time.Second

	router.set(loud)
	stopNoise := playAlarmThroughSpeakers(t)
	time.Sleep(6 * time.Second)
	loud.span = 6 * time.Second
	stopNoise()

	router.set(quietB)
	time.Sleep(3 * time.Second)
	quietB.span = 3 * time.Second

	st := mic.Stats()
	t.Logf("AC#2 counters at the seam: levels=%d frames_sent=%d frames_dropped=%d last_err=%q",
		ra.levels.Load(), st.FramesSent, st.FramesDropped, st.LastError)
	for _, p := range []*livePhase{quietA, loud, quietB} {
		t.Logf("AC#2 reading: %s", p.String())
	}

	if loud.count() == 0 || quietA.count() == 0 || quietB.count() == 0 {
		t.Fatalf("a phase recorded nothing: %s | %s | %s", quietA, loud, quietB)
	}
	if want := 1 / audio.FrameDuration.Seconds(); loud.perSecond() < want*0.6 || loud.perSecond() > want*1.6 {
		t.Fatalf("level cadence %.2f/s is not the frame cadence %.2f/s (P8 cost clause: delivery on the pinned thread)",
			loud.perSecond(), want)
	}
	if loud.mean() <= quietA.mean()*2 || loud.max() <= quietA.max() {
		t.Fatalf("the loud phase is not distinguishable from silence: %s | %s | %s"+
			" (speakers off, or routed to an endpoint this microphone cannot hear)",
			quietA, loud, quietB)
	}
	if quietB.mean() > loud.mean() {
		t.Fatalf("the control phase stayed as loud as the sound phase, so this is noise not a reading: %s | %s", loud, quietB)
	}

	if recs := shut(); len(recs) != 10 {
		t.Fatalf("shutdown trail = %d records, want the frozen 10", len(recs))
	}
	if tid := mic.ThreadID(); tid != 0 {
		t.Fatalf("capture thread id not released after D38(e) step 4: %d", tid)
	}
}

// requireNoCompetingAudioProcesses is the ticket's §排程 ruler: balldebug.exe and
// wisp.exe must both be absent before a live device test touches the
// microphone - otherwise the reading belongs to somebody else's frames.
func requireNoCompetingAudioProcesses(t *testing.T) {
	t.Helper()
	for _, name := range []string{"balldebug.exe", "wisp.exe"} {
		out, err := exec.Command("tasklist.exe", "/FI", "IMAGENAME eq "+name).Output()
		if err != nil {
			t.Fatalf("tasklist for %s: %v", name, err)
		}
		if strings.Contains(strings.ToLower(string(out)), strings.ToLower(name)) {
			t.Fatalf("%s is running; the live reading would be shared with it (ticket §排程 requires 0):\n%s",
				name, string(out))
		}
	}
}

// playAlarmThroughSpeakers puts a real sound in the room through the render
// endpoint and returns the stop function.
func playAlarmThroughSpeakers(t *testing.T) func() {
	t.Helper()
	media := filepath.Join(os.Getenv("WINDIR"), "Media", "Alarm01.wav")
	if _, err := os.Stat(media); err != nil {
		t.Skipf("no alarm sample to play out loud (%v): AC#2 then needs a person speaking into the microphone", err)
	}
	script := "$ErrorActionPreference='Stop';" +
		"$p=New-Object System.Media.SoundPlayer '" + media + "';" +
		"while ($true) { $p.PlaySync(); Start-Sleep -Milliseconds 250 }"
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", script)
	if err := cmd.Start(); err != nil {
		t.Fatalf("start the sound source: %v", err)
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
