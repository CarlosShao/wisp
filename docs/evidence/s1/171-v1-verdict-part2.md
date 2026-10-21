# 171-v1 裁决 · part2 —— `AC#1` 与 `AC#3` 的逐发逐字读数

本件是 `171-v1-verdict.md` 第 2 节（`AC#1`）与第 3 节（`AC#3`）的**逐字展开**。判语在主件；这里放可复核的每一行。

---

## sec 2 `AC#1`：尾随注释——三发两向逐形

**尺与锚**：`sh .scratch/wisp/probes/171/r1/r1-two-shots.sh both` => **rc=0**、**SECONDS=336**、末行逐字
`r1-two-shots.sh: GREEN - mode=both, every assertion above held.`
改前尺＝`probes/171/r1/gate-unfixed-2a5489d3-parent.sh`（577 行，blob `a5bcc802`）；
改后尺＝`probes/154/gate-clauses.sh`（735 行，blob `22a7fd6b`）。

### 2.1 T1（只有行尾注释：**必须不计**）

```
===== T1: word roots ONLY in a trailing comment (must NOT be counted) =====
  ok    T1 unfixed: the trailing-comment file IS counted (未成对=1 枚) <- the hole is real today
  ok    T1 unfixed: the trailing comment is a run hit (命中行数=1) <- same hole, run leg
  ok    T1 fixed: the trailing-comment file is NOT counted (未成对=0)
  ok    T1 fixed: the trailing comment is NOT a run hit (命中行数=0)
```

 =>  **改前 1／1，改后 0／0**。两条腿（pair 与 run）在同一样本上**同时**从"计入"变"不计入"。

### 2.2 T2（注释＋真违规：**正控必须响**）

```
===== T2: trailing comment PLUS a real violation (the positive control must RING) =====
  pair G6 未成对枚数：unfixed=0   fixed=1
  run  G2 命中行数 ：unfixed=3   fixed=1
  pair G6 点名名册（fixed）  ：
     #   UNPAIRED internal/fake171/t2_comment_plus_real_violation.go
  ok    T2 unfixed: the leak is SILENT (未成对=0 - the comment bought it a close) <- the hole, dangerous direction
  ok    T2 unfixed: the run leg already counts the line (命中行数=3)
  ok    T2 fixed: the leak is NAMED (未成对=1) <- the positive control rang
  ok    T2 fixed: the run numerator is 1 (the call), not 2
```

 =>  两个方向都对：`pair` **0 -> 1**（漏被点名＝正控响）、`run` **3 -> 1**（注释那两行不再算命中，只剩真调用）。
**两条同时成立**才叫"分子只吃代码那一半"。

### 2.3 T3（腿 A 自加的反控：真违规、全无提及）

```
===== T3 (this leg's control): a real violation and NO close mention at all - both rulers must count it =====
  (without T3, 'T1 went 1 -> 0' would be indistinguishable from 'this ruler stopped seeing anything')
  pair G6 未成对枚数：unfixed=1   fixed=1
  run  G2 命中行数 ：unfixed=1   fixed=1
  ok    T3 unfixed: the bare violation is counted (未成对=1)
  ok    T3 fixed: the bare violation is STILL counted (未成对=1) <- the fix did not blind the ruler
```

 =>  **两把尺都是 1**。这是整套对照的**承重件**：没有它，"T1 由 1 变 0"与"这把尺从此什么都看不见"
**读数完全一样**。T3 成立 => T1 的 0 是"剔掉了注释"，不是"瞎了"。

### 2.4 射程可比性（本腿独立复核，NOT 未采信自报）

| 检查 | 读数 |
|---|---|
| 改前快照是否＝`2a5489d3^` 字节 | **是**：`git hash-object` ＝ `git rev-parse 2a5489d3^:…gate-clauses.sh` ＝ `a5bcc802` |
| pathspec 承载行数 | 改前 **10** ／ 改后 **10** |
| 两向 `comm -3` | **全空**（逐行打出 10 对 10，逐字节同） |
| 正控：两尺整体文本必须不同 | **成立**（577 vs 735 行；`--numstat` ＝ `192 34`） |

**正控为什么必要**：没有它，"pathspec 相同"可能是"两边都 grep 不到"的假绿。

### 2.5 码侧（自报之外）

`2a5489d3` 对 `gate-clauses.sh` **新增**：`strip_comment()`（按"**空白＋`//`**"切，不是按第一个 `//` 切
 =>  `https://x` 不被切）、`comment_only_lines()`（把**被剔除的行逐枚打出来**，`#   COMMENT-ONLY …`）、
并把 `pair()` 与 `run()` **两族分子一次改齐**＝票面 `AC#1` 那格 09-27 19:1x 追加提醒逐字要求的事。

---

## sec 3 `AC#3`：接 CI——四项逐条

### 3.1 新增那步逐字（`HEAD:.github/workflows/ci.yml`，`:441-443`）

```yaml
      - name: "Attribution instrument self-test with its roster nails (ticket 171 AC#3 - the mutation goes red)"
        if: ${{ !cancelled() }}
        run: sh .scratch/wisp/probes/161/r5/attrib.sh --self-test
```

`:397-440` 是 **44 行纯注释**（为什么加这步／为什么不加 `continue-on-error`／为什么 `if:` 用这个／位置）。

### 3.2 纯新增：删除行 ＝ **0**

```
@@ -396,0 +397,48 @@ jobs:
+ ... 48 added lines ...
```

尺：`git -c core.quotePath=false diff -U0 85edab3d^..HEAD -- .github/workflows/ci.yml | grep '^-' | grep -v '^---'`
 =>  **0 行**。全窗口只有 `85edab3d` 一笔碰过 `ci.yml`。

### 3.3 `continue-on-error`：**解析 YAML 查键**（ NOT 用 grep 数注释）

| 尺 | 改前（`85edab3d^`） | 改后（HEAD） |
|---|---|---|
| **PyYAML 解析出的键数** | **0** | **0** |
| 文本出现次数（含注释） | 17 | 20 |

**文本 17 -> 20 是事实**：HEAD 里 20 处逐行 `TrimStart` 后第一个字符**都是 `#`** => 全是注释行。
**解析出的键改前改后都是 0** => 新步**没有** `continue-on-error`，全文件**任何**一步也没有。
 =>  "17 涨到 20"与"有三个键被加上"是**两件事**；本腿分开报，不把文本增量读成行为变化。
`if:` ＝ `${{ !cancelled() }}`（**不是** `success()`；与同 job 既有 staticcheck／mockllm vet 同形）。

### 3.4 位置：**只在 `lint` job 末尾**

| job | 改前 | 改后 | 差 |
|---|---|---|---|
| **`lint`** | 13 | **14** | **+1（就是这一步，index=13/14＝最后一步）** |
| `lint-frontend`／`slo-full`／`slo-smoke`／`test-core`／`test-windows` | 11／5／6／7／10 | 同左 | **0（逐 step 四元组全同）** |

job 名单零增删 => **触发表与并发组未动**（票 134 地界未被碰）。

### 3.5 名册断言的牙：本腿自跑四发

| 发 | 尺 | 变异 | rc | 读数 |
|---|---|---|---|---|
| (1) 基线 | 仓内件 | 无 | **0** | `cases=10 failures=0`（23 秒） |
| (2) **本腿删行** | 修后尺副本（662 行，`3b3af251`） | 删 777 那行 | **1** | `cases=9 failures=1` |
| (3) **抗装饰** | **改前尺**（`8288d49a`，552 行，`c6e11f0c`） | **同删法** | **0** | `cases=8 failures=0` |
| (4) 未变异对照 | 修后尺副本 | 无 | **0** | `cases=10 failures=0` |

删行逐字（两份副本都验过目标行已消失）：
`before: 345: st_case 1 UNATTRIBUTABLE '.scratch/wisp/probes/777/second-fake-ticket.go'`／
`m1 after deletion: lines=661 target-line occurrences=0`／`m2 after deletion: lines=551 target-line occurrences=0`

**(2) 的完整 FAIL 块**（"逐字点名"的证据）：

```
   FAIL  roster assertion  want ring=5 silent=3 named=3  got ring=4 silent=3
         why (every reason that held, one per line):
           ring-count: registered must-RING cases=4, this file declares 5 (...)
           named-missing: a NAMED must-ring sample is not among the cases that ran (...)
         missing must-ring sample(s), verbatim:
           .scratch/wisp/probes/777/second-fake-ticket.go
attrib.sh: --self-test cases=9 failures=1
```

- 三枚钉子（`ring` 数／`silent` 数／**逐名在场**）**各自独立可红**，**不是 `elif` 链**（源码 `attrib.sh:465-471` 逐字说了为什么）。
- **(3) 是本格抗装饰的核心**：改前尺同一发删行 **rc=0** => **洞在改前真开**，这一步不是装饰。
- 副本位置 `D:\tmp\wisp\20261021\171v1\mut\{m1,m2,ctl}\`；**仓内件一枚未改**。副本用尺自带
  `ATTRIB_ROOT` 接缝（`attrib.sh:137`）指回真仓 —— 不加这个会 rc=2（"no verdict"）。

### 3.6 反向判据：CI 那一步**会红**

CI 那步命令**逐字**就是 `sh .scratch/wisp/probes/161/r5/attrib.sh --self-test`，所以把那命令打在副本上**就是那一步本身**：

| 形状 | CI 步 rc | 读数 |
|---|---|---|
| **只接线、未加钉子**（改前尺 551 行＋删行） | **0** | `cases=8 failures=0` => 接线本身**沉默** |
| **接线＋三枚钉子**（修后尺 661 行＋删行） | **1** | `cases=9 failures=1`，逐字点名 `…/probes/777/second-fake-ticket.go` |

 =>  若只做接线、不做名册钉子，**这一格就是装饰**；本腿交付的形状**会红**。**判 `AC#3`** 〔成立〕。
 NOTE  一条本腿**未**证（ NOT 当"过了"）：那步在 **GitHub 上的真实颜色未观测** —— 本腿证的是"这条命令会因删样本而 rc=1"，
**没有**证"CI 跑过"。按本仓既有账口径这属"未观测"，**不阻塞**本格，但也不能记成"CI 已验证"。
