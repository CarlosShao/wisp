# 票 29 `29-v1` 复核裁决（第二轮）— AC#1 PASS / AC#2 PASS

- 身份：非实现者复核腿（v1）。起手锚 `1637f74598e2920a999acac4ed9ae1ada511102d`。
- 被裁对象：`f92e4b177`（AC#1，真 L1 抽取链 + 宿主 Warm 窗口装配）与
  `bbe7e2bee`（AC#2，`newGoldenExtractor` 换真链 + 10 枚金标 SSE）。
- 本文件覆盖 2026-10-28 的上一版 29-v1 裁决（提交 `e99f6932d`，判 AC#1/AC#2 两格
  “待装配/待 golden 链”）。上一版判词在 git 历史里，未删除；本轮只针对它点名的
  两枚残余现核。
- 结论：**AC#1 PASS，AC#2 PASS**；反向齿咬得动；金标 JSON 未被这两提交改过；
  范围只含预期文件。残余见 §6（均不阻断翻格）。

## 0. 对象与提交读数

| 提交 | 日期 | 文件/行 | 内容 |
|---|---|---|---|
| `f92e4b177b0c36a1928c18d1c962ff90160a24a3` | 2026-11-04 20:41:37 +08 | 7 files, +1654 | `internal/memory/l1_extract.go`、`l1_warm.go`、两个 AC#1 测试、`cmd/wisp` 三处装配 |
| `bbe7e2beeed2101f8ed237255dd0df9b722e9c65d` | 2026-11-04 21:02:01 +08 | 14 files, +631/-258 | 删规则桩、`newGoldenExtractor` 换真链、`golden_sse_replay_test.go`、10 枚 `testdata/golden_sse/*.sse` |
| `1637f7459...`（HEAD） | 2026-11-04 21:04:18 +08 | 1 file, +6 | 仅 `docs/reports/pending-and-issues.md`（A1525），非实现面 |

`git diff --name-only f92e4b177^ HEAD` 共 21 条：20 条实现/测试路径 + 1 条台账；
`internal/memory/testdata/golden_extract.json` **不在其中**。

## 1. 本轮自跑读数（精确 rc）

| # | 命令 | rc | 读数 |
|---|---|---|---|
| 1 | `go build ./internal/memory ./cmd/wisp` | **0** | 无输出 |
| 2 | `go test ./internal/memory -run TestTicket29 -count=1 -timeout 220s -v` | **0** | `ok github.com/CarlosShao/wisp/internal/memory 4.355s`；12 个顶层测试 + 30 个子测试全 PASS |
| 3 | `go test ./internal/memory -count=1 -timeout 320s` | **0** | `ok github.com/CarlosShao/wisp/internal/memory 30.466s` |
| 4 | `gofumpt.exe -l <两提交触碰的 10 个 .go>` | **0** | **空输出**（无未格式化文件） |
| 5 | `sh scripts/d22scan.sh` | **0** | 正控 35 PASS；实仓 `clean - no D22 ban violations`（408 production Go files） |
| 6 | `go test ./internal/session -run TestTicket29 -count=1` | **0** | `TestTicket29WarmHookRunsL1Extraction` PASS（0.13s），补 AC#1 跨包机制 |

读数留档：`D:\tmp\wisp\20261104\29v2\{test-ticket29,test-full,reverse-tooth,test-session-29}.log`（仓外）。

AC#1 实跑测试名（`-run TestTicket29` 命中）：
`TestTicket29_AC1_WarmWindowLatencyIsolation`、
`TestTicket29_AC1_OneCallPerTaskBatchCommit`、
`TestTicket29_AC1_ParserMatchesGoldenContract`（含 10 个 case 子测试）；
跨包另有 `TestTicket29WarmHookRunsL1Extraction`。

AC#2 实跑测试名：
`TestTicket29_AC2_GoldenSetContract`、
`TestTicket29_AC2_GoldenExtractionMatches`（10 个 case 子测试）、
`TestTicket29_AC2_GoldenBytesThroughPipeline`（10 个 case 子测试）、
`TestTicket29_AC2_ComparatorNegativeControls`。

## 2. AC#1 PASS：抽取真在 post-response 的 Warm 窗口里

**判据不是“有钩子接口”，是主机端真的赋了值，且响应路径上零 C5 调用。**

装配链（逐行现核）：

| 位置 | 读数 |
|---|---|
| `cmd/wisp/resident_voice_engines_windows.go:499` | `WarmWindow: r.warmWindow,`（`session.Hooks{}` 内的真实赋值） |
| 同文件 `:530` | `func (r *voiceResidency) warmWindow()` 钩子体；`r.l1Warm` 非 nil 才 `fire()` |
| 同文件 `:552` | `attachMemoryL1(run *agentRuntime)`：用**同一个** `run.store` 与 `run.L1Provider()` 组 `NewL1Extractor` → `NewL1Pipeline` → `NewL1WarmWindow` |
| 同文件 `:588-589` | `run.memoryL1 = warm; rv.residency.attachL1Warm(warm.Fire)`（recorder 与 hook 是同一对象） |
| `cmd/wisp/resident_windows.go:603` | `rv.attachMemoryL1(src.run)`（resident 装配根真正调用） |
| `internal/session/scope.go:477,483` | `warmWork := s.hooks.WarmWindow` … `warmWork()`：只有 `EnterWarm` 会触发钩子 |

post-response 顺序（`cmd/wisp/run.go`）：

- `:1699` `res := bg.Wait()` —— 任务收口、答案已出；
- `:1711` `rt.memoryL1.Record(memoryConversationFromHistory(res.TaskID, loop.History()))`
  —— 在 `bg.Wait()` 之后、presenter/backfill 之前；`Record` 只拷贝转写。

`internal/memory/l1_warm.go:158` 的 `Record` 本体无 `Stream`、无 DB 写；`:199` 的
`Fire` 只 `reg.Spawn` 后立即返回。测试把这两半钉死：

- `l1_warm_29_test.go:171` `TestTicket29_AC1_WarmWindowLatencyIsolation`：
  40 次 `Record` 后断言 `prov.callCount() == 0`；provider 被 block 时
  `Fire()` 必须在 20ms 内返回（实测阈值内），且 `Started=1, Succeeded=0`；
  `close(release)` 后整批一次落库；首 token 中位数 `wired <= base + 5ms`。
- `l1_warm_29_test.go:321` `TestTicket29_AC1_OneCallPerTaskBatchCommit`：
  一次任务三槽位 = **1 次 C5 调用**、一次事务；坏槽位整批拒绝（`pref.tone` 不落）。
- `internal/session/warm_l1_29_test.go:117` 用 `WarmWindow: warm.Fire` 复现宿主
  赋值；`EnterWarm` 前 0 调用、后 1 调用，且 `EnterWarm` 不等抽取完成。

上一版点名“`Hooks{}` 未赋 WarmWindow”的缺口已由 `:499` 实值闭合。**AC#1 PASS。**

## 3. AC#2 PASS：规则桩真删，金标断言真走生产链

**桩确实不在：** `git show bbe7e2bee -- internal/memory/golden_extract_test.go`
显示删除 `goldenStubExtractor`、`regexp` 导入、`goldenStartRe/goldenEndRe/goldenDateRe`
及 `goldenForgetTarget/goldenReplyLanguage/goldenTone/goldenWorkHours/goldenEditor/
goldenFamily/goldenContainsAny`。现仓 `grep -rniE "stub" internal/memory --include=*.go`
只剩 5 处**历史注释**，无符号/类型/函数。

**施工点返回真链：**

- `internal/memory/golden_extract_test.go:43-46`：
  `func newGoldenExtractor(...) { return newGoldenSSEChain(t, set) }`
- 同文件 `:91`（GoldenExtractionMatches）与 `:162`（GoldenBytesThroughPipeline）
  都从这一施工点取 extractor。
- `internal/memory/golden_sse_replay_test.go:67-81` `newGoldenSSEChain`：
  `openaichat.New(...)`（真适配器）+ `NewL1Extractor(provider, ...)`（真链），
  BaseURL 指向 `httptest.NewServer` 回放嵌入的 SSE 字节。
- 回放器在 `:152-158` 强制每个 `.sse` **恰好 1 条 200 response**；`:21-28` 的注释
  与 `:189-193` 的实现拒绝 system 前缀漂移；一个 case 第二次调用会打成 599。

10 枚 `case-01..case-10` 文件均在，`golden_extract.json` 仍是 10 个 case
（`grep -c '"id": "case-'` = 10）。AC#2 的 `GoldenExtractionMatches` 每例断言：
1 次 C5 调用、`CompareExtraction` 全等、经真 DAO/真表读回 projection；
`GoldenBytesThroughPipeline` 再把同一批字节过 `L1Pipeline.RunTask`（一调用一事务）。

**金标 JSON 未被这两提交编辑：**

- `git diff HEAD~2 HEAD -- internal/memory/testdata/golden_extract.json` → **0 行**；
- 更强的 `git diff --name-only f92e4b177^ HEAD` 也不含该文件。

**AC#2 PASS。**

## 4. 反向齿（先咬红，再原样还原）

对象：`internal/memory/testdata/golden_sse/case-01-pref-reply-language.sse`
（唯一一处 `zh` → `zg`，恰改 1 字节）。

| 步 | 读数 |
|---|---|
| 改前 sha256 | `96e441d6fed889eabc4a2fff0341b166f5d7c931a5120204e6fc7fb7073f583a` |
| 改后 sha256 | `956ce00597bd5ada295e3384347729b23dee4e172a943667fe65b46b02099c69` |
| `git diff --stat` | 1 file changed, 1 insertion(+), 1 deletion(-) |
| `go test ./internal/memory -run TestTicket29_AC2 -count=1` | **rc=1（RED）** |
| 红行 | `golden_extract_test.go:115: golden mismatch: memory: golden extraction: slot "pref.language" content "zg", want "zh"`；`GoldenBytesThroughPipeline` 在 `:182` 同一行转红；包级 `FAIL ... 3.057s` |
| 还原 | `git checkout -- <该 sse>`；sha256 回到 `96e441d6...`（与改前逐字节相同） |
| `git status --porcelain -- <该 sse>` | **空**；`git diff --stat -- <该 sse>` 空 |
| 还原后复跑 AC2 | `ok ... 2.851s`，rc=0 |

反向齿有效：改一个期望值字节，真链回放即红并点名槽位与两值；还原后即绿。
（实现腿提交信息写“case-01 … naming pref.tone”，与本轮改 `zh` 得到 `pref.language`
不矛盾，只是各自挑的字节不同。）

## 5. 范围与禁面

- 两提交只碰预期文件：`cmd/wisp` 3 处装配、`internal/memory` 5 个 .go + 10 枚新增
  `.sse`、`internal/session` 1 个测试；HEAD 追加的 1 条是台账 `pending-and-issues.md`。
- 未碰：`internal/memory/testdata/golden_extract.json`、`thresholds.go`、任何 SLO
  阈值/golden 期望、`docs/PLAN.md`、`docs/specs/**`、票面 `.scratch/wisp/issues/29-*.md`。
  （`git diff --name-only f92e4b177^ HEAD | grep -Ei "thresholds|golden_extract.json|PLAN.md|SPEC-|\.scratch"` = 空。）
- d22scan rc=0 clean，无裸 goroutine、无面板侧 L2 allow、无 emoji 违规。
- 本轮写面：仅本文件；临时读数在仓外 `D:\tmp\wisp\20261104\29v2\`；未 push。

## 6. 残余（不阻断本轮翻格，归后续小修）

1. **`internal/memory/golden_extract.go:8-10` 注释已过期**：仍写
   “the L1 extraction pipeline is NOT wired yet (29-v1: the host never assigns
   WarmWindow …)”。f92e4b177 之后主机已赋值 `WarmWindow`，该句与 `:17-37` 的新说明
   自相矛盾。纯注释、不影响断言，但会误导下一读者；建议下一腿顺手改成历史时态。
2. **金标 SSE 是“按评审过的 fixture 钉死的字节”，不是实时端点抓包**
   （`l1_extract.go:17-22` 已明说）。故 AC#2 证明的是“真适配器 + 真解析/归一 +
   真事务链能复现钉死契约”，**不证明 live 模型当天会这样回答**。这是既定 harness
   形状（与 `internal/agent/testdata/golden` 先例一致），不是缺项。
3. **AC#1 的“statistical assertion”是 40 样本中位数 + 5ms 容差**，不是形式化显著性
   检验；真正咬住的是零 C5 调用、`Fire` 非阻塞、宿主赋值、`EnterWarm` 唯一触发这四条
   机制断言。若要更硬的统计口径，可在后续把容差/样本量写进测试注释或加分布断言。
4. 上一版 29-v1 裁决文件被本文件覆盖；旧判词保留在提交 `e99f6932d`，可随时回看。

## 7. 判语

| 格 | 判 | 一句话凭据 |
|---|---|---|
| **AC#1 Latency-isolation** | **PASS** | `resident_voice_engines_windows.go:499` 真赋 `WarmWindow`；`run.go:1711` Record 在 `bg.Wait()`(:1699) 之后；Record 零 C5、Fire 异步 |
| **AC#2 Golden extraction** | **PASS** | `goldenStubExtractor` 已删；`newGoldenExtractor` → 真 `openaichat` + `L1Extractor` 回放 10 枚 SSE；fixture JSON 0 diff；反向齿红→还原绿 |
| AC#3/#4/#5/#6 | 沿用上一版 PASS | 本轮全包复跑 rc=0，未见回归 |

**最终判：AC#1 PASS、AC#2 PASS；两枚 29-v1 具名残余闭合。**
