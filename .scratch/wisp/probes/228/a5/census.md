# 228-a5 — 只读普查：「审批卡片要真能被答复」这一格缺的**入口**名册

射程＝从"有人能按"到"Go 真的收到"的每一环，逐枚点名现在处于**存在／写了没接／一块没写**哪一态。
零产码、零构建、零 Go 命令。本件不落任何裁定（契约面＝编排者裁，C17 面＝人工批准）。

---

## §0 锚点与口径

- 起手 `date` 逐字读数：`2026-10-03 10:01:23 +0800`
- 起手 `git log -1 --format=%H` 逐字读数：`ccd9543b3d86f3f627d331b3c399cb3d4ee99eec`（branch `dev`）
- ⛔ **本程零 Go 命令自证**（同机三枚写腿在飞，任何编译都会互洗它的读数）：
  本腿全程只用 `Read`／`Grep`／`Glob` 与 `grep`／`sed`／`awk`／`ls`／`wc`／`git log`／`git diff --quiet`。
  未跑的字面清单＝`go build`／`go vet`／`go test`／`go run`／`go list`／`build.ps1`／`go mod tidy`（一枚都没跑）。
- ⛔ 零读、零转述 `frontend/**` 与 `design/**`（两层禁令）。
- 路径口径：所有引用报到**包／目录级 + 文件:行号**。
- "零 X" 的入账规矩：每枚负向结论自带**那把尺的字面命令 + 一处正控**（已知命中的行）。

---

## §1 「允许」入口全名册（三态）

链路的六个环（先立环，再逐枚入口点名落在哪一环）：
**手 → 原生事件 → `cmd/wisp` 回调 → `approval` 路由 → 队列落答 → 审计／球态**。

### 六环各是谁（承重行的现读位置）

| 环 | 位置 |
|---|---|
| 原生事件 | `internal/ball/ball_windows.go:656-686`（`wmHotkey` 与 `wmAppTray` 分流）＋`internal/ball/tray_windows.go:72-107`（`showMenu` 建菜单＋取选择） |
| `cmd/wisp` 回调 | `internal/ball/ball_windows.go:50-59`（`Events` 十枚 `func()` 键）→ `cmd/wisp/resident_ball_windows.go:173-184`（常驻宿主填的那张字面量） |
| `approval` 路由 | `internal/agent/approval/replies.go:316`（`Replies.Allow`）→ `:328`（`g.DecideFromNative`）→ `internal/agent/approval/gate.go:718`→`:723`（`g.q.allow`） |
| 队列落答 | `internal/agent/approval/queue.go:396`（`ANSWER-ALLOW … decision=allow`）／`:391`（`allow-session` 那一支） |
| 审计 | 门侧 `queue.go:396` 经 `Options.Logf`：常驻＝`cmd/wisp/resident_approval_windows.go:129`，CLI 腿＝`cmd/wisp/run.go:931`；宿主侧另加一行＝`cmd/wisp/approval_reply.go:232`→`:414` |
| 球态 | `cmd/wisp/resident_approval_windows.go:465-474`（`Update`/`EventDismissed`）→`:483-493`（`settleOrb`：还 Esc ＋ 回 `Sleeping`） |

### 入口逐枚三态

| # | 入口 | 三态 | "接"的凭据（谁调谁，行号）／缺的那一环 |
|---|---|---|---|
| 1 | **控制台 stdin @ `wisp run`**（一发即退那条腿） | **存在** | `cmd/wisp/run.go:802 attachReplyListener(s.reply, s.replyVeto)` → `cmd/wisp/approval_reply.go:451` Spawn `approval-waiter` → `:530 runReplyLoop` → `:562 handle`（动词名册 `:563-583`，`yes`＝`:564`）→ `:213 replySurface.allow` → `replies.go:316 Replies.Allow` → `replies.go:328` → `gate.go:718/723` → `queue.go:396` |
| 2 | **控制台 stdin @ 常驻腿**（owner 双击那枚 `wisp.exe`） | **存在** | `cmd/wisp/resident_windows.go:206 startResidentTaskSource(rt, ra)` → `cmd/wisp/resident_task_source_windows.go:218 interactiveStdin()` → `:288 run.newReplySurface(...)` → `:318 observe.Default.Spawn(taskSourceConsoleGoroutine…)` → `:345 runConsoleLoop` → `:381 src.surface.handle(verb, corr, arg)` → 同 1 的后半条漏斗（同一张动词表，不是副本） |
| 3 | **全局 Esc 热键 @ 常驻腿** | **存在，但只有否决那一支** | `cmd/wisp/resident_windows.go:163 startResidentBall(…, ra.vetoByEsc, …)` → `cmd/wisp/resident_ball_windows.go:177 OnCancelHotkey: func(){ recordCancelHotkey(onCancelEsc) }` → `:279` → `cmd/wisp/resident_approval_windows.go:171 vetoByEsc` → `:179 ra.cards.Veto` → `replies.go:463 g.Veto` → `gate.go:415`。**「允许」在这一支＝一块没写**：`Replies.Allow` 全仓没有任何 GUI 调用者（尺 R8/R10） |
| 4 | **托盘菜单** | **一块没写** | 菜单在册只有四枚：`internal/ball/tray_windows.go:20-24`（id 1..4）＋`:86/:88/:89/:91`（打开面板／静音／暂停唤醒／退出）；分流 `internal/ball/ball_windows.go:676-683`；`Events` 十枚键里**没有一枚通向允许**（`ball_windows.go:50-59`）。⇒ 五环里**前三环全缺**（没有手能按的那一枚、没有 id、没有 `Events` 键），后三环（路由／落答／审计）现成可接 |
| 5 | **单击球（`ChannelBall`）** | **写了没接**（三处都在，中间那一环是空的） | 「写」的部分：`internal/ball/ball_windows.go:632 b.fire(b.opts.Events.OnClickBall)` 有发；`cmd/wisp/resident_ball_windows.go:174 OnClickBall: func(){ recordBallGesture("click") }` 有听。**「没接」的部分**：那枚回调只记账不答复（`ballGestureWhy` 逐字见 `:257-259`），且 `ChannelBall` **从未被 `SetLoaded`**（尺 R50：全仓非测试 `SetLoaded` 调用点＝`approval_reply.go:512`、`resident_approval_windows.go:157`/`:351`，命中的通道只有 `ChannelEsc`／传进来的 `vetoChannel`） |
| 6 | **面板 WebView** | **契约上就不该有这一支**（结构性无门，不是漏接） | `internal/agent/approval/ui.go:167-171` `PanelAPI` 只有 `Reject/Head/View`，**没有 `Allow` 方法**（`:164-166` 逐字「a panel-sourced allow is structurally rejected … it is a method that does not exist」）；C17 入向名册四枚 `internal/panel/bridge.go:42-45`，**无 `approval.decide`**；`cmd/wisp/panel_inbound.go:271-277` 的 `Workspace/Attachment/Message` 全为 `nil` |
| 7 | **KWS 否决词（`ChannelKWS`）** | **一块没写** | 尺 R50：`ChannelKWS` 在非测试码里只出现在定义／名册／文案三处（`internal/agent/approval/approval.go:57/:61/:92/:105`），**零枚 `SetLoaded`、零枚宿主** |
| 8 | **一次性 CLI 子命令（例：`wisp approve <corr>`）** | **一块没写** | 子命令名册＝`cmd/wisp/main.go:89-124` 的十枚（`run/providers/doctor/secret/models/slo/panel-assets/panel-inbound/version/help`），**没有一枚是答复用的**；答复只活在 1／2 那两条 stdin 环里 |

### 特别判死：`Gate.Replay` 有没有产码调用者

**复认腿报＝零枚。成立。**

- 尺 R1 字面命令：`grep -rn "\.Replay(" --include=*.go internal cmd tools scripts`
  全量读数＝**1 行**：`internal/agent/approval/queue_test.go:283`（测试）。
- 第二把尺 R52 防止"绕过方法名"这一形：`grep -rn "\.replay(\|replay(" --include=*.go internal cmd`（剥 `_test.go`）
  读数＝**2 行**：`internal/agent/approval/gate.go:753`（`Gate.Replay` 自己去调 `q.replay`）＋`internal/agent/approval/queue.go:552`（定义）。
  ⇒ 队侧那半（`queue.go:279-291` 的 `keepForReplay`/`maxRepl`、`:113-114 DefaultReplayHistory = 8`）**在生产里一直往里存**，**取的那一枚没人调**。
- 定义与语义：`internal/agent/approval/gate.go:745-748` 逐字「Replay re-displays a refused or expired L2 request under a fresh correlation id (C18 一键重放) … it is NOT an answer」。
  ⇒ 对本格的影响：**被拒／超时的卡今天没有任何一条路能被重新问一次**；"允许"入口即使补上，也只能答**此刻还挂着**的那一枚。
- 正控（证明这把尺不瞎）：同一把 `\.Method\(` 形尺在**同包**能命中生产调用者——
  `grep -rn "g\.Native()" --include=*.go internal cmd` 命中 `internal/agent/approval/replies.go:363`（产码），
  `grep -rn "q\.allow(\|g\.q\.allow(" --include=*.go internal/agent/approval`（R54）命中 `gate.go:625` 与 `gate.go:723`（产码）。
  尺能命中别的方法的生产调用者 ⇒ `Replay` 那个零是**真的零**，不是尺的形状不对。

---

## §2 托盘加一枚「允许一次」的最小改动面

### 2.1 逐处列（改动面＝五处产码 + 一处名册）

| 处 | 文件:行 | 要动什么 | 现读凭据 |
|---|---|---|---|
| ① 菜单项在哪建 | `internal/ball/tray_windows.go:20-24`、`:86-91` | 加一枚命令 id（在册是 `1..4`，下一枚＝`5`）＋一行 `appendItem(id, 标签, false)` | `appendItem` 闭包在 `:79-85`，签名 `(id uintptr, label string, checked bool)`；`showMenu` 只收两枚布尔（`:72`） |
| ② 事件怎么发出来 | `internal/ball/ball_windows.go:669-685` | `switch sel` 加一条 `case` → `b.fire(...)` | 分流今天只有四条：`:676/:678/:680/:682` |
| ③ 回调键 | `internal/ball/ball_windows.go:50-59` | `Events` 加第 11 枚 `func()` 字段 | `fire` 的形＝`func(fn func())`（`:728-732`），只判 nil |
| ④ 常驻宿主填哪 | `cmd/wisp/resident_ball_windows.go:173-184` | 那张字面量加一枚键，值走"注入函数值"形（同 `:177` 的 `onCancelEsc` 形） | 宿主对审批一无所知的分层写在 `:23`、`:56-61`、`:195-203` |
| ⑤ `approval` 侧接哪枚方法 | **现成，不用新增**：`internal/agent/approval/replies.go:316 (*Replies).Allow(ctx, corr) error` | 关联号取 `AwaitingHuman()`（`replies.go:252`，L2 优先的注释在 `:242-251`） | 它自己走 `:328 DecideFromNative` → `gate.go:718/723` → `queue.go:396`；无卡返回 `ErrNoTrackedCard`（`:57`） |
| ⑥ 会被打红的名册 | `cmd/wisp/resident_ball_228_test.go:50-53`（列表）＋`:319-322`（`len(set) > len(list)` 那一支） | 见 §5，落地时**只许扩列表** | 该钉用 AST 扫 `Events` 字面量的 key，两向都判 |

**常驻腿握着的是哪枚对象**：`*residentApproval`——`cmd/wisp/resident_windows.go:123` 由 `newResidentApproval()` 构造，
字段在 `cmd/wisp/resident_approval_windows.go:81-94`：`gate *approval.Gate`（:82）、`cards *approval.Replies`（:83）、`ui *ballCardUI`（:84）。
`cards` 与 `gate` 在 `:115-120` 由 `HostBinding{Gate, VetoChannel: approval.ChannelEsc, NativeSource, PanelSource}` 绑成一对。
⇒ 托盘那一发要够到的就是这枚 `ra.cards`，而 `ra` 今天已经通过 `:163` 的函数值注入够得着球宿主——**同一枚注入形状，不必新开装配点**。

### 2.2 「要不要新增方法名」——只报事实，⛔ 不替编排者裁

三层逐个报名字与它在不在 C17 那张表上：

1. `internal/ball` 层：加的是**一枚 Go 结构体字段**（`Events` 的第 11 枚）与**一枚命令 id 常量**。C17 那张方法表在
   `internal/panel/bridge.go:42-45`（在册只有 4 枚 `panel.*`），**托盘这一支不碰这四个名字中的任何一个**。
2. `cmd/wisp` 层：加的是**一枚函数值**（`startResidentBall` 的入参/回调），不是方法名，也不是线协议名——
   `cmd/wisp/approval_reply.go:561` 逐字「not a wire protocol, and nothing here is a C17 method name」，那枚 `switch verb`（`:563-583`）就是这条判断的现读形状：控制台动词同样不算 C17。
3. `approval` 层：`Allow` / `AllowSession` / `Reject` / `Veto` **四枚名字全在册**
   （`internal/agent/approval/ui.go:146`、`:158`、`:161`；`replies.go:316`、`:351`、`:372`、`:455` 一带），
   托盘**不需要任何一枚新名字**。

⇒ 事实结论：**托盘这一支的三层里没有一处需要新增 `panel.*` 入向方法名**。
⚠ 但有一枚**别的在册闭集**会被"顺手扩"打到：`internal/agent/approval/approval.go:48-50`
逐字「The set is closed: a fifth value means a caller invented a selector, which is exactly the M-7/C-3 failure shape」，
名册在 `:61 var allChannels = []Channel{ChannelBall, ChannelEsc, ChannelPanel, ChannelKWS}`。
⇒ 若实现者想让审计行"带正确 `Channel`"，就会撞上这四枚的闭集；归 **`A498` F8 已裁的那一格（裁甲＝`Source` 标签是标签、不是权威）**，本腿只标出**碰哪一行**（`approval.go:48-50`/`:61`）。

### 2.3 与票 245「稳态不绑裸 Esc、只在确认那两三秒借」的关系——只写碰哪一行／不碰

**不碰**（现读为证，四处）：
- 借／还那对机制的两个作者点：`internal/ball/hotkey_windows.go` 的 `takeEsc`／`releaseEsc` 一支，与
  `internal/ball/ball_windows.go:250` 的 `registerAll(b.hwnd, b.opts.Hotkeys)`（开窗即绑稳态三枚）——
  **托盘菜单项不是键绑定**：`showMenu` 全程只调 `CreatePopupMenu`/`AppendMenuW`/`TrackPopupMenu`（`tray_windows.go:73-102`），零处 `RegisterHotKey`。
- 稳态枚数那枚钉：`internal/ball/live_windows_test.go` 一族锚的是**注册集枚数**，托盘项不进注册集 ⇒ 不会因这枚改动变红。
- 常驻腿那行 `hotkeys live %d/4` 文案（`cmd/wisp/resident_ball_windows.go:207`）：分子是 `len(rep.Live())`，与菜单枚数无关。
- **借期条件本身**：`cmd/wisp/resident_approval_windows.go:440-442` 逐字 `if p.Level == "L1" { b.TakeEscForCancel() }`
  ⇒ **L2 卡根本不借 Esc**，而按 `A498` F6 那一裁，托盘「允许一次」永远作用于队头那枚 **L2** 卡（`AwaitingHuman` 的 L2 优先写在 `replies.go:242-251`）
  ⇒ **两件事作用域不重叠，不冲突。**

**碰**（两处，都是"必须一起改口／一起扩"的账，不是矛盾）：
1. `cmd/wisp/resident_ball_windows.go:173-184` 那张 `Events` 字面量枚数 +1 ⇒ `cmd/wisp/resident_ball_228_test.go:50-53`
   列表要同步扩（`A498` F7 已裁："只许扩列表、不许扩断言维度"）。
2. 托盘一次选择会**把一张 L2 卡答掉**，答掉之后走 `Update(EventDismissed)` → `settleOrb`（`cmd/wisp/resident_approval_windows.go:465-474`、`:483-493`）。
   那枚函数末尾正是"还 Esc ＋回 `Sleeping`"那一支。
   ⇒ **碰的是它的调用者枚数**（今天只有 CLI 答复与退出序列两条会走到这里，`ball_windows`/托盘是第三条），
   **不改它的判据**：`settleOrb` 首行 `if _, awaiting := u.ra.cards.AwaitingHuman(); awaiting { return }` 已经保证"两张卡时第一张落不惊动第二张"。
3. ⚠ 一条**回调线程**的在册坑（票 228 `A493` 已经写死同一形状）：`b.fire(...)` 是**内联在 `ui-sta` 线程上跑**
   （`internal/ball/ball_windows.go:728-732`；同文件 `:734-736` 注释逐字「every Events callback」跑在 UI 线程上），
   而 `Ball.Close()` 那形是 `sta.PostTask + <-done`（`:904` 一带）⇒ **托盘 `Allow` 回调里不许同步等任何球侧收口**；
   合法形状＝今天 `vetoByEsc` 那一形（`cmd/wisp/resident_approval_windows.go:163-171` 逐字「It runs ON the ui-sta thread … so it waits on nothing」）。

---

## §3 「没卡时按允许」该说什么

**结论先报**：这一格**不需要造词**——现成的诚实形状有**五族**，而且**同一枚错误码已经被钉过**。
下面逐枚给位置，⛔ 不推荐、不裁。

### (a) 语义层：`Allow` 无卡时返回的那一枚 sentinel（今天就已经是这个返回）

- `internal/agent/approval/replies.go:57` 逐字
  `ErrNoTrackedCard = errors.New("approval: 本机账上没有这张卡（从未显示、已答复，或已经离开屏幕）")`
  ——**三支含义它自己列全了**（未显示／已答复／已离屏）。
- 送出点：`replies.go:317-319`（`Look` 不中即返回它）；姊妹两支在同一处开关里：
  `:60` `ErrRouteHasNoAllow`「这一路线只有否决、没有允许」（L1 窗口那支，`:321-322`）、
  `:64` `ErrNoGateAttached`（半接线那支，`:325-326`）。
- **已被看守**（＝不需要新语义的证据）：`internal/agent/approval/replies_201_test.go:93-94`
  断 `Allow(ctx,"never-displayed")` 必 `ErrNoTrackedCard`；`:99-102` 断 L1 窗口上必 `ErrRouteHasNoAllow`；
  `:104-107` 反向断"被拒的 allow 不许顺手把卡片 retire 掉"（`Look("w-1")` 仍在）。

### (b) 宿主成句：控制台那一侧已经把同一枚码翻成给人看的那一句

- `cmd/wisp/approval_reply.go:217-219`：`ErrNoTrackedCard` →
  「没有找到待答复的卡片 %q：它可能已经结束，**从未显示的卡片这里也没有令牌可花**」，
  并且**同时**记一行 `s.record("REFUSED", corr, "", "native/allow", "本机没有这张卡的记录")`（`:218`）。
- 姊妹支：`:220-223`（L1 只有否决）、`:226-229`（令牌失效＝要重新显示才能被允许）。
  ⇒ 这一族的形状＝**"人读的因果句" + "机读的 route 标签"两句一起出**，不是只印一个错误码。

### (c) GUI 同线程同形状的现成实话（离托盘最近的一枚）

- `cmd/wisp/resident_approval_windows.go:171-178` `vetoByEsc()`：`AwaitingHuman()` 不中即返回
  「按下的取消键没有可否决的确认项：本进程此刻没有卡片在等人」，注释逐字
  「Said out loud rather than left silent」。
- **它被钉着**：`cmd/wisp/resident_approval_246_windows_test.go:341-346`
  `TestAC246VetoSentenceWithNoCard` 断那句里必须出现「没有可否决的确认项」，
  并反向断空账上 `AwaitingHuman()`/`WaitingState()` 都不许说"有人在等"。
  ⇒ **"按下那一枚、此刻没有对象"这一形状在本仓已有一条会响的钉**，托盘那枚只是它的第二条。

### (d) **同一枚托盘文件**里的"诚实落空"先例（两枚）

- `cmd/wisp/resident_ball_windows.go:301-312` `recordTrayExit()`：`const why` 一整句因果
  ＋`slog.Warn("tray exit requested", "outcome", "ignored", "why", why)` ＋`fmt.Printf`，
  注释逐字「it is **reported instead of assumed**」。
- `cmd/wisp/resident_ball_windows.go:261-269` `recordBallGesture()`：`slog.Warn("ball gesture arrived with no executor", …)`
  ＋控制台一行，共用的 `why` 常量在 `:249-259`。
  ⇒ **出声的通道只有一条能选**：`b.fire(fn func())`（`internal/ball/ball_windows.go:728-732`）没有回程值，
  所以落空的话**只能由 `cmd/wisp` 侧说**（这与 `A498` F4 的"⛔ 不许让 `showMenu` 的返回值带出错误"是同一枚事实）。

### (e) B1 那一族（通道可用性），⚠ 与 (a)-(d) **不是一族**

- `internal/agent/approval/approval.go:80` 逐字「An unloaded channel is NEVER rendered as available (B1)」，
  `:95-119 unavailableText()`（四枚通道各一句不同文案），`:121-124 ErrChannelUnavailable`
  逐字「a silently ignored cancel attempt is how a fake channel becomes a user-visible promise」。
- 送出点：`internal/agent/approval/gate.go:605` `Channels: g.channels.Statuses()`（随卡片一起给宿主）。
- ⚠ **族的边界**（这一行只报名字在哪，不裁该用哪族）：
  (a)–(d) 说的是**"此刻没有对象"**；(e) 说的是**"这条路不能开"**。
  后者的"被拒"句在 `internal/agent/approval/ui.go:124-136`（`ErrUnknownCorrelation`/`ErrNotPending`/`ErrBadGrant`/`ErrPanelAllow`/`ErrNotAdmitted`）
  与 `internal/agent/approval/gate.go:740`（面板来源的「允许」被服务端 API 直接拒），
  以及 `cmd/wisp/resident_approval_windows.go:292-293`（退出时"未获批准，按拒绝处理"）。
  ⇒ 若落地腿要把"没卡"翻成上面任何一句，**改的就是这一族的语义**。

### 还有一族"入口不存在时说什么"的在册形状（同一条腿、同一批常量）

- `cmd/wisp/resident_task_source_windows.go:84-111`：五枚 claim 常量（`taskEntryDisabledClaim`/`taskEntryConsoleClaim`/
  `taskEntryInjectedClaim`/`taskEntryRefusedClaim`/`pipelineAbsentClaim`）＋四枚 posture 常量，
  每个 return 分支一枚，打印在 `:220-246`。钉在 `cmd/wisp/resident_task_source_246_windows_test.go`。
  ⇒ 本仓对"这一条腿今天能不能被用"已经**按分支各给一句**，不是一句兜底。

---

## §4 审计链与零仪器段

### 4.1 一次 allow 从入口到落审计行经过的处数＝**六处**（现走通的那条＝控制台）

| 序 | 处 | 行号 |
|---|---|---|
| 1 | 入口分发 → 答复面 | `cmd/wisp/approval_reply.go:564`（`yes`）→ `:565`→`:213 (*replySurface).allow` → `:215 s.live.h.Allow(s.ctx, corr)` |
| 2 | 账本 → 路由 | `internal/agent/approval/replies.go:316-335`：`Look`→`Grant` 判→`:328 g.DecideFromNative(Request{CorrelationID, Allow:true, Grant, Source:native})` |
| 3 | 门 → 队列 | `internal/agent/approval/gate.go:718` `DecideFromNative` → `:723 return g.q.allow(r.CorrelationID, r.Grant)` |
| 4 | **第一行审计**（门自己的） | `internal/agent/approval/queue.go:396` `approval: ANSWER-ALLOW corr=… tool=… route=native decision=allow`（会话那一支＝`:391` 的 `decision=allow-session`）；写入器＝`Options.Logf`，两枚宿主各一枚：控制台＝`cmd/wisp/run.go:617`→`run.go:931-935`（stderr `[audit] ` ＋ `spec.sink.logger()` 落盘），常驻＝`cmd/wisp/resident_approval_windows.go:112`→`:129-133`（`slog.Info("audit: "+line)` ＋ **stdout** `fmt.Printf("wisp: [audit] %s")`） |
| 5 | **第二行审计**（宿主自己的） | `cmd/wisp/approval_reply.go:232 s.record("ANSWERED", corr, card.Tool, "native/allow", "原生侧允许")` → `:414-419`（`s.audit` 就是第 4 行那枚 `rt.auditf`） |
| 6 | 账本收口＋球态 | `replies.go:333 r.Forget(corr)`；球侧经 `cmd/wisp/resident_approval_windows.go:465-474 Update(EventDismissed)` → `:483-493 settleOrb()`（还 Esc ＋ 回 `Sleeping`），**这一段零枚审计行**（只有 `slog.Info`/`slog.Warn` 的是别的分支） |

⚠ 一处**读数**（不是判语）：常驻腿里一次 allow 的两行审计**由两枚不同的 writer 出声**——
门那行走 `ra.residentAuditf`（stdout ＋ 当时那枚进程默认 slog），宿主那行走 `rt.auditf`（stderr ＋ 钉死在它自己那枚 `spec.sink`）。
而常驻腿里 `installLogSink` **被调用两次**：`cmd/wisp/resident_windows.go:63` 与嵌套装配 `cmd/wisp/run.go:222`；
`installLogSink` 无条件 `slog.SetDefault(...)`（`cmd/wisp/logsink.go:152`）⇒ 后一次会**替换**进程默认。
⇒ 事实＝两行落在哪个文件／哪个流，取决于第 4 行与第 5 行各自的绑定，今天**没有一枚用例读过这两行的相对位置**。

### 4.2 哪一段有会响的钉（尺＝钉，不是注释）

| 段 | 有钉？ | 具名 |
|---|---|---|
| 控制台 `yes` → 两行审计 | **有** | `cmd/wisp/approval_reply_201_test.go:252-253`（**同时**断门那行 `ANSWER-ALLOW …decision=allow` 与宿主那行 `REPLY ANSWERED route="native/allow"`）；`cmd/wisp/approval_seam_201_test.go:119` |
| 无卡／L1／半接线三种落空各返回什么 | **有** | `internal/agent/approval/replies_201_test.go:93-94`/`:101-102`/`:81`/`:106` |
| grant 花掉后第二次花必失败 | **有** | `cmd/wisp/subagent_selfapproval_197_test.go:450`、`:466-472`；`internal/agent/approval/queue_test.go:161-166` |
| 面板路线的 allow 必被判红 | **有** | `internal/agent/approval/queue_test.go:89`、`cmd/wisp/approval_reply_201_test.go:340`、`cmd/wisp/approval_seam_201_test.go:137` |
| 常驻腿**否决／退出拒绝**那两族的审计行 | **有，但只在 `winlive` 档** | `cmd/wisp/resident_approval_live_246_windows_test.go:164`（`ANSWER-VETO`）、`:265`（`ANSWER-REJECT`）；`cmd/wisp/resident_task_source_live_246_windows_test.go:164`（出货进程日志里找 `ANSWER-VETO`）。文件头 tag＝`//go:build windows && winlive`（现读 `:1`） |

### 4.3 零仪器的三段（只报"哪一段零枚钉"，不写能不能做到）

1. **常驻腿的 `ANSWER-ALLOW` 那一行**：尺字面
   `grep -rn "ANSWER-ALLOW" cmd/wisp/resident_*_test.go` ⇒ **零命中**（正控＝同一把尺在同名文件里能命中 `ANSWER-VETO`／`ANSWER-REJECT`，见 `resident_approval_live_246_windows_test.go:164`/`:265`）。
   ⇒ 今天**没有任何用例读过常驻那条腿的一发"允许"落在审计里长什么样**，因为那条腿上今天也没有允许入口（§1 第 4/5/7 枚）。
2. **"GUI 那一次点击 → 路由"整段**：尺字面
   `grep -rn "打开面板\|暂停唤醒\|静音\|menuExit\|showMenu\|OnTray" --include=*_test.go internal cmd` ⇒ **唯一命中是 `cmd/wisp/resident_ball_228_test.go:52` 那枚字面名字册**（正控＝同尺去掉 `--include=*_test.go` 在产码里命中 12 行，见 R39）。
   ⇒ 菜单枚数／命令 id／标签文案／`showMenu` 入参**四者今天零枚钉**；能被"多一枚键"打红的只有那枚名字册（§5）。
3. **`Gate.Replay` 的审计行 `approval: replay …`（`internal/agent/approval/gate.go:758`）**：
   该方法产码零枚调用者（§1 判死），故这一行**今天不可能出现在任何盘上**；
   唯一覆盖 `Replay` 语义的是**进程内**用例 `internal/agent/approval/queue_test.go:269 TestReplayRedisplaysUnderAFreshGrant`，它不读审计。
   ⚠ 另：`winlive` 那一档在 CI 里**零命中**（尺字面 `grep -n "winlive" .github/workflows/ci.yml` ⇒ 零；正控＝同一条命令换成 `grep -n "winlive\|go test"` 时 `go test` 侧命中 10 行，证明这把尺读得到那个文件）。
   ⇒ 上表第 4 行里"有钉但只在 `winlive`"那几枚，**CI 默认档零读数**。

---

## §5 既有钉名册与射程

（待填：`internal/agent/approval`／`cmd/wisp` 里 allow／replay／tray 相关用例逐枚点名 + 哪几枚盖得到"新增托盘入口"）

---

## §6 量不到的格子

（待填：量不到 + 为什么；⛔ 不用推测填空）

---

## §7 我推翻前人哪几句

（待填：含 `259-a1` 两条线索与编排者 `A4xx` 若干句）
