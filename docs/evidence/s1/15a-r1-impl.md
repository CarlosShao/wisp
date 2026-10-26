# 票 15a 引擎基座腿 `15a-r1` —— 落地件（写腿）

> **尺**：`HEAD=654b6b50d`（`dev`，2026-10-26 17:33 +08，收 `15e-r1` 的发），现跑 2026-10-26 17:34–17:5x +08。★**七项读数逐条在 §1–§7**，全为现跑，⛔ 转述。
> **本片代码 commit ＝ `2c0be38dd`**（★由 `2c0be38dd` **之后的补正 commit 回填** —— 规则禁 `--amend`，要更正就追加）。
> **写面**：`internal/speech/**` —— **新增 2 枚**（`engines.go` 227 行／10,019 B、`engines_test.go` 293 行／10,396 B）＋ **改 1 枚**（`transcriber.go`，`git diff --numstat` ＝ `8 6`）＋ 本文件。
> **零字节**：`asr.go`／`vad.go`／`endpoint.go`／`doc.go`／`fixtures_test.go`／`asr_test.go`／`transcriber_test.go`／`vad_test.go`／`endpoint_test.go`／`cgo_boundary_windows_test.go`（10 枚逐个 `git hash-object` **＝ HEAD 同值**，§6）＋ `internal/speech/cer/**`、`cmd/**`、`testdata/**`、`internal/audio/**`、`internal/models/**`、`scripts/**`、`.github/**`、`frontend/**`、`design/**`、票面、台账、`docs/specs/**`、`docs/PLAN.md`。
> **临时件（只建不删）**：`D:\tmp\wisp\20261026\15ar1\{bin,binspeech,cases,logs,mut,overlay}`。★**五发牙全走仓外件 ＋ `-overlay`** ⇒ **工作树零突变窗口**（§3 还原核）。

## 0. 一句话

`15-a1` 的判语是「`C9` Go 接口在 Go 里**零定义**、`PLAN.md:1365` 只有契约行」。本片把那一行变成 **4 个 Go 接口 ＋ 1 条编译期断言**：`AsrEngine`（★**由已落地的 `15-c` 会话反抽**，⛔ 重设计）／`TtsEngine`／`WakeWordEngine`／`RealtimeEngine`；并把 `Transcriber` 的引擎参数从 `*StreamingAsr` 放宽到 `AsrEngine` —— **唯一一处动已落地代码**，8 增 6 删、零语义变化（§2 的 28→31 ＋ §6 的 blob 双证）。

## 1. ★★★ 读数① `C9` 接口逐条（方法签名 ＋ 对应 `PLAN.md:1365` 哪句）

`PLAN.md:1365` 逐字：`| C9 | AsrEngine TtsEngine WakeWordEngine | 本地 sherpa 与云端 provider 可互换；D47 增 RealtimeEngine（见 C32）：云端 S2S realtime，opt-in、零工具权限 | D5, D47 |`

| 接口（`internal/speech/engines.go`） | 方法签名逐条 | 对应 `:1365` 的哪句 |
|---|---|---|
| `AsrEngine`（:54） | `Begin()` · `Feed(samples []float32)` · `Partial() string` · `Finish() string` · `InFlight() bool` · `Fed() int` · `Close() error` | 「本地 sherpa 与云端 provider 可互换」的 **ASR 半边** |
| `TtsEngine`（:94） | `Speak(text string, emit func(samples []float32) error) error` · `SampleRate() int` · `Close() error` | 同上，**TTS 半边** |
| `WakeWordEngine`（:128） | `Accept(samples []float32, out []KeywordHit) []KeywordHit` · `Close() error`（`KeywordHit{Word string; Veto bool; At int}`） | 同上，**唤醒半边**（本地只有它，`Key constraints` 写死） |
| `RealtimeEngine`（:216） | `Send(samples []float32, out []RealtimeEvent) []RealtimeEvent` · `SampleRate() int` · `Suspend() error` · `Resume() error` · `Close() error` | 「D47 增 `RealtimeEngine`…云端 S2S realtime，opt-in、**零工具权限**」 |

★★ **签名逐条的理由（⛔ 静默）**：
- `AsrEngine` **无 ctx、无 error 返回**：会话由单一 owner goroutine 驱动（D38b）⇒ 取消＝调用者的 ctx、关停＝调用者的等待；解码失败退化成**短/空 final**，而那正是 AC#5「一次重试 → 没听清」唯一键的信号（`asr.go:158-162` 逐字如此）。★后果：**云端 provider 报错与本地模型听错，在缝上是同一件事** ⇒ 调用者永不按 provider 分支。
- `Fed() int` 是**音频钟**（样本数），⛔ 墙钟（D42#9）⇒ `Transcript.Latency` 在任何机器上是同一个数（§2 乙那一发以此断言）。
- `TtsEngine.Speak` 的 **emit-返回错误即中止** 抄的是 C5 `LlmProvider.Stream` 的同一份协作式中止契约 ⇒ Path C 的 barge-in 不用 goroutine、不用定时器。

## 2. ★★★ 读数② 已有实现挂上去的证明（零回归）

**甲：编译期。** `engines.go:78` `var _ AsrEngine = (*StreamingAsr)(nil)` —— `15-c` 的本地 sherpa 会话**不改一行**就挂上了缝。

**乙：驱动期。** `Transcriber` 字段与 `NewTranscriber` 参数 `*StreamingAsr` → `AsrEngine`。这是「互换」从句子变成性质的那一步：`TestTheTranscriberRunsOverAnyAsrEngine` 用**与 sherpa 零共享代码**的 `cloudStubAsr` 走**真判停器**（staging），读数逐字：
`19 partial(s), 1 final(s), final "far-side transcript" closed by end over 28608 samples (1.792s)`；且 `last.Latency == samplesDuration(eng.Fed())` 逐位成立。

**丙：读数不退化。**

| 尺 | 改前 | 改后 |
|---|---|---|
| `internal/speech` staging 全套 | **28 PASS / 0 FAIL / 0 SKIP**，rc 0 | **31 PASS / 0 FAIL / 0 SKIP**，rc 0 |
| 那 28 条**逐名** | —— | 名字集合差集 **missing=0**；新增恰 3 条（`TestC9EngineRostersArePinned`／`TestC9ReservedSlotsAreImplementable`／`TestTheTranscriberRunsOverAnyAsrEngine`） |
| `cmd/wisp` 15e staging（复跑，改后重编） | **6 案 6 过** | **6 案 6 过**，rc 0 |

★**15e 的确定性读数逐位同值**：`openDelay=262ms chunk=640ms ballLatency=902ms`、`partials=5 finals=1 handed=1 handoffs=1`、静音正控 `handoffs=0 finals=0 ball=[] state=Sleeping`、`drained=120/120`。
★**`drained` 不是确定性的**：同一发端到端复跑三遍得 166／171／181（分母 181），`seamWaitToFirstText` 4.838／4.998／5.510s —— 那是 192ms 有界通道在回放解码窗口里丢的帧（15e §1 已记同一现象）。⛔ 当回归，也⛔ 当常量。

## 3. ★★ 读数③ 接口的牙（五发，四发必须编不过 ＋ 一发必须跑红）

变异件在 `D:\tmp\wisp\20261026\15ar1\mut\`，经 `-overlay` 换入编译。**还原核**：`git status --porcelain -- internal/speech` 只有我 3 枚；树内 `MUTANT`／`FedSamples`／`Armed() bool` **0 命中**。

| 发 | 变异 | rc | 逐字（节选） |
|---|---|---|---|
| **A** | `asr.go`：`Fed` → `FedSamples`（真实现少一个方法） | **1** | `engines.go:78:19: cannot use (*StreamingAsr)(nil) … as AsrEngine value …: *StreamingAsr does not implement AsrEngine (missing method Fed)` |
| **B** | 加一枚**整段合成** TTS（sherpa 离线 `Generate` 的形状） | **1** | `wholeClipTts does not implement TtsEngine (wrong type for method Speak)` / `have Speak(string) ([]float32, error)` / `want Speak(string, func([]float32) error) error` |
| **C** | 加一枚**只会报词**的 KWS（无否决位、无音频钟） | **1** | `wordsOnlyKws does not implement WakeWordEngine (wrong type for method Accept)` / `have Accept([]float32) ([]string, bool)` / `want Accept([]float32, []KeywordHit) []KeywordHit` |
| **D** | 拿 `RealtimeEngine` 调工具 | **1** | `e.CallTool undefined (type RealtimeEngine has no field or method CallTool)` ⇒ **零工具权限是结构性的** |
| **E** | 给 `WakeWordEngine` **加**一个 `Armed() bool`（把 stub 一并补上，让编译保持绿） | **0 编过 / 1 跑红** | `engines_test.go:53: WakeWordEngine roster = [Accept Armed Close], want exactly [Accept Close]` |

★★ **E 是这套牙里唯一的"加法"钉子**：Go 的结构化类型对**加方法**是瞎的（不破坏任何实现、不破坏任何调用点），只有**名册钉**（`c9Rosters`，`reflect` 逐方法）会响。⇒ C32 的「零工具权限」不是靠评审，而是**加一个方法就得改一行数据并被拒**。

## 4. ★★ 读数④ 零回归的双向（逐条，改前／改后）

| 尺 | 改前 | 改后 |
|---|---|---|
| 甲 `gofmt -l internal/speech` | 空，rc 0 | **空**，rc 0 |
| 乙 `gofumpt -l internal/speech`（`D:\work\base\gopath\bin\gofumpt.exe`） | 空，rc 0 | **空**，rc 0 |
| `go vet ./internal/speech` | rc 0 | **rc 0** |
| `d22scan`（`cd tools/d22scan && go run . -root <repo>`） | rc 0 `clean`；bans#1-5 internal/=**258**、ban#8 internal/=**664** | **rc 0 `clean`**；bans#1-5 internal/=**259**、ban#8 internal/=**666**（＋2 枚新 `.go`） |

★ 改后 `clean` 是**现跑**，且分母**涨了 2**（新文件被真走到，⛔ 空跑 green）。★`d22scan` 改后**第一次复跑即 rc 0**（本片没有它抓到的违规，⛔ 掩过不说）；另据实核：两枚新文件**全 ASCII**（非 ASCII 字节 `0`／`0`）、语句级 `go` **0 条**、注释里一律写 `->` 而不是箭头。

## 5. ★★ 读数⑤ `C9` 契约行有多种读法 —— 我选哪种、为什么

| 读法 | 内容 | 判 |
|---|---|---|
| 甲 | `C9` 只冻结 3 个名字，`RealtimeEngine` 归 C32，Go 里先不定义 | ⛔ 弃 |
| **乙** | **4 个名字都在 `internal/speech` 定义；`RealtimeEngine` 的语义仍以 C32 为准（`C9` 原文自己写「见 C32」）** | ★**选这个** |
| 丙 | 只定义"今天有实现"的（即只有 `AsrEngine`） | ⛔ 弃 |

**为什么乙**：① `C9` 的**关键内容**列写的就是「本地 sherpa 与云端 provider 可互换」—— 只有一半能换的缝不叫互换；② SPEC-04 §3 给 `C8` 写 `[AecSource] 仅接口位（DEFERRED）`，已立了「**声明但⛔ 实现**」这个家法，后三个接口是同一形状；③ 甲会让这 4 个名字里有一个**继续零定义** —— 正是 `15-a1` 判语的病，换张票就复发。
★**代价写明**：乙把 `RealtimeEngine` 的**方法集**在票 `60` 之前定死了一版。它现在**零实现、零 importer** ⇒ 改它只动一个没人 import 的类型（比票 60 落地后再改便宜一个量级），且 §3-E 的名册钉让"改进来"必须显式过门。

★★ **`AsrEngine` 的方法集是"反抽"不是"设计"**：`15-c` 的会话带实测（首 partial 640ms、flush 1ms），让它去迁就一个假想云端 provider ＝ 拿真数换猜数 ⇒ **不满足这个形状的引擎才是错的**。
★★ **两条我明确没放进接口**：① **opt-in ／「音频将上传」告知 ／ 首次 L2 隐私确认**（C32 逐字要求）—— 那是**宿主策略与配置门**，不是引擎形状；放进接口等于让引擎自己能开麦。② **VAD**：`C9` 的 4 个名字里**没有** `VadEngine`，silero 不可换（SPEC-04 §4 写死）⇒ `Detector` 保持具体类型，⛔ 硬塞一个接口。

## 6. ★ 读数⑥ 四数

| 面 | 命令 | rc | 读数 |
|---|---|---|---|
| `internal/speech` **改前** | `git ls-tree -r --long HEAD -- internal/speech` | —— | **19 文件 ／ 152,258 B**（`transcriber.go` 7,463 B） |
| `internal/speech` **改后** | 盘上同一把尺 | —— | **21 文件 ／ 172,823 B**（＋2 文件 ＋20,415 B；`transcriber.go` 7,613 B ⇒ **＋150 B**） |
| **未改件 10 枚** | `git hash-object`（盘上） vs `git ls-tree HEAD` | —— | **逐枚同值**：`asr.go`／`asr_test.go`／`vad.go`／`vad_test.go`／`endpoint.go`／`endpoint_test.go`／`doc.go`／`fixtures_test.go`／`transcriber_test.go`／`cgo_boundary_windows_test.go` |
| `internal/speech` staging | §8 | **0** | **31 案 31 过 ／ 0 skip** |
| `internal/audio`（**零字节、独立跑**） | `go test -count=1 ./internal/audio` | **0** | 改前 `ok … 15.892s` ／ 改后 `ok … 15.929s` |
| `cmd/wisp` **定向**（`PATH` 带 `third_party\sherpa-onnx`） | `go test -count=1 -run TestAC3LogSinkLandsInsideTheEnvDataRoot -v ./cmd/wisp` | **0** | 改前 `--- PASS … 0.065s` ／ 改后 `--- PASS … 0.065s` |
| `cmd/wisp` **15e staging** | §8 第二段 | **0** | **6 案 6 过** |

## 7. ★ 读数⑦ 格式甲乙两尺（含已知假红）＋ 门禁 rc

| 尺 | 输出 |
|---|---|
| 甲 `gofmt -l internal/speech` | **空** |
| 乙 `gofumpt -l internal/speech`（绝对路径） | **空** |
| **归属判据** `[\\/]` | 两尺对本片 2 枚新文件**单点**跑：甲空、乙空 |
| 全仓 `gofmt -l .` | ★**既存假红原样**：`cmd\wisp\models.go`、`internal\agent\approval\pending_read.go`、`internal\agent\tools.go`（CRLF／`autocrlf` 造）＋ `.scratch/wisp/probes/**` 里一堆**故意坏**的突变件（parse error ⇒ 全仓 `gofmt` rc=2 由它们造成）。★★**`internal/speech` 一枚都不在名单里** —— 派单写的"＋`internal/speech` 的相关件"在本树上**没有实测对象**，据实报空，⛔ 顺手认领别人的红 |
| `d22scan` | **rc 0** `clean - no D22 ban violations` |
| ⛔ 错尺 | **没用 ban#8 扫 `docs/`**（它只扫 `internal/`＋`cmd/` 的 `.go`） |

## 8. staging 跑法（完整命令；★`PATH` 不够）

```
go test -c -o D:\tmp\wisp\20261026\15ar1\binspeech\speech.test.exe ./internal/speech
copy third_party\sherpa-onnx\*.dll D:\tmp\wisp\20261026\15ar1\binspeech
set WISP_SILERO_MODEL=<repo>\third_party\spike-models\silero_vad.onnx
set WISP_SPIKE_MODELS=<repo>\third_party\spike-models
cd D:\tmp\wisp\20261026\15ar1\binspeech   &&   speech.test.exe -test.v        =>  31/31, rc=0
# 15e：go test -c -o ...\bin\wisp.test.exe ./cmd/wisp ＋ 同三枚 DLL ＋ WISP_15E_TMP
#   cd <repo> && ...\bin\wisp.test.exe -test.run=Test15e -test.v                =>  6/6, rc=0
```
★**三层口径（本片复测）**：什么都不给 ⇒ `0xc0000135`（进程起不来）；只给 `PATH` ⇒ 拿 System32 的 1.17.1（`checkOnnxruntime` 具名拒绝／skip）；**staging ⇒ 全绿**。

## 9. 未做／边界（⛔ 越界即作废，逐条自陈）

- ★**「⛔ 改既有语义」这条我是这样守的**：唯一动过的已落地行是 `transcriber.go` 的**两处类型 ＋ 两处注释**，零行为改动；证据＝§2 的 28→31、15e 6/6，＋§6 的 10 枚 blob 同值。★若判"连参数类型也不许动"⇒ 那本片的接口就是**零消费者**的装饰品（正是 `15-a1` 治的病），请编排者裁。
- **未实现**：`TtsEngine`（票 26）／`WakeWordEngine`（票 41）／`RealtimeEngine`（票 60）。本片只冻结缝；`AsrEngine` 的云端实现是票 61。
- **⛔ 加配置键**：本片**零配置面**（引擎选择仍读 `[voice]` 已有字段，方法集里没有一个旋钮）。
- **本片不认领任何 AC 的勾**：`15-a1` 拆片表里 `15-a` 是「AC6 ＋ AC4 机制」，而 AC6 的实测读数与 cgo 边界／单实例已在 `15-b`／`15-c` 落地；AC4 的 slot 互斥实测要 TTS（票 26）。⇒ 本片**没有新增 AC 读数**，⛔ 蹭。
- **零裸 `go`**（语句级 0 条）、**零墙钟超时**（本片没有时长逻辑）、**零新 cgo 面**（`engines.go` 不 import 任何 sherpa 包，纯 Go）。
