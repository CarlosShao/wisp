# 票 395 验收腿 `395-v1` —— D43#15 工具臂与 Acting 家族零生产派发修复（非实现者对抗验收）

> **尺**：被裁实现＝`2216f0d3d0ca6dd42fe2e14b57f15705a10289d4`（`feat(395-r1)`，6 件 **+490/−8**，
> 父＝`c9f6e097f`＝实现件申报开工锚，逐字对上了）；开工 `HEAD`＝`ea41fcab7`（≥派单锚，现跑
> `git rev-parse HEAD` 即此串）；`merge-base --is-ancestor 2216f0d3d HEAD` rc=0。
> 索引干净＝`git status --porcelain` 对**已跟踪文件**零改（满屏 `??` 为各腿既有 untracked 临时件）。
> ★本腿是本波唯一 cmd/wisp 腿，全程单发。
> **方法**：五格判据全部现跑（读数逐条本腿自取，⛔ 转述）；突变在仓外
> `D:\tmp\wisp\20261102\395v1\mut\`＋`-overlay`，共享工作树零突变窗口；还原核在**最后一发变异之后**。
> **staging（sherpa 三层陷阱，全部按票面环境段执行）**：`go test -c -o D:\tmp\wisp\20261102\395v1\wisp.test.exe
> ./cmd/wisp`（仓树编译 rc=0）＋三枚 DLL（`build\{onnxruntime,sherpa-onnx-c-api,sherpa-onnx-cxx-api}.dll`）
> 贴测试二进制同目录＋`WISP_SPIKE_MODELS/WISP_SILERO_MODEL/WISP_15E_TMP` 指主仓资产＋binary 直跑
> cwd=仓根（`voice15eRepoFile` 从 cwd 上爬）＋`-test.timeout=30m`＋GOTMPDIR 每跑置空。
> **零写面**：除本证据件外零改；⛔ thr/golden/frontend/design；⛔ 票面/台账；⛔ push。
> **临时件（只建不删）**：`D:\tmp\wisp\20261102\395v1\{wisp.test.exe,wisp.test.mut.exe,overlay.json,mut\,
> cases\,v1-run1.log,mut-run.log,v1-run2-restore.log,zeroreg-v1.log,zeroreg-384-single-v1.log}`。

## 0. 逐格裁决总表

| 格 | 判据 | 裁决 | 本腿现跑读数（出处 log） |
|---|---|---|---|
| ① | D43 状态机表零改 | **PASS** | commit 只碰 6 件（4 产码＋1 测试＋1 证据），`thresholds/golden/STATE_MACHINE/internal/**` 文件名零命中；`git diff 2216f0d3d..HEAD -- internal/statemachine internal/session internal/risk internal/speech internal/audio docs/PLAN.md docs/specs` **空** |
| ② | speak() HOP5a 真按 res.ToolCalls 分工具臂/非工具臂 | **PASS** | 事实链五环逐行读码核过（§2）；正控产线日志两枚冻结边真走（`v1-run1.log`） |
| ③ | 新测试真跑 | **PASS** | `--- PASS: Test395r1VoiceToolTurnWalksActing (11.81s)` rc=0（`v1-run1.log`；还原后再跑 `v1-run2-restore.log` 又 PASS 10.12s rc=0） |
| ④ | 反向判据（仓外 -overlay 摘工具臂 ⇒ 红 ⇒ 还原核） | **PASS** | 单点突变 diff＝1 插 1 删（`if req.hasToolCall {`→`if false {`）；`wisp.test.mut.exe` 同案 `--- FAIL` rc=1，三读数独立红（§4）；还原核后正控复跑绿 |
| ⑤ | 纪律格（git show --stat 与题面一致、票面未动） | **PASS** | 6 件 +490/−8 逐字对题面；`.scratch/wisp/issues/395-acting-family-dispatch-gap.md` 与 HEAD 逐字节同（diff 空）；台账/push 零触 |

## 1. 格① D43 表零改（现跑）

- `git show --name-status 2216f0d3d`＝`A cmd/wisp/resident_voice_acting_395r1_windows_test.go`＋
  `M cmd/wisp/{resident_voice_engines_windows,resident_voice_windows,resident_windows,run}.go`＋
  `A docs/evidence/s1/395r1-impl.md`，**无第 7 件**。
- `git diff --stat 2216f0d3d..HEAD -- internal/statemachine internal/session internal/risk
  internal/speech internal/audio docs/PLAN.md docs/specs` → **空输出 rc=0**（本 commit 之后也没人动过）。
- `internal/statemachine/table.go` 实质最后一次改动＝`af29b38d0`（D43 表驱动内核原落），其后仅
  `f34241380` gofumpt 空白 sweep（零语义）——与本票无关的两枚都在票前冻结。
- 对照冻结行（本腿现读 `table.go`）：
  - **#15 工具臂**（`:101-103`）：`Thinking + EvFirstToken → Acting`，Guard `first-token-has-tool-call`
    ＝ `f.HasToolCall`。
  - **#15 纯文本臂**（`:105-110`）：`Thinking + EvFirstToken → Speaking`，Guard `first-token-pure-text`
    ＝ `!f.HasToolCall`。
  - **#18**（`:128-130`）：`Acting + EvToolsCompleted → Speaking`，Guard `has-announcement`
    ＝ `f.HasAnnouncement`。
  - speak() 两臂所派 event/facts 与这三行**逐字对上**，零新行、零改行、零绕过（非法转移走
    `dispatch()` 的拒表路径，无 panic 无硬改状态）。**裁决 PASS。**

## 2. 格② 派发逻辑逐行读码（工具臂真按回合事实分岔）

事实链五环，每环本腿读的是**盘上 HEAD 工作树文件**（与 commit diff 逐字一致，restore check 另证）：

1. **事实来源（真数非推断）**：`internal/agent/loop.go:95` `Result.ToolCalls int`，由环路在工具执行/
   失败路径自增（`:775` `res.ToolCalls++` 等三处），＝已执行工具行数。
2. **转交**：`run.go` `execute()`（diff 现位 `+1512-1524`）：`rt.voiceTurnTools != nil` 时，
   `res.Status == StatusCompleted ⇒ rt.voiceTurnTools(res.ToolCalls)`，**否则转 0**（失败行保纯文本臂，
   失败件行为零变）；转交发生在 presenter/notify 呈现**之前**。`run.go:323-330` 新字段注释如实写明
   「the real res.ToolCalls, not an inference; nil when this process assembled no voice chain (CLI shape)」。
3. **挂线**：`resident_windows.go:582-585` boot 挂 `src.run.voiceTurnTools = rv.noteTurnToolCalls`，
   与 wire-r1 `voiceAnnounce` 同段兄弟行、同一 nil-guard（`src != nil && src.run != nil`）。
4. **消费一次**：`resident_voice_windows.go:411-416` 字段 `turnToolCalls int`（announceMu 守护）；
   `resident_voice_engines_windows.go:814-835` `noteTurnToolCalls`（写，nil-rv 早退）/
   `takeTurnToolCalls`（读后立即清零⇒纯文本回合永不吃上一回合残留、每回合一枚 announce 携带）；
   `announce()`（`:861-862`）`hasToolCall := rv.takeTurnToolCalls() > 0` 转上请求。
5. **分岔**：`speak()` HOP 5a（`:1011-1035`）：
   `if req.hasToolCall { dispatch(EvFirstToken,{HasToolCall:true,HasAnnouncement:true}) ⇒ 走通才
   dispatch(EvToolsCompleted,{HasAnnouncement:true}) } else { 改前单跳纯文本臂原样 }`。
   `dispatch()`（`resident_voice_windows.go:1394-1415`）返回真走通的边才续第二跳，表拒→noteProblem+false，
   后续播报失败处理（`newSpeaker` 错 ⇒ `moved` 时补派 #27 EvInterrupt）形制改前后一致。

**逻辑缺口排查**：工具臂只在 #15/#18 两枚冻结行上走；`to==from` 与表拒都不算 moved（不会假报）；
CLI 零回归面（voiceTurnTools nil 早退）。产线正控日志（`v1-run1.log`）真出现
`from=Thinking event=llm.first-token to=Acting why=播报：D43#15 首个 token 带 tool call…` 与
`from=Acting event=tools.completed to=Speaking why=工具完成：D43#18…`——读码与实跑互证。**裁决 PASS。**

## 3. 格③ 新测试真跑（staged 二进制，非 mock）

★**题面命令名勘误（如实报，非缺陷）**：题面 `-run 'TestResidentVoiceActing395r1'` 是**文件名前缀**，
盘上名册（`wisp.test.exe -test.list=395`，rc=0）唯一一案＝**`Test395r1VoiceToolTurnWalksActing`**；
按名册现跑。三层陷阱本腿各撞过一遍并如实录：`go test` 整包形态在包目录跑 binary ⇒ DLL 不贴身 ⇒
System32 onnxruntime 1.17.1 遮蔽（红）；什么都不给 ⇒ `0xc0000135`；只有**贴身＋env＋cwd=仓根**全形才绿。

正控读数（`v1-run1.log`，rc=0，`--- PASS (11.81s)`）：

```
D43_395 Sleeping+session.summon->Listening at_sample=13824
D43_395 Listening+asr.vad-stop->Thinking at_sample=74240
D43_395 Thinking+llm.first-token->Acting at_sample=74240     ← #15 工具臂（AC#1）
D43_395 Acting+tools.completed->Speaking at_sample=74240     ← #18 has-announcement（AC#1/#2）
D43_395 Speaking+tts.speak-done->Warm at_sample=74240        ← #26
BALL_STATES_395 [Listening Thinking Acting Speaking Warm]    ← Acting 真投到球（AC#2）
TRANSCRIPT_395 ["读一下我的笔记"]
TOOL_ROW_395 tool=fs.read level=L0 decision=allow outcome=success corr="…#call_16r1_demo"
```

- 与实现件 §2 记载逐位同形（仅时值差 11.81s vs 8.87s、corr 前缀不同＝每次新 task id，正常）。
- **非 mock 假报核**：回合经真 `runTextTask` 组合根＋C5 golden SSE＋真 `memory.Store` 回读 fs.read 行
  （`ListToolCallsByTask`）＋真 matcha TTS（合成 88934 样本真吐）；注入面只有 ASR 文本与渲染设备两枚，
  均为 `AGENTS §1.3`/wire-r1 既有许可接缝；**Acting 读数来自真 `res.ToolCalls`，不来自注入面**。
- 断言非摆设：三条 `demo16Walked` 独立钉 #15 工具臂/#18/#26 边，`ballStateSequence` 钉连续子序列
  `Thinking→Acting→Speaking→Warm`，tool 行钉 `outcome=success`——摘任一臂都红（§4 实证）。
  **裁决 PASS。**

## 4. 格④ 反向判据（仓外 overlay 摘工具臂 ⇒ 红 ⇒ 还原核）

- **突变最小性**：`git diff --no-index` 原件 vs `mut\resident_voice_engines_windows.go` ＝
  **1 insertion / 1 deletion**（`if req.hasToolCall {` → `if false {`），零其他改动。
- **overlay 构建**：`go test -c -overlay=overlay.json`（Replace 仓内绝对路径→mut 绝对路径）rc=0，
  产 `wisp.test.mut.exe`；仓内文件全程零碰。
- **红读数三枚**（`mut-run.log`，rc=1，`--- FAIL (10.23s)`），各自独立：
  1. timeline：`Thinking+llm.first-token->Speaking`（回 #15 纯文本臂，非 ->Acting）；
  2. `BALL_STATES_395 [Listening Thinking Speaking Warm]`（Acting 从球消失）；
  3. 点名断言红：`no D43#15 TOOL-arm walk Thinking+EvFirstToken->Acting: the turn carried a tool
     call but the announce still walked the pure-text arm`。
- **还原核（最后一发变异之后）**：仓内 `if false {` 命中 **0** 枚、`if req.hasToolCall {` 命中 **1** 枚；
  `git status --porcelain -- cmd/wisp .scratch/wisp/issues docs/evidence/s1/395r1-impl.md` **空**；
  不带 overlay 的正控二进制复跑 `--- PASS (10.12s)` rc=0（`v1-run2-restore.log`）。
  ⇒ 红由变异所致、非抖动、非共享树污染。**裁决 PASS。**

## 5. 格⑤ 纪律格

- `git show --stat 2216f0d3d`＝题面 6 件 +490/−8 逐字一致（4 产码 + 新测试 250 行 + 证据 146 行，
  行数与票面括注全对）。
- 票面 `.scratch/wisp/issues/395-acting-family-dispatch-gap.md`：`git diff HEAD` 空＋不在提交清单
  ⇒ 实现腿未翻框未改票面；本腿亦零触。台账 `pending-and-issues.md`、`docs/PLAN.md`、`docs/specs/**`
  不在改动面；无 push（实现腿与本腿都只 commit）。
- 共享编译面纪律：`cmd/wisp/**` 本波写腿唯一（派单自证），本腿只读＋仓外 overlay。
- **裁决 PASS。**

## 6. SPEC-10 §8 对抗验收七项表（本票射程内逐项落）

| # | §8 项 | 裁决 | 本腿读数 |
|---|---|---|---|
| 1 | stub 扫描 | **PASS** | 五改件 grep `todo!/unimplemented!/not implemented` 仅命中 `:339` **既有注释行文**（非本票新增、非 stub 形状）；无返回 nil 假成功、无 mock 顶真实现（§3 非 mock 核） |
| 2 | DEFERRED 双向对账 | **PASS** | `git show 2216f0d3d` 全文 `DEFERRED` **零命中**＝零新推迟标记＝登记面无需增行、无反向悬空 |
| 3 | SLO 测量完整性 | **PASS** | 六族合跑 `LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` **逐位同值**（`zeroreg-v1.log` 现跑，对照实现件 §4 与 43r1 §表先值）；thresholds/golden 零触（格①） |
| 4 | 契约偏离（D43 表逐条） | **PASS** | 格①/§2：两臂派发只落冻结行 #15(工具臂)/#15(纯文本)/#18/#26，guard 与 Facts 逐字对 `table.go:101-103/105-110/128-130`；D47/D34/D18 面本票零触 |
| 5 | 七条新增禁令静态扫 | **PASS** | 本 commit 新增行 grep：裸 `go func(`／`filepath.Clean·Abs`／明文 key／wall-clock 超时差值／镜像取哈希／面板侧 L2／内部工件受门控 Tool ＝**全零命中**（零新 goroutine、零新定时器，`noteTurnToolCalls` 只取 announceMu） |
| 6 | 场景↔票双向核对 | **PASS** | 票面 AC#1-5 逐条有读数落点（AC1=§3 边读数、AC2=ball+TOOL_ROW、AC3=§4 三红、AC4=§7 名册、AC5=§5）；缺口来源＝16-v2 终判升级发现（语音回合工具结果不走 Acting 呈现，S1/D43 场景链），追得到决策（D43 表内 #15 工具臂＋#18 本已冻结、只缺生产派发位） |
| 7 | 失败预演＋缺口审计 | **PASS** | 失败形状预演过：①非 completed 回合转 0⇒失败行保纯文本臂（run.go else 支，行为零变）；②表拒/`to==from`⇒moved=false 不假报、播报失败补 #27 形制不变；③CLI⇒缝 nil 早退；④残余如实：Acting 驻留为**收尾后按真事实回放**非工具执行实时（实现件 §6.1，`agent.EvToolStart` 全仓零生产者所致，写面外），16-r1 demo 基线未挂新缝仍自洽（§6.2，留编排裁）——两笔本腿核属实，**不构成本票 AC 失败**（AC#1/#2 判据原文只要求工具臂可达＋呈现读数） |

## 7. AC#4 零回归独立复跑（本腿名册，按名作差不看整包 rc）

staged 二进制六族合跑（`zeroreg-v1.log`，`-test.run=^(Test15e|Test26R3|Test41a2r1|Test381r2|Test43|
Test384|Test26P|Test392)`，rc=1 **因** 384 一枚，逐名差如下）：

| 族 | 本腿读数 | 对照实现件 |
|---|---|---|
| Test43 | **9/9 PASS** | 9/9 ✓ |
| Test15e | **6/6 PASS**＋LATENCY `262/640/902` 逐位同值 | 6/6 ✓ |
| Test41a2r1 | **5/5 PASS**（含两枚 KWS 真装载） | 5/5 ✓ |
| Test381r2 | **5/5 PASS**（含真设备 400ms 案） | 5/5 ✓ |
| Test26R3 | **3/3 PASS** | 3/3 ✓ |
| Test384 | 合跑 `Test384ResidentWatchdogEndToEnd` **红 1**：逐字 `never observed the tree back under the Armed cap … within 3s`＝实现件申报**同名同因**；**单发复跑绿** `--- PASS 1.15s`＋`384-READ settle failures=3 …`（`zeroreg-384-single-v1.log`），与 12-r2 §6 同族负载抖动；该案不经本票写面（本票零 goroutine 零定时器） | 申报一致 ✓ |
| Test26P/392 Path-C | **8 枚 SKIP**，逐字 `the AEC3 binding is not in this build (-tags webrtcap …)`＝本机 build 无 AEC3 自跳，非回归；392 纯 JSON 案 PASS | 同形 ✓ |
| 合计 | PASS=31＋384 单发补绿；FAIL 归因＝环境抖动非写面 | — |

## 8. 结论

五格判据（①D43 表零改／②派发逻辑逐行／③新测试真跑／④反向判据三红＋还原核／⑤纪律格）
**全 PASS**；SPEC-10 §8 七项**全 PASS**；六族零回归独立复跑与申报一致；实现件诚实边界两笔核属实、
不触 AC。本腿未发现任何 mock 假报、表绕过、放宽断言或纪律越界。**遗留（归编排者，非本票缺陷）**：
①题面 `-run 'TestResidentVoiceActing395r1'` 系文件名误写作测试名，建议后续派单以名册为准；
②§6.1/§6.2 两笔跨票边界（实时 Acting 事件缝、16-r1 demo 基线挂新缝）留编排裁。

> **是否同意结 395：同意结。**（`2216f0d3d` 可作为结票收件；推送归编排者，本腿⛔未 push。）
