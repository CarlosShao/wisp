# 票 27 标点恢复腿 `27-r1` —— 落地件（写腿）

> **尺**：`HEAD=ba9ca8fa7`（`dev`），现跑 2026-10-26 18:30–19:0x +08。★**八项读数逐条在 §1–§9**，全为现跑，⛔ 转述。
> **写面**：`internal/speech/**` —— **新增 2 枚**（`punctuate.go` 12,807 B／`punctuate_test.go` 27,035 B）＋ **改 1 枚**（`transcriber.go` 7,613→10,137 B，+2,524 B）＋ 本文件。
> **零字节**：`engines.go`／`engines_test.go`／`asr.go`／`vad.go`／`endpoint.go`／`doc.go`／`fixtures_test.go`／`asr_test.go`／`vad_test.go`／`endpoint_test.go`／`transcriber_test.go`／`cgo_boundary_windows_test.go`、`internal/speech/cer/**`、`testdata/**`、`cmd/**`、`internal/**` 其余、`scripts/**`、`.github/**`、`frontend/**`、`design/**`、票面、台账、`docs/specs/**`、`docs/PLAN.md`、`docs/audit/**`。
> **临时件（只建不删）**：`D:\tmp\wisp\20261026\27r1\{bin,mut,punc,logs}`。★**四发变异全走仓外件 ＋ `-overlay`** ⇒ **工作树零突变窗口**（§4）。

## 0 一句话

`internal/speech` 长出**标点恢复级**：`punctuate.go`（CT-Punc 类 CT-Transformer int8，`OfflinePunctuation`，一实例常驻）＋ 把 `Transcriber.final()` 接成**唯一的标点注入点**（partials 结构性不标点），＋ 一条**先于 C API 的 onnx 头检查**——后者是被一次真实进程崩溃逼出来的（§9-①）。

## 1 ★★★ 票面 5 格 AC 逐条（原文 ＋ 判定 ＋ 读数）

| AC | 票面原文（逐字） | 判定 | 读数／边界 |
|---|---|---|---|
| AC1 | `Punctuated E2E: wav → transcript with punctuation reaching agent (mock LLM asserts punctuation present).` | **可翻（带边界）** | §3 的 `TestFinalTranscriptCrossesTheSeamPunctuated`：真判停器 ＋ 非 sherpa 引擎 ＋ 真标点模型 ⇒ final **`今天下午三点提醒我开会。`**、`Punctuated=true`。★**边界**：进 agent 的那一跳是 `voiceTranscriptSink`→`residentTaskSource`，在 `cmd/wisp`（15-e 已测），**本片写面外、零碰** ⇒ 断言只到"到达该 sink 的 Transcript"。 |
| AC2 | `Latency ≤300ms P50/P95 on fixtures, recorded in SLO appendix.` | **可翻（带边界）** | **P50=9.255ms／P95=19.587ms／max=21.9ms**，80 次调用（8 例 × 10 轮，暖机排除），预算 300ms ⇒ **15.3× 余量**。★**边界**：`docs/SLO.md` 在仓里**不存在**（AGENTS.md §4 已记），且 `docs/specs/**` 是本腿禁写面 ⇒ 附录落点记在**本件 §1**，不伪称已进 SLO。 |
| AC3 | `300-char dictation fixture subjective score ≥8/10 (evidence sheet) or blocked-decision logged with next-step.` | **可翻** | §2-⑨：`302 字 ⇒ 26 个标点`，逐字对照后我判 **8/10**（理由在 §2-⑨）。 |
| AC4 | `Toggle off path works; load-failure degrade path works.` | **可翻** | 关掉：`TestPunctuationToggleOffPassesTheTranscriptThroughUnchanged`（原串 ＋ `Punctuated=false`）。加载失败：`TestPunctuationLoadFailureDegradesInsteadOfFailingTheUtterance`，**4 种坏件各给具名错误**、句子照样出来（§5-②）。 |
| AC5 | `Memory accounting: punct model included in Warm residency numbers (sampler).` | **⛔ 翻** | **缺：采样器**。🔴 **`sampler` 在 `cmd/wisp`／`internal/models` 一侧，两处都在本片禁写面** ⇒ 本片只能给"一实例、`Close` 释放"的**结构事实**（§3），给不出 Warm 驻留读数。具名残余，见 §9-④。 |

★**票面自身**：`Status: ready-for-agent`、`Claimed by: —`、`Progress log` 0 行、5 条 `- [ ]` 全空。★**⛔ 本腿改票面**（翻框／`Status` 是编排者收尾动作）。

## 2 ★★★ 逐例：输入一串无标点中文 ⇒ 输出加了哪些标点

夹具：`punctuationSamples`（8 例，正是票面场景③的形状：定时、日程、问句、多分句口令）。全部**输入无标点**，故下列标点**每一个都是模型插的**。读数原样来自 `TestPunctuationIsRestoredPerExample` 的 `t.Logf`：

| # | 输入 | 输出 | **加了** |
|---|---|---|---|
| ① | 今天下午三点提醒我开会 | 今天下午三点提醒我开会**。** | `。` |
| ② | 帮我把明天上午的日程整理一下发给我 | 帮我把明天上午的日程整理一下**，**发给我**。** | `，。` |
| ③ | 这个月的电费是多少钱 | 这个月的电费是多少钱**？** | `？` |
| ④ | 先保存文件然后关闭窗口最后重启电脑 | 先保存文件**，**然后关闭窗口**，**最后重启电脑**。** | `，，。` |
| ⑤ | 明天天气怎么样需要带伞吗 | 明天天气怎么样**？**需要带伞吗**？** | `？？` |
| ⑥ | 打开收音机然后把音量调小一点再帮我记一笔买牛奶 | 打开收音机**，**然后把音量调小一点**，**再帮我记一笔买牛奶**。** | `，，。` |
| ⑦ | 周四下午两点半和三点各有一个会 | 周四下午两点半和三点各有一个会**。** | `。` |
| ⑧ | 我刚才说的那些内容你都记住了 | 我刚才说的那些内容**，**你都记住了**。** | `，。` |
| ⑨ | **302 字口头长稿**（§1-AC3 夹具，全文见日志 `mut/../full4.log`） | 26 处标点，见下 | `，`×若干 `。` `？` |

★★**⛔ 只测"函数没报错"**：三条硬性质逐例断言 —— ①真的插了标点（`len(added)>0`）；②**词一个没动**（把两侧标点全剥掉后逐字相等，`contentRunes`）；③**不是全插同一个记号**（全语料**恰好 3 种**`。，？`，一个"句末一律加句号"的假实现过不了）。首字符是标点也算红（无前导标点）。

★**⑨ 的 8/10 怎么来的**（`TestPunctuationOnThe300CharDictationFixture`，302 字一次调用、**未截断**）：26 处边界里 **25 处我判对**——问句落 `？`（"能不能处理完**？**"）、时间词内不切（"九点，我要"／"十点半跟"）、分句落 `，`、长稿结尾落 `。`。**扣 2 分的一处**：`我还没准备你提醒我周四之前弄完` —— "准备" 与 "你提醒我" 之间漏一个 `，`。⇒ **8/10**，过票面线。

## 3 ★★ 与流式 ASR 的接法（接口 ＋ 为什么在这里）

```go
type Punctuator interface {          // punctuate.go，不进 engines.go（见下）
	AddPunct(text string) string     // 无 error：拒绝与"无话可标"在缝上是同一件事
	Close() error                    // C11，幂等
}
func (p *CTTransformerPunctuator) AddPunct(text string) string
func OpenPunctuation(cfg PunctuationConfig, enabled bool) (Punctuator, error)
func (t *Transcriber) SetPunctuator(p Punctuator)      // nil = 关掉 = 加载失败
```

**`ASR 的 final` 怎么进标点**：`Transcriber.final()` 是**唯一**造 final 的地方 ⇒ 标点只在那里发生：

`asr.Finish()` → `t.punctuate(text)` → `Transcript{Kind: TranscriptFinal, Text: punctuated, Punctuated: done}`。

三条与 `15-b`/`15-c` 同款的 seam 风格：①**单一 owner goroutine**、零 goroutine（D38b）；②**不测时间**、无墙钟超时（`PunctuationLatencyBudget` 是**量尺不是死线**，D42#9）；③**无 error 返回**，理由抄 `AsrEngine.Finish`（§ `asr.go:158-162`）。

★**为什么 `Punctuator` ⛔ 写进 `engines.go`**：`engines.go` 是**冻结的 C9 名册**（4 个名字，`c9Rosters` 钉）。标点不是 `C9` 的 5 个名字之一，往冻结件里塞第 5 个名字＝为形式动冻结面 ⇒ **另开文件、`engines.go` 与 `engines_test.go` 一字节未动**（§7 的 blob 核）。另加一枚 `TestPunctuatorRosterIsPinned`（`reflect` 钉 2 个方法）做同一颗牙。

★**partials 结构性不标点**：`partial()` 一行未改 ⇒ 不是靠纪律、是靠**没那条路径**。测试反向核：任一 partial 报 `Punctuated=true` 或带标点即红。

## 4 ★★ 反向判据：三读数 ＋ 还原核（**时机在最后一发变异之后**）

变异件在 `D:\tmp\wisp\20261026\27r1\mut\`，经 `-overlay` 换入，**工作树零改动**：

| 发 | 变异 | rc | 红的断言（逐字） |
|---|---|---|---|
| **A** | **把标点步骤摘掉**：`transcriber.go` 的 `return out, true` → `return text, false`（阶段照跑、结果丢弃） | **1**（43/1/0） | `TestFinalTranscriptCrossesTheSeamPunctuated`：★**`the final reports Punctuated=false with a stage attached: downstream cannot tell dictation-grade text from a raw hypothesis (final "今天下午三点提醒我开会")`** |
| **B** | **静默返回原串**：`punctuate.go` 的 `return p.punc.AddPunct(text)` → `return text` | **1**（41/3/0） | **3 红**：`case timer: "今天下午三点提醒我开会" came back with no punctuation at all -- the stage ran and did nothing`／300 字夹具／seam 那条 |
| **C** | **线程钉松**：`punctuationThreads` `IntraOpThreads`→`4` | **1**（43/1/0） | `Model.NumThreads = 4, want the literal 1: punctuation inference runs inside the single-core voice chain (D32), and the pin is only real while nothing can raise it` |
| **D** | **拆掉"⛔ 吃掉整句"守卫**：`punctuate()` 只留 `return t.punct.AddPunct(text), true` | **1**（43/1/0） | `TestPunctuationNeverErasesAnUtterance`：`a stage that returned nothing ate the utterance: final "", want "今天下午三点提醒我开会"` |

★★**三读数**：**变异红**（A/B/C/D 各 rc=1，红的正是点名那条）／**还原绿**（无 overlay 重编重跑 **rc=0、55 PASS／0 FAIL／0 SKIP**）／**原样绿**（改前 31 案 → 改后 **44 PASS**，见 §7）。

★★**还原核（时机：四发变异全跑完之后，18:52）**：①`git hash-object` 与**变异前**逐个同值 —— `transcriber.go`=`1029197d3656cbe1…`、`punctuate.go`=`1245f4101d9e71b6…`（变异前也打了一遍，逐字相同）；②`go test -c` **无 `-overlay`** 重编 rc=0、重跑 **rc=0 55/55**；③树内**零变异标记**。

## 5 ★ 正控（⛔ "静默返回原串当成功"）

| 夹具 | 读数 | 判 |
|---|---|---|
| 空串 | `""` → `""` | ✅ 不崩、不造标点 |
| 纯空格 `"   "` | `"   "` → `"   "` | ✅ 逐字不变（空白在缝上直接短路，**不进模型**） |
| 纯标点 `"。。，？"` | → `"。。，？"` | ✅ 逐字不变 |
| 纯英文 `set a timer for three pm tomorrow` | → `set a timer for three pm tomorrow**。**` | ⚠ **如实报**：模型给英文补了一枚**中文句号**。词一个没动（`contentRunes` 相等）、不崩；★这是模型 `punc_list` **只有中文标点**（`config.yaml`：`<unk> _ ， 。 ？ 、`）的必然结果，**⛔ 缺陷但也不是"无幻觉"** ⇒ 英文口述要么换英文标点模型、要么加语言门（§9-③ 具名残余）。 |

★**"不静默成功"的牙**：这两条正控**单独不能证明**任何事 —— 真正防"函数只会返回原串"的是 §4-**B** 那发（返原串 ⇒ 3 条红）。正控只钉"该不动的不动"，不冒充承重证明。

## 6 ★★ staging 跑法（完整命令）＋ **为什么 `PATH` 不够**

```bash
# 1) 模型：仓内 third_party/model-fetch 的 .onnx 是坏件（§9-①），用完好压缩包解开，落仓外
tar -xjf third_party/model-fetch/punc-int8.tar.bz2 -C /d/tmp/wisp/20261026/27r1/punc
export WISP_PUNC_MODEL='D:\tmp\wisp\20261026\27r1\punc\sherpa-onnx-punct-ct-transformer-zh-en-vocab272727-2024-04-12-int8\model.int8.onnx'
# 2) 编译测试二进制，并把三枚 sherpa DLL 放到它自己旁边
go test -c -o D:\tmp\wisp\20261026\27r1\bin\speech.test.exe ./internal/speech
copy third_party\sherpa-onnx\*.dll D:\tmp\wisp\20261026\27r1\bin
# 3) 另两个模型路径（staging 目录里没有 third_party，必须给）
export WISP_SILERO_MODEL='<repo>\third_party\spike-models\silero_vad.onnx'
export WISP_SPIKE_MODELS='<repo>\third_party\spike-models'
# 4) 就地跑（cwd = 二进制所在目录）
cd D:\tmp\wisp\20261026\27r1\bin && .\speech.test.exe -test.v      # => 44/44 (isolated) / 55/55 (with ticket-26 files), rc=0
```

★★**为什么 `PATH` 不够**：Windows 对**依赖 DLL** 的搜索序是「**主模块所在目录** → System32 → 系统目录 → 当前目录 → `PATH`」。`go test` 把二进制建在临时目录，那里没有 `onnxruntime.dll` ⇒ 加载器走到 **System32 拿 1.17.1**；`PATH` 排在 System32 **之后**，顶到最前也**改不了这个序** ⇒ `OrtGetApi(28)` 返回 NULL、C API 解引用、整个二进制崩。**staging 把 1.28.2 放进"主模块所在目录"**，那是唯一排得比 System32 前的位子。
★**三层口径（本片现跑）**：①什么都不给 ⇒ `0xc0000135`（进程起不来）②只给 `PATH` ⇒ 拿 System32 的 1.17.1，`checkOnnxruntime` **具名拒绝→skip**（rc=0，⛔ 红）③**staging ⇒ 全绿**。

## 7 ★ 四数 ＋ 零字节

| 面 | 命令 | rc | 读数 |
|---|---|---|---|
| `internal/speech` **改前** | `git ls-tree -r --long HEAD -- internal/speech` | — | **21 文件／172,823 B** |
| `internal/speech` **改后**（盘上） | `Get-ChildItem -Recurse -Filter *.go` | — | **30 文件／251,952 B** |
| ★**其中我的** | — | — | **＋2 文件＋39,842 B**（`punctuate.go` 12,807／`punctuate_test.go` 27,035）＋`transcriber.go` **＋2,524 B** ＝ **我的净增 42,366 B** |
| ★**其余 ＋7 文件／36,763 B** | — | — | 🔴 **⛔ 我的**：并存腿**票 26（TTS＋slot）**在 18:45–18:52 落进同一目录（`tts.go`／`slot.go`／`rss_other_test.go`／`rss_windows_test.go`／`slot_test.go`／`slot_real_test.go`／`tts_test.go`）。★**本片 ⛔ 碰它们**，也 ⛔ 在 commit 里带它们。 |
| `internal/speech` staging | §6 | **0** | **44 PASS／0 FAIL／0 SKIP**（只有我的 13 案时，18:42）→ **55 PASS／0 FAIL／0 SKIP**（含票 26 的 11 案，18:52） |
| `internal/audio`（**零字节、独立跑**） | `git status --porcelain -- internal/audio`＝**空**；`go test -count=1 ./internal/audio` | **0** | `ok … 16.587s` |
| `cmd/wisp` **定向**（`PATH` 带 `third_party\sherpa-onnx`） | `go test -count=1 -run TestAC3LogSinkLandsInsideTheEnvDataRoot -v ./cmd/wisp` | **0** | `--- PASS … 0.01s`；`ok … 0.166s` |

★**未改件逐枚 `git hash-object` 同值**：`engines.go`／`engines_test.go`／`asr.go`／`vad.go`／`endpoint.go`／`doc.go`／`fixtures_test.go`／`asr_test.go`／`vad_test.go`／`endpoint_test.go`／`transcriber_test.go`／`cgo_boundary_windows_test.go`（12 枚）。

## 8 ★ 格式甲乙两尺 ＋ 门禁 rc

| 尺 | 输出 |
|---|---|
| 甲 `gofmt -l internal/speech` | **空**（rc 0） |
| 乙 `gofumpt -l internal/speech`（`D:\work\base\gopath\bin\gofumpt.exe`） | **空**（rc 0） |
| 甲／乙 全仓（`internal cmd tools scripts frontend`） | 两尺**逐字相同 3 行**：`internal\agent\approval\pending_read.go`／`internal\agent\tools.go`／`cmd\wisp\models.go` ⇒ **已知 3 条既存假红**（CRLF／`autocrlf` 造），**⛔ 本片责任**；**我的 3 枚两尺皆空** |
| `go vet ./internal/speech` | **rc 0** |
| 门禁 `d22scan`（`cd tools/d22scan && go run . -root <repo>`） | **先 rc 1／后 rc 0**：第一发唯一 finding 是 **`internal/speech/punctuate.go:132: [phantom-citation]`** —— ★它抓的是**本文件当时还不存在**（ban#9，票 212）。本件落盘后复跑 **rc 0 `clean`**。★分母：`ban#8 internal/`=**679**、`bans#1-5 internal/`=**263**、`cmd/`=155／ban#8 `cmd/`=155（各含票 26 的在飞文件）。 |

★**⛔ 没用 ban#8 去扫 `docs/`**（它只扫 `.go`，那是错尺）。

## 9 ★ 三条发现 ＋ 边界（⛔ 越界即作废，逐条自陈）

1. 🔴🔴 **①仓内的标点模型是坏件，且"没模型"不等于"动不了"**：`third_party/spike-models/` **没有**标点模型；`third_party/model-fetch/punc-model.int8.onnx` 是 **29 B 的 `Invalid username or password.`**、`punc-tokens.txt` 是一段 HTTP-401 JSON（`fetch.log` 里 `curl: (22) … 401`）。★★**但** `punc-int8.tar.bz2` **完好**：解开后 `model.int8.onnx` sha256 `65a3fb9f5ad7bfb9…`／75,519,198 B、`tokens.json` `c960ab87bccea4aa…`、压缩包本身 `c0d5aa5f8eeb6860…` —— **与 `build/models/manifest.json` 的 pin 三处逐字相同**。⇒ 我**按"模型在"继续**（停在目录名上等于拿形式换价值），并把它记成本片第一条读数。★产品侧该模型**早已登记**（manifest `id=punc-ct-transformer-zh-en-vocab272727`、`purpose: punctuation`、`status: ok`）⇒ 缺的从来不是模型，是**用它的代码**。
2. 🔴🔴 **②给 sherpa 一个非 onnx 文件 ⇒ 进程直接死，⛔ 返回 nil**：本片第一版"加载失败降级"测试**把整个测试二进制打崩了** —— `Exception 0xe06d7363`（MSVC C++ 异常），`signal arrived during external code execution`，Go **无法 recover**。⇒ 光判 `punc == nil` 的降级是**假的**。故加 `checkOnnxModelFile`：先 `os.Stat`（存在／常规文件／非 0 字节）、再读首字节必须 `0x08`（ONNX `ModelProto` 的 `ir_version` 字段 tag；本树三枚真模型全是 `08 07`，而那枚 29 B 坏件首字节 `0x49`＝`'I'`）⇒ **这正是本仓自己的坏件形状**，四类坏件现在各给具名错误（§1-AC4 读数）。
3. **③英文口述会拿到中文标点**（§5）：模型 `punc_list` 只有 `。，？、`。⇒ 具名残余：英文标点要么换模型、要么按 `unicode.Is(unicode.Han)` 加语言门，**本片 ⛔ 做**（票面场景③是中文口述）。
4. **④AC5 的 Warm 驻留读数拿不到**：采样器在 `cmd/wisp`／`internal/models`，**两处禁写**。本片只给结构事实（一实例、`Close` 释放、无 per-utterance 装载）。
5. **边界：⛔ 改 `engines.go` 签名**（§3 已说明，改 2 个新文件实现）；**⛔ 加配置键**——`punctuationThreads`／`PunctuationLatencyBudget`／`punctuationAlphabet` 全是**包内常量**，`[voice] punctuation` 这个开关要动 D36 配置模型，**留单开票**；`OpenPunctuation(cfg, enabled bool)` 把"开/关"做成**实参**，测试直接覆盖两条路。
6. **零裸 `go func(`／零墙钟超时**：本片零 goroutine；`PunctuationLatencyBudget` 只用于**测量**，代码里没有任何按时间放弃的路径（D42#9）。
7. **新 cgo 面 0**：只新增 `sherpa.NewOfflinePunctuation`／`AddPunct`／`DeleteOfflinePunc`，全在 `cmd/wisp` 早已吃进去的同一条 cgo 边界内。
8. **未做**：把标点接进 `residentVoice` 装配（`cmd/wisp`，禁写面）、模型热替换、英文门、Warm 驻留采样、`docs/SLO.md` 附录（仓里不存在）。
