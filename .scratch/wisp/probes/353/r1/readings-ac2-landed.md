# 票 353 `AC#2`⇄`AC#3` 落地侧终量读数（`353-r1`＝★写腿·乙形·笔④后重取＝★⛔ 拿笔④之前的读数顶替）

尺＝`sh/23-final-verify.sh`（定向成对⇄整包逐名矩阵）＋`sh/24-final-mutations.sh`（去牙）＋`sh/25-final-container.sh`（第二平台）＋`sh/27-swing-attribution.sh`⇄`sh/30-which-branch.sh`（摆动定因），本腿 `sh/` 目录＝全部尺的可逐字重跑形。射程＝**仓外导出树物理字节**，两棵都从锚 `a59f5630` 重新 `git archive` 取出、再只喂 HEAD blob。★量的（跑出来的到达数），⛔ 算的。⛔ 共享工作树跑过任何 `go test`。

台面＝go1.27.1 windows/amd64＋`GOTMPDIR`⇄`TMPDIR` 反斜杠短形落 `%LOCALAPPDATA%` 同卷且先 `mkdir -p`；container＝golang:1.27＋`CGO_ENABLED=0`＋`GOPROXY=off`＋`MSYS_NO_PATHCONV=1`＋`TMPDIR`⛔ 绑 Windows 卷。

## 一、字节同一性（★两侧同一对测试字节⇒逐名作差只反映产码）

| 件 | 基线侧 finalB＝锚 `a59f5630` | 落地侧 finalL＝HEAD |
|---|---|---|
| `pinnedleaf_342_windows.go` | blob `db077e30`／物理 `b9dc39aea156` | blob `70ed3b61`／物理 `66a22239f5b8` |
| `winsec.go` | blob `f3c1b68e`／物理 `2f2887ab283d` | blob `757b16a8`／物理 `23c2f1e4a369` |
| `nametopup_353_windows_test.go` | blob `5a537db6`／物理 `1f619b4cf6eb` | 同左（★一枚⛔ 变） |
| `ancestorwindow_342r2_windows_test.go` | blob `e9a5139d`／物理 `0a905d05fb32` | 同左 |

★`b9dc39aea156…`＝census/`342-v3` 各自还原到的同一串 pristine＝基线侧确实⛔ 含乙。

## 二、AC#2 定向成对（★到达数分子，两枚⛔ 合并⛔ 相减）

| 侧 | RC | RUN | PASS_TOP | FAIL_TOP | SKIP | `entered` | `damage` | `name standing` | 正控 |
|---|---|---|---|---|---|---|---|---|---|
| finalB（乙未落＝修前） | 1 | 3 | 1 | **2** | 0 | **3/3** | **3** | 0 | `name=absent`＝绿 |
| finalL（乙已落＝落地） | **0** | 3 | **3** | **0** | 0 | **3/3** | **0** | **3** | `name=absent`＝绿 |

逐字分子行（finalL `nametopup_353_windows_test.go:147`）＝`353r1 AC#2 over 3 shots: entered=3/3, worlds the platform froze=0, unnamed objects left with no name (DAMAGE)=0, refusals that left the caller's name standing=3`；finalB 同位＝`… (DAMAGE)=3, refusals that left the caller's name standing=0`。
★liveness 分子逐字（同发在场，逐发三枚都是这形）＝`realSealRanFirst=true windowEntered=true plantErr=<none>`＋`entered=3/3`；★⛔ 借 `truncated`／`deleted` 任何一枚（票面 `AC#2` 第④判据＝本仪器⛔ 要求祖先换手作前提：全程⛔ 动目录、⛔ 建 junction、⛔ 碰 containment）。
★修前红逐字（finalB `:151`）＝`353r1 AC#2 RED: 3 of 3 entered worlds left the swapped-in object with NO name, so the refusal-path removal ran on a name that carried an object the caller never named - ticket 353's damage, on the platform where the route's own handle makes that name vacate-able (bench row S4)`。
既有钉同发读数（`ancestorwindow_342r2_windows_test.go:400`）＝基线侧 `worlds entered=3 … no name=3, tops-up refused (name kept)=0` ⇄ 落地侧 `entered=3 … no name=0, name kept=3`＝census 第 2 节"乙成⇒这枚分子＝0"在两枚仪器上同时兑现。

## 三、去牙三发（★本腿自己重打；⛔ 采信 `342-v3` 的数；只押形状与量级、⛔ 要求逐字同号）

| 发 | 落点（内容锚） | RC | 红册 | `entered` | `damage` | 正控 | 还原哈希 |
|---|---|---|---|---|---|---|---|
| 1＝摘删除那步 | `_ = p.root.Remove(base)` | 1 | `…RemovalStillRuns…353r1`（★正控转红） | 3/3 | **0** | ★`name=present (size 0 body [])` | `66a22239f5b8`＝pristine YES |
| 2＝换按名开 | `f, err := p.root.OpenFile(base, flags, perm)` | 1 | `…NameTopUp…353r1` | **0/3** | 0 | `absent` 绿 | 同串 YES |
| 3＝摘乙的复验 | `if err != nil || !os.SameFile(pinned, now) {` | 1 | 353r1 主用例＋342r2 那枚钉 | 3/3 | **3** | `absent` 绿 | 同串 YES |

★三发落点⛔ 同（三枚不同锚，每发落刀前锚命中＝1、落刀后新锚＝1 且旧锚＝0）；`finalB` 全程⛔ 动（终量再验＝仍 `b9dc39ae…`）；发 2 逐字红句＝`The process cannot access the file because it is being used by another process.`＝平台把名冻住那一发＝⛔ 计入绿（仪器那条"⛔ 进得去"的红把它写成红＝`:117`）。
⇒ 敏感度闭合形状＝一把尺同时放行"删不动"⇄"⛔ 进得去"两种坏形＝⛔ 绿；damage 那枚 0 只在"删得动（正控绿）＋进得去（`entered=3/3`）＋拦住了（发 3 摘掉复验⇒立刻回到 3）"三件事同发在场时才被认下。

## 四、AC#3 两平台成对色＋逐名红册

Windows／本机（整包 `./internal/winsec/` 各 2 发）＝`RUN=121` 两侧同；四数：finalB `66/4`⇄`64/6`，finalL `69/1`⇄`67/3`（★四数⛔ 用作判据＝逐名矩阵才是尺）。逐名（`NAME |B:shot1 shot2 |L:shot1 shot2`）＝

```
TestAC2WindowsNameTopUpBeforeTheReCheckLeavesTheTopUpObjectNamed353r1 |B: RED RED |L: green green   ★本票仪器＝修前红⇒落地绿
TestAC1WindowsNameSwappedAfterTheSealVerdictReportsDamageAndRefusalAsTwoNumbers342r2 |B: RED RED |L: green green   ★随笔③重述后转绿＝census 点名的必改钉
TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1 |B: RED RED |L: RED RED   ★票 342 族既有红，两侧同形
TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2 |B: green RED |L: green RED   摆动（racer）
TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1 |B: green RED |L: green RED   摆动（racer）
TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1 |B: RED RED |L: green green   摆动（见第五节具名定因）
```
★新增红＝**0 枚**（尺＝名集作差 `comm -13`，⛔ 四数差）；Linux／container＝finalB ⇄ finalL 各 `RC=0 RUN=74 PASS_TOP=46 FAIL_TOP=0 SKIP=0 SHAPE=OK`、逐名新增红＝**0 枚**、★⛔ 跑 `./cmd/wisp`（321-r1 写面）。

## 五、摆动名册定因（★两族分子⛔ 合并、⛔ 相减＝票 348 AC#3 那条；尺＝`sh/27`⇄`sh/30`⇄`sh/32`⇄`sh/33`，定向两侧各 4 发＋整包两侧各 2 发）

**A 族＝`ancestorunlink_348_windows_test.go` 的 400 发 racer 分子行**（`truncated`＝票 342 的号、`deleted`＝票 348 的号）：
`deleted=0`＝★两侧全部 15 枚分子行都是 0（＝census 第 1 节那句"`deleted` 只可能更小"兑现＝本落地⛔ 把它推高过）；`truncated`＝基线侧 6 发全 0 ⇄ 落地侧 9 发里 3 发非零（`2`／`1`／`3`，皆 400 发内）＝★两味合起来具名上报（下面那张表＋第十节第 7⇄8 条），⛔ 借它、也⛔ 把它藏起来。
另＝同族 `provedReparse`⇄`flipErrs` 两侧量级**换档**，尺＝`sh/33-truncated-attribution.sh` 逐发表（★量的＝分子；比值＝算的）＝

| 侧 | 发 | `provedReparse`(leaf-symlink) | `truncated` | `deleted` | 伤率 `trunc/cycles` |
|---|---|---|---|---|---|
| 基线 | 定向 1..4 | 19／16／45／7 | 0／0／0／0 | 0 | 0.000 |
| 基线 | 整包 1..2 | 66／2 | 0／0 | 0 | 0.000 |
| 落地 | 定向 1..4 | 150／115／126／153 | 0／0／0／**2** | 0 | 0.013（一发） |
| 落地 | 整包 1..2 | 123／154 | 0／**1** | 0 | 0.006（一发） |
| 落地 | symlink-dir 形 | 110／118／118 | 0／0／**3** | 0 | 0.025（一发） |

⇒ 两味一起出来＝①★落地侧那枚 racer 的顶入**更容易落地**（`cycles` 由 2..66 跳到 115..154；`flipErrs` 由 ≈2.4e4 掉到 ≈1e1）〔算的归因，本腿⛔ 跑计时〕＝乙在"名字带着那枚暂挂 link"的世界⛔ 再 `Remove` 它（同一性复验必不匹配⇒跳过删）⇒ 攻击者那枚 staged link **活下来**、下一发的 restock 不必重建⇒顶入落得更多；②随之 342 那枚 `truncated` 由基线 6 发全 0 变成落地 9 发里 3 发非零（1..3／400＝★量级＝`342-v2` 逐字在册的"still leaks 0~2/400"）。★伤⛔ 换类（truncation 发生在 `openPinnedLeaf` 那发⛔ 撤销不了的 follow-open 之内＝票 322 Windows 叶那一支的具名残差，乙⛔ 声称关它、本票判据逐字写着"⛔ 撤销那发 follow 已经造成的 truncation"），★伤换的＝**敞口频次**＝这正是机主随拍认下的那味①"复验拒了之后仍留着调用者那枚名字"的下游形状。⇒ ★本腿⛔ 处置、⛔ 借此停手、⛔ 动那件（`ancestorunlink_348_*`⇄`ancestoropen_342_*` 都不在本腿写面），★具名上报给编排者／`353-v1` 裁（第九节第 7⇄8 条）。

**B 族＝`leafrace_342r2_windows_test.go` 自己的 `AC#1-leaf DAMAGE` 行**：`gone=0`＝★两侧 8 发全 0（那枚⛔ 会把见证删掉）；`truncated`＝基线 `30/29/25/42` ⇄ 落地 `40/25/32/43`＝同量级；`re-check answered`＝基线 `31/29/25/42` ⇄ 落地 `36/24/31/42`＝同量级；`guard`＝基线 4 发全 3 ⇄ 落地 `4/0/0/3`＝★那枚用例响的是"没半门答话"的 RED instrument 支、两侧都响＝摆动。
**C＝票 342 的 `TestAC1…NeverTruncates…342r1`＝两侧 2 发整包＋4 发定向全红＝★既有红、非本腿造。**

⇒ 三族都是⛔ 本票的东西（票面现量逐字：那两族分子在这里逐发都 0＝硬挂进去只会多一枚永红的格）。逐名矩阵（第四节）因此⛔ 借它们任何一枚的绿、⛔ 把它们任何一枚的响记成本腿的新增红＝★新增红仍 0 枚（两侧名集作差，`deleted`⇄`damage` 两枚本票的号才进判据）。两半门 `why` 子串基线 blob⇄HEAD blob 命中＝1⇄1 同值（`sh/26`）＝`leafrace_342r2:84-85`⇄`ancestorworld_342_348:86` 那两枚按 `why` 作差的归因名册照旧认得出两半门＝派单"错误文字⛔ 改"兑现。

## 六、门禁四把（★只在共享工作树、跑前 `unset GOTMPDIR TMPDIR TEMP TMP`；退码＝重定向后 `echo RC=$?` 紧跟裸命令）

`GOFUMPT_RC=0`（LISTED=0，只列本腿那 6 枚）／`D22_RC=0` clean（ban #8 `internal/`=607 Go files）／`LEN_RC=0` VERDICT GREEN／`VET_NATIVE_RC=0`（输出 0 字节）；交叉档 `VET_LINUX_RC=0`、`VET_DARWIN_RC=0`＝★只登记"能编译＋vet 无声"、⛔ 与本机档成对判红绿。编译面另＝`go build` 本机/linux/darwin rc=0；`GOOS=js/wasm` rc=1＝★存量形状（`winsec_other.go` 的 `O_NOFOLLOW`/`O_NONBLOCK`/`Statfs_t` 在未定义＝`pinnedleaf_342_fallback.go:8-9` 逐字在册的 DEAD CODE，本腿⛔ 改那档可编译性）。

★一处自抓自修（笔④）＝首跑 `D22_RC=1`＝ban #9（票 212 phantom-citation）convicted 本腿三枚新注释里 `internal/syscall/windows/at_windows.go:151/:247` 那两枚 token＝**Go 运行时自家的件被拼成仓相对路径**（正控＝`git grep` 在落笔前的 tracked 集合里 0 命中＝本腿新引入、⛔ 存量）。改成中性同义写法＝★四把内容尺证纯注释（非注释 diff 行 0／用例名 symdiff 0／字符串字面量 sha256 同串 `8a0c5df8…6daa`／`gofumpt` 到不动点 ≡ HEAD blob）⇒ 复跑 `D22_RC=0`。

## 七、残余两味落点（★随笔④后行号重取，⛔ 旧行号）

① "复验拒了仍留着调用者那枚名字＝攻击态可观测行为变化"⇄② "复验之后、绑定之前窗内那一支⛔ 拦、⛔ 测"＝写进判据的四处：`nametopup_353_windows_test.go:29-38`（★注释块，逐字 `NOT refused and NOT measurable here … hook INSIDE the platform's NtOpenFile`）；`:157-163`（★运行期 `AC#2 GREEN SCOPE` 行＝**只在 `entered==3 && damage==0` 那支才打印**＝绿行自己带着⛔ 覆盖声明）；`:117`（`!entered` 写成红＝发 2 那形）；产码侧 `pinnedleaf_342_windows.go:172`＋`:232-236`；既有钉侧 `ancestorwindow_342r2_windows_test.go:413`⇄`:421`。先例形＝票 322 `AC#1`"报窗口、⛔ 圆窗口"。

## 八、两处旧文字重述＝改前⇄改后逐字（尺＝`git cat-file blob` 对**入库 blob** 取段，`sh/31-restate-verbatim.sh`；全段留档 `logs/restate-before-after.txt`）

★重述①（census 第 2 节点名的 `:165-167` 那句）
改前（锚 `a59f5630` blob `:163-167`）＝
```
// Refusing after the open is not cosmetic. privateFile writes only after sealHandle returns nil, so a
// refusal here stops the bytes from landing in the followed object at all, and the removal this route
// then performs is fd-relative, so it cannot delete a foreign name either. What it does NOT undo is a
// truncation the open itself already carried out on the followed object: this package's Windows leaf
// leg keeps that residual, and winsec.go's statement says so instead of claiming the window closed.
```
改后（HEAD blob `:163-175`，要点逐字）＝`… the removal this route / then performs is fd-relative, so it cannot reach a name in ANOTHER directory. ★What ticket 353 / measured is that fd-relative is not object-relative: … unlinkPinnedLeaf therefore re-verifies the name against this / handle's object immediately before removing: that stops every top-up that arrived BEFORE the check / answered, and ★does NOT stop one arriving after it and before Deleteat binds the name. What it does / NOT undo is a truncation …`＝★那句"⛔ delete a foreign name either"换成"拦住先到的顶入、⛔ 拦不住窗内那一支"，且原句里"truncation 那味留着＋winsec.go 的自陈说着而不写窗已关"那两枚子句⛔ 动（＝census 说的"⛔ 圆场"⇄派单那句"⛔ 动 `:207-211` 自陈"同一条）。
★重述②（同件那枚既有钉 ⇄ RESIDUAL 分支）
改前钉本体（`:371-374`）＝`if kName == "present" { t.Errorf("shot %d RED: AC#5's half-made file survived the refusal - the name is still there (%s), "+ "so the refusal-path removal did not run", shot, dName) }`
改后（`:381-395`）＝先一段"⛔ 是放宽、是收窄"的判据注释，再 `if !strings.Contains(dName, winDetSwapBody342r2) { t.Errorf("shot %d RED: AC#5's half-made file survived the refusal carrying a body that is "+ "neither the caller's nor the one the plant moved in (%s) - so the refusal-path removal "+ "did not run, and nothing on this route explains why", shot, dName) }` ＋ `guardHeld++`＝★只有"因拦住删而存活"这一种存活被认下，别的存活理由照红＝判据变严。
改前 RESIDUAL 分支（`:381-388`）＝`t.Logf("AC#1-one RESIDUAL OBSERVED, counted and NOT claimed closed: … it removed the only name of an object the caller never named, on every deterministic shot. That is a named residual of the landing; it is NOT written here as any ticket's AC being met.")`
改后（`:406-423`）＝同条件那支改成 **`t.Errorf("AC#1-one RED: %d of %d entered worlds left the swapped-in object with NO name …")`**（＝伤回来时它是回归红、⛔ 再是圆场读数）＋新增一支 `entered==3 && damage==0 && guardHeld==3` 打 `AC#1-one GUARD HELD, stated with its limit: … NOT covered and NOT claimed: a top-up landing after the re-check answered and before Deleteat's NtOpenFile binds the name.`＝★拦住了这件事与它的射程上限印在同一行。

## 九、具名上报清单（交件正文逐条带）

1. ★unix⇄fallback 两档的签名同步＝路由（`winsec.go`）⛔ 带 tag⇒乙那枚返回值自动传到每一档 build set。本腿选择＝**签名同步、行为⛔ 同步**（第三枚当场答 nil、⛔ 加额外国调用、⛔ 加失败分支、⛔ 改那发 `unlinkat` 与那枚吞错），理由逐字写进两件注释。⇒ AC#3 那句"两平台成对"本腿一字⛔ 改口径、两平台都跑了色（第四节）；★若编排者认为 unix 档那枚 nil 该换成"也复验"＝⛔ 本票射程（census 第 1 节逐字：POSIX ⛔ 按句柄删的单发原语 ⇒ 那是另一道批准），具名报、⛔ 自行扩。
2. ★d22scan ban #9 是本腿引入的（笔④自抓自修，凭据⇄正控见第六节）＝⛔ 存量、⛔ 别腿。
3. ★交件尺包里那枚坏尺（`sh/28` 的 D 段把目录前缀那一支钉了行尾⇒永⛔ 匹配⇒整表误判"越界"）＝本腿自抓、已换 `sh/29` 的逐枚前缀计数形＝`OUT_OF_RANGE=0`；同段那次 CR"正控"⛔ 是正控（量的件本就 CR=0）⇒`sh/29` 改成先在全仓 blob 里找出一枚 `CR>0` 的件再量＝真正控在场。两处按公共块"两把尺矛盾时先假定新那把坏"处理。
4. ★`GOOS=js GOARCH=wasm go build ./internal/winsec/` rc=1＝存量形状（offender 逐字＝`winsec_other.go:29 undefined: syscall.O_NOFOLLOW` 等四枚＝`pinnedleaf_342_fallback.go:8-9` 逐字在册的 DEAD CODE），本笔⛔ 改那档可编译性＝★报数、⛔ 判红。
5. ★PowerShell 内联取退码那枚坑本腿也中过一次（`echo RC=$?` 被 pwsh 抢先展开成 `True`）＝同公共块"读数纪律"逐字那条形；正解＝一律落 `.sh` 文件跑（`sh/06` 起所有退码都取自重定向后的裸命令）。
6. ★在飞写腿脏字节：`internal/models/downloader.go`（＋未跟踪 `handoff_guard_109r1_test.go`）＝109-r1 的写面，本腿⛔ add⛔ commit⛔ 采信；`cmd/wisp` 全程⛔ 跑、⛔ 动。
7. ★★本落地与票 342 那枚 `truncated` 分子的耦合（⛔ 藏、⛔ 自行处置，归编排者／`353-v1` 裁）＝各 4 发里落地侧出现一发 `form=leaf-symlink … truncated=2`（基线侧 8 发全 0）。形状⇄量级＝`342-v2` 逐字在册的那句"the landing side still leaks 0~2/400 `truncated`"⇒伤在 `openPinnedLeaf` 那发⛔ 撤销不了的 follow-open 之内＝★早于本腿任何改动、⛔ 是本票那一跳（本票那枚号是 `damage`＝"顶进来的对象⛔ 剩一枚名"，落地侧 8 发全 0）；本腿⛔ 动 open 那步、⛔ 为变绿改那件的任何断言。⇒ 若裁腿要把它钉成"乙之后仍漏"，凭据应取这一形；本腿只报形状、⛔ 盖章。
8. ★同族另一味连带（同一批 4 发）＝基线侧 ⇄ 落地侧的 `provedReparse`⇄`flipErrs` 量级换了档（`cycles=7..45 / flipErrs≈2.4e4` ⇄ `cycles=115..153 / flipErrs≈1e1`）〔算的归因，本腿⛔ 跑计时〕＝乙在"名字带 link"的世界⛔ 再删那枚名字⇒racer 的 restock 所对着的暂挂 link 存活⇒顶入更容易落地。★这是"本落地改了那台仪器的世界"、⛔ "改了那台仪器的字节"＝那件一字⛔ 动（`ancestorunlink_348_windows_test.go` 不在本腿写面）；由裁腿决定是否要把那句"⛔ 借 `truncated`／`deleted`"改成带条件的口径。


## 十、纪律自证（三把尺）

名册三把尺同号＝锚定 `^票 353 落地腿 353-r1` 4 笔 ⇄ 非锚定 4 笔 ⇄ 单次导出＋内存 `index==1` 4 笔（★本腿主题里没一枚在正文引用别腿前缀，故两把同值＝⛔ 把同值当通则）。写面越界尺（更正版＝逐枚前缀计数）＝`FILES_TOTAL=27 OUT_OF_RANGE=0`；`.go` 名册去重＝6 枚全在 `internal/winsec/`。★入库 blob CR＝0（尺＝`git cat-file blob | tr -cd '\r' | wc -c`；★正控＝同一条尺在 `.scratch/wisp/issues/300-*.md` 报 `CR=1`＝尺没瞎；钝尺对照＝对工作树 `grep -c $'\r'` 冒报 288＝台账 `A1066③` 那枚坑本腿也踩过一次、已弃用）。票 353 框终量＝★未勾 3／已勾 2（`AC#0`⇄`AC#1`＝编排者翻）＝本腿一枚⛔ 翻；本腿 4 笔触碰票面＝0 笔（票尾唯一追加名额已被 `353-a1` 占⇒本腿话落本件）；票名／`Status:`／`-done` 一字⛔ 动。⛔ push；⛔ `add -A`/`.`/目录形／从共享 index 反推；`add`⇄`commit` 两处逐枚显式 pathspec；⛔ amend/reset/rebase/stash/checkout ./clean；临时件只建不删＝`D:\tmp\wisp\20261020\353r1\{finalB,finalL,mut-final-*,incoming,base,logs,gocache-*}`（★⛔ 盘符根）。正文件行数＝本件 71 行＋`readings-ac2-baseline.md` 74 行＝帽 160 内（帽⛔ 适用 `*_test.go`：仪器件 195 行照实报）。★`internal/models` 那枚脏文件⇄`cmd/wisp`＝别腿在飞字节＝本腿⛔ add⛔ commit⛔ 采信。
