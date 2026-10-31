# 票 41 A2 接线腿 `41a2-r1`（写腿裁决表）— KWS 接进 `cmd/wisp` 真管线，`warmFeed` 升真喂食门

> **尺**：开工 `HEAD=4cd8cb4b1`（`dev`，索引 0 脏、`cmd/wisp` 无人脏；`281862907` 与 381-r1 的 `ae66669f1` 均在祖先链上——开工时 dev 已由编排者推进到派单存档提交本身，非 detached/旧树）→ 收工前 dev 又前进 4 笔（至 `dcf6e910f`，全为编排者台账/派单存档，**`cmd/wisp` 0 触碰**，d22scan 导出树以 `dcf6e910f` 为基）。现跑 2026-10-31 09:55–11:2x +08。★读数逐条现跑，⛔ 转述。
> **写面（三件，逐名）**：`cmd/wisp/resident_voice_windows.go`（+460/−30）、`cmd/wisp/resident_audio_windows.go`（+9，静音手势的 D43 半步）、`cmd/wisp/resident_voice_41a2r1_windows_test.go`（新，441 行）＋本文件。★**零字节**：`internal/speech/**`（`git status` 无它——`kws.go` 一字未碰、⛔ 触发停手条件）、`internal/risk/**`、`internal/statemachine/**`（D43 表 0 字节）、`internal/audio/**`、`internal/config/**`、`frontend/**`、`design/**`。
> **临时件（只建不删）**：`D:\tmp\wisp\20261028\41a2r1\{bin,mut,cases,logs,export}`（前腿 10-28 的 `bin/`＋`cases/` 续用）。★变异全走仓外件＋`-overlay` ⇒ **工作树零突变窗口**（§5 还原核）。
> **装载形状（按派单照抄）**：商店路径＝manifest 拼法五件直填 `speech.KwsConfig`（`encoder/decoder/joiner-epoch-12-avg-2-chunk-16-left-64.int8.onnx`＋`tokens.txt`＋`keywords.txt`）；开工现跑 sha256：epoch-12 与 epoch-99 两套逐位同、五件全对 manifest 钉值（`dd784973…`/`ed834540…`/`f7976005…`/`72316508…`/`46de8c68…`）⇒ 引擎腿读数可继承。

## 0. 一句话

`internal/speech` 的 `KwsEngine` 第一次有了生产调用者：opt-in＋模型五件齐 ⇒ boot 线程构造（VAD 先例）→ **D43#5 Sleeping→Armed**；`run()` 帧管道在 ring 之后喂 `kws.Accept`（phase1 段＋Warm 分支两处，⛔ 插 seam）；**`warmFeed` 从闩升级成真喂食门**（Warm 里 `tr.Accept` 停、KWS 当耳朵）＝**R28-2 落地**；命中派 `EvWakeWord`（#7）/Warm 里走 #29/`EvVeto`（#13/#22＋381 打断通道）；`Facts.KwsLoaded` 有了生产赋值，#8/#10/#32 以真 fact 走通。

## 1. 票面 6 格逐条（AC 原文逐字＋判定＋本腿读数）

| # | 票面原文（逐字） | 判定 | 读数 |
|---|---|---|---|
| AC1 | `Armed RSS/CPU gate measured (sampler) within D32 row; wake latency ≤300ms.` | **半格**：sampler RSS/CPU 门与 300ms 延迟数需要票 42 的 watchdog/sampler 面，本腿不冒称；**唤醒延迟的原料已落**（timeline 的 AtSample 全样本钟） | `KWS-TIMELINE` 三跳 AtSample=0/115200/163840（§2）；Armed 预算由 D32 行与 `IntraOpThreads` 钉（引擎腿），本腿未测 RSS——**具名残余** |
| AC2 | `Keywords: custom keyword list change → reload without restart; hit on fixture audio.` | **半格**：`hit on fixture audio`＝真管线读数（§2，demo 表配对）；`reload without restart` 的词表热更＝**不落**（`[voice] wake_word.keywords/thresholds` 两键仍是 255 钉的零读；消费它们需要词表→keywords 文件的生成器，归 B 腿/config 票） | 管线级：0.wav 0 hit、4.wav 1 hit（蒋友伯）、生产表 7 夹具全 0 hit（§1b 控制案）；`enabled` 键 boot 读一次（热更等重启，与 `mic_muted_default` 同款，证据内具名） |
| AC3 | `Veto: KWS running + Confirming + veto-word wav → L1 call cancelled; KWS absent → UI message asserted.` | **不覆盖（结构就位）**：veto 派发路径已接（`EvVeto` #13/#22＋381 的 `TriggerInterrupt`），但 `kwsVetoWords` 表为空（出货 keywords.txt 无 veto 词）⇒ **本树上不可能产生真 veto 命中**，e2e 无从断言 | `onKwsHit` 的 veto 分支走 `TriggerInterrupt()`＋Listening/Confirming 派 `EvVeto`（代码在、tests 不伪造）；具名：B 腿改表即激活 |
| AC4 | `P14 false-trigger report (thresholds table + recommended values) committed.` | **不覆盖**：阈值仍吃包内常量 `Score=3.0/Threshold=0.1`（API 默认全 0 hit 的防呆未破）；本腿的负控读数是 P14 的原料之一 | 生产表×7 夹具＝0 hit（管线级，§1b）；静音/非关键词 0 hit 在引擎腿已钉 |
| AC5 | `Mute cycle Armed→Muted→Armed with inference stop/start asserted (no model reload).` | **半格（机制落地、e2e 走了 guard 行）**：静音手势→#8（Armed→Muted，kwsFeed 门在 Muted 关＝停推理）、再按→#10（真 fact 决定 Armed）；引擎对象不关＝无重载（`kwsWake` 只在进程退出由 consumer 的 defer Close） | §4：`#8/#10` 双向走通（真 fact、含 kws-not-loaded 分支）；手势→机器的接线在 `toggleMute` 一点（290 的 Sleeping 形状原样保留） |
| AC6 | `Opt-in off → zero mic opens (device-handle probe) + Sleeping budget.` | **半格**：opt-in off ⇒ 不解析、不构造、不派 #5（Sleeping 原样）；设备句柄探针归采集/票 42 面 | `Test41a2r1OptInGateResolvesTheManifestSpelling`：off ⇒ Kws 空＋`KwsNote` 具名「enabled=false」；⇒ 装载计数 0、无 #5 walk |

**★1b 管线级生产表控制案**（防呆的正面落法）：`Test41a2r1ProductionTableRejectsTheDemoWavs`——生产拼法构造的引擎（keywords.txt）对 demo 7 夹具（614,884 样本）**0 hit** rc=0。「每夹具必中」是模型谎言，⛔ 未犯；唤醒案改吃 **demo 表配对**（引擎腿自己的选择，测试内具名注释说明）。

## 2. ★★★ 唤醒端到端（逐跳读数、At=样本钟；`Test41a2r1WakeEndToEndArmedToListening` rc=0，17.3s）

| 跳 | 量 | 读数（逐字） |
|---|---|---|
| ① 装载 | boot 线程构造 | `wake-word engine loaded and armed model=kws-zipformer-wenetspeech-3.3M-2024-01-01 threshold=0.1`（staging `WISP_KWS_DIR`→third_party 目录；manifest 拼法五件） |
| ② D43#5 | timeline | `Sleeping + kws.enabled -> Armed at_sample=0`；球收 Armed |
| ③ 负控（Armed） | 0.wav 全排（176 帧） | `kws_fed` 前进、`kws_hits=0`、机器仍 Armed、`handoffs=0`（官方负样本在真管线里仍是负样本） |
| ④ 命中 | 4.wav 第 ~49 帧 | `wake-word hit opens the chain load word=蒋友伯 engine_at_sample=115200 event=kws.wake-word`；`kws_hits=1` |
| ⑤ D43#7 | timeline | `Armed + kws.wake-word -> Listening at_sample=115200`；球收 Listening；**Armed 内 VAD 不喂**（#4 无行，纯语音不开链） |
| ⑥ 链装载 | hand-off walk | `handoffs=1 handoff_exit=0`（恰一次——中途 dispose 误判曾致 4 次，已修，§5 还原核含此）；`asr_loaded=true`；标点级接上 |
| ⑦ ASR 开拾 | ring 重放→真转写 | `partials=11 finals=1 handed=1 punctuated_finals=1`；`#11 Listening+asr.vad-stop→Thinking at_sample=163840（语音 2816ms）`；session `SCOPE-WARM reason=turn-end engines_hot=2 loads_after_warm=0` |

**读数面事实**：命中后 KWS 立即停喂（#7 的 kws.pause＝门控语义）⇒ 4.wav 的第二词（女儿）的音频不再进引擎——单命中而非裸流的双命中是**正确形状**；唤醒句本身的转写（ring 重放）会作为话轮交给文本路（唤醒词话轮的过滤策略票面/D43 均未定义 ⇒ 具名残余，⛔ 自创过滤）。

## 3. ★★★ `warmFeed` 真喂食门（R28-2 落地；`Test41a2r1WarmGateFeedsKwsNotTheAsr` rc=0，16.5s）

| 量 | 读数 |
|---|---|
| Turn1 | 唤醒→链 live→尾静音 90 帧→final→`#11`→Warm；`warm_feed=true`（MicClose 钩）、`stage=Warm`、`loads_after_warm=0` |
| **两计数（Warm 态）** | 推 40 帧静音：**`kws_fed` 恰 +20480（=40×512，逐帧精确）**、**`asr_fed` +0（冻结）** ⇒ Warm 里 ASR `tr.Accept` 停、KWS 继续吃帧 |
| 静音正控 | 40 帧静音 0 额外命中；门未开 |
| Warm 唤醒 | 同夹具再喂→`kws_hits=2`→`warm_wake=1`→`warm_feed=false`（门重开）→推 8 帧 `asr_fed` 前进 ✓ |
| 零装载 | `loads_after_warm=0`、`acquire_log=[2.9419532s 0s]`（冷→Warm 免装载，B4） |

**R28-2 判**：可落、已落——`resident_voice_windows.go` 的 ASR 段分支 `if rv.warmFeed.Load() && rv.kwsFeed(samples) { continue }`（无 KWS 时退回 wire-r1 闩形状：transcript 路照旧，行为逐字节兼容）。台账 R28-2 的注销归编排者/验收。

## 4. ★★ `Facts.KwsLoaded` 生产赋值＋guard 行 #8/#10/#32 走通（`Test41a2r1GuardRowsWalkWithTheProductionFact` rc=0）

- **生产赋值**：`rv.kwsLoaded()`（引擎对象在否）是唯一生产者；`rv.kwsFacts()` 在每一发走 KwsLoaded 守卫行的派发里真实传参（此前全仓生产赋值=0）。
- **#10 双分支**（真 fact）：`Muted+EvMuteKey`→`Armed`（kws-loaded）／→`Sleeping`（kws-not-loaded）双向走通；**#8** `Armed+EvMuteKey→Muted` 经真实手势入口 `onMuteGesture`（owner 队列投递，D38b）。
- **#32 双分支**：`Settling+EvSettleExpired`→`Armed`（kws-loaded）／→`Sleeping`（kws-not-loaded）。
- **诚实注记**：#10/#32 的生产派发点现状＝静音手势（#8/#10，本腿接）；#32 的事件本体由 `internal/statemachine` 的内部定时器以 **nil facts** 派发（frozen 包、本腿无钩、15e 残余同款）⇒ 该路径落 Sleeping 后由 **onKwsHit 的 #5 重武装**自愈（Sleeping＋引擎在→#5→Armed→#7）。真相源是会话自己的级联（residency stage），球面残余=wire-r1 §10-④ 同族。

## 5. ★★★ 反向判据两发（仓外 `-overlay`，工作树零突变）＋还原核（时机＝最后一发变异之后）

| 发 | 变异（mut\resident_voice_windows.mut-*.go） | 面红 | rc | 红句（逐字） |
|---|---|---|---|---|
| **M-a** | `kwsFeed` 首行 `if true { return false }`（KWS 永不被喂） | `Test41a2r1WakeEndToEndArmedToListening` | **1** | `the engine was not fed in Armed (kws_fed=0 before, 0 after): the Armed gate is closed` |
| **M-b** | ASR 段门 `if false && rv.warmFeed.Load() && rv.kwsFeed(samples)`（Warm 退回闩、tr 照吃） | `Test41a2r1WarmGateFeedsKwsNotTheAsr` | **1** | `kws_fed moved 0, want exactly 20480: the KWS is not the Warm ears`（asr_fed 冻结断言紧随其后同样会红） |
| **还原** | 正树重编重跑两案 | —— | **0** | 两案 PASS（17.12s／15.03s） |

**还原核（三件）**：①树件 `MUTANT M-a|M-b` 计数 **0**；`git status -- cmd/wisp` 恰 3 件（2 M＋1 ??）；②`git diff --no-index --stat`：mut-a 对树件 **3 insertions**、mut-b **1 ins/1 del**（只差变异块本身）；③sha256 台账：树 `45B6C021…340C14`、mut-a `AEEF9EC0…E64ACB`、mut-b `80164029…552427`。⇒ 工作树从未进突变窗口。

**过程中抓到并修掉的两个真缺陷**（变异前自测暴露，非变异产物）：①0.wav 带 LIST chunk，首版帧数读取错位导致夹具被截断（chunk walk 重写）；②dispose 重臂守卫把「首 epoch 未开」误判为「已释放」（stage() 零值碰撞）⇒ 一发唤醒开 4 次链——守卫收紧为 `tr != nil && stage()==StageDisposed` 后 `handoffs=1`。

## 6. ★★ 零回归（staged，`-test.run=Test15e|TestWireR1|Test26R3|Test41a2r1 -test.timeout=30m`）

| 尺 | 读数 |
|---|---|
| 合跑 | **21 PASS／0 FAIL**，rc=0（15e 6＋wire-r1 7＋26R3 3＋本腿 5）；`logs/full-suite-2.log` 全量 |
| 延迟四数 | `LATENCY openDelay=262ms chunk=640ms ballLatency=902ms (openDelay+chunk) seamWaitToFirstText=5.51s finalLatency=4.096s speechMS=3392 closedBy=end` ⇒ **262/640/902 逐位同值** |
| 15e 关键计数 | `drained=181/181 partials=5 finals=1 handed=1 handoffs=1 asrLoaded=true problems=[]` ⇒ 与 28-r1/wire-r1/381-r1 记录逐位同；**本跑零抖动** |
| 26R3 三案 | 3 PASS；`gate_closed=1 gate_reopened=1`（needle 语义未动） |
| 381-r1 打断 | 覆盖 ride 在 WireR1/26R3 套内（381r1-impl §测试计数同款），7 案全绿 |
| **首跑抖动注记** | 首发合跑 `TestWireR1WarmIdleDisposesOnTheOwner` 红、单跑与复跑绿（5.68s/4.68s）⇒ 按抖动报（该案依赖真实 3s 级联定时器，合跑机器负载下超时），⛔ 不改变通过性 |
| `internal/speech` 整包（零字节） | staging 全套 **rc=0，RUN=67／PASS=61／FAIL=0／SKIP=6**——与 41a2 引擎腿四数逐位同；12 案 `*Kws*` 全绿 |
| `go vet ./cmd/wisp` | rc=0 空 |

## 7. ★ 四数／两尺／d22scan＋正控

| 面 | 命令 | 读数 |
|---|---|---|
| `cmd/wisp` 定向（带 PATH，⛔ 整包 MARK299） | `go test -count=1 -run TestAC3LogSinkLandsInsideTheEnvDataRoot -v ./cmd/wisp` | rc=0 PASS |
| `cmd/wisp` staged 定向全套 | §6 | 21/21 rc=0 |
| 甲 `gofmt -l internal cmd tools` | 收工树 | **恰 2 行**＝`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`（既存 CRLF 假红；派单所记第 3 条 `cmd\wisp\models.go` 已不在列——非本腿责任、非本腿所改）；两尺对本腿 3 件单点＝**空** |
| 乙 `D:\work\base\gopath\bin\gofumpt.exe -l` | 同树 | 逐字同上 2 行 |
| `sh scripts/d22scan.sh`（`git clone --shared` 导出树，HEAD `dcf6e910f` 原样 status=0） | —— | **rc=0 clean**；`bans#1-5 cmd/=45`、`ban#8 cmd/=160` |
| 同上＋本腿三件（唯一差异） | —— | **rc=0 clean**；`ban#8 cmd/ 160→161`（增量即本腿新测试件） |
| **正控** | 导出树植入 bare `go func(` | 仪器红句逐字：`zz_positive_control_41a2r1.go:5: [bare-goroutine] bare \`go func(\` is banned (D22/D38b)` ⇒ **仪器有牙**；删除后复跑 rc=0 clean |

## 8. 边界与停手核对（逐枚）

- **⛔ `internal/speech` 零字节**（`kws.go` 无需改 ⇒ 停手条件不触发）；`internal/risk`/`internal/statemachine`/D43 表零字节；`cmd/wisp` 每波一枚（本腿在飞时无他腿脏面）。
- **派单四条未定义即停**：①D43 合法行全用现成（#5/#7/#13/#22/#29/#33/#8/#10/#32）⇒ 未触发；②warmFeed 升门未改 `internal/speech` 语义（门在 `cmd/wisp` 的帧分支）⇒ 未触发；③Armed 性帧来源=C8 帧管道（生产形状可取）⇒ 未触发；④`KwsEngine` 构造在 cmd/wisp 通（boot 线程、staging 实测 ~2-3s/次）⇒ 未触发。
- **防呆三条**：demo 表/生产表配对未颠倒（§1b 控制案钉死）；`Score=3.0/Threshold=0.1` 未动（`kwsConfig()` 显式留空 Threshold 吃包内常量）；命中必 Reset 在引擎缝内、管线未绕（`kwsFeed` 只调 `Accept`）。
- **代码形状**：零裸 `go func(`（唯一新 goroutine 走 `rv.scope.Go(voiceHandOffGoroutine, …)` 既名册入口）、零墙钟超时（`At`/`AtSample` 全样本钟；时长读数只在测试计时）、`.go` 注释全 ASCII、无 emoji（d22scan rc=0 即仪器读数）。
- **配置键 0 新增**：`enabled` 是既有键的首次生产读（boot 一次，热更等重启，`KwsNote` 具名）；`WISP_KWS_DIR` 是环境变量兜底（与 `WISP_SILERO_MODEL` 同一先例）。

## 9. 未做／具名残余（逐条自陈，⛔ 纸糊）

1. **AC1 的 sampler 半格**：Armed RSS/CPU 门、唤醒延迟 ≤300ms 的实测数——需票 42 watchdog/sampler 面；timeline 的样本钟已备好原料。
2. **AC2 的 reload 半格**：`keywords/thresholds` 两键仍零读；词表→keywords 文件生成器＋热更归 B 腿/config 票。
3. **AC3 的 veto e2e**：veto 派发路径已接但 `kwsVetoWords` 空表使真命中在本树不可能；「KWS absent → UI 明示」的 Confirming UI 断言归 Confirming 通道票。
4. **唤醒句话轮的转写会交给文本路**（含唤醒词本身）：票面/D43 未定义过滤策略 ⇒ 未自创；B 腿/P14 若要滤，接 `consume()` 一处即可。
5. **#32 的内部定时器 nil facts**：机器自己的 3s 定时器派发无钩（frozen 包），KwsLoaded 真值走不到该路径 ⇒ 由 #5 重武装自愈＋测试双分支钉守卫语义；若要真值化需 statemachine 开 facts 钩（改冻结面，另票）。
6. **KWS 引擎字节不走 hand-off walk**（存在性检查，与 VAD 同款 split）：manifest sha256 的核是 models.go 的；`wisp models ensure kws-…` 的交还入口未接 KWS（36MB、opt-in，boot 直载）。
7. **发声期（Speaking）KWS 靠 warmFeed/状态双门停喂**；捕获门（26r3）仍是对播报回声的第一道，KWS 不构成第二道（门关闭期帧不到 KWS）。

## 10. 交回父代理（一句话）

**KWS 已接进真管线**：#5/#7 端到端逐跳可复跑、Warm 真喂食门两计数钉死（R28-2 的解落地）、`Facts.KwsLoaded` 有生产赋值且 #8/#10/#32 双分支走通、两发反向判据双红＋还原核三件、零回归 21/21＋speech 67/61/0/6＋四数逐位同＋d22scan 三发（净/带我/正控有牙）；残余七条具名（sampler、reload、veto e2e、唤醒句过滤、#32 facts 钩、models ensure、Speaking 双门），均不阻塞收票。
