# 84 — L2 审批等待"是不是永挂"：**已证伪，有 300 秒上界**（票面前提按实改小；副产品另开票 87）

**Status:** **done-as-refutation**（编排者验收 2026-09-21 16:0x）—— **本票的前提不成立，这不是实现缺陷**：两侧实测等待都有上界（Windows `5m0.0005138s` / Linux `5m0.017199699s`），正是 C18 的"300s 一律判拒绝"（`PLAN.md:1368`、`:3204`、`SPEC-06:104-106/151`），而 `PLAN.md:3143` **早就明文驳回过**"L2 设成无限等待"的提案。8 条"无对应待审批项"的答复路由在空门上 **0s / 700ns–1.9µs** 返回 `ErrUnknownCorrelation`（不 panic、不重试）⇒ 答复侧是诚实的。AC#3 **故意不勾**（契约要求的正是现状，生产码零改动）。本票的净价值：① 一个**能抓到"上界消失"的回归钉**（`84-ac1-bounded-wait.md` + 默认 SKIP 的计量用例）；② 把我 A59② 那句"无上界等待"**证伪并收回**；③ 逼出票 87（人想提前拒做不到，只能等满 300s）。
AC#1/AC#2/AC#4/AC#5/AC#6 已勾，**AC#3 故意不勾**（没有需要修的东西，勾它等于宣称做过一次不存在的修复）。
裁决表：`docs/evidence/s1/84-ac1-bounded-wait.md`。待编排者按 §6 拍板（关票 / 另开 L2 否决路由那张票）。
**Type:** 可用性/正确性缺陷（可能是死等；生产链路会把工具调用挂住）
**Blocks:** nothing（目前只在变异态下被逼出来）· **Blocked by:** nothing
**Packages:** `internal/agent/approval/gate.go`（≈`:473`）与它的调用方 `internal/tools/bridge.go`（≈`:340`）。
**禁改**：`docs/PLAN.md`、`docs/specs/*.md`、`internal/risk/**`（冻结/他人领地）、`tools/d22scan/**`。

## 怎么被发现的（证据链要如实带着走）

`verify75` 在 **Linux（docker `golang:1.27`）** 上做票 75 的变异复现（把 `normalizeLocalUNC` 的守卫退回
"无条件折叠"）时，`internal/tools` 整包 **600s `panic: test timed out`**。它按**测试名**把 600s 拆开归因：

- `TestL1WriteGoesThroughTheRealBlockWindow` **真等满了 C18 的 300 秒窗口**（300.05s）；
- 紧接着 `TestVetoInsideTheWindowWritesNothing` 卡在
  `approval.(*Gate).PendingApproval` **4m48s**（288s），前一行日志是
  **`correlation_id 无对应待审批项`**；
- 300 + 288 ≈ 600 ⇒ 撞的是**测试总闹钟**，不是死循环。另有 3 条从未轮到跑、1 条被打断，
  它在报告里**单列为"未定"、没算进新增红**（这个处理是对的）。
- 完整读数：`docs/evidence/s1/75-linux-mutation-check.md`（含基线 `RUN 169 / PASS 101 / FAIL 8 / SKIP 3`
  与变异后 `+32 顶层红`、`comm -23` 两包皆空 ⇒ 没有旧红被洗绿）。

⚠ **两条边界，别把这张票读歪**：
1. 这是在**变异态**（票 75 的缺陷被人为装回去）下被逼出来的，**基线态 `internal/tools` 是 0 FAIL、9.6s 过**。
   ⇒ 本票要判的是"**没有待审批项时该等多久**"这件事**是否与路径形状无关**。
2. 观察发生在 Linux。**Windows 上是否同样阻塞尚未复现** ⇒ AC#1 必须两侧都量。

## 为什么要单独一张票（不塞回票 75）

"等一个不会来的决定"是**审批交接**的语义问题：如果 `correlation_id` 在门里根本没有对应条目，
那**没有任何人会被叫起来点按钮**，等待就没有上界 ⇒ 生产里等于**一次工具调用把 agent 循环挂死**。
把它留在"让 Linux 变绿"的票里，下一个代理的最短路径会是"把测试的等待时间调短"——那是**把洞藏起来**。

## AC（1:1，裁决表 `docs/evidence/s1/84-*.md`）

- [x] **AC#1** **两侧各量一次"无对应待审批项"的真实等待**：写一条用例，直接对一个**空门**调
  `PendingApproval`（或走 `bridge` 那条真实路径），用**带时间戳的失败**记录它到底
  (a) 立刻返回错误 / (b) 等一个有限窗口后返回 / (c) **无限期阻塞**。
  Windows 与 Linux（docker `golang:1.27`，挂 `git archive` 快照，**不挂工作树**）各一份读数和真实 exit code。
- [x] **AC#2** **定性 + 定契约边界**：查 C18/C19 与 `SPEC-06`/`SPEC-04` 里"待审批窗口"到底约定了什么
  （引用要**自己打开文件核实行号**——本仓已有两起引用腐坏，A30 与 A51⑥）。
  若契约要求"无对应项 ⇒ 立即失败"，这就是**实现 bug**，直接修；
  若契约**没说**，**停下来上报编排者**（不要自己定新语义）。
- [ ] **AC#3** 修法只许走"**快速失败**"这一侧：找不到待审批项 ⇒ 返回一个**具名错误**（不 panic、不静默重试），
  调用方（桥接层）把它变成一条**给模型看得懂的拒绝**。**禁止**用"把超时调短"、`t.Skip`、
  或给测试单独塞小窗口来交差。
- [x] **AC#4** 变异双向：(i) 把快速失败退回原样 ⇒ 用例红；(ii) 把"拒绝"改成"默默放行" ⇒
  必须有用例红（这条是**安全侧**判据：审批门失效只能朝"更保守"方向坏）。
  锚点=承载行为的那一行，同一条 `&&` 链里 grep 自证落地，还原后 `diff -q`/`git diff --quiet` 证干净；
  **编译失败不算行为变异**（本仓为此重跑过）。
- [x] **AC#5** 顺带把 `TestL1WriteGoesThroughTheRealBlockWindow` 那条**真等 300 秒**的用例量一下：
  如果 C18 的窗口时长是可配置的，它应当用短窗口跑同一个语义；**如果不该改**（就是要测真实 300s），
  在票面写一句"这条是有意慢"并说明它该归哪个 CI job。**不许**为了 CI 快就把它删掉或调小真实值。
- [x] **AC#6** 门禁（只跑自己碰的包）：`gofmt -l` 空、`gofumpt -l <files>` 空、`go vet <pkgs>` rc=0、
  `GOOS=linux go vet <pkgs>` rc=0（⚠ 不要用 `GOOS=linux go vet ./...` 当判据，那条命令在 Windows 主机上
  因 CGO=0 排除 sherpa 预编译包而**永远 rc=1**，与你的改动无关 —— **A54③**）、
  `go test -count=2 <pkgs>` rc=0，`--- SKIP`/`--- FAIL` 逐条点名。

## Rules（本仓固定）

15 次工具调用内交第一枚 checkpoint commit；每次 commit 同步票面 Status + 勾框 + 末行 `next=`；
`git commit -q -F - -- <显式路径>` + **带引号** heredoc `<<'MSGEOF'`；禁 `git add -A`/`.`；
commit 前核对 `git diff --cached --name-only`；禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`（共树 A34）；
**不 push**；不在仓内建 worktree（A38④）；票面 Progress log append-only，**要改的那行先读再替换**；
四种假绿逐跑点名；数字难看就如实报，**不许调阈值/不许重测到运气好的那次**。

## Progress log（append-only）

- 2026-09-21 **认领 + 代码读数（checkpoint 1）**：读完 `internal/agent/approval/gate.go`（`PendingApproval`
  在 `:432`，卡点 `:473`）、`internal/agent/approval/queue.go`、`internal/tools/bridge.go:340` 与
  `docs/evidence/s1/75-linux-mutation-check.md` §4.3。**定性中间结论**：`PendingApproval` 不是"查一个已存在的
  待审批项"，它自己 `q.push()` 登记再等；等的上界是 `g.clock.After(g.q.Timeout())`，
  `Queue.Timeout()` 默认 `DefaultApprovalTimeout = 300 * time.Second`（`queue.go:92`，`NewQueue` 里
  `timeout <= 0` 也会被兜回默认值 ⇒ 不存在"0 = 无限"这条路）。`ErrUnknownCorrelation`
  （`internal/agent/approval/ui.go:124`）是**答复侧**（`Veto`/`q.allow`/`q.reject`）的具名错误，
  这三条路径实测是**立即返回**的。⇒ 倾向 **(b) 有界等待**，票面"无上限挂死"的说法待 AC#1 两侧读数裁决后**照实改小**。
  `next=` 先补 AC#1 的两侧读数（Windows 本机 + docker `golang:1.27` 的 `git archive` 仓外快照）。
- 2026-09-21 **AC#1 用例落地 + Windows 侧读数（checkpoint 2）**：新文件
  `internal/agent/approval/ticket84_no_owner_test.go`（5 条用例，全部在 `internal/agent/approval/` 内，
  未动 `bridge.go`、未动任何既有测试）。Windows 本机 `go test -count=1 -timeout 120s -v -run '<5 条命名>'`
  **rc=0**：`unanswered L2 (250ms deadline): returned after 250.9374ms (start=12:54:47, end=12:54:48)`；
  八条"无对应待审批项"答复路由（`Native.Allow`/`Native.Reject`/`Panel.Reject`/`DecideFromNative.{allow,reject}`/
  `DecideFromPanel.reject`/`Veto.unknown`/`Veto.empty`）全部 `in 0s` 且 `errors.Is(err, ErrUnknownCorrelation)`；
  界面不可达那条 `returned after 0s`（立即 fail-closed）。
  **1 条 `--- SKIP` 已点名**：`TestDefaultDeadlineWallClockMeasurement`（默认 300s 的墙钟计量，
  需 `WISP_84_MEASURE=1` 显式开启，本轮两侧各真跑一次，读数进 `docs/evidence/s1/84-ac1-bounded-wait.md`）。
  短窗口那条用例用的是 `Options.ApprovalTimeout`——生产同款旋钮（`cmd/wisp/run.go:259` 从
  `[risk].confirm_timeout_sec` 喂进来），不是给测试单独塞的小窗口。
  `next=` Linux（docker `golang:1.27` + `git archive <sha>` 仓外快照）跑同一批 + 300s 墙钟计量，再核 AC#2 契约行号。
- 2026-09-21 **AC#1 两侧读数完成（checkpoint 3）**：新用例
  `internal/agent/approval/ticket84_no_owner_test.go:TestDefaultDeadlineWallClockMeasurement`
  （`WISP_84_MEASURE=1` 显式开启，默认 SKIP）在**默认配置**（无 `ApprovalTimeout` 覆盖、`SystemClock`）下两侧各跑一次：
  **(a)** 默认门 UI 未接入 ⇒ `answer=reject why="审批界面不可达，已 fail-closed 拒绝执行"`，
  Windows `elapsed=0s`、Linux `elapsed=124.211µs` ⇒ **立刻拒绝**（这条不是失误，是契约 SPEC-06 §7「宿主不可达 fail-closed」的正面证据）；
  **(b)** 卡片已显示但无人答复（75 变异留下的形状）⇒ `answer=timeout`、
  Windows `elapsed=5m0.0005138s`（12:59:27→13:04:27）、Linux `elapsed=5m0.017199699s`（05:00:39Z→05:05:39Z），
  两侧 rc=0。**(c) 无限期阻塞未观测到**（外沿观察预算 400s）。
  答复侧 8 条"无对应待审批项"路由在空门上 Windows `in 0s` / Linux `700ns~1.9µs` 且
  `errors.Is(err, ErrUnknownCorrelation)`。⇒ **答案是 (b)：有上界，上界就是契约钉的 300s。**
  Linux 侧用 `git archive 9816bc1` 解到 `C:/Users/swq/wisp84-snap` 的**仓外快照**（未挂工作树），
  容器 `golang:1.27`；`MSYS_NO_PATHCONV=1` 必需（否则 Git Bash 把 `-w /src` 折成宿主路径，docker 直接报错）。
  短窗口对照（`Options.ApprovalTimeout=250ms`，生产同款旋钮 `cmd/wisp/run.go:259`）：Windows `250.9374ms`。
- 2026-09-21 **AC#2 定契约 + 定性（checkpoint 4，逐处自己开行号）**：`docs/PLAN.md:1368`（C18 行：「超时 = 300s，一律判拒绝」
  +「超时前 30s 醒目提示」+「宿主不可达时 fail-closed」）、`docs/PLAN.md:3143`（§16.9 第 4 条：「L2 审批超时应设为**无限等待**」
  被明文 **驳回**，理由是「永挂并持有路径锁（C20）」）、`docs/PLAN.md:3204`（第四轮给 C18 补数值 300s）；
  `docs/specs/SPEC-06-security-gatekeeping.md:104/105/106/151`（300s 判拒绝 / 宿主不可达 fail-closed /
  回复按 correlationId 路由 / S7 验收项），`SPEC-06:18` 是 L1 的 2–3s 阻止窗口（与 300s 无关，
  `internal/agent/approval/queue.go:105-107` 钉 `MinL1Window/MaxL1Window`）。
  `docs/PLAN.md:1369` 的 **C19 是 `RiskAssessor`，全文不含窗口/超时时长语义**（本票与它只有一处关系：判定器异常 ⇒ fail-closed 升 L2 ⇒ 进 C18 的 300s 队列）；
  `SPEC-04` 只在 `:71/:73` 谈 KWS 否决词与 2–3s 窗口相容，**没有"待审批窗口"时长约定**。
  ⇒ **契约说了，而且说得很硬 ⇒ 走 AC#2 的"实现 bug"分支都不成立：实现与契约一致，本票不是缺陷**，
  无需上报"契约没说"。票面"没有上界 / 把 agent 循环挂死"的账按实改小：真正的账是
  **(b) 300s 有界 + 一律拒绝**，以及 §6.3 那条相邻缺口（L2 卡片上 `Veto` 找不到项 ⇒ 人想拒也拒不掉，只能等 300s）。
- 2026-09-21 **AC#4/AC#5/AC#6（checkpoint 5）**：
  **AC#4 双向变异**：(i) 锚 `gate.go:456` `deadline := g.clock.After(g.q.Timeout())` → `make(chan time.Time)`
  ⇒ 新用例红 `--- FAIL: TestAnUnansweredL2CardResolvesToRejectAtTheQueueDeadline (8.00s)`，红信息带时间戳
  `STILL BLOCKED after 5s (start=13:10:09, now=13:10:14) - 闸门没有上界`；整包跑时既有
  `TestR7SizedBatchThatNobodyAnsweredIsRejectedNotRun (5.00s)` 同红、`TestL2QueueAutoRejectsAt300sKeepsTaskAliveAndWarnsAt270s`
  挂到 `FAIL ... 300.053s`，panic 栈 = `PendingApproval @ gate.go:473` ← `bridge.go:340`（与票 75 那条 288s 同形状）。
  (ii) 锚 `gate.go:469` 的 fail-closed 拒绝 → `tools.AnswerAllow` ⇒ 整包 `RUN 38 / 顶层 PASS 21 / 子 PASS 13 / FAIL 2 / SKIP 1`，
  两条红名 `TestHostUnreachableAndFullQueueFailClosed`（既有）+ `TestAnUnreachablePromptSurfaceIsRefusedImmediatelyNotWaitedOut`（新增）。
  两次还原后 `git diff --quiet -- internal/agent/approval/gate.go` 均成立（`CLEAN-i` / `CLEAN-ii`）。
  **AC#5**：`TestL1WriteGoesThroughTheRealBlockWindow` 基线 Windows 实测 **3.02s**（不是 300s），
  `TestVetoInsideTheWindowWritesNothing` 0.00s，`ok internal/tools 3.067s` rc=0；它测 L1 的 2–3s 窗口，
  而窗口被 `Min/MaxL1Window` 钉在 2–3s ⇒ **不存在"换短窗口跑同一语义"这条路**，那条 300.05s 是变异态把该写操作
  升成 L2 之后进了 300s 队列（**票 75 的缺陷在计时，不是这条用例在计时**）。它归 `test-core`
  （`.github/workflows/ci.yml:126-133`「Portable package tests」含 `./internal/tools/...`）。
  真正"有意慢"的是本票新增的 300s 墙钟计量用例（默认 SKIP，不进任何现有 job，是否给 nightly 待裁决）。
  **AC#6 两侧同数**：`gofmt -l` 空 rc=0、`gofumpt -l <新文件>` 空 rc=0、`go vet ./internal/agent/approval/` rc=0、
  `GOOS=linux go vet ./internal/agent/approval/`（按包作用域，A54③）rc=0、
  `go test -count=2 -v` 两侧均 rc=0 且 `=== RUN 76 / 顶层 PASS 46 / 子 PASS 28 / FAIL 0 / SKIP 2`
  （= `-count=1` 的 38/23/14/0/1 整 2 倍）；2 条 `--- SKIP` 全是 env 门的墙钟计量用例，已在 §AC#1 真跑两侧。
  **AC#3 不勾**：契约要求的就是现状 ⇒ 没有任何生产代码改动（`internal/tools/bridge.go` 一行未动，
  `internal/agent/approval/` 生产码零改动，只有新增测试文件）。
  `next=` 等编排者按裁决表 §6 拍板：(1) 本票按"前提不成立 ⇒ 非缺陷 + 回归钉已落"关票还是改写票面；
  (2) 300s 墙钟计量要不要独立 job；(3) 是否为"L2 卡片上 `Veto` 找不到项 ⇒ 人拒不掉只能等 300s"另开票。
