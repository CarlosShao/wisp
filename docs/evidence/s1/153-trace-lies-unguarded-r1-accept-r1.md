# 153 — 那枚压缩痕的两条谎法 · 对抗验收件 r1（非实现者，只读取证）

被验版本＝**`6de3d1c5`**（票 153 实现件最后一枚）。实现链条＝码 `ff550f3`（AC#1 新用例）＋`b23c7f7`（痕带 taskID）→ 证据件逐格 `ac7fb00`／`f58c513`／`bec1727`／`9a8766e`／`6de3d1c`。
派单＝`.scratch/wisp/dispatches/2026-09-26-0958x-accept-153.md`。实现件＝`docs/evidence/s1/153-trace-lies-unguarded-r1.md`。
本程角色＝非实现者对抗验收（D22 双角色）。**只读取证**：零生产码改动、零工单改动、零台账改动、零 push；只写本文件与 `.scratch/wisp/probes/153/accept-r1/**`。
临时件一律**只建不删**、且建在仓库外（`D:\tmp\wisp153-accept-r1\`）。票面四框本程一枚都不勾，也不替实现方勾。

**总裁决（详见 §11）**：**附条件入账**。三格成立、三格成立但点名待补、一格（攻击点①）本程**没造出**"重写后变弱"的实质缺陷；
本程未造出任何一发"判据失灵"的活证据，故不入"退回"档。三条待补项全是记账级，无一条动代码。

---

## 0. 锚点、取版纪律与工作树现量（本程第一发）

```
$ date '+%Y-%m-%d %H:%M:%S %z'      2026-09-26 10:23:54 +0800
$ git rev-parse --short HEAD        7263456          ← 不是被验版本
$ git log --oneline -3
7263456 evidence(145 §5-§7): 门禁改前改后四数＋两包名册差集为空、没核十条、报回七条
eb38c97 docs(145 载体指针): ...
972ceba evidence(145 §4): AC#3 两角色一次定死 ...
$ git status --porcelain=v1 | wc -l   40             ← 脏的全是别家的
```

⇒ 开工此刻**票 145 的程正在提交**（HEAD 已是它的 `7263456`），脏树里另有 `design/**`＋`frontend/**`（前端会话）与
`.scratch/wisp/probes/152/my152.py`（票 152 的程）。**本程一次都没读脏工作树当被验版本**：

| 取什么 | 怎么取 |
|---|---|
| 证据件三版 | `git show bec1727:<path>` / `git show 9a8766e:<path>` / `git show 6de3d1c:<path>` 各落一份到仓库外 |
| 被验的码 | `git archive 6de3d1c \| tar -x` → `D:\tmp\wisp153-accept-r1\snap-delivered`（27 MB，无 `.git`） |
| "改前"的码 | `git archive 86b0161 \| tar -x` → `...\snap-anchor`（实现程自报锚点，同锚） |
| 变异 | 只发生在上述两枚仓外快照里；每发前 `restore` 并与 `pristine/` 逐枚比字节（打印 `RESTORE OK (3/3 byte-identical to pristine)`） |
| 交付树自证 | `md5sum internal/agent/{compress.go,loop.go,compress_trace_test.go}` 快照＝`git show 6de3d1c:…` 三枚全 MATCH（`fa28dc33` / `0bb2e016` / `e72bfb0c`） |

派单里那条"当日实测的仪器坑"本程**先撞上再核**：`-overlay` 与 `-cover*` 合用时 overlay 被静默忽略 ⇒ 本程**不用 overlay**，
直接改快照里的文件本体，并在每一发变异后打印"新文本出现次数＋文件字节差"，让"变异没落地"这一族当场暴露。

---

## 1. 攻击点① —— "实现程重写了自己已提交的证据件"：逐格原文并排

派单要求：把"重写前"与"重写后"的 §§5–9 原文并排走一遍，专找 ① "本程没测什么"的条目 ② 被改短的红句 ③ 对它自己不利的措辞。

### 1.1 先纠正派单自己的两处前提（这两条不成立，必须先摊开）

**(a) `bec1727` 里根本没有 §§5–9。** 派单写"把 `bec1727` 与 `6de3d1c` 两版的 §§5–9 并排"，实测量：

```
$ git show bec1727:… | 最后一个标题         329:### 4.3 本格自证：还原与"跑的就是交付码"     总行 342
$ git show 9a8766e:… | grep -n '^## '       只到   691:## 5. AC#4 — 契约轴零字节            总行 830
$ git show 6de3d1c:… | grep -n '^## '       346:## 5. AC#4 … / 420:## 6. AC#5 … / 488:## 7. 本程没测什么
                                            511:## 8. 伪授权两栏 / 544:## 9. 收口            总行 569
```

⇒ §5 **首次出现**就是 `9a8766e`（被那次排版事故污染的那枚），§§6–9 **首次出现**就是 `6de3d1c`。
所以可核的"并排"只有一对：**§5 的 `9a8766e` 版 vs `6de3d1c` 版**；§§6–9 不存在"原稿"，也无从"变弱"。
实现件 §9 那句"§5–§9 是从干净版 `bec1727` 起重写"用的词是"从……状态"（把文件退回 `bec1727` 那份只有 §§0–4 的树），
不是"从 `bec1727` 里抄 §§5–9"——**这句话本身没撒谎，是派单把它读成了后者**。

**(b) "行尾"那一半在可核对象上不存在。** 全仓 `.gitattributes` 写着 `*.md text eol=lf`，`core.autocrlf=true`，
而三版 blob 实测：`bec1727 CRLF=0`／`9a8766e CRLF=0`／`6de3d1c CRLF=0`（`od -c` 现量，非 grep）。
⇒ 就算它的工作树缓冲真被写成 CR/CRLF 混排，**git 也存不进去**；留在 `9a8766e` 那枚 commit 里的事故形态只有一种：**空行膨胀**。
比例本程独立复算，与它自报一致：

```
bec1727  总行 342 非空 265 空行 78   空/非空 = 0.29
9a8766e  总行 830 非空 320 空行 511  空/非空 = 1.60   ← 坏的那版（它自报 1.59）
6de3d1c  总行 569 非空 441 空行 129  空/非空 = 0.29   ← 修回来
```

它给的"内容零丢失"算式 `320 = 265 + 55`（§5 的 55 枚非空行）本程复算成立。**注意：本程判这条用的是 blob，不是它的缓冲；
"缓冲里行尾坏过"这件事不可核也不被本程判假——它只是不进判据。**

### 1.2 §§0–4（含 AC#1 红句与全部读数）在重写后**逐字节相同**

```
$ head -342 bec1727.md > q1 ; head -342 6de3d1c.md > q2 ; cmp q1 q2
（无输出，rc=0）      ← 含 CR 的逐字节比，不是"规范化后相同"
```

⇒ 派单点名的三样东西里，**"被改短的红句"这一项在 §§0–4 结构上不可能发生**：§2.2 那两行红句原文（`compress_trace_test.go:458`／`:461`）、
§1.2 那句"后半只在进程内复现、未复跑到盘"、§3.3 那枚"C39 不存在"的答句、§4.1 九发矩阵——**一个字节都没动**。

### 1.3 §5：`9a8766e` 版 vs `6de3d1c` 版，逐格并排（先给尺，再给读数）

本程用的尺＝两版各取 §5 的非空行 → 去掉 `**` → 折成单流 → **词级 `difflib` 全差集**（不是"读一遍觉得没少"）。
读数：`old_tokens=296 → new_tokens=308`（**变多**），非空行 `54 → 56`。全部差集只有下面 12 处，**一处不落抄在下面**：

| # | 位置 | 重写前（`9a8766e` 原文） | 重写后（`6de3d1c` 原文） | 本程判 |
|---|---|---|---|---|
| 1 | 坏尺点名 | 全编队共用 `git config user.name/email` 同一枚身份 | 全编队共用同一枚 git 身份 | 同义，信息不减 |
| 2 | 同上 | 把 152／151／前端／本程全印成同一个名字 | 把 152／151／前端／**编排者**／本程全印成同一个名字 | **变强**（名册多一类） |
| 3 | 同上 | "我这一枚区间里谁提交了什么" | "这一格是谁交的" | 同义 |
| 4 | 名册枚数 | 本程到此共五枚 commit，逐枚名册**（本节这枚证据件是第六枚，不列自己）** | 本程到这一枚为止五枚 commit，逐枚名册 | **唯一一枚消失的措辞**，见 §1.4 |
| 5 | 别家名册 | `5365cb2`、`cdf2471`/`ad9b29f`、`eb4755a`/`4cc85bb`、`45c920e`/`5d46f24`/`4e16976`、`720cae6` | 同上 **＋`8883b3f`（前端）＋`10e3585`/`6550dc4`（票 152）** | **变强**（多三枚有名） |
| 6 | 同上 | "之后区间又长，本程收口时另读到 `eed7c99`/`63e228f`" | "本程收口时区间又长了 `eed7c99`／`63e228f`" | 同义 |
| 7 | 新增段 | （无） | 追加"共享 index 下本程每次 commit 前都现量了 `git diff --cached --name-only`，五枚里没混进别家路径（一枚 `git add` 之后 index 名册见 §9）" | **变强**（新增自证） |
| 8 | 表格拆分 | `\| f58c513 / bec1727 \| AC#2 / AC#3 表 \| 只有那枚证据件 \|` | 拆成两行，各自有名 | 同信息、可核粒度更细 |
| 9 | 禁改面尺 | `grep -E "…approval/\|⏎ ^internal/observe/\|…"`（三行折行，行尾带 `\|`） | 同一条管道**合成一行** | **变强**：折行版若被照抄执行，`\| ^internal/observe` 前会带一枚空格⇒那条支永远不命中；合行后没有这个洞 |
| 10 | 正控注释 | "同一把尺、同一管道上的正控：尺活着" | "正控：尺活着，0 命中不是尺坏" | **变强**（补了"0 命中"的排他解释） |
| 11 | DEFERRED 自打 | "…改成"the D28-1 deferral…"后才跑计数——**这条是靠"改之前先想到"没成的，不是靠读数捞回来的**，记在这里以免读者以为本程有门兜着" | "…改成**转述**（"the D28-1 deferral…"）之后才跑计数——**这一条是"写之前想到"没成的，不是被读数捞回来的**，记在这里，免得读者以为有门兜着" | **同一枚自打，措辞等重**（"顶到 4 枚"→"顶成 4 枚"，纯字） |
| 12 | 12000 名册 | "（139 验收 §3.4 那枚"本票贡献 0 枚"过期陷阱）" | 同句 **＋"本程先量再落笔"** | 这句是**向它自己偏**的一处（本程 §1.9 复算：anchor 9／delivered 9／单枚 1／1 全对，"先量"没有被冒用） |

**三样点名要查的东西，逐样答：**
- ① "本程没测什么"的条目：那一节（§7）**首版就在 `6de3d1c`**，`9a8766e` 里根本没有 ⇒ 不存在被删；且它的七条里第 1、2、7 条都在**主动报自己没读数**（见 §5、§8 本程复判）。
- ② 红句原文：全部在 §§0–4，`cmp` 判逐字节相同；红句里的 `:458`／`:461` 两个行号本程独立复跑**打中同一枚行号**（§3 表）。
- ③ 对它自己不利的措辞：两处自打（`DEFERRED(D28-1)` 差点顶成 4 枚／`gofmt` 与 `gofumpt` 双双点名而 `go vet` 不响）＋一处"未复跑到盘"＋一处"cmd/ 差 1 枚不是我造的"——**全部在场**（前两处在上表 #11 与 §6.2；后两处逐字节未动）。

### 1.4 唯一一枚消失的东西，与本格的判法

`（本节这枚证据件是第六枚，不列自己）`——**全件搜 `不列自己`／`第六枚` 在 `6de3d1c` 版：0 命中**（`grep -c` 现量）。
它的作用是把"名册只列五枚"这件事自我限定住；删掉之后，读者拿 §5 那句"本程到这一枚为止五枚 commit"去对
`git log`（收口时实为七枚：`ff550f3` `b23c7f7` `ac7fb00` `f58c513` `bec1727` `9a8766e` `6de3d1c`）会以为对不上。
**判**：这是一枚**自陈精度的退化**，不是判据、不是读数、不是不利结论——它删不掉任何一发红。
派单写"只要有一格比原稿弱，本格判退回"。**本程不替编排者把这句话自动放大到覆盖一枚括回**：
按字面收，本格的"退回条件"是**判据变弱**（红句变短、没测项被抹、不利措辞消失），这三条本程逐条量过、一条都没发生；
若编排者要按"任何措辞级退化即退回"执行，闭合面＝**把那句括回补回去**（一格一行，不动码、不动读数），本程把它列进 §11 待补项第 3 条。

**本格档位：成立（附一条措辞级待补）。** 理由：作者事后改写自己的证据这件事**确实发生过**（`9a8766e` 的 §5 被整体替换），
但本程把可核的两版并排走完，**没找到一处判据或读数变弱**，且 §§0–4 是逐字节还原——"抽掉可信度地基"这一后果今天没有落地。

---

## 2. 攻击点② —— `DEFERRED(D-xx)` 计数：实现程"commit 前抓到"那句，现量成立（但派单那条判据要连口径一起读）

### 2.1 交付版三枚文件里 `DEFERRED(` 的枚数与位置

```
$ git show 6de3d1c:internal/agent/compress.go | grep -n 'DEFERRED('
27:// DEFERRED(D28-1): move the call into the Warm window hook; keep this
53:// Warm-window hook DEFERRED(D28-1) installs will get it for free); it pairs
$ git show 6de3d1c:internal/agent/loop.go | grep -n 'DEFERRED('
394:                    // DEFERRED(D28-1): the Warm-window hook owns this call; the
$ git show 6de3d1c:internal/agent/compress_trace_test.go | grep -c 'DEFERRED('
0
```

同一把尺打在锚点 `86b0161`：`compress.go` 27 与 **49**、`loop.go` 394、`_test.go` 0。
⇒ 交付版没有多出一枚标记：第二枚只是**从 :49 挪到 :53**（它上面新加了四行注释）。本程第一发把 `git grep -n` 的行号一起
喂进了 `comm`，于是"挪位置"被读成"新增一枚"——那是**本程自己的尺坏了**，不是它的件坏了；改用"只取内容"重量即归零（见 §9）。

### 2.2 全仓两枚计数（口径写死，否则下一位会误判）

| 尺 | `86b0161` | `b23c7f7`（它的码） | `6de3d1c`（交付版） |
|---|---|---|---|
| `git grep -oh 'DEFERRED(D28-1)' … -- '*.go'` | **3** | **3** | **3** |
| `git grep -oh 'DEFERRED(D11-3)' … -- '*.go'` | **1** | **1** | **1** |
| 同一条尺**不带 `-- '*.go'`**（含 `*.md`／`probes/**`） | 39 | 39 | **40** ← 会假报漂移 |

⇒ 那句"计数仍是 3／1"**只在代码标记这个口径上成立**，而 `AGENTS.md §1.1` 的硬规矩写的正是"**代码标记**与 `SPEC-12 §5` 1:1"，
所以本程判它**对**；但差值那一枚本程具名到了行：`40` 比 `39` 多的那一枚＝**它自己证据件 §5 里那句"第一版真的写了字面量 `DEFERRED(D28-1)`"**
（`git show b23c7f7:…md | grep -c` = 0 → `git show 6de3d1c:…md | grep -c` = 1）。
**判**：这是一枚**文档里的提及**，不是代码标记，不顶任何 1:1 的账；但它是"下一位拿无路径过滤的尺去数就得出漂移"的形状，
本程按派单要求的"具名＋给可复算尺"记在这里，不当违规。

### 2.3 那两枚标记在 `SPEC-12 §5` 登记表里到底有没有行（本程顺手核，因为派单点名"1:1 是硬规矩"）

```
$ git grep -n "D28-1\|D11-3" 6de3d1c -- docs/specs/SPEC-12-roadmap-governance.md
（无输出，rc=1）        ← 登记表里两枚都没有行
$ ls .scratch/wisp/issues | grep '^150-'
150-four-deferred-markers-have-no-row-in-the-spec-12-registry-so-the-1-to-1-rule-is-violated-today-and-no-instrument-scans-it.md
```

⇒ 1:1 **今天确实不符**，但那不是我票的事：票 150 的标题就是"四枚 DEFERRED 标记在 SPEC-12 登记表里没有行、且没有仪器在扫"。
**本程判**：票 153 既没新增标记也没让这件事变坏，**本格成立**；顺带回派单那句"这条不是挑刺"——它确实不是，只是**归票 150**。

**本格档位：成立。**

---

## 3. 攻击点③ —— "那枚新用例是不是唯一证人"：本程自己跑的 12 发（不背它的 X 表）

### 3.1 先钉判据与尺

尺＝`.scratch/wisp/probes/153/accept-r1/ruler153.py`：只认锚定的 `^--- FAIL` 行判红（`t.Logf` 也带 `file:line:` 前缀，
不能拿它当红），四数之外还出 `panic` 数与名册；`-count=1`、**跑整包**（不只那几枚）。
每发变异由 `mut153.py` 落地并**强制打印**：锚点字符串出现次数必须恰为 1（否则 `SystemExit`）、落地后新文本次数与文件字节差。
每发之前 `restore` 并把三枚文件与 `pristine/`（＝`git show 6de3d1c:…`，md5 逐枚 MATCH）比字节 ⇒ `RESTORE OK (3/3 byte-identical to pristine)`。

### 3.2 读数（本程亲跑；"它的 X"一列是从它 §4.1 表里抄来对照用的，不是本程的凭据）

| 本程发 | 施了什么 | RUN / FAIL | 红名 | 它的 X | 判 |
|---|---|---|---|---|---|
| `delivered-count1` | 交付码，无变异 | 83 / **0** | — | X0 | 一致 |
| `A1-m5` | 只施 M5（`if rep.Ran {`→`if c.Need(hist) {`，落在 `compress.go:216`） | 83 / **1** | `SilentWhenNothingFoldableOverThreshold` | X1 | **一致** |
| `A2-m5-dropAC1test` | M5 ＋ 摘掉那枚新用例 | 82 / **0** | **NONE ⇒ 逃逸回来** | X2 | **一致** |
| `A3-dropAC1test-only` | **只摘新用例、不改码**（派单的 ①） | 82 / **0** | NONE | — | **派单前提不成立**，见 §3.4 |
| `A4-dropN2-looptag` | 摘 N2（`loop.go:399` 不再打标） | 83 / **1** | `CarriesTheOwningTaskID` | X3 | 一致 |
| `A5-dropN2-dropowntaskidtest` | 摘 N2 ＋ 摘掉那枚证人用例 | 82 / **0** | NONE | X4 | 一致 |
| `A6-dropN3-attrwrite` | 摘 N3（痕不再写 `task` 属性） | 83 / **2** | `CarriesTheOwningTaskID` ＋ `NeverInventsATaskID` | X5 | 一致 |
| `A7-dropN3-dropbothtasktests` | 摘 N3 ＋ 摘掉那两枚证人用例 | 81 / **0** | NONE | X8 | 一致 |
| `A8-m6-unconditional-attr` | 变异：属性**无条件**追加（不判空） | 83 / **1** | `NeverInventsATaskID` | X6 | 一致 |
| `A9-m7-constant-fake-id` | 变异：`traceTaskID` 恒返一枚 uuid 常量 | 83 / **2** | 两枚任务用例都红 | X7 | 一致 |
| `anchor-count1` | 锚点 `86b0161` 原树（四枚钉子、无新用例） | 80 / **0** | — | 它的 §6.1 改前 | 一致（80） |
| `anchor-m5-count1` | **未修码**上只施 M5（139 那一发的复算） | 80 / **0** | 四枚钉子逐枚 PASS | 它的 R1 | **一致：逃逸复现** |

⇒ 它 §4.1 那张九发表**本程独立重走出同一批数**，包括三处最容易糊的地方：X2／X4／X8 是"把证人和被证的东西一起摘掉"才绿，
少摘一枚当场红（A1/A4/A6 就是另一半）；`panic` 全 0、两向名册差集见 §6.1，没有"一条红吞掉整包读数"的形状。

### 3.3 改前红句原文——本程复跑打到同一枚行号

```
A1-m5 读数（交付树＋M5，`-count=1`）：
    compress_trace_test.go:458: over-threshold pass that folded nothing left a trace: agent: history compressed kept_raw_rounds=1 history_changed=false tokens_before=1006 tokens_after=1006 threshold=384 msgs_before=3 msgs_after=3 compressed_msgs=0
    compress_trace_test.go:461: "agent: history compressed" records = 1, want none: … 1 raw round(s) against the KeepRawRounds floor of 3 is exactly the case where no fold can happen (all records: [trace-capture-ruler-control agent: history compressed])
--- FAIL: TestCompressionTraceSilentWhenNothingFoldableOverThreshold (0.00s)
```

与它 §2.2 贴的那两行：行号 `:458`／`:461` **相同**；红句里的**八枚 `键=值` 全相同**（`tokens_before=1006 tokens_after=1006 threshold=384 msgs_before=3 msgs_after=3 compressed_msgs=0 kept_raw_rounds=1 history_changed=false`）；
**唯一不同是打印顺序**——它把两行按 `compressed_msgs` 在前抄出，本程读到的是 `kept_raw_rounds` 在前。
原因本程读到源码层级：`traceCapture.Handle` 把属性收进 `map[string]any`（`compress_trace_test.go:64`／`:67`），
而打印用的 `capturedRecord.flat()` 在 `:117-124` 里 `for k, v := range r.attr` ⇒ **map 遍历顺序，天生不稳定**。
⇒ 判：红句**内容逐对可复算、行号可复算**，"原文照贴"这件事成立；但**顺序不是可复算量**，
下一位若拿"逐字符相同"去对这两行会误判——本程把这条差异写死在这里（不构成退件，见 §9 结论修正）。

### 3.4 派单前提①「只摘新用例、不改码 ⇒ 应红」**结构上产不出它想量的读数**

本程照字面跑了 `A3`：`82 / 0`，**绿**。理由不是环境也不是判据变了，是这句话本身不可满足：
交付树上的生产码守卫是**正的那枚**（`if rep.Ran`），摘掉一枚 `_test.go` 里的用例**不会造出红**——
测试只能拒代码，不能因少一枚而多拒。⇒ 本程**不**按它字面判实现件红，并给出这句话唯一有意义的读法：
**"摘证人"必须与"施变异"同发**，也就是 `A1↔A2`（同一枚 M5，有证人红／摘了证人从此不红）。
这一对已由实现件与本程**各自独立**跑到，故"新用例是 M5 的唯一证人"这句话**入账**。
（派单第 3 点另两句——②"只上 M5 不带新用例应绿"、③"新用例＋M5 应红"——本程复算**都成立**，见 `A2`／`A1`。）

**本格档位：成立。** 新用例是真判据（它单独把 M5 打红），既不是装饰也不是恒红：
`A0`（无变异）绿、`A1`（只变异）红、`A2`（变异＋摘证人）绿——三态齐全，"今天不响、修了才响"那一发本程也替它在未修树上复算过（`anchor-m5-count1`＝139 的逃逸原状）。

---

## 4. 攻击点④ —— "未碰契约面"那三问本程独立重走（＋派单自己那句"8 枚字段名各 0 命中"不成立）

### 4.1 那把尺先证明它能命中，再谈"0"

`docs/PLAN.md`＋`docs/specs/**` 这一范围内，本程先拿**一枚确实被冻结的字段名**打正控：

```
$ git grep -n "keep_transcript" 6de3d1c -- docs/PLAN.md docs/specs
docs/PLAN.md:2735 | `[privacy]` | `redact_paths` … **`keep_transcript`/`keep_audio` 为硬编码 false…**
docs/specs/SPEC-03-config-secrets-envs.md:37  （同一行形状的 section 表）        → 2 命中，尺活着
```

同一把尺打到那枚痕上：

```
$ git grep -n "history compressed\|compressed_msgs\|kept_raw_rounds\|history_changed" 6de3d1c -- docs/PLAN.md docs/specs
（无输出）
$ git grep -n "tokens_before" … → 0    tokens_after → 0    msgs_before → 0    msgs_after → 0
$ git grep -n "threshold"       … → 6 枚命中行，不是 0；逐枚看全是别的东西：
      PLAN.md:2728 `[voice] wake_word{… thresholds[] …}`   PLAN.md:2731 `loop_guard{repeat_thresholds:[3,5,8]}`
      PLAN.md:2738 `[cost] … alert_threshold`              SPEC-03:30 / :33 / :40 同三处
```

⇒ **派单第 4 点那句"那 8 枚字段名与新增的 `task` 键各 0 命中"按字面不成立**：`threshold` 有两枚无关命中，
而 `task` 是日常词（`docs/PLAN.md` 里成百枚），"0 命中"这句话对它毫无意义。
**成立的是它那半句实质**：痕的那张字段表**在 PLAN/specs 里不存在**（七枚独有名 0 命中＋唯一同名的 `threshold` 是配置键）⇒
"这张表是票 139 的实现选择、不是契约"这句话**入账**。本程把它重写成可核的形状：
**"零命中"只对独有名成立，不对通用词成立**，任何下一位拿 `task`/`threshold` 去数都会读出假号。

### 4.2 "没有 C39 这枚契约、契约是 C1–C32、`docs/contracts/` 是 0 个文件"——三句全部复算成立

```
$ git grep -n "^| C12" 6de3d1c -- docs/PLAN.md    → docs/PLAN.md:1362（正控：这张表就是 `^| C<N> ` 形状）
$ git grep -n "^| C39" 6de3d1c -- docs/PLAN.md    → 无输出
$ git grep -oh "C[0-9]\{1,2\} " 6de3d1c -- docs/PLAN.md | sort -u -V | tail -1   → C32
$ git ls-tree -r --name-only 6de3d1c | grep -c "^docs/contracts/"                → 0
```

⇒ 它 §3.3 那半句（"票面问的 `C39` 这枚编号不存在"）**成立**，且它没拿"编号不存在"当免检：它同时答了 `D39` 里有没有这张表。

### 4.3 那笔账归派单方：`D39` 被工单错写成 `C39`

派单第 4 点末要求本程"回查我的工单／派单里有没有把 `D39` 错写成 `C39`"。**有，就在票 153 票面上**：

```
$ git grep -n "C39" 6de3d1c -- docs .scratch
.scratch/wisp/issues/153-…-no-taskid.md:27: ③ 把 taskID 送进那条 Info 要不要改**契约面**（`C12` 状态机、**`C39` 契约深化**里那条痕的字段表）？
.scratch/wisp/probes/153/ac2-q3.txt:6-7:      -- C39 exists as a contract? (contracts run C1-C32) / (no C39 row, rc=1)
docs/evidence/s1/153-trace-lies-unguarded-r1.md:225: 它答这一枚的那一行
```

⇒ 全仓 `C39` 只有三处命中，**第一处就是工单那枚错号**，后两处是实现对它的回答。`D39` 才是"契约深化"那枚决策
（`AGENTS.md §5` 与 `PLAN.md:2926` 同名）。**这笔账按派单说的归我（编排者／开票人），本程在此点亮它**：
不是实现件的缺陷；实现件的处理也**不算蒙混**——它没有用"该编号不存在"把这一问打发掉，而是两读都答了（C39 不存在／D39 里没有这张表）。
⚠ 但工单这一枚错号有实际代价：**"未定义即停"的判据是"要不要动契约面"，编号写错会让下一位在 `C` 表里找一个不存在的东西**，
找不到时最容易滑成"那就当我没这问"。本程把它写进 §11 待办给编排者（改票面一行字，不改判据）。

### 4.4 `C17` 白名单不受影响 ＋ 交付码没碰任何禁改面（本程自己那把尺）

```
$ git grep -c "compression\|compressed" 6de3d1c -- internal/panel docs/specs/SPEC-08-ui-ball-panel.md
（无输出，rc=1）                        ⇒ 一条 Go 侧 slog 属性不过 C17 那条线（既不在白名单、也不需要进）
$ git show --pretty=format: --name-only ff550f3 b23c7f7 | sed '/^$/d'
internal/agent/compress_trace_test.go  internal/agent/compress.go  internal/agent/compress_trace_test.go  internal/agent/loop.go
$ 同一条管道 | grep -E "^docs/PLAN.md|^docs/specs/|^internal/risk/|^internal/panel/|^internal/agent/approval/|^internal/observe/|thresholds\.go|golden|allowlist\.txt|^scripts/slo-check\.ps1|^tools/d22scan/|^frontend/|^design/"
（零命中）  ；正控 grep -cE "^internal/agent/" = 4（同一根管子上）
$ sed -n '/func (c \*Compressor) Need/,/^}/p' 两版 cmp → Need() IDENTICAL
$ "if len(raw) <= c.b.KeepRawRounds"      anchor=1 delivered=1
$ "for c.totalOf(rounds) > c.b.HistoryCompressTokens"  anchor=1 delivered=1
$ git diff --numstat 86b0161 6de3d1c -- internal/agent/budgets.go → 0 行
```

⇒ 派单第 4 点那句"C17 白名单不受影响"**成立**；`Need()`／`KeepRawRounds`／阈值**一字节没动**（这条最要紧，
因为"让常见路径也折叠"是这枚痕唯一的"把谎说圆"捷径，它没走）。`frontend/**`／`design/**` 零触碰、零宣称 ✓。

**本格档位：成立。**

---

## 5. 攻击点⑤ —— `withTraceTask` 是一条新的隐式通道：两个方向本程各自取到读数

### 5.1 方向②（先看硬的）：缺 ID 时那一行**到底**长什么样——到盘读数，本程补上了实现件自己没测的那一环

实现件 §7 第 1 条自报"没把那枚痕复跑到盘、因此没验 `redactHandler` 会不会改写字段名"。本程把这一环补上：
在仓外快照里加一枚**只存在于快照**的探针（`.scratch/wisp/probes/153/accept-r1/zzaccept153_ondisk_test.go`），
走**真的宿主日志管链**（`observe.InitLogWithRegistry` → `redactHandler`（D33 硬编码、不可关）→ `slog.NewJSONHandler` → 滚动文件），
再回读那个文件。读数原文（`probes/153/accept-r1/probe-ondisk.txt`）：

```
{"time":"…","level":"INFO","msg":"agent: history compressed","tokens_before":1236,"tokens_after":780,
 "threshold":384,"msgs_before":18,"msgs_after":10,"compressed_msgs":9,"kept_raw_rounds":3,"history_changed":true}
{"time":"…","level":"INFO","msg":"agent: history compressed",…,"history_changed":true,"task":"70583f25-…"}
```

⇒ **未打标的那一行：`task` 这枚键整个不存在**——不是 `task=""`、不是 `task="unknown"`、没有任何兜底值，
也没有被红删器改名（八枚既有键名逐字在场）。打标那一行：uuid **原样过 `redactHandler`**（键名与值都没被动过）。
⇒ 派单要的"造不出反例就说造不出"：**本程造不出**。本程试了三条造法，全部不响：
(a) 不打标 ⇒ 键不存在（上面第一行读数）；
(b) 空串／`"   "`／`"\t"`／`" \t \n "` 四枚空白 id 逐枚喂进 `withTraceTask` ⇒ **四条痕里 `task` 一枚都没出现**，
    同发里那枚真 id 照常在场（正控，防"尺瞎了"）——读数 `probes/153/accept-r1/probe-blank.txt`，探针源 `zzaccept153_blank_test.go`；
(c) 指望红删器改写它 ⇒ `Redactor` 按**键名词根**命中的三类是 `key/token/secret/password/…`、`audio/pcm/…`、`body/content/…`，
`task` 一枚都不沾；值侧的内嵌密钥形状（`sk-…`／`Bearer …`／`ghp_…`）uuid 也不匹配 ⇒ 规则层面就没有这条路。

**两种假法本程自己造出来、当场红（这才是"没兜底"的证明，不是措辞）**：

```
A8 属性无条件追加（不判空）   → 83/1，红句：
   compress_trace_test.go:552: untagged call produced a "task" attribute "" (string): absence is honest, an invented owner is not
A9 traceTaskID 恒返一枚 uuid 常量 → 83/2（两枚任务用例都红）
```

⚠ 本程在 A9 之后**忘了先 restore 就跑了一次到盘探针**，于是白捡一发读数：那枚**假常量原样落在盘上**
（`…"task":"00000000-0000-4000-8000-000000000153"…`，本程随后 restore 重跑，才有上面那两行干净读数）。
⇒ 这一发把问题的边界钉死了：**盘上没有任何一道工序会拦假归因，拦它的只有那枚用例**。
本程不因此判实现件什么（它挡的是"自己造"，不是"别人不许造"），但下一位若问"要不要在 sink 侧再加一道"，读数在这里。

### 5.2 泄漏面：这枚键今天只被读一次、也只多到一个地方

```
$ git grep -n "ctx.Value(" 6de3d1c -- '*.go' ':!*_test.go'
internal/agent/compress.go:161: s, _ := ctx.Value(traceTaskKey{}).(string)      ← 新增这一枚
internal/tools/cancel.go:49:    h, ok := ctx.Value(cancelKey{}).(cancelHandle)   ← 既有，全仓就这两枚
$ traceTaskID 的读者：除 :237 那一处之外零读者（全量 git grep 核过）；withTraceTask 唯一非测试点＝loop.go:399
```

⇒ 键是未导出类型、读侧一处；被打标的 ctx 会顺路进 `c.summarize(ctx, …)`（`compress.go:198`），
但 Go 的 `context` 没有"枚举值"的 API，下游 provider 拿不到它 ⇒ **不存在"多一个出口"**。
一句话结论：**这条隐式通道的到达面＝同一枚 `Info` 行，没有第二处。**

### 5.3 方向①：taskID 进日志文件这件事，隐私面上到底算不算账——**算，但它没算在这件证据件里**

三笔现量，逐笔给：

1. **不是新增暴露类。** 锚点 `86b0161` 上**已经有三处**把同一枚 id 用同一个键名写进同一枚文件：
   `loop.go:373 "agent: task_log open failed" … "task", taskID`、`loop.go:741 "agent: tool timeout" … "task", req.TaskID`、
   `loop.go:933 "agent: task_log finish failed" … "res.TaskID"`；到 `6de3d1c` 非测试侧带 `"task"` 的行数是 **7**。
   ⇒ 本票做的不是"让 id 第一次进日志"，是"让成功边那一行也进"。这一笔实现件用"键名沿用 `loop.go` 已经在用的 `"task"`"带过，
   **方向对、没展开**。
2. **它算过的那笔账只写在码里、没写在件里。** `compress.go:236` 那句
   "It is a correlation id, not history content, so `[privacy] keep_transcript` does not reach it"
   ——本程复算它**成立**（`keep_transcript` 硬编码 false、写 true 直接报错：`internal/config/validate.go:82-84`＋`validate_test.go:136`
   ⇒ 没有任何配置面能把它变成"转写上盘"），但它就到此为止。⚠ **本程在整份交付版证据件里搜
   `隐私`／`privacy`／`keep_transcript`／`私有数据`／`Q-31`／`SealDir`：0 命中**（`grep -c` 现量）。
   派单写"本仓已定案日志属私有数据（Q-31）"，而**该定案要求的那把锁今天还没上**：
   `winsec.SealDir` 在 `6de3d1c` 的非测试代码里**零调用者**（只有定义与注释），那正是票 132 挂着的一格。
   ⇒ 本程判：**这笔账不是票 153 欠的，但票 153 是新落进这枚"未上锁文件"的第一枚可归因键**，
   证据件里该有一行"落在票 132 的未闭合面上"；现在没有，读者会以为注释里那句"keep_transcript 不到它"＝隐私已审。
3. **命名一致性（本程自己加的一发）**：`SPEC-02:74`／`PLAN.md:2697` 里那枚叫 **`task_id`**（DB 列，`REFERENCES task_log(id)`），
   日志侧叫 **`task`**。⇒ 同一枚值两个面**故意不同名**，本票跟随日志侧是对的；
   但"从日志 join 回 DB"的人要知道这次改名，实现件与 spec 都没写这一句（登记，不判退）。

**本格档位：成立，附两笔待补**（①件内缺那笔隐私账；②日志↔DB 键名差异无人写）。
两笔都不动判据、不动码；且本程**没有**造出"这枚键泄漏到它不该去的地方"的活证据——§5.2 那把尺就是为这一条造的，造不出。

---

## 6. 攻击点⑥ —— 门禁与名册独立复跑（一律逐包，从不 `go test ./...`；版本现读）

### 6.1 四数之外比名册差集（派单点名的这条，本程照做）

| 形 | 采于（仓外快照，`git archive`） | RUN(all) | PASS 顶／子 | FAIL | SKIP | panic | 名册枚数 |
|---|---|---|---|---|---|---|---|
| 改前 `-count=1` | `86b0161` | 80 | 63／17 | **0** | **0** | 0 | 80 |
| 改后 `-count=1` | `6de3d1c` | 83 | 66／17 | **0** | **0** | 0 | 83 |
| 改后 `-count=2` | 同上 | 166 | 132／34 | **0** | **0** | 0 | 83 |

```
$ comm -13 改前名册 改后名册    # 新增，逐名
TestCompressionTraceCarriesTheOwningTaskID
TestCompressionTraceNeverInventsATaskID
TestCompressionTraceSilentWhenNothingFoldableOverThreshold
$ comm -23 改前名册 改后名册    →（空）
$ diff -q 改后-count1名册 改后-count2名册  →（空 ⇒ 两形名册全等，166＝2×83 算术自洽）
$ git grep -n "t.Skip" 6de3d1c -- internal/agent
internal/agent/approval/ticket84_no_owner_test.go:224        ← 唯一一枚，票 84 的"有意慢"闸，非本票
```

⇒ 实现件 §6.1 那三个数（80→83、新增恰三枚、丢失为空、两形名册全等）**逐条复现**，
"摘一枚证人"（§3 的 A2/A5/A7 三发 82／82／81）也都落在同一把尺上——**没有"一条 panic 吞掉整包读数"的形状**。

### 6.2 三件工具（版本本程现读，不背数）

```
$ go version        → go version go1.27.1 windows/amd64
$ gofumpt --version → v0.12.0 (go1.27.1)   （/d/work/base/gopath/bin/gofumpt.exe）
$ go vet ./internal/agent/  → 无输出 rc=0
$ gofmt -l internal/agent   → 无输出
$ gofumpt -l internal/agent → 无输出
$ go test ./internal/agent/ -count=1 -race → ok github.com/CarlosShao/wisp/internal/agent 3.894s
```

⚠ **一处读数本程复现不出来，且它是 CI 同形的那一行**（实现件 §6.2 第三条工具块里）：

```
$ gofumpt -l . tools/d22scan tools/mockllm          # 在 .scratch 之上的仓根形状＝CI ci.yml:114 那一步
.scratch/wisp/probes/152/zz152probe_windows_test.go     ← 1 枚命中，它写的是"无输出"
```

本程核对过这枚文件不是刚出现的：它由 `eb4755a`（票 152 的程，**09:36**）落进 `6de3d1c` 的祖先链，
而它这格读数取于 09:50–09:57；同一枚文件直到 `5429c0d`（票 152 自己"探针源过 gofumpt"，09:57，**与 `6de3d1c` 同一分钟**）才转干净
（逐版重量：`eb4755a` 版被点名／`6de3d1c` 版被点名／`5429c0d` 起不被点名）。
⇒ 两种解释本程都摆出来：**(i)** 它那一行跑在别的 cwd 下（`gofumpt -l .` 从 `internal/` 或快照里跑都不含 `.scratch`）；
**(ii)** 它取数那一刻票 152 已把**工作树**修好而尚未提交（同分钟竞争，无法从历史上排除）。
但**按可核对象——`6de3d1c` 那棵树——读出来是 1 枚命中**，所以这一行作为"无输出"的读数不可复现，需要更正口径。
不是判据失灵，也不是它伸手了别家地界：`gofumpt -l internal/agent`（真正属于它地界的那两条）本程复跑仍为空。
最小闭合＝把那一行改成"仓根整树 1 枚命中，属票 152 的探针，非本程地界"，或直接注明该条跑在只含 `internal/agent` 的目录。
（同族提醒：本程自己的两枚探针源在落盘前先跑过同一把尺，空 ⇒ 不当那个把 CI 格式步点红的来源。）

### 6.3 `sh scripts/d22scan.sh`：rc=0，八枚分母非零，且"别用 `grep -m1 examined`"那条坑本程复算成立

```
$ sh scripts/d22scan.sh ; echo $?            → 0        （末行 "d22scan: clean - no D22 ban violations"）
$ grep -c examined 本程这份日志              → 69       （与实现件同一枚数）
$ grep -m1 examined                          → 第 106 行：scan_test.go:1002: ban #6 examined 7 and fired in all 7 extension classes
                                                ↑ 这就是那把坏尺：取到的是自检 fixture 的 7，不是任何分母
真扫描块＝本程日志的 228–236 行（九行，与实现件说的行号相同）：
   228  examined 228 production Go files under internal/ and cmd/
   229-230  bans #1-5 internal/=205 · cmd/=23      231  ban #6 frontend/=73
   232  ban #7 internal/tools/=18                  233-236 ban #8 design/=39 · frontend/=73 · internal/=413 · cmd/=44
```

分母直接和树自身对齐（本程重量，不引用它的数）：

```
$ git ls-tree -r --name-only 6de3d1c -- internal | grep -c '\.go$' → 413  = ban #8 internal/ ✓
$ git ls-tree -r --name-only 6de3d1c -- cmd      | grep -c '\.go$' →  44  = ban #8 cmd/      ✓
$ git ls-tree -r --name-only b23c7f7 -- cmd      | grep -c '\.go$' →  43                    ← 它当年读到的就是这个 43
```

⇒ 它那句"`cmd/` 差 1 枚不是我造的、因为扫描跑在工作树上、当时 `cmd/wisp/task_scope_close_151_test.go` 还是票 151 的未跟踪文件"
**被本程独立证实**：那枚文件正是 `5d46f24`（subject 就叫 `placeholder`）带进树的，`6de3d1c` 起 `cmd/` 就是 44。
`frontend/` 从它读到的 67 涨到本程的 73 同理（前端会话在写），**别家的树在动不算进本票任何宣称**——这句它写对了。

**本格档位：成立（附一条更正）**——待补那一条就是 §6.2 里那行 `gofumpt -l .` 的口径。
