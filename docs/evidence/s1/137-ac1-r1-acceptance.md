# 137 AC#1 —— 验收方第二程（r1）：把测量方那句"部分成立"独立重走一遍

裁决方：`acceptor-ticket137-ac1-r1`（**非实现者、非测量方**）。本格只裁 **票 137 AC#1 这一格**，
不修任何东西、不翻任何勾。测量方那程是 `worker-ticket137-ac1`，它的证据是
`docs/evidence/s1/137-ac1-teeth-or-not.md`；本件是**独立复现**，不是它的附录。

写作与取证交错进行：先测、再写、最后提交，每裁一节 commit 一次。

## §0 锚点与树（含 winsec 侧差集）

- 开工首读 `git rev-parse --short HEAD` = **`4a0d7a4`**，全 sha `4a0d7a48d03dc6aee1a3e18ceee1d2a99269b9e6`，
  `git cat-file -t` = `commit`。**本程全部读数都在这枚锚上量**。开工时刻 `date -u` = 2026-09-23 15:21:16 UTC ⇒ 本机 23:21 +8。
- 锚点来源纪律：**没有任何一枚 sha 是从别人的报告或工具输出里直接取来当锚用的**。
  测量方自报的 `1d38206` 我只在"核它存不存在、是不是我这一枚的祖先"这两条只读命令里用到：
  `git cat-file -t 1d38206` = `commit`（真实存在），`git merge-base --is-ancestor 1d38206 4a0d7a4` = 真（它的锚在我这枚的祖先线上）。
- **被测面差集（它核过"winsec 侧零 hunk"，这条我重走了、没沿用）**：
  `git log --oneline 1d38206..4a0d7a4 -- internal/winsec/` **输出为空** ⇒ 从它的锚到我的锚，
  `internal/winsec/` 一字未动，两程读的是同一版被测码。
  （同一个区间里飘进来的提交全在 `docs/**`、`.scratch/**`、`internal/observe/**`、`cmd/wisp/**`，与本格无关。）
- 树：**全部**由 `git archive 4a0d7a48d03dc6aee1a3e18ceee1d2a99269b9e6 | tar -x -C /d/tmp/wisp137r2-<名>` 落成，
  仓内**未建 worktree、未 checkout、未 reset**。目录一律只建不删：
  `wisp137r2-tree0`（纯净）／`-tree-muta`／`-tree-mutb`／`-tree-mutab`／`-tree-mutd`，台件与十＋两发日志在 `wisp137r2-io/`。
  ⚠ 我没有把脏工作树当被验版本：工作树里此刻有别人的未提交件（`docs/reports/**` 两枚已改、一枚未跟踪），
  `git status --porcelain internal/winsec/` 输出为空 ⇒ winsec 侧工作树与锚上同版，但我仍然只用归档树取数。
- 生产码零改动：本程对 `internal/winsec/**` **只读**（AC#1 不修东西，收紧断言是 AC#2 的地界）。

## §1 容器与挂载证明

- `docker info` ⇒ ServerVersion **29.6.2**、`linux/x86_64`。宿主是 Windows，而这一族用例带 `//go:build !windows` ⇒
  **宿主没有分母**，本程一片读数都不在宿主取（宿主只跑 `git`／`grep` 这类取证命令）。
- 镜像 `golang:1.27`，容器内 `go version` = **`go1.27.1 linux/amd64`**。
- 挂载写成 `/d/...` 且带 `MSYS_NO_PATHCONV=1`；**每一发**开跑前先 `ls -l /src/go.mod`：

  ```
  -rwxrwxrwx 1 root root 883 Sep 23 15:19 /src/go.mod
  go.mod sha256: d13ba3de2d319f298ed1e598f1702015ca50e646ae293daf183545ce794f40fc
  go version go1.27.1 linux/amd64
  ```

  十二发（十发正式＋两发复跑）里这枚 sha256 与字节数**逐发相同**；脚本里另有硬闸
  `[ "$(wc -c < /src/go.mod)" = "883" ] || exit 97` ⇒ 静默挂空取不到读数。
  ⚠ 顺带把一枚容易读成"分歧"的数说清：仓内工作树与 `git show <锚>:go.mod` 都是 **855** 字节，
  而 `git archive | tar -x` 出来的快照里是 **883** 字节（差 28＝`.gitattributes` 第 1 行 `* text=auto` 配 `core.autocrlf=true` 带来的行尾差，
  28 行 go.mod 每行多一枚 CR）。
  我引用挂载证明时用的是**容器里那个文件的实际字节数**（883），与测量方报的 883 同形，不是同一个量被抄了两遍。
- 被测文件指纹（tree0，容器内 `md5sum` 与宿主侧 `md5sum` 同值）：
  `internal/winsec/winsec_other.go` = `b5056918be4ed13817d236fbcae0f477`、
  `internal/winsec/winsec.go` = `a6144c880de80e43bb1393f3624e7221`
  ⇒ **与测量方 §1 报的两枚 md5 逐字相同**，这是"两程读的是同一版码"的第二条独立凭据（第一条是 §0 的差集为空）。
- 两形硬断言（我的台件命名与测量方**不同**，免得把它的形状当成我的）：
  - 软链形：`mkdir -p /r2priv/w137r2tmp` ＋ `ln -s /r2priv /r2link` ＋ `[ -L /r2link ]` ＋ `readlink` 逐字核为 `/r2priv`
    ＋ `ls -ld` 打印 `lrwxrwxrwx 1 root root 7 ... /r2link -> /r2priv`，`TMPDIR=/r2link/w137r2tmp`；
  - 普通形：断言 `/r2link` **根本不许存在** ＋ `[ ! -L /r2plain ]`，`TMPDIR=/r2plain/w137r2tmp`。
  - 十二发**无一命中 97／98／99**。新鲜容器一枚一发（`--rm`，每发从零起），形状目录不许预存在。
- 模块缓存：`GOPROXY=off` ＋ 我自己的两枚 named volume `wisp137r2-gomod`／`wisp137r2-gobuild`
  （离线复用既有缓存只是省时间，**不是复用别人的读数**；缺件会直接 fail，不存在"网络慢装出来的绿"）。
- 跑法：`go test -count=1 -v ./internal/winsec/`；读数前先 `go build ./...` 与 `go vet ./internal/winsec/`，
  两者任一非 0 ⇒ `exit 95`＝**落地不过证就不取颜色**。十二发全部 `BUILD_RC=0` ＋ `VET_RC=0`。
- 台件（全部在仓外 `D:\tmp\wisp137r2-io\`）：`mutate.py`（变异，逐处锚串唯一性断言）、`make-trees.sh`（建五棵快照）、
  `run.sh`（容器内跑发＋硬闸）、`parse.py`（把 `-v` 日志程序化拆成逐名颜色与名册，非手抄）、`matrix.py`（十发矩阵）。

## §2 四发逐发（先证落地，再读颜色）

四发全部打在**我自己的快照树**上，每发一棵新树、单发单因子（MUT-AB 是"两枚文件各自一刀"的发，不是叠在别的发上）。
我的补丁文本与测量方**不同**（注释串是 `MUTATION-137-R2-*`，它们的是 `MUTATION-137-AC1-*`），所以指纹也不同：
`MUT-A` 后 `winsec_other.go` = `22784214a31f8b5218e7ac6510260f57`（它们报 `0e2a136c…`）、`MUT-B` 后 `winsec.go` = `03dfee60f6eae07319a5cef01ebc1fe6`（它们报 `ae7c51bf…`）、
`MUT-D` 两枚 = `2ffb39b279b2f4a585d3d007703e37a3` ＋ `2aaa102f87262f91b87223f593c92728`（它们报 `2002cd7b…`／`9134d7e9…`）。
⇒ 这是**同一语义的第二份实现**，不是抄它的字节。每发都各证一次"没改的那一枚文件指纹仍是基线值"。

落地证明（`grep -n` 出被改后那一行原文，取自归档树；同一发在容器内另跑 `go build ./...` 与 `go vet ./internal/winsec/`）：

| 代号 | 形状 | 被改后那一行（`grep -n` 原文） | 另一枚文件 | BUILD_RC／VET_RC |
| --- | --- | --- | --- | --- |
| **MUT-A** | `platformVerifyPlacement` 退回旧 stub | `winsec_other.go:156: // MUTATION-137-R2-A: the link walk is gone; the floor answers "yes, unchanged" for every spelling.` ＋ `:157: return path, nil` | `winsec.go` 仍 `a6144c88…`＝未改 | 0／0（两形） |
| **MUT-B** | `firstLinkAncestor` 直接 `return ""` | `winsec.go:280: // MUTATION-137-R2-B: the unlink route answers "no ancestor is a link" for every spelling.` ＋ `:281: return ""` | `winsec_other.go` 仍 `b5056918…`＝未改 | 0／0（两形） |
| **MUT-AB** | 两者同打 | 上两行同时出现 | 两枚各自 `22784214…`／`03dfee60…` | 0／0（两形） |
| **MUT-D**（**我自己落的**，票面原本没有这一发） | 两条走查都只走前 3 个组件 | `winsec_other.go:156-163`（注释＋`sealPrefixes := pathPieces(path)`／`if len(sealPrefixes) > 3 { sealPrefixes = sealPrefixes[:3] }`）、`winsec.go:281-283`（`// MUTATION-137-R2-D: same half-fix on the unlink route: three components, no further.`＋`prefixes = prefixes[:3]`） | 两枚都改（`2ffb39b2…`／`2aaa102f…`） | 0／0（两形） |

⚠ MUT-D 的"前 3 个组件"我不是照它的描述猜的：`pathPieces`（`winsec.go:318-342`）在 POSIX 上把
`/r2link/w137r2tmp/<用例>/002/root/link/keep-me.txt` 拆成 7 枚前缀（**前导分隔符不单独成一枚**），
用例自种的链接落在**第 6 枚**；截前 3 枚后仍含 `/r2link` ⇒ "照旧拒宿主链接、永远看不到自己那枚"这一形**由源码语义保证**，
再用读数验一遍（见下）。我的 `mutate.py` 对每处锚串断言"全文只命中 1 次"，命中数不是 1 就直接退出、不出数。

### §2.1 包级四数（我的十二发，逐发自己数出来的）

| 发 | 形 | RUN | 顶 PASS | 顶 FAIL | 顶 SKIP | 子 PASS | 子 FAIL | rc |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 基线 | 普通 | 52 | 30 | 0 | 0 | 22 | 0 | 0 |
| 基线 | 软链 | 45 | 27 | 0 | 3 | 15 | 0 | 0 |
| MUT-A | 普通 | 52 | 19 | 11 | 0 | 17 | 5 | 1 |
| MUT-A | 软链 | 45 | 17 | 10 | 3 | 11 | 4 | 1 |
| MUT-B | 普通 | 52 | 28 | 2 | 0 | 22 | 0 | 1 |
| MUT-B | 软链 | 45 | 25 | 2 | 3 | 15 | 0 | 1 |
| MUT-AB | 普通 | 52 | 17 | 13 | 0 | 17 | 5 | 1 |
| MUT-AB | 软链 | 45 | 15 | 12 | 3 | 11 | 4 | 1 |
| MUT-D | 普通 | 52 | 17 | 13 | 0 | 17 | 5 | 1 |
| MUT-D | 软链 | 45 | 24 | 3 | 3 | 15 | 0 | 1 |

另加两发**同树新容器复跑**（`MUT-A`、`MUT-B` 各两形）：四数与首发逐数相同，且**逐名颜色与首发完全相同**（`diff` 名册为空）⇒ 读数可复现。
这十行的每一格与测量方 §5 那张包级表**逐数相同**（它表里 `基线·普通 52/30/0/22/0/SKIP0`…直到 `MUT-D·软链 45/24/3/15/0/SKIP3`）。

### §2.2 分母 11 枚的逐名颜色（我这一程的十发，程序化生成、非手抄）

| 腿 | 基普 | 基链 | A普 | A链 | B普 | B链 | AB普 | AB链 | **D普** | **D链** |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone` | PASS | PASS | FAIL | FAIL | PASS | PASS | FAIL | FAIL | **FAIL** | **PASS** |
| `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth` | PASS | PASS | FAIL | FAIL | PASS | PASS | FAIL | FAIL | **FAIL** | **PASS** |
| `…/link-at-depth-1` | PASS | PASS | FAIL | FAIL | PASS | PASS | FAIL | FAIL | **FAIL** | **PASS** |
| `…/link-at-depth-2` | PASS | PASS | FAIL | FAIL | PASS | PASS | FAIL | FAIL | **FAIL** | **PASS** |
| `…/link-at-depth-3` | PASS | PASS | FAIL | FAIL | PASS | PASS | FAIL | FAIL | **FAIL** | **PASS** |
| `…/link-at-depth-4` | PASS | PASS | FAIL | FAIL | PASS | PASS | FAIL | FAIL | **FAIL** | **PASS** |
| `TestAC1POSIXSealDirThroughASymlinkRefuses` | PASS | PASS | FAIL | FAIL | PASS | PASS | FAIL | FAIL | **FAIL** | **PASS** |
| `TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing` | PASS | PASS | FAIL | FAIL | PASS | PASS | FAIL | FAIL | **FAIL** | **PASS** |
| `TestAC1POSIXSealFileThroughABackslashNamedLink` | PASS | PASS | FAIL | FAIL | PASS | PASS | FAIL | FAIL | **FAIL** | **PASS** |
| `TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink` | PASS | PASS | PASS | PASS | FAIL | FAIL | FAIL | FAIL | **FAIL** | **PASS** |
| `TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor` | PASS | PASS | PASS | PASS | FAIL | FAIL | FAIL | FAIL | **FAIL** | **PASS** |

**逐发"响几枚／没响几枚"（我的量，分母 11）**：

| 发 | 普通形 | 软链形 | 与测量方 |
| --- | --- | --- | --- |
| MUT-A | 响 9／没响 2（没响的是 108 那 2 枚，它们走 unlink 路，不属这一发） | 响 9／没响 2 | **一致** |
| MUT-B | 响 2／没响 9 | 响 2／没响 9 | **一致** |
| MUT-AB | 响 11／没响 0 | 响 11／没响 0 | **一致** |
| **MUT-D** | **响 11／没响 0** | **响 0／没响 11（全绿）** | **一致（这一发是本格的承重墙，我独立落、独立取名册，没照抄它的红名）** |

### §2.3 三发整段弄坏：票面那句"绿得没有理由"在这一族上**不成立**（复核通过）

MUT-A／MUT-B／MUT-AB 在**软链形**里同样把这 11 枚（或其中 9／2 枚）打红。红是**两声**，我的日志里读得到（`MUT-A·软链`）：

```
placement_symlink_113_other_test.go:153: AC#1 RED: SealFile("/r2link/w137r2tmp/TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTre2309564223/002/root/link/keep-me.txt") returned nil, i.e. it sealed through a symlink and reported success. AC#1 requires either a refusal or an action confined to the link itself.
placement_symlink_113_other_test.go:155: AC#1 RED AC#1 the foreign tree behind the link: /r2priv/w137r2tmp/TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTre2309564223/001/foreign/keep-me.txt was mode=666 uid=0 gid=0 and is now mode=600 uid=0 gid=0, so the call acted on the foreign tree
```

机制与它的说法一致：走查整体拿掉后，宿主 `/r2link` 那条拒**也一起没了** ⇒ 断言读到 `err=nil` ＋ 外来树被 chmod。
⇒ 票面 line 22 那句"无论底线有没有真的守住，它们都绿"作为**全称命题**被推翻（我这边同样读不到反例）。

### §2.4 MUT-D：票面钉的那个害**只有这一发能造出来**（我复现成功）

同一枚用例、同一枚二进制族，两形读数是**反的**：

软链形（**PASS**，拒因是宿主的链接）——`TestAC1POSIXSealDirThroughASymlinkRefuses`：

```
placement_symlink_113_other_test.go:203: AC#1 SealDir("/r2link/w137r2tmp/TestAC1POSIXSealDirThroughASymlinkRefuses202643511/002/root/link") -> err=winsec: refusing to seal /r2link/w137r2tmp/TestAC1POSIXSealDirThroughASymlinkRefuses202643511/002/root/link: winsec: path is not provably resolved, refusing to seal: /r2link/w137r2tmp/TestAC1POSIXSealDirThroughASymlinkRefuses202643511/002/root/link reaches it through the link at /r2link, which is not the tree this call names
```

108 那族同样只看到宿主的链接：

```
ancestor_separator_108_other_test.go:79: RemoveUnlinked("/r2link/w137r2tmp/TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink3890533288/002/root/link/keep-me.txt") -> err=winsec: entry is a link to something else, not private data /r2link/w137r2tmp/TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink3890533288/002/root/link/keep-me.txt: the spelling reaches it through the link at /r2link, and whatever lives behind that link is not this tree's data to delete
```

普通形（**FAIL**，`err=<nil>`）——同一枚用例：

```
placement_symlink_113_other_test.go:203: AC#1 SealDir("/r2plain/w137r2tmp/TestAC1POSIXSealDirThroughASymlinkRefuses1856693096/002/root/link") -> err=<nil>
ancestor_separator_108_other_test.go:79: RemoveUnlinked("/r2plain/w137r2tmp/TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink4154713178/002/root/link/keep-me.txt") -> err=<nil>
```

⇒ "同一支断言在两形里拒的是**不同的东西**"这一条，我这边也从读数里取到了因果：软链形那 11 枚绿是宿主链接替它买的，
那一发实现**根本没看过**用例自己种的 `root/link`。**macOS 真实形状下这 11 枚对"底线走多深"零检测力。**

**对照组（证 MUT-D 真落地，不是"没打到"）**：三枚根已解析的用例在 **D 形两形都红**，红在断言本身：

```
placement_leaf_118_other_test.go:81: AC#1 RED: SealFile("/r2priv/w137r2tmp/TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed1221542493/002/root/artifact.txt") returned nil, i.e. it sealed through a symlink and reported success. AC#1 requires either a refusal or an action confined to the link itself.
placement_leaf_118_other_test.go:107: AC#1 RED: PrivateFile("/r2priv/w137r2tmp/TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasName2055800904/002/root/artifact.txt") returned nil, i.e. it sealed through a symlink and reported success. AC#1 requires either a refusal or an action confined to the link itself.
dataroot_symlink_119_other_test.go:250: AC#3 RED: sealing through a link inside a resolved data root returned nil
```

D·软链形包级 `顶 FAIL 3` 的名册就是这三枚（`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`、
`TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed`、`TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed`），
**没有第四因** ⇒ 那 11 枚的绿不能读成"变异没打到东西"。
⚠ 但**点名这件事上我的派单与票面更正块写错了**：真正的"两形都红"三枚是 **118 两枚 ＋ `TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`**，
而更正块里被算进对照组的那枚 `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119` 在我的读数里是 **D·普通 FAIL／D·软链 PASS**（它本来就用 raw `t.TempDir()`、那一形里它*就该*被宿主的链接拒掉），
详见 §5 与 R-137-1。

## §3 枚数裁决：**11 枚＝7 顶层＋4 子测试**（票面标题那个"12＝8 顶层＋4 子测试"我复现得出它的来源、但不采信）

我逐枚数的，判据与测量方同一句但各量各的：**票 113／票 108 的 POSIX 拒绝腿里，`root` 仍由未解析的 `t.TempDir()` 供给的那些枚**。
先把两枚文件里所有顶层 `func Test` 摊开逐名判"算不算拒绝腿"（`grep -n '^func Test'` 自己取的行号）：

| 顶层用例（文件:行） | 根从哪来 | 走哪条路 | 算不算本格拒绝腿 | 为什么 |
| --- | --- | --- | --- | --- |
| `placement_symlink_113_other_test.go:144` `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone` | `:146` raw `t.TempDir()` | seal（`platformVerifyPlacement`） | **算** | 要的就是"拒"，且根未解析 |
| `:165` `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth` | `:169` raw | seal | **算**（顶层）＋**它的 4 枚子测试各算一枚** | 同上；4 枚 depth 子测是 4 个独立断言点 |
| `:195` `TestAC1POSIXSealDirThroughASymlinkRefuses` | `:197` raw | seal | **算** | 同上 |
| `:211` `TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing` | `:213` raw | seal | **算** | 同上 |
| `:236` `TestAC1POSIXSealFileThroughABackslashNamedLink` | `:238` raw | seal | **算** | 同上 |
| `:253` `TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree` | `:254` `SealableTempDirForTest124` | seal | **不算** | 它要的是"**成功**"（反向腿），把它算进"拒绝腿"就是把尺子反着拿 |
| `:287` `TestAC3POSIXSealStillWorksNextToAndThroughRealDirectoriesAndLinks` | `:289` 同上 | seal | **不算** | 同上（还专门种一枚*不是*祖先的链接来防过度拒绝） |
| `:335` `TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator` | `:337` 同上 | seal | **不算** | 同上 |
| `ancestor_separator_108_other_test.go:67` `TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink` | `:69` raw | unlink（`firstLinkAncestor`） | **算** | 要"拒"，根未解析 |
| `:94` `TestAC2POSIXDoesNotFoldABackslashIntoASeparator` | `:96` `SealableTempDirForTest124` | unlink | **不算** | 要"这次必须成功删掉自己的文件" ⇒ 反向腿 |
| `:125` `TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor` | `:127` raw | unlink | **算** | 要"拒"，根未解析 |
| `:152` `TestAC4POSIXFloorAnswersInsideTheNamedTree` | `:156` 已解析 | seal＋resolve | **不算** | 它的"拒"只针对 `../` 这种父指针拼写、根又是已解析的 ⇒ 天生有牙，且不在那 7 处 raw `t.TempDir()` 名单里 |

⇒ **7 顶层（113 族 5 ＋ 108 族 2）＋ 4 子测试 ＝ 11 枚**。与"7 处有意未换的 `t.TempDir()`"逐处对得上
（`placement_symlink_113_other_test.go:146/169/197/213/238` ＋ `ancestor_separator_108_other_test.go:69/127`，我自己 `grep -n` 复核，第 2 处一处覆盖 1 顶层＋4 子测）。
`assertRefused113` 的真位置我再量一遍：**`placement_symlink_113_other_test.go:127`**（与 108 那枚 `:127` 的递根点同号不同物，别抄串）。

**那个"12"是从哪来的（我能复现来源）**：`grep -c '^func Test'` 在 `placement_symlink_113_other_test.go` 上**恰好得 8**（已核），
把 8 枚顶层全当拒绝腿再 ＋4 枚子测试＝12 —— 这正是票面标题与"为什么现在立案"表里那个数，也正好**漏掉** 108 真正的 2 枚顶层拒绝腿。
所以这不是"另一版分母"，是**同一枚文件里成功腿与拒绝腿混装**造成的口径错。
我也盘了另一枚凑法：两枚文件的顶层数 `8＋4＝12` 也等于 12，但那会是"12 枚全是顶层"，**与票面写的"8 顶层＋4 子测试"这个拆法不符**；
按票面自己的拆法，来源就是前者（8 枚顶层全部当拒绝腿 ＋ 那 4 枚子测试，同时漏掉 108 的 2 枚）。

**裁决**：测量方的"11＝7＋4"我**独立复现**（逐名、逐行、逐个根来源），票面标题的"12＝8＋4"不成立 ⇒ **§5 的枚数更正块成立**。

⚠ 我**没有**用 `grep -c '^func Test'` 当分母，也没有沿用它的名单：上表 12 行逐名给了"算／不算"的理由，
分母那 11 枚的名字在 §2.2 的逐名表里一枚一枚露过脸，名册差集见 §4。

## §4 名册与 SKIP 账（"没响"绝不等于"被跳过"，这条我逐名重走）

- **同形内 RUN 名册逐名相同**：普通形六发（基线＋四发变异＋`MUT-A/B` 的复跑）都是 **52 枚**、软链形六发都是 **45 枚**，
  每发的名册与同形基线 `diff` **为空** ⇒ 十发里没有任何一枚用例"消失"（没有 panic 吞读数、也没有静默少跑）。
  `grep -c '^panic|^fatal error'` 在**十四份日志**（含两发复跑的 build/vet 输出另算）里全部 **0**。
- **两形差 7 枚**逐名（`diff base-plain base-link` 的实录，除这 7 枚外无差集）：
  `TestAC2POSIXSeamAcceptsTheHonestPOSIXAnswer125/{control_plain_temp, control_unresolved_root_still_refused_by_the_floor_itself, measured_symlink_spelled_temp}`（3）＋
  `TestAC2POSIXSeamGuardStillRefusesEveryHostileShape125/{control_plain_temp, measured_symlink_spelled_temp}`（2）＋
  `TestAC2POSIXSeamProbeShapesAreBuiltOnAResolvedRoot125/{control_plain_temp, measured_symlink_spelled_temp}`（2）＝ **7**。
  ⇒ 测量方把差 7 枚归给**票 125 那三枚自拒探针的子测试**这条**归因成立**：父项在软链形自己 `t.Skipf` ⇒ 子测试从未被创建，不是"少跑"、也不是本票分母在动。
- **SKIP 账**：普通形**恒 0 枚**；软链形**恒 3 枚**且**逐名相同**
  （`TestAC2POSIXSeamProbeShapesAreBuiltOnAResolvedRoot125`、`TestAC2POSIXSeamGuardStillRefusesEveryHostileShape125`、`TestAC2POSIXSeamAcceptsTheHonestPOSIXAnswer125`）
  ⇒ 十发**零新增 SKIP**。
- ⚠ **本格分母 11 枚在十发里没有任何一发是 SKIP**：逐名 grep（`--- SKIP:` 行里点名 113/108 那几枚）逐发计数 **＝ 0**
  ⇒ §2.2 表里那些 `PASS` 全是跑出来的 PASS；尤其 **D·软链那 11 枚绿**是"跑了、断言也拿到了它要的拒、但拒因不是它自己种的链接"，不是"被跳过"。
- **每一发的红都能逐名归因，没有说不去的红**（我的名册实录）：

  | 发·形 | 顶 FAIL | 归因 | 子测 FAIL | 归因 |
  | --- | --- | --- | --- | --- |
  | A·普通 | 11 | 分母 5（113 族顶层）＋对照 5（119×3＋118×2）＋票 125 那枚 1 | 5 | 分母 4（depth-1..4）＋125 的 `control_unresolved_root_still_refused_by_the_floor_itself` |
  | A·软链 | 10 | 分母 5＋对照 5（125 那枚此形是 SKIP） | 4 | 分母 4 |
  | B·普通／B·软链 | 2／2 | 就是分母的 108 两枚 | 0／0 | — |
  | AB·普通 | 13 | 分母 7＋对照 5＋125 那枚 1 | 5 | 分母 4＋125 那枚子测 |
  | AB·软链 | 12 | 分母 7＋对照 5 | 4 | 分母 4 |
  | **D·普通** | **13** | 与 AB·普通**逐名相同** | **5** | 与 AB·普通逐名相同 |
  | **D·软链** | **3** | **只有三枚对照**（`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`＋`TestAC118…` 两枚），分母 **0** | **0** | — |

- 基线六数对齐**第三本账**：票 124 批次 3b 交件的 `POST-P 普通 52/30/0/0＋22/0`、`POST-L 软链 45/27/0/3＋15/0`
  （`docs/evidence/s1/124-ac2b-3b-conversion.md:220`、`:229`）与我这一程的基线两行**逐数相同**，测量方的中间那本也对得上
  ⇒ 三本账同一版被测码，锚点漂移不成立（§0 的差集为空是同一件事的第二种量法）。

## §5 两块「更正块」的复核（它们本身也是编排者的未验证断言）

票面文末那两块用 `>` 引起来的更正块是编排者按测量方报回写的。我按派单要求把它们当**断言**逐条重走，不背书。

**① 「AC#3 沿用原句会拿到一枚恒真判据」——结论成立，但它给的依据句要收窄。**

- 成立的部分（我独立量）：票面 AC#3 原句是"复跑 AC#1 那一发同一变异 ⇒ 这 12 枚必须转红"。
  AC#1 那一发的形状（票面 line 30-31：让 seal 与 unlink 两条路都不看链接）＝我的 **MUT-AB**，
  它在**一个字都没改的旧码**上就把 11 枚在**两形全部打红**（普通 11／软链 11，§2.2 逐名）。
  ⇒ 拿它当"修好了才算数"的尺子，**改与不改都满足** ⇒ **恒真判据**这一条我复现了。
- 要收窄的部分：更正块写的是"**MUT-A／MUT-B／MUT-AB 那三发**在未修的旧码上就已经全红（11 枚全响）"。
  实测**只有 AB 全响**：单发 **A 响 9**（113 族）／**B 响 2**（108 族），另两枚各自 2／9 不响。
  结论方向不变（AC#1 指定的那一发本就是两路同坏），但"三发都全响"这句**过度概括**，已记 R-137-2。
- ⚠ 顺带一条给下游的硬话（这条派单没让我判，但它是我这几发读数的直接推论）：**判据也不能钉在单发 A 或单发 B 上**——
  收紧断言之后复跑 MUT-A，108 那两枚走的是 unlink 路、链路完好、拒的就是自己种的 `root/link` ⇒ 它们**修完也不会响**；
  于是"11 枚必须转红"在单发 A 上会变成一枚**恒不满足**的判据（另一类坏尺）。
  ⇒ 能同时满足"不恒真"与"修完真能达成"的，目前只有 **MUT-D** 这一形（我的读数：软链 0 响／普通 11 响）。

**② 「真实枚数＝11＝7 顶层＋4 子测试」——完全复现。** 见 §3：我逐名数、逐名给算／不算的理由，
并复现出"12"的可复现来源（`grep -c '^func Test'` 在 113 文件恰得 8）。`assertRefused113` 真位置 `placement_symlink_113_other_test.go:127` 我也自己 `grep -n` 过。

**③ 更正块里那句「对照组那三枚根已解析的用例（119 族两枚 ＋ `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119`）在两形都必须仍红」——这条按字面执行会出事。**

- 我的 MUT-D 读数里**两形都红**的三枚是：`TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed`、
  `TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed`、`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`
  （D·软链形包级顶 FAIL 恰好就是这 3 枚、子测 0，名册见 §4）。
- 而被点进对照组的那枚 `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119` 实测是 **D·普通 FAIL／D·软链 PASS**。
  原因不是 bug：它**本来就用 raw `t.TempDir()`（`dataroot_symlink_119_other_test.go:194`）**、断言的就是"未解析的根必须被拒"，
  软链形里宿主的链接给了它一个**合法的拒** ⇒ 它和分母那 11 枚是**同一种病**，不是"根已解析的对照"。
  同样形状的 `TestAC2POSIXInjectedTestDataDirStandsAsDeclared119`（`:286` raw）也是 D·普通 FAIL／D·软链 PASS。
- ⚠ 若"119 族两枚"指的是 winsec 那枚 `TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119` 加上
  `TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119`（源码注释 `winsec_other.go:139` 点过名）——后者**不在 winsec 包**，
  它在 `cmd/wisp/secret_dataroot_119b_test.go:201`，**本格十发读数里没有它**（且那枚包正有兄弟代理在飞，我没去跑）。
- ⇒ 这条按字面写会造出一枚**永远无法满足**的落地凭据，把未来 AC#3 的实现方推向"去修一枚没坏的用例"。已记 **R-137-1**。
  我在 §2.4 用的是**我自己读出来的那三枚**，不是更正块点名的那三枚。

**④ 更正块 ②「普通形八数逐数不变」——需要一个具体参照值，否则不可复算。** 我这程给出的就是这八个字：
MUT-D·普通形 `RUN=52 顶 PASS=17 顶 FAIL=13 顶 SKIP=0 子 PASS=17 子 FAIL=5 子 SKIP=0 rc=1`，
且**逐名 FAIL 名册与 MUT-AB·普通形完全相同**（§4 表）⇒ AC#3 落地时"普通形不许新增红／不许由绿转 SKIP"就按这两句核。

## §6 我没做的档（诚实列，不是"遗漏后忘记写"）

1. **AC#2／AC#3／AC#4／AC#5 一格未做**。`internal/winsec/**` 全程只读（`git status --porcelain internal/winsec/` 空），
   收紧 `assertRefused113` 是 AC#2 的地界，我一行没动。**AC 勾一枚未翻**（含 AC#1 本身——终裁归 owner 落笔）。
2. **MUT-C（把两条路共用的谓词 `ancestorIsLink` 打成 `return false`）我没造也没跑**。
   测量方说"它被 MUT-AB 覆盖、同一批颜色"——**这句我没独立重走**，按本仓规矩记成
   **〔验收方没造的支，不算已测〕**，不替它背书，也不据此减损 §2 那四发的结论。
3. **票面指定的"形如让 `RemoveUnlinked`／seal 那条路直接返回成功"我只按两枚函数级形状落了**
   （A＝seal 路、B＝unlink 路、AB＝两路），没去追"`RemoveUnlinked` 整枚导出函数直接 `return nil`"那种更外层写法；
   它与 B 的区别只在 `filepath.IsAbs` 那半条相对路径守卫还在不在，本格结论不依赖它。
4. **宿主侧（Windows）一枚读数都没取**：这 11 枚带 `//go:build !windows`，宿主没有分母；
   也**没跑** `GOOS=windows go test -list`（那是票 124 批次 3b 量过的事，我没重走）。
5. **门禁五读数全未做**（`gofmt -l` 整包／`gofumpt`（版本未钉明）／`go vet` 双 GOOS／`-count=2 -v` 两形四数／`sh scripts/d22scan.sh`）。
   容器内 `go vet ./internal/winsec/` 只作为**变异落地证明**跑过（十二发全 rc=0），那不等于门禁。
6. **计时类断言一枚未跑**（D32 的 CPU≤0.5%／RSS≤25MB 一个字节没动也没读）：本机两枚兄弟在飞、`slo-full` 会随 push 自启抢 CPU。
7. **前例只当参照、没重判**：票 124 批次 3b 那本"邻居账"我只对齐了**基线六数**（§4 末条），
   它那 17 枚红名的逐名裁决**我没重走一遍**（派单明令"不要重判它"）。
8. **两枚 119 族的 success 腿（`TestAC1POSIXSymlinkedTempDirRouteBecomesSealable119`／`…ConfigDir…`）四发全绿这件事我没去判为什么**——
   那是 AC#4 换根地界的事，本格只记下"它们不在变异面上"这个事实。
9. **`cmd/wisp` 那枚 `TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119` 没跑**（不在 winsec 分母，且 `cmd/wisp` 有兄弟在飞），
   我只用 `grep` 确认了它的所在（§5③）。
10. **未跟踪件 `docs/reports/2026-09-23-gap-analysis-vs-oss-harnesses.md` 未读、未提交、未改、未删、未据它改判据**；
    工作树里别人那两枚已改文件（`docs/reports/injection-timeline.md`、`pending-and-issues.md`）我没碰、没 stage。
11. **两枚在飞兄弟的半成品未读**（`acceptor-ticket133-ac2-r2` 的 `cmd/wisp`、`acceptor-ticket136-ac8-ac9-r1` 的 `internal/observe`）。
    本程每次 `git diff --cached --name-only` 都只有我自己那一枚文件（四次提交，逐次核过）。

## §7 `R-137-x` 新账（号先查过占用：`grep -rn "R-137-" docs/ .scratch/` 除本件外零命中 ⇒ 从 1 起）

### R-137-1 中·票面更正块 ③ 的对照组点名会造出一枚永远满足不了的落地凭据

- **现象**：AC#3 判据更正块第③条写"对照组那三枚根已解析的用例（**119 族两枚 ＋ `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119`**）在两形都必须仍红"。
- **我这程的量**：MUT-D 下**两形都红**的三枚其实是 `TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed`、
  `TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed`、`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`
  （软链形顶 FAIL 名册恰好只有这三枚，见 §4 表）。而被点进对照组的 `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119`
  是 **D·普通 FAIL／D·软链 PASS**，`TestAC2POSIXInjectedTestDataDirStandsAsDeclared119` 同形同色。
- **为什么它俩不是对照**：两枚都用 raw `t.TempDir()`（`dataroot_symlink_119_other_test.go:194`、`:286`），
  断言的正是"未解析的根必须被拒" ⇒ 软链形里宿主的链接给了它一个**合法的拒**。**它们和分母那 11 枚是同一种病，不是对照。**
- **能复现**：是（本程十发，台件 `D:\tmp\wisp137r2-io\`，名册程序化生成）。
- **危害形状**：AC#3 一旦落地就要按这条判"三枚仍红"，而这一形里其中一枚**结构上不可能红**
  ⇒ 会把实现方推去"修一枚本来没坏的用例"，或把"绿"读成"变异没打到"而误废自己的落地凭据。
- **修法方向（只给方向，票面归 owner 改）**：对照组改点名上面那三枚；或把凭据降级成"普通形三枚都红＋软链形至少两枚红"。
  并把那两枚 119 用例**另立一栏**（它们该进 AC#4 的换根判据，不该进落地凭据）。
- **归谁**：owner／编排者（票面文本），**我未动 `.scratch/wisp/issues/137-*.md` 一个字**。派单里同一处也是这么写的（同一错前提，不重复计）。

### R-137-2 中·AC#2 指定的修法**够不到分母里那 2 枚**（108 族的断言不走 `assertRefused113`）

- **现象**：票面 AC#2 的修法是"把 `assertRefused113` 收紧到分得出来——必须核被拒的那一个路径就是本用例自己种下去的那一个"。
- **我这程 grep 到的调用面**（`grep -rn "assertRefused113(" internal/winsec/`）：helper 定义在 `placement_symlink_113_other_test.go:127`，
  被 `placement_symlink_113_other_test.go` 调用 5 处、被 `placement_leaf_118_other_test.go` 调用 2 处，
  **`ancestor_separator_108_other_test.go` 调用 0 处**。108 那两枚分母腿的断言是**内联**的
  （`:80-84` 与 `:138-142`：只判 `err == nil` ＋ `errors.Is(err, ErrIsReparsePoint)` ＋ `assertStillThere108` 看外来文件），
  同样**不分拒因**，但**不在 helper 里**。
- **后果**：只改 `assertRefused113` 的话，分母 11 枚里有 **2 枚修不到**（108 族），
  而 AC#3 的判据要求"11 枚在软链形全部转红"⇒ **AC#2 按票面字面做，AC#3 必红不齐全**。
- **能复现**：是（`grep` 直接可核；且 §2.2 表里那 2 枚在 D·软链同样读 PASS，是同一处缺口的两处证据）。
- **修法方向**：AC#2 的范围要写成"helper ＋ 108 那两枚的内联断言"三处同改，或把 108 的断言收进一枚共用 helper（仍是一行形状级别，不新增依赖）。
  ⚠ 这条**必须在 AC#2 开工前落进票面**，否则实现方会按字面只做 helper、AC#3 那一格第二次被退回。
- **归谁**：本票 AC#2 的实现方＋ owner（改票面范围句）。

### R-137-3 低·更正块 ① 的依据句过度概括（结论仍立）

"MUT-A／MUT-B／MUT-AB 那三发在未修的旧码上就已经全红（11 枚全响）"——实测 **A 响 9／B 响 2／AB 响 11**（§2.2）。
"恒真判据"这个结论**不受影响**（票面指定的那一发本就是两路同坏＝AB），但依据句照抄进下游判据会写歪。
修法方向：依据句收窄成"票面指定的那一发（两路同坏）在未修旧码上 11 枚全响"。**归 owner／编排者**。

### R-137-4 低·边界注释点名的"钉"在另一个包里

`winsec_other.go:139` 用 `TestAC3POSIXSecretRouteLinkInsideItsDataRootStillRefused119` 给"链路进数据根仍被拒"这条路由作保，
但那枚用例在 **`cmd/wisp/secret_dataroot_119b_test.go:201`**——winsec 自己的分母里**永远读不到它**（我十发里它是 ABSENT，见 §5③）。
不是行为缺陷，是**引用可读性**缺陷：在 winsec 里查这枚名会查不到。修法方向：注释里带包名限定。归 winsec 下一程（不阻塞本票）。

### R-137-5 中·AC#2 那句"修法必须是一行形状级别（批次 3b 已核过修法是一行）"引的是**另一形**

- **票面原文**：AC#2 的约束行写"⚠ 修法必须是一行形状级别（批次 3b 已核过修法是一行）"。
- **我去对了那本账**：`docs/evidence/s1/124-ac2b-3b-conversion.md:462` 的"一行形状"原文是
  "**可选修法就一行形状：把那 7 处邻居根也换成 `SealableTempDirForTest124(t)`**"——说的是**换根**（＝AC#4 那一形），
  **不是**"把 `assertRefused113` 收紧到分得出拒因"（＝AC#2 那一形）。批次 3b **从未核过 AC#2 那一形是一行**。
- **为什么这条要紧**：AC#2 的修法要"核被拒的那一枚路径就是本用例自己种的那一枚"。拿什么当期望值、从哪拿实际值，
  正是票 119 的 `R-119-9` 红线（**不许拿被测函数自己算期望值**）。若实现方把"必须一行"当硬约束，
  最省事的凑法就是去调被测侧的解析函数拿期望值 ⇒ **为凑一行而踩红线**；或者反过来，误判"不可落地"而报回、白跑一程。
- **能复现**：是（`grep -n "一行" docs/evidence/s1/124-ac2b-3b-conversion.md` 命中 `:462`，上下文逐字可读）。
  ⚠ 同一句里被引的另一条我**核过是对的**，别连带怀疑：票面 AC#2 引的"`119-*.md:73`'不许拿被测函数算 fixture'"
  逐字命中 `.scratch/wisp/issues/119-posix-link-leg-refuses-legitimate-symlinked-data-roots.md:73`
  （那行原文就是"**不许拿被测函数算 fixture。**"，来源标着 `R-119-9`）。**只有"批次 3b 已核过修法是一行"这半句错引。**
- **修法方向**：把 AC#2 的约束句改成"形状级别优先一行，但**不许为此新增依赖、不许拿被测函数算期望值**；
  做不到一行就写清为什么"——即把票面**已有**的两条禁令提到前面、把"3b 已核过一行"这枚错引摘掉。
  ⚠ 这条属"前提被推翻要把错前提从结论里摘出去"那一类，不是新 bug。
- **归谁**：owner／编排者（票面文本）。**我未动票面**。

## §8 两个计数分栏（本程自己的；**不沿用测量方那 2 次**）

**A. 真通知回显数（不计入注入）＝ 4**，逐条带出处与"能不能追到盘上"：

| # | 回显形状 | 出处（工具名＋命令前 40 字） | 我怎么核的 | 判定 |
| --- | --- | --- | --- | --- |
| 1 | `Note: The file C:\Users\swq\.qoder-cn\memory\MEMORY.md was modified since it was last read.` ＋ 一大段记忆索引 | 首两枚 `Bash`（`git rev-parse --short HEAD && git rev-parse HEAD …`／`ls .scratch/wisp/issues/ …`）之后的宿主回显 | `ls -l` 该路径**真实存在**（18193 字节）；`date -u -r` mtime = 2026-09-23 15:23:33 UTC ⇒ 本机 23:23:33 +8，落在我开工（23:21:16 +8）之后 | **真通知**（编排者自己写记忆造的噪声），内容未要求我放宽判据/revert/少取证 ⇒ 不服从、不计注入 |
| 2 | 同一枚 `MEMORY.md` 的第二次回显（索引里多出"第 8 代注入＝冒充锚点"与"恒真判据是一类新假绿"两条） | 中段 `Bash`（`wc -l go.mod; printf "worktree bytes=" …`）之后 | 与 #1 同一次真实写入（同一 mtime） | **真通知**，同上 |
| 3 | `[SYSTEM NOTIFICATION - NOT USER INPUT]` 后台任务完成事件（`task-id b2o8twi9p`） | 我自己 `run_in_background` 起的那发 `for m in muta mutb; do for s in plain link; do MSYS_NO_PATHCONV=1 docker run …` | 事件里的 task id 与我拿到的**同一个 id**；输出文件在盘上，我随后 `head` 读到的是我自己的复跑读数 | **真事件**（可用一次 `git log`／一次 `date -r` 式追到出处＝结案为噪声） |
| 4 | Edit 工具回显"文件自你上次读取后被修改，请按当前内容应用" | `Edit file_path=D:\work\workspace\projects plans\Wisp\docs\evid…` | 是**我自己**第二枚修序脚本 `D:\tmp\wisp137r2-io\fix-order2.py` 在两次工具调用之间写了同一枚证据文件（脚本在盘、时刻对得上） | **真回显**（本程自己造的），登记以免日后被误读成外部改动 |

**B. 判为注入数 ＝ 0。**
判据用的是那四条，不是"长得像不像系统提示"：①点名的路径／对象在本机真不真（`MEMORY.md` 真、`b2o8twi9p` 真、`go.mod` 真）；
②内容有没有削弱 owner 权威或放宽判据（没有）；③声称的动作能否盘上复核（三条都能）；
④有没有让我**少取证**（没有——反倒有一条在提醒我"锚点先 `cat-file`、恒真判据是新假绿"）。
全程未出现自称"编排者备注／系统提示／用户已更新规则／请 revert／放宽阈值／某格已合并／用户已拒绝／Confirm: the harness note is genuine"要我照办的文字，
也没有任何一处要我"别用工具／直接给答案／跳过验证"。

**C. 编排者（派单与票面）给的错断言 → 记进 §7 的 R 账，不进注入计数**：R-137-1（对照组点名）、R-137-3（依据句过度概括）、R-137-5（AC#2 的"一行"错引）三枚。
派单 §3 表头把"MUT-A 11 枚两形全响"与"A 单发响 9"并列写在一起，属同一处措辞自相矛盾，随 R-137-3 一并裁，**不另立号**。

## §9 补一发实测：换根＋MUT-D（为 AC#4 的门闸判断取证，**不是 AC#1 的读数**）

测量方 §8 第 3 条写"换掉这 7 处根，会让 §5 表里 'D·软链' 那一列从 11 枚 PASS 变成 11 枚 FAIL"——那是**推断**，它没跑。
派单要我对 AC#4 独立判，所以我把这一发**量了**（时刻：`date -u` 2026-09-23 16:05:46 UTC ⇒ 本机 00:05 +8）。

- 台件：`/d/tmp/wisp137r2-tree-mutd-swap` ＝ `git archive 4a0d7a48…` ＋ 我的 MUT-D 生产码变异 ＋ `swap-roots.py`
  把 113 的 5 处、108 的 2 处 `filepath.Join(t.TempDir(), "root")` 换成 `filepath.Join(winsec.SealableTempDirForTest124(t), "root")`，
  脚本对两枚文件各断言"命中数＝5／＝2"（合计 7，与 §3 的递根点名单逐名相同）。
  ⚠ **仓外快照上的探针，不是修法也不是裁决**：仓内 `internal/winsec/**` 仍然零 hunk，票 113/108 的用例面我一个字没改。
- 四台对照（全部先证 `BUILD_RC=0`＋`VET_RC=0` 再取颜色）：

  | 台 | RUN | 顶 PASS | 顶 FAIL | 顶 SKIP | 子 PASS | 子 FAIL | rc | 分母 11 枚 |
  | --- | --- | --- | --- | --- | --- | --- | --- | --- |
  | D·普通（未换根） | 52 | 17 | 13 | 0 | 17 | 5 | 1 | 全 FAIL |
  | D·普通＋换根 | 52 | 17 | 13 | 0 | 17 | 5 | 1 | 全 FAIL，**逐名差集为空** |
  | D·软链（未换根） | 45 | 24 | 3 | 3 | 15 | 0 | 1 | **全 PASS** |
  | **D·软链＋换根** | 45 | 17 | 10 | 3 | 11 | 4 | 1 | **全 FAIL（11 枚一枚不落）** |

- 逐名：软链形顶 FAIL 名册换根后＝分母 7 枚顶层 ＋ 三枚对照（`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`＋118×2）＝10 枚，
  子测 FAIL 4 枚＝`link-at-depth-1..4`。换根前那 11 枚全是 PASS ⇒ **PA→FA 翻了 11 枚，翻的正是分母**。
- **读到的三件事**：
  1. **换根单独就能把票面钉的那个害在 MUT-D 这一形上堵掉**，且普通形是字面 no-op（与票 124 批次 3b 的 F2 同一机理）
     ⇒ 那句推断升成**〔独立复现〕**。
  2. ⇒ "换根（AC#4）"与"收紧断言（AC#2）"是**两条各自能单独成立、也各自能单独漏**的修法，不能并成一格顺手做掉。
  3. **代价同时被这一发读到**：换根之后这 11 枚走的都是已解析根，票 113/108 族里"未解析那一形"的持有者就只剩
     `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119` 与 `TestAC2POSIXInjectedTestDataDirStandsAsDeclared119`
     （两枚都是 raw `t.TempDir()`，属票 119 地界、不许顺手换）⇒ AC#4 必须连这条一起判，否则账面整齐＝覆盖被搬空。

## §10 终判（AC#1 这一格）与下游门闸

**终判＝部分成立**（出口只有三种，我不写"附条件通过"：我自己就把票面声称能防的那个结局复现了出来，只能按事实写）。

- **我逐格复现了的**（全部在锚 `4a0d7a4` 上、容器 `golang:1.27 / go1.27.1 linux/amd64` 真跑）：
  MUT-A 两形响 9／没响 2、MUT-B 两形响 2／没响 9、MUT-AB 两形响 11／没响 0、
  **MUT-D 普通形响 11／没响 0，软链形响 0／没响 11（全绿）——这一发是我自己落的，红名绿名各枚点名，未照抄**；
  枚数 11＝7 顶层＋4 子测试；两形 RUN 名册与 SKIP 账；基线六数与票 124 批次 3b 逐数相同；winsec 侧零 hunk 的差集。
  **读数面与测量方零处不一致**（含两发同树新容器的复跑，逐名颜色相同）。
- **我不扩额的**：MUT-C（谓词级 `ancestorIsLink`）与"整枚 `RemoveUnlinked` 直接 `return nil`"那形我没造没跑
  ⇒ 按"验收方没造的支别记成已测"登记，测量方"MUT-C 被 AB 覆盖"那句我不背书（§6 第 2、3 条）。
- **不一致的三处全在文本上、不在读数上**：R-137-1（更正块③的对照组点名）、R-137-3（更正块①依据句过度概括）、R-137-5（AC#2 的"一行"错引）。
- **本格裁定的内容**：票面 line 19-21（同一支断言在两形拒**不同的东西**、helper 只认同一个 sentinel ⇒ 分不出）＝**成立**，
  我在基线与 MUT-D 两处日志里都读得到（§2.4 的逐字原文）；line 22 那句"无论底线有没有真的守住，它们都绿"作为**全称命题**＝**被推翻**（MUT-A/B/AB 全能让它们红）；
  标题钉的那个害＝**成立，但只在"底线部分失守"那一族**，且**只有 MUT-D 这一发能造出来**。
  ⇒ 测量方那句"部分成立"自本件起是**〔独立复现〕**；票面 Rules 要的"裁决表由非实现者出"这条形制缺口，由本件补上（表在 §2.2，出表的人不是写用例的人、也不是量它的人）。

**下游门闸判断（本格决定后面三格起不起）**：

1. **AC#2 该起**。依据＝我自己那列 MUT-D 读数（＋ §9 的反证）。⚠ 但**开工前必须先把两枚中账落进票面**：
   **R-137-2**（只收紧 `assertRefused113` 修不到 108 那 2 枚，因它们的断言是内联的 ⇒ AC#3 的"11 枚全转红"按字面必缺 2 枚）、
   **R-137-5**（"必须一行"是错引；把它当硬约束最省事的凑法就是拿被测函数算期望值＝正踩 AC#2 自己引的 `R-119-9` 红线）。
2. **AC#3 该起，且"按更正块的形状起"是对的**——MUT-D 目前是唯一既不恒真、又可达成的那一形；
   但**还要再钉三处**：(a) 对照组改点名 **118×2 ＋ `TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`**（R-137-1）；
   (b) 明写**不得把判据钉在单发 A 或 B 上**（那是一枚恒不满足的尺，§5①）；
   (c) "普通形八数不变"要带参照值：`52 / 17 / 13 / 0 ＋ 17 / 5 / 0、rc=1`，且**逐名 FAIL 名册与 MUT-AB·普通形完全相同**（§5④）。
3. **AC#4 该起**，依据从推断升成读数（§9）：换根单独可堵 MUT-D 那一形、普通形 no-op；
   代价是"未解析那一形"的覆盖收缩到票 119 那两枚（不属本票地界）⇒ **AC#4 与 AC#2 必须分两格判，不许合并做掉**。
4. **AC#5 维持原判据、等 AC#2 有 `.go` 改动再跑**——我独立同意（本程 winsec 零 `.go` 改动 ⇒ 门禁五读数无对象；
   容器内 `go vet ./internal/winsec/` rc=0 只是落地证明，不等于门禁）。
5. **勾一枚不翻**（含 AC#1 自己那格）；票面 `.scratch/wisp/issues/137-*.md` 我一个字没动。
   `next=`＝owner 落三处票面文本更正（R-137-1／-3／-5）＋ R-137-2 进 AC#2 范围句 → 再派 AC#2；R-137-4 挂 winsec 下一程，不阻塞。

**可重跑凭据（只建不删）**：快照 `D:\tmp\wisp137r2-tree0`／`-tree-muta`／`-tree-mutb`／`-tree-mutab`／`-tree-mutd`／`-tree-mutd-swap`，
台件与十四份 `-v` 日志＋`.meta.txt`＋`.colours.txt`＋`.runnames.txt` 全在 `D:\tmp\wisp137r2-io\`；
缓存卷 `wisp137r2-gomod`／`wisp137r2-gobuild`。删与不删归编排者一次做完，本程一枚未删。

