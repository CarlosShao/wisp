# 247-a2 — 采集栈接球「四问代价」复量普查件（只读腿）

- 腿：`247-a2`（只读普查，`247-a1` 之后的现树复量）· 票：`.scratch/wisp/issues/247-the-capture-stack-has-zero-importers-so-no-real-microphone-level-ever-reaches-the-ball.md` AC#0
- 本节文件：`§0` 锚点与口径 · `§1` 落点甲乙与依赖边名册 · `§2` 管道逐跳三态 · `§3` 两把开关与默认档自证 · `§4` 真机读数怎么量 · `§5` 量不到的格子 · `§6` 我推翻前人哪几句
- 本件不改任何产码；写点只有这一枚文件。`frontend/**`／`design/**` 未读未引；零 Go 命令（详见 §0）。

## §0 锚点与口径

- 起量 HEAD：`47b38765`（分支 `dev`，`git log --oneline -1` 现跑）；本件自身的骨架 commit＝`fb8e1ff5`（只动 `.scratch/wisp/probes/247/a2/**` 两枚路径，`git show --stat` 证 2 文件 24 增、零产码）。
- 量数时刻：`2026-10-03 09:49–09:57 +0800`（每节末尾 `date` 现插，不手打）。
- module 路径：`github.com/CarlosShao/wisp`（`go.mod:1` 现读）⇒ 所有 import 尺都认这条前缀。
- **本程零 Go 命令自证**：本腿只用过 Read/Grep/Glob 与 `grep`/`sed`/`awk`/`git log`/`git show`/`ls`/`head`；⛔ 未跑 `go build`／`go vet`／`go test`／`go run`／`go list`、未起 sherpa、未碰真机音频。同机三枚写腿（`cmd/wisp`／`internal/panel`／`internal/tools`）在飞，任何编译都会互洗读数，故本件对"包级依赖边"的判断是**结构读**（逐枚 Read `import` 块 + `grep` 路径），并交叉 `247-a1` 的 `go list` 名册（那三枚名册＝`deps-*.txt`，一律标〔待验〕、非本腿现跑）。
- 工作树口径（诚实边界）：现树里带着同机在飞腿的**未提交**改动——`cmd/wisp/config_receipt_255_test.go`（M，测试件，本件未引）、`design/**`（D，本件未读未引，见 §硬约束）、以及海量 `.scratch/**` 未跟踪件。本件凡"现读 file:line"＝上述时刻**文件系统**里那枚文件的行号，不是 HEAD blob；引用的产码文件都不是在飞腿正在改的那枚，故不受其未提交态污染。
- 全称否定一律附尺＋正控：本件出现"零 importer／无产码读取处／没有出口"这类句子时，同处给那把尺的字面命令、`grep_exit` 读数，以及一处**已知命中**的正控（证尺不瞎）。含多分支或中文的 grep 若回空，先复核尺再入账（`247-a1` 的 E2 就是栽在这上面）。

## §1 落点甲乙与依赖边名册

### 1.1 现量：`internal/audio` 的非测试 importer ＝ 0（AC#1 的"接前读数"）

- 尺：`grep -rl "CarlosShao/wisp/internal/audio" --include=*.go . | grep -v "/internal/audio/"` ⇒ 空，`grep_exit=1`。连测试枚一起数、整棵树搜那枚 import 路径，**零命中**（不只是非测试零）。
- 尺（点名 cmd/wisp）：`grep -rln "wisp/internal/audio" --include=*.go cmd/wisp | grep -v _test` ⇒ 空，`cmdwisp_audio_exit=1`。
- 正控（证两把尺都不瞎）：① `grep -rln "^package audio" ./internal/audio` ⇒ `audio.go/device.go/doc.go/gate.go/level.go/…` 命中；② 同形尺换成 `wisp/internal/proc` 搜 cmd/wisp ⇒ 命中 `doctor.go`／`resident_task_source_windows.go`／`resident_windows.go`。⇒ "cmd/wisp 不 import audio"是真的没有，不是尺瞎。
- 结论：**票 33（面板宿主进常驻进程，`13acad46`）／票 246（任务源）／票 248（设置落盘）都并进来之后，AC#1 的"接前读数"此刻仍是 0**，可直接用。`cmd/wisp/resident_ball_windows.go:211` 的日志正文字面写着这条腿"**no microphone**"，与读数一致。

### 1.2 依赖边名册（结构读，非 go list）

- `internal/audio` 的**唯一 wisp 包级依赖是 `internal/observe`**：`grep -rhn "CarlosShao/wisp/internal" internal/audio | grep -v _test | sort -u` ⇒ 只有 `internal/observe`（散在 audio.go:9 等枚文件）；外部只有 `golang.org/x/sys/windows`（`mmdevice_windows.go`／`wasapi_windows.go`）。
- `internal/ball` **不 import** `internal/audio`（尺 `grep -rn "wisp/internal/audio" internal/ball` ⇒ `e=1`）⇒ 今天 audio↔ball 两头都没边。
- `cmd/wisp` 已 import `internal/observe`（海量枚）、`internal/ball`（`resident_ball_windows.go`／`resident_approval_windows.go`）、`internal/proc`。⇒ cmd/wisp 的传递名册里 audio 除 observe 外的东西（stdlib＋x/sys/windows，proc 本就用）**早已在场**。

### 1.3 形甲（采集协程挂常驻腿 `cmd/wisp`，与球同进程）——要动的面与代价

装配次序现读（`cmd/wisp/resident_windows.go` `runResident`）：`proc.Boot`(:39) → `installLogSink`(:63) → `defer rt.Shutdown`(:72) → `newResidentApproval`(:123) → `newResidentPanelManager`/`startResidentPanel`(:142/:148) → **`rb := startResidentBall(rt.Registry, ...)`(:163)** → `ra.bindBallHost(rb)`(:174) → `RegisterShutdownHook(StepCancelTasks,...)`(:186) → `startResidentTaskSource(rt, ra)`(:206) → 开机报告(:210-216) → 事件循环(:231)。

- 新增包级边：**恰好 1 枚** `cmd/wisp → internal/audio`。不新开 `audio→ball`，也不动 audio 的 observe 依赖（observe 已在 cmd/wisp 名册里）。这是结构推读；`247-a1` 的 `comm -23 deps-audio deps-cmdwisp` 只出这一名（那枚名册〔待验〕，非本腿现跑）。
- 要新增/改的文件（甲形）：
  1. **新增** `cmd/wisp/resident_audio_windows.go`（`//go:build windows`）：`NewWASAPIMicrophone()`(`wasapimic_windows.go:52`) → `NewHalfDuplexGate(inner, PathT, WithStartMuted(c.Audio.MicMutedDefault))`(`gate.go:84`/`:57`) → `NewBoundedFrames()`(`audio.go:65`) → `mic.Start(ctx, buf)` → 消费侧读 buf → `FrameLevel(frame)`(`level.go:118`, 返回 `float32`) → `rb.b.SetAudioLevel(v)`(`internal/ball/liquid_windows.go:42`)。球句柄就在同包 `rb.b`（`resident_ball_windows.go:106` `b *ball.Ball`，包 main 可直接取）。
  2. `resident_windows.go`：在 :163 之后接这段，并 `rt.RegisterShutdownHook(proc.StepStopAudio, ...)`（见 1.5 join）。
- **一笔藏不住、`247-a1` 已预警、现树仍在的代价**：`runResident` 手里那枚 `rt` 是 `proc.Runtime`，字段只有 `Env/Layout/Job/Instance/Registry/StartedAt/shutdownHooks`（`boot_windows.go:28-42`），**没有 cfg**；`agentRuntime.cfg`（`run.go:268`）在 `startResidentTaskSource`→run 那条子装配里才取（`run.go:424 rt.cfg=cfg`、`run.go:991 cfg := rt.cfg`）。⇒ 甲形要在 runResident 层读 `audio.mic_muted_default`/`voice.enabled`，得**自己 `config.LoadFile(filepath.Join(rt.Layout.DataDir, configFileName), nil)`**（`configFileName`＝`cmd/wisp/secret.go:63` "config.toml"，同一枚调用形状现成在 `providers.go:98`/`models.go:184`）。这不是装饰，是甲形的第一笔。
- 不需要非 Windows 桩：`cmd/wisp/resident_other.go` 已 `os.Exit(2)` 拒起常驻；但消费体若调 `mic.Err()/ThreadID()/Endpoints()`（这三枚只在 `wasapimic_windows.go:119/130/137`，`wasapi_other.go` 占位体没给）⇒ 那枚消费文件**必须** windows-tagged。

### 1.4 形乙（`internal/audio` 自起协程并推到球）——两支，代价不对称

- 乙-1 直推（audio 里持/调 `*ball.Ball`）：新增 `internal/audio → internal/ball`，**今天不存在**（1.2 已证）。后果两笔：① 采集层从此吃下渲染层＋D43 状态机（ball 的直接依赖含 observe＋statemachine）；② **破 `internal/audio/wasapi_other.go:6-7` 的明文承诺**——"the package must still compile on other GOOS values so accidental cross-packages references fail loudly, not silently"，而 `internal/ball` 是 windows-only。且它**并不省 cmd/wisp**：`mic.Start` 仍得有人调、第 4 步钩子仍得有人注册 ⇒ cmd/wisp→audio 这条边照开，乙-1 实际 ≥2 枚新边＋一条契约破口。
- 乙-2 注入（audio 自起、出口是传入的 `func(float32)` 或 audio 侧 sink 接口，cmd/wisp 把 sink 接到球）：**边数与甲完全相同**（还是 cmd/wisp→audio 一枚，不碰 ball 的任何依赖）。乙-2 相对甲只是把"那枚读帧的协程"的所有权从 cmd 侧搬进 audio 侧——于是 §1.6 的"名册无槽"约束从 cmd 侧搬到 audio 侧，并未消失。**所以乙不比甲省 cmd/wisp 的改动**，只是换个包承担同一枚第二协程难题。

### 1.5 协程 owner 名册（AC#5 禁裸 `go func(`；两形都走现成 owner）

`observe.Root`／`Registry.Spawn` 定义：`internal/observe/goroutine.go:88 type Root`、`:98 NewRoot`、`:103 NewRootFrom`、`:262 func (r *Registry) Spawn(name, owner string, root *Root, fn)`；`Spawn` 是全仓唯一被许可的 `go` 语句，名册外名字在 `:271 slog.Warn("goroutine outside the D38 roster (leak symptom)")` 现形，`RosterReport`(:405) 的 `:422 ResidentOverBaseline = rep.Resident > ResidentBaseline(6)` 判泄漏。

现成可挂 owner（两形都用得上，两侧各至少一枚 file:line）：

- **audio 侧既有 owner**：`internal/audio/wasapimic_windows.go:83 observe.Default.Spawn("audio-capture","audio",nil,...)`（真麦路径）＋ `internal/audio/wavinjector.go:84`（WAV 注入路径，**同名** `audio-capture`）——两枚都**写死 `observe.Default`**。收注册表的那枚接缝 `internal/audio/audio.go:180 SpawnCapture(registry *observe.Registry, run)` 产码**零调用者**（尺：`grep -rn SpawnCapture --include=*.go .` 排 `.scratch`/`_test` ⇒ 只有 audio.go:176/180 定义）。
- **cmd/wisp 侧既有 owner**（现读）：`cmd/wisp/approval_always.go:99 observe.Default.Spawn("approval-waiter",...)`；`cmd/wisp/config_reload.go:119 rt.reloadHandle = observe.Default.Spawn("watchdog","config",rt.reloadRoot,...)`；**票 33 新落**的 `cmd/wisp/panel_resident_windows.go:153 rp.thread = reg.Spawn(panelSTAName, panelSTAOwner, root, rp.loop)`（用**传入的 registry**，即 `rt.Registry`）；球本体经 `resident_ball_windows.go:172 Registry: reg` 起它的 `ui-sta` 线程。
- 裸 `go func(`：产码里 `grep -rn "^\s*go func(" --include=*.go internal cmd tools | grep -v _test` ⇒ 只命中 `cmd/wisp/testdata/esclistener/main.go:252`（testdata fixture，非 app 路径）⇒ 正控说明尺能抓到裸 `go func(`，"产码零裸协程"是真的。

注册表同一性（本腿新读，纠 `247-a1` 的担忧见 §6）：`runResident` 调 `proc.Boot(env)` **不带** `WithRegistry` ⇒ `rt.Registry == observe.Default`（`boot_windows.go:79` 那支赋值；`WithRegistry` 的产码调用者除注释外为零）。⇒ **今天球的 `ui-sta`（走 `rt.Registry`）与 wasapimic 写死的 `observe.Default` 落进同一枚注册表，没有"两枚并存"的假绿**；那枚风险是**潜伏**的，一旦有 caller 传 `WithRegistry`（今只有测试路径）才现形。`A485`-P4 仍要求把采集协程改成吃传入 registry（用现成 `SpawnCapture` 接缝），本腿认可其**方向**，但把"正在假绿"降级为"潜伏、当前不成立"。

### 1.6 D38(e) 十步一字不动下，退出谁 join 它

- 接缝现成、且 `StepStopAudio` 那格**今天无人注册**：`internal/proc/shutdown.go:40 const StepStopAudio`、`:163 run(StepStopAudio, hooks.StopAudio, 0)`；`shutdown_hooks.go:51` 列进可挂钩名册、`:118 case StepStopAudio`。`RegisterShutdownHook` 的产码注册者只有 `resident_windows.go:186`（StepCancelTasks）与 `resident_task_source_windows.go:311`（循环注册），**没有一处是 StepStopAudio** ⇒ 那正是本票要填的格。
- join 动作本就在源里：`wasapimic_windows.go:94 Stop()` → `cancel()` 后 `select { case <-done: case <-time.After(tm.Remaining()): }`，`tm = observe.NewTimeout(2*time.Second)`（`:106`，单调、不是墙钟差）。⇒ 钩子体只需 `gate.Stop()`（`gate.go:115`）/`mic.Stop()`，**不需新造第 11 步**。
- ⛔ 诚实边界：本腿**没跑到**那枚 `case StepStopAudio`（禁 build/test），`shutdown_hooks.go:118` 是读到的分支不是跑到的分支（与 `247-a1` E5 同源）。交 `247-v1` 钉一条：注册后 `StepRecord.Name == "stop-audio"` 且记录不是 skipped。

### 1.7 名册无第二枚电平协程的槽（硬约束，甲乙都撞）

`ResidentNames`（`observe/goroutine.go:43-44`）固定 6 名含 `audio-capture`，`ResidentBaseline=6`（`:40`），`OnDemandNames` 只有 `kws-infer`（`:48`），`PerTaskPrefixes`（`:52`）＝`agent-task-`/`tool-exec-`/`approval-waiter`。⇒ 电平消费者**不能新起一枚有名常驻协程**。audio 包作者已把这行写进代码：`level.go:86-87` 逐字 "a level reader is called from the consumer's own thread, and D38b leaves no spare resident slot for a level-only goroutine"。⇒ 落点无论甲乙，电平都必须在**已有的 `audio-capture` 协程内算完再投递**（甲形：投递到 `rb.b.SetAudioLevel`，它 `:46 b.sta.PostTask` 把渲染甩回 UI 线程，采集侧不碰渲染）。代价：把一次 `FrameLevel`＋一枚跨包投递带进 pinned 循环——而 `wasapimic_windows.go:146-148` 的循环自述只列 open/wait/drain/resample/bounded push/hotplug，**不含算电平**（现读 `:186-234` 循环体确实只 `EncodeFrame`+`meter.push`，无 FrameLevel）。⚠ 这条自述与实现的差要 `247-v1` 复量 pinned 周期后再定。

_§1 量数时刻：`2026-10-03 10:0x +0800`。_

## §2 管道逐跳三态

真麦克风 PCM → 球上那个 `float32` → 液态，一共 5 跳。逐跳三态＝**今天存在**／**差一行赋值**／**一整块没写**。尺：`grep -rn "NewWASAPIMicrophone\|NewHalfDuplexGate\|NewBoundedFrames" --include=*.go cmd internal tools | grep -v "internal/audio/" | grep -v _test` ⇒ `h1_exit=1`（正控：同三名去掉排段后在 `internal/audio/*_test.go` 与 `wavinjector.go`/`wasapimic_windows.go` 大量命中，尺不瞎）。

| 跳 | 这一跳要做什么 | 三态 | 现读凭据 |
|---|---|---|---|
| **H1** | 有人**构造**采集栈：`NewWASAPIMicrophone()`＋`NewHalfDuplexGate(...)`＋`NewBoundedFrames()`，并 `mic.Start(ctx, buf)` | **一整块没写** | 三名产码零构造者（`h1_exit=1`）；`NewWASAPIMicrophone` `wasapimic_windows.go:52`、`NewHalfDuplexGate` `gate.go:84`、`NewBoundedFrames` `audio.go:65` 全在库内、只被测试与包内占位用 |
| **H2** | 有人**读那枚 `buf chan<- []byte`**（消费者侧循环） | **一整块没写** | `buf` 由源写（`meter.push` `audio.go:115`、被 `wasapimic_windows.go:231` 调）；全仓**没有一枚产码 reader**（audio 零 importer，见 §1.1）；pinned 循环体 `wasapimic_windows.go:186-234` 只 open/wait/drain/`EncodeFrame`+push，**不算电平**（尺 `grep -n "FrameLevel\|float32\|level" wasapimic_windows.go` ⇒ `e=1`） |
| **H3** | `[]byte` → 一枚 `float32` 电平 | **函数存在，调用者零** | `FrameLevel(frame []byte) (float32, error)` `level.go:118`；`LevelOfSamples` `:96`（内用 `math.Sqrt` `:105`）；`DecodeFrame` `:137`、`EncodeFrame` `wavinjector.go:138`；分母 `LevelFullScale = 32768.0` `level.go:54`。⇒ 生产者半**已写、已测**（票 241），只差没人调 |
| **H4** | 把 `float32` 交给球：`b.SetAudioLevel(v)` | **方法存在，产码调用者只有 balldebug** | 定义 `liquid_windows.go:42`（`:46 b.sta.PostTask(func(){ b.applyLevel(level) })`）；产码调用者＝`cmd/balldebug/main.go:418`/`:477`；**`cmd/wisp` 里零 SetAudioLevel**（尺 §1.1 同域，`cmdwisp_audio_exit=1` 的另一面） |
| **H5** | 球把电平变成看得见的液态 | **一道独立闸，默认关，非本票地界** | `applyLevel` `liquid_windows.go:51` 第一行 `:52 if !prototypeVisuals { return }`（**在把 `b.liqRaw = level` 之前**，`:55`）；`var prototypeVisuals bool` `statevisual.go:102` 默认关，开关 `EnablePrototypeVisuals` `:106`，产码只有 `cmd/balldebug/main.go:122` 打开；即便打开，`:56 if !liquidDriven(b.curState)` 仍把非会话态的电平丢弃（`liquidDriven` 六态 `liquid.go:62-70`，不含常驻腿能到达的 Sleeping/Confirming） |

### 关键纠偏（这一节最重要）

- **"只差一行赋值"是假象**：H4 那行 `SetAudioLevel` 调用看起来最省，但它**挂在 H1＋H2 之后**——今天 H1、H2 **一整块都没写**。若有人按"补一行"去做，只会补出一枚**喂着假数的死赋值**（要么塞 balldebug 那种合成包络＝票面现量 2 要剪掉的东西，要么写一句编译得过、永远不被调的调用）。⇒ 本腿明确：真正欠的是 H1（构造＋Start＋config＋shutdown 钩子）与 H2（读 buf 的消费者协程，且撞 §1.7 名册无槽），不是 H4。
- **`09-30` 傍晚 241-r1 之后，"全仓零 Sqrt" 是历史读数**：本腿现读 `level.go:105` 就有 `math.Sqrt`，`FrameLevel`/`DecodeFrame`/`EncodeFrame` 与 32768 分母都落地了。⇒ H3 的**生产者半已闭合**，别再引"PCM→电平这一整段都没定义"。欠的只在 H1/H2/H4 接线。
- **证据落点受 H5 牵连**：`applyLevel` 在 `:52` 就 return，连 `b.liqRaw` 都不写。⇒ 想拿"电平进了球"当证据，**不能靠读球内部**（默认档下球里啥也没留），只能在接线侧（H4 之前那一步）把 `float32` 自己打出来。这直接决定 §4 的取证形态。

_§2 量数时刻：`2026-10-03 10:0x +0800`。_

## §3 两把开关与默认档自证

### 3.1 三把开关：字段、默认值、读取处（默认值是代码事实，不是文档口径）

| 语义 | Go 字段 | toml 路径 | 默认值（现读） | 默认值怎么落地 |
|---|---|---|---|---|
| `voice.enabled` | `VoiceSection.Enabled` | `voice.enabled` | **true** | `schema.go:245`（struct 头 `:244`，section 挂点 `:114`），标签 `default:"true"` |
| `wake_word.enabled` | `WakeWord.Enabled` | `voice.wake_word.enabled` | **false** | `schema.go:197`（struct 头 `:196`，嵌在 VoiceSection 里 `:247`），标签 `default:"false"` |
| `audio.mic_muted_default`（第三枚，票面没点名、唯一真挡得住麦的） | `AudioSection.MicMutedDefault` | `audio.mic_muted_default` | **true** | `schema.go:277`（struct 头 `:269`），注释 `:276` "starts every session muted" |

- 默认值由标签单点决定（D36）：`NewDefaults()` `defaults.go:58` → `applyDefaults` `:65` 逐叶读 `default:"…"`，Bool 分支 `setDefault` `:96-101`（`strconv.ParseBool`）。正控：`sed -n '197p;245p;277p' schema.go` 逐字回 `default:"false"` / `default:"true"` / `default:"true"`。
- **本腿零改动**：AC#4 的"一字不改"由"我没动 schema.go 任何一个字节"满足（写点唯一＝本件，`git show --stat` 每发都只有 `.scratch/wisp/probes/247/a2/**`）。

### 3.2 这两把开关今天**没有**任何产码用来决定开麦

- 尺：`grep -rn "Voice\.Enabled\|WakeWord\.Enabled" --include=*.go . | grep -v _test` ⇒ 产码命中只有 `internal/config/manager.go:380`（reload 分层比较）、`:394`、`:395`（把新值抄回活配置）。**用途是"改了要不要重载"，不是"要不要开麦"**。
- 正控：同尺排测试前会多命中 `internal/config/boundary_test.go:91`（`c.Voice.Enabled = false`）；且票 255 自己的读取名册 `cmd/wisp/config_readers_255.go:123` 逐字 *"扫描零命中：nothing outside internal/config reads cfg.Audio"*、`:140-143` 对 `voice.*` 同样判 no-reader ⇒ 与我的 grep 同向、且是**盘上现成的第三方读数**。
- ⚠ 纠 `247-a1`/`A485`：U13/现量引的是 `manager.go:370/:384/:385`，现树是 `:380/:394/:395`（见 §6）。
- `mic_muted_default` 的产码读取处＝**零**（尺 `grep -rn "MicMutedDefault\|mic_muted_default" --include=*.go . | grep -v _test | grep -v ".scratch"` ⇒ 只剩 `schema.go:277` 定义与 `gate.go:55-56` 的**注释**"map … here at boot wiring"；没有任何产码真的把它传给 `WithStartMuted`）。⇒ 三把开关今天对"开麦"**都没有方向**，因为没消费者。

### 3.3 "默认档下这条路根本不启动"在现树里靠哪一行做到

**诚答：现树里它靠的是"没有那一行"，不是靠任一开关。** 决定"不启动"的，是 §1.1/§2-H1 的**缺席**本身——`cmd/wisp/resident_windows.go` 的 import 块（`:14-16`）里**没有 `internal/audio`**，`mic.Start` 产码零调用者。一双击 `wisp.exe` 起的是 `runResident`：它建球（`:163`）、建面板（`:148`）、建任务源（`:206`），**从不构造采集器**。所以默认档（voice=true × wake_word=false）下麦克风不会被这条路打开——因为这条路今天不存在，不是因为那两枚布尔值挡的。

**一旦接线（`247-r1` 之后），"默认档不开麦"这件事必须由下列几枚行的组合来保证**（缺一即破，交 `247-v1` 逐枚钉）：

| 保证 | 决定它的行 | 形状 |
|---|---|---|
| 装了门且 `mic_muted_default` 真传进去 | `gate.go:108`（`if g.effectiveOpen()`）× `gate.go:235`（`effectiveOpen` 首行 `if g.muted { return false }`）× 接线侧调 `WithStartMuted(c.Audio.MicMutedDefault)`（`gate.go:57`） | muted=true ⇒ `openInnerLocked`（`gate.go:247`）不跑 ⇒ 内层 `mic.Start` 不跑 ⇒ `wasapi_windows.go:210 Open` 不跑 |
| `voice.enabled=false` 视为"根本不构造采集器" | **今天没有这一行**——`A485`-P1 要求 `247-r1` **新增**（现树没有任何产码读 `Voice.Enabled` 做开麦决策，见 3.2） | 若不补这行，`voice.enabled=false` 在默认档下挡不住任何东西 |
| ⛔ 反面：不装门、直接 `mic.Start(ctx, buf)` | `wasapimic_windows.go:88`（`return <-started`＝同步等设备开成功） | **一双击就开麦**（`started` 只有在 `run` 里 open 成功后才送 nil） |

⇒ 结论句：**挡住默认档采音的既不是 `voice.enabled=true` 也不是 `wake_word.enabled=false`，而是"没人接线"这件事本身**；接线时唯一自带的默认是 `audio.mic_muted_default=true`，而它**只在门被装上、且有人把它传进 `WithStartMuted` 时才生效**（今天连"传进去"这步都还没有产码做）。这与 `A485`-P1 裁"甲"完全一致，本腿不裁、只把"改哪几行会失效"摆明。

### 3.4 "如果动了默认值会怎样"（本票零改动，只讲清）

- `audio.mic_muted_default` → false：装了门也会在启动时开设备（`WithStartMuted(false)` ⇒ `effectiveOpen()` 真 ⇒ `openInnerLocked` 跑）。三枚里唯一**直接改开麦行为**的默认。
- `voice.enabled` → false：**今天零后果**（无决策读取处，3.2）；只有 `247-r1` 补上"voice.enabled=false 不构造采集器"那行后它才开始挡。它的约束力是**接线造成的，不是默认值自带的**。
- `voice.wake_word.enabled` → true：**今天同样零后果**；KWS 落地属唤醒词票（票面 AC#7 明列唤醒词/ASR/TTS 都不属于本票）。同向证据：进 Armed 的 `EvKwsEnabled`（`events.go:18`）产码零发射者（`table.go:55` 只有表行）。
- ⚠ 别把"零后果"读成"无关"：`voice`/`wake_word` 的 `Enabled` 都在 reload-tier，改值会走重载计划（`manager.go:380` 的比较、`:394/:395` 的抄回），不是设备开麦决策。

_§3 量数时刻：`2026-10-03 10:0x +0800`。_

## §4 真机读数怎么量

### 4.1 今天有没有"把电平打进日志/仪器"的现成出口：**没有**（尺＋正控）

- 唯一**接收**电平的产码口＝`Ball.SetAudioLevel(float32)`（`internal/ball/liquid_windows.go:42`）。它本身**不打日志**：`:46` 只是 `b.sta.PostTask(func(){ b.applyLevel(level) })`；而 `applyLevel` 第一行 `:52 if !prototypeVisuals { return }`，在 `:55 b.liqRaw = level` **之前** ⇒ 默认档（prototypeVisuals 关）下连"最后收到的电平"都不落到球里可读。
- 尺（有没有任何电平日志出口）：`grep -rn "level|Level" internal/ball/liquid_windows.go | grep -i "slog|Print|printf|log\."` ⇒ 空（`level_log_exit=1`）⇒ **没有现成的电平读数出口**。
- 正控（尺能认电平这个词）：同文件 `grep -n "SetAudioLevel\|applyLevel\|liqRaw\|level float32"` 大量命中（`:42/:51/:55`）⇒ 尺不是对 "level" 失明，是真的没有打点。
- `observe` 侧的"事件名"里也**没有一枚载电平**：`ClassAudioDevice = "audio_device"`（`internal/observe/errors.go:25`）、`EvAudioDeviceLost = "audio.device-lost"`（`internal/statemachine/events.go:30`）、`EvKwsEnabled = "kws.enabled"`（`:18`）——全是分类/状态事件，不带 `float32` 电平。

### 4.2 可复用的台架（逐字 flag）——但它喂的是**合成数**，不能当 AC#2 的证

`cmd/balldebug/main.go` 现读 flag 定义：

| flag（逐字） | 行 | 作用 |
|---|---|---|
| `level := flag.Float64("level", 0, "synthetic audio envelope 0..1 pushed to the ball at ~30fps (0 = feed nothing)")` | `:106` | 手填一个合成包络常量 |
| `diffLevel := flag.Float64("diff-level", 0, "with -diff: the -level handed to each child")` | `:107` | -diff 子进程用的 level |
| `frozen := flag.Bool("frozen", false, "render the frozen SPEC-08 §2.1 visuals instead of the ticket 62 prototype")` | `:97` | `:122 EnablePrototypeVisuals(!*frozen)` 是**唯一**打开可见性的产码 |
| `state := flag.String("state", "", "show one state for 2s (name as in D43)")` | `:91` | 把球摆进 liquidDriven 内的态（否则电平被 §2-H5 丢） |
| `tour := flag.Bool("tour", false, ...)` | `:111` | 引导走查 |

- 调用链：`feedLevels`（`:404`）→ `b.SetAudioLevel(v)`（`:418`）；起协程 `spawnLevelFeeder`（`:427`）→ `observe.Default.Spawn(goroutineLevelFeeder, "balldebug", root, ...)`（`:432`），`goroutineLevelFeeder = "balldebug-level-feeder"`（`:46`，名册外名，owner 标 "balldebug"）。
- `SetAudioLevel` 的调用者名册（尺 `grep -rn "SetAudioLevel" --include=*.go internal cmd` 现读）：产码＝`cmd/balldebug/main.go:418`/`:477` 两枚；测试＝`internal/ball/live_windows_test.go:637/648/656/685`（且该文件 `//go:build windows && winlive`）；`cmd/wisp`＝**零**。⇒ 票面现量 2 那句"非测试生产者只有 cmd/balldebug"**现树仍成立**（无翻转）。

### 4.3 "本机可量"与"CI 可覆盖"必须分开——本枚读数属"本机可量"，且今天还量不到

**（甲）真机那一半＝只有本机可量，不是 CI 可覆盖。** 现读的闸：

- `TestLiveWasapiSmoke`（`internal/audio/hotplug_test.go:525`）：`if os.Getenv("WISP_LIVE_MIC") != "1" { t.Skip("live WASAPI smoke requires WISP_LIVE_MIC=1 and a real microphone (真机冒烟待票 16)") }`（`:526-527`）⇒ 默认跳过。
- `TestPinnedThreadStable10s`（`:443`）：`if testing.Short() { t.Skip(...) }`（`:444-445`）⇒ `-short`（CI 常见跑法）跳过。
- 球的 `live_windows_test.go`／`hotkey_live_test.go`／`interaction_live_test.go`／`live_guard_windows_test.go` 全带 `//go:build windows && winlive` ⇒ 不打 `-tags winlive` 连编译都不进。
- ⛔ 本腿**一枚 `go` 命令都没跑**（§0 自证），上面全是读到的闸，不是跑到的闸。

**（乙）AC#2 那枚"说话 vs 不说话 → `SetAudioLevel` 收到的值不同"的具体读数：**

- 它需要（1）H1/H2 的接线存在（今天**没有**，见 §2）＋（2）一枚真读设备、并把测得的 `float32` 在送进球之前先打出来（今天也**没有**，见 4.1）。⇒ **今天根本没有能取这枚读数的地方**；它是 `247-r1` 落地之后、在**本机**、开着 `WISP_LIVE_MIC`/真设备条件下取的一次性人工读数。
- **CI 可覆盖的只有代理形，不是真机形**：audio 侧有 `NewWavInjector(path string, opts...) (*WavInjector, error)`（`internal/audio/wavinjector.go:48`）→ `Start(ctx, buf)`（`:75`），它从**磁盘 wav 文件**喂 PCM、不碰设备。⇒ 一枚 CI 用例可拿"响 wav vs 静音 wav"过 `FrameLevel`（`level.go:118`）→ `SetAudioLevel`，证"管线能把 PCM 折成不同的 `float32`"。但那是**接缝代理**，证的是 H3+投递，**不是"真麦克风被读了"**。
- ⇒ 具名分档：**"真麦 PCM 使 `SetAudioLevel` 收到两个不同值"＝本机可量、且此刻不可量（等接线）**；**"wav 注入使消费者得到两个不同 `float32`"＝CI 可覆盖（代理）**。本腿把 AC#2 判给前者，并明确警告：若有人拿后者的绿去勾 AC#2，就是把"接缝代理"读成"真机结论"——本仓吃过三次"只在本机成立却写成全局结论"，这枚反向亦然（只在 CI 绿不代表真机过）。

_§4 量数时刻：`2026-10-03 10:0x +0800`。_

## §5 量不到的格子

## §6 我推翻前人哪几句（含票面、A5xx、`247-a1`）
