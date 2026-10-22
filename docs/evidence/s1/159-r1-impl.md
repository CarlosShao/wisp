# 票 159 落地腿 159-r1 实现件（part 1/2）：三处时态过期句的就地追加更正

- 身份：写腿（实现者）。NOT 翻任何框、NOT 改票名／`Status:`／`-done`、NOT 改台账、NOT push。
- 起手锚 `ff21530f`（`git rev-parse HEAD`，本腿开工那一刻）。
- 写面：`cmd/wisp/slo_report_144_windows_test.go`（三处更正＋指针表）、
  `cmd/wisp/slo_comment_tense_159_windows_test.go`（新增，尺）、本件（part 1）＋`159-r1-impl-part2.md`。
- 一句话：三处"时态已过期"**全部按内容判为真过期**（不是关键词命中，见 sec 1.1 的判据）；
  就地追加更正，**未删原句一字**；补的尺**七发变异全红、还原全绿**（读数在 part 2）。

---

## 1. 先量：三处过期句的原文逐字 ＋ 每处"今天为什么过期"的盘上依据

### 1.1 我怎么"按内容判"而不是搜关键词（派单点名这一格是难点）

判据是**指涉物（referent）能不能在现盘上被找到**，不是词形：

1. 该句含**时间指涉**（`at this anchor`／`in today`／现在时断言"某缺陷仍在"）；
2. 该指涉物**在写入那一刻**是可解析的（我要能说出它当时指什么）；
3. 该指涉物**在现盘上已不可解析**或已改真  =>  过期。

三处逐条过这三关；不满足第 2 关的（"找不到它当时指什么"）我**不判过期**。
现盘 `grep -n '过时|已过期'` 在本文件**零命中**——派单这句我复核为真，所以确实不能靠关键词。

### 1.2 三处原文（逐字，现盘 `:742`／`:749`＋`:750`／`:752`）

```
742 | // Correction of that triple, measured at this anchor instead of copied - ticket
749 | // block names says 20/13/0 as well. So "=== RUN=21, PASS=14, FAIL=0" is neither
750 | // reading: it is this selector's UNMUTATED baseline, 21/14/0, pasted into the
752 | // 21/14/0 again at this anchor: probes/156/mut-156-r2/r2-asis-152-surface.log).
```

### 1.3 逐处的盘上依据

**第 (1) 处 `:742` ＝ `measured at this anchor`**

- 写入者 ＝ `2ce77e18`（`git log -S 'at this anchor' -- cmd/wisp/slo_report_144_windows_test.go` 唯一命中）。
  它的**父提交** ＝ `0f186525`（`git log -1 --format='%h %p' 2ce77e18`  =>  `2ce77e18 0f186525`）。
  该块自己在 `:743` 把锚点写死成 `At head 0f18652`  =>  **"this anchor" 的指涉物 ＝ `0f18652`**，可解析。
- 今天：`git rev-list --count 0f18652..HEAD` ＝ **4076**；`git diff --stat 0f18652..HEAD -- <本文件>`
  ＝ `56 insertions(+), 1 deletion(-)`  =>  该文件本身也被改过（`1d53526c` 改了 `:739`）。
  =>  读者按字面"at this anchor"去量**今天的树**，量到的不是这句被写下时的树。
  **过期形 ＝ 时态／指涉，不是数值**（数值我复算为真，见 sec 1.4）。

**第 (2) 处 `:749`＋`:750` ＝ 现在时断言"某缺陷仍在"**

- 原句断言：`"=== RUN=21, PASS=14, FAIL=0" is neither reading ... pasted into the sentence that
  promises the mutation's numbers`。写入时（`2ce77e18`）为真：那时 `:739` 里确实是 `21/14/0`。
- 今天为假：`1d53526c`（`git log -1` 主题逐字含"把 case 14 注释里那三个数改成 r4 自己量到的变异读数"）
  把 `:739` 改成了 `20/13/0`。现盘 `sed -n '739p'` ＝
  `// g1-restored-borrowed-value-case14off.log reports === RUN=20, PASS=13, FAIL=0);`
  =>  "the sentence that promises the mutation's numbers" 里**已经是读数**，那句"is neither reading"
  所描述的缺陷**在盘上已不存在**。
- NOTE 验收程 `docs/evidence/s1/156-close-and-load-bearing-r1-accept-r1.md:478` 已自己量过这一半：
  `reports === RUN=21…` 那句句式今天 `grep -c` ＝ 0 命中。我复核：现盘该串**仅作为 `:749` 的引文**存在
  （`grep -c 'RUN=21, PASS=14, FAIL=0'` ＝ 1，位置 `:749`，即引文本身）。

**第 (3) 处 `:752` ＝ `21/14/0 again at this anchor`**

- 同 (1)：同一个 deictic、同一个指涉物 `0f18652`、同一个过期理由。

### 1.4 NOTE 反过来说明"什么没有过期"（这是本票最容易被做错的一半）

三处**数值全部仍为真**。我用盘上日志逐枚重数（不是抄注释）：

```
. scratch/wisp/probes/152/mut-post/g1-restored-borrowed-value-case14off.log  RUN=20 PASS=13 FAIL=0
. scratch/wisp/probes/152/mut-post/g1-restored-borrowed-value.log            RUN=21 PASS=13 FAIL=1
. scratch/wisp/probes/152/mut-post/asis.log                                  RUN=21 PASS=14 FAIL=0
. scratch/wisp/probes/156/mut-156-r2/r2-152g1-case14off.log                  RUN=20 PASS=13 FAIL=0
. scratch/wisp/probes/156/mut-156-r2/r2-152g1-case14live.log                 RUN=21 PASS=13 FAIL=1
. scratch/wisp/probes/156/mut-156-r2/r2-asis-152-surface.log                 RUN=21 PASS=14 FAIL=0
```

六枚**全部 `git ls-files` tracked**（`git ls-files --error-unmatch` 逐枚 rc=0）。
=> 所以本票的修法只能是**改时态**，NOT 改数；这一点我写进更正块，也写进尺子的断言。

---

## 2. 就地追加更正：写了什么、在哪、为什么放这里

**位置**：紧接 `:752` 之后（现盘 `:754` 起，块尾 `:789`），原句 `:742`–`:752` **一字未动**。
**理由**：读 `:749` 的人必须**先遇到**它已过期的更正，再往下一段走；放文件末尾等于没更正。

更正块逐字（现盘 `:754`–`:789`，此处摘结构与判据行，全文见盘）：

```
754 | // Ticket 159 AC#2 (C2), appended by the record-level cleanup leg: all three of the
... | //   （说明：三句 went stale 于 1d53526c 把重测数写进 :739；本块不编辑原句一字）
760 | // 159 (1) :742 - |Correction of that triple, measured at this anchor instead of copied|
... | //   "This anchor" 的指涉物 ＝ 0f18652（2ce77e18 的父），今天 head 已在其后 4076 笔
767 | // 159 (2) :749 - |block names says 20/13/0 as well. So "=== RUN=21, PASS=14, FAIL=0" is neither|
768 | // 159 (2) :750 - |reading: it is this selector's UNMUTATED baseline, 21/14/0, pasted into the|
... | //   两句当时为真、今天为假：promise 数的那句是 :739，1d53526c 已放进 20/13/0
772 | // 159 (3) :752 - |21/14/0 again at this anchor|
... | //   同 (1) 的 deictic、同指涉物、同过期
... | // The numbers did not expire ... 六个数重数仍为 20/13/0、21/13/1、21/14/0 两处
```

**四枚 pill 引文与原件的关系是可断言的**：每枚 `|...|` 是它所指名那一行的**逐字尾部**
（不带前导 `// `）。这条不是我的声称——part 2 sec 4 的尺逐枚断言它。

### 2.1 NOTE 这次插入自己造成的代价，我如实登记并修掉（NOT 藏）

插 56 行在 Case 13 块**之前**  =>  该块里**四条自指行号**全部短了 37 行（`:835`->`:872`、
`:851`->`:888`、`:849`->`:886`、`:857`->`:894`）。这正是本票处理的那一族
（"注释承诺 X、X 被就地改真"）——**我在修这一族时当场又造了一例，所以必须自报**。
处理口径：**NOT 改那四句原文**（那是别人已入库的读数，append-only），
在 Case 13 块尾**追加一张指针表**（现盘 `:928`–`:932`），逐行写 `said`（当时）->`says`（现在）->该行内容。
=> 这就是本票第 (1)/(3) 处同一形状的**第二例**，我把它一并闭在同一把尺下（part 2 sec 3 的第 9 枚子用例）。
