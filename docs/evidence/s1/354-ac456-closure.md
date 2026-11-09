# 票 354 收口腿：AC#4 颜色与门禁 ＋ AC#5 接缝 ＋ AC#6 纪律格（本腿读数全部自跑）

- 日期：2026-11-09 00:0x +08。腿＝`354` 收口写腿（AC#4/AC#5/AC#6 的读数全部本腿现跑；⛔ 信任何前腿转述；零产码零 `.go` 字节）。
- 锚：主树起手 `git log -1`＝`9dbd468e3`；门禁读数全部跑在 `git clone --shared` 导出树
  `/d/tmp/wisp/20261108/354cls/tree-cls`（HEAD＝`cda38fd52`）。起手锚⇄导出树之间漂 **1 笔**（`cda38fd52`＝编排者
  A1648：翻 AC#2/AC#3 勾＋票 410 证据件）；**`git diff --stat 9dbd468e3 cda38fd52 -- internal/winsec/`＝空**
  （本腿现跑）⇒ 对 `internal/winsec` 的读数在两锚同值，导出树读数有效。
- 判据出处（三核①，现读 ⛔ 信转述）：票面 AC#4＝`.scratch/wisp/issues/354-identity-recheck-refuses-the-delete-so-the-callers-name-keeps-an-attackers-object.md:37`；
  AC#5＝`:38-39`；AC#6＝`:40`。证据件（三核②，全文现读）：`docs/evidence/s1/354-a1-residual-readings-r1.md`（55 行）、
  `354-r1-deterministic-teeth.md`（66 行）、`354-v1-verdict.md`（82 行）。
- 三核③：主树 `git status --porcelain` 非 untracked 行＝**0**（写面空闲起跑）；红/绿全部取自导出树。
- 本腿写面＝本件＋`.scratch/wisp/probes/354/cls/logs/{01-gates,02-d22scan,03-winsec-full-v,04-nevertruncates-count3}.txt`＋票尾 Progress log 一行。**`.go` 字节＝0**。
- 跑前口径照派单：`unset GOTMPDIR TMPDIR TEMP TMP`＋`export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`；**⛔ CGO_ENABLED=0（未设）**。

## 一、AC#4 颜色与门禁 —— windows 半边实跑＋unix 半边如实报

### 1. windows 半边（本腿自己跑，导出树 `cda38fd52`）

```
go test ./internal/winsec -count=1 -timeout 25m -v
⇒ rc=0；ok  github.com/CarlosShao/wisp/internal/winsec  27.637s
```

- **四数（含子测试，尺＝`grep -c`，日志 `03`）**：`RUN=129 / PASS=129 / FAIL=0 / SKIP=0`（顶层那把同尺＝78/78/0/0；两把一致：FAIL/SKIP 双 0）。
- **逐名红册＝0 红**：`grep -E '^\s*--- FAIL'` 全文 **0 命中**（含缩进子测试行）⇒ windows 侧今天没有一枚红名可列。
- 旁证（非判据）：今天前腿的三发独立整包全绿——`354-a1` `ok 20.148s`、`354-r1` `ok 37.254s`、`354-v1` `ok 35.199s`——加本腿 `ok 27.637s`＝**四发独立全绿**。

### 2. unix 半边（★如实：本机跑不了，运行时色归 CI）

- 本机是 Windows，**跑不了 linux 测试二进制** ⇒ `*_unix_test.go` 各用例的运行时红/绿本机取不到，**归 CI**。本腿⛔ 谎报跑了。
- 本机可交的两把编译期尺（含 `_test.go` 类型检查 ⇒ unix 半边仪器文件编译通过）：

| 尺 | 读数 |
|---|---|
| `GOOS=linux go vet ./internal/winsec/` | **rc=0**（日志 `01`） |
| `GOOS=darwin go vet ./internal/winsec/` | **rc=0**（同前例 `342-v3` 的交叉档口径） |

### 3. 四道门禁（导出树 `cda38fd52`，日志 `01`/`02`）

| 门禁 | 读数 |
|---|---|
| gofumpt（`/d/work/base/gopath/bin/gofumpt.exe -l internal/winsec/`） | **空，rc=0** |
| `go vet ./internal/winsec` | **rc=0** |
| `go build ./...` | **rc=0** |
| `go test ./internal/winsec -count=1` | **`ok … 27.637s`，rc=0**（＝上面第 1 节同一发） |
| d22scan（导出树现编 `tools/d22scan` 独立 module，`-root .`） | **`clean - no D22 ban violations`，rc=0** |

- ★`./cmd/wisp` 整包＝**没跑**（票面 AC#4 逐字明禁，射程外）。

**判：AC#4 成立**——windows 半边实跑全绿＋逐名红册「0 红」＋四数在场；unix 半边按派单两条路交（编译期 vet 绿＋如实报「运行时归 CI」）。

## 二、AC#5 与上游两票的接缝

票面 `:38-39` 判据＝「那两族的**既有红⇄摆动名**（`…NeverTruncatesA…` 那族）今天什么状态，逐名写清；并写明本票**⛔ 重述**票 118／120-r3 已批的"link 只 unlink、⛔ 跟随"」。

### 1. 342/348/353/354 四族逐名（windows 侧，全部来自本腿整包 `-v` 日志 `03`，今天的色）

| # | 测试名 | 今天 |
|---|---|---|
| 1 | `TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`（＋`/junction`＋`/symlink-dir` 两子形） | **PASS** |
| 2 | `TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1` | **PASS** |
| 3 | `TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1` | **PASS** |
| 4 | `TestAC2WindowsSteadyStateAncestorJunctionReportsDamageAndRefusalAsTwoNumbers342r2` | **PASS** |
| 5 | `TestAC1WindowsNameSwappedAfterTheSealVerdictReportsDamageAndRefusalAsTwoNumbers342r2` | **PASS** |
| 6 | `TestAC1WindowsAncestorMoveIsRefusedWhileAParentHandleIsHeld342r2` | **PASS** |
| 7 | `TestAC2WindowsPinRefusalOnAnAbsentAncestorIsNotFiledAsASealRefusal348r1` | **PASS** |
| 8 | `TestAC2SeamWindowsTheRefusalPathRemovalCannotBeMovedByAnAncestorFlip348r2` | **PASS** |
| 9 | `TestAC2WindowsNameTopUpBeforeTheReCheckLeavesTheTopUpObjectNamed353r1` | **PASS** |
| 10 | `TestAC2WindowsRemovalStillRunsWhenNobodyTopsUpTheName353r1` | **PASS** |
| 11 | `TestAC2WindowsAfterTheRefusalTheCallersNameStandsCarryingTheForeignObject354r1` | **PASS** |
| 12 | `TestAC2WindowsRemovalStillRunsWhenNothingTopsUpTheName354r1` | **PASS** |

（12 枚顶层名＝14 枚 PASS 行；整包红名册 0 枚，无一族例外。）

### 2. `…NeverTruncatesA…` 那族专项（公共块「摆动名先复量 3 发」口径，日志 `04`）

```
go test ./internal/winsec -run 'TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1' -count=3 -timeout 25m -v
⇒ rc=0；ok 5.068s；3 轮 × junction＋symlink-dir 两形全 PASS；truncated=0 ×6 条读数
```

⇒ **今天的状态＝绿、三发零摆动**：342-v3／348-r2 时代那枚「truncated=1..3 of 400」既有红，本腿三发在这台面**一次都没取到**（0/3）。
⇒ 「既有红⇄摆动名」今天的读法＝**该族今天两半都是绿读数，无在飞的摆动名**；历史红名册（342-v3/348-v2/v3 逐名在册）属于票 342/348 的账，今天在这台面不可复现。

### 3. 那族的 unix 半边（如实）

`TestAC1POSIXPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1` 及各 `*_unix_test.go` 用例：本机跑不了（同 AC#4 第 2 节口径），运行时色归 CI；编译期 `GOOS=linux go vet` rc0。⛔ 用 windows 读数代 unix 判断。

### 4. 接缝陈述＋不重述 118/120-r3

- 被量的那一跳（对 `354-r1` 证据件 §四 的码读数）：`internal/winsec/pinnedleaf_342_windows.go` 的 `unlinkPinnedLeaf`——同一性复验 `p.root.Lstat(base)`⇄`os.SameFile(pinned, now)` 不匹配 ⇒ `:277-283` `return` 先于 `:284` `_ = p.root.Remove(base)` ⇒ 调用者名字下留着攻击者对象。仪器经包内**已声明接缝** `applyDescriptor`（`winsec_windows.go:35-38`）驱动，零新接缝。
- **本票不重述票 118／120-r3 已批的"link 只 unlink、⛔ 跟随"**：那两票已 `-done`，本票只测"复验拒删之后那一跳"，`link` 语义一字未碰、一句话未重述（自查：本腿 0 枚 `.go`、342/348 两票原句 0 字节触碰——`git show --name-only` 本腿笔见第五节）。

**判：AC#5 成立。**

## 三、AC#6 纪律格（自查逐条读数）

| # | 要求 | 读数 | 判 |
|---|---|---|---|
| 1 | ⛔ 翻任何框（本票各格由非实现者裁；本腿只加 Progress log 一行） | 票面框尺：起手 `^- \[ \]`=**3**／`^- \[x\]`=**4**；本腿交件后同值 3/4（本腿 0 勺；AC#4/5/6 的翻归编排者凭本件裁） | 过 |
| 2 | ⛔ 改票名／`Status:`／加 `-done` | 票文件名未动；`^- Status:` 命中＝1 行未碰；票名无 `-done` 后缀 | 过 |
| 3 | ⛔ push | 本腿 **0 次 push**（只 commit） | 过 |
| 4 | 改动件 LF | 本件＋票面 `grep -c $'\r'`＝**0** | 过 |
| 5 | `git add`/`commit` 两处显式 pathspec 点名 | 见第五节，两处均逐枚字面路径 | 过 |
| 6 | ⛔ 动 342/348 两票任何原句；⛔ 动产码 | 本腿笔 `git show --name-only`＝证据件＋probes 日志＋票 354 尾一行，**0 枚 `.go`**、342/348 票面 0 字节 | 过 |

**判：AC#6 成立。**

## 四、总裁决

- **AC#4 成立**（windows 半边 129/129/0/0 全绿＋逐名红册「0 红」＋四门禁 rc0；unix 半边编译期两把尺绿＋运行时如实归 CI）；**AC#5 成立**（四族逐名 12 枚全 PASS＋NeverTruncates 族三发零摆动＋unix 如实归 CI＋不重述 118/120-r3 写明）；**AC#6 成立**（上表六条自查全过）。
- ⇒ **票 354 具名判「可结」**：`AC#0`（凭 `211-354-closure.md`）＋`AC#1`（已勾）＋`AC#2`/`AC#3`（凭 `354-v1-verdict.md`，编排者已翻）＋`AC#4`/`AC#5`/`AC#6`（凭本件）——七格无未满足判据。**编排者翻 AC#4/5/6 三勾＋票名加 `-done`**（本腿⛔ 翻框）。
- ★一句如实的边界话：`internal/winsec` 的运行时色今天是 windows 半边的实跑读数；unix 半边本机只有编译期凭据，整票 `-done` 不冒认 unix 运行时已量。
- 撤销口令：revert 本腿这笔 commit 即回原状；导出树与 `/d/tmp/wisp/20261108/354cls/` 全为临时件、不入库。
