# 票 47 — 装配腿（ASSEMBLY）：把真 C20 Scheduler 接进 cmd/wisp 生产装配行

- **腿名**：47 装配腿（实现腿；非实现者裁决表另出，本表只记盘上证据与门读数）
- **日期**：2026-11-07 09:1x +08
- **HEAD**：`9dd310cde`（票 47 实现腿 `1cc67664b`→`4656ec715` 已入库；本腿写面＝`cmd/wisp/**`＋`docs/evidence/s1/47-assembly.md`）
- **禁令遵守**：零 GUI、零 wisp.exe、零 *winlive*；`docs/reports/pending-and-issues.md` 只读一字未动；`frontend/**`、`design/**` 零写面；`internal/tools/**` 本腿一字未改（TaskDeps.Scheduler 与 NewBridgeDeps 均为实现腿已入库物）；无 push；未触碰 `cmd/wisp/panel_resident_*_323*` 两枚他人未跟踪件。
- **环境**：CGO_ENABLED 未置（全程非 0）；GOTMPDIR/TMPDIR/TEMP/TMP 全程 unset；PATH 前置 `build` 与 `third_party/sherpa-onnx`；go1.27.1 windows/amd64。

## 0. 缝的现状与本腿修的是什么

- 票 47 实现腿把 `tools.TaskDeps.Scheduler`（`SchedulerDeps{View,Control}` 两个窄面）、`task.list` 实现、task.cancel 调度器分支、`scheduler.NewBridgeDeps` 适配器全部落了库；但 `cmd/wisp/run.go` 装配行的 `TaskDeps{Roster, Paths}` 没带 Scheduler 字段 ⇒ 生产 task.list 一直答 fail-closed「调度器未接线」，task.cancel 的调度器分支也永不触发。
- bridge.go 自己的 disposition note 写明装配归 cmd/wisp 写面：「the scheduler is created there, the adapter is wrapped around the bridge, the entries register through the existing TaskDeps seams」——本腿就是那个装配半。

## 1. 装配点（全部在 `cmd/wisp/run.go`，逐字三处）

1. `agentRuntime` 新增 `scheduler *scheduler.Scheduler` 字段（组合根独占生命周期）。
2. 装配块（`rt.tasks = tools.NewTaskRoster()` 之后）：`rt.scheduler = scheduler.New(scheduler.Options{})`；`BuiltinTaskEntries` 的 TaskDeps 补 `Scheduler: scheduler.NewBridgeDeps(rt.scheduler)`。Options 零值＝生产形（Registry=nil→observe.Default，与全进程 D38b 名册同一本）；PathLock 刻意不接——生产还没有任何锁的获取方，一把没人 acquire 的锁探不到任何冲突，只会假装接了线（提交注释逐字记录）。
3. `close()`：`rt.scheduler.Close()`（D38e 第 1 步「关门」半；常驻腿自己的 step-1 钩子与其任务源门不受影响）。

生产此刻还没有向调度器 Submit 的入向通道（那是后续票的写面），所以生产空队列的 task.list 答「队列为空：当前没有任何任务。」——真调度器的如实答案，不是 fail-closed 拒绝；task.cancel 对调度器行的父权停法也已可达（数据 join＝SetParent）。

## 2. 门读数（全部本腿现跑，主树）

| 门 | 命令 | 读数 | 判 |
|---|---|---|---|
| gofumpt | `gofumpt -l cmd/wisp/run.go cmd/wisp/task_scheduler_wiring_47_test.go` | 空（rc=0） | ✅ |
| gofmt | `gofmt -l` 同两件 | 空 | ✅ |
| vet | `go vet ./internal/tools ./cmd/wisp` | rc=0 | ✅ |
| build | `go build ./cmd/wisp` | rc=0 | ✅ |
| test 主 | `go test ./internal/agent/scheduler ./internal/tools -count=2 -timeout 20m` | `ok scheduler 1.293s`＋`ok tools 63.395s`（-count=2 两遍全绿） | ✅ |
| test 新 | `go test ./cmd/wisp -run 'Test47Assembly' -count=1 -timeout 20m` | `ok 8.345s`：两枚全 PASS | ✅ |
| d22scan | `sh scripts/d22scan.sh`（Git Bash） | `clean - no D22 ban violations`；bans #1-5 internal/=373、cmd/=58、ban #7 tools/=48、ban #8 internal/=926、cmd/=222 | ✅ |
| 主树净度 | `go vet ./cmd/wisp`（含两枚未跟踪 323 件一同编译） | rc=0——323 件不破编译，故未启用 clone --shared 隔离 | ✅ |

## 3. 新增接线尺（`cmd/wisp/task_scheduler_wiring_47_test.go`，走 runFixture/rtHook 真桥）

| 用例 | 钉什么 | 读数 |
|---|---|---|
| `Test47AssemblyTaskListAnswersTheRealQueue` | 空队列 ⇒ 正常成功＋「队列为空」，回执零「未接线」 | PASS（回执逐字：「调度器队列为空：当前没有任何任务。」） |
| `Test47AssemblyTaskListSeesARowAndCancelStopsIt` | 经**生产 scheduler 实例**真 Submit＋SetParent ⇒ task.list 回执含该行（id＋标题＋行数 1）；task.cancel 以父任务身份（经 admitTask 的宿主身份与 L1 窗口，票 375 形）成功停掉 | PASS（回执逐字：「调度器任务 1 行…9b56563e…『票47装配校验任务』 状态=Thinking」→「已对调度器任务 9b56563e…发出停止」） |

L1 答复源形＝`f.reply = strings.NewReader("")`＋`pinL1WindowCeiling220`（与 panel_l1_blocked_220_test.go 同款；窗口到点无人否决即执行，是 L1 路自己的语义，非给 cancel 开后门）。

## 4. 变异试验（有牙证明）

**突变**：装配行 TaskDeps 摘掉 `Scheduler: scheduler.NewBridgeDeps(rt.scheduler)` 一行（其余不动）。

| 尺 | 突变下读数 | 判 |
|---|---|---|
| `Test47AssemblyTaskListAnswersTheRealQueue` | **FAIL**：回执＝「调度器未接线（fail-closed：读不到任务行，不代表没有任务在跑）」，IsError=true | ✅ 有牙——红句正是 fail-closed 原文 |

**还原**：字段行补回后 `git diff --stat -- cmd/wisp/run.go`＝本腿净改动（+32/−1），重跑 `Test47Assembly` 两枚全 PASS。

## 5. 契约与边界

- 契约零改动：D34 行（task.list L0／task.cancel L1）、C20、D38e 均按 PLAN.md 冻结文字消费；`internal/tools/**` 零字改动。
- 未接 PathLock 的代价与落点：调度器行的 HeldPaths 恒空、路径冲突排队（D43 #20/#39）暂无可达触发方——登记为后续 Submit 入向票的写面，不是本票缺陷。
- 已知边界（如实）：生产 Submit 入向通道未落，故生产 task.list 今天只会答「队列为空」；「行可见」已由 T2 经同一枚生产实例证真。
