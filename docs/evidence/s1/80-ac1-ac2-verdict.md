# 票 80 — AC#1 / AC#2 裁决：`blacklist_overrides` 是死线，但**契约没有写这条能力** ⇒ 不接线、上报编排者

- 票面：`.scratch/wisp/issues/80-blacklist-overrides-never-wired-to-gate.md`
- 起点：分支 `dev` @ `17efc2c`；本轮 commit：`87de494`（认领）、`0f8db02`（AC#1 清单）、本文与 L3（AC#2 裁决）
- **交付性质：判定 + 上报，零 Go 文件改动。** AC#3/AC#4 被 AC#2 的裁决挡住，不是"没做完"。

## 1. 结论摘要

| 门 | 裁决 |
|---|---|
| AC#1 死线证明 | **成立，且比票面更死**：`bOverrides` 全仓无写入者；更进一步，`risk.Gate` 本身**生产零调用点** |
| AC#2 契约判定 | **契约没有把 `blacklist_overrides` 写成"决策时生效的用户能力"** ⇒ 走票面 ⚠ 分支：**停下上报编排者，由编排者判是否 D22** |
| AC#3 / AC#4 | **未做（按裁决不该做）**：本票不产出实现码 |
| AC#5 门禁 | **零改动面 ⇒ 无门禁可闭**；已交 `internal/risk` + `internal/config` 的现状基线（含一条**先存在的脆弱性能用例**） |

## 2. AC#1：读/写穷尽清单（完整版在票面 Progress log L2）

检索式（全仓，含 `cmd/`、`internal/tools/`、`internal/agent/`、测试、`docs/`）：
`BlacklistOverrides|blacklist_overrides|blacklistOverrides`（-i）、`bOverrides`、
`risk\.Classify|risk\.Gate|risk\.Class`、`SensitiveClassifier|TierB`、`\.Risk\.`、`ConfirmLocked`、`Gate\(`。

- `config.RiskSection.BlacklistOverrides`：**写**只有 `internal/config/schema.go:457`（TOML 解码）与
  `internal/config/boundary_test.go:147`（测试）；**读**只有 `internal/config/manager.go:356`（方向审计）。
  `cmd/`、`internal/tools/`、`internal/risk/`、`internal/agent/` 零命中 ⇒ 无消费者。
- `risk.Gate` 的 `bOverrides`：`internal/risk/blacklist.go:76/85/115/116/120/125` 全是**读**；
  写入只存在于测试就地字面量（`pathresolver_anchor_spelling_windows_test.go:179,184,189,194,221,237,273`；
  `pathresolver_junction_windows_test.go:263,291`；`pathshape_portable_test.go:241`）。
- **生产判定链不含 Gate**：`cmd/wisp/run.go:262 tools.New(...)` → `internal/tools/bridge.go:159 NewSensitiveClassifier()`
  → `internal/tools/paths.go:146 risk.Classify(canonical)`；`Classify` 签名（`blacklist.go:69`）**没有 override 形参**。
- **不存在"人点过 L2 之后记住该文件"的容器**：`internal/agent/approval/approval.go:273-320` 的 `grantStore`
  是**一次性 nonce**（F2 第 3 层），与豁免记录无关。
- 该键唯一残存的活效果（变更检测）也是哑的：`manager.go:239` 走的 `m.ConfirmLocked` 钩子
  在生产里**从未被赋值**（全仓非测试命中只有它的声明与这一处调用）。

⇒ 票面"后果 1：配置项说谎"成立；"通过另一个字段名其实接上了"**未找到** ⇒ 本票不转登记为"报告者计数错误"。

## 3. AC#2：契约引用（每条都自己打开核对；票面既有引用 `schema.go:457` / `manager.go:356` 核对为真）

| 编号 | 位置 | 原文承担什么 |
|---|---|---|
| C1 | `docs/PLAN.md:2732` | 🔒 `[risk]` 配置表行列出 `blacklist_overrides[]`，整行标"变更即需 L2 级重新确认" |
| C2 | `docs/specs/SPEC-03-config-secrets-envs.md:34` | 同上：`blacklist_overrides[]=[]`（B 档豁免）｜🔒 放宽需 L2 重新确认 |
| D1 | `docs/specs/SPEC-06-security-gatekeeping.md:67` | "B 档 · 默认拒绝 + 单文件豁免（**豁免 = 一次 L2 强确认 + 写日志**）" |
| D2 | `docs/specs/SPEC-06-security-gatekeeping.md:36` | R3："B 档 → L2（可单文件豁免）" |
| D3 | `docs/PLAN.md:2396-2399` | "**但豁免必须按文件、必须由人点、必须留痕。**"（F1 CRITICAL） |
| D4 | `docs/specs/SPEC-06-security-gatekeeping.md:19` | L2 的「允许」**只接受原生侧来源**（悬浮球点击/原生确认卡按钮/全局快捷键） |
| D5 | `internal/risk/assessor.go:144-146`（冻结，仅读） | "tier B is L2 with per-file exemption（**exemption flow lands with the approval queue, ticket 21**）" |
| D6 | `internal/risk/blacklist.go:13-14, 74-75` | `bOverrides` 的既定生产者 = "individually L2-confirmed B-tier files" |
| — | `docs/specs/SPEC-10-testing-acceptance.md` | **零命中**：无任何验收条目要求该键在决策时生效 |

**判词。** `grep -rn "blacklist_overrides" docs/` 的全部命中就是 C1/C2 两张配置键表；拥有这项能力的规格正文
（SPEC-06 §4.1 / PLAN F1）把豁免**用等号定义成一个运行时事件**（D1），并用"必须由人点"（D3）排除静态预授权，
D5/D6 则把这条水流归给审批队列。没有任何文字说明该键的决策时语义（谁读、key 形状、生效时机、与 L2 确认的先后）。
⇒ 接上它不是"把契约落成实现"，而是**新造一条契约未写的放行侧预授权通道**；票面 AC#3 的 (a) 形状
（配置里写一行 ⇒ 该文件从 B 变放行，全程无人点）恰是 D3 所禁止的那种。

## 4. 交回编排者的四选一（本票不代拍）

- **(A)** 照 AC#3 原样接静态预授权 ⇒ 必须先改 `SPEC-06 §4.1`(D1) 与 `PLAN.md:2399`(D3) 的"必须由人点" ⇒ **需 D22 人工批准**（冻结文档）。
- **(B)** 接成契约真正写的那件事：预登记 + **仍要人点一次**，确认后把 C26 canonical 记进 `bOverrides`
  ⇒ 落点在审批队列，生产链要从 `risk.Classify` 换到 `risk.Gate`，并要动 `rules_gateway.go:84` 的 `case TierB:` /
  `assessor.go:158` 的 seam（**均为冻结面**）⇒ **需 D22**；这本就是 D5 里 ticket 21 的欠账，不该塞进本票。
- **(C)** 不接线，但让该键**响亮地失败**（写非空即报错 / doctor 亮红灯），照既有先例 `verify_signature`
  （`SPEC-03:42` "不可关：写 false 报错而非生效"）⇒ **收紧侧**，不放大放行面，但改 SPEC-03 键语义仍需点头。
- **(D)** 维持现状 + 登记（当前交付）。

## 5. 基线门禁真实数字（本票零 Go 改动；跑前 `git status --porcelain` 只剩票面一处 M）

| 命令 | rc | 真实输出/计数 |
|---|---|---|
| `gofmt -l internal/risk internal/config` | 0 | 空 |
| `go vet ./internal/risk/ ./internal/config/` | 0 | 无输出 |
| `go test -count=2 -v ./internal/config/ ./internal/risk/` | **1** | `ok internal/config 0.975s`；`FAIL internal/risk 17.114s`。`=== RUN` 460；`--- PASS` 267 + `    --- PASS` 190；`--- SKIP` 2（`TestSyncRegistryProbeLive` × count=2）；`--- FAIL` 1 = `TestResolvePerCallBudget`（`pathresolver_budget_norace_test.go:37`，1.640779ms > 1ms 预算） |
| `go test -count=2 -run TestResolvePerCallBudget -v ./internal/risk/` | 0 | 2 RUN / 2 PASS，0.692 与 0.637 ms/op ⇒ 上一条 FAIL 是**负载敏感抖动** |
| `go test -count=2 ./internal/risk/`（单独重跑） | 0 | `ok internal/risk 6.806s` |
| `GOOS=linux go vet ./internal/risk/ ./internal/config/` | 0 | 无输出 |

**假绿点名**：`TestOverrideKeyKeepsPosixBackslashesDistinct`（`internal/risk/pathshape_portable_test.go:192`）
在 `filepath.Separator == '\\'` 时提前 `return` ⇒ **本平台它是空跑**，不可当作"跨文件 bleed 已被守住"的证据；
它的 PASS 也不构成接线判据。`-count=2` 的 N 倍核对以 SKIP 恰好 2 行 + 两包独立重跑一致为准。

## 6. 在外的先存在事项（登记，不在本票改动面）

1. `TestResolvePerCallBudget` 是**墙钟 1ms** 的性能契约（票 18 AC#6）：与整包并发跑会抖到 1.6ms 判红 ⇒
   "两包并跑"的门禁形态本身会假失败。建议单开一票（`testing.Short()` 门控，或改成同机基线的倍数）。**不归票 70**（它只做格式化）。
2. 🔒 锁段配置的"变更即需 L2 重新确认"（`ConfirmLocked` / D36 规则 1）在生产**整体未接线**；
   同族还有 `[risk]` 的 `shell_enabled` / `allow_shell_string` / `shell_allowlist` 在 `internal/config` 之外**零消费者**
   （`grep "\.Risk\."` 全仓只命中 `L1WindowSec`、`ConfirmTimeoutSec`）⇒ 这是一族问题，不止本票这一键。
