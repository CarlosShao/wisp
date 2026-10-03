package main

// Ticket 255 AC#1 - which hot-tier sections this host may call "已立即生效".
//
// WHY THIS FILE EXISTS. cmd/wisp/config_reload.go's reload receipt used to print
// `rep.Hot` verbatim:
//
//	"wisp run: 配置热加载：这些段已立即生效（D36 立即档）：%v\n"
//
// `rep.Hot` is produced by internal/config/manager.go's plan(), and plan() books a
// section there when the *values differ and the memory got overwritten* - nothing
// in that path asks whether any component reads the section. So the sentence
// reported "已立即生效" off a memory write. 票 255's 现量 says the same thing in
// the owner's words: 「它报的是内存里的值换掉了，不是有人按新值做了事」. A hand edit
// of [panel] width therefore printed a claim that cannot be true: the shipped
// panel host never receives a config object at all
// (cmd/wisp/panel_host_windows.go:304-306 hard-codes Width: 420 / Height: 260, and
// NewPanelManager takes no config - that break is ticket 255 AC#4's job, not this
// file's).
//
// WHAT THIS FILE ADDS. One roster, keyed by the SAME paths internal/config's
// TierRegistry uses, saying for each hot-tier row what this host's reader is. The
// receipt may print "已立即生效" for a section only when every hot row behind it is
// verdicted "consumed: ..." - a reader that reads the LIVE config (through
// Manager.Config() at use time) and acts on the value. Everything else hot-applied
// gets its own honest sentence instead: 「值已换进本进程内存，但没有会按新值做事的读
// 者」. The whole line is never deleted - 票 255 AC#1 forbids摘尺当修尺, and the
// [fs] case in config_reload.go already sets the precedent of saying "memory moved,
// this run's verdicts did not" out loud instead of letting the operator infer.
//
// WHY IT ASKS config.TierOf AND NOT A COPY OF THE TIER WORDS. Ticket 255 AC#2-ⓑ
// landed the tier table in internal/config/tiers.go (commit dd92bb92) with TierOf
// as its accessor - and TierOf had ZERO production callers (measured at HEAD
// 8a3790f0: grep -rn "TierOf(" over cmd internal tools scripts returns the
// definition and nothing else). A receipt that hand-listed hot sections would be a
// fourth tier word table (the census's T1-T4 problem). So hotRowsFor() resolves
// each name plan() booked through TierOf/TierRegistry, and manager.go's own
// same-source guard (plan()'s panic when a walked section is not registered "hot")
// stays exactly as it was.
//
// WHAT THE TWO CLASSES OF "有读者" MEAN HERE, STATED SO A READER CAN DISAGREE.
// 票 180's census classes a key "R" when production code reads its value at all.
// By that ruler [agent] and [models] have readers. This file claims "已立即生效"
// only for the narrower class - a reader that re-reads the live Manager - because
// the sentence is about THIS run: cmd/wisp/run.go:423-424 takes `cfg := mgr.Config()`
// once and its own comment says "the runtime's non-mode settings stay frozen at
// what this boot read", and execute() reads that copy (cmd/wisp/run.go:991
// `cfg := rt.cfg`). Calling [agent] 立即生效 on that shape is the same lie as
// [panel], one layer deeper. Both files' evidence is in the roster strings below and
// is pinned by config_receipt_255_test.go, so the disagreement is checkable rather
// than buried.

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/CarlosShao/wisp/internal/config"
)

// Claim prefixes of hotRowClaims. "consumed: " is the only one the receipt may
// print as 已立即生效; the rest are honest verdicts naming what the value's reader
// actually is. Adding a verdict here without adding the row to hotRowClaims fails
// TestTicket255HotRowRosterCoversTheRegistry.
const (
	hotClaimConsumed      = "consumed: "
	hotClaimSnapshotOnly  = "snapshot-only: "
	hotClaimNoReader      = "no-reader: "
	hotClaimOtherProcess  = "other-process: "
	hotClaimDebugHostOnly = "debug-host-only: "
)

// hotRowClaims is the AC#1 roster: one row per TierRegistry path that plan() can
// book into Report.Hot, carrying this host's verdict for it.
//
// Every row cites its evidence as `path/to/file.go:LINE [token]`, and
// TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim re-reads the source and
// requires that line to still contain that token. A file:line in here is therefore
// a checked claim, not a comment.
var hotRowClaims = map[string]string{
	// The one hot section whose value reaches code that acts on it mid-run: the
	// settings leg answers the panel out of the LIVE Manager on every call
	// (cmd/wisp/panel_config_store.go:92 [cfg := s.mgr.Config()], surfaced at
	// :96 [cfg.LLM.Roles.Chat.Model]), so a hand-edited [llm] role/provider
	// changes what the next settings read says without a restart. Caveat this
	// file owns: the boot snapshot still drives the running provider chain
	// (cmd/wisp/run.go:435 [res := llm.NewResolver(cfg, st)]), so "已立即生效"
	// here means "the live reader answers differently", not "the in-flight
	// model swapped".
	"llm": hotClaimConsumed + "cmd/wisp/panel_config_store.go:96 [cfg.LLM.Roles.Chat.Model] - configStore.ReadSettings calls s.mgr.Config() per call",

	// Read in production, but only from the boot snapshot: rt.cfg is taken once
	// at assembly (cmd/wisp/run.go:424 [rt.cfg = cfg]) and every use reads that
	// copy (cmd/wisp/run.go:991 [cfg := rt.cfg]; the field is consumed at
	// cmd/wisp/run.go:1014 [cfg.Agent.PerToolTimeoutMS]). plan() did overwrite
	// the memory, so the honest sentence is "值已换进内存", not "已立即生效".
	"agent": hotClaimSnapshotOnly + "cmd/wisp/run.go:991 [cfg := rt.cfg] - every cfg.Agent read is the boot copy, so this run keeps acting on the old value",

	// Read by a DIFFERENT command's process: openModelStore takes the config
	// (cmd/wisp/models.go:163 [cfg.Models.Mirror]) but its callers are the
	// `wisp models` subcommands, each of which LoadFile's fresh
	// (cmd/wisp/models.go:183 [func loadModelConfig]). The running `wisp run`
	// host never opens the model store, so a mid-run [models] edit changes the
	// next command, not this run.
	"models": hotClaimOtherProcess + "cmd/wisp/models.go:163 [cfg.Models.Mirror] - only the wisp models subcommands read it, in their own process, off a fresh LoadFile",

	// Only reader in the tree is the side debug program
	// (cmd/balldebug/main.go:238 [h := mgr.Config().Hotkey]); the shipped ball
	// binds ball.DefaultHotkeys() instead
	// (cmd/wisp/resident_ball_windows.go:171 [Hotkeys:  ball.DefaultHotkeys(),]).
	"hotkey": hotClaimDebugHostOnly + "cmd/balldebug/main.go:238 [h := mgr.Config().Hotkey] - cmd/balldebug is not the shipped host, and wisp run binds ball.DefaultHotkeys()",

	// Zero production readers (票 180's census classes these keys "D"). Two shapes
	// of evidence, and the roster test demands one of them per row: a cite naming
	// the code that decides the value instead of reading it, or the marker
	// "扫描零命中", which is the reader scan's own result for that field
	// (TestTicket255RosterStillMatchesTheActualReadSites in
	// config_receipt_255_test.go re-measures it on every run).
	"ball":    hotClaimNoReader + "internal/ball/ball_windows.go:64 [SizePx] - the ball's own option carries the size, and nothing outside internal/config reads cfg.Ball",
	"session": hotClaimNoReader + "扫描零命中：nothing outside internal/config reads cfg.Session - the three session timeouts this section names drive nothing yet",
	"audio":   hotClaimNoReader + "扫描零命中：nothing outside internal/config reads cfg.Audio - internal/audio takes its device and rate from elsewhere",
	"privacy": hotClaimNoReader + "internal/observe/logging.go:50 [RedactPaths] - the mirror field exists and is filled by callers, never from cfg.Privacy",
	"memory":  hotClaimNoReader + "扫描零命中：nothing outside internal/config reads cfg.Memory - the L1/L3 knobs have no consumer",
	"cost":    hotClaimNoReader + "扫描零命中：nothing outside internal/config reads cfg.Cost - the C23 budget knobs (ticket 44) are not landed",
	"observe": hotClaimNoReader + "cmd/wisp/logsink.go:149 [Level: logSinkLevel] - the log sinks choose their level themselves, never from cfg.Observe",

	// [panel]: THE section 票 255 was立 for. Zero production readers, and the
	// host that would use width/height hard-codes its own geometry
	// (cmd/wisp/panel_host_windows.go:304 [Width:  420,]). Wiring it is AC#4
	// and needs a real window-width reading, which is why this row is a verdict
	// and not a wire-up done here.
	"panel": hotClaimNoReader + "nothing outside internal/config reads cfg.Panel - cmd/wisp/panel_host_windows.go:304 [Width:  420,] is hard-coded and NewPanelManager receives no config (票 255 AC#4 owns that break)",

	// Per-key sections: planApp/planVoice book the SECTION name into rep.Hot, so
	// the claim must be settled per key row and only claimed when every hot row
	// behind the name is consumed.
	"app.theme":                  hotClaimNoReader + "扫描零命中：nothing outside internal/config reads cfg.App - manager.go's planApp compares and copies it, and no component re-skins from it",
	"voice.tts.speed":            hotClaimNoReader + "扫描零命中：nothing outside internal/config reads cfg.Voice - the TTS knobs are not threaded to the voice path yet",
	"voice.punctuation":          hotClaimNoReader + "扫描零命中：nothing outside internal/config reads cfg.Voice - the punctuation switch has no consumer yet",
	"voice.wake_word.thresholds": hotClaimNoReader + "扫描零命中：nothing outside internal/config reads cfg.Voice - the wake-word tunables are not consumed by the KWS path yet",
	"voice.wake_word.veto_words": hotClaimNoReader + "扫描零命中：nothing outside internal/config reads cfg.Voice - the veto list is not consumed by the KWS path yet",
}

// sectionReadSites pins, per Config field, which production files the reader scan
// in TestTicket255RosterStillMatchesTheActualReadSites is allowed to hit. It is
// the roster's other half: a NEW file reading a section nobody claimed, or a
// claimed file that stopped reading it, reddens that test instead of quietly
// making the receipt untrue. The values are file paths, not line numbers, so an
// unrelated edit above a reader does not churn the pin.
var sectionReadSites = map[string][]string{
	"Ball":    nil,
	"Session": nil,
	"Audio":   nil,
	"Privacy": nil,
	"Memory":  nil,
	"Panel":   nil,
	"Cost":    nil,
	"Observe": nil,
	"App":     nil,
	"Voice":   nil,
	"Agent":   {"cmd/wisp/run.go"},
	"Models":  {"cmd/wisp/models.go"},
	"LLM":     {"cmd/wisp/panel_config_store.go", "cmd/wisp/providers.go", "cmd/wisp/run.go", "internal/llm/resolver.go"},
	"Hotkey":  {"cmd/balldebug/main.go"},
}

// evidenceCite matches the `path.go:LINE [token]` shape every roster row is
// required to carry at least once.
var evidenceCite = regexp.MustCompile(`([\w./-]+\.go):(\d+) \[([^\]]+)\]`)

// hotRowsFor resolves one name plan() booked into Report.Hot to the TierRegistry
// rows behind it, through config.TierOf - the same table plan()'s own guard
// consults. An empty answer means plan() and the registry disagree, which
// splitHotTier reports loudly instead of claiming anything.
func hotRowsFor(name string) []string {
	if tier, ok := config.TierOf(name); ok { // TierOf's first production caller
		if tier != "hot" {
			return nil
		}
		return []string{name}
	}
	// [app] and [voice] carry no section-level row: their tiers are per key, so
	// the rows behind the name are exactly its registered "hot" keys.
	rows := []string{}
	for path, tier := range config.TierRegistry {
		if tier == "hot" && strings.HasPrefix(path, name+".") {
			rows = append(rows, path)
		}
	}
	sort.Strings(rows)
	return rows
}

// hotSectionSplit is one reload's answer to "may I say 已立即生效 about this".
type hotSectionSplit struct {
	claimable []string // every hot row behind it is verdicted consumed
	quiet     []string // hot-applied, but no reader here acts on the new value
	disagree  []string // plan() booked it hot and the registry cannot explain it
}

// splitHotTier sorts Report.Hot by what this host's roster says reads it. It never
// invents a section: a name only reaches claimable when config.TierOf (or its
// per-key rows) says "hot" AND every row behind it is "consumed: ...".
func splitHotTier(hot []string) hotSectionSplit {
	var s hotSectionSplit
	for _, name := range hot {
		rows := hotRowsFor(name)
		if len(rows) == 0 {
			s.disagree = append(s.disagree, name)
			continue
		}
		consumed := true
		for _, row := range rows {
			verdict, ok := hotRowClaims[row]
			if !ok || !strings.HasPrefix(verdict, hotClaimConsumed) {
				consumed = false
			}
		}
		if consumed {
			s.claimable = append(s.claimable, name)
		} else {
			s.quiet = append(s.quiet, name)
		}
	}
	return s
}

// verdictFor names, for the audit trail, which rows a hot section's claim rests on
// and what each row's verdict is. The receipt's stdout line is deliberately short;
// the per-row truth belongs on the ledger where 票 255's 「谁读它」 question can be
// answered after the fact.
func (s hotSectionSplit) verdictFor(name string) string {
	rows := hotRowsFor(name)
	if len(rows) == 0 {
		return "tier_row=[] reason=plan() booked it hot but config.TierRegistry has no hot row for it"
	}
	parts := make([]string, 0, len(rows))
	for _, row := range rows {
		verdict, ok := hotRowClaims[row]
		if !ok {
			parts = append(parts, fmt.Sprintf("%q:no verdict in hotRowClaims (ticket 255 AC#1 roster gap)", row))
			continue
		}
		parts = append(parts, fmt.Sprintf("%q:%s", row, verdict))
	}
	return fmt.Sprintf("tier_row=%v claims=[%s]", rows, strings.Join(parts, " | "))
}

// configSectionsForScan returns the Config fields whose section the roster claims
// to have adjudicated, as Go field names. It walks the schema by reflection so a
// new hot section cannot slip past the roster - the same shape as
// internal/config's tiers_255_test.go rulers.
func configSectionsForScan() []string {
	t := reflect.TypeOf(config.Config{})
	out := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		section := strings.Split(f.Tag.Get("toml"), ",")[0]
		if section == "" || section == "-" {
			continue
		}
		if rosterCoversSection(section) {
			out = append(out, f.Name)
		}
	}
	sort.Strings(out)
	return out
}

// rosterCoversSection answers whether hotRowClaims carries a verdict for this
// section - either the section-level row ([llm], [ball], ...) or at least one key
// row behind it ([app] and [voice], whose tiers are per key).
func rosterCoversSection(section string) bool {
	if _, ok := hotRowClaims[section]; ok {
		return true
	}
	for path := range hotRowClaims {
		if strings.HasPrefix(path, section+".") {
			return true
		}
	}
	return false
}
