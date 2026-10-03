# 票 252 · AC#2「两侧同形」修法 · 非实现者对抗验收表（腿 `252-v1`）

被验收件：实现腿 `252-r1` 的 `internal/tools/paths.go`（`Canonicalize` ＋新函数 `sameFormOfUnresolved`）
与载具 `internal/tools/paths_shortname_252_r1_test.go`；提交链 `bbbb3593`→`34d89191`→`2b1a3071`→`ebd5da2b`。
本腿＝对抗验收腿，**不翻任何 AC 框**；判语只写"成立／不成立／附条件成立"。

---

## §0 起手锚（同发取，逐字）

| 项 | 读数 |
|---|---|
| 取锚时刻 | `2026-10-03 09:38:07 +0800` |
| `git rev-parse --short=8 HEAD` | `1fea1e97` |
| HEAD subject | `145-c2: S4 reconfirm/overturn against r1/r2/r3 (fresh line reads)` |
| 全仓 porcelain 行数 | `405`（`.scratch/wisp/probes/252/v1-tmp-porcelain.txt`，均为他人台件与 design 脏面） |
| 本腿写面基线 `git status --porcelain -- internal tools` | **空（0 行）**＝`.scratch/wisp/probes/252/v1/start-porcelain-internal-tools.txt` |
| 分支 | `dev`（`git rev-parse --abbrev-ref HEAD`） |

四枚被验收提交的存在性与 ancestry：逐枚 `git log -1` 均命中；`git merge-base --is-ancestor ebd5da2b HEAD` ⇒ YES。
`bbbb3593^..HEAD` 里动过 `internal/tools` 的只有本票三枚（`34d89191`/`2b1a3071`/`ebd5da2b` 的 docs 枚）。

---

## §1 逐格判语（每格三选一＋凭据）

| # | 格 | 判语 | 凭据 |
|---|---|---|---|
| 1 | 「两次包含没合一」复量 | <待填> | <待填> |
| 2 | fail-closed `!ok` 未改软＋M2 复跑 | <待填> | <待填> |
| 3 | ★落点偏差是否算遵守裁的形（详见 §2） | <待填> | <待填> |
| 4 | 存在性分岔有没有被改出新岔 | <待填> | <待填> |
| 5 | 正控的恒真自查（换形还绿不绿） | <待填> | <待填> |
| 6 | 越界名册（详见 §3） | <待填> | <待填> |

---

## §2 ★第 3 问：偏差裁决（算不算遵守编排者裁的形）

- 裁语原文出处：`.scratch/wisp/issues/252-…-r2-l2.md:83`
- 实现腿自述落点：`.scratch/wisp/probes/252/r1/impl.md:95-98`、自陈残余不对称 `:121-125`、自请裁定 `:183-185`、`:252`
- 本腿复跑作用面读数：见 §4
- 裁决（一句话）：<待填>

---

## §3 越界名册（`git diff --name-only bbbb3593^..ebd5da2b` 逐枚）

<待填>

---

## §4 突变自证表（种什么→哪枚必须红→红句逐字→还原复跑终值）

<待填>

---

## §5 我建议编排者怎么处置

<待填>

---

## §6 我不确定的地方

<待填>

---

## §7 判语

<待填>
