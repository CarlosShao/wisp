# 240-c1 只读普查：麦克风采集侧 / KWS 侧 / 接线落点 / 半成品名册

- 腿：`240-c1`（只读，不产码、不跑编译与测试）
- 取钟：`2026-09-30 11:26:55 +0800`（`date` 与 `mkdir` 同发，落笔在本节之前）
- 锚点：起手 `git log --oneline -1` = `0589fd9c`（球被窗口切成方形那枚）
- 红线自查：本腿零 `go build`／零 `go test`／零 `gofumpt`／零 `go vet`／零 `d22scan`；
  零读、零写、零转述 `frontend/**` 与 `design/**`；不改 `PLAN.md`／`docs/specs/**`／
  `docs/SLO.md`／`internal/ball/thresholds.go`／golden／`allowlist.txt`／台账／HANDOVER。

---

## §0 现量口径（我这轮跑了哪些尺）

> 每条尺都先在本机真跑，读数照抄。尺失效（rc=1 且期望非空）本身就是读数。

| # | 尺（命令要点） | 本机真实读数 | 口径说明 |
|---|---|---|---|
| R1 | `grep -rn 'CarlosShao/wisp/internal/audio' --include=*.go .`（排 `.scratch`） | **零命中，rc=1** | 全仓模块内没有任何文件 import `internal/audio` |
| R2 | 同上，`internal/speech` | **零命中，rc=1** | `internal/speech` 只有 `doc.go`，无人 import |
| R3 | `grep -rn 'NewWASAPIMicrophone\|WavInjector' --include=*.go .`（排 `.scratch`） | 命中全部落在 `internal/audio/` 自身：定义 + `hotplug_test.go:529` 等测试；`cmd/`、其它 `internal/` 包 **零命中** | 真设备栈有实现、无生产调用者 |
| R4 | `grep -rn 'SetAudioLevel' --include=*.go .`（排 `.scratch`） | 非测试生产者仅 `cmd/balldebug/main.go:418`、`:477`；定义 `internal/ball/liquid_windows.go:42`；测试 `internal/ball/live_windows_test.go:626/637/645/674` | 与编排者 09-30 当面认的缺口一致 |
| R5 | `grep -rn 'sherpa\.' --include=*.go .`（排 `.scratch`） | `cmd/wisp/` 只用 `GetVersion`/`GetOnnxruntimeVersion`（`doctor.go:60/81`、`main.go:165`）；`NewKeywordSpotter` **只出现在 `scripts/spike/`** | 生产码链了 sherpa 但只用来报版本号 |
| R6 | `grep -rniE 'keywordspotter\|kws' --include=*.go .`（排 `.scratch`） | 命中三类：`scripts/spike/{model-residency,speech-baseline}/main.go`（真建 spotter）、`internal/agent/approval/*`（否决通道占位）、`internal/statemachine/*`（事件与守卫已成型） | 待 §2 逐枚拆 |
| R7 | `cat scripts/spike/go.mod` | `module github.com/CarlosShao/wisp/scripts/spike` —— **独立模块**，不在主模块构建图里 | 所以 R5 的"生产零命中"与 spike 里有真代码并存，不矛盾 |
| R8 | `cat go.mod`（主模块） | 直依赖含 `github.com/k2-fsa/sherpa-onnx-go v1.13.8`；`sherpa-onnx-go-windows` 等三平台包是 indirect | cgo 那层＝上游模块，不是本仓 vendored 源码 |
| R9 | `ls third_party/sherpa-onnx/` | 三枚 DLL：`onnxruntime.dll`、`sherpa-onnx-c-api.dll`、`sherpa-onnx-cxx-api.dll` + `.cache-manifest.json`；**无 .go、无 .h** | 本仓不持有 sherpa 源码，只持有运行期 DLL |
| R10 | `sed -n '515,535p' docs/PLAN.md`（只读） | 空闲 ≤25MB/≤0.5%、待唤醒 ≤90MB/≤2% 两档原文确认在位 | 用于 §2 找代码里的对应常量 |

（§0 未完：后续节追加尺子编号从 R11 起。）

### §0 追加尺（§1/§2 现跑，读数照抄）

| # | 尺 | 本机真实读数 | 口径 |
|---|---|---|---|
| R11 | `grep -rn 'CarlosShao/wisp/internal/ball' --include=*.go .`（排 `.scratch`） | **命中 1 行**：`cmd/balldebug/main.go:27` | 悬浮球全包只有调试载体一个 import 者；`cmd/wisp` 不 import ball |
| R12 | `grep -rn 'Sqrt' --include=*.go internal cmd tools` | **零命中（rc=1）** | 主模块生产+测试码里没有任何 RMS/开方运算 |
| R13 | `grep -rn '\.Voice\.' --include=*.go internal cmd tools \| grep -v 'internal/config/'` | **零命中** | `[voice]` 配置被解析与热载，但没有任何非 config 包读过它 |
| R14 | `grep -rn 'C8\b\|AudioSource' --include=*.go internal cmd tools \| grep -v '^internal/audio/'` | **零命中** | C8 契约的名字与接口只在 `internal/audio` 自己体内出现 |
| R15 | `grep -rn 'EvSummon\|EvWakeWord\|EvKwsEnabled' --include=*.go internal cmd tools`（排 `_test.go`、排 `events.go`） | `EvSummon` 的生产触发只有 `cmd/balldebug/main.go:195`（热键回调）→ `:600` `dispatch(...EvSummon, nil)`；**`EvWakeWord`/`EvKwsEnabled` 除 `statemachine/events.go:18,20` 与 `table.go:55,61` 的定义/表行外零命中** | 唤醒词事件与 KWS 启用事件**今天没有任何生产者** |
| R16 | `grep -rn 'hotkey-listener' --include=*.go internal cmd` | 4 命中，全在名册声明/注释（`observe/goroutine.go:39,44`、`goroutine_test.go:245`、`cmd/balldebug/main.go:37`）；**无 Spawn** | D38b 六个常驻名里 `hotkey-listener` 今天也没落地 |
| R17 | `grep -rn 'Spawn("' --include=*.go internal cmd`（排 `_test.go`） | `audio-capture` 三处全在 `internal/audio` 内部（`audio.go:181`、`wasapimic_windows.go:83`、`wavinjector.go:84`）；`kws-infer`/`asr-infer`/`tts-infer` **零 Spawn** | 名册槽位声明了、没人占 |
| R18 | `wc -l third_party/model-fetch/x-kws/.../keywords.txt` + `head -8` | **8 行**，逐行是拼音音素 + `@中文词`（你好军哥/蛋哥蛋哥/小爱同学/你好问问/小艺小艺/小米小米/林美丽/你好西西），**无 "hi wisp"** | 出厂 keywords.txt 是 wenetspeech 中文音素表 |
| R19 | python 解析 `models/manifest.json` 与 `build/models/manifest.json` | 两份各 **6 枚** entry；`kws-zipformer-wenetspeech-3.3M-2024-01-01` purpose=`kws`，archive 内 `keywords.txt` 带 `sha256=46de8c68…`、`size_bytes=286` | C29 清单**逐文件**钉哈希（含 keywords.txt）；`internal/models/Purposes` 含 `"kws"`（`manifest.go:94`） |
| R20 | sherpa Go 绑定现读：`grep -n 'func NewKeywordStreamWithKeywords'`+`type KeywordSpotterResult struct`（module cache v1.13.8） | `NewKeywordStreamWithKeywords(spotter, keywords string)` 在位；`KeywordSpotterConfig` 有 `KeywordsFile` **与 `KeywordsBuf`/`KeywordsBufSize`**（`sherpa_onnx.go:2408/2411` 真往 C 传）；**`KeywordSpotterResult` 只有 `Keyword string` 一个字段，没有置信度** | 自定义词可走内存缓冲、不必改已验签模型目录；但**逐词置信度阈值拿不到**（见 §5） |
| R21 | `grep -rn 'func encodeFrame' --include=*.go .`（排 `.scratch`） | **1 命中**：`internal/audio/wavinjector.go:135`，小写未导出；全仓无 `PCM16ToFloat`/`decodeFrame` 之类的反向工具 | 帧出包是 `[]byte` 小端 int16，包外没有现成解码器 |
| R22 | `sed -n '45,80p' internal/models/downloader.go` → 实读 `VerifyInstalled`（`downloader.go:228`） | `VerifyInstalled` 走 `m.VerifyDir(entry, dir)`，按清单逐文件比哈希 | 与 R19 合起来才成立：**把自定义唤醒词写进模型目录会当场验签失败** |

---

## §1 采集侧：`internal/audio` 与 `internal/speech` 各暴露什么

### 1.1 一句话结论

**`internal/audio` 是写完了但没接线，`internal/speech` 是连包都还是空的。**
两枚包在主模块里** importer 数为零**（R1、R2）——`internal/audio` 的 6 个文件、约 4.4  KB 的
`wasapimic_windows.go` + 14.7 KB 的 `wasapi_windows.go` 全套真实 WASAPI 代码，今天都不在任何
可执行路径上。这不是"缺一块"，是**缺一根线**。

### 1.2 `internal/audio` 暴露的面（逐枚 `文件:行`）

| 名字 | 位置 | 是什么 | 非测试调用者 |
|---|---|---|---|
| `AudioSource`（C8 接缝） | `internal/audio/audio.go:23-33` | `Start(ctx, chan<- []byte) error` + `Stop() error` | **0**（R1/R14） |
| `Stater` | `audio.go:36-39` | `Stats() Stats` | 0 |
| `Stats` | `audio.go:82-89` | 帧/字节收发与丢弃计数、`Reopens`、`LastError` | 0 |
| 帧格式常量 | `audio.go:43-61` | `TargetRate=16000`、`FrameSamples=512`、`FrameBytes=1024`、`FrameDuration=32ms`、`BoundedWindow=192ms`、`BoundedFrameCapacity()`→6 | 0 |
| `NewBoundedFrames()` | `audio.go:65` | 建那条 ≤200ms 的有界通道（D38d） | **0** |
| `Path` / `PathT` / `PathC` | `audio.go:70-78` | D47 两档双工路径 | 0 |
| `meter.push` | `audio.go:115-145` | 非阻塞投递；满了**丢弃＋计数＋限速日志** | 包内私有，被两枚真源用 |
| `SpawnCapture` | `audio.go:180-182` | 占 `audio-capture` 常驻名的派生入口 | 0 |
| `WASAPIMicrophone` | `wasapimic_windows.go:35`（构造 `:52`） | **真麦克风**：共享模式 WASAPI、`runtime.LockOSThread` 钉线程（`:150-151`）、热插拔 reopen-once（`:270`）、占用/权限错误映射 | **0**（R3） |
| `WavInjector` | `wavinjector.go:48` | C8 的测试注入背骨（读 RIFF wav） | 0（仅包内 `gate_test.go:16`、`wavinjector_test.go`） |
| `HalfDuplexGate` | `gate.go:84` | D16/D47 装饰器：`SetMuted`/`SetSpeaking`/`Path`，T 路播报期关麦 | 0 |
| `Resampler`/`ResampleLinear`/`MonoDownmix`/`FloatToPCM16` | `resample.go:33/92/102/119` | 设备原生率 → 16k、多声道 → 单声道、float32 → int16 | 0 |
| `DeviceWatcher`/`streamOpener`/`deviceStream` | `device.go:28/45/51` | 枚举与打开的可注入接缝 | 仅包内与包内测试 |
| `newMMDeviceWatcher` / `wasapiOpener.Open` | `mmdevice_windows.go` / `wasapi_windows.go:210` | **真的在打开设备**：COM vtable 直调（`comCall` `:129`）、`IAudioClient::Initialize`（`init` `:256`）、`GetService(IAudioCaptureClient)`、`WaitEvent`（`:341`）、`Drain`（`:361`）、`convertPacket`（`:394`） | 0 |

**"有没有真的打开设备的代码"：有，而且是全套共享模式 WASAPI + MMAudio 端点枚举，不是桩。**
判定依据（现读）：`wasapi_windows.go:210 func (wasapiOpener) Open(dev DeviceDescriptor)`、
`:256 func (s *wasapiStream) init() error`、`:361 func (s *wasapiStream) Drain() ([]int16, error)`、
`mmdevice_windows.go`（`IMMDeviceEnumerator` + `OnDefaultChanged` 回调）。
没有 `waveIn*`、没有 PortAudio —— 本仓走的是 WASAPI COM，符合 `internal/audio/doc.go:2` 的自述。

### 1.3 `internal/speech` 暴露的面

**只有 `internal/speech/doc.go`（21 行）一个文件。**目录实测：

```
internal/speech: doc.go   (1066 字节)
```

`doc.go:1-21` 冻结了包的职责边界，并且**自己承认引擎未实现**：
`doc.go:19` 逐字 `DEFERRED(engines): implemented by ticket 15 (ASR + CER harness), ticket 26`
`/ (TTS), ticket 41 (KWS). This ticket only freezes the package boundary.`

⇒ 结论口径要说清：**没有 `AsrEngine`/`TtsEngine`/`WakeWordEngine`（C9）任何 Go 类型落地**。
现跑两把尺：`grep -rn 'type .*Engine' --include=*.go internal cmd tools` → **零命中（rc=1）**；
`grep -rn 'WakeWordEngine\|AsrEngine\|TtsEngine'` → **仅 2 命中，都在 `internal/speech/doc.go:1-2`
的包注释里**（只有名字被提到，没有一处声明）。
所以"KWS 引擎"这一侧不是"写了没接"，是**一块都没写**。

### 1.4 从"采到一帧 PCM"到"算出一个 RMS 电平"缺哪几环

先把"已经有了"的那一段写实：`wasapimic_windows.go:222-233` 这七行就是终点前的最后一站——
`stream.Drain()` 出设备原生率单声道 int16 → `res.Process()` 出 16k 样本 → 攒满 512 →
`m.meter.push(buf, encodeFrame(...))` 出**一帧 1024 字节的 `[]byte`**。
球那边收的地方也写实了：`internal/ball/liquid_windows.go:42 SetAudioLevel(level float32)`，
其注释 `:35` 逐字 `RMS of the last 512 samples, normalised 0..1`。

**中间缺五环，逐环给判定：**

| 环 | 现状 | 缺什么 | 补的量级 |
|---|---|---|---|
| ① 通道的消费者 | `NewBoundedFrames()`（`audio.go:65`）**零调用者**；`Start(ctx, buf)` 的 `buf` 今天没有任何人建、没有人读 | 一个 `for frame := range buf` 的循环体＋它的 owner | 缺一整块（但要的循环很小） |
| ② `[]byte` → int16 解码 | `encodeFrame` 在 `wavinjector.go:135`，**未导出**（R21）；包外无 `PCM16ToFloat`/`decodeFrame` | 要么 `audio` 导出一个解码/取样本 helper，要么消费者自己重抄 8 行小端解码 | 补一行级（若导出）／约 10 行（若在包外重抄） |
| ③ RMS 本身 | **R12 全主模块 `Sqrt` 零命中** | 一个逐帧 RMS（512 样本平方和均值开方） | 缺一整块，约 15 行 |
| ④ 归一到 0..1 的标度 | 球的文档写了 `normalised 0..1`，但**全仓没有任何一处定义这个除数**；球侧唯一的相关常数 `SilenceLevelGate = 0.06`（`liquid.go:30`）是**判静默的门限、不是标度** | 一条口径决定：`rms/32767` 线性？还是要增益/压缩？ | **这是判定项，不是代码量**：定不下来就写不出来（见 §5 Q1） |
| ⑤ 投递到球的那根线程 | `SetAudioLevel` 自己 `b.sta.PostTask(...)`（`liquid_windows.go:44`），跨线程安全，不需要新线程 | 只需要消费者在读到帧时调用 | 补一行级 |

**还有第六环，是形状问题不是代码量：电平送进球也不会动。**
`liquidDriven(s)`（`internal/ball/liquid.go:60-68`）只对 `Listening/Thinking/Speaking/
Conversation/Warm/Acting` 返回 true；`Armed`、`Sleeping`、`Muted` 不在表内，
`applyLevel` 对这些状态**记录标量后当场丢弃**（`liquid_windows.go:52-56`，注释逐字
`a hot mic can never wake Sleeping`）。而唤醒词待命态恰好就是 `Armed`。
⇒ 就算 ①—⑤ 全接上，**待唤醒态的波纹仍然是死的**，"波纹控件跟着环境声音轻微起伏"这一条
需要球侧状态白名单扩容（那是 **SPEC-08 冻结视觉表的射程**，见 §3 与 §5）。

**隐私口径（按事实写，不升级定性）**：
- 现象在哪出现：`internal/audio` 今天**没有任何生产 import 者**（R1），进程里没有人打麦克风；
  `cmd/wisp` 无参数常驻路径（`resident_windows.go:24`）不 import `internal/audio`、不 import
  `internal/ball`（R11），只做 `proc.Boot` + 日志 sink + 空事件循环。
- 有没有本机被入侵的证据：**没有，本腿一条都没找到**。
- 最坏后果是什么形状：这条链今天最坏的形状是"接线之后按默认配置仍然不开麦"
  （`internal/config/schema.go:197` `Enabled bool ... default:"false"`，且 R13 显示**今天没人读它**）。
  所以现状是"想开都开不了"，不是"悄悄在听"。
  `PLAN.md:460`（默认不开麦）与 `:462`（音频缓冲永不落盘、不写日志）里，
  `:462` 那一半在代码里是**真做到的**：`audio.go:18-20` 与 `meter.push`（`audio.go:136-141`）
  的日志只带计数与设备名，不带 payload；本腿未见任何写音频文件的路径。

---

## §2 KWS 侧：绑定怎么引的、谁建过 spotter、两档预算的落点、opt-in 与指示

（待写。）

---

## §3 两条链的落点清单

### 3.1 采集 → 电平 → 球

（待写，一张表。）

### 3.2 采集 → KWS → 唤起会话

（待写，一张表。）

---

## §4 半成品名册

（待写：每枚 `文件:行` + 「补一行就通」/「缺一整块」。）

---

## §5 我答不上来的

（待写，本节不许为空。）
