# 248-c3 残余格分诊普查件（只读腿，零 Go 命令）

> 本件是票 248「面板设置那条路」的**残余格分诊**普查腿 `248-c3`。只读源码＋grep/find＋git log/show/diff；
> **一枚 Go 命令都没跑**（`go build`/`vet`/`test` 全禁，见票面 §4 排程与派单第一硬规）。
> 判据一律问调用点枚数；`rc=0` 证明不了接线。凭据值零外泄（只写变量名/字段名/blob 名/env 名）。
> ⛔ `frontend/**`／`design/**` 两层禁令：不读、不引、不转述。

## §0 起手锚

同一发命令取（`2026-10-03 09:25:27+0800`）：
- `date` = `2026-10-03 09:25:27+0800`
- `git log -1 --format="%h %ci"` = `850a76b4 2026-10-03 09:25:19 +0800`
- `git status --porcelain -- cmd internal tools docs scripts .scratch | wc -l` = `367`

⚠ 起手期间三枚写腿/验收腿在飞，工作树活跃漂移：同一批锚在两秒内 `HEAD` 已从 `5c3c22d8 09:24:05` 移到 `850a76b4 09:25:19`、porcelain 计数从 `369`→`367`。这正是本腿禁跑 Go 命令的现场理由（会撞对方编译、洗读数）。本件所有 file:line 一律本腿现读现取，不沿用他腿历史行号。

## §1 残余格逐枚分诊表

> 票面现量：未勾 **5**（AC#2／AC#4／AC#8／AC#9／AC#10），已勾 **7**（AC#0/1/3/5/6/7/11）。本节只分诊那未勾的五枚。
> 所有 file:line 为本腿现取（`git grep … HEAD`／`git show HEAD:`），因工作树在写腿飞动中、行号会漂，⛔ 不沿用 v1c／r1 历史行号。
> 行号更正预告：v1c 引 `manager.go:281`＝llm 热应用锚，本腿现量该表项已漂到 `:284`（`:281` 现为 `hotkey`），见 §5。

**AC#2｜题面逐字**："出向快照里出现'凭据是否已录入／配置是否可读'这一维时，全仓任何日志与快照产物里 grep 不到任何 key 值……这一发要作为常驻用例进仓"。
- 缺的是：【属前端那半】。Go 半边全写好且已接线；唯一没结的是页面对那两枚新键的声明。
- 盘上现证：维度两枚 json tag `internal/panel/composer.go:268-269`（`credentialState`/`credentialKnown`，声明在 `ComposerState` 内＝composer 段、无顶层第五键）；无 reader 时 `:288` 折回 `CredentialUnknown`；reader 钩子 `internal/panel/pump.go:146 Credential func() CredentialState`，`:260` 仅当 reader 非 nil 才填；生产接线点 `cmd/wisp/run.go:710 Credential: rt.settings.credentialStatus`（该 reader 的**唯一生产调用点＝1 枚**）；reader 本体 `cmd/wisp/panel_config_store.go:285 credentialStatus()`。零命中尺那枚常驻用例在 `cmd/wisp/panel_config_248_test.go`（v1c §2 逐名跑过含 `TestAC2*`）。卡住的那道门＝`internal/panel/composer_test.go:48 TestComposerContractTypesMatchFrontend`：它 `os.ReadFile` 一枚前端 TS 类型文件、把 Go emit 的 JSON 键集与前端 `interface ComposerState` 做**双向减集**，Go 多出的 `credentialState`/`credentialKnown` ⇒ 红。**这一发红出自 v1c 那腿跑的读数（本腿禁跑，只读了该 test 的比对机制），不是本腿读数。**
- 要落地最少动哪几行：Go 侧 0 行（已就绪）。差的是页面侧那枚 `interface ComposerState` 补两字段声明 ⇒ 门转绿。抄给页面侧的问句见 §4-Q1。

**AC#4｜题面逐字**："一条完整链＝面板点设置 → 填入参数与凭据 → 保存 → 下一次任务真用上刚配的模型……⛔ 在票 33 的宿主真起来之前这一格不许勾，也不许用'我手工改了 config.toml'来代替那一次点击"。
- 缺的是：【依赖票 33 那块宿主】＋真机读数＋【属前端那半】，三样叠在一格。
- 盘上现证：宿主**线程建了、窗没开**——`cmd/wisp/resident_windows.go:142 newResidentPanelManager`、`:148 startResidentPanel` 各 1 枚生产调用；`cmd/wisp/panel_resident_windows.go:138` 注释逐字 "startResidentPanel builds the thread but NOT the window: nothing is created until the user asks (a ball gesture, a tray item, a test)"。开窗那枚 `bringUp` 定义在 `cmd/wisp/panel_host_windows.go:230`，**除自身文件外零枚生产调用点**（尺见 §2）。前端产物缺：`git ls-files frontend/dist`＝1 枚 `.gitkeep`（见 AC#9）。同一次运行重建那条今天不通：`cmd/wisp/panel_config_store.go:27` 注释逐字 "once at assembly (run.go) and never rebuilt by a reload - OnReload has no [listener]"（J9 现量仍在 HEAD 成立）。
- 要落地最少动哪几行：不是本票/本腿能落的量——需 (a) 票 33 把 `bringUp` 接进"球手势/托盘 → rp.tasks → loop 调 bringUp"那一跳且真机不 panic（A492 那条故意没接的坑）；(b) 前端把 dist 填上一包；(c) 一次真机点击＋下次重启真用上。Go 装配侧本票已备好（setter/reader/派发俱全）。真机读数归谁见 §4。

**AC#8｜题面逐字**："票面与回执文案不许出现'保存即生效'，且'要重启'必须到达页面可见面"；编排者 §126 补："AC#8 那种'不许说谎'的判据遇到'说了另一句不成立的实话'同样算不过"。
- 缺的是：【产码】（回执档位须由登记表同源产出，归票 255-r2）＋【待人裁定】（翻勾交非实现者）；"页面可见"那一半共用 AC#4 的传输路（票 33＋前端），但编排者把本格翻勾条件钉在 255-r2 而非前端。
- 盘上现证：两套档位词表对不上——`cmd/wisp/panel_config_store.go:228` 与 `:275` **无条件** `res.Tier = panel.EffectiveRestart`，前面零按键/零按段判断；档位定义与成句在 `internal/panel/config_handlers.go:199-201`（`EffectiveNow/NextTask/Restart`）与 `:428-434 tierSentence`，但 `EffectiveNow/EffectiveNextTask` **非测试写者 0 枚**（§2）；登记表侧 `internal/config/tiers.go:30 "llm":"hot"`、`:34 "panel":"hot"` 声明这两段是热应用段（`manager.go` 热应用表 `:284 {"llm",…}` 与之同断言）。查档读口 `internal/config/tiers.go:93 TierOf` 在 HEAD **零枚生产调用者**；正解雏形正被写——未提交的 `cmd/wisp/config_readers_255.go:174` 注释逐字 "TierOf's first production caller" ⇒ **255-r2 在飞、未落库，不能替它勾**。
- 要落地最少动哪几行：255-r2 把 `:228/:275` 的硬填换成 `config.TierOf(写入键路径/段)` 得档位（restart→"要重启"、hot→诚实的立即/下次形），再交非实现者裁。本票 0 行独立可动——§127 明令"不许硬写死一句漂亮话，要与登记表同源"。

**AC#9｜题面逐字**："embed 之后的文件系统里条目数与关键入口文件必须真存在（⛔ 一枚 `.gitkeep` 就能让 `go build` 绿，那不算证据）"。
- 缺的是：【属前端那半】。判据问的是 embed 之后文件系统里的产物条数/入口文件，而那在两层禁令里（不读、不引、结论也不落那两层行号）。
- 盘上现证：`git ls-files frontend/dist` 现量 **1 行**、逐名 `frontend/dist/.gitkeep`；`frontend` 总跟踪 85 枚（多属源码非 dist 产物）。名册级事实＝一包页面产物**不在**。embed 那一行不引（禁令层内）。本腿⛔不据"能构建"反推"有内容"（票面 AC#9 原文禁止）。
- 要落地最少动哪几行：Go 侧 0 行——"谁把 dist 填上"没落定前勾不了（票面同判）。抄给页面侧的问句见 §4-Q2。

**AC#10｜题面逐字**："常驻那条腿今天吃不到 `[risk]` 那两项配置，本票要么接上、要么在票面写明它是常量——⛔ 不许留成'配置页改了它就变了'的错觉"；§153 定案 ⓑ："改由设置回执逐字写明 `[risk]` 那两项只作用于跑任务的进程，常驻腿今天用常量 300s／3s"。
- 缺的是：【产码】（ⓑ 那句话由登记表同源产出，归 255-r2，与 AC#8 同一块石头）＋【待人裁定】（选形 ⓑ 已由编排者 §153/台账 A534 裁死、ⓘ 另立票 256，但翻勾要非实现者）；"页面可见"那一半同 AC#8 的传输路。
- 盘上现证：常驻腿 `cmd/wisp/resident_approval_windows.go:109-113` `approval.New(approval.Options{ UI, Channels: NewChannels(), Logf })`——**无 `Window`、无 `ApprovalTimeout`、无 `Grants`** ⇒ 落 `internal/agent/approval/queue.go:107 DefaultApprovalTimeout=300s`、`:116 DefaultL1Window=3s`、`:122 MaxL1Window=3s` 常量；跑任务那腿 `cmd/wisp/run.go:615-616` 才真读 `cfg.Risk.L1WindowSec`/`ConfirmTimeoutSec`。页面写不进 risk：`internal/panel/config_handlers.go:108 lockedFieldFamilies` 含 `"risk."` ⇒ 错觉面在"读侧/display 或未来可写"那一侧，不在 config.set 写侧。GRANT-DROPPED（`Grants` nil）那发的**实跑读数本腿量不到**（要真跑常驻腿，见 §6）。
- 要落地最少动哪几行：255-r2 的登记表同源句产出 ⓑ 文案（或本票把它登记为常量事实），交非实现者裁。撤销口令「248 AC#10 改 ⓘ」仍挂着——不移动 `approval.New`。

## §2 「写了没接」vs「一块没写」两套名册

（取数中）

## §3 依赖票 33 的那几格 vs 今天能独立落的那几格

（取数中）

## §4 真机才量得到的读数清单（要谁做）

（取数中）

## §5 我可能写错的条目（自我对抗）

（取数中）

## §6 量不到的地方（具名）

（取数中）

## §7 交件判语

（取数中）
