# 197 证据件 — 「行已发布、取消句柄未挂」那枚窗口（票 70 step 7 定因的修复）

> **本件是什么**：票 70（CI 红分诊）step 7 逐字读出、判为「真缺陷（窗口型）」的那枚时序缺陷的**实现腿**证据：
> 复核、先造红、修法选型、修后分母、四把门禁尺。
> **本件不是什么**：不是裁决件。裁决者≠实现者（`SPEC-12 §4.3` #1/#3、`issues/README` 硬约束末条、D22 双角色），
> 本腿只把「盘上可核的读数」摆出来，不替非实现者判 AC。
>
> - 起手 HEAD：`99ae90113`（分支 `dev`）；生成时刻 `2026-11-14 12:45 +08`。
> - 派单里那枚缺陷的定因件：`docs/evidence/s1/70-ci-red-triage.md` step 7（下称「triage」）。
> - 所有读数出自本腿亲手跑的命令，原始输出全在 `.scratch/wisp/probes/197pw/`（只建不删）。

---

## ① 复核：四处原文与行号（先核再动）

### 1.1 三处偏差（✍ 写在最前面，都不改机理）

| 派单／triage 写的 | 盘上实际 | 影响 |
|---|---|---|
| `internal/agent/subagent_197.go` | **`internal/tools/subagent_197.go`** | 写面所在包从 `agent` 纠正为 `tools`；`internal/agent/` 下从来没有这枚文件（`ls internal/agent/*197*` 空） |
| `internal/agent/task_cancel_221_legs_test.go` | **`internal/tools/task_cancel_221_legs_test.go`** | 同上；triage step 7 开头已自我更正过这一条，本腿复跑确认 |
| 派单里的 `legs_test.go` | 仓里**没有**裸 `legs_test.go`；它指的就是上面那枚 `task_cancel_221_legs_test.go` | 无 |
| triage 引的 `internal/agent/loop.go:360`（「第一次模型调用之前」那句注释里的行号） | 盘上 `AdmitTask` 的调用是 **`internal/agent/loop.go:409-414`**；`RunAsync` 里 `newTaskID` 是 **`:364`** | 只是**注释行号过期**，机理一字不差；本腿顺手把 `subagent_197.go` 里那两处引用改成盘上真行号（`:364`／`:409-414`），未动任何代码语义 |

**行号本身全部对得上**（判据：`git show HEAD:<file>` 逐行取，不采信任何转述）。

### 1.2 四处原文（逐字，行号＝HEAD `99ae90113`）

```
# internal/tools/subagent_197.go:337-345  ← 子 goroutine 侧：第一次模型调用之前「只发布行」
337: 	opt.AdmitTask = func(taskID string) func() {
338: 		t.d.Roster.PublishSubagent(taskID, parentID, label)
339: 		revoke := baseAdmit(taskID)
...
345: 	}
```

```
# internal/tools/subagent_197.go:373-374  ← 派生 goroutine 侧：发布之后才挂句柄
373: 	t.d.Roster.PublishSubagent(bg.ID, parentID, label)
374: 	t.d.Roster.AttachCancel(bg.ID, cancelChild)
```

```
# internal/tools/task_cancel_221_legs_test.go:132-158  ← 放行点＝第一次模型调用
132: func (x *harness221) spawnHeld(t *testing.T, label string, started chan struct{}) (childID string, parked chan Result) {
...
145: 		t.Fatalf("子代理 %s 没有开始它的模型调用", label)
...
158: }
```

```
# internal/tools/task_cancel_221_legs_test.go:452（`Test221SubagentCannotStopSiblingOrItself`，本枚从 :413 起）
450: 	stopA := x.cancel(t, parent197, aID)
451: 	if stopA.IsError {
452: 		t.Errorf("父任务停甲号孩子被拒（拒绝类正控不许把有权的那一支也一起堵死）：%s", stopA.Text)
# 同族 :322（`Test221ParentStopsItsOwnChildRowAndStreamSettle`，本枚从 :305 起）
322: 		t.Fatalf("父任务停自己的孩子被拒了：%s", res.Text)
```

### 1.3 机理（本腿复核后的读法，与 triage 一致）

- 行有**两个**写者：子侧 `AdmitTask`（`:338`，第一次模型调用之前，**只写行**）与派生侧（`:373` 写行、`:374` 挂句柄，**两次调用**）。
- `TaskRoster.Cancel`（`task.go:519-531`）在 `known==true && cancel==nil` 那一支答的是
  「没有在跑（kind=…，state=Thinking），没有可停的句柄」⇒ 同一行上「在跑」与「没在跑」同时被说出口。
- `spawnHeld` 的放行点是 `started`＝**第一次模型调用已经开始**。那一刻只保证**行**已由子侧 `:338` 落地；
  **不保证**派生侧 `:374` 已执行。负载一高，派生 goroutine 在 `RunAsync` 返回后被抢占一下，窗口就开。
- ⇒ 测试的判据（父任务停自己的孩子必须成功）是对的；病在**产品码的发布顺序**。**不改判据、不放宽断言**。

---

## ② 先造红（读数先于修，全部在 `99ae90113` 那棵树上取）

### 2.1 既有仪器在空载下打不开（与 triage 同结论，给分母）

```
go test -count=40 -run '^(Test221SubagentCannotStopSiblingOrItself|Test221ParentStopsItsOwnChildRowAndStreamSettle)$' ./internal/tools -v
⇒ === RUN 80 发 / --- PASS 80 / --- FAIL 0 / 「没有在跑」0 次
   （.scratch/wisp/probes/197pw/pre-plain-count40.txt）

go test -race -count=25 …（同上两枚）
⇒ --- PASS 50 / --- FAIL 0 / 「没有在跑」0 次
   （.scratch/wisp/probes/197pw/pre-race-count25.txt）
```

⇒ **本机空载打不开这枚窗口**（12 核、无并发负载），与 triage ③ 的 3/3 绿同源：不是 flake，是窗口窄。
所以本腿不拿它当红的证据，只当"照旧绿"的对照。

### 2.2 确定性红（形状）：`Test197SpawnPathPublishesRowAndStopHandleTogether`

这枚是**不靠调度**的红：它读的就是那两条语句本身。修前逐字红句（`.scratch/wisp/probes/197pw/pre-shape-red.txt`）：

```
--- FAIL: Test197SpawnPathPublishesRowAndStopHandleTogether (0.00s)
    subagent_publish_window_197_test.go:106: 派生路径没有调用 PublishSubagentInFlight：行与句柄必须由一次发布落地，否则「行可见 ⇒ 句柄已挂」就只是一个约定而不是不变量
    subagent_publish_window_197_test.go:115: subagent_197.go:338 仍在用裸 PublishSubagent 发行（这一调用不带句柄）：t.d.Roster.PublishSubagent(taskID, parentID, label)
        这一行就是票 70 step 7 读出来的窗口：行先可见、句柄后挂，中间那一瞬state=Thinking 而 Cancel 说「没有在跑」
    subagent_publish_window_197_test.go:115: subagent_197.go:373 仍在用裸 PublishSubagent 发行（这一调用不带句柄）：t.d.Roster.PublishSubagent(bg.ID, parentID, label)
        这一行就是票 70 step 7 读出来的窗口：行先可见、句柄后挂，中间那一瞬state=Thinking 而 Cancel 说「没有在跑」
    subagent_publish_window_197_test.go:120: subagent_197.go:374 仍把 AttachCancel 当第二次发布：t.d.Roster.AttachCancel(bg.ID, cancelChild)
        句柄必须和行进同一个临界区（PublishSubagentInFlight），否则两次发布之间就是一个可观察的中间态
```

### 2.3 症状红（分母）：`Test197InFlightRowNeverLacksItsStopHandle`

一枚复现仪器＝**40 轮 × 4 枚并发**（4＝`MaxConcurrentSubagents`，即"整个池同时派生"），每轮把 4 枚孩子都停在
**第一次模型调用**里，然后让**父任务**用**真 bridge／真 roster／真 `task.cancel`** 去停每一枚，读 `state` 与拒绝理由。
无 sleep、无 deadline、无耗时阈值（D22）；放行点仍是"第一次模型调用"。

```
空载一发：0 / 160 命中（.scratch/wisp/probes/197pw/pre-behaviour-concurrent.txt）

六路并发负载（同时 6 条 go test 跑 Test221ParentStopsItsOwnChildRowAndStreamSettle）下 3 发：
⇒ 1 / 480 命中（.scratch/wisp/probes/197pw/pre-behaviour-underload.txt）
```

**逐字红句**（第 3 发）：

```
--- FAIL: Test197InFlightRowNeverLacksItsStopHandle (0.12s)
    subagent_publish_window_197_test.go:129: 行已发布、句柄未挂：1/160 个在飞子代理被父任务停时回了「没有在跑」（CI 20 发里 2 发命中 = 10%，本地空载一时打不开；这就是那枚窗口）：
        第 18 发 窗口-18-3（5b0c8f82-deca-4ea0-9888-3c75efc6a4ea）：没能停掉 5b0c8f82-deca-4ea0-9888-3c75efc6a4ea：任务 5b0c8f82-deca-4ea0-9888-3c75efc6a4ea 没有在跑（kind=subagent，state=Thinking），没有可停的句柄（这一行没有被改动，别的孩子也没被牵连）
```

⇒ 与 CI 那两句**同一句拒绝文本、同一个自证矛盾**（`state=Thinking` vs 「没有在跑」）；
`kind=subagent` 那一格证明打到的是子代理行，`这一行没有被改动` 那半句证明打的是"真拒绝"而不是别的支路。

> 对照：同一轮六路负载里，**原始那枚 `Test221ParentStopsItsOwnChildRowAndStreamSettle` 跑 150 发全绿**
> （`.scratch/wisp/probes/197pw/pre-load-0..5.txt`）。⇒ 负载不是万能放大器：是"每轮 4 枚并发 + 读完立刻停"
> 这个形状把窗口放大的（每轮 4 个窗口 × 40 轮）。

### 2.4 为什么症状红只能给分母（诚实说明，不粉饰）

窗口的"关闭动作"是派生 goroutine 在 `RunAsync` 返回后的**下一条语句**（`:374`）；那两条语句之间**没有任何接缝**：
`SubagentDeps` 的每个字段（`Roster`／`BaseOptions`／`ParentTools`／`Provenance`／`Stream`）都在这段之外调用。
要让症状红在没有负载时也确定性打开，只有三条路，本腿全部拒了：

1. 在产品码里加 `time.Sleep`／墙钟门闩 —— **D22 明令禁**；
2. 在产品码里加一枚测试专用回调（`publishStep func()`）钉在 `:373` 与 `:374` 之间 —— 把测试脚手架焊进热路径；
3. 改 `spawnHeld` 让它"更早"去停 —— 停得再早也早不过"派生侧执行下一条语句"，改不动概率。

⇒ 走派单允许的第二支：**形状红（2.2，确定性）＋ 症状红（2.3，给分母 480）**，两个红都摆在修前那棵树上。
`spawnHeld` 一个字未动（派单允许的"仅当必须加确定性钩子"这一支**没有用上**：不需要）。

---

## ③ 修（选型与理由）

### 3.1 改了什么

| 文件 | 改动 |
|---|---|
| `internal/tools/task.go` | `Record` 的锁内体抽成 `recordLocked`（同一份"身份三格保留"的写，不复制第二份）；**新增 `TaskRoster.PublishSubagentInFlight(taskID, parentID, label, cancel)`**：`mu.Lock` 内**先写行、后挂句柄、再一次 Unlock**；`PublishSubagent` 保留但注释改成"只给没有句柄的行（夹具／兄弟行），派生路径禁用" |
| `internal/tools/subagent_197.go` | `context.WithCancel(...)` 提到 `AdmitTask` 之前铸造（钩子要用它）；`AdmitTask` 内 `PublishSubagent` → **`PublishSubagentInFlight(…, cancelChild)`**；派生侧 `:373/:374` 两次调用 → **同一次 `PublishSubagentInFlight`**；`agent.New` 失败路径补 `cancelChild()`（句柄已铸、无环可拥有它） |
| `internal/tools/subagent_publish_window_197_test.go`（新） | 2.2／2.3 两枚仪器；文件头写明"症状／形状"两半各自答什么 |

**没有改**：状态名与转移（`subagentStateRunning=Thinking`、`Settling`、`Muted`、`Error` 一字未动，D43 转移表零改动）、
`Cancel` 的措辞与分支、`DetachCancel`、`spawnHeld`、任何阈值／golden／断言。

### 3.2 为什么是这一形（选型理由）

1. **一次发布＝一个临界区**。`Cancel` 的判据是"行在 ∧ 句柄在"；把两句写在同一个 `mu.Lock` 里，
   「行可见 ⇒ 句柄已挂」就变成**表自身的性质**，而不是两个调用点之间的时序约定。这是唯一能把
   "可观察中间态"消灭掉的做法——两处都原子还不够，必须让"行与句柄同生"这件事只由一处代码负责。
2. **`AdmitTask` 当权威发布者**（票面要求，也是本腿复核后的选择）。它是全仓唯一**happens-before 第一次模型调用**
   的发布点（`loop.go:409-414` 在任何 provider/tool 调用之前）⇒ 修后的绿是**因果**，不是"这次运气好"：
   只要孩子在模型调用里，句柄就一定已经挂上。
3. **派生侧那次发布不删、改成同一方法**。原因：票 197 立的性质是"已启动的孩子必有行"（DSH/Step-Code 那类
   "孩子在跑、名册里没有它"的失败模式）。删掉它会让"派生返回到 admission 之间"与"控制词孩子（admission 之前就返回）"
   两处出现没有行的孩子；两条发布**同值同句柄**，last-writer-wins 幂等，谁也擦不掉谁。
   `Test197RowExistsBeforeFirstChildModelCall` 仍绿可证这条性质没丢（§④）。
4. **`nil` 句柄不发布**（fail-closed）：这方法的语义就是"让一行可停"；给一个停不了的行盖上 `Thinking`，
   正是本票要杀的谎。
5. **不把 `PublishSubagent` 直接改签名**：它的其余调用者（夹具、`sibling-221` 那种"谁都不许停我"的行）
   本来就没有真句柄；硬塞一个假句柄会让 `Cancel` 声称一次它做不到的停止——那比"没有句柄"更坏。
   这条拆分由 2.2 那枚形状尺钉住（派生路径禁用裸 `PublishSubagent`、禁用 `AttachCancel` 当第二次发布）。

### 3.3 代价（写清落点，不拿来降级）

- 行第一次可见的时刻从"派生侧发布"变成"派生侧发布 或 admission 发布（谁先谁算）"——**两者都在第一次模型调用之前**，
  且**两个写者都原子**，所以状态机可见语义（`Thinking`＝在飞）不变，只是"中途那一瞬"没了。
- `PublishSubagentInFlight` 在 `mu` 内多写一个 map：一次 spawn 两次，不构成任何热路径成本。

---

## ④ 验（修后读数）

### 4.1 同一枚仪器（分母与修前同形同量）

```
空载 1 发（40 轮 × 4 枚 = 160）：0 命中
   --- PASS: Test197InFlightRowNeverLacksItsStopHandle (0.04s)
   （.scratch/wisp/probes/197pw/final-instrument-idle.txt）

六路并发负载 3 发（480 枚）：0 命中
   --- PASS × 3（.scratch/wisp/probes/197pw/final-instrument-underload.txt）
   同一轮负载里 6 条 Test221ParentStopsItsOwnChildRowAndStreamSettle × 25 发：--- FAIL 0
```

⇒ **480 发 0 命中**（修前同条件 1/480）。且这不是"碰巧"：§3.2-2 的 happens-before 使它在
"孩子已进第一次模型调用"这个观测点上**结构性为真**。

### 4.2 形状钉（确定性那一半）

```
--- PASS: Test197SpawnPathPublishesRowAndStopHandleTogether (0.00s)
```

### 4.3 原两枚用例

```
go test -count=5 -run '^(Test221SubagentCannotStopSiblingOrItself|Test221ParentStopsItsOwnChildRowAndStreamSettle)$' ./internal/tools -v
⇒ --- PASS 10 / --- FAIL 0
   （.scratch/wisp/probes/197pw/final-original-two-count5.txt）
```

### 4.4 全包（四数）

```
go test -count=1 -v ./internal/agent/ ./internal/tools/
⇒ === RUN 767 / --- PASS 584 / --- FAIL 0 / --- SKIP 1
   ok  github.com/CarlosShao/wisp/internal/agent  39.119s
   ok  github.com/CarlosShao/wisp/internal/tools  81.166s
   （.scratch/wisp/probes/197pw/final-full-two-packages.txt）
```

> 数法：`^=== RUN` / `^--- PASS` / `^--- FAIL` / `^--- SKIP` 逐行数（`=== RUN` 含子测试行，故 767 与 584 不同口径，照实给）。

### 4.5 额外（不属门禁，顺手量的）

```
go test -race -count=2 -run '^(Test197|Test221)' ./internal/tools  ⇒ ok  4.359s（无 DATA RACE 报告）
go test -count=3 -run '^Test197' ./internal/tools（六路负载下）     ⇒ 39 枚全 PASS，含
   --- PASS: Test197RowExistsBeforeFirstChildModelCall（"已启动的孩子必有行"这条性质没丢）
```

---

## ⑤ 门禁四把尺（终态，`2026-11-14 12:45:26 +08`）

| 尺 | 读数 |
|---|---|
| `go build ./internal/... ./cmd/... ./tools/...` | **rc=0** |
| `go vet ./internal/agent ./internal/tools` | **rc=0** |
| `gofumpt -l internal/agent internal/tools` | **空**（0 行） |
| `sh scripts/d22scan.sh` | **`d22scan: clean - no D22 ban violations`**，rc=0 |

---

## ⑥ 边界、写面与没做的事

### 6.1 写面（4 枚）

- `internal/tools/task.go`、`internal/tools/subagent_197.go`、`internal/tools/subagent_publish_window_197_test.go`（新）、
  本证据件；另在 `.scratch/wisp/issues/197-…md` 末尾追加一条 Progress log（只追加不删）。
- **零改动**：`cmd/wisp/**`、`internal/panel/**`、`internal/observe/**`、`internal/agent/**`（连注释都没动——
  `loop.go` 的行号只写进 `subagent_197.go` 的注释里，没回改 `loop.go`）、`scripts/**`、`frontend/**`、`design/**`、
  任何阈值／golden／契约、D43 转移表。
- 已知行号偏差只在 ① 里报了，**没有顺手改别人的文件**。

### 6.2 没做／做不到（都在这里交代）

1. **裁决留给非实现者**：本件只给读数与选型理由，不勾 AC、不判"这格完成"。
2. **`cmd/wisp` 的测试在本机跑不起来**：`go test -run '^TestRunPacketMarksTheRosterRowACardIsHolding$' ./cmd/wisp`
   ⇒ `exit status 0xc0000135`（DLL 未暂存；triage step 13 记过同一形态 `:10369`）。本腿对 `cmd/wisp/**` **零写面**，
   `go build ./cmd/...` rc=0；那条读数是**环境**（DLL 暂存脚本），不是本改动的回归。
3. **工作树里有别人的脏文件**（本腿一枚未提交、未暂存）：`internal/observe/quota_44.go`、
   `docs/evidence/s1/44-quota-seedfrom-lock.md`、`.scratch/wisp/issues/70-ci-actually-green.md`（票 44 家族的活）。
   ⚠ 因此 §4.4／§4.5 的绿是在"含别人未提交 quota 改动"的树上量的；它们与本票写面无交集。
4. **Git**：`git commit -- <显式 pathspec>`；未 `add -A`／`add .`；未 amend／reset／rebase／stash／`checkout .`／clean；
   未建 worktree；**未 push**。
