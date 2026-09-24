# 137 AC#4 —— 那 7 处"按地界有意未换"的 raw `t.TempDir()`：逐枚判定＋换根后的 MUT-D 读数

**格子**：票 `.scratch/wisp/issues/137-winsecs-12-posix-rejection-legs-go-green-in-the-symlink-shape-for-the-hosts-own-varlink-not-their-own-planted-link.md`
的 **AC#4 一格**（票面 `:74-76`）。AC#1／AC#2／AC#3／AC#5 **一格未重裁、一字未动**。
**agent**＝`worker-ticket137-ac4`（实现者）。**本件不是裁决表**：AC#4 的勾**没翻**，终裁归非实现者。
**开工锚点**＝`182daed377feeff20f73ea47675c121df5607ca5`（＝当时的 `dev` HEAD）。

三档口径（本件每格都标）：〔独立复现〕＝本程自己在树上跑出来的；〔日志＋归档，抽验〕＝引别程读数但
本程重测过它所依赖的前提；〔仅自述，不背书〕＝只有别人说过、我没复算也不据它下判的。

---

## §0 锚点自量 ＋ 争用闸门读数

### 0.1 锚点（不读脏工作树冒充被验版本）

```
$ git rev-parse HEAD
182daed377feeff20f73ea47675c121df5607ca5
$ git rev-parse --abbrev-ref HEAD
dev
$ git cat-file -t 182daed377feeff20f73ea47675c121df5607ca5
commit
```

工作树**不干净**且**与我的活无关**：`design/**` 16 枚 tracked 文件被 owner 那侧挪成未跟踪的
`design/old/`（另有空的 `design/doubao/`）。本程**一枚未还原、一枚未提交、一枚未删、没过问**；
`git add` 只用显式路径，每枚 commit 前都跑过 `git diff --cached --name-only`（原文在 §0.5、§6）。

### 0.2 争用闸门（本机 self-hosted runner 会抢 CPU；本程不 push）

| 时刻（本机 +08，逐行现量） | `gh run list --branch dev --limit 3` 的状态 | 本程那时刻在做什么 |
|---|---|---|
| `15:09:04` | **`in_progress`** run `35967768017`（`ci`/`dev`/push，已跑 3m24s；上一发同形状跑了 9m3s） | **闸门关**：只读票面与证据、只读源码、建台件目录 ⇒ **零发读数** |
| `15:16:26` | `35967768017` 已 `completed`（failure，8m38s），最近三枚无 `in_progress` | 建快照／码 commit（非读数） |
| `15:19:24` | 同上，无 `in_progress` 的 `ci` run | 落 MUT-D、种模块缓存卷（非读数） |
| `15:20:56` | 同上，最近两枚全 `completed` | **闸门开** ⇒ 读数批 15:21:04 起（`run.sh` 每发容器内另打 `DATE_UTC`） |

- `docker ps` 同刻现量：`union-api-proxy:local`／`clipsync-clipsync`／`minio/minio`／`postgres:15-alpine`
  等**长期驻留**服务（Up 4–14 天），**无第二枚 golang 构建容器** ⇒ 记录在案：宿主进程名门看不见容器负载，
  这一列只是"有没有别人在同类容器里跑"的反扫，不主张"零负载"。
- ⚠ 全程**没有任何一处**用两个时间戳相减算时长；每行时刻只对它自己那一行成立。
- ⚠ **一次闸门自拒（登记，不抹）**：第一遍批跑第 1 发就 `RUN_SH_RC=96` —— 我的两道判据串写错：
  `grep -c "refusalCreditsLink137(err"` 连函数定义行一起数（得 4 不是 3），且 `SWAPPEDROOTS` 的地基值
  在**未换根**那两枚文件里本来就是 2（`placement_symlink_113_other_test.go:388`、
  `ancestor_separator_108_other_test.go:104` 两处早在 `182daed` 上就已走已解析根），我按 0 钉 ⇒ 判据错、**不是树错**。
  **那一遍零发颜色入账**（`run.sh` 在取颜色之前 `exit 96`），四棵树的计数在宿主机逐枚重量过
  （base `RAW=7 SWAP=2 CRIT=3 MUT=0`／swap `0/9/3/0`／两棵 mutd `7/2/3/2` 与 `0/9/3/2`）后才重跑。

### 0.3 仪器与快照（全部仓外，`/d/tmp/...`；Git Bash 的 `-v C:\` 假绿坑不走）

- `docker version --format '{{.Server.Version}}'` ⇒ **`29.6.2`**，`rc=0`；`docker images` 现量
  `golang:1.27 1.31GB`（**本机已有，未 `docker pull`**）、容器内 `go version go1.27.1 linux/amd64`。
- 快照（`git archive <sha> | tar -x -C <dir>`）：
  - `wisp137ac4-tree-base` ← `182daed`（开工锚，**未换根**：那 7 处仍是 raw `t.TempDir()`）
  - `wisp137ac4-tree-swap` ← `9c0f546`（**本程换根后**那一版，见 §0.5）
  - `wisp137ac4-work-base-mutd` ＝ base ＋ MUT-D
  - `wisp137ac4-work-swap-mutd` ＝ swap ＋ MUT-D
- 挂载非空硬证（每发容器里，`gate 97`）：`wc -c /src/go.mod` **＝883** ＋ `md5sum` ＝
  `f6ef661732b1851e5c3db348113cb605` ＋ `ls -l /src/go.mod` 原文进每发 `head.txt`。
- 形状硬断言（`gate 98`，我自己的路径名）：软链形 `ln -s /ac4priv /ac4link` ＋ `[ -L ]` ＋ `readlink` 逐字核 ＋
  `TMPDIR=/ac4link/w137ac4tmp`；普通形**断言 `/ac4link` 根本不许存在** ＋ `/ac4plain` 不是软链 ＋
  `TMPDIR=/ac4plain/w137ac4tmp`。
- 落地先证再取颜色（`gate 95`）：`go build ./...` ＋ `go vet ./internal/winsec/` 任一非 0 ⇒ 不取颜色。
- 树身份闸（`gate 96`）：每发容器里逐枚数 `MUTATION-137-AC4-D` 标记枚数（`mutd`=2／`unmut`=0）、
  两枚文件里 `root := filepath.Join(t.TempDir(), "root")` 与已解析根的枚数、
  **活着的** `[!]refusalCreditsLink137(` 判据分支枚数＝3（AC#2 那把尺在八发里必须一直活着）。
- 模块缓存：`GOPROXY=off` ＋ `GOMODCACHE=/gomod`。卷 `wisp137ac4-gomod` 由**别人的源卷只读拷出**
  （`-v wisp137ac3r2-gomod:/from:ro` → `cp -a`，打印 `726M /from` → `CP_RC=0` → `726M /to`；
  源卷没被写过一枚字节），`wisp137ac4-gobuild` 新建空卷。
- MUT-D ＝ 票 137 AC#1 §4 那一形（两条 POSIX 走查都只看**前 3 个前缀**），本程**自己第四份独立实现**、
  自带标记 `MUTATION-137-AC4-D`，落点：`winsec_other.go:156-162`（`platformVerifyPlacement`）、
  `winsec.go:284-288`（`firstLinkAncestor`）。⚠ **只在仓外快照上**：仓库里 `internal/winsec` 的
  **生产码零改动**（§4 有 `git show --name-only` 与 `git diff 182daed..HEAD --stat` 双证）。

### 0.4 读数的取法与口径

- 四数一律 `-count=1 -v`（`-count=N` 会把数翻倍，本程没有 N>1 的分母读数），格式
  `RUN / 顶 PASS/FAIL/SKIP ＋ 子 PASS/FAIL/SKIP`，由 `parse.py` 从 `-v` 日志**程序化生成**
  （`<发>.v.{run,pass,fail,skip,colour}.txt`），**没有一枚名是手抄**；每发另算
  `grep -cE '^(panic|fatal error)'`（一条用例 panic 会吞掉同包其余几十条，包级 rc 看不出来）。
- 分母 11 枚 ＝ 那 7 处递根点覆盖的 **7 枚顶层 ＋ 4 枚子测试**（第 2 处一处覆盖 1 顶层＋4 子测），
  口径来源＝`docs/evidence/s1/137-ac1-teeth-or-not.md:42/54` 的逐处 `grep -n` 表；本程在开工锚点上
  **重新 `grep -n` 过**（§1.1），没有抄那一版行号。

### 0.5 本程第一枚 commit（生产码零改动；只有那两枚测试件）

```
$ git add internal/winsec/placement_symlink_113_other_test.go internal/winsec/ancestor_separator_108_other_test.go
$ git diff --cached --name-only
internal/winsec/ancestor_separator_108_other_test.go
internal/winsec/placement_symlink_113_other_test.go
$ git commit -q -F - -- internal/winsec/placement_symlink_113_other_test.go internal/winsec/ancestor_separator_108_other_test.go <<'MSGEOF' … MSGEOF
```

`git log --oneline -1`（提交时刻现量 `15:16` 前）：

```
9c0f546 test(winsec/137 AC#4): 那 7 处递根点换成已解析根，并写明未解析形归谁守
```

`git show --name-only HEAD`：

```
commit 9c0f546…（作者 worker-ticket137-ac4）

internal/winsec/ancestor_separator_108_other_test.go
internal/winsec/placement_symlink_113_other_test.go
```

---

## §1 判据①：那 7 处递根点的逐枚判定表

### 1.1 先重数一遍枚数（开工锚点 `182daed` 现量，不抄任何前手行号）

```
$ git grep -n 'root := filepath.Join(t.TempDir(), "root")' 182daed -- \
    internal/winsec/placement_symlink_113_other_test.go internal/winsec/ancestor_separator_108_other_test.go
internal/winsec/placement_symlink_113_other_test.go:197:	root := filepath.Join(t.TempDir(), "root")
internal/winsec/placement_symlink_113_other_test.go:220:			root := filepath.Join(t.TempDir(), "root")
internal/winsec/placement_symlink_113_other_test.go:248:	root := filepath.Join(t.TempDir(), "root")
internal/winsec/placement_symlink_113_other_test.go:264:	root := filepath.Join(t.TempDir(), "root")
internal/winsec/placement_symlink_113_other_test.go:289:	root := filepath.Join(t.TempDir(), "root")
internal/winsec/ancestor_separator_108_other_test.go:69:	root := filepath.Join(t.TempDir(), "root")
internal/winsec/ancestor_separator_108_other_test.go:135:	root := filepath.Join(t.TempDir(), "root")
```

⇒ **113 五处 ＋ 108 两处 ＝ 7 处**，与派单在 `182daed` 量的形状相同、与 AC#1/终裁在 `1d38206`/`4a0d7a4`
量的 `146/169/197/213/238`＋`69/127` **同处不同号**（AC#2 往同一枚文件插了行 ⇒ 裸行号会被自己的修法挪走，
这一条票面 `:68-71` 已钉死）。同两枚文件里**另有 2 处早已走已解析根**（`182daed` 上的 `113:388`／`108:104`，
本程加注释后在 `9c0f546` 上是 `:403`／`:117`；非本格分母），`dataroot_symlink_119_other_test.go` 的 5 处与
`seam_probe_root_125_other_test.go` 的 2 处**不在这 7 处的地界里**（前者见 §3，后者是 125 自拒探针的 fixture）。

### 1.2 逐枚判定（判据①：换 or 不换 ＋ 一句理由 ＋ 未解析形此后由谁守）

| # | 递根点（换根前 → 换根后行号） | 它喂的那枚用例／子测 | 判定 | 一句理由（为什么它**不**离不开未解析的根／或为什么必须换） |
|---|---|---|---|---|
| 1 | `placement_symlink_113_other_test.go:197` → `:212` | `TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone` | **换** | 它的断言对象是**自己种在 `root` 里的 `link`**（`linkTo113` → `SealFile(root/link/keep-me.txt)`）；根一未解析，走查在第 1 个组件就拒掉整条拼写，`refusalCreditsLink137` 判到的是宿主的链接 ⇒ 这一枚**今天响的不是它的尺**（§2 发 `r2` 读数） |
| 2 | `:220` → `:235` | `TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth` ＋ 其 4 枚子测 `link-at-depth-1..4` | **换** | 这枚用例存在的理由是"链接在根下第 2/3/4 层时走查还得走到"；未解析根让那 4 层**全在拒点之下**，depth-2/3/4 三枚子测**根本没被检查过** ⇒ 留着未解析根＝把这枚 mutation-bait 直接注销 |
| 3 | `:248` → `:263` | `TestAC1POSIXSealDirThroughASymlinkRefuses` | **换** | 同 1，入口是 `SealDir(root/link)`；`platformVerifyPlacement` 先撞宿主链接，`SealDir` 那条腿不参与读数 |
| 4 | `:264` → `:279` | `TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing` | **换** | 这枚要钉的是"字节落盘**之前** placement 腿就跑完了"；根未解析时"没写字节"是宿主链接那份拒给的，**不是**这条顺序给的 ⇒ 断言与它声称的因果脱钩 |
| 5 | `:289` → `:304` | `TestAC1POSIXSealFileThroughABackslashNamedLink` | **换** | 被测的是**含反斜杠的那一枚链接名**会不会被切成两截；未解析根让走查在**它之前**就返回 ⇒ 这一形（A74(3) 的 fail-open 钉）恰是"永远走不到"的那一枚 |
| 6 | `ancestor_separator_108_other_test.go:69` → `:82` | `TestAC2POSIXAncestorGuardRefusesASpellingThroughASymlink` | **换** | 走的是 `RemoveUnlinked`→`firstLinkAncestor`（票面 `R-137-2` 已定：这一族的断言是**内联**的、不走 helper）；纯斜杠拼写的祖先腿在未解析根下永远停在宿主链接，AC#2 给它新加的那把尺在软链形**只能响、不能绿** |
| 7 | `:135` → `:148` | `TestAC2POSIXABackslashInALinkNameIsStillALinkAncestor` | **换** | 同 5 ＋ 同 6：反斜杠名＋未解析根＝两重"走不到"；它测的是 `Lstat` 用的必须是**原名字面**，未解析根把那次 `Lstat` 都省了 |

⇒ **7 处全部换。没有一处需要"未解析的根"来完成它自己的断言**：每一枚的拒因都必须落在
**它自己种下去的那枚链接**上（AC#2 之后这就是它们的判据本体），而未解析根恰恰是**把那个拒因换掉**的东西。
这一条不是账面整齐的说法，是 §2 那 8 发读数的读法：`r2`（未换根·软链形·无变异）响 11 枚、
`r4`（换根·软链形·无变异）响 **0** 枚，而 `r6`／`r8`（未换根／换根两棵树上各打一发 MUT-D·软链形）都是 11 枚响
——**换根前"能响"来自宿主链接、换根后"能响"来自被测走查本身**（红句换轨，§2.3），
这两件事必须分开记账（判据③），而"名册没变"那一半的归因要单独更正一次（§2.4）。

**换掉之后，"未解析那一形"由谁守（具名用例，不指票号）**：

1. `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119`
   （`internal/winsec/dataroot_symlink_119_other_test.go:193`）——把 `PrivateDirAll(<base>/varlink119/data)`
   的未解析根本身当被测对象，断言拒＋`ErrUnresolvedPath`＋两枚树都不许被创建；
2. `TestAC2POSIXInjectedTestDataDirStandsAsDeclared119`（同文件 `:285`）——第一腿声明一枚**穿过自己种的链接**
   的根，断言"照声明返回"且被地板拒、不创建任何东西；
3. `TestAC2POSIXSeamAcceptsTheHonestPOSIXAnswer125` 的子测
   `control_unresolved_root_still_refused_by_the_floor_itself`
   （`internal/winsec/seam_probe_root_125_other_test.go:277`）——直接问地板：`builtinVerifier{}.Resolve(<link-spelled>)`
   必须报错；**外加**同文件另两枚自拒探针（`TestAC2POSIXSeamProbeShapesAreBuiltOnAResolvedRoot125`、
   `TestAC2POSIXSeamGuardStillRefusesEveryHostileShape125`）与 `internal/config/c26_seam_posix_125_test.go:118-129`
   的同号族——这几枚**在软链基根下按设计 `t.Skipf`**（`plantTemps125` 看到基根穿链接回自己就自拒），
   所以它们守的是**普通形**那一半。⚠ 由此得出一条本格的**代价读数**，写在 §3。

---

## §2 判据②③：换根前后的两形读数与名册差集（八发，全部本程亲手）

### 2.1 八数一览（`go test -count=1 -v ./internal/winsec/`，一发一枚新容器）

| 发 | 树 | 生产码 | 形 | RUN | 顶 P/F/S | 子 P/F/S | 包级 rc | `panic\|^fatal error` |
|---|---|---|---|---|---|---|---|---|
| `r1-unswap-plain` | `wisp137ac4-tree-base`（`182daed`） | 未变异 | 普通 | 52 | 30/0/0 | 22/0/0 | **0** | 0 |
| `r2-unswap-link` | 同上 | 未变异 | 软链 | 45 | 20/**7**/3 | 11/**4**/0 | **1** | 0 |
| `r3-swap-plain` | `wisp137ac4-tree-swap`（`9c0f546`） | 未变异 | 普通 | 52 | 30/0/0 | 22/0/0 | **0** | 0 |
| `r4-swap-link` | 同上 | 未变异 | 软链 | 45 | **27**/0/3 | **15**/0/0 | **0** | 0 |
| `r5-d-raw-plain` | `wisp137ac4-work-base-mutd` | **MUT-D** | 普通 | 52 | 17/13/0 | 17/5/0 | 1 | 0 |
| `r6-d-raw-link` | 同上 | **MUT-D** | 软链 | 45 | 17/**10**/3 | 11/**4**/0 | 1 | 0 |
| `r7-d-swap-plain` | `wisp137ac4-work-swap-mutd` | **MUT-D** | 普通 | 52 | 17/13/0 | 17/5/0 | 1 | 0 |
| `r8-d-swap-link` | 同上 | **MUT-D** | 软链 | 45 | 17/**10**/3 | 11/**4**/0 | 1 | 0 |

- 落地先证：八发 `BUILD_RC=0` ＋ `VET_RC=0` 全中（`gate 95`），`PANIC_LINES=0` 八份逐份 0；
  `COLOUR_LINES` 每发＝该形 `=== RUN` 枚数（普通 52＝30＋22、软链 45＝30＋15）⇒ 读数是全量、没有被 panic 吞掉的行。
- 与别程参照对点〔日志＋归档，抽验；本程自己重跑，没抄它的日志〕：`r1`＝AC#1 基线普通形
  `52/30/0/0＋22/0`；`r2`＝AC#3 r1 件 §0.7 的 `ac3-tight-link` `45/20/7/3＋11/4/0`；
  `r5`/`r7`＝同件的 `ac3-d-tight-plain` `52/17/13/0＋17/5/0`；`r6`/`r8`＝同件的 `ac3-d-tight-link`。
  **六行全对上，零冲突。**

### 2.2 名册差集（`comm` 双向，四档：RUN / FAIL / SKIP / 全量 colour）

| 对比 | RUN | FAIL | SKIP | 全量 colour | 读法 |
|---|---|---|---|---|---|
| `r1` vs `r3`（普通形，未换根 ↔ 换根，无变异） | 空／空 | 空／空 | 空／空 | **完全相同** | 普通形换根＝**字面 no-op**（机制：`SealableTempDirForTest124` 走 `resolveProbeRoot`，无链接时返回同一串） |
| `r2` vs `r4`（软链形，无变异） | 空／空 | **only-`r2`＝那 11 枚**，only-`r4`＝∅ | 空／空 | 同 11 枚 `FAIL→PASS`，**无第三枚动** | 本格要收的那笔债收掉了；**没有谁消失、没有谁由红转 SKIP** |
| `r5` vs `r7`（普通形，各自 MUT-D） | 空／空 | 空／空 | 空／空 | **完全相同** | ⇒ 派单参照"普通形逐名差集为空（＝no-op）"在本程读数上成立〔独立复现〕 |
| `r6` vs `r8`（软链形，各自 MUT-D） | 空／空 | 空／空 | 空／空 | **完全相同** | ⇒ 参照"换根后 MUT-D 软链形那 11 枚一枚不落全转红"的**枚数与名册**成立〔独立复现〕，但**"转红"这件事对不上归因**，见 §2.4 |
| `r2` vs `r6`（软链形，无变异 ↔ MUT-D，同一棵未换根树） | 空／空 | only-`r6`＝**对照组 3 枚** `PASS→FAIL` | 空／空 | 同 | 变异真落地的凭据（分母在两形都红**不能**自证变异） |
| `r4` vs `r8`（软链形，换根树：无变异 ↔ MUT-D） | 空／空 | only-`r8`＝**分母 11 枚 ＋ 对照组 3 枚＝14 枚** 全 `PASS→FAIL` | 空／空 | 同 | 换根后的树**没有变哑**：同一发 MUT-D 下软链形 14 枚一起响 |

`r1`↔`r3`↔`r5`↔`r7` 的 `=== RUN` 名册两两 `diff` 行数＝0，`r2`↔`r4`↔`r6`↔`r8` 同样＝0 ⇒ 八发里没有"某枚用例
整发不见了"这一档。

### 2.3 换根后仍带牙：红Reason 换了轨（这条是本格真正的凭据，不是"变红了"）

同一枚判据分支（`refusalCreditsLink137`，八发 `ACTIVE_CREDIT_BRANCHES=3` 恒活）下，按**红在哪一句**分：

| 发 | `returned nil, i.e. it sealed through a symlink and reported success` | `which does not credit the link this case planted` | 日志里出现 `through the link at /ac4link` |
|---|---|---|---|
| `r2`（未换根·无变异·软链） | 0 | **10** | 22 |
| `r4`（换根·无变异·软链） | 0 | 0 | **2**（只剩 §3 那两枚 119 用例的日志行） |
| `r6`（未换根·MUT-D·软链） | 2（都是 118 那两枚已解析根的对照） | **10** | 22 |
| `r8`（换根·MUT-D·软链） | **10** | **0** | 2 |
| `r5`／`r7`（普通形·MUT-D，两棵树） | 10 | 0 | 0 |

⇒ **未换根时那 10 句红说的是"拒是拒了，可拒的不是我种的链接"**（宿主 `/ac4link` 在第 1 个组件替它答了），
**换根后同一批用例的红变成"根本没拒"**——也就是 `assertRefused113` 现在读到的是被测走查**自己的**失效，
和票 118 那两枚天生已解析根的对照**同一机制、同一句话**。这正是判据②要的"不再是装饰"的证明。
（10 而非 11：`TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth` 的顶层 FAIL 由它 depth-1 子测的报错行
聚上来，母项自己不打印红句；枚数账仍以 `-v` 的 `--- FAIL` 行为准。）

### 2.4 对参照值的一处**归因不符**（照实报，没改判据去对读数）

派单给的参照（＝票面 `:173` 那块，终裁方在 `4a0d7a4` 上补跑 `tree-mutd-swap` 量的）说：
"换根之后 MUT-D 软链形那 11 枚**应一枚不落全转红**"。**枚数与名册我这程对上了**（`r8` ＝ 那 11 枚，
`r6`↔`r8` 差集为空）；但**"转红"在我的树上不是换根造成的**：AC#2 的判据收紧已经落在两棵树共同的地基里，
所以**未换根＋MUT-D 的 `r6` 也是同样那 11 枚红**。那句话在 `4a0d7a4`（AC#2 之前）成立——那时 `r6` 那一形
是 11 枚绿（AC#3 r1 件 §0.7 的 `ac3-d-loose-link` 45/24/3/3＋15/0 就是那一发）。
⇒ **登记为一条口径更正**：换根这一味药在 AC#2 之后的树上**新增的**凭据只有两条——
① `r2→r4`（未变异软链形那 11 枚由红转绿，债收掉）；② §2.3 的**红Reason 换轨**（`r6`→`r8`：credit→nil）。
只报"11 枚全红"会读起来像换根白做，只报"名册差集为空"会读起来像它没有代价，两条都得留。

### 2.5 判据③：本格没有拿"未变异的软链形也 11 枚红"当凭据

`r2` 那一发在票面上（`:140`）与 AC#3 两枚终裁里都写明**修之前就在响**。本程只把它用作**债的起点**（§2.2 第 2 行），
**没有**用它证明任何"修好了"。能证明换根后仍有牙的读数是 `r4`↔`r8` 那一对（同一棵树、同形状、
只差变异）与 §2.3 的红Reason 换轨，两对都是**同树内自比**，不跨形、不跨树凑数。

### 2.6 判据④：125 那三枚自拒探针在八发里的颜色（逐名）

| 发 | `…SeamProbeShapesAreBuiltOnAResolvedRoot125` | `…SeamGuardStillRefusesEveryHostileShape125` | `…SeamAcceptsTheHonestPOSIXAnswer125` |
|---|---|---|---|
| `r1`／`r3`（普通·无变异） | PASS（＋2 子测 PASS） | PASS（＋2 子测 PASS） | PASS（＋3 子测 PASS） |
| `r2`／`r4`（软链·无变异） | **SKIP** | **SKIP** | **SKIP** |
| `r5`／`r7`（普通·MUT-D） | PASS | PASS | **FAIL**，红行＝子测 `control_unresolved_root_still_refused_by_the_floor_itself`（同发另两枚子测 PASS） |
| `r6`／`r8`（软链·MUT-D） | **SKIP** | **SKIP** | **SKIP** |

⇒ 三枚探针**没有一枚被改成 SKIP、也没有一枚从 SKIP 转出来**：软链形四发里名册逐名恒定＝3 枚 SKIP
（`plantTemps125` 看到基根穿链接回自己就自拒，票 124 批次 3b 的既有形状），普通形四发里全部**真跑**。
`r5`/`r7` 那一枚转红**不是探针失守**，恰恰是它的 `control_unresolved_root_still_refused_by_the_floor_itself`
那一腿在**被测地板**上抓到了 MUT-D 破口——它因此是 §1.2 第 3 条"谁守未解析形"里**会响的那一枚**（只在普通形）。

### 2.7 顺带一枚树身份旁证（不是判据，只防"读了错树"）

同一枚 108 文件在两棵树上的报错位不同：`r2` 打印 `ancestor_separator_108_other_test.go:79`（未换根树），
`r4` 打印 `ancestor_separator_108_other_test.go:92`（换根树，本程在那枚文件头加了 13 行说明）⇒
读的确是两棵不同的树，与 `gate 96` 的 `RAWROOTS=7`／`RAWROOTS=0＋SWAPPEDROOTS=9` 一致。

---

## §3 票 119 那三处 sentinel-only 内联断言：判定＋能红能绿的读数（**本程一字未改**）

### 3.1 先现量三处的位置与形状（开工锚点 `182daed`，不抄派单行号）

```
$ git grep -n "errors.Is(err, winsec.ErrUnresolvedPath)" 182daed -- internal/winsec/dataroot_symlink_119_other_test.go
internal/winsec/dataroot_symlink_119_other_test.go:203   （母块 :201 `if err == nil` + :203 `else if !errors.Is(...)`）
internal/winsec/dataroot_symlink_119_other_test.go:251   （母块 :249 + :251）
internal/winsec/dataroot_symlink_119_other_test.go:308   （母块 :306 + :308）
```

枚数＝**3 处，与派单在 `182daed` 量的相同**（同文件另有 `:140`/`:179`/`:204`/`:252`/`:309` 等 `t.Errorf` 行，
不是 sentinel 判据，别数串）。三处的形状确实是"只判 sentinel、不判记名"：拿到 `ErrUnresolvedPath` 就算过，
**不核那枚被记名的链接是不是它自己种的**——与 AC#2 收紧前的 `assertRefused113` 同族（票面 `R-137-2` 说的"不走 helper 的内联腿"就是这一族）。
本文件 5 处递根点（`:128`/`:154`/`:194`/`:220`/`:286`）本程**一枚未动**。

### 3.2 逐处判定（换根之后它们还有没有区分力）

| 处 | 它属于哪枚用例 | 它喂的根 | 判定：本程换根后它失去区分力了吗 | 读数（§2 里哪一发证的） |
|---|---|---|---|---|
| `:203` | `TestAC1POSIXUnresolvedSymlinkedRootStillRefused119`（`:193`） | **raw `t.TempDir()`**（`:194`） | **是——软链形失去，普通形没失** | `r5`/`r7`（普通形·MUT-D）**FAIL**；`r6`/`r8`（软链形·MUT-D）**PASS**，且四发都在 `=== RUN` 名册里（不是 SKIP）⇒ 同一发变异下它只有一形会响＝在 macOS 那一形**零区分力**。`r4` 日志里它自己就把拒因打印成 `through the link at /ac4link`（宿主链接），不是它种的 `varlink119` |
| `:308` | `TestAC2POSIXInjectedTestDataDirStandsAsDeclared119`（`:285`） | **raw `t.TempDir()`**（`:286`） | **是——同上，软链形失去** | 与 `:203` 逐发同色：`r5`/`r7` FAIL、`r6`/`r8` PASS；`r4` 日志同样打印 `through the link at /ac4link`（它种的是 `injlink119`） |
| `:251` | `TestAC3POSIXLinkInsideAResolvedDataRootStillRefused119`（`:219`） | **已解析根**（`:226` 走 `cleanSpelling119`，不是 raw `t.TempDir()`） | **没失去——两形都还带牙** | `r5`/`r6`/`r7`/`r8` **四发全 FAIL**（＝本程的对照组之一），未变异四发全 PASS ⇒ 它不属"随另两枚一起改"的那一族；把它一并改是**扩大 finding 面**，本程不这么报 |

### 3.3 finding（要改的东西，交回由编排者在票 119 上另开一格承接；**本格不吞、不改码**）

**建议的那一格**：把 `dataroot_symlink_119_other_test.go:203` 与 `:308` 两枚内联判据按 AC#2 的形状收紧——
共用现成的 `refusalCreditsLink137(err, planted)`（同属 `package winsec_test`，函数体在 `9c0f546` 的
`placement_symlink_113_other_test.go:160`，**零新增依赖**、**零生产码改动**），并**同时**把这两枚用例的
`base := t.TempDir()`（`:194`、`:286`）换成 `SealableTempDirForTest124(t)`。

**为什么两味要同批**（这条是本程读出来的，不是推的）：只收紧判据不换基根，判据会在软链形**恒红**
——那正是 AC#2 单独落地时那 11 枚的处境（`r2`），而换根那半恰是让它"能绿"的那半（`r4`）。
`:251` 不在这一格里（它已走已解析根）。

**代价（本格换根之后已经落在盘上的账，票 119 那格必须接住）**：`r4` 与 `r8` 的日志里，
"未解析根被拒"这条在**软链形**只剩**能打印**、**不能区分**的两枚 sentinel 断言在守（§3.2 前两行），
外加 125 那三枚在这一形**按设计 SKIP**（§2.6）⇒
**票面 `:175` 那句"换根之后未解析那一形只剩票 119 那两枚 raw `t.TempDir()` 用例持有"成立，
本程把它从"持有"再读成"持有但不带区分力"**：在 macOS 那一形，目前**没有任何一枚**断言能指出
"地板拒的是宿主链接、不是它自己种的链接"。这条不是本格的破口（本格的判据是逐枚判定 7 处），
但它是本格换根的**已知代价**，必须有名有姓地落在票 119 的一格里，否则读者会以为那两枚还在守同一件事。

**票 125 那三枚探针在普通形带牙**（§2.6 的 `r5`/`r7`），所以"未解析形"不是全仓无人守——
失守面**只在软链形那一半**。

