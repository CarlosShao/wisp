package main

// Ticket 223 - the hot-reload wiring: the tick that re-reads config.toml, the
// D36 L2 re-confirmation it raises for a loosening, and the loud sentence for a
// change that has to wait for a restart.
//
// WHY THIS FILE EXISTS. Three of config.Manager's hooks had no production
// assignment point before this file (measured at HEAD 7ffa9520, census
// .scratch/wisp/probes/223/c1/census.md):
//
//   - Manager.CheckAndReload - its only non-test caller was cmd/balldebug, a
//     side program. `wisp run` never re-read config.toml, so D36's three-tier
//     semantics (立即/重载/重启, PLAN.md:2721) had an engine and no driver.
//     This file added the driver: reloadOnce() below calls CheckAndReload.
//   - Manager.ConfirmLocked - before this file, zero assignments outside
//     *_test.go; manager.go snapshots the hook under mu (confirm := m.ConfirmLocked)
//     and a nil hook denies (approved := confirm != nil && ...), so an unwired
//     hook was a silent fail-closed DENY of every [fs] loosening: D33's
//     「热加载放宽必须触发 L2 级重新确认」 was attached to a dead wire.
//     startConfigReload assigns it below, before the tick is spawned.
//   - Manager.OnRestartPending - zero assignments outside *_test.go too before
//     this file, so the restart tier said nothing at all: a hand edit of [app]
//     simply did not happen, silently, which is the one shape 票 223 AC#2 forbids
//     (不许用「静默不生效」充当这一档). startConfigReload assigns it below,
//     before the tick is spawned.
//
// TRIGGER SHAPE, NAMED (AC#1). A resident tick - not a file notification, and
// not only an explicit command: one goroutine spawned through
// observe.Default.Spawn under the D38b roster name "watchdog", the slot SPEC-03
// §4.3 already reserves for this very tick ("config-file mtime polling reusing
// its loop", internal/watchdog/doc.go:3) and which nothing spawned before this
// file. It wakes on a time.Ticker at 1s and calls CheckAndReload, which is the
// documented idempotent poll: unchanged mtime+size returns (nil, nil) without
// re-running the load pipeline. Building the real watchdog loop (SLO sampling,
// thresholds, the D32 table, the DEFERRED marker internal/watchdog/doc.go:18
// still carries) stays ticket 42's job; this file takes only the config-poll
// half and says so out loud.
//
// WHY THIS IS NOT A WALL-CLOCK TIMEOUT (AGENTS.md §1.2's banned shape). The tick
// measures nothing. There is no elapsed-time comparison anywhere in this file: no
// `.Sub(time.Now())`, no `time.Now().Unix(`, not even `time.Since` - the ticker's
// channel is the only scheduler, and every decision on this path is either "the
// file's mtime+size differs from what this Manager adopted" (a content fingerprint
// per reload, not a duration) or "a human answered this card" (a channel result).
// The one deadline that does exist - how long an unanswered L2 card may sit - is
// not this file's either: it is the C18 gate's ApprovalTimeout, which
// internal/agent/approval runs on its own clock. Nothing here can even express
// "too late", which is why an unanswered card lands as a deny rather than as an
// silently expired grant.
//
// WHAT A CONFIRMATION COSTS, AND WHY IT IS SAFE NOW. ConfirmLocked is reached
// from inside Manager.CheckAndReload. Until this ticket that call held the
// Manager's mutex, so a hook that read its own config deadlocked and an
// unanswered card froze every Config() reader for up to 300s; manager.go's
// plan/commit split (same commit as this file) moved the ask outside the lock.
// That is what lets the hook below do the two things a real confirmation needs:
// call rt.permissionMode(), which reads Manager.Config(), and park on a human.
//
// WHAT THIS FILE DOES NOT FIX (named, not hidden). The C26 canonicalizer is built
// once, at assembly (cmd/wisp/run.go:390), and both the bridge and the registered
// fs tools hold that pointer; swapping rt.paths mid-run would need an exported
// reset on internal/tools, which this leg may not mint (no new exported names).
// So an APPROVED [fs] loosening changes this process's memory and its next
// start's judgements, but not this run's path verdicts - and the line below says
// exactly that, instead of letting 「已确认」 read as 「已生效」.

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/CarlosShao/wisp/internal/config"
	"github.com/CarlosShao/wisp/internal/observe"
	"github.com/CarlosShao/wisp/internal/risk"
	"github.com/CarlosShao/wisp/internal/tools"
)

// The card identity for a reload-triggered confirmation, in the shape
// confirmModeSwitch (run.go) and the ticket 201 widening card already use: a
// synthetic task id the gate has never seen because no model call raised it,
// admitted for exactly the length of the confirmation, plus a
// 'namespace.action' name that names the door being widened rather than a tool
// the model can call. Neither is a C4 registration and neither is a C17 method.
const (
	configReloadTaskID = "host:config-reload"
	configReloadTool   = "config.reload"
	// configReloadPollInterval is SPEC-03 §4.3's 1s watchdog tick.
	configReloadPollInterval = time.Second
)

// hotReloadDisabledPanelInbound is AC#4's fourth sentence, and it is a real
// production line, not a test-only branch: `wisp panel-inbound` opens the same
// config.toml through config.NewManager and owns no tick, because it has no
// approval gate to raise a card from (its perm.Options.Confirm is nil on
// purpose - see that file's header). A host that cannot ask must not silently
// widen, so this host says plainly that nothing it holds will pick up an edit.
const hotReloadDisabledPanelInbound = "config: HOT-RELOAD state=disabled host=panel-inbound " +
	"detail=\"本宿主没有接 config.toml 轮询，这次运行期间手改配置不会生效；" +
	"它也没有可弹卡的面，所以放宽本来就无法被确认（fail-closed）。要生效请重启进程并用 wisp run。\""

// startConfigReload wires this run to D36's three tiers and starts its tick.
// Called from the assembly root, after the gate exists: a confirmation needs a
// gate to be raised on, and wiring the hook before rt.gate exists would hand
// Manager a hook that can only ever answer "no card surface available".
func (rt *agentRuntime) startConfigReload() {
	if rt == nil || rt.mgr == nil {
		return
	}
	rt.reloadRoot = observe.NewRoot("config-reload")
	// The two hooks ticket 223 found with zero production assignment points.
	// (Mutation readings m1/m2 were taken against commented-out copies of these
	// three lines at 15:16:52 and 15:19:20; docs/evidence/s1/223-hot-reload-wiring-r1.md
	// carries the red readings, and both lines are back.)
	rt.mgr.ConfirmLocked = rt.confirmLockedLoosening
	rt.mgr.OnRestartPending = rt.reportRestartPending
	// AC#1's production caller. Owner names the subsystem that spawned it; the
	// recover boundary and the roster entry come from Registry.Spawn, never from
	// a local recover (D22 ban #1: bare `go` has no legal annotation).
	rt.reloadHandle = observe.Default.Spawn("watchdog", "config", rt.reloadRoot, rt.configReloadTick)
	rt.auditf("config: HOT-RELOAD state=armed tick=%s path=%q goroutine=watchdog owner=config "+
		"d36_confirm=on restart_notice=on detail=%q",
		configReloadPollInterval, filepath.Join(rt.spec.dataDir, configFileName),
		"本进程会每 tick 重读 config.toml：手改的可热加载段立即生效，锁定段的放宽要先过一张 L2 卡，重启档会明确告知不生效")
	fmt.Fprintf(rt.stdout,
		"wisp run: 配置热加载已接管（每 %s 检查一次 config.toml）。手改会按 D36 三档处理："+
			"可热加载段立即生效；[risk]/[fs]/[net]/[plugins] 的放宽要先答一张 L2 卡，不答按拒绝保留旧值；"+
			"重启档的改动本次不生效，会另有一句告诉你为什么不生效。\n", configReloadPollInterval)
}

// configReloadTick is the 1s poll body. It ends when the run's reload root ends.
func (rt *agentRuntime) configReloadTick(ctx context.Context) {
	t := time.NewTicker(configReloadPollInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			// Distinct from "disabled" on purpose: this host DID arm a tick and
			// it is going away with the process, so an edit after this point is
			// not being ignored by design - it is arriving after the run ended.
			rt.auditf("config: HOT-RELOAD state=stopped reason=run-root-cancelled detail=%q",
				"轮询协程已退出（本次进程正在收口）；此后对 config.toml 的改动要下一次启动才会被读到")
			return
		case <-t.C:
			rt.reloadOnce()
		}
	}
}

// reloadOnce runs one idempotent poll and turns its outcome into sentences. A
// nil report means "nothing changed" and is silent by contract (manager.go:62);
// every other branch says something.
func (rt *agentRuntime) reloadOnce() {
	rep, err := rt.mgr.CheckAndReload()
	if err != nil {
		rt.auditf("config: HOT-RELOAD state=not-applied %s", describeReloadFailure(err))
		return
	}
	if rep == nil {
		return
	}
	rt.reportReload(rep)
}

// reportReload books the tier verdicts this reload produced. One line per
// locked section, so "did my edit take effect" is a reading off the ledger and
// never an inference; one user-visible line for the tiers an operator can act on.
//
// Ticket 255 AC#1 narrowed what the hot line is allowed to say. rep.Hot is booked
// by internal/config/manager.go's plan() off "the values differed and the memory
// got overwritten" - it never asks whether anything reads the section. Printing it
// verbatim claimed 「已立即生效」 for [panel] (zero production readers, and
// panel_host_windows.go hard-codes its own geometry), which is the sentence this
// ticket was立 for. The hot list is now filtered by cmd/wisp/config_readers_255.go's
// roster, resolved through config.TierOf rather than a copied word table: sections
// whose hot rows all name a live reader keep this line, the rest get their own
// honest sentence. Neither line is ever dropped, because taking the ruler down is
// not fixing it.
func (rt *agentRuntime) reportReload(rep *config.Report) {
	rt.auditf("config: HOT-RELOAD state=applied hot=%v reload=%v restart=%v locked=%d",
		rep.Hot, rep.Reload, rep.Restart, len(rep.Locked))
	if len(rep.Hot) > 0 {
		split := splitHotTier(rep.Hot)
		for _, name := range rep.Hot {
			rt.auditf("config: HOT-RELOAD-READER section=%s %s", name, split.verdictFor(name))
		}
		if len(split.claimable) > 0 {
			fmt.Fprintf(rt.stdout, "wisp run: 配置热加载：这些段已立即生效（D36 立即档）：%s\n", bracketed(split.claimable))
		}
		if len(split.quiet) > 0 {
			fmt.Fprintf(rt.stdout,
				"wisp run: 配置热加载：这些段的值已换进本进程内存，但本宿主没有会按新值做事的读者，"+
					"本次运行不会因此改变行为（票 255 AC#1：这一半不许说成「已立即生效」；"+
					"逐段的读者判定见 HOT-RELOAD-READER 行）：%s\n", bracketed(split.quiet))
		}
		for _, name := range split.disagree {
			// plan() hot-applied a name the tier registry cannot explain. The
			// manager.go guard makes this unreachable for whole sections; it is
			// loud anyway so a future per-key split cannot land as silence.
			rt.auditf("config: HOT-RELOAD state=tier-disagreement section=%s %s",
				name, split.verdictFor(name))
			fmt.Fprintf(rt.stdout,
				"wisp run: 配置热加载：[%s] 被热应用了，但 config.TierRegistry 里查不到它的 hot 行，"+
					"回执没有资格判断它是否真有人读（票 255 AC#2-ⓑ：登记表与实现不符，请上报）。\n", name)
		}
	}
	for _, d := range rep.Locked {
		effect := "applied"
		switch {
		case d.Direction == config.DirLoosen && !d.Approved:
			effect = "kept-old-values"
		case d.Direction == config.DirLoosen && d.Approved:
			effect = "applied-after-L2"
		}
		rt.auditf("config: D36-SECTION section=%s direction=%s approved=%v effect=%s keys=%v",
			d.Section, d.Direction, d.Approved, effect, d.Keys)
		if d.Direction != config.DirLoosen {
			continue
		}
		if !d.Approved {
			rt.auditf("config: HOT-RELOAD state=denied section=%s keys=%v detail=%q",
				d.Section, d.Keys,
				"放宽没有被确认（卡片被拒绝、被否决或没人答），内存里继续用旧值；配置文件里的那一行没有被丢掉，也没有被偷偷生效")
			fmt.Fprintf(rt.stdout,
				"wisp run: [%s] 的放宽本次没有生效（原因：那张 L2 重新确认没有拿到「允许」）。"+
					"旧的严格值仍在用，你写的这一行还在文件里：%v\n", d.Section, d.Keys)
			continue
		}
		fmt.Fprintf(rt.stdout,
			"wisp run: [%s] 的放宽已经过 L2 重新确认并写进本进程内存（D33：放宽不得静默生效）：%v\n",
			d.Section, d.Keys)
		if d.Section == "fs" {
			// The honesty line this file owes: memory moved, this run's path
			// verdicts did not. Say it instead of letting the operator infer.
			fmt.Fprintf(rt.stdout,
				"wisp run: 注意 - 本次运行的路径判定仍按启动时建好的 C26 名单，"+
					"新放宽的目录要重启进程才参与判定（原因：canonicalizer 只在装配时构造一次）。\n")
			rt.auditf("config: D36-SECTION section=fs effect=memory-only-this-run "+
				"detail=%q", "本次运行的 C26 canonicalizer 是启动时建的，放宽只对下一次启动生效")
		}
	}
}

// confirmLockedLoosening is D36 rule 1's production L2 re-confirmation: the one
// card a hand-edited loosening of a locked section must cost. It runs with the
// Manager's lock released (manager.go's plan/commit split), which is what lets it
// both read the live config and wait for a human - the two things that were
// impossible while the hook was called under the mutex.
//
// Every way out of here except an explicit allow is a deny, and a deny keeps the
// previous values: no answer, no gate, no root, an abandoned card, a refused
// card. That direction is D33's, and it is the reason an unwired hook (the state
// ticket 223 found) was already fail-closed - the bug was that it was fail-closed
// silently, with no card and no sentence.
func (rt *agentRuntime) confirmLockedLoosening(section string, keys []string) bool {
	if rt == nil {
		// No owner, so no card and no verdict: deny, which is the direction D33
		// requires and never a silent widening.
		return false
	}
	if rt.gate == nil {
		rt.auditf("config: D36-CONFIRM state=no-gate section=%s keys=%v result=deny "+
			"detail=%q", section, keys,
			"这台宿主没有审批 gate，放宽无法被确认；保留旧值（不是静默生效）")
		return false
	}
	if rt.reloadRoot == nil {
		rt.auditf("config: D36-CONFIRM state=no-root section=%s keys=%v result=deny "+
			"detail=%q", section, keys,
			"热加载轮询没有拥有者根，卡片无处挂超时与收口；保留旧值")
		return false
	}
	revoke := rt.gate.AdmitTextTask(configReloadTaskID)
	defer revoke()
	// The mode this run is screening under, read through Manager.Config() while
	// the card below is still unanswered. It is only safe because ticket 223
	// moved the ask off the Manager's mutex; on the old shape this line
	// deadlocked the whole config layer.
	mode := rt.permissionMode()
	ans, why := rt.gate.PendingApproval(rt.reloadRoot.Ctx, tools.Decision{
		TaskID: configReloadTaskID,
		Tool:   configReloadTool,
		Level:  risk.L2,
		Reason: fmt.Sprintf("config.toml 里手改的这一条会放宽 [%s]：%v。"+
			"D33/D36 要求放宽必须重新确认一次才生效；拒绝或不做声都保留旧的严格值，"+
			"本次运行不会因为没答而偷偷变宽。（这张卡由配置轮询签发，与任何一次工具调用无关）",
			section, keys),
		Mode: mode,
	})
	if ans != tools.AnswerAllow {
		if why == "" {
			why = string(ans)
		}
		rt.auditf("config: D36-CONFIRM state=answered section=%s keys=%v answer=%s result=deny detail=%q",
			section, keys, ans, why)
		return false
	}
	rt.auditf("config: D36-CONFIRM state=answered section=%s keys=%v answer=%s result=allow",
		section, keys, ans)
	return true
}

// reportRestartPending is AC#7's exit: the restart tier finally has somewhere to
// say what it is. Before this function existed, rep.Restart's contents were read
// by nobody and an [app] edit simply did nothing, silently - which is the shape
// 票 223's second tier forbids.
//
// It says the three things the operator needs and cannot derive from silence:
// which sections, that this run will not use the new values, and WHY (the keys
// hand their value to the platform layer once, at start-up). It also says what
// did NOT happen to the file, because "不生效" and "被丢了" are different fears.
func (rt *agentRuntime) reportRestartPending(sections []string) {
	if rt == nil {
		return
	}
	rt.auditf("config: HOT-RELOAD state=restart-pending sections=%v effect=next-process-start tier=restart",
		sections)
	rt.auditf("config: RESTART-PENDING detail=%q keys=%q",
		"这些键只在进程启动时被读一次并交给平台层（开机自启注册、单实例锁、界面语言），本次运行没有重新加载它们的路径；"+
			"内存里继续用旧值，文件里的新值没有被丢掉，重启进程后生效",
		restartTierKeys)
	fmt.Fprintf(rt.stdout,
		"wisp run: 这些段的改动本次运行不会生效（D36 重启档，需要重启进程）：%v。"+
			"原因：这几枚键在进程启动时被读一次就交给平台层（开机自启注册、单实例锁、界面语言），"+
			"本次运行没有重新注册它们的路径。涉及：%s。"+
			"注意两件事都没发生：文件里的新值没有被丢掉，内存里的旧值也没有被偷偷替换。\n",
		sections, strings.Join(restartTierKeys, " / "))
}

// restartTierKeys names the keys planApp books as restart-tier (D36, PLAN.md
// :2715's [app] split). Kept as data so the sentence above and manager.go's
// split can be compared by a reader - and by a test - instead of drifting.
// That "and by a test" was a phantom claim until ticket 255 (ledger A549 ruled
// this cleanup onto 255-r2): the reader was zero outside this file. It is now
// TestTicket255RestartTierKeysAreBackedByATest in restart_tier_keys_255r2_test.go,
// which drives each named key through the real Manager and requires it to land in
// Report.Restart and never in Report.Hot, pins the list name-for-name, and requires
// every named key to be registered "restart" in config.TierRegistry.
var restartTierKeys = []string{
	"app.language", "app.autostart", "app.single_instance",
}

// describeReloadFailure turns "the reload did not happen" into one sentence per
// cause. 票 223 AC#4 forbids collapsing these: a missing file, a file with a
// syntax error, a file this process may not read and a host that never armed a
// tick are four different things an operator can fix in four different ways, and
// one shared "配置未生效" line would be a lie by omission (the same discipline
// 票 216 set). Each branch is exercised by a planted artifact in
// config_reload_223_test.go, and each case asserts the OTHER sentences are
// absent, so the four stay four.
//
// Order carries weight: the io/fs sentinels are checked first, because a failed
// read also arrives wrapped in a config-class error, and "the file is gone" must
// not be reported as "the file is invalid".
func describeReloadFailure(err error) string {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return "cause=missing detail=\"" +
			"config.toml 读不到：文件不存在（这一条只说缺失，不说语法、不说权限）。本次运行继续用内存里的旧配置；" +
			"文件回来之后要它的 mtime 或大小变过才会被下一 tick 重读\""
	case errors.Is(err, fs.ErrPermission):
		return "cause=permission detail=\"" +
			"config.toml 读不到：这个进程没有读它的权限（文件在，也读得开名字，只是不让读）。" +
			"本次运行继续用内存里的旧配置；这一条不是语法错，也不是文件缺失\""
	}
	var oe *observe.Error
	if errors.As(err, &oe) {
		d := oe.Detail
		switch {
		case strings.HasPrefix(d, "config.toml parse"):
			return "cause=syntax detail=\"" +
				"config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。" +
				"本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读\""
		case strings.Contains(d, "unknown key"):
			return "cause=unknown-key detail=\"" +
				"config.toml 语法没问题，但里面有这份 schema 不认的键（拼错的键会被这样拒绝，而不是被忽略）。" +
				"本次运行继续用内存里的旧配置\""
		case strings.Contains(d, "cannot migrate from schema version") ||
			strings.Contains(d, "no migration registered"):
			// Its own sentence because it has its own fix: this file declares a
			// version this build cannot carry forward, and the config layer's
			// message already promises the file was left untouched. Measured
			// boundary (票 223 r1, narrowed by 票 223 r2): a file declaring a
			// version BELOW this build's current one and broken further down
			// lands HERE - loader.go routes it through the migration pipeline,
			// which is what migrate_test.go:123 requires. r1 had drawn the
			// line at "any readable declared version"; 223-v1's 17-shape
			// overlay showed that made four declared-CURRENT shapes (whose
			// only fault was a broken line) borrow this branch's following
			// 语法没问题 sibling instead of 语法错. Since r2 a file declaring
			// the current-or-newer version with an unparseable body is named
			// a syntax error; this branch keeps exactly the shape its own
			// fix applies to.
			return "cause=migration detail=\"" +
				"config.toml 声明了一个这份 Wisp 不会迁移的 schema_version（文件被原样留着，不会被重置）。" +
				"本次运行继续用内存里的旧配置；升级 Wisp 或恢复备份才会读它\""
		case strings.Contains(d, "was written by a newer build"):
			// Ticket 231 AC#2: a file declaring a version ABOVE this build's is
			// its own shape, and until this branch it had no sentence. It cannot
			// fall through to the two branches below: the prefix branch books it as
			// "值不合法或引用解不开", which is false for it -
			// loader.go:120 returns before decodeStrict and before validate, so the
			// file was never validated at all - and the migration branch's wording
			// ("不会迁移的 schema_version") describes a version BELOW this build,
			// which is the one shape that does have a migration to run.
			//
			// Position is load-bearing: this case must stay written BEFORE the
			// strings.HasPrefix(d, "config.toml:") branch, because loader.go:122
			// opens its detail with exactly that prefix. A copy placed after it is
			// dead code and the new row in config_sentences_223r2_test.go reddens
			// on its prefix ruler with "the two sentences swapped again".
			//
			// It claims a phrase, not a prefix: "was written by a newer build" is
			// loader.go:122's own wording (the cause token comes from the same
			// source), so the branch is narrower than the prefix it is racing and
			// cannot reach any of the other 47 details in internal/config that
			// also open with "config.toml: ".
			//
			// AC#4's registration (this branch still classifies prose; why today
			// nothing else moves) is in .scratch/wisp/probes/231/r1/evidence.md §4.
			return "cause=newer-build detail=\"" +
				"config.toml 是由一个更新的 Wisp 写出来的（它声明的 schema_version 比这份程序懂得的高；" +
				"这一条不说语法错，也不说值不合法，因为它还没走到校验）。本次运行继续用内存里的旧配置；" +
				"升级 Wisp 或恢复备份才会读它\""
		case strings.HasPrefix(d, "config.toml:"):
			return "cause=invalid detail=\"" +
				"config.toml 语法没问题，但内容被校验拒绝（值不合法或引用解不开）。" +
				"本次运行继续用内存里的旧配置\""
		}
	}
	// The stat failure that is neither missing nor permission, a migration
	// failure, and anything a later loader adds all land here. It quotes the
	// cause instead of pretending to have classified it - the alternative is a
	// line that reads like an answer.
	return fmt.Sprintf("cause=unclassified detail=%q",
		"config.toml 没有被重读成功，原因没有归入已知四类（缺失/语法/权限/schema）；本次运行继续用内存里的旧配置: "+
			err.Error())
}
