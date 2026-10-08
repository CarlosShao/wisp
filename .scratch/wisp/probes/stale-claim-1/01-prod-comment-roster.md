# stale-claim-1 / 01 — 表一：产码注释名册（HEAD `6547fd30`）

射程：`HEAD` 的产码·非测试 Go 文件（剥 `.scratch`、剥 `_test.go`）。
内容扫 101 行命中 -> 逐枚读整行/整块之后，**判为状态级断言的是 39 枚**（本表 P01–P38 ＋ 交件后补的 P39），另有 **4 枚在 `_test.go`**（本表末「表一·补」）。
分档计数（含 P39）：**已过期 15 ／ 判不动 9 ／ 仍成立 15**。
（甲节标"14 枚"是 P39 补入前的数；P39 是**半枚过期**，计入已过期一侧，故甲节实际 15 枚。）
每枚的"尺"都是 §1.2 那把调用形状尺，命令原文在 `00-ruler-and-inventory.md`；引文一律是 `git show HEAD:<file>` 的整行，未截断。

图例：`包内`=非测试同包调用点 ／ `包外`=非测试跨包调用点 ／ `入口`=到得了 `main` 的哪条生产入口。

---

## 甲 · 已过期（14 枚）——今天拿去派单会直接害事的

### P01 `HEAD:internal/panel/composer_dispatch.go:50`
> `// belongs to the host; and it has NO production caller yet. The honest state of`

- 断言符号：`(*panel.ComposerDispatch).Handle`
- 尺：`git grep -nE "ComposerDispatch|\.Handle\(" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go' '*.go'`
- 输出摘要：非测试调用点 **3 处**。`cmd/wisp/panel_inbound.go:163 reply, err := disp.Handle(ctx, raw)`；`cmd/wisp/panel_host_windows.go:821 return m.disp.Handle(ctx, raw)`；`cmd/wisp/panel_resident_windows.go:234 disp, err := newResidentComposerDispatch(dataDir, auditf)`。
- 三层：包内 0 ／ 包外 3（全在 `cmd/wisp`，package main）／ 入口 `main.go:115 case "panel-inbound" -> cmdPanelInbound`（`main.go:120 os.Exit(cmdPanelInbound(args[1:], panelInboundIO{}))`）+ 面板宿主腿 + 常驻腿。
- **分档：已过期〔已证〕。现在到得了它的是 `HEAD:cmd/wisp/panel_inbound.go:163`（经 `HEAD:cmd/wisp/main.go:115`）**，另有 `HEAD:cmd/wisp/panel_host_windows.go:821`、`HEAD:cmd/wisp/panel_resident_windows.go:234`。
- ⚠ 同一枚断言的行号在台账/票面里被引成 `composer_dispatch.go:48` —— HEAD 上 `:48` 不是这句（**漂 2 行**）；且 `internal/agent/composer/` 这个包路径在 HEAD **不存在**（`git show HEAD:internal/agent/composer/composer_dispatch.go | wc -l` = **0**）。

### P02 `HEAD:internal/agent/approval/doc.go:26`
> `// Wiring still owed by ticket 12 (this package has NO production caller yet,`
> `// so nothing here runs in cmd/wisp today):`

- 断言符号：整包 `internal/agent/approval`（含 `approval.New`、`Gate`）
- 尺：`git grep -nE '"github[^"]*internal/agent/approval"' HEAD -- …` ＋ `git grep -nE "approval\.New\(" HEAD -- …`
- 输出摘要：非测试 importer **4 枚**（`cmd/wisp/approval_reply.go:67`、`cmd/wisp/resident_approval_windows.go:61`、`cmd/wisp/resident_task_source_windows.go:88`、`cmd/wisp/run.go:47`）；非测试 `approval.New(` **2 枚**（`cmd/wisp/resident_approval_windows.go:369`、`cmd/wisp/run.go:612`）。
- 三层：包内 有 ／ 包外 有（4 处 import）／ 入口 `runTextTask`（`run.go:181`）与常驻腿（`resident_windows.go:260 -> startResidentTaskSource`）**两条都到**。
- **分档：已过期〔已证〕。到得了它的是 `HEAD:cmd/wisp/run.go:612` 与 `HEAD:cmd/wisp/resident_approval_windows.go:369`。**
- 加重项：这句带 **`today`** 一词，且它下面列的五条"still owed"里至少两条（`approval.New` 作为 Gate、`Queue.Native()` 交给宿主）**已经落地**（见 P12、表三 T-04）。⛔ 我不动这段字。

### P03 `HEAD:cmd/wisp/panel_inbound.go:11`
> `// line and hands the bytes to (*panel.ComposerDispatch).Handle. Nothing else in`
> `// this repository calls Handle today, and the judgement in`

- 断言符号：`(*panel.ComposerDispatch).Handle`（"本仓除本腿外没人调 Handle"）
- 尺：同 P01
- 输出摘要：`cmd/wisp/panel_host_windows.go:821` 是一枚**除本腿之外**的调用点。
- 三层：包内（`cmd/wisp` 自身）有 ／ 包外 不适用（同包）／ 入口 = WebView2 宿主腿，见末节的看得见/看不见。
- **分档：已过期〔已证〕。多出来那一处 = `HEAD:cmd/wisp/panel_host_windows.go:821`。**
- ⚠ 这句自陈"judgement in `panel_inbound_33_test.go` reads THAT fact off the disk"（表一·补 T02）——即注释与仪器**互相咬合在同一枚已经动了的事实上**，这一族最脆的形状。

### P04 `HEAD:internal/panel/git.go:76`
> `// ComposerRequest) has zero production callers in this tree - measured, with the`
> `// grep lines, in census §⑤. A panel that drew a working switcher on top of that`
> `// would be a button wired to nothing.`
> （整块 `:72-78`；同一枚还带 `:79-80` `The day the inbound hop lands, this constant stops being true and ticket 186 owns replacing it - not this file.`）

- 断言符号：票 33 的 WebView2 宿主 + postMessage 接收器 + `ComposerRequest` 路由器**这一整条入向通路**
- 尺：P01 的尺（`ComposerDispatch`/`.Handle(`）＋ `git grep -nE "panel-inbound|cmdPanelInbound" HEAD -- …`
- 输出摘要：路由器与宿主**都已落**（`main.go:115`、`panel_inbound.go:163`、`panel_host_windows.go:821`、`panel_resident_windows.go:234`）。
- 三层：包外 有 ／ 入口 **到得了 `main`**（`wisp panel-inbound` 是注册在 `main.go` 的子命令）。
- **分档：已过期〔已证〕。到得了它的是 `HEAD:cmd/wisp/main.go:115` + `HEAD:cmd/wisp/panel_host_windows.go:821`。**
- ⚠ 这枚自己写了"到期就该由票 186 换掉、不是本文件换"——**到期条件已满足，替换未发生**（`GitSwitchBlockedReason` 在 HEAD 上原文还在，见 P05）。

### P06 `HEAD:cmd/wisp/config_readers_255.go:43-45`
> `// as its accessor - and TierOf had ZERO production callers (measured at HEAD`
> `// 8a3790f0: grep -rn "TierOf(" over cmd internal tools scripts returns the`
> `// definition and nothing else). A receipt that hand-listed hot sections would be a`

- 断言符号：`config.TierOf`
- 尺：`git grep -nE "TierOf\(" HEAD -- ':(exclude).scratch' '*.go'`
- 输出摘要：非测试调用点 **1 枚**，且它就在**同一个文件的 `:235`**（`if tier, ok := config.TierOf(name); ok { // TierOf's first production caller`），定义 `internal/config/tiers.go:93`。
- 三层：包外（`cmd/wisp` -> `internal/config`）有 ／ 入口 同一枚函数 `hotRowsFor` 在本文件 `:266`、`:292` 被调（回执路径）；`hotRowsFor` 的更上游入口我未追 -> 见末节 B-06。
- **分档：已过期〔已证〕。到得了它的是 `HEAD:cmd/wisp/config_readers_255.go:235`（同文件自盖）。**

### P07 `HEAD:cmd/wisp/config_reload.go:12`
> `//   - Manager.CheckAndReload - its only non-test caller was cmd/balldebug, a`

- 断言符号：`(*config.Manager).CheckAndReload`
- 尺：`git grep -nE "CheckAndReload\(" HEAD -- ':(exclude).scratch' '*.go'`
- 输出摘要：非测试调用点 **2 枚**：`cmd/wisp/config_reload.go:153 rep, err := rt.mgr.CheckAndReload()` 与 `cmd/balldebug/main.go:243`（侧程序，仍在）。定义 `internal/config/manager.go:142`。
- 三层：包外 有 ／ 入口 `run.go:813 rt.startConfigReload()` -> `config_reload.go:105` -> `:153`。**到得了 `main`。**
- **分档：已过期〔已证〕。到得了它的是 `HEAD:cmd/wisp/config_reload.go:153`（入口 `HEAD:cmd/wisp/run.go:813`）。**

### P08 `HEAD:cmd/wisp/config_reload.go:14`
> `//   - Manager.ConfirmLocked - zero assignments outside *_test.go, and`

- 断言符号：`Manager.ConfirmLocked`（函数字段）
- 尺：`git grep -nE "ConfirmLocked|OnRestartPending" HEAD -- ':(exclude).scratch' '*.go'`
- 输出摘要：非测试赋值点 **1 枚** = `cmd/wisp/config_reload.go:114 rt.mgr.ConfirmLocked = rt.confirmLockedLoosening`；读侧 `internal/config/manager.go:172 confirm := m.ConfirmLocked`。
- 三层：包外 有 ／ 入口同 P07。
- **分档：已过期〔已证〕。到得了它的是 `HEAD:cmd/wisp/config_reload.go:114`。**

### P09 `HEAD:cmd/wisp/config_reload.go:18`
> `//   - Manager.OnRestartPending - zero assignments outside *_test.go too, so the`

- 尺/三层：同 P08。赋值点 `cmd/wisp/config_reload.go:115 rt.mgr.OnRestartPending = rt.reportRestartPending`；读侧 `internal/config/manager.go:201-202`。
- **分档：已过期〔已证〕。到得了它的是 `HEAD:cmd/wisp/config_reload.go:115`。**
- ⚠ P07/P08/P09 三枚都在同一个 `WHY THIS FILE EXISTS` 块里（`:7` 起 `Three of config.Manager's hooks had no production / assignment point before this file (measured at HEAD 7ffa9520, census .scratch/wisp/probes/223/c1/census.md)`）。`:7` 那枚**过去式＋带 ref＋带日期**，判 **仍成立（历史）**，见 P38；`:12/:14/:18` 三条子弹是**现在式口吻的事实句**，判已过期。**同一个块里两档，是最容易被下游只读子弹不读框架的一族。**

### P10 `HEAD:cmd/wisp/firstrun.go:11`
> `// loader.go:238), and writeguard.go:125-131 already proves a missing file is`
> `// creatable through them ("writing it creates what first-run did not"). What`
> `// no production path ever did was CALL that pair once, at the moment a user`

- 断言符号：`config.NewDefaults` + `config.SaveFile` 这一对
- 尺：`git grep -nE "NewDefaults\(|SaveFile\(" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go' '*.go'`
- 输出摘要：非测试点 `cmd/wisp/firstrun.go:82 if err := config.SaveFile(cfgPath, config.NewDefaults()); err != nil`；入口 `cmd/wisp/run.go:245 if _, frErr := ensureFirstRunConfig(s.dataDir, s.stderr); frErr != nil`。`ensureFirstRunConfig` 定义 `firstrun.go:72`，非测试调用点**只有 run.go:245 一枚**。
- 三层：包外 有 ／ 入口 `main.go:91 -> cmdRun(:151) -> runTextTask(:159)`；T1 尺（`awk 'NR>=176 && NR<=260 && /^func /'` 只回 `:181 func runTextTask`）**证实 `:245` 落在 `runTextTask` 函数体内**。
- **分档：已过期〔已证〕（作为现状读）。到得了它的是 `HEAD:cmd/wisp/firstrun.go:82`，入口 `HEAD:cmd/wisp/run.go:245`。**

### P11 `HEAD:cmd/wisp/resident_approval_windows.go:282`
> `// wins) but has no production caller; a second assembly in one process belongs`

- 断言符号：`(*residentApproval).bindResidentGrantLedger`
- 尺：`git grep -nE "bindResidentGrantLedger" HEAD -- ':(exclude).scratch' '*.go'`
- 输出摘要：非测试调用点 **1 枚** = `cmd/wisp/resident_task_source_windows.go:309 ra.bindResidentGrantLedger(run.session)`；定义 `resident_approval_windows.go:284`。
- 三层：包内（同 package main）有 ／ 入口 `resident_windows.go:260 src := startResidentTaskSource(rt, ra)` -> `resident_task_source_windows.go:228` -> `:309`。**到得了常驻腿入口。**
- **分档：已过期〔已证〕。到得了它的是 `HEAD:cmd/wisp/resident_task_source_windows.go:309`。**
- ⚠ **同文件 `:322` 的注释已经改口了**（`session ledger by bindResidentGrantLedger, called by startResidentTaskSource`）—— 一枚文件内自相矛盾：`:282` 说没有生产调用者，`:322` 说调用者是 `startResidentTaskSource`。

### P12 `HEAD:cmd/wisp/approval_reply.go:11-16`
> `// switch raises a host-side L2 card - while the ANSWERING side,`
> `// Gate.DecideFromNative / Gate.DecideFromPanel / Gate.Veto, had zero callers`
> `// outside _test.go. Every card this process shows was therefore answered by`
> …
> `// from the other end ("三种人能答复的入口 - 生产零调用者").`

- 断言符号：`Gate.DecideFromNative` / `Gate.DecideFromPanel` / `Gate.Veto`
- 尺：`git grep -nE "\.DecideFromNative\(|\.DecideFromPanel\(|\.Veto\(" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go' '*.go'`
- 输出摘要（非测试）：`.Veto(` -> `cmd/wisp/approval_reply.go:363 err := s.live.h.Veto(corr)`、`cmd/wisp/resident_approval_windows.go:612 if err := ra.cards.Veto(card.CorrelationID); err != nil`。`DecideFromNative/Panel` -> `internal/agent/approval/replies.go:328/:408/:413/:437`（**同包生产调用**，经 `Replies` 面），而 `Replies` 被 `cmd/wisp/approval_reply.go:302 s.live.h.PanelReject(...)` / `:331 s.live.h.PanelAllow(...)` 驱动。
- 三层：包内 有 ／ 包外 有（经 `Replies`）／ 入口 `runTextTask`（`run.go:612` 建 gate）与常驻腿都到。
- **分档：已过期（作为现状读）〔已证〕。到得了它的是 `HEAD:cmd/wisp/approval_reply.go:363`、`HEAD:cmd/wisp/resident_approval_windows.go:612`、`HEAD:internal/agent/approval/replies.go:328/:408/:413/:437`。**
- **降档保护（写清楚）**：这句**自带日期与出处**（`:7` `The 09-29 census of the approval reply surface (docs/evidence/s1/219-approval-reply-surface-c1.md §②A) measured…`），且是过去式。按尺判"**历史陈述仍成立、现在式读法已过期**"；分档记已过期是因为 `:16` 那句**引用票面时用现在式**（`"三种人能答复的入口 - 生产零调用者"`）。读法风险＝高。

### P13 `HEAD:internal/risk/assessor.go:28-29`
> `// Absent dependencies (PathCanonicalizer / SensitiveClassifier / TaintDetector`
> `// not yet wired by tickets 18/19) leave the corresponding rule DORMANT — an`

- 断言符号：`RiskAssessor.WithCanonicalizer` / `WithSensitiveClassifier` / `WithTaintDetector`（"还**没接**"）
- 尺：`git grep -nE "WithCanonicalizer|WithSensitiveClassifier|WithTaintDetector" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go' '*.go'`
- 输出摘要（非测试调用形状）：`internal/tools/bridge.go:212-213`（`NewBridge` 默认装配，canonicalizer + classifier）、`internal/tools/bridge.go:925-927`（`assessorFor`，三枚全上，含 `WithTaintDetector(b.prov.Detector(taskID))`）、`cmd/wisp/panel_assets.go:66`（`assessor = assessor.WithTaintDetector(det)`）。定义在 `assessor.go:193/:200/:207`。
- 三层：包外 有 ／ 入口 `run.go:544` 一带的 bridge 装配 -> `agent.Loop` -> `l.opt.Tools.Execute(ctx, req)`（`internal/agent/loop.go:737`）。到得了生产入口。
- **分档：已过期〔已证〕。到得了它们的是 `HEAD:internal/tools/bridge.go:212-213` 与 `HEAD:internal/tools/bridge.go:925-927`。**
- ⚠ 这枚是**风险最高的一枚**：票 18/19 早已落，而**这句在 assessor 的头文档里**，正是"某条 R 规则是不是 DORMANT"的权威读面。`DORMANT` 这个词本身（`:29`）仍是合法的行为描述——**过期的是 `not yet wired by tickets 18/19` 这一枚括号**。

### P14 `HEAD:cmd/wisp/resident_task_source_windows.go:13-14`
> `// 22:5x found AskOnTaskRoot (resident_approval_windows.go:261) and`
> `// askConfirmation (:205) with ZERO product callers - only this package's own`

- 断言符号：`(*residentApproval).AskOnTaskRoot` / `(*residentApproval).askConfirmation`
- 尺：`git grep -nE "AskOnTaskRoot" HEAD -- ':(exclude).scratch' '*.go'`（全谱，含测试）＋ `git grep -nE "askConfirmation" HEAD -- …`
- 输出摘要：`AskOnTaskRoot` 非测试调用点 **0**（只在 `resident_approval_windows.go:688` 注释、`:697` 定义、`resident_task_source_windows.go:13` 注释、`resident_windows.go:248` 注释 + 3 枚 `_test.go:115/:220/:311`）。`askConfirmation` 非测试调用点 **1** = `resident_approval_windows.go:700 return ra.askConfirmation(taskCtx, c)`。
- 三层＋**行号漂移两处**：`AskOnTaskRoot` 现在在 `:697`（注释写 `:261`）、`askConfirmation` 现在在 `:641`（注释写 `:205`）。
- **分档：已过期〔已证〕（`askConfirmation` 那半句 + 两处行号）。现在到得了 `askConfirmation` 的是 `HEAD:cmd/wisp/resident_approval_windows.go:700`——但那一枚住在 `AskOnTaskRoot` 体内，而 `AskOnTaskRoot` 生产调用点为 0 ⇒ `askConfirmation` = 〔建了但没接〕（编译进产物、生产走不到）。**

### P15 `HEAD:cmd/wisp/resident_windows.go:248`
> `// the way out - and AskOnTaskRoot / askConfirmation still had zero product`
> `// callers, so nothing was going to raise one. This call is the caller: it reads`

- 断言符号：同上两枚 ＋ **正向断言"这一发就是那个调用者"**
- 尺：同 P14 ＋ `git show HEAD:cmd/wisp/resident_windows.go | sed -n '252,275p'`
- 输出摘要：这一块说的"这一发"其实是 `resident_windows.go:260 src := startResidentTaskSource(rt, ra)`。`startResidentTaskSource` 定义在 `resident_task_source_windows.go:228`（非测试），但**它不调 `AskOnTaskRoot`**（U2 尺：`AskOnTaskRoot(` 非测试 0 命中）。
- **分档：已过期〔已证〕（"`This call is the caller`" 这一枚）。**票 246 AC#7 要补的那一跳（`AskOnTaskRoot` 的生产调用者）**今天仍然是 0**；补上的是 `bindResidentGrantLedger` 那一枚（P11）和 `askConfirmation` 的包内一处。⇒ 这一枚是**"票已勾、断言未真"**的形状，害度最大。

---

## 乙 · 判不动（9 枚）——缺哪一发读数写在 `04` 里，逐枚具名

### P05 `HEAD:internal/panel/git.go:82`（`GitSwitchBlockedReason` 的**用户可见文案**本体）
> `const GitSwitchBlockedReason = "切换分支／切换工作树今天不可用：网页到宿主的那一跳还没有落地（无 postMessage 接收器、无 router、无 WebView2 宿主），面板只有快照这一条出向通道"`

- 这枚不只是注释：**它进快照**（`git.go:160 SwitchBlocked: GitSwitchBlockedReason`），是给用户看的一句。
- 尺：P01/P04 的尺。`无 router`、`无 WebView2 宿主` 两半按 `panel_host_windows.go:152/:213/:821` + `main.go:115` **判过期**；`无 postMessage 接收器` 那半**判不动** —— 与之直接冲突的另一枚读数在 `_test.go`：`HEAD:internal/panel/composer_dispatch_test.go:440` 逐字 `no native host / WebView2 message channel is attached in this tree`（表一·补 T03）。两把尺一真一假不能同存，我这一程⛔ 不能跑仪器去裁。
- **分档：判不动（半过期 + 仪器与产码互相矛盾）。**

### P16 `HEAD:cmd/wisp/models.go:36-39`
> `// window is *reachable*. There is still no reader of model bytes anywhere in`
> ``cmd/wisp`'s dependency graph - this process never opens a file under`
> `<store>/<id>/, it only re-hashes it - so the segment has no end point in the`
> `// product today. It gets an end point the moment an engine loader is wired, and`

- 断言符号：模型字节的生产读侧（"still no reader … today"，**现在式**）
- **分档：判不动。**尺不是调用形状能给的：要的是"谁在产码里 `os.Open`/`os.ReadFile` `<store>/<id>/`"的**文件访问扫**，加上 `cmd/wisp` 依赖图闭包。缺的读数＝`04` B-01。

### P17 `HEAD:internal/risk/provenance.go:64`
> `// DisposalScope" half, because no *plugin.DisposalScope reaches the plugin agent command`
> `// in production today (160-c1 §2.2) and laying that pipe means touching`

- **分档：判不动。**这句自带出处（`160-c1 §2.2`）＝〔仅注释 + 具名读数件〕，我这一程没复跑那一把尺。缺的读数＝`04` B-02（`DisposalScope` 的调用形状）。

### P18 `HEAD:internal/tools/bridge.go:872`
> `// Today no production code dispatches on the bridge except loop.go, so being`

- 断言符号：`*Bridge` 的派发入口（**现在式 + 具名例外**）
- 尺：我的 `.Dispatch(` 扫**不可用** —— 命中里 `statemachine.Machine.Dispatch`、`windows pump.Dispatch` 是同形不同符号（`internal/models/bridge.go:40/46/59/64`、`internal/statemachine/machine.go:202`、`cmd/wisp/panel_host_windows.go:613`、`cmd/wisp/panel_resident_windows.go:360`、`cmd/balldebug/main.go:618`），必须先把 `func (b *Bridge) Dispatch|Execute` 的真实方法名抽出来再用带接收者的形状量。唯一确定的正向点是 `internal/agent/loop.go:737 return l.opt.Tools.Execute(ctx, req)`（接口形状，不是 `*Bridge` 具名调用）。
- **分档：判不动。**缺的读数＝`04` B-03。

### P19 `HEAD:internal/winsec/resolve.go:231`
> `// measure instead of assert: no caller's path passes through here; ResolvePath`

- 断言符号：`resolverProbeRoot` / `resolveProbeRoot`（"没有**生产调用者的路径**经过这里"）
- 尺：`git grep -nE "resolverProbeRoot" HEAD -- ':(exclude).scratch' '*.go'` -> 非测试调用点 **2 枚**（`resolve.go:198`、`resolve.go:340`），定义 `:242`；`resolveProbeRoot(os.TempDir())` 在 `:243`。
- ⇒ 这不是"零调用者"，而是"这两枚调用者都在 conformance probe 体内、真发 OS 路径不流经此处"——**该命题的形状不是调用计数**，我的尺给不出答案。
- **分档：判不动。**缺的读数＝`04` B-04（`:198`/`:340` 的宿主函数名 + 它们是否在非 probe 路径上被调）。

### P20 `HEAD:internal/ball/ball_windows.go:808-811`
> `// A host that rebinds while a card is waiting therefore has to call`
> `// TakeEscForCancel again on the next state sync - no caller does that today`
> `// because the config poll and the card live on different legs, and the residual`

- 断言符号：`Ball.TakeEscForCancel`（在 rebind 之后再调一次）＋ 理由分句"配置轮询与卡片不在同一条腿"
- 尺：`git grep -nE "TakeEscForCancel" …` -> 非测试调用点 `cmd/wisp/resident_approval_windows.go:877`、`cmd/balldebug/main.go:635`；`git grep -nE "RebindHotkeys" …` -> 非测试调用点 `internal/ball/hotkey_reload.go:96`。
- ⇒ "配置轮询与卡片不在同一条腿"这枚**理由分句**可疑：常驻腿今天同时握着两件事（`resident_windows.go:205 hotReload258 := func() ball.HotkeyConfig` / `:217 startResidentBall(…, hotReload258, …)` 与 `resident_approval_windows.go:877 b.TakeEscForCancel()`）。但我没读调用顺序与闭包同一性。
- **分档：判不动**（倾向过期）。缺的读数＝`04` B-05。

### P21 `HEAD:internal/ball/hotkey_windows.go:262`
> `// ReleaseEscAfterSession puts it back. The other three slots are never`
> `// standby: summon / mute / panel are modifier combinations by default and`
> `// are meant to be live the whole time the ball exists.`

- 这不是调用者计数断言，是**关于默认绑定 + 状态赋值**的断言（`HotkeyStandby` 从不被这三枚 slot 取到）。
- **分档：判不动。**需要的尺＝`HotkeyStandby` 的赋值点全谱 + `DefaultHotkeys()` 三枚绑定是否真为组合键（两把都不是我这把尺）。

### P29 `HEAD:internal/config/unwired.go:81`
> `missing: "nothing populates the bOverrides map risk.Gate reads: ticket 90 gave Gate a production call site (internal/tools readBlacklist, which feeds the card and the security log), but the confirmations that would fill this map are minted only by an L2 answer, and that flow is ticket 21's approval queue",`

- ⚠ 这一枚**自己就把"已过期的一半"写在句子里**（"ticket 90 gave Gate a production call site" 是正向、"nothing populates the bOverrides map" 是状态）。
- 尺：`git grep -nE "risk\.Gate\(" …` -> 非测试调用点 `internal/tools/mode.go:98 d := risk.Gate(c, overrides)`。**`overrides` 从哪来、生产里能不能非空，我没量。**
- **分档：判不动。**

### P30 `HEAD:internal/config/validate.go:112`
> `// approval.DefaultApprovalWarning (30s) whenever no producer passed a lead`
> `// (queue.go:88-89; ticket 267's census: no production caller passes one at`
> `// all). So any timeout at or below that lead does not fail loudly - it`

- 断言符号：`approval.Options.WarningLead` 的**生产者**
- 尺：`git grep -nE "WarningLead" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go' '*.go'` -> 非测试：`gate.go:35`（字段）、`gate.go:158`（`NewQueue(o.ApprovalTimeout, o.WarningLead, o.MaxPending, logf)`）、`queue.go:143-144`（getter）。⇒ 我**没有量 `WarningLead:` 这个组合字面量赋值形状**（那才是"有没有生产者传"的尺）。
- **分档：判不动**（倾向仍成立）。缺的读数＝`git grep -nE "WarningLead:" HEAD -- ':(exclude).scratch' '*.go'`。

---

## 丙 · 仍成立（15 枚）——按同一把尺复跑，今天还站得住

| # | HEAD 位置 | 断言（逐字关键行） | 符号 | 尺 -> 输出摘要 | 三层 / 分档 |
|---|---|---|---|---|---|
| P22 | `HEAD:cmd/wisp/config_readers_255.go:140`（同族 `:187/:188/:189/:190`） | `"memory":  hotClaimNoReader + "扫描零命中：nothing outside internal/config reads cfg.Memory - the L1/L3 knobs have no consumer"` | `cfg.Memory` / `cfg.Voice` 的跨包读者 | `git grep -nE "\.Memory\b\|\.Voice\b" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go' '*.go'` -> **剥掉 `internal/config` 自己之后只剩 `config_readers_255.go` 那 5 行（就是这几枚断言本体）**；`voice.*` 四枚同尺同数 | 包外 0 ／ **仍成立〔已证〕**，⚠ 射程＝"字段访问"这一把尺；换"构造参数名"那把尺读数可能不同 |
| P23 | `HEAD:cmd/wisp/panel_pump.go:80-82` | `// reparse/rewrite account live on the request path (RequestWorkspaceSwitch),` `// which the tree still has no caller for, so the snapshot carries what the` | `panel.RequestWorkspaceSwitch` | `git grep -nE "RequestWorkspaceSwitch" HEAD -- ':(exclude).scratch' '*.go'` -> 非测试只有 **定义 `internal/panel/workspace.go:111` + 5 处注释/文档行**（`bridge.go:17`、`composer_dispatch.go:74`、`git.go:268`、`instructions_200.go:112`、本行）；调用点全在 `_test.go` | 包外 0 ／ 入口不达 ／ **仍成立〔已证〕**（现在式 + 判对了）|
| P24 | `HEAD:cmd/wisp/resident_approval_windows.go:128` | `// the shipped process runs: this field has no production writer.` | `residentApproval.cancelKeyRead` | `git grep -nE "cancelKeyRead" HEAD -- ':(exclude).scratch' '*.go'` -> 赋值点**只在 `resident_cancel_key_label_260r4_windows_test.go:87/:101/:129` 与 `…_wording_260r3_windows_test.go:97/:119/:139/:162/:177`**；非测试只有读侧 `:169-170` | 包内读、无写 ／ **仍成立〔已证〕**（这是**设计成如此**的测试接缝，`:128` 明说 `nil is the production reader`）|
| P25 | `HEAD:internal/ball/d2d_windows.go:19` | `//   - SetOpacity/SetDpi/SetTransform-with-floats are never called; opacity` | `SetOpacity` / `SetDpi` / `SetTransform`（vtable 槽位） | `git grep -nE "SetOpacity\|SetDpi\|SetTransform" HEAD -- ':(exclude).scratch' '*.go'` -> **零枚调用形状**；命中只有槽位常量 `slotRTSetTransform = 30`（`:56`）与注释 `slotSolidSetColor = 8 // after ID2D1Brush (SetOpacity=4..GetTransform=7)`（`:69`） | 零调用 ／ **仍成立〔已证〕** |
| P26 | `HEAD:internal/config/unwired.go:69` | `missing: "risk.Facts.ShellString has no production writer, so R6's string-mode verdict never runs"` | `risk.Facts.ShellString` 的**写入点** | 尺＝写入形状 `git grep -nE "ShellString:\|ShellAllowlist:\|Allowlist:\|BlacklistOverrides:" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go' '*.go'` ＋ 非测试 `risk.Facts{` 字面量全谱：`cmd/wisp/panel_assets.go:73`（只填 `Declared/Paths/Irreversible`）、`internal/tools/fs_write.go:640/:663/:700`、`internal/tools/bridge.go:824`（后四枚只填 `Irreversible`） | **零枚写入点** ／ 读侧 `internal/risk/rules_shell.go:24` 恒 false ／ **仍成立〔已证〕** |
| P27 | `HEAD:internal/config/unwired.go:75` | `missing: "risk.Facts.ShellAllowlist has no production writer, so R6's allowlist verdict never runs"` | 同上，`ShellAllowlist` | 同尺：写入点 **0**；读侧 `rules_shell.go:41` | **仍成立〔已证〕** |
| P28 | `HEAD:internal/config/unwired.go:87` | `missing: "risk.NetTarget.Allowlist has no production writer, and no web/open tool exists to build one, so R5's domain gate never runs"` | `risk.NetTarget.Allowlist` | 同尺：`Allowlist:` 非测试写入 **0**；`assessor.go:112 Network *NetTarget` 字段在生产里也不见被填 | **仍成立〔已证〕** |
| P31 | `HEAD:internal/agent/spill.go:224-225` | `if s.judge == nil {` / `return "；注意：路径授权判定者未接线（fail-closed：C26 没接进来，这条路径是否还读得回来无法核实，按读不到处理）"` | `Spiller.judge`（`PointerJudge`）的生产注入 | 尺＝`git grep -nE "PointerJudge" HEAD -- ':(exclude).scratch' 'cmd/wisp/*.go'` -> **零命中**；`git grep -nE "NewSpiller\|WithPointerJudge" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go'` -> `cmd/wisp/run.go:1109 Spills: agent.NewSpiller(filepath.Join(rt.spec.dataDir, "artifacts"), loop.Budgets())`（**不带 judge**）、`internal/agent/loop.go:245 NewSpiller(…).WithPointerJudge(opt.Config.PointerJudge)`（该字段无人填） | 生产走**未接线那一支** ／ **仍成立〔已证〕**，⚠ 后果是**这句 fail-closed 文案每天在给模型看**，不是"以后会用到" |
| P32 | `HEAD:cmd/wisp/panel_pump.go:100` | `// is the hop 200-r1 left unconnected: the carrier field existed, no production` `// reader did, so the wire key could never arrive - …` | `agent.loop` 的项目说明读者 | `git grep -nE "setInstructionLoader" HEAD -- ':(exclude).scratch' '*.go'` -> 非测试调用点 **1**：`cmd/wisp/run.go:1070 rt.setInstructionLoader(instrLoader)`；读者 `panel_pump.go:96 instructionBundle()` | 入口到得了 ／ **仍成立（作为历史陈述）〔已证〕**，同文件自盖，读法风险＝低-中 |
| P33 | `HEAD:internal/panel/subagent_roster_197.go:9-11` | `// roster had no reader on the pump, and the three truncation accessors` `// (Truncated / ElidedRunes / DroppedKeys) had zero production callers, so` | `StreamLog.Truncated/ElidedRunes/DroppedKeys` | 尺：`git grep -nE "\.Truncated\(\)\|\.ElidedRunes\(\)\.DroppedKeys\(" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go'` -> **3 枚**：`cmd/wisp/panel_pump.go:203 state.StreamTruncated = rt.stream.Truncated()`、`:204 …ElidedRunes()`、`:205 …DroppedKeys()`（另 `:172 rt.stream.TruncationFor(key)`） | 包外 有、入口到得了（panel pump）／ **仍成立（过去式＋具名 `197-r2 §6.2`）〔已证〕** |
| P34 | `HEAD:cmd/wisp/providers.go:6` | `// Before this file llm.RunProbeSuite had nine green tests and no caller: the` | `llm.RunProbeSuite` | `git grep -nE "RunProbeSuite" HEAD -- ':(exclude).scratch' '*.go'` -> 非测试调用点 **1**：`cmd/wisp/providers.go:190 rep, err := llm.RunProbeSuite(ctx, prov, llm.ProbeSuiteOptions{`；定义 `internal/llm/probe_health.go:217` | 入口 = `wisp providers`/`probe` 子命令 ／ **仍成立（`Before this file` 框架）〔已证〕**；⚠ **"nine green tests" 那半句＝〔仅注释〕**——本程不跑 `go test`，枚数未复量 |
| P35 | `HEAD:cmd/wisp/run.go:11` 与 `:12-13` | `//	A11 llm.RunProbeSuite had no caller at all (ticket 11 AC#6)` ／ `//	A13 approval.Gate and tools.Bridge had zero importers outside their own` `//	packages, so no L1 window and no native card could ever open` | `llm.RunProbeSuite`／`approval.Gate`／`tools.Bridge` 的跨包 import | 同 P34 ＋ A4 尺：`internal/agent/approval` 非测试 importer **4**（见 P02） | **仍成立（过去式块 `run.go:5-8` `Before this file the three capabilities that landed this week were each green inside their own package and unreachable from production`）〔已证〕** |
| P36 | `HEAD:internal/tools/task.go:30` | `// promising task.cancel to the model while nothing registered it and` `// TaskRoster.Cancel had zero production callers. Registration and implementation` | `(*TaskRoster).Cancel` | 尺：`git grep -nE "\.Cancel\(\|func \(r \*TaskRoster\) Cancel" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go'` -> `TaskRoster` 族调用点**恰 1 枚** `internal/tools/task.go:740 stopped, why := t.d.Roster.Cancel(target)`；其余 15 枚 `.Cancel(` 是 `observe.Root`/context 根（同形异符号，剥掉） | **仍成立（过去式，票 221 的旧账）〔已证〕**；现行正向断言"ONLY one"见表三 |
| P37 | `HEAD:internal/tools/paths.go:25` | `// still had R2/R3 DORMANT because no production composition fed them. The` | `PathCanonicalizer` 作为 `risk` 缝的实现 | 尺同 P13：非测试 wire 点 3 枚（`bridge.go:212-213`、`bridge.go:925-927`、`panel_assets.go:66`）；类型定义 `paths.go:31`；生产构造 `cmd/wisp/run.go:501 rt.paths = tools.NewPathCanonicalizer(allowed, cfg.FS.ReparsePointExceptions)` | **仍成立（`Until this type existed` 框架）〔已证〕**，⚠ 句中 `still` 一词读法风险＝中（会被读成现状） |
| P38 | `HEAD:cmd/wisp/config_reload.go:7` | `// WHY THIS FILE EXISTS. Three of config.Manager's hooks had no production` `// assignment point before this file (measured at HEAD 7ffa9520, census` `// .scratch/wisp/probes/223/c1/census.md):` | 同 P07/P08/P09 | 同 P07-P09 尺 | **仍成立（过去式＋带 ref `7ffa9520`＋具名读数件）〔已证〕** —— 与 P07/P08/P09 的**唯一**差别就是这个框架句 |

---

## 表一·补 · `_test.go` 里同族的断言（4 枚，与产码分开报）

测试注释不是产码声明，但它会被派单当现状读，所以**单列不混**。

| # | HEAD 位置 | 逐字关键行 | 分档 | 依据 |
|---|---|---|---|---|
| T01 | `HEAD:internal/tools/task_output_canonicalize_fail_174_test.go:21-22` | `// Which arm production actually walks today is the unwired one` `// (cmd/wisp/run.go builds TaskDeps without Paths - 174-c2 §1); wiring is NOT` | **已过期** | `HEAD:cmd/wisp/run.go:544` 逐字 `for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: rt.tasks, Paths: rt.paths})` —— `Paths` **已经填**，`rt.paths` 在 `run.go:501` 由 `tools.NewPathCanonicalizer` 赋。⇒ `tools` 侧（`internal/tools/task.go:831` 的同款"路径授权判定者未接线"）走的是 wired 支；⚠ **`agent` 侧仍是 unwired**（P31，两枚不能压成一枚） |
| T02 | `HEAD:cmd/wisp/panel_inbound_33_test.go:190` | `// a production caller nobody dispatches is the same hole one level up.` | 判不动 | 依赖 P03/P05 的 `postMessage 接收器`那半枚未裁读数 |
| T03 | `HEAD:internal/panel/composer_dispatch_test.go:440` | `t.Errorf("no native host / WebView2 message channel is attached in this tree: no production source carries a " + …)` | 判不动（与产码冲突） | 与 `HEAD:cmd/wisp/panel_host_windows.go:821 return m.disp.Handle(ctx, raw)` **直接矛盾**。本程⛔ 不能跑仪器去裁谁错 —— 见 `04` B-07 |
| T04 | `HEAD:cmd/wisp/providers_test.go:4` | `// production caller, so the 「声明 PASS / 实测 FAIL」 event could never fire on a` | 仍成立（过去式） | 同 P34 |

---

## 表一·追 · 交件后补的 1 枚（W 尺，同分档口径，总数按 39 枚读）

### P39 `HEAD:cmd/wisp/run.go:333-334`（半枚过期）
> `// modeWrites is ticket 114 AC#2's gate: the one handler a panel mode request`
> `// is answered by, assembled with the SAME L2 leg the Store above was handed.`
> `// Nothing calls it yet - the WebView2 "event -> ParseComposerRequest" hop does`
> `// not exist in this tree (tickets 33/35) - so this is the gate's other half,`
> `// not a path: the host that gets wired later cannot assemble a wider档 without`
> `// a confirmation leg.`

- 两个分句，两把尺，**两档**：
  - **`Nothing calls it yet`（指 `rt.modeWrites` 这个字段）-> 仍成立〔已证〕**。尺：`git grep -nE "modeWrites" HEAD -- ':(exclude).scratch' '*.go'` -> 非测试命中 **4 行**：`cmd/wisp/panel_inbound.go:248 modeWrites := &panel.ModeWriteHandler{`（**局部变量，同名不同符号**）、`cmd/wisp/panel_inbound.go:272 Mode: modeWrites,`（喂的是那个局部量）、`cmd/wisp/run.go:337`（字段定义）、`cmd/wisp/run.go:674 rt.modeWrites = &panel.ModeWriteHandler{`（**赋值点**）。⇒ **`rt.modeWrites` 的读侧 = 0**。这是"建了但没接"档，注释说得对。
  - **`the WebView2 "event -> ParseComposerRequest" hop does not exist in this tree (tickets 33/35)` -> 已过期〔已证〕**。尺：`git grep -nE "ParseComposerRequest" HEAD -- ':(exclude).scratch' '*.go'` -> 非测试**调用形状 1 枚** = `internal/panel/composer_dispatch.go:155 req, err := ParseComposerRequest(raw)`（住在 `(*ComposerDispatch).Handle` 体内），而 `Handle` 的生产调用点是 P01 那 3 枚；`main.go:115 case "panel-inbound" -> cmdPanelInbound`；`cmd/wisp/panel_host_windows.go:403` 的注释还逐字承认 `panel.ParseComposerRequest inside Handle`。
- ⇒ **这一枚是本族最典型的害形**：**过期的是"理由分句"，不是"结论分句"**。结论（`rt.modeWrites` 没有读者）今天仍对，支撑它的理由（那一跳不存在）**已经不存在了**。后续程若照这句派单，会得出"所以本仓还没有入向通路"，而通路已经有了三处入口。

---

## 表一·附 · 一枚"产码侧无断言、票面侧有断言"的对照（`SealDir`）

`SealDir` 的"零生产调用者"这句话**不在产码注释里**（`git grep` 在 `internal/winsec/*.go` 非测试文件里搜不到该断言），**只在票 132 的票面**（见表二）。产码侧的**事实读数**在此具名：

- 尺：`git grep -nE "SealDir\(|func SealDir" HEAD -- ':(exclude).scratch' ':(exclude)*_test.go' '*.go'`
- 输出：**1 命中 = 定义行本身** `HEAD:internal/winsec/winsec.go:222 func SealDir(path string) error {`
- ⇒ **生产调用点 = 0（三层全 0：包内 0／包外 0／入口 0）** —— 票 132 那句读数以这把尺**仍成立**；票未 `-done`（`git ls-tree` 里文件名无 `-done` 后缀）。
