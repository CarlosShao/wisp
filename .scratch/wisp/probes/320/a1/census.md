# 320-a1 只读设计核查 —— 模型页 v3 Go 侧（字段集/新方法名/执行面/面序/风险）

- 腿：`320-a1`（只读，零产码）；锚 `36cf5d06`（`git rev-parse --short HEAD` 起手现读，与派单无漂移）。
- 行号＝`36cf5d06` 的**读数**（会烂，逐处自己现读；token 才是判据，本仓 255 AC#6 口径）。
- 本刻在飞 `319-r2`：现读 `git status --short -- cmd/wisp internal/proc internal/panel internal/config` **零改动**（它还没落笔；它的面＝319 测试件＋必要时 319 生产件，与我列的文件零同名）。
- 现跑门：`sh scripts/d22scan.sh` = clean（internal/+cmd/ 生产 Go 277 枚）；`sh scripts/check-path-length-budget.sh` = GREEN（60=60，最长 185）。

## 0. 结论先行

1. **三个写入字段全部落在 313 名册**（`internal/panel/config_write_313.go` 的 `fields313`＋twin 三 map）：`provider_billing` / `provider_row_create_custom`（名待裁）/ `model_row_create`。248 七枚**一枚都不许动**（两处 `want` 钉死在 `config_route_248_test.go:271-277`＋`config_313_test.go:44-59`）。
2. 值格式：`provider_billing` 值＝`plan|pay-per-token`（字面量出自 `schema.go:80-81` 的 `BillingPlan`/`BillingPayPerToken`）；`provider_row_create_custom` 值＝`<protocol>,<base_url>`（逗号连接，链值已有先例）；`model_row_create` **不带值**（条目从 `default:"..."` 标签播种，同 `settings_313.go:169-176`）。
3. 新方法名候选甲＝`config.models.list`（推荐；`config.*` 命名空间已存在、不带 `panel.` 前缀）。**代价是全票最重的一枚**：新 Method 常量 ⇒ `frontend/src/lib/panel.ts` 必须同批出现该字面量（`bridge_test.go:184` 的 `TestEveryWhitelistedInboundNameIsAlsoSentByThePage` 逐枚读 bridge.go 的 Method* 常量），而 `frontend/**` 是本编队零写面 ⇒ **需要机主前端 agent 落一行**，否则 `./internal/panel` 必红。
4. 拉取执行面＝**已有生产路径**：`llm.DiscoverModels`（`internal/llm/discover.go:42`）＋`secret.Store.Resolve`（`internal/secret/store.go:89`，env/dpapi 两族都在它里面）；`cmd/wisp/providers.go:116-149` 是现成样板。零新网络码、零落盘、零明文。
5. AC#5 我裁：**不落**预设默认模型目录（`defaults.go:12-14` 逐字写着 presets 只供 protocol/base_url，模型来自用户配置或 /v1/models）。
6. AC#3 我裁：**建条目单独一枚新字段**＋复用既有 `model_context_window/price_in/price_out` 三枚；不并成一枚（并一枚＝一字段写多键，违 `settings.go:25-28` 的 one-key-per-call）。
7. 面序：`319-r2`（在飞）→ `320-r1`（先 `internal/config`＋`internal/panel` 半，再 `cmd/wisp` 半）→ `135-r1`；`cmd/wisp/**`、`internal/panel/**` 共享编译面**每波只许一枚写腿**。

## 1. 字段集定稿（AC#1/#2/#3）

| 字段名（定稿） | 值格式 | NeedsProvider | NeedsModel | TakesValue | 落哪张表 | config 侧 setter（新） | ApplySetting case |
|---|---|---|---|---|---|---|---|
| `provider_billing` | `plan` 或 `pay-per-token`（枚举） | ✓ | ✗ | ✓ | `fields313`＋`fields313NeedsProvider`＋`fields313TakesValue` | `SetProviderBilling(provider, billing)`，键 `llm.providers.<p>.billing`，写前枚举校验 | 新 |
| `provider_row_create_custom`〔名待裁，见 §3〕 | `<protocol>,<base_url>`；`configModel` 可选带首枚模型 id | ✓ | ✗（可选照走 `SettingWrite.Model`，`config_handlers.go:386-411` 现读如此） | ✓ | 同上三 map（takesValue ✓） | `CreateCustomProviderRow(name, protocol, baseURL, model)`：拒预设名（指回 `provider_row_create`）、拒空 base_url、拒枚举外 protocol、拒已存在行；一次 guarded 写（行＋可选模型） | 新 |
| `model_row_create` | 无值 | ✓ | ✓ | ✗ | 同上（provider+model 两 map） | `CreateModelRow(provider, model)`，键 `llm.providers.<p>.models.<id>`；ModelSpec 由 `applyDefaults` 播种（Enabled=true / billing=pay-per-token）；未知 provider / 已存在条目一律拒 | 新 |

- 逐落点（现读）：`internal/panel/config_write_313.go:58-80`（名册＋三 map）、`:85/:92/:104/:113/:117/:121`（accessor 与四枚 union helper 都要加第三表追问）；`internal/panel/config_handlers.go:425`（`knownWritableField` 加第三表）；`cmd/wisp/panel_config_store.go:202-255`（ApplySetting 三枚新 case）；`internal/config/` 新增 `settings_320.go`（三个 setter，骑 `settings.go:259` 的 `writeOneKey`）。
- **必须同批改的测试（逐文件）**：
  - `internal/panel/config_313_test.go`：往 `fields313Cases()`（`:28-37`）**追加**三枚 case（否则"每枚恰好到达腿一次"那半不覆盖新字段）；`:129-143` 的命名禁令测试读 live 名册，自动覆盖；`:262-324` 载荷键钉**不要动**（拉取结果不进 config.get 载荷）。现读**无任何 313 名册计数钉**（`AdditionalWritableFields` 只被三个文件按 accessor 读）。
  - `cmd/wisp/panel_config_313_test.go`：`Test313WritableUnionEqualsTheTwoRosters`（`:225-286`）读 live accessor＋逐枚过门 ⇒ 自动通过（三枚新字段在"provider+model+value=x"下都过门；具体被腿拒是 spy 地盘）。**若改成第三张表**，则该测试 `:249` 的 `want` 必须加第三 accessor，否则拒绝文与 `want` 不等 ⇒ 红。
  - `internal/config/settings_313_test.go`（或新 `settings_320_test.go`）：三枚 setter 的落地/拒因测试；`CreateProviderRow` 的"拒非预设"钉（`:80-111`）不许放宽。
  - `internal/panel/config_route_248_test.go`：**零改动**（前提＝不加 `ConfigStore` 方法，见 §4）。`cfg248Leg`（`:26-56`）是接口实现者，接口长一枚就编译红。
  - `cmd/wisp/panel_config_store.go` 自己的用例（新文件 `cmd/wisp/panel_config_320_test.go` 或 313 测试尾追加）。

## 2. 新方法名（AC#4）

**候选**（三枚都过两枚词表，见下）：

- **甲（推荐）`config.models.list`**：`config.*` 已是设置页命名空间（`bridge.go:66-67`）；读"某 provider 的模型目录"；无 `panel.` 前缀。
- 乙 `models.list`：**SPEC-08 §5.2:172 已carries**（`models.list / models.delete`，invoke，无需原生授权）——但那行的语义更像 C29 本地签署模型仓（`cmd/wisp/models.go` 的 `wisp models list`／D26），拿它装"问供应商 /v1/models"语义相撞 ⇒ **要裁**，不推荐直接用。
- 丙 `config.models.discover`／`config.models.fetch`：与甲同形，动词贴近 `llm.DiscoverModels` 自己的词。

**注册四点＋两点（一枚新名要动的全部地方）**：
1. `internal/panel/bridge.go:41-68`（Method* 常量）＋`:146-152`（`knownComposerMethod` 两处 case 是**同一条逗号子句**，照 248 形状接排）；
2. `internal/panel/composer_dispatch.go:181-216`（case 调 `d.Config`，**写字面量会撞** `TestDispatcherSpellsNoRouteLiteralOfItsOwn`：`composer_dispatch_test.go:374`，只用常量）；
3. `internal/panel/config_handlers.go:324-335`（`HandleConfigRequest` 的 switch 加一枚分支）；
4. `internal/panel/inbound_roster_253_test.go`（**必改**：`wantFullInboundRoster253` `:70-77` 加名、`wantFullInboundRosterSize253` 6→7（`:56`）、`wantNonPanelPrefixedInbound253` 2→3（`:63`））；
5. `frontend/src/lib/panel.ts`：**必须出现该字面量**（`bridge_test.go:184`）；且前端**只能复用既有 `sendRequest`**——`bridge_test.go:161-164` 钉 `postMessage` 调用点恰 2 处，`composer_test.go:535` 的 renderer 门只放行 `panel.*` 字面量（`routeLiteralRe` `:454`）⇒ 新名用 `config.*` 反而不会被它拦。
6. ⛔ **不许用 `panel.` 前缀**：`git_test.go:385-391` 钉死 bridge.go 里 `panel.*` 字面量恰四枚（extractor `:409` 只认 `panel.`）。248 的两枚名正是为躲它才叫 `config.*`（`bridge.go:50-60` 逐字）。

**载荷/回执形状**：散文一句＋**一枚标记行**（照 `SettingsPayloadMarker` 模式，`config_handlers.go:546-555` 的注释逐字要求"不带点、不带 `panel.` 前缀"⇒ 候选 `wisp-models-list:`），JSON 例：`{"provider":"stepfun","models":[{"id":"…","ownedBy":"…"}]}`；失败＝散文里给 D37 分类（`llm` 已分好类），**不回显密钥**。⛔ 不往 config.get 载荷加键（`:262-324` 键集冻结）。

**冻结件口径逐字核**（`internal/panel/l2_grant_boundary_test.go`，零改动）：
- 路由词表 `grantRouteWords` `:192-195`（11 枚，含 approval/grant/allow/decide/verdict…），匹配前**去掉 `_`/`-`/空格/点**（`carriesGrantWord` `:242-250`）⇒ `config.models.list` 归一化 `configmodelslist`，**0 命中**；
- 字段词表 `grantFieldWords` `:183-188`（18 枚，多 `outcome/override/bypass/allowonce`）只扫**"JSON 解码落点类型"的 JSON 键**（`scanGrantBoundary` `:1163-1208`＋`inboundEnvelopes` `:1006-1031`）⇒ 三枚新**字段名是 Go 字符串常量、不是 JSON 键**，不进它；
- 该仪器还要求 `answered == 声明的 Method* 常量`（`:1293-1303`）与 route 池双向相等（`:1125-1142`、`:1149-1157`）⇒ 新常量走常量即自动成立。
- `config_313_test.go:129-143` 的"字段名不许带审批词/不许进 locked 家族"（11 词＋`refusedLockedFamily` `:116-137`）：新三枚名 0 命中；locked 前缀是 `risk./fs./net./plugins./privacy./models./audio.`（**带点**）⇒ `model_row_create` 不落 `models.` 前缀，安全；**别给新字段起带点的名**（`models.x` 会红）。

## 3. 执行面（AC#4 的"最小且不泄密"）

- **密钥**：`secret.Store.Resolve(ref)`（`store.go:89-123`）——`env:NAME` 走 `os.LookupEnv`，`dpapi:<id>` 读 blob＋DPAPI unprotect；解析不了**不发明文、不回退**（`providers.go:124-133` 是"解析不了就不发未认证请求"的样板）。常驻进程手里已有同一个 store（`cmd/wisp/panel_inbound.go:273-279` → `newConfigStore`；`run.go:434` 同源）。
- **HTTP 层**：**不落 `internal/panel`**（该包按纪律只声明 socket，不 import `internal/llm`）；执行面落 `cmd/wisp` 的设置腿：`llm.DiscoverModels`（`discover.go:42-96`，GET `<base>/models`、有 key 才带 Bearer、错误经 `ClassifyTransportError` `:68`／`NewHTTPError` `:73` 分类、错误体 8KB／响应 4MB 封顶、只排序 id、**零落盘**）；client＝`llm.NewDefaultHTTPClient(cfg.Net.Proxy.Mode, cfg.Net.Proxy.URL)`（`errors.go:358`），与 CLI 同源。
- **超时**：`context.WithTimeout`（context deadline＝单调；禁令禁的是 `time.Now()` 差值——`tools/d22scan/main.go:785-790`）＋`DiscoverOptions.Timeout`；⛔ 不许起 goroutine：d22scan ban#1 连**具名 `go f()`** 都扫（`main.go:721-733`），合法 spawn 只有 `observe.Registry.Spawn`。
- **socket 形状**：`ConfigWriteHandler` **另加一枚可选接口字段**（如 `Fetch ModelFetcher`，nil＝按名字拒绝），⛔ 不往 `ConfigStore` 加方法——那会连带 `internal/panel/config_route_248_test.go` 的 `cfg248Leg` 与 `cmd/wisp/panel_config_313_test.go:203` 的 `rosterSpy313` 编译红（248 面本来可以零改动）。
- ⚠ **未定（两条候选＋判据缺口）**：`dispatchRaw` 在 WebView2 回调线程上**同步**跑（`panel_control_windows.go:711-800` msgcb→callbinding→绑闭包），一次同步 GET 会在超时区间里冻窗。甲＝短预算（≤10s）＋页面 loading 态；乙＝挪到有回复通道的 worker（今天 RPC 管道没有异步回执面）。缺的判据＝票面没给超时值、也没有"冻窗可接受上限"的裁决。
- ⚠ **未定**：`DiscoverModels` 会把上游错误体的**首行**带进错误串（`discover.go:72-77`），而回执是页面可见面——若供应商把自己的 401 体回显成含 key 的文本，理论上会到页面。甲＝照 CLI 先例原样；乙＝设置腿在回执前截断/只留分类。缺的判据＝没有裁决说必须净化。密钥本身**不经过**任何回显路径（`llm` 只把它放进请求头）。

## 4. AC#3 与 AC#5 取舍

- **AC#3**：建条目＝新字段 `model_row_create`（§1 第三行）＋**复用** `model_context_window`/`price_in`/`price_out`（三枚已存在、已钉、按"已有条目"语义设计，248 钉名不动）。并成一枚＝一字段写 4 键（违 one-key-per-call，`settings.go:25-28`）且要新造值语法 ⇒ 不推荐。
- **AC#5**：**不落**。理由三条（现读）：① `defaults.go:12-14` 逐字"presets supply protocol/base_url only；model entries are user-configured or discovered via /v1/models"；② 落预设目录要动 `loader.go:199` 的 load-time preset merge ⇒ 改变**所有既有配置**的 `validateCatalog` 结果与 `Configured` 语义（面比本票大得多）；③ 静态名单会随供应商上新而腐烂，"拉"（AC#4）就是真名单。若机主坚持离线优先，须另裁"仅在行内无模型时合并"的规则（契约味，属人工批准档）。
- 顺带现读：`llm.ImportDiscovered`（`discover.go:102`）**今天零生产调用者**；若将来做"整单导入"，它是第一枚生产 caller——本票 AC#3 是逐枚录入，**不顺手接线**（会扩面）。

## 5. 面序与冲突

- 逐枚文件面：
  - `internal/config`：新 `settings_320.go`（＋测试件）。与 319/135 **零重叠**。
  - `internal/panel`：`config_write_313.go`（名册＋helper＋注释改口）、`config_handlers.go`（门＋handler 分支＋新标记/载荷函数）、（仅新方法时）`bridge.go`、`composer_dispatch.go`、`inbound_roster_253_test.go`；测试：`config_313_test.go`（追加）、新 `config_320_test.go`。
  - `cmd/wisp`：`panel_config_store.go`（3 case＋拉取腿）、新 `panel_config_320_test.go`；仅新方法时 `panel_inbound.go`（socket 装配）。
  - `frontend/src/lib/panel.ts`（仅新方法；**机主外部 agent 的面**）。
- **串行序建议**：① `319-r2` 落地（它现在占 `cmd/wisp` 319 面）；② `320-r1` 一发做完 `internal/config`＋`internal/panel` 两半（这两包不碰 319/135 的面；⛔ `internal/panel` 也算共享编译面，别再并行第三腿）；③ `cmd/wisp` 半＋新方法前端字面量（前端那枚不落则 `./internal/panel` 红——**把这一枚当交付前提写进派单**）；④ `135-r1`。⛔ 320 写腿与 319-r2／135-r1 **不许同时在 `cmd/wisp` 写**；`go test ./cmd/wisp` 同一时刻只跑一枚。

## 6. 风险清单（逐条"会不会被本票触发"）

| 项 | 现读判断 |
|---|---|
| 冻结三件（`internal/panel/tokens_fourway_test.go`/`internal/panel/l2_grant_boundary_test.go`/`internal/perm/ticket90_persist_test.go`，名册出自 223/224/225 票面）＋票 320 点名的那枚 | **预期零改动**：新名/新字段名对三枚词表 0 命中（§2 已算）；该仪器只要求"新常量进 case 便双向相等"。若环境变了让它红 ⇒ **停手上报**，不许动它 |
| 审批词汇扫描（冻结件＋`config_313_test.go:129`＋`config_route_248_test.go:259`） | 0 命中（11 词 vs 三枚名，逐一算过）；⚠ 起名时避开 `approve/allow/decide/verdict` 与 `outcome/override/bypass` |
| 313 名册测试的读法 | `fields313Cases()` **手写追加**；`config_313_test.go` 无计数钉；`cmd/wisp` 并集钉读 live accessor ⇒ 自动。**改成第三张表**才会多两处改（拒绝文 `want` ＋ helper） |
| d22scan ban#6 | 射程＝`frontend/` 文本（`main.go:417`），整体是"面板侧 allow 门"词形；本票 Go 侧不写 `frontend/`；前端那枚字面量（若走新方法）是 `config.*`，不进它的 `panel.*` 词形 |
| ban#7 | 射程＝`internal/tools/`（`:433`），零碰 |
| ban#8 | 射程含 `internal/`＋`cmd/` 的 Go（`emojiScopes` `:577-584`），**注释豁免、字符串不豁免**⇒ 新产码/测试注释与字符串零 emoji（含 `≤`/`✓` 这类）；现跑基线 clean |
| bans #1–5／#9 | internal/＋cmd/ 生产 Go 全走（`:409-413`）：⛔裸 `go`（含具名调用）·⛔`filepath.Clean/Abs` ·⛔明文密钥字面量 ·⛔`time.Now()` 差值做超时 ·⛔镜像＋哈希同现；#9 phantom-citation 也扫 internal/＋cmd/ 的**生产**注释 ⇒ 新注释不许引不存在的仓内路径 |
| 长度门 | 现跑 GREEN（60=60，最长 185）；新文件名短；票文件名 58 字符＜100 帽 |
| `cmd/wisp/config_readers_255.go`＋`config_receipt_255_test.go:256` | 拉取腿落 `panel_config_store.go`（已是 `LLM` 行的读者）；`[net]` **不在** `sectionReadSites` 覆盖内 ⇒ 预期零改。若换去新文件读被覆盖的段 ⇒ 必改该 roster。⚠ `llm` 行 token `[ChatModel:      cfg.LLM.Roles.Chat.Model]` 在该文件里**必须保持唯一**（行号漂移无碍） |
| 248 钉 | 新字段**不许**进 `allWritableFields`/`WritableFields`（两处 `want` 恰七枚） |
| config.get 载荷键钉 | 拉取结果**不进** config.get 载荷（`:262-324` 逐键冻结） |

## 7. 未定／停手上报（⛔ 我不猜）

1. **AC#2 形状**：甲＝扩 `provider_row_create`（值在非预设时必填＝protocol，base_url 还得第二写 ⇒ 丢崩溃原子性、且"值必填"是条件式、装不进布尔 map）；乙＝新字段 `provider_row_create_custom`（值＝`protocol,base_url`，一次写）。我推荐乙，但"新字段名/值语法"是接口面 ⇒ 交编排者/机主裁。判据缺口：票面只说"值格式/字段名由设计腿定稿"，没给单选裁决面。
2. **AC#4 超时值与同步/异步**（§3 两条候选）。
3. **上游错误体首行进回执**（§3）。
4. **`models.list` 与 SPEC-08 §5.2:172 的语义相撞**：若采候选乙须先裁"那一行属于谁"（本地签署仓 vs 供应商 /v1/models）。
5. **AC#5 若改成"落"**：合并规则（仅空目录时？）属 [llm] preset 继承语义，建议另裁。

---
具名没做：没跑任何 Go 测试/构建（只读腿，四门只跑 d22scan 与长度门两把**只读**脚本）；没碰 `frontend/**`（读了 `panel.ts` 三处字面量以核 `bridge_test.go` 的门）；没动台账/票面框/`Status:`；没碰 313/248 的任何文件。
