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

## 五、摆动名册定因⇄★本落地的连带（两味分开）（★两族分子⛔ 合并、⛔ 相减＝票 348 AC#3 那条；尺＝`sh/27`⇄`sh/30`⇄`sh/32`⇄`sh/33`⇄`sh/35`＝定向两侧各 4 发＋整包两侧各 2 发＋同时段交替 12 发）

**A 族＝`ancestorunlink_348_windows_test.go` 的 400 发 racer 分子行**（`truncated`＝票 342 的号、`deleted`＝票 348 的号）：
`deleted=0`＝★两批共 26 枚分子行都是 0（＝census 第 1 节那句"`deleted` 只可能更小"兑现＝本落地⛔ 把它推高过；本票那枚号⛔ 在这一族里）。`truncated`＝一批里基线 6 发全 0 ⇄ 落地 9 发里 3 发非零，★另一批（同时段交替）11 枚逐发都 0＝⛔ 复现⇒★本腿⛔ 据它盖章（笔⑥ 当场撤回该条），两批的数都摆着交裁腿复量＝详见下面"换档"那两味（第九节第 7⇄8 条具名上报）。
另＝同族 `provedReparse`⇄`flipErrs` 两侧量级**换档**＝★这味随字节走、⛔ 随时段走。尺＝`sh/35-alternating-attribution.sh`＝两侧在同一时段**交替**各打（⛔ 先打完一侧再打另一侧＝那形会把台面负载漂移误当成字节差），`form=leaf-symlink` 那枚逐字＝

| 交替发序 | 树 | `provedReparse` | `flipErrs` | `truncated` | `deleted` |
|---|---|---|---|---|---|
| r1-1 | finalB | 2 | 31366 | 0 | 0 |
| r1-2 | **finalL** | 102 | **16** | 0 | 0 |
| r1-3 | **finalL** | （那发⛔ 打印分子） | — | — | — |
| r1-4 | finalB | 32 | 21815 | 0 | 0 |
| r2-1 | finalB | 98 | 9130 | 0 | 0 |
| r2-2 | **finalL** | 152 | **6** | 0 | 0 |
| r2-3 | finalB | 54 | 20929 | 0 | 0 |
| r3-1 | finalB | 4 | 29885 | 0 | 0 |
| r3-2 | **finalL** | 153 | **14** | 0 | 0 |
| r3-3 | **finalL** | 129 | **8** | 0 | 0 |
| r3-4 | finalB | 3 | 30048 | 0 | 0 |

⇒ 分离**干净到没有交叠**＝基线 6 发 `flipErrs` ∈ 9,130..31,366、落地 4 发 ∈ 6..16（★三个数量级、两侧无交叠、交替打仍如此）。归因〔算的，本腿⛔ 跑计时〕＝乙之前，那一路由在"名字带着攻击者那枚暂挂 link"的世界照样 `Root.Remove(base)`（`Deleteat` 带 `FILE_OPEN_REPARSE_POINT`）⇒**把攻击者的 staged link 删掉**⇒racer 下一发必须重建它＝那两三万枚 `flipErrs` 的来源；乙之后同一性复验⛔ 匹配⇒跳过删除⇒★那枚 link 活着⇒restock ⛔ 必重建⇒`provedReparse` 上到 102..153。⇒ ★这正是机主随拍认下的那味①"复验拒了之后仍留着调用者那枚名字"的下游形状，且它落在⛔ 属本票的世界＝留着的那枚⛔ 可能是攻击者的对象。本腿⛔ 处置（改成"也删 link"＝重述票 118/120-r3 已批的"link 本身⛔ 跟随、只 unlink"那批形＝另批、⛔ 本腿权）、⛔ 借此停手＝具名上报交编排者／`353-v1` 裁。

★同批的 `truncated`＝**⛔ 成立为"上了台阶"，本腿当场撤回那条盖章**＝上面交替的 11 枚分子逐发 `truncated=0`、`deleted=0`，⛔ 复现早前 `sh/32`⇄`sh/33` 那批里落地侧 9 发中 3 发非零（`2`/`1`/`3`，皆 400 发内；基线那批 6 发全 0）。两批冲突时⛔ 能用"先假定新那把坏"裁＝两把都是本腿的尺、只差时段⇄负载，⇒★本腿⛔ 据任何一批说"342 那枚漏得更多了"，只把两批各自的数摆出来＋交裁腿用 `sh/35` 的交替形自行复量。`deleted`＝两批全部分子行都是 0＝census 第 1 节那句"只可能更小"兑现＝本落地⛔ 把它推高过；本票那枚号是 `damage`（第二、三节＝落地侧⛔ 论几发都 0、去牙发 3 摘掉复验立刻回到 3）＝★本票判据自始至终⛔ 读那两族的任何数。

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
7. ★★本落地⇄票 342/348 那两族仪器的耦合＝**两味分开报**（⛔ 藏、⛔ 自行处置，归编排者／`353-v1` 裁）：
   · ★**成立的那一味＝`flipErrs`／`provedReparse` 换档**（量的＝`sh/35` 交替那张表 11 枚分子；算的＝归因，⛔ 跑计时）：基线 6 发 `flipErrs` ∈ 9,130..31,366 ⇄ 落地 4 发 ∈ 6..16＝★三个数量级、两侧无交叠、★同一时段交替打仍如此⇒随字节走、⛔ 随时段漂移。机制＝乙之前那一路由在"名字带着攻击者那枚暂挂 link"的世界照样 `Root.Remove(base)`（`Deleteat` 带 `FILE_OPEN_REPARSE_POINT`）⇒把攻击者的 staged link 删掉⇒racer 下一发必重建它＝那两三万枚 `flipErrs` 的来源；乙之后复验⛔ 匹配⇒跳过删除⇒★那枚 link 活着⇒`provedReparse` 上到 102..153。⇒ ★这正是机主随拍认下的那味①"复验拒了仍留着调用者那枚名字"的下游形状、且落在⛔ 属本票的世界＝留着的那枚⛔ 可能是攻击者的对象；改成"也删 link"＝重述票 118/120-r3 已批那批"link 本身⛔ 跟随、只 unlink"＝另批、⛔ 本腿权⇒★只报不处置。
   · ★**⛔ 成立、本腿当场撤回盖章的那一味＝"`truncated` 敞口上了台阶"**：`sh/32`/`sh/33` 那批里落地侧 9 发有 3 发非零（`2`/`1`/`3`）、基线 6 发全 0；同一时段交替的 `sh/35` 那 11 枚分子逐发 `truncated=0`＝⛔ 复现⇒★本腿⛔ 据它说"342 那枚漏得更多了"，只把两批的数都摆出来，交裁腿用 `sh/35` 的交替形自行复量。`deleted`＝两批共 26 枚分子★全 0＝census 第 1 节"只可能更小"兑现＝⛔ 任何一族的数被本落地推高过。
8. ★那两族仪器件（`ancestorunlink_348_*`⇄`ancestoropen_342_*`⇄`leafrace_342r2_*`）都不在本腿写面＝一字⛔ 动、⛔ 为变绿改任何断言；★票面 AC#2 那句"⛔ 借 `truncated`／`deleted` 任何一枚"本腿⛔ 改口径、也⛔ 需要改＝本票判据（`damage`⇄`entered`＋正控）自始至终⛔ 读那两族的任何数（第二／三节可逐字复核）＝所以第七节这两味＝★连带读数上报、⛔ 判据的一部分。


## 十、纪律自证（三把尺）

名册三把尺＝锚定 `^票 353 落地腿 353-r1` ⇄ 非锚定子串 ⇄ 单次导出 `%s`＋内存 `index(...)==1`（尺＝`sh/34-final-roster.sh`，可逐字重跑；★枚数随 HEAD 漂⇒本件不写死枚数、只认逐枚 `%s` 起头＝公共块"名册尺"那条第五形）。写面越界尺＝`sh/29` 的逐枚前缀计数形（★⛔ 用 sh/28 那把坏尺＝它把目录前缀钉了行尾⇒永⛔ 匹配⇒整表误判"越界"，本腿自抓已换）＝`OUT_OF_RANGE=0`；`.go` 名册去重＝6 枚全在 `internal/winsec/`。★入库 blob CR＝0（尺＝`git cat-file blob | tr -cd '\r' | wc -c`；★正控＝同一条尺在 `.scratch/wisp/issues/300-*.md` 报 `CR=1`＝尺没瞎、0 是真 0；钝尺对照＝对工作树 `grep -c $'\r'` 冒报 288＝台账 `A1066③` 同一条坑，本腿也踩过一次、已弃用）。票 353 框终量＝★未勾 3／已勾 2（`AC#0`⇄`AC#1`＝编排者翻）＝本腿一枚⛔ 翻；本腿各笔触碰票面＝0 笔（票尾唯一追加名额已被 `353-a1` 占⇒本腿话落本件⇄交件正文）；票名／`Status:`／`-done` 一字⛔ 动。⛔ push；⛔ `add -A`/`.`/目录形／从共享 index 反推；`add`⇄`commit` 两处逐枚显式 pathspec；⛔ amend/reset/rebase/stash/checkout ./clean；临时件只建不删＝`D:\tmp\wisp\20261020\353r1\{finalB,finalL,mut-final-*,incoming,base,logs,gocache-*}`（★⛔ 盘符根）。正文件行数＝本件⇄`readings-ac2-baseline.md` 都⛔ 超帽 160（帽⛔ 适用 `*_test.go`：仪器件 195 行照实报）。★`internal/models` 那枚脏文件⇄`cmd/wisp`＝别腿在飞字节＝本腿⛔ add⛔ commit⛔ 采信。
