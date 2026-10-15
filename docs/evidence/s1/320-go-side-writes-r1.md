# 320-r1 落地证据 —— 模型页 v3 的 Go 侧三个写入（计费模式／自定义供应商行／模型条目）

- 腿：`320-r1`（实现者；判语归非实现者）。起手锚 `33e429c1`；派单存档 `.scratch/wisp/probes/orch/20261015/txt/320r1-dispatch.txt`（逐字同收到的那份）。
- 依据：票 `.scratch/wisp/issues/320-model-page-v3-go-side-three-writes-and-a-fetch-port.md`；设计定稿 `.scratch/wisp/probes/320/a1/census.md`；编排者五裁 `A972`（AC#2 采乙形；AC#5 不落）。
- **本腿不含 AC#4**（"拉取模型列表"端口＝`320-r2`，待前端字面量先落）。
- 本腿四笔（逐格一笔＋一笔测试补强，均显式 pathspec、只 commit 不 push）：

| sha | 内容 |
|---|---|
| `a2d29798` | AC#1 `provider_billing` |
| `a920c4d4` | AC#2 `provider_row_create_custom` |
| `84fc3f20` | AC#3 `model_row_create` |
| `527d483f` | AC#1 补强：拒因分层断言（追加笔，非 amend） |

## 0. 结论先行

1. 三个写入全部落 **313 名册**（`internal/panel/config_write_313.go` 的 `fields313`＋twin 三 map），248 的七枚一枚未动（两处 `want` 逐字核过）；新字段是**字符串常量＋`ApplySetting` case**，`ConfigStore` 接口零新方法 ⇒ `internal/panel/config_route_248_test.go` 整件零改动。
2. 三个 setter 落 `internal/config/settings_320.go`，全部骑 `settings.go` 的 `writeOneKey`（one-key-per-call；唯 AC#2 是"一行＋可选首枚模型"的**一次** guarded 写，与 313 `CreateProviderRow` 同形）。
3. 枚举/协议词表**读 `internal/config` 自己的 `billingModes`/`protocols`**（不抄字面量）；写前预检 + `validate()` 双闸。
4. `config.get` 载荷键集**未动**（拉取结果不在本腿）；inbound 名册/`bridge.go` 零改动（无新方法名）。
5. 门禁：`go build ./...` rc=0、`go vet ./...` rc=0、gofumpt（本腿六件）零输出、`sh scripts/d22scan.sh` clean、长度门 GREEN；`./internal/config` 全包 rc=0；`./internal/panel` 的 5 枚失败与**把本腿文件 overlay 回 `33e429c1`** 后逐枚同名（下方 §4），非本腿引入。

## 1. AC#1 `provider_billing`（sha `a2d29798`，补强 `527d483f`）

- **改了什么／落点**：`internal/panel/config_write_313.go`（`FieldProviderBilling` 常量＋`fields313`＋`fields313NeedsProvider`＋`fields313TakesValue`）· `internal/config/settings_320.go`（`SetProviderBilling`，键 `llm.providers.<名>.billing`，写前 `isBillingMode` 枚举预检，拒因句点名 billing）· `cmd/wisp/panel_config_store.go`（`ApplySetting` 第三枚 case）。
- **两向三态读数**（`.scratch/wisp/probes/320/r1/01-ac1-two-way.txt`）：
  - 接受：`plan` 落盘并 `LoadFile` 逐字回读；`pay-per-token` 随后是**真变更**（文件已载 plan），回读对上。
  - 拒绝一：`flat-rate` 被**字段自己的预检**拒（句子含 billing、**不含** `第 3 种拒因`），文件字节不变。
  - 拒绝二：未声明 provider（`ghostco`）＝第 2 种拒因（`行不存在`），文件字节不变。
  - 对照：两拒之后同一字段仍写得进（真腿 `cmd/wisp` 用例复证同形）。
- **变异红名逐字**（`04-mutations-red-names.txt`）：
  - M1（把 `p.Billing = v` 改成 no-op）⇒ `--- FAIL: Test320SetProviderBillingWritesTheEnumLeafAndRefusesOthers` ＋ `cmd/wisp` 侧 `--- FAIL: Test320ProviderBillingThroughTheRealLeg`（收据"没有任何键路径发生变化"逐字）。
  - M1b（把预检改成只看非空）**首读仍绿**——`validate()` 第二道闸也拒同一值且句子含 billing；这是 `527d483f` 加分层断言的原因，加后逐字红 `settings_320_test.go:60: the out-of-enum value reached the whole-file gate instead of the field's own pre-check`。
- **门禁 rc**：build 0 / vet 0 / gofumpt 零输出 / d22scan clean / 用例 `-run Test320SetProviderBilling` ok。

## 2. AC#2 `provider_row_create_custom`（sha `a920c4d4`）

- **改了什么／落点**：`config_write_313.go`（常量＋三 map，值桶标 takesValue）· `settings_320.go`（`CreateCustomProviderRow(name, protocol, baseURL, model)`：先按 `applyDefaults` 播默认再覆写 protocol/base_url，可选首枚模型同批；**四拒**＝预设名〔具名指回 `provider_row_create`〕／空 base_url／枚举外 protocol／行已存在；空名＝第 2 种拒因）· `cmd/wisp/panel_config_store.go`（case＋`splitCustomRowValue` 按**首个逗号**切两半并 trim，切法归腿、与 313 `splitChainValue` 同族）。
- **两向三态读数**（`02-ac2-two-way.txt`）：接受 `local-gw`（节＋protocol＋base_url＋首枚模型落盘；收据一枚键 `llm.providers.local-gw`）；接受"只建行不带模型"（对照）；五条拒绝逐条具名且文件字节不变（子测逐枚 PASS 逐字在件内）；真腿值 `openai-chat, http://127.0.0.1:8080/v1`（逗号后带空格）落盘后 `config.get` 载荷读到 `modelCount=1`、`hasBaseURL=true`；两拒后对照行 `local-gw-four` 仍可建。
- **变异红名逐字**：
  - M2（删掉预设名拒绝）⇒ `--- FAIL: Test320CreateCustomProviderRowLandsAndRefusesByName/a_preset_name_belongs_to_provider_row_create: settings_320_test.go:152: the create was accepted (reported [llm.providers.stepfun])`（后四枚子测级联同红，逐字在件内）。
  - M2b（腿把切点从 `,` 改成 `;`）⇒ `--- FAIL: Test320CustomProviderRowThroughTheRealLeg: panel_config_320_test.go:109: config.set was refused: ... 建 [llm.providers.local-gw] 的值必须是「协议,base_url」两半 ...`。
- **门禁 rc**：build 0 / vet 0 / gofumpt 零输出 / d22scan clean / `-run Test320CreateCustomProviderRow` ok。

## 3. AC#3 `model_row_create`（sha `84fc3f20`）

- **改了什么／落点**：`config_write_313.go`（常量＋`fields313`＋NeedsProvider＋NeedsModel，**不进** takesValue）· `settings_320.go`（`CreateModelRow(provider, model)`，键 `llm.providers.<p>.models.<id>`；ModelSpec 由**既有 `applyDefaults`** 播种：`enabled=true`／`billing=pay-per-token`；拒未知 provider 与已存在条目）· `cmd/wisp/panel_config_store.go`（case）。
- **两向三态读数**（`03-ac3-two-way.txt`）：接受新条目落节并回读 `enabled=true`／`billing=pay-per-token`；随后 **既有三枚** `context_window`／`price.in`／`price.out` 在同一枚新条目上落盘（这是本步存在的理由：313 的叶子字段只能改"已有条目"）；三条拒绝具名且文件字节不变；`config.get` 名册回读两条按序；对照第二枚条目落盘且不扰第一枚的叶子。
- **变异红名逐字**：M3（不播默认标签、用零值 ModelSpec）⇒ `--- FAIL: Test320CreateModelRowSeedsFromDefaultsAndRefusesTwice: settings_320_test.go:197: the created entry is not enabled: an entry without enabled=true is treated as off and is refused at start-up`（＋`:200` billing 空）。
- **门禁 rc**：build 0 / vet 0 / gofumpt 零输出 / d22scan clean / `-run Test320CreateModelRow` ok。
- ⚠ **实读抖动一条**：`cmd/wisp -run Test320ModelRowCreateThroughTheRealLeg` 首抓一趟 rc=1 且**无 `--- FAIL` 行**（进程未具名即死，派单已点名的 `./cmd/wisp` 抖动）；同参复跑两趟 rc=0 逐字 PASS，件内两趟都记着（未粉饰）。

## 4. 门禁汇总与基线对比（探针 `00-…`／`05-…`）

| 门 | 读数 |
|---|---|
| `go build ./...` | rc=0 |
| `go vet ./...` | rc=0 |
| gofumpt（本腿六件） | 零输出（`gofumpt -l`） |
| `sh scripts/d22scan.sh` | `clean`；内部 238／cmd 40／emoji 覆盖面 560+132 全走 |
| 长度门 | `VERDICT GREEN`（60=60，最长 185） |
| `./internal/config` 全包 | ok（rc=0，零 FAIL） |
| `./internal/panel` 全包 | rc=1，5 枚失败：`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`／`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`／`TestComposerContractTypesMatchFrontend`／`TestC21DesignTokensFourWayAgree` |
| 同上，但**把本腿六件 overlay 回 `33e429c1`** | **逐枚同名 5 枚**（`gate_overlay_baseline.txt`）⇒ 非本腿引入；归因：`frontend/**` 与 `design/**` 的机主外部面（工作树里 `design/assets/tokens.css` 等处于删除态，`git ls-files` 可证其被跟踪） |
| `./cmd/wisp` 带窗 255 段 | 带腿与 overlay 回退**都** rc=1 且都无 `--- FAIL` 行（整进程死，末行同为 `MARK299 wndproc`）⇒ 与派单所述 314 死区同形 |
| `./cmd/wisp` 无窗段（313/320＋255 名册扫描＋248 凭据） | rc=0（`ok`） |

## 5. 面核（`git diff --name-only` 逐枚）

- 四笔共动 **六件**：`internal/panel/config_write_313.go`、`internal/panel/config_313_test.go`、`internal/config/settings_320.go`、`internal/config/settings_320_test.go`、`cmd/wisp/panel_config_store.go`、`cmd/wisp/panel_config_320_test.go`（`05-…` 逐 sha 列表在件内）。
- 冻结/硬钉面**零改动**（`33e429c1..HEAD`）：`l2_grant_boundary_test.go`、`tokens_fourway_test.go`、`internal/perm/ticket90_persist_test.go`、`config_route_248_test.go`、`cmd/wisp/panel_config_313_test.go`、`config_handlers.go`、`bridge.go`、`composer_dispatch.go`、`inbound_roster_253_test.go`、`frontend/`、`design/`、`docs/PLAN.md`、`docs/specs/`、`tools/d22scan/allowlist.txt`、`.github/`（输出为空；同名区间里 `.github/workflows/ci.yml` 的改动出自 `757f08fc`＝编排者 111 收口，非本腿）。
- 新增字段名三枚**不带点、不带审批词**；`config_313_test.go` 的 live 命名禁令测试自动覆盖（逐枚 0 命中）。

## 6. 具名没做（本腿）

- **AC#4**"拉取模型列表"端口：**没做**（`320-r2` 的活；前端 `frontend/src/lib/panel.ts` 字面量先落才可 self-green）。
- **AC#5** 预设默认模型目录：**不落**（`A972` 第 5 裁逐字）。
- `config.models.list` 或**任何新 socket 名**：**没加**（不加 `ConfigStore` 方法＝248 面零改的前提）；`inbound_roster_253_test.go` 计数未动。
- `docs/PLAN.md`／`docs/specs/**`／阈值／golden／`allowlist.txt`／`.github/**`／`frontend/**`／`design/**`：**零碰**。
- 票面 `Status:`／框／票名：**未动**（只追加 1 条 bullet）。
- 台账 `pending-and-issues.md`：未动（判语归非实现者）。
