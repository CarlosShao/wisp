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

## §2 KWS 侧：绑定怎么引的、谁建过 spotter、两档预算、opt-in 与指示

### 2.1 sherpa-onnx 的 Go 绑定在本仓是怎么引的（cgo 那层在哪）

三层，逐层给出处：

1. **Go 绑定＝上游模块，不在本仓。** `go.mod` 直依赖第 1 行就是
   `github.com/k2-fsa/sherpa-onnx-go v1.13.8`（R8），`sherpa-onnx-go-{linux,macos,windows}`
   三枚是 indirect。`sherpa-onnx-go@v1.13.8/sherpa_onnx/` 里只有三个转发文件
   （`sherpa_onnx_windows.go` 等，`type KeywordSpotter = sherpa.KeywordSpotter` 这种别名）。
2. **cgo 那层在上游的平台模块里，不在本仓。** 实现在
   `~/go/pkg/mod/github.com/k2-fsa/sherpa-onnx-go-windows@v1.13.8/sherpa_onnx.go`
   （`NewKeywordSpotter` 在 `:2396` 段、`GetResult` 在 `:2471`），链接开关在
   同目录 `build_windows_amd64.go:5`，逐字：
   `#cgo LDFLAGS: -L ${SRCDIR}/lib/x86_64-pc-windows-gnu -lsherpa-onnx-c-api -lonnxruntime`。
   ⇒ **要点：链接期用的是 module cache 里那三枚 DLL，本仓 `third_party/` 那份不参与链接。**
3. **本仓 `third_party/sherpa-onnx/` 只有运行期 DLL。** 实测（R9）：
   `onnxruntime.dll`(17,799,168B) + `sherpa-onnx-c-api.dll`(4,605,952B) +
   `sherpa-onnx-cxx-api.dll`(259,584B) + `.cache-manifest.json`，**零 `.go`、零 `.h`**。
   它的作用是"加载期 PATH 上要有 DLL"，`cmd/wisp` 的一堆测试头注释就是这个
   （例：`cmd/wisp/always_write_no_clobber_226_test.go:24-25` 逐字
   `this package's test binary links sherpa-onnx and dies at load (0xc0000135) unless
   third_party/sherpa-onnx is on PATH`）。
4. **版本钉与自校验在位**：`deps.toml` 里 `sherpa-onnx` 与 `go-binding-sherpa-onnx` 两段
   被 `cmd/wisp/doctor.go:166-177` 解析，`doctor.go:60-68` 拿运行期
   `sherpa.GetVersion()` 与 `buildinfo.SherpaOnnxVersion` 对撞，不等就 FAIL。
   ⇒ **这是全仓唯一一处生产代码真的调 sherpa 的地方**（R5：`doctor.go:60/81`、`main.go:165`，
   全是"报版本号"，没有一处推理）。

**模型侧的分发已经就绪**（这块别漏）：C29 清单 `models/manifest.json` 有 6 枚条目、
其中 `kws-zipformer-wenetspeech-3.3M-2024-01-01` purpose=`kws`，带 archive sha256
`b2f7c896…`、`size_bytes=32654866`，并在 archive 内**逐文件**钉了
`encoder/decoder/joiner.int8.onnx` + `tokens.txt` + **`keywords.txt`（sha256 `46de8c68…`，286 字节）**
（R19，两份清单 `models/manifest.json` 与 `build/models/manifest.json` 各 6 枚）。
`internal/models/manifest.go:94` 的合法 purpose 表含 `"kws"`。

### 2.2 全仓有没有任何一处创建 keyword-spotter 的代码

**有，两处；两处都不在生产、也不在测试。** 按"生产 / 测试 / 其它"三分：

| 分类 | 命中 | 判定 |
|---|---|---|
| **生产码**（`internal/**`、`cmd/**`，排除 `_test.go`） | **0 处** | `NewKeywordSpotter` 零命中；`internal/speech` 连类型都没声明（§1.3） |
| **测试码**（`_test.go`） | **0 处** | 本腿 `grep -rn 'KeywordSpotter' --include='*_test.go' internal cmd tools` 零命中（口径限制：`.scratch` 下别人的探针件不算，也未被本腿转述） |
| **spike 码**（独立模块 `scripts/spike`，R7） | **2 处，且是完整能跑的真代码** | 见下 |

spike 两处的精确位置：

- `scripts/spike/model-residency/main.go:114 func kwsConfig(modelsDir string) (*sherpa.KeywordSpotter, string)`，
  `:116` 真建 spotter（`KeywordsFile`、`KeywordsThreshold: 0.25`、`KeywordsScore: 2.0`、
  `MaxActivePaths: 4`、`NumThreads: 1`、`Provider: "cpu"`），`:142 openKws()` 里 `:143 NewKeywordStream`
  → `:145-148` 1 秒静音 warmup（`AcceptWaveform(16000, chunk)` + `IsReady`/`Decode` 循环）
  → `:151-152` `DeleteOnlineStream`/`DeleteKeywordSpotter`。
- `scripts/spike/speech-baseline/main.go:124-167`：同一套配置，`:147 NewKeywordStream`，
  `:149-155` 喂 **3 秒** 512 样本块并解码（注释自称 `realistic Armed state`），
  `:164-167` 释放；`:131` 行注释逐字 `NumThreads: 1, // D32: intra_op_num_threads pinned to 1`。

⇒ 这两处的价值要如实说：**它们是本仓唯一一份"怎么把 KWS 跑起来"的可抄答案**，
包括模型文件名拼法、四个阈值常数、chunk 尺寸（512 = `audio.FrameSamples`）、
以及 `IsReady`/`Decode` 的取结果节奏。但它们**不读 config、不产事件、不进名册**，
而且**在另一个 module 里**——`internal/speech` 引不到它们，只能照抄形状。

### 2.3 `PLAN.md:527/528` 两档预算在代码里对应哪些常量

对应关系是**精确到数字**的，位置全部在 `internal/observe/thresholds.go`：

| PLAN.md 那一档 | 原文（本腿 `sed` 现读） | 代码常量 | 实测换算 |
|---|---|---|---|
| `:527` 空闲态（KWS 关） | `≤ 25MB` / `≤ 0.5%（1min 均值）` / 麦克风关闭·零网络长连接·零周期性磁盘写入·悬浮球静态不做动画 | `thresholds.go:19 memCapSleeping int64 = 25 << 20`；`:26 cpuLimitSleeping = 0.5`；`:36 goroutineLimitSleeping = 6`；`:191-205 sleepingConstraints`（`disk_write_ops == 0`、`tcp_connections == 0`，均 `Gate: true`） | `25 << 20 = 26,214,400 字节 = 25 MiB`（口径：`mb()` 在 `:208` 用 `1<<20`，即 MiB 而非 MB） |
| `:528` 待唤醒态（KWS 开） | `≤ 90MB` / `≤ 2%` / 音频缓冲不落盘（D16） | `thresholds.go:20 memCapArmed int64 = 90 << 20`；`:27 cpuLimitArmed = 2.0`；`:38 goroutineLimitArmed = 7 // + kws-infer` | `90 << 20 = 94,371,840 字节`；`7 = 常驻 6（goroutine.go:41 ResidentBaseline）+ 1` |

判定：**数字全在、闸门也真的会红，但"待唤醒"那一档今天测不到、也进不去。** 现跑尺（R：
`grep -rn 'SLOArmed' --include=*.go cmd tools internal | grep -v '^internal/observe/'` → **零命中，rc=1**）：

- `SLOArmed` 只在 `internal/observe` 自己体内出现（`sampler.go:43` 声明、`sampler.go:51` 进
  `SLOStates` 表、`thresholds.go:49/66/171` 查表）；
- 跑 SLO 的 `cmd/wisp/slo_windows.go`：`-state` 一枚旗子**是收 Armed 这个名的**
  （`:149` usage 逐字 `-state <Sleeping|Armed|Warm|Conversation|PanelOpen|WorkPeak>`，
  校验走 `SLOState(*state).Valid()`），**但该文件零 import ball、零 SetState**，
  进程侧也没有任何事件能进 Armed ⇒ **旗子只是标签，采到的窗口不是真 Armed**（详见 §4 F1）。
- 而 `Facts.KwsLoaded`/`EvKwsEnabled` 生产零生产者（R15）⇒ 进程今天**无法进入 `Armed` 态**。

⚠ 另外这枚文件路径要说清，否则会找错：任务书与部分文档把它写成
`internal/ball/thresholds.go`——**该路径不存在**（本腿实测 `sed: can't read
internal/ball/thresholds.go: No such file or directory`）。真身在
**`internal/observe/thresholds.go`**。本腿两枚都没动（红线一致：一字节都不许改）。

`:514` 那行"看门狗阈值必须按态查表：在 `Armed` 态用 25MB 阈值会触发卸载 KWS，直接杀掉唤醒词"
的落点也就是 `thresholds.go:47 stateMemCap`/`:64 stateCPULimit` 这两个按 `SLOState` 查表的函数
——**表已经画对了**；`internal/watchdog/` 本身**只有 `doc.go` 一个文件**
（`doc.go:16-17` 逐字 `DEFERRED(watchdog loop/thresholds): implemented by ticket 42`），
所以"按态查表去执行动作"那一半没人做。

### 2.4 `PLAN.md:460`「默认不开麦（opt-in）」与 `:462`「KWS 激活时常亮可见指示」有没有落点

任务书要求分"配置项 / 状态位 / 界面指示"三处各判。逐处：

**(a) 配置项——三件套里 `:460` 有完整落点，`internal/observe`/`cmd` 无人执行。**

| 落点 | 位置 | 状态 |
|---|---|---|
| 默认 false 的开关 | `internal/config/schema.go:197` `Enabled bool` + tag `toml:"enabled" default:"false"`（**这枚是 `wake_word.enabled`**） | ✅ 在 |
| ⚠ 同一棵树上的另一枚开关 | `schema.go:245` `VoiceSection.Enabled` + tag `toml:"enabled" default:"true"`（`voice.enabled` 默认 **true**） | 在，但**方向相反**：`:460` 那句"默认不开麦"精确对应的是 `wake_word.enabled=false`（常开听音），**不是** `voice.enabled`。接线腿必须自己裁定"哪一枚门控哪一种开麦"，否则默认配置下会读到一个 true（见 §5 D-2） |
| 挂在 `[voice]` 树下 | `schema.go:247` 字段 `WakeWord`，tag `toml:"wake_word"`；类型声明在 `schema.go:196 type WakeWord struct` | ✅ 在 |
| 关键词与逐词阈值 | `schema.go:199 Keywords []string`、`:201 Thresholds []float64`、`:203 VetoWords`（默认 `取消,停下,别`） | ✅ 在 |
| 热载分级 | `internal/config/manager.go:370-375`（`WakeWord.Enabled`/`Keywords` 变更算 **reload-tier**）、`:385-386`/`:398-399`（写回）、`:405-407`（`Thresholds`/`VetoWords` 算 hot-tier） | ✅ 在，且判据齐全 |
| **读它并决定开不开麦的消费者** | **无**（R13：`internal/config/` 之外 `.Voice.` 零命中） | ❌ 缺 |

⇒ 精确说法：`:460` 的"默认不开麦"今天是**被动成立**的——不是因为有代码守住它，
而是因为**根本没人能开麦**。这一点对接线腿很要命：接线之后这句会从"碰巧成立"变成
"必须被实现"，落点就是那个缺失的消费者。

**(b) 状态位——两套都建好了，两套都零生产者。**

- `internal/statemachine`：`events.go:18 EvKwsEnabled`（D43 #5）、`events.go:20 EvWakeWord`（#7）、
  `events.go:70 Facts.KwsLoaded`、`table.go:55`（`Sleeping --EvKwsEnabled--> Armed`，副作用
  `kws.load`）、`table.go:61`（`Armed --EvWakeWord--> Listening`，副作用
  `session.scope-create` + `speech.load-vad-asr` + `kws.pause`）、`table.go:76/80` 与
  `:237/238` 两处 `kws-loaded`/`kws-not-loaded` 守卫（#10、#32）、`table.go:66` `kws.stop-inference`、
  `table.go:71` `kws.keep-alive-alert`（#9 Armed 自环，注释逐字 `D32③: KWS must NOT be unloaded`）。
  ⇒ **转移表这半是完整的**（D43 冻结的就是这张）。缺的全在"谁发事件"：R15 实测
  `EvKwsEnabled`/`EvWakeWord` 生产零命中。
- `internal/agent/approval`：`approval.go:57 ChannelKWS = "kws"`、`:61` 进 `allChannels`、
  `:92` 文案 `说取消词`、`:99 DEFERRED(kws-veto, B1)`、`:105` 的 `case ChannelKWS:`、
  `:161` 逐字 `the panel (ticket 37) and KWS (ticket 41) are not`（`DefaultChannels()` 只放
  Ball+Esc，`:162`）、`:176-178 SetLoaded`（注释逐字 `the KWS loader calls this when its model
  is actually resident`）、`gate.go:137` 暴露 `Channels()` 给"the KWS loader (ticket 41)"。
  ⇒ 唯一生产 `SetLoaded` 调用在 `cmd/wisp/approval_reply.go:416`，它翻的是
  `run.go:618` 传进来的 `s.replyVeto`，**不会是 KWS**；`run.go:456-457` 逐字说明控制台腿
  `NewChannels() stays empty and runSpec.replyVeto stays unset`。
  **也就是说：今天没有任何一条路径会把 KWS 通道谎报成可用**，`window_test.go:92-100` 还把
  这条正反两向钉着（`:97` 逐字 `KWS 未加载却报为可用（这是在假装支持）`，`:100` 要求文案含
  「语音取消不可用」）。这是好事，接线时别把它改松。

**(c) 界面指示——只有"半透明静态球"这一形，且它同时是冻结视觉表规定的形状。**

| 候选落点 | 位置 | 是不是 `:462` 要的"常亮可见指示" |
|---|---|---|
| 冻结表 `Armed` 行 | `internal/ball/statevisual.go:158 case statemachine.StateArmed: v.Opacity = 0.6`（对应 `docs/specs/SPEC-08-ui-ball-panel.md:57` 逐字 `Armed \| opacity 0.6，直径 44px，静态`） | **是冻结表对 Armed 的全部规定**；它是"待命态可辨"，不是"麦克风正开"的交代 |
| 原型视觉 `Armed` 行 | `statevisual.go:288 v.Opacity = 0.92` | 同上，只差透明度 |
| 图标 | `statevisual.go:11-21` 的 `IconKind` 枚举里**没有 Armed 专用图标**（`IconMic` 归 Conversation、`IconAudioLines` 归 Listening）；`Armed` 分支不设 `v.Icon`，取零值 `IconNone` | ❌ 缺 |
| 动画 | `internal/ball/anim.go:62` 把 `Armed` 落进 `default: return AnimPolicy{AnimNone, 0}`，注释逐字 `Sleeping, Armed, Muted, ...: static` | ❌ 明确静态 |
| "麦开"的唯一交代在别处 | `SPEC-08:71` 逐字 `Conversation \| **danger 常亮环，不得渐隐、不得弱于 Confirming**（麦克风开着的唯一交代）` | 那是**会话态**，不是 KWS 待命态 |
| 一键静音快捷键 | `internal/ball/hotkey_reload.go`、`internal/ball` 的 `EvMuteKey` 链路（D43 #8/#10） | ✅ `:462` 的后半有落点；且 `observe/goroutine.go:44` 的 `hotkey-listener` 常驻名今天**零 Spawn**（R16） |

⇒ `:462` 前半句（常亮可见指示）**今天没有专属落点**：能拿的只有 `Armed` 的透明度差；
`SPEC-04:47` 与 `SPEC-04 §6:80` 两处把这句要求各写了一遍，代码侧只有 opacity。
另外这枚指示在 owner 的"平时不挂界面"新形态里**没有载体**（编排者已在 A461 把这句列为
"新形态第一枚判据"），本腿不重复裁它，只把代码事实摆出来。

**隐私三件套（`:460`/`:462`/`:463`）按事实三行定性：**
- 现象在哪出现：`:463`（音频缓冲永不落盘/不写日志）在 `internal/audio` 内部**是做着的**——
  `audio.go:18-20` 明文规则、`audio.go:136-141` 的丢弃日志只带计数/窗口/帧长、
  `audio.go:152` 的失败日志只带 source 名与 error 串；本腿在 `internal/` 全树未见任何写音频
  文件的路径。`:460`（默认不开麦）今天的成立方式是"没人能开麦"（R1/R13）。
- 有没有本机被入侵的证据：**没有。本腿零证据。** 本腿也没跑任何进程。
- 最坏后果是什么形状：现状最坏＝"想接还接不上"（缺引擎、缺消费者、缺电平）。
  接线之后最坏＝按默认配置仍不开麦（`schema.go:197` default false），
  要开必须用户显式改 `[voice] wake_word.enabled` 并触发 reload-tier（`manager.go:370`）。
  ⇒ **不写成"窃听风险"这类没有证据的定性；本腿判"未定性"，事实如上。**

---

## §3 两条链的落点清单

### 3.0 先摆"本仓已存在的纪律与禁令"（同类形状现读，逐条带出处）

这两条链不是"想怎么接就怎么接"，下面每一枚都会真实挡住某个写法。全部是本腿现读的仓内文本：

| 号 | 纪律（原文或精确转述） | 出处 | 对这两条链的直接后果 |
|---|---|---|---|
| P1 | 「**零新增依赖边**」是被当作**验收理由**写下来的，不是风格偏好。票 223 票面逐字：`后者零新增依赖边：run.go 已同时持有 config 与 gate`；124 的裁决逐字：`⇒ 无环、无新依赖边（§7.4 双 GOOS vet 各 rc=0 为证）`，同段还把"提成生产包会新增一条 `memory/config → proc` 的生产依赖边"当作**否决理由** | `.scratch/wisp/issues/223-checkandreload-has-zero-production-callers-hot-reload-never-runs.md:15`；`docs/evidence/s1/124-ac2b-4-conversion.md:72` | 装配根（`cmd/wisp/run.go:328 assembleRuntime`）已经同时持有 config 与 gate；**两条链的接法优先"装配根注入函数/句柄"，不要为了少写一行而新开包间直连** |
| P2 | 状态机的副作用**设计上就是由消费者注入的**：逐字 `executes them itself - consumers (session, speech, panel, tickets 08+) wire a Sink. The zero sink is an explicit no-op`；`Sink` 的约束逐字 `Must not block and must not re-enter Dispatch synchronously` | `internal/statemachine/machine.go:24-25`、`:31-34`、`:64-68`（nil → `func(Effect){}`） | 这就是"panel/装配根注入、不新开依赖边"的**同一形状的仓内实现**：KWS 的 `kws.load`/`kws.pause`、语音的 `speech.load-vad-asr`、采集的 `audio.stop-capture` 全等着一个 Sink，**今天生产两处建机器都没给 Sink**（§3.3） |
| P3 | `audio-capture`（pinned）线程：逐字 `runtime.LockOSThread()`；**该线程不得跑任何其他 Go 代码** | `docs/specs/SPEC-01-architecture.md:113`；实现 `internal/audio/wasapimic_windows.go:150-151` | RMS/电平/投递**能不能算在采集线程上**是一枚要裁的问题（见 §5 D-3）；现状实现里那条线程只做 open/wait/drain/resample/push |
| P4 | goroutine 花名册：常驻**上限 6**、按需 `kws-infer`**仅 Armed 态**、`总数上限 = 常驻 6 + 每任务 3；超出即泄漏征兆` | `SPEC-01:119-124`；代码 `internal/observe/goroutine.go:41 ResidentBaseline = 6`、`:44` 六枚名、`:47 OnDemandNames = []string{"kws-infer"}` | **给电平消费者新起一条常驻 goroutine ＝ 动 D38b 名册**（契约面）；名册里已经**给 KWS 留了槽**（`kws-infer`），R17 实测该槽今天零 Spawn |
| P5 | 禁止裸 `go func()`：必须有名字/owner/退出条件/recover 边界 | `SPEC-01:125`；AGENTS.md §1.2（`tools/d22scan` 自动扫） | 两条链新增的任何 goroutine 一律走 `observe.Registry.Spawn` |
| P6 | 接缝**全部复用冻结契约，不新开接缝**：逐字 `本项目接缝已由冻结契约天然划定，全部复用，不新开接缝` | `docs/specs/SPEC-10-testing-acceptance.md:8` | 注入面只有 C8(`AudioSource`) / C5(`LlmProvider`) / C17(`PanelBridge`) / CLI `wisp run`；**别为电平或 KWS 造新接缝** |
| P7 | 依赖白名单（新增须人批准）：`sherpa-onnx` Go 绑定（cgo）**已在名单内**；`webrtc-audio-processing` 也随 D47/P15 批准加入 | `SPEC-01:96-100` | 接 KWS **不需要**人批新依赖（模块已在 `go.mod:8`）；只有真要上 AEC 才碰那枚 |
| P8 | C25 污染面：逐字 `This is the only render-side audio input in the project: no samples and no transcript text cross into the ball` | `internal/ball/liquid_windows.go:38-40` | 电平链**只许一个 float32 过界**；把 PCM 帧或转写喂进球＝违规，别为了"波纹更准"破例 |
| P9 | D38(e) 关停 10 步：`2. 停热键与唤醒词监听`、`4. 停音频采集线程 → 关设备 ← 必须在释放 ASR 之前（否则 ASR 拿到半截 buffer）`、`5. 释放 ASR/TTS/KWS session（DisposalScope 逆序销毁）` | `SPEC-01:139-143`；实现 `cmd/wisp/resident_windows.go:66-75`（`rt.Shutdown`） | 两条链各要在关停表里挂上；KWS loader 的释放点在第 5 步，采集消费者停在第 4 步 |
| P10 | 球的动效白名单：liquid 的过渡定时器 `may only run in a state the frozen SPEC-08 animation whitelist already grants a timer to (see transitionDriven)`；`Sleeping appears in neither list` | `internal/ball/liquid_windows.go:16-19`、`:24-27`；`internal/ball/anim.go:40-70` | 想让 `Armed` 也有动效＝同时动**冻结视觉表**（SPEC-08 §2.1 `:57`）与**定时器白名单**两枚，不是接线量级 |
| P11 | 默认视觉开关 `prototypeVisuals` **仍是 OFF**，且注释自陈"翻默认是票 68 AC#2"、`cmd/balldebug is today the only caller that turns this on, so a library consumer (the panel, a future wisp GUI host) still gets the frozen look` | `internal/ball/statevisual.go:88-102`（`var prototypeVisuals bool` `:102`、`EnablePrototypeVisuals` `:106`）；实测调用者只有 `cmd/balldebug/main.go:122` | **整条电平链在今天默认视觉下是空转的**：`applyLevel` 第一行 `if !prototypeVisuals { return }`（`liquid_windows.go:52`） |

⚠ **上面这两枚（"panel 为真相源""不新开 `tools → panel` 这类依赖边"）本腿找到了逐字原文，但不在 `docs/specs/**`、也不在工单池 README，而在台账**：
`docs/reports/pending-and-issues.md:8786-8792`（票 197 那一段）。逐字摘录：

> `SubagentStreamKeyPrefix`／`SubagentStreamKey` **在两包里各定义了一份**（`internal/tools/subagent_197.go:43-48,91` 与 `internal/panel/pump.go:375-379`），
> 我现量两条 grep 确认**双向都不 import**（`grep -rn "wisp/internal/panel" internal/tools/`＝空、反向亦空）⇒ **装配根 `cmd/wisp` 是唯一的接缝**。
> **裁定＝`internal/panel` 那一枚是真相源** …… `internal/tools` 那份**改成由装配根注入**……
> **不新开 tools→panel 的依赖边**（那会把"视图层"变成"工具层"的下家，方向装反）。

⇒ 这条对本表的直接含义（本腿的推广，不是台账原话）：**"两包都要用的东西"一律由 `cmd/wisp` 装配根注入，
被注入的那一枚定为真相源**；照此形——电平消费者要球 ⇒ 装配根把"投递函数"注给消费者，
`internal/audio` 不 import `internal/ball`；KWS 要状态机 ⇒ 装配根把 machine 注给 loader，
`internal/speech` 不 import `internal/statemachine` 之外再顺手够到 UI。

扫法（负向尺，本腿真跑）：`grep -rn -iE 'panel.{0,40}真相源|真相源.{0,40}panel' docs/specs/ .scratch/wisp/issues/README.md docs/evidence/s1/`
→ **零命中**；`grep -rn '依赖边' docs/ .scratch/wisp/issues/README.md` → **3 命中**
（`docs/evidence/s1/124-ac2b-1-conversion.md:196`、`docs/evidence/s1/124-ac2b-4-conversion.md:72`、
`docs/reports/pending-and-issues.md:8791`）。⇒ **结论：禁令的权威文本在台账与裁决表里，不在 spec 里**，
引用它必须点名这两枚出处，不能写成"spec 规定"。

### 3.1 表一：采集 → 电平 → 球

| # | 落点 | 现状位置 | 要动什么 | 新依赖边？ | 量级 |
|---|---|---|---|---|---|
| 1 | **常驻 GUI 腿得先存在** | `cmd/wisp/resident_windows.go:81` 逐字 `empty event loop running; the floating ball arrives in ticket 07`、`:83 rt.RunEventLoop()`；循环本体 `internal/proc/boot_windows.go:120-131` | 在这条腿上建球＋建消费者；`internal/ball` 目前**只有 `cmd/balldebug/main.go:27` 一枚 import 者**（R11） | **新**：`cmd/wisp → internal/ball` | **缺一整块**（票 07 的地界，比电平链本身还大） |
| 2 | 装配根 | `cmd/wisp/run.go:328 assembleRuntime`（run 腿已有）；GUI 腿**没有对应的 assemble** | GUI 腿需要一处"同时握 ball / audio / config / machine"的装配点 | — | 缺一整块（可复用 `assembleRuntime` 的形状） |
| 3 | 建那条 ≤192ms 有界通道 | `internal/audio/audio.go:65 NewBoundedFrames()` | **补一行调用**（今天零调用者） | 无 | 补一行就通 |
| 4 | 起真麦克风 | `internal/audio/wasapimic_windows.go:52 NewWASAPIMicrophone()`、`:69 Start(ctx, buf)` | **补一行**；opt-in 与静音要先过门（下一行） | **新**：`cmd/wisp → internal/audio` | 补一行就通 |
| 5 | 静音/半双工门 | `internal/audio/gate.go:84 NewHalfDuplexGate(inner, path, opts...)`、`:130 SetSpeaking`、`:168 SetMuted`、`:209 Open` | **补一行**把 mic 包进门里，门再给消费者 | 无（同包） | 补一行就通 |
| 6 | 帧解码 `[]byte → int16` | `encodeFrame` 未导出：`internal/audio/wavinjector.go:135`；包外零反向工具（R21） | 在 `internal/audio` 导出一个取样本/直接给电平的口（推荐），避免包外重抄小端 | 无 | 约 5–10 行 |
| 7 | **RMS 本体** | **全主模块 `Sqrt` 零命中（R12）** | 新写逐帧 RMS（512 样本） | 无 | **缺一整块**（约 15 行＋用例；`FrameSamples` 已有） |
| 8 | **0..1 归一标度** | 球侧文档 `liquid_windows.go:35` 逐字 `RMS of the last 512 samples, normalised 0..1`；**除数全仓无定义**；球唯一相关常数 `liquid.go:30 SilenceLevelGate = 0.06` 是门限不是标度 | 需要一条口径定案（§5 D-1）后落成常量 | 无 | **判定项**，不是代码量 |
| 9 | 投进球 | `internal/ball/liquid_windows.go:42 SetAudioLevel(level float32)`（内部 `:44 b.sta.PostTask`） | **补一行**调用 | 消费者侧需握 `*ball.Ball`：由装配根注入，**不要让 `internal/audio` import ball**（撞 P1/P8） | 补一行就通 |
| 10 | 消费者的那条 goroutine | 名册无对应槽（P4：常驻 6 已排满、按需只有 `kws-infer`） | 要么复用采集线程（撞 P3，需裁），要么动名册（D38 契约面） | 无 | **缺一整块＋可能碰契约** |
| 11 | 视觉开关 | `statevisual.go:102` 默认 false；`liquid_windows.go:52` 直接早退 | 装配根显式 `ball.EnablePrototypeVisuals(true)`，或等票 68 AC#2 翻默认 | 无 | 补一行就通（但**默认翻转不归写腿**） |
| 12 | 状态白名单（`Armed` 不动） | `liquid.go:60-68 liquidDriven` 六态表，**不含 Armed**；`anim.go:62` Armed 落 default 静态 | 若要求"待唤醒态波纹跟环境声动"＝扩白名单 | 无 | **缺一整块，且撞 P10（冻结视觉表＋定时器白名单）** |
| 13 | 关停挂载 | `SPEC-01:142` step 4；`cmd/wisp/resident_windows.go:66-75` | 消费者先停、再关采集 | 无 | 补一行级 |
| 14 | 测试注入面 | `internal/audio/wavinjector.go:48 NewWavInjector(path)` + `:40 WithInjectorPace`；现成 wav 素材 `third_party/model-fetch/x-kws/.../test_wavs/*.wav`（7 枚） | 用例走 C8 注入，不许 mock 代替真源 | 无 | 补一行就通 |

**表一结论（一句话）**：14 格里 **#3/#4/#5/#9/#11/#13/#14 是"补一行就通"，#6/#7 是小块新写，#8 是判定，#1/#2/#10/#12 各是一整块**；
其中 **#1（GUI 腿根本没有球）比 #7（没有 RMS）更靠前**——现在这条链不是"最后一公里没接"，
是**两端都还没有家**（采集包没人 import，球包只被调试载体 import）。

### 3.2 表二：采集 → KWS → 唤起会话

| # | 落点 | 现状位置 | 要动什么 | 新依赖边？ | 量级 |
|---|---|---|---|---|---|
| 1 | **引擎包本体** | `internal/speech/` **只有 `doc.go`（21 行）**；`grep 'type .*Engine'` 主模块零命中 | 声明 C9 `WakeWordEngine` 并实现（load/infer/release、C11 DisposalScope 内释放不卸载） | — | **缺一整块**（`doc.go:19 DEFERRED(engines) … ticket 41`） |
| 2 | sherpa 那层的引入 | 上游模块 `sherpa-onnx-go v1.13.8`（`go.mod:8`）；cgo 在上游 windows 子模块（`build_windows_amd64.go:5`）；白名单已含（P7） | 在 `internal/speech` 直接 import 绑定 | **新**：`internal/speech → github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx`（**不需人批**，P7） | 补一行就通 |
| 3 | spotter 构造的可抄答案 | `scripts/spike/model-residency/main.go:114-152`、`scripts/spike/speech-baseline/main.go:124-167`（**独立模块 R7，引不到**） | 把形状搬进 `internal/speech`：四个常数 `KeywordsThreshold 0.25 / KeywordsScore 2.0 / MaxActivePaths 4 / NumThreads 1`＋`IsReady/Decode` 节奏＋`Delete*` 释放 | 无 | 缺一整块，但**零设计** |
| 4 | 模型路径与验签 | `models/manifest.json` 里 `kws-zipformer-wenetspeech-3.3M-2024-01-01`（R19，archive 逐文件钉哈希、含 `keywords.txt`）；`internal/models/manifest.go:94 Purposes` 含 `kws`；`downloader.go:228 VerifyInstalled` → `VerifyDir` 逐文件比哈希（R22） | 走 `Manager.Ensure` 拿目录；⚠ **自定义词不许写进模型目录**（当场验签失败） | **新**：`internal/speech → internal/models`（装配根注入目录路径可避免） | 补一行就通（前提是 #3 存在） |
| 5 | **keywords 从哪来** | `config.WakeWord.Keywords []string`（`schema.go:199`，人话词表）vs sherpa 要 `KeywordsFile` **或** `KeywordsBuf`（R20）；出厂表 `keywords.txt` 8 行全是**拼音音素＋@中文**（R18） | 缺"文本 → 模型音素行"的那一步（PLAN.md:168 说"只需 keywords 文本文件"，但文件内容是音素化拼写） | 无 | **缺一整块**（且是"能不能拼出 hi wisp"的可行性问题，见 §5 D-4） |
| 6 | 帧格式转换 | 采集出 `[]byte` 小端 int16（`wavinjector.go:135`）；sherpa 要 `[]float32` + sampleRate（`AcceptWaveform`，绑定 `sherpa_onnx.go:310`）；本仓只有**反向**的 `FloatToPCM16`（`resample.go:119`） | 补 `PCM16 → float32`（约 6 行）＋ 复用帧解码 #6(表一) | 无 | 约 10 行 |
| 7 | 推理线程 | 名册槽 `kws-infer` 已声明：`observe/goroutine.go:32`（注释 `kws-infer, Armed state only`）、`:47 OnDemandNames`；`thresholds.go:38 goroutineLimitArmed = 7 // + kws-infer`；R17 实测**零 Spawn** | 用这个槽 Spawn，**不要新起名**（P4） | 无 | 补一行就通 |
| 8 | 半双工下的 KWS 关/停 | `table.go:62 kws.pause`、`:66 kws.stop-inference`、`:71 kws.keep-alive-alert`（#9 自环，注释逐字 `D32③: KWS must NOT be unloaded`）；`PLAN.md:453` 播报期关麦&关 KWS；`audio/gate.go:130 SetSpeaking` | 副作用要有人执行——**见 #9** | 无 | 缺一整块 |
| 9 | **副作用执行器（Sink）** | 机制在位：`machine.go:26-34 Effect/Sink`、`:162-164 for _, name := range effects { m.sink(...) }`；`table.go:36` 字段注释 `effect hook names, fired in order`；生产两处建机器**都没传 Sink**：`cmd/balldebug/main.go:188`、`cmd/wisp/models.go:303`（→ 走 `machine.go:64-68` 的显式 no-op） | 在装配根接一枚 Sink，把 `kws.load/pause/stop-inference`、`speech.load-vad-asr`、`audio.stop-capture`、`audio.discard-buffer` 分派到真子系统 | 无（Sink 由装配根注入，正合 P1/P2） | **缺一整块，但这是全表最"值"的一格**：机制、事件名、表行、测试全都在，只缺那一个接线者 |
| 10 | 命中 → 唤起会话 | `statemachine/events.go:20 EvWakeWord`、`table.go:61`（`Armed --EvWakeWord--> Listening`，副作用 `session.scope-create` + `speech.load-vad-asr` + `kws.pause`）；R15 **生产零生产者** | KWS 命中后 `Dispatch(EvWakeWord, &Facts{KwsLoaded:true})` | 由 #9 的同一处持有 machine ⇒ **不新开 speech→statemachine 之外的边** | 补一行就通（依赖 #1/#3/#7） |
| 11 | 进入 Armed（opt-in 生效） | `events.go:18 EvKwsEnabled`、`table.go:55`（副作用 `kws.load`）；`Facts.KwsLoaded`（`events.go:70`）；**`schema.go:197 wake_word.enabled default false`，R13 显示 `.Voice.` 在非 config 包零读者** | 需要一个读 opt-in → `Dispatch(EvKwsEnabled)` → 翻 `KwsLoaded` 的门控者；同时裁定与 `schema.go:245 voice.enabled default true` 谁门控哪种开麦（§5 D-2） | 无 | 缺一整块（约一小函数＋一枚定时器复用 reload tick） |
| 12 | 语音否决通道翻可用 | `approval/approval.go:57 ChannelKWS`、`:161` 逐字 `the panel (ticket 37) and KWS (ticket 41) are not`、`:176-178 SetLoaded`（注释逐字 `the KWS loader calls this when its model is actually resident`）、`gate.go:137` 暴露 `Channels()`；生产唯一 `SetLoaded` 在 `cmd/wisp/approval_reply.go:416`，翻的是 `run.go:618` 传入的 `s.replyVeto`，`run.go:456-457` 逐字 `replyVeto stays unset` | KWS loader 真 resident 后调 `SetLoaded(ChannelKWS,true)` | 无 | 补一行就通 |
| 13 | 否决词（veto words） | `schema.go:203 VetoWords` 默认 `取消,停下,别`；`approval.go:92` 文案 `说取消词`、`:99 DEFERRED(kws-veto, B1)`；`docs/specs/SPEC-04-voice-pipeline.md:70-71` 逐字要求 `Confirming 态例外：KWS 保持运行以识别否决词` | 需要第二组 keywords + 置信度判定；**绑定拿不到置信度**（R20：`KeywordSpotterResult` 只有 `Keyword`） | 无 | **缺一整块，且撞工具天花板**（§5 D-5） |
| 14 | 常亮指示（`:462`） | 见 §2.4(c)：只有 `statevisual.go:158/288` 的 opacity，`Armed` 无图标、`anim.go:62` 静态 | 要么按 SPEC-08 现有 `Armed` 行交付，要么新增指示元素＝**视觉契约面** | 无 | 视取舍：照现有＝补一行；要新元素＝改冻结表（人批） |
| 15 | SLO 那档今天进不去 | R：`SLOArmed` 在 `internal/observe/` 之外零命中；`cmd/wisp/slo_windows.go:149` 的 `-state` 旗子收 Armed 这枚名，但该文件零 import ball／零 SetState ⇒ 标签能贴、状态进不去 | 接线后要让 Armed 行**真的被测到**，否则 `≤90MB/≤2%` 永远是纸面 | 无 | 缺一整块 |

**表二结论（一句话）**：这一条比表一**更靠前**——表一缺的是"没人按开关"，表二缺的是"开关后面还没有电器"。
唯一能立刻兑现的高性价比格是 **#9（副作用 Sink 接线者）** 和 **#12（SetLoaded）**；
真正的重活集中在 **#1/#3（引擎本体）** 与 **#5（唤醒词文本→音素）**。

### 3.3 两表共同的第一枚前置：装配根今天没有 GUI 腿

现读：`cmd/wisp/run.go:328 assembleRuntime` 是 **run 腿**的装配根，
`cmd/wisp/resident_windows.go:81/83` 的 GUI 腿只到 `rt.RunEventLoop()`
（实现在 `internal/proc/boot_windows.go:120`，循环体只有激活事件与信号），
**没有 assemble、不建球、不建采集**。而"同时握着 ball + machine + config"的那段装配代码
今天**只存在于调试载体** `cmd/balldebug/main.go:188-247`。
⇒ 两条链无论先做哪一条，第一个动作都是**把 balldebug 里那段装配形状搬到生产 GUI 腿**
（那是票 07 的地界）。这件事本腿只报事实，不替它排期。

---

## §4 半成品名册

### 4.0 口径先钉住

"零调用者"这一列本腿是**真跑出来的**，不是推断。尺形：

```
grep -rn "\b<符号>\b" --include='*.go' internal cmd tools | grep -v '_test.go' | grep -v '^internal/audio/' | wc -l
```

即"包外、非测试"计数。`internal/audio` 的 20 枚导出名（含常量）**全部＝0**，
逐枚读数在 §4.1 表里合并成一行写（因为 20 枚都是同一个 0，逐行列 20 遍是噪声不是证据）；
其余每一枚单独跑。⛔ 本腿**没跑** `go vet`／编译器来找"未使用的导出符号"，
所以这份名册的完整性只能到"本腿想到的候选"为止，缺口如实写在 §5。

本仓这条形状**不是新发现，是第六次**：票 **101**（perm 存储无生产 import 者）、
**105**（C26 改写台账无生产读者）、**114**（composer 请求无生产调用者）、
**117**（安全告警在生产里没听众）、**223**（`CheckAndReload` 零生产调用者、热重载从未跑）
——票面标题就是这么写的（`.scratch/wisp/issues/101-perm-store-has-no-production-importer-done.md`
等一批）。票 11 那段还把这条立成了规律：
`A8、A11 两张外来框都是"测试里已证明、装配根从未调用"，票 12 规划时应把"把已存在的东西接上"当主活，`
`不要当集成收尾`（`.scratch/wisp/issues/11-llm-adapters-rest-done.md:67`）。
⇒ **语音这条链是这个模式目前最大的一枚实例**，前面的票都在"一函数一档"量级，这枚是"两个包"。

### 4.1 名册（逐枚带 `文件:行` 与判定）

| # | 名字 | 位置 | 包外非测试调用者（实测） | 判定 |
|---|---|---|---|---|
| A1 | `internal/audio` **整包的 20 枚导出名**：`AudioSource`/`Stater`/`Stats`/`NewBoundedFrames`/`BoundedFrameCapacity`/`SpawnCapture`/`Path`+`PathT`/`PathC`/`NewWASAPIMicrophone`/`NewWavInjector`/`WithInjectorPace`/`NewHalfDuplexGate`/`WithStartMuted`/`WithGateEvents`/`NewResampler`/`ResampleLinear`/`MonoDownmix`/`FloatToPCM16`/`DeviceError`/帧常量 | `internal/audio/*.go`（声明处见 §1.2） | **全部 0**（逐枚跑过，见 4.0；根因＝R1 无 importer） | **缺一整块调用方**——实现本体是全的，且不是桩 |
| A2 | `SpawnCapture`（自称 capture 线程的"唯一 sanctioned entry point"） | 声明 `audio.go:180-182` | **全仓 0，连本包内都没用它**；两枚真源各自直接 `observe.Default.Spawn`（`wasapimic_windows.go:83`、`wavinjector.go:84`） | **补一行就通，且顺手修一处自相矛盾**：文档说唯一入口，实现绕过它；且这里写死 `observe.Default`，而球侧是注入式（`ball_windows.go:170 opts.Registry.Spawn`） |
| A3 | `WASAPIMicrophone.Endpoints()`（D47 留给 P15 AEC spike 的采集/渲染端点对） | `wasapimic_windows.go:137` | 0 | 补一行就通（但要等 AEC，不着急） |
| A4 | `HalfDuplexGate` 全套（`SetSpeaking`/`SetMuted`/`Open`/`Muted`/`Path`/`Stats`） | `gate.go:130/168/209/202/216/221` | 0 | **补一行就通**：包内测试已在跑同一条逻辑（`gate_test.go:11-16` 走真 `WavInjector`） |
| B1 | `Ball.SetAudioLevel` | `liquid_windows.go:42` | **1，且是调试载体**：`cmd/balldebug/main.go:418`、`:477`（R4；flag 自述逐字 `synthetic audio envelope 0..1 pushed to the ball at ~30fps`，`:396` 逐字 `syllabified synthetic envelope`） | **补一行就通**（前提是表一 #6/#7/#8 那三格先有电平值） |
| B2 | `EnablePrototypeVisuals` / `PrototypeVisualsEnabled` | `statevisual.go:106/109` | **各 1，都在 balldebug**：`cmd/balldebug/main.go:122`、`:451` | 补一行就通；但**默认翻不翻不归写腿**（`statevisual.go:88-96` 自陈翻默认＝票 68 AC#2） |
| B3 | `LookNames` / `SetLook`（皮肤 seam，编排者 A461 已点名 `tokens.go:316/332`） | `internal/ball/tokens.go`（`LookNames`、`SetLook`） | `ball.LookNames` 2、`ball.SetLook` 1，**全在 balldebug 的 flag 与开关** | 补一行就通（样式维度 owner 已说"以后再说"，本腿不排它） |
| B4 | `NewHotkeyReloader` | 声明 `internal/ball/hotkey_reload.go:73`（`:16` 是它自己的用法注释） | 2，**全在 `cmd/balldebug/main.go:237`** | 补一行就通（GUI 腿要有 config 驱动的快捷键时接） |
| C1 | `EvKwsEnabled` / `EvWakeWord`（D43 #5/#7） | `events.go:18/20`，表行 `table.go:55/61` | **0 生产者**（R15） | **缺一整块**：要发它，先得有 KWS loader |
| C2 | `Facts.KwsLoaded` | `events.go:70` | 生产 0；只有测试在填（`internal/ball/hotkey_live_test.go:238/254/272`、`window_test.go:134`） | 补一行就通（loader 起来后填） |
| C3 | **`Options.Sink` 副作用注入口**（`Effect{D43,Name}` 全套机制） | `machine.go:26-34/39/64-68`，发射 `:162-164` | **生产两处建机器都没传**：`cmd/balldebug/main.go:188`、`cmd/wisp/models.go:303` ⇒ 走显式 no-op | ⭐ **这是全表最"值"的一格：补一个装配根 Sink 就通**，被空掉的副作用名含 `kws.load`、`kws.pause`、`kws.stop-inference`、`kws.keep-alive-alert`、`speech.load-vad-asr`、`audio.stop-capture`、`audio.discard-buffer` |
| C4 | 上述七枚副作用名的**执行器** | 表行见 `table.go:56`（`kws.load`）、`:62`（`kws.pause`＋`speech.load-vad-asr`）、`:66`（`kws.stop-inference`）、`:71`（`kws.keep-alive-alert`）、`:52`（`speech.load-vad-asr`）、`:87`（`audio.stop-capture`）、`:94`（`audio.discard-buffer`） | 0 | **缺一整块**（C3 是它的投递机制） |
| D1 | `ChannelKWS` + `SetLoaded` 的 KWS 分支 | `approval/approval.go:57/92/99/105/176`、`gate.go:137` | 生产 `SetLoaded` 有 1 处（`cmd/wisp/approval_reply.go:416`）但翻的是 `s.replyVeto`，`run.go:456-457` 逐字 `replyVeto stays unset` ⇒ **KWS 那一支 0** | **补一行就通**；⚠ 不许顺手把 `DefaultChannels()` 改成含 KWS（`approval.go:161` 那句诚实陈述与 `window_test.go:97` 的钉会被打红，那是"假装通道存在"） |
| E1 | `WakeWord` 四字段（`Enabled`/`Keywords`/`Thresholds`/`VetoWords`） | `schema.go:196-203`；热载分级 `manager.go:370-375/385-386/398-399/405-407` | **0**（R13：`internal/config/` 之外 `.Voice.` 零命中） | **缺一整块读者**；配置侧本身**不用新写一行** |
| E2 | `[voice]` 其余小节：`ASRConfig`/`TTSConfig`/`AECConfig`/`RealtimeConfig`/`ConversationMode`/两条云链 | `schema.go:207-236/240-250` | **同上，0** | 缺一整块读者（ASR/TTS 那两块的"识别"至少还在 schema 里，引擎侧同样零） |
| E3 | `localSherpaProvider = "local-sherpa"` | `config/catalog.go:20`、`:57`（校验里拒它）、`schema.go:210/216` 默认值 | 只有**校验路径**用它（拒绝用），**没有任何一处用它去起本地引擎** | 补一行就通（但要等 E1/E2 的读者） |
| F1 | `SLOArmed` 那一档（`memCapArmed`/`cpuLimitArmed`/`goroutineLimitArmed`） | `observe/thresholds.go:20/27/38`；状态声明 `sampler.go:43`、进表 `:51` | **`internal/observe` 之外 0**（实跑 `grep 'SLOArmed' --include=*.go cmd tools internal \| grep -v '^internal/observe/'` → rc=1）。⚠ 但旗子是收的：`cmd/wisp/slo_windows.go:149` 的 usage 明写 `-state <Sleeping\|Armed\|Warm\|Conversation\|PanelOpen\|WorkPeak>` 并用 `SLOState(*state).Valid()` 校验；同文件**零 import ball、零 SetState**（`grep 'ball\.' cmd/wisp/slo_windows.go`＝空，R11 亦证 ball 只有 balldebug 一枚 import 者）⇒ **`-state Armed` 只是给采样窗口贴标签，并不能把进程带进 Armed** | **缺一整块**（真能进 Armed 的前提＝C1/C3；在那之前 Armed 行的读数一律是"标签读数"，别拿它当验收） |
| F2 | `captureBufferLimitConversation = 20 << 20` | `thresholds.go:42` | 注释自陈 `not measurable until the speech pipeline (tickets 15/26) runs, so it is recorded as a note, never as a fake pass` ⇒ **不进门禁** | 缺一整块（诚实挂着，别去翻它） |
| G1 | `OnDemandNames = {"kws-infer"}` | `observe/goroutine.go:32/47` | **零 Spawn**（R17） | **补一行就通**（槽是给 KWS 留的） |
| G2 | `ResidentNames` 里的 `hotkey-listener` | `goroutine.go:44`；`SPEC-01:119` | **零 Spawn**（R16） | 补一行就通（GUI 腿接快捷键时） |
| H1 | KWS 模型本体与分发 | `models/manifest.json` 里 `kws-zipformer-wenetspeech-3.3M-2024-01-01`（R19，逐文件哈希含 `keywords.txt`）；`manifest.go:94 Purposes` 含 `kws`；`downloader.go:228 VerifyInstalled` | **生产码里无人按这个 ID 请求它**（实跑 `grep 'kws-zipformer-wenetspeech' --include=*.go internal cmd tools`（排测试）→ 零命中） | **补一行就通**（`Manager.Ensure(ctx, id)` 已在 `handOffModel` 那条腿上通用） |
| I1 | `internal/speech` **整包** | 只有 `doc.go`（21 行），`:19 DEFERRED(engines) … ticket 41 (KWS)` | 0 importer（R2）；`grep 'type .*Engine'` 主模块 0 命中 | **缺一整块**（不是半成品，是空地基） |
| I2 | `internal/watchdog` **整包** | 只有 `doc.go`，`:16-17 DEFERRED(watchdog loop/thresholds): implemented by ticket 42`，正文逐字含 `never unload KWS in Armed` | 0 | **缺一整块**；D43 #9 的 `kws.keep-alive-alert` 与 PLAN.md:514 的"按态查表"最终都要落在这（票 223 票面 `:15` 已把这圈 tick 登记成推迟件） |
| J1 | **spike 里那两枚真 KWS 会话** | `scripts/spike/model-residency/main.go:114-152`、`scripts/spike/speech-baseline/main.go:124-167` | 生产 0、测试 0；**在另一个 module 里**（R7） | **不能"补一行"**：跨模块引不到，只能照抄形状进 `internal/speech` |
| K1 | D11 控制词的宿主挂钩 `ControlHandler` / `Options.Control` | `internal/agent/control.go:69-73`、`loop.go:158` | **从未被赋值**（实跑 `grep 'Control:'` 生产与测试都无一处设它）⇒ 永远走 `defaultControl`（`loop.go:524`），`repeat/louder/confirm` 恒 `unhandled`（`loop.go:522-523` 注释逐字） | **补一行就通**；ASR 一通，`voice.wake_word.veto_words`（E1）与这枚是近邻，别各造一套 |

### 4.2 名册统计（口径写清，免得被当成"裁过"）

- 表内 **24 行**（A1-A4、B1-B4、C1-C4、D1、E1-E3、F1-F2、G1-G2、H1、I1-I2、J1、K1）。
- 判"**补一行就通**"（含"实现齐、只缺调用方"）：**11 行**＝A2、A3、A4、B1、B2、B3、B4、C2、D1、G1、G2。
  ⚠ 其中 **B1 与 C2/D1 是有前提的"一行"**：B1 要 A1/A2 那侧先有电平值，D1/C2 要先有 loader。
- 判"**缺一整块**"：**13 行**。
- **本表不是穷举**：穷举需要"未使用导出符号"仪器（编译器级），本腿按红线没跑 go 工具链，
  所以这是"想到的候选全查完"，不是"全仓无漏"（这条写进 §5）。

---

## §5 我答不上来的

（待写，本节不许为空。）
