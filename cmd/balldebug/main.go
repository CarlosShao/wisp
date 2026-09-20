// Command balldebug is the ticket 07 debug harness: it boots the real ball
// window + state machine and either walks all 20 states on a timer (for the
// visual evidence script scripts/dev/ball-cycle.ps1) or runs interactively
// (hotkeys + tray + click all live, wired through the D43 machine).
//
// It prints the process handle count per step so the window-stack SLO gate
// (<600, docs/SLO.md) is measurable, and on exit destroys every window and
// COM object it created.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/statemachine"
)

// allStates is the D43 state order used by the visual cycle (boot -> idle ->
// session -> failure helpers), matching SPEC-08 §2.1 table order.
var allStates = []statemachine.State{
	statemachine.StateFirstRun,
	statemachine.StateSleeping,
	statemachine.StateArmed,
	statemachine.StateMuted,
	statemachine.StateListening,
	statemachine.StateThinking,
	statemachine.StateActing,
	statemachine.StateConfirming,
	statemachine.StateAwaitingApproval,
	statemachine.StateSpeaking,
	statemachine.StateSettling,
	statemachine.StateWarm,
	statemachine.StateConversation,
	statemachine.StateDownloading,
	statemachine.StateError,
	statemachine.StateNoNetwork,
	statemachine.StateUnconfigured,
	statemachine.StateWatchdogAlert,
	statemachine.StateQueued,
	statemachine.StateStuck,
}

var modPSAPI = windows.NewLazySystemDLL("kernel32.dll")

var procGetProcessHandleCount = modPSAPI.NewProc("GetProcessHandleCount")

func handleCount() uint32 {
	var n uint32
	r, _, _ := procGetProcessHandleCount.Call(uintptr(windows.CurrentProcess()), uintptr(unsafe.Pointer(&n)))
	if r == 0 {
		return 0
	}
	return n
}

func main() {
	stay := flag.Bool("stay", false, "interactive mode: hotkeys/tray/click live until tray Exit or Ctrl+C")
	shotDir := flag.String("shots", "", "Go-side composite captures per state into this dir")
	var seq int
	state := flag.String("state", "", "show one state for 2s (name as in D43)")
	cycleMs := flag.Int("cycle-ms", 2000, "per-state dwell for -state / full cycle")
	posX := flag.Int("x", -1, "debug: place the ball window at this x (with -y)")
	posY := flag.Int("y", -1, "debug: place the ball window at this y (with -x)")
	sizePx := flag.Int("size", 0, "configured orb size px 44..72 (0 = ball default)")
	look := flag.String("look", "", "liquid treatment: "+strings.Join(ball.LookNames(), "|")+" (empty = default)")
	frozen := flag.Bool("frozen", false, "render the frozen SPEC-08 §2.1 visuals instead of the ticket 62 prototype")
	hold := flag.Bool("hold", false, "with -state: stay in that state until Ctrl+C / tray Exit")
	statusPath := flag.String("status", "", "write 'ready=1 timers=<bool>' to this file (used by -diff)")
	diffDir := flag.String("diff", "", "differential evidence dir: alive-vs-dead screenshots + resource table")
	diffStates := flag.String("diff-states", "Sleeping", "comma-separated states for -diff")
	diffSample := flag.Duration("diff-sample", 3*time.Second, "CPU sampling window per state (ball alive)")
	diffMargin := flag.Int("diff-margin", 120, "screenshot padding around the window rect (px)")
	diffAmp := flag.Int("diff-amp", 6, "gain applied to the |alive-dead| diff image")
	diffDwell := flag.Duration("diff-dwell", 6*time.Second, "time the ball stays in state before each shot")
	level := flag.Float64("level", 0, "synthetic audio envelope 0..1 pushed to the ball at ~30fps (0 = feed nothing)")
	diffLevel := flag.Float64("diff-level", 0, "with -diff: the -level handed to each child")
	dock := flag.String("dock", "", "dock the ball to this edge the way a drag ending there would: "+
		"left|right|top|bottom (no timer: the dock is carried by mouse messages)")
	diffDock := flag.String("diff-dock", "", "with -diff: also shoot every state docked to this edge, as an extra row")
	tour := flag.Bool("tour", false, "guided walk for owner sign-off, in one run: liquid flow -> voice-driven "+
		"rotation -> idle border -> static Sleeping -> docked on all four edges -> popped back -> yours to play with")
	tourDwell := flag.Duration("tour-dwell", 5*time.Second, "with -tour: how long each numbered step holds")
	configPath := flag.String("config", "", "drive the ball's hotkeys from a real config.toml [hotkey] "+
		"section, live-reloaded by the ticket 64 bridge (empty = compiled defaults)")
	flag.Parse()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

	// The ticket 62 prototype is what this harness exists to show; -frozen
	// renders the still-frozen SPEC-08 §2.1 table for comparison.
	ball.EnablePrototypeVisuals(!*frozen)
	if *look != "" {
		if !ball.SetLook(*look) {
			fmt.Fprintf(os.Stderr, "balldebug: unknown -look %q (choose one of %s)\n",
				*look, strings.Join(ball.LookNames(), ", "))
			os.Exit(2)
		}
	}

	// The edge-dock flag is validated up front so a typo cannot silently leave
	// the ball floating in the middle of the screen while evidence is taken.
	var dockEdge ball.Edge
	if *dock != "" {
		e, ok := ball.EdgeByName(*dock)
		if !ok {
			fmt.Fprintf(os.Stderr, "balldebug: unknown -dock %q (left|right|top|bottom|none)\n", *dock)
			os.Exit(2)
		}
		dockEdge = e
	}
	if *diffDock != "" {
		if _, ok := ball.EdgeByName(*diffDock); !ok {
			fmt.Fprintf(os.Stderr, "balldebug: unknown -diff-dock %q (left|right|top|bottom)\n", *diffDock)
			os.Exit(2)
		}
	}

	// -diff is the parent role: it spawns one child balldebug per state, shoots
	// the same region alive and dead, and measures. No ball lives here.
	if *diffDir != "" {
		err := runDiff(diffOpts{
			dir:     *diffDir,
			states:  splitStates(*diffStates),
			dwell:   *diffDwell,
			sample:  *diffSample,
			x:       *posX,
			y:       *posY,
			margin:  *diffMargin,
			amplify: *diffAmp,
			size:    *sizePx,
			look:    *look,
			frozen:  *frozen,
			level:   *diffLevel,
			dock:    *diffDock,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "balldebug: %v\n", err)
			os.Exit(1)
		}
		return
	}

	fmt.Printf("balldebug: start handles=%d\n", handleCount())

	var (
		b    *ball.Ball
		m    *statemachine.Machine
		errs []error
		// Hotkey bridge teardown (only set with -config): cancel, then join.
		stopBridge func()
		joinBridge chan struct{}
	)

	// Tray Exit and Ctrl+C both land here; the teardown is common.
	exit := make(chan string, 1)

	m = statemachine.New(statemachine.Options{Initial: statemachine.StateSleeping})
	b, err := ball.New(ball.Options{
		SizePx:      *sizePx,
		Initial:     statemachine.StateSleeping,
		WindowTitle: os.Getenv(envBallTitle), // private title when spawned by -diff
		Events: ball.Events{
			OnClickBall:    func() { gesture(b, m, ball.EvClickBall) },
			OnSummonHotkey: func() { gesture(b, m, ball.EvSummonHotkey) },
			OnMuteHotkey:   func() { gesture(b, m, ball.EvMuteHotkey) },
			OnCancelHotkey: func() { gesture(b, m, ball.EvCancelHotkey) },
			OnPanelHotkey:  func() { gesture(b, m, ball.EvPanelHotkey) },
			OnTrayPanel:    func() { fmt.Println("tray: open panel (stub, ticket 33)") },
			OnTrayMute:     func() { gesture(b, m, ball.EvMuteHotkey) },
			OnTrayPauseWake: func() {
				fmt.Println("tray: pause-wake toggled (stub, ticket 41)")
			},
			OnTrayExit: func() {
				fmt.Println("tray: exit requested")
				exit <- "tray"
			},
			OnDragEnd: func() { fmt.Println("ball: position persisted") },
		},
		Registry: observe.Default,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "balldebug: ball.New failed: %v\n", err)
		os.Exit(1)
	}

	// Every [hotkey] outcome is printed, not swallowed: "occupied by another
	// program" and "never attempted" are different sentences for the user, and
	// a demo harness that hides them is how ticket 64's silent failure survived
	// a whole sign-off run (A1b).
	for _, line := range b.HotkeyReport().Problems() {
		fmt.Printf("balldebug: %s\n", line)
	}
	fmt.Printf("balldebug: hotkeys live=%d/4 %s\n", len(b.HotkeyReport().Live()), hotkeySummary(b))

	// -config: drive the ball's hotkeys from a REAL config.toml through the
	// ticket 64 bridge, so "the user edited [hotkey]" is testable on a running
	// instance instead of being taken on trust (A1). Without -config the ball
	// keeps the compiled defaults and no config poll runs.
	if *configPath != "" {
		mgr, errC := config.NewManager(*configPath, nil)
		if errC != nil {
			fmt.Fprintf(os.Stderr, "balldebug: config.NewManager(%s): %v\n", *configPath, errC)
			b.Close()
			os.Exit(1)
		}
		bridge := ball.NewHotkeyReloader(b, b.ConfiguredHotkeys(), func() ball.HotkeyConfig {
			h := mgr.Config().Hotkey
			return ball.ApplyHotkeyDefaults(ball.HotkeyConfig{
				Summon: h.Summon, Mute: h.Mute, Cancel: h.Cancel, Panel: h.Panel})
		})
		bridge.Refresh = func() error { _, err := mgr.CheckAndReload(); return err }
		mgr.OnReload = bridge.OnReload()
		pollCtx, stopPoll := context.WithCancel(context.Background())
		// runHotkeyBridge is panic-isolated and returns when ctx is cancelled;
		// the join below keeps the harness free of leaked goroutines (main owns
		// the goroutine, so stopBridge+joinBridge are its roster entry).
		pollDone := make(chan struct{})
		go runHotkeyBridge(bridge, pollCtx, time.Second, pollDone)
		stopBridge, joinBridge = stopPoll, pollDone
		fmt.Printf("balldebug: hotkey bridge polling %s\n", *configPath)
	}
	if *posX >= 0 && *posY >= 0 {
		debugMoveWindow(b.DebugHWND(), int32(*posX), int32(*posY))
	}
	// -dock docks through the very same commit path a drag ending at that edge
	// runs, so the evidence shows the shipped behaviour and not a mock. It
	// starts no timer, which is what the timers column of a docked row proves.
	if *dock != "" {
		if !b.DebugDock(dockEdge) {
			fmt.Fprintf(os.Stderr, "balldebug: -dock %s did not dock the ball\n", *dock)
			b.Close()
			os.Exit(1)
		}
		_, p, owns := b.Docked()
		fmt.Printf("balldebug: dock=%s owns=%v progress=%.2f handles=%d\n", *dock, owns, p, handleCount())
	}
	fmt.Printf("balldebug: ball up handles=%d\n", handleCount())

	// Optional synthetic voice envelope: only the scalar level crosses into
	// the ball (SetAudioLevel is the project's single render-side audio seam,
	// C25). The feeder is owned by main and joined before teardown. In -tour
	// the walk feeds each step on its own, so the quiet steps really are quiet.
	feedStop := make(chan struct{})
	feedDone := make(chan struct{})
	if *level > 0 && !*tour {
		go feedLevels(b, float32(*level), feedStop, feedDone)
	} else {
		close(feedDone)
	}

	dwell := time.Duration(*cycleMs) * time.Millisecond

	switch {
	case *tour:
		runTour(b, *tourDwell, float32(*level), exit)
	case *state != "":
		s := statemachine.State(*state)
		b.SetState(s)
		timers := b.DebugTimersAlive() // read on the UI thread, post-apply
		fmt.Printf("balldebug: state=%s handles=%d timers=%v\n", s, handleCount(), timers)
		writeStatus(*statusPath, s, timers)
		if *hold {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt)
			select {
			case <-sig:
				fmt.Println("exit: signal")
			case <-exit:
			}
		} else {
			time.Sleep(dwell)
		}
	case *stay:
		// Interactive: hotkeys/tray/click live; the machine's own timeout
		// timers drive the session tail (Warm 90s -> Settling -> ...).
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		select {
		case <-sig:
			fmt.Println("exit: signal")
		case <-exit:
		}
	default:
		// Full 20-state visual cycle (scripts/dev/ball-cycle.ps1 screenshots
		// each step by polling the window title marker).
		for _, s := range allStates {
			fmt.Printf("balldebug: state=%s handles=%d\n", s, handleCount())
			b.SetState(s)
			if s == statemachine.StateAwaitingApproval {
				b.SetBadge(3)
			}
			if s == statemachine.StateDownloading {
				b.SetProgress(0.69)
				b.SetBadgeText("69%")
			}
			if s == statemachine.StateConfirming {
				b.SetBadgeText("3")
			}
			time.Sleep(dwell)
			if *shotDir != "" {
				goScreenShot(b.DebugHWND(), filepath.Join(*shotDir, fmt.Sprintf("%02d-%s.png", seq, s)))
				seq++
			}
		}
	}

	// Stop feeding before the zero-timer assertions so no posted level task
	// interleaves with the final state switch.
	close(feedStop)
	<-feedDone

	// Stop polling the config before the final assertions: a rebind mid-flight
	// would unregister/re-register under the zero-timer check below.
	if stopBridge != nil {
		stopBridge()
		<-joinBridge
	}

	// Zero-timer assertion in Sleeping (ticket acceptance): back to Sleeping
	// and check no animation timer survives.
	b.SetState(statemachine.StateSleeping)
	time.Sleep(200 * time.Millisecond)
	if b.DebugTimersAlive() {
		errs = append(errs, fmt.Errorf("animation timer alive in Sleeping"))
	}

	b.Close()
	final := handleCount()
	fmt.Printf("balldebug: closed handles=%d\n", final)
	if final > 600 {
		errs = append(errs, fmt.Errorf("handle gate exceeded: %d > 600", final))
	}
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "balldebug: %v\n", e)
		}
		os.Exit(1)
	}
	fmt.Println("balldebug: OK")
}

// splitStates parses a comma-separated -diff-states value.
func splitStates(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// writeStatus publishes "the state is applied and rendered" to a file, which
// is how the -diff parent synchronises with a child it cannot signal.
func writeStatus(path string, s statemachine.State, timers bool) {
	if path == "" {
		return
	}
	body := fmt.Sprintf("state=%s ready=1 timers=%v at=%s\n", s, timers, time.Now().Format(time.RFC3339Nano))
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "balldebug: -status write: %v\n", err)
	}
}

// feedLevels pushes a syllabified synthetic envelope into Ball.SetAudioLevel
// (the render-side audio seam) at the capture cadence, so the liquid has a
// voice to follow. Owned by main: it returns when stop closes, and reports
// through done. A panic here must never take the harness down with it.
func feedLevels(b *ball.Ball, amp float32, stop <-chan struct{}, done chan<- struct{}) {
	defer close(done)
	defer func() {
		if r := recover(); r != nil {
			slog.Error("balldebug: level feeder panicked", "panic", r)
		}
	}()
	tk := time.NewTicker(33 * time.Millisecond) // ~30fps, the model's cap
	defer tk.Stop()
	var n int
	for {
		select {
		case <-stop:
			return
		case <-tk.C:
			n++
			v := amp
			if n%12 == 0 {
				v = amp / 4 // a breath between phrases: the border must glide in
			}
			b.SetAudioLevel(v)
		}
	}
}

// runTour is the owner-facing walkthrough for ticket 62 sign-off: one run, one
// numbered step at a time, in the order the behavior was asked for - summoned
// liquid flow, voice-driven rotation from a synthetic level, the idle border
// transition, the static Sleeping body, the tab on each of the four edges, the
// pop back out - and then it leaves the ball docked so the pointer can do the
// last step for real. Every dock and pop below goes through the shipped message
// paths (DebugDock runs the drag-end commit, DebugPop drives the hover ramp), so
// what the owner watches is the code that ships, not a mock.
//
// With -frozen the same walk shows the still-frozen SPEC-08 visuals instead:
// the liquid and the edge-dock are prototype-only, so the steps say so rather
// than pretending, which is what makes the two runs comparable on one desktop.
func runTour(b *ball.Ball, dwell time.Duration, level float32, exit chan string) {
	prototype := ball.PrototypeVisualsEnabled()
	amp := level
	if amp <= 0 {
		amp = 0.7 // a default voice to follow: -level overrides it
	}
	n := 0
	say := func(cue string) {
		n++
		fmt.Printf("\ntour %02d. %s\n", n, cue)
		fmt.Printf("tour %02d. holding %s - tray Exit or Ctrl+C stops the walk at any time\n", n, dwell)
	}
	report := func(tag string) {
		e, p, owns := b.Docked()
		fmt.Printf("tour: %s timers=%v dock edge=%s progress=%.2f owns=%v handles=%d\n",
			tag, b.DebugTimersAlive(), edgeLabel(e), p, owns, handleCount())
	}
	// feed runs the synthetic envelope for one step, owned and joined here.
	feed := func(d time.Duration) {
		stop := make(chan struct{})
		done := make(chan struct{})
		go feedLevels(b, amp, stop, done)
		time.Sleep(d)
		close(stop)
		<-done
		b.SetAudioLevel(0)
	}

	fmt.Printf("balldebug: tour start mode=%s dwell=%s level=%.2f handles=%d\n",
		map[bool]string{true: "prototype", false: "frozen SPEC-08"}[prototype], dwell, amp, handleCount())

	if !prototype {
		for _, s := range []statemachine.State{
			statemachine.StateSleeping, statemachine.StateArmed, statemachine.StateListening,
			statemachine.StateSpeaking, statemachine.StateWarm, statemachine.StateError,
		} {
			say(fmt.Sprintf("frozen SPEC-08 look, state %s: dot + status ring, no liquid, no dock.", s))
			b.SetState(s)
			report(string(s))
			time.Sleep(dwell)
		}
		say("end of the frozen walk. Run the same command without -frozen to see the ticket 62 prototype.")
		waitForExit(exit)
		return
	}

	b.SetState(statemachine.StateSleeping)
	time.Sleep(300 * time.Millisecond)
	say("AT REST (Sleeping). This is the ball sitting on the desktop doing nothing: a static glass orb, " +
		"no animation, no timers. Look at how clearly you can read it against the wallpaper.")
	report("sleeping")
	time.Sleep(dwell)

	b.SetState(statemachine.StateListening)
	say("SUMMONED (Listening). This is what the hot key does: the liquid inside starts flowing and the ball " +
		"brightens. Nothing is being said yet, so the flow is the model's own burst, and it retires itself.")
	report("listening")
	time.Sleep(dwell)

	b.SetState(statemachine.StateSpeaking)
	say(fmt.Sprintf("SPEAKING WITH A VOICE FED IN (synthetic level %.2f). The liquid rotates and swells with the "+
		"level; every third phrase is quieter, so you can watch it follow the envelope rather than spin blindly.", amp))
	feed(dwell)
	report("speaking")

	b.SetState(statemachine.StateListening)
	say("STOPPED TALKING, STILL LISTENING. Nothing is being said, so instead of cutting, the liquid settles and " +
		"the border glides in around the ball. That arriving border is the transition the owner asked for.")
	time.Sleep(dwell * 2)
	report("idle-border")

	b.SetState(statemachine.StateSleeping)
	say("BACK TO REST (Sleeping). The border and the liquid both release: one static frame again, zero timers. " +
		"If this frame looks different from step 01, something is wrong.")
	report("sleeping-again")
	time.Sleep(dwell)

	for _, name := range []string{"left", "right", "top", "bottom"} {
		edge, _ := ball.EdgeByName(name)
		say("DOCKED ON THE " + strings.ToUpper(name) + " EDGE. The ball has squeezed itself into a tab against " +
			"the screen boundary - only a sliver of it is left showing, which is the part you can still hit.")
		if !b.DebugDock(edge) {
			fmt.Println("tour: the ball would not dock there (the monitor layout moved?) - skipping the rest of the dock steps")
			break
		}
		report("docked-" + name)
		time.Sleep(dwell)

		say("POPPED BACK OUT OF THE " + strings.ToUpper(name) + " EDGE. Same ramp the pointer running over the " +
			"tab drives, walked to its free end: the ball is round again and wholly on screen, so it can never " +
			"leave you with a part of itself you cannot reach.")
		if p := b.DebugPop(); p != 0 {
			fmt.Printf("tour: the pop-out stopped short at level %v\n", p)
		}
		report("popped-" + name)
		time.Sleep(dwell)
	}

	say("NOW YOUR TURN, FOR REAL. The ball is docked on the right edge as a tab: move the pointer onto the " +
		"sliver and it pops out, move away and it closes again, click it and it stays out, drag it into the " +
		"middle of the screen and let go and it is free.")
	if !b.DebugDock(ball.EdgeRight) {
		fmt.Println("tour: could not leave the ball docked for the interactive part")
	}
	report("yours")
	waitForExit(exit)
}

// waitForExit blocks until the tray Exit item or Ctrl+C lands, whichever the
// owner reaches for first.
func waitForExit(exit chan string) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	select {
	case <-sig:
		fmt.Println("exit: signal")
	case why := <-exit:
		fmt.Printf("exit: %s\n", why)
	}
	signal.Stop(sig)
}

// edgeLabel names an edge for the tour log (ball.Edge is an int with no String).
func edgeLabel(e ball.Edge) string {
	switch e {
	case ball.EdgeLeft:
		return "left"
	case ball.EdgeRight:
		return "right"
	case ball.EdgeTop:
		return "top"
	case ball.EdgeBottom:
		return "bottom"
	}
	return "none"
}

// gesture translates a ball gesture into machine events by context (the
// ticket 07 wiring: click = summon/veto, mute toggles, cancel vetoes).
func gesture(b *ball.Ball, m *statemachine.Machine, k ball.EventKind) {
	switch k {
	case ball.EvClickBall, ball.EvSummonHotkey:
		switch m.State() {
		case statemachine.StateListening, statemachine.StateConfirming:
			dispatch(b, m, statemachine.EvVeto, nil)
		case statemachine.StateSpeaking:
			dispatch(b, m, statemachine.EvInterrupt, nil)
		default:
			dispatch(b, m, statemachine.EvSummon, nil)
		}
	case ball.EvMuteHotkey:
		dispatch(b, m, statemachine.EvMuteKey, nil)
	case ball.EvCancelHotkey:
		if m.State() == statemachine.StateConfirming {
			dispatch(b, m, statemachine.EvVeto, nil)
		} else {
			dispatch(b, m, statemachine.EvVeto, nil) // Listening veto
		}
	case ball.EvPanelHotkey:
		fmt.Println("hotkey: panel (stub, ticket 33)")
	}
}

// dispatch runs one machine event and mirrors the resulting state (and the
// B1 Esc takeover window) onto the ball.
func dispatch(b *ball.Ball, m *statemachine.Machine, ev statemachine.Event, f *statemachine.Facts) {
	to, err := m.Dispatch(ev, f)
	if err != nil {
		fmt.Printf("machine: rejected %s in %s\n", ev, m.State())
		return
	}
	syncMachineToBall(b, m)
	if to != statemachine.StateConfirming {
		_ = to
	}
}

// syncMachineToBall mirrors machine state + overlay facts to the ball.
func syncMachineToBall(b *ball.Ball, m *statemachine.Machine) {
	s := m.State()
	b.SetState(s)
	// B1: Esc is the cancel key while Confirming, returned afterwards.
	if s == statemachine.StateConfirming {
		b.TakeEscForCancel()
		b.SetBadgeText("3")
	} else {
		b.ReleaseEscAfterSession()
	}
}

// hotkeySummary prints one line per [hotkey] slot: the configured binding and
// what Win32 actually made of it. The harness needs it because a ball that came
// up with none of its four keys registered looked exactly like a working demo
// (ticket 64 A1b): registerAll used to slog.Warn and drop the id, and nothing
// downstream read the report.
func hotkeySummary(b *ball.Ball) string {
	lines := make([]string, 0, 4)
	for _, bd := range b.HotkeyReport().Bindings() {
		binding := bd.Binding
		if binding == "" {
			binding = "(unset)"
		}
		lines = append(lines, fmt.Sprintf("%s=%s %s", bd.Name, binding, bd.Status))
	}
	return strings.Join(lines, "\n  ")
}

// runHotkeyBridge is the -config poll goroutine's body. main owns it (cancel
// through the stop func, join through done) and a panic in one tick is reported
// instead of taking the harness down; HotkeyReloader.Run returns when ctx is
// cancelled, so the join cannot hang.
func runHotkeyBridge(r *ball.HotkeyReloader, ctx context.Context, every time.Duration, done chan<- struct{}) {
	defer close(done)
	defer func() {
		if rec := recover(); rec != nil {
			slog.Error("balldebug: hotkey bridge panicked; hotkeys unchanged", "panic", rec)
		}
	}()
	r.Run(ctx, every)
}
