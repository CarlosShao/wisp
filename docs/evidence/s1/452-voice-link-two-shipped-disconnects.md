# 票 452 — 语音链交接之后的两处出厂断线（读数与裁定）

- 票：`.scratch/wisp/issues/452-voice-link-two-shipped-disconnects.md`
- 腿：writing leg（实现腿），2026-11-15，本机 `D:\work\workspace\projects plans\Wisp`
- 代码与本文同一个 commit（`git log -1 --oneline` 取 hash）；未 push。
- 本文件里的每条读数都是本机现跑的；不能现跑出来的那份，**明确写成"编不出来"**，不写成"测过"。

规范 §0（读数位置）：
- AC#1(a) 的四条编译形态读数：`.scratch/wisp/probes/452/ac1-a-measure.txt`（原文另抄在 §1）
- AC#1(b)/AC#3 的红→绿：§3、§5；gate 读数：§6
- AC#2 的裁定与代价：§2

---

## 1. AC#1(a)：出厂构建到底带不带 `-tags webrtcap`，`NewApm` 两种形态各返回什么

### 1.1 出厂构建行本身（`scripts/build.ps1:270`，逐字）

```
& $go.Source build -trimpath -ldflags $ldflags -o (Join-Path $outDir 'wisp.exe') ./cmd/wisp
```

`scripts/build.ps1` 全文件 `webrtcap` 命中数：**0**（`Select-String -Pattern webrtcap`）。这一行里也没有 `-tags`。

### 1.2 没打标签：`internal/audio/aec` 编译成什么（现跑）

```
$ go list -f 'GoFiles={{.GoFiles}}; CgoFiles={{.CgoFiles}}' ./internal/audio/aec
GoFiles=[aec.go aec_stub.go]; CgoFiles=[]
```

桩文件 `internal/audio/aec/aec_stub.go:20`：`func NewApm(rate int) (*Apm, error) { return nil, ErrUnavailable }`，
`ErrUnavailable`（`aec.go:76`）= `"aec: webrtc-audio-processing binding not built (requires cgo + -tags webrtcap; see deps.toml [webrtc-audio-processing])"`。

**没打标签的 NewApm 实测答案：`nil, ErrUnavailable`，永远是它。**
本条已用一枚随包测试钉住（`internal/audio/aec/aec_shape_452_test.go`，跑过）：

```
$ go test ./internal/audio/aec/ -run TestAC452 -v
=== RUN   TestAC452NewApmAnswersMatchTheBuildShape
--- PASS: TestAC452NewApmAnswersMatchTheBuildShape (0.00s)
ok  	github.com/CarlosShao/wisp/internal/audio/aec	0.053s
```

### 1.3 打了标签：本机**编不出来**（这就是"打了标签"的实测答案）

```
$ go list -tags webrtcap -f 'GoFiles=...; CgoFiles=...; Ignored=...' ./internal/audio/aec
GoFiles=[aec.go]; CgoFiles=[aec_webrtc_windows.go]; Ignored=[aec_stub.go]

$ go build -tags webrtcap ./internal/audio/aec
# github.com/CarlosShao/wisp/internal/audio/aec
wap_shim.cpp:23:10: fatal error: api/scoped_refptr.h: No such file or directory
   23 | #include "api/scoped_refptr.h"
      |          ^~~~~~~~~~~~~~~~~~~~~
compilation terminated.
rc=1
```

**打了标签的 NewApm 测不到**——不是"没测"，是**这个仓库在这台机器上没有可链接的 webrtc-audio-processing 树**：
`third_party/webrtc-audio-processing` 不存在，`deps.toml:61` 只登记了源码包（`webrtc-audio-processing-2.1.tar.xz` + abseil 子项目），
`scripts/fetch-deps.ps1` 只下载 deps.toml 里**预编译的 sherpa DLL**，没有任何构建 webrtc 的路径；meson/ninja 也不在 PATH 里。

打了标签、但 cgo 关掉，编译**通过**（标签单独不产生能力）：

```
$ CGO_ENABLED=0 go build -tags webrtcap ./internal/audio/aec
rc=0
```

⇒ 于是"打了标签的 NewApm 返回 `&Apm{...}, nil`"这一句只能引**源码**（`internal/audio/aec/aec_webrtc_windows.go:34-44`，
`wap_new` 成功即 `return &Apm{rate: rate, frame: rate / 100, sess: sess}, nil`），本文把它标成**源码事实、非本机读数**；
本机的可编译读数是 1.3 的编译失败。票面允许的"两条实测"＝**1.2 的现跑**＋**1.3 的编译失败**，两条都在这里。

### 1.4 这条断线的产品后果（AC#2 裁定的事实基础）

`cmd/wisp/resident_voice_bargein_windows.go` 的 `arm()` 在 AEC 构建失败时把原因写成
`"AEC3 不可用（本构建没有 -tags webrtcap 的绑定）：" + aerr.Error()`，即**开公告时才说**。
在此之前，面板的设置项 `SettingsView.Voice.AecEnabled`（`cmd/wisp/panel_config_store.go:102-120`）回答的是
**配置文件里 `[voice] aec.enabled` 的值**，不是这个二进制有没有这个能力。两者在出厂构建里恒不相等。

---

## 2. AC#2：裁定 — 选 ②「出厂不带标签；产品明说"本构建无回声消除"，界面判得了」

### 2.1 三个选项的代价（按"最有价值 + 最主流 harness"，不是"风险最小"）

| 选项 | 做法 | 真实代价 | 买到什么 |
|---|---|---|---|
| ① 出厂带标签 | `build.ps1` 加 `-tags webrtcap` | **标签单独无用**（§1.3：`CGO_ENABLED=0` 时加了标签仍是桩）。要真能力还必须：新写一条 webrtc-audio-processing 的**原生构建流水线**（源码包 → meson/ninja → MinGW 修 → cgo flags → DLL 进 `deps.toml` 的 SHA256 钉），本机没有 meson/ninja、没有树、`fetch-deps.ps1` 也没有这条路。且**今天加上去 = 唯一的出厂构建直接编不过**（§1.3 的 fatal error），发货件归零 | 出厂真的能 AEC（但要先有上面整条流水线） |
| **② 出厂不带；产品明说；界面判得了** | 代码里加机器可读的能力位 + 启动时说一句 + 用测试钉住；读数改标"专用构型" | 出厂**没有**回声消除，Path C 插话在出厂构建里**不可能发生**（半双工 D16 照旧）。恢复能力＝①，要走上面那条流水线 | 产品不再对用户撒谎；任何 UI/日志/读数都能一眼判"这枚二进制没有 AEC"；≤400ms 那类读数不许再当出厂成绩 |
| ③ 换不依赖标签的路径 | 换引擎/换算法 | 换引擎＝重做 389/392/401 三轮已验证的组合门（AEC、VAD、双讲判据全部要重测），且拿不到 ① 已有的 400ms 读数 | 不解决"出厂没有 AEC"，只是绕过标签；代价最大、收益最不确定 |

**裁定：②。** 理由一句话：**出厂构建今天就没有回声消除**（§1.2 是现跑，§1.3 是"连编都编不出来"），
① 在这一次交付里不是"多写一行标签"，而是要新建一条原生依赖流水线（并且今天加了就编不过），
③ 拿已验证的三轮读数去换一个不解决出厂事实的东西。② 是唯一**今天就能让出厂件说真话**的选项，
且它给 ① 留了口子：能力位一翻，同一句产品话术自动改成"可用"。

**⛔ 不做**（票面写法）：即"出厂不带标签、又把 ≤400ms 的读数当出厂成绩挂着"，本文与测试一起把它钉死（§2.3）。

### 2.2 ② 落地了什么（都在写面内）

1. `internal/audio/aec/aec_stub.go:17` `func BindingBuilt() bool { return false }` /
   `aec_webrtc_windows.go:29` `func BindingBuilt() bool { return true }`——**机器可读的能力位**，两条编译形态各一份，
   不用构造引擎就能回答"这枚二进制能不能消回声"。
2. `cmd/wisp/resident_voice_bargein_windows.go` 新增常量 `bargeinShippedPosture452`（产品原话，逐字）：

   > 本构建无回声消除（没有 -tags webrtcap 的 AEC3 绑定）：说话即插话（Path C barge-in）不可用，公告期间麦克风按 D16 半双工关闭；D47 的 400ms 以内插话读数只属于专用构型，不是出厂成绩

   `newVoiceBargeinMonitor` 在**构造时**（不是首次 arm 时）就把它写进 `reason`，
   于是任何在公告之前读姿态的地方都拿到真话。
3. 启动时印一次（`startResidentVoice`，`rv.start` 之前）：
   `wisp: 本构建无回声消除（…）：说话即插话（Path C barge-in）不可用…`；姿态行 `posture()` 新增 `插话=%s` 字段。
4. `internal/audio/aec/aec_shape_452_test.go` 钉住两种形态的 `NewApm` 答案与 `BindingBuilt()` 一致。
5. `cmd/wisp/resident_audio_device_lost_452_windows_test.go:TestAC452ShippedBuildSaysItHasNoEchoCancellation` 钉住：
   `build.ps1` 的 `go build` 行不含 `-tags`、全脚本 0 处 `webrtcap`、`NewApm` 是桩、`reason` 含"本构建无回声消除"。

### 2.3 "界面判得了"——诚实边界（不许夸大）

- **面板设置项（`AecEnabled`）今天答的是配置值，不是能力值**，且 `internal/panel/**` 与 `frontend/**` 都在本腿写面之外
  （`cmd/wisp/panel_config_store.go` 的 round-trip 已被 `panel_local_audio_r_test.go:102-103,145` 钉住）。
- 因此本腿交付的是**产品话术 + 机器可读位 + 启动读数**；**"面板上显示'本构建无回声消除'"这一半不在本腿写面内**，
  移交给面板/前端腿的落点是：面板 `SettingsView.Voice` 增加一个**能力字段**（如 `AecBuildAvailable bool`），
  数据源直接用 `aec.BindingBuilt()`，与用户可改的 `AecEnabled` **并列且分开**（一个说"这枚二进制能不能"，一个说"你想不想开"）。
  在它落地之前，能力事实在**启动 stdout + 姿态行 + 日志**三处都已可判。
- **"球与面板都离开 Listening"**：D43 状态只有球带着（`residentBall.setVoiceState` 是唯一非测试落点），
  面板包（`panel packet`）里今天**没有** D43 状态字段 ⇒ 球离开 Listening 是事实，面板侧要等上面那个字段一起做。
  本文不写"面板也离开了"。

> **2026-11-15 14:5x 补记（第二腿）**：上面第一条与第二条的那个"写面外"缺口**已经补上**——`SettingsView.Voice` 的
> `AecBuildAvailable`（`*bool`，数据源 `aec.BindingBuilt()`）已落，`cmd/wisp` 组装处逐行填真值，并有**会红的线材测试**
> 与**前端契约**（字段名/语义）。落点、突变读数、门禁与前端契约见本文 **§9**。第三条（面板包里的 D43 状态字段）
> 仍**未做**，不因本腿而变。

---

## 3. AC#1(b)：改代码**之前**的两条红句（逐字，从载具日志抄）

载具：`bash scripts/wisp-cli-staged-tests.sh -test.run=^TestAC452`（`go test -c ./cmd/wisp/` + 暂存 sherpa DLL）。
改动前跑出：

```
=== RUN   TestAC452ShippedBuildSaysItHasNoEchoCancellation
    resident_audio_device_lost_452_windows_test.go:136: 452 AC#1(a)/AC#2 RED: this build has no AEC3 binding (aec.NewApm answers ErrUnavailable, scripts/build.ps1 passes no -tags webrtcap), and the barge-in monitor's posture at construction is "": the product states neither that this binary cannot cancel echo nor that Path-C barge-in is therefore unavailable. A build without the capability may not stay silent about it (AC#2 form ②), and the 400ms-class readings belong to the dedicated -tags webrtcap staging, not to what ships
--- FAIL: TestAC452ShippedBuildSaysItHasNoEchoCancellation (0.00s)
=== RUN   TestAC452LostMicWalksD43FromListeningIntoError
    resident_audio_device_lost_452_windows_test.go:204: 452 AC#1(b) RED: the microphone was unplugged after the device was handed over (mic.Err() now answers "audio_device: capture device lost on device USB Mic 452 (0x88890004): device invalidated (unplugged or disabled)") while the machine sat in Listening, and 3s later the state is still Listening. Nothing in the production tree reads the capture leg's terminal error, so D43#14 (Listening + EvAudioDeviceLost -> Error, table.go:97) has zero dispatchers and the ball stays in Listening with a microphone that no longer exists. Ball states: []; problems: []
--- FAIL: TestAC452LostMicWalksD43FromListeningIntoError (3.42s)
FAIL	github.com/CarlosShao/wisp/cmd/wisp	4.00s
```

两句红句的两个要点：
1. **姿态在构造时是空串 `""`** ⇒ 出厂构建对"没有回声消除"一言不发（AC#2 ② 的形状被否定）。
2. **拔掉麦克风 3 秒后机器还是 Listening，球的投影历史 `[]`（一次都没推过）** ⇒ D43#14 在出厂树里**零派发者**。

---

## 4. AC#3：把 `EvAudioDeviceLost` 接上派发（怎么接的，为什么这么接）

改动点（全在写面内：`internal/audio/**`、`cmd/wisp/**`）：

1. **粘性终态读**：`cmd/wisp/resident_audio_windows.go` 新增 `residentAudio.captureErr()`。
   源头的 `Err()` 是**破坏性一次读**（`wasapimic_windows.go` 的 `errC` 是一格通道，`postErr` 非阻塞写），
   谁先读谁拿走 ⇒ 由采集腿缓存首个非 nil 答案，之后所有读者（健康探针、静音手势的句子、退出行）答同一个事实。
   `toggleMute` 与启动 switch 里原先的两处 `mic.Err()` 改走它（否则探针会把错误"吃掉"）。
2. **探针**：`resident_voice_windows.go` 新增 `checkAudioDevice()` + 常量 `voiceDeviceHealthInterval = 250 * time.Millisecond`。
   在 `rv.run()` 的两个主 select（阶段一、阶段二）各加一条 `case <-health.C: rv.checkAudioDevice()`。
   **不新增 goroutine**（`internal/observe` 名册不在写面内，`TemporaryNames` 里没有可用名）：ticker 骑在既有的 owner goroutine 上，
   派发也在这条 owner 上，和本腿其它 D43 走动同一条线程。
3. **派发守表**：只有 `StateListening` 才派发（`table.go:97` 的 #14 只有一个 From）。
   其它状态只记一次事实（`deviceLostNoted`），ticker 继续问——公告播完回到 Listening 时再走。
   走过一次后 `deviceLostDelivered` 锁死，不重复派发。
4. **接线**：`startResidentVoice` 里 `rv.deviceErr = ra.captureErr`，**在 `rv.start()` 之前**——
   这既是 happens-before 边（消费侧只读），也让第一次 tick 就知道这次启动有没有交接设备。
   测试对这行做了**生产行存在性断言**（见 §5 的 M2），防止"测试自己接线、生产没人接"的空绿。
5. **恢复路径不在本片**：重新插回/换设备之后球停在 Error，等 `EvErrorAck`（D43#38）。这句话写进了 `noteProblem` 的句子和球的状态里，不假装有自动恢复。

`resident_voice_bargein_windows.go` 的 AEC 缺席句、`resident_audio_windows.go` 的启动句子都按 AC#4 重写（§5.3）。

---

## 5. 绿读数 + 突变

### 5.1 绿（载具，`-test.run=^TestAC452`）

```
=== RUN   TestAC452ShippedBuildSaysItHasNoEchoCancellation
--- PASS: TestAC452ShippedBuildSaysItHasNoEchoCancellation (0.00s)
=== RUN   TestAC452LostMicWalksD43FromListeningIntoError
time=... level=ERROR msg="voice: the capture device died after the handover" err="audio_device: capture device lost on device USB Mic 452 (0x88890004): device invalidated (unplugged or disabled)" state=Listening
time=... level=INFO msg="voice: D43 walked" from=Listening event=audio.device-lost to=Error why="麦克风在设备交接之后丢失（…）：D43#14 Listening → Error，球离开 Listening"
--- PASS: TestAC452LostMicWalksD43FromListeningIntoError (0.53s)
ok  	github.com/CarlosShao/wisp/cmd/wisp	1.00s
```

绿的三件事同时成立：机器走 `Listening + audio.device-lost → Error`；`timeline` 里有这条 From/To；
球的 `voiceStatesSnapshot()` 末项是 `Error`。健康期（一个探针间隔 400ms）**没有**误报（机器仍是 Listening）。

### 5.2 突变（哪一个测试变红）——两次，逐字

**M1：把 `checkAudioDevice` 里的派发关掉**（`if false && rv.dispatch(...)`，只此一处）：

```
--- FAIL: TestAC452LostMicWalksD43FromListeningIntoError (3.43s)
    resident_audio_device_lost_452_windows_test.go:218: 452 AC#1(b) RED: … 3s later the state is still Listening … D43#14 … has zero dispatchers … Ball states: []; problems: []
```

⇒ 行为测试**真的**盯着派发，不是盯着一个常量。

**M2：删掉 `startResidentVoice` 里的生产接线行 `rv.deviceErr = ra.captureErr`**：

```
--- FAIL: TestAC452LostMicWalksD43FromListeningIntoError (0.02s)
    resident_audio_device_lost_452_windows_test.go:183: 452 AC#3 RED: resident_voice_windows.go has no `rv.deviceErr = ra.captureErr` line, so no boot actually hands the health probe a capture leg: the machine below would only walk because the TEST wired the seam
```

⇒ 被人删掉生产接线时，这枚测试**不会**因为"测试自己接了线"而继续绿。

两次突变都当场还原（`git diff` 里没有 `if false`、`rv.deviceErr = ra.captureErr` 计数为 1），还原后 §5.1 复跑仍绿。

### 5.3 AC#4：那句已经不成立的注释（改前 / 改后）

改前（`cmd/wisp/resident_audio_windows.go` 启动分支的 `fmt.Printf`，逐字）：

```
wisp: 这一条不推状态、也不拒绝启动：D43 的 EvAudioDeviceLost 只有从 Listening 出发的合法边，启动期没有；票 128 只定了「没有数据根」这一种拒绝启动
```

改后：

```
wisp: 启动期这一条不推状态、也不拒绝启动：D43 的 EvAudioDeviceLost 只有从 Listening 出发的合法边，而启动期的机器还在 Sleeping（票 128 只定了「没有数据根」这一种拒绝启动）；运行期的设备丢失已接派发（票 452：设备交接之后丢失，语音链的健康探针走 D43#14 Listening → Error，球离开 Listening）
```

文件头那条 "It pushes NO state" 同样改成"**启动期**不推状态"＋一段说清运行期那一半（`resident_audio_windows.go:34-45`），
并把"这一段写在 452 之前"的事实留在注释里，避免下一个人再照着旧句读。

---

## 6. AC#5 gate 读数（本机现跑）

| gate | 命令 | 读数 |
|---|---|---|
| 格式 | `gofumpt -l cmd/wisp/ internal/audio/aec/` | 空（先出一枚 `resident_voice_bargein_windows.go` 多余空行，`-w` 后空） |
| vet | `go vet ./internal/audio/aec/ ./cmd/wisp/` | `rc=0` |
| D22 扫描 | `bash scripts/d22scan.sh` | `OK - control rc=0, scan rc=0`；`cmd/ 292 Go files, comments and _test.go included` 全扫，0 违规 |
| 路径长度 | `bash scripts/check-path-length-budget.sh` | `VERDICT GREEN`，`not in roster=0` |
| 单元 | `go test ./internal/audio ./internal/statemachine` | `ok internal/audio 18.252s`、`ok internal/statemachine 0.263s` |
| 新读数 | `go test ./internal/audio/aec -run TestAC452` | `ok 0.053s` |
| 整包 | `bash scripts/wisp-cli-staged-tests.sh` | 见 §6.1 |

### 6.1 整包（`./cmd/wisp`，载具一次一发）

两次（含最终树）：`626 PASS / 17 SKIP / 2 FAIL`，两枚红都是已知本机红（AC4 gap33r5、AC246 dev 实例），
第三条已知红 320 两次都过；未被 native AV 终止。读数原文与 FAIL 行：§8。

---

## 7. 本腿**没做**的事（写明，不假装）

1. **重新插回/换设备之后的恢复**：不在本片。球停在 Error 等 `EvErrorAck`（D43#38），`noteProblem` 的句子明写。
2. **面板能力字段**（§2.3）：`internal/panel/**`、`frontend/**` 出写面 ⇒ 移交前端/面板腿，落点已写明。
3. **`internal/audio/aecspike` 的整包测试**在本机 `exit status 0xc0000135`（缺 sherpa DLL 的加载期失败，非本腿改动；不在 AC#5 的三条 gate 里）。
   本腿只改了 `internal/audio/aec`（`go test ./internal/audio/aec` 绿）。
4. **带 `-tags webrtcap` 的构建**在本机编不出来（§1.3）⇒ 本腿**没有**交付"AEC 可用"的读数，只有"出厂没有 AEC"的读数。

---

## 8. 整包实测（`./cmd/wisp`，载具一次一发，跑了两次）

摘要件：`.scratch/wisp/probes/452/cmdwisp-full-summary.txt`（两次的 FAIL 行、PASS/SKIP 计数、本票两枚测试的行）。
原始日志（本机，未入库）：`.scratch/wisp/probes/452/cmdwisp-full.log`（第一次，1.28MB）、`…-full-run2.log`（第二次，1.28MB）。

| 次 | 起跑 | PASS | SKIP | FAIL | 末行 |
|---|---|---|---|---|---|
| 1（注释修订前的树） | 13:47 | 626 | 17 | 2 | `FAIL github.com/CarlosShao/wisp/cmd/wisp 1341.00s` |
| 2（**最终树**，含两处注释修订） | 14:14 | 626 | 17 | 2 | `FAIL github.com/CarlosShao/wisp/cmd/wisp 1190.00s` |

两次的 FAIL 都**只有**这两枚，且都在票面/派单预告的已知本机红清单里：

```
--- FAIL: TestAC4FocusReturnToPriorWindowGap33r5 (4.87s)
    panel_host_windows_test.go:1124: the panel did not take the foreground on Show: foreground 0x5701544, panel hwnd 0x76207cc (AC#4 says only the panel takes focus when shown)
--- FAIL: TestAC246DevLegIgnoresTheTestTaskInjection (37.20s)
    resident_task_source_246_windows_test.go:484: AC#7 RED (desktop state, not our code): a dev Wisp is already running in this session, so this leg handed its activation over and exited. Stop it and re-run.
```

- 第三条已知红 `Test320FetchLegWritesTheKeyNowhereOnDisk` 两次都**通过**（它扫 `%TEMP%`，这次没被别人的临时件命中）。
- 两次都**没有被 native AV 终止**（`0xC0000005` 未出现）：整包跑到了结论行。
- 本票两枚测试两次都是 `--- PASS`（第二次 `0.00s` / `0.52s`）。


---

## 9. AC#2 的"界面可判"那一半（第二腿，2026-11-15 14:4x–15:0x）——面板字段 + 组装接线 + 会红的线材测试 + 前端契约

§2.3 与 §7 第 2 条把这一半划在写面外；这一腿把它落了，**不删不改上面那两处原文**，只在这里写落到哪儿。
（同票面"收件"一节：这是 AC#2 剩下的第 1 件；第 2 件＝前端界面，仍归机主外部 agent，本腿**零字节**碰 `frontend/**`。）

### 9.1 两个布尔，各答一件事（这三条是判据，不是描述）

| 字段 | 类型 | 答什么 | 谁写 | 能不能改 |
|---|---|---|---|---|
| `VoiceSettingsView.AecEnabled` | `*bool`（既有） | **用户想不想开**：`[voice.aec] enabled` 的配置值 | `cmd/wisp/panel_config_store.go` 从 `cfg.Voice.AEC.Enabled` 读 | **可写**（`config.set` 的 `voice_aec_enabled`） |
| `VoiceSettingsView.AecBuildAvailable` | `*bool`（**本腿新增**） | **这枚二进制能不能消回声**：`aec.BindingBuilt()` 的编译期事实 | 同一个 `ReadSettings`，从 `aec.BindingBuilt()` 读 | **只读**：`config.set` 名册里**没有**它，也不该有 |

- **⛔ 合成一个布尔的后果**，逐字写清：出厂构建（`scripts/build.ps1` 不传 `-tags webrtcap`）**不能**消回声，而
  `[voice.aec] enabled` 的默认值是 **on** ⇒ 只有一个布尔时页面会印"回声消除：开"，而那枚二进制一次也消不了。
  这就是票面 AC#2 要杀的那句假话。
- **指针仍是闸**（与同结构里三枚 `本机音频` 字段同一条口径）：`nil` ⇒ 这条腿没答，JSON 键**根本不在线上**；
  `&false` ⇒ **答了**"本构建没有"；`&true` ⇒ 答了"有"。所以"没人问"和"问了、答案是 false"在线上不可互换。
- **与启动印同源**：面板这一格读的就是构造期姿态与启动印那一行读的**同一个** `aec.BindingBuilt()`
  （`cmd/wisp/resident_voice_bargein_windows.go:332-333` 构造期写 `reason`、`:349-354` `bargeinAECPosture452()`
  答姿态行）⇒ 页面、姿态行、启动 stdout 三处**不可能分叉**：只要构建形态没变，它们说的是同一句真话。
- `ReadSettings` 的两个早退（`mgr == nil` / `cfg == nil`）**不填**这个字段：那两条是"这条腿没有配置可读"，
  它们的 view 里 `Readable=false`，页面整体进等待态；而那两条路径在本进程里没有可达构造（`Manager.Config()` 恒返回
  `deepcopy(m.cur)`，`m.cur` 在 `internal/config` 里除回滚外不会被置 nil）⇒ **不写一条测不到的赋值**
  （A1880：摘掉它 0 枚红的线＝装饰线）。这一点具名在此，不靠读者猜。

### 9.2 接线点（`file:line`，现读）

- **产码源头**：`cmd/wisp/panel_config_store.go:117` `aecBuildAvailable := aec.BindingBuilt()`
  （local 与同函数三枚既有 local 同因：view 比这次调用活得久，交出去的是快照不是 per-call 的地址）
- **赋值点（就是它）**：`cmd/wisp/panel_config_store.go:137` `AecBuildAvailable: &aecBuildAvailable,`
  在 `ReadSettings` 的 `Voice: panel.VoiceSettingsView{...}` 组装块里（与 `AecEnabled`/`InputDevice`/`SampleRate` 并列）
- **契约字段**：`internal/panel/config_handlers.go:236`（`AecBuildAvailable *bool`，注释在 `:217-235`）
- **上线（JSON）**：`internal/panel/config_handlers.go:670`（view→payload 透传）＋ `:735`
  `AecBuildAvailable *bool` 带 tag `json:"aecBuildAvailable,omitempty"`

### 9.3 会红的测试：两次突变，红句逐字

**测的是线上的字节，不是源码里的字符串**（A1880 的判据）：测试
`cmd/wisp/panel_voice_aec_452_test.go:TestAC452PanelVoiceTellsThePageWhatTheBuildCanDo` 用
`newRealSettingsLeg` 起**真**腿（真 `config.Manager` + 真 `secret.Store` + 真 `panel.ComposerDispatch`），
走**门**发 `config.get`／`config.set`，读回页面真拿到的那行 JSON。判据是
`voice.aecBuildAvailable` **在线上且等于 `aec.BindingBuilt()`**（⛔ 不是断言字面 `false`——零值会让它在静默里绿）。

**M1：摘掉那格赋值**（`AecBuildAvailable: &aecBuildAvailable,`）。Go 不接受"局部变量/导入声明了不用"，
所以这一发连 `:117` 的 local 与 `:43` 的 import 一起摘（三处，全在 `panel_config_store.go`），`go build ./cmd/wisp` rc=0 后才跑：

```
    panel_voice_aec_452_test.go:66: the settings payload's voice half carries NO aecBuildAvailable key: the page cannot tell that this build cannot cancel echo, so it prints the user's switch as if the capability existed - which is ticket 452 AC#2's undone half (voice = map[aecEnabled:false cloudASRChain:<nil> cloudTTSChain:<nil> inputDevice:default realtimeEnabled:false realtimeModel: realtimeProvider: sampleRate:16000])
--- FAIL: TestAC452PanelVoiceTellsThePageWhatTheBuildCanDo (0.02s)
FAIL  github.com/CarlosShao/wisp/cmd/wisp	0.00s
```

**M2：把能力位折进用户开关**（只改一行，`:117` 变成
`aecBuildAvailable := cfg.Voice.AEC.Enabled || aec.BindingBuilt()`，编译 rc=0）：

```
    panel_voice_aec_452_test.go:86: voice.aecBuildAvailable = true (present=true) after the user switched AEC on, want false: the capability bit is being derived from the user's switch instead of the build, which is the one-boolean shape ticket 452 forbids
--- FAIL: TestAC452PanelVoiceTellsThePageWhatTheBuildCanDo (0.03s)
```

⇒ 两发打的是**不同半边**：M1 打"键不在线上"（缺接线），M2 打"值跟着用户开关走"（合成一个布尔）。
★**同时点名一件不遮的事**：`internal/panel` 侧那枚在册的键集冻结
（`Test313PayloadKeySetIsFrozenAndCarriesNoValueBearingKey`）在 M1 下**照样绿**——它直接构造 view 交给渲染器，
看不见组装有没有接线。⇒ 光靠它＝A1880 说的"盘上字面钉"，所以**牙必须长在 `cmd/wisp` 的线材测试上**，本腿就是这么放的。

两次突变都当场还原：`panel_config_store.go` 前后 sha256 同为
`EB28195ABA18D7D755CCE160ABFCF73DBA4C1FF44849BEF6E1582C70CDB78B62`，还原后 `go build ./cmd/wisp` rc=0，
定向两发复跑全绿（见 §9.4）。备份件 `.scratch/wisp/probes/452/panel-half/panel_config_store.go.pristine`。

### 9.4 绿读数

| 面 | 命令 | 读数 |
|---|---|---|
| 面板包全量 | `go test ./internal/panel -count=1` | `ok github.com/CarlosShao/wisp/internal/panel 28.787s` |
| 面板侧新钉 | `go test ./internal/panel -run 'Test452\|Test313' -count=1` | `ok 0.072s` |
| 线材测试（载具） | `bash scripts/wisp-cli-staged-tests.sh -test.run '^TestAC452PanelVoiceTellsThePageWhatTheBuildCanDo$'` | `--- PASS: … (0.04s)`，`ok github.com/CarlosShao/wisp/cmd/wisp 0.00s`，rc=0 |
| 452 家族（载具） | `bash scripts/wisp-cli-staged-tests.sh -test.run '^TestAC452'` | 三枚全绿：`TestAC452ShippedBuildSaysItHasNoEchoCancellation` / `TestAC452LostMicWalksD43FromListeningIntoError (0.52s)` / 本腿那枚 |

### 9.5 前端契约（给机主外部 agent；本腿 ⛔ 碰 `frontend/**`）

前端要显示的东西，逐字如下，**照这四行做即可**：

1. **字段名**：`voice.aecBuildAvailable`——JSON **布尔**，在 `config.get` 回执里
   `wisp-config-view:` 标记行之后那串 JSON 的 `voice` 对象里（与 `voice.aecEnabled` 同一个对象）。
2. **语义**：`true` ＝ 这枚二进制的构建里**带** AEC3 绑定（能消回声）；`false` ＝ **本构建没有**回声消除
   （出厂形态：`scripts/build.ps1` 的 `go build` 不传 `-tags webrtcap`）；**键缺席** ＝ 这条设置腿没回答
   （生产里不该出现：`cmd/wisp` 的 `ReadSettings` 每次都填；出现即腿侧缺陷，页面按等待态渲染，⛔ 不要当成 false）。
3. **与 `voice.aecEnabled` 的关系**：**并列、独立、不可合并**。`aecEnabled` 是用户开关（可写；写字段
   `config.set` 的 `voice_aec_enabled`），`aecBuildAvailable` 是二进制能力（**只读**，`config.set` 名册里没有它、
   ⛔ 不要给它造写路）。两者可以同时是"用户要开 + 构建不能"，那正是出厂态。
4. **页面规则一句话**：`aecBuildAvailable === false` 时，回声消除这一栏必须显示"**本构建无回声消除**"
   并让开关表现为**打不开/无效**（⛔ 不能只显示 `aecEnabled: true` 就当能力存在）；`true` 时正常渲染开关；
   键缺席时两样都别断言。

**为什么不用改前端就能先跑起来**：`frontend/src/lib/config-view.ts:105-108` 逐字写着 unknown fields are
ignored on purpose，且本仓 `*.ts` 里今天 `aecEnabled` 的读者是 0 命中 ⇒ 新键对页侧是**容错**的（不破 §2.3 那句
"payload 新增键对 config-view.ts 是容错的"）。**页面那一半没做，本腿不声称做完。**

### 9.6 门禁（本机现跑）

| gate | 命令 | 读数 |
|---|---|---|
| 格式 | `gofumpt -l internal/panel cmd/wisp`（`D:\work\base\gopath\bin\gofumpt.exe`） | **空** rc=0（起手有一枚 `config_handlers.go` 对齐，`-w` 后空） |
| vet | `go vet ./internal/panel ./cmd/wisp` | **rc=0** |
| build | `go build ./internal/panel ./cmd/wisp` | **rc=0** |
| D22 扫描 | `bash scripts/d22scan.sh` | `OK - control rc=0, scan rc=0`；`cmd/ 293 Go files`、`internal/ 1019`，0 违规 |
| 路径长度 | `bash scripts/check-path-length-budget.sh` | **VERDICT GREEN**，`over the hat=68` / `roster=68` / not in roster=0 |
| 跨平台 | `GOOS=linux go build ./internal/audio/aec ./internal/panel` | **rc=0**（新 import 不引入平台面）。`GOOS=linux go build ./cmd/wisp` **rc=1**，但红源是既有事：`github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files`（`scripts/wisp-cli-tests.sh` 头部早记过），⛔ 本腿新增 |

### 9.7 整包（`./cmd/wisp`，载具一次一发，两发）

命令（同一个载具，串行两发，⛔ 并行）：

```
D:\work\soft\Git\bin\bash.exe -c "cd '/d/work/workspace/projects plans/Wisp' && bash scripts/wisp-cli-staged-tests.sh"
```

| 发 | 起 / 末 | 载具判行 | 顶层 RUN | PASS | FAIL | SKIP | 没拿到结果的那一枚 |
|---|---|---|---|---|---|---|---|
| 1 | 14:49:25 / 14:56:41 | `FAIL … 437.00s` | 339 | 338 | 0 | 0 | `TestAC4PriorFocusSurvivesARefusedPanelSample` |
| 2 | 14:59:49 / 15:06:13 | `FAIL … 386.00s`，exit 1 | 339 | 337 | 1 | 0 | 同上（同一枚） |

**两发都死在票 450 在册的那个形状上**（整进程消失、无 `--- FAIL`、最后一条 `=== RUN` 就是
`TestAC4PriorFocusSurvivesARefusedPanelSample`，437s／386s）：票 450 的现量表逐字记过同一枚
（`PASS 334 / FAIL 0 / 403s`，死在同一枚用例半途）。**⛔ 不报绿**：这不是"整包绿"，是"整包在本机跑不完"。

- 第 2 发那 1 枚 FAIL 是**本机已知红**：`TestAC4FocusReturnToPriorWindowGap33r5 (5.31s)`（派单点名的三枚之一，非本腿）。
  另两枚已知红（`TestAC246DevLeg…`／`Test320Fetch…`）在死点**之后**的文件序里 ⇒ 这两发**根本没跑到**，既不算绿也不算红。
- ★**必须点名的一条覆盖事实**：`cmd/wisp/panel_voice_aec_452_test.go` 按文件序排在那枚死点**之后**，
  所以**两发整包对本腿这枚新测试的覆盖＝0**。本腿的绿读数是 §9.4 那一发**定向载具**（`-test.run '^TestAC452…'`，
  同一个载具、同一份二进制构建），不是整包里的某一行。⛔ 不拿"整包 FAIL 0/1"当覆盖证明——这正是票 450
  存在的理由（"看起来绿、实际没跑"）。
- 原始日志（本机、未入库）：`.scratch/wisp/probes/452/panel-half/cmdwisp-full-452panel.log`（第 1 发）、
  `…-run2.log`（第 2 发）；突变红句 `mut-m1-red.txt`／`mut-m2-red.txt`；突变说明 `notes-mutations.md`。

### 9.8 本腿**没做**的（写明，不假装）

1. **前端界面**：`frontend/**` 零字节改动，只有 §9.5 的契约。页面今天**仍不会**显示"本构建无回声消除"。
2. **面板包里的 D43 状态字段**（§2.3 第三条）：仍未做，球离开 Listening 是事实、面板侧还看不见。
3. **`ReadSettings` 两个早退路径不填能力位**：理由见 §9.1 末段（那两条在本进程里测不到，写了就是装饰线）。
