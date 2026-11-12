# 票 429 — 384 看门狗回落读数（实现腿 429-384）

> 写面：`cmd/wisp/resident_watchdog_windows_test.go`（**台面侧，一处**）＋本件＋票面 Progress log 一条。
> 产码零字节：`internal/**` 与 `cmd/wisp/**` 的非测试文件全部未动。
> 临时件（只建不删）：`D:\tmp\wisp\20261112\429-384\`。
> 仪器：`go test -overlay`（拷贝在仓外），跑测试一律经 `D:\work\soft\Git\bin\bash.exe`。

## 0. 一句话结论

**既不是甲（窗口太短），也不是乙（64MB 被抓住）——是第三条：断言量的是「整个进程的历史」，不是「这一次释放」。**

释放本身在两种语境下都**完美工作**（六族合跑里树从 167.0MB 掉到 100.0MB，掉的正是那 64MB）。
红的原因是**基线**：六族合跑时本进程在做本用例之前就已经是 **100.0MB**，
比 80MB 这条线高 20MB、比 Armed cap 还高 10MB —— 因为语音族（15e / 26R3 / 41a2r1 / 390r1）
把真 sherpa-onnx ASR/TTS/KWS 引擎装进了同一个进程并留在常驻。那是**活的 cgo 状态，不是垃圾**，
任何 GC 都还不了。80MB 这条线低于该进程能到达的地板。

## 1. AC#0 只读取证（⛔ 零写）

### 1.1 复跑定性（≥2 次）

| 发 | 命令 | 结果 |
|---|---|---|
| 1 | 单发 `^Test384ResidentWatchdogEndToEnd$` | **PASS 1.36s** |
| 2 | 单发同上 | **PASS 1.55s** |

⇒ 票面「恒红」**单发不可复现**。红只在**合跑**出现（round 88 整包 `:230` 3.36s、390/392/394/395/397/400/401/402 各件同形）。

### 1.2 三个数值（overlay 探针，仓外）

| 读数 | 单发 | 六族合跑（`Test43\|Test15e\|Test41a2r1\|Test381r2\|Test390r1\|Test26R3\|Test384`） |
|---|---|---|
| **释放前** `PrivateWorkingSetBytes` | **92,311,552 B = 88.0MB** | **175,157,248 B = 167.0MB** |
| **释放后 3 秒** | **13,688,832 B = 13.1MB**（under80=true，overBy −70,197,248） | **109,498,368 B = 104.4MB**（under80=**false**，overBy **+25,612,288**） |
| **cap（`memCapArmed` 冻结值）** | 94,371,840 B = **90MB** | 同左 |
| 本用例动手**之前**的基线 | **13,905,920 B = 13.3MB** | **104,820,736 B = 100.0MB** |

### 1.3 判据：插 `runtime.GC(); debug.FreeOSMemory()` 后是否立刻降？

| 语境 | 强制释放前 | 强制释放后 | 实际放掉 | 降到 80MB 以下？ |
|---|---|---|---|---|
| 单发 | 14,196,736 B (13.5MB) | 13,529,088 B (12.9MB) | 667,648 B | **是**（本来就在线下） |
| 六族合跑 | 109,965,312 B (104.9MB) | **104,919,040 B (100.1MB)** | **5,046,272 B** | **否** |

**这是决定性的一发**：合跑里强制 GC＋FreeOSMemory 只放掉 **5.0MB**，树仍停在 **100.1MB**，
离 80MB 还差 **20.1MB**。再多等几秒也不会有用——那些页不是垃圾。

### 1.4 甲／乙判定

- **⛔ 不是甲（窗口太短）**：强制释放后仍差 20.1MB。3 秒改 30 秒同样不降，因为页不属于本用例。
- **⛔ 不是乙（64MB 被抓住）**：hold 在两种语境下**都还回来了**。合跑里 167.0 → 100.0MB，
  掉 **67.0MB**＝64MB hold ＋ 记账取整。释放机制本身**零缺陷**。
- **✓ 是基线**：合跑时本用例还没分配任何东西，进程已是 100.0MB。
  同一个二进制、同一枚用例、同一发 64MB hold：单发 **13.1MB**，合跑 **104.4MB**。

断言的形状问题：它是**绝对读数**，所以量到的是「这个进程之前跑过什么」，
而不是本用例拥有的那条机制。

## 2. AC#1 修（只修一处）

### 2.1 选形理由

修法＝**把回落读数放进一个干净的子进程**（同一个二进制的全新实例，只跑本用例）。

为什么这是最有价值的形状（按主流 harness 的做法，⛔ 不是「动作最小」）：

1. **阈值与窗口一个字节没动**：`80<<20` 与 `3*time.Second` 原样保留。
   票面禁止的两条（放宽阈值／只放大窗口）都没有碰。
2. **承诺保住了**：本用例存在的理由——「释放之后能看到一个真的 within-cap 读数」——
   仍然被逐字执行，只是**在它能成立的地方**执行。子进程里那条断言照旧会红（见 §3 M2）。
3. **主流 harness 的形状**：Go 生态里「需要干净进程基线」的用例一律**重入自身二进制**
   （`exec.Command(os.Args[0], "-test.run=…")`），本仓已有同形先例
   （`slo_exit_os_156_windows_test.go:322`、`secret_argv_windows_test.go:455`）。
   这不是新发明，是把仓里既有的那把手尺用到该用的地方。
4. **不是把红藏起来**：父进程**要求子进程退出码为 0 且输出里有 `--- PASS`**，
   两条任一不满足即 `t.Fatalf` 并把子进程全文贴出来。隔离是**测量载体，不是跳过断言的通道**。

### 2.2 改了什么

- `Test384ResidentWatchdogEndToEnd` 变成两臂：父臂 `run384IsolatedChild`，子臂 `run384EndToEnd`（原体一字未改）。
- 新增 `env384IsolatedChild` 常量与 `run384IsolatedChild`；子进程读数**回贴父日志**（成功时也看得见）。
- import 加 `os`、`os/exec`。

## 3. AC#1 突变控制（两发，真跑，仓外 overlay）

| 突变 | 内容 | 读数 | 判 |
|---|---|---|---|
| **M1** 摘掉隔离 | `if os.Getenv(env384IsolatedChild) == "" {` → `if false {`（＝回到修前形状，父进程自己量） | **FAIL 3.41s**，红句**逐字同原案**：`resident_watchdog_windows_test.go:288: 384: never observed the tree back under the Armed cap (a real within-cap reading) within 3s`；同族三读数照常（tick=3／名册=1／unload 趟=1） | **必红 ✓** 证明红是这一处造成的 |
| **M2** 让 hold 真被抓住 | `hold64 = nil; _ = hold64` → `keepAlive384 = hold64`（包级 sink 保活，任何 GC 都放不掉） | **FAIL 3.57s**，父臂红句：`resident_watchdog_windows_test.go:211: 384: the isolated child that measures the fallback reading failed (exit status 1)`；子进程红句：`:289: 384: never observed the tree back under the Armed cap … within 3s` | **必红 ✓** 证明 80MB 断言在子进程里**仍然有牙**，且父进程如实传播 |
| **还原** | 撤 overlay，跑修后文件 | **PASS**（1.44s，五条 `384-READ` 全在） | **必绿 ✓** |

⇒ 修的那一处**既必要又充分**：删掉它回到原红（同句同点），恢复它即绿；
而把「页真没还」这个真缺陷注入进去，隔离**不会**把它洗绿。

## 4. AC#2 越界

`git diff --name-only`（见 §6 逐行）只含本腿三件。⛔ 未出现 `docs/PLAN.md`／`docs/specs/**`／
`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**`。
`build/mut_rvw.go` 未删、未动。

## 5. AC#3 门禁四读数

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/wisp/` | **rc=0**（无输出） |
| vet | `go vet ./cmd/wisp/` | **rc=0**（无输出） |
| gofumpt | `$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp/resident_watchdog_windows_test.go` | **零输出，rc=0** |
| d22scan | `scripts/d22scan.sh` | **rc=0 clean**：`d22scan: clean - no D22 ban violations`；分母 `bans#1-5 internal/=403 cmd/=77`、`#6 frontend/=123`、`#7 internal/tools/=52`、`#10 internal/risk/=21`、`#8 design/=39 frontend/=123 internal/=1004 cmd/=276` |

附带复跑（确认没连坐）：
- 六族合跑**修后绿**（EXIT=0），修前同命令同二进制**红**。
- `TestAC1AC2DispatchHopGate133` **PASS**：`leg slo … covered=test Test384ResidentWatchdogEndToEnd drives settleBeforeSample`
  —— 用例体搬进 helper 后，133 的闭包（`testClosure133` 走本目录测试 helper）**仍然认账**，账本行未失效。
- `TestAC4EveryLegIsNailedOrRuled` / `TestAC3SelfProof133*` **PASS**。

## 6. AC#4 与票 384／273 的关系

本票**只**处理这一枚的回落读数。⛔ 未翻票 384／273 任何框；
⛔ 不宣称看门狗整族全绿（`Test384StartResidentWatchdogRegistersRealLoop` 与
`Test384PanelWebviewComponentArms` 本就是绿的，与本票无关）。
生产装配、D43 #9／#35 映射、白名单 unload 趟——**一字未动**。

## 7. 诚实边界（本票没做的）

1. **没有让合跑基线降下来**。六族合跑时本进程仍是 ~100MB，那是语音引擎常驻的**真实成本**，
   不是本票的病。要动它得走「语音引擎是否该常驻」的票，⛔ 不在本票射程。
2. **子进程多花约 1.3 秒**（用例从 ~1.3s 变 ~1.4s＋子进程启动）。代价换来的是
   「这条断言在任何合跑顺序下都量同一件事」。
3. 票面提的「或等看门狗自己那条 `FreeOSMemory` 级联」**没有采用**：那条级联属 settle 路径，
   在 `Sleeping` 无 D43 行时不会为这条读数触发；改成等它等于把读数接到另一条语义上，
   会把「释放后回落」偷换成「级联跑过后回落」。子进程隔离不改语义。
