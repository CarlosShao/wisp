# T13 对抗验收报告（T13-adv 独立执行）

> 执行者：T13-adv 子代理（实现者为 T13-impl，相互独立）。
> 时间：2026-09-19T15:20Z 前后（本地 2026-09-19 深夜）。
> 环境：go1.27.1 windows/amd64（CGO_ENABLED=1，GOPROXY=goproxy.cn），Windows 11 x64 真机
> （Realtek 麦克风/扬声器在位）。并行保护：T09-impl 在同工作树开发 internal/llm、tools/mockllm，
> 本验收全程未触碰，git 操作均为显式路径。

## 1. 全量复跑（gates）

| 命令 | 结果 | 证据 |
|---|---|---|
| `go build ./...` | PASS | BUILD_EXIT=0 |
| `go test ./...` | PASS | TEST_EXIT=0；audio/ball/buildinfo/config/llm/golden/openaichat/memory/observe/plugin/proc/secret/statemachine 全 ok |
| `go test -race -count=1 -v ./internal/audio/...` | PASS | 26/26 测试 PASS（TestLiveWasapiSmoke 因未设环境变量 SKIP，预期），17.030s，无 race 报告 |
| `go vet ./internal/audio/...` | PASS | VET_AUDIO_EXIT=0 |
| `go vet ./...` | 存疑→非 T13 缺陷 | 仅 2 条诊断：`internal\llm\openaichat\adapter_test.go:190/238 self-assignment of events`——该文件是 **T09-impl 未跟踪的在途 WIP**（git status `??`），非 T13 代码；T13 包 vet 干净。T13 收口前若需全仓 vet 绿，须待 T09 提交后复验 |

MINOR：实现交接声称 "audio 21 tests"，实测包内 26 个测试函数（少计，全部绿，无风险）。

## 2. 真机冒烟复跑（独立复现）

`WISP_LIVE_MIC=1 go test ./internal/audio/ -run TestLiveWasapiSmoke -v -count=1` → **PASS**

```
capture="麦克风 (Realtek(R) Audio)" render="扬声器 (Realtek(R) Audio)"
live smoke OK: 25 frames in 847ms, stats {FramesSent:25 FramesDropped:0 BytesSent:25600 BytesDropped:0 Reopens:0}
```

25 帧/847ms ≈ 32ms/帧节奏 + 线程启动开销（测试上限 1500ms）；帧长恒 1024B。独立复现实现方声称的 25 帧/820ms。

## 3. 代码审计（关键行为断言非常量）

**a) drop 计数真实递增 + 日志节流 — PASS**
`audio.go:115-145 meter.push`：非阻塞 send，`default` 分支真实递增 `dropped/bytesDropped`；节流 `total==1 || time.Since(lastWarn)>=1s`（首丢必记，D42#9）。测试 `TestWavInjectorBackpressureDropCounted` 用 flood 注入 + 200ms 慢消费者，断言 `FramesDropped>0`、`BytesDropped==FramesDropped*1024`、日志含 `audio frames dropped` 与 `dropped_frames_total=`（计数可见，非静默丢弃）。

**b) 半双工门 SetSpeaking(true) 是真实 Stop — PASS**
`gate.go:151-155/255-259`：`!shouldOpen && g.open` → `closeInnerLocked()` → `inner.Stop()`；gate 内**无任何数据过滤/丢弃路径**。`TestGateHalfDuplexClosedZeroFrames` 做的是冻结断言而非常量断言：`frozen := inner.Stats().FramesSent; sleep(150ms); got != frozen → fail`——证明内层采集真实停转（设备关），而非数据被丢弃式过滤。reopen 后帧恢复。

**c) 重采样分块不变性逐位对比 — PASS**
`TestResamplerChunkInvariant`：4823 个随机样本，随机分块（1–700 样本/块）流式处理，与 one-shot **逐样本 `!=` 比较（即逐位）**，长度也逐位相等。`TestResamplerLatencyBounded` 进一步逐样本馈送仍逐位一致（1 输入样本延迟界）。实现为 32.32 定点绝对坐标累加器（`resample.go:49-70`），分块无漂移由构造保证并由测试钉死。

**d) 热插拔"重枚举恰一次" + stale-handle 具名报错 — PASS**
`wasapimic_windows.go:270-280 reopenOnce`：单次 `openCurrent()`（一次枚举 + 一次 Open），无重试环；失败即终态 `observe.Wrap(..., "hotplug reopen failed; keeping no stale handle (previous device "+prev.String()+")")`——命名被拒绝保留的旧设备。
- `TestHotplugReenumerateOnce`：注入 OnDefaultDeviceChanged 后断言 `opens==2`（初始 1 + 重开 1）、`Reopens==1`、第二端点数据经真实管线流出；再次注入 → `Reopens==2/opens==3`（每事件恰一次）。
- `TestHotplugStaleHandleFailsLoudly`：8 周期后设备失效 + 重开被占用 → 终态错误 class=audio_device、消息含设备名 "Dying Mic" 与 in-use 引导文案、`opens==2`（恰一次重开尝试）、`Stats().LastError` 非空、终止后通道无帧流出（不静默保留旧句柄）。
- `TestHotplugReenumerateFailureNamesDevice`：枚举本身失败同样终态并命名 prev 设备 "Unplugged Mic"。
- D37 语义核对：`internal/observe/errors.go:109-110` `ClassAudioDevice → RetryOnce`，与实现一致。

**e) 音频数据通道零 mock — PASS**
grep 非测试代码：`mock|fake|stub` 仅命中注释（device.go/wasapimic_windows.go 的 seam 说明）。假件（fakeWatcher/fakeOpener/fakeStream）全部位于 `hotplug_test.go`，注入点是**枚举/打开 seam**；数据全部流经真实捕获循环（WaitEvent→Drain→Resampler→512 帧→metered push）。`TestCaptureLoopWavIntegrity` 用生产 `parseWav` 解析真实 wav fixture、经真实循环、逐帧 `bytes.Equal` 比对。

**vtable 勘误独立核对 — PASS**
`wasapi_windows.go:422-430` 槽位 Initialize=3/GetMixFormat=8/GetDevicePeriod=9/Start=10/Stop=11/SetEventHandle=13/GetService=14，与 audioclient.h MIDL 声明序（IUnknown 0-2 后）一致；IAudioCaptureClient GetBuffer=3/ReleaseBuffer=4/GetNextPacketSize=5、IMMDeviceEnumerator GetDevice=5、IMMDevice Activate=3、Register/UnregisterNotificationCallback=6/7 均核对无误。`AUDCLNT_E_BUFFER_SIZE_NOT_ALIGNED`（0x88890014）按 MS 文档配方以设备周期对齐重试（`alignedBufferHns`，周期查询失败回退系统默认）。

## 4. D16 合规 — PASS

- 日志路径 grep：包内全部 6 处 slog（audio.go:136/152、mmdevice_windows.go:60/65、wasapimic_windows.go:112/277）仅携带计数、时长、设备名、错误串——**零帧 payload**（`encodeFrame` 数据从不进日志）。
- 落盘扫描：非测试代码无 `os.WriteFile/CreateFile/OpenFile`（唯一文件访问是 WavInjector 的测试 wav 读入）。音频缓冲无任何持久化。
- `keep_audio` 硬编码 false 在 config 层强制：`internal/config/validate.go:82-84` 拒绝 `privacy.keep_audio=true` 加载。
- `[audio] mic_muted_default` 布线点在位：`internal/config/schema.go:276-277`（default true）→ gate `WithStartMuted`（`TestGateStartMuted` 钉死启动静音）。

## 5. 越界扫描 — PASS

| commit | 触碰文件 | 判定 |
|---|---|---|
| 2aeae00 | internal/audio/{audio,resample,resample_test,wavinjector,wavinjector_test}.go + 票13 | 合规 |
| 4e3c132 | internal/audio/{gate,gate_test}.go + 票13 | 合规 |
| 77538f1 | internal/audio/{device,doc,hotplug_test,mmdevice_windows,wasapi_other,wasapi_windows,wasapimic_windows}.go + 票13 | 合规 |
| e1a64cf | 票13（handoff 日志 1 行） | 合规 |

4 个 commit 零卷入 internal/llm、tools/mockllm（mockllm 提交 0839407 属 T09）；emoji 扫描 4×0 行命中；密钥模式（api_key/secret/password/BEGIN RSA/sk-）0 命中；go.mod 无新增依赖（x/sys 为既有依赖，无 cgo）。

## 6. 票据对照（6 条验收逐条）

| # | 验收条件 | 裁决 | 证据（独立复跑 -race 全绿） |
|---|---|---|---|
| 1 | WavInjector→channel 精确帧 + 背压丢帧计数可见 | PASS | TestWavInjectorFrameExact（逐帧 bytes.Equal）+ TestWavInjectorBackpressureDropCounted（Dropped>0、字节账目一致、日志含 dropped_frames_total） |
| 2 | 钉线程 10s 不迁移、无外来工作 | PASS | TestPinnedThreadStable10s（10.07s：线程 id 集合大小==1；4 个采样 goroutine 证实无他者落在钉定线程） |
| 3 | 48k→16k 重采样 SNR + 延迟界 | PASS | TestResample48kTo16kLengthExact + TestResamplerSineSNR48k（1kHz SNR≥40dB）+ SNR44k1（非整比 2.75625）+ TestResamplerLatencyBounded |
| 4 | 热插拔注入→重枚举恰一次；stale-handle 具名报错 | PASS | TestHotplugReenumerateOnce + TestHotplugStaleHandleFailsLoudly + TestHotplugReenumerateFailureNamesDevice（见 §3d） |
| 5 | 占用/权限 fixture → Error(audio_device) 带引导文案 | PASS | TestOpenOccupiedAndPermissionDenied（occupied→inUseGuidance；permission→privacyGuidance 含 Windows 设置路径；均含设备名，Start 即失败不静默） |
| 6 | 半双工门：关闭→零帧；重开恢复 | PASS | TestGateHalfDuplexClosedZeroFrames（真实 Stop + 冻结断言 + reopen）；Path C 直通、Muted 双路径、Stop 幂等各有专测 |

**真机冒烟待办登记**：已登记——票 13 Progress log 14:58:45Z "LIVE smoke on real Realtek mic PASS (WISP_LIVE_MIC=1 hook retained for ticket 16)"；测试注释 "真机冒烟待票 16"；票 16 What-to-build 含 real-device smoke、pre-mortem 清单含 mic occupied / permission denied / hotplug mid-utterance 三条物理场景（即归票 16 的人工物理场景）。
**Progress log**：5 条 append-only 完整（claim → c8-seam → half-duplex-gate → wasapi+hotplug+errors → handoff）。
**Status**：in-progress，Claimed by orchestrator → T13-impl；未自行标 done（收口决策归编排者，符合协议）。

## 问题分级

- **BLOCKER**：无
- **MAJOR**：无
- **MINOR**：
  1. 实现交接声称 "audio 21 tests"，实测 26 个测试函数（少计，全绿，无风险）。
  2. 实现交接声称 "internal/audio/ 9 文件"，实测 10 个源文件 + 4 个测试文件（少计）。
  3. 当前工作树 `go vet ./...` 红——仅因 T09-impl 未跟踪在途文件 `internal/llm/openaichat/adapter_test.go`（:190/:238 self-assignment of events）；非 T13 缺陷，但 T13 正式收口时若以全仓 vet 为 gate，需在 T09 提交后复验一次。

## 裁决

六条验收标准全部由真实行为断言钉死并独立复跑通过（-race）；真机实麦冒烟独立复现通过；越界、隐私、mock 纪律全部干净。

**VERDICT: PASS**
