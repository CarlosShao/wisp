# 65 — 常驻 leg（wisp.exe）快照读口接线 · 证据件

**腿名** resident-pump-wiring（wave77 C 腿，`cmd/wisp/**` 独占波次）
**日期** 2026-11-08 +08
**派单**：把常驻 leg 的面板快照读口接上（机主诊断：`wisp.exe` 面板快照不发 `ballStyle`/`modelSwitch`/`commands`/`plugins`/`run` 五枚 reader-gated 键；真机设置→外观 ball.style 行显示「缺省」）。

---

## 1. 复核结论：机主诊断「结论对、机制不准」，真根因共三层

派单要求 ⛔ 信诊断、现跑复核。复核读数（全部现跑）：

### 1.1 机主诊断的修正

- **"NewSnapshotPump 唯一生产调用在 run.go:998"** ✅ 实测成立（`grep -rn NewSnapshotPump --include=*.go cmd/ internal/ | grep -v _test` ＝ `cmd/wisp/run.go:998` ＋ `internal/panel/pump.go:313` 定义，仅两行）。
- **"常驻 leg 走 NewPanelManager，没有 pump 字面量"** —— **结论对、机制不准**。pump 不是 NewPanelManager 的职责：常驻 leg 的面板快照泵唯一来源是任务管线（`startResidentTaskSource` → `assembleRuntime`，run.go:998 字面量在装配根里，`cmd/wisp/resident_task_source_windows.go:290` 实调）。**出货形态走不到那里**：
  - `scripts/build.ps1` 出货 `wisp.exe` 带 `-H=windowsgui`（票 244；`resident_approval_windows.go:544` 逐字："A -H=windowsgui double-click has no console"）；
  - 无控制台 ⇒ `interactiveStdin()`（`approval_reply_stdin_windows.go:41`，GetConsoleMode 判据）返 nil ⇒ `startResidentTaskSource` 在**装配之前**就返 nil（`resident_task_source_windows.go:243-250`）⇒ **整个 runtime（连同 pump）从不装配** ⇒ 常驻进程零快照、零推送。
  - 即：不是"常驻 leg 建了 pump 但缺五枚 reader"，而是"**出货形态根本没有 pump**"。五枚 reader 在 run.go:998 那份共享字面量里**本来就全接了**（Commands :1041、BallStyle :1049、Run :1059、ModelSwitch :1071、Plugins :1104——Plugins 是 215a 在飞 WIP，本腿未触碰该文件）。
- **第二层（与泵无关的独立阻断）**：页面接收器缺失。`window.wispOnPanelSnapshot` 在 `frontend/src` **零命中**、在嵌入产物 `frontend/dist/assets/index-Dn9XNmNv.js` + `dist/index.html` **零命中**（实测该 bundle 的 `window.*` 赋值仅 `ErrorEvent`/`MotionIsMounted`，`wisp` 族仅 `wispBridge`/`wispDispatch`）。`.scratch/wisp/probes/329/r1/page-contract.md` §4 白纸黑字："页面接收器今天缺位＝现量……页面若未定义该全局，调用被 WebView2 静默丢弃——这正是今天的事实"。**机主看到的「快照里还没有球风格这一节（本机没有装配球风格读口）」＋「缺省」＝ App.tsx:315 `snapshot.ballStyle ?? null` 对 EMPTY 默认快照的诚实渲染——即便泵完美接线，屏幕在本波之后仍会显示「缺省」**。
- **第三层（触发缺失）**：`Show()` 从不推快照（`panel_host_windows.go:549` 无 publish 驱动；证据台账 13822 行早有登记"Show() 不推"）。冻结契约 `PLAN.md:2982`/`SPEC-08 §5.1`（"面板每次 show 必须发 panel.resync，由 Go 侧全量推状态"）此前无人兑现。有控制台的形状里泵虽在，但只在状态移动时 publish——开机后不跑任务就开设置页，页面照样一个包都没收到过。

### 1.2 真机日志佐证（机主运行形状）

`%APPDATA%\wisp-dev\logs\wisp-20261108-001.jsonl`（本腿实弹实例）：windowsgui 形态启动 → 发 Ctrl+Alt+P → show 那一刻：

```
10:07:12 [INFO] panel: SNAPSHOT at=2026-11-08T02:07:12Z depth=0 pending= mode=ask_every_step ws=unset results=0/0 done bytes=3656 sha256=1e4c49c7f08e4e0d
```

`mode=ask_every_step` ＝ 档位 reader 活着（链上 perm.Store 真读数）；`ws=unset`＝该形状无路径范围、诚实缺席。该实例数据根 `config.toml` `[ball] style='jarvis'`。

---

## 2. 设计裁定：甲（常驻 leg 自建泵），乙排除——理由

**裁定＝甲**：无管线形状由常驻 leg 自建 `panel.NewSnapshotPump`（`panel_resident_pump.go`），发布器在装配根按形状二选一；**乙（抽 `buildPumpSources` 两腿共用）排除**。理由（铁律 ⓪⓪⓪，判据＝最有价值＋主流 harness，非风险最小）：

1. **乙的前提不存在**：乙防的是"两处抄写"——实测两条 leg **共用同一份字面量**（console 形态常驻也走 assembleRuntime），仓库里本就没有第二份 PumpSources 可去重。为"第二个调用方"抽公共构造，第二个调用方却在另一条不经过 run.go 的装配路径上——抽出来的函数接不上（常驻 boot 时刻管线对象不存在：无 gate 注入对象之外的 runtime、无 runFeed、无 stream）。
2. **甲买到的东西乙买不到**：出货形态（windowsgui，机主真实使用形状）需要一个**不依赖任务管线的读口**，读这个 leg 真正持有的对象（ra.gate、链上 config Manager、命令名册、ingest 台账、链上 perm store）。这是主流 harness 的形状：**面板宿主自己持有面板读面**，而不是借用隔壁功能进程的装配残骸。
3. **乙现在物理落不了地**：run.go 载着 215a 腿的未提交 hunks（Plugins reader），本腿 commit 若含 run.go 会吞他腿 WIP（纪律 ⛔）——乙必须动 run.go，甲一个字节都不用。
4. **重复的代价被压到最小并具名**：Verdicts 的映射体抽成自由函数 `nativeVerdictsFromGate`（两腿共用，run.go 方法值引用零改动）；ModelSwitch/Plugins/BallStyle 的 config→panel 映射在 `panel_resident_pump.go` 落了**具名 helper**（`ballStyleSectionFromConfig`/`modelSwitchSectionFromConfig`/`pluginsSectionFromConfig`），run.go 的内联副本**在 run.go 解冻后应收敛到这三个 helper**（一行换调用，具名欠账，见 §6）。

**形状选择**：boot 里选**恰好一个**发布器——管线在（console 形态）→ `src.run.publishPanelSnapshot`（全量包，含 run/results/tasks，无包抖动）；管线不在（出货形态）→ 常驻泵（§3 的键面）。同一时刻一个进程只有一份发布器，无第二真相源。

---

## 3. 落地面（逐件）

| 件 | 改动 |
|---|---|
| `cmd/wisp/panel_resident_pump.go`（新） | 常驻泵：真源 readers（Verdicts←ra.gate、Mode←链上 perm.Store、Attachments/AttachmentReason←ingest 台账、Commands←commandRosterSection213、BallStyle/ModelSwitch/Plugins←链上 Manager）＋ Out＝329 跳原语（`panelSnapshotEvalJS` 逐字节 Eval＋与 run 台账同款的 `panelSnapshotSummary` 日志行）。**run/results/tasks/instructions 无源不冒充**（nil reader，键缺席）。 |
| `cmd/wisp/panel_pump.go` | `liveVerdicts` 映射体抽成 `nativeVerdictsFromGate`（自由函数，两腿共用；run.go:999 的方法值引用零改动）。 |
| `cmd/wisp/panel_inbound.go` | `newComposerDispatchChain` 返回值加宽 `(*panel.ComposerDispatch, *config.Manager, error)`——常驻泵读的 Manager **就是设置写入通路那一个**（链内单例，零第二装载，SPEC-03 §3.1 单真相规则）。CLI seam 丢弃该返回。 |
| `cmd/wisp/panel_resident_windows.go` | `newResidentComposerDispatch`/`newResidentPanelManager` 穿出 cfgMgr；`residentPanel` 加 `snapshotPublish` 钩子＋`setSnapshotPublish`/`publishOnShow`；`showOnThread` 成功 Show 后触发发布。 |
| `cmd/wisp/resident_windows.go` | boot 装配点：发布器二选一接上（管线泵 vs 常驻泵）；`panelHolder` 持有 publisher 并盖到每个新 spawn 的线程（319 重建形状不丢发布器）；分支选择行 slog＋stdout 双落。 |
| 测试件 | `panel_resident_pump_windows_test.go`（新，四牙见 §4）＋既有测试调用点随签名加宽机械更新：**四枚在库件**（`panel_attachment_214_test.go`、`panel_mode_confirm_114_windows_test.go` ×3）随本腿 commit；**两枚是他腿未跟踪在飞件**（`panel_resident_323_windows_test.go`、`panel_resident_assembly_323_test.go`，git 历史零记录＝323 腿未提交），本腿只在其盘上做了返回值加位（零断言改动），⛔ 不收入本腿 commit，随 323 腿自己的提交走（本腿 commit 后的库态不含这两枚，编译面自洽：新签名的全部在库调用方均已更新）。 |

**触发时机**：`publishOnShow`＝冻结契约 resync-on-show（PLAN.md:2982 / SPEC-08 §5.1）第一次有机制。状态移动型 publish 在无管线形状仍是具名非目标（该形状唯一可动的状态＝orb 卡，下次 show 即达页面）。

---

## 4. 判据（四牙，全绿）

`go test ./cmd/wisp -run 'Test65|TestBallStyle|TestResident' -count=1` ＝ **ok 0.207s**（含既有 65/Resident 全族）。新牙逐枚：

1. **Test65ResidentPumpCarriesConfigTruth**（①②）：走 boot 同款装配 `newResidentPanelManager(dir, nil)`，jarvis 经**真 setter**（`SetBallStyle`）在装配前写入，泵 publish 后**读推送字节**（recorder 捕 `window.wispOnPanelSnapshot(...)` 的逐字 JS）：
   - `ballStyle` 键在，`style=="jarvis"`（＝config 当前值），vocabulary 非空（schema 词表原样）；
   - `modelSwitch` 键在，Current 对上 `ChatSwitchOptions()` 真读数；
   - `commands`/`plugins` 键在；`plugins.engineState`＝`PluginEngineNotWired215`（诚实句原样）；
   - ④ `run`/`tasks`/`instructions` **键缺席**（注释写明：该 leg 无任务环/装载器/名册，假 reader＝常数冒充；`results` 不在此列——它非 reader-gated，泵恒发该键，空列表＝"没流过"的真话）。
2. **Test65ResidentPumpMutationStrippedBallStyleIsAbsent**（③突变）：取**真装配产出的 sources**、只摘 BallStyle reader → publish → 快照节与推送字节里 `ballStyle` 双双缺席（键确是 reader-gated，不是管道恒写）。
3. **Test65ResidentPumpNilManagerOmitsConfigKeys**（④）：cfg=nil（链未装配）→ 构造器**不挂** config 族 reader → 三键缺席（零值节比缺键更谎）。
4. **Test65ResidentShowTriggersResyncPins**（源钉）：`showOnThread` 必调 `rp.publishOnShow()`；boot 必含两分支（管线分支 `src.run.publishPanelSnapshot`、无管线分支 `newResidentPanelPump(ra.gate, cfgMgr65, rp.disp, attachmentDrafts, rp).publish`）；holder 必 `stampPublishLocked()`（319 重建不丢）。摘任一行即红。

---

## 5. 门禁读数（2026-11-08 +08 全部本腿现跑）

| 门禁 | 读数 |
|---|---|
| `gofumpt -l`（10 件点名，绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`） | 空，rc0 |
| `go vet ./cmd/wisp ./internal/panel` | rc0 |
| `go build ./cmd/wisp` | rc0（另按 build.ps1 同款 `-ldflags "-H=windowsgui"` 出 `build/wisp65pump.exe` 亦 rc0） |
| `go test ./cmd/wisp -run 'Test65\|TestBallStyle\|TestResident' -count=1 -timeout 30m` | **ok 0.207s** |
| `go test ./internal/panel -count=1` | rc1，**唯一红＝`TestC21DesignTokensFourWayAgree`**（ball tokens 面：`c21-native-tokens.md` 的 GlassRing/GlassHi/GlassHiSoft 在 `internal/ball/tokens.go` 缺——既有红，非本腿写面，如实登记） |
| `d22scan` | **clean**（internal 963＋cmd 243 等，rc0） |
| LF | 七枚产码/测试件逐一 CR=0（按字节扫） |

## 6. 真机验收（实弹）＋如实残余

- 构建 `build/wisp65pump.exe`（`-H=windowsgui`，与出货同旗标）→ 启动（PID 3684，数据根 `%APPDATA%\wisp-dev`）→ 键盘注入 Ctrl+Alt+P（全局热键，keybd_event 注入）→ 面板窗起（`wisp-panel-host`，rect 312,312-952,572）→ **日志即刻落 `panel: SNAPSHOT ... mode=ask_every_step bytes=3656 sha256=1e4c49c7f08e4e0d`**（实弹截图 11 张在 `.scratch/wisp/probes/65/panel-65-*.png`：面板初开/命令面板尝试/抽屉/设置页/外观页/外观滚动各态）。该数据根 `[ball] style='jarvis'`——show 时推的包携带的 Appearance 节即 jarvis（wire 形状由牙①在同构装配上钉死）。第二次 show（关窗再 Ctrl+Alt+P）再推一包（sha 随 generatedAt 变化，发布计数真在走）。
- **⛔ 如实残余：屏幕上的 `ball.style` 行本波之后仍显示「缺省」＋「快照里还没有球风格这一节（本机没有装配球风格读口）」。** 这不是 Go 半边的缺陷：包已在 show 时推到宿主（日志 sha 为证），但页面没有 `window.wispOnPanelSnapshot` 接收器（§1.1 第三层；`frontend/dist` 实测零命中；page-contract §4 交机主外部前端 agent 的那半）。Go→页最后一环要等前端接收器落地，届时本腿的包一到、行即显示当前值。**本腿零改 `frontend/**`、零重跑 npm build。**
- 实弹实例已停（PID 3684 kill）；`build/wisp65pump.exe` 留档可复跑。

## 7. 具名欠账（后续腿各一行可了结）

1. run.go:998 字面量的 BallStyle/ModelSwitch/Plugins 内联映射收敛到 `panel_resident_pump.go` 的三个 helper（run.go 现载 215a 腿未提交 hunks，本腿 ⛔ 未触碰；收敛＝三处换调用，零语义变化）。
2. 页面接收器 `window.wispOnPanelSnapshot`（前端面，page-contract.md §4 合同）。
3. 无管线形状的状态移动 publish（模式放松等）——本腿只兑现了 show 触发，具名非目标。

## 8. 边界自报

- 写面全部在派单清单内（`cmd/wisp/**`＋新测试件＋本证据件＋票 65 log 一行）；`internal/panel/**` 零触碰（无需改 pump.go——reader-gated 语义现成）；`frontend/**`、`internal/ball/**`、`internal/config/**`、`internal/tools/**`、`design/**` 零触碰。
- 未 commit 挪用他腿 WIP：`run.go`（215a Plugins reader）与 `approval_reply.go`（219 when= 戳）的在飞 hunks 原样未动，本腿 commit pathspec 严格点名避开。
- 签名加宽波及的既有测试调用点：在库四枚随本腿更新提交；他腿未跟踪两枚（323 腿在飞件）仅盘上加位、不收入本腿 commit（§3 表格）。
