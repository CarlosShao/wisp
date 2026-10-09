# 282-v1 · 按住的 AC#1 与 AC#3（派单约束，⛔ 不是本腿能力问题，也⛔ 不是本腿的欠账）

## 为什么按住（串行铁律）

- 这两格的判据都**必须跑 Go 测试并在 `internal/tools/**` 磁盘上种突变**：AC#1 要把 `internal/tools/cancel.go` 的 `TaskID(ctx)` 种坏（返回空／回退到 corr）后看指名用例红不红；AC#3 要"改读访问器前后"各一发读数，而"前"那一发只能靠把 `internal/tools/subagent_197.go:262`／`internal/tools/task.go:708` 的读者**改回** `CorrelationID(ctx)` 再跑一次得到（没有第三条取数路径——同一二进制里编不进两个版本的同一枚符号）。
- 起手一刻在飞＝`242-v2`（HEAD `7789f953` 即其交件，写面 `internal/agent/approval/**`，同属 Go 导入图）。同一时刻只许一条腿在导入图里种突变 ⇒ 我种的坏会被它编进测试二进制、它种的坏会洗掉我的读数，两边读数互不可信。
- 本腿实到动作面（可核）：**零 `go` 命令、零 `scripts/*.sh`、零 `.go` 字节改动、零票面改动**；写面只有 `.scratch/wisp/probes/282/v1/` 三枚 `.md`。⇒ 这两格在票面里保持未勾，⛔ 不许按"本腿没做"记失败。

## AC#1 缺的读数（具名，四件套逐件列明缺哪一半）

> 四件套＝① 种前 hash ② 改的那行 `sed` 复量 ③ 红句 `file:line` ④ 还原 hash 逐字回＋`git status --porcelain` 空。下面每一发都缺**全套四件**。

- **缺第 1 发（空形）**：种 `internal/tools/cancel.go` 的 `func TaskID(ctx context.Context) string` ⇒ `return ""`。指名必须红的候选面（本腿只做静态名册，未跑）：`internal/tools/ticket283_corr_identity_rulers_test.go:40`（单元尺）／`:194`／`:206`（整链尺）；`internal/tools/subagent_197_test.go`（`task.spawn` 侧）；`internal/tools/task_cancel_221_legs_test.go`（`task.cancel` 侧）。⇒ 读数缺失：**到底哪几枚红、红句原文**。
- **缺第 2 发（回退形，真正验"有没有牙"的那一发）**：同处 ⇒ `return h.corr`。这一发的意义＝区分"读的是 taskID"与"读的是 corr 的回落"，而 `bridge.go:528` 的 `corr: orDefault(req.CorrelationID, req.TaskID)` 加上夹具普遍 `TaskID == CorrelationID`，正是最容易被洗掉的一形。⇒ 读数缺失：**红／全绿**；若全绿 ⇒ 具名写"该支今天没有仪器"。
- **本腿额外具名的一发（静态扫出来的松面，交派单方决定要不要种）**：`internal/agent/loop.go:603-607` `callCorr` 的 **id-less 回落支** `fmt.Sprintf("%s#call-%d", taskID, index)`。静态尺＝`git grep -n 'callCorr' HEAD -- '*.go' ':(exclude).scratch'` ⇒ 除定义与调用点外**零断言**（`corr_percall_242_test.go:139` 与 `loop_approval_test.go:220` 都只钉 `HasPrefix`，不钉形状）。⇒ 可种形状：让 index 跨轮复用（同任务第 2 轮第 0 枚与第 1 轮第 0 枚同值）⇒ 两枚调用共用路由键。**读数缺失：会不会有用例红**。这一发同时是 AC#2 里 S-4"更松"那一行的补牙，归 AC#1 的射程。

## AC#3 缺的读数（具名，"同一形状的前后两发"）

- **缺第 3 发＝深度判定那一跳的前后对**：同一形状（父任务派生子任务、断言"深度"读数）跑两遍——
  - 后（现状）：`internal/tools/subagent_197.go:262` `parentID := TaskID(ctx)` 一发。
  - 前（还原写法）：同处改回 `parentID := CorrelationID(ctx)` 一发。
  ⇒ 缺的读数＝**这两发在同一 fixture 下的 PASS/FAIL 对**。判的是"静默失效"到底被修掉了、还是只是换了写法。
- **缺第 4 发＝父子停机那一跳的前后对**：同一形状（父停子）跑两遍——
  - 后（现状）：`internal/tools/task.go:708` `caller := TaskID(ctx)`。
  - 前：同处改回 `caller := CorrelationID(ctx)`。
  ⇒ 缺的读数＝**同样两发**。
- ⚠ 方法学提示（省得下一枚踩）：**夹具必须让 `TaskID != CorrelationID`**，否则前后两发必然同绿。静态证据：`internal/tools/ticket283_corr_identity_rulers_test.go:5-16` 的抬头注释已具名"既有用例的派发 fixture 一律 `TaskID == CorrelationID`，分不出两者"；`internal/tools/subagent_197_test.go:217`／`subagent_carrier_197_test.go:173`／`task_cancel_221_legs_test.go:119`／`failclosed_236_teeth_test.go:145` 逐枚实测两枚 id 同值 ⇒ 拿这些夹具跑前后对**得到的"全绿"不是读数**。可指的真尺只有 `ticket283_corr_identity_rulers_test.go`（`039ec93c`，227 行新建、零删除），而它是**票 283 的产物、晚于 `bd124b2a`** ⇒ AC#3 那句"改读访问器前后"在 242 交付那一刻**根本没有可用夹具**，这一点可以作为 AC#3 的静态结论先记，但**红绿读数仍缺，本腿不代判**。

## 本腿已替这两格消掉的部分（免得重复劳动）

- AC#1 的"指名用例名册"已静态交齐（上面 3 枚候选文件＋行号）。
- AC#3 的"有没有可指夹具"已静态答：**242 交付时没有，283 才补上**（凭 `039ec93c` 的 `git show --stat` 与 `git log --oneline -- internal/tools/ticket283_corr_identity_rulers_test.go` 名册）。
- AC#2 的**形状表**已经把空／旧形等值／无前缀三形判红（纯布尔，不需种坏）⇒ AC#1 的两发只剩"验证测试二进制真会红"这一层，⛔ 不是重新判方向。
