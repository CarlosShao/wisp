# 136 — AC#14 `Gate` 落地件（`136-ac14b-impl.md`）的非实现者终裁 r2 —— **acceptance**

> **这份文件是什么**：票 136 `AC#14` 的**第二张**裁决表。前一张 `136-ac14-r1-acceptance.md`（429 行／7 枚
> commit）裁的是**前一程** `worker-ticket136-ac14`（`aef82f5`）并给出附条件＝编排者批准翻布尔；本表裁的是
> **那个批准落地之后**的盘上终态（落地程 `worker-ticket136-ac14b`，证据 `136-ac14b-impl.md`，码 `52191ce`）。
> **本程身份**：`acceptor-ticket136-ac14b-r2`，**只读裁决者**，不是落地程、不是前一程实现者、不是编排者。
> 本程**不翻勾、不改码、不 revert**；发现的一切缺陷一律报回，不就地修。
> 本文件只引用变量名与文件名，不含任何凭据值（本程没有读到过任何凭据）。

---

## §0 锚点与口径

| 项 | 读数（本程自己量） |
| --- | --- |
| `git rev-parse HEAD` | **`4ecc284d1dd7b2b39f35bf54ae2222ec8d837548`**（简报说的 `4ecc284` 或其后代 ⇒ 实测正是 `4ecc284`，未漂） |
| 被验码 | `52191ce feat(136 AC#14 Gate): flip settleCoverageRowGates true + rewrite 2 coverage legs …`；`52191ce^` = `ade897c` |
| 取数范围 | `git diff 5a946d3..HEAD`（`5a946d3` ＝ 落地程 §0 自量锚点，也是编排者落批准那枚 commit） |
| 工具链 | `go version go1.27.1 windows/amd64`；`CGO_ENABLED=1`；本程自建二进制见 §5 |
| 工作树里别人的东西 | `git status --porcelain` ＝ 16 枚 `design/**` 删除 ＋ 2 枚未跟踪目录（`design/old/`、`design/doubao/`）⇒ owner 的，**本程一枚没碰、没还原、没代提交、没 stage** |
| 临时件落点 | 全部在 `D:\tmp\wisp141r2-ac14b-002e0e3e\`（`head/` `snapA/` `snapC/` `bin/`）与 `D:\tmp\ac14b-r2-136\`；**仓库目录内没有新建任何东西**；只建不删 |
| 快照怎么来的 | `git archive HEAD \| tar -x`（仓外）⇒ 三形快照：`head`＝终态、`snapA`＝只还原测试改写（gate 仍 true）、`snapC`＝`head` 再摘掉门行构造点（M1）。**全部在仓外，未在仓库内建 worktree/checkout** |

---

## §1 判据①：改动范围 —— 〔独立复现〕**成立**

### 1.1 `internal/observe/` 内的 diff 逐字节

```
$ git diff --name-only 5a946d3..HEAD -- internal/observe/
internal/observe/sampler.go
internal/observe/sampler_settle_coverage_136_test.go
$ git diff --stat 5a946d3..HEAD -- internal/observe/
 internal/observe/sampler.go                        |  2 +-
 .../observe/sampler_settle_coverage_136_test.go    | 56 ++++++++++++++++++++--
 2 files changed, 53 insertions(+), 5 deletions(-)
```

`sampler.go` 全量 diff ＝ **恰好一枚字面量、1 增 1 删**，站点 `:556`：

```
-const settleCoverageRowGates = false
+const settleCoverageRowGates = true
```

⚠ 简报里那句"`rep.Pass = foldSettlePass(...)` 那一行**若已在**就不该再动"—— **盘上确认它本来就在**：
`sampler.go:537`，`git blame -L 537,537` ⇒ `aef82f5d`（前一程 09-24 19:08）、`git log -S "foldSettlePass(memOK"`
⇒ 唯一引入枚 `aef82f5`；本范围 `5a946d3..HEAD` 对 `sampler.go` 的 diff 只有上面那一行（1 增 1 删）⇒ 落地程**没动它**，
与授权面（"只 `:556` 一枚布尔"）逐格一致。同理 `:529` `rep.Verdicts = buildSettleVerdicts(*rep)` 也 blame 到 `aef82f5d`、
本范围未触碰——这一行是 §4(C) 那一发变异的靶子。

### 1.2 冻结面 0 行

```
$ git diff 5a946d3..HEAD -- internal/observe/thresholds.go | wc -l
0
$ git diff --name-only 5a946d3..HEAD | grep -icE 'golden|testdata|thresholds'
0
```

全仓 golden/testdata 清单在本范围**逐枚 0 行**（`internal/llm/testdata/golden/*.sse`、
`internal/agent/testdata/golden/*.sse`、`internal/llm/golden/**` 一字节未动）。

### 1.3 没有别的包动过

```
$ git diff --name-only 5a946d3..HEAD | grep '\.go$' | xargs -n1 dirname | sort -u
internal/observe            ← 全范围唯一动过的 Go 包
$ git diff --name-only 5a946d3..HEAD
.scratch/wisp/issues/136-…-54-tests-stay-green-without-it.md
docs/evidence/s1/136-ac14b-impl.md
docs/evidence/s1/141-ac2-stock-inventory-r1.md
internal/observe/sampler.go
internal/observe/sampler_settle_coverage_136_test.go
```

### 1.4 逐枚 commit 的 `git show --name-only`（本范围 17 枚，每枚**只带一枚路径**）

| sha | 带的路径 | 在授权面内？ |
| --- | --- | --- |
| `52191ce` | `internal/observe/sampler.go` ＋ `sampler_settle_coverage_136_test.go` | **是**（两枚，唯一码 commit） |
| `b9ca0b0` `1e620d6` `1657135` `115173b` `92dd40f` | `docs/evidence/s1/136-ac14b-impl.md`（各一枚） | **是** |
| `a4b5deb` | `​.scratch/wisp/issues/136-….md`（票面 delivery note） | **是** |
| `ade897c` `5dfb8ba` `e0c4ada` `64664d6` `d036929` `a2d84ac` `94f4e33` `68ff486` `877f979` `4ecc284` | `docs/evidence/s1/141-ac2-stock-inventory-r1.md`（各一枚） | **否 ⇒ 见下** |

**越界旗（按简报字面清单如实报）**：本范围 17 枚里有 **10 枚**（`ade897c`、`5dfb8ba`、`e0c4ada`、`64664d6`、
`d036929`、`a2d84ac`、`94f4e33`、`68ff486`、`877f979`、`4ecc284`）落在简报给的三面授权清单**之外**——它们全是
`docs/evidence/s1/141-ac2-stock-inventory-r1.md`（票 141 盘点程 r1 的第 1—10 节）。
**判**：**不是本程要裁的越界**。简报正文自己就预告了"`docs/evidence/s1/141-*.md` 与 `design/**` 的 churn 是*别的*活"，
且这 10 枚**没有一枚碰过 `internal/`、`cmd/`、阈值、golden、票 136 的任何一面**；票 136 的写者只有
`52191ce`/`b9ca0b0`/`1e620d6`/`1657135`/`115173b`/`92dd40f`/`a4b5deb` 七枚。**结论：落地程写集与批准面逐格吻合；范围判据成立。**

**档位：〔独立复现〕**（`git diff`/`git show --name-only` 本程自己跑，不抄 §1.1—§1.4 任何转述）。

---

## §2 判据②：独立复现全绿 —— 〔独立复现〕**成立，四数与名册逐格对上**

本程在**跟踪树 `4ecc284`** 上自己跑 `go test -count=2 -v ./internal/observe/`（原始件
`D:\tmp\ac14b-r2-136\test-head.txt`）：

| 口径 | 本程实测 | 落地件 §1.3 自报 | 判 |
| --- | --- | --- | --- |
| `^=== RUN` | **142** | 142 | 对上 |
| `^--- PASS`（顶层） | **142** | 142 | 对上 |
| `^--- FAIL` | **0** | 0 | 对上 |
| `^--- SKIP` | **0** | 0 | 对上 |
| 子测试 `^    --- PASS` | **0**（顶层即全集） | — | 无嵌套稀释 |
| 去重顶层名数 | **71**（每枚恰好 ×2，`uniq -c` 只有 `2` 这一档） | 71 | 对上 |
| `grep -ci panic`（名带 Panic） | **8** | 8 | 对上 |
| 真 `^panic:` | **0** | 0 | 对上 |
| 包级行 | `ok github.com/CarlosShao/wisp/internal/observe 6.626s`，`rc=0` | — | — |

**名册是集合差集，不是只看四数**（简报明令不许只收包级 rc）：本程把自己 71 枚名册与落地件留在盘上的两枚
原始名册各做一次 `comm -3`——

```
$ comm -3 <本程 71 枚> D:\tmp\wisp141gate\logs\base.names   → 0 行
$ comm -3 <本程 71 枚> D:\tmp\wisp141gate\logs\post.names   → 0 行
```

（那两枚文件是 `PASS: TestX` 形，本程先剥前缀再比；剥后各 71 枚。）
⇒ **71 枚逐名三方守恒**：基线（gate=false、旧披露腿）、落地后（本程自己量的 `4ecc284`）、落地程自报名册，**同一份**。
无一名消失、无一名转 SKIP、无改名。

**"PASS 掉而 FAIL 不涨"这一维本程专门看了**：`--- SKIP` 两味（顶层与 `=== SKIP`）皆 **0**、`=== PAUSE` **0**、
RUN=PASS=142 三者相等 ⇒ **没有任何一枚用 SKIP 换色，也没有任何一枚被包级 `rc` 吞掉**。
本程**不以任何人的包级 rc 作判**：上表四数逐味是本程从自己的 `-v` 输出里 `grep -c` 出来的。

**档位：〔独立复现〕**。

---

## §3 判据③ 定性：那两处改写是收紧还是放水？—— 〔独立复现〕**收紧，四味新钉逐味打出"还能红"**

### 3.1 盘面事实（本程自己 `git diff` ＋ `grep -c`，不抄 §1.2 表）

| 维度 | 改前（`aef82f5`） | 改后（`4ecc284`） | 判 |
| --- | --- | --- | --- |
| `t.Skip` | **0** | **0** | 无 Skip 换色 |
| `t.Fatalf` 枚数 | **43** | **49** | **＋6**（两枚腿各 −1 ＋4） |
| `t.Error` | 0 | 0 | 无"红转警告"降级 |
| 被删的行（`git diff` 全量，逐字） | `if !rep.Pass {` ×2 ＋ 其红句 2 行（`this leg pins disclosure, not the verdict; a covered-enough window must still pass: %+v` / `report=%+v`） | — | **删掉的正是"要求这窗仍 pass"那两枚**，各有 4 枚替代钉 |
| 新增的 `||` | — | 2 处，形如 `if !Contains(Measured," errors") \|\| !Contains(Note,"sample_errors=") { t.Fatalf }` | **不是放宽**：`||` 在**否定侧**，任一缺失即红 ⇒ 两个条件**都**必须成立，比单条件更严 |
| 阈值／`thresholds.go`／golden | — | §1.2 已证 0 行 | 无放宽 |

**"有没有断言被删掉却没有'仍能失败'的替代"？——没有。** 两枚腿各得 4 枚替代钉，且本程**逐味用变异打出红**
（全在仓外快照 `D:\tmp\wisp141r2-ac14b-002e0e3e\snap*`，`-count=1`，顶层 71）：

| 本程这一发 | 改了什么（只一枚 `sampler.go` 行） | 四数 RUN/PASS/FAIL/SKIP | 两枚改写腿红不红 | 证到哪一枚钉是活的 |
| --- | --- | --- | --- | --- |
| **snapC＝M1 摘门行本体** | 删 `:529` `rep.Verdicts = buildSettleVerdicts(*rep)`；597→**596** 行；构造点 0、只剩 `:460`/`:540`/`:557` 注释＋`:563` 定义 | **71/65/6/0**，真 `^panic:`=0 | **两枚都红** | `if coverage == nil` 活 |
| **snapD＝M2 让门行不看丢读** | `:565` `covered := rep.SampleErrors == 0 && len(rep.Samples) > 0` → `covered := len(rep.Samples) > 0` | **71/67/4/0**，panic 0 | **两枚都红** | `if coverage.Pass` 活 |
| **snapE＝M4 摘"红句点名"** | `:570` 失败分支 `sample_errors=%d` → `dropped=%d` | **71/67/4/0**，panic 0 | **两枚都红**（站点 `:257`/`:324`，红句逐字 `the failing row has to print the count it failed on`） | `Measured`/`Note` 那枚 `||` 钉活 |
| **snapF＝M3a′ 撤折叠** | `:537` `rep.Pass = foldSettlePass(...)` → `rep.Pass = memOK && backInTime && releaseOK`，`:556` 保持 `true` | **71/66/5/0**，panic 0 | **两枚都红**（站点 `:261`/`:328`，红句逐字 `a failing gate row and pass=true at once: the fold rule is not applied`） | 折叠一致性那枚 `for` 循环钉活 |

⇒ **四枚新钉，四味各自被一发变异打红 ⇒ 这两枚腿"改完还能失败"是读数，不是自述。**
按简报那句"改完再也红不了的测试＝自动退回"——**不触发**。

### 3.2 与本前一程 §4 的对照：五发变异的差值恒为 ＋2

| 发 | r1 表在**未改写**树上的红数 | 本程在**终态**树上的红数 | 差 |
| --- | --- | --- | --- |
| M1 摘门行 | 4 | **6** | ＋2 |
| M2 不看丢读 | 2 | **4** | ＋2 |
| M4 摘点名 | 2 | **4** | ＋2 |
| M3a′ 撤折叠（gate=true） | 3 | **5** | ＋2 |

⇒ 两枚腿从"只核 `rep.Pass`"迁到"核门行存在＋门行自陈不 pass＋印 `sample_errors`＋折叠一致"，
**门行一旦被摘，它们就跟着红**。**没有一发红数下降**。

**档位：〔独立复现〕**（snapC/D/E/F 四发是本程自己 `cp -r` ＋ `sed` ＋ `go test`，原始件 `D:\tmp\ac14b-r2-136\snap{A,C,D,E,F}-v.txt`）。

---

## §4 判据④ 三张名册在**本程自己的仓外快照**里复现 —— 〔独立复现〕**A/B 全对，C 数变了、方向更强**

跟踪树 `4ecc284` 全程未参与变异；每发跑完回跟踪树复核 `git status --porcelain internal/observe` ⇒ **0 行**。

### (A) flip-only：只还原测试改写、门仍 `true`（快照 `snapA`）

构造：`git archive HEAD` → 仓外 `head/`，再把 `sampler_settle_coverage_136_test.go` **一字换回 `aef82f5` 版**；
`sampler.go` 保持终态。落地凭据：`grep -n 'const settleCoverageRowGates' snapA/…sampler.go` ⇒ `556: … = true`；
`grep -c 'this leg pins disclosure' snapA/…coverage_136_test.go` ⇒ **1**（旧披露腿在位）。

`go test -count=1 -v ./internal/observe/` ⇒ **RUN=71 / PASS=69 / FAIL=2 / SKIP=0**，真 `^panic:`=0，包级行 `FAIL … rc=1`。

| # | 本程实测红名 | 落地件点名的红名 | 判 |
| --- | --- | --- | --- |
| 1 | `TestCheckSettleSingleTrustworthyReadReportsItsLoss` | `…SingleTrustworthyRead…` | 同 |
| 2 | `TestCheckSettleZeroFootprintDropsAreCountedToo` | `…ZeroFootprintDrops…` | 同 |

⇒ **恰好那两枚、无第三枚转红、也没有任一枚该红的仍绿。名册与落地件 §2(A)／r1 表 §2(A) 逐名同。**

### (B) 必须仍绿的那批 —— 逐名点名（本程两味各测一遍）

| # | 用例 | `snapA`（gate=true＋旧腿） | `4ecc284 -count=2`（终态） | 判 |
| --- | --- | --- | --- | --- |
| B1 | `TestCheckSettleFullyMeasuredWindowReportsNoLoss`（健康窗） | PASS ×1 | PASS ×2 | 绿 |
| B2 | `TestCheckSettleHalfTheReadsFailedReportsItsLoss`（`:208` 已改写腿） | PASS ×1 | PASS ×2 | 绿 |
| B3 | `TestSettleCoverageRowExistsAndPassesWhenFullyMeasured` | PASS ×1 | PASS ×2 | 绿 |
| B4 | `TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed` | PASS ×1 | PASS ×2 | 绿 |
| B5 | `TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured` | PASS ×1 | PASS ×2 | 绿 |
| B6 | `TestFoldSettlePassOnlyGateRowsVeto` | PASS ×1 | PASS ×2 | 绿 |
| B7 | `TestSettleReportPassNeverContradictsItsGateRows` | PASS ×1 | PASS ×2 | 绿 |
| B8 | `TestStateReportVerdictBuilderStaysSinglePurpose`（冻结面守卫） | PASS ×1 | PASS ×2 | 绿 |

⚠ **纠简报一处数**：简报说"the **9 named** must-stay-green tests"。**盘上只有 8 枚具名**——
r1 表 §2(B) 的 `B1…B8` 是 8 个具名用例，**第 9 行 `B9` 不是名字**，它写的是"其余 **60** 枚与 settle/gate 无关的顶层用例"
（一个用排除法兜住的**集合**，`71 − 2 红 − 8 点名 = 61`… 见下）。⇒ **"9 枚具名"这句在纸上不成立，应为"8 枚具名 ＋ 1 个其余集合"**。
本程按盘上口径核：**8 枚具名全绿**（上表），且**全体 71 枚**在终态全绿（§2）、在 `snapA` 只红那 2 枚（(A) 表）
⇒ 排除法也成立：**不存在任何一枚"B 里该绿却红了"**。
顺带纠 r1 表 `B9` 自己那行算术：它写 `71 − 2 − 9 = 60`，**把非名的 `B9` 也减进名册**了；盘上无红态下的正确式是
`71 − 2 红 − 8 具名 = 61 其余`。不影响判据，但**今后引"9 枚"必须连这条口径一起引**。

### (C) 那个唯一变了的数：M1 摘门行 4 ⇒ 6 —— 本程裁定：**更强，不是越界**

本程在**终态快照**上实测（`snapC`，见 §3.1 表）红名册 **6 枚，逐名**：

1. `TestCheckSettleSingleTrustworthyReadReportsItsLoss` ← 本程核算是"改写带来的增量"
2. `TestCheckSettleHalfTheReadsFailedReportsItsLoss` ← r1 表 M1 原有
3. `TestCheckSettleZeroFootprintDropsAreCountedToo` ← 改写带来的增量
4. `TestSettleCoverageRowExistsAndPassesWhenFullyMeasured` ← r1 表 M1 原有
5. `TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed` ← r1 表 M1 原有
6. `TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured` ← r1 表 M1 原有

⇒ **与落地件 §3 自报的 6 枚逐名全等**，r1 表的 4 枚是它的**真子集**。

**裁定（简报要本程明说的那一句）**：**判"更强"，不判"越界"。** 三条理由：
① **授权面没变**——批准内容是"`:556` 一枚布尔 ＋ 按 `:208` 形状改写那 2 枚断言"，本程只动了这两处（§1 已证
`internal/observe/` 的 diff 恰两枚文件、`sampler.go` 恰一行），**没有新增用例、没有新造文件、没碰阈值**；
"红数从 4 涨到 6"是**同一次授权**的**必然副产物**——把断言从"核 `rep.Pass`"换成"核门行存在"，摘掉门行时它当然转红。
② **方向是收紧**——多出来的 2 枚红＝多出来的 2 枚钉在守卫上的消费者；§3.2 显示**五发变异无一红数下降**，
差值恒 ＋2。放水只会让红数下降，不会上升。
③ **"越界"在票上的定义是碰没授权的面**（生产码语义、阈值、别的包、别的票），那五样本程在 §1 逐样量过、**全 0**。
⇒ **4→6 这一枚数不是缺陷，是这次改写的收益**；引用它时必须连"改前树/改后树"一起引（落地件 §6.4 末行已自报此点，本程确认）。

**档位：〔独立复现〕**。

---
