# 票 290 AC#3 两形真机读数（实现腿 `290-ac3-live`）— 默认档一字不改：形① 出厂档不开设备 ＋ 形② 用户开门之后球拿到电平

**性质**：这是**实现腿交付的读数件**，不是裁决表（`docs/evidence/s1/` 的裁决权在非实现者验收腿/编排者）。本件只放**现跑读数与逐字引文**，⛔ 不翻 AC#3 的框，⛔ 不改票面一字，⛔ 零 push。
**钟点**：2026-11-14 10:16–10:41 +08（`time` 现读）。**树锚**：主仓 `7bdf09042`（branch `dev`）。
**仪器**：`cmd/wisp/resident_audio_290_live_windows_test.go`（`//go:build windows && winlive`，本腿新建，三枚用例）。它复用仓里既有 carrier：`buildWispForTest`（`cmd/wisp/secret_argv_windows_test.go:161`）、`bootResidentLeg`/`pollUntil127`/`readSinkTail258`（`resident_sink_nail_127_windows_test.go:224`/`:324`、`resident_hotkey_live_258_windows_test.go:296`）、`findBallWindows228`/`live228R3*`（`resident_ball_live_228_windows_test.go:60`、`resident_tray_stop_live_228_windows_test.go`）、票 247 的台件家族 `bootAudioRuntime`/`writeAudioConfig`/`assembleCapture`/`livePhase`/`phaseRouter`/`requireNoCompetingAudioProcesses`。
**原始件**：`.scratch/wisp/probes/290/live/raw-live-1.log`（10:24 第一次两形现跑，形②成功那一发）·`raw-live-2.log`/`raw-live-3.log`（形②输入投递诊断）·`raw-live-4.log`（终态三枚用例一发，rc=0）。仪器源拷贝：`probes/290/live/resident_audio_290_live_windows_test.go`。

## 起手三核（现跑，⛔ 信转述）

1. **两把 tasklist 尺＝0**（10:16 现读）：`tasklist /FI "IMAGENAME eq wisp.exe"` ⇒ `INFO: No tasks are running`；`balldebug.exe` 同句。跑完（10:41）复读仍为 0，无孤儿。
2. **默认档前提现读**：机主自己那份 `%APPDATA%\wisp-dev\config.toml:69` 逐字 `mic_muted_default = true`；`[voice] enabled = true`（`:33`）、`wake_word.enabled = false`（`:40`）。编译默认表 `internal/config/schema.go:299` 逐字 `MicMutedDefault bool \`toml:"mic_muted_default" default:"true"\``；装配 `cmd/wisp/resident_audio_windows.go:314` 逐字 `gate := audio.NewHalfDuplexGate(mic, audio.PathT, audio.WithStartMuted(c.Audio.MicMutedDefault))`。
3. **默认值一字未动**：`go test ./cmd/wisp -run 'TestAC247ShippedDefaultsAreTheOnesThisLegReads|TestAC290' -count=1` ⇒ **PASS 7 枚 / 0 红，rc=0**（含三枚默认值现读＋既有六枚 290 行为用例）；`gofumpt -l cmd/wisp/resident_audio_290_live_windows_test.go` ⇒ 空。

## 形① — 默认双击启动（出厂默认档：`mic_muted_default` 仍 true、装配走 `WithStartMuted(true)`、verdict 仍打「设备未打开」）

仪器＝**真发运 exe**（`go build ./cmd/wisp` 的产物，无参数＝双击那条命令行），data dir＝空临时目录（**没有 config.toml** ⇒ 编译默认表决定姿态，等价于"第一次双击"），10:33:2x（`raw-live-4.log`）：

| # | 逐字读数（出处行） |
|---|---|
| 1 | `wisp: 采集腿：采集腿已装配、门处于静音：设备未打开（[audio] mic_muted_default=true，来源 compiled defaults (config.NewDefaults)）` —— `raw-live-4.log` 的 `resident_audio_290_live_windows_test.go:307` 那行；格式串本体 `cmd/wisp/resident_audio_windows.go:339` |
| 2 | `wisp: resident event loop running (... hotkeys live 3/4: summon=Ctrl+Alt+Q live, mute=Ctrl+Alt+M live, ...)` —— `raw-live-4.log:308`（同一句里给出 mute 槽真注册上了，形②的前提） |
| 3 | `wisp: audio capture leg stopped: levels_delivered=0 frames_sent=0 frames_dropped=0 reopens=0` —— `raw-live-4.log:325`；该行格式串 `resident_audio_windows.go:429` |
| 4 | `{"time":"2026-11-14T10:33:27.7914548+08:00","level":"INFO","msg":"audio: capture armed but muted at boot; device not opened","config_source":"compiled defaults (config.NewDefaults)","path":"T"}` —— 盘上日志 sink 记录，`raw-live-4.log:330` |

**装配级三判据（同一棵树、同一发 `-run TestAC290Live` 的第三枚用例里现读，rc=0）**：`ra.mutedAtBoot==true`（＝交进 `WithStartMuted` 的那一枚值）／`ra.gate.Muted()==true && ra.gate.Open()==false && ra.started.Load()==false`／`rt.Registry.CountByName("audio-capture")==0`（采集线程一枚都没起）／开机窗内球缝收到 0 枚标量。**`frames_sent=0` 是这一格的硬读数**：设备压根没被交接，所以一帧 PCM 都不存在。

## 形② — 用户主动开门之后（闸放行、球拿到电平）

### 2a 发运进程上的开门读数（真机、真手势、进程自己打的计数）

10:24:39–40（`raw-live-1.log`，`-run 'TestAC290Live'` 第一发；该发后续因为我自己写错的一条断言被判红，**读数本身与断言无关**，见下"具名事故"）：

| # | 逐字读数（出处行） |
|---|---|
| 1 | 开门前同形①：`wisp: 采集腿：采集腿已装配、门处于静音：设备未打开（[audio] mic_muted_default=true，来源 compiled defaults (config.NewDefaults)）` —— `raw-live-1.log` 的 `:222` 行 |
| 2 | **手势之后（球侧热键→门）**：`wisp: ball mute-hotkey: 已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）` —— `raw-live-1.log:243` |
| 3 | **进程自己的退场计数**：`wisp: audio capture leg stopped: levels_delivered=129 frames_sent=6 frames_dropped=123 reopens=0` —— `raw-live-1.log:261` |

`levels_delivered` 的产码路径逐字：`resident_audio_windows.go:148 ra.levels.Add(1)` → `:165 rb.b.SetAudioLevel(level)`（球缝就是 `internal/ball` 的 `SetAudioLevel`；`resident_windows.go:657 rb.attachMuteGate(raudio.toggleMute)` 是手势执行者的注入点，`:182 func (ra *residentAudio) toggleMute()` 是它本体）。⇒ **用户开门之后，球缝真拿到 129 枚电平标量**，出处＝发运进程自己的第 4 步退场行。

### 2b 同一路径上的逐值读数（生产装配＋生产执行者＋真麦）

10:33:5x（`raw-live-4.log`，第三枚用例 PASS 5.41s）。装配＝`assembleCapture`（真 `WASAPIMicrophone`＋真 `HalfDuplexGate`），config＝`writeAudioConfig(t, nil)`（**出厂默认表，`mic_muted_default` 原值 true，⛔ 没有 `MicMutedDefault=false` 那枚临时配置**），开门＝`residentAudio.toggleMute()`（球侧热键与托盘静音**两支手势共用的那一个执行者**，⛔ 不是改配置）：

| # | 逐字读数（出处行） |
|---|---|
| 1 | `AC#3 assembly verdict: 采集腿已装配、门处于静音：设备未打开（[audio] mic_muted_default=true，来源 …\config.toml）` —— `raw-live-4.log:505` |
| 2 | `AC#3 gesture open: executed=true outcome="已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）"` —— `:524` |
| 3 | `AC#3 counters after open: levels_delivered=93 frames_sent=6 frames_dropped=87 last_err=""` —— `:544` |
| 4 | `AC#3 scalars after open: after-user-open: samples=93 mean=0.003697 max=0.057523 levels_per_s=31.00` —— `:546` |
| 5 | `AC#3 gesture close: executed=true outcome="已静音：采集已关闭，设备未打开（再按一次取消静音）"` —— `:555` |
| 6 | `AC#3 scalars after close: after-close: samples=0 mean=0.000000 max=0.000000 levels_per_s=0.00` —— `:564` |

断言同发全绿（rc=0）：开门前 `muted/open/started`＝true/false/false、`audio-capture` 线程 0、球缝 0 枚；开门后 `gate.Open()==true && gate.Muted()==false`、线程 1；3s 窗 93 枚（31.00/s，与 P8 甲的 32ms 帧节奏一致）；关回后**再收 0 枚**、线程回到 0；`shut()` 的 D38(e) 记录 10 枚。

**⛔ 不声称"随声音变化"**：mean 0.003697／max 0.057523 是**安静房间的底噪**，不是"说话时升高"。这台机 2026-11-08 的麦侧四形读数（`probes/291/mic/00-four-shape-readings.md:24-41`）已经量过：安静 mean 0.003067、外放音调麦读 max 0.002530＝扬声器→麦通路无效；且尺不剥直流是另一个票（297）的定义面。⇒ 本件只声称**闸放行＋球缝收到标量**（AC#3 的"采样数与出处"那一半），"数随人声变"要机主在场对麦说话，归那一发。

### 2c 具名事故与残余（⛔ 不遮）

1. **形② 在发运进程上的可复现性：1 成 / 4 败**。10:24 那一发落地（上表 2a）；随后 10:26、10:29、10:31、10:33 四发，同一个 chord（`keybd_event` 注入 Ctrl+Alt+M）都不再落到进程里。诊断读数（`raw-live-3.log:375-417`、`raw-live-4.log:378-425`）：
   - `AC#3 shape2 input desktop before the chord: Default`（不是锁屏的 Winlogon）；
   - `AC#3 shape2 chord ownership: held by another process (RegisterHotKey answered ERROR_HOTKEY_ALREADY_REGISTERED while the child is alive)`（⇒ 子进程**真持有**那枚组合键，`mute=Ctrl+Alt+M live` 那句不假）；
   - 3 次按键后仍无句子；退到托盘那条：`AC#3 shape2 tray: no popup menu is owned by pid … after the tray callback`（`#32768` 弹窗没出现）。
   ⇒ 结论按证据写：**这是本会话输入注入的可靠性问题，不是门的缺陷**（同一枚 `toggleMute` 在 2b 里每次都能开、能关、能收数）。仪器把这一支写成 **SKIP-LOUD**（`raw-live-4.log:425` 逐字：`SKIP-LOUD: no user gesture reached the shipped process in this session: 3 synthesised chord press(es) missed …`），⛔ 不把它当通过、⛔ 不据此判本格失败。**真手按键那一发**仍欠机主在场。
   - ⚠ 我自己的事故（如实入账）：第一发（`raw-live-1.log`）形②被判红，红的是**我写的一条错断言** `frames_dropped > frames_sent`——那个形状恰恰是"本 harness 没有 C8 信道消费者、有界信道满"的已知形（同 `probes/291/mic/00-four-shape-readings.md:35` 的 `levels=1158 frames_sent=6 frames_dropped=1152`），不是设备故障；我把它改成"60..250 枚的节奏带"那条尺。读数本身没受影响。
2. **发运进程上没有逐值打印器**：`levels_delivered` 是**计数**，不是值；逐值只有台件侧（2b）能打。⇒ 对 AC#3"采样数与出处"够用；若要"发运进程上的逐值轨迹"，那是新的可观测面（本腿⛔ 未改产码去加它，归编排者归口，本件不立票）。

## 与 AC#3 两条的对照（只列读数落在哪，⛔ 不判框）

| AC#3 子句（票面原文） | 本件读数 | 落点 |
|---|---|---|
| 默认档 `mic_muted_default` 仍 true、装配走 `WithStartMuted(true)`、verdict 仍打「设备未打开」 | 形① 表 1–4 ＋ `mutedAtBoot==true` ＋ `frames_sent=0` | 已量（发运 exe，空 data dir＝出厂档） |
| 用户主动开门之后 `SetAudioLevel` 真收到数（采样数＋出处） | 2a：**129 枚**（发运进程退场行）／2b：**93 枚、31.00/s、mean 0.003697、max 0.057523**，关回即 0 枚 | 已量（`toggleMute` 是两支手势的生产执行者） |
| "随声音变化" | ⛔ 未量 | 欠机主在场对麦说话那一发；今日扬声器→麦通路无效（291-mic 已量） |

## 门禁读数（本腿现跑）

- `go vet -tags winlive ./cmd/wisp` ⇒ **rc=0**；`go vet ./cmd/wisp`（默认 tag，新件被 tag 挡住）⇒ 随整包编译未单独跑，无产码改动。
- `gofumpt -l cmd/wisp/resident_audio_290_live_windows_test.go` ⇒ **空、rc=0**；`gofumpt -w` 无改动。
- `sh scripts/d22scan.sh` ⇒ **clean、rc=0**（`bans #1-5 cmd/` examined **77** production Go files；`ban #8 cmd/` examined **278** Go files＝含本枚新 test 件；无 emoji/裸 `go func(` 违规）。
- `go test -tags winlive ./cmd/wisp -run 'TestAC290Live' -count=1 -v` ⇒ **rc=0**：`--- PASS: TestAC290LiveShapeOneDefaultBootOpensNoDevice (19.36s)`／`--- SKIP: TestAC290LiveShapeTwoTheMuteHotkeyOpensTheGate`（SKIP-LOUD，理由见 2c）／`--- PASS: TestAC290LiveScalarsAtTheBallSeamAfterTheUserPath (5.41s)`；`ok … 54.189s`。
- `go test ./cmd/wisp -run 'TestAC247ShippedDefaultsAreTheOnesThisLegReads|TestAC290' -count=1` ⇒ **PASS 7 / 0 红，rc=0**。
- 整包 `go test ./cmd/wisp`：本腿**未跑**（派单限定同一时刻只许一发；另有多腿在飞，整包名册归编排者）。本腿只新增一枚 `winlive` tag 件，默认编译面零变化。

## 机主在场那一发（若要做"真手"那一格）

一条命令，键盘上的手按一下即可（10:24 那一发就是这么成的）：

```
cd /d/work/workspace/projects plans/Wisp
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
export WISP_LIVE_MIC=1
go test -tags winlive ./cmd/wisp -run 'TestAC290LiveShapeTwo' -count=1 -v
```

流程：exe 起来约 3-4 秒后屏幕上出现球（托盘同时挂上）；**看到球就按一次 `Ctrl+Alt+M`**，等 4 秒后按 `Ctrl+C`（测试自己会给子进程发 break，无需手按）；读数＝`AC#3 shape2 exit line` 那行与 `levels_delivered`。⚠ 若该 chord 被别的程序占着，boot 行会写 `mute=…` 不是 live，用例会 SKIP-LOUD 并打全 console。
