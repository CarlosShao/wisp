# 358 — 压缩"该压但压不动"这一形今天不留痕（按主流 harness 口径补读数，⛔ 不动 D15 预算契约）

- Status: **ready-for-agent**（`AC#0` 是只读现量，可先派；`AC#2` 动码、须先过 `AC#1` 裁形）
- 立案：2026-10-21 14:4x +08 编排者立。**机主拍板原话**：`Q-312` 答「我觉得还是得动，千万不要因为有风险就不去做，因为这个契约是早先定下的，我不知道符不符合主流harness，你最好还是按照主流harness的方式去搞吧」。
- ★**本票的立票前提里，有一处我把机主的话往下拽了一格，具名记在这里**：
  机主说的是"动契约"，但我**到盘上核完之后**发现**这个缺口不必动契约就能修**（理由见下「技术前提」第 1 条），
  所以本票按"**不动 D15 预算语义、只补可观测面**"立；**D15 契约文字要不要跟着改，单列 `AC#4` 摆机主**，⛔ 我替机主拍。
- 非实现者签名＝`139-v1` 判语第 5 件 (a)（`docs/evidence/s1/139-v1-verdict.md`）：它判"⛔ 立实现票、⛔ 挂回票 139 `AC#2`"、只登已知残缺＋开 `Q##`。
  机主 `Q-312` 答"动"之后，本票就是那条残缺的实现票。
- 关联：票 139（已 `-done`，痕的产码面在那里结案：`compress.go:240` 那枚 `slog.Info` ＋ 7 枚用例）· 票 355（痕到得了盘但盘上零读者）· `Q-312`

## 技术前提（全部我现读盘取的，⛔ 凭印象）

**第 1 条（决定本票形状）＝真实"压不动"的位置在 `Compress` 内部，不在 `Need`／`KeepRawRounds` 的语义里。**
`internal/agent/compress.go:180-183` 逐字：
```
	for c.totalOf(rounds) > c.b.HistoryCompressTokens {
		raw := ...
		if len(raw) <= c.b.KeepRawRounds {
			break // D15(4): stop at the last 3 raw rounds
		}
```
⇒ "超了阈值但只剩那 3 轮必须保留的原文"这个形状，是**循环自己 `break` 出来的**，
而不是 `Need()` 判错、也不是 `KeepRawRounds` 的取值错。
⇒ **补一条痕不需要改 D15 的任何语义**（⛔ 改 `Need()` 的判据、⛔ 改 `KeepRawRounds=3`、⛔ 改 12000／4000 两个阈值）。

**第 2 条＝`Need()` 今天用的是**本地估算**，不是 provider 回报的实际数。**
`compress.go:119-121` 逐字 `func (c *Compressor) Need(hist []llm.Message) bool { return c.TotalTokens(hist) > c.b.HistoryCompressTokens }`，
而 `TotalTokens` 走 `ApproxTokensOf`（`budgets.go:132`）。
⇒ 主流 harness 今天普遍改用 **provider 回报的实际 input token** 当触发基数（Cline 的 `PR #14195` 就是这个改动；
OpenAI Codex 侧有两条 issue 记的正是"用陈旧 token 数触发压缩导致不可恢复的上下文溢出"）。
★**但这一条属**触发面**、不属**可观测面**——它要动 `Need()` 的输入，风险与收益都另算 ⇒ **本票⛔ 做它**，单列 `AC#5` 摆机主。

**第 3 条＝D15 那条"必须按 `context_window` 等比缩小"今天**已实现**（我核过，⛔ 当缺口）。**
`budgets.go:10-18` 有 `ReferenceContextWindow = 128000` 与缩放函数 `s()`（`b.HistoryCompressTokens = s(refHistoryTokens)`）。

## 主流 harness 的三条做法（本票照它定形状）

1. **痕要"事前＋事后"两个数**，不只事后一个：主流做法是把"压前占多少、压后占多少"都端出来（Claude Code 的 context-window 文档把这件事写成用户可见的 compaction 边界）。⇒ 本票要求三味：**压前 token／压后 token／是否真动了历史**（与票 139 `AC#2` 同口径）。
2. **"没压动"必须能与"没触发"分开**：主流实现里 compaction 有独立的"attempted / skipped / compacted"三态；今天 Wisp 只有 `Ran` 一个布尔（`compress.go:62` `Ran bool`）⇒ **attempted-but-skipped 这一态不存在**。
3. **读数落点必须在**已有出口**上**，⛔ 新造一条前端通道（主流的读数是日志／状态栏／诊断包，⛔ 新事件类型）⇒ 与票 355 的落点结论同源。

## 结案判据

- [ ] **AC#0 现读＋归因（只读，⛔ 动任何产码字节）**：把"该压但压不动"这一形**量成可复跑的读数**——
      造一个"超过阈值、但可压部分已空、只剩 `KeepRawRounds` 那 3 轮"的历史，
      现量：`Need()` 返回什么、`Compress` 返回的 `CompressionReport` 是什么（三味各是什么）、**痕有没有落**、落了几条。
      ★必须先跑**正控**（一个"真压得动"的历史 ⇒ 痕必须落）与**负控**（一个"没到阈值"的历史 ⇒ `Need()` 必须假），否则你的台子测不出发育不良。
- [ ] **AC#1 裁形（编排者＋机主）**：三形各写代价＝〔甲〕**只给 `Ran` 加一个伴生字段**（如 `Attempted`／`Stop`），⛔ 碰 `Need()`；
      〔乙〕把三态（attempted／skipped／compacted）显式做成 `CompressionReport` 的一枚枚举；〔丙〕只在日志那一行上追一个 `skipped_reason` 键。
      ★**未裁形任何实现都⛔ 开工**。判据：三形都要回答"**主流 harness 是哪种、我们为什么选这个**"。
- [ ] **AC#2 落地＋牙（写腿）**：按裁定那一形做；每处改动配**一枚今天响、不修不响**的用例；
      ★反向判据＝把新加的那一跳摘掉 ⇒ 用例必须转红并**逐名报红行**；先证变异落地（`grep -n` ＋ `go build` rc=0）再读数，三读数逐字贴。
      ★**⛔ 动 `Need()` 的判据、⛔ 动 `KeepRawRounds` 的值、⛔ 动 12000／4000 两个阈值**（那是 D15 预算语义）。
- [ ] **AC#3 门禁**：`GOFLAGS= go build ./...`／`go vet ./internal/agent/`／`go test -count=1 ./internal/agent/`（四数分开报、⛔ 相减）／`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh` 各 rc。
- [ ] **AC#4 D15 契约文字要不要跟着改（**摆机主**，⛔ agent 自决）**：本票做的是"补可观测面"，D15 的**预算语义一字不动** ⇒
      契约文字**原则上不需要改**。若机主认为"留痕"本身应当写进 D15，则另开一条契约修订（须人工批准原文）。
- [ ] **AC#5 触发基数要不要改用 provider 回报的实际 token（**摆机主**，⛔ 本票做）**：今天 `Need()` 走本地估算（见技术前提第 2 条）；
      主流做法在往"用 provider 实际 input token"走。★这是**另一条链**（触发面），改它要动 `Need()` 的输入契约；
      ⇒ 本票只**登记**，⛔ 顺手做。
- [ ] **AC#6 纪律格**：⛔ 翻任何框（本票各格由非实现者裁）、⛔ 改票名／`Status:`／`-done`、⛔ push、⛔ 碰票 139／355 任何原句。

## 硬约束

- ★**⛔ 改 D15 的任何预算语义**：`Need()` 判据、`KeepRawRounds = 3`（`budgets.go:42`）、`refHistoryTokens = 12000`、`refTokenBudget = 200000`、`ReferenceContextWindow = 128000`（`budgets.go:16-18`）—— ⛔ 碰这些数字与判据。
- ⛔ 改 `docs/PLAN.md`／`docs/specs/**`（契约文字改动须人工批准、另票）。
- ⛔ 新造 Go→前端事件／新 `panel.*` 方法／新 JSON 键（`SPEC-08` 第 5.2 节白名单＝契约变更）⇒ 碰到就停手上报。
- ⛔ `internal/risk/**`、`tools/d22scan/**`、`allowlist.txt`、`thresholds.go`、任何 golden、`frontend/**`、`design/**`。
- ★**写面互斥**：`internal/agent/**` 同包只许一枚写腿。

## 排程

`358-a1`（只读＝`AC#0` 三态现读）→ 编排者＋机主裁 `AC#1` 三形 → `358-r1`（写腿＝`AC#2`）→ `358-v1`（非实现者裁）。
`AC#4`／`AC#5` 两格各自摆机主，⛔ 与实现同批。

**Status 追加（2026-10-21 14:4x）**：号位＝编排者定（`354`/`355`/`356`/`357` 已占，`358` 现量空号）。⛔ 零翻框、零 push。

## Progress log

（追加式；每条 `- [日期 时刻 +08] agent=… did=… next=…`，钟点由 `date` 插值，不手打。）

- [2026-10-21 14:4x +08] agent=`编排者` did=立案（机主 `Q-312` 答"动"；★我把"动契约"往下拽了一格＝查完发现**不必动契约就能修**，理由逐字写进技术前提第 1 条；另两条技术前提：`Need()` 今天走本地估算、D15 的 `context_window` 等比缩小**已实现**）next=派 `358-a1` 只读出三态现读；`AC#1` 三形表回来后摆机主裁形。
