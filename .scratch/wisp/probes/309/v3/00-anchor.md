# 309-v3 起手锚件（非实现者终裁腿，⛔ 产码／⛔ 改断言／⛔ 翻框）

## 0. 锚（本腿起手现跑，⛔ 抄派单）

尺＝`git rev-parse --short HEAD`＝**`c27928d6`**
派单写的锚＝**`04f2ecbc`**（`git rev-list --count 04f2ecbc..HEAD` 见下）⇒ **两号不同**，本腿一律以 `c27928d6` 为锚；
凡引用「哪一笔改了什么」的读数，仍按派单点名的历史枚（`df8d517d`／`a1f84d39`）具名打，那些是**历史对象**、⛔ 随 HEAD 漂。
`04f2ecbc`⇄`df8d517d`⇄`a1f84d39`⇄`6fed247e`⇄`715f63cb` 六枚对 `c27928d6` 全部 `git merge-base --is-ancestor … HEAD`＝**真**（件 `logs/gitlog.txt`，另 `logs/ancestry.txt`）。

起手台面（件 `logs/anchor-state.txt`）：
- `git diff --cached --name-only | wc -l` ＝ **0**
- `git status --porcelain -- internal cmd tools scripts | wc -l` ＝ **0**
- 票面框数尺＝`grep -cE '^[[:space:]]*- \[x\]'`＝**4** ⇄ `- \[ \]`＝**2**（与派单 §0「4 勾／2 未勾」对得上；未勾＝`AC#2`／`AC#3`）

## 1. 射程与台面

- 本腿射程＝`internal/tools`（派单 §2）；⛔ 跑整包 `cmd/wisp`；与 `304-v1`（`internal/llm`／`internal/proc`）互⛔ 遮蔽。
- 突变＝**仓外成对导出树**，模板六件套齐（`go.mod go.sum internal cmd tools scripts`），每棵独立 `mktemp -d`，先数 `.go` 证非空再取数。
- ⛔ 引用编排者 §3 的读数当判语：那三发本腿**自己重打**（`10-teeth-pairs.md`）。
- ⛔ push／⛔ 翻框／⛔ 动 `Status:`／⛔ 改任何断言；还原＝`mv` 挪开⛔ 删＋`cmp`⇄`sha256sum` 双证。

## 2. 交付名册（计划）

1. `00-anchor.md`（本件）
2. `10-teeth-pairs.md`＝`AC#2` 三对（new／make([]T,n)／转换形）逐对「红句原文是哪一行」＋pristine 绿控＋还原双证
3. `20-cells.md`＝两格判语＋必答 (I)(II)＋4 枚 0 字节件的处理结论
4. `docs/evidence/s1/309-v3-verdicts.md`＝裁决表本体
5. 票尾一条 bullet（1 加／0 删）

## 3. 写面声明（三件套）

① `.scratch/wisp/probes/309/v3/**`；② `docs/evidence/s1/**`（AGENTS §1.5，⛔ 当越界）；③ 本票末尾一条 append-only bullet。
件＝`.md`／`.txt`；⛔ `.out`；⛔ 在 `probes/**` 建 `.go`；每把门禁件自落一行 `rc=N`；⛔ 0 字节件当读数。
