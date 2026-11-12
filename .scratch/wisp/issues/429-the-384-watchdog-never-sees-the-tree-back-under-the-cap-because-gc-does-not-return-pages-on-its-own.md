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
