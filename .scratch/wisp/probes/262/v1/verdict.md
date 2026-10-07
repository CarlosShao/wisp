# 票 262 门禁仪器 — 独立验收裁决表（腿号 `262-v1`）

**本件是什么**：非实现者验收腿 `262-v1` 对票 262 落地的「全仓跟踪路径长度门禁」的**逐格独立取证**。
⛔ **本腿不勾任何 AC 框**（翻框权在编排者）；⛔ 零产码／零脚本／零 CI／零票面／零台账改动。

- 验收腿：`262-v1`（写面唯一＝`.scratch/wisp/probes/262/v1/**`）
- 起手 HEAD：`dea66c57`（`git log --oneline -1` 复认，`dev`）
- 起手锚时刻：`2026-10-07 14:09:10 +08`（`date` 现取）
- 票面行数（起手）：`83`（`wc -l`）
- 被验收仪器：`scripts/check-path-length-budget.sh`（**577 行**，只读）
- CI 接线：`.github/workflows/ci.yml`（941 行，只按行区间读）
- 落地腿读数件：`.scratch/wisp/probes/262/r1/census.md`（530 行，按节 `sed -n` 取）
- 仓外工作拷贝（本腿 AC#3/#4 突变面）：`D:/tmp/262v1/`（⛔ 不入库、不删）

---

## 0. 锚读数（本腿自己现跑，取数时刻 2026-10-07 14:09:16→14:09:20 +08）

命令逐字：

```
sh scripts/check-path-length-budget.sh --with-self-test
```

原始输出＝`logs/gate-head.txt`（19 行）。关键行：

- `rc=0`／`VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`
- `denominator: tracked paths=7405  over-budget=57  covered by roster=57  not in roster=0`
- `longest=180 chars relative (.scratch/wisp/issues/252-the-two-spellings-of-one-path-split-inside-the-allowlist-check-short-name-stays-short-long-name-folds-long-so-in-root-writes-read-as-out-of-bounds-r2-l2.md)`
- `worst full path on the self-hosted runner=224 chars (hat budget 165, wall open interval (206,217])`
- `bands: over the hat=57  of which in the 122..180 middle=57  past the old debt line(180)=0  in the wall interval=0  roster entries=57`
- `unit=characters (probe: length of the 5-byte sample "a U+25FF b" = 3; LANG=C.UTF-8 LC_ALL=C.UTF-8)`
- `hat: rule 9 name cap=100 + issues dir prefix=21 -> relative hat=121 ; worst checkout prefix=44 -> full-path budget=165`
- `unit cross-check: the byte view and the character view select the SAME 57 paths on this tree, recomputed just now in both units`
- 内建正控三发：`control 1/3 ok - bench baseline green`／`control 2/3 ok - the planted over-budget tracked path is rejected and named verbatim:` + 一行 `RED - over budget and NOT in the roster: .scratch/wisp/issues/seeded-…over-the-hat.md  (relative 125 chars, name 104 chars, full path on the self-hosted runner 169 chars, budget 165)`／`control 3/3 ok - green again once the planted path is gone` ⇒ `positive control PASSED`
- bench 落点：`positive control bench=/tmp/tmp.Eeu8TV0TzH (kept on disk; this project never deletes temp artifacts)`

**与编排者 13:5x 那发的对照**：分母四项（7405／57／57／0）与本腿 14:09 现量**逐字一致** ⇒ 本腿不推翻该数；
但本腿把它当**此刻读数**而非常量记账（见 §AC#2 的"分母会动"条款）。

---

## 1. 八格裁决（骨架；取证进行中）

| 格 | 判语 | 凭据 |
|---|---|---|
| AC#0 | 〔进行中〕 | §0＋§2 |
| AC#1 | 〔进行中〕 | §3 |
| AC#2 | 〔进行中〕 | §0＋§4 |
| AC#3 | 〔进行中〕 | §5（仓外 clone 自种自拆） |
| AC#4 | 〔进行中〕 | §6（含两发恒真检查） |
| AC#5 | 〔预判：待取数〕 | §7（10-07 机主"暂时不推远程"） |
| AC#6 | 〔进行中〕 | §8（`./cmd/wisp/` 归编排者补跑） |
| AC#7 | 〔进行中〕 | §9（11＋3＝14 枚逐枚对） |

（以下各节随取证填充。）
