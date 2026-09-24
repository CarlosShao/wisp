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
