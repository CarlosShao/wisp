# 审计：编排者 2026-11-14 结的 5 张票（441／442／274／443／444）——只信盘上读数

- **审计者**：`agent-26`（非实现者；这 5 张票的实现/收件都不是本腿）
- **派单存档**：`.scratch/wisp/probes/orch/20261114/txt/audit-closures-20261114-dispatch.txt`
- **起点 HEAD**：`3e9d9779f`（2026-11-14 17:33:55 +08）；**落笔时 HEAD** `81ba68874`、`origin/dev` `5ceb5af6d`（中间只多 2 笔 probe/台账 commit，**5 张票的 blob 一字未变**，见 §5）
- **判据**：只认盘上读数——`git cat-file`／`git ls-tree`／`git merge-base`／文件字节与行号／`gh run view --json jobs` 的逐 job 逐 step 结论；票面自述与完成通知一律当**未验证**
- **写面**：本件 ＋ `.scratch/wisp/probes/audit-20261114/**`。⛔ 未改票面、未改产码、未动阈值/golden/契约、未 push

---

## 0. 一句话结论

**应当复开 2 张：441、274。** 442／443 保持已结（442 补两处行号更正、443 无待补）；444 本来就未结，**且它现在是五张里唯一没有把本地绿当 CI 绿的一张**——它的 AC#5 盘上已有 CI 读数（见 §3.5），可以结，但那是编排者的动作，本件不改票面。

---

## 1. 总表（一格一行，逐条可复跑）

| 票 | 票面状态 | AC 框（翻／未翻） | 它引的 commit 真在？ | 它引的证据件（盘上字节） | file:line 抽核 | 有没有把「本地绿」说成「CI 绿」 | 它引的 CI 说法有 run 撑？ | HEAD 同票号文件数 | 判 |
|---|---|---|---|---|---|---|---|---|---|
| **441** | 已结 | 5／0 | `8e1367bfd`／`b412f50e4`／`952f4462b` 均 `commit` ✓ | `docs/evidence/s1/441-instrument-abstention.md` **19324 B** ✓ | ABSTAIN 实体在 `posix_seal_window_120_other_test.go:328`；票面写的 `:105`／`:296` **都只是注释**；`WISP_120_ATTEMPTS:142` ✓ | 不是冒充，是**引错对象**：它给的两发不含修复 | **无**（见 §2.2，字面为假） | 1 ✓（历史 6 笔曾 =2） | **复开** |
| **442** | 已结 | 5／0 | `6949484e1` = `commit` ✓（18 文件 `+2369/−146`，我现量） | `442-runtime-criterion.md` **14139 B** ✓；测试件 **11536 B** ✓；seam 件 **1851 B** ✓ | `:199`／`:216`／`:248`／`:254` ✓；**`:224`／`:235` 不在它说的地方**（真守卫在 `:231`／`:240`，见 §3.2）；seam 在 `panel_host_windows.go:332-337` ✓ | 票面**没有**任何 CI 主张；门禁全部写成"我自己跑了"（本地） | 票面没主张；**盘上反而有**（`43575142588`／`43575466311` 两发 6/6 全绿，见 §3.2） | 1 ✓（历史无重复） | 保持已结（补两处行号） |
| **274** | 已结 | 12／0 | `45c68fa79` = `commit` ✓（4 文件 `+28/−24`） | `274-c1-three-cells-verdict.md` **25483 B**／`274-pa-ci-nail.md` **28239 B**／`274-w1-provenance-nail-census.md` **8779 B** ✓ | `build.ps1` 的 `:47/:49/:158/:165/:209/:212/:215/:240/:270` ✓；`SPEC-11:44/:45/:69` ✓；`SPEC-01:70` ✓ | **命中**：结案段写"path-length VERDICT GREEN"，而**结票动作自己把 CI 的 lint 弄红了 2 发、27 分钟**（见 §3.3） | 门禁那句是本地读数；CI 反证在盘上（`43575830504`／`43575854808`） | 1 ✓（历史 1 笔曾 =2） | **复开** |
| **443** | 已结 | 5／0 | `d78a0ae5c`／`310b38cf4` = `commit` ✓ | `443-instrument-abstention.md` **19537 B**（＝票面自报，逐字节同数）✓ | ABSTAIN 在 `command_52_test.go:811` ✓、注释 `:738` ✓、真红 `RED(443:…tree-survived)` 在 `:858` ✓ | 票面**只写本地**，一个"CI 绿"的字都没有；我自己复跑 `ok 18.701s` | 票面没主张；**事后盘上够**：`43577965575` 的 plugin 步 success（含修复）、`43578030753` 6/6 全绿 | 1 ✓（历史无重复） | 保持已结 |
| **444** | **未结** | 4／**1**（AC#5） | `d857ae57d` = `commit` ✓ | `444-instrument-abstention.md` **23277 B** ✓ | 7 条 token 逐条对：`:138`／`:76`／`:172`／`:128`／`:170`／`:192`／`:229` ✓ | **没有**——票面 AC#5 原话就是"⛔ 不拿本地绿冒充 CI 绿" | 无主张（AC#5 留白等 CI） | 1 ✓（历史无重复） | 保持未结 |

---

## 2. 三件专项

### 2.1 HEAD 双计数

**现 HEAD 的读数（判据＝同一 commit 的 tree 里同票号文件数必须 =1）**

```
git ls-tree -r --name-only HEAD -- .scratch/wisp/issues | grep -E '/441-'
```
441=1、442=1、274=1、443=1、444=1；磁盘上同样每票 1 枚。⇒ **当前树不欠账。**

**历史窗口（我逐 commit 扫 `b062d52a7^..HEAD`，^ 之前 441 还不存在）**

| 票 | 双计数窗口 | 持续 | 制造者 | 修 |
|---|---|---|---|---|
| 441 | `e4f5baa56` 15:42:45 → `936bbaa5b` 16:00:32（**6 笔已推送的 commit**） | **17 分 47 秒**（到修复 commit `ae4e506ee` 16:01:48 是 19 分 03 秒） | 结案 commit 只 `A` 了新 `-done` 路径，没 `D` 旧路径 | `ae4e506ee` 16:01:48 显式删旧路径 |
| 274 | `813bf8574` 16:26:12（1 笔） | **20 秒** | 同上 | `c0e48dbc4` 16:26:32 显式删旧路径 |

442／443／444 **从未双计数**（442 的 `-done` 是 `R097` 同一笔里两路径都写了，这是正确形状）。

⇒ 编排者今天说的"441／274 都犯过"属实，而且时间线可核；**现在已修好**，但 441 的账在**已推 origin** 的 6 笔 commit 里挂了 **17 分 47 秒**——谁在那段时间里数票池都会多数一枚。这属记录层缺陷，不是当前树缺陷，不单独构成复开理由，记在这里是因为它随时会再来一次。

### 2.2 overclaim：五张里只命中 441 一句

**被命中的句子**（441 票面 `### CI 两发连续绿（本票修复的最终验证）` 整段；行号以 `:56` 起算，但该票有裸 CR，**行号不可靠**，按标题找；结案 commit `e4f5baa56` 里就已经这么写，我核过那一版）：

> ### CI 两发连续绿（本票修复的最终验证）
> - `eb9cf0e63`：**success**
> - `a46352fd3`（id `43568643680`）：**6/6 job 全绿** …
> ⇒ 原来红的 `test-core` 在**含本修复的两发**里都绿。

**盘上读数（决定性）**

```
git log -1 --format='%p' 8e1367bfd   ->  eb9cf0e63
```

441 的修复主体 `8e1367bfd` 的**父提交就是 `eb9cf0e63`**；`a46352fd3` 更早（14:55:16，修复是 15:23:38）。⇒ **这两发都在修复之前，按定义不可能"含本修复"**。run 本身是真的、绿的（`43570834384` 6/6 success、`43568643680` 6/6 success），但**它们撑的是修复前那棵树**。

**它想说的那两发，盘上也真有**（结论还站得住，只是引错了证）：

| 含修复的 head | run | 结论 | 逐 job |
|---|---|---|---|
| `952f4462b`（修复三笔的 head） | `43571235710`（15:28 起） | **success** | lint／test-core／test-windows／slo-smoke／slo-full／lint-frontend 全 success |
| `2e8233aeb`（它的下一笔） | `43571511743`（15:31 起） | **success** | 同上 6/6 |

**更硬的一条（完成时刻，不是创建时刻）**：结案时（15:42:45），那一刻**已跑完**的两发是 `eb9cf0e63`（15:41:19 done，success）与 `650eff6cf`（15:42:09 done，success）——**两发都在修复之前**（链：`a46352fd3 → 650eff6cf → eb9cf0e63 → 8e1367bfd`）。第一发**含修复且已跑完**的 run 是 `952f4462b`，**15:47:00 才完成**，比结案晚 4 分 15 秒；`2e8233aeb` 15:47:39 完成；随后的 `b30cc7cf2` 15:57:33 完成且**红**（test-windows）。⇒ **结案那一刻，这枚修复在 CI 上"零完成读数"**；票面把修复前那两发当成了"本票修复的最终验证"。这不是"两发绿不真"，是"绿的两发不可能是它"。15:47 那两发后来确实把结论补上了（见上一段），但**结案时并没有**。

**其余四张**：442／443／274／444 票面上**没有**"全绿／已绿／CI 绿／两发连续绿"这类整句。274 有 run 级数字（AC#0/AC#2，分项），444 反向写明"不拿本地绿冒充"。⇒ **这一类错今天只剩 441 一处。**

### 2.3 「本地绿冒充 CI 绿」：逐票标

| 票 | 结案凭据的成色 | 标注 |
|---|---|---|
| 441 | CI 读数有，但对象错；本机 POSIX 跑不了（票面自己写了"这一格是我的核查缺口"） | **CI（引证错）**，见 §2.2 |
| 442 | `wisp-cli-staged-tests.sh` RUN 48／PASS 48、vet／gofumpt／staticcheck／d22scan／path-length —— **全是本地**，票面没有一个"CI"字 | **仅本地**（但事后 CI 有绿，见 §3.2） |
| 274 | 结案段"四门禁 …——腿跑的"＝本地读数，且**被结票动作自己推翻** | **仅本地 + 读数停在前动作** |
| 443 | 三条定向测试串行全绿、vet／gofumpt／d22scan／path-length —— 全本地；票面明写"本票结"，没声称 CI | **仅本地（如实）** |
| 444 | 明确拒绝："⛔ 不拿本地绿冒充 CI 绿"，AC#5 留空等 CI | **未主张** |

⇒ 真正该改的是 441 的引证与 274 的时点；442／443 只是"只有本地读数"，而这一点它们都没有粉饰。

---

## 3. 逐票细节（总表之外、判据所依赖的读数）

### 3.1 票 441 —— 修复是真的，票面与引证是坏的

**（a）修的在不在**：`8e1367bfd` 非 probe/docs 文件只有 **1 枚** `internal/winsec/posix_seal_window_120_other_test.go` ⇒ 产码 0 字节 ✓；`git show … | grep '^+.*t\.Skip'` = **0 行** ✓（没偷用 Skip）；`internal/winsec/winsec_other.go` 现读 SHA256 `19F3C211055011CC699A8D8027EF49FF543EDE0DA5C863F58AAAAD21A640D3E2`，与票面写的前缀 `19F3C211…` 同串 ✓。

**（b）为什么它现在不再假红**：`judgeRun120` 把"窗口到达过吗"与"密封还封得住吗"拆成两个载体——`refused>0 && sealedOwn==0` 走 `ABSTAIN(441:<shape>:raced-quiet)`（带 attempts/samples/refused 数字，`:328`）；控制腿跑不动仍真红（`:323`）；`attempts==0` 仍真红（`:311`）；`refused==0` **仍真红**（`:315`，这是刻意的：删掉 racer 也会让 refused=0，"零命中"不能变绿）。⇒ CI 那条红（`refused=3..4, sealed-own=0`）正好落进弃权支，设计对得上。

**（c）census 那个数是真的**：票面 AC#4 写 "37 行 / 13 文件"。我在 `8e1367bfd` 上复算：

```
git grep -n -E 'RED instrument|instrument is dead|RED INSTRUMENT' 8e1367bfd -- internal
-> 37 行（全是 *_test.go）、13 个文件
```

与票面**逐数对上**。当前 HEAD 同一把尺是 **43 行／13 文件**（444 那波往 winsec 7 枚里加了行），⇒ 37 是当时真数、不是抄的。

**（d）票面字节缺陷（盘上缺点，不是判断）**：这张票在 HEAD 上含 **3 枚 NUL**（字节偏移 213／726／854）＋ **3 枚裸 CR**（`\r` 转义留下的，导致不同读法数出来的行号不一样）：

- `\05fb21317\` → `<NUL>5fb21317`（**两处**：出处行与"触发这一发的 commit"行）、`\069f97cac\` → `<NUL>69f97cac`（一处）——CI head 的首位 `0` 被吃；
- `\bash -n\` → `\ash -n\`（**正好落在 AC#5「门禁四数」那一行**）、`\test-core\` → `<TAB>est-core`（落在"同一套 test-core 在三发都绿"那一行）；

后果三条，都是现读：
1. `git show --stat 8e1367bfd` 对这张票显示 **`Bin 3031 -> 4211 bytes`** —— git 当它是二进制，**逐行 diff 没了**，后续 commit（含 `952f4462b` 的"改写 Progress log"）都只能整块替换；
2. 任何工具直接读会拒（本审计第一次 `read` 就被判 binary）；
3. 同一份文本按 CRLF 切是 68 行、按裸 CR 也切就会多出断行 ⇒ **票面自身不可按行号引用**（本审计只能用标题／上下文定位），而票内证据又全按行号引用别处。

**同样的病今天在 445 上被发现并修了**（`4184c4e16` 的 commit message 逐字讲了 backtick 转义把出口码变成 NUL），**441 没修**。⇒ 这一格不是"判断不同"，是盘上缺陷。

**（e）判**：修复不需重做（CI 两发绿在盘上）；但票面①唯一那句"最终验证"引的是不含修复的两发，②正文在 git 里是二进制、命令与 head 被吃掉首字母。**复开，要补的那一格＝引证 + 字节**（见 §4）。

### 3.2 票 442 —— 证据齐、CI 反而比它自己说的强；两处行号指错

**（a）字节**：`442-runtime-criterion.md` 14139 B、`cmd/wisp/panel_data_path_distinct_442_windows_test.go` 11536 B、`cmd/wisp/panel_data_path_seam_windows.go` 1851 B —— 与票面自报逐枚同数 ✓。

**（b）file:line**：`SAME data path` 判据本体在 **`:254`** ✓；`:199`（记录数＝点数×轮数）✓、`:216`（点数）✓、`:248`（空值）✓。**两处错**：

| 票面说 | 盘上那一行实际是 | 真守卫在哪 |
|---|---|---|
| `:224`（期望文件恰一个点） | `for _, point := range points {` | **`:231`** `if len(found) != 1 {` |
| `:235`（每点每轮都构造） | `}` | **`:240`** `if len(values) != rounds {` |

差 7／5 行，属"指到脚手架没指到守卫"。该文件自 `6949484e1` 起**一字未改**（`git log --oneline -- <file>` 只有那一笔）⇒ 不是行号腐烂，是当时就写错。

**（c）产码面**：`6949484e1` 有 **2 枚非测试产码件**——`panel_host_windows.go` **`+6/−0`**（3 行注释＋`if panelDataPathSeam != nil { … }` 3 行，现在在 `:332-337`）＋新件 `panel_data_path_seam_windows.go`。票面 §边界 写的是"⛔ 产码"，松这条的是票内 15:5x 的「甲＋落①」裁定（`:61-65`）。两段现在都在盘上，读票的人必须一起读；**建议补一句"边界以 15:5x 裁定为准"**，不算错。

**（d）AC#2 的字面偏离已具名**：突变做法是临时改仓内 `panel_unbound_253r6_windows_test.go:71` 再还原，**不是**票面写的"仓外导出树"；收件人自己写在 `:90`，并把 `C9A14528…40551` 同串、`git status` 0 行、复绿都记了。实质（能编译的探针／真红并点名／还原后 SHA 同串）成立 ⇒ **不构成 overclaim**（因为披露了）。

**（e）CI**：票面没主张。盘上 `6949484e1`（`43575142588`）与 `021f15661`（`43575466311`）两发 **6/6 全绿**，`cmd/wisp CLI tests` 步 success ⇒ 新判据在 CI 上真跑过（该步 RUN 846／PASS 603）。但**结案那一笔自己的 run 是红的**：`312a58612`／`43575529807` = failure，红在 `cmd/wisp CLI tests` 的 `TestWireR1AnnounceEndToEnd`（与 442 无关，后续发次复绿）。如实记：**结案时点的 CI 是红的**，只是红因不在本票。

**（f）`path-length` 我现跑＝GREEN**（`over the hat=69`／`roster entries=69`／`not in roster=0`，rc=0）⇒ 442 的 AC#5 那一条现在仍成立。

**（g）判**：保持已结。补两处行号（`:224→:231`、`:235→:240`）＋边界一句。

### 3.3 票 274 —— 结票动作自己把门禁弄红了 27 分钟，票面写的是绿

**（a）批准链在盘上（不是无授权改契约文字）**

- `docs/reports/pending-and-issues.md:37119`：`### Q-440 裁定（2026-11-14 15:5x 编排者记）：机主批「都按推荐」⇒ 走 P1`，撤销口令同段；
- 6 处改文与 `docs/evidence/s1/274-c1-three-cells-verdict.md §3.2` 的提案逐条一致；我抽核四处现文：`SPEC-11:44`（"…npm ci＋npm run build…没有 node/npm 或构建失败＝整条构建失败（无跳过路径）"）、`SPEC-11:45`（"embed **frontend/dist** 已就位，见 frontend/embed.go"）、`SPEC-11:69`（`### 3.3 compose.dev.yml…`，§3.2 整节已删）、`SPEC-01:70`（一行 `frontend/embed.go  # //go:embed all:dist（…不起 localhost HTTP 服务；D29）`）——**与收件段引的现文逐字一致**。
- ⇒ 规格/`.gitignore` 是真改了的：`45c68fa79` = 4 文件 `+28/−24`（规格 3 枚 `+4/−24`，第 4 枚是本票 Progress log `+24`）。

**（b）但 AC#5 的框与自己的 diff 相反**

票面 AC#5 现文（`:52`，逐字）：**AC#5 零新契约面**：⛔ 不新造 `C##`，⛔ 不往 `C17` 方法白名单加名字，**⛔ 不动 D29/D41/SPEC-11 的文字**（`SPEC-11 §2.2` 那句"embed lands in S5"要不要改＝人工批准，**本票只上报不修**）——而本票 `45c68fa79` 改了 **SPEC-11 `+3/−20`、SPEC-01 `+1/−3`、`.gitignore −1`**。可以辩"AC#8 被 Q-440 授权后 AC#5 那句作废"，但**票面没有一个字这么写**（收件段只写"AC#8 翻，本票 12/12"）。⇒ 这是"框是勾的、框里的句子与自己的 diff 相反"。

**（c）门禁那句是错的（这是复开的主因）**

票面 `:310` 收件段逐字：**"四门禁：`d22scan` clean、`path-length` **VERDICT GREEN**（over 69 / roster 69 / not-in-roster 0）——腿跑的"**。

盘上事实链（单变量）：

| 时刻 | head | run | lint 步 `Tracked path-length budget (ticket 262)` |
|---|---|---|---|
| 16:22:28 | `312a58612` | `43575529807` | success（结票前，GREEN） |
| **16:26:12** | **`813bf8574`（274 首次带 `-done`）** | **`43575830504`** | **failure** |
| **16:26:32** | `c0e48dbc4`（删旧路径） | `43575854808` | **failure**（新名本身 102 字符） |
| 16:53:20 | `6236a6110`（改名 102→69） | 下一发 `43577965575` | success（该发只红 test-windows） |

`43575830504` 的 lint 日志逐字：

```
check-path-length-budget.sh: denominator: tracked paths=20028  over-budget=70  covered by roster=69  not in roster=1
check-path-length-budget.sh: RED - over budget and NOT in the roster: .scratch/wisp/issues/274-no-nail-requires-shipped-exe-to-carry-page-build-ps-step-2-comment-is-false-in-both-halves-done.md  (relative 123 chars, name 102 chars, …)
check-path-length-budget.sh: VERDICT RED
```

⇒ **红因就是本票结案时加的那 5 个字符（`-done`）：名字 102 字符 > 规矩 9 的 `HAT_NAME_LIMIT=100`**。红窗口 **16:26:12 → 16:53:20（27 分 08 秒）、2 发 CI**，由 `6236a6110` 改名修掉。274 自己的收件段**没有**这段红；它只出现在 443／444 的票面里（两张都如实写了"红不是我的，是 274 的 `-done` 名"）。

**（d）一条对它有利的更正**：票面 AC#0(b) 留着"口径差一行"（"CI 从不跑 `panel-assets`"，理由是 `ci.yml:1587` 只是注释）。盘上现在不是这样了：`274-paci` 加的那一步 **`Shipped wisp.exe carries the page bundle (panel-assets -check on the artifact)`** 在 `43575142588`／`43575466311`／`43578030753`／`43578228598` **四发里都是 success** ⇒ 这一格现在有 CI 读数，票面那句是**保守未更新**，不是虚报。

**（e）判**：保持"12 格都做了"的实质没问题；问题是**结票后的门禁状态与票面写反**，以及 AC#5／AC#8 的授权关系没落到票面。复开，要补的那一格见 §4。

### 3.4 票 443 —— 最干净的一张；票面只写本地，事后 CI 也够

**（a）立票事实可核**：`e1cc304e0`（`43567331249`）failure、`1f0d04714`（`43568581591`）failure、`a46352fd3`（`43568643680`）success ⇒ **3 发红 2** ✓；`1f0d04714` 的 `--name-status` = 3 枚 probes 件＋1 枚 `314-shape-ruler.md`，**没有产码** ⇒ "只改了 probes 和证据件" ✓。

**（b）字节与行号**：`443-instrument-abstention.md` **19537 B**＝票面自报 ✓；`command_52_test.go:811` 是 `t.Logf("AC#2 ABSTAIN(443:timeout-tree:raced-no-pid) …` ✓；`:738` 是注释 ✓；真红 `RED(443:timeout-tree:tree-survived)` 在 `:858` ✓。

**（c）我自己复跑（本机 Windows，HEAD 树）**：

```
go test ./internal/plugin/manifestslot -run 'Test52|Test138' -count=1
ok  github.com/CarlosShao/wisp/internal/plugin/manifestslot  18.701s
```

⇒ 票面"尺3 自己复跑"这一格**可独立复现**（0 FAIL）。

**（d）日志里的数也对得上**：`probes/443/` 里 `03-natural-25` 3 条 abstain、`13-natural-25-postedit` 0 条 ⇒ 50 发 3 弃权=6% ✓；`11-amplified-5-delay1500` 5 条 ⇒ 5/5 ✓；`04-amplified-20-delay3000` 20 条 ⇒ 20/20 ✓；`05-mut-singlekill-raced` 1 abstain＋1 RED＋FAIL ⇒ 牙 ✓。

**（e）CI**：票面**没主张**。盘上事后有两发含修复的读数——`43577965575`（`45ccaf186`，含 `d78a0ae5c`）的 **`Portable windows tests (proc/secret/config/risk/ball/perm/plugin/llmrecord)` = success**（该发整体红在 445 的 `Test390r1`），`43578030753`（`721a1d8a0`＝**443 自己的结案 commit**）**6/6 全绿**。⇒ 443 的修复在 CI 上是绿的；**票面当时只能算"仅本地"，但它如实这么写了，没有冒充。**

### 3.5 票 444 —— 正确保持未结；AC#5 现在盘上已经够了

**（a）框**：4 翻／1 未翻，**唯一没翻的是 AC#5**（`:29`），它的字面就是"CI 上 `test-core` 与 `test-windows` 各两发连续绿"，并且同段写着"⛔ 不拿本地绿冒充 CI 绿"。

**（b）token 与"牙"**：7 条 `ABSTAIN(444:…)` 逐条在盘（`openleaf_322_other:229`／`ancestoropen_342_unix:138`／`ancestoropen_342_windows:76`／`ancestorunlink_348_unix:172`／`ancestorunlink_348_windows:128`／`deleteseam_348r2_windows:170`、`:192`）✓。断言数（`t.Errorf|t.Fatalf` 计数，改前→改后）：`3→3`、`2→3`、`7→8`、`15→16`、`14→23`、`12→16`、`4→7`、`8→13` —— **没有一枚下降**；`ancestoropen_342_unix` 里被删的 `return` 我也抽读了，确实让下面的 damage 断言继续跑（不是跳过）。

**（c）两处精度问题（不影响 AC）**：
- 票面写"基线／突变日志 **22 枚**入库"；`probes/444/**` 在 HEAD 是 **22 枚文件**，其中 `logs/` 下是 **20 枚**（另 2 枚是 `commit-msg.txt`／`repro-444.sh`）⇒ 类别差一枚，属"读数要能照着复现"的精度。
- 票面写"path-length 那条红是我自己造的（274 结票文件名 102 字符），我已改短"——时间线成立（见 §3.3c），但**它是在 444 的票里修 274 的账**，这也是 274 该复开的旁证。

**（d）AC#5 的 CI 读数（我现取）**：`d857ae57d` 的下一笔 `45ccaf186` 是红的（445 的 `Test390r1`）；此后 **`43578030753`（`721a1d8a0`）／`43578127657`（`5b56792c6`）／`43578228598`（`c64a4fa22`）三发 6/6 全绿**，三发都含修复链，完成时刻 17:14:21／17:14:15／17:17:55。⇒ **当时"等 CI"是对的（三发在 17:1x 之后才完成）；现在盘上已经满足 AC#5**。这属于"可以结但还没结"，**不是复开对象**，本件不动票面。

---

## 4. 结论（逐张）

**保持已结：2 张**

1. **442** —— 证据件三枚字节全对；CI 两发 6/6 全绿且跑过它那条新判据；`path-length` 现跑 GREEN。**只补两处行号**（`:224→:231`、`:235→:240`）＋在边界段加一句"以 15:5x 裁定为准"。
2. **443** —— 立票事实、证据件字节、token 行号、日志数字全部对上；我独立复跑 `ok 18.701s`；票面只写本地、没冒充 CI；事后 CI 也有绿。**无待补。**

**保持未结：1 张**

3. **444** —— AC#5 的留白是对的（写下"等 CI"时三发绿 run 还没跑完）；现在三发 6/6 已绿、含修复，**够勾**。结不结由编排者拍，本件只报读数。

**应当复开：2 张**

4. **441 —— 复开（引证 + 字节）**。理由：①票面唯一那句"CI 两发连续绿（本票修复的最终验证）"点的 `eb9cf0e63`／`a46352fd3` **按 ancestry 证明都不含修复**（`eb9cf0e63` 就是修复的父提交），而结案时（15:42:45）已跑完的两发 `eb9cf0e63`／`650eff6cf` 也都在修复之前，**第一发含修复且跑完的 run 要到 15:47:00**——即结案那一刻这枚修复在 CI 上零完成读数；②票面本身在 git 里是**二进制**（`Bin 3031 -> 4211 bytes`），含 3 枚 NUL（吃掉两个 CI head 的首位 `0`）＋ 3 枚裸 CR ＋ 两处被转义吃字母（`\bash -n\`→`\ash -n\` 在 AC#5 行、`\test-core\`→`<TAB>est-core`）。**要补的那一格**：把该段引证改成 `952f4462b`（`43571235710`）／`2e8233aeb`（`43571511743`）这一对，并把两发 head 与"含修复"的 ancestry 关系写清；另起一笔把 3 枚 NUL／3 枚裸 CR／两处被吃字母修回（⚠ 按 `2934bd2d5` 的教训：只改那几处字节，⛔ 不许整段 rewrite 把 Progress log 吞掉）。**修复本身不需要重做。**

5. **274 —— 复开（门禁时点 + 授权关系）**。理由：①结案段写"`path-length` VERDICT GREEN"，而**结票动作（给名字加 `-done`，102 字符 > 帽 100）自己让 CI 的 lint 红了 2 发、27 分钟**（`43575830504`／`43575854808` 的 `Tracked path-length budget` = failure，`6236a6110` 改名后才绿）——赔的是别人的腿和别人的红；②AC#5 的框里写着"⛔ 不动 SPEC-11 的文字"，而本票自己的 `45c68fa79` 改了 SPEC-11 `+3/−20`（授权链＝Q-440「都按推荐」，盘上有，但票面没把这条与 AC#5 的关系写死）。**要补的那一格**：AC#5 改写成"结票后读数"，并把 16:26–16:53 的两发 lint RED（含 `6236a6110` 的修）记为已发生事实；同段加一句"P1 六处改文由 Q-440 授权、AC#5 的'不修 spec'自本条裁定起只约束产码面"。**6 处规格改文不回退**（Q-440 已批，撤回口令在台账）。

> 若只允许复开一张：**441**——它是五张里唯一一句"CI 绿"是假的；274 的性质是"读数停在动作之前"，代码与规格本身没走样。

---

## 5. 复跑命令（逐条可粘；Windows 上 `sh` 要走 Git bash：`D:\work\soft\Git\bin\bash.exe`）

```sh
# 0) 树锚
git rev-parse HEAD                                  # 落笔时 81ba68874；票面未变
git diff --name-only 3e9d9779f HEAD                  # 只有 probes/orch + 台账两笔

# 1) HEAD 双计数（每票必须 =1）
git ls-tree -r --name-only HEAD -- .scratch/wisp/issues | grep -E '/(441|442|274|443|444)-'

# 2) 历史双计数窗口（脚本见 .scratch/wisp/probes/audit-20261114/dupescan.ps1）
#    判据：同一 commit 的 tree 里同票号文件数 >1

# 3) 441 的引证错在哪
git log -1 --format='%h %p' 8e1367bfd               # parent=eb9cf0e63
git cat-file -t eb9cf0e63 a46352fd3 952f4462b 2e8233aeb
gh run view 43571235710 --repo CarlosShao/wisp --json conclusion,jobs
gh run view 43571511743 --repo CarlosShao/wisp --json conclusion

# 4) 441 的字节缺陷
#    （NUL 位置 213/726/854；git 当二进制）
git show --stat 8e1367bfd | grep 441-posix          # -> Bin 3031 -> 4211 bytes

# 5) 274 的红与修
gh run view 43575830504 --repo CarlosShao/wisp --json jobs   # lint: Tracked path-length budget = failure
git show --name-status 6236a6110                    # R100 102 -> 69 字符

# 6) 274 的 paci 步（AC#0(b) 的 CI 读数）
gh run view 43575142588 --repo CarlosShao/wisp --json jobs   # Shipped wisp.exe carries the page bundle = success

# 7) 443 的本机复跑
go test ./internal/plugin/manifestslot -run 'Test52|Test138' -count=1

# 8) path-length 现跑（应 GREEN over 69 / roster 69）
sh scripts/check-path-length-budget.sh

# 9) 441 census 复算（应 37 行 / 13 文件）
git grep -n -E 'RED instrument|instrument is dead|RED INSTRUMENT' 8e1367bfd -- internal
```

---

## 6. 本审计没做的（边界，⛔ 别把这件读成"全验过"）

- **没跑 POSIX 测试**（本机 Windows；441 票面自己也写了"这一格是我的核查缺口"）。441 的运行时读数我只核了代码形状、census、CI 与字节，**没有**自己复现那条弃权。
- **没复跑 `cmd/wisp` 整包**（需 sherpa DLL 暂存；CI 里有 RUN 846／PASS 603 的读数）。442 的 48 例、444 的 9 发争用与 7 发突变，我只核了证据件、日志与 CI，**没有重造突变**。
- **没裁决 Q-440 的"都按推荐"是否真的覆盖那 6 处改文**——那句话只在会话里，不在盘上；本件只报"盘上有这条裁定＋改文与提案逐字一致"。
- **没动任何票面／产码／阈值／golden／契约**；未 `add -A`、未 amend/reset/rebase/stash/checkout、未 push。
- HEAD 在审计期间前进了 2 笔（编排者并行写的 probe/台账）；5 张票的 blob 未变，本件所有票面读数在 `3e9d9779f` 与 `81ba68874` 上同值。
