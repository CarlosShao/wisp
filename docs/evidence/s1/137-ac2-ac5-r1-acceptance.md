# 137 AC#2／AC#5 —— 验收方第一程（r1）：两格各裁一次，读数全部我自己打

裁决方：本程（**非实现方**；实现方是 `worker-ticket137-ac2`，它的自述是
`docs/evidence/s1/137-ac2-ac5-impl.md`——**那是被审对象，不是本件的依据**）。
判据依据＝票面 AC#2／AC#5 ＋ `docs/evidence/s1/137-ac1-r1-acceptance.md` 里终裁方写的
**R-137-1／R-137-2／R-137-3**（原文核，不采纳实现方对判据的再解释）。

**只裁 AC#2 与 AC#5 两格。AC#3／AC#4 一格未裁、一格未做。**

## §0 锚点、树、容器与三档口径

- **被验版本＝`4e66817`**。这枚 sha 的来历：派单点名它，我没有直接拿它当锚，先现核
  `git cat-file -t 4e66817` = **`commit`**；再自量 `git rev-parse --short HEAD`（开工时＝`4e66817`，
  本程中途 HEAD 被兄弟推进到 `f6d21fe`／`191f0d6`，**我的所有读数仍在 4e66817 的归档树上，不受影响**）。
  顺手补一条归因账：`git merge-base --is-ancestor 4e66817 HEAD` = 真，且
  `git diff --name-only 4e66817 HEAD -- internal/winsec/` **输出为空**（0 枚）⇒ 我这一程的结论对当前线仍然成立。
- **没读脏工作树**：开工时 `git status --porcelain` 只有 `docs/evidence/s1/136-ac12-ac13-r1-acceptance.md`
  一枚已改（别人的活）；取件一律
  `mkdir -p /d/tmp/wisp137r2-tree-tight && git archive 4e66817 | tar -x -C /d/tmp/wisp137r2-tree-tight`，
  之后**所有读数、所有变异、所有门禁都在仓外**（仓内未建 worktree、未 checkout／switch／stash／reset／
  --amend／rebase／clean，**仓内一枚 `go build`／`go test` 都没跑**）。
- **本程建的临时件（只建不删，路径全列）**：

  | 路径（`D:\tmp\` 下） | 内容 | 用来量什么 |
  | --- | --- | --- |
  | `wisp137r2-tree-tight` | `4e66817` 纯净树（＝收紧后、生产码未变异） | AC#5 五组门禁；`tight-plain-1/2`、`tight-link-1/2` 四发 |
  | `wisp137r2-ac2-tree-anchor` | `f53ad5c` 纯净树（AC#2 之前那版） | 门禁与 d22scan 的"锚点侧各 scope 数"对照 |
  | `wisp137r2-ac2-tree-loose` | 4e66817 树＋把三枚测试件还原成 `f53ad5c` 版（未收紧） | 还原态核字（`diff -q` 见 §2） |
  | `wisp137r2-ac2-tree-d-loose` | 上树＋**我自己的 MUT-D** | `d-loose-plain`／`d-loose-link` |
  | `wisp137r2-ac2-tree-d-tight` | 4e66817 树＋**我自己的 MUT-D** | `d-tight-plain`／`d-tight-link` |
  | `wisp137r2-ac2-tree-naive-d` | 4e66817 树＋MUT-D＋把 `refusalCreditsLink137` 换成裸 `strings.Contains`（**验收侧探针 PROBE-F，不是修法**） | §4 的恒真判据 |
  | `wisp137r2-ac2-tree-ab-tight` | 4e66817 树＋**我自己的 MUT-AB**（两条走查彻底不看链接） | §4 第二问 |
  | `wisp137r2-ac2-tree-e-tight` | 4e66817 树＋**我自己的 MUT-E**（只挪两枚生产文案里那截标记，错误身份与效果不动） | §4 第三问（fail-closed） |
  | `wisp137r2-ac2-fakeroot` | 手工搭的假仓根（`frontend/` 空） | §6.5 证"空 scope 是 fatal" |
  | `wisp137r2-ac2-io\` | 台件 `mutate.py`／`run.sh`／`go.sh`／`batch.sh`／`parse.py`／`compare.py`／`gates.sh` ＝7 枚，
    加 12 发原始 `-v` 日志（`*.v.log`）＋容器头（`*.head.txt`）＋名册三件套（`*.run/.fail/.skip/.colour.txt`）
    ＋门禁日志（`gates.out`、`gates-<树>.{pkgs,vet-native,vet-cross,d22scan}.txt`） | 全部读数可复核 |

  ⚠ 前缀按要求是 `wisp137r2-`。**注意这一族目录名上一程终裁方也用过**（`wisp137r2-tree0`／`wisp137r2-io` 等），
  我一枚未覆盖、一枚未删：我的树全部另起 `-ac2-`／`-tight` 后缀，`wisp137r2-io`（上一程的台件）只读不写。
- **实现方指纹的抽验（不是拿它的数当我的数）**：我这棵纯净树里三枚被测文件的 `md5sum` =
  `cb8350cb5e5917a685c10b39d57d5275`／`20c141df792de655408138cdaebf633b`／`50d4e54c72167d969b361c092cb3e644`，
  与它自述打印的三个值**逐字相同**；未变异的四发里两枚生产码 = `a6144c880de80e43bb1393f3624e7221`（`winsec.go`）／
  `b5056918be4ed13817d236fbcae0f477`（`winsec_other.go`），与终裁方 §1 与它 §0 报的两个值同 ⇒
  **三程读的是同一版底线码**〔日志＋归档，抽验〕。
  另记它 §0 一句**说错的话**：它把还原版（`ce181f23…`／`0ae03ee4…`／`b84c85c8…`）解释成"CRLF 版、行尾差同 go.mod 那 28 字节"。
  我现量三枚还原版文件 `CR=0`（纯 LF，见 §2 表）⇒ **那三枚的差异是内容（＝AC#2 那一刀本身），不是行尾**；
  这不影响任何读数，但别再顺着"CRLF"这句话往下推。
- **容器**：镜像 `golang:1.27`（image id `3680233e3204`，本机已有、未拉网络件），容器内
  `go version go1.27.1 linux/amd64`；宿主 `go version go1.27.1 windows/amd64`。
  挂载 `MSYS_NO_PATHCONV=1` ＋ `/d/...`，**每一发**开跑前硬闸 `[ "$(wc -c < /src/go.mod)" = "883" ] || exit 97`
  （归档树 883 字节＝上一程记的行尾差，我复现到同值），并 `ls -l` ＋ `md5sum /src/go.mod`；
  两形形状硬断言（软链形 `ln -s /r2ac2priv /r2ac2link` ＋ `readlink` 逐字核；普通形断言 `/r2ac2link` 根本不许存在）＝
  98；**落地不过证就不取颜色**（`go build ./...` 与 `go vet ./internal/winsec/` 任一非 0 ⇒ exit 95）。
  **12 发无一命中 97／98／95，12 发全 `BUILD_RC=0` ＋ `VET_RC=0`**〔独立复现〕。
  模块缓存 `GOPROXY=off` ＋ 复用既有卷 `wisp137r2-gomod`（离线缓存，缺件直接 fail），构建缓存用我自己新建的
  `wisp137r2-ac2-gobuild`。
- **三档口径**：本件里
  - 〔独立复现〕＝我自己在这棵树上跑出来的（§1 的逐枚给号、§2 的六发、§4 的三探针、§6 的五组门禁全是这一档）；
  - 〔日志＋归档，抽验〕＝引别人读数但我抽过（§0 的三枚文件指纹、§6.4 的基线六数与终裁方逐数对点）；
  - 〔仅自述，不背书〕＝只有实现方说过、我没复算的（§6.6 的 CI 步归因、它 §1.4 的"软链形只存在于容器台件里"那句宿主无分母的话）。

## §1（甲）判据点是不是真的全覆盖 —— 逐枚现量，不漏 108 那两枚

判据（终裁方原文，不采实现方再解释）：**R-137-2** "收紧后的判据必须同时覆盖 helper 与 108 那两枚内联断言
（否则按票面字面只改 helper，AC#3 要求的 11 枚软链形全转红必缺 2 枚）"；票面 AC#2 本体
"它必须核被拒的那一个路径**就是本用例自己种下去的那一个**"；三条禁令＝不新增依赖／不放宽／不拿被测函数算期望值。

我在 `4e66817` 纯净树上逐枚 `grep -n`（**行号是我现量的，不是抄它的**）：

| 判据点 | 它声称 | 我读到的原文（截断） | 判定 |
| --- | --- | --- | --- |
| `placement_symlink_113_other_test.go:131` | 共享尺的标记常量 | `const linkAttributionMarker137 = " through the link at "` | 是 |
| `:145` | 共享尺本体 | `func refusalCreditsLink137(err error, planted string) bool {` | 是 |
| `:174` | helper 加第四参数 | `func assertRefused113(t *testing.T, what, spelled, planted string, err error) {` | 是 |
| `:186` | 新分支 | `if !refusalCreditsLink137(err, planted) {`（红行在 `:187`） | 是 |
| `:204/:237/:254/:271/:296` | 113 五个调用点 | 五枚 `assertRefused113(t, "...", spelled/link, **link**, ...)`，第 4 实参逐枚都是本用例自己种的 `link` | 是 |
| `placement_leaf_118_other_test.go:81/:107` | 118 两处 | `assertRefused113(t, "SealFile", link, link, winsec.SealFile(link))`／同形 `PrivateFile` | 是 |
| `ancestor_separator_108_other_test.go:84` → 红行 `:90` | 第一枚内联腿 | `} else if !refusalCreditsLink137(err, link) {` → `t.Errorf("AC#2 RED: RemoveUnlinked(%q) refused with %v, which does not credit the link this case planted at %q...` | 是 |
| `:151` → 红行 `:156` | 第二枚内联腿 | 同上形状（`link := filepath.Join(root, linkName)` 在 `:140` 提出来了） | 是 |

**"有没有漏"我用两条独立查法，不靠上面那张表自证**〔独立复现〕：

1. 调用面普查：`grep -rn "assertRefused113" internal/winsec/` ⇒ 全仓命中＝1 枚定义（`:174`）＋
   **113 五处 ＋ 118 两处**，`ancestor_separator_108_other_test.go` 里只有 `:85` 那行**注释**提到它（"never called
   assertRefused113"），**不是调用**；`grep -rn "refusalCreditsLink137" internal/winsec/` ⇒ 定义 `:145` ＋
   helper `:186` ＋ **108 的 `:84`／`:151`**，共四处消费。⇒ 判据面 ＝ helper（1）＋ 108 内联（2）＝**三处，与 R-137-2 同形**。
2. 分母闭合：终裁方 §3 定的 11 枚（7 顶层＋4 子测试，`root` 仍来自 raw `t.TempDir()`）逐枚对上断言点 ——
   113 族 9 枚（`:204`、`:237` 的顶层＋4 枚 depth 子测、`:254`、`:271`、`:296`）全走 helper，
   108 族 2 枚走内联 ⇒ **11/11 有着陆**，没有第 12 处需要改而没改。
   并且我在 §2 的 `dtight-link` 那一发里按名册把它验了一遍：**红名 11 枚一枚不落、名册差集里没有"该红而没红"**。

**三条禁令的独立核**〔独立复现〕：

- **零新增依赖**：`git diff --numstat f53ad5c 4e66817` 里 `go.mod`／`go.sum` **零命中**；
  winsec 侧只有那三枚 `*_test.go`（`19/3`、`2/2`、`59/8`）；新增 import 只有标准库 `strings`（`:34`）。
  ⚠ 同区间还飘着 `internal/observe/sampler_settle_coverage_136_test.go`（7/2），那是兄弟票 136 的，不是它的。
- **未放宽任何断言**：`git diff f53ad5c 4e66817 -- internal/winsec/ | grep '^-'` **一共 13 行**，我逐行读完：
  3 行是 `t.Logf`（换成带 planted 的长版）、7 行是调用点（换成长一参的同名调用）、
  1 行是 helper 签名（加一参）、1 行是 `os.Symlink(foreignDir, filepath.Join(root, linkName))`（提成 `link` 变量再用）、
  1 行是 `linkTo113(t, root, ...`（返回值从丢弃提成变量）⇒ **没有任何一条判据被删除、被改成 `||`、被改成 Skip**；
  既有的"必须拒"＋"必须名 `ErrUnresolvedPath`／`ErrIsReparsePoint`"两条在两处都**一字未动**（`113:177/182`、`108:80/82`、`108:147/149`）。
- **期望值不来自被测函数（R-119-9）**：`planted` 逐枚都是 `linkTo113`／`filepath.Join(root, ...)` 的返回值（fixture 侧）；
  被测侧产物只作为**被读的那一个**进入比较。那截短语 `" through the link at "` 是**抄生产文案的字面量**
  （`winsec_other.go:158`、`winsec.go:262`），不是调用被测函数算出来的 ⇒ 这条红线没踩。

**登记一格本票地界外的同族形状**（不是 AC#2 的缺牙，是给编排者的料）〔独立复现〕：
`dataroot_symlink_119_other_test.go:203`／`:251`／`:308` 三处内联断言仍只判 sentinel、不判记名，
其中 `:203`（`TestAC1POSIXUnresolvedSymlinkedRootStillRefused119`）与 `:308`（`...InjectedTestDataDirStandsAsDeclared119`）
的根就是 raw `t.TempDir()`（`:194`／`:286`）⇒ 与 R-137-1 说的是同一种病；它们是票 119 的用例、不在本票 11 枚分母里，
**AC#2 不该被拿去替它们负责，但 AC#3/AC#4 记这笔账时别把它们当对照**。

**(甲) 结论：判据点全覆盖，零漏点 ⇒ 这一味齐。**〔独立复现〕

## §2（乙）MUT-D 三态：我一枚都不引用实现方的读数，全部自己重跑

判据（终裁方 R-137-1 原文的形状）：只有 **MUT-D** 那一形能分辨"改了"与"没改"——两条走查各只走前 3 个组件，
**仍拒宿主那枚链接、永远走不到用例自己种的链接**；不许钉在单发 A 或 B 上（那是恒不满足的坏尺）。

### 2.1 我这发变异长什么样（先证落地，再取颜色）

我自己写 `mutate.py`（锚串唯一性断言，命中数≠1 直接 exit 3），**逐处 `grep -n` 出被改后那一行的原文**：

- `internal/winsec/winsec_other.go:156-163`：注释 `// MUTATION-137-R2AC2-D: half-fix on the seal route - the walk stops after`
  ＋ `pieces := pathPieces(path)` ＋ `if len(pieces) > 3 { pieces = pieces[:3] }` ＋ `for _, prefix := range pieces {`
- `internal/winsec/winsec.go:284-288`：`// MUTATION-137-R2AC2-D: the same three-component stop on the unlink route.`
  ＋ `if len(prefixes) > 3 { prefixes = prefixes[:3] }`

改完之后 `pathPieces` 给 `/r2ac2link/w137r2ac2tmp/<用例>/002/root/link/keep-me.txt` 的 7 枚前缀里，
用例种的链接在**第 6 枚**，截到 3 枚之后两条路仍看得见第 1 枚 `/r2ac2link` ⇒
"照旧拒宿主链接、永远不碰自己那枚"由源码语义保证，再由 2.3 的读数验一遍。
**我的补丁字节与两枚前手都不同**（同一语义的第三份实现）：MUT-D 后
`winsec_other.go` = `29d9c8f5f779a035fe2ec080a6e17d2d`、`winsec.go` = `0a7113f72356cca70fd41cf9f4dc96f6`
（终裁方 `2ffb39b2…`／`2aaa102f…`，实现方 `9589636253…`／`8548aeb1…`）〔独立复现〕。

树的洁净度我用 `diff -q` 逐枚钉死（不是"我以为没改"）〔独立复现〕：

- 纯净树 `wisp137r2-tree-tight` 的五枚文件与 `git show 4e66817:<路径>` **逐字节相同**（`PRISTINE_MATCH` ×5）；
- 还原态树 `wisp137r2-ac2-tree-loose` 的三枚测试件与 `git show f53ad5c:<路径>` **逐字节相同**（`RESTORE_MATCH` ×3）
  ⇒ "未收紧"那一发用的确实是 AC#2 之前那一版断言；
- 三棵变异树（`d-loose`／`d-tight`／`naive-d`）的**生产码 md5 逐字相同**（上两枚值）⇒ 态与态之间唯一的差异
  就是"测试件是哪一版"，构造上排除了"上一发没还原干净"。

### 2.2 四发读数（`go test -count=1 -v ./internal/winsec/`，容器 `go1.27.1 linux/amd64`）

| 发 | 断言版 | 生产码 | 形 | RUN | 顶 PASS/FAIL/SKIP | 子 PASS/FAIL/SKIP | rc | 分母那 11 枚 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `d-loose-plain` | 未收紧（`f53ad5c` 版） | MUT-D | 普通 | 52 | 17/13/0 | 17/5/0 | 1 | **11 枚全 FAIL** |
| `d-loose-link` | 未收紧 | MUT-D | 软链 | 45 | 24/3/3 | 15/0/0 | 1 | **11 枚全 PASS＝零检测力（本票钉的那个害，在我自己这棵树上复现）** |
| `d-tight-plain` | 收紧（`4e66817` 版） | MUT-D | 普通 | 52 | 17/13/0 | 17/5/0 | 1 | 11 枚全 FAIL |
| `d-tight-link` | 收紧 | MUT-D | 软链 | 45 | 17/10/3 | 11/4/0 | 1 | **11 枚全 FAIL（逐名，见 2.4）** |

**四数只是入场券，比的是名册差集**〔独立复现〕：

- 软链形：`SETDIFF d-loose-link vs d-tight-link (FAIL) only-A=0 only-B=11` ⇒ 收紧之后**多出的红恰好 11 枚**、
  `only-A` 侧为空＝**没有任何一枚由红转绿、也没有任何一枚由绿转 SKIP**；
  `(RUN) only-A=0 only-B=0`（两版 45 枚名册 `diff -q` 无输出）＝**没有一枚用例消失**；
  `(SKIP) only-A=0 only-B=0` 且名册只有票 125 那三枚自拒探针。
- 普通形：`SETDIFF d-loose-plain vs d-tight-plain` 的 RUN/FAIL/SKIP **三档全部 only-A=0 only-B=0**，
  八数逐数相同（`52/17/13/0 ＋ 17/5/0、rc=1`）⇒ **这一刀在普通形是字面 no-op**，与终裁方 §5④ 给的参照值对点成功。
- panic 账：`grep -c '^panic|^fatal error'` 在**全部 12 份 `-v` 日志**里逐份 **0** ⇒
  没有"一条用例 panic 吞掉同包其余几十条读数"这种形状，那些数是全量读数不是残局。
- 跨树 RUN 名册：普通形四棵树（`d-loose`／`d-tight`／`naive-d`／纯净）52 枚**逐名相同**；
  软链形四棵树 45 枚**逐名相同** ⇒ 分母在态与态之间没有动过。

### 2.3 变异真落地（对照组＋反向证据）

- 三枚根已解析的对照（按 **R-137-3** 改点名后的那三枚，不是更正块原本点错的三枚）
  在 `d-loose-plain`／`d-loose-link`／`d-tight-plain`／`d-tight-link` **四发全红**〔独立复现〕：
  `TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed`、
  `TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed`、
  `TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`（红行原文可读，如
  `placement_leaf_118_other_test.go:81: AC#1 RED: SealFile("/r2ac2priv/...") returned nil, i.e. it sealed through a symlink...`、
  `dataroot_symlink_119_other_test.go:250: AC#3 RED: sealing through a link inside a resolved data root returned nil`）。
  ⇒ `d-loose-link` 那 11 枚的绿**不能**读成"变异没打到东西"。
- 按 R-137-1／R-137-3 我另外单独量了那两枚"名字里有 Still Refused 但同病"的 119 用例，
  四发颜色：`TestAC1POSIXUnresolvedSymlinkedRootStillRefused119` 与
  `TestAC2POSIXInjectedTestDataDirStandsAsDeclared119` 都是 **普通形 FAIL／软链形 PASS**（四发一致），
  与终裁方那句"它们自己也用 raw `t.TempDir()`"逐发相同 ⇒ 我**没有**拿它们当落地凭据，
  也没被这一格的新判据动到（收紧前后同色）。

### 2.4 `d-tight-link` 的 11 枚红名逐名，以及"红在哪一行"

红名册（`--- FAIL:` 行原文取名，四枚子测各算一枚；与 `d-loose-link` 的名册差集＝恰好这 11 枚）：

1. `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone` → 红行 `placement_symlink_113_other_test.go:204`
2. `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth`（顶层；它自己那条红是从子测派生的）
3-6. `.../link-at-depth-1`…`link-at-depth-4` → 红行 `:237`（`t.Helper()` 生效 ⇒ 报在调用点，与终裁方同一机理）
7. `TestAC1POSIXSealDirThroughASymlinkRefuses` → `:254`
8. `TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing` → `:271`
9. `TestAC1POSIXSealFileThroughABackslashNamedLink` → `:296`
10. **`TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink` → 红行 `ancestor_separator_108_other_test.go:90`（内联腿）**
11. **`TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor` → 红行 `:156`（内联腿）**

归因账（这条必须算清，否则"红"可能是别的因）〔独立复现〕：

- 那一发里带新判据文案（`which does not credit the link this case planted at`）的失败行**恰 10 枚**
  ＋ 顶层母项 1 枚派生 ＝ **11 枚，逐枚对号**；
- 同发里 108 那两枚的红**落在它们自己文件的 `:90`／`:156`**，不在 helper 里 ⇒ 若只按票面字面改 helper，
  这两枚没有红可落（R-137-2 的那 2 枚现在盘上有对点）；
- 没有"说不去的红"：那一发里 `f.untouched` 的读数是**before=after**（例
  `placement_symlink_113_other_test.go:206: ... before=mode=666 uid=0 gid=0 after=mode=666 uid=0 gid=0`），
  即外来树**没有**被 chmod、也没有被删 ⇒ 这 11 枚的红全部由那条新判据判出来，不是附带 Damage；
- 反向对照（收紧前那一发为什么是绿）：`d-loose-link.v.log:61` 的原文里，用例 PASS，而它拿到的拒记名的是宿主链接
  ——`... through the link at /r2ac2link ...`（`grep -o` 只取到 `/r2ac2link`，因为 MUT-D 把走查截到第 3 枚组件，
  `firstLinkAncestor` 压根没看 `root/link`）⇒ 绿得没有理由，这条我在自己日志里读到了原文。

**(乙) 结论：三态我自己重跑成功，与实现方那本账逐数相同、名册差集同形 ⇒ 这一味齐。**〔独立复现〕
（实现方自述的 `45/24/3/3＋15/0` → `45/17/10/3＋11/4`、普通形 `52/17/13/0＋17/5/0、rc=1` 三处，
我复算到**逐数相同**；这是抽验过的对点，不是拿它的数当我的数。）

## §3（丙）那八个字的出处：标签错位是我的派单造的，不是它造的

派单把 `52/17/13/0＋17/5/0、rc=1` 标成"AC#5 未变异 `-count=2` 普通形的参照值"，实现方回说这八个字在终裁证据里
是 **§5④ 的 MUT-D·普通形**。我独立判这一处置对不对，分三步：

1. **去原始出处核它当时是哪一发**〔独立复现，`grep -n` 现量〕：这八个字在
   `docs/evidence/s1/137-ac1-r1-acceptance.md` 只出现两处 —
   - `:286`（§5④正文）：`MUT-D·普通形 RUN=52 顶 PASS=17 顶 FAIL=13 顶 SKIP=0 子 PASS=17 子 FAIL=5 子 SKIP=0 rc=1`
   - `:463`（§10 下游门闸第 2 条 (c)）：**"普通形八数不变"要带参照值：`52 / 17 / 13 / 0 ＋ 17 / 5 / 0、rc=1`**
   两处的口径都是 **AC#3 的"普通形八数不变"参照值**，跑形是 **MUT-D·普通形·`-count=1`**；
   终裁证据里**从未**把它写成 AC#5 的 `-count=2` 未变异参照值（那里未变异普通形基线是
   `:91` 的 `52/30/0/0＋22/0`）。票面文末编排者那块（`:166`）同引一处，也是给 AC#3 的。
   ⇒ **实现方那句"值不错、标签错"在出处这一层成立。**
2. **自己采一发 AC#5 真正的未变异 `-count=2` 普通形**〔独立复现〕：`tight-plain-2`
   （纯净树 `4e66817`、无变异、`go test -count=2 -v ./internal/winsec/`）＝
   **`RUN=104 顶 60/0/0 子 44/0/0、rc=0`**，逐名跨两遍 `REPETITION_DRIFT=0`、`-count=1` 那一发（`tight-plain-1`）
   ＝ `52/30/0/0＋22/0、rc=0`（＝基线的整两遍，也与终裁方 `:91` 那行数逐数相同）。
   ⇒ 与实现方报的 `104/60/0/0＋44/0/0、rc=0` **逐数相同**。
3. **两组数冲不冲突**：不冲突。八个字属于"MUT-D·普通·`-count=1`"那一发（我自己两发都量到 `52/17/13/0＋17/5/0、rc=1`，
   收紧前后同数同名册，见 §2.2），`104/60/0/0＋44/0/0、rc=0` 属于"未变异·普通·`-count=2`"那一发。
   **数值层面零处不一致 ⇒ 不触发"立刻停下报回"那条**（那是给"两组数其实打架"准备的）。

**但要记一笔后果**（不是它的错，是这枚标签的杀伤力）：如果有人**照派单字面**去执行，AC#5 的 `-count=2` 普通形
会被拿去和 `52/17/13/0…` 对点 ⇒ 必然"对不上"，而最省事的凑法就是把 `-count=1` 的读数当 `-count=2` 交、
或者去动断言。这次的处置（两种读法都量、两边都留）是**对的**，也是唯一安全的走法。
⇒ 编排者该把票面/派单里这枚标签改掉：**八个字属 AC#3 的普通形参照值，AC#5 的 `-count=2` 普通形参照值是 `104/60/0/0＋44/0/0、rc=0`**。

**(丙) 结论：实现方的处置成立；数没错、名字错，错在我这边的派单。**

## §4（丁）"恒真判据是一类新假绿"：我用三枚反向探针独立重走

实现方自述排掉一枚新假绿：`strings.Contains(err.Error(), planted)` 在软链形**恒真**。这句话本身是它的断言 ⇒ 我重走。

**探针 PROBE-F（决定性那一发）**：在我自己的快照树里，把 `refusalCreditsLink137` 的**函数体整块换成**
`return strings.Contains(err.Error(), planted)`（落地原文＝
`placement_symlink_113_other_test.go:145-150`，注释标 `ACCEPTANCE PROBE 137-R2AC2-F`；生产码 md5 与
`d-tight`／`d-loose` 两棵**逐字相同**＝同一发 MUT-D），其余一概不动。

| 发（同一棵变异树，只差那一枚函数体） | RUN | 顶 P/F/S | 子 P/F/S | rc | 分母 11 枚 |
| --- | --- | --- | --- | --- | --- |
| `d-loose-link`（未收紧） | 45 | 24/3/3 | 15/0/0 | 1 | 11 枚 PASS |
| **`naive-d-link`（裸 Contains）** | 45 | 24/3/3 | 15/0/0 | 1 | **11 枚 PASS（枚枚点名，见名册）** |
| `d-tight-link`（收紧版，即进树那一版） | 45 | 17/10/3 | 11/4/0 | 1 | 11 枚 FAIL |

⇒ **裸 Contains 在软链形确实恒真**：它把读数打回"未收紧"那一发的同一形状（四数逐数相同、
FAIL 名册＝只有三枚对照、分母 11 枚全绿、RUN 名册 45 枚逐名相同），
而**进树的那一版在完全相同的输入上把 11 枚判红** ⇒ 新判据带的是信息，不是装饰。
反向也核了：`naive-d-plain`（普通形）＝ `52/17/13/0＋17/5/0、rc=1`，RUN/FAIL/SKIP 三档名册与
`d-tight-plain`／`d-loose-plain` **差集全空** ⇒ Contains 那一版只在软链形失去分辨力，不是"哪儿都红/哪儿都不红"，
排除了"两版其实是同一把尺"这种读法。**它自述的那枚假绿我复现到了，它确实绕开了。**〔独立复现〕

**派单点的那一发字面问法（"让 guard 根本不检查链接，看那枚新断言响不响"）我也量了**，答案是
"**不响，而且不该响**"：MUT-AB 那一发（`platformVerifyPlacement` 整段 `return path, nil` ＋
`firstLinkAncestor` 整段 `return ""`，我的字节 `59946c9a…`／`2534912d…`）·收紧版·软链形 =
`45/15/12/3＋11/4、rc=1`，`grep -c "does not credit the link"` = **0** ⇒
那 11 枚红在**既有**那两条判据上（`err=nil` 时 helper 在 `:177` 就 `return`，108 的 `else if` 链也进不去），
新判据在这一形**根本不会被走到**。这不是新判据没牙，而是它本来就是**第三条判据**、不是替换品；
真正的牙在 PROBE-F 那一发（err≠nil、拒的是宿主链接，只有分得出来才看得见）。
顺带这条也给 AC#3 提个醒：**别拿"MUT-A/MUT-B 那一形里新判据响不响"当凭据**（R-137-1 禁钉单发 A／B 就是这个机理）。〔独立复现〕

**MUT-E（它自加的探针，我也自己落一发）**：只把两枚生产文案里那截标记挪走
（`winsec_other.go:158`／`winsec.go:262` 改成 `via an ambient link: …`，
落地证据＝两枚文件里 ` through the link at ` 命中数 **0**、三枚测试件 md5 仍是进树版），
错误身份与效果一字未动。**未变异之外的这一发在普通形**（本应全绿的形状）读到
`52/21/9/0＋18/4/0、rc=1`，红名册 13 枚＝分母 11 枚＋`TestAC118POSIX…SealFile`／`…PrivateFile`（共用 helper），
而不经这条判据的 `TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119` 仍 PASS。
⇒ 文案一挪这条尺**响成一片、不会静默放行**（`refusalCreditsLink137` 取不到标记即 `return false`＝fail-closed，
`placement_symlink_113_other_test.go:147-150` 我逐行读过）。它的代价也说清：这一形里三枚对照组成了不了照。〔独立复现〕

**(丁) 结论：恒真那条不成立在进树版上；实现方"排掉一枚新假绿"这句我独立复现。**

## §5 那条未验证断言：未变异的收紧态在软链形（它决定 AC#3 会不会是"恒不满足"）

它自述：未变异的收紧态在软链形里也是这 11 枚红，而同形里根已解析的三枚对照仍 PASS ⇒ 这把尺可满足。
这一句是 AC#3 的地基，我单独复现：

| 发 | 树 | 形 | RUN | 顶 P/F/S | 子 P/F/S | rc |
| --- | --- | --- | --- | --- | --- | --- |
| `tight-plain-1` | 4e66817 纯净 | 普通 | 52 | 30/0/0 | 22/0/0 | **0** |
| `tight-plain-2` | 同上 | 普通 `-count=2` | 104 | 60/0/0 | 44/0/0 | **0** |
| `tight-link-1` | 同上 | 软链 | 45 | 20/7/3 | 11/4/0 | 1 |
| `tight-link-2` | 同上 | 软链 `-count=2` | 90 | 40/14/6 | 22/8/0 | 1 |

逐名（`tight-link-1`／`tight-link-2` 两发同形）〔独立复现〕：

- **分母 11 枚：全 FAIL**（名册与 §2.4 那 11 枚**逐名相同**，没有第 12 枚红）；
- **三枚根已解析的对照：全 PASS，而且不是 SKIP** ——
  `TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`、
  `TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed`、
  `TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed`
  两枚 118 走的正是被收紧过的 helper（`placement_leaf_118_other_test.go:81/:107`），
  它们在这一形满足"记名者＝自己种的链接" ⇒ **这把尺在软链形里可被满足，不是恒不满足**；
- SKIP 名册只有票 125 那三枚自拒探针（`-count=2` 下 6＝3×2），**分母 11 枚在 12 发里一枚都不是 SKIP**；
- 普通形未变异四数与终裁方/票 124 批次 3b 的基线**逐数相同** ⇒ 这一刀没有把普通形推红一寸（`rc=0` 照旧）。

⇒ **我不需要写"AC#3 这一格今天不可满足"**：结论与它的自述一致（软链形可满足），
并且今天这一格"未变异软链形也 11 枚红"这件事的成因是那 7 处 raw `t.TempDir()`（＝**AC#4 换根**那一格的地界，
**我一字未动、也未裁**）。这条账留给 AC#3/AC#4，不构成本格缺陷。

## §6 AC#2 一格终判：**成立（PASS、无附条件）**

票面 AC#2 的判据逐条结清（每条给档）：

| 判据（票面＋终裁 R-137-2） | 我的量法 | 读数 | 档 |
| --- | --- | --- | --- |
| 它必须核"被拒的那一个路径**就是本用例自己种下去的那一个**" | §1 逐枚给号＋调用面普查 | 共享尺 `:145`/`:186` 只在标记之后取记名者、要求等于 `planted` | 独立复现 |
| 判据点必须**同时覆盖 helper 与 108 那两枚内联断言**（R-137-2） | §1 第 108 行两枚 | `:84`→红行 `:90`、`:151`→红行 `:156`，且 §2.4 里红**落在 108 自己文件上** | 独立复现 |
| 分母闭合（11 枚一枚不落） | §2.4 名册差集 | `d-loose-link` vs `d-tight-link` FAIL 差集＝**恰好 11 枚**、RUN/SKIP 差集全空 | 独立复现 |
| 不许新增依赖 | `git diff --numstat f53ad5c 4e66817` | `go.mod`/`go.sum` 零命中；只 import 标准库 `strings` | 独立复现 |
| 不许放宽成"两形都算过" | `grep '^-'` 逐行读 13 枚删除行 ＋ §2 两形读数 | 无一条判据被删/被 Skip；两形用**同一句**要求，后果正是软链形 11 枚红 | 独立复现 |
| 不许拿被测函数算期望值（R-119-9） | §1 第 3 条 | `planted` 全来自 fixture；标记是抄生产文案的字面量 | 独立复现 |
| 不得是一枚恒真判据（一族新假绿） | §4 三枚反向探针 | PROBE-F：裸 Contains 让 11 枚回到全绿；进树版在同一输入上判红。MUT-E：文案一挪响成一片（fail-closed） | 独立复现 |
| 不得是一枚恒不满足的坏尺（R-137-1 反向那一族） | §5 | 未变异软链形三枚对照 PASS（非 SKIP）、普通形未变异 `rc=0` 全绿 | 独立复现 |

**出口＝成立，不是"附条件通过"**：四问（甲／乙／丙／丁）各成一节、全部我自己取证，没有一格需要"等下游补料才算数"。
两处**不翻账的登记**留给编排者（都不构成本格缺陷）：
①实现方 §0 那句"CRLF 版"是误解释（实为内容差）；
②`dataroot_symlink_119_other_test.go:203/:251/:308` 三处仍是 sentinel-only 断言（票 119 地界、本票分母外），
AC#3/AC#4 记料时别把它们当对照（R-137-1 已点到其中两枚）。



