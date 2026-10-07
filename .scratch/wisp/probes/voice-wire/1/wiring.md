# voice-wire-1 — 语音链路逐跳接线普查（只读）

> 腿 id `voice-wisp-1`（工单 `voice-wire-1`）。**只读**：本件之外零改动。
> 射程＝麦克风进来 → 送进模型 → 声音出去，七跳，逐跳量「存在吗／被谁调用／缺的是一行还是一整块」。
> 结论三档：**〔已写零调用＝接线〕／〔写了一半〕／〔一块没写〕**，每档都带尺（命令＋读数）。

---

## 0. 起手锚

| 项 | 读数 | 尺 |
|---|---|---|
| HEAD | `82a10da4`（branch `dev`） | `git rev-parse --short HEAD` |
| 工作树脏行 | **752** | `git status --porcelain \| wc -l` |
| 本腿写点 | `.scratch/wisp/probes/voice-wire/1/**`（本件＋`logs/` 两枚临时件） | — |

> ⚠ **起手锚是"我进来时"的 HEAD，不是"我交件时"的 HEAD。** 本腿取数期间 `68ff68d5`／`7ec2f518`／`9aef5939` 三枚别人的探针件落了盘（`git log --oneline -3` 现读）⇒ §1 跳 1 那两把"零 importer"尺读数**落笔时未过期，但引用前必须重跑**（`272-r2` 此刻正在写 `cmd/wisp` 测试，见 §4 第 4 条）。

**⛔ 我没跑 `go test`／`go build`／`go vet`／任何 `./...`／任何突变。** 理由具名：此刻 `272-r2` 正在 `cmd/wisp` 跑 Go 突变、`111-r5` 在写 `.github/workflows/ci.yml`，任何整包测试都会洗掉 `272-r2` 的红名册读数。本件**全部**结论来自 `grep`／`ls`／`wc`／`sed` 读原文，尺的口径锚在**调用形状**（带括号的 `x.Method(` 或 `pkg.Fn(`），不是裸符号名。

**⚠ 与编排者转述不符之处（先行具名，详证见 §1/§5）：**
1. 「`cmd/wisp/*.go` 里有 **3 处提到 `internal/audio`**」——实际 `grep -rn 'internal/audio' cmd/wisp/` 命中 **1 处**（`cmd/wisp/config_readers_255.go:133`）。`cmd/wisp/models.go:18` 与 `:25` 那两行提到的是 **`internal/speech`**，不是 `internal/audio`（原文：`:18` "internal/speech is a doc.go-only boundary stub"、`:25` "internal/speech owns adding its own handOffModel call"）。
2. `internal/audio/` = 16 枚 `.go` ✓，但其中 **5 枚是 `_test.go`**（`gate_test`/`hotplug_test`/`level_test`/`resample_test`/`wavinjector_test`）⇒ 非测试产码 **11 枚／3849 行（含测试）**。
3. 「全仓零 importer」✓ 复现为 0，且口径可再收紧：**含测试的全口径也是 0**（不只是非测试）。
4. `internal/speech/` = 只有 `doc.go` ✓（1066 字节）。
5. **新发现（编排者未提）**：主模块 `go.mod:8` 已直接依赖 `github.com/k2-fsa/sherpa-onnx-go v1.13.8`（＋`:21-23` 三枚平台模块），且 `cmd/wisp/doctor.go:60/:81`、`cmd/wisp/main.go:171` 在**生产产码里真调用** `sherpa.GetVersion()`／`sherpa.GetOnnxruntimeVersion()`。⇒ sherpa 的 cgo 绑定**已经链进出货二进制**，缺的只是推理，不是链接。

---

## 1. 跳表

> 每跳固定五栏。**尺一律锚在调用形状**；「追到 main」＝能写出 `main.go:行 → … → 那一行` 的完整链；追不到的「调用者」标〔非生产／测试载具〕。

### 跳 1 · 采集／电平（麦克风进来，有没有人读电平）

- **存在吗**：**在，而且是全链最完整的一块**。`internal/audio/` 非测试 11 枚：`audio.go`(182, C8 `AudioSource` 接口＋有界帧道 `NewBoundedFrames`/`BoundedFrameCapacity`＋`SpawnCapture`)、`wasapimic_windows.go`(288, `WASAPIMicrophone` 共享模式＋pinned 线程)、`wasapi_windows.go`(430, `wasapiOpener.Open`)、`mmdevice_windows.go`(283, 设备枚举＋热插拔)、`device.go`(132)、`level.go`(147, `LevelOfSamples:96`/`FrameLevel:118`/`DecodeFrame:137`)、`gate.go`(259)、`resample.go`(139)、`wavinjector.go`(213, C8 测试载具)、`wasapi_other.go`(29)、`doc.go`(29)。
- **调用形状尺**：
  - `grep -rn '\baudio\.[A-Z][A-Za-z]*(' --include=*.go . | grep -v '^\./internal/audio/' | grep -v '_test\.go' | grep -v '^\./\.scratch'` ⇒ **0**（包限定调用点，全仓产码零枚）。
  - `grep -rln '"github.com/CarlosShao/wisp/internal/audio"' --include=*.go .` ⇒ **0**（含测试的全口径 importer）。
  - 电平生产者的包内自调：`internal/audio/level.go:118 FrameLevel(` → `:123 DecodeFrame(` → `:127 LevelOfSamples(`——**只有 level.go 自己这一条内部链**，没有外部消费者。
  - `SpawnCapture(` 全仓命中 **1**＝定义处 `audio.go:180`；**调用者 0**。
  - 球那一侧入口 `SetAudioLevel(`：定义 `internal/ball/liquid_windows.go:42`；产码调用点只有 `cmd/balldebug/main.go:418`／`:477` ⇒〔**非生产／测试载具**〕（`cmd/balldebug` 是调试 exe，不是 `cmd/wisp`）。
- **追到 main 吗**：**不能**。`cmd/wisp/main.go` 的 `switch args[0]`（`:89-124`）里没有任何采集/音频分支，`cmd/wisp` 也不 import `internal/audio`。
- **三档**：**〔已写零调用＝接线〕**。
- **票**：**247**（OPEN，判据 **10** 枚未勾／0 已勾）——本跳已被裁完：`A485` 八问（P1 默认档谁挡麦／P4 协程 owner／P6 电平在哪侧算／P8 名册无槽…）逐条给了甲乙，落地腿 `247-r1` 已批准、排在 `246-r2` 之后。上游＝票 **13**（DONE）。⚠ 票 247 明写「唤醒词／ASR／TTS **都不属于本票**」，且 `AC#7` 把 `internal/speech` 列为越界路径。

### 跳 2 · 唤醒词（KWS）

- **存在吗**：**引擎一块没写**；**外围三面都已写好**：
  - 状态机面：`internal/statemachine/table.go` 行 #5 `kws.load`(`:56`)、#7 `EvWakeWord`+`kws.pause`(`:61-62`)、`kws.stop-inference`(`:66`)、`kws.keep-alive-alert`(`:71`)，守卫 `KwsLoaded`(`:76/:80/:232/:238`)——表和守卫**在**。
  - 配置面：`internal/config/schema.go:196 WakeWord`、`:247`、`tiers.go:54-55/73-74`、`manager.go:380-416`（reload/hot 分层）——**在**。
  - 否决通道面：`internal/agent/approval/approval.go:57 ChannelKWS`、`:110 "说取消词"`、`replies.go` 的 `vetoChannel`——**在**。
- **调用形状尺**：
  - `grep -rn 'WakeWordEngine' --include=*.go internal cmd tools` ⇒ **2**，两枚都在 `internal/speech/doc.go:1-2` 的**注释**里 ⇒ 产码 **0**。
  - `grep -rn 'sherpa\.NewKeywordSpotter(' --include=*.go .` ⇒ **3**，全在 `scripts/spike/speech-baseline/main.go:128`、`scripts/spike/model-residency/main.go:116/142` ⇒〔**非生产**：`scripts/spike/` 是独立模块（`scripts/spike/go.mod`）的 RSS 基线测量件〕。
  - `EvKwsEnabled`/`EvWakeWord` 的产码发射者：`grep -rn 'EvKwsEnabled\|EvWakeWord' --include=*.go internal cmd tools | grep -v '^internal/statemachine/' | grep -v '_test\.go'` ⇒ **0**。
  - `Facts.KwsLoaded` 的产码写入者 ⇒ **0**（`events.go:70` 只有字段定义，读侧全是 table）。
  - `SetLoaded(` 产码命中 **3**（`cmd/wisp/approval_reply.go:512`、`resident_approval_windows.go:567/:787`），三枚传的都是 **Esc**/`vetoChannel`，**没有一枚传 `ChannelKWS`**。
  - 配置读取者：`cmd/wisp/config_readers_255.go:169-170` 自己写着 `hotClaimNoReader "扫描零命中：nothing outside internal/config reads cfg.Voice"`。
- **追到 main 吗**：**不能**。
- **三档**：**引擎〔一块没写〕**；状态机/配置/否决三面〔已写零调用＝接线〕。**缺的是"谁去 `sherpa.NewKeywordSpotter` 并把 `EvWakeWord` 打进 Machine"这一整块**，不是一行。
- **票**：**41**（OPEN，**6** 枚未勾）；`internal/speech/doc.go:20` 逐字把 KWS 归口票 41。

### 跳 3 · VAD／分段

- **存在吗**：**一块没写**（连包边界注释都在互相推）。
  - ⚠ 文档级冲突具名：`internal/audio/doc.go:3` 声称本包负责 "…VAD, and the half-duplex switch"，但 `internal/audio/` 目录里**没有 VAD 文件**；`:22` 又写 "no ASR/TTS/KWS model inference (speech)"；`internal/speech/doc.go` 未列 VAD；而 `internal/statemachine/table.go:52/:62` 的副作用名是 **`speech.load-vad-asr`**、`models/manifest.go:94` 把 `"vad"` 登记成 speech 侧的 model purpose。⇒ **VAD 的归属在文档里是两说的**，接线前必须先定这一枚（见 §2 与 §4）。
- **调用形状尺**：
  - `grep -rn 'Vad\|VAD' --include=*.go internal cmd tools` ⇒ 产码命中全在**注释**（`audio/audio.go:41,45`、`audio/doc.go:3`、`audio/level.go:112`）；唯一非注释产码相关＝`internal/statemachine/events.go:25 EvVadStop` 与 `table.go:85` 的行本身。
  - `EvVadStop` 产码发射者（排除 `internal/statemachine/` 与 `_test`）⇒ **0**。
  - `grep -rn 'sherpa\.NewVoiceActivityDetector(' --include=*.go .` ⇒ **2**，全在 `scripts/spike/model-residency/main.go:157/179` ⇒〔非生产〕。
  - 反向对照（证明这把尺命中得了真名）：`internal/observe/thresholds.go:29 cpuLimitConversation = 5.0 // VAD always-on` 说明 VAD 在 SLO 里已被计过，而 `internal/models/manifest_real_test.go:118 TestRealDownloadVadThroughPipeline` 证明 **silero VAD 那 0.64MB 模型的下载链已经跑通**（真的拉到、真的过 hash）。⇒ 缺的只有推理侧。
- **追到 main 吗**：**不能**。
- **三档**：**〔一块没写〕**。
- **票**：**15**（OPEN，**6** 枚未勾，标题逐字 "Speech engines: **VAD** + streaming ASR via sherpa, engine slot mutex, CER harness"）。

### 跳 4 · ASR（本地 onnx 路径 vs 云端；`third_party/sherpa-onnx` 在本机的角色）

- **存在吗**：**两条路径的引擎都没写；运行时候件与依赖已经就位，且已链进出货二进制。**
  - `third_party/sherpa-onnx/` ＝ **3 枚 DLL**（`onnxruntime.dll`／`sherpa-onnx-c-api.dll`／`sherpa-onnx-cxx-api.dll`），目录内 `.go` 文件 ⇒ **0**。它的角色是**运行时分发＋`wisp doctor` 的"与 exe 同目录/版本对得上"自检项**（`cmd/wisp/doctor.go:71-81`），不是绑定层；绑定层是 `go.mod:8 sherpa-onnx-go v1.13.8`（Go 侧包，已 vendor 在模块缓存里）。
  - 模型分发侧**已完备**：`internal/models/manifest.go:94 Purposes = {kws, vad, asr-streaming, asr-offline, punctuation, tts}`，`internal/models/downloader.go:140 Manager.Ensure` 返回落盘目录。票 **14**（DONE）。
- **调用形状尺**：
  - `grep -rn 'AsrEngine' --include=*.go internal cmd tools` ⇒ **2**，全在 `internal/speech/doc.go` 注释 ⇒ 产码 **0**。
  - `grep -rn 'sherpa\.[A-Za-z]+(' --include=*.go internal cmd tools scripts` 非测试 ⇒ **生产 3 枚**：`cmd/wisp/doctor.go:60 sherpa.GetVersion()`、`:81 sherpa.GetOnnxruntimeVersion()`、`cmd/wisp/main.go:171 sherpa.GetVersion()`；**推理 0 枚**。
  - `grep -rn 'sherpa\.NewOnlineRecognizer(' --include=*.go .` ⇒ **4**，全在 `scripts/spike/model-residency/main.go:184` 与 `scripts/spike/xy-verdict/session_cgo.go:33` ⇒〔非生产〕。
  - **正向对照（证明这把尺命中得了真名，也证明链接是通的）**：`wisp doctor` 的链 `cmd/wisp/main.go:95 case "doctor" → :97 cmdDoctor() → cmd/wisp/doctor.go:60 sherpa.GetVersion()` ——这条**能追到 main**，但它只打印版本号，不推理。⇒ 「sherpa 不可用」是假话，「sherpa 只被用来打印版本」才是现状。
  - 云端 ASR（C9／票 61）：`grep -rn 'CloudASRChain' --include=*.go internal cmd` ⇒ **7**，全部落在**配置面**（`config/schema.go:259-263`、`catalog.go:30-31` 校验、`manager.go:386/:403` reload 比对、`tiers.go:71`）⇒ **消费端 0**。`ls internal/llm/` 只有 `openaichat`/`openairesponses`/`anthropic`/`chain.go`/`probe.go` 这类**文本 LLM** 适配器，没有任何 ASR/TTS 传输层。
- **追到 main 吗**：**引擎侧不能**；只有版本号打印那一支能（且它不产生转写）。
- **三档**：本地 ASR **〔一块没写〕**；云端 ASR **〔一块没写〕（只有配置面〔已写零调用〕）**；运行时/依赖/分发 **〔已写〕**。
- **票**：**15**（本地＋CER harness，6 枚未勾）／**61**（云端 C9 级联，6 枚未勾）／**27**（标点复原，5 枚未勾，`Purposes` 里的 `"punctuation"` 等它）／**14**（模型分发，DONE）。

### 跳 5 · 送进模型（`internal/speech` 与 `cmd/wisp/models.go` 那枚 handOff 桩的关系）

- **存在吗**：**两个半边都在，但接的不是对方。**
  - 消费端（模型入口）**已接好且能追到 main**：`internal/agent/loop.go:339 Loop.Run(ctx, input string)` / `:328 RunAsync`。生产调用点：`cmd/wisp/run.go:1099 bg := loop.RunAsync(ctx, task)`，链 `cmd/wisp/main.go:89 case "run" → :91 cmdRun(args[1:]) (run.go:151) → … → run.go:1099`（另一支内部调用者 `internal/tools/subagent_197.go:344 child.RunAsync(childCtx, prompt)`）。⚠ **`input` 今天的来源是命令行 argv 的文本，不是声音**——全仓没有一枚"转写文本 → `Loop.Run`"的调用点。
  - 交还端（模型字节）**也已接好、能追到 main**：链 `cmd/wisp/main.go:102 case "models" → :109 cmdModels(...) → cmd/wisp/models.go:293 store.handOffModel(io_, id, logf) → :302/:303 statemachine.New → :304 models.WireDownloading → :311 bridge.Run(ctx, id) → internal/models/bridge.go:58 VerifyInstalled`（`bridge.go:40/:46/:59/:64` 是产码里**唯一**的 `.Dispatch(` 命中，4 枚，全是下载事件）。
  - **两者的关系＝还没有关系**：`cmd/wisp/models.go:25` 逐字写着 "internal/speech owns adding its own handOffModel call; this file cannot become that reader by wishing"；`:18-19` 逐字 "That capability does not exist in this tree: internal/speech is a doc.go-only boundary stub (DEFERRED: ticket 15/26/41 per its own doc.go)"。`internal/speech` 的 importer ⇒ **0**。
  - 更硬的一条：`cmd/wisp/models.go:36-41` 自己声明 "There is still no reader of model bytes anywhere in `cmd/wisp`'s dependency graph - this process never opens a file under `<store>/<id>/`, it only re-hashes it"。尺的复现：`grep -rn '^func ' internal/models/*.go`（非测试）里 **没有** 任何 load/engine/session/infer 语义的函数；`Manager.Ensure` 返回的是目录路径 `string`，**没有下一跳**。
- **三档**：**文本入口〔已接线〕／模型字节消费端〔一块没写〕**。缺的不是"某一行的桩"，是**两枚新物件**：① `internal/speech` 里那枚吃帧道＋模型目录的引擎；② 把转写文本交进 `Loop.Run(ctx, input)` 的那一次调用。
- **票**：**15/26/41**（引擎本体）＋**121**（DONE，交还链进二进制＋装配可达性门）。⚠ **「谁把 ASR 终稿喂进 `Loop.Run`」这一枚，票池里我没有对应到号**（见 §4）。

### 跳 6 · TTS／出声

- **存在吗**：**引擎与播放流**都没写；**状态机副作用名＋风险通道名**已写好。
  - `internal/audio/` **没有播放实现**：目录里唯一 render 相关的产码是**端点查询**——`device.go:31-32 DefaultRenderDevice()` 接口、`mmdevice_windows.go:108-111 w.defaultEndpoint(eRender)`、`wasapi_windows.go:89 eRender = 0`（常量）。`wasapimic_windows.go:142 render, err = m.watcher.DefaultRenderDevice()` 是**采集协程去查 render 端点**（D47 的 capture/render 配对／P15 AEC 用），**不是往扬声器写**。`internal/audio/doc.go:27` 逐字：`DEFERRED(playback): TTS output lands with ticket 26 on the same WASAPI technique.`
  - 状态机面：`table.go:192 tts.stop`/`audio.release-output`、`:201/:207 tts.stop-bargein-400ms`/`asr.exclude-playback`、`events.go:46 EvSpeakDone`/`:47 EvInterrupt`/`:48 EvBargeIn`——**名在，执行器不在**。
  - 风险面：`internal/risk/provenance.go:122 ChTTS = "tts.announce"`，`:823` 逐字 "no tool name yet: TTS announce text (ticket 26)"。
- **调用形状尺**：
  - `grep -rn 'TtsEngine' --include=*.go internal cmd tools` ⇒ **2**，全在 `internal/speech/doc.go` 注释。
  - `grep -rn 'sherpa\.NewOfflineTts(' --include=*.go .` ⇒ **1**，在 `scripts/spike/model-residency/main.go:220` ⇒〔非生产〕。
  - `EvSpeakDone`/`EvInterrupt`/`EvBargeIn` 产码发射者（排除 `internal/statemachine/`、`_test`）⇒ **0**。
  - 产码里向扬声器写帧的调用形状（`GetAudioClient`/`IAudioRenderClient` 一类）在 `internal/audio/` ⇒ 未见 render 流；只有 `Activate`+`eCapture` 路径（`mmdevice_windows.go:105 defaultEndpoint(eCapture)`）。
- **三档**：**〔一块没写〕**（引擎＋播放流两块都缺，缺一整块）。
- **票**：**26**（OPEN，**7** 枚未勾，标题逐字 "TTS output: sherpa matcha engine, serial slot, **Speaking pipeline**, P7 gate"）／**61**（云端 TTS）。

### 跳 7 · 半双工 vs 全双工（D47：半双工只约束干活路径，陪聊路径全双工＋可选云端 realtime 大脑）

- **存在吗**：**两条路径的形状都写在文档与表里，实现只有"半双工的那枚门"存在，且没人用**。
  - **Path T（干活路径／半双工）**：`internal/audio/gate.go`（259 行）**已写完**——`HalfDuplexGate` 包 `AudioSource`，`:14` 逐字 "Path T (work path): while TTS is speaking the gate CLOSES capture"，`:130 SetSpeaking()`／`:168 SetMuted()`／`:57 WithStartMuted()`／`:63 WithGateEvents()`／`:234 effectiveOpen()`／`:247 openInnerLocked()`；`audio.go:70 type Path string`、`:68-69` 逐字 "half-duplex by design (D16); C = conversation path, full duplex with AEC (AEC itself is ticket 26/59 - the gate only stops policing it)"。⇒ 与编排者给的 D47 口径**一致**（门只切换 Path T 的采集；Path C 不由它管静音）。
  - **Path C（陪聊路径／全双工）**：**一块没写**。`doc.go:24` 逐字 "no AEC source (Path C AEC is tickets 26/59; the gate interface carries the Path switch, not the processing)"。
  - **云端 realtime 大脑（D47 的"可选"那一支／C32）**：配置面**已写完**——`config/schema.go` `voice.realtime`、`catalog.go:121-125 validateRealtime`、`loader.go:33-34 RealtimeKey`＋`:238-244` 真的把 `APIKeyRef` 解析成密钥、`manager.go:133 return &Resolved{ProviderKeys: pk, RealtimeKey: …}`、`tiers.go` `voice.aec`/`realtime.*` reload 分层。**消费端 0**。尺：`grep -rn 'RealtimeKey' --include=*.go internal cmd`（非测试）⇒ **4**，四处全在 `internal/config` 自己内部（定义/赋值/返回），**包外读取者 0**。`AsrEngine/TtsEngine/WakeWordEngine/RealtimeEngine` 四字合并命中 ⇒ **2**，全在 `internal/speech/doc.go` 注释 ⇒ **C32 产码 0**。
  - **会话保活（全双工的落地载体／C31）**：`internal/session/` 只有 `doc.go`＋`session.go`(ID `Mint`/`Valid`)＋`grants.go`(Ledger)。`grep -rn 'Warm\|Conversation' internal/session/`（非测试）⇒ 只有 `doc.go:2/:9` 两句**声明**，产码 0。
- **调用形状尺（门那侧）**：`grep -rn 'NewHalfDuplexGate(' --include=*.go internal cmd tools`（非测试）⇒ **1**＝定义处 `gate.go:84`；**调用者 0**（`SetSpeaking(`/`SetMuted(` 产码调用者同样 0，命中全在 gate.go 自身与注释）。
- **追到 main 吗**：**都不能**。
- **三档**：Path T 的门 **〔已写零调用＝接线〕**；Path C／AEC／realtime 大脑／会话保活 **〔一块没写〕**（其中 realtime 是"配置面〔已写零调用〕＋引擎〔一块没写〕"）。
- **票**：**59**（P15 AEC spike，**7** 枚未勾）／**60**（C32 RealtimeEngine，**7** 枚未勾）／**28**（C31 SessionScope Warm/Conversation，**6** 枚未勾）／**247**（半双工门所依附的采集腿，10 枚未勾）。
- **⚠ 全链公共缺口（不属于任何一跳、但每一跳都会撞上）＝状态机副作用没有执行器**：`internal/statemachine/table.go:15` 逐字 "SideEffects are event names only: the machine fires them on the EffectSink"；`machine.go:23-25` 逐字 "Effects are EVENTS: the machine never executes them itself - consumers (session, speech, panel, tickets 08+) wire a Sink. The zero sink is an explicit no-op."；`Options.Sink`(`machine.go:39`) `nil = explicit no-op`；`machine.go:156-166` 把 41 张副作用名（含 `speech.load-vad-asr`／`kws.load`／`tts.stop-bargein-400ms`）投给 sink。尺：`grep -rn 'statemachine\.New(' --include=*.go internal cmd tools`（非测试）⇒ **2**：`cmd/wisp/models.go:303`（`Options{Initial: …}`，**没传 Sink**，且是一枚 `defer machine.Close()` 的一次性机器）与 `cmd/balldebug/main.go:188`〔非生产〕。⇒ **出货进程里没有任何一枚带 Sink 的 Machine，也没有任何一枚常驻 D43 机器**（此读与票 **246** 自述的"常驻腿不持有 Machine"一致）。**`Options.Sink` 的产码传参点 ⇒ 0。**

---

## 2. 接线顺序（先接哪一跳能让后面少返工；⛔ 不新造契约）

**第一跳＝跳 1（采集／电平，票 247）。** 三条理由都落在"少返工"上，不落在"哪一跳最容易"：

1. **它是七跳里唯一一处判据已经裁完的纯接线**。票 247 的 `A485` 已把四个会返工的点逐条写死：落点（`cmd/wisp` 拥有采集，新增包级边恰好 1 枚）、协程 owner（P4：改吃传入 registry，`audio.go:180 SpawnCapture` 就是那枚现成接缝）、电平在哪侧算（P6：采集腿内侧 `FrameLevel` → 只过一枚 `float32`）、名册无槽怎么办（P8：在**已有**的 `audio-capture` 协程内算完再投递，D38(b) 六名零膨胀）。任何**别的**跳先做，都得把这四个决定**再重做一遍**。
2. **它产出的是每一跳共用的载具**：pinned 采集线程＋C8 有界帧道（`audio.go:61/:65`）＋被 `HalfDuplexGate` 包好的 `AudioSource`。KWS/VAD/ASR 的输入都是这枚帧道；先立载具，后面每枚引擎落地时是"往已有帧道上挂一个消费者"；反过来先写 ASR，它今天只能吃 `WavInjector`，接真麦时要**重包一层门＋重挂一次 owner＋重走一遍降级**（票 247 AC#6 那三形）。
3. **它不碰任何一块没写的包**：票 247 `AC#7` 把 `internal/speech` 明列为越界退回路径，所以第一跳的 diff 与第二跳之后的 diff **天然不重叠**，串行代价最小（本票当前也正按 `246-r2`→`247-r1` 排着）。

**第二跳＝给常驻进程立一枚带 `Options.Sink` 的 D43 机器（跳 7 末尾那枚公共缺口）。** 理由：`speech.load-vad-asr`／`kws.load`／`tts.stop` 这三条**已经写在冻结的 D43 表里**，Sink 不落地时，每一枚引擎都只能各自发明一条"从状态到动作"的私路，引擎写完再往表上收，是四遍返工。⚠ 这一枚**票池里没有对应号**（见 §4），且它**不动转移表一字**（只是把那 41 张名投给一个执行器），所以不需要人工批契约——但要不要新立一枚票由编排者定，本腿⛔不建票。

**第三跳＝跳 4 的本地 ASR（票 15，连带跳 3 的 VAD）。** 理由：有了帧道＋Sink，ASR 才有可挂的输入与可驱动的装载副作用；且票 15 自带 CER harness（量尺先于实现落），此时接能立刻验。**跳 4 之前不要接跳 2**——KWS 与 ASR 共用 `sherpa` 会话层与"engine slot mutex"，先做 KWS 会把会话生命周期写两遍。

**第四跳＝跳 6（TTS，票 26）**，因为它同时是 `HalfDuplexGate.SetSpeaking` 的**唯一自然调用者**（`gate.go:128` 逐字 "SetSpeaking is driven by the TTS side"）：门在跳 1 已经立起来，但只有 TTS 落地才有人按它，这也是一条"越晚接越不用重测门"的顺序理由。

**第五跳＝跳 7 的 Path C（59→60→28）与跳 4 的云端（61）。** 理由：D47 把全双工定成"陪聊路径＋可选云端大脑"，它是**增强**不是主干；AEC（59）是 spike 前置，realtime（60）是 gated；把这两块排在主干之后，主干的任何返工都不会传染它们。

**顺序第一跳之外的一条硬前置（不是跳，是归属问题）**：接 VAD 之前必须先定 **VAD 归 `internal/audio` 还是 `internal/speech`**（`audio/doc.go:3` 说 audio 负责 VAD，`table.go:52` 的副作用名说 `speech.load-vad-asr`，`speech/doc.go` 未列 VAD）。两说会导致同一枚引擎被写两遍或写进被 `AC#7` 禁掉的包里。⛔ 本腿不裁这一枚，交编排者（§4）。

---

## 3. 查重结论（每一跳 ↔ 已有票号；⛔ 本腿不建票）

| 跳 | 已有票 | 状态 | 未勾 AC 枚数 | 已勾 | 覆盖了这跳的什么 |
|---|---|---|---|---|---|
| 1 采集／电平 | **13** | DONE | — | — | 采集栈本体（C8/WASAPI/pinned/热插拔/电平生产者半） |
| | **247** | OPEN | **10** | 0 | **接线本体**（＋AC#1 非测试 importer≥1、AC#3 只过 float32、AC#10 "接上≠看得见"） |
| | 68（AC#2） | — | 未量 | — | `prototypeVisuals` 默认值的翻不翻（票 247 AC#10 具名归口给它） |
| 2 唤醒词 KWS | **41** | OPEN | **6** | 0 | 引擎＋Armed 态＋keywords/veto 词＋静音＋opt-in 隐私 |
| 3 VAD／分段 | **15** | OPEN | **6** | 0 | 标题逐字含 VAD（与 ASR 同票） |
| 4 ASR（本地） | **15** | OPEN | **6** | 0 | streaming ASR via sherpa＋engine slot mutex＋CER harness |
| | **14** | DONE | — | — | 模型分发（六枚 purpose 全部可下载可校验） |
| | **27** | OPEN | **5** | 0 | 标点复原（ASR 尾巴） |
| 4 ASR（云端 C9） | **61** | OPEN | **6** | 0 | 级联 ASR/TTS from configured models |
| 5 送进模型 | **121** | DONE | — | — | 交还链进二进制＋装配可达性门 |
| | **12** | OPEN | 未量 | — | `wisp run` 文本路径（＝今天唯一的进模型路径） |
| | — | | | | **⚠「ASR 终稿 → `Loop.Run(ctx, input)`」这一枚没有票**（§4） |
| 6 TTS／出声 | **26** | OPEN | **7** | 0 | matcha 引擎＋serial slot＋Speaking 管线＋P7 门＋`DEFERRED(playback)` |
| 7 Path T 半双工 | **247** | OPEN | 10 | 0 | 门所依附的采集腿（门本体已由 13 交付） |
| 7 Path C 全双工 | **59** | OPEN | **7** | 0 | AEC spike（许可/绑定/回声质量/CPU/误触发阈） |
| | **60** | OPEN | **7** | 0 | C32 RealtimeEngine（云端 S2S，gated） |
| | **28** | OPEN | **6** | 0 | C31 SessionScope Warm/Conversation 保活 |
| 全链（Sink 执行器） | **246** | OPEN | 未量 | — | 只裁了"常驻腿不持有 Machine"这件事的**否决通道**一半 |
| | **07** | DONE | — | — | 逐字声明 "SessionScope creation etc. are no-ops here"＝当时的**设计内** no-op |
| | — | | | | **⚠「给常驻进程装一枚带 `Options.Sink` 的 Machine＋41 张副作用名的执行器」没有票**（§4） |

**语音链路上 OPEN 票的未勾 AC 合计（只算本腿量到的 8 枚）**：247=10 ＋ 15=6 ＋ 26=7 ＋ 41=6 ＋ 61=6 ＋ 59=7 ＋ 60=7 ＋ 28=6 ＋ 27=5 ＋ 165=5 ⇒ 10 枚票共 **65** 枚未勾 AC／0 已勾。票池总 open 数＝**166**（`ls | grep -v -- '-done' | wc -l`）。
另有一枚**契约级、立而不派**的相关票：**165**「先出方案、你点头再动手」那一档（5 枚未勾，标题逐字标 **契约级·立而不派**）——它约束的是语音入口之后的行动门，不是链上的跳。

---

## 4. 判不动（缺什么料，老实列）

1. **`wisp run` 的 `task` 到底从 argv 还是 stdin 取**——我只读到 `run.go:1099 bg := loop.RunAsync(ctx, task)`，没有回溯 `task` 的赋值行（怕整段读进上下文超预算）。⇒ 影响"跳 5 今天这条文本入口是不是已经能吃非命令行来源"的判定；需要一枚 `cmd/wisp/run.go` 里 `task :=` 的行号级定位。
2. **VAD 的包归属两说**（`internal/audio/doc.go:3` vs `statemachine/table.go:52` ＋ `models/manifest.go:94`）——我判不动哪一说是权威的，因为两枚文件都是产码注释、票池里 15 号票的正文我没有逐字读完（只读了标题行）。⇒ 这一枚不定，第三跳会写两遍或写进被禁的包。
3. **「ASR 终稿 → `Loop.Run`」与「常驻 `Options.Sink` 执行器」这两枚缺口，是真的没票，还是被我扫漏了**——我用的口径是：票池 260 枚文件名＋标题级 grep（`EffectSink`/`执行器`/`statemachine.New`/`sink`/`handOff`/`ASR`/`transcri`），命中者已列进 §3。**没有**逐字读完 166 枚 open 票的正文。⇒ 要定案需要一枚把 open 票正文全扫一遍的活（本腿受输出预算限制未做）。
4. **`272-r2` 正在改的 `cmd/wisp` 测试文件会不会新增音频/语音 import**——我只在 `82a10da4`＋脏树上量了一次"audio importer=0"。它此刻可能已经变。⇒ 任何下游引用这一读数前，**重跑 §1 跳 1 那两把尺**。
5. **票 12/246/68 的未勾 AC 枚数我没有量**（§3 标"未量"）。它们不是语音主票，本腿不去数，避免与别的腿的取数窗口打架。
6. **真机侧一律未证**：麦克风是否可用、`wisp doctor` 现在的 PASS/FAIL、`Armed` 态是否可达——本腿没跑任何产码，也没开任何窗口，没加 `-tags winlive`。票 247 的 `A485/P2` 已经裁过"Armed 今天不可达（`EvKwsEnabled` 零发射者）"，我这把尺独立复现了"零发射者"，但**没有**独立复现"Armed 不可达"那一半（它依赖转移表可达性分析，属别的腿的射程）。
7. **D47 原文我没有逐字读**（`docs/PLAN.md` 只按 §1 引了 `internal/audio` 内对 D47 的复述）。⇒ 本件里"半双工只约束干活路径"这一句的**权威出处是 PLAN.md 的 D47 那一节**，不是我的复述；要写进票面请回原文。
