# 票 171 落地腿 171-r1 实现记录（第一部分：前提与 AC#1）

- 腿：`171-r1`（写腿·实现者）。写面天花板＝`.scratch/wisp/probes/171/**` ＋
  `.scratch/wisp/probes/161/**`（若洞在 `attrib.sh` 本体）＋ `.github/workflows/ci.yml` 的**新增步骤**
  ＋ 本件 `docs/evidence/s1/171-r1-impl*.md`。
- 分工硬规矩：本件**不判自己过**、不给任何 AC 翻框、不碰票面 `Status:`／`-done`、不改台账。
- 取数时刻：`2026-10-21 09:1x–11:0x +08`。开工锚 `b482968f`；取数期间 HEAD 前移到 `8288d49a`。
- 凡写"响／不响"的，一律是当次现跑的读数；凡推的、量不到的，逐条写在第二部分 SEC10。
- 本件分三枚：本件＝SEC0-SEC4（前提、AC#1）；`171-r1-impl-part2.md`＝SEC5-SEC6（AC#3、AC#9）；`171-r1-impl-part3.md`＝SEC7-SEC11（门禁、写面、码位、边界）。

---

## 0. 开工前三步

### 0.1 写面临场（`git status --porcelain` 限三条写面）

```
?? .scratch/wisp/probes/161/r2/__pycache__/
?? .scratch/wisp/probes/161/r2/ctl/
?? .scratch/wisp/probes/161/r5/negative-control/
?? .scratch/wisp/probes/161/r6/logs/flip-7.txt
```

四条全部 `??`（未跟踪）且**全部早于本腿**（`probes/171` 开工时为空）。开工那一刻**没有别人的未提交改动**，
不构成停手条件。（取数中途出现并发改动，见 SEC5.0 与第二部分 SEC5.0。）

### 0.2 票面前提新鲜度（尺 v4 的(8)）：逐条现跑

| 票面／派单断言 | 现跑命令 | 今天的读数 | 判定 |
|---|---|---|---|
| `attrib.sh --tracked-only` 今天 rc=1、`tracked-dirty=35` | `sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only` | **rc=1**，末行 `attrib.sh: RED (tracked-only) - tracked-dirty=35 files_dirty=35.` | **未过期** |
| `attrib.sh --self-test` 应 `cases=9 failures=0` | 同上一枚加 `--self-test` | **rc=0**，`attrib.sh: --self-test cases=9 failures=0` | **未过期** |
| CI 跑的那行＝`ci.yml:240` 的 `--tracked-only` | `Select-String -Path .github/workflows/ci.yml -Pattern 'attrib\.sh'` | `240: run: sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only`（逐字） | **未过期** |
| 仓里**没有** `scripts/attrib.sh` | `Test-Path scripts/attrib.sh` | `False`（不存在） | **未过期**（派单提醒成立） |
| AC#3＝"自检不在 CI" | 同第二枚：`attrib.sh` 在 `ci.yml` 里只出现 1 次（`:240`，且是 `--tracked-only`）；字串 `--self-test` 出现 **0** 次 | 成立 | **未过期** |

### 0.3 两条**过期**（派单说"任一条过期即停手上报"，已当场上报）

**过期甲：AC#1（尾随注释）与 AC#9（表被核到几分）的码早已落地，不在"未闭"状态。**

- `171-r3`（2026-10-03，提交 `f4d64fcc`／`2a5489d3`／`7cb9d570`）已把 `probes/154/gate-clauses.sh` 的
  分子换成 `call_lines()`（`strip_comment()` ＋ `comment_only_lines()`），**pair 与 run 两族一次改齐**；
  并把受影响的 `want_n` 基线在同一格复算（G2 `2->1`、G7 `3->1`、G7neg `4->2`）。
- 同一程把 AC#9 第 1 项「核那两列」做进 `probes/161/r6/flip-declaration.sh:202`-`:223`
  （`want_row="基线=${new_n}枚 实测=${m0}枚 差=$((new_n-m0))枚"`，`grep -qF` 核）。
- AC#9 第 2 项「删样本不删腿」正面台件由 `171s`（10-14，随票 178 的 piles 口径）重推。
- **票面 `Progress log` 里没有 `171-r3`／`171s` 两行**（现跑 `Select-String` 0 命中）⇒ 票面读到的是 09-27 的旧状态。
- 我今天现跑三枚台件，**全部 GREEN**：
  - `sh .scratch/wisp/probes/171/r3/r3-trailing-comment.sh post` ⇒ **rc=0**
  - `sh .scratch/wisp/probes/171/r3/r3-sample-deleted.sh` ⇒ **rc=0**
  - 列核在盘上：`Select-String ... flip-declaration.sh -Pattern 'the row prints its own numbers'` ⇒ **3 处**

处置：本腿**没有**重造 AC#1／AC#9 的落地码（重做已落地的活会与 `171-r3` 打架）。本腿做的是派单点名的
**两发对照台件**（我自己这一版），以及 AC#9 第 1 项的**两向对照读数**（含一枚我自造的假红，见第二部分 SEC6.3）。

**过期乙：AC#3 的前提成立，但票面给的反向判据今天答不出——这才是真正还开着的那一格。**

- 票面 AC#3 反向判据逐字：*"把自检里『必须响』的那枚样本摘掉 ⇒ CI 那一步必须红"*。
- 我现造 `probes/171/r1/attrib-mut-noring.sh`（只删 `attrib.sh:291` 那一行
  `st_case 1 UNATTRIBUTABLE '.scratch/wisp/probes/777/second-fake-ticket.go'`，`diff` 逐字 1 行），
  跑 `--self-test` ⇒ **`cases=8 failures=0`、`SELF-TEST GREEN`、rc=0**。
- ⇒ **只加 CI 步骤答不出这条判据**：`--self-test` 只数"跑了几条"，对"名册该有几条"没有任何断言。
  这正是票面那句 *"答不出＝这一格白接"* 的形状。

---

## 1. AC#1（尺洞二＝尾随注释）：两发对照台件

### 1.1 台件与手法

台件＝`.scratch/wisp/probes/171/r1/r1-two-shots.sh`（＋`fixtures/t1..t3*.go.txt`）。
手法照 `probes/171/r3/r3-trailing-comment.sh` 的形状（丢弃树在 `$TMPDIR`，仓内零删除、零 checkout），
但台件是**我自己这一版**：派单要的两发对照 ＋ 一发我加的反控。

| 发 | 形状 | 判据 |
|---|---|---|
| **T1** | 词根**只**出现在一行真代码的**行尾注释**里，全文件没有开方调用 | 必须**不计入**（派单"只有同行注释"那一发） |
| **T2** | **真开方调用** ＋ 合方**只写在同一行的行尾注释**里 | 未修尺读成"已合"⇒**静默的活漏**；修完必须**点名**（派单"注释＋真违规"那一发＝正控） |
| **T3** | 真开方调用，**通篇没有合方**（代码里没有、注释里也没有） | 两把尺**都必须计入** |

**T3 是我加的反控，理由写进台件头**：没有它，"T1 从 1 变成 0"与"这把尺什么都看不见了"两种情形读数
**一模一样**——那是恒真的另一面。T3 把那个方向堵上。

**两把尺**：未修尺＝`gate-unfixed-2a5489d3-parent.sh`（`git show 2a5489d3^:` 的字节，577 行）；
修后尺＝工作副本（735 行）。取 `2a5489d3^` 而非 `HEAD`：`171-r3` 一提交，`HEAD` 就不再是未修尺，
整个对照会**悄悄变成一把尺自己跟自己比**。台件里有 `cmp` 守卫：两尺字节相同即 rc=2 不给读数。

### 1.2 读数（`sh r1-two-shots.sh both` ⇒ **rc=0 GREEN**，199.1s）

逐字取自 `logs/gate-{unfixed,fixed}-T{1,2,3}.txt`：

| 发 | pair `G6 未成对枚数` 未修 -> 修后 | run `G2 命中行数` 未修 -> 修后 | 判 |
|---|---|---|---|
| T1 | **1 -> 0** | **1 -> 0** | 未修尺把行尾注释当调用点（洞是真的）；修完不计入 |
| T2 | **0 -> 1** | **3 -> 1** | 未修尺静默了一枚真漏；修完点名＝正控响了 |
| T3 | **1 -> 1** | **1 -> 1** | 两把尺都计入＝修完没有把尺变瞎 |

被判红的行，**逐字**：

```
#   UNPAIRED internal/fake171/t1_only_trailing_comment.go
#   UNPAIRED internal/fake171/t2_comment_plus_real_violation.go
#   UNPAIRED internal/fake171/t3_real_violation_no_mention.go
#   COMMENT-ONLY ...:internal/fake171/t1_only_trailing_comment.go:23:	cfg.Name = "retention-job" // scheduled through OpenTask by the caller, never here
```

第一行＝T1 未修尺点名（洞）；第二行＝T2 修后尺点名（正控响）；第三行＝T3 两尺都点名；
第四行＝修后尺把 T1 判为"不计入"的**理由本体**，逐枚打出来给复核，不藏在算式里。

### 1.3 射程没动（我自己的控制发）

台件尾部比两把尺的 pathspec 行：**unfixed=10 行、fixed=10 行、两向逐字节相同**；
另有两条反控：一、两尺全文**必须**不同（577 对 735 行，实测不同）；
二、那味 grep **必须**打到东西（10 行），否则"相同"是空口恒真。

---

## 2. 本腿自己踩的两枚仪器坑（记账，不抹掉）

**坑一：gofumpt 不在 shell 的 PATH 上。** 第一版台件读 `command -v gofumpt` 为空，于是打
`gofumpt not on PATH; sample formatting left unchecked` ——那是**具名跳过，仍是跳过**（`171-r3` SEC7 同形）。
改法：照 `attrib.sh:139`-`:142` 的取法，从 `go env GOPATH` 推 `bin/gofumpt.exe`，推不到就判 **FAIL**
而不是打一行字。

**坑二：丢弃树里的 CRLF 让每个样本都被判"不干净"。** git 在丢弃树里把样本重写成 CRLF
（`warning: LF will be replaced by CRLF`），`gofumpt -l` 于是因行尾把三枚样本全列出来。
那不是样本的毛病，是丢弃树的毛病。改法：在一份 `tr -d '\r'` 的**副本**上读格式，并把这件事写进
`ok` 行原文；同时补一枚真因：三枚 fixture 原本**缺尾随换行**（`\ No newline at end of file`），
补上后 LF 副本 `gofumpt -l` 列 **0** 枚（`v0.12.0 (go1.27.1)`）。

**坑三（第二部分 SEC6.3 详述）：AC#9 台件第一版把门禁的 baseline 格砍掉**，导致"诚实尺也红"，
那红是我自造的，不是门的病。口径：**"诚实尺也红"先怀疑台件**，与票面"零命中先怀疑尺"同一条。
