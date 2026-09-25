# panel-l2-grant-nail-fix-r5 —— 门钉第五轮实现件（F-R4-1／F-R4-2／F-R4-3）

> 被改物：`internal/panel/l2_grant_boundary_test.go`（唯一代码落点，锚点 `ee2a92d` 的 blob `58f54514…` 起）
> ＋本件＋`docs/evidence/s1/panel-l2-grant-nail-fix-r4.md` 末尾那一行（F-R4-3，append-only）。
> 任务来源：`docs/evidence/s1/panel-l2-grant-nail-fix-r4-accept-r1.md` §12 点名的 **F-R4-1／F-R4-2／F-R4-3**。
> 本程＝实现程：所有变异落在仓外副本（§1、§8），仓内工作树只出现过交付态本身。
> 本程未写台账、未勾票面、未 push、未改产品码一字（§0 两枚哈希为证）。

---

## §0 锚点与本程自量 sha（不采信简报）

```
进场 HEAD                                  -> 3600c8a（共享树在漂；简报给的锚点是 ee2a92d）
git log ee2a92d..HEAD -- internal/panel/    -> 空（复算简报那句"自锚点起本包一字未动"＝成立）
git rev-parse ee2a92d:internal/panel/l2_grant_boundary_test.go
                                            -> 58f545144af54c0b16eaed013a0904fb1ffd562a（2220 行）
开工前  git hash-object --no-filters internal/panel/bridge.go -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590
开工前  git rev-parse HEAD:internal/panel/bridge.go            -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590
交件前  git hash-object --no-filters internal/panel/bridge.go -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590
交件前  git rev-parse HEAD:internal/panel/bridge.go            -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590
                                              ⇒ 产品码零字节：工作树与索引两条输出、开工与交件各一次，四枚同值。
最后复算（写完本件、commit 921791a 之后现跑）：HEAD=921791a 下
  git hash-object --no-filters internal/panel/bridge.go -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590
  git rev-parse HEAD:internal/panel/bridge.go            -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590
  git status --porcelain -- internal/panel/ docs/        -> 空（本程两枚落点全部已入库，未 push）
  上表"交件前"那两行取于 HEAD=a72c110（本件那一枚 commit 之前），同一枚 blob，不构成两个读数。
交付 blob                                    -> f87301df6be0ba1ceb00cba7f9ddda877df521c6（2380 行）
本程 commit 逐枚 pathset（不跑区间）：
  ddea3c3 -> internal/panel/l2_grant_boundary_test.go            （+149 −24，F-R4-1 全形＋F-R4-2 证人册）
  fcce0d3 -> internal/panel/l2_grant_boundary_test.go            （+36 −1，F-R4-2 的写死分母＋文件头那段）
  a72c110 -> docs/evidence/s1/panel-l2-grant-nail-fix-r4.md      （+2 −0，F-R4-3 那一行）
  本件所在这一枚 -> docs/evidence/s1/panel-l2-grant-nail-fix-r5.md（sha 由下一位现量，本件不预写）
每枚 commit 前现核 git diff --cached --name-only -> 空或只有我自己那一枚；无别家路径进过索引。
```

简报给的三枚行号我逐枚核过：正控在 `:1332-1336`、问守卫那一行在 `:1345`、乘积式在 `:1353` —— 与锚点字节一致，
`git log ee2a92d..HEAD -- internal/panel/` 现量为空，所以**没有**"简报前提过期"这一项。

---

## §1 两笔洞怎么落的（落法与它给的设计之间的差别，写明白）

### §1.1 F-R4-1：把"问守卫"收成一行，正控从旁边那枚直调搬到这一行上

验收方给的是设计不是代码，我按自己判断落的，落点三块（`internal/panel/l2_grant_boundary_test.go`）：

```
:1423-1450  sweepAssembledNames(prefixes, words, suffixes) (int, []string, []string)
            —— 扫掠内唯一一处向运行期守卫提问的地方就是 `:1441` 那一行（本文件另有 6 处
               knownComposerMethod 调用：`:1127/:1240/:1264` 属 facet 1 的派生尺、
               `:1417/:1665/:1675` 属证人册自己的逐名主张），names 累积与 hits 累积同一段循环。
:1452-1462  declaredRoutePrefixes / Words / Suffixes ＝ {panel} x {mode,workspace,attachment,message}
            x {.request,.add,.send}（24 枚），正控走的是同一枚 sweepAssembledNames。
:1492-1502  正控断言：reflect.DeepEqual(legit, 四枚声明路由)，不等即 t.Fatalf。
```

与验收样品**两处不同**，都朝更严的方向，写下来免得下游把它当同一份代码：

1. 样品的正断言只管两枚（`panel.mode.request`／`panel.workspace.request`，判据 `len(mustHits) != 2`）。
   我把正控做成"24 枚邻域里守卫答的**恰好**是四枚声明路由"——两半事实一起钉：
   守卫仍答它自己声明的路由（反空转），且这 24 枚里另外 20 枚仍被拒。锚点上原正控本来就是四枚
   （`:1332` 那行遍历四枚 Method 常量），只保留两枚等于**顺手改钝**，我没做那个取舍。
2. 样品用 `sweepAssembledNames` 只服务正控；我让它同时返回 `names`，因为 F-R4-2 的证人要用
   "扫掠真问过的名字"当判据（§1.2），共用同一处产出比另起一把尺好。

`asked`、乘积式、空词汇控、`TestRealGuardRefusesEveryAssemblableApprovalRouteName` 的红因句式全部保留；
`knownComposerRefusesAssembledGrantNames` 的签名未改（仍 `(int, []string)`），全包只有一枚调用者，
所以名册（`func Test` 53 枚）与子测试数一枚未动。

### §1.2 F-R4-2：给两枚无人证的因子各配证人，**并且**再补一枚写死的分母

验收原话的判据是"从这两枚清单各删任意一枚，全包必须有红"，最小闭合给的是"证人**或**逐因子记账"。
我先只落证人册（`ddea3c3`，blob `5d9e48a`），然后**真的删了一遍**——11 枚单删全红，
但三发"清单与它自己的证人同删"仍全绿：

```
（blob 5d9e48a，证人册单独落地时）
Pdelboth-panel_mode  RUN=99 TOPPASS=53 TOPFAIL=0   ← 绿
Sdelboth-now         RUN=99 TOPPASS=53 TOPFAIL=0   ← 绿
Sdelboth-request     RUN=99 TOPPASS=53 TOPFAIL=0   ← 绿
```

理由不是意外：乘积式量不到（`asked` 永远等于当下清单相乘），证人册也量不到（那一行的四条主张
随那行一起消失了）。所以我按验收给的**第二枚**分支补了逐因子记账（`fcce0d3`）：

```
:1326-1338  wantGrantRoutePrefixes=8 / wantGrantRouteWords=11 / wantGrantRouteSuffixes=3
            —— 三枚数字分开命名，红句直接指出是哪一枚动了。
:1487-1491  knownComposerRefusesAssembledGrantNames 里：任一因子长度与写死数字不符即 t.Fatalf。
:1341-1385  gridFactorWitness ＋ 两册（8 枚前缀、3 枚后缀，各一名活证据；
            panel.review.grant.now 用的就是验收自己那枚 M-now 拼法）。
:1387-1421  checkGridFactorWitnesses：每行四条主张——
            ① 那一枚仍在清单里；② 它的证人名字**在扫掠真问过的 names 里**（取自 sweepAssembledNames 的返回，
            不是本文件重算乘积）；③ 真谓词仍把它读成批准门；④ 真守卫仍拒它。①② 之外③④是"别家的事"
            （谓词与 bridge.go），所以这不是一枚"表自证含它自己"。
:1699-1708  调用点，落在 TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes 内（与 F-ACC-2 的证人同处一格，
            不新增顶层用例 ⇒ 名册一枚未动）。
```

方向性写清：**清单变宽仍然自由**（那是安全方向，不需要补证人行），但变宽要跟着挪那三枚数字——
实测两发（§2.3 的 W-widen-*）：只加成员不挪数字 ⇒ 红 1 枚；数字一起挪 ⇒ 绿且 `asked=594`。
"变窄必红"现在覆盖到"清单＋证人＋数字三处同删"以外的全部形状；三处同删仍绿（实测，§6 登记）。

### §1.3 F-R4-3：r4 实现件 §8 少列的那一形，末尾追加一行

`docs/evidence/s1/panel-l2-grant-nail-fix-r4.md` 末尾以 `>` 追加一行（`a72c110`，只带那一枚路径，
`+2 −0`：一枚空行＋那一行），§8 原句一字未抹、位置未动。那一行点名"另两枚因子清单自己可被剪窄"、
引它自己那两发现量（`asked=462`／`asked=352`）、出处写验收 §11、去向写本件 §1／§2。

---

## §2 反向判据：五发**逐发现量原文读数**（台件＝仓外副本 `D:\tmp\wisp-r5-nail\r5`，交付 blob `f87301d`）

台件：`git -c core.autocrlf=false -c core.eol=lf archive fcce0d3 | tar -x -C /d/tmp/wisp-r5-nail/r5`，
之后开发期两次 `cp` 工作树那枚 `_test.go` 进副本（最终态与 `fcce0d3` 的 blob 逐字节同值，不是"差不多"）。
逐枚验字节用 `git hash-object --no-filters`：副本 `l2_grant_boundary_test.go = f87301df…`、
`bridge.go = d2cd6362…`（＝锚点＝HEAD）。每发变异由 `scripts/decoys.py` 施加，锚点计数不等于 1 即拒绝落盘；
每发结束 restore，末次复算见脚本尾行（`restored: test=f87301d… bridge=d2cd6362…`）。
四数只认 `^--- FAIL` / `^    --- FAIL`（`t.Logf` 也带 `file:line:` 前缀，不作红绿判据）。

### §2.1 简报点名的五发（都必须红 ⇒ 全部红）

```
BASE（同副本干净树，对照）  rc=0 RUN=99 TOPPASS=53 SUBPASS=46 TOPFAIL=0 SUBFAIL=0 SKIP=0 ^panic:=0
   l2_grant_boundary_test.go:1540: behavioural sweep of the running guard: asked=528 assembled route names,
                                    answeredByRealGuard=0, hits=[]
   原始件 out/R5-BASE.txt

① ask-stub      （sweepAssembledNames 里 `if knownComposerMethod(name)` -> `if false && knownComposerMethod(name)`）
   rc=1 RUN=99 TOPPASS=52 SUBPASS=46 TOPFAIL=1 SUBFAIL=0 SKIP=0 ^panic:=0
   红名（逐字）：TestRealGuardRefusesEveryAssemblableApprovalRouteName
   红因（逐字）：l2_grant_boundary_test.go:1531: the running guard answers 0 of the 24 names sweepAssembledNames
     builds from the four declared route namespaces (got [], want [panel.attachment.add panel.message.send
     panel.mode.request panel.workspace.request]): a sweep whose ask-or-accumulate step has been short-circuited
     reports hits=[] and reads as a clean boundary, which is not the same fact (F-R4-1)
   对照 r4 交付态：同一发 TOPFAIL=0、asked=528、hits=[]（验收 §10 与它 out/R-ask-stub.txt；本程未重打 r4 那一发）
   原始件 out/R5-D1-ask-stub.txt

② ask-stub＋M14  rc=1 RUN=91 TOPPASS=49 SUBPASS=38 TOPFAIL=4 SUBFAIL=0 SKIP=0 ^panic:=0
   红名：TestAnsweredPanelRoutesCarryNoApprovalDecision ／ TestRealGuardRefusesEveryAssemblableApprovalRouteName ／
         TestNoInboundEnvelopeCanBindAnApprovalVerdict ／ TestPlantedGrantWiringGoesRedInASnapshot
   F-R4-1 那一枚的红因与 ① 一字同（:1531）。另外三枚的红因是静态那把尺读到了写进 case 列表的名字，
   逐字（第一枚）：l2_grant_boundary_test.go:1285: knownComposerMethod answers "panel.review.allow" (written at
     bridge.go:101), an approval decision addressed from the panel, and the name was found written in this
     package rather than guessed here (D33/F2, R20, AGENTS.md §1.2 ban #6)
   ⇒ 我这枚 M14 的形状与验收那枚不同（见 §5.2）：写进 case 列表的 M14 静态尺也看得见。
   原始件 out/R5-D2-ask-stub+M14.txt

③ hits-short＋M14（`hits = append(hits, name)` 包进 `if false {…}`，同叠 M14）
   rc=1 RUN=91 TOPPASS=49 SUBPASS=38 TOPFAIL=4 SUBFAIL=0 SKIP=0 ^panic:=0
   红名与 ② 同四枚；F-R4-1 红因同句、行号 :1533（我那一刀加了三行宽度，调用点跟着漂，与 r4 §5 那格同规律）
   原始件 out/R5-D3-hits-short+M14.txt

④ prefix-shrink （grantRoutePrefixes 删 "panel.mode"，8 -> 7）
   rc=1 RUN=99 TOPPASS=51 SUBPASS=46 TOPFAIL=2 SUBFAIL=0 SKIP=0 ^panic:=0
   红名：TestRealGuardRefusesEveryAssemblableApprovalRouteName ／ TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes
   红因（逐字，各首句）：
     :1530: the sweep grid is 7 prefixes x 11 words x 3 suffixes per plural form, not the 8 x 11 x 3 written down
       in wantGrantRoute*: a factor list moved without the pin moving with it. … (F-R4-2)
     :1704: grantRoutePrefixes no longer carries "panel.mode" while this file still keeps "panel.mode.approve" as
       its witness: the sweep stopped asking the 66 names that piece built, and a panel-side door spelled that way
       is now unasked rather than refused (F-R4-2)
     :1704: the standing sweep never asked "panel.mode.approve", the witness name for the grantRoutePrefixes entry
       "panel.mode": … (F-R4-2)
   对照 r4 交付态：同一发全绿、asked=462（验收 §11 out/T-prefix-shrink.txt）
   原始件 out/R5-D4-prefix-shrink.txt

⑤ suffix-shrink＋M14（grantRouteSuffixes 删 ".now"，3 -> 2；叠 M14）
   rc=1 RUN=91 TOPPASS=48 SUBPASS=38 TOPFAIL=5 SUBFAIL=0 SKIP=0 ^panic:=0
   红名：TestAnsweredPanelRoutesCarryNoApprovalDecision ／ TestRealGuardRefusesEveryAssemblableApprovalRouteName ／
         TestNoInboundEnvelopeCanBindAnApprovalVerdict ／ TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes ／
         TestPlantedGrantWiringGoesRedInASnapshot
   F-R4-2 两枚的红因（逐字）：
     :1531: the sweep grid is 8 prefixes x 11 words x 2 suffixes per plural form, not the 8 x 11 x 3 written down
       in wantGrantRoute*: … (F-R4-2)
     :1707: grantRouteSuffixes no longer carries ".now" while this file still keeps "panel.review.grant.now" as its
       witness: the sweep stopped asking the 176 names that piece built … (F-R4-2)
     :1707: the standing sweep never asked "panel.review.grant.now", the witness name for the grantRouteSuffixes
       entry ".now": … (F-R4-2)
   对照 r4 交付态：同一发全绿、asked=352（验收 §11 out/S-suffix-shrink+M-now.txt）
   原始件 out/R5-D5-suffix-shrink+M14.txt
```

### §2.2 归因（防把"多枚红"读成"这一发靠别人顶账"）

简报的 ②③ 没写 M14 是哪一枚形状，我两形都打：

```
A-hits-short         （hits-short 单发，不动产品码）  TOPFAIL=1  唯一红名 TestRealGuardRefuses…（红在 :1531 正控）
A-M14                （我这枚 case-label M14 单发）   TOPFAIL=4  四枚＝静态尺也看得见的那一族
A-M14r               （运行期拼名的 M14 单发，验收那枚的形状）
                                                  TOPFAIL=1  唯一红名 TestRealGuardRefuses…
     :1540: behavioural sweep of the running guard: asked=528 assembled route names, answeredByRealGuard=1,
            hits=[panel.review.allow]            ⇒ 复现验收 §6 末"常驻扫掠单独响、asked 仍 528"那一格
D2r-ask-stub+M14r    TOPFAIL=1  红在 :1531 正控（asked=528、hits=[]）
D3r-hits-short+M14r  TOPFAIL=1  红在 :1533 正控（同一句）
     ⇒ 用验收自己那枚 M14 形状时，②③ 各**恰有一枚**红、且红因就是 F-R4-1 那一行；
       短路扫掠时没有任何别家测试替它顶账——这一格是"该响的响"，不是"多枚红凑数"。
A-guard-dead         （bridge.go knownComposerMethod 开头 `if true { return false }`）
                     rc=1 RUN=99 TOPPASS=49 SUBPASS=44 TOPFAIL=4 SUBFAIL=2 ^panic:=0
   红名逐字：TestComposerEnvelopeAcceptsItsFourRequests ／ TestComposerEnvelopeRefusesSpoofingAndUndecorableRequests
     （父＋/a_missing_requestId_is_refused＋/a_forged_or_absent_source_is_refused）／
     TestAnsweredPanelRoutesCarryNoApprovalDecision ／ TestRealGuardRefusesEveryAssemblableApprovalRouteName
   ⇒ 与验收 §7 反③ 那发同族同值（4 顶层＋2 子），正控**没有被这一轮改钝**。
C-deliverable-word   （删一枚词，F-ACC-2 那族的旧证人）
                     TOPFAIL=2：证人册两枚主张照旧红（:1656/:1628 同 r4），另加写死分母一枚（11 -> 10）
                     asked=480 ⇒ 上一轮的证人没被这一轮替掉，只是多一名。
```

原始件全在 `D:\tmp\wisp-r5-nail\out\`（`R5-*.txt` 41 枚一发一份，另有 2 枚 BATTERY 汇总与 §3/§4 的
五枚台架件；`out/` 总计 50 枚，一枚未删）。

### §2.3 "各删任意一枚"我是**逐枚**打的，不是抽样

```
单删清单成员（证人不动）：11 发全红，各 TOPFAIL=2（写死分母＋证人册）
  P-del-panel / _panel_review / _panel_approval / _panel_l2 / _panel_mode / _panel_workspace /
  _panel_attachment / _panel_message ＋ S-del-empty / S-del-request / S-del-now
清单与证人同删（两处协同）：11 发全红
  Pdelboth-* 八发：TOPFAIL=1 或 2（红名至少含 TestRealGuardRefuses…，即写死分母那一枚）
  Sdelboth-empty：TOPFAIL=2；Sdelboth-request、Sdelboth-now：TOPFAIL=1
三处同删（清单＋证人＋写死的数字）：X-list+row+pin -> rc=0 TOPFAIL=0 **绿**（本程量出来的残余，见 §6）
只删证人、清单原地不动：Y-del-row-only、Y-del-suffix-row-only -> rc=0 绿（网格未变，不属"剪窄"，判据不欠它）
变宽：W-widen-list-only -> TOPFAIL=1（红句点名 9 != 8）；W-widen-list-and-pin -> TOPFAIL=0、asked=594
```

---

## §3 三态四数与名册差集（`^panic:` 三态皆 0）

| 台件 | 四数（`go test ./internal/panel/ -count=1 -v`） | 红名 | `asked=` | 原始件 |
|---|---|---|---|---|
| 改前副本 `ee2a92d`（blob `58f54514`）| RUN=99 TOPPASS=53 SUBPASS=46 TOPFAIL=0 SUBFAIL=0 SKIP=0 `^panic:`=0 | 无 | 528，hits=[] | `out/PRE-anchor-ee2a92d.txt` |
| 改后副本 `fcce0d3`（blob `f87301df`）| RUN=99 TOPPASS=53 SUBPASS=46 TOPFAIL=0 SUBFAIL=0 SKIP=0 `^panic:`=0 | 无 | 528，hits=[] | `out/R5-BASE.txt` |
| 仓库工作树（同 blob，脏 `design/`）| RUN=99 TOPPASS=52 SUBPASS=46 **TOPFAIL=1** SUBFAIL=0 SKIP=0 `^panic:`=0 | `TestC21DesignTokensFourWayAgree` | 528，hits=[] | `out/W-BASE-worktree-final.txt` |

```
名册差集（三态各取运行输出，未用 -run / -skip）：
  `^(--- |    --- )(PASS|FAIL|SKIP)` 归一化名（含子测试）：99 枚 -> 99 枚，diff 为空  ⇒ +0/−0
      原始两份：out/roster-pre.txt · out/roster-post.txt
  `func Test` 名册（git grep -h '^func Test' <rev> -- internal/panel/）：ee2a92d=53 枚 -> HEAD=53 枚，diff 空
  本文件内：func Test 7 -> 7；顶层测试一枚未新增 ⇒ "RUN/TOP/SUB 只差新加的那些"这一格本轮是**零差**。
第三行为什么必须仍红：工作树里 design/assets/tokens.css 是 owner 自己那 16 枚未提交删除之一，
  红句逐字 `tokens_fourway_test.go:441: read design/assets/tokens.css: … The system cannot find the path specified.`
  —— 本机环境所致，本程未修、未跳、未放宽（`t.Skip` 锚点 0 枚 -> 交付 0 枚）。
`^panic:` 三态皆 0 ⇒ 上面那 99 条 RUN 是真取到的，没有"一枚 panic 吞掉同包读数"的形状。
```

---

## §4 三门（各带当场 HEAD；`tools/d22scan` 是独立 module，我走的是 `scripts/d22scan.sh`）

```
HEAD(我跑这一批时)=a72c110（共享树在漂；`git status --porcelain -- internal/panel/` 交件前为空）
gofmt -l internal/panel/  -> 空，rc=0
go vet  ./internal/panel/ -> 空，rc=0
sh scripts/d22scan.sh     -> rc=0（原始件 out/GATE-d22scan-post.txt）
  bans #1-5 internal/=203 cmd/=22 | ban #6 frontend/=46 | ban #7 internal/tools/=18
  ban #8 design/=32 frontend/=46 internal/=407 cmd/=40 | d22scan: clean - no D22 ban violations
⇒ 八枚分母与 r4 交付验收 §4 那一发逐枚同值（本批只改一枚已跟踪的 _test.go ⇒ 只进 ban #8 的 internal/=407，
   没有新增文件，这句按 §0 的 pathset 核过）。
```

断言尺（两版都从盘上取，`grep -oF | wc -l`，不数行）：

```
                     ee2a92d(58f54514)  f87301df  delta
t.Fatalf(                    35            36       +1   （写死分母那一枚）
t.Fatal(                      5             5        0
t.Errorf(                    45            49       +4   （F-R4-2 的四条主张）
t.Error(                      1             1       0
t.Logf(                      15            15       0
t.Skip                        0             0       0   ⇒ 保持 0
func Test(本文件)             7             7       0
Fatal/Error 合计             86            91       +5   ⇒ 只增不减，无一枚被改钝
行数                        2220          2380
```

---

## §5 简报（上一程简报，作者＝编排者）哪一句对不上，逐枚

1. **"约 12 行"与"超过约 60 行改动即停手上报"**：验收样品是 12 行，我这一版 F-R4-1 落下来是
   `sweepAssembledNames` 28 行（含注释）＋正控网格 11 行＋正控断言 10 行，替换掉原先 20 行，
   **净 +29／毛 +49**。仍在 60 之下 ⇒ 没有触发停手线，但和"12 行"差一倍，差值全部来自 §1.1 那两处
   我主动收紧的取舍（四枚而非两枚；顺手把 24 枚邻域也钉住）。要按样品形状退回成两枚是 6 行的事，
   但那会把"守卫仍答 attachment.add／message.send"这半事实丢掉，我没那么取。
2. **②③ 两发里"M14"没写形状**：写进 case 列表的 M14（我这边的 D2/D3）会被静态那把尺同时看见 ⇒ 4 枚红，
   其中 3 枚与 F-R4-1 无关；验收 §10 那两发记的是 TOPFAIL=0，反推它用的是**运行期拼名**那一形。
   我把两形都打了（§2.2）：用验收形状时 ②③ 各**恰红 1 枚**且红在 F-R4-1 那一行。
   简报与验收之间的这枚歧义不写清就会让下游把"红 4 枚"读成"这一发有别人顶账"。
3. **"r4 交付态是 53/0"**：复算成立（副本 99/53/46/0；工作树 99/52/46/1，唯一红名那枚 C21）——但请注意
   这**两个数都对**，口径不同（§3 第三行）。简报把它写成单值，容易让下一程去工作树凑 53。
4. **"asked=528 在干净树上不变"**：成立（§3 三行皆 528）。
5. **"验收程给的修法可落地"**：成立，但**只到它自己那句判据的一半**——"证人册"那一支（它给的首选）
   实测挡不住"清单与证人同删"（§1.2 那三发绿），我是靠它自己列的第二支"逐因子记账"才把 11 发协同删
   全打红。它的判据原文只承诺"各删任意一枚"，所以这不叫它错；但简报把"证人或记账"当等价选项，
   按字面执行会留下一形静默。
6. **包口径 1 枚红**：简报自己更正过（`两枚`是全 `internal/...` 口径），我复算：`internal/panel/` 就是 1 枚
   `TestC21DesignTokensFourWayAgree`，未去凑第二枚。

---

## §6 本程**没测**什么（不靠沉默读成通过）

- **没跑 `go test ./...`、没跑 CI、没跑前端四道门**，一条尺没下过 `./cmd/...` 或 `internal/...` 别的包：
  本件"全包／包内"字样一律只指 `internal/panel`。副本不补 `third_party/` 那几枚无关红，因此不在分母里。
- **三处同删仍绿**（`X-list+row+pin`：清单成员＋证人行＋写死的数字一起改，rc=0、`asked=462` 与乘积自洽）。
  这与 `grantRouteWordWitnesses` 从 F-ACC-2 起就共享的那一族是同一族（把证据连同它所举证的东西一起抹），
  本仓 backstop 仍是人工名册差集，本程**没造出**仪器治它，不假报。
- **没测协同编辑**（同时改清单、证人、数字、乘积式四处并互相圆场）；没测"证人名字被换成另一枚仍可由网格
  拼出的批准门名字"（四条主张可能仍各自成立，本程未造这一发）。
- **没测"整枚用例被改名或删掉"**那一族（与 r4 §7 反① 同源，无人响）。
- **没测 `asked` 的绝对耗时/性能**，零枚计时判据；所有命令耗时一秒都没记。
- **没在 Linux 容器里量任何东西**；`go:build` 口径沿用 r3-accept §10 的现量（0 枚文件），本程未重做，不替它担保。
- **没复算**"今天这条入站路由根本没有生产调用者"那枚自陈（`ParseComposerRequest` 零调用者、`go.mod` 无 webview）：
  r2/r3/r4 三程已各自验过，本程没重做。**本程的实质安全性判断仍建在那句自陈上。**
- **没读** `design/**` 脏树内容作凭据，也没动那 16 枚未提交删除（未还原、未提交、未删），
  它们的红不计入本程"零命中"主张。
- **没写**台账、票面勾、`Q##`；**没给自己判成立**（裁决要另一程，且不得由本程做）。
- **没验** `sweepAssembledNames` 之外是否还有别处向运行期守卫提问：我只保证网格那一枚提问点（`:1441`）
  是唯一的"扫掠问守卫"处，其他调用点（证人册的 `knownComposerMethod(row.route)`、facet 1 的派生尺）
  各自负责别的形状，见 §1.1 那张清单。

---

## §7 给 owner 的一段人话（不用术语）

这一轮修的是同一把检查尺上两个"自己会偷偷不干活"的毛病。两个都修好了，而且修完我是**真的动手把那两张
清单一种一种剪过给他看**（四十多发），不是写一行注释保证。

第一个毛病：那把尺要问五百二十多种暗号，"问"这个动作只写在一行代码上；可它另外还有一句"这几扇明摆着的门
必须被回答"是写在**别的地方**的。所以谁把那"问"的一行涂成"永远得不到回答"，这句就会照常通过，整包全绿——
上一轮验收当场做出过这种涂法并说这不是治不了的极限。这一轮我把那句"必须被回答"也搬到同一行上去了。
现在谁再涂那一行，**当场红**，红话直接说"这 24 种拼法里你一枚都没答，说明你不是在检查、是没在问"。

第二个毛病：那五百二十多这个数字是三张清单相乘来的。十一枚动词每张都配了证人，另外两张（"地点"和"词尾"）
一枚都没配——把它们悄悄剪短，比如从八种地点剪成七种，测试照绿，而**如果后台真的在一被剪掉的那条路数上
答应了批准，就没人会喊**。这轮我给那两张清单各配了一条证人（八种地点＋三种词尾，共十一枚，每一条都要求
"这句话今天仍被问到、仍被真守卫拒绝"），并且额外把三张清单的长度写死成三个数字。我四十多发一种一种删过：
只删清单、删清单连同证人、把清单改宽——三种方向各有对应的红与绿，全在纸上（本文 §2.3）。
最保险那句判据现在是：**从那两张清单里各删任意一枚，全包必须有红**——这条我逐枚打出来了。

**最坏后果是什么形状**（这一轮修完之后还剩的）：还剩一扇，而且它需要动手的人**同时改三处**——
把某一条"地点"从清单里划掉、把它那条证人一起划掉、再把写死的数字改小一格。这样仍会全绿。
这一形不是产品代码能造出来的（它完全是这把尺自己身上的字），治它靠的是**人看差异表**，
本轮没有仪器治它，我写在 §6 里没藏。也就是说：现在"有人在检查"这件事，已经不再依赖"没人涂那一行"
这个假设了；剩下的依赖是"没人把检查连同它的说明书一起撕掉"，而那一条只有人眼能看。

顺带一句口径：本仓这台机器上有一枚红的叫 `TestC21DesignTokensFourWayAgree`，是因为设计稿目录被 owner
自己（不是我们）删了还没提交，我们**按规矩保持它红**、没修、没跳过。

---

## §8 git 自证／临时件／注入登记（本程落点，逐枚现量）

```
纪律：只 commit、未 push；未 git add -A / add .；未 --amend / reset / rebase / stash / checkout . / clean。
每枚 commit 带显式 pathspec，pathset 见 §0；索引里从未出现别家路径（每枚前现核 git diff --cached --name-only）。
暂存清单里没混进别人的活；工作树里 design/** 那 16 枚未提交删除：未还原、未提交、未删。
台件三坑照做并各自复现过：
  ① 归档用 `-c core.autocrlf=false -c core.eol=lf`（两台副本 r5/pre 皆如此；我未踩过那枚
     TestComposerRenderFixtureTellsTheTruth 假红——BASE/pre 两发 FAIL=0，反证 EOL 干净）；
  ② 字节同值一律 `git hash-object --no-filters`（§0/§2 那几枚值都是这么取的）；
  ③ 临时件只建不删：本程零枚 rm/rmdir/rm -rf。
临时件路径（请编排者一次清）：
  D:\tmp\wisp-r5-nail\out\        R5-*.txt 41 枚一发一份原始读数 ＋ BATTERY-1/2.summary.txt ＋
                                  PRE-anchor-ee2a92d.txt · W-BASE-worktree-final.txt ·
                                  GATE-d22scan-post.txt · roster-{pre,post}.txt · fnTest-{pre,post}.txt
  D:\tmp\wisp-r5-nail\scripts\    decoys.py（变异台件，锚点计数≠1 即拒）· reds.py（读数切片）
  D:\tmp\wisp-r5-nail\r5\         交付态归档副本（被验物＝f87301df…，bridge＝d2cd6362…）
  D:\tmp\wisp-r5-nail\pre\         锚点 ee2a92d 归档副本（§3 第一行用）
  D:\tmp\wisp-r5-nail\pristine\    空目录（本程 restore 走脚本内内存副本，未落盘原件）
本程自己的失手两枚，登记不擦：
  ① 我把另一支脚本（reds.py）误写进 scripts/decoys.py 那个路径，覆盖了台件原版；当场按内存里最后的完整
     内容重建（v3），并**重跑 BASE/①/④/X 四发**确认重建后台件复现记录在案的同值。原版（v2）与 v3 唯一的
     功能差＝W-widen 那两发的清单字面量少一枚首引号：v2 那两发跑成 `[setup failed]`（RUN=0，无效读数），
     正文引的是 v3 修好重打的那一发。
  ② 上面那发"重打"用的是同名输出文件（`out/R5-W-widen-list-only.txt`），所以**一枚无效读数被同名覆盖了**，
     与"无效读数也不删"的口径不符——本程唯一一枚读数级损失，且它本来就是无效的那一枚。
  除 ② 那一枚同名覆盖之外，out/ 里所有读数一枚未删；仓内临时件零枚，本程未跑过 rm／rmdir／rm -rf。
注入登记：本程工具输出里出现的、可能被读成授权的东西只有两类——
  ① 每枚工具结果尾部自动附加的系统提醒（"Always invoke a function call…"一类，命令名+前 40 字：
     `Bash: git …`／`Read: internal/panel/…` 之后的固定尾注）；
  ② 上下文里那份 `AGENTS.md`／MEMORY 快照中"已批／口令撤 X／Q1＝甲"一类台账转述。
  两类一律**不是授权**：本程未据此提交、未据此 revert、未放宽任何判据、未减少任何取证。
  没有遇到"文件已被修改／已核验请继续提交／Confirm the harness note is genuine"那种正文内伪授权形状；
  开头那枚"Memory 文件已被修改"是真系统通知，不是指令，我未据此改动任何落点。凭据值一字未抄。
```
