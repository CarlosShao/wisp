# evidence-close-7（搬运腿 `backfill-7a`）— 状态件

**起手时刻**＝2026-10-06（`date -Iseconds` 见本件末尾「起手现量」一节）
**任务性质**＝把已在盘的凭据**指进**对应工单名下（搬运）；⛔ 零新读数、零判断、零产码、零翻勾。
**授权来源**＝编排者裁定动作清单（动作 1＝4 枚票各追加一节「凭据归口（搬运笔，非新验收）」；动作 2＝给 `docs/evidence/s1/closed-tickets-evidence-index.md` 的 §3.1 表下方插一段更正说明）。

## 1 存在性复量（本腿自跑 `test -f`／`wc -l`／`ls`，2026-10-06，一次跑齐）

| 编排者给的路径 | 在盘 | 复量行数／枚数 | 与编排者口径 |
|---|---|---|---|
| `docs/evidence/s1/75-independent-verification.md` | 是 | 1052 行 | 一致（1052） |
| `docs/evidence/s1/81-adversarial-acceptance.md` | 是 | 56 行 | 一致（56） |
| `docs/evidence/s1/82-adversarial-acceptance.md` | 是 | 511 行 | 一致（511） |
| `docs/evidence/s1/85-preflight-staticcheck.md` | 是 | 505 行 | 一致（505） |
| `docs/evidence/s1/ci-step-readings-2026-09-22.md` | 是 | 343 行 | 一致（343） |
| `docs/evidence/s1/248-settings-write-path-v1.md` | 是 | 279 行 | 一致（279） |
| `docs/reports/HANDOVER.md` | 是 | 1836 行 | 编排者未给行数，本腿现量 |
| `.scratch/wisp/probes/256/v1/verdict.md` | 是 | 223 行 | 一致（223） |
| `.scratch/wisp/probes/254/r1b/logs/` | 是 | 14 份文件 | 一致（14） |

**枚数级不一致一枚（⛔ 本腿不据此改任何票面句子，只登记）**：票 254 那一条编排者写「`:68`／`:278` **两处**出现 `254-r1b` 名号」，
本腿 `grep -n '254-r1b' docs/evidence/s1/248-settings-write-path-v1.md` 命中 **3 行**（`:68`／`:179`／`:278`）。
⇒ 给定的两行**逐字成立**（存在性对得上），差的是"两处"这枚计数；多出的 `:179` **不列进本票凭据**（⛔ 不许自己新增别的凭据），只在票面归口节里具名报这一枚计数差。

## 2 做到哪一枚／还剩哪几枚（截至本件落笔）

| 动作 | 对象 | 状态 |
|---|---|---|
| 动作 1-① | 票 78 `.scratch/wisp/issues/78-linux-vet-buildtags-unblocks-d22-gate-done.md` | **已追加归口节**（178→198 行，＋20 行／＋1,191 字节），读回已过 |
| 动作 1-② | 票 115 `.scratch/wisp/issues/115-…-done.md` | **已追加**（461→479 行，＋18 行／＋1,726 字节），读回已过 |
| 动作 1-③ | 票 254 `.scratch/wisp/issues/254-fifth-tier-winsec-no-ci-caller-done.md` | **已追加**（47→72 行，＋25 行／＋2,174 字节），读回已过 |
| 动作 1-④ | 票 267 `.scratch/wisp/issues/267-…-done.md` | **已追加**（87→107 行，＋20 行／＋1,744 字节），读回已过 |
| 动作 2 | `docs/evidence/s1/closed-tickets-evidence-index.md` §3.1 表下方插更正引用块 | **已插入**（354→361 行，＋7 行／＋1,041 字节），删除列＝0，表里 17／74／2／93 一字未动 |
| 提交 | 6 枚文件逐枚点名 commit（⛔ 不 push） | 本件＝**第 1 笔（状态件闸门）**；5 枚交付件的第二笔待落；落完再追加第 3 笔登记 commit 号 |

**剩余工作清单（若本腿此刻死掉，接手腿照这个做）**：
1. `git commit -F <msg.txt> -- <5 枚点名路径>`（⛔ 不含本 status 件亦可，或同笔带上；禁 `add -A`／`.`／`--amend`）。
2. 把两笔 commit 号回填进本件 §5，并 `git log --oneline -3` 复认。
3. 回报编排者：两处对不上的地方（票 254 的计数差、票 267 的行号差）——已各自写进票面，不需另立件。

## 3 闸门自查

- ⛔ 零 Go 命令、零 `sh scripts/*.sh`（本腿不跑任何门禁）。
- ⛔ 不碰任何 `- [ ]`／`- [x]`；不改 `-done` 文件名；不动台账／`AGENTS.md`／`PLAN.md`／`docs/specs/**`。
- ⛔ 禁读面：`.scratch/wisp/probes/232/**`／`111/ci-read-1/**`／`pool-validity/**`／`frontend/**`／`design/**`（本腿未读、未引）。
- 本件＝硬预算闸门的落点：任何时刻被打断，本节即"做到哪一枚／还剩哪几枚"的答案，⛔ 不留「（待填）」。

## 4 起手现量（逐字，2026-10-06）

- `date -Iseconds`（第一次）＝`2026-10-06T18:47:20+08:00`；（复量）＝`2026-10-06T18:47:50+08:00`
- `git log -1`＝`257275be 2026-10-06T18:39:10+08:00 probes(收 236-v1＋翻票 236 三格＋立票 270/271＋A643)`，分支＝`dev`
- `mkdir -p` 建了本腿台件目录 `.scratch/wisp/probes/evidence-close/7/logs/` 与 `D:/tmp/wispbf7a/`（临时件**只建不删**）
- `git status --porcelain -- .scratch/wisp/probes/evidence-close docs .scratch/wisp/issues` 起手可见：
  `evidence-close/1/**` 五枚未入库件（前一枚腿 `evidence-close-6` 的台件，⛔ 不属本腿、本腿不代提交）＋`?? .scratch/wisp/probes/evidence-close/7/`（本腿自己）。
- 本腿落笔纪律：中文正文一律走 `Write`／`Edit` 工具（⛔ 不用未加引号的 heredoc）；commit message 先 `Write` 成 `.txt` 再 `git commit -F`；逐枚 pathspec 点名。
