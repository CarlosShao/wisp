# AC1ResidentLeg 三枚真红归因（只读腿 · 零写面）

> **尺**：`HEAD=dd8286a6b`（`dev`），现跑 2026-11-12 19:08–19:30 +08。
> **本腿写面**：只有本文件一枚。`*.go` 零改动、断言零改动、AC 框零翻动、台账零改动。
> **载体**：`scripts/wisp-cli-staged-tests.sh`（`bash -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC1ResidentLeg'"`）。
> ⚠ 本机 `bash` 在 `PATH` 上解析到 `C:\Windows\System32\bash.exe`（WSL 转发器，报 `execvpe(/bin/bash) failed: No such file or directory`，rc=1，**一枚测试都没跑**）。改用 `D:\work\soft\Git\bin\bash.exe` 执行同一条命令；脚本本身一字未改。这是一条**环境事实**，不是本腿的结论。
> **临时件（只建不删）**：`D:\tmp\wisp\20261112\resident-ac1\`。

## 0. 一句话

三枚**同一根因**：常驻腿从进程启动到打印 `wisp: resident event loop running` 之间，在**本机**要花 **3.4–5.6 秒**（大头是 `speech.NewDetector` 装载 silero VAD，见 §2.0）；而这三枚断言给出的预算分别是 **10s（够）／1s（不够）／8s（够但撞上另一个更早的窗口）**。前两枚是**预算量错**，第三枚是**两个 defer 的登记顺序被票 43 改过**、而断言仍按老顺序读最后一条记录 —— 它们**共用同一段被测代码（`runResident` 的 boot→loop→defer 链）与同一台慢机**，但**不是同一行产码**。三根因落在同一段链上，**收敛成一件事：`resident_sink_nail_127_windows_test.go` 的三枚钉把"这台机器多快"和"defer 怎么排"当成了产品契约**。

## 1. 三枚逐字的红句（`文件:行`）

载体日志 `.scratch/wisp/stage-cli/wisp-cli-tests.log`（首发）与 `D:\tmp\wisp\20261112\resident-ac1\stage2\wisp-cli-tests.log`（私域 stage 复跑）逐字一致：

| # | 测试名 | 逐字红句 | 耗时（首发／复跑） |
|---|---|---|---|
| 1 | `TestAC1ResidentLegInstallsItsLogListenerOnDisk` | `resident_sink_nail_127_windows_test.go:498: the child never reached the event loop, so the disk read above proves little.` | 5.95s／7.71s |
| 2 | `TestAC1ResidentLegOutlivesItsOwnLogFailure` | `resident_sink_nail_127_windows_test.go:539: the refusal stopped the app: the child never reached the event loop.` | 6.64s／9.91s |
| 3 | `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` | `resident_sink_nail_127_windows_test.go:635: the last record in the sink is "wisp: D38(e) step 10 (exit): ok elapsed=0s", want a shutdown-step record: the shutdown must be logged before the sink closes (records: [...45 条...])` | 8.29s／10.44s |

三枚的 `--- FAIL:` 行（逐字）：
```
--- FAIL: TestAC1ResidentLegInstallsItsLogListenerOnDisk (5.95s)
--- FAIL: TestAC1ResidentLegOutlivesItsOwnLogFailure (6.64s)
--- FAIL: TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink (8.29s)
```
> 第 3 枚的 `:635` 红句全长 **6729 字节**（内含 45 条记录全量打印），上面只抄了判语头部；尾部逐字为 `... msg="wisp: D38(e) step 10 (exit): ok elapsed=0s" dir= step=0}])`。

**同族另有两条 `t.Logf`（不是红句，但读数关键）**：
```
resident_sink_nail_127_windows_test.go:433: resident sink ...\logs: 1 file(s), 12 record(s)
resident_sink_nail_127_windows_test.go:640: RESIDENT LEG RECORDED ON DISK: 5 shutdown record(s), steps [1 2 5 6 7], first="shutdown step skipped (module not present)" last="shutdown step skipped (module not present)"
resident_sink_nail_127_windows_test.go:571: resident leg exit: exited 0
```
⇒ **第 1 枚的磁盘半边全绿**（12 条记录、install 记录在场、位置正确），红只红在"子进程有没有走到 loop"这一条；**第 3 枚的 trail 也在场**（5 条、steps 1/2/5/6/7），红只红在"最后一条是不是 shutdown 记录"。**两枚都不是能力缺失。**

## 2. 逐枚定性

### 2.0 先量一个共同的分母：boot → loop 到底多久

同一枚 `wisp.exe`（`go build -o ... ./cmd/wisp` ＋ `third_party/sherpa-onnx/*.dll` 同目录），`WISP_ENV=test` ＋ `WISP_TEST_DATA_DIR` 私有目录，从 `Start-Process` 到 stdout 出现 `resident event loop running`，现跑：

| 环境 | 轮次 | boot→loop |
|---|---|---|
| 五个 voice 变量全给（与载体同形） | v1／v2 | **5.61s／3.45s**（另两发 4.2s／3.42s） |
| 只给 `WISP_SILERO_MODEL` | 1 | **4.29s** |
| 五个全不给 | n1／n2 | **2.87s／2.86s**（另两发 1.77s／1.88s） |

**钱花在哪**：stderr 时间戳逐行（voice 全给那发）：
```
19:15:28.603  winsec: sealing path resolver installed
19:15:28.628  wisp: persistent log sink installed          ← 磁盘 sink 已开（+25ms）
19:15:29.826  voice: the store carries no VAD model; using WISP_SILERO_MODEL
19:15:31.657  voice: consumer booked on the D38b roster    ← 中间 1.83s
19:15:31.666  lifecycle: WM_QUERYENDSESSION 门已建
```
那 **1.83 秒**落在 `cmd/wisp/resident_voice_windows.go:713` 的 `speech.NewDetector(speech.DetectorConfig{ModelPath: rv.plan.VadModelPath})` 上（它构造在 boot 线程、`rv.scope.Go` 之前，`:766-769` 才 book）。**这是产品行为、不是台面**，而且是**刻意的**（`:748-750` 逐字："a construction inside the loop would drop the first frames of whatever the microphone was already saying"）。
⇒ **分母事实：磁盘 sink 在 ~1.5–2.0s 就绪，event loop 要到 ~2.9–5.6s 才打印。两者之间有 1–3 秒的真空，正是三枚钉各自踩中的地方。**

### 2.1 第 1 枚 `:498` —— **产品行为**（但红句判的是"仪器够不够快"，不是产品坏了）

- **断言**：`if !leg.stdout.has(residentReachedLoop)` ⇒ 子进程 stdout 里必须出现 `resident event loop running`。
- **它断言的是产品行为**：`residentReachedLoop` 的判语在 `:95-119` 逐字写着它是"the sink landmark"的配对物 —— 只有走过 `installLogSink` 与 shutdown defer **之后**的分支才会打印。**删掉 install 块，这条就该红**，这正是票 127 存在的理由。
- **产码责任人**：`cmd/wisp/resident_windows.go:851`（`fmt.Printf("wisp: resident event loop running ...")`）。
- **为什么这次红**：测试自己的时序是
  1. `:424` `pollUntil127(200, ...)` 等磁盘上出现第一条记录（预算 200×50ms = **10s**）；
  2. `:433` 读盘；
  3. `:483` `leg.stop()` —— **无条件 Kill**；
  4. `:497` 才检查 `leg.stdout.has(residentReachedLoop)`。

  第 1 步在**磁盘 sink 一就绪**（实测 1.49–2.04s）就返回，第 3 步立刻 Kill；而 loop 行要到 **2.9–5.6s** 才打印 ⇒ **子进程在第 3 步就被打死，第 4 步永远读不到那一行**。现跑复现（同一 exe、同环境、照抄上面 4 步）：
  ```
  poll_returned_at=1.49s (sink record present)
  loop_line_present_BEFORE_kill=False
  loop_line_present_AFTER_kill=False
  ```
  子进程 stderr 跨度只有 `19:29:07.554 → 19:29:07.853`（**299ms**，12 行），停在 `ball: [hotkey] ...` —— 离 `resident event loop running` 还差整整一个 VAD 构造。
- **定性结论**：**断言的产品主张成立且是真钉（⛔ 不许放宽）**，但**它被写在了"磁盘先于 loop"的假设上**，而本机 VAD 构造把 loop 推到 sink 之后 1–3 秒 ⇒ **红的是"这台机器的时序"**。判据：把 `leg.stop()` 挪到 loop 轮询之后再 Kill，本条即绿（本腿未改，仅陈述判据）。

### 2.2 第 2 枚 `:539` —— **产品行为**（红句判的是"仪器预算 1s 够不够"）

- **断言**：`if !sawLoop` ⇒ `pollUntil127(20, func() bool { return leg.stdout.has(residentReachedLoop) })`，预算 **20×50ms = 1 秒**。
- **它断言的是产品行为**：`residentRefusalPrefix` 已在场（`:533` 的 `sawRefusal` 通过，子进程 stderr 逐字有 `wisp: 持久日志未启用（... logs is not a directory）：安全告警只会到 stderr，不会落盘`），本条要证的是**拒装日志不等于拒启动** —— 这是产品主张，且是这条钉的**全部意义**（`:509-519` 注释）。
- **产码责任人**：`cmd/wisp/resident_windows.go:851`（同一行）；拒装句在 `:178`。
- **为什么这次红**：**1 秒预算对这台机器不够**。现跑照抄 `:529`（先 10s 等 refusal）＋ `:530`（再 1s 等 loop）：
  ```
  voice   : refusal_at=1.23s loop=NOT_WITHIN_1s
  voice   : refusal_at=1.12s loop=NOT_WITHIN_1s
  novoice : refusal_at=1.61s loop=NOT_WITHIN_1s
  novoice : refusal_at=1.41s loop=NOT_WITHIN_1s
  ```
  拒装句自己在 1.1–1.6s 才出现，而 loop 要 2.9–5.6s ⇒ **1s 窗口结构性读不到**。
- **★这枚是 flake，必须按 flake 记**：单独复跑 `^TestAC1ResidentLegOutlivesItsOwnLogFailure$` **5 发得 3 PASS / 2 FAIL**（`-test.v` 那 4 发：`PASS 6.78s`／`FAIL 6.47s`／`PASS 5.78s`／`PASS 6.48s`）。**过的时候是 6.78s／5.78s／6.48s，红的时候是 6.47s** —— 同一量级，边界上翻面。⇒ **⛔ 不判它真红；记"偶发（边界 flake），复跑 PASS 6.78s / 5.78s / 6.48s，复跑 FAIL 6.47s"**。
- **定性结论**：**产品主张成立（拒装不阻启动，逐字有拒装句＋loop 照跑），红因是 1s 预算 vs 本机 3.4–5.6s 的 boot**。判据：把 `:530` 的 `20` 抬到与 `:529` 同量级（200）即绿（本腿未改，仅陈述判据）。

### 2.3 第 3 枚 `:635` —— **产品行为**（★红因是**两个 defer 的登记顺序**，与时长无关）

- **断言**：`last := recs[len(recs)-1]`，要求 `strings.Contains(last.Msg, residentShutdownRecord)`（`residentShutdownRecord = "shutdown step"`，`:131`）。
- **它断言的是产品行为**：`:554-564` 逐字 —— "the two defers in `resident_windows.go` are registered in a specific order so that the D38(e) sequence runs while the sink is still open and the close flushes its records"。这是**排序**主张，是真钉。
- **为什么这次红（关键，与时长无关）**：盘上**最后一条记录是 `wisp: D38(e) step 10 (exit): ok elapsed=0s`**，而不是 `shutdown step ...`。
  - `wisp: D38(e) step N ...` 由 `cmd/wisp/resident_windows.go:246`（`slog.Info(fmt.Sprintf("wisp: D38(e) step %d (%s): %s elapsed=%s", ...))`）打印，**在 defer 尾部把十步全量回放**；
  - `shutdown step ...`（`shutdown step skipped (module not present)`）由 `internal/proc/shutdown.go:130` 打印，**在序列执行当时**；
  - ⇒ 十步回放**必然排在**序列当时那些记录**之后**，所以**文件最后一条一定是 `D38(e) step 10 (exit)`，`shutdown step` 永远不可能在末位**。
  - **这是票 43 引进的**：`5056a365c`（2026-10-31，"票 43 落地腿 43-r1"）在 `resident_windows.go` 的 defer 里**新加了** `:246` 这个逐步回放循环（`git show 5056a365c -- cmd/wisp/resident_windows.go` 的 `@@ -116,14 +126,79 @@` 块里 `+ slog.Info(fmt.Sprintf("wisp: D38(e) step %d ...`）。断言 `:635` 与它的 `t.Logf`（`:640`）都是 **2026-09-22/23** 的（`git blame`：`8663a39967`／`9b5d64d8ba`），**比票 43 早一个月**。
  - ⇒ **断言没跟上产码**：它读"末位"，而票 43 之后"末位"被 D38(e) 回放占满。**时长、负载、环境都不参与这一枚**（现跑 5/5 稳定红：`7.80s`／`6.96s`／`6.79s`／`5.82s`／`5.97s`，voice 与 novoice 都红）。
- **旁证（这条尤其硬）**：`t.Logf` 自己逐字说 `RESIDENT LEG RECORDED ON DISK: 5 shutdown record(s), steps [1 2 5 6 7]` —— **5 条 shutdown 记录确实在盘上**，`len(trail) != 0` 那关（`:629`）**过了**；红的只是"它在不在末位"。
- **定性结论**：**产品行为**，**真红且稳定**，且**是"断言落后于产码"而不是"产品退化"**。判据：`last` 改成"`trail` 的最后一条必须是 `recs` 里最后一条 `shutdown step` 记录"、或改成"`trail` 非空且其步序单调"即绿（本腿未改，仅陈述判据）。

## 3. 台账既有归因

`grep -n "AC1ResidentLeg" docs/reports/pending-and-issues.md` 命中 11 处（4594／6276／6424／9897／12025／12056／12067／14865／18192／20057／23256）。逐条对号：

| 行 | 逐字（节录） | 对本轮三枚的效力 |
|---|---|---|
| 4594 | "第 1 批三发红是既有 flake `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`（`0xc000013a`，同命令四发 3 红 1 绿）" | **与本次不同**：那次红是 `0xc000013a`（进程被 CTRL_BREAK 打死在 flush 之前，exit code 非 0）；**本次 `leg exit: exited 0`、45 条记录全在盘上、红在末位断言**。同一枚钉，两种红，⛔ 不可互相引用。 |
| 6276 | "`go test ./internal/panel/ ./cmd/wisp/` 两包合跑时 `TestAC1ResidentLegBooks...` 报 `0xc000013a` 红，单独跑就过" | **与本次不同**：本次是**逐包单跑**（载体只跑 `./cmd/wisp`），且红因不是 `0xc000013a`。 |
| 6424 | "它只 call 了 `./cmd/wisp/` 这一包、用 CI 同形脚本，仍然撞上那枚假红 ...（`PASS=69 FAIL=1`；终态重跑 `PASS=70 FAIL=0`）" | **与本次不同**：那次自愈；**本次第 3 枚 5/5 稳定红**，不属"编队静不静"那一族。 |
| 9897 | "球那条路径实测 268 毫秒 ... 这段窗口里按 Ctrl+C 会以 0xc000013a 打死进程 ... 它一度因此把票 127 那枚钉改红" | **已知、且是本轮的远因之一**：它记的是"最早可中断时刻"这个窗口；本次第 3 枚**不是**这个窗口（exit 0）。 |
| 12025 | "`--- FAIL` 共 1 枚＝`TestAC1ResidentLegInstallsItsLogListenerOnDisk`（`resident_sink_nail_127_windows_test.go:433`／`:498` 逐字 "the child never reached the event loop"）" | **★已知，同一枚、同一行 `:498`、同一句逐字**。台账当时把它记为"与那 16 枚名册不相交 ⇒ 新增 0 枚"，**没定性**。 |
| 12056 | "`cmd/wisp` 那两枚带载偶发（`TestAC1ResidentLegInstallsItsLogListenerOnDisk`／`TestAC14GoSideEvalPushReachesThePage`）" | **已知**：把它列进"带载偶发"。**本轮把它落到了机制上**（§2.1 的 kill-before-loop），不再只是"带载偶发"四个字。 |
| 12067 | "在册那两枚带载红（`TestAC1ResidentLeg…`／`TestAC14GoSideEvalPush…`）**两发都没红**（grep 零命中），它没追没修＝处置正确" | **已知**：确认它会自己变绿。 |
| 14865 | "`289-v1` 第二发翻出 `TestAC1ResidentLegBooks...`，`resident_sink_nail_127_windows_test.go:565`，失败读数 `:605/:571 exit status 0xc000013a`" | **与本次不同**：`0xc000013a` 族；本次 exit 0、红在 `:635`。 |
| 18192 | "19 枚红逐名全属本机已知噪声族（... `TestAC246*`＋`TestAC228*`＋`TestAC1ResidentLeg*`＝机主实例在跑那族 ...）" | **与本次不同**：那次归因是"机主实例在跑"；**本次运行时 Wisp 实例未跑**（`single instance = false`，逐字在子进程 stdout 里）。 |
| 20057 | "它报的 6 枚红（... `TestAC1ResidentLegInstallsItsLogListenerOnDisk` ...）无一枚属 SLO/settle 族" | **已知**：只做了排除，未定性。 |
| 23256 | "红名 `TestAC1ResidentLegInstallsItsLogListenerOnDisk`（`:498` "never reached the event loop"）—— ★**它说 `HEAD` 上单跑 3/3 也红（本笔 1/3 红）⇒ 环境时序 flake ⇒ 新增红 0** ⇒ ★**我让终裁核这个归因**" | **★已知，且与本次一致**：它已判"环境时序 flake"；**本轮把这个判语坐实到机制上**（§2.1：`leg.stop()` 早于 loop 打印）。**但要注意口径差别**：它说"单跑 3/3 也红"，**本腿单跑 3/3 也红**（`5.72s`／`5.84s`／`11.49s`，三发全红）⇒ **对第 1 枚而言"单跑"不救它**，它不是"编队静不静"那族，是**结构性时序**。 |

**第 2 枚 `TestAC1ResidentLegOutlivesItsOwnLogFailure` 在台账里零命中**（`grep "OutlivesItsOwnLogFailure" docs/reports/pending-and-issues.md` ⇒ 无输出）。⇒ **无既有归因，本文件是它的首份定性**。
**第 3 枚在台账里的每一条命中都是 `0xc000013a` 族或"合跑/带载"族，没有一条对上本次的"末位被 D38(e) 回放占满"** ⇒ **本文件是它这一种红的首份定性**。

## 4. 一句话结论：**同根因还是三根因**

**三根因，落在同一段链上。** 共同面是同一枚被测物（`runResident` 的 boot→loop→defer 链）与同一台慢机；分叉面是三处不同的"仪器假设"：**第 1 枚**假设"磁盘有记录时 loop 已经跑起来"（本机 1.5s vs 3–5.6s ⇒ Kill 抢在 loop 前）；**第 2 枚**假设"1 秒够走到 loop"（本机要 2.9–5.6s，且它是 flake，5 发 3 绿 2 红）；**第 3 枚**假设"末位一定是 `shutdown step`"（票 43 的 `5056a365c` 把十步回放塞进 defer 尾部后，末位必然是 `D38(e) step 10`，**与时长无关、5/5 稳定红**）。⇒ **⛔ 不可用一个修法糊三枚**：前两枚是预算/时序（`leg.stop()` 的位置、`pollUntil127` 的 attempts），第三枚是断言与产码脱节（末位判据要改成 trail 内部判据）。三枚的**产品主张都成立、都是真钉、⛔ 都不许放宽成"存在即可"**。

---

### 附：本腿跑过的尺（可复跑）

| 尺 | 命令 | 读数 |
|---|---|---|
| 载体（首发） | `bash -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC1ResidentLeg'"` | 3 RUN / 3 FAIL（5.95s / 6.64s / 8.29s） |
| 载体（私域 stage 复跑） | 同上，`WISP_CLI_STAGE_DIR=D:\tmp\...\stage2` | 3 RUN / 3 FAIL（7.71s / 9.91s / 10.44s），红句逐字同 |
| 第 2 枚单跑 ×5（voice 全给） | `wisp.test.exe -test.count=1 -test.run '^TestAC1ResidentLegOutlivesItsOwnLogFailure$'` | **3 PASS（6.78s / 5.78s / 6.48s）／2 FAIL（6.47s / 10.20s）** ⇒ flake |
| 第 1 枚单跑 ×3 | 同上，换名 | 3 FAIL（5.72s / 5.84s / 11.49s）⇒ 稳定红 |
| 第 3 枚单跑 ×2（voice）／×2（novoice） | 同上，换名 | 4/4 FAIL（7.80s / 6.96s / 6.79s / 5.82s）⇒ 稳定红 |
| boot→loop 计时 | 手起 `wisp.exe`，轮询 stdout | voice 3.42–5.61s；只 VAD 4.29s；无 voice 1.77–2.87s |
| sink 文件 vs loop 行 | 同上，同时轮询 `logs\*.jsonl` 与 stdout | sink 1.86–2.04s／loop 2.86–5.61s ⇒ **sink 恒早于 loop** |
| 第 1 枚四步复现 | 照抄 `:424/:433/:483/:497` | `poll=1.49s`、`loop_before_kill=False`、`loop_after_kill=False` |
| 第 2 枚预算复现 | 照抄 `:529(200)/:530(20)` | refusal 1.12–1.61s、loop **1s 内读不到**（4/4） |
