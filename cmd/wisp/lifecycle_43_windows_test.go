//go:build windows

package main

// Ticket 43, the suspend/resume/session-end sequences, driven through the real
// components the resident boot assembles:
//
//   - a REAL approval gate (approval.New) with an injected UI that records the
//     card the way ballCardUI does - so the pending L2 card is a REAL queue
//     item, and the void really walks the queue's refusal funnel (D43 row 24),
//     not a ledger-shaped stand-in;
//   - a REAL audio HalfDuplexGate over a counting source, so "capture stopped
//     before sleep" and "capture reopened after wake" are the gate's own
//     Start/Stop calls, counted;
//   - a REAL session.NewScope residency over a counting ASR engine, so
//     "no models resident during sleep" is EnginesHot=0 plus the engine's own
//     Close having run - the sampler face of the ticket's cell.
//
// The dispose runs ON THE VOICE OWNER (the same queue.run marshalling the
// production owner uses), so the test also pins the threading shape 381/15e
// made load-bearing.

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/audio"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/session"
	"github.com/CarlosShao/wisp/internal/speech"
	"github.com/CarlosShao/wisp/internal/tools"
)

// ---- the counting ASR engine -----------------------------------------------

type fakeAsr43 struct {
	closed atomic.Int64
}

func (a *fakeAsr43) Begin()          {}
func (a *fakeAsr43) Feed([]float32)  {}
func (a *fakeAsr43) Partial() string { return "" }
func (a *fakeAsr43) Finish() string  { return "" }
func (a *fakeAsr43) InFlight() bool  { return false }
func (a *fakeAsr43) Fed() int        { return 0 }
func (a *fakeAsr43) Close() error    { a.closed.Add(1); return nil }

// ---- the counting capture source -------------------------------------------

type fakeCaptureSource43 struct {
	starts atomic.Int64
	stops  atomic.Int64
}

func (s *fakeCaptureSource43) Start(ctx context.Context, buf chan<- []byte) error {
	s.starts.Add(1)
	return nil
}
func (s *fakeCaptureSource43) Stop() error        { s.stops.Add(1); return nil }
func (s *fakeCaptureSource43) Stats() audio.Stats { return audio.Stats{} }
func (s *fakeCaptureSource43) Err() error         { return nil }

// ---- the recording UI (the real queue-item path) ----------------------------

type recordingUI43 struct {
	ra *residentApproval
}

func (u *recordingUI43) Prompt(_ context.Context, p approval.Prompt) error {
	u.ra.cards.Record(p)
	return nil
}

func (u *recordingUI43) Update(_ context.Context, e approval.Event) error {
	if e.Kind == approval.EventDismissed || e.Kind == approval.EventStarted {
		u.ra.cards.Forget(e.CorrelationID)
	}
	return nil
}

// ---- the fixture ------------------------------------------------------------

type suspendFixture43 struct {
	ra     *residentApproval
	raudio *residentAudio
	rv     *residentVoice
	asr    *fakeAsr43
	source *fakeCaptureSource43
	askRc  chan tools.Answer
	askWhy chan string
}

func newSuspendFixture43(t *testing.T) *suspendFixture43 {
	t.Helper()
	root, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	f := &suspendFixture43{asr: &fakeAsr43{}, source: &fakeCaptureSource43{}}

	f.ra = &residentApproval{root: root, cancel: cancel}
	f.ra.ui = &ballCardUI{ra: f.ra}
	f.ra.grants = &residentGrantHolder{}
	f.ra.gate = approval.New(approval.Options{
		UI:       &recordingUI43{ra: f.ra},
		Channels: approval.NewChannels(),
		Logf:     func(string, ...any) {},
	})
	f.ra.cards = approval.NewReplies()
	f.ra.cards.Attach(approval.HostBinding{Gate: f.ra.gate})

	// The capture leg: the real gate over the counting source, started the way
	// assembleCapture starts it (ctx, bounded frames, WithStartMuted=false so
	// the inner source really opens).
	frames := audio.NewBoundedFrames()
	ctx, cancelCap := context.WithCancel(context.Background())
	t.Cleanup(cancelCap)
	f.raudio = &residentAudio{frames: frames, cancel: cancelCap}
	f.raudio.gate = audio.NewHalfDuplexGate(f.source, audio.PathT)
	if err := f.raudio.gate.Start(ctx, frames); err != nil {
		t.Fatalf("gate.Start: %v", err)
	}
	f.raudio.started.Store(true)

	// The voice residency: a real session scope over the counting engine, and
	// the owner goroutine that drains the same queue production drains.
	f.rv = &residentVoice{queue: newVoiceOwnerQueue()}
	f.rv.residency = newVoiceResidency(voicePlan{AsrDir: "stub-43"}, f.rv.queue, nil)
	f.rv.residency.asrLoader = func() (speech.AsrEngine, error) { return f.asr, nil }
	observe.Default.Spawn("lifecycle-43-owner", "test", nil, func(ctx context.Context) {
		wake := f.rv.queue.wakeChan()
		for {
			select {
			case <-ctx.Done():
				return
			case <-wake:
				f.rv.queue.run()
			}
		}
	})

	f.askRc = make(chan tools.Answer, 1)
	f.askWhy = make(chan string, 1)
	return f
}

// raiseL2 pushes ONE REAL queue item through the gate: PendingApproval is the
// exported production entry (askConfirmation's own L2 branch), the injected UI
// records the card, and the ask parks on the queue's answer channel exactly as
// it does under a live ball.
func (f *suspendFixture43) raiseL2(t *testing.T, corr string) {
	t.Helper()
	taskID := "task-43-" + corr
	// The D47 admit guard is the production precondition for an L2 card (a
	// task the text loop never admitted cannot be asked about) - the fixture
	// registers its task exactly the way execute() does.
	f.ra.gate.AdmitTextTask(taskID)
	d := tools.Decision{
		TaskID: taskID, Tool: "fs.write", Level: risk.L2,
		Reason:        "票 43 挂起作废用例的待确认操作",
		CorrelationID: corr,
	}
	observe.Default.Spawn("lifecycle-43-ask", "test", nil, func(ctx context.Context) {
		ans, why := f.ra.gate.PendingApproval(ctx, d)
		f.askRc <- ans
		f.askWhy <- why
	})
}

func waitPending43(t *testing.T, ra *residentApproval, want int) []approval.ReplyCard {
	t.Helper()
	var cards []approval.ReplyCard
	deadline := observe.NewTimeout(10 * time.Second)
	for {
		cards = ra.cards.Pending()
		if len(cards) >= want {
			return cards
		}
		if deadline.Expired() {
			t.Fatalf("pending cards = %d after 10s, want %d", len(cards), want)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// ---- the cases --------------------------------------------------------------

// Test43SuspendSequence is the ticket's cell 1 and cell 3 in one sequence:
// pending L2 voided + counted, capture stopped, engines released pre-sleep -
// and the resume half reopens the capture.
func Test43SuspendSequenceVoidsDisposesAndStopsCapture(t *testing.T) {
	f := newSuspendFixture43(t)

	// One REAL pending L2 card (queue item, nonce, the works).
	f.raiseL2(t, "corr-43-suspend")
	cards := waitPending43(t, f.ra, 1)
	if cards[0].Level != "L2" {
		t.Fatalf("pending card level = %q, want L2", cards[0].Level)
	}

	// One L1 window on the ledger (its route has no reject verb - the same
	// shape cancelTaskRoots abandons).
	f.ra.cards.Record(approval.Prompt{CorrelationID: "corr-43-l1", Tool: "fs.write", Level: "L1"})
	waitPending43(t, f.ra, 2)

	// An engine is hot before the sleep: the whole point of the cell is what
	// happens to it.
	if _, _, err := f.rv.residency.acquireTurn(session.EngineASR, "43-suspend"); err != nil {
		t.Fatalf("acquireTurn: %v", err)
	}
	if got := f.rv.residency.stats().EnginesHot; got != 1 {
		t.Fatalf("EnginesHot before suspend = %d, want 1 (a hot model)", got)
	}
	// The sampler face: the engine's bytes are in this process now; the reading
	// after the sequence is what "no models resident" means.
	if f.asr.closed.Load() != 0 {
		t.Fatal("the engine was closed before the suspend sequence ran")
	}

	voided := runSuspendSequence(f.ra, f.raudio, f.rv, "43-suspend")
	if voided != 2 {
		t.Fatalf("voided = %d, want 2 (one L2 refused + one L1 window voided)", voided)
	}

	// Cell 3: the L2 card is GONE from the ledger, and the parked ask received
	// the refusal - the queue's own funnel ran, not a ledger edit.
	if left := f.ra.cards.Pending(); len(left) != 0 {
		t.Fatalf("pending cards after suspend = %d, want 0 (voided by the suspend)", len(left))
	}
	select {
	case ans := <-f.askRc:
		if ans != tools.AnswerReject {
			t.Fatalf("parked ask answered %s, want reject", ans)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("the parked PendingApproval never came back from the suspend void")
	}
	if why := <-f.askWhy; !strings.Contains(why, "挂起") {
		t.Fatalf("refusal reason = %q, want the suspend wording", why)
	}

	// Cell 1, the dispose half: the scope is empty and the engine's Close ran.
	st := f.rv.residency.stats()
	if st.EnginesHot != 0 {
		t.Fatalf("EnginesHot after suspend = %d, want 0 (no models resident during sleep)", st.EnginesHot)
	}
	if f.asr.closed.Load() != 1 {
		t.Fatalf("engine closes = %d, want exactly 1 (released once, pre-sleep)", f.asr.closed.Load())
	}
	if f.source.stops.Load() != 1 {
		t.Fatalf("capture stops = %d, want 1 (device off before sleep)", f.source.stops.Load())
	}
	if f.source.starts.Load() != 1 {
		t.Fatalf("capture starts after suspend = %d, want 1 (no reopen happened during suspend)", f.source.starts.Load())
	}

	// The resume half: capture reopened through the same gate.
	if !runResumeSequence(f.raudio, "43-resume") {
		t.Fatal("runResumeSequence reported no reopen")
	}
	if f.source.starts.Load() != 2 {
		t.Fatalf("capture starts after resume = %d, want 2 (device reopened post-wake)", f.source.starts.Load())
	}

	// The resume-e2e face, in-process: after the suspend cleared the epoch, the
	// NEXT acquire really reopens a fresh epoch through the cold path - the
	// chain the wake depends on is intact, not merely empty. The strong counter
	// is the residency's loaderCallCount (it spans epochs): one load before the
	// suspend, one after the resume.
	if _, _, err := f.rv.residency.acquireTurn(session.EngineASR, "43-resume-reacquire"); err != nil {
		t.Fatalf("post-resume acquireTurn: %v", err)
	}
	if got := f.rv.residency.loaderCallCount(); got != 2 {
		t.Fatalf("loader calls after the post-resume acquire = %d, want 2 (one pre-suspend, one post-resume cold load)", got)
	}
	st = f.rv.residency.stats()
	if st.ColdLoads != 1 || st.EnginesHot != 1 {
		t.Fatalf("new-epoch stats = ColdLoads %d EnginesHot %d, want 1/1 (a fresh epoch, loaded once)", st.ColdLoads, st.EnginesHot)
	}
}

// Test43SessionEndSetsTheFastExitFlag pins the session-end half: the flag the
// boot's deferred rt.Shutdown reads is set by the sink the door calls, and the
// stop is a REQUEST on the runtime's own bit (the D38(e) chain still hangs off
// the boot's defer, never off the callback).
func Test43SessionEndSetsTheFastExitFlag(t *testing.T) {
	var fast atomic.Bool
	sink := wireLifecycleSink(nil, nil, nil, nil, &fast)
	if fast.Load() {
		t.Fatal("the fast flag started set")
	}
	sink.sessionEnd("43-test")
	if !fast.Load() {
		t.Fatal("the fast flag was never set: the deferred rt.Shutdown would run the slow path")
	}
}

// Test43PowerDeliveryFunnelIsSerial drives both power doors' funnel: the
// callback seam and the events it posts are the production body, and the pump
// executes the sink in order.
func Test43PowerDeliveryFunnelRunsTheSink(t *testing.T) {
	var fast atomic.Bool
	var suspended, resumed atomic.Int64
	sink := lifecycleSink{
		suspend:    func(string) { suspended.Add(1) },
		resume:     func(string) { resumed.Add(1) },
		sessionEnd: func(string) { fast.Store(true) },
	}
	le := &lifecycleEvents{sink: sink, fastExit: &fast, events: make(chan lifecycleEvent, 8)}

	le.deliverPowerEvent(pbtApmSuspend)
	le.deliverPowerEvent(pbtApmResumeAuto)
	le.run(<-le.events)
	le.run(<-le.events)

	if suspended.Load() != 1 || resumed.Load() != 1 {
		t.Fatalf("suspended=%d resumed=%d, want 1/1 (both doors funnel through one channel)", suspended.Load(), resumed.Load())
	}
	if le.suspendDepth.Load() != 1 || le.resumeDepth.Load() != 1 {
		t.Fatalf("depth counters suspend=%d resume=%d, want 1/1", le.suspendDepth.Load(), le.resumeDepth.Load())
	}
}
