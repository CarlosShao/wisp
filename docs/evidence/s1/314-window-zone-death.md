# 票 314 — `cmd/wisp` 窗口区整进程死亡：**测量腿读数表**（2026-11-13）

- **腿**：314 测量/定性腿（⛔ 非写腿；零产码改动）
- **时刻**：2026-11-13 20:50–21:2x +08
- **★锚点**：派单写「当前 HEAD（`02939b8a2`）」，**本机现读 HEAD＝`397cd9e9db1ebb898f629671c408795f48898d6f`**，
  `git merge-base --is-ancestor 02939b8a2 HEAD`＝真、`git rev-list --count 02939b8a2..HEAD`＝**1**
  ⇒ 派单那个 sha 是本机 HEAD 的**父提交**（差 1 枚 `A1853`）。本表所有读数取自 **`397cd9e9d`**，⛔ 不是 `02939b8a2`。
- **跑法**：`scripts/wisp-cli-staged-tests.sh`（carrier；`CGO_ENABLED=1`；PATH 含 `D:\work\base\go\bin`、`D:\work\soft\Git\bin`、`E:\work\base\msys64\mingw64\bin`）
- **原始日志根**：`D:\tmp\wisp\20261113\314\`（只建不删）；副本 `.scratch/wisp/probes/314/logs/`
- **门禁**：`sh scripts/d22scan.sh` clean（见 §6）

---

## 0. 桌面/会话状态（候选 (a) 的代理变量，每发前后各测一次）

判据＝`OpenInputDesktop` 的桌面名（`Default`＝未锁屏；`Winlogon`＝锁屏）＋ session 1 的 `LogonUI` 进程数＋前台窗口。
⛔ 不猜：这台机器 session 2/3/4 是**断开的 RDP 会话**、各自持有 1 枚 `LogonUI`（pid 24096/27932/50412），
**session 1（console）零枚** —— 所以「有 LogonUI 在跑」不等于「本机桌面锁着」，必须按 session 分。

| 时刻 | input_desktop | logonui(session1) | 前台窗口 | CPU load | cores |
|---|---|---|---|---|---|
| 20:56:12（run1 前） | `Default` | 0 | `44296:Weixin` | 100 | 12 |
| 21:07:16（iso1 前） | `Default` | 0 | `44296:Weixin` | 67 | 12 |
| 21:07:46（iso2 前） | `Default` | 0 | `44296:Weixin` | 100 | 12 |
| 21:08:12（iso3 前） | `Default` | 0 | `44296:Weixin` | 97 | 12 |
| 21:09:40（run2 前） | `Default` | 0 | `44296:Weixin` | 71 | 12 |

⇒ **本轮每一发都在「未锁屏＋有前台窗口（微信）＋机器繁忙（CPU 67–100%）」的桌面上跑**。
⇒ 候选 (a) 的机制（**锁屏/空闲桌面**）**本轮不在场**，因此本轮**无法用它解释**任何一发读数；也**无法证伪**它（没造出锁屏条件）。

---

## 1. 第 1 步：逐发读数（退出码／耗时／四数／最后一行／是否 139）

四数口径逐字＝`scripts/portable-tests.sh:1049-1053`：`^=== RUN` / `^--- PASS` / `^--- FAIL` / `^--- SKIP`（顶层、行首无空白）。

| 发 | 形 | 退出码 | 耗时 | 四数（RUN/PASS/FAIL/SKIP） | 最后一行 | 139? |
|---|---|---|---|---|---|---|
| **run1** | carrier 全量（无 `-test.run`） | **无**（进程未退出，被我 kill） | 8m31s 后被杀 | **11 / 10 / 0 / 0**（只跑了 11 枚） | `scope disposed scope=retention-job-host …`（20:57:05） | **否**（不是死亡，是**卡死**） |
| **iso1** | carrier `-test.run '^(5 枚死亡点)$'` | **0** | 9s | **5 / 5 / 0 / 0** | `ok github.com/CarlosShao/wisp/cmd/wisp 6.00s` | 否 |
| **iso2** | 同上复跑 | **0** | 18s | **5 / 5 / 0 / 0** | `ok … 7.00s` | 否 |
| **iso3** | 同上复跑 | **0** | 13s | **5 / 5 / 0 / 0** | `ok … 8.00s` | 否 |
| **iso219** | carrier `-test.run '^Test219AC2RevokeRoundTripByteForByte$'` | **0** | 5s | **1 / 1 / 0 / 0** | `ok …` | 否 |
| **run2** | carrier 全量 + `-test.timeout=10m` | **2**（撞 test 超时 panic） | **604s** | **380 / 295 / 2 / 0** | `panic: test timed out after 10m0s`（最后 `=== RUN`＝`TestRunBooksWithASnapshotOfItsLiveQueue`，无 verdict） | 否 |

- run2 的两枚具名红（**断言红，非死**）：`TestAC2SealNoticeLandsInTheRunLegLogFile` 0.20s、`TestAC4FocusReturnToPriorWindowGap33r5` 5.21s。
- run2 **终点锚 `TestZZ315A1` 命中 0** ⇒ 按 `A976` 口径，**该发四数不作覆盖面分母**。
- run2 全日志 `panic:` 只 1 处＝那行 `test timed out`；**`signal`／`exit status`／`fatal error` 零命中**。

### 1.1 run1 的「卡死」读数（形态与票面**不同**）

- 卡死点：`=== RUN   Test219AC2RevokeRoundTripByteForByte`（第 11 枚 `=== RUN`；`cmd/wisp/approval_always_219ac2_test.go:230`）。
- **进程活着、但不推进**：`wisp.test.exe` pid 33580 一直在；**CPU 增量 30 秒＝0.0156 s**、8 秒＝0.00、20 秒＝0.00；
  **日志 8 分 31 秒零增长**（45,121 B 恒定，最后写 20:57:05）；无 TCP 连接（`Get-NetTCPConnection -OwningProcess` 空）。
- **⛔ 不是票面死法**：票面签名是「整进程**消失**、无 verdict」；这一发是「整进程**还在**、零 CPU、无 verdict」。
- **⛔ 也不是环境问题**：同一枚测试单跑（iso219）**5 秒 PASS**；2026-11-11 整包日志里同一枚 **2.67s PASS**（`D:\tmp\wisp\20261111\314rerun\fullpkg2.log`）。
  ⇒ **顺序依赖**：单跑活、整包卡。

### 1.2 iso1–iso3 的读数（死点 5 枚单跑，3 发）

5 枚＝`TestAC13BringUpSurvivesAReusedThreadQuit`、`TestAC13BringUpRefusesAThreadWithAQueuedClose`、
`TestAC14AwaitedBindingReplyReachesThePage`、`TestAC14GoSideEvalPushReachesThePage`、`TestPanelThreadIsSTAAndExitsCleanly`。

- **3/3 全 PASS、rc=0**；逐枚耗时 0.97–1.76s。
- 每发都真的建了窗：`panel thread up and down: hwnd 0x…, exits observed, shows=1`、`AC#13 reused-thread … created=true`。
- ⇒ **在本轮桌面状态下，死亡点 5 枚单跑不死**（n=3，m=0）。

### 1.3 run2 的卡死点：**具名死锁**（goroutine dump 逐字）

run2 撞满 10 分钟 test 超时，Go 打出全量 goroutine 栈。**卡点逐字**：

```
panic: test timed out after 10m0s
	running tests:
		TestRunBooksWithASnapshotOfItsLiveQueue (29s)

goroutine 3361 [chan receive]:
github.com/CarlosShao/wisp/cmd/wisp.TestRunBooksWithASnapshotOfItsLiveQueue.func1(...)
	D:/work/workspace/projects plans/Wisp/cmd/wisp/panel_pump_test.go:177 +0x453
github.com/CarlosShao/wisp/cmd/wisp.runTextTask(...)
	D:/work/workspace/projects plans/Wisp/cmd/wisp/run.go:284 +0x415
```

- **`panel_pump_test.go:177` 逐字＝`<-done`**：这枚测试在 `rtHook` 里另起一个 goroutine 去跑
  `executeOn145(... "fs.read" ...)`（`:163`），主 goroutine 在 `:177` 等它回来。
  那枚 `fs.read` 是**越出 allowlist 的读**、判据是 L2 ⇒ **进审批队列、等一张卡片被回答**；
  卡片起来了（`:169` 的 `rt.ui.shown()==0` 那道闸**过了**），但**回答永远没来** ⇒ `executeOn145` 不返回 ⇒ `<-done` 永久阻塞。
- **★这枚测试自己就在窗口区**：它断言 `rt.ui.shown()`、`rt.publishPanelSnapshot()`、`rt.windowCount()`——
  是**面板窗口 + 审批队列**的路子，⛔ 不是纯逻辑测试。
- **同一张 dump 里还有一枚 3 分钟前起的泄漏 goroutine**（`goroutine 2778 [select, 3 minutes]`）：
  `newSuspendFixture43.func3`（`lifecycle_43_windows_test.go:136`，`observe.Default.Spawn("lifecycle-43-owner", …)`）
  ⇒ **前一枚测试起的 owner goroutine 没有被 cancel 掉**，与 `A1746` 的「累积状态」同族。
- **⇒ run2 的卡死是「窗口区 + 审批回执丢失」的具名死锁，⛔ 不是票面那种「整进程消失」。**

---

## 2. 第 2 步：`rc=139` 与票面「整进程死亡」是不是同一件事

### 2.1 本机真实看到的数字（逐字）

| 来源 | 数字（逐字） | 出处 |
|---|---|---|
| 本轮 run1 | **无退出码**（进程未退出；`job_kill` 后才消失） | 本表 §1.1 |
| 本轮 iso1/iso2/iso3/iso219 | **`rc=0`** | 本表 §1 |
| 本轮 run2 | **`rc=2`**（`panic: test timed out after 10m0s`） | 本表 §1.3 |
| 仓内产品码注释（HEAD `397cd9e9d`） | **`rc=139`** | `cmd/wisp/panel_host_windows.go:581-586` |
| 仓内证据件 | **`SIGSEGV，rc=139`** | `docs/evidence/s1/33-panel-host-c27-r7.md:99` |
| 仓内台账（`A1746` 引用的 native 码） | **`-1073741819` = `0xC0000005`** | `docs/reports/pending-and-issues.md`（A1746 族）＋ `314-ac0-direct-run.md:60` |

### 2.2 判据（⛔ 不是理论，是仓内逐字材料＋算术）

1. **产品码自己把 `rc=139` 等同于「整枚进程倒下」**——`panel_host_windows.go:581-586` 逐字：
   > 「the alternative measured here is a half-initialised controller, which on this box can take the **whole process down** (rc=139 reading in the evidence file's ①)」
   ⇒ 在**同一块面板宿主代码**里，`rc=139` 就是「整进程死」的读数，⛔ 不是别的东西。
2. **证据件 `33-panel-host-c27-r7.md:99` 逐字**：「五种 mode 挤在一发进程里跑：第 2 发直接 **SIGSEGV，rc=139** … ⇒ 一次 panic 之后留下的半初始化 controller 能让**整枚测试二进制**没掉，不是"这枚用例红"」。
3. **算术**：bash 的 `139 = 128 + 11 = 128 + SIGSEGV`。Windows 上进程因未处理异常 `0xC0000005`（`STATUS_ACCESS_VIOLATION`）死亡时，
   MSYS2/Git-Bash 把它翻译成 `SIGSEGV` ⇒ `$?` 读 **139**。⇒ **`0xC0000005` 与 `rc=139` 是同一件事的两种写法**（一个是 native 异常码，一个是 shell 退出码）。
4. **与票面死法签名的关系**：票面记「**无 Go panic 输出**、无 `--- FAIL`、无 `exit status` 行——整进程消失」。
   `0xC0000005` 若发生在 **Go 侧**，Go runtime 会打印 `fatal error: fault`／`unexpected fault address`；
   若发生在 **原生侧**（WebView2／onnxruntime 的 C 代码），Go 的 vectored handler 不接，进程直接被 NT 杀掉 ⇒ **零 Go 输出**。
   ⇒ 票面那个「安静消失」的形状与 `0xC0000005` **一致**。

### 2.3 判定

**★同一件事（同一类现象），但有一个必须说清的边界。**

- **是同一类**：`rc=139` ＝ 未处理访问违例（`0xC0000005`）⇒ 整枚测试二进制倒下、零 Go 输出。
  这与票面「整进程死亡（安静、无 verdict、无 panic）」是**同一形态**；新情报「沙盒树 CLI 全量偶发 `rc=139`」与本票主诉**相容**。
- **边界（诚实）**：**本轮我没有亲眼看到 139**。本轮两发全量**都不是死亡**（run1 卡死、run2 撞 test 超时 panic `rc=2`），
   因此**没有「死亡时的退出码」可抄**。
   票面 2026-10-13 那 5 发死亡当时也**没有留下退出码**（原始日志根 `orch311-q-20261013\logs` 现读 MISSING，`314-diagnosis-v1.md:9` 已记）。
  ⇒ 我能证明的是「**139 ＝ 整进程死」这个等式在仓内有具名出处**；⛔ 不能证明「票面那 5 发当时的退出码就是 139」。
- **★本轮新增的一条反向读数**：run2 证明了这个包**也会以 `rc=2`（Go test 超时 panic）收场**——
  即「整包不完成」**不只有「无声消失」一种形态**，还有「卡到超时、响亮 panic」。
  ⇒ 票面那个「无 `exit status` 行、无 panic」的签名**只是其中一种**，⛔ 不是唯一。

---

## 3. 第 3 步：三个候选的可判别实验

### (a) 桌面/会话状态 —— **本轮不可验（既不能排除、也不能坐实）**

- **本轮读数**：§0 那张表——每一发都在 `input_desktop=Default`（**未锁屏**）、session 1 零 `LogonUI`、前台是微信、CPU 67–100%。
- **代理变量**：锁屏判据（`OpenInputDesktop` 名字）、session 1 的 `LogonUI` 计数、前台窗口、CPU 负载。已逐发记录。
- **结论**：**候选 (a) 的机制（锁屏/空闲桌面）本轮根本不在场**，所以本轮没有任何一发能用它解释、也没有任何一发能证伪它。
- **⛔ 不需要机主配合也做不到**：把这台机器弄成「锁屏」是**机主动作**（Win+L），腿⛔ 不能替他锁。
- **需要机主做什么（一句话）**：**「离开机器前按一次 Win+L 锁屏，然后说一声，我复跑同一发全量。」**
  （锁屏期间 session 1 的 `LogonUI` 会 +1、`input_desktop` 变 `Winlogon`——这两个数就是判据。）

### (b) 宿主/前端负载 —— **部分判别：单跑不死、整包仍失败 ⇒ 需要「包内上下文」**

- **读数 1（隔离）**：死点 5 枚单跑 **3/3 全 PASS、rc=0**（§1.2）。⇒ **「这 2 枚测试本身有毒」被削弱**。
- **读数 2（整包）**：run1 全量**没跑到窗口区就卡死**在第 11 枚（`Test219AC2…`，审批流，非窗口区）；
  run2 全量跑到第 380 枚，**卡在 `TestRunBooksWithASnapshotOfItsLiveQueue`**（`panel_pump_test.go:177`，**窗口区 + 审批队列**），
  撞满 10 分钟 test 超时；历史整包（2026-11-11 `314-ac3`）跑到 418 枚后**死在 `panel_run_145` 族**。
- **读数 3（顺序依赖）**：`Test219AC2RevokeRoundTripByteForByte` 单跑 5s PASS、整包卡死 8.5 分钟。
- **读数 4（goroutine dump）**：run2 的卡点是**具名死锁**（`<-done` 等一张永不回答的 L2 卡片），
  **不是**「负载把进程压垮」的那种随机时刻死。
- **结论**：**「死亡/卡死」不是某枚测试的属性，是「整包跑起来之后的累积状态」的属性**。
  这**与候选 (b)（包内负载/累积）同向**，但**⛔ 不足以坐实**——因为「累积状态」也可能是 (a) 桌面会话或 (c) 随机的表现。
  ★另有一条**与 (b) 相反方向**的读数：run2 的卡点是一枚**确定性的通道死锁**（`<-done`），
  而「负载」通常给出**随机时刻**的死/卡；两发全量卡在**两枚不同的测试**上（219 审批流 vs panel_pump 窗口区），
  ⇒ 更像**多枚测试各自带一个「等回执」的闸**，谁先等到谁卡住。
- **⛔ 未验的部分**：新前端页面字节（`frontend/**`）在真窗区的负载**没法在本轮单独隔离**（要改测试才能关掉页面加载，⛔ 写面外）。

### (c) 纯随机 —— **⛔ 不能排除，也⛔ 不能坐实（统计不足）**

| 发类 | n | 死（整进程消失） | 卡（不完成） | 活/完成 |
|---|---|---|---|---|
| 死点 5 枚单跑（本轮） | 3 | 0 | 0 | 3 |
| 全量（本轮） | 2（run1/run2） | **0** | **2** | 0 |
| 全量（历史，2026-10-13 同机） | 5 | 5 | 0 | 0 |
| 全量（历史，2026-10-13 同日早些时候） | 4 | 0 | 0 | 4 |
| 全量（历史，2026-10-15 CI） | 5 | 2 | 0 | 3 |
| 全量（历史，2026-11-11 本机 cli 形） | 1 | 1 | 0 | 0 |
| 全量（历史，2026-11-12 本机） | 1 | 0 | 1 | 0（40 分钟仍在跑，被 kill） |

- **读数**：本轮 **n=2 发全量、死 0 发、卡 2 发**；死点单跑 **n=3、死 0 发**。
- **判断**：**⛔ 不能排除随机**——本轮全量只有 2 发，且两发都是**卡**不是**死**，
  「卡」与「死」在票面里是两种形态（虽然都叫「整包不完成」）。
  票面自己给过算术（5 死 4 活若 p=0.5 是 1/512），但那是**跨天混合样本**，本轮⛔ 不能沿用。
- **要到定性线**：需要**同一台机器、同一桌面状态、连续 ≥10 发全量**，每发抄四数＋退出码＋终点锚＋卡点栈。
  **谁提供**：再跑 N 发（本轮时间不够，见 §5）。

---

## 4. 第 4 步：裁定建议（一行）

> **定性尚不清，且本轮把「死法」扩成了两种形态**：本轮**没有复现票面的「整进程死亡」**——
> 两发全量都是**卡住**（run1 卡在 `Test219AC2…` 审批流；run2 卡在 `TestRunBooksWithASnapshotOfItsLiveQueue`
> `/panel_pump_test.go:177` 的 `<-done`，撞满 10m test 超时 `rc=2`），死点 5 枚单跑 **3/3 全活**。
> **`rc=139` ＝未处理访问违例整进程死** 在仓内有具名出处（产品码注释＋`33-panel-host-c27-r7.md:99`），
> 但**本轮⛔ 没亲眼见到 139**，也没拿到「死亡那一发」的退出码。
> **还缺的读数**：① **同一桌面状态下连续 ≥10 发全量**的退出码＋四数＋终点锚＋**卡点/死点栈**
> （判 (c)，并把「卡」与「死」分开计数）；② **死亡真发生的那一发**的退出码原值（判 `rc=139` 与票面死亡是否同一件）；
> ③ **锁屏桌面**下 1–2 发全量（判 (a)）；④ 全量跑到**终点锚 `TestZZ315A1*`** 的读数（判「四数能不能当分母」）。
> **谁提供**：①②④ 由腿再跑（本轮时间不够，见 §5）；③ 由**机主**锁屏配合（Win+L）。
> **最小修票面（若 ① 坐实「整包不完成」仍在）**：⛔ 本腿不建议在定性前动 `cmd/wisp/**` 夹具——
> 先按 314-ac3 §3 的建议给 `scripts/portable-tests.sh` 的 cli 分支加**终点锚断言**（脚本面，非产品码），
> 让「早停/卡住」在 CI 与推送前门上**响亮变红**，再谈夹具唯一化。

---

## 5. 未完成，卡在哪

- **(c) 的统计不足**：全量一发 6–20 分钟，本轮时间只够 **2 发全量**（run1 卡死、run2 撞 test 超时）；
  **要达到 n≥10 需要单独排一轮**（每发抄四数＋退出码＋终点锚＋卡点栈）。
- **(a) 的锁屏实验⛔ 本轮做不了**：锁屏是机主动作，腿⛔ 不能代劳（理由见 §3(a)）。
- **⛔ 本腿没有复现票面的「整进程死亡」**：两发全量都是「卡」，不是「消失」。
  要拿到「死亡时的退出码」，必须在**死亡真发生的那一发**上取 —— 见 §4 的 ①。

---

## 6. 门禁与边界

- `sh scripts/d22scan.sh` ⇒ **`clean - no D22 ban violations`**（2026-11-13 21:15 现跑，rc=0；
  正控 `runtests.sh -C tools/d22scan ./...`＝`PASS=45 FAIL=0 SKIP=0`、`=== RUN=93`；
  扫描分母 `internal/=403`、`cmd/=77`、`frontend/=123`、`internal/tools/=51`、`internal/risk=21`、`design/=39`、
  `internal/=1005`、`cmd/=276`。**扫描器不扫 `docs/` 与 `.scratch/`** ⇒ 本件的 emoji 不受影响）。
- 写面：本件 ＋ `.scratch/wisp/probes/314/**` ＋ 票 314 Progress log。
- ⛔ 未改任何 `cmd/**`、`internal/**`、`scripts/**`、`.github/**`、`models/**`、`thresholds.go`、golden、契约。
- 临时件：`D:\tmp\wisp\20261113\314\`（只建不删）。

## 7. 撤销口令

删本件 ＋ `.scratch/wisp/probes/314/**` 即复原（零产码、未动任何被跟踪的产品/脚本/契约字节）。
