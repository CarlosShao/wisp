package config

// Ticket 39 AC#1's tier-events teeth at the layer that owns the write path
// (2026-11-18). The route-level teeth live in cmd/wisp (tier_events_39_test.go
// through the real panel settings door); these prove the guarded section-write
// path itself honors the registered D36 tier of the key it writes:
//
//	restart (app.language)  -> file written, memory keeps the boot value, the
//	                           next reload books Report.Restart and fires
//	                           OnRestartPending exactly once (the 需要重启
//	                           surface an external edit lands on).
//	reload  (voice.aec.enabled) -> memory applies AND Manager.OnReload fires
//	                           once with the owning section, values already
//	                           applied - the subscriber reads the new value
//	                           off Config() inside the event.
//	hot     (ball.style)    -> memory applies immediately, no events at all.
//
// THE TIER COMES FROM THE REGISTRY IN EVERY CASE: each subtest asks
// TierOf(key) and branches on the registry's own answer; a key re-tiered in
// tiers.go fails LOUDLY naming the mismatch instead of silently running the
// wrong branch (that is the ⛔ no-hardcoded-tiers rule enforced, not dodged).

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const tierBase39 = `schema_version = 2

[app]
theme = "dark"
language = "zh-CN"

[ball]
size = 56

[voice]
enabled = true

[voice.aec]
enabled = false
`

func tierManager39(t *testing.T) (*Manager, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(tierBase39), 0o600); err != nil {
		t.Fatal(err)
	}
	m, err := NewManager(path, nil)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	return m, path
}

// Test39TierRestartWriteKeepsTheOldValueAndTheNextReloadSaysSo: the restart
// tier's two observables through a PROGRAMMATIC write (ticket 223's external
// edit shape, now reached via the guarded write path).
func Test39RestartTierWriteLandsOnDiskButNotInMemory(t *testing.T) {
	m, path := tierManager39(t)
	restarts := [][]string{}
	reloads := [][]string{}
	m.OnRestartPending = func(sections []string) { restarts = append(restarts, sections) }
	m.OnReload = func(sections []string) { reloads = append(reloads, sections) }

	// The tier judgment is the registry's own row, never this test's word.
	tier, ok := TierOf("app.language")
	if !ok || Tier(tier) != TierRestart {
		t.Fatalf("TierOf(app.language) = (%q,%v), want the registry's restart row; re-tier the key in tiers.go and move this test with it", tier, ok)
	}

	res, err := m.SetAppLanguage("en-US")
	if err != nil {
		t.Fatalf("SetAppLanguage: %v", err)
	}
	if !reflect.DeepEqual(res, []string{"app.language"}) {
		t.Fatalf("written = %v, want [app.language]", res)
	}

	// Observable 1: the running config keeps the boot value (不许声称已生效).
	if got := m.Config().App.Language; got != "zh-CN" {
		t.Fatalf("Config().App.Language = %q, want zh-CN: a restart-tier write hot-applied to memory (the restart-当-hot lie)", got)
	}
	// Observable 2: the file carries the new value for the next start.
	disk, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(disk), "en-US") {
		t.Fatalf("the file does not carry the written language:\n%s", disk)
	}
	// No events fired by the write itself: the notice is the next reload's.
	if len(restarts) != 0 || len(reloads) != 0 {
		t.Fatalf("the write itself fired restarts=%v reloads=%v; the D36 notice is the reload's", restarts, reloads)
	}

	// The next reload books the change in Report.Restart and says so once.
	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("CheckAndReload: %v", err)
	}
	if rep == nil || !reflect.DeepEqual(rep.Restart, []string{"app"}) || len(rep.Hot) != 0 {
		t.Fatalf("report = %+v, want app restart only (the restart write must land in the restart bucket, not hot)", rep)
	}
	if len(restarts) != 1 || !reflect.DeepEqual(restarts[0], []string{"app"}) {
		t.Fatalf("OnRestartPending = %v, want exactly one [app] notice (静默不生效 is the shape this tooth forbids)", restarts)
	}
	if len(reloads) != 0 {
		t.Fatalf("a restart-tier change fired OnReload %v: the reload event would claim a subsystem rebuilt", reloads)
	}
	// And the old value STILL stands after the reload (planApp keeps it).
	if got := m.Config().App.Language; got != "zh-CN" {
		t.Fatalf("after the reload Config().App.Language = %q, want zh-CN (restart tier keeps the old value until process start)", got)
	}
}

// Test39ReloadTierWriteAppliesAndSignalsTheOwningSubsystem: apply + event, the
// reload tier's own pair of observables, reached through a programmatic write.
func Test39ReloadTierWriteFiresOnReloadWithValuesApplied(t *testing.T) {
	m, _ := tierManager39(t)
	restarts := [][]string{}
	var reloadReads []bool
	m.OnRestartPending = func(sections []string) { restarts = append(restarts, sections) }

	tier, ok := TierOf("voice.aec.enabled")
	if !ok || Tier(tier) != TierReload {
		t.Fatalf("TierOf(voice.aec.enabled) = (%q,%v), want the registry's reload row; re-tier the key in tiers.go and move this test with it", tier, ok)
	}

	// The subscriber: the reload event's own contract says values are already
	// applied, so the rebuild reads the NEW value off Config() inside the hook.
	m.OnReload = func(sections []string) {
		reloadReads = append(reloadReads, m.Config().Voice.AEC.Enabled)
		_ = sections
	}

	res, err := m.SetVoiceAecEnabled(true)
	if err != nil {
		t.Fatalf("SetVoiceAecEnabled: %v", err)
	}
	if !reflect.DeepEqual(res, []string{"voice.aec.enabled"}) {
		t.Fatalf("written = %v, want [voice.aec.enabled]", res)
	}
	if len(reloadReads) != 1 || !reloadReads[0] {
		t.Fatalf("OnReload fired %d time(s), reads %v - want exactly one event and the subscriber reading the NEW value (applied before the event)", len(reloadReads), reloadReads)
	}
	if len(restarts) != 0 {
		t.Fatalf("a reload-tier write fired OnRestartPending %v", restarts)
	}
	// The live reader answers the new value without any restart.
	if got := m.Config().Voice.AEC.Enabled; !got {
		t.Fatalf("Config().Voice.AEC.Enabled = false after a reload-tier write: the apply half is dead")
	}

	// The no-op save fires nothing: the file already carries the value.
	reloadReads = nil
	if _, err := m.SetVoiceAecEnabled(true); err != nil {
		t.Fatalf("no-op SetVoiceAecEnabled: %v", err)
	}
	if len(reloadReads) != 0 {
		t.Fatalf("a no-op reload write fired OnReload %d time(s): nothing changed, so there is nothing to reload", len(reloadReads))
	}
}

// Test39HotTierWriteAppliesImmediatelyWithNoEvents: the hot tier's observable
// (改了 -> 不重启 -> 行为变了) through the same guarded write path.
func Test39HotTierWriteAppliesImmediatelyWithoutEvents(t *testing.T) {
	m, _ := tierManager39(t)
	restarts := [][]string{}
	reloads := [][]string{}
	m.OnRestartPending = func(sections []string) { restarts = append(restarts, sections) }
	m.OnReload = func(sections []string) { reloads = append(reloads, sections) }

	// [ball] is a WHOLE-section hot row: TierOf answers for the section.
	tier, ok := TierOf("ball")
	if !ok || Tier(tier) != TierHot {
		t.Fatalf("TierOf(ball) = (%q,%v), want the registry's hot row; re-tier the section in tiers.go and move this test with it", tier, ok)
	}

	if _, err := m.SetBallStyle("jarvis"); err != nil {
		t.Fatalf("SetBallStyle: %v", err)
	}
	if got := m.Config().Ball.Style; got != "jarvis" {
		t.Fatalf("Config().Ball.Style = %q after a hot-tier write, want jarvis (改了 -> 不重启 -> 值已换)", got)
	}
	if len(restarts) != 0 || len(reloads) != 0 {
		t.Fatalf("a hot-tier write fired restarts=%v reloads=%v: a hot key carries neither notice", restarts, reloads)
	}
}
