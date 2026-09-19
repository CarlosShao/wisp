// Command wisp is the single Wisp binary (SPEC-01 §3): no-args starts the GUI
// resident process, `wisp run` and `wisp doctor` are console subcommands whose
// output goes through AttachConsole.
//
// Ticket 03 scope: the no-args path boots the runtime skeleton (Job Object,
// per-session single instance, goroutine registry - the init-time self-checks)
// and runs an empty event loop until an exit signal, then exits through the
// frozen D38(e) 10-step shutdown order. The floating ball GUI is ticket 07.
// `wisp run` still echoes its task text as a placeholder (agent loop: ticket
// 10); `wisp doctor` reports the toolchain/DLL self-check (ticket 01).
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/proc"
	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
)

const usage = `wisp - personal voice agent for Windows

Usage:
  wisp             GUI resident process (boots the runtime skeleton, empty
                   event loop; the floating ball window is ticket 07)
  wisp run "task"  run a task from the command line (ticket 01: echo placeholder)
  wisp doctor      toolchain and native-DLL self-check, prints PASS/FAIL
  wisp slo         SLO sampling driver (ticket 08): one state per run, JSON
                   verdict; driven by scripts/slo-check.ps1
  wisp version     print version information
  wisp help        show this help

Environment: WISP_ENV in {prod|dev|test}, default from build (SPEC-03 §5).
`

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		// GUI entry point: boot the resident runtime and park in the empty
		// event loop until an exit signal arrives. Printed proof includes the
		// sherpa-onnx runtime version, which can only appear if the cgo link
		// against the sherpa-onnx C API actually works.
		attachParentConsole()
		runResident()
		return
	}
	switch args[0] {
	case "run":
		attachParentConsole()
		cmdRun(args[1:])
	case "doctor":
		attachParentConsole()
		if !cmdDoctor() {
			os.Exit(1)
		}
	case "slo":
		os.Exit(cmdSLO(args[1:]))
	case "version", "--version", "-v":
		attachParentConsole()
		printVersions("")
	case "help", "-h", "--help":
		attachParentConsole()
		fmt.Print(usage)
	default:
		attachParentConsole()
		fmt.Fprintf(os.Stderr, "wisp: unknown command %q\n\n", args[0])
		fmt.Print(usage)
		os.Exit(2)
	}
}

// runResident is the no-args path: boot (env layout, Job Object, single
// instance, goroutine registry), run the empty event loop, then exit through
// the D38(e) shutdown order. A second launch in the same session signals the
// running instance's activation event and exits (D42#7).
func runResident() {
	printVersions("")

	env, err := buildinfo.ResolveEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "wisp: %v\n", err)
		os.Exit(2)
	}

	rt, err := proc.Boot(env)
	if errors.Is(err, proc.ErrAlreadyRunning) {
		layout, lerr := proc.DefaultLayout(env)
		if lerr == nil && layout.MutexEnabled {
			_ = proc.SignalExistingInstance(layout.ActivateEventName)
		}
		fmt.Println("wisp: another instance is running in this session; activated it; exiting")
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "wisp: boot failed: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		records := rt.Shutdown(false)
		failed := 0
		for _, rec := range records {
			if rec.Err != nil {
				failed++
				fmt.Printf("wisp: shutdown step %d (%s) FAILED: %v\n", rec.Step, rec.Name, rec.Err)
			}
		}
		fmt.Printf("wisp: exited through the D38(e) shutdown order (10 steps, %d failed)\n", failed)
	}()

	sum := rt.Layout.Summary()
	fmt.Printf("wisp: resident runtime booted (%s, data dir = %s, portable = %v, job object = on, single instance = %v)\n",
		sum.EnvBadge(), sum.DataDir, sum.Portable, rt.Instance != nil)
	fmt.Printf("wisp: empty event loop running; the floating ball arrives in ticket 07 (Ctrl+C exits cleanly)\n")

	reason := rt.RunEventLoop()
	fmt.Printf("wisp: event loop ending (%s); running the D38(e) shutdown order\n", reason)
}

// cmdRun echoes the task text. Placeholder until the agent loop (ticket 10).
func cmdRun(args []string) {
	task := "no task text given"
	if len(args) > 0 {
		task = args[0]
	}
	printVersions("")
	fmt.Printf("wisp run: task text accepted: %q\n", task)
	fmt.Printf("wisp run: agent loop is ticket 10; this is the CLI plumbing placeholder\n")
}

// printVersions writes the common version block (prefix used by callers).
func printVersions(prefix string) {
	fmt.Printf("%swisp %s (%s, built %s)\n", prefix, buildinfo.Version, buildinfo.Commit, buildinfo.BuildDate)
	fmt.Printf("%sWISP_ENV=%s (data dir rules: SPEC-03 §5)\n", prefix, buildinfo.EnvString())
	fmt.Printf("%ssherpa-onnx runtime version: %s\n", prefix, sherpa.GetVersion())
}
