# 258-a2 只读普查 — 热键接线最小改动面 / rebind×borrow / 今天可达性

> 腿：`258-a2`（只读普查）。工作树共享，同机有三枚写腿在飞 ⇒ **本腿零 Go 命令**（只读源码＋grep＋find＋git）。
> 票面：`.scratch/wisp/issues/258-the-schema-calls-hotkey-hot-tier-but-the-resident-leg-builds-the-ball-from-default-hotkeys-and-the-only-reloader-caller-is-balldebug.md`
> 票 258 AC#0 已由编排者裁**形 A**（§7 第 4 点，账 `A538`）；撤销口令「258 改形 B」。本腿**不选形**。

## §0 起手锚

- 起手时刻：`2026-10-03T09:11:03+08:00`；起手 HEAD `8a3790f0`（分支 `dev`）。
- `git status --porcelain -- cmd internal tools docs .scratch | wc -l` ＝ **353**（同机三枚写腿在飞，共享工作树）。
- 对照锚点：票 258 立票锚 `9588138a`（10-02 15:0x）、AC#0 收档与选形锚 `3df8b82b`（10-02 16:4x）、
  前一普查腿 `245-c1` 起手 `19b63425`。⇒ **本件所有 file:line 以 `8a3790f0` 时点为准**，
  引用旧件行号前先复核（并发腿会推移行号）。
- 本腿边界（全程遵守）：⛔ 零 Go 命令（不跑 `go test`/`build`/`vet`/`run`/`gofumpt`/`d22scan`/`staticcheck`）；
  ⛔ 不改任何产码；⛔ 票面 AC 勾选框一枚不碰；`frontend/**`／`design/**` 两层禁令（未进入、不转述）；
  `grep`/`find` 显式根＝`cmd internal tools docs scripts .scratch`；只 commit 不 push；
  commit 带显式 pathspec 且 add 与 commit 同发一条命令。
- 本腿与 `245-c1` 的分工：`245-c1` 量的是「借/还那一对的收尾仪器」，顺带给出 §1.3
  rebind×borrow 的现量与「今天产线不可达」判语；本腿**独立复量**它的三个支点
  （热键表持有者／rebind 是全拆还是增量／借用那一跳是否走同一条路径），判语在 §2。

## §1 接线最小改动面（逐处 file:line，现在是什么／要变成什么）

> 本节**不选形**（票 258 §7 第 4 点已由编排者裁形 A，账 `A538`）。我只把"形 A 那台机器要落进常驻腿，
> 手上的螺丝孔有哪几个"逐枚量出来，含每一孔"现在是什么／要变成什么"，以及票 §3 的排程约束下
> 哪些孔会撞别人正在写的文件。

### 1.1 那台重载器的真实签名与全套把手

- **签名本体**：`internal/ball/hotkey_reload.go:73`
  `func NewHotkeyReloader(binder HotkeyBinder, current HotkeyConfig, src HotkeySource) *HotkeyReloader`
  —— 三参数：binder／"球当前持有的那一份"（作为 diff 基线）／取新值的闭包。
  `:70-72` 注释逐字：构造**本身不重注册**（`applied` 起点＝`current`，`exists: true`，`:74`）。
- 它依赖的三个类型/字段：
  - `internal/ball/hotkey_reload.go:44-46` `HotkeyBinder` 接口＝单方法
    `RebindHotkeys(cfg HotkeyConfig) HotkeyReport`；真身 `*Ball` 实现在
    `internal/ball/ball_windows.go:813`。
  - `internal/ball/hotkey_reload.go:50` `type HotkeySource func() HotkeyConfig`。
  - `internal/ball/hotkey_reload.go:61` `Refresh func() error` 字段（`:82-86` 每次 Check 先跑它，
    失败只 Warn 并保留旧键；nil＝宿主自己刷）。
- 驱动把手（三选一，全是现成的）：`:81 Check() (bool, HotkeyReport)`（拉一次）、
  `:126 Run(ctx, every)`（自带轮询环）、`:110 Sections([]string)`＋`:120 OnReload()`
  （推形；`:26-32` 注释具名解释**为什么这台机器不能只靠 `OnReload`**：`[hotkey]` 是 hot 档，
  而 `OnReload` 只对 reload 档响）。
- 记账把手：`:154 Applied()`／`:161 Report()`／`:169 Rebinds()`（`:167-168` 注释：`Rebinds` 是给
  "不许每 tick 洗一次 Win32" 用的测试缝）。
- ⚠ **文件头示例已过期**：`internal/ball/hotkey_reload.go:16` 写的是 `b.HotkeyConfig()`，
  仓里**没有这枚方法**；真实取"球被告知要持有的那一份"是
  `internal/ball/ball_windows.go:840 ConfiguredHotkeys()`（两个真调用点都用它，见 §1.2）。
  落地腿照抄头注释会编译不过——这一处属"可顺手改的注释"，⛔ 但它在 `internal/ball`，
  本票写面只到 `cmd/wisp`（票 258 §48 行末：`写面＝cmd/wisp＋internal/ball（只读后者）`），
  ⇒ **改不得，只可在 cmd/wisp 侧照正确形状写**。

### 1.2 它现在的调用者名册（逐枚）

尺＝`grep -rn "NewHotkeyReloader" cmd internal tools --include=*.go`（本腿 09:1x 自跑，读数全文）：

| # | 位置 | 性质 |
|---|---|---|
| 1 | `cmd/balldebug/main.go:237` | **全仓唯一非测试调用点** |
| 2 | `internal/ball/hotkey_live_test.go:159` | winlive 测试（真 `config.Manager`，见 `:132`） |
| 3 | `internal/ball/hotkey_status_test.go:458` | 确定性层（假 binder `recorder`，`:444`） |
| 4 | `internal/ball/hotkey_status_test.go:511` | 同上，`Sections` 那一支 |
| 5 | `internal/ball/hotkey_reload.go:16/:70/:73` | 注释与定义本体，非调用 |

⇒ `cmd/wisp/**` 里 `NewHotkeyReloader` **零枚**；再补一把尺：
`grep -rn "RebindHotkeys" cmd internal tools --include=*.go` ⇒ `cmd/wisp` 也是**零枚**
（产码侧只有 `internal/ball/hotkey_reload.go:96` 这一处驱动者，测试侧
`internal/ball/hotkey_live_test.go:247/272/297` 与假件 `hotkey_status_test.go:444`）。
**这条比票面现量 3 更硬**：不只"没接重载器"，是常驻腿连 `RebindHotkeys` 的边都没碰过。

balldebug 那一套全套把手（形 A 的现成样板，逐行）：
`cmd/balldebug/main.go:231`（`config.NewManager(*configPath, nil)`）→ `:237-242`（构造＋src 闭包：
`mgr.Config().Hotkey` 四枚字段搬进 `ball.HotkeyConfig` 再过 `ball.ApplyHotkeyDefaults`）→
`:243`（`Refresh` 挂 `CheckAndReload`）→ `:244`（`mgr.OnReload = bridge.OnReload()`）→
`:250-253`（`observe.NewRoot` ＋ `observe.Default.Spawn(goroutineHotkeyBridge, "balldebug", root, Run)`；
名字常量在 `cmd/balldebug/main.go:45`＝`"balldebug-hotkey-bridge"`）→ `:254`（cancel＋join 句柄）→
`:255`（打出"bridge polling <path>"这一句）。

### 1.3 构造期那一改（球的第一口键从哪来）

| 处 | 现在是什么 | 要变成什么 |
|---|---|---|
| `cmd/wisp/resident_ball_windows.go:171` | `Hotkeys:  ball.DefaultHotkeys(),` | 吃配置那一份：`ApplyHotkeyDefaults(...)`  over `mgr.Config().Hotkey` 四枚字段（映射形状照 `cmd/balldebug/main.go:238-241`；`internal/ball/hotkey_windows.go:87` 是真身） |
| `cmd/wisp/resident_ball_windows.go:158` | `func startResidentBall(reg *observe.Registry, onCancelEsc escVetoFunc, hooks ...ballHostHook) *residentBall` — 参数表里**没有任何配置把手** | 要么加一枚参数，要么走**同一文件已有的 variadic hook 通道**：`ballHostHook` 定义 `:63-68`、承载结构 `panelHostHooks` `:70-74`、应用循环 `:160-163`、"不动既有调用点"的设计意图逐字写在 `:76-81`（"so the three existing two-argument call sites of startResidentBall keep compiling untouched"） |

- ⚠ 那枚 hook 通道今天**只能塞进 `panelHostHooks`**，而 `panelHostHooks` 的字段是给"开面板"用的；
  热键是**构造入参**（`ball.Options`，`internal/ball/ball_windows.go:67`），必须在 `ball.New`
  （`cmd/wisp/resident_ball_windows.go:165`）**之前**拿到 ⇒ 复用 variadic 时，hook 必须
  在 `:160-163` 那个循环里、`:165` 之前生效。这条是**次序事实**，不是选形。
- `startResidentBall` 全部调用点（改了参数表就要一起看）：产码 1 枚
  `cmd/wisp/resident_windows.go:163`；测试 4 枚
  `cmd/wisp/resident_approval_246_windows_test.go:314`、
  `cmd/wisp/resident_approval_live_246_windows_test.go:92 / :202 / :296`。
  ⚠ 这 4 枚测试文件按票 §33 的串行约束属别人写面（`cmd/wisp` 测试面被 `253-r1` 占着，票 §48 括号）——
  **枚枚具名，落地腿开工前先确认没人在写**。
- **"配置缺失退回 DefaultHotkeys 并说得出那句话"落点**（AC#1 的后半）：
  `internal/ball/ball_windows.go:153-154` 已经有
  `if opts.Hotkeys == (HotkeyConfig{}) { opts.Hotkeys = DefaultHotkeys() }`——
  ⇒ 若 cmd/wisp 侧传了零值，**库内会静默补默认**，那句"退回 DefaultHotkeys"的话**库不会替你说**。
  要说得出来，只能落在 cmd/wisp：现成的句子挂载点在
  `cmd/wisp/resident_ball_windows.go:202-208`（`Problems()` 逐行打印＋`rb.verdict` 那发
  `"hotkeys live %d/4: %s"`），或 `hotkeySummary` `:316-326`（它已能把空串说成 `(unset)`）。

### 1.4 桥的装配位（常驻腿）——三个硬事实先摆出来

1. **对象不同层（票 §44 已具名的那笔代价）**：建球在 `cmd/wisp/resident_windows.go:163`，
   而常驻进程里**第一枚** `config.Manager` 直到 `cmd/wisp/resident_windows.go:206`
   （`startResidentTaskSource`）→ `cmd/wisp/resident_task_source_windows.go:265`
   （`assembleRuntime`）→ `cmd/wisp/run.go:409` 才存在。`:163` 与 `:206` 之间没有任何 mgr。
2. **那枚 mgr 还可能压根不存在**：`cmd/wisp/resident_task_source_windows.go:221-231`
   在没有可交互控制台且没有注入任务时**直接 `return nil`**（连 `assembleRuntime` 都不进）；
   即便进去了，`assembleRuntime` 在 `cmd/wisp/run.go:410-419`（`config.NewManager` 失败）
   与 `:437-446` 等多处 `return rt, 2`，而重载 tick 挂在函数尾部的
   `cmd/wisp/run.go:813`（`rt.startConfigReload()`）——**早退＝tick 根本没武装**。
   ⇒ 结论：`rb`（`:163`）与"会变的 mgr"（`:206` 之后，且可能没有）**不是同生共死的两枚对象**，
   桥的装配位必须处理"后到"与"永不到"两种形状。
3. **第二枚 mgr 不 tick**：常驻面板链自己新建一枚
   （`cmd/wisp/panel_resident_windows.go:166-171` → `cmd/wisp/panel_inbound.go:228-233`），
   并且它自己把这事写在脸上：`cmd/wisp/panel_resident_windows.go:162-164`
   "this leg does not tick config.toml either; the reload tick lives in `wisp run`"。
   ⇒ 258-a1 §7 第 1 点补的"常驻进程里有两枚 Manager"复认成立（`run.go:409` ＋ `panel_inbound.go:230`），
   且两枚都进不了 `startResidentBall` 的参数表也复认成立。

### 1.5 驱动位与收口位（形 A 特有的两处，逐行）

- **驱动位（谁每 tick 拉一次）**两种既有形状都在仓里，⛔ 我不裁：
  - (i) 桥自带环：`internal/ball/hotkey_reload.go:126 Run` ＋ `:144 checkSafe`（recover 在环内），
    挂法见 `cmd/balldebug/main.go:250-253`。要在常驻腿新起一枚 goroutine ⇒
    过 `internal/observe/goroutine.go:262 Registry.Spawn`；名字**不在**冻结六枚名册里
    （`internal/observe/goroutine.go:40 ResidentBaseline = 6`、`:43-45 ResidentNames`），
    后果只有 `:270-273` 那发 `slog.Warn` 与 `RosterReport.Unknown`（`:419`）——
    `ResidentOverBaseline`（`:422`）不会响，因为 Unknown 不计数为 Resident（`:418-420`）；
     boot 期那把尺 `internal/proc/boot_windows.go:113-118` 跑在任何 spawn 之前。
    **既有判例＋既有钉**：`cmd/wisp/panel_resident_windows.go:76-82`（明知 CategoryUnknown 仍不塞名册）
    与 `cmd/wisp/panel_resident_windows_test.go:917-926`（把它钉成"往名册加名字＝人工批准的契约变更"）。
  - (ii) 蹭既有 tick：`internal/ball/hotkey_reload.go:24-25` 逐字欢迎这一形
    （"on a host that already polls config every tick (the wisp watchdog), call r.Check() from that tick"）；
    常驻腿**确实有**那枚 tick：`cmd/wisp/config_reload.go:119`（`Spawn("watchdog", "config", ...)`）、
    环体 `:131-147`、`reloadOnce` `:152-162`、`rt.mgr.CheckAndReload()` 在 `:153`。
    ⚠ 这一形的挂点在 `agentRuntime` 上，而票 §44 把"要动 `agentRuntime` 结构"记在**形 B 的代价**里
    ⇒ 我把它当**两形的共用事实**报出来，不当新形：`rb.b` 要能被那一跳看见，
    就必须给 `agentRuntime`（定义在 `cmd/wisp/run.go`，字段表见 `rt.mgr` 赋值处 `cmd/wisp/run.go:425`）
    加一枚球侧把手，或反过来让桥自己带环（(i)）。
- **收口位（必须先于 Close）**：`cmd/wisp/resident_ball_windows.go:237-244 stop()` 里
  `rb.b.Close()`（`internal/ball/ball_windows.go:935`）。这是**真会挂人的地方**：
  - `Ball.uiRun`（`internal/ball/ball_windows.go:741-752`）非本线程时是 post-and-wait（`:751 <-done`）；
  - `staThread.PostTask`（`internal/ball/sta_windows.go:233-249`）在 `hwnd == 0` 时
    `:243-246` **把任务删掉并 return**——闭包不执行、`close(done)` 不发生 ⇒ `uiRun` 永久阻塞；
    （那几行的注释写"run inline as last resort"，**代码里没有 inline**，这是一处注释与实现不符，
    具名记在 §4）；
  - `hwnd` 归零有两条路：`Ball.Close` 自己 `b.sta.forgetWindow()`（`internal/ball/ball_windows.go:964`
    → `internal/ball/sta_windows.go:169-173`），线程收口 `releaseThread`（`:147-151`）。
  ⇒ 桥的 cancel/join 必须排在 `rb.stop()` 之前；现成的机制就是本装配根已经在用的 defer LIFO：
  `cmd/wisp/resident_windows.go:166`（`defer rb.stop()`）之后注册的那条**先跑**，
  `:168-173`（`defer ra.detachBall()`）是这一形状的**既有判例**，
  它防的正是同一件事（注释逐字："after Ball.Close posts WM_QUIT there is no STA thread left to answer"）。
  - 参照的"运行前重读把手、把手为 nil 就早退"写法：
    `cmd/wisp/resident_approval_windows.go:393-397 releaseBall` / `:405-409 currentBall` /
    `:429-434`（Prompt 取局部副本，nil 即 fail-closed）/ `:487-490`（settleOrb 同形）。
    桥要走的正是这条：拿 `*ball.Ball` 的局部副本＋nil 检查，⛔ 不许跨 Win32 post 持锁
    （`cmd/wisp/resident_approval_windows.go:380-382` 的字段注释逐字定了这条规矩）。

### 1.6 落地会撞到的既有尺（只报射程，⛔ 不改）

- `cmd/wisp/resident_ball_228_windows_test.go:37-49` 把那三句 boot 文案**抄成常量**
  （`:25-28` 注释自陈理由："reword one and a case here goes red"）⇒ **改 verdict 文案要连带改这里**。
- `cmd/wisp/resident_ball_live_228_windows_test.go:143-147`（只 `t.Logf` 抽 `"hotkeys live"`，**无断言**）、
  `:160-162`（钉 `0/4`）、`:163-167`（钉 `4/4`）⇒ 复认 258-a1 §3："三枚既有钉盖不到这一格"：
  一个"配置里 summon 写成了被占的组合、于是今天 3/4 变 2/4"的形状**三枚全不响**。
- `internal/ball/hotkey_status_test.go:171-196 TestDefaultHotkeysIdlePassHoldsNoEsc` 的**注释**
  `:172-174` 逐字写着 "DefaultHotkeys() is what cmd/wisp's resident leg hands ball.New" ⇒
  形 A 落地后这句成**过期陈述**（断言本身测的是 `registerAllWith(0, DefaultHotkeys(), ...)`，
  不会红；**红的是注释的真话**）。⚠ 该文件在 `internal/ball`＝本票只读面 ⇒
  **落地腿改不了它**，只能把这一处当"过期注释"具名上报，别当沉默。

## §2 rebind×borrow 现量与那枚自钉要长在哪

（取数中）

## §3 今天可达路径名册

（取数中）

## §4 我可能写错的条目（自我对抗）

（取数中）

## §5 量不到的地方（具名）

（取数中）

## §6 交件判语

（取数中）
