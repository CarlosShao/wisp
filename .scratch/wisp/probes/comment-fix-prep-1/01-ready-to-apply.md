# 01-ready-to-apply — 13 枚〔只欠文案〕逐处可直接粘贴的替换文本

- **腿**：`comment-fix-prep-1`（准备料腿）。⛔ 本料**不改任何产码**；落地由后续窄腿执行。
- **现量基点**：HEAD `9e8477edc690c46f3540f97069a1009c05ea2fc4`（2026-10-08 18:27 +0800）。全部锚行用 `git show HEAD:<file> | grep -n '' | sed -n 'N,Mp'` 现取；调用形状尺＝`git grep -nE '<形状>' HEAD -- '*.go' ':(exclude)*_test.go' ':(exclude).scratch'`；**零 Go 命令**。
- **料源**：`.scratch/wisp/probes/comment-truth-2/01-triage.md`（入库 `762b694e`）。归口照该表原样，本料不重裁。
- **范围**：13 枚〔只欠文案〕＝P01 P02 P03 P04 P06 P07 P08 P09 P10 P11 P12 P13 P39。⛔ P14／P15（〔要动产码〕）**未碰**。
- **13 枚锚行号复测**：与 01-triage 的 15 枚锚**逐枚一致、零漂移**（P07 也在 `:11`）。
- **与 01-triage 的口径差（具名报回，以盘上为准）**：
  1. **P02**：triage 只点"第 1 条 bullet 已在盘上"；现量 HEAD 下 **bullet 2（AdmitTextTask）与 bullet 4（D31 applied-steps 缝）也已落产码**（`run.go:842/:955`、`config_reload.go:272`、`resident_approval_windows.go:666`；`run.go:749 Cancel: rt.gate.ToolsCancelBus()`，`fs_write.go` 多处填 `AppliedSteps`）。两程间隔 22 个 commit，按现量写。
  2. **P03/P04**：补量到另一处 `Handle` 产码调用点＝`cmd/wisp/panel_host_windows.go:821`（C27 宿主 `dispatchRaw`，:803 由 wispDispatch 门喂入）。triage 的 P04 凭据没列它，不影响档位。
- **同文件多处提醒**：`cmd/wisp/config_reload.go` 的 P07/P08/P09 三块，**从后往前贴（P09→P08→P07）**，否则锚共漂。
- **在飞登记（影响落地前复核）**：`cmd/wisp/panel_host_windows.go` 工作树有改动（起手 `git status --porcelain` 的 ` M`）。P01/P03/P04/P39 的替换文本点到它，落地时若该文件已被改动，行号/措辞按当时重核。

---

## P01 — internal/panel/composer_dispatch.go

1. **目标**：`internal/panel/composer_dispatch.go:50`（现量＝`// belongs to the host; and it has NO production caller yet. The honest state of`；替换块＝**47–52**）。
2. **原文逐字**（块 47–52，整行不截断）：
```go
// WHAT THIS DOES NOT DO, stated so nobody infers it from the file's existence:
// it opens no window and imports no WebView2 symbol (slice A is the half of the
// hop that needs neither); it writes no reply back to the page, which is H10 and
// belongs to the host; and it has NO production caller yet. The honest state of
// that question is recorded in docs/evidence/s1/33-minimal-inbound-hop-r1.md
// rather than in a comment that reads better than it is.
```
3. **替换文本**（换 47–52；6 行 → 8 行，可直接粘贴）：
```go
// WHAT THIS DOES NOT DO, stated so nobody infers it from the file's existence:
// it opens no window and imports no WebView2 symbol (slice A is the half of the
// hop that needs neither); it writes no reply back to the page, which is H10 and
// belongs to the host. Slice A shipped with zero production callers; both callers
// exist today - cmd/wisp/panel_host_windows.go's dispatchRaw and
// cmd/wisp/panel_inbound.go (`wisp panel-inbound`, one envelope per stdin line) -
// and the slice-A reading, with its grep lines, is in
// docs/evidence/s1/33-minimal-inbound-hop-r1.md §④.
```
4. **风险 / 现验**：
   - 受扫文件：`composer_dispatch_test.go:375` 的 `"panel\.[^"]*"` 字面量尺（**只扫非注释行**）＋ `:708-716` 整份文件**含注释**扫五枚禁词尺。替换文本为注释行、不含五枚禁词（见末尾体检）。
   - 点名对象现验：`docs/evidence/s1/33-minimal-inbound-hop-r1.md` 存在（147 行），其 §④（`:55`）逐字记着 slice A 零调用者的读数。两处调用点现量：`panel_host_windows.go:821`、`panel_inbound.go:163`。

## P02 — internal/agent/approval/doc.go

1. **目标**：`internal/agent/approval/doc.go:26`（现量＝`// Wiring still owed by ticket 12 (this package has NO production caller yet,`；替换块＝**26–37**）。
2. **原文逐字**（块 26–37）：
```go
// Wiring still owed by ticket 12 (this package has NO production caller yet,
// so nothing here runs in cmd/wisp today):
//   - compose approval.New(...) as tools.Options.Gate (replacing NoGate);
//   - call Gate.AdmitTextTask(taskID) once per TEXT-loop task and defer its
//     revoke - unadmitted tasks are refused by design (D47), so forgetting it
//     fails closed, it does not fall open;
//   - hand Queue.Native() to the ball click / native card / hotkey handlers
//     and Queue.Panel() to the ticket 37 panel bridge;
//   - have the fs.write family (ticket 20 segment 2) poll Gate.Bus() and fill
//     Result.AppliedSteps;
//   - delete loop.decideRisk's declared-L1/L2 refusal, which currently kills
//     fs.write before the bridge ever assesses it.
```
3. **替换文本**（换 26–37；12 行 → 11 行）：
```go
// Wiring owed by the ticket 12 list. Landed since this paragraph was written:
// approval.New is composed in cmd/wisp (run.go:612, resident_approval_windows.go:369)
// and handed to the tool bridge as tools.Options.Gate (run.go:748); Gate.AdmitTextTask
// has production call sites (run.go's mode-switch card and admitTask, the config-reload
// tick, the resident card path); and the D31 ledger seam is wired (tools.Options.Cancel
// takes Gate.ToolsCancelBus, and the fs.write family fills Result.AppliedSteps).
// Still owed:
//   - hand Queue.Native() to the ball click / native card / hotkey handlers
//     and Queue.Panel() to the ticket 37 panel bridge;
//   - delete loop.decideRisk's declared-L1/L2 refusal, which currently kills
//     fs.write before the bridge ever assesses it.
```
4. **风险 / 现验**：
   - 现验：`approval.New(` 产码两处＝`run.go:612`、`resident_approval_windows.go:369`；`Gate: rt.gate`＝`run.go:748`（包在 `tools.New(tools.Options{...})` 内）；`AdmitTextTask` 产码＝`run.go:842`（模式切卡）、`run.go:955`（admitTask）、`config_reload.go:272`、`resident_approval_windows.go:666`（另有 `approval_always.go:113`）；`run.go:749` `Cancel: rt.gate.ToolsCancelBus()`；`fs_write.go` 多行填 `AppliedSteps`。
   - "Still owed"两条独立现验：`.Native()`/`.Panel()` 产码调用只在 `internal/agent/approval/replies.go` 包内，无任何 ball click／native card／hotkey 侧接点 ⇒ 仍欠；`loop.go:793-797` decideRisk 的 L1/L2 拒绝（`AdmitTask == nil` 分支）仍在 ⇒ 仍欠。
   - 原文 13 条逐字串在 `*_test.go` 0 命中（triage 已测，本程复测见末尾）；无同族仪器扫 `doc.go`。

## P03 — cmd/wisp/panel_inbound.go

1. **目标**：`cmd/wisp/panel_inbound.go:11`（现量＝`// line and hands the bytes to (*panel.ComposerDispatch).Handle. Nothing else in`；替换块＝**11–14**）。
2. **原文逐字**（块 11–14）：
```go
// line and hands the bytes to (*panel.ComposerDispatch).Handle. Nothing else in
// this repository calls Handle today, and the judgement in
// panel_inbound_33_test.go reads THAT fact off the disk rather than off this
// paragraph, so deleting the line below reddens a case instead of passing quietly.
```
3. **替换文本**（换 11–14；4 行 → 5 行）：
```go
// line and hands the bytes to (*panel.ComposerDispatch).Handle. A second
// production caller now exists - cmd/wisp/panel_host_windows.go's dispatchRaw,
// the C27 host's dispatch door - and the judgement in panel_inbound_33_test.go
// reads the cmd/wisp call sites off the disk rather than off this paragraph,
// so deleting the line below reddens a case instead of passing quietly.
```
4. **风险 / 现验**：
   - 现验：`panel_host_windows.go:821` 是第二处 `Handle` 产码调用（经 `dispatchRaw`）；仪器 `TestAC9ComposerDispatchHasAProductionCaller`（`panel_inbound_33_test.go:182`）扫**本目录（cmd/wisp）非测试文件**的调用形状，**覆盖不了原文说的 "this repository"**——所以新句改成"cmd/wisp call sites"（与仪器射程一致）。
   - 撞钉风险：AST 尺只认调用形状，注释文本不参与；五枚禁词 0 命中。

## P04 — internal/panel/git.go

1. **目标**：`internal/panel/git.go:76`（现量＝`// ComposerRequest) has zero production callers in this tree - measured, with the`；替换块＝**75–81**）。
2. **原文逐字**（块 75–81）：
```go
// on (the WebView2 host of ticket 33, its postMessage receiver, and a router for
// ComposerRequest) has zero production callers in this tree - measured, with the
// grep lines, in census §⑤. A panel that drew a working switcher on top of that
// would be a button wired to nothing.
//
// The day the inbound hop lands, this constant stops being true and ticket 186
// owns replacing it - not this file.
```
3. **替换文本**（换 75–81；7 行 → 10 行）：
```go
// on (the WebView2 host of ticket 33, its postMessage receiver, and a router for
// ComposerRequest) has since landed in code - the C27 host, the postMessage
// transport and the composer router are in this tree (cmd/wisp/panel_host_windows.go,
// internal/panel/composer_dispatch.go), and the router has production callers on
// the host and stdin legs. What is still not built is the switching dimension
// itself - ticket 186's work, gated on Q-69 - so a panel that drew a working
// switcher on top would still be a button wired to nothing.
//
// The constant's own hop sentence is now stale (it still says 无 postMessage
// 接收器、无 router、无 WebView2 宿主); re-wording it is ticket 186's call, not
// this file's.
```
4. **风险 / 现验**：
   - 现验（现量）：`cmd/wisp/panel_host_windows.go` 在 HEAD 是 C27 WebView2 宿主（`:1 //go:build windows`、`:64` 引 `go-webview2`、安装 postMessage 转发钩子），resident 腿在 `resident_windows.go:151` 实际构建它（`panel_resident_windows.go:253`）；router＝`internal/panel/composer_dispatch.go`；调用者＝`panel_host_windows.go:821`、`panel_inbound.go:163`；Q-69 出自 `composer_dispatch_test.go:716` 的禁词尺文案。
   - ⚠ **邻近同族过期句（不在 13 枚内，本程未动、仅点名）**：`git.go:82` 的 const 串逐字仍写「还没有落地（无 postMessage 接收器、无 router、无 WebView2 宿主）」——与 HEAD 冲突；**const 是否同步改＝产码面，本料不裁**。另 `cmd/wisp/panel_pump.go:13`「no WebView2 host (ticket 33 is unclaimed)」同族。
   - 撞钉：`TestGitDimensionHasNoModelCallableTool`（`git_test.go:355/:364`）钉的是维度形状不是这几行注释；五枚禁词 0 命中（"switching/switcher" 不是禁词）。

## P06 — cmd/wisp/config_readers_255.go

1. **目标**：`cmd/wisp/config_readers_255.go:43`（现量＝`// as its accessor - and TierOf had ZERO production callers (measured at HEAD`；替换块＝**42–45**）。
2. **原文逐字**（块 42–45）：
```go
// landed the tier table in internal/config/tiers.go (commit dd92bb92) with TierOf
// as its accessor - and TierOf had ZERO production callers (measured at HEAD
// 8a3790f0: grep -rn "TierOf(" over cmd internal tools scripts returns the
// definition and nothing else). A receipt that hand-listed hot sections would be a
```
3. **替换文本**（换 42–45；4 行 → 5 行）：
```go
// landed the tier table in internal/config/tiers.go (commit dd92bb92) with TierOf
// as its accessor - and TierOf had ZERO production callers at that point (measured
// at HEAD 8a3790f0: grep -rn "TierOf(" over cmd internal tools scripts returned
// the definition and nothing else). This file is its first production caller
// since. A receipt that hand-listed hot sections would be a
```
4. **风险 / 现验**：
   - 现验：commit `dd92bb92` 与 HEAD `8a3790f0` 均在（`git log -1` 各成）；`TierOf(` 产码＝定义 `internal/config/tiers.go:93` ＋ 唯一调用 `config_readers_255.go:235`（`hotRowsFor`，原代码注释已自称 "TierOf's first production caller"）。新句只把"当时零"说清、把"今天已有"补上，时态分界＝`at that point`。
   - 同文件无逐字钉；原文串在 `*_test.go` 0 命中。

## P07 — cmd/wisp/config_reload.go（三处之第 1，⚠ 从后往前贴）

1. **目标**：`cmd/wisp/config_reload.go:11`（现量＝`//   - Manager.CheckAndReload - its only non-test caller was cmd/balldebug, a`；替换块＝**11–13**）。
2. **原文逐字**（块 11–13）：
```go
//   - Manager.CheckAndReload - its only non-test caller was cmd/balldebug, a
//     side program. `wisp run` never re-read config.toml, so D36's three-tier
//     semantics (立即/重载/重启, PLAN.md:2715) had an engine and no driver.
```
3. **替换文本**（换 11–13；3 行 → 4 行）：
```go
//   - Manager.CheckAndReload - its only non-test caller was cmd/balldebug, a
//     side program. `wisp run` never re-read config.toml, so D36's three-tier
//     semantics (立即/重载/重启, PLAN.md:2721) had an engine and no driver.
//     This file added the driver: reloadOnce() below calls CheckAndReload.
```
4. **风险 / 现验**：
   - **原引 `PLAN.md:2715` 已漂**：HEAD 上 `:2715` 是别的内容；D36 标题现实位于 **`docs/PLAN.md:2721`**（`### D36 — 配置模型（config.toml 全量 section 树 + 三档生效级别）`，三档明细在 `:2726`）。替换文本已换 2721（同一对象＝D36 标题行）。
   - 现验：产码调用点＝`config_reload.go:153`、`cmd/balldebug/main.go:243`；`reloadOnce` 在本文件 `:152` 定义。同文件三块，**按 P09→P08→P07 顺序贴**。

## P08 — cmd/wisp/config_reload.go（三处之第 2）

1. **目标**：`cmd/wisp/config_reload.go:14`（现量＝`//   - Manager.ConfirmLocked - zero assignments outside *_test.go, and`；替换块＝**14–17**）。
2. **原文逐字**（块 14–17）：
```go
//   - Manager.ConfirmLocked - zero assignments outside *_test.go, and
//     manager.go reads it as `approved := m.ConfirmLocked != nil && ...`, so an
//     unwired hook is a silent fail-closed DENY of every [fs] loosening: D33's
//     「热加载放宽必须触发 L2 级重新确认」 was attached to a dead wire.
```
3. **替换文本**（换 14–17；4 行 → 6 行）：
```go
//   - Manager.ConfirmLocked - before this file, zero assignments outside
//     *_test.go; manager.go snapshots the hook under mu (confirm := m.ConfirmLocked)
//     and a nil hook denies (approved := confirm != nil && ...), so an unwired
//     hook was a silent fail-closed DENY of every [fs] loosening: D33's
//     「热加载放宽必须触发 L2 级重新确认」 was attached to a dead wire.
//     startConfigReload assigns it below, before the tick is spawned.
```
4. **风险 / 现验**：
   - **原引语已非逐字**：`manager.go` 现在把读法拆成两步——`:172 confirm := m.ConfirmLocked`（mu 下快照）、`:180 approved := confirm != nil && confirm(c.section, c.keys)`。替换文本按现量改引（`git grep -F 'ConfirmLocked != nil'` 今天只命中这句旧注释本身）。
   - 现验：赋值点＝`config_reload.go:114`，在 `:119` 起 tick 之前（同函数 `startConfigReload`）。

## P09 — cmd/wisp/config_reload.go（三处之第 3）

1. **目标**：`cmd/wisp/config_reload.go:18`（现量＝`//   - Manager.OnRestartPending - zero assignments outside *_test.go too, so the`；替换块＝**18–21**）。
2. **原文逐字**（块 18–21）：
```go
//   - Manager.OnRestartPending - zero assignments outside *_test.go too, so the
//     restart tier said nothing at all: a hand edit of [app] simply did not
//     happen, silently, which is the one shape 票 223 AC#2 forbids
//     (不许用「静默不生效」充当这一档).
```
3. **替换文本**（换 18–21；4 行 → 5 行）：
```go
//   - Manager.OnRestartPending - zero assignments outside *_test.go too before
//     this file, so the restart tier said nothing at all: a hand edit of [app]
//     simply did not happen, silently, which is the one shape 票 223 AC#2 forbids
//     (不许用「静默不生效」充当这一档). startConfigReload assigns it below,
//     before the tick is spawned.
```
4. **风险 / 现验**：
   - 现验：赋值点＝`config_reload.go:115`；读侧＝`manager.go:201-202`（`m.OnRestartPending != nil` 才调）。归口＝票 223 `AC#7`（其判据今天已被产码满足却仍未勾——**勾框不是本程的事**，triage 已记）。

## P10 — cmd/wisp/firstrun.go

1. **目标**：`cmd/wisp/firstrun.go:11`（现量＝`// no production path ever did was CALL that pair once, at the moment a user`；替换块＝**8–15**）。
2. **原文逐字**（块 8–15）：
```go
// (NewDefaults, defaults.go:58 - the `default:"..."` tags in schema.go are the
// single source, D36 rule 3) and the atomic sealed writer (SaveFile,
// loader.go:238), and writeguard.go:125-131 already proves a missing file is
// creatable through them ("writing it creates what first-run did not"). What
// no production path ever did was CALL that pair once, at the moment a user
// first asks `wisp run` for work. On a fresh machine the run leg therefore
// died at run.go's 配置未就绪 branch with exit 2, quoting a missing file as
// the user's only experience of the product.
```
3. **替换文本**（换 8–15；8 行 → 10 行）：
```go
// (NewDefaults, defaults.go:58 - the `default:"..."` tags in schema.go are the
// single source, D36 rule 3) and the atomic sealed writer (SaveFile,
// loader.go:252), and writeguard.go:125-131 already proves a missing file is
// creatable through them ("writing it creates what first-run did not"). What
// no production path had done before this file was CALL that pair once, at the
// moment a user first asks `wisp run` for work: ensureFirstRunConfig below calls
// config.SaveFile(cfgPath, config.NewDefaults()), reached from the run entry
// (run.go). Before it, a fresh machine's run leg died at run.go's 配置未就绪
// branch with exit 2, quoting a missing file as the user's only experience of
// the product.
```
4. **风险 / 现验**：
   - **原引 `loader.go:238` 已漂**：`func SaveFile` 现实在 `internal/config/loader.go:252`（triage 记漂 14 行，本程复核一致）；替换文本已换 252。
   - 现验：首建调用＝`firstrun.go:82 config.SaveFile(cfgPath, config.NewDefaults())`；入口＝`run.go:245 ensureFirstRunConfig(...)` 调用；`配置未就绪` 仍在 `run.go:417`；`NewDefaults`＝`defaults.go:58`；`writeguard.go:125-131` 缺文件分支与逐字引语仍在（`:127`）。
   - 仪器：`firstrun_257_test.go:430-438` 钉 `\nfunc ` 计数＝1 与 `func ensureFirstRunConfig(dataDir string, stderr io.Writer) (bool, error) {` 签名逐字——替换文本**只增注释行**，不含新增 `func`、不含该完整签名（现验）。

## P11 — cmd/wisp/resident_approval_windows.go

1. **目标**：`cmd/wisp/resident_approval_windows.go:282`（现量＝`// wins) but has no production caller; a second assembly in one process belongs`；替换块＝**279–283**）。
2. **原文逐字**（块 279–283）：
```go
// One holder, one bind site, and it happens before any task can be submitted -
// so no card this process shows can be answered while the holder is unbound for
// a reason nobody named. Re-binding is allowed by construction (last writer
// wins) but has no production caller; a second assembly in one process belongs
// to the ⓐ-Ⅲ shape this ruling refused.
```
3. **替换文本**（换 279–283；5 行 → 6 行）：
```go
// One holder, one bind site, and it happens before any task can be submitted -
// so no card this process shows can be answered while the holder is unbound for
// a reason nobody named. Re-binding is allowed by construction (last writer
// wins) and has no production caller - the single bind is startResidentTaskSource's
// (resident_task_source_windows.go:309); a second assembly in one process
// belongs to the ⓐ-Ⅲ shape this ruling refused.
```
4. **风险 / 现验**：
   - 现验：`bindResidentGrantLedger` 产码调用**恰一处**＝`resident_task_source_windows.go:309`（`if run != nil { ra.bindResidentGrantLedger(run.session) }`）；同文件 `:322` 已写"called by startResidentTaskSource"（原句与之自相矛盾）。替换文本把"唯一绑点"点名，同时保留"re-binding 无生产调用者"这一仍成立的事实（一次绑定 ≠ 重绑）。
   - 仪器：`resident_grant_writer_265_windows_test.go:465` 起是 **AST 读 bind 位**——注释不进 AST，贴注释不影响。归口票 246（`AC#7` 已勾、`AC#8` 不承接，照 triage）。

## P12 — cmd/wisp/approval_reply.go

1. **目标**：`cmd/wisp/approval_reply.go:12`（现量＝`// Gate.DecideFromNative / Gate.DecideFromPanel / Gate.Veto, had zero callers`；替换块＝**12–16**）。
2. **原文逐字**（块 12–16）：
```go
// Gate.DecideFromNative / Gate.DecideFromPanel / Gate.Veto, had zero callers
// outside _test.go. Every card this process shows was therefore answered by
// nobody: an L2 card waited out the full C18 deadline and auto-rejected, and an
// L1 window could not be opposed at all. 票 201's 现量 table names the same hole
// from the other end ("三种人能答复的入口 - 生产零调用者").
```
3. **替换文本**（换 12–16；5 行 → 7 行）：
```go
// Gate.DecideFromNative / Gate.DecideFromPanel / Gate.Veto, had zero callers
// outside _test.go at census time. Every card this process showed was therefore
// answered by nobody: an L2 card waited out the full C18 deadline and
// auto-rejected, and an L1 window could not be opposed at all. 票 201 :10 names
// the same hole from the other end, row verbatim: | 三种"人能答复"的入口 | **生产零调用者** |.
// This file (with the Replies routes it drives) is what closed it - .Veto,
// .DecideFromNative and .DecideFromPanel all have production call sites today.
```
4. **风险 / 现验**：
   - **原引文是近似转写**（triage 已点名）：票面 `:10` 逐字行＝`| 三种"人能答复"的入口 | **生产零调用者** | \`grep ...\` ⇒ 只剩定义行 |`；替换文本引其两格逐字并标 "row verbatim"。票面 `:73` 已自行登记"该尺读数已过期"（**不算本程新发现**）。
   - 现验：`.Veto(` 产码两处＝`approval_reply.go:363`、`resident_approval_windows.go:612`；`.DecideFromNative/.DecideFromPanel` 产码＝`internal/agent/approval/replies.go:328/:408/:413/:437`。时态改为 `at census time`/`showed`。

## P13 — internal/risk/assessor.go

1. **目标**：`internal/risk/assessor.go:29`（现量＝`// not yet wired by tickets 18/19) leave the corresponding rule DORMANT — an`；替换块＝**28–32**）。
2. **原文逐字**（块 28–32）：
```go
// Absent dependencies (PathCanonicalizer / SensitiveClassifier / TaintDetector
// not yet wired by tickets 18/19) leave the corresponding rule DORMANT — an
// integration gap, not a judged result. A WIRED dependency that panics or
// errors is a judged failure: panic -> R9 fail-closed L2 (recovered here);
// canonicalization error -> R2 fail-closed L2.
```
3. **替换文本**（换 28–32；5 行 → 7 行）：
```go
// The tickets 18/19 dependencies (PathCanonicalizer / SensitiveClassifier /
// TaintDetector) are wired at the assembly sites today (internal/tools/bridge.go's
// New and assessorFor; cmd/wisp/panel_assets.go for the C25 detector); where an
// assembly leaves one absent, the corresponding rule stays DORMANT — an
// integration gap, not a judged result. A WIRED dependency that panics or
// errors is a judged failure: panic -> R9 fail-closed L2 (recovered here);
// canonicalization error -> R2 fail-closed L2.
```
4. **风险 / 现验**：
   - 现验：`WithCanonicalizer`/`WithSensitiveClassifier` 装配＝`internal/tools/bridge.go:211-213`（`New` 内）与 `:924-926`（`assessorFor`）；`WithTaintDetector`＝`bridge.go:927` ＋ `cmd/wisp/panel_assets.go:66`（条件装配）。票 18/19 均已 `-done`（triage 记）。
   - ⚠ **该句位于 `assessor.go:9-36` 的 "FROZEN CONTRACT — C19" 横幅内**。规则 ID/语义未动，只动依赖状态描述；**"横幅内描述句是否算契约文本"的边界我判不动**（缺的尺＝人工/编排者裁定）。
   - 无同族仪器扫本文件；原文串在 `*_test.go` 0 命中。

## P39 — cmd/wisp/run.go

1. **目标**：`cmd/wisp/run.go:333`（现量＝`// Nothing calls it yet - the WebView2 "event -> ParseComposerRequest" hop does`；替换块＝**331–336**，行前有 **1 个 tab**）。
2. **原文逐字**（块 331–336，逐字含前导 tab）：
```go
	// modeWrites is ticket 114 AC#2's gate: the one handler a panel mode request
	// is answered by, assembled with the SAME L2 leg the Store above was handed.
	// Nothing calls it yet - the WebView2 "event -> ParseComposerRequest" hop does
	// not exist in this tree (tickets 33/35) - so this is the gate's other half,
	// not a path: the host that gets wired later cannot assemble a wider档 without
	// a confirmation leg.
```
3. **替换文本**（换 331–336；6 行 → 9 行，粘贴时保留前导 tab）：
```go
	// modeWrites is ticket 114 AC#2's gate: the one handler a panel mode request
	// is answered by, assembled with the SAME L2 leg the Store above was handed.
	// Nothing reads this field yet: the "event -> ParseComposerRequest" hop now
	// exists (composer_dispatch.go:155, reached from `wisp panel-inbound` and the
	// panel host), but every leg that reaches it assembles its own handler chain
	// (cmd/wisp/panel_inbound.go, the resident leg's newResidentComposerDispatch)
	// and nothing attaches this one - so this is the gate's other half, not a
	// path: the leg that gets wired to it later cannot assemble a wider档 without
	// a confirmation leg.
```
4. **风险 / 现验**：
   - 现验：`modeWrites` 在产码只有 定义（`run.go:337`）＋ 赋值（`run.go:674`）＋ 注释；**读侧 0**（`panel_inbound.go:248/:272` 是**同名局部量**，不读该字段）。过期的是**理由分句**（"hop does not exist"）：`ParseComposerRequest(` 产码调用＝`internal/panel/composer_dispatch.go:155`；到达它的腿各自装自己的 handler 链（`panel_inbound.go:248`/`:271`；resident 腿 `panel_resident_windows.go:211 newResidentComposerDispatch`）。结论分句（该字段无人读）**今天仍成立**。
   - 撞钉：五枚禁词 0 命中；`panel_inbound_33_test.go` 的 AST 尺不受注释影响。

---

## 汇总表（13 行）

| 出处（HEAD 现量） | 归票（照 01-triage） | 撞禁词 | 引用未验对象 |
|---|---|---|---|
| P01 `internal/panel/composer_dispatch.go:47-52` | 票 33 `AC#9`（未勾，承接格） | 否（0/五枚） | 否（证据件 147 行＋§④ 现验） |
| P02 `internal/agent/approval/doc.go:26-37` | 票 12（未勾 3 格无一承接） | 否 | 否（5 处产码点全现量） |
| P03 `cmd/wisp/panel_inbound.go:11-14` | 票 33 `AC#9` | 否 | 否（仪器件＋第二调用点现验） |
| P04 `internal/panel/git.go:75-81` | 票 186（未勾 9，`AC#2` 不承接） | 否 | 否（C27 宿主/router/Q-69 现验；⚠ `git.go:82` const 同族过期未动） |
| P06 `cmd/wisp/config_readers_255.go:42-45` | 票 255（未勾 2，非注释格） | 否 | 否（commit+HEAD+tiers.go 现验） |
| P07 `cmd/wisp/config_reload.go:11-13` | 票 223（未勾 2，非注释格） | 否 | 否（**原引 PLAN.md 行已修 2715→2721**） |
| P08 `cmd/wisp/config_reload.go:14-17` | 票 223 | 否 | 否（**原引 manager.go 读法已修成现行 verbatim**） |
| P09 `cmd/wisp/config_reload.go:18-21` | 票 223 `AC#7`（判据已满足未勾，勾归非实现者） | 否 | 否（`manager.go:201-202` 现验） |
| P10 `cmd/wisp/firstrun.go:8-15` | 票 198（已 `-done`，**不可承接**，去向待编排者） | 否 | 否（**原引 loader.go 行已修 238→252**；首建点现验） |
| P11 `cmd/wisp/resident_approval_windows.go:279-283` | 票 246（`AC#8` 不承接） | 否 | 否（`:309` 唯一绑点现验） |
| P12 `cmd/wisp/approval_reply.go:12-16` | 票 201（未勾 5，`AC#1` 承接） | 否 | 否（**票面 :10 逐字行已现取改正**） |
| P13 `internal/risk/assessor.go:28-32` | 票 18/19（均 `-done`，去向待编排者） | 否 | 否（bridge/panel_assets 装配点现验） |
| P39 `cmd/wisp/run.go:331-336` | 票 114（未勾 9）＋票 33/35 | 否 | 否（`composer_dispatch.go:155` 等现验） |

**我判不动的格子**（3 处，均写明缺哪把尺）：
1. **P13**：替换句在 `assessor.go:9-36` 的 C19 FROZEN CONTRACT 横幅内——"横幅内描述句可不可以只改文案"的边界我判不动；缺的尺＝人工/编排者裁定（照 triage 本枚已判〔只欠文案〕）。
2. **P04**：注释下方 `git.go:82` 的 const 串同族过期（逐字冲突已点名），**const 是否连带改＝产码面**我判不动；缺的尺＝编排者/票 186 口径。
3. **P04/P39 的"hop 已落地"粒度**：我只按代码事实写（宿主/传输/router/调用点都在 HEAD）；"落地"到用户可见/票面验收的粒度判不动，缺的尺＝票 33/35 的验收口径（非本程）。

**体检读数**（本料末尾三条均为现跑，非转述）：
- 禁词体检：五枚禁词（名册照编排者转述，本表不逐字复写）对 `01-ready-to-apply.md` 逐枚 `grep -F` ⇒ **全部 0 命中**；替换文本块内亦 0（本料设计上不含）。
- 原句撞钉复测：13 条锚句对 `git grep -F ... HEAD -- '*_test.go'` ⇒ **0 命中（rc=1）**，与 01-triage 的 15 枚全 0 一致。
- 新文本关键串对 `*_test.go` 的命中：`ensureFirstRunConfig`／`startResidentTaskSource`／`bindResidentGrantLedger`／`dispatchRaw`／`reloadOnce`／`newResidentComposerDispatch` 等会命中既有判据件（属**仪器邻接**，非禁词；明细见各块"风险"）。四处"改文案有形状约束"的仪器（`composer_dispatch_test.go:375/:708`、`firstrun_257_test.go:430`、`resident_grant_writer_265_windows_test.go:465`）逐把结论已写在各块。

⛔ 本料不提"该不该改"（编排者的事）；P14/P15 未碰；产码/票面/台账零改动。
