# 票 352 终裁腿 352-v1 判语（★非实现者·只读·裁决者≠实现者）

被裁对象＝实现腿 `352-r1` 的 11 笔（`0fb399c3`①`ecffb93b`②`aacde632`③`c16ff47e`④`06aa2abf`⑤`3ecdea30`⑥`5c9402f8`⑦`1c787282`⑧`8a10a7e3`⑨`e74762b5`⑩`a651c086`⑪）。
本腿⛔ 跑任何 `go test`／`go build`、⛔ 动任何 `.go` 字节、⛔ 翻任何框、⛔ push、⛔ 改写历史。
它交件自报的每一把数一律当未验证＝下面每个数都是本腿另起的尺；尺件⇄逐字输出＝`.scratch/wisp/probes/352/v1/sh/01-my-rulers.sh`＋`logs/01-my-rulers.txt`⇄`sh/02-quote-provenance.sh`＋`logs/02-quote-provenance.txt`。
判据字面以票面为准（`.scratch/wisp/issues/352-…-private-route.md` :14..:17）。
起手锚＝`95691aa4`；交件锚＝本腿末笔（★本趟期间 HEAD 被别腿推进＝凡 blob 逐锚报清，见下）。

## AC#0 现状钉死（改前读数）＝**成立**（内容），唯"逐字⇄出处"那味撞出一处不符（＝第五处，具名下）

本腿尺（量的）＝对入库底片 `logs/ac0-{col,host}-base-s{1,2}-*.raw.txt` 逐枚 `git cat-file blob HEAD:… | grep -oE 'R[0-9]+ [^:]+: isErrNotExist=… isNotSealable=…'` 现量。
现量＝POSIX 六跳 R1..R6 全 `false ⇄ true`（两发逐字同）；Windows R1..R6 全 `false ⇄ true`（R1⇄R2 另带 `isUnresolved=true`），两发逐字同；控制跳＝R8／R9／R12 两侧都 `false/false`、R13 两侧 `false/true`、★R10 平台形不同＝POSIX `false/false` ⇄ Windows `false/true`（两发各自逐字同）。
⇒ 票面 :14 那两枚真值（现语义＝`ErrNotExist` false＋`ErrNotSealable` true）在两平台都量到、逐字红句在库＝**成立**。
它自报"POSIX R1..R6 六枚＋Windows R3..R6 四枚"＝本腿复量＝那张表实列 12 行（Windows 侧另有 R1⇄R2 两行、同行已标"⛔ 经 `wrapPath`"＝它⛔ 把这两枚算进"经 `wrapPath`"那四枚）＝枚数⇄表自洽＝**⛔ 冒报**。

## AC#1 落地那一行＝**成立**（判据＝两分类并存，不是换成另一种不可穿）

本腿现读（量的，⛔ 抄它的表）＝`git cat-file blob HEAD:internal/winsec/winsec.go` :430..432 逐字＝
`func wrapPath(path string, err error) error {` ⇄ `return fmt.Errorf("%w: %s: %w", ErrNotSealable, path, err)` ⇄ `}`。
blob 逐锚＝`bf509bb8`＝`07300bff` → `ecffb93b`⇄`a651c086`＝`f3c1b68e` → 现 HEAD＝`757b16a8`。
★那枚漂移＝别腿 `353-r1` 票〔乙〕（`fa395b33` 起）＝6 行增删，本腿尺 `git diff ecffb93b HEAD -- …| 含 wrapPath 行的增删枚数`＝**0**＝它⛔ 碰那一枚格式串；漂移⛔ 本票射程、具名登记即可。
`git show -U0 ecffb93b` 净变动（量的）＝可执行 **1 增 1 删**＝2 行；函数签名行改前⇄改后 SAME；`sealError` 体 435..446 改前后 symdiff＝**0** 行；⛔ 新增导出名（diff 只那一行）。`go.mod` 那枚＝`go 1.27`（双 `%w` 合法）。
两分类并存（量的，尺＝对入库落地底片逐名现量）＝POSIX `R1..R6` 六枚全 `isErrNotExist=true isNotSealable=true`；Windows `R3/R4/R5/R6` 四枚全 `true/true`；★`R13` 两平台仍 `false/true`＝幂等判定未动、⛔ 造出第三种分类。
★它自报那味"没修到的地方"（Win `R1⇄R2` ⛔ 经 `wrapPath`）＝**如实报了、⛔ 报歪**：本腿源码现读坐实 `pinRefusal`→`placementRefusalWindows`→`preOpenRefusal`（`winsec_windows.go:402`）体内＝`fmt.Errorf("%w: %s: … (%w)", ErrUnresolvedPath, …, ErrNotSealable)`＝双 sentinel 自带且⛔ 带 `fs.ErrNotExist`＝与票面现量第 11 行"⛔ 第二处"逐字同形。它那句"Windows 落地那 4 枚 false/true＝R1/R2/R10/R13"＝本腿逐名现量＝**那四枚名一字不差**（含"⛔ 一枚经改动而翻转分类"＝对：R10⇄R13 改前⇄改后逐值同）＝这条自陈⛔ 遮短板。

## AC#2 消费者面＝**部分成立**＝缺"本腿⛔ 独立复跑实测跑色"这一味（★⛔ 圆判）

本腿⛔ 跑 `go test`（派单预判为真＝起手现量到交件期间 `353-r1`〔winsec〕⇄`321-r1`⇄`109-r1` 在落笔＝跑了必读到半保存态假红）。⇒ 只走既有底片⇄blob。
本腿复算（量的，尺＝对入库 raw blob **整枚文件级** `^=== RUN`／`^--- PASS`／`^--- FAIL`／`^--- SKIP` 计数＋逐名 `comm`，⛔ 它的逐包切块）＝
本机 pair 两侧 `805/497/9/2` ⇄ `805/497/9/2`；容器 `734/455/9/2` ⇄ `734/457/7/2`＝**与实现腿逐包四数之和逐枚相等**（＝它那把包归属尺没串位）。
逐名作差（本腿）＝本机档新增红 1 枚＝`…NeverDeletesAForeignFile342r1`、消失红 1 枚＝同族 `…NeverTruncates…342r1`；容器档新增红 **0**、消失红 2。
★那枚新增红住 `internal/winsec`＝⛔ 票面点名那五包之列⇒**票面"五包两档 12 格逐名新增红 0"在本腿整枚尺下也成立**（五包名集两侧相同）。
竞态族四发并集作差（本腿尺＝首发 raw 的 `^--- FAIL` ＋复量三枚 `.rednames.txt`，先 `grep -v '^#'`）＝B 并集 3 ⇄ L 并集 3、两向 comm 皆 **0**；config 那族 2 ⇄ 2、两向 0。⇒ 归因"两侧并集相同＝⛔ 我所造"在底片上复现（本腿第一版没剥那行 `# RED_NAMES=0` 说明、并集假成 3＝已改）。
具名连带（本腿现读）＝`git cat-file blob HEAD:internal/winsec/private_fail_test.go | grep -n 'errors.Is(err, ErrNotSealable)'`＝**5 枚、行号 37/50/61/72/75 逐字对**＝票面那五个行号真存在；宿主 `TestAC5FailedSealRefusesTheWrite`＋四枚子测试⇄`TestAC5FailureIsNotSwallowedByTheHappyPath`⇄`…EveryProductionCaller322r1` 在本机两侧逐字 PASS；`openleaf_322`/`placement_symlink_113` 具名族在容器两侧逐字同集 PASS。
★缺哪一味＝① 落地后的**现势跑色**（本腿禁跑）；② 仓外⛔ 一枚覆盖落地后的 CI 底片可用（尺＝`git merge-base --is-ancestor`：实现腿 11 笔在 `origin/dev`(`6f808a83`) 上的＝**0 枚**＝`ecffb93b` 根本还没推⇒CI 从未见过落地字节⇒这味天然缺）。⇒ **该由编排者那窗补**（`cmd/wisp` 整包本票明令归它）＝三枚写腿收线＋推送之后跑一次成对色⇄逐名作差（或取那一跑之后的 CI 日志），本判语只登记"既有底片自洽＋成对逐名作差成立"、⛔ 把它当"实测过现势"。

## AC#3 纪律格＝**成立**（三味小账具名、判据内容未被歪）

名册两把独立尺同号（量的）＝锚定行首 `--grep='^票 352 落地腿 352-r1'`＝**11** ⇄ 本腿内存内 awk `index($2,P)==1`＝**11**；非锚定形＝11。★正控＝同两把尺数票 321＝8⇄8（两把都灵、锚尺没空转）。
⚠ 它两件里把第二把尺写作 `index($2,P)==0`（`readings-ac3-discipline.md:12`⇄笔⑨正文）＝本腿读了尺件本体 `sh/06`＝那是 **perl**（perl `index()==0` 正是"行首"）＝实现对、叙述形错一味（awk 里同写法语义相反）。本腿不靠它这句，另起 awk==1 形独立同号。
写面（算的，尺＝逐枚 `git show --name-only`）＝11 笔里 `.go` **只 1 枚**＝`internal/winsec/winsec.go`，且只在 `ecffb93b`；越界 **0**；`pinnedleaf_342_windows.go` 命中 **0** 笔。
禁区逐枚 blob 对撞（量的）＝`bf509bb8`⇄HEAD 下 `winsec_other.go`⇄`winsec_windows.go`⇄`pinnedleaf_342_windows.go`⇄`pinnedleaf_342_fallback.go` 全 **UNCHANGED**；`winsec.go:118..140` 自陈段改前⇄改后逐字节 IDENTICAL。
票面尺（量的）＝29 行、`Status:` 仍在 `:22`、未勾 **4**／已勾 **0**／总框 4（正控票 348＝6/1/7＝尺灵）；`bf509bb8`⇄HEAD 删除行 **0**、numstat `5 0`＝纯追加、前 24 行逐字节 IDENTICAL＝⛔ 改写历史；票尾 bullet＝2 枚（一条自述＋一行 Progress log）＝"尾只许追加一条"的额度已被 `1c787282` 占＝**本腿⛔ 补第二条**、话落本件。
⚠ `8a10a7e3` 正文 `ANCHORED_COUNT=8`＝其时末笔自数（＝不含它自己那笔；本腿逐锚复量：到 `1c787282` 为止＝8、到 `8a10a7e3` 为止＝9、现 11）＝点时读数＋自数差一枚、⛔ 造假（现量随 HEAD 漂＝公共块点名的那族）。

## 四处"入库正文枚数⇄盘上现量"不符＝本腿独立坐实（枚数全对得上）＋层面判语

| # | 笔 | 入库正文那句 | 盘上现量（本腿尺＝逐枚 `git show --name-only` 后 `grep -c .`） |
|---|---|---|---|
| ① | `0fb399c3` | 本笔 3 枚文件 | **5**（漏 `logs/void-fours.tsv`＋第二枚 `sh/01-build-trees.sh`＝件内首行自称"已取代的草稿"） |
| ② | `aacde632` | 本笔 8 枚文件 | **9**（漏 `probe-copy-of-stage.txt`＝★见下第五味） |
| ③ | `06aa2abf` | 1 报告＋27 散件＝28 | **41**（1 报告＋38 logs＋2 sh） |
| ④ | `5c9402f8` 的报告件 `readings-ac3-discipline.md:45` | 票尾 numstat＝`3 增 0 删` | 票尾真值 **5 增 0 删** |

★层面判语（本腿独立读正文、⛔ 采信它的自陈）＝四处**全在"枚数叙述"层**＝正文那句自述的文件计数⇄git 里的真实枚数；四处都⛔ 碰四格判据本体：①②③ 那三笔的判据产出（改前读数⇄那一行落地⇄消费者面作差）本腿各自另尺复跑皆对得上；④ 那句里的方向性断言"纯追加／0 删"＝本腿现验 numstat 删列＝**0**、前 24 行 IDENTICAL＝**判据未被歪，只有"增行枚数"这一枚数写错**。
⇒ 本腿判＝**四格判据内容⛔ 被这四处歪**；但这四处是"入库正文自述⇄盘上不符"这一味纪律缺陷（第④枚还在报告件里＝会误导下一程引用者），⛔ amend、⛔ 改写历史＝由编排者按台账只追加登记。

## ★本腿新撞的第五处（实现腿笔⑪那四处⛔ 覆盖＝具名上报、⛔ 替它圆）

`readings-ac0-before.md:62`⇄`:64` 那两枚标称"**逐字**"、且标注出处＝`logs/ac0-host-base-s1-mut-1-treeB.raw.txt` 的 Windows 行（R1 与 R4），**不在任何一枚入库件里**。
尺＝逐字节 `grep -qF` 整行，逐枚扫 HEAD 下 `probes/352/r1/**` 全 **93** 枚件（剥掉报告件自身＝它当然含自己；含 md／sh／tsv，⛔ 只扫 logs）＝命中 **0** 枚（见 `logs/02-quote-provenance.txt` A 段）；同尺数 POSIX 那三枚引文（`:54/:55/:56`）＝逐字节在它所标注的那枚底片里＝对得上；AC#1 报告件六枚引文逐枚对撞标注底片＝全 YES。
拼接坐实（量的）＝被引那两行里 run 根 `352r1-20261020-160844-1453` 只出现在 `ac0-host-base-s1`，而测试目录名 `TestP352R1ErrClassify1725796622` 只出现在 `ac0-host-base1-mut`（＝它自己登记为**作废发**那一枚；s1 里那两跳的目录名＝`…1935309345`）＝一行里两半各来自一枚不同底片＝该行是手拼的，⛔ 任何一发的原始字节。
判据内容有没有被歪＝**⛔**：本腿另把三枚候选底片里 R1⇄R4 的四个分类位逐枚现量＝全同（`false/true`＋`isUnresolved=true` ⇄ R4 `false/true`＋`isUnresolved=false`）＝改前那两个真值由三枚独立底片同时承载＝AC#0 判据稳；歪的是"逐字"＋出处指针这一味＝下一程若拿它当凭据会追不到原件。
另两处它判"非不符"＝本腿复核＝**它判对**：笔⑥"8 行 CRLF"⇄入库 blob 利尺 `tr -cd '\r' | wc -c`＝8＝对（钝尺 `grep -c $'\r'` 本机报 1948/3470＝行计数，本腿另造 2-CR 件正控同形）；笔⑨那两枚 8＝点时读数随 HEAD 漂（上面 AC#3 那条已具名）。
★但笔⑥那句"CR＝原始读数照实保留"的**层面**要改口（量的，尺件 `sh/02` F 段）＝入库 blob 那两枚 CR＝**0**、工作树那两枚 CR＝8（`core.autocrlf=true`＋`.gitattributes` 那枚 `* text=auto`）＝CR 活在导出树⇄工作树那一层、入库时已被归一化＝"照实保留"在入库件层⛔ 立；判据承载性＝那 8 行全非 `^--- FAIL`／`^=== RUN`／`^--- PASS` 形（本腿现量＝0 枚）＝⛔ 歪判据。

## 另两处本腿具名（⛔ 影响四格判据、⛔ 算冒报，登记给编排者）

① `probe-copy-of-stage.txt`（＝四处不符里"漏数 8⇄9"那枚件）首行逐字自陈＝"只存在于 `D:\tmp` 的拷贝树，⛔ 入库、⛔ 任何跟踪字节"⇄尺＝`git ls-files --error-unmatch`＝该件在 HEAD **TRACKED**（`aacde632` 首次加入）＝**件内自我声明与实际入库相反**；它与 `probe/probe-352r1-errclass_test.go.txt` 逐字节差＝0＝同一件两份拷贝＝⛔ 新字节风险、但"⛔ 入库"这句写进件里会误导下一程（尺件 `sh/02` E 段）。
② `readings-ac3-discipline.md:12`⇄`8a10a7e3` 正文把第二把名册尺写作 `index($2,P)==0`＝本腿读尺件本体 `sh/06`＝那是 **perl**（perl `index()==0` 恰＝行首）＝实现对、⛔ 形写错一味（awk 里同写法＝"找不到"）；本腿⛔ 采它这句，另起 awk `==1` 形独立同号 11⇄11。

## 写面自证⇄本腿自报

写面＝只两枚根＝`docs/evidence/s1/352-v1-verdict.md`＋`.scratch/wisp/probes/352/v1/**`＝共 5 枚件（尺件 2⇄尺输出 2⇄本判语 1），★每落一枚就一笔 commit、`add`⇄`commit` 两处逐枚写显式 pathspec。★⛔ 翻框、⛔ 动票面、⛔ 补票尾第二条 bullet、⛔ 动 `Status:`、⛔ push、⛔ `add -A`／amend／reset／rebase／stash／checkout／clean。
本腿自抓两味照实入档＝① 第一版尺件用 `cut -c1-90` 截 commit 主题＝按字节切中文＝整枚输出坏成非法 UTF-8（`iconv` 在 line 2 报 cannot convert＝公共块 台面四味 ④ 那味、本编队今天第二次撞上）＝改成 `--format` 分行；② 一把 grep 的 `R[0-9]+ [a-zA-Z -]+` 漏咬 tag 里带 `.` 的 `R3 wrapPath direct fs.ErrNotExist`＝改 `[^:]+` 后重跑（改前那发读数照实留在本件叙述里＝⛔ 拿"0 命中"当结论＝正是公共块"报 0 命中前先跑正控"那味，我这次靠正控＋复跑抓到）。
终量对回＝票 352 未勾 **4**／已勾 **0**（本腿一枚没翻）；票面 **29** 行、`Status:` 仍在 **`:22`**、`bf509bb8`⇄HEAD 删除行 **0**＝⛔ 改写历史。

## 追加节（笔⑥⑦⑧＝⛔ 改写上面已入库段、只追加）＝本腿自证违令登记⇄写面枚数更正

★本腿自己一枚违令（具名、⛔ 替自己圆、⛔ amend＝只追加）＝笔⑤ `a0bed7c0` 的 **commit 正文**含被禁字面字符（码位 U+00A7＝那枚分节符，我写 `SPEC-12` 后面那味时把字符本身带进了正文）＝**1 枚命中**。
尺＝`.scratch/wisp/probes/352/v1/sh/03-self-purity-gate.sh`＋逐字输出 `logs/03-self-purity-gate.txt`：正控＝同尺数 `docs/reports/HANDOVER.md`＝304 在场＝尺灵；笔①②③④ 正文各 **0** 枚 ⇄ 笔⑤ 正文 **1** 枚；★交付物本体（本判语件入库 blob）该码位＝**0**、CR＝0、76 行／14082 字节＝干净＝**违令只在 commit 正文那一层、⛔ 在判语内容那一层**。
同族前例＝实现腿 `3ecdea30` 自抓的"入库正文净化"整改＝本编队这条禁令今天已咬两枚腿（编排者一枚、本腿一枚）＝定式＝凡写正文，落笔前先跑那把 `grep -o` 那对字节 C2 A7 的尺＋已知含量正控（本腿这次靠正控才没把"笔⑤＝0"当干净交出去）。
与实现腿那四处枚数不符的关系＝⛔ 相互抵消＝两回事＝它四处＝"正文自述枚数⇄盘上枚数"（本腿判＝错在枚数叙述层、四格判据内容未被歪）；本腿这枚＝"入库正文纯度"（本腿⛔ 违令）。本腿⛔ 拿"我也犯了"替它圆、也⛔ 拿它那份错盖自己这份。
写面枚数更正（上面那段"共 5 枚件"＝到笔⑤为止的真值，现量随 HEAD 漂＝公共块点名的那族）＝本趟到笔⑧ 共 **7** 枚件＝尺件 3（`sh/01`／`sh/02`／`sh/03`）⇄尺输出 3（`logs/01`／`logs/02`／`logs/03`）⇄判语件 1；`add`⇄`commit` 两处仍逐枚显式 pathspec、越界 0、★全部笔的 `.go` 枚数＝**0**、票 352 仍一枚没翻（未勾 4／已勾 0）。
名册尺（本腿自己那枚前缀 `票 352 终裁腿 352-v1`，两把独立形）＝到笔⑦ 为止锚定行首 **7** ⇄ awk `index==1` **7**＝同号（笔⑥ 那发量到 5＝点时读数）。
