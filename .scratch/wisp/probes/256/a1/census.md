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

**§1 小结（只报形状）**：常驻腿在 `assembleRuntime` 之前**真的读那枚 `*Gate` 指针**的地方共 **3 处**＝第 1（`:116` 绑账本）、第 2（`:157` 加载 Esc）、第 5（`:351` 卸载 Esc，被第 2 处的结果挡着）；**经账本间接读门**的 **2 处**＝第 3（热键答复）、第 4（退出第 3 步）；**不读门**的 **3 处**＝第 6、第 7、第 11。**票面点名的"托盘「允许一次」"＝0 处**（第 10 行现读不存在）。
⇒ 挪动的最小爆炸半径读数是**第 1、2 两行**：`:116` 与 `:157` 都要求"门在 boot 时点已经存在"，而 `assembleRuntime` 在 `resident_windows.go:206` 之后才可能被叫到（`resident_task_source_windows.go:265`）。

## 2. `Confirming` 那一维今天由什么决定（以及门在它生命周期里被读的时刻）

本节结论：`Confirming` 在常驻腿**不是由 D43 转移表决定的**——那张表在本仓产码里只有一枚 `statemachine.New` 的调用点，且不在常驻腿。挪动**不需要新增状态词、不需要改转移表**；但量出来一处**更硬的形状**（门的"存在时刻"与"无控制台"那一支），具名在 §2④。

## 3. `Grants` 为 nil 这一条的真实边界

本节结论：「本会话内允许」在常驻腿今天有**三道前置门槛**，`Grants` 为 nil 只是第三道；这一条**确实有一行审计记着**，产生点与落盘点本腿都读到了源码行（`file:line` 在下面），但**本腿没跑过它**，所以那是"这一行会写这句话"的读数、不是"今天真写过"的读数。

## 4. 我可能判错的条目

见 §4 正文（六条，逐条带"错在哪／用什么尺能验"）。

## 5. 判不动／没测到的地方

见 §5 正文（八条，含"必须跑测试才拿得到"的那一类，本腿一律不做）。

## 6. 我推翻／更正编排者哪一句

见 §6 正文。要点：编排者给的两处现读**都还成立**（逐字复认），但**第一处的行号已漂 1 行**、**票面 AC#0 里"托盘「允许一次」"这一处今天不存在**（它是票 244 的射程，不是常驻腿的现状）。

## 末节（一句话，只报形状，不替编排者裁）

见文件末行。
