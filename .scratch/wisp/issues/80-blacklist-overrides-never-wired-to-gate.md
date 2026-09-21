# 80 — `blacklist_overrides` 配置被解析、被方向审计，但**从来没有接到 `risk.Gate` 上**

**Status:** AC#1 已裁决（死线确认，且比票面更死：`risk.Gate` 生产零调用点）；正在跑 AC#2 契约判定
**Type:** 契约与实现脱节（一个面向用户的开关是死线）
**Blocks:** nothing · **Blocked by:** nothing（`internal/risk` 此刻无人写：票 75 已交、票 72 只碰过 `pathresolver.go`/`blacklist.go`）
**Packages:** `internal/config/`、`risk.Gate` 的构造点、以及新建的接线测试。**禁改冻结面**：
`internal/risk/assessor.go`、`internal/risk/pathresolver*.go`、`rules_gateway.go`、`docs/PLAN.md`、`docs/specs/*.md`。
**发现的**：票 75 的接续代理（2026-09-21 收尾报告第 6 节第 1 条）。

## 缺陷

票 75 的实测链条：
- `internal/config/schema.go:457` 有 `BlacklistOverrides`，会被解析；
- `internal/config/manager.go:356` 甚至对它做了**方向审计**（说明写它的人预期它生效）；
- 但 `risk.Gate` 的 `bOverrides` map **在生产路径上没有任何调用者往里灌数据** ——
  那个"单个文件的 B 档豁免"能力**只活在测试里**。

后果有两种方向都不好的读法：
1. 用户/文档以为可以豁免某条 B 档，实际**静默无效**（配置项说谎）；
2. 反过来，若哪天接上去，等于**没人审过的放行通道**被打开。

## 判据前先做的一件事（AC#1，只读）

**"grep 零命中"不等于"没有保护"，也不等于"没有生效"**。所以 AC#1 必须先穷尽列举
`BlacklistOverrides` / `bOverrides` / 任何等价的 override 容器在**全仓**的每一次读与写
（含 `cmd/`、`internal/tools/`、桥接层、以及测试里的构造），并给出**逐条位置 + 该行是读还是写**的清单，
然后才允许下"它是死线"的结论。若清单证明**其实有生效路径**（例如通过另一个字段名接上），
那本票就地关闭并登记为"报告者计数错误"，**不许为了有活干而造修复**。

## AC（1:1，一框一行裁决表，验收在 `docs/evidence/s1/80-*.md`）

- [x] **AC#1** 上面那张"每次读/写"清单进本票 Progress log，并明确结论：死线 / 已生效（后者 ⇒ 本票转登记，不写码）。
- [ ] **AC#2** **契约判定**：在 `docs/PLAN.md` 与 `docs/specs/SPEC-06*.md`（C19/C26 相关条目）里查
  `blacklist_overrides` 是否被写成**用户可用能力**。
  - 若**是**（契约要求它生效）⇒ 接线是**修 bug**，本票继续 AC#3。
  - 若**契约没有这条**、或它属于"放行侧放宽"⇒ **停下来上报编排者**，由编排者判是否 D22。
    ⚠ 不许用"先接上再说"绕过：放行侧的开关是安全边界，不是配置便利。
- [ ] **AC#3** 接线落地 + 一条**真跑在 `risk.Gate` 上**的用例：同一份配置里写 override，
  证明 (a) 被点名的那一个文件从 B 变成放行，(b) **同目录下其它文件仍然 B**，
  (c) override 的形状必须是 C26 交回的 canonical 形式（票 72 的 R17 裁定：
  **放行侧只认解析后的形式**，`bOverrides` 的 key 不接受拼写变体）。
- [ ] **AC#4** 变异双向：(i) 把接线断开 ⇒ AC#3 必须红；(ii) 把 (b) 的"其它文件仍 B"断言指向 override
  ⇒ 必须红。两次都要在**同一条 `&&` 链里先 grep 证明改动真落地**再跑，还原后 `git diff --quiet` 证干净。
- [ ] **AC#5** 门禁（**只跑自己碰到的包**，共树不跑整仓）：`gofmt -l <pkgs>` 空、
  `go vet <pkgs>` rc=0、`go test -count=2 <pkgs>` rc=0，且**逐跑点名** `--- SKIP`/`--- FAIL` 的行数与测试名；
  `GOOS=linux go vet <pkgs>` rc=0（Linux 编译哑弹是 A44/A49 的老坑）。

## Rules

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc，A31）；**禁** `git add -A`/`.`；
每次 commit 前 `git diff --cached --name-only` 核对清单；
**禁** `--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34，共树）；**不 push**（编排者的活）；
不在仓目录里建 worktree（A38④），要纯净树用 `git archive HEAD | tar -x -C /tmp/...`；
不碰别人的票面（尤其票 75/72/79 的 `.md`）；15 次工具调用内交第一枚 checkpoint commit，
之后每次 commit 同步 Status + 勾框 + `next=`。

## Progress log（append-only，最新在下）

### L1 — 认领 + checkpoint（2026-09-21）

已把 Status 从 `open` 改为 `claimed`。本轮次序：AC#1 只读穷尽清单 ⇒ AC#2 契约判定 ⇒ 只有两者都过才允许写实现码。
`next=` 跑 `BlacklistOverrides` / `bOverrides` / 等价 override 容器在全仓的每次读与写清单。

### L2 — AC#1 只读穷尽清单（2026-09-21，全部位置由 grep + 逐条打开核对）

检索式（全仓，含 `cmd/`、`internal/tools/`、`internal/agent/`、测试、`docs/`）：
`BlacklistOverrides|blacklist_overrides|blacklistOverrides`（大小写不敏感）、`bOverrides`、
`risk\.Classify|risk\.Gate|risk\.Class`、`SensitiveClassifier|TierB`、`\.Risk\.`、`ConfirmLocked`。

**容器 1 —— `config.RiskSection.BlacklistOverrides []string`**

| 位置 | 读/写 | 说明 |
|---|---|---|
| `internal/config/schema.go:456` | —（注释） | "B-tier blacklist exemptions (per entry)" |
| `internal/config/schema.go:457` | **写**（唯一真实写入者） | 字段声明；由 TOML 解码器填充。无 `default:` tag（默认零值 nil） |
| `internal/config/manager.go:356` | **读**（`old` 与 `new` 各一次） | `setDirection(...)` → 只产出 loosen/tighten 键路径；票面引用正确 |
| `internal/config/boundary_test.go:147` | **写**（测试） | `c.Risk.BlacklistOverrides = []string{"B-1"}`，只测方向审计 |
| 全仓其它位置 | 无 | `cmd/`、`internal/tools/`、`internal/risk/`、`internal/agent/` **零命中**（含任何别的字段名转发） |

**容器 2 —— `risk.Gate` 的 `bOverrides map[string]bool`**

| 位置 | 读/写 | 说明 |
|---|---|---|
| `internal/risk/blacklist.go:76` | 形参声明 | `func Gate(canonical string, bOverrides map[string]bool)` |
| `internal/risk/blacklist.go:85` | **读** | 传给 `overrideApplies` |
| `internal/risk/blacklist.go:115,116,120,125` | **读** | `len(...)`、`bOverrides[cand.raw]`、`bOverrides[cand.real]`；**没有任何写入** |
| 生产调用点 | **0 个** | `grep -rn "Gate(" --include=*.go \| grep -v _test.go` 只剩定义本身与同名无关符号（`audio.NewHalfDuplexGate`、`provenance.writeGate`） |
| `internal/risk/pathresolver_anchor_spelling_windows_test.go:63` | 传 `nil` | — |
| 同文件 `:179,184,189,194,221,237,273` | **写**（调用点就地字面量 map） | 全部 `map[string]bool{strings.ToLower(<path>): true}`，值来自测试自己构造的路径 |
| `internal/risk/pathresolver_junction_windows_test.go:82,283` | 传 `nil` | — |
| 同文件 `:263,291` | **写**（字面量 map） | 同上 |
| `internal/risk/pathshape_portable_test.go:241`（`ovr :=` 构造）→ `:242,245` 读用 | **写**（测试局部变量） | key = `normPath(res.Canonical)`；注释自称"the shape config's risk.blacklist_overrides will feed"——**这是意图陈述，不是实现**。注：同函数在 `:192` 对 `filepath.Separator == '\\'` 提前 `return`，**该用例在 Windows 上是空跑** |

**容器 3 —— 等价物排查（生产实际走的判定链）**

| 位置 | 读/写 | 说明 |
|---|---|---|
| `internal/tools/paths.go:136-152` | — | 生产侧敏感路径判定 = `sensitiveClassifier.Classify()` → `risk.Classify(canonical)`。**`Classify` 根本没有 override 形参**（`blacklist.go:69`） |
| `internal/risk/assessor.go:158-160` | — | 冻结的 seam：`type SensitiveClassifier interface { Classify(canonical string) Tier }`，**一参、无上下文** |
| `internal/risk/rules_gateway.go:84-90` | — | 冻结：`case TierB:` → 一律 L2，reason 写"可单文件豁免"，**无处查豁免** |
| `internal/tools/bridge.go:159` | — | `classifier: NewSensitiveClassifier()`，**构造时不吃任何配置** |
| `cmd/wisp/run.go:223-227, 264-268` | — | 组合根只喂 `cfg.FS.AllowedDirs`/`ReparsePointExceptions`/`cfg.FS.DeleteEnabled`/`cfg.Risk.L1WindowSec`/`ConfirmTimeoutSec`；**从不引用 `cfg.Risk.BlacklistOverrides`** |
| `internal/agent/approval/approval.go:273-318` | — | `grantStore` 是**一次性 nonce**（F2 第 3 层），不是持久豁免记录 ⇒ 全仓不存在"人点过 L2 之后把该文件记进某张 map"的容器 |

**AC#1 结论：死线确认成立，且比票面描述更死。**

1. `bOverrides` 在全仓**没有任何写入者**（除测试就地字面量）；
2. 更强的事实：`risk.Gate` 本身**在生产路径上零调用点** —— 生产走的是 `risk.Classify`，那条函数**签名里就没有 override**；
   ⇒ 因此"接线"不是"把 config 灌进 Gate 的 map"这一步，而是要先在**判定链**上开一个能携带豁免集合的位置；
3. `BlacklistOverrides` 唯一仍在起作用的后果是**变更检测**：往该键加元素会把整节 `[risk]` 的改动判成 loosen
   （`manager.go:239`）。但该分支调的 `m.ConfirmLocked` **在生产里也从没被赋值**（全仓唯一非测试命中就是它的声明与调用点）
   ⇒ 于是当前真实行为是：写了 `blacklist_overrides` ⇒ 触发不了 L2 重确认、也豁免不了任何东西 ⇒ **配置项说谎**（票面后果 1 成立）。
4. 未发现"通过另一个字段名接上"的生效路径 ⇒ 本票**不转登记**，继续 AC#2。

附带登记（**不在本票改动面**，交回编排者）：`[risk]` 的 `shell_enabled` / `allow_shell_string` / `shell_allowlist`
在 `internal/config` 之外同样**零消费者**（`grep "\.Risk\."` 全仓只命中 `L1WindowSec` 与 `ConfirmTimeoutSec`），
且 `ConfirmLocked`（D36 规则 1 的 L2 重确认钩子）未接线 —— 这是**一整族**"锁段配置未落地"，不是本票单点。

`next=` AC#2 契约判定（`docs/PLAN.md` 的 F1 段 + `docs/specs/SPEC-06*` §4.1 + `SPEC-03` 的 🔒 `[risk]` 行）。
