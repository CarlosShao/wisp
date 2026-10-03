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

> 本节按任务书三问之问二逐行量：**谁持有热键表／rebind 是全拆重装还是增量改／借用那一跳走的
> 是不是同一条路径**，然后回答"要把这枚钉装上需要哪几枚零件、今天可达吗"。
> 对前一腿 `245-c1`（`.scratch/wisp/probes/245/c1/census.md`，起手 `19b63425`）的判语在 §2.6。

### 2.1 热键表有两层，两层的持有者不是同一枚对象

| 层 | 持有者 | 现量 |
|---|---|---|
| Win32 真表（桌面级，按 `id + 调用线程` 生效） | **ui-sta 那枚线程**；`RegisterHotKey` 全部用同一个 `b.hwnd` 作 host | `internal/ball/hotkey_windows.go:389-401`（registerer）/`:407-409`（unregisterer）；id 名册 `:29-34`（`hkSummon=1 / hkMute=2 / hkCancel=3 / hkPanel=4`），注册顺序表 `:38-46` |
| 球内记账（5 枚字段） | `*Ball` 结构体字段，**没有 mutex** | `internal/ball/ball_windows.go:122-126`：`registeredHotkeys / hotkeyReport / boundCfg / escTakenOver / cancelBinding`；构造期直接写（不经 uiRun）在 `:253-256` |

- 这 5 枚字段的**唯一并发保护机制**＝"只在 ui-sta 上碰"：读写的公开入口
  （`HotkeyReport` `internal/ball/ball_windows.go:830-834`、`ConfiguredHotkeys` `:840-844`、
  `RegisteredHotkeys` `:849-853`、`EscTakenOver` `:910-914`、`RebindHotkeys` `:813-825`、
  `TakeEscForCancel` `:873-889`、`ReleaseEscAfterSession` `:895-907`）**全部**套在
  `b.uiRun` 里，而 `uiRun`（`:741-752`）非本线程时是 post-and-wait（`:746-751`）。
  ⇒ 结论先摆明：**rebind 与 borrow 之间不会撕裂、不会 data race，会的是"谁先排进那条线程"。**
- ⚠ 一处注释与实现不符（影响落地腿的判断，具名记在 §4）：
  `internal/ball/hotkey_reload.go:78-80` 说 `Check()` 在球自己的 UI 线程上调用"would deadlock"；
  `uiRun` 的实现（`internal/ball/ball_windows.go:742-745`）在**同一线程时是 inline 执行**，
  不会 deadlock。真正会挂的是**另一条**：`staThread.PostTask`
  （`internal/ball/sta_windows.go:233-249`）在 `hwnd == 0` 时 `:243-246` 把任务删掉就 return
  （注释写"run inline as last resort"，代码里没有 inline），闭包不执行 ⇒ `uiRun` 的
  `<-done`（`internal/ball/ball_windows.go:751`）**永久阻塞**。
  这条才是"桥必须早于 `Ball.Close` 收口"的真理由（Close 会 `forgetWindow`：
  `internal/ball/ball_windows.go:964` → `internal/ball/sta_windows.go:169-173`）。

### 2.2 rebind ＝ 全拆重装，不是增量改（逐行）

`internal/ball/ball_windows.go:813-825`，整段在**一枚** `uiRun` 闭包里，故对 ui-sta 是原子的：

1. `:816 unregisterAll(b.hwnd)` → `internal/ball/hotkey_windows.go:503-511`：
   对 `hkNames` **全部四枚 id 逐个 `unreg`**，其中含 `hkCancel(3)`。
   `:498-502` 注释逐字承认这一刀的目的就是"连我们不跟自己记账的那枚 takeover id 一起拆掉"
   ——**takeover id 就是借用的裸 Esc**。
2. `:817 registerAll(b.hwnd, cfg)` → `:416-418` → `registerAllWith :453-496`；
   `:457-463` 的 `hkCancel` 分支**只分类、continue、绝不注册**（idle 集里没有 cancel）。
3. `:818-822` 换记账：`hotkeyReport`／`boundCfg`／`registeredHotkeys = rep.Live()`／
   `cancelBinding = cfg.Cancel`／**`escTakenOver = false`**。
   ⇒ `:822` 那一行就是"丢借用"的**记账面**；`:816` 是它的 **Win32 面**。两处同闭包，所以两者必然同时发生。
4. 增量改的可能性：`registerAllWith` 只吃 `HotkeyConfig` 全量四枚（`:455` 把它摊成
   `[]string{Summon, Mute, Cancel, Panel}` 按 `hkNames` 顺序对位），**没有**"只改动了的那一枚"
   这种入口 ⇒ 库层面 rebind 只有全拆重装一形。

### 2.3 借用那一跳走的是同一条路径吗：同槽同原语，不同入口，且**键名不受配置影响**

| 对照项 | rebind | borrow（借） | return（还） |
|---|---|---|---|
| 公开入口 | `ball_windows.go:813` | `:873` | `:895` |
| 用到 hkCancel(3) | `unreg(3)`（经 `unregisterAll`）**从不 reg(3)** | `reg(3, escBorrowAcc())` | `unreg(3)` |
| 底层原语 | `hotkey_windows.go:503-511 / 453-496` | `:531-539 takeEscWith`（`:532` 先 `unreg(3)` 清残、`:533` 再 reg） | `:553 releaseEscWith`（只 `unreg(3)`） |
| 注册的那把键 | **cfg 里那枚（cancel 除外）** | **永远裸 Esc**：`:524 escBorrowAcc() = {modNoRepeat, vkEscape}`，`:519` 注释逐字 "the bare Esc, whatever the configured binding says" | — |
| 改 `escTakenOver` | `:822` 无条件 false | `:885` true（`:875-877` 幂等早退） | `:901` false（`:897-899` 幂等早退） |
| 报告行来源 | 整份新 `HotkeyReport`（cancel＝`cancelIdleLine`） | `:886 withCancel(cancelBorrowedLine())`；被拒时 `:881 cancelFailedLine` | `:904 withCancel(cancelIdleLine(b.cancelBinding))` |

⇒ **判语（任务书问"是不是同一条路径"）**：三者操作的是**同一枚 Win32 槽位与同一组 5 枚记账字段**，
所以任何两跳的先后**会互相覆盖**；但它们**不是同一条代码路径**（rebind 走 `registerAllWith` 的
skip 分支＋盲拆，借还走 `takeEscWith/releaseEscWith` 这对专用原语）。
"不是同一条路径"正是丢借用的**机制**：rebind 没有任何一行读过 `escTakenOver` 再决定要不要保留它。

★ **本节最要紧的一条外延（245-c1 未说，且它直接改写 AC#1 的射程）**：
`[hotkey] cancel` 那枚值在**产线任何路径上都不会被注册成键**——
idle 遍跳过它（`hotkey_windows.go:457-463`），借用那一跳硬编码裸 Esc（`:524/:533`），
归还再显式拒绝重绑配置值（`:547-553` 注释逐字："re-binding the configured default is the
defect this ticket is about"）。⇒ 形 A 落地后，用户把 `cancel` 改成 `Ctrl+Alt+X` 的**全部效果**
＝报告里那行 standby 的拼写变了；idle 时 Ctrl+Alt+X 按了没事、Confirming 时裸 Esc 照样被借走。
这一格不是"接线漏了一枚"，是**hot 档承诺在第四枚上今天就没有兑现面**，
票 258 AC#1 那句"四枚热键以 `config.toml` 为准"按现读码只对三枚成立 ⇒ **必须上报，⛔ 我不自填修法**。

### 2.4 交错矩阵：rebind 撞上借用在飞会怎样（四种时序，逐行给依据）

前提（产线形 A 落地后）：借＝`cmd/wisp/resident_approval_windows.go:441`（`Prompt` 内，仅
`p.Level == "L1"`，`:440`）；还＝`settleOrb` `:491`；rebind＝桥的 `Check()`
（`internal/ball/hotkey_reload.go:96`）→ `RebindHotkeys`。窗口＝3s
（`internal/agent/approval/queue.go:116 DefaultL1Window`，界 `:120/:122` = 2s/3s）。

| # | 时序 | 结果（依据） |
|---|---|---|
| 1 | **借已完成、窗口还开着，rebind 后到**（任务书直问的那一格） | `unreg(3)`＋`escTakenOver=false`（`ball_windows.go:816/822`）⇒ 窗口剩余时间里**没有 cancel 键**；报告回 standby 行（`:817-820`）⇒ `Problems()` 为空（`hotkey_windows.go:365` 具名把 Standby 排除在问题行之外）⇒ **症状为零的降级**。之后 `settleOrb` 的 `ReleaseEscAfterSession` 走 `:897-899` 幂等早退，**不会二次 unreg、也不会把报告写烂**（这条不变量值得钉）。 |
| 2 | rebind 的闭包与借的闭包**同时在排队** | 二者都经 `uiRun` post 到同一条 ui-sta 队列 ⇒ **串行，不分先后顺序地二选一**。rebind 先排 ⇒ `Prompt` 里 `:441` 之后 `:454 !EscTakenOver()` 读出 true（`:885` 已置位）⇒ 有键；borrow 先排 ⇒ rebind 把它抹掉 ⇒ 无键。**同一发配置改动可以两种结局**，胜负手＝两次 post 的落地顺序，产线无人可预测、也无人记录。 |
| 3 | 借在飞、rebind 之后**又来一次借**（第二张 L1 卡） | `:875` 的幂等位已被 `:822` 清成 false ⇒ 第二张卡**能正常借到**（`:532` 先清残），这一形反而健康。 |
| 4 | 桥的 `Refresh`（`:82-86`）与 `Check` 撞上球的收口 | 见 §2.1 末条：`PostTask` 在 `hwnd==0` 时**丢任务不执行** ⇒ `uiRun` 永久阻塞 ⇒ 挂住的正是退出序列（`cmd/wisp/resident_ball_windows.go:237-244`）。**这跟 rebind×borrow 是同一枚雷的两半**，落地腿的收口次序必须一起处理。 |

**★ #1 的真实代价比注释里说的重一档**（现读码，245-c1 未写）：cancel 键没了之后，那张 L1 卡
**不会永远挂着**——它到点**执行**。`internal/agent/approval/gate.go:319 case <-deadline:` →
`:332` 记 `ANSWER-EXPIRED decision=timeout->execute` → `:338 return tools.AnswerTimeout`；
同文件 `:243` 注释逐字 "the window ending unopposed means EXECUTE (AnswerTimeout)"，
`gate.go:320-323` 再补一句"Timeout MEANS EXECUTE on the L1 route (SPEC-06 §2)"。
⇒ 一次撞在借用在飞时的 rebind，等价于**把一张本可否决的确认卡变成必然放行**，而且用户
看不出区别（报告干净、控制台无话）。这条应写进票 258 落地腿的判据射程；
⛔ 我不改票面、不自填修法。

### 2.5 全仓现役钉的射程（复量"零钉"）

尺＝`grep -rn "RebindHotkeys" cmd internal tools --include=*.go`（全文读数见 §1.2）。
逐枚核过、确认**没有一枚把 rebind 与在飞借用并置**：

- `internal/ball/hotkey_live_test.go` `TestLiveHotkeyOccupiedVsNotAttempted`：rebind 在
  `:247 / :272 / :297`，借还在 `:310-313`（`requireIdleRoster → TakeEscForCancel →
  requireEscBorrowed → ReleaseEscAfterSession → requireEscReturned`）——**严格串行，
  最后一次 rebind（:297）之后才借**。逐行读过 `:294-313`，中间没有任何 borrow 挂起中的 rebind。
- `internal/ball/hotkey_status_test.go`：`TestHotkeyReloaderRebindsOnConfigChange` 用假 binder
  （`:444 recorder.RebindHotkeys`）⇒ 碰不到 `escTakenOver`；`TestCancelBorrowRoundTrip`
  （`:202` 起）走的是 `registerAllWith/takeEscWith/releaseEscWith` 原语，**不经 RebindHotkeys**。
- `cmd/wisp/resident_approval_live_246_windows_test.go` 三枚：`:72` 的借还在无 rebind 的窗口里；
  `:195` 是 L2（不借，`:240`）；`:289` 是退出弃窗。⇒ **零钉复认成立**。

### 2.6 ★ 对 `245-c1` 那句结论的判语：复认其三、外延其二、修正其一

`245-c1`（`.scratch/wisp/probes/245/c1/census.md` §3.3／§4.1）的结论拆开是三句：

1. **"借还那一对已有最强形钉，复用不新增"——复认。**
   `cmd/wisp/resident_approval_live_246_windows_test.go:72 TestLive246ConfirmingCardBorrowsEscVetoesAndReturns`
   仍在，形状未变：正控 `-steal`（`:76-81`）＋无 Wisp 基线（`:83-86`）＋idle（`:108`）＋
   借（`:136 EscTakenOver` / `:139 len(rep.Live())!=4`）＋桌面被借（`:145-147 keydown_esc==0`）＋
   否决（`:157 AnswerVeto`）＋还（`:173`、`:176` idle 名册、`:177-179 keydown==1`）＋
   无残留（`:183`）。**全仓唯一带第二进程桌面读数的借还钉**，新腿照抄不得、复用即可。
2. **"rebind×borrow 零钉"——复认**（§2.5 用今天的 HEAD 重新逐枚数过，一枚不多一枚不少）。
3. **"今天产线不可达"——复认**，而且**更强**：不止"没接重载器"，是 `cmd/wisp` 连
   `RebindHotkeys` 的边都为零（§1.2 那把尺）；两条既有驱动路径都断在装配层
   （`config_reload.go` 的 tick 只到 `rt.mgr.CheckAndReload()`，见 §3）。
4. **外延 A（245-c1 未量到）＝§2.3 那条 ★**：`[hotkey] cancel` 在产线任何路径上都不成键，
   所以"rebind 会不会影响借用那把键"这个问题的答案今天**恒为"不会"**——借用那把键
   从来不看配置。这条**加重**而非推翻它的结论。
5. **外延 B＝§2.4 #1 的"到点执行"**：丢借用的后果不是"一张卡少个键"，是"一张卡必然放行"。
   这条同样**加重**。
6. **修正（只修可装性，不修结论）＝"要等形 A 落地才谈得上钉"不成立**：`245-c1` §1.3 末句
   说"A539 形 A 落地后常驻腿才长出这个面，残余才从注释里的话变成可发生的事"——
   对**产线**成立；对**钉**不成立，见 §2.7。

### 2.7 这枚自钉要长在哪：零件表与今天可达性（任务书点名要答的那半）

**需要的五枚零件，逐枚验货（全部现成、全部在票 258 落地腿的写面 `cmd/wisp` 内）**：

| # | 零件 | 现量（file:line） | 今天可用？ |
|---|---|---|---|
| 1 | 真球＋真 `residentBall` 句柄 | `cmd/wisp/resident_ball_windows.go:158 startResidentBall`；既有非 winlive 用法判例 `cmd/wisp/resident_approval_246_windows_test.go:314` | 是（但见下方"tag 陷阱"） |
| 2 | 真 L1 卡把借举起来 | `cmd/wisp/resident_approval_windows.go:139 bindBallHost` / `:261 AskOnTaskRoot`；`Prompt` 内 `:440-441` 借、`:454` 读回 | 是（246 族全套跑通） |
| 3 | 一次 rebind，且**要用形 A 那台机器本身**而不是裸调 `RebindHotkeys` | 导出的 `internal/ball/hotkey_reload.go:73 NewHotkeyReloader`（binder 是导出接口 `:44-46`，`*Ball` 已实现）；样板闭包 `cmd/balldebug/main.go:237-243` | **是——不需要等形 A**：这枚桥是导出 API，测试可在 `cmd/wisp` 内自己 new 一枚挂到 `rb.b` 上，"改 config.toml ⇒ Check() ⇒ rebind"整条链今天在测试进程里就能闭合 |
| 4 | 桌面侧读数（判"真注册／真注销"，不信报告行） | `cmd/wisp/resident_approval_live_246_windows_test.go:405 buildEscListener246` / `:444 observeEsc246`（`-steal`/`-watch` 两形，字段读数 `:502-503`） | 是（同包可复用） |
| 5 | 收尾不变量的既有尺 | `:350 requireIdleCancelSlot246`（内部 `:352 HotkeyReport`、`:357 EscTakenOver`）、`:534 waitFor246`、`:365 captureAudit246` | 是 |

**钉该断言什么（三头，形状而已，⛔ 不写码）**：
N1 借在飞时 rebind ⇒ `EscTakenOver()==false` 且 `RegisteredHotkeys()` 里**没有 id 3**
（⛔ 别只读 `HotkeyReport()`：`245-c1` §4.3 已警告"报告烂了骗得过去"，本腿复认——
`requireIdleCancelSlot246:352` 读的正是报告；纯 Win32 名册面在 `internal/ball/ball_windows.go:849`）；
N2 **窗口仍开着**（`AwaitingHuman()` 仍 true，`:125` 同形）——这一句才是这枚钉区别于 246 第一例的地方；
N3 到点之后**执行了**（`AnswerTimeout` 那一路，audit 行 `ANSWER-EXPIRED ... timeout->execute`，
尺＝`captureAudit246`）＋ `Problems()` 仍为空（把"症状为零"这件事本身钉住，
`internal/ball/hotkey_windows.go:365` 是它之所以无声的原因）。
红绿两向：残余未修时 N1/N3 就是**具名红灯**；将来票 245 那笔在册残余真去修（rebind 保留借用或重借）时，
同一枚钉翻期望值即成绿钉——**一枚钉承担两档**，不需要新增第二枚。

**"能不能在今天可达的路径上装"——分两层答，⛔ 混说会误导落地腿**：

- **装得上**：五枚零件全现成，且**不依赖形 A 先落地**（零件 3 就是那一处修正）。
- **装在哪一层有讲究（tag 陷阱）**：
  - 装进 `//go:build windows && winlive`（`cmd/wisp/resident_approval_live_246_windows_test.go`，`:1`）
    ⇒ 基建全套复用、零新造，**但 CI 不跑它**：windows job 那枚 cmd/wisp 步骤走的是
    `scripts/wisp-cli-tests.sh`（`.github/workflows/ci.yml:455-475`），该脚本不含 winlive tag
    （尺＝`grep -n "tags" scripts/wisp-cli-tests.sh .github/workflows/ci.yml` 均无 winlive 命中）。
    ⇒ 这是一枚"要人跑才会红"的钉，**属自钉、不属门禁**。
  - 装进 `//go:build windows`（会进 CI 分母）⇒ 必须先解决"桌面不在场时它算什么"：
    既有那枚同 tag 的球测试**故意容忍无桌面**
    （`cmd/wisp/resident_approval_246_windows_test.go:309-310` 逐字 "On a machine with no
    desktop both legs collapse to the second reading"，断言的是 flag 而不是窗口存在），
    而这枚新钉的三个断言**全部**以"真借到键"为前提 ⇒ 无桌面时它只能 `t.Fatal`（＝CI 永久红）
    或 `t.Skip`（＝票 258 禁区明文禁止，且 `tools/d22scan/runtests.sh` 把 SKIP 判红）。
    ⇒ **本腿判：这枚钉今天该装在 winlive 层，不装在 CI 层**；若编排者要它进门禁，
    那是"CI windows runner 有没有可交互桌面"这台机器的读数，⛔ 我量不到（见 §5）。
- **要不要新造第六枚零件**：不要。唯一"新"的东西是零件 3 那枚**测试内自搭的桥**，
  它用的是导出 API，⛔ 不需要给 `internal/ball` 加把手（那一层是本票只读面）。

## §3 今天可达路径名册

> 判据拆成三跳，逐跳量"今天有没有人在跑"：**跳 1 文件→内存**（hot 档应用）、
> **跳 2 内存→球**（把新值交给宿主）、**跳 3 球→Win32**（真注册）。
> ⚠ 本节所有 `cmd/wisp` 行号是**工作树现读**（09:3x）：`config_reload.go` 此刻是 dirty
> （`git status --porcelain` 读数 `M cmd/wisp/config_reload.go`＋`?? cmd/wisp/config_readers_255.go`，
> 一枚 255 的写在飞的腿正在这个包里补回执），**半前一节的函数起始行未漂**
> （复量：`startConfigReload` 仍在 `:105`、Spawn watchdog 仍在 `:119`、`reloadOnce` 仍在 `:152`、
> `rt.mgr.CheckAndReload()` 仍在 `:153`——与票 258 现量 4 逐字同），漂的是 `reportReload` 的函数体。

### 3.1 跳 1（文件 → 内存）：引擎活着，但**跑不跑取决于这个常驻进程是从哪拉起来的**

- 应用表本体：`internal/config/manager.go:281`
  `{"hotkey", &cur.Hotkey, &fresh.Hotkey, func() { cur.Hotkey = fresh.Hotkey }}`
  ⇒ 值会跟着文件换（票面现量 1 复认，**但行号从票面的 `:278` 漂到今天的 `:281`**，
  漂因＝票 255 AC#2-ⓑ 把 hot 名册改成从 `TierRegistry` 派生并加了 `:293-300` 的同源 panic）。
- 唯一的驱动原语：`internal/config/manager.go:142 CheckAndReload`（mtime+size 指纹早退在 `:156`）。
  产码调用点**全仓只有两枚**（尺＝`grep -rn "CheckAndReload()" cmd internal tools --include=*.go`
  去掉 `_test.go`）：
  1. `cmd/wisp/config_reload.go:153`（watchdog tick，**下面第 3 层**），
  2. `cmd/balldebug/main.go:243`（桥的 `Refresh`，旁支调试进程）。
- ⚠ **`OnReload` 这条路对 `[hotkey]` 结构性不通**（不是"没接"，是接不上）：
  `internal/config/manager.go:198` 的触发条件是 `len(rep.Reload) > 0`，字段注释 `:52-54`
  具名"reload-tier sections"；`[hotkey]` 是 hot 档、永不进 `rep.Reload` ⇒
  谁把 `bridge.OnReload()` 挂上去都不会为热键改动响。这一条与
  `internal/ball/hotkey_reload.go:26-32` 的自我说明一致，本腿独立复核成立。
- **tick 的武装条件**（票面现量 4 说"重载的引擎在常驻腿是活的"——**这句要加限定**）：
  武装点只有一枚 `cmd/wisp/run.go:813 rt.startConfigReload()`，它在 `assembleRuntime`
  （`cmd/wisp/run.go:382`）**函数尾部**，前面任何一条 `return rt, 2` 都＝不武装。而
  `assembleRuntime` 在常驻腿只被 `cmd/wisp/resident_task_source_windows.go:265` 调一次，
  调用点在 `:221-231` 那道门**之后**：
  - 那道门的谓词＝`interactiveStdin()`（`cmd/wisp/approval_reply_stdin_windows.go:41-52`：
    标准输入不是控制台输入缓冲就返回 nil，`:47-49` 具名"a pipe, a file, or no console"）；
  - 常驻那条腿自己的场景自述正是这件事：`cmd/wisp/resident_windows.go:52-57`
    逐字 "the resident process is the leg owner actually uses - double click the icon,
    no terminal attached, stderr going nowhere"，`main.go:60-67` 的无参数分支先
    `attachParentConsole()` 再 `runResident()` ⇒
    **真·双击起来的常驻进程：既没有 `rt.mgr`、也没有 tick、连"内存里的 `[hotkey]` 会跟着文件变"这一句都不成立**；
  - 反过来，**从终端里裸跑 `wisp`（无参数）**这一形有控制台 ⇒ 过门 ⇒
    `assembleRuntime` ⇒ `run.go:409` 建 mgr ⇒ `run.go:813` 武装 tick ⇒ 跳 1 活。
    ⛔ 这一条我只按读码给形状，**"我的常驻腿今天到底是被谁拉起来的"是实机读数，见 §5**。
- 常驻进程里另一枚 mgr（面板链）从不 tick：`cmd/wisp/panel_resident_windows.go:162-164`
  逐字自陈，构造点在 `cmd/wisp/panel_resident_windows.go:166-171` → `cmd/wisp/panel_inbound.go:228-233`。
- **谁来写 `[hotkey]`**：全仓设置写面只有 provider/model/role 五枝
  （`cmd/wisp/panel_config_store.go:173-217` 的 `ApplySetting` 分支表，逐条是
  `SetProviderBaseURL / SetProviderAPIKeyRef / SetModelContextWindow / SetModelPriceIn /
  SetModelPriceOut / SetRoleChatModel`）⇒ **产码里没有任何一条路径写 `[hotkey]`**，
  改动只能来自**人手编辑 config.toml**（票面说的就是这一形）。

### 3.2 跳 2（内存 → 球）：**今天整条为零**

- 球侧唯一的"交给新值"入口＝`internal/ball/ball_windows.go:813 RebindHotkeys`；
  产码调用者只有桥（`internal/ball/hotkey_reload.go:96`）⇒ 于是问题归到"谁 new 了桥"：
  **`cmd/balldebug/main.go:237` 一枚**（§1.2 那把尺的全文读数）。
- 桥自己还挂在 `cmd/balldebug/main.go:230` 的 `-config` 分支里 ⇒ **不带 `-config` 连它也不活**。
- `cmd/wisp` 侧对 `[hotkey]` 的读取次数：**零**（尺＝
  `grep -rn "\.Hotkey\b" cmd internal tools --include=*.go` 去测试 ⇒ 命中只有
  `cmd/balldebug/main.go:238`、`cmd/wisp/config_readers_255.go:110/:113`（这两处是**引用句**、
  不是读值）、`internal/config/manager.go:281`（hot 应用表）、以及两处注释）。
  ⇒ **258-a1 §7 第 1 点补的"`wisp run` 不建球"复认且加码**：
  `cmd/wisp/run.go` 里 `ball.` 出现次数＝**0**（尺＝`grep -c "ball\." cmd/wisp/run.go`），
  全仓 `cmd` 下 import `wisp/internal/ball` 的产码文件只有
  `cmd/balldebug/main.go`、`cmd/wisp/resident_approval_windows.go`、`cmd/wisp/resident_ball_windows.go` 三枚。

### 3.3 跳 3（球 → Win32）：今天只有两条，且都与 `[hotkey]` 无关

- 启动那一遍：`internal/ball/ball_windows.go:255`（`registerAll(b.hwnd, b.opts.Hotkeys)`），
  而 `b.opts.Hotkeys` 在常驻腿是 `cmd/wisp/resident_ball_windows.go:171` 那枚写死值；
  ⚠ 另有 `internal/ball/ball_windows.go:153-154` 的零值兜底（§1.3 已具名）。
- 借用那一跳：`internal/ball/hotkey_windows.go:531-539`，**恒为裸 Esc**（`:524`）⇒
  它注册的是"球自己决定的那把键"，不是配置里那把（§2.3 ★）。
- 归档旁证（球"从来没吃过配置"的机制面，258-a1 §7 第 1 点的②复认）：
  常驻进程里两枚 mgr（`run.go:409` ＋ `panel_inbound.go:230`）
  都不在 `startResidentBall` 的参数表（`cmd/wisp/resident_ball_windows.go:158`）上。

### 3.4 名册（任务书点名"哪条路今天真能触发、哪条不能"）

| # | 路径 | 三跳状态 | 现量 | 今天能不能让改过的 `[hotkey]` 变成手上真按得动的键 |
|---|---|---|---|---|
| 1 | **`wisp`（无参数，终端里起）＝常驻腿** | 跳 1 活 / **跳 2 零** / 跳 3 只吃写死值 | `resident_windows.go:163→206`、`run.go:813`、§3.2 | **不能**。改完文件内存会换，**没人次日或当场把它交给球**；重启进程才会（下一次 `ball.New` 仍吃 `DefaultHotkeys()` ⇒ **连重启都不解决**，这一格比票面更狠） |
| 2 | **`wisp`（真·双击/GUI 无控制台）** | **跳 1 也不活** / 跳 2 零 / 跳 3 写死 | `resident_task_source_windows.go:221-231`＋`approval_reply_stdin_windows.go:41-52` | **不能**，且比 #1 更空：这个进程里连 mgr 与 tick 都没有（`[hotkey]` 的**内存值**都不换） |
| 3 | **`wisp run <文本>`** | 跳 1 活 / 跳 2 零 / **跳 3 不存在（不建球）** | `main.go:89-91`→`run.go:181→249→813`；`grep -c "ball\." cmd/wisp/run.go`＝0 | **不能**（258-a1 的"对 `wisp run` 是空集"复认）。但它会**把这句老实话打给你看**，见 §3.5 |
| 4 | **`wisp panel-inbound`** | 跳 1 不武装 / 无球 | `config_reload.go:91-99`（`hotReloadDisabledPanelInbound` 自陈） | **不能**，且自己会说明不能 |
| 5 | **常驻面板链** | 有自己的 mgr、从不 tick | `panel_resident_windows.go:162-164` | **不能** |
| 6 | **`cmd/balldebug -config <path>`** | 三跳全活 | `main.go:230-255` | **能——全仓唯一一条今天真能触发 rebind 的路**。⚠ 它是旁支调试进程：球在它自己的进程里、`observe.Default`、无审批门 ⇒ 它**证明机器能跑，不证明出厂腿接上了** |
| 7 | **手改 config.toml 之后不重启** | 见 #1/#2 | — | 与票面题目同义：**热改了没人接** |

### 3.5 今天这条缺口**已经有一枚回执在替它说实话**（票面未列，本腿新量）

`cmd/wisp` 里刚落地的票 255 AC#1 回执名册**已经为 `[hotkey]` 写了一行判定**：
`cmd/wisp/config_readers_255.go:113` ＝ `hotClaimDebugHostOnly` ＋
`"cmd/balldebug/main.go:238 [h := mgr.Config().Hotkey] - cmd/balldebug is not the shipped host, and wisp run binds ball.DefaultHotkeys()"`。
⇒ 效果：手改 `[hotkey]` 后不再说"已立即生效"，改说"值已换进本进程内存，
但本宿主没有会按新值做事的读者"（那句在 `cmd/wisp/config_reload.go:189-194` 的 quiet 分支，
工作树现读）。**缺口对用户不再静默**，且这行话**常驻腿也打得出来**：
名册 #1（终端里裸跑 `wisp`）走的是同一枚 `agentRuntime`，
`cmd/wisp/resident_task_source_windows.go:252-264` 把 `stdout: os.Stdout` 交给了装配根，
所以 quiet 那句会打到操作者的控制台上；只有名册 #2（真·双击、无控制台）
是"打了但没人看"（`cmd/wisp/resident_windows.go:52-57` 自陈 stderr going nowhere）。
⛔ 两形共同点：**这句话改不了任何一枚键的注册**——票 258 要修的仍是那半。

⚠ 顺带一处**该句自身的现量错误**（新落地的文件，趁早报）：`:113` 写的是
"`wisp run` binds `ball.DefaultHotkeys()`"，而 §3.4 #3 量到 `wisp run` **零 `ball.` 引用、不建球**；
真正吃 `DefaultHotkeys()` 的是**常驻腿** `cmd/wisp/resident_ball_windows.go:171`——
同一枚文件的 `:112` 其实引对了行（它引 `resident_ball_windows.go:171`），
**句子里的主语却写成了 `wisp run`**。这是 `cmd/wisp` 在飞文件里的一处措辞缺陷，⛔ 不改、只具名。

### 3.6 ★ 形 A 落地**必配**的两处 255 侧改写（票面 §15-20 与 §48 都没列这一族）

`cmd/wisp/config_readers_255.go` 的两枚字段是被**测试**当"可核断言"用的
（该文件 `:78-81` 与 `:142-147` 自陈两枚尺名：
`TestTicket255RosterEvidenceLinesStillSayWhatTheyClaim`、
`TestTicket255RosterStillMatchesTheActualReadSites`；⚠ **这两枚测试此刻还没落盘**——
`grep -rn "TestTicket255" cmd tools --include=*.go` 只命中这三处**注释**，
说明 255 的测试腿与回执腿都在飞）：

1. `cmd/wisp/config_readers_255.go:112-113` 引
   `cmd/wisp/resident_ball_windows.go:171 [Hotkeys:  ball.DefaultHotkeys(),]` 作**证据行**
   ⇒ 形 A 把那一行改掉的那一刻，"证据行仍说着它说的话"这枚尺**必红**（若它要求行号也命中，
   连 `:171` 这个数都要跟着重算）。
2. `cmd/wisp/config_readers_255.go:162` 的 `sectionReadSites["Hotkey"] = {"cmd/balldebug/main.go"}`
   ⇒ 常驻腿一旦真的开始读 `[hotkey]`，读值文件集合就**多出一枚**，"名册与实际读点对齐"这枚尺**必红**。

⇒ **给落地腿的最小改动面补两枚**（都在写面 `cmd/wisp` 内、都是这张表的行）：
`:110-113` 的 hotkey 判定要从 `debug-host-only` 换成 `consumed:` 那一档
（并照 `:92` 那行的规矩带上"读值处 file:line ＋ `[token]`"），`:162` 的读点名册要加上新落点。
⚠ 这是**回执与实现必须同步**的形状，不是可选清洁工；漏改＝要么 255 的尺红、
要么 `[hotkey]` 从此被回执说成"已立即生效"而**读者仍然不存在**（那正是票 255 AC#1  forbids 的
"摘尺当修尺"的反向版本）。
⚠ 另：**这两枚尺今天还不存在**⇒ 落地腿开工时必须先复量 255 交没交件（§5 具名）。

## §4 我可能写错的条目（自我对抗）

（取数中）

## §5 量不到的地方（具名）

（取数中）

## §6 交件判语

（取数中）
