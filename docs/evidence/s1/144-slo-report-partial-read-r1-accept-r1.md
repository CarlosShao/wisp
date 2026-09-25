# 144 — 对抗验收 r1（非实现者）：`collectReport` 分"没写完 / 坏了"这一味

- 验收程：本件作者，**不是** `95885fb` / `f63f0e3` / `a91d7c2` 的作者，进场前对本票无先验上下文。
- 被验锚点：**`1e94672`**（本程自己 `git cat-file -t 1e94672` = `commit`；`git diff --stat 95885fb 1e94672 -- cmd/wisp/slo_windows.go` 现量为空 ⇒ 收笔那句"本票零改动 `slo_windows.go`"**成立**）。
- 本程进场时刻仓库 `HEAD` = `b694378`（分支 `dev`，本程自己 `rev-parse` 量的）——**HEAD 已被别家程推走**，
  所以本件所有读数出自 `1e94672` / `95885fb^` 的**仓外纯净副本**，不出自工作树。
- 派单里编排者写的每一句（含"读数是成立的"那几句）在本件里一律按**未验证断言**处理；
  下面每一格给"现量命令 + 读数 + 判定"，判定档位只有 **成立 / 退回 / 退回（附条件入账、文件保留、不许勾）**。

## 0　副本与两件自证

### 0.1　纯净副本配方（两个 `-c` 都带，并自证生效）

```
$ cd "D:\work\workspace\projects plans\Wisp"
$ git -c core.autocrlf=false -c core.eol=lf archive 1e94672 | tar -x -C D:\tmp\wisp-144acc1\snap
$ cd /d/tmp/wisp-144acc1/snap && gofumpt -l . ; echo rc=$?
rc=0            ← 输出为空 ⇒ 副本没有被烙上 CRLF，两个 -c 确实都必要
```

建了四棵仓外树（全部在 `D:\tmp\wisp-144acc1\` 下，**只建不删**）：

| 树 | 来源 | 用途 |
|---|---|---|
| `snap/` | `1e94672` archive + 本程那枚验收驱动 | 修后形状，跑本程自己的差分 |
| `pre/` | `95885fb^` archive + 同一枚驱动 | 修前形状，跑**同一发**差分 |
| `prer/` | `95885fb^` archive，**无驱动** | AC#4 改前名册（`ls prer/cmd/wisp \| grep -c slo_report_144` = **0** ⇒ 真改前树） |
| `post/` | `1e94672` archive，**无驱动** | AC#4 改后名册 |

DLL 不在 git 里（现量 `git ls-files third_party` 为空）⇒ 四棵树各从工作树 `cp third_party/sherpa-onnx/*.dll`，
`PATH` 用 `/d/…` 形（前一程 §3.4 那发 127 的坑本程照它躲过了，未"顺手优化"那枚脚本）。

### 0.2　五枚 AC 框此刻的勾态与本程之前的表

```
$ grep -nE '^\- \[[ x]\] \*\*AC#' .scratch/wisp/issues/144-*.md
52:- [x] **AC#1 …    58:- [x] **AC#2 …    61:- [x] **AC#3 …    64:- [x] **AC#4 …    70:- [x] **AC#5 …
$ ls docs/evidence/s1/ | grep -i 144
144-slo-report-partial-read-r1.md          ← 只有实现方自证，无 accept 件
$ ls docs/evidence/s1/144-slo-report-partial-read-r1-accept-r1.md
No such file or directory
```

⇒ 编排者那句断言**复算相符**：五格全 `[x]`、非实现者的表在本程之前**不存在**。
本件就是那枚表；五格在下面**逐格出档**，不做整体判断。票面的勾本程**不动**（地界外）。

### 0.3　先量最承重的那一发：真修前树 vs 真修后树，同一枚真半截文件

前一程 §3.4 的那张差分表用的是 **667/668 字节**（它自己那枚夹具）。本程**不拿它的表当凭据**，
另写一枚只调**生产入口** `sloSubject.collectReport()` 的驱动（`acc144driver_windows_test.go`，
同文件同内容放进两棵树，仓内一字未改），自己造真半截文件：

**这份半截文件"是不是真的半截"——不是本程挑的形状，是主体自己的序列化器产的**：

```
ACC1 full doc via writeSLO = 1950 bytes; last byte = 0x7d (0x0a would mean trailing newline)
ACC1 0-byte file exists and is readable (0 bytes) => a reader CAN see a file before its bytes arrived
```

`writeSLO` 现量（`95885fb^` 树 `cmd/wisp/slo_windows.go:628-638`）= `json.MarshalIndent` 整块进内存 →
**一次** `os.WriteFile(out, data, 0o644)`，无临时文件、无 rename、尾字节 `0x7d`（不是换行）
⇒ 读者期间可见的**只有那份文档的前缀**，本程取 `full[:len(full)-1]` 就是**它能出现的最真的一形**。
前置自检：`bytes.HasPrefix(full, pre)` 为真、`json.Valid(pre)` 为**假**（否则这发不携带信息）。

**同一发在两棵树上的读数（逐字）**：

```
### pre/  (95885fb^)
ACC2 raw json.Unmarshal on the 1949-byte prefix = unexpected end of JSON input
ACC2 tree prefix: file=1949 bytes of a 1950 byte document
ACC2 collectReport error = wisp slo: subject report: unexpected end of JSON input
ACC2 elapsed=0s (budget constant … =33s; a first-read death returns far under it)

### post/ snap (1e94672)
ACC2 collectReport error = wisp slo: subject 4242 never wrote a complete report within 33s (last read: 1949 bytes read, document still open at offset 0: the tail had not arrived)
ACC2 elapsed=33.02s
```

⇒ §3.4 那张表的**句子形状**本程复算相符（字节数不同只因夹具不同），且**"修前第一发即死"是真的**
（`elapsed=0s` vs `33.02s`，两形用的是同一枚预算常量）。

⚠ 本程自己那发 `--- FAIL: TestACC2StaticPrefix` 是**本程驱动的断言写坏了**
（`elapsed > budget` 漏算最后一枚 `time.Sleep(poll)`，33.018s > 33s），**不是被验对象的红**；
它反过来是一枚正控：**预算常量确实没被挪**（否则那一发不会正好落在 33s 那一格上）。

**本程补的一发（前一程没取、且它才是"值不值"那一问的答案）—— ACC3**：
半截文件先在盘上，尾巴 **400ms 后**到达；只调生产入口，问"这份报告收不收得回来"：

```
pre/   ACC3 VERDICT=error     : wisp slo: subject report: unexpected end of JSON input
snap/  ACC3 VERDICT=collected  mode=subject-in-tree pass=true
```

⇒ **摘掉这一味，外部可见的不只是"句子变短"，是"一份真会到达的报告被丢掉、整发红"**。
这条是本票"承重"的操作定义下**唯一需要的外部读数**，见 §3（AC#3 那一格）。

### 0.4　本程撞到的、两程都没量的一件事：那一发竞态在本机**撞得出来**

前一程 §3.3 的结论是"无争抢时读者落不进那次写"，续程据此把外部读数判成"两版相同"。
本程拿同一份真报告做 reader-vs-writer 硬撞（写侧循环 `os.WriteFile` + `os.Remove`，读侧循环 `os.ReadFile`，
各 3 秒），并带**已知正控**（一枚手截 100 字节的文件确实可被 `os.ReadFile` 读到）：

```
pre/   ACC4 positive control: a hand-truncated 100-byte file IS readable … json.Unmarshal says unexpected end of JSON input
pre/   ACC4 partial read: *json.SyntaxError len=0 x998
pre/   ACC4 partial reads observed during 3s of concurrent writeSLO = 998
snap/  ACC4 partial read: *json.SyntaxError len=0 x871
snap/  ACC4 partial reads observed during 3s of concurrent writeSLO = 871
```

⇒ 三点，前两点是**新增读数**、不是复算：
① 那一形**在这台机器上可重放**（871–998 次/3 秒），所以票面"AC#1 自证作废"那一支**不必走**，
   前一程 §6 第 1 条"从没在真进程上把那一发竞态召出来"是**没做**、不是**做不到**；
② 观察到的**唯一**部分形状是 **`len=0`**（文件已建、字节还没灌），本程**没有**撞到 1..1949 的中间前缀
   ⇒ 生产里真正救回来的是 `io.EOF`-on-0-bytes 那一支（`reportUnwritten` 的 0 字节分支），
   `io.ErrUnexpectedEOF` 那一支在本机今天走不到；两支同归 `reportUnwritten`，所以修法仍然对得上，
   但**"半截文件"这个票面措辞与实测量到的形状不完全是一回事**，登记给归因；
③ 这两发**没拿"多少秒"当判据**：871/998 是**次数**，且两版数量同级 ⇒ 它证的是"这形造得出"，不是"哪版快"。

## 4　AC#4 门禁与名册 —— **成立**

本程**没有**采信前一程 §4.1/§4.6 的表，也没有"同树取两次"：改前用 `95885fb^` 的 archive 现造一棵
不含本票测试的树（`prer/`，现量 `ls prer/cmd/wisp | grep -c slo_report_144` = **0**），
改后用 `1e94672` 的 archive（`post/`，**不带本程那枚驱动**，避免把我自己的 4 枚用例混进名册）。
两发都只叫 `./cmd/wisp/`，**没有**一次把 `./internal/panel/` 放进同一次调用。

```
$ cd prer  && PATH=<tree>/third_party/sherpa-onnx:$PATH go test -count=1 -v ./cmd/wisp/   → rc=0
$ cd post  && PATH=<tree>/third_party/sherpa-onnx:$PATH go test -count=1 -v ./cmd/wisp/   → rc=0
$ grep -oE '\-\-\- (FAIL|PASS|SKIP): [A-Za-z0-9_/]+' <log> | sort -u > <t>.roster      # 票面指定的那把尺

prer:  RUN=116  --- PASS=116(含子格)  --- FAIL=0  --- SKIP=0  ^panic|^\[signal=0  roster=106
post:  RUN=130  --- PASS=130(含子格)  --- FAIL=0  --- SKIP=0  ^panic|^\[signal=0  roster=120
$ comm -23 prer.roster post.roster | wc -l   → 0        ← 消失 0 枚
$ comm -13 prer.roster post.roster           → 14 枚，逐枚 --- PASS: TestSLO144*
```

新增那 14 枚逐枚在册（顶格 7 ＋ `ReportsThatContradictThemselves…` 的子格 7），
**没有一枚 `SKIP`、没有一枚 `FAIL`、两遍 `^panic|^\[signal` 均 0**
⇒ 票面担心的"一枚用例 panic 吞掉同包几十条读数"这一形**今天没发生**，所以也没有"未取到"的债。

⚠ **判红绿只认 `--- FAIL:` 那一行**：`post-full.log` 里带 `^\s+<file>.go:<NN>:` 前缀的行有 **28** 枚（`t.Logf` 也带这个前缀），
而 `^ *--- FAIL` 是 **0** 枚 ⇒ 按"有没有 `file:line:`"判红会把 28 枚 Logf 读成红。两棵树的日志本程都按这一条判。

**CI 同形那发（`bash scripts/wisp-cli-tests.sh`，只 call `./cmd/wisp/`，在 `post/` 纯净树上）**：

```
runtests.sh: OK - packages=[./cmd/wisp/ -count=1 -skip ^(TestDefaultDeadline…|…|TestSyncRegistryProbeLive)$]
             top-level: PASS=70 FAIL=0 SKIP=0, === RUN=130, '[no tests to run]'=0
portable-tests.sh: four numbers (from -v): === RUN=130  --- PASS=70  --- FAIL=0  --- SKIP=0
```

⇒ 与 §4.4 自述的终态读数**逐字相符**（PASS=70／FAIL=0／RUN=130）。
本程**没有**在这一发里复现出续程报过的那枚并发假红
（`TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`，续程那一发是 `PASS=69 FAIL=1`）——
这一条与本票归因有关，单独在 §8 回答，不在这里记账。

**那把"0 命中"的尺先拿正控打过**（本仓当天两枚空核验的教训）：
- `grep -c TestSLO144 post.roster` = **14**（尺能命中本票新增）；
- `grep -c TestSecretOverwriteIsAnnounced prer.roster post.roster` = **1 / 1**（尺能命中既有用例，两册都在）
⇒ `comm -23` 读空是"确实没消失"，不是"尺不响"。

**门禁其余四发（全部在 `post/`＝`1e94672` 纯净树上现量，不采信自证件）**：

```
$ gofmt   -l cmd/wisp/                                → 空
$ gofumpt -l cmd/wisp/                                → 空
$ gofumpt -l . tools/d22scan tools/mockllm            → 空        ← CI 那枚全树尺（两把尺都跑，不是一把）
$ go vet ./cmd/wisp/                                  → rc=0 空
$ sh scripts/d22scan.sh                               → rc=0
    runtests.sh: OK - packages=[./...] top-level: PASS=30 FAIL=0 SKIP=0, === RUN=70   ← d22scan 自带正控
    d22scan: examined 228 production Go files under internal/ and cmd/
    d22scan: scope ban #8 cmd/  examined  43 Go files, comments and _test.go included  ← 本票新增那枚测试确实在射程里
    d22scan: clean - no D22 ban violations（真树 8 枚作用域 examined 全非零）
```

⇒ §4.6 那三行（改前 106／改后 120／消失 0／新增 14）**与本程现量逐格相符**；
四道门禁相符；唯一一处与本票叙述不符的数字差异是 `ban #8 design/`＝**30**（自证件写 32）——
原因是本程那棵树是 archive 出来的**已提交树**，而工作树里 `design/**` 有 16 枚别家程的未提交删除/移动，
**与本票无关，不记在本票账上**（按"引一枚红要连红句"的规矩，这一发根本不是红：rc=0、clean）。

**`slo_windows.go` 的改动面（本程现量，票面地界核销）**：
`git diff --name-status 95885fb^ 1e94672 -- cmd/wisp/` = `A slo_report_144_windows_test.go` ＋ `M slo_windows.go` 两枚，
**没有任何既有测试文件被改**（`| grep -v 'slo_windows.go\|slo_report_144'` 读空，
且该 grep 的正控＝去掉 `-v` 后它确实吐出那两枚路径）。
预算三常量 `subjectGrace=3s` / `subjectReportBudget=30s` / `subjectPollInterval=20ms`
在两棵树里**定义逐字相同**（`pre:74/81/82`、`post:77/84/85`，只是行号平移），
`git diff -U0 95885fb^ 1e94672 | grep '^[+-].*(subjectGrace|subjectReportBudget|subjectPollInterval)'`
只命中 4 行、**全是引用处不是定义处**（两 `-` 两 `+`，即"把两枚数挪进 `collectReportWithin` 实参"那一步）
⇒ **没有用放宽预算当修法**。

## 2　AC#2 修法只有一个方向 —— **成立**

票面这一格的四条约束，本程**逐条**用自己的现量判，不采信自证件的表格：

| AC#2 的约束 | 本程现量 | 判定 |
|---|---|---|
| 只许"没写完留预算里重读 / 坏了当场红"这一个方向 | §0.3 ACC2（修后 spending 33s 重读、修前一发即死）＋ ACC3（尾巴真到达时修后收到、修前丢掉）＋ ACC5 三态 tally | 成立 |
| 两者**分得开且有凭据** | `readSubjectReport` 的边界是 `errors.Is(err, io.EOF)/io.ErrUnexpectedEOF`（本程现量在 `post/cmd/wisp/slo_windows.go:625`，**类型化错误、不是字符串比对**）；凭据是 `subjectReportRead{bytes, offset}` 随最后一发一路带进错误句 | 成立 |
| 不许"解析失败就 `continue` 到预算尽头"而不留最后一发读数 | 见下面"AC#1/AC#3 变异复量"那一节的 N2 一发：本程把那一支真造出来量过 | 成立 |
| 不许加 `t.Skip`、不许动阈值/预算常量当修法 | `grep -c 't\.Skip' cmd/wisp/slo_report_144_windows_test.go` = **0**；那把尺的正控＝同包另外三枚文件里 `t.Skip` 命中 1/1/2 枚（尺能命中）；三枚预算常量在两棵树定义逐字相同（§4 末） | 成立 |
| fail-closed 不许放宽 | 最外层那句 `in-tree record unavailable (fail-closed, no silent downgrade to the tree basis)` 现量在 `post:394`，`git diff 95885fb^ 1e94672 -- cmd/wisp/slo_windows.go` 里**既无 `-` 也无 `+` 命中它**（唯一命中 `fail-closed` 的是新注释那一行 `+`）；`internal/observe/**`、`*thresholds.go`、`*golden*` 在本票区间 `git diff --name-only` **读空** | 成立 |

**"分得开"这句话本程不是读注释判的，是把两侧各问一遍判的**：
正控侧 = §0.3 的 1949 字节真前缀 ⇒ `unwritten`；
反控侧 = ACC5 顺手问的一发（`<html>…`，17 字节）⇒

```
ACC5 corrupt branch: state=corrupt offset=0 summary="17 bytes read, report corrupt:
   17 bytes contradict a subject report at offset 0: invalid character '<' looking for beginning of value"
```

⇒ 两态给的是**两个不同的句子、两个不同的动作**（一个回预算、一个当场红），
且红句点名"这些字节自相矛盾"，不是"再等等"。这一格本程没有异议。

⚠ 一处**本票自己已在册、本程复算相符**的软处：`subjectReportRead.summary()` 的 `unwritten` 那一支
写的 `…at offset %d…` 里的 `offset` **在 unwritten 分支上不携带信息**（详见 §6 那一格，本程用全 1978 枚前缀量过）。
它**不构成 AC#2 违背**（AC#2 要的凭据是"分得开"，AC#1ⓐ 要的是"预算＋字节数"，两样都在），
但它让 AC#5 新注释的最后那句（见下一格）说了一件实现没做到的事。

## 5　AC#5 那句注释 —— **成立（但新注释自己多带了一枚过期承诺，见 §6）**

```
$ grep -rn "never races" prer/cmd/wisp/      → :522 "// window, so this never races the observer's last sample."   hits=1
$ grep -rn "never races" snap/ --include=*.go → 0 命中                                                        hits=0
```

⇒ 被点名那句**已从代码里删掉**，且那把尺先拿"它确实在修前树里"打过正控（不是"正则恒不匹配"那种空核验）。

新措辞现量在 `post/cmd/wisp/slo_windows.go:663-681`（`collectReportWithin` 头上）。票面 AC#5 要的两件**都在**：
① **分得开的是哪两种状态** ⇒ "reportUnwritten stays inside the budget and is re-read.
reportCorrupt returns an error on the spot"；
② **为什么"文件存在"不等于"报告完整"** ⇒ "'the file is there' and 'the bytes are all there' are two different
states (os.WriteFile creates the file, then fills it)"，并且把原来那对错了的观察者换成了
"the race this loop has to handle is not observer-vs-observer, it is reader-vs-writer"。
⇒ **AC#5 成立**。

⚠ 但同一段注释的最后一句写着：*"if this gives up, the sentence names the budget it spent,
**the byte count of its last read and where inside those bytes the document stopped**"*。
本程量到的第三样**不成立**：`unwritten` 分支上"文档停在哪一字节"恒为 0（§6 那一格给全 1978 枚的分布）。

## 1　AC#1 钉成能响的用例 —— **成立**（7 枚全响、0 枚装饰、0 枚恒绿）

先确认"这一形造得出可重放用例"⇒ 不走票面那条"AC#1 自证作废"：本程 §0.4 已经在真进程＋真文件系统上
撞到 871–998 次部分读，且 §1 下面这 7 枚在**未变异的本票树上全绿**、在**单点回退下会响**。
基线（`snap/` = `1e94672` archive，挂 DLL）：

```
$ go test -count=1 -v -run TestSLO144 ./cmd/wisp/     → rc=0，7 枚顶格 PASS ＋ 7 枚子格 PASS，FAIL 0 SKIP 0
```

**恒真判据（本仓固定动作：单点回退）**——本程不只复算前一程的变异 A/B，另建两批共 13 枚单点变异，
每枚只改一处，问"撤掉这处，哪条用例变得不响？"。落点全在仓外 `snap/`，每枚跑完**二进制级还原**
（`shutil.copyfile` 回 pristine，`cmp -s` 现量 RESTORED-IDENTICAL，且 `gofumpt -l .` 仍空）：

| 变异（各只撤一处） | 红的用例（数） | 判定 |
|---|---|---|
| MA 两支 unwritten 全撤（＝修前判别形状） | 5 红／2 绿 | 与 §3.2 自述 `PASS=2 FAIL=5` **相符** |
| MA2 只撤 `io.ErrUnexpectedEOF` 那一支 | 5 红／2 绿 | 该一支独立承重 |
| N1 只撤 `io.EOF` 那一支（0 字节／纯空白） | 3 红／4 绿（`slo_report_144_windows_test.go:108 / :200 / :340`） | 该一支独立承重 |
| N2 吞掉 `case reportCorrupt: return`（AC#2 明禁那形，落到预算尽头） | 2 红／5 绿（`:250` 与 `:270`） | **与本程独立复算的 §3.2 变异 B 逐字相符** |
| MC 撤"值已闭合后面还有内容"那一判 | 1 红（`…ContradictThemselves…`） | 承重 |
| MD 撤 `rep.Report == nil` 那一判 | 2 红（`…ContradictThemselves…` ＋ `…CorruptReportIsJudgedOnTheFirstRead`） | 承重 |
| ME／N3 撤"到点红句里的字节数" | 2 红（`:230`／`:325` 点名 `44 bytes`／`8 bytes` 缺失） | **AC#1ⓐ 要的那两样里，字节数是被钉住的** |
| MF 撤"把最后一发读数带进放弃句"（`last = obs`） | 2 红 | 承重 |
| **MG 只撤 `offset` 那一半** | **0 红／7 绿，rc=0** | **不响——见 §6** |
| **N4 把 corrupt 那一支 summary 的字节数换成 err 里的重复** | **0 红／7 绿，rc=0** | **不响（但字节数仍由 `o.err` 带进句子，信息未丢）** |
| N5 **组合**：MG ＋ N1 同时 | 3 红（与 N1 同三枚） | **两枚各自半盲的洞合体没有放行真答案**（本仓有此先例，故专门造了一发） |

⇒ **AC#1 那 7 枚里没有一枚装饰**：每一枚都至少被上面一发行点名变红
（1↔MA/MA2/N1、2↔MC/MD、3↔MA/N1、4↔N2/ME/N3/MF、5↔N2/MD、6↔MA/ME/N3/MF、7↔MA/N1）。
⇒ **也没有一枚"在修前旧码上也响"**：判红绿只认 `^--- FAIL:` 那一行，两批共 26 发里 `NOT-RUN` 只在
电池 1 的 `continue` 版变异 B 出现过一次——**那是本程自己写坏的变异**（`continue` 会跳过状态 switch
**下方**的预算检查，所以那一形不是"重试到预算尽头"而是**永不结束**，`panic: test timed out after 10m0s`，
4 枚用例记**未取到**）。电池 2 的 N2 按票面原话重做（删掉那行 `return` 而不是插 `continue`），
一次跑通、**0 枚 NOT-RUN**，前一程 §3.2 的两枚红句（`:250` 的 `within 20s … 25 bytes …` 与 `:270` 的
`read #2 of a 1-reading script`）本程**逐字复现**。
⚠ 顺手留一条仪器事实给以后写派单的人：**`continue` 与"删掉那行 return"在这段循环里不是同一枚变异**，
前者永不结束、后者到点红；按票面字面写"改成 continue"会量出一个根本不等价的对象。

**"0 命中"这类读数本程都先拿已知正控打过那把尺**：`grep -c 't.Skip'` 在本票测试＝0、在同包另三枚文件＝1/1/2；
`never races` 在修前树＝1、修后＝0；名册尺见 §4；MG／N4 的"0 红"是**跑通了、7 枚都在册**的 0 红（`rc=0`、
`GREEN=7`），不是尺没响。

## 3　AC#3 变异自证两向 —— **成立**（含"值不值"那一问，本程给出比自证件更硬的答案）

票面 AC#3 三问，本程各给一发：

**(i) "主体写的报告确实坏了（非空、语义错）"⇒ 红且红因点名。** 本程**不采信**前一程"喂一串手写坏字节"那一形，
按票面原话走**主体自己的写入侧**：`writeSLO`（真序列化器）落真文件 → **生产入口** `collectReport()`（真预算、真
`os.ReadFile`）去收。本程 ACC1 现量 writeSLO 产的完整文档 = **1950 字节**（前一程那枚夹具 1978，两枚不同形状
同一条构造路径）；ACC5 反控那一发 `<html>` 17 字节 ⇒ `state=corrupt`、红句点名 `17 bytes … at offset 0:
invalid character '<'`。⇒ **成立**。

**(ii) 摘掉 AC#2 那一味 ⇒ AC#1 从绿变红。** 见 §1 那张表：MA（=判别侧的修前形状）红 5 枚。
**并且本程补了一发前一程没有的、按"承重"操作定义最硬的一发（§0.3 ACC3）**：
同一枚半截文件（1949/1950 真前缀）在盘上、尾巴 400ms 后到达，只调生产入口——

```
pre  (95885fb^) : ACC3 VERDICT=error     wisp slo: subject report: unexpected end of JSON input
snap (1e94672)  : ACC3 VERDICT=collected mode=subject-in-tree pass=true
```

⇒ **摘掉这一味，外部行为从"收到一份真会到达的报告"变成"丢掉它并红"**。这是绿↔红，不是句子好看。

**"另问那句"（摘掉它有没有任何外部读数变过）——本程的答案与自证件两程都不同，且比它们强：**
前一程 §3.3 答"外部读数一个字都没变"（两枚二进制各一发 `wisp slo` rc=0、JSON 里 `"observer"` 都在）；
续程 §3.4 改口"半截文件真被读到时必变"，并给了 667/668 字节的差分表。
本程复算：① §0.3 那一发**确实变**（红句整条换掉 + 从一发即死变成花完预算重读）；
② 更关键的是 ACC3 那一发——**外部可见的不是文本，是"这份报告收不收得回来"**；
③ 续程 §4.4 那两版 `state … pass=` 逐字相同，本程认为那是**那道门禁根本没造出前提条件**，
   而不是"这一味对外不可见"——理由见 §6 末与下面的 §8（那枚前置**见到任何一枚 `go.exe` 就拒采样**）。
⇒ 编排者那一问（"既然外部读数一样，这枚修法买到的到底是什么？"）的正确答案是：
**在无争抢的门禁路径上确实买不到外部读数差；在"半截文件真被读到"这条路径上买到的是行为差，不是措辞差。**
所以**不许**用"红句更好看"来答这一格，也**不许**用"两版 `pass=` 相同"来判这一味不承重。

**FAIL/SKIP 面**：两批变异里 `--- SKIP` 计数 **0**；`NOT-RUN` 只出现在本程自己写坏的那一枚变异（已按"未取到"记，
没有改成 Skip、没有拿它当绿）。

## 6　那一格"没顺手改"的新发现（`offset` 恒 0）—— **读数成立；"不改"这个处置不成立，退回**

分两问判，本程分开出档。

**读数那一问：真。** 续程说"对全部 1978 枚前缀全量量过、恒为 0、只有完整文档报 1978"。
本程**不采信它那张表**，另写一枚探针（`acc144probe_windows_test.go`，只落仓外 `snap/`），
用 `slo144Report` 那枚**夹具本体**（不是本程另造的文档）逐枚问 `0..len` 全量，并给这把尺装了正控：

```
ACC5 fixture length = 1978 bytes
ACC5 state tally over 1979 readings (0..len inclusive): map[complete:1 unwritten:1978]
ACC5 distinct offset values = 2 -> [0 1978] (counts map[0:1978 1978:1])
ACC5 state="unwritten" offset values=[0]
ACC5 state="complete" offset values=[1978]
ACC5 positive control: a 44-byte prefix gives state=unwritten offset=0 bytes=44
```

⇒ 三件事同时被钉住：① 夹具确实 1978 字节、② 全 1978 枚前缀确实全判 `unwritten`（前一程 §1 那格的**分母**成立）、
③ `offset` 在 `unwritten` 分支上确实只有一个取值 0。
**正控不空转**：这把尺能报出 `1978`（完整文档那一发），所以"恒 0"不是"尺坏了我没发现"。
顺带本程量到 `corrupt` 分支也报 0——**`offset` 这一味只在 `complete` 时携带信息**。

**"不改"那一问：理由不成立。** 续程给的理由是
*"改错误文本会让 §3.2/§3.4 已落盘的变异红句逐字变样，而那批红句正是变异自证的凭据"*。本程造出发得出这个结论的尺：

- **MG：把 `unwritten` 那一支 summary 里的 `at offset %d` 整段拿掉 ⇒ `rc=0`、7 枚用例全绿。**
  ⇒ **本票那 14 枚名册里没有任何一枚钉住 `offset`**。改它**不改变任何一条断言的颜色**，
  "会弄乱凭据"只对**证据件里的文字转录**成立，而转录是**可复跑再生成**的（本程两批变异就是再生成了一遍）。
- 对照：**ME／N3（拿掉字节数）⇒ 2 枚红**、**MF（不带最后一发）⇒ 2 枚红**
  ⇒ 真正被钉住的凭据（预算＋字节数）与本格无关。**被"少了一半"的恰好是没钉的那一半。**
- 后果不只是审美：§5 已量明 `collectReportWithin` 头上那枚**新注释**逐字承诺红句会点名
  "where inside those bytes the document stopped"。**恒 0 ⇒ 那半句承诺今天不成立**，
  而这枚票的起因就是"注释把一件实现没保证的事写成不变式"。
  对一枚 1949 字节的真前缀，那句红今天读起来像"文档停在第 0 字节"（§0.3 本程那一发的原文就有这一串），
  **这会把读红句的人往"根本没开始写"那个方向带**，而真答案是"停在第 1949 字节"。

⇒ **判定：读数成立；`不改` 这个处置退回。** 最小闭合集合（**不在本票内做**，本程一字未改码）：
① 把 `unwritten`／`corrupt` 两支的 `offset` 要么改成携带真值（`dec.InputOffset()` 在报错那一发上给 0 是 Go 解码器的行为，
   要真值得自己算，例如取"已消费字节数"），要么**从红句里删掉那一半**（MG 已量明删了不掉任何覆盖）；
② 同步把 `:578-580` 那枚字段注释与 `collectReportWithin` 头上那句"where inside those bytes the document stopped"
   改成实现真给的东西；③ 改完**重跑两批变异、再生成红句转录**，别保留旧转录当凭据。
④ 归哪一格请编排者裁：它**不是 AC#1／AC#2／AC#5 的违背**（三格各自的要求本程都量明满足了），
   它是"新注释自己刚写下一句实现没做的承诺"——**同一族缺陷的第二发**，
   本程倾向**记台账另开票**（与本票 5 格的勾无关），但**不接受**用"改了会弄乱凭据"把它留在原处。

**另外两件与这一格同批登记、本程判"不算缺陷"的**：
- `subjectReportBudget+subjectGrace` 在红句里以 `33s` 出现——本程现量 `:77 = 3 * time.Second`、`:84 = 30 * time.Second`，
  两枚定义在修前/修后**逐字相同**（§4 末），§0.3 那一发实测确实花到 `elapsed > 33s` 才放弃
  ⇒ **`33s` 是被花掉的预算常量，不是计时读数**，写手这一条自证相符，不记账。
- 续程 §4.4 那句"那道前置筛的是机器不是这一味"——**本程复算成立，而且理由更强**：
  `scripts/slo-check.ps1` 的采样有效性前置现量是**两味**，都不看本票这一味——
  ① `:153-156` 那 15 枚进程名（`go.exe / gofmt.exe / cgo.exe / compile.exe / asm.exe / link.exe / gcc / g++ /
  cc1 / cc1plus / as / ld / wisp.exe / wisp-cli.exe / staticcheck.exe`）里**任何一枚**不在本进程树内出现，
  就记一枚 offender（`:179` 的 reason 逐字是 `foreign toolchain/wisp process present`），
  外加 `:160-168` 那条"任何跑在 runner work root 下的进程"；
  ② `:211` `$cpuBusyPct = 50` ⇒ 整机 CPU 两发 1s 采样的最大值 ≥50% 也拒（`:226`）。
  任一成立就走 `:250-256` 的 `NO CONCLUSION (machine-contended)`（owner 在 Q-36 定案"拒采样不是红"，
  注释 `:238-241` 逐字写着 `What is NOT changed here … WHEN we refuse` 的三条永不放宽）。
  ⇒ **只要这台机器上还有第二枚程在跑 `go.exe`，这道前置就不产出任何 state 行**，
  两版对照因此**结构上取不到"这一味的外部读数"**，与谁对谁错无关。
  ⚠ 本程自己那一发实测见 §8 末（含它到底拒在哪一味）。
  ⇒ 建议从今往后对这枚仪器的派单措辞改成：**"两版 `state … pass=` 相同不构成任何一格的外证，
  它最多构成"这道门没被本票改动"的外证"**——续程那句"筛的是机器不是这一味"是对的，
  但它**不该**被写成"取到了两版对照"，那一发的信息量是 0。
⇒ 这是**同一枚票刚修掉的那类毛病的自家用版**——注释写了一件实现没保证的事。
它不推翻 AC#5（AC#5 只要求"改掉 never races ＋说明那两件事"，两条都满足），
但它把 §6 那一格里"要不要修 offset"从" cosmetic 取舍"抬成"注释与实现不一致"，
**本程因此不采纳"改了只会弄乱凭据"作为不改的充分理由**（理由本身合不合规矩，见 §6 的判定）。

## 7　常规必查两轴 ＋ 一枚与本票叙述不符的 commit 形状

**轴① 已有断言有没有被动过方向**（只许加、不许减、不许把极性改成永真）：
本程不读自证件的自述，直接用树对树：
`git diff --name-status 95885fb^ 1e94672 -- cmd/wisp/` 只有 `A`（新增那枚测试）＋ `M slo_windows.go` 两枚
⇒ **没有任何既有测试文件被改**，既有的 `if rep.Report == nil` 那一腿在 `post:644` 仍在（本程现量、
且 MD 那一发变异证明它今天仍在响）。`thresholds.go`／golden／`internal/observe/**` 本票区间零改动（§2 表末行）。
⇒ **轴①：没有放水**。

**轴② 用的 helper 是不是原有那枚**：夹具走**主体自己那枚序列化器**
（`json.MarshalIndent` 的入参是 `sloRun`+真 `observe.StateReport`，与 `writeSLO` 同一调用；本程 ACC1 用真
`writeSLO` 落盘又量了一遍同一形状，两者形状一致）；循环那一侧 `TestSLO144LoopGiveUpSentencesOnRealFiles`
用**真文件＋真 `os.ReadFile`**（走 `collectReportWithin`，只是把预算收进参数），
`…CorruptReportIsJudgedOnTheFirstRead` 用**新增的 `readReportFile` 接缝**——
本程核过那一枚接缝**在生产里是 `nil`**（`post:686-688` `if read == nil { read = os.ReadFile }`）
⇒ 不是"用 mock 代替真的"那一形。**轴②：过。**

**⚠ 与本票叙述不符的一处（这是本程造出来的洞，不是怀疑）：`a91d7c2` 不是"只有空白"那一枚 commit。**

```
$ git diff -w a91d7c2^ a91d7c2 -- cmd/wisp/slo_report_144_windows_test.go
   → 非空：新增 6 行注释 ＋ 把夹具从 2 字段扩到含 Samples/Verdicts/SubjectPID 等，
     并把原来那一发 `if !strings.Contains(doc, `"report"`)` 换成
     `len(doc) < 300` 前置 + 5 枚字面量循环
$ git diff a91d7c2^ a91d7c2 --numstat     → 42 插入 / 13 删除
$ <两版各 tr -d 所有空白后 cmp>           → differ: byte 1406   ← 剥掉全部空白也不相同
```

票面（Progress log 续程段）与证据件 §4.2 都写着 **"改动只有空白（前后两版 `diff -w` 逐字相同）"**，
commit `a91d7c2` 的标题句是 **"复量与票面不符，只改对齐"**。**两句都不成立。**
真实形状：`a91d7c2` 把**前一程未提交的夹具升级（+36/−7）与 `gofumpt -w` 的空白改动混进了同一枚 commit**，
所以"只有空白"是对**半枚**改动的描述，被当成了对**整枚**改动的描述。

**实质影响本程判过：没有放水。** 那处替换把一发断言换成**更严**的三件
（长度下限 ＋ 5 枚字面量，原来那枚 `"report"` 仍在集合里），极性没变、没减；
且 §0.3/§4 里"1978 字节"这个分母**正是这次升级的产物**——
本程在 `1e94672` 树上现量夹具 = **1978 字节**（ACC5），与全件所有引用相符。
**被记进的账**：一处"自证凭据是空白-only"的陈述为假 ⇒
它同时是 §4.2 那三行"gofmt/gofumpt 全空"读数的**适用范围问题**（见下），
所以**AC#4 那一格本程按"成立"出，但把这一处按缺陷退回给编排者**（缺陷在证据件与票面的措辞，不在码）。

顺带把 §4.2 那三行的适用范围钉清楚（本程现量，且是**改后现树**）：
`gofmt -l cmd/wisp/` 空、`gofumpt -l cmd/wisp/` 空、`gofumpt -l . tools/d22scan tools/mockllm` 空——
**三发都在 `1e94672` 的 archive 树上取的**，不是"升级夹具之前"。
前一程那三行读数确实过期（续程自己也这么写），但**过期已被现量补上**，这一处不再挂账。
