//go:build windows

package main

// Ticket 258 verification leg 258-v1: an independently written positive
// control for AC#1 (the ticket sentence: "种一发'改 [hotkey] 里 summon 的组合键
// ⇒ 下一次建球注册的就是新值'"). It is a rewrite in the same shape as
// Test258BridgeRebindsLiveKeysFromConfigEdit (the dead leg's in-process case)
// and MUST NOT be cited as that case's reading - the probe owns its own
// fixture, its own expectations (planted Z then W), and its own poll loop, so
// its readings are this leg's own, not a transcription of the dead leg's.
//
// The case runs the FULL production wiring: startResidentBall (assembly-root
// shaped, live src closure reading a real config.toml on disk), a real ball
// window, an on-disk [hotkey] edit, a bridge poll hop driven through the seam,
// and reads what Win32 holds after the rebind. A second case pins AC#1's
// fallback sentence end to end through the chain function itself.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
)

// v1probeWriteHotkey writes a [hotkey]-only config.toml body to dir. The
// [hotkey]-only body is deliberate: it isolates the section under adjudication
// from any other section's behavior.
func v1probeWriteHotkey(t *testing.T, dir string, summon string) {
	t.Helper()
	path := filepath.Join(dir, configFileName)
	body := "[hotkey]\nsummon = \"" + summon + "\"\nmute = \"Ctrl+Alt+M\"\ncancel = \"Esc\"\npanel = \"Ctrl+Alt+P\"\n"
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write [hotkey] config: %v", err)
	}
}

// Test258V1ProbeSummonEditRebindsLiveBall is AC#1's positive control, written
// by the verification leg: boot on Ctrl+Alt+Z, edit to Ctrl+Alt+W, poll hop,
// Win32 must answer 'W'.
func Test258V1ProbeSummonEditRebindsLiveBall(t *testing.T) {
	dir := t.TempDir()
	v1probeWriteHotkey(t, dir, "Ctrl+Alt+Z")
	src := func() ball.HotkeyConfig {
		c, _, err := config.LoadFile(filepath.Join(dir, configFileName), nil)
		if err != nil || c == nil {
			return ball.HotkeyConfig{}
		}
		return ball.HotkeyConfig{Summon: c.Hotkey.Summon, Mute: c.Hotkey.Mute, Cancel: c.Hotkey.Cancel, Panel: c.Hotkey.Panel}
	}

	rb := startResidentBall(observe.NewRegistry(), nil, src, src)
	if rb == nil || rb.b == nil {
		t.Fatalf("258-v1 PROBE: no ball window on this host - the positive control cannot run headless (verdict line: %q)", rb.verdictStatusLineFor258())
	}
	defer rb.stop()
	if !rb.hotkeyBridgeArmed {
		t.Fatalf("258-v1 PROBE: bridge did not arm with a live src")
	}

	// BOOT: Win32 must hold 'Z' (the planted value, not any compiled literal).
	if got := rb.b.RegisteredHotkeys()[hkSummon258].VK; got != 'Z' {
		t.Fatalf("258-v1 PROBE BOOT RED: summon VK = 0x%X, want 'Z' (planted Ctrl+Alt+Z)", got)
	}

	// THE EDIT: summon -> Ctrl+Alt+7 (digit 0x37; deliberately rare - the first
	// draft of this probe used Ctrl+Alt+W and the rebind was refused with
	// ERROR_HOTKEY_ALREADY_REGISTERED: W is held by another program on this
	// desktop, which is itself the AC#3 occupied semantics working), on disk;
	// the src re-reads it per tick.
	const newVK258 uint32 = 0x37
	v1probeWriteHotkey(t, dir, "Ctrl+Alt+7")

	// Bounded retries through the seam's synchronous hop, never wall-clock
	// arithmetic (D42#9): the src read is synchronous, so a handful of hops
	// must be enough if the wiring is real.
	arrived := false
	for i := 0; i < 50 && !arrived; i++ {
		rb.hotkeyBridgeCheck258()
		arrived = rb.b.RegisteredHotkeys()[hkSummon258].VK == newVK258
		if !arrived {
			time.Sleep(20 * time.Millisecond)
		}
	}
	if !arrived {
		live := rb.b.RegisteredHotkeys()
		t.Fatalf("258-v1 PROBE RED: the [hotkey] edit never reached Win32: summon VK = 0x%X (%+v), want 0x37", live[hkSummon258].VK, live)
	}

	// The rebind's report must name the NEW binding; the verdict keeps the
	// boot provenance and the live-count shape (the boot sentence is not
	// rewritten by a rebind).
	if got := rb.b.HotkeyReport().Live()[hkSummon258].VK; got != newVK258 {
		t.Errorf("258-v1 PROBE: post-rebind report summon VK = 0x%X, want 0x37", got)
	}
	if rb.b.ConfiguredHotkeys().Summon != "Ctrl+Alt+7" {
		t.Errorf("258-v1 PROBE: ConfiguredHotkeys().Summon = %q, want Ctrl+Alt+7", rb.b.ConfiguredHotkeys().Summon)
	}
	if !strings.Contains(rb.verdict, "hotkeys from config") {
		t.Errorf("258-v1 PROBE: boot verdict lost the config provenance: %q", rb.verdict)
	}
	if !strings.Contains(rb.verdict, "hotkeys live 3/4") {
		t.Errorf("258-v1 PROBE: boot verdict lost the live-count shape: %q", rb.verdict)
	}
}

// Test258V1ProbeConstructionMissingFileNamesTheFallback drives the chain
// function the way the assembly root's closure does (LoadFile over a dir with
// no config.toml) and requires the fallback to be NAMED, per AC#1's
// "说得出这句话" clause.
func Test258V1ProbeConstructionMissingFileNamesTheFallback(t *testing.T) {
	// (1) missing file, through the chain:
	dir := t.TempDir()
	hotCfg := func() ball.HotkeyConfig {
		c, _, err := config.LoadFile(filepath.Join(dir, configFileName), nil)
		if err != nil || c == nil {
			return ball.HotkeyConfig{}
		}
		return ball.HotkeyConfig{Summon: c.Hotkey.Summon, Mute: c.Hotkey.Mute, Cancel: c.Hotkey.Cancel, Panel: c.Hotkey.Panel}
	}
	got, provenance := residentBallHotkeyChain258(hotCfg)
	if provenance != hotkeyProvenanceDefaults {
		t.Fatalf("258-v1 PROBE: missing-file provenance = %q, want %q (AC#1 forbids the silent swap)", provenance, hotkeyProvenanceDefaults)
	}
	if got != ball.DefaultHotkeys() {
		t.Fatalf("258-v1 PROBE: missing-file chain = %+v, want DefaultHotkeys() %+v", got, ball.DefaultHotkeys())
	}

	// (2) planted config: the chain must answer the planted value with
	// provenance "config" - the same drive the boot leg does.
	v1probeWriteHotkey(t, dir, "Ctrl+Alt+Z")
	got, provenance = residentBallHotkeyChain258(hotCfg)
	if provenance != hotkeyProvenanceConfig {
		t.Fatalf("258-v1 PROBE: planted provenance = %q, want %q", provenance, hotkeyProvenanceConfig)
	}
	if got.Summon != "Ctrl+Alt+Z" {
		t.Fatalf("258-v1 PROBE: planted chain summon = %q, want Ctrl+Alt+Z", got.Summon)
	}
}
