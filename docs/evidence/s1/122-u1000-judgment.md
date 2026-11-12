# 票 122 · 波 1 · U1000 逐符号判表（linux/CI 形状 47 条，本波零删）

- 依据：staticcheck 2026.2.2 (0.8.2)，linux（CI）形状，HEAD `095427398`，2026-11-12 现跑。
- 判例三类：**真死码(候选删)** / **等着接线** / **留**（具名为什么工具看不见：GOOS 构建标签 / 反射 / 接口 / embed）。
- 关键背景：本仓 GOOS 双形状的 U1000 集合**几乎不相交**（linux 47 条 vs windows 85 条）——
  linux 形状看不见 `*_windows.go`/`*_windows_test.go`，windows 形状看不见 `!windows` 文件。
  **CI 门（linux）会永远报这 47 条里的 44 条"留"项**；若要把 CI lint 的 U1000 段清绿，路径不是删它们，
  而是编排者裁定读数口径（例如按 GOOS 拆分母）。本表只判归属，不裁决口径。
- 本波**零删**（派单明令，删留波 2）。

## 判：真死码（候选删）——3 条，双 GOOS 形状都报 unused、全仓 `git grep` 零引用

| 包 | 符号 | file:line | 证据与删除后的绿路径 |
|---|---|---|---|
| internal/agent | `(*fakeReminderClock).releasedCount` | internal/agent/reminder_timer_test.go:109 | 测试 helper，`git grep releasedCount -- '*.go'` 仅定义+注释两行；linux 与 windows 两形状均报 unused。波 2 删后跑 `go test -count=1 ./internal/agent/`（GOOS=windows）应绿 |
| internal/memory | `l1WaitFor` | internal/memory/l1_warm_29_test.go:151 | 测试轮询 helper，全仓零引用（仅定义+注释）；两形状均报。波 2 删后跑 `go test -count=1 ./internal/memory/` 应绿 |
| internal/tools | `(*scriptGate).levels` | internal/tools/web_bridge_test.go:49 | 测试 helper，全仓零引用；两形状均报。波 2 删后跑 `go test -count=1 ./internal/tools/` 应绿 |

## 判：留——44 条，全部是 GOOS 形状差（工具在 linux 形状下看不见 windows 消费者）

### internal/ball · jarvis 视觉族 · 38 条（jarvis.go:17-228）

全部被 `internal/ball/jarvis_windows.go`（生产渲染路径：`jarvisLayout`/`jarvisFamilyColors` 等）与
`internal/ball/jarvis_windows_test.go`（结构性测试：粒子预算、确定性、环布局）消费；
windows 形状 U1000 名单里**一条 jarvis 符号都没有**（复跑核过）。留。

| 包 | 符号 | file:line |
|---|---|---|
| internal/ball | `jarvisRingInner` | jarvis.go:17 |
| internal/ball | `jarvisRingMain` | jarvis.go:18 |
| internal/ball | `jarvisRingCore` | jarvis.go:19 |
| internal/ball | `jarvisRingField` | jarvis.go:20 |
| internal/ball | `jarvisRingInnerFrac` | jarvis.go:27 |
| internal/ball | `jarvisRingMainIdleFrac` | jarvis.go:28 |
| internal/ball | `jarvisRingMainThinkFrac` | jarvis.go:29 |
| internal/ball | `jarvisRingMainAlertFrac` | jarvis.go:30 |
| internal/ball | `jarvisRingCoreFrac` | jarvis.go:31 |
| internal/ball | `jarvisCoreDotFrac` | jarvis.go:32 |
| internal/ball | `jarvisRingTiltFrac` | jarvis.go:33 |
| internal/ball | `jarvisFieldMinFrac` | jarvis.go:34 |
| internal/ball | `jarvisFieldMaxFrac` | jarvis.go:35 |
| internal/ball | `jarvisGlowRadiusFrac` | jarvis.go:36 |
| internal/ball | `jarvisCoreGlowFrac` | jarvis.go:37 |
| internal/ball | `jarvisDotRadiusPx` | jarvis.go:38 |
| internal/ball | `jarvisDotBrightPx` | jarvis.go:39 |
| internal/ball | `jarvisBrightShare` | jarvis.go:40 |
| internal/ball | `jarvisRingJitterMain` | jarvis.go:41 |
| internal/ball | `jarvisRingJitterInner` | jarvis.go:42 |
| internal/ball | `jarvisParticleCountBase` | jarvis.go:43 |
| internal/ball | `jarvisParticleCountMin` | jarvis.go:44 |
| internal/ball | `jarvisParticleCountMax` | jarvis.go:45 |
| internal/ball | `jarvisListenPulseFrac` | jarvis.go:46 |
| internal/ball | `jarvisScatterPerLevel` | jarvis.go:47 |
| internal/ball | `jarvisParticleSeed` | jarvis.go:48 |
| internal/ball | `jarvisRng`（type） | jarvis.go:55 |
| internal/ball | `newJarvisRng` | jarvis.go:57 |
| internal/ball | `(*jarvisRng).f64` | jarvis.go:59 |
| internal/ball | `jarvisDot`（type） | jarvis.go:71 |
| internal/ball | `jarvisParticleCount` | jarvis.go:82 |
| internal/ball | `jarvisLayout` | jarvis.go:98 |
| internal/ball | `jarvisMainRingFrac` | jarvis.go:141 |
| internal/ball | `jarvisRingSpin` | jarvis.go:155 |
| internal/ball | `jarvisListenPulse` | jarvis.go:164 |
| internal/ball | `jarvisScatterAmp` | jarvis.go:172 |
| internal/ball | `jarvisClamp01` | jarvis.go:176 |
| internal/ball | `jarvisFamilyColors` | jarvis.go:228 |

（jarvis.go 是 `!windows` 排他文件还是无标签通用件以 `head -1` 现读为准；判定不受影响——
消费面在 windows 独占件里，linux 形状必然整簇报 unused。）

### internal/audio · 设备选择簇 · 3 条

| 包 | 符号 | file:line | 为什么留 |
|---|---|---|---|
| internal/audio | `captureDeviceLister`（interface） | device.go:60 | 被 `selectCaptureDevice`（device.go:86 的类型断言）用；传递根因是 `selectCaptureDevice` 的调用方全在 windows 面（`input_device_180k_windows_test.go:45/64/76/95`）。windows 形状不报 |
| internal/audio | `inputDeviceDefault` | device.go:66 | 被 `selectCaptureDevice`（device.go:79）用，同上传递 |
| internal/audio | `selectCaptureDevice` | device.go:77 | windows 测试在用（`input_device_180k_windows_test.go` 4 处）；windows 形状不报 |

### internal/tools · 参数 helper · 2 条

| 包 | 符号 | file:line | 为什么留 |
|---|---|---|---|
| internal/tools | `argInt` | system_bridge.go:274 | `system_backend_windows.go` 12+ 处调用（:889/:899/:906/:986/:1073-1075/:1177/:1289/:1315-1318/:1406…）；windows 形状不报。票 122 批 2 已判过一次（留），本波复核同判 |
| internal/tools | `argBool` | system_bridge.go:297 | 同上（:899） |

### internal/plugin/manifestslot · 1 条

| 包 | 符号 | file:line | 为什么留 |
|---|---|---|---|
| internal/plugin/manifestslot | `childJob.exited`（field） | command_job_other.go:13 | `//go:build !windows` 占位件的镜像字段：windows 侧 `command_job_windows.go:27` 同名字段是活的（:82/:132/:155-157）；此件的文件头注释明说它是"非 Windows 占位，让包在 cross-vet 腿下能编译"（D7 Windows-only）。字段是对齐 windows 变体形状的刻意镜像，不是漏删 |

## 波 2 建议摘要（供编排者）

1. 删 3 条真死码（上表），每条删后跑对应包 `go test -count=1`，纯机械。
2. linux 形状剩 44 条 U1000 全是 GOOS 形状差，**删不得**；CI lint 的 U1000 段清绿只能走读数口径裁定
   （windows-only 符号在 linux 分母里报 unused 属仪器盲区，不是欠账）。
3. windows 形状另有 85 条 U1000 + 7 条非-U1000（SA4004/SA4010/ST1005×2/ST1012×2/S1016）待单独成波
   （其中 internal/ball 40+ 条、cmd/wisp 若干，写面要与 ball/cmd 单写腿规矩协调）。
4. SA5000（panic_isolation_51_test.go:62，有意为之的崩溃夹具）与 SA9009（frontend/embed.go:4，已知假阳性）
   两条需要编排者裁定读数口径豁免，否则 linux 非-U1000 段永远剩 2。
