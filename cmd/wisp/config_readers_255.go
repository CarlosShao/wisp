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
// of [panel] width therefore printed a claim that could not be true: at the time
// this file was written the shipped panel host received no config object at all.
// The evidence is quoted as CONTENT, not as a line number: at HEAD 67ab595d the
// create call inside bringUp spelled its geometry as the literals `Width: 420` /
// `Height: 260` (grep those two in
// `git show 67ab595d:cmd/wisp/panel_host_windows.go`), and NewPanelManager took no
// config parameter. No line number is quoted here on purpose - that cold-start
// block has moved twice since (its tail went into coldStartPageHandover in 33-r5,
// and 33-v4 named the rot on 10-08), so a bare file:LINE would have been a stale
// pointer by now. 票 255 AC#4 has since landed in this same package: the
// resident assembly root hands the host a per-creation geometry closure, and the
// [panel] row below is what that changed in THIS file's claim. What did not
// change is the judgement this file exists to enforce - a section is only claimed
// 立即生效 when a reader in this host acts on the new value, and the panel host is
// still not that reader in `wisp run`.
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
// as its accessor - and TierOf had ZERO production callers at that point (measured
// at HEAD 8a3790f0: grep -rn "TierOf(" over cmd internal tools scripts returned
// the definition and nothing else). This file is its first production caller
// since. A receipt that hand-listed hot sections would be a
// fourth tier word table (the census's T1-T4 problem). So hotRowsFor() resolves
// each name plan() booked through TierOf/TierRegistry, and manager.go's own
// same-source guard (plan()'s panic when a walked section is not registered "hot")
// stays exactly as it was.
//
// WHAT THE TWO CLASSES OF "有读者" MEAN HERE, STATED SO A READER CAN DISAGREE.
// 票 180's census classes a key "R" when production code reads its value at all.
// By that ruler [agent] and [models] have readers. This file claims "已立即生效"
// only for the narrower class - a reader that re-reads the live Manager. Until
// ticket 233 execute() read the boot copy (cmd/wisp/run.go:1023 `cfg := rt.cfg`,
// taken once at run.go:426-427), which is what kept [agent] snapshot-only while
// its own comment said "the runtime's non-mode settings stay frozen at what this
// boot read" - calling [agent] 立即生效 on that shape was the same lie as [panel],
// one layer deeper. Ticket 233 re-pointed execute() at the live Manager
// (run.go:1252 `cfg := rt.mgr.Config()`), so [agent] now carries "consumed" with
// the new read site as its cite; the boot snapshot rt.cfg still exists, still
// frozen, and the 223/227/255 cases still read it as the "before". Both files'
// evidence is in the roster strings below and is pinned by
// config_receipt_255_test.go, so the disagreement is checkable rather than
// buried.

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
// Every row cites its evidence as `path/to/file.go:LINE [token]`. The token is the
// anchor and the :LINE is a reading: 票 255 AC#6 (10-10) fixed the criterion as
// 「token 唯一命中当判据，行号只当读数」, and
// TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim requires the token to
// occur EXACTLY ONCE in the named file. An insertion above the line moves the
// reading but leaves the token alone, so it must not churn this roster (before
// AC#6 the ruler reddened on exactly that, which punished an unrelated edit); a
// changed token reddens it, naming the row and the line it was recorded at. A
// file:line in here is therefore a checked claim, not a comment - and each token
// is sharpened until it pins one line ("SizePx   int", not the bare "SizePx"
// which 13 lines carry), so the anchor can never silently mean a lookalike.
var hotRowClaims = map[string]string{
	// The one hot section whose value reaches code that acts on it mid-run: the
	// settings leg answers the panel out of the LIVE Manager on every call
	// (cmd/wisp/panel_config_store.go:92 [cfg := s.mgr.Config()], surfaced at
	// :98 [ChatModel:      cfg.LLM.Roles.Chat.Model]), so a hand-edited [llm] role/provider
	// changes what the next settings read says without a restart. Caveat this
	// file owns: the boot snapshot still drives the running provider chain
	// (cmd/wisp/run.go:438 [res := llm.NewResolver(cfg, st)]), so "已立即生效"
	// here means "the live reader answers differently", not "the in-flight
	// model swapped".
	//
	// 2026-11-05 re-adjudication (ticket 98's C2): the reader scan also names
	// three files that take a *config.Config and reach into this section, and
	// they are pinned in sectionReadSites below rather than cited here because
	// each one's verdict depends on WHO hands it the config - 票 44's
	// internal/agent/costbook_44.go + internal/agent/quota_44r4.go (price
	// cards, provider cards and quota params; execute() hands them the LIVE
	// per-task cfg) and 票 61's internal/speech/cloudcatalog.go (cloud chain
	// elements, resolved from a snapshot the caller supplies; its
	// ResolveCloudChain still has zero non-test callers). None of the three
	// changes this row's verdict: the claim still rests on the settings leg's
	// per-call re-read that this cite names.
	"llm": hotClaimConsumed + "cmd/wisp/panel_config_store.go:98 [ChatModel:      cfg.LLM.Roles.Chat.Model] - configStore.ReadSettings calls s.mgr.Config() per call; the scan also names internal/agent/costbook_44.go, internal/agent/quota_44r4.go and internal/speech/cloudcatalog.go, which read the section off a config they are handed (pinned in sectionReadSites)",

	// Since ticket 233 this is a live reader, and the row moved to "consumed":
	// execute() (cmd/wisp/run.go:1252 [cfg := rt.mgr.Config()]) now reads the
	// LIVE Manager at the start of every task, so a hot-applied [agent] edit
	// reaches the NEXT task's PerToolTimeoutMS / SteeringEnabled /
	// ProjectInstructionsEnabled instead of every task acting on the boot copy
	// forever. Before that line the reads went through the boot snapshot
	// (cmd/wisp/run.go:1023 [cfg := rt.cfg]), which is what kept this row
	// snapshot-only; rt.cfg still exists and is still the frozen boot snapshot
	// its own pin (and the 223/227/255 cases) read it as, but execute() no
	// longer sources its config from it. This is the one hot section whose
	// claim upgraded with ticket 233 - [llm] is the other "consumed" row, its
	// live reader unchanged (cmd/wisp/panel_config_store.go:92).
	"agent": hotClaimConsumed + "cmd/wisp/run.go:1263 [cfg = rt.mgr.Config()] - execute() reads the live Manager per task since ticket 233, so the next task after a hot-applied [agent] edit acts on the new values (PerToolTimeoutMS / SteeringEnabled / ProjectInstructionsEnabled, and since ticket 423 also MaxRounds / TokenBudget / LoopGuard.RepeatThresholds, which execute() hands to agent.Config at the loop's construction - the ticket-10 tail the census caught feeding from nowhere)",

	// Read by a DIFFERENT command's process: openModelStore takes the config
	// (cmd/wisp/models.go:163 [cfg.Models.Mirror]) but its callers are the
	// `wisp models` subcommands, each of which LoadFile's fresh
	// (cmd/wisp/models.go:183 [func loadModelConfig]). The running `wisp run`
	// host never opens the model store, so a mid-run [models] edit changes the
	// next command, not this run.
	"models": hotClaimOtherProcess + "cmd/wisp/models.go:163 [cfg.Models.Mirror] - only the wisp models subcommands read it, in their own process, off a fresh LoadFile",

	// Since ticket 258 (form A): the resident leg's ball host takes [hotkey]
	// from config (cmd/wisp/resident_ball_windows.go:322 [Hotkeys:  cfg,],
	// mapped by the chain in cmd/wisp/resident_windows.go) and its ticket-64
	// bridge rebinds the ball on a real diff. The reader is real, but it lives
	// in the RESIDENT `wisp` process - this receipt prints in `wisp run`, which
	// builds no ball - so the row is other-process, the same shape [panel] was
	// adjudicated into when its reader moved to the resident panel host. The
	// rebind drops an in-flight Esc borrow (internal/ball/ball_windows.go:807-812);
	// that residual stays ticket 245's, named rather than hidden.
	"hotkey": hotClaimOtherProcess + "cmd/wisp/resident_ball_windows.go:322 [Hotkeys:  cfg,] - the resident ball host binds from [hotkey] and its reloader bridge rebinds on change (cmd/wisp/resident_windows.go), but the ball lives in the resident process, not in this `wisp run` host",

	// Ticket 419 re-adjudicated this row, and it stays no-reader while the two
	// fields behind the section carry different facts, so the verdict text says
	// each field separately instead of one sentence covering both. The SIZE
	// half is still what this row has always described: ball_windows.go:64's
	// SizePx is the option the configured size would ride, and no code feeds it
	// (ticket 413 is the leg about to feed it), so the hot-applied size still
	// reads from nowhere. The STYLE half is the change that reddened the scan:
	// ticket 65's leg gave cfg.Ball.Style two production read sites outside
	// internal/config, so the old "nothing outside internal/config reads
	// cfg.Ball" is no longer sayable. Neither site is a live reader - each is
	// one fresh config.LoadFile (the resident ball host's boot seed, and
	// `wisp config get` in its own process), neither re-reads the live Manager,
	// and this `wisp run` host builds neither - so a hot-applied [ball] edit
	// reaches the next boot or the next CLI invocation, never a running
	// process, which is why the verdict does not move to other-process the way
	// [audio]'s did when its resident reader landed.
	// Ticket 413 re-adjudicated the SIZE half of this row, and the verdict
	// stays no-reader. The gap this row described is closed at the feed side:
	// cmd/wisp/ball_size_413_windows.go:48 [c.Ball.Size] is a production read
	// site now - the resident ball host's boot seed, one fresh config.LoadFile
	// whose answer the assembly root hands to ball.New's SizePx
	// (cmd/wisp/resident_ball_windows.go [SizePx:   hs.sizePx,]) - so the
	// slider's 44-160 reaches the window at construction instead of every boot
	// building at the package default. What keeps the verdict where 419 put it:
	// the read is one construction-time LoadFile, nothing re-reads the live
	// Manager (the ball carries no live re-size, so a hot-applied [ball] size
	// waits for the next boot), and the reader lives in the resident `wisp`
	// process while this receipt prints in `wisp run`, which builds no ball.
	// The STYLE half keeps its own sentence: cmd/wisp/ball_style_65_windows.go:50
	// [c.Ball.Style] (the boot style seed) and cmd/wisp/config_cmd_65.go:235
	// [c.Ball.Style, true, nil] (`wisp config get`, a process of its own) are
	// boot/invocation-time LoadFile passes too.
	"ball": hotClaimNoReader + "internal/ball/ball_windows.go:64 [SizePx   int // configured orb size 44..160 (0 = default 56)] - the option the configured size rides; ticket 413 landed the feed: cmd/wisp/ball_size_413_windows.go:48 [c.Ball.Size] (the resident ball host's boot seed) hands the file's number to ball.New through cmd/wisp/resident_ball_windows.go, so the size reaches the window at construction; the verdict stays no-reader because the read is one boot-time LoadFile that never re-reads the live Manager and the ball has no live re-size, and the style field keeps its own sentence: cmd/wisp/ball_style_65_windows.go:50 [c.Ball.Style] (the boot style seed) and cmd/wisp/config_cmd_65.go:235 [c.Ball.Style, true, nil] (`wisp config get`, a process of its own) are boot/invocation-time LoadFile passes that never re-read the live Manager, and this `wisp run` host reads neither",
	// Ticket 233 re-adjudicated this row: since ticket wire-r1 the resident
	// voice leg reads one [session] key at boot (cmd/wisp/resident_voice_windows.go
	// [c.Session.WarmTimeoutSec], handed to the voice plan's warm timeout), so
	// the old 扫描零命中 reading is no longer true. Two facts keep the verdict in
	// other-process, the same shape [audio] carries: the read is one fresh
	// config.LoadFile at the voice leg's construction (a changed [session]
	// waits for the next process; nothing re-reads the live Manager), and the
	// reader lives in the resident `wisp` process - `wisp run` builds no voice
	// leg. Ticket 423 wired the other two [session] keys the same way:
	// settling_sec and conversation_idle_sec ride the plan into
	// session.ScopeOptions' SettleStatic / ConversationIdle slots (which had
	// stood empty since ticket 28 with "0 = default" contracts), so all three
	// keys now reach the scope at its construction.
	"session": hotClaimOtherProcess + "cmd/wisp/resident_voice_windows.go:1329 [c.Session.WarmTimeoutSec] - the resident voice leg reads the warm timeout once at its boot LoadFile, so a hand edit reaches the next process, never this run; ticket 423 feeds the other two keys the same way (cmd/wisp/resident_voice_windows.go:1848 [c.Session.ConversationIdleSec] and cmd/wisp/resident_voice_windows.go:1849 [c.Session.SettlingSec], both handed to session.ScopeOptions at the scope's construction)",
	// Ticket 247 moved [audio] out of the reader-less class: the resident
	// process's capture assembly reads mic_muted_default at boot and hands it to
	// audio.WithStartMuted, which is the sentence internal/audio/gate.go:55-56 wrote
	// when it asked for that boot wiring ("map [audio] mic_muted_default here at
	// boot wiring") and got none. Two facts keep this row out of "consumed":
	// the reader is a fresh config.LoadFile at boot, not a re-read of the live
	// Manager (a changed [audio] waits for a restart), and it lives in the
	// resident `wisp` process, which builds the capture leg - `wisp run` does
	// not. The device/rate half of the section still reads from nowhere:
	// internal/audio takes input_device and sample_rate from the OS default
	// endpoint and the seam constant, not from cfg.Audio.
	"audio":   hotClaimOtherProcess + "cmd/wisp/resident_audio_windows.go:256 [ra.mutedAtBoot = c.Audio.MicMutedDefault] - the resident capture leg is [audio]'s first production reader and it reads the file once at boot, so a hot reload of this section waits for a restart instead of acting in this run",
	"privacy": hotClaimNoReader + "internal/observe/logging.go:51 [RedactPaths bool] - the mirror field exists and is filled by callers, never from cfg.Privacy",
	"memory":  hotClaimNoReader + "扫描零命中：nothing outside internal/config reads cfg.Memory - the L1/L3 knobs have no consumer",
	// Ticket 44 (C23 money path) landed its composition-root wiring after this
	// row was written, so the 2026-11-05 re-adjudication moved [cost] out of the
	// reader-less class: execute() reads the LIVE Manager at the top of every
	// task (cmd/wisp/run.go:1524 [cfg = rt.mgr.Config()]) and hands that same
	// cfg to agent.CostWiring (run.go:1570) and, through quotaDispatchChain, to
	// agent.QuotaWiring - and both read the four [cost] leaves below. So a
	// hot-applied [cost] edit reaches the NEXT task's meter and quota gate with
	// no restart, the same per-task shape [agent] carries (and the same caveat:
	// the task in flight keeps the caps it started with). The old verdict
	// ("扫描零命中 ... the C23 budget knobs (ticket 44) are not landed") became
	// false the day 44-r3/44-r5 landed, which is exactly the drift ticket 98's
	// census caught in this roster.
	"cost":    hotClaimConsumed + "internal/agent/costbook_44.go:183 [b, err := observe.BudgetFrom(cfg.Cost.DailyBudget, cfg.Cost.MonthlyBudget,] - execute() hands the live per-task cfg to agent.CostWiring (cmd/wisp/run.go:1570), whose budget half is built from these four [cost] leaves; agent.QuotaWiring (internal/agent/quota_44r4.go:135) reads the same caps, and both run once per task",
	"observe": hotClaimNoReader + "cmd/wisp/logsink.go:149 [Level: logSinkLevel] - the log sinks choose their level themselves, never from cfg.Observe",

	// [panel]: the section 票 255 was 立 for, and AC#4 landed while this file was
	// open. The host used to hard-code its own geometry - at HEAD 67ab595d the
	// create block inside bringUp spelled it as the literals `Width: 420` /
	// `Height: 260` (a CONTENT anchor, grep it in
	// `git show 67ab595d:cmd/wisp/panel_host_windows.go`; this file quotes no line
	// number for it because that block has moved twice since) - and received no
	// config at all. Both halves of that are gone, and what replaced them is still
	// NOT a reader in this host:
	//
	// cmd/wisp/panel_resident_windows.go's assembly root now hands the host a
	// closure that re-reads [panel] width/height at every window creation, and the
	// host asks it for the number the window is born with. So a value DOES act,
	// and the row says which value acts where rather than claiming a reader this
	// command never had. Three facts the wording has to carry, all measured:
	//   - the reader lives in the resident `wisp` process, which is the only
	//     process that builds a panel host at all (NewPanelManager's one non-test
	//     call site, cmd/wisp/panel_resident_windows.go:273). `wisp run` never
	//     opens a panel window, so on THIS seam nothing acts on the new number -
	//     the same shape [models] is booked in above, which is why this row is
	//     other-process and not consumed.
	//   - "close it and reopen it" (关窗再开) still follows the new value, because
	//     the create re-reads the source. Since 票 255-r1 an ALREADY-created window
	//     is covered too: Show's re-show branch asks the same source once more and
	//     posts the resolved pair through the library's SetSize on Dispatch
	//     (panel_host_windows.go's requestGeometryOnReshow). That hop is a
	//     CLIENT-area request while the create is an OUTER-FRAME one, so the same
	//     width value is not the same pixels on screen - the frame delta is
	//     〔仅本机可量〕 and no ruler here claims the on-screen width equals
	//     [panel] width.
	//   - what still does NOT move is the tick: nothing calls that on the instant
	//     config.toml is saved. internal/config/manager.go's OnReload only fires for
	//     Reload-tier sections while internal/config/tiers.go registers "panel" as
	//     "hot", and startConfigReload's only non-test caller is cmd/wisp/run.go -
	//     the resident leg never starts it. So the honest tense is "the next show
	//     request picks it up", never "改完立刻自己变". Pinned by
	//     TestTicket255PanelRosterVerdictIsTheHonestShape and by
	//     TestTicket255r1ReshowSendsTheResolvedPairThroughDispatch.
	"panel": hotClaimOtherProcess + "cmd/wisp/panel_resident_windows.go:207 [cfg.Panel.Width] - the resident panel host's assembly root re-reads [panel] at every window creation and the window is built from that number (cmd/wisp/panel_host_windows.go:262 [Width:  uint(width)], reached from the create at :392); 面板关窗再开即跟上新值, and since 票 255-r1 a re-show of an already-created window posts this host's currently resolved pair through the library's SetSize on Dispatch - a CLIENT-area request, not the OUTER-FRAME one the create makes, so the same width is not the same on-screen pixels; nothing presses it when config.toml is saved, so the size arrives on the next show request, and this `wisp run` process builds no panel host at all",

	// Per-key sections: planApp/planVoice book the SECTION name into rep.Hot, so
	// the claim must be settled per key row and only claimed when every hot row
	// behind the name is consumed.
	//
	// Ticket 233 re-measured the reader list these rows cite: ticket wire-r1's
	// resident voice leg (cmd/wisp/resident_voice_windows.go) joined the
	// section's production readers at its boot LoadFile - it reads [voice]
	// enabled (:1315) and asr.provider (:1321) to decide whether the voice
	// chain is built at all, and [session] warm_timeout_sec (:1329) for the
	// warm plan. The four hot KEYS below still read from nowhere, so the rows
	// stay no-reader; only the reader-list half of their verdict text changes.
	"app.theme": hotClaimNoReader + "扫描零命中：nothing outside internal/config reads cfg.App.THEME - manager.go's planApp compares and copies it, and no component re-skins from it. The SECTION does have one production reader since 票 43: cmd/wisp/autostart_windows.go reads App.Autostart (the restart-tier key) once per boot to make the registry Run key agree with it - a reader of a different key, so this row stays no-reader",
	// 2026-11-05 re-adjudication (ticket 98's C2): ticket 61's cloud catalog
	// added a production READ SITE for voice.tts.speed, so the old 扫描零命中
	// wording is no longer literally true - but it did not add a reader that
	// anything reaches: internal/speech/cloudcatalog.go copies the knob onto a
	// CloudEndpoint inside ResolveCloudChain, and that function still has zero
	// non-test callers in the repository. The row stays no-reader (nothing that
	// RUNS acts on the value) and now names the dead read site instead of
	// denying it exists; the other three hot [voice] keys have no read site at
	// all.
	"voice.tts.speed": hotClaimNoReader + "没有一条运行路径按这个值做事：唯一的生产读点是 internal/speech/cloudcatalog.go:119 [ep.Speed = cfg.Voice.TTS.Speed]，而它住在票 61 的 ResolveCloudChain 里——那枚函数全仓非测试调用者 0 枚（cmd/internal 里只有定义与注释）⇒ 读点存在、没人调用，TTS 的 speed 仍未接进运行的语音链 (the section's production readers answer other keys: cmd/wisp/resident_audio_windows.go:257 [ra.voiceEnabled = c.Voice.Enabled] decides at boot whether the capture leg is built at all, ticket wire-r1's cmd/wisp/resident_voice_windows.go:1315 [c.Voice.Enabled] makes the same boot decision for the voice leg, and since 票 313 the settings leg cmd/wisp/panel_config_store.go and the panel view internal/panel/config_handlers.go read the cloud-chain/realtime keys off the live Manager - the files sectionReadSites pins)",
	// Ticket 423 re-adjudicated this row: the switch used to ride a hardcoded
	// `true` into speech.OpenPunctuation, so no production line read the key at
	// all (the row carried the 扫描零命中 wording then). The boot LoadFile now
	// carries it into voicePlan.PunctEnabled, which gates both the stage's load
	// (the ASR epoch acquires no punctuation engine when it is false) and the
	// stage's OpenPunctuation enabled argument - whose false contract is the
	// "switched off" state 27-r1 pinned (nil stage, raw finals), reached by
	// configuration instead of by a missing model. The verdict moves to
	// other-process, NOT consumed: the read is one boot-time LoadFile that never
	// re-reads the live Manager (a hot-applied edit waits for the next process,
	// even though the key itself is hot-tier), and the reader lives in the
	// resident `wisp` process while this receipt prints in `wisp run`, which
	// builds no voice leg.
	"voice.punctuation":          hotClaimOtherProcess + "cmd/wisp/resident_voice_windows.go:1850 [c.Voice.Punctuation] - the resident voice leg carries the switch in its boot plan (voicePlan.PunctEnabled) and the punctuation stage opens and attaches under it (cmd/wisp/resident_voice_engines_windows.go:445's OpenPunctuation enabled argument, gated by the ASR-epoch acquire); the section's other production readers answer other keys: cmd/wisp/resident_audio_windows.go:257 [ra.voiceEnabled = c.Voice.Enabled] decides at boot whether the capture leg is built at all, ticket wire-r1's cmd/wisp/resident_voice_windows.go:1315 [c.Voice.Enabled] makes the same boot decision for the voice leg, and since 票 313 the settings leg cmd/wisp/panel_config_store.go and the panel view internal/panel/config_handlers.go read the cloud-chain/realtime keys off the live Manager - the files sectionReadSites pins",
	"voice.wake_word.thresholds": hotClaimNoReader + "扫描零命中：nothing outside internal/config reads this [voice] key (the section's production readers answer other keys: cmd/wisp/resident_audio_windows.go:257 [ra.voiceEnabled = c.Voice.Enabled] decides at boot whether the capture leg is built at all, ticket wire-r1's cmd/wisp/resident_voice_windows.go:1315 [c.Voice.Enabled] makes the same boot decision for the voice leg, and since 票 313 the settings leg cmd/wisp/panel_config_store.go and the panel view internal/panel/config_handlers.go read the cloud-chain/realtime keys off the live Manager - the files sectionReadSites pins) - the wake-word tunables are not consumed by the KWS path yet",
	"voice.wake_word.veto_words": hotClaimNoReader + "扫描零命中：nothing outside internal/config reads this [voice] key (the section's production readers answer other keys: cmd/wisp/resident_audio_windows.go:257 [ra.voiceEnabled = c.Voice.Enabled] decides at boot whether the capture leg is built at all, ticket wire-r1's cmd/wisp/resident_voice_windows.go:1315 [c.Voice.Enabled] makes the same boot decision for the voice leg, and since 票 313 the settings leg cmd/wisp/panel_config_store.go and the panel view internal/panel/config_handlers.go read the cloud-chain/realtime keys off the live Manager - the files sectionReadSites pins) - the veto list is not consumed by the KWS path yet",
}

// sectionReadSites pins, per Config field, which production files the reader scan
// in TestTicket255RosterStillMatchesTheActualReadSites is allowed to hit. It is
// the roster's other half: a NEW file reading a section nobody claimed, or a
// claimed file that stopped reading it, reddens that test instead of quietly
// making the receipt untrue. The values are file paths, not line numbers, so an
// unrelated edit above a reader does not churn the pin.
var sectionReadSites = map[string][]string{
	// Ticket 419 gave this list its first two entries (the style half); ticket
	// 413 added the third: the [ball].size boot seed is a production read site
	// now (see the [ball] roster row above). The verdict itself stays
	// no-reader: all three reads are boot / invocation-time config.LoadFile
	// passes, not live re-readers, and this `wisp run` host builds neither.
	"Ball": {"cmd/wisp/ball_size_413_windows.go", "cmd/wisp/ball_style_65_windows.go", "cmd/wisp/config_cmd_65.go"},
	// Ticket 233: the resident voice leg's boot LoadFile reads WarmTimeoutSec
	// (see the roster row above), so [session] left the reader-less class.
	"Session": {"cmd/wisp/resident_voice_windows.go"},
	// Ticket 247's capture assembly is the first production reader of both
	// sections: it reads [audio] mic_muted_default into audio.WithStartMuted and
	// [voice] enabled as "build the collector at all, or not". Both reads are
	// one fresh config.LoadFile per boot in the resident process (the same
	// per-boot shape [hotkey] has used since ticket 258), so neither lands in
	// the "consumed by a live re-reader" class - see the [audio] row above.
	"Audio":   {"cmd/wisp/resident_audio_windows.go"},
	"Privacy": nil,
	"Memory":  nil,
	// 2026-11-05 re-adjudication (ticket 98's C2): ticket 44 landed the C23
	// money path, and both of its adapters read the four [cost] leaves off a
	// *config.Config they are handed - agent.CostWiring (costbook_44.go, called
	// from execute() with the LIVE Manager config, cmd/wisp/run.go:1570) and
	// agent.QuotaWiring (quota_44r4.go, reached through quotaDispatchChain).
	// The same two files also read [llm] provider cards, which is why they
	// appear in both lists.
	"Cost":    {"internal/agent/costbook_44.go", "internal/agent/quota_44r4.go"},
	"Observe": nil,
	// Read by the resident panel host's assembly root, per window creation, since
	// 票 255 AC#4. The host file itself is NOT in this list and must stay out of
	// it: it receives a function value, and cmd/wisp/panel_host_windows.go
	// importing internal/config or naming config.toml would be the dependency edge
	// AC#4 forbids (TestTicket255HostStillDoesNotParseConfigItself checks both).
	"Panel": {"cmd/wisp/panel_resident_windows.go"},
	"App":   {"cmd/wisp/autostart_windows.go"},
	// 票 43: the resident boot's autostart sync is [app]'s first production
	// reader outside internal/config - it reads App.Autostart (restart tier,
	// ticket 233) once per boot and makes HKCU's Run key agree with it. The
	// restart tier is why the read is boot-time only: the key's own verdict
	// (需要重启进程生效) is what makes this the honest consumer shape.
	// 票 313 added two more [voice] readers to the scan's list: the settings leg
	// cmd/wisp/panel_config_store.go reads the cloud chains and realtime keys off
	// the LIVE Manager (ReadSettings), and internal/panel/config_handlers.go
	// renders that answer - its `v.Voice.*` lines read panel.SettingsView, so the
	// scan keys on the field name and this list records what it measures. Ticket
	// wire-r1's voice leg added the fourth: it reads [voice] enabled/asr.provider
	// and [session] warm_timeout_sec at its boot LoadFile (see the per-key rows
	// above). The per-key rows stay "no-reader" because no file reads THOSE keys -
	// a verdict about the keys, not about these files.
	// 2026-11-05 re-adjudication (ticket 98's C2): 票 61's cloud catalog is the
	// fifth file the scan names - ResolveCloudChain reads voice.cloud_*_chain
	// plus tts.voice / tts.speed off its snapshot argument. It is a READ SITE
	// whose function still has zero non-test callers, which is exactly why
	// voice.tts.speed's row was re-worded (the read exists, the caller does
	// not) rather than promoted to consumed.
	"Voice":  {"cmd/wisp/panel_config_store.go", "cmd/wisp/resident_audio_windows.go", "cmd/wisp/resident_voice_windows.go", "internal/panel/config_handlers.go", "internal/speech/cloudcatalog.go"},
	"Agent":  {"cmd/wisp/run.go"},
	"Models": {"cmd/wisp/models.go"},
	// 票 320's fetch port added a fifth live reader of [llm], and 票 324 AC#2 books it
	// here: cmd/wisp/panel_models_320.go's HandleModelsListRequest takes `cfg :=
	// h.mgr.Config()` on the request hop and reads `cfg.LLM.Providers[provider]` off
	// that answer, so it reads the LIVE Manager at use time - the same class as
	// panel_config_store.go, which is why the row above stays in the hotClaimConsumed
	// tier instead of moving to snapshot-only or other-process. The tier was not picked
	// for this edit: it was adjudicated by 票 255 AC#1 and this reader agrees with it
	// (a hand-edited provider row changes what the NEXT models.list hop asks for, with
	// no restart). Its caveat is the row's existing caveat, unchanged - the boot
	// snapshot still drives the in-flight chain.
	// 2026-11-05 re-adjudication (ticket 98's C2): three more files read this
	// section off a config value they are handed, and are pinned here rather
	// than cited in the row: 票 44's internal/agent/costbook_44.go (price card
	// for one provider/model) and internal/agent/quota_44r4.go (provider cards,
	// model quota cards, quota params) - both fed the LIVE per-task cfg by
	// execute() - and 票 61's internal/speech/cloudcatalog.go (the provider
	// entry a cloud chain element names). The verdict does not move: the claim
	// still rests on the live re-readers already cited.
	"LLM": {"cmd/wisp/panel_config_store.go", "cmd/wisp/panel_models_320.go", "cmd/wisp/providers.go", "cmd/wisp/run.go", "internal/agent/costbook_44.go", "internal/agent/quota_44r4.go", "internal/llm/resolver.go", "internal/speech/cloudcatalog.go"},
	// Ticket 258 moved the shipped reader into the resident leg: the
	// construction chain lives in resident_windows.go (it maps [hotkey] into
	// ball.HotkeyConfig), while resident_ball_windows.go only receives the
	// mapped value and installs the reloader bridge. balldebug keeps reading it
	// too, so it stays.
	"Hotkey": {"cmd/wisp/resident_windows.go", "cmd/balldebug/main.go"},
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

// bracketed renders a hot-tier name list as `[panel] [llm]`, one bracket per
// section, instead of `%v`'s `[panel llm]`.
//
// WHY, since the old line used `%v`: the receipt's claim is checked per section
// (票 255 AC#1's judgement is literally "那句话里不许出现 `[panel]`"). With `%v` a
// multi-element list reads `[ball session panel]`, where a section is present but
// un-bracketed - so a reader or a grep looking for `[panel]` misses it, which is
// how a lying claim survives a spot check. The restart sentence already names
// sections one bracket at a time (「[fs] 的放宽…」), so this is the house shape too.
func bracketed(names []string) string {
	if len(names) == 0 {
		return "[]"
	}
	parts := make([]string, 0, len(names))
	for _, n := range names {
		parts = append(parts, "["+n+"]")
	}
	return strings.Join(parts, " ")
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
