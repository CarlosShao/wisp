package main

// Ticket 39 AC#1's tier-events route teeth (2026-11-18): one representative
// key per D36 tier, written through the REAL panel settings door
// (ConfigWriteHandler -> newConfigStore -> config's guarded, tier-aware write
// path), the tier judgment asked of internal/config's TierRegistry every time:
//
//	Test39TierRestartWriteDemandsRestartOnThePage   (app.language, restart)
//	    the save is accepted, the receipt carries 需要重启 verbatim and never
//	    claims an immediate effect, the reload event says tier=restart, the
//	    running config keeps the boot value, the file carries the new one, and
//	    the next reload books Report.Restart + one OnRestartPending notice.
//	Test39TierReloadWriteFiresTheReloadEventAndTheSubscriberReadsIt
//	    (voice.aec.enabled, reload) the accepted save carries the
//	    wisp-config-reload: marker event, Manager.OnReload fires exactly once
//	    with the owning section, and the subscriber (the event's own contract:
//	    values already applied) reads the NEW value off Config() inside the
//	    event.
//
// The hot tier's representative (ball.style with the live applier) lives in
// tier_events_39_windows_test.go - the applier seam is the ball host's.
//
// THE TIER COMES FROM THE REGISTRY IN EVERY TEST: each one asks
// config.TierOf(key) and branches on the registry's own answer; a key
// re-tiered in tiers.go fails loudly naming the mismatch instead of silently
// running the wrong branch. The production leg stamps the receipt tier through
// the same table (registryEffectiveTier39, cmd/wisp/config_tier_39.go) - a
// hand-stamped tier there is what the mutation teeth catch.

import (
	"context"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
)

const tierEventsBase39 = `schema_version = 2

[app]
theme = "dark"
language = "zh-CN"

[voice]
enabled = true

[voice.aec]
enabled = false
`

func tierEventsRig39(t *testing.T) (*panel.ConfigWriteHandler, *configStore, *config.Manager, string) {
	t.Helper()
	dir := t.TempDir()
	path := dir + "/config.toml"
	write39(t, path, tierEventsBase39)
	mgr, err := config.NewManager(path, nil)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	leg := newConfigStore(mgr, nil, nil, "test-39-tier")
	h := &panel.ConfigWriteHandler{Store: leg, Actor: "test-39-tier"}
	return h, leg, mgr, path
}

// tierFingerprint39 is the page half of the drift handshake: read the config,
// take the fingerprint the read face answered.
func tierFingerprint39(t *testing.T, leg *configStore) string {
	t.Helper()
	view, err := leg.ReadSettings(context.Background())
	if err != nil {
		t.Fatalf("ReadSettings: %v", err)
	}
	if view.BaseSHA == "" {
		t.Fatal("the read face answered no fingerprint; the page has nothing to echo")
	}
	return view.BaseSHA
}

func Test39TierRestartWriteDemandsRestartOnThePage(t *testing.T) {
	h, leg, mgr, path := tierEventsRig39(t)

	// The tier judgment is the registry's own row (tiers.go's per-key restart
	// entry for app.language), never this test's word.
	if tier, ok := config.TierOf("app.language"); !ok || tier != string(config.TierRestart) {
		t.Fatalf("TierOf(app.language) = (%q,%v), want the registry's restart row; re-tier the key in tiers.go and move this test with it", tier, ok)
	}

	restarts := [][]string{}
	reloads := [][]string{}
	mgr.OnRestartPending = func(sections []string) { restarts = append(restarts, sections) }
	mgr.OnReload = func(sections []string) { reloads = append(reloads, sections) }
	base := tierFingerprint39(t, leg)

	reply, err := h.HandleConfigRequest(context.Background(), panel.ComposerRequest{
		Method: panel.MethodConfigSet, RequestID: "rid-restart-39", Source: panel.ComposerRequestSource,
		ConfigField: panel.FieldAppLanguage, ConfigValue: "en-US", ConfigBaseSHA: base,
	}, "")
	if err != nil {
		t.Fatalf("config.set app_language was refused: %v", err)
	}

	// The receipt: 需要重启 VERBATIM, the page-visible restart-required line,
	// and no claim that the save is already in force.
	if !strings.Contains(reply, "需要重启") {
		t.Fatalf("the restart receipt does not carry 需要重启 verbatim:\n%s", reply)
	}
	if strings.Contains(reply, "立即生效") || strings.Contains(reply, "保存即生效") {
		t.Fatalf("a restart-tier receipt claims the save is in force:\n%s", reply)
	}
	// The event: the write's own tier word is the registry's.
	ev := event39(t, reply)
	if ev.Tier != "restart" {
		t.Fatalf("the reload event's tier = %q, want restart (the registry's own word)", ev.Tier)
	}
	if len(ev.Written) != 1 || ev.Written[0] != "app.language" {
		t.Fatalf("the event's written = %v, want [app.language]", ev.Written)
	}
	if len(ev.Sections) != 1 || ev.Sections[0] != "app" {
		t.Fatalf("the event's sections = %v, want [app]", ev.Sections)
	}

	// 不许声称已生效: the running config keeps the boot value...
	if got := mgr.Config().App.Language; got != "zh-CN" {
		t.Fatalf("Config().App.Language = %q after a restart-tier save, want zh-CN (the running process must not claim the new value)", got)
	}
	// ...the page's next read says the same (the re-read the event triggers).
	if view, err := leg.ReadSettings(context.Background()); err != nil || !view.Readable {
		t.Fatalf("re-read failed: %v", err)
	} else {
		for _, row := range view.Sections {
			if row.Key == "app.language" {
				if row.Value != "zh-CN" || row.Tier != "restart" {
					t.Fatalf("the section tree answers app.language = (%q, tier %q), want (zh-CN, restart)", row.Value, row.Tier)
				}
			}
		}
	}
	// ...and the file carries the new value for the next start.
	if disk, err := os.ReadFile(path); err != nil {
		t.Fatal(err)
	} else if !strings.Contains(string(disk), "en-US") {
		t.Fatalf("the file does not carry the saved language:\n%s", disk)
	}

	// No events fired by the save itself; the notice is the next reload's.
	if len(restarts) != 0 || len(reloads) != 0 {
		t.Fatalf("the save itself fired restarts=%v reloads=%v", restarts, reloads)
	}

	// The next reload books the change in Report.Restart and says so once.
	rep, err := mgr.CheckAndReload()
	if err != nil {
		t.Fatalf("CheckAndReload: %v", err)
	}
	if rep == nil || !reflect.DeepEqual(rep.Restart, []string{"app"}) || len(rep.Hot) != 0 {
		t.Fatalf("report = %+v, want app restart only", rep)
	}
	if len(restarts) != 1 || !reflect.DeepEqual(restarts[0], []string{"app"}) {
		t.Fatalf("OnRestartPending = %v, want exactly one [app] notice (静默不生效 is the shape this tooth forbids)", restarts)
	}
	if len(reloads) != 0 {
		t.Fatalf("a restart-tier change fired OnReload %v: the reload event would claim a subsystem rebuilt", reloads)
	}
}

func Test39TierReloadWriteFiresTheReloadEventAndTheSubscriberReadsIt(t *testing.T) {
	h, leg, mgr, _ := tierEventsRig39(t)

	// The tier judgment is the registry's own row (tiers.go's per-key reload
	// entry for voice.aec.enabled), never this test's word.
	if tier, ok := config.TierOf("voice.aec.enabled"); !ok || tier != string(config.TierReload) {
		t.Fatalf("TierOf(voice.aec.enabled) = (%q,%v), want the registry's reload row; re-tier the key in tiers.go and move this test with it", tier, ok)
	}

	// The subscriber: the reload event's own contract says values are already
	// applied, so the rebuild reads the NEW value off Config() inside the hook.
	var reloadReads []bool
	restarts := [][]string{}
	mgr.OnReload = func(sections []string) {
		reloadReads = append(reloadReads, mgr.Config().Voice.AEC.Enabled)
	}
	mgr.OnRestartPending = func(sections []string) { restarts = append(restarts, sections) }
	base := tierFingerprint39(t, leg)

	reply, err := h.HandleConfigRequest(context.Background(), panel.ComposerRequest{
		Method: panel.MethodConfigSet, RequestID: "rid-reload-39", Source: panel.ComposerRequestSource,
		ConfigField: panel.FieldVoiceAecEnabled, ConfigValue: "true", ConfigBaseSHA: base,
	}, "")
	if err != nil {
		t.Fatalf("config.set voice_aec_enabled was refused: %v", err)
	}

	// The marker event: an accepted save carries the wisp-config-reload: line
	// with the key it wrote and the section that owns it.
	ev := event39(t, reply)
	if len(ev.Written) != 1 || ev.Written[0] != "voice.aec.enabled" {
		t.Fatalf("the event's written = %v, want [voice.aec.enabled]", ev.Written)
	}
	if len(ev.Sections) != 1 || ev.Sections[0] != "voice" {
		t.Fatalf("the event's sections = %v, want [voice]", ev.Sections)
	}

	// The Manager-level reload event: exactly one, with the new value already
	// readable by the subscriber inside it (重建后读新值).
	if len(reloadReads) != 1 || !reloadReads[0] {
		t.Fatalf("OnReload fired %d time(s), subscriber reads %v - want exactly one event carrying the NEW value", len(reloadReads), reloadReads)
	}
	if len(restarts) != 0 {
		t.Fatalf("a reload-tier save fired OnRestartPending %v: a reload key carries no restart notice", restarts)
	}

	// The live reader answers the new value without any restart, and the file
	// agrees (the next poll is a no-op, so the event cannot double-fire).
	if !mgr.Config().Voice.AEC.Enabled {
		t.Fatal("Config().Voice.AEC.Enabled = false after a reload-tier save: the apply half is dead")
	}
	rep, err := mgr.CheckAndReload()
	if err != nil {
		t.Fatalf("CheckAndReload: %v", err)
	}
	if rep != nil {
		t.Fatalf("CheckAndReload after the save reported %+v, want nil (the write is its own reload; the next poll must be a no-op)", rep)
	}
	if len(reloadReads) != 1 {
		t.Fatalf("OnReload fired %d time(s) after the poll: a reload event doubled", len(reloadReads))
	}
}
