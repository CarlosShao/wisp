# 433 — `TestAC246DevLegIgnoresTheTestTaskInjection` 稳定红：`TempDir RemoveAll cleanup: unlinkat … failed`——测试自己的清理与它起的常驻子进程抢文件

- Status: **ready-for-agent**（写面＝`cmd/wisp/resident_task_source_246_windows_test.go` 的**台面侧**；⚠ 先只读取证，判清是"清理顺序错"还是"子进程没退干净"，再定改法）
- 来源：round 89 整包（`D:\tmp\wisp\20261112\round89\full.log`）＋**我单独复跑**（`D:\tmp\wisp\20261112\clean2\run.log`，**FAIL 10.56s**，桌面干净、无并行腿）
- 关联：票 246（resident 任务源）· 票 127（resident leg 台面）· 台账 `A1806`／`A1807`

## 这是什么（人话）

有个测试要证明：**开发版（dev leg）会忽略测试注入的任务入口**。
它起一个**真实的常驻子进程**，然后检查那个子进程**没有**去吃注入的任务。

测试结束时 Go 的 `t.TempDir()` 会自动清理临时目录，**清理失败了**：
> `testing.go:1617: TempDir RemoveAll cleanup: unlinkat C:\Users\swq\AppData\Local\Temp\TestAC246DevLegIgnoresTheTestTaskInjection3034784080\002\wisp\orb\20432-2fc31ac0\boot-1\EB… failed`

**这是测试自己的清理缺陷，⛔ 产品缺陷**——但它**让整枚测试红**，而且**稳定红**。

## 现量（逐字，可复跑）

**红句**（`testing.go:1617`，Go 测试框架自己抛的）：
```
TempDir RemoveAll cleanup: unlinkat C:\Users\swq\AppData\Local\Temp\TestAC246DevLegIgnoresTheTestTaskInjection3034784080\002\wisp\orb\20432-2fc31ac0\boot-1\EB… failed
```
**耗时 10.56s**（我单独跑）。⇒ **稳定红，⛔ 偶发**（我跑一次、round 89 整包一次，两次都红）。

## ★与另一枚的区别（⛔ 混为一谈）

同轮还有 `TestAC4FocusReturnToPriorWindowGap33r5` 也红，但**那枚我单独跑 PASS 1.17s** ⇒ **那枚是桌面焦点争抢（环境件）**。
**本票这一枚单独跑仍红** ⇒ **是台面自身的缺陷**。两枚根因不同。

## 机制链（现读，供你起步；⛔ 别照抄结论）

- 测试文件：`cmd/wisp/resident_task_source_246_windows_test.go`
- `:375`／`:436` 是 `t.Logf("resident sink %s: …")` —— **⛔ 断言，只是日志**（我一开始误读成红因了，别重复这个错）。
- `:432-434`：`if !leg.exitedWithin(400) { leg.stop(); t.Fatalf("AC#7 RED: the leg did not leave through its own path.\n%s", leg.console()) }` —— **本枚的断言在这条**（走自己的路退出）。
- `:436-441`：读 sink、检查有没有 `"shutdown step failed"`／`"abandoned wait"`。
- **红因不在任何断言里，而在测试结束后的 `t.TempDir()` 自动清理**。
- 清理失败的路径：`…\002\wisp\orb\20432-2fc31ac0\boot-1\EB…` —— 看着像 **`orb` 子目录下 `boot-1` 里的一个 EB 开头文件**（可能是 WebView2 的 `EBWebView` 目录）。

## ⚠ 必须自己判清的核心问题（⛔ 不许我替你定）

**"文件被谁占着"？三种可能，修法不同：**

- **甲：子进程没退干净**——`leg.stop()` 只是发了停止信号，**进程可能还没真正退出**就开始清理。
  判据：在清理之前加一个"等子进程真正退出"的等待（`leg.exitedWithin(...)` 已有，看它是否被正确使用）。
- **乙：WebView2 / 别的后台线程占着 `EBWebView`**——这类目录 Windows 上常有句柄残留。
  判据：看那个 EB 开头的东西是不是 `EBWebView`；若是，可能需要 `t.Cleanup` 里先 `runtime.GC()`＋等句柄释放。
- **丙：测试用了 `t.TempDir()` 但子进程的 cwd/temp 指向那里**，父进程清理时子进程还活着。

**你必须先跑出**：清理失败的那个**具体文件名**（日志被截断了，`EB…` 后面是什么），以及**那一刻子进程是否已退出**。

## 要做的（按序）

- [ ] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
  1. **复跑 ≥2 次**确认稳定红（我已有 1 次 10.56s）：
     ```
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC246DevLegIgnoresTheTestTaskInjection$'"
     ```
     ⚠ **必须用 `D:\work\soft\Git\bin\bash.exe`**（本机 `bash` 是 WSL 转发器，一枚测试都不跑）。
  2. **抄出完整红句**（日志被截断，`EB…` 后面必须补全——去 `D:\tmp\wisp\<date>\<腿码>\` 的完整日志里找，或用 `-test.v` 重跑）。
  3. 读 `leg.stop()`／`leg.exitedWithin()` 的实现，写清**它们是否保证子进程真退出**。
  4. 判清甲／乙／丙，给一句话结论＋`文件:行`。
  5. **⚠ 确认这枚⛔ 桌面焦点件**（与 `TestAC4Focus*` 不同）：它只起常驻子进程，⛔ 抢前台。

- [ ] **AC#1 按 AC#0 的判定修（只修被判定那一处）**：
  - **甲/丙**：在 `t.Cleanup` 或测试末尾**显式等子进程真退出**再让 `TempDir` 清理（`t.Cleanup` 是 LIFO，注册顺序要对）。
  - **乙**：处理 WebView2 句柄残留（等释放／重试删除）。
  - **⛔ 不许**把 `t.TempDir()` 换成手工 `os.MkdirTemp` + 不清理——那会把"清理失败"藏起来，⛔ 修。
  - **必须配突变控制真跑**：把修好的那处删掉 ⇒ 必须回到红；恢复 ⇒ 绿。

- [ ] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。

- [ ] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · **`scripts/d22scan.sh`**（⚠ ⛔ 用 `go run ./tools/d22scan`，独立 module）。四条读数全贴。
  ⚠ `build/mut_rvw.go` 是未跟踪残留件，会让 `go build ./...` 红。**⛔ 不许删**。

- [ ] **AC#4 与票 246 的关系写清**：本票只处理这一枚的清理缺陷，⛔ 不翻票 246 任何框；⛔ 不宣称 `TestAC246*` 全族绿。

## 本票**不**解决

- 不碰 `TestAC4Focus*`（那是桌面焦点环境件，**已确认单独跑 PASS 1.17s**，需另立"真窗族独占"调度票）。
- 不改 `thresholds.go`／SLO 阈值（冻结件）。

## 派单约束

- 先只读（AC#0 零写），判定后再定写面。
- ⛔ `git add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；`add`／`commit` 两处显式 pathspec。
- 只 commit，⛔ push。
- ⛔ 翻任何 AC 框。
- 临时件只建不删，落 `D:\tmp\wisp\<YYYYMMDD>\433-<腿码>\`。
- **仪器**：运行时行为检查 ⇒ `go test -overlay`（拷贝放**仓外**）有效。
