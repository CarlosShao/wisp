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
	if *posX >= 0 && *posY >= 0 {
		debugMoveWindow(b.DebugHWND(), int32(*posX), int32(*posY))
	}
	fmt.Printf("balldebug: ball up handles=%d\n", handleCount())

	dwell := time.Duration(*cycleMs) * time.Millisecond

	switch {
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
