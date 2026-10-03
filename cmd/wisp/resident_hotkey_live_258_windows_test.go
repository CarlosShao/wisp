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
	"strings"
	"testing"
	"time"
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
// UP on the default summon, then config.toml is edited underneath it; the poll
// (1s) must rebind the live keys inside the same process. Read through the
// process's own summary line, which the rebind updates via HotkeyReport.
func TestLive258ResidentRebindsWithoutRestart(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	const newSummon258 = "Ctrl+Alt+R"
	body := "[hotkey]\nsummon = \"Ctrl+Alt+Q\"\nmute = \"Ctrl+Alt+M\"\ncancel = \"Esc\"\npanel = \"Ctrl+Alt+P\"\n"
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
	if !strings.Contains(out, "summon=Ctrl+Alt+Q") {
		t.Fatalf("AC#2 LIVE RED: the boot summary does not name the planted summon Ctrl+Alt+Q:\n%s", out)
	}

	// THE EDIT: same file, new summon.
	newBody := "[hotkey]\nsummon = \"" + newSummon258 + "\"\nmute = \"Ctrl+Alt+M\"\ncancel = \"Esc\"\npanel = \"Ctrl+Alt+P\"\n"
	if err := os.WriteFile(filepath.Join(dataDir, configFileName), []byte(newBody), 0o600); err != nil {
		t.Fatalf("edit config: %v", err)
	}

	// The rebind books "ball: hotkeys rebound after config change" in the sink
	// and the NEXT summary reads it; the summary itself is printed once at boot,
	// so the durable reading is the log file the sink installed. Read the raw
	// jsonl bytes (the rebind line is a slog record whose msg carries the
	// phrase), with a bounded poll.
	rebound := false
	var sinkTail string
	for i := 0; i < 240 && !rebound; i++ {
		sinkTail = readSinkTail258(t, logSinkDir(dataDir))
		if strings.Contains(sinkTail, "hotkeys rebound after config change") &&
			strings.Contains(sinkTail, "summon="+newSummon258) {
			rebound = true
		}
		if !rebound {
			time.Sleep(50 * time.Millisecond)
		}
	}
	if !rebound {
		t.Fatalf("AC#2 LIVE RED: 12s after the [hotkey] edit the process has not rebound the live keys (no 'hotkeys rebound after config change' with summon=%s in the sink).\nsink tail:\n%s", newSummon258, sinkTail)
	}
	t.Logf("AC#2 LIVE: the same process rebound summon to %s after the file edit", newSummon258)
}

// TestLive258ResidentOccupiedCombinationNamesTheNewValue is AC#3's live form:
// an occupied combination must be said in Problems() with the NEW value named.
// The squat this case uses is the machine's own Ctrl+Alt+U risk
// (internal/ball's live suite measured it free on the dev machine; if some
// other program owns it, the case says so loudly rather than pretending).
func TestLive258ResidentOccupiedCombinationNamesTheNewValue(t *testing.T) {
	exe := buildWispForTest(t)
	dataDir := t.TempDir()
	const occupied258 = "Ctrl+Alt+U"
	body := "[hotkey]\nsummon = \"Ctrl+Alt+Q\"\nmute = \"Ctrl+Alt+M\"\ncancel = \"Esc\"\npanel = \"" + occupied258 + "\"\n"
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
		t.Fatalf("AC#3 RED: an occupied panel combination did not produce the honest Problems() line naming %s:\n%s", occupied258, out)
	}
	if strings.Contains(out, "hotkeys live 3/4") {
		t.Logf("AC#3 LIVE: live count reads 3/4 with panel taken (summon+mute live, cancel standby)")
	}
	if !strings.Contains(out, "panel="+occupied258) {
		t.Errorf("AC#3 RED: the per-slot summary does not name the occupied NEW binding %s (the summary must say what Win32 made of the NEW value, not keep the old one):\n%s", occupied258, out)
	}
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
