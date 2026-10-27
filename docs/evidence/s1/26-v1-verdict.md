# 票 26 终裁腿 `26-v1` —— 七格终裁（非实现者）

> **尺**：`HEAD=f1983bb9c`（`dev`，2026-10-27 19:3x–20:1x +08；父＝编排者 `A1287` 收 26-r3 后的收件笔）。★**导出树**＝`git clone --shared` → `D:\tmp\wisp\20261026\26v1\export`（HEAD 同值、`git status` 0 行）。★**落点前 `ls`**：`docs/evidence/s1/26-v1-verdict.md` 不存在 ⇒ 新建；未过 180 行 ⇒ 不拆 `-partN`。
> **工作树现况（开工现跑）**：脏 1202 条 ＝ `.scratch/**` 1184 ＋ `design/**` 11 ＋ 票 115 的 M（他人）＋ 杂 7；**`internal/**`／`cmd/**`／`docs/**` 零脏行** ⇒ 无在飞码污染读数；★`347-v1`／`233-r1` 在飞件零碰（本腿全部突变走仓外导出树＋`-overlay`）。
> **⛔ 未改**：`internal/**`／`cmd/**`／`testdata/**`／`.scratch/wisp/**`／`docs/reports/**`／`docs/specs/**`／`docs/PLAN.md`／`docs/audit/**`；**⛔ push**。收工核：`git status --porcelain -- cmd/wisp`／`-- internal/speech` 各 **0 行**。

## 0 一句话判语
**七格里 4 格可翻（AC2 新入列）、3 格具名残余**；26-r3 的最小可结条件第一件（`SetSpeaking` 一针＋self-echo 夹具）**真在盘上且牙还在**（M1/M2 我自己复造、全红、还原绿）。★**判"能带残余结案"**——但**带一个新抓的生产缺陷**：`buildResidentVoice` 补了 queue 却仍缺 `provScope`/`newSpeaker` 两字段（**R26-7，本腿新具名**），生产播报"出声"半边在 nil 渲染缝上有 panic 面。★**未定义即停：一处**——tri-v1 最小可结条件写的是"AC7/AC1 各**开承接票**"，而工单池现跑**没有任何 R26 承接票**（唯一 `26-tts-output` 提及＝票 30 被 blocked）；"开票"属编排者动作，本腿按 15-done 先例给"票面具名即结"路线，**是否要先开两票由编排者拍**（§10）。

## 1 票面 7 格现状（现跑，两把尺）
7 框 0 勾（工作树与 HEAD blob 两尺同值，`^\s*- \[[ x]\]`）；`Status: ready-for-agent`／`Claimed by: —`／`Last update: 2026-09-19`（未动）；Progress log 0 行。

## 2 七格逐条判（对照 tri-v1，"变了没"）
| 格 | tri-v1 判 | **现判** | 依据（本腿现跑） |
|---|---|---|---|
| AC1 出声＋球＋打断 | ⛔ | **⛔ 没变** | `EvInterrupt` 生产仍只两处**失败分支**（`resident_voice_engines_windows.go:940` 播报失败／`:970` 播报中断）；热键／球单击／否决词零通道；≤400ms 无读数。26-r3 没碰 ⇒ **R26-1 仍在** |
| AC2 不自激＋采集门 | ⛔ | **✅ 变了，能翻（带两注）** | ①针在盘上：`engines_windows.go:912` `SetSpeaking(true)`（HOP 2b，换族后、首样本前）／`:903` `false`（`gateRelease` defer，六出口恰一次：闸拦 `:868`、无引擎 `:877/:921`、无设备 `:942`、合成失败 `:972`、drain 失败 `:976`、成功）；②夹具在盘上且本轮复绿：`GATE-26R3 closed=1 reopened=1 events=2 mid_speech_transcripts=0 spoken_samples=89389 gate_on=1 gate_off=1`／`ECHO-26R3 gate_state_during_feed=speaking(PathT capture refused) frames_fed=6 transcripts_after=1`／`REOPEN-26R3 finals=2`；③★tri-v1 的"ring 重放带播放期回声（机理未测）"**现判结构性消除**：门闭 ⇒ inner 源 `Stop()`（`gate.go:12-16` 逐字 "the inner source is Stop()ed (device off, D16)"；`gate_test.go:47` 钉死 inner 帧计数冻结）⇒ 播放期零新帧可进 ring（`run():728` 无条件追加的是**采到的**帧）⇒ 下个 epoch `Accept(ring)`（`:680`）无回声可重放。注 1：票面 "13 counter asserts zero frames" 的真身＝`internal/audio/gate_test.go:47`（inner `FramesSent`），席位在 13 不在 26，设备级零帧归生产装配态；注 2：夹具喂的是数字静音 ⇒ "零转写"对照力有限，真回声波形的自激面随 Path C 归 59 |
| AC3 Path C barge-in | ⛔ | **⛔ 没变** | 票 59 现跑 `ready-for-agent`、无 `-done`（blocked by 13）；AEC／P15／≤400ms 全归 59 ⇒ **R26-3 仍在** |
| AC4 首音 800/1600ms | 能翻（半格必带注） | **✅ 判不变** | 26-r3 未碰；暖档 26-r1 p50 181.4ms／p95 195.1ms（10 轮）；★冷档 NOT MET 必注：26-r1 **5.82/6.12s**、装配播报 acquire wire-r1 4.29s／tri-v1 4.82s／**本轮 26R3 案实测 4.0550224s**，对 1600ms 全线 2.5–3.8× ⇒ **D32 16.3.3 退化档照写** |
| AC5 串行驻留 peak=max | 能翻 | **✅ 判不变** | tri-v1 比值 0.41<0.75；本轮 `internal/speech` staging 67/0/0 全绿（含 `Test26TheRealSwitchDoesNotAddTheTwoFamilies` 9.45s PASS） |
| AC6 P7 或 blocked-decision | 能翻（第二支） | **✅ 判不变** | `R5` 云 TTS blocked-decision 已入账；P7 人评表现跑**不存在**（`docs/evidence` grep 0 命中） |
| AC7 ChTTS 闸＋26 侧闸测 | ⛔ | **⛔ 变了半格，仍 ⛔** | 调用点半边已收（`speak():864`，wire-r1-M1 牙＋`TestWireR1AnnounceGateRefusesTaintedText` 本轮绿）；但①**毒源**：生产 `.Mark(` 现跑仍只有 `panel_assets.go:245`（`taintSourceScopeID`）一处，`voice-announce` 的 Mark 全仓唯一写者仍是测试 `resident_voice_wire_r1_windows_test.go:569` ⇒ 闸生产**永不命中**；②**命中无 L2 确认**：`speak()` 只 `noteProblem` 即 return（`:867`），全文件无 confirm/card。26-r3 没碰 `internal/risk` ⇒ **R26-5 两条都仍在** |

## 3 反向判据两发（我自己造；仓外 `D:\tmp\wisp\20261026\26v1\export` @`f1983bb9c` ＋ `-overlay`）
| 发 | 变异 | rc | 逐字红句 |
|---|---|---|---|
| **M1** | `export\cmd\wisp\resident_voice_engines_windows.go` 删 `rv.captureGate.SetSpeaking(true)` 一行（991→990 行，恰 −1） | **1** | `MUTANT 26r3-M1 SHAPE: the announce path never closed the capture gate (closed=0 reopened=1): SetSpeaking(true) is not on the path, so the D16 half-duplex close never happens and the microphone would hear the announcement` |
| **M2** | 同法删 `rv.captureGate.SetSpeaking(false)` 一行（991→990） | **1** | `MUTANT 26r3-M2 SHAPE: the capture gate was left in Speaking after the sentence finished (closed=1 reopened=0 speaking=true): SetSpeaking(false) is not on the path; capture stays closed - stuck in Speaking - and no later utterance can ever be captured again`（★含 "stuck in Speaking"） |
| **还原核**（时机＝M2 之后） | 无 overlay 重编重跑 `Test26R3` | **0** | **3/3 PASS**（17.67s/0.00s/0.07s）；导出树 `git status` **0 行**、HEAD 同值；`MUTANT 26r3` 4 处全在测试件；工作树 `cmd/wisp`/`internal/speech` 零脏 |

★口径更正（清单第 3 例同型）：M1 第一发在 `mut\` 目录跑得 `cannot find build\models\manifest.json` rc 1 —— `voice15eRepoFile` 从 **cwd** 向上找仓根，尺量错对象；改 cwd＝仓根重跑才得真红。★两发红句各含该格的可观测后果（M1=麦克风会听到播报；M2=卡死在 Speaking）。

## 4 生产缺陷核（我自己核的）
- **修复在盘上**：`buildResidentVoice`（`resident_voice_windows.go:1299`）字面量第 **1302** 行 `queue: newVoiceOwnerQueue()`，且在全部 early-return 之前 ⇒ 每条分支（含降级返回）都有队列；`announce()` 首行 `queue==nil` 守卫（`engines:784`）不再拒。
- **构造器普查（grep 现跑）**：`&residentVoice{` 全仓恰 2 处（`:419` 测试／`:1300` 生产）＋ 测试零值字面量 1 处（`selfecho:264`，故意量 nil 安全）；`newResidentVoice` 2 个调用者全在测试（wire-r1:229/563、15e:256）。**两条构造路现在都建 queue** ⇒ 26r3 抓的"测试构造器 vs 生产构造器"那一个差异已消。
- ★★★**但没消干净（本腿新抓，具名 R26-7）**：生产构造器字面量**仍缺 `provScope` 与 `newSpeaker`**（两字段只在测试构造器 `:427-428` 赋值；我核过 wire-r1 原版 `d803344ae` 与 26r3 终版 `06609feaf` 的 blob，两版都如此）。后果两半：①`provScope=""` ⇒ `attachProvenance:447` 对空串 `OpenScope`（`provenance.go:460` 会注册空 id），但 `announce():798-801` 有空串回落 `voice-announce` 的守卫 ⇒ **闸的 scope 面被兜住**；②`newSpeaker=nil` ⇒ 生产 `speak()` HOP 3/4（`engines:935`）对 **nil 函数值调用 ⇒ panic**；owner goroutine 在 `observe.Registry` recover 边界内（`goroutine.go:296`）⇒ 进程不塌、但该次播报死、gateRelease 在 panic 展开时仍会发 `false`（Go defer 语义）⇒ 门不会卡，**出声半边死**。⇒ **"测试构造器 vs 生产构造器是两条路"还有第二处**；测试全走 `newResidentVoice` ⇒ staging 全绿而生产出声路径待验，与 26r3 抓的 queue 缺陷同一家法。**AC2 的门针不受它影响**（true 在 HOP 2b 先于渲染、false 在 defer）。
- **P7 评分表**：不存在（现跑）；**R26 承接票**：工单池 0 张（现跑）。

## 5 正控（我自己复跑）
`Test26R3SilenceNeverTouchesTheGateNeedle`：**120 帧静音 ⇒ `SILENCE-26R3 finals=0 announce_asked=0 spoken=0 samples=0 gate_closed=0 gate_reopened=0`**（`frames_drained=121`）——门针两套计数零触发。外加 15e 侧 `SILENCE drained=120 handoffs=0 finals=0`（同轮合跑）。反向正控：同轮真句 `samples=89389` 真出声 ⇒ "静音不响"不是仪器瞎。

## 6 零回归复跑（我自己复跑，staging）
| 尺 | 读数 |
|---|---|
| `TestWireR1\|Test15e\|Test26R3` 合跑 | **16 PASS／0 FAIL／0 SKIP**，rc 0（15e 6＋wire-r1 7＋26R3 3） |
| 延迟三数 | `LATENCY openDelay=262ms chunk=640ms ballLatency=902ms` ⇒ **逐位同值** |
| 关键计数 | `drained=169/181`（**夹具固有抖动族** 178/179/181，⛔ 当判据；`finals=1 handed=1 handoffs=1 asrLoaded=true problems=[]`） |
| 播报针 | `gate_closed=1 gate_reopened=1`（26R3 案）；wire-r1 七案 `gate_closed=0 gate_reopened=0`（未接门缝的台架不动）＋ `TestWireR1AnnounceEndToEnd` 绿（samples≈89.4k） |

## 7 `d22scan` 导出树＋正控（第 21 处定式）
导出树 @`f1983bb9c`（status 0 行）干净跑 **rc 0 `clean`**；分母 `bans#1-5 internal/=264、cmd/=45、ban#8 internal/=684、cmd/=159`（frontend 116／design 14／tools 25）——与 26r3 终版逐位同。正控两发（注入导出树、验后即移）：A 字符串放 `U+2713` ⇒ **rc 1** `[emoji] zz_poscontrol_26v1.go:5`；B 裸 `go func(` ⇒ **rc 1** `[bare-goroutine] :6`；删两枚复跑 **rc 0**、分母回落、status 0 行 ⇒ **正控可逆**。★⛔ 未拿 ban#8 扫 `docs/`（不在其 scope）。

## 8 四数复跑（staging ＋ 定向）
| 面 | rc | 读数 |
|---|---|---|
| `cmd/wisp` 定向（staged 二进制） | **0** | `--- PASS TestAC3LogSinkLandsInsideTheEnvDataRoot (0.00s)` |
| `internal/speech` staged（带 `WISP_PUNC_MODEL`） | **0** | **67 PASS／0 FAIL／0 SKIP**（无 punc 那发 61/0/6 skip——6 枚标点案具名 skip，非红） |
| `internal/session`（零字节） | **0** | `ok … 2.267s` |
| `internal/audio`（零字节） | **0** | `ok … 16.362s` |
| `go vet ./cmd/wisp ./internal/speech` | **0** | 空 |

## 9 残余清单（`R26-1..6` 逐条现状 ＋ 新增 R26-7）
| 号 | 内容 | 现状 |
|---|---|---|
| R26-1 | 生产打断通道（热键/球单击/否决词→#27）＋≤400ms 读数 | **仍在**（AC1 ⛔） |
| R26-2 | 采集门未驱动＋self-echo 夹具＋ring 重放回声 | **收口**（26r3 两件＋本腿复验；两注见 §2-AC2） |
| R26-3 | Path C/AEC/P15 → 59 | **仍在**（59 未 done） |
| R26-4 | 冷档 4.1–6.1s vs 1600ms、TTS 驻留策略缺 | **仍在**（AC4 必带注；本轮 acquire 4.055s 佐证） |
| R26-5 | 任务毒到不了 `voice-announce`＋命中无 L2 确认 | **仍在**（AC7 ⛔；26-r3 没碰 risk） |
| R26-6 | P7 人评表缺 | **仍在，不挡结案**（`R5` 第二支已收口） |
| **R26-7（本腿新具名）** | `buildResidentVoice` 仍缺 `provScope`/`newSpeaker` ⇒ 生产出声半边 nil 渲染缝 panic 面（进程不塌、该次播报死） | **新立**；修法＝生产构造器字面量补两字段（与 queue 同形），或两构造器合一 |

## 10 能不能带残余结案（对照 tri-v1 最小可结条件）
- **tri-v1 原条件**："把 AC2 的 SetSpeaking 一针接上＋一条 self-echo 夹具（同一张票内可做完）、AC7 的 L2 确认与毒源各开承接票、AC1 的打断通道开承接票；**四条具名后可结**"。
- **现状**：第一件**已做完且本腿复验**（针＋夹具＋M1/M2＋正控＋零回归全绿）；AC7/AC1 的承接票**未开**（工单池现跑 0 张 R26 票）。"具名"的载体按 15-done 先例可落在票 26 票面（Progress log 逐条具名残余＋归口），不必先开票；若编排者按字面要"先开两票再结"，差的就只是那两张票的动作。
- **★判语：能带残余结案**（4 格翻／3 格具名残余／R26-6 不挡），**条件＝结案时在票 26 Progress log 追加一段残余清单**：`R26-1 打断通道（归承接票 NNN 或具名待开）`、`R26-3 归 59`、`R26-4 冷档 NOT MET（归换族/驻留策略票）`、`R26-5 毒源＋L2 确认（归承接票 NNN）`、`R26-6 P7 人评（R5 已收口）`、`R26-7 生产构造器补 provScope/newSpeaker（建议与 R26-1 同一张承接票，都在 cmd/wisp 装配面）`。**若结的逐格写法**：AC1 不勾＝"归 R26-1"；AC2 勾＋两注；AC3 不勾＝"归 59"；AC4 勾＋"cold-load NOT MET (5.82/6.12s vs 1600ms；装配 acquire 4.06–4.82s)"；AC5 勾；AC6 勾（第二支 `R5`）；AC7 不勾＝"归 R26-5"。AC2 的勾框语句必须写"capture gate closed during playback **(13-counter 半句以门态读数等价收口，设备级归生产装配态)**"，⛔ 静默吞掉票面原文。

## 11 提交号 ＋ 逐笔越界
| 笔 | commit | 面 |
|---|---|---|
| 本腿 | 见下方提交号 | **只 `docs/evidence/s1/26-v1-verdict.md`**（`add`/`commit` 两处显式单文件 pathspec；⛔ push） |

★越界核：写面仅本件＋`D:\tmp\wisp\20261026\26v1\**`（仓外）；全部变异/正控在导出树；还原核时机在最后一发变异之后；rc 一律 `$LASTEXITCODE` 非`管道`。
