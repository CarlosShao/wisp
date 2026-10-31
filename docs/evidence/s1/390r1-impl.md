# 票 390 落地腿 `390-r1` —— 球单击打断装配（写腿）

> **身份**：写腿（实现者）自报读数。判定归验收/重裁腿；本腿**零翻框、零改票面、零台账、零 push**。
> **尺**：开工 `HEAD=b109a3010`（`dev`；staged 索引 0 行、`cmd/wisp`／`internal/ball` porcelain 0 行）。共享树他腿零碰。现跑 2026-10-31 18:3x–19:1x +08。
> **staging**：`go test -c -o D:\tmp\wisp\20261031\390r1\bin\wisp.test.exe ./cmd/wisp`＋四 DLL 拷同目录（**含 onnxruntime.dll**，A1326 清单）＋`WISP_SILERO_MODEL/WISP_SPIKE_MODELS/WISP_PUNC_MODEL/WISP_TTS_DIR/WISP_15E_TMP` 五 env（26r3 §7 同款）＋`-test.run` 只定向、`-timeout 30m`、GOTMPDIR 置空。
> **临时件（只建不删）**：`D:\tmp\wisp\20261031\390r1\{bin,mut,head,headtree,export,cases,logs}`。**⛔ push；⛔ 票面/台账/specs/PLAN/testdata/frontend/design/scripts/.github 零字节；⛔ `internal/{ball,risk,session,speech}` 零字节；⛔ 381-r2 打断本体零字节。**

## 0 一句话

381 残余①收口：球单击——interruptCh 注释里点名的第二个触发源、381-r1 自认「预留钩子但未装配执行器」——现在真挂上 `rv.TriggerInterrupt`：**播报中球单击 ⇒ 真机五轮 22–39ms 真停（≤400ms MET）＋门重开**；空闲单击惰性（381-r2 idle 正控同形）；未挂接形状单击照播完（挂接敏感性的测试内反形）；摘挂接突变主案红、Esc 通路照常绿（定向自证）、原树还原绿。零回归按红名集合比零新增（六族 41/42，唯一红 `Test384ResidentWatchdogEndToEnd` 经 HEAD 对照树同样红＝既有负载时序形状，非本腿）。

## 1 ★ 三件结构（逐件：接在哪／为什么／读数）

| 件 | 接在哪 | 为什么是这个形状 |
|---|---|---|
| **①执行器与手势入口** | `resident_ball_windows.go`：`residentBall.voiceInterrupt func()` 字段（`muteMux` 同锁组）＋`attachVoiceInterrupt`／`currentVoiceInterrupt`（唯一写者／锁下读回）＋`clickInterruptGesture()`（有执行器→`trigger()`＋记录；无→`recordBallGesture("click")`＝pre-390 形状）＋`Events.OnClickBall` 闭包改为 `rb.clickInterruptGesture()` | `ball.Events.OnClickBall` **本来就是** internal/ball 的单击事件面（`ball_windows.go`：WM_LBUTTONUP 非拖拽非 dock 时 `b.fire(b.opts.Events.OnClickBall)`）——381-r1 说的「需 internal/ball API 扩展」实际不需要，API 已够，**冻结面零字节**。执行器晚挂沿用本文件自己的先例（`attachMuteGate`/`attachWakePause`，同 `muteMux`）：球先建、语音链后建，闭包构造时拿不到 rv。回调跑在 ui-sta 线程（`Ball.fire` 内联，契约 quick and non-blocking）——**与 Esc cancel key 同一条线程**，381-r2 对 TriggerInterrupt 的线程假设（「runs on the caller's goroutine (the Esc cancel key fires on the ball's UI thread)」）单击完全满足 |
| **②装配行** | `resident_voice_windows.go` `startResidentVoice` 的 `if rb != nil` 块：`rb.attachVoiceInterrupt(rv.TriggerInterrupt)`，紧挨 15e 的 `attachWakePause` | 同一个「这条腿拥有通道 ⇒ 由它把执行器交出去」的所有权理由：interruptCh 是语音链的（381-r2 本体），单击只是第二触发源。装配顺序与生产一致（rv 建好后 attach），与 `attachWakePause` 同一块同形状。**381-r2 本体（TriggerInterrupt/speakAbort/speak 入口清 token）零字节** |
| **③文案随行** | `ballGestureWhy` 常量＋未寄宿手势 bullet＋boot log gestures 行：单击从「无执行者」名单移出（interrupt 半归 390），**summon 半仍在名单**（把单击变任务仍无执行者，不属于本票） | ticket 290 的教训：解释手势为什么没执行的理由必须始终为真。单击现在有 interrupt 执行器（条件＝本进程装配了语音链），但其 summon 半仍无——文案逐处改真 |

## 2 ★★★ ≤400ms 实测（真 waveOut 设备＋真 matcha 合成；触发走球单击手势入口，非裸 TriggerInterrupt；`Test390r1BallClickRealDeviceWithin400ms`）

句形与 381-r2 同（30 字生产句、4 chunks ≈6.4s 音频）。每轮：播报→等 Drain 停入（parked）→**300ms 后球单击**（`ball.clickInterruptGesture()`）→计 click→停。触发链＝`OnClickBall` 闭包同体的 `clickInterruptGesture` → `currentVoiceInterrupt` → `rv.TriggerInterrupt` → `speaker.Abort`。

| 轮 | unstick（click→Drain 作答＋#27） | reset（click→`waveOutReset` 落地） | 对照（无单击） |
|---|---|---|---|
| 记录轮 1 | **33ms** | 33ms | control `parked_to_done=2.096s`（照播完） |
| 记录轮 2 | **39ms** | 39ms | |
| 记录轮 3 | **22ms** | 22ms | |
| 记录轮 4 | **31ms** | 31ms | |
| 记录轮 5 | **39ms** | 39ms | |

**≤400ms：MET**（判格归验收）。五轮分布 22–39ms，与 381-r2 的 Esc 五轮（41–46ms）同量级——同一 abort 半边、第二个触发源。`abort_returned` 全 0（381-r2 已具名的 Close/Abort 抢锁读数波动，非断言）。控制轮同机现测。

## 3 ★★ 正控（「无事发生」的那一半）

- **非播报期球单击 ⇒ 无事发生**：`Test390r1IdleBallClickIsInert` PASS——挂接后空闲单击：机器 Sleeping 不动、timeline 空、problems 空、`ttsAborted/ttsSpoken=0`（381-r2 idle Esc 的同形读数，经球手势入口）；nil host／未挂接 host 单击不 panic。
- **未挂接形状**：`Test390r1UnattachedClickDoesNotInterrupt` PASS——`UNATTACHED spoken=1 aborted=0 problems=0`：pre-390 形状（无执行器→recordBallGesture）播报照播完。与主案的唯一 delta＝`attachVoiceInterrupt` 那一行，这就是挂接敏感性的测试内证明。
- **静音/空文本/良性句**：`TestWireR1SilenceAndEmptyTextNeverSpeak`、`Test26R3SilenceNeverTouchesTheGateNeedle`、`Test381r2VetoSubstringFinalsReachTheAgent` 在合跑中复绿——本腿零新增代码路径被它们触碰。

## 4 ★★★ 反向判据（仓外 `D:\tmp\wisp\20261031\390r1\mut\` ＋ `-overlay`）＋还原核

**M1 突变＝摘球单击挂接**：`clickInterruptGesture` 的 `trigger()` 调用摘除（attach／记录保留）。突变件 SHA256 `D6C8C44B…` ≠ 原件 `E229F9DA…`（互异）。

| 读数 | 发 | rc | 结果 |
|---|---|---|---|
| **摘挂接 ⇒ 球单击不再打断 ⇒ 红** | `Test390r1BallClickPlayingInterruptAbortsAndWalks27` @ M1 二进制 | **1** | 红句逐字：`MUTANT 390r1-M1 SHAPE: the ball click did not stop the announcement in flight (spoken=1, aborted=0): the click gesture did not fire the attached TriggerInterrupt, so the sentence played out - the second trigger source is not wired (ticket 390: 球单击→TriggerInterrupt 真挂)` |
| **定向自证（Esc 通路照常）** | `Test381r2PlayingPhaseInterruptAbortsAndWalks27`＋`Test381r2IdleInterruptIsInert`＋`Test390r1UnattachedClickDoesNotInterrupt` @ 同 M1 二进制 | **0** | 3/3 PASS；`ABORT-27 aborted=1 spoken=0 … gate_closed=2 gate_reopened=2`——摘的是球单击这一跳，interrupt 本体与 Esc 半边分毫未动 |
| **还原核** | 无 overlay 原二进制重跑主案 | **0** | PASS，`CLICK-27 aborted=1 spoken=0 next_spoken=1 probe1_aborts=1 probe2_aborts=0 gate_closed=2 gate_reopened=2` 同读数；`git status --porcelain -- cmd/wisp` 始终恰本腿 3 枚（overlay 只读，工作树零突变窗口） |

## 5 ★★ 零回归＋四数＋两尺＋d22scan

| 尺 | 命令 | 读数 |
|---|---|---|
| **六族合跑**（15e/26R3/41a2r1/381r2/43/384＋本腿 4 案） | staged `wisp.test.exe "-test.run=TestWireR1\|Test15e\|Test26R3\|Test41a2r1\|Test381r2\|Test43\|Test384\|Test390r1" -test.v -test.timeout=30m` ×2 | 两发 **42 RUN／41 PASS／1 FAIL／0 SKIP**（本腿 4 案全绿；`CLICK-27`、`CLICK-CONTROL parked_to_done=2.096s`、五轮 22–39ms、`UNATTACHED` 全在）；唯一 FAIL＝`Test384ResidentWatchdogEndToEnd`（`never observed the tree back under the Armed cap within 3s`） |
| **384 红归因（三发）** | ①同案单独跑 ×3 @ 本腿二进制 ②HEAD 对照树 | ①**三发全绿**（1.57–1.63s）；②**`git archive HEAD` 纯净对照树（零本腿字节）六族合跑 = 38 RUN／36 PASS／2 FAIL**：同案同红点（3.38s）＋`Test41a2r1WarmGateFeedsKwsNotTheAsr`（本腿两次合跑该案均绿）⇒ **合跑负载时序红是 HEAD 既有形状，非本腿回归；本腿红名集合零新增**（133 的先例：这台机器上按红名集合比，不当四数稳定量） |
| **延迟三数** | 同合跑 `Test15eSpeechEndToEnd…` | `LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` **逐位同**（381-v2/26r3/381-r2 基线值） |
| ②`internal/speech` | staged speech.test.exe 全量 | **67 PASS／0 FAIL／0 SKIP，rc 0**（与 41a2r1 记录的 PASS=61+KWS 面同族口径） |
| ③`internal/session` | `go test -count=1 ./internal/session` | `ok … 4.473s`，rc 0 |
| ④`go vet ./cmd/wisp` | —— | rc 0 |
| 甲 `gofmt -l` 本腿 3 枚 | —— | 0 行 |
| 乙 `gofumpt -l`（`D:\work\base\gopath\bin\gofumpt.exe`）本腿 3 枚 | —— | 0 行 |
| `d22scan` | `git clone --shared` 导出树 `b109a3010`＋替入本腿 3 枚，`go run . -root <export>` | **rc 0 clean**；分母 `bans#1-5 internal/=286 cmd/=51、ban#8 internal/=726 cmd/=177、design 14 frontend 116 tools 25`（较 381-r2 时代的 274/45/706/162 为 HEAD 上他腿的正常演进；clean） |
| d22scan 正控 | 注入 `✓`(U+2713) 字符串／裸 `go func(` | **rc 1** `[emoji] …string literals included`／**rc 1** `[bare-goroutine] …use observe.Registry.Spawn`；还原探针复跑 **rc 0 clean**（分母同）＝可逆 |

## 6 ★ 票面五格逐条（实现腿自评；**翻框归验收**）

| 格 | 票面原文 | 本腿给到的东西 | 自评 |
|---|---|---|---|
| **AC#1** | 球单击→TriggerInterrupt 真挂（对照 Esc 通路挂接形状；internal/ball API 若需扩展 ⇒ 具名交回） | `attachVoiceInterrupt(rv.TriggerInterrupt)`（`startResidentVoice`，与 `attachWakePause` 同块）＋`clickInterruptGesture`（`OnClickBall` 闭包同体）；**ball API 无需扩展**——`ball.Events.OnClickBall` 既有事件面直接够用，`internal/ball` 零字节；线程形状与 Esc 同（同在 ball UI 线程内联回调） | **读数齐，建议可翻** |
| **AC#2** | 端到端：播报中球单击 ⇒ TTS 真停 ≤400ms ＋ 门重开（复用 381-r2 staged 测量形状） | 真机五轮 22–39ms（§2 表）＋`gate_closed=6 gate_reopened=6`＋`spoken` 不计被打断句＋下一句照常（`CLICK-27 aborted=1 next_spoken=1 probe2_aborts=0 gate 2/2`） | **读数齐，建议可翻** |
| **AC#3** | 反向判据：摘球单击挂接 ⇒ 球单击不再打断 ⇒ 红（三读数＋还原核） | §4：M1 突变红（具名红句）＋同突变 Esc 通路照常绿（定向自证）＋原树还原核绿＋SHA 互异＋工作树零突变窗口 | **读数齐，建议可翻** |
| **AC#4** | 零回归（六族＋Esc 通路不动） | 六族合跑两发 41/42，唯一红 384 经 HEAD 纯净对照树同红归因（既有负载时序形状）＋speech 67/67＋session ok＋延迟三数逐位同＋Esc 通路案全绿（§5） | **读数齐，建议可翻** |
| **AC#5** | 纪律格：⛔ 翻框/改票面/台账/push；⛔ `internal/{ball,risk,session,speech}/**`（ball API 扩展须具名交回） | 写面恰 3 枚 cmd/wisp＋本件；`internal/**`、testdata、frontend、design、scripts、.github、票面、台账全零字节；d22scan clean＋正控可逆 | 遵守 |

## 7 ★ 未定义即停三条核销

①「internal/ball 无单击事件且加事件＝改冻结面」——**未触发**：`ball.Events.OnClickBall` 既有（`ball_windows.go:50` 字段、`:632` fire 点），零扩展、冻结面零字节；②「球单击线程与 Esc 的 UI 线程形状不同」——**未触发**：单击回调与 Esc cancel key 同在 ball ui-sta 线程内联跑（`Ball.fire` 契约），381-v2 的「abort 必须从触发方线程出手」由 381-r2 的 TriggerInterrupt 非阻塞形状（buffered send＋直接 abort）满足，单击线程完全同形；③「必须动 381-r2 本体」——**未触发**：`TriggerInterrupt`／`speakAbort`／`speak()` 入口清 token 零字节，本腿只挂第二触发源（`resident_voice_windows.go` 只在 `startResidentVoice` 加一行 attach，381-r2 的件未动）。

## 8 ★ 写面申报（逐枚）

1. `M cmd/wisp/resident_ball_windows.go` —— voiceInterrupt 字段＋attachVoiceInterrupt/currentVoiceInterrupt＋clickInterruptGesture＋OnClickBall 闭包＋ballGestureWhy/bullet/boot 行文案随真
2. `M cmd/wisp/resident_voice_windows.go` —— startResidentVoice 里 `rb.attachVoiceInterrupt(rv.TriggerInterrupt)` 一行＋注释（381-r2 本体零字节）
3. `A cmd/wisp/resident_voice_interrupt_390r1_windows_test.go` —— 4 案（M1 具名红句＋真机五轮＋idle 正控＋unattached 反形）
4. `A docs/evidence/s1/390r1-impl.md` —— 本件（证据笔单独提交）

**残余具名（不阻塞）**：①单击的 summon 半（把单击变任务）仍无执行者——`ballGestureWhy` 已随真，非本票射程；②合跑负载下 `Test384ResidentWatchdogEndToEnd` 的 3s 时序窗敏感（HEAD 对照树同红）——归口 384 侧既有形状，本腿零触碰其代码。

**交回**：编排者核验——票面五格读数齐、AC#2 ≤400ms 真机五轮在册、反向三读数＋还原核齐、零回归红名集合零新增——即可派验收/翻框。
