# 票 359 只读普查腿 359-a1 —— AC#0 现读（sameAbsoluteness 的 POSIX 分母）

- 腿 `359-a1`（只读）；写面只有本件 ＋ 仓外 `D:\tmp\wisp\20261021\359a1\**`。读数时刻 2026-10-21 19:4x +08。
- 锚定 `HEAD = e97cfd9cb6913fd4d7ef1de5b516a5b5a378482f`（`dev`）；工作树 `internal/winsec/resolve.go` 与 HEAD blob
  同一（`git hash-object` = `e0b0e22f...` = `git rev-parse HEAD:...`）。读数全部取自 `git archive HEAD`
  （`go.mod`/`go.sum`/`internal`）解出的仓外快照，**不读工作树**；票面行号按内容锚定后现跑复认。

## sec 1 前提新鲜度（逐条对 HEAD blob 复量）

| # | 票面断言 | 现读 | 判 |
|---|---|---|---|
| P1 | 定义在 `internal/winsec/resolve.go` | `HEAD:internal/winsec/resolve.go:519` `func sameAbsoluteness(a, b string) bool` | 成立 |
| P2 | 作为合取项参与某个决策 | 两枚：`:456`（`sameTree`）、`:552`（`answerInsideTree`） | 成立（票面说"某个"，实测两枚） |
| P3 | 注释在 `resolve.go:513-517` | POSIX 那段在 **`:516-518`**（3 行）；`513-517` 是代价段，不含 POSIX 句。自 `a45b2e90` 到 HEAD 该句一直在 516 | **不成立（实为 516-518）** |
| P4 | POSIX 得 104/60/0/0 | 容器 `golang:1.27`（`GOOS=linux`）`-count=2 -v` 得 **104/60/0/0**；`-count=1` 得 **74/46/0/0** | 成立（104/60 是 `-count=2` 的数） |
| P5 | `absoluteness` 在 POSIX 用例命中 0 | 21 枚非 `windows` tag 的 `_test.go` 里命中 **0**；命中全在 53 枚 `windows` tag 件里 | 成立 |
| P6 | 摘合取 => POSIX 零红 | 实测零红，四发都零红（sec 4） | 成立 |
| P7 | 正控拿得到差异 | 拿得到（sec 5） | 成立 |

## sec 2 sameAbsoluteness 定义与全部调用点（逐名，`git grep` 对 HEAD blob）

全仓 `git grep -n sameAbsoluteness HEAD -- '*.go'` **只命中 1 个文件、5 行**，均在 `resolve.go`：

| 行 | 内容 | 种类 |
|---|---|---|
| `:449` | `// Ticket 129 added the second dropped segment, sameAbsoluteness below, on the` | 注释（`sameTree` 文档） |
| `:456` | `if !sameVolume(a, b) \|\| !sameAbsoluteness(a, b) {` | **调用点 1** |
| `:487` | `// sameAbsoluteness reports whether two spellings hang their tails off the same` | 注释（本函数文档） |
| `:519` | `func sameAbsoluteness(a, b string) bool {` | **定义** |
| `:552` | `if !sameVolume(child, parent) \|\| !sameAbsoluteness(child, parent) {` | **调用点 2** |

即 **定义 1 枚、调用点 2 枚**；未导出、包外零命中、测试零直接调用。上游（非测试面）：

| 调用点 | 直接上游 | 再上游 |
|---|---|---|
| `:456` | `:405` `treeOwnershipFailureForPair` 第二见证腿 | `:338` `resolverTreeOwnershipFailure` => `:304` `resolverConformanceFailure` => `:135` `SetPathResolver`（安装期缝守） |
| `:456` | `winsec_windows.go:111` `noticeNamesTree` | 仅 Windows 编译 |
| `:552` | `:391`、`:405` | 同上两处缝守腿 |

## sec 3 那段 POSIX 结论逐字（HEAD `resolve.go:516-518`）

```
516|// On POSIX filepath.VolumeName is empty for every input and this is the same
517|// statement IsAbs makes about a leading separator, so "a/b" and "/a/b" stop
518|// reading as one tree here too - which is the same correction, one root over.
```

## sec 4 变异复现（POSIX 容器，`golang:1.27`，`GOOS=linux`）

改法：仓外快照里把守卫行改回票 129 之前形状（`python` 精确串替换、`assert count==1` 才写盘）：

```
:456  if !sameVolume(a, b) || !sameAbsoluteness(a, b) {              => if !sameVolume(a, b) {
:552  if !sameVolume(child, parent) || !sameAbsoluteness(child,parent) { => if !sameVolume(child, parent) {
```

**落地证明两把**：容器内 `grep -n '|| !sameAbsoluteness' internal/winsec/resolve.go` => **0 行**（rc=1，
锚定树 2 行）；`go build ./internal/winsec/` **rc=0**。`resolve.go` sha256 `8e64efd7…c6e4` => `2c20f5b9…9037`。

四数（`-count=1 -v`；四条 grep **分开报、未相减**）：

| 快照 | rc | `^=== RUN` | `^--- PASS` | `^--- FAIL` | `^--- SKIP` | 日志 |
|---|---|---|---|---|---|---|
| HEAD 原样 | 0 | 74 | 46 | **0** | 0 | `posix-base.log` |
| 只摘 `:456` | 0 | 74 | 46 | **0** | 0 | `posix-mut-tree.log` |
| 只摘 `:552` | 0 | 74 | 46 | **0** | 0 | `posix-mut-inside.log` |
| 两枚同摘 | 0 | 74 | 46 | **0** | 0 | `posix-mut.log` |

**零枚红复现**。RUN 比 PASS 多 28 枚是 Go 的子测试行（`=== RUN   T/sub`）；四数各数各的，不相减。

## sec 5 正控（硬闸）：差异在哪、差多少

主尺（确定性）：同一发变异、同一条固定 9 名的 `-run` 正则，在 **Windows** 上打红 4 枚。

| 快照 | rc | RUN | PASS | FAIL | SKIP |
|---|---|---|---|---|---|
| `win-base` 原样 | 0 | 11 | 9 | **0** | 0 |
| `win-mut` 两枚同摘 | **1** | 11 | 5 | **4** | 0 |

红名 4 枚：`TestAbsolutenessSpellingsAreNotOneTree` / `TestSeamGuardRefusesACandidateWhoseSecondWitnessDiffersOnlyInAbsoluteness` /
`TestAC1SeamVouchesForTwoRealObjectsThatDifferOnlyInAbsoluteness` / `TestAC1SealLandsInAForeignTreeWhenTheSeamAdmittedAnAbsolutenessBlindCandidate`。

分腿（同一正则，各摘一枚）：

| 快照 | 摘哪枚 | 落地后仍在的行 | rc | RUN/PASS/FAIL/SKIP | 红名 |
|---|---|---|---|---|---|
| `win-mut-tree` | `:456` | `552` | 1 | 11/6/**3**/0 | 上表第 1、2、4 名 |
| `win-mut-inside` | `:552` | `456` | 1 | 11/6/**3**/0 | 上表第 1、2、3 名 |

=> 两枚腿各有一枚"只靠它自己才红"的用例，不是一枚顺带钉住另一枚。红例原文（`win-mut-129only.log`）：
`:79 AC#1/AC#3 RED: sameTree("C:wisp129-trees\store-44440\artifact.txt", "C:\wisp129-trees\...") = true`；
`:85 ... answerInsideTree("C:wisp129-trees\...","C:\wisp129-trees\store-44440") = true`。

副口径（全包，附噪声说明）：HEAD 原样在本机 `go test -count=1 -v ./internal/winsec/` 已是
**121/68/2/0 rc=1**（红名 `TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`、
`TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2`，属票 342 系竞态、与合取无关，原样即红）；
两枚同摘后 **121/65/5/0 rc=1**（净 +3 红）。同一棵原样树另一次插桩发读 121/69/1/0 => 342 系红数在 1~2 抖。
**故正控只用主尺（固定 9 名子集），不用全包 FAIL 数。**

## sec 6 Windows 侧口径与读数（两侧都报）

口径：同一份 `git archive HEAD` 快照、同一台机、`GOOS=windows/amd64 go1.27.1`；用只观测不改判据的插桩件
（仓外，`probe-parts2`/`probe-parts4`）给两枚比较各加一个计数器与一个调用者栈帧。

| 读数（Windows 全包 `-count=1 -v`，判据未改） | 值 |
|---|---|
| `sameTree` 被问次数 / 其中两拼写绝对性不一致的对数 | **171 / 5** |
| `answerInsideTree` 被问次数 / 不一致对数 | **40 / 9** |
| 插桩发四数（与未插桩同形） | 121/70/0/0 与 121/68/2/0 两次，差异全在 342 系竞态 |

分歧对（= 合取项在决定判据）逐名归因：

| 调用者 | 对数 |
|---|---|
| `absoluteness_attribution_129_windows_test.go:78/81/84/87/92`（`TestAbsolutenessSpellingsAreNotOneTree`） | 5（比较面） |
| `absoluteness_attribution_129_windows_test.go:261`（`TestSeamGuardRefuses...DiffersOnlyInAbsoluteness`） | 4（缝守裁决面） |
| `export_test.go:49`（`TreeOwnershipProbeForTest`，调用者 `absoluteness_seam_landing_129_windows_test.go:246`，`TestAC1SeamVouchesForTwoRealObjects...`） | 3（真对象面） |
| `absoluteness_seam_landing_129_windows_test.go:370`（`SetPathResolver`，`TestAC1SealLandsInAForeignTree...`） | 2（落点面） |

五个数相加 5+4+3+2 = **14** = 5 + 9，与上面两行计数器一致；四面（比较 / 缝守裁决 / 真对象 / 落点）各有用例。

反面对照（POSIX 同一把插桩，全包）：`sameTree` 被问 **2** 次、不一致 **0**；`answerInsideTree` 被问 **8** 次、
不一致 **0** => POSIX 有调用分母（10 次）、判据分母为 **0**。归因全在两枚用例（`seam_probe_root_125_other_test.go:250/263/291/302`）：
`TestAC2POSIXSeamGuardStillRefusesEveryHostileShape125`、`TestAC2POSIXSeamAcceptsTheHonestPOSIXAnswer125`。

**POSIX 侧必然为 0 的结构原因**（不是"用例没写到"）：`pathComponents`（`winsec.go:405`）在 POSIX 上吃掉前导
分隔符、`filepath.VolumeName` 恒空，故 `"a/b"` 与 `"/a/b"` 的组件集完全相同 `["a" "b"]`。实测
（`posix-shape` 快照，判据打印）：原样 `answerInsideTree("a/b/c","/a/b")=false`、`("a/b/c","/a")=false`；
摘腿后两者均 `true`。**摘腿在 POSIX 上多出 3 处"读成一棵树"，74 枚用例一枚没响。**

## sec 7 注释结论钉没钉住：直接实测

把注释点名的那个对当面问判据（仓外 `zz_claim_359_test.go`，只打印不断言）：

| 快照 | `IsAbs` | `sameTree("a/b","/a/b")` | `answerInsideTree("a/b","/a/b")` |
|---|---|---|---|
| HEAD 原样 | false/true | **false** | false |
| 两枚同摘 | false/true | **true** | false |

=> 注释结论**是真的且可钉**（原样 false、摘腿 true），但**今天没有任何用例钉它**：摘腿后 POSIX 侧零红。
票面"注释声称的结论没有任何用例钉住"**成立，且比票面写的更强**。

## sec 8 门禁（工作树，本腿零写码字节）

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| build | `GOFLAGS= go build ./...`（清 `GOTMPDIR/TMPDIR/TEMP/TMP`） | **0** | 零输出 |
| d22 | `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=244 cmd/=42, #6 frontend/=116, #7 internal/tools/=25, #8 design/=39 frontend/=116 internal/=610 cmd/=145` |

`CGO_ENABLED` 未设 0；本腿零 `.go` 写入故未跑 `gofumpt`。

## sec 9 交付与越界

- 本腿 `.go` 写入 **0 枚**：`git status --porcelain -- 'internal/**/*.go' 'cmd/**/*.go'` = 0 行；
  全仓 `git status --porcelain -- '*.go'` 的 77 行全是既存未跟踪件（均在 `.scratch/` 下），本腿未碰。
- 越界改动 **0 处**：`internal/**`/`cmd/**`/`scripts/**`/`.github/**`/`.scratch/wisp/**`/`docs/reports/**`
  零改动；全部变异落在仓外 `D:\tmp\wisp\20261021\359a1\` 的 `git archive` 快照上（9 枚快照目录，
  只建不删、无 `.out`）。本文引用的日志一律相对 `D:\tmp\wisp\20261021\359a1\logs\`。
- 正控尺（可复跑）：`grep -c '^--- FAIL' D:\tmp\wisp\20261021\359a1\logs\win-base-129only.log` = **0**、同目录 `win-mut-129only.log` = **4**。
- 提交号：见 `git log -1 --format=%H -- docs/evidence/s1/359-a1-census.md`（本件自身那笔，只 commit 不 push）。
