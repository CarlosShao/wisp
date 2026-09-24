# 136 — AC#14 实现件（settle 报告的自陈门行）**impl**

> **这份文件是什么**：票 136 `AC#14`（`.scratch/wisp/issues/136-...md:251-271` 那一格 ＋
> `:273-292` 那枚编排者 09-24 18:4x 的**具名解冻＋四裁定**）的实现件。实现程
> `worker-ticket136-ac14`（本程），非验收方，**本文件不翻任何勾、不构成任何一格的结案依据**。
> **授权**：解冻只按票面 `:276-277` 那两段（`internal/observe/sampler.go:431-456` ＋ `:462-521`），
> 边界按 `:278-279`（`StateReport` 侧 `:189-190`/`:331`/`:341-346` 冻结、`buildVerdicts` 不许两用）。
> **本文件只引用变量名与文件名，不含任何凭据值**（本程也没有读到过任何凭据）。

---

## §0 锚点、口径、取数时刻

| 项 | 读数 | 取数时刻（本机） |
| --- | --- | --- |
| 本程自量锚点 | `git rev-parse --short HEAD` ⇒ **`ca2c55e`**（分支 `dev`） | `2026-09-24 18:50:04 +0800` |
| 工具链 | `go version go1.27.1 windows/amd64`；`gofumpt -version` ⇒ **`v0.12.0 (go1.27.1)`**（盘上现量，未照抄票面旧值 v0.7.0） | 18:52 |
| 工作树里别人的东西 | `git status --porcelain` 整树只有 owner 的 `design/**` 16 枚未提交删除 ＋ `design/doubao/`、`design/old/` 两枚未跟踪目录 ⇒ **一枚没碰、没还原、没代提交** | 18:50 |
| 临时件落点 | `D:\tmp\wisp136ac14\`（`bin/wisp.exe` ＋ 三枚 DLL · `readings/` · `snap/` 纯净快照 · 两份 `-v` 日志），**仓库目录内没有新建任何东西**；`build/wisp.exe`（mtime 09-21 15:06，是预检程 §1.2 引过的那枚参照件）**未被覆盖** | 18:53-18:56 |

### 0.1 基线四数（改动**之前**，`go test -count=2 -v ./internal/observe/` 跑两遍）

| 口径 | run1 | run2 |
| --- | --- | --- |
| `^=== RUN`（顶层） | 130 | 130 |
| `^--- PASS` | 130 | 130 |
| `^--- FAIL` | 0 | 0 |
| `^--- SKIP` | 0 | 0 |
| 去重后的顶层名数 | **65** | **65** |
| `grep -ci panic`（整文件命中行数） | 8 | 8 |
| `^panic:`（真 panic） | 0 | 0 |

⇒ **四数口径**：本包的 "-count=2" 把顶层名册翻倍，**65 枚 × 2 = 130**，报数一律带这个乘子。
**`grep -ci panic` 的 8 行命中全是用例名里带 `Panic` 的两枚**（`TestPanicInWorkerSurvivesAndCancelsRoot`、
`TestFakeTreeEmptyScriptFailsClosedAndNotPanics`）各 ×2（`=== RUN` ＋ `--- PASS`）×2（`-count=2`）⇒ **真 panic=0**，
按"名册差集＋panic 逐名"这一条口径写死在这里，免得下一位把 8 读成 8 发崩溃。
**名册差集**：`grep '^--- ' \| awk '{print $2,$3}' \| sort` 两遍 ⇒ `comm -3` **输出 0 行**。
凭据文件 `D:\tmp\wisp136ac14\base-v-1.txt`／`base-v-2.txt`／`b1.names`／`b2.names`。

### 0.2 门禁基线（改动之前）

| 门 | 读数 |
| --- | --- |
| `gofmt -l internal/observe cmd/wisp` | 空输出（rc=0） |
| `"$(go env GOPATH)/bin/gofumpt.exe" -l internal/observe` | 空输出 |
| `go vet ./internal/observe/` | rc=0 |
| `sh scripts/d22scan.sh` @ `git archive HEAD` 纯净快照（`D:\tmp\wisp136ac14\snap`，树外） | **rc=0**；八 scope 命中数 `203 / 22 / 40 / 18 / 16 / 40 / 404 / 39`（`bans #1-5 internal/`、`#1-5 cmd/`、`#6 frontend/`、`#7 internal/tools/`、`#8 design/`、`#8 frontend/`、`#8 internal/`、`#8 cmd/`）；正对照 `runtests.sh` 21/21/0/0。日志 `D:\tmp\wisp136ac14\d22-base.txt` |

台账口径出处：`docs/reports/pending-and-issues.md:3922`（"台账八 scope … 逐格不降"）、`:4192`（"两处各有出处 ⇒ **一枚不降**，不是'没降是因为没人看'"）、
`:4392`（`ban #8 internal/` 401→402 归因到自己那枚 `_test.go`）。⇒ 引用命中数**必须连归因一起引**，本文照此办（§3）。

---

## §1 `Gate` 那一支：先量的读数，逐次列（票面 `:281-285` 要求的前置）

**裁判问题**：门行做成 `Gate: true` 会把 `wisp slo -settle` 在"部分未测"的真窗口上从 exit 0 变 exit 1，
并经 `scripts/slo-check.ps1:381` → `:396` 打到 CI 的 `slo-smoke`（托管）与 `slo-full`（**这台 6C12T 本机**）。
票面 `:283` 要的实现程第一发＝编队安静窗口内 ≥5 次 `wisp slo -settle`、逐次记 `sample_errors`。

### 1.1 取数前置：编队现量

| 检查 | 读数 |
| --- | --- |
| `scripts/slo-check.ps1:153-155` 名单**盘上现量**（逐字取自该脚本，15 枚，**不是派单简报里的 8 枚**：多 `link.exe`/`g++.exe`/`cc1.exe`/`cc1plus.exe`/`as.exe`/`ld.exe`/`wisp-cli.exe` ⇒ 用更全的那份，只可能更严） | 每发读数**前后各量一次** ⇒ **12 次全空**（`roster_before=` / `roster_after=` 皆无值） |
| `gh run list --limit 6` | 在飞的只有 `35989390166`（`ca2c55e` 那次推送，18:48:47Z 起） |
| 同一发的逐 job | `slo-full`（self-hosted，**就是这台机**）= **completed/success**；`slo-smoke` = success；`lint` = **failure（先于本程，见 1.4）**；`test-core`/`lint-frontend` = success；`test-windows` 在飞但 `runs-on: windows-latest`（`ci.yml:335`，托管，不抢本机） ⇒ **取数窗口内 `slo-full` 不在飞** |

### 1.2 六发读数（≥5，逐次）

二进制：本程 `go build -trimpath -ldflags "-X …buildinfo.Commit=ca2c55e" -o D:\tmp\wisp136ac14\bin\wisp.exe ./cmd/wisp`
（`CGO_ENABLED=1` ＋ `E:/work/base/msys64/mingw64/bin/gcc.exe`，rc=0），**三枚 sherpa/onnx DLL 与 exe 同目录**
（`third_party/sherpa-onnx/{onnxruntime,sherpa-onnx-c-api,sherpa-onnx-cxx-api}.dll`，简报 `⚠` 那条加载期依赖照办了）。
命令逐发：`WISP_ENV=test`（＝`ci.yml:540` 那枚 job env）＋ `wisp.exe slo -settle -out readings/settle-N.json`。
取数脚本 `D:\tmp\wisp136ac14\readings.ps1`（脚本内**不调用任何 go 工具**），日志 `D:\tmp\wisp136ac14\readings/readings.txt`。

| # | 起时 | exit | `sample_errors` | `len(samples)` | `pass` | `back_within_cap_ms` | `final_bytes` | 名单前后 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | 18:54:59 | 0 | **0** | 40 | true | 267 | 10133504 | 空 / 空 |
| 2 | 18:55:13 | 0 | **0** | 40 | true | 263 | 9252864 | 空 / 空 |
| 3 | 18:55:26 | 0 | **0** | 40 | true | 264 | 9064448 | 空 / 空 |
| 4 | 18:55:39 | 0 | **0** | 40 | true | 265 | 9203712 | 空 / 空 |
| 5 | 18:55:52 | 0 | **0** | 40 | true | 265 | 9248768 | 空 / 空 |
| 6 | 18:56:05 | 0 | **0** | 40 | true | 269 | 10350592 | 空 / 空 |

⇒ **六发 `sample_errors` 逐发全 0**（`sample_errors_total=0`；`retries=0` 六次，即没有一发因争用被重取）。
按票面 `:283` 那一支，这枚读数的含义是"**全 0 ⇒ 按 `Gate: true` 落地，这一形今天不会误伤存量合法窗口**"。

### 1.3 为什么本程仍然**没有**把它落成 `Gate: true`（一处前提被盘上推翻，具名停手报回）

票面 `:287` 与派单简报都把"AC#14 要推翻的既有断言"记成**一处**：
`sampler_settle_coverage_136_test.go:208-210` 那一枚 `if !rep.Pass`。预检程 §3.3 末也只点到那一处。
**盘上现量是三处**，同一条"disclosure not verdict"形状，同一枚文件：

| 位置 | 逐字断言 | 它现在要求什么 |
| --- | --- | --- |
| `:143-145`（`TestCheckSettleSingleTrustworthyReadReportsItsLoss`） | `if !rep.Pass {` ⇒ `t.Fatalf("this leg pins disclosure, not the verdict; a covered-enough window must still pass: %+v", rep)` | 10 枚读数里只 1 枚可信（丢了 9 枚）那一形**必须仍然 pass** |
| `:208-210`（`TestCheckSettleHalfTheReadsFailedReportsItsLoss`） | `if !rep.Pass {` ⇒ `t.Fatalf("disclosure leg, not a verdict leg: this window still passes, report=%+v", rep)` | 丢一半那一形**必须仍然 pass** ← **唯一被预先授权改的那一处** |
| `:252-254`（`TestCheckSettleZeroFootprintDropsAreCountedToo`） | `if !rep.Pass {` ⇒ `t.Fatalf("report=%+v", rep)` | 1 枚可信 / 2 枚零足迹（丢了 9 枚）那一形**必须仍然 pass** |

**任何"会抓住丢一半那一形"的门行判据必然同时抓住另两形**（那两窗口的覆盖率严格更差：kept=1 / lost≥9，
对照 kept≈lost≈5），除非把判据写成"只红在恰好一半"这种非单调的任意阈值——那是**新造阈值**，
`AGENTS.md §1.1` 与票面 `:279` 都禁。⇒ `Gate: true` ＋折进总 `Pass` 的**最小**代价是**同改三处断言**，
而本程的授权只覆盖其中一处。**扩授权不是本程能做的动作**（`AGENTS.md §0` 第 1 句"未定义即停"、派单"超出即停手报回，不许自己扩"），
故本程落地形状＝票面 `:284` 给的那一档**记录行**（`gate=false`、不折进总 `Pass`），
把"三处 vs 一处"这枚盘上事实连同 §1.2 的全 0 读数一起交回编排者裁：
**批准另外两处之后，`sampler.go` 里要改的是**一枚具名常量的一个布尔字面量**（`settleCoverageRowGates`），
门行的自陈、红句里的 `sample_errors`、折叠处形状**都已就位**，不需要再动别的。
§4 那两发快照变异给出"批准之后会红哪三枚"的**读数**（不是推理）。

本程**没有**把 `AC#14` 判成完成：按票面 `:262-263` 判据① 的原话与预检程 §3.3 末的读法，
记录行那一支会让"不许只靠 `pass=false` 一个总布尔代答"这一维**落空**——本文件如实把这一格标为**未结**。

### 1.4 顺带量到的、不属本格的账（登记，不动手）

- `lint` 在锚点 `ca2c55e` 上就是 **failure**（`gh run view 35989390166 --json jobs`，18:5x 现量）⇒ 本程交件时 `lint`
  的色**不是本程造的**；本程自己的 gofmt/gofumpt/vet/d22scan 逐门读数在 §3。
- 简报里"进程名单 8 枚"与盘上 15 枚不符（§1.1）⇒ 按更全那份量的，只可能更严，不影响"全 0"这一结论。
- **锚点在程中漂过一枚，且漂来的那枚与本格直接相关**：本程取基线时 `HEAD=ca2c55e`；
  18:5x 编队落了 `2e6d171 docs(evidence/136): SLO 出线契约普查 r1`（＝票面 `:325(b)` 要的那一查），
  本程的码与探针落在 `aef82f5`。普查的判定与本程无关但方向一致：
  `136-slo-output-contract-census-r1.md:80` 逐字 "**普查结论：排除 `SPEC-02 §3` 之后，没有任何一份文档在 `SPEC-12 §4.1` 意义上管 SLO 报告的出线形状**"
  ⇒ 给 `SettleReport` 加 `verdicts` **不触**任何文赋契约的停手线（该文件 `:184` 亦明确"要不要做成 `Gate:true` 仍按真取样读数走人工判断"，与本程 §1 的处置同向；
  它 `:167-174` 把"AC#12 用例与 AC#14 方向相反"记成**一处**并被预授权——这一处正是 §1.3 里被盘上推翻的那枚前提，**是三处不是**一处）。

---

## §2 落点：门行落在哪几段（形状与逐段出处）

本程**只**动了三枚文件（`git show --name-only aef82f5` 原样在 §2.4），逐段落点如下（行号为**改动后** `aef82f5` 上的现量）：

### 2.1 `internal/observe/sampler.go`（解冻面内，7 处改动）

| # | 段（改后行号） | 是什么 | 在票面 `:277` 划的哪一段里 |
| --- | --- | --- | --- |
| 1 | `:455-466` `SettleReport.Verdicts []Verdict json:"verdicts"` ＋注释 | 新字段：出线带自陈行 | `:431-456` 声明块内 ✓ |
| 2 | `:467-469` `Pass bool json:"pass"` 的注释改写（字段本体逐字未动） | 写明"总布尔之外还要被失败 gate 行否决" | 同一声明块 ✓ |
| 3 | `:529` `rep.Verdicts = buildSettleVerdicts(*rep)` | 行的构造点（在样本/丢读计数都定形之后、折叠之前） | `:462-521` 体内（现 `:475-547`）✓ |
| 4 | `:534-539` 折叠注释 ＋ `:540` `rep.Pass = foldSettlePass(memOK && backInTime && releaseOK, rep.Verdicts)` | 折叠处（原 `:520` 那三布尔的与，语义逐字保留，只多一道 gate 行否决） | 同上 ✓ |
| 5 | `:549-552` `const settleCoverageMetric = "sampling"` | 行名（与 `StateReport` 那行同名，同问题） | **文件末尾新增**，见 §2.3 |
| 6 | `:554-566` `const settleCoverageRowGates = false` | **唯一一枚待裁开关**（§1.3） | 同上 |
| 7 | `:568-600` `buildSettleVerdicts` ＋ `:602-612` `foldSettlePass` | 另起一枚构造函数（票面 `:279` 明令）＋折叠规则（可单测） | 同上 |

`Measured` 逐字沿用 `StateReport` 那行的格式串（`sampler.go:336` `fmt.Sprintf("%d valid / %d errors", …)`），
红句（`Note`）把 `sample_errors=` 按**字段名**印出来并带上最后一次丢读原因，绿句也印计数：
判据① 的"由那一行自己说出不 pass"与 `:285` 的"红句必须点名 `sample_errors`"同一条满足。

### 2.2 冻结面：逐字未动（`git diff ca2c55e..aef82f5 -- internal/observe/sampler.go` 全量在上，可复算）

| 冻结坐标（票面 `:278`） | 现量 |
| --- | --- |
| `:189-190` `Verdicts []Verdict` ＋ `Pass bool // all Gate verdicts pass`（`StateReport` 侧） | **不在 diff 的任何 hunk 里**（`git diff -U0` 只有 `@@ -452,7`、`@@ -514,9` 两枚 hunk，起点均在 `SettleReport` 之后） |
| `:331` `rep.Verdicts = buildVerdicts(st, *rep)` | 未动（同上） |
| `:341-346` `Gate && !Pass` 的折叠 | 未动 |
| `internal/observe/thresholds.go`（`buildVerdicts` 的家，`:84`） | `git diff --stat` 里**没有这枚文件**，一字节未动 |
| `cmd/wisp/**`、`internal/proc/**`、`internal/winsec/**`、`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`internal/risk/**`、`tools/d22scan/**`、`scripts/slo-check.ps1`、`.github/workflows/**`、任何阈值／golden | 全部未动（§2.4 的三枚路径就是本程全部写集） |
| `buildVerdicts` 是否被改成两用 | 没有；`TestStateReportVerdictBuilderStaysSinglePurpose` 在**编译期**钉住它的 `(SLOState, StateReport)` 签名，并在断言里挡住"sampling 行/`sample_errors` 字样漏进表行" |

### 2.3 一处坐标偏离，具名登记（不悄改）

票面 `:277` 划的两段是 `:431-456` 与 `:462-521`，而 `:522` 就是 `CheckSettle` 的收尾花括号、`:523` 起是文件边界。
**一枚新的顶层函数在物理上放不进"`:462-521` 体内"**——放不进去又不许改 `buildVerdicts`（`:279`）只剩两条路：
闭包（把构造函数藏进 `CheckSettle`，测试就点不到它，变异也打不中它）或**紧邻 `CheckSettle` 之后追加到文件末尾**。
本程取后者：新增码在 `:549-612`，**超出票面字面坐标 27 行、但不出本格的射程**，
且它正是 `:279` 那句"要给 settle 造行就**另起一枚构造函数**"所要求的东西。登记在此由编排者裁是不是我的错读。

### 2.4 本程写集（`git show --name-only`，逐枚 commit）

```
$ git log --oneline -1        # 0ee68e1 evidence(136 AC#14 §0-§1) …   -> docs/evidence/s1/136-ac14-impl.md
$ git log --oneline -1        # aef82f5 feat(136 AC#14) …            -> 三枚路径，见下
aef82f5 feat(136 AC#14): settle 报告自陈门行 + 六枚探针（记录行形状，gate 常量留一枚待裁的开关）
    internal/observe/sampler.go
    internal/observe/sampler_settle_gate_136_test.go
    internal/observe/sampler_settle_coverage_136_test.go
```

**预先授权的那处既有断言改动（按 `:287` 要求逐字点名）**：
`internal/observe/sampler_settle_coverage_136_test.go` 原 `:208-210` 三行——
逐字 `if !rep.Pass {` / `t.Fatalf("disclosure leg, not a verdict leg: this window still passes, report=%+v", rep)` / `}`——
**被本程改掉**。理由＝票面 `AC#14` 判据①／②（票面 `:262-264`）：那一形逐字要求"丢一半读数仍然 pass"，
与 AC#14 要的方向正相反。改后那段钉的是：门行必须存在、必须**自己**说 not-pass、红句必须印出 `sample_errors`、
且"失败 gate 行与 `pass=true` 并存"这一形必须永不出现在报告里。
⚠ **本程没有动那枚文件的头注释**（`:25-27` "a partially covered window still passes today (changing that verdict is not this cell's job)"）
——它对 `:143-145`/`:252-254` 那两形**今天仍是真的**（门行是记录行），对 `:208-210` 已被改写；
把它改准属 AC#15 的射程（票面 `:288` 明令"AC#15 到时要现量重划"），本程不越界，登记在此。

---

## §3 门禁读数（改动之后，全部本程自己跑）

取数时刻 `2026-09-24 19:0x +08`，树＝`aef82f5`（锚点链 `ca2c55e` → `2e6d171` → `0ee68e1` → `aef82f5`）。

| 门 | 命令 | 读数 |
| --- | --- | --- |
| 四数 run1 | `go test -count=2 -v ./internal/observe/` | **RUN=142 / PASS=142 / FAIL=0 / SKIP=0** ⇒ 顶层 **71 枚 × 2**（基线 65×2=130，净增 6＝本程探针）；`grep -ci panic`=8，`^panic:`=**0** |
| 四数 run2 | 同上，第二遍 | 逐格相同：**142/142/0/0**，panic 8/0 |
| 名册差集（两遍之间） | `grep '^--- ' \| awk '{print $2,$3}' \| sort` → `comm -3 p1 p2` | **0 行** ⇒ 两遍名册逐名相同 |
| 名册差集（基线→改动后） | `comm -23 b1 p1`（谁消失了）／`comm -13 b1 p1`（谁新增） | 消失 **0 枚**（**没有一名转 SKIP、没有一名不见**）；新增恰为本程 6 枚：`TestSettleCoverageRowExistsAndPassesWhenFullyMeasured`、`TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed`、`TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured`、`TestFoldSettlePassOnlyGateRowsVeto`、`TestSettleReportPassNeverContradictsItsGateRows`、`TestStateReportVerdictBuilderStaysSinglePurpose` |
| gofmt | `gofmt -l internal/observe` | 空输出 |
| gofumpt | `"$(go env GOPATH)/bin/gofumpt.exe" -l internal/observe`（盘上现量 **v0.12.0 (go1.27.1)**） | 空输出 |
| vet | `go vet ./internal/observe/` | rc=0 |
| 构建 | `go build ./...` | rc=0 |
| d22scan | `git archive HEAD` → `D:\tmp\wisp136ac14\snap-post`（仓外）→ `sh scripts/d22scan.sh` | **rc=0**；八 scope `203 / 22 / 40 / 18 / 16 / 40 / 405 / 39`，正对照 `runtests.sh` 21/21/0/0 |

**八 scope 的归因（§0.2 台账口径：引命中数必须连归因）**：与基线 `203/22/40/18/16/40/404/39` 逐格比，
唯一变化是 `ban #8 internal/` **404 → 405**，＋1 枚＝本程新增的 `internal/observe/sampler_settle_gate_136_test.go`；
其余七格**一枚不降**，`bans #1-5 internal/`（只数生产码）维持 203——本程那 1 枚新文件是 `_test.go`，不进生产码分母，
这与台账 `:4392`（`ban #8 internal/` 401→402 归因到自己那枚 `_test.go`）是同一形制。

**Emoji / 注释卫生**：新码与新测试文件的注释逐字 ASCII（`AGENTS.md §1.2` 禁的 `U+2190–U+2BFF` 段含 `⇒`、`→`，
本程在 Go 侧一律写 `->` 或改述）；d22scan ban #8 覆盖 `internal/` 405 枚 Go 文件（注释与 `_test.go` 全算）rc=0，
这是"零 emoji"那一维的**外部**读数，不靠本程自述。

**跨包后果自查**（本程改了 `cmd/wisp` 消费的那枚结构体，虽不写它的码）：`grep -rln settle cmd/wisp/*_test.go` ⇒ **零命中**
⇒ `cmd/wisp` 没有任何用例碰 settle 出线，`slo_windows.go:623` 那枚合成报告是具名字面量、加字段不破编译（`go build ./...` rc=0 已证）。
⚠ 另外量到一件不属本格的事，如实登记：宿主上 `go test ./cmd/wisp/` 直接跑会以 **`exit status 0xc0000135`（STATUS_DLL_NOT_FOUND）**
收场（本程现量），把三枚 sherpa/onnx DLL 放进 `cmd/wisp` 目录后**在改动前的快照 `snap/`（ca2c55e）上仍然 FAIL（91.19s）**
⇒ 那枚红**先于本程**、不是本程造的；改动后的逐名对照见 §5。

---

## §5 跨包后果的逐名对照（改动前 vs 改动后，同一台机、同一天、各一发 `-count=1 -v`）

两发都在**仓外快照**上跑（`snap/`＝`ca2c55e`、`snap-post/`＝`aef82f5`），
并把三枚 sherpa/onnx DLL 放进各自 `cmd/wisp/` 目录（不放就是上面那发 `0xc0000135`）。

| 树 | RUN | PASS | FAIL | SKIP |
| --- | --- | --- | --- | --- |
| `ca2c55e`（本程改动**前**）`go test -count=1 -v ./cmd/wisp/` | 95 | 46 | **8** | 0 |
| `aef82f5`（本程改动**后**）同一条命令 | 95 | 46 | **8** | 0 |

**逐名对照**：`comm -13 pre-fail post-fail`（只在新树红的＝本程造的）⇒ **空**；
`comm -23`（只在旧树红的）⇒ **空**；整名册 `comm -3` ⇒ **空**。
⇒ 本程给 `SettleReport` 加 `verdicts` 这一维**对 `cmd/wisp` 的 95 枚读数零影响**（可复算：`grep -rln settle cmd/wisp/*_test.go` ⇒ 零命中）。

那 8 枚既有的红（`ca2c55e` 上就在，本程不背也不修，登记给出站 AC#10 的实现程）：
`TestAC2RealProcessRefusesOnEveryLegWithoutAppData128`、`TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord`、
`TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot`、`TestAC1ResidentLegInstallsItsLogListenerOnDisk`、
`TestAC1ResidentLegOutlivesItsOwnLogFailure`、`TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`、
`TestSecretArgvCarriesNoSecret`、`TestSecretRealBinaryRefusesValueFlag`
（票 128 的 resident/earlylog 族 ＋ winsec/secret 族，落在本程禁改的 `internal/proc/**`、`internal/winsec/**` 那一侧）。

## §6 本程**没有**做的、与残留风险（一格都不许被当结论地基）

1. **没做 `Gate: true`**，因此 AC#14 判据① 的"不许只靠总布尔代答"这一维**只落了一半**：门行会自己说不 pass，
   但它今天不否决总布尔。§1.3/§4.2 给了那一步的全部前置读数与代价（1 枚布尔字面量 ＋ 2 枚断言）。
   **本格按票面口径判"未结"，本程没有翻任何勾。**
2. **容器/linux 分母没复跑**：只做了 `GOOS=linux CGO_ENABLED=0 go build ./internal/observe/` 与 `go vet`（各 rc=0），
   `scripts/portable-tests.sh` 那 13→14 枚 `_test.go` 的真实容器读数本程没取（预检程 §7 第 4 条同样没取，账在它名下也在这名下）。
3. **本程新探针的偶发暴露**（AC#15 的族，本程只量不自修）：`TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed`
   的前提腿是 `kept>=1 && lost>=1`，窗口 200ms/10ms（约 20 发），**要只响 1 发 ticker 才会不成立**；
   AC#12 那两枚既有腿用的是 100ms/10ms 且要求 `reads>=4`，比本程的前提更紧、已被量到 240 发里红 1 发（票面 `:266` 口径）。
   ⇒ 本程的前提在**同一把尺下比已知会红的那枚更松**，但**这不等于它不会红**：六发变异与两遍 `-count=2`（142 发）**没有一发**碰到它，
   样本量小、不外推。AC#15 重划靶子时要把它一起数进去（本程六枚新探针的名字都在 §3 的新增清单里，逐名可取）。
4. **机时卫生的回执**：本程 §1 那六发取样是在 18:53 完成 `go build` 之后、且每发前后各量一次进程名单全空的前提下取的；
   反向账本程也如实说——18:5x 起本程在宿主上多次 `go build`/`go test`（含两发 90s 级 `cmd/wisp` 整包），
   **同一时间窗内别的程若在这台机上取真数，那些数按 `slo-check.ps1:153` 的口径应判"无效样"**。本程没有在那段时间取任何真数（除 §1 那六发，逐发名单为空）。
5. **`slo_windows.go:623` 那枚合成报告今天多了 `verdicts: null`**：本程不碰 `cmd/wisp/**`（AC#10 地界），
   但出线形状因本程而变这一点必须让 AC#10 知道——它追加问② 要区的"根本没测 vs 测满零丢"两形，
   在 settle 侧现在**有了行可依**（`§4` 的 `buildSettleVerdicts(SettleReport{})` 那一形就是"0 valid / 0 errors"），
   而那枚合成报告要落到同一判据得在 `cmd/wisp` 侧补，本程不越界。

## §7 纪律回执

- 改动文件全集＝`internal/observe/sampler.go`、`internal/observe/sampler_settle_gate_136_test.go`（新建）、
  `internal/observe/sampler_settle_coverage_136_test.go`（**只**改预授权那一处）、`docs/evidence/s1/136-ac14-impl.md`（本文件）、
  `.scratch/wisp/issues/136-...md`（**只追加**一条 progress log，不改任何人的历史行）。逐枚 commit 的 `git show --name-only` 原样在下面。
- 未 push；未用 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；`git add` 只用显式路径；未在仓库内建 worktree/临时件；
  `design/**` 与两枚未跟踪目录没碰、没还原、没代提交；`build/wisp.exe` 没被覆盖（本程 exe 落在 `D:\tmp\wisp136ac14\bin\`）。
- 临时件只建不删。
- 注入面两栏计数：**真通知回显 0 条**（本程未收到任何自称"系统提示/编排者备注/文件已被修改/请 revert/阈值已放宽/已解冻"的工具输出）；
  **判为注入 0 条**。全程唯一带"指令"味道的输入是派单本身与票面 `>` 块（都是编排者署名、可复算的授权面）。
- 凭据卫生：本程未读到也未抄写任何凭据值；出现的都是变量名与文件名（`WISP_ENV`、`MINGW64_ROOT`、`GITHUB_WORKSPACE`、`RUNNER_TEMP`）。


---

## §4 变异自证（六发，全部落在仓外快照 `D:\tmp\wisp136ac14\mut\*`，工作树未参与）

每一发都按本仓口径**先证落地再读数**：`grep -n` 出被改那一行（下表"落地凭据"列）＋ `go build ./...` rc=0，然后才读红名。
`-count=1 -v ./internal/observe/`，全 71 枚名册都在，**没有一发用 SKIP 换色**。

| 发 | 改了什么（落地凭据＝改后 grep 出来的那一行） | build | 红名（逐名） | 计数 |
| --- | --- | --- | --- | --- |
| **M1 摘掉门行**（判据②后半"把门行摘掉 ⇒ 钉子必须转红"） | 删掉 `rep.Verdicts = buildSettleVerdicts(*rep)` 整行；改后 `grep -n buildSettleVerdicts sampler.go` ⇒ 只剩 `:563` 的函数定义与三处注释，**构造点为零** | rc=0 | `TestCheckSettleHalfTheReadsFailedReportsItsLoss`（`sampler_settle_coverage_136_test.go:226` "must carry its own sampling verdict row, got verdicts=[]"）、`TestSettleCoverageRowExistsAndPassesWhenFullyMeasured`（`_gate_:162`）、`TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed`（`_gate_:207`）、`TestSettleCoverageRowSeparatesUnmeasuredFromFullyMeasured`（`_gate_:245`） | 71/67/**4**/0 |
| **M2 让门行不再关心丢读**（判据②"一半读数报错 ⇒ 该门行红"） | `covered := rep.SampleErrors == 0 && len(rep.Samples) > 0` → `covered := len(rep.Samples) > 0`（`grep -n "covered :=" ⇒ :565`） | rc=0 | `TestCheckSettleHalfTheReadsFailedReportsItsLoss`（"this window dropped 5 of 10 reads and its own row says it measured enough"）、`TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed`（"a window that dropped 10 of 20 reads must be judged by its own row as not fully measured"） | 71/69/**2**/0 |
| **S1 只翻 `settleCoverageRowGates` 一枚布尔**（＝§1.3 交回的那一发，打在**已含 :208-210 改写**的树上） | `const settleCoverageRowGates = true`（`grep ⇒ :556`），其余一字节未动 | rc=0 | `TestCheckSettleSingleTrustworthyReadReportsItsLoss`（`sampler_settle_coverage_136_test.go:144` "this leg pins disclosure, not the verdict; a covered-enough window must still pass"）、`TestCheckSettleZeroFootprintDropsAreCountedToo`（`:281` "report=&{TargetState:Sleeping …}"） | 71/69/**2**/0 |
| **M3a 单点回退折叠调用点（改动 #4），gate 仍 false** | `rep.Pass = foldSettlePass(…)` → `rep.Pass = memOK && backInTime && releaseOK`（`grep ⇒ :537`） | rc=0 | **无（71/71 全绿）** ⇒ 见下面"哪一处是装饰" | 71/71/0/0 |
| **M3b 同一处回退、gate 翻 true** | 同上两行同时改（`:537` + `:556=true`） | rc=0 | `TestCheckSettleHalfTheReadsFailedReportsItsLoss`（`:236` "a failing gate row and pass=true at once: the fold rule is not applied"）、`TestSettleCoverageRowSaysNotPassWhenHalfTheReadsFailed`（`_gate_:226`）、`TestSettleReportPassNeverContradictsItsGateRows`（`_gate_:325` "half failed: pass=true alongside a failing gate row"） | 71/68/**3**/0 |
| 正向对照（尺是活的） | 未改的 `snap-post` | rc=0 | 无 | 71/71/0/0 |

### 4.1 单点回退的答案：**哪一处是装饰**

七处改动里，**改动 #4（`CheckSettle` 末尾那行折叠调用点）在 HEAD 的落地形状下就是装饰**：M3a 撤掉它，
71 枚**一枚都不变不响**——因为记录行（`gate=false`）按定义不否决总布尔。本程不拿"它有注释"糊这一条，给三样东西：

1. **明说**：#4 今天无测试可感知，它是为 `gate=true` 那一步预备的形状，**不是**已经生效的行为；
2. **让可感知性先就位**：折叠规则本身被拆成 `foldSettlePass` 并直接被 `TestFoldSettlePassOnlyGateRowsVeto`
   的 7 个用例钉住（含"失败 gate 行必须否决 / 失败记录行不得否决"两形），**这条单测在 M3a 下仍绿**
   ⇒ 它钉的是规则，不是调用点，这正是"装饰"的证据形状；
3. **给出它何时不再是装饰**：M3b＝同一处回退 ＋ `gate=true` ⇒ 立刻红 3 枚并逐名可查。
   ⇒ 翻那枚布尔之后，#4 从装饰变承重，**且变承重的瞬间就有钉子看着**，不需要到时再补仪器。

其余几处各有独立感知者：#1/#3/#5/#7 由 M1（红 4 枚）与 M2（红 2 枚）覆盖；#6 是 §1.3 那枚待裁开关本身，
由 S1 量出代价＝**2 枚**命名用例（不是三枚，因为 `:208-210` 那一处已按预授权改成了 gate 无关的形状）；
#2 是注释（按定义装饰，不参与任何读数，故不与 #4 混记）。

### 4.2 §1.3 的那一发在此改写（读数推翻了自己前面的一句话）

§1.3 写"最小代价是同改三处断言"——那是**在改动之前**数的既有断言。S1 打在改动之后，
实测代价是**另外 2 枚**：`TestCheckSettleSingleTrustworthyReadReportsItsLoss`（`:143-145`）与
`TestCheckSettleZeroFootprintDropsAreCountedToo`（`:252-254`）。第 3 枚（`:208-210`）已被本程按预授权改成
"行必须存在＋行必须自己说不 pass＋不得出现失败 gate 行与 pass=true 并存"，那一形**与 gate 取值无关**，
所以翻布尔后它照旧绿——**这不是漏，是设计**：授权范围内的改写刻意写成两种取值都成立，
编排者批不批 `gate=true` 都不需要再回来改这枚文件。⇒ **交回的那一发读数精确化为**：
批准之后要动的只有那 2 枚断言（改法＝本程已用的同款三段），生产码要动的只有 `:556` 那一枚布尔字面量。


