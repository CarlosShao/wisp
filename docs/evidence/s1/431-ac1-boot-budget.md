# 票 431 — 两枚 `TestAC1ResidentLeg*` 红的不是同一个根：一枚稳定红在"等错了时刻"，一枚偶发在"预算写死"

> **尺**：起手 `HEAD=235292472`（`dev`），现跑 `2026-11-12 21:16–22:1x +08`。
> **写面**：`cmd/wisp/resident_sink_nail_127_windows_test.go` 一枚 ＋ 本文件一枚。**产码零字节**。
> **载体**：`D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC1ResidentLeg…$'"`。
> ⚠ 本机 PATH 上的 `bash` 是 WSL 转发器（`C:\Windows\System32\bash.exe`，`execvpe(/bin/bash) failed`），**一枚测试都不跑**。本文所有读数**全部**经 Git bash 取。
> **临时件（只建不删）**：`D:\tmp\wisp\20261112\431-ac1budget\`（`install-run{1..5}.log` / `outlive-run{1..5}.log` / `final-install-run{1..3}.log` / `final-outlive-run{1..3}.log` / `final-mutant-*.log` / `overlay*.json` / `mutant/` / `overlay/` / `probe/` / `bin/`）。

## 0. 一句话（票面假设被推翻的那一条）

票面把两枚写成"①稳定红 ②偶发"，并暗示①红在 **10s 窗口不够**（本机启动 3.4–5.6s）。
**实测不是。** ①的 10s 轮询 **611ms 就成功了**——预算一次都没绑住；它红在**后面那句** `leg.stdout.has(residentReachedLoop)`，因为用例在 ~0.6s 就把子进程 `stop()` 了，而那条启动完成行要到 **2.8–4.3s** 才出现。⇒ ①是"**在事实还不存在的时候去读它**"，②才是"**预算写死成 1 秒**"。两枚**不同根**，但**同一个正确修法**：等事实，不等时间。

## 1. AC#0 只读取证（⛔ 本格零写）

### 1.1 两枚各复跑 5 次（改前，逐次耗时）

命令（每发独立跑，`-test.count=1`）：

```
D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC1ResidentLegInstallsItsLogListenerOnDisk$'"
D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC1ResidentLegOutlivesItsOwnLogFailure$'"
```

| 发 | ① `…InstallsItsLogListenerOnDisk` | ② `…OutlivesItsOwnLogFailure` |
|---|---|---|
| 1 | **FAIL** 12.54s | **FAIL** 7.07s |
| 2 | **FAIL** 8.76s | **FAIL** 7.59s |
| 3 | **FAIL** 8.02s | **FAIL** 7.57s |
| 4 | **FAIL** 8.34s | **FAIL** 10.60s |
| 5 | **FAIL** 6.13s | **FAIL** 8.81s |
| 定性 | **稳定红 5/5** | **稳定红 5/5**（本会话读数；票面 round 88 记的是 3 PASS / 2 FAIL ⇒ 见 §1.5） |

**⚠ 与票面的两处出入，如实报：**
1. **②本会话是 5/5 全红，不是票面写的 3 PASS / 2 FAIL。** 票面那份 5 发是 round 88 整包跑的读数；本会话单枚跑是 5/5 红。两者都成立且不矛盾：**②的红是"1 秒 vs 2.8–4.3 秒"的边界赌**，整包跑时机器负载不同、子进程冷热不同，会偶尔踩进 1 秒内（round 88 那 3 次 PASS 就是这么来的）。所以②的正确定性是**偶发（flaky）**，不是稳定红；本会话 5/5 只是这次全踩在红侧。
2. **①的红因不是 10s 窗口。** 见 §1.2。

### 1.2 ① 的逐字红句，与它真正红在哪

5 发逐字同形（`resident_sink_nail_127_windows_test.go:505`）：

```
    resident_sink_nail_127_windows_test.go:440: resident sink C:\Users\swq\AppData\Local\Temp\TestAC1ResidentLegInstallsItsLogListenerOnDisk72963709\002\logs: 1 file(s), 12 record(s)
    resident_sink_nail_127_windows_test.go:505: the child never reached the event loop, so the disk read above proves little.
--- FAIL: TestAC1ResidentLegInstallsItsLogListenerOnDisk (6.13s)
```

**注意 `:440` 那行是 `t.Logf`，不是错误**：它证明**盘上已经有 12 条记录**，即 10s 轮询那一关**过了**。红的是 `:505`。

**为什么 `:505` 必红——用 `go test -overlay` 打时间戳实测**（拷贝放**仓外** `D:\tmp\…\431-ac1budget\overlay\`，`overlay.json` 映射回去）：

```
resident_sink_nail_127_windows_test.go:431: PROBE431 poll1 result=true after 611.7732ms
resident_sink_nail_127_windows_test.go:490: PROBE431 pre-stop: loop_line=false sink_read_at=2026-11-12 21:40:28.1783072 +0800
resident_sink_nail_127_windows_test.go:506: PROBE431 post-stop: stdout bytes=2375, has loop line=false
```

⇒ 盘上事实 **611ms** 就到了；用例随即 `leg.stop()`（把子进程杀了），而启动完成行要 **2.8–4.3s**。**子进程在事实能出现前约 2.7 秒就被杀了**，`:505` 读的当然永远是 false。

**反证（同一 overlay 里只把"等事实"挪到 `stop()` 之前，断言一字不改）**：

```
resident_sink_nail_127_windows_test.go:491: PROBE431 waited for loop line: true
--- PASS: TestAC1ResidentLegInstallsItsLogListenerOnDisk (8.20s)
```

⇒ ①的病灶是**读取时刻**，不是窗口大小。把 10s 改成 60s、600s 都不会让它变绿（`stop()` 之后 stdout 再也不会长出新行）。

### 1.3 `pollUntil127` 与两枚各自的等待条件（真正在等什么）

`cmd/wisp/resident_sink_nail_127_windows_test.go:322-330`：

```go
// pollUntil127 runs done on a fixed attempt budget (the shape ticket 63's
// pollCommandLine uses) and reports whether it ever said yes.
func pollUntil127(attempts int, done func() bool) bool {
	for i := 0; i < attempts; i++ {
		if done() {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return done()
}
```

⇒ `attempts × 50ms`。`200` = **10s**，`20` = **1s**。**它本来就接受条件函数**——所以"改成等事实"不需要新机制，只需要换掉那个乘出来的时长。

| 用例 | 旧等待 | 真正在等的事实 | 该事实**最早**何时可观测（实测） |
|---|---|---|---|
| ① `…InstallsItsLogListenerOnDisk` | `pollUntil127(200, …)` = 10s，等 `sinkHasARecord(sinkDir)` | 监听器把 **install 记录写进盘**（`wisp: persistent log sink installed`，`cmd/wisp/logsink.go:200`） | **0.61s**（3 发：615.6ms / 610.7ms / 611.3ms） |
| ① 的**第二**个事实（真红点） | 无等待——`leg.stop()` 之后直接 `has(residentReachedLoop)` | 子进程打印 **`resident event loop running`**（`cmd/wisp/resident_windows.go:851`） | **2.36–2.59s**（3 发：2.5946s / 2.3550s / 2.3857s） |
| ② `…OutlivesItsOwnLogFailure` | `pollUntil127(200, …)` = 10s，等 `residentRefusalPrefix` | 日志目录装不上时那句**具名拒绝**（sink 的第一件事） | **100.5ms**（3 发：100.51 / 100.91 / 101.10ms） |
| ② 的**第二**个事实（真红点） | `pollUntil127(20, …)` = **1s**，等 `residentReachedLoop` | 同上：`resident event loop running` | **2.84–2.99s**（3 发：2.9884s / 2.8366s / 2.8854s） |

### 1.4 启动完成行首次可观测的时刻：10 发独立实测

用仓外探针（`D:\tmp\…\431-ac1budget\probe\`，自己起 `wisp.exe` 无参、逐行打 `time.Since(start)`）：

| 发 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 |
|---|---|---|---|---|---|---|---|---|---|---|
| `resident event loop running` | 2.80s | 3.31s | 3.49s | 3.28s | 3.32s | 3.18s | 3.86s | 3.71s | 4.11s | 3.35s |

再 6 发**把 `logs` 挡住**（即②的场景）：**2.77 / 3.93 / 3.73 / 3.46 / 3.12 / 3.77s**。
⇒ 合计 16 发：**2.77–4.32s**，中位数约 **3.4s**。

**阶段成本（同一探针逐行打戳，4 发）**：

```
t=0.07s  wisp: persistent log sink installed        ← ①的第一个事实
t=1.24–1.81s  voice: the store carries no VAD model; using WISP_SILERO_MODEL
t=3.19–4.14s  voice: consumer booked on the D38b roster
t=3.19–4.14s  wisp: resident event loop running     ← 两枚真正要等的那个事实
```

⇒ **`speech.NewDetector`（`cmd/wisp/resident_voice_windows.go:713`）落在 1.24–1.81s 到 3.19–4.14s 之间，实测 2.4–2.6s**，是这段的主导成本。**旁证**：把语音夹具全部撤掉（`WISP_SILERO_MODEL` 等不设）后，同一台机器启动到 1.75 / 1.84 / 1.78s ⇒ **约 1.5–2.5s 的差就是 VAD 构造**。

⚠ 票面写的"`speech.NewDetector` 一枚就 1.83s"：本机实测是 **2.4–2.6s**（阶段差），票面那个数偏小。不影响结论（更慢只会让 1s 预算更不可能够）。

### 1.5 ①与②是否同一个根：**不是**

| | ① | ② |
|---|---|---|
| 红句 | `:505 the child never reached the event loop, so the disk read above proves little.` | `:546 the refusal stopped the app: the child never reached the event loop.` |
| 第一个事实的等待 | 10s，实际 **611ms** 就够（**预算从未绑住**） | 10s，实际 **100ms** 就够（预算从未绑住） |
| 第二个事实的等待 | **没有等待**——`stop()` 之后才读，而事实在 `stop()` 之后 2.2–3.7s 才可能出现 | **1s**，事实要 **2.84–2.99s** |
| 病灶 | **读取时刻错**（把子进程杀在事实出现之前） | **预算写死**（1s 对 2.8–4.3s 的赌） |
| 定性 | 稳定红（5/5） | 偶发（round 88: 3 PASS / 2 FAIL；本会话 5/5 全踩红侧） |

**一句话结论**：①红在 `cmd/wisp/resident_sink_nail_127_windows_test.go:505`（`leg.stop()` 早于事实），②红在同文件 `:546`（`pollUntil127(20,…)` = 1s）。**不同根，同一修法。**

## 2. AC#1 选形：**等事实**（不是加大秒数、不是算系数），理由

### 2.1 三个候选与判据

| 形 | 内容 | 判据 |
|---|---|---|
| 甲 **加大秒数** | `20 → 120`（6s）、`200 → 1200` | ⛔ 仍是魔数。①**根本不会因此变绿**（§1.2 反证：事实在 `stop()` 之后永不到来）；②只是把偶发变少见，换台更慢的机器照红 |
| 丙 **算系数** | 预算 = 实测启动成本 × 安全系数，注释写明依据 | 比甲诚实，但**①照样不变绿**，而且**真故障时仍然要多等那个系数倍** |
| 乙 **等事实**（选中） | 等那条 `resident event loop running` 逐字行 / 等盘上 install 记录；**时长只作上限**，事实一到立刻返回 | ✅ 主流 harness 做法：等**事实**，⛔ 等**时间**。机器慢只是等得久，不会假红；事实永不到来时由上限报红，且**子进程一死立刻报**（见 §2.3） |

**为什么乙是"最有价值"而不是"动作最大"**（AGENTS.md 铁律 2）：
- **甲/丙买不到 ① 的绿**——它们的射程只有"窗口太短"，而①的病不在窗口。选甲就是选一个**不解决目标**的小动作。
- 乙把**同一个机制**用在两枚上，且**不新增魔数语义**：`pollUntil127` 本来就收条件函数，缺的只是"把上限和事实分开"。
- 乙让**真故障更快报**：子进程一死立即返回（§2.3），而甲/丙只会让真故障多等。
- 代价写清：乙多一个 `bootFact127` 结果类型 + 一个 helper（约 90 行含注释），换来"①真绿 + ②真绿 + 突变控制有牙"。**复杂买到了真价值**。

### 2.2 改了什么（断言一字未放宽）

`cmd/wisp/resident_sink_nail_127_windows_test.go`：

1. 新增 `bootFactCeiling127 = 60 * time.Second`（**上限**，非等待；注释写明 16 发实测依据：最慢 4.32s，60s ≈ 14×）。
2. 新增 `waitForBootFact127(leg, budget, fact) bootFact127`：用 `observe.NewTimeout`（D42#9 单调钟，⛔ 墙钟差），**事实一到立刻返回**；两个出口——事实到了 / 子进程死了（`leg.reap(0)`）/ 上限用尽。
3. ①：盘上记录改 `waitForBootFact127`；**并把"等启动完成行"挪到 `leg.stop()` 之前**（这是①真绿的原因）。
4. ②：两个等待都改 `waitForBootFact127`，**共用同一个 `bootCeiling`**（`Remaining()` 传递，⛔ 2×60s）。
5. 失败消息分开写"上限用尽"与"子进程中途死了"，并带上实测耗时。

**⛔ 没有改的**：`residentReachedLoop` 的字面量、`sinkHasARecord` 的判据、`residentRefusalPrefix` / `residentRefusalPromise`、任何 `t.Skip`、任何阈值/golden/`thresholds.go`。要等的**事实**逐字不变。

### 2.3 为什么加"子进程死了立刻报"这一出口

票面明确点名加大秒数的代价之一是"**真故障时多等**"。若只把上限拉到 60s，一个**启动即崩**的产品缺陷就要坐满 60s 才报——那正是票面拒绝的代价。
所以 `waitForBootFact127` 每轮先 `leg.reap(0)`（非阻塞、幂等，不抢 `leg.stop` 的 Wait 结果）；子进程已退且事实仍不成立 ⇒ **立刻**返回 `exitedEarly`，失败消息写"exit code N after Xs"。**上限只由"活着但永不到达"这一种失效模式消耗。**

## 3. AC#1 突变控制（必做，真跑）

**突变体**：把产码里那条"启动完成"标记**删掉**——`cmd/wisp/resident_windows.go:851` 的
`fmt.Printf("wisp: resident event loop running …")` 整句删除，只留 `reason = rt.RunEventLoop()`。
**做法**：拷贝到**仓外** `D:\tmp\…\431-ac1budget\mutant\resident_windows.go`，用 `go test -overlay` 映射回去（`overlay-mutant.json`），**仓库内产码零字节改动**。突变体先验 `go build` rc=0（证明是"删了标记仍能编译"，不是"删坏了"）。

| 用例 | 突变后读数 | 判定 |
|---|---|---|
| ① `…InstallsItsLogListenerOnDisk` | `:541 resident listener booked itself on disk after 609.3864ms`（**盘上事实仍到**）<br>`:633 the child never reached the event loop within 1m0s, so the disk read above proves little.`<br>`--- FAIL (66.69s)` | ✅ **超时红**，不是空转绿 |
| ② `…OutlivesItsOwnLogFailure` | `:689 the refusal reached the console after 151.3452ms`（**拒绝仍到**）<br>`:705 the refusal stopped the app: the child never reached the event loop within 1m0s.`<br>`--- FAIL (66.90s)` | ✅ **超时红** |

⇒ 两枚等的都是**真事实**：标记一删，它们就红在上限上，而**不相关的那个事实（盘上记录 / 拒绝句）照常到达**——这正证明它们各自钉的是自己那句话，不是"随便等到点什么就绿"。

## 4. AC#1 改后绿（每枚 3 发）

| 发 | ① | ② |
|---|---|---|
| 1 | **PASS** 10.65s | **PASS** 10.02s |
| 2 | **PASS** 10.12s | **PASS** 9.32s |
| 3 | **PASS** 10.32s | **PASS** 9.67s |

实测事实耗时（同一发日志里 `t.Logf` 出来的）：
- ①：盘上记录 **615.6 / 610.7 / 611.3ms**；启动完成行 **2.5946 / 2.3550 / 2.3857s**。
- ②：拒绝句 **100.51 / 100.91 / 101.10ms**；启动完成行 **2.9884 / 2.8366 / 2.8854s**。

**同族三枚一起跑**（`-test.run '^TestAC1ResidentLeg'`，39.7s）：

```
--- PASS: TestAC1ResidentLegInstallsItsLogListenerOnDisk (11.33s)
--- PASS: TestAC1ResidentLegOutlivesItsOwnLogFailure (9.84s)
--- PASS: TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink (10.82s)
```

⇒ 票 425 修好的第三枚**仍然绿**（本票未碰它的断言，只共用了新的等待 helper 之外的零处）。

## 5. AC#2 越界检查

`git diff --name-only`（逐行，见 §7）。本票**只**改了 `cmd/wisp/resident_sink_nail_127_windows_test.go` ＋ 本文件 ＋ 票面 Progress log。
`docs/PLAN.md` / `docs/specs/**` / `internal/observe/thresholds.go` / golden / `tools/d22scan/allowlist.txt` / `frontend/**` / `design/**` —— **零命中**。
`build/mut_rvw.go` **仍在**（未删）。

## 6. AC#3 门禁四读数

```
$ go build ./cmd/wisp/          → rc=0（无输出）
$ go vet ./cmd/wisp/            → rc=0（无输出）
$ gofumpt -l cmd internal       → 零输出（rc=0）
$ scripts/d22scan.sh            → rc=0
  runtests.sh: OK - packages=[./...] top-level: PASS=45 FAIL=0 SKIP=0, === RUN=93
  d22scan: examined 480 production Go files under internal/ and cmd/
  d22scan: clean - no D22 ban violations
```

⚠ `gofumpt -l .`（全仓）会吐**大量 `.scratch/wisp/probes/**` 历史探针件**——那是既有残留（`build/probe61r2/main.go` 等），与票 431 无关；按前几腿的射程用 `gofumpt -l cmd internal`，零输出。

## 7. AC#4 与票 425 / 127 的关系

- 本票只处理**预算/等待方式**这两枚（①的读取时刻 ＋ ②的 1s 预算）。
- ⛔ 未翻票 127／425 任何 AC 框。
- ⛔ 不宣称 `TestAC1*` 全族绿：本票只证明 `TestAC1ResidentLeg` 三枚同跑绿（§4），全族未跑。
- 与票 425 的差别：425 的红因是**断言落后于产码**（"末位"结构上不可满足）；本票两枚红因是**等待形状**。**⛔ 混为一谈。**

## 8. 给下一位的两条注意事项

1. **`pollUntil127` 仍有 28 处调用**（`lifecycle_live_43`、`resident_task_source_246` 等），本票**未动**它们。它们多数等的是"某句台词出现"，同样有"预算 vs 事实"的味道，但**不在本票射程**；若将来再红，先按 §1.3 的表格量出"事实最早何时可观测"，再决定。
2. **①的教训比②更贵**：一个"读得太早"的断言会把**产品正常**报成红，而且**加大预算永远不会修好它**。凡"先 `stop()` 再读子进程输出"的用例，都要问一句"我读的那个事实，在 `stop()` 之前出现过吗"。
