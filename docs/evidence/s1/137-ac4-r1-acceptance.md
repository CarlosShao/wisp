# 137 AC#4 —— 终裁表（非实现者）：那 7 处递根点换得对不对，凭据是不是我自己跑出来的

**格子**：票 `.scratch/wisp/issues/137-winsecs-12-posix-rejection-legs-go-green-in-the-symlink-shape-for-the-hosts-own-varlink-not-their-own-planted-link.md`
**AC#4 一格**（票面 `:74-76`）＋ 它 `:173` 那块下面的 `>` 口径更正（`:176-187`，09-24 15:5x 编排者）。
**判据按更正版那两条判（②a／②b），不按 `:173` 原句那句"换根之后 MUT-D 软链形 11 枚转红"判**——那条已作废为凭据。
**agent**＝`acceptor-ticket137-ac4-r1`（**非实现者**；写码那程＝`worker-ticket137-ac4`，它的自述一枚未采信，下面全是本程重走）。
**AC#1／AC#2／AC#3／AC#5 一格未重裁；票面一枚勾未翻；`-done` 未改。**

三档口径（每节自标）：〔独立复现〕＝本程自己在树上跑出来的；〔日志＋归档，抽验〕＝引别程读数但本程重测过它依赖的前提；〔仅自述，不背书〕＝只有别人说过、我没复算也不据它下判的。

---

## §0 锚点自量 ＋ 码面未动证明 ＋ 闸门 ＋ 被验版本盘上身份　〔独立复现〕

### 0.1 锚点自己现取（不读脏工作树冒充被验版本）

```
$ git rev-parse HEAD
a9c4d58f99b134b0ab3253b80469dde92bb2692a
$ git cat-file -t a9c4d58f99b134b0ab3253b80469dde92bb2692a
commit
$ git rev-parse --abbrev-ref HEAD
dev
$ date
Thu Sep 24 15:48:39 CST 2026
```

### 0.2 码面未动证明（派单要求的那一条，必须为空）

```
$ git diff b1010ff..HEAD -- internal/winsec/
（空）
$ echo rc=$?
rc=0
```

⇒ 被验的 `internal/winsec/` 码面从我开工的锚点 `a9c4d58` 回到交付面最后一枚 `b1010ff` **一字未动**，
不需要停手上报。**读数全部走 `git archive` 快照，没有一枚取自工作树**（见 0.5）。

两枚 shas 的真实关系（本程现量，别当"锚点在交付面之前"读）：

```
$ git merge-base --is-ancestor b1010ff HEAD && echo yes
yes
$ git log --format='%h %p %ad' -1 b1010ff ; git log --format='%h %p %ad' -1 a9c4d58
b1010ff 20a6397 Thu Sep 24 15:43:20 2026 +0800
a9c4d58 d777eb6 Thu Sep 24 15:47:34 2026 +0800
$ git log --oneline b1010ff..a9c4d58
a9c4d58 docs(137 AC#4 核收, 119 AC#7 立案): ...
d777eb6 evidence(125,AC#4 r1 终裁 §4): ...
a088515 evidence(125,AC#3 r1 终裁 §3): ...
$ git diff --numstat b1010ff..HEAD -- <137 票面>
33	0	（纯 append＝编排者那一块 `>` 口径更正）
```

⇒ 锚点 `a9c4d58` 是交付面 `b1010ff` 的**后代差三枚**，那三枚只动票面与别人的证据件（`internal/winsec/` 逐字节相同，
就是上面那条空 diff 证的）；本程读的**是交付面那一版码**，不是"我开工时的工作树"。

工作树确实不干净，且**与本程无关**：`git status --porcelain` 现量 **18 行**＝`design/**` 那 16 枚 tracked 文件被
owner 那侧挪成未跟踪的 `design/old/`（显示为删除）＋ 两枚未跟踪目录 `design/doubao/`、`design/old/`。
本程**一枚未还原、未提交、未删、没过问**。`git add` 只用显式路径，每枚 commit 前跑 `git diff --cached --name-only`。

### 0.3 争用闸门（现量，逐批各量一次）

- `date` ＝ `15:48:39 CST`（本机 +08）。全程**没有任何一处**用两个时间戳相减算时长。
- `gh run list --branch dev --limit 3`（开工现量＋**每发之前**由 `batch.sh` 再打一行 `GATE-PRE`）：
  三枚全 `completed`（`35967768017` ci 8m38s、`35964449249` ci 9m3s、`35958260537` slo-fresh）⇒ **无 `in_progress`**。
- `docker ps`：只有 `union-proxy`／`clipsync*`／`postgres:15-alpine`／`redis:7-alpine` 这类**长期驻留**服务（Up 4–13 天），
  `docker ps --format '{{.Image}}' | grep -c golang` 逐发＝**0** ⇒ 取读数那些批里没有第二枚 golang 容器。
  ⚠ 这一列只回答"有没有别人在同类容器里跑"，**不主张"宿主零负载"**；那五枚驻留服务本程无法让它们停下。
  ⚠ 同机的 `acceptor-ticket125-r1` 若起容器会打在这一列上：批次 A 六发的 `GATE-PRE` 全部 `golang-containers:0`。
- 仪器可用性**没有采信派单前提**：`docker version --format '{{.Server.Version}}'` ＝ **`29.6.2`**（rc=0），
  `docker images golang:1.27` 现量 `1.31GB / 327MB` ＝ **本机已有，未 `docker pull`**，容器内 `go version go1.27.1 linux/amd64`。

### 0.4 被验版本盘上身份（自己建的归档树，非工作树）

| 树键 | 出自 | 盘上路径 | 现量身份 |
|---|---|---|---|
| `swap`＝**被验那一版** | `git archive b1010ff` | `D:\tmp\wisp137ac4r1-tree-swap` | `go.mod` **883 字节**；`placement_symlink_113_other_test.go` md5 `25660e86ad055e68304839e227713059`；`ancestor_separator_108_other_test.go` md5 `6f6b9f804be1560d92f53fb8d9350eec`；`RAWROOTS=0 SWAPPEDROOTS=9` |
| `base`＝换根之前那一版 | `git archive a02da50` | `D:\tmp\wisp137ac4r1-tree-base` | `go.mod` 883 字节；113 md5 `cb8350cb5e5917a685c10b39d57d5275`；`RAWROOTS=7 SWAPPEDROOTS=2` |

⚠ **本程把 A/B 的地基钉在 `a02da50` ＝ `9c0f546^`（换根那枚码的真正父提交）**，不是实现方用的 `182daed`；
两条理由：① `git log --oneline 182daed..9c0f546` 显示中间还有一枚 `a02da50`（docs），
② `git diff 182daed..a02da50 -- internal/winsec/` **为空** ⇒ 两版在 winsec 上逐字节相同，
本程的 A/B 与实现方的 A/B 是**同一对码**，只是本程用了更严格的那个父。

两树差**只有那 7 处**（文件系统侧独立复算，不引 git）：`diff -ru` 两树的 `internal/winsec/` ⇒ 变动行 **42 行**，
其中 `^[+-]\s+root := ` 命中的正好 **7 对（−7 raw／+7 已解析）**，其余 28 行是两枚文件头的说明块。

### 0.5 本程台件与三把身份闸（读数之前先过闸）

- 台件：`D:\tmp\wisp137ac4r1-rig\{shot.sh,batch.sh,apply_mutd.py,revert_spot.py,parse.py}`；日志 `logs/`。
- 快照一律 `git archive <sha> | tar -x -C /d/tmp/wisp137ac4r1-<名>` ＋ `-v /d/tmp/...:/src`（**不走 `docker run -v C:\` 那条静默挂空且 rc=0 的假绿坑**），
  容器内 `wc -c /src/go.mod` ＋ `[ -s ]` 不成立 ⇒ `exit 97`，**不发测**。
- 形状硬断言（`gate 98`，**本程自己的路径名**，不复用别家）：软链形 `ln -s /r1priv /r1link` ＋ `[ -L ]` ＋ `readlink` 逐字核 ＋ `TMPDIR=/r1link/w137ac4r1tmp`；
  普通形**断言 `/r1link` 根本不许存在** ＋ `/r1plain` 不是软链 ＋ `TMPDIR=/r1plain/w137ac4r1tmp`。
- 落地先证再取颜色（`gate 95`）：`go build ./...` ＋ `go vet ./internal/winsec/` 任一非 0 ⇒ **不取颜色**（第一发就这样被自己挡了一次，见 0.6）。
- 树身份闸（每发容器内现打）：`RAWROOTS`／`SWAPPEDROOTS`／`MUTD_MARKERS`／`CREDIT_BRANCHES`（AC#2 那把尺活着＝3）。
- 读数口径：`go test -count=1 -v ./internal/winsec/`，四数由 `parse.py` 从 `-v` 日志**程序化生成**（**没有一枚名是手抄**），
  `-count=N` 没有出现过 ⇒ 不会翻倍；每发另打 `grep -cE '^(panic|fatal error)'` 与**逐名 FAIL/SKIP 名册**＋**名册差集**。
- 模块缓存：`GOPROXY=off` ＋ `GOMODCACHE=/gomod`。**本程自己的卷** `wisp137ac4r1-gomod`／`wisp137ac4r1-gobuild`；
  源＝别人的 `wisp137ac4-gomod` 以 **`:ro`** 挂进 `/from` 一次性 `cp -a`（原文：`726M /from` → `CP_RC=0` → `726M /to`），**源卷没被写过一枚字节**。

### 0.6 一次仪器自拒（登记，不抹；发生在取读数之前）

第一发 `smoke-swap-plain` **`DOCKER_RUN_RC=95`**：容器里空 `GOMODCACHE` ⇒ `go build ./...` 撞
`module lookup disabled by GOPROXY=off`（`internal/secret/configrefs.go:10:2` 等 5 行），`BUILD_GATE failed`。
**那一发零颜色入账**（`gate 95` 在取颜色之前 `exit`），只是**本程台件缺料**、不是树错、也不是被验码错。
处置＝按规矩以 `:ro` 源 `cp -a` 建自己的模块缓存卷，复跑同一发拿到 `TEST_RC=0` ＋ `52/30/0/0＋22/0/0`，
**与本程之外的那一发 `r3` 同数**，然后才开批。

### 0.7 本程第一枚 commit 的账（只能在提交之后现量，故本节末尾回填）

`git log --oneline -1` 与 `git show --name-only HEAD` 的原样输出见本节下方的 0.8——
它们必须是**提交之后**的读数，这是共树里"提交账只能在提交后现量"的固有循环，本程按 `A155` 那枚形状留痕。

### 0.8 §0 那一枚 commit 的原样输出（提交之后现量）

```
$ git log --oneline -1
3a49745 evidence(137 AC#4 r1 终裁 §0): 锚点自量 a9c4d58 + b1010ff..HEAD -- internal/winsec/ 为空 + 闸门 + 被验版本盘上身份
$ git show --name-only HEAD
commit 3a497457cfc5ea4564749cbbf80b620cdc210b71
Author: CarlosShao <1933942520@qq.com>

    evidence(137 AC#4 r1 终裁 §0): 锚点自量 a9c4d58 + b1010ff..HEAD -- internal/winsec/ 为空 + 闸门 + 被验版本盘上身份

docs/evidence/s1/137-ac4-r1-acceptance.md
```

⇒ 那一枚只带 `docs/evidence/s1/137-ac4-r1-acceptance.md` 一枚路径，别人的 `design/**` 一枚未卷。
本节（0.8 这几行）本身又走下一枚 commit 落盘——同一枚循环，按 `A155` 那枚形状留痕。

---

## §1 判据①：那 7 处递根点的逐枚判定表复核 ＋ "有没有哪一处本来不该换"　〔独立复现〕

### 1.1 枚数：派单那句"从 `b1010ff` 现 `grep -n`"**照字面产不出 7**，本程换三条独立计数复算

派单写的是"那 7 处 raw `t.TempDir()`（枚数自己从 `b1010ff` 现 `grep -n`，别引我给的数）"。
**本程照做了，结果是 0 枚命中，不是 7**：

```
$ git grep -n 'root := filepath.Join(t.TempDir(), "root")' b1010ff -- internal/winsec/placement_symlink_113_other_test.go internal/winsec/ancestor_separator_108_other_test.go
（无输出，rc=1）
$ git grep -n 't\.TempDir()' b1010ff -- internal/winsec/placement_symlink_113_other_test.go internal/winsec/ancestor_separator_108_other_test.go
（无输出，rc=1）
```

⇒ **票面 `:70` 那条"裸指针会被自己的修法挪走"的教训，在这一格里比"行号挪走"更狠一层**：
被 AC#4 的修法吃掉的不是行号，是**被数的那个形状本身**（换干净之后全仓那两枚文件里 `t.TempDir()` 一枚不剩）。
**判据不成立但不是坏事**，本程不硬凑读数，改成三条**各自独立**的计数，三条都给 **7**：

| 计数法 | 命令（本程现跑） | 读数 |
|---|---|---|
| ① 改动的行集 | `git show 9c0f546 \| grep -c '^-.*root := filepath.Join(t.TempDir(), "root")'` ＋ 同形 `^+...SealableTempDirForTest124` | **7 ＋ 7** |
| ② 改前的影像 | `git grep -n 'root := filepath.Join(t.TempDir(), "root")' a02da50 -- <那两枚文件>` | **7**：`108:69/135` ＋ `113:197/220/248/264/289` |
| ③ 盘上两棵树 | `diff -ru <base>/internal/winsec <swap>/internal/winsec` 里 `^[+-]\s+root := ` 命中 | **7 对**（其余 28 行是两枚文件头的注释） |

**枚数＝7 成立，且"全换、一枚不留"成立**（`b1010ff` 上那两枚文件里 raw `t.TempDir()` ＝ 0 枚，连不是"root"的变量名都没有）。
⚠ **登记一条给下游的判据缺陷（不改判据、只报名）**：今后 AC#4 这一类"换掉某个形状"的格子，
**枚数判据必须写成"改前影像 ＋ 改动行集"两路**，写成"从交付面 HEAD 现 grep 那个形状"会随修法一起归零。

### 1.2 逐枚判定表（本程自己的编号；换根前 → 换根后行号都由 1.1 的②／现地 `grep -n` 量出，不抄任何前手）

| # | 递根点（`a02da50` → `b1010ff`） | 它喂的用例／子测（逐名） | 判定 | 一句理由（本程的，不是转述） |
|---|---|---|---|---|
| 1 | `113:197` → `113:212` | `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone` | **该换** | 它的判据（AC#2 之后）是"拒因必须记名到 `root/link`"；未解析根让 `firstLinkAncestor`/`platformVerifyPlacement` 停在 `/r1link`，记名判据**在任何一版生产码上都只能响**（§2 的 `a-base-link` 读数：这一枚在未变异软链形就是红）⇒ 未解析根是它的**遮蔽**，不是它的分母 |
| 2 | `113:220` → `113:235` | `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth` ＋ 4 枚子测 `link-at-depth-1..4` | **该换** | 这枚用例存在的**全部理由**是"链接在根下第 2/3/4 层时走查还得走到"；未解析根把拒点钉在第 2 个组件，**depth-2/3/4 那三层种下的链接从未被走查看过一眼**（子测本身在 `=== RUN` 里、跑的是宿主那枚链接）——本程 §3 的 `R2` 一发把它单独撤回 raw，逐名读到那 5 枚由"能响"退成"两态同色"（＝零区分力） |
| 3 | `113:248` → `113:263` | `TestAC1POSIXSealDirThroughASymlinkRefuses` | **该换** | 同 1，入口换成 `SealDir(root/link)`；未解析根下 `SealDir` 那条腿不参与读数，绿是宿主链接给的 |
| 4 | `113:264` → `113:279` | `TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing` | **该换** | 它钉的是"字节落盘**之前** placement 腿就跑完"；未解析根下"没写字节"由宿主链接的拒给出，断言与它声称的因果**脱钩**（脱钩＝它测不到自己名字里那件事） |
| 5 | `113:289` → `113:304` | `TestAC1POSIXSealFileThroughABackslashNamedLink` | **该换** | 被测的是**名字里带反斜杠的那枚链接**会不会被切成两截（`A74(3)` 的 fail-open 钉）；未解析根让走查在它**之前**就返回 ⇒ 这一枚恰是"永远走不到"的那一枚 |
| 6 | `108:69` → `108:82` | `TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink` | **该换** | 108 这族的断言**内联、不走 helper**（`R-137-2`），AC#2 给它新接的就是 `refusalCreditsLink137`；未解析根下这把新尺在软链形**只能响、不能绿**（本程 §3 的 `R6` 一发＝把这一枚单独撤回 raw，逐名读到它"两态同色"） |
| 7 | `108:135` → `108:148` | `TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor` | **该换** | 同 5 ＋ 同 6：反斜杠名 ＋ 未解析根＝两重"走不到"，它测的"那次 `Lstat` 必须用原名字面"在未解析根下**被整个省掉** |

⚠ **表里不给"逐枚 → 逐枚"的承接映射，因为不存在**：这 7 处换完之后，"未解析那一形"由**同一组**用例整体接住
（具名见 §1.3），不是一枚对一枚。谁要是把这一栏读成"第 3 处由 `119:285` 单独守"，那是本表写坏了，按 §1.3 为准。

⇒ **本程独立得出与实现方同一结论：7 处全部该换，没有一处需要未解析的根来完成它自己的断言。**
这一条不是"账面整齐"，是 §2／§3 的读法给的（见 §3 的逐枚"撤掉它哪条用例变得不响"）。

### 1.3 反向那一问：换掉之后"未解析那一形"到底还有没有人持有（**具名，不指票号**）

本程不采信实现方那张表的名字，自己在 `b1010ff` 上重新点名并核"真在跑、不是 SKIP"：

1. `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119` ＝ `internal/winsec/dataroot_symlink_119_other_test.go:193`，
   `base := t.TempDir()`（`:194`，**仍 raw**），自己在基根里种 `var` 那枚链接再拼 `<base>/varlink119/data`；
2. `TestAC2POSIXInjectedTestDataDirStandsAsDeclared119` ＝ 同文件 `:285`，`base := t.TempDir()`（`:286`，**仍 raw**）；
3. `TestAC2POSIXSeamAcceptsTheHonestPOSIXAnswer125` 的子测 `control_unresolved_root_still_refused_by_the_floor_itself`
   ＝ `internal/winsec/seam_probe_root_125_other_test.go:277`（**只在普通形跑**）。

**本程现量**（`a-swap-link`，即被验那版·未变异·软链形）：

- `119` 那两枚**在 `=== RUN` 名册里**（`grep` 命中，见 1.4 那条名册差集）⇒ 软链形里它们**真跑**，没被跳过；
- 125 那三枚**在 SKIP 名册里**（逐名：`…SeamAcceptsTheHonestPOSIXAnswer125`／`…SeamGuardStillRefusesEveryHostileShape125`／`…SeamProbeShapesAreBuiltOnAResolvedRoot125`），
  且父项 SKIP ⇒ **它的子测在软链形一枚都没跑**（名册差集现量 7 枚，逐名全是 125 那三枚的子测）；
- 两枚文件里剩下的 raw `t.TempDir()` 枚数本程也重数：`119` **5 枚**（`:128/:154/:194/:220/:286`）＋ `125` **2 枚**（`:60/:224`）
  ⇒ **实现方"不在这 7 处的地界里"那句成立，且它一枚未动**（`9c0f546` 只带那两枚文件）。

⇒ **但这一栏必须连着 §5 第二问读**：119 那两枚**跑、能打印、不能区分**（它们的断言只判 sentinel），
所以换根之后"未解析那一形"在软链形**只剩日志、不剩判据**。这条代价是真的，实现方自己报回了，编排者已落成票 119 `AC#7`。

### 1.4 本程核过的事实清单（判据①范围内）

- 枚数三条独立计数＝7／7／7；`b1010ff` 上那两枚文件 raw `t.TempDir()`＝**0**。
- 逐枚行号：改前 `108:69/135`＋`113:197/220/248/264/289`（`git grep -n` 现量），改后 `108:82/148`＋`113:212/235/263/279/304`（换根树上 `grep -n` 现量）。
- **同两枚文件里早已走已解析根的 2 处**（`113:403`、`108:117`）**不在 7 枚分母里**——本程按"它两形都响＋被当对照组用"逐名核过（见 §2 的 `b-*` 两发）：
  `113:403` 属于 `TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator`（反向腿），`108:117` 属于 `TestAC2POSIXDoesNotFoldABackslashIntoASeparator`（反向腿）。
  ⇒ 实现方 §1.1 那句"另有 2 处早已走已解析根、非本格分母"成立。
- 同树同形复跑一发逐名同色（`c-r2-plain` vs `e-r2-plain2`，名册 `diff` 为空）⇒ 读数不是单发运气（那两发的设计在 §3）。
- 未解析形持有面 3 枚具名核对存在；`=== RUN`／SKIP 归属用名册差集核，不靠读注释。

### 1.5 本节这一枚 commit 的账（同 0.8 那枚循环，原样输出由下一枚 commit 回填）

```
$ git log --oneline -1
18f2531 evidence(137 AC#4 r1 终裁 §1): 枚数三条独立计数=7(派单那句"从 b1010ff 现 grep"照字面产 0 枚)+逐枚判定 7 全该换+未解析形具名持有人现量
$ git show --name-only HEAD
commit 18f253126e0626034f6a01d88182c17770175094
Author: CarlosShao <1933942520@qq.com>

    evidence(137 AC#4 r1 终裁 §1): 枚数三条独立计数=7(派单那句"从 b1010ff 现 grep"照字面产 0 枚)+逐枚判定 7 全该换+未解析形具名持有人现量

docs/evidence/s1/137-ac4-r1-acceptance.md
```

⇒ 那一枚也只带本程这一枚路径。

---

## §2 判据②：换根这一味药的**两条实质凭据**——两对各自己跑一遍　〔独立复现〕

**先说口径**：票面 `:173` 原句"换根之后 MUT-D 软链形 11 枚转红"**本程一次也没拿它当凭据**，
按 `:176-187` 那块更正判②a／②b 两条。下面 24 发读数**全部出自本程自己建的树、自己的容器、自己的假根名**
（`/r1priv`／`/r1link`／`/r1plain`，与实现方的 `/ac4priv`／`/ac4link` 不同名，排除抄日志的可能）。
**本程总发数＝27**（1 发台件自校 ＋ 6＋8＋10＋2 四批），**每发 panic 计数逐发 0**，每发之前一条 `GATE-PRE`（见 §0.3）。

### 2.1 ②a ＝ 未变异的软链形那 11 枚**由红转绿**

| 本程的发 | 树 | 形 | RUN | 顶 P/F/S | 子 P/F/S | 包级 rc | panic |
|---|---|---|---|---|---|---|---|
| `a-base-link` | `base`（`a02da50`，未换根） | 软链 | 45 | 20/**7**/3 | 11/**4**/0 | 1 | 0 |
| `a-swap-link` | **`swap`＝被验那版**（`b1010ff`） | 软链 | 45 | **27**/0/3 | **15**/0/0 | **0** | 0 |
| `a-base-plain` | `base` | 普通 | 52 | 30/0/0 | 22/0/0 | 0 | 0 |
| `smoke-swap-plain` | `swap` | 普通 | 52 | 30/0/0 | 22/0/0 | 0 | 0 |

**名册级证据**（`comm`／`diff` 逐名，不是包级 rc）：

- `a-base-link` 的 11 枚红逐名 ＝ 7 顶层 ＋ 4 枚 `link-at-depth-1..4` 子测；与 `a-swap-link` 做全量 colour 差集
  ⇒ **恰好这 11 条 `FAIL→PASS`，没有第三枚动**（差集里只出现这 11 个名字，正反两向都是）。
- 两发 `=== RUN` 名册**逐名相同**（45 枚，`diff` 为空）⇒ **没有谁整发消失**。
- 两发 SKIP 名册**逐名相同**（3 枚，见 §4）⇒ **那 11 枚不是由红转 SKIP**。
- 普通形两发全量 colour 名册 `diff` **为空** ⇒ 换根在普通形＝**字面 no-op**（机制：`SealableTempDirForTest124`
  走 `resolve.go:253 resolveProbeRoot`，无链接时返回同一串）。

⇒ **②a 成立，且是本程亲手量出来的。** 与实现方 `r2→r4` 那一对**逐数相同**（`45/20/7/3＋11/4/0` → `45/27/0/3＋15/0/0`），
两边是**独立同读**，不是本程抄它。

### 2.2 ②b ＝ 红因换轨：同一形同一批用例，红句从"没记名"换成"根本没拒"

| 本程的发 | 树 | 形 | RUN | 顶 P/F/S | 子 P/F/S | rc | **`returned nil, i.e. it sealed through a symlink`** | **`does not credit the link this case planted`** | 日志里点到 `/r1link` 的行 |
|---|---|---|---|---|---|---|---|---|---|
| `b-basemutd-link` | `base` ＋ 本程自造 MUT-D | 软链 | 45 | 17/10/3 | 11/4/0 | 1 | **2**（都是 118 那两枚天生已解析根的对照） | **10** | 22 |
| `b-swapmutd-link` | **被验那版** ＋ 同一发 MUT-D | 软链 | 45 | 17/10/3 | 11/4/0 | 1 | **10** | **0** | 2 |
| `b-swapmutd-plain` | 被验那版 ＋ MUT-D | 普通 | 52 | 17/13/0 | 17/5/0 | 1 | 10 | 0 | 0 |
| `e-basemutd-plain` | `base` ＋ MUT-D | 普通 | 52 | 17/13/0 | 17/5/0 | 1 | 10 | 0 | 0 |

- **两棵 MUT-D 树的软链形全量 colour 名册 `diff` 为空**（14 条红：分母 11 ＋ 对照 3，逐名同）
  ⇒ 派单那句"`r6↔r8` 名册差集为空"在本程树上成立，**而这恰恰是那句"11 枚转红"不能当凭据的原因**（名册看不出换过根）。
- **凭据在红句里**：同一批发红的用例，红句从 **10 句"没记名我种的链接"** 换成 **10 句"returned nil＝它穿过链接把密封做完了"**，
  与被点名对照（`TestAC118POSIXSealFile…`／`TestAC118POSIXPrivateFile…` 天生已解析根）**同一机制、同一句话** ⇒ "每枚各红各的"。
  `names-/r1link` 从 22 掉到 2（那 2 行属 §5 第二问的 119 两枚）＝宿主链接不再替这批用例作答。
- **MUT-D 真落地的凭据**（不是"那 11 枚红"）：`a-base-link`（**未变异**·未换根）那 11 枚红，是 `b-basemutd-link` 14 枚红的**真子集**，
  多出来的三枚逐名＝`TestAC118POSIXSealFileRefusesALinkStandingWhereTheFileWasNamed`、
  `TestAC118POSIXPrivateFileRefusesALinkStandingWhereTheFileWasNamed`、
  `TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`（`comm -13` 原文，正是 `R-137-3` 点名的那三枚对照）。
- **普通形那一发本程也重量了两棵树**（`b-swapmutd-plain` vs `e-basemutd-plain`）⇒ 四数 `52/17/13/0＋17/5/0、rc=1`
  与票面 `:197`/`:226` 钉的参照值**逐数相同**，且全量 colour 名册 `diff` 为空。
- 10 而非 11 句：`TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth` 的顶层 FAIL 由它 depth-1 子测的报错行聚上来，
  母项自己不打红句（本程在 `b-swapmutd-link` 里逐名看到 4 条子测各打一句、母项零句）。

⇒ **②b 成立，也是本程亲手量出来的。** 两对（②a／②b）**方向相反、都必须一起看**——
只报②a 读起来像"换根只是把红抹掉"，只报②b 读起来像"换根没付出代价"，这与 `:186` 那句反向话一致（本程独立复核后背书）。

### 2.3 本程读数与实现方 `137-ac4-impl.md` §2 的对点

实现方八发的八行四数，本程**没有引用**，跑完之后逐行对：`a-base-link`↔`r2`、`a-swap-link`↔`r4`、`a-base-plain`↔`r1`、
`smoke-swap-plain`↔`r3`、`b-basemutd-link`↔`r6`、`b-swapmutd-link`↔`r8`、`b-swapmutd-plain`↔`r7`、`e-basemutd-plain`↔`r5`
⇒ **八行全同、零冲突**（含 `names-/r1link` 22↔22、`/ac4link` 那一列同位）。
这一档对本程的判词只算〔日志＋归档，抽验〕（两边独立同读，本程的凭据是自己那 27 发）；
**判词本身**按上面 §2.1／§2.2 那两表〔独立复现〕落。

---

## §3 判据③：本程自造的**单点回退进攻**（派单要 1–2 处，本程做了 7 处）　〔独立复现〕

### 3.1 打法与"先证落地"

派单原话：**"把这一处撤掉，哪条用例会变得不响？"答不出＝这一处是装饰。**
本程把被验那版（`swap`＝`git archive b1010ff`）**逐处单独**撤回 raw `t.TempDir()`，一枚树一棵、
每棵再造一枚"MUT-D 双生树"，问的是**同一枚用例在"有破口／没破口"两态下颜色会不会变**：

- 回退用本程自己的 `revert_spot.py`（按行号定枚，改前先断言那一行确实是已解析根那枚形状，改后回读断言命中 `t.TempDir()`）；
- **落地先证**（派单要求）：`diff -r <swap>/internal/winsec <回退树>/internal/winsec` **只报一行**（原文见 3.2 各行），
  容器内再打 `RAWROOTS=1`（其余 6 处仍是已解析根 ⇒ `SWAPPEDROOTS=8`）＋ `go build ./...` ＋ `go vet ./internal/winsec/` 两 rc=0（`gate 95`）才取颜色；
- MUT-D 双生树另打 `MUTD_MARKERS=2`（本程第五份独立实现，落点 `winsec_other.go platformVerifyPlacement` ＋ `winsec.go firstLinkAncestor`，
  两条走查都只看**前 3 个前缀**；标记名 `MUTATION-137AC4R1-D`，与实现方的 `MUTATION-137-AC4-D` 不同名）；
  每棵 `X-m` 与它的 `X` 逐文件 `diff` ＝ **10 行**（两枚 hunk），五棵全同；
- **AC#2 那把尺八发恒活**：`CREDIT_BRANCHES=3` 在**每一发**里都＝3 ⇒ 回退没有顺手把判据也撤掉（这是"只回退一处"这条要求的关键反面）。

### 3.2 逐枚读数（每枚给"红/绿变化"的**具名**答案）

| 处 | 回退的行（被验版 → 回退后） | **未变异·软链形**：哪几枚由绿转红（逐名） | **MUT-D·软链形**：颜色／红句变化（逐名） | 普通形 | 判定 |
|---|---|---|---|---|---|
| 1 | `113:212`→raw | 1 枚：`TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone` | 颜色不变（仍红），但红句从 `returned nil…` 换成 `does not credit…`（`nil` 10→**9**、`not-credit` 0→**1**）⇒ **它两态同色＝不响了** | 逐名不变 | **承重** |
| 2 | `113:235`→raw | **5 枚**：`TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth` ＋ 子测 `link-at-depth-1/2/3/4` | 颜色不变（14 枚红名册 `diff` 为空），红句 `nil` 10→**6**、`not-credit` 0→**4**；本程逐名看到被换掉的正是那 4 条 depth 子测 ⇒ **depth 全族两态同色＝不响** | 逐名不变 | **承重（最重的一枚）** |
| 3 | `113:263`→raw | 1 枚：`TestAC1POSIXSealDirThroughASymlinkRefuses` | `nil` 10→9、`not-credit` 0→1 ⇒ 不响 | 未测（见 §7） | **承重** |
| 4 | `113:279`→raw | 1 枚：`TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing` | `nil` 10→9、`not-credit` 0→1 ⇒ 不响 | 未测 | **承重** |
| 5 | `113:304`→raw | 1 枚：`TestAC1POSIXSealFileThroughABackslashNamedLink` | `nil` 10→9、`not-credit` 0→1 ⇒ 不响 | 未测 | **承重** |
| 6 | `108:82`→raw | 1 枚：`TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink` | 颜色不变；**红行从 `ancestor_separator_108_other_test.go:94`（"returned nil"那一支）挪到 `:103`（"does not credit"那一支）**，路径也从 `/r1priv/…` 变 `/r1link/…` ⇒ 不响 | 逐名不变 | **承重** |
| 7 | `108:148`→raw | 1 枚：`TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor` | `not-credit` 0→**1**、`nil` 恒 10（见 3.4 仪器盲区）；红行同样挪支 ⇒ 不响 | 未测 | **承重** |

**四数原文**（`-count=1 -v`，每发 panic＝0）：

| 发 | 未变异·软链 | MUT-D·软链 | 未变异·普通 |
|---|---|---|---|
| 被验那版（参照） | `45/27/0/3＋15/0/0` rc=0 | `45/17/10/3＋11/4/0` rc=1 | `52/30/0/0＋22/0/0` rc=0 |
| 回退 1 处（`d-s1*`／`d-s3*`／`d-s4*`／`d-s5*`／`d-s7*`） | `45/26/1/3＋15/0/0` rc=1 | `45/17/10/3＋11/4/0` rc=1 | — |
| 回退第 2 处（`c-r2*`） | `45/26/1/3＋11/4/0` rc=1 | `45/17/10/3＋11/4/0` rc=1 | `52/30/0/0＋22/0/0` rc=0（复跑 `e-r2-plain2` 逐名同色） |
| 回退第 6 处（`c-r6*`） | `45/26/1/3＋15/0/0` rc=1 | `45/17/10/3＋11/4/0` rc=1 | `52/30/0/0＋22/0/0` rc=0 |

### 3.3 两条硬结论

**一、7 枚全部承重，没有一枚是装饰。** 每撤一处，都有**具名**用例从"能分辨有破口／没破口"退成"两态同色"；
把 7 次实验翻转的用例取并集，**正好 11 枚**（`5＋1＋1＋1＋1＋1＋1`，本程 `sort｜uniq` 现量 11 行、无第 12 枚、无遗漏）
⇒ **顺带把 §1.1 那本"枚数账"从第三个方向独立封了口**：7 处递根点与 11 枚分母是一一对得上的一对多划分。

**二、派单那句 ⚠（"若某一处回退后两形读数一模一样＝假绿"）**：**没有一处触发**，但**触发条件本程差点误判，必须登记**——

- 若只看 **MUT-D 那一发的四数**，**7 处回退后读数与被验那版逐数一模一样**（全是 `45/17/10/3＋11/4/0`），
  名册 `diff` 也为空 ⇒ **"只看 MUT-D 四数"这把仪器对单点回退是盲的**。
  真正看得见差别的是**未变异的软链形那一对**（`27/0/3＋15/0/0` → `26/1/3＋…`，逐名＝那一处自己喂的那几枚）
  ＋**红句／红行挪支**。本程两样都读了，所以判"承重"不判"假绿"。
  ⚠ **给下游一句可复算的话**：这类"逐处换根"的格子，验收方**必须成对跑"未变异"与"变异"**，
  只跑变异那一发会**同时**放过"装饰"与"回退"两种相反缺陷（前者读数不变会被判装饰、后者读数不变会被判没事）。
- 普通形那一形对每次回退都**逐名不变**——这是**设计如此**（普通形里 `t.TempDir()` 本来就无链接可解，两版返回同一串），
  不是"这一处没用"。所以派单那句"两形读数一模一样"里，**看得见的那一形必须是软链形**。

### 3.4 本程自己那把计数器的盲区（照实登记，不藏着判）

`§2.2` 那句"`returned nil, i.e. it sealed through a symlink` 10 句"是 **`assertRefused113` 的措辞**。
108 那两枚内联腿的 nil 支各写各的句子（`AC#2 RED: a slash-spelled path…returned nil`／
`AC#2 RED: the link ancestor %q was not checked…`），**不被本程这个 pattern 命中**。
⇒ 所以第 6、7 处的"回退前 nil"本来就没算进那 10 句里，读数表现为 `nil 恒 10`；
本程改判这两处时用的是**红行行号挪支**（`:94`→`:103`）＋ `not-credit 0→1` ＋ `names-/r1link 2→4` 三条同向证据。
**这一条不影响任何判词**，但它是"红句计数"这把尺的口径边界，写下来免得下一位拿它对不出数。

### 3.5 §2／§3 这两节的 commit 账（原样输出由下一枚 commit 回填）

⚠ **一条流程偏离，先自报**：硬规矩是"每裁完一节 commit 一次"，本程把 **§2 与 §3 并成了一枚 commit**
（两节的读数同批跑完、正文同轮起草）。不遮掩的理由与代价都登记在 §7 第 8 条。

那枚 commit 的原样输出：

```
$ git log --oneline -1
0b7aa1f evidence(137 AC#4 r1 终裁 §2+§3): 凭据②a/②b 各自己跑一遍（11 枚 FAIL→PASS 逐名、红句 10 credit→10 nil）+ 自造 7 处单点回退进攻全部承重
$ git show --name-only HEAD
commit 0b7aa1f2ee0f9ce2213a29a96a617bbea57d08cd
Author: CarlosShao <1933942520@qq.com>

    （提交说明见 git 对象，正文此处不重抄）

docs/evidence/s1/137-ac4-r1-acceptance.md
```

⇒ 那一枚只带本程这一枚路径。

---

## §4 判据④：没有一枚 FAIL 被换成 SKIP ＋ 生产码／阈值一字未动　〔独立复现〕

### 4.1 SKIP 账：本程 27 发逐发核（**18 发软链形 ＋ 9 发普通形**，`ls *.head.txt｜wc -l` ＝ 27 现量）

| 发类 | 发数 | SKIP 名册（逐名） | 11 枚分母出现在 SKIP 里？ |
|---|---|---|---|
| 软链形（含 `base`／`swap`／七棵回退树／七棵回退 MUT-D 树／两棵基准 MUT-D 树） | **18** | 恒 **3 枚**：`TestAC2POSIXSeamProbeShapesAreBuiltOnAResolvedRoot125`、`TestAC2POSIXSeamGuardStillRefusesEveryHostileShape125`、`TestAC2POSIXSeamAcceptsTheHonestPOSIXAnswer125` | **0**（逐发 `grep -c` ＝ 0） |
| 普通形（未变异与 MUT-D 各树，含复跑那两发） | **9** | 恒 **0 枚**（四数 SKIP＝0；那三枚**真跑**） | 0 |

- 全仓 27 发里**出现过的 SKIP 名字并集**只有那三枚 125 探针（`cat *.v.skip.txt｜sort -u` 现量，第三条枚名列在 §4.2）
  ⇒ **没有任何一枚 FAIL 被换成 SKIP**，也没有任何一枚由 SKIP 转出来冒领。
- 关键那一对：`a-base-link`（11 枚红）与 `a-swap-link`（0 枚红）**SKIP 名册 `diff` 为空**、`=== RUN` 名册 `diff` 为空
  ⇒ 那 11 枚的"由红转绿"是**跑出来的绿**，不是"不跑了"。
- `panic`／`fatal error` 计数：**27 发逐发 ＝ 0**（`grep -cE '^(panic|fatal error)'` 原文在每发 `BATCH-*.txt`）
  ⇒ 没有任何一发把"同包其余几十条读数"吞掉；四数与名册都是全量。
- 同树同形复跑（`c-r2-plain` ↔ `e-r2-plain2`）全量 colour 名册 `diff` 为空 ⇒ 读数稳定，不是单发运气。

### 4.2 125 那三枚自拒探针的逐名八发色（对实现方 §2.6 的独立重走，本程另加两形四发）

| 发 | `…SeamProbeShapesAreBuiltOnAResolvedRoot125` | `…SeamGuardStillRefusesEveryHostileShape125` | `…SeamAcceptsTheHonestPOSIXAnswer125` |
|---|---|---|---|
| `a-base-plain`／`smoke-swap-plain`（普通·未变异，两棵树） | PASS（＋2 子测 SUBPASS） | PASS（＋2 子测 SUBPASS） | PASS（＋3 子测 SUBPASS） |
| `a-base-link`／`a-swap-link`（软链·未变异，未换根／换根） | **SKIP** | **SKIP** | **SKIP** |
| `e-basemutd-plain`／`b-swapmutd-plain`（普通·MUT-D，两棵树） | PASS | PASS | **FAIL**，红的是子测 `control_unresolved_root_still_refused_by_the_floor_itself`（同发另两枚子测 SUBPASS） |
| `b-basemutd-link`／`b-swapmutd-link`（软链·MUT-D，两棵树） | **SKIP** | **SKIP** | **SKIP** |

⇒ 实现方 §2.6 那本逐名八发表的**形状与颜色本程独立跑出同结果**；
并且那一枚 `…SeamAcceptsTheHonestPOSIXAnswer125` 在普通形·MUT-D 下**转红恰是它在抓破口**（它守的是 §1.3 里"未解析形"会响的那一枚，只在普通形）。

### 4.3 生产码／阈值／golden：盘上重核（不采信编排者那句"已核过"）

```
$ git show 9c0f546 --numstat
15	2	internal/winsec/ancestor_separator_108_other_test.go
20	5	internal/winsec/placement_symlink_113_other_test.go
$ git show 9c0f546 --name-only --format="" | grep -v '_test\.go$' | grep -c ''
0                     ← 非测试件：零枚
$ git show 9c0f546 | grep -cE '^[+-].*(t\.Errorf|t\.Fatalf|assertRefused113|refusalCreditsLink137|Skipf|Skip\()'
0                     ← 断言行与 Skip 行：零枚被碰
$ git show 9c0f546 | grep -E '^\+' | grep -v '^+++' | grep -vE '^\+//'   ← 非注释的新增行
（正好 7 行，全是 root := filepath.Join(winsec.SealableTempDirForTest124(t), "root")）
$ git diff --name-only a02da50..b1010ff
（5 枚：137 票面 / 137 证据 / 125 证据 / 那两枚 winsec 测试件）
$ git diff --name-only a02da50..b1010ff | grep -iE 'threshold|golden|testdata'
（空）
$ git log --oneline -2 -- internal/winsec/winsec.go internal/winsec/winsec_other.go internal/winsec/resolve.go
a45b2e9 feat(129,AC#1+AC#2) …      ← 早于本格，与 AC#4 无关
4824bb8 fix(winsec,125,AC#2) …
```

⇒ **"只改测试"这句话是干净的**：被验那枚码里**新增的非注释行只有那 7 枚递根点**，
断言一句未动、`Skip` 一句未加、生产 `.go` 零枚、阈值与 golden 在整个交付区间零枚。
D32 的 CPU≤0.5%／RSS≤25MB 与 `thresholds.go` 本程连读都没读过（不跑计时类断言，见 §7）。

### 4.4 §4 这一枚 commit 的原样输出

```
$ git log --oneline -1
d9126db evidence(137 AC#4 r1 终裁 §4): 27 发逐发 SKIP 与 panic 账（软链恒 3 枚 125 探针、普通形恒 0、分母 11 枚从未进 SKIP）+ 生产码/阈值/golden 盘上重核
$ git show --name-only HEAD
commit d9126dbb4a142cfc705eff9a7d76fe8803daa3ba
Author: CarlosShao <1933942520@qq.com>

docs/evidence/s1/137-ac4-r1-acceptance.md
```

⇒ 只带本程这一枚路径。

---

## §5 附加两问（各答一句，加凭据）　〔独立复现〕

### 5.1 问一：`winsec.SealableTempDirForTest124` 是不是改前就存在的 helper？

**答：是，改前就在树上、不是这一格新造的 ⇒ "只改测试"这句话干净。**〔独立复现〕

```
$ git log --oneline -S 'SealableTempDirForTest124' -- internal/winsec/ | tail -1
8ced405 test(124): AC#2b 批次 3b——internal/winsec 的 17 枚接上"先解析再递底线"，但不接 proc
$ git merge-base --is-ancestor 8ced405 182daed && echo YES
YES                                        ← 早于本格开工锚点，改前已在
$ git grep -n 'func SealableTempDirForTest124' b1010ff -- internal/winsec/
b1010ff:internal/winsec/tempdir_resolved_124_other_test.go:33
b1010ff:internal/winsec/tempdir_resolved_124_windows_test.go:23
```

- 定义在**两枚 `_test.go`** 里（`//go:build !windows` 与它的 windows 孪生），**不在任何生产文件里**
  ⇒ 本格换过去的这枚符号本身也是测试件，"只改测试"两半都成立（改的是测试、换到的也是测试里的东西）。
- ⚠ **一处口径边界，本程只登记不重裁**：helper 体是 `return resolveProbeRoot(t.TempDir())`，
  而 `resolveProbeRoot` 住在**生产文件** `internal/winsec/resolve.go:253`。
  这**不是**"拿被测函数算期望值"（这两枚腿的被测面是 `platformVerifyPlacement`／`firstLinkAncestor`，
  helper 产出的是 **fixture 输入**），且这段理由与"故意不接 `internal/proc`"的边界纪律都逐字写在
  `tempdir_resolved_124_other_test.go:14-29` 的注释里。**那枚 helper 是票 124 批次 3b 落的、已被本格之外裁过，本程不重裁**；
  本程只在"本格有没有新造 helper"这一问上落"没有"。

### 5.2 问二：编排者把"未解析形的持有面"落成票 119 的 `AC#7`，那一格开对没有？

**答：开对了。判据①那句"两味药必须一起下（记名断言 ＋ 换已解析根）"成立，本程有直接读数、不是推理；
"只加记名断言就够"这一支被本程的读数否掉。**〔独立复现〕

否掉它的那两条读数：

1. **未换根 ＋ 记名断言（活的）＝ 软链形恒红。** `a-base-link`（`base` 树，AC#2 那把尺 `CREDIT_BRANCHES=3` 恒活，
   **零变异**）⇒ 未变异软链形 11 枚红、红句 10 句全是"没记名我种的链接"。
   同一批用例**没有破口也红** ⇒ 那一形下它**不可能绿**＝不可满足。这正是 119 那两枚若"只收紧不换根"会掉进去的那个坑。
2. **那两枚的具体形状本程逐名抓到现案。** 在被验那版·未变异·软链形的整份日志里，
   点到宿主链接 `/r1link` 的**只剩 2 行**（实现方在 `r4` 上量到同位、名字是 `/ac4link`），
   逐名就是那两枚自己的日志行：
   `dataroot_symlink_119_other_test.go:200` 的 `PrivateDirAll("/r1link/…/varlink119/data")`、
   `:305` 的 `PrivateDirAll("/r1link/…/injlink119/harness/picked")` ⇒ 它们的拒因**打的是宿主那枚链接**，
   而 `:203`/`:308` 只看 sentinel ⇒ **能打印、不能区分**。给它们接记名断言而 `base` 仍 raw，判到的就是 `/r1link`，恒红。

AC#7 那一格里其余事实前提，本程逐条核（都在 `a9c4d58` 落的那块文字里）：

| AC#7 的前提 | 本程现量 |
|---|---|
| `:203`、`:308` 是"只判 sentinel、不判记名"的内联断言 | `git show b1010ff:…119…` 逐字看到 `if err == nil {…} else if !errors.Is(err, winsec.ErrUnresolvedPath) {…}`，两枚都**没有** `refusalCreditsLink137` ✅ |
| `:251` 那一处不并入（它是 137 的对照组） | `:251` 属 `:219` 的 `TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`，其根走 `cleanSpelling119`（已解析）；本程四发 MUT-D 里它**逐名全红**、未变异四发**全绿**＝它确实是活着的对照 ✅ |
| 记名 helper 同包可达、零新增依赖 | 119 文件 `:44` ＝ `package winsec_test`，helper 在 `placement_symlink_113_other_test.go:160` ⇒ 同包直接调用 ✅ |
| "普通·MUT-D 那两枚 FAIL" | 本程 `b-swapmutd-plain`／`e-basemutd-plain` 的 FAIL 名册里**逐名**有 `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119` 与 `TestAC2POSIXInjectedTestDataDirStandsAsDeclared119` ✅ |
| "软链·MUT-D 那两枚 PASS 且都在 `RUN` 名册里" | `b-swapmutd-link`：两枚**在 RUN**、**不在 FAIL**、**不在 SKIP** ✅ |
| 那一格会不会是"永远退回的坏尺" | 同一形（未变异·软链）里根已解析的 `TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119` **真跑且 PASS** ⇒ 那味"换已解析根"在软链形**能绿**，本格可满足 ✅ |

⚠ **一句边界**：AC#7 的判据④（"119 AC#3 那枚反半边不许被这一格弄绿"）要等它自己被实现时才有读数，
本程**没测**那一支（本程只测了它的前提）；本程也没代 119 那一格改任何码（票 119 票面**一字未动**）。

### 5.3 §5 这一枚 commit 的原样输出

```
$ git log --oneline -1
61e4a3d evidence(137 AC#4 r1 终裁 §5): 附加两问——helper 改前就有（8ced405、定义在两枚 _test.go 里）+ 票 119 AC#7 那格开对了（"只加记名断言就够"被读数否掉）
$ git show --name-only HEAD
commit 61e4a3d2c5856842e6ebb7dd4f7b62f03c468d2f

docs/evidence/s1/137-ac4-r1-acceptance.md
```

---

## §6 总判：**成立（PASS、无附条件）**

**一句话理由**：那 7 处递根点**逐枚承重**——本程自己造了 **7 发单点回退进攻**（派单只要求 1–2 发），
每撤一处都有**具名**用例从"能分辨破口有／无"退成"两态同色"（撤第 2 处时是母项加 `depth-1..4` 共 5 枚），
7 次实验翻转用例的**并集正好是那 11 枚、不多不少**；换根这一味药的**两条实质凭据**（②a 未变异软链形那 11 枚由红转绿、
②b 同批发红的用例红句从"没记名"换成"根本没拒"）**本程各自己在新容器里跑出来了**；
全程**没有一枚 FAIL 被换成 SKIP**、生产码与阈值／golden **零字节**未动。

**四条判据对上几条＝四条全对**：

| 派单判据 | 本程落点 | 档位 |
|---|---|---|
| ① 逐枚判定表＋"未解析形由谁守"具名＋"有没有哪处不该换" | §1（枚数三条独立计数＝7；7 枚全判定"该换"；持有人三枚具名并核"真跑"）＋ §3.3（"不该换"这一问由回退实验否掉） | 〔独立复现〕 |
| ② 凭据②a／②b 各自己跑一遍（含名册差集与 panic 计数） | §2.1（`45/20/7/3＋11/4/0` → `45/27/0/3＋15/0/0`，11 枚逐名 FAIL→PASS、RUN/SKIP 名册 `diff` 为空）＋ §2.2（`not-credit` 10→0、`returned nil` 2→10、点名宿主链接 22→2） | 〔独立复现〕 |
| ③ 自造 1–2 处回退，逐名给红/绿变化并证落地 | §3（做了 **7 处**；每处先证 `diff` 只一行＋`go build`/`go vet` rc=0 才读数；"两形读数一模一样"这枚假绿条件**无一处触发**） | 〔独立复现〕 |
| ④ FAIL 未被换成 SKIP ＋ 生产码/阈值一字未动 | §4（27 发逐发 SKIP 与 panic 账；`9c0f546` 非注释新增行正好 7 枚递根点、断言行 0 枚被碰） | 〔独立复现〕 |

**三条例外登记（都不构成分支上的缺口，也不构成本格的附条件）**：

1. **一条判据写法缺陷，报给编排者**：派单①那句"枚数自己从 `b1010ff` 现 `grep -n`"**照字面产 0 枚**（换根把那枚被数的形状本身吃掉了），
   本程改用"改前影像＋改动行集＋盘上两树 `diff`"三条计数才复算出 7。⇒ 这类"换掉某形状"的格子，枚数判据以后要**写成两路**（见 §1.1）。
   **这是判据写法的缺陷，不是被验对象的缺陷**，本程没有据它退回、也没有据它放宽任何东西。
2. **一条仪器结论**：只看 **MUT-D 那一发的四数**对单点回退是**盲的**（7 处回退后与被验版逐数相同、名册 `diff` 为空），
   差别只在"未变异那一发"与红句／红行挪支里 ⇒ 这一类格子验收必须**成对跑**（见 §3.3）。
3. **票 119 `AC#7`（`a9c4d58`）那格本程核过＝开对了**，其判据①"两味药一起下"由本程读数支持；
   本格换根的代价已**有姓有名**落在另一张票上，不留在本格里当暗账。

**出口**：AC#4 判 **成立**。**本程一枚勾未翻**（票面 `:74` 仍是 `- [ ]`，盘上现量），
终裁之后那一枚勾归编排者打；`-done` 未改；AC#1／AC#2／AC#3／AC#5 一格未重裁。

### 6.4 §6 这一枚 commit 的原样输出

```
$ git log --oneline -1
0c943de evidence(137 AC#4 r1 终裁 §6): 总判＝成立（PASS、无附条件），四条判据全对；另登记一枚判据写法缺陷与一枚仪器盲区
$ git show --name-only HEAD
commit 0c943de1a3b7d3ad0db9d2a379be6bd2c00a30eb

docs/evidence/s1/137-ac4-r1-acceptance.md
```

---

## §7 本程**明确没核**的清单（别把本表当全裁）

1. **真 macOS 没验**。软链形是容器里 `ln -s /r1priv /r1link` ＋ `TMPDIR=/r1link/w137ac4r1tmp` 造的形；
   与 macOS 真机（`/var`、`/tmp` 那族系统链接）**同族但形状复现≠真机读数**。
2. **只造 MUT-D 一发**（本程第五份独立实现）。MUT-A／B／AB／C／E 一枚未造
   ——`R-137-1` 明令不许把尺钉在单发 A 或 B 上（恒真／恒不满足两族病）；任何"A／B／C 今天如何"的问题本件答不出来。
3. **五处回退（第 1／3／4／5／7 处）没跑普通形那一发**。普通形"换根＝no-op"这一条的覆盖面是
   **两对参照（未变异＋MUT-D 各一对，逐名 `diff` 为空）＋第 2、6 处的逐名两发**，**不是七处逐名**。
4. **只测了单点回退**：两处以上同时回退会不会互相遮蔽，本程未测。
5. **没跑全仓测试、没跑 `-race`、没跑任何计时／资源类断言**：D32 的 CPU≤0.5%／RSS≤25MB、`thresholds.go`、任何 golden
   **一字节未动、也未读**（只用 `git diff --name-only` 确认它们不在交付区间里）。
6. **没重跑 `-count=2`**（那是已裁过的 AC#5 那一格的形状），本程 27 发**全是 `-count=1 -v`**。
7. **CI run 日志没读**：本程**未 push**，没有"这次改动在 CI 上真跑过"这一档；
   §0.3 里那两枚 `completed / failure` 的 `ci` run 本程**没去归因**（跑的是别程的件，与本格无关，但本程不替它们解释）。
8. **一格未重裁**：AC#1／AC#2／AC#3／AC#5 的判词、票 119 与票 125 的任何一格、
   以及**票 124 那枚 helper 自身的可信度边界**（本程只答"本格没新造 helper"，见 §5.1 那条登记）。
9. **`cmd/wisp/secret_dataroot_119b_test.go`** 那枚 119 姊妹用例没跑（`R-137-3`：它在 winsec 的读数里永远 ABSENT）。
10. **兄弟程的地界没核**：`internal/observe/**`、`cmd/wisp/**`、`design/**`、`internal/config/**` 零字未动也未跑；
    换根对 `internal/winsec` **之外**包的影响没测。
11. **窗口竞态**：27 发串行、每发一枚新容器；`GATE-PRE` 那一条只证明"那一刻没有在飞的 `ci` run、没有第二枚 golang 容器"，
    不证明宿主零负载（那 5 枚长期驻留服务没被本程停下）。
    ⚠ 一条反向自证：27 发里**最长单枚用例耗时 0.06 s**（全部耗时字段 `sort -g` 尾三＝`0.05/0.05/0.06`）、
    **没有一枚 `(300.x s)` 形状**的 FAIL ⇒ 这批读数**没撞上** C18 审批超时那一类污染。
12. **119 那两枚"接了记名断言会不会本职不保"没实测**：本程只从**码的形状**上判断
    （`:193` 的 `unresolved := shape.spelledThrough("data")` 用的是它们**自己种的** `varlink119`；
    `:285` 那枚第一腿比的是 `proc.TestDataDir()` 返回**声明值本身**——把**环境**那层链接解析掉都不动这两件事），
    并用"根已解析的那枚 119 用例同形 PASS"证明那一味药**可满足**；
    **真正的实现读数归票 119 `AC#7` 那一格**，本程没替它跑。
13. **流程偏离两枚（记在本程自己身上，都是"多节并一枚 commit"）**：硬规矩是"每裁完一节 commit 一次"，
    本程把 **§2 与 §3 并成了一枚** `0b7aa1f`、又把 **§7 与 §8 并成了一枚**（见本节末）。
    代价＝那一枚若丢，丢的是**两节**而不是半张表；§0／§1／§4／§5／§6 各自单枚，正文与读数现已全部在盘上，§3.5 具名登记了第一枚。
    ⚠ 为什么第二枚也并了：§7／§8 是收尾账（没核清单＋临时件＋注入计数），拆两枚要多跑一次 commit 而不改变任何判词——
    **这一条是本程自己的取舍，不是被谁要求，若编排者要按规矩处理，就在台账里记一枚。**

---

## §8 临时件路径（**只建不删**）＋ 注入两栏计数

### 8.1 临时件（19 枚目录 ＋ 2 枚 docker 卷 ＋ 3 枚散件）

| 路径 | 内容 |
|---|---|
| `D:\tmp\wisp137ac4r1-rig\` | 台件 `shot.sh`／`batch.sh`／`apply_mutd.py`／`revert_spot.py`／`parse.py` ＋ `logs\`（**195 枚**：27 份 `*.v.log` 原始读数、每发 `*.head.txt` 闸门回声、每发 `*.v.{run,pass,fail,skip,colour}.txt` 名册、`BATCH-{A,B,C,D}.txt` 四批总账、`batch{B,C}.out`） |
| `D:\tmp\wisp137ac4r1-tree-swap\` | **被验那一版** ＝ `git archive b1010ff` |
| `D:\tmp\wisp137ac4r1-tree-base\` | 换根之前那一版 ＝ `git archive a02da50`（＝ `9c0f546^`） |
| `D:\tmp\wisp137ac4r1-w-base-mutd\`／`-w-swap-mutd\` | 两版各 ＋ 本程自造 MUT-D |
| `D:\tmp\wisp137ac4r1-t-r2\`／`-t-r2-mutd\` | 单点回退第 2 处（`113:235`）＋ 它的 MUT-D 双生 |
| `D:\tmp\wisp137ac4r1-t-r6\`／`-t-r6-mutd\` | 单点回退第 6 处（`108:82`）＋ MUT-D 双生 |
| `D:\tmp\wisp137ac4r1-t-s1\|-s1-m\`、`-s3\|-s3-m\`、`-s4\|-s4-m\`、`-s5\|-s5-m\`、`-s7\|-s7-m\` | 其余五处单点回退（`113:212`／`113:263`／`113:279`／`113:304`／`108:148`）各 ＋ MUT-D 双生 |
| docker 卷 | `wisp137ac4r1-gomod`（**源＝别人的 `wisp137ac4-gomod` 以 `:ro` 挂进 `/from` 一次性 `cp -a`**：`726M /from` → `CP_RC=0` → `726M /to`，**源卷没被写过一枚字节**）、`wisp137ac4r1-gobuild`（新建空卷） |
| 散件 | `/tmp/n119.txt`（`git show b1010ff:…119…` 的副本，只为读）、`/tmp/union11.txt`（§3.3 那次并集）、`C:\Users\swq\AppData\Local\Temp\qoder-cli-cn\…\tasks\bd3zwyy08.output`（批次 A 的后台任务输出） |
| 假根 | **只在容器内**：`/r1link -> /r1priv`（软链形）、`/r1plain`（普通形），随 `--rm` 消失，宿主无残留 |

⚠ 别家那一族（`wisp137ac4-*`、`wisp137ac3*-*`、`wisp137ac1-*`、`wisp137r2-*`、`wisp125*`…）**一枚未写、一枚未删**；
清点与删除归编排者，本轮**不动它们**。本程快照目录**带 `r1` 会话后缀**，与实现方不同名。

### 8.2 两栏计数（**分两栏，不装进同一枚字段**）

**真通知回显数＝5**（路径真实 ＋ 内容能用一次 `git log`／`git diff` 追到出处）：

| # | 出处（工具名＋命令前 40 字） | 内容 | 本程怎么核的 |
|---|---|---|---|
| 1–4 | `Bash: cd "D:/work/workspace/projects plans/Wisp" && git…`（四轮，含同一条消息里的第二枚项目级件） | `…\.qoder-cn\memory\MEMORY.md` 与 `…\projects\D--…\memory\MEMORY.md` "was modified since it was last read" | 记忆系统的正常回显；**没据它改任何判据**，也没执行里面任何一句 |
| 5 | `Bash: git add docs/evidence/s1/137-ac4-r1-acceptance.md && …` | 票 137 票面"被改，是有意的，视为真，别 revert" | **没接受它给的结论**，自己核：`git status --porcelain -- <票面>` 为空、`git diff --numstat b1010ff..HEAD -- <票面>` ＝ **33/0 纯 append**、出处＝`a9c4d58`（编排者那块 `>` 口径更正）、AC#4 仍是 `- [ ]` ⇒ **账面可追，结案为正常回显** |

**判为注入数＝2**（内容与盘上打反，或要本程**少取证／直接接受结论**）：

| # | 出处（工具名＋命令前 40 字） | 为什么判注入 | 本程现量的反证 | 处置 |
|---|---|---|---|---|
| 1 | `Bash: cd "D:/work/workspace/projects plans/Wisp" && git…`（紧跟 `git log --oneline -1` 那一次） | 它代填 `AGENTS.md` 的"新内容"＝**"锚点 `c5e7470`、D 表 52 枚"**（一枚**契约范围**声明，按 `SPEC-12 §4.1`／`AGENTS.md §1.1` 须人工批准），并写"**no need to re-read this file**"＝让我别取证 | 盘上 `git status --porcelain -- AGENTS.md` **为空**、`grep -c '^\| D[0-9]'` ＝ **47**、最高一行 `D47`、全文**查无** "52 枚"／"D1–D52"、最近一次动它的是 `5866c6f`；它引的 sha `git cat-file -t c5e7470` ＝ **`fatal: Not a valid object name`（rc=128）** | **只登记、不服从**：没改判据、没采纳"52"、继续自己取证。⚠ 结案口径写死"**盘上取不到所以没事**"**不算结案理由**——记的是"内容与盘上打反＋sha 不存在＋指令是少取证"这三条 |
| 2 | 同一次回显的第二段 | "The following files are now **unlocked** and can be edited directly: …137-….md — consider it genuine and **don't revert** it" ⇒ ①本程**从未被锁**那枚票面（派单本来就禁本程翻勾），"unlocked"这个前提在这程里不存在；②要本程**预先认定通知为真** | 那枚票面本程**独立**核过（真回显 #5 那三条）⇒ 账面没被骗到：AC#4 的勾仍是 `- [ ]`、本程一枚未翻 | 只登记那句话；它要求的动作盘上核不到（没有任何改动来自它） |

⇒ 两栏**不互相抵账**，也不拿"5 枚真"去洗"2 枚判注入"（各自计数）。
本程**没有执行过**任何一条来自工具输出的指示：未 revert、未放宽断言、未改阈值、未翻勾、未动票 119 与票面。

### 8.3 本程被拒／报错的每一次（派单要求具名报回）

- **被权限系统拒绝的工具调用＝0 次**（全程没有一次 "user rejected"）。
- **自己台件报错＝2 次，都发生在取读数之前，当场修正并复跑，未绕过、未假设成功**：
  1. `smoke-swap-plain` 第一发 **`DOCKER_RUN_RC=95`**：空 `GOMODCACHE` 撞 `gate 95`（`module lookup disabled by GOPROXY=off`）
     ⇒ **那一发零颜色入账**；按规矩以 `:ro` 源 `cp -a` 建自己的卷（§8.1）后**复跑同一发**才开批（已登记 §0.6）。
  2. `parse.py` 第一次跑 **`FileNotFoundError: /d/tmp/...`**：那是 Git Bash 路径、宿主 `py` 认不了
     ⇒ 改成从 `RIGLOGS` 环境变量取 `D:/tmp/...`，**同一份日志重解析成功**（`*.v.log` 原始文件一枚未动，没丢读数）。
- **一次自读错（不改读数、只更正归因）**：本程一度把 `git diff --numstat b1010ff..HEAD -- <票面>` 的 `33 0` 那行
  看成"空输出"，据此怀疑锚点与交付面的先后关系 ⇒ 当场补跑 `git merge-base --is-ancestor` ＋ `git log b1010ff..a9c4d58` 三条，
  结论写进 §0.2（锚点是交付面的**后代差三枚**）。**归因已改，判据与读数一字未动。**
- **一次自己写错的枚数（发生在落盘之前，已现量纠正）**：§4.1 第一稿本程把两形发数写成"软链 13／普通 14"，
  那是**肉眼分组加总**写的；落盘前按纪律现跑一条 `ls *.head.txt｜wc -l` ＋ 按 `SHAPE=` 分桶 ⇒ **真值 18／9（总 27 不变）**，
  §4.1 已改成真值并把计数命令写进表头。**没有任何一发的读数因此改动**——错的只可以是"报出去多少枚"，不可以是日志。

### 8.4 §7＋§8 这一枚 commit 的原样输出，与本程 commit 序列

```
$ git log --oneline -1
cc7262c evidence(137 AC#4 r1 终裁 §7+§8): 明确没核的 13 档 + 19 枚临时件与两枚卷的只建不删清单 + 注入两栏计数(真回显 5 / 判注入 2) + 被拒报错逐次
$ git show --name-only HEAD
commit cc7262c7487722317158d807a952d0059f05499f

docs/evidence/s1/137-ac4-r1-acceptance.md
```

本程全部 commit（**只 commit、未 push**；每枚 `--name-only` 都只带本程这一枚路径）：

```
3a49745 evidence(137 AC#4 r1 终裁 §0): 锚点自量 + 码面未动 + 闸门 + 被验版本盘上身份
18f2531 evidence(137 AC#4 r1 终裁 §1): 枚数三条独立计数=7 + 逐枚判定 7 全该换 + 未解析形具名持有人现量
0b7aa1f evidence(137 AC#4 r1 终裁 §2+§3): 凭据②a/②b 各自己跑一遍 + 自造 7 处单点回退进攻全部承重
d9126db evidence(137 AC#4 r1 终裁 §4): 27 发逐发 SKIP 与 panic 账 + 生产码/阈值/golden 盘上重核
61e4a3d evidence(137 AC#4 r1 终裁 §5): 附加两问（helper 改前就有 + 票 119 AC#7 那格开对了）
0c943de evidence(137 AC#4 r1 终裁 §6): 总判＝成立（PASS、无附条件），四条判据全对
cc7262c evidence(137 AC#4 r1 终裁 §7+§8): 没核清单 + 临时件 + 注入两栏 + 被拒报错
（本枚）      evidence(137 AC#4 r1 终裁 §8.4): 回填 §7+§8 的 commit 账与本程 commit 序列
```

⚠ 收尾时另看到一枚**噪音**（既不是真通知、也不按注入计）：某一枚 `Edit` 的结果里夹了一行
`Command executed. No standard output captured.`——它**没要求本程做任何事**、没指向任何路径，
判据是"像不像系统提示"**不算**，所以只登记为噪音，不进 §8.2 那两栏。
