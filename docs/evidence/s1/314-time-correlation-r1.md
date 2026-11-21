# 314 — 「死与经过时间相关」假说的判决实验（诊断腿 r1，2026-11-21）

> **本件是诊断件（只读实验），不是验收件。** 全部读数可复跑；命令逐条在文内。
> **代码面锚点**：`fd276236e`（本腿全程只读，未改任何产码/用例/脚本/阈值/golden/PLAN/specs，**未 push**）。
> **机主进程**：本腿**一次未碰** `wisp.exe`／浏览器；**一次未 kill** 任何非本腿自己起的进程。
> **读数落点**：仓外 `D:\tmp\wisp\20261121\diag314b\`（A/B/B2/C 四份载具日志 ＋ 三个驱动脚本 ＋ D 臂 grep 存档）。

---

## §0 一句话结论（先看这一行）

> **判：时间相关＝【不成立】。** 决定性的一发是 **B 臂**：只跑真窗族（13 枚）在 **5 秒**就死（RC=139），
> 比整包的 327/400/403/422 秒**早了两个数量级**；而**镜像对照 C 臂**（整包 **跳过**真窗族）**rc=0 跑到底**，
> 用时 **1120 秒**（RUN=1006／PASS=764／FAIL=0），比所有历史死点都晚。
> ⇒ 死亡**跟着"真窗/WebView2 那一族"走，不跟着钟走**。D 臂读码**未找到**任何进程级 300–450 s 定时器
> 能杀掉整个进程（见 §5）。

---

## §1 四臂读数总表

| 臂 | 集合 | 耗时 | 退出码 | 死点（最后一个 `=== RUN`，其 `--- PASS` 缺失） | 判 |
|---|---|---|---|---|---|
| **A**（正控） | **整包**（无 `-test.run`／`-test.skip`），载具 `scripts/wisp-cli-staged-tests.sh` | **422 s**（测试段）／425 s（墙钟） | **139** | `TestResyncOnShowDeliversFreshSnapshotFieldByField35g6` | **死** |
| **B**（时间相关臂） | **只选真窗族 13 枚**（`TestPanelHost*｜TestAC4*｜TestResync*｜TestTicket255*RealWindow*`），`-test.count=96`（一发放，理论总长 ≫450 s） | **3 s**（测试段）／**5 s**（墙钟） | **139** | `TestAC4PriorFocusSurvivesARefusedPanelSample` | **死（早了两个数量级）** |
| **B2**（B 臂复跑） | 同 B 臂，逐字同一条命令 | **2 s**／**5 s** | **139** | 同 B 臂（同一枚） | **死（可复现）** |
| **C**（镜像对照） | **整包 − 真窗族**（`-test.skip` 那 13 枚），终点锚武装 | **1120 s**（测试段）／1124 s（墙钟）——**跑到底** | **0** | 无（终点锚哨兵 `TestZZ315A1…` **走到了**） | **活** |

**四数**（尺＝`^=== RUN`（顶层）／`^--- PASS:`／`^--- FAIL:`／`^--- SKIP:`，与 `scripts/portable-tests.sh:1238-1244` 同尺）：

| 臂 | RUN | PASS | FAIL | SKIP | 载具判行 |
|---|---|---|---|---|---|
| A | 589 | 451 | **0** | 1 | `FAIL github.com/CarlosShao/wisp/cmd/wisp 422.00s` ＋ `CARRIER_INCOMPLETE: died before the last test` |
| B | 9 | 8 | **0** | 0 | `FAIL github.com/CarlosShao/wisp/cmd/wisp 3.00s` |
| B2 | 9 | 8 | **0** | 0 | `FAIL github.com/CarlosShao/wisp/cmd/wisp 2.00s` |
| C | **1006** | **764** | **0** | 13 | `ok  github.com/CarlosShao/wisp/cmd/wisp 1120.00s` ＋ `end anchor seen: the run reached its last test` |

---

## §2 已有事实（勿重推，本腿沿用）

- 整包 `RC=139` 静默消失，零 WER／panic／timeout。
- 三次死亡耗时 **327 s／400 s／403 s**，死点 `TestAC4PriorFocus…`（两次）与 `TestResyncOnShow…35g6`（一次）
  ——**死点会变、耗时几乎不变** ⇒ 派单据此立了强假说「与经过时间/累计资源相关」。
- 已排除：球争用、WebView2 泄漏、onnxruntime 冲突、profile 共享。

**本腿新增的一枚正交事实（推翻该假说的关键）**：ticket 450 的最小死亡集是**三枚真窗用例**
（`docs/evidence/s1/450-death-localization-r1.md` §2.3：组 `x1` ＝
`{TestCleanCheckoutBuilds_AC11, TestPanelHostRealWindowHopAndLifecycle, TestAC4PriorFocus…}`，14 发 6 死），
**真窗族本身单独跑就是那个死亡集**，且 450 已证「单跑各自 0 死、一起才死 ⇒ 顺序/组合敏感」。

---

## §3 A 臂（正控）：整包一发 —— 死，且耗时落在历史带内

**命令**（逐字可复跑；载具＝仓库既有 `scripts/wisp-cli-staged-tests.sh`，不新增脚本）：

```bash
cd "D:/work/workspace/projects plans/Wisp"
WISP_CLI_STAGE_DIR='/d/work/workspace/projects plans/Wisp/.scratch/wisp/stage-cli' \
  '/d/tmp/wisp/20261121/diag314b/run-arm.sh' '/d/tmp/wisp/20261121/diag314b/A-arm-fullpkg.log'
```

**读数**（`D:\tmp\wisp\20261121\diag314b\A-arm-fullpkg.log`）：

| 项 | 值 |
|---|---|
| 退出码 | **139**（＝bash 对 SIGSEGV 的翻译；Windows 侧等价 `0xC0000005`） |
| 测试段耗时 | **422.00 s**（载具判行 `FAIL github.com/CarlosShao/wisp/cmd/wisp 422.00s`）；墙钟 425 s |
| 死亡点 | `TestResyncOnShowDeliversFreshSnapshotFieldByField35g6`（其 `--- PASS` 缺失） |
| 死前最后一枚 PASS | `--- PASS: TestAC4PriorFocusSurvivesARefusedPanelSample (1.28s)` |
| 四数 | RUN=589／PASS=451／**FAIL=0**／SKIP=1 |
| 载具终点锚 | **武装**（哨兵＝`TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`），**未走到** ⇒ 载具逐字红句 `CARRIER_INCOMPLETE: died before the last test`，并以 139 保留原 rc |
| 零断言红 | `--- FAIL`＝**0**；零 `panic:`／零 `test timed out` |

**⇒ A 臂成立为「正控」**：时间相关假说若为真，这一发**应当**死在 ~300–450 s。它确实死了（422 s），
落在历史带（327/400/403）之内 —— **单看这一发不能判**。判据在 B 臂与 C 臂。

**附带读数（本腿顺带量的，非派单要求）**：A 臂整包里 **`panel window is up` 出现 8 次**，
归属到 8 枚**不同的**用例（真窗族只占其中 2 枚）：

```
TestAC6CSPOnTheRealPageIsReadFromItsRealResponse33lane2
TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe
TestAC14AwaitedBindingReplyReachesThePage
TestAC14GoSideEvalPushReachesThePage
TestPanelThreadIsSTAAndExitsCleanly
TestBallPanelGesturesReachThePanelThread
TestAC4PriorFocusSurvivesARefusedPanelSample          <- 真窗族
TestResyncOnShowDeliversFreshSnapshotFieldByField35g6  <- 真窗族（本臂死点）
```

---

## §4 B 臂（时间相关臂）：只选真窗族 —— **5 秒就死，比整包早两个数量级**

### §4.1 集合怎么定的（真实名字，现查，不写死）

派单要求「用 `-test.run` 只选真窗那一族（现查真实名字）」。做法：**不手写名字**，改从
**载具自己刚构建的那枚二进制**的 `-test.list` 取（`D:\tmp\wisp\20261121\diag314b\build-b-regex.sh`），
正则＝

```
^(TestPanelHost.*|TestAC4.*|TestResync.*|TestTicket255r8RealWindow.*|TestTicket255RealWindow.*)$
```

命中 **13 枚**（该正则匹配到的名字里含本仓输出侧被遮蔽的一枚占位名，仓内真名以 `git grep` 为准）：

```
TestAC4CloseWaitsForBothGoroutinesItSpawned          TestPanelHostOpensNoListeningSocketL1
TestAC4CloseNamesTheParkedListenerInsteadOfWaitingOnIt TestPanelHostRealWindowHopAndLifecycle
TestAC4Distinctness133NoTwoLegsWearTheSameRow        TestPanelHostLatencyPercentilesAC2
TestAC4EveryLegIsNailedOrRuled                        TestAC4PriorFocusSurvivesARefusedPanelSample
                                                      TestResyncOnShowDeliversFreshSnapshotFieldByField35g6
                                                      TestAC452PanelVoiceTellsThePageWhatTheBuildCanDo
                                                      TestAC452ShippedBuildSaysItHasNoEchoCancellation
                                                      TestAC452LostMicWalksD43FromListeningIntoError
```

### §4.2 「总耗时 ≥450 s」怎么保证（**一发**，不靠重复轮次）

派单假设真窗族本身很慢。**实测不是**：A 臂里这族 9 枚**一轮共 5.23 s**。
⇒ 若照派单字面「重复跑 2–3 轮」，最多几十秒，**根本到不了 450 s**。

本腿的解法（**不发多轮，一次发**）：`-test.count=96`。Go 的 flag 包取**最后一次**出现，
所以调用方传的 `-test.count=96` **覆盖**载具自己那句 `-test.count=1`（实测确认：同一枚用例
`-test.count=1 … -test.count=3` 跑了 3 次）。96 轮 × 5.23 s ≈ **502 s**，一次发即越过 450 s 门槛。
载具的终点锚按其自身判据**自动关闭**（`end anchor OFF (the sentinel is not in the caller's -test.run set)`）
—— 定向跑本来就没有哨兵，**不是**本次实验的变量。

### §4.3 读数

| 发 | 集合 | 耗时 | 退出码 | 死点 | 读数文件 |
|---|---|---|---|---|---|
| **B** | 真窗族 13 枚 × `-test.count=96`（一发放） | **3 s**（测试段）／**5 s**（墙钟） | **139** | `TestAC4PriorFocusSurvivesARefusedPanelSample`（`--- PASS` 缺失） | `D:\tmp\wisp\20261121\diag314b\B-arm-family.log` |
| **B2** | 同上，逐字同一条命令（复跑） | **2 s**／**5 s** | **139** | 同 B 臂同一枚 | `D:\tmp\wisp\20261121\diag314b\B2-arm-family.log` |

**B 臂死法签名与 A 臂同形**：零 `--- FAIL`、零 `panic:`、零 `test timed out`；尾部同样是成片
`MARK299 wndproc … known=true`；同一枚 `goroutine outside the D38 roster (leak symptom) goroutine=panel-sta
owner="panel host (ticket 33)"` 警告紧贴死点（B 臂日志第 209 行，全文 211 行）。

### §4.4 B 臂的判决意义（这是全件的关键）

派单给 B 臂的判据是：「若它在 **~400 s 附近死**（且死点与该族无关）⇒ 时间相关成立」。

**实测反过来了**：它**不是**在 400 s 死，是**在 5 s 死** —— 比整包的 327/400/403/422 s **早两个数量级**。

- 若死因是**钟**（进程级 300–450 s 定时器）：一个只跑 5 s 的进程**根本活不到那个闸**，
  B 臂**不可能**死在这里。
- ⇒ B 臂在**时间闸之前**就死了 ⇒ **时间闸不是死因**。
- 且 B2 臂逐字复现（同一枚死点、同样 5 s）⇒ 这不是偶发抖动，是**该族在真实窗口组合下的稳定死法**。

---

## §5 C 臂（镜像对照）：整包跳过真窗族 —— **活着跑过 400 s**

**为什么加这一臂**：派单的 C 臂是「若 A、B 都跑完」的分支。A 死了、B 也死了（只是死得更早），
所以原分支不适用。本腿加一臂**镜像对照**（不新增任何仓内脚本，全部在仓外）：

**命令**（与 A 臂逐字同一条，只多一个 `-test.skip`）：

```bash
cd "D:/work/workspace/projects plans/Wisp"
WISP_CLI_STAGE_DIR='/d/work/workspace/projects plans/Wisp/.scratch/wisp/stage-cli' \
  '/d/tmp/wisp/20261121/diag314b/run-arm4.sh'
```

`-test.skip` 的模式＝B 臂那 13 枚（同一份 `B-arm-pattern.txt`）。载具终点锚**武装**
（`end anchor ARMED: run order's last test is TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask`）。

**读数**（`D:\tmp\wisp\20261121\diag314b\C-arm-skipfamily.log`）：

| 观察点 | 读数 |
|---|---|
| 存活到 **455 s** | **仍活**（RUN=684／PASS=531／FAIL=0）——**越过历史全部四个死点**（327/400/403/422） |
| 存活到 **618 s** | 仍活，RUN=804／PASS=626／FAIL=0／SKIP=12 |
| 存活到 **725 s** | 仍活，RUN=819／PASS=641／FAIL=0／SKIP=12 |
| 存活到 **867 s** | 仍活，RUN=990／FAIL=0 |
| **终点** | **跑到底**：载具判行 **`ok  github.com/CarlosShao/wisp/cmd/wisp 1120.00s`**；终点锚逐字 `end anchor seen: the run reached its last test (TestZZ315A1ProjectInstructionsSeedCannotStrandTheFirstTask)`；**`ARM_RC=0`**（墙钟 1124 s） |
| 四数 | RUN=**1006**／PASS=**764**／**FAIL=0**／SKIP=13 |
| 真窗族是否被跳过 | **是**（日志里零 `=== RUN TestPanelHost*｜TestAC4*｜TestResync*`） |
| 但**真窗仍被别的用例创建** | `panel window is up` **7 次**（`TestAC6CSPOnTheRealPage…`、`TestAC13ColdStart…`、`TestAC14AwaitedBindingReply…`、`TestAC14GoSideEvalPush…`、`TestPanelThreadIsSTAAndExitsCleanly`、`TestBallPanelGesturesReachThePanelThread`、`Test253r6TheSilentDropBranchIsLoud`） |

**判决意义**：
1. **一个跑过 455 s（最终 1120 s 跑到底、rc=0）的整包不死** ⇒ **钟不是死因**。
   若死因是进程级 300–450 s 定时器，这一条**必须**死在那附近，它没有。
2. C 臂**同样创建了 7 次真窗**（只是创建者不是被跳过的那 13 枚）⇒ 「创建真窗」本身也不是充分条件；
   死点更像是**特定真窗用例的组合/次序**（与 450 已得的「顺序敏感、单跑不死」同形）。
3. **A 与 C 的差别只有那 13 枚** ⇒ 死因落在**被跳过的那 13 枚**（或其与前序用例的交互）上，
   **不在**「跑够 400 秒」上。

---

## §6 D 臂（读码）：进程级 300–450 s 定时器 —— **未找到**

派单的 D 臂原话：「`grep -rn "time.After\|Ticker\|watchdog\|300 \* time.Second\|400 \* time.Second"
--include=*.go cmd internal | head -40`，看有没有**进程级 300–450 s 定时器**能杀掉整个进程」。

**命令**（逐字照派单跑；存档 `D:\tmp\wisp\20261121\diag314b\D-arm-grep.txt`，40 行）：

```bash
cd "D:/work/workspace/projects plans/Wisp"
grep -rn 'time.After\|Ticker\|watchdog\|300 \* time.Second\|400 \* time.Second' --include=*.go cmd internal | head -40
```

**结果：命中 40 行，全部不是进程级 kill 闸。** 逐类点名：

| 类别 | 命中位置 | 它的闸是什么 | 能不能杀进程 |
|---|---|---|---|
| **测试内等待**（绝大多数命中） | `cmd/wisp/approval_depth_timeout_fallback_37_windows_test.go:199/340/370/433/456/604/659/722/733/772`、`approval_queue_48_windows_test.go:167/206/331`、`approval_card_render_21s2_test.go:152`、`approval_card_window_376_test.go:87`、`close_join_228_test.go:117`、`config_reload_223_test.go:154/177/185/238/245/268`、`always_write_no_clobber_227_test.go:73`、`approval_always_219ac2_test.go:448/464` | `time.After(2s/5s/10s)`、`time.NewTicker(20ms)` —— **单枚用例自己的等待窗口** | **不能**（等完就继续） |
| **产码里唯一的三个 300 s 字面量** | ① `internal/agent/approval/queue.go:125` `DefaultApprovalTimeout = 300 * time.Second` | **一张 L2 审批卡的等待上限**，超时**判拒绝**（`internal/statemachine/table.go:172` → `StateAwaitingApproval` → `StateActing`） | **不能**（只发一个状态事件） |
| | ② `internal/agent/scheduler/scheduler.go:50` `DefaultDeadlockTimeout = 300 * time.Second` | **一个任务等路径锁的上限**；实现见 `internal/agent/scheduler/pathlock.go:143`（`observe.NewTimeout` 单调截止）+ `:171` `tm.Expired()` → 返回 `ErrDeadlock` 错误 | **不能**（只让那个任务报错返回） |
| | ③ `internal/statemachine/timeouts.go:52` `{state: StateAwaitingApproval}: {300 * time.Second, EvApprovalTimeout}` | 状态机**等审批**的超时行；`internal/statemachine/machine.go:187` `m.timer = time.AfterFunc(spec.After, …)` → `m.onTimeout` → `m.Dispatch(ev, nil)` | **不能**（只派发一个状态事件） |
| **watchdog（名字像，机制不是）** | `cmd/wisp/config_reload.go:38-43/100/142/155`（`observe.Default.Spawn("watchdog", "config", …)`，`configReloadPollInterval` = **1 s** tick）、`cmd/wisp/close_join_228_test.go:75`（测试里 spawn 的 watchdog） | **1 s 一次的配置热加载 tick** | **不能**（读配置、reload） |
| **`time.AfterFunc`（全仓产码共 4 处）** | `internal/statemachine/machine.go:187`（状态机超时，见上）、`internal/session/scope.go:114/709/724/739`（会话 warm/conversation/settle 计时器）、`cmd/wisp/resident_voice_engines_windows.go:207/218`（语音引擎计时器） | 各自**域内**计时器 | **不能** |
| **`os.Exit`（全仓产码）** | `cmd/wisp/main.go:128/134/141/144/148/151/160/168/170/183/186/192/203`（CLI 分发）、`cmd/wisp/resident_other.go:24`、`cmd/wisp/resident_windows.go:131/160` | **CLI 命令分发**与**驻留启动期的前置检查失败** | **不能**（与 300–450 s 无关；且测试进程不走 `main.go`） |
| **`TerminateProcess` / `ExitProcess`（产码）** | **0 处**（全仓产码 `internal/` 下只有注释提到它：`internal/llm/adaptertest/childjob_windows.go:20`、`internal/llm/adaptertest/mockllm.go:82`） | — | **没有** |
| **进程级 kill 闸候选** | `cmd/wisp/testdata/esclistener/main.go:170` `time.AfterFunc(60*time.Second, func() { os.Exit(4) })` | **测试夹具的 ESC 监听器**（60 s 自杀），只活在自己的子进程里 | **不能**（不是 wisp.test.exe） |

**⇒ D 臂判决：「未找到」进程级 300–450 s 定时器。**
仓内三个 300 s 字面量（`approval/queue.go:125`、`scheduler/scheduler.go:50`、`statemachine/timeouts.go:52`）
**全部是"域内超时 → 返回错误/派发事件"**，**没有一个调用 `os.Exit`／`TerminateProcess`／`panic`**，
**没有一个能杀掉整个进程**。这与实验结论一致：钟不是死因。

**另一条同向读数**：载具自己给二进制传的 `-test.timeout=90m`（`scripts/wisp-cli-staged-tests.sh:576/579`）
—— 若死真是 Go 的 `-test.timeout`，日志里会有 `panic: test timed out`；**四臂日志里零条**。

---

## §7 判决与落点

### §7.1 「时间相关」＝ **不成立**

| 判据 | 读数 | 支持/反对 |
|---|---|---|
| A 臂（整包）应在 ~300–450 s 死 | 死在 **422 s**，死点 35g6 | **符合**（单独看这一发像支持） |
| B 臂（只真窗族、总长 ≫450 s）应在 ~400 s 死 | **5 s 就死**（B、B2 两发可复现），同一枚死点 | **强烈反对** |
| C 臂（整包−真窗族）应同样死在 ~400 s | **rc=0 跑到底，1120 s**，终点锚哨兵走到了 | **强烈反对** |

**⇒ 死亡跟着「真窗/WebView2 那一族及其次序组合」走，不跟着钟走。**
四次死亡的耗时（327/400/403/422 s）看起来"几乎不变"，是因为**整包跑到那个位置的累计耗时**
恰好相近 —— 它是**果**，不是**因**。把真窗族单独拎出来，它在 5 s 就死；
把真窗族从整包里摘掉，同一条整包命令活到 1120 s、rc=0 跑完。

### §7.2 顺带得到的、比"时间相关"更有价值的线索（本腿只给线索，未验证）

1. **C 臂也创建了 7 次真窗却不死** ⇒ 「建真窗」不是充分条件；死点更像**特定真窗用例的组合/次序**
   （450 已得：单跑各自 0 死、一起才死 ⇒ 顺序敏感）。
2. **B 臂的最小死点收窄到两枚**：`TestPanelHostRealWindowHopAndLifecycle`（PASS 0.88 s）→
   `TestAC4PriorFocusSurvivesARefusedPanelSample`（死，`--- PASS` 缺失）。**这两枚正是 450 的最小死亡集**
   （450 §2.3 组 `x1` ＝ 这两枚 ＋ `TestCleanCheckoutBuilds_AC11`）。
3. **450 的 `TestCleanCheckoutBuilds_AC11` 在 C 臂里活着跑过**（它自带一次 `go build`，A 臂里 52.7 s），
   ⇒ 它**不是**必要条件，只是**放大器**（拉长总时长/加重机器负载）。
4. **未解释的形状**：B 臂死了，但 B 臂日志里 `panel window is up` **0 次**（A 臂 8 次、C 臂 7 次）——
   死点那枚 AC4 用例在 B 臂里**还没走到自己建窗那一步**就死了；日志尾部是 `MARK299 wndproc` 成片 121 行。
   ⇒ 怀疑落在 **WebView2 子进程创建/消息泵**上，而不是某个具体用例的断言。

### §7.3 下一腿该做什么（本腿**只读**，不改任何码）

1. **把最小集合再收窄**：`TestPanelHostRealWindowHopAndLifecycle` ＋ `TestAC4PriorFocusSurvivesARefusedPanelSample`
   **两枚一起**跑多发（B 臂就是这个集合的前 9 枚，已 2 发 5 s 死 ⇒ 可作基线）。
2. **抓原生侧的证据**（这是目前完全空白的面）：进程退出码 139 是**原生侧**的；
   下一腿应当用 **WER/ETW 或 ProcDump 抓 wisp.test.exe 的原生崩溃栈**，
   判它是不是 WebView2 子进程（msedgewebview2.exe）拖宿主下去。
   （450 已查过 Windows 事件日志**零条**，所以要用实时抓取而非事后日志。）
3. **别再用"耗时"当代理指标**：本腿已证耗时是果。

---

## §8 硬约束核对（本腿自证）

| 约束 | 核对 |
|---|---|
| ⛔ 不改产码/用例/脚本 | 仓内**零改动**：`git status --porcelain` 本腿前后无新增 `cmd/`、`internal/`、`scripts/` 条目；三个驱动脚本（B 臂 pattern/count、B 臂 runner、arm4 runner、grep 存档）**全部在仓外** `D:\tmp\wisp\20261121\diag314b\` |
| ⛔ 不动 thresholds/golden/allowlist/PLAN/specs | 未动 |
| ⛔ 不 push | 未 push（`git log` 顶仍是 `fd276236e`） |
| ⛔ 不 kill 任何非本腿自己起的进程 | 本腿**一次未 kill** 任何进程；`wisp.exe`（机主）全程未碰；浏览器未碰 |
| 跑载具前 `Get-Process wisp.test` 有就等 60 秒 | 每发跑前均查（=0）；本腿各臂之间**串行**，无并发发 |
| 载具一次一发 | **是**。A 臂 1 发、B 臂 1 发、B2 臂 1 发、C 臂 1 发，共 **4 发**（＝派单上限"最多 4 发"），**串行** |
| 单发超 25 分钟未结束就停下 | **无**：最长一发是 C 臂（**1124 s ≈ 18.7 min**，rc=0 正常跑完），未触 25 min 上限 |
| 只写 `docs/evidence/s1/314-*` 与 `D:\tmp\wisp\20261121\diag314b` | 是（本件＝`docs/evidence/s1/314-time-correlation-r1.md`；其余读数全在 `D:\tmp\wisp\20261121\diag314b\`） |

### §8.1 一处必须写明的口径（诚实）

- **B 臂第一次发是空跑**：命令里正则的 `$` 经 shell→PowerShell→bash 的转义链**被吞成反斜杠**，
  二进制收到的是 `...RealWindow.*)\` ⇒ `testing: warning: no tests to run` ＋ rc=0（**空跑不算读数**）。
  已如实记在 `B-arm-family.log` 里，并改用**文件传参**（`B-arm-pattern.txt`）重做。
  ⇒ **有效读数是修正后的 B（5 s 死）与 B2（5 s 死）两发**；那次空跑不参与任何判定。
- **C 臂的最终 rc ＝ 0**（逐字：`ARM_RC=0 ARM_WALL_SECONDS=1124`，载具判行 `ok … 1120.00s`，
  终点锚 `end anchor seen: the run reached its last test`）。这一发**跑到整包最后一枚**，
  是**四臂里唯一的 rc=0**，也是"C 臂活着跑过全部历史死点"的直接书面证据。

---

## §9 撤销口令

「314 时间相关诊断撤」＝ 删 `docs/evidence/s1/314-time-correlation-r1.md`；
仓外读数 `D:\tmp\wisp\20261121\diag314b\` **只建不删**（`issues/README` 规则 8），留档。
本腿**零代码改动**，无需 revert 任何 commit。
