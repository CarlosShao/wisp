# 241 — 没有任何一处把 PCM 算成"音量有多大"：全主模块 `Sqrt` 零命中（电平链生产者半）

Status: OPEN（编排者 09-30 11:5x 立，料全部出自 `240-c1` 的现量普查）
实现者：`241-r1`（写码腿） · 裁决者：`241-v1`（必须≠实现者） · 归口：语音链路（票 41/26/15 的前置，非其替代）

## 现量（全部 `240-c1` 或编排者本轮自己跑的尺，逐条带出处）

| # | 事实 | 出处 |
|---|---|---|
| 1 | 全主模块 `Sqrt` **零命中**＝压根没算过音量 | `240-c1` §1 尺 R:，rc 实跑 |
| 2 | `0..1` 这个标度**全仓无定义**；`SilenceLevelGate = 0.06` 是"判静默的门限"，不是量程 | ⚠ **本行出处我抄错了（`241-r1` 当场顶正，票面已改）**：那枚常量**不在** `internal/audio/gate.go`，真身在 **`internal/ball/liquid.go:30`**（判语本身不受影响，`internal/audio` 里今天只有 `level.go` 的注释在引它）。普查 §1.4 写对了，是我抄票时抄错 |
| 3 | 采集栈**写完了零 importer**：真 WASAPI 共享模式在位（`wasapi_windows.go:210 Open`、`:256 init`、`:361 Drain`），但全仓零文件 `import internal/audio` | `240-c1` §1 R1/R3 |
| 4 | `NewBoundedFrames`（`audio.go:65`）**零消费者**；`encodeFrame`（`wavinjector.go:135`）**未导出**，包外拿不到解码后的帧 | `240-c1` §1 |
| 5 | 球那边唯一的入口 `SetAudioLevel`（`internal/ball/liquid_windows.go:42`）今天**非测试生产者只有 `cmd/balldebug`**（`main.go:418/:477`，flag 文案逐字 `synthetic audio envelope`）⇒ 现在屏幕上音量是**命令行里手填的数** | A460、A462、`240-c1` R4 |
| 6 | `internal/audio` 里 `RMS`／`Level`／`Peak` 符号**零命中**（编排者本轮 `grep -rn "func .*RMS\|func .*Level\|Peak\|Sqrt" internal/audio/` 空）⇒ 本票新符号不撞钉 | 编排者 09-30 11:5x 现跑 |
| 7 | 包基线全绿：`go test ./internal/audio/ -count=1` ⇒ **ok 15.457s**，26 枚用例（名册见下"起跑名册"） | 编排者 09-30 11:5x 现跑 |

## 本票要做的（只做生产者半，一格都不许多）

把"一串 PCM 帧"变成"一个 `0..1` 的电平数"，并且**把这个标度本身定义出来**（事实 2 缺的就是这个）。落点在 `internal/audio` 包内：帧进来 → 一个数出去，可被任何消费者拿去用。

**明确不属于本票**（多做了算越界，`241-v1` 直接退回）：
- ❌ 不接消费者、不改 `cmd/wisp`（`224-r1` 正在那枚包里改 `run_mode101_test.go`，同包两枚写腿会把红的归因搅浑）；
- ❌ 不动 `internal/ball`（那是票 228／239 的地界，且 `liquid_windows.go:38-40` 逐字 `no samples and no transcript cross into the ball` ⇒ 过界的只许有一个 `float32`）；
- ❌ 不碰唤醒词、不碰 KWS、不碰 ASR/TTS、不引 `sherpa`；
- ❌ 不解决"第六环形状"：`liquidDriven(liquid.go:60-68)` 不含 `Armed` 那一格归票 228 之后的消费腿。

## 起跑名册（`240-c1`／编排者现跑，写腿开工前先自复一遍）

`go test ./internal/audio/ -count=1` 今天**全绿**，26 枚：`TestGateHalfDuplexClosedZeroFrames` `TestGatePathCFullDuplex` `TestGateMutedEvents` `TestGateStartMuted` `TestGateStopIdempotent` `TestHotplugReenumerateOnce` `TestHotplugStaleHandleFailsLoudly` `TestHotplugReenumerateFailureNamesDevice` `TestOpenOccupiedAndPermissionDenied` `TestEndpointsPairQueryable` `TestPinnedThreadStable10s` `TestLiveWasapiSmoke` `TestCaptureLoopWavIntegrity` `TestResample48kTo16kLengthExact` `TestResamplerSineSNR48k` `TestResamplerSineSNR44k1` `TestResamplerChunkInvariant` `TestResamplerLatencyBounded` `TestResamplerPassthroughAndFlush` `TestMonoDownmixAndFloatConvert` `TestWavInjectorFrameExact` `TestWavInjectorRateConversion` `TestWavInjectorBackpressureDropCounted` `TestWavInjectorCtxCancel` `TestWavInjectorBadFile` `TestBoundedWindowInvariants`。
⚠ 其中 `TestLiveWasapiSmoke` 与 `TestPinnedThreadStable10s` 碰真设备：跑之前 `tasklist //FI "IMAGENAME eq balldebug.exe"` 必须为 **0**。

## 判据（⛔ 框归编排者，产码／验收腿一枚都不许碰）

- [x] **AC#1 标度是真的，不是装饰**：新落点必须让"同一个输入进去，出来是同一个数"可被非浮点噪声地钉住——静默帧 ⇒ 严格 `0`（或明确定义的地板值）、满幅方波 ⇒ `1.0` 或具名上界、~~半幅正弦 ⇒ 与 `1/√2` 的偏差有具名容差~~ ⚠ **这半句我写错了口径（`241-r1` 顶正，票面已改）**：满幅正弦才读 `1/√2`，**半幅**正弦读 `1/(2√2)=0.35355`。判据按**两头发各自钉**收（本票终态：两头都钉住了，无论原意是哪一头判据都成立）。三发缺一发算不成立。
- [x] **AC#2 攻"生产者有没有真消费者路径"**：把 `240-c1` 认定的两处断点（`NewBoundedFrames` 零消费者、`encodeFrame` 未导出）**逐条判**：本票要么补上其中一环并给出包内测试，要么具名写清"这一环留给哪枚票"。不许默默留空。
- [x] **AC#3 攻恒真**：任选本票新增的一枚断言做定向突变（把 RMS 分子改成常数、或把容差放宽到任何输入都能过）⇒ **指名用例必须红**。绿着＝那枚断言是装饰，具名登记。
- [x] **AC#4 越界检查**：`git diff` 里出现 `cmd/wisp`／`internal/ball`／`internal/speech`／`PLAN.md`／`docs/specs/**`／`thresholds.go` 任一路径 ⇒ 直接退回（禁区不是建议）。
- [x] **AC#5 门禁四数**：`go build ./...`、`gofumpt -l`、`go vet ./internal/audio/`、`tools/d22scan/d22scan.exe` 终态逐名照抄（⚠ `d22scan` 是**独立模块**，只能跑 `./tools/d22scan/d22scan.exe`，`go run ./tools/d22scan` 必失败）。

## 禁区（实现腿与验收腿共用）

零阈值变更（`internal/observe/thresholds.go`／golden／`docs/SLO.md` 一字节不许动）；⛔ 不许为了变绿放宽任何断言；不动 `PLAN.md`／`docs/specs/**`／`allowlist.txt`；不动三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）；`frontend/**`／`design/**` 零读零写零转述；**本腿不跑 `go test` 之外的门禁、不推送**（⚠⚠ **这半句与 AC#5 打架＝我写票时自相矛盾，`241-r1` 具名上交、我此刻裁**：AC#5 那四数**照跑**，`build`／`gofumpt`／`vet`／`d22scan` 本来就不是 `go test`；**被禁的是"跑别的包的 `go test`"**，那才是本意射程。它按最窄一致解执行（跨包 `go test` 一枚未跑、`gofumpt`/`vet` 收窄到 `internal/audio`）＝**读成合规，不读成越界**；`241-v1` 按 AC#5 判四数）（`internal/audio` 之外的包一律不跑，因为有并行写腿的半成品在里面，红了归不清是谁的）；只 commit、显式 pathspec、禁 `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`。

## 编排者增量 — 09-30 **14:46:54**（标题的钟由 date 的 stdout 插值，不手打）：收「241-v1」⇒ **五枚 AC 全部成立，本票翻满五格并结案**

裁决表「docs/evidence/s1/241-audio-level-producer-v1.md」**603 行／63,642 字节、§0-§8 齐、真占位符 0 枚**（表里那一处「待填」是它自己在数占位符时引用的模式串）。判语：**AC#1 成立（附两枚文档精度缺陷）／AC#2 成立（附两枚具名条件）／AC#3 成立（十七发突变＋注释级正控）／AC#4 成立（附一枚尺口径纠正）／AC#5 成立（四发全 rc=0，每发带"仪器真会响"的正控）**，⛔ **没有一枚需要退回**；交件现跑票面仍是「未勾 5、已勾 0」⇒ **翻框动作归我，此刻五枚全翻**。
**这枚腿值得当范本的地方**：它**没有复用 241-r1 的读数**，自己重搭突变台（M1-M15＋注释级正控）；AC#5 它跑了三遍——初跑、全部突变与写盘动作之后一遍、交件前最后一遍——因为票面要的是"终态"；§8 是一枚会自己抓自己的尺（数本件占位符，rc=1 当场翻译成"没找到＝好消息"）。

### AC#2 那枚条件里，最重的不是它的判语，是它顺手挖出的一条**过期落点指认链**

「cmd/wisp/resident_windows.go:81」现读逐字「the floating ball arrives in ticket 07」——**票 07 早就关闭了**（盘上文件名「07-ball-state-machine-core-done.md」，我现跑），可那行注释还说球"要等票 07 才来"；普查「240-c1 §3.1」照抄了它，「241-r1」再照抄普查 ⇒ **一条过期的指认在三程之间被原样传递，而今天没有任何仪器看得见它**。后果：按这条指认排程会把"球进干活进程"当成别人的活，而那个"别人"已经结案。
**处置（当场改，不留口头）**：
1. 票面正文那两处「欠票 07 的 GUI 腿」的归属**作废**，本票那笔债今天起的有效归属＝**票 228**（球与托盘不在跑任务的那个进程里，未结案）。
2. 「resident_windows.go:81」那行注释的过期性**具名登记给 228／224 地界**（「241-v1」按硬约束没碰「cmd/wisp」，只上交原文），落点＝票 228 面新增一格：改注释为带条件的事实句，并把"谁真的把球装配进来"写成可判的东西。⛔ 我不顺手改——那是别人的地界，且改了要有人判。
