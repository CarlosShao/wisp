# 票 136 AC#10 只读量腿（136-r1b）读数表

- 角色：只读腿。未落任何 `.go` 字节、未改任何断言、未碰票面/台账/CI、未造钉子。
- 锚点（本程自量）：`1d5c0706ee24f8336683e49782c3bbfa5c00e4f3`，`go version go1.27.1 windows/amd64`。
- 工作树对 `cmd/wisp/`、`internal/observe/`、`scripts/`、`.github/` 零漂移（`git status --porcelain -- <四树>` 输出空）
  ⇒ 本程 `go test` 读到的就是 HEAD 源码，不是别人未入库的活。整树 `??` 1237 条、` M` 16 条全在临时件，
  **被改的 `.go` 追踪文件数 = 0**。
- 口令遵守：`unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`、`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`
  （Bash 形态直接前置，未用 `cygpath -wp` 转；未建 worktree、未用 `-overlay`、未动 `.gitattributes`）。

## 0. 一句话结论（正文见 1/2/3 节）

**这枚接缝今天「不响」，且不是「量不到」——是「没有那枚用例」**：`run.Pass` 那条链在生产码里真实存在
（HEAD `cmd/wisp/slo_windows.go:399` 与 `:336-339`），但 `cmd/wisp` 包内 **0 枚用例读它、0 枚用例调 `cmdSLO`**。
⇒ AC#10 仍是未实现的未勾格；本程**没有造钉子、没有写用例、没有做变异**。

## 1. 前提新鲜度（先验，逐枚）

### 1.1 落件在不在（`git cat-file -e HEAD:<path>`，13 枚逐枚，IN = 命中）

`cmd/wisp/slo_windows.go` IN · `cmd/wisp/slo_other.go` IN · `cmd/wisp/main.go` IN ·
`cmd/wisp/slo_report_144_windows_test.go` IN · `cmd/wisp/slo_exit_os_156_windows_test.go` IN ·
`internal/observe/sampler.go` IN · `internal/observe/sampler_test.go` IN ·
`internal/observe/sampler_settle_coverage_136_test.go` IN · `internal/observe/sampler_zerosample_136_test.go` IN ·
`internal/observe/goroutine.go` IN · `scripts/slo-check.ps1` IN · `.github/workflows/ci.yml` IN ·
`docs/evidence/s1/134-adversarial-acceptance.md` IN。**0 枚 OUT。**

### 1.2 那句「可判定的断言」对 HEAD blob 复量（未对工作树、未凭票面自述）

票面 `:167` 原句：`grep -rn "run\.Pass" cmd/wisp/*_test.go` ⇒ **0 命中**。复量：

- `git grep -n "run\.Pass" HEAD -- "cmd/wisp/*_test.go"` ⇒ **rc=1，0 行** ⇒ **断言成立（新鲜）**。
- 放宽到 `\.Pass`：`cmd/wisp/*_test.go` 只 3 行，全在 `slo_report_144_windows_test.go:125/:209/:210`，
  读的是测试自己写的 fixture（`slo144Report` 里 `Pass: true`）经 `collectReportWithin` 解析出的 `sloRun.Pass`。
  ⇒ 它**不**经过 `rep.Pass && observer.Pass -> run.Pass -> exitCode` 这条链，**不是**这条接缝的见证者。
- 生产侧实名：`run.Pass` 在 HEAD 出现 4 处——`:300`（settle 支）、`:336`（`if !run.Pass`）、
  `:399`（`rep.Pass && observer.Pass`）、`:443`（subject 支）。

### 1.3 坐标腐烂（如实登记：断言成立，但行号大面积过期）

本程按内容重新定位，未照抄票面行号。票面 -> HEAD 实名（漂 = 行号过期、内容仍在；过期 = 所指之物已不是票面描述的那样）：

| 票面说 | HEAD 实际 | 判 |
|---|---|---|
| `:374` 采样 | `:385` `observe.NewSampler(proc.NewExternalSampler(...))` | 漂 |
| `:386` `run.Pass = rep.Pass && observer.Pass` | `:399` | 漂 |
| `:323-325` `!run.Pass` 则 `exitCode=1` | `:336-339` | 漂 |
| `:270` `NewTreeSampler` 注入点 | `:283` | 漂 |
| `:223` `if !*settle && *state == "" {` | `:236-238`（`-state is required` 在 `:237`） | 漂 |
| `:623` 合成的失败报告 | 该文件 `:620-628` 是 `subjectReportRead.summary()`；**全文件 `back_within_cap_ms`/`cap_bytes`/`final_bytes` 0 命中**，唯一 `SettleReport` 字面量在 `:988` | **过期** |
| `:477` 那段初始化 | `cmd/wisp` 内无此 shape；`BackWithinCapMS: -1` 的初始化在 `internal/observe/sampler.go:489` | **过期** |
| `slo-check.ps1:364-370` leak 自检 | `:581-597`（`$leakFlipped = ($leakCode -eq 1)` 在 `:594`） | 漂 |
| `slo-check.ps1:153-155` 争用名单 / `ci.yml:479-500` `slo-smoke` | `ps1:361` 起（`$loadNames`） / `ci:836-837`（`runs-on: windows-latest`） | 漂 |
| `134-adversarial-acceptance.md:538` 那条命令 / `:643-645` 「这条路不通」 | `:538` / `:643` 在 `7.2` 第 1 条 | 不漂 |

判：**实质前提全部成立**（接缝在、无读者、`slo-smoke` 是 CI 层第二 catcher、零样本在 CLI 上不可达），
**但票面 `:623` 与 `:477` 两处坐标所指之物在 HEAD 上已不是票面描述的那样**。这两处属 AC#10「追加两问」，
本程按「未定义即停」处理：**只登记、不据此推导任何 AC#10 追加问的答案**。

### 1.4 本程顺带核到、票面未点名的一条

`cmdSLO` 的 reader 注入面：`grep 'readReportFile|treeReader|inject' cmd/wisp/slo_windows.go` 仅
`:135`/`:142`/`:795` 三处，都是 `sloSubject` 的报告**文件**读取钩子，**没有**任何把「树读数」注进 `cmdSLO` 的面
⇒ 票面 `:170` 那句「`cmdSLO` 没有任何注入面」**成立**。
## 2. 现量：这条端到端接缝今天响不响

### 2.0 先答一个前置事实（这条决定后面怎么读）

**AC#10 全格（票面 `:149-176`）没有点名任何一枚 Go 用例名**：`sed -n '149,176p' | grep -E 'Test[A-Za-z0-9_]+'`
⇒ **rc=1，0 命中**。AC#10 是一格**待写**的 AC，不是一枚已存在的用例。
⇒ 「`-run '<AC#10 点名的用例名>'`」**无值可填**；本程改用 AC#10 改写块 `:166` 自己点名的那把尺
（`go test ./cmd/wisp/`）＋最贴族 `^TestSLO` 作定点读数，并各带正/负控（见 2.2）。

### 2.1 四数（原样，未相加、未相减；每格独立）

尺 A：`go test -run "^TestSLO" -count=1 -v ./cmd/wisp/`，`rc=0`

| 口径（正则原文） | 枚数 |
|---|---|
| `^=== RUN` | **26** |
| `^--- PASS` | **19** |
| `^--- FAIL` | **0** |
| `^--- SKIP` | **0** |

闭合说明（不合并成新数）：另有缩进行 `^[[:space:]]+--- PASS` = **7**（全是
`TestSLO144ReportsThatContradictThemselvesAreCorruptNow/` 的 7 枚子测试）。26 枚 RUN = 19 枚顶层（各带顶层 PASS）
+ 7 枚子测试（各带缩进 PASS）。本程**不**把 19 与 7 相加当「通过总数」报。

尺 B（AC#10 改写块 `:166` 真正点名的那把尺）：`go test -count=1 -v ./cmd/wisp/`，`rc=1`

| 口径 | 枚数 |
|---|---|
| `^=== RUN` | **344** |
| `^--- PASS` | **260** |
| `^--- FAIL` | **6** |
| `^--- SKIP` | **3** |

闭合说明：缩进行 `^[[:space:]]+--- PASS` = **74**。344 = 260+6+3（顶层 269）+ 74（子测试）+ **1 枚被 10 分钟
超时 panic 吃掉、无判决行**（`comm -23` 逐名比对：`TestTicket101SessionGrantDoesNotCrossRestart`）。
**注意：这把尺这次不干净**——`panic: test timed out after 10m0s`（panic 时正在跑
`TestTicket101SessionGrantDoesNotCrossRestart (1m6s)`）⇒ **B 尺的 FAIL=6 只覆盖跑完的 343 枚，不能当整包结论**。
A 尺不受影响（1.369s 跑完、rc=0）。### 2.2 正控／负控（证明 `-run` 不是空转、尺不瞎）

| 发 | 命令（`-run` 原文） | RUN | PASS | 判 |
|---|---|---|---|---|
| 正控（已知存在） | `^TestSLO144CorruptReportIsJudgedOnTheFirstRead$` | **1** | **1** | 尺**不瞎**，命中即跑 |
| 负控（确不存在） | `^TestNoSuchCaseExistsIn136r1b$` | **0** | **0** | 选择器有效，0 命中 = 真没命中 |
| 跨包探（`TestCheckSettle` 族） | `^TestCheckSettle` | **0** | **0** | 该族 8 枚在 `internal/observe`，**不在** `cmd/wisp` |

负控末行原样带 `[no tests to run]` ⇒ 0 命中是可区分的「没有」，不是「跑挂了」。

### 2.3 红名册（逐名，不只报枚数）

B 尺那 6 枚，逐名＋耗时＋文件（**无一枚是 SLO／settle／`run.Pass` 家族**）：

| 红名 | 耗时 | 文件 |
|---|---|---|
| `TestAC14GoSideEvalPushReachesThePage` | 0.56s | `cmd/wisp/panel_resident_windows_test.go` |
| `Test258BridgeMutationNoSrcKeepsOldBinding` | 0.50s | `cmd/wisp/resident_hotkey_258_windows_test.go` |
| `Test258OccupiedCombinationNamesTheNewValue` | 0.17s | `cmd/wisp/resident_hotkey_258_windows_test.go` |
| `Test258V1ProbeSummonEditRebindsLiveBall` | 0.12s | `cmd/wisp/resident_hotkey_v1probe_test.go` |
| `TestAC1ResidentLegInstallsItsLogListenerOnDisk` | 11.19s | `cmd/wisp/resident_sink_nail_127_windows_test.go` |
| `TestAC246DevLegIgnoresTheTestTaskInjection` | 38.13s | `cmd/wisp/resident_task_source_246_windows_test.go` |

SKIP 3 枚（供对照，非红）：`TestAC247LiveMicrophoneLevelsReachTheBallSeam`、
`Test296BridgeKeepsThreeLiveKeysWithOwnerEmptySlots`、`Test296BridgeKeepsThreeLiveKeysWhenConfigUnreadable`。
**A 尺红名册 = 0 枚。**

### 2.4 这条接缝今天到底响不响（与「我读不到」分开写）

- **生产码在**：`:399` 与 `:336-339` 两跳逐字在 HEAD 上（见 1.2）⇒ **不能写「它不存在」**。
- **见证者在不在**：`git grep -nE "cmdSLO\(|runOutOfTree\(|runSubject\(|runSettle\(" HEAD -- "cmd/wisp/*_test.go"`
  ⇒ **rc=1，0 命中**；再查「有没有用例把 `wisp slo` 当子进程跑」
  `git grep -nE "\"slo\",[[:space:]]*\"-state\"|...\"-settle\"|...\"-subject\"" HEAD -- "cmd/wisp/*_test.go"`
  ⇒ **rc=1，0 命中**。⇒ **0 枚用例能观察到 `:399`/`:336-339` 的改动**。
- **诚实边界**：「拆掉它也不会红」这一步本程**没有用变异去读**（只读腿，不落 `.go`）。本程能读的是
  「读 `run.Pass` 的用例 = 0 枚、调 `cmdSLO` 的用例 = 0 枚」这个结构事实，即**无可观察者**。
  这与「我读不到」是两件事，也与「生产码不存在」是两件事，三者本程分开登记。

## 3. 命令末行逐字（每发一条）

1. `git rev-parse HEAD` -> `1d5c0706ee24f8336683e49782c3bbfa5c00e4f3`
2. `go version` -> `go version go1.27.1 windows/amd64`
3. A 尺 -> `ok  	github.com/CarlosShao/wisp/cmd/wisp	1.369s`
4. 正控 -> `ok  	github.com/CarlosShao/wisp/cmd/wisp	0.270s`
5. 负控 -> `ok  	github.com/CarlosShao/wisp/cmd/wisp	0.243s [no tests to run]`
6. 跨包探 `-run "^TestCheckSettle"` -> `ok  	github.com/CarlosShao/wisp/cmd/wisp	0.131s [no tests to run]`
7. B 尺 -> 末行 `FAIL`；其上一行 `FAIL	github.com/CarlosShao/wisp/cmd/wisp	600.267s`

原始日志（只建不删，未用 `.out` 后缀）：
`D:\tmp\wisp\20261021\136r1b\logs\{run-slo,poscontrol,negcontrol,fullpkg,regex-checkSettle}.txt`
＋ HEAD blob 快照 `D:\tmp\wisp\20261021\136r1b\blobs\**`。

## 4. 结论（一句）与交件核验

**这枚接缝今天不响——不是「不存在」也不是「我读不到」：生产码两跳在 HEAD 上真实存在，而 `cmd/wisp` 里读
`run.Pass` 的用例 0 枚、调 `cmdSLO` 的用例 0 枚，故它今天无可观察者；本程未造钉子、未写用例、未做变异。**

`.go` 枚数 **0**：见本次提交的 `git show --name-only`（逐枚读数见提交正文）。
