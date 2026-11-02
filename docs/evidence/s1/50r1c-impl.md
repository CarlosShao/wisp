# 50-r1c — 落地腿实现证据（票 50 · 只做 AC1 + AC4）

- 腿：`50-r1c`（写腿·收窄）｜工作目录 `internal/plugin/**` + 本文档｜⛔ push（只 commit）
- 基线 HEAD：`a18ff45`（编排者 A1428 派单存档 `.scratch/wisp/probes/orch/20261102/txt/50r1c-dispatch.txt`）
- 交付 commit：`73ab5da6`（实现三件）＋ `635a2707`（定向测试两件）
- 票面：`.scratch/wisp/issues/50-tier1-manifest-plugins.md`

## 0. 起手三核（现跑）

| 核 | 读数 |
|---|---|
| HEAD 现盘最新 | `a18ff4546df66755ab524e4fad62f11da310210f` Mon Nov 2 21:22:18 +0800 |
| 索引干净 | `git diff --cached --stat` 空；`git status --porcelain -- internal/plugin` 空（本腿写面无他脏） |
| C14 字段表已定义？ | **已定义 ⇒ 不停手**。权威三处：`docs/PLAN.md:1370`（C14 契约行，含 D46 增项）、`docs/specs/SPEC-07-tools-and-plugins.md:77-111`（§4 完整 TOML 字段表）、`SPEC-07:128-139`（§6 生命周期含加载时哈希校验）；capability 词表＝`PLAN.md` C3（11 项）/`SPEC-07 §2`；工具名形状＝C1 `namespace.action`；风险档位＝C2 `L0/L1/L2` |

未发明任何字段：字段集完全取自上面四处。票 51/52 拥有的 C14 字段（`entry`／`kind/exe/exe_hash/argv_template/extract/timeout_ms/env_allowlist`）**本腿不实现、也不当作未知字段**——见 §3。

## 1. AC1 — C14 manifest 解析 + schema 字段校验（未知 key 拒绝）

落件：`internal/plugin/manifest.go`（结构体 + 闭合字段表 + `ParseManifest`）、`internal/plugin/validator.go`（逐字段规则 + `ValidationErrors`）。

三道闸，顺序固定（先拒形、再解型、后语义）：

1. **闭合 key 审计 `auditKeys`**：TOML → 泛型树 → 每个 key 必须落在 C14 Tier-1 表内。
   - 顶层允许：`schema_version id host_api name version tools`
   - `[[tools]]` 允许：`name description risk capabilities net_allowlist parameters`
   - 未知 key ⇒ 拒绝，报文点名 `manifest.<key>` / `tools[i].<key>` + "unknown key … never ignored"；**绝不静默忽略**。
   - `parameters` 是不透明子表（JSON Schema 自己的语法），审计**不 descend**（`$schema`/`additionalProperties`/`maxLength` 等不会被当未知字段）。
2. **TOML 类型解码**：形状不符（`tools` 不是数组表、数组元素不是表）→ `(document)` 级拒绝。
3. **逐字段语义 `Validate`**：`schema_version` 必须＝1；`id` 必须匹配 `[A-Za-z0-9._-]+` 且非 `.`/`..`（它同时是 `plugins\<id>\` 目录名，SPEC-07 §6）；`host_api`/`name`/`version` 必填且无控制字符；`tools[i].name` 必须 `namespace.action`（C1）；`description` 必填；`risk ∈ {L0,L1,L2}`（C2；声明只是 R1 下界，抬高由 C19，本腿不判）；`capabilities` 逐项必须 ∈ C3 11 项且不得重复/空串；**声明 `net` ⇒ `net_allowlist` 必填非空**（SPEC-07 §4 明文方向），条目必须是裸主机名（无 scheme/端口/路径/通配符）；`parameters` 必须是非空表。
   所有发现一次全报（`ValidationErrors` 列表），改一版文件就能看到全部问题。

## 2. AC4 — tamper 校验（翻转即拒 + 告警）

落件：`internal/plugin/tamper.go`。契约出处：SPEC-07 §6「manifest 哈希不符 → 拒绝加载并告警（防文件被就地替换）」＋票 50 AC4。

- `ManifestHash(raw)` ＝ `sha256:<小写hex>`（与 `plugin_state.hash`／C14 `exe_hash` 同形状）；用 FIPS "abc" 公开向量钉死摘要本身，不只自洽。
- `VerifyBytes(raw, pinned)` 纯函数核；哨兵错误 `ErrTampered` / `ErrNoPinnedHash` / `ErrMalformedHash`（`errors.Is` 可分派，⛔ 字符串匹配）。
- `LoadInstalled(path, LoadOptions{PinnedHash, PluginID, Alerter})`：**先核哈希，后解析**。
  - 哈希不符 ⇒ ①同步 `Alerter.AlertTamper`（含 kind/plugin_id/path/expected/actual）②返回 `ErrTampered` ③返回零值 Manifest。
  - **没有 pinned hash ⇒ 也拒**（不可验证＝不可信，⛔ "load anyway"）；malformed pin ⇒ 拒＋告警。
  - 默认 `LogAlerter` 打稳定事件名 `plugin_tamper`（结构化 Error 日志，⛔ 静默）；`CaptureAlerter` 供测试与后续审计/面板接线。
  - 哈希对但 schema 不合法 ⇒ 走 C14 失败路径，**不算 tamper、不告警**（两条判据不混淆，有定向用例钉住）。
  - 路径按原样使用，本腿不做 `Clean/Abs`（FS 决策属 `risk.PathResolver`，AGENTS §1.2）。

## 3. 读数

| 项 | 读数 |
|---|---|
| 定向用例（AC1+AC4） | `PASS 68 / FAIL 0`（清单 `D:\tmp\wisp\20261102\50r1c\ac-readings.txt`） |
| AC1 逐字段矩阵 `TestFieldMatrix` | **34 例＝9 正 / 25 反**（schema_version 3、id 5、host_api 2、name/version 2、risk 6、capabilities 5、net_allowlist 6、tool name/description 4、parameters 1） |
| AC1 未知 key 反向判据 | `TestUnknownKeysRejected`（顶层 `auto_approve`／工具级 `silent_shell`／错拼 `net_allow_list`）全红→绿 ✓ |
| AC1 归属命名判据 | `TestDeferredKeysNameOwner`：`entry` 报文点名 ticket 51、`exe` 点名 ticket 52，且**不得**谎报 unknown ✓ |
| AC1 不透明子表 | `TestParametersKeysAreOpaque`：JSON Schema 内层 key 未被当 manifest key 拒 ✓ |
| AC4 反向判据 | `TestLoadInstalled_ManifestBitFlip` 4 变体（1-bit 翻转 / `L0→L2` 合法但不同字节 / 追加 `[[tools]]` 外发工具 / 截断文件）⇒ 全部 `ErrTampered` ＋ 恰好 1 条告警（expected/actual/path 齐全）＋ 返回零值 ✓ |
| AC4 其余 | `TestVerifyBytes`（含大小写不敏感、no-pin、4 种 malformed pin）、`TestManifestHash_Format`（FIPS 向量）、`TestLoadInstalled_HappyPath`（干净装载零告警）、`_HashOkSchemaBad`（不混淆）、`_MissingFile`（普通读错误，不告警）、`_DefaultAlerterDoesNotPanic` ✓ |
| **变异核（反向判据实跑）** | MUT-A 把 key 审计短路（`if true \|\| allowed[k]`）⇒ `TestUnknownKeysRejected`＋`TestDeferredKeysNameOwner` **变红**；MUT-B 把哈希比较短路（`if false && …`）⇒ `TestLoadInstalled_ManifestBitFlip` **变红**；两枚还原后 `git diff --stat -- internal/plugin` 空、`git status --porcelain -- internal/plugin` 空（字节级还原），整包 `ok` ✓ |
| 包级 / 仓级 | `go test ./internal/plugin/ -count=1` → `ok`（含前腿 disposal 用例不回归）；`go vet ./internal/plugin/` 干净；`go build ./...` 干净 |
| 格式 | `/d/work/base/gopath/bin/gofumpt.exe -l -w internal/plugin` 已跑（5 文件），复跑 `-l` 空 |
| 写面纪律 | 只新增 `internal/plugin/{manifest,validator,tamper}.go` + `{manifest,tamper}_test.go`；⛔ 未碰 `internal/models`、`cmd/wisp`、`thresholds.go`/golden；`git add` 逐枚字面 pathspec；⛔ `add -A`/amend/reset；⛔ push |

## 4. 后续登记（本腿具名交出，⛔ 未假完成）

| 项 | 归属 | 缺口现状 / 接法 |
|---|---|---|
| **AC2** host_api `^1.0` vs host 1.x 装载、vs 2.x 明确拒绝（C16 闸门） | 票 50 后续腿 | 本腿只校 `host_api` 存在＋无控制字符（`validator.go` 已写注释指明属 AC2）；semver range 语法解析与主版本闸门未实现 |
| **AC3** 安装/卸载 e2e（来源显示、能力清单＋来源显式授权、落盘 `plugins\<id>\`、`plugin_state` 写入、卸载 DisposalScope 零残留） | 票 50 后续腿 | 本腿只提供 `LoadInstalled`/`VerifyAndParse` 这道闸；`internal/memory.PluginState.Hash/ExeHash` 列已存在，安装腿把 pinned hash 传进来即可复用 |
| **AC6**（票面 AC5/AC6 列）插件工具进 `list_tools`/discovery＋risk pill、经 bridge 执行记录 plugin id、capability 越界拒调用、net_allowlist 名单外→L2（R5） | 票 50 后续腿 + `tools`/bridge 片 | 本腿不改 bridge，也不在 plugin 侧判权限（SPEC-07 §2「`plugin` 不得自行判权限」）；`ToolDef` 已带 capabilities/net_allowlist 原样数据供其消费 |
| 票 51：`entry`（Tier-2 goja） | 票 51 | 现在报 "not implemented by the Tier-1 base loader; owner: ticket 51" |
| 票 52：D46 `kind/exe/exe_hash/argv_template/extract/timeout_ms/env_allowlist` ＋ 执行前 `exe_hash` 校验 | 票 52 | 现在点名 owner: ticket 52 拒装；哈希工具 `ManifestHash`/`VerifyBytes` 可直接复用 |
| 票 57：registry 元数据 / 发布校验 | 票 57 | — |

### 待定案（本腿未自行扩契约，留给裁决者/编排者拍）

1. **空 `[[tools]]`** 是否合法：契约无明文 ⇒ 本腿**不拒**（拒＝发明规矩）。注册表侧若要求 ≥1 工具，请在 AC3/票 57 定。
2. **`net_allowlist` 存在但未声明 `net`**：契约只写了一个方向（net ⇒ 必带白名单），反向未写 ⇒ 本腿不拒，已在 `validateNetAllowlist` 注释标记。
3. **`version` 严格 semver 解析**、**`host_api` range 语法**：都属 C16/AC2 那格，本腿只做存在性。
4. 工具名**全局唯一**（跨插件）与域名**匹配语义**（精确 vs 后缀）：分别属 registry 与 net 执行腿（AC5），本腿只做字符集合规。

## 5. 失败预演（本腿射程内的最小三条）

- 攻击者就地替换 manifest（加一个 `capabilities=["shell"]` 的工具）⇒ 哈希不符 ⇒ 拒载＋告警；即使替换内容本身 schema 合法也一样拒（`swap_risk_level`／`append_extra_tool` 用例钉住）。
- 安装记录丢了 hash（或被人写成非 `sha256:` 串）⇒ `ErrNoPinnedHash`/`ErrMalformedHash` ⇒ 拒载＋告警，⛔ 降级为"先装着"。
- 插件作者塞未知字段（`auto_approve = true`）⇒ 拒装并点名 `manifest.auto_approve`；错拼契约字段（`net_allow_list`）⇒ 同样拒，⛔ 静默丢掉白名单。
