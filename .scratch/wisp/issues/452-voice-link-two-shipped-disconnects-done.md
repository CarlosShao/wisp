# 票 452 — 语音链路两个**出厂构建**里的断点：barge-in 根本不在 `go build` 里；拔麦不推状态

**状态**：待派（两处都经验收腿现读现跑发现，编排者已复核）

## 断点一：barge-in 不在出厂构建里
`scripts/build.ps1:270` 的 `go build` **不带 `-tags webrtcap`** ⇒ 进的是 `aec_stub`：`NewApm` 恒返回 `ErrUnavailable`、monitor 不武装。
⇒ 票面/证据里那些 **≤400ms barge-in 读数全部来自专用构型**，出厂用户拿到的那份**没有这个能力**。
⇒ 而且 26p-v2 复跑 6 窗只 **1 火** ⇒ 「每轮必火」也不成立。

## 断点二：拔麦克风不推状态，球可能停在 Listening
`EvAudioDeviceLost` 在产码里的全部命中 = **4 处**，其中 **2 处是注释**：
- `internal/statemachine/events.go:30` 常量声明
- `internal/statemachine/table.go:97` D43 转移表条目（`From: StateListening → To: StateError`）
- `cmd/wisp/resident_audio_windows.go:39` **注释**
- `cmd/wisp/resident_audio_windows.go:358` **注释**，逐字「这一条不推状态、也不拒绝启动」

⇒ **零派发**：语音侧从不读 `mic.Err()`，因此**表里有那条边、码里没人走那条边**。用户拔掉麦克风，球可以一直停在 Listening。这正是票 32 AC#5 说的 stuck state。

## AC

- [x] **AC#1 先把两个断点都变成可判事实**：(a) 跑 `go build -n`／读 `build.ps1` 现证**没有** `webrtcap` 标签，并给"打了标签时 `NewApm` 返回什么、没打时返回什么"两条实测；(b) 写一枚测试：把 `mic.Err()` 置成设备丢失，断言状态机**收到** `EvAudioDeviceLost` 并走到 `StateError`。⛔ 先改产码 —— **它必须先红**。
- [x] **AC#2 决定 barge-in 的出厂形态并落成事实**（⓪⓪⓪：按最有价值＋主流 harness 选，写理由）：三选一并写明代价 —— ①**出厂即带 `webrtcap`**（代价：体积/构建依赖）②**出厂不带，但产品明说"此构建无声学回声消除"**并让界面可判 ③**换成不依赖该标签的默认路径**。⛔ 不选"什么都不做"；⛔ 不许一边不带标签一边继续引用 ≤400ms 那些读数当成绩。
- [x] **AC#3 把 `EvAudioDeviceLost` 接上派发**：语音侧检测到设备丢失就派发该事件（走 `observe.Registry.Spawn` 若要新协程），球与面板都离开 Listening；拔麦测试由红转绿。
- [x] **AC#4 两处注释同步**：`resident_audio_windows.go:358` 那句「这一条不推状态」必须改成事实句（接上之后就变成假的）。⇒ 这是**同一条改动里的必做项**，⛔ 留一句过期注释。
- [x] **AC#5 门禁**：`gofumpt` 空／`go vet` rc=0／`d22scan` clean／path-length GREEN／`go test ./internal/audio ./internal/statemachine ./cmd/wisp` 全绿（整包跑不完如实写）。

## 边界
⛔ 动 `thresholds.go`／golden／`allowlist.txt`／`docs/PLAN.md`／`docs/specs/**`／D43 转移表（**表里那条边已经对，别改表**）／`frontend/**`／`design/**`；⛔ 放宽既有断言；⛔ 把 `ErrUnavailable` 改成"假装成功"；⛔ add -A/--amend/reset/rebase/stash/checkout .；⛔push。
写面：`internal/audio/**`、`internal/speech/**`、`internal/statemachine/**`（只读该表）、`cmd/wisp/**`（每波一枚写腿）、`scripts/build.ps1`、`docs/evidence/s1/452-*`。

## 撤销口令
「452 撤」＝删本票与其证据件（不 reset/amend）。

## Progress
- 2026-11-15 写腿（agent-51，`cmd/wisp` 本波唯一写腿）：
  - AC#1(a) 现跑：`scripts/build.ps1:270` 的 `go build -trimpath -ldflags …` 无 `-tags`，全脚本 `webrtcap` 命中 0；
    `go list ./internal/audio/aec` = `GoFiles=[aec.go aec_stub.go]` ⇒ `NewApm` 恒 `nil, ErrUnavailable`（已由 `internal/audio/aec/aec_shape_452_test.go` 钉住）。
    带标签：`go build -tags webrtcap ./internal/audio/aec` **编不过**（`api/scoped_refptr.h: No such file or directory`，本机没有 webrtc 树/meson）；
    `CGO_ENABLED=0 go build -tags webrtcap` rc=0（标签单独不给能力）。读数原文 `.scratch/wisp/probes/452/ac1-a-measure.txt`。
  - AC#1(b) 先红：两枚测试写在 `cmd/wisp/resident_audio_device_lost_452_windows_test.go`，改产码**之前**跑出
    「姿态是空串 `""`」与「拔麦 3s 后仍 Listening、球投影 `[]`」两条红句（逐字抄进证据 §3）。
  - AC#2 裁定 **②**（出厂不带标签 + 产品明说 + 界面判得了）：`aec.BindingBuilt()` 能力位、`bargeinShippedPosture452` 常量在
    **构造期**就写进 `reason`、启动印一次、`posture()` 增 `插话=`。①的代价＝整条原生构建流水线（本机无树、无 meson/ninja、`fetch-deps.ps1` 无此路），
    ③的代价＝重做 389/392/401 三轮读数；面板能力字段出写面 ⇒ 已写明落点移交（证据 §2.3）。
  - AC#3 已接：`residentAudio.captureErr()`（粘性，源头 `Err()` 是破坏性一次读）+ `rv.checkAudioDevice()` +
    `rv.run()` 两个 select 各一条 ticker case（**不新增协程**，骑既有 owner/`asr-infer` 线程）+ 仅从 `StateListening` 派发 D43#14（表未动）。
    接线 `rv.deviceErr = ra.captureErr` 在 `rv.start()` 之前；测试对该生产行做存在性断言。
  - AC#4 已改：`resident_audio_windows.go` 启动分支的「这一条不推状态…」与文件头 "It pushes NO state" 改成「启动期/运行期」两句。
  - 突变：M1 关掉派发 ⇒ 行为测试红；M2 删掉生产接线行 ⇒ 接线断言红。两次都还原。
  - 证据：`docs/evidence/s1/452-voice-link-two-shipped-disconnects.md`。
  - 未做：重新插回/换设备后的恢复（球停 Error 等 `EvErrorAck`，D43#38）；面板能力字段（出写面，已列落点）。
---

## 收件（2026-11-15 14:4x，`269665ace`）：**四格翻，AC#2 只翻一半 —— 界面那半是我派单时划掉的**

**尺1**：`269665ace` 在，11 件：产码 `cmd/wisp/resident_audio_windows.go`／`resident_voice_bargein_windows.go`／`resident_voice_windows.go`／`internal/audio/aec/{aec_stub.go,aec_webrtc_windows.go}`，测试两枚新件，证据件＋两枚读数件。
**尺2**：`docs/evidence/s1/452-voice-link-two-shipped-disconnects.md`；读数 `.scratch/wisp/probes/452/{ac1-a-measure.txt,cmdwisp-full-summary.txt}`。
**尺3 我自己复跑**：`TestAC452LostMicWalksD43FromListeningIntoError` ⇒ **PASS 0.52s**；`TestAC452NewApmAnswersMatchTheBuildShape` ⇒ **PASS 0.00s**（`./internal/audio/aec`）。

### AC#1 翻 —— 两条红句都是真读数
**AC#1(a)**：`this build has no AEC3 binding (aec.NewApm answers ErrUnavailable, scripts/build.ps1 passes no -tags webrtcap), and the barge-in monitor's posture at construction is "": the product states neither that this binary cannot cancel echo nor that Path-C barge-in is therefore unavailable…`
**AC#1(b)**：`the microphone was unplugged … while the machine sat in Listening, and 3s later the state is still Listening. Nothing in the production tree reads the capture leg's terminal error, so D43#14 … has zero dispatchers and the ball stays in Listening with a microphone that no longer exists.`
**四条实测（我逐条读过原文）**：未打标签 `go list` ＝ `GoFiles=[aec.go aec_stub.go]`／打标签 ＝ `GoFiles=[aec.go] CgoFiles=[aec_webrtc_windows.go] Ignored=[aec_stub.go]`／**`go build -tags webrtcap` ⇒ `wap_shim.cpp:23:10: fatal error: api/scoped_refptr.h: No such file or directory`，rc=1**／`CGO_ENABLED=0` 加标签 ⇒ **rc=0（标签单独不产生能力，仍是桩）**。

### AC#3 翻 —— 拔麦测试先红后绿
`residentAudio.captureErr()` **粘性读**（源头 `Err()` 是破坏性一次读，探针读了会让静音手势/退出行答 nil——这个细节是真功夫）→ `rv.checkAudioDevice()` ＋ 250ms ticker 挂在 `rv.run()` **两个既有 select** 上 ⇒ **不新增协程**（`observe` 名册不在写面），派发跑在既有 owner/`asr-infer` 线程。
突变两发都验：关掉派发 ⇒ `TestAC452LostMicWalksD43FromListeningIntoError` **行为红**；删掉生产接线行 ⇒ **接线断言红**（两次都还原）。⇒ 这一格是**行为牙**，不是字面钉（A1880）。

### AC#4 翻 —— 那句注释真改了，而且是"缩小射程"而非"改口"
`resident_audio_windows.go:401` 现在是 `// TICKET 452 rewrote this sentence: it used to say "这一条不推状态" as a …`；`:407` 的运行期句子改成「**启动期**这一条不推状态…」⇒ **句子变真了**（启动期那一段确实不推，运行期 ticker 才推）。⛔ 不是把假话删掉了事。

### AC#2 ⛔ 只翻一半 —— 原因是我，不是它
**已做**（真落地）：`aec.BindingBuilt()` 能力位（两种编译形态各一份）／`bargeinShippedPosture452` 在 monitor **构造期**就写进 `reason`（不是等首次 arm）／启动印一次 `wisp: 本构建无回声消除…`／`posture()` 增 `插话=` 字段／测试钉住 `build.ps1` 那行无 `-tags`。
**没做**：面板能力字段 —— `SettingsView.Voice` 要加 `AecBuildAvailable bool`（数据源 `aec.BindingBuilt()`，与用户可改的 `AecEnabled` **并列分开**），落点已写进证据 §2.3。
⇒ **是我在派单里把 `internal/panel/**` 从写面划掉了**，而 AC#2 原文要求"让界面可判"。⇒ 这一格**判部分满足**，剩下的是：
1. `internal/panel` 那个字段（**我们可派**，下一波一枚写腿）
2. 界面那一半（`frontend/**` ⇒ **机主外部前端 agent**）

### AC#5 翻 —— 整包两发都跑完
| 发 | 结果 |
|---|---|
| 13:47 树 | 红 **2 枚**（`TestAC4FocusReturnToPriorWindowGap33r5` 5.29s／`TestAC246DevLegIgnoresTheTestTaskInjection` 36.07s），`PASS=626 SKIP=17`，1341.00s |
| 14:14 最终树 | 同上**只有那 2 枚**，`PASS=626 SKIP=17`，1190.00s |
⇒ **第三枚红（`Test320Fetch…`）这两发没出现** ⇒ 反证它确实是**互扰型**（A1876 记的形状）。

### 一处它主动收窄的声称（我认）
面板包里今天**没有 D43 状态字段**，所以它只写"球离开 Listening"，**没写"面板也离开"**。⇒ ⛔ 不夸大。

### 这一张不结案
AC#2 的界面半 ⇒ 本票**保持未结**，剩两件（`internal/panel` 字段 ＋ 前端界面），已在上面具名。
---

## AC#2 补齐（2026-11-15 15:2x，`aaa40b15d`）：**5/5，本票结案**

**尺1**：`aaa40b15d` 在，6 件：`cmd/wisp/panel_config_store.go`、`cmd/wisp/panel_voice_aec_452_test.go`、`internal/panel/config_handlers.go`、`internal/panel/config_452_test.go`、`internal/panel/config_313_test.go`、证据件。
**尺2**：证据件在（`docs/evidence/s1/452-voice-link-two-shipped-disconnects.md` 追加）。
**尺3 我自己复跑**：`-test.run '^TestAC452'` ⇒ **`RUN=3 PASS=3 FAIL=0`**（`TestAC452PanelVoiceTellsThePageWhatTheBuildCanDo` 0.04s／`TestAC452ShippedBuildSaysItHasNoEchoCancellation` 0.00s／`TestAC452LostMicWalksD43FromListeningIntoError` 0.53s）。
**赋值点我核**：`cmd/wisp/panel_config_store.go:137 AecBuildAvailable: &aecBuildAvailable,`，真值源 `:117 aecBuildAvailable := aec.BindingBuilt()`；面板字段 `internal/panel/config_handlers.go:236 AecBuildAvailable *bool`。

### 这一格的设计（值得记）
- **指针即闸**：`nil` ＝ 腿没答（**键不在线上**）／`&false` ＝ **答了**"本构建没有"。⛔ 没有和用户开关合成一个布尔 —— 注释 `config_handlers.go:226` 逐字写着「`AecEnabled` = 你想不想开（可写），`AecBuildAvailable` = 这枚二进制能不能」。
- **两发突变都验，且都还原**（sha256 前后同 `EB28195A…`）：
  - **M1 删接线** ⇒ 红句「the settings payload's voice half carries NO aecBuildAvailable key: the page cannot tell that this build cannot cancel echo…」
  - **M2 把能力位并入用户开关** ⇒ 红句「the capability bit is being derived from the user's switch instead of the build, which is the one-boolean shape ticket 452 forbids」
- **它自己点出的 A1880 事实**：面板侧那枚**键集冻结**测试在 M1 下**仍然通过**（它直接构造 view，看不见组装）⇒ 所以行为测试必须放在 `cmd/wisp` 的**线材**测试上。它就是这么放的。⇒ 这是"字面钉 vs 行为牙"的一个活样本。

### 界面那一半的去向（⛔ 不算我们没做）
`frontend/**` 是本编队**零写面**（机主外部 agent）。本腿交付的是**契约**：
> **`voice.aecBuildAvailable`** —— JSON 布尔。`true`＝这枚二进制能消回声；`false`＝本构建没有；**键缺席**＝腿没答。**只读**，⛔ 永不与用户开关 `voice.aecEnabled` 合成一个。

⇒ 前端要做的就是按这个键渲染，契约已写明，**交给机主外部前端 agent**。

### ⚠ 整包两发：都**静默死了**（这是本票最要紧的披露）
14:49→437s（339 起／338 PASS／0 FAIL）与 14:59→386s（337 PASS／1 FAIL）。两发**整进程消失、无 `--- FAIL`**，最后一条 `=== RUN` 都是 **`TestAC4PriorFocusSurvivesARefusedPanelSample`** ⇒ **就是票 450 在册的那个形状**。
★**它主动点名的覆盖事实**：本票新测试按文件序排在**死点之后** ⇒ **两发整包对它覆盖＝0**，绿读数是定向那一发。
⇒ 我把这条写成台账 A1882：**静默死不给后面的用例打 SKIP**，所以 PASS 数看着"正常"，实际**被截断**。⇒ 任何"整包绿"的读数**必须报告终点锚有没有命中**（票 314 用的 `TestZZ315A1`），⛔ 只报 PASS 数。

### 本票 5/5，结案
剩一件**不在本编队写面**的已具名交接（前端渲染 `voice.aecBuildAvailable`）。
