# 136 — AC#14 `Gate` 那一支的落地件（翻 `settleCoverageRowGates` + 按 `:208` 形状改写两枚断言）**impl**

> **这份文件是什么**：票 136 `AC#14` 的 **`Gate` 那一支**（编排者 09-24 19:5x 批准的落地程，票面
> `.scratch/wisp/issues/136-...md:313-325` 那枚 `>` 块＝本程唯一授权来源与判据清单）的**实现件**。
> 前一程 `worker-ticket136-ac14`（码 `aef82f5`）＋ 非实现者终裁（`136-ac14-r1-acceptance.md`，429 行）已把
> "翻完之后必须仍绿的有哪几枚"这一半量出来；本程把那枚批准落到盘上。
> **本程不翻任何勾**（`AC#14` 的勾归编排者在终裁落地后翻）。
> **批准内容（票面 `:314` 逐字）**：动 1 枚布尔（`sampler.go:556` `settleCoverageRowGates` false→true）
> ＋ 按 `:208` 已建立的形状改写另 2 枚断言（`sampler_settle_coverage_136_test.go` 现盘 `:143-145`、`:280-282`）。
> **本文件只引用变量名与文件名，不含任何凭据值**（本程没有读到过任何凭据）。

---

## §0 锚点、口径、被验码关系

| 项 | 读数（本程自己量） | 取数时刻（本机） |
| --- | --- | --- |
| 开工自量锚点 | `git rev-parse --short HEAD` ⇒ **`5a946d3`**（分支 `dev`；简报写的 `30e19ef` 是编排者落笔时的读数，已漂） | 开工 |
| 被验码与 HEAD 的关系 | `git diff --stat aef82f5..HEAD -- internal/observe/` ⇒ **空** ⇒ 在 HEAD `5a946d3` 上量 `./internal/observe/` ＝ 在 `aef82f5` 上量（逐字节同码）；开工前 `git status --porcelain internal/observe` ⇒ 0 行 | 取基线前 |
| 两枚被改文件的 `git log -1` | `internal/observe/sampler.go` 与 `sampler_settle_coverage_136_test.go` 末次改动均为 **`aef82f5`（09-24 19:08）**；本程在其上落 `:556` 一枚字面量 ＋ 那两处断言 | 取基线前 |
| 工具链 | `go version go1.27.1 windows/amd64`；`go env CGO_ENABLED`=**1**，`CC`=**gcc**；`gofumpt` 盘上现量见 §5（不抄任何人的版本号） | 取基线前 |
| 工作树里别人的东西 | `git status --porcelain` 除本程两枚 `internal/observe/**` ＋ 本证据文件外，只有 owner 的 `design/**` 未提交删除与两枚未跟踪目录 ⇒ **一枚没碰、没还原、没代提交** | 全程 |
| 临时件落点 | 全部在 `D:\tmp\wisp141gate\`（`bin/wisp.exe` + 三枚 sherpa/onnx DLL · `readings/` · `logs/`）；**仓库目录内没有新建任何东西**；`build/wisp.exe`（09-21 参照件）不覆盖 | 全程 |

**基线四数（改动之前，`go test -count=2 -v ./internal/observe/`）**：RUN=**142** / PASS=**142** / FAIL=**0** / SKIP=**0**，
去重顶层 **71** 枚（×2 乘子），`grep -ci panic`=**8**、真 `^panic:`=**0**。
⇒ 与实现件 §0.1（65×2 基线＋本前一程净增 6=71）与终裁表 §1（142/142/0/0、71、panic 8/0）**逐格对上**；
本程自带一棵基线名册（`D:\tmp\wisp141gate\logs\base.names`）用于 §1 的守恒差集，不靠任何人的转述。
凭据文件：`D:\tmp\wisp141gate\logs\base-v.txt`。

---

## §1 判据① 翻布尔前 + 先证落地再读数；判据② 改后四数与名册守恒

### 1.1 判据①：翻之前的那一行 ＋ 基线构建（票面 `:324` 判据①）

**翻之前 `grep -n` 出被改那行**（证明确实是从 `false` 起步）：

```
$ grep -n "const settleCoverageRowGates" internal/observe/sampler.go
556:const settleCoverageRowGates = false
```

**基线 `go build ./...`**（未改树）⇒ **rc=0**（`CC=gcc`、`CGO_ENABLED=1`）。

### 1.2 三处改动落地凭据（先证落地，才读 §1.3 的数）

本程写集**恰好两枚文件**（`git diff --numstat` 现量）：

```
1	1	internal/observe/sampler.go
52	4	internal/observe/sampler_settle_coverage_136_test.go
```

**改动 A：翻一枚布尔（只 `:556` 一行，`git diff internal/observe/sampler.go` 全量在下面）**——

```
@@ -553,7 +553,7 @@ const settleCoverageMetric = "sampling"
-const settleCoverageRowGates = false
+const settleCoverageRowGates = true
```

改后盘上现量：`grep -n "const settleCoverageRowGates" internal/observe/sampler.go` ⇒ `556:const settleCoverageRowGates = true`。
⚠ **登记一处会腐坏的注释（本程按"只 `:556` 一行"的授权面**没有**去改它）**：紧邻常量上方的注释块
（`sampler.go:546-555`）仍逐字写着 "It is false at HEAD because ... the ticket pre-authorised one of the THREE existing
assertions"——翻完之后这句对盘上的值已不成立。派单/票面 `:314` 把本程授权**限定在 `:556` 一枚布尔**，改注释属扩面，
故本程只把这一处**如实报回编排者**、不擅改（那一枚布尔值本身与相邻注释的口径要一并收口，交给翻勾那一程）。

**改动 B：按 `:208` 已建立形状改写两枚断言（逐字点名，票面 AC#14 判据①/② 预先授权）**——

| 枚（改后站点） | 原句（被本程改掉，逐字） | 改后钉什么（与 `:208` 同形、gate 无关） |
| --- | --- | --- |
| `sampler_settle_coverage_136_test.go:143-145`（`TestCheckSettleSingleTrustworthyReadReportsItsLoss`） | `if !rep.Pass {` ＋ `t.Fatalf("this leg pins disclosure, not the verdict; a covered-enough window must still pass: %+v", rep)` ＋ `}` | 门行必须存在（否则红 "must carry its own sampling verdict row"）／门行必须**自己**说 not-pass（`coverage.Pass` 为真即红）／红句必须含 `sample_errors=` 且 `Measured` 含 " errors"／失败 gate 行与 `pass=true` 不得并存 |
| `:280-282`（`TestCheckSettleZeroFootprintDropsAreCountedToo`） | `if !rep.Pass {` ＋ `t.Fatalf("report=%+v", rep)` ＋ `}` | 同上四道（"zero-footprint-drops settle window" 那一枚措辞） |

⚠ **是"换成会自己判失败"，不是删掉、不是放宽、不是 `t.Skip`**：两枚都从"要求这窗仍 pass"改成
"要求门行自己说不 pass 并点名 `sample_errors`"，断言强度**只增不减**（各新增 4 道钉）。改法逐字复刻前一程在
`:208-236`（`TestCheckSettleHalfTheReadsFailedReportsItsLoss`）已建立并被终裁表 §2(B) 验证过"两味取值都成立"的
gate 无关形状——翻不翻布尔它都成立，所以 §1.3 的全绿不是靠放宽换来的。
**没有碰的第三枚 `if !rep.Pass`**：`TestCheckSettleFullyMeasuredWindowReportsNoLoss`（现 `:355`，
红句 `a measurable, settled, released window must pass`）——那是终裁表 (B) 里"健康窗口翻 Gate 后必须仍绿"那一枚，
本程**一字未动**（它的窗口 `sample_errors=0` ⇒ 门行自陈 pass ⇒ 折叠无从否决 ⇒ 翻布尔后照旧绿）。

**落地后 `go build ./internal/observe/` rc=0 ＋ `go build ./...` rc=0**（三条都在，才读下面的数）。

**代码 commit 回显**（原样贴，供核"每枚 commit 只带本程自己那两枚路径"）：

```
$ git log --oneline -1
52191ce feat(136 AC#14 Gate): flip settleCoverageRowGates true + rewrite 2 coverage legs to gate-independent verdict-row shape
$ git show --name-only HEAD
internal/observe/sampler.go
internal/observe/sampler_settle_coverage_136_test.go
```

### 1.3 判据②：改后 `go test -count=2 -v ./internal/observe/` 四数 + 名册守恒

| 口径 | 本程改后一发 | 基线（§0） | 终裁表 §1 |
| --- | --- | --- | --- |
| `^=== RUN` | **142** | 142 | 142 |
| `^--- PASS` | **142** | 142 | 142 |
| `^--- FAIL` | **0** | 0 | 0 |
| `^--- SKIP` | **0** | 0 | 0 |
| 去重顶层名数 | **71**（×2＝142） | 71 | 71 |
| `grep -ci panic`（名带 Panic） | **8** | 8 | 8 |
| 真 `^panic:` | **0** | 0 | 0 |

**名册差集（除四数外那一维）**：基线名册（71 行）与改后名册（71 行）
`comm -3 base.names post.names` ⇒ **输出 0 行** ⇒ 71 枚逐名**完全守恒**，无一名消失、无一名转 SKIP、无改名。
⇒ **判据②「必须仍 71 守恒、0 红 0 跳」成立**。翻布尔＋改写两枚断言后，整包全绿；那两枚断言是**换成会自己判失败**，
不是被删或被 Skip 换色。`-count=2` 的 ×2 乘子写清：顶层 71 枚 × 2 轮 = 142 行 RUN/PASS。
凭据文件：`D:\tmp\wisp141gate\logs\post-v.txt`。

**档位**：判据① 与 判据② 本程**自量**（基线名册、改后名册、四数、panic 两味全出自本程在跟踪树 `5a946d3` 上自己跑的读数）。

---

## §2 两张名册复算（翻布尔后红的确切是哪几枚 / (B) 里有没有意外变红）

**为什么要复算**：本程改写的是终裁表 §2(A) 点名那两枚；要证"改对了、不多不少"，最硬的一发是把
**前一模的状态**（gate=false、那两枚断言**还没改写**的树）单独翻布尔，看会红的恰好是不是那两枚。
快照 `git archive aef82f5`（前一程实现件那棵树，`:556=false`、`:143`/`:280` 仍是旧披露腿）→
仓外 `D:\tmp\wisp141gate\snap-pre`，**只 `sed -i '556s/false/true/'`，两枚断言一字未动**。

**先证落地再读数**：`grep -n "const settleCoverageRowGates"` ⇒ `556: ... = true`；`go build ./...` rc=0；
`go test -count=1 -v ./internal/observe/`（顶层 71，不带 ×2 乘子）⇒ **RUN=71 / PASS=69 / FAIL=2 / SKIP=0**。

**(A) 红名册恰好 2 枚（与终裁表 §2(A) 逐名同）**：

| # | 红名（逐名） | 站点 |
| --- | --- | --- |
| A1 | `TestCheckSettleSingleTrustworthyReadReportsItsLoss` | `sampler_settle_coverage_136_test.go:144` |
| A2 | `TestCheckSettleZeroFootprintDropsAreCountedToo` | `sampler_settle_coverage_136_test.go:281` |

⇒ **本程改写的正是这两枚**（§1.2 的两处），一枚不多一枚不少。**不是"红了一片"**。

**(B) 点名必须仍绿的 9 枚 —— 本程实测全部仍绿**（FAIL 总数＝2 且都是 (A)，按排除法 (B) 全体绿；
两枚关键再逐名点名）：`TestCheckSettleFullyMeasuredWindowReportsNoLoss`（健康窗，翻 Gate 后必须仍绿，实测 `--- PASS`）、
`TestCheckSettleHalfTheReadsFailedReportsItsLoss`（`:208` 已改写腿，实测 `--- PASS`）。另 60 枚与 6 枚 settle/gate 无关的顶层
用例照常全绿，`SKIP=0`、真 `^panic:=0`。**没有任何 (B) 里的枚意外变红。**
凭据文件：`D:\tmp\wisp141gate\logs\fliponly-v.txt`。**档位：〔本程自量〕。**

---

## §3 判据③ `M1` 摘门行 ⇒ 钉子必须转红（复核终裁表自报的"4 枚"）

**这一发打在**本程**改后的树**（`git archive HEAD`＝`52191ce` → 仓外 `snap-m1`；`:556` 已是 true、两枚断言已改写），
落终裁表 §4 同一发变异：删掉 `sampler.go:529` 那行 `rep.Verdicts = buildSettleVerdicts(*rep)`（门行构造点）。

**先证落地再读数**：改后 `grep -n buildSettleVerdicts internal/observe/sampler.go` ⇒ 只剩 `:460`/`:540`/`:557`
三处注释 ＋ `:563` 函数定义，**构造点 0**；文件 597→**596** 行；`go build ./...` rc=0。
`go test -count=1 -v ./internal/observe/` ⇒ **RUN=71 / PASS=65 / FAIL=6 / SKIP=0**，顶层 71 守恒、真 `^panic:=0`。

**M1 红名册（6 枚，逐名）**：

| # | 红名 | 是不是本程新引入的消费者 |
| --- | --- | --- |
| 1 | `TestCheckSettleSingleTrustworthyReadReportsItsLoss` | **是**——本程 §1.2 改写的两枚之一 |
| 2 | `TestCheckSettleHalfTheReadsFailedReportsItsLoss` | 否（前一程 `:208` 腿） |
| 3 | `TestCheckSettleZeroFootprintDropsAreCountedToo` | **是**——本程 §1.2 改写的两枚之一 |
| 4 | `TestSettleCoverageRowExistsAndPassesWhenFullyMeasured` | 否（前一程 gate 探针） |
| 5 | `TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed` | 否（前一程 gate 探针） |
| 6 | `TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured` | 否（前一程 gate 探针） |

**复核结论**：终裁表 §4-M1 自报红 **4 枚**（= 本表第 2/4/5/6 枚），**本程改后为 6 枚**——差的两枚正是
本程 §1.2 改写的那两枚断言。**这是更强不是更弱**：那两枚从"只核 `rep.Pass`"改成"核门行存在 + 门行自陈不 pass"后，
把门行摘掉时它们也会转红（`coverage == nil ⇒ t.Fatalf`）。判据③"摘门行⇒钉子必须转红"**成立**，
且本程的改写让"钉子的覆盖面"从 4 增到 6。健康窗 `TestCheckSettleFullyMeasuredWindowReportsNoLoss` **未被扫进红名**
（它不核行、只核 `rep.Pass`，而 nil 门行经 `foldSettlePass` 不否决 ⇒ 仍绿）。凭据：`D:\tmp\wisp141gate\logs\m1-v.txt`。
**档位：〔本程自量〕。**

---

## §4 判据④（本程真正的产出）：用**含门行的新 exe** 跑 ≥6 发 `wisp slo -settle`

**为什么必须重建再量**：编排者 19:5x 那六发复现用的 `D:\tmp\wisp141gate\bin\wisp.exe` 前身是从 `ca2c55e` 编的
（早于 AC#14 的码），它的 JSON **连 `verdicts` 键都没有** ⇒ 那六发只答了"真窗口会不会丢采样"这个物理问题，
**没答"带闸的终态二进制会不会变红"**。本程从**改后的 HEAD**（`52191ce` 的码，锚在 `1e620d6`）重建二进制：
`CGO_ENABLED=1 go build -trimpath -o D:\tmp\wisp141gate\bin\wisp.exe ./cmd/wisp` ⇒ rc=0，
三枚 sherpa/onnx DLL 同目录（否则加载期 `0xc0000135`）。

**命令形状（照实现件 §1.2）**：`WISP_ENV=test` ＋ `wisp.exe slo -settle -out readings/settle-N.json`（N=1..6）。
CLI 出线是**信封**：顶层 `pass`/`mode`/`settle`，`SettleReport` 嵌在 `.settle` 下（`slo-check.ps1` 读 `.settle.free_os_memory_count` 同一形制）。

**争用自查**：驱动脚本 `D:\tmp\wisp141gate\slo-shots.ps1` 每发前后各扫一次 `slo-check.ps1:153-155` 那份
**盘上现量 15 枚**名单（`go`/`gofmt`/`cgo`/`compile`/`asm`/`link`/`gcc`/`g++`/`cc1`/`cc1plus`/`as`/`ld`/`wisp`/`wisp-cli`/`staticcheck`），
且扫前先等自家子 `wisp.exe` 退净（最多 30s）⇒ **六发 `roster_before=[]` / `roster_after=[]` 全 12 次空、无一发 ABORT**。
本程**不 push**（`slo-full` 是 `push`/schedule 触发，编排者这轮不推）⇒ 取数窗内无并发取样。

**逐发读数（含 exit code，逐枚取自 `readings/settle-N.json` 的 `.settle`）**：

| # | 起时(本机) | exit | `sample_errors` | `len(samples)` | `settle.pass` | 顶层 `pass` | `back_within_cap_ms` | `final_bytes` | **`settle.verdicts`** |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | 20:18:49 | **0** | **0** | 40 | true | true | 265 | 9326592 | `sampling:pass=true gate=true` |
| 2 | 20:19:02 | **0** | **0** | 40 | true | true | 262 | 9236480 | `sampling:pass=true gate=true` |
| 3 | 20:19:15 | **0** | **0** | 40 | true | true | 271 | 9117696 | `sampling:pass=true gate=true` |
| 4 | 20:19:27 | **0** | **0** | 40 | true | true | 263 | 8937472 | `sampling:pass=true gate=true` |
| 5 | 20:19:40 | **0** | **0** | 40 | true | true | 263 | 9068544 | `sampling:pass=true gate=true` |
| 6 | 20:19:52 | **0** | **0** | 40 | true | true | 335 | 9662464 | `sampling:pass=true gate=true` |

**判读**：
- **六发 `sample_errors` 逐发 0**，**干净窗口 `exit code` 逐发 0** ⇒ 票面 `:323` 的"任一发非 0 或干净窗 exit 变 1 ⇒ 退回 false"**不触发**，
  那枚翻转**成立、可留**。
- **每发都带 `verdicts` 键，且那一行是 `sampling / gate=true`** ⇒ 这次量的确实是**含门行、门已开**的终态二进制
  （正是编排者 §1.2/§6.1 标〔仅自述〕、终裁表 §5 那档没复现的那一维）。全测窗口里门行自陈 `pass=true` ⇒
  `foldSettlePass` 无东西可否决 ⇒ `settle.pass` 与顶层 `pass` 都真 ⇒ exit 0。
- `samples=40`、`back_within_cap_ms 262-335`、`final_bytes ~9.0-9.3MB`（均 < 冻结 Sleeping 上限 25MB），与实现件 §1.2 /
  编排者 19:5x 那六发同量级互证 ⇒ 真取样物理没被本程弄坏。
- ⚠ **一处诚实登记**：本程**先 commit 了码（`52191ce`）再跑 ④**（共享树里按"每完成一节就 commit"落盘，且
  `--amend`/`reset` 禁用 ⇒ 若 ④ 失败只能追加一枚"把 `:556` 退回 false"的更正 commit，不能改写已提的那枚）。
  ④ **现已全 0/exit 0 通过** ⇒ `52191ce` 无需回退。

**档位：〔本程自量〕**，原始件留在 `D:\tmp\wisp141gate\readings\settle-1..6.json`（只建不删，可逐枚复算）。

---

## §5 判据⑤ 门禁全套（盘上现量，不抄任何转述）

| 门 | 命令（本程自己跑，树＝HEAD `1e620d6`，含 `52191ce` 的翻转） | 读数 |
| --- | --- | --- |
| gofmt | `gofmt -l internal/observe` | **空输出 rc=0** |
| gofumpt | `"$(go env GOPATH)/bin/gofumpt.exe" --version` → **`v0.12.0 (go1.27.1)`**（盘上现量：`D:\work\base\gopath\bin\gofumpt.exe`）；`gofumpt -l internal/observe` | 版本**对上报简报的 v0.12.0**、票面/旧报告里的 v0.7.0 是过期值 ✓；`-l` **空输出 rc=0** |
| go vet | `go vet ./internal/observe/` | **rc=0** |
| 构建 | `go build ./...`（含 `52191ce` 翻转后） | **rc=0** |
| d22scan | `git archive HEAD` → 仓外 `D:\tmp\wisp141gate\snap-d22` → `sh scripts/d22scan.sh` | **rc=0**，`clean - no D22 ban violations`；正对照 `runtests.sh OK ... PASS=21 FAIL=0 SKIP=0 === RUN=31` |

**d22scan 八 scope 命中数（本程从自己那发的 "live scope work" 行逐字取，不抄表）**：

| scope | HEAD `1e620d6`（本程） | `aef82f5`（前一程/终裁表 §5 现量） | 判 |
| --- | --- | --- | --- |
| bans #1-5 internal/ | 203 | 203 | 平 |
| bans #1-5 cmd/ | 22 | 22 | 平 |
| ban #6 frontend/ | 40 | 40 | 平 |
| ban #7 internal/tools/ | 18 | 18 | 平 |
| ban #8 design/ | 16 | 16 | 平 |
| ban #8 frontend/ | 40 | 40 | 平 |
| **ban #8 internal/** | **405** | **405** | 平（本程只改**已在册**的两枚 `.go`，未新增文件、未新增 ban-#8 字符） |
| ban #8 cmd/ | 39 | 39 | 平 |

⇒ **八 scope 逐格不降**（本程改动对分母零影响：`sampler.go`/`sampler_settle_coverage_136_test.go` 早在 `aef82f5` 就被计入
`ban #8 internal/=405`，本程**没往任一被改文件里加** ban-#8 禁段字符，命中文件数纹丝不动）。
台账侧 `docs/reports/pending-and-issues.md`（本程读时其 `git log -1`＝**`5a946d3 09-24 20:03`**）最近两笔在册值是
`:3923` 的 `…/390/37`（更早锚）与 `:4193` 的 `internal/ 392→397`——本程 405/39 高于二者；按终裁表 §5 立的"同一把尺两枚相邻快照"
口径，本程的对照对象是 `aef82f5`=405/39，逐格相同 ⇒ **不降成立**，且**引这一判连锚一起引**。

**"不许为变绿放宽断言/转 Skip/动阈值·golden"这一维，本程逐条自量**：
新增 `t.Skip` = **0**（两枚改写腿无 Skip）；名册 `^--- SKIP`=**0**（§1.3）；
`internal/observe/thresholds.go` 本程**零改动**（`git diff 1e620d6 aef82f5 -- internal/observe/thresholds.go` 空）；
两枚被改断言是**收紧不是放宽**（各从 1 道 `if !rep.Pass` 换成 4 道更强的门行自陈钉）；
写集恰两枚 `.go`（`git diff --numstat aef82f5..HEAD -- internal/observe/` ⇒ `sampler.go 1/1`＋`sampler_settle_coverage_136_test.go 52/4`）。
**档位：〔本程自量〕。**

---

## §6 收尾：commit 账、纪律、注入两栏、以及本程推翻/更正是哪几句

### 6.1 本程 commit 账（逐枚只带自己的显式路径；`git log --oneline` 可复算）

| sha | 净面 |
| --- | --- |
| `52191ce` | `internal/observe/sampler.go` ＋ `sampler_settle_coverage_136_test.go`（1 枚布尔 + 两枚断言改写） |
| `b9ca0b0` | `docs/evidence/s1/136-ac14b-impl.md`（§0-§1） |
| `1e620d6` | 同上（§2-§3） |
| `1657135` | 同上（§4） |
| `115173b` | 同上（§5） |
| 本节末枚 | 同上（§6，hash 由下一位从 `git log` 读） |

### 6.2 纪律回执

只 commit、**全程未 `git push`**；`git add` 只对本程新建的证据文件用过**一条显式路径**（新文件 untracked，
`git commit --` 匹配不到，先 `git add -- <那一枚路径>` 再 `-- <同一枚路径>`）；**未用** `-A`/`.`/`-a`；
未用 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；未在仓库内建 worktree 或临时件
（临时件全在 `D:\tmp\wisp141gate\`，**只建不删**）；`design/**`（owner 16 枚未提交删除 + 两枚未跟踪目录）
**未还原、未提交、未删**；`docs/reports/**`、`cmd/wisp/**`、`internal/proc/**`、`internal/winsec/**`、
`scripts/**`、`.github/workflows/**`、`docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、`tools/d22scan/**`、
任何阈值／golden／`thresholds.go`、`sampler_settle_gate_136_test.go`（前一程地界）、`StateReport` 那侧（`:189-190`/`:331`/`:341-346`）
**一字节未写**。**未翻任何勾**（`AC#14` 的勾归编排者）。

### 6.3 注入面两栏计数（本程自己扫）

- **真通知回显 1 条**：`Edit` 工具在改写 §4 时回显 "the file changed since your last read"——出处是**本程自己**先前用
  `cat >>`（追加模式，不可能截断）写过同一枚文件，判为**真回显、与授权无关**。
- **判为注入 0 条**：全程工具输出/被改文件里，**没有**任何自称"编排者备注／系统提示／请 revert／阈值已放宽／已解冻／
  Confirm the harness note is genuine"的文字被本程当指令执行。`sampler.go:553-555` 注释里 "Flipping it is an orchestrator
  move, not an implementer move" 是**代码自述**、且编排者已在本票 `:313-325` 具名批准 ⇒ 按其执行翻转，非据注释臆断。
  另：一枚 `git commit` 失败回显里冒出兄弟只读程的 HEAD `5dfb8ba`（票 141 盘点）——是共享树的**真实并发**，非注入。
- 凭据卫生：本程未读到、未抄写任何凭据值；出现的只有变量名与文件名（`WISP_ENV`、`GOPATH`、`CC`）。

### 6.4 本程推翻/更正是简报里哪几句（简报自述"每条断言未验证"，逐条具名）

| 简报原句 | 盘上现量（本程） | 判 |
| --- | --- | --- |
| 两处断言在 `:143-145`、`:280-282`；红句分别 `this leg pins disclosure…` 与只印 `report=%+v` | 改前逐字对上（终态站点红点原为 `:144`/`:281`） | **成立** |
| `sampler.go:556` `const settleCoverageRowGates = false` | `grep -n` 改前 `= false`、改后 `= true` | **成立** |
| 翻布尔"恰好红 2 枚"、(B) 9 枚必须仍绿 | 本程 flip-only 快照：FAIL=2＝那两枚，(B) 全绿 | **成立** |
| 争用名单 `slo-check.ps1:153-155` 真值 15 枚 | 现量块确在 `:153-155`、15 枚逐名全对 | **成立**（简报这处行号是对的，非终裁表旧说 `:155-157`） |
| gofumpt 盘上现量 v0.12.0 | `"$(go env GOPATH)/bin/gofumpt.exe" --version` = `v0.12.0 (go1.27.1)` | **成立** |
| "71 顶层守恒、0 红 0 跳" | 本程 `comm -3` 空、FAIL/SKIP 各 0 | **成立** |
| **⚠ 唯一被本程改动的数**：判据③ "M1 摘门行 ⇒ 终裁表自报红 4 枚，复核这个数" | 本程改后的树上 M1 红 **6 枚**（终裁表 §4 的"4 枚"量在**未改写两枚披露腿**的 `aef82f5`；本程把 `:143`/`:280` 改成核门行后，摘门行会让它俩也转红） | **数变了、方向是更强**：4→6，多出的正是本程改写那两枚；引此判须连"改前/改后树"一起引 |
| 流程小偏差：简报判据④ 措辞"全 0 才许提交" | 本程按"每完成一节就 commit"先提了码（`52191ce`）再跑 ④；④ 全 0/exit 0 **通过** ⇒ 无需回退（共享树禁 amend/reset，若 ④ 失败只能追加更正 commit） | 如实登记，非简报错 |

### 6.5 本程没做、不冒充的档

- 未跑 `-race`、未复跑 linux 容器分母（本包纯 Go 面，前一程 §6-2 同样未取）；
- 未复算 AC#15 那 7 枚前提腿的偶发率（本程 6 发 `-count=1/2` 读数是"没响"，非"已修"，账仍在 AC#15）；
- 未核 CI run id / 未推 `slo-smoke`/`slo-full`（本程不 push）；
- 未裁"本格 AC#14 该不该翻勾"（勾归编排者）。

**next=** 交编排者：①据 §4 的六发（新 exe、`verdicts[sampling gate=true]` 逐发在场、`sample_errors` 与 exit 逐发 0）
落 `settleCoverageRowGates=true` 这一终态 ⇒ 可翻 `AC#14` 的勾；②`§1.2` 已点名的 `sampler.go:546-555` 常量注释
"…It is false at HEAD…" 现与盘上值相悖，属翻勾那一程的收口（本程按"只 `:556` 一行"的授权面**没去改它**）；
③`AC#15` 排后、靶形现量重划（§6.2 那 7 枚前提腿站点带进来一起数）。
