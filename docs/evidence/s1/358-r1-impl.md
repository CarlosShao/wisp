# 358-r1 落地件：压缩三态 与 触发基数换成 provider 真值

身份＝写腿（实现者）。不出终裁、不翻框、不改票名/Status:/-done/台账。
派单＝`.scratch/wisp/probes/orch/20261021/txt/358r1-dispatch.txt`（`c50218bf`）｜锚＝`c50218bf`｜钟点 `2026-10-21 16:0x-17:2x +08`。

## 0 开工前（不许跳步）

- `git status --porcelain -- internal/agent/ cmd/wisp/ docs/evidence/s1/` 结果＝**空**（写面内无别人未提交改动）。
- 票面三条前提**逐条现跑复认，全成立**：(1) 真实"压不动"位置＝`if len(raw) <= c.b.KeepRawRounds { break }`（票写 `:180-183`，今 `:343-347`，漂 +163），是循环自己 `break`，NOT `Need()` 判错、NOT `KeepRawRounds` 取值错；(2) `Need()` 走 `TotalTokens` 得 `ApproxTokensOf`（`:242-244`）；(3) `Need` 在 `loop.go:413` 调用时**上一轮 `turn.Usage` 已在手边**（`:467` 才用掉），接线可行，NOT 需新造管路。

## 1 三态今天各是什么读数（现量，-overlay 探针，不落产码字节）

台子＝`BudgetsFor(4096)` 得阈值 **384**（从预算读回，不手写）。探针 `D:\tmp\wisp\20261021\358r1\zz358r1probe_test.go`，原文 `logs/probe-after.txt`。

| 形 | 历史 | Need() | Ran | Before 到 After | 痕（旧） | 痕（新） |
|---|---|---|---|---|---|---|
| (1) 没到阈值 | 2 轮 | **false** | false | 32 到 32 | **0 条** | **0 条**（skipped 不留痕＝正控） |
| (2) 真压了 | 6 轮 | true | **true** | 1236 到 780 | 1 条 | 1 条 `compacted` |
| (3) 超阈但只剩 KeepRawRounds | 3 轮 | **true** | **false** | 618 到 618 | **0 条** | **1 条 `history compress skipped`** |
| (3b) 单轮（票 153 那形） | 1 轮 | true | false | 1006 到 1006 | **0 条** | 1 条 `skipped` |
| (3c) provider 说超/估算没超 | 2 轮 | true | false | 32 到 32 | **0 条** | 1 条（reason=under_fold_budget） |
| (4) 折了但仍超（贴地板） | 10 轮 | true | true | 10060 到 3359 | 1 条 | 1 条 `compacted` |

**病根**：形 (1) 与 (3) 改前**完全同形**（Ran 都 false、痕都 0 条），于是"该压但压不动"在**每个出口**都读成"没触发"。

## 2 选了哪一形 与 为什么

按已拍的 **乙案三态枚举**，不重议、不降级。`CompressionState` 三值 `skipped`/`attempted`/`compacted`；
**空串不是第四态**：`Valid()` 返 false，未填的 `CompressionReport{}` 读成"压根没跑"，NOT 读成"看了看决定不压"。
**为什么甲案不行（今天可核）**：`Ran` 已被 **11 处**读（`loop.go:415`、`run.go:1318`、7 枚 139/355 钉子）；加伴生布尔表达三态会留下**非法组合** `Ran=true 且 Skipped=true`，且无一处能拒它，编译器不管；枚举让三态**类型上互斥**。
**`Ran` 保留为 compacted 的镜像**（`rep.Ran = rep.State == CompressionCompacted`），由 `TestCompressionStateRanIsTheMirrorOfCompacted` 钉住，故旧读者语义不变、权威归枚举。

## 3 改了哪几行（逐行）

`internal/agent/compress.go`（+267/-2）：`:42-79` 三值枚举 加 `Valid()`/`String()`｜`:82-134` `TriggerSource` 两值｜`:136-147` `TriggerReading{Need,Tokens,Threshold,Source,Observed}`｜`:149-160` 两个 reason 常量｜`:162-194` 报告加 `State`/`TriggerSource`/`TriggerTokens`（`Ran` 注释改镜像说明）｜`:252-282` `Need()` 保留（循环不再用）加新增 `Trigger(hist, observed)`｜`:334-352` `CompressOption`/`WithTrigger`（可变参，旧调用点零改动）｜`:354-397` 算 `entered`（调用方 trigger **或**本趟估算过预算）｜`:441-460` **三态判定唯一一处**（Ran 得 compacted，entered 得 attempted，否则 skipped；随后 `rep.Ran` 由 State 反写）｜`:461-500` compacted 痕加 `state`/`trigger_source`/`trigger_tokens`（余键原样）｜`:501-535` **新增 attempted 痕** `agent: history compress skipped` 加 `reason` 加全量计数 加 trigger 两键 加 `task`。

`internal/agent/loop.go`（+34/-6）：`:382-387` `var prevUsage llm.Usage`（首轮零值故自动落估算支）｜`:405-437` `Need(hist)` 换成 `Trigger(hist, prevUsage.InputTokens)`、`Compress(..., WithTrigger(trig))`、**非 Ran 也写 `res.Compression`**（attempted 才上得了出口）｜`:467-480` `turn.Usage.InputTokens > 0` 时 `prevUsage = turn.Usage`。

**没动的**：`Need()` 阈值语义、`KeepRawRounds=3`、`refHistoryTokens=12000`、`refTokenBudget=200000`、`ReferenceContextWindow=128000`（`budgets.go` **零字节**）；`cmd/wisp/**` 最终**零改动**（见 sec 7）；`docs/PLAN.md` 与 `docs/specs/**` 零字节；`internal/{tools,config,winsec,risk}/**`、`tools/d22scan/**`、`allowlist.txt`、`thresholds.go`、golden、`frontend/**`、`design/**` 零字节。未新造事件、未加 `panel.*` 方法、未加 JSON 键。

## 4 新用例逐名 与 为什么今天响

`internal/agent/compress_state_358_test.go`（8 枚）：`TestCompressionStateAttemptedIsNotSkipped`（改前无 State，形(3)落 skipped）｜
`TestCompressionStateSkippedStaysSilent`（正控反向：形(1)必须 skipped 且**两条痕都空**，防"无条件留痕"恒真）｜
`TestCompressionAttemptedLeavesItsOwnTrace`（改前形(3) trace_landed=0，断言"恰好 1 条"必红；另按名钉 12 键 加隐私哨兵）｜
`TestCompressionStateRanIsTheMirrorOfCompacted`（三形各一枚，钉镜像关系）｜
`TestCompressionTriggerPrefersProviderUsage`（改前无 Trigger；**两向**都测：估算低/真值高、估算高/真值低，单边实现过不了）｜
`TestCompressionTriggerFallsBackToEstimateOnFirstTurn`（observed=0 必落估算，且与 `Need()` **逐值相等**）｜
`TestCompressionAttemptedWhenProviderSaysOverButFoldHasNothing`（第二种 attempted 形，reason=under_fold_budget）｜
`TestCompressionAttemptedWithoutCallerTrigger`（裸 `Compress` 也 attempted 且自报 estimate）。

`internal/agent/compress_trigger_358_test.go`（5 枚）：`TestLoopTriggerUsesProviderReportedInputTokens`（真 `Loop.run` 加 `tool-then-text.sse`，窗口 1024 得阈值 96；种子估算 32 < 96 < provider 150）｜
`TestLoopTriggerFallsBackToEstimateWhenNothingObserved`（128k 加 15 轮重历史，首轮无观测故 estimate 且 TriggerTokens==TokensBefore）｜
`TestLoopTriggerBasisCanChangeBetweenTurns`（同一历史两基数必须**给出不同读数**，防"算一次钉死"）｜
`TestCompressionStateStrings`（空串 `Valid()==false`）。

## 5 三读数逐字 与 sha256 还原

每枚突变**先证落地**（`grep -n` 原文 加 `go build` rc=0），还原＝按字节拷回快照再核 sha256。

| 变 | 突变（grep -n 原文） | build | 变异红（逐名） | 还原绿 |
|---|---|---|---|---|
| M1 | `446: _ = entered`（三臂改两臂，attempted 臂删） | **0** | StateAttemptedIsNotSkipped / AttemptedLeavesItsOwnTrace / StateRanIsTheMirror / AttemptedWhenProviderSaysOver / AttemptedWithoutCallerTrigger / LoopTriggerUsesProviderReported | 0 |
| M2 | `371: // a trace of its own (...)`（skip 痕 1188 字节删） | **0** | AttemptedLeavesItsOwnTrace / AttemptedWhenProviderSaysOver / LoopTriggerUsesProviderReported | 0 |
| M3 | `269: if false && observed > 0 {` | **0** | TriggerPrefersProviderUsage / AttemptedWhenProviderSaysOver / LoopTriggerUsesProviderReported / LoopTriggerBasisCanChange | 0 |
| M4 | `prevUsage = turn.Usage` 删（grep 零命中） | **0** | LoopTriggerUsesProviderReported | 0 |
| M5 | `418: trig := l.comp.Trigger(hist, prevUsage.OutputTokens)` | **0** | LoopTriggerUsesProviderReported | 0 |

原样绿（未突变基线）＝`m0-baseline rc=0`，同一把尺同一选择器 **12 枚全 PASS**。

**红行逐字（M1）**：
```
compress_state_358_test.go:81: State = "skipped", want "attempted": the history is over the 384-token budget and nothing was foldable, which is an ATTEMPT that folded nothing - not a pass that never triggered
compress_state_358_test.go:92: State collapsed onto skipped, so '该压但压不动' is again indistinguishable from '没触发'
compress_state_358_test.go:159: "agent: history compress skipped" records = 0, want exactly 1 for the over-budget/no-fold shape (all: [trace-capture-ruler-control])
compress_trigger_358_test.go:91: Compression.State = "skipped", want "attempted": the provider's 150 cleared the budget over a history the estimate calls small, so the pass tried and had nothing foldable
```
**红行逐字（M3，触发基数那跳）**：
```
compress_state_358_test.go:270: Source = "estimate", want "provider_usage" when a provider count is available
compress_state_358_test.go:273: Tokens = 32, want the provider's 768 (NOT the local estimate 32)
compress_trigger_358_test.go:73: Compression.TriggerSource = "", want "provider_usage" on a run whose provider reported 330 input tokens
compress_trigger_358_test.go:196: second turn basis = "estimate", want provider_usage once a count exists
```
**红行逐字（M4/M5，循环接线那跳）**：`compress_trigger_358_test.go:80: Compression.TriggerTokens = 0, want the provider's count`

**sha256 还原证明**（逐字，M0 快照 与 五次来回之后）：
```
M0 快照  0bc4c7c4452d0ce26b07e4b09582cab6a3055008a2473c9d5d4a8a948f6758c7 *internal/agent/compress.go
M0 快照  e14f73b1458c411511523ed19e91f7f19a79ce0c863d5abd6358605fc476df72 *internal/agent/loop.go
还原后   0bc4c7c4452d0ce26b07e4b09582cab6a3055008a2473c9d5d4a8a948f6758c7 *internal/agent/compress.go
还原后   e14f73b1458c411511523ed19e91f7f19a79ce0c863d5abd6358605fc476df72 *internal/agent/loop.go
```
这两枚**本票要改**，故"还原"是回**本票改后快照**（上）；它们与 `HEAD:` blob **必然不等**（差即本票改动，sec 3 已逐行）：
`compress.go` worktree `66393358...f382` vs HEAD `dc345599...408c`；`loop.go` worktree `31dec50a...9c16` vs HEAD `ae8c25c7...d432`。
**"sha256 回 HEAD blob"实测那一枚＝`cmd/wisp/run.go`**（本票一度改过、后整枚退回）：`git hash-object` = `4865cfde6ef39bf30ad33a2c9a2636b82cb0d3ae` = `git rev-parse HEAD:cmd/wisp/run.go`，**逐字相等**。

## 6 门禁 rc（自己现跑，逐条分开报）

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| build | `GOFLAGS= go build ./...` | **0** | 无输出 |
| vet | `go vet ./internal/agent/ ./cmd/wisp/` | **0** | 无输出 |
| agent | `go test -count=1 -v ./internal/agent/` | **0** | 顶层 **RUN 102** / **PASS 102** / **FAIL 0** / **SKIP 0**（四数各自独立、不相减）加子测试 PASS 29；RUN 与终态名册**两向差皆空** |
| cmd/wisp | `go test -count=1 -timeout 30m -v ./cmd/wisp/` | **1** | 顶层 RUN **210**（共 336）/ **PASS 155** / `^--- FAIL` **0** / SKIP 0，**整进程死**，死点 TestAC14GoSideEvalPushReachesThePage |
| d22scan | `sh scripts/d22scan.sh` | **0** | `clean`；分母 `bans#1-5 internal/=244 cmd/=42`、`ban#6 frontend/=116`、`ban#7 internal/tools/=25`、`ban#8 design/=39 frontend/=116 internal/=610 cmd/=144` |
| path length | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`；`tracked=17507 over-budget=60 covered=60 not-in-roster=0` |
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe -l <4 件>` | **0** | 列表**空** |

gofumpt 用**绝对路径**（裸名 rc=127，不据此判"本机没有"）；不用 `CGO_ENABLED=0`；`cmd/wisp` 整包**抬 `-timeout 30m`**。

## 7 cmd/wisp rc=1 为什么不是本票 与 新增红名册

**记我自己一次返工**：我一度改 `run.go` 的 `compressionLine` 让 attempted 也出一行，撞票 355 钉子
`TestCompressionTraceSilentWhenNothingCompressed`（`:171-173` 要求形(3) **静默**），红句逐字：
`compression_trace_355_test.go:172: a pass that cleared the threshold but folded nothing printed "wisp run: 压缩痕：该压但压不动（…）" (ok=true)`。
我无权改 355 的钉子、也不能为变绿放宽它，故**整枚退回 `run.go`**（sha256 已证回 HEAD blob，sec 5）。
**痕的出口＝已有出口**：`installLogSink` 把 slog 默认改成 **JSONL 落盘 加 stderr 双写**（`logsink.go:144-161`），
故 `agent: history compress skipped` **自然**到得了盘也到得了终端，不需新造前端通道（正合票面"读数落点必须在已有出口上"）。

**rc=1 归因＝先于我存在的整进程死亡（票 314 在册）**。同尺对照＝把**只有本票动过的两枚文件**用 `-overlay`
逐字节换回 HEAD blob 跑同一发整包，**同样 rc=1、同样死在真窗族**，故不是本票造成：
```
current tree   FAIL github.com/CarlosShao/wisp/cmd/wisp 307.757s   rc=1
HEAD baseline  FAIL github.com/CarlosShao/wisp/cmd/wisp 331.376s   rc=1
```
**被死点吃掉的枚名（派单要求登记）**：当前树＝**TestAC14GoSideEvalPushReachesThePage**（有 `=== RUN` 无终态行）；
HEAD 基线＝**TestPanelThreadIsSTAAndExitsCleanly**。两发**死点不同名**，该族随会话漂（票 314 与台账 `:14768`），
故不拿一次读数指控另一次，也不加 `-skip` 变绿。

**新增红＝0 枚**（同一把尺、带 `-v`、**按名作差**，不相减枚数）：
```
HEAD pass=155  current pass=155
NEW RED (passed on HEAD, not on current tree):   (end)   <- 空
RECOVERED:                                       (end)   <- 空
```
**红名册（逐名）**：`TestAC14GoSideEvalPushReachesThePage`（HEAD 基线明确 `--- FAIL`，当前树被死点吃掉）＝台账
`:12027`/`:14768`/`:15340` 已具名 WebView2 真窗族，**与本票无关**，**非新增**。
本票相关的三枚 355 钉子当前树**全 PASS**（逐字）：
```
--- PASS: TestCompressionTraceThreeFlavoursReachTheClosingLine (0.00s)
--- PASS: TestCompressionTraceSilentWhenNothingCompressed (0.00s)
--- PASS: TestCompressionTraceReachesStdout (0.00s)
```
**补充读数（不当门禁、不当放宽）**：整包 `-skip` 掉台账在册 9 枚真窗族后**仍 rc=1**（死点前移
`TestBookedSnapshotReachesHostBitForBit329`，RUN 207 / PASS 151 / FAIL 0 / SKIP 1），印证"死点会漂"；
整包那 **126 枚从未跑到**的用例本票无法用整包这一发判色，**属票 314 射程，不是本票的伤**。

## 8 未定义即停 与 越界自查

乙案**没有做不到的地方**：三态在 `CompressionReport` 落得下，两处出口（`Result.Compression` 与 slog 痕）都接得住；
`cmd/wisp` 那一格因票 355 钉子归零改动（sec 7），**不因此降级成甲案或丙案**。越界 **0 枚**，`git show --name-only` 逐枚见提交正文；写面外零字节。
