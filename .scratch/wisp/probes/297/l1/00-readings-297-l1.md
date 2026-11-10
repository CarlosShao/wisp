# 票 297 落地腿 `297-l1` 读数档 — AC#0 仪器＋AC#1 注入面成对读数

- 腿：`297-l1`（票 297 排程节写的是 `297-a1`＝落地腿；该名已被 2026-10-10 只读普查腿占用，本腿改名 `297-l1`，读数落 `probes/297/l1/`）。
- 锚：起手 HEAD＝`a2495ae1`（`git rev-parse` 现量；派单前锚已从票面 `e96da9f4` 漂移多波，本腿一切读数以 `a2495ae1` 树为准）。
- 写面：`internal/audio` 仅新增 3 枚 `*_test.go`＋`probes/297/l1/**`。⛔ 产码 0 字节改动（`git show --stat` 可核）。
- 前置核过：票 300 已 `-done`（parseWaveFormat @24 修复已在树）、票 293 的 `293-v1` 验收已收（AC#0..4 勾）、`tasklist /FI "IMAGENAME eq wisp.exe"`＝0。

## 1. 加了什么（三枚仪器，全部 windows-tagged、包内）

| 件 | 内容 | 判据 |
|---|---|---|
| `internal/audio/convertpacket_297_windows_test.go` | `TestAC297ConvertPacketSyntheticByteFace`：Go 自有 `[]byte` 取址喂 `convertPacket`（5 面）＋`parseWaveFormat` 四元组打印 | S1 float 路径 RMS＝amp/√2 在 `SineLevelTolerance` 内；S2 同一串 float 字节按 int16 切＝逐字节 int16 对均值（精确相等）；S3 silent 全零；S4 int16 方波＝精确幅度；S5 四元组 tag=3/1 |
| `internal/audio/amplitudepair_297_windows_test.go` | `TestAC297KnownAmplitudePairSeparates`：真采集环＋假 opener＋真 wav 字节（票 247 台件同形），两窗已知幅度 0.05/0.4 | 逐帧 level 与源窗 `LevelOfSamples` 精确相等；mean 比＝幅度比（±0.01） |
| `internal/audio/livepacket_297_windows_test.go` | `TestAC297LivePacketByteFace`：真设备在程仪器（`WISP_LIVE_MIC=1` 门，同 `TestLiveWasapiSmoke` 形状），打印 ⓐ四元组 ⓑ首包 int16 头 16 ⓒ同字节 float32 头 8 ⓓ窗级 FrameLevel＋64 窗统计＋判别句 | 四行都有；COM 失败/无包才 Fatal；判别句只打印不判案（AC#2 归编排者） |

## 2. AC#0 读数（真机两发，verbatim 日志在 `live-instrument-run1.log` / `run2.log`）

设备：`麦克风 (Realtek(R) Audio)`，id `{0.0.1.00000000}.{1115a660-460f-4ff0-9e2f-bdf480451e19}`。

| 面 | run1 | run2 |
|---|---|---|
| ⓐ mix format 四元组 | **tag=3 channels=2 rate=48000 bits=32**（IEEE float stereo 48k） | 同 |
| ⓐ2 convertPacket 路径 | `floating=true`（正确 float 路径） | 同 |
| ⓑ 首包 int16 头 16 | 全 0（安静房间，float 路径正确时首包静默） | 全 0 |
| ⓒ 同字节 float32 头 8 | 全 0 | 全 0 |
| ⓓ 64 窗×512 样本 | **mean=0.003513 max=0.078400 min=0.000000** silentPackets=0 | **mean=0.002175 max=0.030298 min=0.000000** |

判别句（仪器打印）：`no signature matched`——ⓑⓒ 都平滑、电平 0.002–0.0035，**不是** 0.37。

## 3. 决定性对拍（三读数同页）

| 读数 | 数值 | 来源 |
|---|---|---|
| 票 247 真机六窗底（2026-10-09，修复前） | mean=0.367691（RMS ≈ 12,050 counts） | 票面现量逐字 |
| 合成「float32 按 int16 切」形状（本腿 S2/D2） | **level=0.383045（RMS 12,552 counts）** | `convertpacket_297_windows_test.go` D2 行 verbatim |
| 修复后真机安静底（本腿 live 两发） | mean 0.003513 / 0.002175 | `livepacket_297_run*.log` |

⇒ 合成切错字节形状的 RMS（12,552 counts）落在旧底（12,050 counts）±4% 内；而修复后同一台设备、同一把尺读 0.002–0.0035。**旧 0.37 底＝float32 尾数噪声被按 int16 切（AC#2 形②），且票 300 的 @24 修复已把它治掉**——形③（设备本底 0.37）被本腿两发真机读数排除。AC#2 择一处置归编排者，本腿只交料。

## 4. AC#1 注入面读数（verbatim）

```
AC#1 injection face: window A (amp 1638 counts) mean=0.049988, window B (amp 13107 counts) mean=0.399994
AC#1 injection face: mean ratio = 8.0018, amplitude ratio = 8.0018 -- the ruler SEPARATES at the injection seam
```

成对读数：**注入面 8.0018 比＝幅度比（分得开）；真麦面修复前 1.0402（票 247 六窗，分不开）**；修复后真麦面底 0.002–0.0035（本腿 §2）。`LevelOfSamples` 尺本身无恙——票面那句"拿到之后才允许说问题在设备/转换那一侧"的两个前提都齐了。

## 5. 门禁读数（每把门一行 rc）

- `go test ./internal/audio/ -count=1` 改前两发：`ok ... 17.531s` / `ok ... 17.586s`（rc=0；`baseline-audio-run{1,2}.log`）
- `go test ./internal/audio/ -count=1` 改后两发：`ok ... 19.081s` / `ok ... 18.374s`（rc=0；`after-audio-run{1,2}.log`）⇒ 新增红 0 枚
- `go test ./cmd/wisp/ -count=1` 改前两发：两发红名册**同集**＝`TestTicket198FirstRunCallerIsTheRunEntryOnly` / `Test313WritableUnionEqualsTheTwoRosters` / `TestLockedSuffixNeverTakesTheLockItself33r11`（448.972s / 406.526s；`baseline-cmdwisp-run{1,2}.log`）——cmd/wisp 在册历史红，与本腿写面（internal/audio 仅 `*_test.go`，不编入 cmd/wisp 二进制）无关
- `go test ./cmd/wisp/ -count=1` 改后两发：两发红名册与基线**同集**（同上 3 枚；342.046s / 359.000s；`after-cmdwisp-run{1,2}.log`）⇒ 两发取交集、逐名作差＝**新增红 0 枚**
- `go vet ./internal/audio/` rc=0
- `gofmt -l cmd internal tools` 空（rc=0）；`gofumpt -l cmd internal tools` 空（rc=0；`D:\work\base\gopath\bin\gofumpt.exe`）
- `GOFLAGS= go build ./cmd/... ./internal/... ./tools/...` **rc=0**（`gate-build-scoped.log`；第二次现跑复坐实）；整树 `go build ./...` rc=1 但失败点＝**gitignore 目录 `build/` 里一枚未跟踪遗留分片 `build/mut_rvw.go`**（package main、未定义符号 11 处，`.gitignore:14` 现量坐实 `build/` 整目录被忽略）——预存本机状态，与本腿零关系（本腿零文件落 build/），两读数都如实落 `gate-build.log` / `gate-build-scoped.log`
- `sh scripts/d22scan.sh` rc=0，`d22scan: clean - no D22 ban violations`（458 production Go files examined；`gate-d22scan.log`）
- `sh scripts/check-path-length-budget.sh` rc=0，`VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`（`gate-path.log`）
- `tasklist /FI "IMAGENAME eq wisp.exe"` 起手＝`INFO: No tasks are running which match the specified criteria.`

## 6. 真机面收官：票 247 六窗用例在修复后树上复跑一发（`live-247-sixwindow-after.log`，verbatim）

尺＝`go test ./cmd/wisp/ -run TestAC247LiveMicrophoneLevelsReachTheBallSeam -count=1 -v`（`WISP_LIVE_MIC=1`，外放闹铃 ~6s；判据行未动一字＝`resident_audio_247_live_windows_test.go:194/195`）：

```
AC#2 counters at the seam: levels=377 frames_sent=6 frames_dropped=371 last_err=""
AC#2 reading: quiet-A: samples=93 mean=0.003003 max=0.116988 levels_per_s=31.00
AC#2 reading: sound: samples=190 mean=0.001358 max=0.002714 levels_per_s=31.67
AC#2 reading: quiet-B: samples=94 mean=0.001287 max=0.001533 levels_per_s=31.33
the loud phase is not distinguishable from silence: ... (speakers off, or routed to an endpoint this microphone cannot hear)
```

**三形对照（本腿读数给出的支持行）**：

- **形②（字节切错）＝旧 0.37 底的成因，且已被票 300 的 @24 修复治掉**：修复前 quiet mean=0.3677（RMS≈12,050 counts）↔ 本腿合成切错字节形状 level=0.383045（RMS 12,552 counts，`convertpacket_297_windows_test.go` D2 行）↔ 修复后真机 quiet mean=0.003003/0.001287（本节 verbatim）＋仪器 64 窗 mean=0.003513/0.002175。三读数同页对齐。
- **形③（设备本底 0.37）＝被排除**：同一台 Realtek 麦、同一把尺，修复后安静底 0.0013–0.0035（−50 dBFS 量级），不再是 0.37。
- **形①（声音没到麦）＝现下唯一活着的红句成因**：sound 窗 mean=0.001358 **低于** quiet-A（0.003003），max 仅 0.002714——闹铃外放期间麦克风听到的是静默。capture/render 端点名已打出（同一 Realtek 设备族），故嫌疑落在"扬声器当时物理静音/系统音量/外放路由"这个事实问句上——正是红句括号里点名的候选，要机主的耳朵或一次点名音量才能定。⚠ 本腿⛔ 不裁：这一发与立票时那发（0.3677 六窗）的区别恰恰说明底是切错字节、不是设备本底；"响得起来"那格是否可判，归编排者按 AC#2 择一处置。

`resident_audio_247_live_windows_test.go:194`/`:195` 一字未动；本腿零产码改动。
