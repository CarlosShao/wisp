# 322-r2 头号活读数：`cmd/wisp` 那一发（五档，全部本腿现跑，⛔ 转引）

## 0. 起跑前置（编排者裁定①的注意①＝锚漂移复跑）

- 起跑时刻 `11:08:05 +08`，盘上 `HEAD=1b9f602c`（我起树之后又落了 336-r1／166-r1／337-a1／228-r3c／编排者记账若干笔）。
- 漂移尺（`cmd-wisp.ps1` 起跑前那一版先跑）：`cmd/wisp` 树对象
  `c11163c8` ＝ 现 `HEAD` ＝ **`54efbbbf2017…`** 同一枚；
  我这棵导出树 `cmd/wisp` 与现 HEAD 逐枚比 ＝ **139 枚文件、`Compare-Object` 差 0 行**
  （尺＝`git ls-tree -r --format='%(objectname) %(path)' HEAD -- cmd/wisp` ⇄ 树内逐枚 `git hash-object`）
  ⇒ 旧树仍等于现锚的 `cmd/wisp` 字节，读数有效，⛔ 重新认锚。原件＝`logs/cmdwisp-drift-check.txt`。
- 台面＝仓外导出树 `D:\tmp\wisp\20261018\322r2\head.RJNBgS`（⛔ 共享工作树），
  `PATH` 前置 `third_party/sherpa-onnx` ＋ `build` 两串 dll（树内各 3 枚 dll 现量在位），
  `-count=1 -timeout 30m -v`；`GOTMPDIR`＝`gtD`（D 台）⇄ 未设＝默认 `TMP`（C 台）。
- 四数口径具名：`RUN=` 数 `^=== RUN`（**含子测试行**）、`PASS_TOP/FAIL_TOP/SKIP` 只数顶格 `^--- …`，**⛔ 混用⛔ 相减**。

## 1. 整包那一发（默认档，两台成对；派单裁定②⒜＝只作台面终态记录）

```
cwisp-full-D  rc=1  RUN=448 PASS_TOP=322 FAIL_TOP=8 SKIP=3   483.689s  GOTMPDIR=gtD(D)
cwisp-full-C  rc=1  RUN=448 PASS_TOP=323 FAIL_TOP=7 SKIP=3   445.592s  GOTMPDIR=默认(C 卷)
两台 package-line 均为  FAIL  github.com/CarlosShao/wisp/cmd/wisp
台面故障尺：0xc0000135=0  build failed=0  no tests to run=0  setup failed=0（两台都是）
```

**这一回两台都跑到终点、没有整进程死亡**（对照 228-r3b 具名那一格）。相抵尺（两台同一形）：
`=== RUN` 行 **448** ＝ 顶格 **333** ＋ 子测试（名字含 `/`）**115**；
判定行 ＝ 顶格 `^--- …` **333** ＋ 缩进 `    ^--- …` **115** ＝ **448**，
且顶格 RUN 名去重 333 ⇄ 顶格判定名去重 333、**未判定 0 枚** ⇒ 两把尺逐枚相抵，
**零枚失踪＝不是死亡、不是漏数**（尺见 §5；⛔ 拿 `RUN=448 − PASS_TOP=322` 那种相减当"失踪 126 枚"）。

两台差集（红名册）＝**只有一枚**：`TestAC2SealNoticeLandsInTheRunLegLogFile` 只在 D 台红，
C 台绿 ⇒ 随台面摆动的记录数那一格（红句 `logsink_windows_test.go:175: records naming "winsec: seal cleared
principals that stood on this object" = 2, want 1`）。SKIP 名册两台**逐字相同** 3 枚。

## 2. estate 那一发（裁定②⒜＝AC#4 以此为准）

```
cwisp-estate-resident  rc=0  RUN=31 PASS_TOP=22 FAIL_TOP=0 SKIP=0   36.692s   -run 'Resident'
cwisp-estate-live      rc=0  RUN=0  PASS_TOP=0  FAIL_TOP=0 SKIP=0   0.093s    -run 'TestLive'  ← ⛔ 不是读数
cwisp-estate-deathpoints rc=0 RUN=2  PASS_TOP=2  FAIL_TOP=0 SKIP=0   2.338s    -run 'TestPanelThreadIsSTAAndExitsCleanly|TestAC13BringUpSurvivesAReusedThreadQuit'
```

- **estate-resident 那一发的口径先算清**（⛔ 拿 RUN⇄PASS 的差当失踪）：`RUN=31` ＝ 顶格 22 ＋ 子测试 9，
  顶格 22 枚判定名与 RUN 名**两把尺相抵＝未判定 0 枚**，`PASS_TOP=22 FAIL_TOP=0 SKIP=0` ⇒ rc=0。
  ⚠ 派单/判语里那句 "`resident_*` 在整包从不执行" **在本机今天这一发不成立**：整包 D 台里顶格
  `=== RUN …Resident` 命中 **32 行**、判定行 **22 行**、其中 FAIL **0** ⇒ 这一族在整包里既被跑到也全绿，
  与定向那一发的 22 枚同数（具名交编排者：那一格是 228-r3b 台面上的形态，⛔ 本腿复认、⛔ 据此判任何格）。
  ★差额具名（28 grep ⇄ 22 实跑）＝未跑 6 枚**全部**住在 `//go:build windows && winlive` 那三件
  （`resident_ball_live_228`×1／`resident_task_source_live_246`×1／`resident_hotkey_live_258`×4），
  默认档不编译 ⇒ 两把尺对得上：**28 ＝ 22（默认档实跑）＋ 6（winlive 档）**。
  ⚠ 具名口径：这 6 枚与 §2 那六枚 `winlive` 红**不是同一集合**（只共享同一枚 build tag；
  交集＝`TestLive228ResidentLegOwnsABallWindowOnTheDesktop` 与 `TestLive258ResidentOccupiedCombinationNamesTheNewValue` 两枚），
  ⛔ 把两处六枚当同一处数读。`-tags winlive` 那一发顶格 12 枚＝6 绿＋6 红，与 `RUN=12` 相抵。
- **estate-live 第一版那一发＝台面故障、⛔ 计入判据**：`ok … [no tests to run]`＋`noteststorun=2`。
  根因现读＝`TestLive*` 六件文件首行全是 `//go:build windows && winlive`
  （`panel_transport_live_35v2`／`resident_approval_live_246`／`resident_ball_live_228`／
  `resident_hotkey_live_258`／`resident_task_source_live_246`／`resident_tray_stop_live_228`），
  我第一版台件⛔ 带 `-tags winlive` ⇒ 过滤器选不到任何用例＝零读数，不是全绿。
- 补正那一发（同树同台面，只多 `--tags winlive`）：

```
cwisp-estate-live-winelive  rc=1  RUN=12 PASS_TOP=6 FAIL_TOP=6 SKIP=0  93.399s
  --- FAIL: TestLive246ConfirmingCardBorrowsEscVetoesAndReturns (4.62s)
  --- FAIL: TestLive246ExitRefusesAHangingL2Card              (0.10s)
  --- FAIL: TestLive246ExitAbandonsAHangingL1Window           (0.09s)
  --- FAIL: TestLive228ResidentLegOwnsABallWindowOnTheDesktop (7.13s)
  --- FAIL: TestLive258ResidentOccupiedCombinationNamesTheNewValue (5.34s)
  --- FAIL: TestLive228R3TrayExitStopsTheShippedProcessWithCodeZero (14.59s)
```

  红句逐字（三枚 246 同一句）：`resident_approval_live_246_windows_test.go:108: idle live hot key set = 0
  ([{Name:summon ID:1 Binding:Ctrl+Alt+Q Status:occupied by another app Acc:{Mods:16387 VK:81}
  Err:Hot key is already registered. Note:} …` ＝**本机有常驻实例占着热键**（尺＝`Get-Process wisp`
  现量 PID 21892、`StartTime 2026/10/15 19:03:14`，跑在我这一发之前、非本腿所起）。
  `228-resident-ball-v1.md:151-153` 把 "hotkeys live 0/4" 写成一枚**会红的判据**（＋`:121` 那一枚在无桌面时报红），
  台账 `:17005`／`:18192` 已把这一族具名成"**本机既有真窗／真麦／驻留族**…⛔ 由本票归因"。
  ⇒ 本腿⛔ 判生死、⛔ 为它动字节，只交数＋归位；`-tags winlive` 在 CI 侧的射程按 228-v1 §F-6 仍是零读数。

## 3. 整包八枚红的逐名归位（必答 (i) 的料）

| # | 用例 | 本腿现跑红句要点 | 归位（凭据＝盘上，⛔ 我裁） |
|---|---|---|---|
| 1 | `Test258BridgeMutationNoSrcKeepsOldBinding` | 热键族 | **别票存量**＝编排者 10-16 那发名册里就有（`probes/orch/20261016/logs/orch-322v1-ac4-cmdwisp-full.txt`，`430/305/6/3`） |
| 2 | `Test258OccupiedCombinationNamesTheNewValue` | 同上 | 同上（那一发六枚之一） |
| 3 | `Test258V1ProbeSummonEditRebindsLiveBall` | 同上 | 同上 |
| 4 | `TestAC14GoSideEvalPushReachesThePage` | `panel_resident_windows_test.go:946: Go's Eval push did not reach the document … title as ""` | **别票存量**＝同一枚 10-16 在册 ＋ `228/r3b/logs/surface-c.log` 同名 |
| 5 | `TestAC246DevLegIgnoresTheTestTaskInjection` | 驻留实例族 | **别票存量**＝10-16 在册 ＋ CI `run-37166458550-failed.log` |
| 6 | `TestAC4FocusReturnToPriorWindowGap33r5` | `panel_host_windows_test.go:948: the panel did not take the foreground on Show` | **别票存量**＝10-16 在册 ＋ `228/r3b` 的 c／d 两面同名 |
| 7 | `TestCleanCheckoutBuilds_AC11` | `panel_host_gate_test.go:404: git rev-parse --show-toplevel: exit status 128` | **台面/装置红**＝导出树没有 `.git`；台账 `:15164` 已把这一枚定式化为"4 枚真窗族＋**1 枚装置红**"（`A802` 第二面：⛔ 只缺 DLL，还缺 git 元数据）⇒ 派单"凡红句含这些形状判台面故障"那一支的具体化身 |
| 8 | `TestAC2SealNoticeLandsInTheRunLegLogFile` | `logsink_windows_test.go:175: records … = 2, want 1` | **随台面摆动**＝两台差集里唯一一枚（D 红 C 绿），且 10-04 CI 失败件与 `228/r1/baseline-test.txt` 都红过＝别腿基线既有形状 |

⇒ **本票新增红＝零枚**（第 7 枚属我这棵导出树的台面形状、⛔ 是产品字节；第 8 枚属摆动）。
对照尺：编排者 10-16 那发的六枚名册 ⊂ 我今天 D 台八枚（`in-orch-not-D` 差集＝**空**）。

## 4. 本腿那一枚 `.go` 对 `cmd/wisp` 的因果尺（机制级，⛔ 靠推理）

- `sealhandle_privateroute_120_other_test.go` 首行＝`//go:build !windows`；
  `go list -e -f IgnoredGoFiles`（GOOS=windows）命中该件 **1 枚**、同档 `TestGoFiles` 命中 **0 枚**；
  GOOS=linux 时进 `TestGoFiles`（命中 1）。⇒ 本机 `cmd/wisp` 这一发**根本不含这枚文件**。
- 就算含，`_test.go` 只编进**自己包**的测试二进制：`go list -e -f TestGoFiles ./cmd/wisp/` ＝ **89 枚**，
  其中来自 `internal/winsec` 的＝**0 枚**。⇒ 我这枚字节**机制上无法**改变 `cmd/wisp` 的颜色。
- 生产字节尺＝`cmd/wisp` 树对象在起树锚与现 HEAD 同一枚 `54efbbbf`（§0）。

## 5. 死/失踪尺（我这版自己写的，两枚台各自跑过）

`=== RUN` 行 448 ＝ 顶格 333 ＋ 子测试 115 ⇄ 判定行 448 ＝ 顶格 333 ＋ 缩进 115
（顶格 RUN 名去重 333 ⇄ 顶格判定名去重 333 ⇒ **未判定 0 枚**），三把尺互相咬合，零枚静默死亡。
⚠ 具名口径纠正：派单色尺的 `RUN=` 是**含子测试**的那一把（448），`PASS_TOP/FAIL_TOP/SKIP` 只数顶格
（322＋8＋3＝333＝顶格判定行枚数）⇒ 顶格判定 333 ＝ PASS_TOP 322 ＋ FAIL_TOP 8 ＋ SKIP 3，
这一枚相抵是**顶格内部**的账，⛔ 与 448 那把尺相减。
`dead=` 那把尺命中 2 行逐字是装置自陈（`panel_host_gate_test.go:404 … exit status 128`、
`slo_exit_os_156_windows_test.go:206 Wait on a killed child reports the kill as an error (exit status 1);
that is the fixture, not a failure`），⛔ 读成整进程死亡。

## 6. 本腿★同一笔形状违规第二次再犯（`9c372800`，具名自抓，⛔ 回改历史）

`9c372800` 越界 **2 枚**＝`.scratch/wisp/probes/orch/20261018/txt/{commit-40em.sh,commit-corrections.sh}`
（编排者已 `git add` 未 commit 的两枚新件）。
机制与 `7f805583` **完全同一形**：我的 `git add` 逐枚字面 pathspec 合规（18 枚全在我写面内），
但 `git commit` 我写成 `-- $s`，那个 `$s` ＝ `git diff --cached --name-only` ＝**共享 index 的全集**
⇒ 把别人 staged 的东西一起提走。**我在同一笔的正文里刚写着"pathspec 两处仍字面逐枚写死"，
那一句话当场就是假的**＝不是知识缺口，是我没把尺用在最后一步上。

盘上事实（三把尺）：两枚件在 `9c372800^` **不存在**（`git log --diff-filter=A` 显示由我这一笔首现）；
我笔 blob ＝ 工作树 `git hash-object` 现量逐字相同（`b29ca514…`／`1d096171…`）
⇒ **内容没丢、没改坏，只是代提在 `322-r2` 名下**；我笔之后共享 index 剩 0 枚。
定式（升级，⛔ 再靠"我记得"）：**`git commit` 的 pathspec 必须逐枚手打完，与 `git add` 那一串同形，
两处的字面路径都不许来自任何 `git diff --cached`／`git status` 输出**——`$s` 那形一次都不许用。
上报编排者：这两枚请照 `7f805583` 那次的处理＝记账"由 `9c372800` 代提、字节同原稿"，
编排者下一笔⛔ 重复 `add`、⛔ 判"丢了件"。

## 7. 前两版台件的故障（⛔ 计入判据，件与正文只建不删）

`cwisp-*.{FAULT-recursion.log}` ＋ `cmdwisp-summary.FAULT-recursion.txt`：五发 `RUN=0`、2 秒"跑完"。
根因＝函数名 `Go` 撞上外部命令 `go`（PowerShell 函数名大小写不敏感）⇒ 函数体那行递归调自己、
没碰工具链；`0xc0000135=2` 是嵌套 summary 被自己数到的**自指产物**（两串 dll 现量都在树里）。
对照凭据＝同目录 `host-winsec.ps1` 用不撞名的 `HostRun`，那一族五发都出了真数（`RUN=110` 等）。
本版改名 `Run-Shot` ＋ 外部命令写 `& 'go'` ＋ 语法自检 `ParseFile` 零错。
