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
