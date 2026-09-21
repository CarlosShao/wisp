# 84 — 审批门"等一个永远不会来的决定"：`PendingApproval` 无对应待审批项时**阻塞到超时**而不是快速失败

**Status:** open（**排队**：写码代理槽位已满，交回一张再派一张）
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

- [ ] **AC#1** **两侧各量一次"无对应待审批项"的真实等待**：写一条用例，直接对一个**空门**调
  `PendingApproval`（或走 `bridge` 那条真实路径），用**带时间戳的失败**记录它到底
  (a) 立刻返回错误 / (b) 等一个有限窗口后返回 / (c) **无限期阻塞**。
  Windows 与 Linux（docker `golang:1.27`，挂 `git archive` 快照，**不挂工作树**）各一份读数和真实 exit code。
- [ ] **AC#2** **定性 + 定契约边界**：查 C18/C19 与 `SPEC-06`/`SPEC-04` 里"待审批窗口"到底约定了什么
  （引用要**自己打开文件核实行号**——本仓已有两起引用腐坏，A30 与 A51⑥）。
  若契约要求"无对应项 ⇒ 立即失败"，这就是**实现 bug**，直接修；
  若契约**没说**，**停下来上报编排者**（不要自己定新语义）。
- [ ] **AC#3** 修法只许走"**快速失败**"这一侧：找不到待审批项 ⇒ 返回一个**具名错误**（不 panic、不静默重试），
  调用方（桥接层）把它变成一条**给模型看得懂的拒绝**。**禁止**用"把超时调短"、`t.Skip`、
  或给测试单独塞小窗口来交差。
- [ ] **AC#4** 变异双向：(i) 把快速失败退回原样 ⇒ 用例红；(ii) 把"拒绝"改成"默默放行" ⇒
  必须有用例红（这条是**安全侧**判据：审批门失效只能朝"更保守"方向坏）。
  锚点=承载行为的那一行，同一条 `&&` 链里 grep 自证落地，还原后 `diff -q`/`git diff --quiet` 证干净；
  **编译失败不算行为变异**（本仓为此重跑过）。
- [ ] **AC#5** 顺带把 `TestL1WriteGoesThroughTheRealBlockWindow` 那条**真等 300 秒**的用例量一下：
  如果 C18 的窗口时长是可配置的，它应当用短窗口跑同一个语义；**如果不该改**（就是要测真实 300s），
  在票面写一句"这条是有意慢"并说明它该归哪个 CI job。**不许**为了 CI 快就把它删掉或调小真实值。
- [ ] **AC#6** 门禁（只跑自己碰的包）：`gofmt -l` 空、`gofumpt -l <files>` 空、`go vet <pkgs>` rc=0、
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

（空）
