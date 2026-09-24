# 137 AC#3 —— 验收方第二程（r2）：只裁 AC#3 那一格，读数全部我自己打

> 本格裁决者＝**非实现者**（`SPEC-12 §4.3` #1/#3、D22 双角色）。
> 被验版本＝**`6f3817a`**（AC#2 的修法，且 `git diff --name-only 4e66817 HEAD -- internal/winsec/` 为空 ⇒ 与终裁件 `137-ac2-ac5-r1-acceptance.md` 读的是同一版码）。
> 本件**只裁 AC#3**。票 137 其余各格（AC#1／AC#2／AC#4／AC#5）一枚未裁、一枚未翻勾。
> 生产码零改动、票面零改动、`docs/reports/**` 零改动。

---

## §0 争用闸门逐轮原样／锚点自量／取件／三档口径

### 0.1 争用闸门（本机 self-hosted runner 与 CI 会抢 CPU；**本程不 push**）

**闸门第 1 轮**（`date` 原文 `Thu Sep 24 13:34:37 CST 2026`，开工后第一件事）：

```
$ tasklist.exe //FI "IMAGENAME eq Runner.Worker.exe"
INFO: No tasks are running which match the specified criteria.
$ tasklist.exe //FI "IMAGENAME eq go.exe"
INFO: No tasks are running which match the specified criteria.
$ tasklist.exe //FI "IMAGENAME eq compile.exe"
INFO: No tasks are running which match the specified criteria.
$ tasklist.exe //FI "IMAGENAME eq cgo.exe"
INFO: No tasks are running which match the specified criteria.

$ gh run list --limit 5 --json databaseId,status,headSha
[{"databaseId":35958260537,"headSha":"c8967b8bb90f93dc0d7e14a8825ce714fdcc8364","status":"completed"},
 {"databaseId":35953906746,"headSha":"c8967b8bb90f93dc0d7e14a8825ce714fdcc8364","status":"completed"},
 {"databaseId":35949710286,"headSha":"45623e406dc90cd74de0b2262bddcdf6c4dc9e03","status":"completed"},
 {"databaseId":35948947994,"headSha":"8369b24e20eb85b5c9afb45198499dc4b9d9d804","status":"completed"},
 {"databaseId":35947125867,"headSha":"191f0d6844bd177f17a95d3fee9cf84475994ac6","status":"completed"}]
```

（上面那串 JSON 是 `gh` 的单行原样输出，为排版在 `},{"` 处折行，字段一字未改。**`gh` 这一程可用**：
第 1、2 轮都取到了返回值 ⇒ 不需要降级成 `CLEAR-LOCAL-ONLY`；`in_progress` 命中 **0** 枚。）

**闸门第 2 轮**（`date` 原文 `Thu Sep 24 13:40:08 CST 2026`，批量之前）：

```
$ gh run list --limit 5 --json databaseId,status,headSha | grep -c in_progress
0
$ gh run list --limit 3 --json databaseId,status
[{"databaseId":35958260537,"status":"completed"},{"databaseId":35953906746,"status":"completed"},{"databaseId":35949710286,"status":"completed"}]
```

**闸门第 3 轮**（`date` 原文 `Thu Sep 24 13:40:28 CST 2026`，与第 2 轮同一批、补 runner 工作目录那一行）：

```
$ for p in Runner.Worker.exe go.exe compile.exe cgo.exe docker.exe; do tasklist.exe //FI "IMAGENAME eq $p" | grep -c $p; done
Runner.Worker.exe: 0
go.exe: 0
compile.exe: 0
cgo.exe: 0
docker.exe: 1        <- 是我自己那一发 `docker volume ls` 的残留 CLI 进程，不是 CI；已如实记下，未据此判"编队为空"
$ ls -dl /e/work/base/actions-runner/_work
drwxr-xr-x 1 swq 197609 0 Sep 24 12:02 /e/work/base/actions-runner/_work
$ find /e/work/base/actions-runner/_work -maxdepth 2 -newermt "-40 minutes"   # 空输出
（无输出）
```

⇒ 三轮结论：**没有 CI 在飞**（进程 0 枚、`in_progress` 0 枚、runner `_work` 最后动过是 12:02、近 40 分钟零改动）。
批量读数**开局即三行都清**，之后 §2 每一格之前另补一轮并原样贴在那一格末尾。
⚠ 第 1 轮**漏了 runner 工作目录那一行**（当时只查了进程与 `gh`），第 3 轮才补上 ⇒ 如实登记，不回头把第 1 轮写成"三行齐"。

### 0.2 锚点自量＋引用到的每枚 sha 逐枚 `git cat-file -t`

```
$ git rev-parse HEAD                 # 开工时自量
21fed319685ae7b9a429b1dfb410b6b942d48c4e
$ git status --porcelain=v1          # 开工时工作树
（空）
$ git branch --show-current
dev
```

| sha | 我要用它做什么 | `git cat-file -t <sha>` 现量 |
|---|---|---|
| `21fed31` | 开工锚点（我自己 `rev-parse` 量的，不是简报给的） | `commit` |
| `6f3817a` | **被验版本**＝AC#2 的修法 | `commit`（全名 `6f3817ab9110aa3618779696f3bd590b09427f3b`） |
| `f53ad5c` | 对照版本＝AC#2 之前那版（"未收紧"） | `commit` |
| `4e66817` | 上一程终裁件的锚，用它核"winsec 未再动过" | `commit` |
| `1d38206` | AC#1 那程的锚（只出现在我引用的旧 log 里） | `commit` |
| `2d029f9` | 台账 `A155` 点名的上一枚文档 commit | `commit` |
| `99319c7` | **作废的那一程**自述"逐枚 commit 到 99319c7" | `fatal: Not a valid object name 99319c7` ⇒ **不存在**，与 `A155` 的复量一致 |

版本同一性三条（我现量，不是引简报）：

```
$ git diff --name-only 4e66817 HEAD -- internal/winsec/
（空 ⇒ 终裁件那批归档读数的就是今天这版 winsec 码）
$ git show --name-only 6f3817a | tail -4
internal/winsec/ancestor_separator_108_other_test.go
internal/winsec/placement_leaf_118_other_test.go
internal/winsec/placement_symlink_113_other_test.go
$ git diff --name-only f53ad5c 6f3817a | grep internal/winsec
（只有上面那三枚 *_other_test.go ⇒ AC#2 没动生产码）
$ head -1 internal/winsec/placement_symlink_113_other_test.go   # 从 6f3817a 的归档树里读
//go:build !windows
```

⇒ **AC#3 的每一枚读数只能在 linux 侧取**（宿主上这三枚文件连编译都不发生）；这与简报里那句前提一致，且我这一程**全部读数出自 `golang:1.27` 容器**。

### 0.3 取件（严格按 sha 归档，**没把工作树当被验版本**）

```
$ mkdir -p /d/tmp/wisp137ac3r2-tree-loose /d/tmp/wisp137ac3r2-tree-tight
$ git archive f53ad5c | tar -x -C /d/tmp/wisp137ac3r2-tree-loose
$ git archive 6f3817a | tar -x -C /d/tmp/wisp137ac3r2-tree-tight
$ wc -c < .../tree-loose/go.mod  -> 883        # 挂载硬闸用的就是这个字节数（不是上一程的 883 巧合：两棵树都是 883）
$ wc -c < .../tree-tight/go.mod -> 883
```

五枚依赖文件的 `md5sum`（我现量；`git cat-file -s <sha>:go.mod` 的裸 blob 是 855 字节，`git archive` 落盘后 883 ⇒ 那 28 字节是行尾，**同一棵树的两个形状我逐发都对 883**）：

| 文件 | loose（`f53ad5c`） | tight（`6f3817a`，＝HEAD 的 winsec） |
|---|---|---|
| `winsec.go` | `a6144c880de80e43bb1393f3624e7221` | `a6144c880de80e43bb1393f3624e7221`（同） |
| `winsec_other.go` | `b5056918be4ed13817d236fbcae0f477` | `b5056918be4ed13817d236fbcae0f477`（同） |
| `placement_symlink_113_other_test.go` | `ce181f23bafdf99b7f1513e1ee4ccaef` | **`cb8350cb5e5917a685c10b39d57d5275`** |
| `ancestor_separator_108_other_test.go` | `0ae03ee40b9d44ff0b4e5334d27af1bc` | `20c141df792de655408138cdaebf633b` |
| `placement_leaf_118_other_test.go` | `b84c85c87fcb2fbbce1393cc5fcabb87` | `50d4e54c72167d969b361c092cb3e644` |

⇒ 与终裁件 §0 报的四个值逐字相同（tight 三枚＋ loose 三枚它都列了）⇒ **三程读的是同一版码**〔日志＋归档，抽验，我这轮的 md5 是自己打的〕。
生产码两枚在 loose／tight 完全相同 ⇒ **AC#2 那一刀只落在测试件**，这一程的"收紧 vs 未收紧"两列之间**唯一的自变量就是那三处新判据分支**。

变异树（**原件不动**，全部在 `cp -a` 出来的工作副本上做）：

| 树 | 怎么来的 | 生产码两枚 md5 | 判据分支活数 |
|---|---|---|---|
| `wisp137ac3r2-work-d-loose` | loose ＋ 我的 MUT-D | `62da6bcd2cee05d293251d2f1b8f3bbf`／`5361bc3f79a95dddfa3f32dd0228c70d` | 0（`f53ad5c` 本来就没有新判据） |
| `wisp137ac3r2-work-d-tight` | tight ＋ 我的 MUT-D | 同上（逐字相同） | 3 |
| `wisp137ac3r2-work-relax` | tight ＋ **同一发** MUT-D ＋ 把三处新判据分支删掉 | 同上（逐字相同） | 0 |

`mutate.py` 打印的落地指纹（逐字，节选）：

```
PATCHED .../work-d-tight/internal/winsec/winsec_other.go
PATCHED .../work-d-tight/internal/winsec/winsec.go
MD5 62da6bcd2cee05d293251d2f1b8f3bbf    10445 internal/winsec/winsec_other.go
MD5 5361bc3f79a95dddfa3f32dd0228c70d    17537 internal/winsec/winsec.go
MD5 cb8350cb5e5917a685c10b39d57d5275    17797 internal/winsec/placement_symlink_113_other_test.go
...
RELAXED .../work-relax/internal/winsec/placement_symlink_113_other_test.go branches_removed=1
RELAXED .../work-relax/internal/winsec/ancestor_separator_108_other_test.go branches_removed=2
RELAX_TOTAL_BRANCHES=3 (expected 3)
```

（`d-loose` 那 113 枚文件的 md5 停在收紧前的 `ce181f23…`、`d-tight` 停在 `cb8350cb…` ⇒ 两列**只有判据分支这一枚变量**；
`relax` 的 113／108 md5 各变成 `efb6bc46…`／`bf55e86d…` 而 118 与两枚生产码不变。）

### 0.4 容器（照上一程终裁件 §0 那台件抄，数字我自己量）

- 镜像 `golang:1.27`（`docker image ls golang` 现量 id `3680233e3204`、DISK USAGE 1.31GB、本机已有、**未拉网络件**）；容器内 `go version go1.27.1 linux/amd64`。
- `docker version --format 'SERVER={{.Server.Version}}'` ⇒ `SERVER=29.6.2`、`rc=0`（简报那句"server 29.6.2"我重跑过，成立）。
- 挂载 `MSYS_NO_PATHCONV=1` ＋ `/d/...`；**每一发**开跑前 `[ "$(wc -c < /src/go.mod)" = "883" ] || exit 97` ＋ `ls -l`＋`md5sum /src/go.mod`。
- 两形硬断言＝**我自己的**路径名（`/r2ac3link` → `/r2ac3priv`、普通形断言 `/r2ac3link` 根本不许存在）＝ 98。
- **落地不过证就不取颜色**：`go build ./...` 与 `go vet ./internal/winsec/` 任一非 0 ⇒ `exit 95`。
- ⚠ 我比上一程多加一道**变异状态硬闸**（`exit 96`）：每发开跑前在容器里数三样东西 —— `MUTATION-137-AC3R2-D` 标记枚数、
  `RELAX-137-AC3R2` 标记枚数、**活着的** `[!]refusalCreditsLink137(` 判据分支枚数，必须等于该发声明的形状
  （`mutd`＝2/0/3、`unmut`＝0/0/3、`mutd-relax`＝2/3/0）。这一道是为了防"读了错树却当成某一发记账"。
- 模块缓存：`GOPROXY=off`，**不复用别人的卷**（`wisp137r2-*` 那一族按简报只读）。我新建
  `wisp137ac3r2-gomod`（一次性 `docker run --rm -v wisp137r2-gomod:/from:ro -v wisp137ac3r2-gomod:/to ... cp -a` 拷出来，
  源卷挂 `:ro`）与 `wisp137ac3r2-gobuild`（空卷）。拷贝那发打印：`726M /from` → `CP_RC=0` → `726M /to`。

### 0.5 本程建的临时件（**只建不删**，路径全列；`wisp137r2-*` 那一族一枚未写一枚未删）

| 路径（`D:\tmp\` 下） | 内容 |
|---|---|
| `wisp137ac3r2-rig/` | 我的台件 `mutate.py`／`run.sh`／`go.sh`／`batch.sh`／`parse.py` ＋ 每发容器头 `s*.head.txt` ＋ `-v` 日志 `s*.v.log` ＋名册 `s*.{run,fail,skip,colour}.txt` ＋ `batch.out`／`parse.out` |
| `wisp137ac3r2-tree-loose` | `f53ad5c` 纯净树（未收紧） |
| `wisp137ac3r2-tree-tight` | `6f3817a` 纯净树（＝HEAD 的 winsec） |
| `wisp137ac3r2-work-d-loose` | 上 loose ＋ MUT-D |
| `wisp137ac3r2-work-d-tight` | 上 tight ＋ MUT-D |
| `wisp137ac3r2-work-relax` | 上 tight ＋ MUT-D ＋ 删三处新判据分支（**验收侧反向探针，不是修法**） |
| `wisp137ac3r2-io/` | ⚠ **一枚我自己的操作失误**：建目录后我把 `git archive f53ad5c` 错解进了这枚本该放台件的目录（随后台件改放 `wisp137ac3r2-rig/`）。它现在是"另一份 loose 树"，**没有任何读数取自它**；按"临时件只建不删"留着不删，在此登记。 |
| docker 卷 | `wisp137ac3r2-gomod`（从 `wisp137r2-gomod` 只读拷出）、`wisp137ac3r2-gobuild`（新建空卷） |

### 0.6 三档口径声明（本件每一格都会标）

- 〔**独立复现**〕＝我自己在这棵／这发容器里跑出来的读数（§2 全部十发属于这一档）。
- 〔**日志＋归档，抽验**〕＝引上一程／实现方读数，但我抽过且能对上自己的指纹（§0.3 的 md5 对点、§1 的票面原文引用）。
- 〔**仅自述，不背书**〕＝只有别人说过、我没复算也不打算据此下判的（本件如果出现会逐条点名）。
- ⚠ `A155` 那一程（`acceptor-ticket137-ac3-r1`）的**任何**自述（"47 棵树"／"66 发读数"／"逐枚 commit 到 `99319c7`"／"this host has no Docker"）
  **一枚不引用、不背书、不当地基**；`git cat-file -t 99319c7` 我这一程又量了一次＝不存在（见 §0.2）。

---

## §1 判据本体（一律引票面原文，不引转述）

票面＝`.scratch/wisp/issues/137-winsecs-12-posix-rejection-legs-go-green-in-the-symlink-shape-for-the-hosts-own-varlink-not-their-own-planted-link.md`。

### 1.1 AC#3 原句（`:41-42`，**已被下面的更正块改形，不单独生效**）

> - [ ] **AC#3** 修完之后**复跑 AC#1 那一发同一变异** ⇒ 这 12 枚必须**转红**（红名逐名），并且**普通形零新增红**（八数逐数相同）。
>       这一格是本票的牙：**没有 AC#3 的 AC#2 只是把断言改了个形状**。

### 1.2 判据更正块（`:44-53`，**这一块才是有效判据；原句"复跑 AC#1 那一发同一变异"已判为恒真判据，本程不复活它**）

> **AC#3 判据更正（编排者 09-23 23:2x，来源＝`worker-ticket137-ac1` 交件；append-only、上面原文不抹，冲突以本块为准）**
> 原句"复跑 AC#1 那一发同一变异"是一枚**恒真判据**：AC#1 实测 MUT-A／MUT-B／MUT-AB 那三发在**未修的旧码上就已经全红**
> （11 枚全响）⇒ 拿它当"修好了才算数"的尺子，**改与不改都会绿**，本格会变成一个装饰。
> **正确形状（终裁与后续实现方一律按这条判）**：AC#3 必须复跑 **MUT-D 那一形**（两条走查都只走前 3 个组件、
> 仍然拒宿主的 `/varlink`、但**永远走不到用例自己种下去的那枚链接**）⇒
> ①**软链形 11 枚必须全部转红**（红名逐名）；②**普通形八数逐数不变**（不许新增红、也不许由绿转 SKIP）；
> ③对照组那三枚根已解析的用例（119 族两枚 ＋ `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119`）在两形都必须仍红，
> 用来证明变异**真的落地了**、那 11 枚的绿不是"没打到"。
> 依据：AC#1 量到 MUT-D 在**未修码**上的形状恰好是"软链形响 0／没响 11（全绿），普通形响 11／没响 0"——
> **这正是本票标题钉的那个害**，所以只有它能把 AC#2 的修法与"没修"区分开。

⚠ 这个块里 ③ 那句**对照组点名是错的**，被同一枚票的 `R-137-3` 更正式点改过（见 1.4）。**本程按 `R-137-3` 点名，不按 ③ 的原句。**

### 1.3 枚数更正块（`:55-60`，有效）

> **枚数更正（同上来源；本票标题与正文里的"12 枚＝8 顶层＋4 子测试"是错的，AC#1 正文那句"11 枚"反而是对的）**
> 真实枚数＝**11 枚 ＝ 7 顶层 ＋ 4 子测试**。那个"8"可复现但**不可采信**：
> `placement_symlink_113_other_test.go` 确有 8 枚顶层 `func Test`，但其中 **3 枚是要"成功"不要"拒"的反向腿**（归 AC#3 那一族），
> 而 108 真正的 **2 枚**顶层拒绝腿被漏掉 ⇒ 7＋4＝11。`assertRefused113` 的真位置＝`placement_symlink_113_other_test.go:127`。

⚠ 这一程我另外读到一处**行号腐坏**（不影响判据、影响引用）：票面写 `assertRefused113` 在 `:127`，
在**被验版本 `6f3817a`** 上它的真位置是 `placement_symlink_113_other_test.go:174`（函数签名行），
`:127` 落在收紧前那版文件的注释区里。逐枚 `grep -n` 见 §2.1。

### 1.4 三条钉（票面 log `:198-204`，**逐条照抄进本程判据，不自创第四条**）

> **一，只钉 MUT-D 这一发**（票面 `:45` 那条 `>` 更正块已把"复跑 AC#1 同一变异"判成恒真判据，别复活它）；
> **二，软链形"转红"今天不能当凭据**——终裁方独立量到：**未变异的收紧态在软链形里也是那 11 枚红**
> （宿主的链接在第 2 个组件就被拒掉整条拼写）⇒ 转红这件事**修之前就在响**，拿它当"修好了的证据"就是装饰；
> **三，对照组按 `R-137-3` 点名**，并**另禁一条**：不许拿"单发 A／B 里那枚新判据响不响"当凭据。

`R-137-3`（票面 `:148-152`）对对照组的更正式点名：

> **R-137-3（改我上一块更正里的对照组，我又错一次）** —— 我把 `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119` 点进"两形都红的三枚"，**实测不成立**：
> 它在 MUT-D 下 **普通形 FAIL／软链形 PASS**，因为**它自己也用 raw `t.TempDir()`**、与分母那 11 枚同病 ⇒ 拿它当"变异真落地"的凭据会**指错方向**。
> 真正的对照组＝`TestAC118POSIX…SealFile`、`TestAC118POSIX…PrivateFile`、`TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`。

`R-137-1`（票面 `:144-147`）对"禁钉单发"的原文：

> **R-137-1（收窄我"恒真判据"那条的依据句）** —— 结论**立**（票面指定那形 MUT-AB 在未修旧码上确实 11 枚两形全响 ⇒ AC#3 沿用原句必成装饰），
> 但我那句"**三发都已全响**"**过度概括**：实测 **MUT-A 响 9／没响 2、MUT-B 响 2／没响 9、MUT-AB 响 11**。
> ⇒ 另加一条终裁给的硬话并写进判据：**AC#3 也不许钉在单发 A 或 B 上**——收紧断言后 108（或 113）那几枚修完**永远不会响**，
> 那是一枚**恒不满足的坏尺**（与"恒真"是同一族病、方向相反）。**只有 MUT-D 那一形两边都合法。**

### 1.5 本程据此定的结案格（缺一格不下"成立"）

| 格 | 判据出处 | 我这程怎么量 |
|---|---|---|
| ① 软链形 11 枚全部转红、**逐名**、并核红在 `file:line` | 更正块 ①＋三条钉之一 | 发 `s04-dtight-link`（新判据在位）＋红行出自 `placement_symlink_113_other_test.go:187`／`ancestor_separator_108_other_test.go:90`／`:156` |
| ② 普通形零新增红：八数逐数相同 **且** RUN/FAIL/SKIP 三档名册双向差集为空 | 更正块 ② | `s01-dloose-plain` ↔ `s03-dtight-plain`（同一发 MUT-D 两侧）＋ `s06-relax-plain` 参看 |
| ③ 变异真落地：三枚根已解析对照**两形都仍红**；**再加反向证据**＝同一 MUT-D 下"未收紧"版软链形 11 枚全 PASS | 更正块 ③（按 `R-137-3` 点名）＋简报第 ③ 格 | `s02-dloose-link`（未收紧侧）↔ `s04-dtight-link`（收紧侧）＝**判别对**；对照组名从 `-v` 日志逐名取 |
| ④ panic 账＝0（逐份 `-v` 日志），且说明是全量读数 | 简报第 ④ 格 | 每发日志 `grep -cE '^(panic|fatal error)'` ＋ `COLOUR_LINES` 与两形分母对账 |
| ⑤ 恒真自查正面作答 | 三条钉之二＋`R-137-1` | `s07-tight-link`（未变异＋收紧态）与 `s02-dloose-link`（未收紧＋变异）逐名对点，明确哪几发**只登记、不当修复凭据** |

⚠ 分母 11 枚的**逐名表**我从自己的 `-v` 日志生成（`parse.py` 写 `*.fail.txt`／`*.colour.txt`，不手抄），
名册的构成（哪 7 枚顶层、哪 4 枚子测试）在 §2.1 用 `grep -n` 从被验版本现立。

---

## §2 读数

### 2.0 十发批量、闸门、以及**我这一程被自己的硬闸拒掉的 4 发**

**闸门第 4 轮**（batch 2 之前，`date` 原文 `Thu Sep 24 13:43:18 CST 2026`）：

```
$ for p in Runner.Worker.exe go.exe compile.exe cgo.exe; do tasklist.exe //FI "IMAGENAME eq $p" | grep -c $p; done
Runner.Worker.exe: 0
go.exe: 0
compile.exe: 0
cgo.exe: 0
$ gh run list --limit 5 --json status | grep -c in_progress
0
$ ls -dl /e/work/base/actions-runner/_work
drwxr-xr-x 1 swq 197609 0 Sep 24 12:02 /e/work/base/actions-runner/_work
```

⇒ 十发读数的**两批**（batch 1 用闸门第 2/3 轮、batch 2 用第 4 轮）开局三行都清，**无一批发落在有 CI 在飞的窗口里**。

**先记我自己的一个错（它同时是本程最有价值的一条仪器账）**：batch 1 里那四发"未收紧侧"的读数
（`s01-dloose-plain`／`s02-dloose-link`／`s09-loose-plain`／`s10-loose-link`）**被我自建的 96 号硬闸当场拒掉、零颜色入库**。
原因＝我把"期望的判据分支枚数"写成了 `3`，而 `f53ad5c` 那棵树是 AC#2 **之前**的版本、本来就**一枚新判据都没有**。
`-v` 容器头原样（节选）：

```
TREE=wisp137ac3r2-work-d-loose RUN=s01-dloose-plain SHAPE=plain COUNT=1 EXPECT=mutd
GO=go version go1.27.1 linux/amd64
DATE=2026-09-24 05:40:38 +0000
-rwxrwxrwx 1 root root 883 Sep 24 00:53 /src/go.mod
f6ef661732b1851e5c3db348113cb605  /src/go.mod
5361bc3f79a95dddfa3f32dd0228c70d  /src/internal/winsec/winsec.go
62da6bcd2cee05d293251d2f1b8f3bbf  /src/internal/winsec/winsec_other.go
ce181f23bafdf99b7f1513e1ee4ccaef  /src/internal/winsec/placement_symlink_113_other_test.go
0ae03ee40b9d44ff0b4e5334d27af1bc  /src/internal/winsec/ancestor_separator_108_other_test.go
b84c85c87fcb2fbbce1393cc5fcabb87  /src/internal/winsec/placement_leaf_118_other_test.go
STATE MUT_TAG=2 RELAX_TAG=0 ACTIVE_CRITERION_BRANCHES=0
STATE-GATE-FAILED expected=mutd
```

⇒ 这 4 发**没有产生任何读数、也没被记账成读数**；改形态量（`mutd-loose`／`unmut-loose`）后在 batch 2 重跑，
成品名册即下面的 `s01b`／`s02b`／`s09b`／`s10b`。**这条 96 号闸是这一程新加的**（上一程终裁件只有 97/98/95 三道），
它的作用正是防"读了错树却当成某一发记账"——第一次用就拦住了我自己。

**十发台账**（每发＝一枚**全新** `golang:1.27` 容器；`-count=1 -v ./internal/winsec/`；四数由 `parse.py` 从 `-v` 日志生成，非手抄）：

| 发名 | 树 | 形 | 状态闸打印（`MUT/RELAX/CRIT`） | BUILD／VET／TEST | 四数（RUN＋顶层 PASS/FAIL/SKIP＋子测 PASS/FAIL） | panic |
|---|---|---|---|---|---|---|
| `s01b-dloose-plain` | loose＋MUT-D | 普通 | `2/0/0` | 0／0／**1** | 52 ＋ 17/13/0 ＋ 17/5 | 0 |
| `s02b-dloose-link` | loose＋MUT-D | 软链 | `2/0/0` | 0／0／**1** | 45 ＋ 24/3/3 ＋ 15/0 | 0 |
| `s03-dtight-plain` | tight＋MUT-D | 普通 | `2/0/3` | 0／0／**1** | 52 ＋ 17/13/0 ＋ 17/5 | 0 |
| `s04-dtight-link` | tight＋MUT-D | 软链 | `2/0/3` | 0／0／**1** | 45 ＋ 17/10/3 ＋ 11/4 | 0 |
| `s05-relax-link` | tight＋MUT-D＋**删三处新判据** | 软链 | `2/3/0` | 0／0／**1** | 45 ＋ 24/3/3 ＋ 15/0 | 0 |
| `s06-relax-plain` | 同上 | 普通 | `2/3/0` | 0／0／**1** | 52 ＋ 17/13/0 ＋ 17/5 | 0 |
| `s07-tight-link` | tight（生产码未变异） | 软链 | `0/0/3` | 0／0／**1** | 45 ＋ 20/7/3 ＋ 11/4 | 0 |
| `s08-tight-plain` | tight（未变异） | 普通 | `0/0/3` | 0／0／**0** | 52 ＋ 30/0/0 ＋ 22/0 | 0 |
| `s09b-loose-plain` | loose（未变异） | 普通 | `0/0/0` | 0／0／**0** | 52 ＋ 30/0/0 ＋ 22/0 | 0 |
| `s10b-loose-link` | loose（未变异） | 软链 | `0/0/0` | 0／0／**0** | 45 ＋ 27/0/3 ＋ 15/0 | 0 |

⇒ **十发无一命中 97／98／96／95**（batch 2 那四发与 batch 1 那六发全部 `BUILD_RC=0`＋`VET_RC=0` 之后才取颜色）。
⇒ 与旧账对点〔日志＋归档，抽验 → 我这列数字是**独立复现**〕：
`s02b` 的 `45/24/3/3＋15/0` 与实现方 AC#2 自证的"未收紧＋D 软链"逐数相同；
`s04` 的 `45/17/10/3＋11/4` 与它的"收紧＋D 软链"逐数相同；
`s01b`/`s03`/`s06` 的 `52/17/13/0＋17/5/0、rc=1` 与票面 `:172` 点名的 **AC#3 普通形参照值**逐数相同；
`s08`/`s09b` 的 `52/30/0/0＋22/0、rc=0` 与 `s10b` 的 `45/27/0/3＋15/0` 与 AC#1 那程的基线六数逐数相同。

### 2.1 分母 11 枚的**逐名**构成（在**被验版本 `6f3817a`** 上现立）

`grep -n "^func Test"`（归档树 `wisp137ac3r2-tree-tight`，即 `6f3817a`）读到的 7 枚顶层拒绝腿 ＋
`placement_symlink_113_other_test.go:217` 那个 `for _, depth := range []int{1, 2, 3, 4}` 展开的 4 枚子测试：

| # | 名字（`-v` 日志里的原样） | 断言走哪条 | 源码位置 |
|---|---|---|---|
| 1 | `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone` | helper | `placement_symlink_113_other_test.go:195`，调用点 `:204` |
| 2 | `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth` | （父：由 4 枚子测汇总而红） | `:216` |
| 3–6 | 同上 `/link-at-depth-1`…`/link-at-depth-4` | helper | `t.Run` `:218`，调用点 `:237` |
| 7 | `TestAC1POSIXSealDirThroughASymlinkRefuses` | helper | `:246`，调用点 `:254` |
| 8 | `TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing` | helper | `:262`，调用点 `:271` |
| 9 | `TestAC1POSIXSealFileThroughABackslashNamedLink` | helper | `:287`，调用点 `:296` |
| 10 | `TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink` | **108 内联**（不走 helper） | `ancestor_separator_108_other_test.go:67`，新分支 `:84`、红行 `:90` |
| 11 | `TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor` | **108 内联** | `:133`，新分支 `:151`、红行 `:156` |

⇒ 7 顶层＋4 子测＝**11**，与票面枚数更正块一致（标题那句"12 枚＝8＋4"是错的，本表不采）。
⇒ 新判据的**三处落点**（`placement_symlink_113_other_test.go:186` 的 helper 分支 ＋ 108 的 `:84`／`:151` 两枚内联分支）
在**被验版本**上逐枚 `grep -n '[!]refusalCreditsLink137('` 数到 **3 枚**，与容器状态闸 `ACTIVE_CRITERION_BRANCHES=3` 一致；
`f53ad5c` 那棵树同一命令数到 **0 枚**（这就是 batch 1 被 96 号闸拒掉的原因）。
⚠ 顺带一条行号腐坏（不改判据、只改引用）：票面 `:58` 说 `assertRefused113` 在 `:127`，被验版本上它的签名在 `:174`（`:127` 是收紧前那版的注释区）。

### 2.2 结案格 ①：软链形 11 枚**全部转红**——逐名 ＋ 红在哪一行（发 `s04-dtight-link`）

`s04` 的 FAIL 名册共 **14 枚**＝本格分母 **11** ＋ 对照组 **3**（见 §2.4）。分母那 11 枚逐名（名字取自 `--- FAIL:` 原文行）：

| # | 用例 | `s02b`（未收紧，同发变异） | **`s04`（被验版本）** | `s05`（反向探针：删掉新判据） | `s07`（未变异＋收紧态） | `s10b`（两版都未动） |
|---|---|---|---|---|---|---|
| 1 | `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone` | PASS | **FAIL** | PASS | FAIL | PASS |
| 2 | `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth`（父，因子测而红） | PASS | **FAIL** | PASS | FAIL | PASS |
| 3 | `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth/link-at-depth-1` | PASS | **FAIL** | PASS | FAIL | PASS |
| 4 | `…/link-at-depth-2` | PASS | **FAIL** | PASS | FAIL | PASS |
| 5 | `…/link-at-depth-3` | PASS | **FAIL** | PASS | FAIL | PASS |
| 6 | `…/link-at-depth-4` | PASS | **FAIL** | PASS | FAIL | PASS |
| 7 | `TestAC1POSIXSealDirThroughASymlinkRefuses` | PASS | **FAIL** | PASS | FAIL | PASS |
| 8 | `TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing` | PASS | **FAIL** | PASS | FAIL | PASS |
| 9 | `TestAC1POSIXSealFileThroughABackslashNamedLink` | PASS | **FAIL** | PASS | FAIL | PASS |
| 10 | `TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink`（108 内联） | PASS | **FAIL** | PASS | FAIL | PASS |
| 11 | `TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor`（108 内联） | PASS | **FAIL** | PASS | FAIL | PASS |

（11 行＝§2.1 名册的 11 枚，逐名对得上 `s04-dtight-link.v.fail.txt` 那 14 枚里的 11 枚，另 3 枚是对照组。）

**红在哪一行**（`s04` 里带新判据文案 `does not credit the link this case planted` 的 `-v` 原行，按 `file:line` 计数）：

```
      1 placement_symlink_113_other_test.go:204
      4 placement_symlink_113_other_test.go:237        <- 四枚深度子测（t.Helper() ⇒ 行号落在调用点）
      1 placement_symlink_113_other_test.go:254
      1 placement_symlink_113_other_test.go:271
      1 placement_symlink_113_other_test.go:296
      1 ancestor_separator_108_other_test.go:90
      1 ancestor_separator_108_other_test.go:156
```

⇒ **10 条新判据文案 ＋ 1 枚父用例（`:216`，因子测全红而汇总红）＝ 11 枚**，一枚不多一枚不少，
且**行号逐枚落在那 5 处 helper 调用点与 108 那 2 枚内联腿上**（＝AC#2 那一刀的落点，见 §2.1 表最后一列）。
⇒ 排除"红出自别的分支"：同一份日志里
`returned nil, i.e. it sealed through a symlink` 出现 **2** 次（都在 118 那两枚对照，见 §2.4），
`which does not name ErrUnresolvedPath` 出现 **0** 次 ⇒ 分母 11 枚没有一枚是靠"没拒"或"拒错类型"红的，
**全部红在新加的"这枚链接是不是我自己种的"那一支上**。

红行原文一枚（`s04`，逐字截断到能看清机制为止）：

```
    ancestor_separator_108_other_test.go:90: AC#2 RED: RemoveUnlinked("/r2ac3link/w137ac3rtmp/TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink517751341/002/root/link/keep-me.txt")
      refused with winsec: entry is a link to something else, not private data …: the spelling reaches it through the link at /r2ac3link, and whatever lives behind that link …
      , which does not credit the link this case planted at "/r2ac3link/w137ac3rtmp/…/002/root/link": an ambient link above the tree can answer for the refusal while the guard under test never reaches this spelling
```

⇒ 机制与票面钉的**一模一样**：被记名的是**容器自己那枚** `/r2ac3link`（软链形里 TMPDIR 的第一段），
而用例自己种在 `…/002/root/link` 的那一枚**从未被走过**。这就是 MUT-D 那一形，也就是本票标题钉的害。

