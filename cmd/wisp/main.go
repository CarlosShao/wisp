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
	"fmt"
	"os"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
)

const usage = `wisp - personal voice agent for Windows

Usage:
  wisp             GUI resident process (boots the runtime skeleton, empty
                   event loop; the floating ball window is ticket 07)
  wisp run "task"  run one task end to end through the agent loop (ticket 12:
                   the S1 text path - streamed reply + notification + exit code
                   that reflects the error class)
  wisp providers   provider catalog path (ticket 12, ruling A11): discover
                   <provider> lists /v1/models; probe <provider>/<model> runs
                   the capability probe suite into provider_health
  wisp secret      credential entry (ticket 63): set/get/list/unset a DPAPI
                   blob without the key ever entering argv, a log line or chat
  wisp models      signed model store (C29, ticket 121 AC#2): list the verified
                   manifest, re-check an installed model, or make one available
                   through the Downloading hand-off and its re-verification
  wisp doctor      toolchain and native-DLL self-check, prints PASS/FAIL
  wisp panel-assets  embedded panel bundle (ticket 77): -manifest lists what the
                   binary carries, -render <path> writes the bytes the WebView2
                   host would serve (proof the UI needs no node and no network)
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
	// Three branches below dispatch a command that no case in this package drives.
	// They are recorded here, in the file that owns the dispatch, because ticket
	// 133's leg census requires a dispatched command to be nailed, driven by a test
	// of its own, or ruled out in words by somebody who looked:
	//
	// WISP-LEG-COVERAGE-RULING: version prints buildinfo and the sherpa-onnx
	// runtime version and exits 0. It resolves no data root, opens no file and
	// books no record, so the failure this family keeps catching - a listener or a
	// verdict deleted in silence - has no surface here. scripts/build.ps1's own
	// version block is its consumer.
	// WISP-LEG-COVERAGE-RULING: help prints the usage block and exits 0. The block
	// is itself reconciled against the dispatch census by censusVsUsage133 in
	// leg_dispatch_gate_133_test.go, so a drift in what this branch promises is
	// red in the other clause rather than unwatched here.
	// WISP-LEG-COVERAGE-RULING: default is the refusal path: unknown command,
	// usage, exit 2. Driving it from a case here means an in-process call that
	// would exit the test binary, and the real-process form is ticket 128's
	// TestAC2RealProcessRefusesOnEveryLegWithoutAppData128, which covers the other
	// legs' refusals rather than this one.
	switch args[0] {
	case "run":
		attachParentConsole()
		os.Exit(cmdRun(args[1:]))
	case "providers":
		attachParentConsole()
		os.Exit(cmdProviders(args[1:], providersIO{}))
	case "doctor":
		attachParentConsole()
		if !cmdDoctor() {
			os.Exit(1)
		}
	case "secret":
		os.Exit(cmdSecret(args[1:]))
	case "models":
		// Ticket 121 AC#2: the model hand-off chain (models.Manager ->
		// DownloadingBridge.Run -> VerifyInstalled) reaches the product through
		// this line. attachParentConsole because, like `wisp secret`, it is a
		// command an operator runs from an Explorer-launched process; the
		// refused-hand-off verdict has to be readable somewhere.
		attachParentConsole()
		os.Exit(cmdModels(args[1:], modelsIO{stdout: os.Stdout, stderr: os.Stderr}))
	case "slo":
		os.Exit(cmdSLO(args[1:]))
	case "panel-assets":
		attachParentConsole()
		os.Exit(cmdPanelAssets(args[1:]))
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

// runResident lives in resident_windows.go (with a fail-closed counterpart in
// resident_other.go): it boots the Windows resident runtime - Job Object,
// per-session single instance, event loop - none of which exist off Windows.

// cmdRun is the S1 text path (ticket 12): runTextTask assembles the whole
// stack - credential store, config, provider, approval gate, host bridge,
// agent loop - and the exit code reflects the task's error class.
func cmdRun(args []string) int {
	printVersions("")
	return runTextTask(runSpec{argv: args, stdout: os.Stdout, stderr: os.Stderr})
}

// printVersions writes the common version block (prefix used by callers).
func printVersions(prefix string) {
	fmt.Printf("%swisp %s (%s, built %s)\n", prefix, buildinfo.Version, buildinfo.Commit, buildinfo.BuildDate)
	fmt.Printf("%sWISP_ENV=%s (data dir rules: SPEC-03 §5)\n", prefix, buildinfo.EnvString())
	fmt.Printf("%ssherpa-onnx runtime version: %s\n", prefix, sherpa.GetVersion())
}
