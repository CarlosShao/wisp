# 429 — `Test384ResidentWatchdogEndToEnd` 恒红：释放 64MB hold 后 3 秒内私有工作集没降到 80MB 以下

- Status: **ready-for-agent**（写面＝`cmd/wisp/resident_watchdog_windows_test.go` 的**台面侧**；⚠ 先只读取证，判清是"窗口不够"还是"页真没还"，再定改法）
- 来源：round 88 整包复跑（`D:\tmp\wisp\20261112\fullclean\cmdwisp.log`，`PASS=593 / FAIL=14`）
- 关联：票 384（看门狗宿主装配，本测试宿主）· 票 273（D43 机器 Sink）· D38(e)／`FreeOSMemory` 级联 · 台账 `A1801`／`A1803`

## 这是什么（人话）

看门狗测试要证明：**进程吃太多内存时它会报警**，然后**内存降回来时它认得回来**。
测试做法：造一块 64MB 的"占用"→ 看门狗报超限 → 释放这块 → 等 3 秒，要求进程私有工作集**真的降到 80MB 以下**（"回到 cap 之内"的真实读数）。
现在：**3 秒内没降到** ⇒ 红。

## 现量（可逐字复跑）

**逐字红句**（`cmd/wisp/resident_watchdog_windows_test.go:230`，经 `waitUntil384` 抛出）：

> `never observed the tree back under the Armed cap (a real within-cap reading)`

**同族前序读数（本轮整包日志里已 PASS 的那些 `384-READ`）**：
```
:193  384-READ real tick lines=3 (the loop goroutine is alive and ticking)
:195  384-READ D38b roster CountByName(watchdog)=1
:217  384-READ Sleeping unload passes=1, unload-action lines=0, no-row alert receipts=2
```
⇒ **看门狗本身是活的**（tick 在跑、名册在、unload pass 跑过、报警收到了）。红的只是**最后这一步回落读数**。

**稳定性**：本轮整包一次红。⚠ **只跑过一次，⛔ 尚未确认稳定红还是偶发** —— 你必须先复跑 ≥2 次定性。

## 机制链（逐行可核）

测试文件：`cmd/wisp/resident_watchdog_windows_test.go`
- `:210` `hold64 := holdTouched384(64 << 20)` —— 造 64MB 并**触碰每一页**（保证真占物理内存）
- `:211-214` 断言 `len(hold64) == 64<<20`
- `:215-216` 等 `watchdog: unload pass complete` ≥1（**PASS 过**）
- `:221-223` 断言 Sleeping 阶段超 cap **不许** dispatch D43 行
- `:224-226` 断言机器仍在 `Sleeping`
- `:227` `hold64 = nil` —— 引用死掉
- `:229` `settleBeforeSample(200)` —— **等 200ms 让采样器看到**
- `:230-233` **红的这一步**：
  ```go
  waitUntil384(t, 3*time.Second, "the tree back under the Armed cap (a real within-cap reading)",
      func() bool {
          return currentTree384(t, rt).PrivateWorkingSetBytes < 80<<20
      })
  ```

## ⚠ 必须自己判清的核心问题（⛔ 不许我替你定）

**"3 秒窗口不够" vs "页真的没还给 OS"——这两件事的修法完全不同：**

- **甲：窗口问题**——Go 把页还给 OS 需要 **GC + 归还** 两步，`hold64 = nil` 只是让引用死掉，
  **不保证** `runtime.GC()` 跑过、更不保证 `debug.FreeOSMemory()` 被调用。
  3 秒内私有工作集不降是**正常的 Go 行为**，不是产品缺陷。
  ⇒ 修法＝**台面侧**：在 `:229` 之后显式触发一次 `runtime.GC()`＋`debug.FreeOSMemory()`
  （或等看门狗自己那条 `FreeOSMemory` 级联跑到），**再**开始 `waitUntil384`。
  ⛔ **不许**直接把 `80<<20` 阈值放宽——那会把"看门狗认得回来"这个承诺废掉。
- **乙：页真没还**——若显式 `FreeOSMemory()` 之后仍然不降，说明**真有东西抓着那 64MB**
  （内存泄漏／引用未断／缓存未清），那是**产品缺陷**，要修产码。

**判据（你必须实测）**：在 `:229` 之后插一发 `runtime.GC(); debug.FreeOSMemory()`（**仓外 overlay，这枚是运行时行为检查，overlay 有效**），
看私有工作集**是否立刻降到 80MB 以下**。
- 降了 ⇒ **甲**（窗口问题），修台面。
- 没降 ⇒ **乙**（真泄漏），去查谁抓着。

**同时记录**：`currentTree384(t, rt).PrivateWorkingSetBytes` 在**释放前**与**释放后 3 秒**的**具体数值**（不是"没降到"这句话）。
有数值才能判"差多少"。

## 要做的（按序）

- [ ] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
  1. **复跑 ≥2 次**定性：
     ```
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^Test384ResidentWatchdogEndToEnd$'"
     ```
     ⚠ **必须用 `D:\work\soft\Git\bin\bash.exe`**（本机 `bash` 是 WSL 转发器，一枚测试都不跑）。
  2. **贴出具体数值**：释放前 `PrivateWorkingSetBytes`、释放后 3 秒的读数、cap 是多少（去读 `currentTree384` 与 cap 常量）。
  3. 用**仓外 overlay** 在 `:229` 后插 `runtime.GC(); debug.FreeOSMemory()`，再跑一次，**看是否立刻降到 80MB 以下**。
  4. 判清甲／乙，给一句话结论＋逐字读数。

- [ ] **AC#1 按判定修（只修被判定那一处）**：
  - **甲**：修台面——`waitUntil384` 之前显式触发 GC＋FreeOSMemory（或等看门狗那条级联），
    **⛔ 不许放宽 `80<<20` 阈值、⛔ 不许把 `3*time.Second` 改大当唯一手段**（那只是把偶发变少见，不是修）。
  - **乙**：修产码让页真还。
  - **必须配突变控制真跑**：把修好的那处删掉 ⇒ 必须回到红；恢复 ⇒ 绿。

- [ ] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。

- [ ] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · **`scripts/d22scan.sh`**（⚠ ⛔ 用 `go run ./tools/d22scan`，独立 module，从根模块跑会失败）。四条读数全贴。
  ⚠ `build/mut_rvw.go` 是**未跟踪残留件**，会让 `go build ./...` 红。**⛔ 不是你的写面，⛔ 不许删**。

- [ ] **AC#4 与票 384／273 的关系写清**：本票只处理这一枚的回落读数，⛔ 不翻票 384／273 任何框；⛔ 不宣称看门狗整族全绿。

## 本票**不**解决

- 不碰 `Test26P*`／`Test365*`／`Test12r2*`（另计）。
- 不改 `thresholds.go`／SLO 阈值（冻结件）。
- 不改 `docs/PLAN.md` 与 `docs/specs/**` 一字。

## 派单约束

- 先只读（AC#0 零写），判定后再定写面；**产码与台面二选一，⛔ 两边都改**。
- ⛔ `git add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；`add`／`commit` 两处显式 pathspec。
- 只 commit，⛔ push。
- ⛔ 翻任何 AC 框。
- 临时件只建不删，落 `D:\tmp\wisp\<YYYYMMDD>\429-<腿码>\`。

## Progress log

### 2026-11-12 20:26–20:5x +08 — 腿 `429-384`（⛔ 未翻任何 AC 框）

**判定＝第三条，票面给的甲／乙两岔都不成立。** 一句话：**释放本身是好的，量错的是「哪个进程」**。

- **AC#0 只读取证（零写，`go test -overlay`，拷贝在仓外 `D:\tmp\wisp\20261112\429-384\`）**：
  - **复跑定性**：单发 `^Test384ResidentWatchdogEndToEnd$` **PASS 1.36s／1.55s** ⇒ 票面「恒红」**单发不可复现**；
    红只在**合跑**出现（round 88 `:230` 3.36s，与 390/392/394/395/397/400/401/402 各件同形）。
  - **三个数值**（overlay 探针）：
    | 读数 | 单发 | 六族合跑 |
    |---|---|---|
    | 释放前 `PrivateWorkingSetBytes` | **92,311,552 B = 88.0MB** | **175,157,248 B = 167.0MB** |
    | 释放后 3 秒 | **13,688,832 B = 13.1MB**（under80=true） | **109,498,368 B = 104.4MB**（under80=**false**，超 **25,612,288 B**） |
    | cap（`memCapArmed` 冻结值） | 94,371,840 B = **90MB** | 同左 |
    | 本用例动手**之前**的基线 | **13,905,920 B = 13.3MB** | **104,820,736 B = 100.0MB** |
  - **判据实测（强制 `runtime.GC(); debug.FreeOSMemory()`）**：单发 13.5→12.9MB（放 667,648 B，本就在线下）；
    **合跑 104.9→100.1MB，只放 5,046,272 B，仍差 80MB 线 20.1MB**。
  - **结论**：⛔ 不是甲（窗口不够——强制释放后仍差 20.1MB，再等也不降）；
    ⛔ 不是乙（64MB 被抓住——hold 在两种语境下都还回来了，合跑里 167.0→100.0MB 掉 **67.0MB**＝整发 hold）；
    ✓ **是基线**：合跑时本用例还没分配任何东西，进程已 100.0MB——语音族（15e/26R3/41a2r1/390r1）
    把真 sherpa-onnx ASR/TTS/KWS 引擎装进同一进程并常驻，那是**活的 cgo 状态不是垃圾**，任何 GC 都还不了。
    同一二进制、同一枚用例、同一发 64MB hold：单发 **13.1MB** vs 合跑 **104.4MB**。
- **修形：只改台面**（产码零字节）。回落读数放进**同一个二进制的干净子进程**
  （`exec.Command(os.Args[0], "-test.run=^Test384…$")`，本仓既有同形先例 `slo_exit_os_156_windows_test.go:322`／
  `secret_argv_windows_test.go:455`）。**`80<<20` 与 `3*time.Second` 一个字节没动**；
  父臂**要求子进程退出码 0 且输出含 `--- PASS`**，两条任一不满足即 `t.Fatalf` 并贴子进程全文
  ⇒ 隔离是测量载体，⛔ 不是跳过断言的通道。
- **突变控制两发（仓外 overlay 真跑）**：
  - **M1 摘掉隔离**（`if false {`，回到修前形状）⇒ **FAIL 3.41s**，红句**逐字同原案**
    `:288: 384: never observed the tree back under the Armed cap (a real within-cap reading) within 3s`，
    同族三读数照常（tick=3／名册=1／unload 趟=1）。
  - **M2 让 hold 真被抓住**（包级 sink 保活，任何 GC 都放不掉）⇒ **FAIL 3.57s**，
    父臂 `:211: the isolated child … failed (exit status 1)`＋子进程红句 `:289: never observed … within 3s`
    ⇒ **80MB 断言在子进程里仍然有牙**，隔离不会把真缺陷洗绿。
  - **还原** ⇒ **PASS 1.44s**，五条 `384-READ` 全在。
- **六族合跑**（`Test43|Test15e|Test41a2r1|Test381r2|Test390r1|Test26R3|Test384`）：
  **修前红（EXIT=1，同句同点）／修后绿（EXIT=0）**，同二进制同命令。
- **133 账本未失效**：`TestAC1AC2DispatchHopGate133` **PASS**，
  `leg slo … covered=test Test384ResidentWatchdogEndToEnd drives settleBeforeSample`
  （用例体搬进 helper 后，`testClosure133` 走本目录测试 helper 仍认账）；
  `TestAC4EveryLegIsNailedOrRuled`／`TestAC3SelfProof133*` 全 PASS。
- **门禁**：`go build ./cmd/wisp/` rc=0 · `go vet ./cmd/wisp/` rc=0 · `gofumpt -l` 零输出 ·
  `scripts/d22scan.sh` rc=0（`clean - no D22 ban violations`）。
- **证据件**：`docs/evidence/s1/429-384-watchdog-fallback.md`。
- **越界**：`git diff --name-only` 只有 `cmd/wisp/resident_watchdog_windows_test.go`（+58/−0）与证据件；
  `thresholds.go`／golden／`docs/PLAN.md`／`docs/specs/**`／`allowlist.txt`／`frontend/**`／`design/**` 全未动，
  `build/mut_rvw.go` 未删。
- **★诚实残余（本票不解决）**：①合跑基线仍是 ~100MB——那是语音引擎常驻的**真实成本**，
  要动得走「引擎是否该常驻」的票，⛔ 不在本票射程；②子进程多花约 1.3 秒；
  ③票面提的「等看门狗那条 `FreeOSMemory` 级联」**未采用**——那条级联属 settle 路径、
  `Sleeping` 无 D43 行时不为这条读数触发，等它等于把「释放后回落」偷换成「级联跑过后回落」。
