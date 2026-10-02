# 普查 256-a1 — 「把 `approval.New` 挪进 `assembleRuntime` 会不会破 `Confirming` 那一维」的现读凭据

**本腿身份**：只读普查腿，交件＝现读凭据。⛔ 本文件不裁形（ⓘ／ⓑ 两选一由编排者落具名 `A##`）、⛔ 不提修法、⛔ 不改产码。
**票面**：`.scratch/wisp/issues/256-...-deferred-until-confirming-is-measured.md` 的 **AC#0**（本腿只答 AC#0；AC#1/AC#2/AC#3 一字节未动）。
**读法**：只有 `grep`／`sed`／`ls`／`find`／`git log`／`git show`／`git status`。**⛔ 全程未跑 `go build`／`go vet`／`go test`／`go list`／任何 `./...`**（同机有腿 `248-r2` 在跑那三包，跑一把就洗）。
⇒ 因此本文件**没有任何一行是"实跑出来的行为读数"**；凡涉及"运行时会打印／会落盘／会钳位"的句子，本腿只读的是**产生它的那一行源码**，并在各节把这一点说在句子旁边。

---

## 0. 起手锚（逐字）

| 项 | 现值 |
|---|---|
| `date -Iseconds`（起手第一发） | `2026-10-02T13:12:00+08:00` |
| `git log -1 --oneline` | `cf87f7a2 e2e-panel-1 收尾：交件态落款＋一发说快了的假话就地更正` |
| 分支 | `dev` |
| `git status --porcelain -- cmd internal \| wc -l` | `0` |
| 写点 | `.scratch/wisp/probes/256/a1/census.md`（本文件，目录本腿自建） |
| **交件时刻补记**（最后一发 `date -Iseconds` 现取） | `2026-10-02T13:12:00+08:00`；`git status --porcelain -- cmd internal \| wc -l` 复测＝**0**（＝本腿读过的两棵子树到交件此刻仍未被别人改动，§1–§3 的行号可直接复核）。⚠ 本腿**两发** `date -Iseconds` 返回**同一时刻** ⇒ 钟点在本会话里不可分辨，这两枚时间戳只当"起手／交件各一发的逐字读数"用，**不当区间长度读数** |

⚠ 共享工作树整树永远非空——本腿只按上面的**显式 pathspec 计数**判"我读的两棵子树在我读的时刻是干净的"，不当异常、也不当作"别人没在动"的证据。
票面落款锚点是 `4a9851d6`，起手 HEAD 已是 `cf87f7a2` ⇒ **票面行号自锚点之后可能被别人挪过**，本文件所有行号＝**本腿在 `cf87f7a2` 上现读**，不抄票面。

## 1. 谁在这枚门建好之前就需要它（逐处 `file:line` ＋ 挪后拿不拿得到）

**先钉两个时点**（这是本节全部判语共用的一把尺）：

- **被挪的东西现在长在哪**：`cmd/wisp/resident_approval_windows.go:109` 的 `ra.gate = approval.New(approval.Options{`，由 `cmd/wisp/resident_windows.go:123` 的 `ra := newResidentApproval()` 驱动。
- **挪进去的那个函数的调用时刻**：`assembleRuntime` 定义在 `cmd/wisp/run.go:382`；常驻腿里它**只有一个调用点**＝`cmd/wisp/resident_task_source_windows.go:265`（`run, code := assembleRuntime(spec)`），而这一行由 `cmd/wisp/resident_windows.go:206` 的 `src := startResidentTaskSource(rt, ra)` 驱动。
- ⇒ 尺＝`grep -n "newResidentApproval\|bindBallHost\|detachBall\|cancelTaskRoots\|startResidentTaskSource\|vetoByEsc\|residentStatusLine" cmd/wisp/resident_windows.go`，命中行号 **123 / 163 / 173 / 174 / 186 / 206 / 216**（本腿现读）。**206 之后才有门**＝这一列行号本身就是结论：123、163、173、174、186、216 六处全在 206 之前或与之同刻求值。

**逐处**（判语只用三个词：拿得到／拿不到／量不到；每处带尺）

| # | 那一处 | 现读 `file:line` | 它读门的哪一半 | 挪进 `assembleRuntime` 之后拿不拿得到 | 尺 |
|---|---|---|---|---|---|
| 1 | 账本绑门（**最硬的一处**） | `cmd/wisp/resident_approval_windows.go:114-120`：`ra.cards = approval.NewReplies()` 后紧接 `ra.cards.Attach(approval.HostBinding{Gate: ra.gate, VetoChannel: approval.ChannelEsc, NativeSource: …, PanelSource: …})` | 直接把 `*Gate` 当值递给账本 | **拿不到**。`HostBinding.Gate` 是在 :116 求值的**指针值**，不是延迟读取；挪后此刻那枚指针为 nil。`internal/agent/approval/replies.go:155-168` 的 `Attach` 只有"整包覆盖"语义、没有"稍后再补门"的入口 ⇒ 未绑的账本对任何答复路由都回 `ErrNoGateAttached`（`replies.go:61-64` 定义，命中行 `:326`、`:361`、`:403`、`:435`、`:461`） | `grep -n "Attach\|HostBinding" internal/agent/approval/replies.go` ＋ `grep -n "ErrNoGateAttached" internal/agent/approval/replies.go`（本腿现数：定义 1 枚、注释 1 枚、返回点 5 枚） |
| 2 | 绑球＋加载 Esc 否决通道 | `cmd/wisp/resident_windows.go:174` `ra.bindBallHost(rb)` ⇒ `cmd/wisp/resident_approval_windows.go:139-161`，其中 **`:157` `ra.gate.Channels().SetLoaded(approval.ChannelEsc, true)`** | `Gate.Channels()`（注册表），不是答复路由 | **拿不到**。:157 与 :351 是全仓产码里 `SetLoaded(ChannelEsc, …)` 的**唯一两枚**调用者（尺：`grep -rn "SetLoaded(" cmd internal tools --include='*.go' \| grep -v _test.go` ⇒ 两枚都在这一个文件里）。⇒ 挪后 boot 那一刻无法加载 ⇒ :153-160 这段"只有真有球窗、真有取消执行者才 advertise"的形状整段失去意义，而 :145-152 那枚"advertise 一枚按不动的键＝B1 禁止的形状"的守卫（它读的是 `rb.cancelHosted`，不读门）**今天仍在**，只是它下面那一行没了 | 见上；另 `cmd/wisp/resident_ball_windows.go:121` 的 `startResidentBall` 签名（第三枚入参 `onCancelEsc escVetoFunc`）证明球只需要一枚函数值，不需要门 |
| 3 | 热键答复 | `cmd/wisp/resident_windows.go:163` 把 `ra.vetoByEsc` 递进 `startResidentBall` ⇒ 本体 `cmd/wisp/resident_approval_windows.go:171-184`：`ra.cards.AwaitingHuman()` → `ra.cards.Veto(card.CorrelationID)` | 只读账本，由账本去读门（`replies.go:455-469`：`g == nil` ⇒ `ErrNoGateAttached`） | **函数值拿得到、答复拿不到**。`ra.vetoByEsc` 是方法值，捕获的是 `*residentApproval` 而非 `*Gate` ⇒ 编译期与注册期都不需要门；但按下时若门未绑，`Veto` 在 `replies.go:460-462` 直接回错，`resident_approval_windows.go:179-181` 于是把那句印成"按 Esc 否决卡片 … 被拒"。⇒ 这一处的现读形状＝**"卡片在等人、键按下去被拒"**，不是"卡片没进 Confirming" | `grep -n "func (ra \*residentApproval) vetoByEsc" -A 14 cmd/wisp/resident_approval_windows.go` ＋ `sed -n '455,470p' internal/agent/approval/replies.go` |
| 4 | 停机第 3 步（取消任务根） | `cmd/wisp/resident_windows.go:186` `rt.RegisterShutdownHook(proc.StepCancelTasks, ra.cancelTaskRoots)` ⇒ 本体 `cmd/wisp/resident_approval_windows.go:283-331`：`ra.cards.Pending()`（:289）、L2 走 `ra.cards.Reject`（:292）、L1 走 `ra.cards.Forget`（:307）＋本文件自写 `RESIDENT-WINDOW-ABANDONED`（:304-306），最后 `ra.cancel()`（:311） | 同样只经账本；`Reject` 的门检查在 `replies.go:403` 那一支（`decide`） | **注册拿得到**（方法值，:186 不求值门）；**执行时刻拿不拿得到＝取决于装配有没有成**。`ra.cancel()` 与 :285 `ra.closed = true` 两行**不读门**，所以"D38(e) 第 3 步跑了"这句话不会因挪而变谎；变的是 :292 那一支会不会回 `ErrNoGateAttached` 并让 :294-296 写出"退序列拒绝对待卡片 … 失败"。⇒ 量不到"今天真失败过"的读数（本腿不跑） | `grep -n "cancelTaskRoots\|RegisterShutdownHook" cmd/wisp/resident_windows.go cmd/wisp/resident_approval_windows.go` |
| 5 | 分离球窗（退出时那一步的另一半） | `cmd/wisp/resident_windows.go:173` `defer ra.detachBall()` ⇒ 本体 `cmd/wisp/resident_approval_windows.go:344-353`，其中 **`:351` `ra.gate.Channels().SetLoaded(approval.ChannelEsc, false)`**，外面套着 `:350 if was`（`was = ra.escLoad`） | `Gate.Channels()` | **拿得到但会空转**。defer 语句注册时不求值门；:345 `u.releaseBall()` 与 :346-349 的 `escLoad` 都不读门；只有 :351 读，而它被 :350 的"当年真加载过"守卫挡着 ⇒ 第 2 处没加载成功时这一行今天就不执行。⇒ 判语：**挪动不破这一步，但这一步保护的正是第 2 处那枚加载；第 2 处没了，它守的东西也没了** | `sed -n '344,353p' cmd/wisp/resident_approval_windows.go` |
| 6 | 启动报告那一行 | `cmd/wisp/resident_windows.go:216`（`ra.residentStatusLine()`）⇒ 本体 `cmd/wisp/resident_approval_windows.go:359-368` | **不读门**：只读 `ra.escLoad`（:361）与 `len(ra.cards.Pending())`（:367） | **拿得到**（零门依赖）。但它打印的句子由 `escLoad` 决定：`if !loaded { return "审批门未装配（…）" }`（:363-365）⇒ 挪后第 2 处失效 ⇒ 这行报告**必然**改口成"未装配"，且这是**报告与事实同向**的改口（门确实不在 boot 时点上）。另 `cmd/wisp/resident_windows.go:230` 的第二发行只读 `src.taskPosture()`／`rb.statusLine()`，不读门 | `sed -n '359,368p' cmd/wisp/resident_approval_windows.go` |
| 7 | 卡片答复面（控制台那侧） | `cmd/wisp/resident_task_source_windows.go:288` `src.surface = run.newReplySurface(src.root, ra.consoleVetoChannel(), false)` ⇒ `cmd/wisp/approval_reply.go:476-480` 的 `if rt == nil \|\| rt.gate == nil { return nil }` | 读的是 **`rt.gate`**（装配产物的字段），不是 boot 时点的那枚指针 | **拿得到，而且是唯一一处"挪了才顺"**：今天它拿到的 `rt.gate` 本来就是靠 `resident_task_source_windows.go:256-259` 把 boot 那枚门注入进去才与 boot 同一枚（`run.go:606 rt.gate = s.gate`）。⇒ `ra.consoleVetoChannel()`（同文件 `:520-524`）只读 `escLoad`，不读门 | `grep -n "newReplySurface" cmd/wisp/approval_reply.go cmd/wisp/resident_task_source_windows.go` |
| 8 | 注入守卫（决定"挪之后走哪一支"） | `cmd/wisp/resident_task_source_windows.go:252-259` 组 `runSpec{gate: ra.gate, ui: ra.ui, cards: ra.cards, taskCtx: ra.root}` ⇒ `cmd/wisp/run.go:600 if s.gate != nil {`，:601-605 半注入即拒（退出码非 0），:610-627 else 支自建 | 读 `s.gate` 是否为 nil | **拿不到门＝会掉进 else 支**。`s.gate` 为 nil ⇒ `run.go:610-612` `rt.liveCards = newNativeCards()` / `rt.ui = &consoleApprovalUI{…}` / `rt.gate = approval.New(...)` ⇒ 同进程里出现**第二枚门与第二张答复面**（一张打印到 stdout 的门），`resident_approval_windows.go:428` 那枚 `ballCardUI` 再也不是这枚门的 UI。⇒ 这一行是"挪"最直接的代价读数：**不是破 `Confirming` 的名字，而是卡片换了一张没人看的脸** | `sed -n '596,630p' cmd/wisp/run.go` |
| 9 | 举卡的入口（`AskOnTaskRoot`／`askConfirmation`） | `cmd/wisp/resident_approval_windows.go:261-265`、`:205-250`，门调用在 `:230 AdmitTextTask`、`:242 PendingWindow`、`:244 PendingApproval` | 三门全读 | **量不到**（现状即无读数）：本腿现读 **产码调用者枚数＝0**。尺＝`grep -rn "AskOnTaskRoot\|askConfirmation" cmd tools internal --include='*.go' \| grep -v resident_approval_windows.go`，非测试命中只有注释三处（`cmd/wisp/resident_task_source_windows.go:14`、`cmd/wisp/resident_windows.go:194`）与 `_test` 五处（`resident_approval_246_windows_test.go:97,161`、`resident_approval_live_246_windows_test.go:115,220,311`）。⇒ 这三行读门的代码今天**没有一条产码路径能走到**，"挪会不会破它"这一格只能记作量不到 | 见上两条 grep |
| 10 | 托盘「**允许一次**」（票面 AC#0 点名的一处） | **不存在**。托盘菜单项在 `internal/ball/tray_windows.go:86-91`，四枚 id 定义在 `:20-23`：`menuOpenPanel`／`menuMute`／`menuPauseWake`／`menuExit`，标签逐字「打开面板」「静音」「暂停唤醒」「退出」。常驻腿的挂接点在 `cmd/wisp/resident_ball_windows.go:179-182`（`OnTrayPanel`／`OnTrayMute`／`OnTrayPauseWake`／`OnTrayExit`），**四枚里没有一枚碰门、账本或答复面**（`OnTrayExit` 只做 `recordTrayExit`，本体 `:308-312` 只写日志） | 无 | **量不到——这一处今天没有**。尺：`grep -rn "允许\|Allow" internal/ball/*.go \| grep -v _test.go` ⇒ **零命中**；`grep -rn "\.Allow(" cmd internal --include='*.go' \| grep -v _test.go` ⇒ 唯一命中 `cmd/wisp/approval_reply.go:215`（控制台答复面的 `yes` 支，不是托盘）。⇒ 票面把"托盘「允许一次」"和"卡片 UI／热键答复"并列为三处现状，**这一并列不成立**：那一处是**待建**（票 244 的 J9 射程），不是**待挪** | 见上三条 grep |
| 11 | 原生侧 `NativeSource` | `cmd/wisp/resident_approval_windows.go:62-65` 两枚常量 → `HostBinding` 的 `:118-119` → `replies.go:116-124`（注释：它们是 LABELS，不是 authority） → 落到 `gate.go:719`、`:733` 两处 `logf` | 只把门当**日志汇**用 | **不读门本体**（`r.Source` 在全仓产码里没有任何分支读它，尺：`grep -rn "\.Source\b" internal/agent/approval/gate.go internal/agent/approval/queue.go` ⇒ 只有 :719/:733 两枚 `logf` 参数位）。⇒ 挪动对这一处的影响**等价于第 1 处**：门没绑 ⇒ 那句 audit 行根本写不出来。⚠ 引文性质＝**内容引用**（`gate.go:705-707` 的 "Source is ADVISORY ONLY … no branch in this package reads it"） | 见上 |

**§1 小结（只报形状，11 行全部对上号）**：常驻腿在 `assembleRuntime` 之前**真的读那枚 `*Gate` 指针**的地方共 **3 处**＝第 1（`:116` 绑账本）、第 2（`:157` 加载 Esc）、第 5（`:351` 卸载 Esc，被第 2 处的结果挡着）；**经账本间接读门**的 **2 处**＝第 3（热键答复）、第 4（退出第 3 步）；**不读门**的 **3 处**＝第 6（启动报告）、第 7（控制台答复面，读的是 `rt.gate` 字段）、第 11（`NativeSource` 只是审计标签）；**只读"nil 与否"、决定走哪一支**的 **1 处**＝第 8；**量不到**的 **1 处**＝第 9（举卡入口产码调用者＝0）；**票面点名但今天不存在**的 **1 处**＝第 10（托盘「允许一次」）。
⇒ 挪动的最小爆炸半径读数是**第 1、2 两行**：`:116` 与 `:157` 都要求"门在 boot 时点已经存在"，而 `assembleRuntime` 在 `resident_windows.go:206` 之后才可能被叫到（`resident_task_source_windows.go:265`）。

## 2. `Confirming` 那一维今天由什么决定（以及门在它生命周期里被读的时刻）

### ① 权威表里 `Confirming` 只有三枚行，而**这三枚行在常驻腿今天没有执行者**

D43 转移表里凡涉 `Confirming` 的行，全表共三枚（尺：`grep -n "StateConfirming" internal/statemachine/table.go` ⇒ 三枚命中）：

| 行 | 现读 | `file:line` |
|---|---|---|
| #17 入 | `From: StateActing, Event: EvApprovalNeeded, To: StateConfirming, Guard: "level-L1"` | `internal/statemachine/table.go:117-121`（行体 `:118`） |
| #21 出（超时＝执行） | `From: StateConfirming, Event: EvConfirmExpired, To: StateActing, SideEffects: tools.execute` | `internal/statemachine/table.go:150-151` |
| #22 出（否决） | `From: StateConfirming, Event: EvVeto, To: StateActing` | `internal/statemachine/table.go:154-155` |

**"由谁发射"这一问，本腿量到的是：产码里没人发射。**三把尺：

1. `statemachine.New` 全仓产码**只有一枚调用者**＝`cmd/wisp/models.go:303`（`Initial: statemachine.StateFirstRun`，模型下载那一趟 D43 #2/#37 的走查）。尺：`grep -rn "statemachine\.New\|statemachine\.Machine" cmd internal tools --include='*.go' \| grep -v _test.go` ⇒ **仅此一命中**（外加 `internal/statemachine/machine.go:60` 的定义自身）。⇒ 常驻腿与 `wisp run` 腿**都不持有一枚 Machine**，#17/#21/#22 这三行在两条腿里都不可能被"走"一遍。
2. `EvApprovalNeeded` 的产码发射者＝**0 枚**（尺：`grep -rn "EvApprovalNeeded" cmd internal tools --include='*.go' \| grep -v _test.go` ⇒ 只有 `internal/statemachine/events.go:35` 的定义与 `table.go:118/:123` 的两行表体）。
3. `EvVeto` 的产码发射者＝**只有 `cmd/balldebug/main.go:596`、`:606`、`:608` 三枚**（调试台件，不是产品进程）。`EvConfirmExpired` 的表内超时是 `internal/statemachine/timeouts.go:48`（`{state: StateConfirming}: {3 * time.Second, EvConfirmExpired}`），它属于那枚 Machine 的 timeout 机制，常驻腿不持有 Machine ⇒ 不吃这一行。

### ② 常驻腿里 `Confirming` 今天真正的决定者：三行

| 时刻 | 现读 `file:line` | 那句决定 |
|---|---|---|
| 进入 | `cmd/wisp/resident_approval_windows.go:443` `b.SetState(stateForCardLevel(p.Level))`；名字产自 `:501-506`（`if level == "L1" { return statemachine.StateConfirming }`） | **`p.Level` 这枚字符串就是全部判据**。`p` 由门在 `internal/agent/approval/gate.go:287` `p := g.promptFor(d, corr, g.window, 0, bv, "")` 组好、`:288` 交进 UI |
| 退出（正常） | `cmd/wisp/resident_approval_windows.go:465-473`（`Update` 的 `EventDismissed`/`EventStarted` 两支）→ `:469 u.settleOrb()` → `:483-493`，其中 **`:492 b.SetState(statemachine.StateSleeping)`**，前置守卫 `:484` 是 `cards.AwaitingHuman()` | 退出**不经表**：它经账本"还有没有在等人的卡"这一枚事实 |
| 退出（兜底） | `cmd/wisp/resident_approval_windows.go:216-228`（`askConfirmation` 的 `defer`，注释自陈"任务 ctx 被取消时门不发任何 dismissal 事件"）→ 同一枚 `settleOrb()` | 这一支存在，正是因为 `gate.go:340 case <-ctx.Done()` **不**调 `g.ui.Update`（尺：`sed -n '288,345p' internal/agent/approval/gate.go` 里 `Update` 只出现在 veto/expired 两支） |

★ **球这一侧不校验转移表**（这是本节最重要的一枚尺）：`internal/ball/ball_windows.go:311-313` 的 `SetState` 只是 `b.sta.PostTask(func() { b.applyStateLocked(s) })`，而 `applyStateLocked`（`:342` 起）第一行是 `prev := b.curState; b.curState = s`（`:343-344`）——**没有 `Fire`、没有 `Transition`、没有查表**。⇒ 现读结论：**"Sleeping 直落 Confirming" 今天就能发生，且不违反任何会被执行的东西**（表里 #17 的 FROM 是 Acting）。这条现状与"挪不挪"**无关**，本腿只登记它，因为编排者问的是"会不会破那一维"，而这一维今天根本不在表的射程里。
⛔ 本腿因此**不提**"该不该让球走表"——那是契约变更（D43 射程），超出 AC#0。

### ③ 那枚门在 `Confirming` 生命周期里被读的**时刻**（逐行，尺＝`grep -n` 现取）

`Confirming` 的整条生命都在 `Gate.PendingWindow` 里（`internal/agent/approval/gate.go:246` 起），门自身被读的时刻：

| 时刻 | `file:line` | 读的是门的哪一半 | 值何时被定死 |
|---|---|---|---|
| 窗口长度 | `gate.go:285` `deadline := g.clock.After(g.window)` | `g.window` | **`approval.New` 的那一发**：`gate.go:137-145` 的钳位（`win <= 0 → DefaultL1Window`；`< MinL1Window → Min`；`> MaxL1Window → Max`）。常量在 `internal/agent/approval/queue.go:116`（3s）／`:121`（2s）／`:122`（3s） ⇒ **构造之后再无任何入口能改它**（`Window()` 只有 getter，`gate.go:170`） |
| 卡片内容 | `gate.go:287` `g.promptFor(...)` → 本体 `:591`，其中 `:603 Window: win`、`:605 Channels: g.channels.Statuses()` | `g.window` ＋ **通道注册表** | 同上：注册表的 `Loaded` 位由 `SetLoaded` 改，常驻腿里改它的唯一地方是 `resident_approval_windows.go:157`/`:351` |
| 交付 UI | `gate.go:288` `g.ui.Prompt(ctx, p)` | `g.ui`（＝注入的那枚 UI） | `Options.UI`，同样在构造时定死（`gate.go:129-131` 的 nil 兜底：`:131 ui = UIFuncs{}`） |
| 否决可用否 | `gate.go:299` `if err := g.channels.check(v.Channel); err != nil` | 通道注册表 | ⚠ **这一行是 `Confirming` 那一维里唯一被"实时"读的门内状态**：不在 `escLoad`/`SetLoaded` 里定死，就在 veto 落下的那一刻判"未知/未加载通道" |
| 超时＝执行 | `gate.go:324` `g.markStarted(...)` ＋ `:332` 的 `ANSWER-EXPIRED` 行（`after=%s` 打的正是 `g.window`） | `g.window`、`g.logf` | 构造时 |
| 弃等 | `gate.go:340` `case <-ctx.Done()` | ctx（非常量、可后绑） | 常驻腿里这枚 ctx 是 `ra.root` 派生（`resident_approval_windows.go:262`），`ra.root` 在 `newResidentApproval` 的 `:106` 建，**与门同批但不同物** |

### ④ 量出来的硬形状：这一维**不破**，但会**整维消失**——⛔ 本腿在此停手

把 ①②③ 合起来读，"挪进 `assembleRuntime` 会不会破 `Confirming`"这一问的现读答案是**两段的**：

- **表的那一半：不破，也不需要新增状态词或改转移表。**三枚理由：常驻腿不持有 Machine（②③尺 1）；球不查表（②★ 尺）；`Confirming` 这个名字的产生点是 `stateForCardLevel`（`:501`）与 `Replies.WaitingState`（`internal/agent/approval/replies.go:286-302`，其中 `:298 return statemachine.StateConfirming, true`），两者都只吃**卡片自己的 `Level` 字符串**，不吃门的构造时刻。⇒ **⛔ 本票不必上报"要改 D43"，那一格量到的结果是"不动表也挪得动"。**
- **存在的那一半：会消失，条件是这台机器上有没有交互控制台。**尺：`cmd/wisp/resident_task_source_windows.go:224-231`——`if console == nil && injectedText == "" { …; return nil }`，而 `assembleRuntime` 在同一函数的 **`:265`**。`interactiveStdin()`（`cmd/wisp/approval_reply_stdin_windows.go:41`）在 stdin 不是控制台输入缓冲时回 nil（管道／文件／**无 handle**）。⇒ **挪动之后，一枚由 Explorer 双击拉起、没有控制台输入的常驻进程（正是 D2／票 228 那个"用户真正启动的进程"形状）永远走不到 `:265` ⇒ 永远没有门 ⇒ 没有 `PendingWindow` ⇒ 没有 `Prompt` ⇒ 球永不进 `Confirming`、Esc 永不借、`:157` 永不加载。**
  这一支的形状不是"那一维变弱"，是"**那一维的载体不存在**"。⇒ 编排者要的"量不到的那一处要具名说量不到"本腿没有；要的是"**量到了一处会让两选一里 ⓘ 那一支改变产品形状的前置条件**"＝就是这一条。⛔ **本腿到此停手，不裁、不提修法**（"把门的构造挂到别处""让无控制台支也装配一次"这类都不写，那属 AC#1 之后由具名 `A##` 批准的射程）。
- 与之同族、但**不属本票**的一枚既有账：`docs/evidence/s1/246-resident-task-source-v2.md:212` 与台账 `docs/reports/pending-and-issues.md:10107` 都记过"窗口那项即便接了也钳在 3s ⇒ 只有超时是真差异"。本腿现读**复认其钳位形状**（`gate.go:143-144` 的 `win > MaxL1Window → MaxL1Window` ＋ `queue.go:122` 的 3s），故 `l1_window_sec` 即便接上也只能在 [2s,3s] 内动；本腿**不重跑**、不引其结论之外的东西。

## 3. `Grants` 为 nil 这一条的真实边界

### ① 这枚 nil 是从哪儿进来的（一行）

`cmd/wisp/resident_approval_windows.go:109-113` 的 `approval.New(approval.Options{` 只填 `UI`／`Channels`／`Logf` 三枚字段 ⇒ `internal/agent/approval/gate.go:154` 的 `grants: o.Grants` 收下一枚 nil。
尺（构造点的字段枚数）：`sed -n '109,113p' cmd/wisp/resident_approval_windows.go` ⇒ 三行字段、无 `Grants`、无 `Window`、无 `ApprovalTimeout`。

### ② `g.grants` 在整个包里有几处被读——**两行、一支**

尺：`grep -n "g\.grants" internal/agent/approval/gate.go` ⇒ **两枚命中，同在一条分支里**：`gate.go:667` 的 `if g.grants == nil {` 与 `:678` 的 `id, err := g.grants.Record(ctx, tool, p)`。两枚都在 `(*Gate).allowSession`（本体 `gate.go:652` 起）里。⇒ **`Grants` 为 nil 不改任何判定，只改"落不落盘＋写哪一句审计"**，与本票 §0 的姿态一致（只读、不改）。

### ③ 「本会话内允许」在常驻腿今天真实的前置门槛是**四道**，`Grants` 为 nil 只是第四道

| 门槛 | 现读 `file:line` | 常驻腿今天的状态 |
|---|---|---|
| 甲：只能是 **L2** 卡 | `internal/agent/approval/replies.go:356-358` `if card.Grant == "" { return ErrRouteHasNoAllow }`；而 L1 那一支的 grant 参数在 `gate.go:287` 就是空串（`g.promptFor(d, corr, g.window, 0, bv, "")`） | 成立：常驻腿能举 L1 也能举 L2（`resident_approval_windows.go:242`／`:244`），所以这一道**只挡 L1**，不是"常驻腿没有会话档"的理由 |
| 乙：账上要有这张卡且门已绑 | `replies.go:351-362`（`ErrNoTrackedCard` → `ErrNoGateAttached`）；卡由 `resident_approval_windows.go:439 u.ra.cards.Record(p)` 记 | 成立（今天门在 boot 就绑了，`:115-120`） |
| 丙：得有人**说得出** `session` 这个词 | 动词表 `cmd/wisp/approval_reply.go:566 case "session":` → `:567 return s.session(corr)` → 本体 `:257`；循环 `runReplyLoop` 定义在 `:530`。**常驻腿里 `AllowSession` 的产码调用链只有这一条**（尺：`grep -rn "AllowSession" cmd internal tools --include='*.go' \| grep -v _test.go` ⇒ 除 `replies.go`／`gate.go` 的定义与注释外，产品侧命中只有 `approval_reply.go`） | **这一道才是真边界**：常驻腿的答复面在 `resident_task_source_windows.go:288` 才建，而它前面是同一函数 `:224-231` 的"无控制台 ⇒ return nil"与 `:265` 的装配失败分支。⇒ **无交互控制台的常驻进程今天根本没有 `session` 这个词的入口**，`Grants` 是不是 nil 对它无关 |
| 丁：`Grants` 非 nil | `gate.go:667` | **不成立**（§3①）。⇒ 前三道全过之后，才轮到这一道把规则吞掉 |

### ④ "有没有一行日志记着这件事"——有，两形，且本腿读的是产生那行的源码（⛔ 不是实跑读数）

| 形 | `file:line` | 触发条件（现读） |
|---|---|---|
| 第一形 | `internal/agent/approval/gate.go:667-670`：`g.logf("approval: GRANT-DROPPED corr=%s tool=%s paths=%d (本机没有接入会话授权记账，本次按「仅本次」放行，没有落盘任何规则)", …)` | `g.grants == nil`，**且在放行已经成功之后**（`gate.go:664` 的 `allowScoped` 先跑；顺序由 `gate.go:642-649` 的注释钉着：先放行、后记账） |
| 第二形 | `internal/agent/approval/gate.go:673-675`：`g.logf("approval: GRANT-DROPPED corr=%s (卡片主题在答复前已离开队列，未落盘任何规则)", corr)` | `tool == ""`，即 `:657` 的 `g.q.sessionSubject(corr)` 返回 `!ok`。**这一形与 `Grants` 无关**——见下面 ★★ |

落盘路径（常驻腿侧的 `Logf` 是谁）：`resident_approval_windows.go:112 Logf: ra.residentAuditf` → 本体 `:129-133`，它**同时**做 `slog.Info("audit: " + line)` 与 `fmt.Printf("wisp: [audit] %s\n", line)`。slog 那条能到硬盘，是因为 `installLogSink` 在 `cmd/wisp/resident_windows.go:63` 就装上了，**早于** `:123` 建门。⇒ 现读结论：**两形都既有 stdout 也有落盘审计行**，不是"只到 stderr 的哑支"。
⚠ 但这两句是**"这一行源码会写这句话"**的读数；本腿禁跑，**没有**"今天真写过一发"的实跑凭据。这一格与既有账同向：`docs/evidence/s1/248-settings-write-path-v1.md:277` 逐字记着"本腿…没读到用户按了本次会话内允许、今天不生效的实跑一行"（**内容引用**，非本腿实跑）。

★★ **一枚会影响 AC#2 判据形状的量读数**（本腿不改判据、只报形状）：AC#2 写的是"`GRANT-DROPPED` 那行**必须不再出现**"。上表第二形（`gate.go:673-675`）的产生条件里**不含** `Grants`——它是"卡片主题在答复前已离开队列"，挪门挪不掉它。⇒ 判据字面读下来，一枚接了 `Grants` 的常驻腿仍可能因第二形写出 `GRANT-DROPPED`；**"不再出现"要不要区分两形，是 AC#1/AC#2 的定案题，本腿只登记**。（既有仪器对这两形**不加区分**：尺 `internal/agent/approval/ticket224_reply_grant_test.go:283` 只查 `f.log.has("approval: GRANT-DROPPED")` 这个前缀。）

### ⑤ 有没有一枚钉守住"常驻腿这一格"——现读：没有

尺：`grep -rn "GRANT-DROPPED" cmd/wisp/*_test.go internal/agent/approval/*_test.go` ⇒ 三枚命中，**没有一枚在常驻腿**：
- `cmd/wisp/ticket224_assembly_test.go:194`：反向钉，"**有** ledger 的宿主不许出现 GRANT-DROPPED"；
- `internal/agent/approval/ticket224_reply_grant_test.go:128`、`:283`：包内仪器，测的是"没给 recorder 时要说真话"，与 `cmd/wisp` 的装配次序无关。
另：`grep -rn "Grants" cmd/wisp/*_test.go` 的命中全在 `run_mode101_test.go`／`ticket224_assembly_test.go` 的 **run 腿**侧。⇒ **常驻腿的"会话档不落盘"这一格今天无钉**：它只被两处**注释**与三处**文档**记着（`resident_task_source_windows.go:43`、`run.go:596`、`pending-and-issues.md:10107`）。

## 4. 我可能判错的条目

1. **§1 的"两个时点"只覆盖两枚产码调用者**。尺：`grep -rn "assembleRuntime(" cmd internal tools --include='*.go' \| grep -v _test.go` ⇒ 今天**只有两枚**：`cmd/wisp/run.go:249`（`wisp run` 腿）与 `cmd/wisp/resident_task_source_windows.go:265`（常驻腿）。⇒ 若票 244 的 GUI 子系统那一串后来在 boot 更早处多装配一次，我的"六处全在 `resident_windows.go:206` 之前"要重画（爆炸半径那两行 `:116`/`:157` 大概不变，行号会变）。
2. **§2① 的"产码无人发射 D43 那三行"可能因 grep 面窄而漏**。我只扫了 `statemachine.New`／`statemachine.Machine`／`Fire(` 三组模式，`Fire(` 那一组我只扫到 `cmd/wisp/approval_reply.go`、`internal/ball`、`internal/agent` 三处附近（**未做全仓 `Fire(` 穷举**）。⇒ 若某处以接口值或别名导入的方式走表（例如 `sm "…/statemachine"`），①要改写。能验的尺：`grep -rn "\.Fire(" cmd internal tools --include='*.go' \| grep -v _test.go` 全量。
3. **§2④ 那枚"必然"是静态推断，不是实跑**。它成立的前提是"双击／Explorer 拉起的常驻进程 `interactiveStdin()` 回 nil"，而我只读了 `cmd/wisp/approval_reply_stdin_windows.go:41` 那一段 Win32 判定（`GetStdHandle` ＋ `GetConsoleMode`）。⇒ 若 `-H=windowsgui` 那一支（票 244 K1/AC#5 在册）另有形状，"必然"应降级为"取决于链接旗标与启动方式"。**本腿无法把这两者分开**（不能跑）。
4. **§1 第 5 处（detachBall 空转）判语依赖 `escLoad` 只有两枚写点**。尺（本腿已跑）：`grep -n "escLoad" cmd/wisp/resident_approval_windows.go` ⇒ 命中 **5 行**：`:93` 声明、`:154 = true`、`:347` 读、`:348 = false`、`:361` 读。⇒ 写点确实只有两枚，判语成立；但若将来有人在别的文件里给这枚字段加第三枚写点（同包可写），该行判语作废。
5. **§3③ 丙"会话档入口链只有一条"用的是字面 grep**（模式 `AllowSession`）。若有一枚经 `NativeAPI` 接口值的间接调用不出现方法名，我会漏。能验的尺：`grep -rn "Native()\.AllowSession\|NativeAPI" cmd internal tools --include='*.go' \| grep -v _test.go`。
6. **§3④"审计行能落盘"我只量到了装配次序**（`installLogSink` 在 `resident_windows.go:63`、建门在 `:123`），**没有量到 `slog.Info` 走的 handler 是不是那枚 sink**（`installLogSink` 本体在 `cmd/wisp/logsink.go`，本腿只读它的注释射程）。⇒ 若 handler 不是全局替换，"落盘"那半句要撤回，只保留 stdout 那半句。
7. **本文件全部行号钉在起手 HEAD `cf87f7a2`**。本腿自己已在此之上产生多枚 commit（全是新增普查件，未碰 `cmd`／`internal`）；交付前那枚 `git status --porcelain -- cmd internal` 若不再为 0，说明有别的腿动了我读的树，§1–§3 的行号需按新 HEAD 复量。
8. **§1 第 8 行"掉进 else 支"是编译期分支推断，不是运行读数**：我没有量"`s.gate` 为 nil 时 `resident_windows.go` 还会不会把 `ra.ui` 递给别处"——递 UI 这件事本身发生在 `resident_task_source_windows.go:257`，挪动之后那一行填什么**属未定**（属 AC#1 射程，本腿不猜）。

## 5. 判不动／没测到的地方

1. **一切"真打印／真落盘／真钳位"的运行读数**：本腿禁跑 `go build`／`vet`／`test`／`list`（同机 `248-r2` 在跑那三包）。⇒ `GRANT-DROPPED` 有没有在谁的机器上真写过一发、`ANSWER-EXPIRED` 的 `after=` 真值、球真进 `Confirming` 的帧，三格**全空**。
2. **挪之后门的"存活时刻分布"**：`assembleRuntime` 自身可能先失败（`run.go:392-395` 的 `secret.NewStore` 失败 ⇒ `return nil, 2`；`run.go:417-418` 配置未就绪 ⇒ `return rt, 2`），⇒ 门的存活窗口在不同失败支下不是同一形状，量它要真跑常驻 exe。
3. **D43 表侧的守钉**：`internal/statemachine/table_test.go` **本腿未读**（不属 AC#0 射程）。若那里含"产码可达性"级别的断言，§2① 的"0 枚执行者"要与之对表。⚠ 我刻意不去动那枚文件的判语，因为表本身是冻结契约（D43）。
4. **`nativeCards` 与 `Replies` 在常驻腿是不是同一枚账本**：静态读下来 `run.go:608 rt.liveCards = &nativeCards{h: s.cards}` 让 `rt.liveCards.h == ra.cards`，但这是**指针同一性的推断**，本腿没有仪器凭据（钉它属 AC#1 之后）。
5. **干净机器上那枚 panel 句真会打印**：我只静态读到链条 `resident_windows.go:142 → panel_resident_windows.go:189 → panel_inbound.go:230-233 → internal/config/manager.go:102-104 → internal/config/loader.go:67-69（os.ReadFile 失败即返错）`。⇒ "缺 `config.toml` ⇒ 那一句必然出现在启动输出里"是**链条读**，不是实跑读。
6. **票 244 J9（托盘「允许一次」）有无实现腿已落**：§1 第 10 行量的是起手 HEAD 的现状。若那格已在别的分支/工作树里落码而没进 `cmd`／`internal` 的我这一份读面，我的"不存在"要重判（本腿 `git status -- cmd internal` 为 0，看不到未提交的他腿产码）。
7. **真机 Esc 占用形状**：`resident_approval_windows.go:454-458` 那支（借键被 Win32 拒）在本腿射程内只读到源码，没读运行时。
8. **`248-r2` 那枚腿的门禁读数**：与本腿无关，也读不到（它的产物不在我读的树里）。

## 6. 我推翻／更正编排者哪一句

**复认（两处现读都还在，都是成立断言）**：

- **read #1 成立，但区间起点偏一行**。票面 `§现量.1` 写的是 `cmd/wisp/resident_approval_windows.go:108-113`；现读 **`:109` 才是 `ra.gate = approval.New(approval.Options{`**，`:108` 是 `ra.ui = &ballCardUI{ra: ra}`，三枚字段在 `:110 UI`／`:111 Channels`／`:112 Logf`，闭括号在 `:113`。⇒ 断言内容（"只给 UI／Channels／Logf 三项，没有 `Grants`、没有超时、没有 L1 窗口"）**逐字成立**，尺＝`sed -n '109,113p' cmd/wisp/resident_approval_windows.go`。⚠ 既有账里 `docs/evidence/s1/248-settings-write-path-v1.md:163` 用的是 `:109-113`，与本腿一致（**内容引用**）。
- **read #2 成立，且我把"由哪个条件决定"量全了**。句子逐字还在：`cmd/wisp/resident_windows.go:146` 的 `fmt.Printf("wisp: panel host unavailable (%v): the panel hot key and the tray item will be recorded, not executed\n", rpErr)`。**决定它的是 `:144` 的 `if rpErr != nil`**（不是 `panel.statusLine()`、也不是球窗那一枚），`rpErr` 来自 `:142` 的 `newResidentPanelManager(rt.Layout.DataDir)`，该函数**唯一**的非 nil error 出口是 `panel_resident_windows.go:189-192`（转 `:166 newResidentComposerDispatch` → `panel_inbound.go:228`，两枚 return：`:230-233` 配置未就绪、`:237-247` 权限档位不可用）。⇒ 顺带量到两枚下游（同一枚 nil `panel`）：boot 报告 `resident_windows.go:216` 打 `"this process has NO panel thread (the panel host never started)"`（`panel_resident_windows.go:494-496`）；热键／托盘那一按走 `resident_ball_windows.go:92-97` 的"**executor 回了 false**"支（warn ＋ "the panel thread took no request…"），**不是** `:89-91` 的 `recordBallGesture` 支——所以"recorded" 一词在这句里指的是 `:95` 那枚 `slog.Warn`，与票 33 之前"根本没有 executor"那种记录不是同一支。**这是对编排者那句读数的一次精化，不是推翻。**

**推翻／具名不成立（三条）**：

1. **票面 AC#0 的三处并列不成立**：`常驻卡片 UI／热键答复／托盘「允许一次」各一处` —— 第三处**今天不存在**（§1 第 10 行，尺：`internal/ball/tray_windows.go:86-91` 的四枚菜单项 ＋ `grep -rn "Allow" internal/ball/*.go \| grep -v _test.go` 零命中）。⇒ AC#0 能交的"现状逐处"只有两处，第三处只能是"待建，且它的归口在票 244 J9"。
2. **票面 `§现量.2` 把"钳位在 `gate.go:139-145`"标成母票 `246-v2` 的行号、属待验**——本腿复认：**内容成立、行号是 `gate.go:137-145`**（`:137 win := o.Window`，`:139-140` 是 `win <= 0 → DefaultL1Window`，`:141-144` 是 `<Min`／`>Max` 两枚钳位）。常量三枚在 `internal/agent/approval/queue.go:116`（3s）／`:121`（2s）／`:122`（3s），`DefaultApprovalTimeout` 在 `:107`（300s）——**与票面写的 300s／3s／3s 一致**。
3. **"ⓘ 那一支会不会破 `Confirming`"这一问，本腿把问题换了**：表的那一半**没有可破的东西**（§2①②：常驻腿不持有 Machine、球不查表、名字只吃 `p.Level`）；真有形状变化的是"**门在哪一刻存在**"（§2④：`resident_task_source_windows.go:224-231` 那一枚提前 return 让 `assembleRuntime` 在无交互控制台时永不被叫）。⇒ 编排者若按原问法裁"不破就可以挪"，会裁在一枚**没被那一问问到**的前置条件上。**本腿不替编排者补这一裁。**

**顺带复认（未推翻）**：既有账里"`[risk]` 两项在产码里的读取点只有 `cmd/wisp/run.go:615-616`"这一格，本腿现读**逐字命中**（`:615 Window: time.Duration(cfg.Risk.L1WindowSec) * time.Second`、`:616 ApprovalTimeout: time.Duration(cfg.Risk.ConfirmTimeoutSec) * time.Second`），且这两枚只在 `run.go:610` 的 else 支里（＝没注入门时才走到）。⚠ 性质＝**内容引用＋行号复认**，非本腿实跑。

## 末节（一句话）

**现读形状：表那一维不动也不破（不需要新增状态词、不需要碰 D43 转移表）；但"挪"会把那枚门的存在时刻从"进程启动即存在"改成"仅当这台机器上这发常驻进程有交互控制台、且装配根没先失败"——在无控制台那一支里 `Confirming` 不是变弱，是没有载体（`resident_task_source_windows.go:224-231` 挡住 `:265`）。**
⇒ 三选一本腿报：**要动的不是转移表，是"装配次序 × 入口条件"那一格**；⛔ 本腿不裁ⓘ／ⓑ。
