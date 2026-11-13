# 314 卡死根因诊断（只读腿）— `cmd/wisp` 整包跑「卡死」的五问取证

- 腿：314 只读诊断腿（⛔ 非写腿、⛔ 非实现者）
- 时间：2026-11-13 21:2x–21:4x +08
- 派单：`.scratch/wisp/probes/orch/20261113/txt/314-hang-rootcause-dispatch.txt`
- 测量腿的 HEAD：`397cd9e9db1ebb898f629671c408795f48898d6f`（本腿接手时仓 HEAD 已到 `c1ff25028`，但 `397cd9e9d..HEAD` **零代码改动**，只有 docs/scratch —— `git diff --name-only 397cd9e9d..HEAD` 实测）
- 证据根：`.scratch/wisp/probes/314/**`（本腿追加 `logs/repro-run1-subset.log`）
- 门禁：`sh scripts/d22scan.sh` = clean（见文末）
- ⛔ 本腿零代码/测试/脚本改动；⛔ 未跑全量（只跑了 11 枚子集，见 §0）

---

## 0. 本腿新做的那一发测量（因为缺一块证据）

派单说「run2 的 dump 已给具名栈，别重跑全量」。本腿补的是**另一头**：run1 的卡死**没有** dump，只有「最后一行 + CPU 30s 增量 0.0156s + 日志 8m31s 零增长」。为把 run1 的卡死变成具名栈，本腿复跑了 **run1 的前 11 枚**（不是全量）：

- 命令形：`.scratch/wisp/stage-cli/wisp.test.exe -test.count=1 -test.timeout=100s -test.v -test.run '^(TestAC1AlwaysBranch…|…|Test219AC2RevokeRoundTripByteForByte)$'`（11 枚，工作目录 `cmd/wisp`）
- 读数：**11/11 PASS，rc=0，约 25s**（`logs/repro-run1-subset.log`；逐枚 0.00–5.91s，`Test219AC2RevokeRoundTripByteForByte` 3.56s）
- 与 run1 的差：run1 在第 11 枚处日志停在 `scope disposed` 后 8m31s 零增长、CPU 近 0；本腿同一前缀 25s 全过。

⇒ **run1 的卡死不是「前 11 枚确定性地走不到头」**，是一发**竞态**（同前缀可复现不了）。这条读数直接影响 Q4：卡点漂移不是单调累积的共享状态，而是概率性的。

---

## Q1　`panel_pump_test.go` 的 `done`：声明 / 关闭点 / 关闭条件 / 永不关的分支

**声明**：`cmd/wisp/panel_pump_test.go:155`　`done := make(chan struct{})`

**关闭它的 goroutine**：`cmd/wisp/panel_pump_test.go:158` `go func() {` … **`cmd/wisp/panel_pump_test.go:159` `defer close(done)`** … `cmd/wisp/panel_pump_test.go:165` `}()`

**关闭条件**：该 goroutine 的函数体返回，即 `cmd/wisp/panel_pump_test.go:163` 的
`out, err = executeOn145(t, rt, "pump-task", "pump-corr", "p1", "fs.read", …)` 返回。
`executeOn145`（`panel_pump_test.go:110-118`）= `rt.bridge.Execute(ctx, …)`，`ctx` 是 `panel_pump_test.go:112` 的 `context.WithTimeout(context.Background(), 45*time.Second)`。

**等待点**：`cmd/wisp/panel_pump_test.go:177` `<-done` —— **裸接收，无超时**。

### 永不关闭的分支清单（逐支到行）

1. **`bridge.Execute` 在进入 `PendingApproval` 的 select 之前阻塞**
   - `cmd/wisp/bridge.go:425` → `cmd/wisp/bridge.go:540`（L2 分支）→ `internal/agent/approval/gate.go:623`。
   - `gate.go:628` `g.q.push(d)` → `internal/agent/approval/queue.go:159-187`：卡在 `q.mu`（被别的 goroutine 长期持有）。
   - `gate.go:633` `g.q.grantNonce(it)` → `queue.go:359-372`：`mintGrant()`（RNG）或 `q.mu`。
   - `gate.go:657` `g.ui.Prompt(ctx, p)` → `cmd/wisp/run.go:2739` `consoleApprovalUI.Prompt`：`run.go:2764-2770` 往 `u.out`（= `s.stdout`）写卡片；**stdout 若被下游停止消费就永久阻塞**（carrier 的 `tee` 停读即此形）。
2. **`PendingApproval` 的 select 永不命中**（`gate.go:663-703`）
   - `gate.go:665` `case a := <-it.answer`：**没人答**（见 Q2）就永远不命中。
   - `gate.go:692` `case <-deadline`：`deadline` 来自 `gate.go:647` `g.clock.After(g.q.Timeout())`；`g.clock` 非真实时钟（fake/never）时永不命中。
   - `gate.go:697` `case <-ctx.Done()`：`ctx` 无超时就永不命中。
   - `gate.go:675`/`:694` `g.ui.Update` 阻塞。
3. **goroutine panic**：`defer close(done)` 是函数体第一条（`panel_pump_test.go:159`），panic 也会关掉 `done`；但 goroutine 内 panic 无 recover ⇒ **整枚测试二进制崩**，不是静默卡死。
4. **`rt.bridge` 为 nil**：`executeOn145` 里 `rt.bridge.Execute` 空指针 panic ⇒ 同上，崩而非卡。

★**本例的诚实边界**：本例 `ctx` **有** 45s 上限（`panel_pump_test.go:112`），`deadline` 是真实时钟（`approval.New` 未注入 fake clock，`gate.go:149` 起默认 `SystemClock`），`q.Timeout()` = 31s（`panel_pump_test.go:136` 写进 config 的 `confirm_timeout_sec = 31`）。
⇒ 对 **panel_pump 这一枚**来说，`<-done` 是**有界等待（≤31–45s）**，不是无限。run2 的 dump 正好落在 29s（见 Q2/Q4）。

---

## Q2　`panel_pump_test.go:177` 等的那张 L2 卡：谁发、谁入队、谁答、done 谁关

### 链（逐环 文件:行）

| 环 | 位置 | 事实 |
|---|---|---|
| 发卡（入队） | `internal/agent/approval/gate.go:623` `PendingApproval` → `gate.go:628` `g.q.push(d)` → `internal/agent/approval/queue.go:159-187` | 建 `qitem`，分配 `corr`，进 `q.pending`/`q.byID`；`queue.go:162` 队列满（默认 8，`queue.go:130`）即 fail-closed 拒绝 |
| 发卡（上屏） | `gate.go:655` `promptFor` → `gate.go:657` `g.ui.Prompt(ctx,p)` → `cmd/wisp/run.go:2739` `consoleApprovalUI.Prompt` | `run.go:2741` `u.cards++`；`run.go:2749` `u.live.record(p)` 把一次性 grant 记进原生账；`run.go:2777` 打印「卡片编号：」；`run.go:2783` `u.publish()` 推一次面板快照 |
| 谁该答 | **只有 reply 监听器**：`cmd/wisp/run.go:1458-1460` `if s.reply != nil { rt.attachReplyListener(s.reply, s.replyVeto) }` → `cmd/wisp/approval_reply.go:452-475` → `approval_reply.go:547` `runReplyLoop`（`bufio.Scanner` 逐行读操作员输入） | ★**本例 `s.reply == nil`**：`cmd/wisp/run_test.go:89-135` `newRunFixture` 从不设置 `f.reply`（`run_test.go:81` 字段注释逐字：「stays nil - "nobody can answer"」）。⇒ **这一枚用例的卡按设计就没人答** |
| 无人答时的收敛 | `gate.go:692` `case <-deadline` → `gate.go:693` `g.q.expire(it)` → `gate.go:694` `ui.Update(EventDismissed)` | 到期自动**拒绝**（L2 的极性 ≠ L1，`cmd/wisp/bridge.go:545-549` 注释逐字） |
| done 关 | `cmd/wisp/panel_pump_test.go:159` `defer close(done)` | 即 `Execute` 返回后 |

### ★顺序依赖的候选机制：**无据**（本腿逐条查过）

派单问「是否有别的用例留下一张未答复的 L2 卡，或占着同一个单例/队列/监听器，导致本用例的卡进不了可答状态」。本腿的读数是**反的**：

- **队列/门/界面都是每枚 runtime 私建的，不是单例**：`cmd/wisp/run.go:1052-1053` `rt.liveCards = newNativeCards()` / `rt.ui = &consoleApprovalUI{…}`，`run.go:1054` `rt.gate = approval.New(…)` —— 全在 `assembleRuntime` 里逐次新建。`internal/agent/approval/queue.go:77-97` 的 `Queue` 无包级实例；`approval` 包内**没有**包级 `var` 队列（`internal/agent/approval/*.go` 声明面实测）。
- **监听器逐 runtime**：`run.go:1458` 的 `attachReplyListener` 只在 `s.reply != nil` 时挂；`panel_pump` 的 fixture 没挂，`newWidenRun` 那族才挂（`approval_always_201_test.go:90` `w.h.reply = pr`）。两者不共享。
- 所以「别的用例的未答复卡占住本用例的队列」在**结构上不成立**：本用例的卡进不了可答状态，是因为**它自己没接答复源**，与别的用例无关。

**真共享的东西（有据，但不是队列）**：
- `internal/observe.Default`（进程级注册表）：`internal/observe/goroutine.go:262` `Spawn` 累积 `r.live`；泄漏的 owner goroutine 会跨用例存活（见 Q3）。
- `internal/risk` 的包级 `var Logf`（`internal/risk/blacklist.go:54`）：**每次装配都被覆写**（`cmd/wisp/run.go:1050` `risk.Logf = rt.auditf`）。这是货真价实的跨用例全局可变状态，但它是日志函数，不构成阻塞。
- `internal/risk/pathresolver.go:397-400` `anchorCache sync.Map`（有界 128）。
- `internal/llm/adaptertest/mockllm.go:40` `mockllmRegistry = observe.NewRegistry()`（进程级）。

---

## Q3　`lifecycle_43` owner goroutine 活 3 分钟：**后果**，不是原因

**依据（dump + 源码逐行）**：

1. dump 具名栈：`.scratch/wisp/probes/314/logs/run2.log:3147-3153`
   `goroutine 2778 [select, 3 minutes]` → `cmd/wisp/lifecycle_43_windows_test.go:136` → `internal/observe/goroutine.go:340` → `created by observe.(*Registry).Spawn`。
2. 它是谁建的：`cmd/wisp/lifecycle_43_windows_test.go:133`
   `observe.Default.Spawn("lifecycle-43-owner", "test", nil, func(ctx context.Context){ … })` —— 第三个参数 `root = nil`。
3. 为什么没退：`internal/observe/goroutine.go:286-292` —— `root == nil` 时 `ctx = context.Background()`。**`context.Background()` 永不 Done**。而 fixture 里那句 `t.Cleanup(cancel)`（`lifecycle_43_windows_test.go:100-101`，cancel 的是 fixture 自己的 `root`）**cancel 不到这枚 owner**：它没有把 root 传进 `Spawn`。
4. 它现在能干什么的证据：dump 显示它停在 `select`（`lifecycle_43_windows_test.go:136`），既不在 `ctx.Done` 分支也不在 `wake` 分支，**没持任何锁、没在跑 `queue.run()`**。⇒ 它不可能阻塞别人。
5. 正常时候它靠**进程退出**回收；本趟因为整包卡到超时（`run2.log:3123` panic），它才「活过 3 分钟」被 dump 抓到。

⇒ **定性：后果（泄漏症状），不是原因。** 它是 D38 名额内的 owner（`lifecycle-43-owner`），所以 `ClassifyGoroutine` 不报 unknown；但它是**没绑测试生命周期的 fixture goroutine**，属于「进程活着它就活着」的泄漏。要修也是**测试夹具**层面（把 root 传进 `Spawn`，或 `t.Cleanup` 里 cancel 那枚 owner），**与卡死无关**。

⚠ 反向也记一笔：这枚泄漏**证明**了「进程级 observe 注册表会跨用例累积存活 goroutine」这条共享通道真实存在（Q2 列的那条），但它不是阻塞点。

---

## Q4　卡点为什么会漂：**不是一枚测试写坏了，也不是单调累积的共享状态被坐实；已证的是「竞态 + 时钟 cutoff」**

### 已证的三条

1. **同一前缀可复现不了**：run1 的前 11 枚，本腿复跑 **11/11 PASS、25s**（§0）。⇒ 卡死不是「走到第 11 枚必然卡」。
2. **单跑全活**（测量腿 iso1/2/3 3/3、iso219 5s；本腿 11 枚前缀亦活）。
3. **run2 的「卡点」是有界等待被时钟切掉**：
   - `run2.log:3123` `panic: test timed out after 10m0s`；`run2.log:3125` `running tests: TestRunBooksWithASnapshotOfItsLiveQueue (29s)`。
   - 该枚的 L2 卡：队列到期 31s（`panel_pump_test.go:136` + `gate.go:647`），调用方 ctx 45s（`panel_pump_test.go:112`）。
   - 该枚 21:20:12 起（`run2.log:3101-3102`），panic 21:20:39（`run2.log:3235-3236` 起 21:10:35→21:20:39）。⇒ **dump 抓到的是 29s 处，离 31s 到期只差约 2s**。它不是「永不回答」，是「还没到点」。
   - ⇒ **run2 的「卡点」是红鲱鱼**：是整包 10m 预算到期时正好在跑的那一枚，不是死点。

### 那「漂移」的真正机制（有据）

- **整包时长本身在漂**，而整包里有**自起嵌套 `go test` 的变异用例**：`run2.log:241-242` `Test219AC2MutationRevokeWriteMustPersistForBytesToComeBack (39.70s)` —— 它 `exec.Command("go","test","./cmd/wisp","-run","^Test219AC2RevokeRoundTripByteForByte$","-count=1","-timeout","20m","-overlay=…")`（`cmd/wisp/approval_always_219ac2_test.go:379-382`），**重编整包**。这类用例每枚数十秒，且与机器负载强相关。
- 于是「10m 预算砍在哪一枚」随负载/编译缓存漂移 ⇒ run1 被外部 kill 在第 11 枚、run2 被 panic 砍在第 380 枚，**两个「卡点」都不是机制**。
- run1 那一发是**真死锁**（CPU 30s 增量 0.0156s、日志 8m31s 零增长、进程活着无 TCP），但**本腿未能复现**（§0），所以**run1 的确切阻塞栈缺失**，本腿不猜。

### 共享状态：能坐实的只有这些（都不构成阻塞）

| 候选 | 位置 | 是不是「漂移的共享状态」 |
|---|---|---|
| `observe.Default` 存活 goroutine 累积 | `internal/observe/goroutine.go:262,275-282`；实例＝Q3 的 `lifecycle-43-owner` | **是共享、会累积**，但阻塞不了别人 |
| `risk.Logf` 全局覆写 | `internal/risk/blacklist.go:54`；`cmd/wisp/run.go:1050` | 是共享可变，只是日志函数 |
| `risk.anchorCache` | `internal/risk/pathresolver.go:397-400` | 有界 128 |
| `adaptertest.mockllmRegistry` | `internal/llm/adaptertest/mockllm.go:40` | 进程级注册表 |
| 每枚 runtime 的 gate/queue/ui | `cmd/wisp/run.go:1052-1054` | **不是共享**（Q2） |

★结论：**「漂移」的最佳解释＝竞态 + 嵌套变异用例把整包时长抬到 6–20min，使外部时钟（kill / `-test.timeout`）落在哪一枚随负载而变**；不存在一枚「写坏了」的用例（单跑全活），也没有坐实一个单调累积且会阻塞的共享状态。

---

## Q5　最小修法建议（⛔ 只给方案，本腿不动手）

> ⚠ 写面冲突：`cmd/wisp/**` 与 `scripts/portable-tests.sh` 现在都有别的腿在改。本腿只出方案。

**(a) 用例间显式隔离** —— 代价：**中**。需要动 `cmd/wisp/**` 夹具：把 `observe.Default.Spawn(..., nil, …)` 改成传测试 root（如 `lifecycle_43_windows_test.go:133`），并给每个 fixture 的 reply root / retention scope 一个 `t.Cleanup` 顺序契约。它**不解决** run1 那类竞态卡死（本腿没复现出栈），只解决 goroutine 累积。

**(b) 给裸 `<-ch` 加有界超时＋失败即报** —— 代价：**小**。本例 `panel_pump_test.go:177`、`approval_always_219ac2_test.go:242-243`（`<-w.done` / `<-w.reply`）、`approval_always_219ac2_test.go:242` 前后的 `<-done`（`panel_pump_test.go:177`）都是裸接收；改成 `select { case <-done: case <-time.After(60*time.Second): t.Fatalf("…") }`。它把「无限卡」变成「60s 响亮红」，**直接消灭「卡死」这一形态**，且立刻让 run1 那种「无栈可查」变成「有栈可查」（`t.Fatalf` 会打测试名+行）。

**(c) `scripts/portable-tests.sh` cli 分支加终点锚断言** —— 代价：**最小**（脚本面，1 处）。cli carrier（`scripts/wisp-cli-staged-tests.sh:491-495` 那两行 `go test` 包装）跑完后断言 `-test.v` 输出里出现终点锚 `TestZZ315A1`；没出现即 rc≠0 并打印「早停/卡住」。它**不修卡死**，只保证 CI 不会把「早停的四数」当分母（票 314 编排者已立新口径 `A976`）。

### 推荐：**(b) 为主刀，(c) 为立即上线的廉价 tripwire，(a) 排后**

理由（按「最有价值 + 最符合主流 harness」）：
- 主流 harness 的第一戒是**测试不许静默挂起**：有界等待＋失败即报是 Go 测试的标准形状（`go test` 的 `-timeout` 只在前台 panic，且**栈是全包的、定位差**——run2 就是证明）。(b) 把「8m31s 零 CPU 无栈」变成「60s 具名红」，**直接买到可诊断性**，代价最小、落点最准。
- (c) 是**便宜的第一刀**，测量腿的判断本腿同意：它让「卡住」在 CI 上**响亮变红**，且不改任何测试语义、不动阈值。**建议今天就上**（但写面在另一条腿手里，需编排者排队）。
- (a) 是**真清理**但收益最间接（Q3 的泄漏不阻塞），且与 (b) 同写面；**排在 (b)(c) 之后**，别用它当第一刀。

★ 若只能做一件：**先 (c)**（脚本面、零语义风险、当天可落），**再 (b)**（测试面、把裸 `<-ch` 全部加上界）。

---

## 门禁与纪律

- `sh scripts/d22scan.sh`：**clean**（本文件命令的读数见提交信息／票面 Progress log）。
- ⛔ 零代码/测试/脚本改动；⛔ 未跑全量（仅 11 枚子集，§0）；⛔ 未 push。
- 写面仅：本文件（新）、`.scratch/wisp/probes/314/logs/repro-run1-subset.log`（追加副本）、票 314 Progress log。
