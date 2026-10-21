# 票 359 写腿 359-r1 —— POSIX 侧补一枚用例当 `sameAbsoluteness` 合取的分母（AC#2）

- 腿 `359-r1`（写腿）；裁定形＝〔甲〕补一枚用例（票面 AC#1 编排者 2026-10-21 20:1x 拍）。
- 锚定 `HEAD = 3469d88ddd7555f20eeb41e513f3cb4c88f0ee0f`（`dev`）；开工时
  `git status --porcelain -- internal/winsec/` = **0 行**（无别人的未提交改动），`git hash-object
  internal/winsec/resolve.go` = `e0b0e22f...` = `git rev-parse HEAD:internal/winsec/resolve.go`。
- 本文所有数**本腿自己现跑**；普查件 `docs/evidence/s1/359-a1-census.md` 只当输入、不当判据。

## sec 1 票面前提新鲜度（逐条对 HEAD blob 复量；行号按内容锚定）

| # | 票面断言 | 本腿现读 | 判 |
|---|---|---|---|
| P1 | `sameAbsoluteness` 全仓 5 行 | `git grep -n sameAbsoluteness HEAD -- '*.go'` **5 行**：注释 `:449`／调用点 `:456`／注释 `:487`／定义 `:519`／调用点 `:552` | 成立（计数与逐名均同） |
| P2 | 注释在 `resolve.go:513-517` | POSIX 那句在 **`:516-518`**（`git show HEAD:...` 逐字读出）；`513-517` 是代价段、不含 POSIX 句 | **不成立（实为 516-518）**，与普查腿一致；本腿按内容锚定、未信票面行号 |
| P3 | 该合取在 POSIX 上摘掉零枚红 | 本腿复现：摘 `:456` / 摘 `:552` / 两枚同摘 => 既有 74 枚**零红**（sec 3 表） | 成立 |
| P4 | `sameTree("a/b","/a/b")` 原样 false／摘腿 true | 原样 `false`（本腿新用例绿）；摘 `:456` 后 **`true`**（新用例第 67 行红） | 成立，且本腿把它钉进了用例 |
| P5 | `pathComponents` 吃掉前导分隔符（结构原因） | `winsec.go:405` 逐字现读：`i := len(filepath.VolumeName(path))` 后 `for i < len(path) && isSep(path[i]) { i++ }` | 成立 |

## sec 2 新用例（逐名＋为什么今天响）

新增件：`internal/winsec/sameabsoluteness_359_other_test.go`（96 行，`//go:build !windows`，`package winsec`）。
**未导出面就在本包内**：`sameTree` / `sameAbsoluteness` / `sameVolume` / `pathComponents` / `answerInsideTree`
全部包内直接调用，**没有**改用例形状、没有走 `export_test.go`。

| 用例 | 断什么 | 为什么今天必需 |
|---|---|---|
| `TestSameAbsolutenessPOSIXIsTheOnlyLegBetweenTwoRoots359` | `sameTree("a/b","/a/b") == false`（第 67 行）＋镜像序（第 70 行）；`answerInsideTree("a/b/c","/a/b") == false`（第 73 行）＋镜像（第 76 行）；另有 6 条 CONTROL（第 83/86/89/92/95/100 行） | 这是本票**唯一**一枚问这条合取在 POSIX 上要什么判据的用例；摘任一处合取它就转红（sec 3），而既有 74 枚一枚不响 |

**为什么两枚调用点都要问**：那句话在 `resolve.go` 里出现两次（`sameTree` 守卫 `:456`、`answerInsideTree`
守卫 `:552`）。只问 `sameTree` 的话，摘 `:552` 时本用例**仍然绿**——正是普查腿记下的"合取没了、POSIX 全绿"那一形。
本腿实测两枚各有一枚只靠它自己才红的断言（第 67/70 行 vs 第 73/76 行）。

**先证"什么都没植"再断**（第 47-64 行，六条 `Fatalf` 全部是插桩自检、不是判据）：
`IsAbs(rel)=false`／`IsAbs(abs)=true`／`VolumeName` 两边相等／`sameAbsoluteness` 两边**不等**／
`sameVolume` 两边**相等**／`pathComponents` 两边**逐段相等**。最后两枚是关键：它们证明"这条对的判据
只能由 absoluteness 那一条腿给"，否则本用例测的是 volume 面或组件走查、不是合取。

CONTROL 六条（第 83-100 行）挡住"直接 `return false` 也能过"：`a/b` 与自身、`/a/b` 与自身、
`a/b` 与 `a//b`、两种括号内的正常包含，都必须为 `true`；`a/b` 与 `/a/c` 必须为 `false`（尾巴不同本该由走查拒）。
=> 本件是"两个根不是一棵树"，**不是**"拒一切相对输入"。

## sec 3 反向判据（本票的牙）：三读数逐字＋sha256 三处相同

口径：`golang:1.27` 容器、`GOOS=linux/amd64`（容器内 `CGO_ENABLED=1`，非本腿设 0）、`GOFLAGS=` 空；
快照＝`git archive HEAD go.mod go.sum internal` 解到仓外 `D:\tmp\wisp\20261021\359r1\snap` ＋拷入新用例。
每发 mutant 都先验落地（容器内 `grep -n '|| !sameAbsoluteness'` ＋ `go build` rc），**读完立刻单独还原**。

`resolve.go` sha256：

| 状态 | sha256 |
|---|---|
| 提交版（工作树／`git show HEAD:` blob／还原台） | `8e64efd7ed203156cd3de28c6062eea39346abba64de12127b4cd7cab399c6e4` |
| 只摘 `:456` | `01bdf424957b3c1940e60b9de1d5be9790ec7f6c28cea13b7d9db7e183514806` |
| 只摘 `:552` | `25b840d65ee6405ab388472bc7b4fd797c703a71ebec963c5d07e7e3b0a456f7` |

**三处逐字相同**（提交版／还原台／`git show HEAD:` blob）＝`8e64efd7...c6e4`，三次独立现跑同值。
还原方式是 `tar -xf base.tar internal/winsec/resolve.go`（从同一份 `git archive` 解回），非手改。

四发读数（`-count=1 -v`，四数各数各的、**未相减**）：

| 快照 | 容器内 `grep` 命中 | BUILD | rc | RUN | PASS | FAIL | SKIP | 红名（逐名） |
|---|---|---|---|---|---|---|---|---|
| `posix-asis` 原样 | 2 行（`:456`/`:552`） | 0 | **0** | 75 | 47 | **0** | 0 | 无 |
| `posix-mut-tree` 只摘 `:456` | 1 行（只剩 `:552`） | 0 | **1** | 75 | 46 | **1** | 0 | `TestSameAbsolutenessPOSIXIsTheOnlyLegBetweenTwoRoots359` |
| `posix-mut-inside` 只摘 `:552` | 1 行（只剩 `:456`） | 0 | **1** | 75 | 46 | **1** | 0 | `TestSameAbsolutenessPOSIXIsTheOnlyLegBetweenTwoRoots359` |
| `posix-restored` 还原 | 2 行 | 0 | **0** | 75 | 47 | **0** | 0 | 无 |

**红行逐字**（`posix-mut-tree`，摘 `:456`）：

```
    sameabsoluteness_359_other_test.go:67: AC#2 RED: sameTree("a/b", "/a/b") = true, but a relative spelling hangs its tail off the process's own directory and an absolute one off the filesystem root - two objects. This is the verdict the resolve.go paragraph states for this platform
    sameabsoluteness_359_other_test.go:70: AC#2 RED: sameTree("/a/b", "a/b") = true (mirror order of the leg above), so the refusal is not symmetric in its arguments
--- FAIL: TestSameAbsolutenessPOSIXIsTheOnlyLegBetweenTwoRoots359 (0.00s)
```

**红行逐字**（`posix-mut-inside`，摘 `:552`）：

```
    sameabsoluteness_359_other_test.go:73: AC#2 RED: answerInsideTree("a/b/c", "/a/b") = true, so a relative child is read as sitting inside the absolute tree "/a/b" and the seam guard's containment witness would vouch for a seal across two roots
    sameabsoluteness_359_other_test.go:76: AC#2 RED: answerInsideTree("/a/b/c", "a/b") = true (mirror direction): an absolute child is read as sitting inside the relative parent "a/b"
--- FAIL: TestSameAbsolutenessPOSIXIsTheOnlyLegBetweenTwoRoots359 (0.00s)
```

***分母形状对照**：既有 74 枚在四发里**全部同数**（RUN 74／PASS 46／FAIL 0，四发一致），
所以上表 RUN 75／PASS 47 比 74／46 多出来的**正好是本腿新增的这 1 枚**（`-run` 未加过滤，整包跑）。
换话说：**本腿把"摘了无人察觉"变成"摘了必有一枚红、且红行指到两处合取各自的语义"**。

## sec 4 两侧读数

**POSIX 侧（真跑了容器，非"假装"）**：见 sec 3 四发表。原样 **75/47/0/0 rc=0**、还原同值；
两发变异各 **75/46/1/0 rc=1**、红名 1 枚＝本腿新用例。容器内 `go vet ./internal/winsec/` 四发均 rc=0。

**Windows 侧（本机 `go1.27.1 windows/amd64`，工作树本体）**：

固定 9 名 `-run` 子集（名字取自普查件 sec 5，逐名对齐）：

```
^(TestAbsolutenessSpellingsAreNotOneTree|TestAttributionFaceNeverSeesAMixedAbsolutenessPair|TestSeamGuardRefusesACandidateWhoseSecondWitnessDiffersOnlyInAbsoluteness|TestHonestPipelineKeepsPassingItsOwnProbe|TestCrossVolumeSpellingsAreNotOneTree|TestAC1DriveRelativeAndAbsoluteSpellingsNameTwoObjectsOnThisBox|TestAC1SeamVouchesForTwoRealObjectsThatDifferOnlyInAbsoluteness|TestAC2NoSealEverActsOnAGloballyDifferentAbsoluteness|TestAC1SealLandsInAForeignTreeWhenTheSeamAdmittedAnAbsolutenessBlindCandidate)$
```

| 快照 | rc | RUN | PASS | FAIL | SKIP |
|---|---|---|---|---|---|
| 工作树原样（含本腿新件，本件在 Windows 构建里**不入列**） | **0** | 11 | 9 | **0** | 0 |

九名逐名 PASS：`TestAbsolutenessSpellingsAreNotOneTree`／`TestAttributionFaceNeverSeesAMixedAbsolutenessPair`／
`TestSeamGuardRefusesACandidateWhoseSecondWitnessDiffersOnlyInAbsoluteness`／`TestHonestPipelineKeepsPassingItsOwnProbe`／
`TestCrossVolumeSpellingsAreNotOneTree`／`TestAC1DriveRelativeAndAbsoluteSpellingsNameTwoObjectsOnThisBox`／
`TestAC1SeamVouchesForTwoRealObjectsThatDifferOnlyInAbsoluteness`／`TestAC2NoSealEverActsOnAGloballyDifferentAbsoluteness`／
`TestAC1SealLandsInAForeignTreeWhenTheSeamAdmittedAnAbsolutenessBlindCandidate`。**与普查件 sec 5 原样读数 11/9/0/0 rc=0 逐数相同。**

新件不在 Windows 构建里（不是断言，是 `go list` 现读）：`GOOS=windows` 下它出现在
`IgnoredGoFiles`（同列 20 枚既有 `_other`/`_unix` 件）；`GOOS=linux` 下出现在 `TestGoFiles`。
=> Windows 侧读数**不受本腿新增件影响**，两侧分母互不污染。

**全包（副口径，仅登记、NOT 当判据）**：本机 `go test -count=1 -v ./internal/winsec/` 两次
**121/67/3/0 rc=1** 与 **121/68/2/0 rc=1**，红名逐名：
`TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`、
`TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1`、
`TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1`（第三次那枚）。
三枚全在**票 342 系竞态**、原样即红、与本合取无关；两次红数 3 与 2 之差即该族在 1~2 抖。
与普查件 `121/68/2/0 rc=1` 同形。***故正控只用固定 9 名子集，不用全包 FAIL 数。**

***一处环境坑，登记给下一条腿**：本机若把 `TEMP`/`TMP` 清空（"清 `GOTMPDIR TMPDIR TEMP TMP`"那条按字面执行），
Windows 上 `os.TempDir()` 返回**空串**，`t.TempDir()` 因此给出 `\TestXxx...` 这种**没有盘符**的路径，
`filepath.IsAbs` 判否 => 固定 9 名子集里 4 枚假红（`TestAttributionFaceNeverSees...`／
`TestAC1DriveRelative...`／`TestAC1SeamVouches...`／`TestAC1SealLands...`），与本票无关。
本腿的处置：只清 `GOTMPDIR`/`TMPDIR`、`TEMP`/`TMP` 保留环境原值（`go test` 自身仍按 `GOTMPDIR` 走）。

## sec 5 门禁（本腿现跑，工作树）

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| build | `GOFLAGS= go build ./...` | **0** | 零输出 |
| vet | `go vet ./internal/winsec/` | **0** | 零输出 |
| test | `go test -count=1 ./internal/winsec/` | **1** | 见下四数（既有 342 系红，非本腿引入） |
| d22 | `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations`；bans #1-5 internal/=244 cmd/=42、#6 frontend/=116、#7 internal/tools/=25、#8 design/=39 frontend/=116 internal/=611 cmd/=145 |
| pathlen | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`，denominator=17528、over-budget=60/rostered=60/not in roster=0 |
| gofumpt | `& 'D:\work\base\gopath\bin\gofumpt.exe' -l <新件>` | **0** | 零输出（gofumpt `v0.12.0 (go1.27.1)`） |

`go test -count=1 ./internal/winsec/` 四数**分开报**（`-v` 计数）：两次 **121/67/3/0** 与 **121/68/2/0**，rc 均 1。
**按名作差：既有红没有增。** 三枚红全属票 342 系、原样即红；本腿新增件是 `//go:build !windows`，
在 Windows 构建里不入列（`go list` 现读，见 sec 4）=> **本腿对 Windows 侧红名册贡献 0 枚**。
341／342 系以外**新增红 0 枚**。*未用"两次相减"当判据，只按名逐个对齐。

`CGO_ENABLED` 未设 0（`go env` 现读默认 1）；本仓 sherpa 依赖需要它。

## sec 6 纪律与越界

- 写面：**只新增 1 枚测试件**（`internal/winsec/sameabsoluteness_359_other_test.go`）。既有测试件**零改动**，
  `resolve.go` **零字节改动**（sha256 三处相同已证），`internal/winsec` 其它生产码零改动。
- NOT 未改 `sameAbsoluteness` 判定语义、NOT 未改 `resolve.go:516-518` 注释（甲形不改注释）、NOT 未改 `pathComponents`。
- 全部变异落在仓外 `D:\tmp\wisp\20261021\359r1\**`（只建不删、无 `.out`；日志在 `logs\`）。
- `.gitattributes` 要求 `*.go text eol=lf`：新件 CR 字节 **0**；写盘走 `UTF8Encoding($false)`／tar 解回，无 CRLF 带入。
- 提交正文与本件均**无被禁码位**（现扫 U+00A7/U+2713/U+2264/U+2260/U+2192/U+21D2/U+26D4/U+26A0/U+FE0F/
  带圈数字/emoji 区/U+2200-22FF/U+2600-27BF/U+2B00-2BFF，命中 **0**）。
- 提交：只 commit、**NOT push**；`add`/`commit` 两处显式 pathspec（新件 ＋ 本件两枚），NOT `add -A`/`.`；
  NOT amend/reset/rebase/stash/checkout ./clean。
