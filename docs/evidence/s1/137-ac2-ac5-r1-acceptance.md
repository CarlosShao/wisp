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
    加 13 发原始 `-v` 日志（`*.v.log`）＋容器头（`*.head.txt`）＋名册三件套（`*.run/.fail/.skip/.colour.txt`）
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
  **13 发无一命中 97／98／95，13 发全 `BUILD_RC=0` ＋ `VET_RC=0`**〔独立复现〕。
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
