# 220 AC#5 — 整包终态读数 ＋ 具名「这一格修完之后界面上还差哪一跳」（收口腿 `220-ac5`，2026-11-08）

- 角色：复核/收口腿（非实现者）。⛔ 零产码改动；本件只记读数与具名，AC 勾与 Status 留给编排者裁。
- 台面：`git clone --shared` 仓外导出树两棵——**改后**＝`/d/tmp/wisp/20261108/220ac5-tree`（HEAD＝`6e2eb5845`），**改前**＝`/d/tmp/wisp/20261108/220ac5-pre-tree`（锚＝`00d1e02e0`＝`ce693ce4^`，即票 220 第一笔 `ce693ce4` 落地之前的世界；`git merge-base --is-ancestor` 两枚都验过在 HEAD 祖先链上）。
- 建树补件（未跟踪运行件，主工作树手工拷入）：`build/*.dll` ×3 ＋ `third_party/sherpa-onnx/*.dll` ×3 ＋ `frontend/dist` ＋ **`third_party/spike-models/` ×8（第 8 件，本腿新学：缺它 `internal/speech` 整族 ~30 枚 0.00s 假红，红句自述 `silero model not found … looked for third_party\spike-models\silero_vad.onnx`；补上即整族转绿＝这把尺自证）**。
- 命令（两棵树同法，跑前 `unset GOTMPDIR TMPDIR TEMP TMP`，`export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，无 `CGO_ENABLED=0`）：
  - `go test ./internal/... -count=1 -timeout 40m` —— **两侧都跑到终态**（逐包 ok/FAIL 收口，零超时 panic）。
  - `go test ./cmd/wisp -count=1 -timeout 40m`（发1）／`-timeout 60m`（发2、发3）—— **三发全部被 `-timeout` 掐死，⛔ 终态**（详见 §二）。
  - 日志：`.scratch/wisp/probes/220/ac5-post-internal-retry.log`／`ac5-pre-internal.log`／`ac5-post-full.log`（发1）／`ac5-post-cmdwisp-retry.log`（发2）／`ac5-post-cmdwisp-retry2.log`（发3）。

## 一、`./internal/...` 终态红名册（逐名）与差集

**改前（`00d1e02e0`）红名 8 枚**（终态，rc=1，14:31:49 收口）：

| # | 红名 | 包 | 归因 |
|---|---|---|---|
| 1 | TestApprovalCardViewJSONKeysMatchFrontendTypes | internal/panel | 改前锚时代页面类型缺声明（红句自述 `Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`）＝票 329 面的"页面半"未落地形；改后绿 |
| 2 | TestComposerContractTypesMatchFrontend | internal/panel | 同上族（同句形）；改后绿 |
| 3 | TestComposerRenderFixtureTellsTheTruth | internal/panel | 渲染夹具面；改后仍红（两侧同红） |
| 4 | TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme | internal/panel | C21 色面；改后绿 |
| 5 | TestC21DesignTokensFourWayAgree | internal/panel | **已知既有红＝ball tokens 面（票 65 深色腿族），非本票面**（派单明令如实登记的这枚）；两侧同红 |
| 6 | TestExternalSamplerTracksSubjectGrowth | internal/proc | proc 面；改后绿 |
| 7 | TestResolvePerCallBudget | internal/risk | 性能读数形（改前锚时代 3.33s>阈值；220-v1 亦见过此形＝负载相关，⛔ 动 thresholds）；改后绿 |
| 8 | TestTicket252R1AlignmentAddsNoAuthorization | internal/tools | 252 授权面；两侧同红 |

**改后（`6e2eb5845`）红名 9 枚**（终态，rc=1，14:29:32 收口，71 秒——首跑 41 分钟是被同发 cmd/wisp 40 分钟并行负载饿出来的慢，机器空出来后同命令 71s）：

| # | 红名 | 包 | 归因 |
|---|---|---|---|
| 1 | TestOrbAssetShippedUnmodified | internal/ball | **改前锚无此用例**（65 jarvis 腿新增）＝orb 资产钉面，非本票面（红句：`embedded orb page sha256 = 4d9e291b…` ＝ 改后构建产物与钉值不一致；钉值本体在 `internal/ball/orb_test.go:225`，本件不转录） |
| 2 | TestOrbServingConstants | internal/ball | 同族（65 面） |
| 3 | TestJarvisAssetShippedAsBuilt | internal/ball | 同族（65 面） |
| 4 | TestComposerRenderFixtureTellsTheTruth | internal/panel | 渲染夹具面（两侧同红；红句＝夹具缺三个 block，`composer_test.go:873/902`） |
| 5 | TestC21DesignTokensFourWayAgree | internal/panel | 已知既有红＝ball tokens 面，非本票面 |
| 6 | Test52SchemaViolationsAreRejected | internal/plugin | **改前锚无此用例**（52 面新增）；非本票面 |
| 7 | Test52CommandKindValidates | internal/plugin | 同族 |
| 8 | TestCrossVetForLinux | internal/proc | 跨平台 vet 面（两侧同红） |
| 9 | TestTicket252R1AlignmentAddsNoAuthorization | internal/tools | 252 面（两侧同红） |

**差集（贴齐）**：
- **只在改前红 → 改后转绿 5 枚**：TestApprovalCardViewJSONKeysMatchFrontendTypes、TestComposerContractTypesMatchFrontend（＝页面半落地后契约尺合拢）、TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme、TestExternalSamplerTracksSubjectGrowth、TestResolvePerCallBudget（性能形，空载绿）。
- **只在改后红 → 改前不存在该用例 5 枚**：TestOrbAssetShippedUnmodified／TestOrbServingConstants／TestJarvisAssetShippedAsBuilt（65 jarvis 腿新增的资产钉面）、Test52SchemaViolationsAreRejected／Test52CommandKindValidates（52 面新增）。**`grep -rln "func <名>(" 改前树＝NOT-PRESENT 逐枚验过**，是"新用例自己的面"，不是改前绿改后红的回归。
- **两侧同红 4 枚**：TestComposerRenderFixtureTellsTheTruth、TestC21DesignTokensFourWayAgree、TestCrossVetForLinux、TestTicket252R1AlignmentAddsNoAuthorization。
- ★**判据面结论：两侧红名册里 220 面零枚**——`internal/panel` 的 220 判据用例（`subagent_blocked_220_test.go` 三枚含 `#序号` 两形、`subagent_roster_197_test.go`）与 `internal/agent/approval` 的 `window_read` 族（L1 枚举口）在改后终态里**全绿**（agent/approval 包 ok 6.586s；改前锚时代该包亦 ok）；改前锚不存在 220 用例（`ls internal/panel | grep -c 220`＝0）＝差集里没有"本票弄红/弄绿"的混淆项。

## 二、`./cmd/wisp` 整包：三发均⛔ 终态（如实记，非本票面）

| 发 | 时刻 | timeout | 结果 | 占死用例（panic dump `running tests:`） | panic 前红名 |
|---|---|---|---|---|---|
| 1（`ac5-post-full.log`） | 11:38:55 | 40m | `panic: test timed out after 40m0s`（2400.877s） | `TestAC294MuteHotkeyRegistersLiveOnARealBallWindow`（27m19s） | 9 枚（=发2 的 7 枚＋TestAC1AlwaysBranchDoesNotRevertAHandEditedKey＋TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card，后两枚 123s/121s 长时形在发2/3 绿＝负载形） |
| 2（`ac5-post-cmdwisp-retry.log`） | 12:21:34 | 60m | `panic: … 1h0m0s`（3600.143s） | `Test296BridgeKeepsThreeLiveKeysWithOwnerEmptySlots`（53m4s） | 7 枚 |
| 3（`ac5-post-cmdwisp-retry2.log`） | 13:27:45 | 60m ＋ **全新 `LOCALAPPDATA`**（球/面板 WebView2 UDF 全落新目录） | `panic: … 1h0m0s`（3600.149s） | `Test258BridgeMutationNoSrcKeepsOldBinding`（53m42s） | 7 枚（与发2 逐名同） |

- **稳定红 7 枚（发2＝发3 逐名同）**：Test313WritableUnionEqualsTheTwoRosters（313 门面）、TestAC247HandingTheLevelToTheSeamIsNotVisibility（247 面）、TestAC4FocusReturnToPriorWindowGap33r5（真窗前台面，10-16 已知六红族）、TestLockedSuffixNeverTakesTheLockItself33r11（33-r11 命名普查面）、TestTicket198FirstRunCallerIsTheRunEntryOnly（198 面）、TestTicket255RosterStillMatchesTheActualReadSites（255 面）、TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims（256 面）。**零枚属 220 判据文件**；归口编排者。
- **机理论证（为什么挂、挂在哪）**：三发三枚**不同**挂点，全部在球窗拉起同一条路上（`ball.New`→`createOnSTA`→`newOrbWindow`→WebView2 `Embed`）。发1 11:53 有一条被 `observe` recover 的 panic：`edge.(*Chromium).EnvironmentCompleted(..., 0x0)`——`res` 非 `S_OK` 失败会走 `log.Fatalf`（`forks/go-webview2/pkg/edge/chromium.go:172-174`），崩的是 **`res`=成功而 `env`=NULL** 时 `env.vtbl` 解引用（`:177`）；observe 收走 panic 后 STA 协程死亡，`ball.New` 等 done 永不回来。而球窗的 WebView2 user-data-dir 是**固定共享路径** `%LOCALAPPDATA%\wisp\orb`（`internal/ball/orb_windows.go:411-417`）——整包连跑几十枚球窗用例同挤一份浏览器态，前一枚的 msedgewebview2 未退净、后一枚拿回 `S_OK + NULL env` 即触发。**发3 换全新 `LOCALAPPDATA`（UDF 全新）仍挂**＝排除了"陈旧 UDF 脏态"假说，坐实是"整包连跑 churn 固定 UDF"这一运行形状本身。⇒ 判**本机台面系统性形，与票面无关**（与 10-16 的六红族同地界但更深一层：那次是活实例锁热键，这次是整包球窗自挤）。
- **改前侧 cmd/wisp 整包未跑**，理由如实写：挂点是台面形状（球窗 UDF churn）与本票面无关，发3 已带全新 UDF 仍挂，预树再烧一小时属foregone conclusion；本票判据面的两侧覆盖由 §一 `./internal/...` 终态＋§三补充读数承担。
- 环境旁证：跑时 `Get-Process wisp`＝0（无活实例，Q-85 要的世界成立）；42 枚 msedgewebview2 全属 10-25 起的系统壳进程（`--webview-exe-name=SearchHost.exe` 族），与本跑无关；console 会话（swq）Active 但 idle 8+15:43。

## 三、补充读数（明标：⛔ 不算终态——AC#5 原话"-run 单跑不算终态"，此发只为判据面补一枚绿证）

`go test ./cmd/wisp -run TestRunPacketLightsTheRosterRowOfItsLiveL1Window -count=1 -timeout 15m` ⇒ **`--- PASS (11.48s)`**。红句零、日志活字节＝L1 窗口行 `"blockedOnApproval":true` 且同包 `"pending":[]`（AC#2 那半在现 HEAD 的真 run 上仍点亮；红名册另证：cmd/wisp 三发 panic 前 220 用例无一红）。

## 四、门禁（改后导出树）

- `gofumpt -l cmd internal`＝**空**（`/d/work/base/gopath/bin/gofumpt.exe`）
- `go vet ./cmd/wisp ./internal/...`＝**rc0**
- `go build ./cmd/... ./internal/...`＝**rc0**
- `tools/d22scan`＝**clean**（`bans #1-5 internal/=387, cmd/=64, ban #8 internal/=963, cmd/=241, frontend/=119, design/=14 …`，rc0）

## 五、具名：「这一格修完之后界面上还差哪一跳」（界面地界，只写票面、⛔ 代做）

**现读（2026-11-08 现树，非旧结论）**——链路逐跳验过：

| 跳 | 现状 | 尺 |
|---|---|---|
| Go 侧算格 | ✅ 已通 | `internal/panel/subagent_roster_197.go:268`（`BlockedOnApproval` 双查）＋ L1 枚举口 `internal/agent/approval/window_read.go:63` → `cmd/wisp/panel_pump.go:128` reader → `cmd/wisp/run.go:1036`（`L1Windows: rt.liveL1WindowWaits`） |
| 传输跳 | ✅ 已通 | resident pump（`a5f84886f`）＋页面接收器 `window.wispOnPanelSnapshot`（`33b264d70`）；frontend 现 33 枚命中（src＋dist） |
| 页面类型 | ✅ 已声明 | `frontend/src/lib/panel.ts:213`（`TaskRowView.blockedOnApproval`）＋ `:414`（`PanelSnapshot.tasks?: TaskRosterSection`）；Go↔页面双向契约尺今日绿 |
| **渲染跳** | ⛔ **整节断的＝还差的那一跳** | **`snapshot.tasks` 在 `frontend/src` 零消费**（grep `snapshot\.tasks`＝0 枚消费点，只有类型声明）：`frontend/src/App.tsx:246` 任务监控 popover 喂的是写死的 `agents={[]}`；且 `MonitorAgentRow`（`frontend/src/components/harness/monitor-popover.tsx:28`）的状态联合只有 `"running"|"done"|"failed"`——**"等批准"这一态连载体都没有** |

**人话**：这一格修完之后，"哪个任务卡在等批准"到达页面时字节是全的（类型都对得上、今日契约尺绿），但**没有任何一块屏幕读它**——名册行到 `App` 为止就断了。界面今天对"卡在等批准"可见 **0 维**：L2 卡本体与「待审批 N」计数可见（`App.tsx:206-211`，读的是 `snapshot.pending`），但名册那一格没有任何像素。

**拆成两枚前端面小跳（归前端会话，本票只具名）**：
1. **取数跳**：把 `snapshot.tasks?.rows` 映射进任务监控 popover 的 `agents`（`App.tsx:246` 那枚写死空数组换成真读数）。
2. **画格跳**：行上把 `blockedOnApproval=true` 画成**附加的一维**（徽标/状态点/文案，与运行态**并存**）——票 220「外部对照」约束原文仍然有效：⛔ 做成"等批准覆盖还在跑"的单值优先级（OpenChamber 那枚谎不复制）；`MonitorAgentRow.state` 三态枚举怎么扩、双维怎么并排，属前端会话的设计地界，本票不替它拍。

**另一记（Go 侧零欠账确认）**：票 328 已把 L2 卡 corr 形状修到 join 认得（`splitClashSuffix`＝`subagent_roster_197.go:349`，末 `#` 后全 ASCII 数字的边界是 328 丙形加的）；`queue.go:170-172` 撞名改写形仍在、有判据。叠改端到端＝票 340 已结案。

## 六、账（交编排者）

1. **台面定式清单加一枚**：导出树建树七件套之外还需 `third_party/spike-models/`（第 8 件；缺它 speech 整族假红）。
2. **cmd/wisp 整包终态名册在本机今天造不出来**：三发三挂点取证如 §二；要终态需在球窗用例间独占 UDF（＝改测试/产码，本腿⛔）或另想台面。Q-85 的"无活实例"世界今天成立但不够——**新结论：阻塞不在活实例，在整包球窗自挤固定 UDF**。本票 AC#5 的 `./internal/...` 半已终态覆盖判据面；cmd/wisp 半以三发取证＋补充绿证如实交回，翻不翻勾由编排者裁。
3. cmd/wisp 稳定红 7 枚（255/313/198/256/33r11/247/AC4 各自票面）＋ `./internal/...` 改后新面红 5 枚（65 资产钉 ×3、52 ×2）＋两侧同红 4 枚（渲染夹具、C21、crossvet、252）——均非 220 面，逐名在册，归口编排者。
