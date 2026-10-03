# 228-a4 — 票 228 那 11 枚未勾格的只读普查（本腿＝`228-a4`，只读，零 Go 命令）

> 本件射程：票 228 的 11 枚未勾格，逐枚判「缺的是读数／是产码／是装配连线／待人裁定／属前端那半」，
> 并把三类（写了没接线 / 只有真机量得到 / 属前端那半）分开说。
> ⛔ 第一硬规：本腿一枚 Go 命令都没跑（`go build`/`vet`/`test`/`build.ps1`/`go mod tidy` 全零）。
> ⛔ `frontend/**`／`design/**` 两层禁令：未读、未转述、本件零引用。
> ⛔ 冻结件一字未动（本件全程只读）。AC 复选框一枚未碰。凭据值零外泄（只写变量名/字段名/blob 名）。

---

## §0 起手锚（同一发命令取，逐字）

一发命令：`date "+%Y-%m-%d %H:%M:%S%z"` ＋ `git log -1 --format="%h %ci"` ＋ `git status --porcelain -- cmd internal tools docs scripts .scratch | wc -l`

- `date` 逐字读数：`2026-10-03 09:17:38+0800`
- `git log -1` 逐字读数：`67fcc145 2026-10-03 09:17:34 +0800`（HEAD 锚点）
- `git status --porcelain -- cmd internal tools docs scripts .scratch | wc -l` 逐字读数：`358`
  （脏行枚数；三枚写腿/验收腿在飞，本腿只读，不对这个数做任何归因。）

票面复尺（本腿现跑，只读）：`grep -n "^- \[" <ticket>` ⇒ 顶层框共 16 枚＝**5 勾／11 未勾**（与编排者 09-30 17:0x 那句「现读 5 勾／11 未勾」对上）。
未勾的 11 枚行号＝`39 40 41 42 43 66 82 84 86 89 91`；已勾 5 枚＝`37 38 65 79 80`。

---

## §1 票面 11 枚未勾格逐枚表（三类分开）

> ⚠ 本腿引的 `cmd/wisp` / `internal/agent/approval` / `internal/tools` / `internal/risk` / `internal/config` 五包**此刻有写腿在飞**，
> 下面这些行的行号是**本腿 09:17 现读、会漂**；`internal/ball` / `internal/proc` 不在飞、行号稳。落地前一律现读。
> 三类分开判据：**①写了没接线**＝函数在场、生产调用点枚数＝0（本腿数到的枚数）；
> **②只有真机量得到**＝真窗口／真托盘点击／真热键／真进程退出码；**③属前端那半**＝需要一个页面侧消费者（本腿没读 `frontend/**`，只在 Go 侧点这一格的存在＋写成一句问句）。

### 逐枚结论（一行一枚，题面逐字抄进 §1.0）

**格 1｜AC#2（票面 `:39`）托盘审批按钮各映射现成枚举＋一枚真执行者**
- 缺什么：**产码＋装配连线（方向已由编排者裁死：F4 永远显示＋落空出声、F6 拿 `AwaitingHuman()` L2 优先、F7 新增第五枚 `Events` 键、F8 判据改写为"带正确 `Source` 标签"）**；不是读数、不是待裁。
- 现在盘上证据：托盘今天只有四枚 id `internal/ball/tray_windows.go:20-23`（`menuOpenPanel=1 … menuExit=4`，**没有「允许一次／拒绝／长期允许」任何一枚**）；分流 `internal/ball/ball_windows.go:674` `sel := showMenu(...)`＋`switch sel`（`:675-684` 只四支）；常驻腿 `Events` 字面量 `cmd/wisp/resident_ball_windows.go:173-184` 把 `OnTrayMute/OnTrayPauseWake` 接到 `recordBallGesture(...)`（只记不办）、`OnTrayPanel` 接到 `hs.requestPanelOpen`（真执行者）、`OnTrayExit` 接到 `recordTrayExit`（只 `outcome=ignored`，`:308-311`）。执行体侧 `Replies.Allow/Reject` 在场（`internal/agent/approval/replies.go`，本腿现读函数在 `:316`/`:372` 一带，写腿在飞会漂）；`replySurface.always` 在 `cmd/wisp/approval_always.go`。
- **①写了没接线**：托盘那三枚审批按钮**根本还没建**（零枚 tray→Allow/Reject/always 的调用点）；另 `SetTrayChecks`/`SetTrayTip`（`internal/ball/ball_windows.go:919`/`:927`）**生产调用者＝0 枚**（本腿 `grep cmd internal | grep -v _test | grep -v probes` 只命中两枚定义本身，复认 a3 R16）。
- **②真机量得到**：判据形状"点下去⇒`Gate` 收到带正确 `Source` 的答复"这一发，`winlive` 测不了（托盘选择只从 `TrackPopupMenu` 同步返回拿、不走 `WM_COMMAND`，`internal/ball/tray_windows.go:15-16`＋`:101-106`，票面 `:100` 已逐字关掉"投 `WM_COMMAND`＋id"这一支）⇒ "真点一次允许"必须由 owner 眼睛签收或另造导出面；"回调接反/接空判红"那半可走进程内注入函数值直调（照 `A495` ⓑ-1 形）。
- **③属前端那半**：这一格是**托盘（Go 侧）**，不是前端；但"卡片上的允许/拒绝按钮"另有一枚**页面侧消费者**要接（见格 5 AC#5 与本腿末的前端问句）。
- 要落地最少动哪几行：`internal/ball/tray_windows.go`（新增第五枚 id＋`appendItem` 一项）＋`ball_windows.go:675-684` `switch sel` 加三支＋`Events` 加一键（F7 甲）＋`cmd/wisp/resident_ball_windows.go:173-184` 把新键接到真执行者（调 Allow/Reject/always）＋自带一枚正控钉（接反/接空必红）。

**格 2｜AC#3（票面 `:40`）等待态真被驱动（只做 AC#6b「球真进那一态」，AC#6a 已成立不许重做）**
- 缺什么：**产码＋装配连线＋一枚新正控钉**（把驱动那一态的调用拿掉必须红——今天这一格是空的、任何突变都不会红）。
- 现在盘上证据：`bookWaitingState`（`cmd/wisp/approval_always.go:171`）函数体 `:178-184` 只做 `rt.auditf(...)`（写审计行），**一次都不碰球**；`waitingStateName()`（`:190`）导出但**生产调用者＝0**（本腿 `grep` 只命中定义）。球此刻**已在常驻腿里被起来**（`cmd/wisp/resident_ball_windows.go:165 b, err := ball.New(...)`），但没有一条线把 `Confirming`/`AwaitingApproval` 推给它。
- **①写了没接线**：`waitingStateName()`/`liveCards.waitingState()` 有名字、没使用者；`internal/ball/statevisual.go:177/:181` 那两态渲染器在，但无驱动者。
- **②真机**：球真的肉眼进"等人"那一态属真窗口，需 owner 看。
- **③前端**：面板侧那一态的可见性另属页面消费者。
- 最少动：常驻腿在卡片显示/撤回处调 `ball.SetState(statemachine.Confirming/AwaitingApproval)`（用现成 `waitingStateName()`）＋一枚"抽掉这行就红"的正控。
- ⚠ 顺带（不在 11 格内、只登记）：`approval_always.go:163-164` 注释逐字"the ball is built only by cmd/balldebug"**今天已过期**（228-r1 把球接进了 `cmd/wisp`）——本腿不改票面也不改产码，记给编排者。

**格 3｜AC#4（票面 `:41`）goroutine 收口先补旧账：`rt.close()` join `replyHandle`＋`reloadHandle`**
- 缺什么：**产码**（`close()` 里 join 两枚句柄，超时走 `context`/单调 `Duration` 不写减法，避 ban #4）＋具名写"谁 join 谁"。
- 现在盘上证据：`cmd/wisp/run.go:902 func (rt *agentRuntime) close()` → `:904 rt.replyRoot.Cancel()`、`:910 rt.reloadRoot.Cancel()`；`:896` 注释逐字**"The cancel is not a join"**；两枚句柄字段 `replyHandle`（`:318`）/`reloadHandle`（`:326`）**从不被等**；任务侧的 `bg.Wait()`（`:1106`）是另一码事。
- **①写了没接线**：不是没建，是装配处**只 Cancel 不 join**（复认票面 `:41`）。a1 §8 有限尺结论：`cmd/wisp` 测试层**零枚仪器**断这两枚被 join（只扫那一层）。
- **②真机**：join 是否真发生靠进程内测试（可测），不必真窗口。
- 最少动：`run.go close()` 在 `Cancel()` 之后 `select` 等 `replyHandle`/`reloadHandle` 的 `Done()`（带超时）＋一枚"不 join 就红"的正控；⛔ 不许为便利放宽 D38(c)。

**格 4｜AC#5（票面 `:42`）面板侧来源的 L2「允许」仍必须被判红（加了原生入口也不松）**
- 缺什么：**读数（`go test` 跑那枚冻结边界仪器到终态）＋属别人的地界（`internal/agent/approval` 的 `gate` ＋ `internal/panel` 冻结件）**——本票**不新增产码**，它是一枚防回归的负保证。
- 现在盘上证据：`internal/agent/approval/gate.go:730 func (g *Gate) DecideFromPanel` → `:740` 直接回 `ErrPanelAllow`（逐字"面板来源的「允许」被服务端 API 直接拒绝（F2 第三层）"）；冻结件 `internal/panel/l2_grant_boundary_test.go`（本腿只 `ls` 见其在位、109,415 字节、**未读内容**）在场守这一维。
- **②真机/量不到**：本腿一枚 `go test` 都没跑，**读不到"这枚正控今天绿"**；要落地这格由跑测试的腿做（`正控＝从界面方法名送一发"允许"，路由必须拒`，票面 `:42`）。
- 最少动：**无**（已实现＋已被冻结件守）；本票若动了 `DecideFromPanel`/`PanelAllow` 才算踩雷（票面 `:42` ⛔ 不许改宽）。

**格 5｜AC#6（票面 `:43`）整包终态**
- 缺什么：**纯读数，且本腿结构性量不到**——它**就是**"跑 `go test` 到终态＋逐名比红名册"那一发，必须交给能跑测试的腿。
- 现在盘上证据：票面 `:43` 给了尺（`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/... -count=1`）；历史在册 `internal/ball` 1＋`internal/panel` 4 属别人地界。
- **②真机**：本票动 Windows 原生面，POSIX 半边只能在全仓仪器里量到（票面 `:43` 原话），"编译过 ≠ 跑过"。
- 最少动：无产码；**这一格＝那一发 `go test`**（本腿零发，见 §5）。

**格 6｜票面 `:66`（AC#7「射程更正」那一框）**——题面自陈"这条是上面那一格的正文补，不是第二格"
- 缺什么：**待人裁定／记账**，**不欠产码**。它把 AC#7 的射程从"`resident_windows.go:81` 那一行"撑宽到"三枚人读字符串（`main.go` 那枚 usage／`resident_windows.go` boot 那枚 `Printf`／`internal/proc/boot_windows.go` 那枚 `slog.Info`）"；三枚的活由已勾的 AC#7（票面 `:65`、凭据 `:80` 那批 `v1` 现读）承接。
- 现在盘上证据：`cmd/wisp/main.go:24-29` 的 `const usage` 今天已是带条件事实句（本腿现读 `:25-28` 逐字"hosts the floating ball window, its tray icon and its four global hot keys in this same process …"）；`internal/proc/boot_windows.go:139` 的 `slog.Info("activation requested by second launch"...)` 也已不带"欠票 07"字样；余下 25 枚"还在承诺未来"的注释→**另立票 243**（票面 `:72`），不在本格。
- 最少动：无；要么编排者按 AC#7 已翻勾顺手勾这框、要么留作票 243 射程说明。⛔ 本腿不动任何 AC 框。

**格 7｜AC#8（票面 `:82`）`cmd/wisp/main.go:8` 那句 godoc 今天还是假的**
- 缺什么：**产码（一行注释）**——球已在同进程被装起来，这句还在承诺已结案的票 07。
- 现在盘上证据：`cmd/wisp/main.go:8` 逐字`// frozen D38(e) 10-step shutdown order. The floating ball GUI is ticket 07.`（本腿现读，与票面 `:82` 引的一致）。
- 属哪类：①②③都不属——就是一枚**注释假话**。⚠ 注意：`cmd/wisp/console_other.go:7`／`console_windows.go:26`／`notify_windows.go:12` 也还有"ticket 07"字样，但那三枚**不在本票**（票面 `:72` 明写其余注释归**票 243**）；本格只钉 `main.go:8` 这一枚。
- 最少动：`main.go:8` 改成带条件事实句（球由本文件所在常驻腿装配，见 `resident_ball_windows.go`）；⛔ 不新增扫注释票号的仪器。

**格 8｜AC#9（票面 `:84`）球起不来时那句实话被说两遍**
- 缺什么：**产码**（两声只留一位作者）＋**产码（给失败支一枚能被自然走到的注入面或台件）**；⚠ 且"失败支今天从未被自然走到"这一维**只有真机/突变现形、CI 零读数**（②）。
- 现在盘上证据：失败支 `cmd/wisp/resident_ball_windows.go:187-189` `rb.verdict = ... "this process has NO floating ball window: %v"`＋`slog.Error`＋`fmt.Printf` 已出声；boot 报告 `cmd/wisp/resident_windows.go:227`（空窗支）与 `:229`（正常支）各再 `fmt.Printf(... rb.statusLine())`，而 `statusLine()`（`resident_ball_windows.go:219-221`）在 `rb.b==nil` 时返回"this process has NO floating ball window (the ball host never ran)"⇒ **失败时同一句出现两声、成功时一声**（成功支 `:207` 把 verdict 设成"window is up"）。`ball.New` 在 `:165` 写死、无注入面（`:158` 那枚 `hooks ...ballHostHook` 是 `withPanelHost` 递面板能力、**不是**造球失败的缝，见 §2 复认 a3）。
- **②真机**：这条只有突变现形（票面 `:84` 逐字），且 `winlive` 在 `ci.yml` 零命中⇒ CI 里零读数；"console 恰好说 1 次"那枚不变式要靠跑测试/真机量。
- 最少动：两声择一（`statusLine()` 或失败支 `Printf`，另一处改成引用/删句）＋造一枚能让 `ball.New` 自然返回 err 的注入缝或台件＋自带 CI 能走的钉。

**格 9｜AC#10（票面 `:86`）球那条腿三处"我没有任务管线"自述里两处已变假话**
- 缺什么：**产码（三处逐处改带条件事实句）**；关键不是注释过期，是 `:160` 那枚**真打印进台账**的句子把"此刻属性"说成"永久属性"。⛔ 不许靠"把打印挪到装配之后"糊（会破 D38(e) 时序）。
- 现在盘上证据（行号本腿现读，已漂）：`cmd/wisp/resident_ball_windows.go:22`（注释"…pipeline and no microphone, so this file starts no state machine of its own…"）；`:211`（真 `slog` 字段 "…recorded only except the cancel key: this leg has no task pipeline and no microphone, …"）；`:257`（`const ballGestureWhy = "this process has no task pipeline and no microphone, so the gesture has no executor here; …"`）。而票 246 之后**同进程确实起了任务源**：`cmd/wisp/resident_windows.go:205 src := startResidentTaskSource(rt, ra)`⇒ `:211`/`:257` 两句作为永久声明已是假话、`:22` 关于"本文件自身不起状态机"半真半过期。
- 属哪类：①②③都不属核心——是**过期声明**；⚠ 打印时刻（球起来 `resident_windows.go:~229`）早于装配时刻（`:205`），故那一瞬不算谎、但它声明的是属性。
- 最少动：三处逐处改成带条件句（"此刻还没接任务源；任务源由同进程 `resident_windows.go` 之后装配"）；⛔ 不新增扫票号仪器；票面 `:87` 说**与 AC#8 同批改**、排 AC#2/AC#3 之后。

**格 10｜票面 `:89`（AC#2 射程被票 244 撑宽，并进 AC#2）「退出」是 GUI 进程未来的唯一停机源**
- 缺什么：**产码＋装配连线**（就是 AC#11 那枚"退出真执行者"的前置，见格 11）；外加一格**②量不到**（票 244 切 `-H=windowsgui` 之前必须先落地）。
- 现在盘上证据：唯一真停机入口＝`cmd/wisp/resident_windows.go:106-108 signal.Notify(bootExit, os.Interrupt, syscall.SIGTERM)`（消费点 `:225 case sig := <-bootExit`／`:231 reason = rt.RunEventLoop()`）；`recordTrayExit`（`cmd/wisp/resident_ball_windows.go:308-311`）只打 `outcome=ignored`⇒ **今天托盘「退出」不通电**。复认票面 `:89`。
- **②量不到**：`buildWispForTest`（`cmd/wisp/secret_argv_windows_test.go:174`）逐字 `go build -o exe ./cmd/wisp`、**零 `-ldflags`**⇒ 台件二进制恒 CUI，真机侧看不到 GUI 子系统那一支（a2 F4 本腿亦未现读全共用者名册）。
- 最少动：并入 AC#11 的乙-1 落地；⛔ 编排者 `:89` 已写"用哪一条形由我在 `244-a3` 读数回来后裁"——那条裁**已下**（`:91` 裁乙-1），故此格不再欠裁。

**格 11｜AC#11（票面 `:91`）「退出」真执行者只能长成"记一笔请求、立刻返回，由常驻腿自己结束环路"（乙-1）**
- 缺什么：**产码＋装配连线＋两枚新钉**；**方向已由编排者裁死＝乙-1（`internal/proc` 的 stop-request 位），⛔ 协程形（乙-2）被三读数否掉、且它要动 D38(b) 名册＝人工批准**。⚠ 有一子格踩"待人裁定"：若落地发现必须动 `hookableRoster` 闭集八枚，就从"加触发"滑成"改契约"，**当场停手上报**（票面 `:95`）。
- 现在盘上证据：`internal/proc/boot_windows.go:127 RunEventLoop` 唯一定义、唯一停机输入 `:128 signal.NotifyContext(os.Interrupt, SIGTERM)`；全仓 `RequestStop|StopRequest|shutdownRequest|stopRequest` **零命中**（本腿现跑，复认 a2 R3）⇒ AC#11 要的位是**全新面**。有序退出全挂在 `runResident` 正常返回的 defer 上（`cmd/wisp/resident_windows.go:70 defer sink.close()`／`:72-82 defer rt.Shutdown(false)`／`:137 defer rb.stop()`／`:144 defer ra.detachBall()`）⇒ 只需让 `:231 rt.RunEventLoop()` 返回。`recordTrayExit` 注释 `cmd/wisp/resident_ball_windows.go:300-307` 逐字自称"两条补法都是新契约、AC#1 不授权任一"——票面 `:94 ⓒ` 判这句**不准确**（协程形单独完不成停机），落地随 AC#10 改口。
- **②真机量得到**：完成判据 ⓑ-1（进程内十步走满）可测（照 `cmd/wisp/resident_approval_246_windows_test.go:145-158`）；ⓑ-2（"那一次触发导致进程以退出码 0 结束＋审计仍有 install 记录"）**只准钉这一半**——"审计 1..10 全序"在真机盘面上**结构性读不到**（`internal/proc/shutdown.go` 只为 skipped/失败/fast 出声 `:130/:145/:147/:172`，成功步与第 8/10 步零日志，a2 D-3；本腿复认 `shutdown.go` 行号在 `internal/proc` 未飞）。⛔ 投 `WM_COMMAND`＋菜单 id 模拟托盘选择那一支**已被票面 `:100` 逐字证伪**（`tray_windows.go:15-16`）。
- **③前端**：不属（这是 Go 侧停机）。
- 最少动：`internal/proc` 加一枚 stop 位并**同时接进 `RunEventLoop` 内 `:145` 那个 select 与 `cmd/wisp/resident_windows.go:195-203`/`:106-108` 空窗 select**（票面 `:94 ⓐ`：只接环路会吞掉"球起来后、环路进入前"那一下）；`recordTrayExit` 改"记请求即返回"（不 `rt.Shutdown`＋`os.Exit`＝跳 defer，票面 `:92` 反面判据）；ⓑ-1 进程内钉＋ⓑ-2 `winlive` 钉；boot_windows.go:122-126 触发名册 doc 改口＋落一枚具名 `A##`（＝`A493`）。

---

## §2 复认／推翻 a1–a3 三节（具名说哪条被推翻）

（取数中）

---

## §3 停机那一格（AC#11）与「允许一次」那一格（AC#2／F4）的最少改动面

（取数中）

---

## §4 我可能写错的条目（自我对抗，不许留空）

（取数中）

---

## §5 量不到的地方（具名，⛔ 不推测填空）

（取数中）

---

## §6 交件判语

（取数中）
