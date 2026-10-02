# e2e-panel-1 — 「从双击到看见设置页并存下去」跳表（只读取证腿）

> 本件唯一问题：机主那句「我要能看到主面板，我要点击设置，自己配置模型这些参数」到今天之间**差几跳**。
> 本腿**不裁任何 AC、不判任何票好还坏、不提任何修法**。逐跳四问：这一跳是什么／今天通没通／承重处 `file:line`／不通缺哪一环、归哪张票。
> 写点只有本文件。产码／测试／票面／`docs/**` 零字节改动。禁令那两棵树的**内部具体路径与行号本件一律不出现**（本件只把它们作为"禁令名"提及，与票 255／票 253 票面同口径）；页面内容不是本腿射程。

## §0 起手锚（进场第一发，逐字）

```
$ date -Iseconds
2026-10-02T11:34:18+08:00

$ git log -1 --format='%H %ad %s'
21a15d2704d72e7469423b72aaf2741c899614c4 Fri Oct 2 11:33:27 2026 +0800 ledger(A529)＋票 254 结案改名 -done：…（完整标题见 git log）

$ git rev-parse --abbrev-ref HEAD
dev

$ git status --porcelain -- cmd internal | wc -l
0
```

- 代号复认：`.scratch/wisp/probes/` 起手 listing 共 98 枚，无 `e2e-panel-1`；本目录由本腿新建（⛔ 只建不删）。
- ⚠ 共享工作树整树非空是常态，本腿只在 `cmd internal` 两枚根上取分母；`0` 是**起手时刻**读数，不是恒量。
- 中途复量（写 §4/§5 之前）：`date -Iseconds` = `2026-10-02T12:11:34+08:00`，`git status --porcelain -- cmd internal | head -20` **仍为空**。
- ⚠ `248-v1c` 正在往跟踪文件种临时变异 ⇒ 本件的 `file:line` 一律**双口径**给出：工作树读数 + `git cat-file blob HEAD:<path>` 读数。两把尺已对 `panel_config_store.go`／`panel_host_windows.go`／`manager.go` 三枚承重件逐条复认（见 §6）。
- **HEAD 在本腿写作期间移动过两发**，第二发复量：`date -Iseconds` = `2026-10-02T12:46:21+08:00`，`git status --porcelain -- cmd internal` 行数 = **0**；`fb1d9c19`（ledger A530＋票 248 翻格＋新立票 256）与 `e7521610`（248-r2 骨架）两发的 `--stat` 我逐枚读过，**只动 `.scratch/wisp/issues/**` 与 `docs/reports/**`，零产码** ⇒ 本件全部 `file:line` 在新 HEAD 上**逐条重量过一遍，一字未变**（`res.Tier = panel.EffectiveRestart` 仍 `:228`／`:275`，`func NewPanelManager(` 仍 `:175`，`Width: 420,`／`Height: 260,` 仍 `:304`／`:305`，`{"llm", ...}` 仍 `:281`）。
- 同两发里有一条与 §6 第 1 条同源、由编排者自己在 `fb1d9c19` 的正文里复述过的裁定：**AC#8 那句"面板对 `[llm]` 说重启才生效、`manager.go:281` 却把它列在热应用表里"归 `255-r2`**（本腿 §1 的 H13 与 §3 同格归口与之一致，非本腿新裁）。

## §1 跳表（本腿正身）

> 尺＝**调用点枚数**（谁产出／谁投递／谁落盘三处各有没有栈），⛔ 不是 `go build` 通没通。
> 「通」＝**非测试调用者 ≥1 枚**；「接线通」＝栈齐但结论只有真跑才能定；「不通」＝某一处栈没有。
> 依赖顺序排（H1 在前 H14 在后），⛔ 不含修法、⛔ 不裁 AC。机主那句话的原话射程＝**H1→H6→H9→H10→H13**。

| 跳 | 这一跳是什么 | 今天 | 承重处（工作树＝HEAD `21a15d2`，已双口径） | 不通的话缺哪一环、归哪张票 |
|---|---|---|---|---|
| H1 | 双击 `wisp.exe`（无 argv）落进常驻进程 | **通** | `cmd/wisp/main.go:58-67`（`len(args)==0` → `attachParentConsole()` → `runResident()`）；`runResident` 在 `cmd/wisp/resident_windows.go`（windows 构建） | — |
| H2 | 常驻进程手里**有** `config.toml`（面板装配的前置） | **不通** | 唯一产码建件调用点＝`cmd/wisp/run.go:245 ensureFirstRunConfig`（在 `runTextTask`，即 `wisp run`）；常驻面板链先读盘＝`cmd/wisp/panel_inbound.go:229-233`（`config.NewManager`），缺文件在 `internal/config/loader.go:64-68` 直接返 `ClassConfig` 错 ⇒ `cmd/wisp/panel_resident_windows.go:189-191` 返错 ⇒ `cmd/wisp/resident_windows.go:142-146` 打印「panel hot key and the tray item will be recorded, not executed」 | 缺「常驻/双击那条入口上有人建首份配置」这一环。归**票 198**（其票面 `cmd/wisp/firstrun.go:16-29` 自陈建件只挂 `wisp run`，并点名两枚测试钉禁止常驻腿建文件） |
| H3 | 面板宿主对象被**真造出来**（`NewPanelManager` 有产码调用者） | **通**（前置在 H2） | 定义 `cmd/wisp/panel_host_windows.go:175`（全仓唯一定义，复认成立）；产码调用点 **1 枚**＝`cmd/wisp/panel_resident_windows.go:204`，其建造者 `newResidentPanelManager` 被 `cmd/wisp/resident_windows.go:142` 调，线程由 `:148 startResidentPanel` 起 | — |
| H4 | 机主的**手势**→面板出现 | **部分通** | 通的两形：托盘「打开面板」`internal/ball/tray_windows.go:86` → `internal/ball/ball_windows.go:672/677` → `cmd/wisp/resident_ball_windows.go:179`；全局热键（默认 `Ctrl+Alt+P`，`internal/ball/hotkey_windows.go:69`）→ `internal/ball/ball_windows.go:665` → `cmd/wisp/resident_ball_windows.go:178`。两形都汇到 `panel_resident_windows.go:378 RequestToggle`→`:329 RequestShow`→`panel_host_windows.go:395 Show`→`:230 bringUp` | 不通的那形＝**点球**：`cmd/wisp/resident_ball_windows.go:174 OnClickBall` 只 `recordBallGesture("click")`，不接面板；事件表 `internal/ball/ball_windows.go:50-59` **没有双击槽位**（`grep -i "doubleclick\|DBLCLK"` 产码零命中）。⇒「双击/点击就能看到主面板」这句话缺**手势那一环**。归**票 33**（AC#1 那格接进常驻时只覆盖热键＋托盘两形）；⚠ 双击本身**票面没写** ⇒ 按未定义即停在件里具名上报，不替它填射程 |
| H5 | 窗口里出现的是**页面**，不是探测页也不是「未建」告示页 | **不通（干净检出必不通）** | 时序已翻：`cmd/wisp/panel_host_windows.go:340` 先 `firstRoundTripLocked`、`:342` 后 `serveEntry`；`serveEntry:368-371` 要求 `assets.Built()`；`Built` 由 embed 树里能否 `Stat("index.html")` 决定（`internal/panel/assets.go:54-60`）；失败落 `:343 serveNotBuiltNoticeLocked`（`:356-365` 自造告示页）。仓库级 census：embed 目标目录**跟踪件＝1 枚**（派单锚 3 的那枚占位件），**ignored 件＝6249 枚** | 缺「页面产物进树」这一环。归**票 33 AC#13**（次序＋"会响的最终文档断言"两半）＋产物那侧（⛔ 界面侧不在本腿射程，具名归口见票 253 末段那句「由 owner 带话」） |
| H6 | 页面上有「设置」这条路、点了发得出**认得的方法名** | Go 侧**通**，页面侧**判不了** | 两枚名进白名单：`internal/panel/bridge.go:66-67`（`config.get`／`config.set`），收录于 `knownComposerMethod`（`bridge.go:146-152`）；选择器字段在 `bridge.go:117-120`；封套三门禁（方法名/来源/`requestId`）＝`bridge.go:126-144` | 页面里有没有「设置」入口与控件：`frontend/**`／`design/**` 既不读也不引 ⇒ **判不了**（§5 第 6 条）。归**票 248**（票面标题那一跳的界面半）。⚠ 旧封套形状（不走绑定）时 Go 侧零现场那一支归**票 253 AC#3** |
| H7 | 页面发的封套被 Go **收到并按方法分流** | **通** | 产码投递栈齐：绑定＝`cmd/wisp/panel_host_windows.go:319-322`（名字在本机常量化 `:82`）→ `dispatchRaw:544-552` → `internal/panel/composer_dispatch.go:153-165 Handle` → `:197-206` 两个 `case` → `d.Config.HandleConfigRequest`（产码调用点 2 枚）；第二条栈＝控制台腿 `cmd/wisp/panel_inbound.go:163`（`main.go:115-120` 派发的 `wisp panel-inbound`） | — |
| H8 | 「自己配模型这些参数」能被**点名**（字段名册＋校验前置） | **通** | 名册 7 枚＝`internal/panel/config_handlers.go:57-69`，选择器表 `:74-94`，`WritableFields():99-103`；逐条闸门在 `write():307-365`（锁定族 `:315`、未列名 `:321`、必指服务商 `:328`、必指模型 `:334`、空值 `:349`）；锁定族清单 `:108`（`risk.`／`fs.`／`net.`／`plugins.`／`privacy.`／`models.`／`audio.`＋`permission_mode`／`allowed_dirs` 两枚叶子名） | — |
| H9 | 读侧答得出「配了没配」 | **通（形状是一句人话，不是字段集）** | 产出＝`cmd/wisp/panel_config_store.go:88-133 ReadSettings`（经 `mgr.Config()` 读 `[llm]`，逐服务商问 `refRecorded:152-167`）；聚合口径 `aggregateCredential:137-146`；投递＝`internal/panel/config_handlers.go:292-303` → `renderSettingsView:445-462` → H7 的回话 | 形状那一格我不裁（§4 第 6 条）；若判据要结构化字段，那是 C17 契约面，⛔ 未定义即停（票 255 禁区段自己点名 `SettingsView`／`renderSettingsView`／`Snapshot` 三处属这一面） |
| H10 | 写侧**落盘**（一键一路径、验后写、写后按文件报） | **通** | 落盘栈齐：`configStore.ApplySetting`（`panel_config_store.go:173-230`，6 枚字段各一枚 setter 调用 `:184/194/201/208/215/217`）→ `internal/config/settings.go:199-249 writeOneKey`（validate 前置 `:230`、内存+文件 `:237-246`、失败回滚内存 `:243`）→ `writeguard.go:111/160 mergeWrite→SaveFile`（temp+rename）；凭据那支＝`panel_config_store.go:238-279 StoreCredential`（`secret.Store` 后写引用 `:267`，回话只回引用 `:276-277`）；回执键路径按**写后文件 diff** 产（`writtenKeyPaths:254-268`） | — |
| H11 | 干净机器上这些字段**其实写哪都写不进去**（选择器前置对面） | **不通** | 首建默认表里服务商注册表是**空的**：`internal/config/defaults.go:77-78`（reflect 走默认时对 `reflect.Map` 一律「leave nil」）＋`internal/config/schema.go:419`（`Providers map[string]Provider` 无 `default` 标签）；六枚 provider/model 级字段都要求**已存在的行**（`settings.go:310-314 unknownProviderErr`、`:294-308 requireCatalogEntry`）；`role_chat_model` 落 `catalog.go:97-101`（model 有、provider 空 ⇒ 拒），而名册里**没有**写 provider 的那枚字段（`config_handlers.go:57-69`）；`firstrun.go:93-96` 自陈「未替你选任何模型」 | 缺「首份配置里就有一行可被选择的服务商/模型」或「名册里有一条创建路径」这一环（**只列缺口，不列修法**）。归**票 248**（名册范围）＋**票 198**（首建内容那一半） |
| H12 | 回执**真到达页面**（点了要看得见） | **接线通、运行时判不了** | 回话＝`Handle` 的返回值原样给绑定（`panel_host_windows.go:320-321`），文案由 `renderSettingReceipt:381-396`／`renderCredentialReceipt:401-424`／`tierSentence:428-441` 产；投递依赖库的 `Run()` 泵（`cmd/wisp/panel_resident_windows.go:250`；该文件 `:21-27` 自陈「队列唯一读者是 `Run()`」——这句是注释自陈，本腿没读第三方模块复认） | 「到达」只有真页面 `await` 才能定 ⇒ §5 第 3 条。归**票 33 AC#14**（那一格判据写的就是"回话真到达页面"） |
| H13 | 存下去之后**真的生效**（机主那句话的落点） | **不通** | 生效说法相反：面板那条写路径对任何写入恒答 restart＝`panel_config_store.go:228`（`res.Tier = panel.EffectiveRestart`；`switch w.Field` 在 `:182-223` 只分派 setter、不分派档位）＋凭据那支 `:275` 同一枚常量；热应用段表把 `llm` 列在 hot＝`internal/config/manager.go:281`（段表注释 `:271`「Everything else is hot-tier」）；两套之间**零行码互相翻译**。宿主拿不到配置对象＝`NewPanelManager` 签名无配置参数（`:175`），窗口尺寸唯一来源是 `:304-305` 写死的 `420`／`260`，而 `[panel]` 五枚键（`internal/config/schema.go:529-539`，`width` 的 `default:"640"`）非测试读者 **0 枚**（`grep "cfg.Panel\.\|\.Panel\.Width"` 产码零命中）。tick 装配只在 `cmd/wisp/run.go:813 rt.startConfigReload()`（属 `assembleRuntime`，其内 `config_reload.go:153 CheckAndReload`、`:171` 那句「这些段已立即生效」只写 `rt.stdout`）；面板链那枚 Manager 从不被 tick；常驻任务腿无控制台时 `resident_task_source_windows.go:224-231` 直接 `return nil` | 缺「按键产出档位」＋「把配置交进宿主」＋「面板写入与管理器段表同一口径」三环。归**票 255**（AC#1 那句误报／AC#4 真生效／AC#5 两套相反）；`OnReload` 出货进程没听众那一格票面明写**归票 42、不许并票** |
| H14 | 面板窗口能显示「当前模型／凭据状态」这类快照维度 | **不通（面板侧无泵）** | 快照泵非测试构造点 **1 枚**＝`cmd/wisp/run.go:699 panel.NewSnapshotPump`，`Credential` 读者挂 `run.go:710 rt.settings.credentialStatus`（读者本体 `panel_config_store.go:285-291`）；常驻面板链（`newComposerDispatchChain`）里没有泵、也没有 Go→页推送调用者（`grep "\.Eval(\|ExecuteScript"` 产码 **0 命中**）；`cmd/wisp/panel_pump.go:8-20` 自陈"字节记进 run 的台账、最后一公里仍开" | 缺「面板宿主这一侧的泵装配与出向投递」这一环。归**票 145**（载体字段那一格）；⚠ 泵/出向那半在码内注释里被命名为票 35 的射程（`run.go:691`、`panel_inbound.go:56-57`），**本腿没读那张票面** ⇒ 具名上报、不裁归口（§5 第 9 条） |

## §2 今天已经通的跳（每跳一把调用点尺＋真实读数）

> 读数口径统一：`grep -rn <符号> cmd internal tools --include=*.go | grep -v _test.go` ⇒ 命中数拆成「定义 1 枚＋产码调用点 N 枚」。
> 下面每一枚 N 都是本腿自己现跑的，⛔ 没有一条是"应该有"。时刻 `12:11:34` 起（HEAD `21a15d2`），全部行号在 `12:46:21` 的新 HEAD 上逐枚重量过，一字未变（见 §0 第三发）。

| 跳 | 尺（问的那一句） | 现量读数（时刻 `12:11:34` 起，全部 HEAD `21a15d2`，`cmd internal` 分母 0） |
|---|---|---|
| H1 | 无 argv 分支有没有把常驻叫起来 | `main.go:60-67`：`len(args)==0` → `runResident()`，`return`。产码路径 **1 枚**，且 `usage` 块（`main.go:24-31`）逐字写着「GUI resident process … hosts the floating ball window, its tray icon and its four global hot keys」 |
| H3 | `NewPanelManager` 全仓非测试调用点几枚 | 定义 1（`panel_host_windows.go:175`）＋**产码调用点 1**（`panel_resident_windows.go:204`）；它的建造者 `newResidentPanelManager` 非测试调用点 1（`resident_windows.go:142`）；线程启动 `startResidentPanel` 非测试调用点 1（`resident_windows.go:148`）。⇒ 台件里那枚「0 枚产码构造点」的具名 skip（`panel_host_gate_test.go:177-189`）走的是**另一支**（`ctorHits==0` 才跳） |
| H4（两形） | 手势到 `Show` 之间**每一环**有没有产码调用者（逐环尺，不是一枚总数） | `requestPanelOpen(` 命中 3＝定义 1（`resident_ball_windows.go:88`）＋产码调用点 **2**（`:178` 热键、`:179` 托盘）；`RequestToggle(` 命中 2＝定义 1（`panel_resident_windows.go:378`）＋产码调用点 **1**（`resident_windows.go:164`，经 `withPanelHost` 注入）；`RequestShow(` 命中 2＝定义 1（`:329`）＋产码调用点 **1**（`:387`，`RequestToggle` 的未显示分支回落）；`showOnThread(` 命中 2＝定义 1（`:359`）＋调用点 1（`:334` 闭包）；`rp.mgr.Show(` 调用点 1（`:360`）；球侧事件源 2 枚（`internal/ball/ball_windows.go:665` 热键、`:672/677` 托盘）；托盘那枚菜单项**真在菜单里**（`internal/ball/tray_windows.go:86`） |
| H7 | 谁产出封套／谁投递／谁分流 | 产出＝页面绑定 `panel_host_windows.go:319-322`（绑定名常量 `:82`）＋控制台腿 `panel_inbound.go:163`；投递＝`dispatchRaw:544-552` → `composer_dispatch.go:153-165`；分流＝`:197-206` 两个 `case`，`HandleConfigRequest` 产码调用点 **2 枚**。⚠ 这一格三处栈**都在**——不是只有"有函数" |
| H8 | 名册是否被"活的调用点"读，而不是只被定义 | `knownWritableField`（`config_handlers.go:368-375`）在 `write():321` 被调；`fieldNeedsProvider/Model/TakesValue` 三张表在 `:328/334/349` 被读；`lockedFieldFamilies:108` 在 `refusedLockedFamily:115-136` 被读，后者在 `:315` 被调 ⇒ 四道闸门各 1 个产码调用者 |
| H9 | 读侧：谁产出／谁投递／谁落可见面 | 产出＝`configStore.ReadSettings`（`panel_config_store.go:88`），非测试调用者 **2 枚**：`config_handlers.go:293`（走页面回话）＋`panel_config_store.go:286`（快照读者 `credentialStatus`）；配置来源＝`mgr.Config()`（`:92`）；可见面＝`renderSettingsView`（`config_handlers.go:302` 调用） |
| H10 | 写侧：谁产出／谁投递／谁落盘 | 产出＝`ApplySetting`/`StoreCredential` 各 **1 枚**非测试调用者（`config_handlers.go:355`／`:340`）；投递到 setter＝`panel_config_store.go` 里 `s.mgr.Set*` **7 枚调用点**（`:184/194/201/208/215/217` ＋ 凭据那支 `:267`）；落盘＝`settings.go:242 mergeWrite` → `writeguard.go:160 SaveFile`；写后回执＝`:248 writtenKeyPaths`（再读文件 diff，`:261-268`）。内存同写＝`:241 apply(m.cur)`，失败回滚＝`:243` |
| H12（仅接线） | 回话的产与投有没有栈 | 产＝`config_handlers.go:381-424` 两枚 receipt 渲染器，调用点在 `:346`/`:364`；投＝`panel_host_windows.go:320-321` 把 `Handle` 返回值原样 `return` 给绑定。**投递的运行时态未量**（§5 第 3 条），这一行只算"接线通" |

## §3 今天不通的跳（缺哪一环＋归哪张票；⛔ 不含修法）

按依赖顺序，逐跳只写两件事：**缺的那一环在哪**、**票面上哪张票的名字覆盖它**。修法归编排者与后续腿。

- **H2 首启建 `config.toml`（双击这条入口上没人建）**
  缺的那一环：建件调用点在产码里只有 1 枚，挂在 `wisp run`（`cmd/wisp/run.go:245`）；常驻/双击那条路在**同一条链的最前面一步**就因为缺文件返错（`cmd/wisp/panel_inbound.go:230` → `internal/config/loader.go:64-68`），且这枚返错发生在 `secret.NewStore`（`:262`）**之前**，所以连"数据根顺手建出来"都没有。归**票 198**（`cmd/wisp/firstrun.go:16-29` 自陈只挂 `wisp run` 那一枚产码调用链，并具名两枚测试钉禁止常驻腿建件）。
  后果链上还有两枚下游：H3 的前置、H4 的两形手势（此时只剩"recorded, not executed"，`cmd/wisp/resident_windows.go:146`）。

- **H4 机主那句"点击"那一形**
  缺的那一环：球身点击事件源不接面板（`cmd/wisp/resident_ball_windows.go:174` 只记手势），且球侧事件表（`internal/ball/ball_windows.go:50-59`）里没有双击这一枚。归**票 33**（常驻接入那一格只落热键＋托盘两形）。⚠ **双击该不该算票 33 射程，票面没覆盖 ⇒ 本腿停下来具名上报，不自行判定**（AGENTS.md §2「未定义即停」第 1 条口径）。
  同形另一枚已量到的事实（不是修法，是射程外的账）：常驻那条球腿的热键是**硬编码默认表**（`cmd/wisp/resident_ball_windows.go:171 Hotkeys: ball.DefaultHotkeys()`），把 `[hotkey]` 段接进来的那枚 reloader 在非测试码里只有 1 枚调用者且它在调试件里（`cmd/balldebug/main.go:237 NewHotkeyReloader`）⇒ 机主就算把面板拿到了，改 `[hotkey]` 换那枚面板键位在常驻进程里也没有人执行；这一格本腿只报形状，归口不明（票 255 的 AC#3 段外名册那一族？票面没写这一枚 ⇒ 具名上报，不替它归口）。

- **H5 窗口里是页面而不是告示页**
  缺的那一环：`serveEntry` 的前置 `assets.Built()` 由 embed 树里有没有入口文件决定（`internal/panel/assets.go:54-60`），而仓库级 census 给的是跟踪件 1 枚（占位件）⇒ 干净检出这一档**必然**落到 `panel_host_windows.go:343` 的告示页。归**票 33 AC#13**（次序那一半码上已翻，"会响的最终文档断言"那一半本腿不裁）＋产物那侧（⛔ 界面侧不读不引，具名归口按票 253 末段那句"由 owner 带话"）。
  ⚠ 本机工作树另有一枚 ignored 件计数（6249），它与"今天这棵树里 `Built()` 究竟是哪一档"不是一回事，本腿判不了，落在 §5 第 2 条。

- **H6 页面里的"设置"那条路**
  Go 侧那半是通的（两枚名进白名单），缺的是**界面那一半**：设置入口/控件/路由是否存在本腿不读不引。归**票 248**（票面标题「the panel has no settings route」那一跳）。另有一支独立缺口：页面若按旧封套形状直接发，Go 侧**零现场**——归**票 253 AC#3**（本腿只确认 Go 侧的拒绝都会 `record`：`composer_dispatch.go:157`/`:220`、`config_handlers.go:495-505`，够不着"没进 Go 的那一发"）。

- **H11 干净机器上七枚字段全部走不通**
  缺的那一环：六枚 provider/model 级字段要求配置里**已有**那一行（`internal/config/settings.go:310-314`、`:294-308`），而首建默认表把服务商注册表留成 nil（`internal/config/defaults.go:77-78` ＋ `internal/config/schema.go:419` 无 `default` 标签）；剩下的 `role_chat_model` 一枚会撞 `internal/config/catalog.go:97-101`（model 有 provider 空 ⇒ 拒），而名册里没有可写 provider 的字段（`internal/panel/config_handlers.go:57-69`）。归**票 248**（名册范围）＋**票 198**（首建内容那一半，`cmd/wisp/firstrun.go:93-96` 自陈「未替你选任何模型」）。
  ⚠ 这一格是"机主自己配参数"最硬的那一枚死结：不是校验太严，是**没有一行可被选**。本腿只报形状，不提该往哪边松。

- **H13 存下去之后生效**
  缺三环，逐枚列：① 档位不按键产出（`cmd/wisp/panel_config_store.go:228` 与 `:275` 对任何写入恒 `EffectiveRestart`，`switch w.Field`（`:182-223`）只分派 setter 不分派档位），与热应用段表（`internal/config/manager.go:281` 把 `llm` 列 hot）答相反，**两套之间零行码互相翻译**；② 宿主拿不到配置对象（`panel_host_windows.go:175` 签名无配置参数，尺寸唯一来源 `:304-305`，而 `[panel]` 五枚键非测试读者 0 枚，`internal/config/schema.go:529-539` 写着 `width default:"640"`）；③ 面板链那枚 Manager 不被 tick（tick 只在 `cmd/wisp/run.go:813`，句"这些段已立即生效"只写 stdout：`cmd/wisp/config_reload.go:171`）。归**票 255**（AC#1／AC#4／AC#5 三格各对一环）；`OnReload` 出货进程无听众那一格票面明写**归票 42、不许并票**。

- **H14 面板侧没有快照泵**
  缺的那一环：泵的非测试构造点唯一挂在 `cmd/wisp/run.go:699`，常驻面板链里没有；Go→页推送在产码里零调用者（`.Eval(`／`ExecuteScript` 产码 0 命中）。归**票 145**（载体字段那一格）；泵/出向那半的归口在码内注释里被命名为票 35，**本腿没读那张票面** ⇒ 具名上报、不裁（§5 第 9 条）。

## §4 我可能判错的条目

1. **「通」的尺可能偏宽**。本件判「通」＝**非测试调用者 ≥1 枚**。但面板那条链上多枚承重件带 `windows` 构建标签（`cmd/wisp/panel_host_windows.go:1`、`cmd/wisp/panel_resident_windows.go:1`、`cmd/wisp/resident_windows.go`），非 Windows 构建图里它们**不存在**。我按 D7（平台范围＝Windows）把这族判成「通」，如果编排者要的尺是「跨平台可编译可跑」，我这几格一律要降级。
2. **行号可能在我读数之后漂移**。起手 `cmd internal` 分母为 0，12:11:34 复量仍为 0；但 `248-v1c` 的种变异窗口不落 `git status` 之外的时间。我只对三枚文件做了 HEAD-blob 双口径复认，其余（`bridge.go`／`config_handlers.go`／`settings.go`／`resident_ball_windows.go`／`main.go`／`loader.go`／`manager.go` 的段表）**只有工作树单口径**。
3. **`res.Tier` 那枚行号我读成 `:228`，派单与票 255 的 AC#5 都写 `:227`**。两把尺（工作树 grep 与 HEAD blob grep）**一致给 228**（且 `:227` 是 `res.Written = written`）。我倾向认为不是漂移而是票面/派单少一行，但也可能是我对「`:222-229` 前面没有任何按键/按段判断」这句话的射程理解偏窄——那一段里确实有一枚 `switch w.Field`（按键分派），只是它不分派**档位**。这条留给编排者定文字，我不改票面一字。
4. **「机主能存下去」的判据我按落盘取**。`writeOneKey` 同时改内存 `m.cur` 与文件（`internal/config/settings.go:237-246`），我据此判「落盘＝通」。若编排者的尺是「下一次任务真的用上这个值」，那一跳我今天只在 §3 里列成跳，没有判它通。
5. **6249 这个数字我怀疑但不敢用**。我用仓库级 census 量到 embed 目标目录下 **跟踪件 1 枚**（复认派单锚 3），同一目录口径下 **ignored 件 6249 枚**。我只数了路径尾带 `/dist/` 的条目，**没有**（也不能）点名那枚目录，因此**无法排除**其中混有第三方的 `dist/`。这条直接决定 §1 里「页面产物」那一跳是「本机已备、干净检出仍缺」还是「本机也缺」——见 §5 第 2 条。
6. **`config.get` 的回话形状我判成「一句人话，不是字段集」**（`internal/panel/config_handlers.go:445-462` 产 `renderSettingsView` 字符串，宿主把 `Handle` 的返回值原样交给绑定回话）。如果机主那句「看到设置页、自己配参数」**不要求**结构化字段（页侧把这句话渲染出来就算满足），我把它列成跳是**过度解读**；反过来如果我列轻了，那是把 C17 契约面漏了。这一格我不裁，只标尺不确定。
7. **常驻进程里有几枚 `config.Manager` 我判不动**。尺上 `config.NewManager` 非测试调用点是 `cmd/wisp/panel_inbound.go:230`（常驻面板链经 `newComposerDispatchChain` 走它）与 `cmd/wisp/run.go:409`；但 `run.go:409` 所在的 `assembleRuntime` 只在 `interactiveStdin()` 非 nil 时才被常驻腿调到（`cmd/wisp/resident_task_source_windows.go:218-231` 那一支 `return nil`）。⇒ **纯双击（无控制台）时常驻进程可能只有 1 枚 Manager**，这时「两枚各说各话」不成立，成立的是「任务管线根本没装配」。我没跑，判不了，列进 §5 第 4 条。
8. **票面 `-done` 状态我只看了文件名**。本腿只 `ls` 出六枚票（33/145/198/248/253/255）都**不带** `-done` 后缀，池内另有 84 枚带后缀。`-done` 是防重领键、不等于验收结论，我不拿它当「好/坏」判据，也**没有**逐枚读六张票的 AC 勾选态（33 的票面我只读了 AC#13/#14 与裁定段，见 §5 第 5 条）。
9. **锚 1 里「窗口尺寸写死在 `:304-306`」**：我复认 `:304 Width: 420,`／`:305 Height: 260,`，`:306` 是 `},`。射程一致，只是分界行差一枚；我在跳表里按 `:304-305` 写。

## §5 判不动／没测到的地方

1. ⛔ **本腿零跑**：`go build`／`go vet`／`go test`／`./...`／`go list -deps` 一律未执行（`248-v1c` 正在整包跑 `cmd/wisp`）。⇒ 本件所有「通/不通」都是**调用点枚数＋码内分支阅读**的读数，**不含任何执行证据**；凡是「只有真跑才能定」的格子，我在 §1 的「通」列里写的是**接线态**，不是**运行时态**。
2. **页面产物在运行时到底是哪一档，我测不到**。可读的码侧事实：`internal/panel/assets.go:54-60` 用 `fs.Stat(tree, "index.html")` 决定 `built`；`:75-78` 未 built 时 `Resolve` 直接 `errNotBuilt`；`cmd/wisp/panel_host_windows.go:342-344` 供页失败就走 `serveNotBuiltNoticeLocked`（`:356-365` 那枚自造告示页）。至于**本机今天**embed 里到底有没有入口字节，需要跑一次 `panel.BuiltinAssets()` 或读那枚目录的文件名——前者被禁、后者属禁区。⚠ 现存两枚**互相冲突的台件读数**（都在 `cmd/wisp` 与 `docs/evidence/s1/` 里，不是我推的）：一枚记录工作树形如「有页面包」，另一枚的具名 skip 文案写着「embed 解析不出入口，本节在本树无对象」（`cmd/wisp/panel_resident_windows_test.go:317`）。我**不裁**哪枚对。
3. **「回执真到不到页面」测不到**。码侧接线在：`cmd/wisp/panel_host_windows.go:319-322` 绑定回话＝`dispatchRaw` 的返回串；`cmd/wisp/panel_resident_windows.go:250` 把线程泵交进库的 `Run()`（`Run()` 是派发队列唯一读者——这句是 `:21-27` 注释自陈，我没读第三方模块源码复认）。判「到达」要一发真页面 `await`，本腿无此射程。
4. **常驻进程运行态的对象数量判不了**（§4 第 7 条）：面板写入用的是链上那枚 Manager，热加载 tick 只在 `cmd/wisp/run.go:813 rt.startConfigReload()` 装（属 `assembleRuntime`），而 tick 报的那句「这些段已立即生效」只写 `rt.stdout`（`cmd/wisp/config_reload.go:171`）。⇒「谁在什么进程里读到面板刚写的值」我没跑，判不了。
5. **票 33 的 AC 勾选态与裁定 2/3/4 的落地程度我没读全**（票面 342 行，我只读了 AC#13/#14 与随后裁定段）。凡我把某一跳「归票 33」的地方，归的是**那一格的名字与射程**，不是「这格还没做完」的判语。
6. **界面侧我一律不量**：设置页存在与否、页面上有没有输入框、`wispDispatch` 有没有被页面调、页面里那批方法名字面量——`frontend/**`／`design/**` 禁令覆盖，具名口径按票 253「需要界面侧配合的那一件由 owner 带话」执行，本腿不转达、不推断。
7. **凭据面我只写字段名/常量名**：`credentialValue`（`cmd/wisp/panel_config_store.go:53`）、`FieldProviderCredential`、`provider_api_key_ref`、引用形状 `dpapi:`/`env:` 前缀名。⛔ 零凭据值，零本机 blob 名。
8. **`d22scan` 会不会因这棵树今天的形状判红，我判不了**：`tools/d22scan/runtests.sh:98-102` 的尺我读到了（`skipped != 0` ⇒ `exit 1`，「SKIP 不是过」），但**它今天到底跳过哪几枚**只有跑才知道。派单锚 3 那句「AC#13 因此必跳」我**只能确认尺的形状**，不能确认它今天真跳。
9. **票 145（快照字段）在本跳表里的射程**：我量到快照泵的非测试构造点只有 `cmd/wisp/run.go:699`（`panel.NewSnapshotPump`），常驻面板链**没有**泵；`Credential` 那枚读者挂在 `cmd/wisp/run.go:710 rt.settings.credentialStatus`。⇒「面板窗口能收到快照」这一跳我列进 §3，但**票 145 的题面是载体字段**，泵装配的归口我判不动，标在这里。
10. **未定义即停（本腿没有自填的格）**：我没有碰到需要新造规矩的情况；遇到「两套说法相反」「两枚读数冲突」一律**上报不裁**（§4 第 3 条、§5 第 2 条、§6）。

## §6 我推翻编排者哪一句

> 只碰派单那四条锚的复认结果，⛔ 不引申、不裁别家的票。

1. **锚 2 的那枚行号差一行（唯一实质复认不上的点）**。
   派单原话：「`cmd/wisp/panel_config_store.go:227` 逐字 `res.Tier = panel.EffectiveRestart`」。
   我的两把尺**都给 `:228`**：工作树 `grep -n "res.Tier = panel.EffectiveRestart"` ⇒ `228` 与 `275`；`git cat-file blob HEAD:cmd/wisp/panel_config_store.go` 同口径 ⇒ 也是 `228`／`275`。`:227` 那行逐字是 `res.Written = written`。
   同一枚差一行出现在**票 255 的 AC#5 现量段**（票面写的也是 `:227`），出处具名：`docs/reports` 不是我射程，我只说这一处**票面文字与本腿读数差一行**，⛔ 不改票面一字（改契约＝人工批准）。
   「`:222-229` 前面没有任何按键/按段判断」这半句我**判成成立但措辞可争**：那一段里有一枚 `switch w.Field`（`:182-223`），它按**键**分派 setter、不按键分派**档位**，`:224-229` 之间对档位的写法是无条件（`res.Tier = EffectiveRestart` 在 `:228`，`res.Written` 在 `:227`）。如果票面那句的射程是"没有按档位的判断"，我复认；如果是"那一段里没有 switch"，我不复认（switch 在同函数 `:182` 起）。

2. **锚 1 全部复认上**，只有分界行差一枚：`NewPanelManager` 全仓唯一定义＝`:175`，签名逐字 `(disp *panel.ComposerDispatch, assets *panel.Assets, dataPath string)`，**无配置参数**；`Width: 420,` 在 `:304`、`Height: 260,` 在 `:305`，`:306` 是 `},`。⇒ 说「写死在 `:304-306`」射程对，逐字两行在 `:304-305`。
   ⚠ 锚 1 没提、与本跳表直接有关的一条：**同一枚 `WindowOptions` 之外，`[panel]` 段的 schema 默认宽度是 640**（`internal/config/schema.go:532 default:"640"`），而 `[panel]` 五枚键非测试读者 0 枚 ⇒ 写死的 420 比配置默认值还小 220px。这是**新增读数**，不是推翻。

3. **锚 2 的另一半复认上**：`internal/config/manager.go:281` 逐字 `{"llm", &cur.LLM, &fresh.LLM, func() { cur.LLM = fresh.LLM }},`，位于 `:271` 那句「Everything else is hot-tier」起的段表里（`{"panel", …}` 在 `:285`）。⇒ "两套说法相反、零行码互相翻译"这一句我复认成立。

4. **锚 3 复认上，但要加一条本腿量到的分叉**。
   跟踪口径：`git ls-files | grep -c "/dist/"` ⇒ **1**（派单说的那枚占位件）。
   同一口径的未跟踪态：`git ls-files --others --exclude-standard | grep -c "/dist/"` ⇒ **0**，`git ls-files --others --ignored --exclude-standard | grep -c "/dist/"` ⇒ **6249**。
   ⇒ 「`dist` 目录里今天只有一枚 `.gitkeep`」在**跟踪件**这一档成立，在**本机工作树**这一档不成立（有一批被忽略的产物字节在树里）。`go:embed` 是**构建期读工作树**，所以"今天这台机器上造出来的二进制有没有带页面"我判不了（判它要跑，或要看那批文件的形状 ⇒ 双双禁区）。这一条改的是锚 3 那句「票 33 的 AC#13 因此必跳」的**前提适用面**：干净检出必跳，本机不一定。⛔ 我不点名那枚目录，⛔ 我不读页面内容。
   `tools/d22scan/runtests.sh` 那把尺复认：`skipped != 0` ⇒ `exit 1`，红句逐字含「SKIP is not a pass (ticket 71 AC#3)」（尺在 `:98-102`，不是我推的读数）。

5. **锚 4（并发腿的写面）我不裁**，只交三发卫生读数：`11:34:18`／`12:11:34`／`12:46:21` 三发 `git status --porcelain -- cmd internal` **都是 0 行**。`fb1d9c19` 正文自陈 `248-v1c` 已交件、残留突变尺（它的口径是 `cmd internal tools scripts`）＝0 行——那是**它的**读数，我这只复量了我自己那两枚根。本腿全程未碰 `docs/evidence/s1/248-settings-write-path-v1.md`，未碰 `.scratch/wisp/probes/248/`，未改任何产码／测试／票面（写点只有本件＋三枚 commit-msg 临时件，均在 `.scratch/` 下，只建不删）。

## §7 一句话结论

**从今天到"机主能点开设置页并存下去"，按依赖顺序最少还差这六跳（只列跳）：**

1. **H2** 双击那条入口上有人把首份 `config.toml` 建出来（今天建件只挂 `wisp run`）——否则面板压根装不起来，后面五跳全部不可达。
2. **H11** 首份配置里就有一行**可被选择**的服务商／模型，或名册里有一条机主走得通的落点（今天七枚字段在干净机器上全部被选择器前置挡死）。
3. **H5** embed 里真的带着一包页面产物（跟踪件 1 枚 ⇒ 干净检出必显示"未建"告示页），且 AC#13 那枚"最终文档含入口真内容"的会响断言落上。
4. **H6** 页面里那条"设置"路存在并且按 `config.get`／`config.set` 这两枚名发得出（Go 侧两枚名已在，界面那半本腿射程外，由 owner 带话）。
5. **H4** 机主手上的手势形：今天只有热键与托盘两形，"点击/双击球"那一形不接面板（＋双击是否属票 33 射程＝未定义即停，待裁）。
6. **H13** 写下去之后那套"什么时候生效"的说法只剩一套（恒 restart vs `llm` 在 hot 段表）、宿主拿到配置对象、`[panel]` 那五枚键有读者；H14（面板侧快照泵与 Go→页推送）在它之后才有意义。

一句话：**H2 与 H11 是"今天连门都进不去"那两跳，H5 与 H6 是"进去了也看不见那页"那两跳，H4 与 H13 是"看见了也存不成机主要的那种生效"那两跳。**

---

### 交件态（本腿自己量到的收尾，供编排者拿三把尺复量）

- 写点一枚：`.scratch/wisp/probes/e2e-panel-1/readiness.md`。两发 commit 都带**显式 pathspec**、只含这一枚文件：`4a9851d6`（§0＋§4/§5 先写满）、`41bea161`（§1/§2/§3/§6/§7 正身）。⛔ 未 push。
- 零跑复述：`go build`／`go vet`／`go test`／任何 `./...`／`go list -deps` **一次都没执行**（本腿全程只用 `grep`／`sed`／`ls`／`wc`／`git log`／`git show --stat`／`git status`／`git cat-file blob HEAD:`／`git ls-files`）。
- 零污染自证（三发现量，逐发带时刻，⛔ 不是一句"0 行"糊过去）：
  `git status --porcelain -- cmd internal tools scripts` 在 `12:52`／`12:57:25`／`12:59:36` 三发**都是 0 行**。
  `docs/**` 那一枚根另有一发读数：`12:57:25` 量到 **1 行** ` M docs/evidence/s1/248-settings-write-path-r1.md`，
  `12:59:36` 复量该根 **0 行**（其间那枚文件被别的腿自己 commit 掉了）。⇒ **那枚 M 不是我写的**（本腿写点只有 §下面那枚 md，
  两发 commit 的 `--name-only` 合起来也只有它），但它证明一件事：**本仓的 `docs/evidence/s1/` 此刻确实有另一枚腿在写**，
  我第一发把这条写成"交件前现量 0 行"是**说快了的假话**，就地改成带时刻的分根读数。
- 冻结件（`docs/PLAN.md`／`docs/specs/**`／`docs/BUILD.md`／`docs/SLO.md`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`internal/panel/tokens_fourway_test.go`）**一字未动**；
  其中 `manager.go` 与 `runtests.sh` 那两枚属"允许读"射程，本腿**只读不写**（引用为内容引用的地方逐处可见：§1 H13、§6 第 3/4 条）。
- 未定义即停共四处，全部只在件里具名上报、⛔ 没有一处被我自己填掉：**H4 的双击那一形是否票 33 射程**（§3 H4）、**H9 的回话形状要不要结构化字段**（C17 面，§4 第 6 条）、**H14 出向那一半的归口（码内注释叫票 35，本腿没读那张票面）**（§3 H14／§5 第 9 条）、**改 `[hotkey]` 没人执行这一格归谁**（§3 H4 末段）。
- 我只量到 H14 为止：H1..H14 十四跳每跳都给了四问，⛔ 没有任何一跳是"应该有"；判不了的那几格一律落在 §5 并带**为什么判不了**的具名尺（零跑／禁令树／需真机时序），没有一格被含糊过去。
