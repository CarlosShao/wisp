# 票 243 ／ 腿 `243-c1` — 「ticket 07」25 处注释逐处判定表（只读普查，零产码）

- 腿：`243-c1`（只读普查）。工作树锚点：**`90cf65b2`**＝**本腿唯一一枚现读的 HEAD**（`git log --oneline -1`，落盘后跑的）。这枚 HEAD 本身已经是写腿 `228-r1` 的骨架件，说明**本腿全程与它在同一棵动着的树上**；因此 §0 那把尺**在这枚锚点上重跑过一次**：R1 仍＝**28**、`cmd/wisp` 对 `internal/ball` 的 import 行数仍＝**0**、全仓 import 该包的仍只有 `cmd/balldebug/main.go` ⇒ 本表读数是**同一棵树上的**，非跨头拼凑。（⚠ 我没有在更早的提交上现读过任何数，所以不写「从 X 推进到 Y」——那是别人的读数。）
- ⚠ 自我纠正一条：初稿这里写了一枚**没有现读过的锚点号**（正是本票要治的那类「按印象报状态」），现换成实测值。
- 本表纪律：只 `grep`／`sed`／`Read`／只读 `git log`。**零 `.go` 改动、零注释改动、零 `go test`、零 `go build`、零新仪器**。
- ⛔ 本文件每一节起手即写满，**不含任何占位符词**（那三个常被用来标记未完成面的中文／英文词，本文件里一个也没写，故下游 `grep` 应零命中）。

---

## 0. 我自己跑的尺与读数（不抄票面）

| # | 尺（逐字） | 我的读数 | 与票面是否一致 |
|---|---|---|---|
| R1 | `grep -rn "ticket 07" --include=*.go internal cmd scripts tools \| wc -l` | **28** | 一致 |
| R2 | R1 结果里落在**字符串字面量**（非注释）的行数 | **3** | 一致 |
| R3 | 28 − 3 | **25** | 一致 |
| R4 | 按目录分布（R1 全量 28） | `internal/ball` 9／`internal/statemachine` 5／`cmd/wisp` **6**／`internal/proc` 3／`internal/audio` 2／`cmd/balldebug` 2／`internal/observe` 1 | ⚠ 见下方「票面勘误」 |
| R5 | 注释堆 25 处的目录分布 | `internal/ball` 9／`internal/statemachine` 5／`cmd/wisp` **4**／`internal/proc` **2**／`internal/audio` 2／`cmd/balldebug` 2／`internal/observe` 1 ＝ **25** | 一致 |
| R6 | `grep -rn "ticket 07" --include=*_test.go internal cmd \| wc -l` | **6**（`liquid_test`／`live_windows_test`／`tokens_test`×2／`sampler_test`／`machine_test`） | 一致 |
| R7 | `grep -rniE "arrives in ticket\|lands with ticket\|is ticket [0-9]+" --include=*.go internal cmd \| wc -l` | **132** | 一致 |

**票面勘误（具名，不改票面）**：票 243「现量」表把注释堆写成「`internal/proc` 2／`cmd/wisp` 4」以外的组合时是**按去掉 3 字符串之后**写的（那一行确实写了 `cmd/wisp` 4，与我的 R5 一致），但它同时把 R1 的全量 28 与 `cmd/wisp` 4 并列排在全仓命中行里，容易读成「`cmd/wisp` 全量 4」。我实测 **`cmd/wisp` 全量＝6**（4 注释 ＋ 2 字符串：`main.go:25` usage 正文、`resident_windows.go:81` Printf）。⇒ 编排者那句「25／3」的读数**我复现成功了**，只有目录行的歧义值得在票面补半句「`cmd/wisp` 4 是去字符串后的数」。

**判据（我按什么分三类，写在前面免得事后追认）**
- **历史出处**：句子把票 07 当作**一件今天已经在代码里的东西的来源**（「Implemented by ticket 07」「ticket 07 acceptance」「ticket 07 constraint ＋ 一条真实存在的仪器」）。判定动作＝去代码里找那样东西在不在，在＝过。
- **还在承诺未来**：句子用**交付／归属动词**（`is`／`arrives in`／`lands with`／`subscribes`／`is …'s job`）把一件**今天不在这条路径里**的活指给票 07，而票 07 已带 `-done`。判定动作＝① 那件活今天在不在被指认的位置（import／调用者／构建旗标），② 票 07 盘上是否 `-done`。
- **说错了事实**：句子对**今天可测量的树**断言了一件可证伪的事，而那件事**不成立**（不论有没有票号）。
- 三类**互斥**，每行给一枚主判；兼类时在主判后括注第二支，避免下游把一行当两行用。

---

## 1. ⛔ 不写进 25 行表的那 3 处＝运行时字符串／用户可见文案（AC#2：全部 🔴，当场跨引票 228 AC#7）

这 3 处**已由另一枚写腿 `228-r1` 在改**（本腿⛔未碰任何一字节）。本腿只登记。归属凭据：`.scratch/wisp/issues/228-ball-and-tray-are-not-in-the-process-that-runs-tasks.md` 的 **AC#7 射程更正段（账 `A471`）**，其第 ① 条逐字把这三处列为「本格的完成判据」；第 ② 条逐字把「其余 25 处注释」推给票 243。**两票不重叠、本表与 AC#7 不抢活。**

| 🔴 位点 | 载体（现读） | 现读原文逐字 | 今天为何是假话／误导 | 归属 |
|---|---|---|---|---|
| `cmd/wisp/main.go:25` | ⚠ `const usage = ` 反引号原始串＝**`wisp -h` 打给用户的正文** | `                   event loop; the floating ball window is ticket 07)` | 用户今天真读得到；票 07 已 `-done`，而 `cmd/wisp` **零 import `internal/ball`**（实测：全仓 import 该包的只有 `cmd/balldebug/main.go`）⇒ 「球窗口＝等票 07」把一个不存在的未来念给用户听 | 票 228 AC#7 ① |
| `cmd/wisp/resident_windows.go:81` | `fmt.Printf` 运行时输出 | ``	fmt.Printf("wisp: empty event loop running; the floating ball arrives in ticket 07 (Ctrl+C exits cleanly)\n")`` | 就是骗过三程的那一枚（`240-c1` §3.1 → `241-r1` → `241-v1`）。`arrives in` ＝ 交付动词 ＋ 已结案票号 | 票 228 AC#7 ① |
| `internal/proc/boot_windows.go:127` | `slog.Info` 运行时日志 | ``				slog.Info("activation requested by second launch (ball bring-to-front lands with ticket 07)")`` | 跨了 `internal/proc` 地界；`lands with` ＝ 交付动词 ＋ 已结案票号。该分支今天**只打日志再 `ResetActivation()`**，没有任何 bring-to-front | 票 228 AC#7 ①（⚠ 动 `internal/proc` 一行，交件须具名；本腿只登记） |

**这 3 处与表内第 13／14／24／25 行是同一件事的两面**（注释侧说「球的活＝票 07」，字符串侧对用户说同样的话）。⇒ 写腿改字符串时必须**连带**改那几行注释，否则同一地界留下两份互相矛盾的真相。

---

## 2. AC#1 主表：25 行（＝R3 的 25，行数恰好；正文逐字取自我自己跑的 R1）

列：**｜#｜位点｜现读原文逐字｜判定｜若为承诺：这笔活今天归哪张票（按盘上 `-done` 核）**

| # | 位点 | 现读原文逐字（grep 输出原样，含行首缩进） | 判定 | 今天归属／判定依据 |
|---|---|---|---|---|
| 1 | `internal/audio/gate.go:21` | `//     events hook the mute hotkey path into the Muted state (ticket 07).` | **历史出处** | 两侧都活在今天：gate 侧 `GateEvent`＋`SetMuted`（同文件 :165 注明 driven by the mute hotkey）存在，ball 侧 `EvMuteHotkey`／`OnMuteHotkey`（`internal/ball/ball_windows.go:38,52`）存在，`Muted` 是 D43 状态。句子在解释事件契约的**来处**，没许诺谁去做 |
| 2 | `internal/audio/gate.go:62` | `// the ball/state machine subscribes at wiring time, ticket 07/16).` | **还在承诺未来** | **票 228**（`228-...-not-in-the-process-that-runs-tasks.md`，盘上**无** `-done`）。实测 `WithGateEvents` 全仓**只有 `gate_test.go` 两枚调用者、零生产调用者** ⇒ 「ball/状态机在 wiring 时订阅」今天没发生；07 已结案不欠这笔，16（`16-s2-acceptance.md`，**无** `-done`）是**验收票**、不产实现，指它做归属是空的 |
| 3 | `internal/ball/anim.go:5` | `// Animation discipline (SPEC-08 §2, ticket 07):` | **历史出处** | 下面逐条列的纪律就由本文件实现（`AnimationPolicy` 被 `liquid.go:83`、`tokens_test.go:314` 真实调用）⇒ 引来源，非许诺 |
| 4 | `internal/ball/doc.go:24` | `// Implemented by ticket 07 (S1: static states + basic animations; full` | **历史出处** | 过去时。⚠ 但它的**下一行**（:25「visual polish gate is the human acceptance at ticket 12」）指的票 12 盘上**无** `-done`＝活票，合法。包文档（godoc 面）可见，非运行时串 |
| 5 | `internal/ball/liquid.go:76` | `// exceptions in anim.go - so ticket 62 never arms a timer in a state ticket 07` | **历史出处** | 对 07 是「已冻结的策略」＝来处；句中另一枚票 62（`62-liquid-glass-ball-visuals.md`）盘上**无** `-done`＝**真在飞的活票**，被约束方指认正确。`transitionDriven` 由 `AnimationPolicy` 现算，不是空头承诺 |
| 6 | `internal/ball/liquid_test.go:227` | `// transition timer may only ever exist where ticket 07's frozen policy already` | **历史出处** | 「frozen policy」＋同测试真在跑（`TestTransitionTimerSetIsTheFrozenTimerSet` 逐状态对算） |
| 7 | `internal/ball/live_windows_test.go:11` | `// windows. The ticket 07 acceptance items exercised here: window stack` | **历史出处** | 列的是**本文件正在执行**的验收项（tag `windows && winlive`），＝历史验收记录，无未来指认 |
| 8 | `internal/ball/tokens.go:10` | `// literal in any other ball file is a review-rejecting violation (ticket 07` | **历史出处** | 括注后半句点了**真实存在的仪器** `TestNoHardcodedColorsInBallPackage`（`tokens_test.go:86`）⇒ 约束今天有牙，不是等别人来加 |
| 9 | `internal/ball/tokens.go:442` | ``	// Frame cap: <=30fps (ticket 07 animation discipline).`` | **历史出处** | 常量 `MaxAnimFPS = 30` 就在下一行，纪律已落地 |
| 10 | `internal/ball/tokens_test.go:84` | `// is a review-rejecting violation (SPEC-08 §2, ticket 07 constraint).` | **历史出处** | 该注释就是紧挨着那枚仪器函数的 doc，引来源 |
| 11 | `internal/ball/tokens_test.go:309` | `// TestAnimationPolicyZeroTimerInSleeping is the ticket 07 acceptance: the` | **历史出处** | 「是票 07 的验收」＝断言归属记录，函数体就在下面 |
| 12 | `internal/observe/sampler_test.go:62` | ``	// vocabulary must not drift from the 20-state table (ticket 07).`` | **历史出处** | 「20-state」现读为真：`states.go:3` 逐字写「one of the 20 BallState values (D43; C12 frozen)」，测试真的把 `SLOState` 逐个映射到 `statemachine.State*` 常量。**注意权威是 D43/C12 冻结契约，票 07 只是次要来处**，别把这句改成「票 03 的表」（词面反噬，见 §4） |
| 13 | `internal/proc/singleinstance_windows.go:17` | `// (bringing the ball to the front is the ball's job, ticket 07) and exits.` | **还在承诺未来** | **票 228**。三条件全中：①交付归属句（`is … job`）②那件活今天不在这条路径（`cmd/wisp` 零 import `internal/ball`）③07 已 `-done`。现读第二实例分支只 `slog.Info` 一下（表外 `boot_windows.go:127`），没有任何 bring-to-front |
| 14 | `internal/proc/singleinstance_windows.go:84` | `// loop (ball, ticket 07). Handle with care: do not close it.` | **还在承诺未来** | **票 228**。函数 `ActivateEvent()` 与其消费者 `RunEventLoop` 今天都在（`boot_windows.go:126` 真在 `WaitForSingleObject`），但括注把**消费端**认成「ball, ticket 07」；ball 今天不是这条进程的一部分，07 已结案。⛔ 半句「do not close it」是**真约束**，改注释放大这半句、别丢 |
| 15 | `internal/statemachine/doc.go:19` | `// Implemented by ticket 07; ticket 03 pinned only the State vocabulary.` | **历史出处** | 过去时且可复验：表（`doc.go:7` 指 `table.go`，40 行 ＋ #41/#42）与 `timeouts.go` 今天都在；票 03（`03-skeleton-runtime-rules-done.md`）**有** `-done`，两句都是史 |
| 16 | `internal/statemachine/machine.go:46` | `// timeout row (notably Sleeping: the zero-timer discipline, ticket 07).` | **历史出处** | 纪律有牙：`machine.go:180` 的 `rearmLocked` 按 `timeoutKey{state,phase}` 才装表，`TestSleepingHasZeroTimers`（`machine_test.go:62`）钉住 |
| 17 | `internal/statemachine/machine.go:65` | ``		// Explicit no-op side-effect sink (ticket 07: hooks fire as events,`` | **历史出处** | 「hooks fire as events」＝07 已交付的机制（`Effect`/`Sink` 存在且默认 no-op 就是它的设计）。⚠ 第二行「consumed by later tickets」是**无指代未来**——它没点已结案票号，故不在本票三类里；但它是 §4 说的第二种形状（见 §4 附言），登记不判 |
| 18 | `internal/statemachine/machine_test.go:61` | `// TestSleepingHasZeroTimers is the ticket 07 acceptance assertion: after any` | **历史出处** | 同 #11，断言就在下面函数体里 |
| 19 | `internal/statemachine/states.go:5` | `// timeouts are implemented by ticket 07.` | **历史出处** | 全句（:3-5）是「vocabulary only in ticket 03 : the transition table, guards and per-state timeouts are implemented by ticket 07」——过去时，三样（表／guard／`timeouts.go`＋`DefaultTimeouts()`@`machine.go:70`）今天都在 |
| 20 | `cmd/balldebug/main.go:3` | `// Command balldebug is the ticket 07 debug harness: it boots the real ball` | **历史出处** | 现在时，且**为真**：`balldebug` 是全仓唯一 import `internal/ball` 的 main（实测），`scripts/dev/ball-cycle.ps1:30` 真在 build 它 |
| 21 | `cmd/balldebug/main.go:590` | `// ticket 07 wiring: click = summon/veto, mute toggles, cancel vetoes).` | **历史出处** | 该 wiring 今天**真在**：`gesture()`（:591 起）逐 case `dispatch(b, m, statemachine.EvVeto/EvInterrupt/EvSummon)`。句子描述的是自己脚底下那段码 |
| 22 | `cmd/wisp/console_other.go:7` | `// The Windows build links as a GUI-subsystem binary (ticket 07), which starts` | **说错了事实**（兼：承诺） | 现读证伪：全仓（排除 `frontend`／`design` 零读区）**唯一**出现 `windowsgui` 的地方就是注释本身；`scripts/build.ps1:103-110` 的 `$ldflags` 只有 6 枚 `-X` 版本注入，`-H` ＝ 0 枚；`git log -S"windowsgui" -- cmd/wisp/` 只回 `bcc892c9`（＝build-chain S0 那次提交，票 01 已 `-done`）⇒ 这句话既**说错今天的构建**，又把补齐指给已结案的 07。**危险见 §3** |
| 23 | `cmd/wisp/console_windows.go:26` | `// the windowsgui subsystem (the final GUI build, ticket 07): such processes` | **说错了事实**（兼：承诺） | 同 #22 同一把尺证伪；且 `the final GUI build` ＝ 一件尚未发生的产品事实被写成现在时前提。**后果**：`attachParentConsole` 的 no-op 分支（:32「console-subsystem builds」）今天恰是**真分支**，而注释说它不是 |
| 24 | `cmd/wisp/main.go:8` | `// frozen D38(e) 10-step shutdown order. The floating ball GUI is ticket 07.` | **还在承诺未来**（兼：错事实） | **票 228**。`is` ＋ 不在这条路径的活 ＋ 已结案票；实测 `cmd/wisp` 零 import `internal/ball`。这句与 §1 的 `main.go:25` usage 串是**同一文件同一主张的两个受众**（godoc 面／用户面）——只改串不改注释＝留两份真相 |
| 25 | `cmd/wisp/notify_windows.go:12` | `// behind it (the resident process's own tray icon is ticket 07/62's, in` | **还在承诺未来** | **票 228**（tray 装进跑任务的进程＝228 的标题本体）。逐字辨：「nothing here touches it」那半句**是真话**（CLI 侧确实自起自灭 message-only 窗），假的是「resident 进程自己的 tray icon ＝ 票 07/62 的」——`internal/ball/tray_windows.go` 存在于 ball 里，但那条进程今天不载 ball；07 已结案，62 未结案 |

**行数自证**：`25` 行 ＝ R3 的 25；表体逐行对应 R1 去掉 §1 三行后的每一条，无增无减。
**判定统计**：**历史出处 18** ／ **还在承诺未来 5**（#2 #13 #14 #24 #25）／ **说错了事实 2**（#22 #23）。
（「历史出处」一支非空 ✓，满足票面 AC#1 的完成判据；主判互斥，#22 #23 #24 的兼类只作括注、不重复计数。）

---

## 3. 最危险的三行（为什么危险＝它骗人的成本，不是它错得多离谱）

1. **`cmd/wisp/console_other.go:7` ＋ `cmd/wisp/console_windows.go:26`（#22／#23，成对算一枚）**
   现读：「The Windows build links as a GUI-subsystem binary (ticket 07), which starts with no console at all」。
   为什么最危险：**它今天就是假的，而且它一错就连环错三样东西**——① `attachParentConsole` 的**整段存在理由**（注释把防守的前提说成既成事实）；② 该函数 :32 那个 no-op 早退分支被注释判定为「非最终形态」，可读成**死代码**；③ 下一枚普查腿会照抄成「GUI 构建已落地（票 07）」，于是**没有任何构建脚本需要补 `-H windowsgui` 这笔活**——一笔真正的缺口被一句注释**注销**了。这正是票 243 根因链条那一类（三程照抄、零枚仪器看得见），但比 `resident_windows.go:81` 更糟：那枚骗的是「还要做的事」，这两枚骗的是「已经做过的事」。**盘上核**：票 07 `-done`，全仓无 `-H`/`windowsgui` 旗标（尺见 §2 第 22 行）。
2. **`cmd/wisp/main.go:8`（#24）**
   现读：「The floating ball GUI is ticket 07.」
   为什么危险：它坐在**包文档的第一屏**（`Command wisp` doc），是任何人／任何腿理解「这个二进制里有什么」时读到的**第一句**。它把球 GUI 说成票 07 的现在时归属，而 `cmd/wisp` 对 `internal/ball` **零 import**——于是票 228 的整笔活在 `cmd/wisp` 的门面文档上**看起来不存在**。同文件 :25 那枚 usage 串（§1 🔴）只是把同一句谎话念给用户听；**改串不改这行注释＝谎话少了一半受众、仍然在场**。
3. **`internal/proc/singleinstance_windows.go:17`（#13）**
   现读：「(bringing the ball to the front is the ball's job, ticket 07) and exits.」
   为什么危险：它是 `boot_windows.go:127` 那枚 `slog.Info`（票 228 AC#7 ① 已认领）**在注释侧的双胞胎**。票 228 会把字符串改成带条件的事实句；如果没人同时改这一行，**同一地界第二天就留下一处已更正、一处仍承诺**。并且它用「is the ball's job」把责任推给一个**角色**（ball）＋一枚**已结案票号**——今天 ball 根本不被那条进程装载，所以「谁都会以为有人欠这笔活、而没人欠」。

（次危险但不入选的理由：#2／#14／#25 同样是承诺，但它们骗的是「接线尚未做」，读者至少不会以为已完成；#22／#23 骗的是「已完成」，会直接掐掉一笔活。）

---

## 4. AC#3 边界结论（不许为空）：什么形状的句子必须改成带条件的事实句

**先给读数，再给规矩。**

- R7 的 132 处按**动词框架**拆开：`is ticket N` **125** ／ `lands with ticket` **6** ／ `arrives in ticket` **1**。
- 132 处按**盘上 `-done`**拆开：**只指已结案票号 50 处**／只指未结案票号 80 处／两者混指 2 处。

**结论一：分界不能落在「有没有票号」，也不能落在「票号是否已结案」。**
50 处已指结案票号，其中绝大多数是**占有／来处框架**——「which is ticket 124's batch conversion」、「that is ticket 102's account」——它们在说**某件东西的出处或归属是哪张票的档案**，读它的人不会被驱动去做任何事。若把「不许引用已结案票号」写成规矩，这 50 处全部要改，其中 `internal/ball`、`internal/statemachine` 里那批「ticket 07 acceptance / frozen policy」会**被改坏**：它们的价值正是「这条不变量由哪一次验收钉住」，抹掉之后下一位没有回查落点。**这就是票 228 AC#7 第 ③ 条禁词面型仪器的理由，本表复述并支持它。**

**结论二：必须改成带条件事实句的，是「交付／归属动词框架」＋「那件活今天不在被指认的位置」这两条同时成立的句子。** 我的尺（可操作、无需新仪器，人工逐处判）：

> **A 柱（动词）**：句子把一件活**交给**一枚票号——`arrives in ticket N` / `lands with ticket N` / `is implemented by ticket N`（现在时）/ `is ticket N`（＝归属指派）/ `is N's job` / `will be built by`。
> **B 柱（位置）**：那件活**今天不在被指认的位置**——判法不是读注释，而是去盘上核：该包的 import 集、该函数的生产调用者数、该构建旗标是否存在、该票面是否带 `-done`。
> **A ∧ B ⇒ 必须重写为带条件的事实句**：说清「今天这里没有 X；补上 X 是票 M 的活（M 盘上未结案）」，条件写进句子，而不是把票号换个数字继续挂着。
> **A ∧ ¬B ⇒ 改成现在时为真的描述句**（例 #20／#21：`balldebug is the ticket 07 debug harness` 保留，因为它现在真的是；出处可另附一句「由票 07 交付」）。
> **¬A ⇒ 不动**，即使它引用了已结案票号（那 50−7 里的绝大多数）。

**为什么本表 25 行里只有 7 枚落进 A∧B／A∧¬B 的 A 柱，其余 18 枚不动**：因为 `ticket 07` 的 28 处命中里，**只有 5 处出现在 R7 的动词框架里**（`boot_windows.go:127`、`resident_windows.go:81`、`main.go:8`/`notify_windows.go:12` 的 `is ticket`、`singleinstance_windows.go:17` 的 `is … job`），其余 20 余处是「SPEC-08 §2, ticket 07」这类**括注式来处**——形状上根本不驱动任何人去做事。⇒ **R1 那把尺（找 `ticket 07`）与 R7 那把尺（找动词框架）量的不是同一件事**，这也正是「25 ≠ 132 的全部」的原因。

**结论三（把边界算成账，供编排者决定是否另立票；本腿不扩范围）**：把结论二**逐字**套到 R7 的 132 处，落到 A 柱（交付动词框架）的只有 **7 行**，其余 125 行是 `is ticket N` 的占有／来处框架、按本规矩**一律不动**。这 7 行的处置分三档，⚠ **本腿第一版在这里判错过一次并当场纠正**（详见末段）：

- **档 1｜A ∧ B ＋ 票面已结案＝一笔活被注销（3 行，已归票 228 AC#7 ①）**：`main.go:25` / `resident_windows.go:81` / `boot_windows.go:127`——就是 §1 那三枚 🔴 串。
- **档 2｜A ∧ B ＋ 票面已结案＝注释侧同形（本表 #2／#13／#14／#24／#25 五行 ＋ #22／#23 两枚错事实）**：这才是本票的真正产出，全部落在 `ticket 07` 那一族里（球的装载、gate 的订阅、GUI 构建旗标）。**它们的共同点是「那件活今天确实不在树里」，不是「注释时态不好看」。**
- **档 3｜A ∧ ¬B＝活已交付、只有时态过期（3 行，⛔ 不在票 243 范围，本腿只登记）**：
  1. `internal/llm/catalog_test.go:40`「`Protocol: "anthropic", // adapter lands with ticket 11`」——`11-llm-adapters-rest-done.md` 带 `-done`，且**活真在**：`internal/llm/anthropic/adapter.go:1-2` 逐字「Package anthropic is the Anthropic Messages adapter (C5 implementor, ticket 11)」。⇒ 只需把 `lands with` 改成过去时，**不该重指派票号**。
  2. `internal/observe/goroutine.go:25`「AST scan lands with ticket 08; a grep-level regression test runs meanwhile).」
  3. `internal/observe/nobarego_test.go:16`「The real AST scan lands with ticket 08; this grep-level check」——②③同一笔活：`08-...-done.md` 带 `-done`，AST 扫描**今天真存在**，只是在另一个包里：`tools/d22scan/main.go:98` import `go/ast`、:700-710 是 ban #1 `bare-goroutine` 的 AST 判据（:700 逐字解释「matching only FuncLit left the gate blind to」）。⇒ 同样是**时态过期 ＋ 位置写错**（承诺在 `internal/observe`、落在 `tools/d22scan`），不是缺活。

**⚠ 本腿的一次误判与纠正（写进交付物，不藏）**：第一版这里把上述 3 行判成了「A ∧ B＝三笔活被注销」，并写了「`internal/llm/anthropic/` 目录不存在」「AST 扫描从未落地」两条**未核断言**，另引用了一枚**根本不存在的文件**`internal/config/flatten.go:41` 作反例。复核动作与结果：`ls -d internal/llm/anthropic` → 目录在（adapter.go／request.go／stream.go）；`grep -rn "go/ast" tools/d22scan/main.go` → AST 扫描在；`grep -rn "func ApplyPortableOverride" internal/proc` → 真身在 `internal/proc/envfork.go:211`。⇒ **三条断言全部作废，已换成现读盘上位点。** 这个错恰好是本票根因链条的形态（「按印象报状态」，与票 64 那次四态漏一态同族），所以留在文件里当反面凭据。

**反例留在边界外、证明尺没有过宽**：`internal/audio/doc.go`「TTS output lands with ticket 26」（`26-tts-output.md` **未结案**＝真票，合法承诺）；`internal/risk/provenance.go`「wiring lands with tickets 20/21/22/26」（四枚**全部未结案**）；`internal/proc/boot_windows.go:46`「portable-mode override is ticket 06 and will use the same seam」（`06-...-done.md` 已结案，而那件活今天在 `internal/proc/envfork.go:52-53` 逐字写着「applied by DefaultLayout / ApplyPortableOverride」、函数体在 :211，`internal/config/schema.go:149` 亦有 `Portable bool` ⇒ **A ∧ ¬B，只改时态**）。这枚反例重要：它说明**同一动词框架遇到已结案票号，两种正确处理（改句子 vs 换票号）都会发生**，所以这条边界**只能人工逐处判、不能词面成门**——与票 228 AC#7 ③ 的禁令一致。

**这条边界今天算出来的净账（给编排者的硬数）**：132 处里 A 柱只有 **7 行**，其中**真正「活不在树里、票却已结案」的＝本表的 2 枚错事实 ＋ 5 枚承诺 ＋ §1 的 3 枚串＝10 处**，全部集中在 `ticket 07` 这一族；剩下 3 行（档 3）只是时态过期。**⇒ 「过期票号注释」这个形状在本仓的真实密度比 R7 的 132 这个数低两个数量级**，这正是票面「132 行人工判定，成本远大于本票收益，明确不做」的量化理由，也是不必新增仪器的理由：按 A 柱筛（1 枚 `grep -riE "arrives in ticket\|lands with ticket"` ＝ **7 行**）就能人肉看完，不需要门。

**附言：形状之二（本表唯一一处，#17）——「consumed by later tickets」**。它没有票号，所以**不在本票射程内**，也不该被任何词面尺抓到（它和结论一的 125 处「is ticket N」一样，改成带条件句的收益＝0）。登记在此只为让台账知道：**无指代未来**是另一族，量级远大于票号族，别指望这条边界能覆盖它。

---

## 5. 本腿自证与交件清单

- **改了什么**：只有本文件一枚（路径 `.scratch/wisp/probes/243/c1/census.md`，⛔ 无空格／无全角字符）。commit 带显式 pathspec、**未 push**、无 amend／reset／rebase／stash／checkout／restore／clean、无 worktree、仓内未删任何文件。
- **没碰**：任何 `.go`（零字节，含注释零改动）、`docs/reports/pending-and-issues.md`、工单任何 `- [ ]` 框、`frontend/**`、`design/**`（**零读零转述**）、`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／三枚冻结件。
- **没跑**：`go test`／`go build`（写腿 `228-r1` 在飞，按硬规矩只用 `grep`／`sed`／`Read`／只读 `git log`）。
- **没新增**：任何仪器、门、脚本、CI 步。
- **判不动的行**：**无**——25 行全部给出主判；两处给了兼类括注（#22 #23 #24）与两处次级风险登记（#4 的票 12、#17 的「later tickets」），均已写明依据。
- **R1 这把尺的一处已知漏洞（登记，不补尺＝不改任何文件）**：`internal/audio/gate.go:167` 逐字「the Muted/Unmuted events are the hook into the Muted ball state (07).」——同一枚票**省略了 `ticket` 一词**，故 R1（`grep "ticket 07"`）抓不到它，它也不在 25 行表内。它落在 #1 那一枚 `SetMuted` 的正下方 3 行，判型与 #1 相同（历史出处：`Muted` 态与事件都在）。⇒ **任何后续普查若想把这类裸编号形状算全，得另跑 `grep -rnE "\(0?7\)"`，本腿没跑、也不建议为它建门**（同一枚禁词面尺的理由）。
- **本腿的自我纠正记录（必须随件交回）**：§4 结论三第一版把 3 行判成「活被注销」，并含两枚未核断言与**一枚不存在的文件引用**（`internal/config/flatten.go:41`）；复核后全部作废并换成现读盘上位点，纠错过程留在 §4 正文内没有删除。
- **待编排者动作（本腿无权做，列出即可）**：① §0 末「票面勘误」那半句（`cmd/wisp` 全量 6 vs 去串 4）；② §4 档 3 那 3 行**只是时态过期**，按本票边界**不该另立实现票**，是否顺手改句子由写腿定；③ §3 第 1 项那对 `console_*.go`（#22／#23）——若判为缺口，「补齐 GUI 子系统构建」这笔活既不在票 243 也不在票 228 现有 AC 上，**需要一枚新票或明确判「不做并删掉这句前提」**，这是本表唯一一笔无人认领的活；④ #22／#23 的正确修法有两种（改注释成「今天仍是 console 子系统，`attachParentConsole` 走的是 :32 的真分支」vs 真加 `-H windowsgui`），**选哪种是产品决定，本腿不替它拍**。
