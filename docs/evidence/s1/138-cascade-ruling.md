# 138 裁定书 r1 — 「级联」今天的形状、可走路径、AC#2 修法裁定与 AC#3 突变读数

- 工单：`.scratch/wisp/issues/138-cancelling-a-task-kills-nothing-the-task-itself-spawned-single-jobscope-no-terminatejobobject.md`
- 本件：实现腿（Wisp 编队写腿，票 138 步1–步5）的裁定与读数落盘；按派单要求**先裁再写、每步一 commit**。
- 时间：2026-11-07 15:4x–16:0x +08（每步读数现跑，不背旧数）。
- 被验基线：进园时 `HEAD = 48abebc41`（226 AC#5 revoke 语法），工作树 **tracked 零脏**。
- LF：本件与全部改动文件 LF；`git add`/`commit` 两处显式 pathspec；只 commit 不 push。

---

## §0 起手三核（现跑读数，2026-11-07）

### 0.1 ① 核票面点名的串/行还在不在

| 票面锚点 | 现跑命令 | 现读数 | 判定 |
|---|---|---|---|
| `internal/agent/loop.go:300` `func (t *RunningTask) Cancel() { t.root.Cancel() }` | `grep -n "func (t \*RunningTask) Cancel" internal/agent/loop.go` | **`:342`**，函数体逐字同形（只有 `t.root.Cancel()`） | 行号**腐烂 +42**；「Cancel 只有 context 取消、无进程动作」这一句**在 agent 层仍然为真** |
| `git grep 'type JobScope'` ⇒ 仅 `internal/proc/jobscope_windows.go:62` 一枚声明 | `grep -rn "type JobScope" internal/` | `internal/proc/jobscope_windows.go:62` 唯一，未变 | 成立（但见 §1.3：它已不是唯一的 Job Object 形状） |
| `git grep 'TerminateJobObject'` ⇒ 全仓 0 处 | `grep -rn "TerminateJobObject" internal/` | **2 处**：`internal/plugin/manifestslot/command_job_windows.go:146`（生产调用）＋ `command_138_cancel_test.go:12`（上一轮遗留测试件的注释） | **票面此锚已失效**：生产码里今天有 1 枚真调用 |

### 0.2 ② 核引文出处

- 票面引 GAP-09「JobScope 全进程一枚、无 TerminateJobObject ⇒ 任务级级联确无归属」——出处
  `docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md`，裁决 `…gap-analysis-audit-verdict.md` §2 第 09 行。
  该裁决判的是 **09-26 之前的世界**；引文本体无误，但「无走路径」的前件随票 52 落地而失效（§1.4）。
- 上一轮验收程 `docs/evidence/s1/138-cancel-kills-children-r1-accept-r1.md` §0 现读：被验版本 `3f6322f`、
  「AC#1 判"给未来功能预留的洞"、AC#2/AC#3 前件判假未触发」——与票面 Progress log 09-26 行逐字相符。出处核对通过。

### 0.3 ③ 红/绿是不是别人的在飞件造成的

- `git status --porcelain -- internal cmd docs tools scripts go.mod go.sum` ⇒ **空**（tracked 零脏）。
- 未跟踪文件里与本票相关的：`internal/plugin/manifestslot/command_138_cancel_test.go`（3657 B，
  **上一轮 138 腿的遗物**，自述「records the mutated line, the red line and the restored diff」于
  `docs/evidence/s1/138-cancel-cascade-ac2-ac3-r2.md`——**该证据件不存在**，即上一轮声称做了突变但只留下测试件、
  没留下证据与 commit，与派单「上一轮两条腿返回 null」相符）。
- 其余未跟踪件（`.scratch/**`、`cmd/wisp/*323*test.go`、`docs/evidence/s1/{195,351,378}*`）与本票零交集，
  本腿一枚不碰、不 commit。
- ⇒ 本腿读数基线干净；后来的红/绿不存在"别人的在飞 tracked 改动"这一干扰源。

### 0.4 上轮遗物的可用性初判

`command_138_cancel_test.go` 依赖的全部助手（`installFixtureCLI`／`stableExeDir`／`fixtureTool`／`runtimeFor`／
`readPIDFile`／`stillAlive`／`fixtureHandshakeVar`／`fixtureModePrefix`／`ModuleNameEnv`）在票 52 的
`command_52_test.go`／`command_52_main_test.go`／`command.go` 里**全部在盘**（行号见 §1.2 表）；
`def.TimeoutMs = plugin.MaxCommandTimeoutMS / 5 = 120000ms`（120s 预算 ⇒ 测试期内只有调用方 cancel 能结束调用，
预算不是演员——设计正确）。初判：**可编译、可采纳**；步3 以它为底改写并落 commit。

---

## §1 步1：「级联」今天到底是什么形状、可走路径在哪（现跑读数）

### 1.1 任务侧的 ctx 级联链（agent 层）

| 环 | 出处（现跑行号） | 内容 |
|---|---|---|
| 1 | `internal/agent/loop.go:342` | `func (t *RunningTask) Cancel() { t.root.Cancel() }` —— 用户按"停"＝取消任务根 ctx |
| 2 | `loop.go:589` → `loop.go:774` | `executeCalls` 把每个工具调用 spawn 在 `root` 之下：`l.reg.Spawn("tool-exec-"+taskID, "agent", root, …)` |
| 3 | `loop.go:855-863` | `dispatch`：`context.WithTimeout(ctx, timeout)`（C22 每工具预算）→ `l.opt.Tools.Execute(tctx, req)` |
| 4 | `internal/tools/bridge.go:595` → `:621` | 桥层 `run`：`ectx, cancel := b.execContext(ctx, timeout)`（= `context.WithTimeout(调用方 ctx)`，`:719-723`）→ `entry.Tool.Execute(ectx, req.Args, onUpdate)` |

⇒ **任务 ctx 死 ⇒ 工具执行拿到的 ctx 死**。链条前三环今天逐环在盘。

### 1.2 命令插件侧的进程级联（票 52 落地的那半，2026-11-03 commit `5438d104e`）

| 环 | 出处（现跑行号） | 内容 |
|---|---|---|
| 5 | `manifestslot/command.go:209` ＋ `stopwatch.go:21-24` | `withMonotonicBudget(ctx, budget)` = `context.WithDeadline(parent, …)` —— D46 单调预算，**挂在调用方 ctx 之下** |
| 6 | `command.go:218` | `proc, err := c.spawner.Spawn(childCtx, spec)` —— 子进程在预算 ctx 存续期间启动 |
| 7 | `command.go:231` | `kill := func() { killOnce.Do(proc.KillTree) }` —— 幂等杀树 |
| 8 | `command.go:232-238` | C11 scope 守望：`scopeCtx.Done()` ⇒ `kill()` |
| 9 | `command.go:243-251` | 预算守望：`childCtx.Done()` ⇒ `proc.MarkTimedOut()` + `kill()`。**`childCtx` 的 Done 有两个触发源：预算到期、以及父 ctx（= 任务 ctx）被取消** —— `context.WithDeadline(parent,…)` 的取消是自上而下传播的 |
| 10 | `command.go:517-532` | `osChild.KillTree()`：有 Job ⇒ `job.Kill()`；无 Job ⇒ 退化为杀直接子进程 |
| 11 | `command_job_windows.go:141-148` | `childJob.Kill()` ⇒ `windows.TerminateJobObject(h, 1)` —— **全仓唯一生产调用，`:146`** |
| 12 | `command_job_windows.go:31-44, 65-86` | 每子进程一枚 Job Object：`CreateJobObject` + `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`；子进程 `CREATE_SUSPENDED` → `AssignProcessToJobObject`（`:91`）→ `ResumeThread` —— 子进程第一条指令运行前已在 Job 内，孙进程随 Job 继承 |

### 1.3 Job Object 的归属形状（对票面 §0 的修正）

- `proc.JobScope`（`internal/proc/jobscope_windows.go:62`）今天**仍存在**，但调用者面是：球/面板/LLM 适配器测试
  （`internal/llm/adaptertest/childjob_windows.go:33` `proc.OpenJobScope()`）、停机步 9（`internal/proc/shutdown.go:88`）、
  `internal/tools/web.go`（C30 拉起前自检）。**它不再是"唯一的 Job Object"**。
- D46 命令插件另有一族**每子进程一枚**的 `childJob`（`manifestslot/command_job_windows.go`），
  归属单位＝一次命令调用，不是整机。⇒ 票面标题里「唯一的 JobScope 是整机一枚」对 `proc.JobScope` 仍真，
  对「仓里的 Job Object」已不真（票 52 之后每条命令子进程都有自己的 Job）。

### 1.4 可走路径判定（AC#2 的前件）

- 09-26 验收程判「今天无入口踩得到」时的世界：`shell.exec` 未注册、无生产 exec 路径。
- 今天：D46 命令插件**已落地并注册**（`internal/plugin/manifestslot/command.go` 是生产运行时，
  票 52 合同测试 17 条全绿），`Test52TimeoutKillsTheProcessTree`（`command_52_test.go:708`）已经钉住
  **预算到期**这半边的杀树。
- ⇒ **AC#2 前件「如果 AC#1 判定确有可走路径」⇒ 成立**。可走路径＝任务 ctx 级联（§1.1 四环）
  × 命令插件进程级联（§1.2 八环），两段在 `bridge.run` 的 ectx 处焊接。
- 今天真实存在的残余缺口（本票该修/该钉的）：
  1. **没有一枚"任务级取消（非预算到期）杀树"的用例**——预算那半边有牙（52 的 timeout 用例），
     取消那半边没牙；上一轮留下一枚没落盘的候选（§0.4）。
  2. **非 Windows 半边 `childJob.Kill()` 是静默 no-op**（`command_job_other.go:23`）——
     AC#2③ 名字点过的形状；`errNoJobObject` 在 `job_other.go:11` 有话（"degrades to the direct child"），
     但 `Kill()` 本体零字，`command.go:444` 的注释还写着 "a Windows Job Object" 而非"仅 Windows 有"。
- 138 r1 证据件 §AC#1③ 量过「builtin 工具族零子进程」，本件不重算；本件增量＝票 52 之后的世界。

---

## §2 步2：AC#2 裁定——修法形状与理由（先裁再写）

### 2.1 裁定

**不新建"每轮任务一枚 Job Object"，不新增任何级联通道；修法＝三件小实体：**

1. **钉牙（核心）**：新增 `Test138TaskCancelKillsTheProcessTree`（`internal/plugin/manifestslot/command_138_cancel_test.go`），
   钉住「调用方 cancel（= 任务取消的形状）⇒ 被包 CLI 的**整棵进程树**（含孙进程）死」。
   这是级联的取消半边的第一枚牙（预算半边的牙已由票 52 的 `Test52TimeoutKillsTheProcessTree` 提供）。
2. **非 Windows 半边去静默（AC#2③ 点名的形状）**：
   - `command_job_other.go` 的 `childJob.Kill()` 从静默 no-op 改为**留一行可见 stderr 告警**
     （"no Job Object on this platform; the tree kill degrades to the direct child"）；
   - `command.go` 的 Spawn fallback（`openChildJob` 报 `errNoJobObject` 后走 `cmd.Start()` 处）加同一句可见告警——
     这是非 Windows 上**真实可达**的路径，告警落在真实路径才算"可见"；
   - `command.go:442-444` 的 OSCommandSpawner doc 注释校准（"a Windows Job Object" → 准确表述两平台语义）。
3. **本裁定零改动面**：`internal/agent/loop.go`、`internal/proc/**`、`internal/risk/**`、`cmd/wisp/**` 一概不动；
   `cancelled` 文案与错误分类零触碰（AC#2④ 前件不成立，不触发停手）。

### 2.2 为什么不是"每轮任务一枚 Job Object"（票面修法族的字面）

- 票面写「修法只许落在"每轮任务一枚自己的 Job Object、取消时级联"这一族」——那是 **09-24 开票时**的形状，
  当时仓里零枚 `TerminateJobObject`、ctx 级联没有执行器。**票 52 落地后该族已经被更细的形状占据**：
  每条命令子进程一枚 Job（归属单位＝一次命令调用，比"每轮任务"更细），杀树执行器＝`TerminateJobObject` 在位。
- 在它之上再叠一枚 per-task Job 买到的新价值＝**0**（同一棵树、同一个取消源、同一次 kill），
  付出的复杂度是真金白银：一个进程同时挂 per-call Job 与 per-task Job 需要 Windows 8+ 的 Job 层级
  （`AssignProcessToJobObject` 嵌套、`KILL_ON_JOB_CLOSE` 在层级下的继承/breakaway 语义）——
  **复杂买不到真价值，不选**（铁律③：复杂本身不加分）。
- 「取消时级联」这半边今天已经存在且形状正确：Go context 的自上而下取消（§1.1 四环）
  × 预算守望的 `childCtx.Done() → kill()`（§1.2 第 9 环）。**它缺的只是证明**——这正是 AC#3 的牙。

### 2.3 为什么非 Windows 半边选"显式无此能力＋可见告警"、不选"等价实现"

- AC#2③ 给两条路：等价实现（Unix 进程组杀树）或明确声明无此能力＋可见告警。
- 等价实现（`Setpgid` + `syscall.Kill(-pid)`）是主流形状，但**本宿主（Windows）无法给它配牙**：
  AC#3 的标准是"突变转红、逐名记红行"——一个我无法在本宿主跑出红读数的实现是**没有牙的代码**，
  与本仓"不许一步存在却从不产出结论"的文化冲突。D7 只 ship Windows，非 Windows 半边今天的真实消费者
  是交叉 vet/build。⇒ 选显式声明＋可见告警；等价实现留给真有 Unix 交付面的一天（届时必须带牙）。

### 2.4 AC#2 四条件逐一过

| 条件 | 过法 |
|---|---|
| ① 不动 D4/C19 门控判定（risk 侧零字） | 修法面＝`command_job_other.go`＋`command.go` 注释/告警＋新测试件；`internal/risk/**`、`rules_gateway.go` 零触碰（§4 门禁读数复核） |
| ② 不新增重量级依赖 | 零新依赖：标准库 `fmt`/`os`＋已有 `golang.org/x/sys/windows` |
| ③ Windows API 留 `_windows.go`；非 Windows 有等价或可见告警 | Windows API 全部在 `command_job_windows.go`（未动）；非 Windows 从静默 no-op 改为可见告警（§2.1 第 2 条） |
| ④ 若要改 `cancelled` 文案/错误分类 ⇒ 停手 | 判定：不需要（级联已通，本票只钉牙＋去静默，零文案改动）——不触发 |

---

## §3 步4：AC#3 的牙与突变读数（★本票核心）

### 3.1 牙是什么

`Test138TaskCancelKillsTheProcessTree`（`internal/plugin/manifestslot/command_138_cancel_test.go`）：
真进程用例——fixture CLI（测试二进制自体）以 `hang` 模式被包运行，先拉起一个**孙进程**（写 pid 文件）再永久阻塞；
用例用**调用方 cancel**（120s 预算 ≠ 演员，只有 cancel 能结束调用），断言：

1. Exec 在 15s 内返回（`:87-95`）；
2. **孙进程 pid 已死**（`:97-109`）——只杀直接子进程的假实现恰好会留下这个 pid。

时间线：先证孙进程活着（`:70-82`，"从未咬过活树的牙证明不了任何事"），再 cancel（`:85`）。
它钉的是级联的**取消半边**；预算半边已有牙（`Test52TimeoutKillsTheProcessTree`，`command_52_test.go:708`）。

### 3.2 突变 M1：把级联摘掉（预算从父 ctx 上摘下）

| 步 | 读数（现跑，2026-11-07 16:17–16:20 +08） |
|---|---|
| 突变前尺 | `git grep -n 'context.WithDeadline(parent, total.Deadline())' -- internal/plugin/manifestslot/stopwatch.go` ⇒ **`stopwatch.go:23`** |
| pristine 存档 | `.scratch/wisp/probes/138/r1/stopwatch.go.pristine`（只建不删） |
| 突变落地 | `:23` `context.WithDeadline(parent, …)` → **`context.WithDeadline(context.Background(), …)`**（摘掉父链；预算仍会自己到期，但调用方 cancel 不再传播） |
| 落地证明 | `git grep -n` ⇒ `stopwatch.go:23` 新形状在盘；`git diff` 五行 hunk 在案；`go build ./internal/plugin/manifestslot` **rc=0**（票面 AC#3"先证变异落地再读数"过） |
| 突变下跑牙 | `go test -count=1 -run Test138TaskCancelKillsTheProcessTree -v ./internal/plugin/manifestslot` ⇒ **FAIL 15.244s，rc=1** |
| **红在哪一行** | **`command_138_cancel_test.go:94`** —— `t.Fatalf("Exec did not return within 15s of the caller's cancel: the cancel cascade never reached the child (stopwatch budget detached from its parent?)")`。机制：摘掉父链后 `childCtx.Done()` 对调用方 cancel 永不触发 ⇒ 预算守望的 `kill()` 不被调用 ⇒ Exec 永不返回 ⇒ 15s 守望开火。**这正是"被停的任务留下挂死的 CLI 在跑"的形状**（孙进程也还在，只是用例先死在第一断言） |
| 还原 | `Copy-Item pristine → stopwatch.go` ⇒ `git diff -- internal/plugin/manifestslot/stopwatch.go` **空**；`git status` tracked 干净 |
| 复绿 | 同命令 ⇒ **PASS 1.02s，rc=0**（16:19:58） |

### 3.3 过程干扰登记（三核③ 逐条）

- 突变期间 `go build ./internal/...` 一度红在 `internal/ball/orb_windows.go:246`（undefined: orbWindowEdgeBase）——
  **另一条在飞写腿**当时正在改 `internal/ball/orb_windows.go`＋`orb_live_windows_test.go`（` M` 在案，我零触碰）；
  数十秒后复跑 rc=0。突变判定不受影响：突变红绿读数全部取自 `./internal/plugin/manifestslot` 包作用域。
- 上轮遗留的未跟踪测试件（3657 B）已在步3 收编进 tracked（commit `bf224be19`），其证据件引用指向本裁定书 §3。

### 3.4 牙的约束自查（票面 AC#3 警句）

- 无 `t.Skip`、无阈值放水、无断言迁就：断言写的是"孙进程 pid 必须死"，突变后红因是**实现失去级联**，非用例松动。
- 牙的窗口参数：孙进程存活确认轮询 ≤10s、cancel 后返回守望 15s、pid 死亡确认 ≤2s（100×20ms）——全部在用例内可见。

---


