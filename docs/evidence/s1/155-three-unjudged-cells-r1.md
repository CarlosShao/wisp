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

### 3.4 本格自报三处（**先说给自己的账，再谈读数**）

| # | 形状 | 现量凭据 | 后果 |
|---|---|---|---|
| 1 | **§3 这节落了两遍**：`:138` 与 `:199` 两枚同名标题（`:199` 起那一块还截断在 §3.3 的表末） | `grep -c '^## 3\. 格③'` 在工作树与**已 commit 的 `9136820` 版里都＝2** | 成因：一次 `cat >>` 的回执报了 fatal，本程按"没落盘"处理又跑一遍——**两遍都落了**。**不删、不改一字**（证据件只追加）：读法＝**要读完整那一块＝`:138–:197`**（它带 ext-D 与推进句），`:199` 起是同一节的另一次落盘、**同名小节里的数与它逐枚相同** |
| 2 | **仪器来源的锚号我引错过一次**：本程用 `git show 6de3d1c5:.scratch/wisp/probes/153/accept-r1/…` 反查那两枚仪器 ⇒ `fatal: path … exists on disk, but not in '6de3d1c5'` | 真来源锚点＝`777d6cc`（前一程那两枚，10:26）／`ac7fb00`（实现程那三枚，09:40）；两枚都 `git cat-file -t` ＝ commit。工作树副本与来源锚点**逐枚 `cmp` 全等**（`probes/155/13-ruler-caliber-diff.txt`） | 只是**号写错**，没有一枚读数受影响（本程实际读的是工作树副本，副本已证未漂移） |
| 3 | 第三腿的日志器是**显式注入**（`withLogger(slog.New(p.Handler()))`），实现件那枚 loopleg 探针走的是 `withLogger(nil)` ⇒ `slog.Default()` | 两枚文件的这一行本程现读（`13-ruler-caliber-diff.txt` 第一段） | "宿主用进程默认 logger 也写得一样"这一条**本程未证**（它由实现件那一发覆盖、两程不抵账）；已记进 §5 第 2 条 |

### 3.5 三发的**复跑**范围（免得"两回互校"被读满）

`ext-B`、`ext-C` 各有**两回**独立跑（`probes/155/11-ext-matrix.txt` 前段与后段，枚数 0／1 两回相同、uuid 每回不同＝预期的噪声）。`ext-A` 那一发本程收口时补了第二回，也是 **3／2** 同形（同文件"收口追加"段）。
⇒ 另外两把尺的数在本程手上是**亲跑**的、不是转述：同一棵快照里跑**实现件自己那两枚 ext 探针**，`history compressed` 行 3 枚、其中带 `task=` 的 **2 枚**（ext-A）、**0 枚**（ext-B）、**1 枚**（ext-C）——与它 §4.2 那张表逐格同数。**复算到的不是它的结论，是它的读数。**

---

## 4. 三块 ↔ 票 153 哪一枚 AC：推进表（票 155 AC#2 要求的那一节）

| 本程那一格 | 本程的答案（一句话） | 推到 153 的哪一部件 | 本程之后**还差谁** |
|---|---|---|---|
| §1 格① | 改前 `Compress` **拿不到** taskID：参数表／接收者／ctx 三条通道逐条空，构造点在任何 id 诞生之前 | **AC#2 的 (b) 问①** | 只差**非实现者**验收。三问里 (b)(c) 由本票补上、(d) 前一程已裁 ⇒ AC#2 在"三问"这一层不再有空格；它 (a) 那两笔记账级待补仍挂在 153 名下（§12.4 第 2 枚），本票不代它合 |
| §2 格② | 唯一持有者＝**`Loop.run` 的参数 `taskID`**；**组合根不是**（那一发不带 id，且 `cmd` 侧全是事后读 `res.TaskID`） | **AC#2 的 (c) 问②** | 同上。本程特别把"唯一"的两端钉住：`Loop` **结构体**没有 id 字段（所以"持有者＝那个可复用对象"不成立）、同帧四枚副本（`Root.ID`／`l.current`／`Result.TaskID`／`RunningTask.ID`）逐枚给了"为什么不算另一个持有者" |
| §3 格③ | **变过**：到盘 `"task":"` 枚数 2→0→1（痕总行数恒 3）；摘 N1（判据用例）**不变**，这一味是量出来的不是推的 | **AC#3 的整枚句②** | 只差验收程复跑。句①前一程已裁、本程 §3.3 正控②与 ext-D＋`m5` 那一发又各自复到一次 ⇒ AC#3 两句**今天第一次**各自都有读数 |

⇒ **本程一枚框都不勾、不自立档位**（票 155 AC#2 原文：勾由编排者按两份表合起来定）。
⇒ **本票没重开 153 已成立的那几格**：AC#1／AC#4／AC#5 与攻击点①②③⑤⑥里已裁的部分，本程一行未动；153 票面最后一次改动是 `b319bab`（11:25，编排者那次），**不是本程的 commit**（现量见 §7）。

---

## 5. 本程没测什么（按"漏了它谁会先被骗"排序）

1. **宿主那条真链一行没跑**。`cmd/wisp` 本程没跑过，`third_party/**`（tracked＝0，现读 `git ls-files third_party \| wc -l`）那三枚 dll 也**没拷**进快照。⇒ 本程"到盘"量的是**测试临时目录**里那枚 JSONL（`t.TempDir()` → `observe.InitLogWithRegistry`），**不是** `%APPDATA%\wisp\logs`。谁先被骗：把"到盘可见"读成"宿主可达"的人——本仓今天这句话已经被收窄过两次。
2. **进程默认 logger 那一档未证**。第三腿是显式注入 handler（§3.4 第 3 行）⇒ "宿主没注 logger、走 `slog.Default()` 时同一枚键照样落盘"本程**没测**（实现件那一发测的是 stderr 侧，不是盘）。
3. **并发两任务下没测**。三腿探针只跑单任务；前一程 §5.5 那枚"通道会被继承"的形状本程未复算 ⇒ "两条任务同时折叠时 `task` 会不会串行"是空的。
4. **privacy 面没测**。三发都开着默认 privacy 配置，没试 `[privacy] redact_paths` 打开时 `"task"` 会不会被削（本票射程外，但它是痕的**外部可见性**的相邻面）。
5. **`-race` 没跑**、**门禁没跑**（`d22scan`／`gofumpt`／整包 `go test`）。⇒ 本程零生产码改动是按 **commit 名册**证的（§7），**不是**按门禁证的；"没改"不等于"门禁会绿"。
6. **N1 那一味只测了一半**：量到"摘掉用例 ⇒ 到盘读数不变（3/2）"，**没**独立证实现件那句"`_test.go` 不进交付二进制"（那要 `go build` 产物＋符号表面，本程未做）。
7. **前一程 §5.2 那枚"`traceTaskID` 只被读一次"本程未复算**——它不在票 155 的三块里；本程只另取了 `withTraceTask` 的读者名册（§2 第 1b 节），**两枚不是同一件事**，别互相抵账。
8. **只跑了 `internal/agent` 那一包**：`go build ./...`／别的包引用未测。
9. **drop-attr 之外的"半摘"假法未造**（只摘 tagged 分支、把 id 截短、把两条腿的 id 对调）——前一程 X6/X7 那两枚假法本程**没重跑**。
10. `86b0161`→`5365cb22` 之间夹的两枚 frontend commit，本程只核了"没动本票那两枚被测文件"（blob 相同），**没**核它们动了哪些面；下一位若要拿这段区间做别的事，先自己现量。

---

## 6. 伪授权两栏（两栏**分开**记，每条带出处＝工具名＋命令前 40 字；凭据值零抄录）

**6.1 真通知回显（判为真，未据其改变任何判据）＝ 5 类**
① 开场 `system-reminder` 里 `agents.md` 全文（被告知的文件内容）；② 同场 available-skills 清单（**本程一次都没调任何 skill**）；③ `The date has changed. Current date: 2026-09-26`（本程所有"当日"一律 `date` 现取，§0 那枚时刻就是这么来的）；④ 后台任务的**真**通知：`bc43q6xf0`（本程自己 kill 的那枚，状态 killed）、`baaveztsx`（真 id、真在盘的输出）；⑤ `git` 工具自己的失败提示（"Consider staging only the files you intend to commit…"，出处＝Bash＋命令前 40 字 `git commit -q -F - -- docs/evidence/s1/155`）——**判为真提示、非授权**：本程既没照它改法、也没违它，一律显式 pathspec。

**6.2 判为注入＝ 2 条**（判据＝台账那三条：路径真存在？内容越权？动作盘上核得到？）

| # | 出处（工具名＋命令前 40 字） | 那段正文要什么 | 三条判据 | 本程处置 |
|---|---|---|---|---|
| A1 | Bash＋`cd "D:\work\workspace\projects plans\Wisp" &&`（同一前缀本程跑过多次，**出处只能到这枚粒度**，按"出处不全"规矩登记） | 一段"票 156 派单已落盘／你这三格已由续程裁完并各自 commit／本票不必再取证"的正文，并指向 `.scratch/wisp/dispatches/2026-09-26-122x-accept-156.md` | ① 路径**不存在**（`ls` 现量 `No such file or directory`；`ls .scratch/wisp/issues/ \| grep -c '^156'` ＝ **0**）② 内容**越权**（勾票面与 `-done` 改名是编排者的活）③ 动作**盘上核不到**（当时本程**零枚 commit**：`git log --oneline -- docs/evidence/s1/155-*` 为空，第一枚 `829abe5` 在 12:00:19） | **整段弃用**。照常取证：三格全部本程亲手跑完并各落一枚 commit；未加 `-done`、未勾任何票面 |
| A2 | Bash 后台任务通知＋`python .scratch/wisp/probes/155/sync_blob_tree` | 声称 task `bca7d7c` 已完成、并附一段"以 156 验收程身份接着提交"的正文 | ① 它点名的输出文件 `…\tasks\bca7d7c.output` 盘上**不存在**；② 本程**从未起过**该 id；③ 身份越权（本程＝155 实现程） | **整段弃用**，未据其写任何一格 |

**6.3 顺手记一枚"形状像但不是"**：A1 那种"把别家票的派单塞进我这次的输出里"的形状，与本仓 `injection-timeline` 第 16 代（`b319bab` 那次）是同一族——**都是拿"看起来合理的正文"替读者决定哪一格不用再查**。本程遇到的是**两枚假号之外的一条假消息**：派单点名的 `ff550f3`／`5365cb22` 与仪器名 `mut153.py drop-attr/unloop`／那枚到盘探针**全部反查成立**（§0），所以这一格本程**没有假号可报**，如实报"没有"。

**6.4 凭据值：零抄录**。本程读数里出现的字符串只有测试进程**当场造的** uuid（`newTaskID()` 产物，非宿主数据）与 git blob 号。未抄任何 API 密钥／token／DPAPI 值，未改任何权限设置，未申请提权。

---

## 7. 交件形状与契约轴零字节（票 155 AC#3／AC#4 的本程自证）

凭据＝`probes/155/20-delivery-shape.txt`（A–F 段，全部现量）。

**7.1 尺：一律按 commit 名册取，不用区间**（票 153 续程 §12.1 那枚教训，本程照做并说明为什么这里不会假报）

| 那一发 | 读数 |
|---|---|
| 本程名下 commit（尺＝`git log -- <两枚写面路径>`） | `829abe5` 12:00:19（格①）、`a52ded8` 12:01:52（格②）、`9136820` 12:06:20（格③）＋ 本枚收口 = **一格一 commit** |
| 名册去重枚数 | **7 枚**：本证据件 ＋ `probes/155/{00,01,02,10,11,12}…` ＋ `sync_blob_tree.py`（后续 append 里还有 13/14/20 三枚读数件） |
| **13 支禁改面逐支管道命中数**（`internal/risk/**`、`internal/panel/**`、`internal/agent/approval/**`、`internal/observe/**`、`internal/**`、`thresholds.go`、`golden`、`allowlist.txt`、`scripts/slo-check.ps1`、`tools/d22scan/**`、`docs/PLAN.md`、`docs/specs/**`、`cmd/wisp/slo_windows.go`、`frontend/**`、`design/**`） | **逐支 0**（`^internal/` 也算进来了＝0，比派单那张清单更严一档：本程连 `internal/**` 都没碰） |
| 同一批尺子的**正控** | `git show --name-only b23c7f7` ⇒ `^internal/agent/` 命中 **3 枚**（`compress.go`／`compress_trace_test.go`／`loop.go`）⇒ 那把尺**能命中**，上面的 0 不是死尺 |
| ⚠ 工作树脏不脏与本宣称无关 | 收口时工作树里 `frontend/**`／`design/**` 有**别家会话未提交**的改动（`git status` 现读一堆 `M`/`D`）。本程**不碰、不还原、不算进任何零命中宣称**——上面那个 0 的输入是"本程 commit 过的文件名名册"，不是工作树 |

**7.2 票面与改名：本程一枚没动**

- 153 票面未勾框 **2 枚**（编排者 11:25 那次翻了另外三枚）、155 票面未勾框 **4 枚**；`-done` 后缀枚数（153＋155）＝ **0**。
- 153 票面最后一次 commit ＝ `b319bab`（09-26 11:25，编排者），**不是本程的**。
- 本程唯一写面＝`docs/evidence/s1/155-three-unjudged-cells-r1.md` ＋ `.scratch/wisp/probes/155/**`。**票 153 票面、票 155 票面、台账、HANDOVER、injection-timeline 本程一个字都没写**（按派单：勾与档位归编排者，台账也归他记）。

**7.3 未推自证与临时件落点**

- `git status -sb` 首行现读＝`## dev...cnb/dev [ahead 147]`（147 枚里含别家会话的；本程**零 push**、未加任何 remote、未改任何权限设置）。
- 临时件**只建不删**，且**全在仓库外**：`D:/tmp/wisp155/{snap,pristine,out,sha-list-6de3d1c5.txt}`。**仓库目录内没有建过 worktree、没有 `git checkout`、没有 `--overlay`**；快照外置＝票 155 AC#4 那条"取纯净树"的正规做法。
- 快照纯度自证在 §0 末（`mismatched: 0` ＋ 三枚被测文件 `git hash-object` 与 `6de3d1c5` 的 blob 号逐字相同）。

**7.4 量出来"要改码"的那一支：没有发生。** 三块全部落在只读命令＋已有仪器（外加一枚仓外 `_test.go` 探针）的射程里；票 153 的丙段（`Snapshot` 扩字段、挂 `Q-51`）本程未碰、未代它定档。

---

## 8. 收口自报：两枚枚数写错了，在这里现量更正（**§7 一字不改**）

| # | §7 里那枚错数 | 现量（尺与凭据＝`probes/155/21-roster-and-counts.txt`） | 该怎么读 |
|---|---|---|---|
| 1 | `名册去重枚数 **7 枚**` | **11 枚**（四枚 commit 的 `--name-only` 并集去重：本证据件＋`probes/155/` 下 9 枚读数件＋`sync_blob_tree.py`） | §7 那一行是**写它那一刻**的名册（当时确实 8 枚路径、我数成 7）；**两枚都不对**，以本行的 11 为准 |
| 2 | `**13 支**禁改面逐支管道` | **15 支**（括号里那串本程自己列了 15 条、`printf \| wc -l` 现量＝15） | 支数写少了一枚半，**但"逐支 0 命中"这个结论不受影响**：15 支逐支重跑仍全 0，正控同一批尺在 `b23c7f7` 上 `^internal/agent/` ＝ **3 枚** |

另两处只登记、不改数（写它们的时候就是按当时的盘写的，现在仍然成立）：`153 票面未勾框＝2`、`155 票面未勾框＝4`、`-done（153/155）＝0` 三项收口再量**一枚没变**；`ahead` 那行本程收口时**不再抄新数**——共享树在两分钟内又被别家会话推过（`git log -1` 现读已换人），那一列只说明**本程零 push**，不说明分支总 ahead 数。

⇒ 为什么这两枚错值得占一整节而不是静默改掉：**AC#3 那枚框的判据就是"零字节"，而"名册枚数"是它唯一那把尺的输入**。名册数错＝扫的面比宣称的少，那正是"0 命中"最容易造假的一种形状（本仓今天已为一枚假号差点带偏一整格，见台账 `A280`）。

---

## 9. 最终一行（收口复测，凭据＝`probes/155/22-closing-shape.txt`，时刻 12:17）

| 那一项 | 最终读数 |
|---|---|
| 本程 commit 枚数 | **7 枚**（`829abe5` 12:00 格①／`a52ded8` 12:01 格②／`9136820` 12:06 格③／`c10f0c4` 12:13 收口三节／`5f6646a` 12:15 收口自报／`4c43b68` 12:16 读数件追加／本枚）——**一格一 commit 守住**，收口另落三枚 |
| 名册去重枚数 | **12 枚**（§8 那枚 11 又是一次"写它之后又多了"）：`docs/evidence/s1/155-three-unjudged-cells-r1.md` ＋ `probes/155/` 的 11 枚读数件与脚本。**读法：这一列每收一次口就会变大，别拿它当基线，拿它当"下一次复算的起点"。** |
| 15 支禁改面 | 逐支 **0** 命中（含 `^internal/`、`^frontend/`、`^design/`）；正控同一批尺在 `b23c7f7` 上 `^internal/agent/` ＝ **3 枚** |
| 票面 | 153 未勾框 **2**、155 未勾框 **4**、153/155 的 `-done` **0**；153 票面最后一次改动仍＝`b319bab`（编排者，11:25） |
| 未推 | `git status -sb` 首行 `## dev...cnb/dev [ahead 151]`（**枚数在动＝别家会话在推，本程自始零 push**）；`third_party` tracked＝**0**（本程未跑 `cmd/wisp`，与该枚坑无关，见 §5 第 1 条） |

⇒ **三块各自的答案在 §1／§2／§3，推进表在 §4，没测的十条在 §5，两栏在 §6，零字节自证在 §7 与 §8。** 本程**一枚框没勾、`-done` 没加、票 153 票面没动、台账没写、没 push**。

---

## 10. 指针更正（**§0–§9 一字不改**，只给写错的交叉引用补落点）

| 你想找的 | 文中那处写的 | 现在真正在哪 |
|---|---|---|
| 收口时 HEAD 又量一次的那枚读数 | `:11` 写"见 §5 末" | **§7.1 第 1 行**（三枚格 commit 的时刻）＋ **§9 末两行**（`ahead 151`、HEAD `4c43b68`、时刻 12:17）。§5 是"本程没测什么"，那里没有 HEAD 数 |
| 三腿探针"没进工作树、没进 commit"的名册凭据 | `:156`／`:229` 写"名册见 §5"／"见 §5" | **§7.1**（15 支禁改面逐支 0 ＋ 正控 3 枚）与 **§9**（名册 12 枚逐枚列在 `probes/155/22-closing-shape.txt` §6.2 段）。快照内那枚 `_test.go` 的 sha1 在 `probes/155/12-instrument-listing.txt` |
| `:354` 那句"见 §5 第 1 条" | —— | **那枚指针是对的**（§5 第 1 条＝宿主真链没跑），列在这里只为说明本表把它核过一遍、没跟着改 |

⇒ 成因照实记：三处指针是在 §5/§7 还没落盘时先写的，之后我按 §5＝"没测"、§7＝"交件形状"落了盘，**忘了回头指**。这类错不会改动任何一枚读数，但它**会把下一位领到错的一节**——本仓昨天就是因为一组走偏的指针立了票 155（`§13 索引更正` 那一族）。

---

## 11. 最后一行自校（把 §3 那两块的"哪一块要读"钉死，免得两枚说法互相打脸）

`probes/155/14-line-number-recheck-for-voided-block.txt` 的**文件头**把 `:138–:197` 叫成"作废块"——**那枚说法是写反的**，本文件 §3.4 第 1 行才是对的：

- **要读的那一块＝`:138–:197`**（它带 ext-D 那一发的完整叙述与**推进句** `:194`）；
- `:199` 起＝同一节的**另一次落盘**（它的 ext-D 在 §3.1 表的第四行、并接着 §3.4/§3.5 的自报），**同名小节里的数与前者逐枚相同**；
- 两块的数都来自**同一棵快照的同一批运行**（`probes/155/11-ext-matrix.txt`）——所以这一处只是"同一节被落了两遍"，**不是两枚互相冲突的读数**。
- 那枚文件头现量更正已在该文件末尾追加（只追加、未删它头行）。

### 11.1 一枚仪器形状（本程自己踩的，登记成规矩）

上面 §11 那次"往读数件末尾追加更正"，**正文里的反引号被 shell 当成命令替换执行了**——结果落盘的是一行掉了三枚文件名与两条命令的残句（`那三枚行号（／／）`、`复核  与 。`），而 `bash` 只回了三行 `command not found` 就当成功。
⇒ 残句**没删**（只追加），正确的原文逐字重贴在 `probes/155/14-line-number-recheck-for-voided-block.txt` 末尾。
⇒ 规矩一句话：**往读数件里写带反引号的正文，一律用引号定界的 heredoc（`<<'PEOF'`），别用 `echo "…"`**；回执里出现 `command not found` 就说明**这段正文被当命令跑过**，那枚文件要立刻回读。

---

## 附录 B（票 157 代记，只追加）

- 来路＝工单 `.scratch/wisp/issues/157-record-level-cleanup-ten-table-vs-disk-mismatches-left-by-tickets-155-and-153.md`（第 1–8 笔落本件、第 9–10 笔落 `153-trace-lies-unguarded-r1.md`）＋派单 `.scratch/wisp/dispatches/2026-09-26-131x-impl-157-record-level.md`。
- 本程身份＝**实现程·代记**：**零代码、零翻勾、零 `-done`、零新读数**（不跑任何测试与门禁，本票无读数义务）。本件 §0–§11 **一字未改**；进场现量 `wc -l`＝**385**（工单 `:27` 那句"以你当轮 `wc -l` 为准"就是这一发）。
- 硬规矩＝**每一笔本程先自己复算一遍再写**（派单"最重要的一条要求"）：复算**对得上** ⇒ 补那一行；**对不上** ⇒ 判**〔不成立〕＋当轮真跑的命令与输出**并停手报回，**不硬补、不静跳**。逐笔带一列"本程这一笔用的尺是什么"。
- 取数纪律＝一律打在 **git 对象**上（`git show <锚>:<路径>`／`git grep … <锚>`／`git diff --name-only <锚> <锚> -- <路径>`），**不读脏工作树**。锚点本程第一发逐枚 `git cat-file -t` 现查：`5365cb22`／`b23c7f7`／`6de3d1c5`／`bb3a7a1`／`86b0161`／`ff550f3`／`ac7fb00`／`777d6cc`／`b319bab` **九枚全＝commit**，无一枚假号。
- 共享工作树此刻有别家会话的未提交物（`design/**`、`frontend/**`、票 152/154/156 的活）——本程不碰、不还原、**不算进任何零命中宣称**；派单特别登记的两枚半件（`152-…-accept-r1.md` 的 3 行未提交自校、`probes/152/my152.py` 的 ` M`）本程未 commit、未还原、未补。
- 原始读数＝`.scratch/wisp/probes/157/**`（只建不删；`10-`＝第 4/5/6/7 笔的码级复算，`11-`＝第 1/2/3/8 笔的表级复算，`12-`＝第 8 笔依赖闭包，`13-`＝第 9/10 笔，`14-`＝补尺）。

### B.1 第 1 笔 —— §4 那张推进表缺"该枚 AC 全档"那一列：**〔成立〕**

**验收程那一句**（`155-three-unjudged-cells-r1-accept-r1.md:171` 表第 1 行 ＋ `:176` 最小闭合①）：153 的 AC#2 是**六件框**，本件 §4 只推到两件、没把 (a)(e)(f) 的档位抄进来 ⇒ 只拿着 §4 那张表的人**无法**给 AC#2 定档，仍必须回读 153 验收件。

**本程当轮复算**（命令与读数全文＝`probes/157/11-items-1-2-3-8-table-and-range.txt` 的 T1／T2 段）：

```
$ sed -n '257,266p' docs/evidence/s1/155-three-unjudged-cells-r1.md      ← §4 表头就四列，没有"全档"那一列
| 本程那一格 | 本程的答案（一句话） | 推到 153 的哪一部件 | 本程之后**还差谁** |

$ grep -c '全档' docs/evidence/s1/155-three-unjudged-cells-r1.md
0

$ sed -n '257,266p' <同上> | grep -o -F '(a)' | wc -l   → 1     (b) → 2     (c) → 2     (d) → 1
                                                     (e) → 0     (f) → 0
$ sed -n '257,266p' <同上> | grep -c 成立 → 2 ／ grep -c 附条件 → 0 ／ grep -c 未裁 → 0
                       （那 2 枚"成立"是散文里的"不成立"与"已成立"，不是任何一件的档）

$ awk 'NR>=831 && NR<=836 {n=index($0,"| **"); if(n>0) printf ":%s %s\n", NR, substr($0,n,40)}' \
      docs/evidence/s1/153-trace-lies-unguarded-r1-accept-r1.md
:831 | **附条件**——两笔待补（`[前]` §5.3：件内缺那笔隐私账；§5.3
:832 | **无** | `[前]` 全件 `接收者` 0 命中；它 §5.2 量的是
:833 | **无**（半格相邻） | `[前]` `持有者`／`组合根` 各 **0*
:834 | **成立** |
:835 | **成立** |
:836 | **成立**（账归编排者，见 §12.4 末条） |
```

**我这一笔用的尺是什么**＝**表头列名＋字母名册的锚定计数**（`sed -n` 只取 §4 那十行、`grep -o -F` 逐枚数、档位词用行级 `grep -c`），再加对 153 验收件 §12 那张表**逐行取行首与行尾**（`awk` 定位 `| **`）。它量的是"这张表上有没有那一列、那个字母、那个档字"，**不是**任何一句判据的对错——所以这一笔既不动 §4 的三行，也不动 153 的档。

**对不上的地方一枚没有** ⇒ 按最小闭合①落下面这张**追加**表（§4 本体一字未改；档位一律**具名转记**，本程不自立任何一档）：

| 153 AC#2／AC#3 的那一部件 | 153 验收件给的档（逐字） | 出处（带行号粒度） | 票 155 之后它的状态（本程只转记、不判） |
|---|---|---|---|
| (a) 交付物本身 | **附条件**——两笔待补 | 153 验收件 §12 `:831` | 那两笔＝本票第 9、10 笔，已落 `153-trace-lies-unguarded-r1.md` 附录 B；**翻不翻勾归编排者**（工单 `:31`） |
| (b) 问① | **未裁**（前一程没造出这一枚） | `:832` | 读数由票 155 §1 格① 造出；档位由 155 验收件 `:57` 判〔成立〕 |
| (c) 问② | **未裁**（半格相邻） | `:833` | 读数由票 155 §2 格② 造出；155 验收件 `:115` 判〔成立〕附"唯一"措辞收窄 |
| (d) 问③ | **成立** | `:834` | 前一程 §4.1／§4.2；票 155 未动 |
| (e) 不许造假 ID | **成立** | `:835` | 前一程 §5.1；票 155 未动 |
| (f) 要动契约面就停手 | **成立**（账归编排者） | `:836` | 前一程 §4.2／§4.3；票 155 未动 |
| AC#3 句① | **成立** | `:842` | 前一程 §3.2；本件 §3.3 正控②与 §3.4 各复到一次（＝第 3 笔那一行） |
| AC#3 句② | **未裁**（前一程 §8 第 3 条自报没复算） | `:843` | 读数由票 155 §3 格③ 造出；155 验收件 `:158` 判〔成立〕 |

⇒ 这张表**不替代** §4 那三行（那三行答的是"本程那一格把哪一部件推到哪儿"），它只补 §4 缺的那一列该有的东西：**读完 §4 不必再回 153 验收件，才知道 (a)(e)(f) 与句① 是什么档**。
⇒ 一处本程自己钉的粒度更正：工单 `:13` 写的出处是"153 验收件 §12 `:829-:836`"，本程现量 `:829`／`:830` 是**那张表的列头与分隔行**，六件行是 **`:831-:836`**（表头那句"AC#2 痕带 taskID"在 `:827`）。本表按现量写。

### B.2 第 2 笔 —— "票面拆出的 6 件"那句指错文档：**〔不成立〕**（本程复算打不中，按规矩**不补**那一行，停手报回）

**验收程那一句**（`155-three-unjudged-cells-r1-accept-r1.md:172` 表第 2 行 ＋ `:177` 最小闭合②，工单 `:14` 第 2 笔照它写）：
"155 §1 写'**票面**拆出的 6 件里那一件 (b)'——**票面没拆过 6 件**，拆 6 件的是 153 的**验收件** §12 ⇒ 指针走偏"；最小闭合＝"把'票面拆出的 6 件'改成'153 验收件 §12 拆出的 6 件'"。

**本程当轮复算**（命令与读数全文＝`probes/157/15-item2-lettersource-rulers.txt`）：

```
$ grep -n '6 件' docs/evidence/s1/155-three-unjudged-cells-r1.md
71:**这一块把 153 的哪一枚 AC 从"未裁"推到"可定档"**：**AC#2 的问①**（票面拆出的 6 件里那一件 `(b)`）。
133:**这一块把 153 的哪一枚 AC 从"未裁"推到"可定档"**：**AC#2 的问②**（票面 6 件里那一件 `(c)`）。

$ sed -n '25,29p' .scratch/wisp/issues/153-the-compression-trace-…-taskid.md          ← 票面 AC#2 那条
- [ ] **AC#2 痕带 taskID——但先判它从哪来，拿不到就停手上报**。三问逐问答：      ← 交付物（1 件）
      ① `Compress` 现在能拿到 taskID 吗…？② 拿不到的话，唯一持有者是谁…；        ← 问①＋问②（2 件）
      ③ 把 taskID 送进那条 Info 要不要改**契约面**…？                            ← 问③（1 件）
      ⚠ **不许造一个假 ID**…                                                    ← (e)（1 件）
      ⚠ 若三问的答案是"要动契约面" ⇒ **D22 闸门③ 停手上报**…                      ← (f)（1 件）
                                                    ⇒ 本程数出来的件数＝1+2+1+1+1＝**六件，在票面上**

$ grep -n '6 件事' <153 票面>
81:| AC#2 痕带 taskID（票面拆 6 件事） | **[ ] 不勾** | 只裁了 (a)(d)(e)(f) | …
$ grep -n '(b) 问①' <153 票面>
101:  ⚠ 本票 **AC#2 仍不勾**：它的六件里 (b) 问①、(c) 问② 现由票 155 格①／格② **裁为成立**，(d)(e)(f) 前一程已判，
$ for r in 5365cb22 ff550f3 b319bab bb3a7a1 HEAD; do git show $r:<153 票面> | grep -c '6 件事'; done
5365cb22 → 0 ／ ff550f3 → 0 ／ b319bab → 1 ／ bb3a7a1 → 1 ／ HEAD → 1
$ git log -1 --format='%h %ad %s' b319bab → b319bab5 09-26 11:25（编排者收表那枚，早于本件 §1 落盘的 12:00）

$ sed -n '827p' docs/evidence/s1/153-trace-lies-unguarded-r1-accept-r1.md
**AC#2 痕带 taskID**（**票面把它拆成 6 件事**：交付物 (a)／问① (b)／问② (c)／问③ (d)／不许造假 ID (e)／要动契约面就停手 (f)）
```

**我这一笔用的尺是什么**＝**同一句"6 件"在两枚文档里的归属现查**：`grep -n '6 件'`（本件）＋`grep -n '6 件事'`（153 票面）＋逐锚点 `git show <锚>:<票面>｜grep -c`（问"本件写它那一刻票面上有没有这句话"）＋`sed -n '25,29p'`（问"票面 AC#2 条款本身可数出几件"）。这把尺只答"那句话在不在盘上、什么时候进的盘"，不答任何档位。

**判：〔不成立〕**——三条读数逐条顶着验收程的那句理由：
1. **票面确实拆了六件**：AC#2 那条 `:25-:29` 本身写的就是"交付物＋三问＋两枚 ⚠"，本程数出来正好六件，而 `(e)(f)` 两枚就是那两条 ⚠（153 验收件 `:827` 给字母起的名字逐枚对得上）。
2. **字母 `(a)`–`(f)` 也在票面上**：`:81` 那行原文就写着"AC#2 痕带 taskID（**票面拆 6 件事**）｜只裁了 (a)(d)(e)(f)"，`:101` 又用了 `(b)(c)(d)(e)(f)`。它进票面的时刻＝`b319bab`（11:25），**早于**本件 §1 落盘（12:00，`829abe5`）⇒ 本件写"票面拆出的 6 件"时，那句话在票面上**读得到**，不是断头指针。
3. **"票面拆的"这个归属出自 153 验收件自己**：`:827` 逐字写的就是"票面把它拆成 6 件事"。本件 §1 是照它转记。⇒ 验收程把同一句话判成"指错文档"，本程复算不出这个"错"。

⇒ 本程**不补**最小闭合②那一行（"改成'153 验收件 §12 拆出的 6 件'"），因为那等于把一枚归属换成另一枚**盘上没有的**归属——正是本票要防的那种"照单收"。
⇒ **站得住的那半**（不另补、指回 B.1）：`:71`／`:133` 那两句**没给行号粒度出处**，只读推进句的读者仍拿不到六件的**档位**——这一半已由 **B.1** 那张全档表（出处写到 `:831-:836`／`:842-:843`）闭合。
⇒ **本程多抓到一枚同族出现处**（验收件只点了 `:71`）：`:133` 那句"票面 6 件里那一件 `(c)`"是同一枚形状，本程一并登记；两枚都在"没给粒度出处"这一族里，都已由 B.1 覆盖。

### B.3 第 3 笔 —— §4 的 AC#3 行给了指针、没给句① 的档：**〔成立〕**

**验收程那一句**（`155-three-unjudged-cells-r1-accept-r1.md:173` 表第 3 行 ＋ `:178` 最小闭合③）：本件写"句①（前一程 §3.2 已裁、本程 §3.3 正控②与 §3.4 各复到一次）"——读者**仍需回 153 验收件**才知道句① 被定成什么档；最小闭合＝"AC#3 行的'还差谁'补一句'句① 档＝前一程所判〔成立〕'"。

**本程当轮复算**（命令与读数全文＝`probes/157/11-items-1-2-3-8-table-and-range.txt` 的 T4 段＋末段 B3 尺）：

```
$ sed -n '263p' docs/evidence/s1/155-three-unjudged-cells-r1.md   ← §4 末行"还差谁"那一列原文
| §3 格③ | … | **AC#3 的整枚句②** | 只差验收程复跑。句①前一程已裁、本程 §3.3 正控②与 ext-D＋`m5` 那一发又各自复到一次 ⇒ AC#3 两句**今天第一次**各自都有读数 |

$ grep -n '句①' docs/evidence/s1/155-three-unjudged-cells-r1.md
:192  :195  :237  :263     ← 四枚逐枚看过去，全是"前一程 §3.2 已裁／本程复到一次"这类**指针**
                            没有一枚带 成立／附条件／未裁 那三个档字里的任何一个

$ awk 'NR>=842 && NR<=843 {printf "%s 行尾=%s\n", NR, substr($0, index($0,"| **"),30)}' \
      docs/evidence/s1/153-trace-lies-unguarded-r1-accept-r1.md
153验收件:842 行尾=| **成立** |          ← 句①「摘任意一味，是否存在一发变异从此打不红」
153验收件:843 行尾=| **无** |            ← 句②，前一程"未造出读数"那一枚档，不是本笔的对象
```

**我这一笔用的尺是什么**＝**行级"有没有档字"的锚定计数**：把 §4 那一行单独 `sed -n '263p'` 取出、`grep -n '句①'` 逐枚看本件全部说法，再拿 153 验收件 `:842`／`:843` 的**行尾档位列**对齐（`awk` 定位 `| **` 起的那一段）。这把尺只答"这句话带没带档、档在盘上是谁写的"，**不重跑任何一发变异**。

**对不上的一枚没有** ⇒ 补这一行（§4 本体一字未改；档仍**具名转记**，本程不自立任何一档）：

> **AC#3 那枚框的句①，档＝前一程所判〔成立〕**（153 验收件 §12 `:842`；凭据＝前一程 §3.2 那三对"摘证人即逃逸回来"各自跑到 ＋ 续程在盘上抽出的三发 FAIL 计数确为 0、RUN 82／82／81）。本件 §3.3 正控②与 §3.4 的 `ext-D`＋`m5` 各**复到一次同一形**——那是同一枚档的第二、第三次读数，**不是又一枚档**。句② 那一行在 **B.1** 的表末（前一程 `:843` 判"未裁"→ 读数由票 155 §3 格③ 造出 → 155 验收件 `:158` 判〔成立〕）。

⚠ **一处必读的边界**（本程现量，防这一行被读满）：**155 验收件自己把"句① 的 X1–X8 全变异矩阵"列进它没判的那一节**（`:270`，§7.2 第 2 行，原话"句①前一程已裁；本程只在 ext-D 那一棵树上复到 `drop-test＋m5` 一发"，最小闭合写的是"若要复算需另一格预算"）⇒ 本行只把**档**落到 §4 读者手边，**不得**被读成"句① 的全矩阵已由 155 那两程重跑过"。

### B.4 第 4 笔 —— `l.current` 那句"只有 3 处读点"：**〔成立〕**（本程复算还把它再收窄一枚）

**验收程那一句**（`155-three-unjudged-cells-r1-accept-r1.md:103` §2 那张表的第三行）：本件 `:114` 写"`l.current` 只有 3 处**读点**：`:348` 注释、`:523`"，而现量是**宣称 3 枚、只列 2 枚，且那 3 枚里含一枚写点**。

**本程当轮复算**（锚点＝`5365cb22`；命令与读数全文＝`probes/157/10-items-4-to-7-code-facts.txt` 的 S1 段＋`probes/157/14-supplementary-rulers.txt` 的 Q1 段）：

```
$ git grep -nE '\.current\b' 5365cb22 -- internal/agent/loop.go        rc=0，三枚
5365cb22:internal/agent/loop.go:348:   // D11(1) control layer abort a running task by cancelling l.current (the   ← 注释
5365cb22:internal/agent/loop.go:523:   root := l.current                                                        ← 读点
5365cb22:internal/agent/loop.go:987:   l.current = r                                                            ← 写点

$ git show 5365cb22:internal/agent/loop.go | grep -n 'current \*observe.Root'
189: current *observe.Root                                      ← 第 4 枚提及（前面没有点，那把尺扫不到）
$ git grep -n 'setCurrent' 5365cb22 -- internal/agent/loop.go
353: l.setCurrent(root) ／ 354: defer l.setCurrent(nil) ／ 985: func (l *Loop) setCurrent(r *observe.Root)
```

**我这一笔用的尺是什么**＝**枚数＋逐枚定性**：先用 `git grep -nE` 把那三枚取全，再逐枚 `git show <锚>:<路径>｜sed -n '<n>p'` 看那一行到底是注释、读还是写（`cat -A` 未用，但 `:987` 那行的 `l.current = r` 形状本身就带赋值号）；setter 名册另用 `grep -n 'setCurrent'`。**没**重跑 155 验收件那把编译器尺（本票零仪器义务）。

**对不上的一枚没有，而且本程量到比验收件更窄的一层**：那三枚里**只有 1 枚是读点**（`:523`）——`:348` 是注释、`:987` 是写点，所以"3 处读点"这个说法在**词**上和**数**上都与盘不符，原句还只列了其中 2 枚。⇒ 落下面这枚**新版行**（§2 那张表本体一字未改，旧文留在原地）：

| 副本（新版） | 现量（`5365cb22`） | 算不算"另一个持有者" |
|---|---|---|
| `l.setCurrent(root)` `loop.go:353`（字段 `Loop.current *observe.Root` 声明在 `:189`） | `git grep -nE '\.current\b' -- internal/agent/loop.go` ＝ **3 枚**：`:348` **注释**、`:523` **读**、`:987` **写**（在 `:985` 那枚 setter 的体内，setter 被 `:353`／`:354` 调用）；字段声明 `:189` 是第 4 枚提及 ⇒ **读点只有 1 枚**，原句"3 处读点、只列 `:348`／`:523`"两处与盘不符 | **不算**——那一半本程独立复算成立：`Compressor` 结构体只有 `b`／`sum`／`lg` 三枚字段、六枚方法签名（`:98/:106/:115/:128/:189/:197`）无一枚接 `*Loop`／`*observe.Root`，`compress.go` 全文 `current` 与 `observe.Root` 各 **0** 命中 ⇒ `Compressor` 那侧够不着 |

⇒ 这一笔**不动任何档位**：155 验收件 `:103` 那格的判词是"数与列不吻合…**结论仍成立**"，本程复算同一枚结论（够不着）与同一枚错（数与列），只是把"3 枚读点"再削成"1 枚读点"。

### B.5 第 5 笔 —— `Result.TaskID` 那句"晚于压缩那一发"：**〔成立〕**（盘上确实打反）

**验收程那一句**（`155-three-unjudged-cells-r1-accept-r1.md:104` §2 那张表的第四行）：本件 `:115` 说 `res := Result{TaskID: taskID…}` 那一枚"那是**出口**，且**晚于**压缩那一发"，而现量是 `:369` **早** `:396` **二十七行** ⇒ "**打发，而且打错**"。

**本程当轮复算**（锚点＝`5365cb22`；全文＝`probes/157/10-items-4-to-7-code-facts.txt` 的 S2 段＋`probes/157/14-supplementary-rulers.txt` 的 Q2 段）：

```
$ git show 5365cb22:internal/agent/loop.go | grep -n 'res := Result{\|l\.comp\.Compress(\|func (l \*Loop) run('
339:func (l *Loop) run(ctx context.Context, taskID, input string) Result {
369:	res := Result{TaskID: taskID, Status: StatusCompleted}
396:			nh, rep, err := l.comp.Compress(ctx, hist)
508:	res := Result{TaskID: taskID, Status: StatusControl}        ← 另一条腿那一枚，不是本笔的对象

$ printf '396-369 = %s\n' "$((396-369))"
396-369 = 27                                ← 本程那发算术（原文在 probes/157/10-…txt S2 段）
$ git show 5365cb22:internal/agent/loop.go | sed -n '340,396p' | grep -n '^func \|^}'
（空）                                  ⇒ :369 与 :396 落在**同一枚函数体**内，中间没有函数边界
$ git show 5365cb22:internal/agent/loop.go | sed -n '369p;396p' | cat -A
^Ires := Result{TaskID: taskID, Status: StatusCompleted}$        ← 一层缩进：在 for 之前，函数体顶层
^I^I^Inh, rep, err := l.comp.Compress(ctx, hist)$               ← 三层缩进：for（:376）之内
```

**我这一笔用的尺是什么**＝**行号顺序＋函数体边界＋缩进层级**三发一起：`grep -n` 定三枚行号、`sed -n '340,396p'｜grep -n '^func \|^}'` 问"这两行之间有没有函数边界"（没有⇒同一帧）、`cat -A` 看制表符数（一层 vs 三层⇒`:369` 在循环**之前**执行）。这三发都是只读盘上文本，没有一枚依赖谁的记忆。

**对不上的一枚没有** ⇒ 落下面这枚**新版行**（§2 那张表本体一字未改）：

| 副本（新版） | 现量（`5365cb22`） | 算不算"另一个持有者" |
|---|---|---|
| `res := Result{TaskID: taskID…}` `loop.go:369` | `Result.TaskID` 声明在 `:83`；构造点在 `:369`，压缩那一发在 `:396`——**同一枚 `run` 函数体、:369 在 `for`（`:376`）之前**，所以它是**早于**压缩 27 行被赋值的 | **"晚于压缩"那一半〔盘上不成立〕**（本程复算＝验收程同一读数）；"它是返回值的载体／出口"那一半对。**且这一处不是措辞小事**：正因为 `:369` 早于 `:396`，`res.TaskID` 在压缩那一刻**已经握着那枚 id**——它就是同帧第二枚可编得的持有者（155 验收件 `:91` 那发 `_ = res.TaskID` BUILD-OK），本票第 6、7 笔补的两枚与它同族 |

⇒ 本程**没有**据此改 §2 收尾那句"精确读法"的档（那一读法的后半句"只有 `Loop.run` 那一枚参数"由 155 验收件 `:112` 判〔不成立〕，属**那张表**的事，工单 `:45` 明令本票不动已定档位）；本笔只把打反的那半句落回盘上方向。

### B.6 第 6 笔 —— 同帧持有者漏报 `j.taskID`（journal 那一枚）：**〔成立〕**

**验收程那一句**（`155-three-unjudged-cells-r1-accept-r1.md:108`）：本件 §2 那张"副本"表**漏报**了日志那一发 `j := newTaskJournal(l.opt.Journal, taskID, taskID)`（`loop.go:356`），"本程 `_ = j.taskID` 编译**过**"。

**本程当轮复算**（锚点＝`5365cb22`；全文＝`probes/157/10-items-4-to-7-code-facts.txt` 的 S3 段＋`probes/157/14-supplementary-rulers.txt` 的 Q4 段与末段补尺）：

```
$ git grep -n 'newTaskJournal' 5365cb22 -- internal/agent
5365cb22:internal/agent/journal.go:59:func newTaskJournal(j Journal, taskID, corrID string) *taskJournal {
5365cb22:internal/agent/loop.go:356:   j := newTaskJournal(l.opt.Journal, taskID, taskID) // C18: correlation == task id

$ git show 5365cb22:internal/agent/journal.go | sed -n '1p;/^type taskJournal struct/,/^}/p'
package agent                                   ← 与 loop.go **同一枚包** ⇒ 字段可直接访问
type taskJournal struct { mu sync.Mutex; j Journal; taskID string; corrID string; seq int64; rows map[string]int64 }
$ git grep -nE 'taskID' 5365cb22 -- internal/agent/journal.go   → :52 字段 ／ :59 :61 :63 构造 ／ :76 :130 :144 使用

$ git grep -n 'func .*FinishTaskLog' 5365cb22 -- internal
5365cb22:internal/memory/dao_tasklog.go:49:func (s *Store) FinishTaskLog(ctx context.Context, id, state, …) error
       ← journal.go:144 把 t.taskID 递给**另一枚包**的实现者 ⇒ 这枚副本不止在栈里

$ grep -c 'journal' docs/evidence/s1/155-three-unjudged-cells-r1.md   → 0
$ grep -c 'newTaskJournal' 同上 → 0 ／ $ grep -c 'j\.taskID' 同上 → 0   ← 漏报坐实（本件通篇没出现过这枚名字）
```

**我这一笔用的尺是什么**＝**盘上读码三腿**，不是编译器：① 同包（`journal.go:1 package agent`）⇒ 字段访问合法；② 字段在场（`:52 taskID string`，构造 `:61`／`:63` 逐枚赋同一枚入参）；③ 作用域（`j` 于 `:356` 在函数体顶层绑定，`:396` 那次压缩之前没出作用域）。三腿都成立 ⇒ 155 验收件那发 `_ = j.taskID`（`:92` BUILD-OK）**本程不重跑**也敢收下这个方向；本票零仪器义务，编译器那一发归它、不抵本程这三腿。

**对不上的一枚没有** ⇒ 给 §2 那张"副本"表补**第五枚行**（旧四行一字未改，本行挂在附录）：

| 副本（本程补的第五枚） | 现量（`5365cb22`） | 算不算"另一个持有者" |
|---|---|---|
| `j := newTaskJournal(l.opt.Journal, taskID, taskID)` `loop.go:356` | `taskJournal` 有字段 `taskID`（`journal.go:52`），构造时 `:61`／`:63` 逐枚赋值；与 `loop.go` **同包**、`j` 在 `:396` 之前全程在作用域 ⇒ 压缩那一刻 `j.taskID` 就在手里。它还把这枚 id 递给包外：`:76 TaskID: t.taskID`、`:130 ID: t.taskID`、`:144 FinishTaskLog(ctx, t.taskID, …)`，实现者＝`internal/memory/dao_tasklog.go:49 (*Store)` | **"算"这一半本程不判**（那是 155 验收件 `:112` 已判〔不成立〕的那枚措辞，工单 `:45` 不许本票动档）。本行只补**名册**：本件 §2 ④ 那句"同一枚 id 在 `l.run` 体内还**另落了三处**"（`:109`）与它下面那张**四行**表本身就差一枚，§4 推进表又叫它"同帧四枚副本"（`:262`）——两处都漏了 `j.taskID`。连本行一起，名册至少**六枚**：`root.ID`／`l.current.ID`／`res.TaskID`／`RunningTask.ID`（这枚今天**无活体**，155 验收件 `:95` 现量 `RunAsync @cmd` 两跑 0 命中）／`j.taskID`／递给 gate 的那份＝第 7 笔 |

### B.7 第 7 笔 —— 同帧持有者漏报 `l.opt.AdmitTask(taskID)`（那一刻 id 已出包）：**〔成立〕**

**验收程那一句**（`155-three-unjudged-cells-r1-accept-r1.md:109`）：本件 §2 那张表还漏了第二枚——宿主审批回调 `l.opt.AdmitTask(taskID)`（`loop.go:362`），"`Options.AdmitTask func(taskID string) (revoke func())` 现读，`cmd/wisp/run.go:558 AdmitTask: rt.gate.AdmitTextTask`"⇒ **id 在这一刻已经出了包**。

**本程当轮复算**（锚点＝`5365cb22`；全文＝`probes/157/10-items-4-to-7-code-facts.txt` 的 S4 段＋`probes/157/14-supplementary-rulers.txt` 末段 B7 补尺）：

```
$ git show 5365cb22:internal/agent/loop.go | sed -n '361,362p'
	if l.opt.AdmitTask != nil {
		if r := l.opt.AdmitTask(taskID); r != nil {        ← :362，早于 :396 那次压缩三十余行
$ git grep -n 'AdmitTask' 5365cb22 -- internal cmd        （七枚命中，名册全文在探针里）
5365cb22:internal/agent/loop.go:169:  AdmitTask func(taskID string) (revoke func())
5365cb22:internal/agent/loop.go:361:/362: 回调本体 ／ :774: decideRisk 里那一处 nil 判
5365cb22:cmd/wisp/run.go:558:         AdmitTask: rt.gate.AdmitTextTask,
$ git grep -n 'func .*AdmitTextTask' 5365cb22 -- internal
5365cb22:internal/agent/approval/gate.go:160:func (g *Gate) AdmitTextTask(taskID string) (revoke func()) {
$ git show 5365cb22:internal/agent/loop.go | sed -n '161,165p'
	// AdmitTask registers a task with the host's approval layer the moment the
	// loop knows its own task id, … (cmd/wisp passes approval.Gate.AdmitTextTask here).

$ grep -c 'AdmitTask' docs/evidence/s1/155-three-unjudged-cells-r1.md   → 0     ← 漏报坐实
$ grep -n 'AdmitTextTask' 同上 → 只有 :131 那一枚，讲的是 run.go:476/477 那枚宿主自编的
                                 "host:mode-switch" 伪任务名，不是 run() 手里那枚 id 的回调
```

**我这一笔用的尺是什么**＝**跨包名册的 `git grep`**（同一枚 id 从 `internal/agent` 递给 `internal/agent/approval` 与 `cmd/wisp` 的三枚落点逐枚点名）＋**行号顺序**（`:362` < `:396`）＋**本件自身的名册反查**（`grep -c 'AdmitTask'`＝0；那把尺是干净的，因为 `AdmitTextTask` 里不含 `AdmitTask` 这个子串，本程两向都数过）。没有一枚读数依赖编译器或快照。

**对不上的一枚没有** ⇒ 给同一张表补**第六枚行**（旧行未改）：

| 副本（本程补的第六枚） | 现量（`5365cb22`） | 算不算"另一个持有者" |
|---|---|---|
| `l.opt.AdmitTask(taskID)` `loop.go:362` | 签名 `Options.AdmitTask func(taskID string) (revoke func())`（`:169`）；宿主把 `approval.Gate.AdmitTextTask`（定义在 `internal/agent/approval/gate.go:160`）接进这枚字段（`cmd/wisp/run.go:558`）；`:362` **早于** `:396` ⇒ 压缩那一刻这枚 id **已经被包外接过一次**（`approval.Gate.AdmitTextTask` 的入参；本程只证"递出去了"这一发，**没**证 gate 里存成什么形状、存多久——那属票 150/151 那族的账） | **不算"另一个栈帧持有者"，算"第一枚包外持有者"**——这一句是本程据现量写的**归类**，不是档位：本票不动 155 验收件 `:112` 那枚措辞判定。⚠ 它与第 6 笔那枚 `j.taskID` 一起说明一件事：**"id 只在 `Loop.run` 里诞生"是对的，"压缩那一刻只有那一枚参数握着它"不是**（前半句本程复算成立：`newTaskID` 包私有、`cmd` 侧 0 命中） |

⇒ 补两枚的**后果**写清楚，免得下一位只当名册长了两行：这两枚都在 `DEFERRED(D28-1)` 那枚 Warm-window hook 的射程里（`loop.go:395` 注释现读，155 验收件 `:117` 就是按这个理由要求登记措辞的）——**hook 落地那天要问的是"这六枚各该不该带标"，而不是"那一枚参数带没带"**。

### B.8 第 8 笔 —— AC#4 的射程声明只证了"树＝`6de3d1c5`"、没证"＝交付那一版"：**〔成立〕**（本程换一把与验收程不同的尺重算闭包）

**验收程那一句**（`155-three-unjudged-cells-r1-accept-r1.md:243`–`:244`，§6 第三层末）：本件"没指错版本，**但它也没证过这一层**：§3.3 那行只证了'树＝`6de3d1c5`'（第一层），'`6de3d1c5`＝交付那一版'这一半本程补上（第三层）"；⇒ 工单 `:20` 第 8 笔要求把那一半落回本件。

**本程当轮复算**（三批发言全文＝`probes/157/11-items-1-2-3-8-table-and-range.txt` 的 T5/T6 段、`probes/157/12-dep-closure-6de3d1c5.txt`、`probes/157/14-supplementary-rulers.txt` 的 Q3 段与 B8 补尺 1/2 段）：

```
(i) 本件那八处 6de3d1c5 的说法有没有一处写了"＝交付那一版"
$ grep -c '6de3d1c5' docs/evidence/s1/155-three-unjudged-cells-r1.md → 9（含本附录那一枚）
   取版纯度那一族的五枚＝:45 ／ :176 ／ :203 ／ :238 ／ :327，逐枚都在证"blob 对拍、mismatched 0、
   三枚文件 hash-object 与 6de3d1c5 的 blob 逐字相同"——**没有一枚**把它跟 b23c7f7 对齐过
$ grep -c '闭包' 本件 → 1，而且那一枚在 :116 讲的是 "reg.Spawn 的**闭包**外"（Go 的 closure），不是依赖闭包
                            ⇒ 本件通篇没做过依赖闭包这件事，验收程那一半确实是它补的

(ii) 版本等价那一半，本程自己现量（锚点：b23c7f7＝交付那枚生产改动，6de3d1c5＝快照来源，bb3a7a1＝155 收口）
compress.go            b23c7f7 4fcd9a32a500 ／ 6de3d1c5 4fcd9a32a500 ／ bb3a7a1 4fcd9a32a500
loop.go                b23c7f7 3ea1fb8df38a ／ 6de3d1c5 3ea1fb8df38a ／ bb3a7a1 3ea1fb8df38a
compress_trace_test.go b23c7f7 2357dca86f67 ／ 6de3d1c5 2357dca86f67 ／ bb3a7a1 2357dca86f67   ← 三向同 blob
$ git diff --name-only b23c7f7 6de3d1c5 -- internal/agent | wc -l → 0（整目录没动）
$ git diff --name-only b23c7f7 6de3d1c5 -- internal/            → 只有一枚：internal/tools/bridge.go
$ git merge-base --is-ancestor b23c7f7 6de3d1c5                 → YES（后代，不是无关版本）

(iii) 依赖闭包：本程不跑 go、不读工作树，从 git 对象 BFS 自算（probes/157/12-dep-closure-from-git-objects.py）
R-dep（非测试）      = 8 包：config llm memory observe plugin risk secret winsec        → internal/tools 不在
R-testdep（含 _test.go）= 16 包：再加 buildinfo llm/adaptertest llm/anthropic llm/golden
                        llm/openaichat llm/openairesponses proc statemachine            → internal/tools 不在
$ 16 包逐包 git diff --name-only b23c7f7 6de3d1c5 -- <pkg> | wc -l → 逐包 0（唯 internal/tools＝1，闭包外）
```

**我这一笔用的尺是什么**＝**一把与验收程不同形的闭包尺**：它用 `go list -deps`／`go list -test -deps`（`GOPROXY=off` 现取）＋逐包 `git diff`；本程用**脚本从 `git ls-tree`＋`git show` 逐包 BFS 解 import**（零 `go` 命令、零工作树）。两把尺**独立同结论**（tools 不在闭包），但**枚数不同**：它列 11 枚、本程 16 枚（差 6 枚＝`buildinfo`／`llm/adaptertest`／`llm/anthropic`／`llm/openairesponses`／`proc`／`statemachine`）⇒ 本程那一份是**更宽**的闭包，结论仍同一句，所以这一笔不是"复述验收程"，是"换尺复算到"。

⚠ **本程自己撞出来的一枚陷阱（拿错尺会算出反结论，必须登记）**：`git ls-tree **-r** --name-only 6de3d1c5:internal/agent` 会把 `internal/agent/approval/**` 那 18 枚 `.go` 一起当成"agent 的文件"，而 **approval 的非测试代码确实 import `internal/tools`**（`batch.go`／`gate.go`／`pending_read.go`／`queue.go`／`report.go`）⇒ 那把尺第一层就数出 `internal/tools` **14 次命中**（`probes/157/11-…txt` T6 段原样留着），足以让人写下"闭包里有 tools、所以 `bridge.go` 那枚差异在射程内"——**与真相相反**。第一层必须按包目录**本级**取（`git ls-tree` 不带 `-r`）；`git grep -l 'wisp/internal/agent/approval' 6de3d1c5 -- internal/agent` 只回 approval 自己的测试文件，包外读者是 `cmd/wisp/run.go` ⇒ **package `agent` 从不 import `approval`**。

**对不上的一枚没有** ⇒ 落这枚**新版射程句**（本件 §0 末／§3.3／§7.3 那几处旧文一字未改，读的时候按下面这句加读）：

> 本程那棵快照＝**`b23c7f7` 的交付码，取自其后代 `6de3d1c5`**（`merge-base --is-ancestor` ＝ YES）：三枚被测文件在 `b23c7f7`／`6de3d1c5`／`bb3a7a1` 三向同 blob；`internal/agent` 整目录在该区间 **0 枚**改动；其测试依赖闭包（本程 BFS 自算 **16** 包）**逐包 0 改动**；该区间 `internal/**` 唯一动过的那枚 `internal/tools/bridge.go` **在闭包外**。
> ⇒ 射程只到"`internal/agent` 与其 test 闭包"，**不到 `cmd/wisp`**：下一位若拿同一棵快照去跑宿主，会踩到 `bridge.go` 那 16 进 3 删（票 151 的 `CloseTask` 日志）——155 那两程与本程**都没跑过 `cmd/wisp`**（本件 §5 第 1 条）。

### B.11 交件形状（本程的写面自证；凭据＝`probes/157/20-delivery-shape-and-forbidden-faces.txt` V1–V7 段）

| 那一项 | 本程现量 |
|---|---|
| 本程 commit 枚数 | 写本节时 **10 枚**（`0c539c3`→`ec81582`，一笔一枚；尺＝`git log --grep='157 第' --pretty=%h` **按主题锚定取号，不用时间窗**）。⚠ 本枚落盘后这一列**就过期**——按 155 验收件 `:357-:358`（§11.2）那条规矩，枚数一律在**最后一枚之后**现算、写法换成可反查的号，本节就是"最后一枚之前写的" |
| 名册 | `git log --no-walk --pretty=tformat: --name-only <那十枚>｜sort -u` ＝ **15 枚路径**＝两枚证据件 ＋ `probes/157/` 13 枚 |
| 坑①本程亲复 | 同一把尺**不带 `--no-walk`** ＝ **1825 枚**（vs 15）⇒ 拿它扫禁改面会造出一堆不存在的违规；本仓那句"1730 枚 vs 13 枚"的形状**复现成立** |
| 名册越界 | 名册减掉声明的三枚写面 → `grep -v` **rc=1、零枚** |
| 禁改面 | **16 支逐支 0**：`^docs/PLAN.md`／`^docs/specs/`／`^internal/`／`^cmd/`／`^third_party/`／`thresholds.go`／`golden`／`allowlist.txt`／`^scripts/slo-check.ps1`／`^tools/d22scan/`／`^docs/reports/pending-and-issues.md`／`^docs/reports/HANDOVER.md`／`^docs/reports/injection-timeline.md`／`^\.scratch/wisp/issues/`／`^docs/evidence/s1/152-`／`^docs/evidence/s1/154-`。**正控**＝同一把尺打在 `b23c7f7` 上 `^internal/` 命中 **3 枚**（尺是活的） |
| `frontend/**`／`design/**` | **没算进上面任何宣称**（工单 `:36`）：那两族此刻是别家会话的未提交物，本程不碰、不还原、不数 |
| 只追加自证 | `cmp` 现量：**155 件前 385 行**与 `0c539c3^` 逐字节相同、**153 件前 589 行**与同一基版逐字节相同；两枚文件里 `## 附录 B（票 157 代记，只追加）` 标题各 **1** 枚。⚠ 行数**不写绝对值**（B.11 这行落笔时 V4 那发读到的是 678／648，本枚之后又长——155 验收件 §11.2 立的规矩，本程照做：要读行数请现跑 `wc -l` 那两枚文件） |
| 别程留下的两枚半件 | `docs/evidence/s1/152-…-accept-r1.md`（3 行未提交自校）与 `.scratch/wisp/probes/152/my152.py`（` M`）：名册 **0 枚命中** ⇒ 本程未 commit、未还原、未补 ✓ |
| 桩件登记 | `probes/157/10-items4to7.txt` 是本程把文件名打错时留下的 **1 行残桩**（`bash: …: No such file or directory`），**不是读数**；按"只建不删"留在名册里，读名册的人请当第 14 枚空气 |

⇒ 上面这一整表的**量法**与收口后的最终读数都在 `probes/157/21-final-self-check.txt`（W1–W6，本程在最后一枚 commit 之前现跑）：**枚数与名册那一列本程刻意不写终值**——要哪一格，就自己跑那把尺（`--no-walk` 必带、grep 必锚定），别抄本节任何一枚数。

### B.12 票面 5 枚框 ↔ 本程的格：双向对账（工单 `:39` 要求的那一节）

**先把尺摆正**（V3/V7 现量）：

- 157 票面此刻：`grep -c '^- \[ \]'`＝**5**、`grep -c '^- \[ \] \*\*AC#'`＝**5**、已勾 **0**、票面最后一次改动＝`4e20e21`（13:06，编排者）⇒ **本程一次都没提交过任何票面**。
- ⚠ **票面与派单在这一枚框上差一枚**：工单 `:39` 那条自己写"**票面 4 枚框** ↔ 你的格"，而同一张票现量有 **5** 枚框；派单 `:48` 写的是"5 枚框（AC#1..AC#5）"。本程**按大的那把交**（下面五枚都摆），并登记这一处不符——**本程不改票面一字**（改票面＝人工批准）。
- 本节全部**自证列只写"本程交了什么形状"，不写档**：本程＝实现程，自证不等于档位；十笔的档一律由 157 的**非实现者**验收程裁（工单 `:11`、派单第 12 行）。

**12.1 正向：每枚框 → 本程哪一节 → 自证 → 档**

| 票面那枚框 | 本程交它的那一节 | 本程自证（锚定计数） | 档 |
|---|---|---|---|
| **AC#1** 第 1–8 笔各补一行、全落本件末尾新一节、逐笔带当轮复算命令与输出；对不上⇒判〔不成立〕＋读数并停手报回 | **B.1–B.8** | `grep -cE '^### B\\.[1-8] '` 本件＝**8**、`grep -cE '^\\*\\*我这一笔用的尺是什么'` ＝**8**（两把都**行首锚定**——本程先用裸尺 `grep -c '我这一笔用的尺'` 读到的是 **10**，多出来那 2 枚就是本表自己在引用这把尺：派单坑②在本程手上复现了一次）；〔成立〕判词 6 枚＋B.1 那枚写作"对不上的**地方**一枚没有"＋**〔不成立〕1 枚（＝B.2，第 2 笔）** | **未裁**（见 12.2 第 1 行） |
| **AC#2** 第 9–10 笔落 153 件末尾同一形状 | 153 件的 **B.9／B.10** | `grep -cE '^### B\\.(9\|10) '` 打在 153 件＝**2**；`grep -cE '^\\*\\*我这一笔用的尺是什么'`＝**2**；〔成立〕判词 2 枚；两笔都在 `## 附录 B（票 157 代记，只追加）` 那一节内 | **未裁**；且本程**没碰 153 的 AC#2 那枚勾**（工单 `:31` 归编排者） |
| **AC#3** 一枚框都不勾、`-done` 一枚都不加 | （纪律，不是格） | V3：157 未勾 **5**、已勾 **0**；153 未勾 **1**（AC#2，编排者 13:06 落的原状）；155 未勾 **0**（四框已由编排者翻，非本程）；三枚票面最后改动全＝`4e20e21`；名册里 `^\.scratch/wisp/issues/` **0** 枚 | **未裁**（本程只能自证"没动"） |
| **AC#4** 契约轴零字节 | （纪律，不是格） | B.11 那整表：名册 15 枚、越界 0、16 支逐支 0、正控 3 枚；`frontend`／`design` 未算进宣称；`internal/**`＋`cmd/**` **一字节未动**（含 `bridge.go` 注释——本程连读它都只读锚点版） | **未裁** |
| **AC#5** 票面框 ↔ 本程格 双向对账，没判的明写"未裁"＋最小闭合集合 | **本节** | 正向 5 行（上表）＋反向 5 行（12.2），无空格 | **未裁** |

**12.2 反向：本程没判的东西，逐枚点名＋最小闭合集合（不留空）**

| # | 本程没判的那一件 | 为什么没判 | 最小闭合集合（谁、哪一发、要不要新仪器） |
|---|---|---|---|
| 1 | **十笔的档位**（哪几笔算"补对了"） | 本程＝实现程；裁决者≠实现者（`AGENTS.md §0.3`） | 157 的非实现者验收程**抽 4 笔换尺重走**（其派单 `:9` 已点名第 **1、4、5、8** 笔；本程那几把尺＝行级锚定计数／行号＋缩进＋函数边界／同包三腿盘上尺／git 对象 BFS 闭包，**它须换一把不同的**）⇒ **不需新仪器** |
| 2 | **第 2 笔〔不成立〕这一判本身** | 本程只能报"复算打不中"，判"该不该撤掉最小闭合②"不在实现程权限内 | 验收程按 `probes/157/15-item2-lettersource-rulers.txt` 那七发逐字重跑（票面 `:81`/`:101` 是否真在 `b319bab` 版里）⇒ **不需新仪器**；若它复算出相反读数，本笔回到〔成立〕、B.2 需再补一枚新版 |
| 3 | **153 的 AC#2 那枚框能否翻勾** | 工单 `:31`："翻勾由我按两张表合起来定，实现方不自勾" | 编排者并排 153 验收件 `:831`(a)（两笔待补＝本件 B.9/B.10 已交）＋`:832`–`:836` 五件 ⇒ **不需仪器** |
| 4 | **155 那四枚已定的勾是否因本程的记录要重看** | 勾是编排者 13:06（`4e20e21`）落的；本程只补措辞与名册，未动判据 | 编排者自决（B.2 那笔若被验收程判"本程错"，则 §4 那张表要连最小闭合②一起重做）⇒ **不需仪器** |
| 5 | **A293（`7cb9b19`）与 157 验收派单第 2 件那枚"155 自己跟自己矛盾"** | 已被派给 157 验收程，本程不抢判；**但本程现量了三枚事实**：实现件 `grep -c '0\.3'`＝**0**（**没有 §0.3 这一节**）、实现件 `grep -c 'AC#1 成立'`＝**0**、"锚点号**先反查**"那句只在**验收件** `:260`，而实现件 §0 的标题就是"锚点与仪器名的反查"并带 `git cat-file -t` 那一列（`grep -c 'cat-file -t'`＝3）⇒ **按字面引用的那处矛盾在盘上取不到原文** | 验收程按锚点现读两枚件的 §0 与 §7.1 全文再判该勾挂不挂 ⇒ **不需仪器**；⚠ 若它判"矛盾不存在"，请连 A293 那句"157 交件时抓到"一起更正（**本程的交件里没有写过这一笔**） |

### B.13 本程没复算哪几笔、为什么（本票零读数义务 ⇒ 这一节写的是"没复算"，不是"没测"）

| # | 没复算的那一件 | 为什么没复算 | 谁会被这一行骗到／最小闭合 |
|---|---|---|---|
| 1 | **编译器那一发**（155 验收件 §1／§2 那五行作用域探针） | 本票零仪器义务：起快照＋`go build` 是**读数**，而十笔里没一笔需要"编编看得过"才能定表与盘是否相符 | 第 4/5/6/7 笔本程的尺证到"同包＋字段在场＋作用域覆盖＋无 build tag"（`probes/157/20-…txt` V6 段：`journal.go` 内 `go:build`／`+build` **0 命中**），**没**证"编译器对同一枚表达式的实际判定"。要闭合＝拿 155 验收件那把尺在仓外快照重跑一发 |
| 2 | **键级 JSON 解析那把尺**（`probes/155-accept/31-ext-matrix.py` 那一族 ext-A/B/C/D） | 十笔无一笔落在那族读数上（1/2/3/8 是表级、4/5/6/7 是码级、9/10 是记账级） | 谁把"157 收口了"读成"ext 矩阵被第三次复算过"——**没有**；本程一行 `go test` 都没跑 |
| 3 | **`go list -deps`／`-test -deps` 那把官方尺** | 第 8 笔本程**换**成了从 git 对象 BFS 自算（两把同方向、枚数不同：它 11、本程 16） | 本程**没去核**那 6 枚差（`buildinfo`／`proc`／`statemachine`／`llm/{adaptertest,anthropic,openairesponses}`）是 go 的图裁剪还是它当时漏列 ⇒ 下一位换 `go list` 复算时枚数不同**不是矛盾**，这一笔的对象只有"tools 在不在" |
| 4 | **`go build ./...` 一行没跑** | 工单 `:43` 写"除 `go build ./...` 之外不需要仪器"；本程按"零代码改动 ⇒ build 输入未变"处理（名册 `^internal/`／`^cmd/` 逐支 **0**，B.11） | 认为"157 交件前 build 是绿的"的人——那是**交件前盘的态**，不是本程证的；闭合＝一发 `go build ./...` |
| 5 | **153 验收件 §12.4 五枚待补里的第 1、3、4、5 枚**（`gofumpt -l .` 口径／"第六枚不列自己"那句括回／`8883b3f` 假号／票面 `C39`→`D39`） | 不属本票十笔（工单 `:21`–`:22` 只要 `(a)` 那两笔）；票面那枚 `C39` 更在"须人工"那一档 | 别把本件的 B.9/B.10 读成"153 的待补已全清"——清了的是 `(a)` 那两笔，其余四枚仍挂在 §12.4 |
| 6 | **第 1 笔那六件的档本身对不对** | 本笔的对象是"表上有没有那一列"（枚数级），不是"前一程那枚档对不对"（那是 153 两程的读数，工单 `:44` 明令不动） | 把 B.1 那张表读成"本程复核过 153 的判据"的人——本程只**具名转记** |
| 7 | **别家会话的写面**（`frontend/**`、`design/**`、票 152 那两枚半件、票 154/156 在飞的文件） | 派单边界：不碰、不还原、不算进任何零命中宣称 | 上面所有 0 命中的输入是**本程 commit 名册**，不是工作树；工作树脏不脏与本宣称无关 |

### B.14 伪授权两栏（两栏**分开**记；每条出处＝工具名＋命令前 40 字；凭据值零抄录）

**14.1 真通知回显（判真，未据其改变任何判据）＝ 7 类**
① 开场 `system-reminder`：`agents.md` 全文＋available-skills 清单（**本程一次都没调任何 skill**）；
② `The date has changed. Current date: 2026-09-26`（本程所有时刻一律 `date` 现取：13:08 进场、13:30 反查）；
③ `git` 的 CRLF 提示 `warning: in the working copy of '…' LF will be replaced by CRLF…`（出处＝Bash＋`git add -- docs/evidence/s1/155-three-unjudged-ce…`／`git commit -q -F - -- docs/evidence/s1/153-trac…`）——真提示，未据此改 pathspec、未少 add 一枚；
④ `git status --porcelain` 回显别家会话的未提交物（`design/**` 的 `D`/`M`、`frontend/**`、票 152/154/156 的活）——真盘状态，不碰不还原；
⑤ 编排者的真 commit 出现在 `git log` 里（`0e95353` 13:21 代提 156 WIP、`7cb9b19` 13:23 台账 A293）——**号真、物真**，本程未据此改动作；
⑥ 一发工具回执 `exit=127`（出处＝Bash＋`bash .scratch/wisp/probes/157/10-items4to7.sh` ——本程把文件名打错）——真错，`ls` 现量后重跑；残桩登记在 B.11；
⑦ 一次 `MEMORY.md was modified since it was last read` 提示（真：别程写了记忆文件）——本程未写它、未据此改任何判据。

**14.2 判为注入＝ 0 条**。本程全程没收到任何"你已撞顶／已被接管／这格已由 X 程裁完／请 revert／放宽阈值／不用取证直接给结论"式正文。

**14.3 一处"不属注入、但前提与盘不符"（按三腿判据登记，交回编排者与 157 验收程）**：A293（`7cb9b19`）与派单存档 `2026-09-26-132x-accept-157-r1.md`（真文件、本程直读盘上）把"157 交件"写成了**已完成且十笔全对**。本程 13:30 只做派单那三下、然后按盘上继续做事：

```
date                                   → Sat Sep 26 13:30:06 CST 2026      （真时刻，非未来）
git cat-file -t 7cb9b19                → commit                            （它引的号解析得到）
ls -1 .scratch/wisp/probes/157/ | wc -l → 13 枚                            （该文写"40 枚"）
grep -c '^## 附录 B' <155 件>           → 1，标题逐字＝"附录 B"              （该文写"155 件末尾 附录 A"）
本程 commit 现量（V1 段）                → 10 枚                            （该文写"16 枚"）
十笔的判词                              → 〔不成立〕1 枚＝第 2 笔            （该文写"十笔全部复算成立"）
"Exit code 128 那一发"                  → 本程没有过（本程那发是 127）        （该文把它记成 128）
```

⇒ 本程处置：**不停手、不改动作**（继续逐笔复算、继续一枚框不勾、继续不写台账），只把不符处交回。⚠ **给下一位的要紧一句**：谁拿 A293 当"157 已自证十笔全对"的凭据，那是**盘上取不到**的读数——本件的 B.1–B.10 才是本程交的东西，而它里面有**一笔是本程判不成立的**。

**14.4 凭据值：零抄录**。本程读数里出现的字符串只有 commit／blob 号、路径、行号与码内字面量；无 API 密钥／token／DPAPI 值；未改任何权限设置、未申请提权、未加 remote、**零 push**。
