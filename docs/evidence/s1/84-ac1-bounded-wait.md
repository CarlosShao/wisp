# 票 84 裁决表 — AC#1 两侧读数：`PendingApproval` 的等待**有界**（(b)），不是无限期阻塞（(c)）

写码代理，2026-09-21。用例：`internal/agent/approval/ticket84_no_owner_test.go`（commit `e5e5eb7` + `9816bc1`）。
本票票面标题的说法是"阻塞到超时而不是快速失败"。**实测结论：它有超时上界，且上界正是契约钉的 300s；
票面"等待没有上界 / 一次工具调用把 agent 循环挂死"的账要改小。**

## 1. AC#1：三种可能里到底是哪一种

`PendingApproval`（`internal/agent/approval/gate.go:432`）**不是**"查一个已存在的待审批项"——它自己
`q.push()` 登记一项再等（`:437`），所以"空门调用"这个形状里**永远有且只有一个属于本次调用的待审批项**。
真正能出现"无对应待审批项"的是**答复侧**（`Veto` / `Native.Allow` / `Native.Reject` / `Panel.Reject` /
`DecideFrom*`），而那些路径实测**立即返回具名错误**。于是 AC#1 拆成两个真实形状分别量：

| 形状 | 平台 | 读数 | 判定 | 真实 exit code |
|---|---|---|---|---|
| **(a)** 默认门、UI 未接入（`Options{}` ⇒ `UIFuncs{}` 的 `Prompt` 返回 `errNoUI`） | Windows 主机 | `answer=reject why="审批界面不可达，已 fail-closed 拒绝执行" elapsed=0s`（start=end=`2026-09-21T12:59:27+08:00`） | **立刻返回拒绝**（fail-closed，比"返回错误"更强：不给执行） | 用例 `--- PASS`；整条命令 rc=0 |
| **(a)** 同上 | Linux（docker `golang:1.27`，`git archive 9816bc1` 解到 `C:/Users/swq/wisp84-snap`，**未挂工作树**） | `elapsed=124.211µs`，同一 why | 同上 | 同上，rc=0 |
| **(b)** 卡片已显示但**无人答复**（75 变异留下来的形状：答复侧的 key 落空） | Windows 主机 | `answer=timeout why="审批超时（300 秒未确认），C18 一律判拒绝，已自动拒绝" elapsed=5m0.0005138s`（12:59:27 → 13:04:27） | **等一个有限窗口后返回拒绝**，窗口 = `Queue.Timeout()` = 300s | `--- PASS`，`ok ... 300.045s`，`WINDOWS_MEASURE_RC=0` |
| **(b)** 同上 | Linux（同上快照） | `elapsed=5m0.017199699s`（05:00:39Z → 05:05:39Z），同一 answer/why | 同上 | `--- PASS`，`ok ... 300.033s`，`LINUX_MEASURE_RC=0` |
| **(c)** 无限期阻塞 | 两侧 | **未观测到**：外沿观察预算 400s，两侧都在 300.0s 处返回 | 不成立 | — |

短窗口对照（同一语义，`Options.ApprovalTimeout = 250ms`，`SystemClock`）：
Windows `returned after 250.9374ms`，Linux 同用例 PASS（整包 rc=0）。⇒ 上界**跟着配置的 deadline 走**，
不是测试里写死的等待。`Options.ApprovalTimeout` 是生产旋钮（`cmd/wisp/run.go:259` 从
`[risk].confirm_timeout_sec` 喂进来），**不是给测试单独塞的小窗口**。

"无对应待审批项"答复侧 8 条路由（`Native.Allow` / `Native.Reject` / `Panel.Reject` /
`DecideFromNative.reject` / `DecideFromNative.allow` / `DecideFromPanel.reject` / `Veto.unknown` /
`Veto.empty`）在空门上：Windows 全部 `in 0s`、Linux `700ns ~ 1.9µs`，且
`errors.Is(err, approval.ErrUnknownCorrelation)` 全部成立（`ui.go:124`）。**不 panic、不静默重试。**

## 2. AC#2：契约到底约定了什么（每处引用都自己打开核过行号）

| 引用 | 核实位置（本代理自己读到该行） | 原文要点 |
|---|---|---|
| **C18** | `docs/PLAN.md:1368` | 「**超时 = 300s，一律判拒绝**（第四轮补数值，原文无数值 → agent 会自己拍一个或写成无限等待，见 §16.9 第 4 条）；**超时前 30s 醒目提示**；拒绝后**任务 root ctx 不取消，可一键重放**；宿主不可达时 **fail-closed**」 |
| §16.9 第 4 条 | `docs/PLAN.md:3143` | 有人提议「L2 审批超时应设为**无限等待**」⇒ 裁决 **驳回**，理由正是「无限等待会让任务**永挂并持有路径锁（C20）**」 |
| 第四轮补数值表 | `docs/PLAN.md:3204` | C18 原文「超时一律判拒绝」**无数值** ⇒ 补 **300s** + 提前 30s 提示 + 拒绝后不取消 root ctx |
| **C19** | `docs/PLAN.md:1369` | C19 是 `RiskAssessor`（R1–R9 中心推断、判定器异常 ⇒ fail-closed 升 L2），**不含任何窗口/超时时长语义**；本票与 C19 的关系只有"异常判定会升 L2 ⇒ 会进 C18 的 300s 队列" |
| SPEC-06 §7 | `docs/specs/SPEC-06-security-gatekeeping.md:104` | 「**超时 300s 一律判拒绝**；超时前 30s 醒目提示；拒绝后任务 root ctx 不取消、可一键重放」 |
| SPEC-06 §7 | 同文件 `:105` | 「宿主不可达时 **fail-closed**（拒而非放）」——实测形状 (a) 正是它 |
| SPEC-06 §7 | 同文件 `:106` | 「回复按 correlationId 路由（点击不可能落到别的请求上）」——落空时的具名错误即 `ErrUnknownCorrelation` |
| SPEC-06 §7 | 同文件 `:151` | S7 验收：并发两任务同时 L2 ⇒ correlationId 路由正确、队头单显、**超时判拒绝**、重放可续 |
| SPEC-06 §2 | 同文件 `:18` | L1 是「执行前阻止窗口（**2–3s**）」，与 300s 无关（`MinL1Window/MaxL1Window` = 2s/3s，`internal/agent/approval/queue.go:105-107`） |
| SPEC-04 | `docs/specs/SPEC-04-voice-pipeline.md:71`、`:73` | 只约定 KWS 否决词与 2–3s 阻止窗口相容、否决通道清单；**SPEC-04 全文没有"待审批窗口/超时"的时长语义**（grep `300|超时|审批` 命中项均为 ASR/TTS/VAD 延迟预算） ⇒ 审批窗口的契约归属是 SPEC-06 §7 / PLAN C18，不在 SPEC-04 |

**定性结论：契约说了，而且说得非常硬** —— 300s 有限窗口 + 一律判拒绝 + 「无限等待」被明文驳回（`PLAN.md:3143`）。
现实现（`gate.go:456` 用 `g.clock.After(g.q.Timeout())` 武装上界；`queue.go:70-87` 的 `NewQueue` 把
`timeout <= 0` 兜回 `DefaultApprovalTimeout = 300s`，`queue.go:92`；答复侧一律 `ErrUnknownCorrelation` 立即返回）
**与契约一致 ⇒ 本票不是实现 bug，不需要按 AC#3 动刀，也无需上报"契约没说"**（AC#2 的两个分支里走的是前一个）。

## 3. AC#4：双向变异（证明这两条保证是承重的）

| 变异 | 锚点（承载行为那一行） | 结果 |
|---|---|---|
| (i) 把上界摘掉 ⇒ 退回"没有上界" | `internal/agent/approval/gate.go:456`：`deadline := g.clock.After(g.q.Timeout())` → `deadline := make(chan time.Time)`（编译通过，纯行为变异） | **红**：`--- FAIL: TestAnUnansweredL2CardResolvesToRejectAtTheQueueDeadline (8.00s)`，失败信息带时间戳：`STILL BLOCKED after 5s (start=13:10:09, now=13:10:14) - 闸门没有上界`。整包跑时另外两条既有用例同样被打回原形：`--- FAIL: TestR7SizedBatchThatNobodyAnsweredIsRejectedNotRun (5.00s)`，且 `TestL2QueueAutoRejectsAt300sKeepsTaskAliveAndWarnsAt270s` 挂到 300s 闹钟 `FAIL ... 300.053s`，panic 栈正是 `approval.(*Gate).PendingApproval` @ `gate.go:473` ← `internal/tools/bridge.go:340`——**与票 75 证据里那条 288s 等待同一形状** |
| (ii) 把 fail-closed 拒绝改成默默放行 | `gate.go:469`：`return tools.AnswerReject, "审批界面不可达，已 fail-closed 拒绝执行"` → `return tools.AnswerAllow, "界面不可达，放行（MUTANT-84ii）"` | **红**（安全侧判据咬人）：整包 `-count=1 -v` ⇒ `=== RUN 38 / --- PASS 21 / --- PASS(子) 13 / --- FAIL 2 / --- SKIP 1`，红名 `TestHostUnreachableAndFullQueueFailClosed`（既有）+ `TestAnUnreachablePromptSurfaceIsRefusedImmediatelyNotWaitedOut`（本票新增）；基线同命令是 `38 RUN / 23 顶层 PASS / 14 子 PASS / 0 FAIL / 1 SKIP` |
| 还原自证 | 两次变异后各自 `python` 反向替换 → `grep` 已无 `MUTANT-84*` | `git diff --quiet -- internal/agent/approval/gate.go` **成立**（两条都打了 `CLEAN-i` / `CLEAN-ii`），`git status --short internal/agent/approval/` 空 |

## 4. AC#5：`TestL1WriteGoesThroughTheRealBlockWindow` 不测 300s，它测的是 3s

- Windows 基线实测：`--- PASS: TestL1WriteGoesThroughTheRealBlockWindow (3.02s)`、
  `--- PASS: TestVetoInsideTheWindowWritesNothing (0.00s)`、`ok github.com/CarlosShao/wisp/internal/tools 3.067s`，rc=0。
- 它走的是 **L1 执行前阻止窗口**（SPEC-06 §2 `:18` 的 2–3s），`compose()` 用的是
  `approval.New(approval.Options{UI: ui})` ⇒ `Window` 落在 `DefaultL1Window=3s`，
  并被 `MinL1Window=2s / MaxL1Window=3s`（`queue.go:105-107`）钉死 ⇒ **窗口时长"可配"也只能在 2–3s 之间**，
  它 `time.Since(start) < 2*time.Second` 的断言正是在钉这个下界。所以**没有"用短窗口跑同一语义"这条路**：
  2s 已经是契约允许的最短。
- 票面/evidence 里那条 **300.05s** 是**票 75 变异态**的产物：路径折叠把该写操作从 L1 升成 L2 ⇒
  进了 C18 的 300s 队列 ⇒ 该用例等满 300s 后失败（`wiring_test.go:124` 打印的
  `审批超时（300 秒未确认）` 就是 L2 的 why，不是 L1 的窗口）。**基线态它 3.02s。**
- 归属：`internal/tools/...` 在 CI 的 `test-core` job ⇒ `.github/workflows/ci.yml:126-133`
  步骤 `Portable package tests (...)`。这条 3.02s 的用例**不是"有意慢"**，无需为它单开 job。
- 真正"有意慢"的是本票新增的 `TestDefaultDeadlineWallClockMeasurement`（300s 墙钟计量）：
  **默认 `--- SKIP`**，只在 `WISP_84_MEASURE=1` 时跑；不进任何现有 CI job（否则 `test-core` 白增 300s×平台数）。
  要不要给它一个 nightly job，属编排者裁决（见 §6）。

## 5. AC#6：门禁（只跑本票碰的包 `internal/agent/approval`）

| 命令 | Windows 主机 | Linux（docker `golang:1.27`，仓外快照 `9816bc1`） |
|---|---|---|
| `gofmt -l internal/agent/approval/` | 空，rc=0 | 空，`GOFMT_RC=0` |
| `$(go env GOPATH)/bin/gofumpt -l internal/agent/approval/ticket84_no_owner_test.go` | 空，rc=0 | 同（快照内 `gofmt` 已判空） |
| `go vet ./internal/agent/approval/` | rc=0 | `VET_RC=0` |
| `GOOS=linux go vet ./internal/agent/approval/`（**按包作用域**，未用 `./...`；A54③） | `linux_vet_rc=0` | —（原生即 linux） |
| `go test -count=2 -timeout 300s -v ./internal/agent/approval/` | rc=0：`=== RUN` **76**、顶层 `--- PASS` **46**、子 `--- PASS` **28**、`--- FAIL` **0**、`--- SKIP` **2** | rc=0：`=== RUN` **76**、`--- PASS` **46**、子 **28**、`--- FAIL` **0**、`--- SKIP` **2**；`ok ... 0.559s` |

N 倍核对：`-count=1` 基线是 `RUN 38 / 顶层 PASS 23 / 子 PASS 14 / SKIP 1 / FAIL 0` ⇒ `-count=2` 每一项都恰好 ×2（76/46/28/2/0）。
**四种假绿逐条点名**：① `--- SKIP` 有 2 条，全部是 `TestDefaultDeadlineWallClockMeasurement`（env 门，两侧真跑过，读数见 §1）；
② 未使用 `-run` 空匹配（本表是整包无 `-run`）；③ `-count=2` 的 `=== RUN` 已按 N 倍核对；④ 无步骤被静默跳过（两侧数字逐项相同）。

## 6. 需要编排者拍板

1. **本票账要改小**：票面标题/`Type`/「为什么要单独一张票」都建立在"等待没有上界"上，实测是 **(b) 有界 300s**
   且这正是 `PLAN.md:3143` 明文选择的语义。⇒ 建议按"**前提不成立 ⇒ 非缺陷，转为回归钉**"关票
   （回归钉已落：`ticket84_no_owner_test.go` 的 5 条用例，两侧 rc=0）。AC#3 框**未勾**：没有需要修的东西，
   勾它等于宣称做过一次不存在的修复。
2. 300s 墙钟计量用例要不要进某个 nightly/慢 job（默认 SKIP 的话，只有人工跑才拿到数）。
3. 一处**相邻的真实可用性缺口**（本票没动，因为它属审批交接的 UX 而非"无界等待"）：
   L2 卡片显示后，用户的否决若走 `Gate.Veto`（球/Esc/KWS 那条 L1 路由）**找不到 L2 项**，
   只会拿到 `ErrUnknownCorrelation` 并**让卡片继续等到 300s**——即"人想拒绝，但拒绝按钮不在 L2 队列上"。
   票 75 变异态观察到的 288s 就是这个形状的放大版。要不要为它开一张新票（把球/Esc 的否决接到 L2 队列头，
   或在 Veto 落空时把该 corr 的 L2 项按拒绝处理），请裁决。
4. 共用工作树的时序坑（本票遇到一次）：checkpoint 1 的票面编辑被编排者的 `2e20920` 顺手带走，
   我自己 `git commit` 返回 rc=1「nothing added to commit」。内容没丢，但"HEAD 变了"这条判据在共树下
   需要额外确认 commit 作者是本代理。

## 7. 逐字复跑命令

```
# Windows（工作树）
WISP_84_MEASURE=1 go test -count=1 -timeout 900s -run TestDefaultDeadlineWallClockMeasurement -v ./internal/agent/approval/

# Linux（仓外快照，不挂工作树）
git archive 9816bc1 | tar -x -C C:/Users/swq/wisp84-snap
cd /c/Users/swq && MSYS_NO_PATHCONV=1 docker run --rm \
  -v "C:/Users/swq/wisp84-snap:/src" -v wisp75-gomod:/go/pkg/mod -v wisp75-gobuild:/root/.cache/go-build \
  -w /src -e WISP_ENV=test -e GOFLAGS=-buildvcs=false -e WISP_84_MEASURE=1 golang:1.27 \
  bash -c 'go test -count=1 -timeout 900s -run TestDefaultDeadlineWallClockMeasurement -v ./internal/agent/approval/'
```

（`MSYS_NO_PATHCONV=1` 是必需的：Git Bash 会把 `-w /src` 折成 `D:/work/soft/Git/src` 之类的宿主路径，
docker 直接报 `the working directory ... is invalid`。）
