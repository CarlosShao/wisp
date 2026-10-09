# 248 — 面板里没有「设置」这条路：入向白名单只有四枚方法、**零枚 config/凭据方法**；就算宿主起来了也点不出能自己录 key 的界面

Status: OPEN（编排者 09-30 23:1x 立，来路＝owner 当场提的功能要求，原话见下面第一行）
**owner 原话（逐字）**：「问题我早就说过了，要录入key，起码我要能看到主面板，**我要点击设置，自己配置模型这些参数**，直接给你就太不合适了」
实现者：`248-r1`（写码腿，**按住**，见「排程与串行」） · 前段普查：`248-a1`（只读，在飞） · 裁决者：`248-v1`（必须 ≠ 实现者） · 归口：面板（与票 33 宿主、票 35 桥、票 114 原生侧门同一条链）

## 0. 这一票为什么存在（不是我推断的需求，是他提的功能）

今天这台机器上"录一次模型 key"只有一条路：**命令行隐藏输入**（`cmd/wisp/secret.go`，里面有 `termHiddenReader`／`readHidden` 那套）。他明确说了这条路不合适——他要**看得见主面板、点设置、自己填模型参数**。凭据本身我永远不碰、也不许出现在对话／日志／快照里（只引用变量名）。

**这一条不是"宿主起来就顺带有了"**：入向通道今天的白名单**只有四枚方法**，里面没有任何一枚是配置或凭据。所以这是**两层缺口**，要分两票做：票 33 让窗口真存在，本票让"设置"那件事真有一条从页面走到磁盘的路。

## 1. 现量（每条带尺与取数时刻；⚠ 引用前重跑）

| # | 事实 | 尺与出处 |
|---|---|---|
| 1 | 入向白名单＝**四枚**：`panel.mode.request`／`panel.workspace.request`／`panel.attachment.add`／`panel.message.send`（`internal/panel/bridge.go:42-45`）；全仓扫出来的 `panel.x.y` 方法名去重后**也是这四枚**＝**零枚 config/secret 方法** | 编排者 09-30 23:1x 现跑 `grep -rn "panel\.[a-z]*\.[a-z]*" --include=*.go internal/panel/ \| grep -v _test \| grep -o "panel\.[a-z]*\.[a-z]*" \| sort -u` |
| 2 | 派发器把另外三枚槽留成 nil：`cmd/wisp/panel_inbound.go:238` 一带（`Workspace`／`Attachment`／`Message` 各归票 186／92／35），**只有 `Mode` 真接了** | `panel_inbound.go` 现读＋票 33 AC#9 |
| 3 | 今天录凭据的唯一生产路径＝CLI：`cmd/wisp/secret.go`（隐藏输入、不落中间文件，由 `TestSecretFromStdinWritesNoIntermediateFile` 钉着）；`secret.NewStore` → DPAPI（`internal/secret/dpapi_windows.go:24`） | 现读；与台账 `A476` 更正一致（DPAPI 已把密钥绑死他那个账号） |
| 4 | **C17（PanelBridge 方法白名单）是冻结契约面**：`AGENTS.md §0.2`／`SPEC-12 §4.1` 定"改契约＝人工批准"。⇒ 本票要新增方法名，**批准记录由编排者落一枚 `A##` 并逐字引用上面那句原话**；⛔ 但"改了契约"只决定要不要先登记，**永不决定这功能做不做**（owner 09-28 已就这一点当场改判过一次） | `AGENTS.md` §0/§2；`SPEC-12 §4.1`；台账 `A368` 同形先例 |
| 5 | 面板出向快照今天只有四键那一段（`panel_pump.go`），**没有任何"是否已录入凭据"的维度**；快照里绝不许出现 key 值 | 现读；票 35／145／146 那一串表的既有判语 |

## 2. 判据（⛔ 框归编排者，产码腿与验收腿一枚都不许碰）

- [x] **AC#0 先把三问摆开（不许直接开写）**：① **新增哪几枚入向方法名**、各自的参数字段与错误形状（候选形状：一枚读、一枚写；⛔ **不许把凭据塞进 `panel.message.send` 那条对话通道**——消息文本会进上下文预算与日志，那是一枚 key 最坏的去向）；② **写路径落在谁手里**（`secret.NewStore` 是唯一持有者？面板侧只许"发起请求"，照票 114 对权限档定过的那条规则推广过来，并**具名写成这是我的推广**）；③ **读侧怎么回答"配了没配"**而不泄露值。完成判据＝三问各带现读凭据（文件:行＋尺读数），写点只准 `.scratch/wisp/probes/248/a1/census.md`；⛔ 普查腿不许改任何产码。
- [x] **AC#1 白名单真扩**：新增的方法名进 `internal/panel/bridge.go` 的白名单与 `ParseComposerRequest`，且**负向判据配正控**（种一枚不在白名单的名字必被拒、种一枚在白名单但没人处理的那枚要响亮返回 `ErrNoHandlerAttached`）。⛔ 判据一律问能力，不许做成扫注释词面。
- [ ] **AC#2 快照那一维**：出向快照里出现"凭据是否已录入／配置是否可读"这一维时，**全仓任何日志与快照产物里 grep 不到任何 key 值**。尺＝把一枚可识别的哨兵值喂进写路径，然后 grep 快照／日志／持久 sink ⇒ **必须零命中**，且这一发要作为**常驻用例**进仓（不是表里的一次性读数）。
- [x] **AC#3 只走一枚凭据存储**：写路径复用现成 `secret.NewStore`／DPAPI，⛔ 不许新造第二套存储、不许把 key 落进 `config.toml` 明文、不许新增任何"把值回显给页面"的方法（`api_key_ref` 才是可以上页面的东西）。
- [ ] **AC#4 真机那一发（依赖票 33）**：一条完整链＝**面板点设置 → 填入参数与凭据 → 保存 → 下一次任务真用上刚配的模型**。⛔ **在票 33 的宿主真起来之前这一格不许勾**，也不许用"我手工改了 config.toml"来代替那一次点击。
- [x] **AC#5 越界检查**：`git diff` 出现 `frontend/**`／`design/**`（两层禁令：不许读、不许写）／`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／三枚冻结件 任一路径 ⇒ 直接退回。页面长什么样、按钮摆哪、文案怎么写**都不归本编队**（owner 自己带给他用的那枚 agent；本仓 09-28 为这事发过第三次火）。
- [x] **AC#6 门禁四数**：`GOFLAGS= go build ./...`、`gofumpt -l <动过的目录>`、`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/panel/ -count=1`、`./tools/d22scan/d22scan.exe`（独立模块，只能跑那枚 exe），逐名照抄终态。**〔14:1x 编排者现跑自勾：四把尺我自己复跑过，读数与本票末节那条对齐；凭据署名两枚——`248-r2` 的 §4.1–§4.5 与编排者这一发〕**

## 2b. 编排者裁定（10-01 00:0x，台账 `A487`）：`248-a1` 的九问逐条判完

普查件 `.scratch/wisp/probes/248/a1/census.md`＝**432 行／59,623 字节**，占位符现量 **0**，42 把尺（含它自己具名申报的"R34 跑空"与"R42 我读错后重测"），票面 AC 框**一枚没碰**（我现量：未勾 7、已勾 0）。它 ⑥#1 自己推翻了一条**不属于本仓的引文**（一处 `want 4: the whitelist changed size` 全仓 grep 不到），这一笔按"自己抓到自己的假引文"记功，处置＝逐条以真实行号为准（我复认过 `:2051`／`:2139` 两行原文，见下表 J1）。

| 问 | 裁 | 边界与依据 |
|---|---|---|
| **J1 扩白名单必撞冻结件** | **甲：人工批准，只动那两处锚点，且只许变强** | 现读两枚锚逐字：`internal/panel/l2_grant_boundary_test.go:2051` 是 `guardAnchor := "case MethodModeRequest, MethodWorkspaceRequest, MethodAttachmentAdd, MethodMessageSend:"`（**拼写钉**），`:2139` 是 `if got := len(sortedSet(pkg.answered)); got != 4`（**尺寸钉**）。⇒ 解冻授权落 **`A487`**（文件／行／理由／边界／撤销口令「**248 别动冻结件**」）；⛔ 只许改这两处与其直连 helper，改法＝把"数一枚 4"换成**能力形**（路由标签必须出现在 guard 名册里），并**自带正控**：种一枚没进名册的路由标签 ⇒ 指名用例必红。**不许把钉删掉或放宽成"≥4"。** |
| **J2 `git_test.go` 那两枚写死"四枚"的钉** | **甲＋乙合并** | `internal/panel/git_test.go:385`／`:517`（＋`:514`）同 commit 改写，含**报错文案**一起改；形状取乙（能力形：问"这张名册等于白名单吗"）。⚠ **这是票 181 的地界**，我按"把票 181 的判据形状推广到本票"处理，并**在票 181 面追加一行归口**，⛔ 不许悄悄改别人的判据。 |
| **J3 那一维放快照哪里** | **甲：进 `composer` 段** | 现成先例 `internal/panel/composer.go:255-256`（一枚事实＋一枚 known 位）。⛔ 不加顶层第五键（那会同时牵 `Q-51` 与票 145 AC#3 的旧裁，两枚都还没定）。⚠ 它给的现读更正我收下并写进 AC#2 口径：**`Snapshot` 不是"四枚键"，是 4 恒发＋2 可选（`composer.go:58-91`）**。 |
| **J4 凭据值的形状** | **甲：值绝不进那枚共用封套** | 现读：`internal/panel/bridge.go:69-79` 是**所有已答方法共用**的回执封套，而它的裁决词表不含 value/secret ⇒ **"靠词表拦值"拦不住**。⇒ 凭据写入走**单独一枚只写请求**（回执只回状态），配正控：哨兵串喂进写路径 ⇒ 快照／`[audit]` 行／持久 sink grep **必零命中**，并升成**常驻用例**（AC#2）。 |
| **J5 `permission_mode` 能不能由配置写方法改** | **甲：显式拒写 `risk.*` 全族** | 改档只能走既有那条带 L2 卡的路（`cmd/wisp/run.go` 的 `confirmModeSwitch` ＋ `panel.ModeWriteHandler`）。⛔ 不做乙（把"方向不对称"推广到锁定段——那是它的推广，我不采纳为规矩）；⛔ 不做"直接写"（正面撞 D33 与 `AGENTS.md` §1.2"面板侧来源的 L2 允许"禁令）。 |
| **J6 通用 `config.set(key,value)` 还是具名枚** | **具名枚（白名单式），名字直接用规格里已有的两枚** | 现读：`docs/specs/SPEC-08-ui-ball-panel.md:168` **已经写着 `config.get`／`config.set`**（文档先行、代码未跟上，那张表自署"提案，S5 定稿走契约批准"）⇒ 我不新造名，`A484` 的射程**收窄并重述**：只有这两枚名，且 `config.set` 的参数是**具名字段枚举**，⛔ 不许做成任意 key 透传（那等于把 `config.toml` 整张开给页面，`permission_mode`／`allowed_dirs` 全在射程里）。 |
| **J7 写完的回执形状** | **甲：回执带"何时生效"＋页面必须念出来** | 带 tier（立即／下次任务／需重启）；且**"热加载对这条腿无效"那句现在只在 stderr 念**（`hotReloadDisabledPanelInbound`）⇒ 归进 AC#1 完成判据：这条损失要出现在**页面可见的回执**里，不是只落日志。 |
| **J8 文档谁权威、代码扩名后要不要动文档** | **不动文档一字** | `SPEC-08 §5.2` 是方法名出处（已写了＝不需要我补）；`SPEC-03 §3` 缺 `permission_mode` 那处**已漂**，登记不修（`docs/specs/**` 是禁区，定稿走 owner 那边的契约批准流程）；⛔ 任何腿不许写"我补了规格"。 |
| **J9 AC#4 今天能不能满足** | **乙＋新补一格 AC#8** | 现读结论我接受：**同一次运行里"保存后下一次任务真用上刚配的模型"今天不通**（endpoint 只在装配时构造 `run.go:410-421`；`OnReload` 在 `cmd/wisp` **零生产赋值点**；`[llm]` 名义 hot、实则本进程不重建）。⇒ AC#4 改按**"重启后真用上"**判；⛔ 本票不许顺手加重建那一跳（那是票 223 的热加载地界）。 |

**另有一条它没问、我 10-01 00:0x 自己量到并要写死的**：`frontend/dist` 在库里**只跟踪着一枚 `.gitkeep`**（尺：`git ls-files frontend` ⇒ `frontend/dist/.gitkeep`），而 `frontend/embed.go:19` 的活模式是 `//go:embed all:dist`。`all:` 前缀连点文件也算进去 ⇒ **今天"能构建"很可能是因为目录里躺着一枚占位文件**：构建过得了，**却没有任何页面产物可发**。这条同时把票 33 的 AC#11 逼成两问（能建／有内容），并归口在下面 AC#9。

**还有一条它翻出来、我认为是本票最硬的实现约束**：**写路径今天没有写前校验闸门**——`SaveFile`（`internal/config/loader.go:238-249`）**不跑 `validate()`**（那只在上游 `loader.go:123`／`migrate.go:77` 跑）。⇒ 设置页一发非法值就能把**下一次启动**变成 Unconfigured。这条我写成下面 AC#7，⛔ 不许"先写了再说"。

## 2c. 新补判据（10-01 00:0x 编排者追加；勾仍归非实现者）

- [x] **AC#7（写前必校验，来路＝`248-a1` ①.5）**：任何经面板发起的配置写，**落盘前必须过一遍启动时那同一套校验**（`validate()` 的现成实现），校验不过 ⇒ 拒写＋把原因回给页面＋留审计。判据＝**种一发非法值**（例如把上下文档位写成非数）⇒ 磁盘上的 `config.toml` **必须逐字节不变**、下一次启动仍起得来。⚠ 它申报这条是**从函数体读出的结构判断、没实跑过**（本腿禁跑测试），所以读数必须由落地腿自己种一次才算数。
- [ ] **AC#8（归口"重建那一跳"）**：`[llm]`／endpoint 类字段在同一次运行内不重建（J9 现量），故 AC#4 只按"重启后生效"判；本格判据＝票面与回执文案**不许出现"保存即生效"**，且"要重启"必须到达**页面可见面**。⛔ 本票不许新增重建跳（票 223 地界）。
- [ ] **AC#9（宿主在、内容不在这一族）**：设置页要真能点出来，前提是 `//go:embed all:dist` 真匹配到**一包页面产物**。判据形状＝能力型、不问文案、**不问构建退出码**：embed 之后的文件系统里**条目数与关键入口文件必须真存在**（⛔ 一枚 `.gitkeep` 就能让 `go build` 绿，那不算证据）。⚠ 产物由界面侧那枚 agent 产出，本编队不写 `frontend/**` ⇒ 这一格在"谁把 dist 填上"落定前**勾不了**，具名交 owner 带话（要带的话在 §5）。

- [ ] **AC#10（新补，10-01 09:5x 编排者追加，来路＝`246-v2` §3.B＋§5 第 6 条 ⇒ 账 `A488`）**：**常驻那条腿今天吃不到 `[risk]` 那两项配置，本票要么接上、要么在票面写明它是常量——⛔ 不许留成"配置页改了它就变了"的错觉**。现量（`246-v2` 自取，我复认）：`approval/gate.go:139-145` 的钳位 `win <= 0 → DefaultL1Window`；`approval/queue.go:107/115/122`＝`DefaultApprovalTimeout = 300s`／`DefaultL1Window = 3s`／`MaxL1Window = 3s`。常驻那枚 `approval.New`（`resident_approval_windows.go:109`）建在会话账本之前 ⇒ 两处真实代价：① `Options.Grants` 为 nil ⇒「本会话内允许」**放行但不落盘**并写 `GRANT-DROPPED`（这条与票 224 那套会话授权直接冲突——用户在球上按了"本次会话内允许"，今天**不生效**，只有日志里有）；② `l1_window_sec`／`confirm_timeout_sec` 改动对常驻腿**无影响**（窗口那项即便接了也钳在 3s ⇒ **只有超时是真差异**）。
  **本格只许两选一，不许第三形**：ⓘ 把 `approval.New` 移进 `assembleRuntime` ⇒ **这是改票 246 AC#1 裁过的乙形次序**，动手前必须先由编排者落一枚具名 `A##` 批准记录（文件／行／理由／边界／撤销口令），且要兜住"移动会不会破 `Confirming` 那一维"；ⓑ 不移动，那么在设置页那几项的**回执文案**里逐字写明"这两项只作用于跑任务的进程，常驻腿今天用常量 300s／3s"——⛔ 不许只写"已保存"就把差异藏起来（与 AC#8 那条"不许出现'保存即生效'"同族）。
  ⚠ 判据形状＝**反向＋正控**：种一发 `confirm_timeout_sec` 改动 ⇒ 若走 ⓑ，页面必须显式声明不影响常驻腿；若走 ⓘ，`GRANT-DROPPED` 那行**必须不再出现**且"本次会话内允许"要真落一行（凭据复用票 224 r2 那套仪器，⛔ 不许新造）。撤销口令「248 别动常驻门」。

- [x] **AC#11（新补，10-01 10:2x 编排者追加，来路＝派腿前我自己跑的撞钉预检；这一格决定 AC#7 落在哪枚函数上）**：**配置写只能走"按键的受保护写"那一族，⛔ 不许用整份快照序列化器**。现量凭据（`internal/config`，全路径＋行号）：① 受保护原语＝`(*Manager).mergeWrite(ownedKey string, setOn func(base *Config))`（`internal/config/writeguard.go:111`，票 226 立的；它**先重读文件**、只替换自己那枚键，三形都不静默：文件读不回来⇒拒写并报错、文件已是新值⇒一字不写、真写了⇒报告"哪些键路径变了"＋"哪些外来手改被保留"）；② **既有先例三枚**（都长这个形状，照着写就不算自创）：`internal/config/allowdirs.go:62 AddAllowedDir`／`internal/config/allowdirs.go:109 SetAllowedDirs`／`internal/config/permmode.go:70 SetPermissionMode`；③ ⛔ **`SaveFile`（`internal/config/loader.go:238-249`）是 serializer 不是 merge**——`writeguard.go:5-17` 逐字写着"每一枚程序化写入都用本进程启动时加载的值重写 `config.toml`……操作者手改的任意一枚键会被一发聊天答复悄悄复原，而那枚写还会把文件新的 mtime＋size 认领成'我自己的'，于是连重载路径都看不见它刚毁了什么"＝**票 226 存在的理由**。用它就等于把票 226 修好的洞重新打开。
  **本格三条判据（缺一格不算完成）**：
  1. **正向**：种一发并发形状——面板侧写 A 键的同时，另有一枚**进程外手改**的 B 键已经在文件里 ⇒ 落盘后 **B 键必须逐字还在**（凭据＝写前写后各读一次文件比对，⛔ 不许只看 `rc=0`）；负向＝那一枚 setter 写失败时 A 键不许出现"半落"。
  2. **多枚字段一次保存＝逐枚键各一发受保护写（裁甲）**，⛔ 本票不许自创"多键原子写"新原语（那是票 226 地界，要另批另裁）；中途失败即停，**回执必须逐枚列出实际落盘的键路径**——`mergeWrite` case 3 已经给得出这份报告（`diffKeyPaths`，`writeguard.go:201`），⛔ 不许压平成一句"已保存"（与 AC#8 那句"不许出现'保存即生效'"同族）。
  3. **AC#7 那枚"写前必校验"的落点＝新增那枚导出 setter 内、进 `mergeWrite` 之前**。现量理由（别按"应该已经校验了"来写）：`grep -n "validate" internal/config/writeguard.go`＝**零命中**⇒ 校验今天不在这条写路径上，AC#7 是真缺口；⛔ 不许改 `mergeWrite` 本体（票 226 地界，且它会波及既有三枚 setter）。
  ⚠ **一枚本来会挡住本票的钉，射程我量窄了——具名撤除顾虑**：票 224 的 N#9 记着"那枚钉是词面型（禁 `Set*`），日后任何与授权无关的合法 setter 都会被它打红"。我 10:2x 现读真身 `internal/tools/ticket224_setter_scope_test.go:80-97`：它扫的是 **`reflect.PointerTo(bTyp)` 那一枚 `(*tools.Bridge)` 的方法集**（`:93` 逐字 `if strings.HasPrefix(m.Name, "Set")` ⇒ 红句是「(*Bridge).%s is an exported setter on the enforcement layer」），**不覆盖 `internal/config.Manager`**。⇒ 本票新增 `Manager.Set*` 那形**今天不会撞它**，N#9 保持待裁、不再算本票的挡路项。撤销口令「248 别动受保护写」。

## 3. 禁区（实现腿与验收腿共用）

- ⛔ `frontend/**` 与 `design/**` **两层禁令**：不许读、结论也不许引到它们身上。
- ⛔ **任何凭据值都不许进入对话／日志／表／快照**；需要提及时只写变量名或字段名。
- ⛔ 不改 `PLAN.md`／`docs/specs/**`／`docs/SLO.md`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／三枚冻结件；不许为变绿放宽任何断言。
- ⛔ 裸 `go func(`（ban #1）；`risk.PathResolver` 之外用 `filepath.Clean|Abs` 做路径决策；墙钟时间差实现超时；面板侧来源的 L2「允许」。
- git：只 commit 不 push；commit 必带显式 pathspec；禁 `add -A`／`amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；仓内不删文件。

## 4. 排程与串行

- ⛔ **`248-r1` 按住**：本票要动 `internal/panel/bridge.go`＋`cmd/wisp/panel_inbound.go`，而 `cmd/wisp` 此刻由写腿 **`246-r2`** 占着（同包两枚写腿＝归因搅浑，本仓 09-30 已付过两次代价）。顺序＝**246-r2 交件 → 票 33 宿主段（要动 `cmd/wisp`＋新依赖）→ 248-r1**。
- ⚠ **票 33 的 `Blocked by: 07-ball-state-machine-core` 已经消了**（球与托盘进常驻进程由票 228 落地、非实现者验收 `A477` 裁成立）。⇒ 按 owner 这句功能要求，**票 33 的"宿主真起来"那段提到队列头**，排在 `246-r2` 之后；那段要新增 WebView2 依赖＝动构建链，与票 244（GUI 子系统构建）同一批改会互相看不清，**分开发**。
- ⚠ 只读普查腿（`248-a1`）⛔ 禁跑 `go build`／`go vet`／`go test`／任何带 `./...` 的命令，理由同上（会吃到写到一半的码）。

## 5. 给界面侧的话（由 owner 自己带给他用的那枚 agent；本编队不碰 `frontend/**` 一字）

要让他那句「我要点击设置，自己配置模型这些参数」真成立，界面侧需要有两块，**Go 侧这头已经准备好接**（本票落地后即接上）：

1. **一个「设置」入口**（左竖条那排图标里的一枚，按 `Q1` 已批的口径顶一枚不贴切的也行，留 INTERIM 标记），点开是一页能填的表：**服务商地址／模型名／上下文窗口／价格／权限档位**，外加一枚**「录入密钥」的输入框**。
2. **那枚密钥输入框必须是"只写不回显"**：填进去之后页面**任何时候都不得再显示它的值**，只能显示"已录入／未录入"。理由写在票 248 `AC#2`/`AC#3`：值一旦回到页面，就会进日志、进快照、进模型上下文。
3. ⚠ 界面侧**不要**把密钥塞进聊天消息里发出来（那是最坏的去向，`AC#0` 已把它禁掉）；保存动作要走本票新增的那枚"配置写"方法。

本编队这侧对应要交付的：入向新增"配置读／写"两枚方法（批准记录＝台账 `A484`，撤销口令「248 撤 C17」）、出向快照补"凭据是否已录入"这一维、写路径复用现成 DPAPI 存储（`internal/secret/dpapi_windows.go:24`）。窗口本身能否存在＝票 33（本机 WebView2 Runtime 已装，缺的是我们这侧的宿主代码）。

⚠ **10-01 10:3x 编排者更新（宿主代码已落地之后，两句话给界面侧）**：

4. **今天仍然开不出窗口，但那一格不归界面侧，别在那边等**。宿主代码已经进仓（`cmd/wisp/panel_host_windows.go`，依赖是经 `go get` 正规解析的、没有手抄哈希），**但还没接进常驻那条腿**——现量 `grep -rn "PanelManager" cmd/wisp/main.go cmd/wisp/resident_windows.go cmd/wisp/run.go`＝**0 命中**。原因是 `go-webview2` 开窗用的是**阻塞嵌套的消息泵**，从球的 `ui-sta` 那条泵再入一次会 panic；实现腿按裁定**没有**擅自接那一跳、也**没有**为绕开它偷偷加第二枚线程（台账 `A492`）。这一格现在是 Go 侧的待裁项，验收腿 `33-v1` 正在判它报的那两枚依赖边界是真是伪。
5. **一条对页面形状的硬约束（现在就能确定，请让那枚 agent 按这个打）**：今天离线供给走的是把**一整包字节**塞进 `SetHtml`，而"按 URL 去取一堆分散资源文件"那条路（资源过滤器 `AddWebResourceRequestedFilter`）在主模块里**接不上**——依赖里那枚控制器的 `PutBounds` 吃一枚模块私有类型，外部 import 不到。⇒ **请把设置页做成单文件**（HTML＋内联 CSS/JS，或一枚能内联的 bundle），别假设运行时会自动去加载分散的 `.js`／`.css`／字体文件。⛔ 这条的性质要标清：〔读码推的＋实现腿自述，验收腿复认中〕——如果它被判不成立，我会再来追加一句更正，⛔ 不许界面侧自己按"应该能多文件"来做。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）

- [2026-10-01 11:57:07 +08] agent=248-a2（只读普查腿，零产码/零测试/零构建） did=交件 `.scratch/wisp/probes/248/a2/census.md`：逐格分类 AC#1..AC#11（Go 侧可独立=AC#1/3/5/7/11；两栖=AC#2/8；要界面或票33宿主=AC#4/9；AC#6 交落地腿自跑；AC#10 选形交编排者）＋路由/写侧现量带 file:line＋撞钉预检点名 T1-T8。未答(标〔推〕)=whitelistMethodsFromSource/composerRouteLiterals 函数体未读、cmd/wisp T4/T7 断言原文未逐枚读、冻结件内部语义按禁令未引。next=编排者裁 §⑦-3（config.* 前缀是否惊动票181 git_test.go）与 §⑦-2（grant-door 探测器是否误判 config.set）
- [2026-10-01 12:04:14 +08] agent=248-a2（续发）did=逐枚读完非冻结行为钉的断言原文并定稿 §④/⑤/⑥/⑦：坐实 T1 whitelistMethodsFromSource(git_test.go:407) 只抽 panel. 前缀⇒J6 的 config.* 命名不红；T2 composer_test.go:502 与 T4 panel_inbound_33 与 T5 leg_dispatch_gate_133 与 T7 writeguard_226 对 Go 侧独立写均不红；AC#2 被 composer_test.go:48 双向对账门洗红⇒Go 半枚做不到 self-green；pump_test.go:123 钉死顶层四键⇒新维度进 composer 段。净结论:会叫红的只有两格=T3(l2_grant 冻结件 :2051/:2139,无论命名必红,走 J1 解冻改法) 与 AC#2 对账门。仍判不了并交裁: grant-door 探测器会不会误认 config.set(§⑦-2)、AC#10 选形(§⑦-1)、fs.* 能否由 config.set 写(§⑦-4)。next=等编排者裁 §⑦-2



---

## 编排者裁（10-01 12:0x，收 `248-a2` 撞钉预检，普查件 `.scratch/wisp/probes/248/a2/census.md`＝**166 行／39,649 字节、零占位**，两枚提交 `06acaca5`→`139ca815`；⛔ AC 框一枚未碰，现读 **0 勾／12 未勾**。账 `A499`）

- **§⑦-2 我裁了，凭据是我自己现读那枚冻结件内部的词表**（`internal/panel/l2_grant_boundary_test.go`，**只读、一字不改**——"三枚冻结件一字不动"禁的是改，不是读；这一格只有能读它内部的人才答得了，所以归我）：那枚"授权门探测器"的判据形状＝**"被真守卫应答过的名字里，不许有任何一枚的拼写是授权形的"**（注释 `:1047` 逐字「Being route-shaped is only the filter. What decides anything is whether…」，函数 `poolJudgedByRealGuard` 在 `:1125`，红句在 `:1533`）。**"授权形"由两枚表决定**：词表 `grantRouteWords`（`:192-195`，**十枚**：`approve`/`approval`/`grant`/`allow`/`permit`/`ratify`/`authorize`/`authorised`/`decide`/`decision`/`verdict`）＋前缀表 `grantRoutePrefixes`（`:1316-1319`，**八枚**，全是 `panel*`）＋后缀表 `grantRouteSuffixes`（`:1324`，`""`/`.request`/`.now`）。⇒ **判定：`config.set`／`config.get` 两样都不落进那三枚表里 ⇒ 不会被误认成授权门，探测器照绿。**
  **两条命名禁令（写进落地腿派单，比判据本身更容易被忘掉）**：⛔ (a) 设置那一路的路由名**不许带词表里那十枚字的任何一枚**（尤其不许顺手写成 `config.approve`／`config.grant`／`config.decide`——那正是本仓反对的"面板侧批准"形状）；⛔ (b) **也不许为了躲开词表去动那三枚表**（冻结件，一字不动）；⚠ (c) 若产品将来真要一枚"面板侧改配置"的名字带授权味，那一格**不是命名问题、是 `AGENTS.md` §1.2 禁止清单那一档**，先停手上报。
- **逐格分类我照收（并把"哪半今天做"钉死）**：**Go 侧可独立做完并交凭据＝AC#3／AC#5／AC#7／AC#11**；**AC#1 可做，但落地即撞唯一一枚"无论怎么命名都会红"的冻结钉 ⇒ 走 `:2051`／`:2139` 那两枚**已具名解冻的锚**（台账里那枚解冻记录五样齐：文件／行／理由／边界／**撤销口令「248 别动冻结件」**；⚠ 我不在这里复述它是哪一枚 `A##`，落地腿自己 `grep -n "别动冻结件" docs/reports/pending-and-issues.md` 取那一条为准）；**边界原文＝"只许改这两处与其直连 helper"，改法只能是把"数一枚 4"换成能力形**——除那两枚锚之外这枚文件**一字不许动**。；**两栖、Go 半枚做不到 self-green＝AC#2／AC#8**（AC#2 被 `internal/panel/composer_dispatch_test.go:48` 那枚双向对账门洗红＝方法名册与页面名册要对得上，界面侧没动之前它必然差一枚）；**必须等界面侧或票 33 宿主＝AC#4／AC#9**（＋AC#10 的"页面可见"那一半）；AC#6 交落地腿自跑；**AC#10 那两个选择由我裁、不许实现腿选**（详见下条）。
- **三枚我复认成立的"不红"**（腿把前稿两处〔推〕自己改成了〔量〕，方向对）：T1 `git_test.go` 只抽 `panel.` 前缀 ⇒ 叫名 `config.*` 不红；T4 `cmd/wisp/panel_inbound_33_test.go` 那五枚钉的是"有生产调用者／`panel.review.allow` 不在名册⇒拒／空槽响亮拒"，**不钉"四槽里恰好三枚 nil"** ⇒ 给 `ComposerDispatch` 加一枚 `Config` 槽并赋非 nil 处理器**不触它**；T5 票 133 那枚派发闸门枚举的是 `func main` 的 argv 分支＋与 `usage` 常量双向对账 ⇒ 配置写走 panel inbound 就不动腿名册，**不红**（⚠ 唯一会红的做法＝新加一枚 `wisp xxx` 子命令，本票**不许**这么走）。
- **AC#10 那一格我这轮不裁**，理由具名：它两选一的另一半（把 `approval.New` 移进 `assembleRuntime`）会动票 246 AC#1 裁过的装配次序，而**此刻票 228/票 33 那一串还没落**，先裁就等于在没量的地上选路。⇒ 本轮只做"**⛔ 不许实现腿自己选**"这条纪律登记（撤销口令不变：「248 别动常驻门」）。
- **排程**：`248` 落地腿此刻**料齐**（四格可独立＋命名两禁＋T3 走那两枚锚），但仍**按住**——它要动 `cmd/wisp`（`newPanelInboundDispatch`）与 `internal/panel`，而 `33-r4` 正在 `cmd/wisp` 里跑测试并量面板耗时（同包写读重叠一律串行）。⚠ 另记一条争用事实：`33-r4` 在测冷/热启动耗时（冷 ≤1500ms／热 ≤200ms），我这轮**不派任何会吃 CPU 的腿**（本仓已知 `slo-full` 跑在同一台机器上会抢 CPU 那一族坑，不再自造一次）。队列：`33-r4` → `33-p1` 探针 → `33-r2` 功能腿 → **票 248 落地腿** → 票 228 AC#2／AC#11 那一发 → 197-r3 → 245 AC#6..9 → 224-r3 → 242 → 票 244 → 票 247。

## 编排者裁定与翻勾（2026-10-02 12:24:00+0800 现量，台账 `A530`；判语出自非实现者裁决腿 `248-v1c`，交件件我按盘上验过：`docs/evidence/s1/248-settings-write-path-v1.md` 279 行／55,821 字节／§1–§6 六节**没有一节是空的**／占位符 0／四枚号 `b86452c3`→`90576dce`→`8e8c4e7c`→`f1ebb711` 逐枚 `git log -1` 验存在；`git status --porcelain -- cmd internal tools scripts`＝**0 行**＝四发变异全还原了）

**翻了 6 格**（裁决腿判"成立"的那六格，我逐格对着它的凭据看，⛔ 没替它判过任何一格）：**AC#0／AC#1／AC#3／AC#5／AC#7／AC#11**。
**没翻的 6 格与去向**：
- **AC#6＝不成立 ⇒ 退回的是"证据件"那一格，不是产码**（裁决腿 §4 原话）。缺的是**实现者自己的门禁读数**：实现件 `docs/evidence/s1/248-settings-write-path-r1.md` 的 §4 至今是占位，而它 §7 表里 AC#6 那行写着"四数在 §4"＝**指着一节空的凭据**（这句由 `248-v1` 立案、`248-v1b` 续立、`248-v1c` 复认）。⇒ 派 **`248-r2`**：只准写那一枚证据件，四把尺**它自己现跑**（⛔ 不许拿裁决腿 §2 的读数当自己的——那会把"谁量的"洗混），并更正 §7 那句。
- **AC#2＝附条件成立**：Go 侧那两枚键声明要与页面侧对账才转绿 ⇒ 界面侧那一半**不在我射程**，写进本节由机主带给他用的那枚 agent（`frontend/**` 既不读也不引）。
- **AC#4＝不成立·挂账**（票面自己明令"票 33 的宿主真起来之前不许勾"）；**AC#9＝判不了**（禁令面）。
- **AC#8＝附条件成立**，并且裁决腿把我派单里那一问当场写进 §5 交裁——**我现在裁**：面板回执对 `[llm]` 说"重启才生效"而 `manager.go:281` 把它列在热应用表里，**这不是"没说'保存即生效'就过"**；AC#8 那种"不许说谎"的判据**遇到"说了另一句不成立的实话"同样算不过**。⇒ 该缺陷与票 255 新开那格 **AC#5 是同一块石头**（回执那句"什么时候生效"必须由同一份登记产出），**归 `255-r2` 落地**，本票 AC#8 那格**不因此勾**，等 `255-r2` 交完由非实现者一并裁。
- **AC#10＝不成立·两形都没落，选形归我 ⇒ 裁 ⓑ**（不移动 `approval.New`，改在设置那几项的回执里逐字写明"这两项只作用于跑任务的进程，常驻腿今天用常量 300s／3s"）。⛔ **不裁 ⓘ**，理由写在**新立的票 256** 里：ⓘ 要改的是票 246 AC#1 裁过的乙形次序，而"移动会不会破 `Confirming` 那一维"**今天没人量过**（裁决腿 §6 第 4 条自陈没跑，前腿与实现腿也没跑）⇒ 未量不裁、未定义即停；票 256 的第一格就是先把那一问量出来，量完才轮到我落那枚具名 `A##`。⚠ ⓑ 那句话的**产出方式**并入 `255-r2`（不许硬写死一句漂亮话，要与登记表同源），⛔ 不许只写"已保存"把差异藏起来。

---

## 6. `248-r2` 收档＋编排者现跑四把尺 ⇒ **AC#6 翻勾**（2026-10-02 14:1x，锚 `bf498283`）

**1. 交件核过（尺＝盘上）**：证据件 `docs/evidence/s1/248-settings-write-path-r1.md` §4 七小节填实、`grep -c '（待填）'`＝0；四发号都存在（`e7521610`→`69165e09`→`e7599f32`→`0b5b2a0d`），去重路径名册＝**只有那一枚 md**（票面／产码／测试／台账一枚未碰）。§7 那句"四数在 §4（build/vet/d22scan/gofumpt 全 rc=0）"的**原句逐字留在 §4.7**、作废理由两条我认：它指的 §4 当时是空的，而"全 rc=0"那个值两枚独立腿都给不出来。

**2. 编排者现跑（14:0x–14:1x，HEAD `bf498283`；可比性尺＝`git diff --name-only e7599f32..HEAD -- cmd internal tools` 只有两枚 docs 路径 ⇒ Go 面零漂移）**：
- `GOFLAGS= go build ./...` ⇒ **rc=0**（输出 0 字节，`14:04:01`）。
- `gofumpt -l cmd/wisp internal/config internal/panel` ⇒ **rc=0、list-count=0**（`$(go env GOPATH)/bin/gofumpt.exe`，不在 shell PATH 上）。⚠ 口径差异具名：`248-r2` 那一格量的是"本腿动过 0 枚 Go 文件＝**无射程**"，我量的是"本票动过的三目录格式化没有"——**两道命题不同**，我这一发补的是后者。
- `go test ./cmd/wisp ./internal/config ./internal/panel -count=1` ⇒ `internal/config` **ok 1.006s**／`internal/panel` **FAIL（4 枚红，逐名与 `248-v1c`、`248-r2` 三发同册：`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`）**／`cmd/wisp` **ok 198.911s**。
- `./tools/d22scan/d22scan.exe` ⇒ **rc=0**，逐字 `clean - no D22 ban violations`、承重行 `examined 262 production Go files under internal/ and cmd/`。
⇒ **票面点名那四发的终态都齐了，且不是同一枚腿独一份的自述** ⇒ AC#6 勾。⚠ 这枚勾的**边界**：AC#6 要的是"四数照抄终态"，它**从不要求全绿**——那 4 枚 panel 红的归因在别的票面上（两界对账门＝界面侧那一半、`tokens_fourway`＝已知常红、第二样式源＝非本票地界），⛔ 不许被读成"票 248 全绿"。

**3. 那枚第五红我今天多一个读数**：`TestTicket223ModeLooseningChangesTheRunningModeAfterAllow`（红句落点 `cmd/wisp/config_reload_223_test.go:445`）——`248-r2` 整包跑到它红、隔离复跑 3 发全绿；**我这一发整包（198.911s）它绿**。⇒ 目前 5 发读数里"整包红"只出现 1 次，定性继续按**既有间歇族**（该文件最后一次被改是 `5a755c3c`，标题逐字"AC#5 间歇红改有界轮询"），⛔ 不读成回归、⛔ 不读成"永远不红"。

**4. 我这一轮裁掉的两问（都归我、都不摆给机主）**：
- **`§6-5` 的边界：实现腿把 `fs.` 连同 `net./plugins./privacy./models./audio.` 一起列进拒写族（`config_handlers.go:108`）⇒ 我判【维持保守形，不改】**。理由：`fs.allowed_dirs` 是**判级输入**（批准的 L2 卡本来就能读根外，把它做成执行时硬边界是另一码事、已定过不许顺手做），从界面写它＝**让模型可触面在没有 L2 卡的情况下变宽**，与 J5 锁 `risk.*` 同族。窄不出事故、宽出事故，方向不改。⇒ 要放宽必须**具名解冻＋落 `A##`**；机主没要过"从界面改允许根"，这一格不进他的清单。
- **`§5-10` 的价格只做了 5 枚叶子中的 2 枚（in/out）⇒ 我判【够，本期不补】**。判据原文只写"价格"；另 3 枚（`cached`/`audio_in`/`audio_out`）零覆盖记成**残余**，等真接计费那台件时一起补，⛔ 现在不派腿（多开＝多改契约面，实现腿那句是自划线、不是我漏的洞）。
- **`§5-1` 那枚新仪器（`guardRosterOf` 假设路由常量都在 `bridge.go` 一枚文件的 const 块里）我没判**——它正是"仪器自己会不会假红、红了的是冻结件"那一类，交回**下一枚非实现者验收腿（`248-v2`）**，与 §5-9 那枚控制支的读反风险同批。⛔ 我不既当写这批派单的人又当判它的人。

**5. 一处过期指认（票面自己的话，不改原句、就地记）**：`§6-2` 报——票面 §5.4 那句"宿主还没接进常驻那条腿（`grep PanelManager`＝0 命中）"**今天不成立**：`cmd/wisp/resident_windows.go` 命中 2 行、`cmd/wisp/panel_resident_windows.go` 命中 6 行，且 `newResidentComposerDispatch`（`:166-170`）复用的正是本票动过的 `newComposerDispatchChain`。⇒ **票 33 那半边比票面写的更靠前**；但"窗口建得出"≠"机主点得出设置"，缺的是**页面产物**那一格（`frontend/dist` 只有 `.gitkeep`），仍归 AC#4／AC#9。

**6. 仍未翻的 5 格与去向（14:1x 现量：未勾 5／已勾 7）**：AC#2（界面侧声明 `credentialState`/`credentialKnown`，由机主带话）／AC#4（真机那一发，缺页面产物＋入口）／AC#8（归 `255-r2` 同源登记表）／AC#9（禁令面，判不了）／AC#10（等 `256-a1` 量完才落具名 `A##`）。⇒ **本票今天不 `-done`**。

**7. AC#10 定案（14:3x，`256-a1` 的量回来了）**：ⓑ／ⓘ 那一问**定＝ⓑ**——⛔ 不移动 `approval.New`，改由设置回执逐字写明"`[risk]` 那两项只作用于跑任务的进程，常驻腿今天用常量 300s／3s"，且**必须由登记表同源产出**（归 `255-r2`＝与 AC#8 同一块石头，本格仍不勾）。**理由不是我的判断，是现量**：ⓘ 若做，门的存在时刻变成"仅当这发常驻进程有交互控制台"（`cmd/wisp/resident_task_source_windows.go:218→:230 return nil→:265`）⇒ **双击／Explorer 拉起那一支永远没有门、球永不进 `Confirming`**——不是变弱，是载体不存在。五样齐的具名裁定在台账 **`A534`**，全部现读凭据与"我票面被推翻的三句"在**票 256 第 7 节**；撤销口令**「248 AC#10 改 ⓘ」**。

**8. 残余五格分诊收档（10-03 09:5x，只读腿 `248-c3`＝`.scratch/wisp/probes/248/c3/census.md` 143 行／29,246 字节，七枚 pathspec commit 最末 `203014a1`，占位符 grep＝0；编排者复量四处）**：⛔ 本格只改排程与措辞，**一枚 AC 框都不翻**。
- ★**推翻我自己票面里"AC#4＝窗没接"那一读**（我复量：`cmd/wisp/resident_windows.go:148` `panel = startResidentPanel(rt.Registry, rp)`＝1 枚生产调用、`:164`（⚠ 腿报 `:163`，差一行）把 `panel.RequestToggle(via)` 交给球宿主 ⇒ 开窗那一跳**今天已接线**，`bringUp` 由热键／托盘**传递可达**、"nothing is created until the user asks"＝设计上延迟到手势，不是没接）。AC#4 真缺的是三样叠一格：①真机才知"窗真开时会不会撞嵌套消息泵重入"，②前端 dist 产物（`git ls-files frontend/dist` 我复量＝**1 行**、逐名 `frontend/dist/.gitkeep`＝AC#9 同一块石头），③一次真机点击＋下次重启真用上。
- **AC#8／AC#10 的凭据归 `255-r2`，不归票 33、不归前端**（我复量：`cmd/wisp/panel_config_store.go:228` 与 `:275` **无条件** `res.Tier = panel.EffectiveRestart`，前面零按键判断；`internal/panel/config_handlers.go:199-200` 的 `EffectiveNow`／`EffectiveNextTask` 非测试**写者 0 枚**；`internal/config/tiers.go:30 "llm":"hot"`／`:34 "panel":"hot"` 在）。⇒ 这两格翻勾的起跑判据＝`255-r2` 把 `:228/:275` 的硬填换成 `config.TierOf(...)`。
- **AC#2 卡在一枚"读前端 TS 的差集尺"**（我复量机制：`internal/panel/composer_test.go:48 TestComposerContractTypesMatchFrontend` 用 `os.ReadFile` 读前端类型文件做双向减集）；Go 半边四样全在（`internal/panel/composer.go:268-269` 两枚 json tag `credentialState`／`credentialKnown`、`pump.go` 的 reader 钩子、`cmd/wisp/run.go` 那一枚生产接线＝唯一调用点）。⚠ 界面侧那半**只写进票面由 owner 自己带给他用的那枚 agent**，⛔ 本编队永不转达。
- 一句〔腿报，未复核〕留着：`TierOf` 的生产调用者枚数（起手段 0／现段 1）由 `255-c2` 量，等 `255-r2` 终态我整包复跑再认。

## 8. AC#8 第二半的**两道外部前置**（2026-10-09 11:5x 编排者追加；⛔ 本节不翻任何框、⛔ 不改 §要建什么 里 AC#8 的判据原句）

来路＝只读腿 `248-w1`（件 `.scratch/wisp/probes/248/w1/00-findings.md` 126 行，commit `9b029e61`；编排者对拉过它的字样族读数——在它自己的锚 `c2b422a4` 上**逐枚精确复现**，见台账 `A769`）。

第 126 行那节我把 AC#8 判成「**归 `255-r2` 落地**，等 `255-r2` 交完由非实现者一并裁」。**这句现在不完整**，照实补齐（⛔ 不是推翻，是少说了一半）：

- **`255-r2` 只闭合第一半**＝"回执那句『什么时候生效』必须由同一份登记表产出"（票 255 `AC#1`＋`AC#5`，同一块石头）。
- **第二半「"要重启"必须到达页面可见面」今天有两道外部前置，`255-r2` 一道都不解**：
  1. **Go→页那一跳的返回值没人接**：Go 侧已经把那句话交回绑定（`cmd/wisp/panel_host_windows.go:802-804` 逐字 `reply, _ := m.dispatchRaw(ctx, raw)` / `return reply`），断点在页面侧——`frontend/src/lib/panel.ts:211` 的 `sendRequest(...): void` 不取返回值、`:140-142` 桥接口只有 `postMessage(message: string): void`、`internal/panel/composer.go:58-91` 的快照没有回执字段。⇒ **这一半的正解落在 `frontend/**`，按本编队现行口径（`A102` 族＋票面规则 7）⛔ 不由我写**，界面侧需求只能写进票面/台账由机主自己带给他在用的那枚 agent。
  2. **Go 侧"请求↔回执"配对这一块本身还没写**：票 35 的 C17 面 `:44` 那格已登记为**"一块没写"**（`ComposerDispatch` 结构体零请求态容器；`internal/panel/bridge.go:80 NewRequestID()` 非 test 调用者 **0**；答复回程走 `msgcb → Eval`、不走包裹器返回值）。⇒ 归属＝**票 35**，⛔ 本票不重复施工、⛔ 不在这里另立一张（同一物理缺陷链只记一次）。
- **因此 AC#8 那格现在的诚实状态**＝"⛔ 不勾，且**不能只等 `255-r2`**"：登记表那一半做完之后，这格仍会因"到不了页面可见面"而不成立。谁把哪一半补上，上面两条具名。
