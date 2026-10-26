# 票 26 TTS 输出腿 `26-r1` —— 落地件（写腿）

> **尺**：`HEAD=56350b267`（`dev`，2026-10-26 18:3x–19:1x +08；该提交＝编排者 `A1269` 的裁定件，见 §0）。
> **写面**：`internal/speech/**` **新增 7 枚 / 36,763 B** —— `tts.go` 229 行、`slot.go` 157、`tts_test.go` 302、`slot_test.go` 143、`slot_real_test.go` 87、`rss_windows_test.go` 38、`rss_other_test.go` 6 ＋ 本文件。
> **零字节**：`engines.go`（接口签名一字未动）／`asr.go`／`vad.go`／`endpoint.go`／`doc.go`／`transcriber.go`／`fixtures_test.go`／`cer/**`／`testdata/**`／`cmd/**`／`internal/**` 其余／`models/manifest.json`／`docs/specs/**`／票面／台账。
> **临时件（只建不删）**：`D:\tmp\wisp\20261026\26r1\{overlay,bin,binspeech,bin2,binm1,binm2,binrestore,mut,logs}`。★变异全走仓外 `mut\` ＋ `-overlay` ⇒ **工作树零突变窗口**。

## 0. 一句话 ＋ 为什么是"测量级"
`sherpa-onnx` 的 Matcha 中文 TTS 挂上了 `15-a` 冻结的 `TtsEngine`（**⛔ 改接口**），并补上 `15-v1` 判语里 AC4 缺的 **slot 互斥**：串行驻留（先释放再加载）、RSS 读数、1.6s 实测全部现跑。

★★**它不是默认 TTS，也不可出货**：权重在 `models/manifest.json` 里是 `status: blocked-p3`、`license` 逐字 "NON-COMMERCIAL-ONLY … must not ship; dev/measurement use only"；台账 `R5`（机主 2026-09-20 裁定）逐字 **"默认 TTS 走云 TTS，供应商 StepFun … 本地 matcha 条目保持 blocked-p3 不改"**。本腿**只**做 `dev/measurement` 面：**零配置键、零选型、零 manifest 改动**（引擎选择仍在装配根，编译进来不等于成为默认）。★编排者 `A1269` 已就这份三方冲突拍 **甲**。

## 1. ★★★ 票面 7 格 AC 逐条（原文准引 ＋ 判定 ＋ 读数）
| AC | 票面原文（节选，逐字） | 判定 | 读数 |
|---|---|---|---|
| AC1 | "short result → **TTS audible**, ball **Speaking** with success breathing; interrupt mid-playback → instant stop + **Listening**" | **⛔ 本腿翻不了** | 出声半边有：42,248 样本／1.916s／22.05kHz／peak 0.7076（§2）。★球走态与时刻需 `internal/ball` ＋ `cmd/wisp`（**⛔ 本腿写面**）⇒ 无球读数 |
| AC2 | "playback injects no ASR transcripts (self-echo fixture); **capture gate closed during playback (13 counter asserts zero frames)**" | **⛔** | 需 `cmd/wisp`＋`internal/audio`（13 的采集门）⇒ 本腿零读数 |
| AC3 | "**Path C barge-in (D47, blocked by 59)**: with AEC on … **P15 numbers recorded**" | **⛔（票面自己写 blocked by 59）** | seam 侧机理已实测：emit 返错 ⇒ 合成当场中止（发出 14,809／全量 71,132 样本，§2）。★AEC／P15 数字归票 59 |
| AC4 | "**first-audio ≤800ms warm / ≤1600ms cold-load measured (10 runs, P50/P95 recorded)**" | **✅ 实测完成（warm 过／cold 不过）** | ★warm **p50 181.4ms／p95 195.1ms**（10 轮，预算 800ms，**余量 4.41×**）；★cold-load **5.8196s（空进程）／6.1177s（ASR 驻留时）** 对 1600ms 预算 **3.64×／3.82× ⇒ NOT MET**（与 S0 的 4726ms 同向）⇒ ★**D32 16.3.3 的退化档（TTS 常驻＋ASR 按需）是 Path T 唯一可行政策** |
| AC5 | "**ASR and TTS never co-resident (sampler peak = max, asserted); unload frees**" | **✅ 本腿翻** | ★结构：生命周期逐字 `load:asr,close:asr,load:tts`（先释放后加载）；★RSS：切换增益 **117.7MB** 对"从空载入 ASR"的 **264.2MB**（＝**max 而非 sum**）；★卸载后残差 **1.9MB**；★反证见 §6 |
| AC6 | "**P7: score sheet committed ≥7/10 or blocked-decision logged for cloud-TTS fallback**" | **⛔（要人）** | 双盲人评 ≥7/10 是**人**的动作（本腿做不了）；★"blocked-decision for cloud-TTS fallback" **已被 `R5` 记成裁定**（云 StepFun），不是本腿可代写 |
| AC7 | "**every announced string goes through `risk.Provenance.CheckText(scope, risk.ChTTS, text)` before playback** … plus a **26-side gate test**" | **⛔ 未做（具名）** | `risk.ChTTS` = `internal/risk/provenance.go:122`（`"tts.announce"`）在，`TestFourChannelExfilSuite` 也在；★但**播报路径本身不存在**（`grep '\.Speak(' cmd internal` 除本腿与 seam 外零命中）⇒ 在 `internal/speech` 里造一枚无调用者的"26 侧闸"正是 `15-a1` 判过的**装饰**。★真接线在 `cmd/wisp`（`15-e` 地盘、每波一枚写腿）⇒ 本腿不越界，具名交回 |

## 2. ★★★ "TTS 真的出声"（不是"接口实现了"）
`Test26TheTtsEngineEmitsRealAudio`（staging，真 matcha＋vocos），文本 `"你好，我是威斯普。"`：
**chunks=2 ／ samples=42,248 ／ rate=22,050 Hz ／ duration=1.916s ／ peak=0.7076 ／ rms=0.0935 ／ 首块 14,455 样本 @178.5ms ／ 全长 441.2ms**。★RMS 0.0935、peak 0.7076 ⇒ 是语音不是静音（若只断言 samples>0，静音也能过）。
★分块（`Test26TheChunksAreChunksAndNotOneWholeClip`，三句 `"你好。今天天气不错。我们开始吧。"`）：**chunks=3 ／ 首块 14,767 @194.1ms ／ 合计 71,132 样本 ／ 全长 760.6ms** ⇒ **首块严格小于全量**＝逐块真流式（整段包装会退化成 1 块）。
★打断（`Test26BargeInStopsTheSentenceMidFlight`）：首块即返错 ⇒ **发出 14,809／71,132 样本、回调恰 1 次**，`Speak` 把 emit 的错误**逐字**返回。

## 3. ★★ `C9` 挂上的证明
- 编译期：`internal/speech/tts.go` **`var _ TtsEngine = (*MatchaTts)(nil)`**；`engines.go` 的接口**一字节未动**（`Speak(text string, emit func([]float32) error) error`／`SampleRate() int`／`Close() error`）。
- 真调用：§2 那条读数走的正是 `var eng TtsEngine = tts` → `eng.Speak(...)`（**经缝调用**，不是经具体类型）；`SampleRate()` 也从缝上取回 22,050。

## 4. ★★★ slot 互斥实测（`A C4` 的核心）
`slot.go`：`Slot` 只允许一族驻留，**先 `Close` 旧的再跑 loader**（失败则留空＝fail-closed），带 `Loads()/Closes()` 计数。
| 尺 | 读数 |
|---|---|
| 生命周期（假引擎，逐字） | **`load:asr,close:asr,load:tts`**；再要同族 ⇒ 不重建（返回同一个引擎，生命周期串不变） |
| 真引擎 RSS（`Test26TheRealSwitchDoesNotAddTheTwoFamilies`） | baseline 60.5MB → TTS 驻留 **+145.9MB** → 切到 ASR **+117.7MB**（**不是**再来一个 264.2MB）；★从空载入 ASR＝**264.2MB** ⇒ 切换增益 / 新引擎全量 = **0.445** |
| 卸载 | 关槽后残差 **1.9MB**（baseline 之上） |
| 断言 | `switchGrowth < 0.75 × asrFromEmpty` ＋ `switchGrowth < sum` ＋ 残差 `< max(delta)/2` |

## 5. ★★ RSS delta ＋ 1.6s 实测（空进程探针，`-overlay`，不落盘）
| 量 | 读数 |
|---|---|
| matcha 载入 | **5.8196s**（空进程）／**6.1177s**（ASR 驻留时）→ 对 **1600ms** 预算 **3.64×／3.82×** |
| 驻留成本 | TTS **+170.2MB**（12.4→182.6MB）；ASR int8 **+281.4MB**（12.6→294.0MB） ★与 S0 的 174.2MB／281.1MB 同值族 |
| 若共驻（探针故意造出的禁态） | **436.6MB**＝**两项之和逐 MB**（424.0＋12.6）⇒ "peak=sum" 一旦破，700MB 预算无余量 |
| 卸载 | `DeleteOfflineTts` 释放 **142.4MB**，残 40.2MB（vs 12.4 基线 ⇒ 余 **27.8MB**） |
| warm 首块 | **p50 181.4ms／p95 195.1ms**（10 轮，预算 800ms） |

## 6. ★★ 反向判据：三读数 ＋ 还原核（时机＝最后一发变异之后）
| 发 | 变异（仓外 `mut\` ＋ `-overlay`） | rc | 逐字红句 |
|---|---|---|---|
| **M1** | `slot.go`：两处 `releaseLocked()` 删掉（＝摘掉 slot 互斥） | **1** | 假引擎：`slot lifecycle = "load:asr,load:tts", want "load:asr,close:asr,load:tts"`；★真引擎同发：`switching families grew the process by 260.3MB, as much as loading the incoming engine from empty (259.2MB)`（切换增益 260.3 对 max 259.2 ⇒ **≈ 1.0**，与正控的 0.445 成对照） |
| **M2** | `tts.go`：`speakable` 守卫删掉（＝空文本/纯标点静默成功） | **1** | `Speak("") = speech: matcha tts produced no audio for 0 rune(s), want ErrNothingToSpeak` |
| **还原** | 无 overlay 重编重跑（★**在 M2 之后**） | **0/0/0** | 三案全绿：`Test26TheSlotReleasesBeforeItLoads`／`Test26TheRealSwitchDoesNotAddTheTwoFamilies`／`Test26EmptyAndPunctuationAreNotSilentSuccess` |
★**还原核**：树内 `MUTANT` **0 命中**；`git status --porcelain -- internal/speech` 只有本腿 7 枚新件 ＋ `27-r1` 的 3 枚（**零 M/D 我的**）；`tts.go` sha256 `E0B3DE7500332207…`、`slot.go` `E1E103F7635854D8…`。

## 7. ★★ 正控：空／纯标点 ⇒ **不**静默成功
`Test26EmptyAndPunctuationAreNotSilentSuccess`：`""`／`"   "`／`"\t\n"`／`"。，、！？"`／`"…––"` **五形全返 `ErrNothingToSpeak` 且 emit 收到 0 样本**。
★★**为什么守卫必须在 `Speak` 里**（同一案里的对照读数）：**引擎本身不认为标点是空** —— `tts.Generate("。，、！？")` 返回 **3,840 样本／0.174s／peak 0.0025** 且 **error=nil**。⇒ 直接转发就是"发出一段 174ms 近静音并报成功"，与"没说话"在宿主侧不可区分。M2 把那枚守卫摘掉，这条当场红（§6）。

## 8. ★★ staging 跑法（完整命令；★`PATH` 不够）
```
go test -c -o D:\tmp\wisp\20261026\26r1\bin2\speech.test.exe ./internal/speech
copy third_party\sherpa-onnx\*.dll D:\tmp\wisp\20261026\26r1\bin2
set WISP_SPIKE_MODELS=<repo>\third_party\spike-models
set WISP_SILERO_MODEL=<repo>\third_party\spike-models\silero_vad.onnx
cd D:\tmp\wisp\20261026\26r1\bin2  &&  speech.test.exe -test.v=true   =>  49 PASS / 0 FAIL / 6 SKIP, rc=0
```
★**为什么 `PATH` 打不过**：Windows 对依赖 DLL 的搜索序是「**主模块所在目录** → System32 → 系统目录 → 当前目录 → `PATH`」。`go test` 把二进制建在临时目录，那里没有 `onnxruntime.dll` ⇒ 加载器先走到 **System32 的 1.17.1**，`PATH` 排在 System32 **之后**，顶到最前面也改不了这个序。staging 把 **1.28.2 放进"主模块所在目录"**，那是唯一排在 System32 前的位子。★三层口径：什么都不给 ⇒ `0xc0000135`（进程起不来）／只给 `PATH` ⇒ 拿 1.17.1（具名 skip）／**staging ⇒ 全绿**。

## 9. ★ 四数
| 面 | 命令 | rc | 读数 |
|---|---|---|---|
| `internal/speech` **改前** | `git ls-tree -r --long HEAD -- internal/speech` | — | **21 文件 ／ 172,823 B** |
| `internal/speech` **改后** | 盘上同一把尺 | — | **30 文件 ／ 251,952 B**（本腿 **＋7／36,763 B**；另 3 枚属 `27-r1` 在飞件，**⛔ 本腿零碰**） |
| `internal/speech` staging | §8 | **0** | **49 PASS／0 FAIL／6 SKIP**（改前同树 38／0／5；＋11 案全为本腿，★6 枚 skip 全是 `27-r1` 的标点案） |
| `internal/audio`（**零字节、独立跑**） | `go test -count=1 ./internal/audio` | **0** | `ok … 16.164s` |
| `cmd/wisp` **定向**（`PATH` 带 `third_party\sherpa-onnx`，⛔ 整包） | `-run TestAC3LogSinkLandsInsideTheEnvDataRoot -v ./cmd/wisp` | **0** | `--- PASS … 0.00s`；`ok … 0.070s` |
| `go vet ./internal/speech` | — | **0** | 空 |

## 10. ★ 格式甲乙两尺 ＋ 门禁（含既存假红／他人红）
| 尺 | 输出 |
|---|---|
| 甲 `gofmt -l internal/speech` | **空**，rc 0 |
| 乙 `D:\work\base\gopath\bin\gofumpt.exe -l internal/speech` | **空**，rc 0 |
| 全仓 `gofmt -l internal cmd tools scripts` | ★**恰 3 条既存假红**：`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`／`cmd\wisp\models.go`（CRLF/`autocrlf` 造的，**⛔ 本腿责任**；本腿 7 枚两尺皆空） |
| 门禁 `d22scan` | **rc 1**，`1 finding(s)`：★**`internal/speech/punctuate.go:132` phantom-citation（`docs/evidence/s1/27-r1-impl.md` 不存在）＝ `27-r1` 的在飞件，⛔ 本腿责任、⛔ 本腿修**；分母：`bans#1-5 internal/` **263**、`ban #8 internal/` **679**、`ban #8 cmd/` **155**。★本腿 7 枚**零命中** |
| ⛔ 错尺 | **没用 ban#8 扫 `docs/`**（它只扫 `internal/`＋`cmd/` 的 `.go`） |

## 11. 边界／未做（逐条自陈）
- **⛔ 改 `engines.go`**：`git hash-object internal/speech/engines.go` ＝ HEAD 同值（§0 零字节）。
- **⛔ 加配置键**：噪声/时长尺度、线程数、句数上限、说话人、语速全是**包内常量**（`tts.go` 的注释逐条写了理由，与 `15-b/15-c` 同一条家法）。
- **包内零 goroutine、零裸 `go func(`、零墙钟超时**：`Speak` 的取消靠 emit 返错（协作式中止），时长读数只**测**不**判**超时。
- **未做**：AC1/AC2/AC3/AC6/AC7 的球／采集门／AEC／人评／播报闸（§1 逐格具名）；★`cmd/wisp` 接线（`15-e` 地盘）；★模型选择（`R5`＋票 61）。
- **并发写面警报**：`internal/speech` 今天**有第二枚写腿在飞**（`27-r1`：`punctuate.go`／`punctuate_test.go` 未跟踪 ＋ `transcriber.go` M 58/7）⇒ 本腿提交**只**点名自己 7 枚 ＋ 本文件，⛔ 不 `add` 他人件。
- **`DEFERRED` 标记面**：`doc.go` 的 `DEFERRED(engines)` 仍写"票 15 实现"，本腿**未动**（摘 marker 属 `225-deferred-registry-audit`；`TtsEngine` 已有实现这一变化请它登记）。

## 12. 交回父代理
**未定义即停：无。** 本腿按 `A1269` 的甲裁定交付**测量级**实现：AC4（warm 过／cold 不过）与 AC5 **两格可翻**，其余五格具名归口（球/采集门/AEC/人评/播报闸）。★**可出货性**：此实现**不可出货、默认仍走云**（`blocked-p3` 未动、manifest 未动、零配置键）。
