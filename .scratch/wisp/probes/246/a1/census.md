# 票 246 — 前段只读普查 `246-a1`：卡片／审批门怎么进常驻进程（甲形 vs 乙形代价表）

> 本文件是 AC#0 的交件：**只量代价、不写产码、不改任何 `- [ ]` 框**。
> 两形各一行（甲＝常驻腿自己起 loop／审批门；乙＝装配根 `cmd/wisp` 注入），逐形给 ①-⑥ 六格。
> ⛔ `frontend/**`／`design/**` 零读零写零转述；`docs/PLAN.md`／`docs/specs/**`／冻结件一字不动。

---

## §0 起手名册与复认凭据

### §0.1 起手名册（闸门＝终态等于此刻，非"必须为空"）

- 取数时刻：`2026-09-30 18:29:07 +0800`
- 起手锚点（`git log --oneline -1`）：`31fa4d23 ledger(A480 收 245-v1 …)`
- 分支：`dev`
- 起手名册全量已落盘为可 diff 基线：`.scratch/wisp/probes/246/a1/roster-start.txt`（`git status --porcelain` 原文，取数即存）
- `git status --porcelain` 行数：**176**（取数 18:29；落 roster 文件时已含本腿刚提交的骨架 commit，故 roster-start.txt 计 178＝原 176 ＋ 本腿 `msg-skeleton.txt` ＋ `roster-start.txt` 自身两枚未跟踪件）
- **本腿负责的写面**：只允许 `.scratch/wisp/probes/246/a1/**`（本文件＋roster＋msg 件）＋票面 Progress log 末尾一行。其余一律别人的脏件。
- `internal`／`cmd` 目录脏件数（起手现跑）：**0**（`git status --porcelain -- internal cmd` 为空）——即此刻产码树对普查腿是干净的，但本腿一枚产码都不碰。
- **终态闸门算法**：`diff <(git status --porcelain) .scratch/wisp/probes/246/a1/roster-start.txt` 的差集只允许落在 `.scratch/wisp/probes/246/a1/` 内；`internal`／`cmd`／`docs/PLAN.md`／`docs/specs/**`／冻结件必须逐字节等于起手（本腿没碰过＝天然成立，收尾再复跑一把坐实）。

### §0.2 三条现读凭据的独立复认（编排者给过，我不照抄、现跑）

| # | 编排者给的凭据 | 我的独立尺（命令） | 我的现读结果 |
|---|---|---|---|
| 1 | `resident_ball_windows.go:22/:69-70/:126/:163` 逐字自陈"没有管线/麦克风/审批门/面板宿主" | `grep -nE "no task pipeline\|no microphone\|no approval gate\|no panel host\|has no executor here" cmd/wisp/resident_ball_windows.go` | **复认成立**：`:22` 「pipeline, no microphone, no approval gate and no panel host today」；`:70` 「this leg has no card path (no approval gate, no state machine)」；`:126` 「this leg has no task pipeline, no microphone and no approval gate」；`:163` 「…no approval gate and no panel host, so the gesture has no executor here」 |
| 2 | `TakeEscForCancel`／`ReleaseEscAfterSession` 非测试消费者只有 `cmd/balldebug/main.go:635/:638` | `grep -rnE "TakeEscForCancel\|ReleaseEscAfterSession" --include=*.go internal cmd \| grep -v _test` | **复认成立**：非注释命中里，**函数定义**在 `internal/ball/ball_windows.go:873/:895`，其余 `:247/:810`、`hotkey_windows.go:254/:256` 全是注释；**唯一生产调用者**＝旁支调试程序 `cmd/balldebug/main.go:635/:638`。产品进程（`cmd/wisp`）零调用者 |
| 3 | 契约要求四条否决通道：`PLAN.md:3082`（D43，冻结）＋ `approval.go:90` | `sed -n '3082p' docs/PLAN.md` ＋ `sed -n '88,93p' internal/agent/approval/approval.go` | **复认成立**：`PLAN.md:3082` 转移表第 22 行逐字「`Confirming`(L1) \| 否决（**单击球 / `Esc` / KWS 否决词 / 面板拒绝**，B1）」；`approval.go:89-92` `channelNames` 四枚俱在（`ChannelBall: 单击悬浮球`／`ChannelEsc: 按 Esc 键`／`ChannelPanel: 面板拒绝`／`ChannelKWS: 说取消词`）。⛔ 两处只读、一字不改 |

---

## §1 现状：常驻腿到底有什么、缺什么（现读）

- 常驻那条腿的进程拓扑（`cmd/wisp` 无参分支起什么）：`main.go:64 runResident()` → `resident_windows.go:27 runResident` → `proc.Boot(env)`（`internal/proc/boot_windows.go:57`：env→layout→Job Object→单实例→registry）→ `:60 installLogSink` → `:118 rb := startResidentBall(rt.Registry)`（`internal/ball` 建窗口/托盘/三枚热键）→ `:132 rt.RunEventLoop()`（空事件循环，只等 Ctrl+C/SIGTERM）→ 退出经 `:70 rt.Shutdown(false)` 走 D38(e)。**这条路里没有 agent loop、没有审批门、没有卡片显示、没有麦克风、没有面板宿主**（凭据 1 逐字自陈）。
- `wisp run` 那条腿装配了什么（`assembleRuntime` 全链，`cmd/wisp/run.go:336`）：`secret.NewStore → config.NewManager → llm.NewResolver/BuildEndpointProvider → memory.Open → session.Mint/Ledger → tools.NewPathCanonicalizer → tools.NewRegistry(BuiltinFSEntries/TaskEntries/SubagentEntries) → approval.New(Options{UI:consoleApprovalUI, Channels:NewChannels()空, …})（run.go:523）→ perm.New → panel.NewSnapshotPump → tools.New(gate,cancel,modes,grants) → agent.New(opts)`，`execute()` 里 `loop.RunAsync` 起跑。**这条路有 loop＋gate＋卡片显示＋答复监听（stdin reply），但没有球、没有托盘、没有热键、没有 D38(e) 退出序列**（`wisp run` 是一次性进程，任务跑完就退，不调 `proc.Boot`/`Shutdown`）。
- **两条腿是错开的**：会跑任务、会弹卡片、能否决（若有人按）的是 `wisp run`；带球、带托盘、带 Esc 借用位、走 D38(e) 的是无参常驻腿。票 246 的反向缺口＝**这两件事不在同一个进程里**，于是常驻腿里球进来了却没有东西可否决；`wisp run` 里卡片能进能答却没有任何一票真能落到球/Esc（run.go:501「`NewChannels()` stays empty and `runSpec.replyVeto` stays unset」）。
- 审批门 / agent loop / 卡片显示 / 否决通道在两条腿里的"有/无"矩阵：

| 能力 | `wisp run`（一次性 CLI 腿） | 无参常驻腿（带球的腿） |
|---|---|---|
| agent loop | 有（run.go:896 `agent.New`） | **无** |
| approval.Gate | 有（run.go:523 `approval.New`） | **无** |
| 卡片显示 | 有（`consoleApprovalUI` 打到 stdout） | **无**（无 stdout 消费面、无 WebView 宿主） |
| 答复（allow/reject） | 有（`attachReplyListener` 读 stdin） | **无**（无 stdin 面） |
| 否决通道 loaded | **空集**（run.go:525 `NewChannels()` 无参） | **无门可载** |
| 球/托盘/热键 | **无** | **有**（`startResidentBall`，热键 live 3/4） |
| Esc 借用位 | **无球可借** | **有机制、零调用者**（凭据 2） |
| D38(e) 十步退出 | **无**（不用 proc） | **有**（`rt.Shutdown`，但见下） |

- 四条否决通道各自的执行者现状（球单击 / Esc / KWS / 面板拒绝）：
  - **球单击 `ChannelBall`**：`resident_ball_windows.go:89 OnClickBall: recordBallGesture("click")`——只记账、说"没有执行者"。执行者＝**无**。
  - **Esc `ChannelEsc`**：借用机器 `TakeEscForCancel`/`ReleaseEscAfterSession` 在 `internal/ball` 里造好了，生产里只有 `cmd/balldebug` 调；常驻腿既没卡片也没调它。执行者＝**无**。
  - **KWS `ChannelKWS`**：`approval.go:104` 归 DEFERRED(票 41)，且常驻腿无麦克风。执行者＝**无**。
  - **面板拒绝 `ChannelPanel`**：`approval.go:107` 归票 37，且本树无 WebView2 宿主（`approval_always.go:165` 「links no WebView2 host」）。执行者＝**无**。
  - 结论：**`Replies.Veto`／`Gate.Veto` 这套否决路由在 `internal/agent/approval/replies.go:455`／`gate.go:415` 是齐的、可跑的**（票 87/201 已把它接进 `wisp run` 的 stdin `veto <编号>`），缺的**不是机制，是"带球那个进程里没人把一次球手势/Esc 按成一次 `Veto`"**。

---

## §2 甲形：常驻腿自己起一条 agent loop ＋ 一个审批门

> 一句话定义：常驻那条腿（`resident_windows.go` 无参路径）在**本进程内**从零装配一条 agent loop 与一枚 approval.Gate，
> 卡片、状态机、Esc 借用都挂在这条腿自己身上。

### ① 要动哪几枚文件（逐枚 `file:line` 现读）

- `cmd/wisp/resident_windows.go:118` —— 在 `rb := startResidentBall(rt.Registry)` 前后**新增一段自装配**：常驻腿自己 `approval.New(...)` 造 gate、造一枚球驱动的 `approval.UI`（新类型）、把 `ChannelEsc`/`ChannelBall` 标记 loaded。等于把 `run.go:523` 那段 `approval.New` 在常驻文件里**再写一遍**。
- `cmd/wisp/resident_ball_windows.go:89`（`OnClickBall`）／`:92`（`OnCancelHotkey`）—— 从今天的 `recordBallGesture(...)` 改成**回调进本腿自持的 gate**（`gate.Veto(Veto{Corr, ChannelEsc})`）。
- **新增**一枚球驱动 UI（例：`cmd/wisp/resident_ball_ui_windows.go`）—— 实现 `approval/ui.go:84 UI` 接口，`Prompt` 里调 `rb.b.SetState(Confirming)`＋`rb.b.TakeEscForCancel()`、`Update` 里 `ReleaseEscAfterSession()`。这是把 `run.go:1174 consoleApprovalUI` 的 stdout 版换成球版。
- `cmd/wisp/resident_windows.go:69-79`（退出 defer）—— **必须处理 D38(e) step 3**：见 ③。要么在此扩写 cancel 委托，要么改 `internal/proc/boot_windows.go:151 Shutdown(fast bool)` 让它接受调用方递来的 hooks（`ShutdownHooks{}` 现仅填 `CloseJob`）。

### ② 新增哪几条包级依赖边（现跑 import 清单证，不读注释）

现跑尺：`go list -f '{{join .Imports "\n"}}' ./cmd/wisp`（GOOS=windows）⇒ `cmd/wisp` **已直接 import** `internal/agent`、`internal/agent/approval`、`internal/ball`、`internal/proc`、`internal/statemachine`。`internal/proc` 只 import `buildinfo`+`observe`；`internal/ball` 只 import `observe`+`statemachine`。

- **合规做法（装配仍留在 `cmd/wisp`）＝新增 0 条包级边**。因为 loop/gate/UI/veto 全在 cmd/wisp 里组装，球只调自己包内公开的 `SetState`/`TakeEscForCancel`。
- **危险做法（本普查具名提示，不建议）**：若为了"让 loop 生命周期挂进 D38(e)"而让 `internal/proc` 去认识 scheduler／task-root ⇒ 会新增 `internal/proc → internal/agent`（或 `→ internal/agent/approval`）**正向边**。这正是票 238 第 1 刀「两枚正向依赖边一律不开」明令挡掉的形状，**甲形尤其容易被推向这一支**（因为甲要"常驻腿自己持有 loop 的生命周期"，最省事的想法是把生命周期交给 proc）。

### ③ 会不会破 `internal/proc` 的 D38(e) 十步退出顺序

现读：`internal/proc/boot_windows.go:151-161` 的 `Shutdown` **每次 `ShutdownHooks{}` 就地新建、只填 `hooks.CloseJob`（step 9）**；全仓 `grep` 证 step 1-7 的 hook **零生产者**（`SchedulerClose/StopHotkeyAndKWS/CancelTasks/StopAudio/ReleaseSpeechSessions/DestroyPanel/FlushAndCloseDB` 全 nil）。

- **卡片挂着时收到退出信号走到哪一步**：`RunEventLoop`（boot_windows.go:139）`sigCtx.Done()` → 返回 "signal" → `resident_windows.go` 的 defer LIFO：`rb.stop()`（球/热键先收）→ `rt.Shutdown(false)` 走 10 步，但 **step 3 `cancel-task-roots` 被记为 skipped** → step 9 CloseJob → 退出。
- **有没有人负责把待决卡片判成拒绝**：**没有**。`gate.PendingWindow/PendingApproval` 只有两条出口：任务 ctx 被取消（`ctx.Done()` → `AnswerReject`）或 C18 deadline 到。若 loop 的任务 Root 没人取消（step 3 空转），卡片所在的 goroutine 会在进程退出时**被 OS 掀掉、而不是被判成拒绝**——台账/审计里看不到 `ANSWER-REJECT`，只有 skipped 的 step 3。这跟 `wisp run` 里 `agentRuntime.close()`（run.go:782）主动 `replyRoot.Cancel()/reloadRoot.Cancel()` 把待决卡判成 deny 是**相反的姿态**。
- **甲形怎么补**：甲既然自持 loop，就得在常驻 defer 里**显式 cancel 那枚 loop Root**。但若不把 cancel 做成 proc 的 step 3 hook，**D38(e) 审计就会说"cancel-task-roots 跳过了"而实际取消发生在序列之外**——审计面与行为面分叉，是这型最容易被忽略的代价。要审计诚实，得给 `proc.Shutdown` 一个"调用方递 hooks"的接缝（见 §7 归口：proc 侧的注册机制本就不存在，非本票独有）。

### ④ 与票 228 后续片的先后关系

票 228 后续片（A477/A480 编队里排在 `244 → 票 228 后续`）在**同两枚文件**上还有账：常驻腿整个 `config.toml` 没接、托盘「退出」无执行者（`resident_ball_windows.go:192 recordTrayExit`）、球位置不持久（`:179`）、`ball.New` 失败路径不清理（`:101-106`）。

- **谁先做会弄脏谁**：甲形大改 `resident_windows.go`＋`resident_ball_windows.go`（把 `recordBallGesture`/`recordTrayExit` 这些"只记账"回调改成"真执行"）⇒ **会把 228 后续片"托盘退出无执行者""手势无执行者"这两格读数直接改掉**，228 后续片若按现状写判据会对不上。反过来，228 的 config 接线若先落地，会引入常驻腿读 `config.toml` 的新流程，甲形的 loop 装配点会跟着挪。**两票在 `cmd/wisp/resident_*.go` 上重叠 ⇒ 按 A480 编队"一律串行"，不能并发**。
- 但**本票最短链不依赖 config**（见 ⑥）：只做 gate+卡片+Esc 否决，不建 provider，就不碰 228 的 config 接线那格 ⇒ 串行下可先做最短链、把 config 依赖推给"真跑任务"那一支。

### ⑤ 需要 owner 点头的有几枚、各是"功能"还是"契约"

- **功能级（不是契约、不必先落批准记录，但值得告知 owner）**：
  1. 常驻进程里"起一条 agent loop"——D38b 名册 `internal/observe/goroutine.go:50` 已把 `agent-task-`/`tool-exec-`/`approval-waiter` 列为 **per-task 条目**，起 loop 用的是既有 per-task 槽、**不破 resident baseline=6** ⇒ **非契约**。（⚠ "改契约"永不构成"不做"的理由；这里连契约都不用改。）
  2. 球驱动 UI（卡片→Confirming→借 Esc）——是 `approval.UI` 接口的又一实现 ⇒ 功能。
  3. 把 `ChannelEsc`/`ChannelBall` 标 loaded——用现成 `approval.DefaultChannels()`/`SetLoaded` ⇒ 功能。
- **契约触碰（须先落一条批准记录，但仍要去做）**：
  4. **D38(e) 十步的 step 3 语义**：顺序本身冻结（`shutdown.go:11-31`），但"待决卡片在退出时被谁、在哪一步判成拒绝"目前无实现者。要把它落进 step 3，需要给 `proc.Shutdown` 加"调用方注册 hooks"的接缝——**顺序不动、只是补上注册机制**，属**实现缺口**而非改契约（`shutdown.go:77` 注释本就写"Later tickets register their hooks"）。⇒ 归口见 §7，若 owner 认为"给 proc 加注册入口"算契约面，则落批准记录后照做。
  5. **不新造第五枚否决通道、不改 C12/D43 形状**——本票只实现既有四通道，**零契约改动**（票面禁区已钉死）。
- 汇总甲形：**硬契约改动 ≈ 0**（D43/C12 一字不动），需告知的功能 ≈ 3 枚，需 owner 确认"注册接缝是否算契约"的灰区 ≈ 1 枚。⛔ 不因"要动 proc"就把甲判成"建议不做"。

### ⑥ 最小可跑通的那一发

**最短链＝『卡片进得来＋Esc 真能否决一次』，把面板与麦克风留给后面**——而且**这一发根本不需要 agent loop／provider／config／麦克风**：
- `run.go:732 confirmModeSwitch` 与 `approval_always.go:112 runWidening` 已证明：**不经过模型调用**也能用 `gate.AdmitTextTask(hostTaskID)` + `gate.PendingWindow/PendingApproval` 造出一张**真·待决卡片**。
- 于是常驻腿最短链：装配一枚 gate（球驱动 UI）→ `Channels().SetLoaded(ChannelEsc, true)` → 在一个手势（如球单击或一次测试钩子）上 host-initiated 造一张卡 → `OnCancelHotkey` 里 `gate.Veto(Veto{Corr: liveCards.AwaitingHuman().CorrelationID, Channel: ChannelEsc})`。`replies.go:252 AwaitingHuman`＋`replies.go:455 Veto`＋`gate.go:415 Veto` 这套路由现成可用（票 87 已把"命名卡片上的 Esc"接到 `q.reject`）。
- **判据形状**：造一张卡 → 按 Esc → 卡片应以 `AnswerVeto` 收口、球回 Sleeping、`ReleaseEscAfterSession` 归还 Esc、审计落 `ANSWER-VETO … channel=esc`（`gate.go:314`）。这条链把面板（票 37）、KWS（票 41）、真任务源（语音链/票 228 config）全部**留在射程外**，只补"执行者"这一枚反向缺口。

---

## §3 乙形：由装配根 `cmd/wisp` 注入（照票 238 第 1 刀"注入不开边"＋票 197 段"装配根是唯一接缝"）

> 一句话定义：常驻腿不自己装配；loop／审批门由 `cmd/wisp` 这个装配根的**既有唯一接缝**造好后注入常驻腿，
> 常驻腿只负责把球手势与 Esc 借用接到那枚被注入的 gate 上，并把通道标记 loaded。

### ① 要动哪几枚文件（逐枚 `file:line` 现读）

- `cmd/wisp/run.go:336 assembleRuntime` —— 把"造 gate + 造 approval.UI + 标 loaded + 持 liveCards"这一段从"只为 `wisp run` 的一次性 stdout UI"**抽成可复用**：产出一枚"带 veto 路由的 gate 句柄"，供 run 腿和常驻腿**共用同一个装配点**。consoleApprovalUI（`run.go:1174`）与新的球驱动 UI 都作为 `approval.UI` 的不同实现注入。
- `cmd/wisp/resident_windows.go:118` —— `runResident` 调装配根拿到的 gate 句柄，**把它作为参数/闭包递给** `startResidentBall`（照 `resident_windows.go:118 startResidentBall(rt.Registry)` 已在传的注入姿势，也照票 238 第 1 刀「`proc.BootOption` 的 `WithLayout` 同形递函数」）。
- `cmd/wisp/resident_ball_windows.go:89/:92` —— `OnClickBall`/`OnCancelHotkey` 从 `recordBallGesture` 改成调用**被注入进来的 veto 回调**（`func(channel approval.Channel) { injectedGate.Veto(...) }`）。
- **新增**一枚球驱动 `approval.UI`（同甲形，但由装配根建、注入常驻腿，常驻文件不自己 new gate）。
- `internal/proc/boot_windows.go:151-161` —— 同甲形，step 3 的账还在；乙形更倾向把"取消 loop Root"做成**装配根登记的 hook**（见 §3.1），proc 侧需要一个"注册 hooks"的入口。

### ② 新增哪几条包级依赖边（现跑 import 清单证，不读注释）

- **0 条**（与甲形合规做法同）。注入是**函数值/句柄在 cmd/wisp 内部传递**，跨包 import 图不动：`cmd/wisp` 已 import 全部相关包；`internal/ball`、`internal/proc` 保持叶子。
- 与票 238「两枚正向边一律不开、改注入」、票 197 段「装配根 cmd/wisp 是唯一的接缝」（台账 :8786-8792）**同向**：乙形就是把那条裁定用在"卡片/veto 怎么进常驻腿"上。⚠ 沿用 A464 措辞纪律——把票 197/238 的裁定推广到本票要**具名写成"这是推广"**，不写成"仓里本来就有这条规矩"。

### ③ 会不会破 `internal/proc` 的 D38(e) 十步退出顺序

- 事实与甲形**完全相同**（step 1-7 零 hook 生产者、`Shutdown` 只填 `CloseJob`，见 §2③）：这是 **proc 侧的公共缺口，不是甲/乙任选其一才有的**。
- 乙形的差别只在**谁来补 step 3**：乙把"取消 loop Root / 判待决卡为 reject"做成装配根**登记的 hook**（`hooks.CancelTasks = func(ctx){ loopRoot.Cancel(); 等 <=3s }`），于是审计里 `cancel-task-roots` 是**真执行**而非 skipped，卡片以 `AnswerReject` 有据可查地收口。前提仍是 `proc.Shutdown` 拿到一个"能接收/注册 hooks"的接缝——**顺序一字不改，只补注册入口**（`shutdown.go:77` 注释"Later tickets register their hooks on the same sequence"就是留给这一步的）。

### ④ 与票 228 后续片的先后关系

- 乙形同样改 `resident_windows.go`/`resident_ball_windows.go` ⇒ 与票 228 后续片**在这两枚文件上重叠、按 A480 编队串行**（同 §2④的读脏分析）。
- 乙形**额外**要碰 `run.go:336 assembleRuntime`（把它从"run 腿私有"抬成"两腿共用的接缝"）。票 228 后续片**不动 run.go**（它动的是常驻 config 接线/托盘/球持久化），所以 run.go 这一碰**不与 228 撞**；但 run.go 是**别的票（197/224/201/223）的主战场** ⇒ 乙形动 run.go 时真正要串行的是那些票的在飞腿（A480 队列：197-r3 → 224-r3 → 242 → 244）。**先做 228 后续还是先做 246，取决于先让 run.go 稳定还是先让常驻文件稳定**——两形这点一致，甲形只是把"改 run.go"换成"在常驻文件里复刻 run.go"，反而更晚撞上 197/224 在飞腿、却更早制造第二真相源。

### ⑤ 需要 owner 点头的有几枚、各是"功能"还是"契约"

- 功能级：同甲形的 1/2/3（起 loop 用既有 per-task 槽＝非契约、球驱动 UI、通道标 loaded），**再加一枚**：把 `assembleRuntime` 抬成可复用接缝（**重构，非契约**；不新增包边、不改 D43/C12）。
- 契约灰区：同甲形第 4 枚（`proc.Shutdown` 的 hook 注册入口是否算契约面）。乙形把这一枚**摊得更明显**（因为它正面需要"装配根把 hook 递进 proc"），但性质不变。
- 汇总乙形：**硬契约改动 ≈ 0**，功能 ≈ 4（比甲多"抬接缝"一枚，但这枚是消除第二真相源、把漂移风险关掉的那一枚）。⛔ 同样**不因触碰 proc/run.go 就判"不做"**。

### ⑥ 最小可跑通的那一发

- 与甲形**同一发**（『卡片进得来＋Esc 真能否决一次』、不需要 loop/provider/config/麦克风，见 §2⑥）。乙形的差别只在**这一发的 gate 由装配根建、以句柄注入常驻腿**，常驻文件不自己 new gate；最短链阶段可以**只注入一枚裸 gate（无 loop）**，把 `assembleRuntime` 的抬取留到"真跑任务"那一支再做。

---

## §4 两形对比

| 维度 | 甲形（常驻腿自装配） | 乙形（装配根注入） |
|---|---|---|
| 新增包级依赖边 | **0**（合规做，装配留 cmd/wisp）／危险做＝`proc→agent`（票238 明令挡） | **0**（函数值/句柄内部传递，跨包 import 图不动） |
| D38(e) 十步是否被破 | step 3 公共缺口；甲若用 defer 自行 cancel ⇒ 审计说"skipped"、取消发生在序列外（行为/审计分叉） | step 3 同一公共缺口；乙把它做成装配根登记的 `CancelTasks` hook ⇒ 审计真执行、卡片有据判 reject |
| 与票 228 后续片串行 | 撞 `resident_windows.go`/`resident_ball_windows.go` | 撞同两枚 **＋** `run.go`（197/224/201/223 主战场，另需与在飞腿串行） |
| owner 待批（功能/契约） | 功能 3／灰区契约 1（proc hook 入口）／硬契约 0 | 功能 4／灰区契约 1（同一枚，摊更明显）／硬契约 0 |
| 最短链可跑通性 | 同一发，但 gate 由常驻文件自建 | 同一发，gate 由装配根建、句柄注入；可只注入裸 gate（无 loop） |
| 唯一接缝（票197/238 裁定契合度） | **偏离**：复刻 `assembleRuntime` ⇒ 第二个真相源，两份 gate/UI/生命周期各自漂移 | **契合**：一个装配点喂两腿，把"审批门形状改了另一腿没跟上"这类漂移从结构上关掉 |
| 重复装配面（card-UI/lifecycle） | 高（run 腿一份、常驻腿一份） | 低（一份 gate 装配、UI 只是不同实现注入） |

- **依赖边差集（乙相对甲）**：**两者都 = 0 条新增包级边**（合规前提下）。⇒ 差集为空——**本票的选型不在"谁更不开边"，而在"要不要制造第二个装配真相源"**。甲形若为省事走危险支（让 proc 认识 loop）才会开 `proc→agent` 边，那恰恰是票 238 第 1 刀已经否掉的形状。

---

## §5 结论：我建议哪一形（＋三条最硬理由，各带 `file:line`）

- 建议：**乙形（装配根 `cmd/wisp` 注入）**，先做 §2⑥ 那条"卡片进得来＋Esc 真能否决一次"的最短链。⛔ 这是**代价比较后的选型建议，不是"这样就能做完"**——落地腿由编排者裁后再派，AC 框归编排者。
- 理由 1（带 `file:line`）：**两形新增包级边都是 0，但只有乙形守住"唯一接缝"**。`go list` 证 `cmd/wisp` 已直接 import `agent`/`agent/approval`/`ball`/`proc`；甲形若要"常驻腿自己起 loop"，等于在 `cmd/wisp/resident_windows.go:118` 附近把 `cmd/wisp/run.go:336 assembleRuntime` 复刻第二遍 ⇒ 制造**第二个审批门真相源**（gate/UI/生命周期两份，改一份忘一份）。乙形照台账 A464 刀 1「两枚正向边一律不开、改注入」＋票 197 段「装配根 cmd/wisp 是唯一的接缝」把两腿接在同一装配点上（⚠ 这是把该裁定**推广**到本票，推广动作具名，不写成仓里本就有）。
- 理由 2（带 `file:line`）：**D38(e) step 3 是两形共用的 proc 侧缺口，但乙形能把它落成"真执行"而非"审计说谎"**。现读 `internal/proc/boot_windows.go:151-161` 的 `Shutdown` 每次 `ShutdownHooks{}` 只填 `CloseJob`、step 1-7 全 nil；全仓 `grep` 零生产者。甲形自持 loop 只能在 `cmd/wisp/resident_windows.go:69` 的 defer 里自行 cancel Root ⇒ 审计里 `cancel-task-roots` 记 skipped、取消却发生在序列之外（行为/审计分叉）。乙形把它做成装配根登记的 `hooks.CancelTasks`，沿用 `cmd/wisp/run.go:782 agentRuntime.close()`（run 腿已用 `replyRoot.Cancel()` 把待决卡判成 deny）这一既有姿势 ⇒ 卡片以 `AnswerReject` 有据收口。**顺序一字不动，只补 `shutdown.go:77` 注释里本就留给"later tickets"的注册入口。**
- 理由 3（带 `file:line`）：**最短链两形同发，乙形耦合最低**。`cmd/wisp/run.go:732 confirmModeSwitch`＋`cmd/wisp/approval_always.go:112 runWidening` 证明不用模型也能用 `internal/agent/approval/gate.go:188 AdmitTextTask`+`gate.go:246/496` 造真卡；`internal/agent/approval/replies.go:252 AwaitingHuman`＋`replies.go:455 Veto`＋`gate.go:415 Veto` 把"命名卡片的 Esc"路由成 `q.reject`（票 87 已接）。乙形最短链**可只注入一枚裸 gate（不建 loop）** ⇒ 把票 228 的 `config.toml` 接线、面板（票 37）、KWS（票 41）、真任务源全部推到射程外，本票与它们彻底解耦。

---

## §6 判不动／没查的地方（诚实残余）

- ⛔ **未读 `frontend/**`／`design/**`**（两层禁地）⇒ 面板拒绝 `ChannelPanel` 这条通道我只能证"执行者缺、且本树无 WebView2 宿主"（`cmd/wisp/approval_always.go:165` 「links no WebView2 host」），面板侧要什么**不转述**，按票面 §禁区只写进票与台账、由 owner 带给他自己的前端 agent。
- **未跑 `go test`**（读码段）⇒ 两形"跑起来到底红不红"我只用 `go build`/`go vet` rc=0 侧证可编译，**不背书可测**。A480④ 那枚 `cmd/wisp` 默认层 1 红（用例名丢了）**我没归因、也不该由普查腿归因**（票面 §禁区:27 明写归因是下一枚写腿的活）。
- **「给 `proc.Shutdown` 加 hook 注册入口」到底算不算契约面——判不动**。`shutdown.go:77` 注释写"Later tickets register their hooks on the same sequence"，读起来像实现缺口；但 `Shutdown(fast bool)` 现签名不接收 hooks，改它是否触碰 D38(e) 契约，**须 owner/编排者一句定**。我按"顺序冻结、注册机制是缺口"陈述，最终归类留给裁的人。
- **"常驻进程是不是任务管线的永久之家"——只按票面主张读、未从 D2/D24/D38 正文深挖**。我据票标题"双击图标起来的那个进程，能干活"＋`resident_ball_windows.go:16-19` 注释对 D2 的转述得出"意图是让常驻腿跑任务"，但 **D2 正文（`PLAN.md:73` 一带）我未逐字核**，且那枚注释是被禁止当权威的转述层。若产品其实想让 veto **跨进程** reach 到 `wisp run`，本普查的两形框架要重议——**这条我判不动，上报**。
- **票 228 后续片的确切待改文件集**我只据 A477/A480 编队与票面列举（config/托盘退出/球位置/ball.New 失败清理），**未逐格打开票 228 面核它的 AC 落点**与 246 编辑集做全集交并；串行结论（同撞 `resident_*.go`）建立在 cited 文件上，够用但不 exhaustive。
- **`host-initiated` 卡片在"不能 exit(2) 的 GUI 进程"里的资源面**：最短链绕开了 provider/config，但一旦"真跑任务"要 `assembleRuntime`，那里 `run.go:340-381` 的 secret/config/provider 失败都返回 exit 2 —— 常驻腿不能照抄这个失败方向；这一格属于"真跑任务"那一支，**本普查只点名、不设计**。

---

## §7 归口建议（不属于本票的缺陷，逐条给票号建议；本腿一枚都不改）

- **D-1｜`internal/proc` D38(e) step 1-7 零 hook 生产者（本票只吃 step 3，其余归口）**：`boot_windows.go:151-161` `Shutdown` 每次新建 `ShutdownHooks{}`、仅 `CloseJob` 被填，`SchedulerClose/StopHotkeyAndKWS/CancelTasks/StopAudio/ReleaseSpeechSessions/DestroyPanel/FlushAndCloseDB` 全 nil ⇒ 十步序列今天 7 步空转（审计里全是 skipped）。⛔ **本票（246）只需要 step 3 把待决卡判成拒绝**；step 4（StopAudio，语音链）、step 5（ReleaseSpeechSessions，票 15/28）、step 6（DestroyPanel，票 33/35）、step 1/2（scheduler/hotkey 全量）是别的落地面。**归口建议**：把"proc 需要一个 hooks 注册接缝"单列（或并给**票 228 后续片**／或语音链/面板票各自补自己那一步），**不在 246 内一次性把七步都接上**。
- **D-2｜`internal/proc` 缺"调用方注册 hooks"的公开入口**：`shutdown.go:77` 注释承诺"later tickets register their hooks on the same sequence"，但 `Shutdown(fast bool)` 与 `Boot` 都没有递 hooks 的参数/选项（`BootOption` 现仅 `WithLayout`/`WithRegistry`，`boot_windows.go:47/:52`）。⇒ 任何票要让某个 step 真执行都得先补这枚接缝。**归口建议**：**票 228 后续片**或新立一枚"proc-shutdown-hook-registration"小票；246 落地腿若采乙形会撞上它，届时**具名引用本条**而非在 246 内私改 proc 语义。
- **D-3｜过期指认（只登记，不改）**：`cmd/wisp/main.go:8` 逐字「The floating ball GUI is ticket 07」已被票 228 装球进常驻腿证伪；`cmd/wisp/main.go:24` 用户可见 usage「its four global hot keys」与实测稳态 3/4 不符。**这两条已由台账 A477（归 228 AC#8）/A478（归 245 AC#6）具名分配**，本普查**不重复立案、不顺手改注释**，仅确认落点仍有效。
- **D-4｜A480④ `cmd/wisp` 默认层未归因红**：票面 §禁区:27 已钉"下一枚动 `cmd/wisp` 的腿必须先跑默认层取红名、判是否本票造成"。本普查腿**未跑 `go test`、未归因**（射程外）。**归口建议**：留在 246 落地腿（或先行队列里下一枚碰 `cmd/wisp` 的腿）起手第一件事，**不是本普查的活**。
- **D-5｜`resident_ball_windows.go:101-106` `ball.New` 失败路径**：注释自陈失败即 `rb.verdict` 记一句、不 new state machine；票 228 AC#9 已登记"失败支今天从未被自然走到、只有突变现形"。若本票落地腿把 veto 回调建在 `rb.b != nil` 之前，要防"球没起来却注册了 veto 通道 loaded"这一不一致。**归口建议**：并入**票 228 后续片**（球失败清理），246 落地腿顺带不自造 loaded-but-no-ball 状态即可。

---

## §8 门禁读数（本腿只跑 `go build`／`go vet`，⛔ 不跑 `go test`、不跑 `cmd/wisp` 整包）

| 尺 | 命令 | rc | 备注 |
|---|---|---|---|
| build | `go build ./...` | **0** | 2026-09-30 18:3x 现跑，全树干净可编译 |
| vet | `go vet ./...` | **0** | 同上；⛔ 本腿不跑 `go test`（读码段），`cmd/wisp` 整包那枚未归因红（A480④）不在此射程、也不该由普查腿归因 |

> ⛔ 未按 `go test`。台账 A480④ 记 `cmd/wisp` 默认层「1 红 2 绿、用例名丢了」——**归因是下一枚写腿的活，不是本普查腿的**。本腿只跑 build/vet 两发各一次坐实可编译性。
