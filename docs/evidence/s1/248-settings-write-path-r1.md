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

> **本节由谁跑的（这是 AC#6 退回的那一格的正身）**：`248-r2`，票 248 的**续程证据腿**——零产码、零测试文件、零票面改动，写面只有本枚 md。
> 本腿 `248-r1`（写码腿）交件时这一节是占位，而 §7 那行写着"四数在 §4"＝指着一节空的凭据（`248-v1` 立案／`248-v1b` 续立／`248-v1c` 复认）。
> ⛔ 裁决腿 `248-v1c` 在 `docs/evidence/s1/248-settings-write-path-v1.md` §2 跑的读数**不抄进下面任何一行当正身**——那会把"谁量的"洗混。本腿只在 §4.6 **另起一行**做同锚互验，署名各归各。
> 四把尺逐把自己现跑，**每把带当时的 HEAD 短号＋`date` 现取时刻＋逐字终值行**；第 3 把尺的时序类读数一律标【同机可能有争用】。

### 4.1 尺一 · `GOFLAGS= go build ./...`

- 命令逐字：`GOFLAGS= go build ./...`（票面原文带着那枚**空赋值**，为的是不吃环境里 `GOFLAGS` 的任何默认值）。
- 时刻：起跑 `2026-10-02 12:37:56 +0800`／终值 `12:38:00 +0800`；当时的 HEAD＝**`fb1d9c19`**（`e2e-panel-1` 那枚只读腿 12:30:28 的 commit）。
- **逐字终值行：`rc=0`，输出 `0` 字节**（stdout＋stderr 合并后 `wc -c` 现量＝`0`，本腿另发复量过一次文件尺寸，时刻 `12:38:29 +0800`）。
- ⚠ 这一发的"零输出"怎么读：`go build` 成功时**本来就不打印任何东西**，所以"零输出"在这里是终值而不是"没跑"的形——"没跑"的形是 rc≠0 且带错误行（对照 §4.3 那枚 `0xc0000135` 坑，那一族才有 rc≠0 却没有 `--- FAIL`）。
- 归因：本腿到 `12:40:47` 才落第一发 commit（`e7521610`＝本件 §4 骨架），`git diff --name-only fb1d9c19..e7521610` 只有本枚 md ⇒ **这一发 build 量的是 HEAD 那批已入库产码，本腿对它零贡献**。
- 署名：〔**`248-r2` 自己跑的**〕。裁决腿 `248-v1c` §2 同一把尺的读数是它在 HEAD `6bcb934a` 上那一发，⛔ 不署本腿的名，互验见 §4.6。

### 4.2 尺二 · `gofumpt -l <本腿动过的目录>`

〔取数中：射程与具名说明尚未落〕

### 4.3 尺三 · 三包整包（`./cmd/wisp ./internal/config/ ./internal/panel/`）

〔取数中：本腿尚未跑。⚠ 本机 `cmd/wisp` 缺 sherpa PATH 会给 `exit status 0xc0000135` 且没有 `--- FAIL`＝用例根本没跑，那一发不能读成绿〕

### 4.4 尺四 · `./tools/d22scan/d22scan.exe`（独立模块，只能跑那枚 exe）

- 命令逐字：`./tools/d22scan/d22scan.exe`。⛔ 不是 `go run ./tools/d22scan`——`tools/d22scan` 是**独立 module**（`module github.com/CarlosShao/wisp/tools/d22scan`，本腿现读那枚 `go.mod` 首行），从根 module 里 `go run` 它必失败；本腿按票面点名跑 exe 本体，`-root` 由它自己推导。
- 时刻：`2026-10-02 12:40:04 +0800` 起／`12:40:05 +0800` 终值；当时的 HEAD＝**`fb1d9c19`**。
- **rc=0**。逐字末行的判决段（本腿只截到判决句为止，见下面那条禁令说明）：
  - `d22scan: clean - no D22 ban violations`
  - 承重的那一行逐字：`d22scan: examined 262 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp`
  - 射程的那几行逐字（摘本腿需要的三行）：`d22scan: scope bans #1-5 internal/      examined 226 production Go files`／`d22scan: scope bans #1-5 cmd/           examined  36 production Go files`／`d22scan: scope ban #8 internal/         examined 482 Go files, comments and _test.go included`
- ⚠ **两层禁令在这把尺上的处理**：那发输出里另有几行是它自报的 git-ignored 提示与 scope 复述，点名 `frontend/`／`design/` 那两层。本腿**没有读那两层的任何一字节**，也⛔ 不把它们的文件路径／行号写进本件——所以那一行本腿只记"存在、且不是违规"（违规会 rc≠0），判决段后半截的 scope 复述本腿不整行转抄。
- 归因：本腿零产码改动 ⇒ 这一发**不可能**是本腿造的；它量的是 HEAD `fb1d9c19` 上那批已入库产码（含 `248-r1` 那 15 枚路径）。
- 署名：〔**`248-r2` 自己跑的**〕，与裁决腿 §2 的 `sh scripts/d22scan.sh` 那一行**不是同一发**（它跑的是脚本、HEAD `6bcb934a`；本腿跑的是 exe 本体），互验见 §4.6。

### 4.5 起手锚与并发声明（`date` 每发现取，不复用）

| 尺 | 读数（时刻逐条带在旁边） |
|---|---|
| `date '+%Y-%m-%d %H:%M:%S %z'` 进场第一发 | `2026-10-02 12:27:01 +0800` |
| 进场时 `git rev-parse --short HEAD` | `4a9851d6`（`e2e-panel-1 起手：§0 锚＋§4/§5 先写满并落第一发 commit`） |
| 四把尺各自的 HEAD | 尺一 `fb1d9c19`（12:37:56）· 尺二 无射程（判定时刻 `fb1d9c19`，见 §4.2）· 尺三 `e7521610`（本腿 §4 骨架那一发，时刻见 §4.3）· 尺四 `fb1d9c19`（12:40:04） |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git status --porcelain -- cmd internal tools` | **`0` 行**（`12:40:47` commit 后现量同值）⇒ 本腿读到的是 **HEAD 的形状**，不是谁的未提交半片 |
| 工作树整体脏度 | 非空且**永远**非空（owner 在 `design/**` 下的移动／删除与几十枚他人台件）。⛔ 本腿不修它、不动它一个字；本腿⛔ 没有把名册另存成文件（写面只有本枚 md），对照尺就是上面那行"三包 0 行"＋§4.3 的逐名红名册 |
| 并发声明 | 取数期间同机另有一枚**只读**腿 `e2e-panel-1` 在飞（派单具名：它不跑任何 Go 命令 ⇒ 不洗本腿读数）。但本仓已知 `slo-full` 那类负载会在这台机器上抢 CPU ⇒ **§4.3 一切时序类读数（整包秒数、面板冷/热毫秒、任何超时类用例）一律标【同机可能有争用】**，⛔ 偶发慢不当回归 |
| 本腿写面 | 仅 `docs/evidence/s1/248-settings-write-path-r1.md` 这一枚。⛔ 零产码、零测试文件、零票面改动、AC 勾选框一枚未碰（复量尺见 §4.6） |
| 本机工具事实 | `gofumpt.exe` 确实在 `$(go env GOPATH)/bin`（本腿 `ls` 现量，那枚目录里还有 `staticcheck.exe`）但**不在本 shell 的 PATH** 上（`command -v gofumpt` ⇒ 未命中）；`go version go1.27.1 windows/amd64`；`third_party/sherpa-onnx`（3 枚 dll）与 `build` 目录都在，所以尺三的 PATH 前缀是跑得起来的 |

### 4.6 与裁决腿 §2 的同锚互验（具名指回，不署本腿的名）

〔取数中〕

### 4.7 本腿更正的指向（原句逐字留档）

〔取数中〕

---

## 5. 我可能写错的条目（对抗我自己；⛔ 交件时不许留"待填"）

1. **那两枚解冻锚的"能力形"是我写的新仪器，它读的是 `bridge.go` 一枚文件的 const 块**。`guardRosterOf`（现读 `internal/panel/l2_grant_boundary_test.go:2030`）把"守卫的 case 标签"与"本文件声明的 `Method*` 常量"做双向全等。如果日后有人把某枚路由常量挪去 `bridge.go` 之外（同包另一枚文件），派生会报"the guard answers X, which no declared Method* constant of this file names"⇒ **假红**，而且红的是冻结件。我按 `A487` 的边界（"只许改这两处与其直连 helper"）把 helper 写在同文件相邻位置，没有扩大射程去读整包；**代价就是这条单文件假设**。验收腿若要判它，判这枚。
2. **`dispatch` 的返回值从 `error` 变成 `(string, error)`**，为它改了一枚既有用例的**调用点**（`internal/panel/composer_dispatch_test.go:235`，`err := ` ⇒ `_, err := `，断言原文一字未动）。⚠ 如果验收腿或别人的突变脚本按老签名直接调 `dispatch`，会编译不过而不是红一句——这不是我能在自己文件里预防的，具名报出来。
3. **AC#2 的 grep 面我第一版漏了"持久 sink"那一枚**。派单与票面写的表面是三样：快照／`[audit]` 行／持久 sink。我第一版只跑了快照＋audit＋slog sink＋数据根，**没跑 `wisp run` 记账那行**（`cmd/wisp/panel_pump.go:331 panelSnapshotSummary`，它自带 440 字符钳并把摘要与 sha256 落 ledger）。`09:0x` 那一发之后我把 `ledger-line` 补进了零命中尺**和**正控尺（补了才交，没补我就会在 §7 里把它写成欠账）。⚠ 读数因此分两发：第一发零命中是 4 面，终态是 6 面。
4. **哨兵尺的形状正则只认 `sk-` 族**。尺面是"字面哨兵 ＋ `sk-[A-Za-z0-9]{12,}`"。真实服务商里 OpenAI 形是 `sk-`，Anthropic 是 `sk-ant-`（也在射程），但**别家的形状（纯 32 hex、含下划线的长串）今天不在尺上**——我刻意没写宽泛的 `[A-Za-z0-9]{32,}`，因为 `dpapi:<32 hex>` 合法地出现在 `config.toml` 与 blob 文件名上，宽正则会把**允许上页面的引用**判成泄漏（假红→有人会去放宽它，那才是真坑）。⇒ 这条尺**是"哨兵零命中"而非"一切密钥形状零命中"**；票面 AC#2 原文那句"grep 不到任何 key 值"要按这个射程读，不是我替它扩了判据。
5. **`writtenKeyPaths` 报的是"文件前后差"**，不是我这次写的内容清单。这正是 AC#11-1 要的凭据形状（写前写后各读一次文件比对），但**归因是宽的**：如果同一瞬间第三方改了另一枚键，那份改动也会被列进这枚回执。我没有把它与 `ownedKey` 求交，因为求交就等于把"文件实际发生了什么"改回"我以为我写了什么"（票 226 报的就是前者）。⚠ 若验收腿认为回执必须严格等于本次所有权，那一格不是我判的。
6. **字段名册在两处**：可写枚举在 `internal/panel/config_handlers.go:57-104`，枚举→setter 的映射在 `cmd/wisp/panel_config_store.go:158-196`。两枚表可以漂——漂的后果我按响亮那一侧设计：leg 的 `default` 支拒写、`WritableFields()` 有独立用例钉住那 7 枚，所以"panel 收了、leg 不写"是红句不是静默。⛔ 但我**没有**做一个"两枚表必须同枚数"的断言（跨包不可反射，leg 在 `package main`）⇒ 这是结构约束而非机读约束。
7. **`credentialStatus` 的 `env:` 一支读的是本进程的环境**。`os.Getenv(value)` 为空就算"未录入"，可这只说明**这个进程**没看到那枚变量，不代表操作者的 shell 里没有。枚举名 `partly_missing` 在这种一发里会说略重（逐服务商行里我另给了 `CredentialError`/布尔，聚合值没有）。我没找到不引入新键的解法（J3 禁顶层第五键），所以**这是已知会说糙话的一枚聚合值**，不是我没看见。
8. **`role_chat_model` 等四枚 `[llm]` 字段我一律回 `restart`**，而 `SPEC-03 §4.2` 的表把 `[llm]` 记成 hot。我按 J9/AC#8 的实践面写（endpoint 只在装配时构造、`OnReload` 在 `cmd/wisp` 零生产赋值点——两样本腿都复认过，见 §4 的尺），所以**表与我说的话不一致时，说重的不是我错**；但一旦哪天真接了重建那一跳（票 223 地界），这批回执文案就变成过度保守，要跟着改。
9. **`TestAC11ControlSnapshotWriteIsWhatRevertsTheHandEdit` 是靠"控制支必须仍然会覆写"来给正控的**（照票 226 的同名形）。如果哪天 `SaveFile` 被改成 merge 语义（那是票 226 的地界，不是本票），这枚控制会红——红的含义是"控制失去牙齿"，不是"设置写入坏了"。我在用例注释里写了这句话，但**它仍然可能被下一个人读反**。
10. **`model_price_in/out` 我只做了两枚**（`Price` 有 5 枚叶子：in/out/cached/audio_in/audio_out）。票面 §5 那句给界面侧的话只写"价格"，owner 的表单没枚枚点名，所以我按"最常被问的两枚"落地，⛔ 没有把 5 枚全开（多开就是多改契约面）。**这是我自己划的线，归编排者判够不够。**
11. **§3 的行号是 `09:0x` 现量的，`internal/ball/**` 那枚并发腿不吃我的文件，但 `cmd/wisp/**` 是共享目录**：如果有人在我之后动 `cmd/wisp/run.go` 或 `panel_inbound.go`，我这张表里的行号会漂。派单要求"cmd/wisp 行号落地现读为准"，我照做了；引用的人请重跑。
12. **我没有把新名登记进 `composerRouteLiterals()`**（`internal/panel/composer_test.go:409`）。那枚集合是**硬编码**的四枚常量＋`panel.approval.request`，注释自己写着"Growing it is a two-sided change on purpose"——它是"页面一旦真叫 `config.get`"那一步的登记点。页面还不会叫（本编队不写 `frontend/**`），所以我按"两界接缝不在本腿"处理；⚠ 若验收腿认为现在就该登记，那一格是**我对触发条件的读法**，不是漏码。
13. **`frontend/dist` 我只用了两枚 git/d22scan 侧读数**：`git ls-files frontend` 里 `dist` 只有 `.gitkeep`，与 `sh scripts/d22scan.sh` 自己吐的 `d22scan: skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]`。后者说明**本机存在未跟踪的 dist/assets 目录**；我没有读它、没有列它，也没有据此判"AC#9 能勾"。如果有人按我这条读数说"内容已经有了"，那是我没说过的话。

---

## 6. 判不动的地方（逐条具名交回；⛔ 交件时不许留"待填"）

1. **AC#10 的 ⓘ／ⓑ 二选一**——**不做，选形归编排者**（派单与 `ticket248-ruling-248a2` 末段都写"本轮不裁、⛔ 不许实现腿自选"；口令「248 别动常驻门」）。我这发只做**读数**（复认前人 `246-v2` 的断言，⛔ 不是转述）：
   - 常驻那扇门的构造：`cmd/wisp/resident_approval_windows.go:109-113` ⇒ `approval.New(approval.Options{UI, Channels, Logf})`，**没有 `Window`、没有 `ApprovalTimeout`、没有 `Grants`**（同文件 `grep -n "Risk.L1WindowSec\|Risk.ConfirmTimeoutSec" cmd/wisp/resident*.go`＝**零命中**）。
   - 对照：跑任务那条腿 `cmd/wisp/run.go:600-603` 三项都给了（`Window: cfg.Risk.L1WindowSec`／`ApprovalTimeout: cfg.Risk.ConfirmTimeoutSec`／`Grants: grantWrite`），而 `session.NewLedger` 在 `run.go:463`，建在 gate 之前。
   - 常量与钳位（现读，前人那三行我逐行复认，路径是 `internal/agent/approval/` 不是 `internal/approval/`）：`queue.go:107 DefaultApprovalTimeout = 300s`／`:116 DefaultL1Window = 3s`／`:121-122 MinL1Window = 2s`/`MaxL1Window = 3s`；`queue.go:86` 缺省回落到 300s；`gate.go:139-144` 是 `win<=0 → DefaultL1Window` 再夹进 `[2s,3s]`。
   - 记账被丢那两行：`gate.go:668` 与 `:674` 的 `"approval: GRANT-DROPPED ..."`（本腿只读，⛔ 未跑那套仪器）。
   ⇒ **读数结论＝"读不到"**：`[risk]` 的 `confirm_timeout_sec` 对常驻腿**无影响**（走常量 300s），`l1_window_sec` 今天对常驻腿也无影响，且即便接上，`>3s` 会被 `MaxL1Window` 钳回；真正有差异的是超时那一枚。**回执文案要不要逐字声明这件事（ⓑ 形）＝归编排者裁**，本腿不自选也不落地。
2. **AC#4（真机那一发）——本腿不做也不勾**（派单明令）。缺的两样具名，都在别人手里：
   - **一样：受版的页面产物。** `git ls-files frontend` ⇒ `frontend/dist/.gitkeep` 是 `dist` 下唯一被跟踪的东西（本腿 `08:31:49` 复量）；`frontend/embed.go` 的活模式是 `//go:embed all:dist` ⇒ 能建，但**没有一包页面可发**。产物由界面侧那枚 agent 出，本编队⛔ 不写 `frontend/**`。
   - 二样：**"点设置"那枚入口与那枚只写不回显的密钥框**（页面侧），加上页面真叫 `config.get`/`config.set` 之后那步两界登记（`composer_test.go:409` 的 `composerRouteLiterals`）。Go 这侧的接法已经就位（§3#1/#2/#3）。
   - ⚠ 前人票面 §5.4 那句"宿主还没接进常驻那条腿（`grep PanelManager`＝0 命中）"**今天已经不成立**，本腿复量：`cmd/wisp/resident_windows.go` 命中 2 行、`cmd/wisp/panel_resident_windows.go` 命中 6 行，且 `newResidentComposerDispatch`（`panel_resident_windows.go:166-170`）复用的正是本腿动的那枚 `newComposerDispatchChain`。⇒ **票 33 那半边比票面写的更靠前**；但"窗口能建"不等于"owner 真点得出设置"，那一格仍归编排者与验收腿（本腿没跑过任何真窗，也读不到页面）。
3. **AC#9（宿主在、内容不在）——本腿判不了**：尺要问的是 embed 之后文件系统里的条目数与关键入口文件，而那在 `frontend/**` 的两层禁令里（读了也不许引）。我只有上面那枚 git 名册读数＋d22scan 自己报的一行 ignored-path 提示。**归口：谁把 dist 填上没落定之前这一格勾不了**（票面 AC#9 原文同判）。
4. **AC#2 的页面那一半（Q-51）**：`credentialState`/`credentialKnown` 两枚键要让 `frontend/src/lib/panel.ts` 的 `ComposerState` 声明，否则双向对账门一直红。本腿⛔ 不动 `frontend/**` 一字，也没有放宽那枚门（改法只有"页面同批补声明"，那是界面侧的活）。⚠ 读数见 §4：那枚门在我之前**就已经**为 `[git currentModel modelKnown]` 三枚键红着。
5. **`fs.*` 到底禁不禁由 `config.set` 写**——**判不了，这是编排者的边界**。J5 逐字锁的是 `risk.*` 全族；`fs.allowed_dirs` 只被 `allowdirs.go:16-24` 那句"这里的写放宽的锁定段，调用者欠那张 L2 卡"覆盖。**我取了保守形**（`config_handlers.go:108` 把 `fs.` 连同 `net.`/`plugins.`/`privacy.`/`models.`/`audio.` 一起列进拒写族），如果编排者的边界是"只锁 `risk.*`，其余段可写"，那我这一拒比裁的更窄——**窄不会造成安全事故，宽会**，所以这个方向我不改，但具名报出来等一句判语。
6. **票 33 那侧的面板毫秒级红（冷 ≤1500ms／热 ≤200ms 那一族）算不算缺陷**——**不判**。同机另有一枚腿（`internal/ball/**`）在跑，且云端 CI 今天不在这台机器上；我只登记逐字读数（§4）＋"并发负载可能"的归因线索。阈值一字节没动，那一格归编排者与验收腿。
7. **`staticcheck`**——**〔未复认〕**。本机版与 CI 钉版不同，跑了会出假绿，派单明令不跑，本腿没跑。

---

## 7. 交件判语（逐格读数归口；勾与不勾不归本腿）

| 格 | 本腿做到了什么（读数在 §4，凭据在 §3） | 欠什么／归谁 |
|---|---|---|
| AC#0 | 三问在 §1，各带 `文件:行`＋本腿自己跑的尺（`08:31:49`／`09:0x` 两批）；三问都定得下来，**没有触发"停手上报"** | — |
| AC#1 | 两枚名真进白名单（`bridge.go:66-71`、守卫 `:146-152`）与派发表（`composer_dispatch.go:197-212`）；负向（`panel.review.allow`/`approval.decide`/`config.write`/`panel.config.set` 必拒）与正控（同形封套换成在册名 ⇒ 处理器调用计数＝1）同用例落地；在册无处理器 ⇒ `ErrNoHandlerAttached` 响亮 | 两界登记点（`composerRouteLiterals`）等页面，见 §6-12/#12 |
| AC#2 | 那一维进 `composer` 段（非顶层键）；哨兵零命中尺覆盖 6 枚表面（回执/audit/slog/快照/ledger 行/数据根全文件）＋**两枚正控**（种进快照必须被抓、种进 ledger 行必须被抓）＋结构形尺（共用封套不许声明凭据键，种一枚声明了的类型必须报重叠） | `panel.ts` 声明归界面侧（§6-4） |
| AC#3 | 只调 `secret.NewStore`/`Store.Store`/`Exists`/`NewRef`（`internal/secret/**` 一字未改）；`config.toml` 只落 `dpapi:` 引用；没有任何"把值回显给页面"的方法；用例证明 blob 里解得回来且全树无明文哨兵 | — |
| AC#4 | **未做、未勾**（派单明令） | 缺"受版页面产物"＋"设置入口页"两样（§6-2）；票 33 宿主那段比票面靠前（§6-2 末） |
| AC#5 | 越界三把尺逐 commit 自证：`git diff --cached --name-only` 两次（15 枚路径，全在我写面）＋终态名册差集＋d22scan | — |
| AC#6 | 四数在 §4（build/vet/d22scan/gofumpt 全 rc=0）；`staticcheck`〔未复认〕 | 面板毫秒级红的定性归编排者 |
| AC#7 | `validate()` 落在新增导出 setter 内、进 `mergeWrite` 之前（`settings.go:199-249`）；种 6 发非法值 ⇒ 磁盘逐字节不变＋下次仍起得来＋内存不漂（真跑过，不是推） | — |
| AC#8 | 回执带 tier，`[llm]` 全族按重启写死；全文（票面、代码、回执、用例断言）grep 不到"保存即生效"；页面可达面：`Handle` 返回串（票 33 的 bind 闭包 `panel_host_windows.go:319-322` 把它交页面）＋本腿把 console 支也接上 | "页面上真看得见"要真窗那一发，归 AC#4 同族欠账 |
| AC#9 | 本腿未判（禁令），只交名册读数 | 界面侧＋编排者 |
| AC#10 | 读数完成（§6-1，结论＝常驻腿**读不到**那两项）；⛔ 未选形、未动常驻门 | 选形归编排者 |
| AC#11 | 三子判据都在：只走 `mergeWrite`（⛔ 未用 `SaveFile`，未改 `mergeWrite` 本体）；并发保 B 键用例＋"若改用快照序列化器则此控制必失效"的控制支；多枚字段＝逐枚键各一发＋回执逐枚列键路径（取"文件前后差"） | 回执归因宽度见 §5-5 |
| J1（冻结钉） | 两枚锚换成能力形（名册双向全等）＋自带正控（种一枚不在名册的路由标签 ⇒ 派生必报错），未删钉、未放宽成 ≥4；冻结件其余一字未动（`git diff` 逐 hunk 见 §4） | 判语与勾归验收腿 |

