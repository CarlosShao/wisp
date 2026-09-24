# 票 135 AC#8（家族第 7 发 `M-G`）实现方证据 —— 变异自证三态 ＋ 门禁账

**执行方**：`worker-ticket135-ac8-r2`（**接续程**，非 `worker-ticket135-ac8-mg`）
**被验版本**：commit **`8369b24`**（`8369b24e20eb85b5c9afb45198499dc4b9d9d804`，含上一程的 `fa35557`）
**本格判据来源**：票 135 面 AC#8（09-24 00:0x 编排者追加）＋ 票 133 第二任验收方的
`docs/evidence/s1/133-ac2-r2-acceptance.md` §1.11／§1.12／§3.3／§3.4（`R-133-9`＝`M-G` 的出处）
**本程地界**：一枚新建证据文件（本文件）。**生产码零接触、上一程那把尺零改动**（§1.4 字节为证）。
**本程不翻任何勾**（AC#8 由编排者按本表翻）。

## 0. 争用闸门（每一批读数之前跑一次，逐轮原样登记）

### 0.1 为什么本程把这道闸门当第一条规矩

这台机器上 **self-hosted CI runner 与编队同机**：`.github/workflows/ci.yml:538` 的 `slo-full`
那一枚 job 的 `runs-on` 是 `[self-hosted, wisp-slo]`。编排者在 **10:50:36 +08** 推了一次 dev
（推的正是被验那枚 `8369b24`，`gh run list` 的 `headSha=8369b24e20eb85b5c9afb45198499dc4b9d9d804`），
当场起了一枚 `ci` run ⇒ **"编队里没有别的代理"不等于"机器是空的"**。

### 0.2 固定动作（三行）与本程用的台件

派单写死的三行：

```
tasklist | grep -iE 'Runner.Worker.exe|Runner.Listener.exe'
tasklist | grep -iE '^go\.exe|compile\.exe|cgo\.exe'
gh run list --limit 3 --json status,conclusion,name,createdAt
```

本程把它们包成一枚**有上界的轮询**：`/d/tmp/wisp135mg-r2-gate/poll.sh`
（`INTERVAL=45`s × `MAXROUNDS=26` ≈ 19 分钟一枚上界；命中即 `BUSY` 继续等，
全空即 `CLEAR` 退出 0；到点上界仍忙则退出 3 并照实报。
**没有用 `time.Sleep` 糊窗口**——等的是"条件成立"，且带次数与时长上界。
**没有杀过任何 runner 进程、没有改过 `.github/workflows/**`、没有取消过任何 run**：只观察。）

### 0.3 逐轮原始读数

**手工头两轮（本程自己跑的三行原样）**

```
$ tasklist | grep -iE 'Runner.Worker.exe|Runner.Listener.exe'   2026-09-24 10:51:50 +0800
Runner.Listener.exe          50052 Console                    1     86,552 K
Runner.Worker.exe            36068 Console                    1     98,024 K      ← Worker 活着 = 远程 run 在本机落地
$ tasklist | grep -iE '^go\.exe|compile\.exe|cgo\.exe'          （无输出，rc=1）
$ gh run list --limit 3 --json status,conclusion,name,createdAt
[{"conclusion":"","createdAt":"2026-09-24T02:50:36Z","name":"ci","status":"in_progress"},
 {"conclusion":"failure","createdAt":"2026-09-24T02:24:40Z","name":"ci","status":"completed"},
 {"conclusion":"failure","createdAt":"2026-09-24T02:14:55Z","name":"ci","status":"completed"}]

$ tasklist | grep -iE 'Runner.Worker.exe|Runner.Listener.exe'   2026-09-24 10:55:01 +0800
Runner.Listener.exe          50052 Console                    1     86,936 K      ← Worker 已退，Listener 常驻（本机 runner 空转）
$ gh run list --limit 3 --json status,conclusion,name,createdAt
（同上：仍有一枚 in_progress 的 ci）
```

**轮询逐轮**（全文在 `/d/tmp/wisp135mg-r2-gate/gate-rounds.log`，下面按轮摘 VERDICT 行＋该轮时刻；
每轮的三行完整输出都留在那枚日志里）

```
---- round=1 at=2026-09-24 10:56:55 +0800 ----  worker-hits=0 listener-hits=1 toolchain-hits=0 IN_PROGRESS_COUNT=1  VERDICT=BUSY
---- round=1 at=2026-09-24 10:57:10 +0800 ----  worker-hits=0 listener-hits=1 toolchain-hits=0 IN_PROGRESS_COUNT=1  VERDICT=BUSY
```

⚠ 本节此刻**只有这两轮**：本程后面每一批发作之前都按 0.2 那三行复跑一次，
VERDICT 与该轮时刻记在**对应读数那一节**里（§3 起），不在这里预填。
**读到本文件时若 §2…§7 任一节缺失或未结，那一节的读数就是没采**，不是"采了没写"。

### 0.4 判读规则与本程的处置

- `Runner.Listener.exe` **常驻不算命中**：它是等活的进程，只有 `Runner.Worker.exe` 代表"真有 job 在这台机器上跑"。
  这一条本程按派单原文的形状写（派单第 2 节：「有 Worker 就说明远程 run 活着」）。
- **`gh run list` 里有一枚 `in_progress` 的 ci＝命中**，即使该 run 此刻落在 GitHub 托管的 runner 上
  （本程现量过那枚 run 的 job 分布：`gh run view 35948947994 --json jobs` ⇒ `lint`/`slo-full`/
  `lint-frontend`/`test-core`/`slo-smoke` 已 completed、`test-windows` in_progress，
  `runs-on: windows-latest` 即 GitHub 托管）——
  **本程不赌"它不在这台机器上"，只要闸门命中就重采**。
- ⚠ **本程自己犯过一次、并且当场作废的一发**（不藏）：`§3` 的第一发 `R2-B1`
  在 `GATE-CHECK 2026-09-24 11:02:20 worker=1` 之下取的——编排者 11:01:39 又推了一枚 dev
  （`45623e4`），当场起了新一枚 `ci` run 与一枚 `Runner.Worker.exe`（PID 38372，本程 11:02:47／11:03:01 两次
  `tasklist` 都还在）。⇒ **那一发按派单口径不作数**，本程重新按闸门等到 CLEAR 之后再发整批；
  `R2-B1` 的原始日志留在盘上（`/d/tmp/wisp135mg-r2-out/R2-B1.log` 等）当"这一形在未作废之前长什么样"的对照，
  **但本表任何结论都不引用它**。
- ⚠ 一条本程自己造出来的噪声要交代：**一旦本程开始 `go test`，第 2 行的 `^go\.exe|compile\.exe` 必然命中**
  （那是本程自己的编译）。⇒ 闸门只在**每批读数开始之前**跑一次，批次内的命中不属争用。

**结论（§0）**：闸门**确实命中了**——头两轮一枚 `Runner.Worker.exe`、加上一直挂着的 `in_progress` ci run。
本程因此**先等再采**，等待用的是有上界的轮询，逐轮留档。
〔独立复现〕本程自己跑的三行与 `poll.sh`；日志 `/d/tmp/wisp135mg-r2-gate/gate-rounds.log` 可逐轮重看。

## 1. 锚点、取件、以及"本程没有动那把尺"

### 1.1 锚点（开工第一步自量，不抄派单）

```
$ git rev-parse --short HEAD            8369b24          （开工时；本程第一枚 commit 之后才是 087ef2e）
$ git rev-parse --abbrev-ref HEAD       dev
$ git cat-file -t 8369b24               commit           ← 逐枚现核，见 §8 末
$ git cat-file -t fa35557               commit
$ git merge-base --is-ancestor fa35557 8369b24   真（rc=0）⇒ 被验版本含上一程那把尺
```

被验版本按派单钉在 **`8369b24`**。⚠ 本程跑到一半时 dev 又前进了两次
（编排者 11:0x 提 `45623e4`＝`docs(Q-45 自决,A152)` 并推了，那是他们改 `docs/PLAN.md` 的 commit，
与本票无关、本程未读未改未提交）。**本程全部读数仍在 `8369b24` 的归档树上**，不受影响。

### 1.2 取件（仓外快照，零 worktree／零 checkout／零工作树内编译）

```
$ mkdir -p /d/tmp/wisp135mg-r2-tree && git archive 8369b24 | tar -x -C /d/tmp/wisp135mg-r2-tree
   ARCHIVE-OK
$ sha1sum  /d/tmp/wisp135mg-r2-tree/cmd/wisp/leg_dispatch_gate_133_test.go   23b443ac34787aa9ea60e181b9b8c789f468cd56
$ git show 8369b24:cmd/wisp/leg_dispatch_gate_133_test.go | sha1sum          23b443ac34787aa9ea60e181b9b8c789f468cd56
$ git show 8369b24:cmd/wisp/main.go                     | sha1sum            bba113648726af8d47371a00c17048aa94b1055d
$ sha1sum  /d/tmp/wisp135mg-r2-tree/cmd/wisp/main.go                          同上（逐字同）
$ cp <repo>/third_party/sherpa-onnx/*.dll <tree>/third_party/sherpa-onnx/
     onnxruntime.dll 17799168 / sherpa-onnx-c-api.dll 4605952 / sherpa-onnx-cxx-api.dll 259584
     （三枚 dll 不在归档里＝被 gitignore，缺了测试二进制 0xc0000135，照票 133 验收方 §0.3 的同法补）
```

**没有读脏工作树当被验内容**：工作树里那枚尺与 `8369b24` 的 blob 逐字同（`23b443ac…`，本程
在仓外只读地比过一次），但**所有读数取的都是归档树**，且每发一棵**新**树（`cp -r` 自
`wisp135mg-r2-tree`，树已存在即 `TREE-EXISTS-REFUSING` 拒跑）⇒ 结构上不存在"上一轮的 `.bak`
盖回当前树"那一类作废读数。

### 1.3 旧尺（＝`M-G` 洞的载体）也自取了一份，并核到与前一手同源

```
$ git show fa35557^:cmd/wisp/leg_dispatch_gate_133_test.go > /d/tmp/wisp135mg-r2-oldgate.go
$ sha1sum /d/tmp/wisp135mg-r2-oldgate.go        906201f4a3995d10b0a65910aa4cc68e1f745a9d
$ git show fa35557^:cmd/wisp/leg_dispatch_gate_133_test.go | sha1sum   906201f4a3995d10b0a65910aa4cc68e1f745a9d
```

⇒ 这枚 `906201f4…` 与票 133 第二任验收方 §0.3 记录的
`git show 048a9e4:cmd/wisp/leg_dispatch_gate_133_test.go | sha1sum` **逐字同**。
也就是说：**本程拿来当"未修那一侧"的旧尺，与造出 `p7` 那一发、量到 `100/53/0/0` 全绿的那枚尺，是同一枚字节**
⇒ §4 的旧尺对照发才配当"洞复现"的证据，而不是另一棵树的另一回事。

### 1.4 本程对那把尺做了什么：**零**

- 归档树的尺 `23b443ac…` ＝ 工作树的尺 ＝ `8369b24` 的 blob ⇒ **一字未动**。
- 本程**不判断需要改这把尺**：派单第 1 节那条"若必须改才许结案就停手报回"的停手线**没有被碰到**。
  它引用的证据路径 `docs/evidence/s1/133-ac1-ac2-instrument.md` 本程也现核存在（`ls` 命中）。
- 所有变异都落在**快照副本**里（`/d/tmp/wisp135mg-r2-T*`），生产码与测试码在工作树里零改动。

### 1.5 工作树里那一枚不是本程的改动（登记，不提交）

开工第一刻 `git status --porcelain` 是**空**。本程 §0 提交之后，工作树出现
`M docs/PLAN.md`（4 增 4 删，`D1–D46`→`D1–D47`／`C1–C31`→`C1–C32` 那类指针计数修正），
**不是本程写的**；本程未改它、未提交它，随后它由编排者自己以 `45623e4` 入库。
本程每一枚 commit 都带显式 pathspec（`git add -- docs/evidence/s1/135-ac8-mg-impl.md` ＋
`git commit -q -F - -- <同一路径>`），`git diff --cached --name-only` 每枚都只有本文件。

## 2. 被读的那把尺：形状登记（读数在 §3 起）

`8369b24` 里那枚 `cmd/wisp/leg_dispatch_gate_133_test.go` 的 M-G 判据（本程逐枚 `grep -n` 现量行号）：

```
291    coveredBy []string                                   ← leg133 带上"这行账本 credit 了谁"
1538   const rosterChildEnv135 = "WISP_135_RUN_ROSTER_CHILD" ← 子进程防重入闸
1543   const rosterReadTimeout135 = 2 * time.Minute          ← 上界，不是 sleep
1548   var rosterName135 = regexp.MustCompile(`^Test[A-Za-z0-9_]*$`)
1570   func compiledRunRoster135() (map[string]bool, int, error)
1601   func headRunes135(s string, n int) string
1614   func flagValue135(name string) string
1632   func (p *pkg133) runRosterReds135(...) []string       ← 判据 (c)
1663   func (p *pkg133) runRosterDisclosure135(...) string   ← 替换掉那句 guessed 的失明披露
```

**读数来源是运行中那枚二进制，不是解析器**（这条就是票面禁止"折回只信 AST"的那一条）：

```
1574   exe, err := os.Executable()
1580   cmd := exec.CommandContext(ctx, exe, "-test.list", ".*")
```

全文件里 `//go:build` / `_linux_test.go` / `runtime.GOOS` 只出现在**注释与红名文案**里
（`:462`、`:1560`、`:1569`、`:1651`、`:1912` 的 `goos133()`），**没有任何一处代码去求值构建标签**
⇒ 它没有把第二种"只问签名"换成第三种"只问字节"，而是引入了第三种**独立**读数。本程判：方向合规。

`runRosterReds135` 的两条 fail-closed 也现量到（这两条决定 §3 起怎么读它的红）：
`rerr != nil` ⇒ 直接红（"读不到名册"不是"没人欠覆盖"）；名册里没有本尺自己的名字 ⇒ 也红。

⚠ 本程在这把尺里读到一处**它自己写明的不对称**，登记、不改、不替它辩护：
`covered=test` 那一桶**红**，而 registry 钉（`covered=nail`）落在名册外时**只披露不弄红**
（理由＝那四枚住在 `*_windows_test.go`，在 linux 形上必落名册外，弄红＝在本平台主张一笔覆盖债）。
⇒ 本程的 §5 反向控制与 §7 门禁账都只测 **windows 腿**；linux 腿的运行读数**本程取不到**，
写进 §8 的"没做的档"，不当已验。

〔独立复现〕上面每一行行号与 sha 都是本程自己在 `8369b24` 的归档树／仓外只读命令上量的。
〔日志＋归档，抽验〕"上一程为何这么修"的理由文本来自 `fa35557` 的 commit message 与本文件源码注释，
本程只核了字节与行为，不背书它的措辞。

## 3. 基线（未变异）—— 四数 ＋ 逐名名册 ＋ 差集

台件：`/d/tmp/wisp135mg-r2-run.sh <标签> <树> none none - <go test 参数>`
（每发一棵从 `wisp135mg-r2-tree` 新 `cp` 的树，树名已存在即拒跑；每发先 `go build ./cmd/wisp/`
再 `go build ./...` 各取一次 rc，才读红绿）。
本批四发前的闸门：`prebatch round=4 at=11:11:30 worker=0 toolchain=0 inprog=0 VERDICT=CLEAR`，
逐发 `GATE-CHECK` 全部 `worker=0 toolchain=0`；整包三发在
`round=1 at=11:17:23 … VERDICT=CLEAR` 之后的窗口里。

| 发 | 树 | 命令 | rc | RUN | 顶层 PASS/FAIL/SKIP | 子测试 P/F/S | SKIP 行 | build-failed | panic | 名册枚数 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `R2-BL1` | T02 | `-count=1 -run=^TestAC1AC2DispatchHopGate133$` 新尺 | 0 | 1 | 1/0/0 | 0/0/0 | 0 | 0 | 0 | 1 |
| `R2-BL1o` | T03 | 同上，**旧尺**（`fa35557^`） | 0 | 1 | 1/0/0 | 0/0/0 | 0 | 0 | 0 | 1 |
| `R2-BLF` | T18 | `-count=1`（整包，门开着） | **0** | **101** | **54/0/0** | 47/0/0 | 0 | 0 | 0 | 54 |
| `R2-BLC` | T19 | `-count=1 -skip=^TestAC4EveryLegIsNailedOrRuled$`（门关着） | 0 | 100 | 53/0/0 | 47/0/0 | 0 | 0 | 0 | 53 |
| `R2-RST` | T26 | 还原发，同 `R2-BL1` | 0 | 1 | 1/0/0 | 0/0/0 | 0 | 0 | 0 | 1 |

**`R2-BLF` 的 54 枚逐名名册**（`=== RUN` 去重后逐名，全量列，不缩写；
文件 `/d/tmp/wisp135mg-r2-out/R2-BLF.runnames`）：

```
TestAC1AC2DispatchHopGate133            TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink
TestAC1ResidentLegInstallsItsLogListenerOnDisk  TestAC1ResidentLegOutlivesItsOwnLogFailure
TestAC2AC3DegradedLegsStillDeliverTheirVerdict  TestAC2AuditTrailLandsInTheRunLegLogFile
TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128  TestAC2ModelsLegBooksItsHandOffVerdictOnDisk
TestAC2RealProcessRefusesOnEveryLegWithoutAppData128   TestAC2RefusalMarkersAreNotAShortenableList128
TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128  TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite
TestAC2SealNoticeLandsInTheRunLegLogFile  TestAC2TestDataDirBranchStillResolves128
TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord  TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot
TestAC3EmptyDataRootIsARefusalNotAFallback  TestAC3InstallingTheFileSinkDoesNotSilenceTheConsole
TestAC3LogSinkLandsInsideTheEnvDataRoot  TestAC3SecretLegBooksItsAuditRecordsOnDisk
TestAC4EveryLegIsNailedOrRuled  TestComposedGateBlocksAWriteForTwoSeconds
TestHostDispatchThroughTheAssembledBridge  TestMissingBlobFailsUnconfiguredNeverSilently
TestProcessCommandLineProbeDetectsAPlantedValue  TestProcessCommandLineProbeHelperProcess
TestProvidersDiscoverListsWhatTheServerServes  TestProvidersProbeRecordsMeasuredThinkingFalse
TestProvidersProbeRecordsMeasuredThinkingTrue  TestProvidersProbeUnconfiguredRefIsNotSilentlyKeyless
TestRunTextTaskFailNextIsClassified  TestRunTextTaskKeyResolvesInTheStore
TestRunTextTaskTextPathEndToEnd  TestSecretArgvCarriesNoSecret
TestSecretEndToEndConfigRefResolvesAtRequestTime  TestSecretFailurePathsLogAndPrintNoPlaintext
TestSecretFlagsAreBoolOnly  TestSecretFromStdinWritesNoIntermediateFile
TestSecretOverwriteIsAnnounced  TestSecretPortableModeUsesTicket06Seam
TestSecretRealBinaryRefusesValueFlag  TestSecretSameNameUnderThreeEnvsIsThreeBlobs
TestSecretSetConfirmationMismatchStoresNothing  TestSecretSetGetListRoundTrip
TestSecretSetRejectsBadNamesAndEmptyInput  TestSecretSetWithoutConsolePointsAtFromStdin
TestSecretUnsetRefusesWhileReferenced  TestSecretUsageAndUnknownSubcommand
TestSecretValueCarryingFlagsAreRefusedAndUnechoed  TestTicket101ManualSwitchSurvivesRestart
TestTicket101ModeSwitchUsesTheRealL2Gate  TestTicket101SessionGrantDoesNotCrossRestart
TestTicket101UnreadableModeFailsLoudlyAndStrict  TestTicket101UntouchedConfigRestartsAtDefault
```

`--- FAIL` 名册：**零枚**（基线三发都无红名）。名册差集：`R2-BLF` 54 ＼ `R2-BLC` 53 ＝
只多 `TestAC4EveryLegIsNailedOrRuled`（131 的门被 `-skip` ⇒ **少跑**，`--- SKIP` 行仍 0），
反向差集空。⇒ 基线与 133 验收方 §1.1 第 2 批的 `101/54/0/0`、§1.9 的"关门是少跑不是 SKIP"
**本程独立复现**（本程基线没碰上那枚既有 flake，`R2-BLF` 零红）。

**恒真自查第一条落在这里**：新判据 (c) 在未变异基线上**不响**（`R2-BL1`／`R2-BLF` rc=0，
且 disclosure 明写"1 distinct, 1 of them startable"）⇒ 它不是一枚"今天已经全响"的装饰。

〔独立复现〕本程自己跑的五发（表里 `R2-RST` 是 §6 那枚还原发，顺带落在同一张表里）；
另有 §0.4 记的那一发争用态读数 `R2-B1` **作废、不在本表任何结论里**。

## 4. `M-G` 本发 ＋ 两拍 ＋ 旧尺对照（先证落地，再读数）

### 4.1 种法（本程自造，前两程的名字一枚都没用）

台件 `/d/tmp/wisp135mg-r2-mutate.py`（生产形）、`/d/tmp/wisp135mg-r2-probe.py`（用例形）。
票 133 验收方用的是 `sfx131`／`probe133r2b_linux_test.go`／`TestR2P7LinuxOnlyCaseDrivesTheLeg`，
票 135 上一程用的是 `sfx135g`／`probe135mg_linux_test.go`／`TestMGGp1…` ⇒ **本程一枚都不复用**：
腿 `sfx135r2`、entry `cmdSfx135r2`、落地文件 `sfx135r2leg.go`、
植物文件 `probe135r2m_linux_test.go`、用例名 `TestR2MgLinuxTaggedCaseDrivesThePlantedLeg`。

- `mgleg1` ＝ **第一拍**：新腿被 `func main` 分发、身体里装了监听器、钉表里没有它。
- `mgleg2` ＝ **第二拍**（票面写死那一枚）：把那截 install 整个删掉 ⇒ 落地件里
  `installLogSink` 出现 **0 次**（本程的落地闸要求 `any-hits=0`，注释里提一嘴都不许）。
- `usage` ＝ 只加给人看的那行 usage（对照：这条腿今天**确实**没人覆盖）。
- `mgp1` ＝ `M-G` 本发：签名完全正确、带 import、身体里真断言的 `func Test…(t *testing.T)`，
  坐在**文件名自带隐式 GOOS=linux 约束**的文件里 ⇒ windows 的 `TestGoFiles` 不收它。
- `mgp2` ＝ 反向控制：同一形，写在**本程会编译**的文件里（§5）。

⚠ **本程自己的一发落地失败，登记不藏**：batch b1 的前六发（树 T04…T09）
在 `LANDED` 之后被**本程自己的落地闸**拒掉（`installLogSink-hits=1` ⇒ `MUTATION-WRONG mgleg2`，
驱动 `exit 6`，**没读任何红绿**）。原因是第二拍那枚文件的**注释里**写了安装器的名字。
⇒ 修种法（注释换成"the listener install block"）后**另起 T12…T17 重发**（`R2-*2` 系列），
T04…T09 那六棵作废树留着不删。这一条按派单口径也算一次自证：
**没证到落地的那批发确实一个读数都没产生**。

### 4.2 落地证明（每一发都先看这三行，再看数）

以 `R2-MGn2`（新尺 × 第二拍 × `M-G` 植物）为例，`/d/tmp/wisp135mg-r2-out/batch-b1r.txt` 原文：

```
LANDED-MAIN mgleg2 (dispatch sfx135r2 owns the leg):
  main.go:81: 		case "sfx135r2":
  sfx135r2leg.go lines=19 installLogSink-code-hits=0 installLogSink-any-hits=0
  main.go:44:   wisp sfx135r2    ticket 135 M-G probe leg, not a product command
  probe135r2m_linux_test.go:9: func TestR2MgLinuxTaggedCaseDrivesThePlantedLeg(t *testing.T) {
PROBE mgp1 applied
GATE-FILE sha1=23b443ac34787aa9ea60e181b9b8c789f468cd56
BUILD-cmdwisp R2-MGn2 rc=0
BUILD-all   R2-MGn2 rc=0
```

第一拍那两发的落地行是 `sfx135r2leg.go lines=25 installLogSink-code-hits=1 installLogSink-any-hits=1`
＋ `sfx135r2leg.go:17: sink, err := installLogSink(args[0])`。旧尺那几发在落地行之前多一行
`GATE-FILE-SWAPPED-TO /d/tmp/wisp135mg-r2-oldgate.go`、之后 `GATE-FILE sha1=906201f4…`。
**19 发读数全部 `BUILD-cmdwisp rc=0` ＋ `BUILD-all rc=0`**，无一发出现 `build failed`。

### 4.3 读数（单跑本尺那一形，八数＋红名）

| 发 | 尺 | 种件 | rc | 八数（RUN/顶P/顶F/顶S/子P/子F/子S） | 红名 | 账本里那一行 |
| --- | --- | --- | --- | --- | --- | --- |
| `R2-Uo2` T12 | 旧 | 第二拍＋usage | **1** | 1/0/1/0/0/0/0 | 本尺 | `leg sfx135r2 … covered=RED nothing` |
| `R2-Un2` T13 | 新 | 第二拍＋usage | **1** | 1/0/1/0/0/0/0 | 本尺 | `leg sfx135r2 … covered=RED nothing` |
| `R2-MGo2` T14 | **旧** | 第二拍＋`M-G` | **0** | 1/1/0/0/0/0/0 | **无人** | `leg sfx135r2 … covered=test TestR2MgLinuxTaggedCaseDrivesThePlantedLeg drives cmdSfx135r2` |
| `R2-MGn2` T15 | **新** | 第二拍＋`M-G` | **1** | 1/0/1/0/0/0/0 | 本尺 | 同上那一行（账本一字未变） |
| `R2-1Bn` T10 | 新 | 第一拍＋`M-G` | **1** | 1/0/1/0/0/0/0 | 本尺 | `leg sfx135r2 … installs=true covered=RED sink with no nail` |
| `R2-1Bo` T11 | 旧 | 第一拍＋`M-G` | **1** | 1/0/1/0/0/0/0 | 本尺 | 同上 |

**这一族里最省字的一句**：`R2-MGo2` 与 `R2-Un2` 之间只差那一枚植物，
旧尺上"红 → 绿"（钉一枚不存在的用例名就过了闸），新尺上"红 → 仍红"。

### 4.4 红名逐字（`R2-MGn2.log` 第 3 行，本程 `grep` 出来的原文）

```
leg_dispatch_gate_133_test.go:229: AC#1/#2 RED: leg "sfx135r2" (main.go:81) is booked in the
ledger below as covered by the case "TestR2MgLinuxTaggedCaseDrivesThePlantedLeg", and this
round's test binary has no such case: the roster read from the running binary lists 54
startable cases and "TestR2MgLinuxTaggedCaseDrivesThePlantedLeg" is not one of them, so no
"=== RUN   TestR2MgLinuxTaggedCaseDrivesThePlantedLeg" line exists in this run or can exist
in it. The name is declared at probe135r2m_linux_test.go:9, which is the point - the
declaration is in the sources and the sources are not the build.
```

⇒ 红名**逐名指到那枚名字**、说的是"这轮没有、也不可能有它这一行 `=== RUN`"，
不是"缺标记"（票面 ③ 要的那个口径）。同一发的披露行：

```
run-roster disclosure: GOOS=windows, 54 startable cases read from this binary itself
(`-test.list '.*'`); case names this round's ledger credited through `covered=test`:
2 distinct, 1 of them startable (one outside the roster is red above); this gate's registry:
4 claims, 0 of them outside this round's roster - every registry claim this gate makes is a
case this round's binary can start. This round was narrowed by -test.run=…, so the roster
proves these cases are startable here, not that each one printed === RUN in this process…
```

### 4.5 整包两形（票面那条"第二拍也必须红"必须在整包形上也成立）

| 发 | 尺 | 形 | rc | 八数 | 红名 | SKIP 行／build-failed／panic |
| --- | --- | --- | --- | --- | --- | --- |
| `R2-MGFc` T20 | 新 | 门关着（`-skip=^TestAC4…$`）第二拍＋`M-G` | **1** | 100/52/1/0＋47/0/0 | `TestAC1AC2DispatchHopGate133` | 0/0/0 |
| `R2-MGFo` T21 | **旧** | 同上 | **0** | **100/53/0/0**＋47/0/0 | **无人** | 0/0/0 |
| `R2-MGF-n` T22 | 新 | **摘掉本尺**（`-skip=^TestAC1AC2…$`） | 0 | 100/53/0/0＋47/0/0 | 无人 | 0/0/0 |
| `R2-1BFc` T23 | 新 | 门关着，**第一拍**＋`M-G` | **1** | 100/52/1/0＋47/0/0 | 本尺 | 0/0/0 |

- `R2-MGFo` 那一发**就是票 133 §1.12 那两发 `100/53/0/0` 的独立复现**（本程自己的树、自己的名字、
  自己的日志：`/d/tmp/wisp135mg-r2-out/R2-MGFo.log`，`grep -c '^--- FAIL'` = **0**，
  那枚植物名在整份日志里出现 **1 次**、唯一那次是本尺自己打的账本行 `:85`，
  `=== RUN` 名册里 **0 次**）。
- `R2-MGF-n`：摘掉本尺之后整包 rc=0 ⇒ **131 的门不接这一形**。本程还直接读到它确实看见了这条腿：
  131 的门在同一次运行里打了 `leg sfx135r2 main.go:81 install=false records=false ruled=false nails=- -> no records`
  （`R2-MGF-n.log:94`）然后 **PASS** ⇒ "看见了但不管"是盘上读数，不是本程的推断。
- 名册差集（对 `R2-BLC` 基线）：`R2-MGFc`／`R2-MGFo`／`R2-1BFc` 三发的 `=== RUN` 名册都是 53 枚，
  **与基线逐名相同、零缩小、零新增**；逐名换色只有
  `TestAC1AC2DispatchHopGate133 PASS→FAIL` 一枚（旧尺那发连这一枚都没有）。
  ⇒ 植物那枚用例**没有**进名册（这正是被验的那件事），也**没有**吞掉任何一条别的读数。

### 4.6 恒真判据自查（每一发动手前问的那句"这一发今天响不响"）

| 这一发 | 在**未修**那侧（旧尺）响不响 | 在**修了**这侧（新尺）响不响 | 能不能当"修好了的证据" |
| --- | --- | --- | --- |
| `R2-Un2`／`R2-Uo2`（只加 usage） | 响（rc=1，`covered=RED nothing`） | 响 | **不能单独当**——两側都响，它只证"这条腿今天确实没人覆盖"，是对照不是凭据 |
| `R2-1Bn`／`R2-1Bo`（第一拍） | 响 | 响 | 不能当 (c) 的凭据（它响在 `installs` 那一支，比这次改动老）；它算"两拍里第一拍没被改坏"的证据 |
| **`R2-MGo2`→`R2-MGn2`（第二拍＋`M-G`）** | **不响（旧尺 rc=0 全绿）** | **响（新尺 rc=1，红名逐点到那枚名字）** | **能——这一枚才是"今天不响、修了才响"那一枚** |
| `R2-CPn2`（反向控制，见 §5） | 不响 | **不响**（必须不响） | 它是"不许打死合法形状"的证据，不是修复证据 |

⇒ 本格没有拿"旧码上已经全响"的发当凭据；也没有拿"两拍里第一拍的红"去抵第二拍（票 133 §3.4 那条
"同一枚事实的两种写法不许互相抵账"在这里是可算的：`R2-1Bn` 与 `R2-MGn2` 各各独立、
红名分别落在 `installs` 那一支与 `runRosterReds135`）。

〔独立复现〕以上每一发都是本程自己建树、自己种、自己 `grep` 红名；
日志在 `/d/tmp/wisp135mg-r2-out/batch-b1.txt`（作废六发）、`batch-b1r.txt`、`batch-b2.txt` 与同名 `.log`。

## 5. 反向控制：同一形写在**本会编译**的文件里 ⇒ 必须绿

```
R2-CPn2  新尺 + 第二拍 + mgp2（probe135r2m_test.go，无平台后缀）
         EIGHT rc=0 RUN=1 TOPPASS=1 TOPFAIL=0 TOPSKIP=0
         leg sfx135r2 … covered=test TestR2MgCompiledCaseDrivesThePlantedLeg drives cmdSfx135r2
         run-roster disclosure: GOOS=windows, 55 startable cases …;
         case names this round's ledger credited through `covered=test`: 2 distinct, 2 of them startable
R2-CPo2  旧尺 + 同一发：rc=0（旧尺本来就放行，说明这枚控制不是新尺造出来的偏门）
```

⇒ 新尺的分辨力**只在"这枚名字在不在本轮二进制里"这一维**上：
`mgp1` 那发名册 **54 枚**／`mgp2` 这发名册 **55 枚**（多出的正是本程那枚编译得进去的用例），
一红一绿，其余字节全同（同一棵基线树、同一条腿、同一个 entry 调用、同一行 usage）。
⇒ "不许把合法形状打死"这一条成立。

⚠ 本程**没有**做的事，别被本节误导：这枚控制只测了 windows 腿。把同一形换到
`//go:build` 表达式为假（而不是文件名后缀）那一支、以及 `GOARCH` 那一支，本程未造未读（§8 第 5 条）。

## 6. 还原：证与基线逐字相同、复绿

**本程的还原是结构性的**（派单第 3 节那套规矩的直接后果）：每一发读数都用一棵**新** `cp` 出来的树，
工作树从头到尾没被种过件 ⇒ 没有"还原漏一行"这种失败可发生。可本程仍然按票面把三样都量了。

### 6.1 变异面逐枚列（`diff -rq 快照树 ↔ 那发用的树`）

```
$ diff -rq /d/tmp/wisp135mg-r2-tree /d/tmp/wisp135mg-r2-T15        （M-G 本发那棵树）
DIFF: Files .../cmd/wisp/main.go and .../cmd/wisp/main.go differ        ← usage 一行＋dispatch 两行
DIFF: Only in T15/cmd/wisp: probe135r2m_linux_test.go                   ← 植物
DIFF: Only in T15/cmd/wisp: sfx135r2leg.go                              ← 那条腿
DIFF: Only in T15: wisp.exe                                             ← `go build ./...` 的副产物，非源文件

$ diff -rq /d/tmp/wisp135mg-r2-tree /d/tmp/wisp135mg-r2-T26        （还原发那棵树）
DIFF-T26: Only in T26: wisp.exe            ← 只有构建副产物；源码与快照逐枚同
（除这一行外零 DIFF ⇒ 还原发的树与基线树源码逐字相同）
```

`wisp.exe` 是本程驱动里 `go build ./...` 在**仓外快照树根**落的构建副产物（17 发里每一发都有），
**从没进过仓库**；§6.3 的残留扫描一并覆盖它。

### 6.2 还原发复绿 ＋ 与基线那发逐字比

`R2-RST`（T26，无种件、新尺、`-count=1 -run=^TestAC1AC2DispatchHopGate133$`）
⇒ **rc=0，1/1/0/0＋0/0/0，SKIP 行 0、build-failed 0、panic 0**，红名零枚。

```
$ sed -E 's/\([0-9]+\.[0-9]+s\)//' R2-BL1.log > bl1.norm ; 同一法做 R2-RST.log ; diff -u
@@ -1,4 +1,4 @@          ← 只差这一枚 slog 行的 time= 戳
-time=2026-09-24T11:11:44.827+08:00 level=INFO msg="winsec: sealing path resolver installed" …
+time=2026-09-24T11:29:34.558+08:00 level=INFO msg="winsec: sealing path resolver installed" …
@@ -15,4 +15,4 @@        ← 与这一行 go 自己打的耗时
-ok  	github.com/CarlosShao/wisp/cmd/wisp	0.135s
+ok  	github.com/CarlosShao/wisp/cmd/wisp	0.126s
```

⇒ 除时间戳与构建耗时两行，**还原发的日志与基线发逐字相同**（账本 11 legs、4 claims、
名册枚数、披露行文案全同）。

### 6.3 仓库侧三查

```
$ git status --porcelain                       （空 ⇒ 工作树只有本程那一枚证据文件、且已提交）
$ sha1sum cmd/wisp/leg_dispatch_gate_133_test.go
  23b443ac34787aa9ea60e181b9b8c789f468cd56      ← 与 §1.2 取件时同一枚，一把尺没被动
$ grep -rIl "sfx135r2" . --exclude-dir=.git | wc -l
  0                                              ← 种件的腿名/文件名字在仓库里零残留
```

〔独立复现〕§6 三样都是本程自己量的，命令原样在上面。

## 7. 门禁账

两批门禁各自前的闸门：`GATES post at=2026-09-24 11:31:05 +0800 → GATE-WORKER-HITS: 0`、
`GATES pre at=11:32:07 → GATE-WORKER-HITS: 0`，两批的 `gh run list` 都是三枚 `completed`
（零枚 in_progress）。台件 `/d/tmp/wisp135mg-r2-gates.sh <树> <标签>`，
**两棵树都是仓外归档树**（被审树＝`wisp135mg-r2-tree`＝`8369b24`；对照树＝`wisp135mg-r2-tree-pre`
＝`fa35557^`＝`d3e3a43`，即"上一程写码之前"那一版）。

**工具版本写明**（派单那条硬规矩）：

```
go       = go version go1.27.1 windows/amd64        （宿主）
gofumpt  = v0.12.0 (go1.27.1)  路径 $(go env GOPATH)/bin/gofumpt.exe ＝ D:\work\base\gopath\bin
         ⇒ 用的是宿主现成 binary，本程没有、也禁止任何人执行 `go install mvdan.cc/gofumpt@…`
```

| 项 | 被审树（`8369b24`） | 对照树（`d3e3a43`＝改前） | 判 |
| --- | --- | --- | --- |
| `gofmt -l ./cmd/wisp/` | 0 行 | 0 行 | 净 |
| `gofmt -l .`（全仓） | 0 行 | 0 行 | 净 |
| `gofumpt -l ./cmd/wisp/` | 0 行 | 0 行 | 净 |
| `gofumpt -l .`（全仓） | 0 行 | 0 行 | 净 |
| 宿主原生 `go vet ./cmd/wisp/` | **rc=0**，诊断 0 行 | rc=0，0 行 | 净 |
| 宿主原生 `go vet ./...` | **rc=0**，诊断 0 行 | rc=0，0 行 | 净 |
| `GOOS=linux go vet` 逐包（33 枚） | 31 枚 rc=0 且输出 0 行；**2 枚 rc=1** | 逐包账本**归一树名后 diff 逐字节相同** | 见下面归因 |
| `sh scripts/d22scan.sh` | **rc=0**（含阳性对照 `runtests.sh OK top PASS=21/FAIL=0/SKIP=0、=== RUN=31、'[no tests to run]'>=0`） | rc=0，同对照同数 | 净 |
| 逐枚现量各 scope examined | 见下表 | 见下表 | **差值全 0 ⇒ 不降** |

### 7.1 两枚 rc=1 逐错误行归因（不拿"整树 rc=1＝工具链假象"糊过去）

```
github.com/CarlosShao/wisp/cmd/balldebug  rc=1 lines=1
  package …cmd/balldebug: build constraints exclude all Go files in D:\tmp\wisp135mg-r2-tree\cmd\balldebug
  → 本程逐枚量到该目录三枚文件的第一行全是 `//go:build windows`：
      cmd/balldebug/diff_windows.go  //go:build windows
      cmd/balldebug/main.go          //go:build windows
      cmd/balldebug/shot_windows.go  //go:build windows
    ⇒ 没有一枚 linux 侧文件 ⇒ 这是这个包**按设计**的构建约束，不是编译破口。
github.com/CarlosShao/wisp/cmd/wisp  rc=1 lines=3
  package github.com/CarlosShao/wisp/cmd/wisp
  	imports github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx
  	imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in
  	  D:\work\base\gopath\pkg\mod\github.com\k2-fsa\sherpa-onnx-go-linux@v1.13.8
  → 触发点是本仓的一行：cmd/wisp/doctor.go:14  sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
    死在第三方包的 cgo 构建约束（宿主交叉时 CGO_ENABLED 不为 1），一枚 `cmd/wisp\…go:行` 的诊断都没有。
```

⇒ 两枚都**早于**被审改动、与那把尺无关，凭据是"两棵树的逐包账本把树名归一后 `diff -u` **空输出**"：
`/d/tmp/wisp135mg-r2-out/cross-pre.errs` ↔ `cross-post.errs`（同目录另有两棵树的
`gates-{pre,post}.vet-cross-ledger.txt` 全量 33 枚）。其余 31 枚包（含 `internal/**` 全部、
`internal/proc`、`internal/winsec`、`tools/signmodels`）linux 交叉 vet **全 rc=0 且 0 行**。

### 7.2 d22scan 各 scope 逐枚现量（票面"examined 枚数不降"）

| scope | 被审树 `8369b24` | 改前树 `d3e3a43` | 差 |
| --- | --- | --- | --- |
| bans #1-5 `internal/` | 203 | 203 | 0 |
| bans #1-5 `cmd/` | 22 | 22 | 0 |
| ban #6 `frontend/` | 40 | 40 | 0 |
| ban #7 `internal/tools/` | 18 | 18 | 0 |
| ban #8 `design/` | 16 | 16 | 0 |
| ban #8 `frontend/` | 40 | 40 | 0 |
| ban #8 `internal/` | 404 | 404 | 0 |
| ban #8 `cmd/` | 39 | 39 | 0 |
| 合计（`examined 225 production Go files`） | **225** | **225** | 0 |

（八枚 scope ＋ 一枚合计行；本程现量与票 137 验收方那批逐枚同数，本程不据以自证、只算两棵树的差为 0。）

### 7.3 "空 scope 必须 fatal"＝本程自己现跑，不是引用规矩

```
$ (cd tools/d22scan && go build -o /d/tmp/wisp135mg-r2-out/d22scan-r2.exe .)   BUILT rc=0 sha1=45dd96664752…
$ ./d22scan-r2.exe -root /d/tmp/wisp135mg-r2-tree            → rc=0（同一枚 binary 的正常形，八 scope 见 7.2）
$ 搭一枚假根 /d/tmp/wisp135mg-r2-fakeroot：整树复制后把 frontend/ 清空
  （FAKE-ROOT frontend files=0，目录在、文件零 ⇒ 这是"declared live 但走不到东西"那一形）
$ ./d22scan-r2.exe -root /d/tmp/wisp135mg-r2-fakeroot        → **rc=2**
  d22scan: ban #8 scope frontend/ examined 0 files - it is declared in emojiScopes() but walks
  nothing. Point it at a real tree or delete the entry; never leave a scope pretending to scan
  (ticket 71 AC#4)
```

⇒ 空 scope 不是"静悄悄 0 枚"而是**退出码 2 并点名那一枚 scope**，本程第一手读到。
同一批 `sh scripts/d22scan.sh` 里那枚阳性对照子用例也跑绿了：
`--- PASS: TestBuiltBinaryGoesRedEndToEnd/empty_live_scope_exits_2`（连同
`seeded_violation_exits_1`、`frontend_tree_gone_while_declared_live_exits_2` 等六枚，见
`gates-post.d22scan.txt`）。

### 7.4 `-count=2 -v` 两形四数 ＋ 名册差集（⚠ 四数会翻倍，不当红名变多）

| 发 | 形 | rc | RUN | 顶 P/F/S | 子 P/F/S | 名册枚数 | SKIP 行／build-failed／panic |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `R2-C2N` T24 | 门开着（整包） | 0 | **202** | **108/0/0** | 94/0/0 | 54 | 0/0/0 |
| `R2-C2C` T25 | 门关着（`-skip=^TestAC4…$`） | 0 | 200 | 106/0/0 | 94/0/0 | 53 | 0/0/0 |

- 逐名账：把顶层 `--- PASS/FAIL/SKIP` 按名字聚合计数，`R2-C2N` 里 54 枚**每一枚恰好 2 次**
  （`uniq -c | grep -v '^ *2 '` 空输出）⇒ 零换色、零"第二遍少跑"。
- 名册差集：`R2-C2N.runnames` ↔ 基线 `R2-BLF.runnames` **双向差集为空**；
  `R2-C2C` ↔ `R2-BLC` 同样为空。
- 这四数与票 133 验收方 §1.1 第 2 批那两发（202/108/0/0）**同数**；本程没在这两发里碰上那枚既有 flake。

### 7.5 一票本程自加的额外读数：linux 腿**真跑**（容器，不是交叉编译）

派单没要求这一发，本程加它是因为上一程的 commit message 把 linux 推给了"证据 §5"，
而那一格如果只留"取不到"，读者就永远不知道那把尺在 linux 上是红是绿。

```
$ MSYS2_ARG_CONV_EXCL='*' docker run --rm -v D:/tmp/wisp135mg-r2-tree:/src:ro \
      -v D:/work/base/gopath/pkg/mod:/gomod:ro -w /src -e GOFLAGS=-mod=mod -e GOMODCACHE=/gomod \
      -e GOPROXY=off -e CGO_ENABLED=1 golang:1.27 sh -c 'go test -count=1 -v -run "^TestAC1AC2DispatchHopGate133$" ./cmd/wisp/'
挂载先证（防 Git Bash 把 -v 吞成空挂载那枚假绿）：
  /src/go.mod 883 字节、/src/cmd/wisp 34 枚文件、/src/cmd/wisp/leg_dispatch_gate_133_test.go 75627 字节
  Linux 6.6.114.1-microsoft-standard-WSL2 x86_64 / go version go1.27.1 linux/amd64
读数：LINUX-RC=0（真跑，不是只编译）
  run-roster disclosure: GOOS=linux, 41 startable cases read from this binary itself;
  case names this round's ledger credited through `covered=test`: 1 distinct, 1 of them startable;
  this gate's registry: 4 claims, **4 of them outside this round's roster** -
  TestAC1ResidentLegInstallsItsLogListenerOnDisk, TestAC2ModelsLegBooksItsHandOffVerdictOnDisk,
  TestAC2SealNoticeLandsInTheRunLegLogFile, TestAC3SecretLegBooksItsAuditRecordsOnDisk.
  A registry nail the build does not take is disclosed here and not reddened …
  --- PASS: TestAC1AC2DispatchHopGate133 (0.19s)  →  ok  github.com/CarlosShao/wisp/cmd/wisp 0.272s
```

⇒ 三件事被这一发钉住：**(1)** 名册读数在 linux 上**取到了**（41 枚），
`os.Executable()`＋`-test.list` 那一形跨平台成立；**(2)** 那处不对称（`covered=nail` 落名册外
只披露不弄红）**按它注释说的那样发生**，逐枚点名，本程第一手读到；**(3)** linux 上**没有误红**
（rc=0）。⚠ 这**不**等于"linux 整包没别的问题"——本程只跑了这一枚用例（`-run` 收窄），
那 19 枚 linux 侧既有 FAIL 一枚没碰（§8 第 6 条）。
⚠ 这一发采的时候第 3 行 `gh run list` 连报两次 `EOF`（11:39:45／11:41:20），
前两行（worker=0／toolchain=0）都是 0 ⇒ **这一发的 ci 那一维是"未取到"而不是"确认为空"**，照实登记。

〔独立复现〕§7 每一行都是本程自己跑的；日志：
`gates-post.*`／`gates-pre.*`／`d22-binary-*.txt`／`linux-run.log`／`linux-run2.log`，
全在 `/d/tmp/wisp135mg-r2-out/`。
