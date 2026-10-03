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

## §3 两把开关与默认档自证

## §4 真机读数怎么量

## §5 量不到的格子

## §6 我推翻前人哪几句（含票面、A5xx、`247-a1`）
