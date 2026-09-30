# 241 — 没有任何一处把 PCM 算成"音量有多大"：全主模块 `Sqrt` 零命中（电平链生产者半）

Status: OPEN（编排者 09-30 11:5x 立，料全部出自 `240-c1` 的现量普查）
实现者：`241-r1`（写码腿） · 裁决者：`241-v1`（必须≠实现者） · 归口：语音链路（票 41/26/15 的前置，非其替代）

## 现量（全部 `240-c1` 或编排者本轮自己跑的尺，逐条带出处）

| # | 事实 | 出处 |
|---|---|---|
| 1 | 全主模块 `Sqrt` **零命中**＝压根没算过音量 | `240-c1` §1 尺 R:，rc 实跑 |
| 2 | `0..1` 这个标度**全仓无定义**；`SilenceLevelGate = 0.06` 是"判静默的门限"，不是量程 | `internal/audio/gate.go`（编排者本轮现读） |
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

- [ ] **AC#1 标度是真的，不是装饰**：新落点必须让"同一个输入进去，出来是同一个数"可被非浮点噪声地钉住——静默帧 ⇒ 严格 `0`（或明确定义的地板值）、满幅方波 ⇒ `1.0` 或具名上界、半幅正弦 ⇒ 与 `1/√2` 的偏差有具名容差。三发缺一发算不成立。
- [ ] **AC#2 攻"生产者有没有真消费者路径"**：把 `240-c1` 认定的两处断点（`NewBoundedFrames` 零消费者、`encodeFrame` 未导出）**逐条判**：本票要么补上其中一环并给出包内测试，要么具名写清"这一环留给哪枚票"。不许默默留空。
- [ ] **AC#3 攻恒真**：任选本票新增的一枚断言做定向突变（把 RMS 分子改成常数、或把容差放宽到任何输入都能过）⇒ **指名用例必须红**。绿着＝那枚断言是装饰，具名登记。
- [ ] **AC#4 越界检查**：`git diff` 里出现 `cmd/wisp`／`internal/ball`／`internal/speech`／`PLAN.md`／`docs/specs/**`／`thresholds.go` 任一路径 ⇒ 直接退回（禁区不是建议）。
- [ ] **AC#5 门禁四数**：`go build ./...`、`gofumpt -l`、`go vet ./internal/audio/`、`tools/d22scan/d22scan.exe` 终态逐名照抄（⚠ `d22scan` 是**独立模块**，只能跑 `./tools/d22scan/d22scan.exe`，`go run ./tools/d22scan` 必失败）。

## 禁区（实现腿与验收腿共用）

零阈值变更（`internal/observe/thresholds.go`／golden／`docs/SLO.md` 一字节不许动）；⛔ 不许为了变绿放宽任何断言；不动 `PLAN.md`／`docs/specs/**`／`allowlist.txt`；不动三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）；`frontend/**`／`design/**` 零读零写零转述；**本腿不跑 `go test` 之外的门禁、不推送**（`internal/audio` 之外的包一律不跑，因为有并行写腿的半成品在里面，红了归不清是谁的）；只 commit、显式 pathspec、禁 `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`。
