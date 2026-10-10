# 票 295 · 补位腿 `295-a1b` · 第 4 笔：ⓓ 落地腿落点清单 ＋ ⓔ 反形正控落点预判

HEAD `29081a13…`（branch `dev`）。本腿⛔ 未改任何码、⛔ 未跑任何 go 命令 ⇒ ⓔ 全节一律标〔读码推的〕。
内容锚全部逐字抄自 `git show HEAD:<path>`；行号只作定位辅注，⛔ 绝对行号不许进派单（本仓规矩）。

---

# ⓓ 把 `started` 换成 `atomic.Bool` 之后，`cmd/wisp` 里必须跟着换形状的**每一处**

尺：`git grep -n "\.started" HEAD -- cmd/wisp`（含 `_test.go`、含注释行，目录 `cmd/wisp`）⇒ **HEAD 现量 11 枚点**（不是死腿的 12 枚，见 `20-named-corrections.md` D-1）。
分布＝生产 3 枚（2 写 + 1 读）＋测试 8 枚（全读）。

## ⓓ-1 生产面（`cmd/wisp/resident_audio_windows.go`，4 处含声明）

| 处 | 内容锚（逐字原文） | 线程类别 | 换形状后 |
|---|---|---|---|
| 声明 | `	started bool`（其上一枚字段是 `	levels atomic.Uint64`，同一 struct） | — | `	started atomic.Bool`（⚠ 位置若在 `muteMux`/`sync` 那类字段旁要重排 gofmt，本 struct 里没有锁，不受影响） |
| 写 1 | `	ra.started = gate.Open()`（在 `toggleMute` 体内，该体首行是 `func (ra *residentAudio) toggleMute() (string, bool) {`） | **`ui-sta`** | `ra.started.Store(gate.Open())` ← ★票 295 的正主：跨线程写 |
| 写 2 | `		ra.started = true`（`assembleCapture` 的 `default:` 支；同支紧邻三行是 `		ra.verdict = fmt.Sprintf("采集腿在跑：真麦克风 PCM 进，%s 每帧一枚 float32 交给球（路径 %s）",` / `			audio.FrameDuration.String(), string(gate.Path()))` / `		slog.Info("audio: capture leg running", …)`） | `boot` | `ra.started.Store(true)` |
| 读 1 | `	if !ra.started {`（`posture()` 体内，紧接下一行 `		return ra.unrunningClaim()`） | `boot`（`test` 也直调 `posture()`） | `if !ra.started.Load() {` |

**同批必改的注释（本仓 293 已定"旧句不许与新代码并存"）**：
- `toggleMute` 里写点上方那四行（首行逐字 `	// started mirrors what the gate now reports. It is derived state, not a`）＋声明上方那三行（首行逐字 `	// started records that the gate handed the device to the inner source and`）——
  两处的**措辞本身没问题**（说的是"镜像／派生态"），⚠ 但票面 `AC#1` 若要求写清"为什么这枚必须原子而那三枚不用"，落点就在这四行旁边，落地腿要指名。

**为什么换原子不会引 `copylocks`**（HEAD 现量证据，尺＝`git grep -n -E "residentAudio\{|range .*residentAudio|\*ra\b|ra := " HEAD -- cmd/wisp`）：
`residentAudio` 在整个 `cmd/wisp` 里只有**一枚**构造点 `	ra := &residentAudio{}`（`resident_audio_windows.go`，`assembleCapture` 体内），其余全是 `assembleCapture(...)` 返回的指针；
测试面**零枚** `&residentAudio{…}` 字面量（与 `residentBall` 相反——后者有 13 枚手搓字面量，死腿票 294 已复跑）。⇒ 无按值拷贝，且同 struct 已带 `atomic.Uint64` 这一事实本身就证明它今天没被拷过。

## ⓓ-2 测试面（8 枚读点，逐枚标出**条件位 vs 参数位**）

| 处 | 内容锚（逐字原文） | 位置类别 | 件（按名，不写绝对行号进派单） |
|---|---|---|---|
| 1 | `	if ra.started {`（后一行 `		t.Fatal("a leg with no collector reports itself running")`） | 条件位 | `cmd/wisp/resident_audio_247_windows_test.go`（voice.enabled=false 那枚用例） |
| 2 | `	if ra.started {`（后一行 `		t.Fatal("posture claims a running capture leg with the gate muted")`） | 条件位 | 同上件（默认出厂 posture 那枚用例） |
| 3 | `			if ra.started {`（后一行 `				t.Fatal("posture claims the capture leg is running over a failed device")`） | 条件位 | 同上件（设备失败分类表驱动的**子测试**内，缩进 3 层） |
| 4 | `	if !ra.started {`（后一行 `		t.Fatalf("the live capture leg did not come up on this machine: %s", ra.posture())`） | 条件位（真机件） | `cmd/wisp/resident_audio_247_live_windows_test.go` |
| 5 | `		if !ra.started {`（后一行 `			t.Fatalf("%s: the leg's derived mirror says nothing is running while the gate is open", gesture)`） | 条件位 ★**全场唯一一枚咬"手势路径必须写 started"的断言** | `cmd/wisp/resident_mute_290_windows_test.go`（`TestAC290BothMuteGesturesTurnTheGateAndBack`） |
| 6 | `		if ra.gate.Open() \|\| ra.started {` | 条件位（`\|\|` **右侧**，⚠ 只改左侧会漏） | 同上件、同一 for 循环第二轮 |
| 7 | `			t.Fatalf("%s: capture stayed handed over after muting (open=%v started=%v)", gesture, ra.gate.Open(), ra.started)` | ★**`t.Fatalf` 格式串参数位**（`%v` 直取），⛔ 不是条件位 | 同上件，紧接 #6 那条分支体 |
| 8 | `	if ra.started {`（后一行 `		t.Fatal("the leg's derived mirror says capture is running over a refused device")`） | 条件位 | 同上件（`TestAC290OutcomeIsReadOffTheGateNotOffTheRequest`） |

⚠ **#7 是这一节最容易被漏的一枚**：`atomic.Bool` 丢进 `%v` 会打出 `{true}`/`{false}` 结构体形状（编译过、断言不变红），但**失败日志开始撒谎**——而这枚用例的用途恰恰是把 `started` 的形状写进操作员可见的那一行。⇒ 落地腿必须 `.Load()`（并把 `%v` 收成 `%t` 更严），验收腿要现读那行输出。

## ⓓ-3 不改形状、但读数会被带动的间接面（缺口审计要点数）

`posture()` 是 `started` 在生产里唯一的读者，所以它的调用面全在射程边缘，逐枚列（尺＝`git grep -n "posture()" HEAD -- cmd/wisp`，含 `_test.go`＝**14 命中**，不含定义行 `resident_audio_windows.go:375` 与两处注释）：
- 生产：`cmd/wisp/resident_windows.go` 的 `	fmt.Printf("wisp: 采集腿：%s\n", raudio.posture())`（boot 报告那一行）。
- 测试：`resident_audio_247_windows_test.go` 的 `未构造`／`设备未打开`／`mic_muted_default`／`nilLeg.posture() == ""`／表驱动里 `tc.wantSaidSub` 那五段断言；`resident_audio_247_live_windows_test.go` 的 `t.Fatalf` 里那枚；`resident_mute_290_windows_test.go` 的 `设备未打开` 那枚。
⇒ 这些点⛔ 不需要换形状，但**它们读的是同一个 bool 的派生句子**；`AC#3`/`AC#5` 那类"boot 行不许变"的读数要从这里取。

## ⓓ-4 明确不在 ⓓ 清单里（防落地腿顺手扩）

- `verdict`：⛔ 不同批（`AC#0` 的前件不成立——`toggleMute` 体内对 `ra.verdict` 命中 0，本腿已复跑，见 `10-anchor-recheck.md` §1 #7）。
- `rb.b`：同族裸跨线程字段，⛔ 本票不碰（票面 `AC#1` 只许可 `verdict` 同批，死腿已具名"要修另立一票"）。
- `gate`/`mic`/`cancel`/`voiceEnabled`/`mutedAtBoot`：写点全在挂载前或零读者 ⇒ ⛔ 不换形状。
- ⛔ 不许顺手接 `mutedAtBoot` 当真相源（票 293 禁区①）。

---

# ⓔ 反形正控的落点预判（票面 `AC#2` ⓑ：把 `toggleMute` 里那句镜像写换成不写 ⇒ 指名用例必须红）

**变异体（按内容锚描述，不写行号）**：删掉 `toggleMute` 体内 `	ra.started = gate.Open()` 那一枚赋值，其余一字不动。
**前置事实（HEAD 现量，决定下面每一枚的敏感性）**：`assembleCapture` 只在 `default:` 支写 `started=true`；出厂默认 `mic_muted_default=true` ⇒ 静音 boot 走的是 `		ra.verdict = "采集腿已装配、门处于静音：设备未打开（[audio] mic_muted_default=true，来源 %s）"` 那一支，**不碰 `started`** ⇒ 该 rigs 里 `started` 初值＝零值 false。
〔以下全部＝**读码推的**，⛔ 非实测；真跑那一发归落地腿。〕

## ⓔ-1 真会红的：1 枚用例／1 枚断言

- `TestAC290BothMuteGesturesTurnTheGateAndBack`（`cmd/wisp/resident_mute_290_windows_test.go`）：rig 是 `muteRig`＝出厂默认配置 + 注入式可开门 source ⇒ `started` 初值 false；第一轮 `rb.muteGesture("mute-hotkey")` 把门打开后，断言 `		if !ra.started {` / `			t.Fatalf("%s: the leg's derived mirror says nothing is running while the gate is open", gesture)` **必须红**。
  ⚠ 顺序要紧：这枚断言在同一 for 轮里**排在** `if ra.gate.Open() || ra.started` 那枚**之前** ⇒ 删写变异一进来就先撞它，用例红＝红在这枚，⛔ 别把功劳记给后面那枚。
- **这是全场唯一一枚要求"手势路径写 `started`"的断言** ⇒ `AC#2` ⓑ 的指名用例只许是它。

## ⓔ-2 对这一变异**根本不敏感**的：其余 7 枚读点（判据换成反形也照样绿＝恒绿假钉形状）

| 读点（内容锚） | 为什么恒绿〔读码推的〕 |
|---|---|
| `mute_290` 的 `		if ra.gate.Open() \|\| ra.started {` ＋其 `t.Fatalf` 参数位 | 删写后 `started` 整轮停在 false ⇒ `false \|\| false`＝不成立，分支根本不进。**它对"删写"瞎，但对"恒 Store(true)"有牙**（那样第一轮 `!ra.started` 先红，见 ⓔ-4） |
| `mute_290` 的 `	if ra.started {`（`TestAC290OutcomeIsReadOffTheGateNotOffTheRequest`） | 该 rig 的 source 开局就返回 `DeviceError(0x80070005…)` ⇒ `gate.Open()` 永假，`started` 写与不写都是 false。这枚用例的真牙在 `	if ra.gate.Open() {` 与 outcome 串上的 `设备未交接`／`Privacy & security` |
| `mute_290` 的 `	if !strings.Contains(ra.posture(), "设备未打开")` 那枚 boot 用例（`TestAC290MutedDefaultStillDecidesTheBoot`） | 它**不发手势** ⇒ 变异点在 `toggleMute` 内，碰不到它；`started` 恒 false 走 `unrunningClaim()`，句子照旧 |
| `247_windows` 的 `	if ra.started {` ×3（voice 关／默认静音 boot／设备失败表） | 三枚都**不发手势**，期望的都是"boot 不许写 started"；变异只动 `toggleMute` ⇒ 全绿 |
| `247_live` 的 `	if !ra.started {` | 真机件、rig 是 `c.Audio.MicMutedDefault = false` ⇒ 走 boot 的 `default:` 支 `		ra.started = true`（未被变异），且不发手势 ⇒ 全绿（且这枚今天还依赖机主有麦克风） |

⇒ **本腿的具名警告**：票面 `AC#2` ⓑ 那发如果只报"我删了镜像写、用例红了"⛔ 不足以证明**哪一枚**断言在守这件事；反过来若有人拿 ⓔ-2 那 7 枚里任何一枚当"这条写点有测试在守"的凭据，那是**恒绿假钉**（本仓有先例）。指名只许写 `TestAC290BothMuteGesturesTurnTheGateAndBack`。

## ⓔ-3 换成 `atomic.Bool` 之后的"形状变异"敏感性预判〔读码推的〕

落地腿若要用反形证明"改形状这件事本身不是化妆"，能红的仍只有 ⓔ-1 那一枚；以下形状变异**全绿**，⛔ 别拿它们当正控：
- 把 `ra.started.Store(gate.Open())` 换成 `ra.started.Store(true)`：第一轮 `!ra.started` 绿，第二轮 `gate.Open() || started` 里 `started`＝true ⇒ 该用例**红**（这是 ⓔ-2 第一行说的"有牙"方向，红在 `t.Fatalf` 那条上——⚠ 所以 ⓓ-2 #7 那枚参数位的 `%v` 读数会进日志，落地腿别把它打印成结构体）。
- 把 boot 的 `ra.started.Store(true)` 换成不写：`247_live` 的 `	if !ra.started {` 红（真机件，非门禁日常面）＋ `mute_290` 无感 ⇒ **这条⛔ 不能用作 `AC#2` ⓑ 的正控**（它测的是 boot 腿不是手势腿）。

## ⓔ-4 一句话给落地腿

`AC#2` ⓑ 的正控射程在 HEAD 上只有**一枚断言宽**（`TestAC290BothMuteGesturesTurnTheGateAndBack` 里那句"derived mirror says nothing is running while the gate is open"）⇒ 改完形状后**必须现跑那一枚**（`go test -run` 指名到用例），并把 ⓔ-2 那 7 枚逐枚具名标"对此变异不敏感"，否则缺口审计会以为 `started` 有 8 枚测试在守。

---

## 本笔自报

- ⛔ 零源码改动、⛔ 零翻框（票 295 六格 `AC#0`..`AC#5` 一枚没碰，⛔ 没给自己翻任何一格）、⛔ 未动 `docs/**`／`frontend/**`／`design/**`／`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／allowlist／三枚冻结件。
- go 命令：**零枚**（`go env`／`go list` 均未跑）；编译面/门禁/进程：零枚。
- ⚠ 与本票无关但撞见的：`295-a1` 那份件的 ⓒⓓⓔ 三节标题在原文里不存在（`ls .scratch/wisp/probes/295/a1/` 现跑＝1 个文件）⇒ 编排者"从未产出"的转述属实。
