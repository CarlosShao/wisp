# 票 248 · 落地腿 `248-r1` 证据件 — 面板「设置」那条路：从入向请求一路修到磁盘

> **本腿代号**：`248-r1`（写码腿，Go 半边）。票面＝`.scratch/wisp/issues/248-the-panel-has-no-settings-route-the-inbound-whitelist-has-zero-config-or-credential-methods.md`（**只读，一字未改**）。
> **票面 AC 框一枚未碰**（勾与不勾归编排者与验收腿）。本件引用的票号只有两枚：**票 248**（本票）与**票 33**（面板宿主）；其余票号一律不引（派单射程）。
> **前人料一律当待验断言处理**：`248-a1`／`248-a2` 的普查与 `ticket248-ruling-248a2` 的判语里出现过的行号与指认，本件只写**我自己复认过的那一发读数**，取数时刻逐条带在旁边。
> **凭据纪律**：本件不含任何凭据值，只写字段名／变量名／ref 形状。哨兵串是本腿自造的假值，逐字写在本件 §4，⛔ 它不是也不可能是任何真实密钥。
> **两层禁令**：⛔ 未读未写 `frontend/**`／`design/**`；本件里关于页面产物的唯一事实来自 `git ls-files`（名册口径），不是内容。

---

## 0. 起手锚（同发取，逐字）

| 尺 | 读数 |
|---|---|
| `date '+%Y-%m-%d %H:%M:%S %z'` | `2026-10-02 08:25:26 +0800` |
| `git log -1 --format='%h %ad %s'` | `a35f7f5e Fri Oct 2 00:04:33 2026 +0800 ledger(A511)…`（分支 `dev`） |
| `git status --porcelain` | **非空，且永远非空**＝owner 在 `design/**` 下的未提交移动与未跟踪件＋几十枚他人台件。起手名册已整份存档：`.scratch/wisp/probes/248/r1/start-anchors.txt`（**277 行／12,368 字节**） |
| 闸门口径 | ⛔ 不是"终态必须为空"，而是**终态名册 ＝ 起手名册 ＋ 只本腿自己那几枚路径**（派单明写）。本腿写面：`internal/panel/**` · `internal/config/**` · `cmd/wisp/**` · 本件 · `.scratch/wisp/probes/248/r1/**` |
| `internal/panel/**`·`internal/config/**`·`internal/secret/**`·`cmd/wisp/**` 起手脏度 | `git status --porcelain -- internal/panel internal/config internal/secret cmd/wisp` ⇒ **0 行**（派单 08:2x 的量，本腿同法复量一致） |
| `ls .scratch/wisp/probes/248/` | 只有 `a1`／`a2` ⇒ **`r1` 未用过**，无需换号 |

**本腿不做的两件事（具名）**：⛔ 不动票面 AC 框；⛔ 不碰 `internal/ball/**`（同机另一枚腿在写）。

---

## 1. AC#0 三问（题面明写"不许直接开写"⇒ 本节先于任何产码落盘）

> 三问的答案全部带**现读凭据**（`文件:行` ＋ 尺读数 ＋ 时刻）。三问定完之后才开始改代码；若 a1/a2 的料不足以定这三问，本腿的处置是**停手上报**，不是自己发明新范围。

### 1.1 Q1 新增哪几枚入向方法名、各自参数字段与错误形状

**名（⛔ 不新造，取规格里已有的两枚，J6 已裁）**：`config.get`（读）与 `config.set`（写），共 **两枚**。

- 现量·规格那行确实写着这两个名：`grep -n "config\.get\|config\.set" docs/specs/SPEC-08-ui-ball-panel.md` ⇒ **`:168`**：`| \`config.get\` / \`config.set\` | invoke | 安全节放宽走 L2 重新确认（SPEC-03 §4.2） |`（本腿 `08:31:49` 现跑；**只读该行，一字未改 `docs/specs/**`**）。
- 现量·今天白名单确实只有四枚、零枚 config：`grep -rn "panel\.[a-z]*\.[a-z]*" --include=*.go internal/panel/ | grep -v _test | grep -o … | sort -u` ⇒ 恰 `panel.attachment.add`／`panel.message.send`／`panel.mode.request`／`panel.workspace.request`（本腿 `08:31:49` 复跑，票面现量 #1 成立）。真身＝`internal/panel/bridge.go:41-46`，守卫＝`bridge.go:104-110` 的 `knownComposerMethod`。
- **命名两禁我照做并机读**：`config.set`／`config.get` 都不含那十枚裁决词（`approve/approval/grant/allow/permit/ratify/authorize/authorised/decide/decision/verdict`）、不落 `panel*` 前缀族 ⇒ 冻结件的授权门探测器不该把它认成 grant 门。**这一条我不信判语，本腿种正控**（§1.2 Q2 的 P-命名一发）。
- **第三枚凭据名不新增**：`A484` 的射程被 J6 收窄为"只有这两枚名"⇒ 凭据录入**不另起 method**，它是 `config.set` 的一枚**具名字段**（下条）。

**参数字段（具名字段枚举，⛔ 不是任意 key 透传）**：共享封套 `ComposerRequest`（`bridge.go:69-79`，本腿现读：`Method/RequestID/Source/To/Path/Text/AttachmentPayload/Attachments`）只长**非值**的三枚选择子＋一枚具名字段名：

| 新增字段（都 `omitempty`） | 作用 | 为什么它不算"任意 key 透传" |
|---|---|---|
| `configField` | 具名字段枚举（下表 5 枚＋凭据枚） | 枚举在**代码里**列死，不在页面上；未列名一律拒 |
| `configProvider` / `configModel` | 定位**已存在**的 provider / 目录条目 | 只是选择子；写不出新键路径，非 `section.key` 字符串 |
| `configValue` | 该字段的值（字符串承载，按枚举解析类型） | 只在 `configField` 命中枚举时被读 |

可写枚举（`configField` 的全部合法取值，逐枚对应现成字段与行号）：

| 枚举名 | 落到的结构体字段（现读） | 生效档（代码分档处 `internal/config/manager.go`） |
|---|---|---|
| `provider_base_url` | `Provider.BaseURL`（`internal/config/schema.go:383`） | `[llm]` 名义 hot；实践面不重建（见 Q1-末段） |
| `provider_api_key_ref` | `Provider.APIKeyRef`（`schema.go:387`，只接受 `dpapi:`/`env:` 形状） | 同上；ref 不是秘密，可上页面（`internal/secret/refs.go:70-79` 的自陈） |
| `model_context_window` | `ModelSpec.ContextWindow`（`schema.go:354`，`0 = unknown`） | 同上 |
| `model_price_in` / `model_price_out` | `ModelSpec.Price.In` / `.Out`（`schema.go:337-338`，单位 micro-USD / 1M token，`:335-336`） | 同上 |
| `role_chat_model` | `LLMSection.Roles.Chat.Model`（`schema.go:413` → `:287`） | 同上 |
| `provider_credential` | **不写任何 config 字段**：值交 `secret.NewStore` → DPAPI，落的是 blob，档上写的是它换回来的 `api_key_ref` | 值**不进共享封套**（下段） |

**凭据值不进那枚共用封套（J4 甲，本腿照做）**：`bridge.go:69-79` 是**所有已答方法共用**的请求封套，而它的裁决词表不含 value/secret ⇒ "靠词表拦值"拦不住。⇒ `configField == provider_credential` 那一发的**值不放在 `ComposerRequest` 上**，由配置腿从**原始字节**里解进一封装在 `cmd/wisp` 本地声明的**只写形状**（`credentialWriteRequest`：只有 `provider`＋`value` 两枚 tag），立刻交 `Store.Store`，该结构体不入任何字段、不进快照、不进日志。⇒ 派发那一跳必须能把 raw 递给这枚处理器：`dispatch` 的入向加一枚 raw（本腿写面 `internal/panel/composer_dispatch.go`），其余三枚既有 handler 接口签名不动。

**错误形状（复用既有词汇，⛔ 不造第二套）**：

| 形状 | 现读出处 | 本票哪一发用它 |
|---|---|---|
| `ErrComposerRequest` | `bridge.go:49`，`ParseComposerRequest:84-102` 每一支都带 reason＋requestId | 名不在名册／source 不对／无 requestId |
| `ErrNoHandlerAttached` | `composer_dispatch.go:62`、`unattached():182-187` | `config.*` 在册而这枚装配槽是 nil（响亮拒，AC#1 要的那一声） |
| `ErrRosterMismatch` | `composer_dispatch.go:171-178` | bridge 与派发表不一致（default 支） |
| `observe.ClassConfig`＋`observe.Wrap` | `internal/config/permmode.go:83-85`、`writeguard.go:120-122` | 写前校验不过／受保护写失败（D37 的传播形状） |
| 显式拒写 `risk.*`／`fs.*` | 本腿新增，理由逐字引 J5 与 `allowdirs.go:16-24`（"这里的写放宽的锁定段，调用者欠那张 L2 卡"） | 枚举未列 ⇒ 拒，并点名"档位走 `panel.mode.request`、放宽目录走那条欠 L2 卡的路" |

**写路径落在谁手里（票面 AC#0 的第②问）**：`config.Manager` 仍是**唯一**配置写者、`secret.Store` 仍是**唯一**凭据持有者；面板侧只到"发起请求＋被响亮拒绝"。
- 现读·唯一持有者：`secret.NewStore` 的装配点本腿复量 `grep -rn "secret.NewStore" cmd/wisp/ | grep -v _test`（§3 落档读数），写凭据只调 `Store.Store`（`internal/secret/store.go:64`）；DPAPI 那一发 `internal/secret/dpapi_windows.go:24`。
- 现读·唯一写者：受保护原语 `(*Manager).mergeWrite`（`internal/config/writeguard.go:111`，本腿现读到签名逐字 `func (m *Manager) mergeWrite(ownedKey string, setOn func(base *Config)) error`）；既有三枚先例 `allowdirs.go:62 AddAllowedDir`／`:109 SetAllowedDirs`／`permmode.go:70 SetPermissionMode`。
- **依赖边形状（乙形，a1 ②.2 量过的两接法里那枚零增量的）**：`internal/panel` 本地声明接口（先例 `composer_handlers.go:73-78` 的 `ModeWriter`，注释逐字"Declared locally, **with no import of internal/perm**…"），装配根注入 ⇒ panel 不新增对 `internal/config` 的直连边。**这条是本腿的推广还是规矩？**是**接口形状的照抄**，不是新规矩；本件按"照抄现成先例"记，不写成"仓里禁止正向边"。

**读侧怎么答"配了没配"而不泄露值（票面 AC#0 的第③问）**：四态枚举，⛔ 不返回 `Store.Resolve` 的值。
- `internal/secret/store.go:129 Exists(ref)` **只答 `dpapi:`**、`env:` 直接错 ⇒ 两族分开判（a1 ③.1 的这条本腿复认到 `store.go:129` 的签名与 `refs.go:29 ParseRef`）。
- 出的东西只有：枚举名＋计数＋`api_key_ref` 本身（`refs.go:70-79` 自陈随机 id ⇒ ref 不泄漏它是哪家的）。⛔ 不出"末 4 位"：`internal/secret/redact.go` 的 `RedactSecret` 是 CLI 终端出口，不是页面字段（AC#3 那句"不许新增任何把值回显给页面的方法"）。

**生效语义（AC#8 写死）**：`[llm]`／endpoint 类字段**同一次运行内不重建** ⇒ 一律按"**重启后生效**"回执。现量（本腿自己跑，`08:31:49` 那一批）：`grep -rn "hotReloadDisabledPanelInbound" cmd/wisp/ | grep -v _test` ⇒ 常量在 `cmd/wisp/config_reload.go:97`、念它只有一处 `cmd/wisp/panel_inbound.go:207`（**stderr**，不是页面）；重建那一跳本腿⛔ 不加（AC#8 归口）。页面可见性本腿可测（票 33 那侧的回执通道已在：`cmd/wisp/panel_host_windows.go:320` 的 bind 闭包把 `dispatchRaw` 的返回串交回页面，`dispatchRaw` 在 `:544-552` 调 `Handle`）⇒ **回执文案落进 `Handle` 返回串**，这是 Go 侧能做的那一半。

### 1.2 Q2 每一枚的判据长什么样（负向必配正控；⛔ 判据一律问能力）

| # | 判据（机读断言，不是注释词面） | 正控（"种下去一定被抓"） |
|---|---|---|
| P-名 | 白名单与派发表同名册：`config.get`／`config.set` 都在 `knownComposerMethod` **且**在 `dispatch` 的 case 里 ⇒ parse 接受、处理器被调到**次数＝1** | 种一枚未列名 `panel.review.allow`（沿用现成负控名）⇒ parse 必拒、审计一行、处理器 0 次调用 |
| P-在册无处理器 | `config.set` 在册而装配槽 nil ⇒ 响亮 `ErrNoHandlerAttached` | 同一枚用例里另一发把槽接上 ⇒ 必须走到处理器（证明那一响来自 nil 槽而非别的分支） |
| P-值不落盘 | 哨兵值喂进凭据写路径 ⇒ 对**快照字节／`[audit]` 行／持久 sink** 跑"候选密钥形状"正则，**零命中**（常驻用例，不是表里一次性读数） | 同一条尺把哨兵串直接种进快照构造 ⇒ **必须被抓**（证明尺看得到东西，不是恒真） |
| P-只走一枚存储 | 凭据写完后 `api_key_ref` 形状合法、blob 在 `Store` 目录里 ⇒ 且 `config.toml` 逐字节**不含**哨兵串 | 把 setter 换成"把值写进 config.toml" ⇒ `config.toml` 含哨兵 ⇒ 该用例必红 |
| P-写前校验 | 种一发非法值（`model_context_window = "not-a-number"`／越界的档位）⇒ **磁盘逐字节不变**＋`Manager.CheckAndReload` 之后仍读得出旧值＋下一次 `readConfigFile` 起得来 | 现量前提：`grep -c "validate" internal/config/writeguard.go` ⇒ **0**（本腿 `08:31:49`）⇒ 校验今天确实不在这条路上，AC#7 是真缺口；落点＝新导出 setter 内、进 `mergeWrite` 之前 |
| P-按键受保护写 | 并发形状：面板写 A 键时文件里已有进程外手改的 B 键 ⇒ 落盘后 **B 键逐字还在**（写前写后各读一次文件比对，⛔ 不看 rc） | 把 setter 改走 `SaveFile`（整份快照序列化器）⇒ B 键被复原 ⇒ 该用例必红（票 226 的洞重新打开的形状） |
| P-逐枚键各一发 | 多枚字段一次保存 ⇒ 回执**逐枚列出实际落盘的键路径**（取自 `mergeWrite` case 3 的 `written`／`diffKeyPaths`，`writeguard.go:148-151`） | 中途一枚失败即停 ⇒ 回执里已落盘的枚与未落的枚必须不同；把报告压平成"已保存" ⇒ 必红 |
| P-不越界 | `config.set` 的可写枚举**不含**任何 `risk.*`／`fs.*` ⇒ 点名拒写、处理器未被调用、磁盘不变 | 种一发 `configField` 拼成 `risk.permission_mode` ⇒ 必拒且**不落到** `SetPermissionMode`（用调用计数判，不看日志） |
| P-命名两禁 | `config.*` 两枚名不被冻结件的 grant 探测器判成授权门 ⇒ `go test -run 'TestRealGuardRefusesEveryAssemblableApprovalRouteName|…' ./internal/panel/` 读数绿 | ⛔ 不许为了绿去动那三枚表（冻结件）；若它红 ⇒ **停手上报**，本腿不自创解法 |
| P-常驻腿（AC#10） | 只出**读数**不出选择：`[risk]` 那两项对常驻腿到底吃不吃得到，逐字复认＋机读尺 | 选形归编排者（`ticket248-ruling-248a2` 末段明写本轮不裁 ⇒ 本腿⛔ 不自选） |

### 1.3 Q3 越界检查怎么写（AC#5，且是本腿自己那把尺，不是相信别人）

- **名册尺**：每次 commit 前 `git diff --cached --name-only` 逐行读 ⇒ 出现的每一枚路径必须落在 §0 那张写面表里；出现 `frontend/**`／`design/**`／`docs/PLAN.md`／`docs/specs/**`／`docs/SLO.md`／`internal/observe/thresholds.go`／任何 golden／`tools/d22scan/allowlist.txt`／`.github/**`／`scripts/**`／`.scratch/wisp/issues/**`／`docs/reports/**` ⇒ 本腿自己判**退回**，不等人拦。
- **终态名册尺**：`git status --porcelain` 与 §0 存档的起手名册做**差集** ⇒ 新增的每一枚必须是本腿路径；⛔ 不比"是否为空"（本仓工作树永远不干净，那是错尺）。
- **代码形状尺**（`sh scripts/d22scan.sh`，四数之一）：裸 `go func(` 无 owner/recover、`risk.PathResolver` 之外用 `filepath.Clean|Abs` 做路径决策、明文 API 密钥、墙钟差实现超时、面板侧来源的 L2「允许」、宿主内部 artifacts 写成受门控 Tool、UI emoji（ban #8 射程**含 `_test.go`**、注释豁免字符串不豁免）。
- **冻结件尺**：本腿**只**动 `A487` 具名解冻的那两枚锚（`internal/panel/l2_grant_boundary_test.go:2051`／`:2139`）与其直连 helper；`git diff -- numstat` 上那枚文件的删除列必须只覆盖这两处附近；**其余一字未动**由 §4 那把 `git diff --unified=0` 逐 hunk 尺读数证明。
- **凭据面尺**：交件前后各跑一把 `grep` 全仓（含本件、日志、快照产物、测试数据）找哨兵串 ⇒ 除"本件 §4 里那句哨兵定义"之外**零命中**；真实凭据值本腿**从未读过**（`Store.Resolve` 的返回值本腿不进任何测试断言的比较串）。

---

## 2. 判据进度表（票面 AC#0…AC#3、AC#5…AC#11；AC#4 本腿不做也不勾）

| AC | 本腿要做到什么 | 承载点（写面） | 状态 |
|---|---|---|---|
| AC#0 | 三问落本节，各带现读凭据 | 本件 §1 | **已交**（先于产码） |
| AC#1 | 两枚名进白名单＋`ParseComposerRequest` 的守卫＋派发表；负向配正控 | `internal/panel/bridge.go`、`internal/panel/composer_dispatch.go`、`internal/panel/config_handlers.go`（新）＋用例 | 待做 |
| AC#2 | 快照那一维（状态、非值）＋哨兵零命中**常驻用例**＋正控 | `internal/panel/composer.go`（`ComposerState` 子段，J3 裁：⛔ 不加顶层第五键）、`internal/panel/pump.go`（reader 形状）、新用例 | 待做；⚠ 已知两界对账门会把它洗红（§5） |
| AC#3 | 只走 `secret.NewStore`；不落 `config.toml` 明文；不新增回显值方法 | `cmd/wisp/panel_config_store.go`（新）、`internal/config`（新 setter 只写 ref 形状） | 待做 |
| AC#4 | **本腿不做、不勾**（派单明令）；交回"欠账＋缺哪两样" | 本件 §7 | 欠账（读数见 §4/§7） |
| AC#5 | 越界检查＝§1.3 那三把尺 | 本件 §4 | 逐次 commit 自证 |
| AC#6 | 门禁四数 rc=0＋逐字读数 | 本件 §4 | 待填（第 100 轮前必满） |
| AC#7 | 写前必过同一套 `validate()`；种非法值⇒磁盘逐字节不变 | `internal/config/settings.go`（新 setter 内、进 `mergeWrite` 前） | 待做 |
| AC#8 | 回执带生效档＋"要重启"到页面可见面＋不许出现"保存即生效" | `internal/panel/config_receipt.go`（新）、`cmd/wisp` 装配 | 待做（Go 半枚） |
| AC#9 | embed 后页面产物是否真存在（能力形、不问 rc） | 本腿**不判**（`frontend/**` 两层禁令）；只交名册读数 | 欠账，归口 §7 |
| AC#10 | 常驻腿吃不吃 `[risk]` 那两项：只出机读读数＋归口，⛔ 不自选形 | 本件 §6 | 待读（复认前人断言） |
| AC#11 | 只走按键受保护写；三条子判据（保 B 键／逐枚键各一发＋报告落盘路径／AC#7 落点在 setter 内） | `internal/config/settings.go` ＋用例 | 待做 |
| J1（AC#1 的落地前置） | 把 `:2051` 拼写钉与 `:2139` 尺寸钉换成**能力形**＋自带正控；⛔ 不删钉、不放宽成 ≥4 | `internal/panel/l2_grant_boundary_test.go` 那两枚锚＋其直连 helper（`A487` 边界原文） | 待做 |

---

## 3. 产码改动清单（做完一节更新一行，带 `文件:行`）

| # | 文件 | 改了什么 | 现读锚（改后） |
|---|---|---|---|
| 1 | `internal/panel/bridge.go` | 名册加 `MethodConfigGet="config.get"`／`MethodConfigSet="config.set"`；封套加四枚**非值**选择子（`configField`/`configProvider`/`configModel`/`configValue`，全 `omitempty`）；守卫 case 从四枚变六枚（**同一枚 case 子句、逗号接排**，冻结件那枚能力形锚就是照这个形状读的） | `:41-70`（常量与来历，`MethodConfigGet` 在 `:66`）、`:95-121`（封套新字段，`configField` 在 `:117`）、`:146-152`（`knownComposerMethod`） |
| 2 | `internal/panel/composer_dispatch.go` | 新增 `ConfigRequestHandler`（签名带 raw，理由＝J4 甲：值不进共用封套）＋`ComposerDispatch.Config` 槽＋两枚 case；`dispatch` 返回值改 `(string, error)`（只有设置那两枚门填回执，其余四枚仍返空串，所以 `Handle` 对旧四门的句子里没变） | `:90-113`（接口与理由，声明在 `:110`）、`:134`（槽）、`:175-215`（dispatch 六枚 case＋default） |
| 3 | `internal/panel/config_handlers.go`（**新增**） | 可写字段枚举 7 枚＋锁定族显式拒写（`risk.`/`fs.`/`net.`/`plugins.`/`privacy.`/`models.`/`audio.`）＋选择子必配齐＋空值拒；回执形状带 tier；`ConfigStore` 接口**本地声明、不 import `internal/config`**（照 `ModeWriter` 先例，零新依赖边） | `:57-104`（字段表与 `WritableFields`）、`:108-140`（锁定族与逐族理由）、`:307-365`（write 闸门）、`:428-442`（`tierSentence` 文案律） |
| 4 | `internal/config/settings.go`（**新增**） | 六枚导出 setter（`SetProviderBaseURL`／`SetProviderAPIKeyRef`／`SetModelContextWindow`／`SetModelPriceIn`／`SetModelPriceOut`／`SetRoleChatModel`），全走 `mergeWrite` 那枚按键受保护写；**AC#7 的 `validate()` 闸门落在这里、进 `mergeWrite` 之前**；落盘后按"文件前后差"报实际键路径 | `:60-190`（setter 群）、`:199-249`（`writeOneKey`：预读→check→apply→`validate`→`mergeWrite`→回滚）、`:254-272`（`writtenKeyPaths`） |
| 5 | `cmd/wisp/panel_config_store.go`（**新增**） | 装配根那半：`configStore` 实现 `panel.ConfigStore`；凭据值只由本文件的**只写形状** `credentialWriteRequest` 绑定（JSON 键 `credentialValue` 不在共用封套上）；`secret.NewRef`→`Store.Store`（DPAPI）→`SetProviderAPIKeyRef`；tier 一律 `restart` | `:44-68`（`credentialValueKey` 与只写形状）、`:238-283`（`StoreCredential`）、`:285-292`（`credentialStatus`＝泵的 reader） |
| 6 | `cmd/wisp/panel_inbound.go` | `newComposerDispatchChain` 里接上 `Config` 槽（复用同一枚 `config.Manager` ＋ 现成 `secret.NewStore`）；受理支现在把 `Handle` 的回执打出来（AC#8 的可见面） | `:176-186`（回执支）、`:262-280`（装配，`Config: configWrites` 在 `:279`） |
| 7 | `internal/panel/composer.go` + `internal/panel/pump.go` | 快照那一维进 `composer` 段（`credentialState`＋`credentialKnown`，照 `currentModel`/`modelKnown` 的"一枚事实＋一枚 known 位"形）；`PumpSources.Credential` 是 reader，没接 reader 就报 `unknown`/`false` | `composer.go:258-269`（两枚键）、`:288`（构造器的 unread 支）、`pump.go:146`（reader）、`pump.go:260-268`（填充） |
| 8 | `cmd/wisp/run.go` | 装配根把凭据 reader 递给泵（`rt.settings` 用的就是 `wisp run` 已有的那枚 Manager 与那枚 `secret.Store`，⛔ 没有第二份加载、没有第二套存储） | `:259`（字段）、`:416`（装配）、`:695`（`Credential:` reader） |
| 9 | `internal/panel/l2_grant_boundary_test.go` | **只在 `A487` 具名解冻的那两枚锚上动**：`:2051` 拼写钉 → 从"这枚包声明的 `Method*` 名册"推出的能力形锚（含正控：种一枚不在名册里的路由标签 ⇒ 派生必报错）；`:2139` 尺寸钉 `!= 4` → `sortedSet(pkg.answered)` 与名册值集合全等（⛔ 未删钉、未放宽成 ≥4）；直连 helper＝同文件相邻位置新增 `guardRosterOf`/`dedupeSorted`/`slicesContains` | 改后现读：`:2030-2138`（`guardRosterOf`）、`:2140-2158`（两枚 helper）、`:2206-2226`（锚＋plant B 正控）、`:2308-2313`（名册全等支） |
| 10 | 测试（新增 3 枚文件＋1 处调用点） | `internal/panel/config_route_248_test.go`（AC#1/AC#8，**11 枚顶层用例＋2 枚子用例**，逐枚带正控）、`internal/config/settings_248_test.go`（AC#7/AC#11/AC#3，**7＋1 子**）、`cmd/wisp/panel_config_248_test.go`（AC#2/AC#3 端到端＋结构形正控，**9 枚**）；`composer_dispatch_test.go:235` 只是 `dispatch` 换签名的调用点（断言原文一字未改） | 尺＝`grep -c "^func Test"`／`grep -c "t.Run("`，`09:0x` 现量 |


**没动的东西（具名，供越界尺对照）**：`internal/secret/**` 一字未改（只调用现成 `NewStore/Store/Exists/Resolve/NewRef/ParseRef`，存储格式没动）；`internal/ball/**` 未碰；票面／台账／HANDOVER／规格／阈值／golden／allowlist／CI／scripts 未碰；`SaveFile` 与 `mergeWrite` 本体未碰。


---

## 4. 门禁与读数（AC#6 四数逐字；跑测试前必须先 `export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`）

（待填）

---

## 5. 我可能写错的条目（对抗我自己；⛔ 交件时不许留"待填"）

（待填）

---

## 6. 判不动的地方（逐条具名交回；⛔ 交件时不许留"待填"）

（待填）

---

## 7. 交件判语与欠账（AC#4 缺哪两样在此归口）

（待填）
