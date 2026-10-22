package config

// Ticket 223 AC#2 and AC#7, measured at the layer that owns the tiers.
//
// WHAT WAS ALREADY THERE. manager_test.go pins one change per tier in isolation:
// a [ball] size bump is hot-applied (TestManagerHotTierAppliesImmediately), an
// [app] language edit keeps the old value and fires OnRestartPending
// (TestManagerRestartTierKeepsOldValues), a [voice] model swap applies and fires
// OnReload (TestManagerReloadTierEmitsEvent). What none of them pins is the pair
// of claims AC#2 and AC#7 turn on:
//
//   - the restart tier's notice must be EXPLICIT (silence is the failure this
//     ticket exists for) and it must carry a WHY. The why this layer can back is
//     machine-checkable: TierRegistry says which key is restart-tier, and TierOf
//     hands that answer to a host. In production the hook is assigned at
//     cmd/wisp/config_reload.go:119 (startConfigReload, reached from run.go:861)
//     and the operator sentence is printed at cmd/wisp/config_reload.go:325-330 -
//     both outside this leg's write face, pinned since 票 232 by
//     TestTicket223RestartTierSaysItWillNotApply. This file proves the data that
//     sentence stands on instead of claiming the sentence itself.
//   - the two observables must be DIFFERENT objects, not one test wearing two
//     hats: "the value moved and a reader sees it" (hot) versus "the value did not
//     move and the host was told why" (restart). One edit that moves a hot key and
//     a restart key inside the SAME section is the case that proves they are
//     separable - TestTicket223R3OneEditGivesTwoDifferentObservables.
//
// WHY [llm] IS THE HOT CASE. Ticket 255 measured which hot sections this host may
// call "已立即生效": only [llm] has a reader that re-reads the LIVE Manager per
// call (cmd/wisp/panel_config_store.go:92 [cfg := s.mgr.Config()], surfaced at
// :98 [ChatModel: cfg.LLM.Roles.Chat.Model]). [ball] - the key manager_test.go
// uses - is registered no-reader in that roster, so calling a ball.size bump
// "行为变了" would repeat the claim 票 255 was filed to stop. The hot case below
// therefore edits llm.roles.chat.model and reads it back the way that reader does.
//
// THE THIRD TIER'S BOUNDARY, NAMED. The reload case pins the config layer's half
// (value applied + event emitted). The event's production consumer is NOT here:
// Manager.OnReload has zero assignments in cmd/wisp today (its only non-test
// assignment is cmd/balldebug - a side program), which is 票 233's cell. This file
// says that out loud rather than letting "the event fired" read as "the voice
// pipeline reloaded".

import (
	"reflect"
	"strings"
	"testing"
)

// r3Base223 is a valid config whose [llm] carries a catalog pair and a chat role
// (so a role edit is legal - validateCatalog names an unknown pair), and whose
// [app] carries the hot key (theme) next to one restart key (language), which is
// what the same-section case needs.
const r3Base223 = `
schema_version = 2

[app]
theme = "dark"
language = "zh-CN"

[llm.providers.acme]
protocol = "openai-chat"
base_url = "https://acme.example/v1"

[llm.providers.acme.models.m1]
billing = "pay-per-token"

[llm.providers.acme.models.m2]
billing = "pay-per-token"

[llm.roles.chat]
provider = "acme"
model = "m1"

[voice]
enabled = true

[voice.asr]
provider = "local-sherpa"
model = "parakeet-zh"
`

// r3RestartNotices installs the hook AC#7 is about and returns a reader for what
// it has been told. CheckAndReload calls it synchronously after commit, so a
// plain slice is enough here; no goroutine is involved in these cases.
func r3RestartNotices(m *Manager) func() [][]string {
	var seen [][]string
	m.OnRestartPending = func(sections []string) { seen = append(seen, sections) }
	return func() [][]string { return seen }
}

// TestTicket223R3HotTierReachesALiveReaderWithoutRestart is AC#2's first
// observable: 改了 -> 不重启 -> 行为变了. "Changed behaviour" is read through the
// same call the host's live reader makes (Manager.Config), not through the file.
func TestTicket223R3HotTierReachesALiveReaderWithoutRestart(t *testing.T) {
	m, _, mutate := newTestManager(t, r3Base223)
	notices := r3RestartNotices(m)

	// The reader's answer before the edit: this is the field
	// cmd/wisp/panel_config_store.go:98 hands out, read off s.mgr.Config().
	if got := m.Config().LLM.Roles.Chat.Model; got != "m1" {
		t.Fatalf("before the edit the live reader answers %q, want m1", got)
	}

	mutate(strings.Replace(r3Base223, `model = "m1"`, `model = "m2"`, 1))
	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if rep == nil || !reflect.DeepEqual(rep.Hot, []string{"llm"}) {
		t.Fatalf("report.Hot = %+v, want [llm] (an [llm] edit is D36 immediate-tier)", rep)
	}
	if len(rep.Restart) != 0 || len(rep.Reload) != 0 {
		t.Fatalf("an [llm] edit was also booked restart=%v reload=%v; [llm] is one tier only", rep.Restart, rep.Reload)
	}
	// The behaviour: same call, no restart, a different answer.
	if got := m.Config().LLM.Roles.Chat.Model; got != "m2" {
		t.Fatalf("after one poll the live reader still answers %q, want m2 (改了 -> 不重启 -> 行为变了)", got)
	}
	if got := notices(); len(got) != 0 {
		t.Fatalf("a hot change fired the restart notice %v: that would be a false 本次不生效", got)
	}
	if tier, ok := TierOf("llm"); !ok || tier != "hot" {
		t.Fatalf("TierOf(llm) = %q/%v, want hot - the immediate tier's own claim", tier, ok)
	}
}

// TestTicket223R3RestartTierTellsTheHostWhyAndKeepsTheOldValue is AC#2's second
// observable and AC#7's substance: 改了 -> 不重启 -> 明确告知未生效, plus the why.
// Silence here is the failure, not a passing shape.
func TestTicket223R3RestartTierTellsTheHostWhyAndKeepsTheOldValue(t *testing.T) {
	m, _, mutate := newTestManager(t, r3Base223)
	notices := r3RestartNotices(m)

	mutate(strings.Replace(r3Base223, `language = "zh-CN"`, `language = "en-US"`, 1))
	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if rep == nil || !reflect.DeepEqual(rep.Restart, []string{"app"}) {
		t.Fatalf("report.Restart = %+v, want [app]", rep)
	}
	if len(rep.Hot) != 0 {
		t.Fatalf("a restart-tier edit was also booked hot: %v", rep.Hot)
	}
	// The value did NOT move: that is what makes the notice necessary.
	if got := m.Config().App.Language; got != "zh-CN" {
		t.Fatalf("app.language = %q, want zh-CN (restart tier keeps the old value until process start)", got)
	}
	// Told, not silent, and exactly once for one file change.
	got := notices()
	if len(got) != 1 || !reflect.DeepEqual(got[0], []string{"app"}) {
		t.Fatalf("OnRestartPending = %v, want exactly one [app] notice (静默不生效 is the shape AC#2 forbids)", got)
	}
	// The WHY: the notice is about a key the tier table calls restart, and the
	// same section carries a hot key - so the reason is per key, not "app moved".
	if tier, ok := TierOf("app.language"); !ok || tier != "restart" {
		t.Fatalf("TierOf(app.language) = %q/%v, want restart - this row is the notice's 为什么", tier, ok)
	}
	if tier, ok := TierOf("app.theme"); !ok || tier != "hot" {
		t.Fatalf("TierOf(app.theme) = %q/%v, want hot - otherwise the notice would be a blanket app claim", tier, ok)
	}
}

// TestTicket223R3ReloadTierAppliesAndSignalsTheOwningSubsystem is AC#2's third
// reading, and it has to differ from both others: it applies (like hot) AND it
// signals the owning subsystem (unlike hot). Its boundary is in the file header -
// the consumer of that signal is 票 233's, so this case does not claim it.
func TestTicket223R3ReloadTierAppliesAndSignalsTheOwningSubsystem(t *testing.T) {
	m, _, mutate := newTestManager(t, r3Base223)
	notices := r3RestartNotices(m)
	var reloads [][]string
	m.OnReload = func(sections []string) { reloads = append(reloads, sections) }

	mutate(strings.Replace(r3Base223, `model = "parakeet-zh"`, `model = "paraformer-zh"`, 1))
	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if rep == nil || !reflect.DeepEqual(rep.Reload, []string{"voice"}) {
		t.Fatalf("report.Reload = %+v, want [voice]", rep)
	}
	if got := m.Config().Voice.ASR.Model; got != "paraformer-zh" {
		t.Fatalf("voice.asr.model = %q, want paraformer-zh (reload tier applies immediately)", got)
	}
	if len(reloads) != 1 || !reflect.DeepEqual(reloads[0], []string{"voice"}) {
		t.Fatalf("OnReload = %v, want exactly one [voice] event", reloads)
	}
	if got := notices(); len(got) != 0 {
		t.Fatalf("a reload-tier change fired the restart notice: %v", got)
	}
	if tier, ok := TierOf("voice.asr.model"); !ok || tier != "reload" {
		t.Fatalf("TierOf(voice.asr.model) = %q/%v, want reload", tier, ok)
	}
}

// TestTicket223R3OneEditGivesTwoDifferentObservables is the "two cells, not one"
// pin: one file change moves a hot key and a restart key inside the SAME section,
// and in the same reload the two observables come apart - a live value a reader
// can see (A), and a value that did not move plus an explicit why (B). A refactor
// that merged the tiers, or let the restart notice swallow the hot step, reddens
// one of the two halves below.
func TestTicket223R3OneEditGivesTwoDifferentObservables(t *testing.T) {
	m, _, mutate := newTestManager(t, r3Base223)
	notices := r3RestartNotices(m)

	both := strings.Replace(r3Base223, `theme = "dark"`, `theme = "light"`, 1)
	both = strings.Replace(both, `language = "zh-CN"`, `language = "en-US"`, 1)
	mutate(both)

	rep, err := m.CheckAndReload()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if rep == nil || !reflect.DeepEqual(rep.Hot, []string{"app"}) || !reflect.DeepEqual(rep.Restart, []string{"app"}) {
		t.Fatalf("report = %+v, want the same section in BOTH hot and restart", rep)
	}
	// Observable A: the hot key moved and a reader sees it now.
	if got := m.Config().App.Theme; got != "light" {
		t.Fatalf("app.theme = %q, want light (the hot half did not land)", got)
	}
	// Observable B: the restart key did not move, and the host was told why.
	if got := m.Config().App.Language; got != "zh-CN" {
		t.Fatalf("app.language = %q, want zh-CN (the restart half silently applied)", got)
	}
	if got := notices(); len(got) != 1 || !reflect.DeepEqual(got[0], []string{"app"}) {
		t.Fatalf("OnRestartPending = %v, want one [app] notice", got)
	}
}

// TestTicket223R3EveryAppKeyLandsInTheBucketItsTierNames is the per-key truth
// behind "该键属重启后生效档": the why is only honest if the tier table and the
// manager agree key by key. A key deleted from planApp's restart comparison would
// land in NO bucket (silently ignored) while TierRegistry still called it restart;
// a key moved between tiers would land in the wrong one. Both shapes redden here.
func TestTicket223R3EveryAppKeyLandsInTheBucketItsTierNames(t *testing.T) {
	cases := []struct {
		key  string
		tier string
		edit func(string) string
		read func(*Config) any
	}{
		{"app.theme", "hot", func(s string) string {
			return strings.Replace(s, `theme = "dark"`, `theme = "light"`, 1)
		}, func(c *Config) any { return c.App.Theme }},
		{"app.language", "restart", func(s string) string {
			return strings.Replace(s, `language = "zh-CN"`, `language = "en-US"`, 1)
		}, func(c *Config) any { return c.App.Language }},
		{"app.autostart", "restart", func(s string) string {
			return strings.Replace(s, `language = "zh-CN"`,
				"language = \"zh-CN\"\nautostart = true", 1)
		}, func(c *Config) any { return c.App.Autostart }},
		{"app.single_instance", "restart", func(s string) string {
			return strings.Replace(s, `language = "zh-CN"`,
				"language = \"zh-CN\"\nsingle_instance = false", 1)
		}, func(c *Config) any { return c.App.SingleInstance }},
	}
	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			m, _, mutate := newTestManager(t, r3Base223)
			notices := r3RestartNotices(m)
			want, ok := TierOf(tc.key)
			if !ok || want != tc.tier {
				t.Fatalf("TierOf(%s) = %q/%v, want %q (the table is the why; fix the table first)",
					tc.key, want, ok, tc.tier)
			}
			before := tc.read(m.Config())
			mutate(tc.edit(r3Base223))
			rep, err := m.CheckAndReload()
			if err != nil {
				t.Fatalf("reload: %v", err)
			}
			after := tc.read(m.Config())
			switch tc.tier {
			case "hot":
				if rep == nil || !reflect.DeepEqual(rep.Hot, []string{"app"}) || len(rep.Restart) != 0 {
					t.Fatalf("report = %+v, want app hot and nothing else", rep)
				}
				if reflect.DeepEqual(before, after) {
					t.Fatalf("%s is registered hot but the live value did not move: %v", tc.key, after)
				}
				if got := notices(); len(got) != 0 {
					t.Fatalf("%s is hot but fired the restart notice %v", tc.key, got)
				}
			case "restart":
				if rep == nil || !reflect.DeepEqual(rep.Restart, []string{"app"}) || len(rep.Hot) != 0 {
					t.Fatalf("report = %+v, want app restart and nothing else (a key dropped from planApp lands in no bucket)", rep)
				}
				if !reflect.DeepEqual(before, after) {
					t.Fatalf("%s is registered restart but the live value moved %v -> %v", tc.key, before, after)
				}
				if got := notices(); len(got) != 1 {
					t.Fatalf("%s is restart-tier but the host was told %d times (want 1): %v", tc.key, len(got), got)
				}
			default:
				t.Fatalf("%s tier %q has no case here", tc.key, tc.tier)
			}
		})
	}
}
