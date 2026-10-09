# 真机窗口逐形操作单（leg `livewin-1`，2026-10-09）

> 这张单子只服务一件事：把"机主在场的那一次真机窗口"压到最短。
> 每一形给满五栏：①落到哪枚用例/哪段台件（`文件:行`＋逐字引那一行）②那一发要敲的完整命令行
> ③它会打印什么（逐字格式串）④这一形今天需不需要机主在场 ⑤缺什么仪器才拿得到。
> 本腿零 go 命令（Go 编译/测试面由写腿 `290-r1` 独占，本腿全程未跑 `go build/vet/test/list/env`），
> 所有"读数形状"一律从源码与票面逐字读得；行号全部 `sed -n`/`grep -n` 现跑，不凭记忆。
> 凡引用行号，先记一句**锚**：本件的行号对 `fdee536b`（2026-10-09 15:23:40 +08）有效，
> 那一笔是 `290-r1` 的开门产码，它把 `cmd/wisp` 三枚文件的行号整体往后推（见 §9 第 2 条）。

---

## 起手锚（原样抄自本腿第一条命令的输出）

```
$ date; git log --oneline -1; git status --short -- internal cmd | head
Fri Oct  9 15:14:14 CST 2026
--- log
7c89ba13 票 289 结案（289-v1 四格判语全成立⇒我翻四格＋Status 相抵归位＋改名 -done）＋A783 落账＋立票 292（第五句过期叙述）；★记我一枚自抓缺陷：今天三枚新票名全超 README 规矩 9 的 100 字符帽、把 sh scripts/check-path-length-budget.sh 打成 VERDICT RED，而那条"开票后先跑它"的定式是我 10-07 自己写进 README 的
--- status internal cmd
（空，0 行）
```

- `git status --short -- internal cmd` 在 **15:14 回 0 行**（期望成立）。
- 同一把尺 **15:2x 复跑回 4 行**，来历具名＝`290-r1` 的在飞改动：
  `M cmd/wisp/resident_audio_windows.go`、`M cmd/wisp/resident_ball_windows.go`、`M cmd/wisp/resident_windows.go`、
  `?? cmd/wisp/resident_mute_290_windows_test.go`。本腿没有还原、没有提交它们中的任何一枚。
- **15:23:40 那四行又消失**，因为它们以 `fdee536b`「票 290 / 290-r1 落地腿产码（甲-1：只把门拧开）」入库，
  落在我第 1 笔 `fcc60648` 之上。⇒ **本件的射程因此改变了一件事**：票 290 的"开门"那一形，
  在写这份单子期间从"今天不存在"变成"码在 HEAD 里、只差重新 build"。§3 按后者写，并标出差异。
- 分支 `dev`；`ls -d third_party/sherpa-onnx build` 现跑**两枚都在**：
  `third_party/sherpa-onnx` 里有 `onnxruntime.dll`／`sherpa-onnx-c-api.dll`／`sherpa-onnx-cxx-api.dll`，
  `build` 里有 `wisp.exe` 等。⇒ 下面每一条 harness 路径都是验过存在的。
- 但是（这条不在任何票面上）：`ls -l build/wisp.exe` 现读＝`2026-10-07 11:57`，
  **既不含 247-r1 的接线（10-09）、也不含 290-r1 的开门（10-09 15:23）**。
  ⇒ 凡是"启动 exe 读 boot 行"的形（A2／C1／C2）**都必须先重新 build**，也就是都要排 Go 面，
  不是"机主在场就能敲"。

### 窗口开打前的四把前置尺（零 Go、零机主）

```sh
cd "D:/work/workspace/projects plans/Wisp"
date; git log --oneline -1                                    # 记下锚，行号以它为准
ls -d third_party/sherpa-onnx build                           # harness 两枚目录
tasklist //FI "IMAGENAME eq balldebug.exe"                     # 必须 "No tasks are running"
tasklist //FI "IMAGENAME eq wisp.exe"                          # 必须 "No tasks are running"
```

后两把不是礼节：`cmd/wisp/resident_audio_247_live_windows_test.go:214` 那枚
`requireNoCompetingAudioProcesses` 会**自己**查，命中就直接红——
逐字 `:222-223`：`"%s is running; the live reading would be shared with it (ticket §排程 requires 0):\n%s"`。
⇒ **顺序铁律：A2/C1/C2 起的 `wisp.exe` 必须在跑 A1 之前退干净**，否则 A1 一定红在那两行上。

---

## 0. 一分钟盘点（三桶＋总数）

**枚数口径先说死**：本单**整族穷举、不是抽样**。盘上展开是 **3 枚出处 / 9 形**
（A＝票 247 的 `AC#2`/`AC#4`/`AC#6` 三格；B＝票 291 `AC#2` 的四形；C＝票 290 的开门/关门两形）。
派单里"六族"这个词在盘上找不到对应物（既不是 3 也不是 9），具名报回在 §9 第 1 条。

| 桶 | 形（编号见下文） | 机主占用 |
|---|---|---|
| **需要机主在场** | A1（`247 AC#2`，人声 6 秒）＝B②；A2（`247 AC#4` 真机默认开机，1 分钟，不需要人声但需要他的桌面与配置真相）；A3（`247 AC#6` 真设备两形，禁用设备／撤权限各一次，约 4 分钟，**必须放最后**）；C2（`290` 开门，按一次 `Ctrl+Alt+M`，1 分钟，**与 A2 同一次开机**） | **合计约 7-8 分钟**，其中真正"开口说话"只有 **6 秒** |
| **不需要机主（可提前拿）** | B①（完全安静：A1 那一发里自动出两窗）；C1（关门：与 A2 同一发同一句，差别只是"谁来读 stdout"）；`247 AC#6` 的**装配面**三形（早绿，属历史凭据，不占窗口） | 0 分钟（但 **A2/C1/C2 都要 Go 面先重新 build**，见上面锚） |
| **今天不可能（缺一枚仪器，不是缺人）** | B③ 的"已知幅度/理论 RMS"那一半（没有任何件打印一枚 wav 的 RMS）；B④（没有 DC＋正弦的**合成**夹具，也没有能吃仓内帧的 mean/max 打印器）；B② 的"纯人声窗"（`sound` 窗硬编码同时外放 `Alarm01.wav`，且 wav 缺失时整发 SKIP、连读数都没有）；C2 的"电平随声音变化"那一维（常驻路径零打印器，只有退出行计数）；A3 的"设备被占"真形（本仓共享模式，无独占打开路径） | 不可用（各形缺什么，逐条写在 ⑤ 栏） |

**总数：9 形。机主在场 4 形（A1／A2／A3／C2，≈7-8 分钟，说话 6 秒）；不需要机主 2 形（B①／C1）；今天拿不到 3 形（B②纯人声那一半／B③理论 RMS／B④整形）。**
另外 A3 只需"这台机器＋机主的系统设置"，A1 是唯一需要他**出声**的形。

---

## 1. A 族（票 247 `.scratch/wisp/issues/247-the-capture-stack-has-zero-importers-so-no-real-microphone-level-ever-reaches-the-ball.md`）

### A1 `AC#2` 球上的数来自声音（票面 `:25` 那一格）

**① 落点**
`cmd/wisp/resident_audio_247_live_windows_test.go:128`，逐字：
`func TestAC247LiveMicrophoneLevelsReachTheBallSeam(t *testing.T) {`
门控 `:129-130`，逐字 `if os.Getenv("WISP_LIVE_MIC") != "1" {` /
`t.Skip("AC#2 needs a real microphone: run with WISP_LIVE_MIC=1 after both tasklist rulers read 0")`。
三窗时序：`:167` `time.Sleep(3 * time.Second)`（quiet-A）→ `:171-173`（`playAlarmThroughSpeakers` 后 6 秒，loud）→ `:176-178`（3 秒，quiet-B）。
**这一发不吃 `mic_muted_default` 的门**：它用 `:139` 那枚临时配置主动拧开——
逐字 `dir := writeAudioConfig(t, func(c *config.Config) { c.Audio.MicMutedDefault = false })`
⇒ **A1 与 290-r1 是否入库无关，今天就能敲**（票 290 现量 `:16-18` 也点名了这枚"不是用户能走到的路"）。

**② 要敲的完整命令行**（逐字抄自前人件 `.scratch/wisp/probes/247/r1/20-ac-readings.md:54-55`，本腿未跑）

```sh
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" WISP_LIVE_MIC=1 \
  go test ./cmd/wisp/ -run TestAC247LiveMicrophoneLevelsReachTheBallSeam -count=1 -v -timeout 180s
```

**③ 它会打印什么**（逐字格式串，都出自那枚文件）
- `:156` `t.Logf("AC#2 出处: capture=%q render=%q pinned_thread=%d", capDev.String(), renderDev.String(), mic.ThreadID())`
  ⇒ **`AC#2 出处:` 这一行就是"开始说话"的信号**，中段 6 秒落在它之后约 3 秒。
- `:181-182` `t.Logf("AC#2 counters at the seam: levels=%d frames_sent=%d frames_dropped=%d last_err=%q", ...)`
  历史读数：`levels=375 frames_sent=6 frames_dropped=369`（消费侧为零，`internal/speech` 只有 `doc.go`）。
- `:184` `t.Logf("AC#2 reading: %s", p.String())`，三窗各一行；`String()` 在 `:96-101`，形状逐字：
  `name: samples=N mean=0.369079 max=0.418274 levels_per_s=31.00`（`mean`/`max` 各 6 位小数，`levels_per_s` 2 位）。
- 判绿/判红的三句原话：
  `:194` `if loud.mean() <= quietA.mean()*2 || loud.max() <= quietA.max() {` → `:195-197`
  `"the loud phase is not distinguishable from silence: %s | %s | %s (speakers off, or routed to an endpoint this microphone cannot hear)"`；
  `:199-200` `"the control phase stayed as loud as the sound phase, so this is noise not a reading: %s | %s"`；
  `:190-192` `"level cadence %.2f/s is not the frame cadence %.2f/s (P8 cost clause: delivery on the pinned thread)"`（带宽 0.6x-1.6x，标称 `1/FrameDuration`＝31.25/s，`internal/audio/audio.go:53` 逐字 `const FrameDuration = FrameSamples * time.Second / TargetRate // 32ms`）。
- 环境缺件时不是 FAIL 而是 SKIP：`:234` `t.Skipf("no alarm sample to play out loud (%v): AC#2 then needs a person speaking into the microphone", err)`（源文件＝`%WINDIR%\Media\Alarm01.wav`，`:232`）。
- 两窗静默底噪的历史锚（引用前先重跑）：`247:102` 记 `quietA 0.368462 / sound 0.381574 / quietB 0.382262`；`20-ac-readings.md:23-25` 记 `0.369079 / 0.382697 / 0.382564`。

**④ 需不需要机主**：**需要**，且只需要他**在中段 6 秒正常音量说话**。跑＋编译约 1 分钟占用。
判据原话（`20-ac-readings.md:59-61`）："The box closes when the `sound` window's mean and max are clearly above both quiet windows"。

**⑤ 缺什么仪器**：这一形本身仪器齐。但**"人声"与"外放"分不开**——`sound` 窗里 `Alarm01.wav` 一定在放（`:171` 硬编码），
所以拿到的是"外放＋人声"的混合读数；想要票 291 `AC#2` ② 那种"纯人声窗"，缺一枚**可关闭外放音源的开关或第二枚台件**（今天没有 env）。
不许为此放宽 `:194` 那枚 2x 判据（票 291 禁区 `:23` 逐字禁的就是这一支）。

### A2 `AC#4` 默认档真机开机：设备到底有没有被打开（票面 `:27`）

**① 落点**
产码打印位＝`cmd/wisp/resident_windows.go:320`，逐字
`fmt.Printf("wisp: 采集腿：%s\n", raudio.posture())`
（`fdee536b` 之前那一行在 `:295`；漂＋25）。`posture()` 在门关着时走 `unrunningClaim()`，返回的就是 boot 装配那句：
`cmd/wisp/resident_audio_windows.go:298`，逐字
`ra.verdict = fmt.Sprintf("采集腿已装配、门处于静音：设备未打开（[audio] mic_muted_default=true，来源 %s）", cfgSource)`
（票 290 现量 `:19-20` 引的 `:238` 是**同一句**，在 `fdee536b` 之前的行号；见 §9 第 2 条）。
装配那一步＝同文件 `:273`（旧 HEAD `:213`）逐字
`gate := audio.NewHalfDuplexGate(mic, audio.PathT, audio.WithStartMuted(c.Audio.MicMutedDefault))`。
对照用例（**不需要机主、不需要 Go 面之外的东西**，且已是 PASS 史凭据）：
`cmd/wisp/resident_audio_247_windows_test.go:139` `func TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice(t *testing.T) {`，
它钉的正是 `Muted()==true`／`Open()==false`／`FramesSent==0`／`Registry.CountByName("audio-capture")==0`（`:149-166`）。
验收腿 `247:104` 说的"缺的那发"＝**在真机的真数据根上把那句话读出来**，不是再钉一遍装配面。

**② 要敲的完整命令行**（需要先重新 build，本腿没跑、也不能跑）

```sh
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" ./build/wisp.exe      # 无参数＝常驻入口
```
入口证据：`cmd/wisp/main.go:59-66`，逐字 `if len(args) == 0 {` … `attachParentConsole()` / `runResident()`
（`attachParentConsole` 意味着 boot 行会落到当前终端的 stdout）。
读完 `Ctrl+Alt+M` 那一发（C2）之后退出进程，再跑 A1。

**③ 它会打印什么**（默认档应有的那一行，逐字拼出来）
`wisp: 采集腿：采集腿已装配、门处于静音：设备未打开（[audio] mic_muted_default=true，来源 <cfgSource>）`
`来源` 那半是**自证的**：它打印的是**哪张表做了这个决定**——
- 读到了数据根的 `config.toml` ⇒ 打印那枚路径（`:238-240` 的 `cfgPath := filepath.Join(dataDir, configFileName)`）；
- 读不到 ⇒ 落内置表并把来源打成字面量，`:250` 逐字 `cfgSource = "compiled defaults (config.NewDefaults)"`，
  并且**另有一行 Warn**，`:251-253` 逐字
  `slog.Warn("audio: config unreadable at boot; the capture leg is built from the compiled default table", "path", cfgPath, "err", configErrText(err, c), "mic_muted_default", c.Audio.MicMutedDefault, "voice_enabled", c.Voice.Enabled)`。
两种落子都不算偷开麦：内置表那半是保守的（`mic_muted_default = true`）。
另外两条形（同一段产码，逐字）：
- `voice.enabled=false` 那一形 `:265` `fmt.Printf("wisp: 麦克风采集腿未构造（[voice] enabled=false，来源 %s）：球不会收到任何电平，本进程其余部分照常\n", cfgSource)`；
- 门关着但进程退出时 `:358` `fmt.Sprintf("audio capture leg stopped: levels_delivered=%d frames_sent=%d frames_dropped=%d reopens=%d", ...)`（门关着 ⇒ `levels_delivered=0` 就是"没采到"的旁证）。

**④ 需不需要机主**：**要他这个人，但不要他的声音**。三件事只有他能当场确认：
这台机器的数据根 `config.toml` 有没有被他改过、球/托盘会出现在他桌面上、麦克风隐私设置是谁的。
纯取那一句 stdout 的动作本身，编排者在机主不在场时也能做（前提是允许在他机器上起常驻进程）。

**⑤ 缺什么仪器**：不缺。唯一硬前置＝**必须有一个含 `fdee536b`/247-r1 的新 build**（现成 `build/wisp.exe` 是 10-07 的）。
仓外那枚"Windows 麦克风隐私指示"**不是本仓的尺，不许当凭据**。

### A3 `AC#6` 降级照跑（票面 `:29`，真设备那两形）

**① 落点**
今天唯一**绕过门、直开真设备**的台件＝`internal/audio/hotplug_test.go:525`，逐字
`func TestLiveWasapiSmoke(t *testing.T) {`，门控 `:526-527` 逐字
`t.Skip("live WASAPI smoke requires WISP_LIVE_MIC=1 and a real microphone (真机冒烟待票 16)")`。
它 `:529` `mic := NewWASAPIMicrophone()` → `:530 Endpoints()` → `:538 mic.Start(ctx, buf)`，
**不经过 `HalfDuplexGate`** ⇒ 这一族读数不吃 290 的开门，也不吃 `mic_muted_default`。
错误文案的产地＝`internal/audio/device.go`，三句 guidance 逐字现读：
`:86` `const privacyGuidance = "microphone access denied: allow desktop apps to use the microphone under Windows Settings > Privacy & security > Microphone"`、
`:89` `const inUseGuidance = "capture device is held in exclusive mode by another application; close that application or pick another device"`、
`:103`（`hrAUDCLNTDeviceInv` 那支）`detail+" ("+code+"): device invalidated (unplugged or disabled)"`；
拼接形状＝`what + " on device " + dev.String() + " (" + code + ")"`，所以**HRESULT 一定在字符串里**。

**② 要敲的完整命令行**（每一形各敲一次；本腿没跑）

```sh
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" WISP_LIVE_MIC=1 \
  go test ./internal/audio/ -run TestLiveWasapiSmoke -count=1 -v -timeout 120s
```
形甲（无设备）：在 Windows 声音设置里**禁用默认采集设备**（或拔掉耳机麦）后敲。
形乙（无权限）：设置 > 隐私和安全性 > 麦克风，关掉"允许桌面应用访问麦克风"后敲。**跑完立刻恢复**，否则 A1/B 全拿不到读数。

**③ 它会打印什么**（候选行全在源码里钉死，本机落哪一条＝现场读数，本腿不猜）
- `:532` `t.Fatalf("enumerate endpoints: %v", err)` —— 枚举这一跳就断的形；
- `:539` `t.Fatalf("live Start: %v (device occupied? privacy settings?)", err)` —— 打开这一跳断的形，`%v` 里带 `(0x........)`＋guidance 原句；
- `:554` `t.Fatalf("only %d frames in 2s (mic muted or silent?)", frames)`；
- `:534` `t.Logf("capture=%q render=%q", cap1.String(), render.String())`（正常时的设备名行）；
- `:567` `t.Logf("live smoke OK: %d frames in %v, stats %+v", frames, ..., mic.Stats())`。
票面要的"逐字回串"（`0x8889000A`／`DEVICE_INVALIDATED`／`NOT_INIT` 哪一枚）＝上面 `:532` 或 `:539` 里那段 `%v` 的原文。

**④ 需不需要机主**：**需要**（动的是他的系统设置与物理设备），约 4 分钟，含恢复。**这一族必须排在最后跑。**
`247:106` 与 `20-ac-readings.md:130-132` 都具名说过：今天这些是**装配面拼出来的 HRESULT**，不是真设备的回串。

**⑤ 缺什么仪器**：三形里**"设备被占"今天不可能拿到真读数**——本仓采集是**共享模式**
（票 247 现量 1 逐字指 `wasapi_windows.go:210 Open`），而 `inUseGuidance` 只在 `hrAUDCLNTDeviceInUse` 那一支触发，
仓里没有任何独占模式的打开路径；`hotplug_test.go:378 TestOpenOccupiedAndPermissionDenied` 吃的是**构造出来的**错误。
⇒ 票面 `AC#6` 写的是"三形里任取两形"（`:29`），**取"无设备＋无权限"两形就满足判据**，
"被占"这一形要么另立仪器、要么具名带缺，不许写成"大概能行"。

---

## 2. B 族（票 291 `AC#2` 四形；件 `.scratch/wisp/issues/291-level-ruler-does-not-strip-dc.md:17`）

票面原话：四形各 ≥3 秒并**逐形给 `mean/max`**；判据原话（同一行末）：
"**「说话」与「安静」两形之间存在一个既有的、不动用的判据能分开**——分开不了就具名说'分开不了'，不许挑一个能过的门槛凑绿。"

**全族共用的一件仪器**：仓里**唯一**打印 `mean/max` 的地方就是 A1 那枚台件——
`cmd/wisp/resident_audio_247_live_windows_test.go:184` `t.Logf("AC#2 reading: %s", p.String())`
＋ `:96-101` 的 `String()`（形状 `samples=/mean=/max=/levels_per_s=`）。
窗口长度是**写死的 3s / 6s / 3s**（`:167`、`:172`、`:177`），不能只跑其中一窗。
⇒ **B①②③ 三形的读数全部来自 A1 的那一发**（一发拿四窗），B④ 今天无处可打。

### B① 完全安静
① `…_247_live_windows_test.go:141` `quietA := &livePhase{name: "quiet-A"}` / `:143` `quietB := &livePhase{name: "quiet-B"}`（两枚独立 3 秒窗，各约 93 样本 @31/s）。
② 同 A1 那一发（不需要单独的命令）。③ 同 `:184`，逐字行形如 `AC#2 reading: quiet-A: samples=93 mean=0.369079 max=0.418274 levels_per_s=31.00`。
④ **不需要机主**——要的是"房间里没人说话"，任何人在场都能拿到。
⑤ 不缺。票 291 要 ≥3 秒：正好 3 秒，两窗互相给漂移基线（历史上两静默窗差 `0.10178`，`247:102`）。

### B② 机主正常音量说话
① `…_live_windows_test.go:142` `loud := &livePhase{name: "sound"}`，窗口 `:170-174`。
② 同 A1 那一发；操作＝在 `AC#2 出处:` 打印后**说话 6 秒**。③ 同 `:184`，行形如 `AC#2 reading: sound: samples=189 mean=… max=… levels_per_s=…`。
④ **需要机主**（唯一必须出声的一形，6 秒）。
⑤ **缺仪器（这是本族最硬的一条）**：这一窗同时在外放 `Alarm01.wav`（`:171`，源文件 `:232`），
所以拿不到"纯人声"读数；而 `Alarm01.wav` 缺失时走 `:234` 的 `t.Skipf`＝整发 SKIP、**连安静窗都没有**。
⇒ 要"纯人声窗"必须新增一枚开关或第二枚台件（属写腿活）。不许反过来动 `:194` 的 2x 判据凑绿。

### B③ 外放已知幅度音调（先算出该音调的理论 RMS）
① 外放这一半有仪器：`:230 playAlarmThroughSpeakers` → `:236-239` 的 `powershell.exe` + `System.Media.SoundPlayer` 循环播 `%WINDIR%\Media\Alarm01.wav`。
② 同 A1 那一发（外放自动发生，无人工）。
③ 打印的还是 `:184` 那一行；**音调本身不是已知幅度**，且没有任何一行打印它的理论 RMS。
④ **不需要机主**。
⑤ **缺仪器**：仓里没有一枚件/工具把一枚 wav 的 RMS 算出来并打印——`LevelOfSamples` 只是库函数
（`internal/audio/level.go:96 func LevelOfSamples(samples []int16) float64`，`:106 return rms / LevelFullScale`），
命令行面零调用者。⇒ "先算出理论 RMS"这一半今天**拿不到**，只能另立仪器。
可手算的锚（不需跑）：`LevelFullScale = 32768.0`（`level.go:54`）、满幅整周期正弦读 `1/sqrt(2) = 0.7071067812`（`level.go:35` 逐字）。
历史读数：外放增量 `0.09917` **小于**两静默窗漂移 `0.10178` ⇒ 票 291 裁定 `:37` 逐字结论"按今天的尺，'外放已知音调'与'两次安静'在数值上**分不开**"。

### B④ 人为注入纯 DC＋已知正弦的合成帧（票面逐字：仓内夹具，不许据此判生产）
① 两枚**分开的**夹具：
DC＝`internal/audio/capturelevel_windows_test.go:19` `func TestAC247RealCaptureLoopEmitsLevels(t *testing.T) {`，其 `:23` 逐字
`samples[i] = 12345 // a steady DC level: RMS is exactly this magnitude`；
正弦＝`internal/audio/level_test.go:132` `func TestLevelSineAgainstRootTwo(t *testing.T) {`（`:134 wholeCycleSine(FrameSamples, 32767, 8)`、半幅 `:142`、四分之一幅 `:166`）。
② 若要复跑这两枚（**要 Go 面**，本腿不跑）：
`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/audio/ -run 'TestAC247RealCaptureLoopEmitsLevels|TestLevelSineAgainstRootTwo' -count=1 -v`
③ **它们什么都不打印**：两枚都没有 `t.Logf`（尺＝`grep -rn "t.Logf" internal/audio/*_test.go` 只命中 `hotplug_test.go:534/:567`），
只有失败文案，例如 `capturelevel_windows_test.go:52` 逐字 `t.Fatalf("frame %d: level = %v, want %v", i, l, want)`。
⇒ 拿不到"逐形 mean/max"，除非故意让它红。
④ **不需要机主**（票面自己标的就是"仓内夹具"）。
⑤ **缺两枚仪器**：(a) **没有 DC＋正弦的合成帧夹具**（尺＝`grep -rn "12345" internal/audio/*_test.go` 现读**只有 1 命中**，就在 `:23`；正弦与 DC 从不合流）；
(b) **没有一枚能吃仓内帧的 mean/max 打印器**——唯一的打印器 `livePhase.String()` 长在 `cmd/wisp` 的真机台件里、只收真机投递的 `float32`。
可手算的期望值（当基线用，不当读数）：纯 DC `12345/32768 = 0.376740`，与票 291 现量 `:8` 引用的验收腿现跑值 `0.376740` 逐字一致。
另记尺的定义那句（本票不许顺手动它）：`internal/audio/level.go:37-39` 逐字
`// DC counts: RMS does not remove a DC offset, so a frame held at -32768` / `// reads 1.0. Consumers that need "voice only" gate it themselves (the ball` / `// does, with SilenceLevelGate).`；球侧那道门 `internal/ball/liquid.go:30` 逐字
`SilenceLevelGate   = 0.06 // envelope below this counts as "nobody speaking"`，用在 `:215` `if raw < SilenceLevelGate {`（两枚行号本腿现读吻合，未漂）。

---

## 3. C 族（票 290 `.scratch/wisp/issues/290-nothing-in-production-ever-turns-the-gate-on-so-the-default-config-cannot-deliver-level-to-the-ball.md:36` 的 `AC#3` 两形）

### C1 关门（默认档启动，球收不到电平）
①②③ **与 A2 同一发、同一行**：打印位 `cmd/wisp/resident_windows.go:320`，句子 `cmd/wisp/resident_audio_windows.go:298`；
关门那一支的另一句是同文件 `:320` 逐字 `ra.verdict = "采集腿在跑但设备未交接（gate 未 open）：球不会收到电平"`（票面引的 `:260`，漂＋60）。
④ **不需要机主**（只要有人读 stdout），但需要新 build（Go 面）。
⑤ 不缺。

### C2 开门（用户主动发起之后再按一次拧回）

**① 落点**（**这些行号只在 `fdee536b` 之后存在**，本腿先前 15:1x 读时它们还是"Not Committed Yet"）
- 执行者：`cmd/wisp/resident_audio_windows.go:156` 逐字 `func (ra *residentAudio) toggleMute() (string, bool) {`，
  那一跳本体 `:162` 逐字 `gate.SetMuted(!gate.Muted())`；回读 `:168` `ra.started = gate.Open()`。
- 手势接线：`cmd/wisp/resident_ball_windows.go:321` 逐字 `OnMuteHotkey:    func() { rb.muteGesture("mute-hotkey") },`、
  `:325` 逐字 `OnTrayMute:      func() { rb.muteGesture("tray-mute") },`
  （⇒ 票 290 现量 `:13-15` 引的 `:281 recordBallGesture("mute-hotkey")` 已被这笔推翻，**句面结论"没有任何一条生产路径能把它拧开"现在过期**；尺见 §9 第 3 条）。
- 按键：**`Ctrl+Alt+M`**（默认绑定，`internal/ball/hotkey_windows.go:75` 逐字
  `return HotkeyConfig{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Esc", Panel: "Ctrl+Alt+P"}`；
  若机主的 `config.toml` 改了 `[hotkey] mute`，以配置为准）。托盘那支＝同名"静音"菜单项。
- 配套用例（已在 HEAD，**要 Go 面**才算得数）：`cmd/wisp/resident_mute_290_windows_test.go` 六枚，
  现读名册 `:101 TestAC290BothMuteGesturesTurnTheGateAndBack`、`:169 …OutcomeIsReadOffTheGateNotOffTheRequest`、
  `:196 …NoGateSaysWhichShapeThisProcessIsIn`、`:225 …GestureBeforeTheAttachSaysSo`、
  `:243 …MutedDefaultStillDecidesTheBoot`、`:272 …UnhostedGestureWordingStillSaysTrueThing`。

**② 要敲的完整命令行**（一次开机、两下按键）

```sh
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" ./build/wisp.exe   # 默认档：门关着，读到 C1 那句
# 然后按 Ctrl+Alt+M（开门），再按一次（拧回关门）；最后正常退出，读退出行
```

**③ 它会打印什么**（逐字格式串＋三句成品文案）
打印位＝`cmd/wisp/resident_ball_windows.go:536` 逐字 `fmt.Printf("wisp: ball %s: %s\n", name, outcome)`（未执行那一支在 `:532`，形状同句不同内容）。
- 开门成功 → `%s` 取 `cmd/wisp/resident_audio_windows.go:174` 逐字
  `"已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）"`
  ⇒ 成品行：`wisp: ball mute-hotkey: 已取消静音：设备已交接，采集线程在跑；…不声称）`
- 再按一次关门 → `:172` 逐字 `"已静音：采集已关闭，设备未打开（再按一次取消静音）"`
- 请求了但没交接 → `:191` 逐字 `fmt.Sprintf("已取消静音但设备未交接（gate 未 open，错误分类 %s）：%s；球不会收到电平", class, detail)`
- 按键来得比装配早 → `:524` 打印 `wisp: ball mute-hotkey: no capture leg had been assembled when this key arrived, so there was no gate to turn; the assembly root attaches the mute key to the gate at the end of boot (ticket 290)`
- 这进程根本没有门（`voice.enabled=false` 或腿没装配）→ `:532` 打印 `unrunningClaim()`，即 `:265` 那句"采集腿未构造…"
- 退出时 → `resident_audio_windows.go:358` 那行 `audio capture leg stopped: levels_delivered=%d frames_sent=%d frames_dropped=%d reopens=%d`
  （关门史值 `levels_delivered=0`；开过门再退出会是非零——247 的历史读数是 `376`）。
- 视觉那一维**不在这发的凭据里**：`prototypeVisuals` 默认关、`SetAudioLevel` 在它关着时早退
  （票 247 `AC#10`，`internal/ball/statevisual.go:102`／`internal/ball/liquid_windows.go:42,52-54`），球"呼吸"归票 68。

**④ 需不需要机主**：**需要**（按键要落在他桌面的那个进程上；并且开门＝他主动放行麦克风，属隐私动作）。约 1 分钟，且**与 A2 同一次开机**，不另开机。

**⑤ 缺什么仪器**：`AC#3` ② 那一维——"用户主动开门之后 `SetAudioLevel` 真的收到**随声音变化**的数"——**今天缺仪器**。
常驻路径上电平只进不出：`resident_audio_windows.go:135-140` `func (rb *residentBall) setAudioLevel(level float32)` 直接
`rb.b.SetAudioLevel(level)`，全程**没有任何一处打印逐帧读数**，唯一的计数在退出行（`:358`）。
⇒ 三条可选路（都归写腿/编排者裁，本腿不荐其一）：(a) 给常驻腿加一枚限流的逐窗读数（形状同 A1 的 `String()`）；
(b) 开门那一形改在 `TestAC247LiveMicrophoneLevelsReachTheBallSeam` 里读（但它用的是临时配置 `:139`，票 290 `:18` 已具名说那**不是**"用户能走到的那条路"）；
(c) 承认这一维今天只能取"计数非零"这个弱读数，把"随声音变化"押给 A1。
另外两件的前置不是仪器而是队列：**C2 需要新 build**（现成 `build/wisp.exe` 是 10-07 的，不含开门码），
且**A1 必须排在 C2 之前或之后、不能并存**（`wisp.exe` 在跑 ⇒ A1 直接红在 `:222-223`）。

---

## 4. 建议次序（把机主压到一次、约 8 分钟）

1. **机主不在场**：重新 build 一枚含 `fdee536b` 的 `wisp.exe`＋跑门禁（Go 面；本腿一枚 go 命令都没跑）。
2. **机主在场 1 分钟**：A1（＝B①②③ 一发拿全）——`AC#2 出处:` 出现后说话 6 秒。**先跑它**，因为后面会把麦关掉。
3. **机主在场 2 分钟**：启动新 exe → 读 C1/A2 那句"设备未打开" → 按 `Ctrl+Alt+M` 读 C2 的"已取消静音" → 再按一次读"已静音" → 退出读 `levels_delivered`。
4. **机主在场 4 分钟（放最后，破坏性）**：A3 两形（禁用设备 / 撤桌面应用麦权限，各一发 `TestLiveWasapiSmoke`），**每形跑完立刻恢复**。
5. **不需要机主**：B④ 与 B③ 的理论 RMS——**先补仪器**（写腿），补不出来就在票 291 的 `AC#2` 上具名写"分开不了/拿不到"。

---

## 9. 与票面/派单不符的逐处（本腿现跑所见，⛔ 未改任何票面一字）

1. **"六族"对不上盘**：派单说"六族"，但同一张表只列了 A/B/C 三枚出处；按票面逐格展开是 **9 形**（3＋4＋2）。
   本件按 9 形**整族**写，未抽样。
2. **行号漂（漂因是一笔在飞产码，不是票面写错）**：票 290 现量 `:19-20` 引 `resident_audio_windows.go:238`/`:260`、
   `:8` 引 `:213`，这三枚在 `7c89ba13` 上**逐字正确**（本腿用 `git show HEAD:` 复核过）；
   `fdee536b`（290-r1 开门产码，15:23:40）把它们整体推到 **`:298`／`:320`／`:273`**（＋60，`resident_windows.go` 的 boot 打印 `:295` → `:320`，＋25）。
   ⇒ 派单"290-r1 独占 Go 面、正在改 `cmd/wisp/**`"这一条**在本腿执行期间已完成入库**，本件正文一律用新行号。
3. **票 290 的两句结论已被 `fdee536b` 作废**（照实说，不是缺陷）：
   - 标题与 `:3` 那句"没有任何一条生产路径能把它拧开／`SetMuted` 产码调用者＝0"——现读合尺
     `git grep -n "SetMuted(\|SetSpeaking(" HEAD -- internal cmd | grep -v _test` 回 **4 行**，
     其中 `HEAD:cmd/wisp/resident_audio_windows.go:162: gate.SetMuted(!gate.Muted())` 是**产码调用者**；
     `SetSpeaking` 仍零产码调用者（只剩 `gate.go:97` 注释、`:130`/`:168` 两枚定义行）。
     合尺 3 → 4、单尺 `SetMuted` 2 → 3，与 290-r1 起手件自己具名纠正的口径一致（"派单『由 2 枚变 ≥3 枚』在合尺下少一枚，以原文为准"）。
   - `:13-15` 那句"球上那枚静音热键今天只记一笔"——现读 `resident_ball_windows.go:321`/`:325` 两枚手势**都已有执行者**（`muteGesture`）。
   - `:63` 承诺同批改的 `ballGestureWhy` 那句"no microphone"**已改**：现读 `resident_ball_windows.go:464-466`
     开头逐字 `"no executor was handed to this ball host for this gesture; the gestures with one are the "`，后面已列出两枚静音手势。
4. **票面引用的其它行号，本腿全部现读、逐字吻合、未漂**：
   `internal/audio/device.go:86`/`:89`/`:103`；`internal/ball/liquid.go:30`/`:215`；`internal/audio/level.go:37-39`/`:54`/`:96`/`:106`；
   `cmd/wisp/resident_audio_247_live_windows_test.go:139`；`internal/audio/capturelevel_windows_test.go:21-24`（`12345` 在 `:23`）；
   `internal/ball/hotkey_windows.go:75`（`Ctrl+Alt+M`）；`internal/audio/audio.go:53`（32ms）。
5. **注释/用例名的存在性核验**（按派单"注释里的测试名一律当待验断言"）：
   全部**存在**、无孤儿名——`TestAC247LiveMicrophoneLevelsReachTheBallSeam`、`TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice`、
   `TestAC247DeviceFailureShapesStillBootAndSayTheLoss`、`TestLiveWasapiSmoke`（`hotplug_test.go:525`）、
   `TestOpenOccupiedAndPermissionDenied`（`:378`）、`TestLevelSineToleranceIsNamedAndBounded`（`level_test.go:198`，票 291 `:12` 引它）、
   `TestAC247RealCaptureLoopEmitsLevels`、`TestLevelSineAgainstRootTwo`、六枚 `TestAC290*`。
   ⇒ 派单"不许发明不存在的用例名"这一格：本件**零发明用例名**，唯一零命中的名字是票面自己的
   `recordBallGesture("mute-hotkey")`（已被 `muteGesture` 取代）。
6. **`build/wisp.exe` 是 10-07 11:57 的陈件**（本腿现跑 `ls -l`），不含 247 接线也不含 290 开门。
   任何票面都没写这一条；它决定了 A2/C1/C2 三形的真实前置＝**重新 build（Go 面）**，不是机主。
7. **`third_party/sherpa-onnx` 与 `build` 两枚目录今天都在**（`ls -d` 现跑），所以本件所有 harness 命令行都点得开；
   但**本腿一条都没执行**（`go` 全禁、`wisp.exe` 会占机主桌面，也不属本腿写点）。
8. **票 291 的件名**：派单提醒"15:0x 刚改过短名，别按旧长名找"——现读确实是
   `.scratch/wisp/issues/291-level-ruler-does-not-strip-dc.md`，票面 `:45` 逐字记了改名与 111 字符超帽的事；派单这一条对得上。

---

## 10. 本腿命令与 rc 台账

| 命令 | rc | 备注 |
|---|---|---|
| `date` / `git log --oneline -1` / `git status --short -- internal cmd \| head` | 0 | 起手锚，原样在件首 |
| `ls -d third_party/sherpa-onnx build` | 0 | 两枚都在，harness 可用 |
| `git show HEAD:<path>`（对 5 枚文件）／`git grep`／`git blame -L`／`git ls-files`／`sed -n`／`grep -n` | 0 | 全部只读；`git blame` 那一次是"在飞改动"的唯一证人 |
| `sh scripts/check-path-length-budget.sh`（交件后现跑，15:3x） | **rc=1 VERDICT RED，⛔ 红在本腿以外** | 分母 8559 枚 tracked path、超帽 58、名册覆盖 57、**名册外 1 枚＝`.scratch/wisp/issues/290-nothing-in-production-ever-turns-the-gate-on-so-the-default-config-cannot-deliver-level-to-the-ball.md`（名 106 字符 > 规矩 9 的 100 帽）**⇒ 这一枚是票 290 建票时留下的债（同一族缺陷已由编排者在 `7c89ba13` 的 commit message 里逐字自抓、票 291 的 `:45` 也记了改名），**本腿零改票面、零改名**，只具名报回。**本腿那枚件不是超标项**：相对路径 60 字符、末段名 34 字符（帽 100）。 |
| `go build` / `go vet` / `go test` / `go list` / `go env` | **未跑（不适用）** | **本腿零 Go 命令**：受 `290-r1` 独占 Go 面约束。`sh scripts/d22scan.sh`／`tools/d22scan/d22scan.exe` 也未跑——本腿零产码、零 Go 文件改动，跑它们要占 Go 面，不属本腿射程 |
| 写点 | — | 唯一写点＝本件 `.md`；零 `frontend/**` 零 `design/**`；未碰 `.gitignore`／`design/**`／`probes/{152,161,242,268}`／任何票面；零 AC 框翻动；零 push；两笔 commit 均带显式 pathspec，未用 `add -A`/`amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean` |
