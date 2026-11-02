# 票 395 落地腿 `395-r1` —— D43#15 工具臂与 Acting 家族（#17-20）零生产派发修复（写腿）

> **尺**：开工 `HEAD`＝`c9f6e097f0c38d8fcfa2beff98697b4bfc58fea7`（`dev`；现跑 `git rev-parse HEAD` 即此串，
> 索引对本腿 4 枚改件零已暂存）。★读数逐条现跑，⛔ 转述。
> **写面（5 枚，逐名）**：改 4 枚 —— `cmd/wisp/resident_voice_engines_windows.go`（+62/−8：`voiceAnnounceRequest.hasToolCall`
> 字段、`noteTurnToolCalls`/`takeTurnToolCalls` 两枚接缝方法、`announce()` 消费转交事实、`speak()` HOP 5a 工具臂分支）、
> `cmd/wisp/resident_voice_windows.go`（+5/−0：`residentVoice.turnToolCalls` 字段）、`cmd/wisp/run.go`（+22/−0：
> `agentRuntime.voiceTurnTools` 字段＋`execute()` 呈现前转交 `res.ToolCalls`）、`cmd/wisp/resident_windows.go`
> （+5/−0：boot 挂 `src.run.voiceTurnTools = rv.noteTurnToolCalls`，紧贴 wire-r1 的 `voiceAnnounce` 兄弟行）；
> 新 1 枚 `cmd/wisp/resident_voice_acting_395r1_windows_test.go`（端到端正控＋反向判据常驻牙）＋本文件。
> **零字节**：`internal/statemachine/**`（D43 表/events.go/timeouts.go，现跑 `git diff --stat` 空）、
> `internal/{session,risk,speech,audio}/**`（同上空）、thresholds/golden、`testdata/**`、`frontend/**`、`design/**`、
> 票面、台账、`docs/PLAN.md`、`docs/specs/**`。
> **临时件（只建不删）**：`D:\tmp\wisp\20261102\395r1\{wisp.test.exe,*.dll,cases,run1.log,zeroreg.log,zeroreg-384-single.log,mut\}`。
> 突变在仓外 `mut\`（叠加 `-overlay`），共享工作树零突变窗口。
> **staging**：`go test -c -o D:\tmp\wisp\20261102\395r1\wisp.test.exe ./cmd/wisp`（主仓树 build）＋三 DLL 同目录
> （`build\{onnxruntime,sherpa-onnx-c-api,sherpa-onnx-cxx-api}.dll`）＋ `WISP_SPIKE_MODELS/WISP_SILERO_MODEL/WISP_TTS_DIR/
> WISP_MODELS_MANIFEST/WISP_KWS_DIR/WISP_PUNC_MODEL/WISP_15E_TMP` ＋binary 直跑 cwd=仓根（15e-r1 清单同套）。

## 0. 一句话

产线语音回合的 announce 派发位（`resident_voice_engines_windows.go` `speak()`）此前恒以 `HasToolCall:false`
走 D43#15 纯文本臂，于是带 tool_call 的回合永远到不了 Acting；本腿把**真实回合的工具调用数**（`res.ToolCalls`）
经 `agentRuntime.voiceTurnTools → rv.noteTurnToolCalls → announceRequest.hasToolCall → speak()` 一路交进派发位：
`hasToolCall>0` 时走 #15 **工具臂**（Thinking→Acting）再由 #18 `EvToolsCompleted`(has-announcement) 到 Speaking，
否则保持 #15 纯文本臂（Thinking→Speaking）。⛔ D43 表/statemachine 一字未动，两臂都是冻结行。
正控实测：ball 走位 `[Listening Thinking Acting Speaking Warm]`、fs.read 行真落库为 success；
摘派发（仓外 overlay 变异）⇒ 回纯文本臂、Acting 读数消失、本案红；六族＋LATENCY `262/640/902` 逐位同值零回归。

## 1. ★ 必给读数一：机理复核＋派发位现核（第 16-v2 实锤对照）

- **机理（16-v2 逐字复核）**：demo 回合 SSE 首响应带 `tool_calls`（`voice_e2e_demo_16r1_windows_test.go:66-90`
  的 `demo16Golden`：round1=fs.read tool_call＋`finish_reason=tool_calls`，round2=纯文本 close）；而产线
  announce 派发位硬编码 `HasToolCall:false` ⇒ 带 tool_call 的回合恒落纯文本臂。本腿现跑 `git show HEAD:cmd/wisp/
  resident_voice_engines_windows.go` 复核：改前 `speak()` 内 `rv.dispatch(EvFirstToken, &Facts{HasToolCall:false,
  HasAnnouncement:true},...)` 单枚、无条件分支（原文 `:974-981`）。**确为机制所述派发位。**
- **改动落点（现跑行号）**：`resident_voice_engines_windows.go:1023` `if req.hasToolCall {`；`:1025` 派
  `EvFirstToken{HasToolCall:true,HasAnnouncement:true}`（#15 工具臂 Thinking→Acting，guard `first-token-has-tool-call`）；
  `:1027` 派 `EvToolsCompleted{HasAnnouncement:true}`（#18 has-announcement，Acting→Speaking）；`else` 分支保持
  `EvFirstToken{HasToolCall:false,...}`（#15 纯文本臂）。**对照 D43 冻结行（`table.go:100-134`）逐条对上，#15/#18
  两臂皆表内既有行，零新行、零改表。**
- **事实来源（真，不是推断）**：`res.ToolCalls` 是环路自己统计的已执行工具数；多轮里只有首响应带 tool_call 才会
  进工具轮（首响应纯文本即收尾），故 `res.ToolCalls>0 ⟺ 首响应带 tool_call`，与 #15 的 `first-token-has-tool-call`
  守卫同真值。转交发生在 `run.go:1519-1524`（`execute()` 呈现前）：`completed` 才转交 `res.ToolCalls`，失败行转 0
  （保纯文本臂，失败件行为零变）。`resident_windows.go:585` boot 挂线，与 wire-r1 `voiceAnnounce` 同段兄弟行。

## 2. ★ 必给读数二：AC#1＋AC#2 端到端实测（语音回合带 tool_calls ⇒ Acting，对照 16-r1 demo 形状）

**正控跑（`run1.log`，`-test.run=^Test395r1VoiceToolTurnWalksActing$`，rc=0，`--- PASS` 8.87s）**，逐条读数：

```
D43_395 Sleeping+session.summon->Listening at_sample=13824
D43_395 Listening+asr.vad-stop->Thinking at_sample=74240
D43_395 Thinking+llm.first-token->Acting at_sample=74240        ← #15 工具臂（AC#1）
D43_395 Acting+tools.completed->Speaking at_sample=74240        ← #18 has-announcement（AC#1/#2）
D43_395 Speaking+tts.speak-done->Warm at_sample=74240           ← #26
BALL_STATES_395 [Listening Thinking Acting Speaking Warm]        ← Acting 真投到球（AC#2）
TRANSCRIPT_395 ["读一下我的笔记"]
TOOL_ROW_395 tool=fs.read level=L0 decision=allow outcome=success corr="673ff6de-...#call_16r1_demo"
--- PASS: Test395r1VoiceToolTurnWalksActing (8.87s)
```

- **对照 16-r1 demo 形状**：16-r1 demo 的走位是 `Armed(#5)→Listening(#7)→Thinking(#11)→Speaking(#15 纯文本)→Warm(#26)`，
  **Acting 缺席**（16-v2 明列的缺口）。本案复刻同一「首响应带 tool_call、末轮纯文本 close」的回合形状（`demo16Golden`
  ＋`demo16Reply` 逐字复用），实测走位把它补成 `…Thinking→Acting→Speaking→Warm`：**#15 工具臂与 #18 都是产线
  `speak()` 亲自派发，ball 亲自投影，fs.read 行亲自落库**（不是脚本摆出来的闭合）。本案不武装 KWS，故入口是
  `Sleeping+session.summon→Listening`（#4）而非 demo 的 `Armed(#5)→Listening(#7)`；这是本案与 demo 唯一不同的前置腿。
- **端到端真实性（非 mock 假报）**：回合经真 `runTextTask`（`wisp run` 组合根）＋ C5 golden SSE（`internal/llm/golden`）
  ＋ 真冒烟 fs.read（真 `memory.Store` 落 `tool_call` 行、`outcome=success`）＋ 真 matcha TTS 合成（`voiceWireTtsFiles`）
  ＋ 真 announce 走位。注入面仅 `C8/ASR 文本`（`voiceWireAsr`，wire-r1 已用的导出接缝）与 `渲染设备`（`voiceWireSpeaker`
  记录器，无实卡），均为 `AGENTS §1.3` 许可接缝；Acting 读数不来自注入面，来自真回合的真 `res.ToolCalls`。

## 3. ★ 必给读数三：AC#3 反向判据（摘派发 ⇒ 回纯文本臂 ⇒ 红；三读数＋还原核）

**变异（仓外，`-overlay`，共享树零改）**：把 `resident_voice_engines_windows.go` 复制到
`D:\tmp\wisp\20261102\395r1\mut\`，仅把 `if req.hasToolCall {` 改成 `if false {`（＝摘掉 #15 工具臂派发），
造 `overlay.json`（`{"Replace":{<仓内绝对路径>:<mut 绝对路径>}}`），`go test -c -overlay=overlay.json` 重建跑同一案。
**变异跑（`mut\run-mut.log`，rc=1，`--- FAIL` 8.93s）三读数**：

```
D43_395 Thinking+llm.first-token->Speaking …                   ← 现回 #15 纯文本臂（对照正控的 ->Acting）
BALL_STATES_395 [Listening Thinking Speaking Warm]             ← Acting 从 ball 消失
no D43#15 TOOL-arm walk Thinking+EvFirstToken->Acting: the turn carried a tool call
   but the announce still walked the pure-text arm: […{From:Thinking To:Speaking Event:llm.first-token Say:播报：D43#15 首个纯文本 token…}]
--- FAIL: Test395r1VoiceToolTurnWalksActing (8.93s)
```

三条读数各自独立地红：①timeline 的 #15 边落到 `Speaking` 而非 `Acting`；②ball 状态序列无 `Acting`；
③常驻断言点名「回合带了 tool_call 但走位仍纯文本臂」。**还原核**：变异件全程在仓外 `mut\`；现跑核仓内
`resident_voice_engines_windows.go` 仍含唯一 `if req.hasToolCall {`（1 枚）、无 `AC#3 MUTATION` 串（0 枚），
`git diff` 与本文件 §1 所述改动逐字一致；不带 overlay 的正控二进制 `run1.log` 仍 `--- PASS`（承 §2）。
⇒ 红是变异所致，不是抖动，也不是共享树被污染。

## 4. ★ 必给读数四：AC#4 零回归（六族＋Path-C＋LATENCY 262/640/902）

**六族合跑（`zeroreg.log`，`-test.run=^(Test15e|Test26R3|Test41a2r1|Test381r2|Test43|Test384|Test26P|Test392)`）**：

| 面 | 结果 | 读数 |
|---|---|---|
| Test43（9 案） | **9/9 PASS** | autostart/crash/lifecycle 全绿（0.00–0.40s） |
| 15e（6 案） | **6/6 PASS**；**`LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` 逐位同值** | `Test15eSpeechEndToEnd…` 10.79s 等；LATENCY 行逐字同 `262/640/902` |
| 41a2r1（5 案） | **5/5 PASS** | 含 21.37s/16.97s 两枚 KWS 真装载案 |
| 381r2（5 案） | **5/5 PASS** | 含 `Test381r2RealDevicePlayingPhaseAbortWithin400ms` 35.49s 真设备案 |
| 26R3（3 案） | **3/3 PASS** | 含 17.09s 真装载案 |
| 384（3 案） | **单发 3/3 PASS**（E2E 1.15s） | 合跑首遍 `Test384ResidentWatchdogEndToEnd` 红 1 枚，见下「抖动申报」 |
| 邻接包 | `internal/{statemachine,session,observe,config,watchdog,proc}` **全 ok（rc=0）** | — |

- **抖动申报（非零回归差）**：六族合跑首遍 `Test384ResidentWatchdogEndToEnd` 红 1 枚，逐字
  `resident_watchdog_windows_test.go:226: 384: never observed the tree back under the Armed cap (a real within-cap
  reading) within 3s`——**纯负载敏感读数**（串跑前刚跑完 381r2 真设备 35.49s＋26R3 17.09s 两枚重案）。**单发复跑
  `zeroreg-384-single.log` 绿：`--- PASS` 1.15s，`384-READ settle failures=3 (BackWithinCapMS=32…)`**。与 12-r2 §6
  记录的「384 串跑装载抖动、单发绿」同族；该案不经本腿任何写面（本腿只在 announce 派发位加分支、在 runtime 加
  一枚 `func(int)` 字段，零 goroutine、零定时器）。
- **Path-C 套件（26P 5 案＋392 4 案）＝`SKIP`，非回归**：逐字跳过原因 `the AEC3 binding is not in this build
  (-tags webrtcap + the deps.toml CGO flags)…webrtc-audio-processing binding not built`（26p_echo:109/250/364、
  26p:334/512、392:404/554/647）；本机 build 无 AEC3 绑定，案自跳。392 的纯 JSON 案 `Test392JSONSentinelClampsNonFiniteReadings`
  **PASS**。这些案不经本腿写面，且其 announce 调用在跳前逐条走出纯文本臂（log 中 `Thinking+llm.first-token->Speaking
  …播报：D43#15 首个纯文本 token…`）＝默认 `turnToolCalls=0` 的向后一致。
- **静态四尺**：`gofmt -l cmd/wisp`（空）、`gofumpt -l`（绝对路径，空）、`go vet ./cmd/wisp` rc=0、
  `go test -c ./cmd/wisp` rc=0（测试文件含本案）。

## 5. ★ AC#5 纪律格

- ⛔ **D43 表／statemachine**：`git diff --stat -- internal/statemachine internal/{session,risk,speech,audio}` **空**；
  两臂＃15/＃18 逐条取自冻结表（`table.go:100-134`），零新行零改行。
- ⛔ **翻框／改票面／台账**：票面 `.scratch/wisp/issues/395-*.md`、`docs/reports/pending-and-issues.md`、
  `docs/PLAN.md`、`docs/specs/**` 一字未动（本文件之外零写）。
- ⛔ **push**：本腿只在本工作树产出改动与文件，未执行任何 `git push`（推送归编排者）。
- **临时件只建不删**；变异在仓外 overlay，共享工作树零突变窗口；`git add`/`commit` 走显式 pathspec（见提交）。

## 6. 诚实边界／具名残余（⛔ 藏）

1. **Acting 的时相边界（最重要的一条，如实录）**：本腿在 announce 交接点用**真回合**的 `res.ToolCalls` 派发
   #15 工具臂＋#18，ball 真进 Acting 且真投影；但 `speak()` 在回合收尾后才跑，故 Acting 的**驻留不与工具实际执行
   同刻**（是「按真事实回放两行合法边」，不是「工具运行的实时绿球」）。原因在接缝：语音腿没有工具事件的实时通道
   （`agent.EvToolStart` 全仓零生产者，只有 `EvToolEnd` 由环路 sink 收），而本票写面只 `cmd/wisp/**`。**这不是 mock
   假报**——`res.ToolCalls` 是真行数、#15/#18 是冻结行、fs.read 行真落库——但它**不是**「工具执行期实时 Acting」。
   若日后要实时，可把 `consoleSink` 的首个 `EvToolEnd` 经 owner 队列投一条 #15 工具臂；那要动 run.go 的 sink 接线且
   与 #11 有先后序风险，超出本票「announce 派发位」派单，留给票面/编排者裁定。
2. **16-r1 demo 测试件未改（跨票边界）**：`voice_e2e_demo_16r1_windows_test.go` 的 `onRuntime` 只挂
   `voiceAnnounce`（`:244-247`），未挂新缝 `voiceTurnTools`，故其 `demo16Walked(Thinking,EvFirstToken,Speaking)`
   （`:314`）仍是改前的纯文本臂基线、按默认 `turnToolCalls=0` 自洽通过；但它已不反映「boot 挂了新缝后的产线走位」。
   把该 demo 更新为挂新缝并断言 Acting 是**改另一张票的已收证据**，本腿按纪律不擅改；此取舍留编排者裁。
3. **Path-C 套件在本机 build 无 AEC3 绑定 ⇒ 自跳**（§4 逐字）。带 `-tags webrtcap` ＋ deps 的重建不在本腿 staging。
4. **`res.ToolCalls` 作为 #15 守卫真值**：多轮下等价于「首响应带 tool_call」（首响应纯文本即收尾，进不了工具轮）；
   单轮工具即首响应。此等价是推导非假设，但**不是逐 token 观测**的 first-token 事实；与 §6.1 同源（缺实时事件缝）。
