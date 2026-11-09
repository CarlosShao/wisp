//go:build windows

package main

// The resident leg's snapshot read face (ticket 65, the wisp.exe half).
//
// WHY THIS FILE EXISTS. The panel snapshot pump had exactly one production
// construction site - run.go's PumpSources literal inside assembleRuntime - and
// the ONLY path that reaches it is startResidentTaskSource -> assembleRuntime
// (cmd/wisp/resident_task_source_windows.go:290). That path is gated on a task
// ENTRY, not on the panel: interactiveStdin() answers nil for a
// -H=windowsgui double-click (scripts/build.ps1, ticket 244 - "a
// -H=windowsgui double-click has no console", resident_approval_windows.go:544),
// so the shipping resident shape returns nil from startResidentTaskSource
// BEFORE the pipeline - and with it the pump - is ever assembled. The panel
// window opens, the settings door (config.get/config.set over the chain's own
// Manager) works, and the snapshot side of the page never receives a byte.
//
// WHAT THIS FILE IS: a pump assembled at the resident boot for exactly the
// shape where the pipeline does not exist, fed by the objects this leg REALLY
// runs - never by shaped constants:
//
//   - gate      ra.gate, the ONE approval gate this process built at boot, so
//     live pending cards and the packet's pending section come off the same
//     queue the orb's cards come from (the mapping body is nativeVerdictsFrom-
//     Gate, the SAME function rt.liveVerdicts calls - no second mapping);
//   - cfg       the chain's own config Manager, returned by newComposerDispatch-
//     Chain. This is the SAME Manager the settings leg writes through
//     (config.set ball_style lands on it via panel_config_store.go), so the
//     picker shows what the next write would act on - the exact property the
//     pump.go BallStyle field doc names. No second Manager is built here;
//   - disp      the chain's router, read only for the mode row: the chain's
//     perm.Store is the process's档 truth (disp.Mode.Modes.PermissionMode);
//   - drafts    the ingest ledger newResidentPanelManager already returns - the
//     same object the attachment door records into;
//   - commands  commandRosterSection213, the once-per-process production roster.
//
// WHAT THIS FILE IS NOT, and what the packet therefore honestly leaves out:
// this leg runs no task loop, so there is no RunFeed, no StreamLog, no task
// roster and no instruction loader behind it. A reader invented for any of
// those would be a constant dressed as a read (the ⓪⓪⓪ ban this tree keeps
// hitting), so Run/Results/Tasks/Instructions readers stay nil and the wire
// carries no such keys - the reader-gated absence shape pump.go pins. When a
// console exists the pipeline assembles and the boot wires the FULLER publisher
// instead (resident_windows.go's selection); this pump is the degraded shape's
// face, not a rival of the pipeline's.
//
// THE PUSH. Out rides the ticket 329 hop's own primitives - panelSnapshotEvalJS
// wrapping the pump's marshal bytes verbatim into the named global, through the
// host's evalPanelJS (the nil-window branch skips silently, the same documented
// skip the pipeline's push has) - and logs the SAME one-line summary the run's
// ledger books, so an operator reading the resident log can see a packet was
// built, when, and its sha256. There is deliberately no second booking path
// here: the ledger booking is the run's face, and this leg has no run.
//
// THE TRIGGER is resync-on-show: the frozen contract (PLAN.md:2982, SPEC-08
// §5.1 - "面板每次 show 必须发 panel.resync，由 Go 侧全量推状态") that no leg had
// wired. showOnThread calls the publisher after a successful Show, so opening
// the panel costs one full packet built from current state - which is what the
// settings page's appearance rows read. State-move publishes in the no-pipeline
// shape stay this ticket's named non-goal: the only writer this shape can move
// (a mode loosening through the orb card) reaches the page on the next show.

import (
	"errors"
	"log/slog"

	"github.com/CarlosShao/wisp/internal/agent/approval"
	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/panel"
)

// residentPanelPump is the no-pipeline shape's publisher: one pump over the
// resident leg's real sources plus the host the bytes are pushed to. The pump
// itself owns no state beyond its publish counter, so publish() is safe from
// the panel thread (where showOnShow runs it) and from anywhere else.
type residentPanelPump struct {
	pump *panel.SnapshotPump
}

// newResidentPanelPump assembles the degraded shape's read face.
func newResidentPanelPump(gate *approval.Gate, cfg *config.Manager, disp *panel.ComposerDispatch, drafts *attachmentDraftLedger, host panelSnapshotHost) *residentPanelPump {
	return &residentPanelPump{pump: panel.NewSnapshotPump(residentPumpSources65(gate, cfg, disp, drafts, host))}
}

// residentPumpSources65 is the pump's wiring, as its own function so a test can
// take the REAL sources, strip one reader, and watch the key vanish - a
// mutation of this injection, not of a hand-copied struct. Every reader below
// is a live object handed in by the boot; a nil cfg (a chain that never
// assembled) leaves the config-derived readers off the pump instead of backing
// them with zero values - a packet that says "no reading" beats one that
// invents a style word or an empty model menu.
func residentPumpSources65(gate *approval.Gate, cfg *config.Manager, disp *panel.ComposerDispatch, drafts *attachmentDraftLedger, host panelSnapshotHost) panel.PumpSources {
	src := panel.PumpSources{
		// Live pending cards off the process's one gate. The mapping is the
		// shared one - rt.liveVerdicts calls the same function - so a card on
		// the orb and a card in the packet cannot disagree about what L1 meant.
		Verdicts: func() []panel.NativeVerdict { return nativeVerdictsFromGate(gate) },
		// The ingest ledger the attachment door records into (ticket 214): the
		// rows and the current refusal sentence, read - never re-shaped.
		Attachments:      drafts.rows,
		AttachmentReason: drafts.reason,
		// Ticket 213 交付 3's roster, derived from internal/command's
		// production table once per process - the menu's members.
		Commands: commandRosterSection213,
		Out: func(snap panel.Snapshot, data []byte) error {
			if host == nil {
				return errors.New("panel: 常驻快照没有宿主可推（面板宿主未装配）")
			}
			// Ticket 329 form 甲, verbatim wrapping: the pump's own bytes into
			// the named global, no re-encode. A page without the receiver
			// drops the call - the documented page half of that ticket
			// (.scratch/wisp/probes/329/r1/page-contract.md §4).
			host.evalPanelJS(panelSnapshotEvalJS(data))
			// The same bounded summary the run's ledger books, into the
			// resident log: a packet the operator can audit without a window.
			slog.Info(panelSnapshotSummary(snap, data))
			return nil
		},
	}
	if cfg != nil {
		// Ticket 65's appearance dimension, read off the SAME Manager the
		// settings door writes through. The vocabulary rides verbatim from
		// the schema; nothing here spells glass or jarvis.
		src.BallStyle = func() panel.BallStyleSection { return ballStyleSectionFromConfig(cfg) }
		// Ticket 187's picker row, same Manager, same verbatim rule.
		src.ModelSwitch = func() panel.ModelSwitchSection { return modelSwitchSectionFromConfig(cfg) }
		// Ticket 215a's inventory, same Manager, same pre-shaped rule.
		src.Plugins = func() panel.PluginInventorySection { return pluginsSectionFromConfig(cfg) }
	}
	// The档 row: the chain's perm store is this process's mode truth. Read off
	// the router the chain built, only when the row is really the production
	// handler - anything else stays unread (the pump then reports the mode as
	// unknown rather than a made-up tier).
	if disp != nil {
		if mw, ok := disp.Mode.(*panel.ModeWriteHandler); ok && mw.Modes != nil {
			src.Mode = mw.Modes.PermissionMode
		}
	}
	return src
}

// publish drives one packet through the exit. Called from showOnThread after a
// successful Show (resync-on-show); the nil-window branch inside evalPanelJS
// makes a publish without a window a counted no-op, never an error dialog.
func (r *residentPanelPump) publish() {
	if r == nil || r.pump == nil {
		return
	}
	if _, _, err := r.pump.Publish(); err != nil {
		slog.Error("panel host: 常驻快照未送达", "err", err)
	}
}

// ballStyleSectionFromConfig is the config->panel mapping for ticket 65's
// appearance dimension. The chain Manager is the same object the settings leg
// writes through, so this is the canonical home of the mapping; run.go's
// PumpSources literal still carries an inline copy over the pipeline Manager
// (that file is the shared assembly root and was carrying another leg's work
// when this landed) - converging it onto this helper is the named follow-up.
func ballStyleSectionFromConfig(m *config.Manager) panel.BallStyleSection {
	v := m.BallStyleOptions()
	return panel.BallStyleSection{Style: v.Style, Vocabulary: v.Vocabulary}
}

// modelSwitchSectionFromConfig is the config->panel mapping for ticket 187's
// picker row, over the same Manager. Same convergence note as above.
//
// Ticket 145 AC#2b: the catalog rows' Display names ride here too
// (ChatSwitchOptions carries them off the same walk). The adapter vocabulary
// does NOT: this is the no-pipeline shape, no provider chain was ever built
// in this process, and llm.ThinkingLevels needs a protocol word off a built
// provider - so Adapter stays empty, which is the same honest no-source
// posture this file's log line states (无源不冒充), never a guessed protocol.
func modelSwitchSectionFromConfig(m *config.Manager) panel.ModelSwitchSection {
	v := m.ChatSwitchOptions()
	models := make([]panel.ModelOption, 0, len(v.Models))
	for _, opt := range v.Models {
		models = append(models, panel.ModelOption{Provider: opt.Provider, Model: opt.Model, Display: opt.Display})
	}
	return panel.ModelSwitchSection{
		Current:         panel.ModelOption{Provider: v.CurrentProvider, Model: v.CurrentModel, Display: v.CurrentDisplay},
		CurrentThinking: v.CurrentThinking,
		Models:          models,
		ThinkingLevels:  v.ThinkingLevels,
		Reason:          v.Reason,
	}
}

// pluginsSectionFromConfig is the config->panel mapping for ticket 215a's
// inventory, over the same Manager. The two sentences are the carrier's own
// constants (internal/panel must not import internal/config, so the section
// arrives pre-shaped exactly like Commands/ModelSwitch/BallStyle). Same
// convergence note as above.
func pluginsSectionFromConfig(m *config.Manager) panel.PluginInventorySection {
	c := m.Config()
	sec := panel.PluginInventorySection{
		Tier2Enabled: c.Plugins.Tier2Enabled,
		EngineState:  panel.PluginEngineNotWired215,
		SkillState:   panel.SkillsNotInConfigModel215,
	}
	for id, e := range c.Plugins.Entries {
		sec.Rows = append(sec.Rows, panel.PluginRow{
			ID:           id,
			Enabled:      e.Enabled,
			Capabilities: e.Capabilities,
			NetAllowlist: e.NetAllowlist,
			HostAPI:      e.HostAPI,
		})
	}
	return sec
}
