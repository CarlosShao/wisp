# 258-a1 — 只读普查：`[hotkey]` hot 档 vs 常驻腿建球的热键来源

票面＝`.scratch/wisp/issues/258-...-the-only-reloader-caller-is-balldebug.md`（AC#0 三问）。
本腿⛔不改产码、⛔不选形、⛔不判 AC#1 该不该做。唯一写点＝本文件。
重派接手：上枚 15:17 只落骨架，本枚 16:03 起全部读数自取（无前人读数可继承）。

## 0 起手锚（date/HEAD/porcelain 同发取）

- 接手时：16:03:43，HEAD `315c6f5f`；porcelain 有大量**他人**未提交改动（`.gitignore`、`design/**` 删除、`.scratch/wisp/probes/152|161/**` 等）——本枚只碰自己的 census.md。
- 本枚一发 commit `fa9a1fc6`（骨架改状）；取数完成锚 16:35:07，HEAD `fa9a1fc6`。全程零 Go 命令（同机 198-v2 在测）；工具只用 grep/sed/git/ls/wc/Read。

## 1 构造期：两条入口各自的热键从哪来

**起点三处复认（逐字）**：
- `cmd/wisp/resident_ball_windows.go:171`：`Hotkeys:  ball.DefaultHotkeys(),`（`ball.Options` 字面量内，`ball.New` 在 :165）。
- `internal/config/schema.go:176`：`// HotkeySection is [hotkey]; hot-tier (hotkeys re-register on change).`
- `internal/config/manager.go:278`：hot 应用表内 `{"hotkey", &cur.Hotkey, &fresh.Hotkey, func() { cur.Hotkey = fresh.Hotkey }}`。

**Q1a 两形是否一致（枚数/字段名/缺省值）**：
- `ball.DefaultHotkeys()`＝`internal/ball/hotkey_windows.go:68-70`：`{Summon: "Ctrl+Alt+Q", Mute: "Ctrl+Alt+M", Cancel: "Esc", Panel: "Ctrl+Alt+P"}`（summon 默认出处＝owner ruling R10，:63-67 注释）。
- `HotkeySection`＝`internal/config/schema.go:179-184`：四字段 `Summon/Mute/Cancel/Panel`（toml `summon/mute/cancel/panel`）。
- 枚数与字段名**一致**（4/4、同名同序）；**缺省值不一致**：schema 里只有 `cancel` 带 `default:"Esc"` tag（schema.go:182），summon/mute/panel **无 default tag** ⇒ LoadFile 产物里这三枚是 `""`（`internal/config/defaults.go:73-87` applyDefaults 只填带 tag 的叶）。把 `config.Hotkey` 直接灌球会得到"无 summon/mute/panel"——`ApplyHotkeyDefaults`（`internal/ball/hotkey_windows.go:87-102`）就是为这道缝准备的补缺省入口，注释（:72-80）逐字 "summon/mute/panel default to empty in the schema"。
- **`ApplyHotkeyDefaults` 非测试调用者枚数＝1**：`cmd/balldebug/main.go:239`（尺：`grep -rn ApplyHotkeyDefaults cmd internal tools --include=*.go` 排 `_test.go`；其余命中全是 `internal/ball` 的注释/定义）。

**Q1b internal/ball 有没有"从配置取"的现成入口**：
- **没有**。`internal/ball` 产码零 config 依赖（`config` import 只在 `hotkey_live_test.go:23`）。这是设计：`hotkey_reload.go:11-12` 逐字 "the bridge is deliberately free of any config import - internal/ball owns no business/config knowledge (SPEC-01 §3)"——host 供两个闭包，ball 不认识 config。
- `ConfiguredHotkeys()`（`internal/ball/ball_windows.go:840-844`）**不是**"从配置取"：它返回 `boundCfg`＝"球最后一次被 TOLD 的集合"（注释 :836-839 逐字 "the input of the last registration pass, not its outcome"）。**非测试调用者枚数＝1**：`cmd/balldebug/main.go:237`。
- 球侧唯一的零值兜底：`ball.New` 里 `if opts.Hotkeys == (HotkeyConfig{}) { opts.Hotkeys = DefaultHotkeys() }`（`ball_windows.go:153-155`；`Options.Hotkeys` 注释 :67 "zero value = DefaultHotkeys()"）。

**Q1c 两条入口现状**：
- 常驻腿（`wisp` 无参）：`runResident`（`cmd/wisp/resident_windows.go:30-234`）全文无 config 读取；建球点 :163 `startResidentBall(rt.Registry, ra.vetoByEsc, withPanelHost(...))`——参数表无 config、无 Manager。热键＝写死四枚（:171）。
- `wisp run`：**不建球**。尺：run.go 产码零 `ball.` 命中（注释全在讲 veto 通道缺失；run.go:552 逐字 "console run: no floating ball, no global Esc hook"）。所以"`wisp run` 建球时热键从哪来"今天**没有这格**——run 腿根本没有球。
- 常驻进程里的 config Manager 辨析：常驻进程有两处 Manager——(1) 任务管线 `startResidentTaskSource → assembleRuntime`（resident_task_source_windows.go:265）复用 run.go 装配，`config.NewManager` 在 `run.go:409`、`rt.mgr` 在 :425；(2) 面板链 `newComposerDispatchChain`（`cmd/wisp/panel_inbound.go:228-233`）另建一枚。但两处都**不外露给球**：`startResidentBall` 调用（resident_windows.go:163）发生在任务源装配（:206）**之前**，`PanelManager` 结构（`cmd/wisp/panel_host_windows.go:137-149`）只持 dispatch 不持 mgr。
- **配置缺失/该节缺失时今天会发生什么**：常驻腿不读 config ⇒ "配置缺失"对球零效果，球永远 DefaultHotkeys() 四枚（经 :171 直给，连 :153 的零值兜底都不触发）。若改走 config 路径：缺节＝三枚 `""`＋cancel `"Esc"`（Q1a），`registerAllWith` 对 `""` 分类 `HotkeyDisabled`（`hotkey_windows.go:467-469` "disabled by config: never attempted"）——即"缺节直接变无键"，正是 ApplyHotkeyDefaults 注释点名的坑。

## 2 改过之后谁重新注册（两形代价，不选形）

**Q2a 全仓 `NewHotkeyReloader` 调用点逐枚**（尺：`grep -rn NewHotkeyReloader cmd internal tools docs scripts --include=*.go` 排 `.scratch`/`_test.go`）：
- 非测试调用点＝**1**：`cmd/balldebug/main.go:237` `bridge := ball.NewHotkeyReloader(b, b.ConfiguredHotkeys(), ...)`；`:243` `bridge.Refresh = func() error { _, err := mgr.CheckAndReload(); return err }`；`:244` `mgr.OnReload = bridge.OnReload()`；`:250-253` 桥跑在 `observe.Default.Spawn("balldebug-hotkey-bridge", ...)` 的 1s `bridge.Run(ctx, time.Second)`（花名常量 `main.go:45`）。
- 其余命中全是定义/文档/测试：`internal/ball/hotkey_reload.go:16`（doc 例）、`:70/:73`（定义）、`hotkey_live_test.go:159/:498`、`hotkey_status_test.go:458/:511`。
- `RebindHotkeys` 非测试业务调用者＝0：唯一调用在 `HotkeyReloader.Check` 内部（`hotkey_reload.go:96`）。票面"唯一 reloader 调用者在调试台件"**成立**。

**Q2b `config_reload.go:153` 之后有没有产码把新 `[hotkey]` 交给球**：
- **没有**。链路＝`reloadOnce`（config_reload.go:152-162）→ `rt.mgr.CheckAndReload()`（:153）→ `reportReload(rep)`（:161→167-208）；reportReload 只 `rt.auditf`＋`fmt.Fprintf(rt.stdout)`（:168-171），全文件零 ball 引用。`wisp run` 没有球，无从交起；常驻进程的任务管线虽带出 rt.mgr，同样走这段、同样无球交接。
- `OnReload` 回调（manager.go:198-199）只吃 `rep.Reload`（reload-tier）；`[hotkey]` 在 hot 档进 `rep.Hot`（manager.go:293），**不触发任何回调**。且常驻/面板宿主均未赋 `OnReload`（尺：`grep OnReload cmd/wisp` 非测试仅 panel_config_store.go:27 一句注释）。

**Q2c `manager.go:278` 那个闭包今天到底做什么**：
- 只换内存值：`func() { cur.Hotkey = fresh.Hotkey }`，由 `plan.commit()` 执行（manager.go:193，`CheckAndReload` 内重新上 mu 后）。之后无人被通知：hot 表注释（manager.go:271-289）逐字 "Everything else is hot-tier: apply wholesale on change"，无事件机制；`m.OnReload` 只吃 rep.Reload。内存值跟着文件变、消费方一个没有——`mgr.Config().Hotkey` 唯一非测试读者＝balldebug 的桥闭包（main.go:238）。

**Q2d 两形代价与线程约束（只量，不选形）**：

| 维度 | 形 A：把 reloader 接进常驻腿 | 形 B：reload 那一跳多调一次注册 |
|---|---|---|
| 现成件 | `HotkeyReloader` 全套（diff/幂等/panic 隔离/Problems 落日志，hotkey_reload.go:81-105）＋映射闭包样板（balldebug:237-242）＋`ApplyHotkeyDefaults` | `RebindHotkeys`（ball_windows.go:813-825）＋`ApplyHotkeyDefaults`；diff 由 host 自己做 |
| 可见性缺口 | 需 `*ball.Ball` 与一个 config 源同处可见：rb 建于 resident_windows.go:163（装配早期），mgr 在任务管线/面板链内（Q1c），两对象不同层 | 需 `agentRuntime` 持球句柄：`agentRuntime`（run.go:266-326）无 ball 字段；且 run 腿无球，形 B 在 `wisp run` 是空集，真实落点只能是常驻腿任务管线一侧 |
| 驱动源 | 两条现成路：(a) `bridge.Run(ctx,1s)` 独立 goroutine（balldebug 形状，须走 Registry.Spawn 取 recover 边界，balldebug:250-253 是样板）；(b) 从已有 1s watchdog tick 调 `Check()`（hotkey_reload.go:24-25 逐字支持："on a host that already polls config every tick (the wisp watchdog), call r.Check() from that tick"） | 零新增 goroutine：挂 `reportReload` 的 hot 分支（`rep.Hot` 含 `"hotkey"` 时） |
| 线程约束（只引形状，不判改否） | `Check()` 注释（hotkey_reload.go:79-80）逐字 "Safe from any goroutine EXCEPT the ball's UI thread (the rebind is synchronous on it, so calling from there would deadlock)"——不许从 ui-sta 线程调；`RebindHotkeys` 经 `b.uiRun` post 到球 ui-sta 并等待（ball_windows.go:813-825；uiRun 定义 ball_windows.go:741-752：非 ui-sta 调用 post-and-wait、已在 ui-sta 则内联）。球所在 ui-sta＝D38b 冻结花名册第 1 枚（`internal/observe/goroutine.go:43-45`）。⚠ 票 33 十一裁形状＝球/面板各在自己线程、投 `ui-sta` 会冻外层泵——本枚只引：rebind 的 post-and-wait 目标是**球的** ui-sta，watchdog 与 reload tick 均不在 ui-sta 上 | 同左（RebindHotkeys 同一 uiRun 路径）；`reloadOnce` 跑在 watchdog goroutine、非 ui-sta。静态读数：`CheckAndReload` 的 commit（manager.go:192-196）在 mu 内、返回在锁外，reportReload 侧再触 Rebind 无 config 锁叠加 |
| churn/幂等 | `Check` 先 diff（hotkey_reload.go:93-94 "nothing changed: do not touch Win32"），每 tick 轮询零 Win32 消耗 | `RebindHotkeys` 语义＝unregisterAll+registerAll 全量（ball_windows.go:816-817）；注释 :803-805 逐字 "Rebind is idempotent... HotkeyReloader still diffs so the Win32 churn happens once per change"——形 B 每 tick 盲调会 churn，须 host 侧先比对 |
| 已登记残留（两形共有） | rebind 丢 in-flight Esc borrow：ball_windows.go:807-812 逐字 "A rebind drops an in-flight Esc borrow... A host that rebinds while a card is waiting therefore has to call TakeEscForCancel again on the next state sync - no caller does that today... the residual is named in ticket 245" | 同左 |
| 错误路径 | Refresh 失败→Warn 保留现值（hotkey_reload.go:83-85）；rebind 后 Problems 自动 slog.Error（:98-100） | reportReload 现有 hot 句（config_reload.go:171）只报"已生效"，注册失败的说法要 host 自己接 `rep.Problems()` |

## 3 注册失败时那句话＋三枚钉的射程

**Q3a 读数从哪来**：
- `cmd/wisp/resident_ball_windows.go:194` `rep := b.HotkeyReport()`；`:207-208` verdict＝`fmt.Sprintf("the floating ball window is up in this process (tray icon added, hotkeys live %d/4: %s)", len(rep.Live()), hotkeySummary(b))`（:209-212 同数落 slog）。
- `HotkeyReport()`（`internal/ball/ball_windows.go:830-834`）＝"最后一次注册遍（boot 或 rebind）的产物"，uiRun 同步读 `b.hotkeyReport` 字段（:123）。
- `Live()`（`internal/ball/hotkey_windows.go:313-321`）＝`Status == HotkeyLive` 的 id→Accelerator map。**语义＝Win32 现在真握在手里的枚数**——不是"被告知的集合"（那是 `ConfiguredHotkeys()`，ball_windows.go:836-844），也不是"健康集合"（那是 `AllLive()`，:344-351，对 Live/Disabled/Standby 三态都算好）。

**Q3b 四枚里一枚被占时那句话怎么说**：
- 被占枚在 `registerAllWith` 分类 `HotkeyTaken`（`hotkey_windows.go:479-481`，`ERROR_HOTKEY_ALREADY_REGISTERED` 分支），单枚失败不拦其余（:411-414 注释）。它进 `Problems()`（:371-375，逐字 "hotkey %s = %q is occupied by another program and was NOT registered; pressing it will do nothing until you pick a free combination in [hotkey]"），由 `resident_ball_windows.go:202-205` 在 summary 前逐行 slog.Error＋stdout。verdict 例：summon 被占 ⇒ `hotkeys live 2/4`。idle 理想值 3/4：hkCancel 走 `cancelIdleLine` 从不注册（hotkey_windows.go:456-464；resident_ball_windows.go:127-130 逐字 "THREE at idle"）。
- **射程判断（要害）**：Problems 那句让用户"去 `[hotkey]` 挑一个没被占的组合"——但常驻腿今天**不读 `[hotkey]`**（问 1），这句话指向一扇这条腿没有的门。四枚句子（`hotkey_windows.go:353-382`）覆盖"注册结果"（unparsable/taken/refused 逐枚具名），**盖不到"改配置以谁为准"**：没有任何一句说"现在生效的值来自哪里"（config 还是 DefaultHotkeys）。`hotkeySummary`（resident_ball_windows.go:316-326）打 `name=binding status`，binding 是球被告知的字符串——如实反映"球被告知了什么"，不反映"这个告知该不该跟 config 走"。

**Q3c 三枚钉逐枚读断言内部（非只看符号名）**：
- `cmd/wisp/resident_ball_228_windows_test.go:45`：是**常量定义** `ballStoppedMsg = "ball: resident leg destroyed the ball window, its tray icon and its hotkeys"`。全文件 244 行**无一处** `hotkeys live`/`Live()` 读数；两测试（:53、:184）断言 console 三态句（ballUpClaim/ballAbsentClaim :47-48）与落盘记录（created/refused/stopped/install）的存在性、唯一性与次序（:117-167）。**射程＝"球被建/被拆且说了真话"，与热键绑定值零关系。**
- `cmd/wisp/resident_ball_live_228_windows_test.go:144`：`strings.Index(verdict, "hotkeys live")` 截 80 字符进 `t.Logf`（:147），**无断言**——纯日志增强。
- 同文件 `:160` `strings.Contains(verdict, "hotkeys live 0/4")` → Errorf（注释 :148-151 点名票 64 A1b："a ball that registered nothing still looks like a working window"）。
- 同文件 `:163` `strings.Contains(verdict, "hotkeys live 4/4")` → Errorf（注释 :152-157 点名票 245：idle 球不该四枚全握）。
- **:160/:163 合并射程**＝只钉**枚数**两个极端（0/4 与 4/4），对 1/4、2/4、3/4 **无断言**；且读的是"球自己打印的数"，注释 :158-159 逐字承认 "The count is read from the number the process prints, which is len(HotkeyReport().Live()) - the registration set, not a sentence about it"。两钉盖不到"改配置以谁为准"：就算常驻腿改成读 config，`hotkeys live 3/4` 照样打、两钉照样绿；就算它继续用 DefaultHotkeys，两钉也绿。**它们是"注册结果"的钉，不是"绑定来源"的钉。**
- ⚠ 若未来修法改变 idle 枚数语义或改写那句 verdict 字面，:144/:160/:163 的字面尺会红——票面已勾"本票不许为落地去改它们"，这与 AC#2 正控（种 config 后新值注册）不冲突：正控量的是绑定值，两钉量的是注册结果。

## 4 我可能判错的条目

- `Live()` 零值语义：若 map 为 nil vs 空（ball_windows.go:818 `b.registeredHotkeys = rep.Live()`），我对 "0 枚" 的读法只经 `len(map)`，没追 nil-map 边界——不影响三问结论。
- `deepCopyConfig` 对 `HotkeySection`（纯 string 字段）拷贝无疑义，但未逐行读 copyValue 的 interface 分支。
- 形 B 静态锁读数（"reportReload 侧再触 Rebind 无 config 锁叠加"）只读了 manager.go 的 mu/reloadMu 与 CheckAndReload 返回点，未穷举 confirmLockedLoosening 挂起期间的状态。
- 票 33 十一裁的原文我没有直接读（裁决文件不在本次 grep 根清单内），线程形状一句是从 AGENTS.md 禁区描述＋resident_ball_windows.go:122-125/141-148 注释＋panel_resident_windows.go:207-218（panel 线程 LockOSThread/STA）三角引的。
- winlive 标签（`//go:build windows && winlive`）下的行为我只静态读了源码，未实跑（⛔ 禁 Go 命令）。

## 5 判不动／量不到的地方

- **实跑读数一格都没有**：`hotkeys live N/4` 在真机上的实际 N、被占枚的真实 Win32 返回、`bridge.Run` 的真实轮询行为——全部因 ⛔ 禁 Go 命令量不到，具名弃权，不用"应该没问题"填空。
- 票 33 十一裁"要不要改"的判断：本枚只引形状（Q2d 线程行），⛔ 不判。
- "两形哪个该选"：只交代价表，选形归编排者。
- D43 转移表、SLO/golden、thresholds：只读禁令，未触碰、未判。
- `resident_ball_228_test.go`（非 windows 的 shape walk 文件）只读了 :15/:89/:212 三处 grep 上下文，未全文读——它断言"调用点存在"，与三问正交，风险低但如实登记。

## 6 我推翻／更正编排者哪一句

- **没推翻**。票面五条现量逐条复认全中：schema.go:176/manager.go:278/resident_ball_windows.go:171/balldebug:237/:243/config_reload.go:153 与 :101-108 注释，均逐字核过。
- 两处**补充**（非推翻）：(a) 票面 AC#0① 说"`wisp run` 与常驻腿两条入口建出来的球"——今天 `wisp run` **没有球**（run.go 产码零 ball 引用），AC#1 的"两条入口"对 run 腿是空集，落地时这半句指向的是"未来 run 腿若带球"的形状；(b) 票面现量 2 说"`Config().Hotkey` 从没进过构造参数"——成立，且常驻进程里其实存在两枚 config.Manager（任务管线＋面板链），只是都到不了 `startResidentBall` 的参数表，"从没进过"的机制比"没有 Manager"更具体。

## 末节（一句话）

常驻腿建球时热键写死 `DefaultHotkeys()`（:171）、`internal/ball` 无 config 入口、`[hotkey]` 改动只换内存（manager.go:278）无人重新注册、全仓唯一 reloader 调用者在 balldebug:237、`hotkeys live N/4` 如实报注册结果但三枚钉全不盖"以谁为准"——三问答满，读数全部带 file:line，实跑格全部具名弃权。
