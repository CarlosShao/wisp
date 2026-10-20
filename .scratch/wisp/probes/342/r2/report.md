# 342-r2 交件报告＝票 342＋票 348 Windows 半区续体落地（写腿·仪器件）

腿＝`342-r2`（写·`internal/winsec/**`·产码＝★只新增／修改 `*_test.go`）。起手锚＝本腿第一笔 `git log -1` 现认 `3abbb6391a8a497c0332e5363f24c8a250a7eb78`（2026-10-19 21:00:56 +08）。授权＝台账 `A1046` 第④节「批 342＋348 落乙」＋`A1049` 增量＝那批里⛔ 做完的 Windows 半边，⛔ 新批准。⛔ 翻框⛔ 判语⛔ push。

★**三句话结论（给收件者，按派单"做不完就交已做＋老实报名"）**：
1. 派单①那枚确定性窗口件＝**原句形状（祖先在 open 之后翻）在 Windows ⛔ 能造**，原因＝平台的句柄锁（实测 S2/S2b）＝要造它必须新增一枚产码接缝＝**第二道人工批准＝本腿停手上报、⛔ 自加**（＝本波唯一"没做完"的那格，老实报）。
2. 派单②（两栏分栏）＝**做到了**：damage⇄refusal 在四枚件里各是各的数、件二另挂阳性对照证明那枚 0 是读数。派单③（叶子复验给能答的形）＝**做到了**：`openPinnedLeaf` 的 link 面第一次被量到在路上（49/46/48 与 31/32/36 每 400 发、guard 只 0~1）。
3. 交件＝8 笔（前 7 笔已入库、本笔＝更正笔）、新增 `.go` **2 枚**、**产码字节 0 枚**（⛔ 动 `wrapPath`、⛔ 动 `:118-140`、⛔ 新导出面、⛔ 函数形状变化＝尺见下表）；两票尾各一条纯追加（2/0＋2/0）；⛔ 任何一格被我宣布闭环、⛔ 翻框；★⛔ 报"全绿"＝落地上仍有 r1 三枚竞态钉在册红＋票 348 的 `deleted` 分子仍取⛔ 到（缺世界那格本腿只把它从形容词升成分栏读数，⛔ 圆场）。

## 逐笔提交（每落一枚就 commit＝枚数帽按公共块口径＝参考值）

| 笔 | 提交 | 内容 | name-only 核对 |
|---|---|---|---|
| ① | `4f4f8419` | 新增 `ancestorwindow_342r2_windows_test.go`（671 行＝确定性窗口件三枚 AC#1-one／AC#2-two／AC#1-three＋两台柱尺） | 1 枚 ⊆ 写面 |
| ② | `a782196e` | 新增 `leafrace_342r2_windows_test.go`（341 行＝叶子复验按"同目录裸兄弟名"能答的形，dispatch③） | 1 枚 ⊆ 写面 |
| ③ | `2626119a` | 同件补 damage 分栏（每发采样 truncated/gone/unreadable＋逐字首发＝dispatch②两栏分栏） | 1 枚 |
| ④ | `4ff76a0c` | 名册一行改 ASCII（d22scan 对字符串字面量⛔ 豁免＝U+26D4⛔ 进正码；先修在入库侧） | 1 枚 |
| ⑤ | `0e487f06` | 注释面更正＝probe7 未跑成作废登记＋两处引用改指真凭据（★可执行行变化 0＝尺 `git diff --cached -U0` 每行以 // 开头，numstat 3/1＋11/6） | 2 枚 |

| ⑥ | `2b3b1787` | 两票尾各**一条**只追加 bullet（★numstat 2/0＋2/0＝纯追加；⛔ 动任何 `.go`） | 2 枚 ⊆ 写面 |
| ⑦ | `4b784fea` | 散件入库＝**36 枚**（`report.md`＋`logs/` 30 枚逐字日志＋`sh/` 五枚尺件）；★⛔ 任何 `.go` 字节＝本笔产码改动＝0 枚 | 36 枚 ⊆ 写面 |
| ⑧ | `70aac8c0` | 更正笔（追加）＝⑦正文声称 35⇄实际 36；★该笔正文另有一句"两枚 `.go` CR＝0＝我的新字节全 LF"＝**半错**（尺只打在 `internal/winsec` 那两枚上、⛔ 覆盖散件 36 枚＝我把"我的 `.go` LF"写成了"我的新字节 LF"＝笔⑨改口径） | 2 枚 ⊆ 写面 |
| ⑨ | `3efa4bf0` | 散件真损修复＝重生成 `logs/gates-and-roots.txt`（`cut -c` 字节截断损）＋两处口径更正＋`sh/utf8-and-codepoints.sh` 双尺件与输出入库 | ⛔ `.go` |
| ⑩ | （本笔） | 行数自照＋两枚尺件自身缺陷修复＝本件行数⛔ 再写死在表里（写死＝同一笔内自指＝笔⑨复发根源）＝改指尺的输出行 `REPORT_MD_LINES`（帽＝正件 ≤160 行＝越没越由那行说了算、⛔ 我手涂）。★同笔另修两枚尺件的真缺陷＝`sh/utf8-and-codepoints.sh`⇄`sh/gates-and-roots.sh` 里大括号块 `cd` 进子目录后⛔ 回原处⇒尾部 `cat`／`wc` 拿相对路径失败＝**日志内容全对、脚本报 rc=1**＝一把说不清退码的尺⛔ 能交＝`OUT` 改绝对路径；本笔提交枚数＝5（`report.md`＋两枚 `sh/`＋两枚重跑后的 `logs/`＝`utf8-and-codepoints.txt`／`gates-and-roots.txt`）＝⛔ 任何 `.go` |

★枚数对账尺（本腿自抓两枚＝同族坑，逐枚具名）＝① 交件后现跑 `git log -1`＋`git show --name-only --format= | grep -c .`＝笔⑦声称 35⇄实际 36＝追加更正（⛔ `--amend`⛔ `reset`）；② 散件那枚 `cut -c` 字节截断＝`iconv` 抓到、我原先三把尺（`find -size 0`／码位 perl 尺带 `:encoding(UTF-8)` 只 warn⛔ 停／grep）**全都⛔ 报**＝"尺全绿⇄件真损"那一形在本腿自己散件上复发一次＝已按码位尺＋`iconv` 双尺重做（输出＝`logs/utf8-and-codepoints.txt`）。
★名册尺（锚定行首、⛔ 假设其后任何字符＝A1045）＝`git log --format='%h%x09%s' | grep -P '^\S+\t票 342＋票 348 落地腿 342-r2'`＝笔⑧后＝8 笔（首 `4f4f8419`⇄末 `70aac8c0`），本笔＝第 9 枚＝入库后以本笔现跑为准。

枚数＝新增 `.go` **2 枚**（`ancestorwindow_342r2_windows_test.go` 终态 673 行／`leafrace_342r2_windows_test.go` 终态 393 行＝逐枚 `git show --name-only` 已对回，见 `sh/self-proof-rulers.sh` 的 R2/R5 段）⇄ 修改 `.go`＝同一批 2 枚的后续笔（③④⑤），★产码（非 `_test.go`）改动＝**0 枚**；散件＝36 枚（笔⑦入库，明细见上表与 `logs/`＋`sh/`）。

自证（★尺全部**锚定到本腿自己那五笔 `.go` 提交**＝`3abbb639..HEAD` 这个区间里还落着 `227-v2`／`346-a1`／`296-r2`／编排者的笔＝⛔ 用未锚定区间数＝A1045 那味；尺件＝`sh/self-proof-rulers.sh`，逐枚现跑输出已在本腿终端记录）＝R1 非 `_test.go` 路径 **0 枚**（⇒★⛔ 动 `winsec.go:430` 的 `wrapPath`、⛔ 动 `:118-140` 那段 only-half-closed 自陈＝R1b `winsec.go` 命中 0 行）；R3 顶层大写起头新增＝4 枚、逐枚读形全是 `TestXxx` 用例入口（测试二进制按命名约定必须导出、⛔ 进任何产码二进制），★严格尺（⛔ 用例入口）＝**0**＝⛔ 新导出名、⛔ 函数形状变化；R4 产码 numstat 行＝0；R5＝区间内动过 `internal/winsec/**` 的 5 笔全是本腿的＝写面互斥守住、越界 0。

## ①那枚确定性窗口件——做⛔ 到派单原句的形状，做到的是另一枚（★这就是本腿最该被读的一句话）

派单原句＝"父钉已拿到、祖先在 open 那一跳**之后**才翻"。仓外台架（probe1..6＋probe8，逐字行＝`logs/bench-probes-verbatim.txt`）实测：

- **L＝落地路由握着钉的那两跳（seal 跳＋删除跳），平台⛔ 肯让任何祖先被 rename**（`Access is denied`／被钉目录本身＝另一进程使用中）。＝那一世界在落地上⛔ 存在＝要让它存在得把翻转塞进 `pinnedParent`→`openPinnedLeaf`→删除之间＝**一枚产码内的新接缝（新名字／改函数形状）＝第二道人工批准＝本腿⛔ 自加、⛔ 绕**＝已具名上报（＝lead 那封的"停手上报"项）。
- 基线侧同两跳里祖先**翻得动**（probe4-A＝`os.Remove` 按名穿过翻转的祖先删掉外来叶子＝S6 行＝348 的形状在 Windows 真存在）——但它只在⛔ 钉态可达。
- **做得到的确定性世界＝名字那一跳**（probe6＝落地句柄带 `FILE_SHARE_DELETE`、基线⛔ 带）：seal 判语之后、删除之前，⛔ 人点名的对象可被 rename 顶进调用者的名字；路由的删除按"钉目录里的名字"走＝**3/3 发确定性删掉⛔ 被点名的对象**（件一 AC#1-one＝damage 栏分子＝`worlds entered=3/3, unnamed objects left with no name=3/3`×三发⇄基线 `entered=0, refused=3`×三发）。＝落地的钉锁住了祖先、⛔ 是名字＝本腿新量到的残差，具名写进注释，⛔ 写成任何一票的达成。
- 件三 AC#1-three＝把"平台锁"本身做成成对读数（B 形 rename 祖先 ok ⇄ L 形 Access is denied，两态⛔ 经产码路由＝⛔ 模拟）。

## ②damage ⇄ refusal 两栏（⛔ 混成一个数）

件二 AC#2-two＝稳态祖先 junction（⛔ 接缝⛔ 注入⛔ 撞运气）：REFUSAL 栏三发逐字同一条 C26 拒绝（`risk: path traverses a reparse point (junction/symlink) not covered by reparse_point_exceptions`）＝348-v1 的"缺世界"第一次被写成分栏读数；DAMAGE 栏＝外来见证按⛔ 名读、`truncated=0/gone=0`，且带**阳性对照**（同一翻转拼写下 by-name remove 必须删得动一枚别的名字——三发都删得动＝尺咬得住＝0 才是读数）。
件四（叶子竞态）两栏同发并打＝落地 `re-check=49/46/48` ⇄ `truncated=50/46/49`；基线 `re-check=0/0/0`（那扇门锁在基线产码上根本⛔ 在场）⇄ `truncated=6/4/10`＋`returnedNil=123/112/120`。采样⇄复位相竞＝只⛔ 多算、照实写明。

## ③叶子复验给它能答的形

终裁已证带 reparse 位／父指针的形在 `OpenRoot` 那跳被 containment 面拒。件四按"链接目标＝同目录裸兄弟名"造世界，名册里 containment 行**自门**（命中>0＝仪器红＝世界造错了）。实测＝`openPinnedLeaf` 的 link 面 31/32/36（第一轮）与 49/46/48（分栏轮）/400 发答话、guard 只 0~1＝342-v2 记为"到⛔ 达"的那扇门第一次量到在路上。⛔ 拿它当"两半门都关上了"——SameFile 面 0 命中（probe6 说明它需要⛔ 含 reparse 位的顶名＝那形要的是另一个世界）、且 `:118-140` 的更新权⛔ 在本腿。

## 成对四数与逐名红册（本机 `go1.27.1 windows/amd64`，`-v -count=1 -timeout 30m`，四数只在 `^=== RUN`／`^--- PASS`／`^--- FAIL`／`^--- SKIP` 取）

台面＝自建两棵仓外导出树（`git archive` 基线 `fde1df02` ⇄ 落地＝起手现认 `3abbb639`＋本腿仪器字节；八件套＋`docs/`；`.go` 枚数 766⇄775 非空尺）。临时根＝`%LOCALAPPDATA%\Temp` 短形；★容器侧⛔ 讲＝本机测＝GOTMPDIR 绑 D 卷那轮的两枚 SID 红（`TestGateJudgesThePrivateSetByResolvedSID`／`TestSealNarrowsAndNamesThePrincipalItRemovedBySID`）＝台损族（A1053 同族形）＝短根两轮（`base-full-2/3`、`land-full-2/3`）里 0 命中＋共享树⛔ GOTMPDIR 定向复跑 rc=0＝具名台损、⛔ 计入判据。

| 读数 | RUN/PASS_TOP/FAIL_TOP/SKIP | 红名册（逐名） |
|---|---|---|
| 基线树 full r2 | `113/64/0/0` | — |
| 基线树 full r3（含我全部五笔字节） | `114/64/1/0` | `…LeafRecheck…342r2`＝★尺有牙那枚（基线⛔ 门＝红） |
| 落地树 full r2 | `118/64/3/0` | r1 的三枚竞态钉（truncation／never-deletes／leaf）＝**存量红＝⛔ 我引入** |
| 落地树 full r3 | `119/66/2/0` | r1 的两枚竞态钉＝同上存量 |
| 落地树⛔ 我的字节（`treeLandMinus`） | `115/61/3/0` | 同三枚＝**逐名证明与我的字节无涉** |
| 定向·三枚确定性件×两树×各 3 发（短根轮） | 落地 `3/3/0/0` ⇄ 基线 `3/3/0/0` | 全绿；件一形状**相反**＝落地 `entered=3 refused=0 damage=3/3` ⇄ 基线 `entered=0 refused=3 damage=0`（逐字行＝`logs/land-det-shortroot.log`／`logs/base-det-shortroot.log`） |
| 定向·叶子件×落地 3 发 | `1/1/0/0`×3＝`re-check=49/46/48` | — |
| 定向·叶子件×基线 3 发 | `1/0/1/0`×3＝`re-check=0` | 红句逐字＝`NOT YET MEASURED even on the answerable form` |

★成对读数⛔ 混用⛔ 相减；竞态行⛔ 要求逐字相同＝复现量级与形状；红名册逐名作差＝落地侧"新增红"⛔ 一枚是我的钉（我的四名逐枚在 `--- PASS` 行上）。

### liveness 钥匙照实（派单必答项）

派单点名的两枚钥匙＝`pinnedParentBeforeOpen`／`provedLinkCycles`：尺＝`Select-String` 扫我这两枚新文件＝**各 0 命中**（这两枚名住在 POSIX 族 `ancestorworld_342_348_unix_test.go`／`ancestoropen_342_unix_test.go`／`ancestorunlink_348_unix_test.go`＝本腿⛔ 动⛔ 借它们的绿）。本腿用的是**具名等价钥匙**，逐发点名：

- 件一＝`worlds entered=3/3`（＝窗口真被进入＝plant 在产码路由内落地）＋`realSealRanFirst=true`×3×两树；
- 件二＝阳性对照行 `by-name remove err=<nil>` ＋ `control name is now=absent` ×3×两树（＝damage 栏那把尺咬得住）；
- 件三＝B/L 两态逐字 `err=<nil>` ⇄ `err=… Access is denied.`（两树各 3 发都在）；
- 件四＝`proved flips=55/52/54`（`reparseAt` 紧跟 rename 复证）＋逐发 `flippedCycles>0`；基线侧 `flipped=38/33/38` 同发在场＝红是"门⛔ 答"⛔ 是"世界没动"。

## 去牙⇄反形两发（⛔ 共享工作树＝只打 `git archive` 导出树的物理字节、⛔ `-overlay`、串行）

- **发 1（去牙）**＝摘 `openPinnedLeaf` 的两半复验（锚 `if entryIsLink(now) {` 与 `if !os.SameFile(pinned, now) {` 各带 `_ = f.Close()` 形＝函数体内容锚、⛔ 行号；落刀前 1/1＝期望、落刀后新锚 1/1＋旧锚（仅 gate 的那枚）1＝对账）⇒ 我的叶子件 **2/2 发回落红**（红句逐字＝`AC#1-leaf NOT YET MEASURED ... re-check answered 0 across 400 attempts with 67 proved flips`），件一照绿＝尺指向被摘的那支⛔ 全盘红。`go vet` rc=0。mut 哈希 `AB8AF5F5…36B00DD`。
- **发 2（反形）**＝`sameObjectPinnedWindows` 的调用换成同签名"逢名就拒"孪生（真 gate 留在文件里＝⛔ 死码⛔ 未用变量；锚 1→0／新锚 2）⇒ `TestAC2OverwritePromise…322r1` **六枚子测试逐枚红**（四枚调用者 `secret/store.go:79`／`config/migrate.go:93`／`secret/migrate.go:189`／`agent/spill.go:166`＋entry stays private＋exclusive 仍拒）＝尺⛔ 一把两种形状都放行。mut 哈希 `A77BA4E0…9F0D68`。
- 还原＝两发都逐字回 pristine `B9DC39AEA15670E8F24689FA8D5F8157F2A1FFE310996203670C37DCF4D5643B`（`Copy-Item mut\pinnedleaf.pristin.go →`；★⛔ 碰 `checkout`＝还原一律用本腿自己的 pristine 拷贝＋逐枚哈希对回；还原后 `go vet` rc=0）。
- ⚠ 一枚具名自抓：第一发我误用 `refuseEveryName342r2mut` 顶调用⇒ `declared and not used: pinned`＝四数⛔ 编译＝该轮读数⛔ 计入判据、⛔ 采信；改形后重打才算。＝公共块"一把没咬住的尺"那味在本腿自己手上复现一次、当场认。

## 门禁四数（只在共享工作树，跑前 `unset GOTMPDIR TMPDIR TEMP TMP`，退码＝脚本文件里 `echo RC=$?` 紧跟裸命令）

`gofumpt -l <我这枚两字面路径>` rc=0（空列出＝干净）；`sh scripts/d22scan.sh` rc=0＝clean（具名：本腿曾把一枚 U+26D4 写进字符串字面量＝d22scan 对字符串⛔ 豁免＝入库前已改 ASCII＝笔④，先修在入库侧＝⛔ 拿门禁读数当整改）；`sh scripts/check-path-length-budget.sh` rc=0＝VERDICT GREEN（over-budget=60＝roster=60＝⛔ 新增）；`go vet ./internal/winsec/` 本机 rc=0，★另打交叉档 `GOOS=linux`／`GOOS=darwin`／`GOOS=windows` 各 rc=0（交叉档⛔ 与本机档成对判红绿＝这里只登记"三档都能编译＋vet 干净"，⛔ 用它当任何一格的颜色）。全部四道跑前 `unset GOTMPDIR TMPDIR TEMP TMP`，跑在 HEAD＝`0e487f06`（笔⑤之后）；笔⑥之后票面件改动＝⛔ `.go` 字节⇒四数对他人读数⛔ 变（`227-v2` 那面读的是 `.go`／长度／emoji）。

## 台面纪律对账（派单硬要求逐条）

### 临时根逐枚分类（尺＝`sh/roots-and-faults.sh`，输出＝`logs/`；⛔ 手抄＝直接从日志里 grep 出来的行数）

| 读数 | 实际用的根 | 判读 |
|---|---|---|
| `land/base-det-r1..3`（三枚确定性件短根轮） | `C:\Users\swq\AppData\Local\Temp\…` | ★主读数 |
| `land/base-det-shortroot`（同一对、sh 起的） | `C:\Users\swq\…`（USERPROFILE 直挂＝仍短形⛔ 长 D 卷） | 同形复算＝件三 B/L 两行逐发都在 |
| `leaf-base-r1..3` ⇄ `leaf-land-v2-r1..3`（叶子件两树各 3 发） | `C:\Users\swq\AppData\Local\Temp\…` | ★主读数（派单②③的成对分子） |
| `mutA/mutB-targeted-r1..2`（去牙⇄反形） | 同上短形 | ★主读数 |
| `base/land-full-2`、`-full-3`、`landminus-full-1`（成对四数） | 同上短形（441/448、442/449、435 行） | ★主读数 |
| `land/base-targeted-1`、`base/land-full-1` | `D:\tmp\wisp\20261019\342r2\gotmp`（★长 D 卷＝A1053 那味的高危形） | 只作首轮探索＝其结论已被短根轮取代（件一两树首轮即复现）＝⛔ 拿它们当判据 |

★A1053 族的**建设性检查**（同一把尺顺手扫的）：全部日志里 `mklink /J` 失败行＝0、`cannot construct`（＝⛔ 开发模式下的符号链接构造）行＝0、`--- SKIP`＝0（五枚 full 轮逐枚 0）⇒ 本腿⛔ 一枚"建不出"的假红、也⛔ 一枚靠 SKIP 蒙过去的绿。
★具名一条**别腿仪器的台面故障**（⛔ 我的字节＝⛔ 我改）＝落地 full 轮里 `ancestorunlink_348_windows_test.go:75` 每轮都报 `reparseAt proved 3 flipped leaf cycles across 400 (flipErrs=27491~27835, first "leaf … is not a reparse point after the flip")`＝r1 那枚叶子竞态件在本机仍然空转（＝终裁在册的同一味台损）；★对照＝本腿件四在⛔ 同机器⛔ 同短根下报 `proved flips=52~55`、`flipErrs=9~16`＝裸兄弟名那形建得出来、racer 跑得动＝这条差集是"形"⛔ 是机器＝写进交件供裁腿复核（⛔ 我据此翻任何框）。

- 临时根＝`%LOCALAPPDATA%` 短形（真 junction 建得动＝probe 全程 0 枚"建不出"假红）；仓外树与日志＝`D:\tmp\wisp\20261019\342r2\**`（⛔ 盘符根、⛔ 删除＝只建不删）。
- 取正文⛔ 事后 grep 配错行＝全部 `-v` 落文件再看 `--- FAIL` 上下文；`⛔ 管进 tail`＝退码一律脚本文件＋`echo RC=$?`；★踩过一次＝PowerShell 内联 `$?` 被展开成 `True`（公共块逐字警告过的那枚坑，本腿自己撞在第一发门禁上、换脚本文件重取＝具名）。
- ⛔ `go test ./...`；⛔ 跑 `./cmd/wisp` 整包（＝`227-v2` 持跑者）；⛔ `gofumpt -w` 扫全仓／别腿文件（只我两枚字面路径）；每枚 `.go` 落盘→自验→即时 commit（跨腿耦合＝`227-v2` 的门禁四数吃共享工作树）。
- 摆动名册＝本轮⛔ 撞红需要复量的枚；`342-v2` 点名的两枚 SID 红我按"复量 3 发＋⛔ GOTMPDIR 定向复跑"归因成台损族、⛔ 记成新造红、⛔ 当凭据。

## 本腿⛔ 做／⛔ 宣称的（防过度读数）

⛔ 更新 `winsec.go:118-140`（更新权＝终裁之后）；⛔ 宣布任何一格闭环；⛔ 翻框；⛔ 把件一/件二/件三读成 342 或 348 的 AC 达成；⛔ 把"祖先窗⛔ 存在于落地"读成"祖先支已关"；⛔ 把 `:11-20`/`:158-170` 的机制归因错误改口径（＝终裁在册、修它属产码注释＝另批）；⛔ 用 probe7（⛔ 编译成功过）。

## 交件⑨＝两处口径更正＋散件真损修复（★追加、⛔ 改写历史）

1. **散件里一处真损已修**＝`logs/gates-and-roots.txt` 第二行曾被生成脚本用 `cut -c1-90` 截过＝`cut -c` 数的是**字节**⇒把一枚中文提交主题拦腰截出一枚⛔ 解码字节。★我原先三把尺全都⛔ 报＝`find -size 0`（字节数>0）／码位 perl 尺（带 `:encoding(UTF-8)`＝只 warn⛔ 停、照样给 0）／grep（容忍非法字节）＝"尺全绿⇄件真损"在本腿自己散件上复发＝抓它的是第四把尺 `iconv -f UTF-8 -t UTF-8`。修复＝生成脚本改成整行打印＋该散件重生成；新增双尺件＝`sh/utf8-and-codepoints.sh`→`logs/utf8-and-codepoints.txt`＝逐枚 iconv：正件 11 枚⛔ 损 0 枚、raw 日志 27 枚⛔ 损 0 枚；码位尺逐枚 U+00A7＝0、正控 `HANDOVER.md`＝304（⛔ 与间隔号 U+00B7 混＝同尺并报）。★另具名一枚我自己刚犯的＝同一脚本里 `echo "...`git log -1`..."` 双引号内反引号＝命令替换⇒把整段 git log 泄进散件（已改成单引号并读回生成物头部才放行）。
2. **交件⑧正文一句半错**＝那句"我的新字节全 LF"的尺只打在 `internal/winsec` 两枚 `.go` 上（`git cat-file blob HEAD:<path> | tr -cd CR | wc -c`＝各 0＝这半成立、照旧有效），★⛔ 覆盖散件那 39 枚＝散件里恰恰有一枚 CJK 截断损＝口径改写成"两枚 `.go` LF 尺 0／散件另用 iconv 双尺逐枚验"。
3. **门禁四数在最终 HEAD 重取**（`sh/gates-and-roots.sh` 现跑，输出即那枚重生成散件）＝HEAD＝`70aac8c0`（笔⑧）＝`gofumpt -l 我那两枚字面路径` rc=0／`d22scan` rc=0／长度门 rc=0／`go vet` 本机＋`GOOS=linux`＋`GOOS=darwin` 各 rc=0★且逐枚 `BYTES=0`（＝⛔ 诊断；0 字节日志在"⛔ 0 字节文件"那条规矩里是**工具输出**⛔ 我写的证据件，故具名说明而非入库）；笔⑥以后⛔ `.go` 字节改动⇒对 `227-v2` 那面共享树读数⛔ 变。
4. **枚数照实**＝本腿＝**10 笔**（首 `4f4f8419`⇄笔⑩）、`.go` 新增 2 枚⇄产码改动 0 枚、散件 **39 枚**（36＋笔⑨新增 3：`sh/final-sweep.sh`／`sh/utf8-and-codepoints.sh`／`logs/utf8-and-codepoints.txt`＝`find` 现量、0 字节 0、`.out` 0）、票尾各 1 条（`grep -c 'agent=342-r2'`＝1⇄1）；★整笔交件⛔ "全绿"＝落地 r1 三枚竞态钉照旧红、票 348 `deleted` 分子照旧取⛔ 到、派单①原句形状仍等第二道批准＝rc 全 0⛔ 等于颜色全绿。
5. ★落笔后又抓一枚同族＝那件新入库的 `logs/utf8-and-codepoints.txt` 是在我写这一节**之前**生成的＝它记的 `REPORT_MD_LINES=110` 对当时字节是真的、⇄笔⑨提交进去的 `report.md` 已是 117 行＝同一笔里两枚文件互⛔ 一致（⛔ 帽违规＝两个数都⛔ 越 160＝但"提交后现跑对回"这条我没做到位＝顺序错了）。处理＝★⛔ 手涂那件的历史行（那是仪器输出）、★重跑那把尺量当下提交的字节＝笔⑩（只那一枚文件＋本行）＝本腿在这条规矩上共撞三回（笔⑧枚数、笔⑨散件字节损、笔⑩行数），逐枚具名、⛔ 抹原文。
