# done-debt-1 —— 只读取证：两具 dangling 引用 / 票 230 残余名册复量 / 三枚格的凭据现读

- 锚点：起手 `git rev-parse --short HEAD` = **`795ed767`**（收尾时 HEAD 已漂到 `1d2ad737`＝`235-r1` 的注释落盘；件二那把尺在两个锚点各跑过一遍，见末节第 7 条）
- 时刻：`date` = **2026-09-29 22:31:48 +0800**（本文件最后一次自量 `wc -c` 见文末）
- 本腿性质：**只读**。零 `go test`／零 `go build`／零 `go vet`／零编译／零 commit／零 push；取数全部＝读文件＋`grep`／`awk`＋只读 `ls`／`find`。
- 未碰：`frontend/**`、`design/**`（连目录都没枚举）、`PLAN.md`、`thresholds.go`、golden、`allowlist.txt`、三枚冻结件。
  `docs/specs/SPEC-08-ui-ball-panel.md` **只读只引行号、零字节改动**（下面那一段是它自己的原文，不是我写的）。

---

## 件一　两具 dangling 引用

### 1.1　`docs/specs/SPEC-08-ui-ball-panel.md:42` → `docs/evidence/s1/68-*`

**现读原文（冻结件，一字未改）**：

```
41	    `44` 最可能的来源是把 `BallSizeSmallPx`（**最小可配置尺寸**，`tokens.go:369-371`）误读成体径。
42	    证据与推导：`docs/evidence/s1/68-*` + `internal/ball/tokens_test.go::TestSleepingSizeTruthTable`、
43	    `TestRecordedSleepingDiffBoxIsNotA44pxBody`。registry **A29**。
```

（`SPEC-08-ui-ball-panel.md` 现量 20,241 字节；更正块起于 `:32`「⚠ **更正（2026-09-20 23:11，编排者本人的录入错误…）**」。）

**该名册是否真的零枚——本腿现跑**：

| 命令 | 读数 |
|---|---|
| `ls docs/evidence/s1/ \| grep -c '^68'` | **0** |
| `ls docs/evidence/s1/ \| wc -l` | 319（名册里别的都在，只缺 68 族） |
| `find docs/evidence -iname '*68*' \| wc -l` | **0** |
| `ls docs/evidence/s1/ \| grep -cE '^(68\|115)'` | **0**（票 115 名下同样零枚，与票 230 现量那一行同读数、未变） |

⇒ **契约文字指向一份从没产出过的交付物，这个形状今天仍然成立**（不是已被人补掉了）。

**内容到底落在哪里（现读）**：票 `.scratch/wisp/issues/68-ball-default-visuals-parity.md`（13,238 字节，**文件名不带 `-done`＝仍开放**）的 Progress log `:97`–`:100` 就是那份推导本体：`:98`「did=AC#1 三列尺寸真相表（纯码 + 已录证据，**零新实测**）…」、`:99`「did=把 AC#1 的反证测试从"单次实测"升级成**三次实测**（`TestRecordedSleepingDiffBoxIsNotA44pxBody` 现在遍历 2098/2103/2120…）」。⇒ 34.72 的推导**在盘上、活的，只是没落在契约指的那个路径**。

**这笔账已有的三个登记处（全部现读，行号为本次实测）**：
- `docs/reports/pending-and-issues.md:948` ＝「待 owner 拍板（编号清单 R15，` :932`）」表的**第 11 项**：逐字「**SPEC-08:42 引用的证据目录 `docs/evidence/s1/68-*` 根本不存在**（A33③）。契约是冻结的，改它要 D22。」三选项＝补到该路径／改引用／保持不动，编排者推荐＝**改引用**。
- `docs/reports/pending-and-issues.md:9577`（`A452` 一节）逐字：「**两具 dangling 引用**…⇒ **归口票 230**…⛔ `docs/specs/**` 是冻结件，我不改它一个字，只登记。」
- `docs/evidence/s1/62-visual-spec-draft.md:122`（C-22 行）与 `:160` 也各自记了同一条，并明写「本文件**不改 SPEC**」。
- `docs/reports/HANDOVER.md:1330`「第 11 项＝`SPEC-08:42` 引用了不存在的 `docs/evidence/s1/68-*` 该怎么修」。

**该怎么销（一句）**：**维持不动、只登记——真正的销法要等 owner 在 `R15` 第 11 项上拍一次**（推荐列已写「改引用」，但改 `docs/specs/**` 一字都属 D22 人工批准，本腿不动）；**如果拍"补文件"**，成本极低：把票 68 Progress log `:98`–`:100` 的三列真相表原样落成 `docs/evidence/s1/68-*.md` 即可，内容已在盘上、无需重测。

### 1.2　`pending-human-review.md` 那一具（⚠ 先纠一处转述）

**⚠⚠ 先报派单自己的转述偏差**：任务说「`.scratch/wisp/issues/07-ball-state-machine-core-done.md:14` 指着 `docs/reports/pending-human-review.md`」——**现读不符**。票 07 的 `:14` 逐字是普通正文「DesignTokens, plus the `statemachine` module implementing the 20-state / 40-transition table」。那枚**活的指针在 `docs/evidence/s1/07-adversarial-acceptance.md:14`**；票 07 面上只有两处提到它，且都是**引文不是指针**：`:53`（`done-fix-1` 追加里逐字引了 `:14` 那句）、`:146`（`agent-bookkeeping-1` 的 audit-B 记录「noted, nothing edited」）。

**现读指针原文**（`docs/evidence/s1/07-adversarial-acceptance.md`，14,014 字节）：

```
14	**待人工项**：20 态视觉的主观签收（用户本人）→ 已登记 docs/reports/pending-human-review.md。
```

**该文件在不在——本腿现跑**：

| 命令 | 读数 |
|---|---|
| `ls -la docs/reports/pending-human-review.md` | `No such file or directory` |
| `ls docs/reports/ \| grep -c pending-human` | **0**（同目录 25 份文件里没有这枚；`ls docs/reports/ \| wc -l` = 25） |

**它当初登记的那件"待人项"今天实际记在哪里（具名行号，全部现读）**：

1. **真相源**：`docs/reports/pending-and-issues.md:5`「## 待人工审核（pending-human-review）」这一节，条目 **[H1] 在 `:7`**，逐字「**悬浮球 20 态视觉签收** — **2026-09-20 已改为实况签收**：编排者用 `build\balldebug.exe -stay` 把真悬浮球启到用户桌面…**满意则关闭；不满意提修改意见转新工单**。」（`:9` 还留着备选截图对照 `docs/evidence/s1/ball-states/*.png` 与 `design/screens/ball.html`）
2. **停车点**：`docs/reports/HANDOVER.md:1367`「## 8. 待人项（docs/reports/pending-and-issues.md 同步维护）」，**H1 在 `:1369`**。
3. **今天最活的那一份**：`docs/reports/desktop-signoff-2026-09-30.md`（10,733 字节／83 行）＝明天 10:30 的桌面签收对照表，**第 17 行那一件在 `:66`**「| 17 | **把所有状态挨个演一遍**（一共 20 种…）」——这正是 H1 那句"实况签收"的落地形态（owner 原话入档见台账 `A448`，`:9528` 前后一节）。
4. **缺陷本身也早已被记过两次**：`docs/reports/audit-B-docs.md:22` 与 `:50`；`docs/reports/HANDOVER.md:297`；`docs/reports/pending-and-issues.md:9508`。

**该怎么销（一句）**：**改指向**——把 `docs/evidence/s1/07-adversarial-acceptance.md:14` 那半句改成「已登记 `docs/reports/pending-and-issues.md` §待人工审核 **[H1]（`:7`）**」；那枚文件是裁决表不是冻结契约，改它一行属书记账动作，`audit-B` 早在 `:50` 具名要求过（「内容未丢失，故 MINOR」），**不需要补造一份 `pending-human-review.md`**；本腿只读，动手权在编排者／验收方。

---

## 件二　票 230 的残余名册复量（现跑锚定尺）

**尺**（严格按派单：带前缀锚定，`[[:space:]]` 不用 `[ \t]`）：
`grep -nE '^- \[[ x]\]|^- ⛔' <票面>`；未勾＝`grep -cE '^- \[ \]'`，已勾＝`grep -cE '^- \[[xX]\]'`，搬出勾框＝`grep -cE '^- ⛔'`。

**逐票现量（8 张，文件名全部现跑 `ls ${t}-*-done.md` 取到）**：

| 票 | 文件名（现量） | 未勾 | 已勾 | 搬出（`- ⛔`） | 与 `A449` 名册 |
|---|---|---|---|---|---|
| 07 | `07-ball-state-machine-core-done.md` | **1** | 2 | 3 | 1 ✓ |
| 104 | `104-sealfile-silently-drops-inherited-grants-done.md` | **2** | 3 | 0 | 2 ✓ |
| 105 | `105-c26-rewrite-account-has-no-production-reader-done.md` | **1** | 4 | 0 | 1 ✓ |
| 110 | `110-no-ci-step-runs-internal-winsec-done.md` | **2** | 3 | 0 | 2 ✓ |
| 113 | `113-posix-platformverifypplacement-has-no-link-leg-done.md` | **2** | 4 | 0 | 2 ✓ |
| 115 | `115-seal-notices-carry-the-resolvers-…-done.md` | **2** | 3 | 2 | 2 ✓ |
| 92 | `92-panel-composer-mode-attachments-workspace-done.md` | **3** | 3 | 1 | 3 ✓ |
| 97 | `97-dead-strict-param-and-the-comment-that-invents-a-caller-done.md` | **1** | 3 | 1 | 1 ✓ |
| 合计 | 8 张 | **14** | 25 | 7 | — |（八张相加：1+2+1+2+2+2+3+1 = **14**）

**全池自证（不只这 8 张）**：
- `for f in *-done.md; do n=$(grep -cE '^- \[ \]' "$f"); [ "$n" -gt 0 ] && echo …; done` ⇒ **命中且仅命中上面这 8 张**。
- `grep -h -- '^- \[ \]' *-done.md | wc -l` = **14**（与逐张相加一致）。
- `-done` 文件总数 `ls *-done.md | wc -l` = **79**。
- 搬出勾框那一族全池 `grep -hE '^- ⛔' *-done.md | wc -l` = **17**（其中这 8 张占 7 枚，其余 10 枚在别的已结案票名下）。
- **尺的盲区自查**：`grep -cE '^[[:space:]]+- \[ \]'`（缩进子框）＝**0**、`grep -hE '^- \[X\]'`（大写勾）＝**0** ⇒ 这把锚定尺今天没有漏计。

**结论**：**现量与 `A449`（台账 `:9537`）那行逐字相同＝14 枚／8 张，分布一处不差**。今晚的 `undone-28-1`／`done-fix-1`／票 221 结案**都没有改动这个分母**（票 221 现在 `221-…-done.md`，未勾 0／已勾 5）。

**接收方那一侧也顺手复量了（`A449` 同节写的"逐枚未变"，本腿现跑）**：`64` `un=1/chk=8` ✓、`12` `un=3/chk=6` ✓、`230` `un=5` ✓、`114` `un=9` ✓ —— **四条全对得上**。

**Status 行整批过期（现读为证，本腿不据其判完成度）**：`105` 的 Status 写着 `**rejected-needs-fix**`、`110` 与 `97` 写着 `ready-for-review`，而三张文件名都带 `-done`；`113` 面上有**两行 Status**（`:3` accepted-done、`:19` ready-for-review）；票 07 `:66` 那条 note 还写着「本票是 `-done` 但 **5 个框未勾**」，现量只有 1 枚。⇒ **完成度只认 `-done` 后缀＋现跑计数**这条再次成立。

---

## 件三　三枚开放格的凭据现读（抽法＋逐枚现读）

**抽法（可复算）**：候选＝14 枚未勾格里**声称了具体凭据**的 13 枚（`07:51`／`104:52`／`105:56`／`110:39`／`110:47`／`113:54`／`113:56`／`115:58`／`115:64`／`92:68`／`92:73`／`92:79`／`97:55`），跑
`printf '%s\n' …13 项… | shuf -n 3`（真随机源）⇒ **命中 `07:51`、`92:79`、`92:73`**（另跑了一次定种子的 `shuf --random-source=<(yes 42)`，落点是 `92:68 113:54 92:73`，**本文件只采真熵那一发的三枚**）。

### 3.1　票 07 `:51`「Visual: 20 states…human screenshot review vs design/screens/ball.html」

声称凭据＝`docs/evidence/s1/07-adversarial-acceptance.md:10` 与 `:14`。**现读**：

| 项 | 读数 |
|---|---|
| 文件在不在 | **在**，`wc -c` = **14,014 字节** |
| `:10` 逐字 | `| 3 | 视觉证据 | PASS（人工签收挂起） | docs/evidence/s1/ball-states/ 截图 + c21-native-tokens.md 对照表 |` ⇒ **与票面引文逐字相同** |
| `:14` 逐字 | `**待人工项**：20 态视觉的主观签收（用户本人）→ 已登记 docs/reports/pending-human-review.md。` ⇒ **逐字在盘**（也正是件一 1.2 那具死指针的出处，本腿从另一头复现了它） |
| 它引用的两样东西 | `docs/evidence/s1/ball-states/` **20 枚 `.png`**（`ls | grep -c '\.png$'` = 20，与"20 态"对得上）；`c21-native-tokens.md` **32,212 字节** |

⇒ **结论：凭据为真**（这一格本身仍**不可翻勾**，它缺的是 owner 一次眼，不是文件）。

### 3.2　票 92 `:73`「AC#6 台账与门禁…贴出逐作用域文件数」

声称凭据＝`92-adversarial-acceptance.md:147`「**结论：FAIL** —— 附我本机实测数字」、`92b-adversarial-acceptance.md:138`「AC#6 … **通过（数字全复算）**」、`92b-…:131` 记 `R-92-3` 已清。**现读**：

| 项 | 读数 |
|---|---|
| `92-adversarial-acceptance.md` | **在**，29,822 字节；`:147` 逐字「**结论：FAIL —— 附我本机实测数字。**（AC#6 原文要求 `gofumpt -l` 空，它不空；且 POSIX 读数不可复现。）」✓ |
| `92b-adversarial-acceptance.md` | **在**，44,704 字节；`:138` 逐字「| AC#6 台账与门禁 | **通过（数字全复算）** | 〔独立复现〕 | `gofmt`/`gofumpt` 空、`go vet ./internal/panel/ ./internal/tools/ ./internal/memory/` rc=0、`go test -count=2 -v ./internal/panel/ ./internal/tools/` = **360/238/0/0 rc=0**…」✓ |
| ⚠ **对不上的一处** | 票面写「`R-92-3` 在 `92b-…md:131` 记"已清…"」——**`:131` 逐字其实是 `R-92-1` 门二覆盖面 | **未清**…**」；`R-92-3`「**已清**…快照 / 工作树 / `gofumpt -l .` 三处我都跑，全空」真正在 **`:129`**（另 `:12`、`:35` 也各有一处）⇒ **判语本体在盘，行号指针漂了 2 行** |

⇒ **结论：主判语逐字命中、附带一枚行号指针漂移**（`92b:131` 应为 `:129`）。这一格仍**不可翻勾**的理由不是凭据假，而是票面自己写的"那批数只对 `91b5fc4`/`a8f9459` 两棵树负责、今天的树无效"。

### 3.3　票 92 `:79`「AC#7 负判据：`frontend/` 与 `internal/panel/` 里没有任何分支切换能力入口」

声称凭据＝`internal/panel/composer_test.go:269-270` 的正则、`internal/panel/composer_dispatch_test.go:648` 的 `banned` 名单、`92b-…md:139` 的表侧读数。**现读**：

| 项 | 读数 |
|---|---|
| `internal/panel/composer_test.go` | 在，**32,985 字节**；`:269` 逐字「`` \bgit\s+checkout\b|\bgit\s+switch\b|\bswitchBranch\b|\bcheckoutBranch\b|\bchangeRepo(?:sitory)?\b| ``」、`:270` 「`\brepoPicker\b|\bbranchSelect(or)?\b|\bworktree\b…`」，正则变量在 `:268 var gitSwitchCapabilityRe = regexp.MustCompile(` ⇒ **与票面引文逐字相同** |
| `internal/panel/composer_dispatch_test.go` | 在，**28,026 字节**；`:648` 逐字「`for _, banned := range []string{"checkoutBranch", "changeRepo", "repoPicker", "branchSelect", "vcs.switch"} {`」⇒ **逐字相同** |
| `92b-…md:139` | 逐字「| AC#7 不做 git 切换 | **通过** | 〔独立复现〕 | 包内用例 PASS + 我把 `fixtures/`、`dist/` 也 grep 了一遍 0 命中 |」✓ |
| 碰不到的那半 | `frontend/**` 半边本腿**零读**（派单写死），所以「跨 `frontend/` 的那半」本腿不替任何人判真判假 |

⇒ **结论：盘上可验的凭据全部为真；`frontend/` 那一半仍〔未现验〕**，这一格维持开放是对的。

### 3.4　批次定性

三枚抽到的格：**凭据文件全部存在、字节数全部非 0、主判语全部逐字命中**（`07:51`＝两处行号 `:10`/`:14` 加上两处引用物 `ball-states/` 20 枚 PNG、`c21-native-tokens.md`；`92:73`＝两处主判语 `:147`/`:138` 命中；`92:79`＝两处码 `:269-270`/`:648` 加表侧 `:139`，全命中）。唯一对不上的是 **`92:73` 里一枚附带行号指针**（`92b:131`→实为 `:129`）。
⇒ 本腿**不整批判〔凭转述〕**，但要写清边界：**若按"逐枚行号都必须命中"的严口径，`92:73` 那一格应记〔部分凭转述〕**；其余两枚是现读为真。这与台账 `A449`（`:9537` 前后）与 `A452` 已立的「行号引用不能整体信也不能整体废，**逐枚现读是唯一办法**」完全一致。

---

## 没销完／留给编排者（⛔ 不空）

1. **件一 1.2 的动手权不在本腿**：`docs/evidence/s1/07-adversarial-acceptance.md:14` 那一行改指向（→ `pending-and-issues.md` §待人工审核 **[H1] `:7`**）需要一次书记 commit；本腿只读没动。**另需一并纠正派单口径**：死指针在**裁决表 `:14`**，不在票 07 面的 `:14`——台账 `A452`（`:9577`）与 `HANDOVER:371` 那两处的写法照抄就会让下一位去票面找不到东西（本腿已实测票面 `:14` 是正文）。
2. **件一 1.1 卡在 owner 手上，不在票 230 手上**：`R15` 第 11 项（台账 `:948`）今天仍是**未答**状态，推荐列"改引用"是编排者写的、不是 owner 拍的；而 `docs/specs/**` 冻结⇒本腿一字不碰。**注意一个落空**：`A452` 说这两具引用"归口票 230"，但票 230 的**票面（28 行／6,503 字节）现跑 `grep -nE 'dangling|pending-human|68|SPEC-08'` ＝ 零命中**——它面上**一个字都没写这两具引用**。⇒ 要么给票 230 追加一节（append-only），要么就把这两具从"归口 230"改回"归 R15#11＋书记 commit"，**别留"账上有人归口了、票面却没有那一格"的形状**（这正是票 230 自己要治的病）。
3. **件二分母之外还有一格尺看不见**：全池 79 张 `-done` 里，`- [ ]` 只有 14 枚（已复算），但**表格式"不勾"另有 2 枚**——`153-…-done.md` 的 `:81`「| AC#2 痕带 taskID（票面拆 6 件事） | **[ ] 不勾** |…**两问无人裁**」与 `:82`「| AC#3 承重两句 | **[ ] 不勾** |…**句②…＝未裁**」，而同票 `:25`／`:35` 的 AC#2／AC#3 框**早已是 `[x]`**。⇒ 同一张票里"框已勾／表内自报未裁"并存；要不要把它们挪进分母，是编排者的裁，不是本腿能定的。
4. **票 115 名下裁决表仍 0 枚**（本腿现跑 `find docs/evidence -name '115*'` = 0、`ls | grep -cE '^115'` = 0）：票 230 AC#1 那张七项检查表**一张都没出**，账不在本腿。
5. **票 68 仍开放**（`68-ball-default-visuals-parity.md`，无 `-done`，13,238 字节，`:97`–`:100` 是唯一活着的那份推导）⇒ `A452` 状态行说"另需补派一枚只读腿做票 68 的台件补档"，本腿**没领那枚活**，只做了引用取证；补档要不要顺带解掉 `R15#11`，等 owner。
6. **一句免责**：本腿所有"判语在盘"的结论只证明**文字与文件存在且逐字对得上**，不证明那些读数今天在树上仍可复算（那要跑测试／突变＝`235-r1` 的地界）。
7. **锚点在本腿运行期间漂了一发（如实记）**：起手 `git rev-parse --short HEAD` = **`795ed767`**，收尾时同一把尺 = **`1d2ad737`**（`235-r1` 的 AC#1 注释腿落盘，只动 `internal/tools` 测试文件里那段注释）。⇒ 本腿**漂锚后重跑过件二那把尺**：`grep -h -- '^- \[ \]' *-done.md | wc -l` 仍是 **14**，八张逐票未勾数**一枚未变**（`07(1) 104(2) 105(1) 110(2) 113(2) 115(2) 92(3) 97(1)`）⇒ 件二的结论对这两个锚点都成立。件一／件三所引的 `docs/**` 与票面文件在 `git status --porcelain -- docs/specs docs/reports .scratch/wisp/issues` 下**零改动**（只有本腿这一枚 `?? .scratch/wisp/probes/done-debt/`），故行号不因这一发漂。⚠ **另记一枚不是本腿的脏改动**：`git status --porcelain -- docs` 现量 ` M docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md`（别人名下的验收表，本腿**一枚字节未动、未引用**）⇒ 编排者推送时别把它算进本腿的账。

---

**本文件自量**：`wc -c .scratch/wisp/probes/done-debt/debt.md` = **19,034 字节**（非 0；交付判据取字节数，不取工具回执。真值以编排者现跑为准，本腿只保证「非 0」）。本目录只建不删、本腿零 commit 零 push。
