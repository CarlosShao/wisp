# 433 — `TestAC246DevLegIgnoresTheTestTaskInjection` 稳定红：`TempDir RemoveAll cleanup: unlinkat … failed`——测试自己的清理与它起的常驻子进程抢文件

- Status: **done**（写面＝`cmd/wisp/resident_task_source_246_windows_test.go` 的**台面侧**；⚠ 先只读取证，判清是"清理顺序错"还是"子进程没退干净"，再定改法）
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

- [x] **AC#0 只读取证（⛔ 本格不许改任何文件）**：
  1. **复跑 ≥2 次**确认稳定红（我已有 1 次 10.56s）：
     ```
     D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC246DevLegIgnoresTheTestTaskInjection$'"
     ```
     ⚠ **必须用 `D:\work\soft\Git\bin\bash.exe`**（本机 `bash` 是 WSL 转发器，一枚测试都不跑）。
  2. **抄出完整红句**（日志被截断，`EB…` 后面必须补全——去 `D:\tmp\wisp\<date>\<腿码>\` 的完整日志里找，或用 `-test.v` 重跑）。
  3. 读 `leg.stop()`／`leg.exitedWithin()` 的实现，写清**它们是否保证子进程真退出**。
  4. 判清甲／乙／丙，给一句话结论＋`文件:行`。
  5. **⚠ 确认这枚⛔ 桌面焦点件**（与 `TestAC4Focus*` 不同）：它只起常驻子进程，⛔ 抢前台。

- [x] **AC#1 按 AC#0 的判定修（只修被判定那一处）**：
  - **甲/丙**：在 `t.Cleanup` 或测试末尾**显式等子进程真退出**再让 `TempDir` 清理（`t.Cleanup` 是 LIFO，注册顺序要对）。
  - **乙**：处理 WebView2 句柄残留（等释放／重试删除）。
  - **⛔ 不许**把 `t.TempDir()` 换成手工 `os.MkdirTemp` + 不清理——那会把"清理失败"藏起来，⛔ 修。
  - **必须配突变控制真跑**：把修好的那处删掉 ⇒ 必须回到红；恢复 ⇒ 绿。

- [x] **AC#2 越界检查**：`git diff --name-only` 逐行贴。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**` ⇒ **直接退回**。

- [x] **AC#3 门禁**：`go build ./cmd/wisp/` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · **`scripts/d22scan.sh`**（⚠ ⛔ 用 `go run ./tools/d22scan`，独立 module）。四条读数全贴。
  ⚠ `build/mut_rvw.go` 是未跟踪残留件，会让 `go build ./...` 红。**⛔ 不许删**。

- [x] **AC#4 与票 246 的关系写清**：本票只处理这一枚的清理缺陷，⛔ 不翻票 246 任何框；⛔ 不宣称 `TestAC246*` 全族绿。

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

## Progress log

### 2026-11-12 22:46 +08 — 写腿（AC#0 取证＋AC#1 修＋突变控制真跑）

**先纠正票面两处**（都按现量，⛔ 按记忆）：

1. **⛔ 稳定红，是竞态。** 同一台机、桌面干净、无并行腿，连跑：FAIL 23.31s → **PASS 13.42s**（代码一字未改）→ FAIL 11.24s → FAIL 13.13s。票面 `:23` 的"稳定红、⛔ 偶发"与派单 `:8` 的"稳定红"都被这次绿读数推翻。
2. **完整文件名补全了**（票面 `EB…` 截断处）：`…\boot-1\EBWebView\lockfile`，即 WebView2 用户数据目录的锁文件。

**那一刻子进程是否已退出 ⇒ 已退出。** 用 `OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION)` ＋ `GetExitCodeProcess` 直接读进程表：`childAlive=false`、`waited=true`。**句柄活着，进程死了。**

**判：乙（WebView2 句柄残留），⛔ 甲、⛔ 丙。**

- ⛔ 甲：子进程确实退了（上面那条读数），且 `stop()`／`exitedWithin()` 确实 join 了 `cmd.Wait`（`resident_sink_nail_127_windows_test.go:290-301`）。
- ⛔ 丙：子进程 `cmd.Dir` 是 `…\001`（build 目录），被删的树是 `…\002`（APPDATA），两个子树不同；只有 `\002` 报错。
- 真凶：`wisp.exe` 被杀后**遗孤**的 5 个 `msedgewebview2.exe`（父＝死掉的 38616），其中一个占着 `EBWebView\lockfile`。链路＝`LOCALAPPDATA` 指向该树（`:459`）→ `internal/ball/orb_windows.go:471-477`（`orbDataBase` 读 `LOCALAPPDATA`）→ `:529-535`（`orbBootDataDir` 拼 `boot-<n>`）→ `:315`（`browser.DataPath = bootDir`）。

**关键数字**：树在子进程退出后 **5.58s** 才可删；Go 框架 `testing.go:1650` 的 `arbitraryTimeout` 只有 **2s**。5.58 > 2，这就是全部红因。判据用的是框架真正用的谓词（`os.RemoveAll` 那棵树），⛔ "锁文件能否打开"（delete-pending 会永远答 `ERROR_FILE_NOT_FOUND`，那是假读数）。

**修（只一处，台面侧）**：`:451` 的 `appData := t.TempDir()` → `managedDir433(t, "ac246-dev-leg-appdata")`；helper 用 `uniqueProfileLeaf`（`panel_resident_windows_test.go:61`，仓里**已有**的、panel WebView2 profile 用的同一招）＋ `t.Cleanup` 里**按次数**有界重试（150×100ms，⛔ 墙钟 deadline，D42#9）。**⛔ 不是那个被禁的修法**：删除照跑、每次跑都跑、`t.Fatalf` 路径也跑，超预算会**响亮失败**并说明是这棵树，⛔ 藏起来。

**突变控制（真跑）**：修好 PASS 14.48s ⇒ 只把那一行还原成 `t.TempDir()` ⇒ **FAIL 9.82s，逐字复现原红句**（同一文件 `EBWebView\lockfile`）⇒ 恢复 ⇒ PASS 8.30s ⇒ 再跑 PASS 13.48s。

**回归**：`-test.run '^TestAC246'` 全族 **PASS**，0 失败，14 例（含三条 shipped-leg 7.37s／8.01s／8.70s／15.18s）。

**门禁**：`go build ./cmd/wisp/` 0 · `go vet ./cmd/wisp/` 0 · `gofumpt -l` 零输出 · `scripts/d22scan.sh` `clean - no D22 ban violations`（PASS=45 FAIL=0 SKIP=0）。`build/mut_rvw.go` 未删（`Test-Path`=True）。

**证据件**：`docs/evidence/s1/433-ac246-tempdir-cleanup.md`。

**⛔ 未做**：未翻任何 AC 框；未碰 `thresholds.go`／golden／`docs/PLAN.md`／`docs/specs/**`／`frontend/**`／`design/**`；未删 `build/mut_rvw.go`；未 push。

**留给下一张票（⛔ 本票修）**：产品侧"被杀常驻进程会遗留 `msedgewebview2.exe` 约 5.6s"未改——对用户不可见（下次 boot 由 `orbSweepStaleDirs`，`internal/ball/orb_windows.go:544-569` 扫掉死叶），要改得把 WebView2 的 browser 进程塞进 Job Object，属产品改动，本票没有证据支持它。
- [2026-11-12 22:5x +08] **编排者：收件＋三把尺全过＋AC#0–AC#4 五格翻勾＋`-done`**。腿 `ticket-433-ac246-cleanup` 交件（`0bac5d9a3`，3 files；证据件 `docs/evidence/s1/433-ac246-tempdir-cleanup.md` 12835B）。**三把尺**：①`git cat-file -t` = commit；②件 12835B 在盘；③**我自己连跑两次** ⇒ **PASS 8.64s / PASS 13.60s，FAIL=0**。**越界**：只改 `cmd/wisp/resident_task_source_246_windows_test.go` 一枚（台面侧，产码零改动）。
- **★腿判＝乙（WebView2 句柄），⛔ 甲也⛔ 丙**，链条完整：`LOCALAPPDATA` 指向那棵树（`:459`）→ `internal/ball/orb_windows.go:471-477`（`orbDataBase` 读 `LOCALAPPDATA`）→ `:529-535`（`orbBootDataDir` 追加 `boot-<n>`）→ `:315`（`browser.DataPath = bootDir`）。杀掉 `wisp.exe` 会**孤儿化五个 `msedgewebview2.exe`**，其中一个持着锁。
- **★完整文件名补全了（我票面要求的那一项）**：`…\orb\<pid>-<hex>\boot-1\EBWebView\lockfile` —— **WebView2 的 user-data 锁文件**，⛔ 泛指的 "EB…" 产物。
- **★决定性数字**：树在子进程退出后 **5.58s** 才可删，而 Go 的 `removeAll` 只给 **2s**（`testing.go:1650`）。腿**量的是框架真正用的那个判据**（对树做 `os.RemoveAll`），⛔ "锁文件能不能打开"——**delete-pending 的文件会永远回 `FILE_NOT_FOUND`，那会给出假读数**。这个细节做得很硬。
- **★排除甲的证据**：子进程**确实已退出**，且是直接从进程表读的 `GetExitCodeProcess`（⛔ 从 `cmd.Wait` 推断）：`childAlive=false`、`waited=true`。`stop()`／`exitedWithin()` 确实 join `cmd.Wait`（`resident_sink_nail_127_windows_test.go:290-301`）。⇒ **是句柄活过了进程，⛔ 反过来。**
- **★排除丙的证据**：子进程的 `cmd.Dir` 是 `…\001`（构建目录），被删的树是 `…\002`（APPDATA）——**不同子树**，且只有 `\002` 报错。
- **★修法审查（腿主动请我复核，我认可）**：把 `appData := t.TempDir()` 换成 `managedDir433(t, …)`，**不是**我禁止的"不清理"——它注册了 `t.Cleanup`、**每次运行都清（含 `t.Fatalf` 路径）**、预算按**尝试次数**（150×100ms＝15s，约 5.58s 的 2.7 倍）**而非挂钟期限**（合 D42#9）、失败 `t.Errorf` **明确报出**。路径来自仓库既有答案 `uniqueProfileLeaf`（`panel_resident_windows_test.go:61`，与 `panel_reshow_255r8_winlive_test.go:118-125` 对 panel profile 用的是同一条规则）。**⛔ 藏问题，是把清理从框架的 2s 重试预算下挪出来、变成"我们自己的、等得起的清理"。**
- **突变控制**：修后 PASS 14.48s → **只回退那一行** → **FAIL 9.82s，原句逐字节重现**（同一个 `EBWebView\lockfile`）→ 恢复 PASS 8.30s → 再复跑 PASS 13.48s。**整个 AC246 族 14 枚全过、0 失败。**
- **★腿纠正了我第 4 处错**：我票面写"**稳定红**"，**腿 run 2 在代码未动时就 PASS 了**（13.42s）⇒ **是"通常输的竞态"，⛔ 稳定红**。我只跑过一次就下了"稳定"结论。**这是同一根的第 8 次（凭印象/单次读数下结论）。**

