# 155 — 票 153 剩下的三块"无人裁"·实现件 r1（只读取证，不写生产码）

- 工单：`.scratch/wisp/issues/155-ticket-153-leaves-three-cells-unjudged-…-coordinate-system.md`
- 派单：`.scratch/wisp/dispatches/2026-09-26-115x-impl-155.md`
- 本程身份：**实现方（只读取证）**。票面框**一枚不勾**、`-done` **不加**、**零 push**。
- 写面只有两处：本文件 ＋ `.scratch/wisp/probes/155/**`。
- 档位本程不自立：三块各自"能定档到哪"写在 §4 那张推进表里，**由编排者按非实现者验收表合两份表定**。

## 0. 锚点与仪器名的反查（本程第一发；派单自己抄的号，先证再用）

现量凭据＝`probes/155/00-anchor-reverse-check.txt`；本程时刻 `2026-09-26 11:39 +08`（`date` 现取），分支 `dev`，HEAD 当时＝`317805c`（收口时又量一次，见 §5 末）。

| 派单/票面点名的号 | `git cat-file -t` | 是什么 | 成不成立 |
|---|---|---|---|
| `ff550f3` | commit | `test(153 AC#1)` 那一发证据的 commit | 成立 |
| `ff550f3^` | commit | **展开＝`5365cb22`**（`git rev-parse ff550f3^` 原文） | 成立 |
| `5365cb22` | commit | `docs(派单存档)`，与 `ff550f3^` 是**同一枚** | 成立 |
| 验收件 §12.3 用的 `86b0161` | commit | `feat(frontend 10 首启引导)`——派单那两枚的**父的父** | 成立 |
| 生产改动 `b23c7f7`（`git log -- internal/agent/compress.go` 现读） | commit | `feat(153 AC#2)` 痕带上 taskID | 成立 |
| 到盘复算锚点 `6de3d1c5`（验收件 §12.3 第 3 条） | commit | `evidence(153 AC#5 第 5 格＋收口)`；`merge-base --is-ancestor 6de3d1c5 HEAD` ⇒ 已进主线 | 成立 |

**两枚锚点先对质**（派单说 `ff550f3^`、验收件 §12.3 说 `86b0161`——不是同一枚 commit，本程不猜，先量内容）：

```
file                         anchor      blob(short12)
internal/agent/compress.go   86b0161     75758f64aa43
internal/agent/compress.go   5365cb22    75758f64aa43
internal/agent/compress.go   ff550f3^    75758f64aa43
internal/agent/compress.go   b23c7f7^    75758f64aa43   ← 生产改动那枚的真父
internal/agent/loop.go       （同上四枚） e685a6efd595
```

⇒ 三枚"改前锚点"在**本票要看的那两枚文件上是同一枚内容**，且都等于 `b23c7f7^`；`ff550f3` 本身也在 `b23c7f7` 之前（`merge-base --is-ancestor` ＝ YES）。
⇒ **派单那两枚号可用**，本程按派单的号跑（`ff550f3^`／`5365cb22`），并顺手记下它与验收件那枚同源。中间夹的两枚 commit（`cdf24711`／`ad9b29fd`）只动 frontend，与本票无关——这条也是现量的，不是推的。

**仪器名反查**（派单点名的两枚，都真）：

- `.scratch/wisp/probes/153/mut153.py`（4215 字节）——`drop-attr`／`unloop` 两枚 op **就在这枚文件里**（`:92`／`:106`）。
- ⚠ 同名不同物的第二枚：`.scratch/wisp/probes/153/accept-r1/mut153.py`（5981 字节）**没有** `drop-attr`／`unloop` 这两枚 op 名，它用的是 `A4-dropN2-looptag`／`A6-dropN3-attrwrite`／`A0-baseline-delivered`。本程按派单原文用了 `probes/153/mut153.py` 那枚；两枚的**变异内容**本程逐条比过锚串，是同一处（`withTraceTask(ctx, taskID)` 摘掉／`if task := traceTaskID(ctx); task != "" { … }` 那三行摘掉）。**这一处不报清楚，下一位按名字挑一枚就会挑错文件。**
- `.scratch/wisp/probes/153/accept-r1/zzaccept153_ondisk_test.go`（到盘探针）在场；`.scratch/wisp/probes/153/probe153_line_test.go`／`probe153_loopleg_test.go`（实现件 ext 那两枚）在场。

**本程新撞出来的一枚仪器坑（登记，不是读数）**：`probes/155/sync_blob_tree.py` 第一版把 sha 清单**用管道**喂给 `git cat-file --batch`，父进程写满 stdin 时子进程正堵在 stdout 上——**死锁**，跑了 7 分钟落盘 **0 枚文件**、日志 0 字节。改成"清单先落文件、stdin 用文件重定向"后一次跑完。
⇒ 为什么值得占一行：那枚"跑完 0 产物"的形状**看起来像失败、实际像没跑**，而它当时已经被我误当成"完成"读过一次（后台任务通知写的是 killed，我把它记成了 completed）。**下一位拿这类"零落盘"当结论之前，先 `find <tree> -type f | wc -l`。**

**取版自证（AC#4 第①条）**：快照＝`git ls-tree -r 6de3d1c5` 逐 blob `git cat-file --batch` 落盘，**每枚文件本地重算 `blob <len>\0`＋数据的 sha1 与 ls-tree 号对拍**：`written: 1452  mismatched: 0`（`probes/155/10-snapshot-sync.txt`）。三枚被测文件另用 `git hash-object` 复核，与 `git rev-parse 6de3d1c5:<路径>` **逐字相同**（`4fcd9a32a500…`／`3ea1fb8df38a…`／`2357dca86f67…`）。全程**未用** `git archive`。

---

## 1. 格① AC#2 问① —— 改之前 `Compress` 拿不拿得到 taskID？

**答：拿不到，三条通道逐条是空的。**（凭据＝`probes/155/01-ac2-q1-before.txt`，尺打在 `ff550f3^` 的 git 对象上，不读脏工作树）

| 通道 | 现量（原文在探针文件第 1/2/5/7 节） | 读数 |
|---|---|---|
| **参数表** | `git show ff550f3^:internal/agent/compress.go \| grep -n 'func (c \*Compressor)\|^type Compressor struct'` | `:128 func (c *Compressor) Compress(ctx context.Context, hist []llm.Message) ([]llm.Message, CompressionReport, error)` —— 参数表里**没有** id |
| **接收者** | `sed -n '/^type Compressor struct/,/^}/p'` | `type Compressor struct { b Budgets; sum Summarizer; lg *slog.Logger }` —— 三枚字段，**没有** id |
| **ctx 上能不能反查** | `git grep -n 'withTraceTask\|traceTaskID\|traceTaskKey' 5365cb22 -- .` ＝ **0 命中（rc=1）**；`git grep -n 'ctx.Value' 5365cb22 -- internal/agent` ＝ **0 命中（rc=1）**；全仓非测试的 `context.WithValue` 只有 `internal/tools/cancel.go:45` 一枚（cancelHandle，与任务无关）；`observe.NewRootFrom` 把 id 留在 `Root.ID` 结构体字段上、**没进 ctx** | ctx 上**没有**可读的 id，也没有可借的现成通道 |
| **构造时刻塞进去行不行**（问①的第二读法） | `git show ff550f3^:internal/agent/compress.go \| grep -n 'func NewCompressor\|func WithLogger'` → `NewCompressor(b Budgets, sum Summarizer, opts ...CompressorOpt)`、`WithLogger` 是唯一一枚 opt；`git show 5365cb22:internal/agent/loop.go` 的 `:239 comp: NewCompressor(...)` 在 `agent.New` 里，而 id 到 `:322`／`:333` 才诞生 | 构造器**只认 Budgets／Summarizer／Logger**，且**建在任何 id 诞生之前**——这条路不是"没写"，是**位置上就不成立** |

**这把尺不是死尺（正控，两向）**：
- 同一枚 grep 表达式打在**改后**版 `b23c7f7` 上：`withTraceTask`／`traceTaskID`／`traceTaskKey`／`ctx.Value` **12 行命中**（探针文件第 6 节，含 `:161 s, _ := ctx.Value(traceTaskKey{}).(string)`）。
- 改前那版全文 `ctx` 只有 **5 处**（`:39/:128/:152/:197/:201`），**全是透传给 `SummarizeHistory`**，一枚读取都没有。
⇒ 所以"0 命中"是**码里真没有**，不是尺钝了。

**⚠ 一处仪器抖动（自报，不藏）**：`git grep -n 'context.WithValue' 5365cb22 -- '*.go' ':!*_test.go'` 本程第一次单发跑时**输出被吞了那一枚命中并回 rc=1**，重跑三次（同一命令、带 `| cat`）都是 rc=0 且稳定命中 `cancel.go:45`。本程所有"零命中"宣称**一律重跑并记 rc**；这一处如果照第一次的读数写，就会把"全仓只有一枚 WithValue"写成"零枚"——**那是一枚会往上长的错**。

**本格的额外一发（派单没要求，但它是问①的另一半）**：改前那条痕的原文里**根本没有 task 位**——
`git show ff550f3^:internal/agent/compress.go` 的 `if rep.Ran {` 块是九枚键的字面量参数（`tokens_before … history_changed`），第八枚就是最后一枚。
⇒ 这答的不是"拿不拿得到"，而是"改前那行**外部读数**长什么样"，与格③的 ext-A 基线（三腿到盘第一行**没有** `"task"` 键）互指。

**这一块把 153 的哪一枚 AC 从"未裁"推到"可定档"**：**AC#2 的问①**（票面拆出的 6 件里那一件 `(b)`）。
问①的答案是"不能"，因此票面那条"若三问的答案是要动契约面 ⇒ 停手"**在问①这一支不触发**——通道三条全空、构造点位置不对，**"要动的是实现内部的一枚 ctx 键"这件事是被现量逼出来的，不是自选的**。

---

## 2. 格② AC#2 问② —— 唯一持有者是谁（组合根？`Loop`？）

**答：`Loop.run` 那一枚参数 `taskID`（`internal/agent/loop.go:339`）。组合根不是。**
票面给的两个候选里选**后者**，但"唯一"两个字有边界，本程把它说到哪儿为止写在下面第 4 条。凭据＝`probes/155/02-ac2-q2-holder.txt`（13 节，尺全打在 `5365cb22` 的 git 对象上）。

**① 派单点名的那一发，原文照跑**

```
$ git grep -n 'newTaskID\|withTraceTask' 5365cb22 -- internal/agent
5365cb22:internal/agent/loop.go:322:	id := newTaskID()
5365cb22:internal/agent/loop.go:333:	return l.run(ctx, newTaskID(), input)
5365cb22:internal/agent/loop.go:1102:// newTaskID returns a UUIDv4-shaped task id (the task_log.id column is
5365cb22:internal/agent/loop.go:1104:func newTaskID() string {
rc=0
```

⇒ `newTaskID` 四枚命中**全在同一枚文件**（定义 1＋生产调用 2＋注释 1）；`withTraceTask` **0 枚命中**——它当时还不存在（那是本票改出来的东西）。
**这不是"尺钝了"的 0**：同一把尺打在改后版 `b23c7f7` 上命中 **8 行 / 3 枚文件**（`compress.go:4`、`compress_trace_test.go:2`、`loop.go:2`，探针第 1b 节）。

**② id 出不了这枚包**：`git grep -n 'newTaskID' 5365cb22 -- . ':!.scratch/*' ':!docs/*'` ⇒ **仍然只有那 4 行**（排除工单正文与票 139 的快照副本之后，全仓生产码里没有任何包外调用点）。`newTaskID` 是小写开头＝包私有，组合根**结构上拿不到它**。

**③ 压缩那一发就在持有者的作用域里**（这是问②真正在问的东西）

```
$ git show 5365cb22:internal/agent/loop.go | grep -n 'func (l \*Loop) run(\|l.comp.Compress('
339:func (l *Loop) run(ctx context.Context, taskID, input string) Result {
396:			nh, rep, err := l.comp.Compress(ctx, hist)
```

- 生产侧调用 `Compress` 的点**全仓唯一**：`git grep -n '\.Compress(' 5365cb22 -- internal cmd ':!*_test.go'` ⇒ **1 枚命中**（`loop.go:396`）；同一把尺不排测试时是 8 枚（`compress_test.go:4`＋`compress_trace_test.go:3`＋`loop.go:1`）⇒ 测试腿 7 枚、生产腿 1 枚，本程没把测试调用点算进"持有者"里。
- `:396` 落在 `:339` 那一枚函数体里 ⇒ **参数 `taskID` 当场在作用域内**。这正是票 153 §地界 预留的那一支（"loop.go 只读，除非 taskID 必须由它传进来"）。
- `Loop` **结构体自己没有 id 字段**（`type Loop struct` 原文 11 行：`opt/b/guard/asm/sp/comp/reg/provider/mu/history/steer/current`）⇒ "持有者＝`Loop` 这个对象"这种读法**不成立**；持有者是**这一次调用的栈帧**，不是那枚可被复用的对象。这条与格①"构造时刻也对不上"是同一件事的两面。

**④ "唯一"说到哪儿为止（本程不肯把它说满的那半格）**：同一枚 id 在 `l.run` 体内还另落了三处，逐枚带行号（探针第 10–12 节）：

| 副本 | 现量 | 算不算"另一个持有者" |
|---|---|---|
| `root := observe.NewRootFrom(ctx, taskID)` `loop.go:345` | `Root struct { ID string; Ctx; Cancel; … }`（`observe/goroutine.go:88-95`），id 只在**结构体字段**上 | **不算**：`NewRootFrom` 内部只 `context.WithCancel(parent)`，**没有 `WithValue`** ⇒ id 没进 ctx（格①第三条互指） |
| `l.setCurrent(root)` `loop.go:353`（`Loop.current *observe.Root`） | `l.current` 只有 3 处读点：`:348` 注释、`:523` | **不算**：它在 `Compressor` 那侧不可达（不同种类、无引用） |
| `res := Result{TaskID: taskID…}` `loop.go:369` | `Result.TaskID` 声明在 `loop.go:83` | **不算**：那是**出口**，且晚于压缩那一发 |
| `RunningTask{ID: id}` `loop.go:323`（异步腿独有） | `type RunningTask struct { ID string; root; h; done; result }` | **不算**：它在 `reg.Spawn` 的闭包外，同步腿（`Run`）根本不经过它 |

⇒ 精确读法：**"造 id 的通道唯一（包私有 `newTaskID`）、压缩发生那一刻手里有 id 的只有 `Loop.run` 那一枚参数"**；同帧内的四枚副本都在**这一次调用之内**，没有任何一枚住在 `Compressor` 够得着的地方。

**⑤ 组合根被现量否掉**（票面那个问句的另一半）

```
$ git grep -n 'taskID\|TaskID' 5365cb22 -- cmd ':!*_test.go'
cmd/wisp/run.go:476	const taskID = "host:mode-switch"     ← 宿主给 L2 卡片签发的伪任务名
cmd/wisp/run.go:576	(下一发) res := loop.Run(ctx, task)   ← 进 Loop 那一发不带任何 id
cmd/wisp/run.go:587/:598/:608	res.TaskID …                 ← 全是事后读
cmd/wisp/run.go:729/:756	e.TaskID …                      ← 事件侧，事后
```

- `run.go:576` 原文（`sed -n '570,590p'`）：`res := loop.Run(ctx, task)`——`ctx` 是 `context.WithTimeout(context.Background(), …)` 造的，**没有任何任务身份**；组合根第一次拿到 id 是 `:587` 读 `res.TaskID`，**晚于整次环路**，因此**晚于压缩**。
- `:476` 那一枚 `const taskID = "host:mode-switch"` 本程**特别点名**，因为它是一枚长得像答案的东西：它是宿主为一次权限模式切换**自己编的**字符串（同段 `rt.gate.AdmitTextTask(taskID)`／`tools.Decision{TaskID: taskID…}`），**不是** `newTaskID()` 造的那枚，也不在压缩那条路上。拿它抵"组合根持有 id"这一账，就是本仓那种"相近读数顶真问句"的形状。

**这一块把 153 的哪一枚 AC 从"未裁"推到"可定档"**：**AC#2 的问②**（票面 6 件里那一件 `(c)`）。
⇒ 问②裁"组合根**不是**持有者、`Loop.run` 的参数是"，所以"把 id 送进 `Compress`"**不要求扩任何契约面、也不要求改组合根**：唯一需要动的那一发本来就在票 153 §地界 允许动的两枚文件里（`loop.go` 那一发调用点）。它与格①合起来才等于票面那句"若三问的答案是要动契约面 ⇒ 停手"的**前两问**（第三问 (d) 前一程已裁）。

---

## 3. 格③ AC#3 句② —— 摘掉它，有没有任何**外部可见**读数变过？（复算 ext-A／ext-B／ext-C）

**答：变过，两味各自都变；第三味（N1）量出来是"不变"。** 本程**亲跑复算**，不是引用实现件 §4.2 那张表的自述；前一程 §8 第 3 条自报"没复算"、并明写"两程读数不抵账"，所以本程的读数也是**另立一栏**、不与前两程互抵。

**判据摆清楚**：外部可见＝**落到盘上那条 JSONL 记录**（`observe.InitLogWithRegistry` → `redactHandler` → `rollingWriter` → 回读文件）。这比实现件那三发更强一档——它量的是测试进程 stderr 上 `slog.Default()` 打出来的 `task=`，本程量的是**宿主日志文件里的 `"task":"…"`**。

### 3.1 先说复用里那一处不能糊过去的地方（本程最该被审的一行）

派单与验收件 §12.3 说的是"把前一程那枚到盘探针（`probes/153/accept-r1/zzaccept153_ondisk_test.go`）接上 `mut153.py` 的 `drop-attr`／`unloop` 两发"。本程照做，**并且当场量出那枚探针看不见 `unloop`**：

```
ext-C（摘 N2＝loop 不再打标）同一发里，两枚探针并排：
  复用那枚两腿探针： --- PASS: TestAccept153TraceOnDiskThroughRedactor     ← 它照旧绿
  本程三腿探针：      PROBE-155 COUNT trace_lines=3 with_task=1            ← 少了 loop 腿那一枚
```

**为什么**：那枚两腿探针自己直接 `withTraceTask(context.Background(), id)` 打标，**不经过 `Loop`**；`unloop` 改的是 `loop.go` 那一发调用点 ⇒ 对它**结构上不可见**。
⇒ 所以本程**没有**拿它单独交卷（那会把 ext-C 读成"没变"＝一枚假结论），而是**在仓外快照里新增一枚三腿探针**：腿 1 未打标直调／腿 2 打标直调／腿 3 **真 `Loop.Run`**（用 `newHarness(t, "text-reply", …)` 那枚**交付里已有的** golden-SSE 接缝，本程没造新的注入面）。
⚠ 交付码零字节、既有断言零改动：三腿探针只存在于 `D:/tmp/wisp155/snap/internal/agent/`，**没有**进工作树、**没有**进任何一枚 commit（名册见 §5）。

### 3.2 三发读数（同一次快照，每发：`restore` → 施变异（落地必打印）→ `go test -count=1`）

| 发 | 施了什么（`mut153.py` 的 op 名，落地行原文） | 盘上痕总行数 | 盘上带 `"task":"` 的行数 | 与实现件 §4.2 对照 |
|---|---|---|---|---|
| **ext-A** | 无（交付码） | **3** | **2** | 3／2 ＝ 相同 |
| **ext-B** | `drop-attr` ⇒ `LANDED (internal/agent/compress.go): task attribute append removed` | **3** | **0** | 3／0 ＝ 相同 |
| **ext-C** | `unloop` ⇒ `LANDED (internal/agent/loop.go): loop stopped tagging the ctx` | **3** | **1** | 3／1 ＝ 相同 |

三发的**痕总行数恒 3**——这一条是句②的地基：**"摘掉它"摘掉的是归因，不是那条痕本身**。逐行看（ext-B 的三行原文在 `probes/155/11-ext-matrix.txt`）：ext-A 第 3 行 `"tokens_before":826,…,"task":"0bb07666-…"`，同一发 `leg3_ResultTaskID=0bb07666-…` **逐字符相等**（这就是"生产腿带上了它自己的那次运行的 id"）；ext-C 第 3 行同样 `826/532` 那对计数**在场**、但 `"task"` 键**消失**、同发 `leg3_ResultTaskID=305301cc-…` 仍在 → **归因消失了，运行没消失**。

**两向都不抵账的交叉**：实现件用 `task=`（stderr）、本程用 `"task":"`（JSONL）——两把不同形状的尺在 ext-A/B/C 三发上给出 **2／0／1** 同一串数。本程另把**实现件自己那两枚 ext 探针**在同一棵快照里跑了同一发（`probes/155/11-ext-matrix.txt` 每一节的第三段），stderr 侧读数也是 3 行／2、0、1 —— 也就是**复算到的不是它的结论，是它的读数**。

### 3.3 正控（先证明这棵树活着、再谈读数）与还原自证

| # | 那一发 | 读数 |
|---|---|---|
| 正控①（无变异） | `go test ./internal/agent/ -count=1 -v -run TestCompressionTrace` | **7 RUN／0 FAIL**，七枚名册逐枚在场（含 `CarriesTheOwningTaskID`、`NeverInventsATaskID`） |
| 正控②（已知会红） | 同一发上施 `drop-attr` | **2 枚红**：`--- FAIL: TestCompressionTraceCarriesTheOwningTaskID`、`--- FAIL: TestCompressionTraceNeverInventsATaskID` ＝ 与实现件 X5 的红名逐字相同；**且复用那枚到盘探针也红**：`zzaccept153_ondisk_test.go:70: tagged record lost its task on disk: want "a465ea9d-…" in {…}` |
| 还原 | 每发之前 `mut153.py restore` ＋ 收口逐枚 `cmp` | `RESTORED IDENTICAL compress.go / loop.go / compress_trace_test.go`；收口再 `git hash-object` 三枚＝`4fcd9a32a500…`／`3ea1fb8df38a…`／`2357dca86f67…` ＝ **`6de3d1c5` 的 blob 号逐字相同** |
| 树是纯的 | 快照取版＝blob 级（§0 末），`mismatched: 0` | 全程**未用** `git archive`；`third_party/**` 那三枚 dll **没拷**——因为本程**一行 `cmd/wisp` 都没跑**（`internal/agent` 包不需要它们），这条在 §5"没测"里另记一笔，不当"跑过了" |
| 工具版本 | `go version` 现读 | **`go1.27.1 windows/amd64`**（与前一程 §6.2 那枚读数同号；本程现取，未背数） |

### 3.4 本程自加的一发（句②剩下的那一味 N1，实现件是"推"的、本程量了）

实现件 §4.2 对 N1（AC#1 那枚判据用例）答的是**推理**——"`_test.go` 不进交付二进制，所以外部读数不会变"，并写明"本程未为它造外部读数"。这一发本程**量了**：

```
$ python mut153.py drop-test:TestCompressionTraceSilentWhenNothingFoldableOverThreshold …
LANDED (internal/agent/compress_trace_test.go): removed TestCompressionTraceSilentWhenNothingFoldableOverThreshold (lines 426-464)
--- 同一棵树上跑三腿到盘探针 ---
PROBE-155 COUNT trace_lines=3 with_task=2        ← 与 ext-A 基线逐枚相同
```

⇒ **摘掉 N1，外部读数一枚都不变**（3／2，与 ext-A 全等）——实现件那条推理**成立，而且现在有读数了**。
⇒ 顺手把句①那一形在本程手上跑到（同发再施 M5）：`RUN` 从 7 掉到 **6**（证人就是被摘掉的那枚）、`ok github.com/CarlosShao/wisp/internal/agent` **全绿** ⇒ 这正是实现件 X2 的"摘证人⇒逃逸回来"，本程独立复现到同形。

**这一块把 153 的哪一枚 AC 从"未裁"推到"可定档"**：**AC#3 的整枚句②**（票面那枚框的第二个分句）。
⇒ 句①（前一程 §3.2 已裁、本程 §3.3 正控②与 §3.4 各复到一次）＋ 句②（本节）合起来，**AC#3 那枚框两句才都有了读数**；本程不自勾，档位由编排者合两份表定。

---

## 3. 格③ AC#3 句② —— 摘掉它，**外部可见读数**到底变没变（ext-A／ext-B／ext-C 复算）

**答：变了，两枚生产味各自都变；第三枚（判据用例 N1）不变——这一枚本程是量出来的，不是推出来的。**
凭据＝`probes/155/11-ext-matrix.txt`（两回独立跑的读数都在里面，逐枚互校）＋ `probes/155/12-instrument-listing.txt`（仪器名册＋逐枚 sha1）。
快照＝`D:/tmp/wisp155/snap`（**blob 精确树**，`6de3d1c5`，1452 枚全数对拍、`mismatched: 0`；**未用** `git archive`，见 §0 末）。外部读数一律走**到盘**那条链：`observe.InitLogWithRegistry` → `redactHandler` → JSONL → **回读文件**，再按 `"task":"` 数行数（比实现件那三发的"stderr 上数 `task=`"强一档：那是**落在宿主日志文件里**的东西）。

### 3.1 三发对照（痕总行数恒 3，变的是带 `task` 的那几枚）

| 发 | 施了什么（`mut153.py` 的 op 名，落地行原文见探针文件） | 盘上痕行 | 带 `"task":"` 的行 | 第三腿（真 `Loop.Run`）那一行 | 判 |
|---|---|---|---|---|---|
| **ext-A** | 无（交付码基线） | **3** | **2** | `"task":"0bb07666-…"` 且 `leg3_ResultTaskID=0bb07666-…` **逐字符相等** | 基线 |
| **ext-B** | `drop-attr`（摘 N3：痕不再写属性） `LANDED (internal/agent/compress.go): task attribute append removed` | **3** | **0** | `history_changed:true` 在场、**整枚 task 键不存在** | **变过**（2→0）：肉眼可辨 |
| **ext-C** | `unloop`（摘 N2：loop 不再打标） `LANDED (internal/agent/loop.go): loop stopped tagging the ctx` | **3** | **1** | **只剩直接打标那一发带 task**；Loop 腿那行 `826/532` 计数在场但**无归因**，而 `leg3_ResultTaskID=305301cc-…` 明明在场 | **变过**（2→1）：生产腿的归因消失 |
| **ext-D**（本程自加） | `drop-test:TestCompressionTraceSilentWhenNothingFoldableOverThreshold`（摘 N1，即 AC#1 那枚判据用例） `LANDED (internal/agent/compress_trace_test.go): removed … (lines 426-464)` | **3** | **2** | 与 ext-A 逐枚同形 | **不变** |

⇒ 与实现件 §4.2 那张表的三枚数（2／0／1）**逐枚相同**，且**痕总行数恒 3** 这一条也复到了（摘的是归因，不是那条痕本身）。
⇒ 但**这不是"抵账"**：实现件那三发数的是测试进程 stderr 上的 `task=`，本程数的是**盘上 JSONL 里的 `"task":"`**，两把尺不同形、同结论。前一程 §8 第 3 条写的是"没复算"，所以这一格今天**第一次**有到盘读数。

### 3.2 复算途中量出来的一处**仪器形状**（不写清楚下一位会误用这枚探针）

派单说"复用前一程那枚到盘探针接 `drop-attr`／`unloop`"。本程照做，并且**当场发现那枚探针看不见 N2**：

```
ext-C（摘 N2）同一棵树上并排两枚探针：
  复用那枚两腿探针 zzaccept153_ondisk_test.go   --- PASS   两行、第二行仍带 "task":"ed57ab98-…"
  本程三腿探针     zzprobe155_ext_ondisk_test.go            trace_lines=3 with_task=1
```

**为什么**：`zzaccept153_ondisk_test.go` 自己直接 `withTraceTask(context.Background(), id)` 打标，**不经过 `Loop`**；`unloop` 改的是 `loop.go` 那一发调用点 ⇒ 对它**结构上不可见**。
⇒ 所以"只复用已有仪器"这一条**不足以**答 ext-C：本程**在仓外快照里新增一枚三腿探针**（第三腿＝`newHarness(t,"text-reply", …)` 真跑一次 `Loop.Run`），复用那枚**一字未改、留在树里继续跑**（它在 ext-B 会红、ext-C 不红——这枚"半瞎"形状本身就是读数）。
⚠ 边界：**没有**因此得到任何"改交付码"的需要 ⇒ **本程零生产码改动**（§3 全程只往快照里加 `_test.go`，工作树 `internal/**` 一枚字节没动，见 §5）。

### 3.3 正控与还原（AC#4 第③条：开测前先跑一条已知会红的）

| # | 那一发 | 读数 |
|---|---|---|
| 正控①（这棵树活着） | 无变异 `-run TestCompressionTrace -count=1 -v` | **7 RUN／0 FAIL**，七枚 PASS 名册逐枚在场 |
| 正控②（尺有牙） | 同一把尺施 `drop-attr` | **2 枚红**：`--- FAIL: TestCompressionTraceCarriesTheOwningTaskID`、`--- FAIL: TestCompressionTraceNeverInventsATaskID` ⇒ 红名与实现件 X5 逐字相同；复用那枚到盘探针**也红**：`zzaccept153_ondisk_test.go:70: tagged record lost its task on disk: want "a465ea9d-…" in {…}` |
| 顺带复现句①那一形 | ext-D 同一棵树上再施 `m5`（守卫→`if c.Need(hist)`） | RUN 从 7 掉到 **6**（少的正是被摘掉的证人）、`ok github.com/CarlosShao/wisp/internal/agent 0.056s` **全绿** ⇒ 这就是"摘证人⇒逃逸回来"（实现件 X2 的本程版，独立跑到） |
| 还原自证 | 每发之前 `restore`、收口逐枚 `cmp` | `RESTORED IDENTICAL compress.go / loop.go / compress_trace_test.go`；再 `git hash-object` 三枚＝`4fcd9a32a500…`／`3ea1fb8df38a…`／`2357dca86f67…`＝**与 `6de3d1c5` 的 blob 号逐字相同** |
| 仪器没跑到的第三种形状 | `third_party/**` tracked＝0（本程现读 `git ls-files third_party \| wc -l` ＝ **0**） | 本程**一行 `cmd/wisp` 都没跑**、那三枚 dll **没拷** ⇒ 这一格没有"8 条 no native DLLs 当成读数"的问题；`internal/agent` 包不依赖它们。**这一条同时记进 §5"没测"**：宿主那条真链本程未证。 |
| 版本现读 | `go version` ＝ `go1.27.1 windows/amd64`（本程现取，未背前一程的数） | 未用 `-overlay`、未用 `-cover*`（两坑都在 AC#4 第③条里，本程绕开了） |
