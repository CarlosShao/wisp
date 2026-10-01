# 228-a3 普查：托盘「允许一次」在"此刻没有卡在人等"时该长成什么样（只读，零产码）

## 起手锚点

- `date "+%Y-%m-%d %H:%M:%S %z"` 逐字读数：`2026-10-01 11:27:59 +0800`
- `git log -1 --format=%H` 逐字读数：`eed229e48cb9599c29f72fea6a1170739a5f2b41`
- 分支：`dev`（`git rev-parse --abbrev-ref HEAD` ＝ `dev`；`git status -sb` 首行 `## dev...cnb/dev [ahead 169]`，"未推枚数"是起手那一瞬的读数，会漂）
- 本腿代号 `228-a3`，只读；并发的写腿 `33-r3` 在 `cmd/wisp`，另枚只读腿 `33-a2` 在 `internal/ball` 的 STA／泵面。

---

## ① 菜单侧拿不拿得到"有卡吗"

（骨架，正文见下方补完）

## ② 三枚候选的射程与撞钉

（骨架，正文见下方补完）

## ③ 诚实 advertise 的先例与同形判断

（骨架，正文见下方补完）

## ④ 会被叫红的钉（逐枚读断言原文）

（骨架，正文见下方补完）

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

> ⚠ 本节里凡是"补完时逐字登记"的字样，都在同一枚文件的后续提交里被换成现读数；
> 若某格最终仍是空的，那代表**本腿没跑到那一发**，按 §交件判语 认它的射程，不许读成"跑过、结果为空"。

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

---

## ⑦ 判不动的地方（逐条甲／乙／不做＋现量）

> 派单纪律：**最后一枚用哪一形由编排者裁，本腿不许自选**。下面每条只列代价与现量，不写推荐。

- **F1 "有卡吗"这一维怎么进菜单构造**（① 那一跳）。现量：`showMenu` 签名 `tray_windows.go:72` 只有两枚布尔；`Ball` 上的字段由 `SetTrayChecks(muted, pausedWake bool)`（`ball_windows.go:919`）写入；`wmAppTray` 分流在 wndproc 内（`ball_windows.go:669-686`），而 wndproc 跑在 `ui-sta` 线程里。
  - **甲（推：像 `SetTrayChecks` 那样加第三枚状态位）**：改 `SetTrayChecks` 签名或新增一枚 setter ⇒ 动 `internal/ball` 的公开面＋`cmd/wisp` 的调用者；⚠ 若新增 `Events`／公开方法，会不会撞 `cmd/wisp` 那枚"名册／回调必须有执行者"的钉之外还撞别的名册钉，本腿不判（见 §④）。时效性代价：状态由别人**推**过来，推的那一次没跑＝菜单上是旧值。
  - **乙（拉：`showMenu` 里同步调一枚注入的只读闭包）**：`internal/ball` 需要多一个注入点（放在 `Options`／`Events`／还是参数上＝三选一的契约面，本腿不裁），`cmd/wisp` 那侧提供闭包读 `Replies.AwaitingHuman()`（`replies.go:252`）。代价：审批状态住在 `cmd/wisp`，读它要过 `Replies.mu`（`replies.go:105`）与 `Gate.Queue().Depth()`（`:311`），**在 ui-sta 线程里同步读会不会与那条线程互相等，本腿没量**（见 F2）。
  - **丙（不显示那一枚项，只在菜单被打开时决定要不要摆）**：＝乙 的时机变体，代价同乙，只是 `appendItem` 前多一次判断。
  - **不做**：托盘第五枚项不落地 ⇒ AC#2 那一格按票面 `:39` 判不过（常驻腿仍无"允许"入口）。现量：`recordTrayExit` 只打日志（R7），产码里 `Allow` 的托盘调用者＝0（R10）。
- **F2 ui-sta 里同步读审批状态的安全性**。现量：票面 `:92` 逐字承认"在回调里同步收口＝自等待死锁"这一族风险（`Ball.Close()` 是 `sta.PostTask(…) + <-done`），但它讲的是收口不是读；本腿**没读 `Gate.Queue().Depth()` 的锁序**，因此"读一把队列深度会不会撞上正在 prompt 的那条线程持有的锁"判不了。⇒ 要么派一枚专门读锁序的腿，要么由写腿在自己那发里量。**本腿不许用 `go test` 去量它**（派单禁令）。
- **F3 "门拦下的文案"与"权限判定的文案"是否已被复用**。现量：`Allow` 的三支文案逐字见 R4，全部带 `approval:` 前缀；本腿**未打开 `gate.go` 的 `DecideFromNative` 函数体**，也未打开 `internal/perm` 的拒绝文案。
  ⇒ 「(a) 那一形会不会撞『门拦下的文案不许复用权限判定的文案』那条硬规矩」这一格**没答**。判据：`grep -rn` 两套文案词面是否相交（本腿没跑那一发就跑完了预算的话，请由下一枚腿补，别把这里读成"已核不复用"）。
- **F4 无卡时"允许一次"该给哪种诚实**（不显示／置灰／点了出声）。三支都合法、代价不同，且**要不要出声归票面 AC#2 的"每枚按钮断一条带正确 Channel 的答复"这一句怎么读**——那是编排者的裁，不是盘上现成的答案。现量：`appendItem` 今天不会置灰（R6），`showMenu` 的返回语义里 0＝被关掉（`:70-71` 注释逐字"Returns the selected command id (0 = dismissed)"）。
- **F5 L1 那一发的枚举面**。现量：票面 `:19` 逐字写「**L1 的 `Gate.windows` 全仓零枚枚举口**」，并指名"要枚举 L1 那一发就撞票 220 的 AC#2 甲形，先裁甲乙"。
  ⇒ 若第五枚项要对 L1 出声，本票就要先造枚举口＝撞票 220 的"两腿各造一次＝两份枚举器"（票面 `:60`）。**这一格不在本腿射程内，只登记它挡住了哪一支。**
- **F6 「允许一次」拿哪枚 corr**。现量：`Allow(ctx, corr)` 需要关联号（`replies.go:316`），`AwaitingHuman()` 只给一枚卡（`:252`），`Pending()` 按显示序给全部（`:227`）⇒ 若同时有两枚卡在人等，托盘那一枚"允许一次"允许的是**哪一枚**没有定义。甲＝用 `AwaitingHuman` 的那枚（L2 优先，见 `:242-245` 注释原文）；乙＝菜单里列出全部（＝第四枚以上的新面，动 `showMenu` 结构）；不做＝本票只支持单枚并发。**本腿不选**——票面 `:39` 那句"每枚按钮断一条带正确 Channel 的答复"没说并发。

---

## 交件判语

（骨架：终态在下面补完，先钉三格边界）
- 射程：① 菜单侧状态入口的现量与代价表；② 三枚候选要动的文件／函数与今天有没有仪器看得见；③ 先例的共同形状与同形／相反判断；④ `cmd/wisp/resident_ball_228_test.go` 两枚钉的断言语义（不引行号）。
- 没答：F2（ui-sta 同步读的锁序安全）、F3（文案复用与否，未打开 `gate.go`／`internal/perm`）、F5（L1 枚举口归属）、F6（并发卡的 corr 语义）。
- 本腿**零产码、零脚本、零测试、零构建**：没跑 `go build`／`go vet`／`go test`／`scripts/build.ps1`／`go mod tidy`；只写了这一枚文件。
