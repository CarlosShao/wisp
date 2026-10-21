# 171-v1 裁决 · part3 —— `AC#9`／门禁／日志完整性 的逐字读数与自证

本件是 `171-v1-verdict.md` 第 4／5／6／7／10 节的**逐字展开**。判语在主件。

---

## sec 4 `AC#9`：两枚子项的今日读数

### 4.0 先说不可得的那两发（ NOT 把没跑的当跑了）

派单记载两腿都报"AC#9 两向对照未取到今日读数（单发 80 分钟～4 小时，都 kill 了）"。本腿复核并**自己也撞上同一堵墙**：

| 台件 | 本腿现跑结果 | 处置 |
|---|---|---|
| `sh probes/161/r6/flip-declaration.sh`（整份门） | **rc=124 @ 900 秒**（`timeout` 打掉），只走完 **2/10** 格 | **缺读数** |
| `sh probes/171/r3/r3-flip-columns.sh post`（整份） | 超 **1500 秒**被叫停（step(1) 完整，step(2)/(3) 未到） | **缺读数** |

 =>  本腿**换尺**：把 `AC#9` 真正要的那两枚判据各用**能在预算内产出读数**的定向发打。

### 4.1 `AC#9-1`（"实测／差"两列必须被核）—— **成立**

**判据原文**：让新格连那两列一起 grep；判据是"**那一列写错时新格必须红**"，且"别把 (1) 做成『又多印一行字』"。

**做法**：把门上那段断言**逐字抽出**（`sed -n '202,223p'`），喂它两种行；输入仍是门自己读的两个文件
（`flip-baseline.txt` 学 `实测`、`flip-$seq.txt` 读那一行） =>  测的是**在册逻辑**。先证抽出来的是在册字节：
门 blob `4a79fb54` ＝ 工作副本 `git hash-object` **同 sha**；在册门里 `want_row=` 构造 **1** 次。关键是这三行：
`want_row="基线=${new_n}枚 实测=${m0}枚 差=$((new_n - m0))枚"` ＋ 读那一行 ＋ `grep -qF "$want_row"`。

**两发读数**：

| 发 | 喂给门那一行的 `实测／差` | 门那一格输出 | 判 |
|---|---|---|---|
| (i) **诚实行** | `基线=9枚 实测=8枚 差=1枚` | `ok    the row prints its own numbers: # STALE 腿=G5neg … 实测=8枚 差=1枚 …` | **GREEN** |
| (ii) **造假行**（171-v2 的 (3)-limit） | `基线=9枚 实测=9枚 差=0枚` | `FAIL  the row prints its own numbers - row got [… 实测=9枚 差=0枚 …] want [基线=9枚 实测=8枚 差=1枚]` | **RED** |

- (ii) 的 `实测` 从**基线跑**读出真值（8） =>  门把"9/9/0"判假 => **核的是列的一致性，不是"表有没有被打印出来"**。
- (i) 仍 GREEN => 新格**不是恒红形状**（`A320`／`A321` 那扇门没被开回来）。
- **与在册读数同向同形**：`r3/logs/harness-flip-columns-pre.txt` 是 `columns-claim=NO-CLAIM`、`-post.txt` 是
  `columns-claim=RED`（在册那发锚 `4f2a777b`，10-03；**本腿的是今天的**）。

### 4.2 `AC#9-2`（"删样本不删腿"正面台件）—— **成立，且今天真跑通**

**先更正票面一句已过期的登记**：票面第 91 行（`2026-10-14`）写票 178 改形后本票的 `r3/r3-sample-deleted.sh`
**钉 `want G6` 形 => rc=2**（逐字 `no 'want G6' declaration - no verdict`）。**本腿现跑 => rc=0、201 秒、GREEN。**
成因：`171s`（`0b638c5e`，10-14，52 枚）已把这枚钉按 178 的**四堆口径**重推，其 `s/00-anchor-and-gates.md:5`
逐字写 "rc 2->0（D2／D1／D1c／D0 四态全绿）"；现读 `gate-clauses.sh:594`＝`want G6 quiet piles`、`:595`＝`want_n 0`
 =>  钉**找得到**。而**票面没有 `171s` 这一行** => 与主件 sec 1a 同一处腐坏。

**四形读数（逐字摘）**：

```
===== D2: both registered samples present (the honest state) =====
   G6 row : # ok   腿=G6 声明=ring 基线=2枚 实测=2枚      表内枚数=0  聚合退码=6
===== D1: ONE sample deleted, the leg and its declaration untouched =====
   G6 row : # SHR  腿=G6 声明=ring 基线=2枚 实测=1枚 注=读数变好了，基线过期…
   STALE  : # STALE 腿=G6 声明=ring 基线=2枚 实测=1枚 差=1枚 处置=把 want_n 核下来…
   表内枚数=1  聚合退码=5
  ANSWER  D1 不响退码（SHR）：腿 G6 一枚样本被摘掉（(2) 堆 2 枚 -> 1 枚）之后不在 BAD 名册里…
===== D1c: one opening sample + one CLOSE-ONLY sample ((3) 只关 must NOT reach the measured column) =====
   G6 row : # SHR  腿=G6 声明=ring 基线=2枚 实测=1枚 …
  ok    D1c: the close-only file did not enter the measured column (票 178：(3) 只关不计进这一列)
===== D0: the LAST sample deleted (leg still declared, still run) =====
   G6 row : # BAD  腿=G6 声明=ring 基线=2枚 实测=0枚 因=空/非空那一味与声明不符（票 161 AC#6(1)）
   BAD 名册：G2 G5 G5neg G6 G7 G7neg
  ANSWER  D0 响（BAD 空/非空那一味）：腿 G6 的样本被删空（腿还在跑、声明还在） =>  它进了 BAD 名册…
```

 =>  票面 `AC#9(2)` 要的那句"造一发台件，并答未修码上响不响"，**四形都有读数**，且**答案分形**：
**D1／D1c 不响退码**（SHR 按 `A320`／`A321` **刻意**不进，只在【基线过期】表里看得见）；**D0 响**。

### 4.3 残余：**NOT 不记"已收"，移交**

`AC#9(2)` 还剩一形：**删掉一枚样本、而剩下的枚数恰好等于 `want_n`** => 门说 `ok`、表里也没有它（枚数基线这一形在
**数学上**就看不见它）。`r3/impl.md:417-421` 按**发现**登记、NOT 不记"已收"；要看得见它得让每一腿**再登记
"样本名册逐枚是谁"**，那是**新的一味** => 这一形**不在** `AC#9(2)` 判据文字里（那格只要求"造一发台件＋答响不响"）
 =>  **不阻塞**本格，但**必须作残余移交**（主件 sec 9 第 2 条）。
另 `171s` 自述（`s/00-anchor-and-gates.md:60-61`）：`AC#1` 正面两发对照"未做"、`AC#3`／`AC#5`／`AC#9` 其余格
"按派单押后" => **本腿用现跑补上了 AC#1／AC#3／AC#9**，不再依赖它们的自报。

---

## sec 5 写面：`attrib.sh` 那 4 行删除的逐字分类

四行删除逐字与性质：(1) `# back to being believed rather than checked.`＝**注释**；
(2) `#   sh …/attrib.sh --self-test   # classifier both directions + the ticket 275 roster control`＝**注释**；
(3) `echo "attrib.sh: SELF-TEST GREEN - …"`＝**打印绿提示语**；(4) `echo "attrib.sh: SELF-TEST RED - …"`＝**打印红提示语**。
 =>  注释 2 ＋ `echo` 2，`OTHER -> inspect` **0 行**；删行里含 `ST_FAIL`／`exit `／`[ `／`failures=`／`$((`／`st_case`／
`grep -q`／`cmp -s` ＝ **0 行**；三处机制行（`ST_FAIL=$((ST_FAIL + 1))`／`cases=$ST_RUN failures=$ST_FAIL`／
`ST_RUN=$((ST_RUN + 1))`）**deleted=0, added=0**；`exit` 只出现在**新增的注释句**里，**无一行被删或改**。
114 枚新增按首 token 分：**CODE 53／COMMENT 59／ECHO 2**。**正控**：含 `[ "$X" = 1 ] || exit 1` 的样本喂同一分类器
 =>  被判 `OTHER -> inspect` => 分类器**能**抓真断言删除。
**判**：四行删除全是**绿/红提示语改写**（旧句只提 classifier 与 275 门，新句补上 171 钉子），**无断言被改弱**。

---

## sec 6 门禁：四把尺的 rc 与末行逐字

跑前：`unset GOTMPDIR TMPDIR TEMP TMP`；`GOFLAGS=`（空）—— 环境复核逐字 `GOFLAGS=''`、
`GOTMPDIR='<unset>' TMPDIR='<unset>' TEMP='<unset>' TMP='<unset>'`。

| 门 | rc | 秒 | 输出行数 |
|---|---|---|---|
| `go build ./...` | **0** | 20 | 0 |
| `sh scripts/d22scan.sh` | **0** | 142 | 252 |
| `sh scripts/check-path-length-budget.sh` | **0** | 3 | 10 |
| `sh probes/161/r5/attrib.sh --tracked-only` | **1** | 38 | 42 |

末行逐字：(1) `go build ./...` => **（空，0 行输出）**
(2) `d22scan` => `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=244, bans #1-5 cmd/=42, ban #6 frontend/=116, ban #7 internal/tools/=25, ban #8 design/=39, ban #8 frontend/=116, ban #8 internal/=608, ban #8 cmd/=143; ban #8 emoji coverage: design/ 39 text files; frontend/ 116 text files; internal/ 608 Go files, comments and _test.go included; cmd/ 143 Go files, comments and _test.go included`
(3) `path-length-budget` => `check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree`
(4) `attrib.sh --tracked-only` => `attrib.sh: RED (tracked-only) - tracked-dirty=35 files_dirty=35. The CI gofmt step and this step read the same denominator; this one additionally refuses to answer when handed 0 files.`

**`tracked-dirty=35` 照实报、NOT 洗**。末三枚样本行逐字：
`TRACKED-DIRTY internal/agent/approval/pending_read.go (last commit touching it: d03d166f 10-04 10:37 CarlosShao)`／
`TRACKED-DIRTY internal/agent/tools.go (last commit touching it: 5413f46d 10-04 09:19 CarlosShao)`／
`TRACKED-DIRTY internal/risk/provenance.go (last commit touching it: db2f9eab 10-19 12:15 CarlosShao)`

 =>  三枚"最后碰它的提交"是 **10-04／10-19 的他票提交**，**不是**本票四笔（四笔 `.go`＝0）。
这一红是**开工前就在的既有账**，照实记、**不洗**。**前三把绿、第四把红**。

---

## sec 7 日志完整性（尺＝`git diff --numstat`，NOT status 计数）

**内容差**：开工前 **0 行**（本腿自落基线） ->  复核时 **15 行**（全部本腿自己跑台件造成） ->  还原后**全树 0 行**。
 =>  基线不是假设；且它**会**动（弄脏后立刻变 15 行） =>  这把尺不恒零。

**`status` 与 `diff` 的分歧本身就是读数**（三枚 stat 假脏全部 blob 相同、行数相同）：
`probes/161/r6/logs/flip-1.txt` `9f5b2dca`＝`9f5b2dca`（834 行＝834 行）／`flip-baseline.txt` `7f1c791e`＝`7f1c791e`／
`issues/115-…-done.md` `13a1ce79`＝`13a1ce79`。
 =>  **编排者另抓到的那 2 枚确实是"已复原但 stat 缓存未刷新"**，本腿逐字节证实 => 这就是用 `numstat` 的理由。

**本腿自己造成的脏（照实登记，全部已还原）**：跑那四把台件会**往 tracked 日志写**（不尊重"别人的日志"边界
——正是被裁腿 B 自报那件事的同一成因）。被覆写并已还原共 **21 枚**：`r1/logs/gate-{fixed,unfixed}-T{1,2,3}.txt` **6**／
`161/r6/logs/flip-{1,2,baseline}.txt` **3**／`r3/doors/logs/flip-{1,2,3,baseline}.txt` **4**／
`r3/logs/lie-door-unmodified.txt` **1**／`r3/logs/samples2-{D0,D1,D1c,D2,X0,X1,X2}.txt` **7**。
**还原方式**：`git show HEAD:<path> > <path>`；逐枚 `git hash-object` 与 `git rev-parse HEAD:` **比对相等**才算完。示例：
`flip-1.txt` `1214 -> 834` 行、`HEAD=9f5b2dca`＝`work=9f5b2dca`；`lie-door-unmodified.txt` `24 -> 66` 行、`bc9b6b0a`＝`bc9b6b0a`。

**"还原动作本身"是否合规 —— 合规**：允许形 `git show HEAD:<path> > <path>`；被禁形 `checkout --`／`restore`／
`reset`／`stash`／`clean`。扫 `probes/171/**` 下**全部 33 枚 `.sh`** 里**作为命令执行**的禁用动词（剔 `#` 注释行）
 =>  **命中 0 枚 / 33 枚**。被裁腿那枚 `restore-my-damage.sh:43` 逐字 `git show "HEAD:$f" >"$f"` => **正是允许形**
（`:16-18` 自写"NOT `checkout --`/`restore`/reset/stash/clean"；`:12` 登记截断 `1081 -> 586`／`66 -> 59`／`11 -> 2`）。
**正控（防扫描哑）**：自造含 `git checkout -- somefile` ＋ `git reset --hard` 的样本 => 同扫描**命中 2**。
 NOTE  **本腿自己踩的坑**：更宽的粗筛（把 `restore` 当普通词）命中 **7 枚脚本**的 `printf '  ok  restore matches…'`
之类——**那是打印文案，不是 git 动词**。据此把判据收窄到"真的 `git <动词>`"，NOT 用粗筛指控任何一枚腿。

---

## sec 10 本腿自证

- **提交**：`171-v1-verdict.md`＋`-part2`＋`-part3` 三枚。**NOT 未 push**；`git add`／`git commit` **两处都带显式 pathspec**。
- **写面**：上述三枚 ＋ `D:\tmp\wisp\20261021\171v1\**`（仓外）。`internal/**`／`cmd/**`／`scripts/**`／`.github/**`／
  `.scratch/wisp/**`／`docs/reports/**`／`docs/PLAN.md`／`docs/specs/**` **零字节**。
- **本腿落下的 `.go`**：**0 枚**。**越界**：**0 枚**。**正控**（同尺打 `a8eb0d51`） =>  `.go` **2 枚** => 尺**不恒零**。
- **格子**：`AC#1`／`AC#3`／`AC#5`／`AC#9` 的框**一枚未碰**；`Status:` 未碰；无 `-done`；台账未碰。
- **禁令码位**：一律"U+00A7"这类**名称形**；无字符本体。
- **还原**：本腿覆写的 21 枚 tracked 日志**全部**按 `git show HEAD:` 还原、逐枚 blob 相等；收尾**全树 `numstat` ＝ 0**。
- **未定义即停**：未遇需自行假设的方案空白；两处"方案没覆盖"（主件 sec 9 那两条账）已**上报**而非自行处置。
