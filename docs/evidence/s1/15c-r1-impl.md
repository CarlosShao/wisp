# 票 15c 流式 ASR 腿 `15c-r1` —— 落地件（写腿）

> **尺**：`HEAD=1f2bae2a7`（`dev`），现跑 2026-10-26 16:31–16:56 +08。★**九项读数逐条在 §1–§9**，全为现跑，⛔ 转述。
> **写面**：`internal/speech/**` 新增 **4 文件、39,257 B** —— `asr.go` 7,465／`transcriber.go` 7,463／`asr_test.go` 9,622／`transcriber_test.go` 14,707，＋本文件。
> **零字节**：`doc.go`／`endpoint.go`／`vad.go` 与三枚 `15-b` 测试件（`vad.go` 树内 blob `af5fc14a…` **＝ HEAD 同值**）、`cmd/**`、`internal/**` 其余、`scripts/**`、票面、台账。
> **临时件（只建不删）**：`D:\tmp\wisp\20261026\15cr1\{bin,probe,logs,overlay,mut}`。★**两发变异全走仓外件 ＋ `-overlay`** ⇒ **工作树零突变窗口**。

## 0. 一句话

`internal/speech` 长出**流式 ASR 会话**：`asr.go`（流式 Paraformer int8，一个 recognizer 常驻 ＋ 每句租一条 stream）＋ `transcriber.go`（**把 `15-b` 的 `Detector` 判停与本 ASR 焊成一条缝**：句内吐 partial、判停出 final）。★**判停仍在 Go**：`EnableEndpoint: 0`、`NumThreads: IntraOpThreads` —— 与 `15-b` 把 silero 的 `min_silence/min_speech` 钉 0 是**同一条理由**（C 侧不许多一个沉默的真相源）。

## 1. ★★★ partials 真的流出来 —— 时序表

夹具＝`15-b` 的 `synthSpeech`（120Hz 谐波栈、3.5Hz 包络）3s，逐 **512 样本＝32ms** 帧经 `Transcriber.Accept` 喂入。读数原样来自 `TestPartialsStreamWhileTheUtteranceIsOpen` 的 `t.Logf`：

| 第 N 个 32ms 帧 | fixture 时刻 | 吐第几个 partial | 文本 | **相对本句首样本** | 累计喂入 |
|---|---|---|---|---|---|
| **44** | 1408ms | **#1** | `哦` | **640ms** | 10,240 |
| **63** | 2016ms | **#2** | `哦哦哦哦` | **1248ms** | 19,968 |
| **82** | 2624ms | **#3** | `哦哦哦哦哦哦` | **1856ms** | 29,696 |

末态：`final "哦哦哦哦哦哦" over 35328 samples`。★**三条性质**：①**partial 在句子还开着的时候就到了**（#1 时本句只喂了 640ms）；②**只增不减**（1→4→6 字，测试里逐对断言"后一个不比前一个短"）；③**节奏 608ms＝19 窗**，那是 **Paraformer 的 chunk 粒度**（模型属性，不是我能调的旋钮）。

★**本句起点可由 `fed` 反推（不是估的）**：`Begin` 那一帧一次性回放 `ep.Total＝5120` 样本（10 窗），此后每帧 ＋512 ⇒ 帧 44 时 `10240 ＝ 5120 ＋ 10×512` ⇒ **起点＝第 34 帧（1088ms）**；帧 63 校核 `5120 ＋ 29×512 ＝ 19968` ✅。⇒ ★**silero 在这条合成夹具上要把 300ms 最短语音攒够，花了 1088ms 才开句**（`15-b` 记的"silero 过报 ~400ms"的直接后果，归 `15-e` 的球延迟算总账）。

## 2. ★★ 首字延迟 ≤800ms（两种口径都给）

| 口径 | 读数 | 预算 | 判 |
|---|---|---|---|
| **首个 partial**（相对本句首样本，**音频钟**，确定性） | **640ms** | 800ms（`SPEC-04:114` 流式档） | ✅ 余量 **160ms** |
| **VAD 判停 → final**（`SPEC-04:114` 自己的措辞，墙钟） | **1,048,600 ns ＝ 1.05ms** | 800ms | ✅ 余量 ~760× |
| 同上，但先给个尺子自证 | 本机单调钟**量化 ~0.57ms**（连读 `time.Now()` 差值只有 0 与 568500ns）⇒ 单发读数 <1 tick 一律记作"低于 tick" | —— | 已核 |

★**为什么"判停→final"几乎为 0**：流式解码是**边到边解**的，判停发生时音频早已解完，`Finish` 只需 flush 特征器 ＋ 取结果。★**这不是没牙**：同一行代码在整段变异下读 **214.6ms**（§4 发 B）。★**首字延迟用音频钟量、⛔ 用墙钟**（D42#9）⇒ 同一夹具在任何机器上都是 640ms。

## 3. ★★ 词错率（合成夹具上；**⛔ 真实语料**）

★**必须先说清口径**：合成夹具**没有词汇内容**，所以"WER"只能退化成两种可算的量，两种都报：

| 夹具 | 参考（合成器的标称脚本） | 假设 | S | D | I | WER |
|---|---|---|---|---|---|---|
| `E-vowel-a`（/a/ 共振峰，3.5 音节/s ⇒ **10** 音节） | `啊啊啊啊啊啊啊啊啊啊` | `啊啊啊啊啊啊啊啊` | 0 | 2 | 0 | **0.20** |
| `F-vowel-i`（/i/ 同理 ⇒ 10 音节） | `一一一一一一一一一一` | `一一一一一一一一` | 0 | 2 | 0 | **0.20** |
| **脚本组合计（E,F）** | **20 词** | —— | **0** | **4** | **0** | **0.20** |
| `A-harmonic-stack` | （无脚本） | `哦哦哦哦哦哦哦哦` | —— | —— | **8** | 插入率 8 |
| `B-silence` | （无脚本） | `""` | —— | —— | **0** | 0 |
| `C-white-noise` | （无脚本） | `""` | —— | —— | **0** | 0 |
| `D-speechband-noise` | （无脚本） | `""` | —— | —— | **0** | 0 |

★**这不是词级准确率**：真词级 WER 要真语料，那是 **`15-d`** 的基线集（`SPEC-10:57-60`），本片**明收在写面外**。★能钉住的是两件真事：①模型对**有标称音素的合成音**给出一致的元音（/a/→啊、/i/→一，**0 替换**），只是**少数 2 个音节**；②对**无音素内容**的静音／白噪／语音带噪，**一个词都不吐**（3 夹具 × 0）。对齐用的是测试内 20 行 Levenshtein（仅为让上表可复算），**⛔ 不是 CER 量具**。

## 4. ★★ 反向判据：三读数 ＋ 还原核（**时机在最后一发变异之后**）

| 发 | 变异（仓外 `mut\` ＋ `-overlay`） | rc | 红的断言（逐字） |
|---|---|---|---|
| **A** | `IntraOpThreads` **1→4**（`vad.go`，`15-b` 的常量） | **1** | **3 条红**：★**我的** `TestAsrRecognizerConfigPinsTheShape`「`ModelConfig.NumThreads = 4, want the literal 1`」／`15-b` 的 `TestIntraOpThreadShape`「`intra_op=1 burned 3.79 cores (cpu/wall), so the pin is not holding`」／`TestPinnedIntraOpThreads` |
| **B** | **流式 → 整段**（`asr.go`：`Feed` 只攒不喂、`Finish` 才全部灌进去，4 处 7 行） | **1** | **3 条红**：`TestPartialsStreamWhileTheUtteranceIsOpen`「**`0 partial(s) over 3s of speech, want at least 2`**」／`TestConsumeDrainsTheSeamUntilTheChannelCloses`／`TestSpeakPauseSpeakYieldsTwoFinals`「`no partials were produced while the utterances were open`」 |
| **C** | **还原**（无 overlay，重编重跑） | **0** | **28 PASS ／ 0 FAIL ／ 0 SKIP** |

★★**发 B 的"首字延迟必须爆"（这就是"流式真的在流"的证明）**：

| 量 | 流式（C 行） | 整段（B 行） | 倍数 |
|---|---|---|---|
| 首个 partial（音频钟） | **640ms** | **没有** —— 整句期间 0 个 | ∞ |
| 首个文本出现在**句中**？ | 是（句长 1788ms 时 640ms 就有） | **否**（句闭才出） | —— |
| 首字延迟（句起点→首文本） | **640ms** ✅ | **1788ms（句闭）＋ 214.6ms（flush）＝ 2002.6ms** ❌ | **3.1×** |
| 判停→final 的 flush 腿 | **1.05ms** | **214,595,600 ns ＝ 214.6ms** | **205×** |
| 对 800ms 预算 | 0.80× | **2.50×（爆）** | —— |

★**发 A 没红的那一半要说明**：`IntraOpThreads` **只**由 `asrRecognizerConfig` 一个点消费，所以变异只打中"接线"和"观测量"两条，**不影响 partial 时序**（B 行里时序与 C 行一字不差 ⇒ 也证明时序是确定性的，不是抖出来的）。

★**还原核（时机：`mutB3` 之后）**：①`restore3` 单独重编**无 overlay** 的二进制跑全套 ⇒ **rc=0、28/28**；②`vad.go` 的 `git hash-object` ＝ `af5fc14a…` **＝ `HEAD:internal/speech/vad.go`**，sha256 `E6AD28EF7423F888…5D0FF657` 与 `15-b` 记的 `E6AD28EF…F657` **逐字同值**；③`asr.go` 树内 sha256 `51B6547152CC90B5…` ≠ 变异件 `311DED2BA4F15CB1…`；④树内 `MUTANT` 标记 **0 命中**、`IntraOpThreads = 4` **0 命中**。

## 5. ★★ "⛔ 碰 `cmd/wisp`" 的证明

```
git show --name-only --format= <本片 commit>   => 恰好 5 行（4 枚 .go ＋ 本文件）
git show --name-only --format= <本片 commit> | grep -c '^cmd/'   =>  0
git status --porcelain -- cmd/                  =>  空
```

★**并且本片不需要 `cmd/wisp` 就拿到了全部九项读数**（未触发停手②）。★`internal/speech` 仍是**零 importer**：接线（`models.go:25` 的 `handOffModel`）在 **`15-e`**，本片不碰。

## 6. ★★ staging 跑法 ＋ **为什么 `PATH` 不够**

```
# 1) 编译测试二进制，并把三枚 sherpa DLL 放到它自己旁边
go test -c -o D:\tmp\wisp\20261026\15cr1\bin\speech.test.exe ./internal/speech
copy third_party\sherpa-onnx\*.dll D:\tmp\wisp\20261026\15cr1\bin
# 2) 两个模型路径（staging 目录里没有 third_party，所以必须给）
set WISP_SILERO_MODEL=<repo>\third_party\spike-models\silero_vad.onnx
set WISP_SPIKE_MODELS=<repo>\third_party\spike-models
# 3) 就地跑（cwd = 二进制所在目录）
D:\tmp\wisp\20261026\15cr1\bin\speech.test.exe -test.v      =>  28/28, rc=0
```

★★**为什么 `PATH` 打不过 `C:\Windows\system32`**：Windows 对**依赖 DLL** 的搜索顺序是「**主模块所在目录** → System32 → 系统目录 → 当前目录 → `PATH`」。`go test` 把二进制建在临时目录里，那里没有 `onnxruntime.dll` ⇒ 加载器一路走到 **System32，拿到 1.17.1**，`PATH` 排在 System32 **之后**，把 `third_party\sherpa-onnx` 顶到最前面也**改变不了这个序**。1.17 的 `OrtGetApi(28)` 返回 NULL，C API 直接解引用 ⇒ **`0xc0000005` 崩整个测试二进制**。**staging 把 1.28.2 放进"主模块所在目录"**，那才是唯一排得比 System32 前的位子。现跑实证：staging 的映像路径 `D:\tmp\…\bin\onnxruntime.dll`（`FileVersion 1.28.2`），System32 那枚 `1.17.260311-1434.1…`。
★**裸跑（无 staging）不是红**：`checkOnnxruntime()` 前置拒绝 ⇒ **具名 skip**，rc=0、**10 过 / 18 skip / 0 崩**，skip 原文「`speech: onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires…`」。★注：**连 `PATH` 都不给**时是 `0xc0000135`（`STATUS_DLL_NOT_FOUND`，sherpa 的 C API 本身都找不到）—— 那是另一回事，⛔ 与 ORT 遮蔽混为一谈。

## 7. ★★ 正控：纯静音不产出 partial

| 层 | 夹具 | 读数 |
|---|---|---|
| **整条缝**（`Detector`＋`Transcriber`） | 3s **数字静音**，逐 32ms 帧 | **0 partial ／ 0 final ／ 0 transcript**（`endpoint` 也不开） |
| **纯 recognizer**（绕过 VAD） | 同上，直接 `StreamingAsr` | **每帧 partial 都为空**，`final ""` |

★**双保险且两层都实测**：静音**开不了句**（`15-b` 的状态机），就算开了句，模型**也不吐词**。★★**⛔ 把静音解出幻觉词**这条钉子有两颗。

## 8. ★ 四数

| 面 | 命令 | rc | 读数 |
|---|---|---|---|
| `internal/speech` **改前** | `git ls-tree -r --long HEAD -- internal/speech` | —— | **7 文件 ／ 49,825 B** |
| `internal/speech` **改后** | 同上 ＋ 盘上 | —— | **11 文件 ／ 89,082 B**（**＋4 文件 ／ ＋39,257 B**，`15-b` 的 7 枚**一字节未动**） |
| `internal/speech` staging | §6 | **0** | **28 案 28 过 ／ 0 skip** |
| 同上**裸跑**（`PATH` 带 sherpa、无 staging） | `go test -count=1 ./internal/speech` | **0** | **10 过 ／ 18 具名 skip ／ ⛔ 崩** |
| `internal/audio`（**零字节改动、独立跑**） | `git status --porcelain -- internal/audio` ＝ **空**；`go test -count=1 ./internal/audio` | **0** | `ok … 15.946s` |
| `cmd/wisp` **定向**（`PATH` 带 sherpa） | `go test -count=1 -run TestAC3LogSinkLandsInsideTheEnvDataRoot -v ./cmd/wisp` | **0** | `--- PASS` |

## 9. ★ 格式甲乙两尺 ＋ 门禁

| 尺 | 命令 | 输出 |
|---|---|---|
| 甲 `gofmt` | `gofmt -l internal/speech` | **空** |
| 乙 `gofumpt`（绝对路径） | `D:\work\base\gopath\bin\gofumpt.exe -l internal/speech` | **空** |
| 归属判据（`[\\/]`） | 两尺输出按 `[\\/]cer[\\/]` 切：**命中 cer/ 的 0 条、我 4 枚 0 条** | ★中途一次 `gofumpt` 曾报 `internal\speech\cer\norm_test.go` —— **那是并存腿 `15-d` 的在飞文件**（不在本片写面、不在本片 commit），末次复跑已空 |
| `go vet ./internal/speech` | —— | **rc 0** |
| `d22scan` | `go run . -root <repo>`（`tools/d22scan`） | **rc 0** `clean - no D22 ban violations`；`internal/` ban#1-5 257 枚 ＋ ban#8 659 枚 Go 全覆盖 |

★**⛔ 没用 `d22scan` 的 ban#8 去扫 `docs/`**（它只扫 `internal/`＋`cmd/` 的 `.go`）—— 那是错尺，本片未犯。
★**⚠ 已知 3 条既存假红**（`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`cmd/wisp/models.go`）＝CRLF/`autocrlf` 造的，**⛔ 本片责任**：本片 4 枚新文件 `gofmt`/`gofumpt` 双尺皆空。

## 10. 未做／边界（⛔ 越界即作废，逐条自陈）

- **AC5 的"一次重试 → `asr` 错误 →『没听清』"策略未做**：空洞穴已有载体（静音句出来的是 `final ""`，策略按它键），但**重试与错误类映射**不在本片派单九项读数内，留 `15-e`／续片。
- **⛔ 加配置键**：800ms 预算、80 维特征、`greedy_search` 全是**包内常量**；`[voice].asr.*` 要动 D36 配置模型，留单开票。
- **未做**：模型热替换、真语料 CER（`15-d`）、`15-e` 的 `handOffModel` 接线与球走态、TTS（26）、KWS（41）、AEC（59）。
- 包内**零 goroutine**（`Transcriber.Consume` 在调用者 goroutine 跑，取消靠 ctx）、**零裸 `go func(`**、**零墙钟超时**（时长一律按样本计，音频钟）—— 三句都是 `d22scan` rc 0 ＋ 现跑 grep，不是声明。
- **新 cgo 面 0**：只新增 `sherpa.NewOnlineRecognizer`／`NewOnlineStream`／`AcceptWaveform`／`Decode`／`GetResult`／`Delete*`，全在 `cmd/wisp` 早已吃进去的同一条 cgo 边界内。

## 11. ★★ 两条对新下游腿（`15-e`／`15-d`）有直接后果的发现

1. ★★★**Paraformer 的 chunk 粒度是 608ms（19 窗）**，且**一个词都不会在 ~600ms 之前出来** ⇒ **首字延迟的下界是模型给的 640ms，不是代码能压的**。`15-e` 算球延迟（`Listening→Thinking`）时**必须**把「silero 开句迟 ~1088ms（合成夹具实测）＋ chunk 640ms」两段加进去，⛔ 只算后者；★这两段**都是实测**，且**都会**随真语料变（开句迟那一段尤其），`15-e` 复测时别把它当常量。
2. ★★**判停→final 只有 ~1ms**（流式边到边解完了） ⇒ `SPEC-04:114` 那条「VAD-stop→ASR-final ≤800ms」**在本引擎上几乎没有成本**，真正的延迟全在 ①silero 攒够 300ms ②chunk 640ms。**离线档（1500ms）省不下这条**，反而会把首字延迟推到句尾 —— 这是 `15-a1` 选流式的实测支撑。
