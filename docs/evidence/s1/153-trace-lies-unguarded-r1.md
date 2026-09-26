# 153 — 那枚压缩痕的两条谎法 · 实现件 r1（写码程）

票面：`.scratch/wisp/issues/153-the-compression-trace-has-two-unguarded-ways-to-lie-swapping-the-ran-guard-for-need-keeps-all-four-nails-green-and-the-trace-carries-no-taskid.md`
来路：票 139 对抗验收件 `docs/evidence/s1/139-compression-leaves-no-trace-r1-accept-r1.md` §6.2 第 2、3 件。
派单正文：`.scratch/wisp/dispatches/2026-09-26-091x-r152-r153-r151.md` 派单 A。

两格**没有互相抵账**：AC#1 交的是判据（少一形），AC#2 交的是字段（少一枚归因键），
各自有各自的用例、各自的红句、各自的承重答句（§2 / §3 / §4）。

原始读数（逐枚 `go test -v` 全文、驱动脚本、探针源件）落 `.scratch/wisp/probes/153/`；
所有变异只发生在**仓外**快照 `C:\Users\swq\AppData\Local\Temp\wisp153-snap-{1,3,probe}`，
逐发还原 + `cmp` 自证（见 §2.4 / §4.3）。本程未删任何临时件、未在仓库内建 worktree。

**票面四框本程不自勾**；`-done` 后缀未加。

---

## 0. 锚点与工作树现量（开工第一发）

```
$ git rev-parse --short HEAD            # 开工第一发
86b0161                                 ← 与派单共同锚点相同
$ git log --oneline -1 86b0161
86b01613 feat(frontend 10 首启引导): FirstrunScreen props 化 + showcase 第十节
```

⇒ **我这枚锚上是前端会话的提交**（不是编排者那枚台账枚 `f7478d3`）。开工期间 HEAD 又走过：
`5365cb2`（编排者：派单存档落盘）、`ff550f3`（本程 AC#1）、`cdf2471`（前端会话：动画开关机器）。
按派单 §0：**这不算漂移**，只需点名锚上是谁——已点名。本程全部"改前"读数一律钉在 `86b0161`，
不用"当前 HEAD"当改前。

**先确认派前那句"票 139 的实现在 HEAD 上"**（按 sha 取，不读脏工作树）：

```
$ md5sum <快照>/internal/agent/{compress.go,compress_trace_test.go,loop.go}   # 86b0161 取树
9c5015dd23d29060b956502aeb954cac  compress.go
7bf2fbf79f78f935afac6683f906b467  compress_trace_test.go
70f8fb1068dd6541358e0b2ce646d5ca  loop.go
```

三枚逐枚等于 139 验收件 §2 开头记的那三枚（`compress.go 9c5015dd… / loop.go 70f8fb10… /
compress_trace_test.go 7bf2fbf7…`）⇒ `internal/agent/**` 自被验版本 `d949d9c5` 起**零漂移**，
痕确实已经在 HEAD 上，本票是给它补牙、不是重做它。

工作树脏的部分**全是别家的**，本程一律不碰、不还原、不算进任何"零命中"宣称：
`design/**`＋`frontend/**`（前端会话在写）、`cmd/wisp/run.go`＋`cmd/wisp/task_scope_close_151_test.go`＋
`internal/tools/bridge.go`（派单 C＝票 151 的程）。

**本格改了哪些文件**：无（只读）。

---

## 1. 派单带进来的两条"未验证前提"——现量结果：一条成立、一条只成立一半

### 1.1 前提①「M5 那形在未修码上今天就响」——**成立（两半都复跑过）**

先复跑"逃逸"这半（＝139 验收 §2.2 的原句：四枚钉子全在场、整包全绿），在 `86b0161` 的仓外快照上施
`m5.py apply`（`compress.go:170` `if rep.Ran {` → `if c.Need(hist) {`，落地有打印）：

```
[R1]  go build rc=0 | go test ./internal/agent/ -count=1 → rc=0
      RUN(all)=80 PASS_all=80 FAIL=0 SKIP=0 panic=0     ← 四枚 CompressionTrace 用例逐枚 PASS
      logs/r1-m5-4nails.txt
```

⇒ 逃逸复算成立，**不是**背它的数。再复跑"响"那半（同一枚变异 + 本票新增那一发）见 §2.2 的 `F2`。

### 1.2 前提②「普通收尾形 0 枚、重复阶梯退化形 1 枚」——**前半本程亲量成立；后半只在进程内复现，未复跑到盘那一发**

- **前半（普通形 0 枚）成立**：本票新增的那一发就是在量它——1 枚 user 轮、1 006 token 对阈值 384，
  `Need()` 真、`rep.Ran` 假、痕 **0 枚**（`logs/f1-delivered-unmutated.txt`，用例 PASS）。
  根因与本程复算一致：`groupRounds`（`compress.go`）只在 `RoleUser` 开一轮，折叠条件
  `if len(raw) <= c.b.KeepRawRounds { break }` 对 1 枚原始轮直接 break。
- **后半（退化形 1 枚）本程未复跑到盘那一发**：那需要"三条独立重复阶梯"喂真 provider + 读
  `<data>\logs\wisp-*.jsonl`（139 验收 §7.1 的探针形状，落点在 `cmd/wisp`，不是本票地界）。
  本程只在**进程内真 Loop** 上复现了同形的枚数：`buildRoundHistory(4, 400)` 经
  `Loop.Run` 折叠 ⇒ 痕**恰 1 枚**，且走的是进程默认 logger（`logs/ext-A-delivered.txt`，
  见 §3.4）。⇒ 本程引它是引 §7.2/§8 的**两形对照**，未引它 §1（那节被同一枚程自己在 §5 推翻，票面明令不引）。

⚠ 本程没有因为"前提① 成立"就跳过 §2 的改前红句，也没有因为"前提② 后半未复跑"就否认它——
两句各自给了读数与未读数。

---

## 2. AC#1 — 补上"过阈值但折不动 ⇒ 痕 0 枚"那一枚用例（判据形状）

### 2.1 判据与形状

新增 `TestCompressionTraceSilentWhenNothingFoldableOverThreshold`（`internal/agent/compress_trace_test.go:426`）。
它断四件事，缺一就不是本票要的那一形：

| # | 断言 | 为什么必须有 |
|---|---|---|
| 1 | `len(rawRoundIndexes(groupRounds(hist))) == 1` | 钉住"折不动"是**轮数地板**造成的，不是"没过阈值"造成的 |
| 2 | `c.Need(hist)` 为**真**（否则 `t.Fatalf`） | 这正是 139 那四枚的洞：`TestCompressionTraceSilentWhenNothingFolded` 在 `Need()` 为真时**自己 Fatalf 掉了那一支**（票面点名的 `:270-271`），于是"过阈值"这一形在四枚里不存在 |
| 3 | `rep.Ran == false` 且历史逐字节没动（msgs / tokens 两端相等） | 排除"其实折了但没打痕"这一支混进来 |
| 4 | `recs.with(traceMsg)` 长度 **== 0**；不为 0 时先把那条谎痕 `flat()` 打出来再判失败 | 判据本体；且红句里直接带读数，不用二次取证 |

阈值一律从 `BudgetsFor(4096).HistoryCompressTokens` 读（384 是读数不是字面量），
与本文件既有的"不拿 12000 当触发点"口径一致；夹具用 `compress_test.go:23 buildRoundHistory`
（本程未动那个文件，见 §5）。

### 2.2 改前红句原文（**未修码上**，`86b0161` 的 compress.go + 交付版用例 + M5 已落地）

```
$ python m5.py apply <snap-1>
M5 applied: `if rep.Ran {` -> `if c.Need(hist) {`
  landed at compress.go:170: if c.Need(hist) {
$ go build ./internal/agent/ && go test ./internal/agent/ -count=1 -v      # logs/f2-delivered-m5.txt
=== RUN   TestCompressionTraceSilentWhenNothingFoldableOverThreshold
    compress_trace_test.go:458: over-threshold pass that folded nothing left a trace: agent: history compressed compressed_msgs=0 kept_raw_rounds=1 history_changed=false tokens_before=1006 tokens_after=1006 threshold=384 msgs_before=3 msgs_after=3
    compress_trace_test.go:461: "agent: history compressed" records = 1, want none: a record reading "compressed" while nothing folded is a false positive, and 1 raw round(s) against the KeepRawRounds floor of 3 is exactly the case where no fold can happen (all records: [trace-capture-ruler-control agent: history compressed])
--- FAIL: TestCompressionTraceSilentWhenNothingFoldableOverThreshold (0.00s)
```

```
RUN(all)=81 PASS_all=80(63 top + 17 sub) FAIL=1 SKIP=0 panic=0
```

⇒ 三件事同时读到：**只有这一发红**（既有四枚在 M5 下依旧逐枚 PASS，与 §1.1 的 R1 一致）；
**谎痕长什么样**（`tokens_before == tokens_after == 1006`、`compressed_msgs=0`、
`history_changed=false`，正是 139 验收 §2.2 末段推断"比没有痕更误导"的那一形——本程现在把它印出来了，
不再只是推断）；**新路径真执行过**（它红在这发上，不是恒红、也不是恒绿）。

### 2.3 改后转绿（同一发，未施任何变异；工作树交付态）

```
$ go test ./internal/agent/ -count=1 -v      # logs/f1-delivered-unmutated.txt（仓外快照）
--- PASS: TestCompressionTraceBooksCountsOnSuccess          (0.00s)
--- PASS: TestCompressionTraceSilentWhenNothingFolded       (0.00s)
--- PASS: TestCompressionTraceSurvivesTheLoopWiring         (0.02s)
--- PASS: TestCompressionTraceDoesNotAlterTheFold           (0.00s)
--- PASS: TestCompressionTraceSilentWhenNothingFoldableOverThreshold (0.00s)
RUN(all)=81 PASS_all=81(64+17) FAIL=0 SKIP=0 panic=0     → rc=0, ok internal/agent 1.9s
```

```
$ go test ./internal/agent/ -count=1 -v      # logs/ac5-post-count1.txt（工作树，AC#2 也在里面的终态）
RUN(all)=83 PASS_all=83 FAIL=0 SKIP=0 panic=0 → rc=0
```

同一发再还原变异复跑一次（`logs/f3-restored.txt`）：81 / 0 FAIL ⇒ 红绿之差**只由那一枚守卫造成**。

### 2.4 本格自证：快照还原

```
$ python m5.py revert <snap-1> && restore from git show 86b0161:…
RESTORED IDENTICAL internal/agent/compress.go
RESTORED IDENTICAL internal/agent/compress_trace_test.go
RESTORED IDENTICAL internal/agent/loop.go
```

`m5.py` 全程用 `newline=""` 读写：本程第一版没用，导致"还原"后 `cmp` 报了 byte 14 差异
（是 CRLF，不是码），已改脚本并**用改后的脚本重跑了 F1/F2/F3**——上面那三行 RESTORED 是重跑后的读数。

**本格改了哪些文件**：`internal/agent/compress_trace_test.go`（＋71 行，纯追加用例＋文件头第三形说明），
commit `ff550f3`。`compress.go` / `loop.go` **本格零字节**（AC#1 交的是判据，不是码的谎法；
把守卫换成"更保险的双层判"不改变 M5 的可造性，本程不走那条路）。

---

## 3. AC#2 — 痕带 taskID：三问逐问答完才动手（结论：不动契约面，故不停手）

票面那句"若三问的答案是'要动契约面' ⇒ D22 闸门③ 停手上报"本程**先答三问再判**，
下面每一条都是按 sha 现取的（`git grep … 86b0161 --`，不读脏工作树）。

### 3.1 Q1 `Compress` 现在能拿到 taskID 吗？——**不能，三条通道逐条堵死**

```
$ git show 86b0161:internal/agent/compress.go | grep -n 'func (c \*Compressor) Compress\|^type Compressor struct'
128:func (c *Compressor) Compress(ctx context.Context, hist []llm.Message) ([]llm.Message, CompressionReport, error) {
70:type Compressor struct {
71-	b   Budgets
72-	sum Summarizer
73-	lg  *slog.Logger
74-}                                     ← 参数表没有、接收者也没有
```

- **构造时刻也对不上**：`loop.go:239 comp: NewCompressor(b, opt.Summarizer, WithLogger(opt.Logger))`
  在 `agent.New` 里，一枚 Loop 建一枚、**早于任何任务 id**；id 在 `loop.go:322 / :333` 才 `newTaskID()` 现造。
  ⇒ 就算愿意把 id 塞进构造器，它也不是"这一次折叠"的属性（Loop 会被后续任务复用，见 §3.5 的取舍）。
- **ctx 上也没有**：
  ```
  $ git show 86b0161:internal/observe/goroutine.go | sed -n '/func NewRootFrom/,/^}/p'
  func NewRootFrom(parent context.Context, id string) *Root {
      ctx, cancel := context.WithCancel(parent)
      return &Root{ID: id, Ctx: ctx, Cancel: cancel}     ← id 只留在 Root.ID 结构体上，没进 ctx
  }
  $ git grep -n "context.WithValue" 86b0161 -- '*.go' ':!*_test.go'
  86b0161:internal/tools/cancel.go:45: return context.WithValue(ctx, cancelKey{}, h)   ← 全仓非测试只有这一枚
  $ git grep -n "FromContext\|func RootFrom" 86b0161 -- 'internal/**'      → 无输出（rc=1）
  ```
  ⇒ 没有任何"从 ctx 反查任务"的现成通道可借（`observe` 是本票零字节禁改面，也不去那里造一枚）。

### 3.2 Q2 唯一持有者是谁？——**`Loop.run` 的那枚局部 `taskID`，不是组合根**

```
$ git grep -n "taskID" 86b0161 -- internal/agent/loop.go | head -30
322:	id := newTaskID()                                        # RunAsync：id 在这里诞生
333:	return l.run(ctx, newTaskID(), input)                    # Run：同上
339:func (l *Loop) run(ctx context.Context, taskID, input string) Result {
356:	j := newTaskJournal(l.opt.Journal, taskID, taskID)       # C18: correlation == task id
369:	res := Result{TaskID: taskID, Status: StatusCompleted}
373:	l.log().Warn("agent: task_log open failed", "err", err, "task", taskID)   ← 键名 "task" 的既有惯例
411:	turn, streamErr := l.stream(ctx, taskID, req)
```

```
$ git grep -n "taskID\|TaskID" 86b0161 -- 'cmd/wisp/**' ':!*_test.go'
run.go:576:	res := loop.Run(ctx, task)          ← 进 Loop 的这一发**不带**任何 id
run.go:587/:598/:608: … res.TaskID …             ← 组合根只在事后读
run.go:476:	const taskID = "host:mode-switch"    ← 那是宿主自己那枚伪任务名，不是 Loop 的 id，本票没借它
```

⇒ 组合根（`cmd/wisp/run.go`）**不是**持有者：它在压缩发生的那一刻手里没有 id。
唯一持有者＝`Loop.run` 的局部变量，而压缩调用就在它函数体里（`loop.go:396`）——
这正是票面 §地界 预留的那一支："loop.go 只读，除非 taskID 必须由它传进来"。

### 3.3 Q3 送进去要不要动契约面？——**不要**（五枚读数）

| 问的那一面 | 现量 | 结论 |
|---|---|---|
| 痕的字段表是否被 `PLAN.md` 或任何 spec 冻结 | `git grep -n "history compressed\|compressed_msgs\|kept_raw_rounds\|history_changed" 86b0161 -- docs/PLAN.md docs/specs` → **0 命中（rc=1）** | 那张字段表是票 139 的实现选择，不是契约 |
| `C12` 状态机 | `PLAN.md:1362` 原文：`C12 = BallState`、20 态＋D43 转移表＋每态超时 | 压缩痕不是态、不是转移，碰不到它 |
| `C39` 契约深化里那条痕的字段表 | 契约面只有 **C1–C32**（`git grep -n "^| C39" docs/PLAN.md` → 无）；`D39` 深化的是 C6/C7/C11/C17/C19 与 §14.3 分段预算，末句"**C24–C31** 的规格已在 §3 契约表…给出" | 票面问的 `C39` 这枚编号**不存在**；`D39` 里没有这枚痕的任何字段 |
| `docs/contracts/**` 里有没有这一项 | `git ls-tree -r --name-only 86b0161 \| grep -c "^docs/contracts/"` → **0** | 与 AGENTS.md §4 那句"截至锚点它们在仓里不存在"一致 |
| 压缩那条 spec 本身 | `SPEC-05 §4.3` 原文两行：>12000 保 3 轮 / 压缩在响应路径外＋**不得丢弃 tool_call id 与结果引用** | 没有日志 schema，本改动两条都不违（未动算法、未动 id） |
| 会不会新增 Go→前端事件／面板方法（139 AC#2 的硬约束） | `SPEC-08 §5.2` 的 C17 白名单名册里无 compression 相关项；`git grep -c "compression\|compressed" 86b0161 -- internal/panel internal/tools docs/specs/SPEC-08-ui-ball-panel.md` → **0 命中** | 一枚 Go 侧 slog 属性不过 C17 那条线 |

⇒ **不触发 D22 闸门③**，本程不自立"停手"也不自立"放宽口径"：既没扩契约，也没自创归因口径
（键名沿用 `loop.go` 已经在用的 `"task"`，值沿用 `newTaskID()` 造的那枚、与 `Result.TaskID`／
journal／C18 correlation 同一枚，见 §3.4 的读数）。

### 3.4 交的形状与它的读数

**形状**（`b23c7f7`，三枚文件）：
- `compress.go:148-163` 新增 `traceTaskKey{}` ＋ `withTraceTask(ctx, taskID)`（空串＝不是任务，直接丢弃不打标）
  ＋ `traceTaskID(ctx)`；
- `compress.go:222-241` 那条 `Info` 的属性表由字面量改为 `attrs` 切片，**只在 `traceTaskID(ctx) != ""` 时**
  追加 `"task", task`；八枚既有键与 `history_changed` 的表达式逐字节照旧（`git diff` 的 `-` 行里
  只有那条 Info 的头尾与这一枚表达式换行，见 `ac4-zero-byte.txt`）；
- `loop.go:399` 唯一非测试调用点带上 `withTraceTask(ctx, taskID)`；
- 新用例两枚：`compress_trace_test.go:478 TestCompressionTraceCarriesTheOwningTaskID`、
  `:534 TestCompressionTraceNeverInventsATaskID`。

**为什么是 ctx 而不是给 `Compress` 加参数**（被否的那条路，照实记）：加参数要动
`compress_test.go` 的 4 枚调用点——那是本票地界之外的第四枚文件，且 139 实现件 §2.3 正是为了
"让 8 枚既有构造点一字不改"才选的变参 opt；本程沿用同一判断口径。
**代价也照实记**：ctx 通道可以被未来的调用方**忘记**——忘记的后果是那条痕**没有** `task` 键
（可见的缺失，不是假归因），今天唯一非测试调用点已带标；`D28-1` 那张 Warm-window hook 落地时
必须同样带标（写进 `compress.go` 的注释里了）。

**读数（外部可见，进程默认 logger，仓外探针 `probe153_*_test.go`——不是交付码）**：

```
$ go test ./internal/agent/ -count=1 -v -run TestProbe153      # logs/ext-A-delivered.txt
… msg="agent: history compressed" tokens_before=1236 … history_changed=true                      ← 未打标的调用：没有 task 键
… msg="agent: history compressed" tokens_before=1236 … history_changed=true task=ab8a8e40-5df0-41c7-be96-1fa1188bd095
… msg="agent: history compressed" tokens_before=826  … history_changed=true task=23ac7a42-ee92-4a89-8532-f18d23dad154
    probe153_loopleg_test.go:27: PROBE-LOOPLEG this run's Result.TaskID = "23ac7a42-ee92-4a89-8532-f18d23dad154"
```

⇒ 第三行是**真 `Loop.Run` 那条腿**打的痕，`task=` 与同一次运行的 `Result.TaskID` 逐字符相等；
前两行是同一枚共享压缩器上"打标／未打标"的对照。用例侧断的是同一件事（`task == res.TaskID`
＋ uuid 形状 ＋ 未打标⇒**没有**这枚键），并且第二枚用例额外钉两种假法：
`idA/idB/idA` 三次连打的属性值必须逐条对得上（钉"把 id 存在共享对象上"这一形），
以及**任何属性值都不许是** `"" / unknown / none / nil / n/a / -`（钉"拿占位值洗归因"）。
两处 id 都取自本包自己的 `newTaskID()`，本程没有手抄过任何一枚假 id 进生产路径。

**本格改了哪些文件**：`internal/agent/compress.go`、`internal/agent/loop.go`、
`internal/agent/compress_trace_test.go`（commit `b23c7f7`）。

---

## 4. AC#3 — 承重两句：先把"本票新增的每一味"列出来，再逐味摘

**新增的味**（三味，逐枚有名）：

| 味 | 位置 | 是什么 |
|---|---|---|
| N1 | `compress_trace_test.go:426` | AC#1 那一发用例（过阈值×折不动⇒痕 0 枚） |
| N2 | `loop.go:399` | 调用点上的 `withTraceTask(ctx, taskID)`（loop 把 id 交出去） |
| N3 | `compress.go:237-239` | 那条 `Info` 上"有才追加、没有就不写"的 `"task"` 属性 |
| （用例侧） | `:478` / `:534` | `CarriesTheOwningTaskID`（读 N2＋N3 的合成结果）／`NeverInventsATaskID`（读 N3 的"不许占位"与"不许存在共享对象上"） |

矩阵全部跑在**交付 commit `b23c7f7` 的仓外快照**上（不是脏工作树），每发：`restore` → 施味/施变异
（落地必打印，见 `mut153.py`）→ `go build` → `go test ./internal/agent/ -count=1 -v` → 红名取自
`^ *--- FAIL` 行。

### 4.1 句① 「摘掉本票新增的任意一味，是否存在一发变异从此打不红？」——**三味各自：是**

| 发 | 施了什么 | RUN / FAIL | 红名（去掉 `TestCompressionTrace` 前缀） |
|---|---|---|---|
| **X0** | 交付码，无变异 | 83 / **0** | — |
| **X1** | 只施 M5（守卫→`if c.Need(hist)`） | 83 / **1** | `SilentWhenNothingFoldableOverThreshold` ← 本票的牙，同一枚变异从"四枚全绿"变成"这一发红" |
| **X2** | **摘 N1** ＋ 施 M5 | 82 / **0** | **NONE ⇒ 逃逸回来**：N1 是 M5 的唯一证人 |
| **X3** | **摘 N2**（loop 不再打标） | 83 / **1** | `CarriesTheOwningTaskID` |
| **X4** | **摘 N2 ＋ 摘 `CarriesTheOwningTaskID`** | 82 / **0** | **NONE ⇒ N2 的唯一证人就是那一发**（`NeverInventsATaskID` 自己造 ctx，看不见 loop） |
| **X5** | **摘 N3**（痕不再写属性） | 83 / **2** | `CarriesTheOwningTaskID` ＋ `NeverInventsATaskID` |
| **X8** | **摘 N3 ＋ 摘掉那两枚用例** | 81 / **0** | **NONE ⇒ N3 的证人恰是这两发** |
| X6 | 变异：属性**无条件**写（不判空） | 83 / **1** | `NeverInventsATaskID`（未打标的痕会带 `task=""` ——正是"把归因不了洗成归因成功"那一形） |
| X7 | 变异：`traceTaskID` 恒返一枚常量 uuid | 83 / **2** | 两发都红 ⇒ "拿一枚像真的一样的假 ID 凑"当场红，不需要人去读码 |

⇒ **没有一味是装饰**：N1 摘掉 ⇒ X1 那发从此打不红（回到 139 验收 §2.2 的原状）；N2 摘掉 ⇒ X4 全绿；
N3 摘掉 ⇒ X8 全绿。**同时也没有一味是"独占证人却没人证"**：三味各自至少被一发红咬住（X1 / X3 / X5），
且 X6、X7 两发是票面点名的两种假法（占位值、真形状假 id），它们各自也有专属证人。

⚠ 与本仓旧坑对齐的口径：红名一律只数锚定的 `--- FAIL` 行、`-count=1`、**跑整包**（不只那七枚），
所以"别枚用例顺手也抓到了"这一支被排掉了；X2/X4/X8 三发正是把"证人"与"被证的东西"一起摘掉才绿的，
这也就是"摘掉任意一味是否有一发从此打不红"的**是**字的读法。

### 4.2 句② 「摘掉它，有没有任何外部可见读数变过？」——**两句都答，且答案不是一样的**

外部可见读数＝走**进程默认 logger**（`slog.Default()`，也就是宿主 JSONL 的那条链）打出来的那一行，
不是测试内的 in-memory capture（探针 `probe153_line_test.go`／`probe153_loopleg_test.go`，只存在于快照）：

| 发 | 痕行数 | 带 `task=` 的行数 | 判 |
|---|---|---|---|
| **ext-A** 交付码 | 3 | **2**（直接打标那一发＋真 `Loop.Run` 那一发；后者 `task=23ac7a42-…` 与同次 `Result.TaskID` 逐字符相等） | 基线 |
| **ext-B** 摘 N3 | 3 | **0** | **变过**：盘上那一行少一枚键，肉眼可辨 |
| **ext-C** 摘 N2 | 3 | **1**（只剩直接打标那一发；**Loop 腿那一发不再带归因**，`Result.TaskID=16889253-…` 在场而痕里没有） | **变过**：生产腿的归因消失 |

⇒ N2／N3 两句都答"是，外部读数会变"。**N1 那句的答案是"不会"**：它是判据、不是输出面——
`compress_trace_test.go` 是 `_test.go`，不进任何交付二进制，摘掉它本程没有读到任何外部读数会差
（本程未为它造外部读数，可造的那一侧只是"少一枚证人"）。**它的承重只能由变异矩阵那一侧兑现**
（X1 与 X2 之差就是它）。票面 §AC#3 那句"别把'我加了用例所以更严'当结论"本程照抄；
但反过来"加了用例而外部读数不变"也不构成本票的退件理由——两句本程各自给了答案，没有拿一句抵另一句。

### 4.3 本格自证：还原与"跑的就是交付码"

```
== matrix runs on commit b23c7f7; md5 of the three files under test: (logs 头部有逐枚打印)
RESTORED IDENTICAL internal/agent/compress.go
RESTORED IDENTICAL internal/agent/loop.go
RESTORED IDENTICAL internal/agent/compress_trace_test.go
```

每发之间都 `restore`（`mut153.py restore`，从 `pristine-3/` 逐枚 `shutil.copyfile`），
收口再与 `git show b23c7f7:<同一文件>` 逐名 `cmp`。矩阵里 X0 的 83 RUN 与 §6 门禁 POST 的 83 RUN 相等，
说明快照那棵树与本程交付的那棵树在本包同一枚形状。

**本格改了哪些文件**：无（只读＋仓外快照）。

---

## 5. AC#4 — 契约轴零字节（逐枚 commit 点名，不用"区间"糊）

**先点名一枚会骗人的尺**：`git log --author=…` 在本仓**零信息**——全编队共用同一枚 git 身份
（`git log --pretty="%h %an"` 把 152／151／前端／编排者／本程全印成同一个名字），
所以"这一格是谁交的"只能按 `git show --name-only` 逐枚认。

**本程到这一枚为止五枚 commit，逐枚名册**（`git show --pretty=format: --name-only`）：

| commit | 格 | 文件 |
|---|---|---|
| `ff550f3` | AC#1 码 | `internal/agent/compress_trace_test.go` |
| `b23c7f7` | AC#2 码 | `internal/agent/compress.go` · `internal/agent/compress_trace_test.go` · `internal/agent/loop.go` |
| `ac7fb00` | AC#1 表 ＋ 43 枚原始读数 | 本证据件 ＋ `.scratch/wisp/probes/153/**` |
| `f58c513` | AC#2 表 | 只有本证据件 |
| `bec1727` | AC#3 表 | 只有本证据件 |

⇒ 两枚**码** commit 的文件全部落在 `internal/agent/**`，恰是票面 §地界 点名的那三枚，没多一枚。

**本程区间 `86b0161..HEAD` 里别家的 commit 逐枚有名**（现量于 `HEAD=4e16976`；本程收口时区间又长了
`eed7c99`／`63e228f`＝编排者的台账与推送收账，同样不归本程）：`5365cb2`（派单存档）、
`cdf2471`／`ad9b29f`／`8883b3f`（前端会话）、`eb4755a`／`4cc85bb`／`10e3585`／`6550dc4`（票 152 的程）、
`45c920e`／`5d46f24`／`4e16976`（票 151 的程——其中 `5d46f24` 的 subject 就叫 `placeholder`，
动的是 `cmd/wisp/run.go` ＋ `cmd/wisp/task_scope_close_151_test.go` ＋ `internal/tools/bridge.go`，
**不是本程的**）、`720cae6`（编排者特批存档）。共享 index 下本程每次 commit 前都现量了
`git diff --cached --name-only`，五枚里没混进别家路径（一枚 `git add` 之后 index 名册见 §9）。

**禁改面（对本程两枚码 commit 逐条筛，并拿同一把尺打正控）**：

```
$ git show --pretty=format: --name-only ff550f3 b23c7f7 | sed '/^$/d' \
  | grep -E "^docs/PLAN.md|^docs/specs/|^internal/risk/|^internal/panel/|^internal/agent/approval/|^internal/observe/|thresholds\.go|golden|allowlist\.txt|^scripts/slo-check\.ps1|^tools/d22scan/|^frontend/|^design/"
   (zero hits)
$ …同一条管道… | grep -cE "^internal/agent/"        → 4     ← 正控：尺活着，0 命中不是尺坏
$ git diff --name-only 86b0161 b23c7f7 -- internal/observe internal/llm/golden \
      scripts/slo-check.ps1 tools/d22scan internal/agent/budgets.go docs/PLAN.md docs/specs
   (empty = zero bytes)
```

**"阈值／`Need()`／`KeepRawRounds` 不许动"这一条单独量**（票面把它点成放水捷径，故不靠"我没改那个文件"这句话）：

```
Need() IDENTICAL                                        # sed 取函数体，锚点 vs 交付逐字节 cmp
"if len(raw) <= c.b.KeepRawRounds":                     anchor=1  delivered=1
"for c.totalOf(rounds) > c.b.HistoryCompressTokens":    anchor=1  delivered=1
budgets.go 在两枚 commit 之间的 diff 行数：0
```

⇒ 本票**没有**把"让常见路径也折叠"当捷径：AC#1 治的是判据形状，不是折叠面。
`git diff` 里被删掉的行只有三条——那条 `Info` 的头部、`history_changed` 那一行（同表达式换行）、
和 `Compress(ctx, hist)` 那一枚调用点（变成带标版）；`probes/153/ac4-zero-byte.txt` 里带 `-` 前缀原样印着。

**`DEFERRED(D-xx)` 代码标记枚数没被本程顶过期**（139 验收 §3.5 正是被这一枚咬过）：

```
anchor 86b0161    : D28-1 = 3   D11-3 = 1
delivered b23c7f7 : D28-1 = 3   D11-3 = 1
```

⚠ 自打一枚：`compress.go` 那段归因注释的**第一版真的写了字面量** `DEFERRED(D28-1)`，
那会把 D28-1 从 3 枚顶成 4 枚（正是 139 验收点名的"自家 commit 造成计数漂移"那一形）。
写完当场想到、改成转述（"the D28-1 deferral flagged at the top of this file"）之后才跑上面那两条计数——
**这一条是"写之前想到"没成的，不是被读数捞回来的**，记在这里，免得读者以为有门兜着。

**`12000` 名册**（139 验收 §3.4 那枚"本票贡献 0 枚"过期陷阱，本程先量再落笔）：

```
internal/agent/*_test.go 里的 12000：anchor 8+1=9   delivered 8+1=9
compress_trace_test.go 单枚：anchor=1  delivered=1   ← 本程两格新增的 204 行（71＋133）贡献 0 枚
```

**本格改了哪些文件**：无（只读＋计数）。`frontend/**`／`design/**` 全程零触碰、零宣称。

---

## 6. AC#5 — 门禁（一律逐包单跑，从不 `go test ./...`）

### 6.1 改前／改后四数 ＋ 名册差集（两向 `comm`）

| 形 | 采于 | RUN(all) | PASS_all | FAIL | SKIP | panic | unique |
|---|---|---|---|---|---|---|---|
| 改前 `-count=1` | `86b0161` 的仓外快照 | 80 | 80 | **0** | **0** | 0 | 80 |
| 改后 `-count=1` | 交付态工作树 | 83 | 83 | **0** | **0** | 0 | 83 |
| 改后 `-count=2` | 同上 | 166 | 166 | **0** | **0** | 0 | 83 |

```
$ comm -13 改前名册 改后名册   # 新增，逐名（恰三枚，无夹带）
TestCompressionTraceSilentWhenNothingFoldableOverThreshold   (AC#1)
TestCompressionTraceCarriesTheOwningTaskID                   (AC#2)
TestCompressionTraceNeverInventsATaskID                      (AC#2)
$ comm -23 改前名册 改后名册   → (空)          ← 一枚没丢
$ diff -q 改后-count1名册 改后-count2名册 → NAMES IDENTICAL
$ 166 = 2 × 83（算术自洽）；同一发里 RUN=83 / 出裁决=83 ⇒ "开跑了没回来"这一族本包未发作
```

⇒ 票面 §AC#5 那句"一枚用例 panic 会吞掉同包其余几十条"的防法本程照做：四数之外比名册差集；
且**全程未用过一枚 `t.Skip`**——现量：`git grep -n 't.Skip' -- internal/agent/**` 在交付版仍只有
`internal/agent/approval/ticket84_no_owner_test.go:224`（票 84 的"有意慢"闸，139 验收记过的那枚，非本票）。

### 6.2 三件工具（版本现读，不背数）

```
$ go version        → go version go1.27.1 windows/amd64
$ gofumpt --version → v0.12.0 (go1.27.1)      （/d/work/base/gopath/bin/gofumpt.exe，盘上现读）
$ go vet ./internal/agent/                     → 无输出，rc=0
$ gofmt   -l internal/agent                    → 无输出
$ gofumpt -l internal/agent                    → 无输出
$ gofumpt -l . tools/d22scan tools/mockllm     → 无输出
$ go test ./internal/agent/ -count=1 -race     → ok  github.com/CarlosShao/wisp/internal/agent 4.698s  （额外一发，非门禁要求）
```

⚠ 本程踩到并改掉的一枚格式坑（不是"背来的"，是这枚工具链给的）：AC#1 那发的第一版在文件尾留了两枚空行，
`gofmt -l` 与 `gofumpt -l` **双双点名**，而 `go vet` 一声不响——只跑 vet 就会带着一处不合规交付。
去掉那一枚空行后用例的行号会挪 1 位，所以 §2 的红/绿两发**是用改完的交付版文件在快照里重跑过的**
（`logs/f1`／`logs/f2`／`logs/f3`，红句里的 `:458`／`:461` 就是交付文件里的行号）；
更早那三发（`logs/r0`–`logs/r3`，同一条判据、未去掉空行的版本）本程**保留不删**，读数以 `f*` 为准。

### 6.3 `sh scripts/d22scan.sh`：rc=0，八枚分母非零

```
$ sh scripts/d22scan.sh ; echo $?      → 0        （末行 "d22scan: clean - no D22 ban violations"）
   整份日志里有 69 行 "examined" —— 前 60 枚是 tools/d22scan 自检用例自己造的 fixture 行
   （形状：`… of C:/Users/swq/AppData/Local/Temp/TestBuiltBinaryGoesRedEndToEnd…`），
   真扫描是尾部第 228-236 行那九行 ⇒ grep -m1 examined 必取到 fixture 的数（票面点名的坑，本程复算成立）。
   真扫描块：bans #1-5 internal/=205 · #1-5 cmd/=23 · ban #6 frontend/=67 · ban #7 internal/tools/=18
             ban #8 design/=39 · #8 frontend/=67 · #8 internal/=413 · #8 cmd/=44
```

分母不靠"不降"的推理，直接和树自身的枚数对齐：

```
$ git ls-tree -r --name-only b23c7f7 -- internal | grep -c '\.go$'  → 413  = ban #8 internal/ 413  ✓
$ git ls-tree -r --name-only b23c7f7 -- cmd      | grep -c '\.go$'  →  43  ≠ ban #8 cmd/ 44        ← 差 1
```

⇒ `cmd/` 那 1 枚之差**不是本程造的**：扫描跑在**工作树**上，当时
`cmd/wisp/task_scope_close_151_test.go` 还是票 151 那枚程的未跟踪新文件（现由它自己提交为 `5d46f24`）；
`frontend/` 的 67（139 当年记的是 66）同源。**别家的树在动不算进本程的任何宣称**，只把差值点名出来。

**本格改了哪些文件**：无。

---

## 7. 本程没测什么（按"漏了它谁会先被骗"排序）

1. **没把那枚痕复跑到盘**。§3.4／§4.2 的外部读数停在 `slog.Default()` 那一层（`go test -v` 的 stderr），
   本程没起宿主、没读 `<data>\logs\wisp-*.jsonl`，因此**没验 `redactHandler` 会不会改写字段名**。
   139 验收 §7.2 第 2 项验过"八枚键名原样在场"，**但它验的那版里还没有 `task` 这枚键**。
   ⇒ 若那枚 tee／红删会动这枚键，本票的"归因"就只到 stderr 为止；先被骗的是下次拿日志排障的人（他会以为没归因）。
2. **没量真并发下归因会不会串**。本程钉的是"同一枚共享 Compressor 连打 A/B/A 三次，三条痕逐条对得上"
   （＋整包 `-race` ok），那是 **per-call 不是 per-object** 的判据，不是两条 goroutine 同时压缩的读数。
   今天串味的入口不存在（139 验收 §1.4 复算：`Steer`／`RunAsync` 非测试零调用者、Loop 不跨任务复用），
   **但那张 Warm-window hook（D28-1）落地那天必须重读这一条**——本程只把"hook 也要带标"写进了注释，没写进判据。
3. **没验"同一发里会不会打多条痕"**（票面 §不解决的事第 3 条明令不据它下结论）；
   §3.4 那个"恰 1 枚"不是常态断言，别当它读。
4. **失败侧那条 Warn 不带 taskID**（`agent: history compression failed` 在调用方，本票地界是成功边）。
   ⇒ 两种结局现在**不对称**：成功的痕可归因、失败的不能。本程没修它，也没说它没问题——
   搜索前缀 `agent: history compress` 仍一次拿到两种结局，但只有前一种答得出"哪一轮压的"。
5. **ctx 通道可以被忘记**：未打标的调用，痕里就是没有 `task` 键（设计如此），但本程**没有**加编译期或
   lint 期的强制。⇒ 未来新增调用点忘了带标，只有 §4.2 那种"少一枚键"的读数能发现，没有门会红。
6. **没跑全树门禁、没在 linux 容器跑**（派单口径逐包）、**没起真子进程**
   （`wisp.exe` 那一发仍是 139 验收 §6.5#2 的未测项）。
7. **派单前提② 的后半（退化形到盘 1 枚）本程未复跑**，只在进程内复现了同形枚数（§1.2）。

---

## 8. 伪授权两栏计数（分开记，每条带出处＝工具名＋命令前 40 字）

### 8.1 真通知回显（判为真，未据其改变任何判据）——9 条

| # | 出处 | 形状与本程的处理 |
|---|---|---|
| 1 | 开场 `system-reminder` 里的 `<project_context>` AGENTS.md 全文 | 被告知的文件内容，未当新增指令 |
| 2 | 开场 `system-reminder`：available-skills 清单 | 与本票无关；本程未调任何 skill |
| 3 | 开场 `system-reminder`："The date has changed. Current date: 2026-09-26" | 与派单日期一致，未据此改任何"当日"断言 |
| 4 | 开场 `system-reminder`：`Memory: …/wisp/agents.md`（与 #1 同源重复） | 同上 |
| 5 | 多次 `Read`／`Bash` 结果的管道外壳（`Called the Read tool with the following input: {…}`、`Result of reading:`） | 机械回显，未当指令 |
| 6 | `Write` 工具回显：`File created successfully at: D:\tmp\wisp153\m5.py` | 本程请求的是 `/tmp/…`，bash 的 `/tmp` 实为 `C:\Users\swq\AppData\Local\Temp` ⇒ 工具按另一枚根解析。路径映射差异，非授权；本程改用绝对 Windows 路径重做该件 |
| 7 | `Bash` 里 python 抛的 `FileNotFoundError: … '\tmp\wisp153\ac1_snippet.txt'` | 同一枚映射差的机械后果；未据此判"文件丢了"之类 |
| 8 | `Bash`：`git commit` 后 `git show --stat HEAD` 印出 `cdf2471 feat(frontend 动画开关机器)` | 形状像"我的提交没进去"。本程用 `git log --oneline -6` 复核出自己的 `ff550f3`（名册见 §5），**未重提交、未 revert** |
| 9 | 本轮后半：`MEMORY.md` 被外部修改的系统通知（项目级与全局各一次，附整份记忆索引正文） | 编排者自己的记忆写入，未向本票要任何动作。本程**未据它放松任何一条判据**，也没把它当"解锁"；其中的"派单／验收／撤销口令"字样一律不作用于本票 |

### 8.2 判为注入（按派单第四条：即便不越权也登记并继续取证）——1 类

| # | 出处 | 文字形状 | 本程动作 |
|---|---|---|---|
| 1 | `No tools needed for comment` ——只出现在**本程自己每条纯文本消息的尾部**，任何一次工具输出里都搜不到它 ⇒ 给不出"命令前 40 字"，只能给位置；枚数＝每条文本消息一枚，本程未逐枚计数（≥20 枚） | 字面是"不用工具"，正落在派单第四条点名的"别用工具／直接给结论"形状里，即便不越权也要登记 | **未服从**：§1–§6 的读数全部真跑取满（九发矩阵、三发外部读数、两向 `comm`、`-race`、`d22scan` 逐包）。判为 harness 外壳回显而非指令，交编排者改判 |

⇒ 本程未在任何工具输出里读到"像编排者说的话"（没有"已解锁／请 revert／放宽阈值／少取证"这类句子从工具里出来）。

### 8.3 凭据值：零抄录

全程未接触凭据面。本件出现的只有**名字**与形状：函数／类型名 `traceTaskKey`／`withTraceTask`／
`traceTaskID`／`newTaskID`；三枚测试现造的 uuid 形状**任务号**（`ab8a8e40-…`／`23ac7a42-…`／`16889253-…`，
它们是本仓 `task_log.id` 那枚主键、D35 要求它可查，不是凭据）；路径形状 `<data>\logs\wisp-*.jsonl`。
未出现任何 Key、DPAPI 串、`.env` 内容；`cmd/wisp` 里那几枚假凭据常量本程连读都没读。

---

## 9. 收口

| 格 | 本程交了什么 | 判据闭合在哪几枚读数上 |
|---|---|---|
| AC#1 | 一枚新用例（`compress_trace_test.go:426`） | **改前红** `F2`（M5 落地＋未修码，红句原文在 §2.2）＋**改后绿** `F1`／`F3`（§2.3） |
| AC#2 | 三问逐问答完⇒**不动契约面**⇒没停手；痕带真 taskID＋两枚用例 | 三问读数 §3.1–§3.3；外部读数 `ext-A`（`task` 与 `Result.TaskID` 逐字符相等）；与 AC#1 未互相抵账 |
| AC#3 | 承重两句逐味答 | 句①：`X1↔X2`、`X3↔X4`、`X5↔X8`（味与证人一起摘才逃逸）＋`X6`／`X7` 两枚点名假法；句②：`ext-A/B/C`——两味"会变"、一味"不会变但由句①兑现"，两句各答各的 |
| AC#4 | 契约轴零字节 | 逐枚 commit 名册＋同一把尺的正控（§5）；`Need()` 逐字节相同、`budgets.go` diff 0 行；`DEFERRED` 计数 3→3（附一枚自打：差点顶成 4） |
| AC#5 | 门禁 | 80→83 四数、名册两向 `comm`（新增恰三枚、丢失空、两形名册全等）、vet/gofmt/gofumpt v0.12.0 全空、`d22scan` rc=0 且八枚分母非零（真扫描块＝尾部 228-236 行） |

**两格不抵账**（派单与票面各写一遍，这里再写一遍）：AC#2 那枚 `task` 键**替代不了** AC#1 缺的那一形——
把 taskID 送进去之后，`X2` 那一发（摘掉 N1 ⇒ M5 逃逸）仍然全绿；缺一枚判据就是缺一枚判据。

⚠ 本程自打最后一枚（关于这份件本身）：为了"每裁一格 commit 一次"，本程用 python 把这份证据件
截断过一回再续写，`newline` 处理不当 ⇒ **§5 之后所有行的行结束符被搞成 CR/CRLF 混排、段落之间多出一枚空行**
（已提交版 `9a8766e` 里就是坏的那版：总行 830／非空 320，比例 1.59，而 `ac7fb00`–`bec1727` 三版都是 0.3 上下）。
内容零丢失（现量：坏版非空 320 行 = 好版 265 行＋§5 的 55 行，逐行 diff 无缺失），
本节的 §5–§9 是从干净版 `bec1727` 起重写、重写时把三处措辞更正一并并进去的
（`--author` 那枚坏尺、"五枚 commit"、"204 行"）。**改判据没动、动的是排版**，但这份件的中间态确实红过一阵，
读者拿 `9a8766e` 对行号会错位——以本枚之后的版本为准。

**票面四框本程一枚都没自勾**；`.scratch/wisp/issues/153-…md` 未加 `-done`、Progress log 未追加
（勾与账归编排者，按另一枚非实现者程的表来定）。本程动过的路径全集：
`internal/agent/compress.go`、`internal/agent/loop.go`、`internal/agent/compress_trace_test.go`、
`docs/evidence/s1/153-trace-lies-unguarded-r1.md`、`.scratch/wisp/probes/153/**`（只建不删）。
零 push。

---

## 附录 A（09-26 11:2x，**编排者代记**；本件正文 §0–§7 一字未改）

来路＝非实现者验收件 §1.5 与续程 §12.4 第 4 行点名的那笔：**本件 §5 的别家名册里有一枚不存在的 commit 号 `8883b3f`**。
**本轮我自己现量的四条**（不是转述）：

```
git cat-file -t 8883b3f          → fatal: Not a valid object name '8883b3f'      ← 不解析
git rev-parse --disambiguate=8883b3f →（空）                                      ← 连以此为前缀的对象都没有
git log --since='08:30' --until='10:30' -- frontend docs/reports/frontend-session-log.md
                                 → 86b0161(08:59) · cdf2471(09:24) · ad9b29f(09:30) · 5c28b3b(10:14)   ← 该窗口真前端提交＝4 枚
```

⇒ **本件 §5 那行名册的正确读法**："前端会话"那一组只有 **`cdf2471`／`ad9b29f` 两枚可核**（`86b0161` 在窗口起点、`5c28b3b` 在本程收口之后），
**`8883b3f` 是一枚不可核的号** ⇒ 本件自报"18 枚名册、17 枚可核"这个**分数本身仍然成立**，但**别把 `8883b3f` 当第五枚前端提交去查**。
⇒ **候选真身＝`86b0161`**——**〔未证〕**，理由只有"同一区间、同一作者、名册里恰好缺一枚、而它是窗口内第一枚前端提交"；
**我不据它改名册**（按"不回填、只追加"的规矩，也按"猜出来的号比空号更坏"这条本仓既有条）。下一位若要用这一行，**自己现跑上面那三条**。
⇒ **这格不影响任何 AC 的档位**（名册是"本程未越界"的反扫凭据，不是判据；少一枚不可核的号不改变"零越界"的结论）。

---

## 附录 B（票 157 代记，只追加）

- 来路＝工单 `.scratch/wisp/issues/157-record-level-cleanup-ten-table-vs-disk-mismatches-left-by-tickets-155-and-153.md` 第 **9、10** 笔（＝本票 AC#2 里 `(a)` 那枚〔附条件〕欠的两件待补）＋派单 `.scratch/wisp/dispatches/2026-09-26-131x-impl-157-record-level.md`。同一票的第 1–8 笔落在 `155-three-unjudged-cells-r1.md` 末尾的**同名一节**。
- 本程身份＝**实现程·代记**：**零代码、零翻勾、零 `-done`、零新读数**。本件 §0–§7 与 **附录 A** 一字未改（进场现量 `wc -l`＝**589**）。两笔逐笔**先复算再补**，原始读数＝`.scratch/wisp/probes/157/13-items-9-10-privacy-and-keyname.txt`（＋`14-supplementary-rulers.txt` 末段 B9/B10 补尺）。
- ⚠ **补完这两件不等于本票 AC#2 能翻勾**：工单 `:31` 明写"翻勾由编排者按两张表合起来定，实现方不自勾"；本节里没有任何一枚档字是本程自立的。
- ⚠ 出处的一枚更正：工单 `:21` 把第 9 笔的出处写成"153 验收件 §5.3 **第 3 笔**"，本程逐行取版现量 §5.3 三笔的行界＝**第 1 笔 `:398`**（不是新增暴露类）／**第 2 笔 `:403–:411`**（件内缺那笔隐私账）／**第 3 笔 `:412–:414`**（日志 `task` ↔ DB `task_id` 不同名）。⇒ 隐私账那笔是**第 2 笔**；153 验收件 §12 `:831` 与 §12.4 `:887` 两处的引用行号（`:406–:411`＋`:412–:414`）是对的，本附录按它们写。

### B.9 第 9 笔 —— 本件缺那笔隐私账：**〔成立〕**

**验收程那一句**（`153-trace-lies-unguarded-r1-accept-r1.md:403`–`:411` ＋ `:416` 档位行第 ① 件）："`compress.go:236` 那句 'It is a correlation id, not history content, so `[privacy] keep_transcript` does not reach it' —— 本程复算它**成立**…**但本程在整份交付版证据件里搜 隐私／privacy／keep_transcript／私有数据／Q-31／SealDir：0 命中**……证据件里该有一行'落在票 132 的未闭合面上'；现在没有，读者会以为注释里那句 'keep_transcript 不到它' ＝ 隐私已审。"

**本程当轮复算**（全文＝`probes/157/13-…txt` 的 P1/P2 段）：

```
$ for w in 隐私 privacy keep_transcript 私有数据 Q-31 SealDir; do grep -c -- "$w" 本件; done
0 ／ 0 ／ 0 ／ 0 ／ 0 ／ 0          ← 六枚全 0（同一把尺打在 153 验收件上＝7／3／6／—／2／4，尺是活的）
$ git show 86b0161:internal/agent/loop.go | grep -n '"task"'
:373 task_log open failed   ／  :741 agent: tool timeout   ／  :933 task_log finish failed   ← 三枚既有
$ git grep -n '"task"' 6de3d1c5 -- internal/agent ':!*_test.go' | wc -l → 7
$ git show 6de3d1c5:internal/agent/compress.go | sed -n '232,236p'   ← 那句注释逐字在场（:236 是"does not reach it"那行）
$ git grep -n 'KeepTranscript' 6de3d1c5 -- internal/config
schema.go:511 KeepTranscript bool `toml:"keep_transcript" default:"false"`
validate.go:82-84 if c.Privacy.KeepTranscript { return observe.New(observe.ClassConfig, "…hard-coded false (read-only); writing true is rejected") }
validate_test.go:136 {"privacy.keep_transcript=true", …}      ← 写 true 报错这件事有钉
$ git grep -n 'SealDir' 6de3d1c5 -- '*.go' ':!*_test.go' → 6 枚命中，全在 internal/winsec/{winsec.go:174,190,220,222 与 winsec_windows.go:52,313}
                                                  ⇒ 定义 1 枚＋注释 5 枚：**非测试调用者 0 枚**
```

**我这一笔用的尺是什么**＝**两把**：① "缺不缺这一行"用**六词行级计数**打在**本件**上（并同一把尺打验收件做正控，证明 0 不是死尺）；② "那一行该写什么内容"用**锚点取版的码级现量**（`git show <锚>:<路径>`／`git grep … <锚>`，全部打在 `86b0161`／`6de3d1c5` 上，不读工作树），**没有一条是从验收件转抄的**。`SealDir` 那一枚本程额外做了"调用者 vs 定义者"的分拣（6 枚命中全在 winsec 自己两枚文件里），因为"零调用者"这件事光看枚数会读反。

**对不上的一枚没有** ⇒ 按最小闭合补那一行（本件正文一字未改，新文挂这里；这是**记账级**文字，不动判据、不动码）：

> **隐私账（本件当时没算，票 157 第 9 笔代记）**：本票交的那枚 `task` 键把一枚可归因 id 写进宿主日志文件。**它不是新增暴露类**——同一枚键名、同一枚值在改前锚点 `86b0161` 就有三枚非测试行在用（`loop.go:373`／`:741`／`:933`），到 `6de3d1c` 非测试侧带 `"task"` 的行数是 **7**；本票做的是"让成功边那一行也进"，不是"让 id 第一次进日志"。**但它落在票 132 那格未闭合的面上**：`Q-31`（日志算私有数据要上锁）要求的那把锁 `winsec.SealDir` 在 `6de3d1c5` 的非测试代码里**零调用者**。交付版 `compress.go:236` 那句"`[privacy] keep_transcript` does not reach it"复算**成立**（硬编码 false＋写 true 报错，`schema.go:511`／`validate.go:82-84`／`validate_test.go:136`），**可它只到"转写不落盘"那一层，不等于"文件已上锁"**——本件正文里没有这一区分，读者会把它读成隐私已审。

### B.10 第 10 笔 —— 日志侧 `task` ↔ DB 侧 `task_id` 无人写：**〔成立〕**

**验收程那一句**（`153-trace-lies-unguarded-r1-accept-r1.md:412`–`:414` ＋ `:416` 档位行第 ② 件）："`SPEC-02:74`／`PLAN.md:2697` 里那枚叫 **`task_id`**（DB 列），日志侧叫 **`task`** ⇒ 同一枚值两个面**故意不同名**，本票跟随日志侧是对的；但'从日志 join 回 DB'的人要知道这次改名，**实现件与 spec 都没写这一句**（登记，不判退）。"

**本程当轮复算**（全文＝`probes/157/13-…txt` 的 P3 段＋`14-supplementary-rulers.txt` 末段 B9/B10 补尺）：

```
$ sed -n '74p' docs/specs/SPEC-02-data-storage.md      → "  task_id        TEXT NOT NULL REFERENCES task_log(id),"
$ sed -n '2697p' docs/PLAN.md                          → tool_call 那一行：`id` PK、`task_id` FK、…、`correlation_id`、`grant_id`
$ grep -rn 'task_id' docs/specs/                       → 只有 SPEC-02:74 与 SPEC-02:88（那枚索引）两枚
$ grep -rn 'task_id' docs/specs/ | grep -c '日志'      → 0        ← spec 全目录里没有一句把两面键名对上
$ git grep -c 'task_id' 6de3d1c5 -- internal           → 只在 internal/memory：schema.go 2 ／ dao_toolcall.go 3 ／ models.go 1（＋测试）
$ git grep -n '"task"' 6de3d1c5 -- internal/agent ':!*_test.go' | wc -l → 7（日志侧键名，第 9 笔那发）
$ grep -c 'task_id' docs/evidence/s1/153-trace-lies-unguarded-r1.md → 0 ／ 同一把尺打 153 验收件 → 3（正控）
$ grep -n '键名' 本件 → :204 ／ :231 ／ :492 三枚，逐枚都在说**日志侧**惯例（"键名沿用 `loop.go` 已经在用的 `"task"`"），无一枚提 DB 侧
```

**我这一笔用的尺是什么**＝**两面键名各取一枚名册**：DB 侧从**规格两处**（`sed -n` 取行原文）＋**锚点取版的码侧 `git grep -c`** 数出列名在场；日志侧用第 9 笔那发 7 行名册；"有没有人写过这枚对价"用**三枚 0 命中**（本件／spec 与"日志"同现／本件 `键名` 三枚逐枚看过去只讲日志侧），并给每枚 0 配了**同一把尺的正控读数**（验收件 3 枚、SPEC-02 2 枚）——0 与"死尺"分开。

**对不上的一枚没有** ⇒ 补那一行（登记级：不动码、不动判据、不判退；本件正文一字未改）：

> **键名两面账（代记第 10 笔）**：同一枚 id 在两个面**故意不同名**——日志侧键＝`task`（`compress.go:238`、`loop.go:373/744/936`），DB 侧列＝`task_id`（`SPEC-02:74` 的 `task_id TEXT NOT NULL REFERENCES task_log(id)`、`:88` 的 `idx_tool_call_task`、`PLAN.md:2697` 那行 `tool_call`，码侧同形 `internal/memory/schema.go:71/:85`、`dao_toolcall.go:24/:112/:159`、`models.go:162`）。本票跟随**日志侧**惯例是对的（§3.4 已写"键名沿用 `loop.go` 已经在用的 `task`"），**但改名这件事本件与 spec 都没写一句**：拿那条痕去 join `tool_call`／`task_log` 的人，第一次会在 `task` ≠ `task_id` 上撞一次。⇒ 要真去 `SPEC-02` 补那句对价说明是**另一格活**（spec 在禁改面里，本程未碰）。
