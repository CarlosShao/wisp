# 票 343 终裁判语 · 裁腿 `343-v1`（非实现者；只裁⛔ 实现）

起手现量（本腿第一笔）：HEAD＝`d5ba5352`（主题＝342-v1 交件③）；`git status --porcelain -- internal/risk`＝0 行＝车道干净；
名册尺改正形两把对撞＝`^票 343 落地腿 343-r1` **6 笔** ⇄ 单次 `%H%x09%s` 导出＋`index==0` **6 笔**；同尺**非锚定＝7 笔**＝多出 `83ac5f91`
（编排者自己那笔正文引用了腿前缀）＝第十二次同族形状复现、★枚数随 HEAD 漂⛔ 当通则。台面＝**本腿自己** `git archive` 两棵仓外树
（`29d258f9`⇄`31d4a476`、八件套＋`docs/`＋手工三件、各 766 枚 `.go` 证非空）＝⛔ 复用腿的树、⛔ `-overlay`。
逐字读数＝`.scratch/wisp/probes/343/v1/logs/**`（仓库根相对路径；下列每格末尾指到件名）。

## AC#0 起手现跑＝成立

票面三行我逐枚重跑：注释在场（`git grep -n aEquiv 29d258f9`＝1 命中＝该行自己）；`git grep -c 'func aEquiv'` 两侧＝0（rc=1；正控＝同尺形打 `uncertainAnchorMiss`／`aEq` 命中在场）；成因＝`f1033e1e 2026-09-21`（`git log -S` 复跑）。
★同族尺我加到**三形对撞**：尺③纯 shell（`sed -n 's,^[[:space:]]*//,,p'`＋声明形 grep＝与 `go/ast` 无共享实现路径）、尺①′驼峰 `SA`、尺②′严格调用形 `SB`（`foo(`⛔ 空格）、另词表 `SW`。
⇒ `blacklist.go` 注释面同族＝**2 枚**＝`aEquiv`（`:193`＝本票、已改）＋`fires()`（`:310`；`SA=0 SB=1 SW=0`＝只有调用形尺看得见；包内 `func …fires`＝0；他包真身＝`internal/config/unwired.go` 那枚字段）。
★腿那句"两把尺盲区互补"＝**对**（`SA` 看不见纯小写、`SB` 看不见非调用形驼峰）；`bOverrides` 在我三把尺里都是 `IN-PKG param`＝腿自抓的盲区我也复现。
他包两枚（`provenance.go:54 taskScopeID`／`:427 disposalScope`）现验＝锚均在注释行、`SA=1 SB=0`、全仓 0 定义⇒归口票 349。
凭据＝`logs/census/v3-blacklist-start.txt`、`v2-blacklist-delivery.txt`、`census/ac3-strict.txt`、`census/gauge3/*`、`ac5/ac5-gauges.log`。

## AC#1 改准那一行注释＝成立

①交件侧 `git grep -n aEquiv HEAD -- internal/risk`＝**0**（正控＝起手侧 1 命中）；②`git diff -U0` 改动行 **7 枚逐枚以 `//` 开头、非注释行 0**、`--numstat`＝5 加 2 删；③★**内容说实话**＝我自己读了函数体（⛔ 抄腿）：`aEq:271`／`aUnder:273` 都是一行委托 `aRelates`＝"delegate"对；`aRelates:280`＝`if anchorMatches{true}`→`return uncertainAnchorMiss(...)`＝"报 miss 为 hit"**只属于后者**；`anchorMatches:292` 找不到即 `false`＝"单看不那么干"对；`uncertainAnchorMiss:362` 三条件（两侧 root 非空＋⛔ 都 certain＋两 root 互为祖先后代）＝R17 兜底支；`aRuleReason:327`（`:329` 读那两把）＝"要说清哪一把产出"对；A 档 `switch:222` 在 B 档 `:244` 之前＝"⛔ quietly dropping a path from A into B"顺序也对。
⚠一处如实登记（⛔ 本票射程）：句子**前半**那句既有文字"Both sides of every A-anchor comparison below"未穷尽 A 档两枚 `aAnyForm`（`:225`/`:236`＝形状测试、无锚）＝轻微过称；那半句⛔ 是腿写的⇒⛔ 退件、⛔ 在这里改（⛔ 扩写面）。凭据＝`logs/ac1/ac1-gauges.log`＋上列 census 件。

## AC#2 ⛔ 动任何可执行字节＝成立

blob 起手 `c1799cd24afce19e6accef23db23a7d7ee746dec` ⇄ 交件 `62f2dddbde7c72eb12cbeb4cf1ec1edc2708254b`（我从 object store 现取、`cat-file -t`＝blob）。
★**剥注释＋去空白去空行后两侧可执行字节 sha256 相同＝`b1610acfcfb84851dd53a47c1230026d2b9a28415e1968a16b167a9a29924908`**、各 232 行／6396 字节、`cmp` 同、`diff` 空；第二把独立形状（先删整行 `//` 再剥行尾）＝**同一串**。⇒ 结构上证这枚注释改动⛔ 可能影响任何判级／那 8 枚牙。
原始 blob 两侧 sha256 不同（404⇄407 行、CRLF＝0）＝差异只在注释、`diff` 只落 `193,194c193,197`。★腿报那串 `97796acc…`／265 行我也**独立复现**＝另一种剥注释形状（`grep -v '^[[:space:]]*//'`，10 形扫出的 `v8`）⇒ 两串都成立、⛔ 谁抄谁、⛔ 据此退件。
`git grep -c 'func aEquiv' HEAD`＝0＝⛔ 造符号；受保护仪器起手⇄HEAD 逐枚 blob 等值 13 枚全 SAME（含四枚仪器 `a429cdae`/`f5f888de`/`e7b13c39`/`5147c189`、三枚冻结件、`thresholds.go`）；`git diff --name-only 29d258f9 137044b6 -- internal cmd tools`＝**只 `blacklist.go` 一枚**。凭据＝`logs/ac2-01..04`、`logs/ac2-protection.log`。

## AC#3 上不上仪器＝通用形**不该上**（复现腿方向＋加一条它没有的否证）；有界词表形**可上**、归口＝票 349、★⛔ 本腿造门

★**通用形＝否。** 我自己那把尺、生产件（⛔ `_test.go`）逐枚分类：`SA` 驼峰＝6 红／3 假阳＝**50%**（`hostPathBox`、`assessorFor`＝别包真符号；`macOS`＝平台散文）＝与腿逐枚对上；
`SB` 严格调用形＝6 红／4 假阳＝**67%**（`DEFERRED` 标记约定、`Nutstore` 厂商名、`Defer`＝别包真 method、`CloseScope`＝票 160 故意删的历史名、那句自写"Before it…"＝词面与假话⛔ 可分）；真阳＝`fires`＋**第三枚 `Complete`**（`syncdirs.go:453` 写 `Complete()==false`、真名 `SyncDetectionComplete:472`）＝同族第三枚。
★**最硬那条＝切词自身不稳**：宽松调用形 `SBL`（`foo (`）同包同侧＝**142 枚唯一红**（`the candidate (`、`side (`、`them (` 这类散文）。分子分母随"剥/切"定义漂⇒⛔ 存在通用规则；任何区分都是人写的一条＝退化成手工 allowlist（＝`tools/d22scan/allowlist.txt` 那形＝冻结件）。
**有界词表形 `^a[A-Z][A-Za-z0-9_]*`＝值得上**：生产件里唯 6 枚名（`aAnyForm aEq aEquiv aRelates aRuleReason aUnder`），起手侧红 **1 枚＝`aEquiv`＝真阳、假阳 0**、交件侧红 **0**；代价如实＝抓⛔ `fires()`（纯小写）、抓⛔ `taskScopeID`/`disposalScope`、⛔ 泛化、词表要跟文件族维护；仍属契约邻接＝人工批准。
★立票动作**已完成**＝票 349 `AC#3` 逐字要求"复核此结论＋具名写哪一类词表能有界"⇒ 本票⛔ 欠一张新票、⛔ 顺手做门。凭据＝`logs/census/ac3-strict.txt`、`ac3-table2.txt`、`sb-names.txt`。

## AC#4 颜色与门禁＝成立

Windows（我两棵树、串行两发、`-count=1 -timeout 30m -v`、四数⛔ 混用⛔ 相减）：基线 `29d258f9`＝**247/140/1/1** rc=1 ⇄ 交件 `31d4a476`＝**247/140/1/1** rc=1；红名册作差 **0**、`RUN` 名集作差 **0**（247⇄247）。唯一红＝`TestC26Rewrite…EverySecurityLeg`＝**存量**（红句逐字点名 `internal/panel` 四枚读者＝静态源码审计、读不到本包注释一字；在册 `A380`／票 200／216；与 `339-r1`／`343-r1` 跨腿同名对上）＝⛔ 计入新增红。`TestResolvePerCallBudget` 两发都 PASS（1.34／1.21 s）＝在册争用那条⛔ 触发（同波 `342-v1` 吃 CPU 我按口径盯过、⛔ 记成腿的错）。
★台损如实登记（⛔ 计入判据）：我**第一发**用 POSIX 形 `GOTMPDIR`⇒同包报 5 枚红（`mklink /J D:/tmp/…`＝`Invalid switch - "tmp"`＋前置件两侧分隔符混用）；两侧红名册**逐字相同**⇒可证与腿字节无关；改 `cygpath -w`＋`;` 重跑＝247/140/1/1。两发都在册（`logs/pair/`＝台损发、`logs/pair2/`＝判据发）。
Linux 容器（我两棵树、⛔ 复用腿的树；`golang:1.27`、`CGO_ENABLED=0`、`TMPDIR=/tmp/gt`）：基线 **195/115/1/1** ⇄ 交件 **195/115/1/1**；唯一红＝同名存量；SKIP 两侧同枚；台面故障尺＝none。
★**票 336 那 8 枚牙＝8/8 全红、逐名对得上**（现成台件一字⛔ 动、blob `454da3ee`/sha256 `7346a690` 我用前验过；env 指我自己的 `treeDel`；runner rc=0）：1a git-credentials｜1b wisp-config-under｜1c ＋wisp-config-eq｜1d ssh｜1e aws｜1f kube｜1g dpapi-protect｜1h local-credentials；未突变正控 TEETH rc=0（8 枚 operand 全 `class=A`）／LEG rc=0；每发落刀前 `needle_before=1`⇄后 `needle_after=0`＋`branch_removals_on_disk=1`＋`changed_lines=2`；**还原哈希 11/11 match=YES**；`shared_worktree_dirty_internal_cmd=0`。
门禁四数（只在共享树、跑前 `unset GOTMPDIR TMPDIR TEMP TMP`、env sanity 四枚皆空）＝gofumpt `-l blacklist.go` **rc=0/列 0 枚**（新字节 LF；存量假脏 4 枚具名登记⛔ 我动）／`d22scan.sh` **rc=0 clean**／`check-path-length-budget.sh` **rc=0 GREEN**（16457 跟踪、over-budget 60、roster 60、名单外 0）／`go vet ./internal/risk/` **rc=0**。⛔ `./cmd/wisp`、⛔ `go test ./...`、⛔ `internal/winsec/**`。凭据＝`logs/pair2/pair2-summary.log`、`logs/linux/*`、`logs/teeth-03-runner.log`＋`logs/teeth/*.log`、`logs/gates/gates-runner.log`。

## AC#5／AC#6 纪律格＝成立

名册 6 笔逐枚 `git show --name-only`＝并集 24 枚路径 ⇒ **⊆ 三处写面**、`.go` 恰 **1** 枚＝`blacklist.go`、写面外 **0** 枚；票框终量＝`^- [ ] **AC#`＝**6**／`^- [x]`＝**0**＝一枚没翻；票名⇄`Status:` 逐字未动；票尾 `6 0`＝纯追加。
★腿三处自抓复核＝**登记、⛔ 退件、⛔ 污染任何一格读数**：①超"≤4 枚正文件"走拆具名分件＝现量 3 report＋4 读数件全部 ≤160 行（123/139/48/105/47/156/105）＝口径对；②笔数抄死被自己尺抓到（5⇄6）已改"只留尺"＝对（我这趟同样⛔ 抄死）；③乱码＝`U+FFFD` 0、`??` 0、`iconv` 往返非 UTF-8 0、0 字节件 0、`.out` 0、被禁字面字符（U+00A7）腿写面 **0 枚**（★正控 HANDOVER＝199、票 343＝0）＝与编排者复跑同结论、我⛔ 抄他。
⛔ push（本腿 0 次）。本腿写面＝本判语＋`probes/343/v1/**`＋票尾**一条**。凭据＝`logs/ac5/ac5-gauges.log`、`logs/leg-hygiene.log`。

## 两件"归谁"

① **分票⛔ 并进 343＝对**：票面 `AC#1` 射程逐字＝"那一行注释"＝`:193`；腿⛔ 扩射程、具名交回立票权＝按派单办。票 349 判据形状**大体够、两处要补**（我⛔ 落它一字、⛔ 动 `:310`／`provenance.go`）：
(a) 它 `AC#4` 把台件路径写成 `scripts/probes/336/v1/sh/run-ac2-sweep.sh`＝★**该路径 HEAD 不存在**（现跑 `git cat-file -e` 报 does not exist；真身＝`.scratch/wisp/probes/336/v1/sh/run-ac2-sweep.sh`＝存在）＝同族缺陷"指名一枚不存在的凭据路径"、公共块"列路径前逐枚 `cat-file -e`"正是它；
(b) `AC#6` 自写"本票**六格**"而 `^- [ ] **AC#`＝**7**（AC#0…AC#6）＝枚数与文字不符。★另交回一枚它名册⛔ 含的真阳＝`Complete`（`syncdirs.go:453`，真名 `SyncDetectionComplete:472`）＝同族第三枚、在 349 两枚射程文件**之外**⇒349 扩名册或另立一枚。
② `AC#3` 那一半见上格＝腿方向复现＋加强；有界词表形的立票已由票 349 承接（⛔ 我造门、属契约邻接＝人工批准）。

## 必答 (i)–(v)

(i) 新句⇄盘上机制**逐行对上**（我读函数体）：报 miss 为 hit 只有 `uncertainAnchorMiss:362` 那一支（`aRelates:280` 第二个 return）；`anchorMatches:292` 找不到即 `false`＝"单看不那么干"属实；`aRuleReason:327`/`:329`＝把两把分开读的读者。红/绿无关、语义真值＝**真**。
(ii) 同族终数（三形对撞后）＝`blacklist.go` 注释面⛔ 定义符号**现剩 1 枚**＝`fires()`（`:310`、在册票 349、⛔ 本票射程）；`aEquiv`＝**0**（三把尺都⛔ 再报）。他包两枚现验均在注释行⇒归口 349＝对。
(iii) 独立结论＝**通用形不该上**（50%／67%／切词放宽即 142 枚红＝尺自身不稳）；**有界词表 `^a[A-Z]` 形可上**（起手 1 真阳 0 假阳、交件 0）、⛔ 抓 `fires()`、归口票 349。假阳逐枚＝`hostPathBox`/`assessorFor`/`macOS`/`DEFERRED`/`Nutstore`/`Defer`/`CloseScope`。
(iv) 腿随件入库 14 枚原始 stdout＝**⛔ 缺陷**：公共块"成对读数⇄产生它的机器输出⛔ 拆两笔、拆了复跑对不上"＋`339-r1` 跟踪 26 枚 `.log` 先例＋帽射程只管"我撰写的正文件"（派单 12:2x 明写⛔ 外推）。我这趟同形（`logs/teeth/` 19 枚同批）。反过来才要防＝为凑帽截断 `go test` 输出＝篡改读数。
(v) 结案判据＝六格**全成立** ⇒ 票 343 **可 `-done`**。`AC#3` 本腿判"通用形＝否"，"该上"那一半已由**票 349** 承接⇒票面"若判该上⇒具名立新票"**已满足**、★本票⛔ 欠新票。零未满足判据核对表＝AC#0 成立／AC#1 成立／AC#2 成立／AC#3 成立（结论＋代价已给、⛔ 造门、立票已完成）／AC#4 成立／AC#5 成立。
★两条**不阻结案**请编排者处置＝票 349 的 `AC#4` 台件路径不存在、＋"六格"⇄7 枚不符（我⛔ 落它一字）。翻框只由编排者凭本判语做。
