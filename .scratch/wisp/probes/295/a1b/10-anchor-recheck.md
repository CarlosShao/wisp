# 票 295 · 补位腿 `295-a1b` · 第 2 笔：ⓐⓑ 两张表逐枚锚点复跑（HEAD 对象层）

HEAD `29081a1367fa4024716378f1b18486fc97b6edbf`（branch `dev`）；死腿件＝`.scratch/wisp/probes/295/a1/10-ac0-field-roster.md`（提交 `3f9bedf7fa309da74a2c0320c32937b423c13f6b`，自称锚 `448f5a57`）。
尺与射程见 `00-anchor-and-rulers.md`。**结论先行**：`git diff --stat 448f5a57 HEAD -- cmd/wisp` **输出为空** ⇒ 两锚之间 `cmd/wisp/**` 一字未动 ⇒
死腿 ⓐⓑ 表里的**行号全部逐字对上 HEAD**（"树会呼吸"这一发在两枚定义窗上没咬到人）；对不上的**不是行号，是枚数与两处文件/接收者归属**，逐枚见 `20-named-corrections.md`。

## 0. 并排两把读数（HEAD ↔ 工作树）

| 尺（逐字） | HEAD 读数 | 工作树读数 |
|---|---|---|
| `git grep -n -o "\.started\b" HEAD -- cmd/wisp \| wc -l` | **11** | `git grep -n -o "\.started\b" -- cmd/wisp \| wc -l` = **11** |
| 同尺在死腿自称锚 `448f5a57` | **11**（行号逐枚与 HEAD 同） | — |
| `git status --porcelain -- cmd/wisp` | 空（干净） | 同一条命令现跑＝空 |

## 1. ⓐ 名册 A：`residentAudio`（定义窗 HEAD 现量 `:86-:113` ✓＝`:86 type residentAudio struct {` / `:113 }`）

字段声明行 8 枚全部逐字对上：`:89` `	gate *audio.HalfDuplexGate`、`:93` `	mic    captureSource`、`:94` `	cancel context.CancelFunc`、`:100` `	voiceEnabled bool`、`:101` `	mutedAtBoot  bool`、`:105` `	started bool`（`:102-104` 注释首行逐字 `	// started records that the gate handed the device to the inner source and` ✓）、`:110` `	levels atomic.Uint64`、`:112` `	verdict string`。

| # | 字段 | 写点（内容锚，逐字）→ 线程 | 读点（内容锚，逐字）→ 线程 | HEAD 判定 |
|---|---|---|---|---|
| 1 | `gate` | `resident_audio_windows.go:277` `	ra.mic, ra.gate, ra.cancel = mic, gate, cancel` → `boot`（唯一写点 ✓） | `:157` `	if ra == nil || ra.gate == nil || ra.mic == nil {`、`:161` `	gate := ra.gate` → `ui-sta`；`:350` `	if ra == nil || ra.gate == nil {`、`:356` `	err := ra.gate.Stop()` → `boot`（`stop()`＝`:349`，钩子在 `:336` `	if err := rt.RegisterShutdownHook(proc.StepStopAudio, func(context.Context) error {` 注册）；测试面 10+8 枚逐条对上（见本腿盘上名册） | ✓ 与死腿一致（非原子、无锁、单一写点、跨线程读） |
| 2 | `mic` | 同 `:277` 那一枚三元赋值 → `boot` | `:157` → `ui-sta`；`:180` `		detail := ra.mic.Stats().LastError`、`:181` `		if devErr := ra.mic.Err(); devErr != nil {` → `ui-sta`；`:357` `	st := ra.mic.Stats()` → `boot`；测试 `resident_audio_247_windows_test.go:110` `	if ra.gate != nil \|\| ra.mic != nil {`、`:158` `	if st := ra.mic.Stats(); st.FramesSent != 0 \|\| st.BytesSent != 0 {`、`resident_audio_247_live_windows_test.go:151` `	mic := ra.mic.(*audio.WASAPIMicrophone)` | ✓ 逐枚对上 |
| 3 | `cancel` | `:277` → `boot` | `:353` `	if ra.cancel != nil {`、`:354` `		ra.cancel()` → `boot` | ✓ 死腿的"那一枚不算本类型"结论**成立但理由错了**：`resident_grant_writer_265_windows_test.go:585` `	t.Cleanup(ra.cancel)` 上的 `ra` 由同件 `:584` `	ra := newResidentApproval()` 构造 ⇒ 类型＝`*residentApproval`（`resident_approval_windows.go:365` `	ra := &residentApproval{root: root, cancel: cancel}`）。死腿给的对齐尺指向 `resident_approval_windows.go:167`，HEAD 现量那一行是注释 ⇒ 见更正 D-3 |
| 4 | `voiceEnabled` | `:257` `	ra.voiceEnabled = c.Voice.Enabled` → `boot` | **0 枚**（尺限定 `-- cmd/wisp`、含 `_test.go`：HEAD 只有 `:100` 声明 + `:257` 一处写） | ✓ 死字段成立；⚠ 死腿那把尺写作全仓 `HEAD`，未限定目录 ⇒ HEAD 现量另有 12 枚 `.scratch/**` 与 `docs/**` 的 `.md` 命中（口径病，见 D-8） |
| 5 | `mutedAtBoot` | `:256` `	ra.mutedAtBoot = c.Audio.MicMutedDefault` → `boot` | **0 枚** | ✓ 同上 |
| 6 | `started` | 两枚：`:323` `		ra.started = true` → `boot`（`assembleCapture` 的 `default:` 支，函数体 `:229`–`:347`）；`:168` `	ra.started = gate.Open()` → `ui-sta`（`toggleMute` 体＝`:156`–`:193` ✓） | 生产唯一读点 `:379` `	if !ra.started {` → `boot`（`posture()`＝`:375`；生产调用点唯一＝`resident_windows.go:353` `	fmt.Printf("wisp: 采集腿：%s\n", raudio.posture())` ✓）；测试 8 处读点逐枚对上（`247_windows:113`/`:155`/`:256`、`247_live:148`、`mute_290:124`/`:141`/`:142`/`:174`） | ★成立＝**唯一一枚跨线程裸写字段**；⚠ 枚数：HEAD 现量 **11 命中**（写 2／读 9），死腿写 12／读 10 ⇒ 见 D-1 |
| 7 | `verdict` | 6 枚全在 `boot`、全在 `assembleCapture` 内：`:263`、`:284`、`:298`、`:313`、`:320`、`:324`（逐枚原文均 `	ra.verdict = …` / `			ra.verdict = …`，本腿已核） | `:321` `		slog.Warn("audio: gate reports not open without a device error", "posture", ra.verdict)` → `boot`；`:382` `	return ra.verdict + "；帧消费侧（ASR/KWS，internal/speech 一块没写）不属于本票，" +`；`:389` `	if ra != nil && ra.verdict != "" {` / `:390` `		return ra.verdict`（`unrunningClaim()`＝`:388`）→ `ui-sta`（入口 `:158` `		return ra.unrunningClaim(), false`；★补一枚死腿未点名的第二入口：`:380` `		return ra.unrunningClaim()`，在 `posture()` 内 → `boot`/`test`）；测试读点 `247_windows:147`/`:259`/`:260`/`:262`/`:263`/`:362`、`mute_290:90` ✓ | ✓ **`verdict` 不成立**这一判定成立：`toggleMute` 体（`:156`–`:193`）内 `ra.verdict` 命中＝**0**（HEAD 全文 `verdict` 命中只有上面那批＋注释 `:189` 那句字符串里的"verdict"字样＋`:197-198` 注释 ⇒ 非字段访问） |
| 8 | `levels` | `:122` `			ra.levels.Add(1)` → `capture`（`levelOut` 闭包＝`:119`，电平回调原文锚 `internal/audio/captureopt.go:110` `	c.OnLevel(level)` ✓） | `:359` `		ra.levels.Load(), st.FramesSent, st.FramesDropped, st.Reopens)` → `boot`（格式串在 `:358`）；测试 `247_windows:197`/`:322`/`:366`/`:373`、`247_live:182` ✓ | ✓ 同 struct 内唯一原子 ✓（定义窗 8 枚里只它带 `atomic`） |

**总账复核**：8 枚字段／7 枚非原子 ✓；跨线程被读 3 枚（`gate`/`mic`/`verdict`）✓；跨线程被写 1 枚（`started`）✓。

## 2. ⓐ 名册 B：`residentBall`（定义窗 HEAD 现量 `:133-:180` ✓＝`:133 type residentBall struct {` / `:180 }`）

| # | 字段 | 写点 → 线程 | 读点 → 线程 | HEAD 判定 |
|---|---|---|---|---|
| 1 | `:134` `	b *ball.Ball` | `:338` `	rb.b = b`、`:435` `	rb.b = nil` → `boot`（`stop()`＝`:418`） | `:352` `	if hotReload != nil && rb.b != nil {`、`:431` `	if rb.b == nil {`、`:434` `	rb.b.Close()` → `boot`；`resident_audio_windows.go:136` `	if rb == nil \|\| rb.b == nil {` 与 `:139` `	rb.b.SetAudioLevel(level)` → **`capture`**（`setAudioLevel`＝`resident_audio_windows.go:135`，装配期递出处 `resident_windows.go:311` `	raudio := startResidentAudio(rt, rb.setAudioLevel)` ✓）；`resident_approval_windows.go:547` `	if rb == nil \|\| rb.b == nil {` 与 `:567` `	ra.ui.attachBall(rb.b)` → `boot`（`bindBallHost`＝该文件 `:546`） | ✓ 第二枚同形状裸跨线程字段，**另票**（死腿已具名不扩射程） |
| 2 | `:140` `	cancelHosted bool` | `:297` `	rb := &residentBall{cancelHosted: onCancelEsc != nil}` → `boot` | `resident_approval_windows.go:555` `	if !rb.cancelHosted {` → `boot`；测试真读点＝`resident_approval_246_windows_test.go:315` `	if rb.cancelHosted {`（死腿未列）；构造字面量 HEAD 现量 **4 枚**（`246:331`、`label_260r4:89`、`wording_260r3:99`、`wording_260r3:141`） | ⚠ 死腿把构造当读点、枚数少一枚 ⇒ 见 D-6 |
| 3 | `:146` `	verdict string` | `:332`、`:385` → `boot` | `:334` `		fmt.Printf("wisp: %s\n", rb.verdict)`、`:404` `	return rb.verdict`（`statusLine()`＝`:400`；调用点 `resident_windows.go:364` 与 `:367` ✓）；测试 `resident_hotkey_v1probe_test.go:104`/`:105`/`:107`/`:108` ✓ | ✓ 手势路径不读它：HEAD 现量 `resident_ball_windows.go` 里 `rb.verdict` 只有 `:332`/`:334`/`:385`/`:404` 四枚，`muteGesture`（`:515`）体内 **0** ✓ |
| 4 | `:152` `	hotkeyProvenance string` | `:308` `	rb.hotkeyProvenance = provenance` → `boot` | `:601` `	return rb.hotkeyProvenance`（`hotkeySourceName()`＝`:597`） | ✓ |
| 5 | `:157` `	hotkeyBridgeArmed bool` | `:368` `		rb.hotkeyBridgeArmed = true` → `boot` | 生产 **0 枚**（HEAD 非测试 grep 只命中 `:368`）；测试读点 `258_windows:283`/`:346`、`296_windows:163`/`:208`、`v1probe:62` ✓ 逐枚原文都是 `	if …rb.hotkeyBridgeArmed {` 形 | ✓ |
| 6 | `:161` `	hotkeyBridgeStop func()` | `:358` `		rb.hotkeyBridgeStop = func() {`、`:429` `			rb.hotkeyBridgeStop = nil` → 同 `boot` | `:427` `	if rb.hotkeyBridgeStop != nil {`、`:428` `		rb.hotkeyBridgeStop()` → `boot` | ✓ 两写同线程 ⇒ 不是竞争 |
| 7 | `:165` `	hotkeyBridge *ballHotkeyBridge258` | `:365` `		rb.hotkeyBridge = &ballHotkeyBridge258{check258: func() { bridge.Check() }}` → `boot` | 生产 0 枚；测试缝 `resident_hotkey_258_seam_windows.go:11` `	if rb == nil \|\| rb.hotkeyBridge == nil {`、`:14` `		rb.hotkeyBridge.check258()`（缝本体 `:10` `func (rb *residentBall) hotkeyBridgeCheck258() {` ✓） | ✓ |
| 8 | `:178` `	muteMux  sync.Mutex` | — | — | ✓ 护栏注释 `:173-177` 五行逐字对上，末句 `	// residentAudio handle before any gesture may read it.` ✓ ⇒ 它管发布序、**不覆盖 `:168` 那次挂载后的写**（票面现量第 13 条成立） |
| 9 | `:179` `	muteGate muteGestureFunc` | `:478` `	rb.muteGate = fn` → `boot`（`attachMuteGate`＝`:472`；锁在 `:476` `	rb.muteMux.Lock()` / `:477` `	defer rb.muteMux.Unlock()`；调用点 `resident_windows.go:336` `	rb.attachMuteGate(raudio.toggleMute)` ＋ 测试 `mute_290:93`/`:206`/`:235` ✓） | `:489` `	return rb.muteGate` → `ui-sta`（`currentMuteGate`＝`:483`，同锁 `:487`/`:488`；入口 `:516` `	fn := rb.currentMuteGate()`，两枚手势闭包＝`:321` `			OnMuteHotkey:    func() { rb.muteGesture("mute-hotkey") },` 与 `:325` `			OnTrayMute:      func() { rb.muteGesture("tray-mute") },` ✓） | ✓ 唯一带锁字段；⚠ 死腿点名的绕过锁直读 `mute_290:228` `	if rb.muteGate != nil {` 原文对上 ✓ |

## 3. ⓑ 护栏形状表逐枚复跑（尺＋目录＋含不含测试面都写全）

| 行 | 死腿断言 | HEAD 现量 | 结论 |
|---|---|---|---|
| atomic 家族 | 同 struct 唯一邻居＝`:110` `	levels atomic.Uint64`，写者 `capture`、读者 `boot` | 逐字对上（写 `:122`、读 `:359`） | ✓ **真邻居＝1 枚**成立 |
| 同包 atomic 先例 | `panel_resident_windows.go:136` `	stopRequested atomic.Bool` ＋ `:126-129` "四枚 `atomic.Int*`" | `:136` ✓；`:126`–`:130` HEAD 现量＝**5 枚** `atomic.Int64`（toggles/shows/hides/disposals/failedPost） | ⚠ 行区间少一行 ⇒ D-5 |
| 票面"全仓 5 处 `atomic.Bool`" | 尺 `git grep -n "atomic.Bool" HEAD -- internal cmd`＝5 命中，行号 42/136/85/132/69 | **5 命中逐枚对上**：`cmd/wisp/panel_host_windows_test.go:42` `	stopping atomic.Bool`、`cmd/wisp/panel_resident_windows.go:136`、`internal/audio/hotplug_test.go:85` `	closed atomic.Bool`、`internal/ball/ball_windows.go:132` `	closed     atomic.Bool`、`internal/ball/live_guard_windows_test.go:69` `	subclassActive atomic.Bool`。含 `_test.go`、含注释行 | ✓ 枚数口径（3 枚在测试面 ⇒ 生产先例 2 枚）✓ |
| 带锁 | `muteMux`/`muteGate` 那一对：`:476-478`、`:487-489` | 逐字对上 | ✓ 1 对 |
| 只在装配期写（单线程） | 12 枚（含 2 枚死字段） | 逐枚对文件：residentAudio 6 枚（`gate`/`mic`/`cancel`/`verdict`/`voiceEnabled`/`mutedAtBoot`）＋ residentBall 6 枚（`cancelHosted`/`verdict`/`hotkeyProvenance`/`hotkeyBridgeArmed`/`hotkeyBridge`/`hotkeyBridgeStop`）＝12 ✓ | ✓ 枚数成立 |
| 裸且跨线程 | 2 枚（`ra.started`、`rb.b`）⇒ `AC#1` 真邻居＝1 | HEAD 现量：`ra.started` 写者跨线程（`:168` ui-sta / `:323` boot）、读者 `:379` boot；`rb.b` 写 `:435` boot、读 `resident_audio_windows.go:136`/`:139` capture；除此以外的裸字段写点全在单线程 | ✓ 成立 ⇒ **`AC#0` 只有 `started` 成立、`verdict` 不成立；`AC#1` 真邻居＝1 枚**（死腿两大结论复跑通过） |

## 4. 线程词汇表锚复跑（ⓐⓑ 的线程标签全靠它）

`boot`：`resident_windows.go:110` `	defer func() {` + `:111` `		records := rt.Shutdown(false)` ✓、`:253` `	defer rb.stop()` ✓、`:260` `	defer ra.detachBall()` 行原文 ✓ **但接收者是 `*residentApproval`**（`resident_approval_windows.go:781` `func (ra *residentApproval) detachBall() {`）⇒ D-4。
`ui-sta`：`internal/ball/ball_windows.go:656` `	case wmHotkey:` ✓、`:661` `			b.fire(b.opts.Events.OnMuteHotkey)` ✓、`:669` `	case wmAppTray:` ✓、`:674` `			sel := showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)` ✓、`:678` `			case menuMute:` ✓、`:679` `				b.fire(b.opts.Events.OnTrayMute)` ✓、`:728-731` `func (b *Ball) fire(fn func()) { if fn != nil { fn() } }` ✓（内联同调，不起新线程）✓、`:734-736` `uiRun` 注释三行 ✓。⇒ **"两枚手势写者同线程 ⇒ `started` 只有 write(ui-sta)/read(boot) 一枚对"复跑通过**。
`capture`：`internal/audio/captureopt.go:110` `	c.OnLevel(level)` ✓、`resident_windows.go:311` ✓。
`test`：死腿"三枚读 `started` 的件里没有任何 `go func`/`WaitGroup`/`t.Parallel`"⛔ 本腿未复跑（那把尺在从未产出的 ⓒ 节里，且不影响 ⓓⓔ 的判定；落地腿若要引需自跑）。

## 5. 本腿自报命令面

- go 命令：**零枚**（`go env`、`go list` 都没跑；`go 1.27` 一类事实本件不引用）。
- 编译面／门禁／进程：**零枚**（⛔ `go build`/`vet`/`test`/`run`/`gofumpt`/`d22scan`，⛔ 未启动 `build/**` 任何 exe，未动 `build/**`、`frontend/dist/**`）。
- 取 blob 一律 `git show HEAD:<path>` 管道到 `cat -n`；临时落盘在仓外 `$TEMP/295a1b/`（3 枚文本快照），仓内 ⛔ 无 `.go` 拷贝。
