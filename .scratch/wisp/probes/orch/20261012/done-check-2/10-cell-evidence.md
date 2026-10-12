# `10-cell-evidence.md` — A 节：票 304／票 309 十二格的**在盘性**（本腿自己现跑的尺）

锚＝`49153fa8`（＝票 304 收口笔）⇄ `db1694bd`（＝票 309 收口笔，其父）。provenance 一律＝**工作树**（除注明 `git show HEAD:` 的 blob 面）。
尺（逐字，本腿现跑）＝
```
find .scratch/wisp/probes/<票> -type f -printf '%s %p\n'
find .scratch/wisp/probes/<票> -type f -size 0 -print
grep -nE 'rc=|build_rc=|_rc=' <件>
```
名册分母（本腿现量）＝`probes/304` **154 枚件／其中 0 字节 11 枚**；`probes/309` **515 枚件／0 字节 63 枚**（逐枚名册＝本腿 `logs/roster-zero-byte.txt`，rc=0）。
⛔ 任何一格用"大体齐了"结。下面 **★＝判"缺"或判"凭据是转述"**。

## 票 304（六格）

| 格 | 那格要求的读数／件 | 在盘性（件＋`wc -c`＋是否自落 `rc=`） | 恒真问（换个反形会不会也放行） |
|---|---|---|---|
| `AC#0` | 表①名册＋表③代价＝`probes/304/a1/{10-residual,20-roster,30-costs}.md`（16505／12941／6567 B）＋九把尺 `logs/grep-m1..m9.txt`＋切分件 `partition-{66-18-48,start18,sync48}.txt`＋进程名册 `logs/proc-*` | **在盘、非 0**。表②（三形实测）腿自己交不了 ⇒ 凭据＝编排者亲测 `probes/orch/20261011/logs/orch-A850-shapes.txt`（2617 B）＋`orch-A850-shape-b2.txt`（729 B）＋台件 `shapes2.sh`/`shapesb2.sh`/`mockllm-procs.ps1`。件里有命令口径行（`ruler:`／`target test:`）与逐发 `rc=`（`rc=0`／`rc=1`＋`timeout_panic_lines=1`），**(b) 第一发假读数在件里自落为 `rc=0 … timeout_panic=0` 并被具名作废**、b2 重打才真红 ⇒ 形制对 | 表②不是恒真：b2 里"坏世界证据行"（`panic: test timed out after 3s`）与"残留 PID"同发落盘 ⇒ 反形（帽设太长、panic 没触发）会被 marker 行=0 抓到。**但表①③的裁决者是编排者本人＝实现者自评**（`AC#0` 是只读普查格，票面 `:25` 授权他裁形；⛔ 由此推及"六格都由他判也行"） |
| `AC#1` | 落地笔 `e875f00f`＋件 `probes/304/r1/10-impl.md`（5602 B）；判语＝`304-v1b/10-ac1.md`（3224 B）＋`logs/ac1-{commits-detail,percommit,reuse,scope,bullet}.txt` | **在盘、非 0**。逐笔尺在 `ac1-percommit.txt`（13635 B）＝本腿独立复跑过同一把（见 `20-commits.md`），两枚 spawn 点同批、五条硬约束各自 0 行 | 不恒真：反形＝任一 spawn 点没改 ⇒ 名册差件立刻点名第二枚泄漏源。硬约束反形（动 golden）＝`git diff --name-only` 非 0 |
| `AC#2` | ★要求的形＝**成对两发**＋逐枚具名 PID。r1＝`logs/paired-fwd.txt`（4953 B）／`paired-rev.txt`（6839 B）＋8 枚 `roster-*`/`cmdline-*`/`construct-*`/`watch-*` 逐发件；v1b＝`15-ac2-shots.md`（4392 B）＋`logs/ac2-*`（5 枚 0 字节 `-compile.txt` 的**同伴件带 `compile_rc=`／`go_test_rc=`**，本腿逐枚核过） | **在盘**。反形件里 `pid=25952 ppid=22740 … REPO_GOLDEN PARENT_GONE`＋`rc=1`＋`bad-world-evidence` 行逐发齐；正向 4 发 `count_repo_golden=0` 且同行有世界证据 ⇒ 形交到了 | ⚠ **一处真残余（不是假读数，是射程）**：`paired-*.txt` 的头两行自证 `spawn points: SP1=… :72 via -run` —— SP1 那一发的世界**是从 openaichat 的用例里走 helper**，`Q-304-r1-5`（其余 16 枚走 helper 的测试文件⛔ 逐枚重跑）**仍挂着**，六格齐这句话⛔ 蕴含它已清 |
| `AC#3` | 八把尺 `probes/304/r1/{30-gates.md,logs/g-*.txt}`＋台件 `gates.sh`（3075 B） | **在盘**。`g-build.txt`＝`build_rc=0`（11 B）、`g-vet.txt`＝`vet_rc=0`（9 B）＝**rc 有、命令不在同一枚件里**（命令在 `gates.sh`＋`30-gates.md`）⇒ 按派单口径判"自落 rc 达成、命令靠同目录台件"。三枚 0 字节件（`build-1.txt`／`vet-1.txt`／`vet-reverse.txt`）＝stdout 本体（无诊断即 0 字节），同伴件带 rc ⇒ ⛔ 判"那格没交"，也⛔ 把空件当凭据。census totals `g-summary.txt`（462 B）逐字⇄票面基线相等；格式两把并排 `g-format-{worktree,headblob}.txt`（225／140 B）各落 `format_*_rc=0`＋**射程目录逐枚列出** | ★**恒真面（本腿现证）**：`gofumpt -l` 那两把**没有正控**——件里没有"塞一枚故意没格式的文件必然红"的那一发。⇒ 一枚坏掉的尺（路径打错／二进制缺席被 `|| true` 吞掉）也会交回"名册 0 枚"。**能让它⛔ 红的输入**＝同射程放一枚 `x.go` 内容为 `package p` + 未对齐的 var 块，名册必=1；这一把⛔ 本腿能打（要跑 gofumpt＝Go 面）＝**归编排者** |
| `AC#4` | 正典 `probes/orch/gate-roster-304.md`（4007 B，§1 四条名册＋家族尺＋三问）＋现量一发 `logs/orch-A85x-gate-live.txt`（764 B，四条各 `count=0 tasklist_rc=0` 且带 `first_line="INFO: No tasks…"`＋`gate_rc=0`）＋条件落政件 `logs/orch-A856-gate-family.txt`（995 B） | ★**条件那一发是坏尺造的假 0**：件里 `family_rc=0`，但同一枚件把 PowerShell 的 **ParserError 原文**自己吐在 `names=` 后面（`You must provide a value expression following the '-join' operator`／`Unexpected token '\",\"'`）⇒ 那条家族尺**根本没跑出计数**，它上面的 `count=0` 是失败副产物，且 `count_rc=0` 是外壳 rc⛔ 是 PS 的 rc。A856 §2 那句"当场跑一发＝家族 `count=0`（显式计数落盘，⛔ 拿空输出当 0）"**没有凭据支撑** | ★本腿代跑两把补上（`logs/dc2-R1-summary.txt`／`dc2-R2-summary.txt`，均 rc=0）：**正典 §1 那条命令逐字重跑＝`rc=0`、0 命中、stderr 0 字节**；**同一条管线换成必然存在的名字（`powershell.exe`）＝count=3** ⇒ 正典那把尺**有牙**、今天真值确为 0。⇒ 缺的不是读数而是**"落政那一发"必须用正典逐字命令重打一次**（旧件保留、⛔ 删） |
| `AC#4b` | 正典 §2 三句（非零计数须指到一棵树／⛔ 裸枚数进终态锚／⛔ 由裸枚数生成停法账）＋v1b 判语 `30-cells.md`："我今日四名册全 0＝指针规则未被现形" | 件在盘、判语具名承认**未被现形** | ★**这就是恒真格**：判据只在"全 0"的台面上验过 ⇒ 正典 §2 就算写反了，读数也一模一样。★本腿造出**可现形的输入**并现跑（`logs/dc2-R3-summary.txt`／`dc2-R4-summary.txt`，rc=0）：今天这台机器 `msedgewebview2.exe`＝**24 枚＝4 棵具名树各 5 枚＋4 枚命令行里没有 `--webview-exe-name=`**；那 4 枚逐枚 `parent=msedgewebview2.exe`（父在、需再上一跳才到根）⇒ 这一族正是 §2 那句"指不到树就写归属不明"要咬的形状。**结论：本格由"未现形"升为"现形可测"，但现形那一发是 304 的腿与编排者都没交的（本腿代交，仍归编排者落账）**；同时⚠ 本腿读数与票面 `:17`／A856 的"18 枚＝3 棵×6"不同代＝同一把尺两个世代，⛔ 谁覆盖谁 |

**票 304 的"每格都有凭据"这句**：`AC#0`–`AC#3` 站得住（`AC#3` 补一枚正控、`AC#2` 挂着 `-r1-5` 的射程句）；**`AC#4` 的条件凭据是假清洁**（坏尺造 0）；**`AC#4b` 到本腿为止没有一次现形读数**。两枚都是"格翻早了"那一族，⛔ 一句"大体齐了"能盖。

## 票 309（六格）

| 格 | 那格要求的读数／件 | 在盘性 | 恒真问 |
|---|---|---|---|
| `AC#0` | 三笔 `2a37d612`／`27059960`／`ff301daf` ＋件 `probes/309/a1/**`；裁语＝票面 `:51-59` | 在盘（含 8 枚 tracked `.go` 快照，见 `30-debt-ledger.md` §4）。裁语两半都交：名册核实＋更正（票面追加节）＋两形代价 `21-shape-b-cost.md` | 只读普查格，判语＝编排者（票面授权）。⚠ 它裁 Ⓑ 那句"对空 id 的实际防御＝零"是全票承重句，`AC#1..AC#5`⛔ 闭合它 |
| `AC#1` | 恒真形**双向**都有牙：仪器 `internal/tools/assessorfor_taint_instrument_309_test.go`；腿读数＝`probes/309/r1/logs/mutation-rigs.log`＋`probes/309/v2/logs/rigA.txt`（5963 B）／`rigA-instrument.log`（8789 B）；**编排者自跑**＝`logs/orch-A851-309ac1-rig.txt`（1312 B）＋`orch-A851-SHOT1/2*.out`＋`control.out` | **在盘且是本腿逐字读过的一发**：件里有 pristine `rc=0`、SHOT1 摘肢 `rc=1`＋逐字 `--- FAIL: Test309AssessorForEmptyTaskIDDropsTaintDetector`、SHOT2 `WithTaintDetector(nil)` `rc=1`＋同一枚 FAIL、每发 `restored=IDENTICAL sha256[:16]=995c23d2ef6ea16c`＋`dirty=0`。红句体里带"这一腿的存在就是为了…"的自证 ⇒ 两向都有牙 | 不恒真（反形必红，本腿逐字见过红句全文）。⚠ **已知射程边界**＝两枚装配肢（`b.injected`／`b.prov == nil`）⛔ 具名仪器＝`Q-309-v1-prov-nil-limb-blind`／`Q-309-v2-1-assembly-limb-blind` 原样挂着，编排者 `A847` §4.4 裁"⛔ 造第二枚仪器"⇒ 这是**裁过的范围**，⛔ 再当欠读数，但"六格齐"⛔ 蕴含"三枚前提都有牙"（票面 `:80` 那句标题就是三枚前提） |
| `AC#2` | 枚数钉 `internal/tools/toolrequest_nail_309_test.go`＋三对成对读数（`new`／`make([]T,n)`／定义型转换）＋分母自证 | 在盘：`probes/309/{r2,r3,v2,v3b}/logs/**`；★本腿自量"每一形是否有过带红标记的名册行"＝**7 形全有**（`logs/dc2-R5-label-redshots.txt`，rc=0）。分母 293⇄295 与每根名册在件里 | ★**三处"七／六"数目不一（本腿现跑的尺，⛔ 任何一腿报过）**：①**票面 `:33` 修订句＝"钉所认的七形之内"并逐枚列 7 枚**（限定／同包裸名／别名字面量、var 无初值、`new`、`make`、定义型转换）；②**钉的函数头注释（HEAD blob，本腿取 `logs/nail-head.go.txt`）逐字写 `Six shapes count`（三枚自 r2＋三枚自 r3）**；③**词面尺 `grep -oE 'return "<label>"' | sort -u`＝7 枚**，但那 7 枚是 `{qualified, same-package-bare, alias, literal, new-, make-, conversion-}`＝**含 `:183` 那枚今天打不到的兜底 `"literal"`、缺 `:260` 的 `var-decl-no-initializer`**（v3b 自己在判语里也点了这枚把戏）。⇒ 本格**判据的射程边界⛔ 只有一种数法**；下一腿照任何一个数都会和另外两个打架。能让"七形之内"这判据⛔ 红的输入＝把 7 形逐枚各打一发注入⇒名册从 1→2 且红句点名标签（本腿⛔ 能打＝要跑 `go test`＝**归编排者**） |
| `AC#3` | "只动注释"由尺坐实⇄内容锚在盘⇄枚数带尺重述 | ★**本腿自己那把（⛔ 用编排者与腿报过的数，逐笔现取现比）**：对 `16e86bae`／`d75db9ee`／`a1f84d39` 三笔各取 `git show <h>^:…bridge.go` 与 `git show <h>:…` 两份，**两边先证非空**（剥整行注释后各 **854 行**），再 `cmp`＝**三笔全 IDENTICAL**（`logs/dc2-R10-bridge-comment-only.txt`，rc=0）；同笔 numstat＝`73/0＋32/3＋134/0`、`44/9`、`14/4` ⇒ 增删只在注释行 | 不恒真：本腿第一把也差点踩"剥成空行 vs 丢掉整行"那枚雷（编排者在 `A855` §2 就栽过）；本腿用的是**丢掉整行注释**且两边行数并排落盘 ⇒ 空比对造不出假清洁。⚠ 剩一枚真残余＝`bridge.go:959-963` 那句已收窄，但钉头 `:52-63` 的漏形名册与 `:213` 的 "Six shapes" 互相指不到（上一条） |
| `AC#4` | 三门禁 rc 自落行＋格式名册两把并排＋逐笔 pathspec | 在盘：`probes/309/{v1,v2,r2,r3,v3b}/logs/gate-*.log`（v2 具名 build/vet/d22scan/pathlen 各带 `rc=`、且 build/vet 那两枚 11／9 字节＝只有 rc 行、stdout 本就 0 字节）；v3b 把 r3 的 4＋1 枚 0 字节件**全部重打并落 rc**（件在 `probes/309/v3b/logs/*-reshot*`） | ★口径不齐（编排者 `A851` 末段自己记的）＝`rigB/C/D/E` 的 `.txt` 汇总件里没有 `rc=` 行，退码证据散在 `.code`／逐发 `.log` 两组 ⇒ 同一枚"rigA"在两句里指两件不同文件。这一处⛔ 判违规、判**口径不齐**；能红的反形＝拿一枚不存在的件名当"已重打"（本腿逐枚按 `wc -c` 与名册对过，未发现指空） |
| `AC#5` | 四家（307／40／197／242）分工具名＋`A838` §3 ⛔ 在本票重记 | 在盘＝`probes/309/r1/90-final.md` 末节＋v1 判语（成立，含一枚措辞更正：③ 的文案出处是**产码** `subagent_197.go:263-266`，票 197 正文 `grep`＝0 命中；编排者 `A847` 表里已复跑） | ★**这一格恰恰是"⛔ 同一件事记两次"的自查格，而它⛔ 覆盖新长出来的重复**：`AC#5` 只对了 307/40/197/242 四家，本腿在 C 节抓到的是 **309 自己家里同一件事记了 2–4 遍**（六组），以及 `Q-304-a1-H` 在 304 票面与 310 票面各一份。**⇒ `AC#5` 判"成立"⛔ 蕴含"两票的欠账名册无重复"** |

## 本腿欠的读数（要跑 Go 才能判的，一律具名归编排者）

- 票 304 `AC#3` 格式尺的正控一发（故意没格式的文件必使名册=1）。
- 票 304 `AC#2` 的 `Q-304-r1-5`（其余 16 枚走 helper 的测试文件逐枚重跑）。
- 票 309 `AC#2` "七形之内"逐形各一对注入⇒红／挪开⇒绿（并先把 6⇄7 数法统一）。
- 票 309 `AC#1` 两枚装配肢要不要仪器（编排者已裁⛔ 造＝本腿只复述裁语，⛔ 重开）。

rc=0
