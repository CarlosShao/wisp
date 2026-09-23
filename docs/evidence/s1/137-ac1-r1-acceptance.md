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

