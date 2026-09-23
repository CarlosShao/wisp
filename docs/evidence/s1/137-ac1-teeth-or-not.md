# 137 AC#1 —— winsec 那批 POSIX 拒绝腿在软链形下到底有没有牙（变异取证，不改生产码）

裁决方：`worker-ticket137-ac1`。本格**只做一格**：AC#1「先定成不成立」。生产码零改动，动的全是仓外快照副本里的测试面变异对象。

## §0 锚点与漂移

- 开工首读 `git rev-parse --short HEAD` = **`1d38206`**（全 sha `1d382064b40c6e8de787f37374dfe8e4dd764e30`）。
  **本程全部读数都在这枚锚上量**：五台快照一律 `git archive 1d38206 | tar -x -C /d/tmp/wisp137-ac1-<名>`，仓内未建 worktree、未 checkout。
- 写件时 HEAD 已漂到 `e776224`（兄弟 `worker-ticket133-ac2-fix`／`worker-ticket136-ac8-ac9` 在推）。
  已核 `git diff --name-only 1d38206 e776224 -- internal/winsec/` **输出为空** ⇒ 漂进来的 10 个路径全在
  `.scratch/**`、`docs/**`、`internal/observe/**`，对本格读数**零影响**。票 137 票面在 `1d38206..HEAD` 之间无人动过（`git log -- .scratch/wisp/issues/137-*.md` 空）。

## §1 容器与挂载证明（防"静默挂空 rc=0＝假绿"）

- Docker 可用：`docker info` ⇒ ServerVersion **29.6.2**，`linux/x86_64`。宿主是 Windows，这批用例带 `!windows` 标签 ⇒ **宿主没有分母**，读数只在容器。
- 镜像 `golang:1.27`，容器内 `go version` = **`go1.27.1 linux/amd64`**（与票 124 批次 3b 同字）。
- 挂载用 `/d/...` 形式 + `MSYS_NO_PATHCONV=1`；**每一发**跑前先 `ls -l /src/go.mod`，十发逐字相同：

  ```
  -rwxrwxrwx 1 root root 883 Sep 23 14:34 /src/go.mod
  go version go1.27.1 linux/amd64
  ```

  883 字节＝真件（挂空会读不到文件或读到 0 字节），脚本里另有一枚 `[ "$(wc -c < /src/go.mod)" = "883" ] || exit 97` 的硬闸。
- 每发另外自证被测文件指纹（十发只有两类值）：
  - 基线/树 `tree0`：`winsec_other.go` md5 `b5056918be4ed13817d236fbcae0f477`、`winsec.go` md5 `a6144c880de80e43bb1393f3624e7221`
  - 变异树各自只改一枚文件，指纹见 §4。
- 形状硬断言（沿用票 124 的 `exit 97/98/99`）：
  - 软链形：`ln -s /realpriv /varlink` + `[ -L /varlink ]` + `readlink` 逐字核 + `ls -ld` 打印
    `lrwxrwxrwx 1 root root 9 ... /varlink -> /realpriv`，`TMPDIR=/varlink/w137tmp`；
  - 普通形：`[ ! -e /varlink ] || exit 99`（宿主的链接根本不许存在）+ `[ ! -L /plainroot ] || exit 99`，`TMPDIR=/plainroot/w137tmp`。
  - **十发无一发命中 97/98/99**。
- 模块缓存：容器 `GOPROXY=off` + 预置 named volume（`wisp137ac1-gomod`／`wisp137ac1-gobuild`）。
  这条不是偷懒：第一轮我按默认 GOPROXY 跑，`go build ./...` 卡 15 分钟只推进 2 MB（本机 fake-IP 代理到 `proxy.golang.org` 在抖），
  **卡住的取数会被误读成"跑通了"** ⇒ 改成离线取缓存，缺件就直接 fail，不存在"网络慢装出来的绿"。
- 跑法脚本：`/d/tmp/wisp137-ac1-io/run.sh`（+ 表生成 `table.sh`、归因 `attrib.sh`）；十发的原始 `-v` 日志与 summary 都在
  `/d/tmp/wisp137-ac1-io/*.v.log`、`*.summary.txt`。快照目录 `wisp137-ac1-tree0`／`-tree-muta`／`-tree-mutb`／`-tree-mutab`／`-tree-mutd` **一律只建不删**。

## §2 真实枚数裁决：是 **11**（7 顶层 + 4 子测试），票面标题/正文那个「12＝8 顶层＋4 子测试」写错了

我自己逐枚数的，没沿用简报里任何一个数字。判据＝**票 113／票 108 的 POSIX 拒绝腿里，那些 `root` 仍由未解析的 `t.TempDir()` 供给的枚**，
也就是票 124 批次 3b 说的那「7 处按地界有意未换的 `t.TempDir()`」所覆盖的拒绝腿。7 处递根点逐处 `grep -n` 得到：

| # | 用例文件:行 | 递根点 | 它覆盖的拒绝腿 |
| --- | --- | --- | --- |
| 1 | `placement_symlink_113_other_test.go:146` | `root := filepath.Join(t.TempDir(), "root")` | `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone` |
| 2 | `placement_symlink_113_other_test.go:169` | 同上（在子测试里） | `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth` ＋ 它的 4 枚子测试 |
| 3 | `placement_symlink_113_other_test.go:197` | 同上 | `TestAC1POSIXSealDirThroughASymlinkRefuses` |
| 4 | `placement_symlink_113_other_test.go:213` | 同上 | `TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing` |
| 5 | `placement_symlink_113_other_test.go:238` | 同上 | `TestAC1POSIXSealFileThroughABackslashNamedLink` |
| 6 | `ancestor_separator_108_other_test.go:69` | 同上 | `TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink` |
| 7 | `ancestor_separator_108_other_test.go:127` | 同上 | `TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor` |

⇒ **7 处递根点 = 7 枚顶层 + 4 枚子测试 = 11 枚**（第 2 处一处覆盖 5 枚：顶层 1 ＋ 子测试 4）。
`assertRefused113` 真实位置：**`internal/winsec/placement_symlink_113_other_test.go:127`**（票面没写行号，简报让我自己 grep，别信记忆）。

**哪一处写错了**：票面标题与「为什么现在立案」表里的 **12＝8 顶层＋4 子测试** 是错的，**AC#1 正文那句「这 11 枚带 `!windows` 构建标签」才是对的**。
那个 8 的最可能来源我能复现但不能采信：`placement_symlink_113_other_test.go` 里**确实有 8 枚顶层 `func Test`**，
可其中 3 枚是 AC#3 的反向腿（`TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree`、
`TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks`、`TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator`，
它们要的是"成功"不是"拒"，且根已换成 `SealableTempDirForTest124`），把它们当拒绝腿就会得 8＋4＝12，
而真正的 108 那 2 枚顶层反而被漏掉。§5/§6 的逐名表就是按 11 枚给的，枚枚点名，**没有一个我没数过的数字**。

⚠ 同族**不在**本格分母里的邻居（它们的根已解析 ⇒ 天生有牙，我拿它们当变异落地对照组）：
`TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed`、`TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed`
（`placement_leaf_118_other_test.go:60` 走 `ownTree118`→`SealableTempDirForTest124`）、
`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`（`dataroot_symlink_119_other_test.go:226` 走 `cleanSpelling119`）。

## §3 变异前基线：两形各一枚，且「两形拒不同的东西」这一条**读得到**

`go test ./internal/winsec/ -count=1 -v`，纯净树 `tree0`，两形各一枚新容器：

| 形 | RUN | 顶层 PASS | 顶层 FAIL | 子测 PASS | 子测 FAIL | SKIP | 包级 rc |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 普通形 `TMPDIR=/plainroot/w137tmp` | 52 | 30 | 0 | 22 | 0 | 0 | 0 |
| 软链形 `TMPDIR=/varlink/w137tmp`（`/varlink -> /realpriv`） | 45 | 27 | 0 | 15 | 0 | 3 | 0 |

这两行与票 124 批次 3b 交件的「普通形 52/30/0/0＋22/0、软链形 45/27/0/3＋15/0」**逐数相同** ⇒ 我的锚与它的锚在被测面上同版本。
基线还跑过第二遍（写件时补的），六数与第一遍逐数相同。

逐名归因（`attrib.sh`，取每枚用例自己日志块里第一句 `reaches it through the link at X`）：

| 腿（11 枚） | 普通形它拒的那枚链接 | 软链形它拒的那枚链接 |
| --- | --- | --- |
| `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone` | `<TMP>/…/002/root/link` | **`/varlink`** |
| `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth`（顶层，随 depth-1 子测） | `<TMP>/…/002/root/link` | **`/varlink`** |
| `…/link-at-depth-1` | `<TMP>/…/002/root/link` | **`/varlink`** |
| `…/link-at-depth-2` | `<TMP>/…/002/root/dir1/link` | **`/varlink`** |
| `…/link-at-depth-3` | `<TMP>/…/002/root/dir1/dir2/link` | **`/varlink`** |
| `…/link-at-depth-4` | `<TMP>/…/002/root/dir1/dir2/dir3/link` | **`/varlink`** |
| `TestAC1POSIXSealDirThroughASymlinkRefuses` | `<TMP>/…/002/root/link` | **`/varlink`** |
| `TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing` | `<TMP>/…/002/root/link` | **`/varlink`** |
| `TestAC1POSIXSealFileThroughABackslashNamedLink` | `<TMP>/…/002/root/x\y` | **`/varlink`** |
| `TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink` | `<TMP>/…/002/root/link` | **`/varlink`** |
| `TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor` | `<TMP>/…/002/root/x\y` | **`/varlink`** |
| ── 对照组（不属本格分母）── | | |
| `TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed` | `<TMP>/…/002/root/artifact.txt` | `<REAL>/…/002/root/artifact.txt`（**自己种的**） |
| `TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed` | `<TMP>/…/002/root/artifact.txt` | `<REAL>/…/002/root/artifact.txt`（**自己种的**） |
| `TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119` | `<TMP>/…/001/data/out` | `<REAL>/…/001/data/out`（**自己种的**） |

⇒ **票面这一条 premise 成立**：同一支仪器在两形里拒的是**不同的链接**，11 枚在软链形全部由宿主的 `/varlink` 供给拒因，
而 `assertRefused113`（`:127`）两形都只 `errors.Is(err, ErrUnresolvedPath)` ⇒ 从断言本身分不出来。
对照组三枚在两形都点名自己种的链接，这条差异不是我推的，是日志里读出来的。

## §4 变异落地证明（先证落地，再读颜色）

变异树全部由 `tree0` `cp -a` 得来，**每发一棵新树、单发单因子，无叠发**；生产码在**仓内工作树**里一字未动（见 §7 的 git 面）。

| 代号 | 改了哪一步 | 被改后那一行原文（`grep -n` 摘要） | 文件指纹 |
| --- | --- | --- | --- |
| **MUT-A**（票面指定的那一发，seal 路） | `winsec_other.go:155` `platformVerifyPlacement` 的链接走查整段拿掉，函数体只剩 `return path, nil`（＝票 113 拆掉的那枚旧 stub） | `155:func platformVerifyPlacement(path string) (string, error) {`／`156-157: // MUTATION-137-AC1-A: the link walk is gone, this answers "yes, unchanged" for every spelling. …`／`158:	return path, nil` | `winsec_other.go` md5 `0e2a136cbb55dd9912dfada35b62214e`（`winsec.go` 仍 `a6144c88…`＝未改） |
| **MUT-B**（票面指定的那一发，unlink 路） | `winsec.go:279` `firstLinkAncestor` 的祖先走查拿掉，直接 `return ""`（＝退回票 103 PROBE F 量的"只看叶子"那形） | `279:func firstLinkAncestor(path string) string {`／`280-282: // MUTATION-137-AC1-B: the unlink route answers "no ancestor is a link" for every spelling …`／`283:	return ""` | `winsec.go` md5 `ae7c51bf0a067b017a8685dc31c6238d`（`winsec_other.go` 仍 `b5056918…`＝未改） |
| **MUT-AB** | A＋B 同时（＝"seal 与 unlink 两条路都不看链接"的全拆） | 两文件各自与上行逐字相同 | `0e2a136c…` ＋ `ae7c51bf…` |
| **MUT-D**（**加测**，半修那一发，不是票面指定的） | 两条走查都改成「只走前 3 个前缀」`if len(prefixes) > 3 { prefixes = prefixes[:3] }`：**仍拒宿主的 `/varlink`，但走不到用例自己种的链接** | `winsec_other.go:156-162`（`MUTATION-137-AC1-D` 注释＋截断），`winsec.go:284-288`（unlink 路同一刀） | `winsec_other.go` `2002cd7b128ec59d3c26486b77b10d03`、`winsec.go` `9134d7e9928955832b2bb2813af6614e` |

每发读数前都在**同一枚容器里**跑 `go build ./...` 与 `go vet ./internal/winsec/`，十发全部 **`BUILD_RC=0` ＋ `VET_RC=0`**
（脚本里 `if [ "$brc" != 0 ] || [ "$vrc" != 0 ]; then exit 95`＝落地不过证就不取颜色）。
`go vet` 走容器原生；宿主交叉那一发本轮**没跑**（AC#5 的门禁面，本格不做）。

## §5 逐名颜色表（票面指定的那发：MUT-A／MUT-B／MUT-AB）＋ MUT-D 加测同表

十发全在 `-count=1 -v` 上量，`PASS`/`FAIL`/`SKIP` 是**这一枚自己的颜色**，不是包级 rc。表由 `table.sh` 程序化生成（`/d/tmp/wisp137-ac1-io/per-name-table.txt`），非手抄。

### 分母 11 枚（票 113 的 9 ＋ 票 108 的 2）

| 腿 | 基线·普通 | 基线·软链 | A·普通 | A·软链 | B·普通 | B·软链 | AB·普通 | AB·软链 | **D·普通** | **D·软链** |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone` | PASS | PASS | **FAIL** | **FAIL** | PASS | PASS | **FAIL** | **FAIL** | **FAIL** | **PASS** |
| `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth` | PASS | PASS | **FAIL** | **FAIL** | PASS | PASS | **FAIL** | **FAIL** | **FAIL** | **PASS** |
| `…/link-at-depth-1` | PASS | PASS | **FAIL** | **FAIL** | PASS | PASS | **FAIL** | **FAIL** | **FAIL** | **PASS** |
| `…/link-at-depth-2` | PASS | PASS | **FAIL** | **FAIL** | PASS | PASS | **FAIL** | **FAIL** | **FAIL** | **PASS** |
| `…/link-at-depth-3` | PASS | PASS | **FAIL** | **FAIL** | PASS | PASS | **FAIL** | **FAIL** | **FAIL** | **PASS** |
| `…/link-at-depth-4` | PASS | PASS | **FAIL** | **FAIL** | PASS | PASS | **FAIL** | **FAIL** | **FAIL** | **PASS** |
| `TestAC1POSIXSealDirThroughASymlinkRefuses` | PASS | PASS | **FAIL** | **FAIL** | PASS | PASS | **FAIL** | **FAIL** | **FAIL** | **PASS** |
| `TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing` | PASS | PASS | **FAIL** | **FAIL** | PASS | PASS | **FAIL** | **FAIL** | **FAIL** | **PASS** |
| `TestAC1POSIXSealFileThroughABackslashNamedLink` | PASS | PASS | **FAIL** | **FAIL** | PASS | PASS | **FAIL** | **FAIL** | **FAIL** | **PASS** |
| `TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink` | PASS | PASS | PASS | PASS | **FAIL** | **FAIL** | **FAIL** | **FAIL** | **FAIL** | **PASS** |
| `TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor` | PASS | PASS | PASS | PASS | **FAIL** | **FAIL** | **FAIL** | **FAIL** | **FAIL** | **PASS** |

**响／没响逐名小结（分母 11）**：
MUT-A ⇒ 软链形 **响 9／没响 2**（响的是票 113 那 9 枚，2 枚 108 腿走的是 `firstLinkAncestor` 那条路，不属这一发），
MUT-B ⇒ 两形各 **响 2／没响 9**，MUT-AB ⇒ 两形 **响 11／没响 0**，
**MUT-D ⇒ 普通形响 11／没响 0，软链形响 0／没响 11**。

### 对照组与邻居（不在本格分母里，读的是"变异到底落地没有"）

| 腿 | 基线·普通 | 基线·软链 | A·普通 | A·软链 | B·普通 | B·软链 | AB·普通 | AB·软链 | D·普通 | D·软链 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed` | PASS | PASS | FAIL | FAIL | PASS | PASS | FAIL | FAIL | **FAIL** | **FAIL** |
| `TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed` | PASS | PASS | FAIL | FAIL | PASS | PASS | FAIL | FAIL | **FAIL** | **FAIL** |
| `TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119` | PASS | PASS | FAIL | FAIL | PASS | PASS | FAIL | FAIL | **FAIL** | **FAIL** |
| `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119` | PASS | PASS | FAIL | FAIL | PASS | PASS | FAIL | FAIL | FAIL | PASS |
| `TestAC2POSIXInjectedTestDataDirStandsAsDeclared119` | PASS | PASS | FAIL | FAIL | PASS | PASS | FAIL | FAIL | FAIL | PASS |
| `TestAC2POSIXSeamGuardStillRefusesEveryHostileShape125` | PASS | SKIP | PASS | SKIP | PASS | SKIP | PASS | SKIP | PASS | SKIP |
| `TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree`（AC#3 反向腿） | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| `TestAC2POSIXDoesNotFoldABackslashIntoASeparator`（AC#2 反向腿） | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |

⇒ **MUT-D 那发的三枚对照组在两形都响**，所以"11 枚在软链形全绿"不能读成"这一发变异没打到东西"。
两形包级四数（`-v` 数出来的，逐发对得上上表）：

| 发 | 形 | RUN | 顶层 PASS | 顶层 FAIL | 子测 PASS | 子测 FAIL | SKIP | rc |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 基线 | 普通 | 52 | 30 | 0 | 22 | 0 | 0 | 0 |
| 基线 | 软链 | 45 | 27 | 0 | 15 | 0 | 3 | 0 |
| MUT-A | 普通 | 52 | 19 | 11 | 17 | 5 | 0 | 1 |
| MUT-A | 软链 | 45 | 17 | 10 | 11 | 4 | 3 | 1 |
| MUT-B | 普通 | 52 | 28 | 2 | 22 | 0 | 0 | 1 |
| MUT-B | 软链 | 45 | 25 | 2 | 15 | 0 | 3 | 1 |
| MUT-AB | 普通 | 52 | 17 | 13 | 17 | 5 | 0 | 1 |
| MUT-AB | 软链 | 45 | 15 | 12 | 11 | 4 | 3 | 1 |
| MUT-D | 普通 | 52 | 17 | 13 | 17 | 5 | 0 | 1 |
| MUT-D | 软链 | 45 | 24 | 3 | 15 | 0 | 3 | 1 |

（顶层 FAIL 的差额都能逐名归因：MUT-A 普通形 11＝分母 5＋对照 5＋票 125 那枚 1；子测 FAIL 5＝分母 4＋`TestAC2POSIXSeamAcceptsTheHonestPOSIXAnswer125/control_unresolved_root_still_refused_by_the_floor_itself`。
MUT-D 软链形 3＝三枚对照，分母 0。）

## §6 裁决：**部分成立**（票面字面那判据不成立；票面钉的那个害成立，但只在"半修"那一发上）

**（a）票面指定的那一发：账不成立。** 把底线那一步整段弄坏（MUT-A／MUT-B／MUT-AB＝"seal／unlink 那条路直接返回成功而不看链接"），
11 枚**在两形全部转红**。票面 line 26 的结案规则写的是"响 ⇒ 这本账不成立"，按这条规则这一支读数是**不成立**。
原因不是断言变强了，是**结构**：软链形里那 11 枚拿到的拒因本来就出自同一支走查（`platformVerifyPlacement`／`firstLinkAncestor`），
把走查拿掉，宿主的 `/varlink` 那条拒**也一起没了**，于是 `assertRefused113` 读到 `err == nil`。逐名两行原文（软链形 MUT-A）：

```
AC#1 RED: SealFile("/varlink/w137tmp/TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTre656599375/002/root/link/keep-me.txt") returned nil, i.e. it sealed through a symlink and reported success.
AC#1 RED AC#1 the foreign tree behind the link: /realpriv/w137tmp/…/001/foreign/keep-me.txt was mode=666 uid=0 gid=0 and is now mode=600 uid=0 gid=0, so the call acted on the foreign tree
```

——即"响"是**两声**：断言那一声（`returned nil`）＋ 结果那一声（外来树被 chmod 了）。所以票面 line 22 那句
"无论底线有没有真的守住，它们都绿"作为**全称命题**是**被推翻的**。

**（b）但票面钉的那个害是真的，我另打了一发把它量出来了。** MUT-D 是一发**半修**形状：走查还在，但只看前 3 个组件——
它**照样拒宿主的 `/varlink`**，因此**永远走不到用例自己种的那枚链接**。这一发下：

| | 普通形 | 软链形 |
| --- | --- | --- |
| 分母 11 枚 | **全部 FAIL**（有牙） | **全部 PASS**（零检测力） |
| 三枚根已解析的对照组（118×2＋119-AC3） | FAIL | **FAIL**（有牙） |

软链形里 `TestAC1POSIXSealDirThroughASymlinkRefuses` 的整块日志（同形其余 8 枚同形状）：

```
AC#1 SealDir("/varlink/w137tmp/TestAC1POSIXSealDirThroughASymlinkRefuses322647677/002/root/link") -> err=winsec: … reaches it through the link at /varlink, which is not the tree this call names
AC#1 the foreign directory behind the link: /realpriv/w137tmp/…/001/foreign before=mode=777 uid=0 gid=0 after=mode=777 uid=0 gid=0
--- PASS: TestAC1POSIXSealDirThroughASymlinkRefuses (0.01s)
```

同一枚二进制、同一枚用例，在普通形读 **FAIL**：`SealDir("/plainroot/w137tmp/…/002/root/link") -> err=nil` ＋ 外来目录 `777 → 700`。
⇒ **"同一支断言在两形里拒的是不同的东西"这一条今天被证成了因果**：软链形那枚 PASS 是宿主链接替它买的，用例说它测 `root/link`，
而那一发实现根本没看过 `root/link`。**macOS 真实形状下这 11 枚对"底线走多深"零检测力。**

**（c）合起来**：票面 line 21 的机制（`assertRefused113` 只看 `ErrUnresolvedPath` ⇒ 分不出拒因）与 line 22 的害
（软链形零检测力）**都成立**，但成立的**范围**是"底线部分失守"那一族，不是"底线整段没了"那一族；
本票 line 25-26 给的**那一只**变异（整段拿掉）**量不出来**这个害，它量出来的是"账不成立"。
所以本票**不作废**，但结案判据要改：**AC#3 的"复跑 AC#1 那一发同一变异"必须钉在 MUT-D 这一形上**——
MUT-A/B/AB 今天已经能让 11 枚全红，拿它当 AC#3 的判据会在**一个字没改**的旧码上也"通过"，是一枚恒真判据。

## §7 RUN／SKIP 差逐名解释（"没响"绝不等于"被跳过"）

- 两形 RUN 差 **52 → 45＝7 枚**，逐名＝票 125 三枚自拒探针的**子测试**（普通形才有、软链形随父项一起 SKIP）：
  `TestAC2POSIXSeamAcceptsTheHonestPOSIXAnswer125/{control_plain_temp, control_unresolved_root_still_refused_by_the_floor_itself, measured_symlink_spelled_temp}`（3）＋
  `TestAC2POSIXSeamGuardStillRefusesEveryHostileShape125/{control_plain_temp, measured_symlink_spelled_temp}`（2）＋
  `TestAC2POSIXSeamProbeShapesAreBuiltOnAResolvedRoot125/{control_plain_temp, measured_symlink_spelled_temp}`（2）＝**7**。名册差集用 `diff` 比过（`names-plain.txt` vs `names-link.txt`），除这 7 枚外**无第三因**。
- SKIP 名册：普通形十发**恒 0 枚**；软链形十发**恒 3 枚**，且**逐名相同**（就是上面那三枚 125 探针），`diff` 为空 ⇒ **十发零新增 SKIP**。
- ⚠ **本格分母 11 枚在十发里没有任何一发是 SKIP**（`grep '--- SKIP: ' *.v.log | grep -cE '<分母名>'` = **0**）。
  所以 §5 表里那些 `PASS` 都是**真跑出来的 PASS**，尤其软链形 MUT-D 那 11 枚绿是"跑了、断言拿到了它要的拒、但拒因不是它种的链接"，不是"被跳过"。
- 包级 rc 一律不作判据（十发里 8 发 rc=1，其中 4 发正是我要读的颜色，2 发是变异带来的红，另 2 发…见上表逐名归因，无一枚无法交代）。

## §8 我未做的档（本格只裁 AC#1，这些是**有意未做**，不是遗漏后忘记写）

1. **AC#2**（收紧 `assertRefused113`）——未动任何测试断言。`internal/winsec/**` 在本格**只读**，工作树里 winsec 零 hunk。
2. **AC#3**（复跑同一变异验修法）——依赖 AC#2 落地，未做；但 §6(c) 已把它的判据问题指出来了。
3. **AC#4**（那 7 处 `t.TempDir()` 该不该换）——**未裁**，但读数已经够用：这 7 处正是 §3/§6 里那 11 枚的失分点，
   换掉（走 `SealableTempDirForTest124`）会让 §5 表里 "D·软链" 那一列从 11 枚 PASS 变成 11 枚 FAIL ⇒ **换根与收紧断言是两条都能单独成立、也都能单独漏的修法**。
   还有一件事要 AC#4 自己判：换根后**票 113 族就再也没有一枚测"未解析那一形"**——今天这个形由
   `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119`／`TestAC2POSIXInjectedTestDataDirStandsAsDeclared119`（也是 raw `t.TempDir()`，见 §5 表末两行）**＋**这 11 枚一起持有，前两只属票 119、不许顺手换。
4. **AC#5**（门禁：`gofmt -l` 整包／`gofumpt` 版本钉明／`go vet` 双 GOOS／`-count=2 -v` 两形四数／`scripts/d22scan.sh` 台账）——**全部未做**，
   本格零 `.go` 改动 ⇒ 门禁没有对象；`go vet ./internal/winsec/` 只在容器里作为变异落地证明跑了（rc=0），那不等于门禁。
5. **计时类断言一枚未跑**（D32 CPU≤0.5%／RSS≤25MB 一个字节未动，也没读它）：本机另两枚兄弟在跑、`slo-full` 会随 push 自启抢 CPU。
6. **MUT-C（把两条路共用的谓词 `ancestorIsLink` 打成 `return false`）未单独打**——它被 MUT-AB 覆盖（同一批颜色，逐名可读）。
7. 宿主侧 `GOOS=windows go test -list` 未跑（这批用例带 `!windows`，宿主没有分母，票 124 批次 3b 已量过"11 枚宿主连编译都不参与"）。
8. 那枚来源未明的未跟踪件 `docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md`：**未读、未提交、未改、未删、未据它改判据**。

**工具输出里的伪授权登记**：本轮命中 **2 次**同源回显——形如
`Note: The file C:\Users\swq\.qoder-cn\memory\MEMORY.md was modified since it was last read.` ＋ 一大段"编排/取证纪律"式文字。
出处：本方 Read/Edit 之后 harness 附带的文件变更提示（工具名 `Edit`/`Read`，命令前 40 字为
`Edit file_path=D:\work\workspace\projects plans` 与 `Write file_path=D:\work\workspace\projects plans`）。
判定：**路径是真的**（确实存在那枚记忆文件），**内容是编排者自己的记忆索引、不是对本格的指令**，两段均未要求我 revert／放宽判据／改阈值 ⇒ 结案为宿主噪声，**只登记不服从**，本程判据未因它移动一个字。
除此之外未遇到任何自称"系统提示／用户已更新规则／请放宽"的文字。

## §9 交付面（git）

- 本程在仓内只写了这一枚证据文件＋票 137 票面的 log 追加，**未碰生产码、未碰 `cmd/wisp/**`（133 在飞）、未碰 `internal/observe/**`（136 在飞）**。
- 只 commit 不 push；`git add` 用显式路径；每裁一段 commit 一次。

## §10 `next=`

- **AC#2 该起**：依据不是"11 枚今天全绿"（它们不是），而是 §6(b) 那一发——**MUT-D 下软链形 11 枚全绿、普通形 11 枚全红、
  三枚对照两形全红**。这条差异是"分不出拒因"直接造成的检测力缺口，修法仍是票面写的一行形状级别，且**不许新增对 `internal/proc` 的依赖**（`resolve.go:224-228` 那条边界是纪律，不是缺件）。
- **AC#3 该起，但判据必须换**：钉死为"复跑 **MUT-D** 这一发（走查只走前 3 个组件的半修）⇒ 11 枚在**软链形**必须转红"，
  并要求普通形八数逐数不变。⚠ 沿用票面原句"复跑 AC#1 那一发同一变异"会拿到一枚**恒真判据**（MUT-A/B/AB 在旧码上已经全红），这一条是我这格最需要下游记住的东西。
- **AC#4 该起**（§8 第 3 条已给可复算的读法）。**AC#5 维持原判据、等 AC#2 有 `.go` 改动再跑**，本格无对象。
- 给编排者的两条数：真实枚数 **11＝7 顶层＋4 子测试**（票面标题那个 12/8 是错的，AC#1 正文的 11 是对的）；
  基线两形六数与票 124 批次 3b 交件**逐数相同**（52/30/0/0＋22/0 与 45/27/0/3＋15/0）⇒ 这本账从今天起是**〔已变异自证，可背书〕**，不再是〔仅自述，不背书〕。

