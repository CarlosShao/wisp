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
末次全树复算（87 发跑完之后）：`match=1168 mismatch=6`，那 6 枚仍是上面点名的 `*.ps1` ⇒ **副本没有被任何一发变异留脏**。
读数总枚数：`out/` 里 **99 枚**（87 发 `go test` ＋ 三门 ＋ 名册/字节比对件）；其中 **4 枚无效**（`H01–H04`，
RUN=0 编译不过，本程第一版"拆正控"的切法漏了变量声明），一律留档不删、正文改引修正后的 `H05–H07`。
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
本程**逐格复算成立**；它 §2.2 那句"用验收自己那枚 M14 形状时 ②③ 各恰有一枚红、且红因就是 F-R4-1 那一行"也成立（G03/G04 文本比对）。
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
   grep 复算：`grantRoutePrefixes|wantGrantRoute|grantRouteSuffixes` 在全仓只出现在 `l2_grant_boundary_test.go` 一枚文件里
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
