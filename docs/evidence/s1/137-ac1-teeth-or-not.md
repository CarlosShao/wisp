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
| **MUT-D**（**加测**，半修那一发，不是票面指定的） | 两条走查都改成「只inspect前 3 个前缀」`if len(prefixes) > 3 { prefixes = prefixes[:3] }`：**仍拒宿主的 `/varlink`，但走不到用例自己种的链接** | `winsec_other.go:156-162`（`MUTATION-137-AC1-D` 注释＋截断），`winsec.go:284-288`（unlink 路同一刀） | `winsec_other.go` `2002cd7b128ec59d3c26486b77b10d03`、`winsec.go` `9134d7e9928955832b2bb2813af6614e` |

每发读数前都在**同一枚容器里**跑 `go build ./...` 与 `go vet ./internal/winsec/`，十发全部 **`BUILD_RC=0` ＋ `VET_RC=0`**
（脚本里 `if [ "$brc" != 0 ] || [ "$vrc" != 0 ]; then exit 95`＝落地不过证就不取颜色）。
`go vet` 走容器原生；宿主交叉那一发本轮**没跑**（AC#5 的门禁面，本格不做）。
