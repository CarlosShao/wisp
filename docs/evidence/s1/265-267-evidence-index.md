# 凭据索引 — 票 265 与票 267（已结案，但 `docs/evidence/s1/` 里没有对应的表）

本件不是对抗验收表，也不是缺口审计；它只是凭据索引。这两枚票的非实现者终裁在写本件时尚未发生（票 267 的 `267-v1` 仍在队列、票 265 的验收腿未派），任何一格都不因本件而视为已验收。

> 为什么有这一格：票 265（`2026-10-04 20:39` 改名 `-done`）与票 267（`2026-10-05 10:58` 改名 `-done`）都由编排者本人核过结案、五格全翻（票 265 的 AC#0 见 §1.4 的一处盘-账不符），但 `docs/evidence/s1/` 里没有一枚以它们编号命名的件（`ls docs/evidence/s1/ | grep -E '26[57]'`＝零命中，本腿现量）。它们的凭据散在真相源台账 `A607`／`A615` 两节和 `.scratch/wisp/probes/**` 的实现件里。本件只做一件事：**把"做了的读数在盘上哪一枚文件的哪一节"指出来**，⛔ 不裁任何一格、⛔ 不冒充终裁。

---

## §0 起手锚与本腿读数名册

### 0.1 起手锚（本腿现跑）

| 项 | 读数（逐字，本腿现量） |
|---|---|
| 本腿代号 | `evidence-close-1`（文书腿，只查盘、只登记真存在的凭据） |
| 进场时刻 | `2026-10-05 14:21:18 +0800`（`date` 自取） |
| 进场 HEAD | `1711ad06`（`probes(167-a5 起手)`，2026-10-05 14:20） |
| 分支 | `dev`（`git rev-parse --abbrev-ref HEAD` 现量） |
| 此刻写面（非本腿） | `cmd/wisp/resident_approval_windows.go`＝` M`（在飞腿 `268-r1`，⛔ 本腿不碰）；`scripts/` 干净（`git status --porcelain -- cmd internal tools scripts docs .github` 只此一枚） |
| 本腿写面（授权） | 仅两枚：`docs/evidence/s1/265-267-evidence-index.md` ＋ `.scratch/wisp/probes/evidence-close/1/notes.md`；⛔ 不改票面任何字、不碰台账、不碰代码／测试／规格 |
| 禁跑项 | 全程⛔ 无 `go test/build/vet/env`；只跑派单点名的两把壳尺（见 §0.2 第 11/12 行） |

### 0.2 本腿读数名册（逐条：命令 → 读数 → 时刻）

| # | 命令（形状） | 读数 | 时刻（+08） |
|---|---|---|---|
| 1 | `ls docs/evidence/s1/ \| grep -E '26[57]'` | **0 命中**（exit=1）＝缺口为真 | 14:19 |
| 2 | `grep -cE '^- \[[x]\]' / `^- \[[ ]\]`` 两枚票面 | 票 265：`- [x]` **4** 枚／`- [ ]` **1** 枚；票 267：`- [x]` **5** 枚／`- [ ]` **0** 枚 | 14:20 |
| 3 | `git log -1 --format=... e0790a45` | 票 265 结案／改名 `-done`＝`2026-10-04 20:39`；`git show --name-status` 含 `R090 ...265-...-writer.md → ...-writer-done.md` | 14:20 |
| 4 | `git log \| grep '265-r1c'`（subject 起于 `265-r1c`） | 恰 **7** 枚：`a04a095f`→`cd58b810`→`35b8b613`→`987fd1f7`→`7883f855`→`cfd2722f`→`f8cb18ef`（与 A607"七笔 commit `a04a095f`→`f8cb18ef`"逐枚对上） | 14:20 |
| 5 | `git show --name-only` 各笔 | 7 枚 §3–§9 全写进 `.scratch/wisp/probes/265/r1/evidence.md`（非新文件，复用死腿 `265-r1` 的骨架件） | 14:24 |
| 6 | `wc -l -c` 票 265 三枚凭据件 | `a1/census.md`＝**496 行／89,147 字节**；`a1b/ledger.md`＝**100／14,081**；`r1/evidence.md`＝**140／30,303** | 14:24 |
| 7 | `git cat-file -s HEAD:...a1/census.md` | **89,147**＝与工作树相等 ⇒ 票 265 面 line 32 写的"89,141 字节"是**转述笔误（Δ6 字节）**，行数 496 一致 | 14:49 |
| 8 | `git log --follow` 票 267 `-done` 文件 ＋ `git show --name-status c6cf66e6` | 改名 `-done` 那一笔＝`c6cf66e6`（`2026-10-05 10:58`，ledger A614＋A615），含 `R076 ...can-vanish.md → ...can-vanish-done.md` | 14:24 |
| 9 | `git log \| grep 'ticket 267 r1'`（band 六笔） | `42115f65`→`37f8e5c6`→`ec6a47a8`→`22fc968a`→`873c3063`→`0c2445d1`（恰 6 枚，与 A611"六笔 commit"对上）；种子迁移五笔 `00f0ef97`→`32e74479`→`5fa0d28c`→`aac52ab2`→`c7bb02be` 逐枚 `git log -1` 验真在 | 14:24 |
| 10 | `wc -l -c` 票 267 六枚凭据件 | `a1/census.md`＝81／12,579；`a2/census.md`＝362／77,770；`a3/census.md`＝209／43,614；`r1/evidence.md`＝294／30,943；`r2/evidence.md`＝254／44,771；`gate/cmdwisp-HEAD.log`＝1,575／242,438；`r2/cmdwisp-after.log`＝1,673／258,923 | 14:24 |
| 11 | `sh scripts/d22scan.sh`（**允许跑**） | **rc=0／clean**；分母 `bans #1-5 internal/=228, cmd/=38; ban #6 frontend/=85; ban #7 internal/tools/=23; ban #8 cmd/=104`（⚠ `cmd/=104` 含在飞 `268-r1` 那枚未提交文件＝非静默树读数） | 14:49:04→14:49:32 |
| 12 | `sh scripts/check-path-length-budget.sh --with-self-test`（**允许跑**） | **rc=0／VERDICT GREEN**；positive control PASSED；`tracked=5957, over-budget=57, covered by roster=57, not in roster=0`；`longest=180`（`252-…` 那枚）；跑完工作树仍只 `cmd/wisp M`＝自检临时件无残留 | 14:51:15→14:51:19 |
| 13 | 范本对照：`git log --diff-filter=A docs/evidence/s1/257-...-v1.md` ＋ `grep -cE` 257 票面 | `257-v1` 验收腿首笔＝`65142d9a`（12:21）；票 257 五枚 AC 框全 `- [x]`；`docs/evidence/s1/257-clean-machine-provider-registry-v1.md`＝**真·非实现者终裁表**（本件 §3 用它当形状范本） | 14:24 |

---

## §1 票 265（凭据在哪）

票面：`.scratch/wisp/issues/265-resident-gate-has-no-session-grant-writer-done.md`（本腿⛔未改一字）。

### 1.1 结案那一笔＋交件那一串 commit（逐枚号，本腿现跑 `git log -1`）

- **结案／改名 `-done`＝`e0790a45`（2026-10-04 20:39）**。`git show --name-status e0790a45` 现量含 `R090 ...265-...-writer.md → ...-writer-done.md`；其 subject 逐字起于「265 结案：五格全勾（AC#0 19:3x/AC#1 14:3x/AC#2-4 20:4x 编排者核档翻勾）＋改名 -done」。
- **交件七笔（台账 A607 记「七笔 commit `a04a095f`→`f8cb18ef`」，本腿逐枚验真，subject 皆以 `265-r1c` 起头，恰 7 枚）**：

| 笔 | commit | 时刻 | 内容（subject 摘） |
|---|---|---|---|
| 1 | `a04a095f` | 10-04 19:53 | 收编：前任 `265-r1` 遗留的 ⓐ-Ⅰ 产码四枚**原样入库（收编，非新产）** |
| 2 | `cd58b810` | 10-04 19:58 | §3 判据件 8 枚用例逐枚判读＋256 解冻 diff 逐块核（A602 边界内，越界 0） |
| 3 | `35b8b613` | 10-04 20:02 | §4 四发突变自证全兑现（M1–M4，md5 三枚全等×4） |
| 4 | `987fd1f7` | 10-04 20:03 | §5 两句假话兑现——字节级零改动双尺证 |
| 5 | `7883f855` | 10-04 20:06 | §6 门禁五读数全绿带时刻（含 phantom-citation 自修） |
| 6 | `cfd2722f` | 10-04 20:34 | §7 判不动四格具名归口 |
| 7 | `f8cb18ef` | 10-04 20:35 | §8/§9 污染面名册四枚＋提交名册七枚＋票面 AC 框零改动复证（checkbox diff=0） |

- 这七笔的 §3–§9 **全写进同一枚件** `.scratch/wisp/probes/265/r1/evidence.md`（本腿 `git show --name-only` 逐笔现量；⚠ 它不是新文件，是死腿 `265-r1` 立的骨架件被 `265-r1c` 续写填实）。

### 1.2 凭据件路径与 `wc -l -c`（本腿现量）

| 件路径 | 行 | 字节 |
|---|---|---|
| `.scratch/wisp/probes/265/a1/census.md`（普查件，`265-a1` 13:33 立骨架 `0b938773`、末笔 `265-a1b` 17:56 `89aa72d5`） | 496 | **89,147** |
| `.scratch/wisp/probes/265/a1b/ledger.md`（接续腿自对抗台账） | 100 | 14,081 |
| `.scratch/wisp/probes/265/r1/evidence.md`（`265-r1c` 交件证据件） | 140 | 30,303 |
| `.scratch/wisp/probes/265/a1b/logs/00..19-*.txt`（读数台件 20 枚） | — | — |

### 1.3 票面几枚 AC 框、由谁翻（原话位置给 file:line）

- 票面共 **5 枚 AC 框**（AC#0–AC#4）。本腿现量的**物理框态**：AC#1（line 33）、AC#2（line 36）、AC#3（line 38）、AC#4（line 40）＝ `- [x]`（4 枚）；**AC#0（line 29）＝ `- [ ]`（1 枚未翻）**。
- **全部翻勾者＝编排者本人**（⛔ 无非实现者腿翻过，票 265 的验收腿未派）。翻勾时刻与凭据位置：
  - AC#1＝`14:3x`（`docs/reports/pending-and-issues.md:11701` A601，"三档选「每次都重新问」"）。
  - AC#0＝`19:3x` **仅在票面子弹注记为「AC#0 闭合」**（`.scratch/.../265-...-done.md:32`），**物理框未翻**（见 §1.4）。
  - AC#2／AC#3／AC#4＝`20:4x`（票面 line 37／39／41 各带"✔ 翻勾（20:4x 编排者…）"；同一时刻见 `e0790a45` subject 与 `docs/reports/pending-and-issues.md:11871` A607 §2）。

### 1.4 一处盘-账不符（具名，本腿现量，⛔ 不改票面一字）

1. **census.md 字节数**：票面 line 32 写「496 行／**89,141** 字节」；本腿 `wc -c`＝**89,147**，`git cat-file -s HEAD:...census.md`＝**89,147**（工作树＝HEAD，该件末笔 commit 为 17:56、早于 19:3x 翻勾 ⇒ 现值即当时值）。**行数 496 一致、字节差 6 ⇒ 判为编排者转述笔误，非件被改动。**
2. **AC#0 物理框**：票面 line 29 的框至今是 `- [ ]`，而台账 A607 §2（`:11871`）与结案 commit `e0790a45` 均称「五格全勾」。**差别在形状上**：AC#0 的"闭合"只落在票面 line 32 的**子弹注记**里、编排者**没把那一行 markdown 复选框从 `[ ]` 翻成 `[x]`**。本腿按票面**物理态**如实登记（4×`[x]`＋1×`[ ]`），⛔ 不据台账替它翻、也不据它判未结案。

### 1.5 残余几格、各自归口到哪

出自票面 line 41（AC#4 子弹"判不动四格在 §7 具名"，凭据件 `r1/evidence.md` §7）与 `docs/reports/pending-and-issues.md:11869`（A607 §1 第四条，编排者逐格核过归口）：

- ① `internal/agent/approval/approval_always.go:165` 第三枚过期注释＝不扩包 ⇒ **归口编排者**（改它需另开写面窗，见票 256 收档账 A603 三枚注释清单）。
- ② 票 127 钉 `:497` 不 poll 竞态＝**带载偶发**（隔离 ×4 绿）。
- ③ 223／AC1-shutdown 两枚 `count=3` 复绿＝**带载偶发**。
- ④ **U9 真进程尺与 `-race` 未跑**（真开窗未批＝等 `winlive` 那个词）⇒ A607 §1 把它**升格为 winlive 批准后的第一件**（U9 是 AC#1 判语里"唯一能证伪『每次都问』这一档"的尺，仍未跑）。
- 另：本票**未选 ⓐ-Ⅲ（第二枚 mint）**、**未走 ⓒ（登记 DEFERRED）**——ⓒ 那一格由编排者**追加归票 225**（`docs/reports/pending-and-issues.md` A601 段与票面 line 53），本票不建 `docs/DEFERRED.md`。

---

## §2 票 267（凭据在哪）

票面：`.scratch/wisp/issues/267-confirm-timeout-config-has-no-bounds-so-the-c18-warning-can-vanish-done.md`（本腿⛔未改一字）。

### 2.1 结案那一笔＋交件两串 commit（逐枚验真）

- **改名 `-done`＝`c6cf66e6`（2026-10-05 10:58）**：`git show --name-status` 现量含 `R076 ...can-vanish.md → ...can-vanish-done.md`，subject 逐字起于「ledger(A614＋A615) 票 267 五格全勾改名 -done」。⚠ 票面「结案」节（line 73）自报**锚 `c7bb02be`**（那是 `267-r2` 的第五笔、编排者落笔结案prose时的树顶）——**改名那一笔本身是 `c6cf66e6`**，两枚号都真、各指一事，本腿如实并列。
- **band 六笔**（subject 皆以 `ticket 267 r1:` 起头，与 A611"六笔 commit"对上）：`42115f65`(09:15 骨架)→`37f8e5c6`(09:25 值域进门 shape a)→`ec6a47a8`(09:30 pins 改字面种子)→`22fc968a`(09:33 §0/6/7)→`873c3063`(09:33 外置 config_test 守卫)→`0c2445d1`(09:41 §1–§5)。
- **种子迁移五笔**（subject 皆以 `267-r2` 起头，与 A615 §1"五笔自落 commit"对上）：`00f0ef97`(10:05 骨架)→`32e74479`(10:11 种子迁移，凭 A611 具名解冻)→`5fa0d28c`(10:29 §0–§7)→`aac52ab2`(10:33 §8–§10＋19 枚台件)→`c7bb02be`(10:36 正名笔)。
- 产码面＝`internal/config/**` 五枚（`validate.go`／`validate_test.go`／`validate_267_test.go`／`schema.go` 注释／`unwired.go` 一行）＋`cmd/wisp/**_test.go` 五枚；生产码零改动（见 §2.4 本腿复跑的越界尺）。

### 2.2 凭据件路径与 `wc -l -c`（本腿现量）

| 件路径 | 行 | 字节 | 出处对账 |
|---|---|---|---|
| `.scratch/wisp/probes/267/a1/census.md`（普查，编排者代落盘） | 81 | 12,579 | 票面 line 46 逐字「81 行／12,579 字节」**符** |
| `.scratch/wisp/probes/267/a2/census.md`（§7.1 甲/乙/丙/丁料） | 362 | 77,770 | A612（`:11963`）「362 行／77,770 字节」**符** |
| `.scratch/wisp/probes/267/a3/census.md`（`Q-77` 三维料） | 209 | 43,614 | A613（`:11976`） |
| `.scratch/wisp/probes/267/r1/evidence.md`（band 交件件） | 294 | 30,943 | 票面 line 57／A611「294 行／30,943 字节」**符** |
| `.scratch/wisp/probes/267/r2/evidence.md`（种子迁移交件件） | 254 | 44,771 | A615 §1（`:12020`）「254 行／44,771 字节」**符** |
| `.scratch/wisp/probes/267/gate/cmdwisp-HEAD.log`（基线 16 枚 FAIL 台件） | 1,575 | 242,438 | 票面 line 80／A615 §3「1,575 行」**行符** |
| `.scratch/wisp/probes/267/r2/cmdwisp-after.log`（腿 `-v` 台件，四数 RUN=323/PASS=226/FAIL=1/SKIP=0 的唯一出处） | 1,673 | 258,923 | A615 §3 末段 |
| `.scratch/wisp/probes/267/gate/orch-d22scan.log`／`orch-pathlen.log` | — | — | 编排者 10:5x 四门中的两门读数台件 |

⚠ **本腿复跑的这三件字节数与台账逐字一致**（81/12579、362/77770、294/30943、254/44771），与票 265 census.md 那处 Δ6 不同——票 267 侧的登记没有字节漂移。

### 2.3 票面几枚 AC 框、由谁翻

- 票面共 **5 枚 AC 框**，本腿现量物理态：AC#0–AC#4（line 27–31）**全部 `- [x]`＝5 枚**，`- [ ]` **0 枚**（与台账"五格全勾"一致，无 §1.4 那种框态差）。
- **翻勾者＝编排者本人**（⛔ 非实现者终裁未发生，见 §2.5）：AC#0＝`f761a017`（09:12 翻勾＋裁形 ⓐ）；AC#1／AC#2／AC#3＝`9e24d178`（09:51，凭据＝`267-r1` 证据件，见 A611 `:11937`）；AC#4＝结案那一笔 `c6cf66e6`（10:58，凭据全为编排者现跑，见 A615 §5 `:12032`）。

### 2.4 ★AC#4 那格的凭据＝编排者本人现跑的四门＋越界尺（位置指到＋本腿复跑两把壳尺）

- **原始读数位置**：`docs/reports/pending-and-issues.md:12023`（A615 §2）＝四门（`d22scan.sh` rc0 clean、`check-path-length-budget.sh` rc0 GREEN、`go vet ./cmd/wisp/ ./internal/config/` rc0、`gofumpt` 只剩预存 `cmd/wisp/models.go`）＋越界尺（`git diff --name-only 21bec8a1..HEAD` 禁区命中 0 行）。票面「结案」节 item 2（line 78）同记。⚠ 四门里 `go vet`／`gofumpt` 那两门的**原始读数是编排者现跑**（本腿未复跑，见下）。
- **本腿现跑复认（两把允许重跑的壳尺，见 §0.2 row 11/12）**：
  - `sh scripts/d22scan.sh`（14:49:04→14:49:32）＝**rc=0／clean**，分母 `#1-5 internal/=228、cmd/=38；#6 frontend/=85` **与 A615 §2 逐枚对得上**；⚠ `ban #8 cmd/=104` 含在飞 `268-r1` 那枚未提交文件＝**非静默树读数**。
  - `sh scripts/check-path-length-budget.sh --with-self-test`（14:51:15→14:51:19）＝**rc=0／VERDICT GREEN**，`over-budget=57／covered by roster=57／not in roster=0／longest=180`，positive control PASSED；跑完工作树仍只 `cmd/wisp M`＝自检临时件无残留。
  - 越界尺本腿**独立复跑**（git 只读，`git diff --name-only 21bec8a1..c6cf66e6` 存 `267-boundary-names.txt`）：禁区集（`docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／`frontend/**`／`design/**`／`internal/agent/approval/**`）**命中 0**、`internal/config/**` 恰那 **5 枚** ⇒ A615 §2 越界那半**复现**。（⚠ 本腿把 range 收到本票全链终点 `c6cf66e6`，区间夹了别腿交件共 55 枚文件，与编排者当发"非 `.scratch` 只 11 枚"**不同口径**，但禁区命中同为 0。）
- **量不到的两门，具名归编排者在安静窗口**：`go vet`／`gofumpt` 两门本腿**跑不了**——它们要么是 `go` 命令（本腿全程禁跑）、要么需 `gofumpt` 二进制，而此刻 `268-r1` 正在 `cmd/wisp` 上取整包终态（一发约 8.5 分钟、计时敏感）。并发跑会把别人的读数洗成假红，故⛔ 不跑、如实标"量不到"。
- ⛔ **本腿这些复跑≠非实现者终裁**：AC#4 的原始四门本就是"实现者＝编排者"同体读数，本腿复跑只是**确认读数没漂**，不使 267 因此视为已验收。

### 2.5 `267-v1`（非实现者终裁）仍未发生、仍在队列

- `ls docs/evidence/s1/ | grep -E '26[57]'`＝0 命中（本腿现量）⇒ `docs/evidence/s1/` 里**没有** `267-v1` 表。
- 队列依据（逐字在册）：`docs/reports/pending-and-issues.md:12033`（A615 §5 末段）「本票结案≠没有验收：`267-v1`（非实现者终裁…）此刻在队列里、按推送窗口之后派」；`:12035` 队列序「推送 → `257-r2` → **`267-v1`** → `268-r1` → `167-r2`」。
- 若 `267-v1` 推翻任一格：按 A615 §5「追加更正、不改写已入库的 commit」处置。

### 2.6 残余五格、各自归口（票面 line 87＋A615 §5 `:12032`）

- ① `confirm_timeout_sec` 用"拒载"而 `l1_window_sec` 用"钳位"的**行为不一致**⇒ 归**下一枚治理票**。
- ② 下界抬到 60 s 那支＝要改票 256 的种子钉 ⇒ **待编排者具名解冻**（本票未做）。
- ③ `WarningLead` 生产里恒零值、全靠兜底 ⇒ 归**票 255** 的账。
- ④ 常驻腿把"越界拒载"读成"文件读不到" ⇒ **票 268**（料已由 `268-a1` 量齐：`.scratch/wisp/probes/268/a1/census.md`，见 A614）。
- ⑤ `Q-77`（配置值 vs C18 写死 300 s 谁优先）＝**待机主一句话**，⛔ 任何腿不许自裁。

---

## §3 对照：有终裁的票才有表，没有终裁的票只有索引

同样今天（10-05）结案的票 **257** 有一份**真·非实现者终裁表**在 `docs/evidence/s1/257-clean-machine-provider-registry-v1.md`——本腿现量：该件首笔 commit＝`65142d9a`（2026-10-05 12:21，subject 逐字「票 257 验收腿 `257-v1` 第 1 步」），出自非实现者腿 `257-v1`（`SPEC-12 §4.3` #1／`AGENTS.md` §0.3 的双角色），它按 AC#1–AC#4 逐格出**判语**、种了 M1–M5 五发实弹突变、自跑门禁四门。票 257 票面五枚 AC 框＝全 `- [x]`（本腿现量 line 16–20）。

**把三枚摆在一起看形状差（⛔ 这个差别本身要留在盘上，不许用索引抹平）**：

| 票 | 结案方式 | `docs/evidence/s1/` 里有表吗 | 非实现者终裁 | 本腿这次给什么 |
|---|---|---|---|---|
| **257** | 编排者翻框＋`257-v1` 终裁 | **有**：`257-clean-machine-provider-registry-v1.md`（真裁决表） | **已发生**（`65142d9a` 12:21） | —（不需索引，它有表） |
| **265** | 编排者自核翻框（`e0790a45` 20:39） | **无**（`grep 26[57]`＝0 命中） | **未发生**：验收腿未派 | **本件＝凭据索引** |
| **267** | 编排者自核翻框（`c6cf66e6` 10:58） | **无**（同上） | **未发生**：`267-v1` 仍在队列 | **本件＝凭据索引** |

⇒ 判语：**有终裁的票才有表，没有终裁的票只有索引**。257 之所以能有一枚 AC 格逐条判语＋突变红句的表，是因为有一枚**不是实现者**的腿（`257-v1`）去攻过它；265／267 今天只有**实现者＝编排者同体**的自核读数＋台账登记，⛔ 那还不等于对抗验收。本件把"做了的读数在哪"指到本仓要求的落点（`docs/evidence/s1/`），**但指路不等于裁决**——本件任何一格都不因写了凭据位置而视为已验收（承首段性质声明）。

---

## §4 本索引不解决什么（仍欠的非实现者读数）

未判（下一步填）。
