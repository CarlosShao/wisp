# 结案票里藏着的未勾格 —— 只读盘点 `pool-2`（别名 `done-audit-1`）

> 本件是**只读盘点腿**的交付物：不改任何票名、不改任何票面正文、不跑任何编译或门。
> 所有数字都自带口径与命令；命令的工作目录一律 `.scratch/wisp/issues/`（下表简称 `issues/`）。
> 分支 `dev`。生成时刻见文末 §5。

---

## §0 结论摘要（大白话）

（待填）

---

## §1 分母：我自己量的数

### 1.1 票池枚数

| 量 | 命令 | 输出 |
|---|---|---|
| 池内 `.md` 总枚数 | `ls *.md \| wc -l` | **217** |
| 带 `-done` 后缀（＝已结案） | `ls *-done.md \| wc -l` | **78** |
| 不带 `-done`（＝开放） | `ls *.md \| grep -v -- '-done\.md$' \| wc -l` | **139**（**含 `README.md`** ⇒ 真开放票 **138**） |

校验：`78 + 139 = 217`，两把尺对得上。

### 1.2 未勾格 / 已勾格 / 判据行总数

| 量 | 命令 | 输出 |
|---|---|---|
| 未勾格（严格锚定） | `grep -c -- '^- \[ \]' *.md \| awk -F: '{s+=$2} END {print s}'` | **715** |
| 未勾格（允许行首空格／`*` 子弹） | `grep -c -- '^ *[-*] \[ \]' *.md \| awk …` | **715**（与上完全相同 ⇒ 未勾格**没有**缩进或 `*` 变体，锚定这把尺不漏） |
| 已勾格（严格锚定） | `grep -c -- '^- \[x\]' *.md \| awk …` | **534** |
| 已勾格（允许缩进） | `grep -c -- '^ *[-*] \[[xX]\]' *.md \| awk …` | **535**（多出的 1 枚＝`113-…-done.md:66`，票 113 AC#6 的**追加已勾行**，见 §2.3） |
| 含 `[ ]` 的一切行（宽口径） | `grep -c -- '\[ \]' *.md \| awk …` | **769** ⇒ 比锚定的 715 多 **54** 行；这 54 行是**散文里提到方框**（逐条抽查见 §1.3），**不是判据格**，一律不数 |

**判据行总数（两个口径，都给出）**
- 窄口径＝`^- \[ \]` + `^- \[x\]` = **715 + 534 = 1,249**
- 宽口径（含 113 那枚缩进追加勾）＝**715 + 535 = 1,250**
- 恒等式：未勾 + 已勾 = 判据行总数，**在这个定义下天然成立**（判据行＝勾框行本身）；⚠ 但"判据行"不等于"票面上所有 AC 叙述行"，见 §1.4 与 §5 的偏差声明。

### 1.3 结案票 vs 开放票的分摊（本轮的关键数）

| 分摊 | 命令（在 `issues/` 里循环 `grep -c`） | 未勾 | 已勾 |
|---|---|---|---|
| 78 枚 `-done` 票 | `for f in *-done.md; do grep -c -- '^- \[ \]' "$f"; done \| awk …` | **54** | **407** |
| 139 枚非 `-done`（含 README） | 同上排除 `-done` | **661** | **127** |
| `README.md` 自身 | `grep -c -- '^- \[ \]' README.md` | **0** | **0**（README 不贡献判据格） |

⇒ **我这把尺自己算出来的口径：全池未勾 715，其中 661 枚在开放池、54 枚藏在 14 张 `-done` 票里。**
（编排者给的线索数恰好也是 715/661/54；这**不是抄来的**——上面三行命令是本腿独立跑的，逐枚点名见 §2。）

**14 枚涉事票（按未勾格数降序，命令：逐票 `grep -c -- '^- \[ \]'`）**

| 未勾 | 票 | 同票已勾（`^- [x]`） | Status 行（原文第一枚读数） |
|---|---|---|---|
| 7 | 92 panel-composer-mode-attachments-workspace | （见 §2.1） | accepted-done（附条件，条件已归位） |
| 7 | 115 seal-notices-carry-the-resolvers-answer… | （见 §2.2） | BLOCKED 部分 |
| 6 | 113 posix-platformverifypplacement-has-no-link-leg | （见 §2.3） | accepted-done |
| 5 | 97 dead-strict-param-and-the-comment… | **0** | ready-for-review |
| 5 | 110 no-ci-step-runs-internal-winsec | **0** | ready-for-review |
| 5 | 105 c26-rewrite-account-has-no-production-reader | **0** | **rejected-needs-fix** |
| 5 | 104 sealfile-silently-drops-inherited-grants | **0** | accepted-done |
| 4 | 07 ball-state-machine-core | 2 | done |
| 3 | 89 0600-is-decorative-on-windows-acl… | 17 | accepted-done（复验通过） |
| 3 | 80 blacklist-overrides-never-wired-to-gate | 2 | done |
| 1 | 84 pendingapproval-bounded-wait-refuted | 5 | done-as-refutation |
| 1 | 63 credential-entry-cli | 6 | done |
| 1 | 118 winsec-test-hardening-from-ticket104… | 8 | open |
| 1 | 11 llm-adapters-rest | 6 | done |

> ⚠ 表里最扎眼的一列不是未勾，是 Status：**四枚票（97／110／105／115）的 Status 行至今写着 `ready-for-review` / `rejected-needs-fix` / `BLOCKED`，文件名却已经带 `-done`。**
> 逐格读数与出处见 §2、分类见 §3。

### 1.4 方框形状的坑（我数的时候怎么处理）

（待填：`- [2026-…]` 时间戳行、`〔待填〕`、`[H#]` 指针、表内 `**[ ] 维持未勾**`、AC#5 之类）

---

## §2 逐枚点名：14 张 `-done` 票里的 54 枚未勾格

（待填，逐票一节）

---

## §3 逐格分类（甲／乙／丙／丁／戊／己）

（待填）

---

## §4 处置建议名册（只建议，不处置）

（待填）

---

## §5 诚实末节：没读完的、判不了的、我这把尺的偏差

（待填）
