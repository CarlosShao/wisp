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

未判（下一步填）。

---

## §2 票 267（凭据在哪）

未判（下一步填）。

---

## §3 对照：有终裁的票才有表，没有终裁的票只有索引

未判（下一步填）。

---

## §4 本索引不解决什么（仍欠的非实现者读数）

未判（下一步填）。
