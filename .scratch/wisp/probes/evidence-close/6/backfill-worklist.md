# evidence-close-6 — 真缺凭据的 5 枚（78／115／250／254／267）逐枚补档工作单

> 本腿＝只读普查腿 `evidence-close-6`。唯一问题＝**每一枚要补哪一档、补它要什么形状的证据、能不能用盘上已有件顶**。
> ⛔ 本腿不重做 `evidence-close-5` 的 17 枚分诊（只读了它的**名册**（`5/c17-triage.md` §1，`:133-163`）与**归口表**（同件 §2.18，`:400-418`）两块；
> 它对票面的判语一枚未读，所以下面每一枚都是本腿自己看票。
> ⛔ 本腿零 `go` 命令、票面零改动、不翻勾、不改名、不动台账与 `docs/**` 一字；唯一写面＝本目录 `.scratch/wisp/probes/evidence-close/6/**`。
> 代号 `6` 起手时未用过（`ls .scratch/wisp/probes/evidence-close/`＝`1 2 3 4 5`），未启用 `6b`。

---

## §0 起手锚（三把尺原文＋取数时刻）

- `date '+%m-%d %H:%M'` ＝ **10-06 16:56**（起手）／落笔时复量 **10-06 17:05**；分支 `dev`。
- `git rev-parse --short HEAD` ＝ 起手 **`de084db7`** → 写件时 **`34ae7b05`** ⇒ **普查期间 HEAD 漂了一枚**（有腿在飞）。
  ⚠ 本件所有枚数一律"按落笔时刻引"，不复用别腿的读数。
- `git status --porcelain -- docs .scratch/wisp/probes/evidence-close` 原文（起手 16:56 与 17:05 两遍**逐字相同**）：

```
?? .scratch/wisp/probes/evidence-close/1/267-boundary-names.txt
?? .scratch/wisp/probes/evidence-close/1/gate-d22scan.log
?? .scratch/wisp/probes/evidence-close/1/gate-pathlen.log
?? .scratch/wisp/probes/evidence-close/1/msg-01-skeleton.txt
?? .scratch/wisp/probes/evidence-close/1/msg-02-sec1.txt
?? .scratch/wisp/probes/evidence-close/1/msg-03-sec2.txt
?? .scratch/wisp/probes/evidence-close/1/msg-04-sec3.txt
?? .scratch/wisp/probes/evidence-close/1/msg-05-sec4.txt
```

⇒ `docs/` 零行未提交（本腿对 `docs/**` 零写），evidence-close 下只有 `1/` 那批未跟踪件（别腿的 msg/log），本腿新建的 `6/` 在 17:05 那一遍之后才出现。

### 0.1 五枚票面文件真名（`ls .scratch/wisp/issues/<号>-*.md` 现量，全部已带 `-done`）

| 票 | 票面全路径 | 行数／字节 |
|---|---|---|
| 78 | `.scratch/wisp/issues/78-linux-vet-buildtags-unblocks-d22-gate-done.md` | 178／17,295 |
| 115 | `.scratch/wisp/issues/115-seal-notices-carry-the-resolvers-answer-while-the-cases-compare-caller-spelling-done.md` | 461／82,966 |
| 250 | `.scratch/wisp/issues/250-portable-tests-guard-c-merges-go-list-stderr-into-its-own-denominator-so-a-cold-module-cache-kills-the-whole-core-scope-reading-done.md` | 65／10,525 |
| 254 | `.scratch/wisp/issues/254-fifth-tier-winsec-no-ci-caller-done.md` | 47／7,779 |
| 267 | `.scratch/wisp/issues/267-confirm-timeout-config-has-no-bounds-so-the-c18-warning-can-vanish-done.md` | 87／24,037 |

改名（`-done`）那笔 commit 逐枚现量（`git log --diff-filter=R --all --name-status`）：
78＝`294f35d7`（09-21 12:18）／115＝`f5bbccd8`（09-22 23:09）／250＝`e39386b7`（10-02 09:27）／
254＝`21a15d27`（10-02 11:33，后又由 `46079fcc` 10-04 09:18 缩短 slug）／267＝`c6cf66e6`（10-05 10:58）。

### 0.2 本腿用到的尺（⛔ 不写进被扫文件）

- 框尺四形：A `grep -c '^- \[ \]'`（顶格未勾）／B `grep -cE '^[[:space:]]+- \[[ x]\]'`（**必须有**缩进的容缩进锚，派单原样）／
  C `grep -cE '^[[:space:]]*- \[[ x]\]'`（缩进**可选**的容缩进锚，能同时看见顶格）／D `grep -c '^- \[[ x]\]'`（顶格总框）。
- 凭据存在尺（只按**文件名**，⛔ 不读内容）：`find docs/evidence -name '<号>-*'`（＋`-iname '*<号>*'` 复量）、
  `find docs/evidence -mindepth 2 -name '*.md'`（非 s1 层）、`grep -rl '<号>-' .scratch/wisp/probes/`。
- 独立性尺：`grep -n 'agent=orchestrator' <票面>`（派单指定那把）。
- 界面签收尺：`grep -cE '签收|截图|owner 现场|肉眼|界面看' <票面>`。

### 0.3 ★负向结论前的"尺命中真名"控制（本腿逐条跑过，⛔ 不靠它下判语）

- 凭据名册尺：同一把 `find docs/evidence -name '<号>-*' -name '*.md'` 在**已知在盘**的号上给非零——
  `248→2`／`252→1`／`255→2`／`257→1`／`265→1`；本五枚给 `78→0`／`115→0`／`250→0`／`254→0`／`267→0`
  （`-iname '*<号>*'` 复量同为 0）⇒ **"名下零凭据件"是真负向，不是尺瞎**。
  同形控制另跑一组：`230→0`／`251→0`／`253→0`／`263→0`／`268→0`／`243→0`（这几枚"名下无表"是**已由别腿登记**的账，见 §4）。
- 独立性尺：`grep -c 'agent=orchestrator'` 在早期八枚上逐枚命中——`01→3`／`02→1`／`03→2`／`04→2`／`06→2`／`13→2`／`14→2`，
  全票池命中枚数＝**37 枚文件**（`grep -rl 'agent=orchestrator' .scratch/wisp/issues/ | wc -l`）⇒ 尺有牙；
  本五枚该尺**逐枚 0 命中**（读数见各枚⑤栏）。⚠ 但 0 命中**不等于**"非自裁"，理由写在 §3-6。
- 缩进框尺射程：全票池 `grep -rhE '^[[:space:]]*- \[[ x]\]'`＝**1,486** 行 vs `grep -rhE '^- \[[ x]\]'`＝**1,484** 行（差 2），
  未勾形 777 vs 776（差 1）⇒ **缩进里确实藏框**（藏在 `113-…-done.md`／`62-liquid-glass-ball-visuals.md`），
  本腿据把尺**改形**（C 形，缩进可选）逐枚复量，两形读数都记下（见各枚②栏）。

---

## §1 五枚工作单（六栏逐枚，⛔ 不合并）

### 1.1 票 78 —— 待写

### 1.2 票 115 —— 待写

### 1.3 票 250 —— 待写

### 1.4 票 254 —— 待写

### 1.5 票 267 —— 待写

---

## §2 治理账三把尺（只量读数，⛔ 不下"该不该算违规"的裁）

## §3 判得心虚的枚数与具名理由

## §4 判不动／要 owner／要搬运的清单
