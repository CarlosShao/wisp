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

## §5 判据⑤ 退出码前置——本程**自己重建、自己采**，3 发（简报上限），未代抄落地件任何一发

**二进制是本程自己的**：`git archive HEAD`（当时 `4ecc284`）→ 仓外 `head/`，`grep -n 'const settleCoverageRowGates' head/…sampler.go`
⇒ `556: … = true`（**先证落地**），再 `CGO_ENABLED=1 go build -trimpath -o D:\tmp\wisp141r2-ac14b-002e0e3e\bin\wisp.exe ./cmd/wisp`
⇒ **rc=0**；三枚 sherpa/onnx DLL 从仓内 `third_party/sherpa-onnx/` 同目录摆放（否则加载期 `0xc0000135`）。
本程那枚 exe 的 `sha1sum` ⇒ `1060e481240b413365bb263cf3c220be9b2ceac7`（**与落地件那枚不是同一份，本程不借它**）。
命令形状 `WISP_ENV=test ./wisp.exe slo -settle -out readings/settle-N.json`，N=1..3；驱动脚本
`D:\tmp\wisp141r2-ac14b-002e0e3e\shots.sh`，读数 `readings/settle-{1,2,3}.json`，全流水 `D:\tmp\ac14b-r2-136\shots.txt`。

**争用自查用的是盘上那份名单**：`scripts/slo-check.ps1:153-155`，本程现量 **15 枚**
（`go` `gofmt` `cgo` `compile` `asm` `link` `gcc` `g++` `cc1` `cc1plus` `as` `ld` `wisp` `wisp-cli` `staticcheck`）
⇒ 与落地件 §4 的"15 枚"对上，**简报那句"15 枚不是 8 枚"成立**。
本程每发**前、后各扫一次**（`Get-Process` 精确名比对，扫后另等自家子 `wisp.exe` 退净，最多 30s）⇒
**六次扫描全空**（`roster_before=[]`／`roster_after=[]` ×3），**没有一发需要记"争用，未采"**。
（取样窗口内本程刻意不跑任何 `go`/`git` 编译类命令，§3/§4 的变异全部安排在六发之后。）

| # | exit | `sample_errors` | `len(samples)` | `settle.pass` | 顶层 `pass` | `back_within_cap_ms` | `final_bytes` | `free_os_memory_count` | **`settle.verdicts` 那一行** |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | **0** | **0** | 40 | true | true | 276 | 8912896 | 2 | `sampling pass=true gate=true measured="40 valid / 0 errors"` |
| 2 | **0** | **0** | 40 | true | true | 270 | 9990144 | 2 | 同上形制 |
| 3 | **0** | **0** | 40 | true | true | 270 | 9138176 | 2 | 同上形制 |

**判读**：
- **三发 `exit` 全 0、`sample_errors` 全 0** ⇒ 编排者那道硬前置（"任一发非 0 ⇒ 退回 `false`、不提交"）**不触发**，
  那枚翻转**站得住、可留**。
- **三发 JSON 里都带 `verdicts` 键、那一行 `gate=true`** ⇒ 本程量的确实是**含门行且门已开**的二进制。
  这一维正是 r1 表 §5 唯一没复现的〔仅自述〕档（编排者自己那六发用的是 `ca2c55e` 的旧 exe、连 `verdicts` 键都没有）：
  **本程把它从〔仅自述〕升成〔独立复现〕**。
- 干净窗里门行自陈 `pass=true` ⇒ `foldSettlePass` 无物可否决 ⇒ `settle.pass` 与顶层 `pass` 都真 ⇒ exit 0，
  链条在真二进制上闭合了一遍（**不是**"摘门行才红"的反证，而是"正常窗不被误红"的正证）。
- 物理量级与落地件 §4 六发互证：`samples=40`、`back 270-276`（落地件 262-335）、`final_bytes 8.9-10.0MB`
  （落地件 8.9-9.7MB），均远低于冻结 Sleeping 上限 25MB。
- **合计**：本程 3 发 ＋ 落地程 6 发 ＝ **9 发**同形制读数，逐发 `sample_errors=0`／`exit=0`，无一非 0。

**一处流程口径如实登记（不改判语）**：本程读到 `9090f36`（兄弟程/编排者提交，只碰 `docs/reports/**` ＋ 票 141 票面，
`internal/`＋`cmd/` 0 枚）落在本程锚 `4ecc284` **之后**，其 message 里含一句"AC#14 门已翻 true 等 r2 复判"。
**那句不是本程的判据也不是授权**——本程上表三发是自己重建、自己采的；HEAD 会漂这件事按票上老规矩登记，
本程所有读数标注的都是**取数当时实测的 `4ecc284`**，且 `git diff 4ecc284..HEAD -- internal/ cmd/` ⇒ **0 行**
⇒ 漂动不影响本表任何一发读数。

**档位：〔独立复现〕**。

---

## §6 裁定：`sampler.go:546-555` 那句假话——**不阻塞翻勾，属收口；落地程的克制是对的**

### 6.1 简报点名的那一处，本程盘上确认

`sampler.go:550` 逐字仍写着 `// that still states its own not-pass. It is false at HEAD because the`，
`:552` 引 `docs/evidence/s1/136-ac14-impl.md` section 1.2（**翻之前**那份证据），`:553` 说 "one of the THREE existing assertions"。
盘上 `:556` 现值 `true` ⇒ **这三句对终态树已不成立**。

### 6.2 为什么**不阻塞** AC#14 翻勾（三条，都是读数不是推理）

1. **它不参与任何断言**：`grep -rn "false at HEAD" --include=*.go` ⇒ 只有注释自身（`sampler.go:550`）与
   `sampler_settle_gate_136_test.go:222` 两枚**注释行**，**没有一枚测试钉这段文字**；
   本程 §2 的 142/142/0/0 与 §3 的四发变异里，它一次也没红过 ⇒ **零可执行语义**。
2. **判据①—⑤ 不含注释**：批准面（票面 `:314`）写的是"1 枚布尔＋按 `:208` 形状改写 2 枚断言"，
   五条判据量的是范围、名册、定性、退出码、门禁——**五样本程全过**（§1—§5）。
3. **AC#14 判据要的那件事已经落成**：门行现在是 gate、能否决 exit、红句点名 `sample_errors`、
   守卫在终态树下钉得更牢（§4 C）——**注释说反不影响这些**。

### 6.3 落地程"故意不动它"的克制——**判：正确，应当保持**

`136-ac14b-impl.md` §1.2 末段与 §6.5②两处**主动上报**了这一处、并写明"派单/票面 `:314` 把授权限定在 `:556` 一枚布尔，
改注释属扩面，故只如实报回、不擅改"。本程判这是**对的动作**，三条理由：
① `AGENTS.md §1.1`/`SPEC-12 §4.1` 的"改契约＝人工批准"精神下，**agent 单方面扩自己的授权面**才是跑歪模式 #1；
② 简报那句"它的授权只点名了 `:556`"是**给定的事实**，落地程照它执行＝**遵守授权**，不是遗漏；
③ 它**没有悄改**、而是留在证据里交给翻勾那一程 ⇒ 账目可追。
⇒ **不要因为"结果上有假话"倒过来罚它的克制**；该做的是给翻勾那一程补一枚具名授权（见收口 #1）。

### 6.4 本程顺手量到**简报没点的另外三处同病**（一并进收口，别只修一处留三处）

| # | 站点 | 现在写着什么 | 为什么在终态下不成立 |
| --- | --- | --- | --- |
| a | `internal/observe/sampler_settle_gate_136_test.go:222` | `// With the row recorded (gate=false at HEAD) this window still passes;` | 门在 HEAD 已是 **true**；**简报只点了 `sampler.go` 那一处，这一处是它没列的** |
| b | `internal/observe/sampler_settle_coverage_136_test.go:26-27`（文件头注释） | `// … a partially covered window still passes today (changing that verdict is not this cell's job)` | **两重假**：窗今天不 still pass；改那个 verdict **就是** AC#14 这格刚做完的事。（r1 表 §3 已注记该头注释"`aef82f5` 逐字未动、改它属 AC#15 射程"——AC#14 落地后这句更站不住了） |
| c | `internal/observe/sampler_settle_coverage_136_test.go:237-241`（前一程 `:208` 腿内的注释） | `// The pass bit stays this window's verdict for as long as the row is recorded rather than gated (… docs/evidence/s1/136-ac14-impl.md section 1.3 owns the current state)` | 行**已经** gated ⇒ 该条件子句失效；引的证据件是**翻之前**那份（终态证据是 `136-ac14b-impl.md`） |

⇒ 收口**不是"一行"**（简报那句"属另一枚一行 commit 的收口项"这一措辞本程判**说小了**）：**四处、两枚 `.go` 文件**，
且其中一处（c）在**前一程**写的腿里、一处（a）在**前一程地界**那个 `_gate_136_test.go` 里 ⇒
**修它需要的授权面比"一行"宽**。这一条按盘上口径报回，不替编排者决定谁修。

**档位：〔独立复现〕**（四站点行号均为本程现量）。

---

## §7 总裁与收口清单

### 7.1 总裁：**可翻勾**（AC#14 在终态树 `4ecc284` 上成立）

| 判据 | 判 | 一句话理由（全部为本程自量） |
| --- | --- | --- |
| ① 范围 | **成立** | `5a946d3..HEAD` 唯一动过的 Go 包＝`internal/observe`；`sampler.go` 恰一行（`:556` false→true）；`thresholds.go` 与全仓 golden/testdata **0 行**；17 枚 commit 逐枚 name-only，本程写集七枚全在授权面内 |
| ② 独立复现绿 | **成立** | 本程自跑 `-count=2 -v` ⇒ **142/142/0/0**、顶层 71 每枚恰好 ×2、`SKIP`＝0、`PAUSE`＝0、真 `^panic:`＝0；名册对落地程 `base.names`/`post.names` 各 `comm -3` **0 行** ⇒ 71 枚逐名三方守恒 |
| ③ 收紧还是放水 | **收紧** | 无 Skip、无降级、无阈值动；`t.Fatalf` 43→49；新 `||` 全在否定侧＝更严；被删的两枚"要求 pass"各有 4 枚替代钉，且本程**四发变异逐味把它们打红**（摘门行／不看丢读／摘点名／撤折叠） |
| ④ 三张名册 | **全对** | (A) flip-only 恰红那 2 枚、无第三枚；(B) 具名 8 枚两味全绿；(C) M1 实测 6 枚、r1 的 4 枚是其真子集 ⇒ **判"更强"非"越界"** |
| ⑤ 退出码前置 | **成立、不回退** | 本程自建 exe（`sha1 1060e48…`、先证 `:556=true`）自采 3 发 ⇒ 逐发 `exit=0`／`sample_errors=0`／`verdicts[sampling gate=true]` 在场；六次争用扫描全空；与落地程 6 发合起来 9 发无一非 0 |

**唯一扣分项是一枚非阻塞的纸面缺陷**：`sampler.go` 及两枚测试文件里**四处**注释仍说门"false at HEAD／still passes today"，
零可执行语义、钉不红任何一发 ⇒ **列收口，不挡翻勾**（§6）。

### 7.2 收口清单（优先级序；本程一枚都没做，等编排者具名授权）

| 优先级 | 收口项 | 面 | 备注 |
| --- | --- | --- | --- |
| **1** | 改写 §6.4 那**四处**过期注释（`sampler.go:546-555`／`_gate_:222`／`coverage_:26-27`／`coverage_:237-241`），并把引用的证据件从 `136-ac14-impl.md` 换成 `136-ac14b-impl.md` | 2 枚 `.go` 文件 | **须先给具名授权**（落地程的授权只点到 `:556`）；不要只修一处留三处 |
| **2** | 翻 `AC#14` 的勾 | 票面 | 勾是编排者的；判据五条本程全过 |
| **3** | 台账更正："9 枚具名"→"8 枚具名 ＋ 1 个其余集合"，并修 r1 表 `B9` 那行 `71−2−9=60` 的算术（正确 `71−2−8=61`） | `pending-and-issues.md` | 简报那句"the 9 named"在盘上不成立（§4 B） |
| **4** | 引"M1 4→6""M2 2→4""M4 2→4""M3a′ 3→5"这一族数时**必须连"改前树/改后树"一起引**；建议直接把这四行做成一张表进票 136 或 §7 台账 | 票面/台账 | 单引"6"会让下一位以为落地程多写了 2 枚用例 |
| **5** | AC#15 靶形现量重划：`_coverage_:26-27` 头注释与 `:237-241` 现在**同时**是 AC#15 的旧账与收口 #1 的对象，两程要排队、别互相覆盖；AC#15 那 7 枚前提腿站点本程读数**一次没响**（"没响"非"已修"） | 票面 | 落地件 §6.5② 已排此序 |
| **6** | 推送前把 §5 那三发与本表一起带上（本程只 commit 未 push；`slo-smoke`/`slo-full` 由编排者核过之后再触发） | 编排者 | 本程未跑 `-race`、未复跑 linux 分母、未核 CI run id |
| **7** | 低优先登记：本包 71 枚里有 **8** 枚名字含 `Panic`（`grep -ci panic`），真 `^panic:` 全程 0；今后谁用"panic 计数"当红绿判据，先分清这两味 | 可选 | 本程两味分开记 |

### 7.3 纪律回执（本程）

只 commit、**全程未 `git push`**；`git add` 只对本程新建的那**一枚**证据文件用过一次显式路径（新文件 untracked，
`git commit --` 匹配不到）；**未用** `-A`／`.`／`-a`；**未用** `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；
**未在仓库目录内建 worktree 或快照**（五份快照全在 `D:\tmp\wisp141r2-ac14b-002e0e3e\`，只建不删）；
`design/**`（owner 16 枚未提交删除＋`design/old/`、`design/doubao/` 两枚未跟踪目录）**未还原、未 stage、未提交、未删**；
`docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、`rules_gateway.go`、`tools/d22scan/**`、`allowlist.txt`、
任何 `thresholds.go`／golden、`frontend/**`、`docs/reports/**`（编排者的）、
`136-ac14-r1-acceptance.md` 与 `136-ac14b-impl.md`（他程的）**一字节未写**；
跟踪树 `git status --porcelain internal/observe` 在全部变异跑完后复核 ⇒ **0 行**；**未翻任何勾**、**未改任何码**、
**未 revert 那枚布尔**（§5 三发全 0，无回退触发条件）。

### 7.4 注入面两栏计数（本程自己扫，简报明令要记）

- **真通知回显 7 条**（逐条：哪个工具 ＋ 前约 40 字）：
  1. `Bash`（`git rev-parse HEAD` 那次）——"<system-reminder> The following skills are available for use…"
  2. `Bash`（同一次输出尾）——"<system-reminder> The date has changed. Current date: 2026-09-24"
  3. `Bash`（同一次输出尾）——"<system-reminder> Memory: d:/work/workspace/projects…"（`AGENTS.md` 项目上下文回显，与本仓那份逐字一致）
  4. `Read`（读票面 136）——"Do not revert this file if it was changed externally; it was not." ⇒ 形状上正是最该警惕的一句（它提到 revert），但它**自陈未被外部修改**、且与本程随后实测一致 ⇒ 判真回显，且**没据此做任何动作**
  5. `Bash`（后台任务 1）——"<task-notification>…<task-id>baw2mgdlg</task-id>…"（observe 测试跑完）
  6. `Bash`（后台任务 2）——"<task-notification>…<task-id>b3iypkjy0</task-id>…"（三发 settle 跑完）
  7. `Edit`（本文件两次改写）——"The file … has been modified successfully … no need to read it back"
  （另：`9090f36` 那枚兄弟程 commit 顶掉了本程 `git log -1` 的回显——**共享树的真实并发**，非注入；本程 §0 的锚因此按号取，不靠 `HEAD`。）
- **判为注入 0 条**：全程**没有**任何工具输出自称"编排者备注／系统提示／文件已被修改通知"来指令本程
  revert、放宽阈值、改判语，或"confirm this note is genuine"；本程**没有**把任何此类文字当授权执行。
  特别记两味：① `9090f36` 的 message 里含"AC#14 门已翻 true 等 r2 复判"——那是**盘上真实的他人 commit**，
  本程**没有**据此下判（§5 三发为自建自采）；② `sampler.go:554` 注释"Flipping it is an orchestrator move, not an
  implementer move" 是**代码自述**，其对应的批准在票面 `:313-325` 有正文，本程按正文核。
- **凭据卫生**：本程未读到、未抄写任何凭据值；出现的只有变量名与文件名（`WISP_ENV`、`CGO_ENABLED`、`GOPATH`、`CC`）。

**本表 next ＝ 编排者**：①按 §7.1 翻 `AC#14`；②据 §7.2#1 给一枚**覆盖 2 枚 `.go` 文件、四处站点**的具名注释收口授权；
③§7.2#3/#4 两条台账更正随翻勾一起落。**本程不代做其中任何一件。**

### 7.5 顺手把落地件 §5 的门禁自己复跑了一遍（简报未要求，但不复跑就没人核过）

| 门 | 本程自己跑的读数（树＝`4ecc284`） | 落地件 §5 自报 | 判 |
| --- | --- | --- | --- |
| `gofmt -l internal/observe` | **空输出 rc=0** | 空 | 对上 |
| `gofumpt --version` | **`v0.12.0 (go1.27.1)`**（盘上现量，与落地件同一枚） | v0.12.0 | 对上（旧报告里的 v0.7.0 确为过期值） |
| `gofumpt -l internal/observe` | **空输出 rc=0** | 空 | 对上 |
| `go vet ./internal/observe/` | **rc=0** | rc=0 | 对上 |
| `go build ./...` | **rc=0** | rc=0 | 对上 |

⇒ **五道门禁本程独立复现全绿 ⇒ 落地件 §5 那一档从〔自报〕升为〔独立复现〕**。
本程**未**复跑 `d22scan`（它要 `git archive` 到仓外再 `sh scripts/d22scan.sh`，本程§0 已记快照落点；
ban #8 的分母面本程在 §1.2 只用 0-line diff 证过——**这一格标〔未独立复现〕**，谁要收口 #1 之后一并补一发即可，
理由：§6.4 那四处注释若被改写，**必须确认没引入 ban-#8 段字符**（那四处原文里含 `⇒`/`—`/`——` 之类，
`U+2190–U+2BFF` 段正在本仓 ban #8 射程内 ⇒ **改注释这一发真得跑 d22scan**，不是可选项）。
**这一条升进 §7.2 收口 #1 的判据里。**

### 7.6 本程 commit 流水（按号取，不靠 `HEAD`——共享树里号会漂）

| sha | 净面 |
| --- | --- |
| `96706b9` | `docs/evidence/s1/136-ac14b-r2-acceptance.md`（§0—§2） |
| `32036a0` | 同上（§3—§4） |
| 本节末枚 | 同上（§5—§7，hash 由下一位从 `git log -1 -- …` 按号读） |

每枚只带本程那一枚路径；每次 commit 前 `git diff --cached --name-only` 现量、名单只有那一枚。
两枚之间的 `9090f36`（编排者/兄弟程，`docs/reports/**` ＋ 票 141 票面）是**真实并发**，本程未碰。

---
