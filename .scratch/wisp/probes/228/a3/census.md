# 228-a3 普查：托盘「允许一次」在"此刻没有卡在人等"时该长成什么样（只读，零产码）

## 起手锚点

- `date "+%Y-%m-%d %H:%M:%S %z"` 逐字读数：`2026-10-01 11:27:59 +0800`
- `git log -1 --format=%H` 逐字读数：`eed229e48cb9599c29f72fea6a1170739a5f2b41`
- 分支：`dev`（`git rev-parse --abbrev-ref HEAD` ＝ `dev`；`git status -sb` 首行 `## dev...cnb/dev [ahead 169]`，"未推枚数"是起手那一瞬的读数，会漂）
- 本腿代号 `228-a3`，只读；并发的写腿 `33-r3` 在 `cmd/wisp`，另枚只读腿 `33-a2` 在 `internal/ball` 的 STA／泵面。

---

## ① 菜单侧拿不拿得到"有卡吗"

**今天传进菜单构造的状态＝两枚布尔，没有第三枚。**

- 构造函数：`internal/ball/tray_windows.go:72`
  `func showMenu(hwnd windows.HWND, muted, pausedWake bool) uint32`〔量 R6〕。
- 唯一调用者：`internal/ball/ball_windows.go:674`
  `sel := showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)`，位于 `case wmAppTray:`（`ball_windows.go:669`）之内〔量 R5〕。
  这一支在 wndproc 里，而 `tray_windows.go:69-71` 的注释逐字写「Runs on the STA thread inside the tray message handling」
  ⇒ **回调跑在 `ui-sta` 线程内部**（票面 `:92` 同一口径，它引 `ball_windows.go:669-686`）。
- 那两枚字段的写入者：`func (b *Ball) SetTrayChecks(muted, pausedWake bool)`（`ball_windows.go:919`），体是 `b.sta.PostTask(...)`（`:920-923`）＝**异步推**进 STA 队列，不返回值。
- ⚠⚠ **本腿量到一条会影响"推"这一形成色的事实**：`SetTrayChecks` 与 `SetTrayTip`（`ball_windows.go:919`/`:927`）**全仓零枚调用者**。
  尺＝`grep -rn "SetTrayChecks\|SetTrayTip" --include=*.go cmd internal` ⇒ **只命中这两枚定义本身**（`.scratch/wisp/probes/62/v1/m1/ball_windows.go` 里那两行是别人留下的探针副本，不是产码）。〔量 R16〕
  ⇒ 后果逐字：今天托盘菜单上"静音／暂停唤醒"那两枚勾在产码里**永远是 `false`**，"推一维状态进菜单"这一形**有名字、没有生产使用者**，把它当先例用之前要认清这一点。

**审批状态住在哪一侧、依赖箭头现在朝哪。**

- 状态侧（`cmd/wisp`）：`ra.cards` 是 `*approval.Replies`；"有没有卡在等"今天**已经在同一侧被读过三次**：
  `cmd/wisp/resident_approval_windows.go:172`（`vetoByEsc` 内）、`:484`（`settleOrb` 内）、`cmd/wisp/approval_always.go:181`。〔量 R17〕
  读法＝`AwaitingHuman()`（`internal/agent/approval/replies.go:252`）。
- 球侧（`internal/ball`）：**不持有任何审批对象**。现有的注入方向是 `attachBall`（`resident_approval_windows.go:156` 调 `ra.ui.attachBall(rb.b)`；定义 `:387-391`；读用 `currentBall()` `:405-409`）＝**审批这一侧拿球**，反向今天不存在。〔量 R18〕
- 现成的反向先例＝**函数值注入**：`escVetoFunc`（`cmd/wisp/resident_ball_windows.go:61`，注释 `:56-60` 逐字「the whole of what the ball host knows about the thing that happens when the cancel key is pressed」），
  由 `startResidentBall(reg *observe.Registry, onCancelEsc escVetoFunc)`（`:111`）收，并在 `Events` 字面量里挂上（`:126` `OnCancelHotkey: func() { recordCancelHotkey(onCancelEsc) }`）。〔量 R7〕

**要加"有卡吗"需要动的文件／函数（只列代价，不选形）**

| 路线 | 要动谁 | 代价（读出来的，不是估的） |
|---|---|---|
| 甲 推（`SetTrayChecks` 那一形） | ①`internal/ball/ball_windows.go`：`SetTrayChecks` 签名或新增一枚 setter＋一枚 `Ball` 字段；②`tray_windows.go:72` `showMenu` 多收一维；③`cmd/wisp` 某处必须在卡片进出时**调用**它（今天的候选是 `ballCardUI.Update`／`settleOrb`，`resident_approval_windows.go:463-497`） | 状态由别人推 ⇒ **推的那一次没跑＝菜单上是旧值**；且甲形在生产里今天无人使用（上面那条零调用者读数），等于先立一根没人拉的线 |
| 乙 拉（构造菜单时同步读一枚注入的只读闭包） | ①`internal/ball`：`showMenu` 需要一个读者（放 `Options`／`Events`／还是参数＝**三选一都是契约面**，本腿不裁）；②`ball_windows.go:674` 调用点；③`cmd/wisp`：装配根（`resident_windows.go`）把读 `ra.cards.AwaitingHuman()` 的闭包交给 `startResidentBall`（`resident_ball_windows.go:111` 要扩参） | 读发生在 `ui-sta` 上；锁序见下面那段——本腿**读到的是"不阻塞"的形状**，但**没量** |
| 丙 ＝乙 的时机变体 | 同上，只是把这次读放在 `appendItem` 之前而非之后 | 与乙同价；`showMenu` 里**建菜单在 `TrackPopupMenu` 之前**（`tray_windows.go:73-91` 建、`:101-102` 弹），所以"打开菜单时现取"与"构造时取"在今天的代码里**是同一次调用内的同一时刻**，不构成第三枚独立形状 |

**乙形那把锁的实际读数（〔推〕，未量）**：`AwaitingHuman()` → `Pending()`（`replies.go:227-240`，`r.mu` 内只做 map/slice 拷贝）＋ `queueAwaitsHuman()`（`:306-312`，先 `routes()` 取 `r.mu` 再放，再调 `Queue().Depth()`）；
`Queue.Depth()`（`internal/agent/approval/queue.go:133-137`）持 `q.mu` 只 `return len(q.pending)`；
唯一持 `q.mu` 做 I/O 形状的是 `deliver`（`:298-310`），它里面的发送是 `select { case it.answer <- a: default: }`＝**非阻塞**。
`ballCardUI` 自己的 `u.mu`（`resident_approval_windows.go:387-411`）临界区里也只有赋值。
⇒ 〔推〕**读到的形状是"没有任何一条路径在持锁时等 `ui-sta`"**；但本腿没跑任何测试／构建，所以这一句不是量出来的，且它只覆盖我打开过的那几枚函数（`queue.go` 里 `q.mu.Lock` 还有 `:409/:423/:446/:515/:537/:553` 等我**没有逐枚读体**）⇒ 归 §⑦ F2。

---

## ② 三枚候选的射程与撞钉

### (a) 永远显示、按下去若无卡就报错

**`Replies.Allow` 在没有待决卡时的返回（逐字，`internal/agent/approval/replies.go:316-335`）〔量 R4〕：**

| 触发条件 | 返回 | 文案逐字 | 出处 |
|---|---|---|---|
| `Look(corr)` 不中（**托盘今天没有 corr 可传**；传空串同样不中，`replies.go:199-206` 是 map 查） | `ErrNoTrackedCard` | 「approval: 本机账上没有这张卡（从未显示、已答复，或已经离开屏幕）」 | `:57`／支在 `:319` |
| 卡在册但 `Grant == ""`（＝一枚 **L1 窗口**） | `ErrRouteHasNoAllow` | 「approval: 这一路线只有否决、没有允许」 | `:60`／支在 `:322` |
| gate 未 `Attach` | `ErrNoGateAttached` | 「approval: 答复面未绑定审批 gate，任何决定都送不出去」 | `:64`／支在 `:326` |
| 过路由后队列侧拒 | 透传 `g.DecideFromNative`（`gate.go:718-724`：`q.allow(corr, grant)`） | `ErrUnknownCorrelation`「correlation_id 无对应待审批项」`ui.go:124`／`ErrNotPending`「该审批项已结束，不再接受决定」`:126`／`ErrBadGrant`「原生令牌无效（缺失/已用/与本次请求不绑定）」`:130` | `queue.go:366/:370/:380` |

**文案会不会与"门拒绝"的文案复用——本腿量到的是"三族不相交，但 approval 族内部有一处近碰撞"〔量＋推〕：**
- 审批面：全部带 `approval:` 前缀（上表）。
- **门拦下**（无门兜底，`NoGate`）：`internal/tools/gate.go:118`「L1 确认窗口不可用（票 21 未接入），已拒绝执行」、`:126`「L2 审批通道不可用（票 21 未接入），已拒绝执行」、`:139`/`:144` 同族「尚未接入」两支。
- **权限判定**：`internal/perm/store.go:179`「不是三档之一，拒绝切换（R20/M1：没有第四档）」、`:203`「perm: 切到 auto_approve 被拒绝: %w」。
⇒ 〔量〕三族词面**不相交**（`approval:` vs `…已拒绝执行` vs `perm:`），今天**不存在**"门拦下的文案复用权限判定的文案"这一形状；本腿**没有**打开 `internal/perm` 的判定主路径全文，所以这是"我看到的那几处不相交"，不是"全族已穷举核过"⇒ §⑥ E13。
⚠ 但 **approval 族内部**有两枚几乎同义的文案：`ErrNoTrackedCard`（`:57`）与 `ErrUnknownCorrelation`（`ui.go:124`），
前者＝"本机账上就没这张卡"、后者＝"队列里没有这个 correlation"。(a) 那一形在无卡时报前者、在卡刚被别人答过时报后者，**托盘用户看不出差别**（控制台之所以能分辨是因为它自己另写了两句人话并把两种含义分开，见下面 §③ 的 `approval_reply.go:213-229`）。

**要动谁**：`internal/ball/tray_windows.go`（第五枚 id `menuAllowOnce`，`:20-24` 现在只到 `menuExit = 4`）＋ `internal/ball/ball_windows.go:675-684` 的 `switch sel` 加一支 ＋ 一枚把结果带回 `cmd/wisp` 的通道（`Events` 新键 or 注入闭包，见 §④）＋ `cmd/wisp` 侧真正调 `Allow`。
**今天有什么仪器看得见＝零**：`showMenu`／`wmAppTray`／`OnTrayExit`／`recordTrayExit`／`menuExit` 在全仓 `*_test.go` 里**零命中**
〔量 R20：`grep -rn "showMenu" --include=*_test.go .` ＝ 0；正控＝同一条尺去掉 `--include` 命中 `internal/ball/ball_windows.go:674` 与 `internal/ball/tray_windows.go:72`〕。
**但 `Allow` 无卡那一支今天已被钉住**〔量 R21〕：`internal/agent/approval/replies_201_test.go:93-94`（`Allow(ctx,"never-displayed")` 必为 `ErrNoTrackedCard`）与 `:101-102`（L1 窗口的 `Allow` 必为 `ErrRouteHasNoAllow`）⇒ **(a) 不是"无人看守的语义"，只是"无人看守的那半在托盘"**。
**触碰哪枚冻结契约**：C18（`docs/PLAN.md:1368`）里那句逐字「**「允许」决策只接受原生侧来源（D33/F2，面板来源直接拒绝）**」＋ D43/C12 的 `Confirming`(L1)／`AwaitingApproval`(L2) 两态配对（`PLAN.md:3077` 行、`replies.go:286-302` 的 `WaitingState`）。托盘是**原生侧**⇒ 不改宽 C18；⛔ 若为了"让报错好写"去动 `DecideFromPanel`/`PanelAllow`（`gate.go:730-743`、`replies.go:431-441`）就是撞票面 AC#5（`:42`）那条硬禁。

### (b) 灰掉／不显示

**必须先拿到 ① 那一跳**（无卡与否是画在菜单上的前置），代价见 ① 表。
- **置灰那一支今天做不到"零新增"**：`appendItem` 只认 `checked`（`tray_windows.go:79-85`，只 OR 进 `mfCheckd`），
  而包内 menu flag 常量只有两枚：`internal/ball/win32_windows.go:136 mfCheckd = 0x00000008`／`:137 mfSepart = 0x00000800`；
  `mfGrayed`／`MF_GRAYED`／`mfDisabled` **全仓零命中**〔量 R12，正控＝同尺能命中 `mfCheckd`/`mfSepart`〕⇒ 要置灰就得**新增一枚 Win32 flag 常量**（`win32_windows.go`）＋扩 `appendItem` 参数。
- **不显示那一支不需要新常量**：把 `appendItem(...)` 那行套一层判断即可（`tray_windows.go:86-91` 现形状就是逐条 append）。
- 仪器：**零**（同 (a)，尺与正控同 R20）。
- 冻结契约：**不碰 C18 的文字**（按钮可不可见不是"允许的来源"那一维）；但它**直接改变"人能不能行使原生侧允许"这件事的可见性**，
  ⚠ 而 SPEC-08 §7（`docs/specs/SPEC-08-ui-ball-panel.md:229`）逐字只列了**四枚**右键项（「打开面板/静音/暂停唤醒/退出」）⇒ **第五枚无论用哪一形都是对该 spec 行的增补**（票 228 AC#2 已把它当票内活；本腿只登记这一句的存在，不裁它算不算 spec 变更）。
- 失败预演式风险（读码，非量）：置灰／不显示在 **L1-only** 情况下会让人读成"系统坏了"，因为同一张菜单里没有任何一句说明；见 §⑦ F4。

### (c) 只有当菜单本身被打开时才现取一次状态

**`TrackPopupMenu` 这一支的同步语义（本腿现读）〔量 R6〕**：`tray_windows.go:101-102`
`sel, _, _ := pTrackPopupMenu.Call(menu, tpmRightButton|tpmReturnCMD|tpmNoActivate, …)` ⇒
①调用**阻塞到用户选完或关掉**；②选择值**只从返回值**拿；③文件头 `:15-16` 逐字写明「wParam of WM_COMMAND after TPM_RETURNCMD is **NOT** used; we use TrackPopupMenu's return value inline in the wndproc instead」；
④`:70-71` 逐字「Returns the selected command id (**0 = dismissed**)」⇒ **"被关掉"与"什么都没选"在返回层面是同一个 0**，托盘侧今天无法区分。
- ⚠ 这条与票面 `:100` 完全一致（它逐字关掉「从测试里投一枚 `WM_COMMAND`＋菜单 id 去模拟托盘选择」这一支，并指名凭据就是 `:15-16` 与 `:101-106`）⇒ 本腿复认**行号未漂**。
- **要动谁**：只比 (a)/(b) 少一样——**不需要** `SetTrayChecks` 那条推线；需要 `showMenu` 内有一个读者（＝① 的乙形注入）。
  关键顺序事实：菜单是在**同一次 `showMenu` 调用内、弹出之前**建的（`:73-91` 建 → `:101-102` 弹）⇒ "打开时才取"＝"在 `ui-sta` 上、`wmAppTray` 的右键分支里取一次"，**不是一次跨线程往返**。
- 仪器：**零**（同 (a)(b)）。
- 冻结契约：与 (a)/(b) 同一枚 C18；(c) 本身**不新增决策面**，只新增一次读，所以它是否碰契约取决于它承载 (a) 还是 (b)。

### 第二发补充：AC#2 那句"带正确 `Channel` 的答复"今天**没有对应的值可带**〔量 R31〕

尺：`sed -n '44,125p' internal/agent/approval/approval.go` ＋ `grep -rn "Statuses()\|ChannelStatus" --include=*.go internal cmd | grep -v _test`。

- `Channel` 是**四枚否决通道的闭集**，注释 `approval.go:48-50` 逐字：「The Channel is one of the four L1 **veto** channels. **The set is closed: a fifth value means a caller invented a selector, which is exactly the M-7/C-3 failure shape（a security decision keyed on an unvalidated label）**」；
  四枚值 `:54-57`（`ball`/`esc`/`panel`/`kws`），`allChannels :61`，显示名 `channelNames :88-93`。
- `Veto` 结构体才带 `Channel`（`:73-77`）；**"允许"这一侧没有 `Channel` 字段**：`Allow` 走 `Request{CorrelationID, Allow, Grant, Source}`（`replies.go:328-330`），`Source` 是**标签不是权威**（`replies.go:116-123` 逐字「They are LABELS, not authority (gate.go says no branch reads Request.Source)」，`gate.go:719` 把它记成 `claimed_source=%q`）。
- ⇒ **后果逐字**：票面 AC#2（`:39`）那句"每枚按钮断『点下去 ⇒ `Gate` 收到一条带正确 `Channel` 的答复』"，对**「允许一次／拒绝／长期允许」这三枚**今天**无法按字面满足**——要它满足就得给闭集加第五枚值，而那一句注释正把这件事点名成 M-7/C-3 形状。
  可满足的读法只有两支：**甲**把"正确 Channel"读成 `HostBinding.NativeSource` 标签（现成面，`replies.go:127-135`／`resident_approval_windows.go` 的 `nativeReplySource` 一类）；**乙**把它读成"这枚按钮只出现在 `Veto` 那一侧"（则 AC#2 里那三枚 allow 按钮这一判据今天**结构性写不出来**）。**两支都由编排者裁，本腿不选**（⇒ §⑦ F8）。

---

## ③ 诚实 advertise 的先例与同形判断

**先例一（热键）**：`HotkeyStandby`（`internal/ball/hotkey_windows.go:259`）——注释 `:250-254` 逐字「NOT attempted, by design - ticket 245」；
`Problems()`（`:361-381`）在 `:365` 把 `HotkeyLive, HotkeyDisabled, HotkeyStandby` 一起 `continue`（＝**不报警**），
但同一枚分类**并没有被藏起来**：`:275-279` 给它一句自己的状态话（「not bound while idle (cancel is borrowed only during Confirming)」），
且 `Attempted()`（`:284-287`）把"试过／没试过"分开留成可编程的事实；常驻腿还把总数打出来（`cmd/wisp/resident_ball_windows.go:156-157` 逐字 `hotkeys live %d/4`）。
注释 `:357` 逐字「Standby is deliberately NOT a problem line (ticket 245)」。**刻意不绑＋说清是哪一类＋不报警**。

**先例二（无通道就说无通道）**：`interactiveStdin()`（`cmd/wisp/approval_reply_stdin_windows.go:41`）无控制台／是管道／是文件时 `:44`、`:49` 两处 **`return nil`**；
文件注释 `:7-13` 逐字「gets nil, which means "no answer source" … the card is shown, nobody answers it」，`:30-33` 逐字「a run whose stdin is a pipe cannot be answered, and **says so at startup rather than going quiet**. That is the same shape as every other degraded boot in this file set（**the downgrade has to be visible**）」。
它的两个消费者都把那句话**说出口**：`cmd/wisp/main.go:154-158`（`reply == nil` 时向 stderr 逐字「本机没有可交互控制台，本轮没有人能答复卡片：L2 卡会等到超时后按拒绝处理，L1 窗口没有人能否决」）；
`cmd/wisp/resident_task_source_windows.go:224-230`（`console == nil && injectedText == ""` ⇒ `slog.Warn` 带 `why`/`effect` 两句＋`fmt.Printf` 出 `taskEntryDisabledClaim`，而那枚常量 `:85-88` 的注释逐字「**names the consequence, not just the condition**」）。

**先例三（按下去没卡时的那句话，最贴近本格）**：`vetoByEsc`（`cmd/wisp/resident_approval_windows.go:171`）在 `:173-177` 无卡时**返回一句实话**，逐字
「按下的取消键没有可否决的确认项：本进程此刻没有卡片在等人」，注释 `:174-176` 逐字「**Said out loud rather than left silent**: "the key was taken and there was nothing to cancel" is the reading that tells a user the borrow and the card are two different lifetimes」。
它的另一面是 `bindBallHost`（`:139-158`）：装配根没注入取消执行者时**拒绝**把 Esc 登记成可用（`:150` 逐字「**advertise 一枚按不动的键＝B1 禁止的形状**」），并且那一句是打出来的（`residentAuditf`）。

**先例四（报错要分含义）**：`replySurface.allow`（`cmd/wisp/approval_reply.go:213-236`）拿到 seam 的错误后**没有原样抛出去**，而是 `errors.Is` 分流成**三句不同人话**并且各记一行审计：
`ErrNoTrackedCard`→「没有找到待答复的卡片 %q：它可能已经结束，从未显示的卡片这里也没有令牌可花」＋`record("REFUSED", …, "本机没有这张卡的记录")`；
`ErrRouteHasNoAllow`→「%s 是执行前阻止窗口，它只有否决、没有允许（SPEC-06 §2 B1）；要停下它就答 veto %s」＋`record(…, "该路线没有允许动作")`；
`ErrBadGrant`→再补一句「这张卡需要重新显示才能被允许」。
⇒ 这就是本仓**已经付过一次代价**的那件事：**两种"不许"不许共用一句话**。

**先例五（同一枚托盘文件里已经有一枚"什么都不做但说三遍"的项）**：`recordTrayExit`（`cmd/wisp/resident_ball_windows.go:253-258`）
今天的形状＝①`const why` 一句完整因果（逐字「this process leaves when its event loop returns, which today only a console signal (Ctrl+C) asks for; the tray Exit item has no stop path attached to it」）
＋②`slog.Warn("tray exit requested", "outcome", "ignored", "why", why)`＋③`fmt.Printf` 同一句到 stdout；
注释 `:247-252` 逐字「**so it is reported instead of assumed**」，并把两条补法点名成新契约（与票面 `:89`/`:93` 一致）。〔量 R32〕
⇒ **这是全仓离"托盘按下去没有执行者"最近的一枚诚实先例**：它不隐藏、不假绿，也不静默。

**先例六（可用性是随卡带出去的字段，不是"不画"）**：`ChannelStatus`（`approval.go:80-85`，注释逐字「the honest, user-visible availability line for one channel. **An unloaded channel is NEVER rendered as available (B1)**」）
＋`unavailableText`（`:102-119`，四句各不同文案，KWS 那句是 spec 逐字要求的「语音取消不可用」，注释逐字「the strip **must say this, never stay silent** about a channel that cannot fire」）
＋`ErrChannelUnavailable`（`:121-124`，注释逐字「It is returned, not swallowed: **a silently ignored cancel attempt is how a fake channel becomes a user-visible promise**」）
＋`Statuses()`（`:190-194`）被 `gate.go:605` 逐字 stamp 进卡片的 `Channels` 字段（结构在 `ui.go:40`）。〔量 R33〕
⇒ 这一族的形状＝**能力自己带着"我此刻能不能用"走**，界面按 `Loaded` 画一句不可用的话，而不是把能力藏起来。

**先例七（"投了一票但没人接"已被裁过一次）**：`gate.go:456-463` 逐字记录票 87 那一发：
「a veto naming a card that is ON SCREEN used to fall through to `ErrUnknownCorrelation` here … **The user's vote was simply not heard and the call waited out the full C18 deadline**. Hand it to the queue's refusal funnel instead」。〔量 R33〕
⇒ 本仓对"按下去落空"**已有的判决不是"让它安静"，而是"改路由，让落空本身产生一条看得见的拒绝"**。这一条与 §②(a) 的关系最直接：托盘那枚"允许一次"若无卡，本仓的先例要求它**留下一条可审计的落空**，而不是让菜单 0／被关掉那一支替它背。

**共同形状（一两句）**：把"不响"分成**刻意不绑**（分类有名、不报警）与**按了落空**（必须出声、且出声的句子要说清是哪一种落空）；
**能力入口在拿不到真通道时返回 nil／一句后果句，绝不返回一个看起来能用的假象**（`advertise 一枚按不动的键` 是被点名禁止的那一枚）。

- **同形／相反（只报形状匹配，不报推荐）**：
- **(a) 永远显示＋按下去报错**：与先例三／先例四／先例二**同形**（按下 ⇒ 一句分过含义的实话），
  前提是它真的**把那句说出来**；若实现成"取到错误但没人接"（`internal/ball` 的回调是 `func()`，`b.fire` 只看 nil／非 nil，`ball_windows.go:728-732`），它就同时**与先例相反**——这正是派单点名的"静默失败"那一格。⇒ 本腿把 (a) 拆成 (a1) 出声／(a2) 不出声两支登记，见 §⑥ E14。
- **(b) 灰掉／不显示**：**灰掉**与先例一同形（看得见"这里有一枚能力，此刻不可用"），
  **不显示**与先例**相反**——先例四支全都在"说清是哪一类不响"，而"整枚项消失"不留任何一句，读的人无法区分"这里从来没有这枚按钮"与"有卡但我看不见"。
- **(c) 打开菜单时现取一次**：**同形**于先例一（在展示的那一刻按当下的真实分类决定可见性，且 `HotkeyStandby` 那条先例证明"分类本身要说得出"）；
  但它同时与先例二那条"分类要说出口"的**另一半**有张力：`showMenu` 的返回值里"被关掉"与"没选"都是 0（`:70-71`），所以取数时机再准，**选择这一侧今天仍只有"命令 id"这一维可返回**，报错那一维要靠别处（审计／球）出声。
- **补三条与先例五／六／七的匹配（第二发）**：
  先例五（`recordTrayExit` 的"说三遍"）与 **(a1) 同形**、与 **(a2) 相反**、与 **(b) 的"整枚消失"相反**（它那一支什么都没显示过）。
  先例六（能力自带 `Loaded` 句、"NEVER rendered as available"）与 **(b) 的置灰支同形**、与 **(b) 的不显示支相反**、与 **(a) 同形**（(a) 正是"带着不可用那句话被按下"）；⚠ 但先例六今天**只覆盖 `Veto` 那一侧**（`ChannelStatus` 只有四枚通道，`Allow` 没有对应的 availability 面，见 §② 第二发补充）⇒ 照它的形要**新增一枚 availability 概念**，那是新面不是抄面。
  先例七（票 87 把"落空的一票"改路由成一条看得见的拒绝）与 **(a1) 同形**；⚠ 它**不与 (b) 冲突**（(b) 从一开始不让那一票发生），这一点我不能替编排者判"哪种算落空"，归 §⑦ F4。

---

## ④ 会被叫红的钉（逐枚读断言原文）

> 按派单禁令：**本节不引 `cmd/wisp/*_test.go` 的行号**（写腿 `33-r4` 正在改那批文件，见 §⑤ R23）；只写**名册与断言的语义**，**行号以落地时现读为准**。
> 本节读的原文出自 `cmd/wisp/resident_ball_228_test.go`（函数名与逐字句子是本腿读的）。

- **钉 P1＝名册数超过本文件名册 ⇒ 红**。所在用例：`TestAC228BallHostAnswersEveryGesture`。
  名册是一枚**字面列表** `ballEventCallbacks228`，十枚键名逐字：
  `OnClickBall / OnSummonHotkey / OnMuteHotkey / OnCancelHotkey / OnPanelHotkey / OnTrayPanel / OnTrayMute / OnTrayPauseWake / OnTrayExit / OnDragEnd`。
  注释逐字（为什么是字面列表）：「when internal/ball grows an eleventh gesture, the case below goes red and whoever added it has to say what the resident leg does with it. A nil Events entry is a gesture that vanishes without a word, which is the same silence this family keeps counting.」
  断言逐字（红的时候说什么）：「AC#1 RED (the instrument, not the code): the host sets %d callbacks but this file lists %d; **internal/ball.Events grew and ballEventCallbacks228 did not.**」
  取数方式＝AST 走形：只把球宿主函数体里 `*ast.CompositeLit` 且类型为 `…Events` 的字面量的 **key** 收进集合；`len(set) > len(名册)` 即红。
- **钉 P2＝某枚回调留 nil ⇒ 红**（同一枚用例的另一条断言）。
  断言逐字：「AC#1 RED: the resident ball host leaves %d of internal/ball's %d gesture callbacks unset: %v. **A nil callback is a click or hot key that reaches no one and logs nothing.**」
  注释逐字：「every consumer callback internal/ball offers is given a non-nil executor by the resident host, because a nil entry is a click or a hot key that disappears without a log record」。
- **相邻的两枚同族钉（本腿一并读到，语义登记）**：①「no Windows build of `runResident` found among this …（the instrument, not the code）」＝球宿主必须是常驻腿本身；
  ②「`runResident` defers no ball teardown, so the hot keys and the tray icon it installed are left to process death」＝**球与托盘的收尾必须挂在 `runResident` 的 defer 上**（它句里直接指 `internal/proc/shutdown.go` 的 D38(e) 第二步）。

**三枚候选下谁响谁不响（逐枚）**：

| 候选 | P1（名册数） | P2（留 nil） | 说明（读断言取数方式得出） |
|---|---|---|---|
| (a) 永远显示＋按下去报错 | **响**（若走"新增 `Events` 键并在 `startResidentBall` 的字面量里设上"这一形，`len(set)` 变 11 > 10 ⇒ 逐字「the instrument, not the code」） | **不响**（十枚一枚没少） | ⚠ 反过来：如果新增 `Events` 键而宿主**不设上**它，P1 不响（`set` 只数被设上的）、P2 也不响（它只查名册里那十枚）⇒ **"第五枚手势静悄悄"这一维恰好是这两枚钉都照不到的**，与票面 `:102` 那句逐字读数同口径（「名册不响**不等于**有行为」） |
| (b) 灰掉／不显示 | **同上**（要不要新增键是同一问） | 不响 | "不显示"这一支若靠 `internal/ball` 自己读状态（不经 `Events`），P1/P2 **都不响**，代价是破掉 `resident_ball_windows.go:56-61`/`:23` 那句「球宿主对审批一无所知」的分层（该文件 `:23` 逐字「every gesture the ball …」＋`:195-203` 的 `ballGestureWhy`）⇒ **名册钉与分层纯度是互斥的两条**，归 §⑦ F7 交裁 |
| (c) 打开时现取 | **同上**（取数时机不新增手势） | 不响 | (c) 只改**读的时刻**，不改 `Events` 名册；它与"是否新增键"是正交的 ⇒ 单看 (c)，两枚钉**一声不吭** |

**`internal/ball` 那一侧的 11 枚热键状态用例与 AST 走形用例——它管不管这一格（逐枚现读）**：
- 〔量 R22〕`grep -c "^func Test" internal/ball/hotkey_status_test.go` ＝ **11**（另 `hotkey_live_test.go` 4、`hotkey_test.go` 2）。11 枚逐名读下来断的是：注册集合／默认热键空闲不持 Esc／借用往返／借用失败要成 problem line／`Standby` 不算 Disabled 也不算 problem／失败分族／`UnregisterAll`／状态字符串／reloader 与配置节／`applyHotkeyDefaults` ⇒ **没有一枚碰托盘、菜单或 `showMenu`**。
  再配一把尺〔量 R20 同族〕：`grep -rn "tray\|Tray\|showMenu" internal/ball/*_test.go` ＝ **零命中**（正控＝同尺在产码命中 `tray_windows.go:72`）⇒ **11 枚热键用例不管这一格**，三枚候选都不会让它们红。
- AST 走形那一族在 `internal/ball` 只有一枚文件用（`internal/ball/tokens_table_test.go`，头注释逐字：它钉的是 **C21 设计 token 表**与代码的双向漂移）。⇒ **与本格无关**（它不枚举 `Events`、不读托盘）。
- 真机那一族：`internal/ball/interaction_live_test.go:5-13` 的总纲逐字「Where a clause needs a human hand that no harness can supply (a second monitor, the felt experience of focus) the test says so in this file rather than faking a proof」。
  ⚠ 但**票面与 `228-a2` 引的那句具体措辞（"needs a real mouse in the notification area"）在本腿的 HEAD 上量不到**：`grep -rn "notification area\|real mouse" internal/ball/*_test.go` ＝ **零命中**〔量 R24〕⇒ 见 §⑥ E12。
- **结论（本格最重要的读数）**：三枚候选**共同的仪器现状＝托盘侧零枚仪器**；唯一看守得住的语义在 `internal/agent/approval` 那一侧（`replies_201_test.go:93-102`，R21）。
  ⇒ 谁落地这一格而不自带新钉，落的就是"票面 AC#2 后半句『回调接反／接空必须能判红』今天没有任何东西在守"这一格。⚠ 但新增钉的形状**不许是** `WM_COMMAND` 注入（票面 `:100` 已逐字关掉），本腿只登记，不补缺。

---

## ⑤ 我跑了哪些尺、每条真实读数

> 口径：〔量〕＝本腿在这台机器上**真跑到**的读数（Read／Grep／`git cat-file`），〔推〕＝读码推出、没量的。
> ⚠ 按派单禁令：本腿**没跑**任何 `go build`／`go vet`／`go test`／`scripts/build.ps1`／`go mod tidy`，所以 ⑤ 里**零枚编译或测试读数**，全为静态读数。

| # | 尺（命令／动作，逐字） | 真实读数 |
|---|---|---|
| R1 | `date "+%Y-%m-%d %H:%M:%S %z"` 与 `git log -1 --format=%H` 同发 | `2026-10-01 11:27:59 +0800` ／ `eed229e48cb9599c29f72fea6a1170739a5f2b41` |
| R2 | `git status --porcelain=v1 -b`（只看，不动） | 首行 `## dev...cnb/dev [ahead 169]`；脏项含 ` M .gitignore`、` D design/**`（16 枚）、`.scratch/**` 未跟踪件——**一律未动、未提交、未还原** |
| R3 | `grep -n "func (r \*Replies)\|ErrRouteHasNoAllow\|ErrNoPending" internal/agent/approval/replies.go` | `Allow` 在 `:316`、`AllowSession :351`、`Reject :372`、`PanelReject :378`、`PanelAllow :431`、`Veto :455`；`AwaitingHuman :252`、`Pending :227`、`WaitingState :286`、`queueAwaitsHuman :306`；`ErrRouteHasNoAllow` 定义在 `:60`，被 `:322` 与 `:357` 两处使用；**`ErrNoPending` 零命中（不存在这枚错误名）** |
| R4 | `Read internal/agent/approval/replies.go` `:53-65` 与 `:314-335` | `Allow` 的返回路径四支出声：`:319` `ErrNoTrackedCard`（`Look(corr)` 不中）、`:322` `ErrRouteHasNoAllow`（`card.Grant == ""`）、`:326` `ErrNoGateAttached`（`g == nil`）、`:331` 透传 `g.DecideFromNative` 的错；只有 `:333` `r.Forget(corr)` 之后 `:334` 才 `return nil`。三枚文案逐字：`ErrNoTrackedCard`＝「approval: 本机账上没有这张卡（从未显示、已答复，或已经离开屏幕）」（`:57`）、`ErrRouteHasNoAllow`＝「approval: 这一路线只有否决、没有允许」（`:60`）、`ErrNoGateAttached`＝「approval: 答复面未绑定审批 gate，任何决定都送不出去」（`:64`） |
| R5 | `grep -n "showMenu\|wmAppTray\|recordTrayExit\|muted\|pausedWake\|menuExit" internal/ball/ball_windows.go` | `case wmAppTray:` 在 `:669`；`sel := showMenu(b.hwnd, b.trayMuted, b.trayPauseWake)` 在 `:674`；分流 `menuOpenPanel :676`／`menuMute :678`／`menuPauseWake :680`／`menuExit :682`；`func (b *Ball) SetTrayChecks(muted, pausedWake bool)` 在 `:919`，`:921-922` 写两枚字段。**⇒ 传进菜单构造的状态＝两枚布尔，没有第三枚** |
| R6 | `Read internal/ball/tray_windows.go` 全文（110 行） | 签名 `:72 func showMenu(hwnd windows.HWND, muted, pausedWake bool) uint32`；id 常量 `:20-24`＝`menuOpenPanel 1 / menuMute 2 / menuPauseWake 3 / menuExit 4`（**第五枚 id 今天不存在**）；`:79-85` 局部 `appendItem(id, label, checked)` 只会加 `mfString`／`mfCheckd` **两枚 flag，没有置灰那一支**；`:86-91` 五项（含两枚分隔条）；`:101-102` `pTrackPopupMenu.Call(menu, tpmRightButton|tpmReturnCMD|tpmNoActivate, …)` **同步返回选择值**；`:105` 补一枚 `WM_NULL`；`:15-16` 注释逐字写明"不用 `WM_COMMAND`、选择值就地从 `TrackPopupMenu` 返回值拿" |
| R7 | `grep -n "recordTrayExit\|escVetoFunc\|SetTrayChecks\|OnTrayExit\|go func" cmd/wisp/resident_ball_windows.go` | `type escVetoFunc func() string` 在 `:61`；`startResidentBall(reg, onCancelEsc escVetoFunc)` 在 `:111`；`Events` 名册 `:123-132`（`OnTrayExit: recordTrayExit` 在 `:131`）；`recordTrayExit` 函数体在 `:253`。**`go func` 在该文件零命中**（本腿只读数，不新增裸 `go`） |
| R8 | `grep -n "HotkeyStandby\|Problems\|Attempted" internal/ball/hotkey_windows.go` | `HotkeyStandby` 定义 `:259`（注释 `:250-254` 逐字"NOT attempted, by design - ticket 245"）；`Attempted()` 在 `:284-287`；`func (r HotkeyReport) Problems()` 在 `:361`，`:365` 把 `HotkeyLive, HotkeyDisabled, HotkeyStandby` 三枚一起当"不是问题"；`:357` 注释逐字「Standby is deliberately NOT a problem line (ticket 245)」；`cancelIdleLine` 在 `:428`、`:441` 设 `Status = HotkeyStandby` |
| R9 | `grep -n "consoleVetoChannel\|func \|return nil\|stdin" cmd/wisp/approval_reply_stdin_windows.go` | `func interactiveStdin() io.R`（原文 `io.Reader`）在 `:41`，`:44` 与 `:49` 两处 `return nil`（**没有可用控制台就返回 nil，不假装有通道**）；`:17-20` 注释逐字把"redirected stdin 是别人的字节"列为不许读的理由。**本腿读到的这一枚文件里 `consoleVetoChannel` 这个符号名词面零命中**——见 §⑥ E3，派单给的名字与盘上现名不一致 |
| R10 | `grep -rn "Replies.Allow\|\.Allow(" cmd/ internal/` 的调用者清点（只看引用面） | 见 §②(a) 正文；〔量〕＝这一发在 `cmd/wisp` 生产码里**没有托盘/球侧的 `Allow` 调用者**（`recordTrayExit` 那一支只打日志，`resident_ball_windows.go:253`） |
| R11 | 票 228 的 AC#2 与 AC#11 两格现读（`Read` 行 86-130 与 `grep -n "^#..."`） | AC#2 在票面 `:39`（逐字：「允许一次」→`Replies.Allow`（票面写 `replies.go:313`，**本腿现读 `:316`**，见 §⑥ E1）；判据形状＝"每枚按钮断『点下去 ⇒ Gate 收到一条带正确 Channel 的答复』，并把回调接反/接空必须能判红"）；AC#11 在票面 `:91-102`，其中 `:100` 逐字关掉"从测试里投 `WM_COMMAND`＋菜单 id 模拟托盘选择"这一支，并指名凭据＝`internal/ball/tray_windows.go:15-16` 与 `:101-106`（本腿 R6 复认这两处行号**未漂**）；`:102` 逐字回答 F6＝名册钉的措辞是"每个 gesture callback 都要有执行者"、**不是**"每个执行者都要有行为" |
| R12 | `grep -n "mfGrayed\|MF_GRAYED\|mfDisabled" internal/ball/*.go` | 读数见 §②(b) 正文（这一发的真实命中数**在下面补完时逐字登记**，若零命中即写零） |
| R13 | `git diff --name-only` 起手名册核对 | 本腿只写这一枚文件；其余脏项一律不属于本腿（见 §⑤ R2） |
| R14 | `grep -rn "mfGrayed\|MF_GRAYED\|mfDisabled" internal/ball/` ＋正控 | 第一发**零命中**；正控尺命中 `internal/ball/win32_windows.go:136 mfCheckd = 0x00000008`、`:137 mfSepart = 0x00000800` ⇒ **包内没有置灰 flag**〔量〕 |
| R15 | `grep -n "type Events struct" -A 22 internal/ball/ball_windows.go`；`grep -n "func (b \*Ball) fire" -A 12` | `Events` 定义在 `:49-60`，**十枚键**（`OnClickBall … OnDragEnd`）；`fire` 在 `:728-732`，体只有 `if fn != nil { fn() }`＝**回调签名 `func()`，没有返回值、没有错误通道**〔量〕 |
| R16 | `grep -rn "SetTrayChecks\|SetTrayTip" --include=*.go cmd internal` | **只命中定义本身**：`internal/ball/ball_windows.go:918`（注释）／`:919`／`:926`（注释）／`:927` ⇒ 两枚托盘 setter 的**生产调用者＝0 枚**〔量，正文 §① 加重〕 |
| R17 | `grep -rn "AwaitingHuman()" --include=*.go internal cmd` | 产码调用者三枚：`cmd/wisp/approval_always.go:181`、`cmd/wisp/resident_approval_windows.go:172`、`cmd/wisp/resident_approval_windows.go:484`；定义 `internal/agent/approval/replies.go:252`；内部自用 `:290`。其余命中在 `*_test.go`（`replies_201_test.go`／`approval_seam_201_test.go`／`resident_approval_246_windows_test.go`／`resident_approval_live_246_windows_test.go`）⇒ **"有卡吗"今天已在 `cmd/wisp` 被读，但没有一条线把它送进 `internal/ball`**〔量〕 |
| R18 | `grep -n "attachBall\|currentBall\|func (ra \*residentApproval)" cmd/wisp/resident_approval_windows.go` | `bindBallHost` `:139`、`ra.ui.attachBall(rb.b)` `:156`、`SetLoaded(ChannelEsc,true)` `:157`、`SetLoaded(...,false)` `:351`、`attachBall` `:387`、`currentBall` `:405` ⇒ 依赖箭头今天＝**审批侧持球**，反向不存在〔量〕 |
| R19 | `grep -rn "Err[A-Za-z]* = errors.New" internal/agent/approval/*.go`（排测试）＋`grep -n "func (q \*Queue) allow" -A 30` | 族谱：`ErrChannelUnavailable approval.go:121`／`ErrAlreadyStarted gate.go:402`／`replies.go:57/:60/:64`／`ui.go:124 ErrUnknownCorrelation`、`:126 ErrNotPending`、`:130 ErrBadGrant`、`:133 ErrPanelAllow`、`:136 ErrNotAdmitted`；`Queue.allow :343` → `allowScoped :361`，`:366/:370/:380` 三支出声、`:382` `deliver` 在**已 Unlock 之后**才调〔量〕 |
| R20 | 负向尺：`grep -rn "showMenu" --include=*_test.go .` ＋ `grep -rln "showMenu\|wmAppTray\|recordTrayExit\|menuExit\|SetTrayChecks" --include=*_test.go .`；正控：`grep -rn "showMenu" --include=*.go .` | 两把负向尺**都零命中**；正控命中 `internal/ball/ball_windows.go:674` 与 `internal/ball/tray_windows.go:72`（另 `.scratch/wisp/probes/62/v1/m1/ball_windows.go:675` 是探针副本，不算产码）〔量〕⇒ 另：`grep -rn "tray\|Tray\|showMenu" internal/ball/*_test.go` 亦**零命中** |
| R21 | `grep -rn "ErrNoTrackedCard\|ErrRouteHasNoAllow" --include=*_test.go .` | 命中 `internal/agent/approval/replies_201_test.go:93-94`（`Allow("never-displayed")` 必 `ErrNoTrackedCard`）与 `:101-102`（L1 窗口的 `Allow` 必 `ErrRouteHasNoAllow`），同文件 `:96-97` 钉 `Veto` 同族 ⇒ **无卡时报错这一语义已被看守**，只是守的不是托盘那一发〔量〕 |
| R22 | `grep -c "^func Test" internal/ball/hotkey_status_test.go internal/ball/hotkey_live_test.go internal/ball/hotkey_test.go` | **11 / 4 / 2**；再 `grep -rn "func Test" internal/ball/hotkey_status_test.go` 逐名读（`TestRegisterAllLiveSet`、`TestDefaultHotkeysIdlePassHoldsNoEsc`、`TestCancelBorrowRoundTrip`、`TestCancelBorrowFailureIsAProblemLine`、`TestStandbyIsNotDisabledAndNotAProblem`、`TestRegisterAllSplitsFailureFamilies`、`TestUnregisterAllDropsKnownSlots`、`TestHotkeyStatusString`、`TestHotkeyReloaderRebindsOnConfigChange`、`TestHotkeyReloaderSectionsHookIgnoresOtherSections`、`TestApplyHotkeyDefaults`）⇒ 11 枚**逐枚与托盘无关**（判语见 §④ 末）〔量〕 |
| R23 | `git log -3 --format="%h %ad %s" --date=short` ＋ `git show --name-only f4c583db` | 起手锚 `eed229e4` 之后盘上多了两枚：**编排者的台账发 `f4c583db`**（只动 `.scratch/wisp/issues/33-panel-host-c27.md` 与 `docs/reports/pending-and-issues.md`，**零产码**）与本腿自己的骨架发 `07f74694`；台账 A496 逐字把那枚写腿代号从 `33-r3` **改成 `33-r4`**〔量〕⇒ 见 §⑥ E15/E16 |
| R24 | `grep -rn "notification area\|real mouse" internal/ball/*_test.go cmd/wisp/*_test.go` | **零命中**〔量〕⇒ 票面 `:116` 与前腿 `228-a2` 引的那句措辞在本腿 HEAD 上量不到；本腿读到的是 `internal/ball/interaction_live_test.go:5-13` 那段总纲（"needs a human hand that no harness can supply"）⇒ §⑥ E12 |
| R25 | `grep -n "^| C17\|^| C18" docs/PLAN.md`、`grep -n "C18" docs/PLAN.md`、`grep -rn "托盘：右键菜单" docs/specs/`、`sed -n '3077p'` | `PLAN.md:1367`＝C17 `PanelBridge`；**`PLAN.md:1368`＝C18 `ApprovalQueue`，句末逐字「「允许」决策只接受原生侧来源（D33/F2，面板来源直接拒绝）」**；`PLAN.md:3077`＝D43 第 17 行 `Acting`→`Confirming`(L1)/`AwaitingApproval`(L2)；`docs/specs/SPEC-08-ui-ball-panel.md:229` 逐字「托盘：右键菜单（打开面板/静音/暂停唤醒/退出）；左键 = 打开面板。」（**只列四枚**）〔量，只读未改〕 |
| R26 | `grep -rn "errors.New(\"perm\|拒绝\|不允许\|无权" internal/perm/*.go internal/tools/gate.go`（排测试）＋`grep -n "func (g \*Gate) DecideFromNative" -A 40 internal/agent/approval/gate.go` | `gate.go:718-724` `DecideFromNative`：`!r.Allow` → `q.reject`，否则 `q.allow(corr, grant)`；`:730-743` `DecideFromPanel` 对 `r.Allow` 直接回 `ErrPanelAllow` 并 `revokeGrants`。门拦下文面：`internal/tools/gate.go:118/:126/:139/:144`（四句都是 `AnswerReject` ＋「…未接入…已拒绝执行」）；权限面：`internal/perm/store.go:179`、`:203` ⇒ 三族文案词面不相交（详见 §②(a)，**但只覆盖本腿打开过的那几行**＝§⑥ E13）〔量〕 |
| R27 | `sed -n '140,175p' cmd/wisp/main.go`；`sed -n '205,245p' cmd/wisp/resident_task_source_windows.go`；`grep -n "taskEntryDisabledClaim" -A 10` | `main.go:153` `reply := interactiveStdin()`，`:154-158` nil 时向 stderr 打出后果句（逐字见 §③）；`resident_task_source_windows.go:218` 同一枚 gate，`:224-230` nil＋无注入时 `slog.Warn`（`why`/`effect` 两字段）＋`fmt.Printf`；常量 `:88 taskEntryDisabledClaim = "任务入口未启用（本机没有可交互控制台）"`，注释 `:85-87` 逐字「names the consequence, not just the condition」〔量〕 |
| R28 | `grep -n "q.mu.Lock\|q.mu.Unlock\|q.ui\|deliver(" internal/agent/approval/queue.go`＋`sed -n '296,340p'`＋`sed -n '227,240p' replies.go` | `Depth :133-137` 持锁只 `return len(q.pending)`；`deliver :298-310` 持锁但发送是 `select/default` 非阻塞；`Pending()` `replies.go:227-240` 持 `r.mu` 只做拷贝；`allowScoped` 在 `:377` 先 Unlock、`:382` 才 `deliver` ⇒ §① 末段那句"读到不阻塞的形状"＝〔推〕，未跑测试〔量＋推〕 |
| R29 | **本腿没有跑的尺（点名，免得被读成"跑过、结果为空"）** | `go build`／`go vet`／`go test`／`scripts/build.ps1`／`go mod tidy` **全部零发**（派单禁令：`cmd/wisp` 里正有写腿在跑测试；`go mod tidy` 在 HEAD 上会 `exit 1` 并造出不属于本腿的 diff）⇒ 本文件**没有任何编译、测试、仪器绿灯读数**；§④ 的"会不会红"全部是**读断言取数方式推出**的〔推〕 |
| R30 | `grep -rn "h\.Allow(\|cards\.Allow(\|\.Allow(s\.ctx" --include=*.go cmd internal`（排测试）＋`grep -rn "interactiveStdin()" --include=*.go cmd/wisp \| grep -v _test` | `Replies.Allow` 的**产码调用者＝1 枚**：`cmd/wisp/approval_reply.go:215`（在 `replySurface.allow` 内）；`interactiveStdin()` 的产码消费者＝`cmd/wisp/main.go:153` 与 `cmd/wisp/resident_task_source_windows.go:218`；POSIX 半边 `cmd/wisp/approval_reply_stdin_other.go:23` 逐字 `return nil` ⇒ **GUI 常驻腿里没有控制台答复者，也就没有 `Allow` 的可达调用者**〔量，复认派单背景〕 |
| R31 | `sed -n '44,125p' internal/agent/approval/approval.go`＋`grep -rn "Statuses()\|ChannelStatus" --include=*.go internal cmd`（排测试） | `Channel`＝**四枚 L1 否决通道的闭集**（`:48-50` 逐字「The set is closed: a fifth value means a caller invented a selector … M-7/C-3」）；值 `:54-57`、`allChannels :61`、`channelNames :88-93`、`unavailableText :102-119`、`ErrChannelUnavailable :121-124`、`ChannelStatus :80-85`、`Statuses() :190-194`；消费者只有 `gate.go:605`（进卡片 `Channels` 字段，定义 `ui.go:40`）与 `run.go:1390`／`channelRosterText`（`resident_approval_windows.go:510`）〔量〕⇒ 正文 §② 第二发补充 |
| R32 | `sed -n '240,262p' cmd/wisp/resident_ball_windows.go` | `recordTrayExit` 在 `:253-258`：`const why`（逐字整句）＋`slog.Warn("tray exit requested","outcome","ignored","why",why)`＋`fmt.Printf`；注释 `:247-252` 逐字「it is **reported** instead of **assumed**」〔量〕⇒ §③ 先例五 |
| R33 | `sed -n '452,470p' internal/agent/approval/gate.go`＋`sed -n '598,610p'`＋`sed -n '30,50p' internal/agent/approval/ui.go` | `gate.go:456-463` 逐字记票 87：「The user's vote was simply not heard and the call waited out the full C18 deadline. Hand it to the queue's refusal funnel instead」；`:605` `Channels: g.channels.Statuses()`；`ui.go:40` `Channels []ChannelStatus` 在**卡片视图**结构里（`PanelItem` 的注释自 `:45` 起，逐字「There is deliberately no grant field, no Allow capability and no way to name a source」）〔量〕⇒ §③ 先例六／七 |
| R34 | `grep -rn "nativeReplySource\|NativeSource" --include=*.go cmd/wisp`（排测试） | `approval_reply.go:80-81` 两枚标签常量逐字 `nativeReplySource = "cmd-wisp-console-native"`／`panelReplySource = "cmd-wisp-console-panel-route"`；`replies.go:116-123` 注释逐字「They are LABELS, not authority」〔量，供 §② 第二发的 甲 读法〕 |
| R35 | **收尾复尺**：`git status --porcelain=v1 -- cmd internal` ＋ `git log -4 --format="%h %s"` | 交件时刻：`cmd`/`internal` 两包里脏的**只有两枚 `*_test.go`**（`cmd/wisp/panel_host_gate_test.go`／`cmd/wisp/panel_host_windows_test.go`，写腿 `33-r4` 的活）⇒ **本腿引过的那些非测试 `cmd/wisp` 文件在这三发之间没被动过**，E15 的残余风险当场降一格（但仍只算"这一瞬的读数"，落地时请现读）；HEAD 链上另有编排者两发（`f4c583db` A496、`0ecd725c` A497，逐名核过＝只动 `.scratch/**` 与 `docs/reports/**`，零产码）〔量〕 |

> ⚠ 本节里凡是"补完时逐字登记"的字样，都在同一枚文件的后续提交里被换成现读数；
> 若某格最终仍是空的，那代表**本腿没跑到那一发**，按 §交件判语 认它的射程，不许读成"跑过、结果为空"。
> R12 那一行原本写着"读数在下面补完时逐字登记"——**已由 R14 兑现**（零命中＋正控），此行保留以不改上文任何一格。

---

## ⑥ 我可能写错的条目（对抗我自己）

- **E1 票面行号与盘上行号不一致（本腿采盘上）**：票 228 AC#2 逐字写 `Replies.Allow` 在 `replies.go:313`，本腿现读函数签名在 **`:316`**（`:313` 是它上面那段注释的首行）。
  我可能把"函数体首行／注释首行／签名行"的口径搞混，也可能我读的文件与票面写腿当时读的**不是同一个 HEAD**。⇒ 下面所有 `replies.go` 行号都以 `eed229e` 为准，逐条由 R3/R4 的 grep 反取。
- **E2 「三枚候选」是我按派单列的，不是盘上现存的枚举**：派单给的是 (a) 永远显示＋按下去报错／(b) 灰掉或不显示／(c) 打开菜单时现取一次状态。
  这三枚**互相不排斥**（(c) 其实是 (a) 或 (b) 的取数时机，不是第三种"形状"），我却把它们并排列了。⇒ 编排者读数时要按"取值时机"与"无卡时的可见性"两维去理解，别把 (c) 读成独立第三形。
- **E3 符号名可能来自派单而非盘上**：`consoleVetoChannel` 在 `cmd/wisp/approval_reply_stdin_windows.go` 里**词面零命中**（R9）。
  我按"这一枚文件里那一支返回 nil 的函数"理解它＝`interactiveStdin`（`:41`，`:44`/`:49` 两处 `return nil`）。若真正被叫作 `consoleVetoChannel` 的东西在另一枚文件里，我 §③ 的指认就指错了对象。
- **E4 我没量任何运行期**：② 里所有"按下去会返回什么"都是从 `Allow` 的**返回路径**读出来的，不是真按出来的。
  `DecideFromNative` 的透传那一支（`:331`）我在这一发**没打开 gate.go 的对应函数体**，所以"门拒绝的文案 vs 权限判定的文案"是否复用，我只能给〔推〕，见 §②(a) 与 §⑦ F3。
- **E5 名册／断言原文我读了，但读的是语义不是行号**：`cmd/wisp/resident_ball_228_test.go` 里那两枚钉（名册数、回调留 nil）
  按派单禁令**不许引 `cmd/wisp/*_test.go` 的行号当权威**（写腿 `33-r3` 正在改那些文件）。我 §④ 只写名册与断言的**语义**，行号以落地时现读为准。
  ⚠ 残余风险：写腿可能在我读之后**改了那两枚断言的语义本身**（不只是行号）。那 §④ 全节过期。
- **E6 `AwaitingHuman()` 的"有卡"不等于"能允许"**：R3/R4 读到的两支——L1 卡 `Grant == ""` ⇒ `Allow` 必返 `ErrRouteHasNoAllow`；
  `AwaitingHuman` 在队列有深度但账上找不到 L2 卡时 `:262` 返回 **`ReplyCard{}, true`（corr 为空串）**，此时 `Allow("")` 走 `:319` `ErrNoTrackedCard`。
  我可能把"没有待决卡"与"有待决卡但不许允许"混成一格，而**这两格对诚实的要求不一样**（前者该说不显示，后者该说不许）。⇒ 见 §②(a) 拆三支。
- **E7 菜单"置灰"能力我没验证到 API 层**：`appendItem` 今天只有 `checked` 一维（R6）。
  若包里没有定义 `MF_GRAYED` 的 flag 常量（R12），那 (b) 那一形就要**新增常量**，动的文件比我 §②(b) 写的多一枚。
- **E8 我不许引用 `frontend/**`／`design/**`**：起手 `git status` 里那 16 枚 ` D design/**` 我照派单原样登记在 R2（那是别人的活），**没打开、没读内容**，若 §② 的文案射程与 UI 层有关，我不判。
- **E9 「三枚候选谁与先例同形」是判断不是读数**：§③ 只报"同形／相反"，不报推荐；但我给的判定本身可能被推翻（尤其 (a) 那一支：它其实**既同形又相反**，取决于"报错"是不是当场说给用户听）。我把它拆成两支写，读者若只读标题会以为我自相矛盾。
- **E10 §④ 那张表的"响／不响"是有条件的，条件写在括号里**：`(a)` 会让名册钉红，**只在该形走"新增 `internal/ball.Events` 键并在宿主字面量里设上"这条路**时成立（R15/R17 的量法：那枚钉数的是宿主 `Events{…}` 字面量里的 key）。
  若落地改成"球侧自己读一枚注入的只读闭包、不新增手势回调"，则 (a)(b)(c) **三枚都不响**，而那正好把 §④ 末段那条"分层纯度"的账换了一边（见新增的 §⑦ F7）。⇒ 读表时**别把括号当修饰**，它是那一格的成立条件。
- **E11 `ErrNoTrackedCard` 与 `ErrUnknownCorrelation` 的"近碰撞"是我给的语义定性，不是盘上现成的判决**：两枚文案词面不同、分属账／队两族（`replies.go:57` vs `ui.go:124`），有人完全可以合理地说"没复用"。
  我写的是"托盘用户看不出差别"这一**后果**，其证据只有 R19/R21 两发（两枚错误的产生点＋各自被哪枚用例钉住），**没有**任何真机读数支撑"用户看不出"。⇒ 若编排者判"两族各管各的、不算同一段文案"，§②(a) 那一段的告警要降格为口味注记。
- **E12 前人引的一句我量不到**：票面 `:116`（`228-a2` 的承重现量③）写 `internal/ball/interaction_live_test.go:12` 逐字「needs a real mouse in the notification area」，本腿 R24 在同一 HEAD 家族上 `notification area`/`real mouse` **零命中**。
  ⇒ 要么前人读的是别的 HEAD，要么那是**意译**；我只登记"本腿量不到"，**不**据此判前人写了假读数（那需要我先证明我读的 HEAD 覆盖他写的那一发）。
- **E13 "门拦下的文案不许复用权限判定的文案"这一条我只核了六行**：`internal/tools/gate.go:118/:126/:139/:144` 与 `internal/perm/store.go:179/:203`（R26），加审批族六枚（R19）。
  **没读** `internal/perm` 判定主路径全文、没读 `Decide`/`allowScoped` 之外所有出声点，也**没读** `risk`/`tools` 的其它拒绝文案。⇒ §⑦ F3 从"未答"升成"**部分答，且明写没穷举**"，不许读成"已核不复用"。
- **E14 我自造了 (a1)/(a2) 两支编号**（§③ 末）：派单只给 (a)(b)(c) 三枚。拆两支是因为 `Events` 回调签名是 `func()`、**没有错误通道**（R15），所以"报错会不会被听见"在今天的形状里**不是候选自带的属性**，而是实现选择。
  若编排者要求严格按三枚读，(a) 的判语就取 (a1)／(a2) 两支并列，不许合并成一句。
- **E15 起手锚之后盘上动过两枚**：`f4c583db`（编排者台账发，只动 `.scratch/wisp/issues/33-panel-host-c27.md` ＋ `docs/reports/pending-and-issues.md`，R23 逐名核过＝零产码）与本腿 `07f74694`。
  ⇒ 我 §①–④ 的产码引用在这一发上**没漂**；但**残余风险**：若写腿 `33-r4` 在本文件后续发之间往 `cmd/wisp` 落了新码，我的 `cmd/wisp` 非测试文件引用（`resident_ball_windows.go`／`resident_approval_windows.go`／`approval_reply*.go`）就可能漂——那批文件**不在**写腿的"只改测试"范围内，我防不住。⇒ 所有 `cmd/wisp` 引用请以落地时现读为准。
- **E16 派单给的写腿代号与台账不一致**：派单写 `33-r3`，台账 A496（R23 读到的 commit 标题）逐字说那枚代号已被改名为 `33-r4`，因为 `33-r3` 早先被用过且那一发死掉。我按台账写 §④/§⑤，但**这与本腿纪律无关**（我只读、没跑任何尺在 `cmd/wisp` 的测试上）。
- **E17 「第五枚 id 今天不存在」这句我说的依据只是 `:20-24` 那五行**（R6）；我没有核 `internal/ball` 之外是否有别处定义过同值菜单 id。风险低但没穷举。

---

## ⑦ 判不动的地方（逐条甲／乙／不做＋现量）

> 派单纪律：**最后一枚用哪一形由编排者裁，本腿不许自选**。下面每条只列代价与现量，不写推荐。

- **F1 "有卡吗"这一维怎么进菜单构造**（① 那一跳）。现量：`showMenu` 签名 `tray_windows.go:72` 只有两枚布尔；`Ball` 上的字段由 `SetTrayChecks(muted, pausedWake bool)`（`ball_windows.go:919`）写入；`wmAppTray` 分流在 wndproc 内（`ball_windows.go:669-686`），而 wndproc 跑在 `ui-sta` 线程里。
  - **甲（推：像 `SetTrayChecks` 那样加第三枚状态位）**：改 `SetTrayChecks` 签名或新增一枚 setter ⇒ 动 `internal/ball` 的公开面＋`cmd/wisp` 的调用者；⚠ 若新增 `Events`／公开方法，会不会撞 `cmd/wisp` 那枚"名册／回调必须有执行者"的钉之外还撞别的名册钉，本腿不判（见 §④）。时效性代价：状态由别人**推**过来，推的那一次没跑＝菜单上是旧值。
  - **乙（拉：`showMenu` 里同步调一枚注入的只读闭包）**：`internal/ball` 需要多一个注入点（放在 `Options`／`Events`／还是参数上＝三选一的契约面，本腿不裁），`cmd/wisp` 那侧提供闭包读 `Replies.AwaitingHuman()`（`replies.go:252`）。代价：审批状态住在 `cmd/wisp`，读它要过 `Replies.mu`（`replies.go:105`）与 `Gate.Queue().Depth()`（`:311`），**在 ui-sta 线程里同步读会不会与那条线程互相等，本腿没量**（见 F2；**第二发补了半边读数，见 F2 末段**）。
  - **丙（不显示那一枚项，只在菜单被打开时决定要不要摆）**：＝乙 的时机变体，代价同乙，只是 `appendItem` 前多一次判断。
  - **不做**：托盘第五枚项不落地 ⇒ AC#2 那一格按票面 `:39` 判不过（常驻腿仍无"允许"入口）。现量：`recordTrayExit` 只打日志（R7），产码里 `Allow` 的托盘调用者＝0（R10）。
- **F2 ui-sta 里同步读审批状态的安全性**。现量：票面 `:92` 逐字承认"在回调里同步收口＝自等待死锁"这一族风险（`Ball.Close()` 是 `sta.PostTask(…) + <-done`），但它讲的是收口不是读；本腿**没读 `Gate.Queue().Depth()` 的锁序**，因此"读一把队列深度会不会撞上正在 prompt 的那条线程持有的锁"判不了。⇒ 要么派一枚专门读锁序的腿，要么由写腿在自己那发里量。**本腿不许用 `go test` 去量它**（派单禁令）。
  - **本腿第二发现有补充读数（不改上面那句原判，只把射程往前挪一格）**：`Depth()`（`queue.go:133-137`）持锁只做 `len(q.pending)`；`deliver`（`:298-310`）虽持锁，但里面的答案是 `select { case it.answer <- a: default: }`＝**非阻塞**；`Pending()`（`replies.go:227-240`）与 `routes()`（`:384-389`）都是短临界区。⇒ 〔推〕**在我读过的那几枚函数里不存在"持锁等 ui-sta"的边**。没覆盖的：`queue.go` 里 `q.mu.Lock` 还出现在 `:409/:423/:446/:515/:537/:553`（**本腿没逐枚读体**），且 `gate.go:288/:530` 那两处 `g.ui.Prompt(...)` 是否在持任何锁下调用**没核**。⇒ **F2 仍判不了，但从"零读数"升到"半边有读数"**；要闭环请只读这三发，不要跑测试。
- **F3 "门拦下的文案"与"权限判定的文案"是否已被复用**。现量：`Allow` 的三支文案逐字见 R4，全部带 `approval:` 前缀；本腿**未打开 `gate.go` 的 `DecideFromNative` 函数体**，也未打开 `internal/perm` 的拒绝文案。
  ⇒ 「(a) 那一形会不会撞『门拦下的文案不许复用权限判定的文案』那条硬规矩」这一格**没答**。判据：`grep -rn` 两套文案词面是否相交（本腿没跑那一发就跑完了预算的话，请由下一枚腿补，别把这里读成"已核不复用"）。
  - **本腿第二发把这一格从"没答"改成"部分答"（原判不抹）**：后续发了 `DecideFromNative`（`gate.go:718-724`）、`DecideFromPanel`（`:730-743`）、审批族六枚（R19）、门拦下四句（`internal/tools/gate.go:118/:126/:139/:144`）、权限两句（`internal/perm/store.go:179/:203`）⇒ **在我读到的这十句里三族词面不相交**〔量 R26〕。
    剩下的没答部分是：`internal/perm` 判定主路径与 `tools`/`risk` 其余出声点**没穷举**（E13），且我发现的那处近碰撞在 **approval 族内部**（`replies.go:57` vs `ui.go:124`，E11），不在"门 vs 权限"那条硬规矩的射程里。⇒ **硬规矩这一格：本腿未见违反，也未证其不存在。**
- **F7（第二发新增）名册钉与"球侧对审批一无所知"这两条今天互斥**。现量：
  ① 球宿主对审批的无知是被写下来的边界（`cmd/wisp/resident_ball_windows.go:56-61` 的 `escVetoFunc` 注释、`:23`/`:195-203` 的 `ballGestureWhy`）；
  ② `internal/ball` 的回调签名是 `func()`，`fire` 只判 nil（`ball_windows.go:728-732`）⇒ **错误/结果没有回程**，"按下去若无卡就报错"这件事在今天的形状里只能由 `cmd/wisp` 那侧出声（审计或球状态），不能由 `showMenu` 的返回值带出来；
  ③ 那枚 AST 名册钉数的是宿主 `Events{…}` 字面量里的 key（§④ P1 取数方式）。
  - **甲（新增第五枚 `Events` 键）**：执行者进得来、名册钉会响 ⇒ 落地必须同时扩那枚**别人已勾判据**的名册（票面 `:102` 逐字警告「⛔ 不许为了"让钉响"去扩那枚名册断言到行为维」，那是 AC#1 的形状；扩**列表**与扩**断言维度**是两件事，本腿分不清算不算越界 ⇒ 交裁）。
  - **乙（球侧直接调注入的只读闭包／或直接调 `Allow`）**：名册两枚钉**都不响**，代价是破 ① 那条分层，且把"审批"这一概念引入 `internal/ball`。
  - **不做**：托盘第五枚项不存在 ⇒ 回到 F1 的"不做"那一支。
  ⇒ **本腿不选**：甲乙的取舍是"要钉响还是要分层纯"，两维都有人在守（AC#2 的判据形状 vs 票面 `:102`），只有编排者能同时改。
- **F4 无卡时"允许一次"该给哪种诚实**（不显示／置灰／点了出声）。三支都合法、代价不同，且**要不要出声归票面 AC#2 的"每枚按钮断一条带正确 Channel 的答复"这一句怎么读**——那是编排者的裁，不是盘上现成的答案。现量：`appendItem` 今天不会置灰（R6），`showMenu` 的返回语义里 0＝被关掉（`:70-71` 注释逐字"Returns the selected command id (0 = dismissed)"）。
- **F5 L1 那一发的枚举面**。现量：票面 `:19` 逐字写「**L1 的 `Gate.windows` 全仓零枚枚举口**」，并指名"要枚举 L1 那一发就撞票 220 的 AC#2 甲形，先裁甲乙"。
  ⇒ 若第五枚项要对 L1 出声，本票就要先造枚举口＝撞票 220 的"两腿各造一次＝两份枚举器"（票面 `:60`）。**这一格不在本腿射程内，只登记它挡住了哪一支。**
- **F6 「允许一次」拿哪枚 corr**。现量：`Allow(ctx, corr)` 需要关联号（`replies.go:316`），`AwaitingHuman()` 只给一枚卡（`:252`），`Pending()` 按显示序给全部（`:227`）⇒ 若同时有两枚卡在人等，托盘那一枚"允许一次"允许的是**哪一枚**没有定义。甲＝用 `AwaitingHuman` 的那枚（L2 优先，见 `:242-245` 注释原文）；乙＝菜单里列出全部（＝第四枚以上的新面，动 `showMenu` 结构）；不做＝本票只支持单枚并发。**本腿不选**——票面 `:39` 那句"每枚按钮断一条带正确 Channel 的答复"没说并发。
- **F8（第二发新增）AC#2 那句"带正确 `Channel` 的答复"对三枚 allow 按钮今天没有值可带**。现量：`Channel` 是**四枚否决通道的闭集**，`internal/agent/approval/approval.go:48-50` 逐字禁止第五枚值（「a fifth value means a caller invented a selector, which is exactly the M-7/C-3 failure shape」）；`Allow` 走的 `Request` 只有 `CorrelationID/Allow/Grant/Source`（`replies.go:328-330`），而 `Source` 被注释逐字定性为「LABELS, not authority」（`replies.go:116-123`）。〔量 R31/R34〕
  - **甲**：把"正确 Channel"读成 `NativeSource` 标签（现成面，R34）⇒ AC#2 可按字面写用例，**零新契约**。
  - **乙**：读成"要真给托盘一枚 `Channel` 值"⇒ **动闭集＝动 SPEC-06 §2 的冻结词表**，且正落在 `:48-50` 那句被点名为失败形状的描述上 ⇒ 本腿只登记"这一支要人工批准"，**不当它是可选项**。
  - **不做**：把票面 AC#2 的判据句子改写成"带正确 `Source` 标签的答复"——那是**改票面文字**，归编排者，实现腿与本腿都不许自己动。
  ⇒ 这一格与"有没有卡"**正交**：三枚候选谁都不影响它，但它决定落地腿能不能按票面原话写判据。

---

## 交件判语

**射程（答到了哪一格）**
- ①：菜单侧今天**只收两枚布尔**（`tray_windows.go:72`）、那两枚布尔的推线**在生产里零调用者**（R16）、审批状态住在 `cmd/wisp` 且**已经在该侧被读三次**（R17）、依赖箭头今天只朝一个方向（R18）；三条入口路线（甲推／乙拉／丙＝乙的时机变体）各自要动谁、代价是什么，全部列了而**没选**。
- ②：(a) 的返回语义逐字量全（四支出声＋透传支的三枚队列错误，R4/R19）；**"有卡在等"与"允许得动"是两格**这一点有原文支撑（`replies.go:262` 空 corr、`:322` L1 无 allow）；"门／权限"两套文案在**本腿读到的十句里不相交**（R26），近碰撞在 approval 族内部（E11）；(b) 的置灰需要**新增一枚不存在的 flag 常量**（R14）；(c) 的取数时机今天落在**同一次调用内、弹出之前**（R6），且 `WM_COMMAND` 那一支援射法已被票面 `:100` 逐字关掉、本腿复认其凭据未漂。
- ③：抽出共同形状一句，并给出**七枚**先例的逐字出处（热键 `Standby`／`interactiveStdin` 返 nil／`cmdRun` 后果句／常驻腿的 claim＋`why`＋`effect`／`vetoByEsc` 无卡那句／`bindBallHost`「advertise 一枚按不动的键＝B1 禁止的形状」／`ChannelStatus`＋`unavailableText`＋`ErrChannelUnavailable`＋票 87 的"落空一票改路由"）；同形／相反逐枚报，**(a) 拆成 (a1)/(a2) 两支**（因为回调没有错误回程，R15），**没有推荐**。
- ④：两枚钉逐字读完（名册钉的两条断言＋取数方式），并给出**三枚候选×两枚钉**的响／不响表；`internal/ball` 的 **11 枚热键状态用例逐枚点名＝不管这一格**（R22），`internal/ball` 唯一的 AST 走形用例是 C21 token 表＝不管；**托盘侧仪器现状＝零**（R20 带正控）。

**没答（不许读成答了）**
- **F2 的剩下半边**：ui-sta 上同步读的完整锁序——我只读了 `Depth`/`deliver`/`Pending`/`routes`，`queue.go` 另有 6 处 `q.mu.Lock` 与 `gate.go:288/:530` 两处 `ui.Prompt` 的持锁状态**没读**；没跑任何测试（R29）。
- **F3 的剩下半边**：`internal/perm` 判定主路径与 `tools`/`risk` 其余出声点**没穷举**（E13）⇒ 硬规矩那一格只能写成"本腿未见违反，也未证其不存在"。
- **F5**：L1 那一发的枚举口归票 220 裁（票面 `:19`/`:60`），本腿只登记它挡住了哪一支。
- **F6**：多枚卡并发时"允许一次"允许哪一枚**没有定义**，票面 `:39` 没说并发。
- **F7**：名册钉 vs 球侧分层纯度**互斥**，两维都有人在守，本腿不选。
- **F8**：AC#2 那句"带正确 `Channel` 的答复"对 allow 按钮今天**没有值可带**（闭集，R31）；这一格与本题正交，但**最需要编排者先裁**，否则落地腿会自己去动那个闭集。
- **未答的形状问题本身**：派单问"该长成什么样才叫诚实"。本腿能给的是判据读数——**三枚候选里没有任何一枚能单靠 `showMenu` 自身满足先例五/六/七的"出声"要求**（`fire` 无回程 R15、返回值只有命令 id 且 0 兼作"被关掉" R6），因此"诚实"那一维在今天的形状里**必然落在 `cmd/wisp` 侧（审计行／球状态／tooltip）而不是菜单里**。这是一条**约束读数，不是选形**；选哪一形仍归编排者。

**纪律自证**：本腿零产码、零脚本、零测试、零构建（`go build`／`go vet`／`go test`／`scripts/build.ps1`／`go mod tidy` 全部零发，R29）；只新增并提交了 `.scratch/wisp/probes/228/a3/census.md`；票 228 只**追加** Progress log 行（`git show --numstat` 对票面文件＝**10 插 0 删**）；AC 复选框一枚未碰；`frontend/**`／`design/**` 未读未引；未新增任何裸 `go`；脏项（`design/**` 的 16 枚 ` D`、`.gitignore`、`.scratch/**` 未跟踪件）一律未动；**只 commit、未 push**。
