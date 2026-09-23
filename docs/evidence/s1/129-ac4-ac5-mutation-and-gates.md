# 票 129 — AC#4（变异自证）+ AC#5（门禁）读数表

锚定 sha：`a71b2d8`（`internal/winsec/**` 自 `db9fafc`（AC#3）以来一字未动，实测
`git diff --stat db9fafc a71b2d8 -- internal/winsec` 空输出）。
本段由第二任代理接续（前任最后写入 09:03、静默三小时、无未提交 WIP），**只做 AC#4/AC#5**，
AC#1..AC#3 的判据与代码不在本段改动范围内。

工作面：仓库树零改动，全部变异落在 `git archive a71b2d8` 解出的快照目录
`/tmp/wisp129c-s23-*`（带本会话后缀 `s23`；仓库内零 worktree / 零 checkout，A38④）。

## STEP 0 — 基线重量（不复用前任的数）

`git archive a71b2d8 | tar -x -C /tmp/wisp129c-s23-base`，在快照目录里
`go test -count=2 -v ./internal/winsec/`（rc=0，日志 `/tmp/wisp129c-s23-base-run.log`）：

四数按 `scripts/winsec-tests.sh` 里那四条 grep 的**逐字形状**数（`^=== RUN` / `^--- PASS` /
`^--- FAIL` / `^--- SKIP`，即只数顶层行）：

| 读数 | 值 |
| --- | --- |
| `=== RUN` | 202 |
| `--- PASS` | 116 |
| `--- FAIL` | 0 |
| `--- SKIP` | 0 |
| `grep -c '(cached)'` | 0 |
| 顶层 `--- PASS` 去重名字数 | 58 |

⇒ 与前任在 AC#3 那一格自述的 `RUN=202 PASS=116 FAIL=0 SKIP=0`（单次 101/58）同形，
但我这次的 202/116 与那枚 58 枚的名字集合是**重新量的**（落盘
`/tmp/wisp129c-s23-base-pass.txt`），后面 AC#4 第④条的 diff 全部以这一枚为基线，不引用前任的数。

（AC#4 的三态读数与 AC#5 的门禁读数在下方逐格追加。）

## AC#4 — 变异自证（三态 + 拒绝侧 + 名字集合）

四枚快照，全部 `git archive a71b2d8 | tar -x -C /tmp/wisp129c-s23-<名>`，仓库内零 worktree / 零 checkout；
每一发**先 grep 证变异落地 + `go build` rc=0，再读红名**，做完逐发还原（还原后与锚定树 `cmp` 全等，见文末）。
`internal/winsec/**` 在本段一字未改（`git diff --stat -- internal/winsec` 空输出）。

### MUT-BOTH：把"绝对性"那一段从两枚比较里同时摘掉（＝改前语义）

补丁内容（各 1 处，python 精确串替换、`assert count==1` 后才写盘）：

```
resolve.go:456  if !sameVolume(a, b) || !sameAbsoluteness(a, b) {      ->  if !sameVolume(a, b) {
resolve.go:552  if !sameVolume(child, parent) || !sameAbsoluteness…    ->  if !sameVolume(child, parent) {
```

落地证明：`grep -c "|| !sameAbsoluteness" internal/winsec/resolve.go` = **0**（锚定树是 2），
剩下的三处命中全是函数定义与注释（`:449/:487/:519`）；`go build ./internal/winsec/` **rc=0**。

读数（`go test -count=2 -v ./internal/winsec/`，rc=1，日志 `/tmp/wisp129c-s23-mut1-run.log`）：
`RUN=202 PASS=108 FAIL=8 SKIP=0` ⇒ 顶层红名去重 **4 枚**，且**全部**是跨绝对性用例：

| 红名 | 红名原文（节选，逐字取自日志） |
| --- | --- |
| `TestAbsolutenessSpellingsAreNotOneTree` | `absoluteness_attribution_129_windows_test.go:79: AC#1/AC#3 RED: sameTree("C:wisp129-trees\\store-44440\\artifact.txt", "C:\\wisp129-trees\\store-44440\\artifact.txt") = true, but a drive-relative spelling hangs its tail off the process's current directory on C: and an absolute one off C:\ - two objects, one tree in this comparison's eyes`（同发还点到 `:82` 反序、`:85/:88/:93` 三腿 `answerInsideTree`） |
| `TestSeamGuardRefusesACandidateWhoseSecondWitnessDiffersOnlyInAbsoluteness` | `:279: AC#1/AC#3 RED on leg "second witness vouches for the probe parent's tree in a drive-relative spelling": the seam guard admitted a candidate whose answers name two objects …; it owed a refusal and said "".` 另有 `:274: AC#1 instrument measured the wrong leg …: witness "D:wisp129-seam\\probe-tree" asked=false, wanted true` |
| `TestAC1SeamVouchesForTwoRealObjectsThatDifferOnlyInAbsoluteness` | `:253: AC#1 RED: the seam's containment witness read "C:\\wisp129-vouch-27448\\tree\\leaf" as sitting inside the tree "C:wisp129-vouch-27448\\tree" names, and those are two directories on this box - one under the process's own current directory …, one at the volume root - planted for this leg with different contents` |
| `TestAC1SealLandsInAForeignTreeWhenTheSeamAdmittedAnAbsolutenessBlindCandidate` | `:380: AC#1 RED: the install-time seam ADMITTED a candidate whose second witness vouches for the probe parent's tree only in a drive-relative spelling …` ＋ `:396: AC#1 RED: S-1-1-0 was stripped from …\victim\sub\keep-me.txt … Grants that disappeared: [S-1-1-0]. The install-time seam had just admitted the candidate that answers with that path (admitted=true).`（`:393/:394` 给出 icacls 前后逐条：受害者 `sids=[S-1-5-18 S-1-5-32-544 …-1001]` 丢了 S-1-1-0，调用者自己那枚 blob 仍带 `Everyone`） |

⇒ **①达成**：摘掉那一段就红，且红名逐字点到"改前那枚跨绝对性用例"（比较面 + 缝守裁决面 + 真对象面 + 落点面各一枚）。

### 还原 ⇒ 绿（②）

同一目录 `git archive a71b2d8 internal/winsec/resolve.go | tar -x -C /tmp/wisp129c-s23-mut1` 还原：
`grep -c "|| !sameAbsoluteness"` = **2**、`git show a71b2d8:internal/winsec/resolve.go | cmp - …/resolve.go` 通过（**byte-identical**）、
`go build` rc=0、`go test -count=2 -v ./internal/winsec/` **rc=0**，四数 `RUN=202 PASS=116 FAIL=0 SKIP=0`。

### 分腿定位（同一发变异只摘一枚比较）

| 快照 | 摘掉的那一枚比较 | 四数（`-count=2 -v`） | 顶层红名（去重 3 枚） |
| --- | --- | --- | --- |
| `s23-mut2` | 只 `sameTree` | `RUN=202 PASS=110 FAIL=6 SKIP=0`，rc=1 | `TestAbsolutenessSpellingsAreNotOneTree`、`TestSeamGuardRefusesACandidate…`、`TestAC1SealLandsInAForeignTree…` |
| `s23-mut3` | 只 `answerInsideTree` | `RUN=202 PASS=110 FAIL=6 SKIP=0`，rc=1 | `TestAbsolutenessSpellingsAreNotOneTree`、`TestSeamGuardRefusesACandidate…`、`TestAC1SeamVouchesForTwoRealObjects…` |

两枚红名的**并集恰等于 MUT-BOTH 的 4 枚、交集为 2 枚**（比较面与缝守裁决面各一枚红名是共用的），
独有那枚红名跟着"摘掉哪一枚比较"换人 ⇒ 两枚 leg **各自有一枚只靠它自己才绿的用例**，
不是"其中一枚顺带把另一枚也钉住了"。

### ④ `--- PASS` 名字集合 vs 基线：只增不减（重新量）

还原态顶层 `--- PASS` 去重名字集（`/tmp/wisp129c-s23-mut1-restored-pass.txt`，58 枚）
与本段 STEP 0 基线集（58 枚）对：`comm -13` 增 = **0**，`comm -23` 减 = **0**。
⇒ 本段交件的树与锚定树同物，名字集合既没缩也没长（前任在 AC#3 量的 +8/−0 我**没有引用**）。
（变异态下集合是 **108→54 枚顶层 PASS**，红名把那 4 枚挤出去了——那是变异的预期后果，不进这一格账。）

### ③ 拒绝侧一枚不许变松

拒绝侧用例点名口径（可复跑）：从 winsec 测试源里取每枚顶层 `func Test…` 的函数体，
体内出现 `refus`/`Refus` 者记为拒绝侧成员。改前（`f5bbccd`，票 129 之前）对改后（`a71b2d8`）：

| 读数 | 值 |
| --- | --- |
| 拒绝侧成员 改前 / 改后 | 47 / 55 |
| **改前 − 改后（被删/被改名的拒绝侧用例）** | **0** |
| 改后 − 改前（新增） | 8 枚，**全部落在 `absoluteness_*_129_windows_test.go` 两枚文件里**（名单见下） |
| 票 126 的 `wantRefused:` 腿数 改前 / 改后 | true 3→3、false 4→4（一分未动） |
| 票 129 自己的腿数 | true 3 / false 3 |

新增 8 枚：`TestAbsolutenessSpellingsAreNotOneTree`、`TestSeamGuardRefusesACandidateWhoseSecondWitnessDiffersOnlyInAbsoluteness`、
`TestHonestPipelineKeepsPassingItsOwnProbe`、`TestAttributionFaceNeverSeesAMixedAbsolutenessPair`、
`TestAC1DriveRelativeAndAbsoluteSpellingsNameTwoObjectsOnThisBox`、`TestAC1SeamVouchesForTwoRealObjectsThatDifferOnlyInAbsoluteness`、
`TestAC2NoSealEverActsOnAGloballyDifferentAbsoluteness`、`TestAC1SealLandsInAForeignTreeWhenTheSeamAdmittedAnAbsolutenessBlindCandidate`。

**红名集合的差（改前语义=MUT-BOTH 那一发 vs 改后=还原态，只算拒绝侧成员）**：

- 拒绝侧红名 改前 = 4 枚（上表 MUT-BOTH 的四枚，全部 ∈ 本票新增的 8 枚）；改后 = **0 枚**。
- 改后 − 改前 = **0 枚**（没有任何一枚还原态的红是变异态没有的 ⇒ 改后不比改前"多放过"任何东西）。
- 改前 − 改后 = **4 枚**（这四枚的绿**只由那一段 leg 供给**，摘掉即红）。
- 拒绝侧 55 枚里在**两发里都绿**的 51 枚 = 47 枚存量（票 104/108/112/115/118/126/… 的拒绝腿，逐枚判定未变）
  + 本票 4 枚 CONTROL/代价腿（`TestAC1DriveRelativeAndAbsolute…`、`TestAC2NoSealEver…`、`TestAttributionFace…`、
  `TestHonestPipeline…`）——它们本来钉的就是"不许多拒"，改前改后都该绿，实测都绿。

**这套差值不是空仪器（两发反向对照）**：

- MUT5A（把票 126 一枚**存量**拒绝腿的期望从 `wantRefused: true` 翻成 `false`）⇒ `rc=1`、
  `--- FAIL: TestSeamGuardRefusesACandidateWhoseSecondWitnessNamesAnotherVolume`，
  红名原文逐字点到腿名：`volume_attribution_126_windows_test.go:394: CONTROL RED on leg "second witness names the same tree on another volume": the seam guard now refuses a candidate whose answers name one tree on one volume: …`
  ⇒ **放宽一枚拒绝期望，仪器当场响**。
- MUT5B（把票 129 表里一整枚 `wantRefused: true` 腿删掉，`wantRefused: true` 3→2）⇒
  `rc=0`、`RUN=2 PASS=2 FAIL=0 SKIP=0`、`--- PASS: TestSeamGuardRefusesACandidateWhoseSecondWitnessDiffersOnlyInAbsoluteness`
  **仍然绿** ⇒ 登记为**本包仪器的已知盲区**（与 `A109②` 同形：`winsec` 这两张腿表没有"腿数下限"那一枚断言，
  删腿不会自己变红）。所以"删除式放宽"这一形**上面那发行数读数（0 removed、126 腿 3/4 未动）才是承重的证据**，
  红名差值不是；我没给这枚下限，因为本段不改判据。

