# 票 246 — 落地腿 `246-r1` 交付表：最短链「注入 gate → 真卡片 → `Esc` 否决一次并归还 → 退出序列不撒谎」

> ⚠ **这份表的身份**：出自**实现方本人**（`246-r1` 是产码腿）。按 `AGENTS.md §1.5`／`SPEC-12 §4.3`，
> `docs/evidence/s1/` 的裁决表**必须出自非实现者**才算判语；本文件全部格子的性质是
> **〔实现方自述＋台件〕**，等另一程对抗验收裁，**本表不构成任何一格 AC 的翻勾依据**（AC 框归编排者，本腿一枚未动）。
>
> 判据出处＝`.scratch/wisp/issues/246-the-process-that-runs-tasks-is-not-the-process-holding-the-ball-so-no-veto-channel-has-an-executor.md`
> 「判据（本票现行射程）」AC#0..AC#6。前段代价表＝`.scratch/wisp/probes/246/a1/census.md`（206 行）；
> 编排者裁定＝乙形（装配根注入）＋"钩子注册入口不算契约面"，账 `A481`。
> 取数一律现跑，落盘在 `.scratch/wisp/probes/246/r1/`（只建不删）。

---

## §0 主张与起手名册

### §0.1 这一程主张的产品事实（只这一条，别的都不主张）

**双击图标起来的那个常驻进程，今天有了一枚真审批门；一张真待决卡片能在那里挂起来；`Esc` 在 `Confirming`
期间真能否决它一次并把键归还给桌面；卡片挂着时收到退出信号，D38(e) 第 3 步不再被记成 skipped，
卡片以拒绝收口并留下审计。**

⛔ 本表**不主张**（每一条都是既成事实或射程外）：
- 不主张"用户看得见卡片"——面板今天开不出来（`cmd/wisp/approval_always.go:165` 逐字 this binary links no WebView2 host）；
  今天能证的只有**日志＋状态位＋Win32 那一层的读数**。
- 不主张"`Esc` 安全了"——确认那 2–3 秒仍把键从别的程序那里借走（票 245 乙形的具名残余，`SPEC-12 §5` 已有一行，账 `A479`）。
- 不主张四条否决通道通了——只落 `Esc` 一条（AC#6）。
- 不主张常驻进程能跑任务——这里没有 agent loop、没有 provider、没有麦克风、没有 `config.toml` 接线。

### §0.2 起手名册与锚点（闸门＝终态等于起手那一刻的名册，不是"必须为空"）

| 项 | 读数 |
|---|---|
| 取数时刻 | `2026-09-30 18:42:54 +0800`（`date` 现跑） |
| 起手锚点 | `ce1ade1f` — `ledger(A481 收 246-a1：裁走乙＝装配根注入，票 246 判据落成最短链六格…)`（`git log -1` 自取，⛔ 未采用派单里给的任何 sha） |
| 分支 | `dev` |
| 起手名册全量 | `.scratch/wisp/probes/246/r1/roster-start.txt`（`git status --porcelain` 原文，取数即存） |
| `git status --porcelain` 行数 | **178**（含本腿起手前已建的 `.scratch/wisp/probes/246/r1/`） |
| 起手产码面 | `git status --porcelain -- internal cmd`＝**0 枚**（同发取数）⇒ 本腿起手时 `cmd/wisp`／`internal/proc`／`internal/ball` 写面是空的 |
| 别人的脏件 | `design/**` 的删除与改动、`.gitignore`、`docs/evidence/s1/152-…-accept-r1.md`、各票 `.scratch/wisp/probes/**` 台件、`part*-*.txt` 等 ⇒ 本腿**一枚不碰、不删、不提交** |
| 本腿写面（白名单） | `cmd/wisp/**` ＋ `internal/proc/**` ＋ `docs/evidence/s1/246-*.md` ＋ `.scratch/wisp/probes/246/r1/**` ＋ 票 246 面「Progress log」末尾一行 |
| 桌面安静（跑真机前） | `tasklist //FI "IMAGENAME eq balldebug.exe"`＝**0**／`wisp.exe`＝**0**／`esclistener.exe`＝**0**／`go.exe`＝**0**（19:0x 现跑，两档：起手与 winlive 前各一次） |
| 终态闸门算法 | `diff <(git status --porcelain) .scratch/wisp/probes/246/r1/roster-start.txt` 的差集只允许落在白名单写面内；见 §6 |

### §0.3 起手现跑（写第一枚产码之前）

| 尺 | 命令 | 读数 |
|---|---|---|
| 起手依赖图（AC#1 的"前"） | `GOOS=windows go list -deps ./cmd/wisp` | **276 行** → `deps-cmdwisp-before.txt` |
| 起手逐包 import（"前"） | `GOOS=windows go list -f '{{join .Imports "\n"}}' ./cmd/wisp ./internal/proc ./internal/ball` | **87 行** → `imports-before.txt` |
| **AC#5 起手红归因（第 1 发）** | `PATH="$PWD/third_party/sherpa-onnx:$PATH" go test ./cmd/wisp -count=1` | `ok github.com/CarlosShao/wisp/cmd/wisp 143.600s`，`--- FAIL` **0 枚** → `ac5-start-default-tier.txt` |
| AC#5 复跑（第 2 发，`-v`） | 同上 ＋ `-v` | `ok … 150.544s`，`=== RUN` **200 枚**，`--- FAIL` **0 枚** → `ac5-repeat2-v.txt` |

---

## §1 判语（逐格）

> 姿势：判据 → 落地形状（`file:line`）→ **现跑读数** → 自判。⛔ 全部是〔实现方自述〕。

### AC#1 装配根注入一枚真审批门到常驻腿，零新增包级依赖边 — **自判：成立**

- 判据：由 `cmd/wisp`（装配根）把 gate 递进常驻那条腿；⛔ 不开任何新包级依赖边；尺＝import 差集落地前后**逐名相同**。
- 落地形状（乙形，照 `A481` 裁定）：
  - `cmd/wisp/resident_windows.go:122` `ra := newResidentApproval()` — gate 在装配根建，不在常驻文件里建；
  - `:135` `rb := startResidentBall(rt.Registry, ra.vetoByEsc)` — 递给球宿主的是一枚**函数值**（`escVetoFunc`，`resident_ball_windows.go:61`），球宿主不知道审批是什么；
  - `:144` `ra.bindBallHost(rb)`（`:143` 之前先 `defer ra.detachBall()`，LIFO 保证球窗口收掉之前先解除绑定）；
  - `:156` `rt.RegisterShutdownHook(proc.StepCancelTasks, ra.cancelTaskRoots)`；`:164` 启动报告读 `rt.RegisteredShutdownSteps()`。
- **差集读数（两把尺，同一命令前后各一发）**：

  ```
  GOOS=windows go list -deps ./cmd/wisp        before 276 行 / after 276 行   diff rc=0（零行差异）
  GOOS=windows go list -f '{{join .Imports "\n"}}' ./cmd/wisp ./internal/proc ./internal/ball
                                               before  87 行 / after   87 行   diff rc=0（零行差异）
  ```
  落盘 `deps-cmdwisp-before.txt`／`deps-cmdwisp-after2.txt`／`imports-before.txt`／`imports-after2.txt`。
  ⚠ 中途出现过**一枚**差集非空的读数并已被消掉：`internal/proc` 一度 import 了 `sort`（`imports-after.txt` 88 行），
  已改成按声明顺序走固定名册、零新增 import（`shutdown_hooks.go:148` 的注释具名写了为什么）。
- 出厂进程那一格（不是包内自证）：`TestAC246ShippedResidentProcessOwnsItsCancelStep` 用无参数起真 `wisp.exe`，
  读它**自己的控制台**：同时含 `审批门已装配进本进程` 与 `3:cancel-task-roots`；读它**自己落盘的日志**：13 条记录里含
  `resident-approval: 退出第 3 步完成…`。删掉装配那两行，这两处读数一起消失（突变 M1/M2，见 §2）。

### AC#2 一张真卡片能在常驻进程里挂起来 — **自判：成立**（可见性只到"日志＋状态位＋Win32 读数"那一层）

- 判据：用 `internal/agent/approval` 既有机制造**真**待决项、状态机真进 `Confirming`；⛔ 不许 mock 代真；⛔ 不许写"用户看得见卡片"。
- 落地形状：`cmd/wisp/resident_approval_windows.go:197 askConfirmation` 走的是 `gate.AdmitTextTask`（D47 登记）＋
  `gate.PendingWindow`（L1 窗口，`gate.go:246`）／`gate.PendingApproval`（L2 入 C18 队列，`gate.go:496`）；
  注入的面是 `:417 ballCardUI.Prompt`（真单次令牌由 `Replies.Record` 记账，卡片态由 `:490 stateForCardLevel` 取
  `Confirming`／`AwaitingApproval` 两枚 D43 名字之一，⛔ 不造第五枚名字）。
- **现跑读数（winlive 真机，`mut` 之外的干净树，`winelive-246-first.txt`）**：
  - 状态位：`ra.cards.WaitingState() == Confirming`（L1 卡片挂起期间），`AwaitingHuman()` 返回该卡；
  - 日志：`wisp: 卡片挂起：L1 resident.confirmation（编号 host:live246-esc）`；
  - Win32：`HotkeyReport().Live()` 由 3 枚变 **4 枚**（cancel 位在窗口期内真被注册），`EscTakenOver()==true`；
  - 收口：`answer = veto`，审计逐字
    `wisp: [audit] approval: ANSWER-VETO corr=host:live246-esc tool=resident.confirmation channel=esc decision=veto`。
- **反证（⛔ 不是 mock 的证据）**：没有球窗口时同一张卡**挂不起来**，且是被**真 gate** 拒的——
  `TestAC246CardWithNoWindowFailsClosedThroughTheRealGate` 读到 `answer=reject` 且理由逐字含
  `确认界面不可达（常驻进程没有悬浮球窗口，卡片无处呈现）`：这句由 `gate.go` 的 fail-closed 分支产出，
  只有真的调用了注入的 UI 才可能出现。
- **诚实残余（本腿具名，不等验收腿来挖）**：`AskOnTaskRoot` 在产码里**今天没有生产调用者**（只有本包用例）。
  任务源（真跑出来的卡片）要等票 228 的 `config.toml` 接线与语音链——**这一格交的是"门在、卡进得来"，不是"卡会自己来"**。

### AC#3 `Confirming` 期间 `Esc` 真能否决一次、完事归还 — **自判：成立**（三形读数齐）

- 判据：⚠ 前置＝票 245 AC#9 那套"第二个进程观察 `Esc`"的台件**必须先进仓**；完成判据＝**借到／否决真生效／归还后另一进程真收到** 三形都在；⛔ 不许只测"借"。
- 台件已进仓：`cmd/wisp/testdata/esclistener/main.go`（源头是票 245 验收侧
  `.scratch/wisp/probes/245/v1/esclistener/main.go`，本腿把接口改成**由 stdin 触發注键**——L1 窗口只有 2–3 秒，
  原来那套"起来 400ms 后自己注一发"的秒表赌法在这条链上不可靠）。放 `testdata/` 的理由具名：
  go 工具的包走查与 `tools/d22scan` 的走查都按规则跳过 `testdata`（`main.go:660` 那条），用例显式
  `go build ./cmd/wisp/testdata/esclistener` 造它（实测 rc=0、产出 2,673,152 字节 exe，不落仓）。
- **三形＋正控读数**（`winelive-246-first.txt`，观察者是自己建窗、自己抢前台的**第二个真进程**）：

  | 形 | 观察者命令行 | 逐字读数 | 判 |
  |---|---|---|---|
  | 尺的正控（先跑） | `-steal` | `READING keydown_esc=0 syskeydown_esc=0 wm_hotkey=1 other_keys=2 foreground_is_mine=true steal=true` | 一枚全局裸 `Esc` 注册确实吃掉投递 ⇒ 这把尺在"收不到"那一头不瞎 |
  | 基线（球稳态） | `-watch` | `READING keydown_esc=1 … wm_hotkey=0 foreground_is_mine=true steal=false` | 别的程序**平时**收得到 |
  | **借到** | `-watch`（卡片挂起期间注一发） | `READING keydown_esc=0 … wm_hotkey=0 … steal=false` | 那发被 **Wisp 的球**吃走 ⇒ 真的借到了（`foreground_is_mine=true` 已断言，否则整形作废） |
  | **否决生效** | 同上那一发 | `answer=veto` ＋ `ANSWER-VETO … channel=esc decision=veto`（上一格逐字） | 不是"我们记了一笔"，是那张卡真的被那次注键判了否决 |
  | **归还** | `-watch`（卡片收口之后） | `READING keydown_esc=1 … wm_hotkey=0` | 另一个程序又收到 Esc ⇒ 归还不是自家记账 |
  - 同一发用例的自述行逐字：`AC#3 READING: borrow/veto/return in 132.0591ms; veto=veto; observer keydown 1->0->1 (steal control 0)`。
  - 用例名：`TestLive246ConfirmingCardBorrowsEscVetoesAndReturns`（`--- PASS (6.59s)`）。
- ⛔ 只测"借"的形状在本腿不存在：借／否决／归还三形是**同一枚用例里三次独立的观察者进程**，
  摘掉任何一步都会红（§2 的 M2／M3／M4）。
- 顺带（⛔ 不是本票的格，只是本腿看见并写下来）：这套台件进仓后，票 245 AC#9 的"仓内一条能跑的用例做同样的四形"
  **有了载体**；那一格翻不翻勾归编排者与票 245 的裁决腿，本腿不替它判。

### AC#4 卡片挂着时收到退出信号不许撒谎 — **自判：成立**

- 判据：`internal/proc/boot_windows.go` 里 `Shutdown` 今天逐字 `hooks := ShutdownHooks{}` 只填 `CloseJob` ⇒ 第 1–7 步零生产者；
  把"取消任务根"做成**真注册的钩子**；一发真机：卡片挂起 → 发退出 → 那一步 `StepRecord` **不再是 skipped**、
  待决卡片被判**拒绝＋留下审计**。⛔ 十步顺序一个字不许动。
- 落地形状：
  - 注册入口（灰区由 `A481` 裁过：冻结的是**顺序**，注册入口是留给后续票的）：
    `internal/proc/shutdown_hooks.go:81 Register`（拒空钩子／拒第 8、10 步这种没有钩子位的／拒同位第二个所有者／
    `:134 Hooks()` 一取即封闭）＋ `internal/proc/boot_windows.go:188 RegisterShutdownHook`／`:198 RegisteredShutdownSteps`；
  - `boot_windows.go:167` 读注册快照，**顺序仍由 `RunShutdownSequence` 拥有**（该函数一字未动，
    `TestShutdownOrderAudit`／`TestShutdownFastPathSkipsOnlyStep7`／`TestShutdownNilHooksAreRecordedSkipped` 三枚原钉仍绿）；
  - 钩子本体＝`cmd/wisp/resident_approval_windows.go:275 cancelTaskRoots`：先把待批 L2 卡片走**队列自己的拒绝漏斗**
    （`Replies.Reject`→`DecideFromNative`→`q.reject`，写的就是原生「拒绝」那枚 `ANSWER-REJECT` 行），
    再取消任务根，再有界等（`ctx` 到点即"弃等并记录"，⛔ 不悬挂、⛔ 不写墙钟超时）。
- **一发真机读数（winlive，`TestLive246ExitRefusesAHangingL2Card`，`--- PASS (0.08s)`）**：
  卡片挂着（`AwaitingHuman().Level=="L2"`，`WaitingState()==AwaitingApproval`，且**没有**借 `Esc`——300 秒的期限不配借全局键）
  → `rt.Shutdown(false)` → `records[2].Skipped == false`、`Err == nil`、`Name == "cancel-task-roots"`，
  审计逐字 `wisp: [audit] approval: ANSWER-REJECT corr=approval-1 tool=resident.confirmation decision=reject reason="常驻进程退出：未获批准，按拒绝处理（未执行）"`，
  用例自述行逐字 `AC#4 READING: step3 executed in 10.0134ms, card answered reject, ANSWER-REJECT booked; no-skipped=true`。
- **另一形读数（L1 窗口没有「拒绝」这个动词，SPEC-06 §2 B1）**：
  `TestLive246ExitAbandonsAHangingL1Window` `--- PASS (0.06s)`，审计逐字
  `approval: RESIDENT-WINDOW-ABANDONED corr=host:live246-exit-l1 tool=resident.confirmation decision=reject reason="常驻进程退出，L1 确认窗口未放行，未执行" channel=none`
  （这一支 gate 自己不写行，所以由本文件写，`channel=none` 明说**没有人否决过**）＋ 步 3 `Skipped==false`。
- **出厂进程那一形**：见 AC#1 末（无参数 `wisp.exe` 落盘日志含 `退出第 3 步完成`，控制台含 `3:cancel-task-roots`，
  并以退出码 0 走完 `10 steps, 0 failed`）。
- **"不冒出假账"那一半**：winlive 用例断言 `HotkeyReport().Problems()` 收口后为 **0 枚**（归还后 cancel 位回 standby，
  不是"注册失败"），并且默认层不会把"跳过"当失败：`shutdown step skipped (module not present)` 那六枚是本票不拥有的步，
  老实记着，⛔ 没为变绿去把它们伪造成 executed。⚠ 派单/票面那句 `problems()` 到底指哪一处，本腿按"热键报告不得冒假问题行"判的；
  若指别处，这一格的射程就没盖到 —— 已列入 §5 判不动。

### AC#5 起手第一发先归因那枚红 — **自判：成立（结论是"起手这一发没有红可归"）**

- 判据：⛔ 不许当已知常红略过；⛔ 不许顺手修不属于本票的东西。
- **起手读数（动任何产码之前，锚点 `ce1ade1f`）**：
  `PATH=… go test ./cmd/wisp -count=1` → `ok github.com/CarlosShao/wisp/cmd/wisp 143.600s`，`--- FAIL` 0 枚；
  复跑 `-v` → `ok … 150.544s`，`=== RUN` **200 枚**，`--- FAIL` **0 枚**。
  ⚠ 已排除"根本没跑"那一形：那一形会是 `exit status 0xc0000135` ＋ `0.0xxs` 且无 `--- FAIL`；本次是 143.6s／150.5s 真跑完。
- **归因结论**：台账 `A480` ④ 那枚"1 红 2 绿、用例名丢了"的红，在本腿起手锚点、桌面无同类进程、无并发 `go test` 的条件下
  **未复现（2/2 绿）**。⇒ 既**不是**"起手即在的稳定红"（这两发不支持），也**不可能**是"本票造成"（这两发之前本腿零改动，
  `git diff --name-only` 当时为空）。剩余形状（`cmd/wisp` 大量用秒级 `waitFor` 预算的时序抖动）
  由 `245-v1` 自己具名写过环境差（表 `245-esc-not-a-standby-global-hotkey-v1.md:184`："第 1 发期间我在 17:56 编译过一次台件、CPU 有争用"），
  **本腿不替它定案**，只把后续样本交回（§3 的每一发默认层读数）。
- ⛔ 本腿**没有**去修任何东西（那枚红不在本票射程）。本腿唯一一枚自己造成的红（票 131 那道腿门的"名字融合"）在 §2 末与 §3。

### AC#6 其余三条否决通道不许顺手做 — **自判：成立**

- 判据：单击球 veto／KWS 语音否决词／面板拒绝本票一律不做；完成判据＝票面「Progress log」具名写"只落 `Esc` 一条，另三条各有归口"。
- 现跑读数：`newResidentApproval` 用 `approval.NewChannels()`（**四枚全未加载**起步），只有 `bindBallHost` 在
  **Win32 真的给了窗口**之后 `SetLoaded(ChannelEsc, true)`；默认层用例逐枚断言
  `ChannelBall`／`ChannelKWS`／`ChannelPanel` 在本进程 `Loaded()==false`，且 `gate.Veto` 在未加载通道上返回
  `ErrChannelUnavailable`（判定权在 registry，不在本腿写的分支）。`OnClickBall` 仍走 `recordBallGesture`、
  `recordTrayExit` 仍是"只记账"，`resident_approval_windows.go` 里没有任何指向面板／KWS 的路由。
- 票面 Progress log 那一行：见本表提交后的票面末尾（本腿只追加，⛔ 未碰任何 `- [ ]` 框）。

### 通道名册（收口自查）

| 通道 | 本进程 loaded | 执行者 |
|---|---|---|
| `ChannelEsc` | 有球窗口 ⇒ **true**；无窗口 ⇒ false | 本票落的：`OnCancelHotkey` → `vetoByEsc` → `Replies.Veto` → `Gate.Veto` |
| `ChannelBall` | false | 无（票 228 后续片／D43 同一条） |
| `ChannelKWS` | false | 无（票 41，`approval.go:99-106` 的 DEFERRED 行） |
| `ChannelPanel` | false | 无（票 37；本树无 WebView2 宿主） |

---

## §2 突变与正控名册

姿势：每发先确认锚点**恰好 1 处**，落盘，跑指名的用例（判红绿只认 `--- FAIL`），
跑完 `git cat-file blob HEAD:<path> > <path>` ＋ `git diff --quiet -- <path>` 复认。⛔ 未用 `-overlay`、⛔ 未用 `checkout`／`stash`／`reset`。
驱动与逐发读数：`.scratch/wisp/probes/246/r1/run_mutations.py`／`mut-M*-*.txt`／`mutation-summary.json`／`mutation-driver.log`。

| 编号 | 攻哪一格 | 突变 | 指名用例 | 红？ | 还原 |
|---|---|---|---|---|---|
| M1 | AC#1/AC#4 装配 | 摘掉 `rt.RegisterShutdownHook(StepCancelTasks, …)` 那一行 | `TestAC246ShippedResidentProcessOwnsItsCancelStep`（默认层，出厂进程）＋`TestLive246Exit*` | 见下表 | 见下表 |
| M2 | AC#1 注入 | `startResidentBall(reg, ra.vetoByEsc)` → `… , nil` | `TestLive246ConfirmingCardBorrowsEscVetoesAndReturns` | | |
| M3 | AC#2/AC#3 借 | 摘掉 `b.TakeEscForCancel()` | 同上（(a) 形） | | |
| M4 | AC#3 还 | 摘掉 `b.ReleaseEscAfterSession()` | 同上（(c) 形） | | |
| M5 | AC#4 空钩子 | 钩子照注册、照执行，但**什么都不做** | `TestLive246Exit*`＋出厂进程那一格 | | |
| M6 | census D-5 | 没有球窗口也把 `ChannelEsc` 标 loaded | `TestAC246EscChannelStaysUnloadedWithoutABallWindow` | | |
| M7 | proc 读注册 | `Shutdown` 读了注册面又丢掉（等于回到"只有 CloseJob"） | `TestRuntime*`（internal/proc）＋`TestAC246CancelStepHookRunsOnTheRealShutdownSequence` | | |

> ⚠ **本腿踩到并当场纠正的一枚空尺（具名，因为它曾造出七发"全绿"的假读数）**：
> 第一遍突变驱动用 `subprocess.run(..., shell=True)`，本机 shell 是 cmd.exe ⇒ `-run 'TestAC246'` 里的**单引号被逐字交给 go**，
> 七发全部 `testing: warning: no tests to run` ＋ `ok … [no tests to run]` ＋ **rc=0**。
> 若按"零红"入账，本表会写成"七发突变全部没红 ⇒ 用例不守"，而真相是**用例根本没跑**。
> 处置：驱动改成参数列表直调（不过 shell）＋ 把"出现 no tests to run"判为**驱动失败**而不是绿，
> 旧读数作废并保留在 `mut-*-*.txt`（会被第二遍覆盖，覆盖前已在 `mutation-summary.json` 留第二遍结果）。
> 沿用本仓第 70 条：**凡写进表的尺，先跑一遍，读不到期望就先怀疑尺。**

### §2.1 第二遍（修正驱动后）逐发红句

（取数进行中；本节的每一行只写自己跑出来的东西。）

---

## §3 门禁读数

| 尺 | 命令 | 时刻 | rc | 读数 |
|---|---|---|---|---|
| 默认层（起手，产码前） | `PATH=… go test ./cmd/wisp -count=1` | 18:42 | 0 | `ok … 143.600s`；`--- FAIL` 0 枚 |
| 默认层（起手复跑 `-v`） | 同上＋`-v` | 18:47 | 0 | `ok … 150.544s`；`=== RUN` 200 枚；`--- FAIL` 0 枚 |
| build | `go build ./...` | 19:0x | 0 | 无输出 |
| vet | `go vet ./...` | 19:0x | 0 | 无输出 |
| vet（winlive 一档） | `go vet -tags winlive ./...` | 19:0x | 0 | 无输出（含本腿新增的 winlive 用例文件） |
| vet（台件，显式路径） | `go vet ./cmd/wisp/testdata/esclistener` | 19:0x | 0 | 无输出。⚠ `./...` 走查按规则**不吃** `testdata`，所以这一档必须点名跑 |
| gofmt | `gofmt -l internal/proc cmd/wisp` | 19:2x | 0 | 空（唯一一次报本腿文件是 `resident_approval_246_windows_test.go`，已 `gofmt -w` 之后复跑为空） |
| internal/proc 默认层 | `go test ./internal/proc -count=1` | 19:0x | 0 | `ok … 0.937s`（含本腿 4＋3 枚新用例；首跑曾红在**本腿自己的用例** `TestShutdownHookSetZeroValueIsUsable`——先 `Hooks()` 又 `Register`，是测试写错，不是产码；已改正并复跑绿） |
| 四包默认层 | `go test ./cmd/wisp ./internal/proc ./internal/agent/approval -count=1` | 19:1x | 0 | `ok cmd/wisp 137.099s`／`ok internal/proc 0.967s`／`ok internal/agent/approval 0.340s`；`--- FAIL` **0 枚** |
| winlive（本腿三枚真机用例） | `go test -tags winlive ./cmd/wisp -run TestLive246 -v` | 19:1x | 0 | 3 枚 `--- PASS`；三形读数见 AC#3 表；`FAIL` 0 枚 |
| 依赖边两把尺 | `go list -deps`／逐包 import 差集 | 19:0x | 0 | 276=276、87=87，`diff` 零行（AC#1） |
| d22scan | `scripts/d22scan.sh` | 19:3x | （收尾补跑，见下一行） | 第一步是正控：正控红则真扫描根本不跑 |
| 全仓 `go test ./...` | 未跑 | — | — | **该档未读**（本腿写面是 `cmd/wisp`＋`internal/proc`，其余包的读数由编排者的整包门负责；本腿不冒充跑过） |

> ⚠ 本腿唯一一枚**自己造成**的红（已修，且没去动别人）：改名之前
> `TestAC4EveryLegIsNailedOrRuled`（票 131 的"每条腿要么有钉、要么有裁定"门）红了一发，
> 逐字：`AC#4 RED: leg "panel-inbound" (main.go:113) books records (slog.Info@resident_approval_windows.go:129 | run.go:811)
> through the process logger, installs no listener, and carries no WISP-LEG-SINK-RULING: sentence.`
> 根因不是那条腿丢了钉，是那道门**按函数名**走调用图、并把同名声明**并集**，而本腿新文件里那枚方法恰好也叫 `auditf`
> ⇒ 本文件的记账被算到了 panel-inbound 那枚腿上。
> 处置：把本文件与 `run.go`/`resident_ball_windows.go` 同名的方法全部改名（`auditf`→`residentAuditf` 等 8 枚），
> ⛔ 不碰那道门、⛔ 不给 panel-inbound 写裁定、⛔ 不删任何断言。复跑该门与全包均绿。

---

## §4 残余与归口（本腿做了的、没做的、该归谁的）

| # | 事项 | 归口 |
|---|---|---|
| R-1 | D38(e) 第 1／2／4／5／6／7 步仍零生产者（本票只吃第 3 步） | `246-a1` D-1，`A481` 已分：step 4 语音链、step 5 票 15/28、step 6 票 33/35、step 1/2 scheduler/hotkey 全量＝别的落地面 |
| R-2 | `AskOnTaskRoot` 零生产调用者 ⇒ 常驻进程今天仍"有门无任务" | 票 228 后续片（`config.toml` 接线）＋语音链；本腿不擅自搬 loop/provider 进常驻腿 |
| R-3 | 常驻腿整棵 `config.toml` 没接、托盘「退出」无执行者、球位置不持久、`ball.New` 失败路径不清理 | 票 228 后续片（`246-a1` D-4／D-5） |
| R-4 | 过期注释两处（`cmd/wisp/main.go:8`／`:24` 那族句子） | 票 228 AC#8／票 245 AC#6①（账 `A477`／`A478`）。⚠ 本腿确实让 `resident_ball_windows.go` 那三句"本腿没有审批门"变成假话，**已在原地改成带条件的事实句**（`:22`、`:197` 与 `gestures` 字段），没留给别人 |
| R-5 | 面板拒绝 `ChannelPanel` 执行者（面板侧要什么） | ⛔ `frontend/**`／`design/**` 零读零转述；只写进票面与台账，由 owner 自己带给他的前端 agent |
| R-6 | KWS 否决词 `ChannelKWS` | 票 41（`approval.go:99-106` 的 DEFERRED 行） |
| R-7 | 单击球否决 `ChannelBall` | 票 228 后续片（本票 AC#6 明令不做） |
| R-8 | 借 `Esc` 那 2–3 秒仍吞别处的 `Esc` | 票 245 的具名残余（`SPEC-12 §5` 已有一行，账 `A479`），本票**没有**改变它，只是让那 2–3 秒在常驻进程里第一次真的会出现 |
| R-9 | D43 第 22 行的否决目标态是 `Acting`；本腿收口后回 `Sleeping` | 见 §5 第 3 条（判不动，等裁） |
| R-10 | 票 245 AC#9 的"台件进仓"这一步**顺带满足了它的载体条件** | 翻勾归编排者与票 245 的裁决腿；本腿不替它判，也不碰它的框 |

---

## §5 判不动的地方（本腿不自裁，具名交回）

1. **"常驻进程是不是任务管线的永久之家"**：`246-a1` §6 具名判不动（D2 正文 `PLAN.md:73` 未逐字核，本腿也没核）。
   本腿按票面最短链做：注入 gate ＋ 一张真卡 ＋ 一次 `Esc` 否决 ＋ 退出不撒谎，⛔ 没把 agent loop／provider／config 搬进常驻腿。
   若产品其实想让 veto **跨进程** reach 到 `wisp run`（"跑任务的进程"与"带球的进程"仍是两个），本腿的形状要重议。
2. **"卡片由谁在什么产品动作上举起"**：本腿只有 host-initiated 一发（与 `confirmModeSwitch`／`runWidening` 同形）。
   真任务源出来的卡片属票 228 之后。⇒ 任何人把这一格读成"用户会看到卡"都超出证据。
3. **收口后回到哪一枚状态**：D43 第 22 行 `Confirming` 经否决的目标是 `Acting`（取消该调用并回给 LLM），
   而这条腿没有 loop 可回。本腿选择回 `Sleeping`（`resident_ball_windows.go:83` 那句"Sleeping 是本腿唯一可主张的态"的延伸），
   因为假装 `Acting` 就是票 228 表头警告过的那类"声称自己没有的路径"。
   ⚠ 但 `Confirming → Sleeping` 这条转移**不在冻结表上**，是否算越界，本腿判不动，等裁。
4. **AC#4 里那句 `problems()` 的所指**：本腿按"球的热键报告不得冒假问题行"来判（winlive 断言 `Problems()==0`）。
   若原意指别处（例如退出审计里"有产物未注册"那一类），本腿这一格的射程没盖到，需裁决腿补问。
5. **winlive 这一族在 CI 有没有面**：本票的真机判据（借／否决／归还、第 3 步真执行）与票 245 同族，
   `ci.yml` 里没有 `-tags winlive` 这一档 ⇒ 回归只在这台机器上手跑时才有读数。⛔ 本腿不改 CI、不造门（票 245 表 R-3 同口径）。
6. **`ShutdownHookSet` 该不该允许一步多钩子**：本腿按"一个位一个所有者，第二个拒收"实现（静默覆盖＝丢一次 teardown）。
   将来某票想在同一步挂两件事（例如第 2 步既收球又停 KWS），需要**链式注册**或调用方自带复合钩子——
   那是接口形状的选择，本腿不自裁扩。
7. **本机没有第二个真账号、也没有真任务的实测面**：与票 238 那批残余同口径，本票未涉及，登记以免被当成漏做。

---

## §6 收尾名册对撞

（终态 `git status --porcelain` 与 `roster-start.txt` 的差集、逐枚 commit 号、写面是否清空、临时件落盘清单——
收尾现跑后逐字填，⛔ 不在跑之前预填读数。）
