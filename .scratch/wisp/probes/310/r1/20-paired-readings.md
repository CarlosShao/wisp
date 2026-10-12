# `20-paired-readings.md` — 票 310 `AC#2` 反恒真成对两发（腿 `310-r1`）

窄用例＝`internal/config` 的 `TestRoundTripLoadMarshalLoad`（定义在 HEAD blob `internal/config/boundary_test.go:186`，那枚 `sealableTempDir124(t)`→`t.TempDir()`）。名册尺（每一发同批跑）＝

```
powershell -NoProfile -Command 'Get-ChildItem -Path $env:TEMP -Directory | Where-Object { $_.Name -like "Test*" } | ForEach-Object { $_.Name }' | sort
comm -13 <before> <after>      # 具名，⛔ 只报枚数
```

★规矩：`go test` 的输出**一律走独立文件**（⛔ 进读数件——编排者 10-12 踩过那枚"表头被自家输出冲掉"）；⛔ 拿"0 枚"当读数——每一发都同批发"那一形真发生了"的证据行；两发都⛔ 成立的才标〔作废〕并留非零 rc。

## 1. (a) 撤掉那一步 ⇒ `%TEMP%` 顶层当场多出一枚（同一枚用例、同一条形）

本腿起手的仓＝HEAD `7e6c1579`，那一枚步骤**本来就不存在**，所以 (a) 侧⛔ 需要变异，只需⛔ 设 `GOTMPDIR`（`env -u GOTMPDIR`）。

| # | 发 | 名册 before→after | 顶层新增（具名） | **"那一形真发生了"的同行证据** | rc | 件 |
|---|---|---|---|---|---|---|
| a1 | 超时 panic（冷台面，首试） | 55→55 | 0 枚 | 有 `panic: test timed out after 2s`（件内 `:2`），**但 `run_lines=0`、栈里 `TestRoundTripLoadMarshalLoad (0s)`** ⇒ 帽落在**第一次 `t.TempDir()` 之前**（冷编译吃掉了 2 s） | `shot_rc=1` | `logs/10a-shot.txt` |
| | | | | ★**这一发自标〔作废〕**：坏世界⛔ 被造出来（⛔ 拿它那枚 0 当任何一边的读数） | | |
| a2 | 超时 panic（先暖台面再打，正发） | **55→56** | **`TestRoundTripLoadMarshalLoad2461500944`** | 逐字 `panic: test timed out after 2s`（`logs/10a-shot2.txt:517`）＋同批 `run_lines=258`／`pass_lines=257` ⇒ **帽落在第 258 遍的运行之中**（前 257 遍各自收尾删干净） | `shot_rc=1` | `logs/10a-warm.txt`（`warm_rc=0`、`--- PASS`）＋`logs/10a-shot2.txt`＋`logs/10a-before.txt`／`-after.txt`／`-new.txt` |
| a3 | 外力 `TerminateProcess`（复用普查腿台件原形 `probes/310/a1/logs/20-shape-c-watch.ps1`） | 56→57 | `TestRoundTripLoadMarshalLoad2236627062`（由它自报的 `test_star_dirs_now=57` 与外部门册作差对齐） | 同进程逐字 `world=CONSTRUCTED dir=TestRoundTripLoadMarshalLoad2236627062`／`external_kill_issued at=11:49:06.508`／`after_kill target_pid_still_present=0` | `watch_rc=1`（该台件自报 `verdict=NO-LEAK`） | `logs/10c-watch.txt` |
| | | | | ★⚠ 那枚台件的自报判语**读错了对象**（普查腿自己在 `20-shapes.md` §3.3 具名认过同一件事：它拿"盯到的那枚目录"判存在性，而真正漏的是被杀时正在跑的那一发）⇒ 本腿⛔ 沿用它的判语、只取它的**动作行**（同一进程内盯＋下手＋`target_pid_still_present=0`），漏与不漏由外面那把 `comm -13` 尺定 | | |
| a4 | 外力杀（本腿改形的台件，见 §3） | **57→58** | **`TestRoundTripLoadMarshalLoad826293033`** | `c1_noredirect world=CONSTRUCTED`／`guard_name_is_test_exe=True guard_cmdline_has_workdir=True guard_parent_is_this_leg=True`／`external_kill_issued at=…`／`after_kill target_pid_still_present=0 … leaked_count=1 leaked_names=[TestRoundTripLoadMarshalLoad826293033]`／`os_temp_test_star_now=58` | `c1_rc=0`（`verdict=DIR-LEAKED-in-watch-root shape_c_reproduced`） | `logs/10c-c1.txt`＋`logs/10c-roster-start.txt`／`-end.txt`／`10c-topnew.txt`（差集里**只有**这一枚名字） |

⇒ **(a) 侧＝死亡那一形当场让 `%TEMP%` 顶层多出一枚，两形（超时 panic／外力杀）各一枚，逐枚具名、同行证据齐。**

## 2. (b) 还原（＝本腿那一步在位）⇒ 顶层不多，残渣落在专用根里（逐字路径）

| # | 发 | 根（逐字） | 名册 before→after | 顶层新增 | 根里残渣（逐字） | 同行证据 | rc | 件 |
|---|---|---|---|---|---|---|---|---|
| b1 | 超时 panic，第一稿的根（仓所在卷 `D:`） | `D:/tmp/wisp-test-gotmp/probe-114728-128` | 56→56 | **0 枚** | `D:\tmp\wisp-test-gotmp\probe-114728-128\TestRoundTripLoadMarshalLoad2086954487` | `panic: test timed out after 2s`（`logs/10b-shot.txt:1669`）＋`run_lines=417`／`pass_lines=416` | `shot_rc=1` | `logs/10b-shot.txt`＋`10b-root-entries.txt`＋`10b-new.txt`（空） |
| b2 | 超时 panic，**出厂形**（`portable-tests.sh` 自己打出来的那枚根） | `C:/Users/swq/AppData/Local/wisp-test-gotmp/portable-20261012-115401-2013` | 首试：58→58、`run_lines=244`／`pass_lines=244`、根里 `root_entries=0` ⇒ 帽落在两发迭代之间、**⛔ 一个正在跑的 `t.TempDir()`** ⇒ **〔作废〕**（⛔ 拿它的 0 当凭据） | — | — | `logs/10b2-shot.txt` | `shot_rc=1` | `logs/10b2-*` |
| b3 | 同一形重打（出厂根，正发） | 同 b2 那枚根 | **58→58** | **0 枚**（`logs/10b2-r1-new.txt` 为空，`empty above＝0 枚`） | `C:/Users/swq/AppData/Local/wisp-test-gotmp/portable-20261012-115401-2013/TestRoundTripLoadMarshalLoad1761995879` | `attempt=1 rc=1 run=240 pass=239 panic=1 top_new=0 root_entries=1` ⇒ **第 240 遍正在跑时被砍、那一遍的目录留在根里、顶层一枚没多** | `rc=1` | `logs/10b2-r1-shot.txt`／`-before`／`-after`／`-new.txt`／`-root.txt` |
| b4 | 外力杀（形 (c) 的还原侧） | `D:/tmp/wisp-test-gotmp/shapec-115000` | 顶层 58→**58**（c1 之后、c2 之后同一枚数，`logs/10c-roster-end.txt`＝58） | **0 枚**（两轮合起来差集只含 a4 那一枚，`logs/10c-topnew.txt`） | `TestRoundTripLoadMarshalLoad489708636`（`logs/10c-rootentries.txt` 逐字） | `c2_redirected world=CONSTRUCTED`／`after_kill target_pid_still_present=0 dirs_before=0 dirs_after=1 leaked_count=1`／`os_temp_test_star_now=58` | `c2_rc=0` | `logs/10c-c2.txt` |

⇒ **(b) 侧⛔ 只是"顶层 0 枚"：同一批发同时给出 ①`panic: test timed out`／`target_pid_still_present=0` 那种"那一形真发生了"的行，②残渣在专用根里的逐字路径。**两发正控齐，缺一的本腿自己标了〔作废〕（a1、b2 两发）。

## 3. 台件（⛔ 重造，但⛔ 改别家的跟踪件）

- 复用：`probes/310/a1/logs/20-shape-c-watch.ps1`（普查腿交付件，先 `wc -l`＝63、`head` 摸口径后**原样跑**一发＝a3；它的自报判语按 §1 那一条具名不采信）。
- 派单点名的 `probes/orch/20261011/shape-c-orch.ps1`（38 行）⛔ 适用本票：它的"世界被造出来"判据＝"盯到一枚命令行含本仓 `internal/llm/testdata/golden` 的 `mockllm.exe` 孩子"（票 304 的进程面形状），而 `internal/config` 那一枚用例⛔ 起子进程 ⇒ 拿它跑只会读回 `world=NOT-CONSTRUCTED`。具名报回＝`Q-310-r1-1`。
- 本腿自家派生件＝`logs/10c-shapec-watch.ps1`：同一套"同进程起＋盯＋`Stop-Process`"动作，两处改动**写在件头注释里**（① 盯的根变成参数 `$WatchRoot`，② 漏不漏用"起进程之前的名册快照 ⇄ 杀完之后快照"的**逐名差集**判，⛔ 盯"碰巧先看到的那一枚"）。件⛔ 改别家那份。
- ★"盯＋下手必须在同一个进程里"这条按派单口径执行（bash↔PS 往返 0.5–1 s、短命子进程只活 1–2 s），a3/a4/b4 三发的 `external_kill_issued` 与 `world=CONSTRUCTED` 都出自同一枚 PowerShell 进程（`guard_parent_is_this_leg=True`）。

## 4. ★本腿这一程造的目录（归我、可证＝本仓，⛔ 删、⛔ 让任何人顺手删）

`%TEMP%` 顶层：起手 **55** 枚 → 交件 **58** 枚，多出来的三枚**逐枚都是本腿的**（名字首段＝本仓用例名）：

| # | 目录 | 哪一发 | 位置 |
|---|---|---|---|
| 1 | `TestRoundTripLoadMarshalLoad2461500944` | a2（超时 panic，撤掉那一步） | `%TEMP%` 顶层 |
| 2 | `TestRoundTripLoadMarshalLoad2236627062` | a3（外力杀，撤掉那一步） | `%TEMP%` 顶层 |
| 3 | `TestRoundTripLoadMarshalLoad826293033` | a4（外力杀，撤掉那一步，改名册尺那一发） | `%TEMP%` 顶层 |
| 4 | `TestRoundTripLoadMarshalLoad2086954487` | b1（还原，第一稿的卷） | `D:\tmp\wisp-test-gotmp\probe-114728-128\` |
| 5 | `TestRoundTripLoadMarshalLoad1761995879` | b3（还原，出厂根） | `C:\Users\swq\AppData\Local\wisp-test-gotmp\portable-20261012-115401-2013\` |
| 6 | `TestRoundTripLoadMarshalLoad489708636` | b4（外力杀，还原侧） | `D:\tmp\wisp-test-gotmp\shapec-115000\` |

另：本腿建过、当时为空的专用根若干枚（`…/portable-20261012-114637-408`、`…/winsec-20261012-115046-1584`、`…/winsec-20261012-115521-883`、`…/cli-20261012-115545-1029`、`C:\Users\swq\AppData\Local\wisp-test-gotmp\probe-115234`）与两枚仓外导出树（`/d/tmp/310r1-tree.Ttzxdu`、`/d/tmp/310r1-tree2.dieduj`，后者多塞了三枚 `third_party/sherpa-onnx/*.dll` 副本）、一枚编译产物 `D:\tmp\310r1-shapec\config-310r1.test.exe`。全部**只建不删**，路径逐字留在本件；在册那 55 枚存量（含 `TestProbeStepByStep3843095677`）本腿一枚都没碰。
