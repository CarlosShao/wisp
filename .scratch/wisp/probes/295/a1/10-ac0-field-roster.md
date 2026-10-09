# 票 295 · 只读普查腿 `295-a1` · ⓐⓑ 字段名册（AC#0 / AC#1 的"真邻居"有几枚）

**立件时刻** `2026-10-09 21:44:53 +0800`（`date` 的 stdout，非手打）
**锚点** `448f5a57`，branch `dev`

## 射程声明（⚠ 本轮新规矩，逐字写在最前）

- **本文所有行号一律＝HEAD 对象层，不是工作树。** 尺＝`git show HEAD:<path> | cat -n | sed -n '<A>,<B>p'` 与 `git grep -n -o "<形状>" HEAD -- <路径>`。
- **盘上现量我没取，也不声称。** 起手现跑 `git status --porcelain`＝工作树是脏的（`M .gitignore`、`M .scratch/wisp/probes/**` 多枚、`D design/assets/*` 等 20+ 行），编队里确有写腿在飞。⇒ 任何人拿本文行号去对工作树，漂移不代表本文错。
- **零 go 命令**（`go build`/`vet`/`test`/`list` 一枚没跑；**`go env` 也没跑**，本文不需要它）。`go 1.27` 那一枚是从 `git show HEAD:go.mod` 第 3 行读的，不是从 `go env`。
- ⛔ 零源码改动、⛔ 零翻框、⛔ 未动 `docs/**`。
- 收件前先跑：`wc -c` 本件、`git log -1` 本腿那两笔。

## 线程词汇表（后文"线程"列全部用它，⛔ 不引入新名）

| 代号 | 是哪条线程 | 尺（HEAD 逐字） |
|---|---|---|
| `boot` | `runResident` 自己那条 goroutine；装配期写、boot 报告读、退出期两枚 defer 也在这条上跑 | `cmd/wisp/resident_windows.go:111` 逐字 `		records := rt.Shutdown(false)`（在 `:110` 的 defer 里）· `:253` 逐字 `	defer rb.stop()` · `:260` 逐字 `	defer ra.detachBall()` |
| `ui-sta` | internal/ball 的球 STA 线程；两枚静音手势都在这条上同步开火 | `internal/ball/ball_windows.go:661` 逐字 `			b.fire(b.opts.Events.OnMuteHotkey)`（`case wmHotkey:` 支，`:656`）· `:679` 逐字 `				b.fire(b.opts.Events.OnTrayMute)`（`case wmAppTray:` → `showMenu` 返回后，`:674`/`:678`）· `:728-731` 的 `func (b *Ball) fire(fn func()) { if fn != nil { fn() } }` ＝**内联同调，不起新线程** |
| `capture` | 电平那一侧的采集线程 | `internal/audio/captureopt.go:110` 逐字 `	c.OnLevel(level)`；`cmd/wisp/resident_windows.go:311` 逐字 `	raudio := startResidentAudio(rt, rb.setAudioLevel)`；本仓 commit `5d407e77` 标题正文逐字 "level sink emits one float32 inside the audio-capture goroutine" |
| `test` | 用例自己那条 goroutine（现量：三枚读 `started` 的件里**没有任何** `go func`/`WaitGroup`/`t.Parallel`，尺见 ⓒ 件 §1） |

⚠ **`ui-sta` 这一判我判得了，靠的是三处逐字**：`cmd/wisp/resident_ball_windows.go:245-246`（`// ball's global hotkeys on the ui-sta thread internal/ball owns (that thread is` / `// the ball's own Registry.Spawn with the frozen D38b resident name "ui-sta",`）、`:501`（`// Threading, because it is not free: callbacks run ON the ui-sta thread and`）、`internal/ball/ball_windows.go:734-736`（`// uiRun runs fn on the UI thread and waits for it to finish. A caller that is` / `// ALREADY on the UI thread - which is every Events callback, ...`）。
⚠ **托盘那枚 `OnTrayMute` 我也判成 `ui-sta`**：尺＝`b.fire` 在 `wmAppTray` 的 WndProc 里（`:669`/`:679`），WndProc 跑在球自己的 STA 线程上 ⇒ **两枚手势写者同线程 ⇒ `started` 没有 write/write 竞争，只有 write(ui-sta)/read(boot) 一枚对**。这条是给 `295-r1` 的形状前提。

---

# ⓐ 名册 A：`residentAudio`（定义窗 `cmd/wisp/resident_audio_windows.go:86-113`）

尺＝`git show HEAD:cmd/wisp/resident_audio_windows.go | cat -n | sed -n '86,113p'` 整窗读；调用点尺＝`git grep -n -o "\bra\.<字段>\b" HEAD -- cmd/wisp` **再加**一把"其他持有者变量名"尺（本类型的生产持有者叫 `raudio`，见 `resident_windows.go:311`；测试持有者叫 `ra`/`nilLeg`）。

⚠ **锚点陷阱（本腿现撞到，具名报回）**：**接收者名 `ra` 在本包里被两枚不同结构体共用**——`residentAudio`（resident_audio_windows.go）与 `residentApproval`（resident_approval_windows.go）。⛔ 只锚 `ra.<字段>` 会串号：实测 `ra.gate`／`ra.cancel`／`ra.verdict` 三个符号名在 `resident_approval_windows.go` 里各有大量命中，那是**另一枚 struct 的同名字段**。⇒ 下表每一行的射程都限定在"类型确实是 `*residentAudio` 的那些点"，尺＝文件＋定义窗＋调用点人工对齐（对齐依据：`resident_approval_windows.go` 里 `ra` 的声明是 `:167` 逐字 `	ra := newResidentApprovalWithConfig(rt.Layout.DataDir)`，`resident_windows.go:311` 才是 `raudio`）。

| # | 字段（定义行逐字） | 写点 → 线程 | 读点 → 线程 | 本行判定 |
|---|---|---|---|---|
| 1 | `:89` `	gate *audio.HalfDuplexGate` | `:277` `	ra.mic, ra.gate, ra.cancel = mic, gate, cancel` → `boot`（唯一写点） | `:157` `	if ra == nil || ra.gate == nil || ra.mic == nil {`、`:161` `	gate := ra.gate` → `ui-sta`（`toggleMute` 体内）；`:350`、`:356` `	err := ra.gate.Stop()` → `boot`（`stop()` 走 `:336` 注册的 StepStopAudio 钩子，钩子由 `:111` defer 的 `rt.Shutdown` 在 boot 线程上走）；测试面多处 → `test`（`resident_mute_290_windows_test.go:89`/`:108`/`:118`/`:121`/`:138`/`:141`/`:142`/`:171`/`:192`/`:252`、`resident_audio_247_windows_test.go:110`/`:146`/`:149`/`:152`/`:187`/`:253`/`:298`/`:361`） | **非原子、无锁**，但"写全在 attach 之前 + attach 走 `muteMux`"给了发布序 ⇒ 今天不构成竞争。**读点跨线程（`ui-sta`）这件事本身成立**，只是不再被写 |
| 2 | `:93` `	mic    captureSource` | `:277`（同上一条那一枚三元赋值）→ `boot` | `:157` → `ui-sta`；`:180` `		detail := ra.mic.Stats().LastError`、`:181` `		if devErr := ra.mic.Err(); devErr != nil {` → `ui-sta`；`:357` `	st := ra.mic.Stats()` → `boot`；测试 `resident_audio_247_windows_test.go:110`/`:158`、`resident_audio_247_live_windows_test.go:151` `	mic := ra.mic.(*audio.WASAPIMicrophone)` | 同上：**非原子、无锁、单一写点、靠发布序**。`mic` 背后的 `Stats()`/`Err()` 自带锁（夹具那侧逐字见 `resident_mute_290_windows_test.go:37-72`，真源自 `internal/audio/wasapimic_windows.go`），⛔ 别把"方法里有锁"读成"字段读有护栏" |
| 3 | `:94` `	cancel context.CancelFunc` | `:277` → `boot` | `:353` `	if ra.cancel != nil {`、`:354` `		ra.cancel()` → `boot`（`stop()`，同 #1 的钩子路径）；测试面 `resident_grant_writer_265_windows_test.go:585` `	t.Cleanup(ra.cancel)` → ⚠ **这一处我判不了**：那行上的 `ra` 是 `residentApproval` 还是 `residentAudio`？我按类型对齐尺判＝`resident_grant_writer_265_windows_test.go` 里没有任何 `assembleCapture`/`startResidentAudio` 调用，故**这一枚不是本类型的点**，不计入本行 | **单线程（boot 写、boot 读）** |
| 4 | `:100` `	voiceEnabled bool` | `:257` `	ra.voiceEnabled = c.Voice.Enabled` → `boot` | **0 枚（全仓零读者）** | 死字段。尺＝`git grep -n -E "voiceEnabled\|mutedAtBoot" HEAD` ⇒ 源码命中只有 `:100`/`:101` 声明与 `:256`/`:257` 两处写，读＝0。**★这不是本腿新发现的**：`293-a1` 已具名 census 过（件 `.scratch/wisp/probes/293/a1/30-ac0c-second-surfacer-scan.md:68` 与 `.scratch/wisp/probes/293/a1/90-unrun-rulers.md:56` 第 17 条），编排者已把它写成票 293 禁区①（`.scratch/wisp/issues/293-tray-mute-checkmark-never-follows-the-gate.md:48` 逐字 `① ⛔ **不许用 `ra.mutedAtBoot` 当真相源**（死字段、`toggleMute` 不更新它`）。本表只登记它在 ⓐ 名册里的位置＝**非原子、零竞争面** |
| 5 | `:101` `	mutedAtBoot  bool` | `:256` `	ra.mutedAtBoot = c.Audio.MicMutedDefault` → `boot` | **0 枚** | 同上，死字段；⛔ 不得被 295 顺手接成真相源（票 293 禁区① 已钉） |
| 6 | `:105` `	started bool`（`:102-104` 注释逐字开头 `// started records that the gate handed the device to the inner source and`） | **两枚写点**：`:323` `		ra.started = true` → `boot`；`:168` `	ra.started = gate.Open()` → **`ui-sta`**（`toggleMute` 体 `:156-193` 内，`:156` 逐字 `func (ra *residentAudio) toggleMute() (string, bool) {`） | **唯一生产读点**：`:379` `	if !ra.started {` → `boot`（`posture()`，`:375` 逐字 `func (ra *residentAudio) posture() string {`；生产调用点唯一＝`resident_windows.go:353` `	fmt.Printf("wisp: 采集腿：%s\n", raudio.posture())`）。测试读点 → `test`：`resident_mute_290_windows_test.go:124`/`:141`/`:142`/`:174`、`resident_audio_247_windows_test.go:113`/`:155`/`:256`、`resident_audio_247_live_windows_test.go:148` | ★**本名册里唯一一枚"两线程都有写/读、无原子无锁"的字段**＝票 295 的那一发。尺＝`git grep -n -o "\.started\b" HEAD -- cmd/wisp`＝**12 命中**（声明外：写 2 枚、读 10 枚，其中生产读 1 枚） |
| 7 | `:112` `	verdict string` | **6 枚写点，全在 `boot`、全在 `assembleCapture` 返回之前**：`:263`、`:284`、`:298`、`:313`、`:320`、`:324` | 读点：`:321` `slog.Warn(..., "posture", ra.verdict)` → `boot`；`:382` `	return ra.verdict + "；帧消费侧（ASR/KWS，internal/speech 一块没写）不属于本票，" +`（`posture()`）→ `boot` + `test`；**:389-390** `	if ra != nil && ra.verdict != "" {` / `		return ra.verdict`（`unrunningClaim()`，`:388`）→ **`ui-sta`**（`toggleMute:158` 逐字 `		return ra.unrunningClaim(), false` 是唯一入口）；测试读点 `resident_audio_247_windows_test.go:147`/`:259`/`:260`/`:262`/`:263`/`:362`、`resident_mute_290_windows_test.go:90` | ★**AC#0 那一句必答题的答案＝手势路径**不**写 `verdict`，只**读**它。** 尺＝上面 6 枚写点逐枚对文件（`:263`–`:324` 全在 `assembleCapture` 内，`toggleMute` 体 `:156-193` 里对 `verdict` 的命中＝**0**，只有 `:158` 那枚间接读）⇒ **`verdict` 不满足票面 `AC#1` 那句"若判出也会被手势路径写 ⇒ 同批改"的前件 ⇒ 295 不该同批改它**。它的跨线程**读**今天由 `muteMux` 的发布序盖住（同 #1） |
| 8 | `:110` `	levels atomic.Uint64` | `:122` `			ra.levels.Add(1)` → `capture`（`levelOut` 的闭包被 `internal/audio/captureopt.go:110` 调） | `:359` `		ra.levels.Load(), st.FramesSent, st.FramesDropped, st.Reopens)` → `boot`（`stop()`）；测试 `resident_audio_247_windows_test.go:197`/`:322`/`:366`/`:373`、`resident_audio_247_live_windows_test.go:182` | **已护**＝本 struct 里**唯一**一枚原子；票面 `:110` 那一枚邻居行号在 HEAD **逐字对上** |

**"所有非原子字段"一句话总账**：`residentAudio` 共 8 枚字段，其中 7 枚非原子（`gate`/`mic`/`cancel`/`voiceEnabled`/`mutedAtBoot`/`started`/`verdict`），只有 `levels` 是原子家族。⇒ 但"非原子"≠"有竞争"：**跨线程被读**的有 3 枚（`gate`、`mic`、`verdict`，全在 `ui-sta` 的 `toggleMute` 路径上读）**且它们的写点全在挂载之前**；**跨线程被写**的只有 1 枚＝`started`。

---

# ⓐ 名册 B：`residentBall`（定义窗 `cmd/wisp/resident_ball_windows.go:133-180`）

尺＝整窗读 + `git grep -n -o "rb\.[a-zA-Z]*" HEAD -- cmd/wisp`（⚠ 这一把尺的**噪声面**：同形还命中 `cmd/wisp/secret_test.go`/`run_mode101_test.go` 里另一枚 `rb`（`rb.String`），那与本类型无关，已剥）。

| # | 字段（定义行逐字） | 写点 → 线程 | 读点 → 线程 | 本行判定 |
|---|---|---|---|---|
| 1 | `:134` `	b *ball.Ball` | `:338` `	rb.b = b` → `boot`；**`:435` `	rb.b = nil`** → `boot`（`stop()`，注册在 `resident_windows.go:253` `	defer rb.stop()`）；测试面 13 枚手搓 `&residentBall{...}` 字面量（票 294 已复跑枚数） | `:352` `	if hotReload != nil && rb.b != nil {` → `boot`；`:431`/`:434` → `boot`（stop）；**`cmd/wisp/resident_audio_windows.go:136` `	if rb == nil || rb.b == nil {` 与 `:139` `	rb.b.SetAudioLevel(level)`** → **`capture`**；`resident_approval_windows.go:547`/`:567`（`bindBallHost`）→ `boot`；测试大量 → `test` | ★**本名册里第二枚同形状跨线程裸字段**，且**不是票 290 造的**（`setAudioLevel` 是票 247 那批，见 commit `5d407e77`）。竞争窗＝`ui-sta`/`boot` 侧 defer 先跑 `rb.b = nil`（`:435`），而**采集线程要到 D38(e) 第 4 步才停**（`resident_audio_windows.go:336` 注册 `proc.StepStopAudio`；`resident_ball_windows.go:412-414` 逐字 `// runResident registers it as the defer that runs BEFORE rt.Shutdown walks the` / `// frozen sequence, so the hot keys are released ahead of the Job Object close`）⇒ defer 序（LIFO：`:260` detachBall → `:253` rb.stop → `:111` rt.Shutdown）在 HEAD 上确实把"置 nil"排在"停采集"之前。⇒ **具名报编排者：这枚不属本票射程**（票面 `AC#1` 只许可 `verdict` 同批，且 `AC#1` ⛔ 不许自造第三种包装）；要修得另立一票 |
| 2 | `:140` `	cancelHosted bool` | `:297` `	rb := &residentBall{cancelHosted: onCancelEsc != nil}` → `boot` | `resident_approval_windows.go:555` → `boot`（`bindBallHost` 体内）；测试（`resident_approval_246_windows_test.go:331`、`resident_cancel_key_*` 两枚字面量） | 单线程装配期；⚠ 测试面有直接构造写（字面量），不构成竞争 |
| 3 | `:146` `	verdict string` | `:332`、`:385` → `boot` | `:334` → `boot`；`:404` `	return rb.verdict`（`statusLine()`）→ `boot`（`resident_windows.go:364`/`:367`）；`resident_hotkey_258_seam_windows.go:19-20` 的 `verdictStatusLineFor258` → `test`；测试直读 `resident_hotkey_v1probe_test.go:104-108` | 单线程（boot 写 boot 读）；**手势路径不读它**（`muteGesture` 体内 `rb.verdict` 命中 0，尺＝上表 #1 那一把 grep 的全部命中对齐文件） |
| 4 | `:152` `	hotkeyProvenance string` | `:308` `	rb.hotkeyProvenance = provenance` → `boot` | `:601` `	return rb.hotkeyProvenance`（`hotkeySourceName()`，`:597`）→ `boot`/`test` | 单线程 |
| 5 | `:157` `	hotkeyBridgeArmed bool` | `:368` `		rb.hotkeyBridgeArmed = true` → `boot` | **源码零读者**；只有测试读：`resident_hotkey_258_windows_test.go:283`/`:346`、`resident_hotkey_296_windows_test.go:163`/`:208`、`resident_hotkey_v1probe_test.go:62` → `test` | 非原子；单写、读全在 test 线程（且是 boot 完成之后直调）⇒ 无竞争面 |
| 6 | `:161` `	hotkeyBridgeStop func()` | `:358` `		rb.hotkeyBridgeStop = func() {` → `boot`；**`:429` `			rb.hotkeyBridgeStop = nil`** → `boot`（stop 内二次写） | `:427`/`:428` → `boot`（stop） | **两枚写点但同线程（boot）**⇒ 不是竞争；⚠ 记一句形状：这枚字段的"写两次"形状与 `started` 的"写两次"形状不同源（同线程 vs 跨线程），别混进同一批改 |
| 7 | `:165` `	hotkeyBridge *ballHotkeyBridge258` | `:365` `		rb.hotkeyBridge = &ballHotkeyBridge258{check258: func() { bridge.Check() }}` → `boot` | `resident_hotkey_258_seam_windows.go:11`/`:14` → **`test`**（那枚文件是测试缝：`:10` 逐字 `func (rb *residentBall) hotkeyBridgeCheck258() {`） | 生产读点 0 枚 ⇒ 唯一读者是测试缝，且 `bridge.Check()` 自己投给 STA（`internal/ball` 那侧），字段读本身在 test 线程 |
| 8 | `:178` `	muteMux  sync.Mutex` | — | — | **它就是 #9 的护栏**；`:173-177` 注释逐字 `	// muteMux is what makes that late attach visible to the ui-sta thread: the` / `	// hot key is live from the moment the window comes up, so a plain field` / `	// would be a race between the boot goroutine writing it and a key press` / `	// reading it. The attach also orders everything the boot wrote into the` / `	// residentAudio handle before any gesture may read it.` ⇒ ★**最后那一句是票 295 `AC#1` 最该引的先例**：那把锁**管的是发布序**，⛔ **它不覆盖挂载之后再写的那一枚**（`started` 的 `:168` 正是挂载之后的写）⇒ 票面现量第 13 条那句"⛔ 它不覆盖 `ra.started` 的写读对"**成立** |
| 9 | `:179` `	muteGate muteGestureFunc` | `:478` `	rb.muteGate = fn` → `boot`（`attachMuteGate`，`:476-477` 先 `rb.muteMux.Lock()`/`defer …Unlock()`；调用点＝`resident_windows.go:336` `	rb.attachMuteGate(raudio.toggleMute)` 与测试 `resident_mute_290_windows_test.go:93`/`:206`/`:235`） | `:489` `	return rb.muteGate` → **`ui-sta`**（`currentMuteGate`，`:487-488` 同锁；入口＝`muteGesture` `:516` `	fn := rb.currentMuteGate()`，被 `resident_ball_windows.go:321`/`:325` 两枚闭包递给球） | **已护（带锁）**＝`residentBall` 里唯一有锁的字段。⚠ 一枚例外要点名：测试 `resident_mute_290_windows_test.go:228` `	if rb.muteGate != nil {` **绕过锁直读那枚字段**——同线程、无竞争，但它是"护栏字段被裸读"的形状，`295-r1` 若在这带加断言别再抄这一形 |

---

# ⓑ 哪些已有护栏 ⇒ `AC#1` 的"照邻居的形状来"到底有几个真邻居

| 护栏形状 | 真邻居（逐字行） | 枚数 |
|---|---|---|
| **`atomic` 家族** | `residentAudio:110` 逐字 `	levels atomic.Uint64`（同 struct 邻居，票面点名的那一枚）⇒ **写者线程 `capture`、读者线程 `boot`，正与 `started` 的目标形状同形** | **1 枚**（同 struct 内唯一） |
| 同包其他 `atomic` 先例 | `cmd/wisp/panel_resident_windows.go:136` `stopRequested atomic.Bool`（`atomic.Bool` 生产形状先例）；`cmd/wisp/panel_resident_windows.go:126-129` 四枚 `atomic.Int*` | 生产 `atomic.Bool` 在本仓＝**2 枚**：`panel_resident_windows.go:136` 与 `internal/ball/ball_windows.go:132` 逐字 `	closed     atomic.Bool` |
| ⚠ 票面 `AC#1` 那句"全仓已有 **5 处**在用 `atomic.Bool`" | 尺我复跑＝`git grep -n "atomic.Bool" HEAD -- internal cmd`＝**5 命中，行号逐枚对上**：`cmd/wisp/panel_host_windows_test.go:42`、`cmd/wisp/panel_resident_windows.go:136`、`internal/audio/hotplug_test.go:85`、`internal/ball/ball_windows.go:132`、`internal/ball/live_guard_windows_test.go:69` | **枚数成立，但口径要钉**：5 命中里 **3 枚在 `_test.go`** ⇒ **生产先例只有 2 枚**。票面那把尺没写含不含测试面（同 294 编排者裁过的"名册尺必写含不含注释行"那一条）。⛔ 这不是票面错，是**射程口径缺失**，报编排者记账 |
| **带锁** | `residentBall.muteMux`/`muteGate` 那一对：`:476-478`、`:487-489` | **1 对**（⛔ 票面 `AC#1` 明令"不加锁"，且这把锁只管发布序） |
| **只在装配期写（单线程）** | `residentAudio`：`gate`/`mic`/`cancel`（写点唯一＝`:277`，boot）、`verdict`（写点 6 枚全在 `assembleCapture` 返回前）、`voiceEnabled`/`mutedAtBoot`（写点各 1，零读者）；`residentBall`：`cancelHosted`(:297)、`verdict`(:332/:385)、`hotkeyProvenance`(:308)、`hotkeyBridgeArmed`(:368)、`hotkeyBridge`(:365)、`hotkeyBridgeStop`(:358/:429 同线程) | **12 枚**（含 2 枚死字段） |
| **裸且跨线程** | ①`residentAudio.started`（写 `:168` ui-sta / 读 `:379` boot）②`residentBall.b`（写 `:435` boot / 读 `resident_audio_windows.go:136`+`:139` capture） | **2 枚**；⇒ **`AC#1` 的"真邻居"＝1 枚（`levels`），形状＝`atomic.Bool`；`AC#0` 候选两枚里只有 `started` 成立，`verdict` 不成立（见名册 A #7）** |

**给 `295-r1` 的三条具名前提（本腿只报形状，不代它裁）**
1. 改 `started` 只需动**一枚 struct + 两枚写点 + 一枚生产读点**：`:105` 声明、`:168`、`:323`、`:379`；测试面另有 **7 处读点**要跟着换形状（`247_windows:113`/`:155`/`:256`、`247_live:148`、`mute_290:124`/`:141`/`:142`/`:174`＝8 处含 `:142` 的 `t.Fatalf` 参数位），⚠ `:142` 那一枚是**格式串参数**不是条件位，`atomic.Bool` 的 `.Load()` 要显式补上，别漏。
2. `verdict` ⛔ 不同批（前件不成立）；`rb.b` 是同族但**另票**（本腿⛔ 不自立票、⛔ 不改码）。
3. 两枚写者里 `:323` 在 `boot`、`:168` 在 `ui-sta` ⇒ 换成 `atomic.Bool` 后**写/write 序仍不存在**（两枚写点谁后跑谁说了算，票面禁区"派生镜像"的含义＝`:168` 之后 `:323` 不该再被重跑；本腿只具名指出"两枚写点无仲裁"这件事今天也在，⛔ 不是本票要修的语义）。
