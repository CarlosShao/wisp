# 票 15b VAD 端点判停腿 `15b-r1` —— 落地件（写腿）

> **尺**：`HEAD=4569ac061`（`dev`），现跑 2026-10-26 11:44–12:0x +08。★**八项读数逐条在 §1–§8**，全为现跑，⛔ 转述。
> **写面**：`internal/speech/**`（新建 6 文件、`doc.go` **零字节**）＋ 本文件。`cmd/**`、`internal/**` 其余、`scripts/**`、`docs/specs/**`、票面、台账 **零字节**。
> **临时件（只建不删）**：`D:\tmp\wisp\20261026\15br1\{probe,bin,logs,overlay,mut}`。★突变全走仓外件＋`-overlay`，**无工作树突变窗口**。

## 0. 一句话

`internal/speech` 从「21 行 `doc.go`、零导出符号、零 importer」变成 **VAD 判停包**：`endpoint.go`（700/300/60s 状态机，纯 Go、无 cgo）＋ `vad.go`（silero 会话＋seam 消费）。★**判停策略全在 Go**，C 侧 `min_silence`/`min_speech` 钉 0 只当逐窗分类器 ⇒ **三个常量是唯一的真相源**，也才可被变异证伪。

## 1. ★★★ AC3 原文 ＋ 逐条判定

票面 `.scratch/wisp/issues/15-speech-engines-cer-harness.md:49-50` 逐字：

> `- [ ] Endpointing: speak-pause-speak fixture splits correctly; <300ms speech ignored; 60s utterance auto-finalizes.`

★与 `15-a1` 转述（`15-a1-census.md:25`）**一致** ⇒ 未触发停手③。

| # | AC3 条 | 我的读数（现跑） | 判 |
|---|---|---|---|
| ① | speak-pause-speak 正确切分 | 状态机：2 开 2 关，关在 **1600ms / 2900ms**，各 `total=1200ms speech=500ms`；真模型：gap 1500/2000ms → `[start end start end]`；sweep 500/900→1 发，1100/1500/2000→2 发 | ✅ |
| ② | `<300ms` 语音被忽略 | 299ms 突发 → **0 转移**；300ms → 1 开；301ms → 1 开 | ✅（口径见 §2 注） |
| ③ | 60s 语句自动判停 | 59.9s → 0 强制；60.0s 起 → 1 强制**恰在 60000ms**；真模型喂 65s → `[start max-duration start end-of-stream]`，关闭件 `Duration=60000ms`、`Samples=960000` | ✅ |

## 2. ★★ 端点状态机三发边界（**两侧都测**）

三类各一发，均**合成夹具**、⛔ 真实语料；读数来自 `TestEndpointSilenceBoundary` / `TestMinSpeechBoundary` / `TestMaxUtteranceBoundary`（1ms 步进，采样级）。

| 常量 | 值／样本数 | 一侧（不触发） | 另一侧（触发） | 命中点 |
|---|---|---|---|---|
| `EndpointSilence` | 700ms ／ 11200 | **699ms → 0 关** | **701ms → 1 关** | 恰 700ms 也关（`>=`，含等号）；关在 **1700ms** |
| `MinSpeech` | 300ms ／ 4800 | **299ms → 0 开**（0 转移） | **301ms → 1 开** | 恰 300ms 开 |
| `MaxUtterance` | 60s ／ 960000 | **59.9s → 0 强制**（`InSpeech` 仍真） | **60.1s → 1 强制** | 恰 60.000s；`Total` 与 `Samples` 都被截到 960000 |
| 正控 | 纯静音 | **5s 纯静音 → 0 转移**（`Finalize` 后仍 0） | —— | §5 |

★**注（AC3② 的口径，必须写清）**：真模型那条腿对这一格**是空的** —— 实测 silero 对 150ms／250ms 的合成突发**根本不报语音**，所以端到端「没有端点」不能证明 300ms 门在起作用。300ms 门**只在状态机那一层是承重的**（变异 300→0 ⇒ 299ms 立刻开话，§4 发 B）。⛔ 把端到端那条当 AC3② 的证据。

★**32ms 量化的算术**（不是测量）：真 seam 窗口 = 512 样本 = 32ms ⇒ 关的判断粒度 32ms，21 窗 = 672ms < 700 不关、22 窗 = 704ms ≥ 700 关。采样级精确性是状态机的性质，**不是** seam 的可观测分辨率。

## 3. ★★ cgo 边界：**单 onnxruntime 实例** ／ **intra_op 1**

★**证法（`cgo_boundary_windows_test.go`，三条腿 ＋ 一条仪器正控）**，以及**为什么它有效**：

| 腿 | 观测量 | 为什么有效 | 读数 |
|---|---|---|---|
| 1 | 进程已映射映像里 base name 含 `onnxruntime` 的张数（`EnumProcessModules`＋`GetModuleFileNameEx`） | 「第二个 onnxruntime 实例」若成立，必然是**第二张映射映像（另一条路径）或另一个进程**；前者在这把尺上直接现形 | **1 张**，路径 `D:\tmp\wisp\20261026\15br1\bin\onnxruntime.dll`（＝仓里 1.28.2）；sherpa 映像 **1 张** |
| 2 | 连续 **8 次** `NewDetector`＋`Close` 后再数 | 每会话各自起运行时的实现会在这里长出第 2 张 | **仍 1 张，路径不变** |
| 3 | **仪器正控**：把同一 DLL 复制成另一文件名再 `LoadLibraryEx` 进本进程 | ⛔ 自证 —— 这一步证明**计数器看得见第二张**，否则腿 1/2 的「1」毫无意义 | **2 张** ⇒ 仪器有牙；`FreeLibrary` 后复原 |
| 4 | `intra_op` 形状：会话创建后线程增量 ＋ 负载期 CPU 秒/墙秒 | 线程池就是线程；单线程钉死时 **≤1 核可忙**，比值应贴 1 | `1`：**线程 +0**，cpu/wall **0.98–1.10**；`4`（正控）：**线程 +3**，cpu/wall **3.77–4.19** |

★**为什么这两条合起来算「单实例＋钉 1」**：腿 1–3 用「第二张映像」把**同一进程内另起运行时**这条唯一现实路径排除；腿 4 用**线程池**把 `intra_op>1` 排除，且 4 线程的正控**量级分离 4 倍**，证明尺子不是摆设。

★★**它不能看见什么（⛔ 含糊）**：C 侧 `Ort::Env` 的引用计数**从 Go 不可观测**，本腿**不宣称**看见它。补强的结构性事实：`NewDetector` 是**唯一**构造点（钉 `IntraOpThreads`），`internal/speech` 里 `"os/exec"` **0 命中**、裸 `go func(` **0 命中** ⇒ 不存在「第二个进程/第二条线程池」。

## 4. ★★ 反向判据：三读数 ＋ 还原核

突变件在仓外（`mut\m700\endpoint.go`、`mut\m300\endpoint.go`），经 `-overlay` 换入编译 ⇒ **工作树零突变窗口**。

| 发 | 变异 | rc | 红的断言（逐字） |
|---|---|---|---|
| **A** | `EndpointSilence` **700→100** | **1** | **5 条红**：`TestEndpointSilenceBoundary`「`pause 699ms: got 1 end(s), want 0; trace=start@300ms end@1100ms`」／`TestSpeakPauseSpeakSplits`（`total 600ms, want 1200ms`）／★**`TestVadPauseSplitPoint`「`900ms pause gave 2 utterance(s), want 1 (bridged)`」（真模型腿也红）**／`TestVadPauseDetail`／`TestPinnedIntraOpThreads`（`policy samples = 1600/4800/960000`） |
| **B** | `MinSpeech` **300→0** | **1** | **3 条红**：`TestMinSpeechBoundary`「`burst 299ms: 1 utterance(s) opened, want 0`」／`TestMaxUtteranceBoundary`／`TestPinnedIntraOpThreads`（`11200/0/960000`） |
| **C** | **还原**（无 overlay，重编重跑） | **0** | **16 PASS / 0 FAIL / 0 SKIP** |

★**还原核时机**：在**最后一发变异（B）之后**单独重编无 overlay 的二进制跑全套 ⇒ C 行即核。另加两条独立核：盘上 `= 100 * time.Millisecond` / `= 0 * time.Millisecond` 残留 grep **0**；`endpoint.go` sha256 `FC846D16…AD86`、`vad.go` sha256 `E6AD28EF…F657` 与变异前同值。

★**A 的行为解读**：判停阈值调小 ⇒ 同一句话被切成两半，699ms 那一发**该不关却关了** —— 这正是「700 这个常量真在起作用」。B 同理：`MinSpeech=0` ⇒ 299ms 就开话。

## 5. ★★ 正控：静音进来不产出端点

| 层 | 夹具 | 读数 |
|---|---|---|
| 状态机 | 5s 纯静音（1ms 步进）＋ `Finalize` | **0 转移**，`InSpeech()=false` |
| 真模型（seam 全通） | 3s 数字静音经 `WavInjector`→`EncodeFrame`→`DecodeFrame`→silero | **0 端点**（`transitions: []`） |

★⛔ 把纯静音判成一句话。★诚实边界：这条正控**钉不住 `MinSpeech`**（变异 300→0 时它仍绿 —— 静音走的是 idle 的 `!isSpeech` 早退），钉它的是 §2 的 299/301。**两条合起来才覆盖**。

## 6. ★★ 「⛔ 碰 `cmd/wisp`」的证明

```
git show --name-only <本腿提交> | grep -c '^cmd/'     =>  0
git grep -n '"github.com/CarlosShao/wisp/internal/speech"'  =>  0 命中（零 importer）
```

★`internal/speech` **零 importer** 是硬事实：`models.go:25` 的 `handOffModel` 是 `package main`，接线在 `15-e`。本片**不需要**动 `cmd/wisp` 就拿到了全部读数（未触发停手②）。

## 7. ★ 四数

| 面 | 命令 | rc | 读数 |
|---|---|---|---|
| `internal/speech`（**新增包**） | `speech.test.exe -test.v`（staging，见 §9.1） | **0** | **16 案 16 过 / 0 skip** |
| 同上，裸 `go test ./internal/speech` | `go test -count=1 ./internal/speech` | **0** | **8 过 / 8 skip**（skip 全是需真会话的，原因＝§9.1 的 ORT 遮蔽；**⛔ 崩进程**） |
| `internal/audio`（**零字节改动**） | `git status --porcelain -- internal/audio` ＝ **空**；`go test -count=1 ./internal/audio` | **0** | `ok … 15.952s` |
| `cmd/wisp` **定向**（`PATH` 带 sherpa） | `go test -count=1 -run TestAC3LogSinkLandsInsideTheEnvDataRoot -v ./cmd/wisp` | **0** | `--- PASS` |

★`GOOS=linux CGO_ENABLED=0 go vet`：`internal/speech` rc 1 —— ★**与未改动的 `cmd/wisp` rc 1 逐字同因**（`sherpa-onnx-go-linux: build constraints exclude all Go files`）⇒ **我没引入新的破法**；`internal/audio` 同配置 rc 0。

## 8. ★ 格式甲乙两尺 ＋ 门禁

| 尺 | 命令 | 输出 |
|---|---|---|
| 甲 `gofmt` | `gofmt -l internal/speech` | **空** |
| 乙 `gofumpt`（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`） | `gofumpt -l internal/speech` | **空**（两尺都空 ⇒ 无需用 `[\\/]` 判归属） |
| `go vet` | `go vet ./internal/speech` | **rc 0** |
| `d22scan` | `go run . -root <repo>`（`tools/d22scan` 模块） | **rc 0** `clean - no D22 ban violations`；`internal/` 252 production ＋ ban#8 652 Go 文件全覆盖 |

★⛔ 用 `d22scan` 的 ban#8 扫 `docs/`（它只扫 `internal/`＋`cmd/` 的 `.go`）—— 那是错尺，本片未犯。

## 9. ★★ 两个新发现（对下游腿有直接后果）

1. ★★★**`C:\Windows\system32\onnxruntime.dll`（1.17.1）遮蔽仓里的 1.28.2**，而 **PATH 打不过 System32**（Windows 对依赖 DLL 的搜索顺序：应用目录 → System32 → … → PATH）。后果：裸跑会拿到 1.17，`NewVoiceActivityDetector` **`0xc0000005` 崩整个测试二进制**（C API 拿到 NULL `OrtApi` 直接解引用）。已做两件事：①生产码 `checkOnnxruntime()` **前置拒绝**并说明原因（返回错误，不崩）；②**staging 跑法**（下面）。★这条对 `15-c`／`15-e` 同样成立：**要真跑会话的腿都不能只靠 PATH**。
   ```
   go test -c -o $bin\speech.test.exe ./internal/speech      # $bin 里放三枚 DLL
   copy third_party\sherpa-onnx\*.dll $bin
   $env:WISP_SILERO_MODEL=<repo>\third_party\spike-models\silero_vad.onnx
   & $bin\speech.test.exe -test.v                            # 16/16
   ```
2. ★★**silero 在合成夹具上有 ~400ms「不听话」**：信号停了仍报 ~720ms 语音、重新起报 ~300ms ⇒ **标称停顿要 ~1100ms 才凑够 700ms 判停**。实测 sweep（`TestVadPauseSplitPoint`，标称 gap → 发数）：`500→1`、`900→1`、`1100→2`、`1500→2`、`2000→2`。★夹具按**测出来的标定**，⛔ 按标称算；两条断言分别压 900（离 700 有 220ms 余量）与 1500（余 400ms）。

## 10. 未做／边界（⛔ 越界即作废，逐条自陈）

- `15-a` 引擎基座、`15-c` 流式 ASR、`15-d` CER 量具、`15-e` 装配 **一律未碰**。
- **⛔ 加配置键**：700ms 是包内常量（`EndpointSilence`）；`[voice].asr.silence_ms` 要动 D36 配置模型，留单开票。
- `doc.go` **零字节**：`DEFERRED(engines)` 原样留着（摘它属 `225-deferred-registry-audit` 第 22 行那条账，不在本片写面）。
- 无新 cgo 面（`cmd/wisp` 早已吃 cgo）；包内**零 goroutine**（`Consume` 在调用者 goroutine 跑，取消靠 ctx），**零墙钟超时**（时长一律按样本计，音频时钟）—— 两句都是实测 grep，不是声明。
- ★未做：`Flush` 之外的模型热替换、真语料 CER、`15-e` 的 `handOffModel` 接线。
