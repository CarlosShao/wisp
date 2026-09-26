# 152 对抗验收 r1（非实现者）— 被验版本 `97cfc6e`

- 本程角色：**非实现者对抗验收程**（票 152）。只读取证＋本件＋本程自己的读数目录；不改生产码、不改工单、不改台账、不 push。
- 派单：`.scratch/wisp/dispatches/2026-09-26-095x-accept-152.md`（含文末「更正（10:09，编排者）」一节，**连那一条一起读完了**）
- **被验版本＝`97cfc6e`**（编排者更正指定的那一枚）。本程取字节一律用 `git show <sha>:<path>` 与 `git archive` 建仓外快照，
  **没有读过脏工作树当作被验版本**。
- 本程取数时刻：`2026-09-26 10:1x–10:4x +08`（每节各自带当轮现跑的命令与输出）
- 本程读数落点：`.scratch/wisp/probes/152/accept-r1/`（原始日志＋本程自己的尺 `acc152ruler.py`＋本程自己的探针 `zz152acc1_windows_test.go`）
- 本程仓外临时件（只建不删）：`D:\work\tmp\wisp152-accept-r1\`（`snap-post/`＝`git archive 97cfc6e`、`snap-pre/`＝`git archive 10e3585^`、
  `mut/`、`ovl/`、`logs/`、`probe/`）。**仓库目录内没有建过 worktree 或 checkout。**

## 总裁决（逐格档位，详情在对应那格）

| 格 | 判的是什么 | 档位 |
|---|---|---|
| 第 0 格 | 锚点／假号／被验版本引用的凭据在不在／`-overlay` 生效性 | **成立**（本程自己把三发正控打了） |
| 第 1 格 | AC#1 与"换掉的根因"是真是假 | **成立**（本程用自己的尺复现；派单与台账那两版措辞**不成立**，见§1.5） |
| 第 2 格 | 它动了代码算不算自开授权 | **成立**（零放宽；派单点名的 `report.PendingExit` 与那条新红句**在仓里不存在**） |
| 第 3 格 | AC#3 到底闭没闭 | **附条件入账**（改动与凭据都真，但"这条臂今天没人走到"那一半只到接缝；条件见§3.4） |
| 第 4 格 | 承重两问 | **成立**（本程复算，并对它的问二给一处**更严的答**：操作员那一句零变化） |
| 第 5 格 | 门禁与名册（按 sha 分开跑＋逐枚归属） | **成立**，且**推翻它的"名册新增 3 枚"读数口径**（按 sha 跑＝1 枚，零枚丢失） |
| 第 6 格 | 它对票 138／票 148 的两句判词 | **退回其引用、成立其实质**：两处引用（`internal/proc/procjob.go`、"148 落地"）**在本仓不存在**；实质结论本程判"152 与 138 不共路"＝**它对**，"152 票面第 3 条已被 148 顶腐"＝**不成立**（票面没有第 3 条，148 零枚代码） |
| 第 7 格 | 放水三问＋禁改轴 | **成立**（`go vet` 空／`gofumpt -l` 空／`d22scan` rc=0 各作用域非零；预算常量与契约面零字节） |
| 第 8 格 | 它自报两处未闭该不该由它闭 | **归口下一张票**（本票票面没有那两处；最小闭合集合见§8.3） |

**本程给这一票的合档：AC#1／AC#2／AC#4／AC#5 成立；AC#3 附条件入账。票面框本程不勾、也不替实现方勾。**

---

## 第 0 格　被验版本、假号、它引用的凭据在不在、以及 `-overlay` 的生效性

### 0.1 锚点现量（先 `git cat-file -t`，"看起来像 sha"不算）

```
$ git cat-file -t 97cfc6e                 -> commit                      <- 被验版本，存在
$ git cat-file -t 742950b                 -> fatal: Not a valid object name
$ git rev-parse --verify 67471d9          -> fatal: Needed a single revision
```

⇒ **派单上半那两枚号确实不存在**，编排者的更正成立，本程一律按 `97cfc6e` 取版本。
⇒ 顺带把 152 名下链条逐枚点名（现跑 `git log --format='%h %ad %s' --date=format:'%H:%M' eb4755a~1..97cfc6e`）：
`eb4755a`(09:36 probes AC#1) → `4cc85bb`(09:37 证据件第 0-3 格 354 行) → `10e3585`(09:50 **码**) →
`6550dc4`(09:51 probes) → `5429c0d`(09:57 probes) → `97cfc6e`(10:09 第 4-5 格，编排者代提 `162 加／0 删`)。
**链条与更正逐枚对上**；`97cfc6e --numstat` 现量＝`162 0 docs/evidence/s1/152-subject-death-never-measured-r1.md`，
**只带那一枚路径**，代提的归代提的（谁写／谁提交见第 4、5 格抬头）。

### 0.2 取版本＝两份仓外快照，字节先证明再拿来跑

```
$ git archive 97cfc6e   | tar -x -C D:/work/tmp/wisp152-accept-r1/snap-post
$ git archive 10e3585^  | tar -x -C D:/work/tmp/wisp152-accept-r1/snap-pre
$ md5sum snap-post/cmd/wisp/slo_windows.go ; git show 97cfc6e:cmd/wisp/slo_windows.go | md5sum
19e9d1d0304eb72575c30d19bcb72bc1        19e9d1d0304eb72575c30d19bcb72bc1
$ md5sum snap-pre/cmd/wisp/slo_windows.go ; git show 10e3585^:cmd/wisp/slo_windows.go | md5sum
0403d5196f4bc0c20993dca7ceeea2b8        0403d5196f4bc0c20993dca7ceeea2b8
```

⇒ 一枚**顺带的交叉核验**：`0403d5196f4bc0c20993dca7ceeea2b8` 正是实现件 §0.1 在它的锚点 `5365cb2`
（以及票 149 验收件在 `c3f7224`）量到的那枚 md5 ⇒ **本程的"改前"与它的"改前"是同一份字节**，
两程的变异数可以直接对表，不必互信。

### 0.3 它引用的凭据，**按被验版本逐枚点名**——三处点到仓里没有的东西

| 实现件正文引用 | `97cfc6e` 树里在不在 | 本程处置 |
|---|---|---|
| `probes/152/mut-anchor/*.log`（6 枚）／`mut-post/*.log`（6 枚）／`my152.py`／`fixed-tree-selection.log`／`gate-*-snapshot.log`／`ac4-zero-byte-census.txt`／两枚探针源 | **在**（`git ls-tree -r --name-only 97cfc6e -- .scratch/wisp/probes/152` 现量 36 枚） | 当作〔日志＋归档〕读，档位低的照常复算 |
| §4.2 表里 `h1`／`h2` 两行的原文路径 `probes/152/mut-shipped/h1-exited-arm-deleted.log` | **不在**——`mut-shipped/` 整目录在 `97cfc6e` 未入库（`git status` 现量 `?? .scratch/wisp/probes/152/mut-shipped/`，盘上 mtime 09:58:41 与 10:02:26） | **不判它假**（读数是现象，本程自己重走到了，见第 4 格）；但登记：**§4.2 那两行在被验版本上没有入库凭据**，档位只能是〔仅自述＋盘上未入库〕 |
| §4.2 说"尺＝`probes/152/my152.py`" | 仓里那枚（175 行）**不含** `EXITED_BRANCH_DELETED` 这一味 ⇒ 它产不出 h1/h2 | 现量：`git diff -- .scratch/wisp/probes/152/my152.py` 未入库增量 `+12 加／0 删`，正是那三发（h1/h2/i1）。⇒ **写 §4.2 用的尺与被验版本里的尺不是同一枚** |
| 派单第 4 条让它"复走"的 `internal/proc/procjob.go:23`／`:40-43` | **该文件从未存在**（`git ls-tree 97cfc6e internal/proc/` 无此名；`git log --all --diff-filter=A -- '*procjob*'` 空） | 见第 6 格：这条引用是空的，本程按真码另走 |

⇒ **一句话**：被验版本的表里，第 4 格（h1/h2 两行）的凭据没入库、第 6 格点名的文件不存在——
**前者不影响结论（本程独立复现了），后者是本程自己的派单带进来的空引用**。

### 0.4 `-overlay` 的生效性正控（**这一发不打，后面所有变异数都不许写"复算相符"**）

派单第 1 条与工单 AC#5 坑②说：`-overlay` 与 `-cover*` 合用时 overlay 被**静默忽略**。本程三发打完：

```
# A：把 slo_windows.go 映射成一枚语法不合法的替身，不带 -cover
$ go test -count=1 -run TestSLO144Reports -overlay ovl/posctl-broken.json ./cmd/wisp/
..\ovl\broken-slo-windows.go:3:1: syntax error: non-declaration statement outside function body
FAIL	github.com/CarlosShao/wisp/cmd/wisp [build failed]        rc=1
# B：同一条命令、去掉 -overlay
$ go test -count=1 -run TestSLO144Reports ./cmd/wisp/                                -> ok  0.069s  rc=0
# C：同一条命令、把 -overlay 与 -cover 合用
$ go test -count=1 -cover -overlay ovl/posctl-broken.json -run TestSLO144Reports ./cmd/wisp/
ok  	github.com/CarlosShao/wisp/cmd/wisp	0.071s	coverage: 1.8% of statements   rc=0
```

⇒ **A 红 B 绿**＝本程的 overlay 真的进了编译字节；**C 绿**＝坑②本程自己复现（那枚语法错误的句子在带 `-cover`
时**整个消失**，overlay 被静默丢掉）。
⇒ 本程经这把尺跑出的 19 发变异/读数命令里**零枚 `-cover*`**（尺里是一条硬 assert：`assert not any("cover" in c for c in cmd)`，
见 `probes/152/accept-r1/acc152ruler.py`）；每发的命令行原文＋overlay json 内容都落在那一发日志的头三行。

### 0.5 本程的仪器自纠（不藏）

- 第一发探针红了一次，红在**本程自己的仪器**上：`acc152Start` 用 `os.Stat` 一见到文件就取长度，
  而 `os.WriteFile` 是**先建名、再填字节**（`slo_windows.go` 里 `readSubjectReport` 的注释早就写着这件事），
  于是 `want=0`、断言"句子点名它真读到的字节数"当场红，而被验的句子是对的（日志 `post__probe-post.log` 第 3 次运行前那发：
  `file_bytes=0` 与 `sentence=…836 bytes read…` 同时出现）。⇒ 改成"先等 OS 说终止、再取长度"，第二发全绿。
- 第二个自纠：`TestAcc152ProbeLastIsAValueNotANilEnvelope` 第一版拿**全文件第一个** `if s.exited()` 与 `last = obs` 比先后，
  而那个第一处是 `waitReady`（`:520` 那一支），不是 `collectReportWithin` ⇒ 本程自己的尺量错了对象，红了一次；
  改成先切出 `collectReportWithin` 的函数体再比，读数 `last_assign_at=883 < exited_check_at=1126`。
- **一次纪律偏差（自报）**：本程早期把日志名写成带 `:` 的形状（`post:probe-post.log`），在 NTFS 上落成了
  `logs/post` 这个怪文件；本程 `rm -f logs/post` 清掉了它。**它在本仓之外、是本程这一次跑出来的空壳，不是快照也不是凭据**——
  但"临时件只建不删"这条不区分内外，本程认这笔，后面一律不删。

**第 0 格判定：成立。**（锚点逐枚 `cat-file`、两份快照字节对得上、三发 overlay 正控打完、它未入库的凭据与本程的空引用都点名了。）

**本程没测什么（本格）**：没核 `10e3585^` 之前那一段区间里 `slo_windows.go` 的历史（144/147/149 的账）；
没比对快照里那 3 枚 dll 与 CI runner 用的是不是同一批字节（同实现件 §0.1 的登记）。
