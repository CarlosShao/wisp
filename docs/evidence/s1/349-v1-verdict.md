# 票 349 终裁裁决表（`349-v1`＝非实现者裁决腿）

- 起手锚（我自己第一笔 `git log -1` 现认）＝`cd7b946a`（13:08:05）；落笔过程中 HEAD 漂到 `9f4e4f69`＋ ⇒ 本表**所有 blob 尺都指具名提交树**，⛔ 工作树字节。
- 两枚算得的锚：base＝`070de7fa`（＝本腿第一枚改码笔 `bf7f13bc` 的父）⇄ del＝`ed8db1ea`（`349-r1` 名册第 7 笔）。★与派单写的起手锚 `7b633fa9`：三枚目标 blob **逐枚同一**（`62f2dddb`/`d634aae0`/`7125b08a`），两对都跑、两对都同。
- 实现腿＝`349-r1`，名册＝两把锚并集**7 笔**（`票 349 落地腿 349-r1` 4 枚 ⇄ `349-r1` 3 枚；7 枚全 `merge-base --is-ancestor` 现验在祖先链）。正控＝同尺非锚定 `--grep` 吃 13 笔＝坏尺形。

## 逐格判语（七格，⛔ 我翻框）

**AC#0 起手现跑＋同族终数＝成立**。我自己那把尺与它不同形（注释全标识符超集 ⊖ `go/ast` 声明名册，★含形参/字段/局部/接收者＝绑定名，另出"像符号"子集 SYM ⇄ 模块级真身 MODHIT）：prodOnly 超集 base 452 ⇄ del 448 行；SYM 63⇄59；★射程三枚文件 **22 ⇄ 18** 行；作差 base-only＝`blacklist.go:310 fires`／`provenance.go:54 taskScopeID`／`provenance.go:427 disposalScope`＋同行 `Defer`／`syncdirs.go:453 Complete`＝**恰票面四枚**（多出的 `Defer` 与四枚同句、非第五枚残留）；del-only 只 1 行＝改准后的真名 `DisposalScope`（别包 `disposal.go:129` 有定义＝我尺的 MODHIT 直接坐实）。**交件侧 TRUE＝0**。

**AC#1 逐枚改准＝成立**。①交件侧四枚死名**精确形**（词边界／调用形）各自 0 命中，正控非零：`SyncDetectionComplete` 26 行／`scopeID` 34 行／`DisposalScope` 全仓 50 行；★坏尺复现＝裸子串 `Complete` 两侧都 27 命中。②`git diff -U0` base⇄del 全 8 枚 `^[+-]` 行都以 `//` 开头＝**非注释行 0**（numstat `1/1`·`2/2`·`1/1`）。③改准后的名字**都指向盘上真身**（我逐枚现验）：`aEq:271`/`aUnder:273`→`aRelates`、`SyncDetectionComplete:472`、`Detector(scopeID)` `provenance.go:858`、`internal/plugin.DisposalScope:129`＋`Defer:186`。

**AC#2 ⛔ 动可执行字节＝成立**。命门＝**剥注释结构尺我自己另起三形**（⛔ 抄编排者两串、⛔ 抄它的 GS1/GS2/GB），★两侧一律 `git cat-file blob`：
- GS-B（整文件 perl 状态机：字符串/字符/原始字面量保留、注释剔除→删全部空白→sha256）＝`9de81d0d…39d07e`／`d3c64f86…8b2655`／`c6bdd86a…de3c5a1`，**base⇄del⇄现 HEAD⇄派单锚 `7b633fa9` 四锚同串**；
- GS-E（第五形＝`go/parser`＋`go/printer`：剥 `File.Comments` 与全部 `Doc/Comment` 字段后重新打印、去空白→sha256，与任何文本剥法结构性不同）＝与 GS-B **逐枚同值**（三枚、四锚全同）⇒ 一把要文件先能被语法解析、一把只当字节流，两把独立尺同串；
- GS-D（只删整行注释、其余字节一字不动，cmp 级最严形，⛔ `tr`）＝三枚 `cmp` IDENTICAL，sha `97796acc…`／`851972c9…`／`68f2d7cc…`＝我从零算出的值（与票 343／`349-r1` 先例同值＝独立复现，⛔ 抄）；
- 非注释非空行数两侧同＝**232／621／435**；★尺灵敏正控＝注一枚 `package riskX` ⇒ 哈希变（`bcdab391…`／codelen 6165）、注一枚纯注释改动 ⇒ 哈希**不变**（仍 `9de81d0d…`）；
- ★两把文本形在此面等价的形状凭据（⛔ 巧合）＝三枚文件的整行注释 142／439／190、含 `//` 的非整行行 5／25／16、原始字面量内整行注释 0／0／0、`/*` 行 1／0／0 ⇒ 差异面恰被 GS-E 那把 AST 形补住。
- ★CRLF 坑复验：`provenance.go` **blob CR=0／LF=1115**，磁盘副本 **CR=1115** ⇒ 拿磁盘比必假红（本腿全程 blob，含 `git grep` 都打具名提交树）。
- 仪器面 base⇄del：`anchor_fallback_teeth_336_{windows,other}_test.go`／`anchor_miss_reason_339_other_test.go`／`anchor_miss_reason_334_test.go` 四枚 blob **UNCHANGED**；三枚冻结件／`internal/observe/thresholds.go`／`go.mod`／`go.sum`／`.gitignore` 全 UNCHANGED；四枚死名"定义形"两侧 0＝⛔ 为凑话造符号。

**AC#3 ⛔ 上通用仪器＝复核成立（我独立跑一把）**。我的 precision（我这把更宽的尺）：射程三文件 **4/22＝18%**、全包 prodOnly **4/63＝6.3%**、未过筛超集 **4/452≈0.9%**；方向与它一致且更硬（尺越宽越差）。最硬否证复现＝★SC 声明形：射程文件侧 **4 命中／GREEN-pkg 0**，含 `_test.go` 全 package **13 命中／GREEN-pkg 2**（两枚真绿是测试件里的 `func Test…` 前缀切词，⛔ 一枚真残留）＝放宽切词只给一把 0 真阳的尺加噪声，与 `343-v1`"放宽即 142 枚红"同向。有界词表形 `^a[A-Z]…` 我只复量⛔ 造门：三枚产码文件注释名册 **5 枚**（`aAnyForm aEq aRelates aRuleReason aUnder`）两侧全有定义＝红 0；★这形⛔ 抓 `fires()` ⇒ 要不要上属契约邻接＝人工批准，载体＝票 350（不阻结案）。

**AC#4 颜色与门禁＝成立（含一处归因订正，见"归因"节）**。成对四数（★我自己的两棵树 `D:\tmp\wisp\20261019\349v1\tree{Base,Del}`，`.go` 各 766 枚，八件套＋docs＋手工三件，⛔ 复用别腿树；`-count=1 -timeout 30m`；四数⛔ 混用⛔ 相减）：
- Windows 台损形（GOTMPDIR POSIX 形）＝base `247/136/5/1` ⇄ del `247/136/5/1`；
- Windows 干净形（GOTMPDIR `cygpath -w`）＝base `247/140/1/1` ⇄ del `247/140/1/1`；
- Linux（docker `golang:1.27`）＝base `195/115/1/1` ⇄ del `195/115/1/1`；
- 红名册逐名作差：★**三形内 base⇄del 全 0 枚**＝本腿新增红 0；唯一合法存量红＝`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（`A380`／票 200／216）；`TestResolvePerCallBudget` 六发全 PASS（1.20–1.34s）＝争用假红没触发；票 346 那两枚摆动枚在我台面绿。
- ★票 336 那 8 枚牙＝**8/8 全红**（`MU-1a…1h` 逐发 `needle_before=1→after=0`、11 发还原哈希 `964ee75a…` **match=YES**；未动形 pristine 8/8 `clause=true`）；台件一字⛔ 动＝blob `454da3ee` ⇄ 交付锚 `c315fa4c` 同值；⛔ `-overlay`、突变只打我仓外树物理字节。
- 门禁四数（只在共享树、跑前 `unset GOTMPDIR TMPDIR TEMP TMP`）＝`gofumpt -l` rc=0（blob 形三枚 listing=0；打工作副本时 `provenance.go` 标脏＝该件磁盘 CRLF **存量**，起手侧同状态、本腿⛔ 顺手重排）／`d22scan` rc=0 `clean`／`check-path-length-budget` rc=0 `VERDICT GREEN`／`go vet ./internal/risk/` rc=0。⛔ `./cmd/wisp` 整包、⛔ `go test ./...`。

**AC#5 禁区＝成立**（票面 4 枚点名路径逐枚 blob 对撞未动；见 AC#2 节名册）。**AC#6 纪律＝成立**：票框终量尺 `^- [ ] **AC#`＝**7 枚未翻**、`[x]`＝0（本腿⛔ 翻）；逐笔 name-only ⊆ 写面，★正控＝同尺打 `349-r1` 并集 7 笔⇒`.go` 恰 3；`add`/`commit` 两处逐枚字面 pathspec（共享 index 上有 `342-r1`）；⛔ push。

## 归因（编排者点名要我独立定性的一枚）

**四象限**（base⇄del × 台损档⇄干净档，各一发，串行）决定性读数：★**跨象限作差＝恰好那 4 枚**（`…MovedAncestor116`／`…HomeRewriteAtVerdict116`／`…ReResolvesThroughJunctionToRealTree`／`…NotExemptedStaysFailClosed`），且在 base 侧与 del 侧**同样 4 枚**；台损档两发的红句里逐字带着 `mklink /J … Invalid switch - "tmp"`（＝POSIX 形 `GOTMPDIR` 传进 `cmd /c mklink`）。⇒ 我的独立定性＝**这 4 枚是台面造的红、⛔ 存量在册红**，与编排者初步裁同向；**但 `349-r1` 的"与 343-v1 同包读数逐字同名＝存量"这顶帽子是错的归因**：它对上的是 `343-v1` 的**台损发**（`pair/`＝247/136/5/1，同一形），而 `pair2/` 干净发＝247/140/1/1 那 4 枚消失；两串数我自己重打复现。★核心结论不受影响：**新增红 0 成立**（我三形两侧逐名同集＝量到，⛔ 算的）。建议具名订正一句进台账：该腿四数没错、错在"存量"这个定性词，且它**只跑了台损那一形**（未做干净档对撞）＝形状欠一味。

## 两处"复核编排者"（派单指定）

① 交回的 4 枚候选（`…junction_windows_test.go:362 and`／`taintmatch.go:256 documents`／`syncdirs_test.go:298 valueName`／`taintmatch_test.go:171 winStrs`）：我逐行现读 blob＋我自己的 SC 名册，**同意 ⛔ 动 ⛔ 立票**。凭据：`and`＝`// Case 10: env-var and ~ expansion feed the pipeline.`（冠词，我 SC 形同样吃到＝同族假阳）；`documents`＝`…the immutability the type documents keeps holding…`（散文动词）；`valueName` 在串位符 `"<key>|<valueName>"` 内、同处 `key`/`data` 同形；`winStrs`＝注释明写 **dead** 的历史名（`git log -S` 现验两笔：`73b60571` 有、`c091ee2b` 后⛔ 定义），句子已自带"removed"限定＝诚实、正是 `AC#1` 要的形状。★我的尺另摊出同族 `the`/`so`/`that`/`whose`/`marker`/`name`＝**同一味 SC 切词假阳**，四枚候选不是孤例、也⛔ 第五枚真残留（⇒ 它"⛔ 属同族"的判我复核成立）。
② 编排者就地改票 349 三处笔误 vs"实质判据⛔ 就地改"：**裁语＝不冲突**。凭据＝`2962d978` 对票面那笔 numstat `+4/-2`，两枚删行正是"假路径"与"六格"那两句、改后就地留着原句自陈（"原句'六格'＝…枚数笔误、见 `A1047`"），名册扩第四枚走的是**追加段**；三处全在"引用件路径真身⇄自身枚数自述"射程＝⛔ 新增判据、⛔ 放宽断言、⛔ 改 `Status:`／票名／翻框（我复量框仍 7/0）。**实质判据的射程只到"判它过要什么为真"**；引用路径与枚数属读数自校正，就地标＋留原句即可。★一条改进建议（⛔ 判据）：今后凡在票上扩射程，追加段首行写明"本段扩射程／不扩射程"，免得与"就地改"混读。

## 自抓与诚实栏

- 我这趟⛔ 抄过编排者串：结构尺三枚哈希、四象限四数、precision、名册枚数、候选判定，全为自跑。★自抓两枚**坏尺形**（都⛔ 污染任何结论、当场改正后重跑）：①`logs/` 生成脚本里我在同一发 shell 内先 `cd` 进仓外目录、后续 `git`／相对路径尺全部空跑打印 `No such file or directory`＝那批"0"是假读数，已拆成两枚脚本重跑（`logs/10-redsets.log` 前的空转读数⛔ 入库）；②第一把作差 `awk` 用空格切列（源件是 tab 分隔）⇒ 得 `base=0 del=0` 的空名册，改成 `-F'\t'` 后才是 `62⇄58` 那份。★另一枚形状坑我在自己尺上摊开＝`gs.txt` 里 GS-B 与 GS-C 三枚两侧同串同值，说明射程三枚文件的注释**全是整行 `//` 形**（无行尾注释、无被注释掉的代码）⇒ 两形在此面等价；我仍另跑第三形 GS-D（⛔ `tr`、cmp 级）做独立对撞。
- ⛔ 造门、⛔ 翻框、⛔ 动 `.go` 字节、⛔ 改票面、⛔ push、⛔ 删临时件。`internal/winsec/**` 我只⛔ 读工作树字节、⛔ 跑该包（现量 `winsec.go` blob `67783ae0` LEG⇄HEAD 同值＝`342-r1` 那趟尚未落该文件）；`internal/risk` 消费者方向实测＝`go list -deps ./internal/risk` 含 winsec（`winsec_c26.go`），全程读具名提交树快照。
- 未定案上报（停手，⛔ 自行假设）：`AC#4` 干净形与台损形并存时该以哪形为"在册色"＝派单已给 A1048 口径，本腿照用；另请编排者定一处**台账具名**承接上面那枚归因订正（本腿只出判语、⛔ 代提）。
