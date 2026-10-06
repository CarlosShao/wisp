# 凭据索引 — 已结案票（`-done`）全池 × `docs/evidence/s1/` 凭据对账

本件不是对抗验收表，也不是缺口审计；它只是凭据索引。⛔ 它不判任何一票"该不该结"、⛔ 不裁任何一格的通过／退回、⛔ 不替任何一票翻勾。**任何一格都不因本件而视为已验收。**

> 为什么有这一件：`.scratch/wisp/issues/` 里挂着 `-done` 后缀的票已有 **93 枚**（本腿现量，见 §2），而 `docs/evidence/s1/` 里以票号命名的件只覆盖其中一部分。编排者需要一张"哪枚结案票在 `docs/evidence/s1/` 有终裁表、哪枚只有旁证、哪枚完全无凭据"的账，用来派后续的非实现者腿。本件只做这件事：**把每枚结案票名下真实存在的凭据文件名指出来**（不存在就写 `NOFILE`），⛔ 不补、不猜、不美化。
>
> ⚠ 本腿是**收口腿**：前三枚同职腿（`evidence-close-1` 出成品 1 枚、`evidence-close-2` 与 `evidence-close-3` 均死在中途）把扫描数据全落在了盘上（`6050e171`／`b3d29b0d` 两笔入库），但没有写成正式交付件。本件把这些数据**选定一把尺**后写成账，⛔ 不重做扫描、⛔ 不把三版尺平均或混引（三版差异的具名裁决见 §1）。

---

## §0 起手锚与本腿读数名册

### 0.1 起手锚（本腿现跑）

| 项 | 读数（逐字，本腿现量） |
|---|---|
| 本腿代号 | `evidence-close-4`（收口文书腿，零 `go` 命令） |
| 进场时刻 | `2026-10-06 09:12:25 +0800`（`date` 自取） |
| 进场 HEAD | `290dca87`（`票 231 AC#2 · 给"版本更高"一条自己的出口（cause=newer-build）`，2026-10-06 09:12） |
| 分支 | `dev`（`git rev-parse --abbrev-ref HEAD` 现量；`ahead 147` 未推） |
| 上一职腿遗产锚 | `b3d29b0d`（`evidence-close-3` 36 枚扫描台件入库，2026-10-06 09:06）／`6050e171`（`evidence-close-2` 29 枚） |
| 本腿写面（授权） | 仅两处：`docs/evidence/s1/closed-tickets-evidence-index.md` ＋ `.scratch/wisp/probes/evidence-close/4/**`；⛔ 不改 `evidence-close/2`、`/3` 任何一枚原件、不改票面一字、不碰台账／产码／`docs/reports/**` |
| 禁跑项 | 全程⛔ 无 `go test/build/vet/env`；⛔ 不读不引 `frontend/**`、`design/**`；⛔ 不引 `.scratch/wisp/probes/231`、`111/r4`、`pool-validity`（并行腿半成品，`231-r1` 是活腿） |
| 在飞面（非本腿，⛔ 不碰） | `cmd/wisp/**`＝` M`（活腿 `231-r1` 正在写）；工作树另有 `design/**` 那批 `D`/`M` 老脏面 |

### 0.2 本腿读数名册（逐条：命令 → 读数 → 时刻）

| # | 命令（形状） | 读数 | 台件 |
|---|---|---|---|
| 1 | `ls .scratch/wisp/issues/*-done.md \| wc -l` | **93** | `denominator.txt` |
| 2 | 同上抽票号 → 排序 | **93 枚、与死腿 `3/nums93.txt` 逐字节 `diff`＝零差**；⚠ `3/live-done-nums.txt`＝**92 枚**（少 `268`，腿在 22:46 前取数、票 268 于 22:56 才改名 `-done`）；`2/done-nums.txt`＝92 枚（同样少 `268`） | `nums4.txt`／`done-files-4.txt`／`mine-nums.txt` |
| 3 | 同名件名册：逐票 `ls docs/evidence/s1/ \| grep -E "^0*<票号>([-/]\|$)"`（本腿自跑，⛔ 不引用死腿 `named-own-pairs.txt`） | **146 对**票号↔文件名，跨 **76 枚**票；与死腿 `3/pairs93.txt`（146 对／76 枚）逐字节相同＝`diff` 零差 | `pairs-own.txt` |
| 4 | 逐件分类（选定尺，见 §1.3）：读第 1 行标题＋前 8 行免责句 | 146 枚里 **VERDICT 85／OTHER 60／DIR 1**；`MISSING`＝**0**（名册里每枚文件都真在盘） | `file-identity-4.txt` |
| 5 | 逐票三档汇总 | **A=74／B=2／C=17**（和＝93） | `ticket-abc-4.tsv` |
| 6 | `bash reruler-v3.sh` 与死腿 `3/v3-a93.txt` 逐票对照（92 枚公共票） | **零差**；死腿 `rows-draft.tsv` 只有 **92 行**，缺的那一枚＝**268**（本腿补齐为 C/NOFILE） | `reruler-v3.sh` |
| 7 | 尺一：`grep -rnE '^- \[ \]' --include='*-done.md' .scratch/wisp/issues/ \| wc -l` | **14 枚**框，分布在 **8 枚票**（`07`=1／`92`=3／`97`=1／`104`=2／`105`=1／`110`=2／`113`=2／`115`=2）；宽松式（任意缩进）复跑＝同为 14 枚⇒ 本池无缩进框形状差 | `ruler1-14-strict.txt`／`ruler1-14-detailed.txt` |
| 8 | 尺一历史复认：`git grep -hE '^- \[ \]' 2a633eb8^ -- '...*-done.md' \| wc -l` | **15 枚**（含 265 一枚）＝编排者 10-05 17:59 台账 `A624` 那句「同形 …`-done` 票 **15 枚**归 evidence-close-2 逐枚裁」的**当时值**；`2a633eb8` 就地补翻 265／268 各 `numstat 1/1` ⇒ 本腿现量 **14**＝15−1（268 当时不带 `-done`，不入那一分母）。**两数不矛盾，差一枚 265、时点不同** | 同上 |
| 9 | 尺二（票内同行）：`grep -rnE '^- \[ \]' … \| grep -E '翻勾'` | `-done` 票内＝**0 命中**；全池（含开放票）＝**3 行／3 枚票**＝票 182 `:26`、票 242 `:28`、票 245 `:62`（三枚都不带 `-done`，逐枚 `ls` 坐实） | `ruler2-flipclaim.txt` |
| 10 | 尺二加强形：`- [ ]` 那一行的**下一行注记**是否主张"本格已翻／可翻勾"（`` `done-?-1` 勾＝``／`翻勾：`／`⇒ **可翻勾**`） | 全池命中 **1 枚**＝票 **255 `:17`**（⚠ 票 255 **不带 `-done`**，不在本件分母里，只作形状登记）；`-done` 票内＝**0 命中** ⇒ **昨天 265/268 那种"注记写翻勾、物理框没翻"的形状，今天在 93 枚结案票里已经没有了** | 同上 |
| 11 | 编排者自裁表扫描（两发）：对 85 枚 VERDICT 件各取**前 12 行**与**前 30 行**做同一把 `grep -oE`（词组＝`执行者[:：] ?orchestrator`／`Acceptor:\*\* orchestrator`／`编排者执行`／`编排者亲自验收`／`编排者代提`／`由编排者`） | **前 12 行发＝10 枚命中**（多出 `119-ac7-r1-acceptance.md`、`80-ac1-ac2-verdict.md`）；**前 30 行发＝8 枚**＝`05`／`07`／`08`／`11`／`161-gates-selftest-accept-r1`／`17`／`18`／`63`。逐枚现读署名行复认（`05-…md:1`「（编排者执行）」＋`:3`「执行者：orchestrator」；`11-…md:3`、`63-…md:3`「**Acceptor:** orchestrator」）⇒ 净名册见 §3.3-1 | `orch-or-selfdecl.txt`（8 枚版） |
| 12 | 与死腿 `3/orch-executed-verdicts.txt`（**9 枚**）对照＋逐枚读署名语境 | 本腿**剔除 1 枚**＝`161-gates-selftest-accept-r1.md`（`:1`「# 161-v1 **非实现者**验收表…」、`:4`「本程性质＝**非实现者裁决**」，死腿命中的是被验对象来历「161-r3 写码＋编排者代提复跑」）；本腿**另排除 2 枚**（死腿未列、但前 12 行那发扫出来的 `119-ac7-r1`／`80-ac1-ac2-verdict`）＝命中词都是"翻框权归属／上报去向"而非署名 ⇒ 净结论＝**7 枚**编排者自裁（§3.3-1） | 同上 |
| 13 | 三档差集稳定性：`v1-a93`／`v2-a93`／`v3-a93` 的 C 档票号集两两 `diff`；B 档票号集逐版列出 | **C 档三把尺完全相同＝同那 17 枚**（两两 diff 零差）。B 档漂：**v1＝{184, 265}（2 枚）＝v3**、**v2＝{11, 63, 80, 84, 146, 153, 184, 198, 212, 265}（10 枚）**；A 档＝v1 74／v2 66／v3 74 | `v1-a93.txt`／`v2-a93.txt`／`v3-a93.txt`（死腿原件，⛔ 未改）＋本腿 `ticket-abc-4.tsv` |
| 14 | 尺漂移量：逐件分类 `diff` | **v1↔v3＝12 枚件改类**、**v2↔v3＝23 枚件改类**、**v1↔v2＝29 枚件改类**；VERDICT 件总数 v1=83／v2=82／v3=**85**；逐票档 A/B/C 差：v1↔v3 **0 枚票**、v2↔v3 **8 枚票**（全为 B→A） | 本件 §1.2 |
| 15 | `wc -l -c` 本腿台件关键三枚 | `ticket-abc-4.tsv`＝93 行；`file-identity-4.txt`＝146 行；`pairs-own.txt`＝146 行 | — |
| 16 | **独立收敛复认**：把死腿 `2/final-table.tsv`（92 行，`evidence-close-2` 自己那套尺的成品表）与本腿 `ticket-abc-4.tsv` 按票号 join 比档 | **92 枚公共票逐票零差**（`join`＋`$2!=$3` 过滤后无输出）；`2/final-table.tsv` 自带档数＝**A 74／B 2／C 16**（少的那一枚仍是 268，它当时不在分母）。⇒ 两枚腿各写各的脚本、各取各的名册，**同一份账**＝本件最强的复认 | 本件 §1.2／§3.3-3 |
| 17 | 旁证层复认：`2/circumstantial.txt`（18 行，死腿 2 自己那一层"跨名引用"账） | 本腿 §3.3-3 那 5 个形状的**逐票 REFS 清单在死腿 2 原件里独立在册**（`115 REFS=…111-ci-step-readings.md…`、`267 REFS=…265-267-evidence-index.md`＋`PROBES=39`、`250/251/254/263 PROBES=29/19/33/198`、`265 REFS=…ci-step-readings-2026-09-22.md`＋`PROBES=33`）⇒ 本腿不是新下判断，是把死腿已量到、但没进交付件的旁证层**接上来** | 同上 |
| 18 | 免责句层独立决定数复算（可复跑尺 `reruler-disclaim-layer.sh`）：对 146 枚名册逐枚 `head -8 \| grep -E "$DISCLAIM"`，再逐枚把命中件回灌标题层（VERD3／IMPL3） | 前 8 行含免责句＝**6 枚**（`119-ac7-impl`〔`实现侧交件`〕／`137-ac4-impl`〔VERD3 不中〕／`141-ac34-positive-controls-r1`〔`正向对照`〕／`146-liveapprovals-fix-r1`〔`修复程`〕／`146-liveapprovals-r1`〔`实现程自证`〕／`265-267-evidence-index`〔`凭据索引`〕）；逐枚回灌后**全部已在标题层落 OTHER** ⇒「VERD3 中 ∧ IMPL3 不中 ∧ 免责句中」的**独立决定枚数＝0** ⇒ 具名写进 §1.2b，作为对本腿自己选尺理由的一条**诚实修正** | `disclaim-layer-audit.txt`（6 行，逐行标"免责句层无活可接"）＋`reruler-disclaim-layer.sh` |

---

## §1 尺的选定：三版同名重做件的差与漂移（具名）

死腿 `evidence-close-3` 为**同一件事**留了三套同名重做件（`v1-*`／`v2-*`／`v3-*`，各三枚：`a93`＝逐票三档、`i93`／`file-identity`＝逐件分类、`ticket-abc`＝三档），外加它复跑死腿 `evidence-close-2` 的 `classify.sh`／`classify2.sh`／`classify3.sh`。三套并存＝**尺本身有歧义**的信号。本节先把差具名说清，再选一把，⛔ 不把三套平均、⛔ 不混引。

### 1.1 三把尺各是什么口径（三套脚本里逐字读出的差别）

| 版 | 出处脚本 | 读取范围 | 否决方式 | 词表差别 |
|---|---|---|---|---|
| **v1** | `2/classify.sh`（复跑＝`3/rerun-rulers93.sh` 第 1 段） | 只看**第 1 行标题** | 标题命中实现侧词 → OTHER（**一票否决，无免责句层**） | VERD1 词表含裸词 `verdict`／`重判`／`复验`，**不含 `acceptor`／`对抗验收`／`裁决`**；IMPL 词表**窄**（9 词，`证据件 —` 要带破折号才中） |
| **v2** | `2/classify2.sh`（复跑＝同脚本第 2 段） | **前 20 行** | 前 20 行命中实现侧词 → OTHER（**否决面最宽**） | IMPL2 词表**最宽**（19 词，含 `证据件`／`普查`／`写码代理`／`凭据索引`／`读数表`／`预做`）；VERD2 词表**不含裸词 `acceptance`／`verdict`／`裁决`**（改要 `验收腿`／`裁决表`／`对抗验收` 等复合形） |
| **v3** | `2/classify3.sh` ＝ `3/rerun-v3.sh`（死腿最终版） | 第 1 行标题**定类**＋**前 8 行**查免责句 | 标题命中 VERD 且**不**命中 IMPL，再被前 8 行免责句拉回 OTHER | VERD3 含 `acceptance`／`acceptor`／`裁决`／`对抗验收`；IMPL3 21 词；DISCLAIM 7 形（`本件不是验收件`／`不是对抗验收表`／`同一程写码又自证`…） |

### 1.2 漂在哪一枚、漂几枚（逐枚具名，本腿现跑 `diff`）

**件级改类数**（同一枚文件在两把尺下类别不同）＝**v1↔v3＝12 枚／v2↔v3＝23 枚／v1↔v2＝29 枚**；名册层 VERDICT 总数 **v1=83／v2=82／v3=85**（同一份 146 对名册，⛔ 分母相同、只差判词）。

- **`v1` ↔ `v3`＝12 枚**（逐枚具名＋方向）：
  - v1 漏、v3 收（7 枚）＝`119-ac1-ac2-r2-acceptance.md`、`137-ac1-r1-acceptance.md`、`137-ac2-ac5-r1-acceptance.md`、`137-ac3-r2-acceptance.md`、`153-ac2-a-closure-r1.md`、`157-record-level-cleanup-r1-accept-r1.md`、`162-fs-edit-accept-r2.md` ⇒ 本腿把 VERD1 逐词套到标题行现算，根因＝**v1 的 VERD1 词表缺 `回判`／`验收方`／`非实现者` 三个词**（v3 齐）：`119-…` 命中的是 `回判`、137 那三枚命中的是 `验收方`（`:1`「# 137 AC#1 —— **验收方**第二程…」，标题里没有英文 `acceptance`，那只在文件名里）、153／157／162 那三枚命中的是 `非实现者` ⇒ v1 一律落"标题无终裁词"默认出口。
  - v1 收、v3 踢（5 枚）＝`131-reaccept-ac4.md`、`142-non-quiescent-index-guard-r1.md`、`146-liveapprovals-r1-accept-r2.md`、`153-trace-lies-unguarded-r1-accept-r1.md`、`67-ac3-ascii-verdict.md` ⇒ 本腿逐枚把两把尺的词表套到标题行上现算，根因**两条都在标题层**：**① v1 的 IMPL1 词表窄**——`142`（`:1` 逐字「…（**实现方自证**；裁决表须出自非实现者）」，IMPL1 只有 `实现者自证`、**抓不到 `实现方`**）、`146-…-r2`（`:1`「…**修复程**补上的那两味尺子：裁决表」，IMPL1 无 `修复程`）、`153-…-r1-accept-r1`（`:1` 尾含 `只读取证`，IMPL1 无此词）三枚因此被 v1 放行；v3 的 IMPL3 三词齐备 ⇒ 踢掉。**② v1 的 VERD1 含两个 v3 已删掉的裸词**——`复验`（→ `131-reaccept-ac4.md` 命中）与 `verdict`（→ `67-ac3-ascii-verdict.md` 命中，而它 `:1` 其实是"把 `wisp providers` 的 verdict 列改成 ASCII `PASS`/`FAIL`"＝一枚讲工具输出的实现侧件）⇒ v1 因词表多收、v3 因词表精而不收。
    ⚠ **本腿更正死腿 `3/rerun-rulers93.sh` 注释里的一处口径宣称**：这 12 枚的改类**没有一枚是"前 8 行免责句"层造成的**（见下 §1.2b）。
- **`v2` ↔ `v3`＝23 枚**（逐枚具名＋方向；**13 枚被 v2 踢掉／10 枚被 v2 误收**）：
  - **v2 踢、v3 收（13 枚）**＝`11-adversarial-acceptance.md`、`63-adversarial-acceptance.md`、`80-ac1-ac2-verdict.md`、`84-ac1-bounded-wait.md`、`146-liveapprovals-r1-accept-r1.md`、`152-subject-death-never-measured-r1-accept-r1.md`、`153-ac2-a-closure-r1.md`、`154-close-gate-never-rings-r1-accept-r1.md`、`161-gates-accept-r2.md`、`162-fs-edit-accept-r1.md`、`162-risk-tier-accept-v3.md`、`198-firstrun-config-v1.md`、`212-comments-phantom-citation-v2.md`。
  - **v2 收、v3 踢（10 枚）**＝`131-followup-1-three-shapes.md`、`131-reaccept-ac4.md`、`134-ac6-contended-no-conclusion.md`、`137-ac4-impl.md`、`141-ac2-stock-inventory-r1.md`、`142-non-quiescent-index-guard-r1.md`、`143-panel-assets-r4-leg-r1.md`、`151-task-scope-never-closed-r1.md`、`152-subject-death-never-measured-r1.md`、`154-host-id-never-closed-r1.md`。
  - **根因两族，本腿逐枚复跑两套词表坐实**：**①族＝v2 的 VERD2 词表不含裸词 `acceptance`／`verdict`**（3 枚：`11-…md`＝`:1` 全英文「# Ticket 11 — adversarial acceptance」、`63-…md` 同形、`80-ac1-ac2-verdict.md`＝标题有"裁决"二字但前 20 行两族命中**都是空**，本腿现量 `IMPL2hit=[]`＋`VERD2hit=[]`）⇒ 直接落 OTHER（"两个都不中"＝OTHER 的默认出口）。
    **②族＝v2 把前 20 行里"描述被验对象"的实现侧词当成件自己的自称**（10 枚：`84`〔命中 `写码代理`〕、`146-…-r1`〔`普查`+`证据件`〕、`152-…-accept-r1`〔`只读取证`〕、`153-ac2-a-closure-r1`／`154-…-accept-r1`／`162-risk-tier-accept-v3`／`198-v1`／`212-v2`〔各命中 `证据件` 1 词〕、`161-gates-accept-r2`〔`普查`〕、`162-fs-edit-accept-r1`〔`写码位`+`证据件`〕）——这几枚的标题逐字都写着「**对抗验收**…**裁决表**」／「**非实现者**验收程」／「**非实现者**腿」，VERD2 其实**全中**（本腿现量逐枚列出，见台件），却被"前 20 行否决"这道**先判**吃掉 ⇒ **一票否决踢掉 10 枚真表**。v3 的否决面只有第 1 行标题，正好避开这两族误伤（这也是 v3 名册 VERDICT 85 ＞ v2 的 82 的全部来源）。
- **`v1` ↔ `v2`＝29 枚**＝两把尺各朝相反方向漂（v1 的漏收＋v2 的过否决叠加），互校时漂得最多 ⇒ 说明"v1 与 v2 不是同一把尺的两个版本，是两把不同的尺"。

### 1.2b 免责句层在本名册上**独立决定数＝0**（对本腿选尺理由的一条诚实修正，具名）

派单提示"尺本身有歧义"，本腿逐枚验到一件比那更具体的事：**v3 比 v1 多出来的那道"前 8 行免责句"否决，在这 146 枚名册上一枚都没有独立判掉过。**现量三步：

1. 全名册 146 枚里，前 8 行含免责句的只有 **6 枚**：`119-ac7-impl.md`（`本件不是裁决表`）、`137-ac4-impl.md`（同）、`141-ac34-positive-controls-r1.md`（`本文件不是裁决表`）、`146-liveapprovals-fix-r1.md`（`本件不是验收件`）、`146-liveapprovals-r1.md`（`同一程写码又自证`＋`本件不是验收件`）、`265-267-evidence-index.md`（`不是对抗验收表`）。
2. 把这 6 枚逐枚再过一遍标题层：`119-ac7-impl`／`141-ac34-positive-controls-r1`／`265-267-evidence-index` 三枚标题即命中 IMPL3（`实现侧交件`／`正向对照`／`凭据索引`）；`146-liveapprovals-fix-r1`／`146-liveapprovals-r1` 两枚标题命中 IMPL3（`修复程`／`实现程自证`）；`137-ac4-impl.md` 两族都不中 ⇒ 走"标题无终裁词"的默认出口落 OTHER。**⇒ 6 枚全部在标题层就已经被判掉，免责句层接手时手上已经是空的。**
3. 反向也量：全名册里满足「VERD3 中 ∧ IMPL3 不中 ∧ 免责句中」的**独立决定枚数＝0**。

这条不是挑 v3 的毛病，是**把 v3 的真实承重结构说清**，防下一位误引：v3 相对 v1 的优势＝**IMPL3 词表宽 12 词＋VERD3 删掉了 `verdict`／`复验`／`重判` 三个裸词**（＝上面那 12 枚的全部来源）；v3 相对 v2 的优势＝**否决面从"前 20 行"收到"仅第 1 行"**（＝避开 ②族那 10 枚误伤）。免责句层是一道**今天没有出过力的保险**——留着（它防的是"标题合规、开头就自陈不是裁决表"这种尚未出现的形状），但⛔ 不许把它当成"v3 比 v1 严"的理由来引，那是错的。

⚠ 同一件事也解释了 §3.3-2 那个盲点为什么**必然**存在：**整把尺的效力全压在"第 1 行"这一行字上**，所以票 84 那枚「`:1` 合规、`:3` 自称写码代理」的件一定漏。要补这一层＝改尺＝按 `AGENTS.md` §0.2「改契约＝人工批准」，⛔ 本腿不改，只具名并把口径写死在原处供裁。

**票级换档数**（A/B/C 三档里跳到别的档）：

- `v1` ↔ `v3`＝**0 枚票**（§0.2 那把 join 尺现量零命中）。那 12 枚件改类之所以**动摇不了档位**，是因为涉及的 9 枚票（`119 131 137 142 146 153 157 162 67`）名下**各有另一枚两把尺都认的终裁件**兜着——⇒ 换 A→A（v1 A＝74＝v3 A＝74）。这条同时也说明：**12 枚的件级漂在票级被吸收了，所以"漂了几枚"必须两粒度都报，只报件级会吓人、只报票级会掩盖。**
- `v2` ↔ `v3`＝**8 枚票**，方向**全部 B→A**（逐枚具名）：`11`、`63`、`80`、`84`、`146`、`153`、`198`、`212` ⇒ 根因就是上面那两族（这 8 枚票**名下只有那一张表**，v2 把它踢掉后票就掉进 B）。
- **C 档（零凭据）三把尺完全一致＝同那 17 枚**（`v1`/`v2`/`v3` 的 C 票号集两两 `diff`＝零差，§0.2 row 13）⇒ **"哪些票在 `docs/evidence/s1/` 里根本没有同名件"这件事不随尺漂**，漂的只是"名下有件、那件算不算终裁表"。⇒ 这条是本腿选尺决策的**承重墙**：三段计数里最需要稳的那一段（17 枚无凭据）三把尺给的是同一个答案，⛔ 换尺不会改变"要派哪些腿"。

### 1.3 本腿选定：**v3**（死腿最终版口径），并把它写死成下面这一段

```
同名件 := ls docs/evidence/s1/ | grep -E "^0*<票号>([-/]|$)"
          ——票号必须在文件名开头，且紧跟 `-`、`/` 或行尾。
          ⇒ 136-ac14-* 归票 136（不是票 14）；245-esc-* 归票 245（不是票 45／票 5）。
          ⇒ README.md 永不进分母：分母只由 `ls .scratch/wisp/issues/*-done.md` 抽票号，
            `README.md` 不带 `-done` 后缀，形状上就进不来（本腿现量：分母 93、其中无 README）。
VERDICT := 第 1 行标题命中 VERDICT 词族 且 标题行不命中实现侧自称词族
            且 前 8 行没有明写的免责句（"本件不是验收件／不是裁决表／同一程写码又自证"）。
OTHER   := 其余同名件（实现侧交件／普查／读数台件／取证件／索引／目录）。
三档    := A 名下有 VERDICT ／ B 名下只有 OTHER ／ C 名下零枚同名件（NOFILE）。
```

选它的理由（三条，都是可核的；⚠ 第 1 条经 §1.2b 修正过口径）：

1. **v3 的两处口径都恰好在"这一名册真正会咬人的地方"**：① 它的**否决面只在第 1 行标题**，而 v2 把否决面铺到前 20 行 ⇒ 把标题逐字写着"对抗验收／裁决表／非实现者"、只是正文顺带提到"证据件／普查"的 **10 枚真表误踢**，连带 **8 枚票掉档**（§1.2 ②族）；② 它的 **IMPL3 词表比 v1 宽 12 词、且 VERD3 删掉了 `verdict`／`复验`／`重判` 三个裸词** ⇒ 避开 v1 那 5 枚多收（§1.2 v1↔v3 ①②两条）。
   ⛔ **但本腿不能把"v3 多一道免责句否决"当理由引**：§1.2b 现量＝那道保险在这 146 枚上**独立决定数＝0**（6 枚含免责句的件全部已在标题层被判掉）。v3 比 v1 严的说法**不成立**，两把尺的差**全在词表**。
2. **v3 是死腿自己的最终版**，且 `rerun-v3.sh` 的头部注释逐字把口径写成散文（三套里唯一把口径写清楚的），本腿能逐字复跑。
3. **复现性**：本腿不引用死腿任何一枚输出，只拿它的脚本口径自己重抽名册（`pairs-own.txt` 146 对）、自己分类（`file-identity-4.txt`）、自己汇总（`ticket-abc-4.tsv`），结果与 `3/v3-a93.txt` **逐票零差**（§0.2 row 6）。

**另三条口径声明（本腿自己拍的，写死在这里，防下一位再漂）**：

- **A 档不看裁决者是谁**。v3 尺只判"这枚件**形状上**是不是一张终裁表"。本件 §3.3-1 另把**编排者自裁那 7 枚**具名列出，作为 A 档里"需打折"的一层注记——⛔ 不因此把它们挪出 A 档（挪了就等于用第二把尺重洗名册，正是本腿要避免的事）。
- **`DIR` 算 OTHER**（`docs/evidence/s1/66` 是一枚**目录**，是台件不是裁决表）。
- **⛔ 本尺只认文件名前缀，不认"内容里提到了哪枚票"**。这直接导致 §3.3 那批"有非实现者凭据、但凭据不在自己名下"的票落进 C 档——这是尺的性质，不是那些票的性质。本件在 §3.3 逐枚具名登记这类凭据，⛔ 不改档。

---

## §2 分母（本腿现量）

| 项 | 读数 |
|---|---|
| 尺 | `ls .scratch/wisp/issues/*-done.md \| wc -l` |
| 本腿现量（`2026-10-06 09:12`） | **93** |
| 取票号形状 | `basename` 后取第一个 `-` 之前的段，保留前导零（`01`…`268`），`sort -u`＝93 枚，无重号 |
| `README.md` | **不在分母里**（它不带 `-done` 后缀）。⛔ 前一腿曾因把它算进去而多计一枚 |
| 台件 | `.scratch/wisp/probes/evidence-close/4/nums4.txt`（93 行）／`done-files-4.txt`／`denominator.txt` |
| 会漂声明 | 分母随结案动作单调上漂。派单给的 10-06 09:0x＝**93**＝本腿现量一致；写本件时以本腿为准，⛔ 不引任何腿的历史值当现值 |

**分母历史复认（三枚腿的取数差，具名）**：`3/live-done-nums.txt`＝92、`3/live-done-nums-v2.txt`＝93、`3/nums93.txt`＝93、`2/done-nums.txt`＝92、`3/rows-draft.tsv`＝92 行。逐枚 `diff` 后的**唯一差集＝票 268 一枚**：`git log --diff-filter=R` 现量 268 改名 `-done` 那一笔＝**`11105265`（2026-10-05 22:56，`R075`）**，而死腿那枚 92 枚的取数发生在 21:47 ⇒ **不是尺错，是取数时刻早于结案**。`rows-draft.tsv` 少 268 同因。本腿分母含 268（§3.3 里它的凭据形态另说）。

---

## §3 三段计数与逐票两列

### 3.1 三段计数（选定尺＝v3，分母＝本腿现量 93）

| 档 | 含义 | 枚数 | 占比 |
|---|---|---|---|
| **A** | 名下有非实现者终裁表（形状判定，裁决者身份见 §3.3） | **74** | 79.6% |
| **B** | 名下只有旁证／索引件（`OTHER`），零枚终裁表 | **2** | 2.2% |
| **C** | 名下完全无凭据（`NOFILE`＝零枚同名件） | **17** | 18.3% |
| 合计 | | **93** | 100% |

名册层：146 对同名件里 **VERDICT 85 枚／OTHER 60 枚／DIR 1 枚**（85 枚终裁件分布在 74 枚票名下，5 枚票名下有 ≥2 张表）。

### 3.2 逐票两列（票号｜命中文件名或 `NOFILE`）

| 票号 | 档 | 命中文件名（`docs/evidence/s1/` 名下，本腿现量；⛔ 不判该不该结） |
|---|---|---|
| 01 | C | `NOFILE` |
| 02 | C | `NOFILE` |
| 03 | C | `NOFILE` |
| 04 | C | `NOFILE` |
| 05 | A | 05-adversarial-acceptance.md ⚠自裁 |
| 06 | C | `NOFILE` |
| 07 | A | 07-adversarial-acceptance.md ⚠自裁 |
| 08 | A | 08-adversarial-acceptance.md ⚠自裁 |
| 09 | A | 09-adversarial-acceptance.md |
| 10 | A | 10-adversarial-acceptance.md |
| 11 | A | 11-adversarial-acceptance.md ⚠自裁 |
| 13 | C | `NOFILE` |
| 14 | C | `NOFILE` |
| 17 | A | 17-adversarial-acceptance.md ⚠自裁 |
| 18 | A | 18-adversarial-acceptance.md ⚠自裁 　·　旁证件 18-perf-bench-addendum.md |
| 19 | C | `NOFILE` |
| 63 | A | 63-adversarial-acceptance.md ⚠自裁 |
| 66 | A | 66-adversarial-acceptance.md 　·　旁证件 `66`（目录）66-mutation-parse-order.md |
| 67 | A | 67-adversarial-acceptance.md 　·　旁证件 67-ac3-ascii-verdict.md 67-emoji-scope-internal-cmd.md |
| 69 | A | 69-adversarial-acceptance.md 　·　旁证件 69-mutation-tokens.md |
| 73 | A | 73-adversarial-acceptance.md 　·　旁证件 73-orphan-staging-sweep.md |
| 74 | A | 74-adversarial-acceptance.md |
| 75 | A | 75-independent-verification.md 　·　旁证件 75-linux-mutation-check.md 75-rootcause-linux-path-shape.md 75-rootcause-locator-report.md |
| 76 | A | 76-adversarial-acceptance.md |
| 78 | C | `NOFILE` |
| 79 | A | 79-adversarial-acceptance.md |
| 80 | A | 80-ac1-ac2-verdict.md |
| 81 | A | 81-adversarial-acceptance.md |
| 82 | A | 82-adversarial-acceptance.md |
| 83 | A | 83-adversarial-acceptance.md |
| 84 | A | 84-ac1-bounded-wait.md 〔§3.3-2 存疑〕 |
| 87 | A | 87-adversarial-acceptance.md |
| 88 | A | 88-adversarial-acceptance.md |
| 89 | A | 89-adversarial-acceptance.md |
| 90 | A | 90-adversarial-acceptance.md |
| 92 | A | 92-adversarial-acceptance.md |
| 93 | A | 93-adversarial-acceptance.md |
| 94 | A | 94-adversarial-acceptance.md |
| 95 | A | 95-adversarial-acceptance.md |
| 96 | A | 96-adversarial-acceptance.md |
| 97 | A | 97-adversarial-acceptance.md |
| 99 | A | 99-adversarial-acceptance.md |
| 101 | A | 101-adversarial-acceptance.md |
| 102 | A | 102-adversarial-acceptance.md |
| 104 | A | 104-adversarial-acceptance.md |
| 105 | A | 105-adversarial-acceptance.md |
| 106 | A | 106-adversarial-acceptance.md |
| 110 | A | 110-adversarial-acceptance.md |
| 113 | A | 113-adversarial-acceptance.md |
| 115 | C | `NOFILE`（凭据在他名下，见 §3.3-3） |
| 116 | A | 116-adversarial-acceptance.md |
| 117 | A | 117-adversarial-acceptance.md |
| 118 | A | 118-adversarial-acceptance.md |
| 119 | A | 119-ac1-ac2-r2-acceptance.md 119-ac2-r2b-acceptance.md 119-ac7-r1-acceptance.md 119-adversarial-acceptance.md 　·　旁证件 119-ac7-impl.md 119-next-audit-r1.md |
| 121 | A | 121-adversarial-acceptance.md 　·　旁证件 121-preflight-deps-reachability.md |
| 126 | A | 126-adversarial-acceptance.md |
| 127 | A | 127-adversarial-acceptance.md |
| 129 | A | 129-adversarial-acceptance.md 　·　旁证件 129-ac4-ac5-mutation-and-gates.md |
| 131 | A | 131-adversarial-acceptance.md 　·　旁证件 131-followup-0-crossplatform.md 131-followup-1-three-shapes.md 131-followup-2-ratification-and-linux.md 131-reaccept-ac4.md |
| 134 | A | 134-adversarial-acceptance.md 　·　旁证件 134-ac2-ac3-trigger-and-freshness-pin.md 134-ac4-machine-contended-readings.md 134-ac5-gate-readings.md 134-ac6-contended-no-conclusion.md |
| 137 | A | 137-ac1-r1-acceptance.md 137-ac2-ac5-r1-acceptance.md 137-ac3-r1-acceptance.md 137-ac3-r2-acceptance.md 137-ac4-r1-acceptance.md 　·　旁证件 137-ac1-teeth-or-not.md 137-ac2-ac5-impl.md 137-ac4-impl.md |
| 141 | A | 141-q46c-accept-r1.md 　·　旁证件 141-ac2-stock-inventory-r1.md 141-ac34-positive-controls-r1.md 141-q46c-impl.md |
| 142 | A | 142-non-quiescent-index-guard-accept-r1.md 　·　旁证件 142-non-quiescent-index-guard-r1.md |
| 143 | A | 143-panel-assets-r4-leg-accept-r1.md 　·　旁证件 143-panel-assets-r4-leg-r1.md |
| 144 | A | 144-slo-report-partial-read-r1-accept-r1.md 　·　旁证件 144-slo-report-partial-read-r1.md |
| 146 | A | 146-liveapprovals-r1-accept-r1.md 　·　旁证件 146-liveapprovals-fix-r1.md 146-liveapprovals-r1-accept-r2.md 146-liveapprovals-r1.md |
| 147 | A | 147-offset-naming-r1-accept-r1.md 　·　旁证件 147-offset-naming-r1.md |
| 149 | A | 149-three-unpinned-outlets-r1-accept-r1.md 　·　旁证件 149-three-unpinned-outlets-r1.md |
| 151 | A | 151-task-scope-never-closed-r1-accept-r1.md 　·　旁证件 151-task-scope-never-closed-r1.md |
| 152 | A | 152-subject-death-never-measured-r1-accept-r1.md 　·　旁证件 152-subject-death-never-measured-r1.md |
| 153 | A | 153-ac2-a-closure-r1.md 　·　旁证件 153-trace-lies-unguarded-r1-accept-r1.md 153-trace-lies-unguarded-r1.md |
| 154 | A | 154-close-gate-never-rings-r1-accept-r1.md 　·　旁证件 154-host-id-never-closed-r1.md |
| 155 | A | 155-three-unjudged-cells-r1-accept-r1.md 　·　旁证件 155-three-unjudged-cells-r1.md |
| 156 | A | 156-close-and-load-bearing-r1-accept-r1.md 　·　旁证件 156-exited-asks-os-r1.md 156-r4-gates-and-cleanup-r1.md |
| 157 | A | 157-ac5-recheck-r1.md 157-record-level-cleanup-r1-accept-r1.md |
| 161 | A | 161-gates-accept-r2.md 161-gates-selftest-accept-r1.md 　·　旁证件 161-attrib-selftest-r5.md 161-doorbell-census-r1.md 161-gate-aggregate-r6.md 161-gate-blindspot-r1.md 161-gate-order-r4.md 161-selftests-r3.md |
| 162 | A | 162-fs-edit-accept-r1.md 162-fs-edit-accept-r2.md 162-risk-tier-accept-v3.md 　·　旁证件 162-fs-edit-r1.md 162-fs-edit-r2.md 162-fs-edit-r3.md 162-risk-tier-and-contract-axis-r4.md |
| 183 | A | 183-pointer-exemption-accept-v2.md 　·　旁证件 183-exemption-why-not-live-a1.md 183-pointer-exemption-r1.md 183-reread-permanent-teeth-r2.md |
| 184 | B | 184-audit-vocabulary-census-c1.md（＝词汇普查件，非表） |
| 198 | A | 198-firstrun-config-v1.md 　·　旁证件 198-firstrun-config-r2.md |
| 212 | A | 212-comments-phantom-citation-v2.md |
| 221 | A | 221-task-cancel-v1.md |
| 235 | A | 235-pool-nail-and-bridge-preread-v1.md |
| 241 | A | 241-audio-level-producer-v1.md 　·　旁证件 241-audio-level-producer-r1.md |
| 243 | C | `NOFILE` |
| 250 | C | `NOFILE` |
| 251 | C | `NOFILE` |
| 254 | C | `NOFILE` |
| 257 | A | 257-clean-machine-provider-registry-v1.md |
| 263 | C | `NOFILE` |
| 265 | B | 265-267-evidence-index.md（＝前一腿 `evidence-close-1` 出的**索引件**，非表） |
| 267 | C | `NOFILE`（凭据在他名下，见 §3.3-3） |
| 268 | C | `NOFILE`（凭据在盘但不落 `s1/`，见 §3.3-3） |

⛔ 本表只列"文件名"，不列判语、不给分数、不标"可结案"。表里**没有**任何一枚票被本件视为已验收。

### 3.3 A 档里的三处水分（具名，⛔ 不改档）

选定尺只判"形状像不像终裁表"，因此下面三类必须在本节**另层具名**，否则这张索引会被误读成"74 枚票都有非实现者终裁"。

1. **编排者自裁那 7 枚**（§0.2 row 11 命中 8 枚 → row 12 逐枚复核**剔除 1 枚** → 另具名**排除 2 枚误计** → 净 **7**；标题／署名行逐字自陈执行者＝orchestrator）：
   `05-…md:1`「# T05 对抗验收报告（**编排者执行**）」＋`:3`「执行者：orchestrator」／`07-…md:1`+`:3` 同形／`08-…md:3`「执行者：orchestrator（实现者 T08-impl 独立）」／`11-…md:3`「**Acceptor:** orchestrator」／`17-…md:1`+`:3`／`18-…md:1`+`:3`／`63-…md:3`「**Acceptor:** orchestrator (not the implementer)」——`63` 那句自陈"不是实现者"，但裁决者＝**编排者本人**：与实现腿不同体，与 `AGENTS.md` §0.3 要求的"**另一个** agent"仍不同体。
   ⚠ **剔除 1 枚**＝死腿 `3/orch-executed-verdicts.txt` 列在第 9 位的 `161-gates-selftest-accept-r1.md`：该件 `:1` 逐字「# 161-v1 **非实现者**验收表…」、`:4`「本程性质＝**非实现者裁决**」；死腿命中的是它 `:4` 里「161-r3 写码＋**编排者代提**复跑（AC#2）」那句**被验对象的来历**，不是裁决者署名 ⇒ 判死腿这一枚为**误计**（本腿复跑口径改成"署名位词组"后它就不中了）。
   ⚠ **另排除 2 枚**＝同一份死腿名册里的 `119-ac7-r1-acceptance.md`（`:1`「# 119 AC#7 —— 裁决表 r1（**非实现者**）」，命中词来自 `:8`「本表**不翻 AC#7 的勾、不改票名**…勾由编排者按本表打」＝**翻框权归属**陈述）与 `80-ac1-ac2-verdict.md`（命中词来自 `:12`「停下上报编排者，由编排者判是否 D22」＝**上报去向**）。两枚都是真·非实现者表，⛔ 不该混进自裁名册。
   ⇒ 那 **7 枚票**（`05 07 08 11 17 18 63`）的 A 档含义＝**有一张表，但表出自编排者之手**。本腿把它们留在 A 档（尺不管署名）并在 §3.2 逐枚标 `⚠自裁`。
2. **`84-ac1-bounded-wait.md`＝尺的一个真盲点**。标题「# 票 84 裁决表 — AC#1 两侧读数…」命中终裁词、不命中实现侧词 ⇒ v3 尺判 VERDICT；而死腿 `3/deep-selfcheck.txt` 现读其正文自陈「**写码代理**，2026-09-21」——本腿复跑坐实该句在 **`:3`**（「写码代理，2026-09-21。用例：`internal/agent/approval/ticket84_no_owner_test.go`…」）。⇒ 它在**前 8 行内**，但 v3 尺在前 8 行只查**免责句词族**（`本件不是验收件`…），"写码代理"**不在 DISCLAIMER 词表里**（只在 IMPLTITLE 词表、而 IMPLTITLE 只作用于第 1 行）⇒ 结构性漏抓。具名登记为尺的盲点，⛔ 不改档（请不请出名册＝裁决者／编排者的判断，不是索引的判断）。
   ⚠ 同族另三枚（`221`／`235`／`241`，死腿 `deep-selfcheck.txt` 一并标了"实现腿"）**本腿判为过度怀疑**：三枚 `:1` 逐字分别是「# 221 — **非实现者**对抗验收表（腿 `221-v1`）」「# 票 235 裁决表 —— 235-v1（**非实现者**验收腿）」「# 241-v1 — **裁决腿**验收…」，正文里的"实现腿"都在**指被验对象**（`241-…md` 甚至逐字写「实现腿：`241-r1`（**本腿不信它的自述**）」）⇒ 与 84 那枚"署名为写码代理"形状不同，本腿保留它们且不标存疑。
3. **C 档不等于"没做过工作"**。17 枚 C 里有 **5 个具名形状**（涉及 8 枚票）的凭据**真存在、只是不落在 `docs/evidence/s1/<票号>-*` 这个名字＋落点形状里**：

| 形状 | 票 | 实际在盘的凭据（本腿现量） | 为什么落 C |
|---|---|---|---|
| ① | **115** | `docs/evidence/s1/ci-step-readings-2026-09-22.md`（一枚非实现者 CI 读数表）；票面 AC#2／AC#3 的 `done-check-1` 注记逐字引它 `:241`／`:243`／`:245` 三行 run 读数，并**自己承认**「本腿的凭据**不是**"115 名下的表"…`find docs/evidence -name '115*'` 本腿现量＝**0**」 | 文件名不以 `115-` 开头 ⇒ 尺看不见它。票 230 `:19` 明确要求"115 要一张真名下的表"，那张表**至今 0 枚**（本腿 `ls docs/evidence/s1/ \| grep -iE '^115'`＝零命中，exit=1；票 230 自身**不带 `-done`**＝还开放） |
| ② | **267** | `docs/evidence/s1/265-267-evidence-index.md`（前一腿 `evidence-close-1` 的**索引件**，`:1` 逐字「# 凭据索引 — 票 265 与票 267…」）＋别票名下表反向引用它（死腿 `3/middle-bucket-raw.txt` 记 `267 xrefs=2`） | `265-267-…` 以 `265` 开头 ⇒ 尺把它记在**票 265** 名下（265 因此落 **B** 而非 C），票 267 名下**零枚**。前一腿 §2.5 已具名：`267-v1` 那张真表此刻**还没出** |
| ③ | **268** | `.scratch/wisp/probes/268/a2/verdict.md`＝**205 行／26,456 字节**（本腿现量），`:1`「# 268-a2 独立复核：对 `268-a1/census.md` 复跑 AC#0 四问那把尺」、`:3`「编队：`268-a2`（**非实现者腿，裁决腿**）」；票 268 AC#0 的翻勾注记逐字引它 | 它**不在 `docs/evidence/s1/` 里**（落在 `.scratch` 台件区），尺的射程只有 `s1/` ⇒ 落 C。⚠ 这是**凭据落点不合 `AGENTS.md` §1.5**（裁决表落点＝`docs/evidence/s1/`），不是"没凭据"。本腿判这一条为**结构性缺陷** |
| ④ | **243** | `.scratch/wisp/probes/243/**` 1 枚台件；死腿 `3/xref-final.txt` 另记 `228-resident-ball-v1.md` 反向引用它（`243 xrefcount=1 probesfiles=1`） | 同上（不在 `s1/`＋名字不对＋那枚台件是实现侧料） |
| ⑤ | **250／251／254／263** | `.scratch/wisp/probes/<票号>/**` 分别 **29／19／33／198** 枚台件（死腿 `3/xref-final.txt` 现量＋本腿 `find` 复量＝同数） | 同上：**实现侧普查台件**，本就不是终裁表，且不在 `s1/` |

   另 **9 枚** C（`01 02 03 04 06 13 14 19 78`）＝**只有"别票名下的表反向提到它"**（死腿 `3/middle-bucket-raw.txt` 逐枚列了引用清单：`02`=2 条、`03`=3、`04`=5、`06`=6、`13`=3、`14`=5、`19`=9、`78`=6、`115`=8）；本腿 `find` 复量：这 9 枚名下连 `.scratch/wisp/probes/<票号>/` 台件目录都是 **0 个**（`01`/`02`/`03`/`04` 命中数 24／35／37／19 是**跨目录模糊匹配**到的别腿散件，不是它们名下的裁决件）。⇒ **真无凭据**，属 S1 早期那批"结案时没有非实现者腿"的历史形状。

⇒ 本节一句话：**A=74 不能被读成"74 枚票过了对抗验收"**（7 枚的表出自编排者、1 枚〔票 84〕正文有实现侧自称）⇒ 真·另一体终裁＝**66 枚**；**C=17 里有 5 个形状盘上其实有东西**，只是不在尺的落点上——票 **268** 是**裁决件落错目录**、票 **267**／**115** 是**凭据挂在别票名下**。

---

## §4 对账清单（给编排者翻框用）

⛔ **本节一枚框都不翻**。翻框归编排者。本腿要防的事已经发生过：10-05 在票 265／268 抓到"注记写翻勾、物理框没翻"（现由 `2a633eb8` 就地补翻，见 §0.2 row 8），那是**账面洞**，不是本腿该替编排者补的洞。

判读只有两档：**〔漏翻物理框〕**＝凭据齐、只差把 `[ ]` 写成 `[x]`；**〔真未闭合〕**＝缺非实现者凭据（缺的是**读数**，不是那个勾）。

### 4.1 尺一：`-done` 票里物理未勾的框，逐枚（本腿现量＝**14 枚／8 枚票**）

尺＝`grep -rnE '^- \[ \]' --include='*-done.md' .scratch/wisp/issues/`（顶格形；任意缩进复跑同为 14 ⇒ 本池无缩进形状差）。台件＝`ruler1-14-strict.txt`／`ruler1-14-detailed.txt`。

| # | 票 | 框行 | AC | 行前 60 字 | 票内注记（file:line）＋台账／表凭据（file:line） | 判读 |
|---|---|---|---|---|---|---|
| 1 | 07 | `:51` | 无 AC 号（自由式 Visual 格） | `- [ ] Visual: 20 states rendered in a debug cycle page/window; h` | 注记 `.scratch/wisp/issues/07-ball-state-machine-core-done.md:53`＝**丁类·待人项，⛔ 不翻勾**（逐字要"owner 一次眼"，第二分句凭据在被禁读的 `design/**`）；台账 `A440` 第④节第 1 条在册；表 `docs/evidence/s1/07-adversarial-acceptance.md:10`「视觉证据 \| **PASS（人工签收挂起）**」＋`:14`「**待人工项**：20 态视觉的主观签收（**用户本人**）」 | **〔真未闭合〕**＝等 owner 一眼；⛔ 不该翻 |
| 2 | 92 | `:68` | AC#5 | `- [ ] **AC#5** 变异三向：(i) 把"面板只能显示"改成"面板能写 mode" ⇒ 必须有用例红；` | 注记 `92-panel-composer-mode-attachments-workspace-done.md:72`＝**己类·判不了，本腿不翻勾**；表 `docs/evidence/s1/92b-adversarial-acceptance.md:137` 逐字「(iii) **未重做**」＋「(iii)〔仅自述，不背书〕」 | **〔真未闭合〕**＝(iii) 那一向缺一次非实现者独立复跑 |
| 3 | 92 | `:73` | AC#6 | ``- [ ] **AC#6** 台账与门禁（只跑自己碰的范围）：`sh scripts/d22scan.sh` 纯净树 rc=0 且贴出**逐作用域文件数**`` | 注记 `92-…-done.md:78`＝**己类·判不了**；表 `92b-…md:138`「**通过（数字全复算）**」〔独立复现〕，**同一行自留 `R-92b-4`「⚠ 整步 `portable-tests.sh` 它未跑」**（同表 `:290` 记为**未答问**、`:223` 记「门禁压根没跑 ⇒ 命中一发」） | **〔真未闭合〕**＝那批数只对 `91b5fc4`／`a8f9459` 两棵树负责，当前 HEAD 缺＋`R-92b-4` 未答 |
| 4 | 92 | `:79` | AC#7 | ``- [ ] **AC#7** **负判据**：把"不做 git 切换"变成可检查的东西——在 `frontend/` 与 `internal/panel/` 里`` | 注记 `92-…-done.md:81`＝**丁类·本腿碰不到全格**；表 `92b-…md:139`「**通过**〔独立复现〕」但其读数范围逐字＝包内用例＋`fixtures/`＋`dist/`，**跨 `frontend/` 那一半无独立复跑记录**；`docs/evidence/s1/92-adversarial-acceptance.md:206` 另记 `R-92-6`＝两类孪生门都 `SkipDir` 了 `fixtures`/`dist` | **〔真未闭合〕**＝需要一枚有 `frontend/**` 读权的腿补那半发 grep |
| 5 | 97 | `:55` | AC#5 | ``- [ ] **AC#5** 门禁（按包）：`gofmt -l`/`gofumpt -l` 空、`go vet ./internal/agent/approval/` rc=0、`` | 注记 `97-dead-strict-param-and-the-comment-that-invents-a-caller-done.md:60`＝**己类·判不了**；表 `docs/evidence/s1/97-adversarial-acceptance.md:17`「**通过**（PASS 相加≠RUN 的写法要更正，数字自洽）」〔独立复现〕，同表 `:41`／`:50` 判「**通过（附条件）**」并留 `R-97-1`／`R-97-2`／`R-97-3`。⚠ **票面注记引的行号是 `18`，本腿现读 `:18`＝空行、真格在 `:17`**＝引用行号大 1（账面小歪，不改判性） | **〔真未闭合〕**＝五发门禁读数在当前 HEAD 缺＋`R-97-3` 未销 |
| 6 | 104 | `:52` | AC#3 | `- [ ] **AC#3** 双向变异：① 把检测退回"只看显式 ACE" ⇒ AC#1 红；` | 注记 `104-sealfile-silently-drops-inherited-grants-done.md:56`＝**己类·判不了，本腿不翻勾**；表 `docs/evidence/s1/104-adversarial-acceptance.md:115` 逐字「**结论：通过。〔独立复现〕**（票面要求的三腿变异 + 我自加两发全部我自己在 `/tmp` 仓外快照重抽、重打、重量）」 | **〔真未闭合〕**＝三腿＋两发在当前 HEAD 纯净快照上缺一次非实现者复跑。⛔ 表里那句"通过"是**当时快照**的通过，本腿不拿它当翻框凭据 |
| 7 | 104 | `:57` | AC#4 | ``- [ ] **AC#4** 回归：`go test -count=2 ./internal/winsec/ ./internal/memory/ ./internal/secret/ ./internal/agent/` rc=0，`` | 注记 `104-…-done.md:62`＝**己类·判不了**；表 `104-…md:162` 逐字「**结论：通过。〔独立复现〕**（每个数字都是我在 `/tmp/ac104-fk`（**`4d43447` 纯净树**）本机重跑的，未照抄自述）」 | **〔真未闭合〕**＝那五发**只对 `4d43447` 那棵树负责**，当前 HEAD 缺 |
| 8 | 105 | `:56` | AC#5 | ``- [ ] **AC#5** 门禁（按包 scope）：`go test -count=2 -v` 各包 rc=0 并逐条点名 SKIP/FAIL（**报 `=== RUN` 行数 == 不同测试名 × 2**；`` | 注记 `105-c26-rewrite-account-has-no-production-reader-done.md:59`＝**己类·判不了**；表 `docs/evidence/s1/105-adversarial-acceptance.md:151`「## AC#5 — 门禁（按票面原文跑）」起逐包给数（`:158` `internal/tools`、`:159` `internal/risk` RUN=328／SKIP=2），口径自陈「全部跑在 `f6818f2` 的**纯净快照**」；总判 `:203`「**FAIL-退回（只退 AC#3 这一格，其余四格通过）**」⇒ AC#5 落在"其余四格通过"里 | **〔真未闭合〕**＝缺的不是裁决而是**时效**：五发＋各 scope 文件数在当前 HEAD 上的一次非实现者读数 |
| 9 | 110 | `:39` | AC#3 | ``- [ ] **AC#3** 这道新步要**自己会红**：在 `/tmp` 快照里把 winsec 某条安全断言人为弄坏（例如私有集改成按名字比）⇒ 新步必须 rc≠0；`` | 注记 `110-no-ci-step-runs-internal-winsec-done.md:42`＝**己类·判不了**；表 `docs/evidence/s1/110-adversarial-acceptance.md:16`「**通过（三发全部独立复现）**」〔独立复现〕，同格「全部在 `/tmp/wisp-ac110-snap`」 | **〔真未闭合〕**＝三发变异在当前 HEAD 纯净快照上缺非实现者复跑 |
| 10 | 110 | `:47` | AC#5 | ``- [ ] **AC#5** 门禁：`bash -n` 改动脚本 rc=0；`sh scripts/d22scan.sh` 纯净快照 rc=0 且台账各 scope 不降；`` | 注记 `110-…-done.md:49`＝**己类·判不了**；表 `110-…md:18`「**通过（附一条后果登记）**」〔独立复现〕＋同一行末逐字「⚠ 代价当场兑现：见 `R-110-2`」；注记另现量 **票 111／112／140 三张都不带 `-done`**（本腿 `ls` 坐实，`un=10`／`un=5`／`un=4`） | **〔真未闭合〕**＝两发执行类读数在当前 HEAD 缺；**且本格后果今天仍挂在开放池**，翻勾≠零后果 |
| 11 | 113 | `:54` | AC#4 | `- [ ] **AC#4** 变异：把新腿关掉 ⇒ AC#1 那条必须红；再把"祖先链只查一层"这种**半修**形状试一发 ⇒ 也要红（证明它咬的是全集不是某一行）。` | 注记 `113-posix-platformverifypplacement-has-no-link-leg-done.md:55`＝**己类·判不了**；表 `docs/evidence/s1/113-adversarial-acceptance.md:16` 判「**绿**」〔独立复现〕，⚠ 同一行末段逐字「**叶子这一维在 `SealFile` 方向没有交付用例**，见下"攻#4"」＝`R-113-C`（真表行 `113-…md:100`） | **〔真未闭合〕**＝四发变异当前 HEAD 复跑缺 **＋** `R-113-C` 未销（连"凭据齐"都谈不上） |
| 12 | 113 | `:56` | AC#5 | ``- [ ] **AC#5** 门禁：容器内 `-count=2 -v ./internal/winsec/ ./internal/memory/ ./internal/risk/` rc=0 且四数逐条点名（报 SKIP 要说是不是 `-v`；`` | 注记 `113-…-done.md:59`＝**己类·判不了**；表 `113-…md:17` 判「**绿（四数与三门），但 ban#8 那格的 sha↔数字配对错了一格**」，同行逐字「**d22：`sh scripts/d22scan.sh` rc=0，但 ban#8 `internal/=371`，不是它报的 368**」 | **〔真未闭合〕**＝**表对这枚格自己的判语里带一条没对上的账**＋五发读数在当前 HEAD 缺。这 14 枚里它最不该翻 |
| 13 | 115 | `:58` | AC#4 | ``- [ ] **AC#4** 变异：把你选的修法退回原状 ⇒ AC#3 新用例必须红；再试一发"只比大小写不敏感"（`EqualFold`）这种**半修** ⇒ 也要红`` | 注记 `115-seal-notices-carry-the-resolvers-answer-while-the-cases-compare-caller-spelling-done.md:60`＝**己类·判不了**，逐字「本票又**零枚裁决表可引**（`find docs/evidence -name '115*'`＝0）」；本腿 §3.2 复量＝票 115 落 **C 档**（`grep -iE '^115'` 零命中，exit=1） | **〔真未闭合〕**＝既无名下表、又无变异复跑 |
| 14 | 115 | `:64` | AC#6 | ``- [ ] **AC#6** 门禁：按包 `-count=2 -v` 四数逐条点名（报 SKIP 要说是不是 `-v` 量的；`-count=2` 不缓存）；`` | 注记 `115-…-done.md:67`＝**己类·判不了**：唯一那批数 `RUN=164 PASS=84 FAIL=0 SKIP=0` 署名**是实现方自己**（票面逐字标〔仅自述〕档），本票零枚裁决表 | **〔真未闭合〕**＝非实现者门禁读数完全缺位 |

**尺一汇总（逐枚判读）**：〔漏翻物理框〕＝**0 枚**；〔真未闭合〕＝**14 枚**。缺口形状分三族（5＋7＋2＝14）：**变异复跑类 5 枚**（#2 票 92 AC#5／#6 票 104 AC#3／#9 票 110 AC#3／#11 票 113 AC#4／#13 票 115 AC#4）、**执行门禁读数类 7 枚**（#3 票 92 AC#6／#5 票 97 AC#5／#7 票 104 AC#4／#8 票 105 AC#5／#10 票 110 AC#5／#12 票 113 AC#5／#14 票 115 AC#6）、**等人或等读权类 2 枚**（#1 票 07 Visual 等 owner 一次眼／#4 票 92 AC#7 等一枚有 `frontend/**` 读权的腿）。

⛔ **本腿在 93 枚结案票里没有找到任何一枚"注记主张已翻勾、物理框却还 `[ ]`"**（§0.2 row 10：把「`- [ ]` 那一行的下一行注记写着 `勾＝`／`翻勾：`／`⇒ 可翻勾`」当尺，`-done` 票内**零命中**；票内那批 `done-check-1` 勾／`done-fix-1` 勾 注记，**配的框都已经是 `[x]`**）。⇒ 昨天 265/268 那个洞**今天在结案池里已清干净**；上面 14 枚全部是"框如实反映缺口"，**翻它们＝把缺口判成已交付**。

### 4.2 尺二：同一行既 `- [ ]` 又含"翻勾"字样（全量扫，本腿现量＝**3 行／3 枚票，且三枚都不带 `-done`**）

尺＝`grep -rnE '^- \[ \]' --include='*.md' .scratch/wisp/issues/ | grep '翻勾'`（全池含开放票；限 `-done` 时命中 **0**）。

| 票 | 行 | `-done`？（本腿现量框态） | 行前 60 字 | 行里的"翻勾"是不是**自称已翻** | 判读 |
|---|---|---|---|---|---|
| **182** | `:26` | ⛔ 否（`un=7`／`[x]=0`） | `- [ ] **AC#3 逐堆指认归属票**（这是本票的主要产出）…也不许顺手替别的票翻勾。` | **否**＝禁令（"不许…翻勾"） | 非账面洞。**开放票、不入本件分母**，只在名册留痕 |
| **242** | `:28` | ⛔ 否（`un=3`／`[x]=0`） | `- [ ] **AC#3 本格不设产码要求**：…若判"够"，本格按"甲落地"翻勾；若判"要乙"，本格不许勾` | **否**＝**条件句**（判够才翻） | 同上：条件未判 ⇒ 框不翻是对的 |
| **245** | `:62` | ⛔ 否（`un=7`／`[x]=8`） | `- [ ] **AC#4 仍不翻**（它判"阻塞于票 228 如实、合格、不许翻勾"…）` | **否**＝逐字**"仍不翻"** | 同上：表已裁"不许翻"，框未翻＝**账框一致** |

⇒ 派单给的已知 3 枚（182/242/245）＝本腿全量复现的同一批，**无新增、无遗漏**；三枚全是"开放票＋禁令或条件句"，⛔ 零枚账面洞。

⚠ 台账 `A624`（`docs/reports/pending-and-issues.md:12235`，10-05 17:59）那句「同形 **182/242/245 三票四枚**＋`-done` 票 **15 枚**归 evidence-close-2 逐枚裁」与本腿现量**两处对不上**，具名登记（⛔ 不改台账一字）：

1. 「三票**四枚**」vs 本腿「三票**三枚**」：逐枚 `grep -cE '^[[:space:]]*- \[ \].*翻勾'`＝182:1／242:1／245:1（含缩进的宽松式同为 1／1／1）⇒ 那"第四枚"本腿量不出来，判为**台账侧计数口径差**（待编排者复核当时用的是哪把尺）。
2. 「`-done` 票 **15 枚**」vs 本腿 **14 枚**：`git grep` 复跑 `2a633eb8^` 那棵树＝**15**（含票 265 一枚），`2a633eb8`（10-05 17:59）就地补翻 265／268 各 `numstat 1/1` ⇒ 今天 **14＝15−1**（268 当时不带 `-done`，本来就不在那一分母里）。**两数不矛盾，差一枚 265、取数时点不同**（§0.2 row 8 坐实）。

### 4.3 一处反向形状（本腿顺手量到的，具名回报，⛔ 不动它）

尺＝「框已是 `- [x]`、紧随其后的注记却写 `⛔ 不翻勾`」＝**账框相反方向**的分歧，全池 `-done` 命中 **1 枚**：

- **票 110 AC#4（`.scratch/wisp/issues/110-no-ci-step-runs-internal-winsec-done.md:43`＝`- [x]`）**：
  - `:44` 注记＝``done-fix-1` 追加（戊类·**附条件的兑现句在码里、但无表侧复算**，⛔ **不翻勾**）…要翻请连那发 `go test -list` 的步级读数一起看。
  - `:45` 注记＝``done-check-1` **勾**＝本格判据是一个「**或**」句…三行定性…**⇒ 可翻勾**…⚠ 与 `done-fix-1` 的分歧：它把"那发步级读数"当本格的缺项，本腿认为那是 `R-110-4` 的账（本格是"或"句，第二支不需要读数）；**分歧交编排者裁，本勾若被撤请连带撤这行的结论**。
  ⇒ 同一枚格上**两枚腿注记互相反对、框却已翻成 `[x]`**（＝取了 `done-check-1` 那一边）。本腿⛔ 不撤、不裁，只把这一格指给编排者：**这枚是"分歧已被单方面销账"还是"编排者已裁过"，需要编排者自己认一次**。

---

## §5 本件不解决什么

1. ⛔ **不裁任何一格**。§3.2 两列只有"文件名"；§4 的判读只在"翻框凭据在不在盘上"这一件事上说话，不构成通过／退回。
2. ⛔ **不使任何一票被视为已验收**。`AGENTS.md` §0.3 要求"缺口审计与对抗验收必须由**另一个** agent 做"——本件恰恰在登记"哪几枚票至今没有那个另一个"。指路≠裁决（承前一腿 `265-267-evidence-index.md` §3 同句）。
3. ⛔ **不补任何凭据**。C 档那 17 枚要的是**新腿的新读数**（票 267 要 `267-v1`、票 115 要一枚名下的真表、票 268 要把 `probes/268/a2/verdict.md` 那份裁决**落到 `docs/evidence/s1/` 名下**），不是把这张索引写得更细。
4. ⛔ **不替 `A##`／`R##` 台账销账**。§4 点名的 `R-92b-4`／`R-97-1..3`／`R-110-2`／`R-113-C`／`R-92-6` 都还在原处，本腿只给 file:line，不动它们一个字。
5. ⛔ **不改尺、不并尺**。§1.3 选定 v3 并把口径写死，是为了下一位不必再漂；本件不给 v1/v2/v3 取平均、不"两种口径都列"——三套并存本身就是死腿的病灶。
6. ⚠ **本件自带的两处局限**（具名，防被当权威）：
   - 尺只读**第 1 行标题＋前 8 行免责句** ⇒ 抓不到前 8 行里的实现侧自称（票 84 名下那枚表 `:3` 的"写码代理"就是这么漏进 VERDICT 名册的，§3.3-2）。
   - 尺只认**文件名前缀**、只扫 `docs/evidence/s1/` 一层 ⇒ 把真存在的跨名凭据（§3.3-3 那 5 个形状）判成"无凭据"。**落点不合规**与**没做过**混在 C 档里，本件靠 §3.3 那张子表分开，⛔ 不靠改档。

**给编排者的三句结论（本腿现量）**

1. **凭据账**：A **74**／B **2**／C **17**，分母 **93**。A 里 **7 枚**的表出自编排者、**1 枚**（票 84）正文有实现侧自称 ⇒ 真·另一体终裁＝**66 枚**。
2. **翻框账**：尺一 **14 枚未勾框全部〔真未闭合〕、0 枚〔漏翻〕** ⇒ 今天没有"账面洞"要补翻；尺二 **3 行**全是开放票的禁令／条件句，0 枚账面洞；另有 §4.3 **1 枚反向格**（票 110 AC#4）待编排者自认。
3. **结构账**：票 **268** 的非实现者裁决件写在 `.scratch/wisp/probes/268/a2/verdict.md` 而**不在 `docs/evidence/s1/`**；票 **267**／**115** 的凭据挂在**别票名下**。这三处是**落点规矩**问题、不是"没做"问题，建议单独立一枚治理票处置（⛔ 本腿不立项）。
