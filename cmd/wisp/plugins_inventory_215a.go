package main

// Ticket 215a's production readers at the composition root: the Tier-2 plugin
// inventory's three missing dimensions, each one a real reader instead of a
// sentence.
//
//   - 名字/版本: pluginManifestInfo215 reads the installed plugin's manifest
//     (plugins\<id>\manifest.toml, SPEC-07 §6's on-disk layout) and hands the
//     inventory the manifest's PARSED name and version - two typed strings,
//     nothing else. Before this file, no production code read a plugin
//     manifest at all (the receipt's own measurement: manifestslot had zero
//     production imports), so a manifest on disk was a file nobody read; the
//     AC#1 judgment ("改前某键零读者 ⇒ 改后 ≥1") is carried by this reader
//     and the wire it feeds.
//   - 这一轮装配没装配: assembledTier2Plugins215 is the assembly probe, and
//     pluginRowState215 derives each row's honest 生效与否 label FROM THE
//     PROBE, never from Enabled. A config that says true with nothing
//     assembled behind it renders 未生效 - the A408 shape ("文案说在跑、其实
//     没人读") made visible per row, which is exactly ticket 215 AC#2's
//     positive control, and the cmd/wisp tests pin it in BOTH directions
//     (probe empty + enabled=true => 未生效; probe present => 本轮生效) so
//     the derivation cannot quietly degenerate into reading the boolean.
//   - the shared mapping: pluginInventorySection215 is the ONE config->panel
//     mapping both assembly roots call (run.go's PumpSources literal and the
//     resident pump's pluginsSectionFromConfig), so the two sites cannot
//     drift the way the receipt found run.go's decorative copy drifting.

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CarlosShao/wisp/internal/buildinfo"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
	"github.com/CarlosShao/wisp/internal/plugin"
)

// pluginsRoot215 is the SPEC-02 §6 plugins directory of one data root: the
// dataDir this run resolved at boot, plus the one fixed layout name. The join
// here mirrors manifestslot's own install layout (<PluginsRoot>\<id>\
// manifest.toml); no path decision beyond the join happens in this file.
func pluginsRoot215(dataDir string) string {
	return filepath.Join(dataDir, "plugins")
}

// residentPluginsRoot215 resolves the plugins root for the degraded (-H=
// windowsgui) shape, whose pump wiring holds no runSpec: the same env -> data
// root resolution `wisp config` does (cmdConfig above), so the resident
// inventory reads the same tree the pipeline shape would. An unresolvable
// environment answers "" and the manifest reader then reports honest absence
// rather than guessing a root.
func residentPluginsRoot215() string {
	env, err := buildinfo.ResolveEnv()
	if err != nil {
		return ""
	}
	layout, err := resolveSecretLayout(env)
	if err != nil {
		return ""
	}
	return pluginsRoot215(layout.DataDir)
}

// assembledTier2Plugins215 is THE ASSEMBLY PROBE: the plugin ids the Tier-2
// engine actually assembled THIS ROUND, keyed by id, true = assembled with
// the current config's enabled face.
//
// The honest answer today is an EMPTY probe, and the reason is written here
// so nobody "fixes" it into a guess: tickets 50/51/52 landed the engine's
// components (Tier2Gate, Tier2Executor, the manifestslot lifecycle), but no
// assembly root constructs them yet - cmd/wisp wires no gate, no executor and
// no manifestslot Manager, and internal/config's unwired.go disposition rows
// still name the [plugins] keys "not built". An empty probe is therefore the
// MEASURED state of this process, not a placeholder: every enabled row
// renders 未生效 until the assembly leg lands, and that leg fills this probe
// from the live gate (plugin.Tier2Gate.Snapshot / the manifestslot Manager's
// own live set) in the same commit it retires the sentence and the disposition
// rows. Handing this probe a hardcoded non-empty map before that commit would
// be the lie ticket 215 was filed to stop.
func assembledTier2Plugins215() map[string]bool {
	return nil
}

// pluginRowState215 derives one row's honest 生效与否 label from the probe.
// The rule, in the order the states are checked:
//
//   - the probe carries the id AND the config still says enabled =>
//     本轮生效 (assembled this round, config agrees with what was assembled);
//   - the probe carries the id but the config moved (switched off after
//     assembly) => 需重启 (the consumer only re-reads config at restart;
//     the row says so instead of claiming the running instance already
//     stopped);
//   - anything else => 未生效, which covers BOTH "switched off" and the
//     "配置开着但没人装配" shape the engine's absence makes universal today.
//
// The id parameter is unused today only because the probe is a set of ids;
// it stays in the signature so the derivation stays keyable per row when the
// probe grows past a set (the probe's own comment names the future shape).
func pluginRowState215(id string, enabled bool, assembled map[string]bool) string {
	if assembled[id] {
		if enabled {
			return panel.PluginRowLive215
		}
		return panel.PluginRowRestart215
	}
	return panel.PluginRowNotLive215
}

// pluginManifestInfo215 reads ONE installed plugin's manifest and returns the
// three facts the inventory may carry: the manifest's parsed name, its
// parsed version, and a note when the manifest exists but refused to load.
//
// The read face is the manifest GRAMMAR's own front door (plugin.ParseManifest,
// the C14 closed-schema audit + typed decode + per-field validation), not a
// hand parse: a manifest that tries to carry keys the C14 schema does not
// declare - a "permissions = allow all" row, an allowed_dirs grant, anything
// the schema never froze - is REFUSED by the audit with the key named, before
// any field of it reaches the wire. That refusal is ticket 215 AC#4's first
// tooth; the second is that this function RETURNS three strings and nothing
// else, so manifest prose (the tool descriptions, the review sheet text) has
// no route onto the wire even when the manifest parses clean - reading a
// plugin's说明书 is not an authorization, and this reader structurally
// cannot turn text into config.
//
// Outcomes, each distinguishable:
//   - no manifest on disk (the normal state today: the installer has no
//     production caller) - three empty values, no note: measured absence,
//     never invented (宁缺毋造);
//   - manifest present but unreadable or C14-refused - empty name/version,
//     note = the refusal: the row does not pretend to know who the plugin
//     is, and it says why;
//   - manifest parses - name/version verbatim, empty note.
func pluginManifestInfo215(root, id string) (name, version, note string) {
	if root == "" || id == "" {
		return "", "", ""
	}
	// The id touches the filesystem here, so the same shape guard
	// manifestslot.safePluginID applies gets re-checked: a config id is
	// operator-controlled text, and a separator or a dot segment in it must
	// refuse loudly instead of naming another directory's manifest.
	if id != filepath.Base(id) || id == "." || id == ".." {
		return "", "", fmt.Sprintf("清单未读：插件 id %q 不是单个安全路径段", id)
	}
	raw, err := os.ReadFile(filepath.Join(root, id, "manifest.toml"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", ""
		}
		return "", "", fmt.Sprintf("清单不可读：%v", err)
	}
	manifest, err := plugin.ParseManifest(raw)
	if err != nil {
		return "", "", fmt.Sprintf("清单拒绝：%v", err)
	}
	if manifest.ID != id {
		return "", "", fmt.Sprintf("清单拒绝：manifest 自报 id %q 与目录 %q 不一致", manifest.ID, id)
	}
	return manifest.Name, manifest.Version, ""
}

// pluginInventorySection215 is the ONE config->panel mapping both assembly
// roots use (run.go's PumpSources literal and the resident pump). The two
// sentences are the carrier's own constants - no forked prose - and every
// row's name/version/state/note comes from the readers above, so the wire
// says what this process MEASURED (manifest on disk or not, probe hit or
// not), never what the config merely asserts.
func pluginInventorySection215(c *config.Config, assembled map[string]bool, root string) panel.PluginInventorySection {
	sec := panel.PluginInventorySection{
		Tier2Enabled: c.Plugins.Tier2Enabled,
		EngineState:  panel.PluginEngineNotWired215,
		SkillState:   panel.SkillsNotInConfigModel215,
	}
	for id, e := range c.Plugins.Entries {
		name, version, note := pluginManifestInfo215(root, id)
		sec.Rows = append(sec.Rows, panel.PluginRow{
			ID:           id,
			Name:         name,
			Version:      version,
			Enabled:      e.Enabled,
			State:        pluginRowState215(id, e.Enabled, assembled),
			Note:         note,
			Capabilities: e.Capabilities,
			NetAllowlist: e.NetAllowlist,
			HostAPI:      e.HostAPI,
		})
	}
	return sec
}
