# 248-a2 普查件 — 设置路由：AC#1..AC#11 里哪几格今天可 Go 侧独立做完

> 只读普查腿（代号 `248-a2`；`248-a1` 是前一枚，已交 432 行/59,623 字节）。⛔ 零产码／零脚本／零测试／零构建／零 `go vet`／零 `go mod tidy`。
> 只用 Read／Grep／Glob／`git cat-file blob`／`git show --numstat`。`cmd/wisp` 有写腿 `33-r4` 正在量面板冷/热启动**耗时**（冷≤1500ms／热≤200ms），本腿全程不吃 CPU、不改盘。
> 两层禁令：不读也不引 `frontend/**`／`design/**`（本件的 frontend 事实**只来自 Go 侧测试文件对它的描述**，绝不引页面内容）。三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）只点名＋只引编排者已裁的 J1/J2 行号，**不把它内容当本腿的凭据**。

## 起手锚点

- 钟点（`date` stdout）：`2026-10-01 11:51:49 +0800`（追加进度行时刷新为 `11:57:07 +08`）
- `git log -1 --format=%H`：`23627aba5d64b09f92ec8687a2716ba0668a6b38`；分支 `dev`
- （两者与 `ls .scratch/wisp/probes/248/`、`ls issues | grep 248` **同发取**：探针目录只有 `a1` ⇒ `a2` 未用过，无需换号；工单唯一文件 `248-the-panel-has-no-settings-route-...md`。）
- 本腿判定范围＝票 248 的 **AC#1..AC#11**（AC#0 是"先把三问摆开"元判据，写点在 `a1/census.md`，非本腿）。

## ① 逐格分类表（Go 侧可独立 / 要界面侧先有东西 / 两栖）

| AC | 分类 | 一句凭据（该 AC 要求的那个动作，今天有没有 Go 侧承载点） |
|---|---|---|
| **AC#1** 白名单真扩（新方法进 `bridge.go` 白名单＋`ParseComposerRequest`＋负向配正控） | **Go 侧可独立（但必撞冻结钉 T3）** | 白名单 `bridge.go:42-45`、`knownComposerMethod:104-110`、派发 `composer_dispatch.go:137-162`——三处纯 Go，可加名；负向正控今天已实现（`ParseComposerRequest:90-92` 拒未列名＋`composer_dispatch.go` `ErrNoHandlerAttached`/`ErrRosterMismatch`）。⚠ **唯一"无论命名都躲不开"的挡路项＝冻结件 `l2_grant_boundary_test.go:2051/:2139`**（answered 尺寸钉），改法受 J1/A487 约束；`git_test.go` 那两枚"数死四枚"的门（T1）**在 J6 的 `config.*` 命名下不红**（详见 §④）。 |
| **AC#2** 快照"凭据已录入／配置可读"那一维＋全仓 grep 不到 key 值（哨兵零命中、升常驻用例） | **两栖（Go 半枚今天做不到 self-green）** | Go 侧能：把那一维放进 `ComposerState`（J3 裁：进 composer 段、形如 `composer.go:255-256` 一枚事实＋一枚 known 位；**非顶层第五键**——`pump_test.go:123` 字节钉死顶层四键）＋把哨兵喂进写路径跑零命中＋升常驻用例。⚠ 但 `composer_test.go:48 TestComposerContractTypesMatchFrontend` 把 `ComposerState` 与 `frontend/src/lib/panel.ts` 的 `ComposerState` **双向对账**（`:73` 报"Go emits key TS does not declare"）⇒ **加 composer 子键即把 `internal/panel` 测试洗红**，必须与界面侧同批改 panel.ts（同 `Tasks` 段 `composer.go:86-91` 那枚先例）。⇒ "哨兵零命中"这一半可独立，"新维度上快照"那一半被两界对账门卡住。 |
| **AC#3** 只走一枚凭据存储（复用 `secret.NewStore`/DPAPI，不落明文、不新增回显值方法） | **Go 侧可独立** | 存储持有者现成＝`cmd/wisp/secret.go`→`secret.NewStore`→DPAPI（`internal/secret/dpapi_windows.go:24`）；写路径复用即可，`api_key_ref` 才上页面。全是 Go 侧构造与约束，无对账门。 |
| **AC#4** 真机那一发（面板点设置→填→保存→下一次任务真用上新模型） | **要界面侧先有东西（且依赖票 33 宿主）** | 需 (a) 窗口真存在（票 33，`grep PanelManager main.go/resident_windows.go/run.go`=0 命中，票面 §5.4）；(b) 页面点得出"设置"。J9 已裁同一次运行不重建⇒改按"重启后真用上"判。**本腿判不了、不做**（§⑦）。 |
| **AC#5** 越界检查（`git diff` 出现禁区路径即退回） | **Go 侧可独立（是护栏非工作量）** | 就是"本编队不碰 `frontend/**`/`design/**`/规格/冻结件"。落地腿严守禁区即自动满足。 |
| **AC#6** 门禁四数（build/gofumpt/go test/d22scan.exe） | **交落地腿自跑——本只读腿无资格判** | 四道尺全吃 CPU／跑测试，本腿**全程禁跑**（会污染 33-r4 的耗时数）。此处登记"由落地腿跑"，**不在本件下红/绿结论**。 |
| **AC#7** 写前必校验（落盘前过启动同一套 `validate()`，种非法值⇒磁盘逐字节不变＋下次仍起得来） | **Go 侧可独立** | 校验实现现成＝`internal/config/validate.go:24`；但它今天只在 `loader.go:123`（读管道）＋`migrate.go:77` 跑、**不在写路径**（`grep validate internal/config/writeguard.go`＝零命中，本腿复认）。落点＝AC#11-3 指定的"新增导出 setter 内、进 `mergeWrite` 之前"。纯 Go 可落；⚠ 它自述"没实跑过"⇒ 种非法值须落地腿自跑才算数。 |
| **AC#8** 归口"重建那一跳"（回执不许"保存即生效"、"要重启"到页面可见面） | **两栖（Go 侧可先做一半）** | Go 侧＝回执形状带 tier（立即/下次任务/需重启）；且"热加载对这条腿无效"今天只落 stderr（`panel_inbound.go:212` 念 `hotReloadDisabledPanelInbound`），写进**回执结构**是 Go 侧能做的另一半。**"到页面可见面"那一半等界面**（accepted 时 `Handle` 返 `("","nil")`、不回页面文本，`composer_dispatch.go:117-119/:130`）。 |
| **AC#9** 宿主在、内容不在（embed dist 真有一包页面产物） | **要界面侧先有东西** | 产物由界面侧那枚 agent 产出；本编队不写 `frontend/**`、且读不了 dist 内容（两层禁令）。票面 §2b 末：`git ls-files frontend`⇒只 `frontend/dist/.gitkeep`。这一格在"谁填 dist"落定前勾不了。 |
| **AC#10** 常驻腿那扇门今天没 `Grants`⇒会话档落 `GRANT-DROPPED`；⛔ 只许 ⓘ/ⓑ 二选一 | **判不了／选形交编排者（且 ⓑ 那一半等界面）** | 本格"只许两选一、不许第三形"。ⓘ 移 `approval.New` 进 `assembleRuntime`（改票 246 AC#1 乙形次序、需先落 `A##`）；ⓑ 不移动、在**回执文案**逐字写"这两项只作用于跑任务的进程，常驻腿用常量 300s/3s"。**选择权归编排者、本腿不选**；且 ⓑ 要落到"页面可见回执"⇒ 另一半等界面（见 §⑦-1）。 |
| **AC#11** 配置写只走"按键受保护写"那族（`(*Manager).mergeWrite`），⛔ 不用 `SaveFile`；三子判据 | **Go 侧可独立** | 受保护原语＝`writeguard.go:111` `func (m *Manager) mergeWrite(ownedKey string, setOn func(base *Config)) error`，三形非静默（§③）；先例三枚 `allowdirs.go:62/:109`／`permmode.go:70` 照抄；`SaveFile`（`loader.go:238-249`）确为 serializer（deepCopy→MarshalCanonical→atomicWrite、不读回、不校验）。三子判据（并发保 B 键／逐枚键各一发＋报告落盘路径／AC#7 落点）全在 Go＋测试射程内；票 226 那两枚测试（T7）**不数 setter**⇒不红。 |

**一句话汇总**：今天**一枚 Go 侧写腿能独立做完并交凭据**＝**AC#3 / AC#5 / AC#7 / AC#11**，外加 **AC#1**（Go 侧可做，但落地即撞唯一那枚必红冻结钉 T3，须走 J1 解冻改法）；**两栖、且 Go 半枚做不到 self-green**＝**AC#2**（被 `composer_test.go:48` 两界对账门卡住）／**AC#8**（页面可见那一半等界面）；**必须先等界面侧／票 33 宿主**＝**AC#4 / AC#9**（及 AC#10 的"页面可见"那一半）；**AC#6 交落地腿自跑、AC#10 选形交编排者**——本普查腿不替它们下结论。

## ② 路由承载点现量

**`cmd/wisp/panel_inbound.go` 的槽位**（Read 现读，非 sed 推行；⚠ `cmd/wisp` 有写腿在改 `*_test.go`，落地时以现读为准）
- `newPanelInboundDispatch`（`:198`）装配链：`config.NewManager`(`:200`) → `perm.New(perm.Options{Manager, Confirm:nil, Logf})`（`:213-220`，**`:218` `Confirm: nil`**）→ `panel.ModeWriteHandler{Modes, Confirm:nil, Audit, Actor}`（`:224-231`，**`:228` `Confirm: nil`**）→ 返回 `&panel.ComposerDispatch{Mode: modeWrites, ...}`。
- 派发器 handler 槽：`:233` `Mode: modeWrites`（唯一真接）；`:236` `Workspace: nil`；`:237` `Attachment: nil`；**`:238` `Message: nil`**。⇒ 今天**没有 `Config` 这一枚槽**（既不在 `ComposerDispatch` struct、也未在此赋值）。

**C17 方法名册（`internal/panel`，invoke/parse 那一族，今天有哪几枚、有没有"设置"相关）**
- 白名单常量：`bridge.go:42-45` = `MethodModeRequest "panel.mode.request"`／`MethodWorkspaceRequest "panel.workspace.request"`／`MethodAttachmentAdd "panel.attachment.add"`／`MethodMessageSend "panel.message.send"`。**四枚，零枚 config/secret 相关。**
- 解析入口：`ParseComposerRequest(raw)`（`bridge.go:84`），能力判定 `knownComposerMethod`（`:104-110`，case 那四枚常量）。⇒ "设置"相关方法名今天**不存在**，落地腿要新增（J6 已裁名用规格里现成的 `config.get`/`config.set`）。
- 路由器：`composer_dispatch.go:137-162` `dispatch` 用 `switch req.Method` 四 case（每 case 引用 bridge.go 常量、不写裸串——头注 `:18-25` 明说写裸 route 串会**故意**触发本包 AST 钉）＋ default `rosterMismatch`（`:159/:167`，`ErrRosterMismatch:178`）。空槽⇒`unattached`⇒`ErrNoHandlerAttached`（`:62/:182`）。

**页面能拿到的返回值形状（回执／错误怎么回给页面）**
- `Handle(ctx, raw) (string, error)`（`composer_dispatch.go:120`）：**被拒**⇒返回 `RefusedEnvelopeForUser(req, err)`（`bridge.go:115`，含 method+requestId+原因）＋err；**被受理**⇒返回 **`("", nil)`**（`:117-119/:130`：切片故意不编成功句）。
- 生产 leg `panel_inbound.go:159-173`：拒句打 stdout；受理只打自造一句"第 N 行已受理，处理器已被调用"——**不是回页面**、是回 console。⇒ **Go→页面的回执通道今天在此树不存在。**
- **出向快照那一维（AC#2）**：`Snapshot`（`composer.go:57-92`）＝顶层四键 `Pending/Results/Composer/GeneratedAt`（`:58-61`）＋两可选 `Instructions/Tasks`（ptr+omitempty）；`pump_test.go:123` 字节钉＝marshal 后 JSON 键必须恰为 `"composer,generatedAt,pending,results"` ⇒ **新维度进 `ComposerState` 子段、不许加顶层键**（与 J3 一致）。`ComposerState` 由 `composer_test.go:48` 与 `frontend/src/lib/panel.ts` 双向对账（见 §④ AC#2 那格）。
- ⚠ **登记不据此下结论**：另一枚腿今量"Go→页面投递队列全仓无取出者"（源码读数在台账 `A497`）。**本腿不因此判"路由做不了"**——该维依赖归票 33 AC#14／35 host。路由（名册＋dispatch）与"把回执投回页面"是两件事：前者 Go 侧可做（AC#1），后者是另一格欠账。

## ③ 写侧落地面现量

**`(*Manager).mergeWrite` 的签名与三个非静默出口**
- 签名（`writeguard.go:111`）：`func (m *Manager) mergeWrite(ownedKey string, setOn func(base *Config)) error`。`ownedKey`＝该 caller 拥有的点分键路径；`setOn`＝把那枚键盖到传入 base 上（`:84-88`）。**caller 持 `m.mu`**。
- 出口 (1) 文件读不回⇒拒写并报错：`os.Stat` 成功但 `readConfigFile` 失败⇒`observe.Wrap(observe.ClassConfig, err, "…refused to write config.toml; … cannot be read back…")`（`:118-123`）；`Stat` 既非 nil 也非 missing⇒另一条 `observe.Wrap`（`:128-131`）。均**一字不写**。
- 出口 (2) 已同值⇒不写：`disk!=nil` 时 `written=diffKeyPaths(disk, base)`，`len(written)==0`⇒`slog.Info("…already carries this key; nothing written")` `return nil`（`:149-155`）。
- 出口 (3) 真写了⇒报告差异键：`SaveFile(m.path, base)`（`:160`）后 `diverged=diffKeyPaths(m.cur, base)`（`:166`）；空⇒`slog.Info(...wrote...)`＋`m.statOwnWrite()`（`:168-176`）；非空⇒`slog.Warn(... kept_in_file_not_in_memory ...)`、**不 adopt stat**（`:180-186`，AC#3 收窄）。

**`diffKeyPaths` 那侧能不能报告"哪些键被写了、哪些外来手改被保留"**
- **能**。`diffKeyPaths(a,b)`（`:201`）走序列化同一套 struct tag 列点分叶子键路径、排序。`written:=diffKeyPaths(disk, base)`＝"这次写实际改了文件的哪几枚键"（**对文件内容比、不对内存比**，`:145-147`）；`diverged:=diffKeyPaths(m.cur, base)`＝"文件现在有、本进程内存没有"＝**被保留的外来手改**（`:163-166`）。⇒ AC#11-2 要的"回执逐枚列落盘键路径"直接取自 `written`。
- 例外跟序列化器一致：`SchemaVersion` 跳过（`:220-223`）、`toml:"-"` 跳过、`Config.Plugins` 手工在 `plugins` 下走（`:194-195/:244-245`）；map 逐 key 命名（`[plugins.<id>]`/`[llm.providers.<name>]`，`:284-314`）；报告有 `maxReportedKeys=24` 上限并显式标 `(+N more)`（`:78/:330-336`）。这份报告被 `writeguard_226_test.go:117 TestAC2GuardedWriteReportsTheKeysItWroteAndTheOnesItKept` 钉着（断 `wrote=[fs.allowed_dirs]` 有、`wrote=[ball.size` 无、`kept_in_file_not_in_memory=[ball.size]` 有、`level=WARN`）。

**配置管理器今天哪几处已在用它（照抄形状）**
- 三枚全走 `mergeWrite`：`permmode.go:79`（`SetPermissionMode`，键 `keyRiskPermissionMode="risk.permission_mode"`）；`allowdirs.go:142`（`writeAllowedDirs` 被 `AddAllowedDir:82`／`SetAllowedDirs:124` 复用，键 `keyFSAllowedDirs="fs.allowed_dirs"`）。键常量声明 `writeguard.go:69-72`。⇒ 新增"配置写 setter"照这三枚形：先改内存、调 `mergeWrite(ownedKey, setOn)`、失败把内存回滚（`permmode.go:82`／`allowdirs.go:143` 都是这么回滚的）。
- ⚠ `SaveFile`（`loader.go:238-249`）＝`deepCopyConfig`→强制 `SchemaVersion`→`MarshalCanonical`→`atomicWrite`，**不读回文件、不跑 validate**——正是票 226 要封的洞（`writeguard.go:5-17` 头注逐字，且 `writeguard_226_test.go:102 TestAC1ControlSnapshotWriteIsWhatRevertsTheHandEdit` 把它当"反面教材"钉住）。AC#11 明禁落地腿用它。
- ⚠ 已知代价（`writeguard.go:45-52` 头注＋`always_write_no_clobber_226_test.go:81` 注）：`MarshalCanonical` 会**整文件重排成 canonical 形式**，注释与手挑键序丢失（**值不丢**、格式丢）。AC#11 的"逐枚键各一发受保护写"仍会重排整文件，这是票 226 已接受的代价，非本票新洞。

**`validate()` 在不在写路径上**
- **不在。** `validate` 定义 `internal/config/validate.go:24`；调用仅两处＝上游读管道 `loader.go:123`（`readConfigFile`）与 `migrate.go:77`。`writeguard.go`／`mergeWrite`／`SaveFile` 全程无 `validate`（本腿 `grep validate internal/config/writeguard.go`＝零命中复认）。⇒ AC#7"写前必校验"是**真缺口**；落点（AC#11-3）＝新增导出 setter 内、进 `mergeWrite` 之前跑 `validate()`（⛔ 不改 `mergeWrite` 本体，会波及既有三枚 setter）。
- 现成校验子函数（落地腿按设置项选用）：`validateRisk:109`／`validateModels:95`／`validateModelSpec:195`／`validateCost:117`／`validateRealtime:124`（`catalog.go`）／`validateAPIKeyRefs:152` 等，对应设置页那几项（模型名／上下文窗口／价格／权限档位）。⚠ 权限档位（`risk.permission_mode`）按 J5 **不许由通用 config.set 写**（见 §⑦-5），故那一项的 setter 不进本票。

## ④ 会被叫红的现成钉（逐枚点名＋断言原文＋判会不会红）

落地腿会新增／改动的符号（据 §②③）：`MethodConfigGet`/`MethodConfigSet`（值 `config.get`/`config.set`，J6 已定名，**不带 `panel.` 前缀**）、`ComposerDispatch.Config` 槽 + `ConfigRequestHandler` 接口 + `HandleConfigRequest`、若干 `Manager.Set*`（走 `mergeWrite`）、`ComposerState` 里"凭据已录入/配置可读"那一维（AC#2）。

> 词面预检：`grep -rn "config.get\|config.set\|MethodConfig\|HandleConfig" --include=*_test.go cmd/ internal/` ⇒ **零命中**（这些名字今天不存在）。**词面零命中≠没钉**——下面 T1–T8 全是**行为型钉**（只写产物／数名册），逐枚读结果如下。除 T3/T8 是冻结件外，其余均为本腿真读到的〔量〕。

**T1 `internal/panel/git_test.go:385/:517`（C17 名册数死四枚）——〔量〕结论：J6 命名下不红**
- 抽取器 `whitelistMethodsFromSource:407` 用 `panelMethodRe = regexp.MustCompile(\`"(panel\.[a-z0-9_.-]+)"\`)`（`:394`、`:415` `FindAllStringSubmatch`）——**只数 bridge.go 里 `panel.` 前缀的字面量**。
- 断言 `:385` `if got := whitelistMethodsFromSource(.../bridge.go); !equalStrings(got, wantMethods)`（`wantMethods`＝那四枚常量，`:381-384` 排序）；`:386` 文案"the C17 panel.* whitelist in bridge.go = %v, want the four methods this ticket may not extend"。正控 `:514` 植第五枚 `MethodWorktreeSwitch = "panel.worktree.switch"` ⇒ 数到 5；`:517` 真 bridge.go `len(real) != 4`。
- **判红**：J6 定名 `config.get`/`config.set` **非 `panel.` 前缀** ⇒ 新常量值不被 `panelMethodRe` 匹配 ⇒ `whitelistMethodsFromSource` 仍返那四枚 ⇒ **`:385/:517` 不红**。**只有把新名起成 `panel.config.*` 才红。** ⇒ 前稿 §⑦-3 那条"要不要惊动票 181 的 git_test.go"由函数体坐实＝**J6 命名下不惊动**。⚠ 副作用：git_test 这道"名册只四枚"的门天然管不到 `config.*`——它守的是票 181 的 panel.\* 读面，不是本票的白名单增长（本票的白名单增长由 T3 那枚 answered 尺寸钉守）。

**T2 `internal/panel/composer_test.go:409/:502`（渲染器只一枚门）——〔量〕结论：Go 侧独立不红，是两界登记点**
- `composerRouteLiterals():409-419` 是**硬编码集合**＝那四枚 `Method*` 常量 + 写死的 `"panel.approval.request"`，**不是**从 bridge.go 常量自动派生 ⇒ 落地腿给 bridge.go 加常量**不会自动进这张 allowed 表**。
- `TestTheRendererHoldsExactlyOneDoorToTheHost:502` 拿 `scanRendererHostDoors` 扫 **`frontend/src`**（`:504`，跳过 `node_modules/.git/dist` `:434`）的 route 字面量，`allowed[lit]` 不中即记 `rep.unknown`；`:521` 红句"the renderer names a route the Go side does not answer"。它看的是页面文件，不看 Go 侧常量。
- **判红**：Go 侧加 config 方法**单独不使此枚红**。**它是"页面一旦能叫 config.get"那一步的登记点**：等界面侧让页面调用新 route，须有人往 `composerRouteLiterals()` 补那枚名（`:406-407` 逐字"Growing it is a two-sided change on purpose"），否则此枚红。⇒ 归 AC#1 的两界接缝，非 Go 独立挡路项。（⚠ `routeLiteralRe` 只匹配 `panel.` 还是全 route，本腿未读到其定义，标〔推〕——落地腿现读；若只匹配 `panel.`，连这层接缝都不触发。）

**T3 `internal/panel/l2_grant_boundary_test.go`（FROZEN）——〔只点名＋引编排者 J1〕结论：无论命名必红，是 J1 已授权那两处锚**
- 本腿 `grep -rn "len(sortedSet(pkg.answered))" l2_grant_boundary_test.go` 复核到 `:2139 if got := len(sortedSet(pkg.answered)); got != 4` **确实存在**；`pkg.answered` 由"跑起来的 guard 应答哪些 route"（＝`knownComposerMethod` case 列）推出（该文件结构，本腿不引其内部文字当凭据）。
- 依 **J1（台账 `A487`）** 现量两枚锚：`:2051` 拼写钉 `guardAnchor := "case MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend:"`；`:2139` 尺寸钉。
- **判红〔推，基于 answered 集＝knownComposerMethod〕**：AC#1 要求 parse 应答 config 方法 ⇒ 必须进 `knownComposerMethod` ⇒ `pkg.answered` 从 4→6 ⇒ `:2139` `!=4` 直红；`:2051` 那枚 `case …四枚…` 字面也变 ⇒ 直红。**这是本票唯一"config.\* 还是 panel.config.\* 都躲不开"的冻结钉**，与命名无关，与"改法"有关——J1 已裁＝甲（人工批准、只动那两处、只许变强、把"数 4"换成能力形＋自带正控）。
- **另一维判不了（→ §⑦-2）**：`answeredGrantDoor`/facet-2 的"approval-shaped/grant-door 探测器"会不会把 `config.set`（写配置键）认成 grant 门——规则在冻结件内部，禁令不让引其内容当凭据。

**T4 `cmd/wisp/panel_inbound_33_test.go`（票 33 slice B）——〔量〕结论：不红（修正前稿〔推〕）**
- `TestAC9ComposerDispatchHasAProductionCaller:182` 是 go/AST 扫：要求本目录某**非测试文件**把 `*panel.ComposerDispatch` 绑到名并 `<name>.Handle(...)`（`inboundHandleCallSites33:278`）；红句 `:185`"no NON-test file in cmd/wisp sends Handle to a *panel.ComposerDispatch"。
- 另四枚跑真 leg：`TestAC9InboundLegFromStdinReachesTheWriteLeg:99`（喂 `panel.mode.request`、断档真落盘）、`TestAC9InboundLegRefusesUnlistedMethodAndAuditsIt:125`（植 `panel.review.allow`、断被拒＋`INBOUND-DISPATCH` 审计＋**未达 mode handler**＋档不变）、`TestAC9InboundLegRefusesRosterMethodWithNoHandler:158`（`panel.workspace.request` ⇒ stdout 含"处理器未接入"）、`TestAC9InboundFlagSurfaceIsNarrow:387`。
- **判红**：这些**不钉"四槽里恰好三枚 nil"**——钉的是"有生产调用者"＋"`panel.review.allow` 不在名册⇒拒"＋"workspace 空槽⇒响亮拒"。落地腿给 `ComposerDispatch` 加 `Config` 槽、并在 `newPanelInboundDispatch` 里赋非 nil config handler ⇒ **不触上述任一枚**（只要 `panel.review.allow` 仍不在名册、workspace/attachment/message 仍 nil）。⇒ 前稿"可能红"作废，坐实＝**不红**。

**T5 `cmd/wisp/leg_dispatch_gate_133_test.go`（票 133 派发闸门）＋`internal/panel/composer_dispatch_test.go:387`——〔量〕结论：不红（只要不新增 `wisp <子命令>`）**
- 此闸门枚举 **`func main` 的 argv 分支**（`enumerateLegs133:1002`：argv switch 的 case 字面标签＋if 比较 argv 的分支），地板 `minLegs133=11`（`:148`）；`censusVsUsage133:1707` 双向对账 main.go 的 `usage` 常量（被派发的腿都要在 usage 里有一行、usage 里的命令都要被派发）；每枚被派发腿须被 nail/test/ruling 覆盖（`coverageReds133:1377`）。
- **判红**：本票 config 方法进的是 `panel.ComposerDispatch`/`bridge.go`，**不是 `func main` 的新 case 标签**；`panel-inbound` 那枚腿已存在、已被 panel_inbound_33 覆盖 ⇒ **腿名册不变 ⇒ 不红**。⚠ 唯一会红的做法＝落地腿新加一枚 `wisp xxx` 子命令（那要同时补 usage 行＋覆盖）——本票不该这么走（配置写走 panel inbound）。
- `composer_dispatch_test.go:387`"constant: this is a second whitelist, and the guard's own case list would never have named it"：只在 dispatch 里写**裸 route 字面量**另立名册才红；照现有四枚形（switch case 引用 bridge.go 导出常量）⇒ 不红。

**T6 `internal/panel/frontend_hygiene_test.go`——〔量〕结论：属界面侧射程，Go 独立枚不红**
- 四把尺全扫 **`frontend/`**：`TestPanelFrontendIsStateless:169`（禁 localStorage/indexedDB/cookie/serviceWorker/…，`bannedStorageAPIs:48`）、`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme:196`（色值只在 `tokens.generated.css`）、`TestFrontendHasNoEmoji:246`（`emojiRangesRe:71`＝d22scan ban #8 字符类副本、**无注释豁免**）、`TestFrontendNeverNamesAnApprovalDecision:281`（`approval.decide` ban #6）。
- **判红**：红不红取决于**界面侧那枚 agent 把设置页写成什么样**，不取决于 Go 侧 config 方法 ⇒ 归 AC#9/界面依赖那一格；本票 Go 独立枚不触它。⚠ `TestFrontendHasNoEmoji` 的字符类起点 `1F000`、含 `2200-22FF`/`2600-27BF`——与 `AGENTS.md §1.2` 记的"仪器实际扫"一致，是**界面侧**文案要避的雷（`→` 界面侧躲得过、`✓`/`≤` 躲不过），与我们这侧无关。

**T7 `internal/config/writeguard_226_test.go`＋`cmd/wisp/always_write_no_clobber_226_test.go`（票 226 不覆写）——〔量〕结论：不红（修正前稿〔推〕）**
- 两枚都**不数 setter、不枚举 Manager 方法**：前者逐枚测 `AddAllowedDir`/`SetPermissionMode`/`SetAllowedDirs` 的**产物**（手改别处键活不活、报告写没写、失败回没回滚、已同值不重写）；`TestDiffKeyPathsNamesRealKeyPaths:416` 对**显式改动的** `Ball.Size`/`LLM.Providers`/`Plugins`/`App.Portable` 断 `want=["ball.size","llm.providers.acme","plugins.p1"]`（`:428`，`App.Portable` 因 `toml:"-"` 不报）。后者端到端驱 `wisp run` 的长期写（`AddAllowedDir` 落盘＋ `[app] theme` 手改存活）。
- **判红**：新增走 `mergeWrite` 的 `Manager.Set*`（写**既有** config 字段，如 llm/model/上下文窗口）⇒ 不改 `mergeWrite` 本体、不新增 Config struct 字段（设置项都是既有字段）⇒ **不红**。⚠ 即便某项要新加 Config struct 字段，`TestDiffKeyPathsNamesRealKeyPaths` 用 a/b 显式对照、新字段两边同默认 ⇒ 仍 len 0，不红；本腿读过 `writeguard_226_test.go` 全文，未见"数 setter 枚数"式断言 ⇒ 无此雷。
- ⚠ **但 AC#11-2 的"逐枚键各一发受保护写"落地的新用例**会**新增测试**（这属落地腿写的东西、不属挡路钉）；票 226 既有测试是"射程参考"，不挡新 setter。

**T8 冻结件 `internal/perm/ticket90_persist_test.go`＋`internal/panel/tokens_fourway_test.go`——〔冻结，只点名〕**
- `ticket90_persist_test.go` 射程＝permission_mode persist（`permmode.go:70` 已走 mergeWrite）。⇒ 落地腿**不动** `SetPermissionMode` 本体、只为其它键加 setter ⇒ 预期不红；若把 `permission_mode` 划进通用 `config.set` 可写键集 ⇒ **正面撞 J5（甲：显式拒写 `risk.*` 全族）**，须走带 L2 卡那条（`cmd/wisp/run.go confirmModeSwitch`＋`panel.ModeWriteHandler`）、不走配置写方法。`tokens_fourway_test.go` 射程＝token 四态，与设置路由无直接交集（预期不红，〔推〕）。

**§④ 净结论（会叫红的只有两格，且命名决定其中一枚的射程）**：
1. **T3（冻结 `:2051/:2139`）——无论命名必红**，须按 J1/A487 的解冻改法处理（改成能力形＋正控、不许删钉或放宽成 ≥4）。这是 AC#1 落地时**第一硬碰**。
2. **AC#2 那一维——`composer_test.go:48` 两界对账（非本票冻结件）会把 `internal/panel` 测试洗红**，直到界面侧补 `frontend/src/lib/panel.ts` 的 `ComposerState`；⇒ AC#2 的 Go 半枚今天做不到 self-green，必须与界面同批。
3. 其余 **T1（J6 命名下）/T2/T4/T5/T7/T8 对"Go 侧独立写"都不红**；T6 只在界面侧红。

## ⑤ 我跑了哪些尺、每条真实读数（全是只读尺；零 CPU／零改盘）

- `ls .scratch/wisp/probes/248/` ⇒ 只有 `a1`（含 59,623 字节 `census.md`）⇒ **`a2` 未用过，不换号**。
- `ls .scratch/wisp/issues/ | grep 248` ⇒ 一枚 `248-the-panel-has-no-settings-route-the-inbound-whitelist-has-zero-config-or-credential-methods.md`（凭名取准，未猜）。
- `date "+%Y-%m-%d %H:%M:%S %z"` ⇒ `2026-10-01 11:51:49 +0800`；`git log -1 --format=%H` ⇒ `23627aba5d64b09f92ec8687a2716ba0668a6b38`；`git rev-parse --abbrev-ref HEAD` ⇒ `dev`（同发取；追加进度行时刷新 `date` ⇒ `11:57:07 +08`）。
- Read `cmd/wisp/panel_inbound.go`（269 行）⇒ `:218 Confirm:nil`／`:228 Confirm:nil`／`:236 Workspace:nil`／`:237 Attachment:nil`／`:238 Message:nil`；**无 Config 槽**；`:159-173` 受理只打 console。
- Read `internal/panel/bridge.go` ⇒ 白名单常量 `:42-45`（四枚，零 config）；`ParseComposerRequest:84`；`knownComposerMethod:104-110`；`RefusedEnvelopeForUser:115`。
- Read `internal/panel/composer_dispatch.go` ⇒ `Handle:120` 返回 `(string,error)`、受理⇒`("","nil")`（`:130`）；`dispatch:137-162` switch 四 case＋default `rosterMismatch:167`；`ErrNoHandlerAttached:62`、`ErrRosterMismatch:178`；头注 `:18-25` 说写裸 route 串会触发本包 AST 钉。
- Read `internal/panel/composer_handlers.go` ⇒ `ModeWriteHandler:86`、`HandleModeRequest:111`（`Confirm:nil`⇒变宽拒 `:133-138`）。
- Read `internal/config/writeguard.go` ⇒ `mergeWrite:111` 签名＋三形（`:118-131` 读不回拒／`:149-155` 同值不写／`:160-186` 真写报告）；`diffKeyPaths:201`；键常量 `:69-72`；`maxReportedKeys=24 :78`。
- Read `internal/config/permmode.go`／`allowdirs.go` ⇒ 三枚先例 `SetPermissionMode:70`/`:79`、`AddAllowedDir:62`、`SetAllowedDirs:109`、`writeAllowedDirs:139`→`mergeWrite:142`；`statOwnWrite:157`；回滚 `permmode.go:82`／`allowdirs.go:143`。
- Read `internal/config/loader.go` ⇒ `SaveFile:238-249`（serializer，无 re-read、无 validate）；`readConfigFile:64`→`validate:123`。
- `grep -n "func validate|validate\(" internal/config/*.go` ⇒ `validate` 定义 `validate.go:24`；调用仅 `loader.go:123`＋`migrate.go:77`；**`writeguard.go` 内零 `validate`**。
- `grep` 名册／尺寸钉（`answered|sortedSet|composerRouteLiterals|whitelist|len(...)` 在 `internal/panel/*_test.go`）⇒ 命中 `git_test.go:385/514/517/518`、`composer_test.go:406/409/426/489/502/525`、`l2_grant_boundary_test.go:2139`、`composer_dispatch_test.go:387`。
- Read `internal/panel/git_test.go:355-528`〔定死 T1〕⇒ `panelMethodRe:394`＝`"(panel\.[a-z0-9_.-]+)"`、`whitelistMethodsFromSource:407-423` 只抽 `panel.` 字面量；`:385` equalStrings、`:514/:517` 正控。
- Read `internal/panel/composer_test.go:1-115`＋`:400-539`〔定死 T2 与 AC#2 对账门〕⇒ `composerRouteLiterals:409-419`（硬编码四常量＋`panel.approval.request`）；`TestTheRendererHoldsExactlyOneDoorToTheHost:502` 扫 frontend/src；`TestComposerContractTypesMatchFrontend:48-80` pairs 表含 `ComposerState↔ComposerState`、`:73/:74` 双向报"Go emits key TS doesn't declare"；`composerModeWriteRe:88-90` 门 2。
- Read `internal/panel/pump_test.go:100-139`〔定死 AC#2 顶层键钉〕⇒ `:123` marshal 后 JSON 键必须＝`"composer,generatedAt,pending,results"`（四枚、加顶层第五键即红）。
- Read `internal/panel/composer.go:40-133`〔定死快照形状〕⇒ `Snapshot:57-92`（4 恒发＋`Instructions`/`Tasks` 两可选 ptr+omitempty）、头注 `:45-56` 记"四把钉"（`composer_test.go:48`/`approval_test.go:105`/`pump_test.go:111-124`/`:270-276`）、`Tasks` 先例 `:86-91`。
- Read `cmd/wisp/panel_inbound_33_test.go` 全文〔定死 T4〕⇒ AST 只要求生产 `.Handle` 调用者存在（`:182`）＋ `panel.review.allow`（`:130`）被拒＋ `panel.workspace.request`（`:160`）"处理器未接入"；**不钉 nil 槽枚数**。
- Read `cmd/wisp/leg_dispatch_gate_133_test.go` 全文〔定死 T5〕⇒ 枚举 `func main` argv 分支（`enumerateLegs133:1002`）、`minLegs133=11`、`censusVsUsage133:1707` 对账 usage；config 方法不进 main switch ⇒ 腿名册不变。
- Read `internal/panel/frontend_hygiene_test.go` 全文〔定死 T6〕⇒ 四尺全扫 `frontend/`（stateless/色值/emoji/`approval.decide`），`emojiRangesRe:71`＝d22scan ban#8 副本。
- Read `internal/config/writeguard_226_test.go`＋`cmd/wisp/always_write_no_clobber_226_test.go` 全文〔定死 T7〕⇒ 两枚都逐枚测产物、不数 setter；`TestDiffKeyPathsNamesRealKeyPaths:416` 用 a/b 显式对照。
- `Glob cmd/wisp/*_test.go`（47 枚）／`internal/config/*_test.go`（12 枚）⇒ 点名 T4–T8 相关。
- `grep -rn "config.get|config.set|MethodConfig|HandleConfig" --include=*_test.go cmd/ internal/` ⇒ 零命中（词面预检，已升级为逐枚读行为钉）。

## ⑥ 我可能写错的条目（对抗我自己）

1. **T3 的"必红"是〔推〕、不是〔量〕**：我坐实了 `l2_grant_boundary_test.go:2139` 那行存在（`grep` 复核），也坐实了"answered 集由 `knownComposerMethod` 推出"是该文件的既定结构（J1/A487 已就此裁过、`composer_dispatch.go:19-25` 头注也指认该 AST 尺），但**我没读冻结件内部**去确认 `pkg.answered` 到底怎么填——若它对 `config.*` 非 `panel.` 前缀名另有一套过滤（只认 `panel.` 入向），则"4→6"这一步可能不成立、T3 可能不红。⇒ 定性标〔推〕，落地腿/验收腿务必现读该件坐实"加 config 方法后 `:2139` 到底红不红"。（即便〔量〕坐实不红，J1 已授权那两处锚的改动仍要做，因为 `:2051` 拼写钉与派发表 case 是两回事。）
2. **`composerModeWriteRe`/`composerRouteLiterals`/`routeLiteralRe` 三者的匹配前缀我只读了两个**：`composerModeWriteRe:88-90`（含 `panel.mode.set`/`permission.mode.set`/`PermissionMode =`）与 `panelMethodRe`（`panel.` 前缀）我坐实了；但 `composer_test.go` 的 `routeLiteralRe` **定义我没读到**——T2"Go 侧独立不红"依赖"该尺只扫 frontend 文件、不扫 Go 常量"（这点坐实），但"界面侧叫 config.get 会不会经 routeLiteralRe 触发两界接缝"取决于 `routeLiteralRe` 是否只匹配 `panel.` ⇒ 标〔推〕。
3. **我把 AC#6 塞进了 §① 逐格分类**：四道尺本只读腿没资格判红/绿，写成"交落地腿自跑"是把它剔出射程，但占了分类表一格——**读者别以为我核过门禁**。
4. **A497（Go→页面投递队列无取出者）我只登记、没据此下结论**，符合派单；但 §② 末段仍可能被读成"回执通道做不了＝路由做不了"。**更正**：路由（名册＋dispatch）与"把回执投回页面"是两件事，前者 Go 侧可做（AC#1），后者是票 33/35 的欠账。
5. **AC#2 我判"两栖"偏乐观**：`composer_test.go:48` 双向对账意味着**只要动 `ComposerState` 子键，`internal/panel` 整包 `go test` 即红**，AC#2 的"Go 半枚"不是"能做但没做完"，而是"做完会让本包测试洗红"——除非界面侧同批改 panel.ts。我在表里已改成"Go 半枚今天做不到 self-green"，但仍要防读者把它读成"Go 先写完放着、回头再补界面"（那样 tree 一直是红的，违 AC#6 门禁）。
6. **panel_inbound.go / 各 cmd-wisp 测试的行号可能因 33-r4 并发改 `cmd/wisp` 而漂**：`panel_inbound.go` 是非测试文件、相对稳，但 `panel_inbound_33_test.go`/`leg_dispatch_gate_133_test.go` 是 `cmd/wisp` 的 `*_test.go`——派单明令"cmd/wisp 行号落地现读为准"，我已逐处标注。
7. **J6 命名（config.\* 非 panel.\*）是我整套 T1 结论的支点**：若编排者最终改判用 `panel.config.set` 那形，T1（`git_test.go:385/:517`）**立刻翻红**、票 181 地界要同批改。我据 J6 现文（`config.get`/`config.set`）判"不红"，这是**条件结论**，不是无条件。

## ⑦ 判不动的地方（逐条甲／乙／不做＋现量＋为什么判不了）

1. **AC#10 的 ⓘ/ⓑ 二选一**——**不做（选形交编排者，派单明令本腿不许选）**。现量：票面 §2c 给 `approval/gate.go:139-145` 钳位、`approval/queue.go:107/115/122`（300s/3s/Max3s）、`resident_approval_windows.go:109`（`approval.New` 建在会话账本前⇒`Options.Grants` nil⇒`GRANT-DROPPED`）。ⓘ 移 `approval.New`＝改票 246 AC#1 乙形次序、需先落 `A##`；ⓑ 不移动＝回执文案逐字声明。**为什么判不了**：选择权与 `A##` 批准记录归编排者，本腿只读、不选、不排落地顺序；且 ⓑ 要落"页面可见回执"⇒ 另一半等界面。
2. **`config.set` 会不会被 `l2_grant_boundary_test.go` 的"approval-shaped/grant-door 探测器"误判成 grant 门**——**判不了（甲/乙/不做三形摆不出来，因证据在冻结件内部）**。现量：`grep` 只坐实 `:2139` 尺寸钉存在；grant-door 匹配（`answeredGrantDoor`/facet-2）在冻结件内部，⛔ 禁令不让我引其内容当凭据。**为什么判不了**：要看那枚探测器按什么词/结构认"grant 门"。⇒ 交编排者（或让非实现验收腿自己取数）。若它红，走 J1 解冻那一族改法、别放宽钉。
3. **AC#4／AC#9 今天能否勾**——**不做**。现量：`grep PanelManager main.go/resident_windows.go/run.go`=0 命中（票面 §5.4）、`git ls-files frontend`⇒只 `frontend/dist/.gitkeep`（票面 §2b 末）；宿主未接进常驻、页面产物未填。J9 另裁 AC#4 改按"重启后真用上"，重建那一跳是票 223 地界。**为什么判不了**：依赖票 33 宿主落地＋界面侧产出 dist；两层禁令禁我读 dist 内容。
4. **新 setter 的可写键集是否含 `allowed_dirs`（`fs.*`）**——**判不了（触 J5 边界）**。现量：J5 裁甲＝配置写方法**显式拒写 `risk.*` 全族**（`risk.permission_mode` 只能走带 L2 卡那条）。但 `fs.allowed_dirs` J5 没逐字点名，是否同样禁由通用 `config.set` 改、还是允许走 `Manager.AddAllowedDir/SetAllowedDirs` 那两枚既有受保护写——是编排者的推广边界。**为什么判不了**：`fs.*` 属"权限档相关但非 permission_mode 本体"，J5 只锁了 `risk.*`。⇒ 本腿只坐实"至少 `risk.*` 全族必拒"。
5. **（已消解）J6 名的 `config.*` 会不会惊动票 181 的 `git_test.go`**——**现由函数体坐实＝不惊动**（T1 只抽 `panel.` 前缀字面量）。⇒ 从"判不了"降级为"有结论的条件项"：只要编排者不改判命名（仍是 `config.get`/`config.set` 而非 `panel.config.*`），票 181 那两枚名册钉就不红；一旦改判用 `panel.` 前缀，回到 T1 红。

## 交件判语

- **射程**：AC#1..AC#11 逐格分类（§①）＋路由承载点／写侧落地现量（§②③，带 `file:line`〔量〕）＋撞钉预检八组尺逐枚点名＋断言原文（§④，除 T3/T8 冻结件外均本腿真读〔量〕）。核心结论：**今天可 Go 侧写腿独立做完并交凭据＝AC#3/AC#5/AC#7/AC#11，加 AC#1**（AC#1 Go 侧可做，但落地即撞唯一那枚"无论命名必红"的冻结钉 T3，须走 J1 解冻改法）；**两栖、Go 半枚做不到 self-green＝AC#2**（被 `composer_test.go:48` 两界对账门洗红）**／AC#8**（页面可见那一半等界面）；**必须等界面侧／票 33 宿主＝AC#4/AC#9**（及 AC#10 的"页面可见"那一半）；**AC#6 交落地腿自跑（本只读腿无资格判门禁）**；**AC#10 选形、grant-door 探测、`fs.*` 边界——一律交编排者裁（§⑦，未替其下结论）**。
- **没答的（不写成答了）**：① T3"必红"是〔推〕——我坐实了 `:2139` 那行存在与"answered 集走 knownComposerMethod"的结构，但没读冻结件内部确认加 `config.*` 后 `:2139` 是否真从 4→6，须非实现腿现读坐实（§⑥-1）；② `routeLiteralRe` 定义未读⇒T2 的界面接缝触发条件是〔推〕（§⑥-2）；③ 冻结件 T3/T8 内部匹配语义按禁令**未引**⇒只点名。
- **最需要编排者裁的那一格**：**§⑦-2（grant-door 探测器会不会把 `config.set` 误认成 grant 门）**——它决定 AC#1 落地时那枚冻结件（T3）是"只按 J1 改尺寸/拼写两枚锚"就够，还是连"新增的 config 方法名会不会被 facet-2 判成授权门"也要一并裁；这一格只有能读冻结件内部的非实现腿/你本人能定。次一位：**§⑦-1（AC#10 选形）**与 **§⑦-4（`fs.*` 是否禁由 config.set 写）**。
