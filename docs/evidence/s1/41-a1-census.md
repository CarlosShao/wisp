# 票 41 KWS 唤醒词普查腿 `41-a1`（只读腿裁决表）— 第一步怎么落

> **尺**：`HEAD=70f020a92`（`dev`，2026-10-27 17:12 +08「编排者A1282：收 347-r1」）· 现跑 2026-10-27 17:16–17:3x +08。★**读数逐条现跑**，⛔ 转述。
> **开工闸**：`git diff --cached --name-status` ⇒ **空**（rc 0）；工作树 `internal/**`／`cmd/**`／`go.mod` **0 行未提交改**（status 全 1202 行都在 `.scratch/`）。
> **写面**：本文件 ＋ `D:\tmp\wisp\20261026\41a1\**`（仓外）。**⛔ 落 `.go`、⛔ go test、⛔ 下载、⛔ push**。
> **先读件**：票面 `41-kws-wake-word.md`（49 行）＋ `tri-v1-verdict.md`（134）＋ `15a-r1-impl.md`（128）＋ `15-a1-census.md`（124）＋ `28-r1-impl.md`（123）＋ `wire-r1-impl.md`（133）——六件全通读。
> **停手条件核对**：①模型**在** ⇒ 不触发；②sherpa-onnx-go KWS 绑定**在**（§③）⇒ 不触发；③票面与 `tri-v1` 无冲突（tri-v1 说"41 未实现"与票面 `Status: ready-for-agent` 同向，R28-2 正是本票要接的债）⇒ **未定义即停未触发**。

## ① 票面 6 格逐条（AC 原文逐字）＋ 哪格是第一步

| # | 票面原文（逐字） | 判 |
|---|---|---|
| AC1 | `Armed RSS/CPU gate measured (sampler) within D32 row; wake latency ≤300ms.` | 需要引擎＋sampler ⇒ 靠后 |
| AC2 | `Keywords: custom keyword list change → reload without restart; hit on fixture audio.` | ★**第一格**：`hit on fixture audio` 不碰装配就能出读数 |
| AC3 | `Veto: KWS running + Confirming + veto-word wav → L1 call cancelled; KWS absent → UI message asserted.` | 依赖 A2 接线＋Confirming 通道 ⇒ 靠后 |
| AC4 | `P14 false-trigger report (thresholds table + recommended values) committed.` | 依赖 A2（TTS/环境音频在跑）⇒ 靠后 |
| AC5 | `Mute cycle Armed→Muted→Armed with inference stop/start asserted (no model reload).` | 依赖 A2 ＋ `Facts.KwsLoaded` 生产赋值 ⇒ 靠后 |
| AC6 | `Opt-in off → zero mic opens (device-handle probe) + Sleeping budget.` | 依赖 A2 ＋ 设备开关面 ⇒ 靠后 |

★**判：AC2 是第一格，它挂在"引擎本体落地"上** ⇒ 第一步＝§⑥ 的 41-a2 引擎腿。**前置现状（现跑）**：`15-…-done.md`（3 未勾，`tri-v1` §11 的先例件）、`28-…-done.md`（`DONE-带残余 2026-10-26`，1 未勾＝AC3 归 59）⇒ **票面两个 blocked-by 都已收口**。

## ② 模型在不在 ⇒ ★**在，且逐位对上 manifest；Apache-2.0，无 P3 问题**

- `ls third_party/spike-models/ | grep -iE 'kws|keyword|zipformer'` ⇒ **`sherpa-onnx-kws-zipformer-wenetspeech-3.3M-2024-01-01`**（另有 matcha／paraformer／silero／vocos）。
- **双份在盘**：`spike-models\` 展开份（**17 文件 ／ 36,213,421 B**）＋ `model-fetch\x-kws\` 同名展开份 ＋ `model-fetch\kws.tar.bz2`（32,654,866 B）。内容：encoder/decoder/joiner 各 **epoch-12 与 epoch-99 两套（fp32＋int8）**、`keywords.txt` 286 B（8 词：你好军哥／蛋哥蛋哥／小爱同学／你好问问／小艺小艺／小米小米／林美丽／你好西西）、`keywords_raw.txt`、`tokens.txt` 1,627 B、`test_wavs\0-6.wav ＋ test_keywords.txt`（8 个测试词）、`README.md`。
- **sha256 五件现跑 vs `models/manifest.json` archive.files**：encoder `dd784973…`／decoder `ed834540…`／joiner `f7976005…`／tokens `72316508…`／keywords `46de8c68…` — **逐位同值**；`kws.tar.bz2`＝`b2f7c896…735f`＝manifest 顶层 `sha256` 同值；`fetch.log:1/111` ＋ `dl.log:13` `OK kws`。⇒ **347-r1/27-r1 的"fetch 包完好但仓内坏件"先例不适用：文件齐、哈希对**。
- **manifest 判**：`id kws-zipformer-wenetspeech-3.3M-2024-01-01`、`purpose kws`、`status "ok"`、`license Apache-2.0`、`quant int8`。
- **许可核（14 的 matcha P3 先例）**：模型卡 README front-matter 逐字 **`license: Apache License 2.0`**（ModelScope pkufool/…，trained on WenetSpeech L）；manifest `license_review.sources` 同判。⇒ **⛔ 无非商用问题**（同一 manifest 里 `tts-matcha-zh-baker` 才是 `blocked-p3`，对照鲜明）。

## ③ `C9` 的 `WakeWordEngine` 接口形状 ＋ sherpa-onnx-go 的 KWS 绑定 ⇒ ★**都在**

- **`internal/speech/engines.go:128-134` 逐字**：`type WakeWordEngine interface { Accept(samples []float32, out []KeywordHit) []KeywordHit; Close() error }`；`KeywordHit{Word string; Veto bool; At int}`（:106-117，`At`＝样本钟 ⛔ 墙钟，D42#9）。名册钉：`engines_test.go:35` `[Accept Close]`、stub `:116`。注释自证形状理由（:123-127）：**KWS 由同一缝喂、512 样本 16k、同 owner goroutine（D38b）**。
- **sherpa-onnx-go v1.13.8 KWS 绑定＝存在**（`go.mod:9` direct）。cgo 实装 `…\sherpa-onnx-go-windows@v1.13.8\sherpa_onnx.go:2322-2478`：`KeywordSpotterConfig{FeatConfig, ModelConfig(OnlineModelConfig.Transducer), MaxActivePaths, KeywordsFile, KeywordsScore, KeywordsThreshold, KeywordsBuf/BufSize}`；`NewKeywordSpotter`（**失败返 nil**，不返 error）、`NewKeywordStream[WithKeywords]`、`IsReady/Decode/Reset/GetResult`（`KeywordSpotterResult{Keyword string}`）、`DeleteKeywordSpotter`；:2466 逐字 **"You MUST call it right after detecting a keyword"**（命中后必 Reset）。平台包装 `sherpa-onnx-go@v1.13.8\sherpa_onnx\sherpa_onnx_windows.go:90-141`（type alias＋var 重导出，与 VAD/ASR 同一形状）。
- **形状映射（我判，实现腿可复核）**：`Accept` 一次喂 ⇒ 流上 `AcceptWaveform`＋`for IsReady { Decode }`＋命中时 `GetResult().Keyword` 填 `KeywordHit{Word, Veto:词在 veto 表, At:Fed 样本钟}`＋`Reset`。
- ⇒ **停手②不触发；①不触发 ⇒ 第一步可以落码**（不是"只报不落"）。

## ④ 两处死结的形状

### ④-1 "关设备会让 warm 唤醒不可达" ⇒ 谁管设备、唤醒走哪条路

- **设备开关的家＝`internal/audio/gate.go`（HalfDuplexGate）＋采集腿**：`SetMuted` :165（mute 手势驱动，**Muted 真关采集双路**）；`WithStartMuted` :56（吃 `[audio] mic_muted_default`，默认 true ⇒ **今天的"麦 OFF"靠静音默认，不靠 KWS-Armed 语义**）；`SetSpeaking` :130 **生产调用 0**（`tri-v1` Q2）；采集腿＝`cmd/wisp/resident_audio_windows.go`（boot 决定建不建，:60-75 注释自认"KWS is still nobody's"）。
- **`wire-r1` 的现状**：`onMicClose`（`resident_voice_windows.go:482`）只 `warmFeed.Store(true)`——**设备没关、门没关**（`tri-v1` :109 逐字"Warm 期间麦克风是开的、音频在进 ASR、设备没关"）。
- **唤醒要走哪条路**：Warm 期间帧仍进 ASR（`run()`:683 `tr.Accept`）⇒ **ASR 是唯一唤醒源** ⇒ 关设备＝唤醒不可达。**KWS 落地后死结的解法形状：KWS 挂在"设备"与"warmFeed"之间**——Warm 门保持开、KWS 吃帧当唤醒源，`warmFeed` 才能升级成**真喂食门**（Warm 里 `tr.Accept` 停、KWS 继续）⇒ D43#26"Warm＝mic CLOSED"在 41 前无法实现（`tri-v1` 判语），**41 就是它的前置**。

### ④-2 `warmFeed` 唤醒闩：`run()` 仍每帧 `tr.Accept` ⇒ KWS 插在哪层

- **四处现跑同值**：声明 `:308`／置真 `:482`／判读 `:830`／置假 `:854`；`consume()` 里 warmFeed 只决定"这句转写要不要先走唤醒"（:830-832 `wakeFromWarm` ⇒ D43#29/#33＋`resume`）。
- **`run()` 两段循环**：phase1 段（:636-655，handOff 前，VAD 判停）＋ASR 段（:660-685，:683 每帧 `tr.Accept`）＋:634 ring 重放。
- ★**插入层判：帧管道（`run()` 的循环体），⛔ 不是 seam 本身**。理由：`WakeWordEngine.Accept` 是**引擎缝**（15a-r1 已冻结名册），而**插入点**是"谁调 Accept"——`phase1` 循环里 samples 进 ring 后喂 `kws.Accept`（Sleeping/Armed 相），hit ⇒ 派 `EvWakeWord`；engine 对象的家＝`voiceResidency`／`speech.Slot`（**`slot.go:31-36` 只有 `FamilyAsr/FamilyTts` ⇒ 需增 `FamilyKws` 枚举＋`Kws()` 访问器，不在共享编译面**）。票面自己的注释已预留（`resident_voice_windows.go:42-43` 逐字："when ticket 41 lands KWS this dispatch becomes EvWakeWord and the table stays byte-identical"）。

## ⑤ 消费侧：唤醒词命中后谁接

- **D43 表（`PLAN.md:3071-3110` ＋ `internal/statemachine/table.go` 现跑，全在、冻结）**：**#5** `Sleeping+EvKwsEnabled→Armed`（table.go:55，"加载 KWS 模型"）｜**#7** `Armed+EvWakeWord→Listening`（:61，KWS 暂停）｜**#8** `Armed+EvMuteKey→Muted`（:65，停推理留模型）｜**#9** Armed 超限自环（:69，不得卸载 KWS）｜**#10** Muted+静音键→Armed/Sleeping（:75/:79，guard＝`Facts.KwsLoaded`）｜**#13** `Listening+EvVeto→Warm`（:93，丢缓冲不落盘）｜**#22** `Confirming+EvVeto→Acting`（:154；PLAN:3092"KWS 未加载时 UI 必须明示语音取消不可用"）｜**#32** Settling 回落（:236-238，KwsLoaded 选 Armed/Sleeping）。
- **事件全存在**：`events.go:18` `EvKwsEnabled`、`:20` `EvWakeWord`（`"kws.wake-word"`）、`Facts.KwsLoaded` :70。
- **生产现状**：`EvKwsEnabled` 生产派发 **0**（唯一 dispatch 在球层测试 `hotkey_live_test.go:408`）；`Facts.KwsLoaded` 生产赋值 **0**（全仓只有 `table_test.go`／`hotkey_live_test.go`）；`resident_voice_windows.go:697/:863` 今天派 `EvSummon`，按 :42-43 注释在 41 落地时换 `EvWakeWord`。
- **票面怎么说**：wake hit → Listening（KWS paused）；veto＝B1 语音取消通道 ~100ms 级、与唤醒词共用同一模型；KWS-not-loaded honesty＝Confirming UI「语音取消不可用」e2e 断言；D47：KWS 只属 Path T，Conversation 结构性不用。

## ⑥ 第一步建议（★最小可派的第一腿）

- ★**腿名 41-a2（A1）：KWS 引擎落地＋fixture 读数**，**⛔ 不碰 `cmd/wisp`**。
- **写面**：`internal/speech/kws.go` ＋ `kws_test.go`（两枚新件）＋ `slot.go` 增 `FamilyKws` 枚举＋`Kws()` nil 安全访问器（与 `Asr()`/`Tts()` 同形；`slot.go` 不在共享编译面）。
- **能出什么读数**：①staging（三 DLL 同目录＋`WISP_SPIKE_MODELS`，同 15a-r1 §8 三层口径）`NewKeywordSpotter` 用 manifest 五件（encoder/decoder/joiner int8＋tokens＋keywords）构造成功；②`test_wavs/0-6.wav`（模型自带夹具）经 `Accept` 得 `KeywordHit{Word=keywords.txt 对应词, At=样本钟, Veto=false}`；③非关键词/静音 wav ⇒ 0 hit（正控）；④`Close` 幂等（C11）；⑤名册钉不动（`engines_test` 仍 `[Accept Close]`）＋d22scan clean＋gofmt/gofumpt 单点空。wav 读取器若包内没有现成的，本腿自带**测试内** reader（⛔ 生产面）。
- **工期：1–2 天**（cgo 配置迭代一次＋夹具齐盘；阈值/词表先吃包内常量，**⛔ 本腿加配置键**——`[voice] wake_word` 消费归接线腿）。
- **后续腿预告（⛔ 本腿不做）**：A2 接线腿（`cmd/wisp`：phase1 循环喂 Accept＋`EvWakeWord` 派发＋Armed 生命周期＋mute 循环＋`Facts.KwsLoaded`）——**必须在 ④-1 的裁定（Warm 门保持开＋warmFeed 升级真喂食门）写清后派**，且排 `cmd/wisp` 共享编译面队列；B 腿（veto words＋L1 veto＋P14）吃 A2 通道。票面并行槽 A 把"KWS engine＋Armed lifecycle"绑一条 ⇒ **判：拆开**，A1 先落（Armed lifecycle 撞 ④-1 的设备开关决策面，先落引擎不亏读数）。

## ⑦ 写面清单 ＋ 冲突

- **本票全量会碰**：`internal/speech`（A1＋B veto）、`cmd/wisp`（A2 接线）、`internal/config`（thresholds/veto/keywords 键消费——**票 255 审计已钉两键零读**：`config_readers_255.go:208-209` 逐字"not consumed by the KWS path yet"）、`internal/models`（KWS 模型 resolve 进 plan，`handOffModel` 同法，`models.go:25`）。
- **⛔ 会碰但本票零改**：`internal/statemachine`（表/事件/Facts 全在、冻结）；`internal/ball`（Armed/Muted 视觉已在 `statevisual.go:158`）；`internal/observe`（Armed 阈值常量已在 `thresholds.go:20/27/37`：**90MB／2.0／7 含 kws-infer**）。
- **现占用（现跑）**：`internal/**`／`cmd/**`／`go.mod` **0 行未提交** ⇒ **今日零冲突**。队列注意：A1282 提交语点名"待派 config 腿修"⇒ 41 的 config 消费腿与它排队；`cmd/wisp` 每波一枚（共享编译面）；`internal/risk` 在飞面已收（347-r1＝A1282 本笔）⇒ ⛔ 扰。

## ⑧ 交回父代理（一句话）

**四全**：模型在（哈希逐位对 manifest、Apache-2.0）、绑定在（v1.13.8 `KeywordSpotter` 全套）、缝在（`WakeWordEngine`/`KeywordHit`/`Slot` 的 `FamilyKws` 空位）、消费表在（#5/#7/#8/#9/#10/#13/#22/#32 全冻结现成）——**第一步＝41-a2 引擎腿**（`internal/speech` 两枚新件＋1–2 天＋fixture 读数）；**两处死结都有形状**：设备开关在 `HalfDuplexGate`（Warm 门保持开、`warmFeed` 升级真喂食门是 R28-2 的解），KWS 插 `run()` 帧管道（phase1 循环）⛔ 不插 seam。
