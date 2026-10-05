# evidence-close-1 检查笔记（本腿台件）

本腿＝文书腿，只查盘、只登记真存在的凭据，产出 `docs/evidence/s1/265-267-evidence-index.md`（凭据索引，⛔ 非裁决表）。

## A. 本腿现跑尺名册（命令 → 读数 → 时刻 +08）

| 时刻 | 命令（形状） | 读数 |
|---|---|---|
| 14:19 | `ls docs/evidence/s1/ \| grep -E '26[57]'` | 0 命中（exit=1）＝缺口为真 |
| 14:20 | `grep -c '^- \[x\]'`／`'^- \[ \]'` 两枚票面 | 265：4 枚 `- [x]`／1 枚 `- [ ]`；267：5 枚 `- [x]`／0 枚 `- [ ]` |
| 14:20 | `git log -1 e0790a45`；`git show --name-status e0790a45` | 265 结案/改名 `2026-10-04 20:39`；`R090 ...writer.md→...writer-done.md` |
| 14:20 | `git log \| grep 265-r1c`（subject 起头） | 恰 7 枚 `a04a095f/cd58b810/35b8b613/987fd1f7/7883f855/cfd2722f/f8cb18ef` |
| 14:24 | `git show --name-only` 各 r1c 笔 | §3–§9 全写入 `probes/265/r1/evidence.md`（复用死腿骨架件） |
| 14:24 | `wc -l -c` 265 三件 | census.md 496/89147；a1b/ledger.md 100/14081；r1/evidence.md 140/30303 |
| 14:49 | `git cat-file -s HEAD:...a1/census.md` | 89147＝工作树同值 ⇒ 票面 line32「89,141」为转述笔误（Δ6） |
| 14:24 | `git show --name-status c6cf66e6` | 267 改名 `-done`＝c6cf66e6 10:58；`R076 ...can-vanish.md→...-done.md` |
| 14:24 | `git log \| grep 'ticket 267 r1'` | band 6 枚 `42115f65/37f8e5c6/ec6a47a8/22fc968a/873c3063/0c2445d1`；种子 5 枚 `00f0ef97/32e74479/5fa0d28c/aac52ab2/c7bb02be`（逐枚 `git log -1` 验真） |
| 14:24 | `wc -l -c` 267 七件 | a1/census 81/12579；a2/census 362/77770；a3/census 209/43614；r1/evidence 294/30943；r2/evidence 254/44771；gate/cmdwisp-HEAD.log 1575/242438；r2/cmdwisp-after.log 1673/258923 |
| 14:24 | `git log --diff-filter=A docs/evidence/s1/257-...-v1.md` | 257-v1 首笔 65142d9a 12:21（非实现者终裁表存在） |
| 14:49:04→14:49:32 | `sh scripts/d22scan.sh`（允许） | rc=0 clean；#1-5 internal/=228/cmd/=38；#6 frontend/=85；#7 internal/tools/=23；#8 cmd/=104（含 268-r1 在飞未提交文件＝非静默） |
| 14:51:15→14:51:19 | `sh scripts/check-path-length-budget.sh --with-self-test`（允许） | rc=0 VERDICT GREEN；positive control PASSED；tracked=5957/over-budget=57/roster=57/not-in-roster=0；longest=180；跑后工作树仍只 cmd/wisp M |

⛔ 未跑：`go vet`、`gofumpt`（267 AC#4 四门另两门）＝非允许尺，具名"量不到，归编排者在安静窗口"。

## B. commit 台账与 numstat 删除列（每笔 commit 后回填）

| 笔 | commit | 时刻 | pathspec | `git diff --numstat HEAD~1..HEAD` 删除列 |
|---|---|---|---|---|
| 1 骨架 | `d4e489fa` | 10-05 15:47 | index + notes | index `63\t0`、notes `32\t0` ⇒ **删除列全 0**（新建两枚文件，无删除） |
| 2 §1 | `f75ca2da` | 10-05 15:54 | index | `50\t1` ⇒ 删除那 **1 行＝本腿 §1 自己的"未判（下一步填）。"占位行**（`git show` 逐字核＝`-未判（下一步填）。`）；⛔ 非吞别人行 |
| 3 §2 | `6ac29ddc` | 10-05 15:58 | index | `52\t1` ⇒ 同上，删的是本腿 §2 占位行 |
| 4 §3 | `84099d58` | 10-05 16:02 | index | `11\t1` ⇒ 同上，删的是本腿 §3 占位行 |
| 5 §4＋notes | 待回填 | | index + notes | |

**关于删除列＝1 的说明（本腿逐笔核过）**：每填一节就替掉本腿骨架里那一行 `未判（下一步填）。`。本索引件从头到尾**只有 `evidence-close-1` 自己的两枚 commit 碰过**（`git log --format='%an' -- <件>` 逐笔＝CarlosShao/evidence-close-1，无别腿），所以那个删除列**结构上不可能是别人的行**——是占位行被填实，属预期，⛔ 非"吞掉别人的行"。真正的卫生判据是"没在别人的文件上落笔、没把别人 staged 的东西带进 pathspec"——本腿每笔 pathspec 只点名 index（＋notes），且 commit 前 `git diff --cached --name-only` 现量为空（共享 index 无别人 staged 内容）。

## C. 交付回报要点（终值以 commit 后现量为准）

- 新建索引件：`docs/evidence/s1/265-267-evidence-index.md`（本腿唯一产品件）。
- 现量到的盘-账不符：**共 2 处**，都在票 265——① census.md 字节 89,147（HEAD blob 同值）vs 票面 line32「89,141」（Δ6，行 496 一致）；② AC#0 物理框 `- [ ]`（line29）vs 台账 A607/e0790a45「五格全勾」。票 267 侧四门/越界尺与 a1/a2/r1/r2 各件字节**全部复跑对上、零处不符**。
- ⛔ 未跑：`go vet`／`gofumpt`／任何 `go` 命令；已跑且仅此两把＝`d22scan.sh`（rc0 clean @14:49:32）、`check-path-length-budget.sh --with-self-test`（rc0 GREEN @14:51:19）。
