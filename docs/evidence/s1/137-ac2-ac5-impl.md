# 137 AC#2／AC#5 —— 实现方第二程：收紧断言（含 108 那两枚内联腿）＋门禁

裁决方：`worker-ticket137-ac2`。本程做两格：**AC#2**（把"被拒的那一个路径就是本用例自己种下去的那一个"写进断言，
范围按 `R-137-2` 含 108 那两枚不走 helper 的内联断言）与 **AC#5**（五组门禁——AC#2 产生了 `.go` 改动，这一格这次才有对象）。
**AC#2／AC#5 两勾均未翻**（票面 AC 归非实现者落笔）。AC#4（换根）是另一格，本程一个字未动。

## §0 锚点与树

- 开工首读（本程第一步，自量）：`git rev-parse --short HEAD` = **`f53ad5c`**
  （全 sha `f53ad5c2c9cdb4786b3c2191756056558594b0aa`，枚举行＝`test(136,AC#12②): 把终裁方那两发探针各钉一枚用例…`）。
  **本程全部读数都钉在这枚自量锚上**；派单时估的"约 `f53ad5c`"与实测同值，`worker-ticket136-ac12-ac13` 的两枚
  （`5c1529a`／`f53ad5c`）已在锚内。**这个自量值是唯一权威锚点**，下面任何 sha 都不替代它。
- ⚠ **凡从别人的报告／工具输出里取 sha，本程逐枚 `git cat-file -t` 现核过存在**（`docs/reports/injection-timeline.md` §11 那串假 sha 的教训）：

  | sha | 来源 | `git cat-file -t` 实测 |
  | --- | --- | --- |
  | `f53ad5c` | 本程自量锚点 | `commit` |
  | `4a0d7a4` | 终裁方 `137-ac1-r1-acceptance.md` §0 的开工锚 | `commit` |
  | `99ac876` | 终裁方末枚（票面"13 枚 commit `ee0a169`…`99ac876`"的尾） | `commit` |
  | `1d38206` | 测量方 `137-ac1-teeth-or-not.md` §0 的锚 | `commit` |
  | `6f3817a` | 本程 AC#2 那枚码提交（自量 `git log --oneline -1`） | `commit` |

- **不写任何耗时／间隔**：本机墙钟在 09-24 凌晨跳过 8h32m（终裁方三法互证，票面更正块末段），
  故本程每个时间戳只对它自己那一行成立；且本程**一枚计时类断言都没跑**（D32 的 CPU≤0.5%／RSS≤25MB 一字节未动、未读）。
- 变异与读数全部落在**仓外纯净树**（`git archive f53ad5c… | tar -x -C /d/tmp/wisp137ac2-<名>`，**只建不删**；
  仓内未建 worktree、未 checkout、未 `git stash`）。四棵：

  | 树 | 内容 | 用途 |
  | --- | --- | --- |
  | `wisp137ac2-tree-0` | 纯锚点（生产码＋测试件全是 `f53ad5c` 版） | `sh scripts/d22scan.sh` 的"锚点侧各 scope 数"对照 |
  | `wisp137ac2-tree-d-orig` | 锚点＋**MUT-D**（三枚测试件仍是锚点版＝未收紧） | 复现"软链形 11 枚全绿"那本害（自证的第一态） |
  | `wisp137ac2-tree-d-tight` | 锚点＋**MUT-D**＋我收紧后的三枚测试件 | 自证主发：收紧后软链形 11 枚必须转红 |
  | `wisp137ac2-tree-tight` | 锚点＋我收紧后的三枚测试件，**无变异** | （建好即用作 overlay 指纹对照；未变异读数走挂载工作树那两发） |

  ⚠ "每发先无条件 restore 再打下一发"在这一程的形状：每一态**一枚新树**，树在任一读数之前一次性建好＋打变异＋核指纹，
  此后不再改动任何一棵；三态之间唯一差异＝**是否 overlay 那三枚测试件**（生产码 md5 在两棵变异树之间逐字相同，见 §1 表），
  所以"上一发没还原干净"这类污染在构造上不可能。
- overlay 的字节同一性（不是我抄的，是 `md5sum` 两侧对点）：working tree 与 `tree-d-tight`／`tree-tight` 内
  `placement_symlink_113_other_test.go` `cb8350cb5e5917a685c10b39d57d5275`、
  `ancestor_separator_108_other_test.go` `20c141df792de655408138cdaebf633b`、
  `placement_leaf_118_other_test.go` `50d4e54c72167d969b361c092cb3e644` **逐枚相同**；
  同一棵树的 overlay 前（＝锚点版）是 `ce181f23…`／`0ae03ee4…`／`b84c85c8…`（CRLF 版，行尾差同 §1 的 go.mod 那 28 字节一说）
  ⇒ 差异只可能来自我那一版断言。
- 生产码在两棵变异树与挂载工作树之间：`winsec_other.go` = `b5056918be4ed13817d236fbcae0f477`、
  `winsec.go` = `a6144c880de80e43bb1393f3624e7221` —— 与终裁方 §1 报的两枚 md5 **逐字相同**（⇒ 本程与 AC#1 那两程读的是同一版底线码），
  MUT-D 之后两枚变成 `958963625351b269dfc04f45d76610d1`／`8548aeb16a12b6ca55390a3b8c3ba971`（两棵变异树同值）。
- 台件全在仓外 `D:\tmp\wisp137ac2-io\`：`mutate.py`（MUT-D，逐处锚串唯一性断言，count≠1 直接 exit）、
  `make-trees.sh`（建树＋overlay＋落地 grep）、`run.sh`（容器内硬闸 97/98/95）、`go.sh`（宿主侧 docker 驱动）、
  `parse.py`（把 `-v` 日志程序化拆成逐名颜色／名册／四数，含 `-count=N` 的逐名跨复读颜色漂移检查），
  八发原始日志 `*.v.log` 与容器头 `*.head.txt` 均在盘上，可复核。

## §1 AC#2 改了什么 ＋ MUT-D 自证三态

### §1.1 改动（三枚文件、四处判据点；生产码零字未动）

1. `internal/winsec/placement_symlink_113_other_test.go`
   - `:131` 新增 `const linkAttributionMarker137 = " through the link at "` —— 这枚短语是**两形拒因共有的**那一截
     （`winsec_other.go:158` 的 seal 侧文案与 `winsec.go:262` 的 unlink 侧文案都带它）。
   - `:145` 新增 `refusalCreditsLink137(err error, planted string) bool`：只在**标记之后**取被记名的那个前缀，
     要求它以 `planted` 开头**且**紧跟 `,`（或到此结束）——`/r137link` 是 `/r137link/w137ac2tmp/…/root/link` 的字符串前缀，
     不作这一步就会把宿主的链接读成本用例的。
   - `:174` `assertRefused113` 增第四参数 `planted string`；`:186` 在既有两条判据（必须拒／必须名 `ErrUnresolvedPath`）之后
     新增第三条 `AC#1 RED … does not credit the link this case planted`。**既有两条一字未放宽**。
   - 调用面五处全部传本用例自己种的那枚链接：`:204`、`:237`（depth 子测共用）、`:254`、`:271`、`:296`；
     其中 `:296` 那枚 `x\y` 用例原本把 `linkTo113` 的返回值丢掉，现提成 `link := linkTo113(t, root, `x\y`, f.dir)` 再用。
2. `internal/winsec/placement_leaf_118_other_test.go`（helper 的另一批调用者，R-137-2 点的"118 两处"）
   - `:81`、`:107` 各传 `link`（叶子链接本身）为 `planted`。
3. `internal/winsec/ancestor_separator_108_other_test.go`（**R-137-2 要求的范围**：这两枚分母腿不走 helper，内联断言）
   - 第一枚 `TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink`：`:79` 的 `t.Logf` 带上 planted，
     `:84` 新增 `} else if !refusalCreditsLink137(err, link) {` ⇒ `:90` 的 `AC#2 RED` 红行。
   - 第二枚 `TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor`：`:140` 把 `filepath.Join(root, linkName)` 提成
     `link` 变量（原来只喂给 `os.Symlink` 就丢了），`:146` Logf、`:151` 新增分支、`:156` 红行。
   - 两条 `else if` 挂在既有"必须拒"＋"必须名 `ErrIsReparsePoint`"之后，**原有两条判据一字未改**。

只改 helper 会漏掉 108 那两枚 ⇒ 本程按 `R-137-2` 把三处判据点全覆盖；AC#3 要的"11 枚软链形全转红"的**分母闭合**由 §1.2 的名册给出。

### §1.2 三条禁令的自证

- **零新增依赖**：`git diff --numstat f53ad5c 6f3817a` 只有三枚 `*_test.go`（`go.mod`／`go.sum` 未动）；
  新代码只 import 标准库 `strings`（同包既有测试件早已用它，如 `reparse_windows_test.go:159`）。多行实现合法（`R-137-5` 已作废"必须一行"）。
- **没放宽成"两形都算过"**：新判据在两形用**同一句**要求（记名者必须等于本用例种的链接）。后果正是软链形那 11 枚变红，见 §1.3。
- **没拿被测函数算期望值（`R-119-9`）**：期望值 `planted` 是 fixture 侧 `filepath.Join`／`linkTo113` 的返回值；
  被测侧的产物（错误串）只作为**被读的那一个**进入比较，且比较用的短语 `" through the link at "` 是我从生产文案里
  **抄下来的字面量**（`winsec_other.go:158`、`winsec.go:262`），不是调用被测函数算出来的。
- ⚠ 顺带排掉一枚**新假绿**：`strings.Contains(err.Error(), planted)` 在软链形**恒真**——被拒的整条拼写里必然含着
  planted 的子串（`/r137link/…/root/link/keep-me.txt` 就写着 `/r137link/…/root/link`）。所以判据只看标记之后的那一截。

### §1.3 MUT-D 自证三态（本程自己落、自己取名册；红名逐名）

MUT-D 形状按票面更正块与 `R-137-1`：两条走查各只走前 3 个组件（`platformVerifyPlacement` 与 `firstLinkAncestor` 各一处），
⇒ 仍拒宿主自己那枚 `/r137link`，**永远走不到用例自种的链接**。
**只用了 MUT-D 这一发**：`R-137-1` 明令 AC#3 不许钉在单发 A 或 B 上（A 只响 9、B 只响 2，收紧后另有几枚恒不响），
本程因此**没造也没跑** MUT-A／MUT-B／MUT-AB，不自证一条我不打算钉的尺。

容器 `golang:1.27`（容器内 `go1.27.1 linux/amd64`），挂载 `MSYS_NO_PATHCONV=1` ＋ `/d/...`，
每发先 `ls -l /src/go.mod` ＋ 字节数硬闸（工作树 855／快照树 883——差 28＝`.gitattributes` 的 `* text=auto` 行尾差，与终裁方 §1 同一解释），
形状硬断言（软链形建 `/r137link -> /r137priv` 并 `readlink` 逐字核；普通形断言 `/r137link` 根本不许存在），
**八发无一命中 97／98／95**；每发读数前先 `go build ./...`＋`go vet ./internal/winsec/`，**八发全 `BUILD_RC=0`＋`VET_RC=0`**。

| 发（树·形·`-count=1 -v`） | RUN | 顶 PASS/FAIL/SKIP | 子 PASS/FAIL/SKIP | rc | 分母 11 枚 |
| --- | --- | --- | --- | --- | --- |
| `dorig-plain-1` 锚点断言＋MUT-D·普通 | 52 | 17/13/0 | 17/5/0 | 1 | **11 枚 FAIL** |
| `dorig-link-1` 锚点断言＋MUT-D·软链 | 45 | 24/3/3 | 15/0/0 | 1 | **11 枚 PASS（零检测力＝本票钉的害，我在自己树上复现）** |
| `dtight-plain-1` 收紧后＋MUT-D·普通 | 52 | 17/13/0 | 17/5/0 | 1 | 11 枚 FAIL |
| `dtight-link-1` 收紧后＋MUT-D·软链 | 45 | 17/10/3 | 11/4/0 | 1 | **11 枚 FAIL（逐名，见下）** |
| `work-plain-1` 收紧后·**无变异**·普通（挂载工作树） | 52 | 30/0/0 | 22/0/0 | **0** | 11 枚 PASS |
| `work-link-1` 收紧后·**无变异**·软链（挂载工作树） | 45 | 20/7/3 | 11/4/0 | 1 | **11 枚 FAIL** |

**收紧前后唯一差异＝软链形那一列**：`dorig-link-1` 45/24/3/3＋15/0 → `dtight-link-1` 45/17/10/3＋11/4，
差值逐名核过（`parse.py`＋集合差，非手抄）：**顶 FAIL 3→10、子 FAIL 0→4，多出的 11 枚正好等于分母**，
`only-in-orig` 侧为空（没有任何一枚由红转绿／由绿转 SKIP），RUN 名册两版逐名相同（45 枚）。
**普通形八数与逐名 FAIL 名册在收紧前后完全相同**（`52/17/13/0＋17/5/0`、`rc=1`，对称差集为空）
⇒ 与终裁方给的参照值一致，且这八个字**不是**我这一版断言造出来的。

`dtight-link-1` 软链形 11 枚红名逐名（`--- FAIL:` 行原文取名，四枚子测各算一枚）：

1. `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone`
2. `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth`
3. `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth/link-at-depth-1`
4. `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth/link-at-depth-2`
5. `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth/link-at-depth-3`
6. `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth/link-at-depth-4`
7. `TestAC1POSIXSealDirThroughASymlinkRefuses`
8. `TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing`
9. `TestAC1POSIXSealFileThroughABackslashNamedLink`
10. `TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink` ← **108 内联腿**
11. `TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor` ← **108 内联腿**

**SKIP 逐发**（"没响"不等于"被跳过"这条本程也逐名给）：普通形三发（`-count=1`）恒 **0 枚**；软链形三发恒 **3 枚**、名册逐发同一组
（`TestAC2POSIXSeamProbeShapesAreBuiltOnAResolvedRoot125`、`TestAC2POSIXSeamGuardStillRefusesEveryHostileShape125`、
`TestAC2POSIXSeamAcceptsTheHonestPOSIXAnswer125`）＝票 125 那三枚自拒探针，与基线同形；
**分母 11 枚在八发里无一枚是 SKIP** ⇒ `dorig-link-1` 那 11 枚绿是跑出来的绿，`dtight-link-1` 的红是断言判出来的红。

**对照组（证 MUT-D 真落地）**：`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`、
`TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed`、`TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed`
三枚在 `dorig-plain-1`／`dorig-link-1`／`dtight-plain-1`／`dtight-link-1` **四发全红**（`dtight-link-1` 顶 FAIL 10 枚里就含它们，无第四因）
⇒ 变异打到了东西，那 11 枚在 `dorig-link-1` 的绿不是"没打到"。
按 `R-137-3` 我只用这三枚作对照，**没有**把 `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119` 算进去，并把它单独量了一列：
`dorig-plain` FAIL／`dorig-link` PASS、`dtight-plain` FAIL／`dtight-link` PASS（`TestAC2POSIXInjectedTestDataDirStandsAsDeclared119` 同形）
⇒ 与终裁方那句"它自己也用 raw `t.TempDir()`、与分母同病"逐发相同，且**我的收紧没动它一根手指**（两版同色）。

### §1.4 一条要主动报给下游的分歧（不是失败，是 AC#4 的账提前露出来）

未变异的收紧态在软链形也是 **11 枚红**（`work-link-1`：45/20/7/3＋11/4、`rc=1`），因为那 11 枚的根仍是 7 处 raw `t.TempDir()`，
宿主的链接在走查第 2 个组件就把拼写拒掉了、根本走不到用例种的链接。⇒ 这不是"坏尺"：
**同一形里根已解析的三枚对照（118×2＋119-AC3）在新判据下仍 PASS**，说明"记名者＝自己种的链接"在软链形可满足；
要让它满足，得先把那 7 处根换成已解析形＝**票面 AC#4 那一格**。按更正块那句"换根与收紧断言各自都能单独成立、也各自都能单独漏 ⇒ 不许并成一格"，
本程**没有**顺手做 AC#4，只把这条读数留给它。宿主 CI 侧无影响：`internal/winsec` 的 CI 步是 Windows 任务（`.github/workflows/ci.yml:345-386`
→ `scripts/winsec-tests.sh`），这三枚文件带 `//go:build !windows`，在那一步里根本不参与编译；软链形只存在于我们的容器台件里。

### §1.5 我自加的一发探针 MUT-E（票面没要求，为的是排除"新判据会静默哑掉"这一形）

新判据读的是生产文案里那截 `" through the link at "`。会哑吗？**MUT-E**：把两处错误的措辞改掉
（`winsec_other.go:158` 与 `winsec.go:262` 改成 `via an ambient link: …`），**拒还是拒、错误身份还是 `ErrUnresolvedPath`／`ErrIsReparsePoint`、
模式与属主的效果一字未动**，只把标记抹掉（`grep -c " through the link at "` 在两枚文件上都是 **0**，落地 grep 见 `tree-e-tight`）。
同一棵树上 overlay 我收紧后的测试件，两形各一发 `-count=1 -v`：

| 发 | RUN | 顶 PASS/FAIL/SKIP | 子 PASS/FAIL/SKIP | rc | 读什么 |
| --- | --- | --- | --- | --- | --- |
| `e-plain-1`（MUT-E·普通形） | 52 | 21/9/0 | 18/4/0 | 1 | **分母 11 枚全 FAIL**＋118 那两枚也 FAIL（它们走同一支 helper），合计顶 7＋2＝9、子 4；`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`（不经我这条判据）仍 PASS |
| `e-link-1`（MUT-E·软链形） | 45 | 18/9/3 | 11/4/0 | 1 | 同一组名册，颜色不变（软链形本来就红，见 §1.4） |

⇒ 生产文案一旦挪动，这条尺**响成一片**而不是悄悄放行（`refusalCreditsLink137` 取不到标记就返回 false＝fail-closed）。
反面形状也读得到：若我按省事写法用 `strings.Contains(err.Error(), planted)`，`e-plain-1` 里那句被拒的整条拼写仍含着 planted 子串
⇒ 那 11 枚会**继续绿**、而且是在标记已经不存在的情况下绿——这就是本程宁可比对"标记之后那一截"的原因。
（这一发的代价也说清：MUT-E 让 118 两枚一起红＝**三枚对照组全成不了照**，它只能用来判"哑不哑"，不能当落地凭据用。）

## §2 108 那两枚内联腿**单独**给读数（`R-137-2` 的那 2 枚，不走 helper）

终裁方给的位是 `ancestor_separator_108_other_test.go:80-84` 与 `:138-142`（锚点版）。我自己复核过：`f53ad5c` 上
`:80-84` 与 `:138-142` 确实是两枚 `if err == nil { … } else if !errors.Is(err, winsec.ErrIsReparsePoint) { … }` 的内联断言，
**一处都不调用 `assertRefused113`**（grep 全仓：helper 调用面只有 113 五处＋118 两处）。收紧后的行号（本程改完）：
第一枚 `:79` Logf／`:84` 新分支／`:90` 红行；第二枚 `:140` 提 `link` 变量／`:146` Logf／`:151` 新分支／`:156` 红行。

### §2.1 逐发颜色（这两枚单独一行，不混在 §1.3 的大表里）

| 发 | `TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink` | `TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor` |
| --- | --- | --- |
| `dorig-link-1`（MUT-D·软链·**锚点断言**） | **PASS** | **PASS** |
| `dtight-link-1`（MUT-D·软链·收紧后） | **FAIL** | **FAIL** |
| `dorig-plain-1`（MUT-D·普通·锚点断言） | FAIL | FAIL |
| `dtight-plain-1`（MUT-D·普通·收紧后） | FAIL | FAIL |
| `work-plain-1`（无变异·普通·收紧后） | PASS | PASS |
| `work-link-1`（无变异·软链·收紧后） | FAIL | FAIL |
| `work-plain-2`／`work-link-2`（`-count=2`，各两遍） | PASS×2／FAIL×2 | PASS×2／FAIL×2 |

⇒ 这两枚**只在收紧那一刀上**由 PASS 翻 FAIL（`dorig-link-1` → `dtight-link-1`，同一棵变异树、只 overlay 测试件），
普通形两版同色（MUT-D 下 `err=<nil>`，红在第一支既有判据上，与我的新分支无关）＝**新增判据没有把普通形推红一寸**。
`-count=2` 的两遍逐名同色（`parse.py` 的 `REPETITION_DRIFT=0`）。

### §2.2 红行的归属：这两枚红在**自己的行上**，不是红在 helper 里

`dtight-link-1.v.log` 里 `grep -oE "ancestor_separator_108_other_test.go:[0-9]+: AC#2 RED"` 命中恰两枚：

```
ancestor_separator_108_other_test.go:90:  AC#2 RED: RemoveUnlinked("/r137link/w137ac2tmp/TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink4055766027/002/root/link/keep-me.txt") refused with winsec: entry is a link to something else, not private data …: the spelling reaches it through the link at /r137link, and whatever lives behind that link is not this tree's data to delete, which does not credit the link this case planted at "/r137link/w137ac2tmp/…/002/root/link": an ambient link above the tree can answer for the refusal while the guard under test never reaches this spelling
ancestor_separator_108_other_test.go:156: AC#2 RED: RemoveUnlinked("/r137link/w137ac2tmp/TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor3982944564/002/root/x\\y/keep-me.txt") … the link at /r137link … does not credit the link this case planted at "/r137link/…/002/root/x\\y" …
```

同一发里 113 族那 9 枚的红分别落在 `placement_symlink_113_other_test.go:204/237(×4 子测)/254/271/296`（＝helper 的五个调用点）。
⇒ 若只按票面字面改 helper，这两枚**没有红可落**（它们的判据在 108 文件里），AC#3 的"11 枚全转红"就必缺 2 枚——这条现在盘上有对点了。

### §2.3 收紧前的"绿得没有理由"在这两枚上读得到（同一发的原文）

`dorig-link-1.v.log:61`（锚点断言版，`t.Logf` 在 `:79`）：

```
ancestor_separator_108_other_test.go:79: RemoveUnlinked("/r137link/…/root/link/keep-me.txt") -> err=winsec: entry is a link to something else, not private data /r137link/…/root/link/keep-me.txt: the spelling reaches it through the link at /r137link, and whatever lives behind that link is not this tree's data to delete
```

这一发该用例 **PASS**，而它拿到的拒记名的是宿主的 `/r137link`（`firstLinkAncestor` 被 MUT-D 截到前 3 个组件，压根没看 `root/link`）。
对照正向：`work-plain-1.v.log` 同一条 Logf 的尾部是

```
… the link at /r137plain/w137ac2tmp/TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink1063965065/002/root/link, and whatever lives behind … (link this case planted: "/r137plain/w137ac2tmp/…/002/root/link")
```

记名者＝自己种的那枚 ⇒ 新分支放行。⇒ 这条尺在普通形是**真被满足**而不是"永远不满足"（恒不满足＝`R-137-1` 那类坏尺）。

## §3 AC#5 五组门禁（原文命令＋读数；对象＝有 `.go` 改动之后的本程工作树）

工作树内容＝`6f3817a`（AC#2 那三枚 `*_test.go`）＋`7d73b5f`／`6fc7f3d`（本文件），生产码＝锚点版
（`winsec_other.go` `b5056918…`、`winsec.go` `a6144c88…`，§0）。⚠ 在飞兄弟的地界本程未碰：
`git diff --cached --name-only` 每枚提交只有我自己的路径，`internal/observe/**` 与 `cmd/wisp/**` 零字未动
（`git status --porcelain` 全程只有那三枚 `internal/winsec/*_test.go` ＋那枚来源未明的未跟踪件——它未读、未提交、未改、未删）。

### §3.1 `gofmt -l`

- 仪器版本：`go version go1.27.1 windows/amd64`（与 `go.mod` 的 `toolchain go1.27.1` 同版）。
- 整包（票面要求的那一枚）：`gofmt -l ./internal/winsec/` → **零行输出、rc=0**。
- 顺手加一层全仓：`gofmt -l .` → **零行输出**。
- 我这轮唯一动过的格式面就是那三枚文件；两枚被测文件在锚点上本来就干净，改完仍干净。

### §3.2 `gofumpt -l`（版本写明）

- 用的就是宿主**现成**那枚：`D:\work\base\gopath\bin\gofumpt.exe --version` → **`v0.12.0 (go1.27.1)`**。
  ⚠ **没有跑** `go install mvdan.cc/gofumpt@latest`（那会把宿主工具升掉，是这仓刚记下的坑；票面 AC#5 那句"CI 装 `@latest`、版本没钉"我只如实记版本，不去动 CI 的那一步）。
- `gofumpt -l ./internal/winsec/` → **零行输出、rc=0**；`gofumpt -l .`（全仓）→ **零行输出、rc=0**。

### §3.3 `go vet` 双 GOOS（逐错误行归因，不整树 rc=1 就归给工具链）

| 发 | 命令（逐字） | rc | 输出 |
| --- | --- | --- | --- |
| 宿主原生（windows＝GOOS 默认） | `go vet ./...` | **0** | 零行 |
| 宿主交叉（linux）·整树 | `GOOS=linux GOARCH=amd64 go vet ./...` | 1 | **3 行**，只有 `cmd/wisp` 那一枚导入链 |
| 宿主交叉（linux）·**逐包**（33 枚逐枚跑） | `GOOS=linux GOARCH=amd64 go vet ./<pkg>` ×33 | 31 枚 rc=0／**2 枚 rc=1** | 见下两条归因 |
| 容器原生（linux，真类型读数） | `docker run … golang:1.27 go vet ./...`（`CGO_ENABLED=1`） | **0** | 零行 |

- 逐包那一发的账本在 `/d/tmp/wisp137ac2-io/vet-linux-cross-perpkg.txt`（33 行，每行 `PKG=… RC=… OUTPUT_LINES=…`）：
  **31 枚 rc=0 且逐枚 `OUTPUT_LINES=0`**（不是"没输出因为整树早退"，是每枚自己跑完的零输出）。
- 两枚 rc=1 的归因，各追到 `file:line`／包约束本身：
  1. `cmd/wisp`：`package github.com/CarlosShao/wisp/cmd/wisp` → `imports github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx`
     → `imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files in
     D:\work\base\gopath\pkg\mod\github.com\k2-fsa\sherpa-onnx-go-linux@v1.13.8`。
     归因＝交叉编译时 `CGO_ENABLED=0`，那枚 linux+cgo 包的文件全被构建约束排除 ⇒ **不是本仓的类型破口**；
     同一条导入链在容器原生那一发里 rc=0（`CGO_ENABLED=1` 实测打印为 1）。
  2. `cmd/balldebug`：`package github.com/CarlosShao/wisp/cmd/balldebug: build constraints exclude all Go files in
     D:\work\workspace\projects plans\Wisp\cmd\balldebug`。错误文本点到**自家仓库路径**，按纪律我没放过：
     逐文件 `head -5 cmd/balldebug/*.go` 量到 `diff_windows.go`／`main.go`／`shot_windows.go` **三枚全带 `//go:build windows`**
     ⇒ 这是一枚按设计只在 Windows 建的调试命令（票 07 的台件），GOOS=linux 下"无文件可建"是它自己的约束，不是回归。
     它也不是本程改动碰到的包（`git diff --numstat f53ad5c HEAD` 零枚 `cmd/**`）。
- ⚠ 一句给下游的口径：**整树那一发 rc=1 只报出 1 枚**（`cmd/wisp` 的链），`cmd/balldebug` 那枚要逐包才看得见 ⇒ "整树红＝只有一处"是错觉。
  派单里那句"失效面只有两枚 cmd 包"在本程**逐包粒度复现成立**（第三枚 `cmd/llmrecord` 逐包 rc=0）。
- 我这轮真正需要的类型读数在容器那一列：三枚改的文件都带 `//go:build !windows`，只有 Linux 原生 vet 会**看到它们**（rc=0）；
  另外 §1.3 那八发每发读数前的 `go vet ./internal/winsec/` 也全 rc=0。

### §3.4 `-count=2 -v` 两形四数（＋四数之外的名册差集）

| 形·发 | RUN | 顶 PASS/FAIL/SKIP | 子 PASS/FAIL/SKIP | rc |
| --- | --- | --- | --- | --- |
| 普通 `work-plain-2`（`go test -count=2 -v ./internal/winsec/`，容器） | 104 | 60/0/0 | 44/0/0 | **0** |
| 软链 `work-link-2`（同左，`TMPDIR` 在 `/r137link` 之后） | 90 | 40/14/6 | 22/8/0 | 1 |
| 参照：普通 `work-plain-1`（`-count=1`） | 52 | 30/0/0 | 22/0/0 | 0 |
| 参照：软链 `work-link-1`（`-count=1`） | 45 | 20/7/3 | 11/4/0 | 1 |

- 四数＝每形那四个数（RUN／顶／子／rc），逐名另算：`-count=2` 两发都是 `-count=1` 的**整两遍**（104＝52×2、90＝45×2），
  `parse.py` 的 `REPETITION_DRIFT=0`＝**没有任何一枚在两遍之间换色**；`RUN` 名册与 `-count=1` 逐名同集（普通 52 枚恒常、软链 45 枚恒常）。
- **SKIP 逐名**：普通形两发恒 **0 枚**；软链形两发恒 **6＝3×2 枚**，名册只有票 125 那三枚自拒探针
  （`…SeamProbeShapesAreBuiltOnAResolvedRoot125`／`…SeamGuardStillRefusesEveryHostileShape125`／`…SeamAcceptsTheHonestPOSIXAnswer125`）；
  **分母 11 枚在两发里都不是 SKIP** ⇒ "判据变严之后这 11 枚由绿转 SKIP"这条退路被读数排除。
- **FAIL 名册**：普通形两发全空；软链形两发恰是那 11 枚（§1.3 逐名表），**没有第 12 枚**。
- ⚠ **与派单给的参照值对不上的是"标签"，不是数**（两边都留）：派单把 `52/17/13/0＋17/5/0、rc=1` 写成 AC#5 `-count=2` 普通形的参照值，
  而这八字在终裁方自己那份证据里是 **§5④ 的 MUT-D·普通形** 读数（未变异的普通形基线是 `52/30/0/0＋22/0`，票 124 批次 3b 与终裁方十发都是这个数）。
  我按两种读法都量了：**MUT-D·普通形**（`dorig-plain-1` 与 `dtight-plain-1`）实测 `52/17/13/0＋17/5/0、rc=1` ＝与参照值**逐数相同**，
  且逐名 FAIL 名册两版相同、也与终裁方那句"与 MUT-AB·普通全同"相符；**未变异·普通形 `-count=2`** 是 `104/60/0/0＋44/0/0、rc=0`（＝基线整两遍）。
  ⇒ 分歧只在"这八个字属于哪一发"，数值本身零处不一致。

### §3.5 一次全仓仪器 `sh scripts/d22scan.sh`

- 命令逐字（仓根，Git Bash）：`sh scripts/d22scan.sh` → **`D22_RC=0`**；日志 `/d/tmp/wisp137ac2-io/d22scan-worktree.txt`。
- 步 1（种子违规阳性对照，走 `tools/d22scan/runtests.sh`，强制 `-count=1`）：
  `runtests.sh: OK - packages=[./...] top-level: PASS=21 FAIL=0 SKIP=0, === RUN=31, '[no tests to run]'=0`
  ⇒ 对照仍活着：这轮"clean"不是"仪器瞎了"。
- 步 2（真扫）各 scope 工作量，与**同一台器在纯锚点树**（`tree-0`＝`git archive f53ad5c | tar -x`，未 overlay 任何东西）
  上跑出来的那一发逐 scope 对点：

  | scope | `tree-0`（纯锚点） | 工作树（AC#2 之后） | 差 |
  | --- | --- | --- | --- |
  | `bans #1-5 internal/` | 203 | 203 | 0 |
  | `bans #1-5 cmd/` | 22 | 22 | 0 |
  | `ban #6 frontend/` | 40 | **43** | **＋3** |
  | `ban #7 internal/tools/` | 18 | 18 | 0 |
  | `ban #8 design/` | 16 | 16 | 0 |
  | `ban #8 frontend/` | 40 | **43** | **＋3** |
  | `ban #8 internal/` | 404 | 404 | 0 |
  | `ban #8 cmd/` | 39 | 39 | 0 |
  | 合计（`internal/`＋`cmd/` 生产 Go） | 225 | 225 | 0 |

  ⇒ **各 scope 一枚不降**（无一枚变小＝没有任何 scope 被静默清空）。两处 ＋3 归因到位（不是源码、不是任何人在飞的东西）：
  `frontend/dist/index.html`、`frontend/dist/assets/index-CEH-Pz8P.css`、`frontend/dist/assets/index-UL9kYvYl.js`
  ＝本机 `npm run build` 的遗留产物，`git check-ignore -v` 三枚逐枚判给 `frontend/.gitignore:12:dist/*`；
  两侧的 `find frontend -type f` 差集里除这三枚外只剩 `frontend/node_modules/**`（扫描器不把它计入 text 工作量，故 43 而非数千）。
  同一枚 ＋3 票 77 那程也量过（`docs/evidence/s1/ci-runner-readings-2026-09-21.md:228-245`），我这里独立重取、不引它的数当自己的数。
- 扫描器自报的 clean 行原文（含逐 scope 的 live work 一串）在日志末行，未截。

