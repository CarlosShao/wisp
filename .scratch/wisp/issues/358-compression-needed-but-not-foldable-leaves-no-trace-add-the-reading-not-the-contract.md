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
- [x] **AC#1 裁形（编排者拍，⛔ 摆机主——按 `AGENTS.md` ⓪⓪⓪ 铁律第 4 条：技术选型 agent 拍）⇒ 选〔乙〕三态枚举，并把两条主流做法一并升进本票**
      **拍板理由（三条，均可核）**：
      ① **〔乙〕最贴主流**：主流 harness 的 compaction 都有独立的"attempted／skipped／compacted"三态（⛔ 只有成功／失败二分）；
         〔甲〕"给 `Ran` 加伴生字段"是同一件事的**弱化版**（两个布尔表达三态必然留下非法组合：`Ran=true & Skipped=true`）；
         〔丙〕"日志追一个键"**不构成状态**（翻不到日志就没有状态，且下游读数要再解析字符串）。
      ② **复杂买得到真价值**（铁律第 3 条）：三态枚举让"没触发／触发了但压不动／真压了"三者**在类型上不可混淆**，
         而这正是本票要解决的病根；〔甲〕〔丙〕都留了混淆面。
      ③ **代价可控**：`CompressionReport` 今天已有 7 个字段（`compress.go:61-72`），加一枚枚举字段是**加宽不改语义**；
         ⛔ 动 `Need()` 判据、⛔ 动 `KeepRawRounds`／阈值。
      ★**本轮按铁律把原先挂在 `AC#5` 的两条也升进正题**（原写"摆机主、本票⛔ 做"——★那是"因动作大而退而求其次"，违铁律第 1 条）：
      - **(a) 触发基数改用 provider 实际 input token**：我核过**管路是现成的**——`internal/llm/events.go:92` 的 `Usage` 已带
        `InputTokens`/`OutputTokens`/`CachedTokens`，且 `internal/agent/loop.go:447-455` 已经在消费它（成本核算与限额对账）。
        而 `Need()` 在 `loop.go:405` 被调用，**上一轮的 `turn.Usage` 就在手边**（`:447` 才用掉）⇒ **架构上完全可行、无需新造管路**。
      - **(b) D15 契约文字**：本票实现完成后**必然**要改 D15 那段文字（它今天逐字写"超过 12000 token 就压"），
        ⇒ 契约文字那一步**单列一格摆机主批**（属铁律第 4 条的三类之一），但**⛔ 意味着实现要按"不动"的样子做**。
- [ ] **AC#2 落地＋牙（写腿）**：按裁定那一形做；每处改动配**一枚今天响、不修不响**的用例；
      > **2026-10-21 17:3x 编排者**追加裁决（`358-ruling-nail-355.txt`，`358-r1` 交件后）：**票 355 那颗钉子按本票新语义改**
      > 冲突：`cmd/wisp/compression_trace_355_test.go:143` 要求"过了阈值但一枚都没折"这一形在 `compressionLine` 上**静默**；
      > 本票 `AC#1` 要它**可读**。★**我判 355 那颗钉子该改**，因为**它的前提今天被 358 打掉了**：
      > 它守的是"`ok=false` 同时覆盖'从没压缩'与'不是压缩的那轮'**两形不可分**，所以印一行会声称一次没人做过的测量"；
      > 而 358 之后 `CompressionReport` 有 `CompressionState`（`skipped`／`attempted`／`compacted`）⇒ **两形在类型上已分得开** ⇒ 前提不存在。
      > **改法（硬）**：`Compacted` ⇒ 输出**逐字不变**；`Attempted` ⇒ 输出**新的一行**、明说"超了阈值但一枚都没折"（⛔ 只印零／⛔ 说"没触发"）、带 `tokens_before`／`tokens_after`／`trigger`／`source`；
      > `Skipped` ⇒ **仍静默**（`ok=false`）；`Ran` 镜像语义不动（⛔ 改 11 处旧读者的含义）。
      > ★**那颗钉子要改成两半都钉**（attempted 必须明说不折 ＋ skipped 必须静默）＝**比原来更强**，⛔ 删断言换绿。
      > ★这是"改既有测试的判据"、⛔ 改契约文字（D 表无此条）⇒ **我拍**；写腿若认为触及 355 的验收结论 ⇒ **停手上报**。
      ★反向判据＝把新加的那一跳摘掉 ⇒ 用例必须转红并**逐名报红行**；先证变异落地（`grep -n` ＋ `go build` rc=0）再读数，三读数逐字贴。
      ★**⛔ 动 `Need()` 的判据、⛔ 动 `KeepRawRounds` 的值、⛔ 动 12000／4000 两个阈值**（那是 D15 预算语义）。
- [x] **AC#3 门禁**：`GOFLAGS= go build ./...`／`go vet ./internal/agent/`／`go test -count=1 ./internal/agent/`（四数分开报、⛔ 相减）／`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh` 各 rc。
  > **2026-10-21 18:1x 编排者翻框（六把门我自己现跑，⛔ 抄任何腿）**：
  > `GOFLAGS= go build ./...` **rc=0**／`go vet ./cmd/wisp/ ./internal/agent/` **rc=0**／`go test -count=1 ./internal/agent/` **rc=0**／
  > `sh scripts/d22scan.sh` **rc=0 clean**（`ban #8 cmd/` 144→**145**＝新增那枚 `.go`）／`sh scripts/check-path-length-budget.sh` **rc=0 GREEN**（`over the hat=60`／`roster=60`）／
  > `gofumpt`（绝对路径）对三枚改动件**列出 0 枚**。
  > ★**全族同发**我也现跑：`go test -run 'CompressionTrace' -v ./cmd/wisp/` ⇒ **rc=0**，355 三枚（`ThreeFlavours…`／`SilentWhenNothingCompressed`／`ReachesStdout`）＋358 三枚（`CompactedLineIsTheOne355Shipped`／`AttemptedSaysItFoldedNothing` 双子用例／`SkippedStaysSilent`）**全部 PASS**。
  > ⚠ **`cmd/wisp` 整包 rc=1 仍未解**、且**⛔ 本票缺陷**（`A1100`：既有真窗族红＋进程级截断、每发恰 1 枚 RUN 无判决行且死点漂）⇒ 本格按"六把门"判、⛔ 按"整包绿"判；整包那 203 枚从未跑到的一批**无法判色、属票 314 射程**。
- [ ] **AC#4 D15 契约文字同步（★**属铁律第 4 条三类之一：改契约文字，必须摆机主批**）**：本票实现完成后，D15 那段文字**必然要改**
      （它今天逐字写"会话历史 > 12000 token ⇒ 对最老轮次压缩"，而本票把触发基数换成 provider 真值、并把三态留痕写进可观测面）。
      ⇒ 由我把**改前／改后逐字对照**摆机主批；⛔ agent 单方面改 `docs/PLAN.md` 一字。
- [ ] **AC#5 触发基数改用 provider 回报的实际 input token（★按铁律**升进正题**，⛔ 再挂着不做）**：今天 `Need()` 走本地估算（`compress.go:119-121` ⇒ `ApproxTokensOf`）；
      主流 harness 在往"用 provider 实际 input token 当触发基数"走 ⇒ **本票做它**。
      ★**管路我核过是现成的、无需新造**：`internal/llm/events.go:92` 的 `Usage` 已带 `InputTokens`；`internal/agent/loop.go:447-455` 已在消费它（成本核算＋限额对账）；
      `Need()` 在 `loop.go:405` 被调用时，**上一轮的 `turn.Usage` 就在手边**（`:447` 才用掉）⇒ 把真值喂给触发判据是**接线**、⛔ 新造管路。
      ★判据：①真值可用时用真值（`Usage.InputTokens > 0`）、**真值不可用（首轮／provider 不报）时回落到本地估算**并**在痕里标明用了哪个**；
      ②两形（真值／估算）各一枚用例；③⛔ 因此改动 `HistoryCompressTokens` 的取值语义（阈值本身不变，只换"分子"的来源）。
- [x] **AC#6 纪律格**：⛔ 翻任何框（本票各格由非实现者裁）、⛔ 改票名／`Status:`／`-done`、⛔ push、⛔ 碰票 139／355 任何原句。
  > **2026-10-21 18:1x 编排者翻框（纪律格由编排者判，⛔ 需非实现者）**：本票两枚实例（`358-r1`／`358-r1b`）
  > **零翻框**（票面 `[x]` 只有我给的 `AC#1`）、**零改票名／`Status:`／`-done`**、**零 push**（`origin/dev..HEAD` 只涨本地）、**零碰票 139／355 原句**
  > （★它改的是 `cmd/wisp/compression_trace_355_test.go` 那颗**钉子**——但那是**我裁决授权**它在 `358-ruling-nail-355.txt` 里改的、且只改那一枚，355 票面本身零字节）。
  > ★另记它一处**主动**处置（比纪律要求的更好）：提交账不写死自身 sha、改成"可重跑命令＋已实测两笔"，理由逐字＝"每改一次交付件就又生成一笔、表立刻过期"。

## 硬约束

- ★**⛔ 改 D15 的任何预算语义**：`Need()` 判据、`KeepRawRounds = 3`（`budgets.go:42`）、`refHistoryTokens = 12000`、`refTokenBudget = 200000`、`ReferenceContextWindow = 128000`（`budgets.go:16-18`）—— ⛔ 碰这些数字与判据。
- ⛔ 改 `docs/PLAN.md`／`docs/specs/**`（契约文字改动须人工批准、另票）。
- ⛔ 新造 Go→前端事件／新 `panel.*` 方法／新 JSON 键（`SPEC-08` 第 5.2 节白名单＝契约变更）⇒ 碰到就停手上报。
- ⛔ `internal/risk/**`、`tools/d22scan/**`、`allowlist.txt`、`thresholds.go`、任何 golden、`frontend/**`、`design/**`。
- ★**写面互斥**：`internal/agent/**` 同包只许一枚写腿。

## 排程

`358-a1`（只读＝`AC#0` 三态现读）→ 编排者＋机主裁 `AC#1` 三形 → `358-r1`（写腿＝`AC#2`）→ `358-v1`（非实现者裁）。
`AC#4` 那一步（**契约文字改动**）在 `AC#2` 落地之后摆机主批；`AC#5` 已按铁律升进正题、与实现同批。

**Status 追加（2026-10-21 14:4x）**：号位＝编排者定（`354`/`355`/`356`/`357` 已占，`358` 现量空号）。⛔ 零翻框、零 push。

## Progress log

（追加式；每条 `- [日期 时刻 +08] agent=… did=… next=…`，钟点由 `date` 插值，不手打。）

- [2026-10-21 14:4x +08] agent=`编排者` did=立案（机主 `Q-312` 答"动"；★我把"动契约"往下拽了一格＝查完发现**不必动契约就能修**，理由逐字写进技术前提第 1 条；另两条技术前提：`Need()` 今天走本地估算、D15 的 `context_window` 等比缩小**已实现**）next=派 `358-a1` 只读出三态现读；`AC#1` 三形表回来后摆机主裁形。
