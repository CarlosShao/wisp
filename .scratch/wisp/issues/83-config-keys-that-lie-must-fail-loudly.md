# 83 — 配置里"看着能用、其实没接线"的键要**响亮地失败**（票 80 裁决 (C)；先例 `SPEC-03:42 verify_signature`）

**Status:** claimed（票 83 由写码代理认领；AC#1 键清单扫描中）
**Type:** 安全可用性/诚实性（一个说谎的配置键）——**不是**新能力
**Blocks:** nothing · **Blocked by:** nothing（`internal/config` 此刻无人写；票 80 已交回且零 Go 改动）
**Packages:** `internal/config/`（校验与加载路径）+ 新建的用例。**禁改**：`docs/PLAN.md`、`docs/specs/*.md`
（冻结面，那两行限定语由编排者落）、`internal/risk/assessor.go`、`internal/risk/pathresolver*.go`、
`rules_gateway.go`、`internal/agent`+`internal/memory` 的测试（票 81）、`internal/risk` 的 sync 家族测试（票 82）。

## 依据：票 80 的穷尽清单（读/写逐条在票面 L2，证据 `docs/evidence/s1/80-ac1-ac2-verdict.md`）

`blacklist_overrides` 的真相**比建票时写的更空**：
- 写：`internal/config/schema.go:457`（TOML 解码）；读：**只有** `manager.go:356` 的方向审计。**没有任何消费者**。
- `risk.Gate` 在生产里**零调用点**。真实判定链是 `cmd/wisp/run.go:262 tools.New` →
  `internal/tools/bridge.go:159 NewSensitiveClassifier()` → `internal/paths.go:146 risk.Classify()`，
  而 `Classify`（`blacklist.go:69`）**签名里没有 override 形参** ⇒ "接线"不是灌一个 map，是要先在判定链上
  **开出能携带豁免集合的位置**。
- 该键唯一残存效果（🔒 变更检测，`manager.go:239`）也是哑的：`ConfirmLocked` 生产从未赋值。

**票 80 的 AC#2 裁决（编排者已签）**：契约**从未承诺**这个键在决策时生效——
`SPEC-06:67` 把豁免定义成**运行时事件**（"豁免 = 一次 L2 强确认 + 写日志"）、`PLAN.md:2396-2399` 明写"**必须由人点**"、
`assessor.go:144-146` 的冻结注释把 exemption flow 推给票 21。⇒ 接静态预授权 = **新造放行侧能力**（选项 A/B，需 D22，
且本来就是票 21 的欠账）。本票走**选项 (C)**：让键**响亮地失败**。

## 为什么 (C) 而不是 (D) 维持现状

现状是"配置项说谎"：用户写下一行他认为在放开某个文件的键，程序默默吃掉它。
先例就在同一份契约里：`SPEC-03:42` 的 `verify_signature` 写 `false` 时报错而非静默生效。
(C) 是四个选项里**唯一收紧侧**的那个，且**不需要**动安全语义、不需要 D22 批准新能力。

## AC（1:1，裁决表 `docs/evidence/s1/83-*.md`）

- [ ] **AC#1** **先把同类键扫全**（不要只修一个键就交）：`grep` 出 `internal/config` 里每个被解析、
  但在 `internal/config` 之外**零消费者**的键。票 80 已经点了三个候选：`[risk] shell_enabled`、
  `allow_shell_string`、`shell_allowlist`（`grep "\.Risk\."` 全仓只命中 `L1WindowSec`/`ConfirmTimeoutSec`），
  以及 🔒 段的 `ConfirmLocked`（D36 规则 1，生产未赋值）。
  **交付物 = 一张表**：键 → 解析处 → 消费者（"无"也要给 grep 证据）→ 处置（响亮失败 / 已有计划接（指名票号）/ 保留但文档说明）。
  ⚠ "grep 零命中"在这里**不足以**定罪：要顺带查是否存在**间接消费**（反射、按名取值的 map、TOML 原样透传）。
- [ ] **AC#2** 校验落地：对表里判为"响亮失败"的键，**加载时**报错并指名键名 + 一句"该能力尚未实现/由票 N 实现"，
  **错误发生在任何副作用之前**（照 `verify_signature` 先例的形状）。
  判据用例必须包含：**不写这个键的用户不受影响**（默认配置仍加载成功）。
- [ ] **AC#3** 双向变异：(i) 把校验去掉 ⇒ 用例红；(ii) 把校验做宽（例如对所有未知键报错）⇒
  必须有用例**因此变红**（防止"用一个大棒假装修好了"）。锚点=承载行为的那一行，同链 grep 自证，还原后 `diff -q`/`git diff --quiet`。
- [ ] **AC#4** 明确**不做**：不接 `risk.Gate`、不给 `Classify` 开 override 形参、不动审批队列。
  在票面写一行"这些属于票 21 / 需 D22"，并把它链接到 A51⑤ 与票 80 的裁决。
- [ ] **AC#5** 门禁（只跑自己碰的包）：`gofmt -l` 空、`gofumpt -l <files>` 空、`go vet ./internal/config/` rc=0、
  `GOOS=linux go vet ./internal/config/` rc=0、`go test -count=2 ./internal/config/` rc=0，
  逐跑点名 `--- SKIP`/`--- FAIL`。**并且**：如果两包并跑出现墙钟抖动（票 80 量到 `TestResolvePerCallBudget`
  在 `config`+`risk` 并跑时红、单跑绿），**如实登记为外项**，不许顺手调那个 1ms 预算。

## Rules（本仓固定）

15 次工具调用内交第一枚 checkpoint commit；每次 commit 同步 Status + 勾框 + 末行 `next=`；
`git commit -q -F - -- <显式路径>` + **带引号** heredoc；禁 `git add -A`/`.`；commit 前核对 `git diff --cached --name-only`；
禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`（共树 A34）；**不 push**；不在仓内建 worktree（A38④）；
票面 Progress log append-only，**要改的那行先读再替换**；四种假绿逐跑点名。
**spec 那两行限定语（`SPEC-03:34` 键表加"未实现 ⇒ 加载报错"）你不要写** —— 把拟好的文字放进票面，编排者落。

## Progress log（append-only）

### L0 — 认领 + checkpoint（2026-09-21）

Status `open` ⇒ `claimed`。本轮次序：AC#1 全量同类键扫描（含间接消费排查）⇒ AC#2 校验落地 ⇒ AC#3 双向变异 ⇒ AC#5 门禁。
零 Go 改动，纯 checkpoint。
`next=` 读 `internal/config/` 全量，产出键 → 解析处 → 消费者 → 处置表。
