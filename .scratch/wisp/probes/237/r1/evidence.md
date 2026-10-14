# 票 237 落地腿 237-r1 交件读数（AC#1 甲／AC#3 甲／AC#2 注释；零产码）

- 写面＝`internal/tools/subagent_197_test.go`（+1 行注释）＋`internal/tools/subagent_222_test.go`（AC#1/AC#3）＋本目录。
  ⛔ `subagent_197.go`／`bridge.go` 一字节未动；三枚冻结件／`thresholds.go`／golden／`allowlist.txt`／`d22scan/**` 零触碰。
- 锚点：起手现跑 HEAD `ab6da884`（派单写 `58366328`，其间只有编排者文档提交；`internal/tools` 最后被 `8509e52e` 碰过）。
  改前 `git hash-object`：197_test `b4336141…`、222_test `7c64505e…`（`logs/00-anchor.txt`）。
- 每发 `-race -count=N` 起止时刻逐发在 §4；台面 09:08（另腿在飞）／09:40 两次 `tasklist`。

## 1. AC#1 —— 走甲：等待条件改成真正依赖的那件事（改前红已复现）

- 甲/乙选择理由：那条红不是"等得不够久"，而是测试等的不是它依赖的事——`parkParents` 收满 n 枚
  "孩子进了第一枚模型调用"后**立刻**把 `len(h.parked)` 当快照断言"父任务已发派生回执"；产码里
  `child.RunAsync`（`subagent_197.go:347`）**早于** `giveBackWhileWaiting`（:368）与 :371 的 delta ⇒ 无蕴含。
  乙（改成不依赖回执数＋归因文案）会把这份**因果等待**整格摘掉；甲保留全部读数，换成"每枚父任务各收一枚
  parked 签"（通道收信＝因果；真没在等的那一枚以自己的结果现身，不再是护栏）。
- 改前红形状（同一锚点，具名命令 `GOFLAGS= go test -race -count=60 -run Test222 ./internal/tools/`）：
  shot1 09:08:45→09:09:11 rc=1，`--- FAIL` **3 枚**（`:514`×2／`:437`×1）全为
  "只有 3/4 枚父任务到达等待点（父任务没派生成功…）"，`DATA RACE` **0**；shot2 rc=1 `--- FAIL` **1 枚**（`:514`）。
  日志 `logs/pre-fix-race60-shot1.txt`（770 B）／`shot2.txt`（383 B）。
- 改后（同命令，N=60≥20）：shot1 09:34:11→09:34:27 rc=0；shot2 09:38:13→09:38:25 rc=0；两发 0 FAIL／0 DATA RACE。
  全包 `go test -count=1 -race ./internal/tools/` 09:37:31→09:38:02 rc=0 `ok 24.773s`。

## 2. AC#3 —— 走甲：收紧成逐枚（撞出的独立真形已在测试面修好并留证）

- 实现 `awaitChildOnBridge222(t, n)`：读满 n 枚（本用例自己那组孩子）到达才返回；观察到的每枚父任务返回
  都先排空到达缓冲再判决（k 枚返回须有 ≥k 枚到达托底，缺口＝有孩子的工具调用没上桥）；收到的父结果**留到退出
  时一次性交回**。⚠ r1 自捉：第一版边收边交回会被 select 二次收取同一枚结果而错计（`-count=40` 曾 48 红全为此形），
  已在注释里具名这段坑。未走 strict（"任一返回早于 n 枚到达就红"）：慢而诚实的子代理可合法晚于别的父任务返回到达，
  strict 会拿调度噪声打红；只有"永远不到的到达"是违规。理由写进注释。
- 判据读数（`-race -count=1 -run Test222 -v`；变异只经 `-overlay`，⛔ 工作树零突变字节）：
  | 形状 | leg1 | leg2 | leg3 |
  |---|---|---|---|
  | (ii) 只有一枚孩子绕开真桥（`x-partial-after`） | 红 3.00s | 红 0.00s | 红 0.00s |
  | 同形 × 收紧前（`x-partial-before`） | 红 30.00s（旧护栏） | **绿 0.00s** | **绿 0.00s** |
  | 全量绕过（`x-full-after`，M3 回归） | 红 0.00s | 红 0.00s | 红 0.00s |
  | 摘掉 OpenTask 布线（`x-noopadmit-after`） | 红 3.00s | 红 0.00s | 红 0.00s |
  ⇒ (i)"逐枚要求不满足即红"与 (ii) 均落：leg2/leg3 从"失明全绿"变 0.00s 红，leg1 从 30s 护栏红变 3.00s 预读红。
- 收紧期间撞出的**独立真形**（新发现；落点判定＝测试面，⛔ 产码零改）：
  三角形归因链（三发 instrumentation overlay，只打印不判定；`logs/dump-leg2.txt` 31,294 B／
  `logs/dump5-leg2.txt` 68,968 B／`logs/dump6-leg2.txt` 19,652 B）：
  ① 红时 4 枚孩子全 Settling 而真桥只跑过 1–3 次 ⇒ 有孩子的调用**没执行**；
  ② `VERDICT237 … level L2 rules [R4] reason R4: 包含来自 unbound-scope…`（dump5）⇒ 风险层判 L2；
  ③ L2 落 `NoGate` 即拒（`gate.go:144`）⇒ 孩子的调用在桥外打回；
  ④ 机制：`subagent_197.go:487` 收口时把结论按 `task.output` 标到**父任务**名下 ⇒ 第一个孩子收口后引擎出现首枚 mark；
  此后任何**未注册** task id 走桥都被 `provenance.go` scopeMarks 的 fail-closed 分支报 unbound-scope ⇒ R4 ⇒ L2。
  修法：harness 的 `AdmitTask` 就是 composition root 的每任务边界（`agent/loop.go:372` 首个工具前交 id、结束时 revoke），
  补 `h.bridge.OpenTask(taskID)`＋revoke 里 `CloseTask(taskID)`＝照 `pointer_185_cli_seam_test.go:346` 先例。
  它不加也不放宽任何断言：`x-noopadmit-after` 变异（布线摘回 no-op）三枚仍全红。
- ⚠ **具名留给编排者的发现（不进本腿写面，本腿未跑 cmd/wisp 面）**：产码侧今天没有任何调用点在任务开始开 scope——
  `cmd/wisp/run.go` 的 `admitTask` 只 `AdmitTextTask`＋revoke 里 `CloseTask`；`OpenTask` 只在 `mark()`（敏感源收口）懒开。
  ⇒ "第二个任务／子代理的第一次工具调用，在别的 scope 已有 mark 时被判 unbound-scope L2"是可在 cmd/wisp 面复量的形状。
  建议落 `A##`/`Q##`；台账不在本腿写面，故只此具名。

## 3. AC#2 —— 守池注释旁 +1 行（零仪器提示）

- `subagent_197_test.go` 注释段尾新增一行（同义提示）：`No instrument watches this prose - only a human reads it (ticket 237
  AC#2): deleting or rewording this paragraph turns nothing red, and the tie from it to the rulers it describes
  (Test222WaitingParentHoldsNoBridgeSlot reads len(bridge.sem), this nail reads the two constants) exists only in a reader's head.`
- 判据：`git diff -- internal/tools/subagent_197_test.go` ＝ **+1／-0**，删除列为零，新增行是 `//` 注释；⛔ 未新写"断言注释文字"的测试。

## 4. 门禁 rc 表 与 diff 形

| 项 | 命令 | rc | 读数 |
|---|---|---|---|
| build | `GOFLAGS= go build ./...` | 0 | — |
| vet | `GOFLAGS= go vet ./internal/tools/` | 0 | — |
| gofumpt | `"$GOPATH/bin/gofumpt.exe" -l <两枚>` | 0 | 零输出 |
| race 口径 | `go test -race -count=60 -run Test222 ./internal/tools/` ×2 | 0／0 | 0 FAIL／0 DATA RACE |
| 全包 | `go test -count=1 -race ./internal/tools/` | 0 | `ok 24.773s` |
| d22scan | `sh scripts/d22scan.sh` | 0 | `d22scan: clean`（含注释与 `_test.go`） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | 0 | VERDICT GREEN；over-budget=58；单位互核同 58 枚 |
| diff 形 | `git diff --stat -- internal/tools` | — | 197：+1/-0；222：+136/-47 |
| 终态 | `git status --porcelain -- internal cmd` | — | 仅本腿两枚 `_test.go`（提交后为空） |

## 5. 盘上日志（大件留盘不入库；入库件＝本目录 4 枚小读数＋本文件）

- `logs/pre-fix-race60-shot1.txt` 770 B ／ `shot2.txt` 383 B（AC#1 改前红）
- `logs/post-fix3-race60-shot1.txt`／`shot2.txt` 各 54 B ／ `fullpkg-race1.txt` 55 B（改后）
- `logs/mut-partial-after.txt` 2,184 B ／ `mut-partial-before.txt` 720 B ／ `mut-full-after.txt` 2,257 B ／ `mut-noopadmit-after.txt` 2,184 B
- `logs/dump-leg2.txt` 31,294 B ／ `dump4-leg2.txt` 57,078 B ／ `dump5-leg2.txt` 68,968 B ／ `dump6-leg2.txt` 19,652 B（归因 instrumentation）
- `logs/leg2-count20-norace.txt` 11,815 B ／ `post-fix2-count40-norace.txt` 30,130 B（r1 自捉的错计形现场）
- `logs/d22scan.txt` 22,948 B ／ `pathlen.txt` 1,457 B（门禁原文）
- 变异件：`mut/x-partial-after_222_test.go`／`x-partial-before_222_test.go`／`x-full-after_222_test.go`／
  `x-noopadmit-after_222_test.go`；overlay 映射 `overlay-*.json`（每形一份）。
