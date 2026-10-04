//go:build windows && winlive

package main

// Ticket 258 AC#1/AC#2/AC#3, live tier: the shipped no-args process reads
// [hotkey] at boot, says where its four bindings came from, rebinds on a hand
// edit without a restart, and names an occupied combination instead of keeping
// the old value quietly.
//
// Run it alone, on a machine with a desktop and no other Wisp ball (the same
// rule resident_ball_live_228_windows_test.go states; the guard is
// internal/ball's live_guard_windows_test.go):
//
//	PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" \
//	  go test -tags winlive ./cmd/wisp -run TestLive258 -v
//
// The physical-press case (does Ctrl+Alt+R really fire after the rebind)
// belongs to internal/ball's live suite, which injects real input; this file
// reads the process's own sentences and Win32's registration set, which is the
// half a shipped process owes the operator.

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/ball"
	"golang.org/x/sys/windows"
)

// TestLive258ResidentBootNamesProvenanceAndDefaults is AC#1's live sentence:
// with NO config.toml in the data dir, the shipped process must fall back to
// the compiled defaults AND say so - "hotkeys from defaults" in the verdict.
// The boot then goes on without the file (a missing config is not a fatal for
// this leg); the case lets the process leave through its own shutdown path.
func TestLive258ResidentBootNamesProvenanceAndDefaults(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	leg := bootResidentLeg(t, exe, dataDir)
	t.Cleanup(leg.stop)

	saw := pollUntil127(200, func() bool {
		out := leg.stdout.String()
		return strings.Contains(out, ballUpClaim) || strings.Contains(out, ballAbsentClaim)
	})
	if !saw {
		t.Fatalf("AC#1 LIVE RED: the process never printed its ball posture.\n%s", leg.console())
	}
	out := leg.stdout.String()
	if strings.Contains(out, ballAbsentClaim) {
		t.Skipf("SKIP-LOUD: this host reports no ball window (%q), so the hotkey sentence has no subject here", ballAbsentClaim)
	}
	if !strings.Contains(out, "hotkeys from defaults") {
		t.Fatalf("AC#1 RED: a boot with no config.toml did not name the fallback in its verdict. Want \"hotkeys from defaults\" in:\n%s", out)
	}
	if i := strings.Index(out, "hotkeys from"); i >= 0 {
		t.Logf("AC#1 LIVE verdict fragment: %s", out[i:min64(len(out), i+120)])
	}
	// The bridge line: with no config file the src answers empty every tick, the
	// bridge keeps the applied set (nothing changed), and the armed line is the
	// one sentence about the wiring itself.
	if !strings.Contains(out, "hotkey reload bridge armed") {
		t.Errorf("AC#1 RED: the bridge-armed sentence is missing on a boot with a live ball:\n%s", out)
	}
	if !strings.Contains(out, "provenance=defaults") {
		t.Errorf("AC#1 RED: the bridge line does not carry the fallback provenance word:\n%s", out)
	}
}

// TestLive258ResidentBootTakesConfigHotkeys is AC#1's live positive control: a
// data dir whose config.toml names non-default combinations boots the ball
// registered on THOSE combinations, provenance "config".
func TestLive258ResidentBootTakesConfigHotkeys(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	const summon258 = "Ctrl+Alt+J"
	body := "[hotkey]\nsummon = \"" + summon258 + "\"\nmute = \"Ctrl+Alt+M\"\ncancel = \"Esc\"\npanel = \"Ctrl+Alt+P\"\n"
	if err := os.WriteFile(filepath.Join(dataDir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatalf("plant config: %v", err)
	}
	leg := bootResidentLeg(t, exe, dataDir)
	t.Cleanup(leg.stop)

	saw := pollUntil127(200, func() bool {
		out := leg.stdout.String()
		return strings.Contains(out, ballUpClaim) || strings.Contains(out, ballAbsentClaim)
	})
	if !saw {
		t.Fatalf("AC#1 LIVE RED: no ball posture sentence.\n%s", leg.console())
	}
	out := leg.stdout.String()
	if strings.Contains(out, ballAbsentClaim) {
		t.Skipf("SKIP-LOUD: no ball window on this host (%q)", ballAbsentClaim)
	}
	if !strings.Contains(out, "hotkeys from config") {
		t.Fatalf("AC#1 RED: a boot with planted [hotkey] values did not name config as the source. Want \"hotkeys from config\":\n%s", out)
	}
	if !strings.Contains(out, "summon="+summon258) {
		t.Fatalf("AC#1 RED: the per-slot summary does not name the planted summon %s:\n%s", summon258, out)
	}
}

// TestLive258ResidentRebindsWithoutRestart is AC#2's live form: the process is
// UP on a planted non-default summon, then config.toml is edited underneath
// it; the poll (1s) must rebind the live keys inside the same process.
//
// THE RULER WAS FIXED IN 258-R2 (258-v1 §5 condition (i)): the judge used to
// grep the console equals shape (summon=VALUE) off the sink, but the sink is a
// slog JSON file and the rebind record spells it "summon":"VALUE" - so the
// read could never hit and the case failed on a working line. The reading now
// goes through sinkRebindLineHas258 (cmd/wisp/resident_hotkey_258_test.go,
// pinned in CI by Test258SinkRebindRulerReadsBothSpellings), which judges the
// record in either spelling and reads "absent" as absent. What cannot be read
// at all - a sink that never opened - fails named, with the tail attached, not
// silently.
func TestLive258ResidentRebindsWithoutRestart(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	const newSummon258 = "Ctrl+Alt+R"
	// Boot summon is Z, a NON-default: with Q the boot summary below could not
	// tell the config path from the fallback, and the ruler would read true on
	// the M3-forever-defaults rot shape too.
	body := "[hotkey]\nsummon = \"Ctrl+Alt+Z\"\nmute = \"Ctrl+Alt+M\"\ncancel = \"Esc\"\npanel = \"Ctrl+Alt+P\"\n"
	if err := os.WriteFile(filepath.Join(dataDir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatalf("plant config: %v", err)
	}
	leg := bootResidentLeg(t, exe, dataDir)
	t.Cleanup(leg.stop)

	saw := pollUntil127(200, func() bool {
		out := leg.stdout.String()
		return strings.Contains(out, ballUpClaim) || strings.Contains(out, ballAbsentClaim)
	})
	if !saw {
		t.Fatalf("AC#2 LIVE RED: no ball posture sentence.\n%s", leg.console())
	}
	if strings.Contains(leg.stdout.String(), ballAbsentClaim) {
		t.Skipf("SKIP-LOUD: no ball window on this host (%q)", ballAbsentClaim)
	}
	out := leg.stdout.String()
	if !strings.Contains(out, "hotkeys from config") {
		t.Fatalf("AC#2 LIVE RED: a boot with the planted non-default summon did not name config as the source (the construction half of this ruler cannot be trusted for the rebind half):\n%s", out)
	}
	if !strings.Contains(out, "summon=Ctrl+Alt+Z") {
		t.Fatalf("AC#2 LIVE RED: the boot summary does not name the planted non-default summon Ctrl+Alt+Z:\n%s", out)
	}

	// THE EDIT: same file, new summon.
	newBody := "[hotkey]\nsummon = \"" + newSummon258 + "\"\nmute = \"Ctrl+Alt+M\"\ncancel = \"Esc\"\npanel = \"Ctrl+Alt+P\"\n"
	if err := os.WriteFile(filepath.Join(dataDir, configFileName), []byte(newBody), 0o600); err != nil {
		t.Fatalf("edit config: %v", err)
	}

	// The rebind books "ball: hotkeys rebound after config change" in the sink
	// (hotkey_reload.go's slog record); the summary line is printed once at
	// boot and is not rewritten. Read the raw jsonl bytes through the fixed
	// matcher, with a bounded poll - a sink that never produced the line ends
	// in a named failure carrying the tail, which is the honest "cannot read"
	// shape (readSinkTail258 says so in the tail itself when no jsonl exists).
	rebound := false
	var sinkTail string
	for i := 0; i < 240 && !rebound; i++ {
		sinkTail = readSinkTail258(t, logSinkDir(dataDir))
		rebound = sinkRebindLineHas258(sinkTail, newSummon258)
		if !rebound {
			time.Sleep(50 * time.Millisecond)
		}
	}
	if !rebound {
		t.Fatalf("AC#2 LIVE RED: 12s after the [hotkey] edit the process has not rebound the live keys (no rebind record naming summon=%s in the sink, in either the JSON colon shape \"summon\":\"%s\" or the console shape summon=%s).\nsink tail:\n%s", newSummon258, newSummon258, newSummon258, sinkTail)
	}
	t.Logf("AC#2 LIVE: the same process rebound summon to %s after the file edit", newSummon258)
}

// TestLive258ResidentOccupiedCombinationNamesTheNewValue is AC#3's live form:
// an occupied combination must be said in Problems() with the NEW value named.
//
// THE RULER WAS FIXED IN 258-R2 (258-v1 §5/§7-R7): it used to BET that some
// third-party program happened to hold Ctrl+Alt+U, with no fallback - on a
// machine where the key was free the child registered it, no occupied line
// existed, and the case went red on a working product. The premise is now
// BUILT: this process registers the combination itself (squatHotkeyForRuler258)
// before the child boots. Either this process holds it or a third party does
// (RegisterHotKey answered ERROR_HOTKEY_ALREADY_REGISTERED) - both readings
// satisfy the premise, and the child's panel slot cannot bind. Anything else
// (a RegisterHotKey failure that is neither) names the ruler red instead of
// judging blind. The only honest skip left is the no-ball-window host, checked
// the same way the rest of the winlive family checks it.
func TestLive258ResidentOccupiedCombinationNamesTheNewValue(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	const occupied258 = "Ctrl+Alt+U"
	body := "[hotkey]\nsummon = \"Ctrl+Alt+Q\"\nmute = \"Ctrl+Alt+M\"\ncancel = \"Esc\"\npanel = \"" + occupied258 + "\"\n"
	if err := os.WriteFile(filepath.Join(dataDir, configFileName), []byte(body), 0o600); err != nil {
		t.Fatalf("plant config: %v", err)
	}

	// Build the premise BEFORE the child exists: whoever ends up holding the
	// combination, the child's first registration attempt must be refused.
	owner := squatHotkeyForRuler258(t, occupied258)
	t.Logf("AC#3 LIVE premise: %s is %s", occupied258, owner)

	leg := bootResidentLeg(t, exe, dataDir)
	t.Cleanup(leg.stop)

	saw := pollUntil127(200, func() bool {
		out := leg.stdout.String()
		return strings.Contains(out, ballUpClaim) || strings.Contains(out, ballAbsentClaim)
	})
	if !saw {
		t.Fatalf("AC#3 LIVE RED: no ball posture sentence.\n%s", leg.console())
	}
	out := leg.stdout.String()
	if strings.Contains(out, ballAbsentClaim) {
		t.Skipf("SKIP-LOUD: no ball window on this host (%q)", ballAbsentClaim)
	}
	if !strings.Contains(out, "hotkeys from config") {
		t.Fatalf("AC#3 LIVE RED: the verdict does not name config as the source:\n%s", out)
	}
	if !strings.Contains(out, "occupied by another program") || !strings.Contains(out, "panel = \""+occupied258+"\"") {
		t.Fatalf("AC#3 RED: the premise was built (%s), yet the occupied panel combination produced no honest Problems() line naming %s:\n%s", owner, occupied258, out)
	}
	// With the premise built, the arithmetic of the boot verdict is pinned,
	// not logged: summon+mute can be the only live slots (cancel is standby by
	// ticket 245, panel was refused). A 3/4 here means the child registered
	// over the holder it should not have been able to bind past.
	if !strings.Contains(out, "hotkeys live 2/4") {
		t.Errorf("AC#3 RED: with %s occupied (%s) the boot verdict does not read hotkeys live 2/4:\n%s", occupied258, owner, out)
	}
	if !strings.Contains(out, "panel="+occupied258) {
		t.Errorf("AC#3 RED: the per-slot summary does not name the occupied NEW binding %s (the summary must say what Win32 made of the NEW value, not keep the old one):\n%s", occupied258, out)
	}
}

// liveSquatID258 is this process's registration id for the built premise. It
// only ever exists inside the test binary, with hWnd 0 (the hotkey is held for
// the locked thread's queue and released in cleanup); nothing pumps it.
const liveSquatID258 = 0x258

var (
	pSquatRegisterHotKey258   = windows.NewLazySystemDLL("user32.dll").NewProc("RegisterHotKey")
	pSquatUnregisterHotKey258 = windows.NewLazySystemDLL("user32.dll").NewProc("UnregisterHotKey")
)

// squatHotkeyForRuler258 turns binding into a BUILT occupancy premise for the
// winlive occupied ruler: it registers the desktop-wide combination in THIS
// process on a locked OS thread (RegisterHotKey is desktop-wide and
// per-process-blind, so the shipped child sees exactly what a third-party
// holder produces). Answers:
//   - registration succeeded: this test process is the holder, released in
//     t.Cleanup on the same thread that registered;
//   - ERROR_HOTKEY_ALREADY_REGISTERED: a third party already holds it, the
//     premise holds with that owner and nothing needs releasing;
//   - any other answer: named t.Fatalf - the ruler will not judge an
//     occupancy sentence on a machine where it could not build (or find) a
//     holder.
//
// It returns a one-line description of who ends up holding the key, for the
// failure messages and the premise log.
func squatHotkeyForRuler258(t *testing.T, binding string) string {
	t.Helper()
	acc, err := ball.ParseAccelerator(binding)
	if err != nil {
		t.Fatalf("258 LIVE RULER RED: parse %s for the built premise: %v", binding, err)
	}
	// The registration and its release must land on one thread: pin this
	// goroutine, register, and only unlock after UnregisterHotKey answered.
	runtime.LockOSThread()
	r, _, callErr := pSquatRegisterHotKey258.Call(0, uintptr(liveSquatID258), uintptr(acc.Mods), uintptr(acc.VK))
	if r == 0 {
		errno, ok := callErr.(syscall.Errno)
		if !ok || errno == 0 {
			errno = syscall.EINVAL
		}
		if errno == windows.ERROR_HOTKEY_ALREADY_REGISTERED {
			runtime.UnlockOSThread()
			return "held by a third-party program on this machine (RegisterHotKey answered ERROR_HOTKEY_ALREADY_REGISTERED before the child was even launched)"
		}
		runtime.UnlockOSThread()
		t.Fatalf("258 LIVE RULER RED: cannot build the occupancy premise for %s: RegisterHotKey answered %v. The pre-r2 ruler bet on the machine and called the bet a reading; this one refuses to judge without its premise.", binding, errno)
	}
	t.Cleanup(func() {
		pSquatUnregisterHotKey258.Call(0, uintptr(liveSquatID258))
		runtime.UnlockOSThread()
	})
	return "held by THIS test process (squatted on a locked thread before the child booted, released in cleanup)"
}

// readSinkTail258 concatenates the raw bytes of every jsonl in the sink dir
// (the rebind line is a slog record whose msg carries the phrase; the raw read
// avoids re-parsing the sink schema, which is the 127 family's subject).
func readSinkTail258(t *testing.T, dir string) string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(dir, "wisp-*.jsonl"))
	if err != nil || len(matches) == 0 {
		return "(no wisp-*.jsonl in " + dir + " yet)"
	}
	var sb strings.Builder
	for _, name := range matches {
		b, err := os.ReadFile(name)
		if err != nil {
			continue
		}
		sb.Write(b)
	}
	return tail258(sb.String(), 6000)
}

// tail258 returns the last n bytes of s for failure messages.
func tail258(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}

// min64 keeps the two-arg form local (the file builds on go1.21+, which has
// min, but the repo's pins stay explicit).
func min64(a, b int) int {
	if a < b {
		return a
	}
	return b
}
