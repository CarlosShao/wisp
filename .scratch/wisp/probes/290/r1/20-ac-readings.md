# 票 290 / 290-r1 — 20 AC 读数（AC#2 枚数尺两列／AC#3 第①形现跑、第②形具名欠）

## AC#2 — 尺＝非测试调用者枚数（票面逐字那条）

尺逐字：`git grep -n "SetMuted(\|SetSpeaking(" <目标> -- internal cmd | grep -v _test`
（改前跑 `HEAD`＝起手锚 `6c976aa6`；改后跑 `HEAD`＝本腿产码笔 `fdee536b`。两次都是对象层尺，与工作树脏否无关。）

| 读数 | 改前（HEAD=`6c976aa6`） | 改后（HEAD=`fdee536b`） |
|---|---|---|
| **合尺**（`SetMuted(` ＋ `SetSpeaking(`，票面 AC#2 那条） | **3 枚** | **4 枚** |
| 单尺（只 `SetMuted(`，票面现量 `:12` 那条的形状） | 2 枚 | 3 枚 |
| 其中"调用者"（非定义行、非注释行） | **0 枚** | **1 枚** |

改前逐字（三枚＝1 句注释＋2 枚定义行，无一枚是调用者）：
```
HEAD:internal/audio/gate.go:97:// source is NOT started; SetMuted(false)/SetSpeaking(false) will start it
HEAD:internal/audio/gate.go:130:func (g *HalfDuplexGate) SetSpeaking(speaking bool) {
HEAD:internal/audio/gate.go:168:func (g *HalfDuplexGate) SetMuted(muted bool) {
```
改后逐字（新增那一枚在 `cmd/wisp` 的生产链路上）：
```
HEAD:cmd/wisp/resident_audio_windows.go:162:	gate.SetMuted(!gate.Muted())          <-- 新增命中
HEAD:internal/audio/gate.go:97:// source is NOT started; SetMuted(false)/SetSpeaking(false) will start it
HEAD:internal/audio/gate.go:130:func (g *HalfDuplexGate) SetSpeaking(speaking bool) {
HEAD:internal/audio/gate.go:168:func (g *HalfDuplexGate) SetMuted(muted bool) {
```

新命中那一枚"在不生产链路上"的两把旁证（⛔ 不是又一枚调试 cmd、⛔ 不是 `cmd/balldebug`）：
1. 调用者所在文件＝`cmd/wisp`（常驻进程本体），函数＝`(*residentAudio).toggleMute`，
   由 `cmd/wisp/resident_windows.go` 的 `rb.attachMuteGate(raudio.toggleMute)` 在装配根注入，
   注入点是 `runResident` 的 boot 序列本身（不是 `cmd/*debug`、不是 `-run` 旗标后面的一段手敲路径）。
2. 球侧两支：`cmd/wisp/resident_ball_windows.go` 的 `OnMuteHotkey` / `OnTrayMute` 现逐字为
   `func() { rb.muteGesture("mute-hotkey") }` / `func() { rb.muteGesture("tray-mute") }`，
   `muteGesture` 唯一做的事是调那枚注入的执行者；`internal/ball` 侧的命中位仍是
   `case hkMute:`（原 `:660-661`）与 `case menuMute:`（原 `:678-679`），本腿⛔ 未动 `internal/ball` 一字。

⚠ **具名分歧（以原文为准，未改票面一字）**：派单与票面 AC#2 写"由 **2** 枚变 ≥3 枚"，
而 AC#2 那条尺（合两枚方法名）在起手锚上现跑＝**3 枚**（注释 `:97` ＋定义 `:130` ＋定义 `:168`）。
"2 枚"只在只量 `SetMuted(` 的单尺上成立。⇒ 本件的判据写作**合尺 3→4、单尺 2→3**，两把都给，
`≥3 枚`那条在合尺读法下今天就已成立（因为分母被少数了一枚）——真正有牙齿的是"调用者 0→1"那一行。

**同一棵树除我以外无人动 Go 面**的尺（防归因漂移）：
`git diff --name-only 6c976aa6..HEAD -- internal cmd | grep -E "\.go$"` ⇒ 4 枚，全是本腿那 4 枚
（中间落的两笔 `7c89ba13` 票 289 结案、`fcc60648` leg `livewin-1` 逐形操作单 ⇒ 只 `.scratch/**`，零 `.go`）。

## AC#3 — 默认档一字不改，两形真机读数

### 第①形（默认双击启动：门在静音位、设备未打开、verdict 仍打那句）＝本腿现跑

真机前置量（进程必须为 0）：
```
tasklist //FI "IMAGENAME eq balldebug.exe" | grep -c "balldebug.exe"  -> 0
tasklist //FI "IMAGENAME eq wisp.exe"      | grep -c "wisp.exe"       -> 0
```
仪器（**票 247 已有的那把出厂表尺**，本腿复跑，读的是真 `WASAPIMicrophone`＋真 `HalfDuplexGate`
＋真 `assembleCapture`，配置＝`config.NewDefaults()` 落成的 `config.toml`，等价于双击启动第一次读到的那一套）：
```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp/ \
  -run 'TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice|TestAC247ShippedDefaultsAreTheOnesThisLegReads|TestAC247HandingTheLevelToTheSeamIsNotVisibility|TestAC247UnmuteReachesTheBallSeam' -count=1 -v
rc=0
--- PASS: TestAC247ShippedDefaultsAreTheOnesThisLegReads (0.00s)
--- PASS: TestAC247DefaultConfigArmsTheGateMutedAndOpensNoDevice (0.02s)
--- PASS: TestAC247UnmuteReachesTheBallSeam (0.01s)
--- PASS: TestAC247HandingTheLevelToTheSeamIsNotVisibility (0.01s)
```
这一发把第①形的三条判据逐条钉住（引行都是该测试自己的断言，⛔ 不是我转述）：
- `mic_muted_default` 仍为 true（`TestAC247ShippedDefaultsAreTheOnesThisLegReads`）；
- 装配走 `WithStartMuted(...)`、`ra.gate.Muted()==true`、`ra.gate.Open()==false`；
- **设备根本没打开**：`ra.mic.Stats()` 的 `FramesSent/BytesSent==0` 且
  `rt.Registry.CountByName("audio-capture")==0`（采集线程一枚都没起，⛔ 不是我推测）；
- verdict 仍打 `resident_audio_windows.go` 那句"采集腿已装配、门处于静音：**设备未打开**（[audio] mic_muted_default=true，来源 %s）"，
  且这句含 `mic_muted_default` 那一枚键名。

本腿自己新增的两枚同形尺（同一次 `-run TestAC290` 里，rc=0）：
- `TestAC290MutedDefaultStillDecidesTheBoot`：三枚默认值现读一字未动＋boot 仍打"设备未打开"
  ＋**注入的源码在 boot 时 start 次数＝0**（离开静音只能由那一发手势做到）；
- `TestAC290BothMuteGesturesTurnTheGateAndBack` 的开头段：静音位起步、`audio-capture` 计数 0（⇒ 本发零新协程）。

⚠ **第①形的残余（具名，⛔ 不自翻 AC#3 框）**：本腿证到的是"默认档装配完 verdict 说设备未打开、
`gate.Muted()` 为真、线程未起"这一**装配级**读数，取的是真源构造（`newRealCaptureSource`）；
"**双击 wisp.exe 那一次**屏幕上打的正是这一句"那一发 console 凭据本腿没跑（见下面"为什么没跑"）。

### 第②形（用户主动开门后 `SetAudioLevel` 真的收到随声音变化的数）＝**具名欠机主在场那一发**

本腿⛔ 未做、⛔ 未冒充，理由逐条：
1. 那一形的判据要**机主在场对着麦克风说话**（票面 AC#3 第②条＋派单第 4 段）；
2. 本腿没有用任何被禁的形状去凑：⛔ 没用外放/合成源/wav 注入器，⛔ 没用
   `cmd/wisp/resident_audio_247_live_windows_test.go:139` 那枚 `MicMutedDefault = false` 临时配置，
   ⛔ 没在测试里 `t.Skip` 之后当作通过（那枚 live 测试在两次整包里都是 `--- SKIP`，
   `TestAC247LiveMicrophoneLevelsReachTheBallSeam`，本腿没有把它读成读数）；
3. ⛔ 本腿没有替机主按下那一枚全局热键：那会在机主不在场时真开这台机器的麦克风，
   并且 `wisp.exe` 常驻会注册全局热键＋Job Object，杀掉它就不是 D38(e) 那十条步的有序退场。
   ⇒ 起手 `tasklist` 两枚读数都是 0，本腿自始至终没有拉起常驻进程（零 `wisp.exe` 生成动作）。

**备好的仪器与命令（给编排者/机主那一发；与 leg `livewin-1` 的操作单
`.scratch/wisp/probes/orch/2026-10-09-live-window-runsheet.md` §3「C 族：票 290 的开门 / 关门两形」并批发，
该单在我产码之前落件，其 §3 仍是"待填"，需要按下面这一版补上"按键之后才有电平"那一跳）**：

```
# 0) 前置（两枚必须为 0）
tasklist //FI "IMAGENAME eq balldebug.exe"
tasklist //FI "IMAGENAME eq wisp.exe"

# 1) 出厂默认档起进程（⛔ 不要改 config.toml：mic_muted_default 保持 true）
GOFLAGS= go build -o build/wisp.exe ./cmd/wisp
build/wisp.exe            # 无参数＝常驻腿：球 + 托盘 + 四枚热键 + 采集腿

# 2) 看第①形：boot 打印的那一行（逐字格式串，改前取自 resident_audio_windows.go:238、本腿落地后住 :298）
#    wisp: 采集腿：采集腿已装配、门处于静音：设备未打开（[audio] mic_muted_default=true，来源 <路径>）
#    同一次 boot 里 [hotkey] 那行会给出 mute= 绑定的是哪一枚键（默认表见 internal/ball 的 DefaultHotkeys）

# 3) 机主按下 mute 那一枚全局热键（或右键托盘 -> 静音那一项，两支行相同）
#    期望逐字（前缀 wisp: ball mute-hotkey: / wisp: ball tray-mute:）：
#      已取消静音：设备已交接，采集线程在跑；电平按票 247 的链路交给球（球屏上会不会呼吸是票 68 那一格，本票不声称）
#    若设备被占用/被拒：改为
#      已取消静音但设备未交接（gate 未 open，错误分类 audio_device）：<指引文本>；球不会收到电平

# 4) 机主对着麦克风说话几秒，再按一次同一枚键（回到静音位）：
#      已静音：采集已关闭，设备未打开（再按一次取消静音）

# 5) Ctrl+C 有序退场，读最后一行（residentAudio.stop 的格式串）：
#      wisp: audio capture leg stopped: levels_delivered=N frames_sent=… frames_dropped=… reopens=…
```
两形的采样出处＝同一枚进程的第 2 行与第 5 行（`levels_delivered`），
"随声音变化"那一半由机主在第 4 步之前/之后的两次读数对照（静音位那一次应当是 0 增量）。

⚠ 本腿把这一格**原样交回**：`AC#3` 的框一枚未勾（框等非实现者裁），
欠的具体那一发＝**第②形的真机读数**（要机主在场），另附带欠第①形那一次真实 `wisp.exe` console 凭据；
票 247 的 `AC#2`/`AC#4`/`AC#6` 三格与票 291 的 `AC#2` 四形按票面排程与这一发**同一次真机窗口**批发。

## AC#4／AC#5 的读数在 `30-gates.md`（逐名 rc＋改前后名册两发并排＋作差）。
