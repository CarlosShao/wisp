# hops.md — 常驻腿 → 面板 这一条链的逐跳普查（leg `panel-resident-1`，只读）

## §0 起手锚

- 仓根：`D:\work\workspace\projects plans\Wisp`
- HEAD（起手现量）：`82a10da4`
- `git status --porcelain | wc -l` = **752**（起手；共享工作树，别人的活在里面）
- **本腿没跑过 `go test`、没跑过 `go build`、没跑过任何突变**（272-r2 在写 `cmd/wisp` 测试，整包跑会洗它的红名册）。
  本腿全部尺 = `grep` / `sed -n` / `ls` / `wc` / `git ls-files` / `git check-ignore`，只读。
- **本腿一行产码都没改**，只新建本件与 `.scratch/wisp/probes/panel-resident/1/logs/`（临时件只建不删）。
- 尺纪律：判"不存在"之前先跑**正向对照**（同一把尺含 `_test.go` 与不含 `_test.go` 各一枚数），
  且"命中数"与"调用者枚数"分开报（定义行不算调用者）。原始输出落 `logs/`。

一句话现状（与编排者开场那句**不一致**，见 §3 冲突 #1）：**面板宿主今天已经接进常驻腿**——
`wisp` 无参数启动时，装配根会建出 `PanelManager`、起一枚专用 STA 线程、跑库的 `Run()` 泵，
球的 `[panel]` 热键与托盘项都能把"开窗"请求投到那条线程上。
**真正一块没写的是"Go 把快照推给页面"那一条边**：任务管线与面板窗口之间今天**零引用**。

---

## §1 跳表 A–F

### A. 常驻腿启动路径：两枚入口各自起了什么，装配根在哪一行

**存在吗**：存在，两枚入口都在，且常驻那枚的装配根是具名的一整段。

| 入口 | 行 | 起了什么 |
|---|---|---|
| `wisp`（无参数＝GUI 常驻） | `cmd/wisp/main.go:60-67` → `runResident()` | 见下 |
| `wisp run "task"` | `cmd/wisp/main.go:89-91` → `cmdRun`（`:151`）→ `runTextTask`（`cmd/wisp/run.go:181`）→ `assembleRuntime`（`run.go:382`） | 任务管线＋快照泵＋**stdin 答复环**；**不开任何窗口** |

装配根＝`cmd/wisp/resident_windows.go:33` `runResident()`，按文件内顺序各跳：

| 建了什么 | 行（现量） |
|---|---|
| `proc.Boot` / 日志 sink / D38(e) 退出 defer | `:42`、`:66`、`:75` |
| **审批门**（票 246 裁的"由装配根注入"那一形） | `:132` `ra := newResidentApprovalWithConfig(rt.Layout.DataDir)` |
| **面板宿主装配**（票 33） | `:151` `rp, rpErr := newResidentPanelManager(rt.Layout.DataDir)` |
| **面板线程起立** | `:157` `panel = startResidentPanel(rt.Registry, rp)`（`:158` `defer panel.stop()`） |
| 球＋热键，并把面板执行器交出去 | `:217-219` `startResidentBall(..., withPanelHost(func(via string) bool { return panel.RequestToggle(via) }))` |
| 球与门的绑定 | `:228` `ra.bindBallHost(rb)` |
| D38(e) 第 3 步钩子 | `:240` `rt.RegisterShutdownHook(proc.StepCancelTasks, ra.cancelTaskRoots)` |
| 任务来源（控制台 stdin 环） | `:260` `src := startResidentTaskSource(rt, ra)` |
| 事件环 | `:285` `rt.RunEventLoop()` |

**尺命令＋读数**：
`grep -rn --include=*.go 'startResidentPanel(' cmd internal tools | grep -v '_test\.go:'` → **2 命中**
（`panel_resident_windows.go:144` 定义 ＋ `resident_windows.go:157` 唯一调用）。
含 `_test.go` 的同尺 = **3**（正向对照：证明这把尺命中得了真名，不是因为文件不存在而空）。
`newResidentPanelManager(` 同形：非 test **2**（定义 `:228` ＋调用 `:151`）。

**追得到 main**：是。`main.go:66 runResident()` → `:151/:157` → `panel_resident_windows.go:154 reg.Spawn(..., rp.loop)` → `:299 w.Run()`。
**三档结论**：**〔已写零调用＝接线〕不适用——这一跳是已接线**。A 跳本身没有缺口。
**票**：票 33（片 A/宿主）＋票 228（球归常驻）＋票 246（门归常驻）。

---

### B. 面板宿主构造：谁 new 出 PanelManager，生产调用者几枚

**存在吗**：存在，**且非 test 调用者恰 1 枚**。

尺（锚在**调用形状**，带括号）：
`grep -rn --include=*.go 'NewPanelManager(' cmd internal tools | grep -v '_test\.go:'`
→ **2 命中**＝`cmd/wisp/panel_host_windows.go:213`（**定义**）＋ `cmd/wisp/panel_resident_windows.go:253`（**唯一的构造调用**）。
含 `_test.go` 的同尺 = **13**（正向对照；差值 11 枚是测试里的构造，不是产码调用者）。
⚠ 这条正是编排者警告的那枚尺坑：裸 grep 符号名会把定义点与同名类型当调用者，所以这里逐枚点名行号。

链上的装配内容（`newResidentPanelManager`，`panel_resident_windows.go:228-253`）：
`newResidentComposerDispatch` → `newComposerDispatchChain(dataDir, auditf, "resident-panel")`（`panel_inbound.go:228`）
→ `panel.BuiltinAssets()`（`:238`；失败 ⇒ `assets = nil`，只 Warn 不致命）
→ `NewPanelManager(disp, assets, dataPath, withGeometrySource(panelGeometrySource(dataDir)))`（`:253`）。

**这一跳真正的风险不在构造，在"有没有页面可发"**：
- `git ls-files frontend/dist` → **只有 `frontend/dist/.gitkeep`**；`git check-ignore -v frontend/dist/index.html` → **命中 `frontend/.gitignore:12` 的 `dist/*`**。
- 本机 `ls -la frontend/dist` → 磁盘上**有** `index.html`（1044 字节）＋ `assets/` 目录（**本腿按边界未读其内容**，只 `ls`）。
- ⇒ 干净检出跑出来的 exe 走的是 `serveNotBuiltNoticeLocked`（`panel_host_windows.go:442`，一发"bundle 没建"的自造页），
  本机现在跑出来的 exe 才可能有真页面。这正是**票 33 AC#12 未勾那一格**（"能构建不等于有页面可发"）。

**三档结论**：构造与装配 = **已接线**；"页面字节从哪来" = **〔写了一半〕**（产码侧齐，缺的是**未跟踪的构建产物**，不是代码）。
**票**：33 AC#11（已勾）／AC#12（未勾，归口＝界面侧那枚 agent 把 dist 填上）。

---

### C. STA 线程与泵：谁拥有线程、`Run()` 从哪儿被调、与 D38 十步有没有冲突

**存在吗**：存在，形状与编排者给的读数 #1/#2/#3 一致（本腿现量复核，非转述）。

| 事实 | 行 |
|---|---|
| 线程由注册表拥有（**不是裸 `go func`**：有 owner、有 recover 归属） | `panel_resident_windows.go:153-154` `observe.NewRoot("panel-host")` ＋ `reg.Spawn(panelSTAName, panelSTAOwner, root, rp.loop)` |
| 锁 OS 线程一辈子 | `:262` `runtime.LockOSThread()` |
| 显式要 STA(`0x2`)，非 S_OK/S_FALSE 就**只作废这一条线程** | `:263` `pnlCoInitializeEx.Call(0, coinitApartmentThreaded)`（常量 `:88` = `0x2`） |
| **`Run()` 的唯一生产调用点** | `:299` `w.Run()` |
| `Dispatch()` 的唯一生产调用点（post 路由） | `:360` `w.Dispatch(fn)` |

尺（`.Run()`）：`grep -rn --include=*.go '\.Run()' cmd internal tools | grep -v '_test\.go:'` → **2 命中**
＝`panel_resident_windows.go:299` ＋ `tools/d22scan/gitignore.go:378`（后者是 `cmd.Run()`，同名单但不同物，
具名排除）⇒ **本仓面板产码里 `Run()` 只有 1 枚调用者**。含 test 的同尺 = 24。

**库侧现量（读数 #2 的成因，本腿自己复跑，不看死腿 logs）**：
`grep -n 'dispatchq\|func (w \*webview) Run\|func (w \*webview) Dispatch\|WMApp' webview.go`
@ `github.com/jchv/go-webview2@v0.0.0-20260205173254-56598839c808/webview.go`
→ `:58 dispatchq []func()`／`:351 func (w *webview) Run()`／`:360-363` 只有 Run 的循环里取 `dispatchq`／`:443-447 Dispatch` 只 append ＋ `PostThreadMessageW`。
⇒ **"不调 `Run()` 队列永远不空"今天仍然为真**，且本仓已经按形 ⓐ 落地（`Run()` 在面板线程上占住那条线程）。

**与 D38 十步**：**不冲突，本腿判"不需要动十步"**。
- 十步序本身未被触碰：面板的 teardown 是 `resident_windows.go:158` 的一枚 defer（LIFO 最后跑），
  两处产码注释逐字声明同一件事——`cmd/wisp/panel_resident_windows.go:41-47`（"What this file does NOT do …
  its teardown is a defer in runResident, the same shape ticket 228 chose for the ball"）与
  `cmd/wisp/resident_windows.go:146-150`（"the frozen D38(e) ten-step order (internal/proc/shutdown.go) is not
  touched: **no new step, no new hook name on that closed roster**"）。
- 线程名 `panel-sta`（`:82`）**刻意不进** `observe.ResidentNames` 那六枚冻结名（加名＝D38b 契约变更），
  未列名按 `CategoryUnknown` 记录、不失败——这条是产码注释里自己写的，与本腿读到的形状一致。
- 唯一"注册进 D38(e)"的动作是 `resident_windows.go:240` 的 `RegisterShutdownHook(proc.StepCancelTasks, ...)`：
  那是**给已有第 3 步填生产者**，不是新增步。票 246 已裁，本腿不动。

**三档结论**：**已接线**。
**票**：33 片 A＋编排者裁定 P1/J1（不在票面上，票面只有 AC#1..AC#14 与 What to build 里那句"created on the shared `ui-sta`"——
⚠ **那句 Key constraints（`33-panel-host-c27.md:58`）与已落地的"专用线程"是矛盾的**，见 §3 冲突 #7）。

---

### D. Go→页面那一跳（推快照）：**今天一块没写**

**泵存在且会被驱动**（这部分是齐的）：

| 事实 | 行 |
|---|---|
| 泵在装配根建出 | `cmd/wisp/run.go:699` `rt.pump = panel.NewSnapshotPump(panel.PumpSources{...})` |
| 触发点 1（UI 侧） | `run.go:735` `rt.ui.publish = rt.publishPanelSnapshot` |
| 触发点 2（sink 侧） | `run.go:998` `consoleSink{... publish: rt.publishPanelSnapshot}` |
| 常驻腿也建了这套 | `resident_task_source_windows.go:278` `run, code := assembleRuntime(spec)` ⇒ 常驻进程**确实持有 rt.pump 并会 Publish** |
| Publish 的出口 | `panel_pump.go:401 publishPanelSnapshot` → `:405 rt.pump.Publish()` → 错误才会打印，成功**只进 `bookPanelSnapshot`** |
| `bookPanelSnapshot` 干了什么 | `panel_pump.go:319-327`：存进 `rt.lastSnap/lastSnapBytes` ＋ 落**一枚摘要日志行**（`panelSnapshotSummary` `:331`，含 bytes/sha256）。**没有任何一发出站投递** |

尺（决定性那把）：
```
grep -rn --include=*.go '\.Eval(' cmd internal tools | grep -v '_test\.go:'   → 0 命中
grep -rn --include=*.go 'PostScript\|InvokeScript'   cmd internal tools        → 0 命中
```
（正向对照：同尺含 `_test.go` 时 `\.Eval(` = **1 命中**，⇒ 尺命中得了真名，0 不是尺坏了。）
**⇒ 产码里今天没有任何"Go 主动往页面塞东西"的调用。**

尺（有没有那条边）：
```
grep -n 'panel\|Panel' cmd/wisp/resident_task_source_windows.go | grep -v ':\s*//'   → 0 命中
grep -n 'PanelManager\|residentPanel' cmd/wisp/run.go                               → 0 命中
```
**⇒ 任务管线（rt.pump 的宿主）与面板窗口（PanelManager 的宿主）之间今天不存在任何引用边。**

尺（AC#7 那枚前置条件）：
`grep -rn --include=*.go 'Marshal()' cmd internal tools | grep -v '_test\.go:'` → **1 命中**＝`internal/panel/pump.go:337`
**它自己的定义行**⇒ **`SnapshotPump.Marshal()` 至今零生产调用者**（`Publish()` 内部走 `pump.go:338` 的 `json.Marshal`，
不调这枚导出方法）。票 33 AC#7 写的触发条件"本票落地 `Marshal()` 的第一枚**生产**调用者之后"**还没到**。

**三档结论**：**〔写了一半〕偏"一块没写"**——包/序列化/摘要/计数全在跑，**投递那一条边整块缺失**。
**票**：票 145（快照载体；未勾 3 枚：`AC#2` 落地有源那一集／`AC#2b`／`AC#6` 反向判据）＋ 票 33 AC#7（未勾，前置未到）
＋ 票 35 的 `Resync`/`Backpressure` 两枚未勾 AC（`35-panel-bridge-c17.md:47`、`:49`）——它们都以"有出站通道"为前提。

---

### E. 页面→Go 那一跳（批准/拒绝入向）：门在、路由在、审批那一支没接；三枚处理器没写

**门（真存在，且是常驻腿唯一的入向门）**：
`panel_host_windows.go:386` `webview2.NewWithOptions(...)` 之后
`:396-405` 一带：hwnd 取回后
`:405` `bindErr := w.Bind(panelDispatchBinding, func(raw string) string { reply, _ := m.dispatchRaw(ctx, raw); return reply })`
（常量 `:80 panelDispatchBinding = "wispDispatch"`）→ `:630 dispatchRaw` → **`:637 return m.disp.Handle(ctx, raw)`**。

`Handle` 的生产调用者尺：
`grep -rn --include=*.go '\.Handle(' cmd internal tools | grep -v '_test\.go:'` → 22 命中，
**其中 `ComposerDispatch.Handle` 只有 2 枚**＝`panel_host_windows.go:637`（窗口那条门）＋ `panel_inbound.go:163`（CLI stdin 那条门）。
其余 20 枚全是同名噪声（`windows.Handle(...)` 类型转换、`logging.go`/`logsink.go` 的 `slog Handler.Handle`）——
⚠ 具名报回：这正是编排者警告的"裸 grep 符号名把同名 OS 类型当调用者"，实测在这把尺上又发生了一次。

**白名单（6 枚方法，`internal/panel/bridge.go:41-68`）**：
`panel.mode.request`／`panel.workspace.request`／`panel.attachment.add`／`panel.message.send`／`config.get`／`config.set`。
**里面没有任何一枚是审批方法**，而且这是刻意的：`bridge.go:53-66` 逐字写明"设置"要拼成 `config.set`
而不是 `config.approve`/`config.grant`，理由＝"AGENTS.md §1.2 keeps panel-side L2 'allow' off the table entirely"。

**链上挂了什么**（`newComposerDispatchChain`，`panel_inbound.go:228`，定义行；常驻与 CLI 共用同一枚函数）：

| 方法 | 处理器 | 常驻链上 |
|---|---|---|
| `panel.mode.request` | `panel.ModeWriteHandler.HandleModeRequest`（`composer_handlers.go:111`） | **已挂**（`panel_inbound.go:247` 构造，`Confirm: nil` 在 `:252`） |
| `config.get` / `config.set` | `panel.ConfigWriteHandler.HandleConfigRequest`（`config_handlers.go:269`） | **已挂**（`configWrites` 构造 `:266`，交进 `Config:` `:279`） |
| `panel.workspace.request` | — | **nil**（`:275` `Workspace: nil`） |
| `panel.attachment.add` | — | **nil**（`:276` `Attachment: nil`） |
| `panel.message.send` | — | **nil**（`:277` `Message: nil`） |

（同一枚链里 `perm.Options{Confirm: nil}` 在 `:242`，`ModeWriteHandler.Confirm: nil` 在 `:252`——两枚 nil 都属刻意，见 F 跳。）

尺（"接口写了、实现没写"那一族，逐枚点名，全部非 test）：
```
grep -rn --include=*.go 'HandleWorkspaceRequest'   cmd internal tools | grep -v _test → 2 命中（composer_dispatch.go:77 接口声明 / :186 switch 调用）
grep -rn --include=*.go 'HandleAttachmentRequest'  …                                    → 2 命中（:84 / :191）
grep -rn --include=*.go 'HandleMessageRequest'     …                                    → 2 命中（:89 / :196）
```
⇒ 这三枚**没有任何实现类型**（只有声明＋路由），到点按名拒（`ErrNoHandlerAttached`，`composer_dispatch.go:230`）。
正向对照＝同一把尺换 `HandleModeRequest` → **3 命中**（多出来那一枚就是实现行 `composer_handlers.go:111`），
证明"2 命中"确实意味着"没实现"，不是尺看不见实现。

**"写好但零（面板）调用"那一族＝审批答复侧**（这是本腿要逐枚点名的那一族）：
`internal/agent/approval/replies.go` 的 4 枚面板向方法，非 test 调用者**只有 `cmd/wisp/approval_reply.go` 一处的动词表**：

| 方法 | 定义 | 唯一非 test 调用者 |
|---|---|---|
| `Replies.PanelReject` | `replies.go:378` | `approval_reply.go:302`（`replySurface.panelReject`） |
| `Replies.PanelAllow` | `replies.go:431` | `approval_reply.go:331`（`replySurface.panelAllow`） |
| `Replies.Head` | `replies.go:473` | `approval_reply.go:375` |
| `Replies.View` | `replies.go:482` | `approval_reply.go:384` |

而这四枚的上游驱动者只有**两枚 stdin 环**，没有一枚来自面板：
`attachReplyListener(` 非 test 命中 2＝定义 `approval_reply.go:442` ＋唯一调用 `run.go:802`（`s.reply` ＝控制台 stdin）；
`runReplyLoop(` 命中 2＝定义 `approval_reply.go:530` ＋唯一调用 `approval_reply.go:456`；
常驻那枚是 `resident_task_source_windows.go:410` `src.surface.handle(verb, corr, arg)`，
上游 `runConsoleLoop`（定义 `:374`，唯一调用在 `:231` 之后），入口读者 `interactiveStdin()`（`approval_reply_stdin_windows.go:41`）。
**⇒ 面板那条 `wispDispatch` 门里，没有任何方法名能走到 `Panel*` 这一族；
而 `Panel*` 这一族今天只能靠"有人在控制台里打字"到达——双击图标的常驻进程里 `interactiveStdin()` 返回 nil（`resident_approval_windows.go:338` 那段注释自己点名"the double-click shape"）。**

**三档结论**：**〔写了一半〕**——门＋路由＋2/6 支已通；
（α）`workspace`/`attachment`/`message` 三支＝**处理器一块没写**（票 186/92/35 各自的地盘）；
（β）审批入向那一支＝**白名单里没有、听众里也没有**，属"一块没写"。
**票**：33 AC#9（**未勾**，票面 `:111` 逐字"**别把这一格当成入向已接线**：H2／H3／H10 仍要真宿主"）
＋ 248（未勾 5 枚，含 `AC#4` 真机那一发、`AC#9` "宿主在、内容不在"）
＋ 186（9 枚全未勾）／92（3 枚未勾）／35（6 枚全未勾，含 `Forged panel allow`、`Resync`）。

---

### F. 审批卡进常驻：门已按 246 注入，卡只到球，**到面板／被面板答复两块没写**

**"门由装配根注入"那一形＝已落地**（票 246 的裁定，本腿现量复核）：
`resident_windows.go:132` 建门 → `:217` 把 `ra.vetoByEsc` 交给球 → `:228` `ra.bindBallHost(rb)` →
`:240` 第 3 步钩子 `ra.cancelTaskRoots` → `:260` `startResidentTaskSource(rt, ra)`（`ra` 交进去当答复侧的门面）。
票 246 现量：**未勾只剩 1 枚**＝`AC#8`（"那枚 boot 时序 flake 要有归属"，修法在常驻启动序，`246-...md:61`）；已勾 8 枚。

**卡的显示面今天只有球**：
`resident_approval_windows.go:365` `ra.ui = &ballCardUI{ra: ra}` → `:370` 交进 `approval.Gate` 的 `UI:` 字段；
`ballCardUI` 的方法表：`:864 Prompt`／`:929 Update`／`:823 attachBall`／`:829 releaseBall`／`:947 settleOrb`。
**整个 988 行的 `resident_approval_windows.go` 里，`panel`/`Panel` 只命中 2 行，且都不是调用**：
`grep -n 'panel\|Panel' cmd/wisp/resident_approval_windows.go` → 8 命中，剔掉注释后只剩
`:76 residentPanelSource = "cmd-wisp-resident-panel-route"`（一枚路由名字符串）＋
`:391 PanelSource: residentPanelSource`（`HostBinding` 的一个字段）。
⇒ **没有任何一行把卡片推给面板窗口。**

**`ra.cards` 上被常驻腿用到的方法**（`grep -n 'ra\.cards\.' cmd/wisp/resident_approval_windows.go`，9 命中全列）：
`:387 Attach`、`:605 AwaitingHuman`、`:612 Veto`、`:725 Pending`、`:728 Reject`、`:743 Forget`、`:875 Record`、`:932 Forget`、`:948 AwaitingHuman`。
**`PanelReject`/`PanelAllow`/`Head`/`View` 在这九枚里一枚都不出现**（与 E 跳的尺一致：它们只被 `approval_reply.go` 的动词表调，
而动词表只被 stdin 喂）。

**面板链里的"卡"是刻意没有的**（不是漏）：
`panel_inbound.go:241` `Confirm: nil`（perm.Options）＋ `:250` `Confirm: nil`（ModeWriteHandler），
`panel_resident_windows.go:221-224` 逐字 "**Confirm stays nil on purpose: this leg has no card of its own**, so a widening request is refused before Store.Set is ever reached"。

**缺的是哪一块（只报缺跳，不选形）**：
1. 出向：面板窗口里没有任何读者会去拿 `ra.cards.Pending()` 或 `rt.pump` 的包 ⇒ **"卡片出现在面板"这块没写**（与 D 跳同一枚缺口的审批子集）。
2. 入向：`wispDispatch` 门后面的 6 枚方法名里**没有一枚通往 `Panel*` 答复族** ⇒ **"面板能答复审批"这块没写**。
   ⛔ 本腿**不**提"给白名单加一枚方法／给 C17 加一个方法"这类形状——那是契约变更、要人工批准；本腿只报"这一跳缺"。
3. 装配根注入那一形本身**不缺**，别把它当缺口。

**三档结论**：**〔写了一半〕**（门已进常驻且能否决／能超时拒绝；卡→面板＝没写；面板→答复＝没写）。
**票**：246（AC#8 未勾，且它管的是 boot 时序 flake 不是面板）＋ 33（AC#9/AC#14 未勾）＋ 145（快照有源）
＋ 167（`AC#4` "停止必须在助手正忙时也能到达"等 6 枚未勾——停止那枚与本跳同一条边）。
⚠ **具名报回**：编排者转述"票 246 已裁'门由装配根注入'，理由＝不要出现第二枚『谁在等批准』的真相源"——
票面／产码两处理由一致（`resident_windows.go:113-118` 逐字点名 "the second truth source 246-a1 §5 reason 1 refuses"），
**该裁定没有过期**。

---

## §2 接线顺序建议（先补哪一跳最省返工）

**第 1 跳补 D（Go→页面的投递那一条边），不要先补 E/F。**

1. **D 是唯一的"第一因"**：A/B/C 三条已经成链（点热键→建窗→`Run()` 泵），机主"我点了，界面有没有出来"
   今天只差两件事：(a) **bundle 真在不在**（B 跳，缺的是**未跟踪产物**不是代码，归界面侧那枚 agent）；
   (b) **页面有没有内容可渲染**（D 跳，缺一整条边）。
   先做 D 不需要动线程归属（读数 #1 那条禁令一寸都不用碰），不需要动 D38 十步（C 跳已核），
   也不需要碰任何契约面（C17/白名单一字不改就能把包推给已存在的窗口）。
2. **第 2 跳补 F 的出向那一半**（把 `ra.cards`/快照的读者接到已存在的窗口上）：
   票 246 的门已经在那儿了，`Pending()`/`AwaitingHuman()` 已经有生产调用者（`:725`/`:605`），
   加一个"给窗口读"的读者**不新增第二枚真相源**——这正是 246 裁定愿意放的形状。
3. **第 3 跳才是 E 的入向答复那一条边**：顺序反过来会造出一枚"能点但看不见最新状态"的门，
   那一形的返工是把整条 UI 判定重写，不是补一行。
4. **E 的三支 nil 处理器（workspace/attachment/message）排最后**：它们各有票（186/92/35），
   跟"面板开不出来"不是同一枚第一因，且 `attachment`/`message` 一旦落地会牵动 D 的包形状（票 145 的未勾 AC 就是这一集），
   先动它们会让 D 白做一遍。

⛔ 建议里**不含**"新增契约/C17 加方法/改白名单"任何一项——那些要人工批准，本腿只报缺跳。

---

## §3 与 D38 十步／票 246 装配根注入那条的冲突点（逐枚具名）

**先说结论：本腿没有读到任何一条"必须动 D38 十步序"才能接线的需求。**
真正要报的是**编排者给的三条读数已过期/需更正**，逐枚：

1. **冲突 #1（重要，直接推翻开场前提）**：开场那句"**面板宿主（WebView2/PanelManager）那套代码已经在仓里，但没有接进常驻腿** ⇒ 用户双击图标后『面板开不出来』"——
   **已过期**。尺与行号：`NewPanelManager(` 非 test **1 枚调用**＝`cmd/wisp/panel_resident_windows.go:253`，
   其上游 `cmd/wisp/resident_windows.go:151`→`:157`，热键/托盘经 `resident_windows.go:218 panel.RequestToggle` 到达；
   `panel_resident_windows.go:141-143` 甚至自己写着"Building the manager is the **first non-test call** of NewPanelManager in this repository"。
   ⇒ "开不出来"今天的原因**不是没接线**，而是 B/D 那两格（bundle 未跟踪、快照无边）。**按开场那句派单会做错活。**
2. **冲突 #2**：票 33 **AC#14**（`:42-43`）里那句"现读：产码 `cmd/wisp/panel_host_windows.go:222` 已经在用 `w.Eval(...)`"——
   **今天 0 命中**（`grep -rn '\.Eval(' cmd internal tools | grep -v _test.go` ⇒ 0；含 test ⇒ 1，尺非坏）。
   且行号漂了：`bringUp` 现在 `panel_host_windows.go:318`；AC#14 列的三形里，**形 ⓐ（换成库的 `Run()`）已落地**，
   形 ⓑ（一律走 `Eval`）**没有**。**AC#14 本体仍未勾**（它的完成判据"页面 `await` 绑定并断言回话真到达"没有会响的用例——见 §4）。
3. **冲突 #3**：`cmd/wisp/panel_resident_windows.go:171` 那行注释指"the panel host is built at `cmd/wisp/resident_windows.go:142`"——
   现量 **`:151`**。同类过期指针还有 AC#13 引的 `:221/:227/:250/:372-375`，本腿现量对应
   **`serveEntry` 调用 `:428`／`firstRoundTripLocked` 调用 `:426`／`serveEntry` 定义 `:454`（入口字节那发 `w.SetHtml(string(data))` 在 `:468`）／探测页那发 `w.SetHtml` 在 `:681`**。
   （次序结论见冲突 #4，行号只是漂了，不是本腿新加的洞。）
4. **冲突 #4（不是矛盾，但要报）**：**票 33 AC#13**（探测页盖掉真页面）**未勾**，但**代码次序已经改了**：
   `panel_host_windows.go:420-425` 的注释逐字 "**AC#13's order, and it is the whole fix**: prove the message channel first, hand the page over LAST"，
   实测次序＝`firstRoundTripLocked`(`:426`) 之后才 `serveEntry`(`:428`)。
   ⇒ 这格现在**欠的是非实现者的表**，不是欠代码；本腿不翻框（边界）。
5. **冲突 #5（读数 #2 与仓里一发自泵并存，需有人裁）**：编排者说"自泵形三枚全 `NOT_RESOLVED`"，
   而**产码里至今留着一发自泵**：`firstRoundTripLocked`（定义 `:656`）bind 一枚 `wispProbeRT`（`:668`）→ 自造探测页 `SetHtml`(`:681`)
   → 循环调 `pnlPumpOnce()`（调用 `:686`，定义 `:707`＝`PeekMessageW(PM_REMOVE)`＋`TranslateMessage`＋`DispatchMessageW`）。
   库尺已证 `dispatchq` 的唯一取出点在 `Run()` 内（`webview.go:351/362-363`）⇒ **这发探测在 `Run()` 形下是否仍然可靠，本腿判不动**（只能真机量，未批）。
   ⚠ 我**不**用它否掉读数 #2；我只报"读数 #2 说的『自泵不兑现』与仓里仍在跑的这发 `pnlPumpOnce` 探测并存"，裁不裁由编排者。
6. **冲突 #6（小，登记性质）**：`go.mod:19` 仍把 `github.com/jchv/go-webview2` 标 `// indirect`，
   而 `cmd/wisp/panel_host_windows.go`／`panel_resident_windows.go` 都**直接** import 它（产码 2 枚 import）。
   这是陈旧标记、不影响构建；归构建链那侧（111-r5 在跑），本腿不动。
7. **冲突 #7（票面内部矛盾，属人工批准域，只报不动）**：票 33 `Key constraints`（`33-panel-host-c27.md:58`）仍写
   "WebView2 created on the **shared `ui-sta`** STA thread (D38a)"，
   而**已落地的形状是"专用 STA 线程"**（`panel_resident_windows.go:5-19` 头部注释点名裁定 P1/J1，且实测反证写在同一处：投进 `ui-sta` 会让外层泵的迭代计数冻住＋嵌套泵提前走 5 枚任务）。
   ⇒ 编排者读数 #1"不许改线程归属"与**票面那句**互斥，与**产码**一致。票面文字要不要改属人工批准，本腿一字未动。

**与 D38 十步的正面核对结论**：无冲突。十步序未被触碰（`internal/proc/shutdown.go` 未被本腿读为"待改"，
`resident_windows.go:75-85` 的退出 defer 仍走 `rt.Shutdown(false)` 并逐失败打印步号）；
面板只有 defer 形态的 teardown（`:158`）与一枚已有步的钩子（`:240`）；`panel-sta` 刻意不挤进 `observe.ResidentNames`（D38b）。

---

## §4 判不动（只能真机量／机主未批，一律未跑）

**闸门**：`winlive` 类（真开窗口）测试**机主至今没批准跑**，且**CI 里没有 tag 档**。
本腿具名声明：**没开过任何真窗口、没给任何 yaml 加 `-tags winlive`、没跑任何 `go test`。**

尺（现量，非转述）：
- `grep -rln 'go:build.*winlive' --include=*_test.go cmd internal | wc -l` = **11**（带 tag 的真窗口测试文件数）
- `grep -rn 'winlive' .github/workflows/*.yml` = **0 命中**（⇒ 读数 #4 仍为真：CI 里没有这一半）

下面每一格都是**〔仅本机可量，未批〕**，本腿没有读数，也不许别人拿"本腿没报"当"它通过"：

| 判不动的那件事 | 为什么只能真机量 | 归哪格 |
|---|---|---|
| 冷开窗 ≤1500ms / 热开 ≤200ms（10 发 P50/P95） | 要真建窗 | 票 33 AC（`:74` 未勾） |
| **第一次真开面板到底成不成**（`showOnThread` 那发 `Println`，`panel_resident_windows.go:411`） | 要真建窗 | 票 33 片 A |
| `bringUp` 的前置拒绝（`errPanelRefusedThread`，返回点在 `:381-383`，紧随 `drainStaleQuitBeforeCreate()`:381）在真实拓扑里会不会被踩到 | 要真消息队列 | 票 33 AC#14 邻域 |
| **页面 `await` 绑定回话真到达**（AC#14 的完成判据：一枚会响的用例） | 要真页面＋真泵 | 票 33 **AC#14**（未勾） |
| 那发 `pnlPumpOnce` 自泵探测在 `Run()` 形下是否仍可靠（§3 冲突 #5） | 同上 | 票 33 **AC#13**（未勾） |
| Show/hide/destroy 生命周期＋进程树子进程数稳定／≤2s 退出 | 要真窗口＋Job Object | 票 33（`:72` 未勾） |
| 焦点回环（编辑器→面板→隐藏→焦点回编辑器） | 要真前台窗 | 票 33（`:76` 未勾） |
| WebView2 运行时缺失回退（改名/mask loader） | 要真加载器 | 票 33（`:77` 未勾） |
| CSP 在真响应上存在、`connect-src 'none'` | 要真响应检查 | 票 33（`:78` 未勾） |
| "无监听端口"那一半（`GetExtendedTcpTable` 且**必须过滤 `dwState==LISTEN`**） | 要本机网络表；票面 `:30` J7 已裁"只留本机证据并登记 CI 里没有" | 票 33 |
| MTA 当场 panic 的**成因**（HRESULT 读不到，要 unexported handler） | 要真窗＋改库可见性 | 票 33 头注释 `:36-39`（产码自己写的"cause still undetermined"） |
| 246 的 boot 时序 flake 归属 | 要真起常驻 | 票 246 **AC#8**（唯一未勾） |

**本腿另外两格"判不动但不是真机原因"**：
1. **`frontend/dist` 里的页面内容真不真**：边界令 `frontend/**` 连读都不许。本腿只量了
   "tracked 只有 `.gitkeep`"＋"`frontend/.gitignore:12` = `dist/*`"＋`ls`（`index.html` 1044 字节、`assets/` 存在，**未读内容、未打开**）。
   ⇒ 归口票 33 AC#12 ＋ 界面侧那枚 agent。
2. **票与产码"谁该改文字"**：§3 冲突 #3/#5/#7 三处都指向票面/注释过期。改票面与改产码都超出只读腿权限，
   本腿只具名，不动一字。
