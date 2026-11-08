// Command wisp is the single Wisp binary (SPEC-01 §3): no-args starts the GUI
// resident process, `wisp run` and `wisp doctor` are console subcommands whose
// output goes through AttachConsole.
//
// The no-args path (ticket 03 first, extended since) boots the runtime
// skeleton (Job Object, per-session single instance, goroutine registry - the
// init-time self-checks) and runs the resident event loop until an exit
// signal or a stop request fires, then exits through the frozen D38(e)
// 10-step shutdown order. The floating ball is not a pending ticket: the
// resident leg of this same process assembles it at startup
// (startResidentBall in resident_ball_windows.go, called from runResident).
// `wisp run` runs one task end to end through the agent loop (ticket 12);
// `wisp doctor` reports the toolchain/DLL self-check (ticket 01).
package main

import (
	"fmt"
	"os"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/buildinfo"
	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
)

const usage = `wisp - personal voice agent for Windows

Usage:
  wisp             GUI resident process: boots the runtime skeleton, hosts the
                   floating ball window, its tray icon and its three always-on
                   global hot keys in this same process (summon, mute, panel;
                   the cancel slot is the fourth key, bound only while a
                   confirmation card is waiting and handed back when it leaves,
                   so an idle wisp holds no bare key against the rest of the
                   desktop), shows confirmation cards on that ball, and takes
                   tasks from the operator's console when this process really
                   owns one (ticket 246 AC#7: no interactive console means the
                   task entry says so out loud and the ball keeps running;
                   Ctrl+C stops the process, and so does the tray's Exit item,
                   which asks the resident loop for the orderly D38(e) shutdown
                   instead of exiting in place (ticket 228 AC#11))
  wisp run "task"  run one task end to end through the agent loop (ticket 12:
                   the S1 text path - streamed reply + notification + exit code
                   that reflects the error class)
  wisp sessions    cross-restart conversation list (ticket 166 outlet #1,
                   wired by 364-r1a): sessions list prints the newest sessions
                   through the existing session.Manager.List - id, updated_at,
                   title
  wisp commands    the command roster (ticket 213 AC#4): prints every command
                   this machine can really execute, one row each - the name, what
                   it does, its argument shape, whether a human approval is
                   involved, and which existing door runs it. The rows come from
                   internal/command's production roster, so a capability with no
                   door cannot appear here
  wisp providers   provider catalog path (ticket 12, ruling A11): discover
                   <provider> lists /v1/models; probe <provider>/<model> runs
                   the capability probe suite into provider_health
  wisp secret      credential entry (ticket 63): set/get/list/unset a DPAPI
                   blob without the key ever entering argv, a log line or chat
  wisp config      config.toml keys from the console (ticket 65): get/set/list
                   over the keys this CLI owns - today ball.style (glass|
                   jarvis, the appearance switch 默认玻璃球), written through
                   the same guarded setter the settings page uses; the receipt
                   names the tier honestly (a CLI write lands in the file, a
                   running wisp picks the new style up at next boot)
  wisp models      signed model store (C29, ticket 121 AC#2): list the verified
                   manifest, re-check an installed model, or make one available
                   through the Downloading hand-off and its re-verification
  wisp doctor      toolchain and native-DLL self-check, prints PASS/FAIL
  wisp panel-assets  embedded panel bundle (ticket 77): -manifest lists what the
                   binary carries, -render <path> writes the bytes the WebView2
                   host would serve (proof the UI needs no node and no network)
  wisp panel-inbound  inbound panel dispatch leg (ticket 33 AC#9): reads one raw
                   composer envelope per stdin line and hands it to the inbound
                   router; -data names the data root (this leg never resolves one)
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
	case "sessions":
		// Ticket 364-r1a: the first production caller of session.Manager.List.
		// Attach the parent console like the other operator-driven legs so the
		// rows and the refusal line are readable from an Explorer start too.
		attachParentConsole()
		os.Exit(cmdSessions(args[1:], sessionsIO{stdout: os.Stdout, stderr: os.Stderr}))
	case "commands":
		// Ticket 213 AC#4: the command roster's first production caller. A console
		// leg that prints internal/command's production table and nothing else - it
		// opens no store, resolves no data root and calls no door, so
		// attachParentConsole is all it needs from this file.
		attachParentConsole()
		os.Exit(cmdCommands(args[1:], commandsIO{stdout: os.Stdout, stderr: os.Stderr}))
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
	case "config":
		// Ticket 65: the operator's own surface for the ball style switch -
		// get/set/list over the closed key roster (today ball.style). A console
		// leg like `wisp secret`: it resolves the env and data root itself
		// (cmdConfig), writes only through the guarded setter, and answers with
		// the tier the write really lands at. Covered by the legCovers133 nail
		// TestConfigLegGetSetGuarded65 -> cmdConfig.
		attachParentConsole()
		os.Exit(cmdConfig(args[1:], configIO{stdout: os.Stdout, stderr: os.Stderr}))
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
	case "panel-inbound":
		// Ticket 33 AC#9: the inbound router's first production listener. Console
		// leg (stdin envelopes in, sentences out), so attachParentConsole like the
		// other operator-driven commands; it writes only the -data root named here.
		attachParentConsole()
		os.Exit(cmdPanelInbound(args[1:], panelInboundIO{}))
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
//
// Ticket 201 added the answer side. attachParentConsole already ran at the call
// site above, so on a GUI-subsystem build the standard handles are the console's
// by the time interactiveStdin looks at them; it returns nil for a pipe, a
// redirected file or a start with no console at all, and that case is said out
// loud rather than left as a silent difference in posture (AC#7's 「明确且响亮」:
// with nobody able to answer, an L2 card resolves as a reject and an L1 window
// has no one to oppose it - both of which the operator is entitled to know
// before the run, not from an exit code afterwards).
func cmdRun(args []string) int {
	printVersions("")
	reply := interactiveStdin()
	if reply == nil {
		fmt.Fprintln(os.Stderr,
			"wisp run: 本机没有可交互控制台，本轮没有人能答复卡片："+
				"L2 卡会等到超时后按拒绝处理，L1 窗口没有人能否决（要能当场答复，请在终端里跑）")
	}
	spec := runSpec{
		argv:   args,
		stdout: os.Stdout,
		stderr: os.Stderr,
		reply:  reply,
	}
	// Ticket 170 form 甲 (owner-approved): a console whose stdin is a real TTY
	// IS a host that owns a cancel transport, so this leg declares one - but
	// ONLY then. The declaration is the host's statement, not the assembly's
	// inference: assembleRuntime cannot tell a TTY from a scripted test reader
	// (the injection seam AGENTS.md §1.3 names for `wisp run`), so the TTY fact
	// is carried by interactiveStdin here and lands on the gate through
	// runSpec.replyVeto, the same one door the ball/Esc-hook leg uses. A pipe,
	// a script or CI keeps the pre-170 posture: no veto channel declared, one
	// loud audit line at assembly (run.go's CLI-VETO-NONE), no pretended 停.
	if reply != nil {
		spec.replyVeto = declareConsoleCancelTransport()
	}
	return runTextTask(spec)
}

// consoleCancelKey is the key THIS host really holds for the cancel slot: on a
// console leg the operator stops an L1 window by feeding the reply stream - a
// bare Enter included (approval_reply.go's live-window veto). The card names
// the key through approval's read-time spelling seam, so what it prints is what
// this process really delivers, never a phantom Esc.
const consoleCancelKey = "回车"

func consoleCancelKeySpelling() string { return consoleCancelKey }

// declareConsoleCancelTransport is the two statements every host owes together
// (approval_reply.go: declaring a transport IS asserting it is up): the channel
// slot is returned for runSpec.replyVeto, and the cancel-key spelling is
// installed here so the card renders 「按 回车 键」 instead of naming a global
// Esc hook this process never registered. The slot itself is one of SPEC-06
// §2's closed four - the cancel-KEY channel (esc) - because the set is frozen
// and inventing a fifth selector is exactly the M-7/C-3 shape the approval
// package refuses; what the console borrows is the SLOT, and the spelling seam
// is what keeps the borrowed label honest.
func declareConsoleCancelTransport() approval.Channel {
	_ = approval.SetCancelKeySpelling(consoleCancelKeySpelling)
	return approval.ChannelEsc
}

// printVersions writes the common version block (prefix used by callers).
func printVersions(prefix string) {
	fmt.Printf("%swisp %s (%s, built %s)\n", prefix, buildinfo.Version, buildinfo.Commit, buildinfo.BuildDate)
	fmt.Printf("%sWISP_ENV=%s (data dir rules: SPEC-03 §5)\n", prefix, buildinfo.EnvString())
	fmt.Printf("%ssherpa-onnx runtime version: %s\n", prefix, sherpa.GetVersion())
}
