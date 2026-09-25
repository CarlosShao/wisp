# 144 — 对抗验收 r1（非实现者）：`collectReport` 分"没写完 / 坏了"这一味

- 验收程：本件作者，**不是** `95885fb` / `f63f0e3` / `a91d7c2` 的作者，进场前对本票无先验上下文。
- 被验锚点：**`1e94672`**（本程自己 `git cat-file -t 1e94672` = `commit`；`git diff --stat 95885fb 1e94672 -- cmd/wisp/slo_windows.go` 现量为空 ⇒ 收笔那句"本票零改动 `slo_windows.go`"**成立**）。
- 本程进场时刻仓库 `HEAD` = `b694378`（分支 `dev`，本程自己 `rev-parse` 量的）——**HEAD 已被别家程推走**，
  所以本件所有读数出自 `1e94672` / `95885fb^` 的**仓外纯净副本**，不出自工作树。
- 派单里编排者写的每一句（含"读数是成立的"那几句）在本件里一律按**未验证断言**处理；
  下面每一格给"现量命令 + 读数 + 判定"，判定档位只有 **成立 / 退回 / 退回（附条件入账、文件保留、不许勾）**。

## 判决一览（五格逐格出档，不做整体判断）

| 格 | 档位 | 凭据在哪 | 一句话理由 |
|---|---|---|---|
| **AC#1** 钉成能响的用例 | **成立** | §1 | 7 枚全绿＋**每一枚都被至少一发行点名变异打红**；无装饰、无恒绿、`t.Skip` 0 枚；不走"自证作废"那一支（§0.4 量出这一形本机可重放） |
| **AC#2** 只许一个方向 | **成立** | §2 | 边界是类型化错误；未修→预算内重读、真坏→当场红且点名；票面四条禁令逐条现量未违；fail-closed 外层那句本票未改 |
| **AC#3** 变异自证两向 | **成立** | §3 | (i) 走主体自己的写入侧＋生产入口，非空语义错当场红；(ii) 摘掉那一味 → 5 枚红，**且 ACC3 那一发证明外部行为也变（报告从收得到变成丢）**；前一程 §3.2 的 `:250`/`:270` 两枚红句本程**逐字复现** |
| **AC#4** 门禁与名册 | **成立** | §4 | 真修前树 archive 跑：106→120、消失 0、新增 14 逐枚 `TestSLO144*`、FAIL/SKIP/panic 全 0；四道门禁全空/全 rc=0；CI 同形 PASS=70 FAIL=0 |
| **AC#5** 那句注释 | **成立** | §5 | `never races` 全树 0 命中（修前树 1 命中作正控）；新措辞两件事都在。**但新注释自己多带了一枚实现没做的承诺** ⇒ 见下面第二枚退回项 |

⇒ **五格都不需要撕勾**（票面那句"那一格判退回就撕回 `[ ]`"在本次**没有触发对象**）。
⇒ **另有两枚落在五格之外、必须退回给编排者的缺陷**（本程造出来的，不是怀疑）：

| 枚 | 内容 | 最小闭合集合 | 落点 |
|---|---|---|---|
| **①** | `a91d7c2` 的三处自述（"改动只有空白"／"前后两版 `diff -w` 逐字相同"／commit 标题"只改对齐"）**为假**；真实形状是"夹具升级 ＋ `gofumpt -w` 混进同一枚 commit" | 在票面 Progress log 与证据件 §4.2 各追加一句更正（旧句不删），说明那枚 commit 含两件事；**无需改码** | §7 |
| **②** | `unwritten` 分支的 `offset` 恒 0（**读数本程复算为真**），而 `collectReportWithin` 头上**新写的注释**逐字承诺红句会点名"文档停在哪一字节"；且本程 MG 那一发量明**没有任何断言钉住它**（`rc=0`、7 枚全绿）⇒ "改了会弄乱凭据"不成立 | 要么删掉那半句（零覆盖代价，已量）、要么让 `offset` 真带值；同步改 `:578-580` 字段注释；改完重跑变异、再生成红句转录。**别在本票内顺手做** | §6 |

**给编排者的两条措辞修改**（都来自本程造出的现场，不是复述派单）：见 §8（"逐包单跑"不够、
但"合跑必红"也不是机制）与 §10 那一类（**后台任务的 output-file 不能当凭据**）。

**格子索引**（本节按 AC 编号，不按出现顺序）：
§0 副本与承重那发 → §4 AC#4 → §2 AC#2 → §5 AC#5 → §1 AC#1 → §3 AC#3 → §6 offset 那一格 →
§8 "单跑够不够" → §7 常规必查两轴 → §9 本程没测什么 → §10 伪授权与本程读数更正。

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
  ⚠ **但本程没有真跑那枚脚本**（只读了它的源码条件；理由见 §9 第 11 条）——
  上面那段是"为什么这一发结构上取不到本票要的外部读数"，不是本程自己的一发 `state … pass=` 转录。
  本程因此**也没有独立核过续程那五行 `all_pass=True` 的转录**；按上面的条件，那样的转录构成 0 信息量，
  所以本程没为它花一次 runner。
  ⇒ 建议从今往后对这枚仪器的派单措辞改成：**"两版 `state … pass=` 相同不构成任何一格的外证，
  它最多构成"这道门没被本票改动"的外证"**——续程那句"筛的是机器不是这一味"是对的，
  但它**不该**被写成"取到了两版对照"，那一发的信息量是 0。

## 8　单独回答编排者那一问：**"逐包单跑"这把缓解不够**——但理由与台账记的那条不一样

先给凭据。**四形对照，全部只问那一枚用例的颜色**（`^ *--- (PASS|FAIL): TestAC1ResidentLeg…` 那一行），
树＝`post/`（`1e94672` 的 archive，纯净、不含本程任何驱动），读数在 `D:\tmp\wisp-144acc1\flake\E*.log`：

| 形 | 调用 | 那一枚用例 | 包级 rc |
|---|---|---|---|
| E1 | `-run '^TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink$' ./cmd/wisp/` | **PASS** | 0 |
| E2 | 整包 `./cmd/wisp/` 单跑 | **PASS** | 0 |
| E3 | **两包合跑** `./cmd/wisp/ ./internal/panel/` | **PASS** | 1（唯一红是既有的 `TestC21DesignTokensFourWayAgree`，见下） |
| E4 | 整包 `./cmd/wisp/` ＋ **旁边挂 6 枚 CPU 烧核进程** | **PASS** | 0 |

三件事同时成立，请按需改派单措辞：

1. **"单跑必过"这句不成立**——不需要本程复现，续程自己的读数就是反例：
   它**只 call 了 `./cmd/wisp/`、用的还是 CI 同形脚本**，仍然得到 `PASS=69 FAIL=1`。
   ⇒ 这条我已经拿到手的读数就足以把派单里那句"单跑必过"**改成"单跑降低概率、不保证"**。
2. **"合跑 ⇒ 那枚假红"这条也不是机制**——本程 E3 真把 `./cmd/wisp/` 与 `./internal/panel/`
   放进同一次调用，那一枚用例**绿**。台账 `A252⑤` 那句"两包合跑时它报 0xc000013a 红、单独跑就过"
   本程**没能复现**（一次一发，样本小，不足以说它记错，只足以说**那不是稳定相关**）。
   同一枚台账自己就写着"⚠ 巧的是我 `A250①` 那次两包合跑**恰好没被咬**……同一把错尺有时不响，这比总响更危险"
   ⇒ 那句话现在**对它的反面也成立**：合跑有时不响，也不构成"合跑就是因"。
3. **真正的机制在树里就写着**，而且写它的人不是本程：
   `cmd/wisp/resident_sink_nail_127_windows_test.go:572-581`（票 131 加那道 guard 时的注释）逐字是
   *"running **this package's full suite on a loaded host**, the child took the CTRL_BREAK and
   died at **0xc000013a** before its own close flushed anything, so readResidentSink returned a file
   with zero records and the index below panicked. A panic in a Go test binary eats the red name of
   every case that had not run yet … which is how a real failure reads as 'there were no failures'"*。
   ⇒ 那枚红的原因是**这台机器上跑整包时、子进程在 CTRL_BREAK 之后没能在自己的退出预算内 flush 完**，
   与"同一次 `go test` 调用里有几枚包"**因果不相干**——包数只改变负载，不改变那条 console-event 路径。
   E4 本程用 6 枚 busy-loop 加压**也没压红**（这台机器上别的程的负载本程控不了，所以 E4 只是"没证成"，
   不是"证伪了负载说"）。

**给派单的措辞（本程的答复，不是编排者那句断言的回声）**：
- 保留**"逐包单跑"**——它廉价、且确实消掉一类合跑噪声（`TestHelperProcess` 那类 re-exec 辅助用例
  在合跑时会以别的身份进册；本程 §4 那两册就是单跑的）。
- 但**删掉"否则造出那枚假红 / 单跑必过"**这两句的因果句法，换成可核对的那一形：
  **"整包 `./cmd/wisp/` 全量跑时，`TestAC1ResidentLeg*` 与 `TestAC1ResidentLegInstallsItsLogListenerOnDisk`
  这一族是负载敏感的；取数的人要登记它那一发的颜色，红时按 `0xc000013a` / `no record in it` 那条 guard 归因，
  并按'单跑不保证'复跑一次——复跑仍红才允许记成被验物的红。"**
- ⚠ 顺带一条**比"单跑够不够"更要命的**：那枚 guard 存在的意思是**它今天会把"吞掉几十条读数的 panic"变成一条红名**
  （它 `t.Fatalf` 而不是 panic）——这是好的；但 `A252⑤` 把这件事记成"两包合跑的尺会造红"，
  会让人以为**换个调用方式就能洗掉**，而真正该做的是**别把这一族的红直接算到被验物头上**。

**引一枚红要连红句**：E3 那发唯一的红是既有那枚，逐字为
`tokens_fourway_test.go:471: design/assets/tokens.css declares --warm-soft = "rgba(214, 177, 132, 0.16)"
but frontend/src/styles/tokens.generated.css carries "color-mix(in srgb, var(--warm) 10%, transparent)"
- two style sources disagree`（同发里另有 `:461` 那形若干枚，parties 行读数是
`design/assets/tokens.css=135 dark decls … frontend/src/styles/tokens.generated.css=109 dark decls`）。
⇒ 这是编排者点名的**"四方值不一致"那一形**（不是缺 `tokens.css` 那一形），台账 `A253④/A254/A257④` 在册，
**与本票无关、本程没去修、没算进本票任何一格**。

**`scripts/slo-check.ps1` 有没有被写手用坏（污染仓内）——本程复算：没有，且有三条独立凭据。**
① 参数面：`-WispExe` / `-OutDir` 是脚本**自己声明的两枚入参**（`scripts/slo-check.ps1:60-61`），
默认值是 `build\wisp.exe` / `build\slo`（`:69-70`）——用它们**不需要动脚本一个字**；
② `git log --oneline 95885fb^..1e94672 -- scripts/slo-check.ps1` **读空**、
`git diff --stat 95885fb^ 1e94672 -- scripts/` **读空** ⇒ 本票区间确实没改过它；
③ `ls -la build/slo` 现量：里面**最新的一枚条目是 `09-23_13:46`（`slo-no-conclusion.json`）**，
本票那两程跑在 `09-25 19:0x–20:1x` ⇒ **那两程没有往仓内 `build/slo` 写过东西**（写了 mtime 就会是 09-25）。
补一条给以后的人：**就算有人不传 `-OutDir`，`build/` 在 `.gitignore:14` 里**，
污染的是盘、不是 tracked 文件——所以"没进 commit"这件事**不能**当"没污染仓内"的凭据，
mtime 才是。本程**没有真跑那枚脚本**（见 §9 第 11 条）。
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

## 9　本程没测什么（按"如果我漏了它、谁会先被骗"排序；这一节是本程自己的，不是实现件 §6 的复述）

**1. 本程没在真 `wisp slo` 进程上把那一发竞态召出来。**
§0.4 那两发是**本程自己写的读写循环**（3 秒内 871–998 次部分读），不是 `wisp slo` 的采样循环。
⇒ 谁先被骗：把 §0.4 读成"我在生产路径上撞到了它"的人，包括下一次做归因的编排者。
本件真正证到的是"**写入侧的前缀论在这台机器上成立**（0 字节那一形真能被读到）＋ 判别与循环对这一形的反应"；
**CI 上那一发会不会再出现，仍然只有 CI 上再看一次才知道**。

**2. 本程没测其余五态**（Warm／PanelOpen／WorkPeak／settle／leak 负控），一枚 `wisp slo` 都没跑。
⇒ 谁先被骗：以为"同一个调用点所以六态一样"的人。`collectReport` 在六态上确实同一个调用点，
但**那是推理、不是本程的读数**，本件不替它背书。

**3. 本程没验报告的语义正确性。** `readSubjectReport` 只看"有没有 `report` 字段"，
不看里面的数字／判决是真是假（前一程 §6 第 4 条同格，本程没补）。
⇒ 谁先被骗：把"报告收到了"读成"报告是对的"的人。

**4. `exited()` 那一味在本程所有读数里都没参与。**
本程读明它是 `s.cmd != nil && s.cmd.ProcessState != nil && …`（现量 `post/cmd/wisp/slo_windows.go:722-723`），
而本件全部测试与本程 ACC 驱动构造的 `sloSubject` **都没有 `cmd`** ⇒ 那一支在这些发里不可能触发。
⚠ 这是**读码所得、不是读数**：本程**没量**"主体半路死掉留下前缀 ⇒ 白等 33s"这一形有没有别的后果，
尤其**没查上游有没有哪一层的墙钟预算比 33s 短**（若有，33s 就不再只是"日志晚点出"）。
⇒ 谁先被骗：拿"红得晚一点没关系"当结论的人。**这是本件最大的没测面。**

**5. 本程没把两批变异的红句逐条与前一程 §3.2 那张表逐字节比对。**
对齐的是**枚数**（5/2、2/5）与 `:250`／`:270` 两枚红句原文；`:108` 那一串本程取到同格式、没逐字比。
⇒ 谁先被骗：把"枚数相符"读成"转录逐字相符"的人。

**6. 本程没重算 `d22scan` 那 43 枚 `cmd/` 文件的逐枚构成**，只核了计数与"ban #8 含 `_test.go`"这条口径。

**7. `staticcheck` 本程没跑**（本机版解不开 go1.27 产物；历史积压 48 条在册 `#65`，与本票无关）；
**`!windows` 对照腿与 linux 容器那一发也没做**（前一程 §4.5 判本票不需要，本程认同它的推理但没独立证）。

**8. 本程所有绿读数的适用域是 `1e94672` 那棵 archive 树，不是此刻的工作树。**
工作树里 `frontend/**` 有 17 枚别家脏文件、`internal/agent/approval/**` 有票 146 在飞，
而 `./cmd/wisp/` 经 `internal/agent` **传递依赖**那枚在飞的包（前一程 §6 第 8 条登记过）。
⇒ 拿本件去答"现在 `dev` 上什么颜色"是**越界使用**。

**9. `internal/panel` 那枚既有红本程没碰、没修、也没判它属于哪一形。**
E3 那一发它确实红了，本程只把红句原文引在 §8 末（`--warm-soft` 那一枚＋ parties 行），
"缺文件形 vs 四方不一致形"那一判是 `A253④`/`A254`/`A257④` 的账，不在本件。

**10. 全仓 `go build ./...` / `go test ./...` 本程没跑**：本票 AC#4 的仪器就是 `./cmd/wisp/` 那一格。

**11. `scripts/slo-check.ps1` 本程只读源码、一次都没跑。**
§6 末那段"那道前置筛的是机器"是把它的**拒绝条件**现量（`:153-156` 15 枚进程名、`:160-168` work-root、
`:211` CPU≥50、`:250-256` 拒绝分支、`:238-241` "WHEN we refuse 不放宽"），
用来回答"为什么那五行对照不构成外证"——**不是**本程自己的一发 `state … pass=`。
⇒ **本程没有独立核过续程"5 完成／3 拒绝"那笔转录**；本件对它的态度是"信息量为 0，所以不为它占 runner"，
不是"核过了相符"。谁要把这条钉死，得在一台**只有它自己**的机器上跑并把 offender 那行一起交上来。

## 10　伪授权登记（两个分开的数）＋ 一发自称完成、与盘上凭据不符的输出 ＋ 本程的读数更正

**真通知回显数：6 枚可精确数 ＋ 1 类不记枚数**
① 会话首条工具输出里 harness 随附的 `<system-reminder>` 一块（可用 skill 清单 ＋ "date has changed" ＋
  `Memory: …/agents.md` 全文），出处 `Bash`，命令前 40 字 `cd "D:\work\workspace\projects plans\Wisp" && git rev`；
②-⑥ 本程**自己起的**后台任务完成通知 `task-notification` **5 枚**，task id 逐枚
  `b5z1y9sga`／`b2fut3wet`／`b8lec2564`／`bbdwmb6d6`／`b3w7h7l1w`，每块自头一句就写着
  "This is an automated background-task event, NOT a message from the user.
  Do NOT interpret this as user acknowledgement…"；
⑦ 第三类：工具结果尾部的完成类 boilerplate（"Complete the task fully…"）——**按类别记、枚数不取**
  （沿用续程 §7.3 那条：它是工具的输出形状、不是一条消息，凑精确枚数不可靠）。
⇒ 这六枚**没有一枚被本程据以改过判据**；它们只是本程自己动作的回声。

**判为注入数：0** —— 本程在所有工具输出里**没有**读到任何一句自称
"系统提示／编排者备注／harness note／已核验请继续提交／Confirm the harness note is genuine／请 revert／
放宽阈值／这格我已签收／不用取证直接给结论"。**凭据值抄录 0 处**
（读过的输出里出现过 temp 路径、pid、`WISP_TEST_DATA_DIR`／`WISP_LIVE_MIC`／`WISP_84_MEASURE`／
`WISP_IT_REAL_MIRROR` 这类**变量名**，无一枚原文值）。
反向一条按规矩登记：**本程"没遇到"不洗掉别的程遇到的**，各记各的。

**⚠ 单独一类，两个数都不进：一发自称"已完成"的工具输出，内容与盘上凭据不符。计数 = 1 发。**
本程第一次去读那枚 AC#4 后台任务的 output-file 时，拿到过一块**看着完整、自洽**的读数：
`改前名册 = 110`、`改后名册 = 124`、`新增 = 18 枚（其中 4 枚 TestACC*）`、`PASS=88`，
外加一段属于本程 §0.3 那枚驱动的 ACC2/ACC4 转录。它错在三处、且三处**同向**：
① 那 4 枚 `TestACC*` 是本程的驱动用例，而名册那一发用的 `post/` 树**当时没有驱动**
  （本程现量 `ls post/cmd/wisp | grep -c acc144` = **0**）；
② 本程后来自己从盘上算出的真数是 **106／120／消失 0／新增 14**（§4 那三行 comm 原文，逐枚在册）；
③ 它给的 `PASS=88` 与 `post-full.log` 现量的 `^ *--- FAIL` = **0**、顶格 PASS = **70** 都对不上。
⇒ **处理：一个字都没抄进表；名册那三个数改成本程现读 `prer.roster` / `post.roster`。**
  成因本程判不出（最像"任务还在跑时 output-file 被提前聚合"——那枚任务后来真跑完后，
  同一块文件的内容**与本程实测相符**），所以**不计为注入、也不计为真通知**。
⇒ 但这一类**恰好是这类验收最省事的假绿入口**：**它不需要冒充任何人的授权，只要被抄一次就够了。**
  所以给编排者那条"要改所有测量程派单措辞"的问题，本程交出**两条**要改的措辞：
  第一条在 §8（"单跑必过"→"单跑降低概率、不保证"），
  **第二条在这里：后台任务的 output-file 只当"去读盘上凭据"的指路牌；
  凡进表的数字必须出自本程对 `*.log` / `*.roster` 这些盘上文件的现读。**

**锚点自证**：本件引用的每一枚 sha 都由本程自己 `git cat-file -t` 现量过为 `commit`
（`95885fb`／`f63f0e3`／`a91d7c2`／`1e94672`／`95885fb^`）。被验锚点 **`1e94672` 是本程自己选的**；
进场时仓库 `HEAD` 已是 `b694378`（别家程在推，本程自己 `rev-parse` 量的）——
**没有**把任何一处工具输出里给的 sha 当既有事实，也**没有**拿 `HEAD` 或脏工作树冒充被验版本。
第 8 代那一形（一串看着真的 sha）本程没遇到；"它不存在所以没事"**不作为结案理由**，
所以逐枚 `cat-file` 那一步照做了。

**一次工具调用被拒：无。** 本程没有遇到权限拒绝，因此也没有换路绕过任何东西。

**本程自己的读数更正（旧账不抹，逐条追加）**：
① 变异电池 1 的 `MA1` 那一行锚点被本程写成 `… && false` ⇒ 它把**两支一起撤了**、与 `MA` 同形，
  **没有量到"只撤 `io.EOF` 那一支"**。发现于回读日志时；电池 2 的 `N1` 才是那一枚的正确版（红 3 枚）。
  §1 那张表**只用 `N1`**；`MA1` 留在 `mut/` 日志里当"本程也写坏过锚点"的凭据。
② 电池 1 的 `MB` 用 `continue` 造 AC#2 明禁那一形，量出来的是**永不结束**
  （`panic: test timed out after 10m0s`，4 枚用例记**未取到**），不是"重试到预算尽头"。
  电池 2 的 `N2`（删掉那行 `return`）才是票面原话那一形，且它**逐字复现**了前一程 §3.2 的两枚红句。
  §1 表里两枚都列、不相混。
③ **第四批那枚 commit 的 message 里有一处"更正"更正的是本程根本没写过的东西**：
  它说"上一版 §7 里那句 post 树 CI 同形 `PASS=69 FAIL=1` 与 offender 读数是本程写坏的"——
  现量 `grep -n 'PASS=69' 本件` 只有第 155 行一处，而那一处**是引用续程的读数、不是本程的**。
  真实情况是：`PASS=69 FAIL=1` 与那枚 offender 转录**从头到尾没进过本件的正文**，
  它们只在本程的推理里出现过、被本程自己在写 §4/§8 时改成了现量值（`PASS=70 FAIL=0`）。
  ⇒ 记这一条是因为**它正是本仓最该登记的那一形**："我没写过的东西被我当成'我写错了'更正了一遍"——
  这种自纠**看起来像诚实、实际上是没核对盘上文本**，下一个人会照着这条去删一个不存在的东西。
  判据留在纸上：**更正之前先 `grep` 自己那份文件。**
④ §4 里"改前／改后 `--- PASS`"本程给的是**含子格**的数（116／130），与自证件那张表的**顶格**数（63／70）
  不是同一把尺。两把尺都跑在同一份日志上；本件引用的名册是 `sort -u` 那把（106／120）。别把两把当互证。
⑤ §1 里"两批共 **26 发**"也是本程写快了的数（同一批里 ③ 已经认过一枚）。现量正确的两个数是：
  **变异 14 枚**（电池 1 的 MA／MA1／MA2／MB／MC／MD／ME／MF／MG ＝ 9 枚，电池 2 的 N1／N2／N3／N4／N5 ＝ 5 枚），
  **`go test` 调用 18 次**（电池 1 含 baseline 与还原后那一发共 11 次、电池 2 共 7 次）。
  `NOT-RUN` 那一条的判定**不受影响**：18 次里只有电池 1 的 `MB` 那一次出现 4 枚未取到，其余每次 7 枚都在册。
  ⇒ 记这一条是因为**同一个"数"本程连错两次**，而两次都是"凭印象乘/加出来的"、不是量出来的；
  留在纸上给下一个写变异电池的人：**枚数与次数要分开数，且都从 `ls mut*/*.log` 现量。**
