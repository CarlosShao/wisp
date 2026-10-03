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
- 缺的是：【依赖票 33 那块宿主】＋真机读数＋【属前端那半】，三样叠在一格。⚠ **不是"窗没接"**——开窗那一跳今天已生产接线（详见下与 §2），本腿初判曾把它读成"线程建了窗没开"，就地更正。
- 盘上现证：
  - 开窗链**已接**：`cmd/wisp/resident_windows.go:148 startResidentPanel`（1 枚生产调用）→ `cmd/wisp/resident_windows.go:163` 把 `withPanelHost(func(via) bool { return panel.RequestToggle(via) })` 交给球宿主；`panel_resident_windows.go:376` 注释逐字 "RequestToggle is what the panel hot key and the tray item drive"；`RequestToggle`(`:378`)→`RequestShow`(`:329`)→投 `rp.tasks`(`:315`)→`loop`(`:232 case fn := <-rp.tasks:`)→`showOnThread`(`:360 err := rp.mgr.Show(...)`)→`PanelManager.Show`(`panel_host_windows.go:395`)→`bringUp`(`:230`, 其直接调用点是同文件 `:402`)。⇒ `bringUp` 的**直接**生产调用点只有同文件 `:402`，但**传递可达**且由热键/托盘驱动，`startResidentPanel` 那句 "nothing is created until the user asks"＝**设计上延迟到手势**，不是没接。
  - 仍缺的三样：① **真机才知**这扇窗真开时会不会撞上嵌套消息泵重入（`:325` 声称 "Safe from a ui-sta callback"、`panel_resident_windows.go:15-16` 却警告嵌套泵会乱序球的手势队列——A492 那条老坑到底解没解，本腿与 v1c §120 "三枚新测试文件没有任何一枚真开窗" 都没量，⛔ 不跑）；② 前端产物缺——`cmd/wisp/panel_host_windows.go:368 serveEntry` 要 `m.assets.Built()` 才推页面，否则 `:362` 只 serve 一枚 "panel assets unavailable: the embedded bundle is not built" 占位，而 `git ls-files frontend/dist`＝1 枚 `.gitkeep`（＝AC#9）；③ 设置页入口＋写回＋真机点击＋重启真用上（`panel_config_store.go:27` 注释逐字 "never rebuilt by a reload - OnReload has no [listener]"，J9 仍在 HEAD 成立）。
- 要落地最少动哪几行：Go 开窗链本票/票 33 已备；差的是 (a) 33-v1/真机确认那一跳不 panic、(b) 前端把 dist 填上一包、(c) 一次真机点击＋下次重启真用上。这些都不是本腿能产码或能跑量的。

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

> 判据＝生产调用点枚数（`git grep HEAD`，非测试文件）。⛔ 本腿一枚 `go build`/`vet` 都没跑，`rc=0` 证明不了接线；下面每枚都只写本腿实际 grep 到的枚数。

**2A 「写了但没接线」（符号在 HEAD 已落库、非测试生产调用点＝0）**
- `internal/config/tiers.go:93 TierOf(path) (tier, ok)` — 尺＝`git grep "TierOf" HEAD -- cmd internal` ⇒ HEAD 只有**定义行 + 它自己的注释行**，生产调用点 **0 枚**（连测试都没有）。它是 AC#8/AC#10 的"按键查档"读口：函数写了、能编译，但没人调 ⇒ 回执档位今天靠硬填。
- 正解雏形此刻在**未提交工作树**里长：`cmd/wisp/config_readers_255.go:174` 注释逐字 "TierOf's first production caller"（该文件 `??` 未跟踪，属飞动的 255-r2）。⇒ 属"另一枚腿正在接、本腿不能替它勾"，不算本票可独立落的格。

**2B 「写了且已接线」（函数在、非测试生产调用点 ≥1，⛔ 别当没接）**
- `cmd/wisp/panel_config_store.go:285 credentialStatus()` — 唯一生产接线点 `cmd/wisp/run.go:710 Credential: rt.settings.credentialStatus`（**1 枚**）⇒ AC#2 的 reader 已接进常驻快照泵。
- `cmd/wisp/panel_host_windows.go:230 bringUp()` — 直接调用点是同文件 `:402`（在 `PanelManager.Show` 内）；`Show` 被 `panel_resident_windows.go:360 showOnThread` 调；`showOnThread` 经 `rp.tasks` 由 `RequestShow`(`:329`)/`RequestToggle`(`:378`) 投；`RequestToggle` 的生产驱动点＝`cmd/wisp/resident_windows.go:163 withPanelHost(...RequestToggle...)`（交给球的热键/托盘）。⇒ **传递可达、非"没接"**，只是创建按设计延迟到手势。AC#4 因此不是"Go 侧一块没接"（见 §1 AC#4 更正）。
- `internal/agent/approval`：`run.go:615-616` 真把 `cfg.Risk.L1WindowSec/ConfirmTimeoutSec` 灌进 `Window`/`ApprovalTimeout`（跑任务那腿，接了）；`resident_approval_windows.go:109-113` 那枚 `approval.New` **故意不接**这三项（`Window`/`ApprovalTimeout`/`Grants` 缺省）⇒ AC#10 的前提，见 2D。

**2C 「写了枚举但零非测试写者」（值定义了、没人赋）**
- `internal/panel/config_handlers.go:199 EffectiveNow`／`:200 EffectiveNextTask` — 尺＝`git grep "EffectiveNow\|EffectiveNextTask" HEAD -- cmd internal | grep -v _test` ⇒ 只命中**两枚常量定义**（199/200）与 `tierSentence` 的 **case 分支**（430/432），**赋 `res.Tier` 为它们的＝0 枚**（回执永远只赋 `EffectiveRestart`，`panel_config_store.go:228/:275`）。⇒ 票 255 那句"别为了好看把这两枚零写者的值用起来＝造第五套"正是针对此形状。

**2D 「一块没写」（HEAD 上该产物/该调用点根本不存在）**
- 回执里那枚"按键/按段查 `TierOf` 得档位"的**调用点**在 HEAD 不存在（只有硬填 `EffectiveRestart`）；被调的函数存在但 0 调用者 ⇒ AC#8/AC#10 缺的是**那一枚调用点**（255-r2 正在写），不是"整块没名字"。
- `cmd/wisp/resident_approval_windows.go:109-113` 里 `approval.New` 的 `Grants`/`Window`/`ApprovalTimeout` 三项赋值＝**一块没写**（ⓑ 裁定不补，改在回执文案里写常量事实）。
- ⛔ 前端两块（`interface ComposerState` 的 `credentialState`/`credentialKnown` 声明；`frontend/dist` 那一包页面产物）本腿**不读禁令层、无法自证其形状**——只从名册级尺（`git ls-files frontend/dist`＝1 枚 `.gitkeep`）与 v1c 跑那枚对账门红的读数**间接**知其未声明，⛔ 不署本腿的名、不引那两层任何行号。

## §3 依赖票 33 的那几格 vs 今天能独立落的那几格

> 硬要求：⛔ 不许把"等 33"当挡箭牌盖住本可以自己落的格。逐枚说死真依赖是谁。

**真·等票 33 的：只有 AC#4 一格。**
- AC#4 的宿主**代码**今天已大体在（`resident_windows.go:163 withPanelHost→RequestToggle→…→bringUp`，见 §2B），但"这扇窗在真机上真点开时会不会撞嵌套消息泵重入 panic"是 33-v1/真机那一发，票 248 这侧替不了。⇒ AC#4 卡的是**票 33 落地＋验收**，不是"Go 一块没写"。

**不·等票 33、而是等票 255-r2 那块产码的：AC#8、AC#10。**
- 这两格的翻勾条件编排者已钉死（§126／§153／台账 A530/A534）＝"回执那句档位/常量事实必须由同一份登记表产出"，而那张表就是票 255 的 `TierRegistry`／`TierOf`。255-r2 正把这枚读口接进回执（未提交的 `config_readers_255.go:174` ＝ TierOf 首个生产调用者）。⇒ 这两格**跟票 33 无关**、跟前端也无关（就其 Go 判据而言），落定就等 **255-r2 提交＋一枚非实现者裁**。⚠ 别把它们误并进"等 33"那堆：那会让 255-r2 交完后没人回头翻这两格。
- ⚠ 两格都还带一枚"页面可见"的传输尾巴（回执字符串要经票 33 的 bind 闭包真渲到窗上）——但那半枚与 AC#4 共用同一发真机读数，编排者没把它单列为 AC#8/AC#10 的翻勾前置（见 §4 归口）。

**不·等票 33、也不等 255、只等前端那半的：AC#2、AC#9。**
- AC#2：Go 半边全写好、已接线（`run.go:710`），唯一拦路是 `composer_test.go:48` 那枚双向对账门——它读一枚前端 TS 类型文件，前端 `interface ComposerState` 补上 `credentialState`/`credentialKnown` 两枚字段声明即转绿。⇒ 与票 33/255 **无关**，前端 agent 今天就能独立做。问句见 §4-Q1。
- AC#9：拦路是"谁把 `frontend/dist` 填上一包页面产物"（现量只剩一枚 `.gitkeep`）。⇒ 与票 33/255 无关，属前端交付物。问句见 §4-Q2。

**本腿对"今天能不能独立落"的总结论：** 五枚里**没有一枚**能由"248 落地腿"单独勾——AC#2/AC#9 要前端（本编队⛔不写 `frontend/**`），AC#8/AC#10 的 Go 句归 255-r2（⛔ 硬写死被 §127 禁），AC#4 要 33＋前端＋真机点击。但"不能由 248 腿独立勾"≠"都在等同一样东西"：真依赖分三类（33／255-r2／前端），上面已逐枚说死，⛔ 不许一句"等 33"糊过去。

## §4 真机才量得到的读数清单（要谁做）

> 凡"面板真点开／真存盘／下次真用上"这类只有真窗口才量得到的，本腿一律只登记、不填值（⛔ 一枚 Go 命令都没跑，也开不了窗）。

**要哪些真机读数：**
1. **AC#4 完整链**：面板点开 → 填参数与凭据 → 保存 → 下一次任务真用上刚配的模型（票面明令 ⛔ 不许用"手工改 config.toml"代替那次点击）。⇒ **owner 眼睛签收**（真桌面交互＋真点击）。前置：票 33 落地＋`frontend/dist` 有一包真页面。
2. **AC#4 的"窗真开不 panic"那一发**：热键/托盘 → `RequestToggle` → `bringUp` 真创建窗时会不会撞上嵌套消息泵重入（`panel_resident_windows.go:15-16` 警告的那条乱序，A492 老坑解没解）。⇒ **编排者本机可自跑**（跑起常驻进程、按热键看窗开没开/进程 panic 没）；但要坐实"机主点得出设置"仍归 1（owner）。本腿与 v1c §120 都没开过窗。
3. **AC#8／AC#10 的"页面可见"尾巴**：回执那句档位/常量文案真渲到窗上、机主看得见。⇒ 依赖 1/2 那扇窗，**owner 眼睛**；但其 Go 判据（回执字符串内容对不对）在派发/回执字符串层就能测（248 自己的用例即走这层，不开窗），编排者把它翻勾钉在 255-r2 的 Go 产出＋一枚非实现者裁，⛔ 没把"真窗可见"当这两格的翻勾前置。
4. **AC#10 的 `GRANT-DROPPED` 实跑**（票面 §59 那族：用户在球上按"本会话内允许"，今天因 `Grants` nil 只放行不落盘）：⛔ 本腿没跑常驻腿、量不到这一发的实跑证据；要真跑常驻腿＋球交互才采得到。⇒ 归 1 同一次真机 session，或票 224 r2 那套授权仪器（票面 ⓑ/ⓘ 判据点名复用，⛔ 不新造）。

**可以抄给页面侧的问句**（owner 自己带给他那枚前端 agent；本腿⛔不读 `frontend/**`、不联系其它会话）：

- **Q1（对应 AC#2）**：Go 侧 `ComposerState` 快照现在多发两枚 JSON 键 `credentialState`（值枚举：已录入/未录入/配置不可读/未知）与 `credentialKnown`（布尔，区分"没人读过"与"没录过"）。那枚由 `TestComposerContractTypesMatchFrontend` 双向对账、解析你们 composer 状态 interface 的门，红句报这两枚"interface 未声明"。⇒ 问：你们那侧描述 composer 状态、被对账门解析的那枚 interface，能否补上 `credentialState`／`credentialKnown` 两枚字段声明（字段名一字不差）？补上对账门即绿，AC#2 转成立。
- **Q2（对应 AC#9）**：`git ls-files frontend/dist` 现在只剩一枚 `.gitkeep`；`//go:embed all:dist` 过得了构建但没有任何页面产物可发。⇒ 问：`frontend/dist` 那一包真页面产物（含关键入口文件）由谁产出、何时填上？设置页要真点得出来，前提就在这包。⛔ 票面 AC#9 原文禁止"能构建"被读成"有内容"。
- **Q3（对应 AC#4 页面半）**：Go 侧入向已备两枚方法 `config.get`／`config.set`（参数字段是**具名枚举**、⛔ 非任意 key 透传；凭据那枚具名字段只写、值绝不回显），出向已带"凭据是否已录入"那一维。⇒ 问：页面侧是否已做票 §5 列的三块——①一个"设置"入口、②能填服务商地址/模型名/上下文窗口/价格/权限档位的表、③一枚"只写不回显"的密钥输入框（填后任何时候不再显示值、只显示已录入/未录入，且⛔不把密钥塞进聊天消息）？保存要走 `config.set` 那枚"配置写"方法。

## §5 我可能写错的条目（自我对抗）

1. **AC#2 那枚对账门的红，我照 v1c＋编排者 §151 的读数、没自己跑。** 我只读了 `composer_test.go:48` 的比对机制（Go emit 键集 ↔ 前端 interface 双向减集）。若前端 agent 在 10-02 之后已把那两枚键声明进 interface，AC#2 现在可能已可翻——**我⛔读不了 `frontend/**` 证实或证伪**。⇒ 下一枚能跑测试的腿：`go test ./internal/panel -run TestComposerContractTypesMatchFrontend` 一发即可定。
2. **"TierOf 零生产调用者"是 HEAD 快照，工作树正在变。** 未提交的 `config_readers_255.go:174` 是它的第一个生产调用者、正被 255-r2 写。若 255-r2 在我交件前落库，AC#8/AC#10 的 Go 半句可能已就绪、"零调用者"这句随之过期。⇒ 本腿的读数只对"我这几发 `git grep HEAD` 的时刻"负责。
3. **AC#4 我判"开窗链已生产接线"是顺着调用点＋注释读出来的推断，没端到端确认那枚回调真被热键/托盘触发。** 我核到 `resident_windows.go:163 withPanelHost(...RequestToggle...)` 把回调交给了 `startResidentBall`、注释说热键/托盘驱动，但**没逐行追到热键注册处确会调那枚回调**。可能形状：接了但回调挂在没人发的 event 上（写了没接的另一种）。⇒ 需要跑／真机的那腿确认；本腿只敢说"代码路径连通、非零调用"。
4. **行号会漂，且我引的那几包正被写。** `manager.go` 的 llm 热应用项我从 v1c 的 `:281` 更正为现量 `:284`（`:281` 现为 `hotkey`）；但 `cmd/wisp`＋`internal/config` 此刻有写腿在飞，我引用的若含工作树行号（尤其未提交的 255-r2 文件）可能随其写入而变。我尽量用 `git show HEAD:`／`git grep HEAD` 取committed 形状，未committed 的 255-r2 两处我**具名标了"未提交、属 255-r2"**。
5. **我把 AC#8 从编排者的"两栖"标签里摘出来、改判"等 255-r2 而非前端"。** 依据是 §126"归 255-r2 落地…等 255-r2 交完由非实现者一并裁"这句把翻勾条件钉在 255-r2＋裁。但"两栖"一词也许本意含"页面可见"要前端——我承认 AC#8 带一枚前端传输尾巴（§4-3），只主张那尾巴不是它的翻勾前置。这是我读法，非定论。
6. **AC#9"只剩一枚 `.gitkeep`"是 `git ls-files` 的名册级读数＝"内容不在 git 跟踪面"，不等于"内容不在磁盘"。** `//go:embed all:dist` 读的是构建时的**真实文件系统目录**，不是 git；若前端 agent 本机 build 出过 dist 内容但没 commit（或被 gitignore），磁盘上可能有包、而我这把尺看不见。⇒ 我⛔不 `ls frontend/dist`（禁令层）、也跑不了 embed——"真有没有内容"这发只能由能碰那层或能跑 embed 的腿/owner 定。**这一条足以推翻"AC#9 判不了"里的"判不了"三字之外的半句**，单列。
7. **`EffectiveNow/EffectiveNextTask` "零非测试写者"我只扫了 `= panel.EffectiveNow` 这类赋值形。** 若 255-r2 用别的写法（如 map 赋值、具名函数返回）给 `res.Tier` 赋 now/next_task，我这把 grep 尺量不到。方向上仍支持"HEAD 只活 restart 一枚"，但非穷尽。
8. **凭据具名字段那套（`provider_credential` 只写不回显）我没逐个 setter 追它的值流向**，只采信 v1c §AC#0/AC#3 的"值不走共用封套、只上 `api_key_ref`"名册级判语（已勾格）。本腿复核集中在未勾的 5 格，勾过的没重审。

## §6 量不到的地方（具名）

（取数中）

## §7 交件判语

（取数中）
