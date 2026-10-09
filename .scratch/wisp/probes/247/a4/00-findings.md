# 247-a4 只读普查：采集栈落点量清（供 247-r1 写腿）

锚 sha `d106fc25`（`git log --oneline -1`）· 分支 `dev` · 时刻 `2026-10-09 11:45 +0800`（`date`）
起手快照 `git status --porcelain -- cmd internal frontend` → **空**（这三棵目录本刻无在飞字节；`.scratch` 有别人在飞的件，一律不碰）。
读法：全部走对象层 `git show HEAD:<path>`；未读 `frontend/**`／`design/**`；**未跑任何 `go` 子命令**。
裁决形状（票面 AC#0 编排者裁决，已读）＝**甲＋乙-2**：新增依赖边恰好 1 枚 `cmd/wisp→internal/audio`；跨 `internal/ball` 只过一枚 `float32`；协程挂现成 `observe.Default.Spawn`（`goroutine.go:262`）；形①只写日志/判语、不开麦。

## Q1 `internal/audio` 对外面 + 电平从哪出来

包导出名（逐枚，来自 `git show HEAD`）：
- `audio.go`：`AudioSource`（C8 缝：`Start(ctx, buf chan<- []byte) error` + `Stop() error`，只吐 PCM **字节**、不吐电平）、`Stater`（`Stats() Stats`）、常量 `TargetRate=16000`/`FrameSamples=512`/`FrameBytes`/`FrameDuration`/`BoundedWindow`、`BoundedFrameCapacity()`、`NewBoundedFrames() chan []byte`、`Path`、`Stats{}`+`Dropped()`、**`SpawnCapture(registry *observe.Registry, run func(ctx)) *observe.Handle`（`audio.go:180-182`）**——它是"audio-capture"协程的**接缝**，今天**非测试调用者＝0**（现尺见下）。
- `level.go`（纯函数、无状态、无 live getter）：`LevelFullScale`/`LevelLSB`/`FullScaleSquareLevel`/`SineLevelTolerance`/`MinLevel`/`MaxLevel`、**`LevelOfSamples([]int16) float64`**、**`FrameLevel([]byte) (float32, error)`**、`DecodeFrame([]byte)([]int16,error)`。
- `device.go`：`DeviceDescriptor`、`DeviceWatcher`、`DeviceError(hr,dev,what) *observe.Error`。
- `gate.go`：`HalfDuplexGate`（也是 AudioSource）、`NewHalfDuplexGate(inner AudioSource, path, opts...)`、选项 `WithStartMuted(bool)`（`:57`）、**`WithGateEvents(fn func(GateEvent))`（`:63`）＝仓里现成的"注入回调"先例**。
- `wasapimic_windows.go`：`WASAPIMicrophone`+`NewWASAPIMicrophone()`；owner 写死在 `:83`。`wavinjector.go`：`WavInjector`+`EncodeFrame([]int16)[]byte`；owner 写死 `:84`。`wasapi_windows.go`：`Drain()([]int16,error)`/`WaitEvent`（未导出类型上的方法，包内）。

电平从哪出：采集循环 `run`（`wasapimic_windows.go:149`）里 `stream.Drain()`→重采样→`m.meter.push(buf, EncodeFrame(frame))`（约 `:228`），**只把 PCM 字节推进消费者 chan，不产出电平**。全仓**没有任何"返回当前电平"的访问器**。

现有 API 够不够：算电平够用（`FrameLevel`），但**没有推送式电平出口**——外部拿不到"当前电平"，除非自己起循环读 chan 再调 `FrameLevel`（那会把 `[]byte`/samples 拖到消费侧，破 AC#3 只过 float32）。
最小新增形状（只给形状，⛔不写码，按 AC#0 乙-2/P8＝电平在**已有 audio-capture 协程内**算完再投递）：仿 `WithGateEvents` 的注入式回调，在采集侧新增一枚导出选项/入参，类型 **`func(float32)`**，内部喂 `FrameLevel`；并把 `wasapimic_windows.go:83`/`wavinjector.go:84` 写死的 `observe.Default.Spawn` 改吃传入 registry（即走现成但零调用的 `SpawnCapture`）。⇒ `internal/audio` 内**不新增包级依赖边**，唯一新边仍是 `cmd/wisp→internal/audio`。
尺＝`git show HEAD:internal/audio/{audio,level,gate,wasapimic_windows,wavinjector}.go | grep -nE '^(func|type|const) '`；`SpawnCapture` 非测试调用者尺＝`git grep -nE '\.SpawnCapture\(|SpawnCapture\(' HEAD -- '*.go' ':!**_test.go'` → 读数仅 1 行＝定义本身（`audio.go:180`），**调用者 0**。

## Q2 球侧入口 `SetAudioLevel`

签名逐字：`func (b *Ball) SetAudioLevel(level float32)`（**`internal/ball/liquid_windows.go:42`**）。
⚠ 任务写的 `internal/ball/table.go:97` **行/文件名皆漂**：HEAD 里 `internal/ball/` **无 `table.go`**（`git ls-tree --name-only HEAD internal/ball/` 无此名）。
生产调用者枚数（调用形状尺）＝**2**，全在 `cmd/balldebug/main.go`：`:418 b.SetAudioLevel(v)`、`:477 b.SetAudioLevel(0)`；喂的是**合成**包络（同文件 `:396` 注释逐字 "feedLevels pushes a syllabified synthetic envelope"）⇒ **真麦克风生产者＝0**（合票面前提）。词面尺（`SetAudioLevel` 非测试全串）另命中定义行 + 若干注释/doc（ball_windows.go:108、liquid.go:13、liquid_windows.go:14/24/35、level.go:112、balldebug 注释）；两尺**分歧只在注释/文档噪声**，对"真实调用点＝2、均为 balldebug 合成"一致。
附带（AC#10）：`SetAudioLevel` 在 `prototypeVisuals` 关时直接 return（`liquid_windows.go:51-54`），默认档接上也不"呼吸"——取证须用读数型证据，⛔ 不许在常驻腿 `EnablePrototypeVisuals(true)`。
尺＝`git grep -nE '[A-Za-z0-9_]\.SetAudioLevel\(' HEAD -- '*.go' ':!**_test.go'`（调用形状）与 `git grep -n 'SetAudioLevel' HEAD -- '*.go' ':!**_test.go'`（词面）并排；均**带 `*.go`+排除测试的 pathspec**，非"全仓"。

## Q3 常驻腿挂载点 + 协程 owner + 配置读取处

`cmd/wisp/resident_windows.go` 今天跑的循环/装配（非测试）：`installLogSink`（`:66`）、`rt.Shutdown(false)`（`:76`）、`newResidentApprovalWithConfig(rt.Layout.DataDir)`（`:132`）、`startResidentPanel(rt.Registry, rp)`（`:157`）、`startResidentBall(rt.Registry, …)`（`:217`）、`rt.RegisterShutdownHook(proc.StepCancelTasks, …)`（`:240`）、`rt.RegisteredShutdownSteps()`（`:264`）、阻塞在 `rt.RunEventLoop()`（`:285`）。协程经 **`rt.Registry`（一枚 `*observe.Registry`）**由 `startResidentPanel/Ball` 内部 Spawn，本文件不直接 `observe.Default`。
owner 逐字定义：`internal/observe/goroutine.go:262` `func (r *Registry) Spawn(name, owner string, root *Root, fn func(context.Context)) *Handle`；包级 var `Default = NewRegistry()`（`goroutine.go:230`）。**既有用法范例（逐字引一处，就在 audio 侧）**：`wasapimic_windows.go:83`
`	observe.Default.Spawn("audio-capture", "audio", nil, func(c context.Context) { … m.run(loopCtx, buf, started) })`
⇒ 采集协程挂这枚现成 owner 即可，⛔ 不裸 `go func(`（ban#1），⛔ 不新造 D38(e) 第 11 步；停机 join 走现成 `proc.StepStopAudio` 钩子位。
配置键与默认档（AC#4：默认值+读取处一字不改，此处只报现量）：
- 默认档逐字：`voice.enabled`＝`Enabled bool \`toml:"enabled" default:"true"\``（**`internal/config/schema.go:245`**，true）；`wake_word.enabled`＝`Enabled bool … default:"false"`（**`schema.go:197`**，false）；第三枚真门 `audio.mic_muted_default`＝`default:"true"`（**`schema.go:277`**）。
- 今天谁读：`Voice.Enabled`/`WakeWord.Enabled` 的**唯一非测试读取处**＝`internal/config/manager.go:380`（reload diff）/`:394`/`:395`（reload apply），属 **reload 分层**（`tiers.go:53-54` 都标 "reload"），**不是开麦决策**；`MicMutedDefault` **无任何产码读取处**，只有 `gate.go:56` 一句 doc 注释点名它。⇒ 证实 AC#0 P1：voice.enabled／mic_muted_default 今天都没接到任何开麦路径；接线须把 `MicMutedDefault` 传进 `WithStartMuted`、并把 `voice.enabled=false` 判为"不构造采集器"。
- 新边落点：cmd/wisp 已 import `internal/config` 与 `internal/observe`，缺的正是 `cmd/wisp→internal/audio`（现量：audio 非测试 importer＝0，票面 现量1）。
尺＝`git show HEAD:cmd/wisp/resident_windows.go | grep -nE 'Spawn|observe\.|rt\.|Registry|RegisterShutdownHook|RunEventLoop'`；`git grep -nE 'voice\.enabled|wake_word\.enabled|mic_muted_default|Voice |WakeWord |MicMutedDefault' HEAD -- 'internal/config/*.go' ':!**_test.go'`；`git grep -nE 'Voice\.Enabled|WakeWord\.Enabled|MicMutedDefault' HEAD -- '*.go' ':!**_test.go'`。

## Q4 降级三形今天分不分得开 + "响亮"先例

`DeviceError(hr, dev, what) *observe.Error`（`device.go:94`）按**具名 HRESULT 常量**分岔，三形今天**分得开**：
- 无权限：`hrEAccessDenied=0x80070005` → `privacyGuidance` 分支（`device.go:86`）。
- 设备被占（独占模式）：`hrAUDCLNTDeviceInUse=0x8889000A` → `inUseGuidance` 分支（`device.go:89`）。
- 无设备/枚举失败：`mmdevice_windows.go:86` `defaultEndpoint` 里 `GetDefaultAudioEndpoint` 失败 → `observe.New(ClassAudioDevice,…)`，再经 `wasapimic_windows.go:261 openCurrent` `observe.Wrap(ClassAudioDevice,…,"capture device enumeration failed")`；运行中被拔＝`hrAUDCLNTDeviceInv=0x88890004` → "device invalidated (unplugged or disabled)" 分支，`isDeviceInvalidated` 靠 `ProviderCode` 识别。
形状：**没有具名哨兵错误值**（无 `ErrNoDevice`/`ErrAccessDenied` 导出）；`AudioSource.Start` 返回的是裸 `error`，**具体类型恒为 `*observe.Error`**（带 `Class==observe.ClassAudioDevice` + `ProviderCode`＝HRESULT hex + guidance 文案）。⇒ "分类＋可见"在第 0 跳就带上，**断点在没人消费它**（importer＝0），正合 AC#0 形①"只把分类＋guidance 响亮打印并落 sink"。
"响亮"现成先例（具名文件:行，本腿现量）：`cmd/wisp/resident_windows.go:68` 注释逐字 "Loud, and it does not stop the app: a log directory that will not open must not become a way to keep Wisp from starting"（`installLogSink` 出错→记日志后回落 stderr、启动继续）；`internal/audio/wasapimic_windows.go:272-276 reopenOnce` 成功 `slog.Info("audio capture reopened…")`、失败具名上一枚设备名后 `postErr`+`meter.fail`。
尺＝`git show HEAD:internal/audio/device.go | sed -n '68,135p'`＋`…wasapimic_windows.go | sed -n '258,295p'`＋`…mmdevice_windows.go | sed -n '76,110p'`；先例尺＝`git grep -nE 'does not stop the app|Loud' HEAD -- 'cmd/wisp/resident*.go' ':!**_test.go'`。

## 没答的 / 量不到的 / 与票面或台账冲突（具名）

- **量不到（须真设备+真跑，归写腿/AC#2）**：真机三形的逐字回串（`0x8889000A`／`DEVICE_INVALIDATED`／`NOT_INIT` 到底哪枚出现）；投递路径进 pinned 线程后的周期读数复量（A485-P8）。本腿⛔未跑任何 `go`。
- **本腿没重开的**：`resident_ball_windows.go` 有 err 也 `return rb` 那条**精确行号**（票面 a1 记 `:139-145`，我没复跑＝标"未复跑"）；`observe.Default.Spawn("watchdog","config",…)` 的**产码行**（我只在 `.scratch/.../231/r1/mutation` 探针拷贝里见过，未开 `config_reload.go` 本体）。
- **与票面/任务冲突（具名，只报不改）**：
  1. 任务 Q2 把 `SetAudioLevel` 写作 `internal/ball/table.go:97`；实测在 `internal/ball/liquid_windows.go:42`。票面 AC#0 裁决又写 `EvAudioDeviceLost` 的边在 `internal/ball/table.go:97`，实测该边在 **`internal/statemachine/table.go:97`**（`D43:14, From: StateListening, To: StateError`；事件名 `internal/statemachine/events.go:30`）。`internal/ball/table.go` 于 HEAD **不存在**。
  2. 存在性陷阱：`git ls-tree HEAD <不存在路径>` 返回 exit-0 且**空输出**，naive `&& echo FOUND` 会**假报 FOUND**；我改用目录清单 + `EvAudioDeviceLost` grep 复核才定案（提醒写腿别拿 ls-tree 判存在）。
  3. 锚点漂移：本次全部读数钉在 `d106fc25`；票面引用的 `bb37fac2`/`4e66817` 等是更早的锚，行号可能已搬。
  4. 起手 `cmd internal frontend` 工作树**本刻干净**，与"写面被另一枚腿独占"的表述不完全一致——在飞字节此刻都在 `.scratch`；我仍一律走 HEAD 读，未依赖工作树。
- **未越界**：未碰 `internal/speech`／`scripts/spike`／`frontend/**`／`design/**`／`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`；未改任何产码/测试/票面/台账；写点唯一＝本文件。
