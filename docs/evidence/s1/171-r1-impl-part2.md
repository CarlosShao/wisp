# 票 171 落地腿 171-r1 实现记录（第二部分：AC#9、AC#3、门禁、写面）

接 `171-r1-impl.md`。本件＝SEC5–SEC11。凡"响／不响"都是当次现跑；量不到的在 SEC10 逐条具名。

---

## 5. AC#3（尺洞一＝接 CI）

### 5.0 停手上报：这一格在 2026-10-21 10:24 由**另一个程序**并发落地

派单开工第 1 条：*"`probes/171/**`、`probes/161/**` 或 `ci.yml` 有**别人的**未提交改动 ⇒ 停手上报。"*
该条件在取数期间成立，**本腿据此停手**（不覆盖、不还原、不认领）：

```
$ git status --porcelain -- .scratch/wisp/probes/161 .github/workflows/ci.yml
 M .scratch/wisp/probes/161/r5/attrib.sh     (+114/-4，552 -> 662 行，mtime 2026-10-21 10:24:08)
 M .github/workflows/ci.yml                  (+48/-0，md5 7CF2BFB3... -> 8DB56645...)
```

两份改动的注释里**逐字自称本腿**（`ticket 171 AC#3`／`leg 171-r1`），但**不是本腿写的**：
本腿对 `ci.yml` 一字节未碰（开工时 md5 `7CF2BFB3821A9A1F2E764A5DB3138C4E`，40 枚步骤名快照在
`D:\tmp\wisp\20261021\171r1\ci-steps-before.txt`），对 `attrib.sh` 只做过 `probes/171/r1/` 下的**副本**变异。
另：HEAD 在取数期间由 `b482968f` 前移到 `8288d49a`。

### 5.1 并发那版改动为什么**不会**让既有 job 变松（逐条，读它的字节得出的）

1. `git diff --stat` 逐字 `1 file changed, 48 insertions(+)` ⇒ **纯新增**，既有步骤零字节改动。
2. **没有** `continue-on-error`（本仓已裁定那是零执行命中形，`ci.yml` 头 `D22` run-away mode 6）。
3. `if: ${{ !cancelled() }}` 是该 job 里 `staticcheck`／`mockllm module vet`／`portable-tests-selftest`
   已经在用的同一枚 guard——它改的是"这一步**报不报**"，不是"这一步**能不能红**"。
4. 位置＝lint job 末尾的下一个空位（`portable-tests-selftest` 之后），**没有越过**任何既有步骤。
5. 它不读 Go、不建任何包；唯一真跑 gofumpt 的那枚用例在 `mktemp` 的**仓外**丢弃树里。

新增那步**逐字**（含缩进）：

```
      - name: "Attribution instrument self-test with its roster nails (ticket 171 AC#3 - the mutation goes red)"
        if: ${{ !cancelled() }}
        run: sh .scratch/wisp/probes/161/r5/attrib.sh --self-test
```

### 5.2 反控发：摘掉"必须响"的样本 ⇒ 那一步必须红

| 尺 | 命令 | 读数 |
|---|---|---|
| **并发版（现工作副本）** | `sh .scratch/wisp/probes/161/r5/attrib.sh --self-test` | **`cases=10 failures=0`**、GREEN，新增 `PASS  want ring=5 silent=3 named=3` |
| **同一把尺，删掉 `attrib.sh:291` 那一行** | `sh probes/171/r1/attrib-mut2-noring.sh --self-test` | **`cases=9 failures=1`、SELF-TEST RED**；红句逐字见下 |
| **改前（HEAD 那版 552 行），同一发** | `sh probes/171/r1/attrib-mut-noring.sh --self-test` | **`cases=8 failures=0`、GREEN、rc=0** ← 这就是洞 |

红句逐字：

```
   FAIL  roster assertion  want ring=5 silent=3 named=2  got ring=4 silent=3
         why (every reason that held, one per line):
           ring-count: registered must-RING cases=4, this file declares 5 (a must-ring sample was deleted or added without moving ST_RING_EXPECT)
```

⇒ 票面 AC#3 的反向判据 *"摘掉必须响的样本 ⇒ 那一步必须红"* **现在答得出**，而**改前答不出**。
两枚变异台件都在 `probes/171/r1/`，各自只删那一行，`diff` 均 1 行。

### 5.3 本腿对 AC#3 的立场

**本腿不认领 `ci.yml` 与 `attrib.sh` 这两处改动。** 理由不是谦虚，是分工硬规矩：实现者给自己那格
出终裁＝倒影。这两处的落地者具名登记为"并发腿（自称 171-r1）"；本腿只交**独立复验的读数**（SEC5.2）。

---

## 6. AC#9（表到底被核到几分）

### 6.1 现量：那一格今天核到什么程度

`probes/161/r6/flip-declaration.sh`（245 行）在 `:202`-`:223` 有列核：

```
want_row="基线=${new_n}枚 实测=${m0}枚 差=$((new_n - m0))枚"
row=$(grep -m1 "^# STALE 腿=$stale_leg " "$logdir/flip-$seq.txt" || true)
if [ -n "$row" ] && printf '%s' "$row" | grep -qF "$want_row"; then ... ok
else ... FAIL  the row prints its own numbers - row got [...] want [...]
```

`Select-String -Pattern 'the row prints its own numbers'` ⇒ **3 处**；改前（`2a5489d3^`）同一命令 ⇒ **0 处**。

### 6.2 两向对照（我这一版）

台件＝`probes/171/r1/r1-flip-columns.sh`。造假尺＝`fixtures/gate-lie-columns-r1.sh`
（只改 STALE 打印那一行：`实测=${n}枚 差=$((COUNT-n))枚` 篡成 `实测=${COUNT}枚 差=0枚`；`diff` 行数核过＝2）。

**为什么只跑 (stale) 那一格、不跑整扇门**：实测一轮完整尺 = **7m20s**，整门 11 轮约 80 分钟，三扇门约 4 小时；
而 AC#9 第 1 项**完全由 (stale) 那一格决定**，它只跑一轮尺。所以台件把门**原件**（不是手抄）裁到
"baseline 格＋stale 格"，每份副本只 patch 两行（`gate_rel=`／`anchor=`，两行都核"补丁必须落得下来"）。

读数（`sh r1-flip-columns.sh both` ⇒ **rc=0 GREEN**，1781.4s）：

| 发 | 门 | 尺 | `row-claim` | `columns-claim` | 判 |
|---|---|---|---|---|---|
| (1) | 未修门 | 造假尺 | GREEN | **NO-CLAIM** | (3)-limit 复现：未修格对那两列**没有任何断言** |
| (2) | 修后门 | **同一造假尺** | GREEN | **RED** | **AC#9 第 1 项的判据成立** |
| (3) | 修后门 | **诚实尺** | GREEN | **GREEN** | 反恒红：那一味不是"永远响" |

(2) 的红句逐字：

```
  FAIL  the row prints its own numbers - row got [# STALE 腿=G5neg 声明=ring 基线=9枚 实测=9枚 差=0枚 处置=把 want_n 核下来；样本一枚没变好而是被摘掉的，先回答谁摘的] want [基线=9枚 实测=8枚 差=1枚]
```

(1) 里那枚 `ok    the stale baseline is its own readable table row: ... 实测=9枚 差=0枚 ...` 是**关键**：
行在、腿名对、枚数 1 都对，**唯一错的就是那两列**，而未修格照样 `ok`。
(3) 里同一格对诚实行 `实测=8枚 差=1枚` 判 `ok`。

(4) 门对照：两门 **198 对 245 行**（必须不同）；列核字串 unfixed=**0** 处、fixed=**3** 处（必须只出现在修后门）。

### 6.3 我自造的假红（记账，不抹掉）

第一版台件把 baseline 格也砍了，读数：

```
(3) FIXED door + the HONEST row
  FAIL  the row prints its own numbers - cannot read G5neg measured from the baseline run, so the columns cannot be checked
  FAIL  (3) expected GREEN ... got RED - the cell became constant-red, which A320/A321 retired
```

那一发**不是门的病**：修后门禁的 `实测` 是从 `$logdir/flip-baseline.txt` 读的（baseline 格写它），
砍掉 baseline 格 ⇒ 那一格**答不出**、诚实尺也红。修法＝保留 baseline 格，并加一枚
`grep -q '== (0) baseline'` 断言（砍掉就 rc=2 而不是红），另加 `rm -f flip-baseline.txt` 让三扇门
不互相偷对方的读数。口径：**"诚实尺也红"必须先怀疑台件**。

### 6.4 AC#9 第 2 项「删样本不删腿」

已由 `171s` 落地并**今天现跑 GREEN**：`sh probes/171/r3/r3-sample-deleted.sh` ⇒ **rc=0**，四形逐形有读数：

- D2 诚实态 ⇒ `# ok 腿=G6 声明=ring 基线=2枚 实测=2枚`，表内 0 枚
- D1 删一枚 ⇒ `# SHR 腿=G6 ... 基线=2枚 实测=1枚`，单列表 1 枚，**不进退码**
- D1c 多一枚【只关】⇒ 不进实测那一列（票 178 的四堆口径）
- D0 删光 ⇒ `# BAD 腿=G6 ... 因=空/非空那一味与声明不符`，退码跟着动

⇒ 票面那句"造一发台件，并答未修码上响不响"的答案：**删一枚＝不响（SHR，只在表里）；删光＝响**。

---

