# panel-l2-grant-nail-fix-r5-accept-r1 —— 非实现者对"门钉 r5"这一批的对抗验收

> 被验物：`internal/panel/l2_grant_boundary_test.go` @ 锚点 **`aeba6ff`**（blob `f87301df6be0ba1ceb00cba7f9ddda877df521c6`，2380 行）
> ＋自述件 `docs/evidence/s1/panel-l2-grant-nail-fix-r5.md`（§0–§8）。
> 任务来源：`docs/evidence/s1/panel-l2-grant-nail-fix-r4-accept-r1.md` §12 点名的 **F-R4-1／F-R4-2／F-R4-3**。
> 本程＝**验收程**：被验物一字未改，只 commit 本件；不勾票面／台账；不 push。
> 本程全部变异落在**仓外副本**（§1），原始读数在 `D:\tmp\wisp-r5acc1\out\`（§14 列路径）。

---

## §0 锚点、口径与"读数出在哪台件"——先把这条写死，因为工作树此刻在别人手里

**简报的硬要求**：另一程正在给 `internal/panel/` 加快照泵（会动 `composer.go`／`bridge.go`）⇒ 所有读数取锚点归档副本，
脏工作树的读数单独标注、且**不作任何一格的判据**。本程执行到位，且这条不是形式：我 15:5x 建副本、
16:0x 复量时工作树已经长成这样（逐字 `git status --porcelain`，只列与本包相关的）：

```
建副本时（15:56）  ?? 无；internal/panel/ 干净
复量时（16:04）    M  cmd/wisp/run.go
                   ?? cmd/wisp/panel_pump.go
                   ?? internal/agent/approval/pending_read.go
                   ?? internal/panel/pump.go
                   ?? internal/panel/pump_test.go
⇒ 本程正文里每一枚"四数／红名／asked="都出自 D:\tmp\wisp-r5acc1\{anchor,pre,witnessonly}（§1）；
   工作树口径的读数只有 §8 表里那一行，标着"工作树"，并且它不参与任何一格的裁。
```

```
进场 HEAD = aeba6ff（＝锚点本身）；交件前复量 HEAD = ddf8b6f（共享树在漂，只有 docs 两枚与泵那一程的未提交件）
git rev-parse aeba6ff:internal/panel/l2_grant_boundary_test.go
   -> f87301df6be0ba1ceb00cba7f9ddda877df521c6（2380 行）＝ fcce0d3 那一枚，与实现件 §0 同值   ✅
git rev-parse aeba6ff:internal/panel/bridge.go   -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590
git rev-parse 88eab34:internal/panel/bridge.go   -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590  ⇒ 产品码零字节（跨 r2..r5）
git rev-parse aeba6ff:internal/panel/composer.go = ee2a92d 同值 -> e444d5d117a073c393d9b43f7e3e98e8ede61755
git diff --numstat ee2a92d aeba6ff -- internal/panel/  ->  185   25   internal/panel/l2_grant_boundary_test.go（本包只有这一枚）
git diff --name-only ddea3c3 aeba6ff                  ->  本件之外只有 docs/**（⇒ 锚点相对 ddea3c3 的代码差＝那三枚写死的数）
三枚 commit 的 pathset（逐枚 git show --numstat，不跑区间）：
   ddea3c3  149  24  internal/panel/l2_grant_boundary_test.go
   fcce0d3   36   1  internal/panel/l2_grant_boundary_test.go
   a72c110    2   0  docs/evidence/s1/panel-l2-grant-nail-fix-r4.md
```

**F-R4-3 那一枚文档账我单独核过**：`git diff ee2a92d aeba6ff -- <那枚 r4 件>` 的输出是 `449a450,451`（纯追加），
`## §8` 标题仍在且只有一枚，追加的那一行以 `>` 起头、内容点名"另两枚因子清单自己可被剪窄"、引它自己那两枚发现量
（`asked=462`／`asked=352`，两枚我本程都独立复算到了，见 §3/§4）、出处与去向都对得上 ⇒ **F-R4-3 结清**。
⚠ 一处**过期风险**记在这里（不改判）：那一行写的是"各删任意一枚时全包全绿"，那是 **r4 交付态**的现状；
r5 之后同一发已红（§4 的 K1 十一发）。该句末尾指了 `panel-l2-grant-nail-fix-r5.md §1、§2` 作去向，
所以不是错句，但下一程若只摘那一行会读到过期现状——本程在 §10 把它列成一句口径提醒。

---

## §1 台件自证（ⓐ）：三条仪器坑本程都自己走了一遍，不复用实现者的脚本

```
正路：git -c core.autocrlf=false -c core.eol=lf archive aeba6ff | tar -x -C /d/tmp/wisp-r5acc1/anchor
逐枚比字节（本程自己的尺）：git ls-tree -r aeba6ff（1174 枚）对比 git hash-object --no-filters <副本>
   -> match=1168  mismatch=6
   那 6 枚逐字：scripts/build.ps1 · scripts/dev/ball-cycle.ps1 · scripts/fetch-deps.ps1
               scripts/sign-models.ps1 · scripts/slo-check.ps1 · scripts/spike/run.ps1
   成因＝.gitattributes:2 逐字 `*.ps1 text eol=crlf`（显式属性压过 core.eol=lf）⇒ 仓自己要求的 CRLF，不是污染；
   grep -rn "\.ps1" internal/panel/ -> 0 枚 ⇒ 本包没有一枚测试读 .ps1 ⇒ 这 6 枚不进本程分母。
被验物逐枚（--no-filters）：
   l2_grant_boundary_test.go -> f87301df6be0ba1ceb00cba7f9ddda877df521c6 ＝ 锚点 ✅
   bridge.go                 -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590 ＝ 锚点 ✅
   composer.go               -> e444d5d117a073c393d9b43f7e3e98e8ede61755 ＝ 锚点 ✅
   frontend/fixtures/composer-states.html -> ab39390c…（6179 字节，CR 计数 0）＝ 锚点 ✅

坑①（本程独立复现，第三程）：只给 core.autocrlf=false、不给 core.eol=lf ——
   git -c core.autocrlf=false archive aeba6ff frontend/fixtures/composer-states.html | tar -x
   -> 680b908b…，6187 字节，CR 计数 8  ⇒ 与锚点不同值。本程副本用了正路，未踩那枚假红。
坑②：比字节一律带 --no-filters（上面每一枚都是这么取的；不带它会把 680b908b 读成 ab39390c）。
坑③：临时件只建不删——本程零枚 rm/rmdir/rm -rf（路径清单在 §14，交编排者一次清）。
```

变异台件：`D:\tmp\wisp-r5acc1\scripts\mut.py`（本程自己写的，未复用实现者的 `decoys.py`）。
每发施加前先把两枚文件从 `pristine\<副本>\` 复原；每一刀的 `old` 串在当前文本里**计数不等于 1 就抛 REFUSED、半行不落**
（本程真实拒过 0 次，但同一把尺挡下过我自己两枚坏切法，见 §14 失手登记）；每发结束再复原一次。
末次全树复算（84 发跑完之后）：`match=1168 mismatch=6`，那 6 枚仍是上面点名的 `*.ps1` ⇒ **副本没有被任何一发变异留脏**。
读数总枚数：`out/` 里 **99 枚**＝**84 发 `go test`**（含三态对照与两台的 BASE）＋ 三门 2 枚 ＋ 字节比对/名册/清单等辅助件；
其中 **4 枚无效**（`H01–H04`，RUN=0 编译不过，本程第一版"拆正控"的切法漏了变量声明），一律留档不删、正文改引修正后的 `H05–H07`。
判红绿只认 `^--- FAIL` / `^    --- FAIL`（`t.Logf` 也带 `file:line:` 前缀，不作判据）。

---

## §2 简报点名的五发反向判据：本程逐枚重打，五发全红，逐格对得上

台件＝`anchor`（交付态副本）。"它记的"＝`panel-l2-grant-nail-fix-r5.md` §2.1 那五行。

| 发（本程编号） | 本程只改哪一处 | 四数（`go test ./internal/panel/ -count=1 -v`） | 红名 | 红因首句（逐字，行号是本程那一刀的落点） | 原始件 |
|---|---|---|---|---|---|
| BASE（对照） | 无 | rc=0 RUN=99 TOPPASS=53 SUBPASS=46 TOPFAIL=0 SUBFAIL=0 SKIP=0 `^panic:`=0 | 无 | `:1540 … asked=528 assembled route names, answeredByRealGuard=0, hits=[]` | `out/BASE-anchor.txt` |
| **① ask-stub** | `:1441` -> `if false && knownComposerMethod(name)` | rc=1 99/**52**/46/**1**/0/0 | `TestRealGuardRefusesEveryAssemblableApprovalRouteName` | `:1531 the running guard answers 0 of the 24 names sweepAssembledNames builds from the four declared route namespaces (got [], want [panel.attachment.add panel.message.send panel.mode.request panel.workspace.request]) … (F-R4-1)` | `out/A01-ask-stub.txt` |
| **② ask-stub＋M14** | ① ＋ `bridge.go:97` 运行期答 `panel.review.allow` | rc=1 99/52/46/**1**/0/0 | 同① | 与①一字同（`:1531`） | `out/A02-ask-stub-M14r.txt` |
| **③ hits-short＋M14** | `:1442` 的累积包进 `if false` ＋ M14r | rc=1 99/52/46/**1**/0/0 | 同① | 同一句，行号 `:1533`（本程那一刀净 +2 行） | `out/A03-hits-short-M14r.txt` |
| **④ prefix-shrink** | `grantRoutePrefixes` 去 `"panel.mode"`（8->7） | rc=1 99/**51**/46/**2**/0/0 | `TestRealGuardRefuses…` ＋ `TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes` | `:1531 the sweep grid is 7 prefixes x 11 words x 3 suffixes … not the 8 x 11 x 3 written down in wantGrantRoute*` ＋ `:1705` 两条证人主张（清单不再含它／常驻扫掠没问过 `panel.mode.approve`） | `out/A04-prefix-shrink.txt` |
| **⑤ suffix-shrink＋M14** | 去 `".now"`（3->2）＋ M14 **标签式** | rc=1 RUN=**91** TOPPASS=48 SUBPASS=38 TOPFAIL=**5** SUBFAIL=0 SKIP=0 | 它列的五枚逐字复现 | `:1531` 分母枚 ＋ `:1707` 两条 `.now` 证人枚 ＋ 标签式 M14 的静态三枚 | `out/A05-suffix-shrink-M14label.txt` |

**归因（防把"多枚红"读成"这一发靠别家顶账"）**：本程把 M14 的**两形都打了**——

```
G01 M14r 单发（运行期拼名，whole name 从不成串写进包里）    TOPFAIL=1  只有常驻扫掠响，hits=[panel.review.allow]
G02 M14label 单发（写进 case 列表）                        TOPFAIL=4  静态那把尺看得见的那一族
G03 ask-stub ＋ M14label   TOPFAIL=4  ⇒ 逐枚比对红因文本：四枚里**恰一枚**带 "(F-R4-1)" 字样（:1531），
                                        另三枚的红因是 :1244/:1285/:1301 那一族（静态尺读到 case 列表里的名字）
G04 hits-short ＋ M14label TOPFAIL=4  同形
G05 hits-short 单发        TOPFAIL=1  红在 :1533（正控）
G06 guard-dead（bridge.go 的 knownComposerMethod 开头 `if true { return false }`）
             TOPFAIL=4 SUBFAIL=2 ⇒ 与 r4 交付验收 §7 反③、实现件 §2.2 那格同族同值 ⇒ **正控这一轮没被改钝**
```

**裁（五发＋归因）**：实现件 §2.1 那五行的红枚数（1/4/4/2/5）、红名、红因、`RUN=91` 那两发的顶层数变化，
本程**逐格复算成立**。写清哪一发对哪一枚，免得下游把两形混成一种：② ③ 的"红 4 枚"是本程用**标签式** M14 复现的（G03/G04），
换成**运行期拼名**那一形时各**恰红 1 枚**（A02/A03，红因逐字就是 F-R4-1 那句）；① 是 A01、④ 是 A04、⑤ 是 A05。
`RUN=91` 的来路也核过：53 枚顶层一枚不少（48 PASS＋5 FAIL），少的是 8 枚子测试（父用例 Fatalf 后不跑）⇒ 48＋5＋38＝91。
一处台件差异写明不作矛盾：它 ③ 记 `:1533` 并自述"那一刀加了三行宽度"，本程那一刀净 +2 行、同样落 `:1533`——
调用点是 `t.Helper()` 报出来的**外层调用行**，净 +2 才该落 `:1533`，它正文里"三行"是**写了几行**不是**净漂几行**。

---

## §3 ⓑ 它自称"只立证人那一支实测不够"：本程在**没有写死分母的副本**上重打——推翻未遂，这句成立

这是本轮最值钱的一条，按简报要求单独一台件打：`witnessonly` ＝ `git archive ddea3c3`（**不是**我构造的"去掉分母"版本，
是它真实提交过的那一版），blob `5d9e48aa508da0914a84479bfd95d20f29f1ddeb`，`grep -c wantGrantRoute` 在该文件里 **＝ 0**
⇒ "只有证人、没有写死的数"这句的前提是**盘上事实**，不是我推出来的状态。

| 发 | 台件 | 四数 | `asked=` | 判 | 原始件 |
|---|---|---|---|---|---|
| WO-B01 `Pdelboth-panel_mode`（清单去 `panel.mode` ＋ 去它的证人行） | witnessonly | rc=0 99/**53**/46/**0**/0/0 | **462** | **绿** | `out/WO-B01-Pdelboth.txt` |
| WO-B02 `Sdelboth-now`（去 `".now"` ＋ 去它的行） | witnessonly | rc=0 99/53/46/0/0/0 | **352** | **绿** | `out/WO-B02-Sdelboth-now.txt` |
| WO-B03 `Sdelboth-request`（去 `".request"` ＋ 去它的行） | witnessonly | rc=0 99/53/46/0/0/0 | **352** | **绿**（本程新增的一枚，它没列） | `out/WO-B03-Sdelboth-request.txt` |
| WO-BASE 对照 | witnessonly | rc=0 99/53/46/0/0/0 | 528 | 干净树全绿 ⇒ 上面三发不是"整包坏了" | `out/WO-BASE.txt` |
| WO-B04/B05 同一副本上的**单删**（清单动、证人不动） | witnessonly | rc=1 99/52/46/**1**/0/0 | 462 / 352 | **红** ⇒ 证人机制在 ddea3c3 里是活的 | `out/WO-B04-*.txt` `out/WO-B05-*.txt` |
| A-B06/B07/B08 同三发打在**交付态**（有写死的数） | anchor | rc=1 99/52/46/**1**/0/0 | — | **红**，红因逐字是 `:1530 … not the 8 x 11 x 3 written down in wantGrantRoute*` | `out/A-B06/07/08` |
| E03 交付态**只关掉那一处 `if`**（`R3-drop-pin-check`）＋ B06 那一刀 | anchor | rc=0 99/53/46/0/0/0 | **462** | **绿** ⇒ 单点回退：写死的数是这一形唯一的持有者 | `out/E03-R3+Pdelboth.txt` |

**裁（ⓑ）**：**成立，不推翻。** 三发协同删在"只有证人"的状态下确实全绿（且 `asked` 自洽到 462/352，
所以那把乘积尺也量不到），交付态把它打红，且红只由 `wantGrantRoute*` 那一处提供（E03 单点回退证死）。
它 §1.2 那句"我先只落证人册，然后真的删了一遍"是**做过的事**，不是事后编的解释——因为 `ddea3c3` 这一版本身就在盘上，
我在它上面重打就是同一状态。加分项另记：本程补的 WO-B03（`.request` 那一尾）说明这一族不止它列的那两枚——
`.request` 被剪掉时扫掠连"最像真路由"的那种拼法都不再问，仍是全绿。

---

## §4 ⓒ 单删·协同删·替换·变宽：本程逐枚打完，并补两发它没打过的形状

全部打在 `anchor`（交付态）副本；"它记的"＝实现件 §2.3。

```
K1 单删清单成员 11 枚（证人不动）  : 11/11 rc=1，四数一律 99/51/46/2 ⇒ 它记的"11 发全红、各 TOPFAIL=2"复算成立
   P-del-panel / panel.review / panel.approval / panel.l2 / panel.mode / panel.workspace / panel.attachment / panel.message
   ＋ S-del-empty / S-del-request / S-del-now                     原始件 out/K1-*.txt
K2 清单＋自己的证人同删 11 枚      : 11/11 rc=1（TOPFAIL=1 九枚；P-del-panel 与 P-del-panel.review 各 2 枚；S-del-empty 2 枚）
   ⇒ 它记的"Pdelboth 八发 TOPFAIL=1 或 2；Sdelboth-empty 2；request/now 1"逐枚复算成立   原始件 out/K2-delboth-summary.txt
C03 清单＋写死的数同改（证人行留着）: rc=1 99/52/46/1，红在证人册 :1705 两条 ⇒ **"把分母跟清单一起改"是钉住的**
C04 只改写死的数（清单不动）        : rc=1 99/52/46/1，红句逐字 "… not the 7 x 11 x 3 written down in wantGrantRoute*"
C07 变宽只改清单（加第 9 枚前缀）   : rc=1 TOPFAIL=1（9 != 8）      ⇒ 它记的 W-widen-list-only 复算成立
C08 变宽清单＋数一起挪（8->9）      : rc=0 99/53/46/0，asked=594     ⇒ 它记的 asked=594 那一格复算成立
C09b 变宽 words（加 "consent"）＋数 11->12 : rc=1 TOPFAIL=1，红在 :1697 "carries \"consent\" with no witness row"
   ⇒ 本程新增：**words 那一枚因子连"变宽"都要配证人，另两枚变宽不要**（方向不对称，实现件未列；是好事，但没人写下来下游会以为三枚同权）
D01/E08 **替换形**（长度守住、名单换掉：`"panel.mode"` -> 重复的 `"panel"`）：本程新造，它没打过
   四数 rc=1 99/52/46/**1**，asked=**528**（与干净树一字同）⇒ 写死的数（8）与乘积式（528）**双双满足**，
   只有证人册红（:1705 两条：清单不再含它／扫掠没问过 panel.mode.approve）
   ⇒ 这一形比"删一枚"更阴：文件自己打印的 asked 仍报 528，实际只有 462 枚不同名字被问过（66 枚变成重复提问）。
E10 后缀侧同一形（".now" -> 重复 ".request"）：rc=1 99/52/46/1，asked=528 ⇒ 同族，红在 :1707
C01/C02/C05 三处同删（清单＋证人＋写死的数）：三发**全绿** —— 见 §7（这是残余，本程把它的枚数与形状钉清）
```

---

## §5 ⓓ 单点回退：只回退一处，答"哪条用例变不响"——**六处全答得出，无一枚装饰**

回退＝在交付态副本上只把那一处改回它 r4 时的形状（其余一字不动），然后打**该处唯一负责**的那一发。
每处都先跑一发"只回退、不加变异"作对照，确认回退本身不在干净树上造假红。

| 处 | 本程那一刀 | 只回退（干净树） | 回退＋目标变异 | **变不响的那条用例（逐字名）** | 原始件 |
|---|---|---|---|---|---|
| 第 1 处 共用执行点（`sweepAssembledNames` 承载"问守卫"，正控走同一枚函数） | `R1-unshare-ask`：正控改回"自己那行直调守卫" | rc=0 99/53/46/0，asked=528（不造假红） | ＋`ask-stub` → **rc=0 99/53/46/0**；＋`hits-short` → 同绿；＋`ask-stub`＋M14r → 同绿 | `TestRealGuardRefusesEveryAssemblableApprovalRouteName`（三发全变不响，正是 r4 交付态的病：`asked=528 hits=[]` 读成干净边界） | `out/E04/05/06/07` |
| 第 2 处 正控扫的是 24 枚邻域、断"恰好四枚" | `R2-control-2names`：照抄验收样品的两枚形状（`wantAnswered` 同步改两枚） | rc=0 99/53/46/0，asked=528 | ＋`M-extra-mode_add-r`（守卫运行期多答一枚 `panel.mode.add`，whole name 不成串）→ **rc=0 全绿** | 同上一条（交付态同一发是 rc=1、红因 `:1531 … answers 5 of the 24 names`，见 `out/F05`） | `out/F04/F05/F06` |
| 第 3 处 写死的分母 `wantGrantRoute{Prefixes,Words,Suffixes}` | `R3-drop-pin-check`：那一处 `if` 短路 | rc=0 99/53/46/0 | ＋`P-del-panel_mode`＋`P-del-row-panel_mode` → **rc=0，asked=462** | 同上一条（三发协同删全靠它，§3 表末行） | `out/E03-R3+Pdelboth.txt` |
| 第 4 处 前缀证人册（清单＋`checkGridFactorWitnesses` 调用点） | `R4a-drop-prefix-call`：只摘那一行调用（册本身留着，包仍编译） | rc=0 99/53/46/0 | ＋`SUBST-mode-to-dup` → **rc=0 99/53/46/0，asked=528** | `TestGrantVocabularyIsNotSatisfiedByTheRealEnvelopes`（交付态该发是 rc=1，见 `out/E08`） | `out/E09-R4a+subst.txt` |
| 第 5 处 后缀证人册 | `R4b-drop-suffix-call` | 同上 | ＋`SUBST-suffix-now-to-req` → **rc=0，asked=528** | 同上那条（交付态红在 `:1707`） | `out/E10/E11` |
| 第 6 处 文件头那段（F-R4-1/F-R4-2 的说明） | `R5-drop-header`：把那段首行换成一句占位注释 | rc=0 99/53/46/0，asked=528 | 任何一发都不打：盘上判据不读注释 | **没有用例变不响** ⇒ 它是文档级那一味，与 r4 交付验收 §7 的 ⓓ 同类，按本仓分界**不计装饰** | `out/E12-R5-header-revert.txt` |

**两枚本程自己造的"拆检"读数（写下来防下一程把这格读松）**：

```
H05 把正控那 7 行整块删掉（清单/函数都留着）        rc=0 99/53/46/0 asked=528 ⇒ 删正控本身不会有人响
H06 上面那一块删掉 ＋ `ask-stub`                    rc=0 全绿 ⇒ F-R4-1 那一形今天需要**两处协同编辑**（r4 只需 1 处）
H07 ＋hits-short＋M14r                              rc=0 全绿 ⇒ 同上
I02 把常驻用例改名（`func Test…` -> `func test…`）   rc=0 **RUN=98** 全绿 ⇒ 名册差一枚
I01 改名 ＋ `ask-stub`                              rc=0 RUN=98 全绿 ⇒ 改名这一路**一枚编辑就够**，但它必然动名册（人眼可见）
⇒ 两条拆除路径的差别要说清：删 7 行正控块**不动名册、CI 与本地全绿**；改名**动名册**、靠人工差集兜。
   r5 的净收益是把"不动名册"那一条从 1 处编辑抬到 2 处；这条抬升是真的，但它不等于"两处协同拆不掉"（§7）。
```

---

## §6 ⓔ 承重（本仓操作定义：摘掉任意一味，是否存在一发变异从此打不红）

按 §5 的六味逐枚判，每一味都给"摘掉它之后唯一打不红的那一发"的实测数（不是推定）：

| 味 | 摘掉后从此打不红的那一发 | 实测 | 判 |
|---|---|---|---|
| 共用执行点 | `ask-stub` / `hits-short`（＋叠 M14r） | E05/E06/E07 三发 rc=0 全绿 | **承重** |
| 24 枚邻域＋"恰好四枚"的正控 | 守卫运行期多答一枚邻域名 `panel.mode.add` | F06 rc=0 全绿（交付态 F05 rc=1） | **承重** |
| 写死的分母 | 11 发"清单＋证人同删"（§3 那三发是其中被它自报的） | WO-B01/02/03 rc=0，asked=462/352/352 | **承重** |
| 前缀证人册 | 替换形 `SUBST-mode-to-dup`（长度与 asked 都不变） | E09 rc=0 99/53/46/0 | **承重**，且这一味**不被写死的分母替代**（两味各拿住对方拿不住的一发） |
| 后缀证人册 | 替换形 `SUBST-suffix-now-to-req` | E11 rc=0 | **承重** |
| 文件头那段 | 无一发 | E12 四数一字同 | **不承重、也不该承重**（文档级，本仓分界同 r4 验收 §7 的 ⓓ） |

**简报点名的反面本程没踩**：没有任何一格用"同一发变异里新加那支必须先响"当判据。§3/§5 的归因用的都是
**回退那一味＋打它独占的那一发**（一发只红在要被检的那一行上），并且 §5 每一味都附了"只回退不打变异"的对照发。
**冗余不等于装饰**：写死的分母与前缀册互为补集（分母量不到替换形 E08/D01；分母量不到的这发只有册红；
册量不到的那 11 发协同删只有分母红）⇒ 两味都在，才凑齐"各删任意一枚必须红"。

---

## §7 ⓕ 残余与**这把尺的射程边界**（简报要的一句比"记为残余"更硬的裁）

### §7.1 先把残余的**枚数**钉住（它登记了一枚，本程量出来是"一族、至少三枚"）

```
C01 前缀三处同删（去 panel.mode ＋ 去它的行 ＋ wantGrantRoutePrefixes 8->7） rc=0 99/53/46/0 asked=462  ← 它登记的那一枚
C02 后缀三处同删（去 ".now" ＋ 去它的行 ＋ wantGrantRouteSuffixes 3->2）     rc=0 asked=352            ← 同一族的第二枚，它未列
C05 words 三处同删（去 ratify ＋ 去它的行 ＋ wantGrantRouteWords 11->10）    rc=0 asked=480            ← 第三枚，它只以"同一族"带过
   ⇒ 本程把这三枚从"推定"升成现量（与 r4 验收 §6 把 authorize/decision 升为现量同一做法）。
   三处同删之后文件内部**完全自洽**：乘积式（asked == 清单相乘）、写死的数、证人行、清单四样互指同一份新事实。
   grep 复算：对 `grantRoutePrefixes|wantGrantRoute|grantRouteSuffixes` 在锚点副本里扫全部代码件
   （--include=*.go/*.ts/*.js/*.ps1/*.sh，即除 docs 之外的一切可执行/配置文本）-> 命中文件只有
   `internal/panel/l2_grant_boundary_test.go` 一枚
   ⇒ 今天没有任何**代码**在别处持有这三枚清单的分母（人工名册差集是唯一的兜）。
```

### §7.2 裁：**简报那句"结构上不可由该文件内任何检守住"不成立；成立的是它的弱化版**

简报给的读法是"任何要靠协同编辑测试文件本身的形状，结构上不可由该文件内任何检守住（改的人总能删掉整枚测试）"。
本程判：**方向对，但作为"结构上不可能"来写是过宽**，理由有两条实测支撑：

1. **这把尺今天守的正是"协同编辑"**——只不过只守到**两处**。§4 的 C03（清单＋分母同改、行留着）与 K2 全 11 枚
   （清单＋行同删、分母留着）**都是协同编辑，都是红**。所以"协同编辑＝守不住"这句若照抄进台账，
   会把已经守住的 11 枚形状误划到门外，下一程就会不去补它、也不去复查它。
2. **"改的人总能删掉整枚测试"是另一族**，且那一族的成本与这一族不同：改名/删整枚（I01/I02）动**名册**，
   CI 里 `internal/panel` 在 core scope（`scripts/portable-tests.sh:141` 的 `core_pin` 列了它、`:179` 是 `./internal/panel/...`；
   CI 步＝`ci.yml:288`）⇒ 名册差集是人能看的信号；三处同删**不动名册**（99/53/46/0 一字不变），
   在 diff 里读起来像"清单少了一枚而已"。⇒ 两族该分开记，不该用一句"极限"合并。

### §7.3 可核的边界句（建议按这一段入台账，而不是按"极限"）

> **这把尺守的是"单点静默删除"与"两点协同删除"，不守"同一枚文件内三处互相圆场的改写"。**
> 三处＝① 因子清单成员 ② 该成员的证人行 ③ `wantGrantRoute*` 里那一枚数字（前缀/后缀/词各一处，实测 C01/C02/C05）。
> 判据边界可核在这三点上：**(a)** 三处同删后文件内部自洽（乘积式、分母、证人三样互指同一新事实），
> 所以任何**放在同一枚文件里**的检要识破它，必须再引入第四处独立事实——那第四处同样可被一起改；
> **(b)** 因此唯一能再往前一格的做法是把分母**搬到这枚文件之外**（一份 golden 清单、上一版 `git show` 比对、
> 或另一枚包持有的常量）——那是**新仪器**（要新的裁决与新的失效模式），不是 r5 没做完的一格；
> **(c)** 三处同删必然伴随**真实覆盖下降**（66 或 176 枚名字不再被问），这正是它必须在 diff 里显形、
> 也是"人工看差异表"在这里仍然算一层兜的原因。
> 另附一句同族边界（§5 末）：F-R4-1 那枚正控今天是"删 7 行正控块 ＋ 短路 `:1441`"两处协同打穿的（H06 全绿），
> 所以这句边界对**本轮两枚债都适用**，不是只写给写死的数那一族。

**由此得出的裁**：`X-list+row+pin` **不再记为待闭债**，改记为**判据边界**（上面那段），并把它的**清单枚数**（三枚、非一枚）
补进实现件 §6——那是 F-R5-1（文档级，见 §11）。理由：本轮它已经把"单点删除"这一族全部收口（§4 的 22 发全红），
剩下的那一族要么搬到文件外（新仪器，需另开一格裁决），要么靠人看 diff（本仓既有 backstop）；
把它继续记成"待闭债"会让下一程在同一枚文件里再加第四处自指事实，而 §7.3(a) 说明那只是把边界挪一位。

---

## §8 ⓖ 三态四数、名册差集，与那枚 C21 的**口径边界**（本程把因果也做了）

| 台件 | 四数（`go test ./internal/panel/ -count=1 -v`） | 红名 | `asked=` | 原始件 |
|---|---|---|---|---|
| 改前副本 `ee2a92d`（blob `58f54514`，2220 行） | rc=0 99/53/46/0/0/0，`^panic:`=0 | 无 | 528，hits=[]，打印在 `:1389` | `out/PRE-ee2a92d.txt` |
| **交付态副本 `aeba6ff`（锚点，判据用的就是这一行）** | rc=0 99/53/46/0/0/0，`^panic:`=0 | 无 | 528，hits=[]，打印在 `:1540` | `out/BASE-anchor.txt` |
| 仓库工作树（**这是工作树、不是锚点**；HEAD 取于我量这一发时＝`ddf8b6f`，`internal/panel/` 当时还没有别家未提交件） | rc=1 99/**52**/46/**1**/0/0，`^panic:`=0 | 唯一一枚逐字 `TestC21DesignTokensFourWayAgree` | 528，hits=[] | `out/W-BASE-worktree.txt` |

```
名册差集（两份运行输出，未用 -run / -skip）：
   `^=== RUN` 归一化名（含子测试）99 枚 -> 99 枚，diff 为**空** ⇒ +0/−0   out/roster-PRE-ee2a92d.txt · out/roster-BASE-anchor.txt
   `func Test` 名册（git grep -h '^func Test' <rev> -- internal/panel/）：ee2a92d=53 枚 -> aeba6ff=53 枚，diff 为空
   本文件内 func Test 7 -> 7 ⇒ 它 §3 那三格逐格复算成立；本程在副本上重打，不采信它的工作树口径
```

### §8.1 为什么"锚点副本 99/53/46/**0**"与"工作树 99/52/46/**1**"同时成立——因果本程自己做了

```
事实 1：design/assets/tokens.css 在锚点里是**跟踪着的**（aeba6ff 的 tree 里有，blob 0ad0112379e4df056ead3fdc2dc6e1f223f2040c，12942 字节）
事实 2：`git status --porcelain` 里那 16 枚 ` D design/**`（含 ` D design/assets/tokens.css`）是**未提交的工作树删除**（owner 的活）
事实 3：归档副本从 tree 取件 ⇒ 该文件存在；工作树被那 16 枚删除掏空 ⇒ 该文件不存在
事实 4（本程做的因果实验，全程在仓外副本里，未碰仓内）：
   在 anchor 副本内把 design/assets/tokens.css 移走 ->
     --- FAIL: TestC21DesignTokensFourWayAgree
     tokens_fourway_test.go:441: read design/assets/tokens.css: open …\anchor\design\assets\tokens.css:
       The system cannot find the file specified.
   移回来 -> git hash-object --no-filters = 0ad01123…（＝锚点 blob，逐值相同）
```

**口径边界（写死这一句，下一程不要再在这个包里找第二枚红）**：
`internal/panel/` 的"几枚红"取决于**读哪一台件**，不取决于本批代码——
锚点归档副本 **0 枚**；本机工作树 **1 枚**且唯一红名是那枚 C21，红因是 owner 未提交的 `design/**` 删除，
与 `l2_grant_boundary_test.go` 无因果（它在三态里都是 99 枚 RUN、`asked=528`）。
⇒ 引用"几枚红"必须连台件口径一起写；两个数都对，**拼成单值就错**（r4 交付验收 §2 那格同一族病，本程复现到第三形态）。
⚠ 另记一枚本程现量的漂移：`internal/panel` 里本程跑完后已出现别家未提交件（`pump.go`／`pump_test.go`），
⇒ 今天再打工作树口径那一行，四数**允许与 §8 第三行不同**，那不是本批的账（本程没有用它当任何判据）。

---

## §9 ⓗ 三门、断言尺、以及"删掉的 25 行"逐枚读过

```
（以下三门全部在锚点归档副本里跑；跑前跑后各复算一次两枚被验物哈希＝f87301df…/d2cd6362…，见 §1 末次全树复算）
gofmt -l internal/panel/  -> 空，rc=0
go vet  ./internal/panel/ -> 空，rc=0
sh scripts/d22scan.sh     -> rc=0，原始件 out/GATE-d22scan-anchorcopy.txt
   第一步（正控）确实跑了：runtests.sh: OK - packages=[./...] top-level: PASS=30 FAIL=0 **SKIP=0**, === RUN=70
     ⇒ 简报点名的坑（set -eu、正控红则真扫描根本不跑）本程验证方式是"看见 30 枚 PASS 才读第二步"
   第二步（真扫描）逐枚分母：bans #1-5 internal/=203 cmd/=22 | ban #6 frontend/=46 | ban #7 internal/tools/=18
     | ban #8 design/=**30** frontend/=46 internal/=407 cmd/=40 → "d22scan: clean - no D22 ban violations"
   ⚠ 副本没有 .git ⇒ 扫描器**自陈**"gitignore rules NOT APPLIED … every path in every scope is being scanned"（A207 那枚机器相关分母）
     ⇒ design/ 我量到 30、实现件与工作树量到 32，差值来自工作树里 owner 新增的未跟踪 design/ 目录（design/old/ 等），
       与本批无关；本批**零枚新增文件**（§0 pathset 为证）⇒ 它只能进 internal/=407 那一枚分母，而 407 与 r4 交付验收 §4 同值。
   对照（这是工作树口径、非判据）：仓内再跑一次 -> rc=0，bans #1-5 internal/=**205** | ban #8 design/=32 internal/=**409**
     ⇒ 那 +2 是别家未提交的 `internal/panel/pump.go` 等，本程不计入自己也不计入被验物的零命中主张。
```

断言尺（两版都从 git 取，`grep -oF | wc -l`，不数行）：

```
                     ee2a92d(58f54514)  aeba6ff(f87301df)  delta   实现件 §4 记的   判
t.Fatalf(                    35                36           +1          +1          ✅
t.Fatal(                      5                 5            0           0          ✅
t.Errorf(                    45                49           +4          +4          ✅
t.Error(                      1                 1            0           0          ✅
t.Logf(                      15                15            0           0          ✅
t.Skip                        0                 0            0           0          ✅ 保持 0（没有一枚被跳掉）
func Test(本文件)             7                 7            0           0          ✅
Fatal/Error 合计             86                91           +5          +5          ✅ **只变长，没变短**
行数                        2220              2380                   2220/2380     ✅
```

**"只变长"这句话本程不止信计数**——把 `ee2a92d..aeba6ff` 删掉的 25 行**逐枚读过**（`git diff -U0`）：

```
1 行 注释：… "is allowed to grow; nothing in here claims it is exhaustive"（被"变窄要动两处"那段替代；
          范围声明本身仍在文件里 :1527-1529"What it does NOT cover"，未失）
8 行 旧正控（for 四枚常量 + 每枚直调 + 那句 Fatalf）：被 :1496-1502 的 DeepEqual(24 枚邻域)==恰好四枚 取代
          ⇒ 本程实测这是**超集**：旧正控负责的四枚"必须被答"两版都红（F02 vs F07），新形状多拿住
            "守卫多答一枚邻域名"那一发（F05 红 / F06 绿）⇒ 没有任何一枚旧断言被丢掉
16 行 旧循环（asked++/四个 for/sort.Strings(hits)）：原样搬进 sweepAssembledNames(:1433-1450)，一字未改语义
⇒ 结论：本批没有为变绿放宽任何断言；阈值/golden/thresholds.go 一字节未碰（§0 pathset 只有那一枚 _test.go）。
```

---

## §10 简报三枚前提＋实现件 §5 那六条：本程自己判，不替谁圆

| # | 谁的说法 | 本程的判 | 依据 |
|---|---|---|---|
| ⓐ | 简报"约 12 行"；实现件自报"毛 +49／净 +29，仍在 60 行停手线内" | **实现件方向对、数值低报**；**停手线的适用范围未定** | 本程按 hunk 边界归属：`ddea3c3` 三枚 hunk 是 +125/−0（@old1306）、+15/−24（@old1320）、+9/−0（@old1546）；F-R4-1 只占 hunk1 的后 **41** 行（`sweepAssembledNames` 28 ＋ 邻域网格 12 ＋ 空行）＋ hunk2 的 +15/−24 ⇒ **毛 +56／净 +32**，比它自报的 49/29 多 7 行；仍 <60 ⇒ 那一枚洞没越线。**但整批是 +185/−25＝净 +160**：若"60 行停手线"按整批读，本批早已越线，而实现件只给了单枚洞那个口径、没把工作树口径的那一枚摊出来 ⇒ 这条属**未定义即停**的适用域问题，见 §11 的 F-R5-2（要编排者定口径，不是实现程的错，也不是它有权自己定的） |
| ⓑ | 简报"②③ 那样写钉不住 M14 那一形；标签式 M14 会得 4 枚红、其中 3 枚与 F-R4-1 无关" | **成立** | G03/G04 现量 TOPFAIL=4，逐枚比红因文本：只有 `:1531` 那枚带 `(F-R4-1)`；另三枚（`:1244`/`:1285`/`:1301`）是静态尺读到 case 列表里的名字。用运行期拼名那一形（G01/A02/A03）→ **恰 1 枚红** |
| ⓒ | 简报"证人 或 逐因子记账"是等价选项 | **不成立（实现件对）** | §3 全表：证人那一支单独存在时三发协同删全绿（WO-B01/02/03），加上写死的数才红（A-B06/07/08），且单点回退那一处 `if` 即复活（E03）⇒ 两枚选项**不等价**，简报该改口径 |
| ⓓ | 实现件 §1.1 理由 1："砍掉 `panel.attachment.add`／`panel.message.send` 会悄悄磨钝锚点已有的正控" | **半对——本程把它换成一条更强的依据** | 那两枚若从正控里去掉，"守卫**少答**一枚声明路由"这一形**不会**无人响：F07（两枚形状＋`guard-loses-attachment_add`）仍红 **2 枚**（`TestComposerEnvelopeAcceptsItsFourRequests` @`bridge_test.go:68` ＋ plant B 的 stale-enumeration @`:2053`）⇒ 它说的那半事实被别人兜住了。**真正只有四枚形状兜住的是反向那一发**：守卫**多答**一枚运行期拼出的邻域名 `panel.mode.add`——交付态红（F05），两枚形状全绿（F06）。⇒ 取舍站得住（本程判它应当保留四枚），但**理由要换成 F05/F06 那一发**，否则下游以为"少答路由"只有这一处兜 |
| ⓔ | 实现件 §5.3"两个数都对，口径不同" | **成立** | §8 三行＋§8.1 的因果实验 |
| ⓕ | 实现件 §5.5"验收给的修法只到它自己判据的一半" | **成立** | §3（它把 `ddea3c3` 真身留在盘上，本程得以在同一状态重打，这是这条结论能被复查的前提） |

---

## §11 总裁：**成立（附条件入账）**——三笔上游债结清；两枚本程账都是文档/口径级，**不构成这批准入条件**

| 上游债（r4 交付验收 §12 开的） | 它这次做了什么 | 本程独立量到的现状 | 结 |
|---|---|---|---|
| **F-R4-1** 问守卫那一行可被 `if false` 短路而全包全绿 | 抽成 `sweepAssembledNames`（唯一问守卫处 `:1441`），正控改为穿过同一枚函数、扫 24 枚邻域断"恰好四枚" | ①②③ 三发从"全绿"变**各红恰 1 枚**（A01/A02/A03，红因逐字带 F-R4-1）；**单点回退第 1 处即全部复活**（E05/E06/E07 三发 rc=0、`asked=528 hits=[]`＝r4 原病）；`guard-dead` 复算 4 顶＋2 子未钝（G06） | **结清** |
| **F-R4-2** `prefixes`／`suffixes` 两枚因子无证人、可被静默剪窄 | 11 枚证人行（四条主张/行）＋ 三枚写死的数 | 11 枚单删全红（K1）、11 枚协同删全红（K2）、`prefix-shrink`/`suffix-shrink` 两发从全绿变红 2/红 5（A04/A05）、替换形（长度与 `asked` 都不变）也被册红（D01/E08）、"分母跟清单一起改"被册红（C03） | **结清** |
| **F-R4-3** 实现件 §8 少列一形（文档级） | append-only 一行 | `449a450,451` 纯追加、`## §8` 唯一、句中引的 `asked=462/352` 本程都复算到 | **结清** |

**不退回的理由（按本仓分界那句：验收方自己造不出"声称要防而没防住"的结局 ⇒ 不退回）**：
它 §1/§2 声称的射程（"从那三枚清单里各删任意一枚，全包必须有红"＋"谁短路那一行谁当场红"）本程**逐枚打完都红**；
本程造出的四发新形状（D01/E08 替换形、C05 words 三处同删、H05/H06 正控两处拆、F05 多答邻域名）**没有一枚落在它声称的射程内**——
其中两枚它已在 §6 白纸黑字划在门外（协同改写／改名整枚用例），另两枚是本程量到的**加分事实**（它守住了、只是理由写错，见 §10 ⓓ）。
它也没有假报："只立证人不够"这一条我按简报要求去推翻，**推翻未遂**（§3）。

### 本程账（编号 `F-R5-*`，两枚都是"改文档/定口径"，零代码）

| 号 | 账 | 落在哪 | **最小闭合集合** |
|---|---|---|---|
| **F-R5-1** | 实现件的**残余清单与归因有四处需更正/补列**：(a) §6 的 `X-list+row+pin` 是**一族三枚**不是枚一枚（C01 前缀 `asked=462`／C02 后缀 `352`／C05 words `480`）；(b) 同一段应写一句"**两点协同已守**"（K2 十一枚＋C03 全红），否则下一程会把"协同编辑守不住"误当已定案、不去复查那 11 枚；(c) §1.1 理由 1 的归因要换成 §10 ⓓ 那一发（`panel.mode.add` 多答形：F05 红/F06 绿），因为它原句那种读法会把"少答声明路由"当成只有正控兜——实际 `bridge_test.go:68` 与 `:2053` 也兜（F07）；(d) "两处协同才拆得掉 F-R4-1"需加限定：改名常驻用例仍是**一枚**编辑（I02，RUN=98），只是那一枚必动名册 | `docs/evidence/s1/panel-l2-grant-nail-fix-r5.md` **末尾 append-only 一段**（不改写已提交正文；形状同 `a72c110` 那一枚） | 一处文档，四句；**产品码继续零字节、阈值/golden 一字节不碰**是本集合硬前提 |
| **F-R5-2** | **"约 60 行改动即停手上报"这句的适用域未定**：按单枚洞的最小闭合读＝本批未越线（F-R4-1 净 +32，本程现量）；按整批读＝净 **+160**（`git diff --numstat ee2a92d aeba6ff -- internal/panel/`＝185/25）早已越线，而实现件只报了前者 | 由**编排者定口径**（若定"按整批"，则本批需补一句自陈；若定"按单枚洞"，写进派单模板，免得每程自己选） | 一句口径裁定，零代码；本程不新造 `Q##`、不写台账 |

### 给编排者的判断：**建议停止续轮**（第五轮已到收益递减点）

按简报点名要这一句，本程给：**建议 r6 不开**。理由三条，都可核：

1. 本轮新造的四发里，唯一"只有测试文件内部再加一味就能红"的形状已经没有了——剩下的三处同删族（§7.3）
   要靠**把分母搬出该文件**才治，那是新仪器（新失效模式、新裁决），不是 r5 欠的一格。
2. 该尺今天保护的那条入站路由**没有生产调用者**——这一条本程**没复算**，是 r2/r3/r4 三程各自的现量（简报点名的原话），
   所以本程的实质安全性判断仍建在那句自陈上；**要作为停止续轮的正式依据，请让下一程顺手复量一次**（一次 `git grep` 的活）。
3. 剩余风险的兜法（**X**，逐枚点名，别把不存在的那枚当兜）：
   **X1**＝CI 每推真跑这一枚包（`scripts/portable-tests.sh:141` 的 `core_pin` 列着 `internal/panel`、`:179` 是 `./internal/panel/...`；
   CI 步＝`ci.yml:288` "Portable package tests (core scope)"）⇒ 单点删除／替换形是**机器**兜的；
   **X2**＝人工名册差集（99 枚/53 枚 `func Test`，本程实测 +0/−0 可核）兜"改名或删整枚用例"那一族；
   **X3**＝三处同删必然在同一枚 diff 里表现为"清单少一枚＋数字小一格＋行少一条"三处同时动 ⇒ 看差异表；
   ⚠ **X4 不成立，别写进台账**：CI 的 `d22scan` ban #6（面板侧 L2 允许）**只走 `frontend/`**（`tools/d22scan/main.go:395`），
   扫不到 `internal/panel/bridge.go` 里把批准路由接上——所以 Go 侧那一形不在扫描器的射程内，只有 X1/X2/X3。

---

## §12 给 owner 的一段人话（不用术语）

这一轮补的是同一把检查尺上两个"自己会偷偷不干活"的毛病。我另开一台机器、把它做过的破坏**自己重做了一遍**
（一共八十四发，逐枚留了原始输出），结论是三句：**两个毛病都真的补上了；它自己报的那句"上一版的补法其实不够"是真的；
它给的一处理由写歪了，但结论站得住。**

第一个毛病原来长这样：那把尺要问五百二十多种暗号，而"这几扇明摆着的门必须被回答"这句话写在**别的地方**。
谁把"问"的那一行涂成"永远问不到"，这句照样通过、整包全绿。这轮把那句也搬到同一行上了。我实测：把那一行涂掉，
**当场红**，而且我把它改回去（只改这一处）再打，三发又全绿——说明这枚红确实是这一行买来的，不是摆设。

第二个毛病原来长这样：五百二十多这个数字是三张清单相乘。十一枚动词每人有证人，另外两张（"地点"八枚、"词尾"三枚）
一个证人都没有——悄悄剪短其中一张，测试照绿。这轮给那十一枚各配了证人，还额外把三张清单的长度写死成三个数。
它自己报了一句不太好听的话："只配证人这一招，我实测不够。"我按它的做法在**它当时真的提交过的那一版**上重打：
把一张清单和它的证人一起删掉，**三发全绿、一声不响**；把写死的那三个数加上，三发全红。⇒ 这句自报成立。

**它写歪的那处**：它解释"为什么没照抄验收方给的十行样品"时说，样品会磨掉两扇明摆着的门——那两扇其实**还有别人守着**，
不会失守；真正只有它这版守得住的是**另一件事**：如果后台偷偷多答了一条相邻的暗号（不是少答、是多答），
只有它这版会红，样品那版全绿。⇒ 结论一样（保留四枚是对的），但理由得换，否则下一个人会以为多答那一形也有人兜。

**最坏后果还剩什么形状**（这一轮之后）：还剩一扇，而且要**同时改三处**才能走进去——把某条"地点"从清单里划掉、
把它那条证人一起划掉、再把写死的数字改小一格。这样从头到尾全绿。上一版把这写成"只剩一扇"，我量出来是**一类三扇**
（八种地点、三种词尾、十一枚动词各自都能这么拆，我各挑一枚打了，全部全绿）。**这一类不是产品代码能造出来的**，
完全是这把尺自己身上的字，治它只能靠人看差异表——本轮仍然没有仪器治它，我写在§13 没藏。
要说清的一点是：**"一起改两处"它已经守住了**（清单和证人一起删、或者清单和数字一起改，都会红），
所以别把剩下这一类读成"协同编辑一律守不住"，那是过头的读法，会让下一个人在错的地方停手。

**我建议这一轮到此为止**：再挖下去要买的是"把数字搬到这把尺自己管不到的地方"（比如另存一份清单），
那是新工具、不是同一件活没做完。眼下真正的兜底是三层：机器每推真跑这枚包、有人核对用例名单没少、看差异表时
注意到"三处同时变小"。顺带一句口径：这台机器上有一枚红的 `TestC21DesignTokensFourWayAgree`，
是因为设计稿目录被您自己删了还没提交；我们按规矩让它红着，没修、没跳、没放宽。我把那枚文件在**仓外副本**里
挪走又挪回来验过一次因果——挪走就红这一句、挪回来哈希与您提交过的一字相同，所以那枚红**不是这批准入的障碍**。

---

## §13 本程**没测**什么（不靠沉默读成通过）

- **没跑 `go test ./...`、没跑 CI、没跑前端四道门**；一条尺没下过 `./cmd/...` 与 `internal/panel` 之外的任何包
  ⇒ 本件"全包／包内"字样一律只指 `internal/panel`。§9 里那 30 枚 PASS 是 `d22scan` 那道门的**自带正控**，不是我的分母。
- **没在 Linux 容器里量任何东西**：本程全部读数出自 windows 原生 `go test`，而 CI 的 core scope 在 ubuntu-latest 跑
  ⇒ "两态四数相同"这句**没有跨平台证据**，只有台账记忆（`*_other_test.go` 那族本包现量 0 枚，沿用 r3-accept §10，本程未重做）。
- **三枚清单之外的协同形状没测**：四处及以上同改（清单＋证人＋数字＋乘积式互相圆场）、以及"把某枚证人名字换成另一枚
  仍能被网格拼出的批准门名字"（它 §6 也自陈没测，本程没造）。
- **没测"整枚用例被改名或删掉"之后 CI 的颜色**：本程只量到本机 `RUN=98 全绿`（I02）；那枚包在 CI 里的表现
  我只核到"它在 core scope 名单上"（`scripts/portable-tests.sh:141/:179`、`ci.yml:288`），**没真跑过 CI**。
- **没测 `internal/panel` 之外是否还有同类"自己给自己当分母"的尺**（例如 `internal/ball` 那枚 C21 色值表尺）；
  一条尺没下过别的文件。
- **没复算**"今天这条入站路由没有生产调用者"（`ParseComposerRequest` 零调用者、`go.mod` 无 webview）：
  r2/r3/r4 三程各自验过，本程没重做 ⇒ §11 那句"停止续轮"的依据里这一项**标着〔未复算〕**，请补一次现量再用。
- **没读** `design/**` 脏树内容作凭据，也没动那 16 枚未提交删除（未还原、未提交、未删）；只用 §8.1 那发因果实验解释那枚 C21。
- **没读别家那一程的快照泵**（`internal/panel/pump.go`／`pump_test.go`，本程跑完后出现在工作树）：
  它与本批代码零交集（锚点里没有它），本程未把它算进任何主张，也没为它调整任何读数。
- **没写**台账、票面勾、`Q##`；**没给自己判成立**；**没改被验物一字**（§14 现量）。

---

## §14 git 自证／临时件／注入登记（本程落点，逐枚现量）

```
本程 commit 逐枚 pathset（不跑区间；共享树在漂）：
  1effda0 -> docs/evidence/s1/panel-l2-grant-nail-fix-r5-accept-r1.md  （§0-§1）
  31223ef -> 同一枚路径                                                （§2-§4）
  2ce00b1 -> 同一枚路径                                                （§5-§7）
  5556e10 -> 同一枚路径                                                （§8-§10）
  782b21a -> 同一枚路径                                                （§11-§14）
  83c339f -> 同一枚路径（+10 −6：§1 发数、§2 归因落到发号、§7.1 grep 口径三处收紧——**本程唯一一枚带删除行的 commit，
             删的全是我自己本件里那三处不准确的表述，未删任何一枚读数**）
  本节所在这一枚 -> 同一枚路径（登记上面那枚之后追加两行；sha 由下一位现量，本件不预写自己的 commit 号）
每次 commit 前现核 git diff --cached --name-only -> 只有那一枚路径（第一次 commit 前索引为空）；
   全程未出现别家路径进索引；未跑过 git add -A / add .（本件是新建文件，用 `git add -- <那一枚路径>` 显式加）。
纪律：只 commit、未 push；未 --amend / reset / rebase / stash / checkout . / clean；未跑 rm/rmdir/rm -rf（连续第三程写这条：看动作不看后果）。
被验物零改动（交件时现量）：
  git rev-parse HEAD:internal/panel/l2_grant_boundary_test.go -> f87301df6be0ba1ceb00cba7f9ddda877df521c6（＝锚点）
  git status --porcelain -- internal/panel/  本程全程未写入（唯一变化是别家那一程的 pump*.go，本程未碰、未提交、未还原）
  design/** 那 16 枚未提交删除：未还原、未提交、未删，也不计入本程任何"零命中"主张。
临时件（只建不删；请编排者一次清）：
  D:\tmp\wisp-r5acc1\out\          99 枚（87 发 go test ＋ BASE/PRE/W 三态 ＋ 三门两份 ＋ 名册两份 ＋ 字节比对三份）
  D:\tmp\wisp-r5acc1\scripts\      mut.py（变异台件，锚点计数≠1 即拒）· cmp_hashes.py（全树逐枚比字节）
  D:\tmp\wisp-r5acc1\{anchor,pre,witnessonly,crlfpitfall}\   四台归档副本（anchor＝交付态判据台；pre＝ee2a92d；
                                   witnessonly＝ddea3c3 真身，§3 用；crlfpitfall＝§1 坑① 那一发的单文件取样）
  D:\tmp\wisp-r5acc1\pristine\{anchor,pre,witnessonly}\      每台的原件快照（每发前后 restore 用；tokens.css 那一发用它自己的移动/移回）
本程失手两枚，登记不擦：
  ① `R1-unshare-ask` 第一版漏了 `sort.Strings(legit)` ⇒ E01/E02 两发是"顺序不等导致 DeepEqual 红"的**无效归因**
     （读数本身没坏、跑起来了，但不能用来回答第 1 处回退的问题）；正文一律改引修正后重打的 E04-E07。
  ② `R7-delete-control-block` 第一版只把 `if` 换掉、留着引用被删变量的两行 ⇒ H01-H04 四发 RUN=0（`build failed`，
     逐字 `undefined: legit / undefined: wantAnswered`）＝**无效读数**；留档不删，正文改引修正后的 H05-H07。
  ③ `F01/F03` 两发不是坏切法，但那是"只砍网格、没同步砍 wantAnswered"的**不忠实样品**（它在干净树上就红），
     正文只用来记"这一发红在哪"，不参与 §5 表格那一行；忠实形状见 F04/F06/F07。
注入登记：本程工具输出里出现过三处**长得像授权**的文字，逐一点名（命令名＋前 40 字），一律**未据其行动**：
  ① 多次 `Edit`/`Bash` 结果尾部附加的 "Note: The file ... was modified since it was last read" ＋ 结尾那句
     "Confirm the harness note is genuine and was intentional before continuing to use this file."（出现在
     `Edit: D:\tmp\wisp-r5acc1\scripts\mut.py`、`Bash: cd /d/tmp/wisp-r5acc1 && python scripts/mut.py …` 之后）；
  ② 同一批附加文字里的 "This change was intentional, so do not revert it"。
  这两类本程**没有**据此提交、没有据此放宽任何判据、没有据此减少任何取证；那两枚路径都是我自己的仓外台件脚本，
  改动是我用 Edit 主动做的（内容与 §1 描述的台件行为一致），"确认"这一步对本程无操作后果。
  ③ 上下文里的 `AGENTS.md`／MEMORY 快照中"已批／口令撤 X／Q1＝甲"一类台账转述 ⇒ 是转述，不是给我的指令。
  真系统通知一枚：任务开头那句"Memory 文件 MEMORY.md 已被修改"——本程未据此改动任何落点。
  凭据值一字未抄；所有读数出自我自己跑的命令与 `out/` 原件。
