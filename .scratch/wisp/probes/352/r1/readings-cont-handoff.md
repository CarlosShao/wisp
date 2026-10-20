# 352-r1 续跑（cont）· 交件自证＝核原 9 笔对票面四格的射程 ⇄ 补齐缺的一味 ⇄ 具名上报

★身份＝票 352 落地腿 `352-r1` 的续跑（写腿·winsec `%w`）；主题起头逐字「票 352 落地腿 352-r1」（⛔ 换形）。
★续跑⛔ 再动任何产码字节：那一行 `%v`→`%w` 已落＝笔 `ecffb93b`（`winsec.go` blob `07300bff`→`f3c1b68e`）；任何新 `internal/winsec/**` 改动都算越界＝本腿零碰（尺见 `logs/13-cont-roster.txt` R4/R5）。
★本腿⛔ 给任何 AC 盖章（AC#0..AC#3 四格全归非实现者裁）；本腿⛔ 翻框（现量 0 勾／4 未勾，见 R7）。
起手锚＝续跑第一笔 `git log -1` 现认＝HEAD（原末笔 `8a10a7e3` 起，编排者 17:2x 现量 9 笔名册对回）；⛔ push（本腿笔全在本地 `dev`）。

## 一、核原 9 笔对票面四格的射程缺口（判定＝只读＋复尺，⛔ 盲信编排者数）

逐枚盘上尺（尺＝`sh/13-cont-roster.sh`，输出 `logs/13-cont-roster.txt`，两把名册尺对撞）：

| 格 | 产出落点（盘上尺现量） | 判定 |
|---|---|---|
| AC#0 改前钉死 | `readings-ac0-before.md` 72 行＋`logs/ac0-*` 4 真发（tsv `SHAPE=OK`/`RESTORE=OK`）＝POSIX R1..R6⇄Windows R3..R6 都 `fs.ErrNotExist=false`/`ErrNotSealable=true` | **齐** |
| AC#1 落地那一行 | `ecffb93b` 那一行（`-U0`＝净 1 增 1 删）＋`readings-ac1-landed.md` 65 行＋`logs/ac1-*`＝落地两分类并存 `true/true`（逐名 `logs/11-ac1-pername.txt`：POSIX R1..R6 六枚、Windows R3..R6 四枚） | **齐** |
| AC#2 消费者面 | `readings-ac2-consumers.md` 92 行＋`logs/perpkg-diff.txt`＋raw＝五包两档 12 格逐名新增红 0 枚；`private_fail_test.go:37/50/61/72/75` 宿主钉两档 PASS 对回 | **齐** |
| AC#3 纪律格 | `readings-ac3-discipline.md` 84 行＋`sh/04-gates.sh`/`sh/06-discipline.sh`＋`logs/discipline.txt` | **缺见下** |

**射程缺口＝1 味**＝票面现量第 9 行逐字留的账「调用点…落地腿逐枚现数并点名」＝原 9 笔盘上⛔ 有这一味（尺＝`grep -rln '调用点' probes/352/r1` 除本腿新件＝0 命中）。已补齐＝`readings-ticket9-callsites.md`（24 枚调用点＋1 枚定义，逐枚 file:line 点名，起手锚⇄落地锚同数＝那一行⛔ 增删任何调用点）。

## 二、成对读数落点（盘上尺，⛔ 采信任何正文里的数）

- AC#0 基线侧四真发 `logs/ac0-{col,host}-base-s{1,2}-mut-*.raw.txt`＝`isErrNotExist=false isNotSealable=true`（col FALSETRUE=7 枚、host=8 枚；两平台各 2 发）。
- AC#1 落地侧 `logs/ac1-{col,host}-s2-mut-2-treeL.raw.txt`＋第二发 `ac1-*2-*`＝POSIX `true/true` 6 枚（R1..R6）＋Windows `true/true` 4 枚（R3..R6）；同发 false/true＝Windows R1/R2/R10/R13（R1⇄R2⛔ 经 `wrapPath`＝⛔ 本票射程）。
- AC#2 两档逐名作差＝`logs/perpkg-diff.txt`（12 格 NEW_RED=0）；`panel` 那 7 枚存量红两侧逐名同批（尺 `logs/08-named-evidence.txt` A 段＝BASE_FAIL=1 LAND_FAIL=1 逐名在场）；容器档 `winsec` 74/46/0/0 两棵全绿。
- 四数 tsv＝`logs/four-counts-errclass.tsv`（含作废两发 RC=0/125 SHAPE=VOID＝台损照实登记、⛔ 计入判据）。

## 三、门禁四数各自 rc（续跑复量，尺＝`sh/19-cont-gates.sh`，只在共享工作树、跑前 unset GOTMPDIR/TMPDIR/TEMP/TMP）

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| 1 | `gofumpt -l internal/winsec/winsec.go` | **0** | 空列出＝干净（只验自己这枚） |
| 2 | `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`（internal/=603、cmd/=143 Go files） |
| 3 | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`（over-budget=60＝roster=60＝not in roster 0＝⛔ 新增） |
| 4 | `go vet ./internal/winsec/` | **0** | 输出 0 字节 |
| 4b | `GOOS=linux`／`GOOS=darwin` `go vet` | **0**／**0** | 只登记"能编译＋vet 无声"，⛔ 与本机档成对判红绿 |
| 存量具名 | `gofumpt -l internal cmd tools` | 0 | 4 枚既有 CRLF 件＝`agent/approval/pending_read.go`⇄`agent/tools.go`⇄`risk/provenance.go`⇄`cmd/wisp/models.go`＝⛔ 我动、⛔ 顺手重排 |

## 四、票 352 框终量（须 4 未勾／0 勾＝本腿一枚没翻）⇄ 票尾 numstat（须纯追加）

- 框尺现量（`logs/13-cont-roster.txt` R7，行首锚＋自校）＝`CHECKED=0 UNCHECKED=4 ANY=4`（SELFCHECK_SUM=4 成立）；正控同两把尺数票 348＝6/1/7＝尺灵。
- 票名⇄`Status:`⇄`-done`（R8）＝对回起手锚 `bf509bb8` 那版逐字同＝`TICKET_NAME_AND_STATUS_UNCHANGED=YES`（`Status:` 仍"未开工"）。
- 票尾 numstat（R8）＝`5 增 0 删`＝`MINUS_CONTENT_LINES=0`＝纯追加；本腿名下票尾 bullet＝1 枚（`MY_TAIL_BULLETS=1`）。
- 续跑本腿对票 352＝⛔ 再追加（原 9 笔已含那条唯一 bullet）；本交件只在 `probes/352/r1/**` 落件。

## 五、★具名上报：四处"入库正文自述枚数 ⇄ 盘上现量"不符（⛔ 我改写历史，只登记）

公共块读数纪律"完成通知一律当未验证＋引数带尺名"＝续跑逐枚 `git show --name-only` 现量对回原 9 笔正文里的**枚数自述**，发现四处数字不符（＝皆枚数叙述⛔ 判据内容，四格产出本身经复核＝齐）：

| 笔 | 正文自述 | 盘上 `git show --name-only` 现量 | 定性 |
|---|---|---|---|
| 笔① `0fb399c3` | "本笔 3 枚文件" | **5 枚**（logs/void-fours.tsv＋probe/…go.txt＋sh/01×2＋sh/02） | 正文少报 2 枚（漏数 logs 1＋sh 1） |
| 笔③ `aacde632` | "本笔 8 枚文件" | **9 枚**（6 ac0 raw＋four-counts.tsv＋probe-copy-of-stage.txt＋readings-ac0） | 正文少报 1 枚（漏数 `probe-copy-of-stage.txt`） |
| 笔⑤ `06aa2abf` | "1 枚报告＋27 枚读数散件"＝28 | **41 枚**（1 报告＋38 logs＋2 sh） | 正文少报 13 枚（散件实为 38 枚 raw＋红名册，另漏数 2 枚 sh） |
| 笔⑦ `5c9402f8` 的报告件 `readings-ac3-discipline.md:45` | "票尾 numstat＝纯追加（`3 增 0 删`…）" | 票尾真 numstat＝**5 增 0 删**（`git show --numstat 1c787282`） | 报告件那句"3 增"⇄真值"5 增"＝少报 2 枚增行（**0 删＝纯追加方向仍成立**＝此格判据未被歪，只是那句把数写小） |

★另两处经复核＝**非**不符（避免冒报）：笔⑥ `3ecdea30` 正文"8 行 CRLF（mklink 自己吐的）"⇄利尺 `tr -cd '\r' | wc -c`＝**8 CR 字节**（钝尺 `grep -c $'\r'` 这台 shell 报行计数＝假 1948/3470＝钝尺自门见 `logs/16-cr-real-and-claims.txt` A 段 CONTROL 行）＝⑥ 那句对得上；笔⑨ `8a10a7e3` 正文"ANCHORED_COUNT=8"⇄现量 9＝**随 HEAD 漂**（笔⑨是其时末笔自数 8，续跑加第 9 枚＝9＝该腿 `logs/handoff-number-reconcile.txt` 已预告此漂）＝⛔ 造假。
⇒ 四处不符＝续跑**具名落此件＋上报编排者**，⛔ 自行改写任何已入库正文（⛔ amend/reset）。

## 六、⛔ 产码越界自证（本腿全部笔的路径并集）

`logs/13-cont-roster.txt` R4＝`OUTSIDE_ALLOWED_WRITEFACE=0`、`GO_FILES_TOUCHED=1`＝只 `internal/winsec/winsec.go`（住笔②）；禁区逐枚命中 0（`pinnedleaf_342_windows.go`/`winsec_{other,windows}.go`/`internal/risk`/`internal/config`/`cmd/wisp`/`internal/agent`/`internal/panel`/`frontend`/`design`/`go.mod`/`go.sum`/`.gitignore`/`ci.yml` 全 0）；写面活体尺 `internal/winsec` 续跑起手 90 分钟窗落笔＝**0**（⛔ 活体、本腿唯一写者）。

## 附·本续跑新落件（只 probes/352/r1/**）

`readings-ticket9-callsites.md`（补齐那味）／`sh/13-cont-roster.sh`＋`logs/13-cont-roster.txt`（名册两把⇄写面⇄blob⇄框尺⇄numstat 一把尺）／`logs/handoff-number-reconcile.txt`（上一腿遗产件，本腿入库＝⛔ 0 字节、⛔ .out、U+00A7=0、CR 利尺 0）／`logs/16-cr-real-and-claims.txt`（CR 三把尺对撞＋逐枚正文⇄盘上枚数对账）。逐枚 LF、UTF-8 VALID、U+00A7=0（尺 `iconv`＋`tr`，见临时根 `sh/17-stage-and-shape.sh`）。
