# 80 — `blacklist_overrides` 配置被解析、被方向审计，但**从来没有接到 `risk.Gate` 上**

**Status:** **done**（编排者验收 2026-09-21 13:0x）—— 本票的交付物是**裁决**，不是码：**零 Go 改动**，
且这一格判得比建票时准。AC#1 的读/写穷尽清单把定性**推翻了编排者自己的 A51⑤**（真相是
`risk.Gate` **生产零调用点**，不是"一根线忘了接"）；AC#2 用契约正文证明**契约从未承诺**这个键生效
⇒ 按票面第二分支停下，没有为了有活干而造修复。**编排者裁决 = 选项 (C)**（见 A53②），
后续工作**全部转移到票 83**；AC#3/AC#4/AC#5 的接手方逐条写在框后。
**Type:** 契约与实现脱节（一个面向用户的开关是死线）
**Blocks:** nothing · **Blocked by:** ~~编排者对 AC#2 第二分支的裁决~~ **已给（A53②：选 (C)，需 D22 的只有键表那半句 ⇒ Q-27）**。
（开票时 `Blocked by: nothing`，`internal/risk` 此刻无人写：票 75 已交、票 72 只碰过 `pathresolver.go`/`blacklist.go`）
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
- [x] **AC#2** **契约判定**：在 `docs/PLAN.md` 与 `docs/specs/SPEC-06*.md`（C19/C26 相关条目）里查
  `blacklist_overrides` 是否被写成**用户可用能力**。⇒ **已判（见 L3）：契约没有这条 ⇒ 停下上报编排者，不写实现码。**
  - 若**是**（契约要求它生效）⇒ 接线是**修 bug**，本票继续 AC#3。
  - 若**契约没有这条**、或它属于"放行侧放宽"⇒ **停下来上报编排者**，由编排者判是否 D22。
    ⚠ 不许用"先接上再说"绕过：放行侧的开关是安全边界，不是配置便利。
- [ ] **AC#3** 接线落地 + 一条**真跑在 `risk.Gate` 上**的用例：同一份配置里写 override，
  证明 (a) 被点名的那一个文件从 B 变成放行，(b) **同目录下其它文件仍然 B**，
  (c) override 的形状必须是 C26 交回的 canonical 形式（票 72 的 R17 裁定：
  **放行侧只认解析后的形式**，`bOverrides` 的 key 不接受拼写变体）。
  ⇒ **未做（被 AC#2 挡住）**：本框的 (a) 形状本身就是 PLAN.md:2399"必须由人点"禁止的那种静态预授权，
  没有编排者裁决（选项 A/B/C，见 L3）之前不许落地。
- [ ] **AC#4** 变异双向：(i) 把接线断开 ⇒ AC#3 必须红；(ii) 把 (b) 的"其它文件仍 B"断言指向 override
  ⇒ 必须红。两次都要在**同一条 `&&` 链里先 grep 证明改动真落地**再跑，还原后 `git diff --quiet` 证干净。
  ⇒ **未做**：AC#3 没有码可断，变异无承载体。
- [ ] **AC#5** 门禁（**只跑自己碰到的包**，共树不跑整仓）：`gofmt -l <pkgs>` 空、
  `go vet <pkgs>` rc=0、`go test -count=2 <pkgs>` rc=0，且**逐跑点名** `--- SKIP`/`--- FAIL` 的行数与测试名；
  `GOOS=linux go vet <pkgs>` rc=0（Linux 编译哑弹是 A44/A49 的老坑）。
  ⇒ **未闭框，但基线已跑**（见 L3）：本票**零 Go 文件改动**，门禁无改动面可验；
  仍把 `internal/risk` + `internal/config` 的现状数字记进票面，顺带交出一条**先存在的脆弱性能用例**。
  ⇒ **编排者补（A53）**：接手方 = **票 83 的 AC#5**（同两个包，而票 83 会真改 Go 文件 ⇒ 门禁有承载体）；
  AC#3/AC#4 的接手方 = **票 21**（审批队列），裁决走 A53② 的选项 (C) ⇒ **本票永远不勾这两框**，
  它们是**被裁决取消**、不是**没做完**。那条 `TestResolvePerCallBudget` 抖动登记为 **A53④**，单独一票，
  **现在不许把 1ms 预算调大**。

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
| `cmd/wisp/run.go:223-227, 243, 258-259` | — | 组合根只喂 `cfg.FS.AllowedDirs`/`ReparsePointExceptions`（`227`）、`cfg.FS.DeleteEnabled`（`243`）、`cfg.Risk.L1WindowSec`/`ConfirmTimeoutSec`（`258-259`）；**从不引用 `cfg.Risk.BlacklistOverrides`** |
| `internal/agent/approval/approval.go:273-320` | — | `grantStore` 是**一次性 nonce**（F2 第 3 层），不是持久豁免记录 ⇒ 全仓不存在"人点过 L2 之后把该文件记进某张 map"的容器 |

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

### L3 — AC#2 契约判定：**契约没有这条** ⇒ 停下上报编排者（2026-09-21）

检索式：`grep -rn "blacklist_overrides" docs/` ⇒ **全 docs 只有两处命中**，都是配置键表；
`grep -rn "单文件豁免\|豁免" docs/PLAN.md docs/specs/*.md` ⇒ 定义豁免语义的正文共四处（见下表 D1-D4）。
下列引用**逐条自己打开核对过**（票面原有的两处引用 `schema.go:457` / `manager.go:356` 也都核对为真）。

| 编号 | 位置（已打开核对） | 原文承担什么 |
|---|---|---|
| C1 | `docs/PLAN.md:2732` | 🔒 `[risk]` 配置表行里列出 `blacklist_overrides[]`，整行标记"变更即需 L2 级重新确认" |
| C2 | `docs/specs/SPEC-03-config-secrets-envs.md:34` | 同上：`blacklist_overrides[]=[]`（B 档豁免）｜🔒 放宽需 L2 重新确认 |
| D1 | `docs/specs/SPEC-06-security-gatekeeping.md:67` | "**B 档 · 默认拒绝 + 单文件豁免**（**豁免 = 一次 L2 强确认 + 写日志**）" |
| D2 | `docs/specs/SPEC-06-security-gatekeeping.md:36` | R3 行："B 档 → L2（可单文件豁免）" |
| D3 | `docs/PLAN.md:2396-2399` | "**但豁免必须按文件、必须由人点、必须留痕。**"（F1 CRITICAL 段） |
| D4 | `docs/specs/SPEC-06-security-gatekeeping.md:19` | L2 行的「允许」**只接受原生侧来源**（悬浮球点击/原生确认卡按钮/全局快捷键） |
| D5 | `internal/risk/assessor.go:144-146`（冻结，仅读） | "tier B is L2 with per-file exemption（**exemption flow lands with the approval queue, ticket 21**）" |
| D6 | `internal/risk/blacklist.go:13-14`、`:74-75` | `bOverrides` 的既定生产者 = "**individually L2-confirmed** B-tier files" |
| — | `docs/specs/SPEC-10-testing-acceptance.md` | **零命中**：没有任何验收条目要求 `blacklist_overrides` 在决策时生效 |

**判定：走票面 AC#2 的第二分支（契约没有这条 ⇒ 停下上报编排者，判是否 D22）。理由三条，都能回指上表：**

1. **契约把 B 档豁免定义成一个"运行时事件"，不是一份"静态预授权"。** D1 用等号写死了它：豁免 **=** 一次 L2 强确认 + 写日志；
   D3 再加三个"必须"，其中"**必须由人点**"直接排除了"TOML 里写一行、以后每次读取都不再有人点"这种形状。
   D5/D6 进一步把这股水流归给**审批队列**（ticket 21），而不是 config。
2. **`blacklist_overrides` 在契约里的全部存在形式就是两张配置键表（C1/C2），没有任何规格文字描述它的决策时语义** ——
   谁读它、key 要什么形状、与 L2 确认是什么先后关系、是否持久。SPEC-06（拥有该能力的规格）**一次都没提到这个键**。
   ⇒ 接上它**不是把已有契约落成实现**，而是新造一条契约没有写过的放行侧通道。
3. **本票 AC#3 的 (a) 形状恰好就是被 D3 禁止的那一种**："同一份配置里写 override ⇒ 被点名的文件从 B 变成放行"，
   全程无人点。⇒ 我不写 AC#3/AC#4 的实现码；这不是"洞没补上"，是**这个洞的正确处置需要先改契约，而契约是冻结面**。

顺带纠正一处票面前提（AC#1 已记，这里给编排者复述一遍）：生产判定链**根本没走 `risk.Gate`** ——
`cmd/wisp/run.go:262 tools.New(...)` → `bridge.go:159 NewSensitiveClassifier()` → `paths.go:146 risk.Classify(canonical)`，
而 `Classify` 的签名里没有 override 形参（`blacklist.go:69`）。所以"接线"要么在 `paths.go` 把 Classify 换成 Gate（非冻结，但
"B 且已豁免"只能塌成 `TierNone`，等于**丢掉 D1/D3 要求的留痕**），要么在 `rules_gateway.go:84` 的 `case TierB:` /
`assessor.go:158` 的 seam 上开一个能携带豁免集合的位置（**两处都是冻结面**）。

**交回编排者三选一（本票不代拍）：**

- **(A) 照 AC#3 原样接静态预授权** ⇒ 必须先改 `SPEC-06 §4.1` 的 D1 等号与 `PLAN.md:2399` 的"必须由人点" ⇒ **要 D22 人工批准**（冻结文档）。
- **(B) 接成契约真正写的那件事**：预登记 + **仍要人点一次**，确认后把 C26 canonical 记进 `bOverrides` ⇒ 落点在审批队列，
  生产链要从 `risk.Classify` 换到 `risk.Gate`，并要动 `assessor.go`/`rules_gateway.go` 的 TierB 语义 ⇒ **也要 D22**，
  且这是 ticket 21 那条"exemption flow lands with the approval queue"的欠账，不该塞进本票。
- **(C) 让这个键响亮地失败**（不接线，但写非空即报错/doctor 亮红灯），照本仓既有先例 `verify_signature`
  （SPEC-03:42"**不可关：写 false 报错而非生效**"）与 `keep_transcript`（写 true 报错）⇒ 这是**收紧侧**，不放大放行面，
  但同样改 SPEC-03 的键语义 ⇒ 需编排者点头。会把票面"后果 1：配置项说谎"就地消掉。
- **(D) 维持现状 + 登记**（当前交付）。

**基线门禁数字（本票零 Go 改动，`git status --porcelain` 在跑门禁前只剩本票面一处 M）**

| 命令 | rc | 真实输出/计数 |
|---|---|---|
| `gofmt -l internal/risk internal/config` | **0** | 空输出 |
| `go vet ./internal/risk/ ./internal/config/` | **0** | 无输出 |
| `go test -count=2 -v ./internal/config/ ./internal/risk/` | **1** | `ok github.com/CarlosShao/wisp/internal/config 0.975s`；`FAIL github.com/CarlosShao/wisp/internal/risk 17.114s`。`=== RUN` 460 行；`--- PASS` 267 + `    --- PASS` 190；`--- SKIP` **2 行**（`TestSyncRegistryProbeLive`，1 个测试 × count=2）；`--- FAIL` **1 行** = `TestResolvePerCallBudget`（`pathresolver_budget_norace_test.go:37` "Resolve averages 1.640779ms per call, budget 1ms"） |
| `go test -count=2 -run TestResolvePerCallBudget -v ./internal/risk/` | **0** | 2 个 `=== RUN` / 2 个 `--- PASS`，实测 0.692 与 0.637 ms/op ⇒ 上一条的 FAIL 是**负载敏感的抖动**，不是断言写反 |
| `go test -count=2 ./internal/risk/`（单独重跑） | **0** | `ok github.com/CarlosShao/wisp/internal/risk 6.806s` |
| `GOOS=linux go vet ./internal/risk/ ./internal/config/` | **0** | 无输出（Linux 编译哑弹已查） |

`-count=2` 的 N 倍核对：`--- SKIP` 恰好 2 行（同一测试 2 次重复）、两次独立 RUN 的 PASS 计数 267+190 与单包重跑一致 ⇒ 未见"`-run` 空匹配 rc=0"或"步骤被静默跳过"。
**唯一假绿风险点名**：`TestOverrideKeyKeepsPosixBackslashesDistinct`（`pathshape_portable_test.go:192`）在 Windows 上提前 `return` ⇒ 本平台**它不覆盖跨文件 bleed**，别把它的 PASS 当豁免接线的证据。

**登记给编排者的先存在外事项**（不在本票改动面，也**不**归票 70 —— 它只做格式化）：
`TestResolvePerCallBudget` 是**墙钟 1ms 预算**的性能契约（票 18 AC#6）。它和整包并发跑会抖到 1.6ms 而判红，
⇒ 今天的门禁只要"两包并跑"就会假失败。建议单独一票（改成 `testing.Short()` 门控 / 或把预算改成同机基线倍数），
本票**不动它**。

`next=` 交回编排者拍 (A)/(B)/(C)/(D)；未拍之前 AC#3、AC#4 保持未闭、本票不产出实现码。
