package manifestslot

// Ticket 52 — the D46 command plugin test-executable harness.
//
// A REAL child process is the only honest witness for four of this ticket's
// acceptance criteria: "no shell was spawned", "the child's environment is the
// allowlist", "the hanging CLI was killed with its tree" and "1MB of stdout was
// killed, not swallowed". So the test binary doubles as the fixture CLI: with a
// handshake in the fixtureHandshakeVar environment it behaves as the wrapped
// program, and every other run (the real test run) behaves as a test binary.
//
// The fixture is deliberately dumb and synchronous (no goroutines, so it cannot
// leave work behind when it is killed):
//
//	echo      prints "env=NAME=value" for every variable it can see, one per
//	          line, then exits 0. The env proof reads the INSTALLED process's
//	          own table, not the []string the host intended to pass.
//	huge      emits 8MB of stdout (past the 1MB cap), then exits 0.
//	hang      starts a grandchild (this same binary, fixtureChildArg) that writes
//	          its pid to file, then blocks. Proves the TREE dies.
//	extract   prints a JSON document (the extraction fixtures).

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

// fixtureHandshakeVar is how a fixture child learns its mode: the manifest
// declares the NAME in env_allowlist and the host supplies the value.
const fixtureHandshakeVar = "WISP52_FIXTURE"

// fixtureModePrefix keeps a stray variable from turning the test binary into a
// fixture: only a value that starts with it is honoured.
const fixtureModePrefix = "wisp52:"

// fixtureChildArg is the argument the fixture passes to its own grandchild.
const fixtureChildArg = "--wisp52-fixture-child"

// fixtureDelay443Var is ticket 443's amplification dial. The "hang" fixture
// spends this many milliseconds before it spawns the grandchild, which is the
// same budget a cold 90MB image (this test binary, re-executed twice) spends on
// its own on a loaded runner. Unset - the committed default - is 0ms, so the
// suite's own timing never depends on it; only the ticket-443 repro sets it.
//
// It exists because the failure ticket 443 fixes is a *scheduling* outcome: on
// a fast bench the fixture reports its grandchild pid every time, and "ran a
// hundred times green" would not be a reproduction of a runner that lost the
// race. The dial moves the fixture's startup across the 500ms budget, so the
// same red sentence CI produced can be produced here on demand.
//
// THE REPRO VALUE IS 1500 (WISP52_443_FIXTURE_DELAY_MS=1500): past the raced
// 500ms budget by 3x, and far enough below the abstain path's 4000ms control
// budget that the control still sees its pid after this 90MB image re-executes
// itself. Measured boundary: at 3000 the control run itself starts racing its
// own budget on a cold bench (first rep printed
// RED(443:timeout-tree:control-no-pid)), so 3000 is too big to be a repro and a
// value in that range fails honestly instead of silently weakening the control.
const fixtureDelay443Var = "WISP52_443_FIXTURE_DELAY_MS"

// fixtureDelay443 reads the dial. Unset, negative or unparsable is 0 (no
// delay), so a stray variable cannot slow the committed suite.
func fixtureDelay443() time.Duration {
	raw := strings.TrimSpace(os.Getenv(fixtureDelay443Var))
	if raw == "" {
		return 0
	}
	ms, err := strconv.Atoi(raw)
	if err != nil || ms <= 0 {
		return 0
	}
	return time.Duration(ms) * time.Millisecond
}

// TestMain runs the fixture when the process was started as one, and the tests
// otherwise.
//
// The handshake is ONE variable carrying "mode[:pidfile]", because an
// environment variable NAME cannot carry a path: the grandchild's pid file
// travels inside the value, so the allowlist (which lists names) has a single
// entry to declare.
func TestMain(m *testing.M) {
	if handshake := os.Getenv(fixtureHandshakeVar); strings.HasPrefix(handshake, fixtureModePrefix) {
		runFixtureCLI(handshake)
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// runFixtureCLI is the wrapped-program side of the test binary. handshake is
// "wisp52:<mode>[:<pidfile>]".
func runFixtureCLI(handshake string) {
	spec := strings.TrimPrefix(handshake, fixtureModePrefix)
	mode, pidFile, _ := strings.Cut(spec, ":")

	if len(os.Args) >= 3 && os.Args[1] == fixtureChildArg {
		// The grandchild: record the pid the host is supposed to reap.
		_ = os.WriteFile(os.Args[2], []byte(strconv.Itoa(os.Getpid())), 0o600)
		block()
	}

	switch mode {
	case "echo":
		for _, kv := range os.Environ() {
			fmt.Println("env=" + kv)
		}
	case "huge":
		// 8 x 1MB written straight through: no buffering games, so the host
		// sees a real stream that passes the cap.
		chunk := []byte(strings.Repeat("A", 1<<20))
		for i := 0; i < 8; i++ {
			if _, err := os.Stdout.Write(chunk); err != nil {
				return
			}
		}
	case "hang":
		if pidFile != "" {
			if delay := fixtureDelay443(); delay > 0 {
				time.Sleep(delay)
			}
			// The grandchild inherits the fixture's environment, which is how
			// it learns the marker; it gets the pid path as an ARGUMENT.
			cmd := exec.Command(os.Args[0], fixtureChildArg, pidFile)
			cmd.Env = os.Environ()
			_ = cmd.Start()
			waitForFile(pidFile)
		}
		block()
	case "extract":
		fmt.Print(`{"data":{"items":[{"summary":"standup 10:00"},{"summary":"review 14:00"}]},"ok":true}`)
	default:
		fmt.Fprintln(os.Stderr, "fixture CLI: unknown mode "+strconv.Quote(mode))
		os.Exit(2)
	}
}

// waitForFile polls for a file the grandchild writes. It is a fixture-side
// wait, not a contract timeout: the host's D46 budget is what ends the call.
func waitForFile(path string) {
	for i := 0; i < 500; i++ {
		if _, err := os.Stat(path); err == nil {
			return
		}
		sleepBrief()
	}
}

// sleepBrief is the fixture's between-polls pause.
func sleepBrief() { time.Sleep(10 * time.Millisecond) }

// block parks the fixture until something kills it.
func block() {
	for {
		time.Sleep(time.Hour)
	}
}
