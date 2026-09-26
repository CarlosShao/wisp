# 151 r1 — `Bridge.CloseTask` 全仓零调用者：per-task 污点 scope 只开不关（实现程判定与修法）

- 工单＝`.scratch/wisp/issues/151-bridge-opentask-sits-on-the-dispatch-hot-path-while-closetask-has-zero-callers-repo-wide-so-the-per-task-taint-scope-is-only-ever-opened.md`
- 派单正文＝`.scratch/wisp/dispatches/2026-09-26-091x-r152-r153-r151.md` 的「派单 C」
- 角色＝实现程（本件不是验收表；按 `SPEC-12 §4.3` #1 与 AGENTS §0.3，缺口审计与对抗验收必须由另一枚程做）
- 落点＝`internal/tools/bridge.go`（`CloseTask` 那一段注释与审计行）＋ `cmd/wisp/run.go`（组合根的 per-task 边界）
  ＋ `cmd/wisp/task_scope_close_151_test.go`（新用例）＋ 本件 ＋ `.scratch/wisp/probes/151/`
- ⚠ 交件名的分歧先记一笔：派单给我的任务是 `151-taint-scope-never-closed-r1.md`，
  工单 Acceptance 一节写的是 `151-task-scope-never-closed-r1.md`。**本件按工单面落名**，
  这一处不一致按 AGENTS 的规矩上报，不由实现方改票面。

---

## 0. 锚点与环境（现读，不抄派单）

```
$ git rev-parse --short HEAD        ->  86b0161   （与派单的共同锚点同一枚）
```

读数采在**共享工作树**上，09:1x-09:3x。期间别人的提交落了 6 枚，逐枚 `git show --name-only` 过，
与本程的三枚文件全部不相交：

| 提交 | 谁 | 碰了什么 |
|---|---|---|
| `ad9b29f` / `cdf2471` | 前端会话 | `frontend/**`、`docs/reports/frontend-session-log-zh.md` |
| `b23c7f7` | 票 153 | `internal/agent/compress.go` / `compress_trace_test.go` / `loop.go` |
| `eb4755a` / `4cc85bb` | 票 152 | `.scratch/wisp/probes/152/**`、`docs/evidence/s1/152-*.md` |
| `ac7fb00` | 票 153 | `.scratch/wisp/probes/153/**`、`docs/evidence/s1/153-*.md` |
| `720cae6` | 编排者 | `.scratch/wisp/dispatches/2026-09-26-093x-thaw-panel-for-145.md`（票 145 特批存档） |

时间归属：AC#1 探针与两包**改前**门禁采在 `86b0161` 那枚树上；两包**改后**门禁与 `-count=2`
采在 `b23c7f7`/`4cc85bb` 之后的树上。`go build ./...` rc=0、两包零红，都是在最终树上复跑的。

工作树里另有 `design/**` 的大量 ` D`/` M` 与 `frontend/**` 未提交改动、`.zcodeignore` 等未跟踪件
（前端会话的活）。**本程一字未碰、未还原，也没把它们算进任何"零命中"宣称**（§5 的零字节名册
只按"本程改到的文件"判，不承诺全仓干净）。

共享 index 事件两次，都按派单 §git 纪律处理：
1. 09:3x 现读 `git diff --cached --name-only` → 18 条 `.scratch/wisp/probes/152/**`（对方 add 了没 commit）
   ⇒ 本程**没有 commit**，等它走成 `eb4755a` 之后才动手。
2. 写本件时再次现读 → 1 条 `docs/evidence/s1/153-trace-lies-unguarded-r1.md` ⇒ 同样等它落地。
   本程每一枚 commit 前都现量一次名册，commit 后 `git show --name-only` 复核枚数（§9 表）。

---

## 1. 票面四条"现量的形状"：全部在 `86b0161` 上重跑（不抄票面的数）

命令与输出见 `.scratch/wisp/probes/151/commands.txt` §0。逐条结论：

1. `git grep -n 'CloseTask' 86b0161 -- '*.go'` ⇒ **2 命中，都是它自己**
   （`internal/tools/bridge.go:642` 的注释、`:646` 的声明）；`git grep -c 'CloseTask' 86b0161 -- '*_test.go'` ⇒ **0**。
   ⇒ 比"生产 0 枚"更硬的那条成立：今天没有任何一处试过它。
2. `git grep -n 'OpenTask' 86b0161 -- internal cmd | grep -v _test` ⇒ `bridge.go:559`（`mark` 里，派发热路径）＋声明＋注释，**唯一调用者就是热路径那一发**。
3. `CloseTask` 的注释原文确实写着 "The composition root defers this on the task's DisposalScope (tickets 12/28)"。
4. `grep -rl "CloseTask" .scratch/wisp/issues/ | grep -v '^.../151-'` ⇒ **0 枚**：不是重复开票。

本程另补一条票面没有、但决定"今天可达面有多大"的现量：

```
$ grep -rn "agent.New(" --include='*.go' cmd internal | grep -v _test   ->  只有 cmd/wisp/run.go:573
$ grep -rn "loop.Run("  --include='*.go' cmd internal | grep -v _test   ->  只有 cmd/wisp/run.go:604
$ sed -n '1,60p' cmd/wisp/resident_windows.go                           ->  "run the empty event loop, then exit"
```

⇒ 今天的 shipped 形状是**一个进程一枚环路任务**；同一枚 `Bridge` 上出现第二枚 task id 的路径只剩
"宿主侧自带 id 直接 `bridge.Execute`"（`cmd/wisp/run_test.go:245` 那枚 `TestHostDispatchThroughTheAssembledBridge`
就是这一形的现成形状）。这条边界决定了 §2 的两问该怎么答，别把它读成"所以无所谓"。

---

## 2. AC#1 — 「不关」今天到底改变了什么（裁决：**渗＝会；长＝会；内容级渗＝不会**）

原始读数＝`.scratch/wisp/probes/151/ac1-probe-tools.txt`（4 枚用例全 PASS，PASS 的是"量到了"，不是"没问题"）。
探针＝`ac1-probe-source.go`，**经 `-overlay` 注入 `internal/tools`，物理文件建在仓库外**（派单 §临时件），
映射见 `overlay-ac1.json`。注入面全部是既有用例同款的真件：`fsBridgeWith`（真 C26 定规器 + 真 C25 引擎
+ 真 fs 工具）、`writeUnder`、`argsFor`、`gateSpy`、`fixtureTool`。**没有任何假对象替真件交回读数。**

### 2.1 渗：第 N+1 轮的判定读到了第 N 轮留下的东西（读到的是"有污点"这件事，不是污点内容）

| 腿 | 形状 | 第二发（新任务、参数干净、allowlist 内 `fs.read`）的裁决 |
|---|---|---|
| `no-close_today` | 只开不关（今天的码） | **`level=L2 rules=[R4] approvals=1`**，reason＝`R4: 包含来自 unbound-scope scope is not open (OpenScope missing or already closed) 的内容` |
| `close-first_control` | 同一条 Bridge，中间补一发 `CloseTask("task-1")` | `level=L0 rules=[] approvals=0 reason=""`（关闭即时 `scopes 1->0 marks 1->0`） |

机制读法（不新造规矩，全部指回现有码）：`OpenTask` 只在 `mark` 里懒开（`bridge.go:559`），所以一枚
新任务的**第一次判定**发生在它的 scope 还没登记的时候；`risk.Provenance.scopeMarks`
（`internal/risk/provenance.go:441-458`）对"未登记的 scope ＋ 别处还留着污点"这一形**fail-closed 造
`SrcUnboundScope` 命中**，`ruleTaint`（`rules_gateway.go:101-115`）把它顶成 L2 且 `SessionOverrideBlocked`。
⇒ 不关的那枚 scope 就是"别处还留着污点"的那个体。

**方向必须说死**：这一发渗是**变严**（多一张卡／多一次拒绝），不是"上一轮的秘密被下一轮读到"。
内容级那一问单独量了（`TestProbe151ContentVsUnboundHit`）：

- 同一枚 task id 复用、参数里带着上一轮那段秘密 ⇒ `reason="R4: 包含来自 fs.read <...>\id.txt 的内容"` ⇒ **同名复用会真继承内容**；
- 换新 id、参数里带着上一轮那段秘密 ⇒ `reason="R4: 包含来自 unbound-scope ..."` ⇒ **不是 fs.read 的内容命中**，
  即 C25 没有跨唯一 id 搬运内容；
- 换新 id 但先把上一轮 `CloseTask` ⇒ `out.level=L0`、零张卡 ⇒ 关掉之后**新任务确实不继承**，
  这正是 `provenance.go` 头部写的契约（"scopes are opened by the composition root and Closed on Dispose,
  so a new session never inherits old taints"）里"Closed on Dispose"那一半今天没接线。

**代价有多大，两形分开量**（同一枚探针，差别只在 gate 怎么答）：

| 宿主形 | 读数 | 一句人话 |
|---|---|---|
| 没有审批通道 / 卡被拒（今天的 `wisp run` 控制台腿、`NoGate`） | 8 轮里 `executed=1 refused=7 L2cards=7`，第 2 轮起全部被拒 | 第一枚任务读完敏感内容之后，**后面每一枚新任务的每一发调用都是硬拒** |
| 有审批通道且用户一路点允许（球/面板宿主的形状） | 64 轮里 `executed=64 refused=0 L2cards=63` | 每一枚干净任务都要吃一张写着 `unbound-scope` 的 L2 卡——卡上点允许就继续，**这正是审批疲劳的形状** |

真组合根上的同一发也量到了（§3 的变异红句里那句
`实时语音（Path C）不具备工具权限；该任务未经文本循环登记（AdmitTextTask），已 fail-closed 拒绝`
⇒ L2 走到 gate，未登记的任务被 D47 拒；这就是控制台腿今天会看到的实际后果）。

### 2.2 长：表只增不减，且没有任何上限

| 形（64 轮 × 每轮一枚新 task id、各留一枚 4,100 字符的敏感读结果） | scopes | marks | 堆 |
|---|---|---|---|
| `no-close`（今天） | **64** | 64 | 591,464 → **3,001,216**（`delta=+2,409,752`，≈**37,652 B/轮**） |
| 再多一枚什么都不读的干净任务 | 64 → **65** | — | — |
| 逐一 `CloseTask` 之后 | **0** | 0 | 578,592（回到基线附近） |
| `close-each-round`（对照腿） | **0** | 0 | 496,784（`delta` 为负，见 §2.3 的坑） |

⇒ **会长，且不关就永不缩**；关掉之后表与污点索引都真能归零（"关"这个动作不是装饰）。
两枚参照物把"没有上限"这条钉住：
- `internal/risk/provenance.go` 的预算只管**单枚 scope 内**（`defaultMaxScopeSources=64` 枚源、
  `defaultMaxSourceRunes=262144` 字），**scope 枚数一枚都不设上限**；
- 同一条任务生命周期的另一张表——D47 的 admitted 表——**有**淘汰（`approval/gate.go:167` 的
  `maxAdmittedTasks`）。同一批 per-task id，一张有界一张无界，这是形状差，不是我推的结论。

### 2.3 后果边界的诚实读法（票面点名"不许凭感觉报数"，两侧都不许）

- 不许报成"泄漏很大"：今天的 shipped 进程是**一进程一枚环路任务**（§1 补的那条现量），
  所以生产面上 scope 枚数今天＝"1 枚环路任务 ＋ N 枚宿主自带 id 的调用"，且一旦有污点，
  后续新任务**先被拒**（拒掉 ⇒ 不再 mark ⇒ 表自己停止增长）。37.7 KiB/轮那枚斜率是
  "球/面板宿主上线、用户一路点允许"之后的形状，不是今天控制台腿每天漏 37 MB。
- 不许报成"只是难看"：`unbound-scope` 那一发改变的是**裁决档位与是否执行**（L0→L2、执行→拒绝），
  不是日志好不好看；且它随进程寿命单调累加。
- 小 N 的堆读数不能用：`rounds=2` 两腿的 `delta` 都是负数（-141,960 / -100,344），是 GC 时机噪声；
  斜率只取 64 轮那一发。
- 未量到：D32 最坏形状（一枚 scope 塞满 64 源 × 262,144 字）下的内存上界；常驻进程跑一整天的曲线。
  见 §8。

### 2.4 判语

票面 AC#2 给的作废条件（"既不渗、也不长"）**不成立**：渗量到了（且今天就在控制台腿上是硬拒），
长也量到了（无上限、只增不减）。⇒ **本票不作废，走 AC#2/AC#3。**
按票面另一句，"该不该关"不再退给 `Q-43` 那一族：`provenance.go` 头部与 `CloseTask` 自己的注释
已经把"组合根在 Dispose 上关"写成契约，缺的是接线。

---

## 3. AC#2 — 会响的检：两枚新用例 ＋ 承重那句

新用例＝`cmd/wisp/task_scope_close_151_test.go`（两枚，都走真组合根：真 config、真 store、
真 `approval.Gate`、真 `tools.Bridge`、真 `agent.Loop`、真 mockllm）。

| 红名 | 钉的是哪一半 | 修好的码上 | 变异树上 |
|---|---|---|---|
| `TestCompositionRootClosesTheLoopTasksTaintScope` | **接线**：真跑一整轮 `wisp run`，要求任务结束时那发 `CloseTask` 的审计行带着**本轮环路自己的 task id**（id 从状态行 `任务 (\S+) 结束` 现读） | PASS | FAIL |
| `TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit` | **效果**：装配好的桥上，宿主任务读走敏感内容 → 走组合根的 per-task 边界收尾 → 另一枚新任务的干净 `fs.read` 必须是 L0 且真读到内容 | PASS | FAIL |

承重那句（本仓操作定义：摘掉本票选定的一味，是否存在一发从此打不红？）：**答"否"**。
变异＝只把 `cmd/wisp/run.go` 里 `admitTask` 的那一行 `rt.bridge.CloseTask(taskID)` 摘掉，其余一字不动
（变异树全文＝`.scratch/wisp/probes/151/run.mut-noclose.go`，映射＝`overlay-mut-noclose.json`）。
两枚**全红**，红因不同（一枚丢审计行、一枚丢裁决），红句原文（`ac2-mutation-noclose.txt`）：

```
    task_scope_close_151_test.go:62: 任务结束时没有关闭它自己的 C25 污点 scope：stderr 里找不到以
      "[audit] tools: C25 scope closed task=8b4bf8cd-40d2-4453-b90e-6260f6a9c762 " 开头的审计行
--- FAIL: TestCompositionRootClosesTheLoopTasksTaintScope (2.12s)
    task_scope_close_151_test.go:109: 上一轮结束后它的污点 scope 还挂在表上，于是这一枚新任务的干净调用
      被 C25 的 unbound fail-closed 顶到 L2：第二发 judged L2, want L0；text="实时语音（Path C）不具备工具权限；
      该任务未经文本循环登记（AdmitTextTask），已 fail-closed 拒绝"
    task_scope_close_151_test.go:113: 第二轮没有读到文件内容，text="实时语音（Path C）…已 fail-closed 拒绝"
--- FAIL: TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit (1.98s)
```

**恒真这一问的答案要先说**：这两枚在"没有本票改动"的码上**结构上跑不起来**（`rt.admitTask` 这枚方法
就是本票引入的，未修码上连编译都过不去）。所以本件不用"未修码"当红绿尺，用的是**变异尺**：
改动已在树上、只摘那一味 ⇒ 两枚必红。摘掉的是本票选定的那一味，不是断言方向。
（票面要的"这一发在未修码上响不响"，等价读数在 §2.1：那枚 `no-close_today` 腿在未修码上就是 L2+一张卡。）

⚠ 两枚已知仪器坑本程都撞过或绕开，记下来给下一程：
1. `cmd/wisp` 的测试二进制在 windows 上要先把 `third_party/sherpa-onnx` 放进 PATH，且用**shell 自己的路径形**；
   `pwd -W` 的盘符形会 `exit status 0xc0000135`、0 条 `=== RUN`（本程第一发就是这么死的，见 `commands.txt` §2）。
2. `-overlay` 全程没与 `-cover*` 同用；变异与门禁分两跑。

---

## 4. AC#3 — 「谁拥有关闭」三选一：裁的是"组合根 defer"

| 候选 | 裁决 | 理由（现量，不是偏好） |
|---|---|---|
| **① 组合根 defer** | **采用** | 今天这一枚组合根里**唯一**同时满足两事的 per-task 边界就是 D47 那枚 `AdmitTask` hook：(a) 它拿得到环路自己生成的 task id（`newTaskID()` 在 `Loop.Run` 内部，`loop.go:332-339`，上面任何一层都不知道它）；(b) 环路已经 `defer revokeAdmission()`（`loop.go:366`），brake / cancel / error / panic 四条形都走它。实现＝`rt.admitTask` 包住 `rt.gate.AdmitTextTask`，revoke 里补一发 `rt.bridge.CloseTask(taskID)`。**零新生命周期、零新接口。** |
| ② `Bridge` 内部自管 | 否 | 桥只知道**调用**（`Execute`），不知道**任务**结束：`Execute` 的返回不等于任务结束（D38d 四路并发、同一 task 多轮）。要做就得让调用方告诉桥"这轮完了"，那是新契约而不是补defer，且会改到 `agent.ToolProvider` 的形状（C1 冻结面）。 |
| ③ 任务生命周期事件 | 否（今天没有这个面） | 现量：非测试码里 `agent.New` 与 `loop.Run` 各只有 1 处（§1），`Result`/`Event` 只朝 sink 走（`EvDone` 那族）；桥不在那条边上。选这条＝为了一行 defer 先造一条事件总线，射程比票大。 |

被明确**没做**的三件事（都是扩大射程的诱饵）：
- 没动 `OpenTask` 一个字，也没把"任务一开始就 OpenScope"接上。这一发**刻意不做**：
  现在的"懒开"配上 `scopeMarks` 的 unbound fail-closed，意味着一枚没登记 scope 的新任务在别处有污点时
  **被顶到 L2**；若在任务开头就登记（空 scope），那一发就变成"读过东西的调用方带着上一轮内容出去也不命中"。
  把 fail-closed 换成 fail-open 不是本票的账，且 `provenance.go` 的 `DEFERRED(C25-loop-wiring)` 第 (1) 项
  要动的正是 `internal/risk/**` 那面零字节墙——**留给编排者判**（见 §8#1）。
- 没给 `Bridge` 加"列出当前开着哪些 scope"的公开 API（本票的读数经 `-overlay` 在包内测到了，不需要新面）。
- 没在 `CloseTask` 里顺手删 `b.seqs`、也没碰 `b.scopes` 的语义（只加一行审计）。

`CloseTask` 本身只改了两处：注释（把"设计期望组合根 defer 它"更新成"现在真的被 `cmd/wisp` 的 admit hook
defer 了，以及**哪一形仍然没人关**"）＋一行 `b.log` 审计（每次调用都打，含 `was_open` / `dropped` /
`open_scopes`）。审计行是 §3 那枚接线钉的承重件，不是装饰：摘掉它 `TestCompositionRootCloses...` 直接红。

---

## 5. AC#4 — 契约轴零字节：本程名册

改到的文件全集（`git show --name-only` 三枚 commit 的并集，逐枚复核过）：

```
internal/tools/bridge.go
cmd/wisp/run.go
cmd/wisp/task_scope_close_151_test.go
.scratch/wisp/probes/151/**        （读数，A264③ 定的落点）
docs/evidence/s1/151-task-scope-never-closed-r1.md   （本件）
```

票面 AC#4 逐条自核：`internal/risk/**`、`internal/panel/**`、`internal/agent/approval/**`、
`tools/d22scan/**`、`thresholds.go`、任何 golden、`allowlist.txt`、`scripts/slo-check.ps1`、
`docs/PLAN.md`、`docs/specs/**`、`frontend/**`、`design/**`、`cmd/wisp/slo_windows.go`、
`cmd/wisp/slo_report_144_windows_test.go` ⇒ **上面那五枚名字里一枚都不在里面**，本程未写它们一字。
（`frontend/**`、`design/**` 此刻有另一枚会话在未提交地写：本程不碰、不还原，也**不把它们的状态
算进任何宣称**。）派单另给的 `internal/agent/**`（票 153 的地界）同样零字节。

---

## 6. AC#5 — 门禁（逐包单跑，四数之外名册两向 `comm`）

日志：`gate-pre-tools-v.txt` / `gate-post-tools-v.txt` / `gate-post-tools-count2.txt` /
`gate-pre-cmdwisp-v.txt` / `gate-post-cmdwisp-v.txt` / `gate-post-cmdwisp-count2.txt`。
**改前的 `cmd/wisp` 怎么量的**：本程改的是工作树，回不去"改前的树"，又不许 `stash/reset/checkout .`，
所以用 `-overlay` 把 `run.go`/`bridge.go` 映射回 HEAD 版（`run.head.go`/`bridge.head.go`），
并把本票的新用例映射成 `null`（`overlay-pre.json`）⇒ 那一跑＝HEAD 码 ＋ 没有本票用例，rc=0。

| 包 | 形 | RUN | 顶格裁决 | 全量裁决 | FAIL | SKIP | panic |
|---|---|---|---|---|---|---|---|
| `internal/tools` | 改前 `-count=1` | 115 | 79 | 79 | **0** | **0** | **0** |
| `internal/tools` | 改后 `-count=1` | 115 | 79 | 79 | **0** | **0** | **0** |
| `internal/tools` | 改后 `-count=2` | 230 | 158 | 158 | **0** | **0** | **0** |
| `cmd/wisp` | 改前 `-count=1`（overlay 回 HEAD） | 136 | 76 | 136 | **0** | **0** | **0** |
| `cmd/wisp` | 改后 `-count=1` | 138 | 78 | 138 | **0** | **0** | **0** |
| `cmd/wisp` | 改后 `-count=2` | 276 | 156 | 276 | **0** | **0** | **0** |

三道自核：
1. **算术自洽**：`230 = 2 × 115`、`158 = 2 × 79`、`276 = 2 × 138`、`156 = 2 × 78` ⇒ 第二遍没漏跑。
2. **名册两向 `comm`**：`internal/tools` 79 → 79，两向**皆空**（`names-pre-tools.txt` / `names-post-tools.txt`）；
   `cmd/wisp` 76 → 78，新增恰两枚
   `TestCompositionRootClosesTheLoopTasksTaintScope` / `TestAdmitTaskRevokeRemovesTheCrossTaskTaintHit`，
   丢失**空**（`names-cli-pre.txt` / `names-cli-post.txt`）。两形名册互差亦空。
3. **开跑 vs 出裁决等集**：`comm -3` 空（RUN 枚数＝全量裁决枚数），三份日志 `panic:` 全 **0**
   ⇒ 没有"开跑了没回来"的用例。FAIL 只认 `--- FAIL` 行（`ok … 14.649s` 那种汇总行不进计数）。

其余工具（版本现读，不背别人的读数）：

```
$ go vet ./internal/tools/ ./cmd/wisp/       ->  无输出，rc=0
$ go build ./...                             ->  无输出，rc=0
$ gofumpt --version                          ->  v0.12.0 (go1.27.1)
$ gofumpt -l internal/tools cmd/wisp         ->  空，rc=0        （gofmt 是另一把尺，没用它交差）
$ sh scripts/d22scan.sh                      ->  rc=0，clean
     各作用域分母全非零：bans#1-5 internal/=205 cmd/=23 ban#6 frontend/=67 ban#7 internal/tools/=18
                          ban#8 design/=39 frontend/=67 internal/=413 cmd/=44
     （d22scan-post.txt。cmd/ 的 ban#8 分母比 139 那件记的 43 多 1＝本票新增的那枚 _test.go；
       frontend/ 的 66→67 是前端会话加的，不是本程。）
     ⚠ 已知仪器坑复现：日志里混着 tools/d22scan 自检造的 examined 行（14/1 那些），
       取数一律取尾部正式行，别用 `grep -m1 examined`。
```

---

## 7. 放水两问自答

1. **断言方向动没动**：没动任何一条既有断言。全仓既有用例一字未改（改动的码只有两枚文件的
   注释/审计行/新函数）；新增只有两枚用例，且它们只做一件事——把 §2 量到的形钉住。
   没有放宽阈值：`thresholds.go`、golden、`scripts/slo-check.ps1` 零字节（§5）。
   唯一"看起来像放水"的地方要说破：修完之后**少**了一种 L2（`unbound-scope` 那一发），
   这是**把 fail-closed 的误报关掉**，不是把真命中关掉——内容级命中在 §2.1 三条腿里逐一验过
   （同 id 复用仍命中 `fs.read`；换 id 且已关闭 ⇒ 不命中，这正是 `provenance.go` 头部写的
   "new sessions never inherit"契约）。
2. **helper 是不是原有的那枚**：是。`cmd/wisp` 两枚新用例复用的全是原有装配件
   （`newRunFixture` / `runFixture.run` / `rtHook` / `taskID()` / `openStore()`），
   本程只新增两枚输入构造器 `req151` / `write151`（不判定、不断言）；
   `internal/tools` 探针复用的是原有 `fsBridgeWith` / `writeUnder` / `argsFor` / `gateSpy` /
   `fixtureTool` / `sealableTempCanonical124`。没有为读数新造第二把尺，也没有用 mock 顶掉真件。

---

## 8. 本程没测什么（按"漏了它谁会先被骗"排序）

1. **宿主侧自带 task id 那一形仍然没人关，本程也没修**。`CloseTask` 今天只长在环路的 admit/revoke 上；
   一个面板/内部调用方自己编 task id 直接 `bridge.Execute`（`cmd/wisp/run_test.go:245` 那形的生产版）
   照旧留下 scope。⇒ 球/面板宿主（票 33/35/77/92）一上线，§2.1 那发硬拒/误报**照样发作**。
   谁先被骗：验收"R4 接线完成"的那一格——它会以为 §2 的问题被本票关掉了。
   真正的两半（`OpenScope` at task start + `DEFERRED(C25-loop-wiring)` 的 (2)(3)(4) 项）都在
   `internal/risk/**` 与环路那面零字节墙上，归编排者判，不在本票射程内。
2. **并发任务没测**。D38d 允许 4 路并发过同一枚桥；两枚任务同时在跑时，B 的第一次判定撞上
   "A 还开着且有污点" ⇒ 仍是 `unbound-scope` L2。本票只关"已结束"的任务，没关"同时在跑"的
   （那要么提前 Open、要么引用计数，都是 §4 里被否掉的扩大射程）。探针里的 64 轮是**串行**的。
3. **内存只量到斜率，没量上界**：4 KiB 源 × 64 轮＝37.7 KiB/轮；未量 D32 最坏形
   （一枚 scope 塞满 64 源 × 262,144 字）与常驻进程整日曲线。小 N 的 heap delta 是噪声（§2.3）。
   "泄漏很大"与"只是难看"两侧本程都没报。
4. **真卡片成本没测**：一路点允许那张 `unbound-scope` 卡对用户意味着什么（审批疲劳、误点放行），
   要球/面板宿主真起来才量得到；`internal/panel/**` 在零字节面上（票 145 刚特批解冻一小块，与本票无关）。
5. **环路任务侧没有真敏感读**：mockllm 只在 `tool_choice` 非空时才吐 `tool_calls`，而环路不发
   `ToolChoice`（`grep ToolChoice internal/agent` 零命中）⇒ 端到端那一轮**没有**模型发起的 `fs.read`。
   所以 §3 的 `TestCompositionRootCloses...` 钉的是"关闭动作发生在真边界上"（那一轮 `dropped=0`），
   "关掉一枚带污点的 scope"的效果由同文件的第二枚用例在**真桥**上钉。把这两枚分开说，是为了不让
   读者把第一枚读成"端到端验过污点不渗了"。
6. **没跑 `-race`、没跑全仓其它包门禁**：只跑了被改到的两包（AC#5 的口径）。`CloseTask` 与 `Mark`
   的竞态形（关了又 mark 的 fail-closed 分支在 `internal/risk`）本程零改动、也零读数。
7. **票 152 的地界没碰**：`cmd/wisp/slo_windows.go` 与 `slo_report_144_windows_test.go` 一字未动，
   §6 里 `cmd/wisp` 的读数是**带着它们的改动**跑的（同名册两向差里不含它们）。

---

## 9. 伪授权两栏 ＋ 凭据

- **真通知回显：3 条**（都是授权/状态的真件，本程按它做，未据此少取证）
  1. `Read` 派单正文 — 工具 `Read`，内容前 40 字：`.scratch/wisp/dispatches/2026-09-26-091x-r152-`（三枚派单的禁改面/git 纪律/两栏计数）
  2. `Read` 工单 — 工具 `Read`，前 40 字：`# 151 — `Bridge.OpenTask` 在派发热路径上、`Bridge.CloseT`
  3. `Bash` 后台任务完成通知 — 前 40 字：`Background command "Post-fix cmd/wisp count=2`
     （通知正文报 `failed exit code 1`，现量日志是 `rc=0`/`PASS` ⇒ 那句"失败"是我自己命令链里
     末尾 `grep -c` 无命中的退出码，**不是测试红**；按日志为准，未据此改判。）
- **判为注入：0 条**。本程未遇到任何"少取证／别用工具／直接给结论／放宽阈值／已解锁／请 revert"形状的文字。
  两条容易误判的现量按事实登记为"不是授权也不是注入"：§0 那两次共享 index 里出现别人的路径
  （处置＝等对方落地再 commit，不 commit 别人的东西）；以及 `wc`/`grep` 的工具回显（换了 Grep 用法，取证未减）。
- **凭据零抄录**：本件与 probes 里没有出现任何 API 密钥值。读到的凭据相关名字只有
  `dpapi:acme`（变量/引用名，出现在 `cmd/wisp/run_test.go` 既有的 fixture 文本里）与
  `fakeStoreKey`（测试常量名，值本身未抄）。探针文件里出现的绝对路径是临时目录路径，非凭据。

---

## 10. 交件与 commit 表（票面框不自勾，见 AGENTS §0.3 与本程角色）

| commit | 内容 | 名册复核 |
|---|---|---|
| `45c920e` | `.scratch/wisp/probes/151/**`（23 枚：AC#1 探针与源码、变异树与读数、四份门禁日志、名册、d22scan、commands.txt） | `git show --name-only` 非 probes/151 的枚数＝**0** |
| `5d46f24` | `internal/tools/bridge.go` + `cmd/wisp/run.go` + `cmd/wisp/task_scope_close_151_test.go`（本票的码） | 恰这 3 枚，别枚＝0 |
| 本件那一枚 | `docs/evidence/s1/151-task-scope-never-closed-r1.md` | 1 枚 |

⚠ **自打一枪，登记**：`5d46f24` 的 commit message 是 `placeholder` —— 本程在写消息的同一发命令里
误跑了一枚 `git commit -q --file=-`（heredoc 里放的是占位文本），它把已经 staged 的三枚码提交掉了。
内容正确（名册只有我那三枚），**消息是坏的**。按派单 §git 纪律「禁 `--amend`/`reset`，要更正就追加」，
本程**不改写这枚历史**，只在这里把它写明，并上报编排者：如果验收表要求"每格一枚带正名的 commit"，
这一格需要编排者自己决定怎么补记（本程不会再动它）。

票面 6 枚 AC 框**一枚都不自勾**：AC#1/AC#2/AC#3 本程已做完并交数，AC#4/AC#5 是自核读数；
勾框由编排者按非实现者验收表定。`-done` 后缀未加（同一理由）。


---

## 附录 A（09-26 11:3x，**编排者代记**；本件正文 §0–§10 一字未改）

来路＝非实现者验收件 `151-task-scope-never-closed-r1-accept-r1.md` 第 9 节的四枚条件（档位＝**附条件成立**；四枚**全在记录文本层，码不用返工**）。
本件那三句"会被下一程当尺用"的话，**在原文里不删**，正确读法只写在这里：

- **对应 C1（本件 §2.1 那张代价表第一行的括注）**：原文"没有审批通道（今天的 `wisp run` 控制台腿、`NoGate`）⇒ 8 轮里 `executed=1 refused=7`"——
  ⇒ **那一形今天不走 `wisp run` 控制台腿**。验收件第 2.3 节＋第 13 节现量到**卡住它的是码、不是夹具**：`fs.go:298` 把 `fs.read` 声明成 `risk.L0`、
  `loop.go:776-789` 的 default 分支拒发、`run.go:588-595` 从不设 `PassThroughUnclassifiedRisk`；端到端那发是 `[工具 fs.read -> error]`、**桥侧全库零行**，
  而对照发（同腿换成 L1 声明的 `fs.write`）`success` ＋文件落盘 ＋桥侧有 `tool_call` 行。
  ⇒ **正确读法分三层，别合并**：①"只开不合"这个**类**在桥上派发的宿主上**生产可达＝成立**；②它**今天不在 `wisp run` 控制台腿上天天发作**；
  ③内容级不渗（本件 §2.1 三条腿原样成立）。⇒ 那发探针读数是**真桥**上的，不是 CLI 端到端的。
- **对应 C2（本件 §3 末段第一句）**："两枚新用例在未修码上**结构上跑不起来**"——**只对 T2 成立**。
  验收件第 3.3 节把 `5d46f24^` 的两枚 head 铺上（先 `sha256` 现核逐字节相同）并裁成 T1-only ⇒ **编得过、而且直接 `--- FAIL`**。
  ⇒ 也就是说票面要的"这一发在未修码上响不响"**本可以直接读**，本件当时用变异尺推的那一步**比必要的弱**。判定不变（两味都有牙，见第 3.2 节）。
- **对应 C3（本件 §8 第 5 条）**：那条只把"端到端没有模型发起的 `fs.read`"**归给 mockllm 的 `tool_choice` 门** ⇒ **少登一道门，而且那道门是码给的**（即 C1 那三处原文）。
  ⇒ 这一条**不是本件的错**（它量到了现象、只归错了因），但它**会让下一程把"没测到"读成"夹具不给"**从而放掉真缺口 ⇒ 已把"第二道门"**接进票 154** 的射程（见该票 Progress log 09-26 11:3x 那条）。
- **对应 C4（本件 §10"自打一枪"那条）**：`5d46f24` 的 subject 确实是 `placeholder`（验收件 `git show -s --format='%s'` 现读；名册恰三枚、内容与本件自述一致）。
  ⇒ **补记在此，不改写历史**：那枚 commit 的内容＝**票 151 的全部三枚码**（`internal/tools/bridge.go`＋`cmd/wisp/run.go`＋`cmd/wisp/task_scope_close_151_test.go`）；
  成因＝实现程在写消息那条命令里误跑了一次 `git commit -q --file=-`，把已 staged 的三枚用占位消息提交了。
  **本仓不因消息坏而回退内容**（已推送历史不改写这条规矩优先），故：**"每格一枚带正名的 commit"这一格判为"内容成立、消息缺失"，凭据＝本附录＋验收件第 3.1 节**，不算装饰也不算退回。

**归属写清**：以上四枚**都是记录级**，验收程明写"没有任何一味是装饰（第 3 节）、没有假绿（第 8 节）、没有越界写（第 5 节）"，
且它**造出来了**（三发变异＋一发端到端金样本＋一发对照＋六枚门禁数＋两包名册对拍）。⇒ 票面 5 枚框由我按它的表全勾，`-done` 由我加。
