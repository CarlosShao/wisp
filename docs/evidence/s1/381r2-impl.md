# 票 381 实现腿 `381-r2` —— 落地件证据表（Abort 接线＋入口清 token＋veto 选删）

> **身份**：写腿（实现者）自报读数。**判定归验收/重裁腿**——本件不翻任何一格，票面六格只给「读数＋实现腿自评」。
> **尺**：开工 `HEAD=3d41c1791`（`dev`，索引 0 staged、`cmd/wisp` 无他人脏；共享树工作区改动 `docs/reports/pending-and-issues.md`＋`internal/config/{loader,manager}.go`＝他腿写面，零碰）。**收工前编排者把共享树推进到 `716874ae1`（385-r1 派单存档，cmd/wisp 零触碰）**，本腿全部读数即在该 tip＋本腿 5 文件的状态上现跑；提交 `77c2cabfb`（5 文件，+756/−54），父提交 `716874ae1`。现跑 2026-10-31 11:33–12:3x +08。
> **staging**：`go test -c -o D:\tmp\wisp\20261031\381r2\bin\wisp.test.exe ./cmd/wisp`＋三 DLL 拷到同目录＋`WISP_SILERO_MODEL/WISP_SPIKE_MODELS/WISP_PUNC_MODEL/WISP_TTS_DIR/WISP_15E_TMP`（26r3 §7 同款；`-test.run` 只定向、`-timeout 30m`、GOTMPDIR 已置空）。
> **临时件（只建不删）**：`D:\tmp\wisp\20261031\381r2\{bin,mut,export,cases,logs,fmt,speechbin}`。**⛔ push；⛔ 票面/台账/specs/PLAN/testdata/frontend/design/scripts/.github 零字节；⛔ internal/** 语义零碰；⛔ D43 表零碰；配置键零新增。**

## 0 一句话

381-v2 点名的三个结构缺口全部补上并各带反向判据：**出声期 Esc 从「6,653ms 照播完」变成真停（真机五轮 41–46ms，≤400ms 达标）**；**残留 token 不再无声杀下一句**；**veto 字符串支路删除（具名选删），良性句 12/12 到达文本路、零吞句零污染**。零回归 27/0＋262/640/902 逐位同。

## 1 ★ 三件结构（逐件：接在哪／为什么／读数）

| 件 | 接在哪 | 为什么是这个形状 |
|---|---|---|
| **①接 `speaker.Abort()`** | `resident_voice_windows.go` `TriggerInterrupt`：发 token 之外 `rv.speakAbort.Load().(func())` 调在飞 speaker 的 `Abort()`；`resident_voice_engines_windows.go` `speak()` 在 `newSpeaker` 成功后 `rv.speakAbort.Store(func(){ speaker.Abort() })`、defer 清空；`voice_out_windows.go` `Abort()` 先置 `abortReq`（atomic.Bool，在锁外）再拿锁做 reset+unprepare | 出声期 owner goroutine 停在 `speaker.Drain()` 里，chunk 边界的 select **结构性零跑**（381-v2 已实测）——唯一能 ≤400ms 停出声的线程是**触发打断的那条**（Esc 在 ball UI 线程），所以 abort 必须从 `TriggerInterrupt` 出手。`abortReq` 放锁外是因为 `Drain` 持锁等待：拿锁才能提请求的话，请求要等整块 chunk 放完才被看见——正是 6,653ms 的形状 |
| **①speaker 侧配套** | `releaseLocked` 等待改**分片**（`voiceOutWaitSlice=100ms`，预算仍是「样本数→毫秒＋2000ms 宽限」，**零墙钟**，D42#9 不破）；片间查 `abortReq`→返回 `errSpeakerAborted`；并加 `WHDR_DONE` 旗标核（driver 先置旗后发信号——顺带修掉既有「WOM_OPEN 残留信号提前 unprepare 仍在播缓冲」的潜在形状）；`Drain` 入口与 `Abort` 全程感知旗标 | 分片让「停在播」的 Drain 在 ≤1 片内作答；sentinel 让 `speak()` 区分「被打断」（走 D43#27、不记 spoken、不记 problem——381-r1 对合成半边的同一裁定）与「设备故障」 |
| **②入口清 token** | `speak()` 第一行：`select { case <-rv.interruptCh: default: }` **恰一次** | 381-v2 §5：空闲期 Esc 留 token ⇒ 下一条播报首个 chunk 回调把它当真打断吃掉（spoken=0、无 problem、无声杀句）。入口清一次＝残留 token 消；入口**之后**到达的 token 是本句真打断，留给 chunk 回调/abort，不被清 |
| **③veto 选删** | `consume()` `TranscriptFinal` 分支：`isVeto := strings.Contains(停/取消/别)` 整块删除，final 一律 `finals.Add`→`finish()`→文本路；头注释具名理由 | **具名选型＝选①删**（派单两案之一）：文本路 veto **做不了它唯一的活**——Speaking 期采集门已闭＋owner 在 speak() 里，veto final 结构性到不了；非 Speaking 期它 12/12 误伤、dispatch 被 D43 拒、`continue` 吞句＋noteProblem 污染（381-v2 §4 三结构性事实全中）。真 veto 通道＝KWS 词级（`onKwsHit`，`EvVeto` 走 #13/#22，与字符串匹配不同源、零碰）＋Esc（`approval.voiceInterrupt`→`TriggerInterrupt`），都不走转写文本路 |

配套读数面：`ttsAborted` 计数器＋`announceReading()` 打断=%d＋`stopLine()` aborted=%d；`resident_windows.go` 里 381-r1 的重复 `raudio.approval = ra` 两行并一行（过期行号注释一并修真）。

## 2 ★★★ ≤400ms 实测（真 waveOut 设备＋真 matcha 合成＋生产 TriggerInterrupt 通路；`Test381r2RealDevicePlayingPhaseAbortWithin400ms`）

句子＝30 字生产句形，matcha 发 **chunks=4 / ≈140,818 样本 ≈6.4s 音频**（比 381-v2 的 1–2 chunk 句更长——abort 不依赖 chunk 边界，chunk 数无关紧要）。每轮：播报→等 Drain 停入（parked）→**300ms 后注 Esc**（`TriggerInterrupt`）→计 trigger→停。

| 轮 | unstick（trigger→Drain 作答＋#27） | reset（trigger→`waveOutReset` 落地） | parked→stop | 对照（无打断） |
|---|---|---|---|---|
| 记录轮 1 | **46ms** | 46ms | 347ms | control `parked_to_done=2.095s`（剩余队列照播完） |
| 记录轮 2 | **45ms** | 45ms | 346ms | （同一 control 句形 381-v2 时代＝照播完 **6,653ms**） |
| 记录轮 3 | **41ms** | 41ms | 341ms | |
| 记录轮 4 | **46ms** | 46ms | 346ms | |
| 记录轮 5 | **46ms** | 46ms | 347ms | |

（前一轮同源二进制另取 4 轮：42/45/55/49ms——**九轮合计全部 ≤57ms**，`logs/381r2-live-rounds-2.txt`、`381r2-zero-regression.txt`。）
**≤400ms：MET**（判格归验收）。`abort_returned` 读数在 0 与 n 之间波动＝Drain 作答后 owner 的 `Close()` 与异步 `Abort()` 抢锁、**两形都落 waveOutReset、出声都停**（round-5 实测 returned=2、round-1..4 returned=0，stop 时延无差）——具名，非缺陷；`dropped` 故只记读数不作断言（断言设计错一次已改，见 §8④）。

**出声期外的正控（同一测试内）**：control 轮不打断 ⇒ spoken=1、#26 落 Warm、门 1/1；五轮打断后累计 `spoken=1 aborted=5 gate_closed=6 gate_reopened=6`（`voice leg stopped` 行原文在 `381r2-zero-regression.txt`）。

## 3 ★★ 正控（「无事发生」的那一半）

- **非播报期 Esc ⇒ 无事发生**：`Test381r2IdleInterruptIsInert` PASS——空闲 TriggerInterrupt：机器 Sleeping 不动、timeline 空、problems 空、`ttsAborted/ttsSpoken=0`；token 入缓冲等下一条播报入口清。nil leg／无通道 leg 触发不 panic。
- **静音/空文本 ⇒ 门针零触发、无 Abort**：既有 `TestWireR1SilenceAndEmptyTextNeverSpeak`＋`Test26R3SilenceNeverTouchesTheGateNeedle` 在合跑中复绿（`SILENCE-26R3 gate_closed=0 gate_reopened=0`、`spoken=0 aborted=0`）——本腿零新增代码路径被静音触碰。
- **良性句不再被吞**：`Test381r2VetoSubstringFinalsReachTheAgent` PASS——12 句含「停/取消/别」子串的良性语料逐句 fresh leg：finals=1、handed=1、problems=0、timeline 无 EvInterrupt；**误伤率 0/12**（381-r1 支路＝12/12）。

## 4 ★★★ 反向判据三发（仓外 `D:\tmp\wisp\20261031\381r2\mut\` ＋ `-overlay`；红句具名在测试件内）＋还原核

| 发 | 变异 | rc | 逐字红句 |
|---|---|---|---|
| **M1** | 摘 `TriggerInterrupt` 的 abort 接线（token 发送保留＝381-r1 原形） | **1** | `MUTANT 381r2-M1 SHAPE: TriggerInterrupt did not abort the speaker in flight (spoken=1, aborted=0): during playback the owner is parked in Drain where the chunk-boundary select never runs, so without the abort hop nothing stops the audio and the sentence plays out (381-v2's 6,653ms shape)` |
| **M2** | 摘 `speak()` 入口清 token 块 | **1** | `MUTANT 381r2-M2 SHAPE: the announcement was killed by a stale interrupt token (spoken=0, chunks=0): speak() does not drain interruptCh at the entry, so a token left by an earlier idle Esc silently kills the next sentence (381-v2 §5: spoken=0, no problem row)` |
| **M3** | 复插 381-r1 veto 子串支路（含 `strings` import） | **1** | `MUTANT 381r2-M3 SHAPE: the veto substring branch came back and swallowed the user's sentence "导航去最近的停车场" (case 0: finals=0, handed=0): strings.Contains 停/取消/别 is not a veto channel - it cannot see the Speaking state, cannot stop the TTS, and its continue eats the final behind a D43 rejection (381-v2 §4: 12/12 false positives)` |
| **还原核** | 无 overlay 重编重跑同四案 | **0** | 4/4 PASS；三源文件 SHA256 与三变异件互异（tree vs M1/M2/M3 differs=True×3）；`git status --porcelain -- cmd/wisp` 始终恰本腿 5 枚（overlay 只读，工作树零突变窗口） |

## 5 ★★ 零回归＋四数＋两尺＋d22scan

| 尺 | 命令 | 读数 |
|---|---|---|
| **合跑** | staged `wisp.test.exe "-test.run=TestWireR1\|Test15e\|Test26R3\|Test41a2r1\|Test381r2" -test.v -test.timeout=30m` | **27 PASS／0 FAIL／0 SKIP，rc 0**（15e 6＋WireR1 7＋26R3 3＋41a2 5＋381r2 6） |
| **延迟三数** | 同上 `Test15eSpeechEndToEnd…` | `LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` **逐位同**（381-v2/26r3 基线值） |
| 26R3 门针 | 同上 | `GATE-26R3 closed=1 reopened=1 events=2 mid_speech_transcripts=0 spoken_samples=89165`／`ECHO…capture refused frames_fed=6 transcripts_after=1`／`REOPEN finals=2`／`SILENCE gate 0/0` 全在 |
| ②`internal/speech` | staged speech.test.exe 全量 | **61 PASS／0 FAIL／6 SKIP，rc 0**（与 41a2r1 记录 RUN=67/PASS=61/FAIL=0/SKIP=6 逐位同） |
| ③`internal/session` | `go test -count=1` | `ok … 2.804s`，rc 0 |
| ④`go vet ./cmd/wisp` | —— | rc 0 |
| 甲 `gofmt -l internal cmd tools scripts frontend` | —— | 恰 4 行：`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`＝**已知 CRLF 假红**（blob 拷出 gofmt 干净、worktree CR=131/225，381-v2 同法复证）；`internal\audio\aecspike\rig_windows{,_test}.go`＝**未跟踪在飞件**（`git show HEAD:` 报 not in HEAD，他腿写面）。**`cmd\wisp` 0 行——`models.go` 确已不在假红名单**（与派单 ⚠ 一致）。本腿 5 枚两尺全净 |
| 乙 `gofumpt -l`（`D:\work\base\gopath\bin\gofumpt.exe` v0.12.0） | —— | 逐字同甲 |
| `d22scan` | `git clone --shared` 导出树 `716874ae1`＋替入本腿 5 文件，`go run . -root <export>` | **rc 0 clean**；分母 `bans#1-5 internal/=274 cmd/=45、ban#8 internal/=706 cmd/=162、design 14 frontend 116 tools 25` |
| d22scan 正控 | 注入 `✓`(U+2713) 字符串／裸 `go func(` | **rc 1** `[emoji] …string literals included`／**rc 1** `[bare-goroutine] …use observe.Registry.Spawn`；删探针复跑 **rc 0 clean**（分母同）＝可逆 |

## 6 ★ 票面六格逐条（实现腿自评；**翻框归验收**）

| 格 | 381-v2 判 | 本腿给到的东西 | 自评 |
|---|---|---|---|
| AC#1 构造器补全 | ✅ 可翻（v2 已核） | 本腿零触碰（`buildResidentVoice` 未动） | —— |
| AC#2 nil-panic 消除 | ✅ 可翻（v2 已核） | 本腿零触碰 | —— |
| **AC#3 打断通道＋≤400ms** | ⛔（R26-1 仍在） | Esc 真通道（381-r1 接线）**＋出声期 Abort 真停**：真机五轮 41–46ms（§2，≤400ms 达标）；合成期 token 通路 M1/M2 验真 | **读数齐，建议可翻** |
| **AC#4 端到端** | ⛔ | 真停（九轮 reset 落地＋`dropped` 读数）＋门放开（打断轮 gate 1/1…6/6）＋下一句能进（`Test381r2PlayingPhaseInterruptAbortsAndWalks27`：aborted=1→next_spoken=1、probe2_aborts=0、gate 2/2）＋残留 token 杀句消（§3/§4 M2） | **读数齐，建议可翻** |
| AC#5 反向判据 | ✅ 可翻（v2 补跑） | 本腿为 **r2 新结构**另跑三发（§4）＋还原核——与 381-r1 时代的三发不同源，防回退 | 维持可翻 |
| AC#6 纪律格 | ✅ | 写面恰 5 枚 cmd/wisp＋本件；`internal/**`、D43、配置键、testdata、票面、台账全零字节；d22scan clean＋正控可逆 | 维持 |

**R26-1 收口判定建议**：≤400ms 已达标（§2）＋真通道三条中 Esc 端到端全绿、KWS veto 通路 41a2 已接、球单击仍未装配（381-r1 自认，非本腿三件）——**R26-1 的「打断」半格具备收口读数**，收口与否归编排者。

## 7 ★ 未定义即停四条核销

①`Abort()` 真停＝waveOutReset 语义相符（九轮实测停播）⇒ 未触发；②≤400ms 修到最优＝41–46ms 远低于线 ⇒ 未触发；③veto 两案＝选删可行且已做 ⇒ 未触发；④未动 `internal/speech` 任何语义（只消费公开面：`speech.Transcript/Endpoint/ErrNothingToSpeak` 等既有面）⇒ 未触发。

## 8 ★ 写面申报（逐枚）＋残余具名

**提交 `77c2cabfb`（5 文件，+756/−54）**：
1. `M cmd/wisp/voice_out_windows.go` —— abortReq/分片等待/DONE 旗核/sentinel/Abort 前置置旗/Drain 入口查
2. `M cmd/wisp/resident_voice_windows.go` —— speakAbort/ttsAborted 字段、TriggerInterrupt abort 半边、consume() veto 删支、strings import 移除、读数行
3. `M cmd/wisp/resident_voice_engines_windows.go` —— speak() 入口清 token、abort 登记/清除、drain-abort 分支（#27）
4. `M cmd/wisp/resident_windows.go` —— 381-r1 重复赋值并一行
5. `A cmd/wisp/resident_voice_interrupt_381r2_windows_test.go` —— 6 案（M1/M2/M3 具名红句＋真机五轮＋正控）
**本件** `docs/evidence/s1/381r2-impl.md`（证据笔单独提交）。⛔ push。

**残余具名（不阻塞，归口在册）**：
① `Close()/Abort()` 抢锁两形＝`dropped` 读数波动，出声两形都停（§2 已具名）；② 出声期 KWS veto 结构性不可能在飞（owner 正忙说话）——KWS veto 的射程是 Listening/Confirming（#13/#22），两通道不混；③ 球单击打断仍未装配（381-r1 自认项，不在本腿三件）；④ 测试自身曾有两处竞态设计错（`dropped≥1` 断言、轮次完成标记早于 #27 dispatch）——已改，教训＝「完成标记必须取 speak() 最末可观测（计数＋状态/门）」；⑤ 381-v2 两条新抓的归口现状：speaker Drain 早回——本腿的 DONE 旗标分片等待已实质收窄（control `parked_to_done=2.095s` vs 381-v2 时代 9.84s 信号/15.1s 真实出声的错位形状），账仍归 26 侧核；`internal/memory/undo.go` gofmt 真红——本腿两尺现读已不点名（已被他人修）。
