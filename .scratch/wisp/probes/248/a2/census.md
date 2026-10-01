# 248-a2 普查件 — 设置路由：AC#1..AC#11 里哪几格今天可 Go 侧独立做完

> 只读普查腿（代号 `248-a2`；`248-a1` 是前一枚，已交）。⛔ 零产码／零脚本／零测试／零构建／零 `go vet`／零 `go mod tidy`。
> 只用 Read／Grep／Glob／`git cat-file blob`／`git show --numstat`。`cmd/wisp` 有写腿 `33-r4` 正在量启动耗时，本腿全程不吃 CPU、不改盘。
> 两层禁令：不读也不引 `frontend/**`／`design/**`。三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）只点名＋只引编排者已裁的 J1/J2 行号，**不把它内容当本腿的凭据**。

## 起手锚点

- 钟点（`date` stdout）：`2026-10-01 11:51:49 +0800`
- `git log -1 --format=%H`：`23627aba5d64b09f92ec8687a2716ba0668a6b38`；分支 `dev`
- （两者与 `ls .scratch/wisp/probes/248/`、`ls issues | grep 248` **同发取**：探针目录只有 `a1` ⇒ `a2` 未用过，无需换号；工单唯一文件 `248-the-panel-has-no-settings-route-...md`。）
- 本腿判定范围＝票 248 的 **AC#1..AC#11**（AC#0 是"先把三问摆开"元判据，写点在 `a1/census.md`，非本腿）。

## ① 逐格分类表（Go 侧可独立 / 要界面侧先有东西 / 两栖）

| AC | 分类 | 一句凭据（该 AC 要求的那个动作，今天有没有 Go 侧承载点） |
|---|---|---|
| **AC#1** 白名单真扩（新方法进 `bridge.go` 白名单＋`ParseComposerRequest`＋负向配正控） | **Go 侧可独立** | 白名单在 `internal/panel/bridge.go:41-46`（四枚常量）与 `knownComposerMethod` `:104-110`、派发在 `composer_dispatch.go:137-162`（`switch req.Method`）——三处都是纯 Go，可加名；正控"种一枚不在白名单必被拒"今天已由 `ParseComposerRequest:90-92` 实现，可照抄。⚠ 但**必撞冻结件** `l2_grant_boundary_test.go:2139`（`len(...answered)!=4`）与 `git_test.go:385/:517`（名册数死四枚），改法受 J1/J2 约束（见 §④）。 |
| **AC#2** 快照"凭据已录入／配置可读"那一维＋全仓 grep 不到 key 值（哨兵零命中、升常驻用例） | **两栖（Go 侧可先做一半）** | Go 侧承载点＝出向快照 `internal/panel/composer.go` 的 `composer` 段（J3 裁定：进 composer 段，形如 `composer.go:255-256` 一枚事实＋一枚 known 位；非顶层第五键）。"喂哨兵→grep 零命中→升常驻用例"的动作全在 Go＋测试里可跑。⚠ 半等界面＝哨兵要从"页面录入路径"喂才贴 AC#4 真链，但**判据本身可在 Go 侧造一发**，不必等界面。 |
| **AC#3** 只走一枚凭据存储（复用 `secret.NewStore`/DPAPI，不落明文、不新增回显值方法） | **Go 侧可独立** | 存储持有者现成：`cmd/wisp/secret.go` 走 `secret.NewStore`→DPAPI（`internal/secret/dpapi_windows.go:24`）；写路径复用即可，`api_key_ref` 才上页面。全是 Go 侧构造与约束。 |
| **AC#4** 真机那一发（面板点设置→填→保存→下一次任务真用上新模型） | **要界面侧先有东西（且依赖票 33）** | 需 (a) 窗口真存在（票 33 宿主，`PanelManager` 现未接进常驻腿——票面 §5.4 现量 `grep PanelManager main.go/resident_windows.go/run.go`=0 命中）；(b) 页面能点出"设置"。J9 已裁：同一次运行不重建，AC#4 改按"重启后真用上"判。**本腿判不了、也不做**（见 §⑦）。 |
| **AC#5** 越界检查（`git diff` 出现禁区路径即退回） | **Go 侧可独立（且是护栏非工作量）** | 就是"本编队不碰 `frontend/**`/`design/**` 等"。落地腿只要严守禁区即自动满足，无需界面。 |
| **AC#6** 门禁四数（build/gofumpt/go test/d22scan.exe） | **要界面侧先有东西——不，本腿无资格判** | 四道尺都要吃 CPU／跑测试，本只读腿**全程禁跑**（会污染 33-r4 的耗时数）。这格是落地腿交件时自跑的，非普查腿射程。此处**登记为"由落地腿跑"**，不在本件下结论。 |
| **AC#7** 写前必校验（落盘前过启动同一套 `validate()`，种非法值⇒磁盘逐字节不变＋下次仍起得来） | **Go 侧可独立** | 校验实现现成＝`internal/config/validate.go:24` `validate(c *Config)`；但它今天只在上游 `loader.go:123`（`readConfigFile`）与 `migrate.go:77` 跑，**不在写路径上**：`grep validate writeguard.go`＝零命中（本腿复认）。落点＝J/AC#11-3 指定的"新增导出 setter 内、进 `mergeWrite` 之前"。纯 Go 可落。⚠ 它自述"从函数体读出、没实跑"⇒ 种一发非法值须落地腿自跑才算数。 |
| **AC#8** 归口"重建那一跳"（回执不许出现"保存即生效"、"要重启"到达页面可见面） | **两栖（Go 侧可先做一半）** | Go 侧＝回执形状带上 tier（立即／下次任务／需重启），且 `hotReloadDisabledPanelInbound`（`panel_inbound.go:212` 念出）这类"热加载对这条腿无效"的措辞今天**只落 stderr**；把损失写进**回执结构**是 Go 侧能做的另一半。**"到达页面可见面"那一半等界面**（页面必须把回执念出来，`panel_inbound.go` 头注：accepted 时 `Handle` 返 `("","nil")`、不回页面文本，H10＝票 35 host）。 |
| **AC#9** 宿主在、内容不在（`//go:embed all:dist` 真匹配到一包页面产物，条目数＋关键入口真存在） | **要界面侧先有东西** | 产物由界面侧那枚 agent 产出；本编队不写 `frontend/**`。票面 §2b 末：`frontend/dist` 只跟踪一枚 `.gitkeep`（`git ls-files frontend`）。⛔ 本编队读不了 dist 内容（两层禁令），只能登记"这一格在'谁填 dist'落定前勾不了"。 |
| **AC#10** 常驻腿那扇门今天没 `Grants`⇒会话档落 `GRANT-DROPPED`；⛔ 只许 ⓘ/ⓑ 二选一 | **要界面侧先有东西——不，是"交编排者裁、非本腿选"** | 本格"只许两选一，不许第三形"。ⓘ 移 `approval.New` 进 `assembleRuntime`（改票 246 AC#1 乙形次序，需先落 `A##`）；ⓑ 不移动、在回执文案逐字写"这两项只作用于跑任务的进程，常驻腿用常量 300s/3s"。两形都要 (ⓑ) 落到"页面可见回执" ⇒ **另一半等界面**；且**选择权归编排者，本腿不选**（见 §⑦，逐条列现量）。 |
| **AC#11** 配置写只走"按键受保护写"那族（`(*Manager).mergeWrite`），⛔ 不用整份序列化 `SaveFile`；三子判据 | **Go 侧可独立** | 受保护原语签名现量：`internal/config/writeguard.go:111` `func (m *Manager) mergeWrite(ownedKey string, setOn func(base *Config)) error`，三形非静默（见 §③）；先例三枚 `allowdirs.go:62 AddAllowedDir`／`:109 SetAllowedDirs`／`permmode.go:70 SetPermissionMode` 照抄即可；`SaveFile`（`loader.go:238-249`）确为 serializer（deepCopy→MarshalCanonical→atomicWrite，不读回、不校验）。三子判据（并发保 B 键／逐枚键各一发＋报告落盘路径／AC#7 落点）全在 Go＋测试射程内。 |

**一句话汇总**：今天**一枚 Go 侧写腿能独立做完并交凭据**的＝**AC#1 / AC#3 / AC#5 / AC#7 / AC#11**（外加 **AC#2、AC#8 的 Go 那一半**）；**必须先等界面侧／票 33 宿主**才有意义的＝**AC#4 / AC#9**（以及 AC#8/AC#10 的"页面可见"那一半）；**AC#6 交落地腿自跑、AC#10 的选形交编排者**，本普查腿都不替它们下结论。

## ② 路由承载点现量

**`cmd/wisp/panel_inbound.go` 的槽位（现读，非 sed 推行；⚠ `cmd/wisp` 有写腿在改 `*_test.go`，落地时以现读为准）**
- `newPanelInboundDispatch`（`:198`）装配链：`config.NewManager`(`:200`) → `perm.New(perm.Options{Manager, Confirm:nil, Logf})`（`:213-220`，**`:218` `Confirm: nil`**）→ `panel.ModeWriteHandler{Modes, Confirm:nil, Audit, Actor}`（`:224-231`，**`:228` `Confirm: nil`**）→ 返回 `&panel.ComposerDispatch{ Mode: modeWrites, ... }`。
- 派发器四个 handler 槽：`:233` `Mode: modeWrites`（唯一真接）；`:236` `Workspace: nil`；`:237` `Attachment: nil`；**`:238` `Message: nil`**。⇒ 今天**没有 `Config` 这一枚槽**（既不在 struct 里、也不在这里赋）。

**C17 方法名册（`internal/panel`，invoke/parse 那一族，今天有哪几枚、有没有"设置"相关）**
- 白名单常量：`bridge.go:42-45` = `MethodModeRequest "panel.mode.request"`／`MethodWorkspaceRequest "panel.workspace.request"`／`MethodAttachmentAdd "panel.attachment.add"`／`MethodMessageSend "panel.message.send"`。**四枚，零枚 config/secret 相关。**
- 解析入口：`ParseComposerRequest(raw)`（`bridge.go:84`），能力判定在 `knownComposerMethod`（`:104-110`，case 那四枚常量）。⇒ "设置"相关的方法名今天**不存在**，落地腿要新增。
- 路由器：`composer_dispatch.go:137-162` `dispatch` 用 `switch req.Method` 四 case（每 case 引用 bridge.go 的常量，不写裸字符串——头注 `:18-25` 明说写裸 route 串会**故意**触发本包 AST 钉）＋ default `rosterMismatch`（`:159/:167`，`ErrRosterMismatch` `:178`）。空槽⇒`unattached`⇒`ErrNoHandlerAttached`（`:62/:182`）。

**页面能拿到的返回值形状（回执／错误怎么回给页面）**
- `Handle(ctx, raw) (string, error)`（`composer_dispatch.go:120`）：**被拒** ⇒ 返回 `RefusedEnvelopeForUser(req, err)`（`bridge.go:115`，含 method+requestId+原因）＋ err；**被受理** ⇒ 返回 **`("", nil)`**（头注 `:117-119`：切片故意不编成功句，"决定页面被告知什么"＝H10 业务）。
- 生产 leg `panel_inbound.go:159-173`：把拒句打 stdout；受理只打自造的一句"第 N 行已受理，处理器已被调用"——**不是回给页面**，是回给 console。⇒ **Go→页面的回执通道今天在此树不存在。**
- ⚠ **登记不据此下结论**：另一枚腿今量"Go→页面投递队列全仓无取出者"（源码读数在台账 `A497`）。**本腿不因此判"路由做不了"**——该维依赖归票 33 AC#14／35 host；本件只记"这一维有已知欠账"。回执结构（AC#8 的 tier）在 Go 侧**能造**；把它**投到页面**那一跳缺取出者，是另一格。

## ③ 写侧落地面现量

**`(*Manager).mergeWrite` 的签名与三个非静默出口**
- 签名（`writeguard.go:111`）：`func (m *Manager) mergeWrite(ownedKey string, setOn func(base *Config)) error`。`ownedKey`＝该 caller 拥有的点分键路径；`setOn`＝把那枚键盖到传入的 base 上（头注 `:84-88`）。**caller 持 `m.mu`**。
- 出口 (1) 文件读不回⇒拒写并报错：`os.Stat` 成功但 `readConfigFile` 失败 ⇒ `observe.Wrap(observe.ClassConfig, err, "…refused to write config.toml; … cannot be read back …")`（`:118-123`）；`Stat` 既非 nil 也非 missing ⇒ 另一条 `observe.Wrap`（`:128-131`）。均**一字不写**。
- 出口 (2) 已同值⇒不写：`disk != nil` 时 `written = diffKeyPaths(disk, base)`，`len(written)==0` ⇒ `slog.Info("…already carries this key; nothing written")` 并 `return nil`（`:149-155`）。
- 出口 (3) 真写了⇒报告差异键：`SaveFile(m.path, base)`（`:160`）后 `diverged = diffKeyPaths(m.cur, base)`（`:166`）；`diverged` 空 ⇒ `slog.Info(... wrote ...)`＋`m.statOwnWrite()`（`:168-176`）；非空 ⇒ `slog.Warn(... kept_in_file_not_in_memory ...)`、**不 adopt stat**（`:180-186`，AC#3 收窄）。

**`diffKeyPaths` 那侧能不能报告"哪些键被写了、哪些外来手改被保留"**
- **能**。`diffKeyPaths(a,b)`（`:201`）走序列化同一套 struct tag 列点分叶子键路径、排序；`written := diffKeyPaths(disk, base)`＝"这次写实际改了文件的哪几枚键"（对**文件内容**比，不对内存比，注释 `:145-147`）；`diverged := diffKeyPaths(m.cur, base)`＝"文件现在有、本进程内存没有"＝**被保留的外来手改**（`:163-166`）。AC#11-2 要的"回执逐枚列落盘键路径"直接取自 `written`（`writeguard.go:201` 那份报告）。
- 例外跟序列化器一致：`SchemaVersion` 跳过（`:220-223`）、`toml:"-"` 字段跳过、`Config.Plugins` 手工在 `plugins` 下走（`:194-195/:244-245`）；map 逐 key 命名（`[plugins.<id>]`／`[llm.providers.<name>]`，`:284-314`）；报告有 `maxReportedKeys=24` 上限并显式标 `(+N more)`（`:78/:330-336`）。

**配置管理器今天哪几处已在用它（照抄形状）**
- 三枚，全走 `mergeWrite`：`permmode.go:79`（`SetPermissionMode`，键 `keyRiskPermissionMode="risk.permission_mode"`）；`allowdirs.go:142`（`writeAllowedDirs` 被 `AddAllowedDir:82`／`SetAllowedDirs:124` 复用，键 `keyFSAllowedDirs="fs.allowed_dirs"`）。键常量声明在 `writeguard.go:69-72`。⇒ 新增"配置写 setter"照这三枚：先改内存、调 `mergeWrite(ownedKey, setOn)`、失败把内存回滚。
- ⚠ `SaveFile`（`loader.go:238-249`）＝ `deepCopyConfig`→强制 `SchemaVersion`→`MarshalCanonical`→`atomicWrite`，**不读回文件、不跑 validate**——正是票 226 要封的洞（`writeguard.go:5-17` 头注逐字）。AC#11 明禁落地腿用它。

**`validate()` 在不在写路径上**
- **不在。** `validate` 定义在 `internal/config/validate.go:24`；调用点仅两处＝上游读管道 `loader.go:123`（`readConfigFile`）与 `migrate.go:77`。`writeguard.go`／`mergeWrite`／`SaveFile` 全程无 `validate` 调用（本腿 `grep validate internal/config/writeguard.go`＝**零命中**复认）。⇒ AC#7"写前必校验"是**真缺口**；落点（AC#11-3）＝新增导出 setter 内、进 `mergeWrite` 之前跑 `validate()`（⛔ 不改 `mergeWrite` 本体，会波及既有三枚 setter）。
- 现成校验子函数（落地腿可选用）：`validateRisk:109`／`validateModels:95`／`validateModelSpec:195`／`validateCost:117`／`validateAPIKeyRefs:152` 等（`validate.go` 一族），对应设置页那几项（模型名／上下文窗口／价格／权限档位）。

## ④ 会被叫红的现成钉（逐枚点名＋断言原文＋判会不会红）

**落地腿大概会新增／改动的符号名**（据此做撞钉预检）：`MethodConfigGet`／`MethodConfigSet`（新入向常量，J6 已裁名用规格里现成的 `config.get`/`config.set`）、`ComposerDispatch.Config`（新 handler 槽）、`ConfigWriteHandler`／`HandleConfigRequest`、`Manager.Set*`（新导出 setter，如 `SetLLMModel`/`SetContextWindow` 之类）、快照 `composer` 段里"凭据已录入"那一维字段。

> 词面预检：`grep -rn "config.get\|config.set\|MethodConfig\|HandleConfig" --include=*_test.go cmd/ internal/` ⇒ **零命中**（这些名字今天不存在）。**但词面零命中≠没钉**——下面这几枚是**行为型钉**，只写产物／数名册，不写常量名。

**T1 `internal/panel/git_test.go:385`（＋`:386` 文案、`:514/:517/:518` 正控）——名册数死四枚**〔非冻结，射程盖到本票〕
- `:361` 注释："the C17 whitelist in bridge.go is still exactly the four panel.* methods."
- `:385` `whitelistMethodsFromSource(t, .../bridge.go) ; !equalStrings(got, wantMethods)` ⇒ `:386` 报错"the C17 panel.* whitelist in bridge.go = %v, want the four methods"。
- 正控 `:514` `len(whitelistMethodsFromSource(fivePath)) != 5`（planted fifth panel.* method ⇒ 5）；`:517` `len(real) != 4` 对真 `bridge.go`。
- **判会不会红**：**取决于 `whitelistMethodsFromSource` 抽的是"`panel.` 前缀字面量"还是"所有 `Method*` 常量"**——⚠ 见 §⑥ 与 §⑦（本腿**未读到** `:403-446` 该函数体就落此骨架；落地腿务必现读）。若按 `panel.` 前缀数 ⇒ 新增 `config.get`/`config.set`（J6 取的名**不带 `panel.` 前缀**）**可能不被计入** ⇒ `:385/:517` 不红；若按 `Method*` 常量名数 ⇒ 加两枚即变 6 ⇒ **必红**。J2 已裁这两枚（＋`:514`）**同 commit 改写**、形状取乙（能力形"这张名册等于白名单吗"），属票 181 地界。

**T2 `internal/panel/composer_test.go:409 composerRouteLiterals()`／`:426`／`:489/:502 TestTheRendererHoldsExactlyOneDoorToTheHost`——渲染器只准留一枚到宿主的门**〔非冻结〕
- `:406` 注释："composerRouteLiterals is the closed vocabulary the renderer may name. Growing…"；`:502` 函数把 `rep.files`/`rep.callSites`/`len(composerRouteLiterals())` 摆一起（`:525`）。
- **判会不会红**：这枚钉"页面侧出现的 `panel.*` route 字面量必须都在 Go 侧应答的名册里"。⚠ 与 T1 同理：**名册若只收 `panel.` 前缀、且由这 4 枚常量派生**，则新增非 `panel.` 前缀的 `config.*` 常量**不进 `composerRouteLiterals`**，需现读确认它是"扫 `panel.` 字面量比对常量"还是"扫全部 Method 常量"。落地腿加常量同时更新这枚派生表即可；但**它一旦红，就是"Go 侧多了页面没对应的门"或反之**——务必与 T1 同批处理。⚠ 本腿未读 `:409-426` 函数体，判语标〔推〕，落地腿现读为准。

**T3 `internal/panel/l2_grant_boundary_test.go`（FROZEN）——answered 名册尺寸钉＋拼写钉**〔冻结件：只点名，凭据引编排者 J1〕
- 依票面 **J1（台账 `A487`）** 现量两枚锚：`:2051` 拼写钉 `guardAnchor := "case MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend:"`；`:2139` 尺寸钉 `if got := len(sortedSet(pkg.answered)); got != 4`。（本腿 `grep -rn ... l2_grant_boundary_test.go` 只复核到 `:2139` 那行确实存在，**不引其余内容当凭据**。）
- **判会不会红**：**会**，且是最硬的一枚。`answered` 集由"跑起来的 guard 应答哪些 route"推出（＝`knownComposerMethod` case 列）；落地腿一旦把 config 方法加进 `knownComposerMethod`，answered 从 4→6 ⇒ `:2139` `!=4` 直红；拼写钉 `:2051` 那枚 `case …四枚…` 的字面也变 ⇒ 直红。**J1 已裁＝甲（人工批准、只动那两处锚、只许变强、把"数一枚 4"换成能力形＋自带正控）**。另有一维本腿**判不了**：`answeredGrantDoor`／"approval-shaped"探测器会不会把 `config.set`（写配置键）误认成 grant 门 ⇒ 见 §⑦（涉冻结件内部语义，不引内容）。

**T4 `cmd/wisp/panel_inbound_33_test.go`——票 33 dispatch 的"名册 vs nil 槽"盘上读数**〔非冻结，cmd/wisp 行号落地现读为准〕
- `panel_inbound.go:13-14` 头注自述：`panel_inbound_33_test.go`"reads THAT fact off the disk rather than off this paragraph, so deleting the line below reddens a case"。⇒ 它很可能钉"四个槽里只有 Mode 非 nil、其余三枚 nil"的装配形状。
- **判会不会红**：若落地腿给 `ComposerDispatch` 加 `Config` 槽并在生产装配里赋非 nil，枚数／nil 形状变 ⇒ **可能红**。需现读该 `_test.go` 断言原文再判（本腿为守"cmd/wisp 落地现读"与工具预算，只点名、标〔推〕）。

**T5 `cmd/wisp/leg_dispatch_gate_133_test.go`＋`internal/panel/composer_dispatch_test.go:387`——"不许有第二套白名单"**〔非冻结〕
- `composer_dispatch_test.go:387` 报错串："…constant: this is a second whitelist, and the guard's own case list would never have named it"。⇒ 若落地腿在 bridge.go 之外另立方法集合（例如在 dispatch 里写 route 字面量），红。
- **判会不会红**：只要新方法常量声明进 `bridge.go` 一处、dispatch 用常量引用（照现有四枚形），**不红**；把名字散落两处才红。

**T6 `internal/panel/frontend_hygiene_test.go`／`assets_test.go`＋`cmd/wisp/panel_assets_143_test.go`——文件清单／embed 资产钉**〔非冻结〕
- AC#9 那一维（embed dist 内容）可能撞"资产清单"类尺；本腿读不了 `frontend/**`（两层禁令），只能标"落地腿务必现读这三枚，确认它们扫的是文件清单还是 embed 结果"。⚠ 与 AC#9"要界面侧先有东西"绑定。

**T7 `internal/config/writeguard_226_test.go`／`cmd/wisp/always_write_no_clobber_226_test.go`——票 226 不覆写钉**〔非冻结〕
- AC#11 复用 `mergeWrite` 即落在它们射程内。判语：只要新 setter 走 `mergeWrite(ownedKey, setOn)`、不改 `mergeWrite` 本体，**不红**（这两枚钉的是"手改的别处键要活着"这一产物形，新增走同一路的子键 setter 不破坏它）。**待落地腿现读确认它们是否数"用 mergeWrite 的 setter 枚数"**——若数枚数，新增 setter 会让计数变 ⇒ 可能红（标〔推〕）。

**T8 冻结件 `internal/perm/ticket90_persist_test.go`＋`internal/panel/tokens_fourway_test.go`**〔冻结件：只点名，不引内容〕
- `ticket90_persist_test.go` 射程＝permission_mode 那枚 persist（`permmode.go:70` 已在用 mergeWrite）。若落地腿**不动** `SetPermissionMode` 本体、只为其它键新增 setter ⇒ 预期不红；若把 `permission_mode` 划进通用 `config.set` 的键集 ⇒ **正面撞 J5（甲：显式拒写 `risk.*` 全族）**，须走带 L2 卡那条路，不走配置写方法。`tokens_fourway_test.go` 射程＝token 四态，与设置路由**无直接交集**（预期不红，标〔推〕）。

## ⑤ 我跑了哪些尺、每条真实读数（全是只读尺；零 CPU／零改盘）

- `ls .scratch/wisp/probes/248/` ⇒ 只有 `a1`（含 59,623 字节 `census.md`）⇒ **`a2` 未用过，不换号**。
- `ls .scratch/wisp/issues/ | grep 248` ⇒ 一枚 `248-the-panel-has-no-settings-route-the-inbound-whitelist-has-zero-config-or-credential-methods.md`（凭名取准，未猜）。
- `date "+%Y-%m-%d %H:%M:%S %z"` ⇒ `2026-10-01 11:51:49 +0800`；`git log -1 --format=%H` ⇒ `23627aba5d64b09f92ec8687a2716ba0668a6b38`；`git rev-parse --abbrev-ref HEAD` ⇒ `dev`（同发取）。
- Read `cmd/wisp/panel_inbound.go`（269 行）⇒ `:218 Confirm:nil`／`:228 Confirm:nil`／`:236 Workspace:nil`／`:237 Attachment:nil`／`:238 Message:nil`；**无 Config 槽**；`:159-173` 受理只打 console。
- Read `internal/panel/bridge.go` ⇒ 白名单常量 `:42-45`（四枚，零 config）；`ParseComposerRequest:84`；`knownComposerMethod:104-110`；`RefusedEnvelopeForUser:115`。
- Read `internal/panel/composer_dispatch.go` ⇒ `Handle:120` 返回 `(string,error)`、受理⇒`("","nil")`（`:130`）；`dispatch:137-162` switch 四 case＋default `rosterMismatch:167`；`ErrNoHandlerAttached:62`、`ErrRosterMismatch:178`。
- Read `internal/panel/composer_handlers.go` ⇒ `ModeWriteHandler:86`、`HandleModeRequest:111`（`Confirm:nil`⇒变宽拒 `:133-138`）。
- Read `internal/config/writeguard.go` ⇒ `mergeWrite:111` 签名＋三形（`:118-131` 读不回拒／`:149-155` 同值不写／`:160-186` 真写报告）；`diffKeyPaths:201`；键常量 `:69-72`。
- Read `internal/config/permmode.go`／`allowdirs.go` ⇒ 三枚先例 `SetPermissionMode:70`/`:79`、`AddAllowedDir:62`、`SetAllowedDirs:109`、`writeAllowedDirs:139`→`mergeWrite:142`；`statOwnWrite:157`。
- Read `internal/config/loader.go` ⇒ `SaveFile:238-249`（serializer：deepCopy→MarshalCanonical→atomicWrite，无 re-read、无 validate）；`readConfigFile:64`→`validate:123`。
- `grep -n "func validate|validate\(" internal/config/*.go` ⇒ `validate` 定义 `validate.go:24`；调用仅 `loader.go:123`＋`migrate.go:77`；**`writeguard.go` 内零 `validate`**。
- `grep` 名册／尺寸钉（`answered|sortedSet|composerRouteLiterals|whitelist|len(...)` 在 `internal/panel/*_test.go`）⇒ 命中 `git_test.go:385/514/517/518`、`composer_test.go:406/409/426/489/502/525`、`l2_grant_boundary_test.go:2139`（＋`:1239/:1304` 名册核对）、`composer_dispatch_test.go:387`。
- `Glob cmd/wisp/*_test.go`（47 枚）／`internal/config/*_test.go`（12 枚）⇒ 点名列 T4–T8 相关（`panel_inbound_33_test.go`、`leg_dispatch_gate_133_test.go`、`writeguard_226_test.go`、`always_write_no_clobber_226_test.go`、`panel_assets_143_test.go`、`frontend_hygiene_test.go`、`assets_test.go`）。

## ⑥ 我可能写错的条目（对抗我自己）

1. **T1/T2 的"会不会红"是我这发最可能错的地方**：`whitelistMethodsFromSource`（`git_test.go:403-446`）与 `composerRouteLiterals`（`composer_test.go:409`）我**没读到函数体**就写了骨架，判语是〔推〕。若它俩**只按 `panel.` 前缀抽名**，则 J6 选的 `config.get`/`config.set`（不带 `panel.` 前缀）**不会被计入 4**⇒ `git_test.go:385/:517`、`composer_test.go:502` 都**不红**，那 §④ 把 T1/T2 列成挡路项就是夸大。**落地腿务必现读这两枚函数体再定性**，别照我这张表决定要不要改票 181 的地界。
2. **J1 已把 T3（冻结件）裁成"甲：只动那两处锚"**，我却仍把它列进"会被叫红的钉"——严格说它**会被叫红、但改法已被授权**。我把"挡路"和"须按裁定改法"混在一张表里，读者可能误读成"落地腿不能碰"。**更正口径**：T3 的红是**预期的、且 J1 给了改法**，不是否决项。
3. **panel_inbound.go 的行号**：我按 Read 的屏号写 `:218/:228/:236/:237/:238`，虽非 sed 推、但仍可能因 `33-r4` 正在改 `cmd/wisp` 而漂（`panel_inbound.go` 是非测试文件、相对稳，但派单明令"cmd/wisp 行号落地现读为准"）⇒ 已在 §②/§④ 逐处标注"落地现读为准"。
4. **AC#8 判成"两栖"我可能偏松**：Go 侧能造带 tier 的回执结构是〔量〕，但"到达页面可见面"整半枚没页面就无从验；若严格看，AC#8 的**完成判据**（回执文案不许出现"保存即生效"且"要重启"到页面）今天只能做到"Go 侧不再只落 stderr"这一半，判"两栖"成立、但**别读成 Go 侧能勾 AC#8**。
5. **AC#6 我不该归进"逐格分类"**：它四道尺全要跑，本只读腿没资格判红/绿；我写成"交落地腿自跑"是把它从射程剔出，但分类表里它占了格——**读者别以为我核过门禁**。
6. **A497（Go→页面投递队列无取出者）我只登记、没据此下结论**，符合派单；但我 §② 末段仍可能让读者以为"回执通道做不了＝路由做不了"。**更正**：路由（方法名册＋dispatch）与"把回执投回页面"是两件事，前者 Go 侧可做（AC#1），后者是票 33/35 的欠账。

## ⑦ 判不动的地方（逐条甲／乙／不做＋现量＋为什么判不了）

1. **AC#10 的 ⓘ/ⓑ 二选一**——**不做（交编排者裁，派单明令本腿不许选）**。现量：票面 §2c 给 `approval/gate.go:139-145` 钳位、`approval/queue.go:107/115/122`（300s/3s/Max3s）、`resident_approval_windows.go:109`（`approval.New` 建在会话账本前⇒`Options.Grants` nil⇒`GRANT-DROPPED`）。ⓘ 移 `approval.New`＝改票 246 AC#1 乙形次序、需先落 `A##`；ⓑ 不移动＝回执文案逐字声明。**为什么判不了**：选择权与 `A##` 批准记录归编排者，本腿只读、不选、不排落地顺序。
2. **`config.set`/`config.get` 会不会被 `l2_grant_boundary_test.go` 的"approval-shaped/grant-door 探测器"误判**——**判不了**。现量：T3 是冻结件，`answeredGrantDoor`/facet 2 的匹配规则在冻结件内部。⛔ 约束禁我把冻结件内容当凭据引。**为什么判不了**：要看那枚探测器按什么词／什么结构认"grant 门"，只能引编排者；本腿只能标"若落地腿被它打红，那是冻结件射程、须走 J1 的解冻改法"。
3. **J6 名的 `config.get`/`config.set` 到底带不带 `panel.` 前缀、会不会绕开 T1/T2 的"数死四枚"**——**判不了（须现读 `whitelistMethodsFromSource`/`composerRouteLiterals` 函数体）**。甲：若按 `Method*` 常量名数⇒加两枚即 6⇒T1/T2/T3 全红，须同批改票 181 地界。乙：若只按 `panel.` 前缀抽⇒`config.*` 可能不进名册计数、T1/T2 不红、唯 T3（answered 集走 knownComposerMethod）仍红。**为什么判不了**：两函数体本腿为守预算未读；且这直接决定"要不要动 git_test.go（票 181）"，属越界判断。
4. **AC#4／AC#9 能否今天勾**——**不做**。现量：`grep PanelManager main.go/resident_windows.go/run.go`=0 命中（票面 §5.4）、`git ls-files frontend`⇒只 `frontend/dist/.gitkeep`（票面 §2b 末）；宿主未接进常驻、页面产物未填。**为什么判不了**：依赖票 33 宿主落地＋界面侧产出 dist；两层禁令禁我读 dist 内容。J9 另裁 AC#4 改按"重启后真用上"，重建那一跳是票 223 地界。
5. **新 setter 的键集合是否含 `permission_mode`/`allowed_dirs`**——**判不了（触 J5 边界）**。现量：J5 裁甲＝配置写方法**显式拒写 `risk.*` 全族**，改档只走带 L2 卡那条（`cmd/wisp/run.go confirmModeSwitch`＋`panel.ModeWriteHandler`）。**为什么判不了**：`allowed_dirs`（`fs.*`）J5 没逐字点名，是否同样禁由配置写方法改，是编排者的推广边界；本腿只标"照 J5 至少 `risk.*` 全族必拒"。

## 交件判语

- **射程**：AC#1..AC#11 逐格分类（§①）＋路由承载点／写侧落地现量（§②§③，带 `file:line`〔量〕）＋撞钉预检点名八组尺（§④）。核心结论：**AC#1/AC#3/AC#5/AC#7/AC#11 今天可 Go 侧写腿独立做完并交凭据**；**AC#2/AC#8 各能做 Go 那一半、另一半等界面可见面**；**AC#4/AC#9 必须等票 33 宿主＋界面侧 dist 产物**；**AC#6 交落地腿自跑（本只读腿无资格判门禁）**；**AC#10 的选形、AC#4/9 的宿主、J6 前缀与 grant-door 探测——一律交编排者裁（§⑦，未替其下结论）**。
- **没答的（不写成答了）**：① `whitelistMethodsFromSource`（`git_test.go:403-446`）与 `composerRouteLiterals`（`composer_test.go:409-426`）函数体**未读**⇒T1/T2"会不会红"是〔推〕，落地腿务必现读再定性；② `cmd/wisp` 里 T4（`panel_inbound_33_test.go`）、T7（`writeguard_226_test.go`/`always_write_no_clobber_226_test.go` 是否数 setter 枚数）断言原文**未逐枚读**⇒标〔推〕，落地时现读为准；③ 冻结件 T3/T8 内部匹配语义按禁令**未引**⇒只点名。
- **最需要编排者裁的那一格**：§⑦-3（`config.*` 前缀决定要不要惊动票 181 的 `git_test.go`）与 §⑦-2（grant-door 探测器会不会误判 `config.set`）——这两格决定 AC#1 落地时要不要同时改名册计数钉／走 J1 解冻改法。
