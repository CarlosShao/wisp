# 票 222（丙形）——等待孩子的那段时间不占桥的执行许可：r1 交付证据件

- 写腿：`222-r1`。开工定锚 `49abb593`（`git log -1 --format=%h` 现取）；本件首次落盘时刻 `2026-09-29 12:12 +08`，`date` 现跑。
- 射程＝工单 `.scratch/wisp/issues/222-spawn-holds-a-bridge-permit-while-waiting-for-its-child.md` 的**丙形**（等孩子的这段时间不占许可）。甲形（`task.spawn` 立即返回句柄）**未做**，等 owner 批。
- ⛔ 本件不新造规矩：D1–D47／C1–C32／SLO／golden／`thresholds.go`／`PLAN.md`／`docs/specs/**` 一字未动。桥的 `MaxToolConcurrency = 4` 与池的 `MaxConcurrentSubagents = 4` **一枚都没改**（§④ 有两枚常量的修复后现读）。

---

## ① 现读链条（行号全部 2026-09-29 12:0x 现跑 `grep -n`，非引用任何归档快照）

| # | 事实 | 现读出处 |
|---|---|---|
| 1 | 桥在调用工具**之前**取许可、整个调用期间握着 | `internal/tools/bridge.go:438-440`（`select { case b.sem <- struct{}{}: defer func() { <-b.sem }()`）；池 `:171` `sem: make(chan struct{}, ceiling)`；常量 `:24` `MaxToolConcurrency = 4` |
| 2 | `task.spawn` 是桥内工具，`Execute` 里阻塞等孩子 | 注册 `internal/tools/subagent_197.go:218`；阻塞 `:356-368` 那个 `select { case <-ctx.Done(): … case res := <-done: … }` |
| 3 | 等待用的 ctx 带 per-tool 超时 | `bridge.go:449` `ectx, cancel := b.execContext(ctx, timeout)`（建在取许可**之后**）→ `:470` `entry.Tool.Execute(ectx, …)`；`timeoutFor` `:550-555`；`TaskSpawnDecl()`（`subagent_197.go:202-211`）没有 `Timeout` 字段 ⇒ 生产走 `cmd/wisp/run.go:562` 的 `DefaultTimeout: time.Duration(cfg.Agent.PerToolTimeoutMS) * time.Millisecond`，未配置＝`bridge.go:28` `DefaultToolTimeout = 30 * time.Second` |
| 4 | 孩子在生产上用同一枚桥递工具 | `cmd/wisp/run.go:581` `ParentTools: rt.bridge` ＋ `subagent_197.go:277` `opt.Tools = newSubagentToolProvider(t.d.ParentTools)` |
| 5 | 孩子的 ctx 没有 deadline | `subagent_197.go:327` `context.WithCancel(context.WithoutCancel(ctx))` |
| 6 | 池帽＝桥帽＝4 | `subagent_197.go:78` ＋ `bridge.go:24` |
| 7 | 测试看不见这一形：孩子拿的是假目录 | `internal/tools/subagent_197_test.go:183` `ParentTools: h.dir`、`:186` `Tools: h.dir`、`:159` `dir: &fake197Dir{}`；父侧 `:216` 却是 `h.bridge.Execute`。全文件 `grep -n "Tools: h.bridge"` ＝ **零命中**（现跑） |

⚠ **对票面 §"推出来的真实形状"的一处现读更正**：父任务真正收到的**不是** `subagent_197.go:360-365` 那段「父任务这一侧已经不等了」文本，而是**桥把它换成了超时文本**——`bridge.go:475-483` 在 `entry.Tool.Execute` 返回**之后**先查 `errors.Is(ectx.Err(), context.DeadlineExceeded) && ctx.Err() == nil`，命中就丢回「工具 task.spawn 超时（1500ms），已协作式中止」。§② 那发红名的原始读数就是这句。结论形状不变（`IsError: true`、父侧拿不到孩子结论），但**用户/模型看到的字面是桥的超时文案**，票面与 `A420` 那句"收到的错误文本就是父任务这一侧已经不等了"按现读应改成上面这句。

## ② 修复前的读数（三枚新用例，产码＝HEAD 未动，`-count=1`）

时刻 `2026-09-29 12:11 +08`（`date` 现跑），命令：

```
export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"
go test ./internal/tools -run 'Test222' -v -count=1
```

逐名结果：

| 用例 | 修复前 | 用时 |
|---|---|---|
| `Test222SpawnConclusionArrivesThroughRealBridgeChildren` | **FAIL** | 1.50s |
| `Test222WaitingParentHoldsNoBridgeSlot` | **FAIL** | 0.00s |
| `Test222CeilingStillCapsExecutedCallsWhileParentsWait` | **FAIL** | 0.00s |

关键原始行（逐字，节选）：

```
subagent_222_test.go:407: 等待孩子的父任务占着 4 枚桥位，want 0：全桥 4 枚，父任务一边干等一边占位＝票 222 现量链的形状，孩子只能排在 sem 后面，直到父任务自己的 per-tool 超时松绑
subagent_222_test.go:446: 等待孩子的父任务占着 4 枚桥位，want 0（反控的前半：先把许可还不回来）
subagent_222_test.go:353: 孩子 c1 的工具调用起跑时已有 2 枚父任务返回，want 0：它是等父任务自己超时才挤上桥的
subagent_222_test.go:353: 孩子 c2 的工具调用起跑时已有 4 枚父任务返回，want 0：它是等父任务自己超时才挤上桥的
subagent_222_test.go:360: 父任务 0 收到的是错误而不是结论："工具 task.spawn 超时（1500ms），已协作式中止"
subagent_222_test.go:378: 结论 "子代理 1 的结论正文 222" 出现在 0 枚父任务回复里, want 1
FAIL	github.com/CarlosShao/wisp/internal/tools	1.532s
```

读数解释（每一枚都问过了"把修复拿掉它会不会红"）：
- `Test222WaitingParentHoldsNoBridgeSlot` 直接读 `len(bridge.sem)`：等待中的 4 枚父任务**占着 4 枚**，全桥就是 4 枚 ⇒ 孩子一枚都拿不到。这一发 **0.00s 判红**，不靠超时、不靠挂死。
- `Test222CeilingStillCapsExecutedCallsWhileParentsWait` 同一枚读数做前半，后半的反控（同时执行数＝4、且宿主调用起跑时 `len(parents) == 0`）在修复前被前半直接 `Fatalf` 挡住。
- `Test222SpawnConclusionArrivesThroughRealBridgeChildren` 量的是**内容**：孩子的工具调用确实起跑不了（`probe.runs` 与起跑时刻父任务已返回数），父任务最终收到桥的超时文案、正文里没有它自己孩子的结论。

（完整原始输出见 §⑤ 的"修复前"那一发，逐字粘贴。）

## ③ 修复形状（丙：等孩子期间不占许可；不动任何冻结数字）

<!-- 待补：bridge.go 的 inFlightSlot 载体 + subagent_197.go 的等待点交接；为什么是单向交还、为什么不新增导出名、为什么不是 SubagentDeps 字段 -->

## ④ 修复后的读数

<!-- 待补：同名三枚用例的绿名 + len(bridge.sem) 的读数 + 两枚常量的现读 -->

## ⑤ 门与原始输出

<!-- 待补：go build ./... / go test ./internal/tools ./cmd/wisp -count=1 / gofumpt / d22scan / -list 计数 -->

## ⑥ 没做的格与卡点

<!-- 待补：AC#3 甲/乙 都要 owner；AC#4 属票 221 同文件串行；AC#5 的定性；AC#6 终态 -->

## ⑦ 新增名字具名清单

<!-- 待补 -->
