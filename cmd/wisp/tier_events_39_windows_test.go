//go:build windows

package main

// Ticket 39 AC#1's tier-events route teeth, HOT tier (2026-11-18): ball.style
// is the representative hot key whose receipt may say 立即生效 - the applier
// seam is installed (the real product function ball.SetStyle, the same rig
// ball_style_65_windows_test.go uses), the write goes through the REAL panel
// settings door, and the subscriber is the ball itself:
//
//   - the receipt carries the now tier sentence and carries NO restart
//     wording (无要重启字样 - a hot key must never say 需要重启);
//   - the reload event rides the reply with tier=now and the key it wrote;
//   - the subscriber (ball.ActiveStyle) and the live config reader both
//     answer the NEW value without any restart.
//
// The tier judgment is the registry's own row (tiers.go's whole-section hot
// entry for [ball]), never this test's word.

import (
	"context"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/ball"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
)

const tierHotBase39 = `schema_version = 2

[ball]
size = 56
`

func Test39TierHotWriteAppliesNowAndSaysNoRestart(t *testing.T) {
	t.Cleanup(func() {
		ball.SetStyle(string(ball.StyleGlass)) // the package default; tests must not leave the family flipped
		ballStyleApplier65 = nil
	})
	dir := t.TempDir()
	path := dir + "/config.toml"
	write39(t, path, tierHotBase39)
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	leg := newConfigStore(mgr, nil, nil, "test-39-tier-hot")
	h := &panel.ConfigWriteHandler{Store: leg, Actor: "test-39-tier-hot"}

	// The tier judgment is the registry's own row: [ball] is a WHOLE-section
	// hot entry in tiers.go.
	if tier, ok := config.TierOf("ball"); !ok || tier != string(config.TierHot) {
		t.Fatalf("TierOf(ball) = (%q,%v), want the registry's hot row; re-tier the section in tiers.go and move this test with it", tier, ok)
	}

	// The live applier (the ball host's own seam, the real product function):
	// this is what makes "now" a measured fact instead of a claim.
	ballStyleApplier65 = ball.SetStyle
	base := tierFingerprint39(t, leg)

	reply, err := h.HandleConfigRequest(context.Background(), panel.ComposerRequest{
		Method: panel.MethodConfigSet, RequestID: "rid-hot-39", Source: panel.ComposerRequestSource,
		ConfigField: panel.FieldBallStyle, ConfigValue: "jarvis", ConfigBaseSHA: base,
	}, "")
	if err != nil {
		t.Fatalf("config.set ball_style was refused: %v", err)
	}

	// 无要重启字样: the hot receipt says the save took effect now and carries
	// no restart wording anywhere.
	if !strings.Contains(reply, "立即生效") {
		t.Fatalf("a hot-tier save with a live applier does not say 立即生效:\n%s", reply)
	}
	if strings.Contains(reply, "重启") {
		t.Fatalf("a hot-tier receipt carries restart wording:\n%s", reply)
	}
	// The event: tier=now (the receipt's own word), the key it wrote.
	ev := event39(t, reply)
	if ev.Tier != string(panel.EffectiveNow) {
		t.Fatalf("the reload event's tier = %q, want %q", ev.Tier, panel.EffectiveNow)
	}
	if len(ev.Written) != 1 || ev.Written[0] != "ball.style" {
		t.Fatalf("the event's written = %v, want [ball.style]", ev.Written)
	}
	if len(ev.Sections) != 1 || ev.Sections[0] != "ball" {
		t.Fatalf("the event's sections = %v, want [ball]", ev.Sections)
	}

	// The subscribers: the ball's process-wide family moved, and the live
	// config reader answers the new value - no restart anywhere.
	if ball.ActiveStyle() != ball.StyleJarvis {
		t.Fatalf("ActiveStyle after the hot save = %v, want jarvis (the subscriber never got the new value)", ball.ActiveStyle())
	}
	if got := mgr.Config().Ball.Style; got != "jarvis" {
		t.Fatalf("Config().Ball.Style = %q after a hot-tier save, want jarvis", got)
	}

	// The control on the same rig: no applier = no live switch = the honest
	// restart word. Delete the measured override in panel_config_store.go and
	// this half reds.
	ballStyleApplier65 = nil
	reply2, err := h.HandleConfigRequest(context.Background(), panel.ComposerRequest{
		Method: panel.MethodConfigSet, RequestID: "rid-hot-39b", Source: panel.ComposerRequestSource,
		ConfigField: panel.FieldBallStyle, ConfigValue: "glass",
	}, "")
	if err != nil {
		t.Fatalf("the no-applier save was refused: %v", err)
	}
	if !strings.Contains(reply2, "重启") {
		t.Fatalf("a hot-tier save with NO live applier must answer the restart word honestly:\n%s", reply2)
	}
}
