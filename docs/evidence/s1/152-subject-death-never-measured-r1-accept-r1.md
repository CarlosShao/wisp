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
| 第 3 格 | AC#3 到底闭没闭 | **成立**（选支、四行判据、"未修码上响不响"全部本程独立复现；那条臂今天没有世界形状走到＝登记，不是本票判据） |
| 第 4 格 | 承重两问 | **成立**（本程复算，并对它的问二给一处**更严的答**：操作员那一句零变化） |
| 第 5 格 | 门禁与名册（按 sha 分开跑＋逐枚归属） | **成立**，且**推翻它的"名册新增 3 枚"读数口径**（按 sha 跑＝1 枚，零枚丢失） |
| 第 6 格 | 它对票 138／票 148 的两句判词 | **退回其引用、成立其实质**：两处引用（`internal/proc/procjob.go`、"148 落地"）**在本仓不存在**；实质结论本程判"152 与 138 不共路"＝**它对**，"152 票面第 3 条已被 148 顶腐"＝**不成立**（票面没有第 3 条，148 零枚代码） |
| 第 7 格 | 放水三问＋禁改轴 | **成立**（`go vet` 空／`gofumpt -l` 空／`d22scan` rc=0 各作用域非零；预算常量与契约面零字节） |
| 第 8 格 | 它自报两处未闭该不该由它闭 | **归口下一张票**（本票票面没有那两处；最小闭合集合见§8.3） |

**本程给这一票的合档：AC#1／AC#2／AC#3／AC#4／AC#5 全部成立；两处缺陷记在账上但不构成退回（§2.2 末类、§4.4）。票面框本程不勾、也不替实现方勾。**

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

---

## 第 1 格　AC#1 与被换掉的那枚根因：本程用自己的尺重走到"哪一句印出来"

### 1.1 本程的判据（跑之前写死）

派单第 1 条要求"自己复走到'摘掉 `s.exited()` 后 `--- FAIL: 0 枚`'那一步"。本程把它拆成三条**可失败**的读数，
全部经本程自己的探针（`probes/152/accept-r1/zz152acc1_windows_test.go`，仓外 overlay 注入，
`ls cmd/wisp/zz152acc1_windows_test.go` → 不存在，`git status --porcelain cmd/wisp` 干净）：

1. 真子进程（`exec.Command(os.Args[0], …)`、真 pid、真退出码 7）、真落在盘上的真文件（子进程自己 `os.WriteFile`）、
   真 `os.ReadFile`（`readReportFile` 留 nil＝生产接线）——**四真一条不退**，与实现件 §1.1 同一条纪律；
2. 外部证人：`windows.OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION)` ＋ `GetExitCodeProcess`，
   **先等 OS 说终止、再开循环**（排除"死得不够快"）；
3. 断言方向取**否证式**：探针 A 在"`exited()` 那一支今天真走到了／句子换了邻居"时红，探针 B 在"收尸之后仍不印 exited 句"时红。

### 1.2 读数（`logs/post__probe-post.log`，RUN=3 全绿；同一发也在 `asisp-probe-only` 重跑过一遍）

```
ACC|setup|pid=57884 file_bytes=836 os_says=terminated exit_code=7 cmd.ProcessState==nil?true
ACC|A|exited_before=false exitCode_before=-1 elapsed=2.0040809s
ACC|A|sentence=wisp slo: subject 57884 never wrote a complete report within 2s (last read: 836 bytes read, document still open at offset 836: the tail had not arrived)
ACC|A|after_loop|exited=false exitCode=-1
ACC|A|last_reading|state=unwritten offset=836 bytes=836 summary="836 bytes read, document still open at offset 836: the tail had not arrived"
ACC|B|after_reap|exited=true exitCode=7
ACC|B|elapsed=565.2µs sentence=wisp slo: subject 61588 exited (code 7) without writing its report (836 bytes read, document still open at offset 836: the tail had not arrived)
ACC|C|zero_value_summary="0 bytes read, document still open at offset 0: the tail had not arrived"
ACC|C|inside_collectReportWithin|last_assign_at=883 exited_check_at=1126
```

⇒ **OS 说它死了（code 7）、`cmd.ProcessState` 仍是 nil、`s.exited()` 假、循环烧满整段预算、印的是邻居那句**：
被验版本在真子进程上复现，**根因换轨那一判＝成立**。翻转它的唯一变量是本程自己补的那一发 `cmd.Wait()`（B：0s 内当场红）。

### 1.3 机制名册（本程现跑，按 97cfc6e 行号）

```
$ grep -n '\.Wait()\|\.Run()\|StartInJob\|ProcessState' cmd/wisp/slo_windows.go
498:	p, err := rt.Job.StartInJob(cmd)
833:	return s.cmd != nil && s.cmd.ProcessState != nil && s.cmd.ProcessState.Exited()
840:	return s.cmd.ProcessState.ExitCode()
855:		_ = s.cmd.Wait()            <- 全文件唯一一处 Wait，在 stop() 里面
$ grep -c "go func" cmd/wisp/slo_windows.go            -> 0    （没有哪条 goroutine 会去收尸）
$ grep -n "defer subject.stop()\|defer instrumented.stop()\|collectReport()" cmd/wisp/slo_windows.go
371:	defer subject.stop()
380:	defer instrumented.stop()
392:	observer, err := instrumented.collectReport()   <- 收尸是 defer，判死在它里面
$ grep -n "s.exited()" cmd/wisp/slo_windows.go -> 520 / 820 / 837 / 851
```

⇒ 实现件 §1.5 那张名册**逐行对得上**（它引的 794/801/816 是本程 833/840/855 加上 `10e3585` 自己插入的 41 行）；
⇒ `s.exited()` 的**两处用法**都在任何 `Wait` 之前（`:520` 属 `waitReady`、`:820` 属本循环），
本程只把 `:820` 量成了读数，**`:520` 那一支本程没重跑**（实现件 §2 有它的一发 60s 读数，档位〔日志＋归档，抽验〕，
本程按静态名册判它同因、不判它已复算）。

### 1.4 AC#1 本身

工单 AC#1 要的是"量成一次读数，而不是再造一个假对象"，并明写"不许预填结论：量出来那条路今天产不出红句，
那正是要的答案"。本程复算：实现件 §1.3 的七格矩阵里，与本程同形的第 1、4 两格（unreaped→timeout 腿、REAPED→exited 腿）
**本程独立走到了**；它 §3 那发端到端（真 `wisp slo`、真 `TerminateProcess`、操作员 stderr 原文、rc=2 不降级）
本程**没有复跑**（理由见§9 第 3 条）。⇒ **AC#1 成立**，且它落在工单预填的"ⓑ/ⓒ 停手上报"那一支上是对的。

### 1.5 两句要更正的措辞——一句是派单的、一句是台账的（**都不是交付表里的**）

| 出处 | 那句话 | 本程现量 |
|---|---|---|
| 派单第 1 条（＝台账 `A275②`）："实现件断言……`saw[1]` 是**空信封**、`last` 被赋成一条都没见过的空报告、**`last` 是 nil**、那一支今天根本不进 ⇒ **摘掉 `s.exited()` 也零枚红**" | ① 被验的表里**零次**出现"空信封／`saw[1]`／`last` 是 nil"（`grep -c` 命中 0，见 §0 现跑）；② `subjectReportRead` 是**值类型**、`last = obs` 在两处 give-up **之前**无条件执行（`last_assign_at=883 < exited_check_at=1126`）⇒ 没有"nil"这个状态可谈；③ 真要谈零值，它印出来的是 `0 bytes read, document still open at offset 0`（本程量到的），那是**指错位置的数**、不是空信封；④ **"摘掉 `s.exited()` 零枚红"直接被本程否掉**：`x1-exited-arm-deleted`＝RUN 21／PASS 13／**FAIL 1**，红名 `TestSLO149ExitedGiveUpSentenceCarriesTheLastReading` | **不成立**。成立的是它的**条件式**：摘掉整条支**并且**摘掉 case 13 ⇒ 0 枚红（`x2`＝RUN 20／PASS 13／**FAIL 0**）——这句话在交付表 §4.2 问一 ② 里就是这么写的，**表是对的、台账与派单把它抄歪了** |
| 任务书/派单第 2 条："为此新增了 `report.PendingExit` 字段和一条新红句" | `git grep -n "PendingExit" 97cfc6e` 全仓命中 **2 处，都在文档里**（这枚派单本身、台账 `A275`）；`*.go` 里 **0 处**，工作树里也 0 处。被验的码改动只新增了一枚**常量** `offsetUnknown` 与一枚函数 `contradictSubjectReportErr`（见第 2 格） | **不成立**：那三枚"工作树里读到的改动"没有进过任何一枚 commit（`A275①` 是 09:5x 按当时脏树逐行读的，`A276①` 之后它复活改成了 AC#3 那一版）。⇒ **本程按 97cfc6e 的真改动裁，不按这句裁** |

**第 1 格判定：AC#1 成立、根因换轨成立（本程复现）；台账 `A275②` 与派单第 1／2 条那两版措辞退回更正。**

**本程没测什么（本格）**：`waitReady`（`:520`）那一支本程没跑真读数；端到端（真 CLI＋真杀）本程没复跑；
`startSubject` 指派失败那一支（`internal/proc/jobscope_windows.go:117-119`，唯一可能出现在 `collectReport` 之前的收尸点）
本程与实现件**都没构造过**——两程一致挂着，没人在这里宣称测过。

---

## 第 2 格　它动了代码——是不是给自己开的授权

**先纠派单那句**（任务书第 2、3 条都按"新增了 `report.PendingExit` 字段＋新红句"来问）：**那两味在被验版本里不存在**（§1.5 表第 2 行）。
本程因此把这一格改成问**真动的那一味**：`10e3585` 交付的改动面。

### 2.1 被验版本的真实改动面（现量）

```
$ git show --numstat 10e3585        -> 88 1  cmd/wisp/slo_report_144_windows_test.go
                                          56 17 cmd/wisp/slo_windows.go
$ git show --name-only 97cfc6e --format=''  -> 只有 docs/evidence/s1/…（第 4-5 格，编排者代提）
```

新增符号只有四枚：常量 `offsetUnknown = -1`、函数 `contradictSubjectReportErr`、`contradictionOffset` 的**签名收窄**
（两参→单参）、用例 `TestSLO152CorruptLegWithNoNamedPositionRefusesToBorrowOne`（case 14）。
**没有新增字段、没有新增状态、没有动 `exited` 那条支的任何一个字**（与它 §4"AC#2 那一半一字节生产码都没改"一致）。

### 2.2 删除行逐枚点名（全部 18 行，本程 `git show 10e3585 | grep -E '^-[^-]' | cat -A` 现跑，分四类）

| 类 | 行 | 是不是断言 | 本程判 |
|---|---|---|---|
| 注释重排 | `subjectReportRead.offset` 语义段 8 行、`Only the -1 leg is never printed…` 2 行、`contradictionOffset` doc 3 行 | 否 | 无放宽；且新注释**没有把没测的说成测过**（它写"case 14 is what makes the new answer checkable at this seam"，seam 这个词是准的） |
| 签名与调用点 | `func contradictionOffset(err error, inputOffset int64) int64 {` 与 `obs.offset = contradictionOffset(err, dec.InputOffset())` | 否 | 同步收窄，全仓无第二个调用者（本程 `grep -n contradictionOffset cmd/wisp/*.go` 只命中该调用点＋case 14） |
| **唯一行为变化** | `return inputOffset` → `return offsetUnknown` | 不是断言，是被测行为 | 见 §2.3 与第 3 格。**没有任何既有断言检查过这条臂**——本程自己在改前字节上复算：`p1`（整臂换成 0）与 `p2`（换成 -1）都是 RUN 20／PASS 13／**FAIL 0**，与 `asis-pre` 四数逐枚相同 ⇒ 它说"装饰腿"是**量出来的**，不是嘴上说的 |
| 句子搬家 | `obs.err = fmt.Errorf("%d bytes contradict a subject report at offset %d: %w", obs.bytes, obs.offset, err)` → `contradictSubjectReportErr(obs.bytes, obs.offset, err)` | 否 | 本程逐字节比：那枚格式串字面量在改前出现 1 次、改后出现 1 次、**完全相同**，只有实参名从 `obs.bytes, obs.offset` 换成 `size, offset` ⇒ 命名臂零变化（case 11 现绿为第二证） |
| **删了没补回的一行** | `// Case 13 - AC#3's ruling is ⓐ: the s.exited() branch SHOULD name where the last` | 注释 | **缺陷登记（本程新增，实现件没报）**：case 13 的注释现在从句中开始（首行变成 `// reading stopped, for the reason…`），而**被删那半句是票 149 那一格裁定（ⓐ）的记录**。它不是放宽断言（票 152 的 AC#2 改判 ⓒ 之后，"ⓐ"这句话本来就过期了），但**别人票面的裁定出处被本票顺手抹了**。归口见 §8.2 |

⇒ 测试文件的删除行**只有那一行注释**；`grep -E '^-' … | grep -E 't\.(Errorf|Fatalf|Skip)'` 只命中生产里那枚 `fmt.Errorf` 搬家行。
**零枚既有断言被删、零枚被放宽、零枚 `t.Skip`、阈值与预算常量零字节**（`git diff 10e3585^ 97cfc6e -- cmd/wisp/slo_windows.go` 里
`subjectReportBudget|subjectGrace|subjectPollInterval|subjectReadyBudget` 命中 **0**）。

### 2.3 三问逐答（派单第 2 条原问的三个）

1. **这条新句有没有偷偷放宽任何既有断言？** 没有（§2.2 逐枚）。被它改的那句 `at an offset the decoder did not name` 只长在一条
   **今天没有任何文档能走到**的臂上（第 3 格条件），而走到过的两枚命名臂字符串逐字节未动。
2. **"subject 没有写出一行报告"那个数是量到的还是造出来的？** 派单问的这句**不在被验版本里**——`exited (code %d) without writing its report (%s)`
   是存量句子（`10e3585^:cmd/wisp/slo_windows.go:782` 就有；钉它的 case 13 出自 `b417d31`＝票 149，句里的 `last.summary()` 出自 `95885fb`＝票 144）。
   本程把它的数**量成了读数**：探针 B 印出的是真 836 字节（`…(%s)` 里的 `836 bytes read, document still open at offset 836`），
   那 836 枚字节是真子进程 `os.WriteFile` 落盘、由真 `os.ReadFile` 读回、`readSubjectReport` 真分类成 `unwritten` 的。⇒ **量到的**。
   顺带把台账那句钉死：**这条支今天在生产接线上不响**，所以"量到的"是**它的形状**，不是**它的到达**——两程都没说过后者。
3. **有没有为新增那一味造一发"今天不响的检"（恒真判据族）？** case 14 **不是恒真**，本程三发变异都能打红它：
   `y1`（臂交回 0）／`y6`（臂交回 1）／`y4`（渲染器删掉承认臂、永远印数字）＝**各 1 枚红，红名都是 case 14**；
   `y7`（渲染器永远不印数字＝命名臂被消音）＝ **2 枚红：case 11 与 case 14**（本程现量：`RUN 21/PASS 12/FAIL 2`）
   ⇒ 与它 §5.3"诚实的一半：那一臂不是 case 14 独占的覆盖面"**同向、且本程自己走到了**。
   但它**响的前提仍然是有人走到那条臂**：本程同样**没有**构造出任何一份真实文档把这枚检叫响过（§3.2）。

### 2.4 授权面

AC#3 的票面文本就是"两选一（删掉／换成会明说自己不知道的形状）"——**动码在授权射程内**，不是自开。
真正越出票面的一点，本程登记在这里：`offsetUnknown` 把**票 147 钉过的 `-1` 语义**（"分类器根本没跑：没读到／没文件／文件读不开"）
**加了一腿**（"解码器跑了但没点名"）。它自己在 §5"没测什么"第 2 条登记了"没做两腿的语义合并验证" ⇒
**诚实、但确实是票面没点头的一小步扩张**。它没有渗进任何契约面（`internal/risk`／`panel`／`approval`／`d22scan`／`thresholds`／golden 零字节，见第 7 格），
所以本程按"票内小扩张＋已登记"记，不按自开授权记。

**第 2 格判定：成立（不是自开授权）。附带两笔要记的：case 13 注释被删那一行（§2.2 末类）、`-1` 语义加腿（§2.4）。**

---

## 第 3 格　AC#3 到底闭没闭——按 §5 与 `10e3585` 的 diff 判，不按任何一方的自述判

编排者的更正里点名："它撞顶前最后一句是'现在我要去实现 AC#3 的改动'，而 AC#3 的码早在 09:50 的 `10e3585` 就落了 ⇒ 别按它的自述算已闭，也别按我的更正算已闭。"
本程的处置：**把 `10e3585` 当作被审对象、把 §5 当作它对那枚 diff 的说法、两样都重走**。

### 3.1 AC#3 那四行判据逐条对

| AC#3 要求 | 被验版本里的东西 | 本程现量 |
|---|---|---|
| "删掉，或者换成会明说自己不知道的形状，两选一、写清选了哪个为什么" | 选了第二支；§5.2 写了为什么不选第一支（Go 每条路径要有返回值，字面删＝编不过） | **本程自己撞了那一发**：`y5-delete-the-fallback-line`（把 `return offsetUnknown` 摘掉）⇒ `slo_windows.go:684:1: missing return`、`=== RUN` 枚数 **0**（＝根本没跑到，不是跑到且红）⇒ "删不掉"是**量出来的**，不是引用 149 的 |
| 三件不算收：加注释 | 除注释外动了 4 枚符号＋1 枚用例 | ✓ 不只是注释 |
| 三件不算收：说成"纵深防御" | 新 doc 里写的是"a decorative leg, not an unmeasured one"，并把自己那两发零读数变化点名（`f1`/`f2`） | ✓ 没走那条说法；且本程在**改前字节**上独立复算到同一形状（`p1`/`p2`＝20/13/0，与 `asis-pre` 四数逐枚相同） |
| 三件不算收："留着反正不害" | 它没这么说，它说的是"能删的只有那句承诺" | ✓ |
| "不许为它开'要它响'的格；判据要能答'这一发在未修码上响不响'" | 答了"响"，凭据 `g1`；并**明写**未修码上叠新用例是编译期就断（`n1`），不以它当读数 | **两半都复现**：本程 `y1`（臂交回 0）→ case 14 红，红句 `…is 0, want -1 (offsetUnknown); any other number is a byte position this error never named…`；`y6`（臂交回 1）→ 同一枚红、数换成 1 ⇒ 判据对**任何**借来的数都响，不是只对 0 响。`pre:n1` → 交付版用例叠到改前字节上 `RUN=0`，编译错原文：`not enough arguments in call to contradictionOffset`／`undefined: offsetUnknown`（两处以上）／`undefined: contradictSubjectReportErr` ⇒ 它那句"编不过是构建状态，不是读数"**站得住** |

### 3.2 那"承认"是不是真承认（本程自己造的三发）

- `y4`（把渲染器里 `offset < 0` 那一臂删掉＝只改值、不改句子）→ **case 14 红，红句原文两条**：
  `no-position corrupt sentence "39 bytes contradict a subject report at offset -1: …" is not the admission shape`
  与 `the admission leaked offsetUnknown as a position`。⇒ 它 §5.2 那句"只在值上加 -1 不够"是**被钉住的**，不是措辞。
- `y7`（反过来：渲染器永远走"承认臂"＝命名臂被消音）→ **2 枚红（case 11 与 case 14）**，
  原文含 `8 of 8 corrupt shapes do not name the byte position the decoder objected at`。
  ⇒ 命名臂**不是** case 14 独占的覆盖面（它 §5.3"诚实的一半"那句，本程同向走到）。
- 但：**今天没有任何一份真实文档把这条臂叫响过**。本程与实现程**都没构造出来**（它 §5"没测什么"第 1 条自己登记了，
  且写得准："见证的是接缝，不是世界"）。本程另加一条静态事实支持这个判断：
  `readSubjectReport` 里能走到 `contradictionOffset` 的错误只有 `dec.Decode` 返回的**非 EOF 类**错误，
  而 `encoding/json` 在那条路上给的就是 `*SyntaxError`／`*UnmarshalTypeError` 两家（两家都点名）——
  本程没有构造出第三家，**也没有证明不存在第三家**（那要枚举 `encoding/json` 的错误全集，含未来 Go 版本）。

### 3.3 那句"现在我要去实现 AC#3 的改动"怎么算

时序现量（本程自己跑的）：`10e3585` 落码 09:50:23 → `6550dc4` 落 12 发变异读数 09:51 → `5429c0d` 落 gofumpt＋交付字节重跑 09:57
→ 证据件 §4/§5 的**盘上写入时间约 10:01**（`find -printf %TH:%TM` 现量该文件 mtime＝10:01，`97cfc6e` 提交 10:09）
⇒ **先有码、有读数，后有那两节的文字**（顺序是正的，不是"预先引用尚未产出的读数"那一族）；
⇒ 那句"现在我要去实现 AC#3"**只能是撞顶前最后一次自述的过期措辞或同格二次自述**，它**不改变 AC#3 已按第二支落地**这件事。
本程据此判：**不按那句话退回**。（它到底指什么，只有那枚程的会话日志知道——本程不替它猜，也不拿它当否证。）

**第 3 格判定：AC#3 成立**（选了哪一支、为什么、四行判据、"未修码上响不响"——本程全部独立复现）。
**残余不是条件，是本格"没测什么"**：那条臂的可达性两程都没量到（§3.2 末行），
`case 14` 的 14b 那三枚真错误在**别的 Go 版本**上是否仍归两家，本程没换版本验（与实现件同一条洞）。

---

## 第 4 格　承重两问——本程的 19 发矩阵，与它对表；并给它问二那一句**更严的答**

本程的尺与它不同源（`acc152ruler.py`，本程自写；每发落盘后**回读**断言"新文本在／旧文本不在"，否则 FATAL 不出数；
命令原文＋overlay 内容印在每发日志头三行；全程零 `-cover`）。被验树＝两份仓外快照，§0.2 已 md5 对表。

### 4.1 逐发对表（左＝本程现量，右＝实现件自报的那一行）

| 本程的发 | 字节基 | RUN | PASS | FAIL | 红名 | 实现件那一行 | 对上没 |
|---|---|---|---|---|---|---|---|
| `asis-post` | 97cfc6e | 21 | 14 | 0 | — | §5.3 `asis` 21/14/0 | ✓ |
| `asis-pre` | 10e3585^ | 20 | 13 | 0 | — | §5.1 `asis` 20/13/0（其锚点 5365cb2） | ✓（且 §0.2 md5 证明两程"改前"同一份字节） |
| `x1-exited-arm-deleted` | post | 21 | 13 | **1** | case 13 | §4.2 `h1` 21/13/1 同红名，红句 837/840/843 三行 | ✓ **逐行对上**（本程读到的是 837/840/843 同三行；重读次数 19418 对它的 18485＝本机差） |
| `x2-…-case13off` | post | 20 | 13 | **0** | — | §4.2 `h2` 20/13/0（逃逸） | ✓ |
| `x3-case13-off` | post | 20 | 13 | 0 | — | §4.2 `e3`（锚点 19/12/0） | ✓ 同形（差一枚 case 14） |
| `z1-exited-loses-last-summary` | post | 21 | 13 | **1** | case 13（4 条红句） | §4.2 `e1`（锚点 20/12/1） | ✓ |
| `z2-…-case13off` | post | 20 | 13 | **0** | — | §4.2 `e2`（锚点 19/12/0 逃逸） | ✓ |
| `y1-fallback-returns-0` | post | 21 | 13 | **1** | case 14 | §5.3 `g1` 21/13/1 | ✓ |
| `y2-…-case14off` | post | 20 | 13 | **0** | — | §5.3 `g1-…-case14off` 表里写 20/13/0＝**对**；**随码交付的注释**写 21/14/0＝**错**（见 §4.4） | 表 ✓／注释 ✗ |
| `y3-case14-off` | post | 20 | 13 | 0 | — | §5.3 `g2` 20/13/0 | ✓ |
| `y4-…-always-numbers` | post | 21 | 13 | **1** | case 14（红句两条） | §5.3 `g3` 21/13/1，引的 `:754` 那句本程逐字复现 | ✓ |
| `y7-renderer-never-numbers` | post | 21 | **12** | **2** | case 11 ＋ case 14 | §5.3 `g4` 21/12/2 | ✓ |
| `y5-delete-the-fallback-line` | post | **0** | — | 编译 | `684:1: missing return` | §5.2 "把 `return offsetUnknown` 摘掉＝missing return" | ✓ |
| `p1/p2`（改前字节把臂换成 0／-1） | pre | 20 | 13 | 0 | — | §5.1 `f1`/`f2` 20/13/0 | ✓ |
| `n1`（交付用例叠改前字节） | pre | **0** | — | 编译 | 3 类未定义／参数数错 | §5.3 末段 `n1` 四行 | ✓ |

⇒ **它两张表里那 14 发变异，本程逐发复算相符**（RUN／PASS／FAIL／红名／红句原文五列都对得上，本机计数差已点名）；
**唯一不符处不在表里，在随码交付的那枚注释里**（§4.4）。
⇒ 加上它没报的 `y7`（本程补的）与 `x1p/asisp` 两发，**它 §4.2／§5.3 两张表的形状本程都独立走到了**——
**h1/h2 那两行在被验版本没有入库凭据（§0.3）这件事，因此不影响本程给它的账**：本程不是拿它的数对它的数。

### 4.2 问一（"摘掉任意一味，是否存在一发变异从此打不红？"）——**存在，本程逐味给发**

| 摘掉的那一味 | 从此打不红的变异 | 本程读数 |
|---|---|---|
| 摘掉 case 13（那条测试内 fake 断言） | `x1`（整条 `if s.exited()` 支删掉） | 21/13/**1** → **20/13/0** |
| 同上 | `z1`（exited 句丢掉 `last.summary()`） | 21/13/**1** → **20/13/0**（`z2`） |
| 摘掉 case 14 | `y1`（把借来的数放回去）／`y4`（只改值不改句子） | 21/13/**1** → **20/13/0**（`y2`） |
| 摘掉 `offsetUnknown` 本身 | 无解——编译期就断（`y5`＝`missing return`） | RUN **0**，属"构建状态"不是"读数" |
| 摘掉"新字段／新判定／新红句"（派单点名的三味） | **仓里没有这三味**（§1.5 表第 2 行） | 不适用，本程不虚构 |

⇒ **case 13 与 case 14 各自是其所钉承诺的唯一持有者**＝两味都**承重**；
但 case 13 承重的那条路**生产接线上今天不通**（第 1 格探针 A），case 14 承重的那条臂**今天没有任何真实文档走到**（§3.2 末行）。
⇒ 本程补一句它没说满的：**承重的强度不一样**。case 14 钉的是"借来的数不许印"，那是一枚**任何值都会被抓住**的判据
（`y1`=0 红、`y6`=1 红）；case 13 钉的是"死者不许被报成慢 writer"，而**今天没有死者能走到那条支**，
所以它的承重只在"有人手工收尸"的前提下成立。两枚都该留，但**只有 case 14 是"改了就会响"**。

### 4.3 问二（"摘掉它有没有任何外部可见读数变过？"）——**分两个口径，本程给它原来那句更正**

- **名册口径：变过。** `x3`／`y3` 摘掉用例 ⇒ `=== RUN` 21→20、顶层 PASS 14→13、名册差集少一枚。它答的就是这一层。
- **操作员口径：一毫米没变（本程现量，它没量过）。** 两发只差"那条支在不在"：

```
asisp-probe-only（支在）   ACC|A|sentence=wisp slo: subject 19080 never wrote a complete report within 2s (last read: 836 bytes read, document still open at offset 836: the tail had not arrived)
x1p-arm-deleted-probe-only（支被删）ACC|A|sentence=wisp slo: subject 17192 never wrote a complete report within 2s (last read: 836 bytes read, document still open at offset 836: the tail had not arrived)
```

⇒ 同一台机、同一枚探针、同一个真子进程形状，**两句逐字符相同**（只差 pid）。
⇒ 反向那一发也量了：`x1p` 里本程的探针 B（**有人收尸**那一形）从"印 exited 句"退化成"烧满 2s 印预算句"，
当场红（`--- FAIL: TestAcc152ProbeReapedTakesTheExitedLeg`，`elapsed=2.0001696s`）。
⇒ **所以**：那条支的全部可观测效果，今天只活在"测试自己收尸"这一个前提下。
它 §4.2 正文其实写了这句（"今天的生产输出，就等于'那条支被删掉'那一发的输出"），
**但那一格的标题答的是"有"**——两种口径在同一格里方向相反，本程按 §9 修正记录第 4 条记它一笔：
**措辞要分口径，否则下一位读者会拿"有"去判那条支今天活着。**

### 4.4 一处随码交付的注释把读数写错了（本程现量，方向已核）

```
$ grep -n -A3 "Case 14's own load-bearing reading" cmd/wisp/slo_report_144_windows_test.go
739:// g1-restored-borrowed-value-case14off.log reports === RUN=21, PASS=14, FAIL=0);
$ 本程数它自己那枚**入库**日志 probes/152/mut-post/g1-restored-borrowed-value-case14off.log
  topRUN=20  PASS=13  FAIL=0  SKIP=0   且名册里 TestSLO152 命中 False
$ 本程同发独立复算 y2-…-case14off（另一台基线、另一把尺）-> RUN=20 PASS=13 FAIL=0
```

⇒ 注释里那三枚数是 **`asis` 那一行的数**（21/14/0）被抄进了"摘掉 case 14"那一句——
**摘掉一枚用例，它自己不可能还在 `=== RUN` 里**。
⇒ **错的方向本程也定了性**：它把这发写得"更像没变异"（好像用例还在场、却零枚红），
比事实**更弱**地描述了那枚守卫的价值 ⇒ **不是假绿，是可复现性受损**：
下一位照注释去对数，会以为自己那发变异没落地（本程第一版 §4.1 就差点这么绕回来，见 §9 修正记录第 3 条）。
⇒ 它的**表**（§5.3 那行 20/13/0）与**它自己的日志**都对 ⇒ **修法只有一处：把那三行注释的读数改成 20/13/0**，
不改码、不改判据、不动用例。
⇒ 本格因此给这一票记两枚**同源同枚 commit（`10e3585`）的注释缺陷**：§2.2 末类那枚被删掉的 case 13 裁定行、与此处这枚写错的读数。

**第 4 格判定：承重两问＝成立（本程 19 发全复现，另补 2 发它没做的操作员口径对照）；它问二那格记一处口径更正；同一枚交付物里另记一处随码注释读数错（§4.4）。**

**本格抬头：第 4 格正文（§4.2 承重两问）由实现程写（盘上 mtime 10:01）、由编排者提交（`97cfc6e`，`162 加／0 删`）；
第 5 格正文同上。** 本程对这两格内容的裁定**不因代提而加重或减轻**——本程全部读数都是自己的尺重跑的，
没有一处引用它的日志当结论（引用它日志的地方都标了〔日志＋归档，抽验〕或"本程复算"）。
