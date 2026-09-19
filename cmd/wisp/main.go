// Command wisp is the single Wisp binary (SPEC-01 §3): no-args starts the GUI
// resident process, `wisp run` and `wisp doctor` are console subcommands whose
// output goes through AttachConsole.
//
// Ticket 01 scope: build-chain proof of life. No-args prints version and
// environment to the console (the floating ball GUI is ticket 07), `wisp run`
// echoes its task text as a placeholder (the real agent loop is ticket 10),
// and `wisp doctor` reports the toolchain/DLL self-check.
package main

import (
	"fmt"
	"os"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
)

const usage = `wisp - personal voice agent for Windows

Usage:
  wisp             GUI resident process (ticket 01: proof-of-life printout)
  wisp run "task"  run a task from the command line (ticket 01: echo placeholder)
  wisp doctor      toolchain and native-DLL self-check, prints PASS/FAIL
  wisp version     print version information
  wisp help        show this help

Environment: WISP_ENV in {prod|dev|test}, default from build (SPEC-03 §5).
`

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		// GUI entry point. Ticket 01 keeps this a console proof-of-life; the
		// floating ball window is ticket 07. Printed proof includes the
		// sherpa-onnx runtime version, which can only appear if the cgo link
		// against the sherpa-onnx C API actually works.
		attachParentConsole()
		runProofOfLife()
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

// runProofOfLife is the no-args placeholder: prove the binary, its environment
// model, and the linked sherpa-onnx C API are alive.
func runProofOfLife() {
	printVersions("")
	fmt.Printf("wisp: GUI proof-of-life (floating ball lands in ticket 07)\n")
	fmt.Printf("wisp: data dir = %s\n", dataDirForDisplay())
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
	fmt.Printf("%sWISP_ENV=%s (data dir rules: SPEC-03 §5)\n", prefix, buildinfo.Env())
	fmt.Printf("%ssherpa-onnx runtime version: %s\n", prefix, sherpa.GetVersion())
}
