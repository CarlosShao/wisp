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

未判（下一步填）。

---

## §3 对照：有终裁的票才有表，没有终裁的票只有索引

未判（下一步填）。

---

## §4 本索引不解决什么（仍欠的非实现者读数）

未判（下一步填）。
