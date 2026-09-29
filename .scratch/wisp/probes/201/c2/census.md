# 201-c2 —— 宿主原生答复入口（球／托盘）接进 `wisp run` 的落点普查

- 本件作者＝只读普查腿 `201-c2`（编排者于 09-29 14:05 派出，打的是票 201 剩余四格 AC#1／AC#2／AC#3a／AC#6b 共同的那枚前置）。
- ⚠ **代落盘说明（编排者，09-29 14:2x）**：该腿的工具集里**没有 Write／Edit**（`Explore` 只读），且系统硬禁建文件／重定向，所以它交回的是**回执正文**而不是这份文件。全文由我逐段落盘，**它一个字没被我改写**；下面第 0 节是我自己复跑的尺，用来交代"哪些句子我信了、哪些我另有一读"。
- 它的锚点：起手 `29daf745`、交件时 `721e3864`，两发之间 `git diff --stat` 零枚 `.go` ⇒ 本件所有行号在两发之间不漂。**⛔ 引用本件行号前请现跑复核**：`internal/config/**` 与 `cmd/wisp/run.go:341` 一带会随票 223／224 落地而漂。

## 0. 编排者自己复跑的尺（09-29 14:1x–14:2x，逐条现跑）

| 该腿的承重主张 | 我自己跑的命令 | 我的读数 | 裁定 |
|---|---|---|---|
| `internal/ball` 生产零引用 | `grep -rln 'wisp/internal/ball' --include=*.go .`（剥 `.scratch`） | **只有 `./cmd/balldebug/main.go` 一枚** | ✅ 复认 |
| `go.mod` 里没有 webview 类依赖 | `grep -icE 'webview\|systray\|walk\|getlantern' go.mod` | **0** | ✅ 复认 |
| `wisp run` 一发即退 | `grep -n 'os.Exit(cmdRun' cmd/wisp/main.go`／`grep -n 'return rt.execute(task)' cmd/wisp/run.go` | `main.go:85`＋`run.go:220` | ✅ 复认 |
| `replyVeto` 生产赋值点＝0 | `grep -rn 'replyVeto' --include=*.go cmd/ internal/ \| grep -v _test` | 只 4 命中：`:140`/`:440` 注释、`:146` 声明、`:601` 读 ⇒ **赋值 0** | ✅ 复认（这条正是票 201 AC#2 改写形不成立的原因） |
| `Answer` 词表只有 4 枚、**没有会话档** | `grep -n 'Answer[A-Za-z]* Answer = ' internal/tools/gate.go` | `:78 allow`／`:80 reject`／`:83 veto`／`:88 timeout` | ✅ 复认 ⇒ "本次会话内允许"在 `Answer` 里确实没有对应值，那一格归票 224 |
| `Gate` 的 `windows` 字段没有枚举口 | `grep -n 'g\.windows' internal/agent/approval/gate.go` | **只有 4 处**：`:321` 查重／`:324` 写／`:330` 删／`:393` 按 corr 查 ⇒ **无"列出此刻在等人的 L1"那一枚** | ✅ 复认 |
| `DefaultChannels()` 生产零调用者，且装配处**主动**传空 | `grep -rn 'DefaultChannels()'`＋`grep -rn 'NewChannels('` | 定义 `approval.go:162`；生产命中只有 `gate.go:108`（`Channels==nil` 才走它）；而 **`run.go:446 Channels: approval.NewChannels()`＝空参**，`approval.go:163` 逐字 `NewChannels(ChannelBall, ChannelEsc)` | ✅ 复认，且**这条比它说得更重**：现成半成品就在场，生产是**显式放弃**了那个兜底，不是"没建" |
| `AddAllowedDir`／`SetAllowedDirs` 的生产调用者 | `grep -rn 'AddAllowedDir(\|SetAllowedDirs(' --include=*.go internal/ cmd/ \| grep -v _test \| grep -v '^internal/config/allowdirs.go'` | `AddAllowedDir`＝**1 枚**（`cmd/wisp/approval_always.go:134`）；`SetAllowedDirs`＝**0 枚** | ✅ 复认（"可撤销"那一半今天真的没人调） |
| 三枚答复函数的生产调用者 | 我用自己的尺跑同一把 grep | **6 命中**：`replies.go:325/372/377/401/427`（＝`Gate` 那三枚）**＋ `cmd/wisp/approval_reply.go:310`** | ⚠ **我与它差一行的那一枚不是矛盾**：多出的 `:310` 是 `s.live.h.Veto(corr)`，调的是 **`(*Replies).Veto`**（`replies.go:419`）这一层，不是 `Gate.Veto`。⇒ 准确说法＝**`Gate` 那三枚的生产调用者 5 枚、全在 `replies.go`；`Replies` 那一层的调用者在 `approval_reply.go`（`:210/:249/:251/:278/:310`）** |
| 托盘现有四枚菜单项 | `sed -n '17,24p;86,91p' internal/ball/tray_windows.go` | id `menuOpenPanel=1`/`menuMute=2`/`menuPauseWake=3`/`menuExit=4`；文案逐字「打开面板」「静音」「暂停唤醒」「退出」＋两条分隔条 | ✅ 复认 |
| 装配处确认钩子仍是 `nil` | `grep -n 'config.NewManager' cmd/wisp/run.go` | **`run.go:341 config.NewManager(cfgPath, nil)`** | ✅ 复认 ⇒ 票 223 AC#6 那枚"第二枚断口"在这儿，行号以我这一读为准 |

**另两格它自己标了"未行级证"的，我也没补上**（`l2_grant_boundary_test.go` 的 AST 扫描根目录、import 环）——**别把这两格当已证**。

---

# ① 球今天是怎么被拉起来的

**链（每一环现读到，〔已证〕）**

`cmd/balldebug/main.go:87 func main()`
→ `:122 ball.EnablePrototypeVisuals`（定义 `internal/ball/statevisual.go:106`）／`:124 ball.SetLook`（`internal/ball/tokens.go:332`）
→ `:188 statemachine.New` → `:189 ball.New(ball.Options{…})`（定义 `internal/ball/ball_windows.go:143`）
→ `ball_windows.go:166 newSTAThread`（`internal/ball/sta_windows.go:38`）→ `:170 opts.Registry.Spawn("ui-sta", "ball", nil, …)` → `sta_windows.go:47 (*staThread).start`
→ `sta_windows.go:52 runtime.LockOSThread` → `:62 pCoInitializeEx(coinitApartmentThreaded)` → `:65 create(s)` = `ball_windows.go:182 createOnSTA` → `:77-93` 消息泵（`GetMessageW`/`TranslateMessage`/`DispatchMessageW`）
→ `createOnSTA` 内部顺序：`:183 ensureFactories`（`d2d_windows.go:181`）→ `:187 registerBallClass`（`ball_windows.go:954`）→ `:208 pCreateWindowExW`（exStyle `:196` = `wsExLayered|wsExTopmost|wsExToolWindow|wsExNoActivate`，style `:212` = `wsPopup`）→ `:229-234 resolveInitial/moveWindow/getDpiForWindow` → `:237 newRenderer`（`renderer_windows.go:85`）→ `:243 addTrayIcon`（`tray_windows.go:30`）→ `:250 registerAll`（`hotkey_windows.go:377`）→ `:254 pShowWindow(swShownoactivate)` → `:258 applyStateLocked`（`ball_windows.go:337`）
→ 回调出窗：`ball_windows.go:535 ballWndProc`，tray 分支 `:664 case wmAppTray`。

**它引了 `internal/ball` 哪几枚导出符号（现跑 `grep -o "ball\.[A-Za-z_]*" cmd/balldebug/*.go`，23 枚，逐枚）**
`Ball` `New` `Options` `Events` `EventKind` `EvClickBall` `EvSummonHotkey` `EvMuteHotkey` `EvCancelHotkey` `EvPanelHotkey` `Edge` `EdgeByName` `EdgeLeft` `EdgeRight` `EdgeTop` `EdgeBottom` `EnablePrototypeVisuals` `PrototypeVisualsEnabled` `SetLook` `LookNames` `HotkeyConfig` `ApplyHotkeyDefaults` `NewHotkeyReloader`。

**`internal/ball` 的生产引用面**：`grep -rln "wisp/internal/ball" --include=*.go .`（剥 `.scratch`）⇒ **只有 `./cmd/balldebug/main.go` 一枚**。〔建了但没接〕（编排者已复跑同一把尺）

**`go.mod` 有没有 WebView2/webview 类依赖（现读，28 行整，逐字结论）**
**没有。** 直接依赖只有 6 枚（`go.mod:7-14`）：`sherpa-onnx-go`／`go-toml/v2`／`x/crypto`／`x/sys`／`x/term`／`modernc.org/sqlite`。`go.sum`（66 行）里 `webview|walk|systray|getlantern|go-ole` 全部 **0 命中**。
⚠ 但**别把这条读成"全仓没有 WebView2 代码"**：`github.com/jchv/go-webview2` **确实存在**，只在**另一枚模块**里——`scripts/spike/go.mod:9`＋`scripts/spike/go.sum:13-14`，消费者 `scripts/spike/webview2-latency/main.go:31`。那是 S0 的 P11 延迟探针，不进生产。另 `internal/panel/composer_dispatch_test.go:458` 拿 `go-webview2` 当**诱饵**造黑屏件（不是真依赖）。

**球在 Windows 上的技术形态（以现读代码为准）**
〔已证〕**不是 WebView2，不是纯 GDI，也不是 XAML**：是**原生 Win32 分层窗口（`WS_EX_LAYERED` + `WS_POPUP` + `Topmost` + `NoActivate` + `ToolWindow`）＋ GDI 内存 DC 上的 top-down 32bpp DIB section ＋ Direct2D（`ID2D1DCRenderTarget`）＋ `UpdateLayeredWindow` 呈现（预乘 alpha）**。逐字证据 `internal/ball/renderer_windows.go:4-9`：「an ID2D1DCRenderTarget bound to a top-down 32bpp DIB selected into a memory DC … the window presents it with UpdateLayeredWindow (premultiplied alpha, SPEC-08 §2)」。COM 指针一律 `unsafe.Pointer`，DWrite 只做文字（`:31-32`）。
⇒ **对写腿的直接含义**：球那一面**零新增依赖**，`go.mod` 不用动；要动的是"谁在 `wisp run` 里 `ball.New`"。

---

# ② 托盘的地形

**现有菜单项名册（`internal/ball/tray_windows.go`，逐枚抄，含 id 与行号）**

| id 常量 | 值 | 声明行 | 追加行 | 中文逐字 |
|---|---|---|---|---|
| `menuOpenPanel` | 1 | `:20` | `:86` | 「打开面板」 |
| 分隔条 | — | — | `:87` | （`mfSepart`） |
| `menuMute` | 2 | `:21` | `:88` | 「静音」（带 `checked`＝`muted`） |
| `menuPauseWake` | 3 | `:22` | `:89` | 「暂停唤醒」（带 `checked`＝`pausedWake`） |
| 分隔条 | — | — | `:90` | |
| `menuExit` | 4 | `:23` | `:91` | 「退出」 |

构建函数 `showMenu(hwnd, muted, pausedWake) uint32` 在 `tray_windows.go:72`，用 `TrackPopupMenu(tpmRightButton|tpmReturnCMD|tpmNoActivate)`（`:101`）取返回值 ⇒ **没有 `WM_COMMAND`/`WM_INITMENUPOPUP` 那条路**，菜单是"一次弹一次读"的同步形状。

**"打开面板"那枚空 stub 的确切位置（⚠ 票面口径错，见下面"票面 vs 现读"第 B-4 行）**
`internal/ball` 里**不存在空 stub**。ball 只做一件事：把选择发成回调——`ball_windows.go:664-681`：`wmLButtonUp ⇒ b.fire(b.opts.Events.OnTrayPanel)`（`:667`），`wmRButtonUp ⇒ showMenu(...) ; case menuOpenPanel: b.fire(b.opts.Events.OnTrayPanel)`（`:669-672`）。
**stub 在消费者那一侧**：`cmd/balldebug/main.go:199`
```
OnTrayPanel: func() { fmt.Println("tray: open panel (stub, ticket 33)") },
```
（同形还有 `:202` `OnTrayPauseWake`、`:208` `OnDragEnd`、`:611` 的 `hotkey: panel (stub, ticket 33)`。）

**菜单点击回调今天最终调用到哪一枚函数：无。**〔已证〕`Events`（`ball_windows.go:49-60`）的十枚函数全是宿主传进来的裸 `func()`，没有任何一处接到 `approval.*`／`Replies`／`Gate`；`internal/ball` 也不 import 审批包。`OnTrayMute` 在 `main.go:200` 只映射到 `statemachine.EvMuteHotkey`，`OnTrayExit`（`main.go:204-206`）只往一个 channel 塞字符串。**没有一条路径通向答复。**

**⚠ 要加四枚菜单项（允许一次／本次会话内允许／长期允许／拒绝）的最小落点**
1. `internal/ball/tray_windows.go:17-24` 加 4 枚 id 常量（现值只到 4），`:86-91` 加 4 枚 `appendItem`，`showMenu` 的签名要从 `(hwnd, muted, pausedWake)` 扩到能拿到"当前待答那一发"——**这是这面地形里唯一的价格**。
2. `internal/ball/ball_windows.go:49-60`（`Events`）＋ `:664-681`（`wmAppTray` 的 `switch sel`）各加 4 个 case 与 4 枚回调位。
3. 消费者侧把 4 枚回调指向 `approval.Replies`。**注意 `ReplyCard` 已经带够了信息**：`CorrelationID`／`Level`／`Tool`／`Grant`／`Paths`（`replies.go:70-87`），四枚按钮语义可直接映射到 `Replies.Allow(corr)`／`Replies.Veto(corr)`／`replySurface.always(corr)`（`approval_always.go:70`）／`Replies.Reject(ctx, corr, reason)`。
4. 现成"当前待答那一发"可读：**有，两枚**——`Replies.AwaitingHuman() (ReplyCard, bool)`（`replies.go:249`）与 `Replies.Pending() []ReplyCard`（`replies.go:224`），且 `Pending()` 按显示先后排序（`r.order`，`:231`），`MaxTrackedCards = 16`（`:49`）天然是托盘子菜单的长度上限。⚠ 但这两枚今天的生产调用者只有 `cmd/wisp`（见⑤）。

---

# ③ `wisp run` 的生命周期能不能挂常驻 GUI

**主循环与退出条件（现读行号）**
`cmd/wisp/main.go:52 func main()` → `:54 if len(args) == 0` ⇒ `:60 runResident()`；`case "run"`（`:83`）⇒ **`os.Exit(cmdRun(args[1:]))`（`main.go:85`）**。
`cmdRun`（`main.go:145`）→ `runTextTask`（`run.go:154`）→ `assembleRuntime`（`:207`）→ **`return rt.execute(task)`（`run.go:220`）** → `execute`（`:751`）：`bg := loop.RunAsync(ctx, task)`（`:848`）→ **`res := bg.Wait()`（`:855`）** → 写 `task_log`（`:869-889`）→ `return code`（`:914`）→ 进程 `os.Exit`。

**结论：`wisp run` 是"跑完一发任务就退"，不是常驻等活。**〔已证〕`ctx` 本身就是一次性超时（`:834 context.WithTimeout(…, cfg.LLM.TimeoutMS)`），没有任何 tick、队列或事件循环。真正的常驻入口是**另一条腿** `runResident()`（`cmd/wisp/resident_windows.go:24`，`proc.Boot` `:33`，`rt.RunEventLoop()` `:83`）——而那条腿**没有 gate、没有 bridge、没有任务**（票 201 §47 与 `cmd/wisp/approval_reply.go:48` 各自的自陈一致）。

**挂球/托盘要插在哪、会不会改 run 的退出语义**
插入点唯一合理处 = `run.go:596-602` 那一段旁边（`if s.reply != nil { rt.attachReplyListener(...) }`），把 `s.reply` 从"stdin"换成／加上"球＋托盘"，并在 `assembleRuntime` 末尾 `ball.New(...)`。
**会改退出语义，而且是三处**：
1. `run.go:220` 一发完就 `return`；球要活着就必须把 `execute` 后面挂一个等待，或者把 `wisp run` 变成 `runResident` 的子命令 —— 这是**改产品拓扑**，不是加一行。
2. **`rt.close()`（`run.go:676-683`）今天不 join 任何 goroutine**：只 `replyRoot.Cancel()`，`replyHandle`（`:265`，Spawn 于 `approval_reply.go:421`）**从不被等**。注释 `:672-675` 自己承认「the goroutine dies with the process」。GUI 线程不能这样收口（`Ball.Close()` 在 `ball_windows.go:904`，且 `staThread.start` 是**阻塞到 `WM_QUIT`** 的，见 `sta_windows.go:47`），所以 D38(e) 第 6 步（销毁面板 WebView→等子进程）在这里是缺的。
3. STA 亲和：`sta_windows.go:52 LockOSThread` 在 `Registry.Spawn("ui-sta","ball",…)` 的回调 goroutine 里执行 —— 该 goroutine 的 owner 目前是字符串 `"ball"`（`ball_windows.go:170`），`observe.ResidentBaseline`（`internal/observe/goroutine.go:44`）里 `ui-sta` 已在册，所以**用现成名额即可，不必新铸**。

**对照两条冻结决策（只读，原文＋行号）**

- **D2（`docs/PLAN.md:73-99`）**
 `:74`「**常驻原生进程树，空闲态无 WebView。**」
 `:83-88`：
 ```
 常驻（主进程，Go）                空闲上限见 D32（25MB / 40MB，视 D25 路径）
   ├── 悬浮球：原生分层窗口（Win32 layered window / macOS NSWindow）
   ├── 全局快捷键注册
   ├── 系统托盘
   ├── Job Object 持有者（C30…）
   └── KWS 唤醒词检测（opt-in…）
 ```
 `:91`「WebView（结果展示 / 配置面板）—— 独立子进程（msedgewebview2.exe × 3–5）」，`:92`「Agent 核心循环（Go，主进程内）」，`:103`「连锁推论：**Agent 核心循环必须在原生侧**」。
 ⇒ **裁定：允许，且是命令式的。** D2 把"球＋托盘＋全局热键"和"Agent 核心循环"画在**同一枚常驻主进程**里，只把 WebView 划到按需子进程。**D2 里没有一处按 `wisp run` 与无名 resident 入口分别说话** —— 它是拓扑图，不是 CLI 子命令许可。今天真正违 D2 的是"球根本不在生产进程里"。

- **D38（`docs/PLAN.md:2818-2848`）**
 `:2824`（表头那行，逐字）「**一个 STA/UI 线程** | 拥有：分层窗口 + 托盘 + **WebView2 窗口** + 所有 Win32 消息循环 | **WebView2 要求创建它的线程是 STA 且有消息泵**；**Direct2D factory 与设备有 thread affinity** → 悬浮球绘制与面板**必须在同一 STA 线程**。决定：**共用一个 STA 线程，不建第二个 D2D factory**」
 `:2831-2833` 常驻 6 枚名额含 `ui-sta`；`:2836` 每任务 3 枚含 `approval-waiter`；`:2838`「**总数上限：常驻 6 + 每任务 3。超出即为泄漏征兆**」；`:2845`「**任务完成必须等所有派生 goroutine 退出（WaitGroup）才算完成**」。
 ⇒ **裁定：允许，但附两枚硬边界，且有一枚今天已经在违约。** 允许：D38(a) 明确要求球/托盘/WebView 共一枚 STA 线程，且 `ui-sta` 名额已备好。边界：① 不许建第二个 D2D factory（`d2d_windows.go:206 FactoriesForPanel()` 就是为此预留的共用出口）；② 常驻上限 6 —— 在 `wisp run` 里加 `ui-sta` 会让这条 CLI 腿的常驻数从 0 变 1，不越限，但**上限是按"常驻"算的，run 腿不是常驻进程，D38 没为它写数**。违约：`:2845` 那句"必须等所有派生 goroutine 退出"对 `replyHandle` **今天就不成立**（`run.go:676-683` 不 join），挂 GUI 前这一格必须先补，否则是新账不是旧账。

---

# ④ 答复侧已有什么

**三枚签名 ＋ 生产调用者枚数（现跑 `grep -rn "DecideFromNative\|DecideFromPanel\|\.Veto(" --include=*.go internal/ cmd/ tools/ | grep -v _test.go`）**

定义（3 处）：
- `internal/agent/approval/gate.go:387 func (g *Gate) Veto(v Veto) error`
- `internal/agent/approval/gate.go:626 func (g *Gate) DecideFromNative(ctx context.Context, r Request) error`
- `internal/agent/approval/gate.go:638 func (g *Gate) DecideFromPanel(ctx context.Context, r Request) error`

**真调用者：生产 5 处，全部在 `internal/agent/approval/replies.go`**
- `replies.go:325` `g.DecideFromNative(...)` ← `(*Replies).Allow`（`:313`）
- `replies.go:377` `g.DecideFromNative(...)` ← `(*Replies).decide`（`:364`）的 native 分支
- `replies.go:372` `g.DecideFromPanel(...)` ← `(*Replies).decide` 的 panel 分支
- `replies.go:401` `g.DecideFromPanel(...)` ← `(*Replies).PanelAllow`（`:395`）
- `replies.go:427` `g.Veto(...)` ← `(*Replies).Veto`（`:419`）

其余命中全是注释（`gate.go:625/634`、`replies.go:7/22/93/94/95/389`、`approval_reply.go:12/22/31/182/274/282`）。
⚠ **编排者补一层（我与该腿差的那一行不是矛盾）**：同一把尺我另抓到 `cmd/wisp/approval_reply.go:310 s.live.h.Veto(corr)`，它调的是 **`(*Replies).Veto`**（`replies.go:419`）这一层，不是 `Gate.Veto`。⇒ **准确说法＝`Gate` 那三枚的生产调用者 5 枚（全在 `replies.go`）；`Replies` 那一层的生产调用者 5 枚（`approval_reply.go:210/249/251/278/310`）**。**票 201 §10 那把尺的读数"只剩定义行"已经过期**（见"票面 vs 现读" B-1）。

**再下一层：`cmd/wisp` 通过 `Replies` 用它们**
`approval_reply.go:210`（`h.Allow`）、`:249`（`h.PanelReject`）、`:251`（`h.Reject`）、`:278`（`h.PanelAllow`）、`:310`（`h.Veto`）——这五枚由 `replySurface.handle`（`:467-491`）分派，由 `runReplyLoop`（`:433`）驱动，由 `observe.Default.Spawn("approval-waiter", "approval", root, …)`（`:421-424`）拉起，拉起者是 `attachReplyListener`（`:389`），调用者 **`run.go:601`，且外面有 `run.go:600 if s.reply != nil` 的闸**，`s.reply` 在生产只有 `main.go:147 interactiveStdin()` 一枚来源。
⇒ **入口枚数 ＝ 1（控制台 stdin）。**〔已证〕

**`internal/agent/approval/replies.go` 导出名册（逐枚＋行号）**
常量 `MaxTrackedCards=16` `:49`；哨兵错误 `ErrNoTrackedCard :57`／`ErrRouteHasNoAllow :60`／`ErrNoGateAttached :64`；类型 `ReplyCard :70`／`Replies :101`／`HostBinding :124`；构造函数 `NewReplies :143`；方法 `Attach :152`、`Record :167`、`Look :196`、`Forget :208`、`Pending :224`、`AwaitingHuman :249`、`WaitingState :283`、`Allow :313`、`Reject :336`、`PanelReject :342`、`PanelAllow :395`、`Veto :419`、`Head :437`、`View :446`、`(ReplyCard) WideningRule :467`。未导出但承重的：`routes :348`、`decide :364`、`queueAwaitsHuman :303`、`isBareDrive :488`。
**文件没有 `//go:build` 标签**（`replies.go:1` 是 `package approval`）⇒ 非 Windows 也能编。〔已证〕

**`replyVeto` 的赋值点今天是否仍为零 ⇒ 是，仍为零**（编排者复跑同一把尺，四条命中逐字对上）。
- `cmd/wisp/run.go:146 replyVeto approval.Channel` —— 结构体字段**声明**
- `cmd/wisp/run.go:601 rt.attachReplyListener(s.reply, s.replyVeto)` —— **读**
- `cmd/wisp/run.go:440` —— 注释自陈「runSpec.replyVeto stays unset」
- `cmd/wisp/run.go:140-145` —— 字段注释
- **生产赋值点：0。** 只有测试在设：`cmd/wisp/approval_reply_201_test.go:97`（测试自己的镜像字段）、`:175`（把它递给 runSpec）、`:574 h.replyVeto = approval.ChannelEsc`。
`cmdRun`（`main.go:153-158`）构造 `runSpec` 时**只给 `argv/stdout/stderr/reply` 四枚**，从不给 `replyVeto`。
⇒ **AC#2 改写形"L1 等待期间能被否决"今天不可达**：`Gate.Veto`（`gate.go:388 g.channels.check(v.Channel)`）走的注册表是 `run.go:446 approval.NewChannels()`（**空参**，全部 `loaded=false`），而不是 `gate.go:108` 的 `DefaultChannels()` 兜底 ——⚠ **编排者把这条加重**：`Options.Channels == nil` 才会拿到 ball+esc 两枚已加载，`run.go` **显式传了空 registry＝主动放弃了那个兜底**（现跑：`grep -rn 'DefaultChannels()'` 生产命中只有 `gate.go:108`；`approval.go:163` 逐字 `return NewChannels(ChannelBall, ChannelEsc)`）。**"半成品就在场，装配处把它关掉了"这一句要写进派单。**

---

# ⑤ "当前待答那一发"能不能被原生面枚举到

**分两层，答案不对称。**

**L2 层：能，且已有两枚只读口。**
- `internal/agent/approval/pending_read.go:104 func (q *Queue) LiveApprovals() []LiveApproval`（`LiveApproval{CorrelationID, Decision(深拷贝), Position}`，`:46-50`；克隆在 `:57-89`）。**生产调用者 1 枚**：`cmd/wisp/panel_pump.go:62 items := rt.gate.Queue().LiveApprovals()`。（编排者复跑同一把尺：非测试命中就这 1 枚＋定义行）
- `internal/agent/approval/replies.go:224 Pending()` / `:249 AwaitingHuman()` / `:283 WaitingState()`：只读、无授权（文件头 `:96` 逐字「read-only, no authority at all」）。生产调用者：`approval_reply.go:162/172`（thin wrapper）→ `approval_always.go:175/178/181`（`bookWaitingState`），驱动者 `run.go:1110`（`ui.Prompt`）与 `run.go:1156`（`ui.Update`）。

**L1 层（`Gate.windows`）：不能，全仓零枚枚举口。**〔已证，编排者复跑 `grep -n 'g\.windows'`＝只有 `:321/:324/:330/:393` 四处〕
`gate.go` 里对 `g.windows` 的操作只有四枚：`:321` 查重、`:324` 写、`:330` 删、`:393` 按 corr 查。**`(g *Gate)` 的 13 枚导出方法**（`:135 Queue`／`:139 Channels`／`:142 Window`／`:160 AdmitTextTask`／`:218 PendingWindow`／`:355 Complete`／`:387 Veto`／`:448 LateVeto`／`:468 PendingApproval`／`:588 Native`／`:592 Panel`／`:626 DecideFromNative`／`:638 DecideFromPanel`／`:657 Replay`）里**没有任何一枚列出 windows**。`Queue` 的导出面只有 `Timeout :126`／`WarningLead :129`／`Depth :133`／`LiveApprovals :104`；`head()/view()` 是未导出的，只能经 `PanelAPI`（`gate.go:602-608`）拿到，而 `PanelAPI` 刻意不含 `Params`／`Grant`。
⇒ **托盘要在球旁边说"现在有人在等的这一发"，L2 够用（`AwaitingHuman`），L1 那一发今天只能从宿主自己的账上读**（`Replies.Pending()` 里有 `Level=="L1"` 的行，因为 `consoleApprovalUI.Prompt` 对两类卡都 `record`，见 `run.go:1102-1103`）——**所以"新面"其实不是必需的那一枚**。

**若仍要新枚只读枚举口的最小形状 ＋ `ModeSource` 样板**
样板：`internal/tools/mode.go:27-31`
```
type ModeSource interface {
	// PermissionMode returns the mode in effect. It must not block, and it must
	// not be able to change anything.
	PermissionMode() risk.Mode
}
```
它的三个特征都可直接照抄：**只有一枚方法、名字里带 Source、注释里明写"不得阻塞、不得改变任何东西"**；消费端 `internal/tools/bridge.go:90 Modes ModeSource`／`:135 modes ModeSource`，读取端 `mode.go:35-47 permissionMode()` 带 recover 且 **fail-closed 到最严档**（`:37 return risk.DefaultMode()`）。票 224 §11 已经用同一枚话头立过例（"必须新增一枚只读 seam（照现成 `ModeSource` 那形），这是新契约面、要先落 `A##`"）。
⇒ **本票对应的新面最小形状**＝在 `internal/agent/approval` 加一枚只读枚举（例如"列出此刻在等人的那一发"），**但它撞票 220 AC#2(甲)**（见⑦）。**建议：走 `replyVeto`＋`Channels.SetLoaded` 那条已备好的门，而不是新造枚举**——`approval.go:147 NewChannels(loaded ...Channel)` 与 `gate.go:139 Channels()` + `approval_reply.go:415-417 if vetoChannel != "" { SetLoaded(true) }` 三件已经齐了，只差**有人把 `runSpec.replyVeto` 设上**。

---

# ⑥ 三枚按钮到现成枚举的映射是否成立

**词表（现读，逐枚抄＋行号）**
`internal/tools/gate.go:75-89`，`type Answer string`（`:74`）共 **4 枚**（编排者复跑：`allow`／`reject`／`veto`／`timeout` 逐枚对上）：
- `AnswerAllow Answer = "allow"` `:78` —— 注释「execute」
- `AnswerReject Answer = "reject"` `:80` —— 「user/queue said no」
- `AnswerVeto Answer = "veto"` `:83` —— 「the L1 pre-execution window was vetoed (B1)」
- `AnswerTimeout Answer = "timeout"` `:88` —— 注释逐字「the L1 window ran out unopposed, which per SPEC-06 §2 **MEANS EXECUTE**; for the L2 queue it means auto-reject … **The branch is the caller's, never the gate's**」

**⚠ 票 219 §49 说"本次＝一条 `AnswerAllow`"，映射成立；但"三枚按钮"里的第二枚（本次会话内同类）在这张词表里没有对应的值**——`Answer` 只有 4 枚，且都不是"会话档"。这一点票 219 §24/§68 与裁决表 §57（`docs/evidence/s1/219-approval-reply-design-adversarial.md:57`）都自己承认：`GrantScopeSession` 是**另一个词的表**（`approval_grant.scope`），不是 `Answer`。

**`GrantScopeSession`：定义在场，生产写手 0。**〔建了但没接〕
- 定义：`internal/memory/models.go:105 const GrantScopeSession = "session"`，注释 `:104` 逐字「the only scope value (SPEC-02 §3 approval_grant.scope)」。
- **唯一读它的是拒绝式守卫**：`internal/memory/dao_misc.go:18 if g.Scope != GrantScopeSession { … }`（`InsertGrant` 的校验，不是写手）。
- **DAO 三枚全部生产零调用者**（现跑 `grep -rn "InsertGrant\|RevokeGrant\|ListGrantsBySession" --include=*.go internal/ cmd/ | grep -v _test`）：只剩 `dao_misc.go:16/17`、`:48/49`、`:77/78` 的注释与定义行本身。⇒ **`approval_grant` 表今天没有任何生产写手**，与票 201 §43 的"AC#5 不翻勾"一致。
- 另有两处**只谈不写**：`internal/config/permmode.go:17`（注释：不持久化 session grant）、`internal/perm/store.go:24`（注释）。

**`internal/config/allowdirs.go`：`AddAllowedDir` / `SetAllowedDirs` 的生产调用者枚数**
- `AddAllowedDir` 定义 `allowdirs.go:62`（注释 `:49`）⇒ **生产调用者 1 枚：`cmd/wisp/approval_always.go:134`**（编排者复跑复认）。测试调用者集中在 `internal/config/writeguard_226_test.go`（`:67/:122/:149/:168/:190/:216/:336/:375/:400`）。
- `SetAllowedDirs` 定义 `allowdirs.go:109`（注释 `:95` 逐字「the 撤销 half of AddAllowedDir」）⇒ **生产调用者 0 枚，但有 2 枚测试调用者**（`writeguard_226_test.go:313`、`:350`）。**票 201 §43 那句"全仓零调用者、零用例"已过期**（见 B-5）。
- `AllowedDirs()`（`allowdirs.go:42`，只读便利）⇒ **全仓零调用者**。生产实际读的是裸字段：`cmd/wisp/run.go:387-389 for _, d := range cfg.FS.AllowedDirs`。
- 内部：`writeAllowedDirs :139`、`statOwnWrite :157`。
- ⚠ 装配侧未变：**编排者现读 `cmd/wisp/run.go:341 mgr, err := config.NewManager(cfgPath, nil)`** ⇒ 确认钩子仍传 `nil`（与票 223 §7 一致，**这一行是票 223 AC#6 的靶心**）。

**映射是否成立的裁定**
- **本次 → `AnswerAllow`：成立**，且已有真执行者（`replies.go:313 Allow` → `gate.go:626` → `queue.go:364 deliver(answer{a: tools.AnswerAllow,…})`）。
- **长期 → 往 `allowed_dirs` 加一行：成立且已接通**（`approval_always.go:70 always` → `:84 WideningRule` → `:99-102` 另起一根 goroutine 打第二张 L2 卡 → `:124` 判 `AnswerAllow` → `:134 AddAllowedDir`）。**可撤销那一半不成立**（`SetAllowedDirs` 零生产调用者）。
- **本次会话内 → `GrantScopeSession`：不成立**，三件全缺（铸造会话身份的写手 0、`approval_grant` 写手 0、路径模式匹配器全仓不存在）。⇒ **这一格归票 224，不该由原生入口腿碰。**

---

# ⑦ 写腿动手前会撞谁

## A. 即将新开的这枚写腿要碰的文件（该腿的预测清单，非现读事实）

| 文件 | 为什么碰 | 现有行数 |
|---|---|---|
| `internal/ball/tray_windows.go` | 加 4 枚 id＋4 枚 `appendItem`，`showMenu` 签名要能拿到待答卡 | 109 |
| `internal/ball/ball_windows.go` | `Events`（`:49-60`）加 4 枚回调位；`wmAppTray`（`:664-681`）加 4 个 case | 972 |
| `cmd/wisp/run.go` | 装配 `ball.New`；设 `s.replyVeto`；改 `Channels`；**改 `rt.close()` 的 join** | 1176 |
| `cmd/wisp/approval_reply.go` | 新增"球/托盘 → `replySurface`"分派；`attachReplyListener` 的 `in io.Reader` 形状要撑住非流式宿主 | 539 |
| `cmd/wisp/main.go` | `cmdRun` 的 `runSpec` 要带上 veto 通道与宿主开关 | 166 |
| **新建** `cmd/wisp/native_host_windows.go`（暂名） | 球/托盘生命周期；必须带 `//go:build windows` | — |
| `cmd/wisp/approval_always.go` | 只可能改 `bookWaitingState` 的下游（真驱动球） | 195 |
| `internal/agent/approval/replies.go` | 只有在要新铸只读枚举面时才碰 | 491 |
| `internal/agent/approval/gate.go` | 若走票 220 甲形（暴露 `windows`）才碰 | 695 |
| `go.mod` | **不需要动**（球是纯 Win32＋D2D；只有走 WebView2 面板才需要 `go-webview2`） | 28 |

## B. 与队列其它票的写面求交集（同文件互斥，具名）

| 别票 | 它要碰 | 与本腿的**同文件互斥** |
|---|---|---|
| **220** | `internal/agent/approval/gate.go`（AC#2 **甲形＝"给 Gate 补一枚只读枚举口把 `windows` 暴露给泵"**，票面 §25）、`pending_read.go`、`internal/panel/subagent_roster_197.go:190-210`、`internal/panel/pump.go`、`cmd/wisp/panel_pump.go` | **最硬的一枚。** 它要在 `gate.go` 上造的新面，正是本腿第⑤问里那枚"最小新面"；两腿各造一次＝两份枚举器。并且票 220 的禁区自己写了「票 201 的写面（`internal/agent/approval`＋`internal/tools`＋`cmd/wisp`）未空出之前不许派本票（同文件互斥）」。⇒ **必须串行，且先裁甲/乙形。** |
| **223** | `cmd/wisp/run.go`（加 tick 调 `CheckAndReload`）、`internal/config/manager.go` | `run.go` 互斥。⚠ 还有语义耦合：票 223 一接上轮询，本腿写的 `allowed_dirs` 行会被"读回"再判一次方向（票 226 §42 已把这条顺序写死：226 必须先于 223 或同批）。 |
| **224** | `internal/tools/bridge.go`（新只读 seam）、`cmd/wisp/run.go`（注 `memory.Store`）、`internal/memory/dao_misc.go` | `run.go` 互斥；且**共用同一枚"照 `ModeSource` 那形新铸只读 seam"的位置**——两枚 seam 若同名或同处，先落的那枚会挡后落的。 |
| **226** | `cmd/wisp/approval_always.go`、`internal/config/allowdirs.go`（`writeguard.go`/`loader.go`/`permmode.go`） | **`approval_always.go` 与 `allowdirs.go` 双双互斥。** 若本腿顺手把"可撤销"接上（调 `SetAllowedDirs`），就直接改在 226 的射程里。 |
| **222** | `internal/tools/subagent_197.go`、`cmd/wisp` | 只撞包不撞文件（本腿不写 `subagent_197.go`）。票面 §46 明写「本票与票 201、票 221 都碰 `internal/tools/subagent_197.go`／`cmd/wisp` ⇒ 一律串行」。 |
| **221** | `internal/tools/task.go`、`internal/panel/subagent_197.go`、`cmd/wisp` | 同上，只撞包。 |
| **213** | `internal/agent/control.go:36` 一带解析器、`internal/panel/composer.go` | **不撞**。 |
| **214** | `internal/panel/attachments.go`、`cmd/wisp/run.go` 装配段、可能加配置键 | `run.go` 互斥。 |

**去重后的硬互斥清单：`internal/agent/approval/gate.go`（220 甲形）、`cmd/wisp/run.go`（223／224／214／220）、`cmd/wisp/approval_always.go`＋`internal/config/allowdirs.go`（226）。** 包级串行的还有 221／222。

## C. `scripts/d22scan.sh` 会天然打到这里的那几枚

先记一句：`scripts/d22scan.sh:51` 会跑 `runtests.sh`（＝`go test -count=1`）、`:54` 会跑 `go run .` —— **本腿两枚都没跑**（硬纪律）。

**ban #1 `bare-goroutine`（一定打中）**
射程＝**任何 `go <anything>`**，闭包与具名调用都算（`tools/d22scan/main.go:697-711`；`:705` 判 `*ast.FuncLit` 给"bare `go func(`"句，`:710` 给具名调用句，R16#1 已把具名算进来）。唯一按文件豁免＝`internal/observe/goroutine.go`（`tools/d22scan/allowlist.txt:7`，逐字「Exempted BY FILE PATH only (never by call shape)」）。
⇒ 球宿主那条 STA 腿**必须**走 `observe.Registry.Spawn`（`ball_windows.go:170` 已经是这个形状，照抄即可），`go` 一个字都不能出现。⚠ 陷阱：`approval_always.go:99` 已经示范了"复用 `approval-waiter` 名额"的写法并注明其代价（`:30-38`），新开宿主 goroutine 若想要独立名额就是 D38b 契约动作。

**ban #4 `wallclock-timeout`（很可能打中）**
实现 `tools/d22scan/main.go:144-146`：
- `wallclockRe = \.Sub\(time\.Now\(\)\)` —— 逐字只匹配 `X.Sub(time.Now())` 这一串（**方向敏感**：`time.Now().Sub(X)` 反而不中）。
- `unixTimeRe = time\.Now\(\)\.(Unix|UnixNano|UnixMilli)\(` **且** `timeoutWordRe = (?i)timeout|deadline|expire|\bttl\b|budget|until` 同在一行 ⇒ 两枚都中才报。
判定循环 `:758-773`：⚠ **它只跳过以 `//` 开头的整行（`:761-763`），代码行尾的注释仍然算进这行** —— 所以"把 `deadline` 写在行尾注释里＋同一行写 `time.Now().Unix()`"会红。
⇒ 托盘/球要做"卡片剩余秒数"倒计时就有风险：`ReplyCard.Deadline`／`.Window` 是 `time.Duration`（`replies.go:85-86`），拿墙钟去算差就是撞这枚。安全形状：用 `context.WithDeadline`/monotonic `Duration`，别写 `.Sub(time.Now())`。（现成参照：`gate.go:310 return tools.AnswerTimeout, ""` 那一支靠 channel/timer，不做减法。）

**ban #8 `emoji`（一定打中，且是最容易白吃的一枚）**
现读 `emojiRe`（`tools/d22scan/main.go:164`，整行逐字）：
```go
var emojiRe = regexp.MustCompile(`[\x{1F000}-\x{1FAFF}\x{2200}-\x{22FF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]`)
```
**实际扫的码段（6 段，逐字结论）**：`U+1F000–U+1FAFF`、`U+2200–U+22FF`（数学符号）、`U+2600–U+27BF`（杂项符号与装饰符）、`U+2B00–U+2BFF`、`U+FE0F`、`U+1F1E6–U+1F1FF`（区域指示符，与首段重叠）。
**刻意不含**（注释 `:157-161` 逐字）：`U+2190–U+21FF`（箭头）、`U+2460–U+24FF`（带圈数字），另 `U+2014 —`、`U+2018/2019`、`U+00A7 §`、全部 CJK 也不在段内。
⇒ **对写腿的可执行结论**：**⚠(U+26A0)、✅(U+2705)、⛔(U+26D4)、✗(U+2717)、✓(U+2713)、☐/☑、❌/⭐ 全部落在扫描段内 ⇒ 出现在非注释文本里就红**；而 `← ⇒ ≥ — 「」§` 安全。**豁免只有注释**（`walkEmoji` `:977-984` 先用 `commentRangesFor` 把注释区间挖空再匹配，Q-46(c)），**字符串不豁免**（`:998-1004` 逐字「a glyph inside any Go string literal always survives」）。
⇒ **最现实的雷**：托盘菜单标签（`tray_windows.go:86-91` 那种 `appendItem(id, "打开面板", false)` 的**字符串字面量**）里放 ⚠/✅/⛔ ⇒ 红。同理"把「⚠ 二次确认」印进 `approval_always.go:103` 那种用户可见字符串"会红。

**顺带会中的两枚（低概率但记一笔）**：ban #2 `pathresolver-bypass`（`main.go:713-721`，`filepath.Clean`/`filepath.Abs`）——若写腿要算菜单/图标路径；ban #7 `internal-artifact-tool`（射程 `internal/tools/`，本腿不碰）。

## D. 冻结件／命名雷（写腿必避）

- `internal/panel/l2_grant_boundary_test.go`（现读存在）—— **禁名词表实际位置 `:1881-1886`**（票 219 §31 写 `:1882-1886`，差一行），`candidates` 24 个串 = **12 词根 × 大小写两拼**：`outcome/allow/allowOnce/approved/grant/verdict/decision/decide/bypass/override/permit/authorize`。AST 扫描的自陈范围是 `internal/panel`（`:1579` 的 `t.Logf` 逐字「AST scan of internal/panel: 0 findings」，另 `:127`）。⚠ **该腿自陈没有追进 `goSourceFiles`（`:360` 定义、调用点 `:726`／`:2024`）拿到 `dir` 字面值 ⇒ "只扫 internal/panel"这一格只到注释自陈级、未行级证**；保守做法＝**新增可解码结构体一律不放 `internal/panel`**，答复字段只叫 `reason`。
- `internal/perm/ticket90_persist_test.go`、`internal/panel/tokens_fourway_test.go` —— 冻结件一字不动。
- `docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`tools/d22scan/allowlist.txt` —— 一字不动。
- `frontend/**`／`design/**` —— 本腿**零读**；只登记"票面引用了 `frontend/src/lib/panel.ts`（票 201 §14、`internal/panel/pump.go:24` 的注释也提到该文件）"这一事实，未打开。⚠ 工作树里 `design/**` 有大量别人的未提交删除（`git status` 可见），不属本腿射程。

---

# 票面写的 vs 我读到的（单列，一律现读）

| # | 票面（何处） | 该腿现读到 | 性质 |
|---|---|---|---|
| B-1 | 票 201 §10 尺：`DecideFrom*`/`.Veto(` **生产零调用者** | 生产真调用 **5 处**，全在 `internal/agent/approval/replies.go:325/372/377/401/427`（编排者复跑另见 `approval_reply.go:310` 属 `Replies` 层，见④） | **已过期**（票面 r1/r2 两节自己解释了原因） |
| B-2 | 票 201 §35：`run.go:593 rt.attachReplyListener` | **`run.go:601`**（外加 `:600` 的 `if s.reply != nil` 闸） | 行号漂 8 |
| B-3 | 票 201 §35：`approval_reply.go:204` 调 `DecideFromNative`、`:265` 调 `DecideFromPanel` | `approval_reply.go` 今天**不调任何 `DecideFrom*`**：`:208 allow()`→`:210 h.Allow`；`:276 panelAllow()`→`:278 h.PanelAllow`。真正的路由已迁到 `replies.go` | **函数已搬家**，按行号办事会指错 |
| B-4 | 票 201 §13／§55：托盘"打开面板"那项是**空 stub** | `internal/ball` 里没有 stub；它只是 `b.fire(Events.OnTrayPanel)`（`ball_windows.go:667/672`）。**stub 在消费者**：`cmd/balldebug/main.go:199` | **口径错**（把"没人实现回调"说成"包里有空 stub"）——**影响写腿对落点的判断**：球那一面已经把选择发出来了，缺的是听众 |
| B-5 | 票 201 §43：`SetAllowedDirs` **全仓零调用者、零用例** | 零生产调用者 ✔，但**有 2 枚用例**：`internal/config/writeguard_226_test.go:313`、`:350` | 用例数已过期（票 226 那批件落的） |
| B-6 | 票 219 §12 自称"现行号" `gate.go:622/610/371` | 现读 **`:638 DecideFromPanel`／`:626 DecideFromNative`／`:387 Veto`** | 票面"更正后"的行号又漂了 8/16/16 |
| B-7 | 票 219 §27：「**没有 `SetAllowedDirs`**，Manager 只有 `SetPermissionMode`」 | `internal/config/allowdirs.go:109` 有 | 已过期 |
| B-8 | 票 219 §27：`config.NewManager(cfgPath, nil)` 在 `run.go:334` | 现读 **`run.go:341`**（另 `cmd/wisp/panel_inbound.go:200` 一处同形） | 行号漂 7；票 219 §63 自己预告过（"312/334/341 三枚在流传"） |
| B-9 | 票 201 §43：`replyVeto`「`run.go:146` 只有声明、`:601` 只有读」 | **逐枚复认成立**（编排者自己跑） | ✅ 一致 |
| B-10 | 票 201 §55：`run.go:1110/:1156 → approval_always.go:171-186` | 调用点 `:1110/:1156` ✔；`bookWaitingState` 函数体是 **`:171-184`**，`:186` 起是另一枚 `waitingStateName` 的注释 | 尾部多读 2 行 |
| B-11 | 票 219 §31：禁名词表在 `:1882-1886` | 现读 `candidates := []string{` 起于 **`:1881`**，表体到 `:1886` | 差一行 |
| B-12 | 票 201 §13：球「20 态视觉齐」 | `cmd/balldebug/main.go:51-72` 的 `allStates` 确为 **20 枚**（`:52-71`），且 `internal/ball/statevisual.go:177`（`StateConfirming`）／`:181`（`StateAwaitingApproval`）两行都渲染 `RingColor=Danger` | ✅ 一致 |
| B-13 | 票 201 §14：`internal/panel/bridge.go` 那 4 枚入向名册 vs 页面发的 `panel.approval.request` | **未验**（要读 `frontend/**`＝禁区） | 未判 |

---

# 我没能量到 / 判不了（该腿自陈，完整清单，不编）

1. **本件的落盘不是我做的**：该腿无 Write／Edit 工具，`wc -c` 自证与 `git log -1` 交件自证因此都无从执行 ⇒ 由编排者代落并代跑尺（见第 0 节）。
2. **零编译期事实。** 按硬纪律没跑 `go build`／`go vet`／`go test`／`gofumpt`／任何门禁。因此下面几格**只到"读到了字面"级、不能算证**：
 a. `cmd/wisp`（`//go:build windows`）里 `import "…/internal/ball"` 会不会成 import 环／撞 tag —— **没有逐枚读 `internal/ball/*.go` 的 import 块求交**（只确认 `ball_windows.go` 引 `observe`/`statemachine`/`x/sys/windows`、`hotkey_reload.go` 引 `config`）。"看起来无环"不等于编译器认。
 b. `Replies` 那套在无 `windows` tag 下能否编（读到 `replies.go` 无 tag、只 import `statemachine`；但没跑）。
 c. `internal/ball` 那 11 枚无 tag 文件（`anim.go`/`doc.go`/`dock.go`/`hit.go`/`liquid.go`/`position.go`/`statevisual.go`/`tokens.go` ＋ 3 枚测试）与 15 枚 windows-tag 文件的切面，在 `GOOS=linux` 下会不会因新宿主件而露出未定义符号。
3. **`internal/panel/l2_grant_boundary_test.go` 的 AST 扫描真实根目录未行级证**（只读到 `goSourceFiles` 定义 `:360`、调用点 `:726`／`:2024`）⇒ 禁名词表"只扫 internal/panel 还是扫到 cmd/"**未判**（第一证据是 `:1579` 与 `:127` 的自陈，属注释级）。
4. **球的真实资源代价未测**（D32 的 25MB/40MB 常驻口径、handle 数 `<600` 那扇 SLO 门）。那要跑进程与 SLO 仪器。`cmd/balldebug/main.go:78-85 handleCount()` 是现成量具，但它在 debug 件里。
5. **`observe.ResidentBaseline` 之外有没有"run 腿"的名额口径未判**：读到 `internal/observe/goroutine.go:44` 常驻 6 枚名单含 `ui-sta`，但**没读到 Registry 是否按"进程是否常驻"分别计上限**，也没读 `RosterReport.Unknown` 的判据 ⇒ "在 `wisp run` 里起 `ui-sta` 会不会被看门狗报泄漏"**判不了**。这条直接决定第③问的代价。
6. **`RunEventLoop()` 里到底能挂什么未读**（归 `internal/proc`，该腿没打开那包）—— 而"球该挂在 run 还是挂在无名 resident 入口"这枚选择**取决于它**；只交得出地形（`wisp run` 一发就退、resident 无 gate），交不出结论。
7. **票 201 §14 的"界面白名单对不上"那一格未验**：要读 `frontend/**`＝禁区，零读零转述。
8. **`internal/agent/approval/ui.go:35/54/99-111/143-158` 只在别的文件的引用句里见过这些行号，没逐行打开 `ui.go`**；`NativeAPI`/`PanelAPI` 的**方法形状**是从 `gate.go:594-608`（实现体）读到的，那个更硬。
9. **`rt.close()` 不 join `replyHandle` 读到了（`run.go:676-683`），但"D38(c) 那句'必须等所有派生 goroutine 退出'在本仓有没有任何仪器在查"未查** ⇒ 第③问里"这一格今天已经违约"是**读到码得出的，不是读到尺得出的**。

---

# 七问各一句话结论（该腿原文）

1. **球**：`cmd/balldebug/main.go:87 main → :189 ball.New → ball_windows.go:143 → :166/:170 observe.Registry.Spawn("ui-sta") → sta_windows.go:47 start（LockOSThread＋CoInitializeEx(STA)＋GetMessage 泵）→ ball_windows.go:182 createOnSTA（分层窗口＋D2D＋托盘＋热键）`；引 `internal/ball` 23 枚导出符号；**`go.mod` 里没有 webview 类依赖（0 命中，`go-webview2` 只活在 `scripts/spike/go.mod:9` 那枚旁支模块）**；球在 Windows 上是**原生 Win32 `WS_EX_LAYERED`+`WS_POPUP` 窗口 + GDI 内存 DC/DIB + Direct2D `ID2D1DCRenderTarget` + `UpdateLayeredWindow`**（`internal/ball/renderer_windows.go:4-9`、`ball_windows.go:196/212`）——**不是 WebView2，也不是纯 GDI**。
2. **托盘**：`internal/ball/tray_windows.go:86/88/89/91` 四枚（打开面板／静音／暂停唤醒／退出，id `:20-23`），点击落到 `ball_windows.go:664-681` 的 `Events.OnTray*`，**最终调用到的答复函数：无**；票面说的"空 stub"在 `cmd/balldebug/main.go:199` 而**不在 ball 包里**。加四枚的最小落点＝`tray_windows.go:17-24` 的 id＋`:86-91` 的 `appendItem`＋`showMenu` 签名扩一枚"待答卡"参数＋`ball_windows.go:49-60/664-681`，**现成可读数据有**：`Replies.AwaitingHuman()`（`replies.go:249`）与 `Pending()`（`:224`）。
3. **run 生命周期**：`wisp run` **一发即退**（`main.go:85 os.Exit(cmdRun())` → `run.go:220 return rt.execute(task)` → `:855 bg.Wait()` → `:914 return`），常驻的是另一条无 gate 的腿（`resident_windows.go:24/:83`）；挂球要插在 `run.go:596-602` 旁边并**改 `rt.close()`（`run.go:676-683`，今天不 join `replyHandle`）** ⇒ **一定改退出语义**；**D2（`PLAN.md:74/83-88/103`）与 D38(a)（`:2824`）都允许、甚至是命令式要求球＋托盘与 Agent 循环同处一枚常驻原生主进程／一枚 STA 线程**，但**两条都没有一句话按"`wisp run` 这条 CLI 腿"分别许可或禁止**，而 D38(c)（`:2845`）"必须等所有派生 goroutine 退出"今天已对答复监听器违约。
4. **答复侧**：`DecideFromNative(ctx,Request)error`（`gate.go:626`）／`DecideFromPanel(ctx,Request)error`（`:638`）／`Veto(Veto)error`（`:387`）—— **`Gate` 那一层的生产调用者 5 枚，全在 `replies.go:325/372/377/401/427`**；`replies.go` 导出名册 21 枚（见④）；**`replyVeto` 生产赋值点今天仍为 0** ⇒ **票面那条断言成立，AC#2 改写形今天仍不可达**。
5. **枚举"当前待答那一发"**：**L2 能**（`Queue.LiveApprovals()` `pending_read.go:104`，生产调用者 `cmd/wisp/panel_pump.go:62`；另有 `Replies.Pending/AwaitingHuman`），**L1 不能**（`Gate.windows` 全仓零枚枚举口）；要新铸的最小面照 `ModeSource`（**`internal/tools/mode.go:27-31`**）那形：单方法、名字带 Source、注释自陈"不阻塞／不改任何东西"、消费端 fail-closed。**但这一格撞票 220 AC#2 甲形，动手前必须先裁甲/乙。**
6. **三枚按钮的映射**：**本次＝`AnswerAllow`（`internal/tools/gate.go:78`）成立且有真执行者**；**长期＝往 `[fs] allowed_dirs` 加一行成立**（`approval_always.go:134 → allowdirs.go:62`，生产调用者 1 枚），**但"可撤销"那一半不成立**（`SetAllowedDirs` `allowdirs.go:109` 生产 0 枚调用者、2 枚用例）；**"本次会话内"在 `Answer` 词表里没有对应值，且 `GrantScopeSession`（`memory/models.go:105`）生产写手 0、DAO 三枚全零调用者** ⇒ 那一格归票 224，不该由本腿碰。
7. **会撞谁**：同文件互斥＝**票 220（`gate.go`）**、**票 223／224／214（`cmd/wisp/run.go`）**、**票 226（`approval_always.go`＋`allowdirs.go`）**，包级串行＝票 221／222（`cmd/wisp`），不撞＝票 213。d22scan 天然打中三枚：**ban #1 任何 `go <anything>`**、**ban #4 墙钟那两形**、**ban #8 emoji**；`emojiRe`（`tools/d22scan/main.go:164`）**实际只扫 6 段：U+1F000–1FAFF、U+2200–22FF、U+2600–27BF、U+2B00–2BFF、U+FE0F、U+1F1E6–1F1FF** ⇒ **⚠/✅/⛔/✓/✗ 写进任何非注释文本（含托盘菜单标签、用户可见字符串）必红，而 `← ⇒ ≥ — 「」§` 与全部 CJK 安全；注释豁免、字符串不豁免。**

---

# 自报：本腿有没有跑过编译/测试（该腿逐字）

**没有跑过任何一枚编译、测试、vet、lint、门禁或 SLO 命令。** 用到的命令全集（只读）：`git status --short`、`git rev-parse`、`git log -1 --oneline`（两次，读到 `29daf745`→`721e3864`）、`git diff --stat 29daf745..721e3864`、`git ls-files`、`ls`、`wc -l/-c`、`grep -rn`、`sed -n`、`head`、`for`+`grep` 做 build-tag 普查。**没有 `mkdir`、没有重定向、没有 `>`/`>>`、没有写任何文件。**（编排者复核：本件目录 `.scratch/wisp/probes/201/c2/` 今天由**我**新建，该腿交件时盘上确实只有 `r2/`。）

> **该腿自己提醒的读数时效**（照抄）：「第①②④⑤⑥节的符号级结论在 `721e3864` 上仍成立；**唯一会漂的是 `internal/config/**` 相关行号（`allowdirs.go:42/62/109/139/157`、`run.go:341`）**，因为票 226／223 那条线正在改那三包。引用第⑥节的 allowdirs 行号请按符号名重取，别按这几个数字办事。」
