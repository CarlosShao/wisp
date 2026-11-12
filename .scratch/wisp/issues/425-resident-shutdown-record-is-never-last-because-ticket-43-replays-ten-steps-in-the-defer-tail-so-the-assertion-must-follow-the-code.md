# 425 — `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` 是「断言落后于产码」：票 43 加了 defer 尾部十步回放，把 `shutdown step` 永远挤出了末位

- Status: **ready-for-agent**（写面＝`cmd/wisp/resident_sink_nail_127_windows_test.go` 的**断言侧唯一一处**；⚠ 本票**不碰**任何产码、不碰 `internal/proc/**`、不碰 `cmd/wisp/resident_windows.go`、不改 `docs/PLAN.md` 与 `docs/specs/**` 一字）
- 来源：round 88 整包复跑（`D:\tmp\wisp\20261112\fullclean\cmdwisp.log`，`PASS=593 / FAIL=14`）＋ 只读腿 `resident-ac1-verdict` 现读取证（交件 `829c6da1c`，证据件 `docs/evidence/s1/resident-ac1-verdict.md` §2.3）
- 关联：票 43（引进回放的产码腿，`5056a365c`）· 票 127（本测试的宿主票）· 票 408（UDF 进程唯一性，同一整包面）· 台账 `A1801`

## 这是什么（人话）

有一枚测试要证明：**常驻进程退出时，先把"关停每一步"写进日志，再关日志本身**。
它的检查办法是——**看日志文件的最后一条**是不是 `shutdown step ...`。

问题：票 43 之后，退出路径上多了一段"把十个关停步骤再逐条回放一遍"的日志（`wisp: D38(e) step N ...`）。
这段回放**排在所有关停记录之后**，所以**文件最后一条永远是 `D38(e) step 10 (exit)`**。
⇒ `shutdown step` **永远不可能**出现在末位 ⇒ 这枚测试**从此必红**，与产品好坏无关。

## 现量（可逐字复跑）

**逐字红句**（`cmd/wisp/resident_sink_nail_127_windows_test.go:635`）：

> `the last record in the sink is "wisp: D38(e) step 10 (exit): ok elapsed=0s", want a shutdown-step record: the shutdown must be logged before the sink closes`

**★旁证（这条最硬）**：同一枚测试自己的 `t.Logf` 逐字说

> `RESIDENT LEG RECORDED ON DISK: 5 shutdown record(s), steps [1 2 5 6 7], first="shutdown step skipped (module not present)" last="sh..."`

⇒ **5 条 shutdown 记录确实在盘上**，`len(trail) != 0` 那关（`:629`）**已经过了**。
红的**只是**"它在不在末位"这一件事。

**稳定性**：现跑 **5/5 稳定红**（7.80s／6.96s／6.79s／5.82s／5.97s），voice 与 novoice 两种台面都红，**时长、负载、机器无关**。

**机制链（逐行可核）**：
- 断言：`:635` 取 `recs[len(recs)-1]`，要求含 `residentShutdownRecord`（`= "shutdown step"`，`:131`）。
- `wisp: D38(e) step N ...` 由 `cmd/wisp/resident_windows.go:246` 打印，位置＝**defer 尾部，十步全量回放**。
- `shutdown step ...` 由 `internal/proc/shutdown.go:130` 打印，位置＝**序列执行当时**。
- ⇒ 回放**必然排在**序列当时那些记录**之后**。
- **引进者**：`5056a365c`（2026-10-31，票 43 落地腿 `43-r1`）在 `resident_windows.go` 的 defer 里新加了 `:246` 这个回放循环（`git show 5056a365c -- cmd/wisp/resident_windows.go` 的 `@@ -116,14` 那段）。**断言没跟上。**

## 要做的（按序，写面只有一处）

- [ ] **AC#0 先复现（⛔ 不许直接改）**：用 carrier 跑一次，贴出逐字红句与耗时。
  ```
  D:\work\soft\Git\bin\bash.exe -c "scripts/wisp-cli-staged-tests.sh '-test.run' '^TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink$'"
  ```
  ⚠ **必须用 `D:\work\soft\Git\bin\bash.exe`**：本机 `bash` 在 PATH 上解析到 `C:\Windows\System32\bash.exe`（WSL 转发器，`execvpe(/bin/bash) failed`，**一枚测试都不跑**）。
  判定通过条件：红句逐字含 `the last record in the sink is "wisp: D38(e) step 10`。

- [ ] **AC#1 改断言，保留判据强度（⛔ 不许放宽成"存在即通过"）**：
  **不许**改成 `len(trail) != 0` 就算过——那关本来就过了，等于把这条测试废掉。
  正确改法（二选一，由你按"最符合主流 harness"自己拍并写理由）：
  - **甲**：断言「`trail` 的最后一条必须是 `recs` 里**最后一条 `shutdown step` 记录**」——即先按 `residentShutdownRecord` 过滤出 `trail`，再要求 `trail` 末位在 `recs` 里的**相对次序**晚于任何非回放记录。
  - **乙**：断言「`trail` 非空 **且** 其步序单调递增 **且** `trail` 末位的记录时刻早于**第一条** `D38(e) step` 记录」——直接钉住"关停记录发生在回放之前"这个真正的产品承诺。
  ⚠ 无论选哪个，**必须写清"为什么它仍然能抓住产品退化"**：把 `internal/proc/shutdown.go:130` 那行 `slog.Info` 删掉 ⇒ 新断言必须变红。**这一发（突变控制）必须实测，不许只写在注释里。**

- [ ] **AC#2 突变控制实测（负向必响）**：用 `go test -overlay`（拷贝放**仓外**，⛔ 不改仓内产码）把 `internal/proc/shutdown.go:130` 的 `slog.Info` 抹掉 ⇒ 新断言**必须红**；恢复 ⇒ **必须绿**。两发读数都贴进 Progress。

- [ ] **AC#3 越界检查**：`git diff --name-only` 只许出现 `cmd/wisp/resident_sink_nail_127_windows_test.go` 一枚。出现 `docs/PLAN.md`／`docs/specs/**`／`internal/observe/thresholds.go`／golden／`tools/d22scan/allowlist.txt`／`frontend/**`／`design/**`／任何 `*.go` 产码（除本测试文件）⇒ **直接退回**。

- [ ] **AC#4 门禁**：`go build ./...` · `go vet ./cmd/wisp/` · `gofumpt -l` 零输出 · `tools/d22scan` 不新增违规。四条读数全贴。

- [ ] **AC#5 与票 127 的关系写清**：本票**只修这一枚的断言**，⛔ 不宣称票 127 的其他格已满足；⛔ 不翻票 127 的任何框。

## 本票**不**解决

- 不修 `TestAC1ResidentLegInstallsItsLogListenerOnDisk`（预算不够，另一根因）与 `TestAC1ResidentLegOutlivesItsOwnLogFailure`（偶发）——那两枚另计。
- 不碰 `resident_windows.go` 的回放顺序（那是票 43 的既定行为，改它＝改产品行为，需单独立票）。
- 不做"把回放挪到 sink 关闭之后"——那会让 `D38(e)` 记录丢失，方向相反。

## 派单约束

- 写面：**仅** `cmd/wisp/resident_sink_nail_127_windows_test.go`。
- ⛔ `git add -A`／`add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；`add`/`commit` 两处显式 pathspec。
- 只 commit，⛔ push（推送由编排者做）。
- 临时件只建不删，落 `D:\tmp\wisp\<YYYYMMDD>\425-<腿码>\`。
- 跑测试**只经 carrier**，⛔ 裸 `go test ./cmd/wisp`（DLL 遮蔽会给出假读数）。
